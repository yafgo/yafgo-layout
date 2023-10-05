// Package jwt 处理 JWT 认证
package jwt

import (
	"context"
	"errors"
	"strings"
	"time"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"github.com/gin-gonic/gin"
	jwtpkg "github.com/golang-jwt/jwt/v5"
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

// JWT 定义一个jwt对象
type JWT struct {
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
	jwtpkg.RegisteredClaims

	// 自定义载荷
	CustomClaims
}

// ParserToken 解析 Token
func (jwt *JWT) ParserToken(tokenStr string) (*CustomClaims, error) {

	// 1. 解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenStr)

	// 2. 解析出错
	if err != nil {
		switch err {
		case jwtpkg.ErrTokenExpired:
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
func (jwt *JWT) ParserTokenFromHeader(c *gin.Context) (*CustomClaims, error) {
	tokenStr, err := jwt.GetTokenFromHeader(c)
	if err != nil {
		return nil, err
	}

	return jwt.ParserToken(tokenStr)
}

// RefreshToken 更新 Token，用以提供 refresh token 接口
func (jwt *JWT) RefreshToken(tokenOld string) (tokenNew string, err error) {

	// 1. 解析用户传参的 Token
	token, err := jwt.parseTokenString(tokenOld)

	// 2. 解析出错，未报错证明是合法的 Token（甚至未到过期时间）
	if err != nil {
		// 满足 refresh 的条件：只是单一的报错token过期 ErrTokenExpired
		if errors.Is(err, jwtpkg.ErrTokenExpired) {
			return "", err
		}
	}

	// 3. 解析 JWTCustomClaims 的数据
	claims := token.Claims.(*_FinalClaims)

	// 4. 检查是否过了『最大允许刷新的时间』
	x := TimenowInTimezone().Add(-jwt.maxRefresh)
	if claims.IssuedAt.After(x) {
		// 修改过期时间
		claims.RegisteredClaims.ExpiresAt = jwt.newExpiresAt()
		return jwt.createToken(*claims)
	}

	return "", ErrTokenExpiredMaxRefresh
}

// RefreshTokenFromHeader 从请求头直接解析并刷新 Token
func (jwt *JWT) RefreshTokenFromHeader(c *gin.Context) (tokenNew string, err error) {
	tokenOld, err := jwt.GetTokenFromHeader(c)
	if err != nil {
		return "", err
	}

	return jwt.RefreshToken(tokenOld)
}

// IssueToken 颁发Token，一般在登录成功时调用
func (jwt *JWT) IssueToken(claims CustomClaims) (token string, err error) {

	// 1. 构造用户 claims 信息(负荷)
	expiresAt := jwt.newExpiresAt()
	issuedAt := jwtpkg.NewNumericDate(TimenowInTimezone())

	finalClaims := _FinalClaims{
		RegisteredClaims: jwtpkg.RegisteredClaims{
			Issuer:    jwt.issuer,            // 签名颁发者
			Subject:   "",                    //
			Audience:  jwtpkg.ClaimStrings{}, //
			ExpiresAt: expiresAt,             // 签名过期时间
			NotBefore: issuedAt,              // 签名生效时间
			IssuedAt:  issuedAt,              // 首次签名时间（后续刷新 Token 不会更新该字段）
		},
		// 自定义载荷
		CustomClaims: claims,
	}

	// 2. 根据 claims 生成token对象
	token, err = jwt.createToken(finalClaims)
	if err != nil {
		ylog.Error(context.Background(), err)
	}

	return
}

// createToken 创建 Token，内部使用，外部请调用 IssueToken
func (jwt *JWT) createToken(claims _FinalClaims) (string, error) {
	// 使用HS256算法进行t生成
	t := jwtpkg.NewWithClaims(jwtpkg.SigningMethodHS256, claims)
	return t.SignedString(jwt.signKey)
}

// newExpiresAt 过期时间
func (jwt *JWT) newExpiresAt() *jwtpkg.NumericDate {
	timenow := TimenowInTimezone()

	return jwtpkg.NewNumericDate(timenow.Add(jwt.expiresIn))
}

// parseTokenString 使用 jwtpkg.ParseWithClaims 解析 Token
func (jwt *JWT) parseTokenString(tokenString string) (*jwtpkg.Token, error) {
	return jwtpkg.ParseWithClaims(tokenString, &_FinalClaims{}, func(token *jwtpkg.Token) (interface{}, error) {
		return jwt.signKey, nil
	})
}

// GetTokenFromHeader 从请求头获取 jwtToken 字符串
//
//	请求头示例: "Authorization:Bearer {jwtToken字符串}"
func (jwt *JWT) GetTokenFromHeader(c *gin.Context) (string, error) {
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
