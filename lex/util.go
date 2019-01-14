package lex

import "strings"

var (
	charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789+-*/().,:;=<> \n	#"
)

func isValid(c byte) bool {
	return strings.Contains(charset, string(c))
}

func getMachineCode(name []byte) int {
	code := CodeMap[string(name)]
	if code != 0 {
		return code
	}
	switch true {
	case IsLetter(name[0]):
		return 18
	case strings.Contains(string(name), "."):
		return 20
	default:
		return 19
	}
}

func IsLetter(c byte) bool {
	return (c >= 65 && c <= 90) || (c >= 97 && c <= 122)
}

func isNumber(c byte) bool {
	return c >= 48 && c <= 57
}

func isLessGreaterColon(c byte) bool {
	return c == '<' || c == '>' || c == ':'
}

func isKeywords(name []byte) bool {
	mcode, exist := CodeMap[string(name)]
	if !exist {
		return false
	}
	return mcode >= 1 && mcode <= 17
}
