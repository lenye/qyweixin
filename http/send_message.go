package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"

	"github.com/lenye/qyweixin/message"
)

// 发送消息
func (s *HttpServer) sendMessage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	opts := s.ctx.app.getOption()
	readMax := s.ctx.app.getOption().MaxMsgSize + 1
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, readMax))
	if err != nil {
		glog.Errorf("failed to ReadFrom request body - %s", err.Error())
		return nil, Err{http.StatusInternalServerError, "INTERNAL_ERROR"}
	}
	if int64(len(body)) == readMax {
		return nil, Err{http.StatusRequestEntityTooLarge, "MSG_TOO_BIG"}
	}
	if len(body) == 0 {
		return nil, Err{http.StatusBadRequest, "MSG_EMPTY"}
	}
	buf := bytes.NewBuffer(body)
	return message.SendMessage(buf, s.ctx.app.accessTokenClient.Load().Ticket, opts.HTTPClientConnectTimeout, opts.HTTPClientRequestTimeout)
}
