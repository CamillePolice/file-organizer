package main

import (
	"file-organizer/src/services"
	"fmt"
)

func main() {
	folderService := services.NewFolderService()
	fileService := services.NewFileService()
	organizerService := services.NewOrganizerService(fileService, folderService)
	path := "../test-folder"

	files, err := folderService.GetAllFiles(path)

	if err != nil {
		fmt.Println(err)
	}

	if len(files) == 0 {
		fmt.Println("No files found")
	}

	folderService.PrepareFolders(path)
	organizerService.OrganizeFiles(files, path)
}
