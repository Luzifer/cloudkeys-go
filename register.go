package main

import (
	"context"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

func registerHandler(c context.Context, res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	var (
		username      = strings.ToLower(r.FormValue("username"))
		password      = r.FormValue("password")
		passwordCheck = r.FormValue("password_repeat")
	)

	if username == "" || password == "" || password != passwordCheck {
		return stringPointer("register.html"), nil
	}

	if storage.IsPresent(c, createUserFilename(username)) {
		(*ctx)["exists"] = true
		return stringPointer("register.html"), nil
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.WithError(err).Error("Could not hash user password")
		(*ctx)["error"] = true
		return stringPointer("register.html"), nil
	}

	d := dataObject{}
	d.MetaData.Password = string(hashedPassword)
	data, _ := d.GetData()

	if err := storage.Write(c, createUserFilename(username), data); err != nil {
		log.WithError(err).Error("Could not write user file to storage")
		(*ctx)["error"] = true
		return stringPointer("register.html"), nil
	}

	(*ctx)["created"] = true
	return stringPointer("register.html"), nil
}
