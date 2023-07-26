package server

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/donwitobanana/fh-sharding/internal/sharding"
)

type Request struct {
	Messages []sharding.Message
}

func handleMessages(svc sharding.Service) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		var req Request
		dec := json.NewDecoder(r.Body)
		if err := dec.Decode(&req); err != nil {
			writeErrorResponse(w, http.StatusBadRequest, err)
			return
		}

		err := svc.ShardData(req.Messages)
		if err != nil {
			writeErrorResponse(w, http.StatusInternalServerError, errors.New("failed to shard data"))
			return
		}

		w.WriteHeader(http.StatusOK)
	}

}

func writeErrorResponse(w http.ResponseWriter, code int, err error) {
	w.WriteHeader(code)
	w.Write([]byte(err.Error()))
}
