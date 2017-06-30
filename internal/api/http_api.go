package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"
)

type APIHandler func(http.ResponseWriter, *http.Request, httprouter.Params) (interface{}, error)

type Decorator func(APIHandler) APIHandler

type Err struct {
	Code int
	Text string
}

func (e Err) Error() string {
	return e.Text
}

func Decorate(f APIHandler, ds ...Decorator) httprouter.Handle {
	decorated := f
	for _, decorate := range ds {
		decorated = decorate(decorated)
	}
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) {
		decorated(w, req, ps)
	}
}

func V1(f APIHandler) APIHandler {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
		data, err := f(w, req, ps)
		if err != nil {
			RespondV1(w, err.(Err).Code, err)
			return nil, nil
		}
		RespondV1(w, 200, data)
		return nil, nil
	}
}

func RespondV1(w http.ResponseWriter, code int, data interface{}) {
	var response []byte
	var err error

	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Content-Type", "application/json; charset=utf-8")

	if code == 200 {
		switch data.(type) {
		case string:
			response = []byte(data.(string))
		case []byte:
			response = data.([]byte)
		case nil:
			response = []byte{}
		default:
			response, err = json.Marshal(data)
			if err != nil {
				code = 500
				data = err
			}
		}
	}

	if code != 200 {
		response = []byte(fmt.Sprintf(`{"message":"%s"}`, data))
	}

	w.WriteHeader(code)
	w.Write(response)
}

func Log(f APIHandler) APIHandler {
	return func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
		start := time.Now()
		response, err := f(w, req, ps)
		elapsed := time.Since(start)
		status := 200
		if e, ok := err.(Err); ok {
			status = e.Code
		}
		glog.Infof("%d %s %s (%s) %s", status, req.Method, req.URL.RequestURI(), req.RemoteAddr, elapsed)
		return response, err
	}
}

func LogPanicHandler() func(w http.ResponseWriter, req *http.Request, p interface{}) {
	return func(w http.ResponseWriter, req *http.Request, p interface{}) {
		glog.Errorf("ERROR: panic in HTTP handler - %s", p)
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{500, "INTERNAL_ERROR"}
		}, V1)(w, req, nil)
	}
}

func LogNotFoundHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{404, "NOT_FOUND"}
		}, V1)(w, req, nil)
	})
}

func LogMethodNotAllowedHandler() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		Decorate(func(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
			return nil, Err{405, "METHOD_NOT_ALLOWED"}
		}, V1)(w, req, nil)
	})
}
