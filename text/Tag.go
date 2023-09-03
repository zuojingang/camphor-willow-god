package text

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Tag 小说/书 标签
type Tag struct {
	storage.BaseStruct
	Name string
}

// NewTag 实例
func NewTag() *Tag {
	return &Tag{}
}

// NewTagByName 实例
func NewTagByName(name string) *Tag {
	return &Tag{Name: name}
}

func (t Tag) TableName() string {
	return "tag"
}

func (t Tag) InsertOrUpdate() {

	db := storage.MysqlDB
	// 新增
	if t.Id == 0 {
		t.Id = identified.IdGenerate()
		db.Create(t)
		return
	}
	db.Updates(t)
}

func PersistentTags(ts *[]*Tag) {
	if len(*ts) == 0 {
		return
	}
	db := storage.MysqlDB
	var names []string
	for _, tag := range *ts {
		if tag.Id != 0 {
			continue
		}
		names = append(names, tag.Name)
	}
	var tagsQueryResult []Tag
	if len(names) > 0 {
		db.Where("name in ?", names).Find(&tagsQueryResult)
	}
	tagNameMap := make(map[string]Tag)
	for _, tagQueryResult := range tagsQueryResult {
		tagNameMap[tagQueryResult.Name] = tagQueryResult
	}
	var tags2Create []Tag
	for _, tag := range *ts {
		tagExist := tagNameMap[tag.Name]
		tag.Id = tagExist.Id
		if tag.Id == 0 {
			tag.Id = identified.IdGenerate()
			tags2Create = append(tags2Create, *tag)
		}
	}
	if len(tags2Create) == 0 {
		return
	}
	db.CreateInBatches(tags2Create, len(tags2Create))
}
