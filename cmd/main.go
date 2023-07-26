package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/donwitobanana/fh-sharding/cmd/config"
	"github.com/donwitobanana/fh-sharding/internal/server"
	"github.com/donwitobanana/fh-sharding/internal/sharding"
	"github.com/donwitobanana/fh-sharding/internal/workers"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	appConfig, err := config.LoadAppConfig()
	if err != nil {
		fmt.Printf("failed to load app config: %s\n", err.Error())
		return
	}

	workerpoolConfig := workers.WorkerpoolConfig{
		WorkersMaxNumber:  appConfig.WorkersMaxNumber,
		WorkersInitNumber: appConfig.WorkersInitNumber,
		WorkersTimeout:    time.Duration(appConfig.WorkersTimeoutInSeconds) * time.Second,
		WorkersBatchSize:  appConfig.WorkersBatchSize,
	}

	errChan := make(chan error)
	workerpoolMap := initWorkerpools(ctx, workerpoolConfig, appConfig.PartitionsNumber, errChan)
	svc := sharding.NewService(workerpoolMap)

	server.RegisterHandlers(svc)
	startServer(errChan)

	err = <-errChan
	if err != nil {
		fmt.Printf("error: %s\n", err.Error())
	}
}

func initWorkerpools(ctx context.Context, config workers.WorkerpoolConfig, partitionsNumber int, errChan chan error) map[int]*workers.Workerpool {
	fmt.Println("initializing workerpools")

	workerpoolMap := make(map[int]*workers.Workerpool, partitionsNumber)
	for i := 0; i < partitionsNumber; i++ {
		config.Partition = i
		workerpool := workers.NewWorkerpool(config)
		go func() {
			err := workerpool.Start(ctx)
			if err != nil {
				errChan <- err
				return
			}
		}()
		workerpoolMap[i] = workerpool
	}

	return workerpoolMap
}

func startServer(errChan chan error) {
	go func() {
		fmt.Println("starting server")
		err := http.ListenAndServe(":8080", http.DefaultServeMux)
		fmt.Printf("server stopped: %s\n", err.Error())
		if err != nil {
			errChan <- err
			return
		}
	}()
}
