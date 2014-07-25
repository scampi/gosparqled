package sparql

import (
	"fmt"
	"math"
	"sort"
	"strconv"
)

const end_symbol rune = 4

/* The rule types inferred from the grammar are below. */
type pegRule uint8

const (
	ruleUnknown pegRule = iota
	rulequeryContainer
	ruleprolog
	ruleprefixDecl
	rulebaseDecl
	rulequery
	ruleselectQuery
	ruleselect
	rulesubSelect
	ruleprojectionElem
	ruledatasetClause
	rulewhereClause
	rulegroupGraphPattern
	rulegraphPattern
	rulegraphPatternNotTriples
	ruleoptionalGraphPattern
	rulegroupOrUnionGraphPattern
	rulebasicGraphPattern
	ruletriplesBlock
	ruletriplesSameSubjectPath
	rulevarOrTerm
	rulegraphTerm
	ruletriplesNodePath
	rulecollectionPath
	ruleblankNodePropertyListPath
	rulepropertyListPath
	ruleverbPath
	rulepath
	rulepathAlternative
	rulepathSequence
	rulepathElt
	rulepathPrimary
	rulepathNegatedPropertySet
	rulepathOneInPropertySet
	ruleobjectListPath
	ruleobjectPath
	rulegraphNodePath
	rulesolutionModifier
	rulelimitOffsetClauses
	rulelimit
	ruleoffset
	rulevar
	ruleiriref
	ruleiri
	ruleprefixedName
	ruleliteral
	rulestring
	rulenumericLiteral
	rulebooleanLiteral
	ruleblankNode
	ruleblankNodeLabel
	ruleanon
	rulenil
	ruleVARNAME
	rulePN_CHARS_U
	rulePN_CHARS_BASE
	rulePREFIX
	ruleTRUE
	ruleFALSE
	ruleBASE
	ruleSELECT
	ruleREDUCED
	ruleDISTINCT
	ruleFROM
	ruleNAMED
	ruleWHERE
	ruleLBRACE
	ruleRBRACE
	ruleLBRACK
	ruleRBRACK
	ruleSEMICOLON
	ruleCOMMA
	ruleDOT
	ruleCOLON
	rulePIPE
	ruleSLASH
	ruleINVERSE
	ruleLPAREN
	ruleRPAREN
	ruleISA
	ruleNOT
	ruleSTAR
	ruleOPTIONAL
	ruleUNION
	ruleLIMIT
	ruleOFFSET
	ruleINTEGER
	rulews

	rulePre_
	rule_In_
	rule_Suf
)

var rul3s = [...]string{
	"Unknown",
	"queryContainer",
	"prolog",
	"prefixDecl",
	"baseDecl",
	"query",
	"selectQuery",
	"select",
	"subSelect",
	"projectionElem",
	"datasetClause",
	"whereClause",
	"groupGraphPattern",
	"graphPattern",
	"graphPatternNotTriples",
	"optionalGraphPattern",
	"groupOrUnionGraphPattern",
	"basicGraphPattern",
	"triplesBlock",
	"triplesSameSubjectPath",
	"varOrTerm",
	"graphTerm",
	"triplesNodePath",
	"collectionPath",
	"blankNodePropertyListPath",
	"propertyListPath",
	"verbPath",
	"path",
	"pathAlternative",
	"pathSequence",
	"pathElt",
	"pathPrimary",
	"pathNegatedPropertySet",
	"pathOneInPropertySet",
	"objectListPath",
	"objectPath",
	"graphNodePath",
	"solutionModifier",
	"limitOffsetClauses",
	"limit",
	"offset",
	"var",
	"iriref",
	"iri",
	"prefixedName",
	"literal",
	"string",
	"numericLiteral",
	"booleanLiteral",
	"blankNode",
	"blankNodeLabel",
	"anon",
	"nil",
	"VARNAME",
	"PN_CHARS_U",
	"PN_CHARS_BASE",
	"PREFIX",
	"TRUE",
	"FALSE",
	"BASE",
	"SELECT",
	"REDUCED",
	"DISTINCT",
	"FROM",
	"NAMED",
	"WHERE",
	"LBRACE",
	"RBRACE",
	"LBRACK",
	"RBRACK",
	"SEMICOLON",
	"COMMA",
	"DOT",
	"COLON",
	"PIPE",
	"SLASH",
	"INVERSE",
	"LPAREN",
	"RPAREN",
	"ISA",
	"NOT",
	"STAR",
	"OPTIONAL",
	"UNION",
	"LIMIT",
	"OFFSET",
	"INTEGER",
	"ws",

	"Pre_",
	"_In_",
	"_Suf",
}

type tokenTree interface {
	Print()
	PrintSyntax()
	PrintSyntaxTree(buffer string)
	Add(rule pegRule, begin, end, next, depth int)
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
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[node.pegRule], strconv.Quote(buffer[node.begin:node.end]))
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
type token16 struct {
	pegRule
	begin, end, next int16
}

func (t *token16) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token16) isParentOf(u token16) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token16) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
}

func (t *token16) String() string {
	return fmt.Sprintf("\x1B[34m%v\x1B[m %v %v %v", rul3s[t.pegRule], t.begin, t.end, t.next)
}

type tokens16 struct {
	tree    []token16
	ordered [][]token16
}

func (t *tokens16) trim(length int) {
	t.tree = t.tree[0:length]
}

func (t *tokens16) Print() {
	for _, token := range t.tree {
		fmt.Println(token.String())
	}
}

func (t *tokens16) Order() [][]token16 {
	if t.ordered != nil {
		return t.ordered
	}

	depths := make([]int16, 1, math.MaxInt16)
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

	ordered, pool := make([][]token16, len(depths)), make([]token16, len(t.tree)+len(depths))
	for i, depth := range depths {
		depth++
		ordered[i], pool, depths[i] = pool[:depth], pool[depth:], 0
	}

	for i, token := range t.tree {
		depth := token.next
		token.next = int16(i)
		ordered[depth][depths[depth]] = token
		depths[depth]++
	}
	t.ordered = ordered
	return ordered
}

type state16 struct {
	token16
	depths []int16
	leaf   bool
}

func (t *tokens16) AST() *node32 {
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

func (t *tokens16) PreOrder() (<-chan state16, [][]token16) {
	s, ordered := make(chan state16, 6), t.Order()
	go func() {
		var states [8]state16
		for i, _ := range states {
			states[i].depths = make([]int16, len(ordered))
		}
		depths, state, depth := make([]int16, len(ordered)), 0, 1
		write := func(t token16, leaf bool) {
			S := states[state]
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int16(depth), leaf
			copy(S.depths, depths)
			s <- S
		}

		states[state].token16 = ordered[0][0]
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
							write(token16{pegRule: rule_In_, begin: c.end, end: b.begin}, true)
						}
						break
					}
				}

				if a.begin < b.begin {
					write(token16{pegRule: rulePre_, begin: a.begin, end: b.begin}, true)
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
					write(token16{pegRule: rule_Suf, begin: b.end, end: a.end}, true)
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

func (t *tokens16) PrintSyntax() {
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

func (t *tokens16) PrintSyntaxTree(buffer string) {
	tokens, _ := t.PreOrder()
	for token := range tokens {
		for c := 0; c < int(token.next); c++ {
			fmt.Printf(" ")
		}
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens16) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token16{pegRule: rule, begin: int16(begin), end: int16(end), next: int16(depth)}
}

func (t *tokens16) Tokens() <-chan token32 {
	s := make(chan token32, 16)
	go func() {
		for _, v := range t.tree {
			s <- v.getToken32()
		}
		close(s)
	}()
	return s
}

func (t *tokens16) Error() []token32 {
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

/* ${@} bit structure for abstract syntax tree */
type token32 struct {
	pegRule
	begin, end, next int32
}

func (t *token32) isZero() bool {
	return t.pegRule == ruleUnknown && t.begin == 0 && t.end == 0 && t.next == 0
}

func (t *token32) isParentOf(u token32) bool {
	return t.begin <= u.begin && t.end >= u.end && t.next > u.next
}

func (t *token32) getToken32() token32 {
	return token32{pegRule: t.pegRule, begin: int32(t.begin), end: int32(t.end), next: int32(t.next)}
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
		token.next = int32(i)
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
			state, S.pegRule, S.begin, S.end, S.next, S.leaf = (state+1)%8, t.pegRule, t.begin, t.end, int32(depth), leaf
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
		fmt.Printf("\x1B[34m%v\x1B[m %v\n", rul3s[token.pegRule], strconv.Quote(buffer[token.begin:token.end]))
	}
}

func (t *tokens32) Add(rule pegRule, begin, end, depth, index int) {
	t.tree[index] = token32{pegRule: rule, begin: int32(begin), end: int32(end), next: int32(depth)}
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

func (t *tokens16) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		for i, v := range tree {
			expanded[i] = v.getToken32()
		}
		return &tokens32{tree: expanded}
	}
	return nil
}

func (t *tokens32) Expand(index int) tokenTree {
	tree := t.tree
	if index >= len(tree) {
		expanded := make([]token32, 2*len(tree))
		copy(expanded, tree)
		t.tree = expanded
	}
	return nil
}

type Sparql struct {
	Buffer string
	buffer []rune
	rules  [88]func() bool
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
	p *Sparql
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

func (p *Sparql) PrintSyntaxTree() {
	p.tokenTree.PrintSyntaxTree(p.Buffer)
}

func (p *Sparql) Highlighter() {
	p.tokenTree.PrintSyntax()
}

func (p *Sparql) Init() {
	p.buffer = []rune(p.Buffer)
	if len(p.buffer) == 0 || p.buffer[len(p.buffer)-1] != end_symbol {
		p.buffer = append(p.buffer, end_symbol)
	}

	var tree tokenTree = &tokens16{tree: make([]token16, math.MaxInt16)}
	position, depth, tokenIndex, buffer, rules := 0, 0, 0, p.buffer, p.rules

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

	add := func(rule pegRule, begin int) {
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

	rules = [...]func() bool{
		nil,
		/* 0 queryContainer <- <(ws prolog query !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !rules[rulews]() {
					goto l0
				}
				{
					position2 := position
					depth++
				l3:
					{
						position4, tokenIndex4, depth4 := position, tokenIndex, depth
						{
							position5, tokenIndex5, depth5 := position, tokenIndex, depth
							{
								position7 := position
								depth++
								{
									position8 := position
									depth++
									{
										position9, tokenIndex9, depth9 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l10
										}
										position++
										goto l9
									l10:
										position, tokenIndex, depth = position9, tokenIndex9, depth9
										if buffer[position] != rune('P') {
											goto l6
										}
										position++
									}
								l9:
									{
										position11, tokenIndex11, depth11 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l12
										}
										position++
										goto l11
									l12:
										position, tokenIndex, depth = position11, tokenIndex11, depth11
										if buffer[position] != rune('R') {
											goto l6
										}
										position++
									}
								l11:
									{
										position13, tokenIndex13, depth13 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l14
										}
										position++
										goto l13
									l14:
										position, tokenIndex, depth = position13, tokenIndex13, depth13
										if buffer[position] != rune('E') {
											goto l6
										}
										position++
									}
								l13:
									{
										position15, tokenIndex15, depth15 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l16
										}
										position++
										goto l15
									l16:
										position, tokenIndex, depth = position15, tokenIndex15, depth15
										if buffer[position] != rune('F') {
											goto l6
										}
										position++
									}
								l15:
									{
										position17, tokenIndex17, depth17 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l18
										}
										position++
										goto l17
									l18:
										position, tokenIndex, depth = position17, tokenIndex17, depth17
										if buffer[position] != rune('I') {
											goto l6
										}
										position++
									}
								l17:
									{
										position19, tokenIndex19, depth19 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l20
										}
										position++
										goto l19
									l20:
										position, tokenIndex, depth = position19, tokenIndex19, depth19
										if buffer[position] != rune('X') {
											goto l6
										}
										position++
									}
								l19:
									if !rules[rulews]() {
										goto l6
									}
									depth--
									add(rulePREFIX, position8)
								}
								{
									position23, tokenIndex23, depth23 := position, tokenIndex, depth
									{
										position24, tokenIndex24, depth24 := position, tokenIndex, depth
										if buffer[position] != rune(':') {
											goto l25
										}
										position++
										goto l24
									l25:
										position, tokenIndex, depth = position24, tokenIndex24, depth24
										if buffer[position] != rune(' ') {
											goto l23
										}
										position++
									}
								l24:
									goto l6
								l23:
									position, tokenIndex, depth = position23, tokenIndex23, depth23
								}
								if !matchDot() {
									goto l6
								}
							l21:
								{
									position22, tokenIndex22, depth22 := position, tokenIndex, depth
									{
										position26, tokenIndex26, depth26 := position, tokenIndex, depth
										{
											position27, tokenIndex27, depth27 := position, tokenIndex, depth
											if buffer[position] != rune(':') {
												goto l28
											}
											position++
											goto l27
										l28:
											position, tokenIndex, depth = position27, tokenIndex27, depth27
											if buffer[position] != rune(' ') {
												goto l26
											}
											position++
										}
									l27:
										goto l22
									l26:
										position, tokenIndex, depth = position26, tokenIndex26, depth26
									}
									if !matchDot() {
										goto l22
									}
									goto l21
								l22:
									position, tokenIndex, depth = position22, tokenIndex22, depth22
								}
								{
									position29 := position
									depth++
									if buffer[position] != rune(':') {
										goto l6
									}
									position++
									if !rules[rulews]() {
										goto l6
									}
									depth--
									add(ruleCOLON, position29)
								}
								if !rules[ruleiri]() {
									goto l6
								}
								depth--
								add(ruleprefixDecl, position7)
							}
							goto l5
						l6:
							position, tokenIndex, depth = position5, tokenIndex5, depth5
							{
								position30 := position
								depth++
								{
									position31 := position
									depth++
									{
										position32, tokenIndex32, depth32 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l33
										}
										position++
										goto l32
									l33:
										position, tokenIndex, depth = position32, tokenIndex32, depth32
										if buffer[position] != rune('B') {
											goto l4
										}
										position++
									}
								l32:
									{
										position34, tokenIndex34, depth34 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l35
										}
										position++
										goto l34
									l35:
										position, tokenIndex, depth = position34, tokenIndex34, depth34
										if buffer[position] != rune('A') {
											goto l4
										}
										position++
									}
								l34:
									{
										position36, tokenIndex36, depth36 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l37
										}
										position++
										goto l36
									l37:
										position, tokenIndex, depth = position36, tokenIndex36, depth36
										if buffer[position] != rune('S') {
											goto l4
										}
										position++
									}
								l36:
									{
										position38, tokenIndex38, depth38 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l39
										}
										position++
										goto l38
									l39:
										position, tokenIndex, depth = position38, tokenIndex38, depth38
										if buffer[position] != rune('E') {
											goto l4
										}
										position++
									}
								l38:
									if !rules[rulews]() {
										goto l4
									}
									depth--
									add(ruleBASE, position31)
								}
								if !rules[ruleiri]() {
									goto l4
								}
								depth--
								add(rulebaseDecl, position30)
							}
						}
					l5:
						goto l3
					l4:
						position, tokenIndex, depth = position4, tokenIndex4, depth4
					}
					depth--
					add(ruleprolog, position2)
				}
				{
					position40 := position
					depth++
					{
						position41 := position
						depth++
						if !rules[ruleselect]() {
							goto l0
						}
						{
							position42, tokenIndex42, depth42 := position, tokenIndex, depth
							{
								position44 := position
								depth++
								{
									position45 := position
									depth++
									{
										position46, tokenIndex46, depth46 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l47
										}
										position++
										goto l46
									l47:
										position, tokenIndex, depth = position46, tokenIndex46, depth46
										if buffer[position] != rune('F') {
											goto l42
										}
										position++
									}
								l46:
									{
										position48, tokenIndex48, depth48 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l49
										}
										position++
										goto l48
									l49:
										position, tokenIndex, depth = position48, tokenIndex48, depth48
										if buffer[position] != rune('R') {
											goto l42
										}
										position++
									}
								l48:
									{
										position50, tokenIndex50, depth50 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l51
										}
										position++
										goto l50
									l51:
										position, tokenIndex, depth = position50, tokenIndex50, depth50
										if buffer[position] != rune('O') {
											goto l42
										}
										position++
									}
								l50:
									{
										position52, tokenIndex52, depth52 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l53
										}
										position++
										goto l52
									l53:
										position, tokenIndex, depth = position52, tokenIndex52, depth52
										if buffer[position] != rune('M') {
											goto l42
										}
										position++
									}
								l52:
									if !rules[rulews]() {
										goto l42
									}
									depth--
									add(ruleFROM, position45)
								}
								{
									position54, tokenIndex54, depth54 := position, tokenIndex, depth
									{
										position56 := position
										depth++
										{
											position57, tokenIndex57, depth57 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l58
											}
											position++
											goto l57
										l58:
											position, tokenIndex, depth = position57, tokenIndex57, depth57
											if buffer[position] != rune('N') {
												goto l54
											}
											position++
										}
									l57:
										{
											position59, tokenIndex59, depth59 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l60
											}
											position++
											goto l59
										l60:
											position, tokenIndex, depth = position59, tokenIndex59, depth59
											if buffer[position] != rune('A') {
												goto l54
											}
											position++
										}
									l59:
										{
											position61, tokenIndex61, depth61 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l62
											}
											position++
											goto l61
										l62:
											position, tokenIndex, depth = position61, tokenIndex61, depth61
											if buffer[position] != rune('M') {
												goto l54
											}
											position++
										}
									l61:
										{
											position63, tokenIndex63, depth63 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l64
											}
											position++
											goto l63
										l64:
											position, tokenIndex, depth = position63, tokenIndex63, depth63
											if buffer[position] != rune('E') {
												goto l54
											}
											position++
										}
									l63:
										{
											position65, tokenIndex65, depth65 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l66
											}
											position++
											goto l65
										l66:
											position, tokenIndex, depth = position65, tokenIndex65, depth65
											if buffer[position] != rune('D') {
												goto l54
											}
											position++
										}
									l65:
										if !rules[rulews]() {
											goto l54
										}
										depth--
										add(ruleNAMED, position56)
									}
									goto l55
								l54:
									position, tokenIndex, depth = position54, tokenIndex54, depth54
								}
							l55:
								if !rules[ruleiriref]() {
									goto l42
								}
								depth--
								add(ruledatasetClause, position44)
							}
							goto l43
						l42:
							position, tokenIndex, depth = position42, tokenIndex42, depth42
						}
					l43:
						if !rules[rulewhereClause]() {
							goto l0
						}
						{
							position67 := position
							depth++
							{
								position68, tokenIndex68, depth68 := position, tokenIndex, depth
								{
									position70 := position
									depth++
									{
										position71, tokenIndex71, depth71 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l72
										}
										{
											position73, tokenIndex73, depth73 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l73
											}
											goto l74
										l73:
											position, tokenIndex, depth = position73, tokenIndex73, depth73
										}
									l74:
										goto l71
									l72:
										position, tokenIndex, depth = position71, tokenIndex71, depth71
										if !rules[ruleoffset]() {
											goto l68
										}
										{
											position75, tokenIndex75, depth75 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l75
											}
											goto l76
										l75:
											position, tokenIndex, depth = position75, tokenIndex75, depth75
										}
									l76:
									}
								l71:
									depth--
									add(rulelimitOffsetClauses, position70)
								}
								goto l69
							l68:
								position, tokenIndex, depth = position68, tokenIndex68, depth68
							}
						l69:
							depth--
							add(rulesolutionModifier, position67)
						}
						depth--
						add(ruleselectQuery, position41)
					}
					depth--
					add(rulequery, position40)
				}
				{
					position77, tokenIndex77, depth77 := position, tokenIndex, depth
					if !matchDot() {
						goto l77
					}
					goto l0
				l77:
					position, tokenIndex, depth = position77, tokenIndex77, depth77
				}
				depth--
				add(rulequeryContainer, position1)
			}
			return true
		l0:
			position, tokenIndex, depth = position0, tokenIndex0, depth0
			return false
		},
		/* 1 prolog <- <(prefixDecl / baseDecl)*> */
		nil,
		/* 2 prefixDecl <- <(PREFIX (!(':' / ' ') .)+ COLON iri)> */
		nil,
		/* 3 baseDecl <- <(BASE iri)> */
		nil,
		/* 4 query <- <selectQuery> */
		nil,
		/* 5 selectQuery <- <(select datasetClause? whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position83, tokenIndex83, depth83 := position, tokenIndex, depth
			{
				position84 := position
				depth++
				{
					position85 := position
					depth++
					{
						position86, tokenIndex86, depth86 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l87
						}
						position++
						goto l86
					l87:
						position, tokenIndex, depth = position86, tokenIndex86, depth86
						if buffer[position] != rune('S') {
							goto l83
						}
						position++
					}
				l86:
					{
						position88, tokenIndex88, depth88 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l89
						}
						position++
						goto l88
					l89:
						position, tokenIndex, depth = position88, tokenIndex88, depth88
						if buffer[position] != rune('E') {
							goto l83
						}
						position++
					}
				l88:
					{
						position90, tokenIndex90, depth90 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l91
						}
						position++
						goto l90
					l91:
						position, tokenIndex, depth = position90, tokenIndex90, depth90
						if buffer[position] != rune('L') {
							goto l83
						}
						position++
					}
				l90:
					{
						position92, tokenIndex92, depth92 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l93
						}
						position++
						goto l92
					l93:
						position, tokenIndex, depth = position92, tokenIndex92, depth92
						if buffer[position] != rune('E') {
							goto l83
						}
						position++
					}
				l92:
					{
						position94, tokenIndex94, depth94 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l95
						}
						position++
						goto l94
					l95:
						position, tokenIndex, depth = position94, tokenIndex94, depth94
						if buffer[position] != rune('C') {
							goto l83
						}
						position++
					}
				l94:
					{
						position96, tokenIndex96, depth96 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l97
						}
						position++
						goto l96
					l97:
						position, tokenIndex, depth = position96, tokenIndex96, depth96
						if buffer[position] != rune('T') {
							goto l83
						}
						position++
					}
				l96:
					if !rules[rulews]() {
						goto l83
					}
					depth--
					add(ruleSELECT, position85)
				}
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					{
						position100, tokenIndex100, depth100 := position, tokenIndex, depth
						{
							position102 := position
							depth++
							{
								position103, tokenIndex103, depth103 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l104
								}
								position++
								goto l103
							l104:
								position, tokenIndex, depth = position103, tokenIndex103, depth103
								if buffer[position] != rune('D') {
									goto l101
								}
								position++
							}
						l103:
							{
								position105, tokenIndex105, depth105 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l106
								}
								position++
								goto l105
							l106:
								position, tokenIndex, depth = position105, tokenIndex105, depth105
								if buffer[position] != rune('I') {
									goto l101
								}
								position++
							}
						l105:
							{
								position107, tokenIndex107, depth107 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l108
								}
								position++
								goto l107
							l108:
								position, tokenIndex, depth = position107, tokenIndex107, depth107
								if buffer[position] != rune('S') {
									goto l101
								}
								position++
							}
						l107:
							{
								position109, tokenIndex109, depth109 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l110
								}
								position++
								goto l109
							l110:
								position, tokenIndex, depth = position109, tokenIndex109, depth109
								if buffer[position] != rune('T') {
									goto l101
								}
								position++
							}
						l109:
							{
								position111, tokenIndex111, depth111 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l112
								}
								position++
								goto l111
							l112:
								position, tokenIndex, depth = position111, tokenIndex111, depth111
								if buffer[position] != rune('I') {
									goto l101
								}
								position++
							}
						l111:
							{
								position113, tokenIndex113, depth113 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l114
								}
								position++
								goto l113
							l114:
								position, tokenIndex, depth = position113, tokenIndex113, depth113
								if buffer[position] != rune('N') {
									goto l101
								}
								position++
							}
						l113:
							{
								position115, tokenIndex115, depth115 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l116
								}
								position++
								goto l115
							l116:
								position, tokenIndex, depth = position115, tokenIndex115, depth115
								if buffer[position] != rune('C') {
									goto l101
								}
								position++
							}
						l115:
							{
								position117, tokenIndex117, depth117 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l118
								}
								position++
								goto l117
							l118:
								position, tokenIndex, depth = position117, tokenIndex117, depth117
								if buffer[position] != rune('T') {
									goto l101
								}
								position++
							}
						l117:
							if !rules[rulews]() {
								goto l101
							}
							depth--
							add(ruleDISTINCT, position102)
						}
						goto l100
					l101:
						position, tokenIndex, depth = position100, tokenIndex100, depth100
						{
							position119 := position
							depth++
							{
								position120, tokenIndex120, depth120 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l121
								}
								position++
								goto l120
							l121:
								position, tokenIndex, depth = position120, tokenIndex120, depth120
								if buffer[position] != rune('R') {
									goto l98
								}
								position++
							}
						l120:
							{
								position122, tokenIndex122, depth122 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l123
								}
								position++
								goto l122
							l123:
								position, tokenIndex, depth = position122, tokenIndex122, depth122
								if buffer[position] != rune('E') {
									goto l98
								}
								position++
							}
						l122:
							{
								position124, tokenIndex124, depth124 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l125
								}
								position++
								goto l124
							l125:
								position, tokenIndex, depth = position124, tokenIndex124, depth124
								if buffer[position] != rune('D') {
									goto l98
								}
								position++
							}
						l124:
							{
								position126, tokenIndex126, depth126 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l127
								}
								position++
								goto l126
							l127:
								position, tokenIndex, depth = position126, tokenIndex126, depth126
								if buffer[position] != rune('U') {
									goto l98
								}
								position++
							}
						l126:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l129
								}
								position++
								goto l128
							l129:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
								if buffer[position] != rune('C') {
									goto l98
								}
								position++
							}
						l128:
							{
								position130, tokenIndex130, depth130 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l131
								}
								position++
								goto l130
							l131:
								position, tokenIndex, depth = position130, tokenIndex130, depth130
								if buffer[position] != rune('E') {
									goto l98
								}
								position++
							}
						l130:
							{
								position132, tokenIndex132, depth132 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l133
								}
								position++
								goto l132
							l133:
								position, tokenIndex, depth = position132, tokenIndex132, depth132
								if buffer[position] != rune('D') {
									goto l98
								}
								position++
							}
						l132:
							if !rules[rulews]() {
								goto l98
							}
							depth--
							add(ruleREDUCED, position119)
						}
					}
				l100:
					goto l99
				l98:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
				}
			l99:
				{
					position134, tokenIndex134, depth134 := position, tokenIndex, depth
					{
						position136 := position
						depth++
						if buffer[position] != rune('*') {
							goto l135
						}
						position++
						if !rules[rulews]() {
							goto l135
						}
						depth--
						add(ruleSTAR, position136)
					}
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					{
						position139 := position
						depth++
						if !rules[rulevar]() {
							goto l83
						}
						depth--
						add(ruleprojectionElem, position139)
					}
				l137:
					{
						position138, tokenIndex138, depth138 := position, tokenIndex, depth
						{
							position140 := position
							depth++
							if !rules[rulevar]() {
								goto l138
							}
							depth--
							add(ruleprojectionElem, position140)
						}
						goto l137
					l138:
						position, tokenIndex, depth = position138, tokenIndex138, depth138
					}
				}
			l134:
				depth--
				add(ruleselect, position84)
			}
			return true
		l83:
			position, tokenIndex, depth = position83, tokenIndex83, depth83
			return false
		},
		/* 7 subSelect <- <(select whereClause)> */
		func() bool {
			position141, tokenIndex141, depth141 := position, tokenIndex, depth
			{
				position142 := position
				depth++
				if !rules[ruleselect]() {
					goto l141
				}
				if !rules[rulewhereClause]() {
					goto l141
				}
				depth--
				add(rulesubSelect, position142)
			}
			return true
		l141:
			position, tokenIndex, depth = position141, tokenIndex141, depth141
			return false
		},
		/* 8 projectionElem <- <var> */
		nil,
		/* 9 datasetClause <- <(FROM NAMED? iriref)> */
		nil,
		/* 10 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position145, tokenIndex145, depth145 := position, tokenIndex, depth
			{
				position146 := position
				depth++
				{
					position147, tokenIndex147, depth147 := position, tokenIndex, depth
					{
						position149 := position
						depth++
						{
							position150, tokenIndex150, depth150 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l151
							}
							position++
							goto l150
						l151:
							position, tokenIndex, depth = position150, tokenIndex150, depth150
							if buffer[position] != rune('W') {
								goto l147
							}
							position++
						}
					l150:
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l153
							}
							position++
							goto l152
						l153:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('H') {
								goto l147
							}
							position++
						}
					l152:
						{
							position154, tokenIndex154, depth154 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l155
							}
							position++
							goto l154
						l155:
							position, tokenIndex, depth = position154, tokenIndex154, depth154
							if buffer[position] != rune('E') {
								goto l147
							}
							position++
						}
					l154:
						{
							position156, tokenIndex156, depth156 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l157
							}
							position++
							goto l156
						l157:
							position, tokenIndex, depth = position156, tokenIndex156, depth156
							if buffer[position] != rune('R') {
								goto l147
							}
							position++
						}
					l156:
						{
							position158, tokenIndex158, depth158 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l159
							}
							position++
							goto l158
						l159:
							position, tokenIndex, depth = position158, tokenIndex158, depth158
							if buffer[position] != rune('E') {
								goto l147
							}
							position++
						}
					l158:
						if !rules[rulews]() {
							goto l147
						}
						depth--
						add(ruleWHERE, position149)
					}
					goto l148
				l147:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
				}
			l148:
				if !rules[rulegroupGraphPattern]() {
					goto l145
				}
				depth--
				add(rulewhereClause, position146)
			}
			return true
		l145:
			position, tokenIndex, depth = position145, tokenIndex145, depth145
			return false
		},
		/* 11 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l160
				}
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l163
					}
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if !rules[rulegraphPattern]() {
						goto l160
					}
				}
			l162:
				if !rules[ruleRBRACE]() {
					goto l160
				}
				depth--
				add(rulegroupGraphPattern, position161)
			}
			return true
		l160:
			position, tokenIndex, depth = position160, tokenIndex160, depth160
			return false
		},
		/* 12 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position165 := position
				depth++
				{
					position166, tokenIndex166, depth166 := position, tokenIndex, depth
					{
						position168 := position
						depth++
						{
							position169 := position
							depth++
							if !rules[ruletriplesSameSubjectPath]() {
								goto l166
							}
						l170:
							{
								position171, tokenIndex171, depth171 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l171
								}
								if !rules[ruletriplesSameSubjectPath]() {
									goto l171
								}
								goto l170
							l171:
								position, tokenIndex, depth = position171, tokenIndex171, depth171
							}
							{
								position172, tokenIndex172, depth172 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l172
								}
								goto l173
							l172:
								position, tokenIndex, depth = position172, tokenIndex172, depth172
							}
						l173:
							depth--
							add(ruletriplesBlock, position169)
						}
						depth--
						add(rulebasicGraphPattern, position168)
					}
					goto l167
				l166:
					position, tokenIndex, depth = position166, tokenIndex166, depth166
				}
			l167:
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					{
						position176 := position
						depth++
						{
							position177, tokenIndex177, depth177 := position, tokenIndex, depth
							{
								position179 := position
								depth++
								{
									position180 := position
									depth++
									{
										position181, tokenIndex181, depth181 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l182
										}
										position++
										goto l181
									l182:
										position, tokenIndex, depth = position181, tokenIndex181, depth181
										if buffer[position] != rune('O') {
											goto l178
										}
										position++
									}
								l181:
									{
										position183, tokenIndex183, depth183 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l184
										}
										position++
										goto l183
									l184:
										position, tokenIndex, depth = position183, tokenIndex183, depth183
										if buffer[position] != rune('P') {
											goto l178
										}
										position++
									}
								l183:
									{
										position185, tokenIndex185, depth185 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l186
										}
										position++
										goto l185
									l186:
										position, tokenIndex, depth = position185, tokenIndex185, depth185
										if buffer[position] != rune('T') {
											goto l178
										}
										position++
									}
								l185:
									{
										position187, tokenIndex187, depth187 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l188
										}
										position++
										goto l187
									l188:
										position, tokenIndex, depth = position187, tokenIndex187, depth187
										if buffer[position] != rune('I') {
											goto l178
										}
										position++
									}
								l187:
									{
										position189, tokenIndex189, depth189 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l190
										}
										position++
										goto l189
									l190:
										position, tokenIndex, depth = position189, tokenIndex189, depth189
										if buffer[position] != rune('O') {
											goto l178
										}
										position++
									}
								l189:
									{
										position191, tokenIndex191, depth191 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l192
										}
										position++
										goto l191
									l192:
										position, tokenIndex, depth = position191, tokenIndex191, depth191
										if buffer[position] != rune('N') {
											goto l178
										}
										position++
									}
								l191:
									{
										position193, tokenIndex193, depth193 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l194
										}
										position++
										goto l193
									l194:
										position, tokenIndex, depth = position193, tokenIndex193, depth193
										if buffer[position] != rune('A') {
											goto l178
										}
										position++
									}
								l193:
									{
										position195, tokenIndex195, depth195 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l196
										}
										position++
										goto l195
									l196:
										position, tokenIndex, depth = position195, tokenIndex195, depth195
										if buffer[position] != rune('L') {
											goto l178
										}
										position++
									}
								l195:
									if !rules[rulews]() {
										goto l178
									}
									depth--
									add(ruleOPTIONAL, position180)
								}
								if !rules[ruleLBRACE]() {
									goto l178
								}
								{
									position197, tokenIndex197, depth197 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l198
									}
									goto l197
								l198:
									position, tokenIndex, depth = position197, tokenIndex197, depth197
									if !rules[rulegraphPattern]() {
										goto l178
									}
								}
							l197:
								if !rules[ruleRBRACE]() {
									goto l178
								}
								depth--
								add(ruleoptionalGraphPattern, position179)
							}
							goto l177
						l178:
							position, tokenIndex, depth = position177, tokenIndex177, depth177
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l174
							}
						}
					l177:
						depth--
						add(rulegraphPatternNotTriples, position176)
					}
					{
						position199, tokenIndex199, depth199 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l199
						}
						goto l200
					l199:
						position, tokenIndex, depth = position199, tokenIndex199, depth199
					}
				l200:
					if !rules[rulegraphPattern]() {
						goto l174
					}
					goto l175
				l174:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
				}
			l175:
				depth--
				add(rulegraphPattern, position165)
			}
			return true
		},
		/* 13 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 14 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 15 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position203, tokenIndex203, depth203 := position, tokenIndex, depth
			{
				position204 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l203
				}
				{
					position205, tokenIndex205, depth205 := position, tokenIndex, depth
					{
						position207 := position
						depth++
						{
							position208, tokenIndex208, depth208 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l209
							}
							position++
							goto l208
						l209:
							position, tokenIndex, depth = position208, tokenIndex208, depth208
							if buffer[position] != rune('U') {
								goto l205
							}
							position++
						}
					l208:
						{
							position210, tokenIndex210, depth210 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l211
							}
							position++
							goto l210
						l211:
							position, tokenIndex, depth = position210, tokenIndex210, depth210
							if buffer[position] != rune('N') {
								goto l205
							}
							position++
						}
					l210:
						{
							position212, tokenIndex212, depth212 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l213
							}
							position++
							goto l212
						l213:
							position, tokenIndex, depth = position212, tokenIndex212, depth212
							if buffer[position] != rune('I') {
								goto l205
							}
							position++
						}
					l212:
						{
							position214, tokenIndex214, depth214 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l215
							}
							position++
							goto l214
						l215:
							position, tokenIndex, depth = position214, tokenIndex214, depth214
							if buffer[position] != rune('O') {
								goto l205
							}
							position++
						}
					l214:
						{
							position216, tokenIndex216, depth216 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l217
							}
							position++
							goto l216
						l217:
							position, tokenIndex, depth = position216, tokenIndex216, depth216
							if buffer[position] != rune('N') {
								goto l205
							}
							position++
						}
					l216:
						if !rules[rulews]() {
							goto l205
						}
						depth--
						add(ruleUNION, position207)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l205
					}
					goto l206
				l205:
					position, tokenIndex, depth = position205, tokenIndex205, depth205
				}
			l206:
				depth--
				add(rulegroupOrUnionGraphPattern, position204)
			}
			return true
		l203:
			position, tokenIndex, depth = position203, tokenIndex203, depth203
			return false
		},
		/* 16 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 17 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		nil,
		/* 18 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position220, tokenIndex220, depth220 := position, tokenIndex, depth
			{
				position221 := position
				depth++
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l223
					}
					if !rules[rulepropertyListPath]() {
						goto l223
					}
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					{
						position224 := position
						depth++
						{
							position225, tokenIndex225, depth225 := position, tokenIndex, depth
							{
								position227 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l226
								}
								if !rules[rulegraphNodePath]() {
									goto l226
								}
							l228:
								{
									position229, tokenIndex229, depth229 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l229
									}
									goto l228
								l229:
									position, tokenIndex, depth = position229, tokenIndex229, depth229
								}
								if !rules[ruleRPAREN]() {
									goto l226
								}
								depth--
								add(rulecollectionPath, position227)
							}
							goto l225
						l226:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							{
								position230 := position
								depth++
								{
									position231 := position
									depth++
									if buffer[position] != rune('[') {
										goto l220
									}
									position++
									if !rules[rulews]() {
										goto l220
									}
									depth--
									add(ruleLBRACK, position231)
								}
								if !rules[rulepropertyListPath]() {
									goto l220
								}
								{
									position232 := position
									depth++
									if buffer[position] != rune(']') {
										goto l220
									}
									position++
									if !rules[rulews]() {
										goto l220
									}
									depth--
									add(ruleRBRACK, position232)
								}
								depth--
								add(ruleblankNodePropertyListPath, position230)
							}
						}
					l225:
						depth--
						add(ruletriplesNodePath, position224)
					}
					if !rules[rulepropertyListPath]() {
						goto l220
					}
				}
			l222:
				depth--
				add(ruletriplesSameSubjectPath, position221)
			}
			return true
		l220:
			position, tokenIndex, depth = position220, tokenIndex220, depth220
			return false
		},
		/* 19 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position233, tokenIndex233, depth233 := position, tokenIndex, depth
			{
				position234 := position
				depth++
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l236
					}
					goto l235
				l236:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
					{
						position237 := position
						depth++
						{
							position238, tokenIndex238, depth238 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l239
							}
							goto l238
						l239:
							position, tokenIndex, depth = position238, tokenIndex238, depth238
							{
								switch buffer[position] {
								case '(':
									{
										position241 := position
										depth++
										if buffer[position] != rune('(') {
											goto l233
										}
										position++
										if !rules[rulews]() {
											goto l233
										}
										if buffer[position] != rune(')') {
											goto l233
										}
										position++
										if !rules[rulews]() {
											goto l233
										}
										depth--
										add(rulenil, position241)
									}
									break
								case '[', '_':
									{
										position242 := position
										depth++
										{
											position243, tokenIndex243, depth243 := position, tokenIndex, depth
											{
												position245 := position
												depth++
												if buffer[position] != rune('_') {
													goto l244
												}
												position++
												if buffer[position] != rune(':') {
													goto l244
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l244
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l244
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l244
														}
														position++
														break
													}
												}

												{
													position247, tokenIndex247, depth247 := position, tokenIndex, depth
													{
														position249, tokenIndex249, depth249 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l250
														}
														position++
														goto l249
													l250:
														position, tokenIndex, depth = position249, tokenIndex249, depth249
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l251
														}
														position++
														goto l249
													l251:
														position, tokenIndex, depth = position249, tokenIndex249, depth249
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l252
														}
														position++
														goto l249
													l252:
														position, tokenIndex, depth = position249, tokenIndex249, depth249
														if c := buffer[position]; c < rune('.') || c > rune('_') {
															goto l247
														}
														position++
													}
												l249:
													goto l248
												l247:
													position, tokenIndex, depth = position247, tokenIndex247, depth247
												}
											l248:
												if !rules[rulews]() {
													goto l244
												}
												depth--
												add(ruleblankNodeLabel, position245)
											}
											goto l243
										l244:
											position, tokenIndex, depth = position243, tokenIndex243, depth243
											{
												position253 := position
												depth++
												if buffer[position] != rune('[') {
													goto l233
												}
												position++
												if !rules[rulews]() {
													goto l233
												}
												if buffer[position] != rune(']') {
													goto l233
												}
												position++
												if !rules[rulews]() {
													goto l233
												}
												depth--
												add(ruleanon, position253)
											}
										}
									l243:
										depth--
										add(ruleblankNode, position242)
									}
									break
								case 'F', 'T', 'f', 't':
									{
										position254 := position
										depth++
										{
											position255, tokenIndex255, depth255 := position, tokenIndex, depth
											{
												position257 := position
												depth++
												{
													position258, tokenIndex258, depth258 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l259
													}
													position++
													goto l258
												l259:
													position, tokenIndex, depth = position258, tokenIndex258, depth258
													if buffer[position] != rune('T') {
														goto l256
													}
													position++
												}
											l258:
												{
													position260, tokenIndex260, depth260 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l261
													}
													position++
													goto l260
												l261:
													position, tokenIndex, depth = position260, tokenIndex260, depth260
													if buffer[position] != rune('R') {
														goto l256
													}
													position++
												}
											l260:
												{
													position262, tokenIndex262, depth262 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l263
													}
													position++
													goto l262
												l263:
													position, tokenIndex, depth = position262, tokenIndex262, depth262
													if buffer[position] != rune('U') {
														goto l256
													}
													position++
												}
											l262:
												{
													position264, tokenIndex264, depth264 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l265
													}
													position++
													goto l264
												l265:
													position, tokenIndex, depth = position264, tokenIndex264, depth264
													if buffer[position] != rune('E') {
														goto l256
													}
													position++
												}
											l264:
												if !rules[rulews]() {
													goto l256
												}
												depth--
												add(ruleTRUE, position257)
											}
											goto l255
										l256:
											position, tokenIndex, depth = position255, tokenIndex255, depth255
											{
												position266 := position
												depth++
												{
													position267, tokenIndex267, depth267 := position, tokenIndex, depth
													if buffer[position] != rune('f') {
														goto l268
													}
													position++
													goto l267
												l268:
													position, tokenIndex, depth = position267, tokenIndex267, depth267
													if buffer[position] != rune('F') {
														goto l233
													}
													position++
												}
											l267:
												{
													position269, tokenIndex269, depth269 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l270
													}
													position++
													goto l269
												l270:
													position, tokenIndex, depth = position269, tokenIndex269, depth269
													if buffer[position] != rune('A') {
														goto l233
													}
													position++
												}
											l269:
												{
													position271, tokenIndex271, depth271 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l272
													}
													position++
													goto l271
												l272:
													position, tokenIndex, depth = position271, tokenIndex271, depth271
													if buffer[position] != rune('L') {
														goto l233
													}
													position++
												}
											l271:
												{
													position273, tokenIndex273, depth273 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l274
													}
													position++
													goto l273
												l274:
													position, tokenIndex, depth = position273, tokenIndex273, depth273
													if buffer[position] != rune('S') {
														goto l233
													}
													position++
												}
											l273:
												{
													position275, tokenIndex275, depth275 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l276
													}
													position++
													goto l275
												l276:
													position, tokenIndex, depth = position275, tokenIndex275, depth275
													if buffer[position] != rune('E') {
														goto l233
													}
													position++
												}
											l275:
												if !rules[rulews]() {
													goto l233
												}
												depth--
												add(ruleFALSE, position266)
											}
										}
									l255:
										depth--
										add(rulebooleanLiteral, position254)
									}
									break
								case '"':
									{
										position277 := position
										depth++
										{
											position278 := position
											depth++
											if buffer[position] != rune('"') {
												goto l233
											}
											position++
										l279:
											{
												position280, tokenIndex280, depth280 := position, tokenIndex, depth
												{
													position281, tokenIndex281, depth281 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l281
													}
													position++
													goto l280
												l281:
													position, tokenIndex, depth = position281, tokenIndex281, depth281
												}
												if !matchDot() {
													goto l280
												}
												goto l279
											l280:
												position, tokenIndex, depth = position280, tokenIndex280, depth280
											}
											if buffer[position] != rune('"') {
												goto l233
											}
											position++
											depth--
											add(rulestring, position278)
										}
										{
											position282, tokenIndex282, depth282 := position, tokenIndex, depth
											{
												position284, tokenIndex284, depth284 := position, tokenIndex, depth
												if buffer[position] != rune('@') {
													goto l285
												}
												position++
												{
													position288, tokenIndex288, depth288 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l289
													}
													position++
													goto l288
												l289:
													position, tokenIndex, depth = position288, tokenIndex288, depth288
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l285
													}
													position++
												}
											l288:
											l286:
												{
													position287, tokenIndex287, depth287 := position, tokenIndex, depth
													{
														position290, tokenIndex290, depth290 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l291
														}
														position++
														goto l290
													l291:
														position, tokenIndex, depth = position290, tokenIndex290, depth290
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l287
														}
														position++
													}
												l290:
													goto l286
												l287:
													position, tokenIndex, depth = position287, tokenIndex287, depth287
												}
											l292:
												{
													position293, tokenIndex293, depth293 := position, tokenIndex, depth
													if buffer[position] != rune('-') {
														goto l293
													}
													position++
													{
														switch buffer[position] {
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l293
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l293
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l293
															}
															position++
															break
														}
													}

												l294:
													{
														position295, tokenIndex295, depth295 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l295
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l295
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l295
																}
																position++
																break
															}
														}

														goto l294
													l295:
														position, tokenIndex, depth = position295, tokenIndex295, depth295
													}
													goto l292
												l293:
													position, tokenIndex, depth = position293, tokenIndex293, depth293
												}
												goto l284
											l285:
												position, tokenIndex, depth = position284, tokenIndex284, depth284
												if buffer[position] != rune('^') {
													goto l282
												}
												position++
												if buffer[position] != rune('^') {
													goto l282
												}
												position++
												if !rules[ruleiriref]() {
													goto l282
												}
											}
										l284:
											goto l283
										l282:
											position, tokenIndex, depth = position282, tokenIndex282, depth282
										}
									l283:
										if !rules[rulews]() {
											goto l233
										}
										depth--
										add(ruleliteral, position277)
									}
									break
								default:
									{
										position298 := position
										depth++
										{
											position299, tokenIndex299, depth299 := position, tokenIndex, depth
											{
												position301, tokenIndex301, depth301 := position, tokenIndex, depth
												if buffer[position] != rune('+') {
													goto l302
												}
												position++
												goto l301
											l302:
												position, tokenIndex, depth = position301, tokenIndex301, depth301
												if buffer[position] != rune('-') {
													goto l299
												}
												position++
											}
										l301:
											goto l300
										l299:
											position, tokenIndex, depth = position299, tokenIndex299, depth299
										}
									l300:
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l233
										}
										position++
									l303:
										{
											position304, tokenIndex304, depth304 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l304
											}
											position++
											goto l303
										l304:
											position, tokenIndex, depth = position304, tokenIndex304, depth304
										}
										{
											position305, tokenIndex305, depth305 := position, tokenIndex, depth
											if buffer[position] != rune('.') {
												goto l305
											}
											position++
										l307:
											{
												position308, tokenIndex308, depth308 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l308
												}
												position++
												goto l307
											l308:
												position, tokenIndex, depth = position308, tokenIndex308, depth308
											}
											goto l306
										l305:
											position, tokenIndex, depth = position305, tokenIndex305, depth305
										}
									l306:
										if !rules[rulews]() {
											goto l233
										}
										depth--
										add(rulenumericLiteral, position298)
									}
									break
								}
							}

						}
					l238:
						depth--
						add(rulegraphTerm, position237)
					}
				}
			l235:
				depth--
				add(rulevarOrTerm, position234)
			}
			return true
		l233:
			position, tokenIndex, depth = position233, tokenIndex233, depth233
			return false
		},
		/* 20 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 21 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 22 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 23 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 24 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			position313, tokenIndex313, depth313 := position, tokenIndex, depth
			{
				position314 := position
				depth++
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l316
					}
					goto l315
				l316:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
					{
						position317 := position
						depth++
						if !rules[rulepath]() {
							goto l313
						}
						depth--
						add(ruleverbPath, position317)
					}
				}
			l315:
				if !rules[ruleobjectListPath]() {
					goto l313
				}
				{
					position318, tokenIndex318, depth318 := position, tokenIndex, depth
					{
						position320 := position
						depth++
						if buffer[position] != rune(';') {
							goto l318
						}
						position++
						if !rules[rulews]() {
							goto l318
						}
						depth--
						add(ruleSEMICOLON, position320)
					}
					if !rules[rulepropertyListPath]() {
						goto l318
					}
					goto l319
				l318:
					position, tokenIndex, depth = position318, tokenIndex318, depth318
				}
			l319:
				depth--
				add(rulepropertyListPath, position314)
			}
			return true
		l313:
			position, tokenIndex, depth = position313, tokenIndex313, depth313
			return false
		},
		/* 25 verbPath <- <path> */
		nil,
		/* 26 path <- <pathAlternative> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l322
				}
				depth--
				add(rulepath, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 27 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position324, tokenIndex324, depth324 := position, tokenIndex, depth
			{
				position325 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l324
				}
			l326:
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l327
					}
					if !rules[rulepathAlternative]() {
						goto l327
					}
					goto l326
				l327:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
				}
				depth--
				add(rulepathAlternative, position325)
			}
			return true
		l324:
			position, tokenIndex, depth = position324, tokenIndex324, depth324
			return false
		},
		/* 28 pathSequence <- <(pathElt (SLASH pathSequence)*)> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				{
					position330 := position
					depth++
					{
						position331, tokenIndex331, depth331 := position, tokenIndex, depth
						if !rules[ruleINVERSE]() {
							goto l331
						}
						goto l332
					l331:
						position, tokenIndex, depth = position331, tokenIndex331, depth331
					}
				l332:
					{
						position333 := position
						depth++
						{
							position334, tokenIndex334, depth334 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l335
							}
							goto l334
						l335:
							position, tokenIndex, depth = position334, tokenIndex334, depth334
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l328
									}
									if !rules[rulepath]() {
										goto l328
									}
									if !rules[ruleRPAREN]() {
										goto l328
									}
									break
								case '!':
									{
										position337 := position
										depth++
										if buffer[position] != rune('!') {
											goto l328
										}
										position++
										if !rules[rulews]() {
											goto l328
										}
										depth--
										add(ruleNOT, position337)
									}
									{
										position338 := position
										depth++
										{
											position339, tokenIndex339, depth339 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l340
											}
											goto l339
										l340:
											position, tokenIndex, depth = position339, tokenIndex339, depth339
											if !rules[ruleLPAREN]() {
												goto l328
											}
											{
												position341, tokenIndex341, depth341 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l341
												}
											l343:
												{
													position344, tokenIndex344, depth344 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l344
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l344
													}
													goto l343
												l344:
													position, tokenIndex, depth = position344, tokenIndex344, depth344
												}
												goto l342
											l341:
												position, tokenIndex, depth = position341, tokenIndex341, depth341
											}
										l342:
											if !rules[ruleRPAREN]() {
												goto l328
											}
										}
									l339:
										depth--
										add(rulepathNegatedPropertySet, position338)
									}
									break
								default:
									if !rules[ruleISA]() {
										goto l328
									}
									break
								}
							}

						}
					l334:
						depth--
						add(rulepathPrimary, position333)
					}
					depth--
					add(rulepathElt, position330)
				}
			l345:
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					{
						position347 := position
						depth++
						if buffer[position] != rune('/') {
							goto l346
						}
						position++
						if !rules[rulews]() {
							goto l346
						}
						depth--
						add(ruleSLASH, position347)
					}
					if !rules[rulepathSequence]() {
						goto l346
					}
					goto l345
				l346:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
				}
				depth--
				add(rulepathSequence, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 29 pathElt <- <(INVERSE? pathPrimary)> */
		nil,
		/* 30 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 31 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 32 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position351, tokenIndex351, depth351 := position, tokenIndex, depth
			{
				position352 := position
				depth++
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l354
					}
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if !rules[ruleISA]() {
						goto l355
					}
					goto l353
				l355:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					if !rules[ruleINVERSE]() {
						goto l351
					}
					{
						position356, tokenIndex356, depth356 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l357
						}
						goto l356
					l357:
						position, tokenIndex, depth = position356, tokenIndex356, depth356
						if !rules[ruleISA]() {
							goto l351
						}
					}
				l356:
				}
			l353:
				depth--
				add(rulepathOneInPropertySet, position352)
			}
			return true
		l351:
			position, tokenIndex, depth = position351, tokenIndex351, depth351
			return false
		},
		/* 33 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position358, tokenIndex358, depth358 := position, tokenIndex, depth
			{
				position359 := position
				depth++
				{
					position360 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l358
					}
					depth--
					add(ruleobjectPath, position360)
				}
			l361:
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					{
						position363 := position
						depth++
						if buffer[position] != rune(',') {
							goto l362
						}
						position++
						if !rules[rulews]() {
							goto l362
						}
						depth--
						add(ruleCOMMA, position363)
					}
					if !rules[ruleobjectListPath]() {
						goto l362
					}
					goto l361
				l362:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
				}
				depth--
				add(ruleobjectListPath, position359)
			}
			return true
		l358:
			position, tokenIndex, depth = position358, tokenIndex358, depth358
			return false
		},
		/* 34 objectPath <- <graphNodePath> */
		nil,
		/* 35 graphNodePath <- <varOrTerm> */
		func() bool {
			position365, tokenIndex365, depth365 := position, tokenIndex, depth
			{
				position366 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l365
				}
				depth--
				add(rulegraphNodePath, position366)
			}
			return true
		l365:
			position, tokenIndex, depth = position365, tokenIndex365, depth365
			return false
		},
		/* 36 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 37 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 38 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position369, tokenIndex369, depth369 := position, tokenIndex, depth
			{
				position370 := position
				depth++
				{
					position371 := position
					depth++
					{
						position372, tokenIndex372, depth372 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l373
						}
						position++
						goto l372
					l373:
						position, tokenIndex, depth = position372, tokenIndex372, depth372
						if buffer[position] != rune('L') {
							goto l369
						}
						position++
					}
				l372:
					{
						position374, tokenIndex374, depth374 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l375
						}
						position++
						goto l374
					l375:
						position, tokenIndex, depth = position374, tokenIndex374, depth374
						if buffer[position] != rune('I') {
							goto l369
						}
						position++
					}
				l374:
					{
						position376, tokenIndex376, depth376 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l377
						}
						position++
						goto l376
					l377:
						position, tokenIndex, depth = position376, tokenIndex376, depth376
						if buffer[position] != rune('M') {
							goto l369
						}
						position++
					}
				l376:
					{
						position378, tokenIndex378, depth378 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l379
						}
						position++
						goto l378
					l379:
						position, tokenIndex, depth = position378, tokenIndex378, depth378
						if buffer[position] != rune('I') {
							goto l369
						}
						position++
					}
				l378:
					{
						position380, tokenIndex380, depth380 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l381
						}
						position++
						goto l380
					l381:
						position, tokenIndex, depth = position380, tokenIndex380, depth380
						if buffer[position] != rune('T') {
							goto l369
						}
						position++
					}
				l380:
					if !rules[rulews]() {
						goto l369
					}
					depth--
					add(ruleLIMIT, position371)
				}
				if !rules[ruleINTEGER]() {
					goto l369
				}
				depth--
				add(rulelimit, position370)
			}
			return true
		l369:
			position, tokenIndex, depth = position369, tokenIndex369, depth369
			return false
		},
		/* 39 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position382, tokenIndex382, depth382 := position, tokenIndex, depth
			{
				position383 := position
				depth++
				{
					position384 := position
					depth++
					{
						position385, tokenIndex385, depth385 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l386
						}
						position++
						goto l385
					l386:
						position, tokenIndex, depth = position385, tokenIndex385, depth385
						if buffer[position] != rune('O') {
							goto l382
						}
						position++
					}
				l385:
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l388
						}
						position++
						goto l387
					l388:
						position, tokenIndex, depth = position387, tokenIndex387, depth387
						if buffer[position] != rune('F') {
							goto l382
						}
						position++
					}
				l387:
					{
						position389, tokenIndex389, depth389 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l390
						}
						position++
						goto l389
					l390:
						position, tokenIndex, depth = position389, tokenIndex389, depth389
						if buffer[position] != rune('F') {
							goto l382
						}
						position++
					}
				l389:
					{
						position391, tokenIndex391, depth391 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l392
						}
						position++
						goto l391
					l392:
						position, tokenIndex, depth = position391, tokenIndex391, depth391
						if buffer[position] != rune('S') {
							goto l382
						}
						position++
					}
				l391:
					{
						position393, tokenIndex393, depth393 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l394
						}
						position++
						goto l393
					l394:
						position, tokenIndex, depth = position393, tokenIndex393, depth393
						if buffer[position] != rune('E') {
							goto l382
						}
						position++
					}
				l393:
					{
						position395, tokenIndex395, depth395 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l396
						}
						position++
						goto l395
					l396:
						position, tokenIndex, depth = position395, tokenIndex395, depth395
						if buffer[position] != rune('T') {
							goto l382
						}
						position++
					}
				l395:
					if !rules[rulews]() {
						goto l382
					}
					depth--
					add(ruleOFFSET, position384)
				}
				if !rules[ruleINTEGER]() {
					goto l382
				}
				depth--
				add(ruleoffset, position383)
			}
			return true
		l382:
			position, tokenIndex, depth = position382, tokenIndex382, depth382
			return false
		},
		/* 40 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position397, tokenIndex397, depth397 := position, tokenIndex, depth
			{
				position398 := position
				depth++
				{
					position399, tokenIndex399, depth399 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l400
					}
					position++
					goto l399
				l400:
					position, tokenIndex, depth = position399, tokenIndex399, depth399
					if buffer[position] != rune('$') {
						goto l397
					}
					position++
				}
			l399:
				{
					position401 := position
					depth++
					{
						position404, tokenIndex404, depth404 := position, tokenIndex, depth
						{
							position406 := position
							depth++
							{
								position407, tokenIndex407, depth407 := position, tokenIndex, depth
								{
									position409 := position
									depth++
									{
										position410, tokenIndex410, depth410 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l411
										}
										position++
										goto l410
									l411:
										position, tokenIndex, depth = position410, tokenIndex410, depth410
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l408
										}
										position++
									}
								l410:
									depth--
									add(rulePN_CHARS_BASE, position409)
								}
								goto l407
							l408:
								position, tokenIndex, depth = position407, tokenIndex407, depth407
								if buffer[position] != rune('_') {
									goto l405
								}
								position++
							}
						l407:
							depth--
							add(rulePN_CHARS_U, position406)
						}
						goto l404
					l405:
						position, tokenIndex, depth = position404, tokenIndex404, depth404
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l397
						}
						position++
					}
				l404:
				l402:
					{
						position403, tokenIndex403, depth403 := position, tokenIndex, depth
						{
							position412, tokenIndex412, depth412 := position, tokenIndex, depth
							{
								position414 := position
								depth++
								{
									position415, tokenIndex415, depth415 := position, tokenIndex, depth
									{
										position417 := position
										depth++
										{
											position418, tokenIndex418, depth418 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l419
											}
											position++
											goto l418
										l419:
											position, tokenIndex, depth = position418, tokenIndex418, depth418
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l416
											}
											position++
										}
									l418:
										depth--
										add(rulePN_CHARS_BASE, position417)
									}
									goto l415
								l416:
									position, tokenIndex, depth = position415, tokenIndex415, depth415
									if buffer[position] != rune('_') {
										goto l413
									}
									position++
								}
							l415:
								depth--
								add(rulePN_CHARS_U, position414)
							}
							goto l412
						l413:
							position, tokenIndex, depth = position412, tokenIndex412, depth412
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l403
							}
							position++
						}
					l412:
						goto l402
					l403:
						position, tokenIndex, depth = position403, tokenIndex403, depth403
					}
					depth--
					add(ruleVARNAME, position401)
				}
				if !rules[rulews]() {
					goto l397
				}
				depth--
				add(rulevar, position398)
			}
			return true
		l397:
			position, tokenIndex, depth = position397, tokenIndex397, depth397
			return false
		},
		/* 41 iriref <- <(iri / prefixedName)> */
		func() bool {
			position420, tokenIndex420, depth420 := position, tokenIndex, depth
			{
				position421 := position
				depth++
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l423
					}
					goto l422
				l423:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
					{
						position424 := position
						depth++
						{
							position427, tokenIndex427, depth427 := position, tokenIndex, depth
							{
								position428, tokenIndex428, depth428 := position, tokenIndex, depth
								if buffer[position] != rune(':') {
									goto l429
								}
								position++
								goto l428
							l429:
								position, tokenIndex, depth = position428, tokenIndex428, depth428
								if buffer[position] != rune(' ') {
									goto l427
								}
								position++
							}
						l428:
							goto l420
						l427:
							position, tokenIndex, depth = position427, tokenIndex427, depth427
						}
						if !matchDot() {
							goto l420
						}
					l425:
						{
							position426, tokenIndex426, depth426 := position, tokenIndex, depth
							{
								position430, tokenIndex430, depth430 := position, tokenIndex, depth
								{
									position431, tokenIndex431, depth431 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l432
									}
									position++
									goto l431
								l432:
									position, tokenIndex, depth = position431, tokenIndex431, depth431
									if buffer[position] != rune(' ') {
										goto l430
									}
									position++
								}
							l431:
								goto l426
							l430:
								position, tokenIndex, depth = position430, tokenIndex430, depth430
							}
							if !matchDot() {
								goto l426
							}
							goto l425
						l426:
							position, tokenIndex, depth = position426, tokenIndex426, depth426
						}
						if buffer[position] != rune(':') {
							goto l420
						}
						position++
					l433:
						{
							position434, tokenIndex434, depth434 := position, tokenIndex, depth
							{
								position435, tokenIndex435, depth435 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l436
								}
								position++
								goto l435
							l436:
								position, tokenIndex, depth = position435, tokenIndex435, depth435
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l437
								}
								position++
								goto l435
							l437:
								position, tokenIndex, depth = position435, tokenIndex435, depth435
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l438
								}
								position++
								goto l435
							l438:
								position, tokenIndex, depth = position435, tokenIndex435, depth435
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l434
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l434
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l434
										}
										position++
										break
									}
								}

							}
						l435:
							goto l433
						l434:
							position, tokenIndex, depth = position434, tokenIndex434, depth434
						}
						if !rules[rulews]() {
							goto l420
						}
						depth--
						add(ruleprefixedName, position424)
					}
				}
			l422:
				depth--
				add(ruleiriref, position421)
			}
			return true
		l420:
			position, tokenIndex, depth = position420, tokenIndex420, depth420
			return false
		},
		/* 42 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				if buffer[position] != rune('<') {
					goto l440
				}
				position++
			l442:
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					{
						position444, tokenIndex444, depth444 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l444
						}
						position++
						goto l443
					l444:
						position, tokenIndex, depth = position444, tokenIndex444, depth444
					}
					if !matchDot() {
						goto l443
					}
					goto l442
				l443:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
				}
				if buffer[position] != rune('>') {
					goto l440
				}
				position++
				if !rules[rulews]() {
					goto l440
				}
				depth--
				add(ruleiri, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 43 prefixedName <- <((!(':' / ' ') .)+ ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* ws)> */
		nil,
		/* 44 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? ws)> */
		nil,
		/* 45 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 46 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? ws)> */
		nil,
		/* 47 booleanLiteral <- <(TRUE / FALSE)> */
		nil,
		/* 48 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 49 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? ws)> */
		nil,
		/* 50 anon <- <('[' ws ']' ws)> */
		nil,
		/* 51 nil <- <('(' ws ')' ws)> */
		nil,
		/* 52 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 53 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 54 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 55 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') ws)> */
		nil,
		/* 56 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') ws)> */
		nil,
		/* 57 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 58 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 59 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 60 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 61 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 62 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') ws)> */
		nil,
		/* 63 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 64 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') ws)> */
		nil,
		/* 65 LBRACE <- <('{' ws)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				if buffer[position] != rune('{') {
					goto l467
				}
				position++
				if !rules[rulews]() {
					goto l467
				}
				depth--
				add(ruleLBRACE, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 66 RBRACE <- <('}' ws)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				if buffer[position] != rune('}') {
					goto l469
				}
				position++
				if !rules[rulews]() {
					goto l469
				}
				depth--
				add(ruleRBRACE, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 67 LBRACK <- <('[' ws)> */
		nil,
		/* 68 RBRACK <- <(']' ws)> */
		nil,
		/* 69 SEMICOLON <- <(';' ws)> */
		nil,
		/* 70 COMMA <- <(',' ws)> */
		nil,
		/* 71 DOT <- <('.' ws)> */
		func() bool {
			position475, tokenIndex475, depth475 := position, tokenIndex, depth
			{
				position476 := position
				depth++
				if buffer[position] != rune('.') {
					goto l475
				}
				position++
				if !rules[rulews]() {
					goto l475
				}
				depth--
				add(ruleDOT, position476)
			}
			return true
		l475:
			position, tokenIndex, depth = position475, tokenIndex475, depth475
			return false
		},
		/* 72 COLON <- <(':' ws)> */
		nil,
		/* 73 PIPE <- <('|' ws)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				if buffer[position] != rune('|') {
					goto l478
				}
				position++
				if !rules[rulews]() {
					goto l478
				}
				depth--
				add(rulePIPE, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 74 SLASH <- <('/' ws)> */
		nil,
		/* 75 INVERSE <- <('^' ws)> */
		func() bool {
			position481, tokenIndex481, depth481 := position, tokenIndex, depth
			{
				position482 := position
				depth++
				if buffer[position] != rune('^') {
					goto l481
				}
				position++
				if !rules[rulews]() {
					goto l481
				}
				depth--
				add(ruleINVERSE, position482)
			}
			return true
		l481:
			position, tokenIndex, depth = position481, tokenIndex481, depth481
			return false
		},
		/* 76 LPAREN <- <('(' ws)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				if buffer[position] != rune('(') {
					goto l483
				}
				position++
				if !rules[rulews]() {
					goto l483
				}
				depth--
				add(ruleLPAREN, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 77 RPAREN <- <(')' ws)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				if buffer[position] != rune(')') {
					goto l485
				}
				position++
				if !rules[rulews]() {
					goto l485
				}
				depth--
				add(ruleRPAREN, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 78 ISA <- <('a' ws)> */
		func() bool {
			position487, tokenIndex487, depth487 := position, tokenIndex, depth
			{
				position488 := position
				depth++
				if buffer[position] != rune('a') {
					goto l487
				}
				position++
				if !rules[rulews]() {
					goto l487
				}
				depth--
				add(ruleISA, position488)
			}
			return true
		l487:
			position, tokenIndex, depth = position487, tokenIndex487, depth487
			return false
		},
		/* 79 NOT <- <('!' ws)> */
		nil,
		/* 80 STAR <- <('*' ws)> */
		nil,
		/* 81 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') ws)> */
		nil,
		/* 82 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') ws)> */
		nil,
		/* 83 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') ws)> */
		nil,
		/* 84 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') ws)> */
		nil,
		/* 85 INTEGER <- <([0-9]+ ws)> */
		func() bool {
			position495, tokenIndex495, depth495 := position, tokenIndex, depth
			{
				position496 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l495
				}
				position++
			l497:
				{
					position498, tokenIndex498, depth498 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l498
					}
					position++
					goto l497
				l498:
					position, tokenIndex, depth = position498, tokenIndex498, depth498
				}
				if !rules[rulews]() {
					goto l495
				}
				depth--
				add(ruleINTEGER, position496)
			}
			return true
		l495:
			position, tokenIndex, depth = position495, tokenIndex495, depth495
			return false
		},
		/* 86 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position500 := position
				depth++
			l501:
				{
					position502, tokenIndex502, depth502 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\v':
							if buffer[position] != rune('\v') {
								goto l502
							}
							position++
							break
						case '\f':
							if buffer[position] != rune('\f') {
								goto l502
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l502
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l502
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l502
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l502
							}
							position++
							break
						}
					}

					goto l501
				l502:
					position, tokenIndex, depth = position502, tokenIndex502, depth502
				}
				depth--
				add(rulews, position500)
			}
			return true
		},
	}
	p.rules = rules
}
