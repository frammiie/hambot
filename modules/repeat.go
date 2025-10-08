package modules

import (
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/frammiie/hambot/db"
	"github.com/frammiie/hambot/db/model"
)

var minLength, _ = strconv.Atoi(os.Getenv("REPEAT_MIN_LENGTH"))
var excludedContent = strings.Split(
	strings.ToLower(os.Getenv("REPEAT_EXCLUDED_CONTENT")), ";",
)
var excludedUsernames = strings.Split(
	strings.ToLower(os.Getenv("REPEAT_EXCLUDED_USERNAMES")), ";",
)

var RepeatModule = CommandModule{
	Commands: []*Command{
		{
			Regex: *regexp.MustCompile("scs$"),
			Handle: func(params *HandlerParams, args ...string) {
				handleQuery(params, nil, false)
			},
		},
		{
			Regex: *regexp.MustCompile("scse$"),
			Handle: func(params *HandlerParams, args ...string) {
				handleQuery(params, nil, true)
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

				handleQuery(params, &searchQuery, false)
			},
		},
		{
			Regex: *regexp.MustCompile("sce$"),
			Arguments: []string{
				"query",
			},
			Required: 1,
			Handle: func(params *HandlerParams, args ...string) {
				searchQuery := strings.ToLower(ConcatArgs(args...))

				handleQuery(params, &searchQuery, true)
			},
		},
	},
}

func handleQuery(params *HandlerParams, searchQuery *string, extended bool) {
	query := db.Instance.
		Table("message_fts").
		Where("channel", params.message.Channel).
		Where("LENGTH(content) >= ?", minLength).
		Where("content NOT LIKE ?", os.Getenv("PREFIX")+"%").
		Order("RANDOM()").
		Limit(1)

	if searchQuery != nil {
		query.Where("content MATCH ?", searchQuery)
	}

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

	var message model.Message
	result := query.First(&message)

	if result.RowsAffected == 0 {
		params.module.Respond(
			params.message, "No messages found 🔎",
		)
		return
	}

	if !extended {
		params.module.Respond(
			params.message,
			message.Content,
		)
	} else {
		var location, _ = time.LoadLocation("Europe/Berlin")
		
		params.module.Respond(
			params.message,
			fmt.Sprintf(
				"%s | @%s, %s",
				message.Content,
				message.Username,
				message.Created.
					Time().
					In(location).
					Format("2006-01-02 15:04:05 MST"),
			),
		)
	}
}
