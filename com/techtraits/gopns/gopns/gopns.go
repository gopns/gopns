package gopns

import (
	config "github.com/gopns/gopns/com/techtraits/gopns/gopnsconfig"
	"github.com/gopns/gopns/com/techtraits/gopns/notification"
	"github.com/stefantalpalaru/pool"
	"runtime"
)

// gopns package level global state
var WorkerPool *pool.Pool = nil
var NotificationSender *notification.NotificationSender = nil

func Start() {

	//setup a generic worker pool
	cpus := runtime.NumCPU()
	runtime.GOMAXPROCS(cpus)
	WorkerPool := pool.New(cpus)
	WorkerPool.Run()

	//setup notification sender
	NotificationSender := &notification.NotificationSender{awsConfig: &config.AWSConfigInstance(), workerPool: WorkerPool}
}

func Stop() {
	WorkerPool.Stop()
}
