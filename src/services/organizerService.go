package services

import (
	"file-organizer/src/models"
	"fmt"
)

type OrganizerService struct {
	fileService   *FileService
	folderService *FolderService
}

func NewOrganizerService(fileService *FileService, folderService *FolderService) *OrganizerService {
	return &OrganizerService{fileService: fileService, folderService: folderService}
}

func (fs *OrganizerService) OrganizeFiles(files []models.File, path string) {
	fileTypes := fs.fileService.GetFileTypes()
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

	for _, file := range files {
		extension := file.Extension
		found := false

		fmt.Println("File", file.Extension)
		for fileType, extensions := range fileTypeMap {
			if _, ok := extensions[extension]; ok {
				fmt.Println(fileType, fileType)
				fs.folderService.moveToFolder(path, file, fileType)
				found = true
				break
			}
		}

		if !found {
			fs.folderService.moveToFolder(path, file, "Others")
		}
	}
}
