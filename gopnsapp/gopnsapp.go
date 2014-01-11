package gopnsapp

import (
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sqs"
	"github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/notification"
	"github.com/stefantalpalaru/pool"
	"log"
	"runtime"
	"time"
)

// gopns package level global state
var NotificationSender *notification.NotificationSender

func Start() {

	initAWS()

	var WorkerPool *pool.Pool = startWorkerPool()
	//setup notification sender
	NotificationSender = &notification.NotificationSender{
		AwsConfig:  gopnsconfig.AWSConfigInstance(),
		WorkerPool: WorkerPool}

	//create a notification consumer
	var NotificationConsumer notification.NotificationConsumer = notification.NewSQSNotifictionConsumer(
		gopnsconfig.AWSConfigInstance().SqsQueueUrl(),
		gopnsconfig.AWSConfigInstance(),
		NotificationSender)

	NotificationConsumer.Start()
}

// for the time being -- will change after aws client changes
func initAWS() {

	// int aws

	err := dynamodb.Initilize(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().DynamoTable(),
		gopnsconfig.AWSConfigInstance().InitialReadCapacity(),
		gopnsconfig.AWSConfigInstance().InitialWriteCapacity())

	if err != nil {
		log.Fatalf("Unable to initilize Dynamo DB %s", err.Error())
	}

	// int sqs
	err, sqsQueue := sqs.Initilize(gopnsconfig.AWSConfigInstance().UserID(),
		gopnsconfig.AWSConfigInstance().UserSecret(),
		gopnsconfig.AWSConfigInstance().Region(),
		gopnsconfig.AWSConfigInstance().SqsQueueName())

	if err != nil {
		log.Fatalf("Unable to initilize SQS %s", err.Error())
	} else {
		log.Printf("Using SQS Queue %s", sqsQueue.QueueUrl)
		gopnsconfig.AWSConfigInstance().SetSqsQueueUrl(sqsQueue.QueueUrl)
	}

}

func startWorkerPool() *pool.Pool {
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
