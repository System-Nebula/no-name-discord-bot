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
			// check user permissions here or in function?
			echoMessage(event, s, strings.Join(msgSplice[1:], " "))
		case ".roles":
			rolesMessage(event, s)
		case ".fbi":
			fbiMessage(event, s)
		}

	}
}

func echoMessage(m *discordgo.MessageCreate, s *discordgo.Session, content string) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	s.ChannelMessageSend(m.ChannelID, content)
}

func IsRoleMember(s *discordgo.Session, UID string, GID string, RoleName string) bool {
	U, _ := s.GuildMember(GID, UID)
	GuildRoles, _ := s.GuildRoles(GID)
	for _, e := range U.Roles {
		for _, r := range GuildRoles {
			if e == r.ID {
				if r.Name == RoleName {
					return true
				}
			}
		}
	}
	return false
}

func fbiMessage(m *discordgo.MessageCreate, s *discordgo.Session) {
	b := IsRoleMember(s, m.Author.ID, m.GuildID, "FBI")
	if b == true {
		s.ChannelMessageSend(m.ChannelID, "Its not safe to talk here..")
	} else {
		s.ChannelMessageSend(m.ChannelID, "Your not with the FBI..")
	}
}

func rolesMessage(m *discordgo.MessageCreate, s *discordgo.Session) {

	if m.Author.ID == s.State.User.ID {
		return
	}
	User, _ := s.GuildMember(m.GuildID, m.Author.ID)
	GuildRoles, _ := s.GuildRoles(m.GuildID)
	for _, e := range User.Roles {
		for _, r := range GuildRoles {
			if e == r.ID {
				s.ChannelMessageSend(m.ChannelID, r.Name)
			}
		}
	}
}
