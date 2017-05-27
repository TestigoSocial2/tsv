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
  
  "github.com/julienschmidt/httprouter"
  "encoding/json"
)

// Instance provides the functionality to implement a Facebook Messenger agent
type Instance struct {
  verificationToken string
  pageToken         string
  IncomingMessages  chan *Callback
  ReceptionErrors   chan error
}

// Config defines the required parameters to start the bot
type Config struct {
  VerificationToken string
  PageToken         string
}

// New will create and return a bot instance
func New(c *Config) *Instance {
  return &Instance{
    verificationToken: c.VerificationToken,
    pageToken:         c.PageToken,
    IncomingMessages:  make(chan *Callback, 10),
    ReceptionErrors:   make(chan error),
  }
}

// Utility method to properly dispatch a request to FB's API
func (b *Instance) dispatch(content []byte) ([]byte, error) {
  // Build request
  url := fmt.Sprintf("%s=%s", "https://graph.facebook.com/v2.6/me/messages?access_token", b.pageToken)
  req, _ := http.NewRequest("POST", url, bytes.NewBuffer(content))
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

// ReceiveMessage handle the message delivery endpoint
func (b *Instance) ReceiveMessage(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
  // Get request body and return response
  body, err := ioutil.ReadAll(r.Body)
  if err != nil {
    b.ReceptionErrors <- err
  }
  
  var cb Callback
  err = json.Unmarshal(body, &cb)
  if err != nil {
    b.ReceptionErrors <- err
  }
  
  b.IncomingMessages <- &cb
}

// DispatchMessage use Facebook's Messenger to reach users
func (b *Instance) DispatchMessage(msg string) ([]byte, error) {
  return b.dispatch([]byte(msg))
}
