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
	*Bgp

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
								if !rules[ruleiri]() {
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
					if !rules[ruleSTAR]() {
						goto l135
					}
					goto l134
				l135:
					position, tokenIndex, depth = position134, tokenIndex134, depth134
					{
						position138 := position
						depth++
						if !rules[rulevar]() {
							goto l83
						}
						depth--
						add(ruleprojectionElem, position138)
					}
				l136:
					{
						position137, tokenIndex137, depth137 := position, tokenIndex, depth
						{
							position139 := position
							depth++
							if !rules[rulevar]() {
								goto l137
							}
							depth--
							add(ruleprojectionElem, position139)
						}
						goto l136
					l137:
						position, tokenIndex, depth = position137, tokenIndex137, depth137
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
			position140, tokenIndex140, depth140 := position, tokenIndex, depth
			{
				position141 := position
				depth++
				if !rules[ruleselect]() {
					goto l140
				}
				if !rules[rulewhereClause]() {
					goto l140
				}
				depth--
				add(rulesubSelect, position141)
			}
			return true
		l140:
			position, tokenIndex, depth = position140, tokenIndex140, depth140
			return false
		},
		/* 8 projectionElem <- <var> */
		nil,
		/* 9 datasetClause <- <(FROM NAMED? iri)> */
		nil,
		/* 10 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position144, tokenIndex144, depth144 := position, tokenIndex, depth
			{
				position145 := position
				depth++
				{
					position146, tokenIndex146, depth146 := position, tokenIndex, depth
					{
						position148 := position
						depth++
						{
							position149, tokenIndex149, depth149 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l150
							}
							position++
							goto l149
						l150:
							position, tokenIndex, depth = position149, tokenIndex149, depth149
							if buffer[position] != rune('W') {
								goto l146
							}
							position++
						}
					l149:
						{
							position151, tokenIndex151, depth151 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l152
							}
							position++
							goto l151
						l152:
							position, tokenIndex, depth = position151, tokenIndex151, depth151
							if buffer[position] != rune('H') {
								goto l146
							}
							position++
						}
					l151:
						{
							position153, tokenIndex153, depth153 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l154
							}
							position++
							goto l153
						l154:
							position, tokenIndex, depth = position153, tokenIndex153, depth153
							if buffer[position] != rune('E') {
								goto l146
							}
							position++
						}
					l153:
						{
							position155, tokenIndex155, depth155 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l156
							}
							position++
							goto l155
						l156:
							position, tokenIndex, depth = position155, tokenIndex155, depth155
							if buffer[position] != rune('R') {
								goto l146
							}
							position++
						}
					l155:
						{
							position157, tokenIndex157, depth157 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l158
							}
							position++
							goto l157
						l158:
							position, tokenIndex, depth = position157, tokenIndex157, depth157
							if buffer[position] != rune('E') {
								goto l146
							}
							position++
						}
					l157:
						if !rules[rulews]() {
							goto l146
						}
						depth--
						add(ruleWHERE, position148)
					}
					goto l147
				l146:
					position, tokenIndex, depth = position146, tokenIndex146, depth146
				}
			l147:
				if !rules[rulegroupGraphPattern]() {
					goto l144
				}
				depth--
				add(rulewhereClause, position145)
			}
			return true
		l144:
			position, tokenIndex, depth = position144, tokenIndex144, depth144
			return false
		},
		/* 11 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position159, tokenIndex159, depth159 := position, tokenIndex, depth
			{
				position160 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l159
				}
				{
					position161, tokenIndex161, depth161 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l162
					}
					goto l161
				l162:
					position, tokenIndex, depth = position161, tokenIndex161, depth161
					if !rules[rulegraphPattern]() {
						goto l159
					}
				}
			l161:
				if !rules[ruleRBRACE]() {
					goto l159
				}
				depth--
				add(rulegroupGraphPattern, position160)
			}
			return true
		l159:
			position, tokenIndex, depth = position159, tokenIndex159, depth159
			return false
		},
		/* 12 graphPattern <- <(basicGraphPattern / (graphPatternNotTriples DOT? graphPattern?))> */
		func() bool {
			position163, tokenIndex163, depth163 := position, tokenIndex, depth
			{
				position164 := position
				depth++
				{
					position165, tokenIndex165, depth165 := position, tokenIndex, depth
					{
						position167 := position
						depth++
						{
							position168 := position
							depth++
							{
								position171 := position
								depth++
								{
									position172, tokenIndex172, depth172 := position, tokenIndex, depth
									{
										position174 := position
										depth++
										{
											position175, tokenIndex175, depth175 := position, tokenIndex, depth
											{
												position177 := position
												depth++
												if !rules[rulevar]() {
													goto l176
												}
												depth--
												add(rulePegText, position177)
											}
											{
												add(ruleAction0, position)
											}
											goto l175
										l176:
											position, tokenIndex, depth = position175, tokenIndex175, depth175
											{
												position180 := position
												depth++
												if !rules[rulegraphTerm]() {
													goto l179
												}
												depth--
												add(rulePegText, position180)
											}
											{
												add(ruleAction1, position)
											}
											goto l175
										l179:
											position, tokenIndex, depth = position175, tokenIndex175, depth175
											if !rules[rulepof]() {
												goto l173
											}
											{
												add(ruleAction2, position)
											}
										}
									l175:
										depth--
										add(rulevarOrTerm, position174)
									}
									if !rules[rulepropertyListPath]() {
										goto l173
									}
									goto l172
								l173:
									position, tokenIndex, depth = position172, tokenIndex172, depth172
									{
										position183 := position
										depth++
										{
											position184, tokenIndex184, depth184 := position, tokenIndex, depth
											{
												position186 := position
												depth++
												if !rules[ruleLPAREN]() {
													goto l185
												}
												if !rules[rulegraphNodePath]() {
													goto l185
												}
											l187:
												{
													position188, tokenIndex188, depth188 := position, tokenIndex, depth
													if !rules[rulegraphNodePath]() {
														goto l188
													}
													goto l187
												l188:
													position, tokenIndex, depth = position188, tokenIndex188, depth188
												}
												if !rules[ruleRPAREN]() {
													goto l185
												}
												depth--
												add(rulecollectionPath, position186)
											}
											goto l184
										l185:
											position, tokenIndex, depth = position184, tokenIndex184, depth184
											{
												position189 := position
												depth++
												{
													position190 := position
													depth++
													if buffer[position] != rune('[') {
														goto l166
													}
													position++
													if !rules[rulews]() {
														goto l166
													}
													depth--
													add(ruleLBRACK, position190)
												}
												if !rules[rulepropertyListPath]() {
													goto l166
												}
												{
													position191 := position
													depth++
													if buffer[position] != rune(']') {
														goto l166
													}
													position++
													if !rules[rulews]() {
														goto l166
													}
													depth--
													add(ruleRBRACK, position191)
												}
												depth--
												add(ruleblankNodePropertyListPath, position189)
											}
										}
									l184:
										depth--
										add(ruletriplesNodePath, position183)
									}
									if !rules[rulepropertyListPath]() {
										goto l166
									}
								}
							l172:
								{
									position192, tokenIndex192, depth192 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l192
									}
									goto l193
								l192:
									position, tokenIndex, depth = position192, tokenIndex192, depth192
								}
							l193:
								depth--
								add(ruletriplesSameSubjectPath, position171)
							}
						l169:
							{
								position170, tokenIndex170, depth170 := position, tokenIndex, depth
								{
									position194 := position
									depth++
									{
										position195, tokenIndex195, depth195 := position, tokenIndex, depth
										{
											position197 := position
											depth++
											{
												position198, tokenIndex198, depth198 := position, tokenIndex, depth
												{
													position200 := position
													depth++
													if !rules[rulevar]() {
														goto l199
													}
													depth--
													add(rulePegText, position200)
												}
												{
													add(ruleAction0, position)
												}
												goto l198
											l199:
												position, tokenIndex, depth = position198, tokenIndex198, depth198
												{
													position203 := position
													depth++
													if !rules[rulegraphTerm]() {
														goto l202
													}
													depth--
													add(rulePegText, position203)
												}
												{
													add(ruleAction1, position)
												}
												goto l198
											l202:
												position, tokenIndex, depth = position198, tokenIndex198, depth198
												if !rules[rulepof]() {
													goto l196
												}
												{
													add(ruleAction2, position)
												}
											}
										l198:
											depth--
											add(rulevarOrTerm, position197)
										}
										if !rules[rulepropertyListPath]() {
											goto l196
										}
										goto l195
									l196:
										position, tokenIndex, depth = position195, tokenIndex195, depth195
										{
											position206 := position
											depth++
											{
												position207, tokenIndex207, depth207 := position, tokenIndex, depth
												{
													position209 := position
													depth++
													if !rules[ruleLPAREN]() {
														goto l208
													}
													if !rules[rulegraphNodePath]() {
														goto l208
													}
												l210:
													{
														position211, tokenIndex211, depth211 := position, tokenIndex, depth
														if !rules[rulegraphNodePath]() {
															goto l211
														}
														goto l210
													l211:
														position, tokenIndex, depth = position211, tokenIndex211, depth211
													}
													if !rules[ruleRPAREN]() {
														goto l208
													}
													depth--
													add(rulecollectionPath, position209)
												}
												goto l207
											l208:
												position, tokenIndex, depth = position207, tokenIndex207, depth207
												{
													position212 := position
													depth++
													{
														position213 := position
														depth++
														if buffer[position] != rune('[') {
															goto l170
														}
														position++
														if !rules[rulews]() {
															goto l170
														}
														depth--
														add(ruleLBRACK, position213)
													}
													if !rules[rulepropertyListPath]() {
														goto l170
													}
													{
														position214 := position
														depth++
														if buffer[position] != rune(']') {
															goto l170
														}
														position++
														if !rules[rulews]() {
															goto l170
														}
														depth--
														add(ruleRBRACK, position214)
													}
													depth--
													add(ruleblankNodePropertyListPath, position212)
												}
											}
										l207:
											depth--
											add(ruletriplesNodePath, position206)
										}
										if !rules[rulepropertyListPath]() {
											goto l170
										}
									}
								l195:
									{
										position215, tokenIndex215, depth215 := position, tokenIndex, depth
										if !rules[ruleDOT]() {
											goto l215
										}
										goto l216
									l215:
										position, tokenIndex, depth = position215, tokenIndex215, depth215
									}
								l216:
									depth--
									add(ruletriplesSameSubjectPath, position194)
								}
								goto l169
							l170:
								position, tokenIndex, depth = position170, tokenIndex170, depth170
							}
							depth--
							add(ruletriplesBlock, position168)
						}
						depth--
						add(rulebasicGraphPattern, position167)
					}
					goto l165
				l166:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
					{
						position217 := position
						depth++
						{
							position218, tokenIndex218, depth218 := position, tokenIndex, depth
							{
								position220 := position
								depth++
								{
									position221 := position
									depth++
									{
										position222, tokenIndex222, depth222 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l223
										}
										position++
										goto l222
									l223:
										position, tokenIndex, depth = position222, tokenIndex222, depth222
										if buffer[position] != rune('O') {
											goto l219
										}
										position++
									}
								l222:
									{
										position224, tokenIndex224, depth224 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l225
										}
										position++
										goto l224
									l225:
										position, tokenIndex, depth = position224, tokenIndex224, depth224
										if buffer[position] != rune('P') {
											goto l219
										}
										position++
									}
								l224:
									{
										position226, tokenIndex226, depth226 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l227
										}
										position++
										goto l226
									l227:
										position, tokenIndex, depth = position226, tokenIndex226, depth226
										if buffer[position] != rune('T') {
											goto l219
										}
										position++
									}
								l226:
									{
										position228, tokenIndex228, depth228 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l229
										}
										position++
										goto l228
									l229:
										position, tokenIndex, depth = position228, tokenIndex228, depth228
										if buffer[position] != rune('I') {
											goto l219
										}
										position++
									}
								l228:
									{
										position230, tokenIndex230, depth230 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l231
										}
										position++
										goto l230
									l231:
										position, tokenIndex, depth = position230, tokenIndex230, depth230
										if buffer[position] != rune('O') {
											goto l219
										}
										position++
									}
								l230:
									{
										position232, tokenIndex232, depth232 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l233
										}
										position++
										goto l232
									l233:
										position, tokenIndex, depth = position232, tokenIndex232, depth232
										if buffer[position] != rune('N') {
											goto l219
										}
										position++
									}
								l232:
									{
										position234, tokenIndex234, depth234 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l235
										}
										position++
										goto l234
									l235:
										position, tokenIndex, depth = position234, tokenIndex234, depth234
										if buffer[position] != rune('A') {
											goto l219
										}
										position++
									}
								l234:
									{
										position236, tokenIndex236, depth236 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l237
										}
										position++
										goto l236
									l237:
										position, tokenIndex, depth = position236, tokenIndex236, depth236
										if buffer[position] != rune('L') {
											goto l219
										}
										position++
									}
								l236:
									if !rules[rulews]() {
										goto l219
									}
									depth--
									add(ruleOPTIONAL, position221)
								}
								if !rules[ruleLBRACE]() {
									goto l219
								}
								{
									position238, tokenIndex238, depth238 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l239
									}
									goto l238
								l239:
									position, tokenIndex, depth = position238, tokenIndex238, depth238
									if !rules[rulegraphPattern]() {
										goto l219
									}
								}
							l238:
								if !rules[ruleRBRACE]() {
									goto l219
								}
								depth--
								add(ruleoptionalGraphPattern, position220)
							}
							goto l218
						l219:
							position, tokenIndex, depth = position218, tokenIndex218, depth218
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l163
							}
						}
					l218:
						depth--
						add(rulegraphPatternNotTriples, position217)
					}
					{
						position240, tokenIndex240, depth240 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l240
						}
						goto l241
					l240:
						position, tokenIndex, depth = position240, tokenIndex240, depth240
					}
				l241:
					{
						position242, tokenIndex242, depth242 := position, tokenIndex, depth
						if !rules[rulegraphPattern]() {
							goto l242
						}
						goto l243
					l242:
						position, tokenIndex, depth = position242, tokenIndex242, depth242
					}
				l243:
				}
			l165:
				depth--
				add(rulegraphPattern, position164)
			}
			return true
		l163:
			position, tokenIndex, depth = position163, tokenIndex163, depth163
			return false
		},
		/* 13 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 14 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 15 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position246, tokenIndex246, depth246 := position, tokenIndex, depth
			{
				position247 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l246
				}
				{
					position248, tokenIndex248, depth248 := position, tokenIndex, depth
					{
						position250 := position
						depth++
						{
							position251, tokenIndex251, depth251 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l252
							}
							position++
							goto l251
						l252:
							position, tokenIndex, depth = position251, tokenIndex251, depth251
							if buffer[position] != rune('U') {
								goto l248
							}
							position++
						}
					l251:
						{
							position253, tokenIndex253, depth253 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l254
							}
							position++
							goto l253
						l254:
							position, tokenIndex, depth = position253, tokenIndex253, depth253
							if buffer[position] != rune('N') {
								goto l248
							}
							position++
						}
					l253:
						{
							position255, tokenIndex255, depth255 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l256
							}
							position++
							goto l255
						l256:
							position, tokenIndex, depth = position255, tokenIndex255, depth255
							if buffer[position] != rune('I') {
								goto l248
							}
							position++
						}
					l255:
						{
							position257, tokenIndex257, depth257 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l258
							}
							position++
							goto l257
						l258:
							position, tokenIndex, depth = position257, tokenIndex257, depth257
							if buffer[position] != rune('O') {
								goto l248
							}
							position++
						}
					l257:
						{
							position259, tokenIndex259, depth259 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l260
							}
							position++
							goto l259
						l260:
							position, tokenIndex, depth = position259, tokenIndex259, depth259
							if buffer[position] != rune('N') {
								goto l248
							}
							position++
						}
					l259:
						if !rules[rulews]() {
							goto l248
						}
						depth--
						add(ruleUNION, position250)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l248
					}
					goto l249
				l248:
					position, tokenIndex, depth = position248, tokenIndex248, depth248
				}
			l249:
				depth--
				add(rulegroupOrUnionGraphPattern, position247)
			}
			return true
		l246:
			position, tokenIndex, depth = position246, tokenIndex246, depth246
			return false
		},
		/* 16 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 17 triplesBlock <- <triplesSameSubjectPath+> */
		nil,
		/* 18 triplesSameSubjectPath <- <(((varOrTerm propertyListPath) / (triplesNodePath propertyListPath)) DOT?)> */
		nil,
		/* 19 varOrTerm <- <((<var> Action0) / (<graphTerm> Action1) / (pof Action2))> */
		nil,
		/* 20 graphTerm <- <((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('<') iri) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral))> */
		func() bool {
			position265, tokenIndex265, depth265 := position, tokenIndex, depth
			{
				position266 := position
				depth++
				{
					switch buffer[position] {
					case '(':
						{
							position268 := position
							depth++
							if buffer[position] != rune('(') {
								goto l265
							}
							position++
							if !rules[rulews]() {
								goto l265
							}
							if buffer[position] != rune(')') {
								goto l265
							}
							position++
							if !rules[rulews]() {
								goto l265
							}
							depth--
							add(rulenil, position268)
						}
						break
					case '[', '_':
						{
							position269 := position
							depth++
							{
								position270, tokenIndex270, depth270 := position, tokenIndex, depth
								{
									position272 := position
									depth++
									if buffer[position] != rune('_') {
										goto l271
									}
									position++
									if buffer[position] != rune(':') {
										goto l271
									}
									position++
									{
										switch buffer[position] {
										case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l271
											}
											position++
											break
										case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l271
											}
											position++
											break
										default:
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l271
											}
											position++
											break
										}
									}

									{
										position274, tokenIndex274, depth274 := position, tokenIndex, depth
										{
											position276, tokenIndex276, depth276 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l277
											}
											position++
											goto l276
										l277:
											position, tokenIndex, depth = position276, tokenIndex276, depth276
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l278
											}
											position++
											goto l276
										l278:
											position, tokenIndex, depth = position276, tokenIndex276, depth276
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l279
											}
											position++
											goto l276
										l279:
											position, tokenIndex, depth = position276, tokenIndex276, depth276
											if c := buffer[position]; c < rune('.') || c > rune('_') {
												goto l274
											}
											position++
										}
									l276:
										goto l275
									l274:
										position, tokenIndex, depth = position274, tokenIndex274, depth274
									}
								l275:
									if !rules[rulews]() {
										goto l271
									}
									depth--
									add(ruleblankNodeLabel, position272)
								}
								goto l270
							l271:
								position, tokenIndex, depth = position270, tokenIndex270, depth270
								{
									position280 := position
									depth++
									if buffer[position] != rune('[') {
										goto l265
									}
									position++
									if !rules[rulews]() {
										goto l265
									}
									if buffer[position] != rune(']') {
										goto l265
									}
									position++
									if !rules[rulews]() {
										goto l265
									}
									depth--
									add(ruleanon, position280)
								}
							}
						l270:
							depth--
							add(ruleblankNode, position269)
						}
						break
					case 'F', 'T', 'f', 't':
						{
							position281 := position
							depth++
							{
								position282, tokenIndex282, depth282 := position, tokenIndex, depth
								{
									position284 := position
									depth++
									{
										position285, tokenIndex285, depth285 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l286
										}
										position++
										goto l285
									l286:
										position, tokenIndex, depth = position285, tokenIndex285, depth285
										if buffer[position] != rune('T') {
											goto l283
										}
										position++
									}
								l285:
									{
										position287, tokenIndex287, depth287 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l288
										}
										position++
										goto l287
									l288:
										position, tokenIndex, depth = position287, tokenIndex287, depth287
										if buffer[position] != rune('R') {
											goto l283
										}
										position++
									}
								l287:
									{
										position289, tokenIndex289, depth289 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l290
										}
										position++
										goto l289
									l290:
										position, tokenIndex, depth = position289, tokenIndex289, depth289
										if buffer[position] != rune('U') {
											goto l283
										}
										position++
									}
								l289:
									{
										position291, tokenIndex291, depth291 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l292
										}
										position++
										goto l291
									l292:
										position, tokenIndex, depth = position291, tokenIndex291, depth291
										if buffer[position] != rune('E') {
											goto l283
										}
										position++
									}
								l291:
									if !rules[rulews]() {
										goto l283
									}
									depth--
									add(ruleTRUE, position284)
								}
								goto l282
							l283:
								position, tokenIndex, depth = position282, tokenIndex282, depth282
								{
									position293 := position
									depth++
									{
										position294, tokenIndex294, depth294 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l295
										}
										position++
										goto l294
									l295:
										position, tokenIndex, depth = position294, tokenIndex294, depth294
										if buffer[position] != rune('F') {
											goto l265
										}
										position++
									}
								l294:
									{
										position296, tokenIndex296, depth296 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l297
										}
										position++
										goto l296
									l297:
										position, tokenIndex, depth = position296, tokenIndex296, depth296
										if buffer[position] != rune('A') {
											goto l265
										}
										position++
									}
								l296:
									{
										position298, tokenIndex298, depth298 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l299
										}
										position++
										goto l298
									l299:
										position, tokenIndex, depth = position298, tokenIndex298, depth298
										if buffer[position] != rune('L') {
											goto l265
										}
										position++
									}
								l298:
									{
										position300, tokenIndex300, depth300 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l301
										}
										position++
										goto l300
									l301:
										position, tokenIndex, depth = position300, tokenIndex300, depth300
										if buffer[position] != rune('S') {
											goto l265
										}
										position++
									}
								l300:
									{
										position302, tokenIndex302, depth302 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l303
										}
										position++
										goto l302
									l303:
										position, tokenIndex, depth = position302, tokenIndex302, depth302
										if buffer[position] != rune('E') {
											goto l265
										}
										position++
									}
								l302:
									if !rules[rulews]() {
										goto l265
									}
									depth--
									add(ruleFALSE, position293)
								}
							}
						l282:
							depth--
							add(rulebooleanLiteral, position281)
						}
						break
					case '"':
						{
							position304 := position
							depth++
							{
								position305 := position
								depth++
								if buffer[position] != rune('"') {
									goto l265
								}
								position++
							l306:
								{
									position307, tokenIndex307, depth307 := position, tokenIndex, depth
									{
										position308, tokenIndex308, depth308 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											goto l308
										}
										position++
										goto l307
									l308:
										position, tokenIndex, depth = position308, tokenIndex308, depth308
									}
									if !matchDot() {
										goto l307
									}
									goto l306
								l307:
									position, tokenIndex, depth = position307, tokenIndex307, depth307
								}
								if buffer[position] != rune('"') {
									goto l265
								}
								position++
								depth--
								add(rulestring, position305)
							}
							{
								position309, tokenIndex309, depth309 := position, tokenIndex, depth
								{
									position311, tokenIndex311, depth311 := position, tokenIndex, depth
									if buffer[position] != rune('@') {
										goto l312
									}
									position++
									{
										position315, tokenIndex315, depth315 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l316
										}
										position++
										goto l315
									l316:
										position, tokenIndex, depth = position315, tokenIndex315, depth315
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l312
										}
										position++
									}
								l315:
								l313:
									{
										position314, tokenIndex314, depth314 := position, tokenIndex, depth
										{
											position317, tokenIndex317, depth317 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l318
											}
											position++
											goto l317
										l318:
											position, tokenIndex, depth = position317, tokenIndex317, depth317
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l314
											}
											position++
										}
									l317:
										goto l313
									l314:
										position, tokenIndex, depth = position314, tokenIndex314, depth314
									}
								l319:
									{
										position320, tokenIndex320, depth320 := position, tokenIndex, depth
										if buffer[position] != rune('-') {
											goto l320
										}
										position++
										{
											switch buffer[position] {
											case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l320
												}
												position++
												break
											case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l320
												}
												position++
												break
											default:
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l320
												}
												position++
												break
											}
										}

									l321:
										{
											position322, tokenIndex322, depth322 := position, tokenIndex, depth
											{
												switch buffer[position] {
												case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l322
													}
													position++
													break
												case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l322
													}
													position++
													break
												default:
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l322
													}
													position++
													break
												}
											}

											goto l321
										l322:
											position, tokenIndex, depth = position322, tokenIndex322, depth322
										}
										goto l319
									l320:
										position, tokenIndex, depth = position320, tokenIndex320, depth320
									}
									goto l311
								l312:
									position, tokenIndex, depth = position311, tokenIndex311, depth311
									if buffer[position] != rune('^') {
										goto l309
									}
									position++
									if buffer[position] != rune('^') {
										goto l309
									}
									position++
									if !rules[ruleiri]() {
										goto l309
									}
								}
							l311:
								goto l310
							l309:
								position, tokenIndex, depth = position309, tokenIndex309, depth309
							}
						l310:
							if !rules[rulews]() {
								goto l265
							}
							depth--
							add(ruleliteral, position304)
						}
						break
					case '<':
						if !rules[ruleiri]() {
							goto l265
						}
						break
					default:
						{
							position325 := position
							depth++
							{
								position326, tokenIndex326, depth326 := position, tokenIndex, depth
								{
									position328, tokenIndex328, depth328 := position, tokenIndex, depth
									if buffer[position] != rune('+') {
										goto l329
									}
									position++
									goto l328
								l329:
									position, tokenIndex, depth = position328, tokenIndex328, depth328
									if buffer[position] != rune('-') {
										goto l326
									}
									position++
								}
							l328:
								goto l327
							l326:
								position, tokenIndex, depth = position326, tokenIndex326, depth326
							}
						l327:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l265
							}
							position++
						l330:
							{
								position331, tokenIndex331, depth331 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l331
								}
								position++
								goto l330
							l331:
								position, tokenIndex, depth = position331, tokenIndex331, depth331
							}
							{
								position332, tokenIndex332, depth332 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l332
								}
								position++
							l334:
								{
									position335, tokenIndex335, depth335 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l335
									}
									position++
									goto l334
								l335:
									position, tokenIndex, depth = position335, tokenIndex335, depth335
								}
								goto l333
							l332:
								position, tokenIndex, depth = position332, tokenIndex332, depth332
							}
						l333:
							if !rules[rulews]() {
								goto l265
							}
							depth--
							add(rulenumericLiteral, position325)
						}
						break
					}
				}

				depth--
				add(rulegraphTerm, position266)
			}
			return true
		l265:
			position, tokenIndex, depth = position265, tokenIndex265, depth265
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
			position339, tokenIndex339, depth339 := position, tokenIndex, depth
			{
				position340 := position
				depth++
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l342
					}
					{
						add(ruleAction3, position)
					}
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					{
						position345 := position
						depth++
						if !rules[rulevar]() {
							goto l344
						}
						depth--
						add(rulePegText, position345)
					}
					{
						add(ruleAction4, position)
					}
					goto l341
				l344:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					{
						position347 := position
						depth++
						if !rules[rulepath]() {
							goto l339
						}
						depth--
						add(ruleverbPath, position347)
					}
				}
			l341:
				if !rules[ruleobjectListPath]() {
					goto l339
				}
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					{
						position350 := position
						depth++
						if buffer[position] != rune(';') {
							goto l348
						}
						position++
						if !rules[rulews]() {
							goto l348
						}
						depth--
						add(ruleSEMICOLON, position350)
					}
					if !rules[rulepropertyListPath]() {
						goto l348
					}
					goto l349
				l348:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
				}
			l349:
				depth--
				add(rulepropertyListPath, position340)
			}
			return true
		l339:
			position, tokenIndex, depth = position339, tokenIndex339, depth339
			return false
		},
		/* 25 verbPath <- <path> */
		nil,
		/* 26 path <- <pathAlternative> */
		func() bool {
			position352, tokenIndex352, depth352 := position, tokenIndex, depth
			{
				position353 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l352
				}
				depth--
				add(rulepath, position353)
			}
			return true
		l352:
			position, tokenIndex, depth = position352, tokenIndex352, depth352
			return false
		},
		/* 27 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position354, tokenIndex354, depth354 := position, tokenIndex, depth
			{
				position355 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l354
				}
			l356:
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l357
					}
					if !rules[rulepathAlternative]() {
						goto l357
					}
					goto l356
				l357:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
				}
				depth--
				add(rulepathAlternative, position355)
			}
			return true
		l354:
			position, tokenIndex, depth = position354, tokenIndex354, depth354
			return false
		},
		/* 28 pathSequence <- <(<pathElt> Action5 (SLASH pathSequence)*)> */
		func() bool {
			position358, tokenIndex358, depth358 := position, tokenIndex, depth
			{
				position359 := position
				depth++
				{
					position360 := position
					depth++
					{
						position361 := position
						depth++
						{
							position362, tokenIndex362, depth362 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l362
							}
							goto l363
						l362:
							position, tokenIndex, depth = position362, tokenIndex362, depth362
						}
					l363:
						{
							position364 := position
							depth++
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l358
									}
									if !rules[rulepath]() {
										goto l358
									}
									if !rules[ruleRPAREN]() {
										goto l358
									}
									break
								case '!':
									{
										position366 := position
										depth++
										if buffer[position] != rune('!') {
											goto l358
										}
										position++
										if !rules[rulews]() {
											goto l358
										}
										depth--
										add(ruleNOT, position366)
									}
									{
										position367 := position
										depth++
										{
											position368, tokenIndex368, depth368 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l369
											}
											goto l368
										l369:
											position, tokenIndex, depth = position368, tokenIndex368, depth368
											if !rules[ruleLPAREN]() {
												goto l358
											}
											{
												position370, tokenIndex370, depth370 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l370
												}
											l372:
												{
													position373, tokenIndex373, depth373 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l373
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l373
													}
													goto l372
												l373:
													position, tokenIndex, depth = position373, tokenIndex373, depth373
												}
												goto l371
											l370:
												position, tokenIndex, depth = position370, tokenIndex370, depth370
											}
										l371:
											if !rules[ruleRPAREN]() {
												goto l358
											}
										}
									l368:
										depth--
										add(rulepathNegatedPropertySet, position367)
									}
									break
								case 'a':
									if !rules[ruleISA]() {
										goto l358
									}
									break
								default:
									if !rules[ruleiri]() {
										goto l358
									}
									break
								}
							}

							depth--
							add(rulepathPrimary, position364)
						}
						{
							position374, tokenIndex374, depth374 := position, tokenIndex, depth
							{
								position375, tokenIndex375, depth375 := position, tokenIndex, depth
								{
									position377 := position
									depth++
									{
										switch buffer[position] {
										case '+':
											{
												position379 := position
												depth++
												if buffer[position] != rune('+') {
													goto l375
												}
												position++
												if !rules[rulews]() {
													goto l375
												}
												depth--
												add(rulePLUS, position379)
											}
											break
										case '?':
											{
												position380 := position
												depth++
												if buffer[position] != rune('?') {
													goto l375
												}
												position++
												if !rules[rulews]() {
													goto l375
												}
												depth--
												add(ruleQUESTION, position380)
											}
											break
										default:
											if !rules[ruleSTAR]() {
												goto l375
											}
											break
										}
									}

									depth--
									add(rulepathMod, position377)
								}
								goto l376
							l375:
								position, tokenIndex, depth = position375, tokenIndex375, depth375
							}
						l376:
							position, tokenIndex, depth = position374, tokenIndex374, depth374
						}
						depth--
						add(rulepathElt, position361)
					}
					depth--
					add(rulePegText, position360)
				}
				{
					add(ruleAction5, position)
				}
			l382:
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					{
						position384 := position
						depth++
						if buffer[position] != rune('/') {
							goto l383
						}
						position++
						if !rules[rulews]() {
							goto l383
						}
						depth--
						add(ruleSLASH, position384)
					}
					if !rules[rulepathSequence]() {
						goto l383
					}
					goto l382
				l383:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
				}
				depth--
				add(rulepathSequence, position359)
			}
			return true
		l358:
			position, tokenIndex, depth = position358, tokenIndex358, depth358
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
			position388, tokenIndex388, depth388 := position, tokenIndex, depth
			{
				position389 := position
				depth++
				{
					switch buffer[position] {
					case '^':
						if !rules[ruleINVERSE]() {
							goto l388
						}
						{
							position391, tokenIndex391, depth391 := position, tokenIndex, depth
							if !rules[ruleiri]() {
								goto l392
							}
							goto l391
						l392:
							position, tokenIndex, depth = position391, tokenIndex391, depth391
							if !rules[ruleISA]() {
								goto l388
							}
						}
					l391:
						break
					case 'a':
						if !rules[ruleISA]() {
							goto l388
						}
						break
					default:
						if !rules[ruleiri]() {
							goto l388
						}
						break
					}
				}

				depth--
				add(rulepathOneInPropertySet, position389)
			}
			return true
		l388:
			position, tokenIndex, depth = position388, tokenIndex388, depth388
			return false
		},
		/* 33 pathMod <- <((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR))> */
		nil,
		/* 34 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			{
				position395 := position
				depth++
				{
					position396 := position
					depth++
					{
						position397, tokenIndex397, depth397 := position, tokenIndex, depth
						{
							position399 := position
							depth++
							if !rules[rulegraphNodePath]() {
								goto l398
							}
							depth--
							add(rulePegText, position399)
						}
						{
							add(ruleAction6, position)
						}
						goto l397
					l398:
						position, tokenIndex, depth = position397, tokenIndex397, depth397
						if !rules[rulepof]() {
							goto l401
						}
						{
							add(ruleAction7, position)
						}
						goto l397
					l401:
						position, tokenIndex, depth = position397, tokenIndex397, depth397
						{
							add(ruleAction8, position)
						}
					}
				l397:
					depth--
					add(ruleobjectPath, position396)
				}
			l404:
				{
					position405, tokenIndex405, depth405 := position, tokenIndex, depth
					{
						position406 := position
						depth++
						if buffer[position] != rune(',') {
							goto l405
						}
						position++
						if !rules[rulews]() {
							goto l405
						}
						depth--
						add(ruleCOMMA, position406)
					}
					if !rules[ruleobjectListPath]() {
						goto l405
					}
					goto l404
				l405:
					position, tokenIndex, depth = position405, tokenIndex405, depth405
				}
				depth--
				add(ruleobjectListPath, position395)
			}
			return true
		},
		/* 35 objectPath <- <((<graphNodePath> Action6) / (pof Action7) / Action8)> */
		nil,
		/* 36 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			position408, tokenIndex408, depth408 := position, tokenIndex, depth
			{
				position409 := position
				depth++
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l411
					}
					goto l410
				l411:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
					if !rules[rulegraphTerm]() {
						goto l408
					}
				}
			l410:
				depth--
				add(rulegraphNodePath, position409)
			}
			return true
		l408:
			position, tokenIndex, depth = position408, tokenIndex408, depth408
			return false
		},
		/* 37 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 38 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 39 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position414, tokenIndex414, depth414 := position, tokenIndex, depth
			{
				position415 := position
				depth++
				{
					position416 := position
					depth++
					{
						position417, tokenIndex417, depth417 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l418
						}
						position++
						goto l417
					l418:
						position, tokenIndex, depth = position417, tokenIndex417, depth417
						if buffer[position] != rune('L') {
							goto l414
						}
						position++
					}
				l417:
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l420
						}
						position++
						goto l419
					l420:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
						if buffer[position] != rune('I') {
							goto l414
						}
						position++
					}
				l419:
					{
						position421, tokenIndex421, depth421 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l422
						}
						position++
						goto l421
					l422:
						position, tokenIndex, depth = position421, tokenIndex421, depth421
						if buffer[position] != rune('M') {
							goto l414
						}
						position++
					}
				l421:
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l424
						}
						position++
						goto l423
					l424:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
						if buffer[position] != rune('I') {
							goto l414
						}
						position++
					}
				l423:
					{
						position425, tokenIndex425, depth425 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l426
						}
						position++
						goto l425
					l426:
						position, tokenIndex, depth = position425, tokenIndex425, depth425
						if buffer[position] != rune('T') {
							goto l414
						}
						position++
					}
				l425:
					if !rules[rulews]() {
						goto l414
					}
					depth--
					add(ruleLIMIT, position416)
				}
				if !rules[ruleINTEGER]() {
					goto l414
				}
				depth--
				add(rulelimit, position415)
			}
			return true
		l414:
			position, tokenIndex, depth = position414, tokenIndex414, depth414
			return false
		},
		/* 40 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position427, tokenIndex427, depth427 := position, tokenIndex, depth
			{
				position428 := position
				depth++
				{
					position429 := position
					depth++
					{
						position430, tokenIndex430, depth430 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l431
						}
						position++
						goto l430
					l431:
						position, tokenIndex, depth = position430, tokenIndex430, depth430
						if buffer[position] != rune('O') {
							goto l427
						}
						position++
					}
				l430:
					{
						position432, tokenIndex432, depth432 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l433
						}
						position++
						goto l432
					l433:
						position, tokenIndex, depth = position432, tokenIndex432, depth432
						if buffer[position] != rune('F') {
							goto l427
						}
						position++
					}
				l432:
					{
						position434, tokenIndex434, depth434 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l435
						}
						position++
						goto l434
					l435:
						position, tokenIndex, depth = position434, tokenIndex434, depth434
						if buffer[position] != rune('F') {
							goto l427
						}
						position++
					}
				l434:
					{
						position436, tokenIndex436, depth436 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l437
						}
						position++
						goto l436
					l437:
						position, tokenIndex, depth = position436, tokenIndex436, depth436
						if buffer[position] != rune('S') {
							goto l427
						}
						position++
					}
				l436:
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l439
						}
						position++
						goto l438
					l439:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
						if buffer[position] != rune('E') {
							goto l427
						}
						position++
					}
				l438:
					{
						position440, tokenIndex440, depth440 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l441
						}
						position++
						goto l440
					l441:
						position, tokenIndex, depth = position440, tokenIndex440, depth440
						if buffer[position] != rune('T') {
							goto l427
						}
						position++
					}
				l440:
					if !rules[rulews]() {
						goto l427
					}
					depth--
					add(ruleOFFSET, position429)
				}
				if !rules[ruleINTEGER]() {
					goto l427
				}
				depth--
				add(ruleoffset, position428)
			}
			return true
		l427:
			position, tokenIndex, depth = position427, tokenIndex427, depth427
			return false
		},
		/* 41 pof <- <('<' ((&('\f') '\f') | (&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))+)> */
		func() bool {
			position442, tokenIndex442, depth442 := position, tokenIndex, depth
			{
				position443 := position
				depth++
				if buffer[position] != rune('<') {
					goto l442
				}
				position++
				{
					switch buffer[position] {
					case '\f':
						if buffer[position] != rune('\f') {
							goto l442
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l442
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l442
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l442
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l442
						}
						position++
						break
					}
				}

			l444:
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l445
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l445
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l445
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l445
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l445
							}
							position++
							break
						}
					}

					goto l444
				l445:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
				}
				depth--
				add(rulepof, position443)
			}
			return true
		l442:
			position, tokenIndex, depth = position442, tokenIndex442, depth442
			return false
		},
		/* 42 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position448, tokenIndex448, depth448 := position, tokenIndex, depth
			{
				position449 := position
				depth++
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l451
					}
					position++
					goto l450
				l451:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
					if buffer[position] != rune('$') {
						goto l448
					}
					position++
				}
			l450:
				{
					position452 := position
					depth++
					{
						position453, tokenIndex453, depth453 := position, tokenIndex, depth
						if !rules[rulePN_CHARS_U]() {
							goto l454
						}
						goto l453
					l454:
						position, tokenIndex, depth = position453, tokenIndex453, depth453
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l448
						}
						position++
					}
				l453:
				l455:
					{
						position456, tokenIndex456, depth456 := position, tokenIndex, depth
						{
							position457 := position
							depth++
							{
								position458, tokenIndex458, depth458 := position, tokenIndex, depth
								if !rules[rulePN_CHARS_U]() {
									goto l459
								}
								goto l458
							l459:
								position, tokenIndex, depth = position458, tokenIndex458, depth458
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l456
								}
								position++
							}
						l458:
							depth--
							add(ruleVAR_CHAR, position457)
						}
						goto l455
					l456:
						position, tokenIndex, depth = position456, tokenIndex456, depth456
					}
					depth--
					add(ruleVARNAME, position452)
				}
				if !rules[rulews]() {
					goto l448
				}
				depth--
				add(rulevar, position449)
			}
			return true
		l448:
			position, tokenIndex, depth = position448, tokenIndex448, depth448
			return false
		},
		/* 43 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position460, tokenIndex460, depth460 := position, tokenIndex, depth
			{
				position461 := position
				depth++
				if buffer[position] != rune('<') {
					goto l460
				}
				position++
			l462:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					{
						position464, tokenIndex464, depth464 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l464
						}
						position++
						goto l463
					l464:
						position, tokenIndex, depth = position464, tokenIndex464, depth464
					}
					if !matchDot() {
						goto l463
					}
					goto l462
				l463:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
				}
				if buffer[position] != rune('>') {
					goto l460
				}
				position++
				if !rules[rulews]() {
					goto l460
				}
				depth--
				add(ruleiri, position461)
			}
			return true
		l460:
			position, tokenIndex, depth = position460, tokenIndex460, depth460
			return false
		},
		/* 44 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iri))? ws)> */
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
		/* 52 VARNAME <- <((PN_CHARS_U / [0-9]) VAR_CHAR*)> */
		nil,
		/* 53 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				{
					position476, tokenIndex476, depth476 := position, tokenIndex, depth
					{
						position478 := position
						depth++
						{
							position479, tokenIndex479, depth479 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l480
							}
							position++
							goto l479
						l480:
							position, tokenIndex, depth = position479, tokenIndex479, depth479
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l477
							}
							position++
						}
					l479:
						depth--
						add(rulePN_CHARS_BASE, position478)
					}
					goto l476
				l477:
					position, tokenIndex, depth = position476, tokenIndex476, depth476
					if buffer[position] != rune('_') {
						goto l474
					}
					position++
				}
			l476:
				depth--
				add(rulePN_CHARS_U, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
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
			position493, tokenIndex493, depth493 := position, tokenIndex, depth
			{
				position494 := position
				depth++
				if buffer[position] != rune('{') {
					goto l493
				}
				position++
				if !rules[rulews]() {
					goto l493
				}
				depth--
				add(ruleLBRACE, position494)
			}
			return true
		l493:
			position, tokenIndex, depth = position493, tokenIndex493, depth493
			return false
		},
		/* 67 RBRACE <- <('}' ws)> */
		func() bool {
			position495, tokenIndex495, depth495 := position, tokenIndex, depth
			{
				position496 := position
				depth++
				if buffer[position] != rune('}') {
					goto l495
				}
				position++
				if !rules[rulews]() {
					goto l495
				}
				depth--
				add(ruleRBRACE, position496)
			}
			return true
		l495:
			position, tokenIndex, depth = position495, tokenIndex495, depth495
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
			position501, tokenIndex501, depth501 := position, tokenIndex, depth
			{
				position502 := position
				depth++
				if buffer[position] != rune('.') {
					goto l501
				}
				position++
				if !rules[rulews]() {
					goto l501
				}
				depth--
				add(ruleDOT, position502)
			}
			return true
		l501:
			position, tokenIndex, depth = position501, tokenIndex501, depth501
			return false
		},
		/* 73 COLON <- <(':' ws)> */
		nil,
		/* 74 PIPE <- <('|' ws)> */
		func() bool {
			position504, tokenIndex504, depth504 := position, tokenIndex, depth
			{
				position505 := position
				depth++
				if buffer[position] != rune('|') {
					goto l504
				}
				position++
				if !rules[rulews]() {
					goto l504
				}
				depth--
				add(rulePIPE, position505)
			}
			return true
		l504:
			position, tokenIndex, depth = position504, tokenIndex504, depth504
			return false
		},
		/* 75 SLASH <- <('/' ws)> */
		nil,
		/* 76 INVERSE <- <('^' ws)> */
		func() bool {
			position507, tokenIndex507, depth507 := position, tokenIndex, depth
			{
				position508 := position
				depth++
				if buffer[position] != rune('^') {
					goto l507
				}
				position++
				if !rules[rulews]() {
					goto l507
				}
				depth--
				add(ruleINVERSE, position508)
			}
			return true
		l507:
			position, tokenIndex, depth = position507, tokenIndex507, depth507
			return false
		},
		/* 77 LPAREN <- <('(' ws)> */
		func() bool {
			position509, tokenIndex509, depth509 := position, tokenIndex, depth
			{
				position510 := position
				depth++
				if buffer[position] != rune('(') {
					goto l509
				}
				position++
				if !rules[rulews]() {
					goto l509
				}
				depth--
				add(ruleLPAREN, position510)
			}
			return true
		l509:
			position, tokenIndex, depth = position509, tokenIndex509, depth509
			return false
		},
		/* 78 RPAREN <- <(')' ws)> */
		func() bool {
			position511, tokenIndex511, depth511 := position, tokenIndex, depth
			{
				position512 := position
				depth++
				if buffer[position] != rune(')') {
					goto l511
				}
				position++
				if !rules[rulews]() {
					goto l511
				}
				depth--
				add(ruleRPAREN, position512)
			}
			return true
		l511:
			position, tokenIndex, depth = position511, tokenIndex511, depth511
			return false
		},
		/* 79 ISA <- <('a' ws)> */
		func() bool {
			position513, tokenIndex513, depth513 := position, tokenIndex, depth
			{
				position514 := position
				depth++
				if buffer[position] != rune('a') {
					goto l513
				}
				position++
				if !rules[rulews]() {
					goto l513
				}
				depth--
				add(ruleISA, position514)
			}
			return true
		l513:
			position, tokenIndex, depth = position513, tokenIndex513, depth513
			return false
		},
		/* 80 NOT <- <('!' ws)> */
		nil,
		/* 81 STAR <- <('*' ws)> */
		func() bool {
			position516, tokenIndex516, depth516 := position, tokenIndex, depth
			{
				position517 := position
				depth++
				if buffer[position] != rune('*') {
					goto l516
				}
				position++
				if !rules[rulews]() {
					goto l516
				}
				depth--
				add(ruleSTAR, position517)
			}
			return true
		l516:
			position, tokenIndex, depth = position516, tokenIndex516, depth516
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
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l524
				}
				position++
			l526:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l527
					}
					position++
					goto l526
				l527:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
				}
				if !rules[rulews]() {
					goto l524
				}
				depth--
				add(ruleINTEGER, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 89 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position529 := position
				depth++
			l530:
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\v':
							if buffer[position] != rune('\v') {
								goto l531
							}
							position++
							break
						case '\f':
							if buffer[position] != rune('\f') {
								goto l531
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l531
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l531
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l531
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l531
							}
							position++
							break
						}
					}

					goto l530
				l531:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
				}
				depth--
				add(rulews, position529)
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
