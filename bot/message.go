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

type sender struct {
  ID string `json:"id,omitempty"`
}

type coordinates struct {
  Lat  float64 `json:"lat,omitempty"`
  Long float64 `json:"long,omitempty"`
}

type location struct {
  Coordinates coordinates `json:"coordinates,omitempty"`
}

type ListElement struct {
  Title         string      `json:"title,omitempty"`
  Subtitle      string      `json:"subtitle,omitempty"`
  DefaultAction URLButton   `json:"default_action,omitempty"`
  Buttons       []URLButton `json:"buttons,omitempty"`
}

type Payload struct {
  URL              string        `json:"url,omitempty"`
  Type             string        `json:"type,omitempty"`
  Text             string        `json:"text,omitempty"`
  Buttons          []URLButton   `json:"buttons,omitempty"`
  Title            string        `json:"title,omitempty"`
  Location         interface{}   `json:"location,omitempty"`
  TemplateType     string        `json:"template_type,omitempty"`
  TopElementStyle  string        `json:"top_element_style,omitempty"`
  Sharable         bool          `json:"sharable,omitempty"`
  ImageAspectRatio string        `json:"image_aspect_ratio,omitempty"`
  Elements         []ListElement `json:"elements,omitempty"`
}

type Attachment struct {
  Type    string  `json:"type,omitempty"`
  Title   string  `json:"title,omitempty"`
  URL     string  `json:"url,omitempty"`
  Payload Payload `json:"payload,omitempty"`
}

type settingsMessage struct {
  SettingType   string       `json:"setting_type,omitempty"`
  ThreadState   string       `json:"thread_state,omitempty"`
  CallToActions []QuickReply `json:"call_to_actions,omitempty"`
}

type messaging struct {
  Recipient Recipient `json:"recipient,omitempty"`
  Sender    sender    `json:"sender,omitempty"`
  Timestamp int64     `json:"timestamp,omitempty"`
  Message struct {
    Mid         string       `json:"mid,omitempty"`
    Seq         int64        `json:"seq,omitempty"`
    Text        string       `json:"text,omitempty"`
    QuickReply  QuickReply   `json:"quick_reply,omitempty"`
    Attachments []Attachment `json:"attachments,omitempty"`
    IsEcho      bool         `json:"is_echo,omitempty"`
    AppID       string       `json:"app_id,omitempty"`
    Metadata    string       `json:"metadata,omitempty"`
  } `json:"message,omitempty"`
  Delivery struct {
    Mids      []string `json:"mids,omitempty"`
    Watermark int64    `json:"watermark,omitempty"`
    Seq       int64    `json:"seq,omitempty"`
  } `json:"delivery,omitempty"`
  Read struct {
    Watermark int64 `json:"watermark,omitempty"`
    Seq       int64 `json:"seq,omitempty"`
  } `json:"read,omitempty"`
  Postback struct {
    Payload string `json:"payload,omitempty"`
    Referral struct {
      Ref    string `json:"ref,omitempty"`
      AdID   string `json:"ad_id,omitempty"`
      Source string `json:"source,omitempty"`
      Type   string `json:"type,omitempty"`
    } `json:"referral,omitempty"`
  } `json:"postback,omitempty"`
  Optin struct {
    Ref string `json:"ref,omitempty"`
  } `json:"optin,omitempty"`
  AccountLinking struct {
    Status            string `json:"status,omitempty"`
    AuthorizationCode string `json:"authorization_code,omitempty"`
  } `json:"account_linking,omitempty"`
}

type entry struct {
  ID        string      `json:"id,omitempty"`
  Time      int64       `json:"time,omitempty"`
  Messaging []messaging `json:"messaging,omitempty"`
}

// Recipient represents the entity receiving the message
type Recipient struct {
  ID          string `json:"id,omitempty"`
  PhoneNumber string `json:"phone_number,omitempty"`
  Name struct {
    FirstName string `json:"first_name,omitempty"`
    LastName  string `json:"last_name,omitempty"`
  } `json:"name,omitempty"`
}

// QuickReply is used to present different actions to the user
type QuickReply struct {
  Title       string `json:"title,omitempty"`
  ContentType string `json:"content_type,omitempty"`
  ImageURL    string `json:"image_url,omitempty"`
  Payload     string `json:"payload,omitempty"`
}

// URLButton is used to present different paths to the user
type URLButton struct {
  Type                string      `json:"type,omitempty"`
  Title               string      `json:"title,omitempty"`
  URL                 string      `json:"url,omitempty"`
  FallbackURL         string      `json:"fallback_url,omitempty"`
  Payload             string      `json:"payload,omitempty"`
  MessengerExtensions bool        `json:"messenger_extensions,omitempty"`
  WebviewHeightRatio  string      `json:"webview_height_ratio,omitempty"`
  WebviewShareButton  string      `json:"webview_share_button,omitempty"`
  CallToActions       []URLButton `json:"call_to_actions,omitempty"`
}

// Persistent menu main structure
type Menu struct {
  Locale                string      `json:"locale,omitempty"`
  ComposerInputDisabled bool        `json:"composer_input_disabled,omitempty"`
  CallToActions         []URLButton `json:"call_to_actions,omitempty"`
}

// Messages received from the messenger platform
type Callback struct {
  Object string  `json:"object,omitempty"`
  Entry  []entry `json:"entry,omitempty"`
}

// MessageBody defines the main content section of a dispatched message
type MessageBody struct {
  Text         string       `json:"text,omitempty"`
  Metadata     string       `json:"metadata,omitempty"`
  QuickReplies []QuickReply `json:"quick_replies,omitempty"`
  Attachment   interface{}  `json:"attachment,omitempty"`
}

// Message dispatched to the messenger platform
type Message struct {
  Recipient        Recipient   `json:"recipient,omitempty"`
  SenderAction     string      `json:"sender_action,omitempty"`
  NotificationType string      `json:"notification_type,omitempty"`
  Message          interface{} `json:"message,omitempty"`
}
