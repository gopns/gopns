package gopnsapp

import (
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sqs"
	config "github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/notification"
	"github.com/stefantalpalaru/pool"

	"log"
	"runtime"
	"time"
)

type GopnsAppStruct struct {
	DynamoClient dynamodb.DynamoClient
	SQSClient    sqs.SQSClient
}
type GopnsApp interface {
	Start()
}

// gopns package level global state
var NotificationSender *notification.NotificationSender

func Initilize() (GopnsApp, error) {

	gopnasapp_ := &GopnsAppStruct{}

	err := gopnasapp_.initilizeDB()
	if err != nil {
		return nil, err
	}

	err = gopnasapp_.initilizeSQS()
	if err != nil {
		return nil, err
	}

	return gopnasapp_, nil
}

func (this *GopnsAppStruct) Start() {

	var WorkerPool *pool.Pool = this.startWorkerPool()
	//setup notification sender
	NotificationSender = &notification.NotificationSender{AwsConfig: config.AWSConfigInstance(), WorkerPool: WorkerPool}
}

func (this *GopnsAppStruct) initilizeDB() error {
	var err error
	this.DynamoClient, err = dynamodb.Initilize(
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

func (this *GopnsAppStruct) initilizeSQS() error {
	var err error
	this.SQSClient = sqs.Initilize(
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

func (this *GopnsAppStruct) startWorkerPool() *pool.Pool {
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
