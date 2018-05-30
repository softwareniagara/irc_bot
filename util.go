package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/whyrusleeping/hellabot"
)

func ParseFlags(msg *hbot.Message, fset *flag.FlagSet) error {
	args, err := shellquote.Split(msg.Content)
	if err != nil {
		return err
	}
	if len(args) == 0 {
		args = []string{""}
	}
	if err = fset.Parse(args[1:]); err == nil {
		return nil
	}
	if err == flag.ErrHelp {
		var usage bytes.Buffer
		fset.SetOutput(&usage)
		fset.Usage()
		return errors.New(usage.String())
	}
	return err
}

func MultiLineReply(bot *hbot.Bot, msg *hbot.Message, s string) {
	s = strings.Replace(s, "\t", "  ", -1)
	r := strings.NewReader(s)
	ss := bufio.NewScanner(r)
	for ss.Scan() {
		bot.Reply(msg, ss.Text())
	}
}
