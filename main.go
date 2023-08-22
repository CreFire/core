package main

import (
	"demo/tools/genProto"
	"encoding/json"
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	lumberjack "gopkg.in/natefinch/lumberjack.v2"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
)

func modifyWithAppend(slice []int) {
	fmt.Printf("before:%p\n", &slice)
	slice = append(slice, 10, 10)
	fmt.Printf("after: %p\n", &slice)
	slice[0] = 100
}

func main() {
	slice := make([]int, 6, 10)
	slice[0] = 1
	slice[1] = 2
	fmt.Printf("before o: %p\n", &slice)
	modifyWithAppend(slice)
	fmt.Printf("after o: %p\n", &slice)
	fmt.Println(slice) // 输出 [100 2 0 0 0 0 10 10]
}
func GetBitUint64(v uint64, offset int) bool {
	if offset > 63 || offset < 0 {
		return false
	}
	l := v >> offset
	return (l & 1) == 1
}

func SetBitUint64(v uint64, t bool, offset int) uint64 {
	if offset > 63 || offset < 0 {
		return v
	}
	if t {
		v = v | (1 << offset)
	} else {
		v = v &^ (1 << offset)
	}
	return v
}
func Demoi() {
	baseDir := "E:\\remoteWork\\server\\NetPb\\sv_pb_main"
	err := os.Chdir(baseDir)
	if err != nil {
		log.Println("err:", err)
	}
	protoPath := path.Join(baseDir, "proto")

	genProto.GetDirs(baseDir)
	curPwd, err := os.Getwd()
	if err != nil {
		log.Println(err.Error())
	}
	outPath := path.Join(curPwd, "./proto")
	cmd := exec.Command(fmt.Sprintf("protoc --go_out=%v --go-grpc_out=%v %v", outPath, outPath, protoPath))
	if cmd.Err != nil {
		log.Println(cmd.Err)
	}
}
func zapInit() {
	// 创建Lumberjack实例
	lumberjack := &lumberjack.Logger{
		Filename:   "/path/to/logfile.log",
		MaxSize:    100, // MB
		MaxBackups: 5,   // 保留5个备份文件
		MaxAge:     30,  // 保留30天内的文件
	}

	// 创建Zap实例
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder := zapcore.NewJSONEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, zapcore.AddSync(lumberjack), zapcore.InfoLevel)
	logger := zap.New(core)

	// 输出日志
	logger.Info("Hello, World!")
}

func ClientApi() {
	proxyUrl, _ := url.Parse("http://127.0.0.1:7890")
	client := &http.Client{Transport: &http.Transport{Proxy: http.ProxyURL(proxyUrl)}}

	url := "https://api.openai.com/v1/models"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("Authorization", "Bearer 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")

	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)
	var response map[string]interface{}
	json.Unmarshal(body, &response)

	responseJson, _ := json.MarshalIndent(response, "", "    ")
	os.WriteFile("response.txt", responseJson, 0644)

}
