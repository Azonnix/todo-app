package service

import (
	"crypto/sha1"
	"fmt"
	"time"

	"github.com/azonnix/todo-app"
	"github.com/azonnix/todo-app/pkg/repository"
	"github.com/dgrijalva/jwt-go"
)

const (
	sult      = "jaodighopaijfa759eoasl"
	signinKey = "nsbuh0839h98h4fgdsnvkjs"
	tocketTTL = 12 * time.Hour
)

type tockerClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
}

type AuthService struct {
	repo repository.Authorization
}

func NewAuthService(repo repository.Authorization) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) CreateUser(user todo.User) (int, error) {
	user.Password = generatePasswordHash(user.Password)
	return s.repo.CreateUser(user)
}

func (s *AuthService) GenerateTocken(username, password string) (string, error) {
	user, err := s.repo.GetUser(username, generatePasswordHash(password))
	if err != nil {
		return "", err
	}

	tocken := jwt.NewWithClaims(jwt.SigningMethodHS256, &tockerClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tocketTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return tocken.SignedString([]byte(signinKey))
}

func generatePasswordHash(password string) string {
	hash := sha1.New()
	hash.Write([]byte(password))
	return fmt.Sprintf("%x", hash.Sum([]byte(sult)))
}
