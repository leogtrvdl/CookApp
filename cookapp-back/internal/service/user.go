package service

import (
	"cookapp/internal/model"
	"cookapp/internal/repository"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService(repo *repository.UserRepository) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) Register(username string, email string, password string) error {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	err = s.repo.CreateUser(username, email, string(hashedPassword))
	return err
}

func (s *UserService) Login(username string, password string) (string, error) {
	secretKey := os.Getenv("JWT_SECRET")

	user, err := s.repo.GetUserByUsername(username)
	if err != nil {
		return "", err
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"id":       user.ID,
		"username": username,
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
	})
	tokenString, err := token.SignedString([]byte(secretKey))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *UserService) GetUserByID(userID int) (*model.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) UpdateUser(userID int, username string, email string, bio string) error {
	err := s.repo.UpdateUser(userID, username, email, bio)
	return err
}

func (s *UserService) DeleteUser(userID int) error {
	err := s.repo.DeleteUser(userID)
	return err
}

func (s *UserService) UpdateAvatar(userID int, avatarPath string) error {
	return s.repo.UpdateAvatar(userID, avatarPath)
}

func (s *UserService) SearchUsers(username string) ([]*model.User, error) {
	user, err := s.repo.SearchUsers(strings.ToLower(username))
	return user, err
}
