package main

import (
	"log"
	"strings"
	"io/ioutil"
	"net/http"
)
const (
	filmFileName = "films.txt"
	passedFileName = "passed.txt"
)
func getFilmNames(filename string) ([]string, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return strings.Split(string(content), "/r/n"), nil
}

func getSiteBody(url, proxy string) ([]byte, error) {
	//
}
func parseSite(body []byte, film string) bool {
	//
}

func main() {
	//
}