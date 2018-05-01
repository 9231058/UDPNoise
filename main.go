/*
 * +===============================================
 * | Author:        Parham Alvani <parham.alvani@gmail.com>
 * |
 * | Creation Date: 01-05-2018
 * |
 * | File Name:     main.go
 * +===============================================
 */

package main

import (
	"net/http"

	"github.com/aut-ceit/UDPNoise/udpnoise"
	"github.com/gin-gonic/gin"
)

var projects map[string]*udpnoise.UDPNoise

func init() {
	projects = make(map[string]*udpnoise.UDPNoise)
}

// handle registers apis and create http handler
func handle() http.Handler {
	r := gin.Default()

	api := r.Group("/api")
	{
		api.GET("/about", aboutHandler)

		api.POST("/project", projectNewHandler)
		api.GET("/project", projectListHandler)
		api.DELETE("/project/:name", projectRemoveHandler)
	}

	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{"error": "404 Not Found"})
	})

	return r
}

func main() {
}

func aboutHandler(c *gin.Context) {
	c.String(http.StatusOK, "18.20 is leaving us")
}

func projectNewHandler(c *gin.Context) {
	var json struct {
		Destination string
		Loss        int
		Name        string
	}
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	u, err := udpnoise.New(json.Loss, json.Destination)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	projects[json.Name] = u

	go u.Run()

	c.JSON(http.StatusOK, u)
}

func projectListHandler(c *gin.Context) {
	c.JSON(http.StatusOK, projects)
}

func projectRemoveHandler(c *gin.Context) {
	name := c.Param("name")

	if p, ok := projects[name]; ok {
		if err := p.Close(); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, name)
}
