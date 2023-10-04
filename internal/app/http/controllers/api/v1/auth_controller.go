package v1

import (
	"yafgo/yafgo-layout/internal/app/http/requests"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/model"
	"yafgo/yafgo-layout/pkg/database"

	"github.com/gin-gonic/gin"
)

type AuthController struct {
	BaseAPIController
}

// LoginByUsername
//
//	@Summary		用户名登录
//	@Description	用户名登录
//	@Tags			Auth
//	@Param			data	body		requests.ReqUsername	true	"请求参数"
//	@Success		200		{object}	any						"{"code": 200, "data": [...]}"
//	@Router			/v1/auth/login/username [post]
//	@Security		ApiToken
func (ctrl *AuthController) LoginByUsername(c *gin.Context) {

	reqData := requests.ReqUsername{}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		ctrl.ParamError(c, err, "请求参数错误")
		return
	}

	ctrl.JSON(c, gin.H{
		"data": reqData,
	})
}

// RegisterByUsername
//
//	@Summary		用户名注册
//	@Description	用户名注册
//	@Tags			Auth
//	@Param			data	body		requests.ReqUsernameRegister	true	"请求参数"
//	@Success		200		{object}	any								"{"code": 200, "data": [...]}"
//	@Router			/v1/auth/register/username [post]
//	@Security		ApiToken
func (ctrl *AuthController) RegisterByUsername(c *gin.Context) {

	reqData := requests.ReqUsernameRegister{}
	if err := c.ShouldBindJSON(&reqData); err != nil {
		ctrl.ParamError(c, err, "请求参数错误")
		return
	}

	user := model.User{
		Username: reqData.Username,
		Password: reqData.Password,
	}
	userDo := g.Query().User.WithContext(c)
	err := userDo.Create(&user)
	if database.IsErrDuplicateEntryCode(err) {
		ctrl.Error(c, err, "用户名已存在")
		return
	}

	ctrl.JSON(c, gin.H{
		"data": user,
		"err":  err,
	})
}
