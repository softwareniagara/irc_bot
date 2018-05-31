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

func ReminderNotify(bot *hbot.Bot, s *store.Store, channel string) error {
	r, err := s.NextReminder(time.Now())
	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	msg := fmt.Sprintf("%s: reminder: %s", r.Nick, r.Msg)
	bot.Notice(channel, msg)
	return s.RemoveReminder(r.RowID)
}

func ReminderNotifyLoop(bot *hbot.Bot, s *store.Store, channel string) {
	for {
		if err := ReminderNotify(bot, s, channel); err != nil {
			log.Println(err)
		}
		time.Sleep(10 * time.Second)
	}
}

func ReminderTrigger(s *store.Store) hbot.Trigger {
	return hbot.Trigger{
		Condition: HasCommand("!remindme"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			if err := s.Authorized(msg.From, store.RoleUser, store.RoleUser); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}

			var wait time.Duration
			fset := flag.NewFlagSet("remindme", flag.ContinueOnError)
			fset.DurationVar(&wait, "wait", 0, "how long to wait before reminding")
			if err := ParseFlags(msg, fset); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}

			err := s.InsertReminder(&store.Reminder{
				Time:    time.Now().Add(wait),
				Created: time.Now(),
				Nick:    msg.From,
				Msg:     strings.Join(fset.Args(), " "),
			})
			if err != nil {
				ErrorReply(bot, msg, err)
				return true
			}

			ReplyTo(bot, msg, "ok")
			return true
		},
	}
}
