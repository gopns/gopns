package gopnsapp

import (
	"github.com/emicklei/go-restful"
	"github.com/gopns/gopns/access"
	"github.com/gopns/gopns/aws/dynamodb"
	"github.com/gopns/gopns/aws/sns"
	"github.com/gopns/gopns/aws/sqs"
	config "github.com/gopns/gopns/gopnsconfig"
	"github.com/gopns/gopns/metrics"
	"github.com/gopns/gopns/notification"
	"github.com/gopns/gopns/pool"
	"github.com/gopns/gopns/rest"
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
	BaseConfig           config.BaseConfig
	AWSConfig            config.AWSConfig
	DeviceManager        access.DeviceManager
	WsContainer          restful.Container
}

func New() (GopnsApp, error) {

	gopnasapp_ := &GopnsApplication{}

	gopnasapp_.AppMode, gopnasapp_.BaseConfig, gopnasapp_.AWSConfig = config.ParseConfig()

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
		PlatformApps: gopnasapp_.AWSConfig.PlatformAppsMap()}

	//create a notification consumer
	gopnasapp_.NotificationConsumer = notification.NewSQSNotifictionConsumer(
		gopnasapp_.AWSConfig.SqsQueueUrl(),
		gopnasapp_.SQSClient,
		&gopnasapp_.NotificationSender)

	//create a device manager
	gopnasapp_.DeviceManager = access.NewDeviceManager(
		gopnasapp_.SNSClient,
		gopnasapp_.DynamoClient,
		gopnasapp_.AWSConfig.DynamoTable())

	//

	gopnasapp_.setupMetrics()

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
		this.AWSConfig.UserID(),
		this.AWSConfig.UserSecret(),
		this.AWSConfig.Region())

	if err != nil {
		return err
	}

	if found, err := this.DynamoClient.FindTable(this.AWSConfig.DynamoTable()); err != nil {
		return err
	} else if found {
		return nil
	} else {

		createTableRequest := dynamodb.CreateTableRequest{
			[]dynamodb.AttributeDefinition{dynamodb.AttributeDefinition{"alias", "S"}},
			this.AWSConfig.DynamoTable(),
			[]dynamodb.KeySchema{dynamodb.KeySchema{"alias", "HASH"}},
			dynamodb.ProvisionedThroughput{
				this.AWSConfig.InitialReadCapacity(),
				this.AWSConfig.InitialWriteCapacity()}}
		return this.DynamoClient.CreateTable(createTableRequest)
	}

	return err
}

func (this *GopnsApplication) setupSQS() error {
	var err error
	this.SQSClient = sqs.New(
		this.AWSConfig.UserID(),
		this.AWSConfig.UserSecret(),
		this.AWSConfig.Region())

	if sqsQueue, err := this.SQSClient.CreateQueue(this.AWSConfig.SqsQueueName()); err != nil {
		log.Fatalf("Unable to initilize SQS %s", err.Error())
	} else {
		log.Printf("Using SQS Queue %s", sqsQueue.QueueUrl)
		this.AWSConfig.SetSqsQueueUrl(sqsQueue.QueueUrl)
	}

	return err
}

func (this *GopnsApplication) setupSNS() error {
	var err error
	this.SNSClient, err = sns.New(
		this.AWSConfig.UserID(),
		this.AWSConfig.UserSecret(),
		this.AWSConfig.Region())

	if err != nil {
		return err
	} else {
		return nil
	}
}

func (this *GopnsApplication) setupMetrics() error {

	//metrics.StartMetricPrinter()

	metrics.StartGraphiteReporter(
		this.BaseConfig.MetricsServer(),
		this.BaseConfig.MetricsAPIKey(),
		this.BaseConfig.MetricsPrefix())

	return nil
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

	//setup a new services container for gopns rest services
	this.WsContainer = *restful.NewContainer()
	this.WsContainer.Filter(rest.ExceptionFilter)

	//ToDo read the gopns rest root path from config (re: embeddable app)
	rootPath := "/gopns/v1/" //without the last slash (e.g., /rest/gopns)

	restful.DefaultResponseMimeType = restful.MIME_JSON

	notificationService := new(rest.NotificationService)
	notificationService.NotificationSender = &this.NotificationSender
	notificationService.DeviceManager = this.DeviceManager
	// register service with our services container
	notificationService.Register(&this.WsContainer, rootPath)

	deviceService := new(rest.DeviceService)
	deviceService.DeviceManager = this.DeviceManager
	// register service with our services container
	deviceService.Register(&this.WsContainer, rootPath)

	log.Printf("start listening on localhost:8080")
	server := &http.Server{Addr: ":" + this.BaseConfig.Port(), Handler: &this.WsContainer}
	log.Fatal(server.ListenAndServe())
}
