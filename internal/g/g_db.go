package g

import (
	"sync"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/database"
	"yafgo/yafgo-layout/pkg/logger"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"gorm.io/gorm"
)

// 全局默认 mysql 对象
func Mysql() *gorm.DB {
	return newMysql()()
}

func newMysql() func() *gorm.DB {
	var once sync.Once
	var db *gorm.DB

	return func() *gorm.DB {
		once.Do(func() {
			gCfg := Cfg()
			gormLogger := logger.NewGormLogger(ylog.DefaultLogger())
			gormDB, err := database.NewGormMysql(gCfg, gormLogger)
			if err != nil {
				panic(err)
			}
			db = gormDB
			// 设置 query 包默认使用的 db 实例
			query.SetDefault(gormDB)
		})

		return db
	}
}

// var MysqlOther *gorm.DB // 另一个数据源的 mysql 对象
