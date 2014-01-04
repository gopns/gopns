package sqs

type SqsQueue struct {
	QueueUrl string `xml:"CreateQueueResult>QueueUrl"`
}
