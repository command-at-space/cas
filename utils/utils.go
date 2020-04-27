/* */

package utils

import (
	"casServer/login/store"
	"encoding/json"
	"log"
	"net/http"
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
