package longpoll

import (
	"fmt"
	"reflect"
	"testing"
)

func TestUpdate_UnmarshalJSON(t *testing.T) {
	update := Update{}

	wantMessage := &NewMessage{
		ID:          1,
		Flags:       2,
		PeerID:      3,
		Timestamp:   4,
		Text:        "567",
		Attachments: map[string]string{"from": "89"},
		RandomID:    10,
	}

	json := fmt.Sprintf(
		`[%d,%d,%d,%d,%d,"%s",{"from":"%s"},%d]`,
		EventNewMessage,
		wantMessage.ID,
		wantMessage.Flags,
		wantMessage.PeerID,
		wantMessage.Timestamp,
		wantMessage.Text,
		wantMessage.Attachments["from"],
		wantMessage.RandomID,
	)

	err := update.UnmarshalJSON([]byte(json))
	if err != nil {
		t.Error(err)
	}

	if update.Type != EventNewMessage {
		t.Errorf("update.Type == %d, want %d", update.Type, EventNewMessage)
	}

	message, ok := update.Data.(*NewMessage)
	if !ok {
		t.Errorf("reflect.TypeOf(update.Data) == %v, want %v", reflect.TypeOf(update.Data), reflect.TypeOf(wantMessage))
	}

	if message.ID != wantMessage.ID {
		t.Errorf("message.ID == %d, want %d", message.ID, wantMessage.ID)
	}

	if message.Flags != wantMessage.Flags {
		t.Errorf("message.Flags == %d, want %d", message.Flags, wantMessage.Flags)
	}

	if message.PeerID != wantMessage.PeerID {
		t.Errorf("message.PeerID == %d, want %d", message.PeerID, wantMessage.PeerID)
	}

	if message.Timestamp != wantMessage.Timestamp {
		t.Errorf("message.Timestamp == %d, want %d", message.Timestamp, wantMessage.Timestamp)
	}

	if message.Text != wantMessage.Text {
		t.Errorf("message.Text == %q, want %q", message.Text, wantMessage.Text)
	}

	if message.Attachments["from"] != wantMessage.Attachments["from"] {
		t.Errorf("message.Attachments[\"from\"] == %q, want %q", message.Attachments["from"], wantMessage.Attachments["from"])
	}

	if message.RandomID != wantMessage.RandomID {
		t.Errorf("message.RandomID == %d, want %d", message.RandomID, wantMessage.RandomID)
	}
}
