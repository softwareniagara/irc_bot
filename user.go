package main

import (
	"flag"
	"strings"

	"github.com/whyrusleeping/hellabot"

	"github.com/softwareniagara/irc_bot/store"
)

type UserManager struct {
	s *store.Store
}

func NewUserManager(s *store.Store) *UserManager {
	return &UserManager{s}
}

func (um *UserManager) isAdmin(msg *hbot.Message) bool {
	u, err := um.s.FindUserByNick(msg.From)
	if err != nil {
		return false
	}
	return u.Role == store.RoleAdmin
}

func (um *UserManager) Trigger() hbot.Trigger {
	return hbot.Trigger{
		Condition: HasPrefix("!user"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {

			if !um.isAdmin(msg) {
				ReplyTo(bot, msg, "no can do")
				return true
			}

			var (
				create bool
				remove bool
				update bool
				role   = store.RoleUser
			)

			fset := flag.NewFlagSet("", flag.ContinueOnError)
			fset.BoolVar(&create, "create", false, "create a new user")
			fset.BoolVar(&remove, "remove", false, "remove an existing user")
			fset.BoolVar(&update, "update", false, "update existing user")
			fset.Var(&role, "role", "admin|user|idiot|banned")

			if err := ParseFlags(msg, fset); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}
			nick = strings.Join(fset.Args(), " ")
			if nick == "" {
				bot.Reply(msg, "you need to provide a nick as a positional argument")
				return true
			}

			if create {
				if err := um.s.InsertUser(&store.User{
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
				u, err := um.s.FindUserByNick(nick)
				if err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				if err := um.s.RemoveUser(u.RowID); err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				ReplyTo(bot, msg, "ok")
				return true
			}

			if update {
				u, err := um.s.FindUserByNick(nick)
				if err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				u.Role = role
				if err := um.s.UpdateUser(u); err != nil {
					ErrorReply(bot, msg, err)
					return true
				}
				ReplyTo(bot, msg, "ok")
				return true
			}

			ReplyTo(bot, msg, "you need to specify the action")
			return true
		},
	}
}
