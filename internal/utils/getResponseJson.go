package utils

import "github.com/gin-gonic/gin"

const (
	SUCCESS     = 200
	FAIL        = 400
	TOKEN_ERROR = 401
	NOT_FOUND   = 404
	SERVER_BUSY = 503
)

type ResponseJson struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func GetResponseJson(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, ResponseJson{
		Code:    code,
		Message: message,
		Data:    data,
	})
}
