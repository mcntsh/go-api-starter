package api

import (
	"encoding/json"
	"github.com/Sirupsen/logrus"
	"net/http"
	"strconv"
)

type Response struct {
	rw http.ResponseWriter `json:"-"`

	Status *ResponseStatus `json:"status"`
	Body   interface{}     `json:"body"`
}

type ResponseStatus struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

func (r *Response) Write() {
	json, err := json.Marshal(r)
	if err != nil {
		logrus.Fatal("Could not marshal JSON in the API writer")
	}

	r.rw.Header().Set("Content-Type", "application/json; charset=UTF-8")
	r.rw.Header().Set("Content-Length", strconv.Itoa(len(json)))
	r.rw.WriteHeader(r.Status.Code)

	r.rw.Write(json)
}

func statusCodeFatal(code int) bool {
	fatalCodes := []int{500, 501, 502, 503, 504, 505, 511}
	for _, c := range fatalCodes {
		if c == code {
			return true
		}
	}

	return false
}

func WriteErrorResponse(w http.ResponseWriter, r *http.Request, code int, err error) {
	res := &Response{rw: w}

	res.Status = &ResponseStatus{
		Code:    code,
		Message: http.StatusText(code),
		Error:   err.Error(),
	}

	if statusCodeFatal(code) == true {
		logrus.WithFields(logrus.Fields{
			"method": r.Method,
			"url":    r.URL,
			"code":   code,
			"error":  err,
		}).Warn("API Handler returned an error!")
	}

	res.Write()
}

func WriteResponse(w http.ResponseWriter, data interface{}) {
	res := &Response{rw: w}

	res.Status = &ResponseStatus{
		Code:    http.StatusOK,
		Message: http.StatusText(http.StatusOK),
	}

	res.Body = data

	res.Write()
}
