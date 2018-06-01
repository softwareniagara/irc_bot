package main

import (
	"log"
	"time"

	"github.com/whyrusleeping/hellabot"

	"github.com/softwareniagara/irc_bot/store"
)

func ActivityTrigger(s *store.Store) hbot.Trigger {
	return hbot.Trigger{
		Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
			return true
		},
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			u, err := s.FindUserByNick(msg.From)
			if err != nil {
				return false
			}
			switch msg.Command {
			case "PART", "QUIT":
				u.Active = false
			default:
				u.Active = true
			}
			u.LastActive = time.Now()
			if err := s.UpdateUser(u); err != nil {
				log.Println(err)
				return false
			}
			return false
		},
	}
}
