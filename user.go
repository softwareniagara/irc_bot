package main

import (
	"github.com/whyrusleeping/hellabot"
	"strings"

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

func (um *UserManager) RemoveUserTrigger() hbot.Trigger {
	return hbot.Trigger{
		Condition: HasPrefix("!user_remove"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			if !um.isAdmin(msg) {
				ReplyTo(bot, msg, "you're not authorized")
				return true
			}
			nick := strings.TrimSpace(strings.TrimPrefix(msg.Content, "!user_remove"))
			u, err := um.s.FindUserByNick(nick)
			if err != nil {
				ErrorReply(bot, msg, err)
				return true
			}
			if err := um.s.RemoveUser(u.RowID); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}
			ReplyTo(bot, msg, "done")
			return true
		},
	}
}

func (um *UserManager) AddUserTrigger() hbot.Trigger {
	return hbot.Trigger{
		Condition: HasPrefix("!user_add"),
		Action: func(bot *hbot.Bot, msg *hbot.Message) bool {
			if !um.isAdmin(msg) {
				ReplyTo(bot, msg, "you're not authorized")
				return true
			}
			nick := strings.TrimSpace(strings.TrimPrefix(msg.Content, "!user_add"))
			u := store.User{
				Nick: nick,
				Role: store.RoleAdmin,
			}
			if err := um.s.InsertUser(&u); err != nil {
				ErrorReply(bot, msg, err)
				return true
			}
			ReplyTo(bot, msg, "done")
			return true
		},
	}
}
