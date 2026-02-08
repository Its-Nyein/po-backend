package services

import (
	"errors"
	"po-backend/helper"
	"po-backend/models"
	"po-backend/repositories"
)

type UserService struct {
	Repo *repositories.UserRepository
}

func NewUserService(repo *repositories.UserRepository) *UserService {
	return &UserService{Repo: repo}
}

func (s *UserService) GetAll() ([]models.User, error) {
	return s.Repo.GetAll(20)
}

func (s *UserService) GetByID(id uint) (*models.User, error) {
	return s.Repo.GetByID(id)
}

func (s *UserService) Register(name, username, bio, password string) (*models.User, error) {
	hashedPassword, err := helper.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &models.User{
		Name:     name,
		Username: username,
		Bio:      bio,
		Password: hashedPassword,
	}

	if err := s.Repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) Login(username, password string) (*models.User, string, error) {
	user, err := s.Repo.GetByUsername(username)
	if err != nil {
		return nil, "", errors.New("incorrect username or password")
	}

	if !helper.CheckPasswordHash(password, user.Password) {
		return nil, "", errors.New("incorrect username or password")
	}

	token, err := helper.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}

func (s *UserService) GetByUsername(username string) (*models.User, error) {
	return s.Repo.GetByUsername(username)
}

func (s *UserService) Search(query string) ([]models.User, error) {
	return s.Repo.Search(query, 20)
}
