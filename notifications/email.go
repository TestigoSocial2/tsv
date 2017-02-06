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
