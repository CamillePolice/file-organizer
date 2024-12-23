package services

import (
	"file-organizer/src/models"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type FolderService struct{}

func NewFolderService() *FolderService {
	return &FolderService{}
}

func (fs *FolderService) GetAllFiles(path string) ([]models.File, error) {
	var files []models.File
	folderPath := path

	err := filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if path == folderPath {
			return nil
		}

		if info.IsDir() {
			for _, folderTypePath := range fs.GetFolderTypesPaths().Paths {
				if strings.HasPrefix(path, filepath.Join(folderPath, folderTypePath)) {
					return filepath.SkipDir
				}
			}
			subFiles, err := fs.GetAllFiles(path)
			if err != nil {
				return err
			}
			files = append(files, subFiles...)
			return nil
		}

		file := models.File{
			Name:      info.Name(),
			Extension: strings.TrimPrefix(filepath.Ext(info.Name()), "."),
			Size:      info.Size(),
			FileType:  http.DetectContentType([]byte(info.Name())),
			Path:      filepath.Join(folderPath, info.Name()),
		}

		files = append(files, file)

		return nil
	})

	if err != nil {
		return nil, err
	}

	return files, nil
}

func (fs *FolderService) GetFolderTypesPaths() models.FolderTypesPaths {
	return models.FolderTypesPaths{
		Paths: map[string]string{
			"Images":    "Images",
			"Videos":    "Videos",
			"Archives":  "Archives",
			"Documents": "Documents",
			"Audio":     "Audio",
			"Code":      "Code",
			"Database":  "Database",
			"Binaries":  "Binaries",
			"Others":    "Others",
		},
	}
}

func (fs *FolderService) PrepareFolders(path string) {
	folderTypesPaths := fs.GetFolderTypesPaths()

	for _, folderPath := range folderTypesPaths.Paths {
		if _, err := os.Stat(filepath.Join(path, folderPath)); os.IsNotExist(err) {
			err := os.MkdirAll(filepath.Join(path, folderPath), os.ModePerm)
			if err != nil {
				fmt.Println("Error creating folder", folderPath)
			}
		}
	}
}

func (fs *FolderService) moveToFolder(path string, file models.File, fileType string) {
	folderTypesPaths := fs.GetFolderTypesPaths()

	oldPath := filepath.Join(path, file.Name)
	newPath := filepath.Join(path, folderTypesPaths.Paths[fileType], file.Name)

	err := os.Rename(oldPath, newPath)

	if err != nil {
		fmt.Println("Error moving file", file.Name)
	}
}
