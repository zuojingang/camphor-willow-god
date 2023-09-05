package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
	"gorm.io/gorm"
)

// BookChapter 章节
type BookChapter struct {
	storage.BaseStruct
	BookId      int64
	VolumeIndex int32
	OriginId    string
	Index       int32
	Name        string
	BookWord    *[]*BookWord `gorm:"-"`
}

// NewBookChapter 实例
func NewBookChapter() *BookChapter {
	return &BookChapter{
		BookWord: new([]*BookWord)}
}

func (c *BookChapter) TableName() string {
	return "book_chapter"
}

func (c *BookChapter) InsertOrUpdate(allowUpdate ...bool) {
	db := storage.MysqlDB
	var existChapter BookChapter
	db.Where("book_id=? and volume_index=? and `index`=?", c.BookId, c.VolumeIndex, c.Index).Find(&existChapter)
	// 新增
	c.Id = existChapter.Id
	if c.Id == 0 {
		c.Id = identified.IdGenerate()
		db.Create(c)
		return
	}
	if len(allowUpdate) == 0 || allowUpdate[0] == false {
		db.Updates(c)
	}
}

func (c *BookChapter) UpdateWordsBookIndex() {
	db := storage.MysqlDB
	var baseWordBookIndex int
	err := db.Model(&BookWord{}).Where("book_id=?", c.BookId).Select("MAX(word_book_index)").Row().Scan(&baseWordBookIndex)
	if err != nil {
		panic(err)
	}
	// 第一章+0索引没问题，后面的章节+0索引会出现索引位重复的问题
	if baseWordBookIndex > 0 {
		baseWordBookIndex++
	}
	db.Model(&BookWord{}).Where("book_id=? and volume_index=? and chapter_index=?", c.BookId, c.VolumeIndex, c.Index).UpdateColumn("word_book_index", gorm.Expr("`index`+?", baseWordBookIndex))
}
