package discord

import (
	"fmt"

	"github.com/coma-toast/notifapi/internal/utils"
	"github.com/disgoorg/disgo/discord"
	"github.com/disgoorg/disgo/rest"
	"github.com/disgoorg/disgo/webhook"
	"github.com/disgoorg/log"
)

type Discord struct {
	URL  string
	Data *utils.DataModel
}

func (d Discord) SendMessage(interests []string, title, body, source string) (string, error) {
	return d.SendMessageFull(interests, title, body, "", source, nil)

}
func (d Discord) SendMessageWithLink(interests []string, title, body, link, source string) (string, error) {
	return d.SendMessageFull(interests, title, body, link, source, nil)
}
func (d Discord) SendMessageFull(interests []string, title, body, link, source string, metadata map[string]interface{}) (string, error) {
	client, err := webhook.NewWithURL(d.URL)
	if err != nil {
		fmt.Println(err)
	}

	message := fmt.Sprintf("%s: %s - %s", source, title, body)

	messageId, err := send(client, message, link)
	if err != nil {
		return "", err
	}

	d.Data.AddNotification(messageId, source, "discord", title, body, interests, metadata)

	return messageId, nil
}

func send(client webhook.Client, payload string, url string) (string, error) {
	var message discord.WebhookMessageCreate
	if url != "" {
		embedData := discord.NewEmbedBuilder()
		embedData.SetURL(url)
		message = discord.NewWebhookMessageCreateBuilder().SetContent(payload).AddEmbeds(embedData.Embed).Build()
	} else {
		message = discord.NewWebhookMessageCreateBuilder().SetContent(payload).Build()
	}
	results, err := client.CreateMessage(message, rest.WithDelay(0))
	if err != nil {
		log.Errorf("error sending message %d: %s", err)
		return "", err
	}

	return results.ID.String(), nil
}
