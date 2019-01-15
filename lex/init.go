package lex

var (
	CodeMap map[string]int

	// temporary variable, deleted when refactoring the code
	ReversCodeMap map[int]string
)

func init() {
	CodeMap = make(map[string]int, 39)
	CodeMap["and"] = 1
	CodeMap["begin"] = 2
	CodeMap["bool"] = 3
	CodeMap["do"] = 4
	CodeMap["else"] = 5
	CodeMap["end"] = 6
	CodeMap["false"] = 7
	CodeMap["if"] = 8
	CodeMap["integer"] = 9
	CodeMap["not"] = 10
	CodeMap["or"] = 11
	CodeMap["program"] = 12
	CodeMap["real"] = 13
	CodeMap["then"] = 14
	CodeMap["true"] = 15
	CodeMap["var"] = 16
	CodeMap["while"] = 17
	CodeMap["标识符"] = 18
	CodeMap["整数"] = 19
	CodeMap["实数"] = 20
	CodeMap["("] = 21
	CodeMap[")"] = 22
	CodeMap["+"] = 23
	CodeMap["-"] = 24
	CodeMap["*"] = 25
	CodeMap["/"] = 26
	CodeMap["."] = 27
	CodeMap[","] = 28
	CodeMap[":"] = 29
	CodeMap[";"] = 30
	CodeMap[":="] = 31
	CodeMap["="] = 32
	CodeMap["<="] = 33
	CodeMap["<"] = 34
	CodeMap["<>"] = 35
	CodeMap[">"] = 36
	CodeMap[">="] = 37
	CodeMap["#"] = 38
	CodeMap["ε"] = 39

	ReversCodeMap = make(map[int]string, 39)
	for k, v := range CodeMap {
		ReversCodeMap[v] = k
	}
}
