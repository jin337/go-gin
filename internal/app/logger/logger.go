package logger

import (
	"io"
	"log"
	"os"
)

/*
	位置：log/run.log

	控制台和日志文件：
	2025/06/30 10:59:12 初始化配置成功
*/

// SetupLog 初始化日志系统
func SetupLog() error {
	// 检查log目录是否存在，如果不存在则创建
	if _, err := os.Stat("log"); os.IsNotExist(err) {
		err := os.Mkdir("log", os.ModePerm)
		if err != nil {
			return err
		}
	}

	// 打开或创建日志文件
	file, err := os.OpenFile("log/run.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}
	// 设置日志输出格式
	log.SetFlags(log.Ldate | log.Ltime)

	// 输出到控制台和文件
	MultiWriter := io.MultiWriter(os.Stdout, file)
	log.SetOutput(MultiWriter)

	log.Println("初始化日志成功")
	return nil
}
