package base

import (
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// A plugin-scoped map for tracking the time of the last activity of a user; limited to message
var lastSeenTimers = make(map[string]time.Time)

// Handle - the word Handle
func Handle(s *discordgo.Session, e *discordgo.Event) {
	if e.Type == "MESSAGE_CREATE" {
		mc := e.Struct.(*discordgo.MessageCreate)
		msgSplice := strings.Fields(mc.Content)

		lastSeenTimers[mc.Author.ID] = time.Now()

		if len(msgSplice) == 0 {
			// if the message is nil, fail safely
			return
		}

		switch msgSplice[0] {
		case ".echo":
			echoMessage(mc, s, strings.Join(msgSplice[1:], " "))
		case ".slap":
			troutSlap(mc, s)
		case ".roles":
			rolesMessage(mc, s)
		case ".fbi":
			fbiMessage(mc, s)
		case ".lastseen":
			lastSeen(mc, s)
		}

	}
}

func echoMessage(m *discordgo.MessageCreate, s *discordgo.Session, content string) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	s.ChannelMessageSend(m.ChannelID, content)
}

// Command: .lastseen
func lastSeen(m *discordgo.MessageCreate, s *discordgo.Session) {
	if m.Author.ID == s.State.User.ID {
		return
	}

		// check to see that targets are present: "@author command @target"
		c := strings.Fields(m.Content)
		if len(c) < 2 {
			uc, err := s.UserChannelCreate(m.Author.ID)
			if err != nil {
				fmt.Println("unable to create DM with user, user=%s", uc)
				return
			}
			s.ChannelMessageSend(uc.ID, "usage: .lastseen @username")
			return
		}
	
		// get user ID from content
		u, err := getUser(c[1], s)
		if err != nil {
			fmt.Println(err)
			return
		}

	t, ok := lastSeenTimers[u.ID]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s has not been seen.", u.String()))
		return
	}
	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s's last message was sent at %s", u.String(), t.UTC()))
}

// Command: .slap
var troutSlapTimers = make(map[string]time.Time)

func troutSlap(m *discordgo.MessageCreate, s *discordgo.Session) {
	timeout := time.Duration(10*time.Second)

	// required check to disallow bot looping
	if m.Author.ID == s.State.User.ID {
		return
	}

	// get debounce timer value for user
	t, ok := troutSlapTimers[m.Author.ID]
	if ok {
		if time.Now().Sub(t) < timeout {
			return
		}
	}

	// check to see that targets are present: "@author command @target"
	c := strings.Fields(m.Content)
	if len(c) < 2 {
		uc, err := s.UserChannelCreate(m.Author.ID)
		if err != nil {
			fmt.Println("unable to create DM with user, user=%s", uc)
			return
		}
		s.ChannelMessageSend(uc.ID, "usage: .slap @username")
		return
	}

	// get user ID from content
	u, err := getUser(c[1], s)
	if err != nil {
		fmt.Println(err)
		return
	}

	s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("_slaps %s around a bit with a large trout._", u.Mention()))
	
	// set debounce timer value for user
	troutSlapTimers[m.Author.ID] = time.Now()
}

func isRoleMember(s *discordgo.Session, UID string, GID string, RoleName string) bool {
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
	b := isRoleMember(s, m.Author.ID, m.GuildID, "FBI")
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

func getUser(data string, s *discordgo.Session) (user *discordgo.User, err error) {
	if !strings.HasPrefix(data, "<@!") {
		return nil, fmt.Errorf("invalid user, value=%s", data)
	}

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