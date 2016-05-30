package vk

import (
	"encoding/json"
)

// Messages https://new.vk.com/dev/messages
type Messages struct {
	client *VK
}

// Send https://new.vk.com/dev/messages.send
func (messages *Messages) Send(params RequestParams) (int64, error) {
	resp, err := messages.client.CallMethod("messages.send", params)

	if err != nil {
		return 0, err
	}

	type JSONBody struct {
		MessageID int64 `json:"response"`
	}

	var body JSONBody

	if err := json.Unmarshal(resp, &body); err != nil {
		return 0, err
	}

	return body.MessageID, nil
}
