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
	CompleteTask(userID int, taskType string, points int) error
	GetUserReferrer(userID int) (int, error)
	GetLeaderboard(limit int) ([]models.User, error)
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

func (r *UserRepository) CompleteTask(userID int, taskType string, points int) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}

	// Записываем выполнение задания
	_, err = tx.Exec(`INSERT INTO tasks (user_id, task_type, points) VALUES ($1, $2, $3)`,
		userID, taskType, points)
	if err != nil {
		tx.Rollback()
		return err
	}

	// Обновляем очки пользователя
	_, err = tx.Exec(`UPDATE users SET points = points + $1 WHERE id = $2`, points, userID)
	if err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit()
}

func (r *UserRepository) GetUserReferrer(userID int) (int, error) {
	var referrerID sql.NullInt64

	err := r.db.QueryRow(`SELECT referrer_id FROM users WHERE id = $1`, userID).Scan(&referrerID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return 0, nil
		}
		return 0, err
	}

	if referrerID.Valid {
		return int(referrerID.Int64), nil
	}
	return 0, nil
}

func (r *UserRepository) GetLeaderboard(limit int) ([]models.User, error) {
	rows, err := r.db.Query(`
		SELECT id, username, points, referrer_id, created_at
		FROM users
		ORDER BY points DESC
		LIMIT $1`, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(&user.ID, &user.Username, &user.Points, &user.ReferrerID, &user.CreatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
