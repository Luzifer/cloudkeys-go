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
	user, _ := checkLogin(r, session)

	if user == nil || !storage.IsPresent(user.UserFile) {
		res.Write(ajaxResponse{Error: true}.Bytes())
		return nil, nil
	}

	userFileRaw, err := storage.Read(user.UserFile)
	if err != nil {
		fmt.Printf("ERR: Unable to read user file: %s\n", err)
		res.Write(ajaxResponse{Error: true}.Bytes())
		return nil, nil
	}

	userFile, _ := readDataObject(userFileRaw)

	res.Write(ajaxResponse{Version: userFile.MetaData.Version, Data: userFile.Data}.Bytes())
	return nil, nil
}

func ajaxPostHandler(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	res.Header().Set("Content-Type", "application/json")
	user, _ := checkLogin(r, session)

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
		fmt.Printf("ERR: Unable to read user file: %s\n", err)
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

	if err := storage.Backup(user.UserFile); err != nil {
		fmt.Printf("ERR: Unable to backup user file: %s\n", err)
		res.Write(ajaxResponse{Error: true, Type: "storage_error"}.Bytes())
		return nil, nil
	}

	userFile.MetaData.Version = checksum
	userFile.Data = data

	d, _ := userFile.GetData()

	if err := storage.Write(user.UserFile, d); err != nil {
		fmt.Printf("ERR: Unable to write user file: %s\n", err)
		res.Write(ajaxResponse{Error: true, Type: "storage_error"}.Bytes())
		return nil, nil
	}

	res.Write(ajaxResponse{Version: userFile.MetaData.Version, Data: userFile.Data}.Bytes())
	return nil, nil
}
