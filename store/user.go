package store

import (
	"fmt"
)

type Role int

const (
	RoleInvalid Role = iota
	RoleAdmin
	RoleUser
	RoleIdiot
	RoleBanned
)

func (r Role) String() string {
	switch r {
	case RoleAdmin:
		return "admin"
	case RoleUser:
		return "user"
	case RoleIdiot:
		return "idiot"
	case RoleBanned:
		return "banned"
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
	case "user":
		return RoleUser
	case "idiot":
		return RoleIdiot
	case "banned":
		return RoleBanned
	default:
		return RoleInvalid
	}
}

type User struct {
	RowID int64
	Nick  string
	Role  Role
}

const createUserSQL = `
	CREATE TABLE IF NOT EXISTS users (
		nick TEXT UNIQUE,
		role INTEGER
	)
`

const insertUserSQL = `
	INSERT INTO users (
		nick,
		role
	) VALUES (?, ?)
`

func (s *Store) InsertUser(u *User) error {
	_, err := s.db.Exec(insertUserSQL, u.Nick, u.Role)
	return err
}

const updateUserSQL = `
	UPDATE users
	SET nick = ?,
	    role = ?
  WHERE ROWID = ?
`

func (s *Store) UpdateUser(u *User) error {
	_, err := s.db.Exec(updateUserSQL, u.Nick, u.Role, u.RowID)
	return err
}

const findUserByNickSQL = `
	SELECT ROWID, nick, role
	FROM users
	WHERE nick = ?
`

func (s *Store) FindUserByNick(nick string) (*User, error) {
	var u User
	if err := s.db.QueryRow(findUserByNickSQL, nick).Scan(
		&u.RowID,
		&u.Nick,
		&u.Role,
	); err != nil {
		return nil, err
	}
	return &u, nil
}

const findUserSQL = `
	SELECT ROWID, nick, role
	FROM users
	WHERE ROWID = ?
`

func (s *Store) FindUser(rowID int64) (*User, error) {
	var u User
	if err := s.db.QueryRow(findUserSQL, rowID).Scan(
		&u.RowID,
		&u.Nick,
		&u.Role,
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
