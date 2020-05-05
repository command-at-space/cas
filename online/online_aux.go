/* */

package online

import (
	"fmt"

	util "casServer/utils"
)

func validateNewAnonymousData(nick string) (ok string) {
	nick = nick[4:len(nick)] // remove tmp_ prefix
	ok = ""
	if len(nick) < 4 || len(nick) > 8 {
		ok += fmt.Sprintf("Name between 4-8 characters\n")
	}
	if !util.CheckValidCharacters(nick) {
		ok += fmt.Sprintf("Name only can contain numbers and letters\n")
	}
	return ok
}
