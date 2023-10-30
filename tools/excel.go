package tools

import (
	"fmt"
	"github.com/core/tools/log"
	"os"
)

func ReadExcel(path string) {
	_, err := os.Open(path)
	if err != nil {
		return
	}
}

func TryCreateDir(path string) {
	if path == "" {
		log.Error("path is nil", log.String("path", path))
		return
	}
	if _, err := os.Stat(path); os.IsNotExist(err) {
		fmt.Printf("目录 %s 不存在，将创建它\n", path)

		// 创建目录
		err := os.MkdirAll(path, os.ModePerm)
		if err != nil {
			fmt.Printf("创建目录 %s 失败：%v\n", path, err)
			return
		}

		fmt.Printf("目录 %s 创建成功\n", path)
	} else {
		fmt.Printf("目录 %s 已经存在\n", path)
	}
}
