package http

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/golang/glog"
)

func Serve(listener net.Listener, handler http.Handler, proto string) {
	glog.Infof("%s: listening on %s", proto, listener.Addr())

	server := &http.Server{
		Handler:      handler,
		ReadTimeout:  20 * time.Second,
		WriteTimeout: 20 * time.Second,
	}
	err := server.Serve(listener)
	// theres no direct way to detect this error because it is not exposed
	if err != nil && !strings.Contains(err.Error(), "use of closed network connection") {
		glog.Errorf("ERROR: http.Serve() - %s", err)
	}
	glog.Infof("%s: closing %s", proto, listener.Addr())
}
