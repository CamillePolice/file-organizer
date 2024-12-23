package main

import (
	"file-organizer/src/services"
	"flag"
	"fmt"
	"os"
)

func main() {
	path := flag.String("path", "", "Path to organize files")
	flag.Parse()

	if *path == "" {
		fmt.Println("Please provide a path using -path flag")
		os.Exit(1)
	}

	if _, err := os.Stat(*path); os.IsNotExist(err) {
		fmt.Printf("Path %s does not exist\n", *path)
		os.Exit(1)
	}

	folderService := services.NewFolderService()
	fileService := services.NewFileService()
	organizerService := services.NewOrganizerService(fileService, folderService)

	files, err := folderService.GetAllFiles(*path)
	if err != nil {
		fmt.Printf("Error getting files: %v\n", err)
		os.Exit(1)
	}

	if len(files) == 0 {
		fmt.Println("No files found")
		return
	}

	folderService.PrepareFolders(*path)
	organizerService.OrganizeFiles(files, *path)
}
