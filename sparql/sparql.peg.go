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
	ruleconstructQuery
	ruleconstruct
	ruledescribeQuery
	ruledescribe
	ruleaskQuery
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
	ruleCONSTRUCT
	ruleDESCRIBE
	ruleASK
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
	"constructQuery",
	"construct",
	"describeQuery",
	"describe",
	"askQuery",
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
	"CONSTRUCT",
	"DESCRIBE",
	"ASK",
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
	rules  [96]func() bool
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
						switch buffer[position] {
						case 'A', 'a':
							{
								position39 := position
								depth++
								{
									position40 := position
									depth++
									{
										position41, tokenIndex41, depth41 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l42
										}
										position++
										goto l41
									l42:
										position, tokenIndex, depth = position41, tokenIndex41, depth41
										if buffer[position] != rune('A') {
											goto l0
										}
										position++
									}
								l41:
									{
										position43, tokenIndex43, depth43 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l44
										}
										position++
										goto l43
									l44:
										position, tokenIndex, depth = position43, tokenIndex43, depth43
										if buffer[position] != rune('S') {
											goto l0
										}
										position++
									}
								l43:
									{
										position45, tokenIndex45, depth45 := position, tokenIndex, depth
										if buffer[position] != rune('k') {
											goto l46
										}
										position++
										goto l45
									l46:
										position, tokenIndex, depth = position45, tokenIndex45, depth45
										if buffer[position] != rune('K') {
											goto l0
										}
										position++
									}
								l45:
									if !rules[rulews]() {
										goto l0
									}
									depth--
									add(ruleASK, position40)
								}
								{
									position47, tokenIndex47, depth47 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l47
									}
									goto l48
								l47:
									position, tokenIndex, depth = position47, tokenIndex47, depth47
								}
							l48:
								if !rules[rulewhereClause]() {
									goto l0
								}
								depth--
								add(ruleaskQuery, position39)
							}
							break
						case 'D', 'd':
							{
								position49 := position
								depth++
								{
									position50 := position
									depth++
									{
										position51 := position
										depth++
										{
											position52, tokenIndex52, depth52 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l53
											}
											position++
											goto l52
										l53:
											position, tokenIndex, depth = position52, tokenIndex52, depth52
											if buffer[position] != rune('D') {
												goto l0
											}
											position++
										}
									l52:
										{
											position54, tokenIndex54, depth54 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l55
											}
											position++
											goto l54
										l55:
											position, tokenIndex, depth = position54, tokenIndex54, depth54
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l54:
										{
											position56, tokenIndex56, depth56 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l57
											}
											position++
											goto l56
										l57:
											position, tokenIndex, depth = position56, tokenIndex56, depth56
											if buffer[position] != rune('S') {
												goto l0
											}
											position++
										}
									l56:
										{
											position58, tokenIndex58, depth58 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l59
											}
											position++
											goto l58
										l59:
											position, tokenIndex, depth = position58, tokenIndex58, depth58
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l58:
										{
											position60, tokenIndex60, depth60 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l61
											}
											position++
											goto l60
										l61:
											position, tokenIndex, depth = position60, tokenIndex60, depth60
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l60:
										{
											position62, tokenIndex62, depth62 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l63
											}
											position++
											goto l62
										l63:
											position, tokenIndex, depth = position62, tokenIndex62, depth62
											if buffer[position] != rune('I') {
												goto l0
											}
											position++
										}
									l62:
										{
											position64, tokenIndex64, depth64 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l65
											}
											position++
											goto l64
										l65:
											position, tokenIndex, depth = position64, tokenIndex64, depth64
											if buffer[position] != rune('B') {
												goto l0
											}
											position++
										}
									l64:
										{
											position66, tokenIndex66, depth66 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l67
											}
											position++
											goto l66
										l67:
											position, tokenIndex, depth = position66, tokenIndex66, depth66
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l66:
										if !rules[rulews]() {
											goto l0
										}
										depth--
										add(ruleDESCRIBE, position51)
									}
									{
										position68, tokenIndex68, depth68 := position, tokenIndex, depth
										if !rules[ruleSTAR]() {
											goto l69
										}
										goto l68
									l69:
										position, tokenIndex, depth = position68, tokenIndex68, depth68
										if !rules[rulevar]() {
											goto l70
										}
										goto l68
									l70:
										position, tokenIndex, depth = position68, tokenIndex68, depth68
										if !rules[ruleiriref]() {
											goto l0
										}
									}
								l68:
									depth--
									add(ruledescribe, position50)
								}
								{
									position71, tokenIndex71, depth71 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l71
									}
									goto l72
								l71:
									position, tokenIndex, depth = position71, tokenIndex71, depth71
								}
							l72:
								{
									position73, tokenIndex73, depth73 := position, tokenIndex, depth
									if !rules[rulewhereClause]() {
										goto l73
									}
									goto l74
								l73:
									position, tokenIndex, depth = position73, tokenIndex73, depth73
								}
							l74:
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruledescribeQuery, position49)
							}
							break
						case 'C', 'c':
							{
								position75 := position
								depth++
								{
									position76 := position
									depth++
									{
										position77 := position
										depth++
										{
											position78, tokenIndex78, depth78 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l79
											}
											position++
											goto l78
										l79:
											position, tokenIndex, depth = position78, tokenIndex78, depth78
											if buffer[position] != rune('C') {
												goto l0
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
												goto l0
											}
											position++
										}
									l80:
										{
											position82, tokenIndex82, depth82 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l83
											}
											position++
											goto l82
										l83:
											position, tokenIndex, depth = position82, tokenIndex82, depth82
											if buffer[position] != rune('N') {
												goto l0
											}
											position++
										}
									l82:
										{
											position84, tokenIndex84, depth84 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l85
											}
											position++
											goto l84
										l85:
											position, tokenIndex, depth = position84, tokenIndex84, depth84
											if buffer[position] != rune('S') {
												goto l0
											}
											position++
										}
									l84:
										{
											position86, tokenIndex86, depth86 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l87
											}
											position++
											goto l86
										l87:
											position, tokenIndex, depth = position86, tokenIndex86, depth86
											if buffer[position] != rune('T') {
												goto l0
											}
											position++
										}
									l86:
										{
											position88, tokenIndex88, depth88 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l89
											}
											position++
											goto l88
										l89:
											position, tokenIndex, depth = position88, tokenIndex88, depth88
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l88:
										{
											position90, tokenIndex90, depth90 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l91
											}
											position++
											goto l90
										l91:
											position, tokenIndex, depth = position90, tokenIndex90, depth90
											if buffer[position] != rune('U') {
												goto l0
											}
											position++
										}
									l90:
										{
											position92, tokenIndex92, depth92 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l93
											}
											position++
											goto l92
										l93:
											position, tokenIndex, depth = position92, tokenIndex92, depth92
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l92:
										{
											position94, tokenIndex94, depth94 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l95
											}
											position++
											goto l94
										l95:
											position, tokenIndex, depth = position94, tokenIndex94, depth94
											if buffer[position] != rune('T') {
												goto l0
											}
											position++
										}
									l94:
										if !rules[rulews]() {
											goto l0
										}
										depth--
										add(ruleCONSTRUCT, position77)
									}
									if !rules[ruleLBRACE]() {
										goto l0
									}
									{
										position96, tokenIndex96, depth96 := position, tokenIndex, depth
										if !rules[ruletriplesBlock]() {
											goto l96
										}
										goto l97
									l96:
										position, tokenIndex, depth = position96, tokenIndex96, depth96
									}
								l97:
									if !rules[ruleRBRACE]() {
										goto l0
									}
									depth--
									add(ruleconstruct, position76)
								}
								{
									position98, tokenIndex98, depth98 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l98
									}
									goto l99
								l98:
									position, tokenIndex, depth = position98, tokenIndex98, depth98
								}
							l99:
								if !rules[rulewhereClause]() {
									goto l0
								}
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruleconstructQuery, position75)
							}
							break
						default:
							{
								position100 := position
								depth++
								if !rules[ruleselect]() {
									goto l0
								}
								{
									position101, tokenIndex101, depth101 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l101
									}
									goto l102
								l101:
									position, tokenIndex, depth = position101, tokenIndex101, depth101
								}
							l102:
								if !rules[rulewhereClause]() {
									goto l0
								}
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruleselectQuery, position100)
							}
							break
						}
					}

					depth--
					add(rulequery, position37)
				}
				{
					position103, tokenIndex103, depth103 := position, tokenIndex, depth
					if !matchDot() {
						goto l103
					}
					goto l0
				l103:
					position, tokenIndex, depth = position103, tokenIndex103, depth103
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
		/* 4 query <- <((&('A' | 'a') askQuery) | (&('D' | 'd') describeQuery) | (&('C' | 'c') constructQuery) | (&('S' | 's') selectQuery))> */
		nil,
		/* 5 selectQuery <- <(select datasetClause? whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position109, tokenIndex109, depth109 := position, tokenIndex, depth
			{
				position110 := position
				depth++
				{
					position111 := position
					depth++
					{
						position112, tokenIndex112, depth112 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l113
						}
						position++
						goto l112
					l113:
						position, tokenIndex, depth = position112, tokenIndex112, depth112
						if buffer[position] != rune('S') {
							goto l109
						}
						position++
					}
				l112:
					{
						position114, tokenIndex114, depth114 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l115
						}
						position++
						goto l114
					l115:
						position, tokenIndex, depth = position114, tokenIndex114, depth114
						if buffer[position] != rune('E') {
							goto l109
						}
						position++
					}
				l114:
					{
						position116, tokenIndex116, depth116 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l117
						}
						position++
						goto l116
					l117:
						position, tokenIndex, depth = position116, tokenIndex116, depth116
						if buffer[position] != rune('L') {
							goto l109
						}
						position++
					}
				l116:
					{
						position118, tokenIndex118, depth118 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l119
						}
						position++
						goto l118
					l119:
						position, tokenIndex, depth = position118, tokenIndex118, depth118
						if buffer[position] != rune('E') {
							goto l109
						}
						position++
					}
				l118:
					{
						position120, tokenIndex120, depth120 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l121
						}
						position++
						goto l120
					l121:
						position, tokenIndex, depth = position120, tokenIndex120, depth120
						if buffer[position] != rune('C') {
							goto l109
						}
						position++
					}
				l120:
					{
						position122, tokenIndex122, depth122 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l123
						}
						position++
						goto l122
					l123:
						position, tokenIndex, depth = position122, tokenIndex122, depth122
						if buffer[position] != rune('T') {
							goto l109
						}
						position++
					}
				l122:
					if !rules[rulews]() {
						goto l109
					}
					depth--
					add(ruleSELECT, position111)
				}
				{
					position124, tokenIndex124, depth124 := position, tokenIndex, depth
					{
						position126, tokenIndex126, depth126 := position, tokenIndex, depth
						{
							position128 := position
							depth++
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
									goto l127
								}
								position++
							}
						l129:
							{
								position131, tokenIndex131, depth131 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l132
								}
								position++
								goto l131
							l132:
								position, tokenIndex, depth = position131, tokenIndex131, depth131
								if buffer[position] != rune('I') {
									goto l127
								}
								position++
							}
						l131:
							{
								position133, tokenIndex133, depth133 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l134
								}
								position++
								goto l133
							l134:
								position, tokenIndex, depth = position133, tokenIndex133, depth133
								if buffer[position] != rune('S') {
									goto l127
								}
								position++
							}
						l133:
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l136
								}
								position++
								goto l135
							l136:
								position, tokenIndex, depth = position135, tokenIndex135, depth135
								if buffer[position] != rune('T') {
									goto l127
								}
								position++
							}
						l135:
							{
								position137, tokenIndex137, depth137 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l138
								}
								position++
								goto l137
							l138:
								position, tokenIndex, depth = position137, tokenIndex137, depth137
								if buffer[position] != rune('I') {
									goto l127
								}
								position++
							}
						l137:
							{
								position139, tokenIndex139, depth139 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l140
								}
								position++
								goto l139
							l140:
								position, tokenIndex, depth = position139, tokenIndex139, depth139
								if buffer[position] != rune('N') {
									goto l127
								}
								position++
							}
						l139:
							{
								position141, tokenIndex141, depth141 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l142
								}
								position++
								goto l141
							l142:
								position, tokenIndex, depth = position141, tokenIndex141, depth141
								if buffer[position] != rune('C') {
									goto l127
								}
								position++
							}
						l141:
							{
								position143, tokenIndex143, depth143 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l144
								}
								position++
								goto l143
							l144:
								position, tokenIndex, depth = position143, tokenIndex143, depth143
								if buffer[position] != rune('T') {
									goto l127
								}
								position++
							}
						l143:
							if !rules[rulews]() {
								goto l127
							}
							depth--
							add(ruleDISTINCT, position128)
						}
						goto l126
					l127:
						position, tokenIndex, depth = position126, tokenIndex126, depth126
						{
							position145 := position
							depth++
							{
								position146, tokenIndex146, depth146 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l147
								}
								position++
								goto l146
							l147:
								position, tokenIndex, depth = position146, tokenIndex146, depth146
								if buffer[position] != rune('R') {
									goto l124
								}
								position++
							}
						l146:
							{
								position148, tokenIndex148, depth148 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l149
								}
								position++
								goto l148
							l149:
								position, tokenIndex, depth = position148, tokenIndex148, depth148
								if buffer[position] != rune('E') {
									goto l124
								}
								position++
							}
						l148:
							{
								position150, tokenIndex150, depth150 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l151
								}
								position++
								goto l150
							l151:
								position, tokenIndex, depth = position150, tokenIndex150, depth150
								if buffer[position] != rune('D') {
									goto l124
								}
								position++
							}
						l150:
							{
								position152, tokenIndex152, depth152 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l153
								}
								position++
								goto l152
							l153:
								position, tokenIndex, depth = position152, tokenIndex152, depth152
								if buffer[position] != rune('U') {
									goto l124
								}
								position++
							}
						l152:
							{
								position154, tokenIndex154, depth154 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l155
								}
								position++
								goto l154
							l155:
								position, tokenIndex, depth = position154, tokenIndex154, depth154
								if buffer[position] != rune('C') {
									goto l124
								}
								position++
							}
						l154:
							{
								position156, tokenIndex156, depth156 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l157
								}
								position++
								goto l156
							l157:
								position, tokenIndex, depth = position156, tokenIndex156, depth156
								if buffer[position] != rune('E') {
									goto l124
								}
								position++
							}
						l156:
							{
								position158, tokenIndex158, depth158 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l159
								}
								position++
								goto l158
							l159:
								position, tokenIndex, depth = position158, tokenIndex158, depth158
								if buffer[position] != rune('D') {
									goto l124
								}
								position++
							}
						l158:
							if !rules[rulews]() {
								goto l124
							}
							depth--
							add(ruleREDUCED, position145)
						}
					}
				l126:
					goto l125
				l124:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
				}
			l125:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l161
					}
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					{
						position164 := position
						depth++
						if !rules[rulevar]() {
							goto l109
						}
						depth--
						add(ruleprojectionElem, position164)
					}
				l162:
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						{
							position165 := position
							depth++
							if !rules[rulevar]() {
								goto l163
							}
							depth--
							add(ruleprojectionElem, position165)
						}
						goto l162
					l163:
						position, tokenIndex, depth = position163, tokenIndex163, depth163
					}
				}
			l160:
				depth--
				add(ruleselect, position110)
			}
			return true
		l109:
			position, tokenIndex, depth = position109, tokenIndex109, depth109
			return false
		},
		/* 7 subSelect <- <(select whereClause)> */
		func() bool {
			position166, tokenIndex166, depth166 := position, tokenIndex, depth
			{
				position167 := position
				depth++
				if !rules[ruleselect]() {
					goto l166
				}
				if !rules[rulewhereClause]() {
					goto l166
				}
				depth--
				add(rulesubSelect, position167)
			}
			return true
		l166:
			position, tokenIndex, depth = position166, tokenIndex166, depth166
			return false
		},
		/* 8 constructQuery <- <(construct datasetClause? whereClause solutionModifier)> */
		nil,
		/* 9 construct <- <(CONSTRUCT LBRACE triplesBlock? RBRACE)> */
		nil,
		/* 10 describeQuery <- <(describe datasetClause? whereClause? solutionModifier)> */
		nil,
		/* 11 describe <- <(DESCRIBE (STAR / var / iriref))> */
		nil,
		/* 12 askQuery <- <(ASK datasetClause? whereClause)> */
		nil,
		/* 13 projectionElem <- <var> */
		nil,
		/* 14 datasetClause <- <(FROM NAMED? iriref)> */
		func() bool {
			position174, tokenIndex174, depth174 := position, tokenIndex, depth
			{
				position175 := position
				depth++
				{
					position176 := position
					depth++
					{
						position177, tokenIndex177, depth177 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l178
						}
						position++
						goto l177
					l178:
						position, tokenIndex, depth = position177, tokenIndex177, depth177
						if buffer[position] != rune('F') {
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
						if buffer[position] != rune('o') {
							goto l182
						}
						position++
						goto l181
					l182:
						position, tokenIndex, depth = position181, tokenIndex181, depth181
						if buffer[position] != rune('O') {
							goto l174
						}
						position++
					}
				l181:
					{
						position183, tokenIndex183, depth183 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l184
						}
						position++
						goto l183
					l184:
						position, tokenIndex, depth = position183, tokenIndex183, depth183
						if buffer[position] != rune('M') {
							goto l174
						}
						position++
					}
				l183:
					if !rules[rulews]() {
						goto l174
					}
					depth--
					add(ruleFROM, position176)
				}
				{
					position185, tokenIndex185, depth185 := position, tokenIndex, depth
					{
						position187 := position
						depth++
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
								goto l185
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
								goto l185
							}
							position++
						}
					l190:
						{
							position192, tokenIndex192, depth192 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l193
							}
							position++
							goto l192
						l193:
							position, tokenIndex, depth = position192, tokenIndex192, depth192
							if buffer[position] != rune('M') {
								goto l185
							}
							position++
						}
					l192:
						{
							position194, tokenIndex194, depth194 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l195
							}
							position++
							goto l194
						l195:
							position, tokenIndex, depth = position194, tokenIndex194, depth194
							if buffer[position] != rune('E') {
								goto l185
							}
							position++
						}
					l194:
						{
							position196, tokenIndex196, depth196 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l197
							}
							position++
							goto l196
						l197:
							position, tokenIndex, depth = position196, tokenIndex196, depth196
							if buffer[position] != rune('D') {
								goto l185
							}
							position++
						}
					l196:
						if !rules[rulews]() {
							goto l185
						}
						depth--
						add(ruleNAMED, position187)
					}
					goto l186
				l185:
					position, tokenIndex, depth = position185, tokenIndex185, depth185
				}
			l186:
				if !rules[ruleiriref]() {
					goto l174
				}
				depth--
				add(ruledatasetClause, position175)
			}
			return true
		l174:
			position, tokenIndex, depth = position174, tokenIndex174, depth174
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position198, tokenIndex198, depth198 := position, tokenIndex, depth
			{
				position199 := position
				depth++
				{
					position200, tokenIndex200, depth200 := position, tokenIndex, depth
					{
						position202 := position
						depth++
						{
							position203, tokenIndex203, depth203 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l204
							}
							position++
							goto l203
						l204:
							position, tokenIndex, depth = position203, tokenIndex203, depth203
							if buffer[position] != rune('W') {
								goto l200
							}
							position++
						}
					l203:
						{
							position205, tokenIndex205, depth205 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l206
							}
							position++
							goto l205
						l206:
							position, tokenIndex, depth = position205, tokenIndex205, depth205
							if buffer[position] != rune('H') {
								goto l200
							}
							position++
						}
					l205:
						{
							position207, tokenIndex207, depth207 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l208
							}
							position++
							goto l207
						l208:
							position, tokenIndex, depth = position207, tokenIndex207, depth207
							if buffer[position] != rune('E') {
								goto l200
							}
							position++
						}
					l207:
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('R') {
								goto l200
							}
							position++
						}
					l209:
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l212
							}
							position++
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if buffer[position] != rune('E') {
								goto l200
							}
							position++
						}
					l211:
						if !rules[rulews]() {
							goto l200
						}
						depth--
						add(ruleWHERE, position202)
					}
					goto l201
				l200:
					position, tokenIndex, depth = position200, tokenIndex200, depth200
				}
			l201:
				if !rules[rulegroupGraphPattern]() {
					goto l198
				}
				depth--
				add(rulewhereClause, position199)
			}
			return true
		l198:
			position, tokenIndex, depth = position198, tokenIndex198, depth198
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position213, tokenIndex213, depth213 := position, tokenIndex, depth
			{
				position214 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l213
				}
				{
					position215, tokenIndex215, depth215 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l216
					}
					goto l215
				l216:
					position, tokenIndex, depth = position215, tokenIndex215, depth215
					if !rules[rulegraphPattern]() {
						goto l213
					}
				}
			l215:
				if !rules[ruleRBRACE]() {
					goto l213
				}
				depth--
				add(rulegroupGraphPattern, position214)
			}
			return true
		l213:
			position, tokenIndex, depth = position213, tokenIndex213, depth213
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position218 := position
				depth++
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					{
						position221 := position
						depth++
						if !rules[ruletriplesBlock]() {
							goto l219
						}
						depth--
						add(rulebasicGraphPattern, position221)
					}
					goto l220
				l219:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
				}
			l220:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					{
						position224 := position
						depth++
						{
							position225, tokenIndex225, depth225 := position, tokenIndex, depth
							{
								position227 := position
								depth++
								{
									position228 := position
									depth++
									{
										position229, tokenIndex229, depth229 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l230
										}
										position++
										goto l229
									l230:
										position, tokenIndex, depth = position229, tokenIndex229, depth229
										if buffer[position] != rune('O') {
											goto l226
										}
										position++
									}
								l229:
									{
										position231, tokenIndex231, depth231 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l232
										}
										position++
										goto l231
									l232:
										position, tokenIndex, depth = position231, tokenIndex231, depth231
										if buffer[position] != rune('P') {
											goto l226
										}
										position++
									}
								l231:
									{
										position233, tokenIndex233, depth233 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l234
										}
										position++
										goto l233
									l234:
										position, tokenIndex, depth = position233, tokenIndex233, depth233
										if buffer[position] != rune('T') {
											goto l226
										}
										position++
									}
								l233:
									{
										position235, tokenIndex235, depth235 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l236
										}
										position++
										goto l235
									l236:
										position, tokenIndex, depth = position235, tokenIndex235, depth235
										if buffer[position] != rune('I') {
											goto l226
										}
										position++
									}
								l235:
									{
										position237, tokenIndex237, depth237 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l238
										}
										position++
										goto l237
									l238:
										position, tokenIndex, depth = position237, tokenIndex237, depth237
										if buffer[position] != rune('O') {
											goto l226
										}
										position++
									}
								l237:
									{
										position239, tokenIndex239, depth239 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l240
										}
										position++
										goto l239
									l240:
										position, tokenIndex, depth = position239, tokenIndex239, depth239
										if buffer[position] != rune('N') {
											goto l226
										}
										position++
									}
								l239:
									{
										position241, tokenIndex241, depth241 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l242
										}
										position++
										goto l241
									l242:
										position, tokenIndex, depth = position241, tokenIndex241, depth241
										if buffer[position] != rune('A') {
											goto l226
										}
										position++
									}
								l241:
									{
										position243, tokenIndex243, depth243 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l244
										}
										position++
										goto l243
									l244:
										position, tokenIndex, depth = position243, tokenIndex243, depth243
										if buffer[position] != rune('L') {
											goto l226
										}
										position++
									}
								l243:
									if !rules[rulews]() {
										goto l226
									}
									depth--
									add(ruleOPTIONAL, position228)
								}
								if !rules[ruleLBRACE]() {
									goto l226
								}
								{
									position245, tokenIndex245, depth245 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l246
									}
									goto l245
								l246:
									position, tokenIndex, depth = position245, tokenIndex245, depth245
									if !rules[rulegraphPattern]() {
										goto l226
									}
								}
							l245:
								if !rules[ruleRBRACE]() {
									goto l226
								}
								depth--
								add(ruleoptionalGraphPattern, position227)
							}
							goto l225
						l226:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l222
							}
						}
					l225:
						depth--
						add(rulegraphPatternNotTriples, position224)
					}
					{
						position247, tokenIndex247, depth247 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l247
						}
						goto l248
					l247:
						position, tokenIndex, depth = position247, tokenIndex247, depth247
					}
				l248:
					if !rules[rulegraphPattern]() {
						goto l222
					}
					goto l223
				l222:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
				}
			l223:
				depth--
				add(rulegraphPattern, position218)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position251, tokenIndex251, depth251 := position, tokenIndex, depth
			{
				position252 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l251
				}
				{
					position253, tokenIndex253, depth253 := position, tokenIndex, depth
					{
						position255 := position
						depth++
						{
							position256, tokenIndex256, depth256 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l257
							}
							position++
							goto l256
						l257:
							position, tokenIndex, depth = position256, tokenIndex256, depth256
							if buffer[position] != rune('U') {
								goto l253
							}
							position++
						}
					l256:
						{
							position258, tokenIndex258, depth258 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l259
							}
							position++
							goto l258
						l259:
							position, tokenIndex, depth = position258, tokenIndex258, depth258
							if buffer[position] != rune('N') {
								goto l253
							}
							position++
						}
					l258:
						{
							position260, tokenIndex260, depth260 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l261
							}
							position++
							goto l260
						l261:
							position, tokenIndex, depth = position260, tokenIndex260, depth260
							if buffer[position] != rune('I') {
								goto l253
							}
							position++
						}
					l260:
						{
							position262, tokenIndex262, depth262 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l263
							}
							position++
							goto l262
						l263:
							position, tokenIndex, depth = position262, tokenIndex262, depth262
							if buffer[position] != rune('O') {
								goto l253
							}
							position++
						}
					l262:
						{
							position264, tokenIndex264, depth264 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l265
							}
							position++
							goto l264
						l265:
							position, tokenIndex, depth = position264, tokenIndex264, depth264
							if buffer[position] != rune('N') {
								goto l253
							}
							position++
						}
					l264:
						if !rules[rulews]() {
							goto l253
						}
						depth--
						add(ruleUNION, position255)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l253
					}
					goto l254
				l253:
					position, tokenIndex, depth = position253, tokenIndex253, depth253
				}
			l254:
				depth--
				add(rulegroupOrUnionGraphPattern, position252)
			}
			return true
		l251:
			position, tokenIndex, depth = position251, tokenIndex251, depth251
			return false
		},
		/* 21 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 22 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position267, tokenIndex267, depth267 := position, tokenIndex, depth
			{
				position268 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l267
				}
			l269:
				{
					position270, tokenIndex270, depth270 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l270
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l270
					}
					goto l269
				l270:
					position, tokenIndex, depth = position270, tokenIndex270, depth270
				}
				{
					position271, tokenIndex271, depth271 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l271
					}
					goto l272
				l271:
					position, tokenIndex, depth = position271, tokenIndex271, depth271
				}
			l272:
				depth--
				add(ruletriplesBlock, position268)
			}
			return true
		l267:
			position, tokenIndex, depth = position267, tokenIndex267, depth267
			return false
		},
		/* 23 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position273, tokenIndex273, depth273 := position, tokenIndex, depth
			{
				position274 := position
				depth++
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l276
					}
					if !rules[rulepropertyListPath]() {
						goto l276
					}
					goto l275
				l276:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
					{
						position277 := position
						depth++
						{
							position278, tokenIndex278, depth278 := position, tokenIndex, depth
							{
								position280 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l279
								}
								if !rules[rulegraphNodePath]() {
									goto l279
								}
							l281:
								{
									position282, tokenIndex282, depth282 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l282
									}
									goto l281
								l282:
									position, tokenIndex, depth = position282, tokenIndex282, depth282
								}
								if !rules[ruleRPAREN]() {
									goto l279
								}
								depth--
								add(rulecollectionPath, position280)
							}
							goto l278
						l279:
							position, tokenIndex, depth = position278, tokenIndex278, depth278
							{
								position283 := position
								depth++
								{
									position284 := position
									depth++
									if buffer[position] != rune('[') {
										goto l273
									}
									position++
									if !rules[rulews]() {
										goto l273
									}
									depth--
									add(ruleLBRACK, position284)
								}
								if !rules[rulepropertyListPath]() {
									goto l273
								}
								{
									position285 := position
									depth++
									if buffer[position] != rune(']') {
										goto l273
									}
									position++
									if !rules[rulews]() {
										goto l273
									}
									depth--
									add(ruleRBRACK, position285)
								}
								depth--
								add(ruleblankNodePropertyListPath, position283)
							}
						}
					l278:
						depth--
						add(ruletriplesNodePath, position277)
					}
					if !rules[rulepropertyListPath]() {
						goto l273
					}
				}
			l275:
				depth--
				add(ruletriplesSameSubjectPath, position274)
			}
			return true
		l273:
			position, tokenIndex, depth = position273, tokenIndex273, depth273
			return false
		},
		/* 24 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position286, tokenIndex286, depth286 := position, tokenIndex, depth
			{
				position287 := position
				depth++
				{
					position288, tokenIndex288, depth288 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l289
					}
					goto l288
				l289:
					position, tokenIndex, depth = position288, tokenIndex288, depth288
					{
						position290 := position
						depth++
						{
							position291, tokenIndex291, depth291 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l292
							}
							goto l291
						l292:
							position, tokenIndex, depth = position291, tokenIndex291, depth291
							{
								switch buffer[position] {
								case '(':
									{
										position294 := position
										depth++
										if buffer[position] != rune('(') {
											goto l286
										}
										position++
										if !rules[rulews]() {
											goto l286
										}
										if buffer[position] != rune(')') {
											goto l286
										}
										position++
										if !rules[rulews]() {
											goto l286
										}
										depth--
										add(rulenil, position294)
									}
									break
								case '[', '_':
									{
										position295 := position
										depth++
										{
											position296, tokenIndex296, depth296 := position, tokenIndex, depth
											{
												position298 := position
												depth++
												if buffer[position] != rune('_') {
													goto l297
												}
												position++
												if buffer[position] != rune(':') {
													goto l297
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l297
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l297
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l297
														}
														position++
														break
													}
												}

												{
													position300, tokenIndex300, depth300 := position, tokenIndex, depth
													{
														position302, tokenIndex302, depth302 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l303
														}
														position++
														goto l302
													l303:
														position, tokenIndex, depth = position302, tokenIndex302, depth302
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l304
														}
														position++
														goto l302
													l304:
														position, tokenIndex, depth = position302, tokenIndex302, depth302
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l305
														}
														position++
														goto l302
													l305:
														position, tokenIndex, depth = position302, tokenIndex302, depth302
														if c := buffer[position]; c < rune('.') || c > rune('_') {
															goto l300
														}
														position++
													}
												l302:
													goto l301
												l300:
													position, tokenIndex, depth = position300, tokenIndex300, depth300
												}
											l301:
												if !rules[rulews]() {
													goto l297
												}
												depth--
												add(ruleblankNodeLabel, position298)
											}
											goto l296
										l297:
											position, tokenIndex, depth = position296, tokenIndex296, depth296
											{
												position306 := position
												depth++
												if buffer[position] != rune('[') {
													goto l286
												}
												position++
												if !rules[rulews]() {
													goto l286
												}
												if buffer[position] != rune(']') {
													goto l286
												}
												position++
												if !rules[rulews]() {
													goto l286
												}
												depth--
												add(ruleanon, position306)
											}
										}
									l296:
										depth--
										add(ruleblankNode, position295)
									}
									break
								case 'F', 'T', 'f', 't':
									{
										position307 := position
										depth++
										{
											position308, tokenIndex308, depth308 := position, tokenIndex, depth
											{
												position310 := position
												depth++
												{
													position311, tokenIndex311, depth311 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l312
													}
													position++
													goto l311
												l312:
													position, tokenIndex, depth = position311, tokenIndex311, depth311
													if buffer[position] != rune('T') {
														goto l309
													}
													position++
												}
											l311:
												{
													position313, tokenIndex313, depth313 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l314
													}
													position++
													goto l313
												l314:
													position, tokenIndex, depth = position313, tokenIndex313, depth313
													if buffer[position] != rune('R') {
														goto l309
													}
													position++
												}
											l313:
												{
													position315, tokenIndex315, depth315 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l316
													}
													position++
													goto l315
												l316:
													position, tokenIndex, depth = position315, tokenIndex315, depth315
													if buffer[position] != rune('U') {
														goto l309
													}
													position++
												}
											l315:
												{
													position317, tokenIndex317, depth317 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l318
													}
													position++
													goto l317
												l318:
													position, tokenIndex, depth = position317, tokenIndex317, depth317
													if buffer[position] != rune('E') {
														goto l309
													}
													position++
												}
											l317:
												if !rules[rulews]() {
													goto l309
												}
												depth--
												add(ruleTRUE, position310)
											}
											goto l308
										l309:
											position, tokenIndex, depth = position308, tokenIndex308, depth308
											{
												position319 := position
												depth++
												{
													position320, tokenIndex320, depth320 := position, tokenIndex, depth
													if buffer[position] != rune('f') {
														goto l321
													}
													position++
													goto l320
												l321:
													position, tokenIndex, depth = position320, tokenIndex320, depth320
													if buffer[position] != rune('F') {
														goto l286
													}
													position++
												}
											l320:
												{
													position322, tokenIndex322, depth322 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l323
													}
													position++
													goto l322
												l323:
													position, tokenIndex, depth = position322, tokenIndex322, depth322
													if buffer[position] != rune('A') {
														goto l286
													}
													position++
												}
											l322:
												{
													position324, tokenIndex324, depth324 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l325
													}
													position++
													goto l324
												l325:
													position, tokenIndex, depth = position324, tokenIndex324, depth324
													if buffer[position] != rune('L') {
														goto l286
													}
													position++
												}
											l324:
												{
													position326, tokenIndex326, depth326 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l327
													}
													position++
													goto l326
												l327:
													position, tokenIndex, depth = position326, tokenIndex326, depth326
													if buffer[position] != rune('S') {
														goto l286
													}
													position++
												}
											l326:
												{
													position328, tokenIndex328, depth328 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l329
													}
													position++
													goto l328
												l329:
													position, tokenIndex, depth = position328, tokenIndex328, depth328
													if buffer[position] != rune('E') {
														goto l286
													}
													position++
												}
											l328:
												if !rules[rulews]() {
													goto l286
												}
												depth--
												add(ruleFALSE, position319)
											}
										}
									l308:
										depth--
										add(rulebooleanLiteral, position307)
									}
									break
								case '"':
									{
										position330 := position
										depth++
										{
											position331 := position
											depth++
											if buffer[position] != rune('"') {
												goto l286
											}
											position++
										l332:
											{
												position333, tokenIndex333, depth333 := position, tokenIndex, depth
												{
													position334, tokenIndex334, depth334 := position, tokenIndex, depth
													if buffer[position] != rune('"') {
														goto l334
													}
													position++
													goto l333
												l334:
													position, tokenIndex, depth = position334, tokenIndex334, depth334
												}
												if !matchDot() {
													goto l333
												}
												goto l332
											l333:
												position, tokenIndex, depth = position333, tokenIndex333, depth333
											}
											if buffer[position] != rune('"') {
												goto l286
											}
											position++
											depth--
											add(rulestring, position331)
										}
										{
											position335, tokenIndex335, depth335 := position, tokenIndex, depth
											{
												position337, tokenIndex337, depth337 := position, tokenIndex, depth
												if buffer[position] != rune('@') {
													goto l338
												}
												position++
												{
													position341, tokenIndex341, depth341 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l342
													}
													position++
													goto l341
												l342:
													position, tokenIndex, depth = position341, tokenIndex341, depth341
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l338
													}
													position++
												}
											l341:
											l339:
												{
													position340, tokenIndex340, depth340 := position, tokenIndex, depth
													{
														position343, tokenIndex343, depth343 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l344
														}
														position++
														goto l343
													l344:
														position, tokenIndex, depth = position343, tokenIndex343, depth343
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l340
														}
														position++
													}
												l343:
													goto l339
												l340:
													position, tokenIndex, depth = position340, tokenIndex340, depth340
												}
											l345:
												{
													position346, tokenIndex346, depth346 := position, tokenIndex, depth
													if buffer[position] != rune('-') {
														goto l346
													}
													position++
													{
														switch buffer[position] {
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l346
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l346
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l346
															}
															position++
															break
														}
													}

												l347:
													{
														position348, tokenIndex348, depth348 := position, tokenIndex, depth
														{
															switch buffer[position] {
															case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l348
																}
																position++
																break
															case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
																if c := buffer[position]; c < rune('A') || c > rune('Z') {
																	goto l348
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('a') || c > rune('z') {
																	goto l348
																}
																position++
																break
															}
														}

														goto l347
													l348:
														position, tokenIndex, depth = position348, tokenIndex348, depth348
													}
													goto l345
												l346:
													position, tokenIndex, depth = position346, tokenIndex346, depth346
												}
												goto l337
											l338:
												position, tokenIndex, depth = position337, tokenIndex337, depth337
												if buffer[position] != rune('^') {
													goto l335
												}
												position++
												if buffer[position] != rune('^') {
													goto l335
												}
												position++
												if !rules[ruleiriref]() {
													goto l335
												}
											}
										l337:
											goto l336
										l335:
											position, tokenIndex, depth = position335, tokenIndex335, depth335
										}
									l336:
										if !rules[rulews]() {
											goto l286
										}
										depth--
										add(ruleliteral, position330)
									}
									break
								default:
									{
										position351 := position
										depth++
										{
											position352, tokenIndex352, depth352 := position, tokenIndex, depth
											{
												position354, tokenIndex354, depth354 := position, tokenIndex, depth
												if buffer[position] != rune('+') {
													goto l355
												}
												position++
												goto l354
											l355:
												position, tokenIndex, depth = position354, tokenIndex354, depth354
												if buffer[position] != rune('-') {
													goto l352
												}
												position++
											}
										l354:
											goto l353
										l352:
											position, tokenIndex, depth = position352, tokenIndex352, depth352
										}
									l353:
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l286
										}
										position++
									l356:
										{
											position357, tokenIndex357, depth357 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l357
											}
											position++
											goto l356
										l357:
											position, tokenIndex, depth = position357, tokenIndex357, depth357
										}
										{
											position358, tokenIndex358, depth358 := position, tokenIndex, depth
											if buffer[position] != rune('.') {
												goto l358
											}
											position++
										l360:
											{
												position361, tokenIndex361, depth361 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l361
												}
												position++
												goto l360
											l361:
												position, tokenIndex, depth = position361, tokenIndex361, depth361
											}
											goto l359
										l358:
											position, tokenIndex, depth = position358, tokenIndex358, depth358
										}
									l359:
										if !rules[rulews]() {
											goto l286
										}
										depth--
										add(rulenumericLiteral, position351)
									}
									break
								}
							}

						}
					l291:
						depth--
						add(rulegraphTerm, position290)
					}
				}
			l288:
				depth--
				add(rulevarOrTerm, position287)
			}
			return true
		l286:
			position, tokenIndex, depth = position286, tokenIndex286, depth286
			return false
		},
		/* 25 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 26 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 27 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 28 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 29 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			position366, tokenIndex366, depth366 := position, tokenIndex, depth
			{
				position367 := position
				depth++
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l369
					}
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					{
						position370 := position
						depth++
						if !rules[rulepath]() {
							goto l366
						}
						depth--
						add(ruleverbPath, position370)
					}
				}
			l368:
				if !rules[ruleobjectListPath]() {
					goto l366
				}
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					{
						position373 := position
						depth++
						if buffer[position] != rune(';') {
							goto l371
						}
						position++
						if !rules[rulews]() {
							goto l371
						}
						depth--
						add(ruleSEMICOLON, position373)
					}
					if !rules[rulepropertyListPath]() {
						goto l371
					}
					goto l372
				l371:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
				}
			l372:
				depth--
				add(rulepropertyListPath, position367)
			}
			return true
		l366:
			position, tokenIndex, depth = position366, tokenIndex366, depth366
			return false
		},
		/* 30 verbPath <- <path> */
		nil,
		/* 31 path <- <pathAlternative> */
		func() bool {
			position375, tokenIndex375, depth375 := position, tokenIndex, depth
			{
				position376 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l375
				}
				depth--
				add(rulepath, position376)
			}
			return true
		l375:
			position, tokenIndex, depth = position375, tokenIndex375, depth375
			return false
		},
		/* 32 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position377, tokenIndex377, depth377 := position, tokenIndex, depth
			{
				position378 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l377
				}
			l379:
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l380
					}
					if !rules[rulepathAlternative]() {
						goto l380
					}
					goto l379
				l380:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
				}
				depth--
				add(rulepathAlternative, position378)
			}
			return true
		l377:
			position, tokenIndex, depth = position377, tokenIndex377, depth377
			return false
		},
		/* 33 pathSequence <- <(pathElt (SLASH pathSequence)*)> */
		func() bool {
			position381, tokenIndex381, depth381 := position, tokenIndex, depth
			{
				position382 := position
				depth++
				{
					position383 := position
					depth++
					{
						position384, tokenIndex384, depth384 := position, tokenIndex, depth
						if !rules[ruleINVERSE]() {
							goto l384
						}
						goto l385
					l384:
						position, tokenIndex, depth = position384, tokenIndex384, depth384
					}
				l385:
					{
						position386 := position
						depth++
						{
							position387, tokenIndex387, depth387 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l388
							}
							goto l387
						l388:
							position, tokenIndex, depth = position387, tokenIndex387, depth387
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l381
									}
									if !rules[rulepath]() {
										goto l381
									}
									if !rules[ruleRPAREN]() {
										goto l381
									}
									break
								case '!':
									{
										position390 := position
										depth++
										if buffer[position] != rune('!') {
											goto l381
										}
										position++
										if !rules[rulews]() {
											goto l381
										}
										depth--
										add(ruleNOT, position390)
									}
									{
										position391 := position
										depth++
										{
											position392, tokenIndex392, depth392 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l393
											}
											goto l392
										l393:
											position, tokenIndex, depth = position392, tokenIndex392, depth392
											if !rules[ruleLPAREN]() {
												goto l381
											}
											{
												position394, tokenIndex394, depth394 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l394
												}
											l396:
												{
													position397, tokenIndex397, depth397 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l397
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l397
													}
													goto l396
												l397:
													position, tokenIndex, depth = position397, tokenIndex397, depth397
												}
												goto l395
											l394:
												position, tokenIndex, depth = position394, tokenIndex394, depth394
											}
										l395:
											if !rules[ruleRPAREN]() {
												goto l381
											}
										}
									l392:
										depth--
										add(rulepathNegatedPropertySet, position391)
									}
									break
								default:
									if !rules[ruleISA]() {
										goto l381
									}
									break
								}
							}

						}
					l387:
						depth--
						add(rulepathPrimary, position386)
					}
					depth--
					add(rulepathElt, position383)
				}
			l398:
				{
					position399, tokenIndex399, depth399 := position, tokenIndex, depth
					{
						position400 := position
						depth++
						if buffer[position] != rune('/') {
							goto l399
						}
						position++
						if !rules[rulews]() {
							goto l399
						}
						depth--
						add(ruleSLASH, position400)
					}
					if !rules[rulepathSequence]() {
						goto l399
					}
					goto l398
				l399:
					position, tokenIndex, depth = position399, tokenIndex399, depth399
				}
				depth--
				add(rulepathSequence, position382)
			}
			return true
		l381:
			position, tokenIndex, depth = position381, tokenIndex381, depth381
			return false
		},
		/* 34 pathElt <- <(INVERSE? pathPrimary)> */
		nil,
		/* 35 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 36 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 37 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l407
					}
					goto l406
				l407:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if !rules[ruleISA]() {
						goto l408
					}
					goto l406
				l408:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
					if !rules[ruleINVERSE]() {
						goto l404
					}
					{
						position409, tokenIndex409, depth409 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l410
						}
						goto l409
					l410:
						position, tokenIndex, depth = position409, tokenIndex409, depth409
						if !rules[ruleISA]() {
							goto l404
						}
					}
				l409:
				}
			l406:
				depth--
				add(rulepathOneInPropertySet, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 38 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position411, tokenIndex411, depth411 := position, tokenIndex, depth
			{
				position412 := position
				depth++
				{
					position413 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l411
					}
					depth--
					add(ruleobjectPath, position413)
				}
			l414:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					{
						position416 := position
						depth++
						if buffer[position] != rune(',') {
							goto l415
						}
						position++
						if !rules[rulews]() {
							goto l415
						}
						depth--
						add(ruleCOMMA, position416)
					}
					if !rules[ruleobjectListPath]() {
						goto l415
					}
					goto l414
				l415:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
				}
				depth--
				add(ruleobjectListPath, position412)
			}
			return true
		l411:
			position, tokenIndex, depth = position411, tokenIndex411, depth411
			return false
		},
		/* 39 objectPath <- <graphNodePath> */
		nil,
		/* 40 graphNodePath <- <varOrTerm> */
		func() bool {
			position418, tokenIndex418, depth418 := position, tokenIndex, depth
			{
				position419 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l418
				}
				depth--
				add(rulegraphNodePath, position419)
			}
			return true
		l418:
			position, tokenIndex, depth = position418, tokenIndex418, depth418
			return false
		},
		/* 41 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position421 := position
				depth++
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					{
						position424 := position
						depth++
						{
							position425, tokenIndex425, depth425 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l426
							}
							{
								position427, tokenIndex427, depth427 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l427
								}
								goto l428
							l427:
								position, tokenIndex, depth = position427, tokenIndex427, depth427
							}
						l428:
							goto l425
						l426:
							position, tokenIndex, depth = position425, tokenIndex425, depth425
							if !rules[ruleoffset]() {
								goto l422
							}
							{
								position429, tokenIndex429, depth429 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l429
								}
								goto l430
							l429:
								position, tokenIndex, depth = position429, tokenIndex429, depth429
							}
						l430:
						}
					l425:
						depth--
						add(rulelimitOffsetClauses, position424)
					}
					goto l423
				l422:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
				}
			l423:
				depth--
				add(rulesolutionModifier, position421)
			}
			return true
		},
		/* 42 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 43 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position432, tokenIndex432, depth432 := position, tokenIndex, depth
			{
				position433 := position
				depth++
				{
					position434 := position
					depth++
					{
						position435, tokenIndex435, depth435 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l436
						}
						position++
						goto l435
					l436:
						position, tokenIndex, depth = position435, tokenIndex435, depth435
						if buffer[position] != rune('L') {
							goto l432
						}
						position++
					}
				l435:
					{
						position437, tokenIndex437, depth437 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l438
						}
						position++
						goto l437
					l438:
						position, tokenIndex, depth = position437, tokenIndex437, depth437
						if buffer[position] != rune('I') {
							goto l432
						}
						position++
					}
				l437:
					{
						position439, tokenIndex439, depth439 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l440
						}
						position++
						goto l439
					l440:
						position, tokenIndex, depth = position439, tokenIndex439, depth439
						if buffer[position] != rune('M') {
							goto l432
						}
						position++
					}
				l439:
					{
						position441, tokenIndex441, depth441 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l442
						}
						position++
						goto l441
					l442:
						position, tokenIndex, depth = position441, tokenIndex441, depth441
						if buffer[position] != rune('I') {
							goto l432
						}
						position++
					}
				l441:
					{
						position443, tokenIndex443, depth443 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l444
						}
						position++
						goto l443
					l444:
						position, tokenIndex, depth = position443, tokenIndex443, depth443
						if buffer[position] != rune('T') {
							goto l432
						}
						position++
					}
				l443:
					if !rules[rulews]() {
						goto l432
					}
					depth--
					add(ruleLIMIT, position434)
				}
				if !rules[ruleINTEGER]() {
					goto l432
				}
				depth--
				add(rulelimit, position433)
			}
			return true
		l432:
			position, tokenIndex, depth = position432, tokenIndex432, depth432
			return false
		},
		/* 44 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				{
					position447 := position
					depth++
					{
						position448, tokenIndex448, depth448 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l449
						}
						position++
						goto l448
					l449:
						position, tokenIndex, depth = position448, tokenIndex448, depth448
						if buffer[position] != rune('O') {
							goto l445
						}
						position++
					}
				l448:
					{
						position450, tokenIndex450, depth450 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l451
						}
						position++
						goto l450
					l451:
						position, tokenIndex, depth = position450, tokenIndex450, depth450
						if buffer[position] != rune('F') {
							goto l445
						}
						position++
					}
				l450:
					{
						position452, tokenIndex452, depth452 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l453
						}
						position++
						goto l452
					l453:
						position, tokenIndex, depth = position452, tokenIndex452, depth452
						if buffer[position] != rune('F') {
							goto l445
						}
						position++
					}
				l452:
					{
						position454, tokenIndex454, depth454 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l455
						}
						position++
						goto l454
					l455:
						position, tokenIndex, depth = position454, tokenIndex454, depth454
						if buffer[position] != rune('S') {
							goto l445
						}
						position++
					}
				l454:
					{
						position456, tokenIndex456, depth456 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l457
						}
						position++
						goto l456
					l457:
						position, tokenIndex, depth = position456, tokenIndex456, depth456
						if buffer[position] != rune('E') {
							goto l445
						}
						position++
					}
				l456:
					{
						position458, tokenIndex458, depth458 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l459
						}
						position++
						goto l458
					l459:
						position, tokenIndex, depth = position458, tokenIndex458, depth458
						if buffer[position] != rune('T') {
							goto l445
						}
						position++
					}
				l458:
					if !rules[rulews]() {
						goto l445
					}
					depth--
					add(ruleOFFSET, position447)
				}
				if !rules[ruleINTEGER]() {
					goto l445
				}
				depth--
				add(ruleoffset, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 45 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position460, tokenIndex460, depth460 := position, tokenIndex, depth
			{
				position461 := position
				depth++
				{
					position462, tokenIndex462, depth462 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l463
					}
					position++
					goto l462
				l463:
					position, tokenIndex, depth = position462, tokenIndex462, depth462
					if buffer[position] != rune('$') {
						goto l460
					}
					position++
				}
			l462:
				{
					position464 := position
					depth++
					{
						position467, tokenIndex467, depth467 := position, tokenIndex, depth
						{
							position469 := position
							depth++
							{
								position470, tokenIndex470, depth470 := position, tokenIndex, depth
								{
									position472 := position
									depth++
									{
										position473, tokenIndex473, depth473 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l474
										}
										position++
										goto l473
									l474:
										position, tokenIndex, depth = position473, tokenIndex473, depth473
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l471
										}
										position++
									}
								l473:
									depth--
									add(rulePN_CHARS_BASE, position472)
								}
								goto l470
							l471:
								position, tokenIndex, depth = position470, tokenIndex470, depth470
								if buffer[position] != rune('_') {
									goto l468
								}
								position++
							}
						l470:
							depth--
							add(rulePN_CHARS_U, position469)
						}
						goto l467
					l468:
						position, tokenIndex, depth = position467, tokenIndex467, depth467
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l460
						}
						position++
					}
				l467:
				l465:
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						{
							position475, tokenIndex475, depth475 := position, tokenIndex, depth
							{
								position477 := position
								depth++
								{
									position478, tokenIndex478, depth478 := position, tokenIndex, depth
									{
										position480 := position
										depth++
										{
											position481, tokenIndex481, depth481 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l482
											}
											position++
											goto l481
										l482:
											position, tokenIndex, depth = position481, tokenIndex481, depth481
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l479
											}
											position++
										}
									l481:
										depth--
										add(rulePN_CHARS_BASE, position480)
									}
									goto l478
								l479:
									position, tokenIndex, depth = position478, tokenIndex478, depth478
									if buffer[position] != rune('_') {
										goto l476
									}
									position++
								}
							l478:
								depth--
								add(rulePN_CHARS_U, position477)
							}
							goto l475
						l476:
							position, tokenIndex, depth = position475, tokenIndex475, depth475
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l466
							}
							position++
						}
					l475:
						goto l465
					l466:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
					}
					depth--
					add(ruleVARNAME, position464)
				}
				if !rules[rulews]() {
					goto l460
				}
				depth--
				add(rulevar, position461)
			}
			return true
		l460:
			position, tokenIndex, depth = position460, tokenIndex460, depth460
			return false
		},
		/* 46 iriref <- <(iri / prefixedName)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l486
					}
					goto l485
				l486:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					{
						position487 := position
						depth++
					l488:
						{
							position489, tokenIndex489, depth489 := position, tokenIndex, depth
							{
								position490, tokenIndex490, depth490 := position, tokenIndex, depth
								{
									position491, tokenIndex491, depth491 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l492
									}
									position++
									goto l491
								l492:
									position, tokenIndex, depth = position491, tokenIndex491, depth491
									if buffer[position] != rune(' ') {
										goto l490
									}
									position++
								}
							l491:
								goto l489
							l490:
								position, tokenIndex, depth = position490, tokenIndex490, depth490
							}
							if !matchDot() {
								goto l489
							}
							goto l488
						l489:
							position, tokenIndex, depth = position489, tokenIndex489, depth489
						}
						if buffer[position] != rune(':') {
							goto l483
						}
						position++
					l493:
						{
							position494, tokenIndex494, depth494 := position, tokenIndex, depth
							{
								position495, tokenIndex495, depth495 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l496
								}
								position++
								goto l495
							l496:
								position, tokenIndex, depth = position495, tokenIndex495, depth495
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l497
								}
								position++
								goto l495
							l497:
								position, tokenIndex, depth = position495, tokenIndex495, depth495
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l498
								}
								position++
								goto l495
							l498:
								position, tokenIndex, depth = position495, tokenIndex495, depth495
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l494
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l494
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l494
										}
										position++
										break
									}
								}

							}
						l495:
							goto l493
						l494:
							position, tokenIndex, depth = position494, tokenIndex494, depth494
						}
						if !rules[rulews]() {
							goto l483
						}
						depth--
						add(ruleprefixedName, position487)
					}
				}
			l485:
				depth--
				add(ruleiriref, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 47 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position500, tokenIndex500, depth500 := position, tokenIndex, depth
			{
				position501 := position
				depth++
				if buffer[position] != rune('<') {
					goto l500
				}
				position++
			l502:
				{
					position503, tokenIndex503, depth503 := position, tokenIndex, depth
					{
						position504, tokenIndex504, depth504 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l504
						}
						position++
						goto l503
					l504:
						position, tokenIndex, depth = position504, tokenIndex504, depth504
					}
					if !matchDot() {
						goto l503
					}
					goto l502
				l503:
					position, tokenIndex, depth = position503, tokenIndex503, depth503
				}
				if buffer[position] != rune('>') {
					goto l500
				}
				position++
				if !rules[rulews]() {
					goto l500
				}
				depth--
				add(ruleiri, position501)
			}
			return true
		l500:
			position, tokenIndex, depth = position500, tokenIndex500, depth500
			return false
		},
		/* 48 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* ws)> */
		nil,
		/* 49 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? ws)> */
		nil,
		/* 50 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 51 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? ws)> */
		nil,
		/* 52 booleanLiteral <- <(TRUE / FALSE)> */
		nil,
		/* 53 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 54 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? ws)> */
		nil,
		/* 55 anon <- <('[' ws ']' ws)> */
		nil,
		/* 56 nil <- <('(' ws ')' ws)> */
		nil,
		/* 57 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 58 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 59 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 60 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') ws)> */
		nil,
		/* 61 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') ws)> */
		nil,
		/* 62 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 63 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 64 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 65 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 66 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 67 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') ws)> */
		nil,
		/* 68 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 69 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') ws)> */
		nil,
		/* 70 LBRACE <- <('{' ws)> */
		func() bool {
			position527, tokenIndex527, depth527 := position, tokenIndex, depth
			{
				position528 := position
				depth++
				if buffer[position] != rune('{') {
					goto l527
				}
				position++
				if !rules[rulews]() {
					goto l527
				}
				depth--
				add(ruleLBRACE, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 71 RBRACE <- <('}' ws)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				if buffer[position] != rune('}') {
					goto l529
				}
				position++
				if !rules[rulews]() {
					goto l529
				}
				depth--
				add(ruleRBRACE, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 72 LBRACK <- <('[' ws)> */
		nil,
		/* 73 RBRACK <- <(']' ws)> */
		nil,
		/* 74 SEMICOLON <- <(';' ws)> */
		nil,
		/* 75 COMMA <- <(',' ws)> */
		nil,
		/* 76 DOT <- <('.' ws)> */
		func() bool {
			position535, tokenIndex535, depth535 := position, tokenIndex, depth
			{
				position536 := position
				depth++
				if buffer[position] != rune('.') {
					goto l535
				}
				position++
				if !rules[rulews]() {
					goto l535
				}
				depth--
				add(ruleDOT, position536)
			}
			return true
		l535:
			position, tokenIndex, depth = position535, tokenIndex535, depth535
			return false
		},
		/* 77 COLON <- <(':' ws)> */
		nil,
		/* 78 PIPE <- <('|' ws)> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				if buffer[position] != rune('|') {
					goto l538
				}
				position++
				if !rules[rulews]() {
					goto l538
				}
				depth--
				add(rulePIPE, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 79 SLASH <- <('/' ws)> */
		nil,
		/* 80 INVERSE <- <('^' ws)> */
		func() bool {
			position541, tokenIndex541, depth541 := position, tokenIndex, depth
			{
				position542 := position
				depth++
				if buffer[position] != rune('^') {
					goto l541
				}
				position++
				if !rules[rulews]() {
					goto l541
				}
				depth--
				add(ruleINVERSE, position542)
			}
			return true
		l541:
			position, tokenIndex, depth = position541, tokenIndex541, depth541
			return false
		},
		/* 81 LPAREN <- <('(' ws)> */
		func() bool {
			position543, tokenIndex543, depth543 := position, tokenIndex, depth
			{
				position544 := position
				depth++
				if buffer[position] != rune('(') {
					goto l543
				}
				position++
				if !rules[rulews]() {
					goto l543
				}
				depth--
				add(ruleLPAREN, position544)
			}
			return true
		l543:
			position, tokenIndex, depth = position543, tokenIndex543, depth543
			return false
		},
		/* 82 RPAREN <- <(')' ws)> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				if buffer[position] != rune(')') {
					goto l545
				}
				position++
				if !rules[rulews]() {
					goto l545
				}
				depth--
				add(ruleRPAREN, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 83 ISA <- <('a' ws)> */
		func() bool {
			position547, tokenIndex547, depth547 := position, tokenIndex, depth
			{
				position548 := position
				depth++
				if buffer[position] != rune('a') {
					goto l547
				}
				position++
				if !rules[rulews]() {
					goto l547
				}
				depth--
				add(ruleISA, position548)
			}
			return true
		l547:
			position, tokenIndex, depth = position547, tokenIndex547, depth547
			return false
		},
		/* 84 NOT <- <('!' ws)> */
		nil,
		/* 85 STAR <- <('*' ws)> */
		func() bool {
			position550, tokenIndex550, depth550 := position, tokenIndex, depth
			{
				position551 := position
				depth++
				if buffer[position] != rune('*') {
					goto l550
				}
				position++
				if !rules[rulews]() {
					goto l550
				}
				depth--
				add(ruleSTAR, position551)
			}
			return true
		l550:
			position, tokenIndex, depth = position550, tokenIndex550, depth550
			return false
		},
		/* 86 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') ws)> */
		nil,
		/* 87 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') ws)> */
		nil,
		/* 88 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') ws)> */
		nil,
		/* 89 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') ws)> */
		nil,
		/* 90 INTEGER <- <([0-9]+ ws)> */
		func() bool {
			position556, tokenIndex556, depth556 := position, tokenIndex, depth
			{
				position557 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l556
				}
				position++
			l558:
				{
					position559, tokenIndex559, depth559 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l559
					}
					position++
					goto l558
				l559:
					position, tokenIndex, depth = position559, tokenIndex559, depth559
				}
				if !rules[rulews]() {
					goto l556
				}
				depth--
				add(ruleINTEGER, position557)
			}
			return true
		l556:
			position, tokenIndex, depth = position556, tokenIndex556, depth556
			return false
		},
		/* 91 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 92 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') ws)> */
		nil,
		/* 93 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') ws)> */
		nil,
		/* 94 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position564 := position
				depth++
			l565:
				{
					position566, tokenIndex566, depth566 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\v':
							if buffer[position] != rune('\v') {
								goto l566
							}
							position++
							break
						case '\f':
							if buffer[position] != rune('\f') {
								goto l566
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l566
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l566
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l566
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l566
							}
							position++
							break
						}
					}

					goto l565
				l566:
					position, tokenIndex, depth = position566, tokenIndex566, depth566
				}
				depth--
				add(rulews, position564)
			}
			return true
		},
	}
	p.rules = rules
}
