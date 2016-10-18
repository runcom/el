package token

type Type string

const (
	TypeEOF         Type = "end of expression"
	TypeName        Type = "name"
	TypeNumber      Type = "number"
	TypeString      Type = "string"
	TypeOperator    Type = "operator"
	TypePunctuation Type = "punctuation"
)

type Token struct {
	Value  interface{}
	Type   Type
	Cursor int
}
