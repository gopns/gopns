package gopnsapp

import (
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/aws/sqs"
	config "github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/notification"
	"github.com/stefantalpalaru/pool"
	"log"
	"runtime"
	"time"
)

type GopnsApp interface {
	Start()
}

type GopnsApplication struct {
	DynamoClient dynamodb.DynamoClient
	SQSClient    sqs.SQSClient
	SNSClient    sns.SNSClient
}

// gopns package level global state
var NotificationSender *notification.NotificationSender

func New() (GopnsApp, error) {

	gopnasapp_ := &GopnsApplication{}

	err := gopnasapp_.setupDynamoDB()
	if err != nil {
		return nil, err
	}

	err = gopnasapp_.setupSQS()
	if err != nil {
		return nil, err
	}

	err = gopnasapp_.setupSNS()
	if err != nil {
		return nil, err
	}

	return gopnasapp_, nil
}

// ToDo return appropriate errors
func (this *GopnsApplication) Start() {

	var WorkerPool *pool.Pool = this.startWorkerPool()
	//setup notification sender
	NotificationSender = &notification.NotificationSender{
		SnsClient:    this.SNSClient,
		WorkerPool:   WorkerPool,
		PlatformApps: config.AWSConfigInstance().PlatformApps()}

	//create a notification consumer
	var NotificationConsumer notification.NotificationConsumer = notification.NewSQSNotifictionConsumer(
		config.AWSConfigInstance().SqsQueueUrl(),
		this.SQSClient,
		NotificationSender)

	NotificationConsumer.Start()
}

func (this *GopnsApplication) setupDynamoDB() error {
	var err error
	this.DynamoClient, err = dynamodb.New(
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region())

	if err != nil {
		return err
	}

	if found, err := this.DynamoClient.FindTable(config.AWSConfigInstance().DynamoTable()); err != nil {
		return err
	} else if found {
		return nil
	} else {

		createTableRequest := dynamodb.CreateTableRequest{
			[]dynamodb.AttributeDefinition{dynamodb.AttributeDefinition{"alias", "S"}},
			config.AWSConfigInstance().DynamoTable(),
			[]dynamodb.KeySchema{dynamodb.KeySchema{"alias", "HASH"}},
			dynamodb.ProvisionedThroughput{
				config.AWSConfigInstance().InitialReadCapacity(),
				config.AWSConfigInstance().InitialWriteCapacity()}}
		return this.DynamoClient.CreateTable(createTableRequest)
	}

	return err
}

func (this *GopnsApplication) setupSQS() error {
	var err error
	this.SQSClient = sqs.New(
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region())

	if sqsQueue, err := this.SQSClient.CreateQueue(config.AWSConfigInstance().SqsQueueName()); err != nil {
		log.Fatalf("Unable to initilize SQS %s", err.Error())
	} else {
		log.Printf("Using SQS Queue %s", sqsQueue.QueueUrl)
		config.AWSConfigInstance().SetSqsQueueUrl(sqsQueue.QueueUrl)
	}

	return err
}

func (this *GopnsApplication) setupSNS() error {
	var err error
	this.SNSClient, err = sns.New(
		config.AWSConfigInstance().UserID(),
		config.AWSConfigInstance().UserSecret(),
		config.AWSConfigInstance().Region())

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (this *GopnsApplication) startWorkerPool() *pool.Pool {
	//setup a generic worker pool
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	WorkerPool := pool.New(cpus)
	WorkerPool.Run()
	log.Println("Worker pool started with", cpus, "workers")
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for _ = range ticker.C {
			status := WorkerPool.Status()
			log.Println(status.Submitted, "submitted jobs,", status.Running, "running,", status.Completed, "completed.")
		}
	}()

	return WorkerPool
}
