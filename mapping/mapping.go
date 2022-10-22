package mapping

import (
	"barckend/controllers/requests/bodies"
	"barckend/controllers/responses"
	"barckend/crud"
)

var Mapper ModelMapper = modelMapperImpl{}

type ModelMapper interface {
	AdminDbToNet(profile *crud.AdminInfo) *responses.ProfileResponse
	OwnerDbToNet(profile *crud.OwnerInfo) *responses.ProfileResponse
	WorkHoursListInToDb(barId uint64, workHours []bodies.WorkHours) []crud.WorkHours
	BarInfoDbToNet(bar *crud.Bar) *responses.BarInfoResponse
}

type modelMapperImpl struct{}

func (m modelMapperImpl) AdminDbToNet(profile *crud.AdminInfo) *responses.ProfileResponse {
	return &responses.ProfileResponse{
		Id:         profile.Id,
		Email:      profile.Email,
		Name:       profile.Name,
		Surname:    profile.Surname,
		Patronymic: profile.Patronymic,
		Phone:      profile.Phone,
		Role:       uint(crud.Admin),
	}
}

func (m modelMapperImpl) OwnerDbToNet(profile *crud.OwnerInfo) *responses.ProfileResponse {
	return &responses.ProfileResponse{
		Id:         profile.Id,
		Email:      profile.Email,
		Name:       profile.Name,
		Surname:    profile.Surname,
		Patronymic: profile.Patronymic,
		Phone:      profile.Phone,
		Role:       uint(crud.Owner),
	}
}

func (m modelMapperImpl) WorkHoursListInToDb(barId uint64, workHours []bodies.WorkHours) []crud.WorkHours {
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

func (m modelMapperImpl) BarInfoDbToNet(bar *crud.Bar) *responses.BarInfoResponse {
	var adminIds = make([]uint64, len(bar.Admins))
	for _, admin := range bar.Admins {
		adminIds = append(adminIds, admin.Id)
	}
	return &responses.BarInfoResponse{
		Id:          bar.Id,
		OwnerId:     bar.OwnerInfo.Id,
		Email:       bar.Email,
		Phone:       bar.Phone,
		Name:        bar.Name,
		Description: bar.Description,
		Address:     bar.Address,
		Admins:      adminIds,
		LogoUrl:     bar.LogoUrl,
		IsVisible:   bar.IsVisible,
	}
}
