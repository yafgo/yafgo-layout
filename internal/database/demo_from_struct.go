package database

import "gorm.io/gen"

// 示例: 从已有的结构体生成
func init() {
	addApplyFuncs(func(g *gen.Generator) {
		// 生成基础功能 model
		// g.ApplyBasic(MyStruct{})

		// 生成带自定义方法的 model
		g.ApplyInterface(func(querierGetByID) {}, MyStruct{})
	})
}

type MyStruct struct {
	ID   int64  `gorm:"column:id;primaryKey;autoIncrement:true" json:"id"`
	Name string `gorm:"column:name;not null;comment:name" json:"name"`
	Age  int64  `gorm:"column:age;not null;comment:age" json:"age"`
}
