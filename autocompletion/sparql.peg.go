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
	rulewhiteSpaces
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
	ruleAction9
	ruleAction10
	ruleAction11
	ruleAction12

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
	"whiteSpaces",
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
	"Action9",
	"Action10",
	"Action11",
	"Action12",

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
	*Scope

	Buffer string
	buffer []rune
	rules  [104]func() bool
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
			p.addPrefix(buffer[begin:end])
		case ruleAction1:
			p.setSubject(buffer[begin:end])
		case ruleAction2:
			p.setSubject(buffer[begin:end])
		case ruleAction3:
			p.setSubject("?POF")
		case ruleAction4:
			p.setPredicate("?POF")
		case ruleAction5:
			p.setPredicate(buffer[begin:end])
		case ruleAction6:
			p.setPredicate(buffer[begin:end])
		case ruleAction7:
			p.setObject("?POF")
			p.addTriplePattern()
		case ruleAction8:
			p.setObject(buffer[begin:end])
			p.addTriplePattern()
		case ruleAction9:
			p.setObject("?FillVar")
			p.addTriplePattern()
		case ruleAction10:
			p.setPrefix(buffer[begin:end])
		case ruleAction11:
			p.setPathLength(buffer[begin:end])
		case ruleAction12:
			p.setKeyword(buffer[begin:end])

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
									position21 := position
									depth++
								l22:
									{
										position23, tokenIndex23, depth23 := position, tokenIndex, depth
										{
											position24, tokenIndex24, depth24 := position, tokenIndex, depth
											{
												position25, tokenIndex25, depth25 := position, tokenIndex, depth
												if buffer[position] != rune(':') {
													goto l26
												}
												position++
												goto l25
											l26:
												position, tokenIndex, depth = position25, tokenIndex25, depth25
												if buffer[position] != rune(' ') {
													goto l24
												}
												position++
											}
										l25:
											goto l23
										l24:
											position, tokenIndex, depth = position24, tokenIndex24, depth24
										}
										if !matchDot() {
											goto l23
										}
										goto l22
									l23:
										position, tokenIndex, depth = position23, tokenIndex23, depth23
									}
									{
										position27 := position
										depth++
										if buffer[position] != rune(':') {
											goto l6
										}
										position++
										if !rules[rulews]() {
											goto l6
										}
										depth--
										add(ruleCOLON, position27)
									}
									if !rules[ruleiri]() {
										goto l6
									}
									depth--
									add(rulePegText, position21)
								}
								{
									add(ruleAction0, position)
								}
								depth--
								add(ruleprefixDecl, position7)
							}
							goto l5
						l6:
							position, tokenIndex, depth = position5, tokenIndex5, depth5
							{
								position29 := position
								depth++
								{
									position30 := position
									depth++
									{
										position31, tokenIndex31, depth31 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l32
										}
										position++
										goto l31
									l32:
										position, tokenIndex, depth = position31, tokenIndex31, depth31
										if buffer[position] != rune('B') {
											goto l4
										}
										position++
									}
								l31:
									{
										position33, tokenIndex33, depth33 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l34
										}
										position++
										goto l33
									l34:
										position, tokenIndex, depth = position33, tokenIndex33, depth33
										if buffer[position] != rune('A') {
											goto l4
										}
										position++
									}
								l33:
									{
										position35, tokenIndex35, depth35 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l36
										}
										position++
										goto l35
									l36:
										position, tokenIndex, depth = position35, tokenIndex35, depth35
										if buffer[position] != rune('S') {
											goto l4
										}
										position++
									}
								l35:
									{
										position37, tokenIndex37, depth37 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l38
										}
										position++
										goto l37
									l38:
										position, tokenIndex, depth = position37, tokenIndex37, depth37
										if buffer[position] != rune('E') {
											goto l4
										}
										position++
									}
								l37:
									if !rules[rulews]() {
										goto l4
									}
									depth--
									add(ruleBASE, position30)
								}
								if !rules[ruleiri]() {
									goto l4
								}
								depth--
								add(rulebaseDecl, position29)
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
					position39 := position
					depth++
					{
						position40 := position
						depth++
						if !rules[ruleselect]() {
							goto l0
						}
						{
							position41, tokenIndex41, depth41 := position, tokenIndex, depth
							{
								position43 := position
								depth++
								{
									position44 := position
									depth++
									{
										position45, tokenIndex45, depth45 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l46
										}
										position++
										goto l45
									l46:
										position, tokenIndex, depth = position45, tokenIndex45, depth45
										if buffer[position] != rune('F') {
											goto l41
										}
										position++
									}
								l45:
									{
										position47, tokenIndex47, depth47 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l48
										}
										position++
										goto l47
									l48:
										position, tokenIndex, depth = position47, tokenIndex47, depth47
										if buffer[position] != rune('R') {
											goto l41
										}
										position++
									}
								l47:
									{
										position49, tokenIndex49, depth49 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l50
										}
										position++
										goto l49
									l50:
										position, tokenIndex, depth = position49, tokenIndex49, depth49
										if buffer[position] != rune('O') {
											goto l41
										}
										position++
									}
								l49:
									{
										position51, tokenIndex51, depth51 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l52
										}
										position++
										goto l51
									l52:
										position, tokenIndex, depth = position51, tokenIndex51, depth51
										if buffer[position] != rune('M') {
											goto l41
										}
										position++
									}
								l51:
									if !rules[rulews]() {
										goto l41
									}
									depth--
									add(ruleFROM, position44)
								}
								{
									position53, tokenIndex53, depth53 := position, tokenIndex, depth
									{
										position55 := position
										depth++
										{
											position56, tokenIndex56, depth56 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l57
											}
											position++
											goto l56
										l57:
											position, tokenIndex, depth = position56, tokenIndex56, depth56
											if buffer[position] != rune('N') {
												goto l53
											}
											position++
										}
									l56:
										{
											position58, tokenIndex58, depth58 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l59
											}
											position++
											goto l58
										l59:
											position, tokenIndex, depth = position58, tokenIndex58, depth58
											if buffer[position] != rune('A') {
												goto l53
											}
											position++
										}
									l58:
										{
											position60, tokenIndex60, depth60 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l61
											}
											position++
											goto l60
										l61:
											position, tokenIndex, depth = position60, tokenIndex60, depth60
											if buffer[position] != rune('M') {
												goto l53
											}
											position++
										}
									l60:
										{
											position62, tokenIndex62, depth62 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l63
											}
											position++
											goto l62
										l63:
											position, tokenIndex, depth = position62, tokenIndex62, depth62
											if buffer[position] != rune('E') {
												goto l53
											}
											position++
										}
									l62:
										{
											position64, tokenIndex64, depth64 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l65
											}
											position++
											goto l64
										l65:
											position, tokenIndex, depth = position64, tokenIndex64, depth64
											if buffer[position] != rune('D') {
												goto l53
											}
											position++
										}
									l64:
										if !rules[rulews]() {
											goto l53
										}
										depth--
										add(ruleNAMED, position55)
									}
									goto l54
								l53:
									position, tokenIndex, depth = position53, tokenIndex53, depth53
								}
							l54:
								if !rules[ruleiriref]() {
									goto l41
								}
								depth--
								add(ruledatasetClause, position43)
							}
							goto l42
						l41:
							position, tokenIndex, depth = position41, tokenIndex41, depth41
						}
					l42:
						if !rules[rulewhereClause]() {
							goto l0
						}
						{
							position66 := position
							depth++
							{
								position67, tokenIndex67, depth67 := position, tokenIndex, depth
								{
									position69 := position
									depth++
									{
										position70, tokenIndex70, depth70 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l71
										}
										{
											position72, tokenIndex72, depth72 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l72
											}
											goto l73
										l72:
											position, tokenIndex, depth = position72, tokenIndex72, depth72
										}
									l73:
										goto l70
									l71:
										position, tokenIndex, depth = position70, tokenIndex70, depth70
										if !rules[ruleoffset]() {
											goto l67
										}
										{
											position74, tokenIndex74, depth74 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l74
											}
											goto l75
										l74:
											position, tokenIndex, depth = position74, tokenIndex74, depth74
										}
									l75:
									}
								l70:
									depth--
									add(rulelimitOffsetClauses, position69)
								}
								goto l68
							l67:
								position, tokenIndex, depth = position67, tokenIndex67, depth67
							}
						l68:
							depth--
							add(rulesolutionModifier, position66)
						}
						depth--
						add(ruleselectQuery, position40)
					}
					depth--
					add(rulequery, position39)
				}
				{
					position76, tokenIndex76, depth76 := position, tokenIndex, depth
					if !matchDot() {
						goto l76
					}
					goto l0
				l76:
					position, tokenIndex, depth = position76, tokenIndex76, depth76
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
		/* 2 prefixDecl <- <(PREFIX <((!(':' / ' ') .)* COLON iri)> Action0)> */
		nil,
		/* 3 baseDecl <- <(BASE iri)> */
		nil,
		/* 4 query <- <selectQuery> */
		nil,
		/* 5 selectQuery <- <(select datasetClause? whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position82, tokenIndex82, depth82 := position, tokenIndex, depth
			{
				position83 := position
				depth++
				{
					position84 := position
					depth++
					{
						position85, tokenIndex85, depth85 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l86
						}
						position++
						goto l85
					l86:
						position, tokenIndex, depth = position85, tokenIndex85, depth85
						if buffer[position] != rune('S') {
							goto l82
						}
						position++
					}
				l85:
					{
						position87, tokenIndex87, depth87 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l88
						}
						position++
						goto l87
					l88:
						position, tokenIndex, depth = position87, tokenIndex87, depth87
						if buffer[position] != rune('E') {
							goto l82
						}
						position++
					}
				l87:
					{
						position89, tokenIndex89, depth89 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l90
						}
						position++
						goto l89
					l90:
						position, tokenIndex, depth = position89, tokenIndex89, depth89
						if buffer[position] != rune('L') {
							goto l82
						}
						position++
					}
				l89:
					{
						position91, tokenIndex91, depth91 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l92
						}
						position++
						goto l91
					l92:
						position, tokenIndex, depth = position91, tokenIndex91, depth91
						if buffer[position] != rune('E') {
							goto l82
						}
						position++
					}
				l91:
					{
						position93, tokenIndex93, depth93 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l94
						}
						position++
						goto l93
					l94:
						position, tokenIndex, depth = position93, tokenIndex93, depth93
						if buffer[position] != rune('C') {
							goto l82
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
							goto l82
						}
						position++
					}
				l95:
					if !rules[rulews]() {
						goto l82
					}
					depth--
					add(ruleSELECT, position84)
				}
				{
					position97, tokenIndex97, depth97 := position, tokenIndex, depth
					{
						position99, tokenIndex99, depth99 := position, tokenIndex, depth
						{
							position101 := position
							depth++
							{
								position102, tokenIndex102, depth102 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l103
								}
								position++
								goto l102
							l103:
								position, tokenIndex, depth = position102, tokenIndex102, depth102
								if buffer[position] != rune('D') {
									goto l100
								}
								position++
							}
						l102:
							{
								position104, tokenIndex104, depth104 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l105
								}
								position++
								goto l104
							l105:
								position, tokenIndex, depth = position104, tokenIndex104, depth104
								if buffer[position] != rune('I') {
									goto l100
								}
								position++
							}
						l104:
							{
								position106, tokenIndex106, depth106 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l107
								}
								position++
								goto l106
							l107:
								position, tokenIndex, depth = position106, tokenIndex106, depth106
								if buffer[position] != rune('S') {
									goto l100
								}
								position++
							}
						l106:
							{
								position108, tokenIndex108, depth108 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l109
								}
								position++
								goto l108
							l109:
								position, tokenIndex, depth = position108, tokenIndex108, depth108
								if buffer[position] != rune('T') {
									goto l100
								}
								position++
							}
						l108:
							{
								position110, tokenIndex110, depth110 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l111
								}
								position++
								goto l110
							l111:
								position, tokenIndex, depth = position110, tokenIndex110, depth110
								if buffer[position] != rune('I') {
									goto l100
								}
								position++
							}
						l110:
							{
								position112, tokenIndex112, depth112 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l113
								}
								position++
								goto l112
							l113:
								position, tokenIndex, depth = position112, tokenIndex112, depth112
								if buffer[position] != rune('N') {
									goto l100
								}
								position++
							}
						l112:
							{
								position114, tokenIndex114, depth114 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l115
								}
								position++
								goto l114
							l115:
								position, tokenIndex, depth = position114, tokenIndex114, depth114
								if buffer[position] != rune('C') {
									goto l100
								}
								position++
							}
						l114:
							{
								position116, tokenIndex116, depth116 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l117
								}
								position++
								goto l116
							l117:
								position, tokenIndex, depth = position116, tokenIndex116, depth116
								if buffer[position] != rune('T') {
									goto l100
								}
								position++
							}
						l116:
							if !rules[rulews]() {
								goto l100
							}
							depth--
							add(ruleDISTINCT, position101)
						}
						goto l99
					l100:
						position, tokenIndex, depth = position99, tokenIndex99, depth99
						{
							position118 := position
							depth++
							{
								position119, tokenIndex119, depth119 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l120
								}
								position++
								goto l119
							l120:
								position, tokenIndex, depth = position119, tokenIndex119, depth119
								if buffer[position] != rune('R') {
									goto l97
								}
								position++
							}
						l119:
							{
								position121, tokenIndex121, depth121 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l122
								}
								position++
								goto l121
							l122:
								position, tokenIndex, depth = position121, tokenIndex121, depth121
								if buffer[position] != rune('E') {
									goto l97
								}
								position++
							}
						l121:
							{
								position123, tokenIndex123, depth123 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l124
								}
								position++
								goto l123
							l124:
								position, tokenIndex, depth = position123, tokenIndex123, depth123
								if buffer[position] != rune('D') {
									goto l97
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
									goto l97
								}
								position++
							}
						l125:
							{
								position127, tokenIndex127, depth127 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l128
								}
								position++
								goto l127
							l128:
								position, tokenIndex, depth = position127, tokenIndex127, depth127
								if buffer[position] != rune('C') {
									goto l97
								}
								position++
							}
						l127:
							{
								position129, tokenIndex129, depth129 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l130
								}
								position++
								goto l129
							l130:
								position, tokenIndex, depth = position129, tokenIndex129, depth129
								if buffer[position] != rune('E') {
									goto l97
								}
								position++
							}
						l129:
							{
								position131, tokenIndex131, depth131 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l132
								}
								position++
								goto l131
							l132:
								position, tokenIndex, depth = position131, tokenIndex131, depth131
								if buffer[position] != rune('D') {
									goto l97
								}
								position++
							}
						l131:
							if !rules[rulews]() {
								goto l97
							}
							depth--
							add(ruleREDUCED, position118)
						}
					}
				l99:
					goto l98
				l97:
					position, tokenIndex, depth = position97, tokenIndex97, depth97
				}
			l98:
				{
					position133, tokenIndex133, depth133 := position, tokenIndex, depth
					{
						position135 := position
						depth++
						if buffer[position] != rune('*') {
							goto l134
						}
						position++
						if !rules[rulews]() {
							goto l134
						}
						depth--
						add(ruleSTAR, position135)
					}
					goto l133
				l134:
					position, tokenIndex, depth = position133, tokenIndex133, depth133
					{
						position138 := position
						depth++
						if !rules[rulevar]() {
							goto l82
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
			l133:
				depth--
				add(ruleselect, position83)
			}
			return true
		l82:
			position, tokenIndex, depth = position82, tokenIndex82, depth82
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
		/* 9 datasetClause <- <(FROM NAMED? iriref)> */
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
		/* 12 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
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
							if !rules[ruletriplesSameSubjectPath]() {
								goto l165
							}
						l169:
							{
								position170, tokenIndex170, depth170 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l170
								}
								if !rules[ruletriplesSameSubjectPath]() {
									goto l170
								}
								goto l169
							l170:
								position, tokenIndex, depth = position170, tokenIndex170, depth170
							}
							{
								position171, tokenIndex171, depth171 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l171
								}
								goto l172
							l171:
								position, tokenIndex, depth = position171, tokenIndex171, depth171
							}
						l172:
							depth--
							add(ruletriplesBlock, position168)
						}
						depth--
						add(rulebasicGraphPattern, position167)
					}
					goto l166
				l165:
					position, tokenIndex, depth = position165, tokenIndex165, depth165
				}
			l166:
				{
					position173, tokenIndex173, depth173 := position, tokenIndex, depth
					{
						position175 := position
						depth++
						{
							position176, tokenIndex176, depth176 := position, tokenIndex, depth
							{
								position178 := position
								depth++
								{
									position179 := position
									depth++
									{
										position180, tokenIndex180, depth180 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l181
										}
										position++
										goto l180
									l181:
										position, tokenIndex, depth = position180, tokenIndex180, depth180
										if buffer[position] != rune('O') {
											goto l177
										}
										position++
									}
								l180:
									{
										position182, tokenIndex182, depth182 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l183
										}
										position++
										goto l182
									l183:
										position, tokenIndex, depth = position182, tokenIndex182, depth182
										if buffer[position] != rune('P') {
											goto l177
										}
										position++
									}
								l182:
									{
										position184, tokenIndex184, depth184 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l185
										}
										position++
										goto l184
									l185:
										position, tokenIndex, depth = position184, tokenIndex184, depth184
										if buffer[position] != rune('T') {
											goto l177
										}
										position++
									}
								l184:
									{
										position186, tokenIndex186, depth186 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l187
										}
										position++
										goto l186
									l187:
										position, tokenIndex, depth = position186, tokenIndex186, depth186
										if buffer[position] != rune('I') {
											goto l177
										}
										position++
									}
								l186:
									{
										position188, tokenIndex188, depth188 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l189
										}
										position++
										goto l188
									l189:
										position, tokenIndex, depth = position188, tokenIndex188, depth188
										if buffer[position] != rune('O') {
											goto l177
										}
										position++
									}
								l188:
									{
										position190, tokenIndex190, depth190 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l191
										}
										position++
										goto l190
									l191:
										position, tokenIndex, depth = position190, tokenIndex190, depth190
										if buffer[position] != rune('N') {
											goto l177
										}
										position++
									}
								l190:
									{
										position192, tokenIndex192, depth192 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l193
										}
										position++
										goto l192
									l193:
										position, tokenIndex, depth = position192, tokenIndex192, depth192
										if buffer[position] != rune('A') {
											goto l177
										}
										position++
									}
								l192:
									{
										position194, tokenIndex194, depth194 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l195
										}
										position++
										goto l194
									l195:
										position, tokenIndex, depth = position194, tokenIndex194, depth194
										if buffer[position] != rune('L') {
											goto l177
										}
										position++
									}
								l194:
									if !rules[rulews]() {
										goto l177
									}
									depth--
									add(ruleOPTIONAL, position179)
								}
								if !rules[ruleLBRACE]() {
									goto l177
								}
								{
									position196, tokenIndex196, depth196 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l197
									}
									goto l196
								l197:
									position, tokenIndex, depth = position196, tokenIndex196, depth196
									if !rules[rulegraphPattern]() {
										goto l177
									}
								}
							l196:
								if !rules[ruleRBRACE]() {
									goto l177
								}
								depth--
								add(ruleoptionalGraphPattern, position178)
							}
							goto l176
						l177:
							position, tokenIndex, depth = position176, tokenIndex176, depth176
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l173
							}
						}
					l176:
						depth--
						add(rulegraphPatternNotTriples, position175)
					}
					{
						position198, tokenIndex198, depth198 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l198
						}
						goto l199
					l198:
						position, tokenIndex, depth = position198, tokenIndex198, depth198
					}
				l199:
					if !rules[rulegraphPattern]() {
						goto l173
					}
					goto l174
				l173:
					position, tokenIndex, depth = position173, tokenIndex173, depth173
				}
			l174:
				depth--
				add(rulegraphPattern, position164)
			}
			return true
		},
		/* 13 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 14 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 15 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l202
				}
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					{
						position206 := position
						depth++
						{
							position207, tokenIndex207, depth207 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l208
							}
							position++
							goto l207
						l208:
							position, tokenIndex, depth = position207, tokenIndex207, depth207
							if buffer[position] != rune('U') {
								goto l204
							}
							position++
						}
					l207:
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('N') {
								goto l204
							}
							position++
						}
					l209:
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l212
							}
							position++
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if buffer[position] != rune('I') {
								goto l204
							}
							position++
						}
					l211:
						{
							position213, tokenIndex213, depth213 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l214
							}
							position++
							goto l213
						l214:
							position, tokenIndex, depth = position213, tokenIndex213, depth213
							if buffer[position] != rune('O') {
								goto l204
							}
							position++
						}
					l213:
						{
							position215, tokenIndex215, depth215 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l216
							}
							position++
							goto l215
						l216:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune('N') {
								goto l204
							}
							position++
						}
					l215:
						if !rules[rulews]() {
							goto l204
						}
						depth--
						add(ruleUNION, position206)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l204
					}
					goto l205
				l204:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
				}
			l205:
				depth--
				add(rulegroupOrUnionGraphPattern, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 16 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 17 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		nil,
		/* 18 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position219, tokenIndex219, depth219 := position, tokenIndex, depth
			{
				position220 := position
				depth++
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					{
						position223 := position
						depth++
						{
							position224, tokenIndex224, depth224 := position, tokenIndex, depth
							{
								position226 := position
								depth++
								if !rules[rulevar]() {
									goto l225
								}
								depth--
								add(rulePegText, position226)
							}
							{
								add(ruleAction1, position)
							}
							goto l224
						l225:
							position, tokenIndex, depth = position224, tokenIndex224, depth224
							{
								position229 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l228
								}
								depth--
								add(rulePegText, position229)
							}
							{
								add(ruleAction2, position)
							}
							goto l224
						l228:
							position, tokenIndex, depth = position224, tokenIndex224, depth224
							if !rules[rulepof]() {
								goto l222
							}
							{
								add(ruleAction3, position)
							}
						}
					l224:
						depth--
						add(rulevarOrTerm, position223)
					}
					if !rules[rulepropertyListPath]() {
						goto l222
					}
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					{
						position232 := position
						depth++
						{
							position233, tokenIndex233, depth233 := position, tokenIndex, depth
							{
								position235 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l234
								}
								if !rules[rulegraphNodePath]() {
									goto l234
								}
							l236:
								{
									position237, tokenIndex237, depth237 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l237
									}
									goto l236
								l237:
									position, tokenIndex, depth = position237, tokenIndex237, depth237
								}
								if !rules[ruleRPAREN]() {
									goto l234
								}
								depth--
								add(rulecollectionPath, position235)
							}
							goto l233
						l234:
							position, tokenIndex, depth = position233, tokenIndex233, depth233
							{
								position238 := position
								depth++
								{
									position239 := position
									depth++
									if buffer[position] != rune('[') {
										goto l219
									}
									position++
									if !rules[rulews]() {
										goto l219
									}
									depth--
									add(ruleLBRACK, position239)
								}
								if !rules[rulepropertyListPath]() {
									goto l219
								}
								{
									position240 := position
									depth++
									if buffer[position] != rune(']') {
										goto l219
									}
									position++
									if !rules[rulews]() {
										goto l219
									}
									depth--
									add(ruleRBRACK, position240)
								}
								depth--
								add(ruleblankNodePropertyListPath, position238)
							}
						}
					l233:
						depth--
						add(ruletriplesNodePath, position232)
					}
					if !rules[rulepropertyListPath]() {
						goto l219
					}
				}
			l221:
				depth--
				add(ruletriplesSameSubjectPath, position220)
			}
			return true
		l219:
			position, tokenIndex, depth = position219, tokenIndex219, depth219
			return false
		},
		/* 19 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 20 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position242, tokenIndex242, depth242 := position, tokenIndex, depth
			{
				position243 := position
				depth++
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l245
					}
					goto l244
				l245:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
					{
						switch buffer[position] {
						case '(':
							{
								position247 := position
								depth++
								if buffer[position] != rune('(') {
									goto l242
								}
								position++
								if !rules[rulews]() {
									goto l242
								}
								if buffer[position] != rune(')') {
									goto l242
								}
								position++
								if !rules[rulews]() {
									goto l242
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
											goto l242
										}
										position++
										if !rules[rulews]() {
											goto l242
										}
										if buffer[position] != rune(']') {
											goto l242
										}
										position++
										if !rules[rulews]() {
											goto l242
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
												goto l242
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
												goto l242
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
												goto l242
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
												goto l242
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
												goto l242
											}
											position++
										}
									l281:
										if !rules[rulews]() {
											goto l242
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
										goto l242
									}
									position++
								l285:
									{
										position286, tokenIndex286, depth286 := position, tokenIndex, depth
										{
											position287, tokenIndex287, depth287 := position, tokenIndex, depth
											if buffer[position] != rune('"') {
												goto l287
											}
											position++
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
										goto l242
									}
									position++
									depth--
									add(rulestring, position284)
								}
								{
									position288, tokenIndex288, depth288 := position, tokenIndex, depth
									{
										position290, tokenIndex290, depth290 := position, tokenIndex, depth
										if buffer[position] != rune('@') {
											goto l291
										}
										position++
										{
											position294, tokenIndex294, depth294 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l295
											}
											position++
											goto l294
										l295:
											position, tokenIndex, depth = position294, tokenIndex294, depth294
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l291
											}
											position++
										}
									l294:
									l292:
										{
											position293, tokenIndex293, depth293 := position, tokenIndex, depth
											{
												position296, tokenIndex296, depth296 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l297
												}
												position++
												goto l296
											l297:
												position, tokenIndex, depth = position296, tokenIndex296, depth296
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l293
												}
												position++
											}
										l296:
											goto l292
										l293:
											position, tokenIndex, depth = position293, tokenIndex293, depth293
										}
									l298:
										{
											position299, tokenIndex299, depth299 := position, tokenIndex, depth
											if buffer[position] != rune('-') {
												goto l299
											}
											position++
											{
												switch buffer[position] {
												case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l299
													}
													position++
													break
												case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l299
													}
													position++
													break
												default:
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l299
													}
													position++
													break
												}
											}

										l300:
											{
												position301, tokenIndex301, depth301 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l301
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l301
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l301
														}
														position++
														break
													}
												}

												goto l300
											l301:
												position, tokenIndex, depth = position301, tokenIndex301, depth301
											}
											goto l298
										l299:
											position, tokenIndex, depth = position299, tokenIndex299, depth299
										}
										goto l290
									l291:
										position, tokenIndex, depth = position290, tokenIndex290, depth290
										if buffer[position] != rune('^') {
											goto l288
										}
										position++
										if buffer[position] != rune('^') {
											goto l288
										}
										position++
										if !rules[ruleiriref]() {
											goto l288
										}
									}
								l290:
									goto l289
								l288:
									position, tokenIndex, depth = position288, tokenIndex288, depth288
								}
							l289:
								if !rules[rulews]() {
									goto l242
								}
								depth--
								add(ruleliteral, position283)
							}
							break
						default:
							{
								position304 := position
								depth++
								{
									position305, tokenIndex305, depth305 := position, tokenIndex, depth
									{
										position307, tokenIndex307, depth307 := position, tokenIndex, depth
										if buffer[position] != rune('+') {
											goto l308
										}
										position++
										goto l307
									l308:
										position, tokenIndex, depth = position307, tokenIndex307, depth307
										if buffer[position] != rune('-') {
											goto l305
										}
										position++
									}
								l307:
									goto l306
								l305:
									position, tokenIndex, depth = position305, tokenIndex305, depth305
								}
							l306:
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l242
								}
								position++
							l309:
								{
									position310, tokenIndex310, depth310 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l310
									}
									position++
									goto l309
								l310:
									position, tokenIndex, depth = position310, tokenIndex310, depth310
								}
								{
									position311, tokenIndex311, depth311 := position, tokenIndex, depth
									if buffer[position] != rune('.') {
										goto l311
									}
									position++
								l313:
									{
										position314, tokenIndex314, depth314 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l314
										}
										position++
										goto l313
									l314:
										position, tokenIndex, depth = position314, tokenIndex314, depth314
									}
									goto l312
								l311:
									position, tokenIndex, depth = position311, tokenIndex311, depth311
								}
							l312:
								if !rules[rulews]() {
									goto l242
								}
								depth--
								add(rulenumericLiteral, position304)
							}
							break
						}
					}

				}
			l244:
				depth--
				add(rulegraphTerm, position243)
			}
			return true
		l242:
			position, tokenIndex, depth = position242, tokenIndex242, depth242
			return false
		},
		/* 21 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 22 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 23 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 24 propertyListPath <- <(((pof Action4) / (<var> Action5) / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			position318, tokenIndex318, depth318 := position, tokenIndex, depth
			{
				position319 := position
				depth++
				{
					position320, tokenIndex320, depth320 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l321
					}
					{
						add(ruleAction4, position)
					}
					goto l320
				l321:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					{
						position324 := position
						depth++
						if !rules[rulevar]() {
							goto l323
						}
						depth--
						add(rulePegText, position324)
					}
					{
						add(ruleAction5, position)
					}
					goto l320
				l323:
					position, tokenIndex, depth = position320, tokenIndex320, depth320
					{
						position326 := position
						depth++
						if !rules[rulepath]() {
							goto l318
						}
						depth--
						add(ruleverbPath, position326)
					}
				}
			l320:
				if !rules[ruleobjectListPath]() {
					goto l318
				}
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					{
						position329 := position
						depth++
						if buffer[position] != rune(';') {
							goto l327
						}
						position++
						if !rules[rulews]() {
							goto l327
						}
						depth--
						add(ruleSEMICOLON, position329)
					}
					if !rules[rulepropertyListPath]() {
						goto l327
					}
					goto l328
				l327:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
				}
			l328:
				depth--
				add(rulepropertyListPath, position319)
			}
			return true
		l318:
			position, tokenIndex, depth = position318, tokenIndex318, depth318
			return false
		},
		/* 25 verbPath <- <path> */
		nil,
		/* 26 path <- <pathAlternative> */
		func() bool {
			position331, tokenIndex331, depth331 := position, tokenIndex, depth
			{
				position332 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l331
				}
				depth--
				add(rulepath, position332)
			}
			return true
		l331:
			position, tokenIndex, depth = position331, tokenIndex331, depth331
			return false
		},
		/* 27 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position333, tokenIndex333, depth333 := position, tokenIndex, depth
			{
				position334 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l333
				}
			l335:
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l336
					}
					if !rules[rulepathAlternative]() {
						goto l336
					}
					goto l335
				l336:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
				}
				depth--
				add(rulepathAlternative, position334)
			}
			return true
		l333:
			position, tokenIndex, depth = position333, tokenIndex333, depth333
			return false
		},
		/* 28 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				{
					position339 := position
					depth++
					{
						position340 := position
						depth++
						{
							position341, tokenIndex341, depth341 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l341
							}
							goto l342
						l341:
							position, tokenIndex, depth = position341, tokenIndex341, depth341
						}
					l342:
						{
							position343 := position
							depth++
							{
								position344, tokenIndex344, depth344 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l345
								}
								goto l344
							l345:
								position, tokenIndex, depth = position344, tokenIndex344, depth344
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l337
										}
										if !rules[rulepath]() {
											goto l337
										}
										if !rules[ruleRPAREN]() {
											goto l337
										}
										break
									case '!':
										{
											position347 := position
											depth++
											if buffer[position] != rune('!') {
												goto l337
											}
											position++
											if !rules[rulews]() {
												goto l337
											}
											depth--
											add(ruleNOT, position347)
										}
										{
											position348 := position
											depth++
											{
												position349, tokenIndex349, depth349 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l350
												}
												goto l349
											l350:
												position, tokenIndex, depth = position349, tokenIndex349, depth349
												if !rules[ruleLPAREN]() {
													goto l337
												}
												{
													position351, tokenIndex351, depth351 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l351
													}
												l353:
													{
														position354, tokenIndex354, depth354 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l354
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l354
														}
														goto l353
													l354:
														position, tokenIndex, depth = position354, tokenIndex354, depth354
													}
													goto l352
												l351:
													position, tokenIndex, depth = position351, tokenIndex351, depth351
												}
											l352:
												if !rules[ruleRPAREN]() {
													goto l337
												}
											}
										l349:
											depth--
											add(rulepathNegatedPropertySet, position348)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l337
										}
										break
									}
								}

							}
						l344:
							depth--
							add(rulepathPrimary, position343)
						}
						depth--
						add(rulepathElt, position340)
					}
					depth--
					add(rulePegText, position339)
				}
				{
					add(ruleAction6, position)
				}
			l356:
				{
					position357, tokenIndex357, depth357 := position, tokenIndex, depth
					{
						position358 := position
						depth++
						if buffer[position] != rune('/') {
							goto l357
						}
						position++
						if !rules[rulews]() {
							goto l357
						}
						depth--
						add(ruleSLASH, position358)
					}
					if !rules[rulepathSequence]() {
						goto l357
					}
					goto l356
				l357:
					position, tokenIndex, depth = position357, tokenIndex357, depth357
				}
				depth--
				add(rulepathSequence, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
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
			position362, tokenIndex362, depth362 := position, tokenIndex, depth
			{
				position363 := position
				depth++
				{
					position364, tokenIndex364, depth364 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l365
					}
					goto l364
				l365:
					position, tokenIndex, depth = position364, tokenIndex364, depth364
					if !rules[ruleISA]() {
						goto l366
					}
					goto l364
				l366:
					position, tokenIndex, depth = position364, tokenIndex364, depth364
					if !rules[ruleINVERSE]() {
						goto l362
					}
					{
						position367, tokenIndex367, depth367 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l368
						}
						goto l367
					l368:
						position, tokenIndex, depth = position367, tokenIndex367, depth367
						if !rules[ruleISA]() {
							goto l362
						}
					}
				l367:
				}
			l364:
				depth--
				add(rulepathOneInPropertySet, position363)
			}
			return true
		l362:
			position, tokenIndex, depth = position362, tokenIndex362, depth362
			return false
		},
		/* 33 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			{
				position370 := position
				depth++
				{
					position371 := position
					depth++
					{
						position372, tokenIndex372, depth372 := position, tokenIndex, depth
						if !rules[rulepof]() {
							goto l373
						}
						{
							add(ruleAction7, position)
						}
						goto l372
					l373:
						position, tokenIndex, depth = position372, tokenIndex372, depth372
						{
							position376 := position
							depth++
							if !rules[rulegraphNodePath]() {
								goto l375
							}
							depth--
							add(rulePegText, position376)
						}
						{
							add(ruleAction8, position)
						}
						goto l372
					l375:
						position, tokenIndex, depth = position372, tokenIndex372, depth372
						{
							add(ruleAction9, position)
						}
					}
				l372:
					depth--
					add(ruleobjectPath, position371)
				}
			l379:
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					{
						position381 := position
						depth++
						if buffer[position] != rune(',') {
							goto l380
						}
						position++
						if !rules[rulews]() {
							goto l380
						}
						depth--
						add(ruleCOMMA, position381)
					}
					if !rules[ruleobjectListPath]() {
						goto l380
					}
					goto l379
				l380:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
				}
				depth--
				add(ruleobjectListPath, position370)
			}
			return true
		},
		/* 34 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		nil,
		/* 35 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			position383, tokenIndex383, depth383 := position, tokenIndex, depth
			{
				position384 := position
				depth++
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l386
					}
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if !rules[rulegraphTerm]() {
						goto l383
					}
				}
			l385:
				depth--
				add(rulegraphNodePath, position384)
			}
			return true
		l383:
			position, tokenIndex, depth = position383, tokenIndex383, depth383
			return false
		},
		/* 36 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 37 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 38 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position389, tokenIndex389, depth389 := position, tokenIndex, depth
			{
				position390 := position
				depth++
				{
					position391 := position
					depth++
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l393
						}
						position++
						goto l392
					l393:
						position, tokenIndex, depth = position392, tokenIndex392, depth392
						if buffer[position] != rune('L') {
							goto l389
						}
						position++
					}
				l392:
					{
						position394, tokenIndex394, depth394 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l395
						}
						position++
						goto l394
					l395:
						position, tokenIndex, depth = position394, tokenIndex394, depth394
						if buffer[position] != rune('I') {
							goto l389
						}
						position++
					}
				l394:
					{
						position396, tokenIndex396, depth396 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l397
						}
						position++
						goto l396
					l397:
						position, tokenIndex, depth = position396, tokenIndex396, depth396
						if buffer[position] != rune('M') {
							goto l389
						}
						position++
					}
				l396:
					{
						position398, tokenIndex398, depth398 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l399
						}
						position++
						goto l398
					l399:
						position, tokenIndex, depth = position398, tokenIndex398, depth398
						if buffer[position] != rune('I') {
							goto l389
						}
						position++
					}
				l398:
					{
						position400, tokenIndex400, depth400 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l401
						}
						position++
						goto l400
					l401:
						position, tokenIndex, depth = position400, tokenIndex400, depth400
						if buffer[position] != rune('T') {
							goto l389
						}
						position++
					}
				l400:
					if !rules[rulews]() {
						goto l389
					}
					depth--
					add(ruleLIMIT, position391)
				}
				if !rules[ruleINTEGER]() {
					goto l389
				}
				depth--
				add(rulelimit, position390)
			}
			return true
		l389:
			position, tokenIndex, depth = position389, tokenIndex389, depth389
			return false
		},
		/* 39 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				{
					position404 := position
					depth++
					{
						position405, tokenIndex405, depth405 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l406
						}
						position++
						goto l405
					l406:
						position, tokenIndex, depth = position405, tokenIndex405, depth405
						if buffer[position] != rune('O') {
							goto l402
						}
						position++
					}
				l405:
					{
						position407, tokenIndex407, depth407 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l408
						}
						position++
						goto l407
					l408:
						position, tokenIndex, depth = position407, tokenIndex407, depth407
						if buffer[position] != rune('F') {
							goto l402
						}
						position++
					}
				l407:
					{
						position409, tokenIndex409, depth409 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l410
						}
						position++
						goto l409
					l410:
						position, tokenIndex, depth = position409, tokenIndex409, depth409
						if buffer[position] != rune('F') {
							goto l402
						}
						position++
					}
				l409:
					{
						position411, tokenIndex411, depth411 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l412
						}
						position++
						goto l411
					l412:
						position, tokenIndex, depth = position411, tokenIndex411, depth411
						if buffer[position] != rune('S') {
							goto l402
						}
						position++
					}
				l411:
					{
						position413, tokenIndex413, depth413 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l414
						}
						position++
						goto l413
					l414:
						position, tokenIndex, depth = position413, tokenIndex413, depth413
						if buffer[position] != rune('E') {
							goto l402
						}
						position++
					}
				l413:
					{
						position415, tokenIndex415, depth415 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l416
						}
						position++
						goto l415
					l416:
						position, tokenIndex, depth = position415, tokenIndex415, depth415
						if buffer[position] != rune('T') {
							goto l402
						}
						position++
					}
				l415:
					if !rules[rulews]() {
						goto l402
					}
					depth--
					add(ruleOFFSET, position404)
				}
				if !rules[ruleINTEGER]() {
					goto l402
				}
				depth--
				add(ruleoffset, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 40 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' whiteSpaces+)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					{
						position421 := position
						depth++
					l422:
						{
							position423, tokenIndex423, depth423 := position, tokenIndex, depth
							{
								position424, tokenIndex424, depth424 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l425
								}
								position++
								goto l424
							l425:
								position, tokenIndex, depth = position424, tokenIndex424, depth424
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l423
								}
								position++
							}
						l424:
							goto l422
						l423:
							position, tokenIndex, depth = position423, tokenIndex423, depth423
						}
						depth--
						add(rulePegText, position421)
					}
					if buffer[position] != rune(':') {
						goto l420
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					{
						position428 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l427
						}
						position++
					l429:
						{
							position430, tokenIndex430, depth430 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l430
							}
							position++
							goto l429
						l430:
							position, tokenIndex, depth = position430, tokenIndex430, depth430
						}
						depth--
						add(rulePegText, position428)
					}
					if buffer[position] != rune('/') {
						goto l427
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l419
				l427:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					{
						position432 := position
						depth++
					l433:
						{
							position434, tokenIndex434, depth434 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l434
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l434
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l434
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l434
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l434
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
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

							goto l433
						l434:
							position, tokenIndex, depth = position434, tokenIndex434, depth434
						}
						depth--
						add(rulePegText, position432)
					}
					{
						add(ruleAction12, position)
					}
				}
			l419:
				if buffer[position] != rune('<') {
					goto l417
				}
				position++
				if !rules[rulewhiteSpaces]() {
					goto l417
				}
			l437:
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if !rules[rulewhiteSpaces]() {
						goto l438
					}
					goto l437
				l438:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
				}
				depth--
				add(rulepof, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 41 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position439, tokenIndex439, depth439 := position, tokenIndex, depth
			{
				position440 := position
				depth++
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l442
					}
					position++
					goto l441
				l442:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
					if buffer[position] != rune('$') {
						goto l439
					}
					position++
				}
			l441:
				{
					position443 := position
					depth++
					{
						position446, tokenIndex446, depth446 := position, tokenIndex, depth
						{
							position448 := position
							depth++
							{
								position449, tokenIndex449, depth449 := position, tokenIndex, depth
								{
									position451 := position
									depth++
									{
										position452, tokenIndex452, depth452 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l453
										}
										position++
										goto l452
									l453:
										position, tokenIndex, depth = position452, tokenIndex452, depth452
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l450
										}
										position++
									}
								l452:
									depth--
									add(rulePN_CHARS_BASE, position451)
								}
								goto l449
							l450:
								position, tokenIndex, depth = position449, tokenIndex449, depth449
								if buffer[position] != rune('_') {
									goto l447
								}
								position++
							}
						l449:
							depth--
							add(rulePN_CHARS_U, position448)
						}
						goto l446
					l447:
						position, tokenIndex, depth = position446, tokenIndex446, depth446
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l439
						}
						position++
					}
				l446:
				l444:
					{
						position445, tokenIndex445, depth445 := position, tokenIndex, depth
						{
							position454, tokenIndex454, depth454 := position, tokenIndex, depth
							{
								position456 := position
								depth++
								{
									position457, tokenIndex457, depth457 := position, tokenIndex, depth
									{
										position459 := position
										depth++
										{
											position460, tokenIndex460, depth460 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l461
											}
											position++
											goto l460
										l461:
											position, tokenIndex, depth = position460, tokenIndex460, depth460
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l458
											}
											position++
										}
									l460:
										depth--
										add(rulePN_CHARS_BASE, position459)
									}
									goto l457
								l458:
									position, tokenIndex, depth = position457, tokenIndex457, depth457
									if buffer[position] != rune('_') {
										goto l455
									}
									position++
								}
							l457:
								depth--
								add(rulePN_CHARS_U, position456)
							}
							goto l454
						l455:
							position, tokenIndex, depth = position454, tokenIndex454, depth454
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l445
							}
							position++
						}
					l454:
						goto l444
					l445:
						position, tokenIndex, depth = position445, tokenIndex445, depth445
					}
					depth--
					add(ruleVARNAME, position443)
				}
				if !rules[rulews]() {
					goto l439
				}
				depth--
				add(rulevar, position440)
			}
			return true
		l439:
			position, tokenIndex, depth = position439, tokenIndex439, depth439
			return false
		},
		/* 42 iriref <- <(iri / prefixedName)> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				{
					position464, tokenIndex464, depth464 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l465
					}
					goto l464
				l465:
					position, tokenIndex, depth = position464, tokenIndex464, depth464
					{
						position466 := position
						depth++
					l467:
						{
							position468, tokenIndex468, depth468 := position, tokenIndex, depth
							{
								position469, tokenIndex469, depth469 := position, tokenIndex, depth
								{
									position470, tokenIndex470, depth470 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l471
									}
									position++
									goto l470
								l471:
									position, tokenIndex, depth = position470, tokenIndex470, depth470
									if buffer[position] != rune(' ') {
										goto l469
									}
									position++
								}
							l470:
								goto l468
							l469:
								position, tokenIndex, depth = position469, tokenIndex469, depth469
							}
							if !matchDot() {
								goto l468
							}
							goto l467
						l468:
							position, tokenIndex, depth = position468, tokenIndex468, depth468
						}
						if buffer[position] != rune(':') {
							goto l462
						}
						position++
					l472:
						{
							position473, tokenIndex473, depth473 := position, tokenIndex, depth
							{
								position474, tokenIndex474, depth474 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l475
								}
								position++
								goto l474
							l475:
								position, tokenIndex, depth = position474, tokenIndex474, depth474
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l476
								}
								position++
								goto l474
							l476:
								position, tokenIndex, depth = position474, tokenIndex474, depth474
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l477
								}
								position++
								goto l474
							l477:
								position, tokenIndex, depth = position474, tokenIndex474, depth474
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l473
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l473
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l473
										}
										position++
										break
									}
								}

							}
						l474:
							goto l472
						l473:
							position, tokenIndex, depth = position473, tokenIndex473, depth473
						}
						if !rules[rulews]() {
							goto l462
						}
						depth--
						add(ruleprefixedName, position466)
					}
				}
			l464:
				depth--
				add(ruleiriref, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 43 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				if buffer[position] != rune('<') {
					goto l479
				}
				position++
			l481:
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
					{
						position483, tokenIndex483, depth483 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l483
						}
						position++
						goto l482
					l483:
						position, tokenIndex, depth = position483, tokenIndex483, depth483
					}
					if !matchDot() {
						goto l482
					}
					goto l481
				l482:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
				}
				if buffer[position] != rune('>') {
					goto l479
				}
				position++
				if !rules[rulews]() {
					goto l479
				}
				depth--
				add(ruleiri, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 44 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* ws)> */
		nil,
		/* 45 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? ws)> */
		nil,
		/* 46 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 47 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? ws)> */
		nil,
		/* 48 booleanLiteral <- <(TRUE / FALSE)> */
		nil,
		/* 49 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 50 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? ws)> */
		nil,
		/* 51 anon <- <('[' ws ']' ws)> */
		nil,
		/* 52 nil <- <('(' ws ')' ws)> */
		nil,
		/* 53 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 54 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 55 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
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
			position506, tokenIndex506, depth506 := position, tokenIndex, depth
			{
				position507 := position
				depth++
				if buffer[position] != rune('{') {
					goto l506
				}
				position++
				if !rules[rulews]() {
					goto l506
				}
				depth--
				add(ruleLBRACE, position507)
			}
			return true
		l506:
			position, tokenIndex, depth = position506, tokenIndex506, depth506
			return false
		},
		/* 67 RBRACE <- <('}' ws)> */
		func() bool {
			position508, tokenIndex508, depth508 := position, tokenIndex, depth
			{
				position509 := position
				depth++
				if buffer[position] != rune('}') {
					goto l508
				}
				position++
				if !rules[rulews]() {
					goto l508
				}
				depth--
				add(ruleRBRACE, position509)
			}
			return true
		l508:
			position, tokenIndex, depth = position508, tokenIndex508, depth508
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
			position514, tokenIndex514, depth514 := position, tokenIndex, depth
			{
				position515 := position
				depth++
				if buffer[position] != rune('.') {
					goto l514
				}
				position++
				if !rules[rulews]() {
					goto l514
				}
				depth--
				add(ruleDOT, position515)
			}
			return true
		l514:
			position, tokenIndex, depth = position514, tokenIndex514, depth514
			return false
		},
		/* 73 COLON <- <(':' ws)> */
		nil,
		/* 74 PIPE <- <('|' ws)> */
		func() bool {
			position517, tokenIndex517, depth517 := position, tokenIndex, depth
			{
				position518 := position
				depth++
				if buffer[position] != rune('|') {
					goto l517
				}
				position++
				if !rules[rulews]() {
					goto l517
				}
				depth--
				add(rulePIPE, position518)
			}
			return true
		l517:
			position, tokenIndex, depth = position517, tokenIndex517, depth517
			return false
		},
		/* 75 SLASH <- <('/' ws)> */
		nil,
		/* 76 INVERSE <- <('^' ws)> */
		func() bool {
			position520, tokenIndex520, depth520 := position, tokenIndex, depth
			{
				position521 := position
				depth++
				if buffer[position] != rune('^') {
					goto l520
				}
				position++
				if !rules[rulews]() {
					goto l520
				}
				depth--
				add(ruleINVERSE, position521)
			}
			return true
		l520:
			position, tokenIndex, depth = position520, tokenIndex520, depth520
			return false
		},
		/* 77 LPAREN <- <('(' ws)> */
		func() bool {
			position522, tokenIndex522, depth522 := position, tokenIndex, depth
			{
				position523 := position
				depth++
				if buffer[position] != rune('(') {
					goto l522
				}
				position++
				if !rules[rulews]() {
					goto l522
				}
				depth--
				add(ruleLPAREN, position523)
			}
			return true
		l522:
			position, tokenIndex, depth = position522, tokenIndex522, depth522
			return false
		},
		/* 78 RPAREN <- <(')' ws)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				if buffer[position] != rune(')') {
					goto l524
				}
				position++
				if !rules[rulews]() {
					goto l524
				}
				depth--
				add(ruleRPAREN, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 79 ISA <- <('a' ws)> */
		func() bool {
			position526, tokenIndex526, depth526 := position, tokenIndex, depth
			{
				position527 := position
				depth++
				if buffer[position] != rune('a') {
					goto l526
				}
				position++
				if !rules[rulews]() {
					goto l526
				}
				depth--
				add(ruleISA, position527)
			}
			return true
		l526:
			position, tokenIndex, depth = position526, tokenIndex526, depth526
			return false
		},
		/* 80 NOT <- <('!' ws)> */
		nil,
		/* 81 STAR <- <('*' ws)> */
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
			position534, tokenIndex534, depth534 := position, tokenIndex, depth
			{
				position535 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l534
				}
				position++
			l536:
				{
					position537, tokenIndex537, depth537 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l537
					}
					position++
					goto l536
				l537:
					position, tokenIndex, depth = position537, tokenIndex537, depth537
				}
				if !rules[rulews]() {
					goto l534
				}
				depth--
				add(ruleINTEGER, position535)
			}
			return true
		l534:
			position, tokenIndex, depth = position534, tokenIndex534, depth534
			return false
		},
		/* 87 whiteSpaces <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position538, tokenIndex538, depth538 := position, tokenIndex, depth
			{
				position539 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l538
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l538
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l538
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l538
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l538
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l538
						}
						position++
						break
					}
				}

				depth--
				add(rulewhiteSpaces, position539)
			}
			return true
		l538:
			position, tokenIndex, depth = position538, tokenIndex538, depth538
			return false
		},
		/* 88 ws <- <whiteSpaces*> */
		func() bool {
			{
				position542 := position
				depth++
			l543:
				{
					position544, tokenIndex544, depth544 := position, tokenIndex, depth
					if !rules[rulewhiteSpaces]() {
						goto l544
					}
					goto l543
				l544:
					position, tokenIndex, depth = position544, tokenIndex544, depth544
				}
				depth--
				add(rulews, position542)
			}
			return true
		},
		nil,
		/* 91 Action0 <- <{ p.addPrefix(buffer[begin:end]) }> */
		nil,
		/* 92 Action1 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 93 Action2 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 94 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 95 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 96 Action5 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 97 Action6 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 98 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 99 Action8 <- <{ p.setObject(buffer[begin:end]); p.addTriplePattern() }> */
		nil,
		/* 100 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 101 Action10 <- <{ p.setPrefix(buffer[begin:end]) }> */
		nil,
		/* 102 Action11 <- <{ p.setPathLength(buffer[begin:end]) }> */
		nil,
		/* 103 Action12 <- <{ p.setKeyword(buffer[begin:end]) }> */
		nil,
	}
	p.rules = rules
}
