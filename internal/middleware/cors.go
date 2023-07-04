package middleware

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func CORS() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowOrigins:  []string{"*"},
		AllowMethods:  []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
		AllowHeaders:  []string{"*"},
		ExposeHeaders: []string{"Content-Type"},
		// 是否允许带cookie, 为true时, allowOrigins 允许的域名中不能写*
		AllowCredentials: false,
		AllowOriginFunc: func(origin string) bool {
			// return origin == "https://github.com"
			return true
		},
		MaxAge: 24 * time.Hour,
	})
}
