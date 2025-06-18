package utils

import "github.com/gin-gonic/gin"

const (
	SUCCESS     = 200
	FAIL        = 400
	TOKEN_ERROR = 401
	NOT_FOUND   = 404
)

type ResponseJson struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// GetResponseJson 发送JSON响应。
//
// 参数:
// - c *gin.Context: Gin框架的上下文对象，用于处理HTTP请求和响应。
// - code int: HTTP状态码，表示响应的状态。
// - message string: 响应的消息，通常用于提供额外的信息。
// - data interface{}: 响应的数据部分，可以是任意类型的对象。
//
// 此函数封装了Gin框架的JSON方法，用于发送一个JSON格式的响应。它将响应的代码、消息和数据封装进一个结构体中，
// 并将其序列化为JSON格式后通过HTTP响应返回给客户端。这种方式提供了一种统一和方便的响应处理方式，
// 使得前端可以更容易地解析和处理后端返回的数据。
func GetResponseJson(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, ResponseJson{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
