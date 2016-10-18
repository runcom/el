package lexer

import (
	"fmt"
	"regexp"
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
	return nil
}

func (ts *TokenStream) EOF() bool {
	return ts.Current.Type == token.TypeEOF
}

func NewTokenStream(tokens []token.Token) *TokenStream {
	// TODO: check len(tokens) > 0
	return &TokenStream{
		Tokens:   tokens,
		Current:  tokens[0],
		Position: 0,
	}
}

var (
	numbersRegexp = regexp.MustCompile(`\A[0-9]+(?:\.[0-9]+)?`)
	stringsRegexp = regexp.MustCompile(`\A("([^"\\\\]*(?:\\\\.[^"\\\\]*)*)"|\A'([^'\\\\]*(?:\\\\.[^'\\\\]*)*)')`)
	// FIXME: golang doesn't support Perl's (?=) see https://github.com/google/re2/wiki/Syntax
	//operatorsRegexp = regexp.MustCompile(`\Anot in(?=[\s(])|\!\=\=|not(?=[\s(])|and(?=[\s(])|\=\=\=|\>\=|or(?=[\s(])|\<\=|\*\*|\.\.|in(?=[\s(])|&&|\|\||matches|\=\=|\!\=|\*|~|%|\/|\>|\||\!|\^|&|\+|\<|\-`)
	operatorsRegexp = regexp.MustCompile(`\A\!\=|\=\=|\>\=|\<\=|&&|\|\||\*|\/|\>|\||\!|\+|\<|\-`)
)

func Tokenize(expression string) (*TokenStream, error) {
	expression = strings.Replace(expression, "\r\t\v\f\n", " ", -1)
	var (
		cursor   int
		tokens   = []token.Token{}
		end      = len(expression)
		brackets = stackgo.NewStack()
	)
	type bracket struct {
		char   byte
		cursor int
	}

	// TODO(runcom): see https://github.com/symfony/expression-language/blob/master/Lexer.php

	for cursor < end {
		if expression[cursor] == ' ' {
			cursor++
			continue
		}
		if m := numbersRegexp.FindStringIndex(expression[cursor:]); len(m) != 0 {
			t := token.Token{
				Value:  expression[cursor+m[0] : cursor+m[1]],
				Type:   token.TypeNumber,
				Cursor: cursor + 1,
			}
			tokens = append(tokens, t)
			cursor = cursor + (m[1] - m[0])
		} else if expression[cursor] == '(' || expression[cursor] == '[' || expression[cursor] == '{' {
			brackets.Push(bracket{char: expression[cursor], cursor: cursor})
			t := token.Token{
				Value:  expression[cursor],
				Type:   token.TypePunctuation,
				Cursor: cursor + 1,
			}
			tokens = append(tokens, t)
			cursor++
		} else if expression[cursor] == ')' || expression[cursor] == ']' || expression[cursor] == '}' {
			if brackets.Size() == 0 {
				return nil, fmt.Errorf("unexpected %c, %d", expression[cursor], cursor)
			}
			b := brackets.Pop()
			br, ok := b.(bracket)
			if !ok {
				// return nil, err
			}
			var closingBracket byte
			switch br.char {
			case '(':
				closingBracket = ')'
			case '[':
				closingBracket = ']'
			case '{':
				closingBracket = '}'
			}
			if expression[cursor] != closingBracket {
				return nil, fmt.Errorf("unclosed %c, %d", br.char, br.cursor)
			}
			t := token.Token{
				Value:  expression[cursor],
				Type:   token.TypePunctuation,
				Cursor: cursor + 1,
			}
			tokens = append(tokens, t)
			cursor++
		} else if m := stringsRegexp.FindStringIndex(expression[cursor:]); len(m) != 0 {
			quotedStr := expression[cursor+m[0] : cursor+m[1]]
			t := token.Token{
				Value:  quotedStr[1 : len(quotedStr)-1],
				Type:   token.TypeString,
				Cursor: cursor + 1,
			}
			tokens = append(tokens, t)
			cursor = cursor + (m[1] - m[0])
		} else if m := operatorsRegexp.FindStringIndex(expression[cursor:]); len(m) != 0 {
			t := token.Token{
				Value:  expression[cursor+m[0] : cursor+m[1]],
				Type:   token.TypeOperator,
				Cursor: cursor + 1,
			}
			tokens = append(tokens, t)
			cursor = cursor + (m[1] - m[0])
		} else { // return unlexable!!!
			cursor++
			//continue
			//return nil, nil
		}
	}
	if brackets.Size() > 0 {
		return nil, fmt.Errorf("wtf")
	}
	return NewTokenStream(tokens), nil
}
