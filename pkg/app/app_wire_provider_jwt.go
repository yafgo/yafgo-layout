package app

import (
	"yafgo/yafgo-layout/pkg/jwt"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

func newJwt(cfg *ycfg.Config) *jwt.JWT {
	j := jwt.NewJWT(
		jwt.WithSignKey(cfg.GetString("jwt.sign_key")),
		// jwt.WithIssuer(AppName()),
	)
	return j
}
