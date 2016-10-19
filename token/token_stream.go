package token

type TokenStream struct {
	// TODO: consider making some of these fields private (?)
	Tokens   []Token
	Current  Token
	Position int
}

func (ts *TokenStream) Next() error {
	// check that ts.Tokens[ts.Position] exists
	ts.Position++
	ts.Current = ts.Tokens[ts.Position]
	return nil
}

func (ts *TokenStream) Size() int {
	return len(ts.Tokens)
}

func (ts *TokenStream) EOF() bool {
	return ts.Current.Type == TypeEOF
}

func NewTokenStream(tokens []Token) *TokenStream {
	// TODO: check len(tokens) > 0
	return &TokenStream{
		Tokens:   tokens,
		Current:  tokens[0],
		Position: 0,
	}
}
