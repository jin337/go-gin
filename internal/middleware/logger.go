package middleware

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

/*
	位置：log/server.log

	控制台和日志文件：
	2025/06/30 11:02:51 | 200 | 5.6025ms | 127.0.0.1 | POST /api/v1/user {"page": 1,"page_size": 10}
	2025/06/30 11:02:51 | 200 | 5.6025ms | 127.0.0.1 | GET /api/v1/user-list page=1&page_size=10
*/

// ANSI 背景色代码
const (
	Reset   = "\033[0m"
	Red     = "\033[31m"
	Green   = "\033[32m"
	Yellow  = "\033[33m"
	Blue    = "\033[34m"
	Magenta = "\033[35m"
	Cyan    = "\033[36m"
	White   = "\033[37m"

	bgRed     = "\033[41m"
	bgGreen   = "\033[42m"
	bgYellow  = "\033[43m"
	bgBlue    = "\033[44m"
	bgMagenta = "\033[45m"
	bgCyan    = "\033[46m"
	bgWhite   = "\033[47m"
)

// 根据状态码选择颜色
func getColorByStatusCode(statusCode int) string {
	switch {
	case statusCode >= 500:
		return bgRed
	case statusCode >= 400:
		return bgYellow
	case statusCode >= 300:
		return bgCyan
	case statusCode >= 200:
		return bgGreen
	default:
		return Reset
	}
}

// 根据请求方法选择颜色
func getColorByMethod(method string) string {
	switch method {
	case "GET":
		return bgGreen
	case "POST":
		return bgBlue
	case "PUT":
		return bgYellow
	case "DELETE":
		return bgRed
	case "PATCH":
		return bgMagenta
	default:
		return Reset
	}
}

// 根据耗时高低选择颜色
func getColorByLatency(latency time.Duration) string {
	if latency < 100*time.Millisecond {
		return Green
	} else if latency < 500*time.Millisecond {
		return Yellow
	} else {
		return Red
	}
}

// 日志中间件
func LoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()

		// 缓存请求体
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		// 获取请求耗时
		latency := time.Since(start)
		// 获取客户端IP
		clientIP := c.ClientIP()
		// 获取请求方法
		method := c.Request.Method
		// 获取状态码
		statusCode := c.Writer.Status()
		// 获取请求路径
		path := c.Request.URL.Path
		// 获取请求参数
		var params string
		switch method {
		case "GET":
			params = c.Request.URL.RawQuery
		case "POST":
			if c.ContentType() == "application/json" {
				var out bytes.Buffer
				if err := json.Compact(&out, bodyBytes); err == nil {
					params = out.String()
				} else {
					params = string(bodyBytes)
				}
			} else {
				c.Request.ParseForm()
				params = c.Request.Form.Encode()
			}
		}

		// 给状态码加上背景色
		statusColored := fmt.Sprintf("%s %d %s", getColorByStatusCode(statusCode), statusCode, Reset)
		// 给请求方法加上颜色
		coloredMethod := fmt.Sprintf("%s %s %s", getColorByMethod(method), method, Reset)
		// 给耗时加上颜色
		coloredLatency := fmt.Sprintf("%s%v%s", getColorByLatency(latency), latency, Reset)
		// 给请求参数加上颜色
		coloredParams := fmt.Sprintf("%s%s%s", Magenta, params, Reset)

		// 格式化输出日志
		consoleLogLine := fmt.Sprintf(
			"%s | %s | %v | %s | %s %s %s",
			time.Now().Format("2006/01/02 15:04:05"),
			statusColored,
			coloredLatency,
			clientIP,
			coloredMethod,
			path,
			coloredParams,
		)

		fileLogLine := fmt.Sprintf(
			"%s | %d | %v | %s | %s %s %s",
			time.Now().Format("2006/01/02 15:04:05"),
			statusCode,
			latency,
			clientIP,
			method,
			path,
			params,
		)

		// 打印到控制台
		fmt.Println(consoleLogLine)
		fmt.Println() // 为了美观，添加额外的换行

		// 输出到日志文件
		const dirName = "log"
		file, err := os.OpenFile(dirName+"/server.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0666)
		if err != nil {
			fmt.Println("打开日志文件时出错:", err)
			return
		}
		defer file.Close()
		fmt.Fprintln(file, fileLogLine)
	}
}
