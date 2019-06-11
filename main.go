package main

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

type comment struct {
	Message   string
	CreatedAt time.Time
}

func main() {
	router := gin.Default()
	router.LoadHTMLFiles("templates/index.tmpl")
	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
			"comments": []comment{
				comment{"deneme", time.Now()},
				comment{"deneme2", time.Now()},
			},
		})
	})
	router.Run(":8080")
}
