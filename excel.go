package main

import (
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"log"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main2() {
	// 获取文件名参数，如果未传递参数则使用默认值
	fileName := "example.xlsx"
	if len(os.Args) > 1 {
		fileName = os.Args[1]
	}
	// 自动添加 .xlsx 扩展名
	if !strings.HasSuffix(fileName, ".xlsx") {
		fileName += ".xlsx"
	}
	// 获取文件路径
	filePath, err := filepath.Abs(fileName)
	if err != nil {
		log.Fatal(err)
	}

	// 打开 Excel 文件
	f, err := excelize.OpenFile(filePath)
	if err != nil {
		log.Fatal(err)
	}

	// 获取第一个工作表的所有行数据
	rows := f.GetRows(f.GetSheetName(1))

	// 创建一个随机数生成器，并设置种子
	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	// 从第一列随机取一行数字并输出
	row := rows[r.Intn(len(rows))]
	fmt.Println(row[0])
}
