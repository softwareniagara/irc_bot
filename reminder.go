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

	"github.com/softwareniagara/irc_bot/store"
)

type ReminderDB struct {
	s       *store.Store
	channel string
}

func NewReminderDB(s *store.Store, channel string) *ReminderDB {
	return &ReminderDB{s, channel}
}

func (rd *ReminderDB) poll(bot *hbot.Bot) error {
	r, err := rd.s.NextReminder(time.Now())
	if err != nil {
		return err
	}
	msg := fmt.Sprintf("%s: reminder: %s", r.Nick, r.Msg)
	bot.Notice(rd.channel, msg)
	return rd.s.RemoveReminder(r.RowID)
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

func (rd *ReminderDB) action(bot *hbot.Bot, msg *hbot.Message) bool {
	var dur time.Duration
	fset := flag.NewFlagSet("remindme", flag.ContinueOnError)
	fset.DurationVar(&dur, "d", 0, "how long to wait before reminding")
	if err := ParseFlags(msg, fset); err != nil {
		ErrorReply(bot, msg, err)
		return true
	}
	if err := rd.s.InsertReminder(&store.Reminder{
		Time:    time.Now().Add(dur),
		Created: time.Now(),
		Nick:    msg.From,
		Msg:     strings.Join(fset.Args(), " "),
	}); err != nil {
		ErrorReply(bot, msg, err)
		return true
	}
	ReplyTo(bot, msg, "ok")
	return true
}

func (rd *ReminderDB) Trigger() hbot.Trigger {
	return hbot.Trigger{
		Condition: HasPrefix("!remindme"),
		Action:    rd.action,
	}
}
