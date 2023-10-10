package middleware

import (
	"context"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/gin-gonic/gin"
)

// BasicAuth 为接口提供 BasicAuth 保护
//
//	配置项示例:
//	basic_auth:
//	  default:
//	    admin: adminpasswd
//	    user1: passwd1
//	    user2: passwd2
//	  the_biz:
//	    admin: thepasswd
//	    user1: thepasswd1
func BasicAuth(cfg *ycfg.Config) func(configKey ...string) gin.HandlerFunc {
	return func(configKey ...string) gin.HandlerFunc {
		var cfgKey = "default"
		if len(configKey) > 0 && configKey[0] != "" {
			cfgKey = configKey[0]
		}
		// 配置项中的key: "basic_auth.default"
		cfgKey = "basic_auth." + cfgKey
		accounts := cfg.GetStringMapString(cfgKey)
		if len(accounts) == 0 {
			accounts["demo"] = "ydaefmgooursr"
		}
		ylog.Warnf(context.Background(), "%s: %+v", cfgKey, accounts)
		return gin.BasicAuth(accounts)
	}
}
