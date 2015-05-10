// +build !debug
package webhooks

import (
	"encoding/json"
	"log"
	"os"
	"time"
)

type GithubRepo struct {
	LanguagesUrl string `json:"languages_url"`
}

func GithubCount(username string) map[string]int {
	var repos []GithubRepo
	token := os.Getenv("GITHUB_TOKEN")

	errChan := make(chan error)
	response := make(chan []byte)

	languages := make(map[string]int)

	respnse_string, err := Get("https://api.github.com/users/" + username + "/repos" + "?access_token=" + token)
	if err != nil {
		log.Fatalln(err)
		return languages
	}
	err = json.Unmarshal(respnse_string, &repos)
	if err != nil {
		log.Fatalln(err)
	}
	reposLeft := len(repos)

	for _, repo := range repos {
		go AsyncGet(repo.LanguagesUrl+"?access_token="+token, response, errChan)
	}

	for {
		select {
		case res := <-response:
			hash := parseLanguages(res)
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

func parseLanguages(data []byte) map[string]int {
	var langs map[string]int
	err := json.Unmarshal(data, &langs)
	if err != nil {
		log.Fatalln(err)
	}
	return langs
}
