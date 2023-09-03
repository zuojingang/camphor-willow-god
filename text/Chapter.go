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
	OriginId  string
	Index     int32
	Name      string
	Paragraph *[]*Paragraph `gorm:"-"`
}

// NewChapter 实例
func NewChapter() *Chapter {
	return &Chapter{
		Paragraph: new([]*Paragraph)}
}

func (c *Chapter) TableName() string {
	return "book_chapter"
}

func (c *Chapter) InsertOrUpdate(allowUpdate ...bool) {
	db := storage.MysqlDB
	var existChapter Chapter
	db.Where("book_id=? and volume_id=? and `index`=?", c.BookId, c.VolumeId, c.Index).Find(&existChapter)
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
