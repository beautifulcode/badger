// +build !debug
package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"
)

type Repo struct {
	LanguagesUrl string `json:"languages_url"`
}

func asyncGet(url string, response chan []byte, errorChannel chan error) {
	// start := time.Now()
	http_response, err := http.Get(url)
	// elapsed := time.Since(start)
	// fmt.Println("Get: ", elapsed)
	if err != nil {
		errorChannel <- err
		return
	}
	defer http_response.Body.Close()
	if http_response.StatusCode != 200 {
		errorChannel <- errors.New("API response:" + http_response.Status + " url:" + url)
		return
	}
	data, err := ioutil.ReadAll(http_response.Body)
	if err != nil {
		errorChannel <- err
		return
	}
	response <- []byte(data)
}

func Get(url string) ([]byte, error) {
	http_response, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer http_response.Body.Close()
	if http_response.StatusCode != 200 {
		return nil, errors.New("API response:" + http_response.Status + " url:" + url)
	}
	data, err := ioutil.ReadAll(http_response.Body)
	if err != nil {
		return nil, err
	}
	return []byte(data), nil
}

func parseMapInt(data []byte) map[string]int {
	var hash map[string]int
	err := json.Unmarshal(data, &hash)
	if err != nil {
		log.Fatalln(err)
	}
	return hash
}

func parseRepo(data []byte) []Repo {
	var array []Repo
	err := json.Unmarshal(data, &array)
	if err != nil {
		log.Fatalln(err)
	}
	return array
}

func repoCount(username string) map[string]int {
	token := os.Getenv("GITHUB_TOKEN")

	errChan := make(chan error)
	response := make(chan []byte)

	languages := make(map[string]int)

	res, err := Get("https://api.github.com/users/" + username + "/repos" + "?" + token)
	if err != nil {
		log.Fatalln(err)
		return languages
	}
	repos := parseRepo(res)
	reposLeft := len(repos)

	for _, repo := range repos {
		go asyncGet(repo.LanguagesUrl+"?"+token, response, errChan)
	}

	for {
		select {
		case res := <-response:
			hash := parseMapInt(res)
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
