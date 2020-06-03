package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/System-Nebula/no-name-discord-bot/plugins/base"
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

	// tomlConfigData := readFromFile("config.toml")
	// mainConfig := config.GetConfig(tomlConfigData)
	// fmt.Println(mainConfig)

	// add a handler - a list of supported event types; possibly all of them at some point.
	d.AddHandler(base.Handle)
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

func readFromFile(file string) string {
	c, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal("unable to load file from disk, err=", err)
	}

	t := string(c)
	return t
}
