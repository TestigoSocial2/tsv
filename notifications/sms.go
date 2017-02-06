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
