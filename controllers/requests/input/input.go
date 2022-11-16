package input

import (
	"barckend/timing"
	"strings"
	"time"
)

type BarIdInPathInput struct {
	BarId int `input:":bar_id;in_path"`
}

type BarIdInQueryInput struct {
	BarId int `input:":bar_id;in_query"`
}

type TableIdInPathInput struct {
	TableId int `input:":table_id;in_path"`
}

type TableIdInQueryInput struct {
	TableId int `input:":table_id;in_query"`
}

type DatesInQueryInput struct {
	Dates string `input:":dates;in_query"`
}

func (datesInput DatesInQueryInput) GetDates() ([]time.Time, error) {
	datesStr := strings.ReplaceAll(strings.ReplaceAll(datesInput.Dates, "[", ""), "]", "")
	splitted := strings.Split(datesStr, ",")
	dates := make([]time.Time, 0, len(splitted))
	for _, dateStr := range splitted {
		date, err := timing.GetConverter().StringToDate(dateStr)
		if err != nil {
			return nil, err
		}
		dates = append(dates, date)
	}
	return dates, nil
}

type ReservationIdInPathInput struct {
	ReservationId int `input:":reserv_id;in_path"`
}
