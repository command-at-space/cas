/* */

package login

const (
	sessionLength int    = 32
	bcryptCost    int    = 12
	cookieName    string = "alphaCAS"
	validChars    string = "abcdefghijklmnopqrstuvwxyz0123456789"
)

// ActiveUsers ...
var ActiveUsers = map[string]string{
	"mainBot": "a",
}
