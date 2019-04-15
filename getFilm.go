package main

import (
	"log"
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"time"
	"sync"
	"github.com/PuerkitoBio/goquery"
	"my.libs/slkclient"
)
const (
	filmFileName = "films.txt"
	passedFileName = "passed.txt"
	proxy = "http://54.37.84.141:3128"
	address = "http://rutor.info"
	search = "/search/0/0/100/0/"
)
type movie struct {
	name string
	show string
}
func loadAlreadyLoadFilms(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return make([]string, 0)
	}
	return strings.Split(string(content), "\r\n")
}

func loadFilmNamesFromFile(filename string) ([]movie, error) {
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
func saveFilmNamesToFile(films []movie, filename string) {
	tempFilm := make([]string, 0)
	for _, film := range films {
		tempFilm = append(tempFilm, film.name + ";" + film.show)
	}
	buffer := []byte(strings.Join(tempFilm, "\r\n"))
	ioutil.WriteFile(filename, buffer, 0666)
}
func saveAlreadyFilms(list []string, filename string) {
	buffer := []byte(strings.Join(list, "\r\n"))
	ioutil.WriteFile(filename, buffer, 0666)
}

func verifyAlready(pattern string, list []string) bool {
	for _, line := range list {
		if pattern == line {
			return true
		}
	}
	return false
}
func verifyIn (film string, list []string) bool {
	var flag bool
	for _, line := range list {
		flag = true
		for _, word := range strings.Split(film, " ") {
			if !strings.Contains(strings.ToLower(line), strings.ToLower(word)) {
				flag = false
				break
			}
		}
		if flag {
			if strings.Contains(strings.ToLower(line), "лицензия") ||  strings.Contains(strings.ToLower(line), "itunes") {}
			return true
		}
	}
	return false
}
func addAlreadyFilms(count int, list *[]string, c chan []byte, wg *sync.WaitGroup) {
	wg.Add(1)
	var msg string
	defer wg.Done()
	for {
		msg = string(<-c)
		if msg == "" {
			count--
			if count == 0 {
				return
			} else {
				*list = append(*list, msg)
			}
		}
	}
}

func parseSite(film *movie, list []string, c chan []byte) {
	log.Println("starting " + film.name)
	defer func() {c <- []byte("")}()
	proxyURL, err := url.Parse(proxy)
	if err != nil {
		log.Println(err)
	}
	transport := &http.Transport{Proxy: http.ProxyURL(proxyURL)}
	client := &http.Client{Transport: transport, Timeout: time.Second * 20}
	request, err := http.NewRequest("GET", address + search + url.PathEscape(film.name), nil)
    if err != nil {
        log.Println(err)
	}
	response, err := client.Do(request)
    if err != nil {
        log.Println(err)
    }
	
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}
	var sendToSlack bool
	doc.Find("tr.gai td:nth-child(2) a:nth-child(3)").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		item := s.Text()
		if !verifyAlready(item, list) {
			c <- []byte(item)
			if verifyIn(film.name, list) {
				sendToSlack = true
				film.show = "0"
			}
		}
	})
	if sendToSlack {
		err := slkclient.SendToSlack(":movie_camera: Фильм", "Появился фильм: " + film.name, "", "", "")
		if err != nil {
			return
		}
	}
}

func main() {
	newFilms := make(chan []byte, 0)
	var wg sync.WaitGroup
	films, _ := loadFilmNamesFromFile(filmFileName)
	listAlready := loadAlreadyLoadFilms(passedFileName)
	go addAlreadyFilms(len(films), &listAlready, newFilms, &wg)
	for i := range films {
		if films[i].show != "0" {
			log.Println("launch " + films[i].name)
			go parseSite(&films[i], listAlready, newFilms)
		}
	}
	wg.Wait()
	saveAlreadyFilms(listAlready, passedFileName)
	saveFilmNamesToFile(films, filmFileName)
}