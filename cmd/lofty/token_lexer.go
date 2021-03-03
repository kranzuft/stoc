package lofty

import (
	"fmt"
	"math"
	"runtime/debug"
	"strings"
)

type stateFn func([]rune, int) (int, token, stateFn)

type status int

const (
	LEXING_ERROR status = iota
	SHUNTING_ERROR
	OK
)

func parseBooleanAlgebraToTokens(raw []rune) (bool, []token) {
	var tokens []token
	var token token

	index := skipWhitespace(raw, 0)
	for state := expressionOrLeftBracket; state != nil; {
		index, token, state = state(raw, index)
		if token.typ == UNKNOWN {
			return false, tokens
		}
		tokens = append(tokens, token)
	}

	fmt.Println(toString(tokens))

	return true, tokens
}

// Lexical left bracket
func lexLeftBracket(rawData []rune, index int) (int, token, stateFn) {
	var tok token
	tok.typ = LBR // we don't need to track that it's a left bracket anymore
	tok.exp = string(rawData[index])
	i := skipWhitespace(rawData, index+1)

	if i < len(rawData) {
		if isExpRune(rawData[i]) {
			return i, tok, lexExpression
		} else if isLeftBracket(rawData[i]) {
			return i, tok, lexLeftBracket
		}
	}

	return printError(rawData, i, EXP)
}

func lexRightBracket(rawData []rune, index int) (int, token, stateFn) {
	var tok token
	tok.typ = RBR // we don't need to track that it's a right bracket anymore
	tok.exp = string(rawData[index])
	i := skipWhitespace(rawData, index+1)

	if i < len(rawData) {
		if isAssociativeOp(rawData[i]) {
			return i, tok, lexOperator
		} else if isRightBracket(rawData[i]) {
			return i, tok, lexRightBracket
		}
	} else {
		return i, tok, nil
	}

	return printError(rawData, i, RBR)
}

// Lexes an expression
func lexExpression(rawData []rune, index int) (int, token, stateFn) {
	var tok token
	tok.typ = EXP
	nextToken := false
	i := index
	lastI := i
	trailingWhitespace := 0

	if string(rawData) == "(fences|fences)" {
		fmt.Println()
	}

	for ; i < len(rawData) && !nextToken; i++ {
		curr := rawData[i]
		if isExpRune(curr) {
			trailingWhitespace = 0
		} else if isWhitespace(curr) {
			trailingWhitespace++
		} else if i != index && (isRightBracket(curr) || isAssociativeOp(curr)) {
			nextToken = true
		} else {
			return printError(rawData, i, EXP)
		}
	}
	if i == len(rawData) && i > 0 && isExpRune(rawData[i-1]) {
		tok.exp = string(rawData[lastI : i-trailingWhitespace])
	} else {
		tok.exp = string(rawData[lastI : i-trailingWhitespace-1])
	}
	i--

	i = skipWhitespace(rawData, i)

	if i < len(rawData) {
		if isAssociativeOp(rawData[i]) {
			return i, tok, lexOperator
		} else if isRightBracket(rawData[i]) {
			return i, tok, lexRightBracket
		}
	}

	return i, tok, nil
}

func lexOperator(rawData []rune, index int) (int, token, stateFn) {
	var tok token

	//if index+1 < len(rawData) {
	//	tok.typ = determineTokenType(rawData[index], rawData[index+1])
	//	tok.exp = string(rawData[index])
	//	if isComplexOp(tok.typ) {
	//		tok.exp += string(rawData[index+1])
	//		index++
	//	}
	//} else
	if index < len(rawData) {
		tok.typ = determineTokenType(rawData[index], ' ')
		tok.exp = string(rawData[index])
	}
	i := skipWhitespace(rawData, index+1)

	if i < len(rawData) {
		if isExpRune(rawData[i]) {
			return i, tok, lexExpression
		} else if isLeftBracket(rawData[i]) {
			return i, tok, lexLeftBracket
		}
	}

	return printError(rawData, index+1, EXP, LBR)
}

func (t tokenType) in(toks []tokenType) bool {
	for _, tok := range toks {
		if tok == t {
			return true
		}
	}
	return false
}

func determineTokenType(character1 rune, character2 rune) tokenType {
	if isAndNot(character1, character2) {
		return ANDNOT
	} else if isOrNot(character1, character2) {
		return ORNOT
	} else if isAnd(character1) {
		return AND
	} else if isOr(character1) {
		return OR
	} else if isExpRune(character1) {
		return EXP
	} else if isRightBracket(character1) {
		return RBR
	} else if isLeftBracket(character1) {
		return LBR
	}

	return UNKNOWN
}

func ifStr(yes bool, str1 string, str2 string) string {
	if yes {
		return str1
	}
	return str2
}

func printError(rawData []rune, index int, expected ...tokenType) (int, token, stateFn) {
	var found tokenType
	index = skipWhitespace(rawData, index)
	if (index+1) < len(rawData) {
		found = determineTokenType(rawData[index], rawData[index+1])
	} else if index < len(rawData) {
		found = determineTokenType(rawData[index], ' ')
	} else {
		found = EOL
	}
	message := tokensToList(expected) + " expected, instead found " + found.toString()
	return errorMessage(rawData, index, message)
}

func errorMessage(rawData []rune, index int, message string) (int, token, stateFn) {
	fmt.Println("ERROR: " + message)
	lenIndexDiff := len(rawData) - index
	start := int(math.Min(20, float64(index)))
	end := int(math.Min(20, float64(lenIndexDiff)))
	fmt.Println(string(rawData[index-start : index+end]))
	fmt.Println(strings.Repeat(" ", start) + "^")
	debug.PrintStack()
	return index, getErrorToken(), nil
}

func skipWhitespace(rawData []rune, index int) int {
	i := index

	for ; i < len(rawData); i++ {
		if !isWhitespace(rawData[i]) {
			return i
		}
	}

	return i
}

func expressionOrLeftBracket(rawData []rune, index int) (int, token, stateFn) {
	if isLeftBracket(rawData[index]) {
		return lexLeftBracket(rawData, index)
	} else if isExpRune(rawData[index]) {
		return lexExpression(rawData, index)
	}

	return printError(rawData, index, LBR, EXP)
}

func isLeftBracket(r rune) bool {
	return r == '('
}

func isRightBracket(r rune) bool {
	return r == ')'
}

func isExpRune(r rune) bool {
	return !(isWhitespace(r) || isKeyword(r))
}

func isAnd(r rune) bool {
	return r == '&'
}

func isOr(r rune) bool {
	return r == '|'
}

func isAndNot(r1 rune, r2 rune) bool {
	return r1 == '&' && r2 == '!'
}

func isOrNot(r1 rune, r2 rune) bool {
	return r1 == '|' && r2 == '!'
}

func isNot(r rune) bool {
	return r == '!'
}

func isAssociativeOp(r rune) bool {
	return isAnd(r) || isOr(r)
}

func isKeyword(r rune) bool {
	return isOr(r) || isAnd(r) || isLeftBracket(r) || isRightBracket(r) || isNot(r)
}

func isWhitespace(r rune) bool {
	return r == ' ' || r == '\n' || r == '\t'
}

func getErrorToken() token {
	return token{UNKNOWN, ""}
}

func tokensToList(expected []tokenType) string {
	var res strings.Builder
	for index, exp := range expected {
		if len(expected) > 1 && index == len(expected)-1 {
			res.WriteString(" or ")
		}
		res.WriteString(exp.toString())
		if index < len(expected)-2 {
			res.WriteString(", ")
		}

	}
	return res.String()
}
