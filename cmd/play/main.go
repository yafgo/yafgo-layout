package main

import (
	"log"
	"os"
	"path/filepath"
)

func main() {
	pg := newPlay("dev")
	rootCmd := pg.PlayCommand()
	filepath.Join("../..")

	// 执行 play 主命令
	if err := rootCmd.Execute(); err != nil {
		log.Fatalf("Failed to run app with %v: %+v\n", os.Args, err)
	}
}
