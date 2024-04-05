package main

import (
	"fmt"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Goscord/goscord/goscord/discord"
	"github.com/Goscord/goscord/goscord/gateway"
	"github.com/Goscord/goscord/goscord/gateway/event"
)

var client *gateway.Session

var greetings = []string{"Hi!", "Hello!", "Hey!"}

func runCommand(command string, msg *discord.Message) {
	switch command {
	case "ping":
		client.Channel.SendMessage(msg.ChannelId, "Pong!")
	case "age":
		age(msg)
	default:
		client.Channel.SendMessage(msg.ChannelId, "I don't know that command!")
	}
}

func age(msg *discord.Message) {
	yob, err := strconv.Atoi(strings.Split(msg.Content, " ")[1])
	if err != nil {
		fmt.Println(err)
	}
	client.Channel.SendMessage(msg.ChannelId, "Your age is "+strconv.Itoa(time.Now().Year()-yob))
}

func main() {
	fmt.Println("Starting...")

	client = gateway.NewSession(&gateway.Options{
		Token:   os.Getenv("BOT_TOKEN"),
		Intents: gateway.IntentsAll,
	})

	client.On(event.EventReady, func() {
		fmt.Println("Logged in as " + client.Me().Tag())
	})

	client.On(event.EventMessageCreate, func(msg *discord.Message) {

		if msg.Author.Bot {
			return
		}

		if msg.Mentions[0].Id == client.Me().Id {

			greeting := greetings[rand.Intn(len(greetings))] + "<@" + msg.Author.Id + ">!"

			client.Channel.SendMessage(msg.ChannelId, greeting)
		}

		if strings.HasPrefix(msg.Content, "!") {
			command := strings.Split(msg.Content, " ")[0][1:]
			runCommand(command, msg)
		}

	})

	client.Login()

	select {}

}
