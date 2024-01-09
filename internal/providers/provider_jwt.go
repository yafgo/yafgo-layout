package providers

import (
	"yafgo/yafgo-layout/pkg/jwtutil"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

func NewJwt(cfg *ycfg.Config) *jwtutil.JwtUtil {
	j := jwtutil.NewJWT(
		jwtutil.WithSignKey(cfg.GetString("jwt.sign_key")),
		jwtutil.WithIssuer(cfg.GetString("appname")),
	)
	return j
}
