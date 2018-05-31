package main

import (
	"flag"
	"strings"

	"github.com/whyrusleeping/hellabot"

	"github.com/softwareniagara/irc_bot/store"
)

func UserTrigger(s *store.Store) hbot.Trigger {
	return hbot.Trigger{
		Condition: HasCommand("!user"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {

			if err := s.Authorized(msg.From, store.RoleAdmin); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}

			var (
				create bool
				remove bool
				update bool
				info   bool
				role   = store.RoleUser
			)

			fset := flag.NewFlagSet("", flag.ContinueOnError)
			fset.BoolVar(&create, "create", false, "create a new user")
			fset.BoolVar(&remove, "remove", false, "remove an existing user")
			fset.BoolVar(&update, "update", false, "update existing user")
			fset.BoolVar(&info, "info", false, "show user info")
			fset.Var(&role, "role", "admin|user|idiot|banned")

			if err := ParseFlags(msg, fset); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}
			nick = strings.TrimSpace(strings.Join(fset.Args(), " "))
			if nick == "" {
				bot.Reply(msg, "you need to provide a nick as a positional argument")
				return true
			}

			if create {
				if err := s.InsertUser(&store.User{
					Nick: nick,
					Role: role,
				}); err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				ReplyTo(bot, msg, "ok")
				return true
			}

			if remove {
				u, err := s.FindUserByNick(nick)
				if err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				if err := s.RemoveUser(u.RowID); err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				ReplyTo(bot, msg, "ok")
				return true
			}

			if update {
				u, err := s.FindUserByNick(nick)
				if err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				u.Role = role
				if err := s.UpdateUser(u); err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				ReplyTo(bot, msg, "ok")
				return true
			}

			if info {
				u, err := s.FindUserByNick(nick)
				if err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				ReplyTo(bot, msg, u.String())
				return true
			}

			ReplyTo(bot, msg, "you need to specify the action")
			return true
		},
	}
}
