package main

import (
	"encoding/json"
	"errors"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Repo struct {
	LanguagesUrl string `json:"languages_url"`
}

func getREST(url string) ([]byte, error) {
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

func getHashNumber(url string) (map[string]int, error) {
	var hash map[string]int
	res, err := getREST(url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &hash)
	return hash, nil
}

func getArray(url string) ([]Repo, error) {
	var array []Repo
	res, err := getREST(url)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(res, &array)
	return array, nil
}

func badger(username string) map[string]int {
	token := "access_token=" + os.Getenv("GITHUB_TOKEN")
	languages := make(map[string]int)

	repos, err := getArray("https://api.github.com/users/" + username + "/repos" + "?" + token)
	if err != nil {
		log.Fatalln(err)
	}
	for _, repo := range repos {
		hash, err := getHashNumber(repo.LanguagesUrl + "?" + token)
		if err != nil {
			log.Fatalln(err)
		}
		for key, value := range hash {
			_, ok := languages[key]
			if ok {
				languages[key] += value
			} else {
				languages[key] = value
			}
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
		res := badger(c.Params.ByName("username"))
		c.JSON(http.StatusOK, res)
	})

	r.Run(":" + os.Getenv("PORT"))
}
