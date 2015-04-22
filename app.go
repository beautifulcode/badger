package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type Repo struct {
	LanguagesUrl string `json:"languages_url"`
}

func getJSONArray(url string) ([]Repo, error) {
	var response_map []Repo

	response, err := http.Get(url)
	if err != nil {
		return response_map, err
	}
	defer response.Body.Close()

	data, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return response_map, err
	}
	err = json.Unmarshal([]byte(data), &response_map)
	return response_map, err
}

func github_badger(username string) map[string]int {
	// token := os.Getenv("GITHUB_TOKEN")
	languages := make(map[string]int)

	repos, err := getJSONArray("https://api.github.com/users/" + username + "/repos")
	if err != nil {
		panic(err)
	}
	for _, repo := range repos {
		fmt.Println(repo)
	}

	return languages
}

func main() {
	github_badger("3zcurdia")
	// r := gin.Default()
	// r.GET("/github/:username", func(c *gin.Context) {
	// 	github_badger(c.Params.ByName("username"))
	// 	c.String(http.StatusOK, "OK")
	// })

	// r.Run(":" + os.Getenv("PORT"))
}
