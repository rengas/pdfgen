package design

import (
	"context"
	"database/sql"
	"errors"
	"github.com/rengas/pdfgen/pkg/pagination"
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
	_, err := r.db.ExecContext(ctx, "INSERT INTO design(id, user_id, name, fields, template ) values($1, $2, $3, $4, $5)",
		p.Id,
		p.UserId,
		p.Name,
		p.Fields,
		p.Template)
	if err != nil {
		return ErrUnableToSaveDesign
	}

	return nil
}

func (r *DesignRepository) GetById(ctx context.Context, userId, designId string) (Design, error) {
	var d Design
	err := r.db.QueryRowContext(ctx, "SELECT id,name,user_id, fields, template FROM design WHERE user_id = $1 and id =$2 and deleted_at is NULL", userId, designId).
		Scan(&d.Id, &d.UserId, &d.Name, &d.UserId, &d.Fields, &d.Template)
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

func (r *DesignRepository) Delete(ctx context.Context, userId, designId string) error {
	_, err := r.db.ExecContext(ctx, "Update design set deleted_at=$2 where user_id=$1 and id=$2",
		userId,
		designId,
		time.Now().UTC(),
	)
	if err != nil {
		return ErrUnableToSaveDesign
	}

	return nil
}

type ListQuery struct {
	Query  string
	UserId string
	Limit  int64
	Page   int64
}

func (r *DesignRepository) ListByUserId(ctx context.Context, lq ListQuery) ([]Design, pagination.Pagination, error) {
	query := `SELECT id, name, user_id, fields, template, created_at, updated_at 
			FROM design 
			WHERE user_id = $1 and deleted_at is NULL
			Limit $2 Offset $3`

	var rows *sql.Rows
	var err error
	offset := lq.Limit * (lq.Page - 1)
	rows, err = r.db.QueryContext(ctx, query, lq.UserId, lq.Limit, offset)
	if err != nil {
		return nil, pagination.Pagination{}, err
	}
	defer rows.Close()

	var ds []Design
	for rows.Next() {
		d := new(Design)
		// works but I don't think it is good code for too many columns
		err = rows.Scan(&d.Id, &d.Name, &d.UserId, &d.Fields, &d.Template, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, pagination.Pagination{}, err
		}
		ds = append(ds, *d)
	}

	qCount := `SELECT count(id)
	FROM design
	WHERE user_id = $1 and deleted_at is NULL`

	var count int64
	err = r.db.QueryRowContext(ctx, qCount, lq.UserId).Scan(&count)
	if err != nil {
		return nil, pagination.Pagination{}, err
	}

	return ds, pagination.Pagination{
		Page:  lq.Page,
		Total: count,
	}, nil
}

func (r *DesignRepository) Search(ctx context.Context, lq ListQuery) ([]Design, pagination.Pagination, error) {
	query := `SELECT id, name, user_id, fields, template, created_at, updated_at
			FROM design
			WHERE user_id = $1
			and  deleted_at is NULL
			and LOWER(name) LIKE '%' || $2 || '%'
			Limit $3 Offset $4`

	var rows *sql.Rows
	var err error
	offset := lq.Limit * (lq.Page - 1)
	rows, err = r.db.QueryContext(ctx, query, lq.UserId, lq.Query, lq.Limit, offset)
	if err != nil {
		return nil, pagination.Pagination{}, err
	}
	defer rows.Close()

	var ds []Design
	for rows.Next() {
		d := new(Design)
		// works but I don't think it is good code for too many columns
		err = rows.Scan(&d.Id, &d.Name, &d.UserId, &d.Fields, &d.Template, &d.CreatedAt, &d.UpdatedAt)
		if err != nil {
			return nil, pagination.Pagination{}, err
		}
		ds = append(ds, *d)
	}

	qCount := `SELECT count(id)
	FROM design
	WHERE user_id = $1
	and  deleted_at is NULL
	and LOWER(name) LIKE '%' || $2 || '%'
	AND name like '%$2'`

	var count int64
	err = r.db.QueryRowContext(ctx, qCount, lq.UserId, lq.Query).Scan(&count)
	if err != nil {
		return nil, pagination.Pagination{}, err
	}
	return ds, pagination.Pagination{
		Page:  lq.Page,
		Total: count,
	}, nil
}
