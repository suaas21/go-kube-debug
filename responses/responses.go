package responses

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status  int         `json:"-"`
	Data    interface{} `json:"data"`
	Message string      `json:"message"`
}

// ServeJSON serves json to http client
func ServeJSON(w http.ResponseWriter, status int, msg string, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	resp := &Response{
		Status:  status,
		Data:    data,
		Message: msg,
	}

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		return err
	}

	return nil
}
