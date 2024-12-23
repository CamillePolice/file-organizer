package services

import (
	"file-organizer/src/models"
	"fmt"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type OrganizerService struct {
	fileService   *FileService
	folderService *FolderService
	workerCount   int
}

type FileTask struct {
	file models.File
	path string
}

type Stats struct {
	succeeded int
	failed    int
	mu        sync.Mutex
}

func NewOrganizerService(fileService *FileService, folderService *FolderService) *OrganizerService {
	return &OrganizerService{
		fileService:   fileService,
		folderService: folderService,
		workerCount:   5,
	}
}

func (org *OrganizerService) OrganizeFiles(files []models.File, path string) error {
	fileTypes := org.fileService.GetFileTypes()
	tasks := make(chan FileTask, len(files))
	stats := &Stats{}
	var wg sync.WaitGroup

	// Start workers
	for i := 0; i < org.workerCount; i++ {
		wg.Add(1)
		go org.fileWorker(tasks, fileTypes, &wg, stats)
	}

	// Send tasks
	for _, file := range files {
		tasks <- FileTask{file: file, path: path}
	}
	close(tasks)

	wg.Wait()
	log.Printf("Organization complete. Succeeded: %d, Failed: %d\n", stats.succeeded, stats.failed)
	return nil
}

func (org *OrganizerService) fileWorker(tasks chan FileTask, fileTypes models.FileTypes, wg *sync.WaitGroup, stats *Stats) {
	defer wg.Done()

	for task := range tasks {
		fileType := org.DetectFileType(task.file, fileTypes)
		if err := org.moveFile(task.path, task.file, fileType); err != nil {
			log.Printf("Error moving file %s: %v\n", task.file.Name, err)
			stats.mu.Lock()
			stats.failed++
			stats.mu.Unlock()
		} else {
			stats.mu.Lock()
			stats.succeeded++
			stats.mu.Unlock()
		}
	}
}

func (org *OrganizerService) DetectFileType(file models.File, fileTypes models.FileTypes) string {
	// Try MIME type detection first
	mimeType := mime.TypeByExtension(filepath.Ext(file.Name))
	if mimeType != "" {
		switch {
		case strings.HasPrefix(mimeType, "image/"):
			return "Images"
		case strings.HasPrefix(mimeType, "video/"):
			return "Videos"
		case strings.HasPrefix(mimeType, "audio/"):
			return "Audio"
		}
	}

	// Fallback to extension mapping
	ext := strings.ToLower(file.Extension)
	fileTypeMap := map[string]map[string]struct{}{
		"Images":    fileTypes.Images,
		"Videos":    fileTypes.Videos,
		"Archives":  fileTypes.Archives,
		"Documents": fileTypes.Documents,
		"Audio":     fileTypes.Audio,
		"Code":      fileTypes.Code,
		"Database":  fileTypes.Database,
		"Binaries":  fileTypes.Binaries,
		"Others":    fileTypes.Others,
	}

	for fileType, extensions := range fileTypeMap {
		if _, ok := extensions[ext]; ok {
			return fileType
		}
	}

	return "Others"
}

func (org *OrganizerService) moveFile(path string, file models.File, fileType string) error {
	const maxRetries = 3
	var err error

	for i := 0; i < maxRetries; i++ {
		err = org.attemptMove(path, file, fileType)
		if err == nil {
			return nil
		}
		if org.IsPermissionError(err) {
			return fmt.Errorf("permission denied: %v", err)
		}
		log.Printf("Retry %d moving file %s\n", i+1, file.Name)
	}
	return fmt.Errorf("failed to move file after %d attempts: %v", maxRetries, err)
}

func (org *OrganizerService) attemptMove(path string, file models.File, fileType string) error {
	oldPath := filepath.Join(path, file.Name)
	newPath := filepath.Join(path, org.folderService.GetFolderTypesPaths().Paths[fileType], file.Name)

	if _, err := os.Stat(oldPath); os.IsNotExist(err) {
		return fmt.Errorf("source file doesn't exist: %v", err)
	}

	return os.Rename(oldPath, newPath)
}

func (org *OrganizerService) IsPermissionError(err error) bool {
	return os.IsPermission(err)
}
