package notifications

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ses"
)

// EmailOptions basic parameters required to dispatch an email notification
type EmailOptions struct {
	To      string `json:"to"`
	From    string `json:"from"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
}

// SendEmail will dispatch an email notification
func SendEmail(opts *EmailOptions) (string, error) {
	sess, err := session.NewSession()
	if err != nil {
		aerr, _ := err.(awserr.Error)
		return "", aerr.OrigErr()
	}

	svc := ses.New(sess)
	params := &ses.SendEmailInput{
		Destination: &ses.Destination{
			ToAddresses: []*string{
				aws.String(opts.To),
			},
		},
		Message: &ses.Message{
			Body: &ses.Body{
				Html: &ses.Content{
					Data:    aws.String(opts.Body),
					Charset: aws.String("utf-8"),
				},
			},
			Subject: &ses.Content{
				Data:    aws.String(opts.Subject),
				Charset: aws.String("utf-8"),
			},
		},
		Source: aws.String(opts.From),
		ReplyToAddresses: []*string{
			aws.String(opts.From),
		},
	}

	resp, err := svc.SendEmail(params)
	if err != nil {
		aerr, _ := err.(awserr.Error)
		return "", aerr.OrigErr()
	}
	return *resp.MessageId, nil
}
