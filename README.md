# 企业微信的access token中控服务和发消息
 
## 命令行参数  
 
 qy-weixin -h 查看参数说明
 
    Usage of qy-weixin:
      -alsologtostderr
            log to standard error as well as files
      -app_id string
            qy weixin app id
      -app_secret string
            qy weixin app secret
      -http_address string
            http server listening (default "127.0.0.1:8000")
      -http_client_connect_timeout duration
            http client connect timeout (default 5s)
      -http_client_request_timeout duration
            http client request timeout (default 5s)
      -log_backtrace_at value
            when logging hits line file:N, emit a stack trace
      -log_dir string
            If non-empty, write log files in this directory
      -logtostderr
            log to standard error instead of files
      -stderrthreshold value
            logs at or above this threshold go to stderr
      -v value
            log level for V logs
      -vmodule value
            comma-separated list of pattern=N settings for file-filtered logging
 
 
### 运行程序:
 
 qy-weixin -app_id=XXXXXXXXXXX -app_secret=XXXXXXXXXXX -logtostderr
   
## 取access-token

    GET /wx/qy/access-token

    {
        "access_token": "qdjWc6kix6RrYfgUFNhDxdcC4EdzXLIUGlKFANJdHvTn5WcxyTyGtJeM2nZPEeHP1SxRWkNj8uqVXp4OIoavCxAD8h_WnR120bv2wDJSOcvfKV8OQaPzjUiI4u6uaelQsi_zOtOhdiFkwgzSeTCcRYWrovmn7KTONcNu-0qPC5Yr8y15FZHM0ol7uuiLocKDO0AMo5jNhBnj2MH1nsfX7xo1sbhyqFju04T7GTRckdko4xtxh8muMteMGAiBB0xNaM4jJHGBWakaaxXMnZgz4MNdb323GELWZDglcoXl8wg",
        "expires_in": 7200,
        "create_at": "2017-06-22T14:56:40.3631589+08:00"
    }

## 发消息

 企业微信的[消息类型及数据格式](http://qydev.weixin.qq.com/wiki/index.php?title=消息类型及数据格式)

    POST /wx/qy/send/message

    {
       "totag" : "1",
       "msgtype" : "text",
       "agentid" : 1000002,
       "text" : {
           "content" : "你的快递已到，请携带工卡前往邮件中心领取。\n出发前可查看<a href=\"http://work.weixin.qq.com\">邮件中心视频实况</a>，聪明避开排队。"
       }
    }

## License

 qyweixin is licensed under the [Apache License 2.0](https://github.com/lenye/qyweixin/blob/master/LICENSE).