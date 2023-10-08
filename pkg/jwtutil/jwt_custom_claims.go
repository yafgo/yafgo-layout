package jwtutil

// CustomClaims 自定义载荷
//
//	可以根据具体业务需求在这里进行修改
type CustomClaims struct {
	UserID int64 `json:"user_id"`
}
