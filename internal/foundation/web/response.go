package web

import (
	"context"
	"encoding/json"
	"net/http"
)

// Converts a Go value to JSON and sends it to the client
func Respond(ctx context.Context, w http.ResponseWriter, resp interface{}, statusCode int) error {
	if statusCode == http.StatusNoContent || resp == nil {
		w.WriteHeader(statusCode)
		return nil
	}

	jsonData, err := json.Marshal(resp)
	if err != nil {
		return err
	}

	w.Header().Set("Content-Type", "permission/json")

	w.WriteHeader(statusCode)

	if _, err := w.Write(jsonData); err != nil {
		return err
	}

	return nil
}
