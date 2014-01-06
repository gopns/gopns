package sqs

type SqsQueue struct {
	QueueUrl string `xml:"CreateQueueResult>QueueUrl"`
}

type SqsSendMessageResponse struct {
	MessageId string `xml:"SendMessageResult>MessageId"`
	BodyHash  string `xml:"SendMessageResult>MD5OfMessageBody"`
}

type SqsReceiveMessageResponse struct {
	SqsMessages []SqsMessage `xml:"ReceiveMessageResult>Message"`
}

type SqsDeleteMessagesResponse struct {
	DeletedMessageIds []string       `xml:"DeleteMessageBatchResult>DeleteMessageBatchResultEntry>Id"`
	Errors            []ErrorMessage `xml:"DeleteMessageBatchResult>BatchResultErrorEntry"`
}

type ErrorMessage struct {
	MessageId    string `xml:"Id"`
	ErrorMessage string `xml:"Message"`
	ServerFault  bool   `xml:"serverFault"`
	ErrorCode    string `xml:"Code"`
}

type SqsMessage struct {
	MessageId     string `xml:"MessageId"`
	Body          string `xml:"Body"`
	BodyHash      string `xml:"MD5OfBody"`
	ReceiptHandle string `xml:"ReceiptHandle"`
}
