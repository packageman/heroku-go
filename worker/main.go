package main

import (
	"log"
	"os"
	"time"

	"github.com/gomodule/redigo/redis"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	log.Printf("starting worker...")
	ticker := time.NewTicker(10 * time.Second)
	for {
		select {
		case <-ticker.C:
			processJob()
		}
	}
}

func processJob() {
	c, err := redis.DialURL(os.Getenv("REDIS_URL"))
	if err != nil {
		log.Printf("failed to get redis connection, err: %s", err.Error())
	}
	defer c.Close()
	job, err := c.Do("LPOP", "queue")
	if err != nil {
		log.Printf("failed to get job, err: %s", err.Error())
		return
	}
	if job != nil {
		log.Printf("got job: %s", job)
		return
	}
	log.Printf("no job found in queue")
}
