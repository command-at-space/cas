/* */

package login

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"golang.org/x/crypto/bcrypt"
)

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

func loadConfig(configjson []byte, c interface{}) {
	err := json.Unmarshal(configjson, &c)
	if err != nil {
		log.Printf("ERROR LoadConfig %s\n", err)
	}
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	re := &requestError{
		Error:      fmt.Errorf("Unexistent Endpoint " + (r.URL).String()),
		Message:    "Bad Request",
		StatusCode: 400,
	}
	log.Println(re.Error)
	sendErrorToClient(w, re)
}

func sendErrorToClient(w http.ResponseWriter, re *requestError) {
	w.WriteHeader(re.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	var dataJSON = []byte(`{}`)
	dataJSON, err := json.MarshalIndent(re, "", " ")
	if err != nil {
		log.Printf("ERROR Marshaling %s\n", err)
		w.Write([]byte(`{}`))
	}
	w.Write(dataJSON)
}

func sendJSONToClient(w http.ResponseWriter, d interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var dataJSON = []byte(`{}`)
	dataJSON, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		log.Printf("ERROR Marshaling %s\n", err)
		w.Write([]byte(`{}`))
	}
	w.Write(dataJSON)
}

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
