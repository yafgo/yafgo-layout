package requests

// ReqUsername 用户名登录
type ReqUsername struct {
	Username   string `json:"username,omitempty" binding:"required"`
	Password   string `json:"password,omitempty" binding:"required"`
	VerifyCode string `json:"verify_code,omitempty"`
}

// ReqUsernameRegister 用户名注册
type ReqUsernameRegister struct {
	Username   string `json:"username,omitempty" binding:"required"`
	Password   string `json:"password,omitempty" binding:"required"`
	VerifyCode string `json:"verify_code,omitempty" binding:"required"`
}
