/* */

package login

import (
	store "cas-server/store"
	"fmt"
	"log"
	"net/http"
	"strings"
)

type requestError struct {
	Error      error  `json:"-"`
	Message    string `json:"message"`
	StatusCode int    `json:"-"`
}

type responseInfo struct {
	IsLogged bool        `json:"isLogged"`
	User     *store.User `json:"profile"`
}

// Login ...
func Login(w http.ResponseWriter, r *http.Request, db *store.DB) {
	e := new(requestError)
	u := store.NewUser()
	r.ParseForm()
	username := r.FormValue("user")
	pass := r.Form.Get("pass")
	u, err := db.Account(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			e.Error = fmt.Errorf(`User %s doesnt exist`, username)
			e.Message = "User doesnt exist or incorrect password"
			e.StatusCode = 401
			sendErrorToClient(w, e)
			return
		}
		e.Error = fmt.Errorf("Error doing login... Please try again")
		e.Message = "Error doing login... Please try again"
		e.StatusCode = 401
		sendErrorToClient(w, e)
		return
	}
	if !comparePass(u.Hash, []byte(pass)) {
		e.Error = fmt.Errorf(`Incorrect Password for %s`, u.Name)
		e.Message = "User doesnt exist or incorrect password"
		e.StatusCode = 401
		sendErrorToClient(w, e)
		return
	}
	cookie := setSessionCookie(w, u.Name)
	sessionID := strings.Split(cookie, ":")[1]
	ActiveUsers[u.Name] = sessionID
	info := &responseInfo{
		true,
		u,
	}
	sendJSONToClient(w, info)
}

// Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Print("No cookie")
		info := &responseInfo{
			false,
			nil,
		}
		sendJSONToClient(w, info)
		return
	}
	user := strings.Split(cookie.Value, ":")[0]
	deleteSessionCookie(w)
	delete(ActiveUsers, user)
	info := &responseInfo{
		false,
		nil,
	}
	sendJSONToClient(w, info)
}

// AutoLogin ...
func AutoLogin(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		badRequest(w, r)
		return
	}
	r.ParseForm()
	cookie, err := r.Cookie(cookieName)
	if err != nil {
		log.Print("No cookie")
		info := &responseInfo{
			false,
			nil,
		}
		sendJSONToClient(w, info)
		return
	}
	user := strings.Split(cookie.Value, ":")[0]
	sessionID := strings.Split(cookie.Value, ":")[1]
	value, ok := ActiveUsers[user]
	if !ok || value != sessionID {
		info := &responseInfo{
			false,
			nil,
		}
		sendJSONToClient(w, info)
		return
	}
	//loginDB.saveSession(u.name, sessionID)
	//ActiveUsers[user] = sessionID
	u := store.NewUser()
	u.Name = user
	u.Hash = sessionID
	info := &responseInfo{
		true,
		u,
	}
	sendJSONToClient(w, info)
}

// SignUp ..
func SignUp(w http.ResponseWriter, r *http.Request, db *store.DB) {
	//fmt.Println("Creating User ...")
	if r.Method != "POST" {
		badRequest(w, r)
		return
	}
	u := store.NewUser()
	r.ParseForm()
	//r.ParseMultipartForm(10000)
	u.Name = r.FormValue("user")
	u.Hash = saltAndHash([]byte(r.Form.Get("pass")), bcryptCost)
	u.Email = r.Form.Get("mail")
	u.Logo = r.Form.Get("logo")
	// already exist user ?
	existUser, err := db.Account(u.Name)
	if err != nil {
		info := struct {
			created bool
			error   string
		}{
			created: false,
			error:   err.Error(),
		}
		sendJSONToClient(w, info)
		return
	}
	if existUser.Name == u.Name {
		e := new(requestError)
		e.Error = fmt.Errorf("User %s already exist", u.Name)
		e.Message = "User already exist"
		e.StatusCode = 400
		sendErrorToClient(w, e)
		return
	}
	//fmt.Printf(`%v`, u)
	err = db.NewAccount(u)
	if err != nil {
		e := new(requestError)
		e.Error = err
		e.Message = "Error doing login... Please try again"
		e.StatusCode = 503
		sendErrorToClient(w, e)
		return
	}
	info := struct {
		created bool
		error   string
	}{
		created: true,
		error:   "",
	}
	sendJSONToClient(w, info)
}
