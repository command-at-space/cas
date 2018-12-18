package login

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

func init() {
	fmt.Println(`Login package init ...`)
	rand.Seed(time.Now().UnixNano())
	loadConfig(mysqljson, &c)
	loadConfig(configjson, &c)
	loginDB.initDB()
	err := loginDB.loadAllSessions()
	if err != nil {
		log.Fatalf("Problem loading Sessions %s\n", err)
	}
	fmt.Println(`USERS ACTIVE `, len(ActiveUsers))
	for user, session := range ActiveUsers {
		//fmt.Printf("%s <-> %s\n", user, session)
		fmt.Sprintln(user, session)
	}

}

// CreateAccount ..
func CreateAccount(w http.ResponseWriter, r *http.Request) {
	var u user
	if r.Method != "POST" {
		e.Text = "GET /v0/create is not a valid endpoint"
		sendErrorToClient(w, e)
		return
	}
	r.ParseForm()
	//r.ParseMultipartForm(10000)
	u.name = r.FormValue("user1")
	u.hash = saltAndHash([]byte(r.Form.Get("pass1")))
	u.email = r.Form.Get("mail1")
	u.logo = r.Form.Get("logo1")
	//fmt.Printf(`%v`, u)
	err := loginDB.insertNewAccount(u)
	if err != nil { // user already exists
		if e.Text != fmt.Sprintf("User %s already exist", u.name) {
			e.Text = err.Error()
		}
		sendErrorToClient(w, e)
		return
	}
	logged(w, r, u)
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	var u user
	if r.Method != "POST" {
		e.Text = "GET /v0/login is not a valid endpoint"
		sendErrorToClient(w, e)
		return
	}
	r.ParseForm()
	username := r.FormValue("user2")
	pass := []byte(r.Form.Get("pass2"))
	//fmt.Printf(`%v`, u)
	u, err := loginDB.getAccount(username)
	fmt.Println(`MEHHH`, u)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			e.Text = fmt.Sprintf(`User %s doesnt exist`, username)
			sendErrorToClient(w, e)
			return
		}
		e.Text = "Error doing login... Please try again"
		sendErrorToClient(w, e)
		return
	}
	if !comparePass(u.hash, pass) {
		e.Text = fmt.Sprintf(`Incorrect Password for %s`, u.name)
		sendErrorToClient(w, e)
		return
	}
	logged(w, r, u)
}

func logged(w http.ResponseWriter, r *http.Request, u user) {
	cookie := setSessionCookie(w, u.name)
	sessionID := strings.Split(cookie, ":")[1]
	loginDB.saveSession(u.name, sessionID)
	ActiveUsers[u.name] = sessionID
	data := struct {
		Location string
	}{
		"//localhost:8080/play.html",
	}
	sendJSONToClient(w, data)
}

// Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//valid := r.Form.Get("test")
	cookie, err := r.Cookie(c.Auth.CookieName)
	if err != nil {
		log.Print("No cookie")
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}
	user := strings.Split(cookie.Value, ":")[0]
	fmt.Println(`LOGOUT => `, user)
	loginDB.deleteSession(user)
	delete(ActiveUsers, user)
	fmt.Println(`USERS ACTIVE `, len(ActiveUsers))
	data := struct {
		Location string
	}{
		"//localhost:8080",
	}
	sendJSONToClient(w, data)
}
