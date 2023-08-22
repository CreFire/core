package core

import (
	"os"
	"path/filepath"
)

// FileExist check file is exist
func FileExist(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// ConvertToAbsolutePath converts a given path (which can be relative) to an absolute path.
func ConvertToAbsolutePath(path string) (string, error) {
	return filepath.Abs(path)
}

// EnsurePathSeparator ensures the correct OS-specific path separator is used.
func EnsurePathSeparator(path string) string {
	return filepath.ToSlash(path) // Use '/' as separator
	// OR
	// return filepath.FromSlash(path) // Use OS specific separator
}

// PathExists checks if a path exists.
func PathExists(path string) bool {
	_, err := os.Stat(path)
	if os.IsNotExist(err) {
		return false
	}
	return err == nil
}

// EnsureDirExists ensures a directory exists. If not, it'll create it.
func EnsureDirExists(path string) error {
	if !PathExists(path) {
		return os.MkdirAll(path, 0755)
	}
	return nil
}
