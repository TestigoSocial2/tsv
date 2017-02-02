package notifications

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
)

// SMSOptions defines the options required to dispatch a SMS notification
type SMSOptions struct {
	To      string `json:"to"`
	Sender  string `json:"sender"`
	Message string `json:"message"`
}

// SendSMS will dispatch an SMS notification
func SendSMS(opts *SMSOptions) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		aerr, _ := err.(awserr.Error)
		return "", aerr.OrigErr()
	}

	svc := sns.New(sess)
	params := &sns.PublishInput{
		Message:     aws.String(opts.Message),
		PhoneNumber: aws.String(opts.To),
		MessageAttributes: map[string]*sns.MessageAttributeValue{
			"DefaultSenderID": {
				DataType:    aws.String("String"),
				StringValue: aws.String(opts.Sender),
			},
		},
	}

	resp, err := svc.Publish(params)
	if err != nil {
		aerr, _ := err.(awserr.Error)
		return "", aerr.OrigErr()
	}
	return *resp.MessageId, nil
}
