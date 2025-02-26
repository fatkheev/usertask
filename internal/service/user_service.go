package service

import (
	"errors"
	"usertask/internal/auth"
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

func (s *UserService) CompleteTask(userID int, taskType string, points int) error {
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	// Выполняем задание и начисляем очки пользователю
	err = s.repo.CompleteTask(userID, taskType, points)
	if err != nil {
		return err
	}

	// Проверяем, есть ли у пользователя реферер
	referrerID, err := s.repo.GetUserReferrer(userID)
	if err != nil {
		return err
	}

	// Если у пользователя есть реферер, начисляем ему бонусные 50 очков
	if referrerID > 0 {
		err = s.repo.UpdateUserPoints(referrerID, 50)
		if err != nil {
			return err
		}
	}

	return nil
}

func (s *UserService) SetReferrer(userID int, referrerID int) error {
	// Нельзя указать самого себя в качестве реферера
	if userID == referrerID {
		return errors.New("нельзя указать самого себя как реферера")
	}

	// Проверяем, существует ли пользователь
	user, err := s.repo.GetUserByID(userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("пользователь не найден")
	}

	// Проверяем, что у пользователя ещё нет реферера
	if user.ReferrerID != nil {
		return errors.New("реферальный код уже установлен")
	}

	// Проверяем, существует ли реферер
	referrer, err := s.repo.GetUserByID(referrerID)
	if err != nil {
		return err
	}
	if referrer == nil {
		return errors.New("реферер не найден")
	}

	// Устанавливаем реферальный код и сразу начисляем бонус рефереру
	err = s.repo.SetUserReferrer(userID, referrerID)
	if err != nil {
		return err
	}

	// Начисляем бонус рефереру
	err = s.repo.UpdateUserPoints(referrerID, 50)
	if err != nil {
		return err
	}

	return nil
}

func (s *UserService) CreateUser(username string) (*models.User, string, error) {
	if username == "" {
		return nil, "", errors.New("имя пользователя не может быть пустым")
	}

	user, err := s.repo.CreateUser(username)
	if err != nil {
		return nil, "", err
	}

	// Генерируем JWT токен
	token, err := auth.GenerateToken(user.ID)
	if err != nil {
		return nil, "", err
	}

	return user, token, nil
}