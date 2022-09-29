package main

import (
	"log"
	"os"
	"strings"

	"github.com/gempir/go-twitch-irc/v3"

	"github.com/frammiie/hambot/db"
	"github.com/frammiie/hambot/modules"
)

type Module interface {
	Hook(*twitch.Client)
}

func main() {
	db.Init()

	client := twitch.NewClient(os.Getenv("USERNAME"), os.Getenv("TOKEN"))

	client.Join(strings.Split(os.Getenv("CHANNELS"), " ")...)

	for _, module := range modules.Instances {
		module.Hook(client)
	}

	client.OnPrivateMessage(func(message twitch.PrivateMessage) {
		for _, module := range modules.Instances {
			module.OnMessage(&message)
		}
	})

	err := client.Connect()
	if err != nil {
		log.Fatal(err)
	}
}
