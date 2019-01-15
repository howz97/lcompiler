package gram

import (
	"strings"

	"github.com/zh1014/lcompiler/lex"
)

func parseNTerm(raw string) (remain string, name string) {
	raw = raw[1:]
	idx := strings.Index(raw, "$")
	nTermName := raw[:idx]
	raw = raw[idx+1:]
	return raw, nTermName
}

func parseTerm(raw string) (remain string, name string) {
	if raw == "" {
		panic("parseTerm: raw is empty")
	}
	for i := len(raw); i > 0; i-- {
		if lex.CodeMap[raw[:i]] != 0 {
			return raw[i:], raw[:i]
		}
	}
	panic("error not match")
}

// type skipDetector struct {
// 	recordTbl map[string]bool
// }

// func newDetector(d string) *skipDetector {
// 	sd := &skipDetector{
// 		recordTbl: make(map[string]bool),
// 	}
// 	sd.recordTbl[d] = true
// 	return sd
// }

// func (sd *skipDetector) detect(d string) bool {
// 	_, exist := sd.recordTbl[d]
// 	if exist {
// 		return true
// 	}
// }

// isSkip check whether we need to call g.getFollowSet() for get f.followSet
func isSkip(g, f *gramer) bool {
	return f.willMeet(g.left)
}

func (g *gramer) willMeet(tgt string) bool {
	if g.left == tgt {
		return true
	}
	if g.holdFollowSet {
		return false
	}
	// force traversal
	for _, gram := range g.lang {
		for _, partGram := range gram.right {
			for i, gramc := range partGram {
				if !gramc.isTerminal && gramc.nt == g.left {
					if i == len(partGram)-1 {
						if gram.left != gramc.nt {
							if gram.willMeet(tgt) {
								return true
							}
						}
					} else {
						if !partGram[i+1].isTerminal {
							if g.lang[partGram[i+1].nt].getFirstSet().where(lex.CodeMap["ε"]) > 0 {
								if g.lang[partGram[i+1].nt].willMeet(tgt) {
									return true
								}
							}
						}
					}
				}
			}
		}
	}
	return false
}
