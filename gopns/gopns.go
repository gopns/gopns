package gopns

import (
	config "github.com/gopns/gopnsconfig"
	"github.com/gopns/notification"
	"github.com/stefantalpalaru/pool"
	"log"
	"runtime"
	"time"
)

// gopns package level global state
var NotificationSender *notification.NotificationSender

func Start() {

	var WorkerPool *pool.Pool = startWorkerPool()
	//setup notification sender
	NotificationSender = &notification.NotificationSender{AwsConfig: config.AWSConfigInstance(), WorkerPool: WorkerPool}
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
