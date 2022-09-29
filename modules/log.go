package modules

import (
	"github.com/frammiie/hambot/db"
	"github.com/frammiie/hambot/db/model"
	"github.com/gempir/go-twitch-irc/v3"
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
		Created:  message.Time,
	})
}
