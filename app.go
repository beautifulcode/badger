// +build !debug
package main

import (
	"github.com/3zcurdia/merithub/parser"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"os"
	"time"
)

func repoCount(username string) map[string]int {
	token := os.Getenv("GITHUB_TOKEN")

	errChan := make(chan error)
	response := make(chan []byte)

	languages := make(map[string]int)

	res, err := parser.Get("https://api.github.com/users/" + username + "/repos" + "?access_token=" + token)
	if err != nil {
		log.Fatalln(err)
		return languages
	}
	repos := parser.ParseRepo(res)
	reposLeft := len(repos)

	for _, repo := range repos {
		go parser.AsyncGet(repo.LanguagesUrl+"?access_token="+token, response, errChan)
	}

	for {
		select {
		case res := <-response:
			hash := parser.ParseMapInt(res)
			for key, value := range hash {
				_, ok := languages[key]
				if ok {
					languages[key] += value
				} else {
					languages[key] = value
				}
			}
			reposLeft--
		case err := <-errChan:
			log.Fatalln(err)
		case <-time.After(1000 * time.Millisecond):
			log.Fatalln("Timeout")
		}
		if reposLeft == 0 {
			break
		}
	}
	return languages
}

func main() {
	r := gin.Default()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, string("Go to /languages/:username"))
	})

	r.GET("/languages/:username", func(c *gin.Context) {
		res := repoCount(c.Params.ByName("username"))
		c.JSON(http.StatusOK, res)
	})

	r.Run(":" + os.Getenv("PORT"))
}
