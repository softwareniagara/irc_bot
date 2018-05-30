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

func Usage(fset *flag.FlagSet) string {
	var usage bytes.Buffer
	fset.SetOutput(&usage)
	fset.Usage()
	return usage.String()
}

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
		return errors.New(Usage(fset))
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
