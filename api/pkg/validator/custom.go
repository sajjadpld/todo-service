package validator

import (
	goValidator "github.com/go-playground/validator/v10"
	"regexp"
	"time"
)

const (
	DateOnly       = "2006-01-02"
	TimeOnly       = "15:04:05"
	HourMinuteOnly = "15:04"
)

func isJwt(fl goValidator.FieldLevel) (res bool) {
	value := fl.Field().String()

	if len(value) == 0 {
		return
	}

	res, err := regexp.MatchString(`^[A-Za-z0-9-_]+?\.[A-Za-z0-9-_]+?\.[A-Za-z0-9-_]+$`, value)
	if err != nil {
		res = false
	}

	return
}

func DateValidator(field goValidator.FieldLevel) bool {
	if dateStr, ok := field.Field().Interface().(string); ok {
		_, err := time.Parse(DateOnly, dateStr)
		return err == nil
	}
	return false
}

func TimeHourMinuteValidator(field goValidator.FieldLevel) bool {
	if field.Field().Interface().(string) == "" {
		return true
	}
	if timeStr, ok := field.Field().Interface().(string); ok {
		_, err := time.Parse(HourMinuteOnly, timeStr)
		return err == nil
	}
	return false
}

func TimeValidator(fl goValidator.FieldLevel) bool {
	t := fl.Field().String()
	_, err := time.Parse(TimeOnly, t)
	if err != nil {
		return false
	}
	return true
}
