/* */

package store

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

// Session ...
type Session struct {
	NickID    string    `json:"-"`
	SessionID string    `json:"-"`
	Expires   time.Time `json:"-"`
}

// NewSession ...
func NewSession() *Session {
	return new(Session)
	//return &Session{}
}

// SaveSession ...
func (loginDB *DB) SaveSession(s *Session) error {
	query := fmt.Sprintf(`
		INSERT 
		INTO %s 
		(
			nickID, sessionID, expires
		) 
		values (?, ?, ?)
		ON DUPLICATE KEY UPDATE sessionID=? , expires= ?
		`,
		sessionsTable)
	stmt, err := loginDB.Prepare(query)
	if err != nil {
		log.Printf("ERROR 1 DB SESSIONS inserting %s -> %s\n", s.NickID, err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		s.NickID,
		s.SessionID,
		s.Expires,
		s.SessionID,
		s.Expires,
	)
	if err != nil {
		log.Printf("ERROR 2 DB SESSIONS inserting %s -> %s\n", s.NickID, err)
		return err
	}
	return nil
}

// UserSession ...
func (loginDB *DB) UserSession(usernameID string) (*Session, error) {
	s := NewSession()
	query := fmt.Sprintf(`
		SELECT
			*
		FROM %s
		WHERE nickID=?
		`,
		sessionsTable,
	)
	stmt, err := loginDB.Prepare(query)
	if err != nil {
		log.Printf("ERROR 3 DB SESSIONS %s\n", err)
		return s, err
	}
	defer stmt.Close()
	err = stmt.QueryRow(usernameID).Scan(
		&s.NickID,
		&s.SessionID,
		&s.Expires,
	)
	if err != nil {
		if err == sql.ErrNoRows { // no result
			//log.Println("NO RESULT", s)
			return s, nil
		}
		log.Printf("ERROR 4 DB SESSIONS %s\n", err)
		return nil, err
	}
	return s, nil
}

// DeleteSession ...
func (loginDB *DB) DeleteSession(usernameID string) error {
	query := fmt.Sprintf(`
		DELETE 
		FROM %s 
		WHERE nickID = ?
		`,
		sessionsTable)
	stmt, err := loginDB.Prepare(query)
	if err != nil {
		log.Printf("ERROR 5 DB SESSIONS deleting %s -> %s\n", usernameID, err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(
		usernameID,
	)
	if err != nil {
		log.Printf("ERROR 6 DB SESSIONS deleting %s -> %s\n", usernameID, err)
		return err
	}
	return nil
}
