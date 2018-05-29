package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"log"
	"strings"

	"github.com/kballard/go-shellquote"
	"github.com/whyrusleeping/hellabot"
)

var (
	host    string
	nick    string
	channel string
)

func init() {
	flag.StringVar(&host, "addr", "irc.freenode.net:6667", "irc server address")
	flag.StringVar(&nick, "nick", "icholy_bot", "nick to use")
	flag.StringVar(&channel, "chan", "#softwareniagara", "channel to join")
	flag.Parse()
}

var Echo = hbot.Trigger{
	Condition: func(bot *hbot.Bot, msg *hbot.Message) bool {
		return strings.HasPrefix(msg.Content, "!echo")
	},
	Action: func(bot *hbot.Bot, msg *hbot.Message) bool {

		var num int

		fset := flag.NewFlagSet("", flag.ContinueOnError)
		fset.IntVar(&num, "n", 1, "number of times to repeat")

		if err := ParseFlags(msg, fset); err != nil {
			MultiLineReply(bot, msg, err.Error())
			return true
		}

		response := strings.Join(fset.Args(), " ")
		for i := 0; i < num; i++ {
			bot.Reply(msg, response)
		}
		return true
	},
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
		var usage bytes.Buffer
		fset.SetOutput(&usage)
		fset.Usage()
		return errors.New(usage.String())
	}
	return err
}

func MultiLineReply(bot *hbot.Bot, msg *hbot.Message, s string) {
	r := strings.NewReader(s)
	ss := bufio.NewScanner(r)
	for ss.Scan() {
		text := strings.Replace(ss.Text(), "\t", "  ", -1)
		bot.Reply(msg, text)
	}
}

func main() {
	bot, err := hbot.NewBot(host, nick)
	if err != nil {
		log.Fatal(err)
	}
	bot.Channels = []string{channel}
	bot.AddTrigger(Echo)
	bot.Run()
}
