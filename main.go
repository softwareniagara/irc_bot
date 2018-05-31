package main

import (
	"flag"
	"log"
	"time"

	"github.com/whyrusleeping/hellabot"

	"github.com/softwareniagara/irc_bot/store"
)

var (
	host    string
	nick    string
	channel string
	dbname  string
)

func init() {
	flag.StringVar(&host, "addr", "irc.freenode.net:6667", "irc server address")
	flag.StringVar(&nick, "nick", "jimmy_38545", "nick to use")
	flag.StringVar(&channel, "chan", "#softwareniagara", "channel to join")
	flag.StringVar(&dbname, "dbname", "bot.db", "database filename")
	flag.Parse()
}

func main() {
	s, err := store.Open(dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	rd := NewReminderDB(s, channel)
	um := NewUserManager(s)

	bot, err := hbot.NewBot(host, nick)
	if err != nil {
		log.Fatal(err)
	}
	bot.Channels = []string{channel}
	bot.ThrottleDelay = time.Second
	bot.AddTrigger(FilterChannel(channel))
	bot.AddTrigger(TellTrigger)
	bot.AddTrigger(EchoTrigger)
	bot.AddTrigger(rd.Trigger())
	bot.AddTrigger(um.AddUserTrigger())
	bot.AddTrigger(um.RemoveUserTrigger())
	bot.AddTrigger(Responder(nick))

	go func() {
		// TODO: find a way to detect when the bot has
		//       successfully connected.
		time.Sleep(time.Second * 10)
		rd.Run(bot)
	}()

	bot.Run()
}
