package lofty

import (
	"fmt"
	"strings"
)

type parsedTokens struct {
	parsed   []token
	operator []token
}

// Shunting algorithm, based on the mathematical implementation available on the shunting algorithm wiki page
//
// A snapshot of that algorithm is below (this implementation is a simplified version with some unincluded additions):
//
// while there are parsed to be read:
//     read a token.
//     if the token is a number, then:
//         push it to the output queue.
//     else if the token is a function then:
//         push it onto the operator stack
//     else if the token is an operator then:
//         while ((there is an operator at the top of the operator stack)
//               and ((the operator at the top of the operator stack has greater precedence)
//                   or (the operator at the top of the operator stack has equal precedence and the token is left associative))
//               and (the operator at the top of the operator stack is not a left parenthesis)):
//             pop operators from the operator stack onto the output queue.
//         push it onto the operator stack.
//     else if the token is a left parenthesis (i.e. "("), then:
//         push it onto the operator stack.
//     else if the token is a right parenthesis (i.e. ")"), then:
//         while the operator at the top of the operator stack is not a left parenthesis:
//             pop the operator from the operator stack onto the output queue.
//         -> If the stack runs out without finding a left parenthesis, then there are mismatched parentheses.
//         if there is a left parenthesis at the top of the operator stack, then:
//             pop the operator from the operator stack and discard it
//         if there is a function token at the top of the operator stack, then:
//             pop the function from the operator stack onto the output queue.
// -> After while loop, if operator stack not null, pop everything to output queue
// if there are no more parsed to read then:
//     while there are still operator parsed on the stack:
//         -> If the operator token on the top of the stack is a parenthesis, then there are mismatched parentheses.
//         pop the operator from the operator stack onto the output queue.
// exit.

func tokenShuntingAlgorithm(toks []token) (bool, []token) {
	var parsed   []token
	var operator []token
	for _, tok := range toks {
		// This implementation does not implement composite functions, functions with variable number of arguments,
		// and unary operators.
		if tok.typ == EXP || tok.typ == TRUE {
			parsed = append(parsed, tok)
		} else if isOp(tok.typ) {
			for len(operator) > 0 && (isLeftAssociative(end(operator).typ) || isComplexOp(end(operator).typ)) {
				parsed, operator = moveEnd(parsed, operator)
			}
			operator = append(operator, tok)
		} else if tok.typ == LBR {
			operator = append(operator, tok)
		} else if tok.typ == RBR {
			for len(operator) > 0 && operator[len(operator)-1].typ != LBR {
				parsed, operator = moveEnd(parsed, operator)
			}

			// If the stack runs out without finding a left parenthesis, then there are mismatched parentheses.
			if len(operator) == 0 {
				return false, parsed
			} else if operator[len(operator)-1].typ == LBR {
				operator = operator[:endIndex(operator)]
			}
		}
	}

	// After while loop, add every operator to output queue
	// (didn't bother popping from stack, just read array from back)
	for i, _ := range operator {
		opI := operator[len(operator)-1-i] // reverse because this is a stack
		// If the operator token on the top of the stack is a parenthesis, then there are mismatched parentheses.
		if opI.typ == LBR || opI.typ == RBR {
			return false, parsed
		}
		parsed = append(parsed, opI)
	}

	fmt.Println(tokensToString(parsed))

	return true, parsed
}

func searchPostfixTokens(search []token, target string) bool {
	var stack []bool
	for _, tok := range search {
		// action := "Apply op to top of stack"
		switch tok.typ {
		case AND:
			res := stack[len(stack)-2] && stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		case OR:
			res := stack[len(stack)-2] || stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		case ANDNOT:
			res := stack[len(stack)-2] && !stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		case ORNOT:
			res := stack[len(stack)-2] || !stack[len(stack)-1]
			stack[len(stack)-2] = res
			stack = stack[:len(stack)-1]
		default:
			// action = "Push num onto top of stack"
			if tok.typ != TRUE {
				stack = append(stack, strings.Contains(target, tok.exp))
			} else {
				stack = append(stack, true)
			}
		}
		// fmt.Printf("%3s    %-26s  %v\n", tok, action, stack)
	}

	return stack[0]
}

/**
 * Utility methods only used by shunting algorithm
 */

func endIndex(tokens []token) int {
	return len(tokens) - 1
}

func end(tokens []token) token {
	return tokens[endIndex(tokens)]
}

func pop(tokens []token) (token, []token) {
	tok := end(tokens)
	return tok, tokens[:endIndex(tokens)]
}

func moveEnd(as []token, bs []token) ([]token, []token) {
	b, bs := pop(bs)
	bs = bs
	as = append(as, b)
	return as, bs
}
