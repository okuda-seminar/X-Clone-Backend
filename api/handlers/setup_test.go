package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"github.com/stretchr/testify/suite"
)

const (
	host     = "localhost"
	user     = "postgres"
	password = "postgres"
	dbname   = "handlers_test"
	port     = "5435"
	dsn      = "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable timezone=UTC connect_timeout=5"

	migrationFilesPath = "../../db/migrations"
)

var (
	opts dockertest.RunOptions
	pool *dockertest.Pool
)

func TestMain(m *testing.M) {
	var err error

	// connect to Docker.
	pool, err = dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	// set up our docker options, specifying the image and so forth.
	opts = dockertest.RunOptions{
		Repository: "postgres",
		Tag:        "14.5",
		Env: []string{
			"POSTGRES_USER=" + user,
			"POSTGRES_PASSWORD=" + password,
			"POSTGRES_DB=" + dbname,
		},
		ExposedPorts: []string{"5432"},
		PortBindings: map[docker.Port][]docker.PortBinding{
			"5432": {
				{HostIP: "0.0.0.0", HostPort: port},
			},
		},
	}

	// run tests.
	code := m.Run()

	os.Exit(code)
}

type HandlersTestSuite struct {
	suite.Suite
	db       *sql.DB
	resource *dockertest.Resource
}

// SetupTest runs before each test in the suite.
func (s *HandlersTestSuite) SetupTest() {
	var err error

	// pulls an image, creates a container based on it and runs it.
	s.resource, err = pool.RunWithOptions(&opts)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container
	// might not be ready to accept connections yet.
	if err := pool.Retry(func() error {
		var err error
		s.db, err = sql.Open("pgx", fmt.Sprintf(dsn, host, port, user, password, dbname))
		if err != nil {
			return err
		}
		return s.db.Ping()
	}); err != nil {
		log.Fatalf("Could not connect to database: %s", err)
	}

	// populate the database with empty tables.
	driver, err := postgres.WithInstance(s.db, &postgres.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	m, err := migrate.NewWithDatabaseInstance(fmt.Sprintf("file://%s", migrationFilesPath), dbname, driver)
	if err != nil {
		log.Fatalln(err)
	}

	m.Up()
}

// TearDownTest runs after each test in the suite.
func (s *HandlersTestSuite) TearDownTest() {
	s.db.Close()

	if err := pool.Purge(s.resource); err != nil {
		log.Fatalf("Could not purge resource: %s", err)
	}
}
