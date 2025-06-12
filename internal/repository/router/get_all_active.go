package router

import (
	"context"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	"github.com/robertobadjio/platform-common/pkg/db"

	"github.com/robertobadjio/tgtime-router-tracker/internal/repository/model"
)

const tableName = "router"

const (
	idColumnName       = "id"
	nameColumnName     = "name"
	addressColumnName  = "address"
	loginColumnName    = "login"
	passwordColumnName = "password"
	statusColumnName   = "status"
)

// GetAllActive ...
func (r *PgRouterRepository) GetAllActive(ctx context.Context) ([]model.Router, error) {
	builder := sq.Select(idColumnName, nameColumnName, addressColumnName, loginColumnName, passwordColumnName).
		PlaceholderFormat(sq.Dollar).
		From(tableName).
		Where(sq.Eq{statusColumnName: true})

	query, args, err := builder.ToSql()
	if err != nil {
		return nil, fmt.Errorf("build query: %w", err)
	}

	q := db.Query{
		Name:     "router_repository.GetAllActive",
		QueryRaw: query,
	}

	var routers []model.Router
	err = r.db.DB().ScanAllContext(ctx, &routers, q, args...)
	if err != nil {
		return nil, err
	}

	return routers, nil
}
