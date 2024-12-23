package services

import (
	"file-organizer/src/models"
)

type FileService struct{}

func NewFileService() *FileService {
	return &FileService{}
}

func (fs *FileService) GetFileTypes() models.FileTypes {
	return models.FileTypes{
		Images:    fs.createSet("jpg", "jpeg", "png", "gif", "bmp", "tiff", "webp", "svg", "raw", "heic"),
		Videos:    fs.createSet("mp4", "avi", "mkv", "mov", "wmv", "flv", "webm", "m4v", "mpeg"),
		Archives:  fs.createSet("zip", "rar", "7z", "tar", "gz", "bz2"),
		Documents: fs.createSet("pdf", "doc", "docx", "xls", "xlsx", "ppt", "pptx", "txt", "csv", "rtf"),
		Audio:     fs.createSet("mp3", "wav", "aac", "flac", "ogg", "m4a"),
		Code:      fs.createSet("dat", "cfg", "manifest", "md", "yml", "log", "json", "xml", "html", "css", "js", "py", "go", "java", "cpp"),
		Database:  fs.createSet("sql", "db", "sqlite"),
		Binaries:  fs.createSet("exe", "dll", "so"),
		Others:    fs.createSet("", "lnk", "msi", "iso", "apk", "crdownload"),
	}
}

func (fs *FileService) createSet(items ...string) map[string]struct{} {
	set := make(map[string]struct{}, len(items))
	for _, item := range items {
		set[item] = struct{}{}
	}
	return set
}
