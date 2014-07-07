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
	rulepathMod
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
	ruleVAR_CHAR
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
	"pathMod",
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
	"VAR_CHAR",
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
	Bgp

	Buffer string
	buffer []rune
	rules  [101]func() bool
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
										if buffer[position] != rune('\'') {
											goto l25
										}
										position++
										goto l24
									l25:
										position, tokenIndex, depth = position24, tokenIndex24, depth24
										{
											switch buffer[position] {
											case ' ':
												if buffer[position] != rune(' ') {
													goto l23
												}
												position++
												break
											case '\'':
												if buffer[position] != rune('\'') {
													goto l23
												}
												position++
												break
											default:
												if buffer[position] != rune(':') {
													goto l23
												}
												position++
												break
											}
										}

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
										position27, tokenIndex27, depth27 := position, tokenIndex, depth
										{
											position28, tokenIndex28, depth28 := position, tokenIndex, depth
											if buffer[position] != rune('\'') {
												goto l29
											}
											position++
											goto l28
										l29:
											position, tokenIndex, depth = position28, tokenIndex28, depth28
											{
												switch buffer[position] {
												case ' ':
													if buffer[position] != rune(' ') {
														goto l27
													}
													position++
													break
												case '\'':
													if buffer[position] != rune('\'') {
														goto l27
													}
													position++
													break
												default:
													if buffer[position] != rune(':') {
														goto l27
													}
													position++
													break
												}
											}

										}
									l28:
										goto l22
									l27:
										position, tokenIndex, depth = position27, tokenIndex27, depth27
									}
									if !matchDot() {
										goto l22
									}
									goto l21
								l22:
									position, tokenIndex, depth = position22, tokenIndex22, depth22
								}
								{
									position31 := position
									depth++
									if buffer[position] != rune(':') {
										goto l6
									}
									position++
									if !rules[rulews]() {
										goto l6
									}
									depth--
									add(ruleCOLON, position31)
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
								position32 := position
								depth++
								{
									position33 := position
									depth++
									{
										position34, tokenIndex34, depth34 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l35
										}
										position++
										goto l34
									l35:
										position, tokenIndex, depth = position34, tokenIndex34, depth34
										if buffer[position] != rune('B') {
											goto l4
										}
										position++
									}
								l34:
									{
										position36, tokenIndex36, depth36 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l37
										}
										position++
										goto l36
									l37:
										position, tokenIndex, depth = position36, tokenIndex36, depth36
										if buffer[position] != rune('A') {
											goto l4
										}
										position++
									}
								l36:
									{
										position38, tokenIndex38, depth38 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l39
										}
										position++
										goto l38
									l39:
										position, tokenIndex, depth = position38, tokenIndex38, depth38
										if buffer[position] != rune('S') {
											goto l4
										}
										position++
									}
								l38:
									{
										position40, tokenIndex40, depth40 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l41
										}
										position++
										goto l40
									l41:
										position, tokenIndex, depth = position40, tokenIndex40, depth40
										if buffer[position] != rune('E') {
											goto l4
										}
										position++
									}
								l40:
									if !rules[rulews]() {
										goto l4
									}
									depth--
									add(ruleBASE, position33)
								}
								if !rules[ruleiri]() {
									goto l4
								}
								depth--
								add(rulebaseDecl, position32)
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
					position42 := position
					depth++
					{
						position43 := position
						depth++
						if !rules[ruleselect]() {
							goto l0
						}
						{
							position44, tokenIndex44, depth44 := position, tokenIndex, depth
							{
								position46 := position
								depth++
								{
									position47 := position
									depth++
									{
										position48, tokenIndex48, depth48 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l49
										}
										position++
										goto l48
									l49:
										position, tokenIndex, depth = position48, tokenIndex48, depth48
										if buffer[position] != rune('F') {
											goto l44
										}
										position++
									}
								l48:
									{
										position50, tokenIndex50, depth50 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l51
										}
										position++
										goto l50
									l51:
										position, tokenIndex, depth = position50, tokenIndex50, depth50
										if buffer[position] != rune('R') {
											goto l44
										}
										position++
									}
								l50:
									{
										position52, tokenIndex52, depth52 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l53
										}
										position++
										goto l52
									l53:
										position, tokenIndex, depth = position52, tokenIndex52, depth52
										if buffer[position] != rune('O') {
											goto l44
										}
										position++
									}
								l52:
									{
										position54, tokenIndex54, depth54 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l55
										}
										position++
										goto l54
									l55:
										position, tokenIndex, depth = position54, tokenIndex54, depth54
										if buffer[position] != rune('M') {
											goto l44
										}
										position++
									}
								l54:
									if !rules[rulews]() {
										goto l44
									}
									depth--
									add(ruleFROM, position47)
								}
								{
									position56, tokenIndex56, depth56 := position, tokenIndex, depth
									{
										position58 := position
										depth++
										{
											position59, tokenIndex59, depth59 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l60
											}
											position++
											goto l59
										l60:
											position, tokenIndex, depth = position59, tokenIndex59, depth59
											if buffer[position] != rune('N') {
												goto l56
											}
											position++
										}
									l59:
										{
											position61, tokenIndex61, depth61 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l62
											}
											position++
											goto l61
										l62:
											position, tokenIndex, depth = position61, tokenIndex61, depth61
											if buffer[position] != rune('A') {
												goto l56
											}
											position++
										}
									l61:
										{
											position63, tokenIndex63, depth63 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l64
											}
											position++
											goto l63
										l64:
											position, tokenIndex, depth = position63, tokenIndex63, depth63
											if buffer[position] != rune('M') {
												goto l56
											}
											position++
										}
									l63:
										{
											position65, tokenIndex65, depth65 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l66
											}
											position++
											goto l65
										l66:
											position, tokenIndex, depth = position65, tokenIndex65, depth65
											if buffer[position] != rune('E') {
												goto l56
											}
											position++
										}
									l65:
										{
											position67, tokenIndex67, depth67 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l68
											}
											position++
											goto l67
										l68:
											position, tokenIndex, depth = position67, tokenIndex67, depth67
											if buffer[position] != rune('D') {
												goto l56
											}
											position++
										}
									l67:
										if !rules[rulews]() {
											goto l56
										}
										depth--
										add(ruleNAMED, position58)
									}
									goto l57
								l56:
									position, tokenIndex, depth = position56, tokenIndex56, depth56
								}
							l57:
								if !rules[ruleiri]() {
									goto l44
								}
								depth--
								add(ruledatasetClause, position46)
							}
							goto l45
						l44:
							position, tokenIndex, depth = position44, tokenIndex44, depth44
						}
					l45:
						if !rules[rulewhereClause]() {
							goto l0
						}
						{
							position69 := position
							depth++
							{
								position70, tokenIndex70, depth70 := position, tokenIndex, depth
								{
									position72 := position
									depth++
									{
										position73, tokenIndex73, depth73 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l74
										}
										{
											position75, tokenIndex75, depth75 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l75
											}
											goto l76
										l75:
											position, tokenIndex, depth = position75, tokenIndex75, depth75
										}
									l76:
										goto l73
									l74:
										position, tokenIndex, depth = position73, tokenIndex73, depth73
										if !rules[ruleoffset]() {
											goto l70
										}
										{
											position77, tokenIndex77, depth77 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l77
											}
											goto l78
										l77:
											position, tokenIndex, depth = position77, tokenIndex77, depth77
										}
									l78:
									}
								l73:
									depth--
									add(rulelimitOffsetClauses, position72)
								}
								goto l71
							l70:
								position, tokenIndex, depth = position70, tokenIndex70, depth70
							}
						l71:
							depth--
							add(rulesolutionModifier, position69)
						}
						depth--
						add(ruleselectQuery, position43)
					}
					depth--
					add(rulequery, position42)
				}
				{
					position79, tokenIndex79, depth79 := position, tokenIndex, depth
					if !matchDot() {
						goto l79
					}
					goto l0
				l79:
					position, tokenIndex, depth = position79, tokenIndex79, depth79
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
		/* 2 prefixDecl <- <(PREFIX (!('\'' / ((&(' ') ' ') | (&('\'') '\'') | (&(':') ':'))) .)+ COLON iri)> */
		nil,
		/* 3 baseDecl <- <(BASE iri)> */
		nil,
		/* 4 query <- <selectQuery> */
		nil,
		/* 5 selectQuery <- <(select datasetClause? whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position85, tokenIndex85, depth85 := position, tokenIndex, depth
			{
				position86 := position
				depth++
				{
					position87 := position
					depth++
					{
						position88, tokenIndex88, depth88 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l89
						}
						position++
						goto l88
					l89:
						position, tokenIndex, depth = position88, tokenIndex88, depth88
						if buffer[position] != rune('S') {
							goto l85
						}
						position++
					}
				l88:
					{
						position90, tokenIndex90, depth90 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l91
						}
						position++
						goto l90
					l91:
						position, tokenIndex, depth = position90, tokenIndex90, depth90
						if buffer[position] != rune('E') {
							goto l85
						}
						position++
					}
				l90:
					{
						position92, tokenIndex92, depth92 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l93
						}
						position++
						goto l92
					l93:
						position, tokenIndex, depth = position92, tokenIndex92, depth92
						if buffer[position] != rune('L') {
							goto l85
						}
						position++
					}
				l92:
					{
						position94, tokenIndex94, depth94 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l95
						}
						position++
						goto l94
					l95:
						position, tokenIndex, depth = position94, tokenIndex94, depth94
						if buffer[position] != rune('E') {
							goto l85
						}
						position++
					}
				l94:
					{
						position96, tokenIndex96, depth96 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l97
						}
						position++
						goto l96
					l97:
						position, tokenIndex, depth = position96, tokenIndex96, depth96
						if buffer[position] != rune('C') {
							goto l85
						}
						position++
					}
				l96:
					{
						position98, tokenIndex98, depth98 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l99
						}
						position++
						goto l98
					l99:
						position, tokenIndex, depth = position98, tokenIndex98, depth98
						if buffer[position] != rune('T') {
							goto l85
						}
						position++
					}
				l98:
					if !rules[rulews]() {
						goto l85
					}
					depth--
					add(ruleSELECT, position87)
				}
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					{
						position102, tokenIndex102, depth102 := position, tokenIndex, depth
						{
							position104 := position
							depth++
							{
								position105, tokenIndex105, depth105 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l106
								}
								position++
								goto l105
							l106:
								position, tokenIndex, depth = position105, tokenIndex105, depth105
								if buffer[position] != rune('D') {
									goto l103
								}
								position++
							}
						l105:
							{
								position107, tokenIndex107, depth107 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l108
								}
								position++
								goto l107
							l108:
								position, tokenIndex, depth = position107, tokenIndex107, depth107
								if buffer[position] != rune('I') {
									goto l103
								}
								position++
							}
						l107:
							{
								position109, tokenIndex109, depth109 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l110
								}
								position++
								goto l109
							l110:
								position, tokenIndex, depth = position109, tokenIndex109, depth109
								if buffer[position] != rune('S') {
									goto l103
								}
								position++
							}
						l109:
							{
								position111, tokenIndex111, depth111 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l112
								}
								position++
								goto l111
							l112:
								position, tokenIndex, depth = position111, tokenIndex111, depth111
								if buffer[position] != rune('T') {
									goto l103
								}
								position++
							}
						l111:
							{
								position113, tokenIndex113, depth113 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l114
								}
								position++
								goto l113
							l114:
								position, tokenIndex, depth = position113, tokenIndex113, depth113
								if buffer[position] != rune('I') {
									goto l103
								}
								position++
							}
						l113:
							{
								position115, tokenIndex115, depth115 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l116
								}
								position++
								goto l115
							l116:
								position, tokenIndex, depth = position115, tokenIndex115, depth115
								if buffer[position] != rune('N') {
									goto l103
								}
								position++
							}
						l115:
							{
								position117, tokenIndex117, depth117 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l118
								}
								position++
								goto l117
							l118:
								position, tokenIndex, depth = position117, tokenIndex117, depth117
								if buffer[position] != rune('C') {
									goto l103
								}
								position++
							}
						l117:
							{
								position119, tokenIndex119, depth119 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l120
								}
								position++
								goto l119
							l120:
								position, tokenIndex, depth = position119, tokenIndex119, depth119
								if buffer[position] != rune('T') {
									goto l103
								}
								position++
							}
						l119:
							if !rules[rulews]() {
								goto l103
							}
							depth--
							add(ruleDISTINCT, position104)
						}
						goto l102
					l103:
						position, tokenIndex, depth = position102, tokenIndex102, depth102
						{
							position121 := position
							depth++
							{
								position122, tokenIndex122, depth122 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l123
								}
								position++
								goto l122
							l123:
								position, tokenIndex, depth = position122, tokenIndex122, depth122
								if buffer[position] != rune('R') {
									goto l100
								}
								position++
							}
						l122:
							{
								position124, tokenIndex124, depth124 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l125
								}
								position++
								goto l124
							l125:
								position, tokenIndex, depth = position124, tokenIndex124, depth124
								if buffer[position] != rune('E') {
									goto l100
								}
								position++
							}
						l124:
							{
								position126, tokenIndex126, depth126 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l127
								}
								position++
								goto l126
							l127:
								position, tokenIndex, depth = position126, tokenIndex126, depth126
								if buffer[position] != rune('D') {
									goto l100
								}
								position++
							}
						l126:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l129
								}
								position++
								goto l128
							l129:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
								if buffer[position] != rune('U') {
									goto l100
								}
								position++
							}
						l128:
							{
								position130, tokenIndex130, depth130 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l131
								}
								position++
								goto l130
							l131:
								position, tokenIndex, depth = position130, tokenIndex130, depth130
								if buffer[position] != rune('C') {
									goto l100
								}
								position++
							}
						l130:
							{
								position132, tokenIndex132, depth132 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l133
								}
								position++
								goto l132
							l133:
								position, tokenIndex, depth = position132, tokenIndex132, depth132
								if buffer[position] != rune('E') {
									goto l100
								}
								position++
							}
						l132:
							{
								position134, tokenIndex134, depth134 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l135
								}
								position++
								goto l134
							l135:
								position, tokenIndex, depth = position134, tokenIndex134, depth134
								if buffer[position] != rune('D') {
									goto l100
								}
								position++
							}
						l134:
							if !rules[rulews]() {
								goto l100
							}
							depth--
							add(ruleREDUCED, position121)
						}
					}
				l102:
					goto l101
				l100:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
				}
			l101:
				{
					position136, tokenIndex136, depth136 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l137
					}
					goto l136
				l137:
					position, tokenIndex, depth = position136, tokenIndex136, depth136
					{
						position140 := position
						depth++
						if !rules[rulevar]() {
							goto l85
						}
						depth--
						add(ruleprojectionElem, position140)
					}
				l138:
					{
						position139, tokenIndex139, depth139 := position, tokenIndex, depth
						{
							position141 := position
							depth++
							if !rules[rulevar]() {
								goto l139
							}
							depth--
							add(ruleprojectionElem, position141)
						}
						goto l138
					l139:
						position, tokenIndex, depth = position139, tokenIndex139, depth139
					}
				}
			l136:
				depth--
				add(ruleselect, position86)
			}
			return true
		l85:
			position, tokenIndex, depth = position85, tokenIndex85, depth85
			return false
		},
		/* 7 subSelect <- <(select whereClause)> */
		func() bool {
			position142, tokenIndex142, depth142 := position, tokenIndex, depth
			{
				position143 := position
				depth++
				if !rules[ruleselect]() {
					goto l142
				}
				if !rules[rulewhereClause]() {
					goto l142
				}
				depth--
				add(rulesubSelect, position143)
			}
			return true
		l142:
			position, tokenIndex, depth = position142, tokenIndex142, depth142
			return false
		},
		/* 8 projectionElem <- <var> */
		nil,
		/* 9 datasetClause <- <(FROM NAMED? iri)> */
		nil,
		/* 10 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position146, tokenIndex146, depth146 := position, tokenIndex, depth
			{
				position147 := position
				depth++
				{
					position148, tokenIndex148, depth148 := position, tokenIndex, depth
					{
						position150 := position
						depth++
						{
							position151, tokenIndex151, depth151 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l152
							}
							position++
							goto l151
						l152:
							position, tokenIndex, depth = position151, tokenIndex151, depth151
							if buffer[position] != rune('W') {
								goto l148
							}
							position++
						}
					l151:
						{
							position153, tokenIndex153, depth153 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l154
							}
							position++
							goto l153
						l154:
							position, tokenIndex, depth = position153, tokenIndex153, depth153
							if buffer[position] != rune('H') {
								goto l148
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
								goto l148
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
								goto l148
							}
							position++
						}
					l157:
						{
							position159, tokenIndex159, depth159 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l160
							}
							position++
							goto l159
						l160:
							position, tokenIndex, depth = position159, tokenIndex159, depth159
							if buffer[position] != rune('E') {
								goto l148
							}
							position++
						}
					l159:
						if !rules[rulews]() {
							goto l148
						}
						depth--
						add(ruleWHERE, position150)
					}
					goto l149
				l148:
					position, tokenIndex, depth = position148, tokenIndex148, depth148
				}
			l149:
				if !rules[rulegroupGraphPattern]() {
					goto l146
				}
				depth--
				add(rulewhereClause, position147)
			}
			return true
		l146:
			position, tokenIndex, depth = position146, tokenIndex146, depth146
			return false
		},
		/* 11 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position161, tokenIndex161, depth161 := position, tokenIndex, depth
			{
				position162 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l161
				}
				{
					position163, tokenIndex163, depth163 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l164
					}
					goto l163
				l164:
					position, tokenIndex, depth = position163, tokenIndex163, depth163
					if !rules[rulegraphPattern]() {
						goto l161
					}
				}
			l163:
				if !rules[ruleRBRACE]() {
					goto l161
				}
				depth--
				add(rulegroupGraphPattern, position162)
			}
			return true
		l161:
			position, tokenIndex, depth = position161, tokenIndex161, depth161
			return false
		},
		/* 12 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position166 := position
				depth++
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					{
						position169 := position
						depth++
						{
							position170 := position
							depth++
							if !rules[ruletriplesSameSubjectPath]() {
								goto l167
							}
						l171:
							{
								position172, tokenIndex172, depth172 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l172
								}
								if !rules[ruletriplesSameSubjectPath]() {
									goto l172
								}
								goto l171
							l172:
								position, tokenIndex, depth = position172, tokenIndex172, depth172
							}
							{
								position173, tokenIndex173, depth173 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l173
								}
								goto l174
							l173:
								position, tokenIndex, depth = position173, tokenIndex173, depth173
							}
						l174:
							depth--
							add(ruletriplesBlock, position170)
						}
						depth--
						add(rulebasicGraphPattern, position169)
					}
					goto l168
				l167:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
				}
			l168:
				{
					position175, tokenIndex175, depth175 := position, tokenIndex, depth
					{
						position177 := position
						depth++
						{
							position178, tokenIndex178, depth178 := position, tokenIndex, depth
							{
								position180 := position
								depth++
								{
									position181 := position
									depth++
									{
										position182, tokenIndex182, depth182 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l183
										}
										position++
										goto l182
									l183:
										position, tokenIndex, depth = position182, tokenIndex182, depth182
										if buffer[position] != rune('O') {
											goto l179
										}
										position++
									}
								l182:
									{
										position184, tokenIndex184, depth184 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l185
										}
										position++
										goto l184
									l185:
										position, tokenIndex, depth = position184, tokenIndex184, depth184
										if buffer[position] != rune('P') {
											goto l179
										}
										position++
									}
								l184:
									{
										position186, tokenIndex186, depth186 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l187
										}
										position++
										goto l186
									l187:
										position, tokenIndex, depth = position186, tokenIndex186, depth186
										if buffer[position] != rune('T') {
											goto l179
										}
										position++
									}
								l186:
									{
										position188, tokenIndex188, depth188 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l189
										}
										position++
										goto l188
									l189:
										position, tokenIndex, depth = position188, tokenIndex188, depth188
										if buffer[position] != rune('I') {
											goto l179
										}
										position++
									}
								l188:
									{
										position190, tokenIndex190, depth190 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l191
										}
										position++
										goto l190
									l191:
										position, tokenIndex, depth = position190, tokenIndex190, depth190
										if buffer[position] != rune('O') {
											goto l179
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
											goto l179
										}
										position++
									}
								l192:
									{
										position194, tokenIndex194, depth194 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l195
										}
										position++
										goto l194
									l195:
										position, tokenIndex, depth = position194, tokenIndex194, depth194
										if buffer[position] != rune('A') {
											goto l179
										}
										position++
									}
								l194:
									{
										position196, tokenIndex196, depth196 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l197
										}
										position++
										goto l196
									l197:
										position, tokenIndex, depth = position196, tokenIndex196, depth196
										if buffer[position] != rune('L') {
											goto l179
										}
										position++
									}
								l196:
									if !rules[rulews]() {
										goto l179
									}
									depth--
									add(ruleOPTIONAL, position181)
								}
								if !rules[ruleLBRACE]() {
									goto l179
								}
								{
									position198, tokenIndex198, depth198 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l199
									}
									goto l198
								l199:
									position, tokenIndex, depth = position198, tokenIndex198, depth198
									if !rules[rulegraphPattern]() {
										goto l179
									}
								}
							l198:
								if !rules[ruleRBRACE]() {
									goto l179
								}
								depth--
								add(ruleoptionalGraphPattern, position180)
							}
							goto l178
						l179:
							position, tokenIndex, depth = position178, tokenIndex178, depth178
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l175
							}
						}
					l178:
						depth--
						add(rulegraphPatternNotTriples, position177)
					}
					{
						position200, tokenIndex200, depth200 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l200
						}
						goto l201
					l200:
						position, tokenIndex, depth = position200, tokenIndex200, depth200
					}
				l201:
					if !rules[rulegraphPattern]() {
						goto l175
					}
					goto l176
				l175:
					position, tokenIndex, depth = position175, tokenIndex175, depth175
				}
			l176:
				depth--
				add(rulegraphPattern, position166)
			}
			return true
		},
		/* 13 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 14 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 15 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position204, tokenIndex204, depth204 := position, tokenIndex, depth
			{
				position205 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l204
				}
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					{
						position208 := position
						depth++
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('U') {
								goto l206
							}
							position++
						}
					l209:
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l212
							}
							position++
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if buffer[position] != rune('N') {
								goto l206
							}
							position++
						}
					l211:
						{
							position213, tokenIndex213, depth213 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l214
							}
							position++
							goto l213
						l214:
							position, tokenIndex, depth = position213, tokenIndex213, depth213
							if buffer[position] != rune('I') {
								goto l206
							}
							position++
						}
					l213:
						{
							position215, tokenIndex215, depth215 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l216
							}
							position++
							goto l215
						l216:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune('O') {
								goto l206
							}
							position++
						}
					l215:
						{
							position217, tokenIndex217, depth217 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l218
							}
							position++
							goto l217
						l218:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if buffer[position] != rune('N') {
								goto l206
							}
							position++
						}
					l217:
						if !rules[rulews]() {
							goto l206
						}
						depth--
						add(ruleUNION, position208)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l206
					}
					goto l207
				l206:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
				}
			l207:
				depth--
				add(rulegroupOrUnionGraphPattern, position205)
			}
			return true
		l204:
			position, tokenIndex, depth = position204, tokenIndex204, depth204
			return false
		},
		/* 16 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 17 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		nil,
		/* 18 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position221, tokenIndex221, depth221 := position, tokenIndex, depth
			{
				position222 := position
				depth++
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					{
						position225 := position
						depth++
						{
							position226, tokenIndex226, depth226 := position, tokenIndex, depth
							{
								position228 := position
								depth++
								if !rules[rulevar]() {
									goto l227
								}
								depth--
								add(rulePegText, position228)
							}
							{
								add(ruleAction0, position)
							}
							goto l226
						l227:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							{
								position231 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l230
								}
								depth--
								add(rulePegText, position231)
							}
							{
								add(ruleAction1, position)
							}
							goto l226
						l230:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							if !rules[rulepof]() {
								goto l224
							}
							{
								add(ruleAction2, position)
							}
						}
					l226:
						depth--
						add(rulevarOrTerm, position225)
					}
					if !rules[rulepropertyListPath]() {
						goto l224
					}
					goto l223
				l224:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
					{
						position234 := position
						depth++
						{
							position235, tokenIndex235, depth235 := position, tokenIndex, depth
							{
								position237 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l236
								}
								if !rules[rulegraphNodePath]() {
									goto l236
								}
							l238:
								{
									position239, tokenIndex239, depth239 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l239
									}
									goto l238
								l239:
									position, tokenIndex, depth = position239, tokenIndex239, depth239
								}
								if !rules[ruleRPAREN]() {
									goto l236
								}
								depth--
								add(rulecollectionPath, position237)
							}
							goto l235
						l236:
							position, tokenIndex, depth = position235, tokenIndex235, depth235
							{
								position240 := position
								depth++
								{
									position241 := position
									depth++
									if buffer[position] != rune('[') {
										goto l221
									}
									position++
									if !rules[rulews]() {
										goto l221
									}
									depth--
									add(ruleLBRACK, position241)
								}
								if !rules[rulepropertyListPath]() {
									goto l221
								}
								{
									position242 := position
									depth++
									if buffer[position] != rune(']') {
										goto l221
									}
									position++
									if !rules[rulews]() {
										goto l221
									}
									depth--
									add(ruleRBRACK, position242)
								}
								depth--
								add(ruleblankNodePropertyListPath, position240)
							}
						}
					l235:
						depth--
						add(ruletriplesNodePath, position234)
					}
					if !rules[rulepropertyListPath]() {
						goto l221
					}
				}
			l223:
				depth--
				add(ruletriplesSameSubjectPath, position222)
			}
			return true
		l221:
			position, tokenIndex, depth = position221, tokenIndex221, depth221
			return false
		},
		/* 19 varOrTerm <- <((<var> Action0) / (<graphTerm> Action1) / (pof Action2))> */
		nil,
		/* 20 graphTerm <- <((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('<') iri) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral))> */
		func() bool {
			position244, tokenIndex244, depth244 := position, tokenIndex, depth
			{
				position245 := position
				depth++
				{
					switch buffer[position] {
					case '(':
						{
							position247 := position
							depth++
							if buffer[position] != rune('(') {
								goto l244
							}
							position++
							if !rules[rulews]() {
								goto l244
							}
							if buffer[position] != rune(')') {
								goto l244
							}
							position++
							if !rules[rulews]() {
								goto l244
							}
							depth--
							add(rulenil, position247)
						}
						break
					case '[', '_':
						{
							position248 := position
							depth++
							{
								position249, tokenIndex249, depth249 := position, tokenIndex, depth
								{
									position251 := position
									depth++
									if buffer[position] != rune('_') {
										goto l250
									}
									position++
									if buffer[position] != rune(':') {
										goto l250
									}
									position++
									{
										switch buffer[position] {
										case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l250
											}
											position++
											break
										case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l250
											}
											position++
											break
										default:
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l250
											}
											position++
											break
										}
									}

									{
										position253, tokenIndex253, depth253 := position, tokenIndex, depth
										{
											position255, tokenIndex255, depth255 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l256
											}
											position++
											goto l255
										l256:
											position, tokenIndex, depth = position255, tokenIndex255, depth255
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l257
											}
											position++
											goto l255
										l257:
											position, tokenIndex, depth = position255, tokenIndex255, depth255
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l258
											}
											position++
											goto l255
										l258:
											position, tokenIndex, depth = position255, tokenIndex255, depth255
											if c := buffer[position]; c < rune('.') || c > rune('_') {
												goto l253
											}
											position++
										}
									l255:
										goto l254
									l253:
										position, tokenIndex, depth = position253, tokenIndex253, depth253
									}
								l254:
									if !rules[rulews]() {
										goto l250
									}
									depth--
									add(ruleblankNodeLabel, position251)
								}
								goto l249
							l250:
								position, tokenIndex, depth = position249, tokenIndex249, depth249
								{
									position259 := position
									depth++
									if buffer[position] != rune('[') {
										goto l244
									}
									position++
									if !rules[rulews]() {
										goto l244
									}
									if buffer[position] != rune(']') {
										goto l244
									}
									position++
									if !rules[rulews]() {
										goto l244
									}
									depth--
									add(ruleanon, position259)
								}
							}
						l249:
							depth--
							add(ruleblankNode, position248)
						}
						break
					case 'F', 'T', 'f', 't':
						{
							position260 := position
							depth++
							{
								position261, tokenIndex261, depth261 := position, tokenIndex, depth
								{
									position263 := position
									depth++
									{
										position264, tokenIndex264, depth264 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l265
										}
										position++
										goto l264
									l265:
										position, tokenIndex, depth = position264, tokenIndex264, depth264
										if buffer[position] != rune('T') {
											goto l262
										}
										position++
									}
								l264:
									{
										position266, tokenIndex266, depth266 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l267
										}
										position++
										goto l266
									l267:
										position, tokenIndex, depth = position266, tokenIndex266, depth266
										if buffer[position] != rune('R') {
											goto l262
										}
										position++
									}
								l266:
									{
										position268, tokenIndex268, depth268 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l269
										}
										position++
										goto l268
									l269:
										position, tokenIndex, depth = position268, tokenIndex268, depth268
										if buffer[position] != rune('U') {
											goto l262
										}
										position++
									}
								l268:
									{
										position270, tokenIndex270, depth270 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l271
										}
										position++
										goto l270
									l271:
										position, tokenIndex, depth = position270, tokenIndex270, depth270
										if buffer[position] != rune('E') {
											goto l262
										}
										position++
									}
								l270:
									if !rules[rulews]() {
										goto l262
									}
									depth--
									add(ruleTRUE, position263)
								}
								goto l261
							l262:
								position, tokenIndex, depth = position261, tokenIndex261, depth261
								{
									position272 := position
									depth++
									{
										position273, tokenIndex273, depth273 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l274
										}
										position++
										goto l273
									l274:
										position, tokenIndex, depth = position273, tokenIndex273, depth273
										if buffer[position] != rune('F') {
											goto l244
										}
										position++
									}
								l273:
									{
										position275, tokenIndex275, depth275 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l276
										}
										position++
										goto l275
									l276:
										position, tokenIndex, depth = position275, tokenIndex275, depth275
										if buffer[position] != rune('A') {
											goto l244
										}
										position++
									}
								l275:
									{
										position277, tokenIndex277, depth277 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l278
										}
										position++
										goto l277
									l278:
										position, tokenIndex, depth = position277, tokenIndex277, depth277
										if buffer[position] != rune('L') {
											goto l244
										}
										position++
									}
								l277:
									{
										position279, tokenIndex279, depth279 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l280
										}
										position++
										goto l279
									l280:
										position, tokenIndex, depth = position279, tokenIndex279, depth279
										if buffer[position] != rune('S') {
											goto l244
										}
										position++
									}
								l279:
									{
										position281, tokenIndex281, depth281 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l282
										}
										position++
										goto l281
									l282:
										position, tokenIndex, depth = position281, tokenIndex281, depth281
										if buffer[position] != rune('E') {
											goto l244
										}
										position++
									}
								l281:
									if !rules[rulews]() {
										goto l244
									}
									depth--
									add(ruleFALSE, position272)
								}
							}
						l261:
							depth--
							add(rulebooleanLiteral, position260)
						}
						break
					case '"':
						{
							position283 := position
							depth++
							{
								position284 := position
								depth++
								if buffer[position] != rune('"') {
									goto l244
								}
								position++
							l285:
								{
									position286, tokenIndex286, depth286 := position, tokenIndex, depth
									{
										position287, tokenIndex287, depth287 := position, tokenIndex, depth
										{
											position288, tokenIndex288, depth288 := position, tokenIndex, depth
											if buffer[position] != rune('\'') {
												goto l289
											}
											position++
											goto l288
										l289:
											position, tokenIndex, depth = position288, tokenIndex288, depth288
											if buffer[position] != rune('"') {
												goto l290
											}
											position++
											goto l288
										l290:
											position, tokenIndex, depth = position288, tokenIndex288, depth288
											if buffer[position] != rune('\'') {
												goto l287
											}
											position++
										}
									l288:
										goto l286
									l287:
										position, tokenIndex, depth = position287, tokenIndex287, depth287
									}
									if !matchDot() {
										goto l286
									}
									goto l285
								l286:
									position, tokenIndex, depth = position286, tokenIndex286, depth286
								}
								if buffer[position] != rune('"') {
									goto l244
								}
								position++
								depth--
								add(rulestring, position284)
							}
							{
								position291, tokenIndex291, depth291 := position, tokenIndex, depth
								{
									position293, tokenIndex293, depth293 := position, tokenIndex, depth
									if buffer[position] != rune('@') {
										goto l294
									}
									position++
									{
										position297, tokenIndex297, depth297 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l298
										}
										position++
										goto l297
									l298:
										position, tokenIndex, depth = position297, tokenIndex297, depth297
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l294
										}
										position++
									}
								l297:
								l295:
									{
										position296, tokenIndex296, depth296 := position, tokenIndex, depth
										{
											position299, tokenIndex299, depth299 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l300
											}
											position++
											goto l299
										l300:
											position, tokenIndex, depth = position299, tokenIndex299, depth299
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l296
											}
											position++
										}
									l299:
										goto l295
									l296:
										position, tokenIndex, depth = position296, tokenIndex296, depth296
									}
								l301:
									{
										position302, tokenIndex302, depth302 := position, tokenIndex, depth
										if buffer[position] != rune('-') {
											goto l302
										}
										position++
										{
											switch buffer[position] {
											case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l302
												}
												position++
												break
											case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l302
												}
												position++
												break
											default:
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l302
												}
												position++
												break
											}
										}

									l303:
										{
											position304, tokenIndex304, depth304 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l304
													}
													position++
													break
												case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l304
													}
													position++
													break
												default:
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l304
													}
													position++
													break
												}
											}

											goto l303
										l304:
											position, tokenIndex, depth = position304, tokenIndex304, depth304
										}
										goto l301
									l302:
										position, tokenIndex, depth = position302, tokenIndex302, depth302
									}
									goto l293
								l294:
									position, tokenIndex, depth = position293, tokenIndex293, depth293
									if buffer[position] != rune('^') {
										goto l291
									}
									position++
									if buffer[position] != rune('^') {
										goto l291
									}
									position++
									if !rules[ruleiri]() {
										goto l291
									}
								}
							l293:
								goto l292
							l291:
								position, tokenIndex, depth = position291, tokenIndex291, depth291
							}
						l292:
							if !rules[rulews]() {
								goto l244
							}
							depth--
							add(ruleliteral, position283)
						}
						break
					case '<':
						if !rules[ruleiri]() {
							goto l244
						}
						break
					default:
						{
							position307 := position
							depth++
							{
								position308, tokenIndex308, depth308 := position, tokenIndex, depth
								{
									position310, tokenIndex310, depth310 := position, tokenIndex, depth
									if buffer[position] != rune('+') {
										goto l311
									}
									position++
									goto l310
								l311:
									position, tokenIndex, depth = position310, tokenIndex310, depth310
									if buffer[position] != rune('-') {
										goto l308
									}
									position++
								}
							l310:
								goto l309
							l308:
								position, tokenIndex, depth = position308, tokenIndex308, depth308
							}
						l309:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l244
							}
							position++
						l312:
							{
								position313, tokenIndex313, depth313 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l313
								}
								position++
								goto l312
							l313:
								position, tokenIndex, depth = position313, tokenIndex313, depth313
							}
							{
								position314, tokenIndex314, depth314 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l314
								}
								position++
							l316:
								{
									position317, tokenIndex317, depth317 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l317
									}
									position++
									goto l316
								l317:
									position, tokenIndex, depth = position317, tokenIndex317, depth317
								}
								goto l315
							l314:
								position, tokenIndex, depth = position314, tokenIndex314, depth314
							}
						l315:
							if !rules[rulews]() {
								goto l244
							}
							depth--
							add(rulenumericLiteral, position307)
						}
						break
					}
				}

				depth--
				add(rulegraphTerm, position245)
			}
			return true
		l244:
			position, tokenIndex, depth = position244, tokenIndex244, depth244
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
			position321, tokenIndex321, depth321 := position, tokenIndex, depth
			{
				position322 := position
				depth++
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l324
					}
					{
						add(ruleAction3, position)
					}
					goto l323
				l324:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					{
						position327 := position
						depth++
						if !rules[rulevar]() {
							goto l326
						}
						depth--
						add(rulePegText, position327)
					}
					{
						add(ruleAction4, position)
					}
					goto l323
				l326:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
					{
						position329 := position
						depth++
						if !rules[rulepath]() {
							goto l321
						}
						depth--
						add(ruleverbPath, position329)
					}
				}
			l323:
				if !rules[ruleobjectListPath]() {
					goto l321
				}
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					{
						position332 := position
						depth++
						if buffer[position] != rune(';') {
							goto l330
						}
						position++
						if !rules[rulews]() {
							goto l330
						}
						depth--
						add(ruleSEMICOLON, position332)
					}
					if !rules[rulepropertyListPath]() {
						goto l330
					}
					goto l331
				l330:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
				}
			l331:
				depth--
				add(rulepropertyListPath, position322)
			}
			return true
		l321:
			position, tokenIndex, depth = position321, tokenIndex321, depth321
			return false
		},
		/* 25 verbPath <- <path> */
		nil,
		/* 26 path <- <pathAlternative> */
		func() bool {
			position334, tokenIndex334, depth334 := position, tokenIndex, depth
			{
				position335 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l334
				}
				depth--
				add(rulepath, position335)
			}
			return true
		l334:
			position, tokenIndex, depth = position334, tokenIndex334, depth334
			return false
		},
		/* 27 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position336, tokenIndex336, depth336 := position, tokenIndex, depth
			{
				position337 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l336
				}
			l338:
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l339
					}
					if !rules[rulepathAlternative]() {
						goto l339
					}
					goto l338
				l339:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
				}
				depth--
				add(rulepathAlternative, position337)
			}
			return true
		l336:
			position, tokenIndex, depth = position336, tokenIndex336, depth336
			return false
		},
		/* 28 pathSequence <- <(<pathElt> Action5 (SLASH pathSequence)*)> */
		func() bool {
			position340, tokenIndex340, depth340 := position, tokenIndex, depth
			{
				position341 := position
				depth++
				{
					position342 := position
					depth++
					{
						position343 := position
						depth++
						{
							position344, tokenIndex344, depth344 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l344
							}
							goto l345
						l344:
							position, tokenIndex, depth = position344, tokenIndex344, depth344
						}
					l345:
						{
							position346 := position
							depth++
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l340
									}
									if !rules[rulepath]() {
										goto l340
									}
									if !rules[ruleRPAREN]() {
										goto l340
									}
									break
								case '!':
									{
										position348 := position
										depth++
										if buffer[position] != rune('!') {
											goto l340
										}
										position++
										if !rules[rulews]() {
											goto l340
										}
										depth--
										add(ruleNOT, position348)
									}
									{
										position349 := position
										depth++
										{
											position350, tokenIndex350, depth350 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l351
											}
											goto l350
										l351:
											position, tokenIndex, depth = position350, tokenIndex350, depth350
											if !rules[ruleLPAREN]() {
												goto l340
											}
											{
												position352, tokenIndex352, depth352 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l352
												}
											l354:
												{
													position355, tokenIndex355, depth355 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l355
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l355
													}
													goto l354
												l355:
													position, tokenIndex, depth = position355, tokenIndex355, depth355
												}
												goto l353
											l352:
												position, tokenIndex, depth = position352, tokenIndex352, depth352
											}
										l353:
											if !rules[ruleRPAREN]() {
												goto l340
											}
										}
									l350:
										depth--
										add(rulepathNegatedPropertySet, position349)
									}
									break
								case 'a':
									if !rules[ruleISA]() {
										goto l340
									}
									break
								default:
									if !rules[ruleiri]() {
										goto l340
									}
									break
								}
							}

							depth--
							add(rulepathPrimary, position346)
						}
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							{
								position357, tokenIndex357, depth357 := position, tokenIndex, depth
								{
									position359 := position
									depth++
									{
										switch buffer[position] {
										case '+':
											{
												position361 := position
												depth++
												if buffer[position] != rune('+') {
													goto l357
												}
												position++
												if !rules[rulews]() {
													goto l357
												}
												depth--
												add(rulePLUS, position361)
											}
											break
										case '?':
											{
												position362 := position
												depth++
												if buffer[position] != rune('?') {
													goto l357
												}
												position++
												if !rules[rulews]() {
													goto l357
												}
												depth--
												add(ruleQUESTION, position362)
											}
											break
										default:
											if !rules[ruleSTAR]() {
												goto l357
											}
											break
										}
									}

									depth--
									add(rulepathMod, position359)
								}
								goto l358
							l357:
								position, tokenIndex, depth = position357, tokenIndex357, depth357
							}
						l358:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
						}
						depth--
						add(rulepathElt, position343)
					}
					depth--
					add(rulePegText, position342)
				}
				{
					add(ruleAction5, position)
				}
			l364:
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					{
						position366 := position
						depth++
						if buffer[position] != rune('/') {
							goto l365
						}
						position++
						if !rules[rulews]() {
							goto l365
						}
						depth--
						add(ruleSLASH, position366)
					}
					if !rules[rulepathSequence]() {
						goto l365
					}
					goto l364
				l365:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
				}
				depth--
				add(rulepathSequence, position341)
			}
			return true
		l340:
			position, tokenIndex, depth = position340, tokenIndex340, depth340
			return false
		},
		/* 29 pathElt <- <(INVERSE? pathPrimary &pathMod?)> */
		nil,
		/* 30 pathPrimary <- <((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA) | (&('<') iri))> */
		nil,
		/* 31 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 32 pathOneInPropertySet <- <((&('^') (INVERSE (iri / ISA))) | (&('a') ISA) | (&('<') iri))> */
		func() bool {
			position370, tokenIndex370, depth370 := position, tokenIndex, depth
			{
				position371 := position
				depth++
				{
					switch buffer[position] {
					case '^':
						if !rules[ruleINVERSE]() {
							goto l370
						}
						{
							position373, tokenIndex373, depth373 := position, tokenIndex, depth
							if !rules[ruleiri]() {
								goto l374
							}
							goto l373
						l374:
							position, tokenIndex, depth = position373, tokenIndex373, depth373
							if !rules[ruleISA]() {
								goto l370
							}
						}
					l373:
						break
					case 'a':
						if !rules[ruleISA]() {
							goto l370
						}
						break
					default:
						if !rules[ruleiri]() {
							goto l370
						}
						break
					}
				}

				depth--
				add(rulepathOneInPropertySet, position371)
			}
			return true
		l370:
			position, tokenIndex, depth = position370, tokenIndex370, depth370
			return false
		},
		/* 33 pathMod <- <((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR))> */
		nil,
		/* 34 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			{
				position377 := position
				depth++
				{
					position378 := position
					depth++
					{
						position379, tokenIndex379, depth379 := position, tokenIndex, depth
						{
							position381 := position
							depth++
							if !rules[rulegraphNodePath]() {
								goto l380
							}
							depth--
							add(rulePegText, position381)
						}
						{
							add(ruleAction6, position)
						}
						goto l379
					l380:
						position, tokenIndex, depth = position379, tokenIndex379, depth379
						if !rules[rulepof]() {
							goto l383
						}
						{
							add(ruleAction7, position)
						}
						goto l379
					l383:
						position, tokenIndex, depth = position379, tokenIndex379, depth379
						{
							add(ruleAction8, position)
						}
					}
				l379:
					depth--
					add(ruleobjectPath, position378)
				}
			l386:
				{
					position387, tokenIndex387, depth387 := position, tokenIndex, depth
					{
						position388 := position
						depth++
						if buffer[position] != rune(',') {
							goto l387
						}
						position++
						if !rules[rulews]() {
							goto l387
						}
						depth--
						add(ruleCOMMA, position388)
					}
					if !rules[ruleobjectListPath]() {
						goto l387
					}
					goto l386
				l387:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
				}
				depth--
				add(ruleobjectListPath, position377)
			}
			return true
		},
		/* 35 objectPath <- <((<graphNodePath> Action6) / (pof Action7) / Action8)> */
		nil,
		/* 36 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l393
					}
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					if !rules[rulegraphTerm]() {
						goto l390
					}
				}
			l392:
				depth--
				add(rulegraphNodePath, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 37 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 38 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 39 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position396, tokenIndex396, depth396 := position, tokenIndex, depth
			{
				position397 := position
				depth++
				{
					position398 := position
					depth++
					{
						position399, tokenIndex399, depth399 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l400
						}
						position++
						goto l399
					l400:
						position, tokenIndex, depth = position399, tokenIndex399, depth399
						if buffer[position] != rune('L') {
							goto l396
						}
						position++
					}
				l399:
					{
						position401, tokenIndex401, depth401 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l402
						}
						position++
						goto l401
					l402:
						position, tokenIndex, depth = position401, tokenIndex401, depth401
						if buffer[position] != rune('I') {
							goto l396
						}
						position++
					}
				l401:
					{
						position403, tokenIndex403, depth403 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l404
						}
						position++
						goto l403
					l404:
						position, tokenIndex, depth = position403, tokenIndex403, depth403
						if buffer[position] != rune('M') {
							goto l396
						}
						position++
					}
				l403:
					{
						position405, tokenIndex405, depth405 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l406
						}
						position++
						goto l405
					l406:
						position, tokenIndex, depth = position405, tokenIndex405, depth405
						if buffer[position] != rune('I') {
							goto l396
						}
						position++
					}
				l405:
					{
						position407, tokenIndex407, depth407 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l408
						}
						position++
						goto l407
					l408:
						position, tokenIndex, depth = position407, tokenIndex407, depth407
						if buffer[position] != rune('T') {
							goto l396
						}
						position++
					}
				l407:
					if !rules[rulews]() {
						goto l396
					}
					depth--
					add(ruleLIMIT, position398)
				}
				if !rules[ruleINTEGER]() {
					goto l396
				}
				depth--
				add(rulelimit, position397)
			}
			return true
		l396:
			position, tokenIndex, depth = position396, tokenIndex396, depth396
			return false
		},
		/* 40 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position409, tokenIndex409, depth409 := position, tokenIndex, depth
			{
				position410 := position
				depth++
				{
					position411 := position
					depth++
					{
						position412, tokenIndex412, depth412 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l413
						}
						position++
						goto l412
					l413:
						position, tokenIndex, depth = position412, tokenIndex412, depth412
						if buffer[position] != rune('O') {
							goto l409
						}
						position++
					}
				l412:
					{
						position414, tokenIndex414, depth414 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l415
						}
						position++
						goto l414
					l415:
						position, tokenIndex, depth = position414, tokenIndex414, depth414
						if buffer[position] != rune('F') {
							goto l409
						}
						position++
					}
				l414:
					{
						position416, tokenIndex416, depth416 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l417
						}
						position++
						goto l416
					l417:
						position, tokenIndex, depth = position416, tokenIndex416, depth416
						if buffer[position] != rune('F') {
							goto l409
						}
						position++
					}
				l416:
					{
						position418, tokenIndex418, depth418 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l419
						}
						position++
						goto l418
					l419:
						position, tokenIndex, depth = position418, tokenIndex418, depth418
						if buffer[position] != rune('S') {
							goto l409
						}
						position++
					}
				l418:
					{
						position420, tokenIndex420, depth420 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l421
						}
						position++
						goto l420
					l421:
						position, tokenIndex, depth = position420, tokenIndex420, depth420
						if buffer[position] != rune('E') {
							goto l409
						}
						position++
					}
				l420:
					{
						position422, tokenIndex422, depth422 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l423
						}
						position++
						goto l422
					l423:
						position, tokenIndex, depth = position422, tokenIndex422, depth422
						if buffer[position] != rune('T') {
							goto l409
						}
						position++
					}
				l422:
					if !rules[rulews]() {
						goto l409
					}
					depth--
					add(ruleOFFSET, position411)
				}
				if !rules[ruleINTEGER]() {
					goto l409
				}
				depth--
				add(ruleoffset, position410)
			}
			return true
		l409:
			position, tokenIndex, depth = position409, tokenIndex409, depth409
			return false
		},
		/* 41 pof <- <('<' ((&('\f') '\f') | (&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))+)> */
		func() bool {
			position424, tokenIndex424, depth424 := position, tokenIndex, depth
			{
				position425 := position
				depth++
				if buffer[position] != rune('<') {
					goto l424
				}
				position++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							goto l424
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l424
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l424
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l424
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l424
						}
						position++
						break
					}
				}

			l426:
				{
					position427, tokenIndex427, depth427 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l427
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l427
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l427
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l427
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l427
							}
							position++
							break
						}
					}

					goto l426
				l427:
					position, tokenIndex, depth = position427, tokenIndex427, depth427
				}
				depth--
				add(rulepof, position425)
			}
			return true
		l424:
			position, tokenIndex, depth = position424, tokenIndex424, depth424
			return false
		},
		/* 42 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position430, tokenIndex430, depth430 := position, tokenIndex, depth
			{
				position431 := position
				depth++
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l433
					}
					position++
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if buffer[position] != rune('$') {
						goto l430
					}
					position++
				}
			l432:
				{
					position434 := position
					depth++
					{
						position435, tokenIndex435, depth435 := position, tokenIndex, depth
						if !rules[rulePN_CHARS_U]() {
							goto l436
						}
						goto l435
					l436:
						position, tokenIndex, depth = position435, tokenIndex435, depth435
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l430
						}
						position++
					}
				l435:
				l437:
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						{
							position439 := position
							depth++
							{
								position440, tokenIndex440, depth440 := position, tokenIndex, depth
								if !rules[rulePN_CHARS_U]() {
									goto l441
								}
								goto l440
							l441:
								position, tokenIndex, depth = position440, tokenIndex440, depth440
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l438
								}
								position++
							}
						l440:
							depth--
							add(ruleVAR_CHAR, position439)
						}
						goto l437
					l438:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
					}
					depth--
					add(ruleVARNAME, position434)
				}
				if !rules[rulews]() {
					goto l430
				}
				depth--
				add(rulevar, position431)
			}
			return true
		l430:
			position, tokenIndex, depth = position430, tokenIndex430, depth430
			return false
		},
		/* 43 iri <- <('<' (!('\'' / '>' / '\'') .)* '>' ws)> */
		func() bool {
			position442, tokenIndex442, depth442 := position, tokenIndex, depth
			{
				position443 := position
				depth++
				if buffer[position] != rune('<') {
					goto l442
				}
				position++
			l444:
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					{
						position446, tokenIndex446, depth446 := position, tokenIndex, depth
						{
							position447, tokenIndex447, depth447 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l448
							}
							position++
							goto l447
						l448:
							position, tokenIndex, depth = position447, tokenIndex447, depth447
							if buffer[position] != rune('>') {
								goto l449
							}
							position++
							goto l447
						l449:
							position, tokenIndex, depth = position447, tokenIndex447, depth447
							if buffer[position] != rune('\'') {
								goto l446
							}
							position++
						}
					l447:
						goto l445
					l446:
						position, tokenIndex, depth = position446, tokenIndex446, depth446
					}
					if !matchDot() {
						goto l445
					}
					goto l444
				l445:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
				}
				if buffer[position] != rune('>') {
					goto l442
				}
				position++
				if !rules[rulews]() {
					goto l442
				}
				depth--
				add(ruleiri, position443)
			}
			return true
		l442:
			position, tokenIndex, depth = position442, tokenIndex442, depth442
			return false
		},
		/* 44 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iri))? ws)> */
		nil,
		/* 45 string <- <('"' (!('\'' / '"' / '\'') .)* '"')> */
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
		/* 52 VARNAME <- <((PN_CHARS_U / [0-9]) VAR_CHAR*)> */
		nil,
		/* 53 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		func() bool {
			position459, tokenIndex459, depth459 := position, tokenIndex, depth
			{
				position460 := position
				depth++
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					{
						position463 := position
						depth++
						{
							position464, tokenIndex464, depth464 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l465
							}
							position++
							goto l464
						l465:
							position, tokenIndex, depth = position464, tokenIndex464, depth464
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l462
							}
							position++
						}
					l464:
						depth--
						add(rulePN_CHARS_BASE, position463)
					}
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					if buffer[position] != rune('_') {
						goto l459
					}
					position++
				}
			l461:
				depth--
				add(rulePN_CHARS_U, position460)
			}
			return true
		l459:
			position, tokenIndex, depth = position459, tokenIndex459, depth459
			return false
		},
		/* 54 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 55 VAR_CHAR <- <(PN_CHARS_U / [0-9])> */
		nil,
		/* 56 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') ws)> */
		nil,
		/* 57 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') ws)> */
		nil,
		/* 58 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 59 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 60 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 61 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 62 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 63 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') ws)> */
		nil,
		/* 64 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 65 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') ws)> */
		nil,
		/* 66 LBRACE <- <('{' ws)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				if buffer[position] != rune('{') {
					goto l478
				}
				position++
				if !rules[rulews]() {
					goto l478
				}
				depth--
				add(ruleLBRACE, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 67 RBRACE <- <('}' ws)> */
		func() bool {
			position480, tokenIndex480, depth480 := position, tokenIndex, depth
			{
				position481 := position
				depth++
				if buffer[position] != rune('}') {
					goto l480
				}
				position++
				if !rules[rulews]() {
					goto l480
				}
				depth--
				add(ruleRBRACE, position481)
			}
			return true
		l480:
			position, tokenIndex, depth = position480, tokenIndex480, depth480
			return false
		},
		/* 68 LBRACK <- <('[' ws)> */
		nil,
		/* 69 RBRACK <- <(']' ws)> */
		nil,
		/* 70 SEMICOLON <- <(';' ws)> */
		nil,
		/* 71 COMMA <- <(',' ws)> */
		nil,
		/* 72 DOT <- <('.' ws)> */
		func() bool {
			position486, tokenIndex486, depth486 := position, tokenIndex, depth
			{
				position487 := position
				depth++
				if buffer[position] != rune('.') {
					goto l486
				}
				position++
				if !rules[rulews]() {
					goto l486
				}
				depth--
				add(ruleDOT, position487)
			}
			return true
		l486:
			position, tokenIndex, depth = position486, tokenIndex486, depth486
			return false
		},
		/* 73 COLON <- <(':' ws)> */
		nil,
		/* 74 PIPE <- <('|' ws)> */
		func() bool {
			position489, tokenIndex489, depth489 := position, tokenIndex, depth
			{
				position490 := position
				depth++
				if buffer[position] != rune('|') {
					goto l489
				}
				position++
				if !rules[rulews]() {
					goto l489
				}
				depth--
				add(rulePIPE, position490)
			}
			return true
		l489:
			position, tokenIndex, depth = position489, tokenIndex489, depth489
			return false
		},
		/* 75 SLASH <- <('/' ws)> */
		nil,
		/* 76 INVERSE <- <('^' ws)> */
		func() bool {
			position492, tokenIndex492, depth492 := position, tokenIndex, depth
			{
				position493 := position
				depth++
				if buffer[position] != rune('^') {
					goto l492
				}
				position++
				if !rules[rulews]() {
					goto l492
				}
				depth--
				add(ruleINVERSE, position493)
			}
			return true
		l492:
			position, tokenIndex, depth = position492, tokenIndex492, depth492
			return false
		},
		/* 77 LPAREN <- <('(' ws)> */
		func() bool {
			position494, tokenIndex494, depth494 := position, tokenIndex, depth
			{
				position495 := position
				depth++
				if buffer[position] != rune('(') {
					goto l494
				}
				position++
				if !rules[rulews]() {
					goto l494
				}
				depth--
				add(ruleLPAREN, position495)
			}
			return true
		l494:
			position, tokenIndex, depth = position494, tokenIndex494, depth494
			return false
		},
		/* 78 RPAREN <- <(')' ws)> */
		func() bool {
			position496, tokenIndex496, depth496 := position, tokenIndex, depth
			{
				position497 := position
				depth++
				if buffer[position] != rune(')') {
					goto l496
				}
				position++
				if !rules[rulews]() {
					goto l496
				}
				depth--
				add(ruleRPAREN, position497)
			}
			return true
		l496:
			position, tokenIndex, depth = position496, tokenIndex496, depth496
			return false
		},
		/* 79 ISA <- <('a' ws)> */
		func() bool {
			position498, tokenIndex498, depth498 := position, tokenIndex, depth
			{
				position499 := position
				depth++
				if buffer[position] != rune('a') {
					goto l498
				}
				position++
				if !rules[rulews]() {
					goto l498
				}
				depth--
				add(ruleISA, position499)
			}
			return true
		l498:
			position, tokenIndex, depth = position498, tokenIndex498, depth498
			return false
		},
		/* 80 NOT <- <('!' ws)> */
		nil,
		/* 81 STAR <- <('*' ws)> */
		func() bool {
			position501, tokenIndex501, depth501 := position, tokenIndex, depth
			{
				position502 := position
				depth++
				if buffer[position] != rune('*') {
					goto l501
				}
				position++
				if !rules[rulews]() {
					goto l501
				}
				depth--
				add(ruleSTAR, position502)
			}
			return true
		l501:
			position, tokenIndex, depth = position501, tokenIndex501, depth501
			return false
		},
		/* 82 QUESTION <- <('?' ws)> */
		nil,
		/* 83 PLUS <- <('+' ws)> */
		nil,
		/* 84 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') ws)> */
		nil,
		/* 85 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') ws)> */
		nil,
		/* 86 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') ws)> */
		nil,
		/* 87 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') ws)> */
		nil,
		/* 88 INTEGER <- <([0-9]+ ws)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l509
				}
				position++
			l511:
				{
					position512, tokenIndex512, depth512 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l512
					}
					position++
					goto l511
				l512:
					position, tokenIndex, depth = position512, tokenIndex512, depth512
				}
				if !rules[rulews]() {
					goto l509
				}
				depth--
				add(ruleINTEGER, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 89 ws <- <((&('\f') '\f') | (&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position514 := position
				depth++
			l515:
				{
					position516, tokenIndex516, depth516 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l516
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l516
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l516
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l516
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l516
							}
							position++
							break
						}
					}

					goto l515
				l516:
					position, tokenIndex, depth = position516, tokenIndex516, depth516
				}
				depth--
				add(rulews, position514)
			}
			return true
		},
		nil,
		/* 92 Action0 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 93 Action1 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 94 Action2 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 95 Action3 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 96 Action4 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 97 Action5 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 98 Action6 <- <{ p.setObject(buffer[begin:end]); p.addTriplePattern() }> */
		nil,
		/* 99 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 100 Action8 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
	}
	p.rules = rules
}
