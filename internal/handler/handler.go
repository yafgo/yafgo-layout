package handler

import (
	"errors"
	"net/http"
	"strings"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/service"
	"yafgo/yafgo-layout/pkg/jwtutil"
	"yafgo/yafgo-layout/pkg/response"
	"yafgo/yafgo-layout/pkg/sys/ylog"
	"yafgo/yafgo-layout/pkg/validators"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *ylog.Logger
	g      *g.GlobalObj
	jwt    *jwtutil.JwtUtil

	svcUser service.UserService
}

func NewHandler(
	logger *ylog.Logger,
	g *g.GlobalObj,
	jwt *jwtutil.JwtUtil,
	svcUser service.UserService,
) *Handler {
	return &Handler{
		logger: logger,
		g:      g,
		jwt:    jwt,

		svcUser: svcUser,
	}
}

// JSON 自定义返回的json结构
func (h *Handler) JSON(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}

// Resp 获取 API 响应处理实例
func (h *Handler) Resp() *response.ApiResponse {
	resp := response.New(h.g.IsDev())
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
