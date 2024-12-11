package api

/* import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)*/

/* func loadRoutes() *chi.Mux {
	r := chi.NewRouter()
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Origin", "X-Requested-With", "Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	// TODO: remove default logger & add custom one.
	r.Use(middleware.DefaultLogger)

	r.Route("/", loadHealthCheckRoutes)

	return r
} */

// func loadHealthCheckRoutes(router chi.Router) {
// TODO: create handlers for  healthcheck and update page
// orderHandler := &

// router.Get("/healthcheck", healthcheck.Health)
// }
