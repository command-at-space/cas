/* */

package store

import (
	"database/sql"
	"fmt"

	// mysql driver
	_ "github.com/go-sql-driver/mysql"
)

type configDB struct {
	DatabaseType string `json:"databaseType"`
	Host         string `json:"host"`
	Port         int    `json:"port"`
	DB           string `json:"db"`
	User         string `json:"user"`
	Password     string `json:"password"`
}

// DB ...
type DB struct {
	*sql.DB
}

// NewDB ...
func NewDB(mode string) (*DB, error) {
	var c configDB
	loadConfigJSON(&c)
	setDBConnConfig(mode, &c)
	connPath := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=true",
		c.User,
		c.Password,
		c.Host,
		c.Port,
		c.DB,
	)
	//fmt.Println("CONNPATH => ", connPath)
	db, err := sql.Open(c.DatabaseType, connPath)
	if err != nil {
		return nil, err
	}
	/*err = db.Ping()
	if err != nil {
		return nil, err
	}*/
	return &DB{db}, nil
}
