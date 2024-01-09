package main

import (
	"yafgo/yafgo-layout/internal/cmd/ymake"
)

func main() {
	cm := ymake.NewCodeMaker()
	cm.Command().Execute()
}
