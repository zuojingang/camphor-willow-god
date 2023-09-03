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

// NewBookWithArgs 实例
func NewBookWithArgs(name string, author string) *Book {
	return &Book{
		Name:       name,
		Author:     author,
		Categories: new([]*Category),
		Tags:       new([]*Tag),
		Volumes:    new([]*Volume)}
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

	var volumes []Volume
	db.Where("book_id=?", b.Id).Find(&volumes)
	if len(volumes) > 0 {
		var volumeIds []int64
		var chapters []Chapter
		for _, volume := range volumes {
			volumeIds = append(volumeIds, volume.Id)
			db.Where("book_id=? and volume_id=?", b.Id, volume.Id).Find(&chapters)
			if len(chapters) > 0 {
				for _, chapter := range chapters {
					db.Where("book_id=? and volume_id=? and chapter_id=?", b.Id, chapter.VolumeId, chapter.Id).Delete(&Paragraph{})
				}
			}
		}
		db.Where("book_id=? and volume_id in ?", b.Id, volumeIds).Delete(&Chapter{})
	}
	db.Where("book_id=?", b.Id).Delete(&Volume{})
	db.Delete(&existBook, existBook.Id)
}
