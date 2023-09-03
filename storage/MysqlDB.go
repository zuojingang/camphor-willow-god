package storage

import (
	"camphor-willow-god/config"
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

var MysqlDB *gorm.DB

// 初始化数据库连接
func init() {
	// Mysql 配置
	mysqlConfig := config.ApplicationConfig.Mysql
	// 连接数据库
	gormDB, err := gorm.Open(mysql.Open(mysqlConfig.DSN), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db, err := gormDB.DB()
	if err != nil {
		panic(err)
	}
	// 设置数据库最大连接数
	db.SetMaxOpenConns(mysqlConfig.MaxOpen)
	// 设置连接最大存活时间
	db.SetConnMaxLifetime(time.Second * time.Duration(mysqlConfig.MaxLifetime))
	// 设置数据库最大空闲连接数
	idle := mysqlConfig.MaxLifetime
	if mysqlConfig.MaxOpen/3 > 10 {
		// 允许1/3空闲连接
		idle = mysqlConfig.MaxOpen / 3
	}
	db.SetMaxIdleConns(idle)
	// 验证连接
	if err := db.Ping(); err != nil {
		panic(err)
	}
	MysqlDB = gormDB
	fmt.Println("connect database success")
}
