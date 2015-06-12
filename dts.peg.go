package main

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 1114112

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulestart
	ruleModuleDeclaration
	ruleModuleBody
	ruleClassDeclaration
	ruleClassBody
	ruleExtendClause
	ruleImplementClause
	ruleVariableDeclaration
	ruleFuncDeclaration
	ruleType
	ruleBasicType
	ruleFuncType
	ruleArrayType
	ruleObjectType
	ruleLiteralType
	ruleVariableDifinition
	ruleTypeSeperator
	ruleFuncReturn
	ruleArgumentSeperator
	ruleDeclarationSeperator
	rulekeywords
	rulemodifier
	ruleidentifier
	ruleseparator
	ruleComment
	ruleLineComment
	ruleBlockComment
	ruleblock_start
	ruleblock_end
	ruleparen_start
	ruleparen_end
	ruleSPACE
	rulews
	rulespacing
	ruleeol
	ruleeof

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"start",
	"ModuleDeclaration",
	"ModuleBody",
	"ClassDeclaration",
	"ClassBody",
	"ExtendClause",
	"ImplementClause",
	"VariableDeclaration",
	"FuncDeclaration",
	"Type",
	"BasicType",
	"FuncType",
	"ArrayType",
	"ObjectType",
	"LiteralType",
	"VariableDifinition",
	"TypeSeperator",
	"FuncReturn",
	"ArgumentSeperator",
	"DeclarationSeperator",
	"keywords",
	"modifier",
	"identifier",
	"separator",
	"Comment",
	"LineComment",
	"BlockComment",
	"block_start",
	"block_end",
	"paren_start",
	"paren_end",
	"SPACE",
	"ws",
	"spacing",
	"eol",
	"eof",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next uint32, depth int)
	Expand(index int) tokenTree
	Tokens() <-chan token32
	AST() *node32
	Error() []token32
	trim(length int)
}

type node32 struct {
	token32
	up, next *node32
}

func (node *node32) print(depth int, buffer string) {
	for node != nil {
		for c := 0; c < depth; c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(string(([]rune(buffer)[node.begin:node.end]))))
		if node.up != nil {
			node.up.print(depth+1, buffer)
		}
		node = node.next
	}
}

func (ast *node32) Print(buffer string) {
	ast.print(0, buffer)
}

type element struct {
	node *node32
	down *element
}

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next uint32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: uint32(t.begin), end: uint32(t.end), next: uint32(t.next)}
}

func (t *token32) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens32 struct {
	tree    []token32
	ordered [][]token32
}

func (t *tokens32) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens32) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens32) Order() [][]token32 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int32, 1, math.MaxInt16)
	for i, token := range t.tree {
		if token.pegRule == ruleUnknown {
			t.tree = t.tree[:i]
			break
		}
		depth := int(token.next)
		if length := len(depths); depth >= length {
			depths = depths[:depth+1]
		}
		depths[depth]++
	}
	depths = append(depths, 0)

	ordered, pool := make([][]token32, len(depths)), make([]token32, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = uint32(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state32 struct {
	token32
	depths []int32
	leaf   bool
}

func (t *tokens32) AST() *node32 {
	tokens := t.Tokens()
	stack := &element{node: &node32{token32: <-tokens}}
	for token := range tokens {
		if token.begin == token.end {
			continue
		}
		node := &node32{token32: token}
		for stack != nil && stack.node.begin >= token.begin && stack.node.end <= token.end {
			stack.node.next = node.up
			node.up = stack.node
			stack = stack.down
		}
		stack = &element{node: node, down: stack}
	}
	return stack.node
}

func (t *tokens32) PreOrder() (<-chan state32, [][]token32) {
	s, ordered := make(chan state32, 6), t.Order()
	go func() {
		var states [8]state32
		for i, _ := range states {
			states[i].depths = make([]int32, len(ordered))
		}
		depths, state, depth := make([]int32, len(ordered)), 0, 1
		write := func(t token32, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, uint32(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token32 = ordered[0][0]
		depths[0]++
		state++
		a, b := ordered[depth-1][depths[depth-1]-1], ordered[depth][depths[depth]]
	depthFirstSearch:
		for {
			for {
				if i := depths[depth]; i > 0 {
					if c, j := ordered[depth][i-1], depths[depth-1]; a.isParentOf(c) &&
						(j < 2 || !ordered[depth-1][j-2].isParentOf(c)) {
						if c.end != b.begin {
							write(token32{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token32{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
				}
				break
			}

			next := depth + 1
			if c := ordered[next][depths[next]]; c.pegRule != ruleUnknown && b.isParentOf(c) {
				write(b, false)
				depths[depth]++
				depth, a, b = next, b, c
				continue
			}

			write(b, true)
			depths[depth]++
			c, parent := ordered[depth][depths[depth]], true
			for {
				if c.pegRule != ruleUnknown && a.isParentOf(c) {
					b = c
					continue depthFirstSearch
				} else if parent && b.end != a.end {
					write(token32{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
				}

				depth--
				if depth > 0 {
					a, b, c = ordered[depth-1][depths[depth-1]-1], a, ordered[depth][depths[depth]]
					parent = a.isParentOf(b)
					continue
				}

				break depthFirstSearch
			}
		}

		close(s)
	}()
	return s, ordered
}

func (t *tokens32) PrintSyntax() {
	tokens, ordered := t.PreOrder()
	max := -1
	for token := range tokens {
		if !token.leaf {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[36m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[36m%v\x1B[m\n", rul3s[token.pegRule])
		} else if token.begin == token.end {
			fmt.Printf("%v", token.begin)
			for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
				fmt.Printf(" \x1B[31m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
			}
			fmt.Printf(" \x1B[31m%v\x1B[m\n", rul3s[token.pegRule])
		} else {
			for c, end := token.begin, token.end; c < end; c++ {
				if i := int(c); max+1 < i {
					for j := max; j < i; j++ {
						fmt.Printf("skip %v %v\n", j, token.String())
					}
					max = i
				} else if i := int(c); i <= max {
					for j := i; j <= max; j++ {
						fmt.Printf("dupe %v %v\n", j, token.String())
					}
				} else {
					max = int(c)
				}
				fmt.Printf("%v", c)
				for i, leaf, depths := 0, int(token.next), token.depths; i < leaf; i++ {
					fmt.Printf(" \x1B[34m%v\x1B[m", rul3s[ordered[i][depths[i]-1].pegRule])
				}
				fmt.Printf(" \x1B[34m%v\x1B[m\n", rul3s[token.pegRule])
			}
			fmt.Printf("\n")
		}
	}
}

func (t *tokens32) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(string(([]rune(buffer)[token.begin:token.end]))))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth uint32, index int) {
	t.tree[index] = token32{pegRule: rule, begin: uint32(begin), end: uint32(end), next: uint32(depth)}
}

func (t *tokens32) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens32) Error() []token32 {
	ordered := t.Order()
	length := len(ordered)
	tokens, length := make([]token32, length), length-1
	for i, _ := range tokens {
		o := ordered[length-i]
		if len(o) > 1 {
			tokens[i] = o[len(o)-2].getToken32()
		}
	}
	return tokens
}

/*func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2 * len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}*/

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type Parser struct {
	Buffer string
	buffer []rune
	rules  [37]func() bool
	Parse  func(rule ...int) error
	Reset  func()
	tokenTree
}

type textPosition struct {
	line, symbol int
}

type textPositionMap map[int]textPosition

func translatePositions(buffer string, positions []int) textPositionMap {
	length, translations, j, line, symbol := len(positions), make(textPositionMap, len(positions)), 0, 1, 0
	sort.Ints(positions)

search:
	for i, c := range buffer[0:] {
		if c == '\n' {
			line, symbol = line+1, 0
		} else {
			symbol++
		}
		if i == positions[j] {
			translations[positions[j]] = textPosition{line, symbol}
			for j++; j < length; j++ {
				if i != positions[j] {
					continue search
				}
			}
			break search
		}
	}

	return translations
}

type parseError struct {
	p *Parser
}

func (e *parseError) Error() string {
	tokens, error := e.p.tokenTree.Error(), "\n"
	positions, p := make([]int, 2*len(tokens)), 0
	for _, token := range tokens {
		positions[p], p = int(token.begin), p+1
		positions[p], p = int(token.end), p+1
	}
	translations := translatePositions(e.p.Buffer, positions)
	for _, token := range tokens {
		begin, end := int(token.begin), int(token.end)
		error += fmt.Sprintf("parse error near \x1B[34m%v\x1B[m (line %v symbol %v - line %v symbol %v):\n%v\n",
			rul3s[token.pegRule],
			translations[begin].line, translations[begin].symbol,
			translations[end].line, translations[end].symbol,
			/*strconv.Quote(*/ e.p.Buffer[begin:end] /*)*/)
	}

	return error
}

func (p *Parser) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *Parser) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *Parser) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens32{tree: make([]token32, math.MaxInt16)}
	position, depth, tokenIndex, buffer, _rules := uint32(0), uint32(0), 0, p.buffer, p.rules

	p.Parse = func(rule ...int) error {
		r := 1
		if len(rule) > 0 {
			r = rule[0]
		}
		matches := p.rules[r]()
		p.tokenTree = tree
		if matches {
			p.tokenTree.trim(tokenIndex)
			return nil
		}
		return &parseError{p}
	}

	p.Reset = func() {
		position, tokenIndex, depth = 0, 0, 0
	}

	add := func(rule pegRule, begin uint32) {
		if t := tree.Expand(tokenIndex); t != nil {
			tree = t
		}
		tree.Add(rule, begin, position, depth, tokenIndex)
		tokenIndex++
	}

	matchDot := func() bool {
		if buffer[position] != end_symbol {
			position++
			return true
		}
		return false
	}

	/*matchChar := func(c byte) bool {
		if buffer[position] == c {
			position++
			return true
		}
		return false
	}*/

	/*matchRange := func(lower byte, upper byte) bool {
		if c := buffer[position]; c >= lower && c <= upper {
			position++
			return true
		}
		return false
	}*/

	_rules = [...]func() bool{
		nil,
		/* 0 start <- <(SPACE (ModuleDeclaration / ClassDeclaration)+ eof)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !_rules[ruleSPACE]() {
					goto l0
				}
				{
					position4, tokenIndex4, depth4 := position, tokenIndex, depth
					if !_rules[ruleModuleDeclaration]() {
						goto l5
					}
					goto l4
				l5:
					position, tokenIndex, depth = position4, tokenIndex4, depth4
					if !_rules[ruleClassDeclaration]() {
						goto l0
					}
				}
			l4:
			l2:
				{
					position3, tokenIndex3, depth3 := position, tokenIndex, depth
					{
						position6, tokenIndex6, depth6 := position, tokenIndex, depth
						if !_rules[ruleModuleDeclaration]() {
							goto l7
						}
						goto l6
					l7:
						position, tokenIndex, depth = position6, tokenIndex6, depth6
						if !_rules[ruleClassDeclaration]() {
							goto l3
						}
					}
				l6:
					goto l2
				l3:
					position, tokenIndex, depth = position3, tokenIndex3, depth3
				}
				{
					position8 := position
					depth++
					{
						position9, tokenIndex9, depth9 := position, tokenIndex, depth
						if !matchDot() {
							goto l9
						}
						goto l0
					l9:
						position, tokenIndex, depth = position9, tokenIndex9, depth9
					}
					depth--
					add(ruleeof, position8)
				}
				depth--
				add(rulestart, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 ModuleDeclaration <- <(modifier? ('m' 'o' 'd' 'u' 'l' 'e') SPACE identifier ModuleBody)> */
		func() bool {
			position10, tokenIndex10, depth10 := position, tokenIndex, depth
			{
				position11 := position
				depth++
				{
					position12, tokenIndex12, depth12 := position, tokenIndex, depth
					if !_rules[rulemodifier]() {
						goto l12
					}
					goto l13
				l12:
					position, tokenIndex, depth = position12, tokenIndex12, depth12
				}
			l13:
				if buffer[position] != rune('m') {
					goto l10
				}
				position++
				if buffer[position] != rune('o') {
					goto l10
				}
				position++
				if buffer[position] != rune('d') {
					goto l10
				}
				position++
				if buffer[position] != rune('u') {
					goto l10
				}
				position++
				if buffer[position] != rune('l') {
					goto l10
				}
				position++
				if buffer[position] != rune('e') {
					goto l10
				}
				position++
				if !_rules[ruleSPACE]() {
					goto l10
				}
				if !_rules[ruleidentifier]() {
					goto l10
				}
				{
					position14 := position
					depth++
					if !_rules[ruleblock_start]() {
						goto l10
					}
				l15:
					{
						position16, tokenIndex16, depth16 := position, tokenIndex, depth
						{
							position17, tokenIndex17, depth17 := position, tokenIndex, depth
							if !_rules[ruleClassDeclaration]() {
								goto l18
							}
							goto l17
						l18:
							position, tokenIndex, depth = position17, tokenIndex17, depth17
							if !_rules[ruleModuleDeclaration]() {
								goto l19
							}
							goto l17
						l19:
							position, tokenIndex, depth = position17, tokenIndex17, depth17
							if !_rules[ruleVariableDeclaration]() {
								goto l20
							}
							goto l17
						l20:
							position, tokenIndex, depth = position17, tokenIndex17, depth17
							if !_rules[ruleFuncDeclaration]() {
								goto l16
							}
						}
					l17:
						goto l15
					l16:
						position, tokenIndex, depth = position16, tokenIndex16, depth16
					}
					if !_rules[ruleblock_end]() {
						goto l10
					}
					depth--
					add(ruleModuleBody, position14)
				}
				depth--
				add(ruleModuleDeclaration, position11)
			}
			return true
		l10:
			position, tokenIndex, depth = position10, tokenIndex10, depth10
			return false
		},
		/* 2 ModuleBody <- <(block_start (ClassDeclaration / ModuleDeclaration / VariableDeclaration / FuncDeclaration)* block_end)> */
		nil,
		/* 3 ClassDeclaration <- <(modifier? (('c' 'l' 'a' 's' 's') / ('i' 'n' 't' 'e' 'r' 'f' 'a' 'c' 'e')) SPACE identifier ExtendClause? ImplementClause? ClassBody)> */
		func() bool {
			position22, tokenIndex22, depth22 := position, tokenIndex, depth
			{
				position23 := position
				depth++
				{
					position24, tokenIndex24, depth24 := position, tokenIndex, depth
					if !_rules[rulemodifier]() {
						goto l24
					}
					goto l25
				l24:
					position, tokenIndex, depth = position24, tokenIndex24, depth24
				}
			l25:
				{
					position26, tokenIndex26, depth26 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l27
					}
					position++
					if buffer[position] != rune('l') {
						goto l27
					}
					position++
					if buffer[position] != rune('a') {
						goto l27
					}
					position++
					if buffer[position] != rune('s') {
						goto l27
					}
					position++
					if buffer[position] != rune('s') {
						goto l27
					}
					position++
					goto l26
				l27:
					position, tokenIndex, depth = position26, tokenIndex26, depth26
					if buffer[position] != rune('i') {
						goto l22
					}
					position++
					if buffer[position] != rune('n') {
						goto l22
					}
					position++
					if buffer[position] != rune('t') {
						goto l22
					}
					position++
					if buffer[position] != rune('e') {
						goto l22
					}
					position++
					if buffer[position] != rune('r') {
						goto l22
					}
					position++
					if buffer[position] != rune('f') {
						goto l22
					}
					position++
					if buffer[position] != rune('a') {
						goto l22
					}
					position++
					if buffer[position] != rune('c') {
						goto l22
					}
					position++
					if buffer[position] != rune('e') {
						goto l22
					}
					position++
				}
			l26:
				if !_rules[ruleSPACE]() {
					goto l22
				}
				if !_rules[ruleidentifier]() {
					goto l22
				}
				{
					position28, tokenIndex28, depth28 := position, tokenIndex, depth
					{
						position30 := position
						depth++
						if !_rules[ruleSPACE]() {
							goto l28
						}
						if buffer[position] != rune('e') {
							goto l28
						}
						position++
						if buffer[position] != rune('x') {
							goto l28
						}
						position++
						if buffer[position] != rune('t') {
							goto l28
						}
						position++
						if buffer[position] != rune('e') {
							goto l28
						}
						position++
						if buffer[position] != rune('n') {
							goto l28
						}
						position++
						if buffer[position] != rune('d') {
							goto l28
						}
						position++
						if buffer[position] != rune('s') {
							goto l28
						}
						position++
						if !_rules[ruleSPACE]() {
							goto l28
						}
						if !_rules[ruleidentifier]() {
							goto l28
						}
						depth--
						add(ruleExtendClause, position30)
					}
					goto l29
				l28:
					position, tokenIndex, depth = position28, tokenIndex28, depth28
				}
			l29:
				{
					position31, tokenIndex31, depth31 := position, tokenIndex, depth
					{
						position33 := position
						depth++
						if !_rules[ruleSPACE]() {
							goto l31
						}
						if buffer[position] != rune('i') {
							goto l31
						}
						position++
						if buffer[position] != rune('m') {
							goto l31
						}
						position++
						if buffer[position] != rune('p') {
							goto l31
						}
						position++
						if buffer[position] != rune('l') {
							goto l31
						}
						position++
						if buffer[position] != rune('e') {
							goto l31
						}
						position++
						if buffer[position] != rune('m') {
							goto l31
						}
						position++
						if buffer[position] != rune('e') {
							goto l31
						}
						position++
						if buffer[position] != rune('n') {
							goto l31
						}
						position++
						if buffer[position] != rune('t') {
							goto l31
						}
						position++
						if buffer[position] != rune('s') {
							goto l31
						}
						position++
						if !_rules[ruleSPACE]() {
							goto l31
						}
						if !_rules[ruleidentifier]() {
							goto l31
						}
						depth--
						add(ruleImplementClause, position33)
					}
					goto l32
				l31:
					position, tokenIndex, depth = position31, tokenIndex31, depth31
				}
			l32:
				{
					position34 := position
					depth++
					if !_rules[ruleblock_start]() {
						goto l22
					}
				l35:
					{
						position36, tokenIndex36, depth36 := position, tokenIndex, depth
						{
							position37, tokenIndex37, depth37 := position, tokenIndex, depth
							if !_rules[ruleVariableDeclaration]() {
								goto l38
							}
							goto l37
						l38:
							position, tokenIndex, depth = position37, tokenIndex37, depth37
							if !_rules[ruleFuncDeclaration]() {
								goto l36
							}
						}
					l37:
						goto l35
					l36:
						position, tokenIndex, depth = position36, tokenIndex36, depth36
					}
					if !_rules[ruleblock_end]() {
						goto l22
					}
					depth--
					add(ruleClassBody, position34)
				}
				depth--
				add(ruleClassDeclaration, position23)
			}
			return true
		l22:
			position, tokenIndex, depth = position22, tokenIndex22, depth22
			return false
		},
		/* 4 ClassBody <- <(block_start (VariableDeclaration / FuncDeclaration)* block_end)> */
		nil,
		/* 5 ExtendClause <- <(SPACE ('e' 'x' 't' 'e' 'n' 'd' 's') SPACE identifier)> */
		nil,
		/* 6 ImplementClause <- <(SPACE ('i' 'm' 'p' 'l' 'e' 'm' 'e' 'n' 't' 's') SPACE identifier)> */
		nil,
		/* 7 VariableDeclaration <- <(modifier? VariableDifinition DeclarationSeperator)> */
		func() bool {
			position42, tokenIndex42, depth42 := position, tokenIndex, depth
			{
				position43 := position
				depth++
				{
					position44, tokenIndex44, depth44 := position, tokenIndex, depth
					if !_rules[rulemodifier]() {
						goto l44
					}
					goto l45
				l44:
					position, tokenIndex, depth = position44, tokenIndex44, depth44
				}
			l45:
				if !_rules[ruleVariableDifinition]() {
					goto l42
				}
				if !_rules[ruleDeclarationSeperator]() {
					goto l42
				}
				depth--
				add(ruleVariableDeclaration, position43)
			}
			return true
		l42:
			position, tokenIndex, depth = position42, tokenIndex42, depth42
			return false
		},
		/* 8 FuncDeclaration <- <(modifier? identifier FuncType DeclarationSeperator)> */
		func() bool {
			position46, tokenIndex46, depth46 := position, tokenIndex, depth
			{
				position47 := position
				depth++
				{
					position48, tokenIndex48, depth48 := position, tokenIndex, depth
					if !_rules[rulemodifier]() {
						goto l48
					}
					goto l49
				l48:
					position, tokenIndex, depth = position48, tokenIndex48, depth48
				}
			l49:
				if !_rules[ruleidentifier]() {
					goto l46
				}
				if !_rules[ruleFuncType]() {
					goto l46
				}
				if !_rules[ruleDeclarationSeperator]() {
					goto l46
				}
				depth--
				add(ruleFuncDeclaration, position47)
			}
			return true
		l46:
			position, tokenIndex, depth = position46, tokenIndex46, depth46
			return false
		},
		/* 9 Type <- <(BasicType (SPACE '|' SPACE BasicType)*)> */
		func() bool {
			position50, tokenIndex50, depth50 := position, tokenIndex, depth
			{
				position51 := position
				depth++
				if !_rules[ruleBasicType]() {
					goto l50
				}
			l52:
				{
					position53, tokenIndex53, depth53 := position, tokenIndex, depth
					if !_rules[ruleSPACE]() {
						goto l53
					}
					if buffer[position] != rune('|') {
						goto l53
					}
					position++
					if !_rules[ruleSPACE]() {
						goto l53
					}
					if !_rules[ruleBasicType]() {
						goto l53
					}
					goto l52
				l53:
					position, tokenIndex, depth = position53, tokenIndex53, depth53
				}
				depth--
				add(ruleType, position51)
			}
			return true
		l50:
			position, tokenIndex, depth = position50, tokenIndex50, depth50
			return false
		},
		/* 10 BasicType <- <(ObjectType / ArrayType / (('n' / 'N') ('u' / 'U') ('m' / 'M') ('b' / 'B') ('e' / 'E') ('r' / 'R')) / (('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L') ('e' / 'E') ('a' / 'A') ('n' / 'N')) / (('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G')) / (('f' / 'F') ('u' / 'U') ('n' / 'N') ('c' / 'C') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N')) / (('a' / 'A') ('n' / 'N') ('y' / 'Y')) / ((&('\'') LiteralType) | (&('\t' | '\n' | '\r' | ' ' | '(' | '/') FuncType) | (&('.' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') identifier)))> */
		func() bool {
			position54, tokenIndex54, depth54 := position, tokenIndex, depth
			{
				position55 := position
				depth++
				{
					position56, tokenIndex56, depth56 := position, tokenIndex, depth
					{
						position58 := position
						depth++
						if !_rules[ruleblock_start]() {
							goto l57
						}
					l59:
						{
							position60, tokenIndex60, depth60 := position, tokenIndex, depth
							if !_rules[ruleVariableDeclaration]() {
								goto l60
							}
							goto l59
						l60:
							position, tokenIndex, depth = position60, tokenIndex60, depth60
						}
						if !_rules[ruleblock_end]() {
							goto l57
						}
						depth--
						add(ruleObjectType, position58)
					}
					goto l56
				l57:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					{
						position62 := position
						depth++
						{
							position63, tokenIndex63, depth63 := position, tokenIndex, depth
							{
								position65, tokenIndex65, depth65 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l66
								}
								position++
								goto l65
							l66:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
								if buffer[position] != rune('N') {
									goto l64
								}
								position++
							}
						l65:
							{
								position67, tokenIndex67, depth67 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l68
								}
								position++
								goto l67
							l68:
								position, tokenIndex, depth = position67, tokenIndex67, depth67
								if buffer[position] != rune('U') {
									goto l64
								}
								position++
							}
						l67:
							{
								position69, tokenIndex69, depth69 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l70
								}
								position++
								goto l69
							l70:
								position, tokenIndex, depth = position69, tokenIndex69, depth69
								if buffer[position] != rune('M') {
									goto l64
								}
								position++
							}
						l69:
							{
								position71, tokenIndex71, depth71 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l72
								}
								position++
								goto l71
							l72:
								position, tokenIndex, depth = position71, tokenIndex71, depth71
								if buffer[position] != rune('B') {
									goto l64
								}
								position++
							}
						l71:
							{
								position73, tokenIndex73, depth73 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l74
								}
								position++
								goto l73
							l74:
								position, tokenIndex, depth = position73, tokenIndex73, depth73
								if buffer[position] != rune('E') {
									goto l64
								}
								position++
							}
						l73:
							{
								position75, tokenIndex75, depth75 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l76
								}
								position++
								goto l75
							l76:
								position, tokenIndex, depth = position75, tokenIndex75, depth75
								if buffer[position] != rune('R') {
									goto l64
								}
								position++
							}
						l75:
							goto l63
						l64:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
							{
								position78, tokenIndex78, depth78 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l79
								}
								position++
								goto l78
							l79:
								position, tokenIndex, depth = position78, tokenIndex78, depth78
								if buffer[position] != rune('B') {
									goto l77
								}
								position++
							}
						l78:
							{
								position80, tokenIndex80, depth80 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l81
								}
								position++
								goto l80
							l81:
								position, tokenIndex, depth = position80, tokenIndex80, depth80
								if buffer[position] != rune('O') {
									goto l77
								}
								position++
							}
						l80:
							{
								position82, tokenIndex82, depth82 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l83
								}
								position++
								goto l82
							l83:
								position, tokenIndex, depth = position82, tokenIndex82, depth82
								if buffer[position] != rune('O') {
									goto l77
								}
								position++
							}
						l82:
							{
								position84, tokenIndex84, depth84 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l85
								}
								position++
								goto l84
							l85:
								position, tokenIndex, depth = position84, tokenIndex84, depth84
								if buffer[position] != rune('L') {
									goto l77
								}
								position++
							}
						l84:
							{
								position86, tokenIndex86, depth86 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l87
								}
								position++
								goto l86
							l87:
								position, tokenIndex, depth = position86, tokenIndex86, depth86
								if buffer[position] != rune('E') {
									goto l77
								}
								position++
							}
						l86:
							{
								position88, tokenIndex88, depth88 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l89
								}
								position++
								goto l88
							l89:
								position, tokenIndex, depth = position88, tokenIndex88, depth88
								if buffer[position] != rune('A') {
									goto l77
								}
								position++
							}
						l88:
							{
								position90, tokenIndex90, depth90 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l91
								}
								position++
								goto l90
							l91:
								position, tokenIndex, depth = position90, tokenIndex90, depth90
								if buffer[position] != rune('N') {
									goto l77
								}
								position++
							}
						l90:
							goto l63
						l77:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
							{
								position93, tokenIndex93, depth93 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l94
								}
								position++
								goto l93
							l94:
								position, tokenIndex, depth = position93, tokenIndex93, depth93
								if buffer[position] != rune('S') {
									goto l92
								}
								position++
							}
						l93:
							{
								position95, tokenIndex95, depth95 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l96
								}
								position++
								goto l95
							l96:
								position, tokenIndex, depth = position95, tokenIndex95, depth95
								if buffer[position] != rune('T') {
									goto l92
								}
								position++
							}
						l95:
							{
								position97, tokenIndex97, depth97 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l98
								}
								position++
								goto l97
							l98:
								position, tokenIndex, depth = position97, tokenIndex97, depth97
								if buffer[position] != rune('R') {
									goto l92
								}
								position++
							}
						l97:
							{
								position99, tokenIndex99, depth99 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l100
								}
								position++
								goto l99
							l100:
								position, tokenIndex, depth = position99, tokenIndex99, depth99
								if buffer[position] != rune('I') {
									goto l92
								}
								position++
							}
						l99:
							{
								position101, tokenIndex101, depth101 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l102
								}
								position++
								goto l101
							l102:
								position, tokenIndex, depth = position101, tokenIndex101, depth101
								if buffer[position] != rune('N') {
									goto l92
								}
								position++
							}
						l101:
							{
								position103, tokenIndex103, depth103 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l104
								}
								position++
								goto l103
							l104:
								position, tokenIndex, depth = position103, tokenIndex103, depth103
								if buffer[position] != rune('G') {
									goto l92
								}
								position++
							}
						l103:
							goto l63
						l92:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
							{
								position106, tokenIndex106, depth106 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l107
								}
								position++
								goto l106
							l107:
								position, tokenIndex, depth = position106, tokenIndex106, depth106
								if buffer[position] != rune('F') {
									goto l105
								}
								position++
							}
						l106:
							{
								position108, tokenIndex108, depth108 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l109
								}
								position++
								goto l108
							l109:
								position, tokenIndex, depth = position108, tokenIndex108, depth108
								if buffer[position] != rune('U') {
									goto l105
								}
								position++
							}
						l108:
							{
								position110, tokenIndex110, depth110 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l111
								}
								position++
								goto l110
							l111:
								position, tokenIndex, depth = position110, tokenIndex110, depth110
								if buffer[position] != rune('N') {
									goto l105
								}
								position++
							}
						l110:
							{
								position112, tokenIndex112, depth112 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l113
								}
								position++
								goto l112
							l113:
								position, tokenIndex, depth = position112, tokenIndex112, depth112
								if buffer[position] != rune('C') {
									goto l105
								}
								position++
							}
						l112:
							{
								position114, tokenIndex114, depth114 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l115
								}
								position++
								goto l114
							l115:
								position, tokenIndex, depth = position114, tokenIndex114, depth114
								if buffer[position] != rune('T') {
									goto l105
								}
								position++
							}
						l114:
							{
								position116, tokenIndex116, depth116 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l117
								}
								position++
								goto l116
							l117:
								position, tokenIndex, depth = position116, tokenIndex116, depth116
								if buffer[position] != rune('I') {
									goto l105
								}
								position++
							}
						l116:
							{
								position118, tokenIndex118, depth118 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l119
								}
								position++
								goto l118
							l119:
								position, tokenIndex, depth = position118, tokenIndex118, depth118
								if buffer[position] != rune('O') {
									goto l105
								}
								position++
							}
						l118:
							{
								position120, tokenIndex120, depth120 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l121
								}
								position++
								goto l120
							l121:
								position, tokenIndex, depth = position120, tokenIndex120, depth120
								if buffer[position] != rune('N') {
									goto l105
								}
								position++
							}
						l120:
							goto l63
						l105:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
							{
								position123, tokenIndex123, depth123 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l124
								}
								position++
								goto l123
							l124:
								position, tokenIndex, depth = position123, tokenIndex123, depth123
								if buffer[position] != rune('F') {
									goto l122
								}
								position++
							}
						l123:
							{
								position125, tokenIndex125, depth125 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l126
								}
								position++
								goto l125
							l126:
								position, tokenIndex, depth = position125, tokenIndex125, depth125
								if buffer[position] != rune('U') {
									goto l122
								}
								position++
							}
						l125:
							{
								position127, tokenIndex127, depth127 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l128
								}
								position++
								goto l127
							l128:
								position, tokenIndex, depth = position127, tokenIndex127, depth127
								if buffer[position] != rune('N') {
									goto l122
								}
								position++
							}
						l127:
							{
								position129, tokenIndex129, depth129 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l130
								}
								position++
								goto l129
							l130:
								position, tokenIndex, depth = position129, tokenIndex129, depth129
								if buffer[position] != rune('C') {
									goto l122
								}
								position++
							}
						l129:
							{
								position131, tokenIndex131, depth131 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l132
								}
								position++
								goto l131
							l132:
								position, tokenIndex, depth = position131, tokenIndex131, depth131
								if buffer[position] != rune('T') {
									goto l122
								}
								position++
							}
						l131:
							{
								position133, tokenIndex133, depth133 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l134
								}
								position++
								goto l133
							l134:
								position, tokenIndex, depth = position133, tokenIndex133, depth133
								if buffer[position] != rune('I') {
									goto l122
								}
								position++
							}
						l133:
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l136
								}
								position++
								goto l135
							l136:
								position, tokenIndex, depth = position135, tokenIndex135, depth135
								if buffer[position] != rune('O') {
									goto l122
								}
								position++
							}
						l135:
							{
								position137, tokenIndex137, depth137 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l138
								}
								position++
								goto l137
							l138:
								position, tokenIndex, depth = position137, tokenIndex137, depth137
								if buffer[position] != rune('N') {
									goto l122
								}
								position++
							}
						l137:
							goto l63
						l122:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
							{
								position140, tokenIndex140, depth140 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l141
								}
								position++
								goto l140
							l141:
								position, tokenIndex, depth = position140, tokenIndex140, depth140
								if buffer[position] != rune('A') {
									goto l139
								}
								position++
							}
						l140:
							{
								position142, tokenIndex142, depth142 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l143
								}
								position++
								goto l142
							l143:
								position, tokenIndex, depth = position142, tokenIndex142, depth142
								if buffer[position] != rune('N') {
									goto l139
								}
								position++
							}
						l142:
							{
								position144, tokenIndex144, depth144 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l145
								}
								position++
								goto l144
							l145:
								position, tokenIndex, depth = position144, tokenIndex144, depth144
								if buffer[position] != rune('Y') {
									goto l139
								}
								position++
							}
						l144:
							goto l63
						l139:
							position, tokenIndex, depth = position63, tokenIndex63, depth63
							if !_rules[ruleidentifier]() {
								goto l61
							}
						}
					l63:
						if buffer[position] != rune('[') {
							goto l61
						}
						position++
						if buffer[position] != rune(']') {
							goto l61
						}
						position++
						depth--
						add(ruleArrayType, position62)
					}
					goto l56
				l61:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					{
						position147, tokenIndex147, depth147 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l148
						}
						position++
						goto l147
					l148:
						position, tokenIndex, depth = position147, tokenIndex147, depth147
						if buffer[position] != rune('N') {
							goto l146
						}
						position++
					}
				l147:
					{
						position149, tokenIndex149, depth149 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l150
						}
						position++
						goto l149
					l150:
						position, tokenIndex, depth = position149, tokenIndex149, depth149
						if buffer[position] != rune('U') {
							goto l146
						}
						position++
					}
				l149:
					{
						position151, tokenIndex151, depth151 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l152
						}
						position++
						goto l151
					l152:
						position, tokenIndex, depth = position151, tokenIndex151, depth151
						if buffer[position] != rune('M') {
							goto l146
						}
						position++
					}
				l151:
					{
						position153, tokenIndex153, depth153 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l154
						}
						position++
						goto l153
					l154:
						position, tokenIndex, depth = position153, tokenIndex153, depth153
						if buffer[position] != rune('B') {
							goto l146
						}
						position++
					}
				l153:
					{
						position155, tokenIndex155, depth155 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l156
						}
						position++
						goto l155
					l156:
						position, tokenIndex, depth = position155, tokenIndex155, depth155
						if buffer[position] != rune('E') {
							goto l146
						}
						position++
					}
				l155:
					{
						position157, tokenIndex157, depth157 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l158
						}
						position++
						goto l157
					l158:
						position, tokenIndex, depth = position157, tokenIndex157, depth157
						if buffer[position] != rune('R') {
							goto l146
						}
						position++
					}
				l157:
					goto l56
				l146:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					{
						position160, tokenIndex160, depth160 := position, tokenIndex, depth
						if buffer[position] != rune('b') {
							goto l161
						}
						position++
						goto l160
					l161:
						position, tokenIndex, depth = position160, tokenIndex160, depth160
						if buffer[position] != rune('B') {
							goto l159
						}
						position++
					}
				l160:
					{
						position162, tokenIndex162, depth162 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l163
						}
						position++
						goto l162
					l163:
						position, tokenIndex, depth = position162, tokenIndex162, depth162
						if buffer[position] != rune('O') {
							goto l159
						}
						position++
					}
				l162:
					{
						position164, tokenIndex164, depth164 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l165
						}
						position++
						goto l164
					l165:
						position, tokenIndex, depth = position164, tokenIndex164, depth164
						if buffer[position] != rune('O') {
							goto l159
						}
						position++
					}
				l164:
					{
						position166, tokenIndex166, depth166 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l167
						}
						position++
						goto l166
					l167:
						position, tokenIndex, depth = position166, tokenIndex166, depth166
						if buffer[position] != rune('L') {
							goto l159
						}
						position++
					}
				l166:
					{
						position168, tokenIndex168, depth168 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l169
						}
						position++
						goto l168
					l169:
						position, tokenIndex, depth = position168, tokenIndex168, depth168
						if buffer[position] != rune('E') {
							goto l159
						}
						position++
					}
				l168:
					{
						position170, tokenIndex170, depth170 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l171
						}
						position++
						goto l170
					l171:
						position, tokenIndex, depth = position170, tokenIndex170, depth170
						if buffer[position] != rune('A') {
							goto l159
						}
						position++
					}
				l170:
					{
						position172, tokenIndex172, depth172 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l173
						}
						position++
						goto l172
					l173:
						position, tokenIndex, depth = position172, tokenIndex172, depth172
						if buffer[position] != rune('N') {
							goto l159
						}
						position++
					}
				l172:
					goto l56
				l159:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					{
						position175, tokenIndex175, depth175 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l176
						}
						position++
						goto l175
					l176:
						position, tokenIndex, depth = position175, tokenIndex175, depth175
						if buffer[position] != rune('S') {
							goto l174
						}
						position++
					}
				l175:
					{
						position177, tokenIndex177, depth177 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l178
						}
						position++
						goto l177
					l178:
						position, tokenIndex, depth = position177, tokenIndex177, depth177
						if buffer[position] != rune('T') {
							goto l174
						}
						position++
					}
				l177:
					{
						position179, tokenIndex179, depth179 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l180
						}
						position++
						goto l179
					l180:
						position, tokenIndex, depth = position179, tokenIndex179, depth179
						if buffer[position] != rune('R') {
							goto l174
						}
						position++
					}
				l179:
					{
						position181, tokenIndex181, depth181 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l182
						}
						position++
						goto l181
					l182:
						position, tokenIndex, depth = position181, tokenIndex181, depth181
						if buffer[position] != rune('I') {
							goto l174
						}
						position++
					}
				l181:
					{
						position183, tokenIndex183, depth183 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l184
						}
						position++
						goto l183
					l184:
						position, tokenIndex, depth = position183, tokenIndex183, depth183
						if buffer[position] != rune('N') {
							goto l174
						}
						position++
					}
				l183:
					{
						position185, tokenIndex185, depth185 := position, tokenIndex, depth
						if buffer[position] != rune('g') {
							goto l186
						}
						position++
						goto l185
					l186:
						position, tokenIndex, depth = position185, tokenIndex185, depth185
						if buffer[position] != rune('G') {
							goto l174
						}
						position++
					}
				l185:
					goto l56
				l174:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					{
						position188, tokenIndex188, depth188 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l189
						}
						position++
						goto l188
					l189:
						position, tokenIndex, depth = position188, tokenIndex188, depth188
						if buffer[position] != rune('F') {
							goto l187
						}
						position++
					}
				l188:
					{
						position190, tokenIndex190, depth190 := position, tokenIndex, depth
						if buffer[position] != rune('u') {
							goto l191
						}
						position++
						goto l190
					l191:
						position, tokenIndex, depth = position190, tokenIndex190, depth190
						if buffer[position] != rune('U') {
							goto l187
						}
						position++
					}
				l190:
					{
						position192, tokenIndex192, depth192 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l193
						}
						position++
						goto l192
					l193:
						position, tokenIndex, depth = position192, tokenIndex192, depth192
						if buffer[position] != rune('N') {
							goto l187
						}
						position++
					}
				l192:
					{
						position194, tokenIndex194, depth194 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l195
						}
						position++
						goto l194
					l195:
						position, tokenIndex, depth = position194, tokenIndex194, depth194
						if buffer[position] != rune('C') {
							goto l187
						}
						position++
					}
				l194:
					{
						position196, tokenIndex196, depth196 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l197
						}
						position++
						goto l196
					l197:
						position, tokenIndex, depth = position196, tokenIndex196, depth196
						if buffer[position] != rune('T') {
							goto l187
						}
						position++
					}
				l196:
					{
						position198, tokenIndex198, depth198 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l199
						}
						position++
						goto l198
					l199:
						position, tokenIndex, depth = position198, tokenIndex198, depth198
						if buffer[position] != rune('I') {
							goto l187
						}
						position++
					}
				l198:
					{
						position200, tokenIndex200, depth200 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l201
						}
						position++
						goto l200
					l201:
						position, tokenIndex, depth = position200, tokenIndex200, depth200
						if buffer[position] != rune('O') {
							goto l187
						}
						position++
					}
				l200:
					{
						position202, tokenIndex202, depth202 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l203
						}
						position++
						goto l202
					l203:
						position, tokenIndex, depth = position202, tokenIndex202, depth202
						if buffer[position] != rune('N') {
							goto l187
						}
						position++
					}
				l202:
					goto l56
				l187:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					{
						position205, tokenIndex205, depth205 := position, tokenIndex, depth
						if buffer[position] != rune('a') {
							goto l206
						}
						position++
						goto l205
					l206:
						position, tokenIndex, depth = position205, tokenIndex205, depth205
						if buffer[position] != rune('A') {
							goto l204
						}
						position++
					}
				l205:
					{
						position207, tokenIndex207, depth207 := position, tokenIndex, depth
						if buffer[position] != rune('n') {
							goto l208
						}
						position++
						goto l207
					l208:
						position, tokenIndex, depth = position207, tokenIndex207, depth207
						if buffer[position] != rune('N') {
							goto l204
						}
						position++
					}
				l207:
					{
						position209, tokenIndex209, depth209 := position, tokenIndex, depth
						if buffer[position] != rune('y') {
							goto l210
						}
						position++
						goto l209
					l210:
						position, tokenIndex, depth = position209, tokenIndex209, depth209
						if buffer[position] != rune('Y') {
							goto l204
						}
						position++
					}
				l209:
					goto l56
				l204:
					position, tokenIndex, depth = position56, tokenIndex56, depth56
					{
						switch buffer[position] {
						case '\'':
							{
								position212 := position
								depth++
								if buffer[position] != rune('\'') {
									goto l54
								}
								position++
								if !_rules[ruleidentifier]() {
									goto l54
								}
								if buffer[position] != rune('\'') {
									goto l54
								}
								position++
								depth--
								add(ruleLiteralType, position212)
							}
							break
						case '\t', '\n', '\r', ' ', '(', '/':
							if !_rules[ruleFuncType]() {
								goto l54
							}
							break
						default:
							if !_rules[ruleidentifier]() {
								goto l54
							}
							break
						}
					}

				}
			l56:
				depth--
				add(ruleBasicType, position55)
			}
			return true
		l54:
			position, tokenIndex, depth = position54, tokenIndex54, depth54
			return false
		},
		/* 11 FuncType <- <(paren_start VariableDifinition? (ArgumentSeperator VariableDifinition)* paren_end ((FuncReturn / TypeSeperator) Type)?)> */
		func() bool {
			position213, tokenIndex213, depth213 := position, tokenIndex, depth
			{
				position214 := position
				depth++
				{
					position215 := position
					depth++
					if !_rules[ruleSPACE]() {
						goto l213
					}
					if buffer[position] != rune('(') {
						goto l213
					}
					position++
					if !_rules[ruleSPACE]() {
						goto l213
					}
					depth--
					add(ruleparen_start, position215)
				}
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					if !_rules[ruleVariableDifinition]() {
						goto l216
					}
					goto l217
				l216:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
				}
			l217:
			l218:
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					{
						position220 := position
						depth++
						if !_rules[ruleSPACE]() {
							goto l219
						}
						if buffer[position] != rune(',') {
							goto l219
						}
						position++
						if !_rules[ruleSPACE]() {
							goto l219
						}
						depth--
						add(ruleArgumentSeperator, position220)
					}
					if !_rules[ruleVariableDifinition]() {
						goto l219
					}
					goto l218
				l219:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
				}
				{
					position221 := position
					depth++
					if !_rules[ruleSPACE]() {
						goto l213
					}
					if buffer[position] != rune(')') {
						goto l213
					}
					position++
					if !_rules[ruleSPACE]() {
						goto l213
					}
					depth--
					add(ruleparen_end, position221)
				}
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					{
						position224, tokenIndex224, depth224 := position, tokenIndex, depth
						{
							position226 := position
							depth++
							if !_rules[ruleSPACE]() {
								goto l225
							}
							if buffer[position] != rune('=') {
								goto l225
							}
							position++
							if buffer[position] != rune('>') {
								goto l225
							}
							position++
							if !_rules[ruleSPACE]() {
								goto l225
							}
							depth--
							add(ruleFuncReturn, position226)
						}
						goto l224
					l225:
						position, tokenIndex, depth = position224, tokenIndex224, depth224
						if !_rules[ruleTypeSeperator]() {
							goto l222
						}
					}
				l224:
					if !_rules[ruleType]() {
						goto l222
					}
					goto l223
				l222:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
				}
			l223:
				depth--
				add(ruleFuncType, position214)
			}
			return true
		l213:
			position, tokenIndex, depth = position213, tokenIndex213, depth213
			return false
		},
		/* 12 ArrayType <- <(((('n' / 'N') ('u' / 'U') ('m' / 'M') ('b' / 'B') ('e' / 'E') ('r' / 'R')) / (('b' / 'B') ('o' / 'O') ('o' / 'O') ('l' / 'L') ('e' / 'E') ('a' / 'A') ('n' / 'N')) / (('s' / 'S') ('t' / 'T') ('r' / 'R') ('i' / 'I') ('n' / 'N') ('g' / 'G')) / (('f' / 'F') ('u' / 'U') ('n' / 'N') ('c' / 'C') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N')) / (('f' / 'F') ('u' / 'U') ('n' / 'N') ('c' / 'C') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N')) / (('a' / 'A') ('n' / 'N') ('y' / 'Y')) / identifier) ('[' ']'))> */
		nil,
		/* 13 ObjectType <- <(block_start VariableDeclaration* block_end)> */
		nil,
		/* 14 LiteralType <- <('\'' identifier '\'')> */
		nil,
		/* 15 VariableDifinition <- <(identifier TypeSeperator Type)> */
		func() bool {
			position230, tokenIndex230, depth230 := position, tokenIndex, depth
			{
				position231 := position
				depth++
				if !_rules[ruleidentifier]() {
					goto l230
				}
				if !_rules[ruleTypeSeperator]() {
					goto l230
				}
				if !_rules[ruleType]() {
					goto l230
				}
				depth--
				add(ruleVariableDifinition, position231)
			}
			return true
		l230:
			position, tokenIndex, depth = position230, tokenIndex230, depth230
			return false
		},
		/* 16 TypeSeperator <- <(SPACE ':' SPACE)> */
		func() bool {
			position232, tokenIndex232, depth232 := position, tokenIndex, depth
			{
				position233 := position
				depth++
				if !_rules[ruleSPACE]() {
					goto l232
				}
				if buffer[position] != rune(':') {
					goto l232
				}
				position++
				if !_rules[ruleSPACE]() {
					goto l232
				}
				depth--
				add(ruleTypeSeperator, position233)
			}
			return true
		l232:
			position, tokenIndex, depth = position232, tokenIndex232, depth232
			return false
		},
		/* 17 FuncReturn <- <(SPACE ('=' '>') SPACE)> */
		nil,
		/* 18 ArgumentSeperator <- <(SPACE ',' SPACE)> */
		nil,
		/* 19 DeclarationSeperator <- <((';' / eol)? SPACE)?> */
		func() bool {
			{
				position237 := position
				depth++
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					{
						position240, tokenIndex240, depth240 := position, tokenIndex, depth
						{
							position242, tokenIndex242, depth242 := position, tokenIndex, depth
							if buffer[position] != rune(';') {
								goto l243
							}
							position++
							goto l242
						l243:
							position, tokenIndex, depth = position242, tokenIndex242, depth242
							if !_rules[ruleeol]() {
								goto l240
							}
						}
					l242:
						goto l241
					l240:
						position, tokenIndex, depth = position240, tokenIndex240, depth240
					}
				l241:
					if !_rules[ruleSPACE]() {
						goto l238
					}
					goto l239
				l238:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
				}
			l239:
				depth--
				add(ruleDeclarationSeperator, position237)
			}
			return true
		},
		/* 20 keywords <- <(('m' 'o' 'd' 'u' 'l' 'e') / ('c' 'l' 'a' 's' 's') / ('i' 'n' 't' 'e' 'r' 'f' 'a' 'c' 'e') / ('e' 'x' 't' 'e' 'n' 'd' 's') / ('i' 'm' 'p' 'l' 'e' 'm' 'e' 'n' 't' 's') / ('b' 'o' 'o' 'l' 'e' 'a' 'n') / ('n' 'u' 'm' 'b' 'e' 'r') / ('s' 't' 'r' 'i' 'n' 'g') / ('v' 'o' 'i' 'd'))> */
		nil,
		/* 21 modifier <- <(((&('s') ('s' 't' 'a' 't' 'i' 'c')) | (&('p') ('p' 'r' 'i' 'v' 'a' 't' 'e')) | (&('e') ('e' 'x' 'p' 'o' 'r' 't')) | (&('d') ('d' 'e' 'c' 'l' 'a' 'r' 'e'))) SPACE)+> */
		func() bool {
			position245, tokenIndex245, depth245 := position, tokenIndex, depth
			{
				position246 := position
				depth++
				{
					switch buffer[position] {
					case 's':
						if buffer[position] != rune('s') {
							goto l245
						}
						position++
						if buffer[position] != rune('t') {
							goto l245
						}
						position++
						if buffer[position] != rune('a') {
							goto l245
						}
						position++
						if buffer[position] != rune('t') {
							goto l245
						}
						position++
						if buffer[position] != rune('i') {
							goto l245
						}
						position++
						if buffer[position] != rune('c') {
							goto l245
						}
						position++
						break
					case 'p':
						if buffer[position] != rune('p') {
							goto l245
						}
						position++
						if buffer[position] != rune('r') {
							goto l245
						}
						position++
						if buffer[position] != rune('i') {
							goto l245
						}
						position++
						if buffer[position] != rune('v') {
							goto l245
						}
						position++
						if buffer[position] != rune('a') {
							goto l245
						}
						position++
						if buffer[position] != rune('t') {
							goto l245
						}
						position++
						if buffer[position] != rune('e') {
							goto l245
						}
						position++
						break
					case 'e':
						if buffer[position] != rune('e') {
							goto l245
						}
						position++
						if buffer[position] != rune('x') {
							goto l245
						}
						position++
						if buffer[position] != rune('p') {
							goto l245
						}
						position++
						if buffer[position] != rune('o') {
							goto l245
						}
						position++
						if buffer[position] != rune('r') {
							goto l245
						}
						position++
						if buffer[position] != rune('t') {
							goto l245
						}
						position++
						break
					default:
						if buffer[position] != rune('d') {
							goto l245
						}
						position++
						if buffer[position] != rune('e') {
							goto l245
						}
						position++
						if buffer[position] != rune('c') {
							goto l245
						}
						position++
						if buffer[position] != rune('l') {
							goto l245
						}
						position++
						if buffer[position] != rune('a') {
							goto l245
						}
						position++
						if buffer[position] != rune('r') {
							goto l245
						}
						position++
						if buffer[position] != rune('e') {
							goto l245
						}
						position++
						break
					}
				}

				if !_rules[ruleSPACE]() {
					goto l245
				}
			l247:
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case 's':
							if buffer[position] != rune('s') {
								goto l248
							}
							position++
							if buffer[position] != rune('t') {
								goto l248
							}
							position++
							if buffer[position] != rune('a') {
								goto l248
							}
							position++
							if buffer[position] != rune('t') {
								goto l248
							}
							position++
							if buffer[position] != rune('i') {
								goto l248
							}
							position++
							if buffer[position] != rune('c') {
								goto l248
							}
							position++
							break
						case 'p':
							if buffer[position] != rune('p') {
								goto l248
							}
							position++
							if buffer[position] != rune('r') {
								goto l248
							}
							position++
							if buffer[position] != rune('i') {
								goto l248
							}
							position++
							if buffer[position] != rune('v') {
								goto l248
							}
							position++
							if buffer[position] != rune('a') {
								goto l248
							}
							position++
							if buffer[position] != rune('t') {
								goto l248
							}
							position++
							if buffer[position] != rune('e') {
								goto l248
							}
							position++
							break
						case 'e':
							if buffer[position] != rune('e') {
								goto l248
							}
							position++
							if buffer[position] != rune('x') {
								goto l248
							}
							position++
							if buffer[position] != rune('p') {
								goto l248
							}
							position++
							if buffer[position] != rune('o') {
								goto l248
							}
							position++
							if buffer[position] != rune('r') {
								goto l248
							}
							position++
							if buffer[position] != rune('t') {
								goto l248
							}
							position++
							break
						default:
							if buffer[position] != rune('d') {
								goto l248
							}
							position++
							if buffer[position] != rune('e') {
								goto l248
							}
							position++
							if buffer[position] != rune('c') {
								goto l248
							}
							position++
							if buffer[position] != rune('l') {
								goto l248
							}
							position++
							if buffer[position] != rune('a') {
								goto l248
							}
							position++
							if buffer[position] != rune('r') {
								goto l248
							}
							position++
							if buffer[position] != rune('e') {
								goto l248
							}
							position++
							break
						}
					}

					if !_rules[ruleSPACE]() {
						goto l248
					}
					goto l247
				l248:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
				}
				depth--
				add(rulemodifier, position246)
			}
			return true
		l245:
			position, tokenIndex, depth = position245, tokenIndex245, depth245
			return false
		},
		/* 22 identifier <- <(((&('.' | '_') ('_' / '.')) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ((&('.' | '?' | '_') ((&('.') '.') | (&('?') '?') | (&('_') '_'))) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*)> */
		func() bool {
			position251, tokenIndex251, depth251 := position, tokenIndex, depth
			{
				position252 := position
				depth++
				{
					switch buffer[position] {
					case '.', '_':
						{
							position254, tokenIndex254, depth254 := position, tokenIndex, depth
							if buffer[position] != rune('_') {
								goto l255
							}
							position++
							goto l254
						l255:
							position, tokenIndex, depth = position254, tokenIndex254, depth254
							if buffer[position] != rune('.') {
								goto l251
							}
							position++
						}
					l254:
						break
					case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l251
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('a') || c > rune('z') {
							goto l251
						}
						position++
						break
					}
				}

			l256:
				{
					position257, tokenIndex257, depth257 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '.', '?', '_':
							{
								switch buffer[position] {
								case '.':
									if buffer[position] != rune('.') {
										goto l257
									}
									position++
									break
								case '?':
									if buffer[position] != rune('?') {
										goto l257
									}
									position++
									break
								default:
									if buffer[position] != rune('_') {
										goto l257
									}
									position++
									break
								}
							}

							break
						case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l257
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l257
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l257
							}
							position++
							break
						}
					}

					goto l256
				l257:
					position, tokenIndex, depth = position257, tokenIndex257, depth257
				}
				depth--
				add(ruleidentifier, position252)
			}
			return true
		l251:
			position, tokenIndex, depth = position251, tokenIndex251, depth251
			return false
		},
		/* 23 separator <- <(':' / ';' / '(' / ')' / '{' / '}' / ',' / '[' / ']' / '=' / '>')> */
		nil,
		/* 24 Comment <- <(LineComment / BlockComment)> */
		nil,
		/* 25 LineComment <- <('/' '/' (!eol .)* eol)> */
		nil,
		/* 26 BlockComment <- <('/' '*' (!('*' '/') .)* ('*' '/'))> */
		nil,
		/* 27 block_start <- <(SPACE '{' SPACE)> */
		func() bool {
			position264, tokenIndex264, depth264 := position, tokenIndex, depth
			{
				position265 := position
				depth++
				if !_rules[ruleSPACE]() {
					goto l264
				}
				if buffer[position] != rune('{') {
					goto l264
				}
				position++
				if !_rules[ruleSPACE]() {
					goto l264
				}
				depth--
				add(ruleblock_start, position265)
			}
			return true
		l264:
			position, tokenIndex, depth = position264, tokenIndex264, depth264
			return false
		},
		/* 28 block_end <- <(SPACE '}' SPACE)> */
		func() bool {
			position266, tokenIndex266, depth266 := position, tokenIndex, depth
			{
				position267 := position
				depth++
				if !_rules[ruleSPACE]() {
					goto l266
				}
				if buffer[position] != rune('}') {
					goto l266
				}
				position++
				if !_rules[ruleSPACE]() {
					goto l266
				}
				depth--
				add(ruleblock_end, position267)
			}
			return true
		l266:
			position, tokenIndex, depth = position266, tokenIndex266, depth266
			return false
		},
		/* 29 paren_start <- <(SPACE '(' SPACE)> */
		nil,
		/* 30 paren_end <- <(SPACE ')' SPACE)> */
		nil,
		/* 31 SPACE <- <spacing*> */
		func() bool {
			{
				position271 := position
				depth++
			l272:
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					{
						position274 := position
						depth++
						{
							switch buffer[position] {
							case '/':
								{
									position276 := position
									depth++
									{
										position277, tokenIndex277, depth277 := position, tokenIndex, depth
										{
											position279 := position
											depth++
											if buffer[position] != rune('/') {
												goto l278
											}
											position++
											if buffer[position] != rune('/') {
												goto l278
											}
											position++
										l280:
											{
												position281, tokenIndex281, depth281 := position, tokenIndex, depth
												{
													position282, tokenIndex282, depth282 := position, tokenIndex, depth
													if !_rules[ruleeol]() {
														goto l282
													}
													goto l281
												l282:
													position, tokenIndex, depth = position282, tokenIndex282, depth282
												}
												if !matchDot() {
													goto l281
												}
												goto l280
											l281:
												position, tokenIndex, depth = position281, tokenIndex281, depth281
											}
											if !_rules[ruleeol]() {
												goto l278
											}
											depth--
											add(ruleLineComment, position279)
										}
										goto l277
									l278:
										position, tokenIndex, depth = position277, tokenIndex277, depth277
										{
											position283 := position
											depth++
											if buffer[position] != rune('/') {
												goto l273
											}
											position++
											if buffer[position] != rune('*') {
												goto l273
											}
											position++
										l284:
											{
												position285, tokenIndex285, depth285 := position, tokenIndex, depth
												{
													position286, tokenIndex286, depth286 := position, tokenIndex, depth
													if buffer[position] != rune('*') {
														goto l286
													}
													position++
													if buffer[position] != rune('/') {
														goto l286
													}
													position++
													goto l285
												l286:
													position, tokenIndex, depth = position286, tokenIndex286, depth286
												}
												if !matchDot() {
													goto l285
												}
												goto l284
											l285:
												position, tokenIndex, depth = position285, tokenIndex285, depth285
											}
											if buffer[position] != rune('*') {
												goto l273
											}
											position++
											if buffer[position] != rune('/') {
												goto l273
											}
											position++
											depth--
											add(ruleBlockComment, position283)
										}
									}
								l277:
									depth--
									add(ruleComment, position276)
								}
								break
							case '\r':
								if buffer[position] != rune('\r') {
									goto l273
								}
								position++
								break
							case '\n':
								if buffer[position] != rune('\n') {
									goto l273
								}
								position++
								break
							case '\t':
								if buffer[position] != rune('\t') {
									goto l273
								}
								position++
								break
							default:
								if buffer[position] != rune(' ') {
									goto l273
								}
								position++
								break
							}
						}

						depth--
						add(rulespacing, position274)
					}
					goto l272
				l273:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
				}
				depth--
				add(ruleSPACE, position271)
			}
			return true
		},
		/* 32 ws <- <(' ' / '\t' / '\n' / '\r')> */
		nil,
		/* 33 spacing <- <((&('/') Comment) | (&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))> */
		nil,
		/* 34 eol <- <'\n'> */
		func() bool {
			position289, tokenIndex289, depth289 := position, tokenIndex, depth
			{
				position290 := position
				depth++
				if buffer[position] != rune('\n') {
					goto l289
				}
				position++
				depth--
				add(ruleeol, position290)
			}
			return true
		l289:
			position, tokenIndex, depth = position289, tokenIndex289, depth289
			return false
		},
		/* 35 eof <- <!.> */
		nil,
	}
	p.rules = _rules
}
