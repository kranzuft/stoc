package processing

import (
	"fmt"
	"lofty/cmd/lofty/types"
	"math"
	"runtime/debug"
	"strings"
)

type stateFn func([]rune, int) (int, types.Token, stateFn)

// BooleanAlgebraLexer looks for boolean operations in the search text, and transforms the operations into a token list.
//
// If an unexpected token type comes up, outside the expected order, false is returned. In this situation, a message
// shows where in the line the exception occurred, with a descriptive message included that states what
// was expected.
//
func BooleanAlgebraLexer(raw []rune) (bool, []types.Token) {
	var tokens []types.Token
	var token types.Token

	if len(raw) == 0 {
		return false, tokens
	}

	index := skipWhitespace(raw, 0)
	for state := determineIfExpressionOrLeftBracketOrNot; state != nil; {
		index, token, state = state(raw, index)
		if token.Typ == types.UNKNOWN {
			return false, tokens
		}
		tokens = append(tokens, token)
	}

	return true, tokens
}

// Lexical left bracket
func lexLeftBracket(rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token
	tok.Typ = types.LBR // we don't need to track that it's a left bracket anymore
	tok.Exp = string(rawData[index])
	i := skipWhitespace(rawData, index+1)

	if i < len(rawData) {
		if isExpRune(rawData[i]) {
			return i, tok, lexExpression
		} else if isLeftBracket(rawData[i]) {
			return i, tok, lexLeftBracket
		} else if isNot(rawData[i]) {
			return i, tok, lexNotSymbolAndCreateTrueToken
		}
	}

	return printError(rawData, i, types.EXP)
}

// Lexical right bracket
func lexRightBracket(rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token
	tok.Typ = types.RBR // we don't need to track that it's a right bracket anymore
	tok.Exp = string(rawData[index])
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

	return printError(rawData, i, types.RBR)
}

// Lexes an expression
func lexExpression(rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token
	tok.Typ = types.EXP
	nextToken := false
	i := index
	lastI := i
	trailingWhitespace := 0

	for ; i < len(rawData) && !nextToken; i++ {
		curr := rawData[i]
		if isExpRune(curr) {
			trailingWhitespace = 0
		} else if isWhitespace(curr) {
			trailingWhitespace++
		} else if i != index && (isRightBracket(curr) || isAssociativeOp(curr)) {
			nextToken = true
		} else {
			return printError(rawData, i, types.EXP)
		}
	}

	if i == len(rawData) && i > 0 && isExpRune(rawData[i-1]) {
		tok.Exp = string(rawData[lastI : i-trailingWhitespace])
	} else {
		tok.Exp = string(rawData[lastI : i-trailingWhitespace-1])
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

func getNextFuncAfterOperator(rawData []rune, index int, tok types.Token) (int, types.Token, stateFn) {
	if index < len(rawData) {
		if isExpRune(rawData[index]) {
			return index, tok, lexExpression
		} else if isLeftBracket(rawData[index]) {
			return index, tok, lexLeftBracket
		}
	}

	return printError(rawData, index+1, types.EXP, types.LBR)
}

func lexNotSymbolAndCreateTrueToken(_ []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token

	tok.Typ = types.TRUE
	tok.Exp = "true"

	return index, tok, lexNotSymbolAndCreateAndNotToken
}

func lexNotSymbolAndCreateAndNotToken(rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token

	tok.Typ = types.ANDNOT
	tok.Exp = "&!"

	i := skipWhitespace(rawData, index+1)

	return getNextFuncAfterOperator(rawData, i, tok)
}

func lexOperator(rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token

	i1 := index
	i2 := skipWhitespace(rawData, i1+1)
	if i2 < len(rawData) {
		tok.Typ = determineTokenType(rawData[i1], rawData[i2])
		tok.Exp = string(rawData[index])
		if types.IsComplexOp(tok.Typ) {
			tok.Exp += string(rawData[i2])
			index = i2
		}
	} else if i1 < len(rawData) {
		tok.Typ = determineTokenType(rawData[index], ' ')
		tok.Exp = string(rawData[index])
	}
	i := skipWhitespace(rawData, index+1)

	return getNextFuncAfterOperator(rawData, i, tok)
}

func determineTokenType(character1 rune, character2 rune) types.TokenType {
	if isAndNot(character1, character2) {
		return types.ANDNOT
	} else if isOrNot(character1, character2) {
		return types.ORNOT
	} else if isAnd(character1) {
		return types.AND
	} else if isOr(character1) {
		return types.OR
	} else if isExpRune(character1) {
		return types.EXP
	} else if isRightBracket(character1) {
		return types.RBR
	} else if isLeftBracket(character1) {
		return types.LBR
	}

	return types.UNKNOWN
}

func printError(rawData []rune, index int, expected ...types.TokenType) (int, types.Token, stateFn) {
	var found types.TokenType
	index = skipWhitespace(rawData, index)
	if (index + 1) < len(rawData) {
		found = determineTokenType(rawData[index], rawData[index+1])
	} else if index < len(rawData) {
		found = determineTokenType(rawData[index], ' ')
	} else {
		found = types.EOL
	}
	message := tokensToList(expected) + " expected, instead found " + found.ToString()
	return errorMessage(rawData, index, message)
}

func errorMessage(rawData []rune, index int, message string) (int, types.Token, stateFn) {
	fmt.Println("ERROR: " + message)
	if len(rawData) > 0 {
		lenIndexDiff := len(rawData) - index
		start := int(math.Min(20, float64(index)))
		end := int(math.Min(20, float64(lenIndexDiff)))
		fmt.Println(string(rawData[index-start : index+end]))
		fmt.Println(strings.Repeat(" ", start) + "^")
		debug.PrintStack()
	}
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

func determineIfExpressionOrLeftBracketOrNot(rawData []rune, index int) (int, types.Token, stateFn) {
	if isLeftBracket(rawData[index]) {
		return lexLeftBracket(rawData, index)
	} else if isExpRune(rawData[index]) {
		return lexExpression(rawData, index)
	} else if isNot(rawData[index]) {
		return lexNotSymbolAndCreateTrueToken(rawData, index)
	}

	return printError(rawData, index, types.LBR, types.EXP)
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

func getErrorToken() types.Token {
	return types.Token{Typ: types.UNKNOWN}
}

func tokensToList(expected []types.TokenType) string {
	var res strings.Builder
	for index, exp := range expected {
		if len(expected) > 1 && index == len(expected)-1 {
			res.WriteString(" or ")
		}
		res.WriteString(exp.ToString())
		if index < len(expected)-2 {
			res.WriteString(", ")
		}

	}
	return res.String()
}
