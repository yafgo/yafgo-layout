package database

import "gorm.io/gen"

// 示例: 从已有的数据库表生成
//
//	示例表 `t_user` 结构如下, 执行生成命令前请保证该表已存在
//	其中 `iEdit`,`iSourceType`,`iSourceID` 为实际业务不需要的字段,
//	通过下面代码配置, 生成的model中会将其忽略
/*
CREATE TABLE `t_user` (
	`id` bigint NOT NULL AUTO_INCREMENT,
	`phone` varchar(20) NOT NULL DEFAULT '' COMMENT '手机号',
	`nickname` varchar(50) NOT NULL DEFAULT '' COMMENT '昵称',
	`avatar` varchar(255) NOT NULL DEFAULT '' COMMENT '头像',
	`name` varchar(30) NOT NULL DEFAULT '' COMMENT '用户姓名',
	`gender` char(5) NOT NULL DEFAULT '' COMMENT '男｜女',
	`created_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP,
	`updated_at` timestamp NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
	`deleted_at` timestamp NULL,
	`iEdit` int(11) DEFAULT NULL,
	`iSourceType` int(11) DEFAULT NULL,
	`iSourceID` varchar(50) DEFAULT NULL,
	PRIMARY KEY (`id`),
	KEY `idx_phone` (`phone`)
) ENGINE = InnoDB DEFAULT CHARSET = utf8mb4 COLLATE = utf8mb4_unicode_ci COMMENT = '用户表';
*/
func init() {
	addApplyFuncs(func(g *gen.Generator) {
		// 用户表
		g.ApplyInterface(
			func(userQuerier) {}, // 附加生成该接口定义的 Query 方法
			g.GenerateModelAs(
				"t_user",                         // 数据库表名
				"User",                           // 生成的 model 名
				gen.FieldIgnore(fieldsIgnore...), // 生成的 model 忽略相应字段
			),
		)
	})
}

type userQuerier interface {
	// where(id=@id)
	GetByID(id int64) (*gen.T, error)

	// where(phone=@phone)
	GetByPhone(phone string) (*gen.T, error)

	// where(username=@username)
	GetByUsername(username string) (*gen.T, error)
}
