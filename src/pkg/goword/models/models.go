package models

// Word holds the metadata of interest for this
// binary. Word can be extended, as long as the
// desire data is available in the wiktionary json
// for the related word
type Word struct {
	Type        []string            `json:"type"`
	IPA         []string            `json:"ipa"`
	Meaning     map[string][]string `json:"meaning"`
	Examples    map[string][]string `json:"examples"`
	Translation []string            `json:"translation"`
	Features    map[string][]string `json:"features"`
}

// Block holds a list of lines with content
// related to a particular section, these can
// be a meaning, translation, type, IPA(pronunciation),
// etc.
type Block struct {
	Lines []string `json:"lines,omitempty"`
	Title string   `json:"title,omitempty"`
}

type WordTypes struct {
	WordType map[string][]*Block `json:"wordtype,omitempty"`
}

// Article holds a word content from wiktionary where
// multiple languages can be defined for the word, and
// each language holds a list of content Blocks(e.g. translation)
type Article struct {
	Language map[string]WordTypes `json:"language,omitempty"`
}

// Audio holds and object containing information regarding
// the pronunciation audio of a word. We use a common json
// struct from wiktionary and auto-generate this object
type Audio struct {
	Query struct {
		Normalized []struct {
			From string `json:"from"`
			To   string `json:"to"`
		} `json:"normalized"`
		Pages struct {
			W struct {
				Ns              int    `json:"ns"`
				Title           string `json:"title"`
				Missing         string `json:"missing"`
				Known           string `json:"known"`
				Imagerepository string `json:"imagerepository"`
				Imageinfo       []struct {
					URL                 string `json:"url"`
					Descriptionurl      string `json:"descriptionurl"`
					Descriptionshorturl string `json:"descriptionshorturl"`
				} `json:"imageinfo"`
			} `json:"-1"`
		} `json:"pages"`
	} `json:"query"`
}
