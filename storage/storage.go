package storage

import (
	"context"
	"frllo_xml/utils"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"sync"
)

var (
	pgOnce sync.Once
)

type Storage interface {
	CreateTemps() error
	GetDocuments(id int64) (pgx.Rows, error)
	GetBenefits(id string) (pgx.Rows, error)
	GetRecipes(ts int64) (pgx.Rows, error)
}

type PGStorage struct {
	ctx        context.Context
	db         *pgxpool.Pool
	t, i, b, r string
}

func NewPGStorage(ctx context.Context, uri string) (*PGStorage, error) {
	var pgStorage *PGStorage
	var (
		t, i, b, r string
	)
	t, i, b, r = utils.GetScripts()
	pgOnce.Do(func() {
		db, err := pgxpool.New(ctx, uri)
		if err != nil {
			log.Fatalf("Unable to connect to database: %v", err)
		}
		pgStorage = &PGStorage{ctx: ctx, db: db, t: t, i: i, b: b, r: r}

	})
	return pgStorage, nil
}

func (d *PGStorage) Close() {
	d.db.Close()
}

func (d *PGStorage) CreateTemps() error {
	query := d.t
	_, err := d.db.Exec(d.ctx, query)
	if err != nil {
		return err
	}
	return nil

}
func (d *PGStorage) GetDocuments(ts int64) (pgx.Rows, error) {
	query := d.i
	rows, err := d.db.Query(d.ctx, query, ts)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
func (d *PGStorage) GetBenefits(id string) (pgx.Rows, error) {
	query := d.b
	rows, err := d.db.Query(d.ctx, query, id)
	if err != nil {
		return nil, err
	}
	return rows, nil
}

func (d *PGStorage) GetRecipes(ts int64) (pgx.Rows, error) {
	query := d.r
	rows, err := d.db.Query(d.ctx, query, ts)
	if err != nil {
		return nil, err
	}
	return rows, nil
}
