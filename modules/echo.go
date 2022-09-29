package modules

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

type EchoModule struct {
	Module
}

type EchoMessage struct {
	Message string
	Time    time.Time
}

var echoCount, _ = strconv.Atoi(os.Getenv("ECHO_COUNT"))
var echoBufferSize, _ = strconv.Atoi(os.Getenv("ECHO_BUFFER_SIZE"))
var echoExpiration, _ = strconv.ParseFloat(os.Getenv("ECHO_EXPIRE_SECS"), 64)

var messages = make([]*EchoMessage, echoBufferSize)
var cursor = 0

func addEchoMessage(message string) {
	if cursor++; cursor > echoBufferSize-1 {
		cursor = 0
	}
	messages[cursor] = &EchoMessage{
		Message: message,
		Time:    time.Now(),
	}
}

func (m *EchoModule) OnMessage(message *twitch.PrivateMessage) {
	if strings.HasPrefix(message.Message, os.Getenv("PREFIX")) {
		return
	}

	addEchoMessage(message.Message)

	if countEqual(messages[cursor]) >= echoCount {
		m.Client.Say(message.Channel, message.Message)
		// Clear
		messages = make([]*EchoMessage, echoBufferSize)
	}
}

func countEqual(message *EchoMessage) int {
	c := 0
	for _, previous := range messages {
		if previous != nil &&
			message.Message == previous.Message &&
			message.Time.Sub(previous.Time).Seconds() < echoExpiration {
			c++
		}
	}

	return c
}
