package lexer

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/alediaferia/stackgo"
	"github.com/pkg/errors"
	"github.com/runcom/el/token"
)

var (
	numbersRegexp = regexp.MustCompile(`\A([0-9]+(?:\.[0-9]+)?)`)
	stringsRegexp = regexp.MustCompile(`\A("([^"\\\\]*(?:\\\\.[^"\\\\]*)*)"|\A'([^'\\\\]*(?:\\\\.[^'\\\\]*)*)')`)
	// FIXME: golang doesn't support Perl's (?=) see https://github.com/google/re2/wiki/Syntax
	//operatorsRegexp = regexp.MustCompile(`\Anot in(?=[\s(])|\!\=\=|not(?=[\s(])|and(?=[\s(])|\=\=\=|\>\=|or(?=[\s(])|\<\=|\*\*|\.\.|in(?=[\s(])|&&|\|\||matches|\=\=|\!\=|\*|~|%|\/|\>|\||\!|\^|&|\+|\<|\-`)
	operatorsRegexp = regexp.MustCompile(`\A(\!\=|\=\=|\>\=|\<\=|&&|\|\||\*|\/|\>|\||\!|\+|\<|\-)`)
	namesRegexp     = regexp.MustCompile(`\A([a-zA-Z_\x7f-\xff][a-zA-Z0-9_\x7f-\xff]*)`)
)

func Tokenize(expression string) (*token.TokenStream, error) {
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
				return nil, errors.Wrap(fmt.Errorf("internal error"), "lexer")
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
		} else if expression[cursor] == '.' || expression[cursor] == ',' || expression[cursor] == '?' || expression[cursor] == ':' {
			t := token.Token{
				Value:  expression[cursor],
				Type:   token.TypePunctuation,
				Cursor: cursor + 1,
			}
			tokens = append(tokens, t)
			cursor++
		} else if m := namesRegexp.FindStringIndex(expression[cursor:]); len(m) != 0 {
			t := token.Token{
				Value:  expression[cursor+m[0] : cursor+m[1]],
				Type:   token.TypeName,
				Cursor: cursor + 1,
			}
			tokens = append(tokens, t)
			cursor = cursor + (m[1] - m[0])
		} else {
			return nil, fmt.Errorf("unlexable %c, %d", expression[cursor], cursor)
		}
	}

	t := token.Token{
		Type:   token.TypeEOF,
		Cursor: cursor + 1,
	}
	tokens = append(tokens, t)

	if brackets.Size() > 0 {
		b := brackets.Pop()
		br, ok := b.(bracket)
		if !ok {
			return nil, errors.Wrap(fmt.Errorf("internal error"), "lexer")
		}
		return nil, fmt.Errorf("unexpected %c, %d", br.char, br.cursor)
	}
	return token.NewTokenStream(tokens), nil
}
