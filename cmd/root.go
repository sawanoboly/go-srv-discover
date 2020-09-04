package cmd

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/go-redis/redis"
)

var ctx = context.Background()

// Execute is wrapped func main()
func Execute() error {
	var (
		srvT        string
		useRedis    bool
		socketRedis bool
	)
	flag.StringVar(&srvT, "srv", "_http._tcp.mxtoolbox.com", "string of SRV RR")
	flag.BoolVar(&useRedis, "redis", false, "use redis to save result. (default false)")
	flag.BoolVar(&socketRedis, "socket-redis", false, "use unix domain socket to connect with redis. (default false)")
	flag.Parse()

	rv := net.Resolver{}
	_, srvs, err := rv.LookupSRV(ctx, "", "", srvT)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}

	for _, srv := range srvs {
		addrs, err := rv.LookupHost(ctx, srv.Target)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
			os.Exit(1)
		}
		// use just one record, enough.
		rr := addrs[0] + ":" + strconv.Itoa(int(srv.Port))
		fmt.Printf("%v\n", rr)

		if useRedis {
			if socketRedis {
				rc := redis.NewClient(&redis.Options{
					Network: "unix",
					Addr:    "/var/run/redis/redis.sock",
				})
				err = rc.Set(ctx, srvT, rr, time.Second*10).Err()
			} else {
				rc := redis.NewClient(&redis.Options{})
				err = rc.Set(ctx, srvT, rr, time.Second*10).Err()
			}

			if err != nil {
				panic(err)
			}
		}
	}

	return err
}
