package vk

import (
	"encoding/json"
	"github.com/go-resty/resty"
	"strconv"
)

const (
	longPollNewMessage = 4
)

// Message's flags
const (
	FlagMessageUnread = 1 << iota
	FlagMessageOutBox
	FlagMessageReplied
	FlagMessageImportant
	FlagMessageChat
	FlagMessageFriends
	FlagMessageSpam
	FlagMessageDeleted
	FlagMessageFixed
	FlagMessageMedia
)

// LPMessage struct
type LPMessage struct {
	ID, Flags, FromID, Timestamp int64
	Subject, Text                string
	Attachments                  map[string]string
}

// EventNewMessage delegate
type EventNewMessage func(*LPMessage)

type longPoll struct {
	client *VK

	chanNewMessage  chan *LPMessage
	eventNewMessage EventNewMessage

	data struct {
		server string
		key    string
		ts     int64
	}
}

func (lp *longPoll) update() error {
	resp, err := lp.client.CallMethod("messages.getLongPollServer", RequestParams{
		"use_ssl":  "0",
		"need_pts": "0",
	})

	if err != nil {
		return err
	}

	type JSONBody struct {
		Response struct {
			Server string `json:"server"`
			Key    string `json:"key"`
			Ts     int64  `json:"ts"`
		} `json:"response"`
	}

	var body JSONBody

	if err := json.Unmarshal(resp, &body); err != nil {
		return err
	}

	lp.data.server = body.Response.Server
	lp.data.key = body.Response.Key
	lp.data.ts = body.Response.Ts

	return nil
}

func (lp *longPoll) process() {
	resp, err := resty.R().
		SetQueryParams(RequestParams{
			"act":  "a_check",
			"key":  lp.data.key,
			"ts":   strconv.FormatInt(lp.data.ts, 10),
			"wait": "25",
			"mode": "2",
		}).
		Get("https://" + lp.data.server)

	if err != nil {
		lp.client.Log("[Error] longPoll::process:", err.Error(), "WebResponse:", string(resp.Body()))
		return
	}

	type jsonBody struct {
		Failed  int64           `json:"failed"`
		Ts      int64           `json:"ts"`
		Updates [][]interface{} `json:"updates"`
	}

	var body jsonBody

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		lp.client.Log("[Error] longPoll::process:", err.Error(), "WebResponse:", string(resp.Body()))
		return
	}

	switch body.Failed {
	case 0:
		for _, update := range body.Updates {
			updateID := update[0].(float64)

			switch updateID {
			case longPollNewMessage:
				message := new(LPMessage)

				message.ID = int64(update[1].(float64))
				message.Flags = int64(update[2].(float64))
				message.FromID = int64(update[3].(float64))
				message.Timestamp = int64(update[4].(float64))
				message.Subject = update[5].(string)
				message.Text = update[6].(string)
				message.Attachments = make(map[string]string)

				for key, value := range update[7].(map[string]interface{}) {
					message.Attachments[key] = value.(string)
				}

				lp.chanNewMessage <- message
			}
		}

		lp.data.ts = body.Ts
	case 1:
		lp.data.ts = body.Ts
		lp.client.Log("ts updated")
	case 2, 3:
		if err := lp.update(); err != nil {
			lp.client.Log("Longpoll update error:", err.Error())
			return
		}
		lp.client.Log("Longpoll data updated")
	}

	lp.process()
}

func (lp *longPoll) processEvents() {
	for {
		select {
		case message := <-lp.chanNewMessage:
			if lp.eventNewMessage != nil {
				lp.eventNewMessage(message)
			}
		}
	}
}
