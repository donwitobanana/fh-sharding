package workers

import (
	"context"
	"fmt"
	"time"

	"github.com/donwitobanana/fh-sharding/internal/storage"
)

type WorkerpoolConfig struct {
	WorkersMaxNumber  int
	WorkersInitNumber int
	WorkersTimeout    time.Duration
	Partition         int
	WorkersBatchSize  int
}

type Workerpool struct {
	workerpoolConfig WorkerpoolConfig
	storageSvc 	 storage.Service
	messageChan      chan interface{}
	workerSem        chan bool
	errChan          chan error
}

func NewWorkerpool(config WorkerpoolConfig, storageSvc storage.Service) *Workerpool {
	return &Workerpool{
		workerpoolConfig: config,
		storageSvc: storageSvc,
		messageChan:      make(chan interface{}),
		workerSem:        make(chan bool, config.WorkersMaxNumber),
		errChan:          make(chan error),
	}
}

func (w *Workerpool) Start(ctx context.Context) error {
	fmt.Printf("starting workerpool, partition %d\n", w.workerpoolConfig.Partition)

	workerConfig := WorkerConfig{
		input:     w.messageChan,
		errChan:   w.errChan,
		sem:       w.workerSem,
		partition: w.workerpoolConfig.Partition,
		batchSize: w.workerpoolConfig.WorkersBatchSize,
		timeout:   w.workerpoolConfig.WorkersTimeout,
	}

	ctx, cancel := context.WithCancel(ctx)

	go w.initWorkerpool(ctx, workerConfig)

	go w.spawnWorkers(ctx, workerConfig)

	for {
		select {
		case <-ctx.Done():
			cancel()
			return nil
		case err := <-w.errChan:
			if err != nil {
				cancel()
				return err
			}
		}
	}

}

func (w *Workerpool) initWorkerpool(ctx context.Context, workerConfig WorkerConfig) {
	for i := 0; i < w.workerpoolConfig.WorkersInitNumber; i++ {
		select {
		case <-ctx.Done():
			return
		case w.workerSem <- true:
			worker := NewWorker(workerConfig, w.storageSvc)
			go worker.Start(ctx)
		}
	}
}

func (w *Workerpool) spawnWorkers(ctx context.Context, workerConfig WorkerConfig) {
	<-time.After(w.workerpoolConfig.WorkersTimeout - time.Second*10)
	for {
		select {
		case <-ctx.Done():
			return
		case w.workerSem <- true:
			worker := NewWorker(workerConfig, w.storageSvc)
			go worker.Start(ctx)
		}
	}
}

func (w *Workerpool) SendMessage(msg interface{}) {
	w.messageChan <- msg
}
