package storage

import (
	"time"
)

// BaseStruct 基础结构 https://gorm.io/zh_CN/docs/models.html
type BaseStruct struct {
	Id         int64     `gorm:"primaryKey;<-:create"`
	CreateTime time.Time `gorm:"<-:false"`
	UpdateTime time.Time `gorm:"<-:false"`
}

// InsertOrUpdate 用来看哪些类实现了这个方法，暂时没别的用途
type InsertOrUpdate interface {

	// InsertOrUpdate allowUpdate default is {true}
	InsertOrUpdate(allowUpdate ...bool)
}
