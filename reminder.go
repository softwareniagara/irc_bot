package main

import (
	"database/sql"
	"time"
)

type Reminder struct {
	RowID   int64
	Time    time.Time
	Created time.Time
	Nick    string
	Msg     string
}

type ReminderDB struct {
	db *sql.DB
}

const createSQL = `
	CREATE TABLE reminders (
		time    TIMESTAMP,
		created TIMESTAMP,
		nick    TEXT,
		msg     TEXT
	)
`

func (rd *ReminderDB) Create() error {
	_, err := rd.db.Exec(createSQL)
	return err
}

const insertSQL = `
	INSERT INTO reminders (
		time,
		created,
		nick,
		msg
	) VALUES (?, ?, ?, ?)
`

func (rd *ReminderDB) Insert(r *Reminder) error {
	_, err := rd.db.Exec(insertSQL, r.Time, r.Created, r.Nick, r.Msg)
	return err
}

const nextSQL = `
	SELECT ROWID, time, created, nick, msg
	FROM reminders
	WHERE time <= ?
	ORDER BY time
	LIMIT 1
`

func (rd *ReminderDB) Next(now time.Time) (*Reminder, error) {
	var r Reminder
	if err := rd.db.QueryRow(nextSQL, now).Scan(
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

const removeSQL = `
	DELETE FROM reminders
	WHERE ROWID = ?
`

func (rd *ReminderDB) Remove(rowID int64) error {
	_, err := rd.db.Exec(removeSQL, rowID)
	return err
}
