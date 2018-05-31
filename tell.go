package main

import (
	"fmt"
	"strings"

	"github.com/whyrusleeping/hellabot"
)

var TellTrigger = hbot.Trigger{
	Condition: HasPrefix("!tell"),
	Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
		content := strings.TrimPrefix(msg.Content, "!tell")
		fields := strings.Fields(content)
		if len(fields) == 0 {
			bot.Reply(msg, fmt.Sprintf("%s: Say what?", msg.From))
			return true
		}
		text := fields[0] + ": " + strings.Join(fields[1:], " ")
		bot.Reply(msg, text)
		return true
	},
}
