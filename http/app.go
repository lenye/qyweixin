package http

import (
	"net"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/golang/glog"

	"github.com/lenye/qyweixin/access-token"
	"github.com/lenye/qyweixin/internal/util"
	"github.com/lenye/qyweixin/internal/version"
)

// QYWeiXinApp 企业微信
type QYWeiXinApp struct {
	sync.RWMutex
	opts              atomic.Value
	startTime         time.Time
	httpListener      net.Listener
	waitGroup         util.WaitGroupWrapper
	accessTokenClient *access_token.AccessTokenClient
}

// NewQYWeiXinApp 创建企业微信
func NewQYWeiXinApp(opt *Option) *QYWeiXinApp {
	p := &QYWeiXinApp{
		startTime: time.Now(),
	}
	p.swapOption(opt)

	glog.Info(version.String("qy-weixin"))

	return p
}

func (p *QYWeiXinApp) getOption() *Option {
	return p.opts.Load().(*Option)
}

func (p *QYWeiXinApp) swapOption(opts *Option) {
	p.opts.Store(opts)
}

// Run 运行
func (p *QYWeiXinApp) Run() {
	ctx := &ContextApp{p}
	opts := p.getOption()

	p.accessTokenClient = access_token.NewAccessTokenClient(opts.HTTPClientConnectTimeout, opts.HTTPClientRequestTimeout)
	p.waitGroup.Wrap(func() {
		p.accessTokenClient.Loop(opts.AppID, opts.AppSecret)
	})

	quitChan := make(chan os.Signal, 1)
	signal.Notify(quitChan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	var httpListener net.Listener
	httpListener, err := net.Listen("tcp", opts.HTTPAddress)
	if err != nil {
		glog.Fatalf("FATAL: listen (%s) failed - %s", opts.HTTPAddress, err)
	}
	p.Lock()
	p.httpListener = httpListener
	p.Unlock()

	httpServer := NewHTTPServer(ctx)
	p.waitGroup.Wrap(func() {
		Serve(p.httpListener, httpServer, "HTTP")
	})

	<-quitChan
	close(p.accessTokenClient.QuitChan)
	if httpListener != nil {
		httpListener.Close()
	}
	glog.Info("exit")
	p.waitGroup.Wait()
}
