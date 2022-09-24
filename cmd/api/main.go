package main

import (
	"bytes"
	"context"
	"flag"
	"github.com/go-chi/chi/v5"
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
	connString      = flag.String("pg-conn-string", "postgres://pdfgen:pdfgen@pg:5432/pdfgen?sslmode=disable", "PostgresSQL server connection string")
)

type ProfileRepository interface {
	Save(ctx context.Context, p account.Profile) error
}

type DesignRepository interface {
	Save(ctx context.Context, d design.Design) error
	GetByID(ctx context.Context, id string) (design.Design, error)
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

	r.Group(func(r chi.Router) {
		r.Use(middleware.NewFirebaseAuth(fAuth).FirebaseAuth)
		r.Post("/profile", profileAPI.CreateProfile)
		r.Post("/design", designAPI.CreateDesign)
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
