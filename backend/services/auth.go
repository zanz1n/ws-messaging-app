package services

import (
	"context"
	"errors"
	"log"

	"github.com/zanz1n/ws-messaging-app/internal/dba"
	"golang.org/x/crypto/bcrypt"
)

var (
	dbCtx = context.Background()
)

type AuthService struct {
	db  *dba.Queries
	jwt *JwtService
}

func NewAuthService(db *dba.Queries, jwt *JwtService) *AuthService {
	return &AuthService{
		db:  db,
		jwt: jwt,
	}
}

func (ap *AuthService) GenerateHash(passwd string) string {
	hash, err := bcrypt.GenerateFromPassword([]byte(passwd), ConfigProvider().BcryptSalt)

	if err != nil {
		log.Printf("error generating hash: %v", err)
	}

	return string(hash)
}

func (ap *AuthService) AuthenticateUser(name string, passwd string) (string, error) {
	user, err := ap.db.GetUserByUsername(dbCtx, name)

	if err != nil {
		return "", errors.New("user and password do not match")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(passwd))

	if err != nil {
		return "", errors.New("user and password do not match")
	}

	token, err := ap.jwt.GenerateToken(&UserJwtPayload{
		ID:       user.ID,
		Username: user.Username,
	})

	return token, err
}

func (ap *AuthService) ValidateJwtToken(token string) (*UserJwtPayload, error) {
	return ap.jwt.ValidateToken(token)
}
