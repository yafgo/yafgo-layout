// Package jwtutil 处理 JWT 认证
package jwtutil

import (
	"context"
	"errors"
	"strings"
	"time"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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

// JwtUtil 自定义一个jwt对象
type JwtUtil struct {
	// 秘钥，用以加密 JWT
	signKey []byte

	// Token 过期时间
	expiresIn time.Duration

	// 刷新 Token 的最大有效时间(即token过期时间在该时间内就可以进行刷新)
	maxRefresh time.Duration

	// jwt RegisteredClaims
	issuer   string
	subject  string
	audience string
}

// _FinalClaims 自定义载荷
type _FinalClaims struct {
	// RegisteredClaims 结构体实现了 Claims 接口继承了  Valid() 方法
	// JWT 规定了7个官方字段，提供使用:
	// - iss (issuer)：发布者
	// - sub (subject)：主题
	// - iat (Issued At)：生成签名的时间
	// - exp (expiration time)：签名过期时间
	// - aud (audience)：观众，相当于接受者
	// - nbf (Not Before)：生效时间
	// - jti (JWT ID)：编号
	jwt.RegisteredClaims

	// 自定义载荷
	CustomClaims
}

// ParserToken 解析 Token
func (ju *JwtUtil) ParserToken(tokenStr string) (*CustomClaims, error) {

	// 1. 解析用户传参的 Token
	token, err := ju.parseTokenString(tokenStr)

	// 2. 解析出错
	if err != nil {
		switch err {
		case jwt.ErrTokenExpired:
			return nil, ErrTokenExpired
		default:
			return nil, ErrTokenInvalid
		}
	}

	// 3. 将 token 中的 finalClaims 信息解析出来和 JWTCustomClaims 数据结构进行校验
	if finalClaims, ok := token.Claims.(*_FinalClaims); ok && token.Valid {
		customClaims := finalClaims.CustomClaims
		return &customClaims, nil
	}

	return nil, ErrTokenInvalid
}

// ParserTokenFromHeader 从请求头直接解析 Token
func (ju *JwtUtil) ParserTokenFromHeader(c *gin.Context) (*CustomClaims, error) {
	tokenStr, err := ju.GetTokenFromHeader(c)
	if err != nil {
		return nil, err
	}

	return ju.ParserToken(tokenStr)
}

// RefreshToken 更新 Token，用以提供 refresh token 接口
func (ju *JwtUtil) RefreshToken(tokenOld string) (tokenNew string, err error) {

	// 1. 解析用户传参的 Token
	token, err := ju.parseTokenString(tokenOld)

	// 2. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		// 满足 refresh 的条件：只是单一的报错token过期 ErrTokenExpired
		if errors.Is(err, jwt.ErrTokenExpired) {
			return "", err
		}
	}

	// 3. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*_FinalClaims)

	// 4. 检查是否过了『最大允许刷新的时间』
	x := TimenowInTimezone().Add(-ju.maxRefresh)
	if claims.IssuedAt.After(x) {
		// 修改过期时间
		claims.RegisteredClaims.ExpiresAt = ju.newExpiresAt()
		return ju.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

// RefreshTokenFromHeader 从请求头直接解析并刷新 Token
func (ju *JwtUtil) RefreshTokenFromHeader(c *gin.Context) (tokenNew string, err error) {
	tokenOld, err := ju.GetTokenFromHeader(c)
	if err != nil {
		return "", err
	}

	return ju.RefreshToken(tokenOld)
}

// IssueToken 颁发Token，一般在登录成功时调用
func (ju *JwtUtil) IssueToken(claims CustomClaims) (token string, err error) {

	// 1. 构造用户 claims 信息(负荷)
	expiresAt := ju.newExpiresAt()
	issuedAt := jwt.NewNumericDate(TimenowInTimezone())

	finalClaims := _FinalClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    ju.issuer,          // 签名颁发者
			Subject:   "",                 //
			Audience:  jwt.ClaimStrings{}, //
			ExpiresAt: expiresAt,          // 签名过期时间
			NotBefore: issuedAt,           // 签名生效时间
			IssuedAt:  issuedAt,           // 首次签名时间（后续刷新 Token 不会更新该字段）
		},
		// 自定义载荷
		CustomClaims: claims,
	}

	// 2. 根据 claims 生成token对象
	token, err = ju.createToken(finalClaims)
	if err != nil {
		ylog.Error(context.Background(), err)
	}

	return
}

// createToken 创建 Token，内部使用，外部请调用 IssueToken
func (ju *JwtUtil) createToken(claims _FinalClaims) (string, error) {
	// 使用HS256算法进行t生成
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString(ju.signKey)
}

// newExpiresAt 过期时间
func (ju *JwtUtil) newExpiresAt() *jwt.NumericDate {
	timenow := TimenowInTimezone()

	return jwt.NewNumericDate(timenow.Add(ju.expiresIn))
}

// parseTokenString 使用 jwt.ParseWithClaims 解析 Token
func (ju *JwtUtil) parseTokenString(tokenString string) (*jwt.Token, error) {
	return jwt.ParseWithClaims(tokenString, &_FinalClaims{}, func(token *jwt.Token) (interface{}, error) {
		return ju.signKey, nil
	})
}

// GetTokenFromHeader 从请求头获取 jwtToken 字符串
//
//	请求头示例: "Authorization:Bearer {jwtToken字符串}"
func (ju *JwtUtil) GetTokenFromHeader(c *gin.Context) (string, error) {
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
