package models

import "gorm.io/gorm"

type Submit struct {
	gorm.Model
	Identity        string  `gorm:"type:varchar(45);column:identity;" json:"identity"`
	ProblemIdentity string  `gorm:"type:varchar(45);column:problem_identity;" json:"problem_identity"`
	Problem         Problem `gorm:"foreignKey:problem_identity;reference:identity"`
	UserIdentity    string  `gorm:"type:varchar(45);column:user_identity;" json:"user_identity"`
	User            User    `gorm:"foreignKey:user_identity;reference:identity"`
	Path            string  `gorm:"type:varchar(255);column:path" json:"path"`
	Status          int     `gorm:"type:tinyint;column:status" json:"status"`
}

func (s *Submit) TableName() string {
	return "submit"
}

func GetSubmitList(offset, size int, problemIdentity, userIdentity string, status int) (st []*Submit, count int64, e error) {
	tx := DB.Model(&Submit{}).Preload("User").Preload("Problem", func(tx *gorm.DB) *gorm.DB {
		return tx.Omit("content")
	}).Offset(offset).Limit(size)

	if problemIdentity != "" {
		tx.Where("problem_identity=?", problemIdentity)
	}
	if userIdentity != "" {
		tx.Where("user_identity=?", userIdentity)
	}
	if status != 0 {
		tx.Where("status=?", status)
	}

	st = make([]*Submit, 0)

	e = tx.Count(&count).Find(&st).Error
	return st, count, e
}
