package main

import (
	"context"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"

	"github.com/gorilla/mux"
	"github.com/pquerna/otp/totp"

	openssl "github.com/Luzifer/go-openssl"
)

type apiHandler func(context.Context, http.ResponseWriter, *http.Request, *sessionData) (interface{}, int, error)

// apiChangeLoginPassword accepts three passwords: old, new and
// new-repeat. It checks the old password, compares the new ones
// and if they match it sets a new user password. In case the user
// has a MFA secret it is re-encrypted with the new password.
func apiChangeLoginPassword(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	var (
		input = struct {
			OldPassword    string `json:"old_password"`
			NewPassword    string `json:"new_password"`
			RepeatPassword string `json:"repeat_password"`
		}{}
		output   = map[string]interface{}{"success": true}
		username = mux.Vars(r)["user"]
	)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, wrapAPIError(err, "Unable to decode login data")
	}

	if state, ok := sess.Users[username]; !ok || state != userStateLoggedin {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("Access to user not logged in"), "Authorization error")
	}

	// Check new passwords matches
	if input.NewPassword != input.RepeatPassword {
		return nil, http.StatusBadRequest, wrapAPIError(errors.New("Password mismatch"), "New passwords do not match")
	}

	// Retrieve data file
	userFile := createUserFilename(username)
	user, err := dataObjectFromStorage(ctx, storage, userFile)
	if err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Unable to retrieve data file")
	}

	// Check bcrypt password and deprecated version of password
	deprecatedPassword := fmt.Sprintf("%x", sha1.Sum([]byte(cfg.PasswordSalt+input.OldPassword))) // Here for backwards compatibility
	if bcrypt.CompareHashAndPassword([]byte(user.MetaData.Password), []byte(input.OldPassword)) != nil &&
		user.MetaData.Password != deprecatedPassword {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("Password mismatch"), "Authorization error")
	}

	// Update user password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Unable to generate bcrypt hash")
	}
	user.MetaData.Password = string(hashedPassword)

	// In case a MFA token is present re-encrypt it with the new password
	if user.MetaData.MFASecret != "" {
		secret, err := openssl.New().DecryptBytes(input.OldPassword, []byte(user.MetaData.MFASecret), openssl.DigestSHA256Sum)
		if err != nil {
			return nil, http.StatusInternalServerError, wrapAPIError(err, "Could not decrypt MFA secret")
		}

		secret, err = openssl.New().EncryptBytes(input.NewPassword, secret, openssl.DigestSHA256Sum)
		user.MetaData.MFASecret = string(secret)
	}

	// Save back the user file
	if err := user.writeToStorage(ctx, storage, userFile); err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Could not write data file")
	}

	return output, http.StatusOK, nil
}

// apiGetUserData retrieves the user file and returns the encrypted
// data together with the current checksum of the content
func apiGetUserData(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	var username = mux.Vars(r)["user"]

	if state, ok := sess.Users[username]; !ok || state != userStateLoggedin {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("Access to user not logged in"), "Authorization error")
	}

	// Retrieve data file
	userFile := createUserFilename(username)
	user, err := dataObjectFromStorage(ctx, storage, userFile)
	if err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Unable to retrieve data file")
	}

	return map[string]interface{}{
		"checksum": user.MetaData.Version,
		"data":     user.Data,
	}, http.StatusOK, nil
}

func apiGetUserSettings(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	// FIXME (kahlers): Implement this
	return nil, http.StatusInternalServerError, wrapAPIError(errors.New("Not implemented yet"), "Not implemented")
}

// apiListUsers returns a dictionary of usernames with their current login state
func apiListUsers(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	return sess.Users, http.StatusOK, nil
}

// apiLogin retrieves an username and a password, loads the user file
// from storage and compares the passwords. If the user has no MFA they
// are logged in, otherwise they require MFA auth. After login the user
// file is automatically migrated to the latest version.
func apiLogin(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	var (
		input = &struct {
			Username string
			Password string
		}{}
		output = map[string]interface{}{"success": true}
	)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, wrapAPIError(err, "Unable to decode login data")
	}

	if _, ok := sess.Users[input.Username]; ok {
		// Already logged in
		return output, http.StatusOK, nil
	}

	userFile := createUserFilename(input.Username)
	if !storage.IsPresent(ctx, userFile) {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("Userfile not present"), "Authorization error")
	}

	// Retrieve data file
	user, err := dataObjectFromStorage(ctx, storage, userFile)
	if err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Unable to retrieve data file")
	}

	// Check bcrypt password and deprecated version of password
	deprecatedPassword := fmt.Sprintf("%x", sha1.Sum([]byte(cfg.PasswordSalt+input.Password))) // Here for backwards compatibility
	if bcrypt.CompareHashAndPassword([]byte(user.MetaData.Password), []byte(input.Password)) != nil &&
		user.MetaData.Password != deprecatedPassword {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("Password mismatch"), "Authorization error")
	}

	// Apply migrations to data file automatically
	if err := user.migrate(ctx, storage, userFile, input.Password); err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Migrating the data file caused an error")
	}

	// Set user as logged in
	if user.MetaData.MFASecret == "" {
		sess.Users[input.Username] = userStateLoggedin
	} else {
		secret, err := openssl.New().DecryptBytes(input.Password, []byte(user.MetaData.MFASecret), openssl.DigestSHA256Sum)
		if err != nil {
			return nil, http.StatusInternalServerError, wrapAPIError(err, "Could not decrypt MFA secret")
		}

		sess.Users[input.Username] = userStateRequireMFA
		sess.MFACache[input.Username] = string(secret)
	}

	return output, http.StatusOK, nil
}

// apiLogoutUser removes the given user from the list of logged in users
// and forgets the MFA secret in case it was still set.
func apiLogoutUser(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	var (
		output = map[string]interface{}{"success": true}
		user   = mux.Vars(r)["user"]
	)

	if _, ok := sess.MFACache[user]; ok {
		delete(sess.MFACache, user)
	}

	if _, ok := sess.Users[user]; ok {
		delete(sess.Users, user)
	}

	return output, http.StatusOK, nil
}

// apiRegister takes an username and two passwords, compares the
// passwords, ensures the username is not already taken and creates
// a new empty user file.
func apiRegister(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	var (
		input = struct {
			Username      string `json:"username"`
			Password      string `json:"password"`
			CheckPassword string `json:"check_password"`
		}{}
		output = map[string]interface{}{"success": true}
	)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, wrapAPIError(err, "Unable to decode request")
	}

	// Check input
	if input.Username == "" || input.Password == "" || input.Password != input.CheckPassword {
		return nil, http.StatusBadRequest, wrapAPIError(errors.New("Invalid input data"), "Invalid input provided")
	}

	// Check username collision
	userFile := createUserFilename(input.Username)
	if storage.IsPresent(ctx, userFile) {
		return nil, http.StatusBadRequest, wrapAPIError(errors.New("User file found"), "Username is already taken")
	}

	// Create user file
	user := newDataObject()
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Unable to generate bcrypt hash")
	}
	user.MetaData.Password = string(hashedPassword)

	// Save the user file
	if err := user.writeToStorage(ctx, storage, userFile); err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Could not write data file")
	}

	// Log-in the newly created user
	sess.Users[input.Username] = userStateLoggedin

	return output, http.StatusOK, nil
}

// apiSetUserData retrieves two checksums and an encrypted data blob.
// It compares the old checksum still matches to ensure no changes
// of another user is overwritten and it compares the received checksum
// of the new data blob to ensure the data integrity is fine. Afterwards
// the user file is backupped and updated.
func apiSetUserData(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	var (
		input = struct {
			Checksum    string `json:"checksum"`
			OldChecksum string `json:"old_checksum"`
			Data        string `json:"data"`
		}{}
		output   = map[string]interface{}{"success": true}
		username = mux.Vars(r)["user"]
	)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, wrapAPIError(err, "Unable to decode request")
	}

	if state, ok := sess.Users[username]; !ok || state != userStateLoggedin {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("Access to user not logged in"), "Authorization error")
	}

	// Retrieve data file
	userFile := createUserFilename(username)
	user, err := dataObjectFromStorage(ctx, storage, userFile)
	if err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Unable to retrieve data file")
	}

	// Check the user is still updating the same version of the data
	if user.MetaData.Version != input.OldChecksum {
		return nil, http.StatusBadRequest, wrapAPIError(errors.New("Update on outdated data"), "Old data checksum does not match")
	}

	// Check we've got the data the user intended to send
	if input.Checksum != fmt.Sprintf("%x", sha256.Sum256([]byte(input.Data))) {
		return nil, http.StatusBadRequest, wrapAPIError(errors.New("Checksum mismatch on input data"), "New data checksum does not match")
	}

	// Create a backup because you know...
	if err := storage.Backup(ctx, userFile); err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Could not create backup, nothin saved")
	}

	user.MetaData.Version = input.Checksum
	user.Data = input.Data

	if err := user.writeToStorage(ctx, storage, userFile); err != nil {
		return nil, http.StatusInternalServerError, wrapAPIError(err, "Could not write data file")
	}

	return output, http.StatusOK, nil
}

func apiSetUserSettings(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	// FIXME (kahlers): Implement this
	return nil, http.StatusInternalServerError, wrapAPIError(errors.New("Not implemented yet"), "Not implemented")
}

// apiValidateMFA retrieves an OTP token and in case of a match with
// the token generated from the stored secret the user is logged in.
func apiValidateMFA(ctx context.Context, res http.ResponseWriter, r *http.Request, sess *sessionData) (interface{}, int, error) {
	var (
		input = struct {
			Token string `json:"token"`
		}{}
		output = map[string]interface{}{"success": true}
		user   = mux.Vars(r)["user"]
	)

	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return nil, http.StatusBadRequest, wrapAPIError(err, "Unable to decode login data")
	}

	state, ok := sess.Users[user]
	if ok && state != userStateLoggedin {
		// Not requesting MFA authorization
		return output, http.StatusOK, nil
	} else if !ok {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("User not logged in"), "Authorization error")
	}

	secret, ok := sess.MFACache[user]
	if !ok {
		return nil, http.StatusInternalServerError, wrapAPIError(errors.New("Missing OTP secret"), "Unable to find OTP secret")
	}
	if !totp.Validate(input.Token, secret) {
		return nil, http.StatusUnauthorized, wrapAPIError(errors.New("OTP token mismatch"), "Invalid OTP token")
	}

	// Remove secret from session and set user logged in
	delete(sess.MFACache, "user")
	sess.Users[user] = userStateLoggedin

	return output, http.StatusOK, nil
}
