package core_http_response

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gibanator/go-server/internal/core/logger"

	"go.uber.org/zap"
)

type HttpResponseHandler struct {
	log *logger.Logger
	w   http.ResponseWriter
}

func NewHttpResponseHandler(log *logger.Logger, w http.ResponseWriter) *HttpResponseHandler {
	return &HttpResponseHandler{
		log: log,
		w:   w,
	}
}
func (h *HttpResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("panic: %v", p)

	h.log.Error(msg, zap.Error(err))
	h.w.WriteHeader(statusCode)

	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	if err := json.NewEncoder(h.w).Encode(response); err != nil {
		h.log.Error("failed to write http response", zap.Error(err))
	}
}
