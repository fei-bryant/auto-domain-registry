package main

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.LoadHTMLGlob("./*.html")
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "submit.html", gin.H{})
	})
	r.POST("/api/v1/registry", func(c *gin.Context) {
		body := c.PostForm("domains")
		domains := strings.Split(body, ",")
		if err := DomainRegister(domains); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"message": err.Error()})
		}
		c.HTML(http.StatusOK, "success.html", gin.H{})
	})
	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
