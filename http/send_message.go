package http

import (
	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/golang/glog"
	"github.com/julienschmidt/httprouter"

	api "github.com/lenye/qyweixin/internal/http"
	"github.com/lenye/qyweixin/message"
)

// 发送消息
func (s *HttpServer) sendMessage(w http.ResponseWriter, req *http.Request, ps httprouter.Params) (interface{}, error) {
	opts := s.ctx.app.getOption()
	readMax := s.ctx.app.getOption().MaxMsgSize + 1
	body, err := ioutil.ReadAll(io.LimitReader(req.Body, readMax))
	if err != nil {
		glog.Errorf("sendMessage: failed to ReadFrom request body - %s", err.Error())
		return nil, api.Err{Code: http.StatusInternalServerError, Text: "INTERNAL_ERROR"}
	}
	msgSize := int64(len(body))
	if msgSize == readMax {
		glog.Errorf("sendMessage: message too big, message size=%d, max size=%d", msgSize, readMax)
		return nil, api.Err{Code: http.StatusRequestEntityTooLarge, Text: "MSG_TOO_BIG"}
	}
	if msgSize == 0 {
		glog.Error("sendMessage: message empty")
		return nil, api.Err{Code: http.StatusBadRequest, Text: "MSG_EMPTY"}
	}
	buf := bytes.NewBuffer(body)
	return message.SendMessage(buf, s.ctx.app.accessTokenClient.Load().Ticket, opts.HTTPClientConnectTimeout, opts.HTTPClientRequestTimeout)
}
