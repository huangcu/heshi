package util

import (
	"crypto/rand"
	"io"
	"log"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

func Encrypt(pass string) string {
	// Generate "hash" to store from user password
	hash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		// TODO: Properly handle error
		log.Fatal(err)

	}
	return string(hash)
	// fmt.Println("Hash to store:", string(hash))
}

func IsPassOK(pass, hash string) bool {
	// Store this "hash" somewhere, e.g. in your database
	// After a while, the user wants to log in and you need to check the password he entered
	// userPassword2 := "some user-provided password"
	// hashFromDatabase := "query db from user get the hash"
	// Comparing the password with the hash
	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(pass)); err != nil {
		// TODO: Properly handle error
		log.Fatal(err)
		if strings.Contains(err.Error(), "not right") {
			return false
		}
	}
	return true
}

var StdChars = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789!@#$%^&*()-_=+,.?/:;{}[]`~")

func NewPassword(length int) string {
	// var username string
	// q := fmt.Sprintf("SELECT password from users where username=%s", username)
	// db.Exec(q)
	return randChar(length, StdChars)
}

func randChar(length int, chars []byte) string {
	newPassword := make([]byte, length)
	randomData := make([]byte, length+(length/4)) // storage for random bytes.
	clen := byte(len(chars))
	maxrb := byte(256 - (256 % len(chars)))
	i := 0
	for {
		if _, err := io.ReadFull(rand.Reader, randomData); err != nil {
			panic(err)
		}
		for _, c := range randomData {
			if c >= maxrb {
				continue
			}
			newPassword[i] = chars[c%clen]
			i++
			if i == length {
				return string(newPassword)
			}
		}
	}
}
