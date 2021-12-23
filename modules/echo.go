package modules

import (
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gempir/go-twitch-irc/v2"
)

type EchoModule struct {
	Module
}

var echoCount, _ = strconv.Atoi(os.Getenv("ECHO_COUNT"))
var echoLength, _ = strconv.Atoi(os.Getenv("ECHO_LENGTH"))
var echoDecay, _ = strconv.Atoi(os.Getenv("ECHO_DECAY"))

var messages = make([]string, echoLength)
var cursor = 0
var mutex = sync.Mutex{}

func addEchoMessage(message string) {
	mutex.Lock()
	if cursor++; cursor > echoLength-1 {
		cursor = 0
	}
	messages[cursor] = message

	defer mutex.Unlock()
}

func (m *EchoModule) Hook(client *twitch.Client) {
	m.Module.Hook(client)

	// Makes sure the oldest messages in memory decay
	ticker := time.NewTicker(time.Duration(echoDecay) * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				addEchoMessage("")
				break
			}
		}
	}()
}

func (m *EchoModule) OnMessage(message *twitch.PrivateMessage) {
	if strings.HasPrefix(message.Message, os.Getenv("PREFIX")) {
		return
	}

	addEchoMessage(message.Message)

	if countEqual(message.Message) >= echoCount {
		m.Client.Say(message.Channel, message.Message)
		// Clear
		messages = make([]string, echoLength)
	}
}

func countEqual(message string) int {
	c := 0
	for _, previous := range messages {
		if message == previous {
			c++
		}
	}

	return c
}
