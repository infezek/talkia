package usecase_category

import "time"

type CreateDtoInput struct {
	Name   string
	Active bool
}

type UpdateDtoInput struct {
	ID     string
	Name   string
	Active bool
}

type CreateDtoOutput struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Active    bool      `json:"active"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UpdateDtoOutput CreateDtoOutput
