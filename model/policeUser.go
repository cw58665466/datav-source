package model

import (
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User 用户模型
type XPoliceUser struct {
	gorm.Model
	Name     string
	Password string
	Email    string
	Status   string
	Mobile   string
}

func (User) TableName() string {
	return "x_police_users"
}

// CheckPassword 校验密码
func (user *XPoliceUser) CheckNewPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	return err == nil
}

// GetUser 用ID获取用户
func PoliceUser(ID interface{}) (XPoliceUser, error) {
	var user XPoliceUser
	result := DB.First(&user, ID)
	return user, result.Error
}
