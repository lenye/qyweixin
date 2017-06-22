package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

type HttpServer struct {
	ctx    *ContextApp
	router http.Handler
}

func NewHTTPServer(ctx *ContextApp) *HttpServer {
	httpStatusLog := Log()
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.PanicHandler = LogPanicHandler()
	router.NotFound = LogNotFoundHandler()
	router.MethodNotAllowed = LogMethodNotAllowedHandler()
	s := &HttpServer{
		ctx:    ctx,
		router: router,
	}

	router.Handle("POST", "/wx/qy/send/message", Decorate(s.sendMessage, httpStatusLog, V1))

	return s
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
