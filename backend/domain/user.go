package domain

import (
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

const passwordCost = 12

// User .
type User struct {
	Model
	Username       string       `yaml:"username" json:"username"`
	Password       string       `yaml:"-" json:"-"`
	PasswordDigest string       `yaml:"password_digest" json:"password_digest"`
	Role           string       `yaml:"role" json:"role"`
	ConfirmCode    string       `yaml:"confirm_code" json:"confirm_code"`
	IsEmailAccount bool         `yaml:"is_email_account" json:"is_email_account"`
	Email          string       `yaml:"email" json:"email"`
	Groups         []Group      `gorm:"many2many:user_groups;" yaml:"-" json:"-"`
	Permissions    []Permission `yaml:"-" json:"-"`
}

// BeforeSave prepare data before create data
func (u *User) BeforeSave(scope *gorm.Scope) error {
	if u.Password == "" {
		return nil
	}

	passwordDigest, err := bcrypt.GenerateFromPassword([]byte(u.Password), passwordCost)
	if err != nil {
		return err
	}

	if err := scope.SetColumn("PasswordDigest", passwordDigest); err != nil {
		return err
	}

	return nil
}
