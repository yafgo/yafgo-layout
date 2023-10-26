package response

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ApiResponse struct {
	debug bool
	resp  response
}

type response struct {
	Success bool   `json:"success"`
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data,omitempty"`

	ErrData   any `json:"error,omitempty"`
	DebugData any `json:"debug,omitempty"`
}

// New 获取ApiResponse实例
func New(debug ...bool) *ApiResponse {
	resp := new(ApiResponse)
	if len(debug) > 0 {
		resp.debug = debug[0]
	}
	return resp
}

func (p *ApiResponse) beforeOutput(ctx *gin.Context) {
	// beforeOutput
}

func (p *ApiResponse) Success(ctx *gin.Context, data ...any) {
	p.SuccessWithMsg(ctx, "", data...)
}

func (p *ApiResponse) SuccessWithMsg(ctx *gin.Context, msg string, data ...any) {
	p.beforeOutput(ctx)
	p.resp.Msg = msg
	if p.resp.Msg == "" {
		p.resp.Msg = "success"
	}
	if len(data) > 0 {
		p.resp.Data = data[0]
	}

	p.resp.Success = true
	ctx.JSON(http.StatusOK, p.resp)
}

func (p *ApiResponse) Error(ctx *gin.Context, err ...error) {
	p.ErrorWithMsg(ctx, "", err...)
}

func (p *ApiResponse) ErrorWithMsg(ctx *gin.Context, msg string, err ...error) {
	p.beforeOutput(ctx)
	p.resp.Msg = msg
	if len(err) > 0 && err[0] != nil {
		p.resp.ErrData = err[0].Error()
	}

	p.resp.Success = false
	ctx.JSON(http.StatusOK, p.resp)
}

func (p *ApiResponse) WithCode(code int) *ApiResponse {
	p.resp.Code = code
	return p
}

func (p *ApiResponse) WithData(data any) *ApiResponse {
	p.resp.Data = data
	return p
}

func (p *ApiResponse) WithMsg(msg string) *ApiResponse {
	p.resp.Msg = msg
	return p
}

func (p *ApiResponse) WithDebug(data any) *ApiResponse {
	if p.debug {
		p.resp.DebugData = data
	}
	return p
}
