package el

import (
	"strings"

	"github.com/alediaferia/stackgo"
	"github.com/runcom/el/token"
)

type TokenStream struct {
	Tokens   []token.Token
	Current  token.Token
	Position int
}

func (ts *TokenStream) Next() error {
	// check that ts.Tokens[ts.Position] exists
	ts.Position++
	ts.Current = ts.Tokens[ts.Position]
}

func (ts *TokenStream) EOF() bool {
	ts.Current.Type == token.TypeEOF
}

func NewTokenStream(tokens []token.Token) *TokenStream {
	// TODO: check len(tokens) > 0
	return &TokenStream{
		Tokens:   tokens,
		Current:  tokens[0],
		Position: 0,
	}
}

func tokenize(expression string) (*TokenStream, error) {
	expression = strings.Replace(expression, "\r\t\v\f\n", "", -1)
	var (
		cursor   int
		tokens   = []token.Token{}
		end      = len(expression)
		brackets = stackgo.NewStack()
	)
	type bracket struct {
		char   rune
		cursor int
	}
	for cursor < end {
		if expression[cursor] == ' ' {
			cursor++
			continue
		}
		if expression[cursor] == '(' || expression[cursor] == '[' || expression[cursor] == '{' {
			brackets.Push(bracket{char: expression[cursor], cursor: cursor})

			// token = new punctuation token ...

			cursor++
		} else if expression[cursor] == ')' || expression[cursor] == ']' || expression[cursor] == '}' {
			if brackets.Size() == 0 {
				// error out!!!! unexpected punctuation
				return nil, nil
			}
			b := brackets.Pop()
			br, ok := b.(bracket)
			if !ok {
				// error out!!!
				return nil, nil
			}
			//br.char
	}
	return nil, nil
}
