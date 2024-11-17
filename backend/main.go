package main

import (
	"fmt"
	"html/template"
	"mime"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	_ "github.com/go-sql-driver/mysql"

	"whoKnows/models"
)

var tmpl = template.Must(template.ParseFiles(
	"../app/frontend/root.html",
	"../app/frontend/search.html",
	"../app/frontend/register.html",
	"../app/frontend/login.html",
	"../app/frontend/about.html",
))

var ENV_MYSQL_USER, _ = os.LookupEnv("ENV_MYSQL_USER")
var ENV_MYSQL_PASSWORD, _ = os.LookupEnv("ENV_MYSQL_PASSWORD")
var ENV_INIT_MODE, _ = os.LookupEnv("ENV_INIT_MODE")
var DATABASE_PATH = ENV_MYSQL_USER + ":" + ENV_MYSQL_PASSWORD + "@(mysql_db:3306)/whoknows"

//var sessionStore = map[string]session{}

// Run the server on port 8080
func main() {
	fmt.Println("Starting server on port 8080")

	fs := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", http.StripPrefix("/static/", fs))

	initDB(ENV_INIT_MODE == "true")

	mux := http.NewServeMux()
	mux.HandleFunc("/weather", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/weather.html")
	})
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/root.html")
	})
	mux.HandleFunc("/search", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/search.html")
	})

	//mux.HandleFunc("/about", aboutHandler)
	//mux.HandleFunc("/login", loginHandler)
	//mux.HandleFunc("/register", registerHandler)
	//mux.HandleFunc("/logout", logoutHandler)
	mux.HandleFunc("/api/search", apiSearchHandler)
	mux.HandleFunc("/api/login", apiLoginHandler)
	mux.HandleFunc("/api/register", apiRegisterHandler)

	mux.Handle("/static/", http.StripPrefix("/static/", fs))
	mux.HandleFunc("/about", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/about.html")
	})
	mux.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/register.html")
	})
	mux.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "/app/frontend/login.html")
	})

	// Apply CORS middleware
	handler := corsMiddleware(mux)

	// Create a non-global registry.
	reg := prometheus.NewRegistry()

	// Create new metrics and register them using the custom registry.
	m := NewMetrics(reg)
	// Set values for the new created metrics.
	m.cpuTemp.Set(65.3)
	m.hdFailures.With(prometheus.Labels{"device": "/dev/sda"}).Inc()

	// Expose metrics and custom registry via an HTTP server
	// using the HandleFor function. "/metrics" is the usual endpoint for that.
	mux.Handle("/metrics", promhttp.HandlerFor(reg, promhttp.HandlerOpts{Registry: reg}))
	http.ListenAndServe(":8080", handler)
}

func init() {
	mime.AddExtensionType(".css", "text/css")
}

func parseSQLCommands(sqlCommands string) []string {
	var commands []string
	var currentCommand strings.Builder
	inSingleQuote := false
	inDoubleQuote := false

	for _, char := range sqlCommands {
		switch char {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case ';':
			if !inSingleQuote && !inDoubleQuote {
				commands = append(commands, currentCommand.String())
				currentCommand.Reset()
				continue
			}
		}
		currentCommand.WriteRune(char)
	}

	// Add the last command if any
	if currentCommand.Len() > 0 {
		commands = append(commands, currentCommand.String())
	}

	return commands
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		User:    &User{Username: "JohnDoe"}, // Example user, replace with actual user data
		Flashes: []string{"Welcome to the About page!"},
	}
	tmpl.ExecuteTemplate(w, "about.html", data)
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		User:    &User{Username: "JohnDoe"}, // Example user, replace with actual user data
		Flashes: []string{"Please log in."},
	}
	tmpl.ExecuteTemplate(w, "login.html", data)
}

func registerHandler(w http.ResponseWriter, r *http.Request) {
	data := models.PageData{
		User:    &User{Username: "JohnDoe"}, // Example user, replace with actual user data
		Flashes: []string{"Please register."},
	}
	tmpl.ExecuteTemplate(w, "register.html", data)
}

func logoutHandler(w http.ResponseWriter, r *http.Request) {
	sessionID, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "No active session", http.StatusBadRequest)
		return
	}

	delete(sessionStore, sessionID.Value)

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		Path:     "/",
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/", http.StatusFound)
}

type metrics struct {
	cpuTemp    prometheus.Gauge
	hdFailures *prometheus.CounterVec
}

func NewMetrics(reg prometheus.Registerer) *metrics {
	m := &metrics{
		cpuTemp: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: "cpu_temperature_celsius",
			Help: "Current temperature of the CPU.",
		}),
		hdFailures: prometheus.NewCounterVec(
			prometheus.CounterOpts{
				Name: "hd_errors_total",
				Help: "Number of hard-disk errors.",
			},
			[]string{"device"},
		),
	}
	reg.MustRegister(m.cpuTemp)
	reg.MustRegister(m.hdFailures)
	return m
}
