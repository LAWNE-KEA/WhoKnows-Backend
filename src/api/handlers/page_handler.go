package handlers

import (
	"html/template"
	"net/http"
)

func ServeHome(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/root.html"))
	tmpl.Execute(w, nil)
}

func ServeAbout(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/about.html"))
	tmpl.Execute(w, nil)
}

func ServeSearch(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/search.html"))
	tmpl.Execute(w, nil)
}

func ServeRegister(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/register.html"))
	tmpl.Execute(w, nil)
}

func ServeLogin(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("pages/login.html"))
	tmpl.Execute(w, nil)
}
