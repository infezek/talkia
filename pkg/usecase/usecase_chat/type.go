package usecase_chat

type CreateDtoInput struct {
	UserID string
	BotID  string
}

type UpdateDtoInput struct {
	ID     string
	Name   string
	Active bool
}

type CreateDtoOutput struct {
	ID string `json:"id"`
}

type UpdateDtoOutput CreateDtoOutput
