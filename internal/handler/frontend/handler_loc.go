package frontend

import (
	"yafgo/yafgo-layout/internal/handler"

	"github.com/gin-gonic/gin"
)

type LocHandler interface {
	List(ctx *gin.Context)
	Create(ctx *gin.Context)
	Get(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
}

func NewLocHandler(
	handler *handler.Handler,
) LocHandler {
	return &locHandler{
		Handler: handler,
	}
}

type locHandler struct {
	*handler.Handler
}

// List implements LocHandler.
//
//	@Summary	Loc list
//	@Description
//	@Tags		API
//	@Success	200	{object}	any	"{"code": 200, "data": [...]}"
//	@Security	ApiToken
//	@Router		/v1/loc/locs [get]
func (h *locHandler) List(ctx *gin.Context) {
	h.Resp().Success(ctx, gin.H{
		"data": "/api/",
	})
}

// Get implements LocHandler.
//
//	@Summary	Loc 查询单条
//	@Description
//	@Tags		API
//	@Param		id	path		int	true	"id"
//	@Success	200	{object}	any	"{"code": 200, "data": [...]}"
//	@Security	ApiToken
//	@Router		/v1/loc/locs/{id} [get]
func (h *locHandler) Get(ctx *gin.Context) {
	id := ctx.Param("id")

	h.Resp().Success(ctx, gin.H{
		"data": id,
	})
}

// Create implements LocHandler.
//
//	@Summary	Loc 新增
//	@Description
//	@Tags		API
//	@Success	200	{object}	any	"{"code": 200, "data": [...]}"
//	@Security	ApiToken
//	@Router		/v1/loc/locs [post]
func (h *locHandler) Create(ctx *gin.Context) {
	h.Resp().Success(ctx, gin.H{
		"data": "/api/",
	})
}

// Update implements LocHandler.
//
//	@Summary	Loc 更新
//	@Description
//	@Tags		API
//	@Param		id	path		int	true	"id"
//	@Success	200	{object}	any	"{"code": 200, "data": [...]}"
//	@Security	ApiToken
//	@Router		/v1/loc/locs/{id} [post]
func (h *locHandler) Update(ctx *gin.Context) {
	id := ctx.Param("id")

	h.Resp().Success(ctx, gin.H{
		"data": id,
	})
}

// Delete implements LocHandler.
//
//	@Summary	Loc 删除
//	@Description
//	@Tags		API
//	@Param		id	path		int	true	"id"
//	@Success	200	{object}	any	"{"code": 200, "data": [...]}"
//	@Security	ApiToken
//	@Router		/v1/loc/locs/{id} [delete]
func (h *locHandler) Delete(ctx *gin.Context) {
	id := ctx.Param("id")

	h.Resp().Success(ctx, gin.H{
		"data": id,
	})
}
