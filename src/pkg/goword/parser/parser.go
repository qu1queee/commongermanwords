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

	"github.com/qu1queee/commongermanwords/src/pkg/goword/models"
)

// GetArticle provides an article based on a word,
// with all related languages blocks.
func GetArticle(word string) (*models.Word, error) {

	url := URLHEAD + url.QueryEscape(word) + URLTAIL

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return GetWordObject(body, word)
}

func convertJSONtoLines(data []byte, word string) (*bufio.Scanner, error) {

	wikiJSONParser := regexp.MustCompile(JSONREGEXPARSER)
	match := wikiJSONParser.FindStringSubmatch(string(data))
	if len(match) == 0 {
		msg := fmt.Sprintf("No wikitext for word '%s'", word)
		return nil, errors.New(msg)
	}

	text, err := strconv.Unquote(`"` + match[1] + `"`)
	if err != nil {
		return nil, nil
	}

	return bufio.NewScanner(strings.NewReader(text)), nil
}

func GetWordObject(data []byte, word string) (*models.Word, error) {

	var auxBlock *models.Block
	var languageName string

	scanner, err := convertJSONtoLines(data, word)
	if err != nil {
		return nil, fmt.Errorf("an error ocurred when retrieving from wiktionary: %v", err)
	}

	//
	// compile required regexp for retrieving language and block types:
	// Wiktionary seems to be making the following distinction
	// double equals sign for languages, see: == hallo ({{Sprache|Deutsch}}) ==
	// triple equals sign for categories(blocks) under the same language,
	// see: === {{Wortart|Interjektion|Deutsch}}, {{Wortart|Grußformel|Deutsch}} ===
	//
	reLanguageType := regexp.MustCompile(LANGUAGETYPEPARSER)
	reBlockType := regexp.MustCompile(BLOCKTYPEPARSER)

	// init Article with the related Language map
	wiktionaryArticle := models.Article{}
	wiktionaryArticle.Language = make(map[string][]*models.Block)

	// read all content lines and categorize by language with multiple blocks.
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if matches := reLanguageType.FindAllStringSubmatch(line, -1); len(matches) > 0 {
			for _, entries := range matches {
				for id, entry := range entries {
					if id == 2 {
						languageName = entry
						wiktionaryArticle.Language[languageName] = []*models.Block{}
					}
				}
			}
		}
		if languageName == "" {
			continue
		}

		if strings.Contains(line, "=") {
			line = strings.Replace(line, "=", "", -1)
		}
		if matches := reBlockType.FindStringSubmatch(line); len(matches) > 0 {
			auxBlock = &models.Block{}
			if matches[2] == "" {
				auxBlock.Title = matches[5]
			} else {
				auxBlock.Title = matches[2]
			}
			wiktionaryArticle.Language[languageName] = append(wiktionaryArticle.Language[languageName], auxBlock)
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
	for _, block := range article.Language[DESIREDLANGUAGE] {
		switch {
		case strings.Contains(block.Title, SIEHEAUCH) || strings.Contains(block.Title, WORDTYPE):
			GetWordType(block.Title, block.Lines, wordObject)
		case block.Title == PRONUNCIATION:
			GetIPASection(block.Lines, wordObject)
		case block.Title == MEANING:
			GetMeaningSection(block.Lines, wordObject)
		case block.Title == EXAMPLES:
			GetUsageSection(block.Lines, wordObject)
		case block.Title == TRANSLATION:
			GetTranslations(block.Lines, wordObject)
		case block.Title == FEATURES:
			GetFeatures(block.Lines, wordObject)
		default:
			// todo: enhance default behaviour
		}
	}
	return wordObject, nil
}

// GetWordType provides the grammatic type of the desired word
func GetWordType(title string, lines []string, wordObject *models.Word) {
	re := regexp.MustCompile(`Wortart\|([a-zA-Zßäüö]{1,})\|Deutsch`)

	if checkTitleMatches := re.FindAllStringSubmatch(title, -1); len(checkTitleMatches) > 0 {
		for _, match := range checkTitleMatches {
			wordObject.Type = append(wordObject.Type, match[1])
			return
		}
	}

	for _, line := range lines {
		if allMatches := re.FindAllStringSubmatch(line, -1); len(allMatches) > 0 {
			for _, match := range allMatches {
				wordObject.Type = append(wordObject.Type, match[1])
			}
		}
	}
}

// GetIPASection consumes all different pronunciations of a word
func GetIPASection(lines []string, wordObject *models.Word) {
	re := regexp.MustCompile(`(Lautschrift\|(.*?)\}\})`)
	for _, line := range lines {
		if allMatches := re.FindAllStringSubmatch(line, -1); len(allMatches) > 0 {
			for _, match := range allMatches {
				if !contains(wordObject.IPA, match[2]) {
					wordObject.IPA = append(wordObject.IPA, match[2])
				}
			}
		}
	}
}

// GetMeaningSection consumes all examples of a word definition in German
func GetMeaningSection(lines []string, wordObject *models.Word) {
	for _, line := range lines {
		if line != "" {
			wordObject.Meaning = append(wordObject.Meaning, replaceIndexWithBrackets(line))
		}
	}
}

// GetUsageSection consumes all examples of a word usage
func GetUsageSection(lines []string, wordObject *models.Word) {
	for _, line := range lines {
		if line != "" {
			wordObject.Examples = append(wordObject.Examples, replaceIndexWithBrackets(line))
		}
	}
}

// GetFeatures consumes all grammatical features of a word usage
// todo: missing unit test
func GetFeatures(lines []string, wordObject *models.Word) {
	for _, line := range lines {
		if line != "" {
			wordObject.Features = append(wordObject.Features, replaceIndexWithBrackets(line))
		}
	}
}

// GetTranslations provides a selected list of translations
// despite multiple available languages
func GetTranslations(lines []string, wordObject *models.Word) {
	esRgx := regexp.MustCompile(`(es\|)([a-z]{1,})(}})`)
	enRgx := regexp.MustCompile(`(en\|)([a-z]{1,})(}})`)

	var englishTranslations, spanishTranslations []string

	for _, line := range lines {
		if line != "" {
			if strings.Contains(line, "{{es}}") {
				spanishTranslations = append(spanishTranslations, getTranslationMatch(esRgx, line)...)
			}
			if strings.Contains(line, "{{en}}") {
				englishTranslations = append(englishTranslations, getTranslationMatch(enRgx, line)...)
			}
		}
	}
	wordObject.Translation = append(wordObject.Translation, fmt.Sprintf("es: %v", strings.Join(spanishTranslations, ", ")))
	wordObject.Translation = append(wordObject.Translation, fmt.Sprintf("en: %v", strings.Join(englishTranslations, ", ")))
}

func replaceIndexWithBrackets(line string) string {
	if strings.HasPrefix(line, ":[") {
		// https://regex101.com/r/JGwzZM/1
		m1 := regexp.MustCompile(`^:\[[\d\s,]{1,}\]\s`)
		return m1.ReplaceAllString(line, "")
	}
	return line
}

func contains(source []string, match string) bool {
	for _, v := range source {
		if v == match {
			return true
		}
	}

	return false
}

func getTranslationMatch(regx *regexp.Regexp, line string) []string {
	var translations []string
	if matches := regx.FindStringSubmatch(line); len(matches) > 0 {
		for index, translation := range matches {
			if index == 2 {
				translations = append(translations, translation)
			}
		}
	}
	return translations
}
