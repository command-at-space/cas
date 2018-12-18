package login

import (
	"log"
	"net/http"
	"strings"
	"time"
)

//type session struct{}

// IsLogged ...
func IsLogged(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		r.ParseForm()
		//valid := r.Form.Get("test")
		cookie, err := r.Cookie(c.Auth.CookieName)
		if err != nil {
			log.Print("No cookie")
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

func setSessionCookie(w http.ResponseWriter, username string) string {
	token, err := generateRandomString(c.Auth.SessionLength)
	if err != nil {
		log.Print("Error Generating Random String")
		token = time.Now().String()
	}
	sessionID := username + ":" + token
	//expire := time.Now().AddDate(0, 0, 1)
	cookie := &http.Cookie{
		Name:   c.Auth.CookieName,
		Value:  sessionID,
		Domain: "localhost",
		Path:   "/",
		//HttpOnly: false,
		//Expires: expire,
	}
	http.SetCookie(w, cookie)
	return cookie.Value
}
