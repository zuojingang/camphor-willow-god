package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// BookVolume 分卷
type BookVolume struct {
	storage.BaseStruct
	BookId   int64
	Index    int32
	Name     string
	Chapters *[]*BookChapter `gorm:"-"`
}

// NewBookVolume 实例
func NewBookVolume() *BookVolume {
	return &BookVolume{
		Chapters: new([]*BookChapter)}
}

func (v *BookVolume) TableName() string {
	return "book_volume"
}

func (v *BookVolume) InsertOrUpdate(allowUpdate ...bool) {
	db := storage.MysqlDB
	var existVolume BookVolume
	db.Where("book_id=? and name=?", v.BookId, v.Name).Find(&existVolume)
	// 新增
	v.Id = existVolume.Id
	if v.Id == 0 {
		v.Id = identified.IdGenerate()
		db.Create(v)
		return
	}
	if len(allowUpdate) == 0 || allowUpdate[0] == false {
		db.Updates(v)
	}
}
