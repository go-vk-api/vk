package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-vk-api/vk"
	lp "github.com/go-vk-api/vk/longpoll/user"
)

func main() {
	client, err := vk.NewClientWithOptions(
		vk.WithToken(os.Getenv("VK_ACCESS_TOKEN")),
		vk.WithHttpClient(http.DefaultClient),
	)
	if err != nil {
		log.Panic(err)
	}

	err = printMe(client)
	if err != nil {
		log.Panic(err)
	}

	longpoll, err := lp.NewWithOptions(client, lp.WithMode(lp.ReceiveAttachments))
	if err != nil {
		log.Panic(err)
	}

	stream, err := longpoll.GetUpdatesStream(0)
	if err != nil {
		log.Panic(err)
	}

	for update := range stream.Updates {
		switch data := update.Data.(type) {
		case *lp.NewMessage:
			if data.Text == "/hello" {
				var sentMessageID int64

				if err = client.CallMethod("messages.send", vk.RequestParams{
					"peer_id":          data.PeerID,
					"message":          "Hello!",
					"forward_messages": data.ID,
					"random_id":        0,
				}, &sentMessageID); err != nil {
					log.Panic(err)
				}

				log.Println(sentMessageID)
			}
		}
	}
}

func printMe(api *vk.Client) error {
	var users []struct {
		Id        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := api.CallMethod("users.get", vk.RequestParams{}, &users); err != nil {
		return err
	}

	me := users[0]

	log.Println(me.Id, me.FirstName, me.LastName)

	return nil
}
