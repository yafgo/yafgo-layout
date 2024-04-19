# Yafgo Basic Layout

> 文档地址: [https://yafgo.pages.dev/](https://yafgo.pages.dev/)

## 功能清单

- [x] make 命令 (开发阶段使用)
  - [x] handler
  - [x] repository
  - [x] service
  - [x] migrations
- [x] play 命令 (开发阶段使用)
- [x] gorm gen 命令 (开发阶段使用)
- [x] migrate 命令
- [x] serve 命令
- [x] swag 文档生成

## 开始使用

### 环境要求

- `git`
- `go` _1.21+_
- `mysql` _5.7+_
- `redis` _6.0+_

### 创建项目

```shell
# 安装
go install github.com/yafgo/yafgo@latest

# 创建新项目
yafgo

# 示例
✔ Project Name: my_project
Use the arrow keys to navigate: ↓ ↑ → ←
Select Template?
  🌶 [Yafgo]    (Yafgo 后端项目模板)
     [YafgoWeb] (Yafgo 前后端项目模板)

# 从模板列表选择一个模板即可
```

### 运行项目

```shell
> ./ycli
[Yafgo-Cli] v1.0.0

Usage:
 ./ycli [command]

Available Commands:
  make      代码生成
  play      代码演练
  orm       生成gorm代码
  migrate   执行db迁移
  doc       更新swagger文档
  serve     启动webServer
```
