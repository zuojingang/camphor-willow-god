package textbook

import (
	"camphor-willow-god/identified"
	"camphor-willow-god/storage"
)

// Category 小说/书 分类
type Category struct {
	storage.BaseStruct
	Parent      int64
	Name        string
	Level       int8
	Description string
}

// NewCategory 实例
func NewCategory() *Category {
	return &Category{}
}

// NewCategoryByName NewCategory 实例
func NewCategoryByName(name string) *Category {
	return &Category{Name: name}
}

func (c *Category) TableName() string {
	return "category"
}

func (c *Category) InsertOrUpdate() {
	db := storage.MysqlDB
	var existCategory Category
	db.Where("name=?", c.Name).Find(&existCategory)
	c.Id = existCategory.Id
	// 新增
	if c.Id == 0 {
		c.Id = identified.IdGenerate()
	}
	db.Create(c)
}
