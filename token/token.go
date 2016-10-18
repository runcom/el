package token

type Type string

const (
	TypeEOF         TokenType = "end of expression"
	TypeName        TokenType = "name"
	TypeNumber      TokenType = "number"
	TypeString      TokenType = "string"
	TypeOperator    TokenType = "operator"
	TypePunctuation TokenType = "punctuation"
)

type Token struct {
	Value  interface{}
	Type   Type
	Cursor int
}
