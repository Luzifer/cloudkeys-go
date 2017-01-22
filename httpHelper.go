package main

import (
	"fmt"
	"net/http"

	"github.com/flosch/pongo2"
	"github.com/gorilla/sessions"
)

type httpHelperFunc func(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error)

func httpHelper(f httpHelperFunc) http.HandlerFunc {
	return func(res http.ResponseWriter, r *http.Request) {
		sess, _ := cookieStore.Get(r, "cloudkeys-go")
		ctx := pongo2.Context{}

		if errFlash := sess.Flashes("error"); len(errFlash) > 0 {
			ctx["error"] = errFlash[0].(string)
		}

		template, err := f(res, r, sess, &ctx)
		if err != nil {
			http.Error(res, "An error ocurred.", http.StatusInternalServerError)
			fmt.Printf("ERR: %s\n", err)
			return
		}

		if template != nil {
			ts := pongo2.NewSet("frontend")
			ts.SetBaseDirectory("templates")
			tpl, err := ts.FromFile(*template)
			if err != nil {
				fmt.Printf("ERR: Could not parse template '%s': %s\n", *template, err)
				http.Error(res, "An error ocurred.", http.StatusInternalServerError)
				return
			}
			out, err := tpl.Execute(ctx)
			if err != nil {
				fmt.Printf("ERR: Unable to execute template '%s': %s\n", *template, err)
				http.Error(res, "An error ocurred.", http.StatusInternalServerError)
				return
			}

			res.Write([]byte(out))
		}
	}
}

func simpleTemplateOutput(template string) httpHelperFunc {
	return func(res http.ResponseWriter, r *http.Request, session *sessions.Session, ctx *pongo2.Context) (*string, error) {
		return &template, nil
	}
}

func stringPointer(s string) *string {
	return &s
}
