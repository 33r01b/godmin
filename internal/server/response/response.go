package response

import (
	"encoding/json"
	"net/http"
)

type Handler interface {
	Respond(w http.ResponseWriter, r *http.Request, code int, data interface{})
	Error(w http.ResponseWriter, r *http.Request, code int, err error)
}

type Response struct {
}

func (res *Response) Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	if data != nil {
		if err := json.NewEncoder(w).Encode(data); err != nil {
			res.Error(w, r, http.StatusInternalServerError, err)
		}
	}
}

func (res *Response) Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	res.Respond(w, r, code, map[string]string{"error": err.Error()})
}

func NewResponse() Handler {
	return &Response{}
}

type Writer struct {
	http.ResponseWriter
	Code int
}

func (w *Writer) WriteHeader(statusCode int) {
	w.Code = statusCode
	w.ResponseWriter.WriteHeader(statusCode)
}
