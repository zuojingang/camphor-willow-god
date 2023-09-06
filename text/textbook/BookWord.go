package textbook

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// BookWord 字符
type BookWord struct {
	storage.BaseStruct
	BookId        int64
	VolumeIndex   int32
	ChapterIndex  int32
	Index         int32
	WordBookIndex int32
	Word          string
}

// NewBookWord 实例
func NewBookWord() *BookWord {
	return &BookWord{}
}

func (bw *BookWord) TableName() string {
	return "book_word"
}

// BatchCreateBookWords 批量创建字符记录
func BatchCreateBookWords(bws *[]*BookWord) {
	if len(*bws) == 0 {
		return
	}
	db := storage.MysqlDB
	for _, bw := range *bws {
		bw.Id = identified.IdGenerate()
	}
	db.CreateInBatches(bws, 500)
}
