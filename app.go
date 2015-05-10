// +build !debug
package main

import (
	"github.com/3zcurdia/merithub/webhooks"
	"github.com/gin-gonic/gin"
	"net/http"
	"os"
)

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, string("Go to /languages/:username"))
	})

	r.GET("/languages/:username", func(c *gin.Context) {
		res := webhooks.GithubCount(c.Params.ByName("username"))
		c.JSON(http.StatusOK, res)
	})

	r.Run(":" + os.Getenv("PORT"))
}
