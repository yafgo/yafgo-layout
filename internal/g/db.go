package g

import (
	"gorm.io/gorm"
)

var Mysql *gorm.DB // 全局默认 mysql 对象

// var MysqlOther *gorm.DB // 另一个数据源的 mysql 对象
