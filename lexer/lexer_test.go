package lexer_test

import (
	"testing"

	"github.com/runcom/el/lexer"
	"github.com/runcom/el/token"
)

func TestTokenize(t *testing.T) {
	ts, err := lexer.Tokenize("")
	if err != nil {
		t.Fatal(err)
	}
	if ts.Size() != 1 {
		t.Fatalf("expected only 1 element, got %d", ts.Size())
	}
	if ts.Current.Type != token.TypeEOF {
		t.Fatalf("expected token type EOF, got %v", ts.Current.Type)
	}
	// FIXME: test this out soon...
	//ts, err := lexer.Tokenize("ciao({} )(199d0) ciao483 '899' 88")
	ts, err = lexer.Tokenize(`({} )(1990) >= +ciao483 "888" 88+11?`)
	if err != nil {
		t.Fatal(err)
	}
	//fmt.Println(ts)
}
