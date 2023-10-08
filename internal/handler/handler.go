package handler

import (
	"errors"
	"net/http"
	"strings"
	"yafgo/yafgo-layout/pkg/sys/ylog"
	"yafgo/yafgo-layout/pkg/validators"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	logger *ylog.Logger
}

func NewHandler(logger *ylog.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}

// JSON 自定义返回的json结构
func (h *Handler) JSON(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}

// JSON 自定义返回的json结构
func (h *Handler) Debug(ctx *gin.Context, data any) {
	ctx.Set("debugData", data)
}

// 种cookie: document.cookie='debug=2k9j38h#4'
func (h *Handler) getDebugData(ctx *gin.Context) (data any) {
	if ck, _ := ctx.Cookie("debug"); ck != "2k9j38h#4" {
		return nil
	}
	if val, ok := ctx.Get("debugData"); ok {
		return val
	}
	return nil
}

func (h *Handler) Success(ctx *gin.Context, data ...any) {
	h.SuccessWithMsg(ctx, "操作成功!", data...)
}

func (h *Handler) SuccessWithMsg(ctx *gin.Context, msg string, data ...any) {
	respData := gin.H{
		"success": true,
		"message": msg,
		"data":    nil,
	}
	if len(data) > 0 {
		respData["data"] = data[0]
	}
	if debugData := h.getDebugData(ctx); debugData != nil {
		respData["debug"] = debugData
	}
	ctx.JSON(http.StatusOK, respData)
}

// Error 传参 err 对象，未传参 msg 时使用默认消息
// 在解析用户请求，请求的格式或者方法不符合预期时调用
func (h *Handler) Error(ctx *gin.Context, err error, msg ...string) {
	var _msg string
	if len(msg) > 0 {
		_msg = msg[0]
	}
	respData := gin.H{
		"success": false,
		"message": _msg,
	}

	if err != nil {
		respData["error"] = err.Error()
	}
	if debugData := h.getDebugData(ctx); debugData != nil {
		respData["debug"] = debugData
	}
	// ctx.AbortWithStatusJSON(http.StatusBadRequest, respData)
	ctx.JSON(http.StatusOK, respData)
}

// Error 传参 err 对象，未传参 msg 时使用默认消息
// 在解析用户请求，请求的格式或者方法不符合预期时调用
func (h *Handler) ErrorWithData(ctx *gin.Context, err error, msg string, data any) {
	respData := gin.H{
		"success": false,
		"message": msg,
		"data":    data,
	}

	if err != nil {
		respData["error"] = err.Error()
	}
	if debugData := h.getDebugData(ctx); debugData != nil {
		respData["debug"] = debugData
	}
	ctx.JSON(http.StatusOK, respData)
}

// ParamError 参数错误
func (h *Handler) ParamError(ctx *gin.Context, err error, msg ...string) {
	errMsgs := validators.TranslateErrors(err)
	if errMsgs != nil {
		var sb strings.Builder
		for _, v := range errMsgs {
			sb.WriteString(v + "\n")
		}
		err = errors.New(strings.TrimSuffix(sb.String(), "\n"))
	}

	h.Error(ctx, err, msg...)
}
