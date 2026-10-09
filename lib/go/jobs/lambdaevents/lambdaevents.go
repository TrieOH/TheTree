// Package lambdaevents serves the worker function's non-HTTP invocations.
// The Lambda Web Adapter forwards every non-HTTP event as a POST to its
// pass-through path (AWS_LWA_PASS_THROUGH_PATH, default /events) and returns
// the response body as the invocation result. Two event shapes arrive:
//
//   - an SQS batch ({"Records": [...]}) from the queue's event source
//     mapping: each record's body is a jobs.Envelope. Failed records are
//     reported back as batchItemFailures (the mapping must enable
//     ReportBatchItemFailures), so only they are redelivered and,
//     past the queue's maxReceiveCount, dead-lettered.
//   - a bare jobs.Envelope from a scheduler target's Input: a non-2xx
//     response fails the invocation, so the scheduler's retry policy
//     applies.
//
// Mount it only on the worker role: on the api role the gateway would
// expose it.
package lambdaevents

import (
	"encoding/json"
	"io"
	"net/http"

	"lib/jobs"
	"lib/telemetry"

	"go.uber.org/zap"
)

// Path is the Lambda Web Adapter's default pass-through path.
const Path = "/events"

// maxEventBytes bounds the request body: an SQS batch is at most 10
// messages of 256KiB, plus the record metadata.
const maxEventBytes = 6 << 20

type sqsRecord struct {
	MessageID   string `json:"messageId"`
	EventSource string `json:"eventSource"`
	Body        string `json:"body"`
}

type event struct {
	jobs.Envelope

	Records []sqsRecord `json:"Records"`
}

type batchItemFailure struct {
	ItemIdentifier string `json:"itemIdentifier"`
}

type batchResponse struct {
	BatchItemFailures []batchItemFailure `json:"batchItemFailures"`
}

// Handler dispatches events into registry.
func Handler(registry *jobs.Registry) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		raw, err := io.ReadAll(io.LimitReader(r.Body, maxEventBytes))
		if err != nil {
			http.Error(w, "read event", http.StatusBadRequest)
			return
		}
		var ev event
		err = json.Unmarshal(raw, &ev)
		if err != nil {
			http.Error(w, "decode event", http.StatusBadRequest)
			return
		}

		if ev.Records != nil {
			resp := batchResponse{BatchItemFailures: []batchItemFailure{}}
			for _, rec := range ev.Records {
				err := dispatchRecord(r, registry, rec)
				if err != nil {
					telemetry.Log().Error("queued job failed",
						zap.String("message_id", rec.MessageID), zap.Error(err))
					resp.BatchItemFailures = append(resp.BatchItemFailures, batchItemFailure{ItemIdentifier: rec.MessageID})
				}
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(resp)
			return
		}

		if ev.Kind == "" {
			http.Error(w, "unrecognized event", http.StatusBadRequest)
			return
		}
		err = registry.Dispatch(r.Context(), ev.Envelope)
		if err != nil {
			telemetry.Log().Error("scheduled job failed", zap.String("kind", ev.Kind), zap.Error(err))
			http.Error(w, "job failed", http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusNoContent)
	})
}

func dispatchRecord(r *http.Request, registry *jobs.Registry, rec sqsRecord) error {
	var env jobs.Envelope
	err := json.Unmarshal([]byte(rec.Body), &env)
	if err != nil {
		return err
	}
	return registry.Dispatch(r.Context(), env)
}
