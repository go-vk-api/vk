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
		vk.WithHTTPClient(http.DefaultClient),
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

	for {
		select {
		case update, ok := <-stream.Updates:
			if !ok {
				return
			}

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
		case err, ok := <-stream.Errors:
			if ok {
				// stream.Stop()
				log.Panic(err)
			}
		}
	}
}

func printMe(api *vk.Client) error {
	var users []struct {
		ID        int64  `json:"id"`
		FirstName string `json:"first_name"`
		LastName  string `json:"last_name"`
	}

	if err := api.CallMethod("users.get", vk.RequestParams{}, &users); err != nil {
		return err
	}

	me := users[0]

	log.Println(me.ID, me.FirstName, me.LastName)

	return nil
}
