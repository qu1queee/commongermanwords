package main

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

func main() {

	os.Exit(0)

}

func getWordDataFromWiktionary(word string) ([]byte, error) {

	urlHead := "https://en.wiktionary.org/w/api.php?action=parse&page="
	urlTail := "&prop=wikitext&format=json"

	url := urlHead + url.QueryEscape(word) + urlTail

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
