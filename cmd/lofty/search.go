package lofty

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func searchPipeMode() {
	reader := bufio.NewReader(os.Stdin)
	command := strings.Join(os.Args[1:], " ")
	cont := true
	for cont {
		line, err := reader.ReadString('\n')
		if err != nil && err == io.EOF {
			cont = false
		}
		if SearchString(command, line) {
			fmt.Println(strings.Trim(line, "\n"))
		}
	}
}

func splitAndSearch(command string, target string, sep string) bool {
	for _, line := range strings.Split(target, sep) {
		SearchString(command, line)
	}
	return true
}

func SearchString(command string, target string) bool {
	var raw = []rune(command)

	lexSuccess, tokens := booleanAlgebraLexer(raw)

	if lexSuccess {
		shuntingSuccess, result := tokenShuntingAlgorithm(tokens)

		if shuntingSuccess {
			return searchPostfixTokens(result, target)
		} else {
			fmt.Println("The calculation was not successful, check that you aren't missing a bracket")
		}
	}

	return false
}
