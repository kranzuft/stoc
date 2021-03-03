package lofty

type tokenType string

const (
	EOL     tokenType = "END_OF_LINE"
	UNKNOWN tokenType = "UNKNOWN" // ignore first value
	AND     tokenType = "AND"
	OR      tokenType = "OR"
	ANDNOT     tokenType = "ANDNOT"
	ORNOT      tokenType = "ORNOT"
	EXP     tokenType = "EXPRESSION"
	LBR     tokenType = "LEFT_BRACKET"
	RBR     tokenType = "RIGHT_BRACKET"
)

type token struct {
	typ tokenType
	exp string
}

func toString(toks []token) string {
	s := ""
	for _, r := range toks {
		s+= r.exp + " "
	}
	return s
}

func (t *tokenType) toString() string {
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
	case EOL:
		return "end of line"
	default:
		return "unknown"
	}
}

func isLeftAssociative(tok tokenType) bool {
	return tok == AND || tok == OR
}

func isAssociative(tok tokenType) bool {
	return tok == AND || tok == OR
}

func isOp(tok tokenType) bool {
	return tok == AND || tok == OR
}
func isComplexOp(tok tokenType) bool {
	return tok == ANDNOT || tok == ORNOT
}
