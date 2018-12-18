package login

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // mysql driver
)

type myDB struct{}

var db *sql.DB

var loginDB myDB

func (loginDB *myDB) initDB() {
	connPath := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", c.Mysql.User, c.Mysql.Password, c.Mysql.Host, c.Mysql.Port, c.Mysql.Db)
	//fmt.Println(connPath)
	var err error
	db, err = sql.Open("mysql", connPath)
	if err != nil {
		log.Printf("ERROR 1 DB %s\n", err)
	} else {
		log.Printf("INFO DB %s AUTH sql.Open() => OK\n", c.Mysql.Db)
	}
}

// AUTH METHODS

func (loginDB *myDB) insertNewAccount(u user) error {
	sql := fmt.Sprintf("INSERT INTO `%s` (username, hash, email, logo) values ('%s', '%s', '%s', '%s')", c.Mysql.Table1, u.name, u.hash, u.email, u.logo)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("ERROR 2 DB %s\n", err)
		e.Text = fmt.Sprintf("ERROR 2 DB %s\n", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("ERROR 3 DB %s\n", err)
		e.Text = fmt.Sprintf("User %s already exist", u.name)
		return err
	}
	return nil
}

func (loginDB *myDB) getAccount(name string) (u user, err error) {
	sql := fmt.Sprintf("SELECT * FROM `%s` WHERE username=?", c.Mysql.Table1)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("ERROR 4 DB %s\n", err)
		return u, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(name).Scan(&u.name, &u.hash, &u.email, &u.logo)
	if err != nil {
		log.Printf("ERROR 5 DB %s\n", err)
		return u, err
	}
	return u, nil
}

// SESSION METHODS

func (loginDB *myDB) saveSession(username, sessionID string) error {
	sql := fmt.Sprintf("INSERT INTO `%s` (username, sessionID) values ('%s', '%s')ON DUPLICATE KEY UPDATE sessionID = '%s'", c.Mysql.Table2, username, sessionID, sessionID)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("ERROR 6 DB %s\n", err)
		e.Text = fmt.Sprintf("ERROR 6 DB %s\n", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("ERROR 7 DB %s\n", err)
		e.Text = fmt.Sprintf("Session from user %s already exist", username)
		return err
	}
	return nil
}

func (loginDB *myDB) deleteSession(username string) error {
	sql := fmt.Sprintf("DELETE FROM `%s` WHERE username = '%s'", c.Mysql.Table2, username)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("ERROR 8 DB %s\n", err)
		e.Text = fmt.Sprintf("ERROR 8 DB %s\n", err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec()
	if err != nil {
		log.Printf("ERROR 9 DB %s\n", err)
		e.Text = fmt.Sprintf("Cant delete session from user %s", username)
		return err
	}
	return nil
}

func (loginDB *myDB) loadAllSessions() error {
	sql := fmt.Sprintf("SELECT * FROM `%s`", c.Mysql.Table2)
	stmt, err := db.Prepare(sql)
	if err != nil {
		log.Printf("ERROR 10 DB %s\n", err)
		e.Text = fmt.Sprintf("ERROR 10 DB %s\n", err)
		return err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		log.Printf("ERROR 11 DB %s\n", err)
		e.Text = fmt.Sprintf("ERROR 11 DB %s\n", err)
		return err
	}
	defer rows.Close()
	var username, sessionID string
	for rows.Next() {
		rows.Scan(&username, &sessionID)
		ActiveUsers[username] = sessionID
	}
	if err = rows.Err(); err != nil {
		log.Printf("ERROR 12 DB %s\n", err)
		e.Text = fmt.Sprintf("ERROR 12 DB %s\n", err)
		return err
	}
	return nil
}

func (loginDB *myDB) existSession() (err error) {
	return nil
}

/*
func dbSearchCookie(user, sessionID string) bool {
	var u, s string
	//fmt.Println(`User, SESSIONID`, user, sessionID)
	row := db.QueryRow("SELECT * FROM book.sessions WHERE Username = ? AND SessionID = ?", user, user+":"+sessionID)
	err := row.Scan(&u, &s)
	if err != nil {
		if err == sql.ErrNoRows {
			fmt.Println("No Records Found")
		} else {
			log.Fatal(err)
		}
	} else {
		return true // exists vote
	}
	return false
}
*/
