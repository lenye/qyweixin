package access_token

import (
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"

	"github.com/lenye/qyweixin/internal/api"
)

const (
	retryInterval = 10 * 1000                                                              //每隔10秒重试
	tokenURL      = "https://qyapi.weixin.qq.com/cgi-bin/gettoken?corpid=%s&corpsecret=%s" //企业微信token url
)

//凭证
type AccessToken struct {
	Ticket    string    `json:"access_token"` //凭证
	ExpiresIn int64     `json:"expires_in"`   //凭证有效时间，单位：秒
	NextGet   int64     `json:"-"`            //下次取凭证时间
	CreateAt  time.Time `json:"create_at"`    //取得凭证的时间
}

//ResponseError 响应操作错误信息
type ResponseError struct {
	ErrorCode    int64  `json:"errcode"`
	ErrorMessage string `json:"errmsg"`
}

type AccessTokenClient struct {
	ticket   atomic.Value
	Client   *api.HttpClient
	QuitChan chan int
}

func NewAccessTokenClient(connectTimeout time.Duration, requestTimeout time.Duration) *AccessTokenClient {
	p := &AccessTokenClient{
		Client:   api.NewHttpClient(connectTimeout, requestTimeout),
		QuitChan: make(chan int),
	}
	p.SwapTicket(&AccessToken{})
	return p
}

//Load
func (p *AccessTokenClient) Load() *AccessToken {
	return p.ticket.Load().(*AccessToken)
}

//swapTicket
func (p *AccessTokenClient) SwapTicket(ticket *AccessToken) {
	p.ticket.Store(ticket)
}

//getAccessToken
func (p *AccessTokenClient) getAccessToken(appId, appSecret string) (*AccessToken, error) {
	accessToken := p.Load()
	//清除过期access-token
	if accessToken.Ticket != "" && int64(time.Now().Sub(accessToken.CreateAt).Seconds()) >= accessToken.ExpiresIn {
		accessToken.Ticket = ""
		accessToken.ExpiresIn = 0
		accessToken.NextGet = 0
		accessToken.CreateAt = time.Now()
		p.SwapTicket(accessToken)
	}
	glog.Infof("old access-token=%+v", accessToken)

	respBody, err := p.Client.HTTPGet(fmt.Sprintf(tokenURL, appId, appSecret))
	if err != nil {
		accessToken.NextGet = retryInterval
		return accessToken, errors.Wrap(err, "getAccessToken failed")
	}

	var response struct {
		ResponseError
		AccessToken
	}
	err = json.Unmarshal(respBody, &response)
	if err != nil {
		accessToken.NextGet = retryInterval
		return accessToken, errors.Wrap(err, "getAccessToken json.Unmarshal failed")
	}

	if response.ErrorCode != 0 {
		accessToken.NextGet = retryInterval
		return accessToken, errors.Wrapf(err, "getAccessToken response errcode: %d, errmsg: %s", response.ErrorCode, response.ErrorMessage)
	}

	response.CreateAt = time.Now()
	//刷新策略
	switch {
	case response.ExpiresIn >= 60*60:
		response.NextGet = (response.ExpiresIn - 30*60) * 1000
	case response.ExpiresIn >= 30*60:
		response.NextGet = (response.ExpiresIn - 10*60) * 1000
	case response.ExpiresIn >= 10*60:
		response.NextGet = (response.ExpiresIn - 60) * 1000
	case response.ExpiresIn <= 6:
		response.NextGet = 100
	default:
		response.NextGet = (response.ExpiresIn - 6) * 1000
	}
	glog.Infof("new access-token=%+v", response.AccessToken)

	accessToken.Ticket = response.Ticket
	accessToken.ExpiresIn = response.ExpiresIn
	accessToken.NextGet = response.NextGet
	accessToken.CreateAt = response.CreateAt
	p.SwapTicket(accessToken)
	return accessToken, nil
}

//Loop
func (p *AccessTokenClient) Loop(appId, appSecret string) {
	var refreshInterval time.Duration

	newAccessToken, err := p.getAccessToken(appId, appSecret)
	if err != nil {
		glog.Error(err)
	}
	refreshInterval = time.Duration(newAccessToken.NextGet) * time.Millisecond
	glog.Infof("next access-token time.NewTicker=%v", refreshInterval)
	timeTicker := time.NewTicker(refreshInterval)

	for {
		select {
		case <-timeTicker.C:
			newAccessToken, err := p.getAccessToken(appId, appSecret)
			if err != nil {
				glog.Error(err)
			}
			refreshInterval = time.Duration(newAccessToken.NextGet) * time.Millisecond
			glog.Infof("next access-token time.NewTicker=%v", refreshInterval)
			timeTicker.Stop()
			timeTicker = time.NewTicker(refreshInterval)
		case <-p.QuitChan:
			goto exit
		}
	}

exit:
	glog.Info("exiting access-token Loop")
	timeTicker.Stop()
}
