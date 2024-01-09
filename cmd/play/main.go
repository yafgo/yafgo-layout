package main

import (
	"yafgo/yafgo-layout/internal/cmd/play"
)

func main() {
	app := play.NewApp().Command()
	app.Execute()
}
