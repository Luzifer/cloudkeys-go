package main

import (
	"context"
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

func loginHandler(c context.Context, res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	var (
		username           = strings.ToLower(r.FormValue("username"))
		password           = r.FormValue("password")
		deprecatedPassword = fmt.Sprintf("%x", sha1.Sum([]byte(cfg.PasswordSalt+r.FormValue("password")))) // Here for backwards compatibility
	)

	if !storage.IsPresent(c, createUserFilename(username)) {
		(*ctx)["error"] = true
		return stringPointer("login.html"), nil
	}

	userFileRaw, err := storage.Read(c, createUserFilename(username))
	if err != nil {
		log.WithError(err).Error("Unable to read user file")
		(*ctx)["error"] = true
		return stringPointer("login.html"), nil
	}

	userFile, _ := readDataObject(userFileRaw)

	bcryptValidationError := bcrypt.CompareHashAndPassword([]byte(userFile.MetaData.Password), []byte(password))
	if bcryptValidationError != nil && userFile.MetaData.Password != deprecatedPassword {
		(*ctx)["error"] = true
		return stringPointer("login.html"), nil
	}

	auth, ok := session.Values["authorizedAccounts"].(authorizedAccounts)
	if !ok {
		auth = authorizedAccounts{}
	}

	for i, v := range auth {
		if v.Name == username {
			http.Redirect(res, r, fmt.Sprintf("u/%d/overview", i), http.StatusFound)
			return nil, nil
		}
	}

	auth = append(auth, authorizedAccount{
		Name:     username,
		UserFile: createUserFilename(username),
	})

	session.Values["authorizedAccounts"] = auth
	if err := session.Save(r, res); err != nil {
		return nil, err
	}

	http.Redirect(res, r, fmt.Sprintf("u/%d/overview", len(auth)-1), http.StatusFound)
	return nil, nil
}

func logoutHandler(c context.Context, res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	session.Values["authorizedAccounts"] = authorizedAccounts{}
	session.Save(r, res)
	http.Redirect(res, r, "overview", http.StatusFound)
	return nil, nil
}

func checkLogin(c context.Context, r *http.Request, session *sessions.Session) (*authorizedAccount, error) {
	vars := mux.Vars(r)
	idx, err := strconv.ParseInt(vars["userIndex"], 10, 64)
	if err != nil {
		return nil, err
	}

	auth, ok := session.Values["authorizedAccounts"].(authorizedAccounts)
	if !ok {
		auth = authorizedAccounts{}
	}

	if len(auth)-1 < int(idx) {
		return nil, nil
	}

	return &auth[idx], nil
}
