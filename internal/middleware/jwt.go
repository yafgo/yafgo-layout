package middleware

import (
	"net/http"
	"yafgo/yafgo-layout/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

// JWTAuth 用于处理 jwt 鉴权，如果未登录则返回错误
//
//	abort: 默认为 true, 验证不通过会 Abort 当前请求并返回 401 错误;
//	为 false 时，验证不通过会返回错误信息，并不会终止请求, 仅仅做了 jwt 解析, 主要用于并不需要必须登录但是又需要解析 token 中用户信息的接口
func JWTAuth(j *jwtutil.JwtUtil, abort ...bool) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 这里 jwt 鉴权从标头 Authorization:Bearer xxxxx 中获取信息，并验证 JWT 的准确性
		// 登录时回返回token信息,前端需要把token存储到cookie或者本地localStorage中 不过需要跟后端协商过期时间 可以约定刷新令牌或者重新登录
		claims, err := j.ParserTokenFromHeader(c)
		if err != nil {
			if len(abort) > 0 && !abort[0] {
				// abort 传false的情况
				// c.Next()
				return
			} else {
				// abort 不传或传true的情况
				if err == jwtutil.ErrTokenExpired {
					unauthorized(c, "授权已过期")
					c.Abort()
					return
				}
				unauthorized(c, err.Error())
				c.Abort()
				return
			}
		}

		// 将用户信息存入 gin.context 里，后续 auth 包将从这里拿到当前用户数据
		c.Set("claims", claims)

		// c.Next()
	}
}

// unauthorized 响应 401，未传参 msg 时使用默认消息
// 登录失败、jwt 解析失败时调用
func unauthorized(c *gin.Context, msg ...string) {
	c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
		"message": defaultMessage("请求解析错误，请确认请求格式是否正确。上传文件请使用 multipart 标头，参数请使用 JSON 格式。", msg...),
	})
}

// defaultMessage 内用的辅助函数，用以支持默认参数默认值
// Go 不支持参数默认值，只能使用多变参数来实现类似效果
func defaultMessage(defaultMsg string, msg ...string) (message string) {
	if len(msg) > 0 {
		message = msg[0]
	} else {
		message = defaultMsg
	}
	return
}
