package store

type Role int

const (
	RoleAdmin Role = iota
	RoleUser
	RoleBanned
	RoleIdiot
)

type User struct {
	RowID int64
	Nick  string
	Role  Role
}

const createUserSQL = `
	CREATE TABLE IF NOT EXISTS users (
		nick TEXT,
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
	_, err := s.db.Exec(removeUserSQL)
	return err
}
