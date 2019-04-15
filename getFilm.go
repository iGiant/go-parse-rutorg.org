package main

import (
	"log"
	"os"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"github.com/PuerkitoBio/goquery"
	_ "my.libs/slkclient"
)
const (
	filmFileName = "films.txt"
	passedFileName = "passed.txt"
	proxy = "http://210.0.128.58:8080"
	address = "http://rutor.info"
	search = "/search/0/0/100/0/"
)
type movie struct {
	name string
	show string
}
func getAlreadyLoadFilms(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return make([]string, 0)
	}
	return strings.Split(string(content), "\r\n")
}

func getFilmNamesFromFile(filename string) ([]movie, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result := make([]movie, 0)
	var tempFilm []string
	films := strings.Split(string(content), "\r\n")
	
	for _, film := range films {
		tempFilm = strings.Split(film, ";")
		result = append(result, movie {name: tempFilm[0], show: tempFilm[1]})
	}
	return result, nil
}
func setFilmNamesToFile(films []movie, filename string) {
	tempFilm := make([]string, 0)
	for _, film := range films {
		tempFilm = append(tempFilm, film.name + ";" + film.show)
	}
	buffer := []byte(strings.Join(tempFilm, "\r\n"))
	ioutil.WriteFile(filename, buffer, 0666)
}
func worker(name string) {
	//
}

func verifyIn(pattern string, list []string) bool {
	for _, line := range list {
		if pattern == line {
			return true
		}
	}
	return false
}
func parseSite(name string, list []string) (bool, error) {
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		return false, err
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport, Timeout: time.Second * 10}
	
	request, err := http.NewRequest("GET", address + search + url.PathEscape(name), nil)
    if err != nil {
        return false, err
	}
    response, err := client.Do(request)
    if err != nil {
        return false, err
    }
	defer response.Body.Close()
	// log.Println(response.StatusCode)
	// text, _ := ioutil.ReadAll(response.Body)
	// log.Println(string(text))
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return false, err
	}
	
	doc.Find("tr.gai td:nth-child(2) a:nth-child(3)").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		item := s.Text()
		log.Println(item)
	})
	return true, nil
}

func main() {
	c := make(chan []string, 0)
	films, _ := getFilmNamesFromFile(filmFileName)
	list := getAlreadyLoadFilms(passedFileName)
	for _, 
	_, err := parseSite("Стекло", list)
	log.Println(err)
	setFilmNamesToFile(films, filmFileName)
}