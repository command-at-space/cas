package login

var configjson = []byte(`
{
	"auth": {
		"version": "0.0.1",
		"bcryptCost": 12,
		"sessionLength": 32,
		"cookieName": "playingCAS"
	},
	"test": {
		"test": "test"
	}
}
`)

// ActiveUsers ...
var ActiveUsers = map[string]string{
	"mainBot": "a",
}
var c configuration
var e myError

type user struct {
	name  string
	hash  string
	email string
	logo  string
}

type configuration struct {
	Auth struct {
		Version       string //`json:"version"`
		BCryptCost    int
		SessionLength int
		CookieName    string
	} //`json:"app"`
	Test struct {
		Test string
	}
	Mysql struct {
		Host     string `json:"host"`
		Port     int    `json:"port"`
		Db       string `json:"db"`
		Table1   string `json:"table1"`
		Table2   string `json:"table2"`
		User     string `json:"user"`
		Password string `json:"password"`
	} `json:"mysql"`
}

type myError struct {
	Text string `json:"Error,omitempty"`
}
