package adapters

import (
	"encoding/json"
	"net/http"
	"strconv"
)

func marshalAndWriteErrorResponse(w http.ResponseWriter, errorMessage string, statusCode int) error {
	msg := map[string]string{
		"message": errorMessage,
	}
	responseBody, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	if err := writeResponse(w, responseBody, statusCode); err != nil {
		return err
	}

	return nil
}

func writeResponse(w http.ResponseWriter, body []byte, statusCode int) error {
	w.Header().Set(`Content-Type`, `application/json`)
	w.Header().Set(`Content-Length`, strconv.Itoa(len(body)))

	w.WriteHeader(statusCode)
	if _, err := w.Write(body); err != nil {
		return err
	}

	return nil
}
