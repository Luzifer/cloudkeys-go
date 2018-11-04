package main

type sessionData struct {
	// The current state of all logged in users
	Users map[string]userState
	// MFA secrets are encrypted with users password, we need to
	// store them inside the encrypted session until the user verified
	// themselves
	MFACache map[string]string
}

func newSessionData() *sessionData {
	return &sessionData{
		Users:    make(map[string]userState),
		MFACache: make(map[string]string),
	}
}
