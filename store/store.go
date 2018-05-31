package store

import (
	"database/sql"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type Reminder struct {
	RowID   int64
	Time    time.Time
	Created time.Time
	Nick    string
	Msg     string
}

type Store struct {
	db *sql.DB
}

func Open(filename string) (*Store, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	s := &Store{db}
	if err := s.CreateReminder(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return s, nil
}

func (s *Store) Close() error {
	return s.db.Close()
}

const createReminderSQL = `
	CREATE TABLE IF NOT EXISTS reminders (
		time    TIMESTAMP,
		created TIMESTAMP,
		nick    TEXT,
		msg     TEXT
	)
`

func (s *Store) CreateReminder() error {
	_, err := s.db.Exec(createReminderSQL)
	return err
}

const insertReminderSQL = `
	INSERT INTO reminders (
		time,
		created,
		nick,
		msg
	) VALUES (?, ?, ?, ?)
`

func (s *Store) InsertReminder(r *Reminder) error {
	_, err := s.db.Exec(insertReminderSQL, r.Time, r.Created, r.Nick, r.Msg)
	return err
}

const nextReminderSQL = `
	SELECT ROWID, time, created, nick, msg
	FROM reminders
	WHERE time <= ?
	ORDER BY time
	LIMIT 1
`

func (s *Store) NextReminder(now time.Time) (*Reminder, error) {
	var r Reminder
	if err := s.db.QueryRow(nextReminderSQL, now).Scan(
		&r.RowID,
		&r.Created,
		&r.Time,
		&r.Nick,
		&r.Msg,
	); err != nil {
		return nil, err
	}
	return &r, nil
}

const removeReminderSQL = `
	DELETE FROM reminders
	WHERE ROWID = ?
`

func (s *Store) RemoveReminder(rowID int64) error {
	_, err := s.db.Exec(removeReminderSQL, rowID)
	return err
}
