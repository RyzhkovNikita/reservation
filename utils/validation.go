package utils

import (
	"encoding/json"
	"fmt"
	"github.com/beego/beego/v2/core/validation"
	"github.com/pkg/errors"
	"strings"
)

func ParseAndValid(model interface{}, body []byte) error {
	if err := json.Unmarshal(body, model); err != nil {
		return errors.Wrap(err, "Input parsing error")
	}
	valid := validation.Validation{}
	isValid, err := valid.Valid(model)
	if err != nil {
		return errors.Wrap(err, "Validation check failed")
	}
	if !isValid {
		builder := strings.Builder{}
		builder.WriteString("Data is invalid.\n")
		for _, er := range valid.Errors {
			builder.WriteString(er.Key)
			builder.WriteString(": ")
			builder.WriteString(er.Message)
			builder.WriteString("\n")
		}
		return fmt.Errorf(builder.String())
	}
	return nil
}
