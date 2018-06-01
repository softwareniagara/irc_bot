package main

import (
	"github.com/whyrusleeping/hellabot"
)

func HelpTrigger() hbot.Trigger {
	return hbot.Trigger{
		Condition: HasCommand("!help"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			ReplyTo(bot, msg, "you're beyond help")
			return true
		},
	}
}
