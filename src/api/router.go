package api

import (
	"net/http"

	"whoKnows/api/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func CreateRouter() http.Handler {

	router := mux.NewRouter()
	setApiRoutes(router)
	setFileRoutes(router)

	corsRouter := corsMiddleware()

	return corsRouter(router)
}

// Middleware to handle CORS
func corsMiddleware() func(http.Handler) http.Handler {

	return cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization", "Content-Type"},
		AllowCredentials: true,
	}).Handler
}

func setApiRoutes(mux *mux.Router) {

	mux.HandleFunc("/api/register", handlers.RegisterHandler)
	mux.HandleFunc("/api/login", handlers.LoginHandler)
	mux.HandleFunc("/api/logout", handlers.LogoutHandler)
	mux.HandleFunc("/api/search", handlers.SearchHandler)
}

func setFileRoutes(mux *mux.Router) {
	mux.HandleFunc("/", handlers.HomeHandler)
}
