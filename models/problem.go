package models

import (
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	Identity         string             `gorm:"type:varchar(45);column:identity;" json:"identity"`
	ProblemCategorys []*ProblemCategory `gorm:"foreignKey:problem_id;reference:id"`
	//Cid        string `gorm:"type:varchar(45);column:cid;" json:"cid"`
	Title      string      `gorm:"type:varchar(255);column:title;" json:"title"`
	Content    string      `gorm:"type:text;column:content;" json:"content"`
	PassNum    uint        `gorm:"type:int unsigned;column:pass_num;" json:"pass_num"`
	MaxRuntime int         `gorm:"type:int;column:max_runtime" json:"max_runtime"`
	MaxMem     int         `gorm:"type:int;column:max_mem" json:"max_mem"`
	TestCases  []*TestCase `gorm:"foreignKey:problem_identity;reference:identity;"`
}

func (p *Problem) TableName() string {
	return "problem"
}

func GetProblemList(offset, siz int, keyword, categoryIdentity string) (data []Problem, count int64, e error) {
	data = make([]Problem, 0)

	//构造查询参数
	argKeyword := "%" + keyword + "%"

	// 错误返回上一层
	tx := DB.Model(&Problem{}).Count(&count).Offset(offset).Limit(siz).Omit("content").
		Where("title like ? or content like ?", argKeyword, argKeyword)

	if categoryIdentity != "" {
		tx.Joins("RIGHT JOIN problem_category pc ON pc.problem_id = problem.id").
			Where("pc.category_id IN (SELECT id FROM category cg WHERE cg.identity = ?)", categoryIdentity)
	}

	e = tx.Preload("ProblemCategorys").Preload("ProblemCategorys.Category").Find(&data).Error

	return data, count, e
}

func GetProblemDetail(identity string) (data []Problem, e error) {
	data = make([]Problem, 0)
	// 预加载
	DB.Preload("ProblemCategorys").Preload("ProblemCategorys.Category")
	if err := DB.Where("identity=?", identity).First(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func CreateProblem(data *Problem) (e error) {
	e = DB.Table("problem").Create(data).Error

	return e
}
