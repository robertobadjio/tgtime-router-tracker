package router

import "github.com/robertobadjio/platform-common/pkg/db"

// PgRouterRepository ...
type PgRouterRepository struct {
	db db.Client
}

// NewPgRepository Конструктор PostgresQL репозитория.
func NewPgRepository(db db.Client) *PgRouterRepository {
	return &PgRouterRepository{db: db}
}
