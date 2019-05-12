package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"time"

	"github.com/bayugyug/gorm-custom-api/api/handler"
	"github.com/bayugyug/gorm-custom-api/configs"
	"github.com/bayugyug/gorm-custom-api/drivers"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// APIRouter the svc map
type APIRouter struct {
	Building *handler.Building
	Mux      *chi.Mux
	Address  string
	Config   *configs.ParameterConfig
}

// Setup options settings
type Setup func(*APIRouter)

// WithSvcOptConfig opts for mux
func WithSvcOptConfig(c *configs.ParameterConfig) Setup {
	return func(args *APIRouter) {
		args.Config = c
	}
}

// WithSvcOptMux opts for mux
func WithSvcOptMux(m *chi.Mux) Setup {
	return func(args *APIRouter) {
		args.Mux = m
	}
}

// WithSvcOptHandler opts for handler
func WithSvcOptHandler(r *handler.Building) Setup {
	return func(args *APIRouter) {
		args.Building = r
	}
}

// NewAPIRouter service new instance
func NewAPIRouter(opts ...Setup) (*APIRouter, error) {

	//default
	svc := &APIRouter{
		Address: ":8989",
	}

	//add options if any
	for _, setter := range opts {
		setter(svc)
	}
	//port
	if svc.Config != nil && svc.Config.Port != "" {
		svc.Address = ":" + svc.Config.Port
	}

	dbh := drivers.NewDBHandle(
		"building-api-db-driver",
		"mysql",
		svc.Config.DSN)

	db := dbh.GetConnection()
	if db == nil {
		return svc, fmt.Errorf("DB Connect failed")
	}
	//db storage
	svc.Building = handler.NewBuilding(db)
	//set the actual router
	svc.Mux = svc.MapRoute()

	//good :-)
	return svc, nil
}

// Run the http server based on settings
func (svc *APIRouter) Run() {

	//gracious timing
	srv := &http.Server{
		Addr:         svc.Address,
		Handler:      svc.Mux,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	//async run
	go func() {
		log.Println("Listening on port", svc.Address)
		if err := srv.ListenAndServe(); err != nil {
			log.Printf("listen: %s\n", err)
			os.Exit(0)
		}

	}()

	//watcher
	stopChan := make(chan os.Signal)
	signal.Notify(stopChan, os.Interrupt)

	<-stopChan
	log.Println("Shutting down service...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	srv.Shutdown(ctx)
	defer cancel()
	log.Println("Server gracefully stopped!")
}

// MapRoute route map all endpoints
func (svc *APIRouter) MapRoute() *chi.Mux {

	// Multiplexer
	router := chi.NewRouter()

	// Basic settings
	router.Use(
		render.SetContentType(render.ContentTypeJSON),
		middleware.Logger,
		middleware.DefaultCompress,
		middleware.StripSlashes,
		middleware.Recoverer,
		middleware.RequestID,
		middleware.RealIP,
	)

	// Basic gracious timing
	router.Use(middleware.Timeout(60 * time.Second))

	// Basic CORS
	cors := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		AllowCredentials: true,
	})

	router.Use(cors.Handler)

	router.Get("/", svc.Building.Welcome)
	router.Get("/health", svc.Building.HealthCheck)

	/*
		@end-points

		GET    /v1/api/building/:id
		POST   /v1/api/building
		PUT    /v1/api/building
		DELETE /v1/api/building/:id

	*/

	//end-points-mapping
	router.Route("/v1", func(r chi.Router) {
		r.Mount("/api",
			func(h *handler.Building) *chi.Mux {
				sr := chi.NewRouter()
				sr.Get("/health", h.HealthCheck)
				sr.Post("/building", h.Create)
				sr.Put("/building", h.Update)
				sr.Patch("/building", h.Update)
				sr.Get("/building", h.GetAll)
				sr.Get("/building/{id}", h.GetOne)
				sr.Delete("/building/{id}", h.Delete)
				return sr
			}(svc.Building))
	})
	//show
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		route = strings.Replace(route, "/*/", "/", -1)
		fmt.Printf("... %s %s\n", method, route)
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		fmt.Printf("Logging err: %s\n", err.Error())
	}
	return router
}
