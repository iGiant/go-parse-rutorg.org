package main

import (
	// "log"
	"strings"
	"io/ioutil"
	//"net/http"
)
const (
	filmFileName = "films.txt"
	passedFileName = "passed.txt"
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
		result = append(result, film {tempFilm[0], "1"})
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

// func getSiteBody(url, proxy string) ([]byte, error) {
// 	//
// }
// func parseSite(body []byte, film string) bool {
// 	//
//}

func main() {
	films, _ := getFilmNamesFromFile(filmFileName)
	
	setFilmNamesToFile(films, filmFileName)
}