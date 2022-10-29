//go:build e2e

package test

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/ory/dockertest/v3"
	"github.com/rengas/pdfgen/pkg/logging"
	"github.com/rengas/pdfgen/pkg/testutils"
	"github.com/rengas/pdfgen/test/models"
	"log"
	"net/http"
	"os"
	"path"
	"testing"
	"time"
)

var httpPort string
var connStr string

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("failed to connect to docker: %s", err.Error())
	}

	network, err := pool.CreateNetwork("pdfgen")
	if err != nil {
		log.Fatalf("failed to create docker network: %s", err.Error())
	}

	// setup mysql container
	user, pass, dbName := "pdfgen", "pdfgen", "pdfgen"
	postgresContainer, err := testutils.CreatePostgres(pool, network, user, pass, dbName)
	if err != nil {
		log.Printf("failed to create postgres: %s", err.Error())
		testutils.Cleanup(1, pool, network, postgresContainer)
	}

	// run migrations
	pwd, err := os.Getwd()
	if err != nil {
		log.Fatalf("failed to get working directory: %s", err.Error())
	}

	projectRoot := path.Join(pwd, "..")
	migrations := path.Join(projectRoot, "migrations")

	connStr = fmt.Sprintf("postgres://%s:%s@localhost:%s/%s?sslmode=disable", user, pass, postgresContainer.GetPort("5432/tcp"), dbName)
	err = testutils.MigratePostgres(connStr, dbName, migrations)
	if err != nil {
		log.Printf("failed to migrate database: %s", err.Error())
		testutils.Cleanup(1, pool, network, postgresContainer)
	}

	// run tests
	time.Sleep(5 * time.Second)
	// setup consumer container
	apiContainer, err := createAPI(pool, network, projectRoot)
	if err != nil {
		log.Printf("failed to create log api: %s", err.Error())
		testutils.Cleanup(1, pool, network, postgresContainer, apiContainer)
	}

	testutils.Cleanup(m.Run(), pool, network, postgresContainer, apiContainer)
}

func createAPI(pool *dockertest.Pool, network *dockertest.Network, projectRoot string) (*dockertest.Resource, error) {
	buildOpts := &dockertest.BuildOptions{
		Dockerfile: "cmd/api/Dockerfile",
		ContextDir: projectRoot,
	}
	runOpts := &dockertest.RunOptions{
		Name:     "api",
		Networks: []*dockertest.Network{network},
	}

	resource, err := pool.BuildAndRunWithBuildOptions(buildOpts, runOpts)
	if err != nil {
		return resource, fmt.Errorf("failed to build and run: %w", err)
	}

	httpPort = resource.GetPort("8080/tcp")

	if err = pool.Retry(func() error {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%s/api/v1/health", httpPort))
		if err != nil {
			return err
		}

		if resp.StatusCode != 200 {
			return errors.New("got http error code")
		}

		return nil
	}); err != nil {
		return resource, fmt.Errorf("failed to connect to container: %w", err)
	}

	return resource, nil
}

type LoginFunc func() models.MainLoginResponse

func validLoginToken() models.MainLoginResponse {
	email := StringPtr(Email())
	password := StringPtr(Password(true, true, true, true, false, 8))

	b, err := json.Marshal(models.MainRegisterRequest{Email: email, Password: password})
	if err != nil {
		logging.Fatal(err.Error())
	}

	resp, err := http.Post(fmt.Sprintf("http://localhost:%s/api/v1/auth/register", httpPort), "application/json", bytes.NewReader(b))
	if err != nil {
		logging.Fatal(err.Error())
	}

	var rs models.MainRegisterRequest
	err = json.NewDecoder(resp.Body).Decode(&rs)

	b, err = json.Marshal(models.MainLoginRequest{Email: email, Password: password})

	resp, err = http.Post(fmt.Sprintf("http://localhost:%s/api/v1/auth/login", httpPort), "application/json", bytes.NewReader(b))
	if err != nil {
		logging.Fatal(err.Error())
	}

	var ls models.MainLoginResponse
	err = json.NewDecoder(resp.Body).Decode(&ls)
	ls.AccessToken = "Bearer" + " " + ls.AccessToken
	ls.RefreshToken = "Bearer" + " " + ls.RefreshToken
	return ls
}
