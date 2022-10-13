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
		validation.Numeric(form.Phone, "phone")
		validation.MinSize(form.Phone, 11, "phone")
		validation.MaxSize(form.Phone, 11, "phone")
	} else {
		validation.MaxSize(form.Email, 30, "email")
		validation.Email(form.Email, "email")
	}
}

type RegisterAdmin struct {
	Email      string `json:"email" valid:"Required; Email; MaxSize(30)"`
	Phone      string `json:"phone" valid:"Required; MinSize(11); MaxSize(11); Numeric"`
	Password   string `json:"password" valid:"Required; MinSize(6); MaxSize(40)"`
	Name       string `json:"name" valid:"Required; MinSize(3); MaxSize(50)"`
	Surname    string `json:"surname" valid:"Required; MinSize(3); MaxSize(50)"`
	Patronymic string `json:"patronymic" valid:"Required; MinSize(3); MaxSize(50)"`
}

type RegisterOwner struct {
	Email      string `json:"email" valid:"Required; Email; MaxSize(30)"`
	Phone      string `json:"phone" valid:"Required; MinSize(11); MaxSize(11); Numeric"`
	Password   string `json:"password" valid:"Required; MaxSize(100)"`
	Name       string `json:"name" valid:"Required; MinSize(3); MaxSize(50)"`
	Surname    string `json:"surname" valid:"Required; MinSize(3); MaxSize(50)"`
	Patronymic string `json:"patronymic" valid:"Required; MinSize(3); MaxSize(50)"`
}

type CreateBarInfo struct {
	Email       string      `json:"email" valid:"Required; Email; MaxSize(30)"`
	Phone       string      `json:"phone" valid:"Required; MinSize(11); MaxSize(11); Numeric"`
	Name        string      `json:"name" valid:"Required; MaxSize(50)"`
	Description string      `json:"description" valid:"Required; MaxSize(400)"`
	Address     string      `json:"address" valid:"Required; MaxSize(100)"`
	WorkHours   []WorkHours `json:"work_hours"  valid:"Required; Length(7)"`
}

type WorkHours struct {
	Weekday            uint    `json:"weekday" valid:"Required; Range(1,7)"`
	From               string  `json:"from" valid:"Required; Match(^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$)"`
	To                 string  `json:"to" valid:"Required; Match(^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$)"`
	MaxReservationTime *string `json:"max_reserv_time,omitempty" valid:"Match(^([0-1]?[0-9]|2[0-3]):[0-5][0-9]$)"`
}

func (form *WorkHours) Valid(validation *validation.Validation) {
	if form.MaxReservationTime != nil {
		regex, err := regexp.Compile("^[0-9]*$")
		if err != nil {
			validation.AddError("internal", "Exception while validating")
		}
		validation.Match(form.MaxReservationTime, regex, "work_hours")
	}
}

type UpdateProfile struct {
	Email      *string `json:"email,omitempty"`
	Phone      *string `json:"phone,omitempty"`
	Name       *string `json:"name,omitempty"`
	Surname    *string `json:"surname,omitempty"`
	Patronymic *string `json:"patronymic,omitempty"`
}

func (form *UpdateProfile) Valid(validation *validation.Validation) {
	if form.Phone != nil {
		validation.Numeric(form.Phone, "phone")
		validation.MinSize(form.Phone, 11, "phone")
		validation.MaxSize(form.Phone, 11, "phone")
	}
	if form.Email != nil {
		validation.Email(form.Email, "email")
		validation.MaxSize(form.Email, 30, "email")
	}
	if form.Surname != nil {
		validation.MinSize(form.Surname, 3, "surname")
		validation.MaxSize(form.Surname, 50, "surname")
	}
	if form.Name != nil {
		validation.MinSize(form.Name, 3, "name")
		validation.MaxSize(form.Name, 50, "name")
	}
	if form.Patronymic != nil {
		validation.MinSize(form.Patronymic, 3, "patronymic")
		validation.MaxSize(form.Patronymic, 50, "patronymic")
	}
}
