/* */
// go build -ldflags="-X 'main.when=$(date -u +%F_%T)'"

package main

import (
	login "cas-server/login"
	store "cas-server/store"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"
)

var version = "v0.0.4"
var when = ""

type app struct {
	Conf struct {
		Mode          string `json:"mode"`
		Host          string `json:"host"`
		Port          int    `json:"port"`
		ErrorsLogFile string `json:"errorsLogFile"`
		InfoLogFile   string `json:"infoLogFile"`
	} `json:"config"`
	db   *store.DB
	iLog *log.Logger
}

type requestError struct {
	Error      error  `json:"-"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

func main() {
	checkFlags()
	rand.Seed(time.Now().UnixNano())

	// Load Conf
	var a app
	loadConfigJSON(&a)
	checkMode(&a)
	fmt.Println(a.Conf.Port)

	// Custom Error Log File + Custom Info Log File
	createCustomInfoLogFile(&a)
	var mylog *os.File
	if a.Conf.Mode == "production" {
		mylog = createCustomErrorLogFile(a.Conf.ErrorsLogFile)
	}
	defer mylog.Close()

	// DataBase
	db, err := store.NewDB(a.Conf.Mode)
	if err != nil {
		log.Fatal("Error connecting DataBase => ", err)
	}
	a.db = db

	// Server
	http.DefaultClient.Timeout = 5 * time.Second
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/login", login.IsNotLogged(
		func(w http.ResponseWriter, r *http.Request) {
			login.Login(w, r, a.db)
		},
	))
	mux.HandleFunc("/auth/autoLogin", login.AutoLogin)
	mux.HandleFunc("/auth/signup", login.IsNotLogged(
		func(w http.ResponseWriter, r *http.Request) {
			login.SignUp(w, r, a.db)
		},
	))
	mux.HandleFunc("/auth/logout", login.IsLogged(login.Logout))
	mux.HandleFunc("/secret", login.IsLogged(secret))
	mux.HandleFunc("/", badRequest)

	server := http.Server{
		Addr:           fmt.Sprintf("%s:%d", a.Conf.Host, a.Conf.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Server up listening %s in mode %s", server.Addr, a.Conf.Mode)
	server.ListenAndServe()
}
