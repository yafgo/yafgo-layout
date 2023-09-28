// Package jwt 处理 JWT 认证
package jwt

import (
	"context"
	"errors"
	"strings"
	"time"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt/v4"
)

var (
	ErrTokenExpired           error = errors.New("令牌已过期")
	ErrTokenExpiredMaxRefresh error = errors.New("令牌已过最大刷新时间")
	ErrTokenMalformed         error = errors.New("请求令牌格式有误")
	ErrTokenInvalid           error = errors.New("请求令牌无效")
	ErrHeaderEmpty            error = errors.New("请先登录")
	ErrHeaderMalformed        error = errors.New("请求头中 Authorization 格式有误")
)

var (
	defaultSignKey    []byte        = []byte("signkey")  // 秘钥，用以加密 JWT
	defaultExpiresIn  time.Duration = time.Hour * 24 * 7 // Token 过期时间
	defaultMaxRefresh time.Duration = time.Hour * 24     // 刷新 Token 的最大过期时间
	defaultIssuer     string        = "yafgo"            // Token 的发行者
)

// JWTCustomClaims 自定义载荷
type JWTCustomClaims struct {
	ID         uint64 `json:"id"`
	BufferTime int64  `json:"buffer_time"`

	// RegisteredClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwtpkg.RegisteredClaims
}

// JWT 定义一个jwt对象
type JWT struct {
	// 秘钥，用以加密 JWT
	signKey []byte

	// Token 过期时间
	expiresIn time.Duration

	// Token 的发行者
	issuer string

	// 刷新 Token 的最大过期时间
	maxRefresh time.Duration
}

type JwtOption func(*JWT)

// WithSignKey 配置 jwt signkey
func WithSignKey(val string) JwtOption {
	return func(p *JWT) {
		if val == "" {
			return
		}
		p.signKey = []byte(val)
	}
}

// WithIssuer 配置issuer
func WithIssuer(val string) JwtOption {
	return func(p *JWT) {
		p.issuer = val
	}
}

// WithExpiresIn 配置过期时间
func WithExpiresIn(val time.Duration) JwtOption {
	return func(p *JWT) {
		p.expiresIn = val
	}
}

// WithMaxRefresh 配置最大refresh次数
func WithMaxRefresh(val time.Duration) JwtOption {
	return func(p *JWT) {
		p.maxRefresh = val
	}
}

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

// ParserToken 解析 Token，中间件中调用
func (jwt *JWT) ParserToken(c *gin.Context) (*JWTCustomClaims, error) {

	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return nil, parseErr
	}

	// 1. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 2. 解析出错
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		if ok {
			if validationErr.Errors == jwtpkg.ValidationErrorMalformed {
				return nil, ErrTokenMalformed
			} else if validationErr.Errors == jwtpkg.ValidationErrorExpired {
				return nil, ErrTokenExpired
			}
		}
		return nil, ErrTokenInvalid
	}

	// 3. 将 token 中的 claims 信息解析出来和 JWTCustomClaims 数据结构进行校验
	if claims, ok := token.Claims.(*JWTCustomClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// RefreshToken 更新 Token，用以提供 refresh token 接口
func (jwt *JWT) RefreshToken(c *gin.Context) (string, error) {

	// 1. 从 Header 里获取 token
	tokenString, parseErr := jwt.getTokenFromHeader(c)
	if parseErr != nil {
		return "", parseErr
	}

	// 2. 调用 jwt 库解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenString)

	// 3. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		validationErr, ok := err.(*jwtpkg.ValidationError)
		// 满足 refresh 的条件：只是单一的报错 ValidationErrorExpired
		if !ok || validationErr.Errors != jwtpkg.ValidationErrorExpired {
			return "", err
		}
	}

	// 4. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*JWTCustomClaims)

	// 5. 检查是否过了『最大允许刷新的时间』
	x := TimenowInTimezone().Add(-jwt.maxRefresh)
	if claims.IssuedAt.After(x) {
		// 修改过期时间
		claims.RegisteredClaims.ExpiresAt = jwt.NewExpiresAt()
		return jwt.CreateToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

// IssueToken 生成  Token，在登录成功时调用
func (jwt *JWT) IssueToken(userID uint64) (token string, claims JWTCustomClaims, err error) {

	// 1. 构造用户 claims 信息(负荷)
	expiresAt := jwt.NewExpiresAt()
	issuedAt := jwtpkg.NewNumericDate(TimenowInTimezone())
	claims = JWTCustomClaims{
		ID:         userID,
		BufferTime: int64(jwt.maxRefresh / time.Second),

		RegisteredClaims: jwtpkg.RegisteredClaims{
			NotBefore: issuedAt,   // 签名生效时间
			IssuedAt:  issuedAt,   // 首次签名时间（后续刷新 Token 不会更新）
			ExpiresAt: expiresAt,  // 签名过期时间
			Issuer:    jwt.issuer, // 签名颁发者
		},
	}

	// 2. 根据 claims 生成token对象
	token, err = jwt.CreateToken(claims)
	if err != nil {
		ylog.Error(context.Background(), err)
	}

	return
}

// CreateToken 创建 Token，内部使用，外部请调用 IssueToken
func (jwt *JWT) CreateToken(claims JWTCustomClaims) (string, error) {
	// 使用HS256算法进行token生成
	token := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return token.SignedString(jwt.signKey)
}

// NewExpiresAt 过期时间
func (jwt *JWT) NewExpiresAt() *jwtpkg.NumericDate {
	timenow := TimenowInTimezone()

	return jwtpkg.NewNumericDate(timenow.Add(jwt.expiresIn))
}

// parseTokenString 使用 jwtpkg.ParseWithClaims 解析 Token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &JWTCustomClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.signKey, nil
	})
}

// getTokenFromHeader 使用 jwtpkg.ParseWithClaims 解析 Token
// Authorization:Bearer xxxxx
func (jwt *JWT) getTokenFromHeader(c *gin.Context) (string, error) {
	authHeader := c.Request.Header.Get("Authorization")
	if authHeader == "" {
		return "", ErrHeaderEmpty
	}
	// 按空格分割
	parts := strings.SplitN(authHeader, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", ErrHeaderMalformed
	}
	return parts[1], nil
}

// TimenowInTimezone 获取当前时间，支持时区
func TimenowInTimezone() time.Time {
	chinaTimezone, _ := time.LoadLocation("Asia/Shanghai")
	return time.Now().In(chinaTimezone)
}
