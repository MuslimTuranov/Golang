package usecase

import (
	"Assignment2/internal/repository"
	"Assignment2/pkg/modules"
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthUsecase struct {
	repo   repository.UserRepository
	ttl    time.Duration
	secret string
}

func NewAuthUsecase(repo repository.UserRepository, secret string, ttl time.Duration) *AuthUsecase {
	return &AuthUsecase{repo: repo, secret: secret, ttl: ttl}
}

func (a *AuthUsecase) Login(ctx context.Context, id int, password string) (string, error) {
	u, err := a.repo.GetUserByID(ctx, id)
	if err != nil {
		return "", modules.ErrUnauthorized
	}
	if password == "" || bcrypt.CompareHashAndPassword([]byte(u.PasswordHash), []byte(password)) != nil {
		return "", modules.ErrUnauthorized
	}

	claims := jwt.MapClaims{
		"user_id": id,
		"exp":     time.Now().Add(a.ttl).Unix(),
	}
	tkn := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tkn.SignedString([]byte(a.secret))
	if err != nil {
		return "", errors.New("token sign failed")
	}
	return s, nil
}
