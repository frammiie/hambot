package modules

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/frammiie/hambot/db"
)

var started = time.Now()

var StatsModule = CommandModule{
	Commands: []Command{
		{
			Regex: *regexp.MustCompile("stats$"),
			Handle: func(params *HandlerParams, args ...string) {
				duration := time.Since(started).Round(time.Minute).String()
				uptime := duration[:strings.IndexByte(duration, 'm')+1]

				var count int
				db.Database.QueryRow("SELECT COUNT(*) FROM message").
					Scan(&count)

				params.module.Respond(
					params.message,
					fmt.Sprintf("Current statistics for hambot🍖"+
						"%d 💬 messages | "+
						"%s ⏳ uptime", count, uptime))
			},
		},
	}}
