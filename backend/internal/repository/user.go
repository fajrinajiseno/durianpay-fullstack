package repository

import (
	"database/sql"

	"github.com/fajrinajiseno/mygolangapp/internal/entity"
)

//go:generate mockgen -source user.go -destination mocks/user_mock.go -package=mock
type UserRepository interface {
	GetUserByEmail(email string) (*entity.User, error)
	GetUserById(id string) (*entity.User, error)
}

type User struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) *User {
	return &User{db: db}
}

func (r *User) GetUserByEmail(email string) (*entity.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, role FROM users WHERE email = ?`, email)
	var u entity.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrorNotFound("user not found")
		}
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}
	return &u, nil
}

func (r *User) GetUserById(id string) (*entity.User, error) {
	row := r.db.QueryRow(`SELECT id, email, password_hash, role FROM users WHERE id = ?`, id)
	var u entity.User
	if err := row.Scan(&u.ID, &u.Email, &u.PasswordHash, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, entity.ErrorNotFound("user not found")
		}
		return nil, entity.WrapError(err, entity.ErrorCodeInternal, "db error")
	}
	return &u, nil
}
