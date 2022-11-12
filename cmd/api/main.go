package main

import (
	"context"
	"encoding/json"
	"flag"
	"github.com/go-chi/chi/v5"
	m "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/rengas/pdfgen/pkg/dbutils"
	"github.com/rengas/pdfgen/pkg/design"
	"github.com/rengas/pdfgen/pkg/logging"
	cmiddleware "github.com/rengas/pdfgen/pkg/middleware"
	"github.com/rengas/pdfgen/pkg/minifier"
	"github.com/rengas/pdfgen/pkg/pagination"
	"github.com/rengas/pdfgen/pkg/password"
	"github.com/rengas/pdfgen/pkg/pdfrender"
	"github.com/rengas/pdfgen/pkg/server"
	"github.com/rengas/pdfgen/pkg/service"
	"github.com/rengas/pdfgen/pkg/token"
	"github.com/rengas/pdfgen/pkg/user"
	"io"
	"log"
	"os"
	"syscall"
	"time"
)

var (
	addr                  = flag.String("addr", ":8080", "Application http server network address")
	shutdownTimeout       = flag.Duration("shutdown-timeout", 30*time.Second, "Graceful shutdown timeout")
	connString            = flag.String("pg-conn-string", "postgres://pdfgen:pdfgen@pg:5432/pdfgen?sslmode=disable", "PostgresSQL server connection string")
	renderDBPath          = flag.String("render-db-path", "/etc/secrets/staging.json", "staging db creds")
	passwordPepper        = flag.String("password-pepper", "secret-random-string", "some random secret string")
	jwtAccessSecretKey    = flag.String("jwt-access-key", "secret-random-access-key", "some random secret key")
	jwtRefreshSecretKey   = flag.String("jwt-secret-key", "secret-refresh-access-key", "some random refresh secret key")
	jwtAccessTokenExpiry  = flag.Int("jwt-access-expiry", 30, "some access token expiry in minutes")
	jwtRefreshTokenExpiry = flag.Int("jwt-refresh-expiry", 2, "some refresh token expiry in hours")
)

type UserRepository interface {
	SaveNewUser(ctx context.Context, p user.User) error
	GetByEmail(ctx context.Context, email string) (user.User, error)
	GetById(ctx context.Context, id string) (user.User, error)
	Update(ctx context.Context, u user.User) error
}

type Bcrypt interface {
	GetHashedPassword(password string) ([]byte, error)
	CompareHashedPassword(found, given string) error
}

type JWTToken interface {
	TokePair(claims map[string]interface{}) (token.TokenDetails, error)
}

type DesignRepository interface {
	Save(ctx context.Context, d design.Design) error
	Update(ctx context.Context, p design.Design) error
	GetById(ctx context.Context, userId, designId string) (design.Design, error)
	Delete(ctx context.Context, userId string, designId string) error
	ListByUserId(ctx context.Context, lq design.ListQuery) ([]design.Design, pagination.Pagination, error)
	Search(ctx context.Context, lq design.ListQuery) ([]design.Design, pagination.Pagination, error)
}

type Minifier interface {
	HTML(s string) (string, error)
}

type Renderer interface {
	HTML(r io.Reader) ([]byte, error)
}

// @title                       Pdfgen.pro API
// @version                     1.0
// @description                 API
// @contact.name                dev
// @contact.email               dev@pdfgen.pro
// @host						https://pdfgen-stg.onrender.com
// @securityDefinitions.apikey  BearerAuth
// @schemes                     http https
// @in                          header
// @name                        Authorization
// @BasePath                    /
func main() {

	log.Println("initialising api...")
	logging.InitDefaultLogger(context.Background())

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
			logging.Fatal("unable to get pg-conn-string from environment")
		}
		*connString = v
	}

	db := dbutils.MustOpenPostgres(*connString)
	designRepo := design.NewDesignRepository(db)
	minify := minifier.NewMinifier()
	renderer := pdfrender.NewPDFRenderer()
	userRepo := user.NewRepository(db)

	designAPI := NewDesignAPI(designRepo, minify)
	generatorAPI := NewGeneratorAPI(designRepo, renderer)

	bcrypt := password.NewBcrypt(*passwordPepper)
	jwt := token.NewJWT(*jwtAccessSecretKey, *jwtRefreshSecretKey, *jwtAccessTokenExpiry, *jwtRefreshTokenExpiry)
	tokenMiddleware := cmiddleware.NewJWTToken(jwt)

	r := chi.NewRouter()
	r.Use(m.RequestID)

	// for more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	logging.Info("initialising routes...")

	authAPI := NewAuthAPI(userRepo, bcrypt, jwt)
	r.Route("/", func(r chi.Router) {
		r.Post("/register", authAPI.Register)
		r.Post("/login", authAPI.Login)
	})

	userAPI := NewUserAPI(userRepo)

	r.Route("/user", func(r chi.Router) {
		r.Use(tokenMiddleware.VerifyToken)
		r.Get("/", userAPI.GetUser)
		r.Put("/", userAPI.UpdateUser)
	})

	r.Group(func(r chi.Router) {
		r.Use(m.Logger)
		r.Use(tokenMiddleware.VerifyToken)
		r.Route("/design", func(r chi.Router) {

			r.Post("/", designAPI.CreateDesign)
			r.Get("/", designAPI.ListDesign)

			r.Route("/{designId}", func(r chi.Router) {
				r.Get("/", designAPI.GetDesign)
				r.Put("/", designAPI.UpdateDesign)
				r.Delete("/", designAPI.DeleteDesign)
			})
		})

		r.Post("/generate", generatorAPI.GeneratePDF)
		r.Post("/validate", designAPI.ValidateDesign)

	})

	r.Group(func(r chi.Router) {
		r.Get("/health", authAPI.Health)
	})

	log.Println("starting api...")
	s := server.NewHTTPServer(*addr, r, *shutdownTimeout)

	s.Start()

	sig := service.Wait(syscall.SIGTERM, syscall.SIGINT)

	logging.WithField(logging.Field{Label: "received signal", Value: sig.String()})

	s.Stop()
}
