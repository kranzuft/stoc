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
	"errors"
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc/lexer"
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc/types"
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
func SearchString(command string, target string) (bool, error) {
	return SearchStringCustom(types.DefaultTokensDefinition, command, target)
}

// SearchStringCustom will search through the contents of target string based on the command string.
// The command string must be a valid condition.
// The target string has no requirements.
//
// The condition must follow the syntax defined by defs arg and match the standard condition grammar.
func SearchStringCustom(defs types.TokensDefinition, command string, target string) (bool, error) {
	preparedTokens, err := LexIntoTokens(defs, command)

	if err == nil {
		return SearchTokens(preparedTokens, target), err
	}

	return false, err
}

// SearchTokens searches target with pre-prepared tokens
func SearchTokens(preparation PreparedTokens, target string) bool {
	return lexer.SearchPostfixTokens(preparation, target)
}

// LexIntoTokens produces postfix tokens from TokensDefinition and raw tokens-to-be command
func LexIntoTokens(defs types.TokensDefinition, command string) (PreparedTokens, error) {
	lexSuccess, tokens := lexer.BooleanAlgebraLexer(defs, []rune(command))

	if lexSuccess {
		result, err := lexer.TokenShuntingAlgorithm(tokens)

		if err == nil {
			return result, err
		} else {
			return nil, err
		}
	}

	return nil, errors.New("couldn't lex command")
}
