package httpprotocol

import (
	"encoding/json"
	"errors"
	"net/http"
)

type Result struct {
	Code   int
	Result interface{}
}

// here use costom define code
func Failed(err error) *Result {
	return &Result{
		Code:   http.StatusInternalServerError,
		Result: err,
	}
}

func Succeed() *Result {
	return &Result{
		Code:   http.StatusOK,
		Result: nil,
	}
}

func SucceedWithResult(result interface{}) *Result {
	return &Result{
		Code:   http.StatusOK,
		Result: result,
	}
}

func EncodeResult(w http.ResponseWriter, result *Result) {
	w.WriteHeader(result.Code)
	if result.Result == nil {
		return
	}

	var body []byte
	if err, ok := result.Result.(error); ok {
		body, _ = json.Marshal(&struct {
			ErrMsg string `json:"error_msg"`
		}{err.Error()})
	} else {
		body, _ = json.Marshal(result.Result)
	}
	w.Write(body)
}

func EncodeForwardSuccessResult(w http.ResponseWriter, body []byte) {
	w.WriteHeader(http.StatusOK)
	if body == nil {
		return
	}

	w.Write(body)
}

func ResponceFailedMsg(w http.ResponseWriter, msg string) {
	EncodeResult(w, Failed(errors.New(msg)))
}

func ResponceSuccess(w http.ResponseWriter) {
	EncodeResult(w, Succeed())
}

func ResponceSuccessWithBody(w http.ResponseWriter, result interface{}) {
	EncodeResult(w, SucceedWithResult(result))
}
