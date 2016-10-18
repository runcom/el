package lexer_test

import (
	"fmt"
	"testing"

	"github.com/runcom/el/lexer"
)

func TestTokenize(t *testing.T) {
	ts, err := lexer.Tokenize("ciao 199d0 ciao483 899 88")
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(ts)
}
