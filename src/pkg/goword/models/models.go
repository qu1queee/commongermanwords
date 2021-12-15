package models

// Word holds the metadata of interest for this
// binary. Word can be extended, as long as the
// desire data is available in the wiktionary json
// for the related word
type Word struct {
	Type        []string `json:"type"`
	IPA         []string `json:"ipa"`
	Meaning     []string `json:"meaning"`
	Examples    []string `json:"examples"`
	Translation []string `json:"translation"`
	Features    []string `json:"features"`
}

// Block holds a list of lines with content
// related to a particular section, these can
// be a meaning, translation, type, IPA(pronunciation),
// etc.
type Block struct {
	Lines []string `json:"lines,omitempty"`
	Title string   `json:"title,omitempty"`
}

// Article holds a word content from wiktionary where
// multiple languages can be defined for the word, and
// each language holds a list of content Blocks(e.g. translation)
type Article struct {
	Language map[string][]*Block
}
