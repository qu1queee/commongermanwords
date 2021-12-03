package parser

import (
	"reflect"
	"testing"

	"github.com/qu1queee/1000germanwords/src/pkg/goword/models"
)

func TestWordParser(t *testing.T) {
	testCases := []models.Article{
		{
			Blocks: []*models.Block{
				{
					Title: "Siehe auch|[[Hallo]]",
					Lines: []string{
						"== hallo ({{Sprache|Deutsch}}) ==",
						"=== {{Wortart|Interjektion|Deutsch}}, {{Wortart|Grußformel|Deutsch}} ===",
						"{{Wortart|Abkürzung|Deutsch}}..{{Wortart|Interjektion|Deutsch}}",
					},
				},
			},
		},
		{
			Blocks: []*models.Block{
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
		{
			Blocks: []*models.Block{
				{
					Title: "Bedeutungen",
					Lines: []string{
						":[1] ''als Interjektion:'' ein [[Anruf]], mit dem man andere, auch Fremde, auf sich aufmerksam machen will",
					},
				},
			},
		},
		{
			Blocks: []*models.Block{
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
		{
			Blocks: []*models.Block{
				{
					Title: "Übersetzungen",
					Lines: []string{
						"*{{en}}: [2] {{Ü|en|hi}}, ''britisch:'' {{Ü|en|hello}}",
						"*{{es}}: [1] {{Ü|es|oiga}}; [2] {{Ü|es|hola}}; [3] {{Ü|es|diga}}",
						"*{{tr}}: [2] {{Ü|tr|merhaba}}",
					},
				},
			},
		},
	}
	expected := []*models.Word{
		{
			Type: []string{"Interjektion", "Grußformel", "Abkürzung", "Interjektion"},
		},
		{
			IPA: []string{"ˈhalo", "haˈloː"},
		},
		{
			Meaning: []string{
				"''als Interjektion:'' ein [[Anruf]], mit dem man andere, auch Fremde, auf sich aufmerksam machen will",
			},
		},
		{
			Examples: []string{
				"''Hallo,'' Max!",
				"''Hallo'', Jana, bist du noch dran?",
				"sample test: : ",
			},
		},
		{
			Translation: []string{
				"en: hi, hello",
				"es: oiga, hola, diga",
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
