package models

import (
	"gorm.io/gorm"
	"xtuOj/helper"
)

type User struct {
	gorm.Model
	Identity  string `gorm:"type:varchar(45);column:identity;" json:"identity"`
	Name      string `gorm:"type:varchar(255);column:name;" json:"name"`
	Password  string `gorm:"type:varchar(255);column:password;" json:"password"`
	Phone     string `gorm:"type:varchar(20);column:phone;" json:"phone"`
	Mail      string `gorm:"type:varchar(100);column:mail;" json:"mail"`
	PassNum   int64  `gorm:"type:int;column:pass_num;" json:"pass_num"`
	SubmitNum int64  `gorm:"type:int;column:submit_num;" json:"submit_num"`
	IsAdmin   int    `gorm:"type:tinyint;column:is_admin" json:"is_admin"`
}

func (u *User) TableName() string {
	return "user"
}

func GetUserDetail(identity string) (data *User, e error) {
	data = new(User)
	if err := DB.Omit("password").Where("identity=?", identity).First(&data).Error; err != nil {
		return nil, err
	}
	return data, e
}

func Login(username, password string) (data *User, e error) {
	data = new(User)
	DB.Table("user")
	e = DB.Where("name=? and password=?", username, password).First(&data).Error

	return data, e
}

func Register(name, password, phone, mail string) (uuidstr string, err error) {
	uuidstr, err = helper.GetUuid()
	if err != nil {
		return "", err
	}

	data := User{
		Model:    gorm.Model{},
		Identity: uuidstr,
		Name:     name,
		Password: password,
		Phone:    phone,
		Mail:     mail,
	}

	if err = DB.Create(&data).Error; err != nil {
		return "", err
	}

	return uuidstr, nil
}

func IsUserExist(email string) (ok bool, e error) {
	var count int64
	data := new(User)
	e = DB.Table("user").Count(&count).Where("mail=?", email).Find(&data).Error

	if e != nil {
		return false, e
	}

	// 如果次数大于0
	if count > 0 {
		return false, nil
	}

	return true, nil
}

func GetRankList(offset, size int) (data []*User, count int64, e error) {
	data = make([]*User, 0)
	e = DB.Table("user").Offset(offset).Limit(size).Count(&count).
		Order("finish_problem_num DESC,submit_num").Omit("password").Find(&data).Error

	return data, count, e
}
