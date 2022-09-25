package account

import (
	"context"
	"database/sql"
	"errors"
)

var ErrUnableToSaveProfile = errors.New("unable to save profile")

type ProfileRepository struct {
	db *sql.DB
}

func NewProfileRepository(db *sql.DB) *ProfileRepository {
	return &ProfileRepository{
		db: db,
	}
}

func (r *ProfileRepository) Save(ctx context.Context, p Profile) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO profile(id, email,firebase_id, provider) values($1, $2, $3, $4)",
		p.Id,
		p.Email,
		p.FirebaseId,
		p.Provider)
	if err != nil {
		return ErrUnableToSaveProfile
	}

	return nil
}

func (r *ProfileRepository) GetByFirebaseId(ctx context.Context, id string) (Profile, error) {
	var p Profile
	err := r.db.QueryRowContext(ctx, "SELECT id, email, firebase_id, provider,created_at  FROM profile WHERE firebase_id = $1 and deleted_at is NULL", id).
		Scan(&p.Id, &p.Email, &p.FirebaseId, &p.Provider, &p.CreatedAt)
	if err != nil {
		return Profile{}, err
	}

	return p, nil
}

func (r *ProfileRepository) GetById(ctx context.Context, id string) (Profile, error) {
	var p Profile
	err := r.db.QueryRowContext(ctx, "SELECT id, email, firebase_id, provider,created_at  FROM profile WHERE id = $1 and deleted_at is NULL", id).
		Scan(&p.Id, &p.Email, &p.FirebaseId, &p.Provider, &p.CreatedAt)
	if err != nil {
		return Profile{}, err
	}

	return p, nil
}
