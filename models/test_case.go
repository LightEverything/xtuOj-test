package models

import "gorm.io/gorm"

type TestCase struct {
	gorm.Model
	Identity        string `gorm:"type:varchar(45);column:identity;" json:"identity"`
	ProblemIdentity string `gorm:"type:varchar(45);column:problem_identity;" json:"problem_identity"`
	Output          string `gorm:"type:text;column:output;" json:"output"`
	Input           string `gorm:"type:text;column:input;" json:"input"`
}

func (t *TestCase) TableName() string {
	return "test_case"
}
func GetTestCaseByProblem(problemIdentity string) (data []*TestCase, err error) {
	data = make([]*TestCase, 0)
	err = DB.Table("test_case").Where("problem_identity=?", problemIdentity).Find(&data).Error
	return data, err
}
