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
	"github.com/bayugyug/gorm-custom-api/services"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/cors"
	"github.com/go-chi/render"
)

// APIRouter the svc map
type APIRouter struct {
	Mux      *chi.Mux
	Address  string
	Config   *configs.ParameterConfig
	DBHandle *drivers.DBHandle
	Handlers *Handler
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
func WithSvcOptHandler(r *Handler) Setup {
	return func(args *APIRouter) {
		args.Handlers = r
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

	svc.DBHandle = drivers.NewDBHandle(
		"building-api-db-driver",
		"mysql",
		svc.Config.DSN)

	db := svc.DBHandle.GetConnection()
	if db == nil {
		return svc, fmt.Errorf("DB Connect failed")
	}
	//handlers
	svc.Handlers = svc.GetHandlers()
	//set the actual router
	svc.Mux = svc.MapRoute()
	//good :-)
	return svc, nil
}

// GetHandlers is the mapping of all the handlers
func (svc *APIRouter) GetHandlers() *Handler {
	//attached all handlers per service
	handlers := &Handler{
		BuildingHandler: handler.NewBuilding(
			svc.DBHandle,
			services.NewBuildingService(),
		),
	}
	return handlers
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

	router.Get("/", svc.Handlers.BuildingHandler.Welcome)
	router.Get("/health", svc.Handlers.BuildingHandler.HealthCheck)

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
			func(h *Handler) *chi.Mux {
				sr := chi.NewRouter()
				sr.Get("/health", h.BuildingHandler.HealthCheck)
				sr.Post("/building", h.BuildingHandler.Create)
				sr.Put("/building", h.BuildingHandler.Update)
				sr.Patch("/building", h.BuildingHandler.Update)
				sr.Get("/building", h.BuildingHandler.GetAll)
				sr.Get("/building/{id}", h.BuildingHandler.GetOne)
				sr.Delete("/building/{id}", h.BuildingHandler.Delete)
				return sr
			}(svc.Handlers))
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

// Handler wraps all handler
type Handler struct {
	BuildingHandler buildingHandler
}

// list of all the handlers mapping
type buildingHandler interface {
	Create(w http.ResponseWriter, r *http.Request)
	Update(w http.ResponseWriter, r *http.Request)
	GetAll(w http.ResponseWriter, r *http.Request)
	GetOne(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	HealthCheck(w http.ResponseWriter, r *http.Request)
	Welcome(w http.ResponseWriter, r *http.Request)
}
