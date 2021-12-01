package models

type Card struct {
	Title  string `json:"title"`
	PageID int    `json:"pageid"`
	Word   *Word  `json:"word"`
}

type Word struct {
	Type        []string `json:"type"`
	IPA         []string `json:"ipa"`
	Meaning     []string `json:"meaning"`
	Examples    []string `json:"examples"`
	Translation []string `json:"translation"`
}

type Block struct {
	Lines []string
	Title string
}

type Article struct {
	Blocks []*Block
}
