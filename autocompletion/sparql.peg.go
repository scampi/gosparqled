package autocompletion

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
	rulepof
	rulevar
	ruleiri
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
	ruleQUESTION
	rulePLUS
	ruleOPTIONAL
	ruleUNION
	ruleLIMIT
	ruleOFFSET
	ruleINTEGER
	rulews
	rulePegText
	ruleAction0
	ruleAction1
	ruleAction2
	ruleAction3
	ruleAction4
	ruleAction5
	ruleAction6
	ruleAction7
	ruleAction8

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
	"pof",
	"var",
	"iri",
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
	"QUESTION",
	"PLUS",
	"OPTIONAL",
	"UNION",
	"LIMIT",
	"OFFSET",
	"INTEGER",
	"ws",
	"PegText",
	"Action0",
	"Action1",
	"Action2",
	"Action3",
	"Action4",
	"Action5",
	"Action6",
	"Action7",
	"Action8",

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
	*Bgp

	Buffer string
	buffer []rune
	rules  [99]func() bool
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

func (p *Sparql) Execute() {
	buffer, begin, end := p.Buffer, 0, 0
	for token := range p.tokenTree.Tokens() {
		switch token.pegRule {
		case rulePegText:
			begin, end = int(token.begin), int(token.end)
		case ruleAction0:
			p.setSubject(buffer[begin:end])
		case ruleAction1:
			p.setSubject(buffer[begin:end])
		case ruleAction2:
			p.setSubject("?POF")
		case ruleAction3:
			p.setPredicate("?POF")
		case ruleAction4:
			p.setPredicate(buffer[begin:end])
		case ruleAction5:
			p.setPredicate(buffer[begin:end])
		case ruleAction6:
			p.setObject(buffer[begin:end])
			p.addTriplePattern()
		case ruleAction7:
			p.setObject("?POF")
			p.addTriplePattern()
		case ruleAction8:
			p.setObject("?FillVar")
			p.addTriplePattern()

		}
	}
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
			fmt.Println("\nqueryContainer")
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !rules[rulews]() {
					fmt.Print("0 ")
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
											fmt.Print("10 ")
											goto l10
										}
										position++
										fmt.Print("9 ")
										goto l9
									l10:
										position, tokenIndex, depth = position9, tokenIndex9, depth9
										if buffer[position] != rune('P') {
											fmt.Print("6 ")
											goto l6
										}
										position++
									}
								l9:
									{
										position11, tokenIndex11, depth11 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											fmt.Print("12 ")
											goto l12
										}
										position++
										fmt.Print("11 ")
										goto l11
									l12:
										position, tokenIndex, depth = position11, tokenIndex11, depth11
										if buffer[position] != rune('R') {
											fmt.Print("6 ")
											goto l6
										}
										position++
									}
								l11:
									{
										position13, tokenIndex13, depth13 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											fmt.Print("14 ")
											goto l14
										}
										position++
										fmt.Print("13 ")
										goto l13
									l14:
										position, tokenIndex, depth = position13, tokenIndex13, depth13
										if buffer[position] != rune('E') {
											fmt.Print("6 ")
											goto l6
										}
										position++
									}
								l13:
									{
										position15, tokenIndex15, depth15 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											fmt.Print("16 ")
											goto l16
										}
										position++
										fmt.Print("15 ")
										goto l15
									l16:
										position, tokenIndex, depth = position15, tokenIndex15, depth15
										if buffer[position] != rune('F') {
											fmt.Print("6 ")
											goto l6
										}
										position++
									}
								l15:
									{
										position17, tokenIndex17, depth17 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											fmt.Print("18 ")
											goto l18
										}
										position++
										fmt.Print("17 ")
										goto l17
									l18:
										position, tokenIndex, depth = position17, tokenIndex17, depth17
										if buffer[position] != rune('I') {
											fmt.Print("6 ")
											goto l6
										}
										position++
									}
								l17:
									{
										position19, tokenIndex19, depth19 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											fmt.Print("20 ")
											goto l20
										}
										position++
										fmt.Print("19 ")
										goto l19
									l20:
										position, tokenIndex, depth = position19, tokenIndex19, depth19
										if buffer[position] != rune('X') {
											fmt.Print("6 ")
											goto l6
										}
										position++
									}
								l19:
									if !rules[rulews]() {
										fmt.Print("6 ")
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
											fmt.Print("25 ")
											goto l25
										}
										position++
										fmt.Print("24 ")
										goto l24
									l25:
										position, tokenIndex, depth = position24, tokenIndex24, depth24
										if buffer[position] != rune(' ') {
											fmt.Print("23 ")
											goto l23
										}
										position++
									}
								l24:
									fmt.Print("6 ")
									goto l6
								l23:
									position, tokenIndex, depth = position23, tokenIndex23, depth23
								}
								if !matchDot() {
									fmt.Print("6 ")
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
												fmt.Print("28 ")
												goto l28
											}
											position++
											fmt.Print("27 ")
											goto l27
										l28:
											position, tokenIndex, depth = position27, tokenIndex27, depth27
											if buffer[position] != rune(' ') {
												fmt.Print("26 ")
												goto l26
											}
											position++
										}
									l27:
										fmt.Print("22 ")
										goto l22
									l26:
										position, tokenIndex, depth = position26, tokenIndex26, depth26
									}
									if !matchDot() {
										fmt.Print("22 ")
										goto l22
									}
									fmt.Print("21 ")
									goto l21
								l22:
									position, tokenIndex, depth = position22, tokenIndex22, depth22
								}
								{
									position29 := position
									depth++
									if buffer[position] != rune(':') {
										fmt.Print("6 ")
										goto l6
									}
									position++
									if !rules[rulews]() {
										fmt.Print("6 ")
										goto l6
									}
									depth--
									add(ruleCOLON, position29)
								}
								if !rules[ruleiri]() {
									fmt.Print("6 ")
									goto l6
								}
								depth--
								add(ruleprefixDecl, position7)
							}
							fmt.Print("5 ")
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
											fmt.Print("33 ")
											goto l33
										}
										position++
										fmt.Print("32 ")
										goto l32
									l33:
										position, tokenIndex, depth = position32, tokenIndex32, depth32
										if buffer[position] != rune('B') {
											fmt.Print("4 ")
											goto l4
										}
										position++
									}
								l32:
									{
										position34, tokenIndex34, depth34 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											fmt.Print("35 ")
											goto l35
										}
										position++
										fmt.Print("34 ")
										goto l34
									l35:
										position, tokenIndex, depth = position34, tokenIndex34, depth34
										if buffer[position] != rune('A') {
											fmt.Print("4 ")
											goto l4
										}
										position++
									}
								l34:
									{
										position36, tokenIndex36, depth36 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											fmt.Print("37 ")
											goto l37
										}
										position++
										fmt.Print("36 ")
										goto l36
									l37:
										position, tokenIndex, depth = position36, tokenIndex36, depth36
										if buffer[position] != rune('S') {
											fmt.Print("4 ")
											goto l4
										}
										position++
									}
								l36:
									{
										position38, tokenIndex38, depth38 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											fmt.Print("39 ")
											goto l39
										}
										position++
										fmt.Print("38 ")
										goto l38
									l39:
										position, tokenIndex, depth = position38, tokenIndex38, depth38
										if buffer[position] != rune('E') {
											fmt.Print("4 ")
											goto l4
										}
										position++
									}
								l38:
									if !rules[rulews]() {
										fmt.Print("4 ")
										goto l4
									}
									depth--
									add(ruleBASE, position31)
								}
								if !rules[ruleiri]() {
									fmt.Print("4 ")
									goto l4
								}
								depth--
								add(rulebaseDecl, position30)
							}
						}
					l5:
						fmt.Print("3 ")
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
							fmt.Print("0 ")
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
											fmt.Print("47 ")
											goto l47
										}
										position++
										fmt.Print("46 ")
										goto l46
									l47:
										position, tokenIndex, depth = position46, tokenIndex46, depth46
										if buffer[position] != rune('F') {
											fmt.Print("42 ")
											goto l42
										}
										position++
									}
								l46:
									{
										position48, tokenIndex48, depth48 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											fmt.Print("49 ")
											goto l49
										}
										position++
										fmt.Print("48 ")
										goto l48
									l49:
										position, tokenIndex, depth = position48, tokenIndex48, depth48
										if buffer[position] != rune('R') {
											fmt.Print("42 ")
											goto l42
										}
										position++
									}
								l48:
									{
										position50, tokenIndex50, depth50 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											fmt.Print("51 ")
											goto l51
										}
										position++
										fmt.Print("50 ")
										goto l50
									l51:
										position, tokenIndex, depth = position50, tokenIndex50, depth50
										if buffer[position] != rune('O') {
											fmt.Print("42 ")
											goto l42
										}
										position++
									}
								l50:
									{
										position52, tokenIndex52, depth52 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											fmt.Print("53 ")
											goto l53
										}
										position++
										fmt.Print("52 ")
										goto l52
									l53:
										position, tokenIndex, depth = position52, tokenIndex52, depth52
										if buffer[position] != rune('M') {
											fmt.Print("42 ")
											goto l42
										}
										position++
									}
								l52:
									if !rules[rulews]() {
										fmt.Print("42 ")
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
												fmt.Print("58 ")
												goto l58
											}
											position++
											fmt.Print("57 ")
											goto l57
										l58:
											position, tokenIndex, depth = position57, tokenIndex57, depth57
											if buffer[position] != rune('N') {
												fmt.Print("54 ")
												goto l54
											}
											position++
										}
									l57:
										{
											position59, tokenIndex59, depth59 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												fmt.Print("60 ")
												goto l60
											}
											position++
											fmt.Print("59 ")
											goto l59
										l60:
											position, tokenIndex, depth = position59, tokenIndex59, depth59
											if buffer[position] != rune('A') {
												fmt.Print("54 ")
												goto l54
											}
											position++
										}
									l59:
										{
											position61, tokenIndex61, depth61 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												fmt.Print("62 ")
												goto l62
											}
											position++
											fmt.Print("61 ")
											goto l61
										l62:
											position, tokenIndex, depth = position61, tokenIndex61, depth61
											if buffer[position] != rune('M') {
												fmt.Print("54 ")
												goto l54
											}
											position++
										}
									l61:
										{
											position63, tokenIndex63, depth63 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												fmt.Print("64 ")
												goto l64
											}
											position++
											fmt.Print("63 ")
											goto l63
										l64:
											position, tokenIndex, depth = position63, tokenIndex63, depth63
											if buffer[position] != rune('E') {
												fmt.Print("54 ")
												goto l54
											}
											position++
										}
									l63:
										{
											position65, tokenIndex65, depth65 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												fmt.Print("66 ")
												goto l66
											}
											position++
											fmt.Print("65 ")
											goto l65
										l66:
											position, tokenIndex, depth = position65, tokenIndex65, depth65
											if buffer[position] != rune('D') {
												fmt.Print("54 ")
												goto l54
											}
											position++
										}
									l65:
										if !rules[rulews]() {
											fmt.Print("54 ")
											goto l54
										}
										depth--
										add(ruleNAMED, position56)
									}
									fmt.Print("55 ")
									goto l55
								l54:
									position, tokenIndex, depth = position54, tokenIndex54, depth54
								}
							l55:
								if !rules[ruleiri]() {
									fmt.Print("42 ")
									goto l42
								}
								depth--
								add(ruledatasetClause, position44)
							}
							fmt.Print("43 ")
							goto l43
						l42:
							position, tokenIndex, depth = position42, tokenIndex42, depth42
						}
					l43:
						if !rules[rulewhereClause]() {
							fmt.Print("0 ")
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
											fmt.Print("72 ")
											goto l72
										}
										{
											position73, tokenIndex73, depth73 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												fmt.Print("73 ")
												goto l73
											}
											fmt.Print("74 ")
											goto l74
										l73:
											position, tokenIndex, depth = position73, tokenIndex73, depth73
										}
									l74:
										fmt.Print("71 ")
										goto l71
									l72:
										position, tokenIndex, depth = position71, tokenIndex71, depth71
										if !rules[ruleoffset]() {
											fmt.Print("68 ")
											goto l68
										}
										{
											position75, tokenIndex75, depth75 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												fmt.Print("75 ")
												goto l75
											}
											fmt.Print("76 ")
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
								fmt.Print("69 ")
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
						fmt.Print("77 ")
						goto l77
					}
					fmt.Print("0 ")
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
			fmt.Println("\nselect")
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
							fmt.Print("87 ")
							goto l87
						}
						position++
						fmt.Print("86 ")
						goto l86
					l87:
						position, tokenIndex, depth = position86, tokenIndex86, depth86
						if buffer[position] != rune('S') {
							fmt.Print("83 ")
							goto l83
						}
						position++
					}
				l86:
					{
						position88, tokenIndex88, depth88 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							fmt.Print("89 ")
							goto l89
						}
						position++
						fmt.Print("88 ")
						goto l88
					l89:
						position, tokenIndex, depth = position88, tokenIndex88, depth88
						if buffer[position] != rune('E') {
							fmt.Print("83 ")
							goto l83
						}
						position++
					}
				l88:
					{
						position90, tokenIndex90, depth90 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							fmt.Print("91 ")
							goto l91
						}
						position++
						fmt.Print("90 ")
						goto l90
					l91:
						position, tokenIndex, depth = position90, tokenIndex90, depth90
						if buffer[position] != rune('L') {
							fmt.Print("83 ")
							goto l83
						}
						position++
					}
				l90:
					{
						position92, tokenIndex92, depth92 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							fmt.Print("93 ")
							goto l93
						}
						position++
						fmt.Print("92 ")
						goto l92
					l93:
						position, tokenIndex, depth = position92, tokenIndex92, depth92
						if buffer[position] != rune('E') {
							fmt.Print("83 ")
							goto l83
						}
						position++
					}
				l92:
					{
						position94, tokenIndex94, depth94 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							fmt.Print("95 ")
							goto l95
						}
						position++
						fmt.Print("94 ")
						goto l94
					l95:
						position, tokenIndex, depth = position94, tokenIndex94, depth94
						if buffer[position] != rune('C') {
							fmt.Print("83 ")
							goto l83
						}
						position++
					}
				l94:
					{
						position96, tokenIndex96, depth96 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							fmt.Print("97 ")
							goto l97
						}
						position++
						fmt.Print("96 ")
						goto l96
					l97:
						position, tokenIndex, depth = position96, tokenIndex96, depth96
						if buffer[position] != rune('T') {
							fmt.Print("83 ")
							goto l83
						}
						position++
					}
				l96:
					if !rules[rulews]() {
						fmt.Print("83 ")
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
									fmt.Print("104 ")
									goto l104
								}
								position++
								fmt.Print("103 ")
								goto l103
							l104:
								position, tokenIndex, depth = position103, tokenIndex103, depth103
								if buffer[position] != rune('D') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l103:
							{
								position105, tokenIndex105, depth105 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									fmt.Print("106 ")
									goto l106
								}
								position++
								fmt.Print("105 ")
								goto l105
							l106:
								position, tokenIndex, depth = position105, tokenIndex105, depth105
								if buffer[position] != rune('I') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l105:
							{
								position107, tokenIndex107, depth107 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									fmt.Print("108 ")
									goto l108
								}
								position++
								fmt.Print("107 ")
								goto l107
							l108:
								position, tokenIndex, depth = position107, tokenIndex107, depth107
								if buffer[position] != rune('S') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l107:
							{
								position109, tokenIndex109, depth109 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									fmt.Print("110 ")
									goto l110
								}
								position++
								fmt.Print("109 ")
								goto l109
							l110:
								position, tokenIndex, depth = position109, tokenIndex109, depth109
								if buffer[position] != rune('T') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l109:
							{
								position111, tokenIndex111, depth111 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									fmt.Print("112 ")
									goto l112
								}
								position++
								fmt.Print("111 ")
								goto l111
							l112:
								position, tokenIndex, depth = position111, tokenIndex111, depth111
								if buffer[position] != rune('I') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l111:
							{
								position113, tokenIndex113, depth113 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									fmt.Print("114 ")
									goto l114
								}
								position++
								fmt.Print("113 ")
								goto l113
							l114:
								position, tokenIndex, depth = position113, tokenIndex113, depth113
								if buffer[position] != rune('N') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l113:
							{
								position115, tokenIndex115, depth115 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									fmt.Print("116 ")
									goto l116
								}
								position++
								fmt.Print("115 ")
								goto l115
							l116:
								position, tokenIndex, depth = position115, tokenIndex115, depth115
								if buffer[position] != rune('C') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l115:
							{
								position117, tokenIndex117, depth117 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									fmt.Print("118 ")
									goto l118
								}
								position++
								fmt.Print("117 ")
								goto l117
							l118:
								position, tokenIndex, depth = position117, tokenIndex117, depth117
								if buffer[position] != rune('T') {
									fmt.Print("101 ")
									goto l101
								}
								position++
							}
						l117:
							if !rules[rulews]() {
								fmt.Print("101 ")
								goto l101
							}
							depth--
							add(ruleDISTINCT, position102)
						}
						fmt.Print("100 ")
						goto l100
					l101:
						position, tokenIndex, depth = position100, tokenIndex100, depth100
						{
							position119 := position
							depth++
							{
								position120, tokenIndex120, depth120 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									fmt.Print("121 ")
									goto l121
								}
								position++
								fmt.Print("120 ")
								goto l120
							l121:
								position, tokenIndex, depth = position120, tokenIndex120, depth120
								if buffer[position] != rune('R') {
									fmt.Print("98 ")
									goto l98
								}
								position++
							}
						l120:
							{
								position122, tokenIndex122, depth122 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									fmt.Print("123 ")
									goto l123
								}
								position++
								fmt.Print("122 ")
								goto l122
							l123:
								position, tokenIndex, depth = position122, tokenIndex122, depth122
								if buffer[position] != rune('E') {
									fmt.Print("98 ")
									goto l98
								}
								position++
							}
						l122:
							{
								position124, tokenIndex124, depth124 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									fmt.Print("125 ")
									goto l125
								}
								position++
								fmt.Print("124 ")
								goto l124
							l125:
								position, tokenIndex, depth = position124, tokenIndex124, depth124
								if buffer[position] != rune('D') {
									fmt.Print("98 ")
									goto l98
								}
								position++
							}
						l124:
							{
								position126, tokenIndex126, depth126 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									fmt.Print("127 ")
									goto l127
								}
								position++
								fmt.Print("126 ")
								goto l126
							l127:
								position, tokenIndex, depth = position126, tokenIndex126, depth126
								if buffer[position] != rune('U') {
									fmt.Print("98 ")
									goto l98
								}
								position++
							}
						l126:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									fmt.Print("129 ")
									goto l129
								}
								position++
								fmt.Print("128 ")
								goto l128
							l129:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
								if buffer[position] != rune('C') {
									fmt.Print("98 ")
									goto l98
								}
								position++
							}
						l128:
							{
								position130, tokenIndex130, depth130 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									fmt.Print("131 ")
									goto l131
								}
								position++
								fmt.Print("130 ")
								goto l130
							l131:
								position, tokenIndex, depth = position130, tokenIndex130, depth130
								if buffer[position] != rune('E') {
									fmt.Print("98 ")
									goto l98
								}
								position++
							}
						l130:
							{
								position132, tokenIndex132, depth132 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									fmt.Print("133 ")
									goto l133
								}
								position++
								fmt.Print("132 ")
								goto l132
							l133:
								position, tokenIndex, depth = position132, tokenIndex132, depth132
								if buffer[position] != rune('D') {
									fmt.Print("98 ")
									goto l98
								}
								position++
							}
						l132:
							if !rules[rulews]() {
								fmt.Print("98 ")
								goto l98
							}
							depth--
							add(ruleREDUCED, position119)
						}
					}
				l100:
					fmt.Print("99 ")
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
							fmt.Print("135 ")
							goto l135
						}
						position++
						if !rules[rulews]() {
							fmt.Print("135 ")
							goto l135
						}
						depth--
						add(ruleSTAR, position136)
					}
					fmt.Print("134 ")
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					{
						position139 := position
						depth++
						if !rules[rulevar]() {
							fmt.Print("83 ")
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
								fmt.Print("138 ")
								goto l138
							}
							depth--
							add(ruleprojectionElem, position140)
						}
						fmt.Print("137 ")
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
			fmt.Println("\nsubSelect")
			position141, tokenIndex141, depth141 := position, tokenIndex, depth
			{
				position142 := position
				depth++
				if !rules[ruleselect]() {
					fmt.Print("141 ")
					goto l141
				}
				if !rules[rulewhereClause]() {
					fmt.Print("141 ")
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
		/* 9 datasetClause <- <(FROM NAMED? iri)> */
		nil,
		/* 10 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			fmt.Println("\nwhereClause")
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
								fmt.Print("151 ")
								goto l151
							}
							position++
							fmt.Print("150 ")
							goto l150
						l151:
							position, tokenIndex, depth = position150, tokenIndex150, depth150
							if buffer[position] != rune('W') {
								fmt.Print("147 ")
								goto l147
							}
							position++
						}
					l150:
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								fmt.Print("153 ")
								goto l153
							}
							position++
							fmt.Print("152 ")
							goto l152
						l153:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('H') {
								fmt.Print("147 ")
								goto l147
							}
							position++
						}
					l152:
						{
							position154, tokenIndex154, depth154 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								fmt.Print("155 ")
								goto l155
							}
							position++
							fmt.Print("154 ")
							goto l154
						l155:
							position, tokenIndex, depth = position154, tokenIndex154, depth154
							if buffer[position] != rune('E') {
								fmt.Print("147 ")
								goto l147
							}
							position++
						}
					l154:
						{
							position156, tokenIndex156, depth156 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								fmt.Print("157 ")
								goto l157
							}
							position++
							fmt.Print("156 ")
							goto l156
						l157:
							position, tokenIndex, depth = position156, tokenIndex156, depth156
							if buffer[position] != rune('R') {
								fmt.Print("147 ")
								goto l147
							}
							position++
						}
					l156:
						{
							position158, tokenIndex158, depth158 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								fmt.Print("159 ")
								goto l159
							}
							position++
							fmt.Print("158 ")
							goto l158
						l159:
							position, tokenIndex, depth = position158, tokenIndex158, depth158
							if buffer[position] != rune('E') {
								fmt.Print("147 ")
								goto l147
							}
							position++
						}
					l158:
						if !rules[rulews]() {
							fmt.Print("147 ")
							goto l147
						}
						depth--
						add(ruleWHERE, position149)
					}
					fmt.Print("148 ")
					goto l148
				l147:
					position, tokenIndex, depth = position147, tokenIndex147, depth147
				}
			l148:
				if !rules[rulegroupGraphPattern]() {
					fmt.Print("145 ")
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
			fmt.Println("\ngroupGraphPattern")
			position160, tokenIndex160, depth160 := position, tokenIndex, depth
			{
				position161 := position
				depth++
				if !rules[ruleLBRACE]() {
					fmt.Print("160 ")
					goto l160
				}
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						fmt.Print("163 ")
						goto l163
					}
					fmt.Print("162 ")
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					if !rules[rulegraphPattern]() {
						fmt.Print("160 ")
						goto l160
					}
				}
			l162:
				if !rules[ruleRBRACE]() {
					fmt.Print("160 ")
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
			fmt.Println("\ngraphPattern")
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
								fmt.Print("166 ")
								goto l166
							}
						l170:
							{
								position171, tokenIndex171, depth171 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									fmt.Print("171 ")
									goto l171
								}
								if !rules[ruletriplesSameSubjectPath]() {
									fmt.Print("171 ")
									goto l171
								}
								fmt.Print("170 ")
								goto l170
							l171:
								position, tokenIndex, depth = position171, tokenIndex171, depth171
							}
							{
								position172, tokenIndex172, depth172 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									fmt.Print("172 ")
									goto l172
								}
								fmt.Print("173 ")
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
					fmt.Print("167 ")
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
											fmt.Print("182 ")
											goto l182
										}
										position++
										fmt.Print("181 ")
										goto l181
									l182:
										position, tokenIndex, depth = position181, tokenIndex181, depth181
										if buffer[position] != rune('O') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l181:
									{
										position183, tokenIndex183, depth183 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											fmt.Print("184 ")
											goto l184
										}
										position++
										fmt.Print("183 ")
										goto l183
									l184:
										position, tokenIndex, depth = position183, tokenIndex183, depth183
										if buffer[position] != rune('P') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l183:
									{
										position185, tokenIndex185, depth185 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											fmt.Print("186 ")
											goto l186
										}
										position++
										fmt.Print("185 ")
										goto l185
									l186:
										position, tokenIndex, depth = position185, tokenIndex185, depth185
										if buffer[position] != rune('T') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l185:
									{
										position187, tokenIndex187, depth187 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											fmt.Print("188 ")
											goto l188
										}
										position++
										fmt.Print("187 ")
										goto l187
									l188:
										position, tokenIndex, depth = position187, tokenIndex187, depth187
										if buffer[position] != rune('I') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l187:
									{
										position189, tokenIndex189, depth189 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											fmt.Print("190 ")
											goto l190
										}
										position++
										fmt.Print("189 ")
										goto l189
									l190:
										position, tokenIndex, depth = position189, tokenIndex189, depth189
										if buffer[position] != rune('O') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l189:
									{
										position191, tokenIndex191, depth191 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											fmt.Print("192 ")
											goto l192
										}
										position++
										fmt.Print("191 ")
										goto l191
									l192:
										position, tokenIndex, depth = position191, tokenIndex191, depth191
										if buffer[position] != rune('N') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l191:
									{
										position193, tokenIndex193, depth193 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											fmt.Print("194 ")
											goto l194
										}
										position++
										fmt.Print("193 ")
										goto l193
									l194:
										position, tokenIndex, depth = position193, tokenIndex193, depth193
										if buffer[position] != rune('A') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l193:
									{
										position195, tokenIndex195, depth195 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											fmt.Print("196 ")
											goto l196
										}
										position++
										fmt.Print("195 ")
										goto l195
									l196:
										position, tokenIndex, depth = position195, tokenIndex195, depth195
										if buffer[position] != rune('L') {
											fmt.Print("178 ")
											goto l178
										}
										position++
									}
								l195:
									if !rules[rulews]() {
										fmt.Print("178 ")
										goto l178
									}
									depth--
									add(ruleOPTIONAL, position180)
								}
								if !rules[ruleLBRACE]() {
									fmt.Print("178 ")
									goto l178
								}
								{
									position197, tokenIndex197, depth197 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										fmt.Print("198 ")
										goto l198
									}
									fmt.Print("197 ")
									goto l197
								l198:
									position, tokenIndex, depth = position197, tokenIndex197, depth197
									if !rules[rulegraphPattern]() {
										fmt.Print("178 ")
										goto l178
									}
								}
							l197:
								if !rules[ruleRBRACE]() {
									fmt.Print("178 ")
									goto l178
								}
								depth--
								add(ruleoptionalGraphPattern, position179)
							}
							fmt.Print("177 ")
							goto l177
						l178:
							position, tokenIndex, depth = position177, tokenIndex177, depth177
							if !rules[rulegroupOrUnionGraphPattern]() {
								fmt.Print("174 ")
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
							fmt.Print("199 ")
							goto l199
						}
						fmt.Print("200 ")
						goto l200
					l199:
						position, tokenIndex, depth = position199, tokenIndex199, depth199
					}
				l200:
					if !rules[rulegraphPattern]() {
						fmt.Print("174 ")
						goto l174
					}
					fmt.Print("175 ")
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
			fmt.Println("\ngroupOrUnionGraphPattern")
			position203, tokenIndex203, depth203 := position, tokenIndex, depth
			{
				position204 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					fmt.Print("203 ")
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
								fmt.Print("209 ")
								goto l209
							}
							position++
							fmt.Print("208 ")
							goto l208
						l209:
							position, tokenIndex, depth = position208, tokenIndex208, depth208
							if buffer[position] != rune('U') {
								fmt.Print("205 ")
								goto l205
							}
							position++
						}
					l208:
						{
							position210, tokenIndex210, depth210 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								fmt.Print("211 ")
								goto l211
							}
							position++
							fmt.Print("210 ")
							goto l210
						l211:
							position, tokenIndex, depth = position210, tokenIndex210, depth210
							if buffer[position] != rune('N') {
								fmt.Print("205 ")
								goto l205
							}
							position++
						}
					l210:
						{
							position212, tokenIndex212, depth212 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								fmt.Print("213 ")
								goto l213
							}
							position++
							fmt.Print("212 ")
							goto l212
						l213:
							position, tokenIndex, depth = position212, tokenIndex212, depth212
							if buffer[position] != rune('I') {
								fmt.Print("205 ")
								goto l205
							}
							position++
						}
					l212:
						{
							position214, tokenIndex214, depth214 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								fmt.Print("215 ")
								goto l215
							}
							position++
							fmt.Print("214 ")
							goto l214
						l215:
							position, tokenIndex, depth = position214, tokenIndex214, depth214
							if buffer[position] != rune('O') {
								fmt.Print("205 ")
								goto l205
							}
							position++
						}
					l214:
						{
							position216, tokenIndex216, depth216 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								fmt.Print("217 ")
								goto l217
							}
							position++
							fmt.Print("216 ")
							goto l216
						l217:
							position, tokenIndex, depth = position216, tokenIndex216, depth216
							if buffer[position] != rune('N') {
								fmt.Print("205 ")
								goto l205
							}
							position++
						}
					l216:
						if !rules[rulews]() {
							fmt.Print("205 ")
							goto l205
						}
						depth--
						add(ruleUNION, position207)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						fmt.Print("205 ")
						goto l205
					}
					fmt.Print("206 ")
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
			fmt.Println("\ntriplesSameSubjectPath")
			position220, tokenIndex220, depth220 := position, tokenIndex, depth
			{
				position221 := position
				depth++
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
								if !rules[rulevar]() {
									fmt.Print("226 ")
									goto l226
								}
								depth--
								add(rulePegText, position227)
							}
							{
								add(ruleAction0, position)
							}
							fmt.Print("225 ")
							goto l225
						l226:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							{
								position230 := position
								depth++
								if !rules[rulegraphTerm]() {
									fmt.Print("229 ")
									goto l229
								}
								depth--
								add(rulePegText, position230)
							}
							{
								add(ruleAction1, position)
							}
							fmt.Print("225 ")
							goto l225
						l229:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if !rules[rulepof]() {
								fmt.Print("223 ")
								goto l223
							}
							{
								add(ruleAction2, position)
							}
						}
					l225:
						depth--
						add(rulevarOrTerm, position224)
					}
					if !rules[rulepropertyListPath]() {
						fmt.Print("223 ")
						goto l223
					}
					fmt.Print("222 ")
					goto l222
				l223:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
					{
						position233 := position
						depth++
						{
							position234, tokenIndex234, depth234 := position, tokenIndex, depth
							{
								position236 := position
								depth++
								if !rules[ruleLPAREN]() {
									fmt.Print("235 ")
									goto l235
								}
								if !rules[rulegraphNodePath]() {
									fmt.Print("235 ")
									goto l235
								}
							l237:
								{
									position238, tokenIndex238, depth238 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										fmt.Print("238 ")
										goto l238
									}
									fmt.Print("237 ")
									goto l237
								l238:
									position, tokenIndex, depth = position238, tokenIndex238, depth238
								}
								if !rules[ruleRPAREN]() {
									fmt.Print("235 ")
									goto l235
								}
								depth--
								add(rulecollectionPath, position236)
							}
							fmt.Print("234 ")
							goto l234
						l235:
							position, tokenIndex, depth = position234, tokenIndex234, depth234
							{
								position239 := position
								depth++
								{
									position240 := position
									depth++
									if buffer[position] != rune('[') {
										fmt.Print("220 ")
										goto l220
									}
									position++
									if !rules[rulews]() {
										fmt.Print("220 ")
										goto l220
									}
									depth--
									add(ruleLBRACK, position240)
								}
								if !rules[rulepropertyListPath]() {
									fmt.Print("220 ")
									goto l220
								}
								{
									position241 := position
									depth++
									if buffer[position] != rune(']') {
										fmt.Print("220 ")
										goto l220
									}
									position++
									if !rules[rulews]() {
										fmt.Print("220 ")
										goto l220
									}
									depth--
									add(ruleRBRACK, position241)
								}
								depth--
								add(ruleblankNodePropertyListPath, position239)
							}
						}
					l234:
						depth--
						add(ruletriplesNodePath, position233)
					}
					if !rules[rulepropertyListPath]() {
						fmt.Print("220 ")
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
		/* 19 varOrTerm <- <((<var> Action0) / (<graphTerm> Action1) / (pof Action2))> */
		nil,
		/* 20 graphTerm <- <((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('<') iri) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral))> */
		func() bool {
			fmt.Println("\ngraphTerm")
			position243, tokenIndex243, depth243 := position, tokenIndex, depth
			{
				position244 := position
				depth++
				{
					switch buffer[position] {
					case '(':
						{
							position246 := position
							depth++
							if buffer[position] != rune('(') {
								fmt.Print("243 ")
								goto l243
							}
							position++
							if !rules[rulews]() {
								fmt.Print("243 ")
								goto l243
							}
							if buffer[position] != rune(')') {
								fmt.Print("243 ")
								goto l243
							}
							position++
							if !rules[rulews]() {
								fmt.Print("243 ")
								goto l243
							}
							depth--
							add(rulenil, position246)
						}
						break
					case '[', '_':
						{
							position247 := position
							depth++
							{
								position248, tokenIndex248, depth248 := position, tokenIndex, depth
								{
									position250 := position
									depth++
									if buffer[position] != rune('_') {
										fmt.Print("249 ")
										goto l249
									}
									position++
									if buffer[position] != rune(':') {
										fmt.Print("249 ")
										goto l249
									}
									position++
									{
										switch buffer[position] {
										case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												fmt.Print("249 ")
												goto l249
											}
											position++
											break
										case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												fmt.Print("249 ")
												goto l249
											}
											position++
											break
										default:
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												fmt.Print("249 ")
												goto l249
											}
											position++
											break
										}
									}

									{
										position252, tokenIndex252, depth252 := position, tokenIndex, depth
										{
											position254, tokenIndex254, depth254 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												fmt.Print("255 ")
												goto l255
											}
											position++
											fmt.Print("254 ")
											goto l254
										l255:
											position, tokenIndex, depth = position254, tokenIndex254, depth254
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												fmt.Print("256 ")
												goto l256
											}
											position++
											fmt.Print("254 ")
											goto l254
										l256:
											position, tokenIndex, depth = position254, tokenIndex254, depth254
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												fmt.Print("257 ")
												goto l257
											}
											position++
											fmt.Print("254 ")
											goto l254
										l257:
											position, tokenIndex, depth = position254, tokenIndex254, depth254
											if c := buffer[position]; c < rune('.') || c > rune('_') {
												fmt.Print("252 ")
												goto l252
											}
											position++
										}
									l254:
										fmt.Print("253 ")
										goto l253
									l252:
										position, tokenIndex, depth = position252, tokenIndex252, depth252
									}
								l253:
									if !rules[rulews]() {
										fmt.Print("249 ")
										goto l249
									}
									depth--
									add(ruleblankNodeLabel, position250)
								}
								fmt.Print("248 ")
								goto l248
							l249:
								position, tokenIndex, depth = position248, tokenIndex248, depth248
								{
									position258 := position
									depth++
									if buffer[position] != rune('[') {
										fmt.Print("243 ")
										goto l243
									}
									position++
									if !rules[rulews]() {
										fmt.Print("243 ")
										goto l243
									}
									if buffer[position] != rune(']') {
										fmt.Print("243 ")
										goto l243
									}
									position++
									if !rules[rulews]() {
										fmt.Print("243 ")
										goto l243
									}
									depth--
									add(ruleanon, position258)
								}
							}
						l248:
							depth--
							add(ruleblankNode, position247)
						}
						break
					case 'F', 'T', 'f', 't':
						{
							position259 := position
							depth++
							{
								position260, tokenIndex260, depth260 := position, tokenIndex, depth
								{
									position262 := position
									depth++
									{
										position263, tokenIndex263, depth263 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											fmt.Print("264 ")
											goto l264
										}
										position++
										fmt.Print("263 ")
										goto l263
									l264:
										position, tokenIndex, depth = position263, tokenIndex263, depth263
										if buffer[position] != rune('T') {
											fmt.Print("261 ")
											goto l261
										}
										position++
									}
								l263:
									{
										position265, tokenIndex265, depth265 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											fmt.Print("266 ")
											goto l266
										}
										position++
										fmt.Print("265 ")
										goto l265
									l266:
										position, tokenIndex, depth = position265, tokenIndex265, depth265
										if buffer[position] != rune('R') {
											fmt.Print("261 ")
											goto l261
										}
										position++
									}
								l265:
									{
										position267, tokenIndex267, depth267 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											fmt.Print("268 ")
											goto l268
										}
										position++
										fmt.Print("267 ")
										goto l267
									l268:
										position, tokenIndex, depth = position267, tokenIndex267, depth267
										if buffer[position] != rune('U') {
											fmt.Print("261 ")
											goto l261
										}
										position++
									}
								l267:
									{
										position269, tokenIndex269, depth269 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											fmt.Print("270 ")
											goto l270
										}
										position++
										fmt.Print("269 ")
										goto l269
									l270:
										position, tokenIndex, depth = position269, tokenIndex269, depth269
										if buffer[position] != rune('E') {
											fmt.Print("261 ")
											goto l261
										}
										position++
									}
								l269:
									if !rules[rulews]() {
										fmt.Print("261 ")
										goto l261
									}
									depth--
									add(ruleTRUE, position262)
								}
								fmt.Print("260 ")
								goto l260
							l261:
								position, tokenIndex, depth = position260, tokenIndex260, depth260
								{
									position271 := position
									depth++
									{
										position272, tokenIndex272, depth272 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											fmt.Print("273 ")
											goto l273
										}
										position++
										fmt.Print("272 ")
										goto l272
									l273:
										position, tokenIndex, depth = position272, tokenIndex272, depth272
										if buffer[position] != rune('F') {
											fmt.Print("243 ")
											goto l243
										}
										position++
									}
								l272:
									{
										position274, tokenIndex274, depth274 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											fmt.Print("275 ")
											goto l275
										}
										position++
										fmt.Print("274 ")
										goto l274
									l275:
										position, tokenIndex, depth = position274, tokenIndex274, depth274
										if buffer[position] != rune('A') {
											fmt.Print("243 ")
											goto l243
										}
										position++
									}
								l274:
									{
										position276, tokenIndex276, depth276 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											fmt.Print("277 ")
											goto l277
										}
										position++
										fmt.Print("276 ")
										goto l276
									l277:
										position, tokenIndex, depth = position276, tokenIndex276, depth276
										if buffer[position] != rune('L') {
											fmt.Print("243 ")
											goto l243
										}
										position++
									}
								l276:
									{
										position278, tokenIndex278, depth278 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											fmt.Print("279 ")
											goto l279
										}
										position++
										fmt.Print("278 ")
										goto l278
									l279:
										position, tokenIndex, depth = position278, tokenIndex278, depth278
										if buffer[position] != rune('S') {
											fmt.Print("243 ")
											goto l243
										}
										position++
									}
								l278:
									{
										position280, tokenIndex280, depth280 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											fmt.Print("281 ")
											goto l281
										}
										position++
										fmt.Print("280 ")
										goto l280
									l281:
										position, tokenIndex, depth = position280, tokenIndex280, depth280
										if buffer[position] != rune('E') {
											fmt.Print("243 ")
											goto l243
										}
										position++
									}
								l280:
									if !rules[rulews]() {
										fmt.Print("243 ")
										goto l243
									}
									depth--
									add(ruleFALSE, position271)
								}
							}
						l260:
							depth--
							add(rulebooleanLiteral, position259)
						}
						break
					case '"':
						{
							position282 := position
							depth++
							{
								position283 := position
								depth++
								if buffer[position] != rune('"') {
									fmt.Print("243 ")
									goto l243
								}
								position++
							l284:
								{
									position285, tokenIndex285, depth285 := position, tokenIndex, depth
									{
										position286, tokenIndex286, depth286 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											fmt.Print("286 ")
											goto l286
										}
										position++
										fmt.Print("285 ")
										goto l285
									l286:
										position, tokenIndex, depth = position286, tokenIndex286, depth286
									}
									if !matchDot() {
										fmt.Print("285 ")
										goto l285
									}
									fmt.Print("284 ")
									goto l284
								l285:
									position, tokenIndex, depth = position285, tokenIndex285, depth285
								}
								if buffer[position] != rune('"') {
									fmt.Print("243 ")
									goto l243
								}
								position++
								depth--
								add(rulestring, position283)
							}
							{
								position287, tokenIndex287, depth287 := position, tokenIndex, depth
								{
									position289, tokenIndex289, depth289 := position, tokenIndex, depth
									if buffer[position] != rune('@') {
										fmt.Print("290 ")
										goto l290
									}
									position++
									{
										position293, tokenIndex293, depth293 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											fmt.Print("294 ")
											goto l294
										}
										position++
										fmt.Print("293 ")
										goto l293
									l294:
										position, tokenIndex, depth = position293, tokenIndex293, depth293
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											fmt.Print("290 ")
											goto l290
										}
										position++
									}
								l293:
								l291:
									{
										position292, tokenIndex292, depth292 := position, tokenIndex, depth
										{
											position295, tokenIndex295, depth295 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												fmt.Print("296 ")
												goto l296
											}
											position++
											fmt.Print("295 ")
											goto l295
										l296:
											position, tokenIndex, depth = position295, tokenIndex295, depth295
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												fmt.Print("292 ")
												goto l292
											}
											position++
										}
									l295:
										fmt.Print("291 ")
										goto l291
									l292:
										position, tokenIndex, depth = position292, tokenIndex292, depth292
									}
								l297:
									{
										position298, tokenIndex298, depth298 := position, tokenIndex, depth
										if buffer[position] != rune('-') {
											fmt.Print("298 ")
											goto l298
										}
										position++
										{
											switch buffer[position] {
											case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													fmt.Print("298 ")
													goto l298
												}
												position++
												break
											case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													fmt.Print("298 ")
													goto l298
												}
												position++
												break
											default:
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													fmt.Print("298 ")
													goto l298
												}
												position++
												break
											}
										}

									l299:
										{
											position300, tokenIndex300, depth300 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														fmt.Print("300 ")
														goto l300
													}
													position++
													break
												case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														fmt.Print("300 ")
														goto l300
													}
													position++
													break
												default:
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														fmt.Print("300 ")
														goto l300
													}
													position++
													break
												}
											}

											fmt.Print("299 ")
											goto l299
										l300:
											position, tokenIndex, depth = position300, tokenIndex300, depth300
										}
										fmt.Print("297 ")
										goto l297
									l298:
										position, tokenIndex, depth = position298, tokenIndex298, depth298
									}
									fmt.Print("289 ")
									goto l289
								l290:
									position, tokenIndex, depth = position289, tokenIndex289, depth289
									if buffer[position] != rune('^') {
										fmt.Print("287 ")
										goto l287
									}
									position++
									if buffer[position] != rune('^') {
										fmt.Print("287 ")
										goto l287
									}
									position++
									if !rules[ruleiri]() {
										fmt.Print("287 ")
										goto l287
									}
								}
							l289:
								fmt.Print("288 ")
								goto l288
							l287:
								position, tokenIndex, depth = position287, tokenIndex287, depth287
							}
						l288:
							if !rules[rulews]() {
								fmt.Print("243 ")
								goto l243
							}
							depth--
							add(ruleliteral, position282)
						}
						break
					case '<':
						if !rules[ruleiri]() {
							fmt.Print("243 ")
							goto l243
						}
						break
					default:
						{
							position303 := position
							depth++
							{
								position304, tokenIndex304, depth304 := position, tokenIndex, depth
								{
									position306, tokenIndex306, depth306 := position, tokenIndex, depth
									if buffer[position] != rune('+') {
										fmt.Print("307 ")
										goto l307
									}
									position++
									fmt.Print("306 ")
									goto l306
								l307:
									position, tokenIndex, depth = position306, tokenIndex306, depth306
									if buffer[position] != rune('-') {
										fmt.Print("304 ")
										goto l304
									}
									position++
								}
							l306:
								fmt.Print("305 ")
								goto l305
							l304:
								position, tokenIndex, depth = position304, tokenIndex304, depth304
							}
						l305:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								fmt.Print("243 ")
								goto l243
							}
							position++
						l308:
							{
								position309, tokenIndex309, depth309 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									fmt.Print("309 ")
									goto l309
								}
								position++
								fmt.Print("308 ")
								goto l308
							l309:
								position, tokenIndex, depth = position309, tokenIndex309, depth309
							}
							{
								position310, tokenIndex310, depth310 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									fmt.Print("310 ")
									goto l310
								}
								position++
							l312:
								{
									position313, tokenIndex313, depth313 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										fmt.Print("313 ")
										goto l313
									}
									position++
									fmt.Print("312 ")
									goto l312
								l313:
									position, tokenIndex, depth = position313, tokenIndex313, depth313
								}
								fmt.Print("311 ")
								goto l311
							l310:
								position, tokenIndex, depth = position310, tokenIndex310, depth310
							}
						l311:
							if !rules[rulews]() {
								fmt.Print("243 ")
								goto l243
							}
							depth--
							add(rulenumericLiteral, position303)
						}
						break
					}
				}

				depth--
				add(rulegraphTerm, position244)
			}
			return true
		l243:
			position, tokenIndex, depth = position243, tokenIndex243, depth243
			return false
		},
		/* 21 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 22 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 23 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 24 propertyListPath <- <(((pof Action3) / (<var> Action4) / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			fmt.Println("\npropertyListPath")
			position317, tokenIndex317, depth317 := position, tokenIndex, depth
			{
				position318 := position
				depth++
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if !rules[rulepof]() {
						fmt.Print("320 ")
						goto l320
					}
					{
						add(ruleAction3, position)
					}
					fmt.Print("319 ")
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					{
						position323 := position
						depth++
						if !rules[rulevar]() {
							fmt.Print("322 ")
							goto l322
						}
						depth--
						add(rulePegText, position323)
					}
					{
						add(ruleAction4, position)
					}
					fmt.Print("319 ")
					goto l319
				l322:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					{
						position325 := position
						depth++
						if !rules[rulepath]() {
							fmt.Print("317 ")
							goto l317
						}
						depth--
						add(ruleverbPath, position325)
					}
				}
			l319:
				if !rules[ruleobjectListPath]() {
					fmt.Print("317 ")
					goto l317
				}
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					{
						position328 := position
						depth++
						if buffer[position] != rune(';') {
							fmt.Print("326 ")
							goto l326
						}
						position++
						if !rules[rulews]() {
							fmt.Print("326 ")
							goto l326
						}
						depth--
						add(ruleSEMICOLON, position328)
					}
					if !rules[rulepropertyListPath]() {
						fmt.Print("326 ")
						goto l326
					}
					fmt.Print("327 ")
					goto l327
				l326:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
				}
			l327:
				depth--
				add(rulepropertyListPath, position318)
			}
			return true
		l317:
			position, tokenIndex, depth = position317, tokenIndex317, depth317
			return false
		},
		/* 25 verbPath <- <path> */
		nil,
		/* 26 path <- <pathAlternative> */
		func() bool {
			fmt.Println("\npath")
			position330, tokenIndex330, depth330 := position, tokenIndex, depth
			{
				position331 := position
				depth++
				if !rules[rulepathAlternative]() {
					fmt.Print("330 ")
					goto l330
				}
				depth--
				add(rulepath, position331)
			}
			return true
		l330:
			position, tokenIndex, depth = position330, tokenIndex330, depth330
			return false
		},
		/* 27 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			fmt.Println("\npathAlternative")
			position332, tokenIndex332, depth332 := position, tokenIndex, depth
			{
				position333 := position
				depth++
				if !rules[rulepathSequence]() {
					fmt.Print("332 ")
					goto l332
				}
			l334:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						fmt.Print("335 ")
						goto l335
					}
					if !rules[rulepathAlternative]() {
						fmt.Print("335 ")
						goto l335
					}
					fmt.Print("334 ")
					goto l334
				l335:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
				}
				depth--
				add(rulepathAlternative, position333)
			}
			return true
		l332:
			position, tokenIndex, depth = position332, tokenIndex332, depth332
			return false
		},
		/* 28 pathSequence <- <(<pathElt> Action5 (SLASH pathSequence)*)> */
		func() bool {
			fmt.Println("\npathSequence")
			position336, tokenIndex336, depth336 := position, tokenIndex, depth
			{
				position337 := position
				depth++
				{
					position338 := position
					depth++
					{
						position339 := position
						depth++
						{
							position340, tokenIndex340, depth340 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								fmt.Print("340 ")
								goto l340
							}
							fmt.Print("341 ")
							goto l341
						l340:
							position, tokenIndex, depth = position340, tokenIndex340, depth340
						}
					l341:
						{
							position342 := position
							depth++
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										fmt.Print("336 ")
										goto l336
									}
									if !rules[rulepath]() {
										fmt.Print("336 ")
										goto l336
									}
									if !rules[ruleRPAREN]() {
										fmt.Print("336 ")
										goto l336
									}
									break
								case '!':
									{
										position344 := position
										depth++
										if buffer[position] != rune('!') {
											fmt.Print("336 ")
											goto l336
										}
										position++
										if !rules[rulews]() {
											fmt.Print("336 ")
											goto l336
										}
										depth--
										add(ruleNOT, position344)
									}
									{
										position345 := position
										depth++
										{
											position346, tokenIndex346, depth346 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												fmt.Print("347 ")
												goto l347
											}
											fmt.Print("346 ")
											goto l346
										l347:
											position, tokenIndex, depth = position346, tokenIndex346, depth346
											if !rules[ruleLPAREN]() {
												fmt.Print("336 ")
												goto l336
											}
											{
												position348, tokenIndex348, depth348 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													fmt.Print("348 ")
													goto l348
												}
											l350:
												{
													position351, tokenIndex351, depth351 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														fmt.Print("351 ")
														goto l351
													}
													if !rules[rulepathOneInPropertySet]() {
														fmt.Print("351 ")
														goto l351
													}
													fmt.Print("350 ")
													goto l350
												l351:
													position, tokenIndex, depth = position351, tokenIndex351, depth351
												}
												fmt.Print("349 ")
												goto l349
											l348:
												position, tokenIndex, depth = position348, tokenIndex348, depth348
											}
										l349:
											if !rules[ruleRPAREN]() {
												fmt.Print("336 ")
												goto l336
											}
										}
									l346:
										depth--
										add(rulepathNegatedPropertySet, position345)
									}
									break
								case 'a':
									if !rules[ruleISA]() {
										fmt.Print("336 ")
										goto l336
									}
									break
								default:
									if !rules[ruleiri]() {
										fmt.Print("336 ")
										goto l336
									}
									break
								}
							}

							depth--
							add(rulepathPrimary, position342)
						}
						depth--
						add(rulepathElt, position339)
					}
					depth--
					add(rulePegText, position338)
				}
				{
					add(ruleAction5, position)
				}
			l353:
				{
					position354, tokenIndex354, depth354 := position, tokenIndex, depth
					{
						position355 := position
						depth++
						if buffer[position] != rune('/') {
							fmt.Print("354 ")
							goto l354
						}
						position++
						if !rules[rulews]() {
							fmt.Print("354 ")
							goto l354
						}
						depth--
						add(ruleSLASH, position355)
					}
					if !rules[rulepathSequence]() {
						fmt.Print("354 ")
						goto l354
					}
					fmt.Print("353 ")
					goto l353
				l354:
					position, tokenIndex, depth = position354, tokenIndex354, depth354
				}
				depth--
				add(rulepathSequence, position337)
			}
			return true
		l336:
			position, tokenIndex, depth = position336, tokenIndex336, depth336
			return false
		},
		/* 29 pathElt <- <(INVERSE? pathPrimary)> */
		nil,
		/* 30 pathPrimary <- <((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA) | (&('<') iri))> */
		nil,
		/* 31 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 32 pathOneInPropertySet <- <((&('^') (INVERSE (iri / ISA))) | (&('a') ISA) | (&('<') iri))> */
		func() bool {
			fmt.Println("\npathOneInPropertySet")
			position359, tokenIndex359, depth359 := position, tokenIndex, depth
			{
				position360 := position
				depth++
				{
					switch buffer[position] {
					case '^':
						if !rules[ruleINVERSE]() {
							fmt.Print("359 ")
							goto l359
						}
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if !rules[ruleiri]() {
								fmt.Print("363 ")
								goto l363
							}
							fmt.Print("362 ")
							goto l362
						l363:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
							if !rules[ruleISA]() {
								fmt.Print("359 ")
								goto l359
							}
						}
					l362:
						break
					case 'a':
						if !rules[ruleISA]() {
							fmt.Print("359 ")
							goto l359
						}
						break
					default:
						if !rules[ruleiri]() {
							fmt.Print("359 ")
							goto l359
						}
						break
					}
				}

				depth--
				add(rulepathOneInPropertySet, position360)
			}
			return true
		l359:
			position, tokenIndex, depth = position359, tokenIndex359, depth359
			return false
		},
		/* 33 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			fmt.Println("\nobjectListPath")
			{
				position365 := position
				depth++
				{
					position366 := position
					depth++
					{
						position367, tokenIndex367, depth367 := position, tokenIndex, depth
						{
							position369 := position
							depth++
							if !rules[rulegraphNodePath]() {
								fmt.Print("368 ")
								goto l368
							}
							depth--
							add(rulePegText, position369)
						}
						{
							add(ruleAction6, position)
						}
						fmt.Print("367 ")
						goto l367
					l368:
						position, tokenIndex, depth = position367, tokenIndex367, depth367
						if !rules[rulepof]() {
							fmt.Print("371 ")
							goto l371
						}
						{
							add(ruleAction7, position)
						}
						fmt.Print("367 ")
						goto l367
					l371:
						position, tokenIndex, depth = position367, tokenIndex367, depth367
						{
							add(ruleAction8, position)
						}
					}
				l367:
					depth--
					add(ruleobjectPath, position366)
				}
			l374:
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					{
						position376 := position
						depth++
						if buffer[position] != rune(',') {
							fmt.Print("375 ")
							goto l375
						}
						position++
						if !rules[rulews]() {
							fmt.Print("375 ")
							goto l375
						}
						depth--
						add(ruleCOMMA, position376)
					}
					if !rules[ruleobjectListPath]() {
						fmt.Print("375 ")
						goto l375
					}
					fmt.Print("374 ")
					goto l374
				l375:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
				}
				depth--
				add(ruleobjectListPath, position365)
			}
			return true
		},
		/* 34 objectPath <- <((<graphNodePath> Action6) / (pof Action7) / Action8)> */
		nil,
		/* 35 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			fmt.Println("\ngraphNodePath")
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if !rules[rulevar]() {
						fmt.Print("381 ")
						goto l381
					}
					fmt.Print("380 ")
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					if !rules[rulegraphTerm]() {
						fmt.Print("378 ")
						goto l378
					}
				}
			l380:
				depth--
				add(rulegraphNodePath, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 36 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 37 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 38 limit <- <(LIMIT INTEGER)> */
		func() bool {
			fmt.Println("\nlimit")
			position384, tokenIndex384, depth384 := position, tokenIndex, depth
			{
				position385 := position
				depth++
				{
					position386 := position
					depth++
					{
						position387, tokenIndex387, depth387 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							fmt.Print("388 ")
							goto l388
						}
						position++
						fmt.Print("387 ")
						goto l387
					l388:
						position, tokenIndex, depth = position387, tokenIndex387, depth387
						if buffer[position] != rune('L') {
							fmt.Print("384 ")
							goto l384
						}
						position++
					}
				l387:
					{
						position389, tokenIndex389, depth389 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							fmt.Print("390 ")
							goto l390
						}
						position++
						fmt.Print("389 ")
						goto l389
					l390:
						position, tokenIndex, depth = position389, tokenIndex389, depth389
						if buffer[position] != rune('I') {
							fmt.Print("384 ")
							goto l384
						}
						position++
					}
				l389:
					{
						position391, tokenIndex391, depth391 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							fmt.Print("392 ")
							goto l392
						}
						position++
						fmt.Print("391 ")
						goto l391
					l392:
						position, tokenIndex, depth = position391, tokenIndex391, depth391
						if buffer[position] != rune('M') {
							fmt.Print("384 ")
							goto l384
						}
						position++
					}
				l391:
					{
						position393, tokenIndex393, depth393 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							fmt.Print("394 ")
							goto l394
						}
						position++
						fmt.Print("393 ")
						goto l393
					l394:
						position, tokenIndex, depth = position393, tokenIndex393, depth393
						if buffer[position] != rune('I') {
							fmt.Print("384 ")
							goto l384
						}
						position++
					}
				l393:
					{
						position395, tokenIndex395, depth395 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							fmt.Print("396 ")
							goto l396
						}
						position++
						fmt.Print("395 ")
						goto l395
					l396:
						position, tokenIndex, depth = position395, tokenIndex395, depth395
						if buffer[position] != rune('T') {
							fmt.Print("384 ")
							goto l384
						}
						position++
					}
				l395:
					if !rules[rulews]() {
						fmt.Print("384 ")
						goto l384
					}
					depth--
					add(ruleLIMIT, position386)
				}
				if !rules[ruleINTEGER]() {
					fmt.Print("384 ")
					goto l384
				}
				depth--
				add(rulelimit, position385)
			}
			return true
		l384:
			position, tokenIndex, depth = position384, tokenIndex384, depth384
			return false
		},
		/* 39 offset <- <(OFFSET INTEGER)> */
		func() bool {
			fmt.Println("\noffset")
			position397, tokenIndex397, depth397 := position, tokenIndex, depth
			{
				position398 := position
				depth++
				{
					position399 := position
					depth++
					{
						position400, tokenIndex400, depth400 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							fmt.Print("401 ")
							goto l401
						}
						position++
						fmt.Print("400 ")
						goto l400
					l401:
						position, tokenIndex, depth = position400, tokenIndex400, depth400
						if buffer[position] != rune('O') {
							fmt.Print("397 ")
							goto l397
						}
						position++
					}
				l400:
					{
						position402, tokenIndex402, depth402 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							fmt.Print("403 ")
							goto l403
						}
						position++
						fmt.Print("402 ")
						goto l402
					l403:
						position, tokenIndex, depth = position402, tokenIndex402, depth402
						if buffer[position] != rune('F') {
							fmt.Print("397 ")
							goto l397
						}
						position++
					}
				l402:
					{
						position404, tokenIndex404, depth404 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							fmt.Print("405 ")
							goto l405
						}
						position++
						fmt.Print("404 ")
						goto l404
					l405:
						position, tokenIndex, depth = position404, tokenIndex404, depth404
						if buffer[position] != rune('F') {
							fmt.Print("397 ")
							goto l397
						}
						position++
					}
				l404:
					{
						position406, tokenIndex406, depth406 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							fmt.Print("407 ")
							goto l407
						}
						position++
						fmt.Print("406 ")
						goto l406
					l407:
						position, tokenIndex, depth = position406, tokenIndex406, depth406
						if buffer[position] != rune('S') {
							fmt.Print("397 ")
							goto l397
						}
						position++
					}
				l406:
					{
						position408, tokenIndex408, depth408 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							fmt.Print("409 ")
							goto l409
						}
						position++
						fmt.Print("408 ")
						goto l408
					l409:
						position, tokenIndex, depth = position408, tokenIndex408, depth408
						if buffer[position] != rune('E') {
							fmt.Print("397 ")
							goto l397
						}
						position++
					}
				l408:
					{
						position410, tokenIndex410, depth410 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							fmt.Print("411 ")
							goto l411
						}
						position++
						fmt.Print("410 ")
						goto l410
					l411:
						position, tokenIndex, depth = position410, tokenIndex410, depth410
						if buffer[position] != rune('T') {
							fmt.Print("397 ")
							goto l397
						}
						position++
					}
				l410:
					if !rules[rulews]() {
						fmt.Print("397 ")
						goto l397
					}
					depth--
					add(ruleOFFSET, position399)
				}
				if !rules[ruleINTEGER]() {
					fmt.Print("397 ")
					goto l397
				}
				depth--
				add(ruleoffset, position398)
			}
			return true
		l397:
			position, tokenIndex, depth = position397, tokenIndex397, depth397
			return false
		},
		/* 40 pof <- <('<' ((&('\f') '\f') | (&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))+)> */
		func() bool {
			fmt.Println("\npof")
			position412, tokenIndex412, depth412 := position, tokenIndex, depth
			{
				position413 := position
				depth++
				if buffer[position] != rune('<') {
					fmt.Print("412 ")
					goto l412
				}
				position++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							fmt.Print("412 ")
							goto l412
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							fmt.Print("412 ")
							goto l412
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							fmt.Print("412 ")
							goto l412
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							fmt.Print("412 ")
							goto l412
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							fmt.Print("412 ")
							goto l412
						}
						position++
						break
					}
				}

			l414:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								fmt.Print("415 ")
								goto l415
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								fmt.Print("415 ")
								goto l415
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								fmt.Print("415 ")
								goto l415
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								fmt.Print("415 ")
								goto l415
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								fmt.Print("415 ")
								goto l415
							}
							position++
							break
						}
					}

					fmt.Print("414 ")
					goto l414
				l415:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
				}
				depth--
				add(rulepof, position413)
			}
			return true
		l412:
			position, tokenIndex, depth = position412, tokenIndex412, depth412
			return false
		},
		/* 41 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			fmt.Println("\nvar")
			position418, tokenIndex418, depth418 := position, tokenIndex, depth
			{
				position419 := position
				depth++
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						fmt.Print("421 ")
						goto l421
					}
					position++
					fmt.Print("420 ")
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if buffer[position] != rune('$') {
						fmt.Print("418 ")
						goto l418
					}
					position++
				}
			l420:
				{
					position422 := position
					depth++
					{
						position425, tokenIndex425, depth425 := position, tokenIndex, depth
						{
							position427 := position
							depth++
							{
								position428, tokenIndex428, depth428 := position, tokenIndex, depth
								{
									position430 := position
									depth++
									{
										position431, tokenIndex431, depth431 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											fmt.Print("432 ")
											goto l432
										}
										position++
										fmt.Print("431 ")
										goto l431
									l432:
										position, tokenIndex, depth = position431, tokenIndex431, depth431
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											fmt.Print("429 ")
											goto l429
										}
										position++
									}
								l431:
									depth--
									add(rulePN_CHARS_BASE, position430)
								}
								fmt.Print("428 ")
								goto l428
							l429:
								position, tokenIndex, depth = position428, tokenIndex428, depth428
								if buffer[position] != rune('_') {
									fmt.Print("426 ")
									goto l426
								}
								position++
							}
						l428:
							depth--
							add(rulePN_CHARS_U, position427)
						}
						fmt.Print("425 ")
						goto l425
					l426:
						position, tokenIndex, depth = position425, tokenIndex425, depth425
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							fmt.Print("418 ")
							goto l418
						}
						position++
					}
				l425:
				l423:
					{
						position424, tokenIndex424, depth424 := position, tokenIndex, depth
						{
							position433, tokenIndex433, depth433 := position, tokenIndex, depth
							{
								position435 := position
								depth++
								{
									position436, tokenIndex436, depth436 := position, tokenIndex, depth
									{
										position438 := position
										depth++
										{
											position439, tokenIndex439, depth439 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												fmt.Print("440 ")
												goto l440
											}
											position++
											fmt.Print("439 ")
											goto l439
										l440:
											position, tokenIndex, depth = position439, tokenIndex439, depth439
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												fmt.Print("437 ")
												goto l437
											}
											position++
										}
									l439:
										depth--
										add(rulePN_CHARS_BASE, position438)
									}
									fmt.Print("436 ")
									goto l436
								l437:
									position, tokenIndex, depth = position436, tokenIndex436, depth436
									if buffer[position] != rune('_') {
										fmt.Print("434 ")
										goto l434
									}
									position++
								}
							l436:
								depth--
								add(rulePN_CHARS_U, position435)
							}
							fmt.Print("433 ")
							goto l433
						l434:
							position, tokenIndex, depth = position433, tokenIndex433, depth433
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								fmt.Print("424 ")
								goto l424
							}
							position++
						}
					l433:
						fmt.Print("423 ")
						goto l423
					l424:
						position, tokenIndex, depth = position424, tokenIndex424, depth424
					}
					depth--
					add(ruleVARNAME, position422)
				}
				if !rules[rulews]() {
					fmt.Print("418 ")
					goto l418
				}
				depth--
				add(rulevar, position419)
			}
			return true
		l418:
			position, tokenIndex, depth = position418, tokenIndex418, depth418
			return false
		},
		/* 42 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			fmt.Println("\niri")
			position441, tokenIndex441, depth441 := position, tokenIndex, depth
			{
				position442 := position
				depth++
				if buffer[position] != rune('<') {
					fmt.Print("441 ")
					goto l441
				}
				position++
			l443:
				{
					position444, tokenIndex444, depth444 := position, tokenIndex, depth
					{
						position445, tokenIndex445, depth445 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							fmt.Print("445 ")
							goto l445
						}
						position++
						fmt.Print("444 ")
						goto l444
					l445:
						position, tokenIndex, depth = position445, tokenIndex445, depth445
					}
					if !matchDot() {
						fmt.Print("444 ")
						goto l444
					}
					fmt.Print("443 ")
					goto l443
				l444:
					position, tokenIndex, depth = position444, tokenIndex444, depth444
				}
				if buffer[position] != rune('>') {
					fmt.Print("441 ")
					goto l441
				}
				position++
				if !rules[rulews]() {
					fmt.Print("441 ")
					goto l441
				}
				depth--
				add(ruleiri, position442)
			}
			return true
		l441:
			position, tokenIndex, depth = position441, tokenIndex441, depth441
			return false
		},
		/* 43 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iri))? ws)> */
		nil,
		/* 44 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 45 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? ws)> */
		nil,
		/* 46 booleanLiteral <- <(TRUE / FALSE)> */
		nil,
		/* 47 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 48 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? ws)> */
		nil,
		/* 49 anon <- <('[' ws ']' ws)> */
		nil,
		/* 50 nil <- <('(' ws ')' ws)> */
		nil,
		/* 51 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 52 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 53 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 54 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') ws)> */
		nil,
		/* 55 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') ws)> */
		nil,
		/* 56 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 57 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 58 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 59 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 60 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 61 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') ws)> */
		nil,
		/* 62 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 63 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') ws)> */
		nil,
		/* 64 LBRACE <- <('{' ws)> */
		func() bool {
			fmt.Println("\nLBRACE")
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				if buffer[position] != rune('{') {
					fmt.Print("467 ")
					goto l467
				}
				position++
				if !rules[rulews]() {
					fmt.Print("467 ")
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
		/* 65 RBRACE <- <('}' ws)> */
		func() bool {
			fmt.Println("\nRBRACE")
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				if buffer[position] != rune('}') {
					fmt.Print("469 ")
					goto l469
				}
				position++
				if !rules[rulews]() {
					fmt.Print("469 ")
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
		/* 66 LBRACK <- <('[' ws)> */
		nil,
		/* 67 RBRACK <- <(']' ws)> */
		nil,
		/* 68 SEMICOLON <- <(';' ws)> */
		nil,
		/* 69 COMMA <- <(',' ws)> */
		nil,
		/* 70 DOT <- <('.' ws)> */
		func() bool {
			fmt.Println("\nDOT")
			position475, tokenIndex475, depth475 := position, tokenIndex, depth
			{
				position476 := position
				depth++
				if buffer[position] != rune('.') {
					fmt.Print("475 ")
					goto l475
				}
				position++
				if !rules[rulews]() {
					fmt.Print("475 ")
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
		/* 71 COLON <- <(':' ws)> */
		nil,
		/* 72 PIPE <- <('|' ws)> */
		func() bool {
			fmt.Println("\nPIPE")
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				if buffer[position] != rune('|') {
					fmt.Print("478 ")
					goto l478
				}
				position++
				if !rules[rulews]() {
					fmt.Print("478 ")
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
		/* 73 SLASH <- <('/' ws)> */
		nil,
		/* 74 INVERSE <- <('^' ws)> */
		func() bool {
			fmt.Println("\nINVERSE")
			position481, tokenIndex481, depth481 := position, tokenIndex, depth
			{
				position482 := position
				depth++
				if buffer[position] != rune('^') {
					fmt.Print("481 ")
					goto l481
				}
				position++
				if !rules[rulews]() {
					fmt.Print("481 ")
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
		/* 75 LPAREN <- <('(' ws)> */
		func() bool {
			fmt.Println("\nLPAREN")
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				if buffer[position] != rune('(') {
					fmt.Print("483 ")
					goto l483
				}
				position++
				if !rules[rulews]() {
					fmt.Print("483 ")
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
		/* 76 RPAREN <- <(')' ws)> */
		func() bool {
			fmt.Println("\nRPAREN")
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				if buffer[position] != rune(')') {
					fmt.Print("485 ")
					goto l485
				}
				position++
				if !rules[rulews]() {
					fmt.Print("485 ")
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
		/* 77 ISA <- <('a' ws)> */
		func() bool {
			fmt.Println("\nISA")
			position487, tokenIndex487, depth487 := position, tokenIndex, depth
			{
				position488 := position
				depth++
				if buffer[position] != rune('a') {
					fmt.Print("487 ")
					goto l487
				}
				position++
				if !rules[rulews]() {
					fmt.Print("487 ")
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
		/* 78 NOT <- <('!' ws)> */
		nil,
		/* 79 STAR <- <('*' ws)> */
		nil,
		/* 80 QUESTION <- <('?' ws)> */
		nil,
		/* 81 PLUS <- <('+' ws)> */
		nil,
		/* 82 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') ws)> */
		nil,
		/* 83 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') ws)> */
		nil,
		/* 84 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') ws)> */
		nil,
		/* 85 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') ws)> */
		nil,
		/* 86 INTEGER <- <([0-9]+ ws)> */
		func() bool {
			fmt.Println("\nINTEGER")
			position497, tokenIndex497, depth497 := position, tokenIndex, depth
			{
				position498 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					fmt.Print("497 ")
					goto l497
				}
				position++
			l499:
				{
					position500, tokenIndex500, depth500 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						fmt.Print("500 ")
						goto l500
					}
					position++
					fmt.Print("499 ")
					goto l499
				l500:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
				}
				if !rules[rulews]() {
					fmt.Print("497 ")
					goto l497
				}
				depth--
				add(ruleINTEGER, position498)
			}
			return true
		l497:
			position, tokenIndex, depth = position497, tokenIndex497, depth497
			return false
		},
		/* 87 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			fmt.Println("\nws")
			{
				position502 := position
				depth++
			l503:
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\v':
							if buffer[position] != rune('\v') {
								fmt.Print("504 ")
								goto l504
							}
							position++
							break
						case '\f':
							if buffer[position] != rune('\f') {
								fmt.Print("504 ")
								goto l504
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								fmt.Print("504 ")
								goto l504
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								fmt.Print("504 ")
								goto l504
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								fmt.Print("504 ")
								goto l504
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								fmt.Print("504 ")
								goto l504
							}
							position++
							break
						}
					}

					fmt.Print("503 ")
					goto l503
				l504:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
				}
				depth--
				add(rulews, position502)
			}
			return true
		},
		nil,
		/* 90 Action0 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 91 Action1 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 92 Action2 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 93 Action3 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 94 Action4 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 95 Action5 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 96 Action6 <- <{ p.setObject(buffer[begin:end]); p.addTriplePattern() }> */
		nil,
		/* 97 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 98 Action8 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
	}
	p.rules = rules
}
