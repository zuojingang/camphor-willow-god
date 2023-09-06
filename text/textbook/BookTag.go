package textbook

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// BookTag 小说/书 标签
type BookTag struct {
	storage.BaseStruct
	BookId int64
	TagId  int64
}

// NewBookTag 实例
func NewBookTag() *BookTag {
	return &BookTag{}
}

// NewBookTagWithArgs 实例
func NewBookTagWithArgs(bookId int64, tagId int64) *BookTag {
	return &BookTag{BookId: bookId, TagId: tagId}
}

func (bt *BookTag) TableName() string {
	return "book_tag"
}

func (bt *BookTag) InsertOrUpdate() {
	db := storage.MysqlDB
	var existBookTag BookTag
	db.Where("book_id=? and tag_id=?", bt.BookId, bt.TagId).Find(&existBookTag)
	// 新增
	bt.Id = existBookTag.Id
	if bt.Id == 0 {
		bt.Id = identified.IdGenerate()
	}
	db.Create(bt)
}
