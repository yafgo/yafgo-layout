package main

import (
	"yafgo/yafgo-layout/internal/cmd/yorm"
)

func main() {
	app := yorm.NewApp()
	app.Command().Execute()
}
