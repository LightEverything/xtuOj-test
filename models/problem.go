package models

import (
	"gorm.io/gorm"
)

type Problem struct {
	gorm.Model
	Identity         string             `gorm:"type:varchar(45);column:identity;" json:"identity"`
	ProblemCategorys []*ProblemCategory `gorm:"foreignKey:problem_id;references:id" json:"problem_categories"`
	Title            string             `gorm:"type:varchar(255);column:title;" json:"title"`
	Content          string             `gorm:"type:text;column:content;" json:"content"`
	PassNum          uint               `gorm:"type:int unsigned;column:pass_num;" json:"pass_num"`
	SubmitNum        int64              `gorm:"type:int;column:submit_num;" json:"submit_num"`
	MaxRuntime       int                `gorm:"type:int;column:max_runtime" json:"max_runtime"`
	MaxMem           int                `gorm:"type:int;column:max_mem" json:"max_mem"`
	TestCases        []*TestCase        `gorm:"foreignKey:problem_identity;references:identity;" json:"test_cases"`
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

	e = tx.Preload("ProblemCategorys").Preload("ProblemCategorys.CategoryDetail").Find(&data).Error

	return data, count, e
}

// 获取问题详细信息
func GetProblemDetail(identity string) (data []Problem, e error) {
	data = make([]Problem, 0)
	// 预加载
	tx := DB.Preload("ProblemCategorys").Preload("ProblemCategorys.CategoryDetail").Preload("TestCases")
	if err := tx.Where("identity=?", identity).First(&data).Error; err != nil {
		return nil, err
	}

	return data, nil
}

func CreateProblem(data *Problem) (e error) {
	e = DB.Table("problem").Create(data).Error

	return e
}

// 修改问题的事务函数
func UpdateProblemTransaction(data *Problem) func(*gorm.DB) error {
	return func(tx *gorm.DB) error {
		// 更新分类表
		// 1. 删除原来的关联表
		if err := tx.Table("problem_category").Where("problem_id IN (SELECT id FROM problem "+
			"WHERE identity=?)", data.Identity).Delete(&ProblemCategory{}).Error; err != nil {
			return err
		}

		// 更新测试集表
		// 1. 删除原有的测试集
		if err := tx.Table("test_case").Where("problem_identity=?", data.Identity).Delete(&TestCase{}).Error; err != nil {
			return err
		}

		// 更新问题表
		if err := tx.Table("problem").Updates(data).Error; err != nil {
			return err
		}
		return nil
	}
}

func UpdateProblem(data *Problem) error {
	return DB.Transaction(UpdateProblemTransaction(data))
}

func GetProblemId(identity string) (id uint, e error) {
	data := &Problem{}
	e = DB.Table("problem").Where("identity=?", identity).First(data).Error

	return data.ID, e
}
