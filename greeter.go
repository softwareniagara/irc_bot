package main

import (
	"database/sql"
	"math/rand"
	"time"

	"github.com/whyrusleeping/hellabot"

	"github.com/softwareniagara/irc_bot/store"
)

func init() {
	rand.Seed(time.Now().Unix())
}

func GreeterTrigger(s *store.Store, nick string) hbot.Trigger {
	return hbot.Trigger{
		Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
			return msg.Command == "JOIN" && msg.From != nick
		},
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			u, err := s.FindUserByNick(msg.From)
			if err != nil {
				if err == sql.ErrNoRows {
					ReplyTo(bot, msg, RandomGreeting())
					return true
				}
				ErrorReply(bot, msg, err)
				return true
			}
			if u.Greeting != "" {
				greeting := u.Greeting
				if greeting == "random" {
					greeting = RandomGreeting()
				}
				ReplyTo(bot, msg, greeting)
			}
			return true
		},
	}
}

func RandomGreeting() string {
	greetings := []string{
		"Hello, sunshine!",
		"Howdy, partner!",
		"Hey, howdy, hi!",
		"What’s kickin’, little chicken?",
		"Peek-a-boo!",
		"Howdy-doody!",
		"Hey there, freshman!",
		"My name's Ralph, and I'm a bad guy.",
		"Hi, mister!",
		"I come in peace!",
		"Put that cookie down!",
		"Ahoy, matey!",
		"Hiya!",
		"‘Ello, gov'nor!",
		"Top of the mornin’ to ya!",
		"What’s crackin’?",
		"GOOOOOD MORNING, VIETNAM!",
		"‘Sup, homeslice?",
		"This call may be recorded for training purposes.",
		"Howdy, howdy ,howdy!",
		"How does a lion greet the other animals in the field? A: Pleased to eat you.",
		"Hello, my name is Inigo Montoya.",
		"I'm Batman.",
		"At least, we meet for the first time for the last time!",
		"Hello, who's there, I'm talking.",
		"Here's Johnny!",
		"You know who this is.",
		"Ghostbusters, whatya want?",
		"Yo!",
		"Whaddup.",
		"Greetings and salutations!",
		"Saying Hello to Your Love",
		"‘Ello, mate.",
		"Heeey, baaaaaby.",
		"Hi, honeybunch!",
		"Oh, yoooouhoooo!",
		"How you doin'?",
		"I like your face.",
		"What's cookin', good lookin'?",
		"Howdy, miss.",
		"Why, hello there!",
		"Hey, boo.",
		"Aloha",
		"Hola",
		"Que pasa",
		"Bonjour",
		"Hallo",
		"Ciao",
		"Konnichiwa",
	}
	n := rand.Int() % len(greetings)
	return greetings[n]
}
