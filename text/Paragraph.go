package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Paragraph 段落
type Paragraph struct {
	storage.BaseStruct
	BookId    int64
	VolumeId  int64
	ChapterId int64
	Index     int32
	Content   string
}

// NewParagraph 实例
func NewParagraph() *Paragraph {
	return &Paragraph{}
}

func (p *Paragraph) TableName() string {
	return "book_paragraph"
}

func PersistentParagraphs(ps *[]*Paragraph) {
	if len(*ps) == 0 {
		return
	}
	db := storage.MysqlDB
	for _, paragraph := range *ps {
		paragraph.Id = identified.IdGenerate()
	}
	db.CreateInBatches(ps, len(*ps))
}
