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
	flag.StringVar(&nick, "nick", "softwareniagara_bot", "nick to use")
	flag.StringVar(&channel, "chan", "#softwareniagara_bot", "channel to join")
	flag.Parse()
}

func main() {
	rd, err := OpenReminderDB("reminders.db", channel)
	if err != nil {
		log.Fatal(err)
	}
	defer rd.Close()

	bot, err := hbot.NewBot(host, nick)
	if err != nil {
		log.Fatal(err)
	}
	bot.Channels = []string{channel}
	bot.ThrottleDelay = time.Second
	bot.AddTrigger(FilterChannel(channel))
	bot.AddTrigger(EchoTrigger)
	bot.AddTrigger(rd.Trigger())

	go func() {
		// TODO: find a way to detect when the bot has
		//       successfully connected.
		time.Sleep(time.Second * 10)
		rd.Run(bot)
	}()

	bot.Run()
}
