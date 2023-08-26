package tools

import "os"

func ReadExcel(path string) {
	_, err := os.Open(path)
	if err != nil {
		return
	}
}
