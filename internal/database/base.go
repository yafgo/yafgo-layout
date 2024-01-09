package database

import (
	"path/filepath"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

var applyFuncs []func(g *gen.Generator)

func addApplyFuncs(funcs ...func(g *gen.Generator)) {
	applyFuncs = append(applyFuncs, funcs...)
}

// RunGenerate 生成 gorm 所需的 model 和 query
func RunGenerate(dsn string) {
	g := gen.NewGenerator(gen.Config{
		OutPath:      filepath.Join("internal", "query"),
		ModelPkgPath: filepath.Join("internal", "model"),
		Mode:         gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	gormdb, _ := gorm.Open(mysql.Open(dsn))
	g.UseDB(gormdb)

	// 处理自定义待生成的 model
	for _, fn := range applyFuncs {
		fn(g)
	}

	g.Execute()
}

// fieldsIgnore 通用忽略字段
var fieldsIgnore = []string{"iEdit", "iSourceType", "iSourceID"}

type querierGetByID interface {
	// where(id=@id)
	GetByID(id int64) (gen.T, error)
}
