package main

import (
	"log"
	"strings"
	"strconv"
	"io/ioutil"
	//"net/http"
)
const (
	filmFileName = "films.txt"
	passedFileName = "passed.txt"
)
type film struct {
	name string
	show int
}
func getFilmNames(filename string) ([]film, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	result := make([]film, 0)
	var tempFilm []string
	films := strings.Split(string(content), "\r\n")
	
	for _, name := range films {
		tempFilm = strings.Split(name, ";")
		show, err := strconv.Atoi(tempFilm[1])
		if err != nil {
			show = 1
		}
		result = append(result, film {tempFilm[0], show})
	}
	return result, nil
}

// func getSiteBody(url, proxy string) ([]byte, error) {
// 	//
// }
// func parseSite(body []byte, film string) bool {
// 	//
//}

func main() {
	films, _ := getFilmNames("films.txt")
	for _, kino := range films {
		log.Printf("Фильм: %s, показывать %b\n", kino.name, bool(kino.show))
	}
}