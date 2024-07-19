package shopify

import (
	"errors"
	"fmt"
)

type UserError struct {
	Code         *string  `json:"code"`
	ElementIndex *int     `json:"elementIndex"`
	Field        []string `json:"field"`
	Message      string   `json:"message"`
}

func (u *UserError) CodeString() string {
	if u.Code == nil {
		return ""
	}
	return *u.Code
}

func (u *UserError) Error() error {
	return fmt.Errorf("UserError: code: %s, field: %v, message: %s", u.CodeString(), u.Field, u.Message)
}

type UserErrors []UserError

func (u UserErrors) Error() error {
	errs := make([]error, 0, len(u))
	for _, userError := range u {
		errs = append(errs, userError.Error())
	}
	return errors.Join(errs...)
}
