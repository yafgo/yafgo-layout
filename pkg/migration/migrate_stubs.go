package migration

import _ "embed"

type MysqlStubs struct {
}

//go:embed stubs/mysql_create_up.sql
var stubCreateUp string

// CreateUp Create up migration content.
func (receiver MysqlStubs) CreateUp() string {
	return stubCreateUp
}

// CreateDown Create down migration content.
func (receiver MysqlStubs) CreateDown() string {
	return "DROP TABLE IF EXISTS `t_DummyTable`;"
}

// UpdateUp Update up migration content.
func (receiver MysqlStubs) UpdateUp() string {
	return "ALTER TABLE `t_DummyTable` ADD column varchar(255) COMMENT '';"
}

// UpdateDown Update down migration content.
func (receiver MysqlStubs) UpdateDown() string {
	return "ALTER TABLE `t_DummyTable` DROP COLUMN column;"
}

type PostgresqlStubs struct {
}

// CreateUp Create up migration content.
func (receiver PostgresqlStubs) CreateUp() string {
	return `CREATE TABLE DummyTable (
  id SERIAL PRIMARY KEY NOT NULL,
  created_at timestamp NOT NULL,
  updated_at timestamp NOT NULL
);
`
}

// CreateDown Create down migration content.
func (receiver PostgresqlStubs) CreateDown() string {
	return `DROP TABLE IF EXISTS DummyTable;
`
}

// UpdateUp Update up migration content.
func (receiver PostgresqlStubs) UpdateUp() string {
	return `ALTER TABLE DummyTable ADD column varchar(255) NOT NULL;
`
}

// UpdateDown Update down migration content.
func (receiver PostgresqlStubs) UpdateDown() string {
	return `ALTER TABLE DummyTable DROP COLUMN column;
`
}

type SqliteStubs struct {
}

// CreateUp Create up migration content.
func (receiver SqliteStubs) CreateUp() string {
	return `CREATE TABLE DummyTable (
  id integer PRIMARY KEY AUTOINCREMENT NOT NULL,
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL
);
`
}

// CreateDown Create down migration content.
func (receiver SqliteStubs) CreateDown() string {
	return `DROP TABLE IF EXISTS DummyTable;
`
}

// UpdateUp Update up migration content.
func (receiver SqliteStubs) UpdateUp() string {
	return `ALTER TABLE DummyTable ADD column text;
`
}

// UpdateDown Update down migration content.
func (receiver SqliteStubs) UpdateDown() string {
	return `ALTER TABLE DummyTable DROP COLUMN column;
`
}

type SqlserverStubs struct {
}

// CreateUp Create up migration content.
func (receiver SqlserverStubs) CreateUp() string {
	return `CREATE TABLE DummyTable (
  id bigint NOT NULL IDENTITY(1,1),
  created_at datetime NOT NULL,
  updated_at datetime NOT NULL,
  PRIMARY KEY (id)
);
`
}

// CreateDown Create down migration content.
func (receiver SqlserverStubs) CreateDown() string {
	return `DROP TABLE IF EXISTS DummyTable;
`
}

// UpdateUp Update up migration content.
func (receiver SqlserverStubs) UpdateUp() string {
	return `ALTER TABLE DummyTable ADD column varchar(255);
`
}

// UpdateDown Update down migration content.
func (receiver SqlserverStubs) UpdateDown() string {
	return `ALTER TABLE DummyTable DROP COLUMN column;
`
}
