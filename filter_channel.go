package main

import (
	"fmt"
	"github.com/whyrusleeping/hellabot"
)

func FilterChannel(channel string) hbot.Trigger {
	return hbot.Trigger{
		Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
			return msg.To != channel
		},
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			if msg.Command == "PRIVMSG" {
				bot.Reply(msg, fmt.Sprintf("I only respond in %s", channel))
			}
			return true
		},
	}
}
