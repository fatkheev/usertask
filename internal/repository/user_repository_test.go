package repository

import (
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{db: db}

	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "points", "referrer_id", "created_at"}).
		AddRow(1, "testuser", 100, nil, createdAt)

	mock.ExpectQuery(`SELECT id, username, points, referrer_id, created_at FROM users WHERE id = \$1`).
		WithArgs(1).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "testuser", user.Username)
	assert.WithinDuration(t, createdAt, user.CreatedAt, time.Second)
}

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := &UserRepository{db: db}

	createdAt := time.Now()

	rows := sqlmock.NewRows([]string{"id", "username", "points", "referrer_id", "created_at"}).
		AddRow(2, "newuser", 0, nil, createdAt)

	mock.ExpectQuery(`INSERT INTO users`).
		WithArgs("newuser").
		WillReturnRows(rows)

	user, err := repo.CreateUser("newuser")
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "newuser", user.Username)
	assert.WithinDuration(t, createdAt, user.CreatedAt, time.Second)
}
