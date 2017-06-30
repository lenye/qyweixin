package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"

	"github.com/lenye/qyweixin/internal/api"
)

type HttpServer struct {
	ctx    *ContextApp
	router http.Handler
}

func NewHTTPServer(ctx *ContextApp) *HttpServer {
	router := httprouter.New()
	router.HandleMethodNotAllowed = true
	router.PanicHandler = api.LogPanicHandler()
	router.NotFound = api.LogNotFoundHandler()
	router.MethodNotAllowed = api.LogMethodNotAllowedHandler()
	s := &HttpServer{
		ctx:    ctx,
		router: router,
	}

	router.Handle("GET", "/wx/qy/access-token", api.Decorate(s.accessToken, api.Log, api.V1))
	router.Handle("POST", "/wx/qy/send/message", api.Decorate(s.sendMessage, api.Log, api.V1))

	return s
}

func (s *HttpServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	s.router.ServeHTTP(w, req)
}
