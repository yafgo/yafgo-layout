package g

import (
	"sync"
	"yafgo/yafgo-layout/pkg/jwtutil"
)

// 全局 Jwt 对象
var instJwt *jwtutil.JwtUtil
var onceJwt sync.Once

// Jwt 获取全局 Jwt 对象
func Jwt() *jwtutil.JwtUtil {
	onceJwt.Do(func() {
		cfg := Cfg()
		instJwt = jwtutil.NewJWT(
			jwtutil.WithSignKey(cfg.GetString("jwt.sign_key")),
			jwtutil.WithIssuer(AppName()),
		)
	})

	return instJwt
}
