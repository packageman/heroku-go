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

	router.GET("/enqueue/:message", func(c *gin.Context) {
		message := c.Param("message")
		conn := pool.Get()
		defer conn.Close()
		_, err := conn.Do("RPUSH", "queue", message)
		if err != nil {
			c.String(http.StatusInternalServerError, "failed, error: %s", err.Error())
		}
		log.Printf("enqueued message: %s", message)
		c.String(http.StatusOK, "enqueued message: %s", message)
	})

	router.Run(":" + port)
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
