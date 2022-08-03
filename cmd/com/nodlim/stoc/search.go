package stoc

import (
	"stoc/cmd/com/nodlim/stoc/lexer"
	"stoc/cmd/com/nodlim/stoc/types"
)

func SearchString(command string, target string) (bool, error) {
	return SearchStringCustom(types.DefaultTokensDefinition, command, target)
}

func SearchStringCustom(defs types.TokensDefinition, command string, target string) (bool, error) {
	var raw = []rune(command)

	lexSuccess, tokens := lexer.BooleanAlgebraLexer(defs, raw)

	if lexSuccess {
		result, err := lexer.TokenShuntingAlgorithm(tokens)

		if err == nil {
			return lexer.SearchPostfixTokens(result, target), err
		} else {
			return false, err
		}
	}

	return false, nil
}
