# go-vk-api
[![GoDoc](https://godoc.org/github.com/urShadow/go-vk-api?status.svg)](https://godoc.org/github.com/urShadow/go-vk-api)
[![Go Report Card](https://goreportcard.com/badge/github.com/urShadow/go-vk-api)](https://goreportcard.com/report/github.com/urShadow/go-vk-api)

Golang wrapper for VK API

## Install

Install the package with:

```bash
go get github.com/urShadow/go-vk-api
```

Import it with:

```go
import "github.com/urShadow/go-vk-api"
```

and use `vk` as the package name inside the code.

## Example

```go
package main

import (
	"go-vk-api"
	"log"
	"strconv"
)

func main() {
	api := vk.New("ru")

	err := api.Init("TOKEN")

	if err != nil {
		log.Fatalln(err)
	}

	api.OnNewMessage(func(msg *vk.LPMessage) {
		if msg.Flags&vk.FlagMessageOutBox == 0 {
			api.Messages.Send(vk.RequestParams{
				"peer_id":          strconv.FormatInt(msg.FromID, 10),
				"message":          "Hello, World!",
				"forward_messages": strconv.FormatInt(msg.ID, 10),
			})
		}
	})

	api.RunLongPoll()
}

```
