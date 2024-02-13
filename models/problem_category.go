package models

import "gorm.io/gorm"

type ProblemCategory struct {
	gorm.Model
	ProblemId  uint `gorm:"type:int;column:problem_id" json:"problem_id"`
	CategoryId uint `gorm:"type:int;column:category_id" json:"category_id"`

	Category *Category `gorm:"foreignKey:category_id;reference:id"`
}

func (pc *ProblemCategory) TableName() string {
	return "problem_category"
}
