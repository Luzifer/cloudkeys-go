package main

import (
	"context"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

func overviewHandler(c context.Context, res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
	user, _ := checkLogin(c, r, session)

	if user == nil || !storage.IsPresent(c, user.UserFile) {
		http.Redirect(res, r, "../../login", http.StatusFound)
		return nil, nil
	}

	frontendAccounts := []string{}
	idx := -1
	for i, v := range session.Values["authorizedAccounts"].(authorizedAccounts) {
		frontendAccounts = append(frontendAccounts, v.Name)
		if v.Name == user.Name {
			idx = i
		}
	}

	(*ctx)["authorized_accounts"] = frontendAccounts
	(*ctx)["current_user_index"] = idx

	return stringPointer("overview.html"), nil
}
