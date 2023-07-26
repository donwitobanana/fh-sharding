package sharding

import "github.com/donwitobanana/fh-sharding/internal/workers"

type Service interface {
	ShardData([]Message) error
}

type service struct {
	workerpoolMap map[int]*workers.Workerpool
}

func NewService(workerpoolMap map[int]*workers.Workerpool) Service {
	return &service{
		workerpoolMap: workerpoolMap,
	}
}

func (s *service) ShardData(messages []Message) error {
	for _, msg := range messages {
		partition := getPartition(msg, len(s.workerpoolMap))
		workerpool := s.workerpoolMap[partition]
		workerpool.SendMessage(msg)
	}

	return nil
}

func getPartition(msg Message, numPartitions int) int {
	return int(msg.Timestamp) % numPartitions 
}