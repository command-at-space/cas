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

// OnlinePJs ...
type OnlinePJs struct {
	Online map[string]string
}

// NewPlayersOnline ...
func NewPlayersOnline() *OnlinePJs {
	var pj OnlinePJs
	pj.Online = make(map[string]string)
	pj.connect("CaS", "adminGameBot")
	return &pj
}

func (pj *OnlinePJs) connect(nick, cookie string) {
	pj.Online[nick] = cookie
}

func (pj *OnlinePJs) disconnect(nick string) {
	if pj.isConnected(nick) {
		delete(pj.Online, nick)
	}
}

func (pj *OnlinePJs) isConnected(nick string) bool {
	_, ok := pj.Online[nick]
	if ok {
		return true
	}
	return false
}

func (pj *OnlinePJs) listAll() {
	for nick, _ := range pj.Online {
		fmt.Printf("%s\n", nick)
		//fmt.Printf("%s:%s\n", nick, cookie)
	}
}

// JoinGame ...
func (pj *OnlinePJs) JoinGame(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	cookie, err := r.Cookie(cookieName)
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
	pj.connect(username, value)
	u := store.NewUser()
	u.Nick = username
	info := &util.ResponseInfo{
		IsLogged: true,
		User:     u,
	}
	util.SendJSONToClient(w, info)
}

// JoinAnonymous ...
func (pj *OnlinePJs) JoinAnonymous(w http.ResponseWriter, r *http.Request) {
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
	token, err := generateRandomString(sessionLength)
	if err != nil {
		log.Print("Error Generating Random String")
		token = time.Now().String()
	}
	sessionID := nick + ":" + token
	expires := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Domain:   util.CheckModeForCookieDomain(), //"localhost",
		Path:     "/",
		HttpOnly: util.CheckModeForCookieHTTPOnly(), //false,
		Expires:  expires,
	}
	http.SetCookie(w, cookie)
	username := strings.Split(cookie.Value, ":")[0]
	value := strings.Split(cookie.Value, ":")[1]
	pj.connect(username, value)
	info := &util.ResponseInfo{
		IsLogged: true,
		User:     u,
	}
	util.SendJSONToClient(w, info)
}
