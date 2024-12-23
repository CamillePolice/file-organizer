package models

type File struct {
	Name      string
	Extension string
	Size      int64
	FileType  string
	Path      string
}

type FileTypes struct {
	Images    map[string]struct{}
	Videos    map[string]struct{}
	Archives  map[string]struct{}
	Documents map[string]struct{}
	Audio     map[string]struct{}
	Code      map[string]struct{}
	Database  map[string]struct{}
	Binaries  map[string]struct{}
	Others    map[string]struct{}
}
