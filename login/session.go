/* */

package login

import (
	store "casServer/login/store"
	util "casServer/utils"
	"log"
	"net/http"
	"strings"
	"time"
)

// IsLogged ...
func IsLogged(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" && r.URL.String() != "/secret" {
			util.BadRequest(w, r)
			return
		}
		r.ParseForm()
		cookie, err := r.Cookie(CookieName)
		if err != nil { // No cookie
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		username := strings.Split(cookie.Value, ":")[0]
		usernameID := strings.ToLower(username)
		s, err := loginDB.UserSession(usernameID)
		if err != nil {
			message := "We are experiencing problems... Please try again later"
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
		if s.SessionID != cookie.Value {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	}
}

// IsNotLogged ...
func IsNotLogged(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			util.BadRequest(w, r)
			return
		}
		r.ParseForm()
		cookie, err := r.Cookie(CookieName)
		if err != nil { // No cookie
			next.ServeHTTP(w, r)
			return
		}
		username := r.FormValue("user")
		usernameID := strings.ToLower(username)
		s, err := loginDB.UserSession(usernameID)
		if err != nil {
			message := "We are experiencing problems... Please try again later"
			http.Error(w, message, http.StatusInternalServerError)
			return
		}
		if s.SessionID == cookie.Value { // User already logged
			message := "User is already logged"
			http.Error(w, message, http.StatusConflict)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func setSessionCookie(w http.ResponseWriter, s *store.Session, nick string) {
	token, err := GenerateRandomString(SessionLength)
	if err != nil {
		log.Print("Error Generating Random String")
		token = time.Now().String()
	}
	s.SessionID = nick + ":" + token
	s.Expires = time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:     CookieName,
		Value:    s.SessionID,
		Domain:   util.CheckModeForCookieDomain(),
		Path:     "/",
		HttpOnly: util.CheckModeForCookieHTTPOnly(),
		Expires:  s.Expires,
	}
	http.SetCookie(w, cookie)
}

func deleteSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   CookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
