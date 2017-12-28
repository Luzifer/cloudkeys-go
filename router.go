package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/assets/").HandlerFunc(serveAssets)

	r.HandleFunc("/register", httpHelper(simpleTemplateOutput("register.html"))).
		Methods("GET")
	r.HandleFunc("/register", httpHelper(registerHandler)).
		Methods("POST")

	r.HandleFunc("/login", httpHelper(simpleTemplateOutput("login.html"))).
		Methods("GET")
	r.HandleFunc("/login", httpHelper(loginHandler)).
		Methods("POST")

	r.HandleFunc("/logout", httpHelper(logoutHandler)).
		Methods("GET")

	r.HandleFunc("/u/{userIndex:[0-9]+}/overview", httpHelper(overviewHandler)).
		Methods("GET")

	r.HandleFunc("/u/{userIndex:[0-9]+}/ajax", httpHelper(ajaxGetHandler)).
		Methods("GET")
	r.HandleFunc("/u/{userIndex:[0-9]+}/ajax", httpHelper(ajaxPostHandler)).
		Methods("POST")

	/* --- SUPPORT FOR DEPRECATED METHODS --- */
	r.HandleFunc("/", func(res http.ResponseWriter, r *http.Request) {
		http.Redirect(res, r, "u/0/overview", http.StatusFound)
	}).Methods("GET")
	r.HandleFunc("/overview", func(res http.ResponseWriter, r *http.Request) {
		http.Redirect(res, r, "u/0/overview", http.StatusFound)
	}).Methods("GET")

	return r
}
