package gram

import (
	"fmt"
	"io/ioutil"
	"strings"

	"github.com/zh1014/data-structure/bitmap"
	"github.com/zh1014/lcompiler/lex"
)

// smallist unit of gramer
type gramChar struct {
	t          int
	nt         string
	isTerminal bool
}

// corresponsd to gramer of language, just like A -> aA|b|c
type gramer struct {
	idx           int
	left          string
	right         [][]*gramChar
	firstSet      sets
	followSet     sets
	holdFollowSet bool

	lang map[string]*gramer // TODO: 设计有问题；为嘛存这个？ 应该把[][]*gramChar换成[][]*gramer 使gramer代表终结符与非终结符
}

func (g *gramer) getFirstSet() sets {
	if g.firstSet == nil {
		g.newFirstSet()
	} else {
		return g.firstSet
	}
	for i0, partRight := range g.right {
		g.firstSet.add(i0, lex.CodeMap["ε"])
		for i := 0; i < len(partRight) && g.firstSet.where(lex.CodeMap["ε"]) >= 0; i++ {
			g.firstSet.remove(lex.CodeMap["ε"])
			if partRight[i].isTerminal {
				g.firstSet.add(i0, partRight[i].t)
			} else {
				g.firstSet.or(i0, g.lang[partRight[i].nt].getFirstSet())
			}
		}
	}
	return g.firstSet
}

type sets []*bitmap.BitMap

func (g *gramer) newFirstSet() {
	firstSet := make([]*bitmap.BitMap, len(g.right))
	for i := range firstSet {
		firstSet[i] = bitmap.NewBitMap(uint32(len(lex.CodeMap) + 1))
	}
	g.firstSet = sets(firstSet)
}

func (g *gramer) newFollowSet() {
	followSet := make([]*bitmap.BitMap, 1)
	followSet[0] = bitmap.NewBitMap(uint32(len(lex.CodeMap) + 1))
	g.followSet = sets(followSet)
}

func (s sets) add(idx int, v int) {
	s[idx].Set(uint32(v))
}

func (s sets) or(idx int, ss sets) {
	for i := range ss {
		s[idx].Or(ss[i].ByteSlice)
	}
}

func (s sets) where(v int) int {
	w := -1
	for i := range s {
		if s[i].Check(uint32(v)) {
			w = i
			break
		}
	}
	return w
}

func (s sets) remove(v int) {
	for i := range s {
		s[i].Unset(uint32(v))
	}
}

func (al *Analyser) getAllFollowSet() {
	// 程序 must generate follow set firstly
	al.lang["程序"].getFollowSet()
	al.lang["程序体"].getFollowSet()
	al.lang["变量说明"].getFollowSet()
	al.lang["变量定义"].getFollowSet()
	al.lang["标识符表"].getFollowSet()
	al.lang["额外标识符"].getFollowSet()
	al.lang["复合句"].getFollowSet()
	al.lang["语句表"].getFollowSet()
	al.lang["执行句"].getFollowSet()
	al.lang["简单句"].getFollowSet()
	al.lang["赋值句"].getFollowSet()
	al.lang["变量"].getFollowSet()
	al.lang["结构句"].getFollowSet()
	al.lang["if句"].getFollowSet()
	al.lang["else句"].getFollowSet()
	al.lang["while句"].getFollowSet()
	al.lang["表达式"].getFollowSet()
	al.lang["算术表达式"].getFollowSet()
	al.lang["算术表达式'"].getFollowSet()
	al.lang["加减"].getFollowSet()
	al.lang["项"].getFollowSet()
	al.lang["项'"].getFollowSet()
	al.lang["乘除"].getFollowSet()
	al.lang["因子"].getFollowSet()
	al.lang["算术量"].getFollowSet()
	al.lang["布尔表达式"].getFollowSet()
	al.lang["布尔表达式'"].getFollowSet()
	al.lang["布尔项"].getFollowSet()
	al.lang["布尔项'"].getFollowSet()
	al.lang["布尔因子"].getFollowSet()
	al.lang["布尔量"].getFollowSet()
	al.lang["标识符或关系表达式"].getFollowSet()
	al.lang["关系表达式"].getFollowSet()
	al.lang["关系运算符"].getFollowSet()
	al.lang["类型名"].getFollowSet()
	al.lang["布尔常数"].getFollowSet()
}

func (g *gramer) getFollowSet() sets {
	if g.holdFollowSet {
		return g.followSet
	}
	g.newFollowSet()
	if g.left == "程序" {
		g.followSet.add(0, lex.CodeMap["#"])
		return g.followSet
	}
	// force traversal
	// fmt.Println("。。。。。。开始找", g.left, "的follow集")
	for _, gram := range g.lang {
		for _, partGram := range gram.right {
			for i, gramc := range partGram {
				if !gramc.isTerminal && gramc.nt == g.left {
					// fmt.Println("乱序遍历到" + k + "开头这一行")
					// fmt.Println("找第", i0+1, "个短句")
					// fmt.Println("正在检测", gramc.nt)
					if i == len(partGram)-1 {
						if gram.left != gramc.nt && !isSkip(g, gram) {
							// fmt.Println("加入", gram.left, "的follow集")
							g.followSet.or(0, gram.getFollowSet())
						}
					} else {
						if partGram[i+1].isTerminal {
							// fmt.Println("加入", partGram[i+1].t)
							g.followSet.add(0, partGram[i+1].t)
						} else {
							// fmt.Println("加入", partGram[i+1].nt, "的first集")
							g.followSet.or(0, g.lang[partGram[i+1].nt].getFirstSet())
							if g.lang[partGram[i+1].nt].getFirstSet().where(lex.CodeMap["ε"]) >= 0 && !isSkip(g, g.lang[partGram[i+1].nt]) {
								// fmt.Println("加入", partGram[i+1].nt, "的follow集")
								g.followSet.or(0, g.lang[partGram[i+1].nt].getFollowSet())
							}
						}
					}
				}
			}
		}
	}
	g.holdFollowSet = true
	return g.followSet
}

// formula position specify a gramer and the index.
// gramer is made up of formulas
type formuPst struct {
	g   string
	idx int
}

// Analyser is a gramer analyser
type Analyser struct {
	lang             map[string]*gramer
	nonTerminalRoot  string
	forecastAnalyTbl [50][39]*formuPst
}

// New generate a gramer analyser instance
func New(gramerfile string) *Analyser {
	al := &Analyser{}
	al.parseGramer(gramerfile)
	al.genForecastAnalyTbl(al.lang)
	return al
}

// Analyse analyse the gramer of token file
func (al *Analyser) Analyse(lexAl *lex.Analyser) []string {
	if !al.isGramPass(lexAl) {
		fmt.Println("\n=======================>> Gramer Analysis NOT PASS <<===================")
		return nil
	}
	fmt.Println("\n=======================>> Gramer Analysis PASS!!! <<===================")
	return nil
}

func (al *Analyser) PrintFirstSet() {
	fmt.Println("\n\n					打印 FIRST SET:")
	for k, gram := range al.lang {
		fmt.Println(k + ": ")
		for i := 1; i <= len(lex.CodeMap); i++ {
			if gram.firstSet.where(i) >= 0 {
				fmt.Print("  '" + lex.ReversCodeMap[i] + "'")
			}
		}
		fmt.Println()
	}
}

func (al *Analyser) PrintFollowSet() {
	fmt.Println("\n\n					打印 FOLLOW SET:")
	for k, gram := range al.lang {
		if gram.followSet == nil {
			panic(k + ".followSet is nil")
		}
		fmt.Println(k + ":")
		for i := 1; i <= len(lex.CodeMap); i++ {
			if gram.followSet.where(i) >= 0 {
				fmt.Print("  '" + lex.ReversCodeMap[i] + "'")
			}
		}
		fmt.Println()
	}
}

func (al *Analyser) PrintFATbl() {
	fmt.Println("\n\n					打印 预测分析表:")
	for k, gram := range al.lang {
		fmt.Println("*****" + k + "*****")
		for j := 1; j < 39; j++ {
			if al.forecastAnalyTbl[gram.idx][j] != nil {
				fmt.Printf("遇到'%s', 选择第%d个产生式\n", lex.ReversCodeMap[j], al.forecastAnalyTbl[gram.idx][j].idx+1)
			}
		}
	}
}

func (al *Analyser) parseGramer(gramerfile string) {
	lang := make(map[string]*gramer)

	file, err := ioutil.ReadFile(gramerfile)
	if err != nil {
		panic("Can not find gramer file: " + gramerfile)
	}

	gramStr := string(file)
	if drop := strings.Index(gramStr, "#"); drop >= 0 {
		gramStr = gramStr[0:drop]
	}
	gramStr = strings.Replace(gramStr, " ", "", -1)
	gramStr = strings.Replace(gramStr, "	", "", -1)

	gramSli := strings.Split(gramStr, "\n") // gramSli 的每一个元素为一行语法
	for i := 0; i < len(gramSli); i++ {
		if gramSli[i] == "" {
			continue
		}

		var nTermName string
		gramSli[i], nTermName = parseNTerm(gramSli[i])
		if gramSli[i] == "" || nTermName == "" {
			panic("error nonterminal symbol missing")
		}
		g := &gramer{
			idx:   i + 1, // start from 1
			left:  nTermName,
			right: make([][]*gramChar, 0, 10),
			lang:  lang,
		}

		if g.idx == 1 {
			al.nonTerminalRoot = nTermName
		}

		rightSli := strings.Split(gramSli[i], "|") // rightSli的每一个元素为一个小短句
		for j := 0; j < len(rightSli); j++ {
			gramChars := make([]*gramChar, 0, 6)
			for k := 0; rightSli[j] != ""; k++ {
				var na string
				gc := &gramChar{}
				if rightSli[j][0] == '$' {
					rightSli[j], na = parseNTerm(rightSli[j])
					gc.nt = na
					gc.isTerminal = false
				} else {
					rightSli[j], na = parseTerm(rightSli[j])
					gc.t = lex.CodeMap[na]
					if gc.t == 0 {
						gc.t = lex.CodeMap["ε"]
					}
					gc.isTerminal = true
				}
				gramChars = append(gramChars, gc)
			}
			g.right = append(g.right, gramChars)
		}
		lang[nTermName] = g
	}
	al.lang = lang
}

func (al *Analyser) genForecastAnalyTbl(lang map[string]*gramer) {
	for _, gram := range al.lang {
		gram.getFirstSet()
	}

	// TODO: print
	al.PrintFirstSet()

	al.getAllFollowSet()

	// TODO: print
	al.PrintFollowSet()

	al.forecastAnalyTbl = [50][39]*formuPst{}
	for k, gram := range al.lang {
		for j := 1; j < 39; j++ {
			switch true {
			case gram.firstSet.where(j) >= 0:
				fp := &formuPst{
					g:   k,
					idx: gram.firstSet.where(j),
				}
				al.forecastAnalyTbl[gram.idx][j] = fp
			case gram.firstSet.where(lex.CodeMap["ε"]) > 0 && gram.followSet.where(j) >= 0:
				fp := &formuPst{
					g:   k,
					idx: gram.firstSet.where(lex.CodeMap["ε"]),
				}
				al.forecastAnalyTbl[gram.idx][j] = fp
			default:
			}
		}
	}

	// TODO: print
	al.PrintFATbl()
}

func (al *Analyser) isGramPass(lexAl *lex.Analyser) bool {
	stack := newAnalyStack(1 << 10)
	stack.pushOne(&gramChar{
		t:          lex.CodeMap["#"],
		isTerminal: true,
	})
	stack.pushOne(&gramChar{
		nt:         al.nonTerminalRoot,
		isTerminal: false,
	})

	for top := stack.pop(); top != nil; top = stack.pop() {
		if top.isTerminal {
			// fmt.Println("弹出：", top.t)
			if !lexAl.IsTokenMatch(top.t) {
				return false
			}
			continue
		}
		// fmt.Println("弹出：", top.nt)
		nonTermIdx := al.lang[top.nt].idx
		// fmt.Print("找预测分析表："+top.nt, nonTermIdx, "行")
		termMachineCode := lexAl.NextCode()
		// fmt.Println(termMachineCode, "列")
		position := al.forecastAnalyTbl[nonTermIdx][termMachineCode]
		if position == nil {
			fmt.Printf("\n[ERROR]: '%s' should not start with '%s' in line %d", top.nt, lex.ReversCodeMap[lexAl.NextCode()], lexAl.Tokens[lexAl.NextToken()].Row)
			return false
		}
		gram := al.lang[top.nt]
		if gram == nil {
			panic("gram == nil")
		}
		stack.push(gram.right[position.idx])
	}
	return true
}
