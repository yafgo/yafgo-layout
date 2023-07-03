package migration

import (
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

type Driver string

func (d Driver) String() string {
	return string(d)
}

const (
	DriverMysql      Driver = "mysql"
	DriverPostgresql Driver = "postgresql"
	DriverSqlite     Driver = "sqlite"
	DriverSqlserver  Driver = "sqlserver"
)
