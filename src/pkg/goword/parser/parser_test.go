package parser

import (
	"reflect"
	"testing"

	"github.com/qu1queee/1000germanwords/src/pkg/goword/models"
)

func TestWordParser(t *testing.T) {
	testCases := []string{
		"{{Siehe auch|[[Hallo]], [[halló]]}}\n== hallo ({{Sprache|Deutsch}}) ==\n=== {{Wortart|Interjektion|Deutsch}}, {{Wortart|Grußformel|Deutsch}} ===\n\n{{Worttrennung}}\n:hal·lo\n\n{{Aussprache}}\n:{{IPA}} {{Lautschrift|ˈhalo}}, {{Lautschrift|haˈloː}}\n:{{Hörbeispiele}} {{Audio|De-hallo.ogg}}, {{Audio|CPIDL German - Hallo.ogg}}\n:{{Reime}} {{Reim|alo|Deutsch}}, {{Reim|oː|Deutsch}}\n\n{{Bedeutungen}}\n:[1] ''als Interjektion:'' ein [[Anruf]], mit dem man andere, auch Fremde, auf sich aufmerksam machen will\n:[2] als Grußwort mit unverbindlichem Charakter\n:[3] [[Grußwort]] am [[Telefon]]\n\n{{Herkunft}}\n:im 15. Jahrhundert von [[mittelhochdeutsch]] ''holā'', ursprünglich an den [[Fährmann]] gerichteter [[Imperativ]] zu ''[[holen]]''<ref>{{Ref-DWDS|hallo}}</ref>\n\n{{Sinnverwandte Wörter}}\n:[1] ''siehe:'' [[Verzeichnis:Deutsch/Grüßen/Begrüßungsformeln|Begrüßungsformeln]]\n\n{{Gegenwörter}}\n:[2] [[tschüss]], [[auf Wiedersehen]], [[ade]]\n\n{{Oberbegriffe}}\n:[2] [[Gruß]]\n\n{{Beispiele}}\n:[1] ''Hallo,'' bist du vollkommen bescheuert?\n:[1] ''Ha - llo,'' wo bist du? <small>(geschrien)</small>\n:[2, 3] ''Hallo,'' Max!\n:[3] ''Hallo'', Jana, bist du noch dran?\n\n{{Redewendungen}}\n:[1] [[aber hallo]]!\n\n{{Wortbildungen}}\n:[[Hallo]], [[hallöchen]]\n\n==== {{Übersetzungen}} ====\n{{Ü-Tabelle|Ü-links=\n*{{sq}}: [2] {{Ü|sq|tungjatjeta}}; [3] {{Ü|sq|alo}}\n*{{bs}}: [2] {{Ü|bs|cao}}\n*{{zh}}: [1] {{Üt|zh|喂|wèi}}; [2] {{Üt|zh|你好|nǐ hǎo}}; [3] {{Üt|zh|喂|wéi}}\n*{{en}}: [2] {{Ü|en|hi}}, ''britisch:'' {{Ü|en|hello}}\n*{{eo}}: [1] {{Ü|eo|hola}}, {{Ü|eo|he}}, {{Ü|eo|hej}}, {{Ü|eo|ha}}; [2] {{Ü|eo|saluton}}; [3] {{Ü|eo|halo}}\n*{{fi}}: [2] {{Ü|fi|terve}}, {{Ü|fi|hei}}, {{Ü|fi|moi}}\n*{{fr}}: [1, 3] {{Ü|fr|allô}}; [2] {{Ü|fr|salut}}\n*{{got}}: [2, 3] {{Üt|got|𐌷𐌰𐌹𐌻𐍃|hails}} ''(nach Genus und Numerus zu deklinieren)''\n*{{el}}: [2] {{Üt|el|γεια|geia}}\n*{{ia}}: [1] {{Ü|ia|hallo}}\n*{{ga}}: [1] {{Ü|ga|Dia duit}}\n*{{is}}: [2] {{Ü|is|halló}}\n*{{it}}: [2] {{Ü|it|ciao}}; [3] {{Ü|it|pronto}}\n*{{ja}}: [1] {{Üt|ja|今日は|konnichiwa}}, [2] <small>am Telefon:</small> {{Üt|ja|もしもし|moshi-moshi}}\n*{{ca}}: [1] {{Ü|ca|hola}}\n*{{hr}}: [2] {{Ü|hr|čao}}\n*{{ku}}:\n**{{kmr}}: [2] {{Ü|kmr|silav}}\n|Ü-rechts=\n*{{lo}}: [02] {{Ü|lo|ສະບາຍດີ}}\n*{{la}}: [2] ''zu einer Person:'' {{Ü|la|salve}}, ''zu mehreren'' {{Ü|la|salvete}}\n*{{lb}}: [1–3] {{Ü|lb|moien}}, {{Ü|lb|bonjour}}\n*{{nds}}: [1–3] {{Ü|nds|moin}}, {{Ü|nds|moin moin}}\n*{{no}}: [1, 2] {{Ü|no|hei}}; [1–3] {{Ü|no|hallo}}\n*{{oc}}: [2] {{Ü|oc|adieu}}\n*{{fa}}: [2] {{Üt|fa|سلام|salām}}\n*{{pl}}: [2] {{Ü|pl|cześć}}, {{Ü|pl|witaj}}, {{Ü|pl|witajcie}}, {{Ü|pl|czołem}}; [3] {{Ü|pl|halo}}, {{Ü|pl|słucham}}\n*{{pt}}: [1, 2] {{Ü|pt|olá}}\n*{{ro}}: [1] {{Ü|ro|alo}}\n*{{ru}}: [1] {{Üt|ru|эй}}; [2] {{Üt|ru|привет}}; [3] {{Üt|ru|алло}}\n*{{sv}}: [1–3] {{Ü|sv|hallå}}, {{Ü|sv|hejsan}}; [2] {{Ü|sv|hej}}\n*{{sr}}: [2] {{Üt|sr|здраво|zdravo}}\n*{{sl}}: [2] {{Ü|sl|zdravo}}\n*{{wen}}:\n**{{dsb}}: [2] {{Ü|dsb|witaj}}, {{Ü|dsb|witajtej}}, {{Ü|dsb|witajśo}}, {{Ü|dsb|halo}}\n*{{es}}: [1] {{Ü|es|oiga}}; [2] {{Ü|es|hola}}; [3] {{Ü|es|diga}}, {{Ü|es|dígame}}\n*{{sw}}: [1, 2] {{Ü|sw|jambo}}\n*{{cs}}: [2] {{Ü|cs|ahoj}}, {{Ü|cs|nazdar}}; [1, 3] {{Ü|cs|haló}}\n*{{tr}}: [2] {{Ü|tr|merhaba}}\n*{{uk}}: [1] {{Üt|uk|ей|ej}}; [2] {{Üt|uk|привіт|pryvit}}; [3] {{Üt|uk|алло|allo}}\n*{{hu}}: [2] {{Ü|hu|szia}}\n|Dialekttabelle=\n*Alemannisch: [?] [[solli]]\n*Allgäuerisch: [?] [[griaßdi]]\n*Bayrisch: [?] [[servus]]\n*Burgenländerisch: [?] [[dere]]\n*Liechtensteinisch: [?] [[hoi]], [[sali]], [[servus]], [[salü]]\n*Norddeutsch: [1, 2] [[moin]] oder [[moin moin]]\n|D-rechts=\n*Rheinhessisch: [?] [[guden]], [[ei gude wie]]\n*Saarländisch: [?] [[Salü]]\n*Schwäbisch: [?] [[hoi]]\n*Schweizerdeutsch: [?] [[hoi]], [[sali]]\n*Südtirolerisch: [?] [[hoila]] (Südtirol - Italien)\n*Trierisch: [?] [[moin]]\n}}\n\n{{Referenzen}}\n:[1, 2] {{Wikipedia|Hallo}}\n:[1–3] {{Ref-DWDS|hallo}}\n:[1] {{Ref-UniLeipzig|hallo}}\n:[1–3] {{Ref-FreeDictionary|hallo}}\n\n{{Quellen}}\n\n{{Ähnlichkeiten 1|[[Halle]], [[Hallodri]], [[Hallore]], [[Halloween]]|Anagramme=[[holla]]}}\n\n{{Ähnlichkeiten 2|[[hall]], [[hillo]]}}",
	}
	expected := []*models.Word{
		&models.Word{
			Type: []string{"Interjektion", "Grußformel"},
		},
	}
	testFuncs := []func(data string) (*models.Word, error){
		GetWordSections,
	}

	for _, testFunc := range testFuncs {
		for index, data := range testCases {
			if res, _ := testFunc(data); !reflect.DeepEqual(res, expected[index]) {
				t.Errorf("got %v, expected %v", res, expected[index])
			}
		}
	}
}
