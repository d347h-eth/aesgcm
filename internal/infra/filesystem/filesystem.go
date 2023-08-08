package filesystem

import "os"

type (
	FileSystem struct{}
)

func NewFileSystem() *FileSystem {
	return &FileSystem{}
}

// Read reads file from FS
func (fs FileSystem) Read(filename string) ([]byte, error) {
	return os.ReadFile(filename)
}

// Write writes file into FS
func (fs FileSystem) Write(filename string, data []byte) error {
	return os.WriteFile(filename, data, 0644)
}

// ResourceExist checks if file exists in FS
func (fs FileSystem) ResourceExist(filename string) bool {
	_, err := os.Stat(filename)
	return !os.IsNotExist(err)
}
