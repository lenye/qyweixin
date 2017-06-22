package http

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/http/httputil"
	"time"

	"github.com/golang/glog"
	"github.com/pkg/errors"
)

const (
	UserAgent        = "Mozilla/5.0 (Windows NT 6.1; Win64; x64; rv:54.0) Gecko/20100101 Firefox/54.0"
	HTTPBodyTypeJSON = "application/json"
)

type deadlinedConn struct {
	Timeout time.Duration
	net.Conn
}

func (c *deadlinedConn) Read(b []byte) (n int, err error) {
	return c.Conn.Read(b)
}

func (c *deadlinedConn) Write(b []byte) (n int, err error) {
	return c.Conn.Write(b)
}

// A custom http.Transport with support for deadline timeouts
func NewDeadlineTransport(connectTimeout time.Duration, requestTimeout time.Duration) *http.Transport {
	transport := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		Dial: func(netw, addr string) (net.Conn, error) {
			c, err := net.DialTimeout(netw, addr, connectTimeout)
			if err != nil {
				return nil, err
			}
			return &deadlinedConn{connectTimeout, c}, nil
		},
		ResponseHeaderTimeout: requestTimeout,
		MaxIdleConns:          100,
		IdleConnTimeout:       8 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}
	return transport
}

type HttpClient struct {
	c *http.Client
}

func NewHttpClient(connectTimeout time.Duration, requestTimeout time.Duration) *HttpClient {
	transport := NewDeadlineTransport(connectTimeout, requestTimeout)
	return &HttpClient{
		c: &http.Client{
			Transport: transport,
			Timeout:   requestTimeout,
		},
	}
}

func (p *HttpClient) HTTPGet(url string) ([]byte, error) {
	var reqDump, respDump []byte
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.Header.Set("User-Agent", UserAgent)

	reqDump, err = httputil.DumpRequest(req, true)
	if err != nil {
		glog.V(4).Info(errors.Wrap(err, "HTTPGet DumpRequest"))
	}

	httpResp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respDump, err = httputil.DumpResponse(httpResp, true)
	if err != nil {
		glog.V(4).Info(errors.Wrap(err, "HTTPGet DumpResponse"))
	}
	glog.V(4).Infof("---------- request -----------\n%s\n---------- response ----------\n%s", string(reqDump), string(respDump))

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	return ioutil.ReadAll(httpResp.Body)
}

func (p *HttpClient) HTTPPostJSON(url string, body *bytes.Buffer) ([]byte, error) {
	var reqDump, respDump []byte
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, err
	}
	req.Proto = "HTTP/1.1"
	req.ProtoMajor = 1
	req.ProtoMinor = 1
	req.Header.Set("User-Agent", UserAgent)
	req.Header.Set("Content-Type", HTTPBodyTypeJSON)
	req.Header.Set("Accept", HTTPBodyTypeJSON)

	reqDump, err = httputil.DumpRequest(req, true)
	if err != nil {
		glog.V(4).Info(errors.Wrap(err, "HTTPGet DumpRequest"))
	}

	httpResp, err := p.c.Do(req)
	if err != nil {
		return nil, err
	}
	defer httpResp.Body.Close()

	respDump, err = httputil.DumpResponse(httpResp, true)
	if err != nil {
		glog.V(4).Info(errors.Wrap(err, "HTTPGet DumpResponse"))
	}
	glog.V(4).Infof("---------- request -----------\n%s\n---------- response ----------\n%s", string(reqDump), string(respDump))

	if httpResp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("http.Status: %s", httpResp.Status)
	}

	return ioutil.ReadAll(httpResp.Body)
}
