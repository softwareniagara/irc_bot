package store

import (
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

const createReminderSQL = `
	CREATE TABLE IF NOT EXISTS reminders (
		time    TIMESTAMP,
		created TIMESTAMP,
		nick    TEXT,
		msg     TEXT
	)
`

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
