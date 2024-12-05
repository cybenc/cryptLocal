package main

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/cybenc/cryptLocal/cmd"
	"github.com/cybenc/cryptLocal/common"
	"github.com/cybenc/cryptLocal/utils"
	"github.com/sirupsen/logrus"
)

func initLog() *os.File {
	// 获取当前日期
	today := time.Now().Format("2006-01-02")
	homeDir := utils.CIns.GetLogDir()
	// 获取今天的日期
	logFilePath := filepath.Join(homeDir, "log_"+today+".log")
	// 打开（或创建）日志文件
	file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("open log file error:", err)
		os.Exit(1)
	}
	// 设置Logrus同时输出到控制台和文件
	//logrus.SetOutput(file)
	logrus.SetOutput(io.MultiWriter(os.Stdout, file))
	// 设置日志格式（可选）
	logrus.SetFormatter(&logrus.TextFormatter{
		TimestampFormat: "2006-01-02 15:04:05",
		FullTimestamp:   true,
	})
	return file
}

func additionWork() {
	// 获取可执行文件名称
	path, _ := os.Executable()
	_, exec := filepath.Split(path)
	if exec != common.ExecuteName {
		common.ExecuteName = exec
	}

	// 判断是否已经初始化过
	common.Initized = utils.CIns.CheckConfig()
}

func main() {
	log_file := initLog()
	defer log_file.Close()

	additionWork()
	cmd.Execute()
}
