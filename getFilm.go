package main

import (
	"strings"
	"io/ioutil"
	"net/http"
	"net/url"
	"sync"
	"time"
	"github.com/PuerkitoBio/goquery"
	"my.libs/slkclient"
)

const (
	filmFileName = "films.txt"
	passedFileName = "passed.txt"
	address = "http://rutor.info"
	search = "/search/0/0/100/0/"
	ua = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) " +
		 "Chrome/73.0.3683.46 Safari/537.36 OPR/60.0.3255.8 (Edition beta)"
)

type movie struct {
	name string
	show string
}

var wg sync.WaitGroup

// loadAlreadyLoadFilms загружает данные из 'filename' о ранее найденных фильмах для их игнорирования
func loadAlreadyLoadFilms(filename string) []string {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return make([]string, 0)
	}
	return strings.Split(string(content), "\r\n")
}

// loadFilmNamesFromFile загружает данные из файла 'filename'
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

// saveFilmNamesToFile пересохраняет список искомых фильмов 'films' в файл 'filename' с учетом найденных
func saveFilmNamesToFile(films []movie, filename string) {
	tempFilm := make([]string, 0)
	for _, film := range films {
		tempFilm = append(tempFilm, film.name + ";" + film.show)
	}
	buffer := []byte(strings.Join(tempFilm, "\r\n"))
	ioutil.WriteFile(filename, buffer, 0666)
}
// saveAlreadyFilms пересохраняет список всех найденных вариантов 'list' в файл 'filename' для последующего игнорирования
func saveAlreadyFilms(list []string, filename string) {
	buffer := []byte(strings.Join(list[1:], "\r\n"))
	ioutil.WriteFile(filename, buffer, 0666)
}

// verifyAlready проверяет текущий найденный фильм 'film' со списом ранее найденных 'list'
func verifyAlready(film string, list []string) bool {
	for _, line := range list {
		if film == line {
			return true
		}
	}
	return false
}

// verifyIn проверяет, соответствует ли найденных фильм 'line' критериям поиска 'film', заданным нами
func verifyIn (film string, line string) bool {
	flag:= true
	for _, word := range strings.Split(film, " ") {
		if !strings.Contains(strings.ToLower(line), strings.ToLower(word)) {
			flag = false
			break
		}
	}
	if flag {
		if strings.Contains(strings.ToLower(line), "лицензия") ||  strings.Contains(strings.ToLower(line), "itunes") {
		return true
		}
	}
	return false
}
// addAlreadyFilms добавляет найденный фильм 'c' в список ранее найденных фильмов 'list'. 'count' - количество горутин, которые парсят сайт
// для завершения программы после окончания их работы
func addAlreadyFilms(count int, list []string, c <-chan string) {
	defer wg.Done()
	startLen := len(list)
	for {
		msg := <- c
		if string(msg) == "-" {
			count--
			if count == 0 {
				break
			}
		} else {
			list = append(list, msg)
		}
	}
	if len(list) > startLen {
		saveAlreadyFilms(list, passedFileName)
	}
}

// parseSite парсит сайт 'address' на появление нужного нам фильма 'film', выполняя проверку на появление его раньше
// (список 'list'). Найденное описание отправляется через 'c' в горутину addAlreadyFilms
func parseSite(film *movie, list []string, c chan<- string) {
	defer func() {c <- "-"}()
	client := &http.Client{Timeout: time.Second * 5}
	request, err := http.NewRequest("GET", address + search + url.PathEscape(film.name), nil)
	request.Header.Set("User-Agent", ua)
    if err != nil {
        return
	}
	request.Header.Set("User-Agent", ua)
	response, err := client.Do(request)
	if err != nil {
        return
    }
	defer response.Body.Close()
	doc, err := goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return
	}
	sendToSlack := false
	doc.Find("tr.gai td:nth-child(2) a:nth-child(3)").Each(func(i int, s *goquery.Selection) {
		// For each item found, get the band and title
		item := strings.TrimSpace(s.Text())
		if !verifyAlready(item, list) && verifyIn(film.name, item) {
			c <- item
			sendToSlack = true
			film.show = "0"
		}
	})
	if sendToSlack {
		err := slkclient.SendToSlack(":cinema: Фильм", "Появился фильм: \"" + film.name + "\"", "", "", "")
		if err != nil {
			return
		}
	}
}

func main() {
	newFilms := make(chan string)
	films, err := loadFilmNamesFromFile(filmFileName)
	if err != nil {
		panic(err)
	}
	listAlready := loadAlreadyLoadFilms(passedFileName)
	index := 0
	client := &http.Client{Timeout: time.Second * 5}
	for i := range films {
		if films[i].show != "0" {
			go parseSite(client, &films[i], listAlready, newFilms)
			index++
		}
	}
	if index > 0 {
		wg.Add(1)
		go addAlreadyFilms(index, listAlready, newFilms)
		wg.Wait()
	}
	saveFilmNamesToFile(films, filmFileName)
}