package message

import (
	"bytes"
	"fmt"
	"time"

	"github.com/lenye/qyweixin/internal/http"
)

const sendMessageURL = "https://qyapi.weixin.qq.com/cgi-bin/message/send?access_token=%s"

// SendMessage 发送消息. 应用支持推送文本、图片、视频、文件、图文等类型
func SendMessage(body *bytes.Buffer, access_token string, HTTPClientConnectTimeout, HTTPClientRequestTimeout time.Duration) ([]byte, error) {
	return http.NewHttpClient(HTTPClientConnectTimeout, HTTPClientRequestTimeout).HTTPPostJSON(fmt.Sprintf(sendMessageURL, access_token), body)
}
