package base

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handle - the word Handle
func Handle(s *discordgo.Session, e *discordgo.Event) {
	if e.Type == "MESSAGE_CREATE" {
		event := e.Struct.(*discordgo.MessageCreate)
		msgSplice := strings.Fields(event.Content)

		switch msgSplice[0] {
		case ".echo":
			echoMessage(event, s, strings.Join(msgSplice[1:], " "))
		}
	}
}

func echoMessage(m *discordgo.MessageCreate, s *discordgo.Session, content string) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	s.ChannelMessageSend(m.ChannelID, content)
}
