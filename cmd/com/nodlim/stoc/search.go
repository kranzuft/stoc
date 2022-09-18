// Package stoc: search text on condition
// Uses boolean algebra to search input text.
// Whether the boolean algebra 'condition' is met determines results: true or false.
//
// stoc is not opinionated regarding the syntax used, only the grammar rules.
// The syntax can be configured using types.TokensDefinition functions alongside the SearchStringCustom() function.
// For an example see types.DefaultTokensDefinition.
//
// The default syntax for a condition can be described as follows:
// - conjunction (and): 			&
// - disjunction (or): 				|
// - inversion 	(not): 				!
// - brackets (groupings): 			(, )
// - expressions (filter-rules): 	list of unicode chars, optionally surrounded by quotes
//
// For the coded version of the default syntax, see types.DefaultTokensDefinition
// For a formalised version of the default syntax (ebnf or railroad), see the design/ folder in the source code
package stoc

import (
	"github.com/kranzuft/boolean-algebra-to-tokens/cmd/com/nodlim/batt/lexer"
	"github.com/kranzuft/boolean-algebra-to-tokens/cmd/com/nodlim/batt/pos_error"
	"github.com/kranzuft/boolean-algebra-to-tokens/cmd/com/nodlim/batt/types"
	"strings"
)

// PreparedTokens tokens that have been pre-prepared and in post-fix notation
type PreparedTokens []types.Token

// SearchString will search through the contents of target arg based on the command arg.
// The command string must be a valid condition.
//
// The condition must follow the syntax of types.DefaultTokensDefinition and match the standard condition grammar.
// The target string has no requirements.
//
// Example arguments returning true: "('dog' or 'cat') and not 'frog'"`, "Is it raining cats and dogs?"
// Example arguments returning false: "('dog' or 'cat') and not 'frog'"`, "It's raining cats, dogs, and even frogs!"
func SearchString(command string, target string) (bool, pos_error.PosError) {
	return SearchStringCustom(types.DefaultTokensDefinition, command, target)
}

// SearchStringCustom will search through the contents of target string based on the command string.
// The command string must be a valid condition.
// The target string has no requirements.
//
// The condition must follow the syntax defined by defs arg and match the standard condition grammar.
func SearchStringCustom(defs types.TokensDefinition, command string, target string) (bool, pos_error.PosError) {
	preparedTokens, err := LexIntoTokens(defs, command)

	if err == nil {
		return SearchTokens(preparedTokens, target), err
	}

	return false, err
}

// SearchTokens searches target with pre-prepared tokens
func SearchTokens(preparation PreparedTokens, target string) bool {
	return SearchPostfixTokens(preparation, target)
}

// LexIntoTokens produces postfix tokens from TokensDefinition and raw tokens-to-be command
func LexIntoTokens(defs types.TokensDefinition, command string) (PreparedTokens, pos_error.PosError) {
	tokens, errLex := lexer.BooleanAlgebraLexer(defs, []rune(command))

	if errLex == nil {
		result, errShunt := lexer.TokenShuntingAlgorithm(tokens)

		if errShunt == nil {
			return result, errShunt
		} else {
			return nil, errShunt
		}
	} else {
		return nil, errLex
	}
}

// SearchPostfixTokens takes the postfix formatted tokens and parses the tokens. It uses the tokens to apply conditions on
// the 'target' string. If the conditions aren't met then
// Recommended to not be called directly. Instead, call stoc.SearchString or stoc.SearchStringCustom.
func SearchPostfixTokens(search []types.Token, target string) bool {
	var stack []bool
	for _, tok := range search {
		// action := "Apply op to top of stack"
		switch tok.Typ {
		case types.AND:
			res := stack[len(stack)-2] && stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		case types.OR:
			res := stack[len(stack)-2] || stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		case types.ANDNOT:
			res := stack[len(stack)-2] && !stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		case types.ORNOT:
			res := stack[len(stack)-2] || !stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		default:
			// action = "Push num onto top of stack"
			if tok.Typ != types.TRUE {
				stack = append(stack, strings.Contains(target, tok.Exp))
			} else {
				stack = append(stack, true)
			}
		}
	}

	return stack[0]
}
