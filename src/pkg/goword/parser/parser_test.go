package parser

import (
	"io/ioutil"
	"os"
	"reflect"
	"strings"
	"testing"

	"github.com/qu1queee/commongermanwords/src/pkg/goword/models"
)

func TestWordParser(t *testing.T) {
	testCases := []models.Article{
		{
			Language: map[string][]*models.Block{
				"Deutsch": {
					{
						Title: "Siehe auch|[[Hallo]]",
						Lines: []string{
							"== hallo ({{Sprache|Deutsch}}) ==",
							"=== {{Wortart|Interjektion|Deutsch}}, {{Wortart|Grußformel|Deutsch}} ===",
							"{{Wortart|Abkürzung|Deutsch}}..{{Wortart|Interjektion|Deutsch}}",
							"{{Wortart|Konjugierte Form|Deutsch}}",
						},
					},
				},
				"Deutschi": {
					{
						Title: "Siehe auch|[[Hallo]]",
						Lines: []string{
							"=== {{Wortart|Fake|Deutsch}}, {{Wortart|Grußformel|Deutsch}} ===",
						},
					},
				},
			},
		},
		{
			Language: map[string][]*models.Block{
				"Deutsch": {
					{
						Title: "Aussprache",
						Lines: []string{
							":{{Reime}} {{Reim|alo|Deutsch}}, {{Reim|oː|Deutsch}}",
							":{{Hörbeispiele}} {{Audio|De-hallo.ogg}}, {{Audio|CPIDL German - Hallo.ogg}}",
							":{{IPA}} {{Lautschrift|ˈhalo}}, {{Lautschrift|haˈloː}}",
						},
					},
				},
			},
		},
		{
			Language: map[string][]*models.Block{
				"Deutsch": {
					{
						Title: "Aussprache",
						Lines: []string{
							":{{IPA}} {{Lautschrift|ˈhalo}}, {{Lautschrift|haˈloː}}, {{Lautschrift|haˈloː}}",
						},
					},
				},
			},
		},
		{
			Language: map[string][]*models.Block{
				"Deutsch": {
					{
						Title: "Bedeutungen",
						Lines: []string{
							":[1] ''als Interjektion:'' ein [[Anruf]], mit dem man andere, auch Fremde, auf sich aufmerksam machen will",
						},
					},
				},
			},
		},
		{
			Language: map[string][]*models.Block{
				"Deutsch": {
					{
						Title: "Beispiele",
						Lines: []string{
							":[2, 3] ''Hallo,'' Max!",
							":[3] ''Hallo'', Jana, bist du noch dran?",
							":[2] sample test: : ",
						},
					},
				},
			},
		},
		{
			Language: map[string][]*models.Block{
				"Deutsch": {
					{
						Title: "Übersetzungen",
						Lines: []string{
							"*{{en}}: [2] {{Ü|en|hi}}, ''britisch:'' {{Ü|en|hello}}",
							"*{{es}}: [1] {{Ü|es|oiga}}; [2] {{Ü|es|hola}}; [3] {{Ü|es|diga}}",
							"*{{tr}}: [2] {{Ü|tr|merhaba}}",
							"*{{es}}: [1, 2, 4] {{Ü|es|en}}, {{Ü|es|junto a}}; [3] ''Zeitpunkt:'' \"an\" wird nicht übersetzt, stattdessen wird der bestimmte Artikel verwendet (zum Beispiel: el {{Ü|es|domingo}}, am Sonntag), ''Zeitraum:'' {{Ü|es|por}} (zum Beispiel por {{Ü|es|navidad}}, zu Weihnachten)",
						},
					},
				},
			},
		},
		{
			Language: map[string][]*models.Block{
				"Deutsch": {
					{
						Title: "Grammatische Merkmale",
						Lines: []string{
							"*1. Person Singular Indikativ Präteritum Aktiv des Verbs '''[[sagen]]'''",
						},
					},
				},
			},
		},
	}

	expected := []*models.Word{
		{
			Type: []string{"Interjektion", "Grußformel", "Abkürzung", "Interjektion", "Konjugierte Form"},
		},
		{
			IPA: []string{"ˈhalo", "haˈloː"},
		},
		{
			IPA: []string{"ˈhalo", "haˈloː"},
		},
		{
			Meaning: []string{
				"[1] als Interjektion: ein Anruf, mit dem man andere, auch Fremde, auf sich aufmerksam machen will",
			},
		},
		{
			Examples: []string{
				"[2, 3] Hallo, Max!",
				"[3] Hallo, Jana, bist du noch dran?",
				"[2] sample test: : ",
			},
		},
		{
			Translation: []string{
				"es: oiga, en",
				"en: hi",
			},
		},
		{
			Features: []string{
				"1. Person Singular Indikativ Präteritum Aktiv des Verbs sagen",
			},
		},
	}
	testFuncs := []func(article models.Article) (*models.Word, error){
		GetSections,
	}

	for _, testFunc := range testFuncs {
		for index, data := range testCases {
			if res, _ := testFunc(data); !reflect.DeepEqual(res, expected[index]) {
				t.Errorf("got %v, expected %v", res, expected[index])
			}
		}
	}
}

func TestKnownWords(t *testing.T) {
	fileBytes, err := ioutil.ReadFile("../../../../wordslist/first1000.yaml")

	if err != nil {
		os.Exit(1)
	}
	sliceData := strings.Split(string(fileBytes), "\n")

	for _, line := range sliceData[2:] {
		line = strings.TrimSpace(strings.Replace(line, "-", "", -1))
		if strings.Contains(line, "�") { //todo
			continue
		}
		t.Log(line)
		if model, _ := GetArticle(line); !isValidObject(model) {
			t.Errorf("got an incomplete word: %v", line)
		}

	}

}

func isValidObject(w *models.Word) bool {
	if w == nil {
		return false
	}
	if len(w.IPA) <= 0 {
		return false
	}
	return true
}
