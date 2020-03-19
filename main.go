package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gomodule/redigo/redis"
	_ "github.com/heroku/x/hmetrics/onload"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}
	pool := newRedisPool()

	router := gin.New()
	router.Use(gin.Logger())
	router.LoadHTMLGlob("templates/*.tmpl.html")
	router.Static("/static", "static")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl.html", nil)
	})

	router.GET("/enqueue/:job", func(c *gin.Context) {
		job := c.Param("job")
		conn := pool.Get()
		defer conn.Close()
		_, err := conn.Do("RPUSH", "queue", job)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed, error: %s", err.Error())
		}
		log.Printf("enqueued job: %s", job)
		c.String(http.StatusOK, "enqueued job: %s", job)
	})

	router.Run(":" + port)
}

func newRedisPool() *redis.Pool {
	redisUrl := os.Getenv("REDIS_URL")
	return &redis.Pool{
		MaxIdle:   5,
		MaxActive: 10,
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialURL(redisUrl)
			if err != nil {
				panic(err.Error())
			}
			return c, err
		},
	}
}
