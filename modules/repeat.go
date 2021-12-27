package modules

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/frammiie/hambot/db"
)

var minLength, _ = strconv.Atoi(os.Getenv("REPEAT_MIN_LENGTH"))
var excludedContent = strings.Split(
	strings.ToLower(os.Getenv("REPEAT_EXCLUDED_CONTENT")), ";",
)
var excludedUsernames = strings.Split(
	strings.ToLower(os.Getenv("REPEAT_EXCLUDED_USERNAMES")), ";",
)

var RepeatModule = CommandModule{
	Commands: []Command{
		{
			Regex: *regexp.MustCompile("scs"),
			Arguments: []string{
				"min length",
			},
			Handle: func(
				params *HandlerParams,
				args ...string) {
				min := minLength
				if len(args) > 1 {
					parsed, _ := strconv.Atoi(args[1])
					if parsed <= maxHeight {
						min = parsed
					}
				}

				query := strings.Builder{}
				query.WriteString(`
					SELECT content FROM message
					WHERE LENGTH(content) >= ?
				`)

				if len(excludedUsernames) > 1 {
					query.WriteString("AND LOWER(username) NOT IN (''")
					for _, username := range excludedUsernames {
						query.WriteString(", LOWER('" + username + "')")
					}
					query.WriteString(")")
				}

				if len(excludedContent) > 1 {
					for _, excluded := range excludedContent {
						query.WriteString(
							" AND LOWER(content) NOT LIKE LOWER('%" +
								excluded + "%')")
					}
				}

				var message string
				err := db.Database.QueryRow(query.String(), min).Scan(&message)

				if err != nil {
					return
				}

				params.module.Client.Say(params.message.Channel, message)
			},
		},
	},
}
