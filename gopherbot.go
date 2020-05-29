package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"
	//_ "github.com/go-sql-driver/mysql"
)

type Credentials struct {
	ConsumerKey      string
	ConsumerSecret   string
	AccessToken      string
	AcessTokenSecret string
}

func main() {
	dg, err := discordgo.New("Bot " + "__________")

	// INITIALZE DISCORD
	if err != nil {
		fmt.Println("error creating Discord session,", err)
		return
	}

	dg.AddHandler(messageCreate)
	//dg.AddHandler(messageHistory)

	err = dg.Open()
	if err != nil {
		fmt.Println("Discord - error opening connection,", err)
		return
	}

	// INITIALIZE TWITTER
	fmt.Println("GopherBot - Twitter proof of concept")

	creds := Credentials{
		AccessToken:      "___________________",
		AcessTokenSecret: "___________________",
		ConsumerKey:      "___________________",
		ConsumerSecret:   "___________________",
	}

	config := oauth1.NewConfig(creds.ConsumerKey, creds.ConsumerSecret)
	token := oauth1.NewToken(creds.AccessToken, creds.AcessTokenSecret)

	httpClient := config.Client(oauth1.NoContext, token)

	client := twitter.NewClient(httpClient)

	demux := twitter.NewSwitchDemux()

	demux.Tweet = func(tweet *twitter.Tweet) {
		//fmt.Println(tweet.User.IDStr, "== 2907774137 is ", (tweet.User.IDStr == "2907774137"))
		if tweet.User.IDStr == "2907774137" {
			fmt.Println(tweet.Text)
			dg.ChannelMessageSend("554837317760712706", tweet.Text)
			// s.ChannelMessageSend("550701854875713565", tweet.Text)
		}
	}

	filterParams := &twitter.StreamFilterParams{
		Follow:        []string{"2907774137", "599536606"},
		StallWarnings: twitter.Bool(true),
	}
	stream, err := client.Streams.Filter(filterParams)
	if err != nil {
		log.Fatal(err)
	}
	go demux.HandleChan(stream.Messages)

	// ch := make(chan os.Signal)
	// signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
	// stream.Stop()

	// START DISCORD
	fmt.Println("Bot is now running.  Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc
	dg.Close()
	stream.Stop()

}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	fmt.Println(m.Author, m.ChannelID, m.Content, m.ch)

	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}

// func messageHistory(s *discordgo.Session, m *discordgo.MessageCreate) {
// 	db, err := sql.Open("mysql", "user:passw0rd12356@tcp(127.0.0.1:3306)/history")
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	currentTime := time.Now()
// 	formattedTime := currentTime.Format("2006.01.02 15:04:05")
// 	channelID := m.ChannelID
// 	channelName, _ := s.State.Channel(channelID)
// 	insert, err := db.Query("INSERT INTO history(user, channelid, channel, message,date) VALUES(?,?,?,?,?)", m.Author.Username, channelID, channelName.Name, string(m.Content), formattedTime)
// 	if err != nil {
// 		panic(err.Error())
// 	}
// 	defer insert.Close()
// }
