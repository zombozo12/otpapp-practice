package http

import (
	"encoding/json"
	"log"
	"net/http"
)

func badRequest(message string, w http.ResponseWriter) {
	if message == "" {
		message = "Bad Request"
	}

	data := map[string]interface{}{
		"error": message,
	}

	response, err := json.Marshal(data)
	if err != nil {
		log.Fatalf("Bad Request Error : %+v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusBadRequest)
	w.Write(response)
}

func okResponse(message string, data interface{}, w http.ResponseWriter) {
	if data == nil {
		log.Fatalf("OK Response cannot be null")
		return
	}

	res := map[string]interface{}{
		"message": message,
		"data":    data,
	}

	resJson, err := json.Marshal(res)
	if err != nil {
		log.Fatalf("Set OK Error: %+v", err)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(resJson)
}
