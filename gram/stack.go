package gram

import (
	"fmt"

	"github.com/zh1014/lcompiler/lex"
)

type analyStack struct {
	chars []*gramChar
	cap   int
}

func newAnalyStack(size int) *analyStack {
	return &analyStack{
		chars: make([]*gramChar, size),
	}
}

func (as *analyStack) pushOne(gc *gramChar) {
	if gc == nil {
		panic("pushing nil gramChar to stack")
	}

	if gc.isTerminal {
		fmt.Println("压入：", gc.t)
	} else {
		fmt.Println("压入：", gc.nt)

	}
	as.chars[as.cap] = gc
	as.cap++
}

func (as *analyStack) push(gcs []*gramChar) {
	if gcs == nil {
		panic("gcs is nil")
	}
	if gcs[0].isTerminal && gcs[0].t == lex.CodeMap["ε"] {
		return
	}
	for i := len(gcs) - 1; i >= 0; i-- {

		// TODO: delete
		// if gcs[i].isTerminal {
		// 	fmt.Println(gcs[i].t)
		// } else {
		// 	fmt.Println(gcs[i].nt)
		// }

		as.pushOne(gcs[i])
	}
}

func (as *analyStack) pop() *gramChar {
	if as.cap < 1 {
		return nil
	}
	as.cap--
	return as.chars[as.cap]
}
