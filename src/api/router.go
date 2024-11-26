package api

import (
	"mime"
	"net/http"
	"path/filepath"

	"whoKnows/api/handlers"
	"whoKnows/monitoring"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func CreateRouter() http.Handler {

	router := mux.NewRouter()
	setApiRoutes(router)
	setFileRoutes(router)

	monitoring.RegisterMetrics()
	monitoring.ExposeMetrics(router)

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
	mux.HandleFunc("/", handlers.ServeHome)
	mux.HandleFunc("/about", handlers.ServeAbout)
	mux.HandleFunc("/search", handlers.ServeSearch)
	mux.HandleFunc("/register", handlers.ServeRegister)
	mux.HandleFunc("/login", handlers.ServeLogin)
	mux.PathPrefix("/static/").Handler(http.StripPrefix("/static/", customFileServer(http.Dir("./static/"))))
}

func customFileServer(root http.FileSystem) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set MIME type based on file extension
		if ext := filepath.Ext(r.URL.Path); ext != "" {
			if mimeType := mime.TypeByExtension(ext); mimeType != "" {
				w.Header().Set("Content-Type", mimeType)
			}
		}
		http.FileServer(root).ServeHTTP(w, r)
	})
}
