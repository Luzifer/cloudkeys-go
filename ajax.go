package main

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
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

func ajaxGetHandler(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	res.Header().Set("Content-Type", "application/json")
	user, err := checkLogin(r, session)
	if err != nil {
		return nil, err // TODO: Handle in-app?
	}

	if user == nil || !storage.IsPresent(user.UserFile) {
		res.Write(ajaxResponse{Error: true}.Bytes())
		return nil, nil
	}

	userFileRaw, err := storage.Read(user.UserFile)
	if err != nil {
		return nil, err // TODO: Handle in-app?
	}

	userFile, err := readDataObject(userFileRaw)
	if err != nil {
		return nil, err // TODO: Handle in-app?
	}

	res.Write(ajaxResponse{Version: userFile.MetaData.Version, Data: userFile.Data}.Bytes())
	return nil, nil
}

func ajaxPostHandler(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	res.Header().Set("Content-Type", "application/json")
	user, err := checkLogin(r, session)
	if err != nil {
		return nil, err // TODO: Handle in-app?
	}

	if user == nil {
		res.Write(ajaxResponse{Error: true, Type: "login"}.Bytes())
		return nil, nil
	}

	if !storage.IsPresent(user.UserFile) {
		res.Write(ajaxResponse{Error: true, Type: "register"}.Bytes())
		return nil, nil
	}

	userFileRaw, err := storage.Read(user.UserFile)
	if err != nil {
		return nil, err // TODO: Handle in-app?
	}

	userFile, err := readDataObject(userFileRaw)
	if err != nil {
		return nil, err // TODO: Handle in-app?
	}

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

	if err := storage.Backup(user.UserFile); err != nil {
		return nil, err // TODO: Handle in-app?
	}

	userFile.MetaData.Version = checksum
	userFile.Data = data

	d, err := userFile.GetData()
	if err != nil {
		return nil, err // TODO: Handle in-app?
	}

	if err := storage.Write(user.UserFile, d); err != nil {
		return nil, err // TODO: Handle in-app?
	}

	res.Write(ajaxResponse{Version: userFile.MetaData.Version, Data: userFile.Data}.Bytes())
	return nil, nil
}
