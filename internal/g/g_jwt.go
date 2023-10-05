package g

import (
	"sync"
	"yafgo/yafgo-layout/pkg/jwt"
)

// 全局 Jwt 对象
var instJwt *jwt.JWT
var onceJwt sync.Once

// Jwt 获取全局 Jwt 对象
func Jwt() *jwt.JWT {
	onceJwt.Do(func() {
		cfg := Cfg()
		instJwt = jwt.NewJWT(
			jwt.WithSignKey(cfg.GetString("jwt.sign_key")),
			jwt.WithIssuer(AppName()),
		)
	})

	return instJwt
}
