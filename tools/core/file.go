package core

import (
	"core/tools/log"
	"fmt"
	"io"
	"os"
	"path/filepath"
)

// WrappedFile 封装了一个 *os.File 和附加的功能
type WrappedFile struct {
	file *os.File
}

// NewWrappedFile 创建一个新的 WrappedFile 实例
func NewWrappedFile(filename string, flag int, perm os.FileMode) (*WrappedFile, error) {
	file, err := os.OpenFile(filename, flag, perm)
	if err != nil {
		return nil, err
	}
	return &WrappedFile{file: file}, nil
}

// Write 封装了原始的 Write 方法并添加了错误检查
func (wf *WrappedFile) Write(b []byte) (int, error) {
	n, err := wf.file.Write(b)
	if err != nil {
		fmt.Printf("Write error: %v\n", err)
		return 0, err
	}
	fmt.Printf("Successfully wrote %d bytes\n", n)
	return n, nil
}

// Read 封装了原始的 Read 方法并添加了错误检查
func (wf *WrappedFile) Read(b []byte) (int, error) {
	n, err := wf.file.Read(b)
	if err != nil && err != io.EOF {
		fmt.Printf("Read error: %v\n", err)
		return 0, err
	}
	fmt.Printf("Successfully read %d bytes\n", n)
	return n, err
}

// Close 封装了原始的 Close 方法并添加了错误检查
func (wf *WrappedFile) Close() error {
	err := wf.file.Close()
	if err != nil {
		fmt.Printf("Close error: %v\n", err)
		return err
	}
	fmt.Println("File successfully closed")
	return nil
}

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

func GetDirs(path string) {

	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		log.Error("err ", log.Err(err))
		return
	}
	defer func(dir *os.File) {
		err = dir.Close()
		if err != nil {
			log.Error("err ", log.Err(err))
		}
	}(dir)

	// 读取目录下的所有文件和子目录
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Error("err ", log.Err(err))
	}

	// 遍历文件和子目录
	for _, fileInfo := range fileInfos {
		// 判断是否为目录
		if fileInfo.IsDir() {
			fmt.Println("Directory:", fileInfo.Name())
		} else {
			fmt.Println("File:", fileInfo.Name())
		}
	}
}
