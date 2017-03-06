package bot

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/Jeffail/gabs"
	"github.com/julienschmidt/httprouter"
)

// Bot provides the functionality to implement a Facebook Messenger agent
type Bot struct {
	verificationToken string
	pageToken         string
	IncomingMessages  chan *IncomingMessage
}

// Config defines the required parameters to start the bot
type Config struct {
	VerificationToken string
	PageToken         string
}

// IncomingMessage are messages received from Facebook users
type IncomingMessage struct {
	User    string
	Content string
}

// New will create and return a bot instance
func New(c *Config) *Bot {
	return &Bot{
		verificationToken: c.VerificationToken,
		pageToken:         c.PageToken,
		IncomingMessages:  make(chan *IncomingMessage, 10),
	}
}

// Verify handle the endpoint verification process
func (b *Bot) Verify(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
	mode := r.FormValue("hub.mode")
	token := r.FormValue("hub.verify_token")
	if mode == "subscribe" && token == b.verificationToken {
		fmt.Fprintf(w, fmt.Sprintf("%s", r.FormValue("hub.challenge")))
		return
	}

	fmt.Fprintf(w, fmt.Sprintf("%s", "INVALID_VERIFICATION_TOKEN"))
}

// ReceiveMessage ...
func (b *Bot) ReceiveMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Get request body and return response
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("error: %s", err)
	}
	fmt.Fprintf(w, "%s", "")

	// Parse and process message contents
	data, _ := gabs.ParseJSON(body)
	if data.Path("object").Data().(string) == "page" {
		entries, _ := data.Search("entry").Children()
		for _, entry := range entries {
			messages, _ := entry.Search("messaging").Children()
			for _, m := range messages {
				b.IncomingMessages <- &IncomingMessage{
					User:    m.Path("sender.id").Data().(string),
					Content: m.Path("message.text").Data().(string),
				}
			}
		}
	}
}

// DispatchMessage use Facebook's Messenger to reach users
func (b *Bot) DispatchMessage(recipient string, content string) error {
	// Build response
	msgResponse := gabs.New()
	msgResponse.Set(recipient, "recipient", "id")
	msgResponse.Set(content, "message", "text")

	// Dispatch response
	url := fmt.Sprintf("%s=%s", "https://graph.facebook.com/v2.6/me/messages?access_token", b.pageToken)
	req, _ := http.NewRequest("POST", url, bytes.NewBuffer([]byte(msgResponse.String())))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return nil
}
