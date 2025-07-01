package logger

import (
	"context"
	"fmt"
	"go-gin/internal/app/config"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"gorm.io/gorm/logger"
)

/*
	位置：log/run.log

	控制台和日志文件：
	2025/06/30 10:59:12 初始化配置成功
*/

var (
	logFile        *os.File
	logInitialized bool
)

// SetupLog 初始化日志
func SetupLog() error {
	// 避免重复初始化
	if logInitialized {
		return nil
	}
	dirName := config.GetGlobalConfig().Log.DirName
	logPath := dirName + "/run.log"

	// 检查log目录是否存在，如果不存在则创建
	if _, err := os.Stat(dirName); os.IsNotExist(err) {
		if err := os.Mkdir(dirName, 0755); err != nil { // 更安全的目录权限
			return err
		}
	}

	// 打开或创建日志文件
	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}

	// 关闭之前打开的日志文件（如果存在）
	if logFile != nil {
		_ = logFile.Close()
	}
	logFile = file

	// 设置日志输出格式
	log.SetFlags(log.Ldate | log.Ltime)

	// 输出到控制台和文件
	multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(multiWriter)

	log.Println("初始化日志成功")
	logInitialized = true
	return nil
}

/*
	位置：log/sql.log

	控制台和日志文件：
	2025/06/30 11:02:51 [3.505ms] [rows:3] SELECT * FROM `users` WHERE `users`.`deleted_at` IS NULL LIMIT 10
*/

type SqlLogger struct {
	file *os.File
	mu   sync.Mutex // 互斥锁
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
	dirName := config.GetGlobalConfig().Log.DirName
	logPath := dirName + "/mysql.log"

	file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return nil, err
	}
	return &SqlLogger{file: file}, nil
}

// Close 关闭日志文件
func (l *SqlLogger) Close() error {
	return l.file.Close()
}

func (l *SqlLogger) LogMode(level logger.LogLevel) logger.Interface {
	return l // 可简化支持所有级别
}

func (l *SqlLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	l.logMessage("[INFO]", msg, data...)
}

func (l *SqlLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	l.logMessage("[WARN]", msg, data...)
}

func (l *SqlLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	l.logMessage("[ERROR]", msg, data...)
}

func (l *SqlLogger) Trace(ctx context.Context, begin time.Time, fc func() (string, int64), err error) {
	sql, rows := fc() // 先执行闭包，避免干扰耗时统计
	elapsed := time.Since(begin)

	now := time.Now().Format("2006/01/02 15:04:05")
	ms := float64(elapsed.Milliseconds()) + float64(elapsed.Nanoseconds()%1e6)/1e6

	fileMsg := fmt.Sprintf("%s [%.3fms] [rows:%d] %s\n", now, ms, rows, sql)
	consoleMsg := fmt.Sprintf(Reset+"%s "+YellowBold+"[%.3fms] "+BlueBold+"[rows:%d]"+Reset+" %s\n", now, ms, rows, sql)

	l.logCurrent(fileMsg, consoleMsg)
}

// logMessage 输出日志
func (l *SqlLogger) logMessage(prefix, msg string, data ...interface{}) {
	logMsg := fmt.Sprintf(prefix+" "+msg+"\n", data...)
	l.logCurrent(logMsg)
}

// Trace 输出日志
func (l *SqlLogger) logCurrent(file string, console ...string) {
	l.mu.Lock()
	defer l.mu.Unlock()

	fmt.Fprintf(l.file, "%s", file)
	for _, msg := range console {
		fmt.Print(msg)
	}
}
