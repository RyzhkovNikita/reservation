package controllers

import (
	"barckend/utils"
	"github.com/beego/beego/v2/core/validation"
	"github.com/pkg/errors"
	"regexp"
)

type BaseBody struct{}

func Require(b any, body []byte) error {
	err := utils.ParseAndValid(b, body)
	if err != nil {
		return errors.Wrap(err, "Input parsing error")
	}
	return nil
}

type LoginForm struct {
	Email    *string `json:"email,omitempty"`
	Phone    *string `json:"phone,omitempty"`
	Password string  `json:"password" valid:"Required; MinSize(6); MaxSize(40)"`
}

func (form *LoginForm) Valid(validation *validation.Validation) {
	if form.Email == nil && form.Phone == nil {
		validation.AddError("login_check", "Phone or email required")
	}
	if form.Phone != nil {
		regex, err := regexp.Compile("^[0-9]*$")
		if err != nil {
			validation.AddError("internal", "Exception while validating")
		}
		validation.Match(form.Phone, regex, "phone")
		validation.MaxSize(form.Phone, 11, "phone")
	} else {
		validation.MaxSize(form.Email, 30, "email")
		validation.Email(form.Email, "email")
	}
}

type RegisterAdmin struct {
	Email      string `json:"email" valid:"Required; Email; MaxSize(30)"`
	Phone      string `json:"phone" valid:"Required; Phone"`
	Password   string `json:"password" valid:"Required; MinSize(6); MaxSize(40)"`
	Name       string `json:"name" valid:"Required; MinSize(3); MaxSize(50)"`
	Surname    string `json:"surname" valid:"Required; MinSize(3); MaxSize(50)"`
	Patronymic string `json:"patronymic" valid:"Required; MinSize(3); MaxSize(50)"`
}

type RegisterOwner struct {
	Email      string `json:"email" valid:"Required; Email; MaxSize(30)"`
	Phone      string `json:"phone" valid:"Required; MinSize(11); MaxSize(11); Phone"`
	Password   string `json:"password" valid:"Required; MaxSize(100)"`
	Name       string `json:"name" valid:"Required; MinSize(3); MaxSize(50)"`
	Surname    string `json:"surname" valid:"Required; MinSize(3); MaxSize(50)"`
	Patronymic string `json:"patronymic" valid:"Required; MinSize(3); MaxSize(50)"`
}

type RegisterBar struct {
	Email       string   `json:"email" valid:"Required; Email; MaxSize(30)"`
	Name        string   `json:"name" valid:"Required; MaxSize(50)"`
	Description string   `json:"description" valid:"Required; MaxSize(400)"`
	Password    string   `json:"password" valid:"Required; MaxSize(100)"`
	Address     string   `json:"address" valid:"Required; MaxSize(100)"`
	WorkHours   []string `json:"work_hours"  valid:"Required; Length(7)"`
}
