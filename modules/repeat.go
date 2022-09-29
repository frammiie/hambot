package modules

import (
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/frammiie/hambot/db"
	"gorm.io/gorm"
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
			Regex: *regexp.MustCompile("scs$"),
			Handle: func(params *HandlerParams, args ...string) {
				query := db.Instance.
					Select("content").
					Table("message").
					Where("channel", params.message.Channel)

				includeConditions(params, query)

				var message string
				query.Find(&message)

				params.module.Client.Say(
					params.message.Channel,
					message,
				)
			},
		},
		{
			Regex: *regexp.MustCompile("sc$"),
			Arguments: []string{
				"query",
			},
			Required: 1,
			Handle: func(params *HandlerParams, args ...string) {
				searchQuery := strings.ToLower(ConcatArgs(args...))

				var message string
				query := db.Instance.
					Select("content").
					Table("message_fts").
					Where("content MATCH ?", searchQuery)

				includeConditions(params, query)

				result := query.Scan(&message)

				if result.RowsAffected == 0 {
					params.module.Respond(
						params.message, "No messages found 🔎",
					)
					return
				}

				params.module.Client.Say(params.message.Channel, message)
			},
		},
	},
}

func includeConditions(params *HandlerParams, query *gorm.DB) {
	query.
		Where("channel", params.message.Channel).
		Where("LENGTH(content) >= ?", minLength).
		Where("content NOT LIKE ?", os.Getenv("PREFIX")+"%").
		Order("RANDOM()").
		Limit(1)

	if len(excludedUsernames[0]) > 0 {
		query.Where("LOWER(username) NOT IN ?", excludedUsernames)
	}
	if len(excludedContent[0]) > 0 {
		for _, excluded := range excludedContent {
			query.Where(
				"LOWER(content) NOT LIKE ?", "%"+excluded+"%",
			)
		}
	}
}
