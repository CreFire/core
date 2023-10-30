package main

import (
	"bufio"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/core/tools/log"
	"os"
	"strings"
)

func main() {
	sensitivewordMap = make(map[string]struct{})
	unlimitwordMap = make(map[string]struct{})
	ReadExcel()
	ReadFile()
	CreateExcel()
}

var sensitivewordMap map[string]struct{}
var unlimitwordMap map[string]struct{}

func ReadExcel() {
	log.Info("start")
	file, err := excelize.OpenFile("./sensitiveword.xlsx")
	if err != nil {
		log.Error("openfile", log.Err(err))
	}

	log.InfoF(file.SheetCount, len(file.Sheet))
	comment := file.GetSheetMap()

	for s, sheet := range comment {
		log.Info("comments", log.Int("index", s), log.String("sheet", sheet))

		cells := file.GetRows(sheet)
		for _, cell := range cells {
			if len(cell) < 2 {
				continue
			}
			if s == 1 {
				sensitivewordMap[cell[1]] = struct{}{}
			} else if s == 2 {
				sensitivewordMap[cell[1]] = struct{}{}
			} else {
				unlimitwordMap[cell[1]] = struct{}{}
			}
		}
	}

}

func ReadFile() {
	// 打开文件
	file, err := os.Open("./blockedword.txt")
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	// 创建一个新的扫描器
	scanner := bufio.NewScanner(file)
	const maxCapacity = 512 * 1024 // 512KB
	scanner.Buffer(make([]byte, maxCapacity), maxCapacity)
	// 逐行读取文件
	for scanner.Scan() {
		line := scanner.Text()

		// 忽略以 '-' 开头的行
		if strings.HasPrefix(line, "-") {
			continue
		}

		// 使用 , ` 和空格作为分隔符分割字符串
		delimiters := []string{",", "`", " ", "、"}
		for _, delimiter := range delimiters {
			line = strings.ReplaceAll(line, delimiter, " ")
		}

		words := strings.Fields(line) // 使用空格分割字符串

		// 打印分割后的单词
		for _, word := range words {
			sensitivewordMap[word] = struct{}{}
		}
	}

	// 检查是否有错误
	if err := scanner.Err(); err != nil {
		fmt.Println("Error reading file:", err)
	}
}

func CreateExcel() {
	xlsx := excelize.NewFile()
	// 创建一个工作表
	xlsx.NewSheet("sensitiveword-敏感词汇表")
	xlsx.NewSheet("unlimitword-白名单词表")

	// 设置单元格的值
	SaveSheet(0, "sensitiveword-敏感词汇表", sensitivewordMap, xlsx)
	SaveSheet(1, "unlimitword-白名单词表", unlimitwordMap, xlsx)
	// 设置工作簿的默认工作表
	xlsx.SetActiveSheet(0)
	// 根据指定路径保存文件
	err := xlsx.SaveAs("./word.xlsx")
	if err != nil {
		log.Error("current", log.Err(err))
	}
	log.InfoF("end")
}
func SaveSheet(index int, name string, maps map[string]struct{}, xlsx *excelize.File) *excelize.File {
	switch index {
	case 1:
		xlsx.SetCellValue(name, "A1", "敏感词id")
		xlsx.SetCellValue(name, "B1", "敏感词字库")
		xlsx.SetCellValue(name, "B2", "sensitiveword")
	case 3:
		xlsx.SetCellValue(name, "A1", "白名单词id")
		xlsx.SetCellValue(name, "B1", "白名单词字库")
		xlsx.SetCellValue(name, "B2", "unlimitword")
	}
	xlsx.SetCellValue(name, "A2", "ID")
	xlsx.SetCellValue(name, "A3", "Int")
	xlsx.SetCellValue(name, "B3", "String")
	n := 4
	for str, _ := range maps {
		xlsx.SetCellValue(name, fmt.Sprintf("A%d", n), n)
		xlsx.SetCellValue(name, fmt.Sprintf("B%d", n), str)
		n++
	}
	log.Info("endSheet", log.Int("index", index))
	return xlsx
}
