package gopnsapp

import (
	"code.google.com/p/gorest"
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/aws/sqs"
	"github.com/gopns/gopns/device"
	config "github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/notification"
	"github.com/gopns/gopns/rest"
	"github.com/stefantalpalaru/pool"
	"log"
	"net/http"
	"runtime"
	"time"
)

type GopnsApp interface {
	Start() error
}

type GopnsApplication struct {
	DynamoClient         dynamodb.DynamoClient
	SQSClient            sqs.SQSClient
	SNSClient            sns.SNSClient
	WorkerPool           pool.Pool
	NotificationSender   notification.NotificationSender
	NotificationConsumer notification.NotificationConsumer
	AppMode              config.APPLICATION_MODE
	DeviceManager        device.DeviceManager
}

func New() (GopnsApp, error) {

	gopnasapp_ := &GopnsApplication{}

	gopnasapp_.AppMode = config.ParseConfig()

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

	err = gopnasapp_.createWorkerPool()
	if err != nil {
		return nil, err
	}

	//setup notification sender
	gopnasapp_.NotificationSender = notification.NotificationSender{
		SnsClient:    gopnasapp_.SNSClient,
		WorkerPool:   &gopnasapp_.WorkerPool,
		PlatformApps: config.AWSConfigInstance().PlatformApps()}

	//create a notification consumer
	gopnasapp_.NotificationConsumer = notification.NewSQSNotifictionConsumer(
		config.AWSConfigInstance().SqsQueueUrl(),
		gopnasapp_.SQSClient,
		&gopnasapp_.NotificationSender)

	//create a device manager
	gopnasapp_.DeviceManager = device.New(
		gopnasapp_.SNSClient,
		gopnasapp_.DynamoClient)

	return gopnasapp_, nil
}

// ToDo return appropriate errors
func (this *GopnsApplication) Start() error {

	err := this.runWorkerPool()
	if err != nil {
		return err
	}

	err = this.NotificationConsumer.Start()
	if err != nil {
		return err
	}

	if this.AppMode == config.SERVER_MODE {
		this.setupRestServices()
	} else if this.AppMode == config.REGISTER_MODE {

	} else if this.AppMode == config.SEND_MODE {

	}

	return nil
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

func (this *GopnsApplication) createWorkerPool() error {
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	this.WorkerPool = *pool.New(cpus)
	log.Println("Worker pool created with", cpus, "workers")
	return nil
}

func (this *GopnsApplication) runWorkerPool() error {
	this.WorkerPool.Run()
	log.Println("Worker pool started")
	ticker := time.NewTicker(time.Second * 60)
	go func() {
		for _ = range ticker.C {
			status := this.WorkerPool.Status()
			log.Println(status.Submitted, "submitted jobs,", status.Running, "running,", status.Completed, "completed.")
		}
	}()

	return nil
}

func (this *GopnsApplication) setupRestServices() {

	notificationService := new(rest.NotificationService)
	notificationService.NotificationSender = &this.NotificationSender
	notificationService.DeviceManager = this.DeviceManager

	deviceService := new(rest.DeviceService)
	deviceService.DeviceManager = this.DeviceManager

	gorest.RegisterService(deviceService)
	gorest.RegisterService(notificationService)
	http.Handle("/", gorest.Handle())
	http.ListenAndServe(":"+config.BaseConfigInstance().Port(), nil)
}