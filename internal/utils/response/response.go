package response

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/go-playground/validator/v10"

)

type errorResponse struct {
	Status string `json:"status"`
	Error  string `json:"error"`
}
type successResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
type successResponsePlusData struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	Data interface{} `json:"data"`
}

const (
	Success = "success"
	Error   = "error"
)

func GeneralError(err error) *errorResponse {
	return &errorResponse{
		Status: Error,
		Error:  err.Error(),
	}
}

func GeneralSuccess(success string) *successResponse {
	return &successResponse{
		Status:  Success,
		Message: success,
	}
}


func GeneralSuccessPlusData (success string,data interface{}) *successResponsePlusData {
	return &successResponsePlusData{
		Status: Success,
		Message: success,
		Data: data,
	}
}

func ValidationErrors(errs validator.ValidationErrors) *errorResponse {
	var errors []string
	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errors = append(errors, fmt.Sprintf("field %s is required", err.Field()))
		case "email":
			errors = append(errors, fmt.Sprintf("field %s must be a valid email address", err.Field()))
		case "min":
			errors = append(errors, fmt.Sprintf("field %s must be at least %s characters long", err.Field(), err.Param()))
		case "max":
			errors = append(errors, fmt.Sprintf("field %s must be at most %s characters long", err.Field(), err.Param()))
		case "oneof":
			errors = append(errors, fmt.Sprintf("field %s must be one of %s", err.Field(), err.Param()))
		case "strongepwd":
			errors = append(errors, fmt.Sprintf("%s should have uppper,lower,number and special chars & min 8 chars", err.Tag()))
		}
	}
	return &errorResponse{
		Status: Error,
		Error:  strings.Join(errors, ","),
	}
}

func IsStrongePassword(f1 validator.FieldLevel) bool {
	pass := f1.Field().String()

	hasUpper := regexp.MustCompile(`[A-Z]`).MatchString(pass)
	hasLower := regexp.MustCompile(`[a-z]`).MatchString(pass)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(pass)
	hasSpecial := regexp.MustCompile(`[!@#$%^&*()_+\-=\[\]{};':"\\|,.<>/?]`).MatchString(pass)

	if len(pass) >= 8 && hasLower && hasNumber && hasSpecial && hasUpper {
		return true
	}
	return false
}





