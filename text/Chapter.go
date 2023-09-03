package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Chapter 章节
type Chapter struct {
	storage.BaseStruct
	BookId    int64
	VolumeId  int64
	OriginId  string `gorm:"-"`
	Index     int32
	Name      string
	Paragraph *[]*Paragraph `gorm:"-"`
}

// NewChapter 实例
func NewChapter() *Chapter {
	return &Chapter{
		Paragraph: new([]*Paragraph)}
}

func (c Chapter) TableName() string {
	return "book_chapter"
}

func (c Chapter) TryInsert() (bool, Chapter) {

	db := storage.MysqlDB
	db.Where("book_id=? and volume_id=? and name=?", c.BookId, c.VolumeId, c.Name).Find(&c)
	// 新增
	if c.Id == 0 {
		c.Id = identified.IdGenerate()
		db.Create(c)
		return true, c
	}
	return false, c
}

func (c Chapter) ClearParagraphs() {
	storage.MysqlDB.Where("book_id=? and volume_id=? and chapter_id=?", c.BookId, c.VolumeId, c.Id).Delete(NewParagraph())
}
