package vk

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/go-resty/resty"
)

// VK struct
type VK struct {
	url     string
	lang    string
	version string
	token   string

	logFile  *os.File
	longPoll *longPoll

	Messages *Messages

	Proxy string
}

// RequestParams struct
type RequestParams map[string]string

// CallMethod calls VK API method
func (client *VK) CallMethod(method string, params RequestParams) ([]byte, error) {
	params["access_token"] = client.token
	params["lang"] = client.lang
	params["v"] = client.version

	if client.Proxy != "" {
		resty.SetProxy(fmt.Sprint("http://", client.Proxy))
	}

	resp, err := resty.R().
		SetQueryParams(params).
		Get(client.url + method)

	if err != nil {
		client.Log("[Error] VK::CallMethod:", err.Error(), "WebResponse:", string(resp.Body()))
		return nil, err
	}

	type JSONBody struct {
		Error map[string]interface{} `json:"error"`
	}

	var body JSONBody

	if err := json.Unmarshal(resp.Body(), &body); err != nil {
		client.Log("[Error] VK::CallMethod:", err.Error(), "WebResponse:", string(resp.Body()))
		return nil, err
	}

	if body.Error != nil {
		if errorMsg, exists := body.Error["error_msg"].(string); exists {
			client.Log("[Error] VK::CallMethod:", errorMsg, "WebResponse:", string(resp.Body()))
			return nil, errors.New(errorMsg)
		}

		client.Log("[Error] VK::CallMethod:", "Unknown error", "WebResponse:", string(resp.Body()))
		return nil, errors.New("Unknown error")
	}

	return resp.Body(), nil
}

// Init sets the token
func (client *VK) Init(token string) error {
	client.token = token

	return nil
}

// RunLongPoll starts longpoll process
func (client *VK) RunLongPoll() {
	if err := client.longPoll.update(); err != nil {
		client.Log("[Error] VK::RunLongPoll:", err.Error())
		return
	}

	client.longPoll.chanNewMessage = make(chan *LPMessage)

	go client.longPoll.processEvents()

	client.longPoll.process()
}

// OnNewMessage sets event
func (client *VK) OnNewMessage(event EventNewMessage) {
	client.longPoll.eventNewMessage = event
}

// SetLogFile sets pointer to logfile
func (client *VK) SetLogFile(logFile *os.File) {
	client.logFile = logFile
}

// Log writes data in stdout and logfile
func (client *VK) Log(a ...interface{}) {
	log.SetFlags(log.LstdFlags)
	log.SetOutput(os.Stdout)
	log.Println(a...)

	if client.logFile != nil {
		log.SetOutput(client.logFile)
		log.Println(a...)
	}
}

// New returns a new VK instance
func New(lang string) *VK {
	vk := new(VK)

	vk.url = "https://api.vk.com/method/"
	vk.lang = lang
	vk.version = "5.52"

	vk.longPoll = &longPoll{client: vk}
	vk.Messages = &Messages{client: vk}

	return vk
}
