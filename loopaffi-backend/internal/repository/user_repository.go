package repository

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/emirrasyad/loopaffi-backend/internal/entity"
	"github.com/jmoiron/sqlx"
)

type UserRepository struct {
	db *sqlx.DB
}

func NewUserRepository(db *sqlx.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) FindByEmail(email string) (*entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, name, email, password_hash, role_id, phone, status, created_at, updated_at FROM users WHERE email = $1", email)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by email: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) FindByID(id string) (*entity.User, error) {
	var user entity.User
	err := r.db.Get(&user, "SELECT id, name, email, password_hash, role_id, phone, status, created_at, updated_at FROM users WHERE id = $1", id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &user, nil
}

func (r *UserRepository) CreateUser(user *entity.User) error {
	if user.ID == "" {
		user.ID = fmt.Sprintf("USR-%d", time.Now().Unix())
	}

	query := `INSERT INTO users (id, name, email, password_hash, role_id, phone, status) 
			  VALUES ($1, $2, $3, $4, $5, $6, $7) 
			  RETURNING created_at, updated_at`
	return r.db.QueryRow(query, user.ID, user.Name, user.Email, user.PasswordHash, user.RoleID, user.Phone, user.Status).Scan(&user.CreatedAt, &user.UpdatedAt)
}

func (r *UserRepository) FindAll() ([]entity.User, error) {
	var users []entity.User
	err := r.db.Select(&users, "SELECT id, name, email, role_id, phone, status, created_at, updated_at FROM users ORDER BY created_at DESC")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch all users: %w", err)
	}
	return users, nil
}

func (r *UserRepository) UpdatePasswordByEmail(email string, hashedPassword string) error {
	result, err := r.db.Exec(
		"UPDATE users SET password_hash = $1, updated_at = NOW() WHERE email = $2",
		hashedPassword,
		email,
	)
	if err != nil {
		return fmt.Errorf("failed to update password: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to check updated rows: %w", err)
	}

	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}
