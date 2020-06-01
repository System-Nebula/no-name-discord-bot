package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/System-Nebula/no-name-discord-bot/config"

	"github.com/bwmarrin/discordgo"
)

func main() {
	t := os.Getenv("DISCORD_BOT_TOKEN")
	if t == "" {
		log.Fatal("no token was provided.")
	}

	d, err := discordgo.New("Bot " + t)
	if err != nil {
		log.Fatal("error establishing discord session, err=", err)
	}

	mainConfig := config.GetConfig()
	fmt.Println(mainConfig)

	// add a handler - a list of supported event types; possibly all of them at some point.
	d.AddHandler(onMsgCreate)

	d.Open()
	fmt.Println("connection established...")
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-c
	fmt.Println("connection closing...")
	d.Close()
	fmt.Println("connection closed.")
}

func onMsgCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	if m.Author.ID == s.State.User.ID {
		return
	}

	fmt.Println(m.Content)
}
