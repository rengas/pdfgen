package testutils

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"

	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/ory/dockertest/v3"
	"log"
	"os"
	"time"
)

func Cleanup(code int, pool *dockertest.Pool, network *dockertest.Network, resources ...*dockertest.Resource) {
	for _, resource := range resources {
		if resource != nil {
			if err := pool.Purge(resource); err != nil {
				log.Fatalf("failed to purge resource: %s", err.Error())
			}
		}
	}
	if network != nil {
		if err := network.Close(); err != nil {
			log.Fatalf("failed to close network: %s", err.Error())
		}
	}
	os.Exit(code)
}

func CreatePostgres(pool *dockertest.Pool, network *dockertest.Network, user, pass, dbName string) (*dockertest.Resource, error) {
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Name:       "pg",
		Repository: "postgres",
		Tag:        "14.2",
		Env: []string{
			"PGUSER=" + user,
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + pass,
			"POSTGRES_DB=" + dbName,
		},
		Networks: []*dockertest.Network{network},
	})
	if err != nil {
		return resource, fmt.Errorf("failed to run postgres: %w", err)
	}

	if err = pool.Retry(func() error {
		connStr := fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, pass, resource.GetPort("5432/tcp"), dbName)
		sql, err := MustOpenPostgres(connStr)
		if sql == nil {
			return errors.New("unable to open postgres")
		}
		return err
	}); err != nil {
		return resource, fmt.Errorf("failed to connect to postgres: %w", err)
	}

	return resource, nil
}

func MustOpenPostgres(connString string) (*sql.DB, error) {
	db, err := sql.Open("postgres", connString)
	if err != nil {
		return nil, err
	}
	if err = Retry(5, time.Second, 2, db.Ping); err != nil {
		return nil, fmt.Errorf("failed to ping postgres server")
	}

	return db, nil
}

type OperationFunc func() error

func Retry(attempts int, delay time.Duration, factor float64, op OperationFunc) error {
	var err error
	f := 1.
	for attempt := 1; attempt <= attempts; attempt++ {
		err = op()
		if err == nil {
			return nil
		}
		time.Sleep(time.Duration(float64(delay) * f))
		f *= factor
	}
	return err
}

func MigratePostgres(connStr, dbName, migrations string) error {
	db, err := MustOpenPostgres(connStr)
	if err != nil {
		return fmt.Errorf("failed to open postgres: %w", err)
	}

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return fmt.Errorf("failed to get driver instance: %w", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://"+migrations,
		dbName,
		driver,
	)
	if err != nil {
		return fmt.Errorf("failed to create migrate instance: %w", err)
	}

	return m.Up()
}
