package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"net/http"
	"os"
)

func github_badger(username string) map[string]int {
	token := os.Getenv("GITHUB_TOKEN")
	languages := make(map[string]int)

	response, err := http.Get("https://api.github.com/users/" + username + "/repos?access_token=" + token)
	if err != nil {
		panic(err)
	}
	defer response.Body.Close()
	if response.StatusCode == 200 {
		content, err := ioutil.ReadAll(response.Body)
		if err != nil {
			panic(err)
		}
		fmt.Println(string(content))
	}

	return languages
}

func main() {
	r := gin.Default()
	r.GET("/github/:username", func(c *gin.Context) {
		github_badger(c.Params.ByName("username"))
		c.String(http.StatusOK, "OK")
	})

	r.Run(":" + os.Getenv("PORT"))
}
