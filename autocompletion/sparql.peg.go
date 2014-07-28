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
	ruleCONSTRUCT
	ruleDESCRIBE
	ruleASK
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
	"CONSTRUCT",
	"DESCRIBE",
	"ASK",
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
	rules  [112]func() bool
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
						switch buffer[position] {
						case 'A', 'a':
							{
								position41 := position
								depth++
								{
									position42 := position
									depth++
									{
										position43, tokenIndex43, depth43 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l44
										}
										position++
										goto l43
									l44:
										position, tokenIndex, depth = position43, tokenIndex43, depth43
										if buffer[position] != rune('A') {
											goto l0
										}
										position++
									}
								l43:
									{
										position45, tokenIndex45, depth45 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l46
										}
										position++
										goto l45
									l46:
										position, tokenIndex, depth = position45, tokenIndex45, depth45
										if buffer[position] != rune('S') {
											goto l0
										}
										position++
									}
								l45:
									{
										position47, tokenIndex47, depth47 := position, tokenIndex, depth
										if buffer[position] != rune('k') {
											goto l48
										}
										position++
										goto l47
									l48:
										position, tokenIndex, depth = position47, tokenIndex47, depth47
										if buffer[position] != rune('K') {
											goto l0
										}
										position++
									}
								l47:
									if !rules[rulews]() {
										goto l0
									}
									depth--
									add(ruleASK, position42)
								}
								{
									position49, tokenIndex49, depth49 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l49
									}
									goto l50
								l49:
									position, tokenIndex, depth = position49, tokenIndex49, depth49
								}
							l50:
								if !rules[rulewhereClause]() {
									goto l0
								}
								depth--
								add(ruleaskQuery, position41)
							}
							break
						case 'D', 'd':
							{
								position51 := position
								depth++
								{
									position52 := position
									depth++
									{
										position53 := position
										depth++
										{
											position54, tokenIndex54, depth54 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l55
											}
											position++
											goto l54
										l55:
											position, tokenIndex, depth = position54, tokenIndex54, depth54
											if buffer[position] != rune('D') {
												goto l0
											}
											position++
										}
									l54:
										{
											position56, tokenIndex56, depth56 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l57
											}
											position++
											goto l56
										l57:
											position, tokenIndex, depth = position56, tokenIndex56, depth56
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l56:
										{
											position58, tokenIndex58, depth58 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l59
											}
											position++
											goto l58
										l59:
											position, tokenIndex, depth = position58, tokenIndex58, depth58
											if buffer[position] != rune('S') {
												goto l0
											}
											position++
										}
									l58:
										{
											position60, tokenIndex60, depth60 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l61
											}
											position++
											goto l60
										l61:
											position, tokenIndex, depth = position60, tokenIndex60, depth60
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l60:
										{
											position62, tokenIndex62, depth62 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l63
											}
											position++
											goto l62
										l63:
											position, tokenIndex, depth = position62, tokenIndex62, depth62
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l62:
										{
											position64, tokenIndex64, depth64 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l65
											}
											position++
											goto l64
										l65:
											position, tokenIndex, depth = position64, tokenIndex64, depth64
											if buffer[position] != rune('I') {
												goto l0
											}
											position++
										}
									l64:
										{
											position66, tokenIndex66, depth66 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l67
											}
											position++
											goto l66
										l67:
											position, tokenIndex, depth = position66, tokenIndex66, depth66
											if buffer[position] != rune('B') {
												goto l0
											}
											position++
										}
									l66:
										{
											position68, tokenIndex68, depth68 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l69
											}
											position++
											goto l68
										l69:
											position, tokenIndex, depth = position68, tokenIndex68, depth68
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l68:
										if !rules[rulews]() {
											goto l0
										}
										depth--
										add(ruleDESCRIBE, position53)
									}
									{
										position70, tokenIndex70, depth70 := position, tokenIndex, depth
										if !rules[ruleSTAR]() {
											goto l71
										}
										goto l70
									l71:
										position, tokenIndex, depth = position70, tokenIndex70, depth70
										if !rules[rulevar]() {
											goto l72
										}
										goto l70
									l72:
										position, tokenIndex, depth = position70, tokenIndex70, depth70
										if !rules[ruleiriref]() {
											goto l0
										}
									}
								l70:
									depth--
									add(ruledescribe, position52)
								}
								{
									position73, tokenIndex73, depth73 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l73
									}
									goto l74
								l73:
									position, tokenIndex, depth = position73, tokenIndex73, depth73
								}
							l74:
								{
									position75, tokenIndex75, depth75 := position, tokenIndex, depth
									if !rules[rulewhereClause]() {
										goto l75
									}
									goto l76
								l75:
									position, tokenIndex, depth = position75, tokenIndex75, depth75
								}
							l76:
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruledescribeQuery, position51)
							}
							break
						case 'C', 'c':
							{
								position77 := position
								depth++
								{
									position78 := position
									depth++
									{
										position79 := position
										depth++
										{
											position80, tokenIndex80, depth80 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l81
											}
											position++
											goto l80
										l81:
											position, tokenIndex, depth = position80, tokenIndex80, depth80
											if buffer[position] != rune('C') {
												goto l0
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
												goto l0
											}
											position++
										}
									l82:
										{
											position84, tokenIndex84, depth84 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l85
											}
											position++
											goto l84
										l85:
											position, tokenIndex, depth = position84, tokenIndex84, depth84
											if buffer[position] != rune('N') {
												goto l0
											}
											position++
										}
									l84:
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
												goto l0
											}
											position++
										}
									l86:
										{
											position88, tokenIndex88, depth88 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l89
											}
											position++
											goto l88
										l89:
											position, tokenIndex, depth = position88, tokenIndex88, depth88
											if buffer[position] != rune('T') {
												goto l0
											}
											position++
										}
									l88:
										{
											position90, tokenIndex90, depth90 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l91
											}
											position++
											goto l90
										l91:
											position, tokenIndex, depth = position90, tokenIndex90, depth90
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l90:
										{
											position92, tokenIndex92, depth92 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l93
											}
											position++
											goto l92
										l93:
											position, tokenIndex, depth = position92, tokenIndex92, depth92
											if buffer[position] != rune('U') {
												goto l0
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
												goto l0
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
												goto l0
											}
											position++
										}
									l96:
										if !rules[rulews]() {
											goto l0
										}
										depth--
										add(ruleCONSTRUCT, position79)
									}
									if !rules[ruleLBRACE]() {
										goto l0
									}
									{
										position98, tokenIndex98, depth98 := position, tokenIndex, depth
										if !rules[ruletriplesBlock]() {
											goto l98
										}
										goto l99
									l98:
										position, tokenIndex, depth = position98, tokenIndex98, depth98
									}
								l99:
									if !rules[ruleRBRACE]() {
										goto l0
									}
									depth--
									add(ruleconstruct, position78)
								}
								{
									position100, tokenIndex100, depth100 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l100
									}
									goto l101
								l100:
									position, tokenIndex, depth = position100, tokenIndex100, depth100
								}
							l101:
								if !rules[rulewhereClause]() {
									goto l0
								}
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruleconstructQuery, position77)
							}
							break
						default:
							{
								position102 := position
								depth++
								if !rules[ruleselect]() {
									goto l0
								}
								{
									position103, tokenIndex103, depth103 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l103
									}
									goto l104
								l103:
									position, tokenIndex, depth = position103, tokenIndex103, depth103
								}
							l104:
								if !rules[rulewhereClause]() {
									goto l0
								}
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruleselectQuery, position102)
							}
							break
						}
					}

					depth--
					add(rulequery, position39)
				}
				{
					position105, tokenIndex105, depth105 := position, tokenIndex, depth
					if !matchDot() {
						goto l105
					}
					goto l0
				l105:
					position, tokenIndex, depth = position105, tokenIndex105, depth105
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
		/* 4 query <- <((&('A' | 'a') askQuery) | (&('D' | 'd') describeQuery) | (&('C' | 'c') constructQuery) | (&('S' | 's') selectQuery))> */
		nil,
		/* 5 selectQuery <- <(select datasetClause? whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position111, tokenIndex111, depth111 := position, tokenIndex, depth
			{
				position112 := position
				depth++
				{
					position113 := position
					depth++
					{
						position114, tokenIndex114, depth114 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l115
						}
						position++
						goto l114
					l115:
						position, tokenIndex, depth = position114, tokenIndex114, depth114
						if buffer[position] != rune('S') {
							goto l111
						}
						position++
					}
				l114:
					{
						position116, tokenIndex116, depth116 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l117
						}
						position++
						goto l116
					l117:
						position, tokenIndex, depth = position116, tokenIndex116, depth116
						if buffer[position] != rune('E') {
							goto l111
						}
						position++
					}
				l116:
					{
						position118, tokenIndex118, depth118 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l119
						}
						position++
						goto l118
					l119:
						position, tokenIndex, depth = position118, tokenIndex118, depth118
						if buffer[position] != rune('L') {
							goto l111
						}
						position++
					}
				l118:
					{
						position120, tokenIndex120, depth120 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l121
						}
						position++
						goto l120
					l121:
						position, tokenIndex, depth = position120, tokenIndex120, depth120
						if buffer[position] != rune('E') {
							goto l111
						}
						position++
					}
				l120:
					{
						position122, tokenIndex122, depth122 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l123
						}
						position++
						goto l122
					l123:
						position, tokenIndex, depth = position122, tokenIndex122, depth122
						if buffer[position] != rune('C') {
							goto l111
						}
						position++
					}
				l122:
					{
						position124, tokenIndex124, depth124 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l125
						}
						position++
						goto l124
					l125:
						position, tokenIndex, depth = position124, tokenIndex124, depth124
						if buffer[position] != rune('T') {
							goto l111
						}
						position++
					}
				l124:
					if !rules[rulews]() {
						goto l111
					}
					depth--
					add(ruleSELECT, position113)
				}
				{
					position126, tokenIndex126, depth126 := position, tokenIndex, depth
					{
						position128, tokenIndex128, depth128 := position, tokenIndex, depth
						{
							position130 := position
							depth++
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
									goto l129
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
									goto l129
								}
								position++
							}
						l133:
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l136
								}
								position++
								goto l135
							l136:
								position, tokenIndex, depth = position135, tokenIndex135, depth135
								if buffer[position] != rune('S') {
									goto l129
								}
								position++
							}
						l135:
							{
								position137, tokenIndex137, depth137 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l138
								}
								position++
								goto l137
							l138:
								position, tokenIndex, depth = position137, tokenIndex137, depth137
								if buffer[position] != rune('T') {
									goto l129
								}
								position++
							}
						l137:
							{
								position139, tokenIndex139, depth139 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l140
								}
								position++
								goto l139
							l140:
								position, tokenIndex, depth = position139, tokenIndex139, depth139
								if buffer[position] != rune('I') {
									goto l129
								}
								position++
							}
						l139:
							{
								position141, tokenIndex141, depth141 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l142
								}
								position++
								goto l141
							l142:
								position, tokenIndex, depth = position141, tokenIndex141, depth141
								if buffer[position] != rune('N') {
									goto l129
								}
								position++
							}
						l141:
							{
								position143, tokenIndex143, depth143 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l144
								}
								position++
								goto l143
							l144:
								position, tokenIndex, depth = position143, tokenIndex143, depth143
								if buffer[position] != rune('C') {
									goto l129
								}
								position++
							}
						l143:
							{
								position145, tokenIndex145, depth145 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l146
								}
								position++
								goto l145
							l146:
								position, tokenIndex, depth = position145, tokenIndex145, depth145
								if buffer[position] != rune('T') {
									goto l129
								}
								position++
							}
						l145:
							if !rules[rulews]() {
								goto l129
							}
							depth--
							add(ruleDISTINCT, position130)
						}
						goto l128
					l129:
						position, tokenIndex, depth = position128, tokenIndex128, depth128
						{
							position147 := position
							depth++
							{
								position148, tokenIndex148, depth148 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l149
								}
								position++
								goto l148
							l149:
								position, tokenIndex, depth = position148, tokenIndex148, depth148
								if buffer[position] != rune('R') {
									goto l126
								}
								position++
							}
						l148:
							{
								position150, tokenIndex150, depth150 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l151
								}
								position++
								goto l150
							l151:
								position, tokenIndex, depth = position150, tokenIndex150, depth150
								if buffer[position] != rune('E') {
									goto l126
								}
								position++
							}
						l150:
							{
								position152, tokenIndex152, depth152 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l153
								}
								position++
								goto l152
							l153:
								position, tokenIndex, depth = position152, tokenIndex152, depth152
								if buffer[position] != rune('D') {
									goto l126
								}
								position++
							}
						l152:
							{
								position154, tokenIndex154, depth154 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l155
								}
								position++
								goto l154
							l155:
								position, tokenIndex, depth = position154, tokenIndex154, depth154
								if buffer[position] != rune('U') {
									goto l126
								}
								position++
							}
						l154:
							{
								position156, tokenIndex156, depth156 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l157
								}
								position++
								goto l156
							l157:
								position, tokenIndex, depth = position156, tokenIndex156, depth156
								if buffer[position] != rune('C') {
									goto l126
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
									goto l126
								}
								position++
							}
						l158:
							{
								position160, tokenIndex160, depth160 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l161
								}
								position++
								goto l160
							l161:
								position, tokenIndex, depth = position160, tokenIndex160, depth160
								if buffer[position] != rune('D') {
									goto l126
								}
								position++
							}
						l160:
							if !rules[rulews]() {
								goto l126
							}
							depth--
							add(ruleREDUCED, position147)
						}
					}
				l128:
					goto l127
				l126:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
				}
			l127:
				{
					position162, tokenIndex162, depth162 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l163
					}
					goto l162
				l163:
					position, tokenIndex, depth = position162, tokenIndex162, depth162
					{
						position166 := position
						depth++
						if !rules[rulevar]() {
							goto l111
						}
						depth--
						add(ruleprojectionElem, position166)
					}
				l164:
					{
						position165, tokenIndex165, depth165 := position, tokenIndex, depth
						{
							position167 := position
							depth++
							if !rules[rulevar]() {
								goto l165
							}
							depth--
							add(ruleprojectionElem, position167)
						}
						goto l164
					l165:
						position, tokenIndex, depth = position165, tokenIndex165, depth165
					}
				}
			l162:
				depth--
				add(ruleselect, position112)
			}
			return true
		l111:
			position, tokenIndex, depth = position111, tokenIndex111, depth111
			return false
		},
		/* 7 subSelect <- <(select whereClause)> */
		func() bool {
			position168, tokenIndex168, depth168 := position, tokenIndex, depth
			{
				position169 := position
				depth++
				if !rules[ruleselect]() {
					goto l168
				}
				if !rules[rulewhereClause]() {
					goto l168
				}
				depth--
				add(rulesubSelect, position169)
			}
			return true
		l168:
			position, tokenIndex, depth = position168, tokenIndex168, depth168
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
			position176, tokenIndex176, depth176 := position, tokenIndex, depth
			{
				position177 := position
				depth++
				{
					position178 := position
					depth++
					{
						position179, tokenIndex179, depth179 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l180
						}
						position++
						goto l179
					l180:
						position, tokenIndex, depth = position179, tokenIndex179, depth179
						if buffer[position] != rune('F') {
							goto l176
						}
						position++
					}
				l179:
					{
						position181, tokenIndex181, depth181 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l182
						}
						position++
						goto l181
					l182:
						position, tokenIndex, depth = position181, tokenIndex181, depth181
						if buffer[position] != rune('R') {
							goto l176
						}
						position++
					}
				l181:
					{
						position183, tokenIndex183, depth183 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l184
						}
						position++
						goto l183
					l184:
						position, tokenIndex, depth = position183, tokenIndex183, depth183
						if buffer[position] != rune('O') {
							goto l176
						}
						position++
					}
				l183:
					{
						position185, tokenIndex185, depth185 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l186
						}
						position++
						goto l185
					l186:
						position, tokenIndex, depth = position185, tokenIndex185, depth185
						if buffer[position] != rune('M') {
							goto l176
						}
						position++
					}
				l185:
					if !rules[rulews]() {
						goto l176
					}
					depth--
					add(ruleFROM, position178)
				}
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					{
						position189 := position
						depth++
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
								goto l187
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
								goto l187
							}
							position++
						}
					l192:
						{
							position194, tokenIndex194, depth194 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l195
							}
							position++
							goto l194
						l195:
							position, tokenIndex, depth = position194, tokenIndex194, depth194
							if buffer[position] != rune('M') {
								goto l187
							}
							position++
						}
					l194:
						{
							position196, tokenIndex196, depth196 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l197
							}
							position++
							goto l196
						l197:
							position, tokenIndex, depth = position196, tokenIndex196, depth196
							if buffer[position] != rune('E') {
								goto l187
							}
							position++
						}
					l196:
						{
							position198, tokenIndex198, depth198 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l199
							}
							position++
							goto l198
						l199:
							position, tokenIndex, depth = position198, tokenIndex198, depth198
							if buffer[position] != rune('D') {
								goto l187
							}
							position++
						}
					l198:
						if !rules[rulews]() {
							goto l187
						}
						depth--
						add(ruleNAMED, position189)
					}
					goto l188
				l187:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
				}
			l188:
				if !rules[ruleiriref]() {
					goto l176
				}
				depth--
				add(ruledatasetClause, position177)
			}
			return true
		l176:
			position, tokenIndex, depth = position176, tokenIndex176, depth176
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					{
						position204 := position
						depth++
						{
							position205, tokenIndex205, depth205 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l206
							}
							position++
							goto l205
						l206:
							position, tokenIndex, depth = position205, tokenIndex205, depth205
							if buffer[position] != rune('W') {
								goto l202
							}
							position++
						}
					l205:
						{
							position207, tokenIndex207, depth207 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l208
							}
							position++
							goto l207
						l208:
							position, tokenIndex, depth = position207, tokenIndex207, depth207
							if buffer[position] != rune('H') {
								goto l202
							}
							position++
						}
					l207:
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('E') {
								goto l202
							}
							position++
						}
					l209:
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l212
							}
							position++
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if buffer[position] != rune('R') {
								goto l202
							}
							position++
						}
					l211:
						{
							position213, tokenIndex213, depth213 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l214
							}
							position++
							goto l213
						l214:
							position, tokenIndex, depth = position213, tokenIndex213, depth213
							if buffer[position] != rune('E') {
								goto l202
							}
							position++
						}
					l213:
						if !rules[rulews]() {
							goto l202
						}
						depth--
						add(ruleWHERE, position204)
					}
					goto l203
				l202:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
				}
			l203:
				if !rules[rulegroupGraphPattern]() {
					goto l200
				}
				depth--
				add(rulewhereClause, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position215, tokenIndex215, depth215 := position, tokenIndex, depth
			{
				position216 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l215
				}
				{
					position217, tokenIndex217, depth217 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l218
					}
					goto l217
				l218:
					position, tokenIndex, depth = position217, tokenIndex217, depth217
					if !rules[rulegraphPattern]() {
						goto l215
					}
				}
			l217:
				if !rules[ruleRBRACE]() {
					goto l215
				}
				depth--
				add(rulegroupGraphPattern, position216)
			}
			return true
		l215:
			position, tokenIndex, depth = position215, tokenIndex215, depth215
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position220 := position
				depth++
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					{
						position223 := position
						depth++
						if !rules[ruletriplesBlock]() {
							goto l221
						}
						depth--
						add(rulebasicGraphPattern, position223)
					}
					goto l222
				l221:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
				}
			l222:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					{
						position226 := position
						depth++
						{
							position227, tokenIndex227, depth227 := position, tokenIndex, depth
							{
								position229 := position
								depth++
								{
									position230 := position
									depth++
									{
										position231, tokenIndex231, depth231 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l232
										}
										position++
										goto l231
									l232:
										position, tokenIndex, depth = position231, tokenIndex231, depth231
										if buffer[position] != rune('O') {
											goto l228
										}
										position++
									}
								l231:
									{
										position233, tokenIndex233, depth233 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l234
										}
										position++
										goto l233
									l234:
										position, tokenIndex, depth = position233, tokenIndex233, depth233
										if buffer[position] != rune('P') {
											goto l228
										}
										position++
									}
								l233:
									{
										position235, tokenIndex235, depth235 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l236
										}
										position++
										goto l235
									l236:
										position, tokenIndex, depth = position235, tokenIndex235, depth235
										if buffer[position] != rune('T') {
											goto l228
										}
										position++
									}
								l235:
									{
										position237, tokenIndex237, depth237 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l238
										}
										position++
										goto l237
									l238:
										position, tokenIndex, depth = position237, tokenIndex237, depth237
										if buffer[position] != rune('I') {
											goto l228
										}
										position++
									}
								l237:
									{
										position239, tokenIndex239, depth239 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l240
										}
										position++
										goto l239
									l240:
										position, tokenIndex, depth = position239, tokenIndex239, depth239
										if buffer[position] != rune('O') {
											goto l228
										}
										position++
									}
								l239:
									{
										position241, tokenIndex241, depth241 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l242
										}
										position++
										goto l241
									l242:
										position, tokenIndex, depth = position241, tokenIndex241, depth241
										if buffer[position] != rune('N') {
											goto l228
										}
										position++
									}
								l241:
									{
										position243, tokenIndex243, depth243 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l244
										}
										position++
										goto l243
									l244:
										position, tokenIndex, depth = position243, tokenIndex243, depth243
										if buffer[position] != rune('A') {
											goto l228
										}
										position++
									}
								l243:
									{
										position245, tokenIndex245, depth245 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l246
										}
										position++
										goto l245
									l246:
										position, tokenIndex, depth = position245, tokenIndex245, depth245
										if buffer[position] != rune('L') {
											goto l228
										}
										position++
									}
								l245:
									if !rules[rulews]() {
										goto l228
									}
									depth--
									add(ruleOPTIONAL, position230)
								}
								if !rules[ruleLBRACE]() {
									goto l228
								}
								{
									position247, tokenIndex247, depth247 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l248
									}
									goto l247
								l248:
									position, tokenIndex, depth = position247, tokenIndex247, depth247
									if !rules[rulegraphPattern]() {
										goto l228
									}
								}
							l247:
								if !rules[ruleRBRACE]() {
									goto l228
								}
								depth--
								add(ruleoptionalGraphPattern, position229)
							}
							goto l227
						l228:
							position, tokenIndex, depth = position227, tokenIndex227, depth227
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l224
							}
						}
					l227:
						depth--
						add(rulegraphPatternNotTriples, position226)
					}
					{
						position249, tokenIndex249, depth249 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l249
						}
						goto l250
					l249:
						position, tokenIndex, depth = position249, tokenIndex249, depth249
					}
				l250:
					if !rules[rulegraphPattern]() {
						goto l224
					}
					goto l225
				l224:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
				}
			l225:
				depth--
				add(rulegraphPattern, position220)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position253, tokenIndex253, depth253 := position, tokenIndex, depth
			{
				position254 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l253
				}
				{
					position255, tokenIndex255, depth255 := position, tokenIndex, depth
					{
						position257 := position
						depth++
						{
							position258, tokenIndex258, depth258 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l259
							}
							position++
							goto l258
						l259:
							position, tokenIndex, depth = position258, tokenIndex258, depth258
							if buffer[position] != rune('U') {
								goto l255
							}
							position++
						}
					l258:
						{
							position260, tokenIndex260, depth260 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l261
							}
							position++
							goto l260
						l261:
							position, tokenIndex, depth = position260, tokenIndex260, depth260
							if buffer[position] != rune('N') {
								goto l255
							}
							position++
						}
					l260:
						{
							position262, tokenIndex262, depth262 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l263
							}
							position++
							goto l262
						l263:
							position, tokenIndex, depth = position262, tokenIndex262, depth262
							if buffer[position] != rune('I') {
								goto l255
							}
							position++
						}
					l262:
						{
							position264, tokenIndex264, depth264 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l265
							}
							position++
							goto l264
						l265:
							position, tokenIndex, depth = position264, tokenIndex264, depth264
							if buffer[position] != rune('O') {
								goto l255
							}
							position++
						}
					l264:
						{
							position266, tokenIndex266, depth266 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l267
							}
							position++
							goto l266
						l267:
							position, tokenIndex, depth = position266, tokenIndex266, depth266
							if buffer[position] != rune('N') {
								goto l255
							}
							position++
						}
					l266:
						if !rules[rulews]() {
							goto l255
						}
						depth--
						add(ruleUNION, position257)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l255
					}
					goto l256
				l255:
					position, tokenIndex, depth = position255, tokenIndex255, depth255
				}
			l256:
				depth--
				add(rulegroupOrUnionGraphPattern, position254)
			}
			return true
		l253:
			position, tokenIndex, depth = position253, tokenIndex253, depth253
			return false
		},
		/* 21 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 22 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position269, tokenIndex269, depth269 := position, tokenIndex, depth
			{
				position270 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l269
				}
			l271:
				{
					position272, tokenIndex272, depth272 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l272
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l272
					}
					goto l271
				l272:
					position, tokenIndex, depth = position272, tokenIndex272, depth272
				}
				{
					position273, tokenIndex273, depth273 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l273
					}
					goto l274
				l273:
					position, tokenIndex, depth = position273, tokenIndex273, depth273
				}
			l274:
				depth--
				add(ruletriplesBlock, position270)
			}
			return true
		l269:
			position, tokenIndex, depth = position269, tokenIndex269, depth269
			return false
		},
		/* 23 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position275, tokenIndex275, depth275 := position, tokenIndex, depth
			{
				position276 := position
				depth++
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					{
						position279 := position
						depth++
						{
							position280, tokenIndex280, depth280 := position, tokenIndex, depth
							{
								position282 := position
								depth++
								if !rules[rulevar]() {
									goto l281
								}
								depth--
								add(rulePegText, position282)
							}
							{
								add(ruleAction1, position)
							}
							goto l280
						l281:
							position, tokenIndex, depth = position280, tokenIndex280, depth280
							{
								position285 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l284
								}
								depth--
								add(rulePegText, position285)
							}
							{
								add(ruleAction2, position)
							}
							goto l280
						l284:
							position, tokenIndex, depth = position280, tokenIndex280, depth280
							if !rules[rulepof]() {
								goto l278
							}
							{
								add(ruleAction3, position)
							}
						}
					l280:
						depth--
						add(rulevarOrTerm, position279)
					}
					if !rules[rulepropertyListPath]() {
						goto l278
					}
					goto l277
				l278:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
					{
						position288 := position
						depth++
						{
							position289, tokenIndex289, depth289 := position, tokenIndex, depth
							{
								position291 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l290
								}
								if !rules[rulegraphNodePath]() {
									goto l290
								}
							l292:
								{
									position293, tokenIndex293, depth293 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l293
									}
									goto l292
								l293:
									position, tokenIndex, depth = position293, tokenIndex293, depth293
								}
								if !rules[ruleRPAREN]() {
									goto l290
								}
								depth--
								add(rulecollectionPath, position291)
							}
							goto l289
						l290:
							position, tokenIndex, depth = position289, tokenIndex289, depth289
							{
								position294 := position
								depth++
								{
									position295 := position
									depth++
									if buffer[position] != rune('[') {
										goto l275
									}
									position++
									if !rules[rulews]() {
										goto l275
									}
									depth--
									add(ruleLBRACK, position295)
								}
								if !rules[rulepropertyListPath]() {
									goto l275
								}
								{
									position296 := position
									depth++
									if buffer[position] != rune(']') {
										goto l275
									}
									position++
									if !rules[rulews]() {
										goto l275
									}
									depth--
									add(ruleRBRACK, position296)
								}
								depth--
								add(ruleblankNodePropertyListPath, position294)
							}
						}
					l289:
						depth--
						add(ruletriplesNodePath, position288)
					}
					if !rules[rulepropertyListPath]() {
						goto l275
					}
				}
			l277:
				depth--
				add(ruletriplesSameSubjectPath, position276)
			}
			return true
		l275:
			position, tokenIndex, depth = position275, tokenIndex275, depth275
			return false
		},
		/* 24 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 25 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position298, tokenIndex298, depth298 := position, tokenIndex, depth
			{
				position299 := position
				depth++
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l301
					}
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					{
						switch buffer[position] {
						case '(':
							{
								position303 := position
								depth++
								if buffer[position] != rune('(') {
									goto l298
								}
								position++
								if !rules[rulews]() {
									goto l298
								}
								if buffer[position] != rune(')') {
									goto l298
								}
								position++
								if !rules[rulews]() {
									goto l298
								}
								depth--
								add(rulenil, position303)
							}
							break
						case '[', '_':
							{
								position304 := position
								depth++
								{
									position305, tokenIndex305, depth305 := position, tokenIndex, depth
									{
										position307 := position
										depth++
										if buffer[position] != rune('_') {
											goto l306
										}
										position++
										if buffer[position] != rune(':') {
											goto l306
										}
										position++
										{
											switch buffer[position] {
											case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l306
												}
												position++
												break
											case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l306
												}
												position++
												break
											default:
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l306
												}
												position++
												break
											}
										}

										{
											position309, tokenIndex309, depth309 := position, tokenIndex, depth
											{
												position311, tokenIndex311, depth311 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l312
												}
												position++
												goto l311
											l312:
												position, tokenIndex, depth = position311, tokenIndex311, depth311
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l313
												}
												position++
												goto l311
											l313:
												position, tokenIndex, depth = position311, tokenIndex311, depth311
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l314
												}
												position++
												goto l311
											l314:
												position, tokenIndex, depth = position311, tokenIndex311, depth311
												if c := buffer[position]; c < rune('.') || c > rune('_') {
													goto l309
												}
												position++
											}
										l311:
											goto l310
										l309:
											position, tokenIndex, depth = position309, tokenIndex309, depth309
										}
									l310:
										if !rules[rulews]() {
											goto l306
										}
										depth--
										add(ruleblankNodeLabel, position307)
									}
									goto l305
								l306:
									position, tokenIndex, depth = position305, tokenIndex305, depth305
									{
										position315 := position
										depth++
										if buffer[position] != rune('[') {
											goto l298
										}
										position++
										if !rules[rulews]() {
											goto l298
										}
										if buffer[position] != rune(']') {
											goto l298
										}
										position++
										if !rules[rulews]() {
											goto l298
										}
										depth--
										add(ruleanon, position315)
									}
								}
							l305:
								depth--
								add(ruleblankNode, position304)
							}
							break
						case 'F', 'T', 'f', 't':
							{
								position316 := position
								depth++
								{
									position317, tokenIndex317, depth317 := position, tokenIndex, depth
									{
										position319 := position
										depth++
										{
											position320, tokenIndex320, depth320 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l321
											}
											position++
											goto l320
										l321:
											position, tokenIndex, depth = position320, tokenIndex320, depth320
											if buffer[position] != rune('T') {
												goto l318
											}
											position++
										}
									l320:
										{
											position322, tokenIndex322, depth322 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l323
											}
											position++
											goto l322
										l323:
											position, tokenIndex, depth = position322, tokenIndex322, depth322
											if buffer[position] != rune('R') {
												goto l318
											}
											position++
										}
									l322:
										{
											position324, tokenIndex324, depth324 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l325
											}
											position++
											goto l324
										l325:
											position, tokenIndex, depth = position324, tokenIndex324, depth324
											if buffer[position] != rune('U') {
												goto l318
											}
											position++
										}
									l324:
										{
											position326, tokenIndex326, depth326 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l327
											}
											position++
											goto l326
										l327:
											position, tokenIndex, depth = position326, tokenIndex326, depth326
											if buffer[position] != rune('E') {
												goto l318
											}
											position++
										}
									l326:
										if !rules[rulews]() {
											goto l318
										}
										depth--
										add(ruleTRUE, position319)
									}
									goto l317
								l318:
									position, tokenIndex, depth = position317, tokenIndex317, depth317
									{
										position328 := position
										depth++
										{
											position329, tokenIndex329, depth329 := position, tokenIndex, depth
											if buffer[position] != rune('f') {
												goto l330
											}
											position++
											goto l329
										l330:
											position, tokenIndex, depth = position329, tokenIndex329, depth329
											if buffer[position] != rune('F') {
												goto l298
											}
											position++
										}
									l329:
										{
											position331, tokenIndex331, depth331 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l332
											}
											position++
											goto l331
										l332:
											position, tokenIndex, depth = position331, tokenIndex331, depth331
											if buffer[position] != rune('A') {
												goto l298
											}
											position++
										}
									l331:
										{
											position333, tokenIndex333, depth333 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l334
											}
											position++
											goto l333
										l334:
											position, tokenIndex, depth = position333, tokenIndex333, depth333
											if buffer[position] != rune('L') {
												goto l298
											}
											position++
										}
									l333:
										{
											position335, tokenIndex335, depth335 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l336
											}
											position++
											goto l335
										l336:
											position, tokenIndex, depth = position335, tokenIndex335, depth335
											if buffer[position] != rune('S') {
												goto l298
											}
											position++
										}
									l335:
										{
											position337, tokenIndex337, depth337 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l338
											}
											position++
											goto l337
										l338:
											position, tokenIndex, depth = position337, tokenIndex337, depth337
											if buffer[position] != rune('E') {
												goto l298
											}
											position++
										}
									l337:
										if !rules[rulews]() {
											goto l298
										}
										depth--
										add(ruleFALSE, position328)
									}
								}
							l317:
								depth--
								add(rulebooleanLiteral, position316)
							}
							break
						case '"':
							{
								position339 := position
								depth++
								{
									position340 := position
									depth++
									if buffer[position] != rune('"') {
										goto l298
									}
									position++
								l341:
									{
										position342, tokenIndex342, depth342 := position, tokenIndex, depth
										{
											position343, tokenIndex343, depth343 := position, tokenIndex, depth
											if buffer[position] != rune('"') {
												goto l343
											}
											position++
											goto l342
										l343:
											position, tokenIndex, depth = position343, tokenIndex343, depth343
										}
										if !matchDot() {
											goto l342
										}
										goto l341
									l342:
										position, tokenIndex, depth = position342, tokenIndex342, depth342
									}
									if buffer[position] != rune('"') {
										goto l298
									}
									position++
									depth--
									add(rulestring, position340)
								}
								{
									position344, tokenIndex344, depth344 := position, tokenIndex, depth
									{
										position346, tokenIndex346, depth346 := position, tokenIndex, depth
										if buffer[position] != rune('@') {
											goto l347
										}
										position++
										{
											position350, tokenIndex350, depth350 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l351
											}
											position++
											goto l350
										l351:
											position, tokenIndex, depth = position350, tokenIndex350, depth350
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l347
											}
											position++
										}
									l350:
									l348:
										{
											position349, tokenIndex349, depth349 := position, tokenIndex, depth
											{
												position352, tokenIndex352, depth352 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l353
												}
												position++
												goto l352
											l353:
												position, tokenIndex, depth = position352, tokenIndex352, depth352
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l349
												}
												position++
											}
										l352:
											goto l348
										l349:
											position, tokenIndex, depth = position349, tokenIndex349, depth349
										}
									l354:
										{
											position355, tokenIndex355, depth355 := position, tokenIndex, depth
											if buffer[position] != rune('-') {
												goto l355
											}
											position++
											{
												switch buffer[position] {
												case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l355
													}
													position++
													break
												case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l355
													}
													position++
													break
												default:
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l355
													}
													position++
													break
												}
											}

										l356:
											{
												position357, tokenIndex357, depth357 := position, tokenIndex, depth
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l357
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l357
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l357
														}
														position++
														break
													}
												}

												goto l356
											l357:
												position, tokenIndex, depth = position357, tokenIndex357, depth357
											}
											goto l354
										l355:
											position, tokenIndex, depth = position355, tokenIndex355, depth355
										}
										goto l346
									l347:
										position, tokenIndex, depth = position346, tokenIndex346, depth346
										if buffer[position] != rune('^') {
											goto l344
										}
										position++
										if buffer[position] != rune('^') {
											goto l344
										}
										position++
										if !rules[ruleiriref]() {
											goto l344
										}
									}
								l346:
									goto l345
								l344:
									position, tokenIndex, depth = position344, tokenIndex344, depth344
								}
							l345:
								if !rules[rulews]() {
									goto l298
								}
								depth--
								add(ruleliteral, position339)
							}
							break
						default:
							{
								position360 := position
								depth++
								{
									position361, tokenIndex361, depth361 := position, tokenIndex, depth
									{
										position363, tokenIndex363, depth363 := position, tokenIndex, depth
										if buffer[position] != rune('+') {
											goto l364
										}
										position++
										goto l363
									l364:
										position, tokenIndex, depth = position363, tokenIndex363, depth363
										if buffer[position] != rune('-') {
											goto l361
										}
										position++
									}
								l363:
									goto l362
								l361:
									position, tokenIndex, depth = position361, tokenIndex361, depth361
								}
							l362:
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l298
								}
								position++
							l365:
								{
									position366, tokenIndex366, depth366 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l366
									}
									position++
									goto l365
								l366:
									position, tokenIndex, depth = position366, tokenIndex366, depth366
								}
								{
									position367, tokenIndex367, depth367 := position, tokenIndex, depth
									if buffer[position] != rune('.') {
										goto l367
									}
									position++
								l369:
									{
										position370, tokenIndex370, depth370 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l370
										}
										position++
										goto l369
									l370:
										position, tokenIndex, depth = position370, tokenIndex370, depth370
									}
									goto l368
								l367:
									position, tokenIndex, depth = position367, tokenIndex367, depth367
								}
							l368:
								if !rules[rulews]() {
									goto l298
								}
								depth--
								add(rulenumericLiteral, position360)
							}
							break
						}
					}

				}
			l300:
				depth--
				add(rulegraphTerm, position299)
			}
			return true
		l298:
			position, tokenIndex, depth = position298, tokenIndex298, depth298
			return false
		},
		/* 26 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 27 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 28 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 29 propertyListPath <- <(((pof Action4) / (<var> Action5) / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l377
					}
					{
						add(ruleAction4, position)
					}
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					{
						position380 := position
						depth++
						if !rules[rulevar]() {
							goto l379
						}
						depth--
						add(rulePegText, position380)
					}
					{
						add(ruleAction5, position)
					}
					goto l376
				l379:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					{
						position382 := position
						depth++
						if !rules[rulepath]() {
							goto l374
						}
						depth--
						add(ruleverbPath, position382)
					}
				}
			l376:
				if !rules[ruleobjectListPath]() {
					goto l374
				}
				{
					position383, tokenIndex383, depth383 := position, tokenIndex, depth
					{
						position385 := position
						depth++
						if buffer[position] != rune(';') {
							goto l383
						}
						position++
						if !rules[rulews]() {
							goto l383
						}
						depth--
						add(ruleSEMICOLON, position385)
					}
					if !rules[rulepropertyListPath]() {
						goto l383
					}
					goto l384
				l383:
					position, tokenIndex, depth = position383, tokenIndex383, depth383
				}
			l384:
				depth--
				add(rulepropertyListPath, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
			return false
		},
		/* 30 verbPath <- <path> */
		nil,
		/* 31 path <- <pathAlternative> */
		func() bool {
			position387, tokenIndex387, depth387 := position, tokenIndex, depth
			{
				position388 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l387
				}
				depth--
				add(rulepath, position388)
			}
			return true
		l387:
			position, tokenIndex, depth = position387, tokenIndex387, depth387
			return false
		},
		/* 32 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position389, tokenIndex389, depth389 := position, tokenIndex, depth
			{
				position390 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l389
				}
			l391:
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l392
					}
					if !rules[rulepathAlternative]() {
						goto l392
					}
					goto l391
				l392:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
				}
				depth--
				add(rulepathAlternative, position390)
			}
			return true
		l389:
			position, tokenIndex, depth = position389, tokenIndex389, depth389
			return false
		},
		/* 33 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position393, tokenIndex393, depth393 := position, tokenIndex, depth
			{
				position394 := position
				depth++
				{
					position395 := position
					depth++
					{
						position396 := position
						depth++
						{
							position397, tokenIndex397, depth397 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l397
							}
							goto l398
						l397:
							position, tokenIndex, depth = position397, tokenIndex397, depth397
						}
					l398:
						{
							position399 := position
							depth++
							{
								position400, tokenIndex400, depth400 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l401
								}
								goto l400
							l401:
								position, tokenIndex, depth = position400, tokenIndex400, depth400
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l393
										}
										if !rules[rulepath]() {
											goto l393
										}
										if !rules[ruleRPAREN]() {
											goto l393
										}
										break
									case '!':
										{
											position403 := position
											depth++
											if buffer[position] != rune('!') {
												goto l393
											}
											position++
											if !rules[rulews]() {
												goto l393
											}
											depth--
											add(ruleNOT, position403)
										}
										{
											position404 := position
											depth++
											{
												position405, tokenIndex405, depth405 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l406
												}
												goto l405
											l406:
												position, tokenIndex, depth = position405, tokenIndex405, depth405
												if !rules[ruleLPAREN]() {
													goto l393
												}
												{
													position407, tokenIndex407, depth407 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l407
													}
												l409:
													{
														position410, tokenIndex410, depth410 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l410
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l410
														}
														goto l409
													l410:
														position, tokenIndex, depth = position410, tokenIndex410, depth410
													}
													goto l408
												l407:
													position, tokenIndex, depth = position407, tokenIndex407, depth407
												}
											l408:
												if !rules[ruleRPAREN]() {
													goto l393
												}
											}
										l405:
											depth--
											add(rulepathNegatedPropertySet, position404)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l393
										}
										break
									}
								}

							}
						l400:
							depth--
							add(rulepathPrimary, position399)
						}
						depth--
						add(rulepathElt, position396)
					}
					depth--
					add(rulePegText, position395)
				}
				{
					add(ruleAction6, position)
				}
			l412:
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					{
						position414 := position
						depth++
						if buffer[position] != rune('/') {
							goto l413
						}
						position++
						if !rules[rulews]() {
							goto l413
						}
						depth--
						add(ruleSLASH, position414)
					}
					if !rules[rulepathSequence]() {
						goto l413
					}
					goto l412
				l413:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
				}
				depth--
				add(rulepathSequence, position394)
			}
			return true
		l393:
			position, tokenIndex, depth = position393, tokenIndex393, depth393
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
			position418, tokenIndex418, depth418 := position, tokenIndex, depth
			{
				position419 := position
				depth++
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l421
					}
					goto l420
				l421:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if !rules[ruleISA]() {
						goto l422
					}
					goto l420
				l422:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
					if !rules[ruleINVERSE]() {
						goto l418
					}
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l424
						}
						goto l423
					l424:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
						if !rules[ruleISA]() {
							goto l418
						}
					}
				l423:
				}
			l420:
				depth--
				add(rulepathOneInPropertySet, position419)
			}
			return true
		l418:
			position, tokenIndex, depth = position418, tokenIndex418, depth418
			return false
		},
		/* 38 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			{
				position426 := position
				depth++
				{
					position427 := position
					depth++
					{
						position428, tokenIndex428, depth428 := position, tokenIndex, depth
						if !rules[rulepof]() {
							goto l429
						}
						{
							add(ruleAction7, position)
						}
						goto l428
					l429:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
						{
							position432 := position
							depth++
							if !rules[rulegraphNodePath]() {
								goto l431
							}
							depth--
							add(rulePegText, position432)
						}
						{
							add(ruleAction8, position)
						}
						goto l428
					l431:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
						{
							add(ruleAction9, position)
						}
					}
				l428:
					depth--
					add(ruleobjectPath, position427)
				}
			l435:
				{
					position436, tokenIndex436, depth436 := position, tokenIndex, depth
					{
						position437 := position
						depth++
						if buffer[position] != rune(',') {
							goto l436
						}
						position++
						if !rules[rulews]() {
							goto l436
						}
						depth--
						add(ruleCOMMA, position437)
					}
					if !rules[ruleobjectListPath]() {
						goto l436
					}
					goto l435
				l436:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
				}
				depth--
				add(ruleobjectListPath, position426)
			}
			return true
		},
		/* 39 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		nil,
		/* 40 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			position439, tokenIndex439, depth439 := position, tokenIndex, depth
			{
				position440 := position
				depth++
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l442
					}
					goto l441
				l442:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
					if !rules[rulegraphTerm]() {
						goto l439
					}
				}
			l441:
				depth--
				add(rulegraphNodePath, position440)
			}
			return true
		l439:
			position, tokenIndex, depth = position439, tokenIndex439, depth439
			return false
		},
		/* 41 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position444 := position
				depth++
				{
					position445, tokenIndex445, depth445 := position, tokenIndex, depth
					{
						position447 := position
						depth++
						{
							position448, tokenIndex448, depth448 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l449
							}
							{
								position450, tokenIndex450, depth450 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l450
								}
								goto l451
							l450:
								position, tokenIndex, depth = position450, tokenIndex450, depth450
							}
						l451:
							goto l448
						l449:
							position, tokenIndex, depth = position448, tokenIndex448, depth448
							if !rules[ruleoffset]() {
								goto l445
							}
							{
								position452, tokenIndex452, depth452 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l452
								}
								goto l453
							l452:
								position, tokenIndex, depth = position452, tokenIndex452, depth452
							}
						l453:
						}
					l448:
						depth--
						add(rulelimitOffsetClauses, position447)
					}
					goto l446
				l445:
					position, tokenIndex, depth = position445, tokenIndex445, depth445
				}
			l446:
				depth--
				add(rulesolutionModifier, position444)
			}
			return true
		},
		/* 42 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 43 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				{
					position457 := position
					depth++
					{
						position458, tokenIndex458, depth458 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l459
						}
						position++
						goto l458
					l459:
						position, tokenIndex, depth = position458, tokenIndex458, depth458
						if buffer[position] != rune('L') {
							goto l455
						}
						position++
					}
				l458:
					{
						position460, tokenIndex460, depth460 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l461
						}
						position++
						goto l460
					l461:
						position, tokenIndex, depth = position460, tokenIndex460, depth460
						if buffer[position] != rune('I') {
							goto l455
						}
						position++
					}
				l460:
					{
						position462, tokenIndex462, depth462 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l463
						}
						position++
						goto l462
					l463:
						position, tokenIndex, depth = position462, tokenIndex462, depth462
						if buffer[position] != rune('M') {
							goto l455
						}
						position++
					}
				l462:
					{
						position464, tokenIndex464, depth464 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l465
						}
						position++
						goto l464
					l465:
						position, tokenIndex, depth = position464, tokenIndex464, depth464
						if buffer[position] != rune('I') {
							goto l455
						}
						position++
					}
				l464:
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l467
						}
						position++
						goto l466
					l467:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
						if buffer[position] != rune('T') {
							goto l455
						}
						position++
					}
				l466:
					if !rules[rulews]() {
						goto l455
					}
					depth--
					add(ruleLIMIT, position457)
				}
				if !rules[ruleINTEGER]() {
					goto l455
				}
				depth--
				add(rulelimit, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
			return false
		},
		/* 44 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position468, tokenIndex468, depth468 := position, tokenIndex, depth
			{
				position469 := position
				depth++
				{
					position470 := position
					depth++
					{
						position471, tokenIndex471, depth471 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l472
						}
						position++
						goto l471
					l472:
						position, tokenIndex, depth = position471, tokenIndex471, depth471
						if buffer[position] != rune('O') {
							goto l468
						}
						position++
					}
				l471:
					{
						position473, tokenIndex473, depth473 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l474
						}
						position++
						goto l473
					l474:
						position, tokenIndex, depth = position473, tokenIndex473, depth473
						if buffer[position] != rune('F') {
							goto l468
						}
						position++
					}
				l473:
					{
						position475, tokenIndex475, depth475 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l476
						}
						position++
						goto l475
					l476:
						position, tokenIndex, depth = position475, tokenIndex475, depth475
						if buffer[position] != rune('F') {
							goto l468
						}
						position++
					}
				l475:
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l478
						}
						position++
						goto l477
					l478:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
						if buffer[position] != rune('S') {
							goto l468
						}
						position++
					}
				l477:
					{
						position479, tokenIndex479, depth479 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l480
						}
						position++
						goto l479
					l480:
						position, tokenIndex, depth = position479, tokenIndex479, depth479
						if buffer[position] != rune('E') {
							goto l468
						}
						position++
					}
				l479:
					{
						position481, tokenIndex481, depth481 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l482
						}
						position++
						goto l481
					l482:
						position, tokenIndex, depth = position481, tokenIndex481, depth481
						if buffer[position] != rune('T') {
							goto l468
						}
						position++
					}
				l481:
					if !rules[rulews]() {
						goto l468
					}
					depth--
					add(ruleOFFSET, position470)
				}
				if !rules[ruleINTEGER]() {
					goto l468
				}
				depth--
				add(ruleoffset, position469)
			}
			return true
		l468:
			position, tokenIndex, depth = position468, tokenIndex468, depth468
			return false
		},
		/* 45 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' whiteSpaces+)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					{
						position487 := position
						depth++
					l488:
						{
							position489, tokenIndex489, depth489 := position, tokenIndex, depth
							{
								position490, tokenIndex490, depth490 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l491
								}
								position++
								goto l490
							l491:
								position, tokenIndex, depth = position490, tokenIndex490, depth490
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l489
								}
								position++
							}
						l490:
							goto l488
						l489:
							position, tokenIndex, depth = position489, tokenIndex489, depth489
						}
						depth--
						add(rulePegText, position487)
					}
					if buffer[position] != rune(':') {
						goto l486
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l485
				l486:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					{
						position494 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l493
						}
						position++
					l495:
						{
							position496, tokenIndex496, depth496 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l496
							}
							position++
							goto l495
						l496:
							position, tokenIndex, depth = position496, tokenIndex496, depth496
						}
						depth--
						add(rulePegText, position494)
					}
					if buffer[position] != rune('/') {
						goto l493
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l485
				l493:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
					{
						position498 := position
						depth++
					l499:
						{
							position500, tokenIndex500, depth500 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l500
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l500
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l500
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l500
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l500
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l500
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l500
									}
									position++
									break
								}
							}

							goto l499
						l500:
							position, tokenIndex, depth = position500, tokenIndex500, depth500
						}
						depth--
						add(rulePegText, position498)
					}
					{
						add(ruleAction12, position)
					}
				}
			l485:
				if buffer[position] != rune('<') {
					goto l483
				}
				position++
				if !rules[rulewhiteSpaces]() {
					goto l483
				}
			l503:
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					if !rules[rulewhiteSpaces]() {
						goto l504
					}
					goto l503
				l504:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
				}
				depth--
				add(rulepof, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 46 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position505, tokenIndex505, depth505 := position, tokenIndex, depth
			{
				position506 := position
				depth++
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l508
					}
					position++
					goto l507
				l508:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
					if buffer[position] != rune('$') {
						goto l505
					}
					position++
				}
			l507:
				{
					position509 := position
					depth++
					{
						position512, tokenIndex512, depth512 := position, tokenIndex, depth
						{
							position514 := position
							depth++
							{
								position515, tokenIndex515, depth515 := position, tokenIndex, depth
								{
									position517 := position
									depth++
									{
										position518, tokenIndex518, depth518 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l519
										}
										position++
										goto l518
									l519:
										position, tokenIndex, depth = position518, tokenIndex518, depth518
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l516
										}
										position++
									}
								l518:
									depth--
									add(rulePN_CHARS_BASE, position517)
								}
								goto l515
							l516:
								position, tokenIndex, depth = position515, tokenIndex515, depth515
								if buffer[position] != rune('_') {
									goto l513
								}
								position++
							}
						l515:
							depth--
							add(rulePN_CHARS_U, position514)
						}
						goto l512
					l513:
						position, tokenIndex, depth = position512, tokenIndex512, depth512
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l505
						}
						position++
					}
				l512:
				l510:
					{
						position511, tokenIndex511, depth511 := position, tokenIndex, depth
						{
							position520, tokenIndex520, depth520 := position, tokenIndex, depth
							{
								position522 := position
								depth++
								{
									position523, tokenIndex523, depth523 := position, tokenIndex, depth
									{
										position525 := position
										depth++
										{
											position526, tokenIndex526, depth526 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l527
											}
											position++
											goto l526
										l527:
											position, tokenIndex, depth = position526, tokenIndex526, depth526
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l524
											}
											position++
										}
									l526:
										depth--
										add(rulePN_CHARS_BASE, position525)
									}
									goto l523
								l524:
									position, tokenIndex, depth = position523, tokenIndex523, depth523
									if buffer[position] != rune('_') {
										goto l521
									}
									position++
								}
							l523:
								depth--
								add(rulePN_CHARS_U, position522)
							}
							goto l520
						l521:
							position, tokenIndex, depth = position520, tokenIndex520, depth520
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l511
							}
							position++
						}
					l520:
						goto l510
					l511:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
					}
					depth--
					add(ruleVARNAME, position509)
				}
				if !rules[rulews]() {
					goto l505
				}
				depth--
				add(rulevar, position506)
			}
			return true
		l505:
			position, tokenIndex, depth = position505, tokenIndex505, depth505
			return false
		},
		/* 47 iriref <- <(iri / prefixedName)> */
		func() bool {
			position528, tokenIndex528, depth528 := position, tokenIndex, depth
			{
				position529 := position
				depth++
				{
					position530, tokenIndex530, depth530 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l531
					}
					goto l530
				l531:
					position, tokenIndex, depth = position530, tokenIndex530, depth530
					{
						position532 := position
						depth++
					l533:
						{
							position534, tokenIndex534, depth534 := position, tokenIndex, depth
							{
								position535, tokenIndex535, depth535 := position, tokenIndex, depth
								{
									position536, tokenIndex536, depth536 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l537
									}
									position++
									goto l536
								l537:
									position, tokenIndex, depth = position536, tokenIndex536, depth536
									if buffer[position] != rune(' ') {
										goto l535
									}
									position++
								}
							l536:
								goto l534
							l535:
								position, tokenIndex, depth = position535, tokenIndex535, depth535
							}
							if !matchDot() {
								goto l534
							}
							goto l533
						l534:
							position, tokenIndex, depth = position534, tokenIndex534, depth534
						}
						if buffer[position] != rune(':') {
							goto l528
						}
						position++
					l538:
						{
							position539, tokenIndex539, depth539 := position, tokenIndex, depth
							{
								position540, tokenIndex540, depth540 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l541
								}
								position++
								goto l540
							l541:
								position, tokenIndex, depth = position540, tokenIndex540, depth540
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l542
								}
								position++
								goto l540
							l542:
								position, tokenIndex, depth = position540, tokenIndex540, depth540
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l543
								}
								position++
								goto l540
							l543:
								position, tokenIndex, depth = position540, tokenIndex540, depth540
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l539
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l539
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l539
										}
										position++
										break
									}
								}

							}
						l540:
							goto l538
						l539:
							position, tokenIndex, depth = position539, tokenIndex539, depth539
						}
						if !rules[rulews]() {
							goto l528
						}
						depth--
						add(ruleprefixedName, position532)
					}
				}
			l530:
				depth--
				add(ruleiriref, position529)
			}
			return true
		l528:
			position, tokenIndex, depth = position528, tokenIndex528, depth528
			return false
		},
		/* 48 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position545, tokenIndex545, depth545 := position, tokenIndex, depth
			{
				position546 := position
				depth++
				if buffer[position] != rune('<') {
					goto l545
				}
				position++
			l547:
				{
					position548, tokenIndex548, depth548 := position, tokenIndex, depth
					{
						position549, tokenIndex549, depth549 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l549
						}
						position++
						goto l548
					l549:
						position, tokenIndex, depth = position549, tokenIndex549, depth549
					}
					if !matchDot() {
						goto l548
					}
					goto l547
				l548:
					position, tokenIndex, depth = position548, tokenIndex548, depth548
				}
				if buffer[position] != rune('>') {
					goto l545
				}
				position++
				if !rules[rulews]() {
					goto l545
				}
				depth--
				add(ruleiri, position546)
			}
			return true
		l545:
			position, tokenIndex, depth = position545, tokenIndex545, depth545
			return false
		},
		/* 49 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* ws)> */
		nil,
		/* 50 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? ws)> */
		nil,
		/* 51 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 52 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? ws)> */
		nil,
		/* 53 booleanLiteral <- <(TRUE / FALSE)> */
		nil,
		/* 54 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 55 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? ws)> */
		nil,
		/* 56 anon <- <('[' ws ']' ws)> */
		nil,
		/* 57 nil <- <('(' ws ')' ws)> */
		nil,
		/* 58 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 59 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 60 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 61 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') ws)> */
		nil,
		/* 62 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') ws)> */
		nil,
		/* 63 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 64 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 65 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 66 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 67 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 68 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') ws)> */
		nil,
		/* 69 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 70 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') ws)> */
		nil,
		/* 71 LBRACE <- <('{' ws)> */
		func() bool {
			position572, tokenIndex572, depth572 := position, tokenIndex, depth
			{
				position573 := position
				depth++
				if buffer[position] != rune('{') {
					goto l572
				}
				position++
				if !rules[rulews]() {
					goto l572
				}
				depth--
				add(ruleLBRACE, position573)
			}
			return true
		l572:
			position, tokenIndex, depth = position572, tokenIndex572, depth572
			return false
		},
		/* 72 RBRACE <- <('}' ws)> */
		func() bool {
			position574, tokenIndex574, depth574 := position, tokenIndex, depth
			{
				position575 := position
				depth++
				if buffer[position] != rune('}') {
					goto l574
				}
				position++
				if !rules[rulews]() {
					goto l574
				}
				depth--
				add(ruleRBRACE, position575)
			}
			return true
		l574:
			position, tokenIndex, depth = position574, tokenIndex574, depth574
			return false
		},
		/* 73 LBRACK <- <('[' ws)> */
		nil,
		/* 74 RBRACK <- <(']' ws)> */
		nil,
		/* 75 SEMICOLON <- <(';' ws)> */
		nil,
		/* 76 COMMA <- <(',' ws)> */
		nil,
		/* 77 DOT <- <('.' ws)> */
		func() bool {
			position580, tokenIndex580, depth580 := position, tokenIndex, depth
			{
				position581 := position
				depth++
				if buffer[position] != rune('.') {
					goto l580
				}
				position++
				if !rules[rulews]() {
					goto l580
				}
				depth--
				add(ruleDOT, position581)
			}
			return true
		l580:
			position, tokenIndex, depth = position580, tokenIndex580, depth580
			return false
		},
		/* 78 COLON <- <(':' ws)> */
		nil,
		/* 79 PIPE <- <('|' ws)> */
		func() bool {
			position583, tokenIndex583, depth583 := position, tokenIndex, depth
			{
				position584 := position
				depth++
				if buffer[position] != rune('|') {
					goto l583
				}
				position++
				if !rules[rulews]() {
					goto l583
				}
				depth--
				add(rulePIPE, position584)
			}
			return true
		l583:
			position, tokenIndex, depth = position583, tokenIndex583, depth583
			return false
		},
		/* 80 SLASH <- <('/' ws)> */
		nil,
		/* 81 INVERSE <- <('^' ws)> */
		func() bool {
			position586, tokenIndex586, depth586 := position, tokenIndex, depth
			{
				position587 := position
				depth++
				if buffer[position] != rune('^') {
					goto l586
				}
				position++
				if !rules[rulews]() {
					goto l586
				}
				depth--
				add(ruleINVERSE, position587)
			}
			return true
		l586:
			position, tokenIndex, depth = position586, tokenIndex586, depth586
			return false
		},
		/* 82 LPAREN <- <('(' ws)> */
		func() bool {
			position588, tokenIndex588, depth588 := position, tokenIndex, depth
			{
				position589 := position
				depth++
				if buffer[position] != rune('(') {
					goto l588
				}
				position++
				if !rules[rulews]() {
					goto l588
				}
				depth--
				add(ruleLPAREN, position589)
			}
			return true
		l588:
			position, tokenIndex, depth = position588, tokenIndex588, depth588
			return false
		},
		/* 83 RPAREN <- <(')' ws)> */
		func() bool {
			position590, tokenIndex590, depth590 := position, tokenIndex, depth
			{
				position591 := position
				depth++
				if buffer[position] != rune(')') {
					goto l590
				}
				position++
				if !rules[rulews]() {
					goto l590
				}
				depth--
				add(ruleRPAREN, position591)
			}
			return true
		l590:
			position, tokenIndex, depth = position590, tokenIndex590, depth590
			return false
		},
		/* 84 ISA <- <('a' ws)> */
		func() bool {
			position592, tokenIndex592, depth592 := position, tokenIndex, depth
			{
				position593 := position
				depth++
				if buffer[position] != rune('a') {
					goto l592
				}
				position++
				if !rules[rulews]() {
					goto l592
				}
				depth--
				add(ruleISA, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 85 NOT <- <('!' ws)> */
		nil,
		/* 86 STAR <- <('*' ws)> */
		func() bool {
			position595, tokenIndex595, depth595 := position, tokenIndex, depth
			{
				position596 := position
				depth++
				if buffer[position] != rune('*') {
					goto l595
				}
				position++
				if !rules[rulews]() {
					goto l595
				}
				depth--
				add(ruleSTAR, position596)
			}
			return true
		l595:
			position, tokenIndex, depth = position595, tokenIndex595, depth595
			return false
		},
		/* 87 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') ws)> */
		nil,
		/* 88 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') ws)> */
		nil,
		/* 89 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') ws)> */
		nil,
		/* 90 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') ws)> */
		nil,
		/* 91 INTEGER <- <([0-9]+ ws)> */
		func() bool {
			position601, tokenIndex601, depth601 := position, tokenIndex, depth
			{
				position602 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l601
				}
				position++
			l603:
				{
					position604, tokenIndex604, depth604 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l604
					}
					position++
					goto l603
				l604:
					position, tokenIndex, depth = position604, tokenIndex604, depth604
				}
				if !rules[rulews]() {
					goto l601
				}
				depth--
				add(ruleINTEGER, position602)
			}
			return true
		l601:
			position, tokenIndex, depth = position601, tokenIndex601, depth601
			return false
		},
		/* 92 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 93 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') ws)> */
		nil,
		/* 94 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') ws)> */
		nil,
		/* 95 whiteSpaces <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l608
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l608
						}
						position++
						break
					case '\n':
						if buffer[position] != rune('\n') {
							goto l608
						}
						position++
						break
					case '\r':
						if buffer[position] != rune('\r') {
							goto l608
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l608
						}
						position++
						break
					default:
						if buffer[position] != rune(' ') {
							goto l608
						}
						position++
						break
					}
				}

				depth--
				add(rulewhiteSpaces, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 96 ws <- <whiteSpaces*> */
		func() bool {
			{
				position612 := position
				depth++
			l613:
				{
					position614, tokenIndex614, depth614 := position, tokenIndex, depth
					if !rules[rulewhiteSpaces]() {
						goto l614
					}
					goto l613
				l614:
					position, tokenIndex, depth = position614, tokenIndex614, depth614
				}
				depth--
				add(rulews, position612)
			}
			return true
		},
		nil,
		/* 99 Action0 <- <{ p.addPrefix(buffer[begin:end]) }> */
		nil,
		/* 100 Action1 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 101 Action2 <- <{ p.setSubject(buffer[begin:end]) }> */
		nil,
		/* 102 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 103 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 104 Action5 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 105 Action6 <- <{ p.setPredicate(buffer[begin:end]) }> */
		nil,
		/* 106 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 107 Action8 <- <{ p.setObject(buffer[begin:end]); p.addTriplePattern() }> */
		nil,
		/* 108 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 109 Action10 <- <{ p.setPrefix(buffer[begin:end]) }> */
		nil,
		/* 110 Action11 <- <{ p.setPathLength(buffer[begin:end]) }> */
		nil,
		/* 111 Action12 <- <{ p.setKeyword(buffer[begin:end]) }> */
		nil,
	}
	p.rules = rules
}
