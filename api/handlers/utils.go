package handlers

import (
	"encoding/json"
	"net/http"

	"go.uber.org/zap"

	"github.com/BugBridge/bugbridge-api/models"
)

// ErrorStatus is a useful function that will log, write http headers and body for a
// given message, status code and error
func ErrorStatus(
	message string,
	httpStatusCode int,
	w http.ResponseWriter,
	err error,
) {
	zap.S().With(err).Error(message)
	w.WriteHeader(httpStatusCode)
	b, _ := json.Marshal(models.ErrorMessageResponse{Response: models.MessageError{Message: message, Error: err.Error()}})
	w.Write(b)
}
