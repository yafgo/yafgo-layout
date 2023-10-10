package migration

import (
	"database/sql"
	"errors"
	"yafgo/yafgo-layout/pkg/sys/ycfg"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mysql"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/database/sqlite3"
	"github.com/golang-migrate/migrate/v4/database/sqlserver"
	"github.com/gookit/color"
	"github.com/spf13/viper"
)

func Setup(conf *viper.Viper) {
}

type migrator struct {
	cfg *ycfg.Config
}

func NewMigrator(cfg *ycfg.Config) *migrator {
	return &migrator{
		cfg: cfg,
	}
}

func (p *migrator) getMigrate() (*migrate.Migrate, error) {
	conf := p.cfg
	tableMigrations := conf.GetString("data.migrate.table")
	dsn := conf.GetString("data.mysql.default")
	driver := conf.GetString(migrateCfgKey + ".driver")

	// 获取 migrate 实例
	return getMigrate(tableMigrations, driver, dsn)
}

// getMigrate 获取 *migrate.Migrate 实例
//
//	tableMigrations: migrations 记录使用的表名
//	driver: migrations 使用的驱动(mysql, postgresql, sqlite, sqlserver)
//	dsn: 数据库连接 dsn
func getMigrate(tableMigrations string, driver string, dsn string) (*migrate.Migrate, error) {
	// migrations 文件所在目录
	dir := "file://./migrations"

	switch Driver(driver) {
	case DriverMysql:
		db, err := sql.Open("mysql", dsn)
		if err != nil {
			return nil, err
		}

		//if err := db.Ping(); err != nil {
		//	return nil, errors.New("Could not ping to database: " + err.Error())
		//}

		instance, err := mysql.WithInstance(db, &mysql.Config{
			MigrationsTable: tableMigrations,
		})
		if err != nil {
			return nil, err
		}

		return migrate.NewWithDatabaseInstance(dir, "mysql", instance)
	case DriverPostgresql:
		db, err := sql.Open("postgres", dsn)
		if err != nil {
			return nil, err
		}

		instance, err := postgres.WithInstance(db, &postgres.Config{
			MigrationsTable: tableMigrations,
		})
		if err != nil {
			return nil, err
		}

		return migrate.NewWithDatabaseInstance(dir, "postgres", instance)
	case DriverSqlite:
		db, err := sql.Open("sqlite3", dsn)
		if err != nil {
			return nil, err
		}

		instance, err := sqlite3.WithInstance(db, &sqlite3.Config{
			MigrationsTable: tableMigrations,
		})
		if err != nil {
			return nil, err
		}

		return migrate.NewWithDatabaseInstance(dir, "sqlite3", instance)
	case DriverSqlserver:
		db, err := sql.Open("sqlserver", dsn)
		if err != nil {
			return nil, err
		}

		instance, err := sqlserver.WithInstance(db, &sqlserver.Config{
			MigrationsTable: tableMigrations,
		})

		if err != nil {
			return nil, err
		}

		return migrate.NewWithDatabaseInstance(dir, "sqlserver", instance)
	default:
		return nil, errors.New("database driver only support mysql, postgresql, sqlite and sqlserver")
	}
}

// MakeMigration 创建迁移文件
func (p *migrator) MakeMigration(name string) error {
	if name == "" {
		color.Redln("Not enough arguments (missing: name)")
		return nil
	}

	color.Green.Printf("make migration: %s...\n", name)

	// 尝试猜测表名, 判断是否为创建操作
	table, create := TableGuesser{}.Guess(name)

	// 写入 migration 文件
	NewMigrateCreator(p.cfg).Create(name, table, create)

	color.Green.Printf("Created Migration: %s\n", name)
	return nil
}

// Version 返回已执行过的迁移版本信息
func (p *migrator) Version() (version uint, dirty bool, err error) {
	m, err := p.getMigrate()
	if err != nil {
		return
	}
	if m == nil {
		color.Yellowln("Please fill database config first")
		return
	}

	return m.Version()
}

// RunMigrate 执行迁移
func (p *migrator) RunMigrate() error {
	color.Greenln("Run migrate...")

	m, err := p.getMigrate()
	if err != nil {
		color.Yellowln(err.Error())
		return err
	}
	if m == nil {
		color.Yellowln("Please fill database config first")
		return nil
	}

	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		color.Redln("Migration failed:", err.Error())
		return nil
	}

	color.Greenln("Migration success")
	return nil
}

// RunRollback 执行迁移回滚
func (p *migrator) RunRollback(step int) error {
	color.Greenln("Run migrate rollback...")

	m, err := p.getMigrate()
	if err != nil {
		color.Yellowln(err.Error())
		return err
	}
	if m == nil {
		color.Yellowln("Please fill database config first")
		return nil
	}

	// 查看当前版本
	_, _, err = m.Version()
	if err == migrate.ErrNilVersion {
		color.Yellowln("No migrations")
		return nil
	}

	// 执行rollback
	_step := step * -1
	if err := m.Steps(_step); err != nil && err != migrate.ErrNoChange && err != migrate.ErrNilVersion {
		switch err.(type) {
		case migrate.ErrShortLimit:
		default:
			color.Redln("Migration rollback failed:", err.Error())
			return nil
		}
	}

	color.Greenln("Migration rollback success")
	return nil
}

// ForceVersion 强制版本
func (p *migrator) ForceVersion(ver int) error {
	color.Greenln("Run force version...")

	m, err := p.getMigrate()
	if err != nil {
		color.Yellowln(err.Error())
		return err
	}
	if m == nil {
		color.Yellowln("Please fill database config first")
		return nil
	}

	if err := m.Force(ver); err != nil && err != migrate.ErrNoChange {
		color.Redln("Migration failed:", err.Error())
		return nil
	}

	color.Greenln("Migration success")
	return nil
}
