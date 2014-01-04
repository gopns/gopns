package sqs

type SqsQueue struct {
	QueueUrl string `xml:"CreateQueueResult>QueueUrl"`
}

type SqsSendMessageResponse struct {
	MessageId string `xml:"SendMessageResult>MessageId"`
	BodyHash  string `xml:"SendMessageResult>MD5OfMessageBody"`
}
