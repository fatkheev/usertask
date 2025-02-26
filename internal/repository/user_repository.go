package repository

import (
	"database/sql"
	"errors"
	"usertask/internal/models"
)

type UserRepositoryInterface interface {
	GetUserByID(userID int) (*models.User, error)
	CreateUser(username string) (*models.User, error)
	UpdateUserPoints(userID int, points int) error
	SetUserReferrer(userID int, referrerID int) error
}

type UserRepository struct {
	db *sql.DB
}

var _ UserRepositoryInterface = (*UserRepository)(nil)

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) GetUserByID(userID int) (*models.User, error) {
	var user models.User
	query := `SELECT id, username, points, referrer_id, created_at FROM users WHERE id = $1`
	err := r.db.QueryRow(query, userID).Scan(&user.ID, &user.Username, &user.Points, &user.ReferrerID, &user.CreatedAt)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(username string) (*models.User, error) {
	query := `INSERT INTO users (username) VALUES ($1) RETURNING id, username, points, referrer_id, created_at`
	var user models.User
	err := r.db.QueryRow(query, username).Scan(&user.ID, &user.Username, &user.Points, &user.ReferrerID, &user.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) UpdateUserPoints(userID int, points int) error {
	query := `UPDATE users SET points = points + $1 WHERE id = $2`
	_, err := r.db.Exec(query, points, userID)
	return err
}

func (r *UserRepository) SetUserReferrer(userID int, referrerID int) error {
	query := `UPDATE users SET referrer_id = $1 WHERE id = $2 AND referrer_id IS NULL`
	result, err := r.db.Exec(query, referrerID, userID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected == 0 {
		return errors.New("реферальный код уже установлен или пользователя не существует")
	}

	return nil
}
