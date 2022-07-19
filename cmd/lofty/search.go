package lofty

import (
	"fmt"
	"lofty/cmd/lofty/lexer"
	"lofty/cmd/lofty/types"
)

func SearchString(command string, target string) (bool, error) {
	return SearchStringCustom(types.DefaultTokensDefinition, command, target)
}

func SearchStringCustom(defs types.TokensDefinition, command string, target string) (bool, error) {
	var raw = []rune(command)

	lexSuccess, tokens := lexer.BooleanAlgebraLexer(defs, raw)

	fmt.Println(tokens)

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
