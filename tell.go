package main

import (
	"strings"

	"github.com/whyrusleeping/hellabot"

	"github.com/softwareniagara/irc_bot/store"
)

func TellTrigger(s *store.Store) hbot.Trigger {
	return hbot.Trigger{
		Condition: HasCommand("!tell"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			if err := s.NotBanned(msg.From); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}
			content := strings.TrimPrefix(msg.Content, "!tell")
			fields := strings.Fields(content)
			if len(fields) == 0 {
				ReplyTo(bot, msg, "Say what?")
				return true
			}
			text := fields[0] + ": " + strings.Join(fields[1:], " ")
			bot.Reply(msg, text)
			return true
		},
	}
}
