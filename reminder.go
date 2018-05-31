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

func (rd *ReminderDB) action(bot *hbot.Bot, msg *hbot.Message) error {
	if err := rd.s.Authorized(msg.From, store.RoleUser, store.RoleUser); err != nil {
		return err
	}
	var wait time.Duration
	fset := flag.NewFlagSet("remindme", flag.ContinueOnError)
	fset.DurationVar(&wait, "wait", 0, "how long to wait before reminding")
	if err := ParseFlags(msg, fset); err != nil {
		return err
	}
	return rd.s.InsertReminder(&store.Reminder{
		Time:    time.Now().Add(wait),
		Created: time.Now(),
		Nick:    msg.From,
		Msg:     strings.Join(fset.Args(), " "),
	})
}

func (rd *ReminderDB) Trigger() hbot.Trigger {
	return hbot.Trigger{
		Condition: HasPrefix("!remindme"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			if err := rd.action(bot, msg); err != nil {
				ErrorReply(bot, msg, err)
			} else {
				ReplyTo(bot, msg, "ok")
			}
			return true
		},
	}
}
