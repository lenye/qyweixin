package http

import (
	"time"
)

// Option 配置
type Option struct {
	HTTPAddress              string
	HTTPClientConnectTimeout time.Duration
	HTTPClientRequestTimeout time.Duration

	AppID     string
	AppSecret string

	MaxMsgSize int64
}
