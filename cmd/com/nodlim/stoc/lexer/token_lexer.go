package lexer

import (
	"fmt"
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc/types"
	"math"
	"runtime/debug"
	"strings"
)

type Definitions interface {
	IsLeftBracket(r []rune, index int) bool

	IsLeftBracketI() int

	IsRightBracket(r []rune, index int) bool

	IsRightBracketI() int

	IsAnd(r []rune, index int) bool

	IsAndI() int

	IsOr(r []rune, index int) bool

	IsOrI() int

	IsNot(r []rune, index int) bool

	IsNotI() int

	IsExpRune(r []rune, index int) bool

	IsAssociativeOp(r []rune, index int) bool

	IsKeyword(r []rune, index int) bool

	IsQuote(r []rune, index int) bool

	IsSingleInvertedComma(r []rune, index int) bool

	IsDoubleInvertedComma(r []rune, index int) bool

	TokToString(t types.TokenType) string
}

type stateFn func(Definitions, []rune, int) (int, types.Token, stateFn)

// BooleanAlgebraLexer looks for boolean operations in the search text, and transforms the operations into a token list.
//
// If an unexpected token type comes up, outside the expected order, false defs.Is returned. In this situation, a message
// shows where in the line the exception occurred, with a descriptive message included that states what
// was expected.
//
func BooleanAlgebraLexer(defs Definitions, raw []rune) (bool, []types.Token) {
	var tokens []types.Token
	var token types.Token

	if len(raw) == 0 {
		return false, tokens
	}

	index := skipWhitespace(raw, 0)
	for state := determineIfExpressionOrLeftBracketOrNot; state != nil; {
		index, token, state = state(defs, raw, index)
		if token.Typ == types.UNKNOWN {
			return false, tokens
		}
		tokens = append(tokens, token)
	}

	return true, tokens
}

// Lexical left bracket
func lexLeftBracket(defs Definitions, rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token
	tok.Typ = types.LBR // we don't need to track that it's a left bracket anymore
	tok.Exp = string(rawData[index : index+defs.IsLeftBracketI()])
	i := skipWhitespace(rawData, index+1)

	if i < len(rawData) {
		if defs.IsExpRune(rawData, i) {
			return i, tok, lexExpression
		} else if defs.IsLeftBracket(rawData, i) {
			return i, tok, lexLeftBracket
		} else if defs.IsNot(rawData, i) {
			return i, tok, lexNotSymbolAndCreateTrueToken
		}
	}

	return printError(defs, rawData, i, types.EXP)
}

// Lexical right bracket
func lexRightBracket(defs Definitions, rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token
	tok.Typ = types.RBR // we don't need to track that it's a right bracket anymore
	tok.Exp = string(rawData[index : index+defs.IsRightBracketI()])
	i := skipWhitespace(rawData, index+1)

	if i < len(rawData) {
		if defs.IsAssociativeOp(rawData, i) {
			return i, tok, lexOperator
		} else if defs.IsRightBracket(rawData, i) {
			return i, tok, lexRightBracket
		}
	} else {
		return i, tok, nil
	}

	return printError(defs, rawData, i, types.RBR)
}

func lexExpressionWithQuotes(defs Definitions, rawData []rune, index int, doubleQuote bool) (int, int, bool) {
	var isQ func([]rune, int) bool
	if doubleQuote {
		isQ = defs.IsDoubleInvertedComma
	} else {
		isQ = defs.IsSingleInvertedComma
	}

	i := index
	for ; i < len(rawData) && !isQ(rawData, i); i++ {
	}

	endExp := i

	nextIndex := i + 1 // move over quote
	// now skip whitespace
	nextIndex = skipWhitespace(rawData, nextIndex)
	if nextIndex == len(rawData) || defs.IsRightBracket(rawData, nextIndex) || defs.IsAssociativeOp(rawData, nextIndex) {
		return nextIndex, endExp, true
	}

	return i, i, false
}

func lexExpressionWithoutQuotes(defs Definitions, rawData []rune, index int) (int, int, bool) {
	i := index
	trailingWhitespace := 0

	foundNextToken := func(raw []rune, cur int) bool {
		return defs.IsRightBracket(raw, cur) || defs.IsAssociativeOp(raw, cur)
	}

	for ; i < len(rawData) && !(i != index && foundNextToken(rawData, i)); i++ {
		if defs.IsExpRune(rawData, i) {
			trailingWhitespace = 0
		} else if types.IsWhitespace(rawData, i) {
			trailingWhitespace++
		} else {
			return i, i, false
		}
	}

	return i, i - trailingWhitespace, true
}

// Lexes an expression
func lexExpression(defs Definitions, rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token
	tok.Typ = types.EXP

	nextI := index
	endExp := index
	success := false
	startExp := index

	if defs.IsSingleInvertedComma(rawData, index) {
		startExp++
		nextI, endExp, success = lexExpressionWithQuotes(defs, rawData, startExp, false)
	} else if defs.IsDoubleInvertedComma(rawData, index) {
		startExp++
		nextI, endExp, success = lexExpressionWithQuotes(defs, rawData, startExp, true)
	} else {
		nextI, endExp, success = lexExpressionWithoutQuotes(defs, rawData, startExp)
	}

	if !success {
		return printError(defs, rawData, nextI, types.EXP)
	}

	tok.Exp = string(rawData[startExp:endExp])

	nextI = skipWhitespace(rawData, nextI)

	if nextI < len(rawData) {
		if defs.IsAssociativeOp(rawData, nextI) {
			return nextI, tok, lexOperator
		} else if defs.IsRightBracket(rawData, nextI) {
			return nextI, tok, lexRightBracket
		}
	}

	return nextI, tok, nil
}

func getNextFuncAfterOperator(defs Definitions, rawData []rune, index int, tok types.Token) (int, types.Token, stateFn) {
	index = skipWhitespace(rawData, index)
	if index < len(rawData) {
		if defs.IsExpRune(rawData, index) {
			return index, tok, lexExpression
		} else if defs.IsLeftBracket(rawData, index) {
			return index, tok, lexLeftBracket
		}
	}

	return printError(defs, rawData, index+1, types.EXP, types.LBR)
}

func lexNotSymbolAndCreateTrueToken(_ Definitions, _ []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token

	tok.Typ = types.TRUE
	tok.Exp = "true"

	return index, tok, lexNotSymbolAndCreateAndNotToken
}

func lexNotSymbolAndCreateAndNotToken(defs Definitions, rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token

	tok.Typ = types.ANDNOT
	tok.Exp = string(rawData[index : index+defs.IsNotI()])

	i := skipWhitespace(rawData, index+defs.IsNotI())

	return getNextFuncAfterOperator(defs, rawData, i, tok)
}

func lexOperator(defs Definitions, rawData []rune, index int) (int, types.Token, stateFn) {
	var tok types.Token
	indexEnd := 0

	tok.Typ, indexEnd = determineTokenType(defs, rawData, index)
	tok.Exp = string(rawData[index : index+indexEnd])

	return getNextFuncAfterOperator(defs, rawData, index+indexEnd, tok)
}

func determineTokenType(defs Definitions, rawData []rune, index int) (types.TokenType, int) {
	typ := types.UNKNOWN
	typeSize := 0

	// first see if it is an operator
	opType, opLen := determineTokenTypeForOperator(defs, rawData, index)

	if opType != types.UNKNOWN {
		typ = opType
		typeSize = opLen
	} else if defs.IsRightBracket(rawData, index) {
		typ = types.RBR
		typeSize = defs.IsRightBracketI()
	} else if defs.IsLeftBracket(rawData, index) {
		typ = types.LBR
		typeSize = defs.IsLeftBracketI()
	} else if defs.IsSingleInvertedComma(rawData, index) {
		typ = types.SQUOTE
	} else if defs.IsDoubleInvertedComma(rawData, index) {
		typ = types.DQUOTE
	}

	return typ, typeSize
}

func determineTokenTypeForOperator(defs Definitions, rawData []rune, index int) (types.TokenType, int) {
	typ := types.UNKNOWN
	typeSize := 0

	if defs.IsAnd(rawData, index) {
		typ = types.AND
		typeSize = defs.IsAndI()
	} else if defs.IsOr(rawData, index) {
		typ = types.OR
		typeSize = defs.IsOrI()
	}

	if types.CouldBeComplexOp(typ) {
		skipForward := skipWhitespace(rawData, index+typeSize)

		if defs.IsNot(rawData, skipForward) {
			if typ == types.AND {
				typ = types.ANDNOT
			} else if typ == types.OR {
				typ = types.ORNOT
			}

			typeSize = (skipForward - index) + defs.IsNotI()
		}
	}

	return typ, typeSize
}

func printError(defs Definitions, rawData []rune, index int, expected ...types.TokenType) (int, types.Token, stateFn) {
	var found types.TokenType
	index = skipWhitespace(rawData, index)
	if (index + 1) < len(rawData) {
		found, _ = determineTokenType(defs, rawData, index)
	} else if index < len(rawData) {
		found, _ = determineTokenType(defs, rawData, index)
	} else {
		found = types.EOL
	}
	message := tokensToList(defs, expected) + " expected, instead found " + defs.TokToString(found)
	return errorMessage(defs, rawData, index, message)
}

func errorMessage(_ Definitions, rawData []rune, index int, message string) (int, types.Token, stateFn) {
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
		if !types.IsWhitespace(rawData, i) {
			return i
		}
	}

	return i
}

func determineIfExpressionOrLeftBracketOrNot(defs Definitions, rawData []rune, index int) (int, types.Token, stateFn) {
	if defs.IsLeftBracket(rawData, index) {
		return lexLeftBracket(defs, rawData, index)
	} else if defs.IsExpRune(rawData, index) {
		return lexExpression(defs, rawData, index)
	} else if defs.IsNot(rawData, index) {
		return lexNotSymbolAndCreateTrueToken(defs, rawData, index)
	}

	return printError(defs, rawData, index, types.LBR, types.EXP)
}

func getErrorToken() types.Token {
	return types.Token{Typ: types.UNKNOWN}
}

func tokensToList(defs Definitions, expected []types.TokenType) string {
	var res strings.Builder
	for index, exp := range expected {
		if len(expected) > 1 && index == len(expected)-1 {
			res.WriteString(" or ")
		}
		res.WriteString(defs.TokToString(exp))
		if index < len(expected)-2 {
			res.WriteString(", ")
		}

	}
	return res.String()
}
