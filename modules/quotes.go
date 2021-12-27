package modules

import (
	"fmt"
	"regexp"

	"github.com/frammiie/hambot/db"
)

var QuoteModule = CommandModule{
	Commands: []Command{
		{
			Regex: *regexp.MustCompile("quotes?$"),
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
			quote.Submitter, quote.Added.Format("2006-02-01"),
		),
	)
}

func quote(params *HandlerParams, args ...string) {
	findQuote(
		params,
		`SELECT
			number, content, author, submitter, added
		FROM quote
		WHERE number = $1`,
		args[0],
	)
}

func addQuote(params *HandlerParams, args ...string) {
	var next int
	err := db.Database.QueryRow(`
		SELECT  number + 1 as number
		FROM    quote q1
		WHERE   NOT EXISTS
				(
				SELECT  NULL
				FROM    quote q2
				WHERE   q2.number = q1.number + 1
				)
		ORDER BY
				number
		LIMIT 1;
	`).Scan(&next)

	if err != nil {
		next = 1
	}

	db.Database.Exec(`
		INSERT INTO quote (
			number, content, author, submitter, added
		) VALUES (
			$1, $2, $3, $4, CURRENT_TIMESTAMP
		)`,
		next, args[1], args[2], params.message.User.DisplayName)

	params.module.Respond(
		params.message, fmt.Sprintf("Added quote %d successfully 📝", next),
	)
}

func deleteQuote(params *HandlerParams, args ...string) {
	db.Database.Exec("DELETE FROM quote WHERE number = $1", args[1])

	params.module.Respond(
		params.message,
		fmt.Sprintf("Deleted quote #%s successfully 💀", args[1]),
	)
}

func searchQuote(params *HandlerParams, args ...string) {
	findQuote(
		params,
		`SELECT
			number, content, author, submitter, added
		FROM quote
		WHERE content LIKE $1`,
		"%"+args[1]+"%",
	)
}
