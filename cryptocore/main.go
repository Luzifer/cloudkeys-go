package main

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"syscall/js"

	openssl "github.com/Luzifer/go-openssl"
)

func main() {
	js.Global().Set("opensslDecrypt", js.NewCallback(decrypt))
	js.Global().Set("opensslEncrypt", js.NewCallback(encrypt))
	js.Global().Set("sha256sum", js.NewCallback(sha256sum))

	// Trigger custom "event"
	if js.Global().Get("cryptocoreLoaded").Type() == js.TypeFunction {
		js.Global().Call("cryptocoreLoaded")
	}

	<-make(chan struct{}, 0)
}

func decrypt(i []js.Value) {
	if len(i) != 3 {
		println("decrypt requires 3 arguments")
		return
	}

	var (
		ciphertext = i[0].String()
		password   = i[1].String()
		callback   = i[2]
	)

	o := openssl.New()

	var (
		err       error
		plaintext []byte
	)
	for _, kdf := range []openssl.DigestFunc{openssl.DigestSHA256Sum, openssl.DigestMD5Sum} {
		plaintext, err = o.DecryptBytes(password, []byte(ciphertext), kdf)
		if err != nil {
			continue
		}

		if plaintext[0] != '[' || plaintext[1] != '{' {
			// This should be the beginning, otherwise the KDF provided a broken key
			err = errors.New("Unexpected output")
			continue
		}

		break
	}

	if err != nil {
		callback.Invoke(nil, fmt.Sprintf("decrypt failed: %s", err))
		return
	}

	callback.Invoke(string(plaintext), nil)
}

func encrypt(i []js.Value) {
	if len(i) != 3 {
		println("encrypt requires 3 arguments")
		return
	}

	var (
		plaintext = i[0].String()
		password  = i[1].String()
		callback  = i[2]
	)

	o := openssl.New()
	ciphertext, err := o.EncryptBytes(password, []byte(plaintext), openssl.DigestSHA256Sum)
	if err != nil {
		callback.Invoke(nil, fmt.Sprintf("encrypt failed: %s", err))
		return
	}

	callback.Invoke(string(ciphertext), nil)
}

func sha256sum(i []js.Value) {
	if len(i) != 2 {
		println("sha256sum requires 2 arguments")
		return
	}

	var (
		plaintext = i[0].String()
		callback  = i[1]
	)

	callback.Invoke(fmt.Sprintf("%x", sha256.Sum256([]byte(plaintext))), nil)
}
