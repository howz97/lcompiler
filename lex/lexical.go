package lex

import (
	"fmt"
	"io/ioutil"
)

// Error is the error found in lexical analysiis
type Error struct {
	id          int
	row         int
	description string
}

// Token is element of token file
type Token struct {
	ID          int
	Name        []byte
	MachineCode int
	Addr        int

	Row int
}

// Symbol is element of symbol table
type Symbol struct {
	ID   int
	Typ  int
	Name []byte
}

// Analyser is a lexical analyser
type Analyser struct {
	sourceCode    []byte
	currentRow    int
	frontIndex    int
	backwardIndex int

	Tokens    []*Token
	tokenNext int
	SymbolTbl []*Symbol
	errs      []*Error
}

// New return a Lexical instance
func New() *Analyser {
	return &Analyser{
		currentRow: 1,
		Tokens:     make([]*Token, 0, 1<<10),
		SymbolTbl:  make([]*Symbol, 0, 1<<9),
		errs:       make([]*Error, 0, 1<<10),
	}
}

// Analyse start to analyse
func (al *Analyser) Analyse(filename string) {
	var err error
	al.sourceCode, err = ioutil.ReadFile(filename)
	if err != nil {
		panic("file not found: " + filename)
	}
	al.sourceCode = append(al.sourceCode, '#')
	for al.frontIndex < len(al.sourceCode) {
		al.scanStartWithAny()
	}
	al.PrintResult()
}

// PrintResult print the result
func (al *Analyser) PrintResult() {
	for i := 0; i < len(al.Tokens); i++ {
		fmt.Printf("{ID:%d, '%s', code in machine:%d, address:%d}", al.Tokens[i].ID, al.Tokens[i].Name, al.Tokens[i].MachineCode, al.Tokens[i].Addr)
		fmt.Println()
	}
	al.PrintErr()
}

func (al *Analyser) PrintErr() {
	for i := 0; i < len(al.errs); i++ {
		fmt.Printf("[ERROR]id=%d: %s in line %d", al.errs[i].id, al.errs[i].description, al.errs[i].row)
		fmt.Println()
	}
	if len(al.errs) == 0 {
		fmt.Println("\n=======================>> Lexical analysis PASS!!! <<===================")
	} else {
		fmt.Printf("\n=======================>> Lexical analysis total %d errors <<===================", len(al.errs))
	}
}

func (al *Analyser) scanStartWithAny() {
	switch true {
	case !isValid(al.sourceCode[al.frontIndex]):
		al.logError("invalid char '" + string(al.sourceCode[al.frontIndex]) + "'")
		al.frontIndex++
		al.backwardIndex = al.frontIndex
	case al.sourceCode[al.frontIndex] == ' ':
		al.frontIndex++
		al.backwardIndex = al.frontIndex
		return
	case al.sourceCode[al.frontIndex] == '	':
		al.frontIndex++
		al.backwardIndex = al.frontIndex
		return
	case al.sourceCode[al.frontIndex] == '\n':
		al.currentRow++
		al.frontIndex++
		al.backwardIndex = al.frontIndex
		return
	case IsLetter(al.sourceCode[al.frontIndex]):
		al.scanStartWithLetter()
	case isNumber(al.sourceCode[al.frontIndex]):
		al.scanStartWithNumber()
	default:
		al.scanStartWithSymbol()
	}
}

func (al *Analyser) scanStartWithLetter() {
	if !isValid(al.sourceCode[al.frontIndex]) {
		al.logError("invalid char '" + string(al.sourceCode[al.frontIndex]) + "'")
		al.frontIndex++
		al.scanStartWithLetter()
	} else if !(IsLetter(al.sourceCode[al.frontIndex]) || isNumber(al.sourceCode[al.frontIndex])) {
		al.genToken(al.sourceCode[al.backwardIndex:al.frontIndex])
		al.backwardIndex = al.frontIndex
	} else if !isKeywords(al.sourceCode[al.backwardIndex : al.frontIndex+1]) {
		al.frontIndex++
		al.scanStartWithLetter()
	} else if IsLetter(al.sourceCode[al.frontIndex+1]) || isNumber(al.sourceCode[al.frontIndex+1]) {
		al.frontIndex++
		al.scanStartWithLetter()
	} else {
		al.genToken(al.sourceCode[al.backwardIndex : al.frontIndex+1])
		al.frontIndex++
		al.backwardIndex = al.frontIndex
	}
}

func (al *Analyser) scanStartWithNumber() {
	switch true {
	case !isValid(al.sourceCode[al.frontIndex]):
		al.logError("invalid char '" + string(al.sourceCode[al.frontIndex]) + "'")
		al.frontIndex++
		al.scanStartWithNumber()
	case isNumber(al.sourceCode[al.frontIndex]):
		al.frontIndex++
		al.scanStartWithNumber()
	case al.sourceCode[al.frontIndex] == '.':
		if !isNumber(al.sourceCode[al.frontIndex+1]) {
			al.logError("there should be digits after '.'")
			al.genToken(al.sourceCode[al.backwardIndex:al.frontIndex])
			al.frontIndex++
			al.backwardIndex = al.frontIndex
		} else {
			al.frontIndex++
			al.scanFloat()
		}
	default:
		al.genToken(al.sourceCode[al.backwardIndex:al.frontIndex])
		al.backwardIndex = al.frontIndex
	}
}

func (al *Analyser) scanFloat() {
	switch true {
	case !isValid(al.sourceCode[al.frontIndex]):
		al.frontIndex++
		al.scanFloat()
	case isNumber(al.sourceCode[al.frontIndex]):
		al.frontIndex++
		al.scanFloat()
	case al.sourceCode[al.frontIndex] == '.':
		al.logError("'.' can not be appended to float")
		al.genToken(al.sourceCode[al.backwardIndex:al.frontIndex])
		al.frontIndex++
		al.backwardIndex = al.frontIndex
	default:
		al.genToken(al.sourceCode[al.backwardIndex:al.frontIndex])
		al.backwardIndex = al.frontIndex
	}
}

func (al *Analyser) scanStartWithSymbol() {
	switch true {
	case !isValid(al.sourceCode[al.frontIndex]):
		al.logError("invalid char '" + string(al.sourceCode[al.frontIndex]) + "'")
		al.frontIndex++
		al.scanStartWithSymbol()
	case isLessGreaterColon(al.sourceCode[al.frontIndex]):
		al.frontIndex++
		if al.sourceCode[al.frontIndex] == '=' || (al.sourceCode[al.frontIndex-1] == '<' && al.sourceCode[al.frontIndex] == '>') {
			al.genToken(al.sourceCode[al.backwardIndex : al.frontIndex+1])
			al.frontIndex++
			al.backwardIndex = al.frontIndex
		} else {
			al.genToken(al.sourceCode[al.backwardIndex:al.frontIndex])
			al.backwardIndex = al.frontIndex
		}
	case al.sourceCode[al.frontIndex] == '\n':
		al.currentRow++
		al.frontIndex++
	default:
		al.genToken(al.sourceCode[al.backwardIndex : al.frontIndex+1])
		al.frontIndex++
		al.backwardIndex = al.frontIndex
	}
}

func (al *Analyser) genToken(name []byte) {
	mcode := getMachineCode(name)
	t := &Token{
		ID:          len(al.Tokens),
		Name:        name,
		MachineCode: mcode,
		Addr:        al.addSymbol(name, mcode),
		Row:         al.currentRow,
	}
	al.Tokens = append(al.Tokens, t)
}

func (al *Analyser) addSymbol(name []byte, mcode int) int {
	if mcode < 18 || mcode > 20 {
		return -1
	}
	for i := 0; i < len(al.SymbolTbl); i++ {
		if string(name) == string(al.SymbolTbl[i].Name) {
			return al.SymbolTbl[i].ID
		}
	}
	sym := &Symbol{
		ID:   len(al.SymbolTbl),
		Typ:  mcode,
		Name: name,
	}
	al.SymbolTbl = append(al.SymbolTbl, sym)
	return len(al.SymbolTbl) - 1
}

func (al *Analyser) logError(desp string) {
	e := &Error{
		id:          len(al.errs),
		row:         al.currentRow,
		description: desp,
	}
	al.errs = append(al.errs, e)
}

func (al *Analyser) IsTokenMatch(c int) bool {
	m := al.Tokens[al.tokenNext].MachineCode == c

	// TODO: for test
	// fmt.Println(c, "meet", al.Tokens[al.tokenNext].MachineCode)
	if !m {
		fmt.Printf("\n[ERROR]: should meet '%s' rather than '%s' in line %d", ReversCodeMap[c], string(al.Tokens[al.tokenNext].Name), al.Tokens[al.tokenNext].Row)
	}

	al.tokenNext++
	return m
}

func (al *Analyser) NextToken() int {
	return al.tokenNext
}

// NextCode return the machine code scaned in gramer analysis
func (al *Analyser) NextCode() int {
	return al.Tokens[al.tokenNext].MachineCode
}
