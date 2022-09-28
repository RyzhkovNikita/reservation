package controllers

import (
	"barckend/utils"
	"github.com/pkg/errors"
)

type RegisterBar struct {
	Email       string   `json:"email" valid:"Required; Email; MaxSize(30)"`
	Name        string   `json:"name" valid:"Required; MaxSize(50)"`
	Description string   `json:"description" valid:"Required; MaxSize(400)"`
	Password    string   `json:"password" valid:"Required; MaxSize(100)"`
	Address     string   `json:"address" valid:"Required; MaxSize(100)"`
	WorkHours   []string `json:"work_hours"  valid:"Required; Length(7)"`
}

func (bar *RegisterBar) ParseFromBody(body []byte) error {
	err := utils.ParseAndValid(bar, body)
	if err != nil {
		return errors.Wrap(err, "Input parsing error")
	}
	return nil
}

type BarInfoResponse struct {
	Id          int64   `json:"id"`
	Email       string  `json:"email"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Address     string  `json:"address"`
	LogoUrl     *string `json:"logo_url"`
}

type LoginForm struct {
	Email    string `json:"email" valid:"Required; Email; MaxSize(30)"`
	Password string `json:"password" valid:"Required; MaxSize(100)"`
}

func (form *LoginForm) ParseFromBody(body []byte) error {
	err := utils.ParseAndValid(form, body)
	if err != nil {
		return errors.Wrap(err, "Input parsing error")
	}
	return nil
}

type AuthorizationPayload struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
