/* */

package login

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func generateRandomString(s int) (string, error) {
	b, err := generateRandomBytes(s)
	return base64.URLEncoding.EncodeToString(b), err
}

func generateRandomBytes(n int) ([]byte, error) {
	b := make([]byte, n)
	_, err := rand.Read(b)
	if err != nil {
		return nil, err
	}
	return b, nil
}

func saltAndHash(pwd []byte, cost int) string {
	hash, err := bcrypt.GenerateFromPassword(pwd, cost)
	if err != nil {
		log.Print(err)
	}
	return string(hash)
}

func comparePass(hashedPass string, plainPass []byte) bool {
	byteHash := []byte(hashedPass)
	err := bcrypt.CompareHashAndPassword(byteHash, plainPass)
	if err != nil {
		//log.Print("Incorrect Password")
		return false
	}
	return true
}

func sha2(str string) string {
	bytes := []byte(str)
	// Converts string to sha2
	h := sha256.New()                   // new sha256 object
	h.Write(bytes)                      // data is now converted to hex
	code := h.Sum(nil)                  // code is now the hex sum
	codestr := hex.EncodeToString(code) // converts hex to string
	return codestr
}
