//go:build integration

package postgres_test

import (
	"context"
	"fmt"
	"os"
	"sync/atomic"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	testpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"
)

type postgresHarness struct {
	admin     *pgxpool.Pool
	base      *pgxpool.Config
	container *testpostgres.PostgresContainer
	sequence  atomic.Uint64
}

var integrationPostgres *postgresHarness

func TestMain(m *testing.M) {
	ctx := context.Background()
	connectionString := os.Getenv("SEREGA_INTEGRATION_DATABASE_URL")

	var container *testpostgres.PostgresContainer

	if connectionString == "" {
		var err error

		container, err = testpostgres.Run(ctx, "postgres:17.6-bookworm",
			testpostgres.WithDatabase("postgres"),
			testpostgres.WithUsername("serega"),
			testpostgres.WithPassword("serega"),
			testpostgres.BasicWaitStrategies(),
		)
		if err != nil {
			fmt.Fprintf(os.Stderr, "start PostgreSQL integration container: %v\n", err)
			os.Exit(1)
		}

		connectionString, err = container.ConnectionString(ctx, "sslmode=disable")
		if err != nil {
			_ = container.Terminate(ctx)

			fmt.Fprintf(os.Stderr, "read PostgreSQL integration connection string: %v\n", err)
			os.Exit(1)
		}
	}

	base, err := pgxpool.ParseConfig(connectionString)
	if err != nil {
		fmt.Fprintf(os.Stderr, "parse PostgreSQL integration connection string: %v\n", err)
		os.Exit(1)
	}

	admin, err := pgxpool.NewWithConfig(ctx, base.Copy())
	if err != nil {
		fmt.Fprintf(os.Stderr, "open PostgreSQL integration administrator: %v\n", err)
		os.Exit(1)
	}

	if err := admin.Ping(ctx); err != nil {
		admin.Close()
		fmt.Fprintf(os.Stderr, "ping PostgreSQL integration administrator: %v\n", err)
		os.Exit(1)
	}

	integrationPostgres = &postgresHarness{admin: admin, base: base, container: container}
	code := m.Run()

	admin.Close()

	if container != nil {
		if err := container.Terminate(ctx); err != nil && code == 0 {
			fmt.Fprintf(os.Stderr, "terminate PostgreSQL integration container: %v\n", err)

			code = 1
		}
	}

	os.Exit(code)
}
