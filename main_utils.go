/* */

package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
)

func checkFlags() {
	versionFlag := flag.Bool("v", false, "Show current version and exit")
	flag.Parse()
	switch {
	case *versionFlag:
		fmt.Printf("Version:\t: %s\n", version)
		fmt.Printf("Date   :\t: %s\n", when)
		os.Exit(0)
	}
}

func loadConfigJSON(a *app) {
	err := json.Unmarshal(getGlobalConfigJSON(), a)
	if err != nil {
		log.Fatal("Error parsing JSON config => \n", err)
	}
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

func createCustomErrorLogFile(f string) *os.File {
	mylog, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("ERROR opening Error log file %s\n", err)
	}
	log.SetOutput(mylog)
	return mylog
}

func createCustomInfoLogFile(a *app) {
	var f = a.Conf.InfoLogFile
	infoLog, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
	if err != nil {
		log.Fatalf("ERROR opening Info log file %s\n", err)
	}
	a.iLog = log.New(infoLog, "INFO :\t", log.Ldate|log.Ltime)
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	re := &requestError{
		Error:      errors.New("Unexistent Endpoint " + (r.URL).String()),
		Message:    "Bad Request",
		StatusCode: 400,
	}
	log.Println(re.Error)
	sendErrorToClient(w, re)
}

func secret(w http.ResponseWriter, r *http.Request) {
	//fmt.Println(`SECRET`)
	w.Write([]byte("SECRET - LOGGED ZONE"))
}
