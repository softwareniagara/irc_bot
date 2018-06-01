package main

import (
	"database/sql"

	"github.com/whyrusleeping/hellabot"

	"github.com/softwareniagara/irc_bot/store"
)

func GreeterTrigger(s *store.Store, nick string) hbot.Trigger {
	return hbot.Trigger{
		Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
			return msg.Command == "JOIN" && msg.From != nick
		},
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			u, err := s.FindUserByNick(msg.From)
			if err != nil {
				if err == sql.ErrNoRows {
					ReplyTo(bot, msg, "hey")
					return true
				}
				ErrorReply(bot, msg, err)
				return true
			}
			if u.Greeting != "" {
				ReplyTo(bot, msg, u.Greeting)
			}
			return true
		},
	}
}
