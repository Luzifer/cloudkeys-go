package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"
	"strconv"

	"github.com/flosch/pongo2"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
)

func loginHandler(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	var (
		username = r.FormValue("username")
		password = fmt.Sprintf("%x", sha1.Sum([]byte(cfg.PasswordSalt+r.FormValue("password"))))
	)

	if !storage.IsPresent(createUserFilename(username)) {
		(*ctx)["error"] = true
		return stringPointer("login.html"), nil
	}

	userFileRaw, err := storage.Read(createUserFilename(username))
	if err != nil {
		fmt.Printf("ERR: Unable to read user file: %s\n", err)
		(*ctx)["error"] = true
		return stringPointer("login.html"), nil
	}

	userFile, _ := readDataObject(userFileRaw)

	if userFile.MetaData.Password != password {
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

func logoutHandler(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	session.Values["authorizedAccounts"] = authorizedAccounts{}
	session.Save(r, res)
	http.Redirect(res, r, "overview", http.StatusFound)
	return nil, nil
}

func checkLogin(r *http.Request, session *sessions.Session) (*authorizedAccount, error) {
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
