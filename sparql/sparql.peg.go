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
							l21:
								{
									position22, tokenIndex22, depth22 := position, tokenIndex, depth
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
										goto l22
									l23:
										position, tokenIndex, depth = position23, tokenIndex23, depth23
									}
									if !matchDot() {
										goto l22
									}
									goto l21
								l22:
									position, tokenIndex, depth = position22, tokenIndex22, depth22
								}
								{
									position26 := position
									depth++
									if buffer[position] != rune(':') {
										goto l6
									}
									position++
									if !rules[rulews]() {
										goto l6
									}
									depth--
									add(ruleCOLON, position26)
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
								position27 := position
								depth++
								{
									position28 := position
									depth++
									{
										position29, tokenIndex29, depth29 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l30
										}
										position++
										goto l29
									l30:
										position, tokenIndex, depth = position29, tokenIndex29, depth29
										if buffer[position] != rune('B') {
											goto l4
										}
										position++
									}
								l29:
									{
										position31, tokenIndex31, depth31 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l32
										}
										position++
										goto l31
									l32:
										position, tokenIndex, depth = position31, tokenIndex31, depth31
										if buffer[position] != rune('A') {
											goto l4
										}
										position++
									}
								l31:
									{
										position33, tokenIndex33, depth33 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l34
										}
										position++
										goto l33
									l34:
										position, tokenIndex, depth = position33, tokenIndex33, depth33
										if buffer[position] != rune('S') {
											goto l4
										}
										position++
									}
								l33:
									{
										position35, tokenIndex35, depth35 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l36
										}
										position++
										goto l35
									l36:
										position, tokenIndex, depth = position35, tokenIndex35, depth35
										if buffer[position] != rune('E') {
											goto l4
										}
										position++
									}
								l35:
									if !rules[rulews]() {
										goto l4
									}
									depth--
									add(ruleBASE, position28)
								}
								if !rules[ruleiri]() {
									goto l4
								}
								depth--
								add(rulebaseDecl, position27)
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
					position37 := position
					depth++
					{
						position38 := position
						depth++
						if !rules[ruleselect]() {
							goto l0
						}
						{
							position39, tokenIndex39, depth39 := position, tokenIndex, depth
							{
								position41 := position
								depth++
								{
									position42 := position
									depth++
									{
										position43, tokenIndex43, depth43 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l44
										}
										position++
										goto l43
									l44:
										position, tokenIndex, depth = position43, tokenIndex43, depth43
										if buffer[position] != rune('F') {
											goto l39
										}
										position++
									}
								l43:
									{
										position45, tokenIndex45, depth45 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l46
										}
										position++
										goto l45
									l46:
										position, tokenIndex, depth = position45, tokenIndex45, depth45
										if buffer[position] != rune('R') {
											goto l39
										}
										position++
									}
								l45:
									{
										position47, tokenIndex47, depth47 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l48
										}
										position++
										goto l47
									l48:
										position, tokenIndex, depth = position47, tokenIndex47, depth47
										if buffer[position] != rune('O') {
											goto l39
										}
										position++
									}
								l47:
									{
										position49, tokenIndex49, depth49 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l50
										}
										position++
										goto l49
									l50:
										position, tokenIndex, depth = position49, tokenIndex49, depth49
										if buffer[position] != rune('M') {
											goto l39
										}
										position++
									}
								l49:
									if !rules[rulews]() {
										goto l39
									}
									depth--
									add(ruleFROM, position42)
								}
								{
									position51, tokenIndex51, depth51 := position, tokenIndex, depth
									{
										position53 := position
										depth++
										{
											position54, tokenIndex54, depth54 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l55
											}
											position++
											goto l54
										l55:
											position, tokenIndex, depth = position54, tokenIndex54, depth54
											if buffer[position] != rune('N') {
												goto l51
											}
											position++
										}
									l54:
										{
											position56, tokenIndex56, depth56 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l57
											}
											position++
											goto l56
										l57:
											position, tokenIndex, depth = position56, tokenIndex56, depth56
											if buffer[position] != rune('A') {
												goto l51
											}
											position++
										}
									l56:
										{
											position58, tokenIndex58, depth58 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l59
											}
											position++
											goto l58
										l59:
											position, tokenIndex, depth = position58, tokenIndex58, depth58
											if buffer[position] != rune('M') {
												goto l51
											}
											position++
										}
									l58:
										{
											position60, tokenIndex60, depth60 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l61
											}
											position++
											goto l60
										l61:
											position, tokenIndex, depth = position60, tokenIndex60, depth60
											if buffer[position] != rune('E') {
												goto l51
											}
											position++
										}
									l60:
										{
											position62, tokenIndex62, depth62 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l63
											}
											position++
											goto l62
										l63:
											position, tokenIndex, depth = position62, tokenIndex62, depth62
											if buffer[position] != rune('D') {
												goto l51
											}
											position++
										}
									l62:
										if !rules[rulews]() {
											goto l51
										}
										depth--
										add(ruleNAMED, position53)
									}
									goto l52
								l51:
									position, tokenIndex, depth = position51, tokenIndex51, depth51
								}
							l52:
								if !rules[ruleiriref]() {
									goto l39
								}
								depth--
								add(ruledatasetClause, position41)
							}
							goto l40
						l39:
							position, tokenIndex, depth = position39, tokenIndex39, depth39
						}
					l40:
						if !rules[rulewhereClause]() {
							goto l0
						}
						{
							position64 := position
							depth++
							{
								position65, tokenIndex65, depth65 := position, tokenIndex, depth
								{
									position67 := position
									depth++
									{
										position68, tokenIndex68, depth68 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l69
										}
										{
											position70, tokenIndex70, depth70 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l70
											}
											goto l71
										l70:
											position, tokenIndex, depth = position70, tokenIndex70, depth70
										}
									l71:
										goto l68
									l69:
										position, tokenIndex, depth = position68, tokenIndex68, depth68
										if !rules[ruleoffset]() {
											goto l65
										}
										{
											position72, tokenIndex72, depth72 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l72
											}
											goto l73
										l72:
											position, tokenIndex, depth = position72, tokenIndex72, depth72
										}
									l73:
									}
								l68:
									depth--
									add(rulelimitOffsetClauses, position67)
								}
								goto l66
							l65:
								position, tokenIndex, depth = position65, tokenIndex65, depth65
							}
						l66:
							depth--
							add(rulesolutionModifier, position64)
						}
						depth--
						add(ruleselectQuery, position38)
					}
					depth--
					add(rulequery, position37)
				}
				{
					position74, tokenIndex74, depth74 := position, tokenIndex, depth
					if !matchDot() {
						goto l74
					}
					goto l0
				l74:
					position, tokenIndex, depth = position74, tokenIndex74, depth74
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
		/* 2 prefixDecl <- <(PREFIX (!(':' / ' ') .)* COLON iri)> */
		nil,
		/* 3 baseDecl <- <(BASE iri)> */
		nil,
		/* 4 query <- <selectQuery> */
		nil,
		/* 5 selectQuery <- <(select datasetClause? whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position80, tokenIndex80, depth80 := position, tokenIndex, depth
			{
				position81 := position
				depth++
				{
					position82 := position
					depth++
					{
						position83, tokenIndex83, depth83 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l84
						}
						position++
						goto l83
					l84:
						position, tokenIndex, depth = position83, tokenIndex83, depth83
						if buffer[position] != rune('S') {
							goto l80
						}
						position++
					}
				l83:
					{
						position85, tokenIndex85, depth85 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l86
						}
						position++
						goto l85
					l86:
						position, tokenIndex, depth = position85, tokenIndex85, depth85
						if buffer[position] != rune('E') {
							goto l80
						}
						position++
					}
				l85:
					{
						position87, tokenIndex87, depth87 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l88
						}
						position++
						goto l87
					l88:
						position, tokenIndex, depth = position87, tokenIndex87, depth87
						if buffer[position] != rune('L') {
							goto l80
						}
						position++
					}
				l87:
					{
						position89, tokenIndex89, depth89 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l90
						}
						position++
						goto l89
					l90:
						position, tokenIndex, depth = position89, tokenIndex89, depth89
						if buffer[position] != rune('E') {
							goto l80
						}
						position++
					}
				l89:
					{
						position91, tokenIndex91, depth91 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l92
						}
						position++
						goto l91
					l92:
						position, tokenIndex, depth = position91, tokenIndex91, depth91
						if buffer[position] != rune('C') {
							goto l80
						}
						position++
					}
				l91:
					{
						position93, tokenIndex93, depth93 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l94
						}
						position++
						goto l93
					l94:
						position, tokenIndex, depth = position93, tokenIndex93, depth93
						if buffer[position] != rune('T') {
							goto l80
						}
						position++
					}
				l93:
					if !rules[rulews]() {
						goto l80
					}
					depth--
					add(ruleSELECT, position82)
				}
				{
					position95, tokenIndex95, depth95 := position, tokenIndex, depth
					{
						position97, tokenIndex97, depth97 := position, tokenIndex, depth
						{
							position99 := position
							depth++
							{
								position100, tokenIndex100, depth100 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l101
								}
								position++
								goto l100
							l101:
								position, tokenIndex, depth = position100, tokenIndex100, depth100
								if buffer[position] != rune('D') {
									goto l98
								}
								position++
							}
						l100:
							{
								position102, tokenIndex102, depth102 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l103
								}
								position++
								goto l102
							l103:
								position, tokenIndex, depth = position102, tokenIndex102, depth102
								if buffer[position] != rune('I') {
									goto l98
								}
								position++
							}
						l102:
							{
								position104, tokenIndex104, depth104 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l105
								}
								position++
								goto l104
							l105:
								position, tokenIndex, depth = position104, tokenIndex104, depth104
								if buffer[position] != rune('S') {
									goto l98
								}
								position++
							}
						l104:
							{
								position106, tokenIndex106, depth106 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l107
								}
								position++
								goto l106
							l107:
								position, tokenIndex, depth = position106, tokenIndex106, depth106
								if buffer[position] != rune('T') {
									goto l98
								}
								position++
							}
						l106:
							{
								position108, tokenIndex108, depth108 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l109
								}
								position++
								goto l108
							l109:
								position, tokenIndex, depth = position108, tokenIndex108, depth108
								if buffer[position] != rune('I') {
									goto l98
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
									goto l98
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
									goto l98
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
									goto l98
								}
								position++
							}
						l114:
							if !rules[rulews]() {
								goto l98
							}
							depth--
							add(ruleDISTINCT, position99)
						}
						goto l97
					l98:
						position, tokenIndex, depth = position97, tokenIndex97, depth97
						{
							position116 := position
							depth++
							{
								position117, tokenIndex117, depth117 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l118
								}
								position++
								goto l117
							l118:
								position, tokenIndex, depth = position117, tokenIndex117, depth117
								if buffer[position] != rune('R') {
									goto l95
								}
								position++
							}
						l117:
							{
								position119, tokenIndex119, depth119 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l120
								}
								position++
								goto l119
							l120:
								position, tokenIndex, depth = position119, tokenIndex119, depth119
								if buffer[position] != rune('E') {
									goto l95
								}
								position++
							}
						l119:
							{
								position121, tokenIndex121, depth121 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l122
								}
								position++
								goto l121
							l122:
								position, tokenIndex, depth = position121, tokenIndex121, depth121
								if buffer[position] != rune('D') {
									goto l95
								}
								position++
							}
						l121:
							{
								position123, tokenIndex123, depth123 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l124
								}
								position++
								goto l123
							l124:
								position, tokenIndex, depth = position123, tokenIndex123, depth123
								if buffer[position] != rune('U') {
									goto l95
								}
								position++
							}
						l123:
							{
								position125, tokenIndex125, depth125 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l126
								}
								position++
								goto l125
							l126:
								position, tokenIndex, depth = position125, tokenIndex125, depth125
								if buffer[position] != rune('C') {
									goto l95
								}
								position++
							}
						l125:
							{
								position127, tokenIndex127, depth127 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l128
								}
								position++
								goto l127
							l128:
								position, tokenIndex, depth = position127, tokenIndex127, depth127
								if buffer[position] != rune('E') {
									goto l95
								}
								position++
							}
						l127:
							{
								position129, tokenIndex129, depth129 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l130
								}
								position++
								goto l129
							l130:
								position, tokenIndex, depth = position129, tokenIndex129, depth129
								if buffer[position] != rune('D') {
									goto l95
								}
								position++
							}
						l129:
							if !rules[rulews]() {
								goto l95
							}
							depth--
							add(ruleREDUCED, position116)
						}
					}
				l97:
					goto l96
				l95:
					position, tokenIndex, depth = position95, tokenIndex95, depth95
				}
			l96:
				{
					position131, tokenIndex131, depth131 := position, tokenIndex, depth
					{
						position133 := position
						depth++
						if buffer[position] != rune('*') {
							goto l132
						}
						position++
						if !rules[rulews]() {
							goto l132
						}
						depth--
						add(ruleSTAR, position133)
					}
					goto l131
				l132:
					position, tokenIndex, depth = position131, tokenIndex131, depth131
					{
						position136 := position
						depth++
						if !rules[rulevar]() {
							goto l80
						}
						depth--
						add(ruleprojectionElem, position136)
					}
				l134:
					{
						position135, tokenIndex135, depth135 := position, tokenIndex, depth
						{
							position137 := position
							depth++
							if !rules[rulevar]() {
								goto l135
							}
							depth--
							add(ruleprojectionElem, position137)
						}
						goto l134
					l135:
						position, tokenIndex, depth = position135, tokenIndex135, depth135
					}
				}
			l131:
				depth--
				add(ruleselect, position81)
			}
			return true
		l80:
			position, tokenIndex, depth = position80, tokenIndex80, depth80
			return false
		},
		/* 7 subSelect <- <(select whereClause)> */
		func() bool {
			position138, tokenIndex138, depth138 := position, tokenIndex, depth
			{
				position139 := position
				depth++
				if !rules[ruleselect]() {
					goto l138
				}
				if !rules[rulewhereClause]() {
					goto l138
				}
				depth--
				add(rulesubSelect, position139)
			}
			return true
		l138:
			position, tokenIndex, depth = position138, tokenIndex138, depth138
			return false
		},
		/* 8 projectionElem <- <var> */
		nil,
		/* 9 datasetClause <- <(FROM NAMED? iriref)> */
		nil,
		/* 10 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position142, tokenIndex142, depth142 := position, tokenIndex, depth
			{
				position143 := position
				depth++
				{
					position144, tokenIndex144, depth144 := position, tokenIndex, depth
					{
						position146 := position
						depth++
						{
							position147, tokenIndex147, depth147 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l148
							}
							position++
							goto l147
						l148:
							position, tokenIndex, depth = position147, tokenIndex147, depth147
							if buffer[position] != rune('W') {
								goto l144
							}
							position++
						}
					l147:
						{
							position149, tokenIndex149, depth149 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l150
							}
							position++
							goto l149
						l150:
							position, tokenIndex, depth = position149, tokenIndex149, depth149
							if buffer[position] != rune('H') {
								goto l144
							}
							position++
						}
					l149:
						{
							position151, tokenIndex151, depth151 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l152
							}
							position++
							goto l151
						l152:
							position, tokenIndex, depth = position151, tokenIndex151, depth151
							if buffer[position] != rune('E') {
								goto l144
							}
							position++
						}
					l151:
						{
							position153, tokenIndex153, depth153 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l154
							}
							position++
							goto l153
						l154:
							position, tokenIndex, depth = position153, tokenIndex153, depth153
							if buffer[position] != rune('R') {
								goto l144
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
								goto l144
							}
							position++
						}
					l155:
						if !rules[rulews]() {
							goto l144
						}
						depth--
						add(ruleWHERE, position146)
					}
					goto l145
				l144:
					position, tokenIndex, depth = position144, tokenIndex144, depth144
				}
			l145:
				if !rules[rulegroupGraphPattern]() {
					goto l142
				}
				depth--
				add(rulewhereClause, position143)
			}
			return true
		l142:
			position, tokenIndex, depth = position142, tokenIndex142, depth142
			return false
		},
		/* 11 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position157, tokenIndex157, depth157 := position, tokenIndex, depth
			{
				position158 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l157
				}
				{
					position159, tokenIndex159, depth159 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l160
					}
					goto l159
				l160:
					position, tokenIndex, depth = position159, tokenIndex159, depth159
					if !rules[rulegraphPattern]() {
						goto l157
					}
				}
			l159:
				if !rules[ruleRBRACE]() {
					goto l157
				}
				depth--
				add(rulegroupGraphPattern, position158)
			}
			return true
		l157:
			position, tokenIndex, depth = position157, tokenIndex157, depth157
			return false
		},
		/* 12 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position162 := position
				depth++
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					{
						position165 := position
						depth++
						{
							position166 := position
							depth++
							if !rules[ruletriplesSameSubjectPath]() {
								goto l163
							}
						l167:
							{
								position168, tokenIndex168, depth168 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l168
								}
								if !rules[ruletriplesSameSubjectPath]() {
									goto l168
								}
								goto l167
							l168:
								position, tokenIndex, depth = position168, tokenIndex168, depth168
							}
							{
								position169, tokenIndex169, depth169 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l169
								}
								goto l170
							l169:
								position, tokenIndex, depth = position169, tokenIndex169, depth169
							}
						l170:
							depth--
							add(ruletriplesBlock, position166)
						}
						depth--
						add(rulebasicGraphPattern, position165)
					}
					goto l164
				l163:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
				}
			l164:
				{
					position171, tokenIndex171, depth171 := position, tokenIndex, depth
					{
						position173 := position
						depth++
						{
							position174, tokenIndex174, depth174 := position, tokenIndex, depth
							{
								position176 := position
								depth++
								{
									position177 := position
									depth++
									{
										position178, tokenIndex178, depth178 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l179
										}
										position++
										goto l178
									l179:
										position, tokenIndex, depth = position178, tokenIndex178, depth178
										if buffer[position] != rune('O') {
											goto l175
										}
										position++
									}
								l178:
									{
										position180, tokenIndex180, depth180 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l181
										}
										position++
										goto l180
									l181:
										position, tokenIndex, depth = position180, tokenIndex180, depth180
										if buffer[position] != rune('P') {
											goto l175
										}
										position++
									}
								l180:
									{
										position182, tokenIndex182, depth182 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l183
										}
										position++
										goto l182
									l183:
										position, tokenIndex, depth = position182, tokenIndex182, depth182
										if buffer[position] != rune('T') {
											goto l175
										}
										position++
									}
								l182:
									{
										position184, tokenIndex184, depth184 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l185
										}
										position++
										goto l184
									l185:
										position, tokenIndex, depth = position184, tokenIndex184, depth184
										if buffer[position] != rune('I') {
											goto l175
										}
										position++
									}
								l184:
									{
										position186, tokenIndex186, depth186 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l187
										}
										position++
										goto l186
									l187:
										position, tokenIndex, depth = position186, tokenIndex186, depth186
										if buffer[position] != rune('O') {
											goto l175
										}
										position++
									}
								l186:
									{
										position188, tokenIndex188, depth188 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l189
										}
										position++
										goto l188
									l189:
										position, tokenIndex, depth = position188, tokenIndex188, depth188
										if buffer[position] != rune('N') {
											goto l175
										}
										position++
									}
								l188:
									{
										position190, tokenIndex190, depth190 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l191
										}
										position++
										goto l190
									l191:
										position, tokenIndex, depth = position190, tokenIndex190, depth190
										if buffer[position] != rune('A') {
											goto l175
										}
										position++
									}
								l190:
									{
										position192, tokenIndex192, depth192 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l193
										}
										position++
										goto l192
									l193:
										position, tokenIndex, depth = position192, tokenIndex192, depth192
										if buffer[position] != rune('L') {
											goto l175
										}
										position++
									}
								l192:
									if !rules[rulews]() {
										goto l175
									}
									depth--
									add(ruleOPTIONAL, position177)
								}
								if !rules[ruleLBRACE]() {
									goto l175
								}
								{
									position194, tokenIndex194, depth194 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l195
									}
									goto l194
								l195:
									position, tokenIndex, depth = position194, tokenIndex194, depth194
									if !rules[rulegraphPattern]() {
										goto l175
									}
								}
							l194:
								if !rules[ruleRBRACE]() {
									goto l175
								}
								depth--
								add(ruleoptionalGraphPattern, position176)
							}
							goto l174
						l175:
							position, tokenIndex, depth = position174, tokenIndex174, depth174
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l171
							}
						}
					l174:
						depth--
						add(rulegraphPatternNotTriples, position173)
					}
					{
						position196, tokenIndex196, depth196 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l196
						}
						goto l197
					l196:
						position, tokenIndex, depth = position196, tokenIndex196, depth196
					}
				l197:
					if !rules[rulegraphPattern]() {
						goto l171
					}
					goto l172
				l171:
					position, tokenIndex, depth = position171, tokenIndex171, depth171
				}
			l172:
				depth--
				add(rulegraphPattern, position162)
			}
			return true
		},
		/* 13 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 14 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 15 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l200
				}
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					{
						position204 := position
						depth++
						{
							position205, tokenIndex205, depth205 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l206
							}
							position++
							goto l205
						l206:
							position, tokenIndex, depth = position205, tokenIndex205, depth205
							if buffer[position] != rune('U') {
								goto l202
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
								goto l202
							}
							position++
						}
					l207:
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('I') {
								goto l202
							}
							position++
						}
					l209:
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l212
							}
							position++
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if buffer[position] != rune('O') {
								goto l202
							}
							position++
						}
					l211:
						{
							position213, tokenIndex213, depth213 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l214
							}
							position++
							goto l213
						l214:
							position, tokenIndex, depth = position213, tokenIndex213, depth213
							if buffer[position] != rune('N') {
								goto l202
							}
							position++
						}
					l213:
						if !rules[rulews]() {
							goto l202
						}
						depth--
						add(ruleUNION, position204)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l202
					}
					goto l203
				l202:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
				}
			l203:
				depth--
				add(rulegroupOrUnionGraphPattern, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 16 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 17 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		nil,
		/* 18 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position217, tokenIndex217, depth217 := position, tokenIndex, depth
			{
				position218 := position
				depth++
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l220
					}
					if !rules[rulepropertyListPath]() {
						goto l220
					}
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					{
						position221 := position
						depth++
						{
							position222, tokenIndex222, depth222 := position, tokenIndex, depth
							{
								position224 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l223
								}
								if !rules[rulegraphNodePath]() {
									goto l223
								}
							l225:
								{
									position226, tokenIndex226, depth226 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l226
									}
									goto l225
								l226:
									position, tokenIndex, depth = position226, tokenIndex226, depth226
								}
								if !rules[ruleRPAREN]() {
									goto l223
								}
								depth--
								add(rulecollectionPath, position224)
							}
							goto l222
						l223:
							position, tokenIndex, depth = position222, tokenIndex222, depth222
							{
								position227 := position
								depth++
								{
									position228 := position
									depth++
									if buffer[position] != rune('[') {
										goto l217
									}
									position++
									if !rules[rulews]() {
										goto l217
									}
									depth--
									add(ruleLBRACK, position228)
								}
								if !rules[rulepropertyListPath]() {
									goto l217
								}
								{
									position229 := position
									depth++
									if buffer[position] != rune(']') {
										goto l217
									}
									position++
									if !rules[rulews]() {
										goto l217
									}
									depth--
									add(ruleRBRACK, position229)
								}
								depth--
								add(ruleblankNodePropertyListPath, position227)
							}
						}
					l222:
						depth--
						add(ruletriplesNodePath, position221)
					}
					if !rules[rulepropertyListPath]() {
						goto l217
					}
				}
			l219:
				depth--
				add(ruletriplesSameSubjectPath, position218)
			}
			return true
		l217:
			position, tokenIndex, depth = position217, tokenIndex217, depth217
			return false
		},
		/* 19 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position230, tokenIndex230, depth230 := position, tokenIndex, depth
			{
				position231 := position
				depth++
				{
					position232, tokenIndex232, depth232 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l233
					}
					goto l232
				l233:
					position, tokenIndex, depth = position232, tokenIndex232, depth232
					{
						position234 := position
						depth++
						{
							position235, tokenIndex235, depth235 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l236
							}
							goto l235
						l236:
							position, tokenIndex, depth = position235, tokenIndex235, depth235
							{
								switch buffer[position] {
								case '(':
									{
										position238 := position
										depth++
										if buffer[position] != rune('(') {
											goto l230
										}
										position++
										if !rules[rulews]() {
											goto l230
										}
										if buffer[position] != rune(')') {
											goto l230
										}
										position++
										if !rules[rulews]() {
											goto l230
										}
										depth--
										add(rulenil, position238)
									}
									break
								case '[', '_':
									{
										position239 := position
										depth++
										{
											position240, tokenIndex240, depth240 := position, tokenIndex, depth
											{
												position242 := position
												depth++
												if buffer[position] != rune('_') {
													goto l241
												}
												position++
												if buffer[position] != rune(':') {
													goto l241
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l241
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l241
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l241
														}
														position++
														break
													}
												}

												{
													position244, tokenIndex244, depth244 := position, tokenIndex, depth
													{
														position246, tokenIndex246, depth246 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l247
														}
														position++
														goto l246
													l247:
														position, tokenIndex, depth = position246, tokenIndex246, depth246
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l248
														}
														position++
														goto l246
													l248:
														position, tokenIndex, depth = position246, tokenIndex246, depth246
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l249
														}
														position++
														goto l246
													l249:
														position, tokenIndex, depth = position246, tokenIndex246, depth246
														if c := buffer[position]; c < rune('.') || c > rune('_') {
															goto l244
														}
														position++
													}
												l246:
													goto l245
												l244:
													position, tokenIndex, depth = position244, tokenIndex244, depth244
												}
											l245:
												if !rules[rulews]() {
													goto l241
												}
												depth--
												add(ruleblankNodeLabel, position242)
											}
											goto l240
										l241:
											position, tokenIndex, depth = position240, tokenIndex240, depth240
											{
												position250 := position
												depth++
												if buffer[position] != rune('[') {
													goto l230
												}
												position++
												if !rules[rulews]() {
													goto l230
												}
												if buffer[position] != rune(']') {
													goto l230
												}
												position++
												if !rules[rulews]() {
													goto l230
												}
												depth--
												add(ruleanon, position250)
											}
										}
									l240:
										depth--
										add(ruleblankNode, position239)
									}
									break
								case 'F', 'T', 'f', 't':
									{
										position251 := position
										depth++
										{
											position252, tokenIndex252, depth252 := position, tokenIndex, depth
											{
												position254 := position
												depth++
												{
													position255, tokenIndex255, depth255 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l256
													}
													position++
													goto l255
												l256:
													position, tokenIndex, depth = position255, tokenIndex255, depth255
													if buffer[position] != rune('T') {
														goto l253
													}
													position++
												}
											l255:
												{
													position257, tokenIndex257, depth257 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l258
													}
													position++
													goto l257
												l258:
													position, tokenIndex, depth = position257, tokenIndex257, depth257
													if buffer[position] != rune('R') {
														goto l253
													}
													position++
												}
											l257:
												{
													position259, tokenIndex259, depth259 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l260
													}
													position++
													goto l259
												l260:
													position, tokenIndex, depth = position259, tokenIndex259, depth259
													if buffer[position] != rune('U') {
														goto l253
													}
													position++
												}
											l259:
												{
													position261, tokenIndex261, depth261 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l262
													}
													position++
													goto l261
												l262:
													position, tokenIndex, depth = position261, tokenIndex261, depth261
													if buffer[position] != rune('E') {
														goto l253
													}
													position++
												}
											l261:
												if !rules[rulews]() {
													goto l253
												}
												depth--
												add(ruleTRUE, position254)
											}
											goto l252
										l253:
											position, tokenIndex, depth = position252, tokenIndex252, depth252
											{
												position263 := position
												depth++
												{
													position264, tokenIndex264, depth264 := position, tokenIndex, depth
													if buffer[position] != rune('f') {
														goto l265
													}
													position++
													goto l264
												l265:
													position, tokenIndex, depth = position264, tokenIndex264, depth264
													if buffer[position] != rune('F') {
														goto l230
													}
													position++
												}
											l264:
												{
													position266, tokenIndex266, depth266 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l267
													}
													position++
													goto l266
												l267:
													position, tokenIndex, depth = position266, tokenIndex266, depth266
													if buffer[position] != rune('A') {
														goto l230
													}
													position++
												}
											l266:
												{
													position268, tokenIndex268, depth268 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l269
													}
													position++
													goto l268
												l269:
													position, tokenIndex, depth = position268, tokenIndex268, depth268
													if buffer[position] != rune('L') {
														goto l230
													}
													position++
												}
											l268:
												{
													position270, tokenIndex270, depth270 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l271
													}
													position++
													goto l270
												l271:
													position, tokenIndex, depth = position270, tokenIndex270, depth270
													if buffer[position] != rune('S') {
														goto l230
													}
													position++
												}
											l270:
												{
													position272, tokenIndex272, depth272 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l273
													}
													position++
													goto l272
												l273:
													position, tokenIndex, depth = position272, tokenIndex272, depth272
													if buffer[position] != rune('E') {
														goto l230
													}
													position++
												}
											l272:
												if !rules[rulews]() {
													goto l230
												}
												depth--
												add(ruleFALSE, position263)
											}
										}
									l252:
										depth--
										add(rulebooleanLiteral, position251)
									}
									break
								case '"':
									{
										position274 := position
										depth++
										{
											position275 := position
											depth++
											if buffer[position] != rune('"') {
												goto l230
											}
											position++
										l276:
											{
												position277, tokenIndex277, depth277 := position, tokenIndex, depth
												{
													position278, tokenIndex278, depth278 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l278
													}
													position++
													goto l277
												l278:
													position, tokenIndex, depth = position278, tokenIndex278, depth278
												}
												if !matchDot() {
													goto l277
												}
												goto l276
											l277:
												position, tokenIndex, depth = position277, tokenIndex277, depth277
											}
											if buffer[position] != rune('"') {
												goto l230
											}
											position++
											depth--
											add(rulestring, position275)
										}
										{
											position279, tokenIndex279, depth279 := position, tokenIndex, depth
											{
												position281, tokenIndex281, depth281 := position, tokenIndex, depth
												if buffer[position] != rune('@') {
													goto l282
												}
												position++
												{
													position285, tokenIndex285, depth285 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l286
													}
													position++
													goto l285
												l286:
													position, tokenIndex, depth = position285, tokenIndex285, depth285
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l282
													}
													position++
												}
											l285:
											l283:
												{
													position284, tokenIndex284, depth284 := position, tokenIndex, depth
													{
														position287, tokenIndex287, depth287 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l288
														}
														position++
														goto l287
													l288:
														position, tokenIndex, depth = position287, tokenIndex287, depth287
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l284
														}
														position++
													}
												l287:
													goto l283
												l284:
													position, tokenIndex, depth = position284, tokenIndex284, depth284
												}
											l289:
												{
													position290, tokenIndex290, depth290 := position, tokenIndex, depth
													if buffer[position] != rune('-') {
														goto l290
													}
													position++
													{
														switch buffer[position] {
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l290
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l290
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l290
															}
															position++
															break
														}
													}

												l291:
													{
														position292, tokenIndex292, depth292 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l292
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l292
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l292
																}
																position++
																break
															}
														}

														goto l291
													l292:
														position, tokenIndex, depth = position292, tokenIndex292, depth292
													}
													goto l289
												l290:
													position, tokenIndex, depth = position290, tokenIndex290, depth290
												}
												goto l281
											l282:
												position, tokenIndex, depth = position281, tokenIndex281, depth281
												if buffer[position] != rune('^') {
													goto l279
												}
												position++
												if buffer[position] != rune('^') {
													goto l279
												}
												position++
												if !rules[ruleiriref]() {
													goto l279
												}
											}
										l281:
											goto l280
										l279:
											position, tokenIndex, depth = position279, tokenIndex279, depth279
										}
									l280:
										if !rules[rulews]() {
											goto l230
										}
										depth--
										add(ruleliteral, position274)
									}
									break
								default:
									{
										position295 := position
										depth++
										{
											position296, tokenIndex296, depth296 := position, tokenIndex, depth
											{
												position298, tokenIndex298, depth298 := position, tokenIndex, depth
												if buffer[position] != rune('+') {
													goto l299
												}
												position++
												goto l298
											l299:
												position, tokenIndex, depth = position298, tokenIndex298, depth298
												if buffer[position] != rune('-') {
													goto l296
												}
												position++
											}
										l298:
											goto l297
										l296:
											position, tokenIndex, depth = position296, tokenIndex296, depth296
										}
									l297:
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l230
										}
										position++
									l300:
										{
											position301, tokenIndex301, depth301 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l301
											}
											position++
											goto l300
										l301:
											position, tokenIndex, depth = position301, tokenIndex301, depth301
										}
										{
											position302, tokenIndex302, depth302 := position, tokenIndex, depth
											if buffer[position] != rune('.') {
												goto l302
											}
											position++
										l304:
											{
												position305, tokenIndex305, depth305 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l305
												}
												position++
												goto l304
											l305:
												position, tokenIndex, depth = position305, tokenIndex305, depth305
											}
											goto l303
										l302:
											position, tokenIndex, depth = position302, tokenIndex302, depth302
										}
									l303:
										if !rules[rulews]() {
											goto l230
										}
										depth--
										add(rulenumericLiteral, position295)
									}
									break
								}
							}

						}
					l235:
						depth--
						add(rulegraphTerm, position234)
					}
				}
			l232:
				depth--
				add(rulevarOrTerm, position231)
			}
			return true
		l230:
			position, tokenIndex, depth = position230, tokenIndex230, depth230
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
			position310, tokenIndex310, depth310 := position, tokenIndex, depth
			{
				position311 := position
				depth++
				{
					position312, tokenIndex312, depth312 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l313
					}
					goto l312
				l313:
					position, tokenIndex, depth = position312, tokenIndex312, depth312
					{
						position314 := position
						depth++
						if !rules[rulepath]() {
							goto l310
						}
						depth--
						add(ruleverbPath, position314)
					}
				}
			l312:
				if !rules[ruleobjectListPath]() {
					goto l310
				}
				{
					position315, tokenIndex315, depth315 := position, tokenIndex, depth
					{
						position317 := position
						depth++
						if buffer[position] != rune(';') {
							goto l315
						}
						position++
						if !rules[rulews]() {
							goto l315
						}
						depth--
						add(ruleSEMICOLON, position317)
					}
					if !rules[rulepropertyListPath]() {
						goto l315
					}
					goto l316
				l315:
					position, tokenIndex, depth = position315, tokenIndex315, depth315
				}
			l316:
				depth--
				add(rulepropertyListPath, position311)
			}
			return true
		l310:
			position, tokenIndex, depth = position310, tokenIndex310, depth310
			return false
		},
		/* 25 verbPath <- <path> */
		nil,
		/* 26 path <- <pathAlternative> */
		func() bool {
			position319, tokenIndex319, depth319 := position, tokenIndex, depth
			{
				position320 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l319
				}
				depth--
				add(rulepath, position320)
			}
			return true
		l319:
			position, tokenIndex, depth = position319, tokenIndex319, depth319
			return false
		},
		/* 27 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position321, tokenIndex321, depth321 := position, tokenIndex, depth
			{
				position322 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l321
				}
			l323:
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l324
					}
					if !rules[rulepathAlternative]() {
						goto l324
					}
					goto l323
				l324:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
				}
				depth--
				add(rulepathAlternative, position322)
			}
			return true
		l321:
			position, tokenIndex, depth = position321, tokenIndex321, depth321
			return false
		},
		/* 28 pathSequence <- <(pathElt (SLASH pathSequence)*)> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				{
					position327 := position
					depth++
					{
						position328, tokenIndex328, depth328 := position, tokenIndex, depth
						if !rules[ruleINVERSE]() {
							goto l328
						}
						goto l329
					l328:
						position, tokenIndex, depth = position328, tokenIndex328, depth328
					}
				l329:
					{
						position330 := position
						depth++
						{
							position331, tokenIndex331, depth331 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l332
							}
							goto l331
						l332:
							position, tokenIndex, depth = position331, tokenIndex331, depth331
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l325
									}
									if !rules[rulepath]() {
										goto l325
									}
									if !rules[ruleRPAREN]() {
										goto l325
									}
									break
								case '!':
									{
										position334 := position
										depth++
										if buffer[position] != rune('!') {
											goto l325
										}
										position++
										if !rules[rulews]() {
											goto l325
										}
										depth--
										add(ruleNOT, position334)
									}
									{
										position335 := position
										depth++
										{
											position336, tokenIndex336, depth336 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l337
											}
											goto l336
										l337:
											position, tokenIndex, depth = position336, tokenIndex336, depth336
											if !rules[ruleLPAREN]() {
												goto l325
											}
											{
												position338, tokenIndex338, depth338 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l338
												}
											l340:
												{
													position341, tokenIndex341, depth341 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l341
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l341
													}
													goto l340
												l341:
													position, tokenIndex, depth = position341, tokenIndex341, depth341
												}
												goto l339
											l338:
												position, tokenIndex, depth = position338, tokenIndex338, depth338
											}
										l339:
											if !rules[ruleRPAREN]() {
												goto l325
											}
										}
									l336:
										depth--
										add(rulepathNegatedPropertySet, position335)
									}
									break
								default:
									if !rules[ruleISA]() {
										goto l325
									}
									break
								}
							}

						}
					l331:
						depth--
						add(rulepathPrimary, position330)
					}
					depth--
					add(rulepathElt, position327)
				}
			l342:
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					{
						position344 := position
						depth++
						if buffer[position] != rune('/') {
							goto l343
						}
						position++
						if !rules[rulews]() {
							goto l343
						}
						depth--
						add(ruleSLASH, position344)
					}
					if !rules[rulepathSequence]() {
						goto l343
					}
					goto l342
				l343:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
				}
				depth--
				add(rulepathSequence, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
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
			position348, tokenIndex348, depth348 := position, tokenIndex, depth
			{
				position349 := position
				depth++
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l351
					}
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if !rules[ruleISA]() {
						goto l352
					}
					goto l350
				l352:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					if !rules[ruleINVERSE]() {
						goto l348
					}
					{
						position353, tokenIndex353, depth353 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l354
						}
						goto l353
					l354:
						position, tokenIndex, depth = position353, tokenIndex353, depth353
						if !rules[ruleISA]() {
							goto l348
						}
					}
				l353:
				}
			l350:
				depth--
				add(rulepathOneInPropertySet, position349)
			}
			return true
		l348:
			position, tokenIndex, depth = position348, tokenIndex348, depth348
			return false
		},
		/* 33 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position355, tokenIndex355, depth355 := position, tokenIndex, depth
			{
				position356 := position
				depth++
				{
					position357 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l355
					}
					depth--
					add(ruleobjectPath, position357)
				}
			l358:
				{
					position359, tokenIndex359, depth359 := position, tokenIndex, depth
					{
						position360 := position
						depth++
						if buffer[position] != rune(',') {
							goto l359
						}
						position++
						if !rules[rulews]() {
							goto l359
						}
						depth--
						add(ruleCOMMA, position360)
					}
					if !rules[ruleobjectListPath]() {
						goto l359
					}
					goto l358
				l359:
					position, tokenIndex, depth = position359, tokenIndex359, depth359
				}
				depth--
				add(ruleobjectListPath, position356)
			}
			return true
		l355:
			position, tokenIndex, depth = position355, tokenIndex355, depth355
			return false
		},
		/* 34 objectPath <- <graphNodePath> */
		nil,
		/* 35 graphNodePath <- <varOrTerm> */
		func() bool {
			position362, tokenIndex362, depth362 := position, tokenIndex, depth
			{
				position363 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l362
				}
				depth--
				add(rulegraphNodePath, position363)
			}
			return true
		l362:
			position, tokenIndex, depth = position362, tokenIndex362, depth362
			return false
		},
		/* 36 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 37 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 38 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position366, tokenIndex366, depth366 := position, tokenIndex, depth
			{
				position367 := position
				depth++
				{
					position368 := position
					depth++
					{
						position369, tokenIndex369, depth369 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l370
						}
						position++
						goto l369
					l370:
						position, tokenIndex, depth = position369, tokenIndex369, depth369
						if buffer[position] != rune('L') {
							goto l366
						}
						position++
					}
				l369:
					{
						position371, tokenIndex371, depth371 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l372
						}
						position++
						goto l371
					l372:
						position, tokenIndex, depth = position371, tokenIndex371, depth371
						if buffer[position] != rune('I') {
							goto l366
						}
						position++
					}
				l371:
					{
						position373, tokenIndex373, depth373 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l374
						}
						position++
						goto l373
					l374:
						position, tokenIndex, depth = position373, tokenIndex373, depth373
						if buffer[position] != rune('M') {
							goto l366
						}
						position++
					}
				l373:
					{
						position375, tokenIndex375, depth375 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l376
						}
						position++
						goto l375
					l376:
						position, tokenIndex, depth = position375, tokenIndex375, depth375
						if buffer[position] != rune('I') {
							goto l366
						}
						position++
					}
				l375:
					{
						position377, tokenIndex377, depth377 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l378
						}
						position++
						goto l377
					l378:
						position, tokenIndex, depth = position377, tokenIndex377, depth377
						if buffer[position] != rune('T') {
							goto l366
						}
						position++
					}
				l377:
					if !rules[rulews]() {
						goto l366
					}
					depth--
					add(ruleLIMIT, position368)
				}
				if !rules[ruleINTEGER]() {
					goto l366
				}
				depth--
				add(rulelimit, position367)
			}
			return true
		l366:
			position, tokenIndex, depth = position366, tokenIndex366, depth366
			return false
		},
		/* 39 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position379, tokenIndex379, depth379 := position, tokenIndex, depth
			{
				position380 := position
				depth++
				{
					position381 := position
					depth++
					{
						position382, tokenIndex382, depth382 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l383
						}
						position++
						goto l382
					l383:
						position, tokenIndex, depth = position382, tokenIndex382, depth382
						if buffer[position] != rune('O') {
							goto l379
						}
						position++
					}
				l382:
					{
						position384, tokenIndex384, depth384 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l385
						}
						position++
						goto l384
					l385:
						position, tokenIndex, depth = position384, tokenIndex384, depth384
						if buffer[position] != rune('F') {
							goto l379
						}
						position++
					}
				l384:
					{
						position386, tokenIndex386, depth386 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l387
						}
						position++
						goto l386
					l387:
						position, tokenIndex, depth = position386, tokenIndex386, depth386
						if buffer[position] != rune('F') {
							goto l379
						}
						position++
					}
				l386:
					{
						position388, tokenIndex388, depth388 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l389
						}
						position++
						goto l388
					l389:
						position, tokenIndex, depth = position388, tokenIndex388, depth388
						if buffer[position] != rune('S') {
							goto l379
						}
						position++
					}
				l388:
					{
						position390, tokenIndex390, depth390 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l391
						}
						position++
						goto l390
					l391:
						position, tokenIndex, depth = position390, tokenIndex390, depth390
						if buffer[position] != rune('E') {
							goto l379
						}
						position++
					}
				l390:
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l393
						}
						position++
						goto l392
					l393:
						position, tokenIndex, depth = position392, tokenIndex392, depth392
						if buffer[position] != rune('T') {
							goto l379
						}
						position++
					}
				l392:
					if !rules[rulews]() {
						goto l379
					}
					depth--
					add(ruleOFFSET, position381)
				}
				if !rules[ruleINTEGER]() {
					goto l379
				}
				depth--
				add(ruleoffset, position380)
			}
			return true
		l379:
			position, tokenIndex, depth = position379, tokenIndex379, depth379
			return false
		},
		/* 40 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position394, tokenIndex394, depth394 := position, tokenIndex, depth
			{
				position395 := position
				depth++
				{
					position396, tokenIndex396, depth396 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l397
					}
					position++
					goto l396
				l397:
					position, tokenIndex, depth = position396, tokenIndex396, depth396
					if buffer[position] != rune('$') {
						goto l394
					}
					position++
				}
			l396:
				{
					position398 := position
					depth++
					{
						position401, tokenIndex401, depth401 := position, tokenIndex, depth
						{
							position403 := position
							depth++
							{
								position404, tokenIndex404, depth404 := position, tokenIndex, depth
								{
									position406 := position
									depth++
									{
										position407, tokenIndex407, depth407 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l408
										}
										position++
										goto l407
									l408:
										position, tokenIndex, depth = position407, tokenIndex407, depth407
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l405
										}
										position++
									}
								l407:
									depth--
									add(rulePN_CHARS_BASE, position406)
								}
								goto l404
							l405:
								position, tokenIndex, depth = position404, tokenIndex404, depth404
								if buffer[position] != rune('_') {
									goto l402
								}
								position++
							}
						l404:
							depth--
							add(rulePN_CHARS_U, position403)
						}
						goto l401
					l402:
						position, tokenIndex, depth = position401, tokenIndex401, depth401
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l394
						}
						position++
					}
				l401:
				l399:
					{
						position400, tokenIndex400, depth400 := position, tokenIndex, depth
						{
							position409, tokenIndex409, depth409 := position, tokenIndex, depth
							{
								position411 := position
								depth++
								{
									position412, tokenIndex412, depth412 := position, tokenIndex, depth
									{
										position414 := position
										depth++
										{
											position415, tokenIndex415, depth415 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l416
											}
											position++
											goto l415
										l416:
											position, tokenIndex, depth = position415, tokenIndex415, depth415
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l413
											}
											position++
										}
									l415:
										depth--
										add(rulePN_CHARS_BASE, position414)
									}
									goto l412
								l413:
									position, tokenIndex, depth = position412, tokenIndex412, depth412
									if buffer[position] != rune('_') {
										goto l410
									}
									position++
								}
							l412:
								depth--
								add(rulePN_CHARS_U, position411)
							}
							goto l409
						l410:
							position, tokenIndex, depth = position409, tokenIndex409, depth409
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l400
							}
							position++
						}
					l409:
						goto l399
					l400:
						position, tokenIndex, depth = position400, tokenIndex400, depth400
					}
					depth--
					add(ruleVARNAME, position398)
				}
				if !rules[rulews]() {
					goto l394
				}
				depth--
				add(rulevar, position395)
			}
			return true
		l394:
			position, tokenIndex, depth = position394, tokenIndex394, depth394
			return false
		},
		/* 41 iriref <- <(iri / prefixedName)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l420
					}
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					{
						position421 := position
						depth++
					l422:
						{
							position423, tokenIndex423, depth423 := position, tokenIndex, depth
							{
								position424, tokenIndex424, depth424 := position, tokenIndex, depth
								{
									position425, tokenIndex425, depth425 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l426
									}
									position++
									goto l425
								l426:
									position, tokenIndex, depth = position425, tokenIndex425, depth425
									if buffer[position] != rune(' ') {
										goto l424
									}
									position++
								}
							l425:
								goto l423
							l424:
								position, tokenIndex, depth = position424, tokenIndex424, depth424
							}
							if !matchDot() {
								goto l423
							}
							goto l422
						l423:
							position, tokenIndex, depth = position423, tokenIndex423, depth423
						}
						if buffer[position] != rune(':') {
							goto l417
						}
						position++
					l427:
						{
							position428, tokenIndex428, depth428 := position, tokenIndex, depth
							{
								position429, tokenIndex429, depth429 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l430
								}
								position++
								goto l429
							l430:
								position, tokenIndex, depth = position429, tokenIndex429, depth429
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l431
								}
								position++
								goto l429
							l431:
								position, tokenIndex, depth = position429, tokenIndex429, depth429
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l432
								}
								position++
								goto l429
							l432:
								position, tokenIndex, depth = position429, tokenIndex429, depth429
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l428
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l428
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l428
										}
										position++
										break
									}
								}

							}
						l429:
							goto l427
						l428:
							position, tokenIndex, depth = position428, tokenIndex428, depth428
						}
						if !rules[rulews]() {
							goto l417
						}
						depth--
						add(ruleprefixedName, position421)
					}
				}
			l419:
				depth--
				add(ruleiriref, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 42 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
			{
				position435 := position
				depth++
				if buffer[position] != rune('<') {
					goto l434
				}
				position++
			l436:
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l438
						}
						position++
						goto l437
					l438:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
					}
					if !matchDot() {
						goto l437
					}
					goto l436
				l437:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
				}
				if buffer[position] != rune('>') {
					goto l434
				}
				position++
				if !rules[rulews]() {
					goto l434
				}
				depth--
				add(ruleiri, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
			return false
		},
		/* 43 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* ws)> */
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
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if buffer[position] != rune('{') {
					goto l461
				}
				position++
				if !rules[rulews]() {
					goto l461
				}
				depth--
				add(ruleLBRACE, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 66 RBRACE <- <('}' ws)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				if buffer[position] != rune('}') {
					goto l463
				}
				position++
				if !rules[rulews]() {
					goto l463
				}
				depth--
				add(ruleRBRACE, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
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
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				if buffer[position] != rune('.') {
					goto l469
				}
				position++
				if !rules[rulews]() {
					goto l469
				}
				depth--
				add(ruleDOT, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 72 COLON <- <(':' ws)> */
		nil,
		/* 73 PIPE <- <('|' ws)> */
		func() bool {
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				if buffer[position] != rune('|') {
					goto l472
				}
				position++
				if !rules[rulews]() {
					goto l472
				}
				depth--
				add(rulePIPE, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 74 SLASH <- <('/' ws)> */
		nil,
		/* 75 INVERSE <- <('^' ws)> */
		func() bool {
			position475, tokenIndex475, depth475 := position, tokenIndex, depth
			{
				position476 := position
				depth++
				if buffer[position] != rune('^') {
					goto l475
				}
				position++
				if !rules[rulews]() {
					goto l475
				}
				depth--
				add(ruleINVERSE, position476)
			}
			return true
		l475:
			position, tokenIndex, depth = position475, tokenIndex475, depth475
			return false
		},
		/* 76 LPAREN <- <('(' ws)> */
		func() bool {
			position477, tokenIndex477, depth477 := position, tokenIndex, depth
			{
				position478 := position
				depth++
				if buffer[position] != rune('(') {
					goto l477
				}
				position++
				if !rules[rulews]() {
					goto l477
				}
				depth--
				add(ruleLPAREN, position478)
			}
			return true
		l477:
			position, tokenIndex, depth = position477, tokenIndex477, depth477
			return false
		},
		/* 77 RPAREN <- <(')' ws)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				if buffer[position] != rune(')') {
					goto l479
				}
				position++
				if !rules[rulews]() {
					goto l479
				}
				depth--
				add(ruleRPAREN, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 78 ISA <- <('a' ws)> */
		func() bool {
			position481, tokenIndex481, depth481 := position, tokenIndex, depth
			{
				position482 := position
				depth++
				if buffer[position] != rune('a') {
					goto l481
				}
				position++
				if !rules[rulews]() {
					goto l481
				}
				depth--
				add(ruleISA, position482)
			}
			return true
		l481:
			position, tokenIndex, depth = position481, tokenIndex481, depth481
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
			position489, tokenIndex489, depth489 := position, tokenIndex, depth
			{
				position490 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l489
				}
				position++
			l491:
				{
					position492, tokenIndex492, depth492 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l492
					}
					position++
					goto l491
				l492:
					position, tokenIndex, depth = position492, tokenIndex492, depth492
				}
				if !rules[rulews]() {
					goto l489
				}
				depth--
				add(ruleINTEGER, position490)
			}
			return true
		l489:
			position, tokenIndex, depth = position489, tokenIndex489, depth489
			return false
		},
		/* 86 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position494 := position
				depth++
			l495:
				{
					position496, tokenIndex496, depth496 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\v':
							if buffer[position] != rune('\v') {
								goto l496
							}
							position++
							break
						case '\f':
							if buffer[position] != rune('\f') {
								goto l496
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l496
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l496
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l496
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l496
							}
							position++
							break
						}
					}

					goto l495
				l496:
					position, tokenIndex, depth = position496, tokenIndex496, depth496
				}
				depth--
				add(rulews, position494)
			}
			return true
		},
	}
	p.rules = rules
}
