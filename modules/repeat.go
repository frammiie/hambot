package modules

import (
	"os"
	"regexp"
	"strconv"

	"github.com/frammiie/hambot/db"
	"github.com/gempir/go-twitch-irc/v2"
)

var minLength, _ = strconv.Atoi(os.Getenv("REPEAT_MIN_LENGTH"))

var RepeatModule = CommandModule{
	Commands: []Command{
		{
			Regex: *regexp.MustCompile("scs"),
			Arguments: []string{
				"min length",
			},
			Handle: func(
				module *CommandModule,
				message *twitch.PrivateMessage,
				args ...string) {
				min := minLength
				if len(args) > 1 {
					parsed, _ := strconv.Atoi(args[1])
					if parsed <= maxHeight {
						min = parsed
					}
				}

				quote := db.Quote{}

				err := db.Database.QueryRow(`
					SELECT content FROM message
					WHERE LENGTH(content) >= $1
					ORDER BY RANDOM()
				`, min).Scan(&quote.Content)

				if err != nil {
					return
				}

				module.Client.Say(message.Channel, quote.Content)
			},
		},
	},
}
