package tests

import (
	"file-organizer/src/services"
	"testing"
)

// TestFileService_GetFileTypes verifies that file extensions are correctly mapped to their categories
func TestFileService_GetFileTypes(t *testing.T) {
	fs := services.NewFileService()
	fileTypes := fs.GetFileTypes()

	tests := []struct {
		name      string
		extension string
		category  map[string]struct{}
		expected  bool
	}{
		{"Image", "jpg", fileTypes.Images, true},
		{"Video", "mp4", fileTypes.Videos, true},
		{"Invalid Image", "invalid", fileTypes.Images, false},
		{"Code", "go", fileTypes.Code, true},
		{"Document", "pdf", fileTypes.Documents, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, actual := tt.category[tt.extension]
			if actual != tt.expected {
				t.Errorf("extension %s: got %v, want %v", tt.extension, actual, tt.expected)
			}
		})
	}
}
