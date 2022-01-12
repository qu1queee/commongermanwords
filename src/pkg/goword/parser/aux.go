package parser

const (
	//
	// WIKTIONARY API URL's
	URLHEAD = "https://de.wiktionary.org/w/api.php?action=parse&page="
	URLTAIL = "&prop=wikitext&format=json"
	//

	//
	// List of current REGEXP
	JSONREGEXPARSER = `\"wikitext\":\{\"\*\":\"(.*?)\"\}\}\}$`
	// https://regex101.com/r/G8C38F/1
	LANGUAGETYPEPARSER = `(^={2}\s.*Sprache\|)(.*)(\}\}\)\s={2}$)`
	// https://regex101.com/r/br7rzZ/1
	BLOCKTYPEPARSER = `(^{\{)(.*)(\}\})|(^\s{\{)(.*)(\}\})`
	WORDTYPEREGEX   = `(Wortart\|)([a-zA-Z\s]{1,})(\|)`
	//

	DESIREDLANGUAGE = "Deutsch"

	//
	//List of categories(blocks) of interest
	WORDTYPE      = "Wortart|"
	SIEHEAUCH     = "Siehe auch|"
	PRONUNCIATION = "Aussprache"
	MEANING       = "Bedeutungen"
	EXAMPLES      = "Beispiele"
	TRANSLATION   = "Übersetzungen"
	FEATURES      = "Grammatische Merkmale"
	//
)
