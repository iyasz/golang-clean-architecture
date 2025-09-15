package helper

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
)

func SuccessResponse(w http.ResponseWriter, data interface{}, statusCode int) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		logrus.Warnf("Failed to encode response: %+v", err)
		ErrorResponse(w, err)
	}

}

func ErrorResponse(w http.ResponseWriter, err error) {
	code := http.StatusInternalServerError
	
	if e, ok := err.(*Error); ok {
		code = e.Code
	}
	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	
	response := map[string]string{
		"error": err.Error(),
	}
	
	if err := json.NewEncoder(w).Encode(response); err != nil {
		logrus.Warnf("Failed to encode response: %+v", err)
	}
}