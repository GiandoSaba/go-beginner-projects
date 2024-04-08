package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/Goscord/goscord/goscord/discord"
	"github.com/Goscord/goscord/goscord/gateway"
	"github.com/Goscord/goscord/goscord/gateway/event"
	"github.com/krognol/go-wolfram"
	"github.com/tidwall/gjson"
	witai "github.com/wit-ai/wit-go/v2"

	"github.com/joho/godotenv"
)

var discordClient *gateway.Session
var witAiClient *witai.Client
var wolframClient *wolfram.Client

func requestToWitAi(msg string, client *witai.Client) string {
	resp, err := client.Parse(&witai.MessageRequest{
		Query: msg,
	})
	if err != nil {
		fmt.Println(err)
	}

	data, _ := json.MarshalIndent(resp, "", "\t")
	rough := string(data[:])
	value := gjson.Get(rough, "entities.wit$wolfram_search_query:wolfram_search_query.0.value")
	answer := value.String()

	res, err := wolframClient.GetSpokentAnswerQuery(answer, wolfram.Metric, 1000)
	if err != nil {
		fmt.Println(err)
	}
	return res

}

func main() {

	godotenv.Load(".env")

	fmt.Println("Starting...")

	discordClient = gateway.NewSession(&gateway.Options{
		Token:   os.Getenv("DISCORD_BOT_TOKEN"),
		Intents: gateway.IntentsAll,
	})

	witAiClient = witai.NewClient(os.Getenv("WIT_AI_TOKEN"))

	wolframClient = &wolfram.Client{AppID: os.Getenv("WOLFRAM_APP_ID")}

	discordClient.On(event.EventReady, func() {
		fmt.Println("Logged in as " + discordClient.Me().Tag())
	})

	discordClient.On(event.EventMessageCreate, func(msg *discord.Message) {

		if msg.Author.Bot {
			return
		}

		if len(msg.Mentions) == 0 {
			return
		}

		if msg.Mentions[0].Id == discordClient.Me().Id {

			messageWithoutMention := strings.Replace(msg.Content, "<@"+discordClient.Me().Id+">", "", -1)

			message := requestToWitAi(messageWithoutMention, witAiClient)

			discordClient.Channel.SendMessage(msg.ChannelId, message)
		}

	})

	discordClient.Login()

	select {}

}
