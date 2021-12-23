package modules

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

var maxHeight, _ = strconv.Atoi(os.Getenv("MAX_PYRAMID_HEIGHT"))

var PyramidModule = CommandModule{
	Commands: []Command{
		{
			Regex: *regexp.MustCompile("pyramid"),
			Arguments: []string{
				"content",
				"height",
			},
			Required: 1,
			Handle: func(
				module *CommandModule,
				message *twitch.PrivateMessage,
				args ...string) {
				height := 3
				if len(args) > 1 {
					parsed, _ := strconv.Atoi(args[1])
					if parsed <= maxHeight {
						height = parsed
					}
				}
				for i := 0; i < height+1; i++ {
					module.Client.Say(message.Channel, strings.Repeat(args[0]+" ", i))
				}
				for i := height - 1; i > 0; i-- {
					module.Client.Say(message.Channel, strings.Repeat(args[0]+" ", i))
				}
			},
		},
	},
}
