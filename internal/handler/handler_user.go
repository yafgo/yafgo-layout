package handler

import (
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
	*Handler
	userService service.UserService

	ju *jwtutil.JwtUtil
}

func NewUserHandler(
	handler *Handler,
	userService service.UserService,
	ju *jwtutil.JwtUtil,
) UserHandler {
	return &userHandler{
		Handler:     handler,
		userService: userService,

		ju: ju,
	}
}

// RegisterByUsername implements UserHandler.
//
//	@Summary		用户名注册
//	@Description	用户名注册
//	@Tags			Auth
//	@Param			data	body		requests.ReqUsernameRegister	true	"请求参数"
//	@Success		200		{object}	any								"{"code": 200, "data": [...]}"
//	@Router			/v1/auth/register/username [post]
//	@Security		ApiToken
func (h *userHandler) RegisterByUsername(ctx *gin.Context) {

	reqData := new(service.ReqRegisterUsername)
	if err := ctx.ShouldBindJSON(&reqData); err != nil {
		h.ParamError(ctx, err, "请求参数错误")
		return
	}

	user, err := h.userService.RegisterByUsername(ctx, reqData)
	if err != nil {
		h.Error(ctx, err, "注册失败")
		return
	}

	// 颁发jwtToken
	token, err := h.ju.IssueToken(jwtutil.CustomClaims{UserID: user.ID})
	if err != nil {
		h.Error(ctx, err, "生成token失败")
		return
	}

	h.SuccessWithMsg(ctx, "注册成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}

// LoginByUsername implements UserHandler.
func (h *userHandler) LoginByUsername(ctx *gin.Context) {
	panic("unimplemented")
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
