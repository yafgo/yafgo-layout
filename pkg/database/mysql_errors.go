package database

import (
	"github.com/go-sql-driver/mysql"
)

// 在代码中怎么判断 Duplicate Entry error:
//   https://github.com/go-gorm/gorm/issues/5144#issuecomment-1066036092

var (
	// ErrDuplicateEntryCode 命中唯一索引
	ErrDuplicateEntryCode = 1062
)

// MysqlErrCode 根据mysql错误信息返回错误代码
func MysqlErrCode(err error) int {
	mysqlErr, ok := err.(*mysql.MySQLError)
	if !ok {
		return 0
	}
	return int(mysqlErr.Number)
}

// IsErrDuplicateEntryCode 是否命中唯一索引错误
func IsErrDuplicateEntryCode(err error) bool {
	return MysqlErrCode(err) == ErrDuplicateEntryCode
}
