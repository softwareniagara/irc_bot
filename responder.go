package main

import (
	"fmt"

	"github.com/whyrusleeping/hellabot"
)

func Responder(nick string) hbot.Trigger {
	return hbot.Trigger{
		Condition: HasPrefix(nick + ":"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			bot.Reply(msg, fmt.Sprintf("%s: never speak to me or my son again", msg.From))
			return true
		},
	}
}
