package modules

import (
	"fmt"
	"regexp"
	"strconv"

	"github.com/frammiie/hambot/db"
)

var QuoteModule = CommandModule{
	Commands: []Command{
		{
			Regex: *regexp.MustCompile("q(uote)?(s)?$"),
			Arguments: []string{
				"number",
			},
			Required: 1,
			Handle:   quote,
			Commands: []Command{
				{
					Regex: *regexp.MustCompile("add$"),
					Arguments: []string{
						"\"content\"",
						"author",
					},
					Required: 2,
					Level:    ModLevel,
					Handle:   addQuote,
				},
				{
					Regex: *regexp.MustCompile("del(ete)?$"),
					Arguments: []string{
						"number",
					},
					Required: 1,
					Level:    ModLevel,
					Handle:   deleteQuote,
				},
				{
					Regex: *regexp.MustCompile("search$"),
					Arguments: []string{
						"query",
					},
					Required: 1,
					Handle:   searchQuote,
				},
			},
		},
	},
}

func quote(params *HandlerParams, args ...string) {
	_, err := strconv.Atoi(args[0])
	if err != nil {
		searchQuote(params, args...)
		return
	}

	findQuote(
		params,
		`SELECT
			number, content, author, submitter, added
		FROM quote
		WHERE
			channel = $1 AND
			number = $2`,
		params.message.Channel, args[0],
	)
}

func findQuote(params *HandlerParams, statement string, args ...interface{}) {
	quote := db.Quote{}
	err := db.Database.QueryRow(
		statement, args...,
	).Scan(
		&quote.Number, &quote.Content, &quote.Author,
		&quote.Submitter, &quote.Added,
	)

	if err != nil {
		params.module.Respond(
			params.message,
			fmt.Sprintf("Quote not found! 👀"),
		)
		return
	}

	params.module.Respond(
		params.message,
		fmt.Sprintf(
			"✒️ #%d - %s - %s (👤 @%s ⏰ %s)",
			quote.Number, quote.Content, quote.Author,
			quote.Submitter, quote.Added.Format("2006-01-02"),
		),
	)
}

func addQuote(params *HandlerParams, args ...string) {
	var next int
	err := db.Database.QueryRow(`
			SELECT  number + 1 as number
			FROM    quote q1
			WHERE
				channel = $1 AND
				NOT EXISTS
				(
					SELECT  NULL
					FROM    quote q2
					WHERE
						channel = $1 AND
						q2.number = q1.number + 1
				)
			ORDER BY
					number
			LIMIT 1;
		`,
		params.message.Channel).
		Scan(&next)

	if err != nil {
		next = 1
	}

	_, err = db.Database.Exec(`
		INSERT INTO quote (
			number, content, author, submitter, channel, added
		) VALUES (
			$1, $2, $3, $4, $5, CURRENT_TIMESTAMP
		)`,
		next, args[0], args[1], params.message.User.DisplayName, params.message.Channel)

	if err != nil {
		params.module.Respond(
			params.message,
			fmt.Sprintf("Failed to add quote :("),
		)
		return
	}

	params.module.Respond(
		params.message, fmt.Sprintf("Added quote %d successfully 📝", next),
	)
}

func deleteQuote(params *HandlerParams, args ...string) {
	db.Database.Exec(`
		DELETE FROM quote
		WHERE
			channel = $1 AND
			number = $2`,
		params.message.Channel,
		args[0],
	)

	params.module.Respond(
		params.message,
		fmt.Sprintf("Deleted quote #%s successfully 💀", args[0]),
	)
}

func searchQuote(params *HandlerParams, args ...string) {
	findQuote(
		params,
		`SELECT
			number, content, author, submitter, added
		FROM quote
		WHERE
			channel = $1 AND
			content LIKE $2`,
		params.message.Channel,
		"%"+args[0]+"%",
	)
}
