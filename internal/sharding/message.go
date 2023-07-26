package sharding

type Message struct {
	Timestamp int64 `json:"timestamp"`
	Value	 string `json:"value"`
}