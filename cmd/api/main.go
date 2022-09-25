package main

import (
	"bytes"
	"context"
	"flag"
	"github.com/go-chi/chi/v5"
	m "github.com/go-chi/chi/v5/middleware"
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
	db := dbutils.MustOpenPostgres(*connString)
	profileRepo := account.NewProfileRepository(db)
	designRepo := design.NewDesignRepository(db)
	minify := minifier.NewMinifier()
	renderer := pdfrender.NewPDFRenderer()

	profileAPI := NewProfileAPI(profileRepo)
	designAPI := NewDesignAPI(designRepo, minify)
	generatorAPI := NewGeneratorAPI(designRepo, renderer)

	b, err := os.ReadFile("./firebase.json")
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

	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.NewFirebaseAuth(fAuth, profileRepo).FirebaseAuth)
		r.Post("/profile", profileAPI.CreateProfile)
		r.Get("/profile", profileAPI.GetProfile)
		r.Post("/generate", generatorAPI.GeneratePDF)

	})

	r.Group(func(r chi.Router) {
		r.Get("/health", profileAPI.Health)
	})

	log.Println("starting api...")
	s := server.NewHTTPServer(*addr, r, *shutdownTimeout)
	s.Start()

	sig := service.Wait(syscall.SIGTERM, syscall.SIGINT)

	log.Printf("recieved signal %s", sig.String())

	s.Stop()
}
