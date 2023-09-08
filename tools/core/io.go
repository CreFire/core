package core

import (
	"bufio"
	"core/tools/log"
	"fmt"
	"io"
	"os"
	"strings"
)

func ReaderFile(file *os.File) {

	reader := bufio.NewReader(file)
	for {
		line, isPrefix, err := reader.ReadLine()
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Errf(err)
			return
		}
		strings.Fields(string(line))
		// 如果行太长，isPrefix 将为 true，你需要自己处理这种情况。
		if isPrefix {
			fmt.Println("Skipping long line")
			continue
		}

	}
}
