package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/whyrusleeping/hellabot"
)

var EchoTrigger = hbot.Trigger{
	Condition: HasPrefix("!echo"),
	Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
		var num int
		fset := flag.NewFlagSet("echo", flag.ContinueOnError)
		fset.IntVar(&num, "n", 1, "number of times to repeat")
		if err := ParseFlags(msg, fset); err != nil {
			MultiLineReply(bot, msg, err.Error())
			return true
		}
		if num > 3 {
			bot.Reply(msg, fmt.Sprintf("%s: fuck off", msg.From))
			return true
		}
		response := strings.Join(fset.Args(), " ")
		for i := 0; i < num; i++ {
			bot.Reply(msg, response)
		}
		return true
	},
}
