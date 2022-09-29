package modules

import (
	"github.com/gempir/go-twitch-irc/v3"
)

type IModule interface {
	Hook(client *twitch.Client)
	OnMessage(message *twitch.PrivateMessage)
}

type Module struct {
	Client *twitch.Client
}

func (m *Module) Hook(client *twitch.Client) {
	m.Client = client
}

func (h *Module) Respond(original *twitch.PrivateMessage, message string) {
	h.Client.Reply(original.Channel, original.ID, message)
}

var Instances = []IModule{
	&LogModule{},
	&QuoteModule,
	&EchoModule{},
	&RepeatModule,
	&StatsModule,
	&PyramidModule,
}
