package main

import (
	"flag"
	"log"
	"time"

	"github.com/whyrusleeping/hellabot"
)

var (
	host    string
	nick    string
	channel string
)

func init() {
	flag.StringVar(&host, "addr", "irc.freenode.net:6667", "irc server address")
	flag.StringVar(&nick, "nick", "icholy_bot", "nick to use")
	flag.StringVar(&channel, "chan", "#softwareniagara", "channel to join")
	flag.Parse()
}

func main() {
	bot, err := hbot.NewBot(host, nick)
	if err != nil {
		log.Fatal(err)
	}
	bot.Channels = []string{channel}
	bot.ThrottleDelay = time.Second
	bot.AddTrigger(EchoTrigger)
	bot.Run()
}
