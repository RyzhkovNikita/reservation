package timing

import (
	"github.com/pkg/errors"
	"time"
)

type Converter interface {
	DateTimeToString(timeObj time.Time) (string, error)
	StringToDateTime(string) (time.Time, error)
	TimeToString(timeObj time.Time) (string, error)
	StringToTime(string) (time.Time, error)
	DateToString(timeObj time.Time) (string, error)
	StringToDate(string) (time.Time, error)
}

func GetConverter() Converter {
	return converterImpl{
		TimeFormat:     "15:04",
		DateFormat:     "02.01.2006",
		DateTimeFormat: "02.01.2006 15:04",
	}
}

type converterImpl struct {
	TimeFormat     string
	DateFormat     string
	DateTimeFormat string
}

var location = time.UTC

func (c converterImpl) DateTimeToString(timeObj time.Time) (string, error) {
	return timeObj.Local().Format(c.DateTimeFormat), nil
}

func (c converterImpl) StringToDateTime(s string) (time.Time, error) {
	parsed, err := time.ParseInLocation(c.DateTimeFormat, s, location)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Error while parsing datetime")
	}
	return parsed, nil
}

func (c converterImpl) TimeToString(timeObj time.Time) (string, error) {
	return timeObj.UTC().Format(c.TimeFormat), nil
}

func (c converterImpl) StringToTime(s string) (time.Time, error) {
	parsed, err := time.ParseInLocation(c.TimeFormat, s, location)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Error while parsing time")
	}
	return parsed, nil
}

func (c converterImpl) DateToString(timeObj time.Time) (string, error) {
	return timeObj.UTC().Format(c.DateFormat), nil
}

func (c converterImpl) StringToDate(s string) (time.Time, error) {
	parsed, err := time.ParseInLocation(c.DateFormat, s, location)
	if err != nil {
		return time.Time{}, errors.Wrap(err, "Error while parsing time")
	}
	return parsed, nil
}
