/* */

package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// User ...
type User struct {
	NickID       string    `json:"-"`
	Nick         string    `json:"nick"`
	PassHashed   string    `json:"-"`
	Email        string    `json:"-"`
	Verified     int       `json:"verified"`
	Logo         string    `json:"logo"`
	SecretQuest  string    `json:"-"`
	SecretHashed string    `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	LastSeen     time.Time `json:"-"`
	Online       int       `json:"online"`
}

// NewUser ...
func NewUser() *User {
	return new(User)
	//return &User{}
}

// Account ...
func (loginDB *DB) Account(name string) (*User, error) {
	//fmt.Println("SEARCHING ...", name)
	u := new(User)
	query := fmt.Sprintf(`
		SELECT
			*
		FROM %s
		WHERE nickID=?
		`,
		usersTable,
	)
	//fmt.Println(query)
	stmt, err := loginDB.Prepare(query)
	if err != nil {
		log.Printf("ERROR 1 DB %s\n", err)
		return u, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(name).Scan(
		&u.NickID,
		&u.Nick,
		&u.PassHashed,
		&u.Email,
		&u.Verified,
		&u.Logo,
		&u.SecretQuest,
		&u.SecretHashed,
		&u.CreatedAt,
		&u.LastSeen,
		&u.Online)
	if err != nil {
		if err == sql.ErrNoRows { // no result
			//log.Println("NO RESULT", u)
			return u, nil
		}
		log.Printf("ERROR 2 DB %s\n", err)
		return nil, err
	}
	return u, nil
}

// NewAccount ..
func (loginDB *DB) NewAccount(u *User) error {
	//fmt.Println("USER => ", u)
	query := fmt.Sprintf(`
		INSERT
		INTO %s
		(
			nickID, nick, passHashed, email, verified, logo, 
			secretQuest, secretHashed, createdAt, lastSeen, online
		)
		values (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
		`,
		usersTable,
	)
	stmt, err := loginDB.Prepare(query)
	//fmt.Println(query)
	if err != nil {
		log.Printf("ERROR 3 DB inserting %s -> %s\n", u.Nick, err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		u.NickID,
		u.Nick,
		u.PassHashed,
		u.Email,
		u.Verified,
		u.Logo,
		u.SecretQuest,
		u.SecretHashed,
		u.CreatedAt,
		u.LastSeen,
		u.Online,
	)
	if err != nil { // for exampke duplicate entry user -> jolav and Jolav
		log.Printf("ERROR 4 DB inserting %s -> %s\n", u.Nick, err)
		return err
	}
	return nil
}

// AccountList ...
func (loginDB *DB) AccountList() ([]*User, error) {
	query := fmt.Sprintf(`
		SELECT
			nickID
		FROM %s
		`,
		usersTable,
	)
	stmt, err := loginDB.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	rows, err := stmt.Query()
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	users := make([]*User, 0)
	for rows.Next() {
		user := new(User)
		err := rows.Scan(
			&user.NickID,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}
