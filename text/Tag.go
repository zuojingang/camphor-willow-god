package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Tag 小说/书 标签
type Tag struct {
	storage.BaseStruct
	Name string
}

// NewTag 实例
func NewTag() *Tag {
	return &Tag{}
}

// NewTagByName 实例
func NewTagByName(name string) *Tag {
	return &Tag{Name: name}
}

func (t *Tag) TableName() string {
	return "tag"
}

func (t *Tag) InsertOrUpdate() {
	db := storage.MysqlDB
	var existTag Tag
	db.Where("name=?", t.Name).Find(&existTag)
	// 新增
	t.Id = existTag.Id
	if t.Id == 0 {
		t.Id = identified.IdGenerate()
		db.Create(t)
		return
	}
	db.Updates(t)
}
