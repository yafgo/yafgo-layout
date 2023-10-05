package v1

import (
	"yafgo/yafgo-layout/internal/app/http/requests"
	"yafgo/yafgo-layout/internal/g"
	"yafgo/yafgo-layout/internal/model"
	"yafgo/yafgo-layout/pkg/database"
	"yafgo/yafgo-layout/pkg/hash"
	"yafgo/yafgo-layout/pkg/jwt"

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

	userQ := g.Query().User
	userDo := userQ.WithContext(c)
	user, err := userDo.Where(userQ.Username.Eq(reqData.Username)).First()
	if err != nil {
		ctrl.Error(c, err, "用户名不存在")
		return
	}

	if !hash.BcryptCheck(reqData.Password, user.Password) {
		ctrl.Error(c, nil, "密码不正确")
		return
	}

	// 颁发jwtToken
	token, err := g.Jwt().IssueToken(jwt.CustomClaims{UserID: user.ID})
	if err != nil {
		ctrl.Error(c, err, "生成token失败")
		return
	}

	ctrl.SuccessWithMsg(c, "登录成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
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

	// 颁发jwtToken
	token, err := g.Jwt().IssueToken(jwt.CustomClaims{UserID: user.ID})
	if err != nil {
		ctrl.Error(c, err, "生成token失败")
		return
	}

	ctrl.SuccessWithMsg(c, "注册成功", gin.H{
		"token": token,
		"user": gin.H{
			"id":       user.ID,
			"username": user.Username,
		},
	})
}
