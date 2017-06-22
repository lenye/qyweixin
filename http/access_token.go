package http

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// 凭证
func (s *HttpServer) accessToken(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	return s.ctx.app.accessTokenClient.Load(), nil
}
