/* */

package login

import (
	"fmt"
	"net/http"
	"regexp"
	"strings"

	store "casServer/login/store"
)

func validateNewUserData(u *store.User) (ok string) {
	ok = ""
	if len(u.NickID) < 4 || len(u.NickID) > 16 {
		ok += fmt.Sprintf("Name between 4-16 characters\n")
	}
	if !checkValidCharacters(u.NickID) {
		ok += fmt.Sprintf("Name only can contain numbers and letters\n")
	}
	if len(u.PassHashed) < 8 || len(u.PassHashed) > 20 {
		ok += fmt.Sprintf("Password between 8-20 any characters\n")
	}
	rxMail := regexp.MustCompile(`^\S+@\S+\.\S+$`)
	if u.Email != "" && !rxMail.MatchString(u.Email) {
		ok += fmt.Sprintf("Please use a valid email address\n")
	}
	if len(u.Email) > 60 {
		ok += fmt.Sprintf("Email adress maximun length 60 characters\n")
	}
	rxLogo := regexp.MustCompile(`\.(jpeg|jpg|gif|png)+$`)
	if u.Logo != "" && !rxLogo.MatchString(u.Logo) {
		ok += fmt.Sprintf("Logo url is not a valid jpeg,jpg,png or gif file\n")
	}
	if u.Logo != "" && !doesReallyLogoURLExists(u.Logo) {
		ok += fmt.Sprintf("Logo url doesnt exist\n")
	}
	if len(u.Logo) > 90 {
		ok += fmt.Sprintf("Logo url maximun length 90 characters\n")
	}
	if len(u.SecretQuest) > 90 {
		ok += fmt.Sprintf("Secret question maximun length 90 characters\n")
	}
	if len(u.SecretHashed) > 20 {
		ok += fmt.Sprintf("Secret response maximun length 20 characters\n")
	}
	if u.SecretQuest != "" && u.SecretHashed == "" {
		ok += fmt.Sprintf("Fill secret response\n")
	}
	if u.SecretQuest == "" && u.SecretHashed != "" {
		ok += fmt.Sprintf("Fill secret question\n")
	}
	return ok
}

func validateNewAnonymousData(nick string) (ok string) {
	nick = nick[4:len(nick)] // remove tmp_ prefix
	ok = ""
	if len(nick) < 4 || len(nick) > 8 {
		ok += fmt.Sprintf("Name between 4-8 characters\n")
	}
	if !checkValidCharacters(nick) {
		ok += fmt.Sprintf("Name only can contain numbers and letters\n")
	}
	return ok
}

func checkValidCharacters(str string) bool {
	str = strings.ToLower(str)
	for _, char := range str {
		if !strings.Contains(validChars, string(char)) {
			return false
		}
	}
	return true
}

func doesReallyLogoURLExists(str string) bool {
	resp, err := http.Head(str)
	if err != nil {
		return false
	}
	if resp.StatusCode != http.StatusOK {
		return false
	}
	return true
}
