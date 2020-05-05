/* */

package utils

import (
	store "casServer/login/store"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
)

// RequestError ...
type RequestError struct {
	Message    string `json:"message"`
	StatusCode int    `json:"statusCode"`
}

// ResponseInfo ...
type ResponseInfo struct {
	IsLogged bool        `json:"isLogged"`
	User     *store.User `json:"profile"`
}

// ResponseAnonymous ...
type ResponseAnonymous struct {
	IsLogged bool   `json:"isLogged"`
	Nick     string `json:"profile"`
}

// SendErrorToClient ...
func SendErrorToClient(w http.ResponseWriter, re *RequestError) {
	w.WriteHeader(re.StatusCode)
	w.Header().Set("Content-Type", "application/json")
	var dataJSON = []byte(`{}`)
	dataJSON, err := json.MarshalIndent(re, "", "  ")
	if err != nil {
		log.Printf("ERROR Marshaling %s\n", err)
		w.Write([]byte(`{}`))
		return
	}
	w.Write(dataJSON)
}

// SendJSONToClient ...
func SendJSONToClient(w http.ResponseWriter, d interface{}) {
	w.Header().Set("Content-Type", "application/json")
	var dataJSON = []byte(`{}`)
	dataJSON, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		log.Printf("ERROR Marshaling %s\n", err)
		w.Write([]byte(`{}`))
	}
	w.Write(dataJSON)
}

// BadRequest ...
func BadRequest(w http.ResponseWriter, r *http.Request) {
	re := &RequestError{
		Message:    "Bad Request",
		StatusCode: 404,
	}
	SendErrorToClient(w, re)
}

// CheckAppMode ...
func CheckAppMode() (mode string) {
	serverName, _ := os.Hostname()
	serverName = strings.ToLower(serverName)
	for _, value := range envsDev {
		if value == serverName {
			return "dev"
		}
	}
	return "production"
}

// CheckModeForCookieDomain ...
func CheckModeForCookieDomain() (domain string) {
	if CheckAppMode() == "dev" {
		return "localhost"
	}
	return domainNameForCookie
}

// CheckModeForCookieHTTPOnly ...
func CheckModeForCookieHTTPOnly() (httpOnly bool) {
	if CheckAppMode() == "dev" {
		return false
	}
	return true
}

const (
	validChars string = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// CheckValidCharacters ...
func CheckValidCharacters(str string) bool {
	str = strings.ToLower(str)
	for _, char := range str {
		if !strings.Contains(validChars, string(char)) {
			return false
		}
	}
	return true
}
