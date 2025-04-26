package discord

import (
	"fmt"

	"github.com/coma-toast/notifapi/backend/pkg/notification"
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

func (d Discord) SendMessage(payload notification.Message) (string, error) {

	client, err := webhook.NewWithURL(d.URL)
	if err != nil {
		fmt.Println(err)
	}

	message := fmt.Sprintf("%s: %s - %s", payload.Server, payload.Title, payload.Body)
	link := payload.Link

	messageId, err := send(client, message, link)
	if err != nil {
		return "", err
	}

	d.Data.AddNotification(payload)

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
		log.Errorf("error sending message %d: %s", results.ID, err)
		return "", err
	}

	return results.ID.String(), nil
}
