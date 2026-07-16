package transport

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/frhan23/dashboard-go/internal/entity"
)

type ErrorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Details any    `json:"details,omitempty"`
}

func WriteJSON(w http.ResponseWriter, code int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(payload)
}

func CodeToStatus(code entity.Code) int {
	switch code {
	case entity.ErrorCodeBadRequest:
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}

func WriteAppError(w http.ResponseWriter, appErr *entity.AppError) {
	status := CodeToStatus(appErr.Code)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	resp := ErrorResponse{
		Code:    string(appErr.Code),
		Message: appErr.Message,
		Details: appErr.Details,
	}
	err := json.NewEncoder(w).Encode(resp)
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}

func WriteError(w http.ResponseWriter, err error) {
	if err == nil {
		w.WriteHeader(http.StatusOK)
		return
	}
	var aErr *entity.AppError
	if errors.As(err, &aErr) {
		WriteAppError(w, aErr)
		return
	}
	// fallback
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	err = json.NewEncoder(w).Encode(ErrorResponse{
		Code:    string(entity.ErrorCodeInternal),
		Message: "internal error",
	})
	if err != nil {
		http.Error(w, "internal server error", http.StatusInternalServerError)
		return
	}
}
