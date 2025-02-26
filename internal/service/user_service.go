package service

import (
	"errors"
	"usertask/internal/models"
	"usertask/internal/repository"
)

type UserService struct {
	repo repository.UserRepositoryInterface
}

func NewUserService(repo repository.UserRepositoryInterface) *UserService {
	return &UserService{repo: repo}
}

func (s *UserService) GetUserStatus(userID int) (*models.User, error) {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("пользователь не найден")
	}
	return user, nil
}

func (s *UserService) CompleteTask(userID int, points int) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	return s.repo.UpdateUserPoints(userID, points)
}

func (s *UserService) SetReferrer(userID int, referrerID int) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	referrer, err := s.repo.GetUserByID(referrerID)
	if err != nil {
		return err
	}
	if referrer == nil {
		return errors.New("реферер не найден")
	}

	return s.repo.SetUserReferrer(userID, referrerID)
}

func (s *UserService) CreateUser(username string) (*models.User, error) {
	if username == "" {
		return nil, errors.New("имя пользователя не может быть пустым")
	}

	user, err := s.repo.CreateUser(username)
	if err != nil {
		return nil, err
	}

	return user, nil
}