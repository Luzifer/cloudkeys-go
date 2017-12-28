package main

import (
	"context"
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
)

type ajaxResponse struct {
	Error   bool   `json:"error"`
	Version string `json:"version"`
	Data    string `json:"data"`
	Type    string `json:"type"`
}

func (a ajaxResponse) Bytes() []byte {
	out, _ := json.Marshal(a)
	return out
}

func ajaxGetHandler(c context.Context, res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	res.Header().Set("Content-Type", "application/json")
	user, _ := checkLogin(c, r, session)

	if user == nil || !storage.IsPresent(c, user.UserFile) {
		res.Write(ajaxResponse{Error: true}.Bytes())
		return nil, nil
	}

	userFileRaw, err := storage.Read(c, user.UserFile)
	if err != nil {
		log.WithError(err).Error("Could not read user file from storage")
		res.Write(ajaxResponse{Error: true}.Bytes())
		return nil, nil
	}

	userFile, _ := readDataObject(userFileRaw)

	res.Write(ajaxResponse{Version: userFile.MetaData.Version, Data: userFile.Data}.Bytes())
	return nil, nil
}

func ajaxPostHandler(c context.Context, res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	res.Header().Set("Content-Type", "application/json")
	user, _ := checkLogin(c, r, session)

	if user == nil {
		res.Write(ajaxResponse{Error: true, Type: "login"}.Bytes())
		return nil, nil
	}

	if !storage.IsPresent(c, user.UserFile) {
		res.Write(ajaxResponse{Error: true, Type: "register"}.Bytes())
		return nil, nil
	}

	userFileRaw, err := storage.Read(c, user.UserFile)
	if err != nil {
		log.WithError(err).Error("Could not read user file from storage")
		res.Write(ajaxResponse{Error: true, Type: "storage_error"}.Bytes())
		return nil, nil
	}

	userFile, _ := readDataObject(userFileRaw)

	var (
		version  = r.FormValue("version")
		checksum = r.FormValue("checksum")
		data     = r.FormValue("data")
	)

	if userFile.MetaData.Version != version {
		res.Write(ajaxResponse{Error: true, Type: "version"}.Bytes())
		return nil, nil
	}

	if checksum != fmt.Sprintf("%x", sha1.Sum([]byte(data))) {
		res.Write(ajaxResponse{Error: true, Type: "checksum"}.Bytes())
		return nil, nil
	}

	if err := storage.Backup(c, user.UserFile); err != nil {
		log.WithError(err).Error("Could not create backup of user file")
		res.Write(ajaxResponse{Error: true, Type: "storage_error"}.Bytes())
		return nil, nil
	}

	userFile.MetaData.Version = checksum
	userFile.Data = data

	d, _ := userFile.GetData()

	if err := storage.Write(c, user.UserFile, d); err != nil {
		log.WithError(err).Error("Could not write user file to storage")
		res.Write(ajaxResponse{Error: true, Type: "storage_error"}.Bytes())
		return nil, nil
	}

	res.Write(ajaxResponse{Version: userFile.MetaData.Version, Data: userFile.Data}.Bytes())
	return nil, nil
}
