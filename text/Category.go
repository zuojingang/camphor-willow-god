package text

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

func (c Category) TableName() string {
	return "category"
}

func (c Category) InsertOrUpdate() {

	db := storage.MysqlDB
	// 新增
	if c.Id == 0 {
		c.Id = identified.IdGenerate()
		db.Create(c)
		return
	}
	db.Updates(c)
}

func PersistentCategories(cs *[]*Category) {
	if len(*cs) == 0 {
		return
	}
	db := storage.MysqlDB
	var names []string
	for _, category := range *cs {
		if category.Id != 0 {
			continue
		}
		names = append(names, category.Name)
	}
	var categoriesQueryResult []Category
	if len(names) > 0 {
		db.Where("name in ?", names).Find(&categoriesQueryResult)
	}
	categoryNameMap := make(map[string]Category)
	for _, categoryQueryResult := range categoriesQueryResult {
		categoryNameMap[categoryQueryResult.Name] = categoryQueryResult
	}
	for index, category := range *cs {
		categoryExist := categoryNameMap[category.Name]
		category.Id = categoryExist.Id
		category.Parent = -1
		if index > 0 {
			category.Parent = (*cs)[index-1].Id
			category.Level = int8(index + 1)
		}
		if category.Id == 0 {
			category.Id = identified.IdGenerate()
			db.Create(category)
		}
	}
}
