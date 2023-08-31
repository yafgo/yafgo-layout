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
	app := new(cmd.Application)
	app.Run()
}
