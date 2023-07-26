package workers

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/google/uuid"
)

const (
	storagePrefix = "storage"
)

type Worker struct {
	id           uuid.UUID
	cfg          WorkerConfig
	batchCounter int
	itemCounter  int
	buffer       []interface{}
}

type WorkerConfig struct {
	input     chan interface{}
	errChan   chan error
	sem       chan bool
	partition int
	batchSize int
}

func NewWorker(cfg WorkerConfig) *Worker {
	id := uuid.New()

	return &Worker{
		id:     id,
		cfg:    cfg,
		buffer: make([]interface{}, 0, cfg.batchSize),
	}
}

func (w *Worker) Start(ctx context.Context) {
	fmt.Printf("starting worker, partition %d, id %s\n", w.cfg.partition, w.id.String())
	for {
		select {
		case <-ctx.Done():
			err := w.Stop()
			if err != nil {
				w.cfg.errChan <- err
				return
			}
			return
		case msg := <-w.cfg.input:
			err := w.Process(msg)
			if err != nil {
				w.cfg.errChan <- err
				return
			}
		}
	}
}

func (w *Worker) Stop() error {
	fmt.Printf("stopping worker, partition %d, id %s\n", w.cfg.partition, w.id.String())
	<-w.cfg.sem
	if len(w.buffer) == 0 {
		return nil
	}

	err := w.writeData()
	if err != nil {
		return err
	}

	return nil
}

func (w *Worker) Process(msg interface{}) error {
	w.buffer = append(w.buffer, msg)
	w.itemCounter++
	if w.itemCounter == w.cfg.batchSize {
		err := w.writeData()
		if err != nil {
			return err
		}
	}

	return nil
}

func (w *Worker) writeData() error {
	fmt.Printf("writing data, partition %d, id %s, batch no %d\n", w.cfg.partition, w.id.String(), w.batchCounter)

	file, err := os.Create(w.getDestination())
	if err != nil {
		return err
	}
	defer file.Close()

	data, err := json.Marshal(w.buffer)
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	if err != nil {
		return err
	}

	w.buffer = make([]interface{}, 0, 5)
	w.itemCounter = 0
	w.batchCounter++
	return nil
}

func (w *Worker) getDestination() string {
	return fmt.Sprintf("%s/%d/%s/%d", storagePrefix, w.cfg.partition, w.id.String(), w.batchCounter)
}
