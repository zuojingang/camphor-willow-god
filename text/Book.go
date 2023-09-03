package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Book 小说/书总览
type Book struct {
	storage.BaseStruct
	Name       string
	Author     string
	Category   int64
	Categories *[]*Category `gorm:"-"`
	Tags       *[]*Tag      `gorm:"-"`
	Volumes    *[]*Volume   `gorm:"-"`
}

// NewBook 实例
func NewBook() *Book {
	return &Book{
		Categories: new([]*Category),
		Tags:       new([]*Tag),
		Volumes:    new([]*Volume)}
}

func (b Book) TableName() string {
	return "book"
}

func (b Book) TryInsert() (bool, Book) {

	db := storage.MysqlDB
	db.Where("name=? and author=?", b.Name, b.Author).Find(&b)
	// 新增
	if b.Id == 0 {
		b.Id = identified.IdGenerate()
		db.Create(b)
		return true, b
	}
	return false, b
}
