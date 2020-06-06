package base

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Handle - the word Handle
func Handle(s *discordgo.Session, e *discordgo.Event) {
	if e.Type == "MESSAGE_CREATE" {
		mc := e.Struct.(*discordgo.MessageCreate)
		msgSplice := strings.Fields(mc.Content)

		if len(msgSplice) == 0 {
			// if the message is nil, fail safely
			return
		}

		switch msgSplice[0] {
		case ".echo":
			echoMessage(mc, s, strings.Join(msgSplice[1:], " "))
		case ".slap":
			troutSlap(mc, s)
		}
	}
}

func echoMessage(m *discordgo.MessageCreate, s *discordgo.Session, content string) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	s.ChannelMessageSend(m.ChannelID, content)
}

func troutSlap(m *discordgo.MessageCreate, s *discordgo.Session) {
	// required check to disallow bot looping
	if m.Author.ID == s.State.User.ID {
		return
	}

	// check to see that targets are present: "command @target"
	c := strings.Fields(m.Content)
	if len(c) < 2 {
		fmt.Println("not enough args.")
		return
	}

	// get user ID from content
	u, err := getUserFromContent(s, c[1])
	if err != nil {
		fmt.Println(err)
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("_slaps %s around a bit with a large trout._", u.Mention()))
}

func getUserFromContent(s *discordgo.Session, data string) (user *discordgo.User, err error) {
	rex, err := regexp.Compile("[^0-9]+")
	if err != nil {
		return nil, err
	}

	pu := rex.ReplaceAllString(data, "")

	uo, err := s.User(pu)
	if err != nil {
		return nil, err
	}

	if uo == nil {
		return nil, fmt.Errorf("unable to find a user, value=%s", data)
	}

	return uo, nil
}
