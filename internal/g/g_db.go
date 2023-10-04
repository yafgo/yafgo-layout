package g

import (
	"sync"
	"yafgo/yafgo-layout/internal/query"
	"yafgo/yafgo-layout/pkg/database"
	"yafgo/yafgo-layout/pkg/logger"
	"yafgo/yafgo-layout/pkg/sys/ylog"

	"gorm.io/gorm"
)

// 全局 mysql 对象
var instMysql *gorm.DB
var onceMysql sync.Once

// Mysql 全局默认 mysql 实例
func Mysql() *gorm.DB {
	onceMysql.Do(func() {
		gCfg := Cfg()
		gormLogger := logger.NewGormLogger(ylog.DefaultLogger())
		gormDB, err := database.NewGormMysql(gCfg, gormLogger)
		if err != nil {
			panic(err)
		}
		instMysql = gormDB
	})

	return instMysql
}

// Query 全局默认 gorm Query 实例
func Query() *query.Query {
	if !query.Q.Available() {
		db := Mysql()
		// 设置 query 包默认使用的 db 实例
		query.SetDefault(db)
	}
	return query.Q
}

// var MysqlOther *gorm.DB // 另一个数据源的 mysql 对象
