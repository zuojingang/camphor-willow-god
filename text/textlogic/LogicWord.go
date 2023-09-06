package textlogic

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// LogicWord 词组
type LogicWord struct {
	storage.BaseStruct
	Head                 string // 词组起始字符
	LogicAssociationWord        // 起始字符的联想词组
	Level                int32  // 词组层级
}

func (lw *LogicWord) TableName() string {
	return "logic_word"
}

func (lw *LogicWord) InsertOrUpdate() {
	if lw == nil {
		return
	}
	db := storage.MysqlDB
	var existLogicWord LogicWord
	db.Model(&LogicWord{}).Where("head=? and distance=? and end=?", lw.Head, lw.Distance, lw.End).Find(&existLogicWord)
	// 如果当前词汇已经存在更高层级的记录，直接忽略
	if existLogicWord.Level >= lw.Level {
		return
	}
	lw.Id = existLogicWord.Id
	if lw.Id == 0 {
		lw.Id = identified.IdGenerate()
	}
	db.Save(lw)
}
