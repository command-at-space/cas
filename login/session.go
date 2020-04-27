/* */

package login

import (
	util "casServer/utils"
	"fmt"
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
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			//log.Print("No cookie")
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		user := strings.Split(cookie.Value, ":")[0]
		sessionID := strings.Split(cookie.Value, ":")[1]
		value, ok := ActiveUsers[user]
		if !ok || value != sessionID {
			log.Printf("No Active User %s with session %s\n", user, sessionID)
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
		cookie, err := r.Cookie(cookieName)
		if err != nil {
			//log.Print("No cookie")
			next.ServeHTTP(w, r)
			return
		}
		username := r.FormValue("user")
		//fmt.Println("HAY COOKIE ---")
		user := strings.Split(cookie.Value, ":")[0]
		sessionID := strings.Split(cookie.Value, ":")[1]
		value, ok := ActiveUsers[user]
		if ok && value == sessionID {
			if user == username {
				e := new(util.RequestError)
				e.Message = fmt.Sprintf(`User %s is already logged`, user)
				e.StatusCode = 401
				util.SendErrorToClient(w, e)
				return
			}
		}
		next.ServeHTTP(w, r)
	}
}

func setSessionCookie(w http.ResponseWriter, username string) string {
	token, err := generateRandomString(sessionLength)
	if err != nil {
		log.Print("Error Generating Random String")
		token = time.Now().String()
	}
	sessionID := username + ":" + token
	expire := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:     cookieName,
		Value:    sessionID,
		Domain:   checkModeForCookieDomain(),
		Path:     "/",
		HttpOnly: checkModeForCookieHTTPOnly(),
		Expires:  expire,
	}
	http.SetCookie(w, cookie)
	return cookie.Value
}

func deleteSessionCookie(w http.ResponseWriter) {
	cookie := &http.Cookie{
		Name:   cookieName,
		Value:  "",
		Path:   "/",
		MaxAge: -1,
	}
	http.SetCookie(w, cookie)
}
