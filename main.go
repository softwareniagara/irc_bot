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
	admin   string
)

func init() {
	flag.StringVar(&host, "addr", "irc.freenode.net:6667", "irc server address")
	flag.StringVar(&nick, "nick", "jimmy_38545", "nick to use")
	flag.StringVar(&channel, "chan", "#softwareniagara", "channel to join")
	flag.StringVar(&dbname, "dbname", "bot.db", "database filename")
	flag.StringVar(&admin, "admin", "", "create an admin user with this nick")
	flag.Parse()
}

func main() {
	s, err := store.Open(dbname)
	if err != nil {
		log.Fatal(err)
	}
	defer s.Close()

	if admin != "" {
		if err := s.InsertUser(&store.User{
			Nick: admin,
			Role: store.RoleAdmin,
		}); err != nil {
			log.Fatal(err)
		}
		return
	}

	bot, err := hbot.NewBot(host, nick)
	if err != nil {
		log.Fatal(err)
	}
	bot.Channels = []string{channel}
	bot.ThrottleDelay = time.Second
	bot.AddTrigger(FilterChannel(channel))
	bot.AddTrigger(ActivityTrigger(s))
	bot.AddTrigger(GreeterTrigger(s))
	bot.AddTrigger(TellTrigger(s))
	bot.AddTrigger(EchoTrigger(s))
	bot.AddTrigger(ReminderTrigger(s))
	bot.AddTrigger(UserTrigger(s))
	bot.AddTrigger(Responder(nick))
	bot.AddTrigger(HelpTrigger())

	go func() {
		// TODO: find a way to detect when the bot has
		//       successfully connected.
		time.Sleep(time.Second * 10)
		ReminderNotifyLoop(bot, s, channel)
	}()

	bot.Run()
}
