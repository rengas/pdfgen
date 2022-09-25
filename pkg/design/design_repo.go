package design

import (
	"context"
	"database/sql"
	"errors"
	"time"
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

func (r *DesignRepository) GetById(ctx context.Context, id string) (Design, error) {
	var d Design
	err := r.db.QueryRowContext(ctx, "SELECT id,name,profile_id, fields, template FROM design WHERE id = $1", id).
		Scan(&d.Id, &d.Name, &d.ProfileId, &d.Fields, &d.Template)
	if err != nil {
		return Design{}, err
	}

	return d, nil
}

func (r *DesignRepository) Update(ctx context.Context, p Design) error {
	_, err := r.db.ExecContext(ctx, "Update design set name=$2, fields=$3, template=$4, updated_at=$5 where id=$1",
		p.Id,
		p.Name,
		p.Fields,
		p.Template,
		p.UpdatedAt,
	)
	if err != nil {
		return ErrUnableToSaveDesign
	}

	return nil
}

func (r *DesignRepository) Delete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, "Update design set deleted_at=$2 where id=$1",
		id,
		time.Now().UTC(),
	)
	if err != nil {
		return ErrUnableToSaveDesign
	}

	return nil
}

type Pagination struct {
	Total int64
	Page  int64
}

type ListQuery struct {
	Query     string
	ProfileId string
	Limit     int64
	Page      int64
}

func (r *DesignRepository) ListByProfileId(ctx context.Context, lq ListQuery) ([]Design, Pagination, error) {
	query := `SELECT id, name, profile_id, fields, template, created_at, updated_at 
			FROM design 
			WHERE profile_id = $1
			Limit $2 Offset $3`

	var rows *sql.Rows
	var err error
	offset := lq.Limit * (lq.Page - 1)
	rows, err = r.db.QueryContext(ctx, query, lq.ProfileId, lq.Limit, offset)
	if err != nil {
		return nil, Pagination{}, err
	}
	defer rows.Close()

	var ds []Design
	for rows.Next() {
		d := new(Design)
		// works but I don't think it is good code for too many columns
		err = rows.Scan(&d.Id, &d.Name, &d.ProfileId, &d.Fields, &d.Template, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, Pagination{}, err
		}
		ds = append(ds, *d)
	}

	qCount := `SELECT count(id)
	FROM design
	WHERE profile_id = $1`

	var count int64
	err = r.db.QueryRowContext(ctx, qCount, lq.ProfileId).Scan(&count)
	if err != nil {
		return nil, Pagination{}, err
	}

	return ds, Pagination{
		Page:  lq.Page,
		Total: count,
	}, nil
}

func (r *DesignRepository) Search(ctx context.Context, lq ListQuery) ([]Design, Pagination, error) {
	query := `SELECT id, name, profile_id, fields, template, created_at, updated_at
			FROM design
			WHERE profile_id = $1
			and LOWER(name) LIKE '%' || $2 || '%'
			Limit $3 Offset $4`

	var rows *sql.Rows
	var err error
	offset := lq.Limit * (lq.Page - 1)
	rows, err = r.db.QueryContext(ctx, query, lq.ProfileId, lq.Query, lq.Limit, offset)
	if err != nil {
		return nil, Pagination{}, err
	}
	defer rows.Close()

	var ds []Design
	for rows.Next() {
		d := new(Design)
		// works but I don't think it is good code for too many columns
		err = rows.Scan(&d.Id, &d.Name, &d.ProfileId, &d.Fields, &d.Template, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, Pagination{}, err
		}
		ds = append(ds, *d)
	}

	qCount := `SELECT count(id)
	FROM design
	WHERE profile_id = $1
	AND name like '%$2'`

	var count int64
	err = r.db.QueryRowContext(ctx, qCount, lq.ProfileId, lq.Query).Scan(&count)
	if err != nil {
		return nil, Pagination{}, err
	}
	return ds, Pagination{
		Page:  lq.Page,
		Total: count,
	}, nil
}
