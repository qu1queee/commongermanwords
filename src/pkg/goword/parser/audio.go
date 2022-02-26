package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/qu1queee/commongermanwords/src/pkg/goword/models"
)

// GetLinkToAudio constructs an endpoint to a ogg file
func GetLinkToAudio(word *models.Word, name string) (string, error) {
	url := fmt.Sprintf("https://de.wiktionary.org/w/api.php?action=query&prop=imageinfo&iiprop=url&format=json&iwurl=l&rawcontinue=&titles=File:De-%v.ogg", name)

	res, err := httpRequest(url)
	if err != nil {
		return "", err
	}

	if res.Body != nil {
		defer res.Body.Close()
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return "", err
	}

	wordAudio, err := getAudioObject(body)
	if err != nil {
		return "", err
	}

	if len(wordAudio.Query.Pages.W.Imageinfo) > 0 {
		return wordAudio.Query.Pages.W.Imageinfo[0].Descriptionurl, nil
	}

	return "", errors.New("requested audio url not found")
}

func getAudioObject(content []byte) (*models.Audio, error) {
	wordAudio := &models.Audio{}
	err := json.Unmarshal(content, wordAudio)
	if err != nil {
		return nil, err
	}
	return wordAudio, nil
}

func httpRequest(url string) (*http.Response, error) {
	httpClient := http.Client{Timeout: time.Second * 2}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	res, err := httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	return res, nil
}
