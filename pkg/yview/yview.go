package yview

import (
	"html/template"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/pkg/errors"
)

type ViewTpl struct {
	Files []string // 模板文件, 如: []string{"public/web/demo/index.tmpl"}
	Name  string   // 模板名, 默认为 Files 中第一个文件的文件名, 或者 Files 所有模板文件中自定义的 {{ define "模板名" }}, 如: {{ define "demo.tmpl" }}
	Data  any      // 模板用到的数据
}

// HandleView 输出页面
func HandleView(ctx *gin.Context, viewTpl ViewTpl) (err error) {
	if viewTpl.Files == nil || len(viewTpl.Files) == 0 {
		err = errors.New("未指定模板文件")
		return
	}
	if viewTpl.Name == "" {
		viewTpl.Name = path.Base(viewTpl.Files[0])
	}
	tpl := template.New(viewTpl.Name)
	tpl, err = tpl.ParseFiles(viewTpl.Files...)
	if err != nil {
		err = errors.Wrap(err, "解析模板出错")
		return
	}
	err = tpl.Execute(ctx.Writer, viewTpl.Data)
	if err != nil {
		err = errors.Wrap(err, "页面渲染出错")
	}
	return
}
