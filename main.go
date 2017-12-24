package main // import "github.com/Luzifer/cloudkeys-go"

//go:generate go-bindata assets

import (
	"crypto/sha1"
	"fmt"
	"mime"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var (
	storage     storageAdapter
	cookieStore *sessions.CookieStore

	cfg     = loadConfig()
	version = "dev"
)

func init() {
	if cfg.VersionAndQuit {
		fmt.Printf("cloudkeys-go %s\n", version)
		os.Exit(0)
	}

	if _, err := cfg.ParsedStorage(); err != nil {
		log.WithError(err).Error("Unable to parse storage URI")
		os.Exit(1)
	}

	if cfg.CookieSigningKey == "" {
		cfg.CookieSigningKey = uuid.NewV4().String()[:32]
		log.Warn("cookie-authkey was set randomly, this will break your sessions!")
	}

	if cfg.CookieEncryptKey == "" {
		cfg.CookieEncryptKey = uuid.NewV4().String()[:32]
		log.Warn("cookie-encryptkey was set randomly, this will break your sessions!")
	}

	cookieStore = sessions.NewCookieStore(
		[]byte(cfg.CookieSigningKey),
		[]byte(cfg.CookieEncryptKey),
	)
}

func main() {
	s, err := getStorageAdapter(cfg)
	if err != nil {
		log.WithError(err).Fatal("Could not instanciate storage")
	}
	storage = s

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

	http.ListenAndServe(cfg.Listen, r)
}

func serveAssets(res http.ResponseWriter, r *http.Request) {
	data, err := Asset(r.RequestURI[1:])
	if err != nil {
		http.Error(res, "Not found", http.StatusNotFound)
		return
	}

	res.Header().Set("Content-Type", mime.TypeByExtension(filepath.Ext(r.RequestURI)))
	res.Write(data)
}

func createUserFilename(username string) string {
	return fmt.Sprintf("%x", sha1.Sum([]byte(cfg.UsernameSalt+username)))
}
