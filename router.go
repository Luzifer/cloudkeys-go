package main

import (
	"net/http"

	"github.com/gorilla/mux"
)

func router() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/").HandlerFunc(gzipFunc(serveAssets))

	registerAPIv2(r.PathPrefix("/v2").Subrouter())

	return r
}

func registerAPIv2(r *mux.Router) {
	r.HandleFunc("/login", apiHelper(apiLogin)).Methods(http.MethodPost)
	r.HandleFunc("/register", apiHelper(apiRegister)).Methods(http.MethodPost)
	r.HandleFunc("/users", apiHelper(apiListUsers)).Methods(http.MethodGet)

	r.HandleFunc("/user/{user}/data", apiHelper(apiGetUserData)).Methods(http.MethodGet)
	r.HandleFunc("/user/{user}/data", apiHelper(apiSetUserData)).Methods(http.MethodPut)
	r.HandleFunc("/user/{user}/logout", apiHelper(apiLogoutUser)).Methods(http.MethodPost)
	r.HandleFunc("/user/{user}/settings", apiHelper(apiGetUserSettings)).Methods(http.MethodGet)
	r.HandleFunc("/user/{user}/settings", apiHelper(apiSetUserSettings)).Methods(http.MethodPatch)
	r.HandleFunc("/user/{user}/password", apiHelper(apiChangeLoginPassword)).Methods(http.MethodPut)
	r.HandleFunc("/user/{user}/validate-mfa", apiHelper(apiValidateMFA)).Methods(http.MethodPost)
}
