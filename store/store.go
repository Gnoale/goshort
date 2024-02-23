package store

import (
	"context"
	"database/sql"
)

type Repository interface {
	CreateURL(context.Context, string) (int64, error)
	GetURL(context.Context, int64) (string, error)
}

type dbRepo struct {
	q *Queries
}

func NewRepository(db *sql.DB) Repository {
	return &dbRepo{New(db)}
}

func (db *dbRepo) CreateURL(ctx context.Context, url string) (int64, error) {
	return db.q.InsertURL(ctx, url)
}

func (db *dbRepo) GetURL(ctx context.Context, id int64) (string, error) {
	return db.q.GetURL(ctx, id)
}
