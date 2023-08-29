package database

import (
	"yafgo/yafgo-layout/pkg/sys/ycfg"

	"github.com/pkg/errors"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

// db 默认配置
/* [yaml 配置项]
data:
  gorm:
    log_mode: 4               #  1-Silent, 2-Error, 3-Warn, 4-Info
    table_prefix: ""          #
    max_idle_conns: 2         #
    max_open_conns: 0         #
  mysql:
    default: "root:123456@tcp(127.0.0.1:3306)/user?charset=utf8mb4&parseTime=True&loc=Local"
*/

var gormConfigKey = "data.gorm"
var gormConfigDefault = map[string]any{
	// 若 NewGormMysql 时传入了 zap 作为 logger 则该 log_mode 设置无效, 以传入的 zap 实例的 logLevel 为准
	"log_mode":       4,  // 1-Silent, 2-Error, 3-Warn, 4-Info
	"table_prefix":   "", //
	"max_idle_conns": 2,  //
	"max_open_conns": 0,  //
}

// NewGormMysql init gorm mysql
func NewGormMysql(conf *ycfg.Config, lg logger.Interface, source ...string) (*gorm.DB, error) {
	// 初始默认配置
	conf.SetDefault(gormConfigKey, gormConfigDefault)

	// 数据源
	var dbSrc = "default"
	if len(source) > 0 && source[0] != "" {
		dbSrc = source[0]
	}
	dsn := conf.GetString("data.mysql." + dbSrc)
	if dsn == "" {
		return nil, errors.New("mysql数据源 " + dbSrc + " 不存在")
	}

	// mysql config
	mysqlCfg := mysql.Config{
		DSN:                       dsn,   // DSN data source name
		DefaultStringSize:         256,   // string default length
		DisableDatetimePrecision:  true,  // 禁用 datetime 精度，MySQL 5.6 之前的数据库不支持
		DontSupportRenameIndex:    true,  // 重命名索引时采用删除并新建的方式，MySQL 5.7 之前的数据库和 MariaDB 不支持重命名索引
		DontSupportRenameColumn:   true,  // 用 `change` 重命名列，MySQL 8 之前的数据库和 MariaDB 不支持重命名列
		SkipInitializeWithVersion: false, // 根据版本自动配置
	}
	// mysql dialector
	dialector := mysql.New(mysqlCfg)

	// gorm config
	var gormCfg = &gorm.Config{
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: conf.GetString(gormConfigKey + ".table_prefix"),
			// SingularTable: true,
		},
		DisableForeignKeyConstraintWhenMigrating: true,
	}
	logMode := logger.LogLevel(conf.GetInt(gormConfigKey + ".log_mode"))
	if lg != nil {
		gormCfg.Logger = lg.LogMode(logMode)
	} else {
		gormCfg.Logger = logger.Default.LogMode(logMode)
	}

	// gormDB
	db, err := gorm.Open(dialector, gormCfg)
	if err != nil {
		return nil, errors.Wrap(err, "MySQL连接失败")
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxIdleConns(conf.GetInt(gormConfigKey + ".max_idle_conns"))
	sqlDB.SetMaxOpenConns(conf.GetInt(gormConfigKey + ".max_open_conns"))
	return db, nil
}
