// Code generated by swaggo/swag. DO NOT EDIT.

package docs

import "github.com/swaggo/swag"

const docTemplate = `{
    "schemes": {{ marshal .Schemes }},
    "swagger": "2.0",
    "info": {
        "description": "{{escape .Description}}",
        "title": "{{.Title}}",
        "contact": {},
        "license": {
            "name": "MIT",
            "url": "https://github.com/yafgo/yafgo/blob/main/LICENSE"
        },
        "version": "{{.Version}}"
    },
    "host": "{{.Host}}",
    "basePath": "{{.BasePath}}",
    "paths": {
        "/": {
            "get": {
                "security": [
                    {
                        "ApiToken": []
                    }
                ],
                "description": "Api Root",
                "tags": [
                    "API"
                ],
                "summary": "ApiRoot",
                "responses": {
                    "200": {
                        "description": "{\"code\": 200, \"data\": [...]}",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/v1/": {
            "get": {
                "security": [
                    {
                        "ApiToken": []
                    }
                ],
                "description": "Api Index Demo",
                "tags": [
                    "API"
                ],
                "summary": "ApiIndex",
                "responses": {
                    "200": {
                        "description": "{\"code\": 200, \"data\": [...]}",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/v1/auth/login/username": {
            "post": {
                "security": [
                    {
                        "ApiToken": []
                    }
                ],
                "description": "用户名登录",
                "tags": [
                    "Auth"
                ],
                "summary": "用户名登录",
                "parameters": [
                    {
                        "description": "请求参数",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.ReqLoginUsername"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 200, \"data\": [...]}",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        },
        "/v1/auth/register/username": {
            "post": {
                "security": [
                    {
                        "ApiToken": []
                    }
                ],
                "description": "用户名注册",
                "tags": [
                    "Auth"
                ],
                "summary": "用户名注册",
                "parameters": [
                    {
                        "description": "请求参数",
                        "name": "data",
                        "in": "body",
                        "required": true,
                        "schema": {
                            "$ref": "#/definitions/service.ReqRegisterUsername"
                        }
                    }
                ],
                "responses": {
                    "200": {
                        "description": "{\"code\": 200, \"data\": [...]}",
                        "schema": {
                            "type": "object"
                        }
                    }
                }
            }
        }
    },
    "definitions": {
        "service.ReqLoginUsername": {
            "type": "object",
            "required": [
                "password",
                "username"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                }
            }
        },
        "service.ReqRegisterUsername": {
            "type": "object",
            "required": [
                "password",
                "username",
                "verify_code"
            ],
            "properties": {
                "password": {
                    "type": "string"
                },
                "username": {
                    "type": "string"
                },
                "verify_code": {
                    "type": "string"
                }
            }
        }
    },
    "securityDefinitions": {
        "ApiToken": {
            "description": "接口请求token, 格式: ` + "`" + `Bearer {token}` + "`" + `",
            "type": "apiKey",
            "name": "Authorization",
            "in": "header"
        }
    },
    "tags": [
        {
            "description": "未分组接口",
            "name": "API"
        },
        {
            "description": "登录相关接口",
            "name": "Auth"
        },
        {
            "description": "设备租赁相关",
            "name": "设备租赁"
        },
        {
            "description": "“我的”相关接口",
            "name": "我的"
        },
        {
            "description": "后台管理相关接口",
            "name": "后台"
        }
    ]
}`

// SwaggerInfo holds exported Swagger Info so clients can modify it
var SwaggerInfo = &swag.Spec{
	Version:          "1.0.0",
	Host:             "",
	BasePath:         "/api",
	Schemes:          []string{"http", "https"},
	Title:            "YAFGO API",
	Description:      "基于 `Gin` 的 golang 项目模板\n- 本页面可以很方便的调试接口，并不需要再手动复制到 postman 之类的工具中\n- 大部分接口需要登录态，可以手动拿到 `登录token`，点击 `Authorize` 按钮，填入 `Bearer {token}` 并保存即可\n- 接口 url 注意看清楚，要加上 `Base URL` 前缀",
	InfoInstanceName: "swagger",
	SwaggerTemplate:  docTemplate,
	LeftDelim:        "{{",
	RightDelim:       "}}",
}

func init() {
	swag.Register(SwaggerInfo.InstanceName(), SwaggerInfo)
}
