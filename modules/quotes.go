package modules

import (
	"fmt"
	"regexp"
	"time"

	"github.com/frammiie/hambot/db"
	"github.com/frammiie/hambot/db/model"
	"github.com/frammiie/hambot/db/types"
	"gorm.io/gorm"
)

var QuoteModule = CommandModule{
	Commands: []*Command{
		{
			Regex: *regexp.MustCompile("q(uote)?(s)?$"),
			Arguments: []string{
				"number",
			},
			Required: 1,
			Handle:   findQuote,
			Commands: []*Command{
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
					Handle: searchQuote,
				},
			},
		},
	},
}

func queryQuote(params *HandlerParams, number *string) *gorm.DB {
	query := db.Instance.
		Where("channel = ?", params.message.Channel)
	if number != nil {
		query = query.Where("number = ?", number)
	}

	return query.Limit(1)
}

func respondQuote(params *HandlerParams, quote *model.Quote) {
	params.module.Respond(
		params.message,
		fmt.Sprintf(
			"✒️ #%d - %s - %s (👤 @%s ⏰ %s)",
			quote.Number, quote.Content, quote.Author,
			quote.Submitter, quote.Added.Time().Format("2006-01-02"),
		),
	)
}

func findQuote(params *HandlerParams, args ...string) {
	quote := &model.Quote{}
	result := queryQuote(params, &args[0]).Find(&quote)

	if result.RowsAffected == 0 {
		params.module.Respond(
			params.message,
			fmt.Sprintf("Quote not found! 👀"),
		)
		return
	}

	respondQuote(params, quote)
}

func addQuote(params *HandlerParams, args ...string) {
	var last int
	result := db.Instance.
		Select("number").
		Table("quote").
		Where("channel = ?", params.message.Channel).
		Order("number DESC").
		Limit(1).
		Scan(&last)

	handleErr := func(err error) {
		params.module.Respond(
			params.message,
			fmt.Sprintf("Failed to add quote :("),
		)
		return
	}

	if result.Error != nil {
		handleErr(result.Error)
		return
	}

	new := &model.Quote{
		Number:    last + 1,
		Content:   args[0],
		Author:    args[1],
		Submitter: params.message.User.DisplayName,
		Channel:   params.message.Channel,
		Added:     types.Timestamp(time.Now()),
	}
	result = db.Instance.Create(new)

	if result.Error != nil {
		handleErr(result.Error)
		return
	}

	params.module.Respond(
		params.message, fmt.Sprintf(
			"Added quote %d successfully 📝",
			new.Number,
		),
	)
}

func deleteQuote(params *HandlerParams, args ...string) {
	queryQuote(params, &args[0]).Delete(&model.Quote{})

	params.module.Respond(
		params.message,
		fmt.Sprintf("Deleted quote #%s successfully 💀", args[0]),
	)
}

func searchQuote(params *HandlerParams, args ...string) {
	quote := &model.Quote{}
	query := ConcatArgs(args...)
	result := db.Instance.
		Table("quote_fts").
		Select("quote.*").
		Joins("INNER JOIN quote on quote.id = quote_fts.id").
		Where("quote.channel = ?", params.message.Channel).
		Where("quote_fts.content MATCH ?", query).
		Scan(&quote)

	if result.RowsAffected == 0 {
		params.module.Respond(
			params.message,
			fmt.Sprintf("Quote with query not found! 👀"),
		)
		return
	}

	respondQuote(params, quote)
}
