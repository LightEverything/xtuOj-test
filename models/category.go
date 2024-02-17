package models

import (
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Identity string `gorm:"type:varchar(45);column:identity;" json:"identity"`
	Name     string `gorm:"type:varchar(255);column:name;" json:"name"`
	Parentid int    `gorm:"type:int;column:parentid;" json:"parentid"`
}

func (c *Category) TableName() string {
	return "category"
}

func GetCategoryList(offset, size int, keyword string) (data []*Category, count int64, err error) {
	// 数据库链接
	data = make([]*Category, 0)
	tx := DB.Debug().Table("category").Limit(size).Offset(offset)
	if keyword != "" {
		tx.Where("name=?", keyword)
	}

	err = tx.Find(&data).Count(&count).Error
	return data, count, err
}

func CreateCategory(category Category) error {
	return DB.Table("category").Create(&category).Error
}

func UpdateCategory(data Category) error {
	return DB.Table("category").Where("identity=?", data.Identity).Updates(&data).Error
}

// 获取该分类下问题的数量
func GetProblemsOfCategoryCount(identity string) (count int64, e error) {
	datas := make([]ProblemCategory, 0)
	e = DB.Table("problem_category").Where("category_id IN (SELECT `id` FROM category WHERE identity=?)", identity).
		Count(&count).Find(&datas).Error
	return count, e
}

func DeleteCategory(identity string) error {
	data := &Category{
		Identity: identity,
	}
	return DB.Model(data).Where("identity=?", identity).Delete(data).Error
}
