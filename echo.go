package main

import (
	"flag"
	"strings"

	"github.com/whyrusleeping/hellabot"
)

var EchoTrigger = hbot.Trigger{
	Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
		return strings.HasPrefix(msg.Content, "!echo")
	},
	Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
		var num int
		fset := flag.NewFlagSet("echo", flag.ContinueOnError)
		fset.IntVar(&num, "n", 1, "number of times to repeat")
		if err := ParseFlags(msg, fset); err != nil {
			MultiLineReply(bot, msg, err.Error())
			return true
		}
		var target string
		if strings.Contains(msg.To, "#") {
			target = msg.To
			if num > 3 {
				bot.Reply(msg, "sending as private message XD")
				target = msg.From
			}
		} else {
			target = msg.From
		}
		response := strings.Join(fset.Args(), " ")
		for i := 0; i < num; i++ {
			bot.Msg(target, response)
		}
		return true
	},
}
