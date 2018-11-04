package main

import (
	"encoding/json"
	"net/http"

	"github.com/gofrs/uuid"
	log "github.com/sirupsen/logrus"
)

type userState string

const (
	userStateLoggedin   userState = "logged-in"
	userStateRequireMFA           = "require-mfa"
)

type apiError struct {
	cause error
	msg   string
}

func (a apiError) Error() string       { return a.msg + ": " + a.cause.Error() }
func (a apiError) Cause() error        { return a.cause }
func (a apiError) UserMessage() string { return a.msg }
func wrapAPIError(cause error, msg string) error {
	if cause == nil {
		return nil
	}
	return apiError{cause, msg}
}

func apiHelper(hdl apiHandler) http.HandlerFunc {
	return func(res http.ResponseWriter, r *http.Request) {
		cookieSession, err := cookieStore.Get(r, "cloudkeys-go")
		if err != nil {
			log.WithError(err).Debug("Session could not be decoded, created new one")
		}

		sess := newSessionData()
		if !cookieSession.IsNew {
			if err := json.Unmarshal(cookieSession.Values["sessionData"].([]byte), sess); err != nil {
				log.WithError(err).Debug("Session cookie contained garbled sessionData")
				// Session data is garbled, create a new session in case
				// something was decoded into the session object
				sess = newSessionData()
			}
		}

		// This is a pure JSON API
		res.Header().Set("Content-Type", "application/json")
		res.Header().Set("Cache-Control", "no-cache")
		res.Header().Set("X-API-Version", version)
		res.Header().Set("Access-Control-Allow-Origin", "*") // FIXME (kahlers): Remove after development

		// Assign an UUID to find potiential errors in the logs
		reqId := uuid.Must(uuid.NewV4()).String()

		var (
			resp   interface{}
			status int
			data   []byte
		)

		// Define a common error handler
		defer func() {
			if err != nil {
				log.WithFields(log.Fields{
					"reqId":  reqId,
					"method": r.Method,
					"path":   r.URL.Path,
				}).WithError(err).Error("API handler errored")

				// Respond with common error format
				res.WriteHeader(status)
				json.NewEncoder(res).Encode(map[string]interface{}{
					"error":   err.(apiError).UserMessage(),
					"reqId":   reqId,
					"success": false,
				})
			}
		}()

		// Do the real work
		resp, status, err = hdl(getContext(r), res, r, sess)
		if err != nil {
			return
		}

		// If there was no error try to marshal the output and wrap error
		// when this fails in order to have a propoer user message
		data, err = json.Marshal(resp)
		err = wrapAPIError(err, "Could not marshal API response")
		if err != nil {
			return
		}

		// Write the session data back to the cookie
		sdata, err := json.Marshal(sess)
		err = wrapAPIError(err, "Failed to encode session")
		if err != nil {
			return
		}

		cookieSession.Values["sessionData"] = sdata
		err = wrapAPIError(cookieSession.Save(r, res), "Failed to save session")
		if err != nil {
			return
		}

		// If no error ocurred, send the response
		res.WriteHeader(status)
		res.Write(data)
	}
}
