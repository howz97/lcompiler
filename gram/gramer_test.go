package gram

import (
	"fmt"
	"testing"

	"github.com/zh1014/lcompiler/lex"
)

func TestGramer(t *testing.T) {
	al := &Analyser{}
	al.parseGramer("./gramer.txt")
	for k, v := range al.lang {
		fmt.Print(k, ":", v.idx, "--->")
		for _, v1 := range v.right {
			for _, v2 := range v1 {
				if v2.isTerminal {
					fmt.Print(v2.t, " ")
				} else {
					fmt.Print(v2.nt, " ")
				}
			}
			fmt.Print("| ")
		}
		fmt.Println()
	}
}

func TestFATbl(t *testing.T) {
	// New("./gramer.txt")
	al := New("./gramer.txt")
	// al.PrintFirstSet()
	// al.PrintFollowSet()
	al.PrintFATbl()
}

func Test_parseTerm(t *testing.T) {
	remain, name := parseTerm("if$布尔表达式$then$执行句$|if$布尔表达式$then$执行句$else$执行句$")
	fmt.Println(remain, name)
}

func Test_firstFollowSet(t *testing.T) {
	al := New("./gramer.txt")
	al.PrintFirstSet()
	al.PrintFollowSet()
}

func Test_Gramer(t *testing.T) {
	lexAl := lex.New()
	lexAl.Analyse("../lex/lex.txt")
	gramAl := New("./gramer.txt")
	gramAl.Analyse(lexAl)
}
