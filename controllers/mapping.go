package controllers

import (
	"barckend/crud"
)

var Mapper ModelMapper = modelMapperImpl{}

type ModelMapper interface {
	AdminDbToNet(profile *crud.AdminInfo) *ProfileResponse
}

type modelMapperImpl struct{}

func (m modelMapperImpl) AdminDbToNet(profile *crud.AdminInfo) *ProfileResponse {
	return &ProfileResponse{
		Id:         profile.Id,
		Email:      profile.Email,
		Name:       profile.Name,
		Surname:    profile.Surname,
		Patronymic: profile.Patronymic,
		Phone:      profile.Phone,
		Role:       uint(profile.User.Role),
	}
}
