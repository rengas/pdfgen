package design

import (
	"context"
	"database/sql"
	"errors"
)

var ErrUnableToSaveDesign = errors.New("unable to save profile")

type DesignRepository struct {
	db *sql.DB
}

func NewDesignRepository(db *sql.DB) *DesignRepository {
	return &DesignRepository{
		db: db,
	}
}

func (r *DesignRepository) Save(ctx context.Context, p Design) error {
	_, err := r.db.ExecContext(ctx, "INSERT INTO design(id, profile_id, name, fields, template ) values($1, $2, $3, $4, $5)",
		p.Id,
		p.ProfileId,
		p.Name,
		p.Fields,
		p.Template)
	if err != nil {
		return ErrUnableToSaveDesign
	}

	return nil
}
