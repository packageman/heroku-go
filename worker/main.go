package main

import (
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	pool := newRedisPool()
	log.Printf("starting worker...")
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			conn := pool.Get()
			defer conn.Close()
			job, err := conn.Do("LPOP", "queue")
			if err != nil {
				log.Printf("failed to get job, err: %s", err.Error())
				continue
			}
			if job != nil {
				log.Printf("got job: %s", job)
				continue
			}
			log.Printf("no job found in queue")
		}
	}
}

func newRedisPool() *redis.Pool {
	redisUrl := os.Getenv("REDIS_URL")
	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 12000,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisUrl)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
