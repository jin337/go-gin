package logger

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"gorm.io/gorm/logger"
)

/*
	位置：log/run.log

	控制台和日志文件：
	2025/06/30 10:59:12 初始化配置成功
*/

// SetupLog 初始化日志
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

/*
	位置：log/sql.log

	控制台和日志文件：
	2025/06/30 11:02:51 [3.505ms] [rows:3] SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT 10
*/

type SqlLogger struct {
	file *os.File
}

const (
	Reset       = "\033[0m"
	Red         = "\033[31m"
	Green       = "\033[32m"
	Yellow      = "\033[33m"
	Blue        = "\033[34m"
	Magenta     = "\033[35m"
	Cyan        = "\033[36m"
	White       = "\033[37m"
	BlueBold    = "\033[34;1m"
	MagentaBold = "\033[35;1m"
	RedBold     = "\033[31;1m"
	YellowBold  = "\033[33;1m"
)

// 设置mysql日志
func SetMySqlLogger() (logger.Interface, error) {
	file, err := os.OpenFile("log/sql.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &SqlLogger{file: file}, nil
}

func (l *SqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l // 可简化支持所有级别
}

func (l *SqlLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	logMsg := fmt.Sprintf("[INFO] "+msg+"\n", data...)
	l.logCurrent(logMsg)
}

func (l *SqlLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	logMsg := fmt.Sprintf("[WARN] "+msg+"\n", data...)
	l.logCurrent(logMsg)
}

func (l *SqlLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	logMsg := fmt.Sprintf("[ERROR] "+msg+"\n", data...)
	l.logCurrent(logMsg)
}

func (l *SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	elapsed := time.Since(begin)
	sql, rows := fc()
	now := time.Now().Format("2006/01/02 15:04:05") // 格式化当前时间
	fileMsg := fmt.Sprintf("%s [%.3fms] [rows:%d] %s\n", now, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	consoleMsg := fmt.Sprintf(Reset+"%s "+YellowBold+"[%.3fms] "+BlueBold+"[rows:%d]"+Reset+" %s\n", now, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	l.logCurrent(fileMsg, consoleMsg)
}

func (l *SqlLogger) logCurrent(file string, console ...string) {
	fmt.Fprintf(l.file, "%s", file)
	for _, msg := range console {
		fmt.Print(msg)
	}

}
