package types

type TokenType string

type TokenInfo struct {
	desc string
	key  []rune
	size int
}

const (
	EOL     TokenType = "END_OF_LINE"
	UNKNOWN TokenType = "UNKNOWN" // ignore first value
	AND     TokenType = "AND"
	OR      TokenType = "OR"
	ANDNOT  TokenType = "ANDNOT"
	ORNOT   TokenType = "ORNOT"
	NOT     TokenType = "NOT"
	EXP     TokenType = "EXPRESSION"
	LBR     TokenType = "LEFT_BRACKET"
	RBR     TokenType = "RIGHT_BRACKET"
	TRUE    TokenType = "TRUE"
)

func IsLeftAssociative(tok TokenType) bool {
	return tok == AND || tok == OR
}

func IsOp(tok TokenType) bool {
	return tok == AND || tok == OR || IsComplexOp(tok)
}

func IsComplexOp(tok TokenType) bool {
	return tok == ANDNOT || tok == ORNOT
}

func CouldBeComplexOp(tok TokenType) bool {
	return tok == AND || tok == OR
}

func IsWhitespace(r []rune, index int) bool {
	return (len(r)-index) > 0 && (r[index] == ' ' || r[index] == '\n' || r[index] == '\t')
}
