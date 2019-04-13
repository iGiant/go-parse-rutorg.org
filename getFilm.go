package main

import (
	// "log"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
)
const (
	filmFileName = "films.txt"
	passedFileName = "passed.txt"
	proxy = "http://54.37.84.141:3128"
)
type film struct {
	name string
	show string
}
func getFilmNamesFromFile(filename string) ([]film, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result := make([]film, 0)
	var tempFilm []string
	films := strings.Split(string(content), "\r\n")
	
	for _, name := range films {
		tempFilm = strings.Split(name, ";")
		result = append(result, film {name: tempFilm[0], show: tempFilm[1]})
	}
	return result, nil
}
func setFilmNamesToFile(films []film, filename string) {
	tempFilm := make([]string, 0)
	for _, kino := range films {
		tempFilm = append(tempFilm, kino.name + ";" + kino.show)
	}
	buffer := []byte(strings.Join(tempFilm, "\r\n"))
	ioutil.WriteFile(filename, buffer, 0666)
}
func worker



func getSiteBody(address string) ([]byte, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return nil, err
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport, Timeout: time.Second * 5}
	request, err := http.NewRequest("GET", address, nil)
    if err != nil {
        return nil, err
    }
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }
    data, err := ioutil.ReadAll(response.Body)
    if err != nil {
        return nil, err
	}
	return data, nil
}
func parseSite(body []byte, film string) bool {
	
	return true
}

func main() {
	films, _ := getFilmNamesFromFile(filmFileName)

	setFilmNamesToFile(films, filmFileName)
}