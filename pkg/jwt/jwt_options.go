package jwt

import "time"

// JwtOption NewJWT 的选项参数
type JwtOption func(*JWT)

// NewJWT 获取一个jwt实例
func NewJWT(opts ...JwtOption) *JWT {
	j := &JWT{
		signKey:    defaultSignKey,
		expiresIn:  defaultExpiresIn,
		maxRefresh: defaultMaxRefresh,
		issuer:     defaultIssuer,
	}

	for _, opt := range opts {
		opt(j)
	}

	return j
}

// WithSignKey 配置 jwt signkey
func WithSignKey(val string) JwtOption {
	return func(p *JWT) {
		if val == "" {
			return
		}
		p.signKey = []byte(val)
	}
}

// WithExpiresIn 配置过期时间
func WithExpiresIn(val time.Duration) JwtOption {
	return func(p *JWT) {
		p.expiresIn = val
	}
}

// WithMaxRefresh 配置最大refresh时限
func WithMaxRefresh(val time.Duration) JwtOption {
	return func(p *JWT) {
		p.maxRefresh = val
	}
}

// WithIssuer 配置 claims 的 issuer
func WithIssuer(val string) JwtOption {
	return func(p *JWT) {
		p.issuer = val
	}
}

// WithSubject 配置 claims 的 subject
func WithSubject(val string) JwtOption {
	return func(p *JWT) {
		p.subject = val
	}
}

// WithAudience 配置 claims 的 audience
func WithAudience(val string) JwtOption {
	return func(p *JWT) {
		p.audience = val
	}
}
