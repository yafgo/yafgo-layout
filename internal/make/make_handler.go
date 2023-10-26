package make

import (
	"fmt"

	"github.com/spf13/cobra"
)

var cmdMakeHandler = &cobra.Command{
	Use:   "handler",
	Short: "Create handler, exmaple: make handler user",
	Args:  cobra.ExactArgs(1), // 只允许且必须传 1 个参数
	Run:   runMakeHandler,
}

func runMakeHandler(cmd *cobra.Command, args []string) {

	handlerName := args[0]
	model := makeModelFromString(handlerName)

	// 组建目标目录
	filePath := fmt.Sprintf("internal/handler/handler_%s.go", model.SnakeName)

	// 基于模板创建文件（做好变量替换）
	createFileFromStub(filePath, "handler", model)
}
