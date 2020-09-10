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
		useRedisL   bool
		socketRedis bool
		cacheTTL    int
	)
	flag.StringVar(&srvT, "srv", "_http._tcp.mxtoolbox.com", "string of SRV RR")
	flag.BoolVar(&useRedis, "redis", false, "use redis to save result. (default false)")
	flag.BoolVar(&useRedisL, "redisl", false, "use redis to load cache. (default false)")
	flag.BoolVar(&socketRedis, "socket-redis", false, "use unix domain socket to connect with redis. (default false)")
	flag.IntVar(&cacheTTL, "ttl", 20, "TTL of redis-cache")
	flag.Parse()

	if useRedisL {
		if socketRedis {
			rc := redis.NewClient(&redis.Options{
				Network: "unix",
				Addr:    "/var/run/redis/redis.sock",
			})
			rr, err := rc.Get(ctx, srvT).Result()
			if err == redis.Nil {
			} else if err != nil {
				os.Exit(0)
			} else {
				// use just one record, enough.
				fmt.Printf("%v", rr)
				return err
			}
		} else {
			rc := redis.NewClient(&redis.Options{})
			rr, err := rc.Get(ctx, srvT).Result()
			if err == redis.Nil {
			} else if err != nil {
				os.Exit(0)
			} else {
				// use just one record, enough.
				fmt.Printf("%v", rr)
				return err
			}
		}

	}

	rv := net.Resolver{}
	_, srvs, err := rv.LookupSRV(ctx, "", "", srvT)
	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		os.Exit(1)
	}

	for _, srv := range srvs {
		addrs, err := rv.LookupHost(ctx, srv.Target)
		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error()+"\n")
			os.Exit(1)
		}
		// use just one record, enough.
		rr := addrs[0] + ":" + strconv.Itoa(int(srv.Port))
		fmt.Printf("%v", rr)

		if useRedis {
			if socketRedis {
				rc := redis.NewClient(&redis.Options{
					Network: "unix",
					Addr:    "/var/run/redis/redis.sock",
				})
				err = rc.Set(ctx, srvT, rr, time.Duration(cacheTTL)*time.Second).Err()
			} else {
				rc := redis.NewClient(&redis.Options{})
				err = rc.Set(ctx, srvT, rr, time.Duration(cacheTTL)*time.Second).Err()
			}

			if err != nil {
				os.Exit(0)
			}
		}
	}

	return err
}
