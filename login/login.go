/* */

package login

import (
	store "casServer/login/store"
	util "casServer/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

const (
	// SessionLength ...
	SessionLength int = 50
	bcryptCost    int = 12
	// CookieName ...
	CookieName string = "alphaCAS"
)

var loginDB *store.DB

func init() {
	db, err := store.NewDB(util.CheckAppMode())
	if err != nil {
		log.Fatal("Error connecting DataBase => ", err)
	}
	loginDB = db
	//users, _ := loginDB.AccountList()
	//fmt.Println(len(users))

}

// Login ...
func Login(w http.ResponseWriter, r *http.Request) {
	e := new(util.RequestError)
	u := store.NewUser()
	r.ParseForm()
	username := strings.ToLower(r.FormValue("user"))
	pass := r.Form.Get("pass")
	u, err := loginDB.Account(username)
	if err != nil {
		if err.Error() == "sql: no rows in result set" {
			e.Message = "User doesnt exist or incorrect password"
			e.StatusCode = 401
			util.SendErrorToClient(w, e)
			return
		}
		e.Message = "Error doing login... Please try again"
		e.StatusCode = 401
		util.SendErrorToClient(w, e)
		return
	}
	if !comparePass(u.PassHashed, []byte(pass)) {
		e.Message = "User doesnt exist or incorrect password"
		e.StatusCode = 401
		util.SendErrorToClient(w, e)
		return
	}
	s := store.NewSession()
	s.NickID = u.NickID
	setSessionCookie(w, s, u.Nick)
	err = loginDB.SaveSession(s)
	if err != nil {
		e.Message = "We are experiencing problems... Please try again later"
		e.StatusCode = 503
		util.SendErrorToClient(w, e)
		return
	}
	info := &util.ResponseInfo{
		IsLogged: true,
		User:     u,
	}
	util.SendJSONToClient(w, info)
}

// Logout ...
func Logout(w http.ResponseWriter, r *http.Request) {
	e := new(util.RequestError)
	r.ParseForm()
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		info := &util.ResponseInfo{
			IsLogged: false,
			User:     nil,
		}
		util.SendJSONToClient(w, info)
		return
	}
	username := strings.Split(cookie.Value, ":")[0]
	usernameID := strings.ToLower(username)
	deleteSessionCookie(w)
	err = loginDB.DeleteSession(usernameID)
	if err != nil {
		e.Message = "We are experiencing problems... Please try again later"
		e.StatusCode = 503
		util.SendErrorToClient(w, e)
		return
	}
	info := &util.ResponseInfo{
		IsLogged: false,
		User:     nil,
	}
	util.SendJSONToClient(w, info)
}

// AutoLogin ...
func AutoLogin(w http.ResponseWriter, r *http.Request) {
	e := new(util.RequestError)
	if r.Method != "POST" {
		util.BadRequest(w, r)
		return
	}
	r.ParseForm()
	cookie, err := r.Cookie(CookieName)
	if err != nil {
		//log.Print("No cookie")
		info := &util.ResponseInfo{
			IsLogged: false,
			User:     nil,
		}
		util.SendJSONToClient(w, info)
		return
	}
	username := strings.Split(cookie.Value, ":")[0]
	usernameID := strings.ToLower(username)
	s, err := loginDB.UserSession(usernameID)
	if err != nil {
		e.Message = "We are experiencing problems... Please try again later"
		e.StatusCode = 500
		util.SendErrorToClient(w, e)
		return
	}
	if s.SessionID != cookie.Value {
		info := &util.ResponseInfo{
			IsLogged: false,
			User:     nil,
		}
		util.SendJSONToClient(w, info)
		return
	}
	u := store.NewUser()
	u.Nick = username
	info := &util.ResponseInfo{
		IsLogged: true,
		User:     u,
	}
	util.SendJSONToClient(w, info)
}

// SignUp ...
func SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != "POST" {
		util.BadRequest(w, r)
		return
	}
	e := new(util.RequestError)

	// create new user validating data
	u := &store.User{
		Nick:         r.FormValue("user"),
		NickID:       strings.ToLower(r.FormValue("user")),
		PassHashed:   r.Form.Get("pass"),
		Email:        r.FormValue("mail"),
		Verified:     0,
		Logo:         r.FormValue("logo"),
		SecretQuest:  r.FormValue("ques"),
		SecretHashed: r.FormValue("secr"),
		CreatedAt:    time.Now().UTC(),
		LastSeen:     time.Now().UTC(),
		Online:       0,
	}
	badData := validateNewUserData(u)
	if badData != "" {
		//fmt.Println(badData)
		e.Message = fmt.Sprint(badData)
		e.StatusCode = 409
		util.SendErrorToClient(w, e)
		return
	}
	u.PassHashed = saltAndHash([]byte(r.FormValue("pass")), bcryptCost)
	if u.SecretQuest != "" {
		u.SecretHashed = saltAndHash([]byte(r.Form.Get("secr")), bcryptCost)
	}
	//fmt.Printf("%+v", u)

	// check username is available
	existUser, err := loginDB.Account(u.NickID)
	if err != nil {
		e.Message = "We have some problems. Please Try Later"
		e.StatusCode = 500
		util.SendErrorToClient(w, e)
		return
	}
	if existUser.NickID == u.NickID {
		e.Message = fmt.Sprintf("User %s already exist", u.Nick)
		e.StatusCode = 409
		util.SendErrorToClient(w, e)
		return
	}

	// insert new user and notify client
	err = loginDB.NewAccount(u)
	if err != nil {
		e.Message = "We are experiencing problems. Try Later please"
		e.StatusCode = 500
		util.SendErrorToClient(w, e)
		return
	}
	info := struct {
		Created bool   `json:"created"`
		Error   string `json:"error"`
	}{
		Created: true,
		Error:   "",
	}

	util.SendJSONToClient(w, info)
}
