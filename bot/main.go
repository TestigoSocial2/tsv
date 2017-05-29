// Copyright © 2016 Transparencia Mexicana AC. <ben@pixative.com>
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

package bot

import (
  "bytes"
  "fmt"
  "io/ioutil"
  "net/http"
  "encoding/json"
)

// Instance provides the functionality to implement a Facebook Messenger agent
type Instance struct {
  verificationToken string
  pageToken         string
  IncomingMessages  chan *Callback
  Errors            chan error
}

// Config defines the required parameters to start the bot
type Config struct {
  VerificationToken string
  PageToken         string
}

type dispatchOptions struct {
  method   string
  endpoint string
  data     []byte
}

// New will create and return a bot instance
func New(c *Config) *Instance {
  return &Instance{
    verificationToken: c.VerificationToken,
    pageToken:         c.PageToken,
    IncomingMessages:  make(chan *Callback, 10),
    Errors:            make(chan error),
  }
}

// Utility method to properly dispatch a request to FB's API
func (b *Instance) dispatch(opts *dispatchOptions) ([]byte, error) {
  // Build request
  url := fmt.Sprintf("https://graph.facebook.com/v2.6/me/%s?access_token=%s", opts.endpoint, b.pageToken)
  req, _ := http.NewRequest(opts.method, url, bytes.NewBuffer(opts.data))
  req.Header.Set("Content-Type", "application/json")
  
  // Send request
  client := &http.Client{}
  resp, err := client.Do(req)
  if err != nil {
    return nil, err
  }
  defer resp.Body.Close()
  
  // Return body
  return ioutil.ReadAll(resp.Body)
}

// Verify handle the endpoint for the verification process
func (b *Instance) Verify(w http.ResponseWriter, r *http.Request) {
  w.Header().Add("Strict-Transport-Security", "max-age=63072000; includeSubDomains")
  mode := r.FormValue("hub.mode")
  token := r.FormValue("hub.verify_token")
  if mode == "subscribe" && token == b.verificationToken {
    fmt.Fprintf(w, fmt.Sprintf("%s", r.FormValue("hub.challenge")))
    return
  }
  
  fmt.Fprintf(w, fmt.Sprintf("%s", "INVALID_VERIFICATION_TOKEN"))
}

// ReceiveMessage handle the MessageBody delivery endpoint
func (b *Instance) ReceiveMessage(w http.ResponseWriter, r *http.Request) {
  // Get request body and return response
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    b.Errors <- err
  }
  
  var cb Callback
  err = json.Unmarshal(body, &cb)
  if err != nil {
    b.Errors <- err
  }
  
  b.IncomingMessages <- &cb
}

// SendMessage use Facebook's Messenger to dispatch a message
func (b *Instance) SendMessage(msg *Message) ([]byte, error) {
  content, _ := json.Marshal(msg)
  return b.dispatch(&dispatchOptions{
    method:   "post",
    endpoint: "messages",
    data:     []byte(content),
  })
}

// SendTextMessage use Facebook's Messenger to send a basic text MessageBody
func (b *Instance) SendTextMessage(user, text string) ([]byte, error) {
  return b.SendMessage(&Message{
    Recipient: Recipient{
      ID: user,
    },
    Message: MessageBody{
      Text: text,
    },
  })
}

// SendImageMessage use Facebook's Messenger to send an image
func (b *Instance) SendImageMessage(user, url string) ([]byte, error) {
  return b.SendMessage(&Message{
    Recipient: Recipient{
      ID: user,
    },
    Message: MessageBody{
      Attachment: Attachment{
        Type: "image",
        Payload: Payload{
          URL: url,
        },
      },
    },
  })
}

// SendAudioMessage use Facebook's Messenger to send an audio MessageBody
func (b *Instance) SendAudioMessage(user, url string) ([]byte, error) {
  return b.SendMessage(&Message{
    Recipient: Recipient{
      ID: user,
    },
    Message: MessageBody{
      Attachment: Attachment{
        Type: "audio",
        Payload: Payload{
          URL: url,
        },
      },
    },
  })
}

// SendVideoMessage use Facebook's Messenger to send a video
func (b *Instance) SendVideoMessage(user, url string) ([]byte, error) {
  return b.SendMessage(&Message{
    Recipient: Recipient{
      ID: user,
    },
    Message: MessageBody{
      Attachment: Attachment{
        Type: "video",
        Payload: Payload{
          URL: url,
        },
      },
    },
  })
}

// SendFileMessage use Facebook's Messenger to send a file
func (b *Instance) SendFileMessage(user, url string) ([]byte, error) {
  return b.SendMessage(&Message{
    Recipient: Recipient{
      ID: user,
    },
    Message: MessageBody{
      Attachment: Attachment{
        Type: "file",
        Payload: Payload{
          URL: url,
        },
      },
    },
  })
}

// SetGetStartedButton will set the postback value used on the 'Get Started' button
func (b *Instance) SetGetStartedButton(payload string) ([]byte, error) {
  msg := settingsMessage{
    SettingType: "call_to_actions",
    ThreadState: "new_thread",
    CallToActions: []QuickReply{
      {Payload: payload },
    },
  }
  content, _ := json.Marshal(msg)
  return b.dispatch(&dispatchOptions{
    method:   "post",
    endpoint: "thread_settings",
    data:     content,
  })
}

// RemoveGetStartedButton will clear any previously set value for the 'Get Started' button
func (b *Instance) RemoveGetStartedButton() ([]byte, error) {
  msg := settingsMessage{
    SettingType: "call_to_actions",
    ThreadState: "new_thread",
  }
  content, _ := json.Marshal(msg)
  return b.dispatch(&dispatchOptions{
    method:   "post",
    endpoint: "thread_settings",
    data:     content,
  })
}

// SetPersistentMenu will configure the persistent menu used by the bot
func (b *Instance) SetPersistentMenu(menu *Menu) ([]byte, error) {
  content, _ := json.Marshal(menu)
  return b.dispatch(&dispatchOptions{
    method:   "post",
    endpoint: "messenger_profile",
    data:     content,
  })
}

// RemovePersistentMenu will clear any previously set persistent menu configuration
func (b *Instance) RemovePersistentMenu(menu *Menu) ([]byte, error) {
  content := []byte("{\"fields\":[\"persistent_menu\"]}")
  return b.dispatch(&dispatchOptions{
    method:   "delete",
    endpoint: "messenger_profile",
    data:     content,
  })
}
