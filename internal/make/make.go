package make

import (
	"embed"
	"log"
	"strings"
	"yafgo/yafgo-layout/pkg/helper/file"
	"yafgo/yafgo-layout/pkg/helper/str"

	"github.com/gookit/color"
	"github.com/iancoleman/strcase"
	"github.com/spf13/cobra"
)

type Model struct {
	SnakeName            string
	SnakeNamePlural      string
	CamelName            string
	CamelNamePlural      string
	LowerCamelName       string
	LowerCamelNamePlural string
}

// stubsFS 方便我们后面打包 stubs 目录下的模板文件

//go:embed stubs
var stubsFS embed.FS

var CmdMake = &cobra.Command{
	Use:   "make",
	Short: "Generate file and code",
}

func init() {
	// 注册 make 的子命令
	CmdMake.AddCommand(
		cmdMakeHandler,
		cmdMakeRepository,
		cmdMakeService,
	)
}

// makeModelFromString 格式化用户输入的内容
func makeModelFromString(name string) Model {
	model := Model{}
	model.CamelName = str.Singular(strcase.ToCamel(name))
	model.CamelNamePlural = str.Plural(model.CamelName)
	model.SnakeNamePlural = str.Snake(model.CamelNamePlural)
	model.LowerCamelName = str.LowerCamel(model.CamelName)
	model.SnakeName = str.Snake(model.CamelName)
	model.LowerCamelNamePlural = str.LowerCamel(model.CamelNamePlural)
	return model
}

// createFileFromStub 读取 stub 文件并进行变量替换
// 最后一个选项可选，如若传参，应传 map[string]string 类型，作为附加的变量搜索替换
func createFileFromStub(filePath string, stubName string, model Model, variables ...any) {

	// 实现最后一个参数可选
	replaces := make(map[string]string)
	if len(variables) > 0 {
		replaces = variables[0].(map[string]string)
	}

	// 目标文件已存在
	if file.Exists(filePath) {
		log.Fatalln(color.Red.Render(filePath + " already exists!"))
	}

	// 读取 stub 模板文件
	modelData, err := stubsFS.ReadFile("stubs/" + stubName + ".gotpl")
	if err != nil {
		log.Fatalln(color.Red.Render(err.Error()))
	}
	modelStub := string(modelData)

	// 添加默认的替换变量
	replaces["{{LowerCamelName}}"] = model.LowerCamelName
	replaces["{{LowerCamelNamePlural}}"] = model.LowerCamelNamePlural
	replaces["{{CamelName}}"] = model.CamelName
	replaces["{{CamelNamePlural}}"] = model.CamelNamePlural
	replaces["{{SnakeName}}"] = model.SnakeName
	replaces["{{SnakeNamePlural}}"] = model.SnakeNamePlural

	// 对模板内容做变量替换
	for search, replace := range replaces {
		modelStub = strings.ReplaceAll(modelStub, search, replace)
	}

	// 存储到目标文件中
	err = file.PutContent(filePath, []byte(modelStub))
	if err != nil {
		log.Fatalln(color.Red.Render(err))
	}

	// 提示成功
	log.Println(color.Success.Sprintf("[%s] created.", filePath))
}
