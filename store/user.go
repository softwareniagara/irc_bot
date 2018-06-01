package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
)

type Role int

const (
	RoleInvalid Role = iota
	RoleAdmin
	RoleRegular
	RoleBanned
	RoleNone
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleRegular:
		return "regular"
	case RoleBanned:
		return "banned"
	case RoleNone:
		return "none"
	default:
		return "invalid"
	}
}

func (r *Role) Set(s string) error {
	*r = RoleFromString(s)
	if *r == RoleInvalid {
		return fmt.Errorf("invalid role: %s", s)
	}
	return nil
}

func RoleFromString(s string) Role {
	switch s {
	case "admin":
		return RoleAdmin
	case "regular":
		return RoleRegular
	case "banned":
		return RoleBanned
	case "none":
		return RoleNone
	default:
		return RoleInvalid
	}
}

type User struct {
	RowID      int64
	Nick       string
	Role       Role
	Active     bool
	LastActive time.Time
}

func (u User) String() string {
	var active string
	switch {
	case u.Active:
		active = "now"
	case u.LastActive.IsZero():
		active = "unknown"
	default:
		active = fmt.Sprintf("%s ago", time.Since(u.LastActive))
	}
	return fmt.Sprintf("nick=%s role=%s active=%s", u.Nick, u.Role, active)
}

const createUserSQL = `
	CREATE TABLE IF NOT EXISTS users (
		nick        TEXT UNIQUE,
		role        INTEGER,
		active      BOOLEAN,
		last_active TIMESTAMP
	)
`

const insertUserSQL = `
	INSERT INTO users (
		nick,
		role,
		active,
		last_active
	) VALUES (?, ?, ?, ?)
`

func (s *Store) NotBanned(nick string) error {
	u, err := s.FindUserByNick(nick)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	if u.Role == RoleBanned {
		return errors.New("banned users can't do that")
	}
	return nil
}

func (s *Store) Authorized(nick string, roles ...Role) error {
	u, err := s.FindUserByNick(nick)
	if err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("%s is not a user", nick)
		}
		return err
	}
	for _, r := range roles {
		if u.Role == r {
			return nil
		}
	}
	return fmt.Errorf("%s users can't do that", u.Role)
}

func (s *Store) InsertUser(u *User) error {
	_, err := s.db.Exec(insertUserSQL, u.Nick, u.Role, u.Active, u.LastActive)
	return err
}

const updateUserSQL = `
	UPDATE users
	SET nick = ?,
	    role = ?,
			active = ?,
			last_active = ?
  WHERE ROWID = ?
`

func (s *Store) UpdateUser(u *User) error {
	_, err := s.db.Exec(updateUserSQL, u.Nick, u.Role, u.Active, u.LastActive, u.RowID)
	return err
}

const findUserByNickSQL = `
	SELECT ROWID, nick, role, active, last_active
	FROM users
	WHERE nick = ?
`

func (s *Store) FindUserByNick(nick string) (*User, error) {
	var u User
	if err := s.db.QueryRow(findUserByNickSQL, nick).Scan(
		&u.RowID,
		&u.Nick,
		&u.Role,
		&u.Active,
		&u.LastActive,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

const findUserSQL = `
	SELECT ROWID, nick, role, active, last_active
	FROM users
	WHERE ROWID = ?
`

func (s *Store) FindUser(rowID int64) (*User, error) {
	var u User
	if err := s.db.QueryRow(findUserSQL, rowID).Scan(
		&u.RowID,
		&u.Nick,
		&u.Role,
		&u.Active,
		&u.LastActive,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

const removeUserSQL = `
	DELETE FROM users
	WHERE ROWID = ?
`

func (s *Store) RemoveUser(rowID int64) error {
	_, err := s.db.Exec(removeUserSQL, rowID)
	return err
}
