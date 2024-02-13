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
	tx := DB.Table("category").Limit(size).Offset(offset).Count(&count)

	if keyword != "" {
		tx.Where("name=?", keyword)
	}

	err = tx.Find(&data).Error
	return data, count, err
}
