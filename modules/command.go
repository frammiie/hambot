package modules

import (
	"os"
	"regexp"
	"strings"

	"github.com/gempir/go-twitch-irc/v2"
)

type CommandModule struct {
	Module
	Commands []Command
}

type Command struct {
	Regex  regexp.Regexp
	Handle func(
		module *CommandModule,
		message *twitch.PrivateMessage,
		args ...string,
	)
	Arguments []string
	Required  int
	Level     int
	Commands  []Command
}

const (
	UserLevel int = iota
	ModLevel
	BroadcasterLevel
)

func (h *CommandModule) OnMessage(message *twitch.PrivateMessage) {
	if !strings.HasPrefix(message.Message, os.Getenv("PREFIX")) {
		return
	}

	// Collect parts of message
	parts := collectParts(message)
	if parts == nil {
		return
	}

	command, depth := matchCommand(parts, 0, &h.Commands)
	if command == nil || command.Handle == nil {
		return
	}

	// Check access level
	if determineAccessLevel(&message.User) < command.Level {
		h.Respond(message, "🤡")
		return
	}

	// Check validity arguments
	if len(parts)-(depth+1) < command.Required {
		h.Respond(message, generateUsage(parts[0], command))
		return
	}

	command.Handle(h, message, parts[1:]...)
}

var partsRegex = regexp.MustCompile("\"([^\"]+)\"|\\S+")

func collectParts(message *twitch.PrivateMessage) []string {
	subs := partsRegex.FindAllStringSubmatch(message.Message, -1)
	if subs == nil || len(subs) == 1 {
		return nil
	}
	parts := make([]string, 0, len(subs))
	for _, s := range subs {
		v := s[1]
		if s[1] == "" {
			v = s[0]
		}
		parts = append(parts, v)
	}
	parts[0] = parts[0][len(os.Getenv("PREFIX")):]
	return parts
}

func matchCommand(query []string, depth int, commands *[]Command) (*Command, int) {
	for _, command := range *commands {
		if !command.Regex.MatchString(query[depth]) {
			continue
		}

		var sub *Command
		var sDepth int
		if command.Commands != nil && len(query)-1 >= depth+1 {
			sub, sDepth = matchCommand(query, depth+1, &command.Commands)
		}
		if sub != nil && sub.Handle != nil {
			return sub, sDepth
		} else if command.Handle != nil {
			return &command, depth
		}
	}

	return nil, 0
}

func generateUsage(name string, command *Command) string {
	var usage strings.Builder
	usage.WriteString("Usage: " + os.Getenv("PREFIX") + name)
	for i, arg := range command.Arguments {
		if i < command.Required {
			usage.WriteString(" [" + arg + "]")
		} else {
			usage.WriteString(" <" + arg + ">")
		}
	}
	return usage.String()
}

var authorized = strings.Split(os.Getenv("AUTHORIZED"), " ")

func determineAccessLevel(user *twitch.User) int {
	for _, displayName := range authorized {
		if user.DisplayName == displayName {
			return BroadcasterLevel
		}
	}

	switch {
	case user.Badges["broadcaster"] == 1:
		return BroadcasterLevel
	case user.Badges["moderator"] == 1:
		return ModLevel
	default:
		return UserLevel
	}
}
