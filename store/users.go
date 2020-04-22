/* */

package store

import (
	"database/sql"
	"fmt"
	"log"
)

// User ...
type User struct {
	Name  string `json:"name"`
	Hash  string `json:"hash"`
	Email string `json:"email"`
	Logo  string `json:"logo"`
}

// AccountList ...
func (db *DB) AccountList() ([]*User, error) {
	query := fmt.Sprintf(`
		SELECT 
			name, 
			hash,
			COALESCE(email, ''),
			COALESCE(logo, '') 
		FROM %s
		`,
		usersTable,
	)
	stmt, err := db.Prepare(query)
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
			&user.Name,
			&user.Hash,
			&user.Email,
			&user.Logo,
		)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

// Account ...
func (db *DB) Account(name string) (*User, error) {
	u := new(User)
	query := fmt.Sprintf(`
		SELECT 
			name, 
			hash,
			COALESCE(email, ''),
			COALESCE(logo, '') 
		FROM %s
		WHERE name=?
		`,
		usersTable,
	)
	//fmt.Println(query)
	stmt, err := db.Prepare(query)
	if err != nil {
		log.Printf("ERROR 1 DB %s\n", err)
		return u, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(name).Scan(
		&u.Name,
		&u.Hash,
		&u.Email,
		&u.Logo,
	)
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
func (db *DB) NewAccount(u *User) error {
	//fmt.Println("USER => ", u)
	query := fmt.Sprintf(`
		INSERT 
		INTO %s 
		(name, hash, email, logo)
		values (?, ?, ?, ?)
		`,
		usersTable,
	)
	stmt, err := db.Prepare(query)
	//fmt.Println(query)
	if err != nil {
		log.Printf("ERROR 3 DB %s\n", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(u.Name, u.Hash, u.Email, u.Logo)
	if err != nil { // for exampke duplicate entry user -> jolav and Jolav
		log.Printf("ERROR 4 DB %s\n", err)
		return err
	}
	//fmt.Println("Ok => ", ok)
	return nil
}

// NewUser ...
func NewUser() *User {
	return new(User)
	//return &User{}
}
