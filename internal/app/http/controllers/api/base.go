package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// BaseAPIController 基础控制器
type BaseAPIController struct {
}

// JSON 自定义返回的json结构
func (ctrl *BaseAPIController) JSON(ctx *gin.Context, data any) {
	ctx.JSON(http.StatusOK, data)
}

// JSON 自定义返回的json结构
func (ctrl *BaseAPIController) Debug(ctx *gin.Context, data any) {
	ctx.Set("debugData", data)
}

// 种cookie: document.cookie='debug=2k9j38h#4'
func (ctrl *BaseAPIController) getDebugData(ctx *gin.Context) (data any) {
	if ck, _ := ctx.Cookie("debug"); ck != "2k9j38h#4" {
		return nil
	}
	if val, ok := ctx.Get("debugData"); ok {
		return val
	}
	return nil
}

func (ctrl *BaseAPIController) SuccessWithMsg(ctx *gin.Context, msg string, data ...any) {
	respData := gin.H{
		"success": true,
		"message": msg,
		"data":    nil,
	}
	if len(data) > 0 {
		respData["data"] = data[0]
	}
	if debugData := ctrl.getDebugData(ctx); debugData != nil {
		respData["debug"] = debugData
	}
	ctx.JSON(http.StatusOK, respData)
}

func (ctrl *BaseAPIController) Success(ctx *gin.Context, data ...any) {
	ctrl.SuccessWithMsg(ctx, "操作成功!", data...)
}

// Error 传参 err 对象，未传参 msg 时使用默认消息
// 在解析用户请求，请求的格式或者方法不符合预期时调用
func (ctrl *BaseAPIController) Error(ctx *gin.Context, err error, msg ...string) {
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
	if debugData := ctrl.getDebugData(ctx); debugData != nil {
		respData["debug"] = debugData
	}
	// ctx.AbortWithStatusJSON(http.StatusBadRequest, respData)
	ctx.JSON(http.StatusOK, respData)
}

// Error 传参 err 对象，未传参 msg 时使用默认消息
// 在解析用户请求，请求的格式或者方法不符合预期时调用
func (ctrl *BaseAPIController) ErrorWithData(ctx *gin.Context, err error, msg string, data any) {
	respData := gin.H{
		"success": false,
		"message": msg,
		"data":    data,
	}

	if err != nil {
		respData["error"] = err.Error()
	}
	if debugData := ctrl.getDebugData(ctx); debugData != nil {
		respData["debug"] = debugData
	}
	ctx.JSON(http.StatusOK, respData)
}
