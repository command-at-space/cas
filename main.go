/* */
// go build -ldflags="-X 'main.releaseDate=$(date -u +%F_%T)'"

package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"time"

	"casServer/chat"
	login "casServer/login"
	online "casServer/online"
	util "casServer/utils"
)

var version = "v0.0.7"
var releaseDate = ""
var iLog *log.Logger

type app struct {
	Conf struct {
		Mode          string `json:"mode"`
		Host          string `json:"host"`
		Port          int    `json:"port"`
		ErrorsLogFile string `json:"errorsLogFile"`
		InfoLogFile   string `json:"infoLogFile"`
		ChatLogFile   string `json:"chatLogFile"`
	} `json:"config"`
}

func main() {
	checkFlags()
	rand.Seed(time.Now().UnixNano())

	// Load Conf
	var a app
	loadConfigJSON(&a)
	a.Conf.Mode = util.CheckAppMode()

	// Custom Error Log File + Custom Info Log File
	createCustomInfoLogFile(a.Conf.InfoLogFile)
	var mylog *os.File
	if a.Conf.Mode == "production" {
		mylog = createCustomErrorLogFile(a.Conf.ErrorsLogFile)
	}
	defer mylog.Close()

	// players Online
	pjs := online.NewPlayersOnline()
	//showList(pjs)

	// chat
	ch := chat.NewChat(a.Conf.ChatLogFile)
	go ch.BroadcastMessages()

	// Server
	http.DefaultClient.Timeout = 5 * time.Second
	mux := http.NewServeMux()

	mux.HandleFunc("/auth/login", login.IsNotLogged(login.Login))
	mux.HandleFunc("/auth/autoLogin", login.AutoLogin)
	mux.HandleFunc("/auth/signup", login.IsNotLogged(login.SignUp))
	mux.HandleFunc("/auth/logout", login.IsLogged(login.Logout))

	mux.HandleFunc("/online/join", login.IsLogged(pjs.JoinGame))
	mux.HandleFunc("/online/anonymous", login.IsNotLogged(pjs.JoinAnonymous))

	mux.HandleFunc("/chat/msg", pjs.IsOnline(ch.HandleChat))

	mux.HandleFunc("/secret", login.IsLogged(secret))
	mux.HandleFunc("/", util.BadRequest)

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

func secret(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("LOGGED ZONE"))
}
