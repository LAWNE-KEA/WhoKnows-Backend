package api

import (
	"fmt"
	"net/http"

	"whoKnows/api/handlers"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

// Router returns a mux router with all routes defined and cors enabled for all routes

// var tmpl = template.Must(template.ParseFiles(
// 	"../app/frontend/root.html",
// 	"../app/frontend/search.html",
// 	"../app/frontend/register.html",
// 	"../app/frontend/login.html",
// 	"../app/frontend/about.html",
// ))

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
	// Serve static files (CSS, JS, images, etc.)
	fs := http.FileServer(http.Dir("./static"))
	mux.Handle("/static/", http.StripPrefix("/static/", fs))

	// Serve HTML files
	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./app/frontend/weather.html")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("Root handler called")
		http.ServeFile(w, r, "./app/frontend/root.html")
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./app/frontend/search.html")
	})
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./app/frontend/about.html")
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./app/frontend/register.html")
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./app/frontend/login.html")
	})
}

// func aboutHandler(w http.ResponseWriter, r *http.Request) {
// 	data := types.ResonseObject{
// 		User:    &models.User{Username: "JohnDoe"}, // Example user, replace with actual user data
// 		Flashes: []string{"Welcome to the About page!"},
// 	}
// 	tmpl.ExecuteTemplate(w, "about.html", data)
// }

// func loginHandler(w http.ResponseWriter, r *http.Request) {
// 	data := types.ResonseObject{
// 		User:    &models.User{Username: "JohnDoe"}, // Example user, replace with actual user data
// 		Flashes: []string{"Please log in."},
// 	}
// 	tmpl.ExecuteTemplate(w, "login.html", data)
// }

// func registerHandler(w http.ResponseWriter, r *http.Request) {
// 	data := types.ResonseObject{
// 		User:    &models.User{Username: "JohnDoe"}, // Example user, replace with actual user data
// 		Flashes: []string{"Please register."},
// 	}
// 	tmpl.ExecuteTemplate(w, "register.html", data)
// }

// func logoutHandler(w http.ResponseWriter, r *http.Request) {

// 	http.SetCookie(w, &http.Cookie{
// 		Name:     "session_id",
// 		Value:    "",
// 		Expires:  time.Now().Add(-1 * time.Hour),
// 		Path:     "/",
// 		HttpOnly: true,
// 		Secure:   false,
// 		SameSite: http.SameSiteLaxMode,
// 	})

// 	http.Redirect(w, r, "/", http.StatusFound)
// }
