package mapping

import (
	"barckend/controllers/requests/bodies"
	"barckend/controllers/responses"
	"barckend/crud"
	"barckend/timing"
	"github.com/beego/beego/v2/client/orm"
	"github.com/pkg/errors"
	"time"
)

var Mapper ModelMapper = modelMapperImpl{}

type ModelMapper interface {
	AdminDbToNet(profile *crud.AdminInfo) *responses.ProfileResponse
	OwnerDbToNet(profile *crud.OwnerInfo) *responses.ProfileResponse
	WorkHoursListInToDb(barId uint64, workHours []bodies.WorkHours) ([]crud.WorkHours, error)
	BarInfoDbToNet(bar *crud.Bar) (*responses.BarInfoResponse, error)
	TableInfoDbToNet(tableInfo *crud.Table) responses.TableInfo
	TableInfoListDbToNet(tableInfoList []*crud.Table) []responses.TableInfo
	ReservationDbToNet(reservation *crud.Reservation) (responses.Reservation, error)
	ReservationListDbToNet(reservations []*crud.Reservation) ([]responses.Reservation, error)
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

func (m modelMapperImpl) WorkHoursListInToDb(barId uint64, workHours []bodies.WorkHours) ([]crud.WorkHours, error) {
	workHoursOut := make([]crud.WorkHours, len(workHours))
	for index, workHourIn := range workHours {
		fromTime, err := timing.GetConverter().StringToTime(workHourIn.From)
		if err != nil {
			return nil, errors.Wrap(err, "Error when tried to parse work hours time")
		}
		toTime, err := timing.GetConverter().StringToTime(workHourIn.To)
		if err != nil {
			return nil, errors.Wrap(err, "Error when tried to parse work hours time")
		}
		workHourOut := crud.WorkHours{
			Weekday: uint(workHourIn.Weekday),
			From:    orm.DateTimeField(fromTime),
			To:      orm.DateTimeField(toTime),
			Bar:     &crud.Bar{Id: barId},
		}
		workHoursOut[index] = workHourOut
	}
	return workHoursOut, nil
}

func (m modelMapperImpl) BarInfoDbToNet(bar *crud.Bar) (*responses.BarInfoResponse, error) {
	var adminIds = make([]uint64, len(bar.Admins))
	for i, admin := range bar.Admins {
		adminIds[i] = admin.Id
	}
	workHours, err := m.WorkHoursListDbToNet(bar.WorkHours)
	if err != nil {
		return nil, errors.Wrap(err, "Error while converting work hours")
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
		WorkHours:   workHours,
	}, nil
}

func (m modelMapperImpl) WorkHoursListDbToNet(workHours []*crud.WorkHours) ([]responses.WorkHours, error) {
	workHoursOut := make([]responses.WorkHours, len(workHours))
	for i, wh := range workHours {
		fromTime, err := timing.GetConverter().TimeToString(wh.From.Value())
		if err != nil {
			return nil, errors.Wrap(err, "Error when tried to format work hours time")
		}
		toTime, err := timing.GetConverter().TimeToString(wh.To.Value())
		if err != nil {
			return nil, errors.Wrap(err, "Error when tried to parse work hours time")
		}
		workHoursOut[i] = responses.WorkHours{
			Weekday: wh.Weekday,
			From:    fromTime,
			To:      toTime,
		}
	}
	return workHoursOut, nil
}

func (m modelMapperImpl) TableInfoDbToNet(tableInfo *crud.Table) responses.TableInfo {
	return responses.TableInfo{
		Id:          tableInfo.Id,
		Name:        tableInfo.Name,
		Description: tableInfo.Description,
		Capacity:    uint8(tableInfo.Capacity),
	}
}

func (m modelMapperImpl) TableInfoListDbToNet(tableInfoList []*crud.Table) []responses.TableInfo {
	result := make([]responses.TableInfo, 0, len(tableInfoList))
	for _, tableinfo := range tableInfoList {
		result = append(result, m.TableInfoDbToNet(tableinfo))
	}
	return result
}

func (m modelMapperImpl) ReservationDbToNet(reservation *crud.Reservation) (responses.Reservation, error) {
	ToString, err := timing.GetConverter().TimeToString(time.Time(reservation.To))
	if err != nil {
		return responses.Reservation{}, err
	}
	fromString, err := timing.GetConverter().TimeToString(time.Time(reservation.From))
	if err != nil {
		return responses.Reservation{}, err
	}
	dateString, err := timing.GetConverter().DateToString(time.Time(reservation.Date))
	if err != nil {
		return responses.Reservation{}, err
	}
	return responses.Reservation{
		Id:          reservation.Id,
		BarId:       reservation.Table.BarInfo.Id,
		TableId:     reservation.Table.Id,
		From:        fromString,
		To:          ToString,
		PersonCount: reservation.PersonCount,
		Date:        dateString,
		Guest:       nil, //TODO
		IsTech:      reservation.Guest == nil,
		Comment:     reservation.Comment,
	}, nil
}

func (m modelMapperImpl) ReservationListDbToNet(reservations []*crud.Reservation) ([]responses.Reservation, error) {
	outReservations := make([]responses.Reservation, 0, len(reservations))
	for _, dbReservation := range reservations {
		reservationOut, err := m.ReservationDbToNet(dbReservation)
		if err != nil {
			return nil, err
		}
		outReservations = append(outReservations, reservationOut)
	}
	return outReservations, nil
}
