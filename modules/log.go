package modules

import (
	"github.com/frammiie/hambot/db"
	"github.com/gempir/go-twitch-irc/v2"
)

type LogModule struct {
	Module
}

func (m *LogModule) OnMessage(message *twitch.PrivateMessage) {
	db.Database.Exec(
		"INSERT INTO message VALUES ($1, $2, $3, $4, $5)",
		message.ID, message.Message, message.User.DisplayName,
		message.Channel, message.Time,
	)
}
