package migration

import (
	"os"
	"path/filepath"
	"strings"
	"time"
	"yafgo/yafgo-layout/pkg/helper/file"
	"yafgo/yafgo-layout/pkg/sys/ycfg"
)

var migrateCfgKey = "data.migrate"
var migrateCfgDefault = map[string]any{
	"driver":  "mysql",   // mysql, postgresql, sqlite, sqlserver
	"charset": "utf8mb4", //
}

type MigrateCreator struct {
	conf *ycfg.Config
}

func NewMigrateCreator(cfg *ycfg.Config) *MigrateCreator {
	return &MigrateCreator{
		conf: cfg,
	}
}

// Create 创建一个 migration
func (m *MigrateCreator) Create(name string, table string, create bool) {
	m.conf.SetDefault(migrateCfgKey, migrateCfgDefault)

	// 获取所需 stub 文件
	upStub, downStub := m.getStub(table, create)

	// 创建 up.sql 文件
	file.Create(m.getPath(name, "up"), []byte(m.populateStub(upStub, table)))

	// 创建 down.sql 文件
	file.Create(m.getPath(name, "down"), []byte(m.populateStub(downStub, table)))
}

// getStub 获取 migration stub 文件
func (m *MigrateCreator) getStub(table string, create bool) (string, string) {
	if table == "" {
		return "", ""
	}

	driver := m.conf.GetString(migrateCfgKey + ".driver")
	switch Driver(driver) {
	case DriverPostgresql:
		if create {
			return PostgresqlStubs{}.CreateUp(), PostgresqlStubs{}.CreateDown()
		}
		return PostgresqlStubs{}.UpdateUp(), PostgresqlStubs{}.UpdateDown()

	case DriverSqlite:
		if create {
			return SqliteStubs{}.CreateUp(), SqliteStubs{}.CreateDown()
		}
		return SqliteStubs{}.UpdateUp(), SqliteStubs{}.UpdateDown()

	case DriverSqlserver:
		if create {
			return SqlserverStubs{}.CreateUp(), SqlserverStubs{}.CreateDown()
		}
		return SqlserverStubs{}.UpdateUp(), SqlserverStubs{}.UpdateDown()

	default:
		if create {
			return MysqlStubs{}.CreateUp(), MysqlStubs{}.CreateDown()
		}
		return MysqlStubs{}.UpdateUp(), MysqlStubs{}.UpdateDown()
	}
}

// populateStub 替换 migration stub 中的占位符
func (m *MigrateCreator) populateStub(stub string, table string) string {
	stub = strings.ReplaceAll(stub, "DummyDatabaseCharset", m.conf.GetString(migrateCfgKey+".charset"))
	stub = strings.ReplaceAll(stub, "DummyDatabaseCollate", m.conf.GetString(migrateCfgKey+".collate"))

	if table != "" {
		stub = strings.ReplaceAll(stub, "DummyTable", table)
	}

	return stub
}

// getPath 获取 migration 的完整路径
func (m *MigrateCreator) getPath(name string, category string) string {
	pwd, _ := os.Getwd()

	return filepath.Join(pwd, "migrations", time.Now().Format("20060102150405")+"_"+name+"."+category+".sql")
}
