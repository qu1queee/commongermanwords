package main

import (
	"os"

	"github.com/davecgh/go-spew/spew"
	"github.com/qu1queee/1000germanwords/src/pkg/goword/parser"
)

func main() {
	if bytes, err := parser.GetCard("auf"); err == nil {
		if wordObject, err := parser.GetWord(bytes, "auf"); err == nil {
			spew.Dump(wordObject)
		}
	}
	os.Exit(0)

}
