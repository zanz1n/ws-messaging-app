package services

import (
	"errors"
	"strings"
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
	token := jwt.NewWithClaims(jwt.SigningMethodHS512, jwt.MapClaims{
		"id":       p.ID,
		"username": p.Username,
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
	})

	tokenString, err := token.SignedString([]byte(j.token))
	return tokenString, err
}

func (j *JwtService) ValidateToken(rp string) (*UserJwtPayload, error) {
	var (
		valId       string
		valUsername string
	)

	token, err := jwt.Parse(rp, func(t *jwt.Token) (interface{}, error) {
		var ok bool

		if _, ok = t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method: " + t.Header["alg"].(string))
		}

		if valId, ok = t.Claims.(jwt.MapClaims)["id"].(string); !ok {
			return nil, errors.New("invalid token payload")
		}

		if valUsername, ok = t.Claims.(jwt.MapClaims)["username"].(string); !ok {
			return nil, errors.New("invalid token payload")
		}

		if t.Claims.(jwt.MapClaims)["exp"].(float64) < float64(time.Now().Unix()) {
			return nil, errors.New("token is expired")
		}

		return []byte(j.token), nil
	})

	if err != nil || !token.Valid {
		return nil, errors.New(strings.ToLower(err.Error()))
	}

	return &UserJwtPayload{
		ID:       valId,
		Username: valUsername,
	}, nil
}
