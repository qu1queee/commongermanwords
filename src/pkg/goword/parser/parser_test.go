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
	}
	expected := []*models.Word{
		{
			Type: []string{"Interjektion", "Grußformel", "Abkürzung", "Interjektion"},
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
