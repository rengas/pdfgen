package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"github.com/go-chi/chi/v5"
	m "github.com/go-chi/chi/v5/middleware"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rengas/pdfgen/pkg/account"
	"github.com/rengas/pdfgen/pkg/dbutils"
	"github.com/rengas/pdfgen/pkg/design"
	"github.com/rengas/pdfgen/pkg/firebase"
	"github.com/rengas/pdfgen/pkg/middleware"
	"github.com/rengas/pdfgen/pkg/minifier"
	"github.com/rengas/pdfgen/pkg/pdfrender"
	"github.com/rengas/pdfgen/pkg/server"
	"github.com/rengas/pdfgen/pkg/service"
	"io"
	"log"
	"os"
	"syscall"
	"time"
)

var (
	addr            = flag.String("addr", ":8080", "Application http server network address")
	shutdownTimeout = flag.Duration("shutdown-timeout", 30*time.Second, "Graceful shutdown timeout")
	connString      = flag.String("pg-conn-string", "postgres://pdfgen:pdfgen@localhost:5432/pdfgen?sslmode=disable", "PostgresSQL server connection string")
	firebasePath    = flag.String("firebase-path", "./firebase.json", "firebase creds")
	renderDBPath    = flag.String("render-db-path", "/etc/secrets/staging.json", "staging db creds")
)

type ProfileRepository interface {
	Save(ctx context.Context, p account.Profile) error
	GetById(ctx context.Context, id string) (account.Profile, error)
}

type DesignRepository interface {
	Save(ctx context.Context, d design.Design) error
	GetById(ctx context.Context, id string) (design.Design, error)
	Update(ctx context.Context, p design.Design) error
	Delete(ctx context.Context, id string) error
	ListByProfileId(ctx context.Context, lq design.ListQuery) ([]design.Design, design.Pagination, error)
	Search(ctx context.Context, lq design.ListQuery) ([]design.Design, design.Pagination, error)
}

type Minifier interface {
	HTML(s string) (string, error)
}

type Renderer interface {
	HTML(r io.Reader) ([]byte, error)
}

func main() {

	log.Println("initialising api...")

	//TODO Too much plumbing, see if there is a better way
	if os.Getenv("env") == "staging" {
		b, err := os.ReadFile(*renderDBPath)
		if err != nil {
			log.Fatal(err)
		}
		var secrets map[string]string
		err = json.Unmarshal(b, &secrets)
		if err != nil {
			log.Fatal(err)
		}
		v, ok := secrets["pg-conn-string"]
		if !ok {
			log.Fatal("unable to get pg-conn-string from environmet")
		}
		*connString = v
		*firebasePath = "/etc/secrets/firebase.json"
	}

	db := dbutils.MustOpenPostgres(*connString)
	profileRepo := account.NewProfileRepository(db)
	designRepo := design.NewDesignRepository(db)
	minify := minifier.NewMinifier()
	renderer := pdfrender.NewPDFRenderer()

	profileAPI := NewProfileAPI(profileRepo)
	designAPI := NewDesignAPI(designRepo, minify)
	generatorAPI := NewGeneratorAPI(designRepo, renderer)

	b, err := os.ReadFile(*firebasePath)
	if err != nil {
		log.Fatal(err)
	}
	fAuth := firebase.New(bytes.NewBuffer(b))

	log.Println("initialising routes...")
	r := chi.NewRouter()
	r.Use(m.RequestID)

	r.Group(func(r chi.Router) {
		r.Use(middleware.NewFirebaseAuth(fAuth, profileRepo).FirebaseAuth)
		r.Use(m.Logger)

		r.Route("/design", func(r chi.Router) {

			r.Post("/", designAPI.CreateDesign)
			r.Get("/", designAPI.ListDesign)

			r.Route("/{designId}", func(r chi.Router) {
				r.Get("/", designAPI.GetDesign)
				r.Put("/", designAPI.UpdateDesign)
				r.Delete("/", designAPI.DeleteDesign)
			})
		})

		r.Get("/profile", profileAPI.GetProfile)
		r.Post("/generate", generatorAPI.GeneratePDF)
		r.Post("/validate", designAPI.ValidateDesign)

	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.NewFirebaseAuth(fAuth, profileRepo).FirebaseAuthWithout)
		r.Post("/profile", profileAPI.CreateProfile)
	})
	r.Group(func(r chi.Router) {
		r.Get("/health", profileAPI.Health)
	})

	log.Println("starting api...")
	s := server.NewHTTPServer(*addr, r, *shutdownTimeout)

	//TODO need to find a way to put this back and remove migrate from docker-compose
	/*if os.Getenv("env") == "staging" {
		driver, err := postgres.WithInstance(db, &postgres.Config{})
		if err != nil {
			log.Fatal(err)
		}

		sourceUrl := fmt.Sprintf("file://%s", "migrations")
		m, err := migrate.NewWithDatabaseInstance(
			sourceUrl,
			"postgres", driver)
		if err != nil {
			log.Fatal(err)
		}
		err = m.Up()
		if err != nil && !errors.Is(err, migrate.ErrNoChange) {
			log.Fatal(err)
		}
	}*/

	s.Start()

	sig := service.Wait(syscall.SIGTERM, syscall.SIGINT)

	log.Printf("recieved signal %s", sig.String())

	s.Stop()
}
