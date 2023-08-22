package genProto

import (
	"fmt"
	"log"
	"os"
)

func GetDirs(path string) {

	// 打开目录
	dir, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer dir.Close()

	// 读取目录下的所有文件和子目录
	fileInfos, err := dir.Readdir(-1)
	if err != nil {
		log.Fatal(err)
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
