package lexer_test

import (
	"fmt"
	"testing"

	"github.com/runcom/el/lexer"
)

func TestTokenize(t *testing.T) {
	//ts, err := lexer.Tokenize("ciao({} )(199d0) ciao483 '899' 88")
	ts, err := lexer.Tokenize(`({} )(1990) 483 "888" 88+11`)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ts)
}
