package tests

import (
	"file-organizer/src/models"
	"file-organizer/src/services"
	"os"
	"path/filepath"
	"testing"
	"time"
)

// TestOrganizerService_OrganizeFiles verifies end-to-end file organization process
func TestOrganizerService_OrganizeFiles(t *testing.T) {
	fileService := services.NewFileService()
	folderService := services.NewFolderService()
	organizerService := services.NewOrganizerService(fileService, folderService)

	// Create temp test directory
	tempDir, err := os.MkdirTemp("", "test-organizer")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	// Prepare category folders
	folderService.PrepareFolders(tempDir)

	// Create test files
	testFiles := []models.File{
		{Name: "test.jpg", Extension: "jpg", Size: 100, FileType: "image/jpeg", Path: filepath.Join(tempDir, "test.jpg")},
		{Name: "test.pdf", Extension: "pdf", Size: 200, FileType: "application/pdf", Path: filepath.Join(tempDir, "test.pdf")},
	}

	// Create physical files
	for _, file := range testFiles {
		err := os.WriteFile(file.Path, []byte("test"), 0644)
		if err != nil {
			t.Fatal(err)
		}
	}

	err = organizerService.OrganizeFiles(testFiles, tempDir)
	if err != nil {
		t.Fatalf("OrganizeFiles failed: %v", err)
	}

	// Give time for async operations to complete
	time.Sleep(100 * time.Millisecond)

	// Verify files were moved to correct folders
	expectedPaths := map[string]string{
		"test.jpg": filepath.Join(tempDir, "Images", "test.jpg"),
		"test.pdf": filepath.Join(tempDir, "Documents", "test.pdf"),
	}

	for fileName, expectedPath := range expectedPaths {
		if _, err := os.Stat(expectedPath); os.IsNotExist(err) {
			t.Errorf("Expected %s to be moved to %s", fileName, expectedPath)
		}
	}
}

// TestOrganizerService_DetectFileType verifies correct file type detection based on extension and MIME type
func TestOrganizerService_DetectFileType(t *testing.T) {
	fileService := services.NewFileService()
	folderService := services.NewFolderService()
	organizerService := services.NewOrganizerService(fileService, folderService)

	tests := []struct {
		name     string
		file     models.File
		expected string
	}{
		{
			name:     "Image file",
			file:     models.File{Name: "test.jpg", Extension: "jpg", FileType: "image/jpeg"},
			expected: "Images",
		},
		{
			name:     "Document file",
			file:     models.File{Name: "test.pdf", Extension: "pdf", FileType: "application/pdf"},
			expected: "Documents",
		},
		{
			name:     "Unknown file",
			file:     models.File{Name: "test.xyz", Extension: "xyz", FileType: "application/octet-stream"},
			expected: "Others",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fileType := organizerService.DetectFileType(tt.file, fileService.GetFileTypes())
			if fileType != tt.expected {
				t.Errorf("Expected %s, got %s", tt.expected, fileType)
			}
		})
	}
}
