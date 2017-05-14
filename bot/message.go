package bot

// Basic structures
type sender struct {
  ID string `json:"id"`
}

type recipient struct {
  ID string `json:"id"`
}

type quickReply struct {
  Payload string `json:"payload"`
}

type location struct {
  Coordinates struct {
    Lat  float64 `json:"lat"`
    Long float64 `json:"long"`
  } `json:"coordinates"`
}

type attachmentPayload struct {
  URL      string   `json:"url"`
  Location location `json:"location"`
}

type attachment struct {
  Type    string            `json:"type"`
  Payload attachmentPayload `json:"payload"`
}

type delivery struct {
  Mids      []string `json:"mids"`
  Watermark int64    `json:"watermark"`
  Seq       int64    `json:"seq"`
}

// MessageReceived represents an incoming message notification
// Endpoint: message
type MessageReceived struct {
  Sender    sender    `json:"sender"`
  Recipient recipient `json:"recipient"`
  Timestamp int64     `json:"timestamp"`
  Message struct {
    Mid         string       `json:"mid"`
    Seq         int64        `json:"seq"`
    Text        string       `json:"text"`
    QuickReply  quickReply   `json:"quick_reply"`
    Attachments []attachment `json:"attachments"`
  } `json:"message"`
}

// MessageDelivered represents an outgoing message notification
// Endpoint: message_deliveries
type MessageDelivered struct {
  Sender    sender     `json:"sender"`
  Recipient recipient  `json:"recipient"`
  Delivery  []delivery `json:"delivery"`
}

// MessageRead represents a 'messages read' notification
// Endpoint: message_reads
type MessageRead struct {
  Sender    sender    `json:"sender"`
  Recipient recipient `json:"recipient"`
  Timestamp int64     `json:"timestamp"`
  Read struct {
    Watermark int64 `json:"watermark"`
    Seq       int64 `json:"seq"`
  } `json:"read"`
}

// MessagePostBack represents a user action
// Endpoint: messaging_postbacks
type MessagePostBack struct {
  Sender    sender    `json:"sender"`
  Recipient recipient `json:"recipient"`
  Timestamp int64     `json:"timestamp"`
  Postback struct {
    Payload string `json:"payload"`
  } `json:"postback"`
}

// AccountLinking represents an authorization action
// Endpoint: messaging_postbacks
type AccountLinking struct {
  Sender    sender    `json:"sender"`
  Recipient recipient `json:"recipient"`
  Timestamp int64     `json:"timestamp"`
  Link struct {
    Status            string `json:"status"`
    AuthorizationCode string `json:"authorization_code"`
  } `json:"account_linking"`
}
