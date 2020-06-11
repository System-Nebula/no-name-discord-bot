package base

import (
	"fmt"
	"math/rand"
	"regexp"
	"strconv"
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
		case ".rest":
			rest(mc, s)
		case ".attack":
			atk(mc, s)
		case ".remindme":
			go reminder(mc, s)
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

func reminder(m *discordgo.MessageCreate, s *discordgo.Session) {

	//message := strings.Fields(m.Content)[2]
	message := strings.Fields(m.Content)

	// take the given length of time and create a timer in seconds for the duration
	stime := strings.Fields(m.Content)[1]
	itime, err := strconv.Atoi(stime)
	if err != nil {
		fmt.Println("problem with atoi")
		s.ChannelMessageSend(m.ChannelID, "he was rigt, atoi err...")

		return
	}
	timer := time.AfterFunc(time.Second*time.Duration(itime), func() {
		fmt.Println("timer stopped")

		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("It's time to do that thing...%s", strings.Join(message[2:], " ")))

	})
	defer timer.Stop()
	s.ChannelMessageSend(m.ChannelID, "Ok i'll remind you later")
	time.Sleep(time.Second*time.Duration(itime) + 1)
}

// Command: .slap
var troutSlapTimers = make(map[string]time.Time)
var restTimers = make(map[string]time.Time)

var maxHP = make(map[string]int)
var userHP = make(map[string]int)

var playerSheets = make(map[string]*charactersSheet)

type charactersSheet struct {
	maxHP  int
	UserHP int
}

func newPlayer(ID string) {
	playerSheets[ID] = &charactersSheet{
		maxHP:  10,
		UserHP: 10,
	}
}

func atk(m *discordgo.MessageCreate, s *discordgo.Session) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	c := strings.Fields(m.Content)
	if len(c) < 2 {
		return
	}

	if len(m.Mentions) == 0 {
		return
	}
	target := m.Mentions[0]
	targetSheet, ok := playerSheets[target.ID]
	if !ok {
		s.ChannelMessageSend(m.ChannelID, "Targeted player does not have a character sheet yet")
	} else {
		// have a valid target
		// see how much health they have
		if targetSheet.UserHP > 1 {
			playerSheets[target.ID].UserHP -= 1
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("DEBUG - Assign %s -1 health", target.Username))
		} else {
			s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("%s be dead already", target.Username))

		}
	}
}

func rest(m *discordgo.MessageCreate, s *discordgo.Session) {
	// TODO mostly redundant boilerplate code (spam filter & timers
	timeout := 60 * time.Second
	if m.Author.ID == s.State.User.ID {
		return
	}

	t, ok := restTimers[m.Author.ID]
	if ok {
		if time.Now().Sub(t) < timeout {
			return
		}
	}
	// TODO Usage for rest

	// Check if player stats exist
	_, ok = playerSheets[m.Author.ID]
	if !ok {
		newPlayer(m.Author.ID)
	}

	// Check if max health or send rest message
	if playerSheets[m.Author.ID].UserHP == playerSheets[m.Author.ID].maxHP {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Why Rest? %s you have full health", m.Author.Username))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("Get some rest %s, you only have %d hp left", m.Author.Username, playerSheets[m.Author.ID].UserHP))
	}
}

func troutSlap(m *discordgo.MessageCreate, s *discordgo.Session) {
	timeout := 10 * time.Second

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
	HitChance := rand.Intn(6-1) + 1
	if HitChance != 1 {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("_slaps %s around a bit with a large trout._", u.Mention()))
	} else {
		s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("_ in an attempt to slap %s around; you fall on your ass and get a mouthful of large trout._", u.Mention()))
	}
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
