// +build !debug
package parser

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
)

type Repo struct {
	LanguagesUrl string `json:"languages_url"`
}

func AsyncGet(url string, response chan []byte, errorChannel chan error) {
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

func ParseMapInt(data []byte) map[string]int {
	var hash map[string]int
	err := json.Unmarshal(data, &hash)
	if err != nil {
		log.Fatalln(err)
	}
	return hash
}

func ParseRepo(data []byte) []Repo {
	var array []Repo
	err := json.Unmarshal(data, &array)
	if err != nil {
		log.Fatalln(err)
	}
	return array
}
