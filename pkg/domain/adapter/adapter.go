package adapter

import (
	"sync"
)

type AdapterImagem interface {
	Upload(file []byte, id, name string, wg *sync.WaitGroup, ch chan error) error
	Read() (url string, err error)
}

type DecodeToken struct {
	Email  string `json:"email"`
	UserID string `json:"user_id"`
}

type AdapterToken interface {
	CreateToken(userID, email string) (string, error)
	DecodeToken(token string) (DecodeToken, error)
}

type Params map[string]string
type AdapterLogger interface {
	AddParam(key, value string)
	AddParams(params Params)
	Info(message string)
	Error(message string)
}
