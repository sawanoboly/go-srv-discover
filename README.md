# go-srv-discover

Returns 1 IPv4 address and port.


```
$ gsr --help
  -redis
        use redis to save result. (default false)
  -socket-redis
        use unix domain socket to connect with redis. (default false)
  -srv string
        string of SRV RR (default "_http._tcp.mxtoolbox.com")
  -ttl int
        TTL of redis-cache (default 20)
```


```
$ gsr -srv service.srv.local
xxx.xxx.xxx.xxx:80 # without linefeed.
```
