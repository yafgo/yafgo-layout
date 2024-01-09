// package ymake 用于代码生成
package ymake

import (
	"bytes"
	"embed"
	"log"
	"os"
	"path"
	"path/filepath"
	"strings"
	"text/template"
	"yafgo/yafgo-layout/pkg/helper/file"
	"yafgo/yafgo-layout/pkg/helper/str"

	"github.com/gookit/color"
	"github.com/pkg/errors"
)

// stubsFS 方便我们后面打包这些 .stub 为后缀名的文件
//
//go:embed stubs
var stubsFS embed.FS

// 生成类型
const (
	TypeHandler    = "handler"
	TypeService    = "service"
	TypeRepository = "repo"
)

// codeMaker 代码生成器
type codeMaker struct {
	dirBase  string // handler/service/repository 等目录所在的目录, 默认 "internal"
	makeType string // 生成类型: handler,service,repo
	makeName string // 用户输入的待生成名称, 如: user, frontend/user
	realName string // 处理后的待生成名称, 如: user
	fileDest string // 待生成文件路径

	vars map[string]any // 自定义模板变量
}

func NewCodeMaker(opts ...ymakeOption) (cm *codeMaker) {
	cm = &codeMaker{
		dirBase: "internal",
	}

	for _, opt := range opts {
		opt(cm)
	}
	return cm
}

type ymakeOption func(*codeMaker)

// WithDirBase 设置基础目录
func WithDirBase(dirBase string) ymakeOption {
	return func(cm *codeMaker) {
		cm.dirBase = dirBase
	}
}

type ymakeModel struct {
	SnakeName            string // 下划线
	SnakeNamePlural      string // 下划线[复数]
	CamelName            string // 大驼峰
	CamelNamePlural      string // 大驼峰[复数]
	LowerCamelName       string // 小驼峰
	LowerCamelNamePlural string // 小驼峰[复数]
}

// makeModelFromString 格式化用户输入的内容
func (p *codeMaker) makeModelFromString(name string) ymakeModel {
	model := ymakeModel{}
	camelName := str.Camel(name)
	model.CamelName = str.Singular(camelName)
	model.CamelNamePlural = str.Plural(model.CamelName)
	model.SnakeName = str.Snake(model.CamelName)
	model.SnakeNamePlural = str.Snake(model.CamelNamePlural)
	model.LowerCamelName = str.LowerCamel(model.CamelName)
	model.LowerCamelNamePlural = str.LowerCamel(model.CamelNamePlural)
	return model
}

type ParamMake struct {
	Type  string         // 生成类型: handler,service,repository
	Model string         // 待生成的结构体名称
	Vars  map[string]any // 自定义模板变量
}

// Make 读取模板 stub 文件并进行变量替换, 生成所需文件
func (p *codeMaker) Make(makeType string, args []string) error {
	if len(args) < 1 {
		return errors.New("缺少参数: 待生成实体名称")
	}
	p.makeType = makeType // 生成类型: handler,service,repo
	p.makeName = args[0]  // 生成实体名

	// 计算目标文件路径
	p.calcFileDest()
	// 目标文件已存在
	if file.Exists(p.fileDest) {
		return errors.New(colorWarn("文件已存在: ") + p.fileDest)
	}

	// 读取模板并生成文件
	err := p.renderTemplate(p.makeType, p.fileDest, p.vars)
	if err != nil {
		return err
	}

	// handler目录特殊处理, 首次创建目录时额外创建 routes.go
	if p.makeType == TypeHandler {
		var dir = path.Dir(p.fileDest)
		var fileRoute = path.Join(dir, "routes.go")
		if !file.Exists(fileRoute) {
			err := p.renderTemplate("routes", fileRoute, p.vars)
			if err != nil {
				return errors.Wrap(err, colorError("生成路由文件出错: ")+fileRoute)
			}
		}
	}

	return nil
}

// renderTemplate 渲染模板并生成渲染后的文件
func (p *codeMaker) renderTemplate(makeType string, fileDest string, data any) error {
	// 读取模板文件stub
	modelData, err := stubsFS.ReadFile(path.Join("stubs", makeType+".gotpl"))
	if err != nil {
		return errors.Wrap(err, colorError("读取模板文件出错: ")+fileDest)
	}
	modelStub := string(modelData)

	// 模板解析
	tmpl, err := template.New("myTpl").Parse(modelStub)
	if err != nil {
		return errors.Wrap(err, colorError("解析模板出错"))
	}
	var buf bytes.Buffer
	if err := tmpl.Execute(&buf, data); err != nil {
		return errors.Wrap(err, colorError("执行模板出错"))
	}

	// 存储到目标文件中
	dir := filepath.Dir(fileDest)
	if !file.Exists(dir) {
		if err := os.MkdirAll(dir, os.ModePerm); err != nil {
			return errors.Wrap(err, colorError("创建目录失败: ")+dir)
		}

	}
	if err := file.PutContent(fileDest, buf.Bytes()); err != nil {
		return errors.Wrap(err, colorError("写入文件失败: ")+fileDest)
	}

	log.Println(colorSuccess("生成成功: ") + fileDest)
	return nil
}

func (p *codeMaker) calcFileDest() error {
	p.makeName = strings.Trim(p.makeName, "/")
	if p.makeName == "" {
		return errors.New("待生成实体名称无效(不能为空)")
	}
	names := strings.Split(p.makeName, "/")
	lenNames := len(names)
	if lenNames == 1 && names[0] == "" {
		return errors.New("待生成实体名称不能为空")
	}
	entityName := names[lenNames-1]

	// 模板变量赋值
	model := p.makeModelFromString(entityName)
	p.realName = model.SnakeName
	if p.vars == nil {
		p.vars = make(map[string]any)
	}
	p.vars["SnakeName"] = model.SnakeName
	p.vars["SnakeNamePlural"] = model.SnakeNamePlural
	p.vars["CamelName"] = model.CamelName
	p.vars["CamelNamePlural"] = model.CamelNamePlural
	p.vars["LowerCamelName"] = model.LowerCamelName
	p.vars["LowerCamelNamePlural"] = model.LowerCamelNamePlural

	// 根据类型区分
	switch p.makeType {
	case TypeRepository:
		p.fileDest = path.Join(p.dirBase, "repository", "repository_"+p.realName+".go")
	case TypeService:
		p.fileDest = path.Join(p.dirBase, "service", "service_"+p.realName+".go")
	case TypeHandler: // handler支持子目录
		var subDir string
		p.vars["PackageName"] = "handler"
		p.vars["HandlerPkg"] = ""
		if lenNames > 1 {
			subPkgName := names[lenNames-2]
			if subPkgName != "" {
				p.vars["PackageName"] = subPkgName
			}
			p.vars["HandlerPkg"] = "handler."
			// 子目录
			subDir = strings.Join(names[0:lenNames-1], "/")
			if subDir != "" {
				subDir += "/"
			}
		}
		p.fileDest = path.Join(p.dirBase, "handler", subDir, "handler_"+p.realName+".go")
	default:
		return errors.New(colorWarn("不支持的生成类型: ") + p.makeType)
	}

	p.fileDest, _ = filepath.Abs(p.fileDest)

	return nil
}

func colorSuccess(str string) string {
	return color.Success.Sprint(str)
}
func colorError(str string) string {
	return color.Error.Sprint(str)
}
func colorWarn(str string) string {
	return color.Warn.Sprint(str)
}
