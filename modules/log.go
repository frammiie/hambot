package modules

import (
	"github.com/frammiie/hambot/db"
	"github.com/frammiie/hambot/db/model"
	"github.com/frammiie/hambot/db/types"

	"github.com/gempir/go-twitch-irc/v4"
)

type LogModule struct {
	Module
}

func (m *LogModule) OnMessage(message *twitch.PrivateMessage) {
	db.Instance.Create(&model.Message{
		Id:       message.ID,
		Content:  message.Message,
		Username: message.User.DisplayName,
		Channel:  message.Channel,
		Created:  types.Timestamp(message.Time),
	})
}
