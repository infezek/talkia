package controller_bots

import "github.com/go-playground/validator/v10"

type Create struct {
	CategoryID  string `json:"category_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Personality string `json:"personality" validate:"required"`
	Description string `json:"description" validate:"required"`
	Location    string `json:"location" validate:"required"`
}

func validateParams(params interface{}) error {
	validate := validator.New()
	err := validate.Struct(params)
	if err != nil {
		return err
	}
	return nil
}
