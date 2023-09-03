package text

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

// NewTag 实例
func NewBookTag() *BookTag {
	return &BookTag{}
}

// NewBookTagWithArgs 实例
func NewBookTagWithArgs(bookId int64, tagId int64) *BookTag {
	return &BookTag{BookId: bookId, TagId: tagId}
}

func (bt BookTag) TableName() string {
	return "book_tag"
}

func (bt BookTag) InsertOrUpdate() {

	db := storage.MysqlDB
	// 新增
	if bt.Id == 0 {
		bt.Id = identified.IdGenerate()
		db.Create(bt)
		return
	}
	db.Updates(bt)
}

func PersistentBookTags(bts *[]*BookTag) {
	if len(*bts) == 0 {
		return
	}
	db := storage.MysqlDB
	var bookId int64
	var tagIds []int64
	for _, bookTag := range *bts {
		if bookTag.Id != 0 {
			continue
		}
		bookId = bookTag.BookId
		tagIds = append(tagIds, bookTag.TagId)
	}
	var bookTagsQueryResult []BookTag
	if len(tagIds) > 0 {
		db.Where("book_id=? and tag_id in ?", bookId, tagIds).Find(&bookTagsQueryResult)
	}
	bookTagTagIdMap := make(map[int64]BookTag)
	for _, bookTagQueryResult := range bookTagsQueryResult {
		bookTagTagIdMap[bookTagQueryResult.TagId] = bookTagQueryResult
	}
	var bookTags2Create []BookTag
	for _, bookTag := range *bts {
		bookTagExist := bookTagTagIdMap[bookTag.TagId]
		bookTag.Id = bookTagExist.Id
		if bookTag.Id == 0 {
			bookTag.Id = identified.IdGenerate()
			bookTags2Create = append(bookTags2Create, *bookTag)
		}
	}
	if len(bookTags2Create) == 0 {
		return
	}
	db.CreateInBatches(bookTags2Create, len(bookTags2Create))
}
