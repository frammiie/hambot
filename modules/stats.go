package modules

import (
	"regexp"
	"time"

	"golang.org/x/text/language"
	"golang.org/x/text/message"

	"github.com/frammiie/hambot/db"
)

var started = time.Now()

var StatsModule = CommandModule{
	Commands: []*Command{
		{
			Regex: *regexp.MustCompile("stats$"),
			Handle: func(params *HandlerParams, args ...string) {
				uptime := time.Since(started).Round(time.Second)

				var count int64
				db.Instance.Table("message").Select("COUNT(*)").Scan(&count)

				printer := message.NewPrinter(language.English)

				params.module.Respond(
					params.message,
					printer.Sprintf("Current statistics for hambot 🍖" +
						"%d 💬 messages | " +
						"%s ⏳ uptime", count, uptime,
					),
				)
			},
		},
	}}
