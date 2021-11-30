package main

import (
	"os"

	"github.com/qu1queee/1000germanwords/src/pkg/goword/parser"
)

func main() {
	if bytes, err := parser.GetCard("hallo"); err == nil {
		parser.GetWord(bytes, "hallo")
	}
	os.Exit(0)

}
