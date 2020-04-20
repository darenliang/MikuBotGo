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
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "An error has occurred")
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
		_ = ctx.Ses.MessageReactionAdd(ctx.Msg.ChannelID, embedMsg.ID, emojis[i])
	}

	defer ctx.Ses.AddHandler(func(_ *discordgo.Session, reaction *discordgo.MessageReactionAdd) {
		lock.RLock()
		defer lock.RUnlock()

		idx = framework.Index(emojis, reaction.Emoji.Name)

		if reaction.MessageID == embedMsg.ID && reaction.UserID != ctx.Ses.State.User.ID && idx != -1 {
			close(callback)
		}
	})()

	select {
	case <-callback:
		if answers[idx] == question.CorrectAnswer {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, "You are correct!")
			return
		} else {
			_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID,
				fmt.Sprintf("You are incorrect. The answer is %s.", html.UnescapeString(question.CorrectAnswer)))
			return
		}
	case <-time.After(config.Timeout * time.Second):
		_, _ = ctx.Ses.ChannelMessageSend(ctx.Msg.ChannelID, fmt.Sprintf(
			"The trivia question timed out. The answer is %s.", html.UnescapeString(question.CorrectAnswer)))
		return
	}
}