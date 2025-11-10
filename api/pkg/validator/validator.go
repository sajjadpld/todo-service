package validator

import (
	"context"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	"log"
)

var validate *validator.Validate

func init() {
	v, ok := binding.Validator.Engine().(*validator.Validate)
	if !ok {
		log.Fatal("[validator] failed to initialize")
	}

	validate = v
	registerCustomValidators()
}

func ValidateStruct(ctx context.Context, s interface{}) error {
	return validate.StructCtx(ctx, s)
}

func ValidateRequestDto(ctx context.Context, s interface{}) (err error) {
	if err = validate.StructCtx(ctx, s); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			for _, validationError := range validationErrors {
				errMsg := fmt.Sprintf("validation failed for the %s field.", validationError.Field())
				err = errors.New(errMsg)
			}

			return
		}

		return
	}

	return
}

// Var validates a single variable using tag style validation
func Var(item interface{}, validation string) (err error) {
	return validate.Var(item, validation)
}

// HELPERS

func registerCustomValidators() {
	var (
		err    error
		errMsg = "[validator] custom register err: %s"
	)

	if err = validate.RegisterValidation("jwt", isJwt); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err = validate.RegisterValidation("date", DateValidator, true); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err = validate.RegisterValidation("time", TimeValidator, true); err != nil {
		log.Fatalf(errMsg, err)
	}

	if err = validate.RegisterValidation("timeHourMinute", TimeHourMinuteValidator, true); err != nil {
		log.Fatalf(errMsg, err)
	}
}
