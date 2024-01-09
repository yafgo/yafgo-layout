package frontend

import (
	"yafgo/yafgo-layout/internal/handler"
	"yafgo/yafgo-layout/internal/service"
	"yafgo/yafgo-layout/pkg/jwtutil"

	"github.com/gin-gonic/gin"
)

type UserHandler interface {
	RegisterByUsername(ctx *gin.Context)
	LoginByUsername(ctx *gin.Context)
	Register(ctx *gin.Context)
	Login(ctx *gin.Context)
	GetProfile(ctx *gin.Context)
	UpdateProfile(ctx *gin.Context)
}

type userHandler struct {
	*handler.Handler
}

func NewUserHandler(handler *handler.Handler) UserHandler {
	return &userHandler{
		Handler: handler,
	}
}

// RegisterByUsername implements UserHandler.
//
//	@Summary		用户名注册
//	@Description	用户名注册
//	@Tags			Auth
//	@Param			data	body		service.ReqRegisterUsername	true	"请求参数"
//	@Success		200		{object}	any							"{"code": 200, "data": [...]}"
//	@Router			/v1/auth/register/username [post]
//	@Security		ApiToken
func (h *userHandler) RegisterByUsername(ctx *gin.Context) {

	reqData := new(service.ReqRegisterUsername)
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		h.ParamError(ctx, err, "请求参数错误")
		return
	}

	user, err := h.SvcUser.RegisterByUsername(ctx, reqData)
	if err != nil {
		h.Resp().ErrorWithMsg(ctx, "注册失败", err)
		return
	}

	// 颁发jwtToken
	token, err := h.Jwt.IssueToken(jwtutil.CustomClaims{UserID: user.ID})
	if err != nil {
		h.Resp().ErrorWithMsg(ctx, "生成token失败", err)
		return
	}

	h.Resp().SuccessWithMsg(ctx, "注册成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// LoginByUsername implements UserHandler.
//
//	@Summary		用户名登录
//	@Description	用户名登录
//	@Tags			Auth
//	@Param			data	body		service.ReqLoginUsername	true	"请求参数"
//	@Success		200		{object}	any							"{"code": 200, "data": [...]}"
//	@Router			/v1/auth/login/username [post]
//	@Security		ApiToken
func (h *userHandler) LoginByUsername(ctx *gin.Context) {

	reqData := new(service.ReqLoginUsername)
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		h.ParamError(ctx, err, "请求参数错误")
		return
	}

	user, err := h.SvcUser.LoginByUsername(ctx, reqData)
	if err != nil {
		h.Resp().ErrorWithMsg(ctx, "登录失败", err)
		return
	}

	// 颁发jwtToken
	token, err := h.Jwt.IssueToken(jwtutil.CustomClaims{UserID: user.ID})
	if err != nil {
		h.Resp().ErrorWithMsg(ctx, "生成token失败", err)
		return
	}

	h.Resp().SuccessWithMsg(ctx, "登录成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

func (h *userHandler) Register(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (h *userHandler) Login(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (h *userHandler) GetProfile(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}

func (h *userHandler) UpdateProfile(ctx *gin.Context) {
	panic("not implemented") // TODO: Implement
}
