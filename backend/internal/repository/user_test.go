package repository

import (
	"database/sql"
	"errors"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func newMockUserRepo(t *testing.T) (*User, sqlmock.Sqlmock, func()) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %v", err)
	}
	repo := NewUserRepo(db)
	cleanup := func() { db.Close() }
	return repo, mock, cleanup
}

func TestGetUserByEmail_Success(t *testing.T) {
	repo, mock, cleanup := newMockUserRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).
		AddRow("u1", "alice@example.com", "$2a$hash", "admin")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, role FROM users WHERE email = ?")).
		WithArgs("alice@example.com").
		WillReturnRows(rows)

	u, err := repo.GetUserByEmail("alice@example.com")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "u1", u.ID)
	assert.Equal(t, "alice@example.com", u.Email)
	assert.Equal(t, "admin", u.Role)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByEmail_NotFound(t *testing.T) {
	repo, mock, cleanup := newMockUserRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, role FROM users WHERE email = ?")).
		WithArgs("missing@example.com").
		WillReturnError(sql.ErrNoRows)

	u, err := repo.GetUserByEmail("missing@example.com")
	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserByEmail_DBError(t *testing.T) {
	repo, mock, cleanup := newMockUserRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, role FROM users WHERE email = ?")).
		WithArgs("err@example.com").
		WillReturnError(errors.New("db fail"))

	u, err := repo.GetUserByEmail("err@example.com")
	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserById_Success(t *testing.T) {
	repo, mock, cleanup := newMockUserRepo(t)
	defer cleanup()

	rows := sqlmock.NewRows([]string{"id", "email", "password_hash", "role"}).
		AddRow("u2", "bob@example.com", "$2a$hash2", "user")

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, role FROM users WHERE id = ?")).
		WithArgs("u2").
		WillReturnRows(rows)

	u, err := repo.GetUserById("u2")
	assert.NoError(t, err)
	assert.NotNil(t, u)
	assert.Equal(t, "u2", u.ID)
	assert.Equal(t, "bob@example.com", u.Email)

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserById_NotFound(t *testing.T) {
	repo, mock, cleanup := newMockUserRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, role FROM users WHERE id = ?")).
		WithArgs("missing").
		WillReturnError(sql.ErrNoRows)

	u, err := repo.GetUserById("missing")
	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "user not found")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}

func TestGetUserById_DBError(t *testing.T) {
	repo, mock, cleanup := newMockUserRepo(t)
	defer cleanup()

	mock.ExpectQuery(regexp.QuoteMeta("SELECT id, email, password_hash, role FROM users WHERE id = ?")).
		WithArgs("errid").
		WillReturnError(errors.New("db fail"))

	u, err := repo.GetUserById("errid")
	assert.Nil(t, u)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "db error")

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Fatalf("unfulfilled expectations: %v", err)
	}
}
