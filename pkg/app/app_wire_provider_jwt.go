package app

import (
	"yafgo/yafgo-layout/pkg/jwtutil"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

func newJwt(cfg *ycfg.Config) *jwtutil.JwtUtil {
	j := jwtutil.NewJWT(
		jwtutil.WithSignKey(cfg.GetString("jwt.sign_key")),
		// jwtutil.WithIssuer(AppName()),
	)
	return j
}
