package main

/*
GOOS=linux GOARCH=amd64 go build -o cas
*/
import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	login "./_login"
	//db "./_db"
	lib "./_lib"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{} // use default options

func init() {
	lib.LoadConfig(configjson, &c)
	if c.App.Mode != "production" {
		c.App.Port = 3000
	}
}

func main() {
	////////////// SEND LOGS TO FILE //////////////////
	if c.App.Mode == "production" {
		var f = c.App.ErrLog
		mylog, err := os.OpenFile(f, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err != nil {
			log.Printf("ERROR opening log file %s\n", err)
		}
		defer mylog.Close() // defer must be in main
		log.SetOutput(mylog)
	}
	///////////////////////////////////////////////////
	//db.MYDB.InitDB()

	http.DefaultClient.Timeout = 10 * time.Second
	mux := http.NewServeMux()
	mux.HandleFunc("/v0/ws", ws)
	mux.HandleFunc("/v0/create", login.CreateAccount)
	mux.HandleFunc("/v0/logout", login.IsLogged(login.Logout))
	mux.HandleFunc("/v0/login", login.Login)
	mux.HandleFunc("/v0/secret", login.IsLogged(secret))
	mux.HandleFunc("/", badRequest)

	server := http.Server{
		Addr:           fmt.Sprintf("localhost:%d", c.App.Port),
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	log.Printf("Server listening ...%s", server.Addr)
	server.ListenAndServe()
}

func secret(w http.ResponseWriter, r *http.Request) {
	fmt.Println(`SECRET`)
	w.Write([]byte("LOGGED ZONE"))
}

func ws(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}

func badRequest(w http.ResponseWriter, r *http.Request) {
	e.Error = fmt.Sprintf("Inexistent Endpoint...%s", r.URL.RequestURI())
	if c.App.Mode != "test" {
		log.Printf("ERROR = %s", e.Error)
	}
	lib.SendErrorToClient(w, e)
}
