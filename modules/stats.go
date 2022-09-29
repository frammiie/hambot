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
				uptime := time.Since(started).Round(time.Minute)

				var count int64
				db.Instance.Table("message").Select("COUNT(*)").Scan(&count)

				params.module.Respond(
					params.message,
					fmt.Sprintf("Current statistics for hambot🍖"+
						"%d 💬 messages | "+
						"%s ⏳ uptime", count, uptime,
					),
				)
			},
		},
	}}
