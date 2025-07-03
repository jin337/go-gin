package middleware

import (
	"go-gin/internal/app/config"
	"go-gin/internal/app/database"
	"go-gin/internal/model"
	"go-gin/internal/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"gorm.io/gorm"
)

// 自定义 Claims 结构体
type Claims struct {
	UserID uint `json:"user_id"`
	jwt.RegisteredClaims
}

func AuthMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Token")
		if authHeader == "" {
			utils.GetResponseJson(ctx, utils.TOKEN_ERROR, "需要授权Token", nil)
			ctx.Abort()
			return
		}

		// 检查 Token 是否在黑名单中
		var blacklistedToken model.BlacklistedToken
		if err := database.DB.Where("token = ?", authHeader).First(&blacklistedToken).Error; err != nil {
			if err != gorm.ErrRecordNotFound {
				// 数据库查询出错
				utils.GetResponseJson(ctx, utils.SERVER_BUSY, "服务器内部错误", nil)
				ctx.Abort()
				return
			}
		} else {
			// Token 在黑名单中
			utils.GetResponseJson(ctx, utils.TOKEN_ERROR, "Token已失效", nil)
			ctx.Abort()
			return
		}

		// 解析 JWT Token
		tokenSecret := config.GetGlobalConfig().Service.TokenSecret
		token, err := jwt.ParseWithClaims(authHeader, &Claims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, jwt.ErrSignatureInvalid
			}
			return []byte(tokenSecret), nil
		})

		if claims, ok := token.Claims.(*Claims); ok && token.Valid {
			ctx.Set("UserID", claims.UserID)
			ctx.Next()
		} else {
			if ve, ok := err.(*jwt.ValidationError); ok {
				switch {
				case ve.Errors&jwt.ValidationErrorExpired != 0:
					utils.GetResponseJson(ctx, utils.TOKEN_ERROR, "Token已过期", nil)
				case ve.Errors&jwt.ValidationErrorNotValidYet != 0:
					utils.GetResponseJson(ctx, utils.TOKEN_ERROR, "Token尚未生效", nil)
				default:
					utils.GetResponseJson(ctx, utils.TOKEN_ERROR, "Token验证失败", nil)
				}
			} else {
				utils.GetResponseJson(ctx, utils.TOKEN_ERROR, "Token声明无效", nil)
			}
			ctx.Abort()
		}
	}
}
