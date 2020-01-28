# go-vk-api
[![GoDoc](https://godoc.org/github.com/go-vk-api/vk?status.svg)](https://godoc.org/github.com/go-vk-api/vk)
[![Go Report Card](https://goreportcard.com/badge/github.com/go-vk-api/vk)](https://goreportcard.com/report/github.com/go-vk-api/vk)

Golang bindings for the VK API

## Install

Install the package with:

```bash
go get github.com/go-vk-api/vk
```

Import it with:

```go
import "github.com/go-vk-api/vk"
```

and use `vk` as the package name inside the code.

## Example

[Full example with errors handling](https://github.com/go-vk-api/vk/blob/master/example/example.go)

```go
package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-vk-api/vk"
	lp "github.com/go-vk-api/vk/longpoll/user"
)

func main() {
	client, _ := vk.NewClientWithOptions(
		vk.WithToken(os.Getenv("VK_ACCESS_TOKEN")),
	)

	_ = printMe(client)

	longpoll, _ := lp.NewWithOptions(client, lp.WithMode(lp.ReceiveAttachments))

	stream, _ := longpoll.GetUpdatesStream(0)

	for update := range stream.Updates {
		switch data := update.Data.(type) {
		case *lp.NewMessage:
			if data.Text == "/hello" {
				_ = client.CallMethod("messages.send", vk.RequestParams{
					"peer_id":          data.PeerID,
					"message":          "Hello!",
					"forward_messages": data.ID,
					"random_id":        0,
				}, nil)
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

	_ = api.CallMethod("users.get", vk.RequestParams{}, &users)

	me := users[0]

	log.Println(me.ID, me.FirstName, me.LastName)

	return nil
}
```
