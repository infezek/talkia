package adapter

import (
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/infezek/app-chat/pkg/config"
	"github.com/infezek/app-chat/pkg/domain/adapter"
)

type AdapterToken struct {
	Cfg *config.Config
}

func NewToken(cfg *config.Config) *AdapterToken {
	return &AdapterToken{
		Cfg: cfg,
	}
}

func (t *AdapterToken) CreateToken(userID, email string) (string, error) {
	claims := jwt.MapClaims{
		"email":   email,
		"user_id": userID,
		"exp":     time.Now().Add(time.Hour * 336).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(t.Cfg.OpenIAToken))
	if err != nil {
		return "", err
	}
	return signed, nil
}

func (t *AdapterToken) DecodeToken(tokenString string) (adapter.DecodeToken, error) {
	tokenSplit := strings.Split(tokenString, " ")
	if len(tokenSplit) >= 2 {
		tokenString = tokenSplit[1]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(t.Cfg.OpenIAToken), nil
	})
	if err != nil {
		return adapter.DecodeToken{}, fmt.Errorf("error parsing token")
	}
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return adapter.DecodeToken{}, fmt.Errorf("error parsing token")
	}
	return adapter.DecodeToken{
		Email:  claims["email"].(string),
		UserID: claims["user_id"].(string),
	}, nil
}
