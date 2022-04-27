package types

type TokenType string

const (
	EOL     TokenType = "END_OF_LINE"
	UNKNOWN TokenType = "UNKNOWN" // ignore first value
	AND     TokenType = "AND"
	OR      TokenType = "OR"
	ANDNOT  TokenType = "ANDNOT"
	ORNOT   TokenType = "ORNOT"
	EXP     TokenType = "EXPRESSION"
	LBR     TokenType = "LEFT_BRACKET"
	RBR     TokenType = "RIGHT_BRACKET"
	TRUE    TokenType = "TRUE"
	FALSE   TokenType = "FALSE"
)

type Token struct {
	Typ TokenType
	Exp string
}

func (t *TokenType) ToString() string {
	switch *t {
	case UNKNOWN:
		return "unknown"
	case AND:
		return "and"
	case OR:
		return "or"
	case LBR:
		return "left bracket"
	case RBR:
		return "right bracket"
	case EXP:
		return "expression"
	case TRUE:
		return "true"
	case FALSE:
		return "false"
	case EOL:
		return "end of line"
	default:
		return "unknown"
	}
}

func IsLeftAssociative(tok TokenType) bool {
	return tok == AND || tok == OR
}

func IsOp(tok TokenType) bool {
	return tok == AND || tok == OR || IsComplexOp(tok)
}

func IsComplexOp(tok TokenType) bool {
	return tok == ANDNOT || tok == ORNOT
}
