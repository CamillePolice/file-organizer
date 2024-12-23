package tests

import (
	"file-organizer/src/services"
	"os"
	"path/filepath"
	"testing"
)

// TestFolderService_GetAllFiles verifies correct file discovery in directories
func TestFolderService_GetAllFiles(t *testing.T) {
	fs := services.NewFolderService()

	// Create temporary test directory
	tempDir, err := os.MkdirTemp("", "test-folder")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Create test files
	testFiles := []string{"test1.txt", "test2.jpg", "test3.pdf"}
	for _, file := range testFiles {
		path := filepath.Join(tempDir, file)
		if err := os.WriteFile(path, []byte("test"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	files, err := fs.GetAllFiles(tempDir)
	if err != nil {
		t.Fatalf("GetAllFiles failed: %v", err)
	}

	if len(files) != len(testFiles) {
		t.Errorf("Expected %d files, got %d", len(testFiles), len(files))
	}
}

// TestFolderService_PrepareFolders verifies creation of category folders
func TestFolderService_PrepareFolders(t *testing.T) {
	fs := services.NewFolderService()

	tempDir, err := os.MkdirTemp("", "test-folders")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	fs.PrepareFolders(tempDir)

	folderTypes := fs.GetFolderTypesPaths()
	for _, folderPath := range folderTypes.Paths {
		path := filepath.Join(tempDir, folderPath)
		if _, err := os.Stat(path); os.IsNotExist(err) {
			t.Errorf("Expected folder %s to exist", folderPath)
		}
	}
}
