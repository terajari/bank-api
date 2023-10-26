package repository

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/terajari/bank-api/model"
)

type UsersRepository interface {
	Create(ctx context.Context, user model.Users) (model.Users, error)
	Get(ctx context.Context, username string) (model.Users, error)
	Update(ctx context.Context, user model.Users) (model.Users, error)
}

type userRepository struct {
	db *sqlx.DB
}

func NewUsersRepository(db *sqlx.DB) UsersRepository {
	return &userRepository{
		db: db,
	}
}

func (u *userRepository) Create(ctx context.Context, user model.Users) (model.Users, error) {
	query := "INSERT INTO users (username, hashed_password, full_name, email) VALUES ($1, $2, $3, $4) RETURNING username, hashed_password, full_name, email, password_changed_at, created_at"
	row := u.db.QueryRowContext(ctx, query, user.Username, user.HashedPassword, user.FullName, user.Email)
	var us model.Users
	err := row.Scan(&us.Username, &us.HashedPassword, &us.FullName, &us.Email, &us.PasswordChangedAt, &us.CreatedAt)
	if err != nil {
		return model.Users{}, err
	}
	return us, nil
}

func (u *userRepository) Get(ctx context.Context, username string) (model.Users, error) {
	query := `SELECT username, hashed_password, full_name, email, password_changed_at, created_at FROM users WHERE username = $1 LIMIT 1`
	row := u.db.QueryRowContext(ctx, query, username)
	var us model.Users
	err := row.Scan(&us.Username, &us.HashedPassword, &us.FullName, &us.Email, &us.PasswordChangedAt, &us.CreatedAt)
	if err != nil {
		return model.Users{}, err
	}
	return us, nil
}

func (u *userRepository) Update(ctx context.Context, user model.Users) (model.Users, error) {
	query := `UPDATE users SET hashed_password = $2, full_name = $3, email = $4 WHERE username = $1 RETURNING username, hashed_password, full_name, email, password_changed_at, created_aUPDATE users
	SET
	  hashed_password = COALESCE($1, hashed_password),
	  password_changed_at = COALESCE($2, password_changed_at),
	  full_name = COALESCE($3, full_name),
	  email = COALESCE($4, email),
	WHERE
	  username = $5
	RETURNING username, hashed_password, full_name, email, password_changed_at, created_at;
	`
	row := u.db.QueryRowContext(ctx, query, user.Username, user.HashedPassword, user.FullName, user.Email)
	var us model.Users
	err := row.Scan(&us.Username, &us.HashedPassword, &us.FullName, &us.Email, &us.PasswordChangedAt, &us.CreatedAt)
	if err != nil {
		return model.Users{}, err
	}
	return us, nil
}
