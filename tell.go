package main

import (
	"strings"

	"github.com/whyrusleeping/hellabot"
)

var TellTrigger = hbot.Trigger{
	Condition: HasCommand("!tell"),
	Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
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
