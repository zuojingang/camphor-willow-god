package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Book 小说/书总览
type Book struct {
	storage.BaseStruct
	OriginId   string
	Name       string
	Author     string
	Category   int64
	Categories *[]*Category   `gorm:"-"`
	Tags       *[]*Tag        `gorm:"-"`
	Volumes    *[]*BookVolume `gorm:"-"`
}

// NewBook 实例
func NewBook() *Book {
	return &Book{
		Categories: new([]*Category),
		Tags:       new([]*Tag),
		Volumes:    new([]*BookVolume)}
}

// NewBookWithArgs 实例
func NewBookWithArgs(name string, author string) *Book {
	return &Book{
		Name:       name,
		Author:     author,
		Categories: new([]*Category),
		Tags:       new([]*Tag),
		Volumes:    new([]*BookVolume)}
}

func (b *Book) TableName() string {
	return "book"
}

func (b *Book) InsertOrUpdate(allowUpdate ...bool) {
	db := storage.MysqlDB
	var existBook Book
	db.Where("name=? and author=?", b.Name, b.Author).Find(&existBook)
	// 新增
	b.Id = existBook.Id
	if b.Id == 0 {
		b.Id = identified.IdGenerate()
		db.Create(b)
		return
	}
	if len(allowUpdate) == 0 || allowUpdate[0] == true {
		db.Updates(b)
	}
}

func (b *Book) Reset() {
	db := storage.MysqlDB
	var existBook Book
	db.Where("name=? and author=?", b.Name, b.Author).Find(&existBook)
	if existBook.Id == 0 {
		return
	}
	b.Id = existBook.Id
	db.Where("book_id=?", b.Id).Delete(&BookTag{})

	var volumes []BookVolume
	db.Where("book_id=?", b.Id).Find(&volumes)
	if len(volumes) > 0 {
		var volumeIndexes []int32
		var chapters []BookChapter
		for _, volume := range volumes {
			volumeIndexes = append(volumeIndexes, volume.Index)
			db.Where("book_id=? and volume_index=?", b.Id, volume.Index).Find(&chapters)
			if len(chapters) > 0 {
				for _, chapter := range chapters {
					for {
						tx := db.Where("book_id=? and volume_index=? and chapter_index=? limit ?", b.Id, chapter.VolumeIndex, chapter.Index, 500).Delete(&BookWord{})
						if tx.RowsAffected == 0 {
							break
						}
					}
				}
			}
		}
		db.Where("book_id=? and volume_index in ?", b.Id, volumeIndexes).Delete(&BookChapter{})
	}
	db.Where("book_id=?", b.Id).Delete(&BookVolume{})
	db.Delete(&existBook, existBook.Id)
}
