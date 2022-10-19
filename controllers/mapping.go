package controllers

import (
	"barckend/crud"
)

var Mapper ModelMapper = modelMapperImpl{}

type ModelMapper interface {
	AdminDbToNet(profile *crud.AdminInfo) *ProfileResponse
	OwnerDbToNet(profile *crud.OwnerInfo) *ProfileResponse
	WorkHoursListInToDb(barId uint64, workHours []WorkHours) []crud.WorkHours
	BarInfoDbToNet(bar *crud.Bar) *BarInfoResponse
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
		Role:       uint(crud.Admin),
	}
}

func (m modelMapperImpl) OwnerDbToNet(profile *crud.OwnerInfo) *ProfileResponse {
	return &ProfileResponse{
		Id:         profile.Id,
		Email:      profile.Email,
		Name:       profile.Name,
		Surname:    profile.Surname,
		Patronymic: profile.Patronymic,
		Phone:      profile.Phone,
		Role:       uint(crud.Owner),
	}
}

func (m modelMapperImpl) WorkHoursListInToDb(barId uint64, workHours []WorkHours) []crud.WorkHours {
	workHoursOut := make([]crud.WorkHours, len(workHours))
	for index, workHourIn := range workHours {
		workHourOut := crud.WorkHours{
			Weekday: workHourIn.Weekday,
			From:    workHourIn.From,
			To:      workHourIn.To,
			Bar:     &crud.Bar{Id: barId},
		}
		workHoursOut[index] = workHourOut
	}
	return workHoursOut
}

func (m modelMapperImpl) BarInfoDbToNet(bar *crud.Bar) *BarInfoResponse {
	return &BarInfoResponse{
		Id:          bar.Id,
		Email:       bar.Email,
		Name:        bar.Name,
		Description: bar.Description,
		Address:     bar.Address,
		LogoUrl:     bar.LogoUrl,
	}
}
