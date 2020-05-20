package cmd

import (
	"fmt"
	"github.com/Necroforger/dgrouter/exrouter"
	"github.com/bwmarrin/discordgo"
	"github.com/darenliang/MikuBotGo/config"
	"github.com/darenliang/MikuBotGo/framework"
	"html"
	"math/rand"
	"sync"
	"time"
)

type TriviaResponse struct {
	ResponseCode int `json:"response_code"`
	Results      []struct {
		Category         string   `json:"category"`
		Type             string   `json:"type"`
		Difficulty       string   `json:"difficulty"`
		Question         string   `json:"question"`
		CorrectAnswer    string   `json:"correct_answer"`
		IncorrectAnswers []string `json:"incorrect_answers"`
	} `json:"results"`
}

var emojis = []string{
	"\x31\xef\xb8\x8f\xe2\x83\xa3",
	"\x32\xef\xb8\x8f\xe2\x83\xa3",
	"\x33\xef\xb8\x8f\xe2\x83\xa3",
	"\x34\xef\xb8\x8f\xe2\x83\xa3",
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Trivia command
func Trivia(ctx *exrouter.Context) {
	var (
		lock     sync.RWMutex
		embedMsg *discordgo.Message
		callback = make(chan struct{})
		idx      int
	)

	triviaReponse := &TriviaResponse{}
	err := framework.UrlToStruct("https://opentdb.com/api.php?amount=1&category=31&type=multiple", triviaReponse)
	if err != nil {
		ctx.Reply(":cry: An error has occurred")
		return
	}

	question := triviaReponse.Results[0]

	answers := make([]string, 0)
	answers = append(answers, question.CorrectAnswer)
	answers = append(answers, question.IncorrectAnswers...)
	rand.Shuffle(len(answers), func(i, j int) { answers[i], answers[j] = answers[j], answers[i] })

	embed := &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Color:  config.EmbedColor,
		Description: fmt.Sprintf(
			"%s\n\n"+
				":one: %s\n"+
				":two: %s\n"+
				":three: %s\n"+
				":four: %s",
			html.UnescapeString(question.Question),
			html.UnescapeString(answers[0]),
			html.UnescapeString(answers[1]),
			html.UnescapeString(answers[2]),
			html.UnescapeString(answers[3]),
		),
		Title: fmt.Sprintf("Anime trivia"),
	}

	embedMsg, _ = ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
	for i := 0; i < 4; i++ {
		ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, embedMsg.ID, emojis[i])
	}

	defer ctx.Ses.AddHandler(func(_ *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		lock.RLock()
		defer lock.RUnlock()

		idx = framework.Index(emojis, reaction.Emoji.Name)

		if reaction.MessageID == embedMsg.ID && reaction.UserID != ctx.Ses.State.User.ID && idx != -1 {
			callback <- struct{}{}
		}
	})()

	embed = &discordgo.MessageEmbed{
		Author: &discordgo.MessageEmbedAuthor{},
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Answer",
				Value:  html.UnescapeString(question.CorrectAnswer),
				Inline: true,
			},
		},
	}

	select {
	case <-callback:
		if answers[idx] == question.CorrectAnswer {
			embed.Title = "Correct"
			embed.Color = 0x4caf50
			ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
			return
		} else {
			embed.Title = "Incorrect"
			embed.Color = 0xf44336
			ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
			return
		}
	case <-time.After(config.Timeout * time.Second):
		embed.Title = "Timed out"
		embed.Color = 0xfdd835
		ctx.Ses.ChannelMessageSendEmbed(ctx.Msg.ChannelID, embed)
		return
	}
}
