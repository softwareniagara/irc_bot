package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/whyrusleeping/hellabot"
)

type Reminder struct {
	RowID   int64
	Time    time.Time
	Created time.Time
	Nick    string
	Msg     string
}

type ReminderDB struct {
	db      *sql.DB
	channel string
}

func OpenReminderDB(filename, channel string) (*ReminderDB, error) {
	db, err := sql.Open("sqlite3", filename)
	if err != nil {
		return nil, err
	}
	rd := &ReminderDB{db, channel}
	if err := rd.create(); err != nil {
		_ = db.Close()
		return nil, err
	}
	return rd, nil
}

func (rd *ReminderDB) Close() error {
	return rd.db.Close()
}

func (rd *ReminderDB) poll(bot *hbot.Bot) error {
	r, err := rd.next(time.Now())
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s: reminder: %s", r.Nick, r.Msg)
	bot.Notice(rd.channel, msg)
	return rd.remove(r.RowID)
}

func (rd *ReminderDB) Run(bot *hbot.Bot) {
	for {
		if err := rd.poll(bot); err != nil {
			if err != sql.ErrNoRows {
				log.Println(err)
			}
			time.Sleep(10 * time.Second)
		}
	}
}

func (rd *ReminderDB) condition(bot *hbot.Bot, msg *hbot.Message) bool {
	return strings.HasPrefix(msg.Content, "!remindme")
}

func (rd *ReminderDB) action(bot *hbot.Bot, msg *hbot.Message) bool {
	var dur time.Duration
	fset := flag.NewFlagSet("remindme", flag.ContinueOnError)
	fset.DurationVar(&dur, "d", 0, "how long to wait before reminding")
	if err := ParseFlags(msg, fset); err != nil {
		MultiLineReply(bot, msg, err.Error())
		return true
	}
	if err := rd.insert(&Reminder{
		Time:    time.Now().Add(dur),
		Created: time.Now(),
		Nick:    msg.From,
		Msg:     strings.Join(fset.Args(), " "),
	}); err != nil {
		MultiLineReply(bot, msg, err.Error())
		return true
	}
	bot.Reply(msg, fmt.Sprintf("%s: ok", msg.From))
	return true
}

func (rd *ReminderDB) Trigger() hbot.Trigger {
	return hbot.Trigger{
		Condition: rd.condition,
		Action:    rd.action,
	}
}

const createSQL = `
	CREATE TABLE IF NOT EXISTS reminders (
		time    TIMESTAMP,
		created TIMESTAMP,
		nick    TEXT,
		msg     TEXT
	)
`

func (rd *ReminderDB) create() error {
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

func (rd *ReminderDB) insert(r *Reminder) error {
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

func (rd *ReminderDB) next(now time.Time) (*Reminder, error) {
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

func (rd *ReminderDB) remove(rowID int64) error {
	_, err := rd.db.Exec(removeSQL, rowID)
	return err
}
