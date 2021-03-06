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
	var currentType string

	scanner, err := convertJSONtoLines(data, word)
	if err != nil {
		return nil, fmt.Errorf("an error ocurred when retrieving from wiktionary: %v", err)
	}

	// compile required regexp for retrieving language and block types:
	// Wiktionary seems to be making the following distinction
	// double equals sign for languages, see: == hallo ({{Sprache|Deutsch}}) ==
	// triple equals sign for categories(blocks) under the same language,
	// see: === {{Wortart|Interjektion|Deutsch}}, {{Wortart|Grußformel|Deutsch}} ===
	reLanguageType := regexp.MustCompile(LANGUAGETYPEPARSER)
	reBlockType := regexp.MustCompile(BLOCKTYPEPARSER)
	reWordType := regexp.MustCompile(WORDTYPEREGEX)
	blocktype := &models.WordTypes{}
	blocktype.WordType = make(map[string][]*models.Block)

	// init Article with the related Language map
	wiktionaryArticle := models.Article{}
	wiktionaryArticle.Language = make(map[string]models.WordTypes)

	// read all content lines and categorize by language with multiple blocks.
	for i := 0; scanner.Scan(); i++ {
		line := scanner.Text()
		if matches := reLanguageType.FindAllStringSubmatch(line, -1); len(matches) > 0 {
			for _, entries := range matches {
				for id, entry := range entries {
					if id == 2 {
						languageName = entry
					}
				}
			}
		}

		// just process contents of Deutsch language
		if languageName == "" || languageName != "Deutsch" {
			continue
		}

		// init map with blocktype content
		if _, ok := wiktionaryArticle.Language[languageName]; ok {
			wiktionaryArticle.Language[languageName] = *blocktype
		} else {
			wiktionaryArticle.Language[languageName] = models.WordTypes{}
		}

		// sanitize per line
		if strings.Contains(line, "=") {
			line = strings.Replace(line, "=", "", -1)
		}

		// when processing new types, init blocktype map with wordtype and a list
		// of blocks
		if matches := reWordType.FindStringSubmatch(line); len(matches) > 0 {
			if matches[2] != "" {
				currentType = matches[2]
				blocktype.WordType[currentType] = []*models.Block{}
			}

		}
		// bail out if wordtype was not found
		if currentType == "" {
			continue
		}
		if matches := reBlockType.FindStringSubmatch(line); len(matches) > 0 {
			auxBlock = &models.Block{}
			if matches[2] == "" {
				auxBlock.Title = matches[5]
			} else {
				auxBlock.Title = matches[2]
			}
			// for current word type, append new blocks
			blocktype.WordType[currentType] = append(blocktype.WordType[currentType], auxBlock)
		} else {
			if auxBlock == nil {
				auxBlock = &models.Block{}
			}
			// for current block, append new lines
			auxBlock.Lines = append(auxBlock.Lines, line)
		}
	}

	return GetSections(wiktionaryArticle)
}

func GetSections(article models.Article) (*models.Word, error) {
	wordObject := &models.Word{}
	wordObject.Meaning = map[string][]string{}
	wordObject.Examples = map[string][]string{}
	wordObject.Features = map[string][]string{}
	for wordtype, s := range article.Language[DESIREDLANGUAGE].WordType {
		for _, block := range s {
			switch {
			case strings.Contains(block.Title, SIEHEAUCH) || strings.Contains(block.Title, WORDTYPE):
				GetWordType(block.Title, block.Lines, wordObject)
			case block.Title == PRONUNCIATION:
				GetIPASection(block.Lines, wordObject)
			case block.Title == MEANING:
				GetMeaningSection(block.Lines, wordObject, wordtype)
			case block.Title == EXAMPLES:
				GetUsageSection(block.Lines, wordObject, wordtype)
			case block.Title == TRANSLATION:
				GetTranslations(block.Lines, wordObject)
			case block.Title == FEATURES:
				GetFeatures(block.Lines, wordObject, wordtype)
			default:
				// todo: enhance default behaviour
			}
		}
	}
	return wordObject, nil
}

// GetWordType provides the grammatic type of the desired word
func GetWordType(title string, lines []string, wordObject *models.Word) {
	re := regexp.MustCompile(`Wortart\|([a-zA-Z\sßäüö]{1,})\|Deutsch`)

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
func GetMeaningSection(lines []string, wordObject *models.Word, wordtype string) {
	for _, line := range lines {
		if line != "" {
			wordObject.Meaning[wordtype] = append(wordObject.Meaning[wordtype], sanitizeLine(replaceIndexWithBrackets(line)))
		}
	}
}

// GetUsageSection consumes all examples of a word usage
func GetUsageSection(lines []string, wordObject *models.Word, wordtype string) {
	for _, line := range lines {
		if line != "" {
			wordObject.Examples[wordtype] = append(wordObject.Examples[wordtype], sanitizeLine(replaceIndexWithBrackets(line)))
		}
	}
}

// GetFeatures consumes all grammatical features of a word usage
// todo: missing unit test
func GetFeatures(lines []string, wordObject *models.Word, wordType string) {
	for _, line := range lines {
		if line != "" {
			wordObject.Features[wordType] = append(wordObject.Features[wordType], replaceAsterisk(sanitizeLine(replaceIndexWithBrackets(line))))

		}
	}
}

// GetTranslations provides a selected list of translations
// despite multiple available languages
func GetTranslations(lines []string, wordObject *models.Word) {
	esRgx := regexp.MustCompile(`(es\|)([a-z\s]{1,})(}})`)
	enRgx := regexp.MustCompile(`(en\|)([a-z\s]{1,})(}})`)

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
	if len(spanishTranslations) > 0 {
		wordObject.Translation = append(wordObject.Translation, fmt.Sprintf("es: %v", strings.Join(spanishTranslations, ", ")))
	}
	if len(englishTranslations) > 0 {
		wordObject.Translation = append(wordObject.Translation, fmt.Sprintf("en: %v", strings.Join(englishTranslations, ", ")))
	}
}

func replaceIndexWithBrackets(line string) string {
	if strings.HasPrefix(line, ":[") {
		line = strings.Replace(line, ":[", "[", 1)
	}
	return line
}

func replaceAsterisk(line string) string {
	if strings.HasPrefix(line, "*") {
		line = strings.Replace(line, "*", "", 1)
	}
	return line
}

// remove double marks and single quotes
func sanitizeLine(line string) string {
	return strings.ReplaceAll(strings.ReplaceAll(strings.ReplaceAll(line, "[[", ""), "]]", ""), `'`, "")

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
	if matches := regx.FindAllStringSubmatch(line, 3); len(matches) > 0 {
		for _, match := range matches {
			for index, translation := range match {
				if index == 2 {
					// limit translations per line to two
					if len(translations) < 2 {
						translations = append(translations, translation)
					}

				}
			}
		}
	}
	return translations
}
