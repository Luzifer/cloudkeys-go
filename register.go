package main

import (
	"crypto/sha1"
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

func registerHandler(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	var (
		username       = r.FormValue("username")
		password       = r.FormValue("password")
		passwordCheck  = r.FormValue("password_repeat")
		hashedPassword = fmt.Sprintf("%x", sha1.Sum([]byte(cfg.PasswordSalt+password)))
	)

	if username == "" || password == "" || password != passwordCheck {
		return stringPointer("register.html"), nil
	}

	if storage.IsPresent(createUserFilename(username)) {
		(*ctx)["exists"] = true
		return stringPointer("register.html"), nil
	}

	d := dataObject{}
	d.MetaData.Password = hashedPassword
	data, _ := d.GetData()

	if err := storage.Write(createUserFilename(username), data); err != nil {
		fmt.Printf("ERR: Unable to write user file: %s\n", err)
		(*ctx)["error"] = true
		return stringPointer("register.html"), nil
	}

	(*ctx)["created"] = true
	return stringPointer("register.html"), nil
}
