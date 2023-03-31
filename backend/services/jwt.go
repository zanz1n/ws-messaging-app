package services

import (
	"errors"
	"log"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type UserJwtPayload struct {
	ID       string
	Username string
}

type JwtService struct {
	token string
}

func NewJwtService() *JwtService {
	return &JwtService{
		token: ConfigProvider().JwtSecret,
	}
}

func (j *JwtService) GenerateToken(p *UserJwtPayload) (string, error) {
	token := jwt.New(jwt.SigningMethodEdDSA)
	claims := token.Claims.(jwt.MapClaims)

	claims["id"] = p.ID
	claims["username"] = p.Username
	claims["exp"] = time.Now().Add(time.Hour * 24 * 7)

	return token.SignedString(j.token)
}

func (j *JwtService) ValidateToken(rp string) (*UserJwtPayload, error) {
	token, err := jwt.Parse(rp, func(t *jwt.Token) (interface{}, error) {
		return "", nil
	})

	if err == nil && token.Valid {
		log.Println(token.Claims)
		return nil, errors.New("invalid jwt token")
	}

	return nil, errors.New("invalid jwt token")
}
