package repositor

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/robertobadjio/platform-common/pkg/db"
	"github.com/robertobadjio/platform-common/pkg/db/pg"
	"github.com/robertobadjio/tgtime-router-tracker/internal/repository/model"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/golang-migrate/migrate/v4"
	postgresMigrate "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	routerRepo "github.com/robertobadjio/tgtime-router-tracker/internal/repository/router"
)

const (
	database = "test"
	user     = "user"
	password = "password"
)

const dbMigrationPath = "./migrations"
const dbSchemaName = "test"

func TestIntegration_GetAllActive(t *testing.T) {
	ctx := context.Background()
	postgresContainer, err := postgres.Run(
		ctx,
		"postgres:16-alpine",
		postgres.WithInitScripts(filepath.Join("./docker-entrypoint-initdb.d", "init.sql")),
		postgres.WithDatabase(database),
		postgres.WithUsername(user),
		postgres.WithPassword(password),
		postgres.BasicWaitStrategies(),
	)
	if err != nil {
		log.Fatal(err)
	}

	state, errGetState := postgresContainer.State(ctx)
	require.NoError(t, errGetState)

	t.Log("State running:", state.Running)

	connStr, errConnectionString := postgresContainer.ConnectionString(ctx, "sslmode=disable", "application_name="+database)
	require.NoError(t, errConnectionString)

	fmt.Println(connStr)

	conn, errOpen := sql.Open(
		"postgres",
		connStr,
	)
	if errOpen != nil {
		log.Fatalf("errOpen: %v", errOpen)
	}

	if errMigration := runMigrations(conn); errMigration != nil {
		log.Fatalf("error migration: %v", errMigration)
	}

	PGRouterRepo := routerRepo.NewPgRepository(DBClient(ctx, connStr))

	routers, err := PGRouterRepo.GetAllActive(context.Background())
	if err != nil {
		t.Errorf("GetAllActive failed: %v", err)
	}

	assert.Len(t, routers, 1)

	expectedRouter := model.Router{
		ID:       1,
		Name:     "Router1",
		Address:  "95.84.134.115:8728",
		Login:    "admin",
		Password: "Vtlcgjgek1",
	}

	assert.Equal(t, expectedRouter, routers[0])

	testcontainers.CleanupContainer(t, postgresContainer)
	require.NoError(t, err)
}

func runMigrations(db *sql.DB) error {
	driver, errWithInstance := postgresMigrate.WithInstance(db, &postgresMigrate.Config{})
	if errWithInstance != nil {
		return errWithInstance
	}

	m, errNewWithDatabaseInstance := migrate.NewWithDatabaseInstance(
		"file://"+dbMigrationPath,
		dbSchemaName,
		driver,
	)
	if errNewWithDatabaseInstance != nil {
		return errNewWithDatabaseInstance
	}

	err := m.Up()
	if err != nil {
		return err
	}

	return nil
}

func DBClient(ctx context.Context, connSrt string) db.Client {
	cl, err := pg.New(
		ctx,
		connSrt,
		1*time.Second,
	)
	if err != nil {
		log.Fatalf("failed to create db client: %v", err)
	}

	err = cl.DB().Ping(ctx)
	if err != nil {
		log.Fatalf("ping error: %s", err.Error())
	}

	return cl
}
