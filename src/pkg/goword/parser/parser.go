package parser

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"

	"github.com/qu1queee/1000germanwords/src/pkg/goword/models"
)

func GetCard(word string) ([]byte, error) {

	urlHead := "https://de.wiktionary.org/w/api.php?action=parse&page="
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

func GetWord(data []byte, word string) (string, error) {
	re := regexp.MustCompile(`\"wikitext\":\{\"\*\":\"(.*?)\"\}\}\}$`)
	match := re.FindStringSubmatch(string(data))
	if len(match) == 0 {
		msg := fmt.Sprintf("No wikitext for word '%s'", word)
		return "", errors.New(msg)
	}

	convertedText, err := strconv.Unquote(`"` + match[1] + `"`)
	if err != nil {
		return "", nil
	}
	GetWordSections(convertedText)

	return "wikitext", nil
}

func GetWordSections(data string) (*models.Word, error) {
	wordContent := &models.Word{}
	re := regexp.MustCompile(`Wortart\|([a-zA-Zß]{1,})\|Deutsch`)
	allMatches := re.FindAllStringSubmatch(data, -1)
	if len(allMatches) == 0 {
		//TODO
	}
	for _, match := range allMatches {
		wordContent.Type = append(wordContent.Type, match[1])
	}
	return wordContent, nil
}
