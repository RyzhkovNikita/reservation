package input

import (
	"github.com/pkg/errors"
	"reflect"
	"strconv"
	"strings"
)
import "github.com/beego/beego/v2/server/web/context"

var tagName = "input"

func ParseInput(input *context.BeegoInput, model interface{}) error {
	v := reflect.ValueOf(model).Elem()
	if v.Kind() != reflect.Struct {
		return errors.New("Function is for structures only")
	}
	for i := 0; i < v.NumField(); i++ {
		structField := v.Type().Field(i)
		tag := structField.Tag.Get(tagName)
		if tag == "" || tag == "-" {
			continue
		}
		params := strings.Split(tag, ";")
		var paramName string
		var findInPath bool
		var findInQuery bool
		for _, str := range params {
			if strings.HasPrefix(str, ":") {
				paramName = str
			} else if str == "in_path" {
				findInPath = true
			} else if str == "in_query" {
				findInQuery = true
			}
		}
		if paramName == "" {
			continue
		}
		var param string
		if findInPath {
			param = input.Param(paramName)
		}
		if param == "" && findInQuery {
			_, after, _ := strings.Cut(paramName, ":")
			param = input.Query(after)
		}
		if param == "" {
			return errors.New("No parameter '" + string(tag) + "' in input")
		}
		field := v.Field(i)
		if !field.IsValid() {
			return errors.New("Invalid field '" + string(structField.Name))
		}
		if !field.CanSet() {
			return errors.New("Can't set field '" + string(structField.Name))
		}

		switch field.Kind() {
		case reflect.Int:
			parsedInt, err := strconv.ParseInt(param, 10, 64)
			if err != nil {
				return errors.Wrap(err, "Error when parsing")
			}
			field.SetInt(parsedInt)
		case reflect.String:
			field.SetString(param)
		default:
			return errors.New("Unsupported field type: " + field.Kind().String())
		}
	}
	return nil
}
