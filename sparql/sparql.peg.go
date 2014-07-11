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
	rulepathMod
	ruleobjectListPath
	ruleobjectPath
	rulegraphNodePath
	rulesolutionModifier
	rulelimitOffsetClauses
	rulelimit
	ruleoffset
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
	rules  [90]func() bool
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
					if !rules[rulevarOrTerm]() {
						goto l222
					}
					if !rules[rulepropertyListPath]() {
						goto l222
					}
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					{
						position223 := position
						depth++
						{
							position224, tokenIndex224, depth224 := position, tokenIndex, depth
							{
								position226 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l225
								}
								if !rules[rulegraphNodePath]() {
									goto l225
								}
							l227:
								{
									position228, tokenIndex228, depth228 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l228
									}
									goto l227
								l228:
									position, tokenIndex, depth = position228, tokenIndex228, depth228
								}
								if !rules[ruleRPAREN]() {
									goto l225
								}
								depth--
								add(rulecollectionPath, position226)
							}
							goto l224
						l225:
							position, tokenIndex, depth = position224, tokenIndex224, depth224
							{
								position229 := position
								depth++
								{
									position230 := position
									depth++
									if buffer[position] != rune('[') {
										goto l219
									}
									position++
									if !rules[rulews]() {
										goto l219
									}
									depth--
									add(ruleLBRACK, position230)
								}
								if !rules[rulepropertyListPath]() {
									goto l219
								}
								{
									position231 := position
									depth++
									if buffer[position] != rune(']') {
										goto l219
									}
									position++
									if !rules[rulews]() {
										goto l219
									}
									depth--
									add(ruleRBRACK, position231)
								}
								depth--
								add(ruleblankNodePropertyListPath, position229)
							}
						}
					l224:
						depth--
						add(ruletriplesNodePath, position223)
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
		/* 19 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position232, tokenIndex232, depth232 := position, tokenIndex, depth
			{
				position233 := position
				depth++
				{
					position234, tokenIndex234, depth234 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l235
					}
					goto l234
				l235:
					position, tokenIndex, depth = position234, tokenIndex234, depth234
					{
						position236 := position
						depth++
						{
							switch buffer[position] {
							case '(':
								{
									position238 := position
									depth++
									if buffer[position] != rune('(') {
										goto l232
									}
									position++
									if !rules[rulews]() {
										goto l232
									}
									if buffer[position] != rune(')') {
										goto l232
									}
									position++
									if !rules[rulews]() {
										goto l232
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
												goto l232
											}
											position++
											if !rules[rulews]() {
												goto l232
											}
											if buffer[position] != rune(']') {
												goto l232
											}
											position++
											if !rules[rulews]() {
												goto l232
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
													goto l232
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
													goto l232
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
													goto l232
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
													goto l232
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
													goto l232
												}
												position++
											}
										l272:
											if !rules[rulews]() {
												goto l232
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
											goto l232
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
											goto l232
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
											if !rules[ruleiri]() {
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
										goto l232
									}
									depth--
									add(ruleliteral, position274)
								}
								break
							case '<':
								if !rules[ruleiri]() {
									goto l232
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
										goto l232
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
										goto l232
									}
									depth--
									add(rulenumericLiteral, position295)
								}
								break
							}
						}

						depth--
						add(rulegraphTerm, position236)
					}
				}
			l234:
				depth--
				add(rulevarOrTerm, position233)
			}
			return true
		l232:
			position, tokenIndex, depth = position232, tokenIndex232, depth232
			return false
		},
		/* 20 graphTerm <- <((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('<') iri) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral))> */
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
									position332 := position
									depth++
									if buffer[position] != rune('!') {
										goto l325
									}
									position++
									if !rules[rulews]() {
										goto l325
									}
									depth--
									add(ruleNOT, position332)
								}
								{
									position333 := position
									depth++
									{
										position334, tokenIndex334, depth334 := position, tokenIndex, depth
										if !rules[rulepathOneInPropertySet]() {
											goto l335
										}
										goto l334
									l335:
										position, tokenIndex, depth = position334, tokenIndex334, depth334
										if !rules[ruleLPAREN]() {
											goto l325
										}
										{
											position336, tokenIndex336, depth336 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l336
											}
										l338:
											{
												position339, tokenIndex339, depth339 := position, tokenIndex, depth
												if !rules[rulePIPE]() {
													goto l339
												}
												if !rules[rulepathOneInPropertySet]() {
													goto l339
												}
												goto l338
											l339:
												position, tokenIndex, depth = position339, tokenIndex339, depth339
											}
											goto l337
										l336:
											position, tokenIndex, depth = position336, tokenIndex336, depth336
										}
									l337:
										if !rules[ruleRPAREN]() {
											goto l325
										}
									}
								l334:
									depth--
									add(rulepathNegatedPropertySet, position333)
								}
								break
							case 'a':
								if !rules[ruleISA]() {
									goto l325
								}
								break
							default:
								if !rules[ruleiri]() {
									goto l325
								}
								break
							}
						}

						depth--
						add(rulepathPrimary, position330)
					}
					{
						position340, tokenIndex340, depth340 := position, tokenIndex, depth
						{
							position341, tokenIndex341, depth341 := position, tokenIndex, depth
							{
								position343 := position
								depth++
								{
									switch buffer[position] {
									case '+':
										{
											position345 := position
											depth++
											if buffer[position] != rune('+') {
												goto l341
											}
											position++
											if !rules[rulews]() {
												goto l341
											}
											depth--
											add(rulePLUS, position345)
										}
										break
									case '?':
										{
											position346 := position
											depth++
											if buffer[position] != rune('?') {
												goto l341
											}
											position++
											if !rules[rulews]() {
												goto l341
											}
											depth--
											add(ruleQUESTION, position346)
										}
										break
									default:
										if !rules[ruleSTAR]() {
											goto l341
										}
										break
									}
								}

								depth--
								add(rulepathMod, position343)
							}
							goto l342
						l341:
							position, tokenIndex, depth = position341, tokenIndex341, depth341
						}
					l342:
						position, tokenIndex, depth = position340, tokenIndex340, depth340
					}
					depth--
					add(rulepathElt, position327)
				}
			l347:
				{
					position348, tokenIndex348, depth348 := position, tokenIndex, depth
					{
						position349 := position
						depth++
						if buffer[position] != rune('/') {
							goto l348
						}
						position++
						if !rules[rulews]() {
							goto l348
						}
						depth--
						add(ruleSLASH, position349)
					}
					if !rules[rulepathSequence]() {
						goto l348
					}
					goto l347
				l348:
					position, tokenIndex, depth = position348, tokenIndex348, depth348
				}
				depth--
				add(rulepathSequence, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
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
			position353, tokenIndex353, depth353 := position, tokenIndex, depth
			{
				position354 := position
				depth++
				{
					switch buffer[position] {
					case '^':
						if !rules[ruleINVERSE]() {
							goto l353
						}
						{
							position356, tokenIndex356, depth356 := position, tokenIndex, depth
							if !rules[ruleiri]() {
								goto l357
							}
							goto l356
						l357:
							position, tokenIndex, depth = position356, tokenIndex356, depth356
							if !rules[ruleISA]() {
								goto l353
							}
						}
					l356:
						break
					case 'a':
						if !rules[ruleISA]() {
							goto l353
						}
						break
					default:
						if !rules[ruleiri]() {
							goto l353
						}
						break
					}
				}

				depth--
				add(rulepathOneInPropertySet, position354)
			}
			return true
		l353:
			position, tokenIndex, depth = position353, tokenIndex353, depth353
			return false
		},
		/* 33 pathMod <- <((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR))> */
		nil,
		/* 34 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position359, tokenIndex359, depth359 := position, tokenIndex, depth
			{
				position360 := position
				depth++
				{
					position361 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l359
					}
					depth--
					add(ruleobjectPath, position361)
				}
			l362:
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					{
						position364 := position
						depth++
						if buffer[position] != rune(',') {
							goto l363
						}
						position++
						if !rules[rulews]() {
							goto l363
						}
						depth--
						add(ruleCOMMA, position364)
					}
					if !rules[ruleobjectListPath]() {
						goto l363
					}
					goto l362
				l363:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
				}
				depth--
				add(ruleobjectListPath, position360)
			}
			return true
		l359:
			position, tokenIndex, depth = position359, tokenIndex359, depth359
			return false
		},
		/* 35 objectPath <- <graphNodePath> */
		nil,
		/* 36 graphNodePath <- <varOrTerm> */
		func() bool {
			position366, tokenIndex366, depth366 := position, tokenIndex, depth
			{
				position367 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l366
				}
				depth--
				add(rulegraphNodePath, position367)
			}
			return true
		l366:
			position, tokenIndex, depth = position366, tokenIndex366, depth366
			return false
		},
		/* 37 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 38 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 39 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position370, tokenIndex370, depth370 := position, tokenIndex, depth
			{
				position371 := position
				depth++
				{
					position372 := position
					depth++
					{
						position373, tokenIndex373, depth373 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l374
						}
						position++
						goto l373
					l374:
						position, tokenIndex, depth = position373, tokenIndex373, depth373
						if buffer[position] != rune('L') {
							goto l370
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
							goto l370
						}
						position++
					}
				l375:
					{
						position377, tokenIndex377, depth377 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l378
						}
						position++
						goto l377
					l378:
						position, tokenIndex, depth = position377, tokenIndex377, depth377
						if buffer[position] != rune('M') {
							goto l370
						}
						position++
					}
				l377:
					{
						position379, tokenIndex379, depth379 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l380
						}
						position++
						goto l379
					l380:
						position, tokenIndex, depth = position379, tokenIndex379, depth379
						if buffer[position] != rune('I') {
							goto l370
						}
						position++
					}
				l379:
					{
						position381, tokenIndex381, depth381 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l382
						}
						position++
						goto l381
					l382:
						position, tokenIndex, depth = position381, tokenIndex381, depth381
						if buffer[position] != rune('T') {
							goto l370
						}
						position++
					}
				l381:
					if !rules[rulews]() {
						goto l370
					}
					depth--
					add(ruleLIMIT, position372)
				}
				if !rules[ruleINTEGER]() {
					goto l370
				}
				depth--
				add(rulelimit, position371)
			}
			return true
		l370:
			position, tokenIndex, depth = position370, tokenIndex370, depth370
			return false
		},
		/* 40 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position383, tokenIndex383, depth383 := position, tokenIndex, depth
			{
				position384 := position
				depth++
				{
					position385 := position
					depth++
					{
						position386, tokenIndex386, depth386 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l387
						}
						position++
						goto l386
					l387:
						position, tokenIndex, depth = position386, tokenIndex386, depth386
						if buffer[position] != rune('O') {
							goto l383
						}
						position++
					}
				l386:
					{
						position388, tokenIndex388, depth388 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l389
						}
						position++
						goto l388
					l389:
						position, tokenIndex, depth = position388, tokenIndex388, depth388
						if buffer[position] != rune('F') {
							goto l383
						}
						position++
					}
				l388:
					{
						position390, tokenIndex390, depth390 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l391
						}
						position++
						goto l390
					l391:
						position, tokenIndex, depth = position390, tokenIndex390, depth390
						if buffer[position] != rune('F') {
							goto l383
						}
						position++
					}
				l390:
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l393
						}
						position++
						goto l392
					l393:
						position, tokenIndex, depth = position392, tokenIndex392, depth392
						if buffer[position] != rune('S') {
							goto l383
						}
						position++
					}
				l392:
					{
						position394, tokenIndex394, depth394 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l395
						}
						position++
						goto l394
					l395:
						position, tokenIndex, depth = position394, tokenIndex394, depth394
						if buffer[position] != rune('E') {
							goto l383
						}
						position++
					}
				l394:
					{
						position396, tokenIndex396, depth396 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l397
						}
						position++
						goto l396
					l397:
						position, tokenIndex, depth = position396, tokenIndex396, depth396
						if buffer[position] != rune('T') {
							goto l383
						}
						position++
					}
				l396:
					if !rules[rulews]() {
						goto l383
					}
					depth--
					add(ruleOFFSET, position385)
				}
				if !rules[ruleINTEGER]() {
					goto l383
				}
				depth--
				add(ruleoffset, position384)
			}
			return true
		l383:
			position, tokenIndex, depth = position383, tokenIndex383, depth383
			return false
		},
		/* 41 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position398, tokenIndex398, depth398 := position, tokenIndex, depth
			{
				position399 := position
				depth++
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l401
					}
					position++
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					if buffer[position] != rune('$') {
						goto l398
					}
					position++
				}
			l400:
				{
					position402 := position
					depth++
					{
						position403, tokenIndex403, depth403 := position, tokenIndex, depth
						if !rules[rulePN_CHARS_U]() {
							goto l404
						}
						goto l403
					l404:
						position, tokenIndex, depth = position403, tokenIndex403, depth403
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l398
						}
						position++
					}
				l403:
				l405:
					{
						position406, tokenIndex406, depth406 := position, tokenIndex, depth
						{
							position407 := position
							depth++
							{
								position408, tokenIndex408, depth408 := position, tokenIndex, depth
								if !rules[rulePN_CHARS_U]() {
									goto l409
								}
								goto l408
							l409:
								position, tokenIndex, depth = position408, tokenIndex408, depth408
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l406
								}
								position++
							}
						l408:
							depth--
							add(ruleVAR_CHAR, position407)
						}
						goto l405
					l406:
						position, tokenIndex, depth = position406, tokenIndex406, depth406
					}
					depth--
					add(ruleVARNAME, position402)
				}
				if !rules[rulews]() {
					goto l398
				}
				depth--
				add(rulevar, position399)
			}
			return true
		l398:
			position, tokenIndex, depth = position398, tokenIndex398, depth398
			return false
		},
		/* 42 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position410, tokenIndex410, depth410 := position, tokenIndex, depth
			{
				position411 := position
				depth++
				if buffer[position] != rune('<') {
					goto l410
				}
				position++
			l412:
				{
					position413, tokenIndex413, depth413 := position, tokenIndex, depth
					{
						position414, tokenIndex414, depth414 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l414
						}
						position++
						goto l413
					l414:
						position, tokenIndex, depth = position414, tokenIndex414, depth414
					}
					if !matchDot() {
						goto l413
					}
					goto l412
				l413:
					position, tokenIndex, depth = position413, tokenIndex413, depth413
				}
				if buffer[position] != rune('>') {
					goto l410
				}
				position++
				if !rules[rulews]() {
					goto l410
				}
				depth--
				add(ruleiri, position411)
			}
			return true
		l410:
			position, tokenIndex, depth = position410, tokenIndex410, depth410
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
		/* 51 VARNAME <- <((PN_CHARS_U / [0-9]) VAR_CHAR*)> */
		nil,
		/* 52 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		func() bool {
			position424, tokenIndex424, depth424 := position, tokenIndex, depth
			{
				position425 := position
				depth++
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					{
						position428 := position
						depth++
						{
							position429, tokenIndex429, depth429 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l430
							}
							position++
							goto l429
						l430:
							position, tokenIndex, depth = position429, tokenIndex429, depth429
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l427
							}
							position++
						}
					l429:
						depth--
						add(rulePN_CHARS_BASE, position428)
					}
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if buffer[position] != rune('_') {
						goto l424
					}
					position++
				}
			l426:
				depth--
				add(rulePN_CHARS_U, position425)
			}
			return true
		l424:
			position, tokenIndex, depth = position424, tokenIndex424, depth424
			return false
		},
		/* 53 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 54 VAR_CHAR <- <(PN_CHARS_U / [0-9])> */
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
			position443, tokenIndex443, depth443 := position, tokenIndex, depth
			{
				position444 := position
				depth++
				if buffer[position] != rune('{') {
					goto l443
				}
				position++
				if !rules[rulews]() {
					goto l443
				}
				depth--
				add(ruleLBRACE, position444)
			}
			return true
		l443:
			position, tokenIndex, depth = position443, tokenIndex443, depth443
			return false
		},
		/* 66 RBRACE <- <('}' ws)> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				if buffer[position] != rune('}') {
					goto l445
				}
				position++
				if !rules[rulews]() {
					goto l445
				}
				depth--
				add(ruleRBRACE, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
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
			position451, tokenIndex451, depth451 := position, tokenIndex, depth
			{
				position452 := position
				depth++
				if buffer[position] != rune('.') {
					goto l451
				}
				position++
				if !rules[rulews]() {
					goto l451
				}
				depth--
				add(ruleDOT, position452)
			}
			return true
		l451:
			position, tokenIndex, depth = position451, tokenIndex451, depth451
			return false
		},
		/* 72 COLON <- <(':' ws)> */
		nil,
		/* 73 PIPE <- <('|' ws)> */
		func() bool {
			position454, tokenIndex454, depth454 := position, tokenIndex, depth
			{
				position455 := position
				depth++
				if buffer[position] != rune('|') {
					goto l454
				}
				position++
				if !rules[rulews]() {
					goto l454
				}
				depth--
				add(rulePIPE, position455)
			}
			return true
		l454:
			position, tokenIndex, depth = position454, tokenIndex454, depth454
			return false
		},
		/* 74 SLASH <- <('/' ws)> */
		nil,
		/* 75 INVERSE <- <('^' ws)> */
		func() bool {
			position457, tokenIndex457, depth457 := position, tokenIndex, depth
			{
				position458 := position
				depth++
				if buffer[position] != rune('^') {
					goto l457
				}
				position++
				if !rules[rulews]() {
					goto l457
				}
				depth--
				add(ruleINVERSE, position458)
			}
			return true
		l457:
			position, tokenIndex, depth = position457, tokenIndex457, depth457
			return false
		},
		/* 76 LPAREN <- <('(' ws)> */
		func() bool {
			position459, tokenIndex459, depth459 := position, tokenIndex, depth
			{
				position460 := position
				depth++
				if buffer[position] != rune('(') {
					goto l459
				}
				position++
				if !rules[rulews]() {
					goto l459
				}
				depth--
				add(ruleLPAREN, position460)
			}
			return true
		l459:
			position, tokenIndex, depth = position459, tokenIndex459, depth459
			return false
		},
		/* 77 RPAREN <- <(')' ws)> */
		func() bool {
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if buffer[position] != rune(')') {
					goto l461
				}
				position++
				if !rules[rulews]() {
					goto l461
				}
				depth--
				add(ruleRPAREN, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 78 ISA <- <('a' ws)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				if buffer[position] != rune('a') {
					goto l463
				}
				position++
				if !rules[rulews]() {
					goto l463
				}
				depth--
				add(ruleISA, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
			return false
		},
		/* 79 NOT <- <('!' ws)> */
		nil,
		/* 80 STAR <- <('*' ws)> */
		func() bool {
			position466, tokenIndex466, depth466 := position, tokenIndex, depth
			{
				position467 := position
				depth++
				if buffer[position] != rune('*') {
					goto l466
				}
				position++
				if !rules[rulews]() {
					goto l466
				}
				depth--
				add(ruleSTAR, position467)
			}
			return true
		l466:
			position, tokenIndex, depth = position466, tokenIndex466, depth466
			return false
		},
		/* 81 QUESTION <- <('?' ws)> */
		nil,
		/* 82 PLUS <- <('+' ws)> */
		nil,
		/* 83 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') ws)> */
		nil,
		/* 84 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') ws)> */
		nil,
		/* 85 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') ws)> */
		nil,
		/* 86 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') ws)> */
		nil,
		/* 87 INTEGER <- <([0-9]+ ws)> */
		func() bool {
			position474, tokenIndex474, depth474 := position, tokenIndex, depth
			{
				position475 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l474
				}
				position++
			l476:
				{
					position477, tokenIndex477, depth477 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l477
					}
					position++
					goto l476
				l477:
					position, tokenIndex, depth = position477, tokenIndex477, depth477
				}
				if !rules[rulews]() {
					goto l474
				}
				depth--
				add(ruleINTEGER, position475)
			}
			return true
		l474:
			position, tokenIndex, depth = position474, tokenIndex474, depth474
			return false
		},
		/* 88 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position479 := position
				depth++
			l480:
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\v':
							if buffer[position] != rune('\v') {
								goto l481
							}
							position++
							break
						case '\f':
							if buffer[position] != rune('\f') {
								goto l481
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l481
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l481
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l481
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l481
							}
							position++
							break
						}
					}

					goto l480
				l481:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
				}
				depth--
				add(rulews, position479)
			}
			return true
		},
	}
	p.rules = rules
}
