package main

import (
	"flag"
	"time"

	"github.com/golang/glog"

	"github.com/lenye/qyweixin/http"
)

var opts http.Option

func init() {
	flag.StringVar(&opts.HTTPAddress, "http_address", "127.0.0.1:8000", "http server listening")
	flag.DurationVar(&opts.HTTPClientConnectTimeout, "http_client_connect_timeout", 5*time.Second, "http client connect timeout")
	flag.DurationVar(&opts.HTTPClientRequestTimeout, "http_client_request_timeout", 5*time.Second, "http client request timeout")

	flag.StringVar(&opts.AppID, "app_id", "", "qy weixin app id")
	flag.StringVar(&opts.AppSecret, "app_secret", "", "qy weixin app secret")
}

func main() {
	flag.Parse()
	defer glog.Flush()

	if opts.AppID == "" || opts.AppSecret == "" {
		glog.Exit("missing app-id or app-secret")
	}
	opts.MaxMsgSize = 3 * 1024
	app := http.NewQYWeiXinApp(&opts)
	app.Run()
}
