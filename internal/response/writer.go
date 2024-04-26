package response

import (
	"encoding/json"
	"net/http"
)

func WriteError(w http.ResponseWriter, err error, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(ErrorCode(err))
	errorMessage := message
	if message == "" {
		errorMessage = err.Error()
	}
	var payload interface{}
	if data != nil {
		payload = map[string]interface{}{
			"message": errorMessage,
			"errors":  data,
		}
	} else {
		payload = map[string]interface{}{
			"message": errorMessage,
		}
	}
	json.NewEncoder(w).Encode(payload)
}

func WriteSuccess(w http.ResponseWriter, data interface{}, code int) {
	w.Header().Set("Content-Type", "application/json")
	if code != 0 {
		w.WriteHeader(code)
	} else {
		w.WriteHeader(http.StatusOK)
	}
	var payload interface{}
	if data != nil {
		payload = data
	} else {
		payload = ResponseMessage{Message: "Success"}
	}
	json.NewEncoder(w).Encode(payload)
}
