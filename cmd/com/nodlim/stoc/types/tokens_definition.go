package types

import (
	"github.com/kranzuft/stoc/cmd/com/nodlim/stoc/commons"
)

type TokensDefinition map[TokenType]TokenInfo

func (td *TokensDefinition) defineTokenInfo(typ TokenType, key string, description string) *TokensDefinition {
	(*td)[typ] = TokenInfo{
		desc: description,
		key:  []rune(key),
		size: len([]rune(key)),
	}
	return td
}

func (td *TokensDefinition) finalise() TokensDefinition {
	return *td
}

func (td TokensDefinition) IsLeftBracket(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[LBR].key)
}

func (td TokensDefinition) IsLeftBracketI() int {
	return td[LBR].size
}

func (td TokensDefinition) IsRightBracket(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[RBR].key)
}

func (td TokensDefinition) IsRightBracketI() int {
	return td[RBR].size
}

func (td TokensDefinition) IsAnd(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[AND].key)
}

func (td TokensDefinition) IsAndI() int {
	return td[AND].size
}

func (td TokensDefinition) IsOr(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[OR].key)
}

func (td TokensDefinition) IsOrI() int {
	return td[OR].size
}

func (td TokensDefinition) IsNot(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[NOT].key)
}

func (td TokensDefinition) IsNotI() int {
	return td[NOT].size
}

func (td TokensDefinition) IsExpRune(r []rune, index int) bool {
	return !(IsWhitespace(r, index) || td.IsKeyword(r, index))
}

func (td TokensDefinition) IsAssociativeOp(r []rune, index int) bool {
	return td.IsOr(r, index) || td.IsAnd(r, index)
}

func (td TokensDefinition) IsKeyword(r []rune, index int) bool {
	return td.IsOr(r, index) || td.IsAnd(r, index) || td.IsLeftBracket(r, index) || td.IsRightBracket(r, index) || td.IsNot(r, index)
}

func (td TokensDefinition) IsQuote(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[SQUOTE].key, td[DQUOTE].key)
}

func (td TokensDefinition) IsSingleInvertedComma(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[SQUOTE].key)
}

func (td TokensDefinition) IsDoubleInvertedComma(r []rune, index int) bool {
	return commons.StartsWith(r, index, td[DQUOTE].key)
}

// TokToString get the TokenInfo description for the token type by active tokens definition
func (td TokensDefinition) TokToString(t TokenType) string {
	return td[t].desc
}

func prepareDefaultTokensDefinition() TokensDefinition {
	def := TokensDefinition{}
	return def.
		defineTokenInfo(AND, "&", "and").
		defineTokenInfo(OR, "|", "or").
		defineTokenInfo(NOT, "!", "not").
		defineTokenInfo(ANDNOT, "&!", "and not").
		defineTokenInfo(ORNOT, "&!", "or not").
		defineTokenInfo(TRUE, "True", "true").
		defineTokenInfo(LBR, "(", "left bracket").
		defineTokenInfo(RBR, ")", "right bracket").
		defineTokenInfo(EOL, "\n", "end of line").
		defineTokenInfo(EXP, "", "expression").
		defineTokenInfo(DQUOTE, "\"", "double inverted comma").
		defineTokenInfo(SQUOTE, "'", "single inverted comma").
		finalise()
}

var DefaultTokensDefinition = prepareDefaultTokensDefinition()
