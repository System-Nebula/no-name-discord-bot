package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

func main() {
	t := os.Getenv("DISCORD_BOT_TOKEN")

	d, err := discordgo.New("Bot " + t)
	if err != nil {
		fmt.Println("error establishing discord session, err=", err)
		os.Exit(1)
	}

	// add a handler
	d.AddHandler(onMsgCreate)

	d.Open()
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-c
	d.Close()
}

func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Content)
}
