package workers

import (
	"context"
	"fmt"
	"time"
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
	taskQueue        chan interface{}
	workerSem        chan bool
	errChan          chan error
}

func NewWorkerpool(config WorkerpoolConfig) *Workerpool {
	return &Workerpool{
		workerpoolConfig: config,
		taskQueue:        make(chan interface{}),
		workerSem:        make(chan bool, config.WorkersMaxNumber),
		errChan:          make(chan error),
	}
}

func (w *Workerpool) Start(ctx context.Context) error {
	fmt.Printf("starting workerpool, partition %d\n", w.workerpoolConfig.Partition)

	workerConfig := WorkerConfig{
		input:     w.taskQueue,
		errChan:   w.errChan,
		sem:       w.workerSem,
		partition: w.workerpoolConfig.Partition,
		batchSize: w.workerpoolConfig.WorkersBatchSize,
	}

	timeoutCtx, timeoutCtxCancel := context.WithTimeout(ctx, w.workerpoolConfig.WorkersTimeout)

	go w.initWorkerpool(timeoutCtx, workerConfig)

	go w.AddWorkers(ctx, workerConfig)

	for {
		select {
		case <-ctx.Done():
			timeoutCtxCancel()
			return nil
		case err := <-w.errChan:
			if err != nil {
				timeoutCtxCancel()
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
				worker := NewWorker(workerConfig)
				go worker.Start(ctx)
			}
		}
}

func (w *Workerpool) AddWorkers(ctx context.Context, workerConfig WorkerConfig) {
	<-time.After(w.workerpoolConfig.WorkersTimeout - time.Second*10)
	for {
		select {
		case <-ctx.Done():
			return
		case w.workerSem <- true:

			workerCtx, cancel := context.WithTimeout(context.Background(), w.workerpoolConfig.WorkersTimeout)
			defer cancel()

			worker := NewWorker(workerConfig)
			go worker.Start(workerCtx)
		}
	}
}

func (w *Workerpool) SendMessage(msg interface{}) {
	w.taskQueue <- msg
}

