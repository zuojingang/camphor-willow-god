package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Volume 分卷
type Volume struct {
	storage.BaseStruct
	BookId   int64
	Index    int32
	Name     string
	Chapters *[]*Chapter `gorm:"-"`
}

// NewVolume 实例
func NewVolume() *Volume {
	return &Volume{
		Chapters: new([]*Chapter)}
}

func (v Volume) TableName() string {
	return "book_volume"
}

func (v Volume) TryInsert() (bool, Volume) {

	db := storage.MysqlDB
	db.Where("book_id=? and name=?", v.BookId, v.Name).Find(&v)
	// 新增
	if v.Id == 0 {
		v.Id = identified.IdGenerate()
		db.Create(v)
		return true, v
	}
	return false, v
}
