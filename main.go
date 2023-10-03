package main

import (
	"yafgo/yafgo-layout/cmd"
)

//go:generate swag fmt
//go:generate swag init --parseInternal --parseDepth=10 -o ./resource/docs -ot "go,json"

// Swagger config
//
//	@title			YAFGO API
//	@version		1.0.0
//	@description	基于 `Gin` 的 golang 项目模板
//	@description	- 本页面可以很方便的调试接口，并不需要再手动复制到 postman 之类的工具中
//	@description	- 大部分接口需要登录态，可以手动拿到 `登录token`，点击 `Authorize` 按钮，填入 `Bearer {token}` 并保存即可
//	@description	- 接口 url 注意看清楚，要加上 `Base URL` 前缀
//	@license.name	MIT
//	@license.url	https://github.com/yafgo/yafgo/blob/main/LICENSE
//
//	@host
//	@BasePath					/api
//	@schemes					http https
//
//	@securityDefinitions.apikey	ApiToken
//	@in							header
//	@name						Authorization
//	@description				接口请求token, 格式: `Bearer {token}`
func main() {
	cmd.RunApp()
}
