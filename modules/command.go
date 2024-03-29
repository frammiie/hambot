package modules

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/gempir/go-twitch-irc/v3"
)

type CommandModule struct {
	Module
	Commands []*Command
}

type HandlerParams struct {
	module  *CommandModule
	message *twitch.PrivateMessage
}

type Command struct {
	Regex  regexp.Regexp
	Handle func(
		params *HandlerParams,
		args ...string,
	)
	Arguments []string
	Required  int
	Level     int
	Cooldown  float64
	LastUsed  time.Time
	Commands  []*Command
}

const (
	UserLevel int = iota
	ModLevel
	BroadcasterLevel
)

var authorized = strings.Split(os.Getenv("AUTHORIZED"), " ")
var defaultCooldown, _ = strconv.ParseFloat(
	os.Getenv("DEFAULT_COOLDOWN_SECS"), 64,
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

	accessLevel := accessLevel(&message.User)

	if activeCooldown(command, accessLevel) {
		return
	}

	if accessLevel < command.Level {
		h.Respond(message, "🤡")
		return
	}

	command.LastUsed = time.Now()

	// Check validity arguments
	if len(parts)-(depth+1) < command.Required {
		h.Respond(message, generateUsage(parts[0], command))
		return
	}

	command.Handle(
		&HandlerParams{module: h, message: message},
		parts[depth+1:]...,
	)
}

var partsRegex = regexp.MustCompile("\"([^\"]+)\"|\\S+")

func collectParts(message *twitch.PrivateMessage) []string {
	subs := partsRegex.FindAllStringSubmatch(message.Message, -1)
	if subs == nil {
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

func matchCommand(query []string, depth int, commands *[]*Command) (*Command, int) {
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
			return command, depth
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

func accessLevel(user *twitch.User) int {
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

func activeCooldown(command *Command, accessLevel int) bool {
	if accessLevel < ModLevel && !command.LastUsed.IsZero() {
		cooldown := command.Cooldown
		if cooldown == 0 {
			cooldown = defaultCooldown
		}
		if time.Since(command.LastUsed).Seconds() < cooldown {
			return true
		}
	}
	return false
}
