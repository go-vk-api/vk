package longpoll

import (
	"encoding/json"
	"fmt"

	"github.com/go-vk-api/vk"
	"github.com/go-vk-api/vk/httputil"
)

const (
	// DefaultVersion is a default version of the VK Long Poll API.
	DefaultVersion = 2
	// DefaultWait is a waiting period. Maximum: 90.
	DefaultWait = 25
	// DefaultMode is an additional answer options.
	DefaultMode = ReceiveAttachments
)

// Mode represents the additional answer options.
type Mode int64

const (
	ReceiveAttachments           Mode = 2
	ReturnExpandedSetOfEvents    Mode = 8
	ReturnPts                    Mode = 32
	ReturnFriendOnlineExtraField Mode = 64
	ReturnRandomID               Mode = 128
)

const (
	eventHistoryOutdated = iota + 1
	keyExpired
	userInformationLost
	invalidVersion
)

// Longpoll manages communication with VK User Long Poll API.
type Longpoll struct {
	client *vk.Client

	Key     string
	Server  string
	Wait    int64
	Mode    Mode
	Version int64

	NeedPts int64
	GroupID int64
}

// UpdateServer updates the longpoll server and returns a new ts.
func (lp *Longpoll) UpdateServer() (newTs int64, err error) {
	params := vk.RequestParams{
		"need_pts":   lp.NeedPts,
		"lp_version": lp.Version,
	}

	if lp.GroupID > 0 {
		params["group_id"] = lp.GroupID
	}

	var body struct {
		Key    string `json:"key"`
		Server string `json:"server"`
		Ts     int64  `json:"ts"`
	}

	if err = lp.client.CallMethod("messages.getLongPollServer", params, &body); err != nil {
		return
	}

	lp.Key = body.Key
	lp.Server = body.Server

	return body.Ts, nil
}

// Poll requests updates starting with the ts and returns a new updates and ts.
func (lp *Longpoll) Poll(ts int64) (updates []*Update, newTS int64, err error) {
	params, err := vk.RequestParams{
		"act":     "a_check",
		"key":     lp.Key,
		"ts":      ts,
		"wait":    lp.Wait,
		"mode":    lp.Mode,
		"version": lp.Version,
	}.URLValues()
	if err != nil {
		return
	}

	rawBody, err := httputil.Post(lp.client.HTTPClient, "https://"+lp.Server, params)
	if err != nil {
		return
	}

	var body struct {
		TS      int64     `json:"ts"`
		Updates []*Update `json:"updates"`
		Failed  *int64    `json:"failed"`
	}

	if err = json.Unmarshal(rawBody, &body); err != nil {
		return
	}

	if body.Failed != nil {
		return nil, body.TS, lp.failedToError(*body.Failed)
	}

	return body.Updates, body.TS, nil
}

func (lp *Longpoll) failedToError(failed int64) error {
	switch failed {
	case eventHistoryOutdated:
		return ErrEventHistoryOutdated
	case keyExpired:
		return ErrKeyExpired
	case userInformationLost:
		return ErrUserInformationLost
	case invalidVersion:
		return ErrInvalidVersion
	}

	return fmt.Errorf("unexpected failed value (%d)", failed)
}

// GetUpdatesStream starts and returns a stream of updates.
func (lp *Longpoll) GetUpdatesStream(ts int64) (*Stream, error) {
	newTs, err := lp.UpdateServer()
	if err != nil {
		return nil, err
	}

	if ts == 0 {
		ts = newTs
	}

	stream := &Stream{
		lp: lp,
		TS: ts,
	}

	if err := stream.Start(); err != nil {
		return nil, err
	}

	return stream, nil
}

// New initializes a new longpoll client with default values.
func New(client *vk.Client) (*Longpoll, error) {
	return NewWithOptions(client)
}

// NewWithOptions initializes a new longpoll client with default values. It takes functors
// to modify values when creating it, like `NewWithOptions(WithMode(…))`.
func NewWithOptions(client *vk.Client, options ...Option) (*Longpoll, error) {
	longpoll := &Longpoll{
		client: client,

		Wait:    DefaultWait,
		Mode:    DefaultMode,
		Version: DefaultVersion,
	}

	for _, option := range options {
		if err := option(longpoll); err != nil {
			return nil, err
		}
	}

	return longpoll, nil
}
