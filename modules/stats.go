package modules

import (
	"fmt"
	"regexp"
	"time"

	"github.com/frammiie/hambot/db"
)

var started = time.Now()

var StatsModule = CommandModule{
	Commands: []Command{
		{
			Regex: *regexp.MustCompile("stats$"),
			Handle: func(params *HandlerParams, args ...string) {
				duration := time.Since(started)
				uptime := fmt.Sprintf(
					"%.1fh%.1fm",
					duration.Hours(),
					duration.Minutes(),
				)

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
