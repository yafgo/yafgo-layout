package handler

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/model"
	"yafgo/yafgo-layout/internal/service"
	"yafgo/yafgo-layout/pkg/jwtutil"
	"yafgo/yafgo-layout/pkg/response"
	"yafgo/yafgo-layout/pkg/sys/ylog"
	"yafgo/yafgo-layout/pkg/validators"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Logger *ylog.Logger
	G      *g.GlobalObj
	Jwt    *jwtutil.JwtUtil

	SvcUser service.UserService
}

func NewHandler(
	logger *ylog.Logger,
	g *g.GlobalObj,
	jwt *jwtutil.JwtUtil,
	svcUser service.UserService,
) *Handler {
	return &Handler{
		Logger: logger,
		G:      g,
		Jwt:    jwt,

		SvcUser: svcUser,
	}
}

// JSON 自定义返回的json结构
func (h *Handler) JSON(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}

// Resp 获取 API 响应处理实例
func (h *Handler) Resp() *response.ApiResponse {
	resp := response.New(h.G.IsDev())
	return resp
}

// ParamError 参数错误
func (h *Handler) ParamError(ctx *gin.Context, err error, msg ...string) {
	errMsgs := validators.TranslateErrors(err)
	if errMsgs != nil {
		var sb strings.Builder
		for _, v := range errMsgs {
			sb.WriteString(v + ";")
		}
		err = errors.New(strings.TrimSuffix(sb.String(), ";"))
	}

	resp := h.Resp()
	if len(msg) > 0 {
		resp.WithMsg(msg[0])
	} else if err != nil {
		resp.WithMsg(err.Error())
	}
	resp.Error(ctx, err)
}

// JwtClaims 从 gin.context 中获取当前 jwtClaims
func (h *Handler) JwtClaims(ctx *gin.Context) (cc *jwtutil.CustomClaims) {
	cc = new(jwtutil.CustomClaims)
	claims, ok := ctx.Get("claims")
	if !ok {
		return
	}
	cc = claims.(*jwtutil.CustomClaims)
	return
}

// CurrentUserID 从 gin.context 中获取当前登录用户 ID
func (h *Handler) CurrentUserID(ctx *gin.Context) int64 {
	claims := h.JwtClaims(ctx)
	if claims == nil {
		return 0
	}
	return claims.UserID
}

// CurrentUserIDStr 从 gin.context 中获取当前登录用户 ID
func (h *Handler) CurrentUserIDStr(ctx *gin.Context) string {
	return strconv.Itoa(int(h.CurrentUserID(ctx)))
}

// 根据请求 jwt 信息从 db 获取当前登录用户
func (h *Handler) CurrentUser(ctx *gin.Context) (u *model.User, err error) {
	uid := h.CurrentUserID(ctx)
	if uid == 0 {
		h.Logger.Errorf(ctx, "从Gin的Context中获取从jwt解析出来的用户ID失败, 请检查路由是否使用jwt中间件")
		return
	}
	return h.SvcUser.GetByID(ctx, uid)
}

// IsLogin 当前是否已登录
func (h *Handler) IsLogin(ctx *gin.Context) bool {
	uid := h.CurrentUserID(ctx)
	return uid > 0
}
