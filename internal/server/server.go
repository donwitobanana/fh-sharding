package server

import (
	"net/http"

	"github.com/donwitobanana/fh-sharding/internal/sharding"
)

func RegisterHandlers(svc sharding.Service) {
	http.HandleFunc("/", handleMessages(svc))
}
