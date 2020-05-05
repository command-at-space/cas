/* */

package online

import (
	login "casServer/login"
	store "casServer/login/store"
	util "casServer/utils"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"
)

// PJs ...
type PJs struct {
	Online map[string]string
}

// UserAnonymous ...
/*type UserAnonymous struct {
	Nick string `json:"nick"`
}*/

// NewPlayersOnline ...
func NewPlayersOnline() *PJs {
	var pjs PJs
	pjs.Online = make(map[string]string)
	pjs.connect("CaS", "adminGameBot")
	return &pjs
}

func (pjs *PJs) connect(nick, cookie string) {
	pjs.Online[nick] = cookie
}

func (pjs *PJs) disconnect(nick string) {
	if pjs.IsConnected(nick) {
		delete(pjs.Online, nick)
	}
}

// IsConnected ...
func (pjs *PJs) IsConnected(nick string) bool {
	_, ok := pjs.Online[nick]
	if ok {
		return true
	}
	return false
}

func (pjs *PJs) listAll() {
	for nick, _ := range pjs.Online {
		fmt.Printf("%s\n", nick)
		//fmt.Printf("%s:%s\n", nick, cookie)
	}
}

// JoinGame ...
func (pjs *PJs) JoinGame(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, err := r.Cookie(login.CookieName)
	if err != nil {
		info := &util.ResponseInfo{
			IsLogged: false,
			User:     nil,
		}
		util.SendJSONToClient(w, info)
		return
	}
	username := strings.Split(cookie.Value, ":")[0]
	value := strings.Split(cookie.Value, ":")[1]
	pjs.connect(username, value)
	u := store.NewUser()
	u.Nick = username
	info := &util.ResponseInfo{
		IsLogged: true,
		User:     u,
	}
	util.SendJSONToClient(w, info)
}

// JoinAnonymous ...
func (pjs *PJs) JoinAnonymous(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	nick := strings.ToLower(r.FormValue("user"))
	u := store.NewUser()
	u.Nick = nick
	// check not exists
	badData := validateNewAnonymousData(u.Nick)
	if badData != "" {
		e := new(util.RequestError)
		e.Message = fmt.Sprint(badData)
		e.StatusCode = 409
		util.SendErrorToClient(w, e)
		return
	}
	token, err := login.GenerateRandomString(login.SessionLength)
	if err != nil {
		log.Print("Error Generating Random String")
		token = time.Now().String()
	}
	sessionID := nick + ":" + token
	expires := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:     login.CookieName,
		Value:    sessionID,
		Domain:   util.CheckModeForCookieDomain(), //"localhost",
		Path:     "/",
		HttpOnly: util.CheckModeForCookieHTTPOnly(), //false,
		Expires:  expires,
	}
	http.SetCookie(w, cookie)
	username := strings.Split(cookie.Value, ":")[0]
	value := strings.Split(cookie.Value, ":")[1]
	pjs.connect(username, value)
	info := &util.ResponseInfo{
		IsLogged: true,
		User:     u,
	}
	util.SendJSONToClient(w, info)
}

// IsOnline ...
func (pjs *PJs) IsOnline(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		cookie, err := r.Cookie(login.CookieName)
		if err != nil { // No cookie
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		username := strings.Split(cookie.Value, ":")[0]
		if pjs.IsConnected(username) {
			next.ServeHTTP(w, r)
			return
		}
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
	}
}
