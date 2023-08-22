package antnet

import (
	"fmt"
	"math/rand"
)

func Println(a ...interface{}) (int, error) {
	return fmt.Println(a...)
}

func Printf(format string, a ...interface{}) (int, error) {
	return fmt.Println(fmt.Sprintf(format, a...))
}

func Sprintf(format string, a ...interface{}) string {
	return fmt.Sprintf(format, a...)
}

const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

func RandStringBytes(n int) string {
	b := make([]byte, n)
	//rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
