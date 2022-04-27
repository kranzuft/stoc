package lofty

import "lofty/cmd/lofty/processing"

func SearchString(command string, target string) (bool, error) {
	var raw = []rune(command)

	lexSuccess, tokens := processing.BooleanAlgebraLexer(raw)

	if lexSuccess {
		result, err := processing.TokenShuntingAlgorithm(tokens)

		if err == nil {
			return processing.SearchPostfixTokens(result, target), err
		} else {
			return false, err
		}
	}

	return false, nil
}
