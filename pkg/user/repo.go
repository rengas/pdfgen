package user

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rengas/pdfgen/pkg/pagination"
	"time"
)

var ErrUserNotFound = errors.New("user not found")

type Repository struct {
	db *sql.DB
}

func NewRepository(db *sql.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) SaveNewUser(ctx context.Context, u User) error {
	stmt, err := r.db.PrepareContext(ctx, "INSERT INTO users(id, email, password_hash, role,first_name, last_name,created_at,updated_at) values($1, $2, $3, $4, $5, $6, $7, $8)")
	if err != nil {
		return err
	}

	if _, err = stmt.ExecContext(ctx, u.Id, u.Email, u.PasswordHash, u.Role, u.FirstName, u.LastName, u.CreatedAt, u.UpdatedAt); err != nil {
		return err
	}
	return nil
}

func (r *Repository) GetByEmail(ctx context.Context, email string) (User, error) {
	var u User
	q := `SELECT id, email, password_hash, role, created_at, updated_at	
           FROM users WHERE email = $1`
	rw := r.db.QueryRowContext(ctx, q, email).
		Scan(&u.Id, &u.Email, &u.PasswordHash, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if rw != nil && rw.Error() != "" {
		return User{}, ErrUserNotFound
	}

	return u, nil
}

func (r *Repository) GetById(ctx context.Context, id string) (User, error) {
	var u User
	q := `SELECT id, email, password_hash, first_name, last_name, role, created_at, updated_at	
           FROM users WHERE id = $1 and deleted_at is NULL`
	rw := r.db.QueryRowContext(ctx, q, id).
		Scan(&u.Id, &u.Email, &u.PasswordHash, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
	if rw != nil && rw.Error() != "" {
		return User{}, ErrUserNotFound
	}

	return u, nil
}

func (r *Repository) Update(ctx context.Context, u User) error {

	q := `UPDATE users SET email=$2,first_name=$4, last_name=$5, updated_at=$6 WHERE id=$1`

	_, err := r.db.ExecContext(ctx, q, u.Id, u.Email, u.FirstName, u.LastName, u.UpdatedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateRole(ctx context.Context, id string, rl Role) error {
	q := `UPDATE users SET role=$2 WHERE id=$1`
	_, err := r.db.ExecContext(ctx, q, id, rl)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) DeleteById(ctx context.Context, id string, deletedAt time.Time) error {

	q := `UPDATE users SET deleted_at=$2 WHERE id=$1`

	_, err := r.db.ExecContext(ctx, q, id, deletedAt)
	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) List(ctx context.Context, limit, page int64) ([]User, pagination.Pagination, error) {
	query := `SELECT id, email, first_name, last_name, role, created_at, updated_at	
           FROM users WHERE deleted_at is NULL
           Limit $1 Offset $2
           `
	var rows *sql.Rows
	var err error
	offset := limit * (page - 1)

	rows, err = r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, pagination.Pagination{}, err
	}
	defer rows.Close()

	var us []User
	for rows.Next() {
		u := new(User)
		// works but I don't think it is good code for too many columns
		err = rows.Scan(&u.Id, &u.Email, &u.FirstName, &u.LastName, &u.Role, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, pagination.Pagination{}, err
		}
		us = append(us, *u)
	}

	qCount := `SELECT count(id)
	FROM users
	WHERE deleted_at is NULL`

	var count int64
	rw := r.db.QueryRowContext(ctx, qCount).Scan(&count)
	if rw != nil && rw.Error() != "" {
		return nil, pagination.Pagination{}, errors.New("unable to get count")
	}

	return us, pagination.Pagination{
		Page:  page,
		Total: count,
	}, nil
}
