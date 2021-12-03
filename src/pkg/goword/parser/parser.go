package parser

import (
	"bufio"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"strings"

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

func GetWord(data []byte, word string) (*models.Word, error) {
	re := regexp.MustCompile(`\"wikitext\":\{\"\*\":\"(.*?)\"\}\}\}$`)
	match := re.FindStringSubmatch(string(data))
	if len(match) == 0 {
		msg := fmt.Sprintf("No wikitext for word '%s'", word)
		return nil, errors.New(msg)
	}

	convertedText, err := strconv.Unquote(`"` + match[1] + `"`)
	if err != nil {
		return nil, nil
	}

	// todo: add description
	scanner := bufio.NewScanner(strings.NewReader(convertedText))
	wiktionaryArticle := models.Article{}

	// https://regex101.com/r/br7rzZ/1
	re = regexp.MustCompile(`(^{\{)(.*)(\}\})|(^\s{\{)(.*)(\}\})`)
	var auxBlock *models.Block

	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if strings.Contains(line, "=") {
			line = strings.Replace(line, "=", "", -1)
		}
		if matches := re.FindStringSubmatch(line); len(matches) > 0 {
			auxBlock = &models.Block{}
			if matches[2] == "" {
				auxBlock.Title = matches[5]
			} else {
				auxBlock.Title = matches[2]
			}
			wiktionaryArticle.Blocks = append(wiktionaryArticle.Blocks, auxBlock)
		} else {
			if auxBlock == nil {
				auxBlock = &models.Block{}
			}
			auxBlock.Lines = append(auxBlock.Lines, line)
		}

	}

	return GetSections(wiktionaryArticle)
}

func GetSections(article models.Article) (*models.Word, error) {
	wordObject := &models.Word{}
	for _, block := range article.Blocks {
		switch {
		case strings.Contains(block.Title, "Siehe auch|"):
			GetWordSection(block.Lines, wordObject)
		case block.Title == "Aussprache":
			GetIPASection(block.Lines, wordObject)
		case block.Title == "Bedeutungen":
			GetMeaningSection(block.Lines, wordObject)
		case block.Title == "Beispiele":
			GetUsageSection(block.Lines, wordObject)
		case block.Title == "Übersetzungen":
			GetTranslations(block.Lines, wordObject)
		default:
			//fmt.Println("todo: nothing to do")
		}
	}
	return wordObject, nil
}

func GetWordSection(lines []string, wordObject *models.Word) {
	re := regexp.MustCompile(`Wortart\|([a-zA-Zßäüö]{1,})\|Deutsch`)
	for _, line := range lines {
		if allMatches := re.FindAllStringSubmatch(line, -1); len(allMatches) > 0 {
			for _, match := range allMatches {
				wordObject.Type = append(wordObject.Type, match[1])
			}
		}
	}
}

func GetIPASection(lines []string, wordObject *models.Word) {
	re := regexp.MustCompile(`(Lautschrift\|(.*?)\}\})`)
	for _, line := range lines {
		if allMatches := re.FindAllStringSubmatch(line, -1); len(allMatches) > 0 {
			for _, match := range allMatches {
				wordObject.IPA = append(wordObject.IPA, match[2])
			}
		}
	}
}

func GetMeaningSection(lines []string, wordObject *models.Word) {
	for _, line := range lines {
		if line != "" {
			wordObject.Meaning = append(wordObject.Meaning, replaceIndexWithBrackets(line))
		}
	}
}

func GetUsageSection(lines []string, wordObject *models.Word) {
	for _, line := range lines {
		if line != "" {
			wordObject.Examples = append(wordObject.Examples, replaceIndexWithBrackets(line))
		}
	}
}

func GetTranslations(lines []string, wordObject *models.Word) {
	esRgx := regexp.MustCompile(`(es\|)([a-z]{1,})(}})`)
	enRgx := regexp.MustCompile(`(en\|)([a-z]{1,})(}})`)
	spanishTranslations := []string{}
	englishTranslations := []string{}
	for _, line := range lines {
		if line != "" {
			if strings.Contains(line, "{{es}}") {
				if matches := esRgx.FindAllStringSubmatch(line, -1); len(matches) > 0 {
					for _, translation := range matches {
						spanishTranslations = append(spanishTranslations, translation[2])
					}
					wordObject.Translation = append(wordObject.Translation, fmt.Sprintf("es: %v", strings.Join(spanishTranslations, ", ")))
				}

			}
			if strings.Contains(line, "{{en}}") {
				if matches := enRgx.FindAllStringSubmatch(line, -1); len(matches) > 0 {
					for _, translation := range matches {
						englishTranslations = append(englishTranslations, translation[2])
					}
					wordObject.Translation = append(wordObject.Translation, fmt.Sprintf("en: %v", strings.Join(englishTranslations, ", ")))
				}
			}
		}
	}
}

func replaceIndexWithBrackets(line string) string {
	m1 := regexp.MustCompile(`^:\[.*\] `)
	return m1.ReplaceAllString(line, "")
}
