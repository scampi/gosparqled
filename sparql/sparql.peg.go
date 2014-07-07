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
	rulePegText

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
	"PegText",

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
	rules  [91]func() bool
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
									position21 := position
									depth++
									{
										position24, tokenIndex24, depth24 := position, tokenIndex, depth
										{
											position25, tokenIndex25, depth25 := position, tokenIndex, depth
											if buffer[position] != rune('\'') {
												goto l26
											}
											position++
											goto l25
										l26:
											position, tokenIndex, depth = position25, tokenIndex25, depth25
											{
												switch buffer[position] {
												case ' ':
													if buffer[position] != rune(' ') {
														goto l24
													}
													position++
													break
												case '\'':
													if buffer[position] != rune('\'') {
														goto l24
													}
													position++
													break
												default:
													if buffer[position] != rune(':') {
														goto l24
													}
													position++
													break
												}
											}

										}
									l25:
										goto l6
									l24:
										position, tokenIndex, depth = position24, tokenIndex24, depth24
									}
									if !matchDot() {
										goto l6
									}
								l22:
									{
										position23, tokenIndex23, depth23 := position, tokenIndex, depth
										{
											position28, tokenIndex28, depth28 := position, tokenIndex, depth
											{
												position29, tokenIndex29, depth29 := position, tokenIndex, depth
												if buffer[position] != rune('\'') {
													goto l30
												}
												position++
												goto l29
											l30:
												position, tokenIndex, depth = position29, tokenIndex29, depth29
												{
													switch buffer[position] {
													case ' ':
														if buffer[position] != rune(' ') {
															goto l28
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l28
														}
														position++
														break
													default:
														if buffer[position] != rune(':') {
															goto l28
														}
														position++
														break
													}
												}

											}
										l29:
											goto l23
										l28:
											position, tokenIndex, depth = position28, tokenIndex28, depth28
										}
										if !matchDot() {
											goto l23
										}
										goto l22
									l23:
										position, tokenIndex, depth = position23, tokenIndex23, depth23
									}
									depth--
									add(rulePegText, position21)
								}
								{
									position32 := position
									depth++
									if buffer[position] != rune(':') {
										goto l6
									}
									position++
									if !rules[rulews]() {
										goto l6
									}
									depth--
									add(ruleCOLON, position32)
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
								position33 := position
								depth++
								{
									position34 := position
									depth++
									{
										position35, tokenIndex35, depth35 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l36
										}
										position++
										goto l35
									l36:
										position, tokenIndex, depth = position35, tokenIndex35, depth35
										if buffer[position] != rune('B') {
											goto l4
										}
										position++
									}
								l35:
									{
										position37, tokenIndex37, depth37 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l38
										}
										position++
										goto l37
									l38:
										position, tokenIndex, depth = position37, tokenIndex37, depth37
										if buffer[position] != rune('A') {
											goto l4
										}
										position++
									}
								l37:
									{
										position39, tokenIndex39, depth39 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l40
										}
										position++
										goto l39
									l40:
										position, tokenIndex, depth = position39, tokenIndex39, depth39
										if buffer[position] != rune('S') {
											goto l4
										}
										position++
									}
								l39:
									{
										position41, tokenIndex41, depth41 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l42
										}
										position++
										goto l41
									l42:
										position, tokenIndex, depth = position41, tokenIndex41, depth41
										if buffer[position] != rune('E') {
											goto l4
										}
										position++
									}
								l41:
									if !rules[rulews]() {
										goto l4
									}
									depth--
									add(ruleBASE, position34)
								}
								if !rules[ruleiri]() {
									goto l4
								}
								depth--
								add(rulebaseDecl, position33)
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
					position43 := position
					depth++
					{
						position44 := position
						depth++
						if !rules[ruleselect]() {
							goto l0
						}
						{
							position45, tokenIndex45, depth45 := position, tokenIndex, depth
							{
								position47 := position
								depth++
								{
									position48 := position
									depth++
									{
										position49, tokenIndex49, depth49 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l50
										}
										position++
										goto l49
									l50:
										position, tokenIndex, depth = position49, tokenIndex49, depth49
										if buffer[position] != rune('F') {
											goto l45
										}
										position++
									}
								l49:
									{
										position51, tokenIndex51, depth51 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l52
										}
										position++
										goto l51
									l52:
										position, tokenIndex, depth = position51, tokenIndex51, depth51
										if buffer[position] != rune('R') {
											goto l45
										}
										position++
									}
								l51:
									{
										position53, tokenIndex53, depth53 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l54
										}
										position++
										goto l53
									l54:
										position, tokenIndex, depth = position53, tokenIndex53, depth53
										if buffer[position] != rune('O') {
											goto l45
										}
										position++
									}
								l53:
									{
										position55, tokenIndex55, depth55 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l56
										}
										position++
										goto l55
									l56:
										position, tokenIndex, depth = position55, tokenIndex55, depth55
										if buffer[position] != rune('M') {
											goto l45
										}
										position++
									}
								l55:
									if !rules[rulews]() {
										goto l45
									}
									depth--
									add(ruleFROM, position48)
								}
								{
									position57, tokenIndex57, depth57 := position, tokenIndex, depth
									{
										position59 := position
										depth++
										{
											position60, tokenIndex60, depth60 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l61
											}
											position++
											goto l60
										l61:
											position, tokenIndex, depth = position60, tokenIndex60, depth60
											if buffer[position] != rune('N') {
												goto l57
											}
											position++
										}
									l60:
										{
											position62, tokenIndex62, depth62 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l63
											}
											position++
											goto l62
										l63:
											position, tokenIndex, depth = position62, tokenIndex62, depth62
											if buffer[position] != rune('A') {
												goto l57
											}
											position++
										}
									l62:
										{
											position64, tokenIndex64, depth64 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l65
											}
											position++
											goto l64
										l65:
											position, tokenIndex, depth = position64, tokenIndex64, depth64
											if buffer[position] != rune('M') {
												goto l57
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
												goto l57
											}
											position++
										}
									l66:
										{
											position68, tokenIndex68, depth68 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l69
											}
											position++
											goto l68
										l69:
											position, tokenIndex, depth = position68, tokenIndex68, depth68
											if buffer[position] != rune('D') {
												goto l57
											}
											position++
										}
									l68:
										if !rules[rulews]() {
											goto l57
										}
										depth--
										add(ruleNAMED, position59)
									}
									goto l58
								l57:
									position, tokenIndex, depth = position57, tokenIndex57, depth57
								}
							l58:
								if !rules[ruleiri]() {
									goto l45
								}
								depth--
								add(ruledatasetClause, position47)
							}
							goto l46
						l45:
							position, tokenIndex, depth = position45, tokenIndex45, depth45
						}
					l46:
						if !rules[rulewhereClause]() {
							goto l0
						}
						{
							position70 := position
							depth++
							{
								position71, tokenIndex71, depth71 := position, tokenIndex, depth
								{
									position73 := position
									depth++
									{
										position74, tokenIndex74, depth74 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l75
										}
										{
											position76, tokenIndex76, depth76 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l76
											}
											goto l77
										l76:
											position, tokenIndex, depth = position76, tokenIndex76, depth76
										}
									l77:
										goto l74
									l75:
										position, tokenIndex, depth = position74, tokenIndex74, depth74
										if !rules[ruleoffset]() {
											goto l71
										}
										{
											position78, tokenIndex78, depth78 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l78
											}
											goto l79
										l78:
											position, tokenIndex, depth = position78, tokenIndex78, depth78
										}
									l79:
									}
								l74:
									depth--
									add(rulelimitOffsetClauses, position73)
								}
								goto l72
							l71:
								position, tokenIndex, depth = position71, tokenIndex71, depth71
							}
						l72:
							depth--
							add(rulesolutionModifier, position70)
						}
						depth--
						add(ruleselectQuery, position44)
					}
					depth--
					add(rulequery, position43)
				}
				{
					position80, tokenIndex80, depth80 := position, tokenIndex, depth
					if !matchDot() {
						goto l80
					}
					goto l0
				l80:
					position, tokenIndex, depth = position80, tokenIndex80, depth80
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
		/* 2 prefixDecl <- <(PREFIX <(!('\'' / ((&(' ') ' ') | (&('\'') '\'') | (&(':') ':'))) .)+> COLON iri)> */
		nil,
		/* 3 baseDecl <- <(BASE iri)> */
		nil,
		/* 4 query <- <selectQuery> */
		nil,
		/* 5 selectQuery <- <(select datasetClause? whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position86, tokenIndex86, depth86 := position, tokenIndex, depth
			{
				position87 := position
				depth++
				{
					position88 := position
					depth++
					{
						position89, tokenIndex89, depth89 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l90
						}
						position++
						goto l89
					l90:
						position, tokenIndex, depth = position89, tokenIndex89, depth89
						if buffer[position] != rune('S') {
							goto l86
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
							goto l86
						}
						position++
					}
				l91:
					{
						position93, tokenIndex93, depth93 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l94
						}
						position++
						goto l93
					l94:
						position, tokenIndex, depth = position93, tokenIndex93, depth93
						if buffer[position] != rune('L') {
							goto l86
						}
						position++
					}
				l93:
					{
						position95, tokenIndex95, depth95 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l96
						}
						position++
						goto l95
					l96:
						position, tokenIndex, depth = position95, tokenIndex95, depth95
						if buffer[position] != rune('E') {
							goto l86
						}
						position++
					}
				l95:
					{
						position97, tokenIndex97, depth97 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l98
						}
						position++
						goto l97
					l98:
						position, tokenIndex, depth = position97, tokenIndex97, depth97
						if buffer[position] != rune('C') {
							goto l86
						}
						position++
					}
				l97:
					{
						position99, tokenIndex99, depth99 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l100
						}
						position++
						goto l99
					l100:
						position, tokenIndex, depth = position99, tokenIndex99, depth99
						if buffer[position] != rune('T') {
							goto l86
						}
						position++
					}
				l99:
					if !rules[rulews]() {
						goto l86
					}
					depth--
					add(ruleSELECT, position88)
				}
				{
					position101, tokenIndex101, depth101 := position, tokenIndex, depth
					{
						position103, tokenIndex103, depth103 := position, tokenIndex, depth
						{
							position105 := position
							depth++
							{
								position106, tokenIndex106, depth106 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l107
								}
								position++
								goto l106
							l107:
								position, tokenIndex, depth = position106, tokenIndex106, depth106
								if buffer[position] != rune('D') {
									goto l104
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
									goto l104
								}
								position++
							}
						l108:
							{
								position110, tokenIndex110, depth110 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l111
								}
								position++
								goto l110
							l111:
								position, tokenIndex, depth = position110, tokenIndex110, depth110
								if buffer[position] != rune('S') {
									goto l104
								}
								position++
							}
						l110:
							{
								position112, tokenIndex112, depth112 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l113
								}
								position++
								goto l112
							l113:
								position, tokenIndex, depth = position112, tokenIndex112, depth112
								if buffer[position] != rune('T') {
									goto l104
								}
								position++
							}
						l112:
							{
								position114, tokenIndex114, depth114 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l115
								}
								position++
								goto l114
							l115:
								position, tokenIndex, depth = position114, tokenIndex114, depth114
								if buffer[position] != rune('I') {
									goto l104
								}
								position++
							}
						l114:
							{
								position116, tokenIndex116, depth116 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l117
								}
								position++
								goto l116
							l117:
								position, tokenIndex, depth = position116, tokenIndex116, depth116
								if buffer[position] != rune('N') {
									goto l104
								}
								position++
							}
						l116:
							{
								position118, tokenIndex118, depth118 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l119
								}
								position++
								goto l118
							l119:
								position, tokenIndex, depth = position118, tokenIndex118, depth118
								if buffer[position] != rune('C') {
									goto l104
								}
								position++
							}
						l118:
							{
								position120, tokenIndex120, depth120 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l121
								}
								position++
								goto l120
							l121:
								position, tokenIndex, depth = position120, tokenIndex120, depth120
								if buffer[position] != rune('T') {
									goto l104
								}
								position++
							}
						l120:
							if !rules[rulews]() {
								goto l104
							}
							depth--
							add(ruleDISTINCT, position105)
						}
						goto l103
					l104:
						position, tokenIndex, depth = position103, tokenIndex103, depth103
						{
							position122 := position
							depth++
							{
								position123, tokenIndex123, depth123 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l124
								}
								position++
								goto l123
							l124:
								position, tokenIndex, depth = position123, tokenIndex123, depth123
								if buffer[position] != rune('R') {
									goto l101
								}
								position++
							}
						l123:
							{
								position125, tokenIndex125, depth125 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l126
								}
								position++
								goto l125
							l126:
								position, tokenIndex, depth = position125, tokenIndex125, depth125
								if buffer[position] != rune('E') {
									goto l101
								}
								position++
							}
						l125:
							{
								position127, tokenIndex127, depth127 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l128
								}
								position++
								goto l127
							l128:
								position, tokenIndex, depth = position127, tokenIndex127, depth127
								if buffer[position] != rune('D') {
									goto l101
								}
								position++
							}
						l127:
							{
								position129, tokenIndex129, depth129 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l130
								}
								position++
								goto l129
							l130:
								position, tokenIndex, depth = position129, tokenIndex129, depth129
								if buffer[position] != rune('U') {
									goto l101
								}
								position++
							}
						l129:
							{
								position131, tokenIndex131, depth131 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l132
								}
								position++
								goto l131
							l132:
								position, tokenIndex, depth = position131, tokenIndex131, depth131
								if buffer[position] != rune('C') {
									goto l101
								}
								position++
							}
						l131:
							{
								position133, tokenIndex133, depth133 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l134
								}
								position++
								goto l133
							l134:
								position, tokenIndex, depth = position133, tokenIndex133, depth133
								if buffer[position] != rune('E') {
									goto l101
								}
								position++
							}
						l133:
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l136
								}
								position++
								goto l135
							l136:
								position, tokenIndex, depth = position135, tokenIndex135, depth135
								if buffer[position] != rune('D') {
									goto l101
								}
								position++
							}
						l135:
							if !rules[rulews]() {
								goto l101
							}
							depth--
							add(ruleREDUCED, position122)
						}
					}
				l103:
					goto l102
				l101:
					position, tokenIndex, depth = position101, tokenIndex101, depth101
				}
			l102:
				{
					position137, tokenIndex137, depth137 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l138
					}
					goto l137
				l138:
					position, tokenIndex, depth = position137, tokenIndex137, depth137
					{
						position141 := position
						depth++
						if !rules[rulevar]() {
							goto l86
						}
						depth--
						add(ruleprojectionElem, position141)
					}
				l139:
					{
						position140, tokenIndex140, depth140 := position, tokenIndex, depth
						{
							position142 := position
							depth++
							if !rules[rulevar]() {
								goto l140
							}
							depth--
							add(ruleprojectionElem, position142)
						}
						goto l139
					l140:
						position, tokenIndex, depth = position140, tokenIndex140, depth140
					}
				}
			l137:
				depth--
				add(ruleselect, position87)
			}
			return true
		l86:
			position, tokenIndex, depth = position86, tokenIndex86, depth86
			return false
		},
		/* 7 subSelect <- <(select whereClause)> */
		func() bool {
			position143, tokenIndex143, depth143 := position, tokenIndex, depth
			{
				position144 := position
				depth++
				if !rules[ruleselect]() {
					goto l143
				}
				if !rules[rulewhereClause]() {
					goto l143
				}
				depth--
				add(rulesubSelect, position144)
			}
			return true
		l143:
			position, tokenIndex, depth = position143, tokenIndex143, depth143
			return false
		},
		/* 8 projectionElem <- <var> */
		nil,
		/* 9 datasetClause <- <(FROM NAMED? iri)> */
		nil,
		/* 10 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position147, tokenIndex147, depth147 := position, tokenIndex, depth
			{
				position148 := position
				depth++
				{
					position149, tokenIndex149, depth149 := position, tokenIndex, depth
					{
						position151 := position
						depth++
						{
							position152, tokenIndex152, depth152 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l153
							}
							position++
							goto l152
						l153:
							position, tokenIndex, depth = position152, tokenIndex152, depth152
							if buffer[position] != rune('W') {
								goto l149
							}
							position++
						}
					l152:
						{
							position154, tokenIndex154, depth154 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l155
							}
							position++
							goto l154
						l155:
							position, tokenIndex, depth = position154, tokenIndex154, depth154
							if buffer[position] != rune('H') {
								goto l149
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
								goto l149
							}
							position++
						}
					l156:
						{
							position158, tokenIndex158, depth158 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l159
							}
							position++
							goto l158
						l159:
							position, tokenIndex, depth = position158, tokenIndex158, depth158
							if buffer[position] != rune('R') {
								goto l149
							}
							position++
						}
					l158:
						{
							position160, tokenIndex160, depth160 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l161
							}
							position++
							goto l160
						l161:
							position, tokenIndex, depth = position160, tokenIndex160, depth160
							if buffer[position] != rune('E') {
								goto l149
							}
							position++
						}
					l160:
						if !rules[rulews]() {
							goto l149
						}
						depth--
						add(ruleWHERE, position151)
					}
					goto l150
				l149:
					position, tokenIndex, depth = position149, tokenIndex149, depth149
				}
			l150:
				if !rules[rulegroupGraphPattern]() {
					goto l147
				}
				depth--
				add(rulewhereClause, position148)
			}
			return true
		l147:
			position, tokenIndex, depth = position147, tokenIndex147, depth147
			return false
		},
		/* 11 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position162, tokenIndex162, depth162 := position, tokenIndex, depth
			{
				position163 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l162
				}
				{
					position164, tokenIndex164, depth164 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l165
					}
					goto l164
				l165:
					position, tokenIndex, depth = position164, tokenIndex164, depth164
					if !rules[rulegraphPattern]() {
						goto l162
					}
				}
			l164:
				if !rules[ruleRBRACE]() {
					goto l162
				}
				depth--
				add(rulegroupGraphPattern, position163)
			}
			return true
		l162:
			position, tokenIndex, depth = position162, tokenIndex162, depth162
			return false
		},
		/* 12 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position167 := position
				depth++
				{
					position168, tokenIndex168, depth168 := position, tokenIndex, depth
					{
						position170 := position
						depth++
						{
							position171 := position
							depth++
							if !rules[ruletriplesSameSubjectPath]() {
								goto l168
							}
						l172:
							{
								position173, tokenIndex173, depth173 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l173
								}
								if !rules[ruletriplesSameSubjectPath]() {
									goto l173
								}
								goto l172
							l173:
								position, tokenIndex, depth = position173, tokenIndex173, depth173
							}
							{
								position174, tokenIndex174, depth174 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l174
								}
								goto l175
							l174:
								position, tokenIndex, depth = position174, tokenIndex174, depth174
							}
						l175:
							depth--
							add(ruletriplesBlock, position171)
						}
						depth--
						add(rulebasicGraphPattern, position170)
					}
					goto l169
				l168:
					position, tokenIndex, depth = position168, tokenIndex168, depth168
				}
			l169:
				{
					position176, tokenIndex176, depth176 := position, tokenIndex, depth
					{
						position178 := position
						depth++
						{
							position179, tokenIndex179, depth179 := position, tokenIndex, depth
							{
								position181 := position
								depth++
								{
									position182 := position
									depth++
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
											goto l180
										}
										position++
									}
								l183:
									{
										position185, tokenIndex185, depth185 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l186
										}
										position++
										goto l185
									l186:
										position, tokenIndex, depth = position185, tokenIndex185, depth185
										if buffer[position] != rune('P') {
											goto l180
										}
										position++
									}
								l185:
									{
										position187, tokenIndex187, depth187 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l188
										}
										position++
										goto l187
									l188:
										position, tokenIndex, depth = position187, tokenIndex187, depth187
										if buffer[position] != rune('T') {
											goto l180
										}
										position++
									}
								l187:
									{
										position189, tokenIndex189, depth189 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l190
										}
										position++
										goto l189
									l190:
										position, tokenIndex, depth = position189, tokenIndex189, depth189
										if buffer[position] != rune('I') {
											goto l180
										}
										position++
									}
								l189:
									{
										position191, tokenIndex191, depth191 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l192
										}
										position++
										goto l191
									l192:
										position, tokenIndex, depth = position191, tokenIndex191, depth191
										if buffer[position] != rune('O') {
											goto l180
										}
										position++
									}
								l191:
									{
										position193, tokenIndex193, depth193 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l194
										}
										position++
										goto l193
									l194:
										position, tokenIndex, depth = position193, tokenIndex193, depth193
										if buffer[position] != rune('N') {
											goto l180
										}
										position++
									}
								l193:
									{
										position195, tokenIndex195, depth195 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l196
										}
										position++
										goto l195
									l196:
										position, tokenIndex, depth = position195, tokenIndex195, depth195
										if buffer[position] != rune('A') {
											goto l180
										}
										position++
									}
								l195:
									{
										position197, tokenIndex197, depth197 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l198
										}
										position++
										goto l197
									l198:
										position, tokenIndex, depth = position197, tokenIndex197, depth197
										if buffer[position] != rune('L') {
											goto l180
										}
										position++
									}
								l197:
									if !rules[rulews]() {
										goto l180
									}
									depth--
									add(ruleOPTIONAL, position182)
								}
								if !rules[ruleLBRACE]() {
									goto l180
								}
								{
									position199, tokenIndex199, depth199 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l200
									}
									goto l199
								l200:
									position, tokenIndex, depth = position199, tokenIndex199, depth199
									if !rules[rulegraphPattern]() {
										goto l180
									}
								}
							l199:
								if !rules[ruleRBRACE]() {
									goto l180
								}
								depth--
								add(ruleoptionalGraphPattern, position181)
							}
							goto l179
						l180:
							position, tokenIndex, depth = position179, tokenIndex179, depth179
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l176
							}
						}
					l179:
						depth--
						add(rulegraphPatternNotTriples, position178)
					}
					{
						position201, tokenIndex201, depth201 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l201
						}
						goto l202
					l201:
						position, tokenIndex, depth = position201, tokenIndex201, depth201
					}
				l202:
					if !rules[rulegraphPattern]() {
						goto l176
					}
					goto l177
				l176:
					position, tokenIndex, depth = position176, tokenIndex176, depth176
				}
			l177:
				depth--
				add(rulegraphPattern, position167)
			}
			return true
		},
		/* 13 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 14 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 15 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position205, tokenIndex205, depth205 := position, tokenIndex, depth
			{
				position206 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l205
				}
				{
					position207, tokenIndex207, depth207 := position, tokenIndex, depth
					{
						position209 := position
						depth++
						{
							position210, tokenIndex210, depth210 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l211
							}
							position++
							goto l210
						l211:
							position, tokenIndex, depth = position210, tokenIndex210, depth210
							if buffer[position] != rune('U') {
								goto l207
							}
							position++
						}
					l210:
						{
							position212, tokenIndex212, depth212 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l213
							}
							position++
							goto l212
						l213:
							position, tokenIndex, depth = position212, tokenIndex212, depth212
							if buffer[position] != rune('N') {
								goto l207
							}
							position++
						}
					l212:
						{
							position214, tokenIndex214, depth214 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l215
							}
							position++
							goto l214
						l215:
							position, tokenIndex, depth = position214, tokenIndex214, depth214
							if buffer[position] != rune('I') {
								goto l207
							}
							position++
						}
					l214:
						{
							position216, tokenIndex216, depth216 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l217
							}
							position++
							goto l216
						l217:
							position, tokenIndex, depth = position216, tokenIndex216, depth216
							if buffer[position] != rune('O') {
								goto l207
							}
							position++
						}
					l216:
						{
							position218, tokenIndex218, depth218 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l219
							}
							position++
							goto l218
						l219:
							position, tokenIndex, depth = position218, tokenIndex218, depth218
							if buffer[position] != rune('N') {
								goto l207
							}
							position++
						}
					l218:
						if !rules[rulews]() {
							goto l207
						}
						depth--
						add(ruleUNION, position209)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l207
					}
					goto l208
				l207:
					position, tokenIndex, depth = position207, tokenIndex207, depth207
				}
			l208:
				depth--
				add(rulegroupOrUnionGraphPattern, position206)
			}
			return true
		l205:
			position, tokenIndex, depth = position205, tokenIndex205, depth205
			return false
		},
		/* 16 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 17 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		nil,
		/* 18 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position222, tokenIndex222, depth222 := position, tokenIndex, depth
			{
				position223 := position
				depth++
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l225
					}
					if !rules[rulepropertyListPath]() {
						goto l225
					}
					goto l224
				l225:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
					{
						position226 := position
						depth++
						{
							position227, tokenIndex227, depth227 := position, tokenIndex, depth
							{
								position229 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l228
								}
								if !rules[rulegraphNodePath]() {
									goto l228
								}
							l230:
								{
									position231, tokenIndex231, depth231 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l231
									}
									goto l230
								l231:
									position, tokenIndex, depth = position231, tokenIndex231, depth231
								}
								if !rules[ruleRPAREN]() {
									goto l228
								}
								depth--
								add(rulecollectionPath, position229)
							}
							goto l227
						l228:
							position, tokenIndex, depth = position227, tokenIndex227, depth227
							{
								position232 := position
								depth++
								{
									position233 := position
									depth++
									if buffer[position] != rune('[') {
										goto l222
									}
									position++
									if !rules[rulews]() {
										goto l222
									}
									depth--
									add(ruleLBRACK, position233)
								}
								if !rules[rulepropertyListPath]() {
									goto l222
								}
								{
									position234 := position
									depth++
									if buffer[position] != rune(']') {
										goto l222
									}
									position++
									if !rules[rulews]() {
										goto l222
									}
									depth--
									add(ruleRBRACK, position234)
								}
								depth--
								add(ruleblankNodePropertyListPath, position232)
							}
						}
					l227:
						depth--
						add(ruletriplesNodePath, position226)
					}
					if !rules[rulepropertyListPath]() {
						goto l222
					}
				}
			l224:
				depth--
				add(ruletriplesSameSubjectPath, position223)
			}
			return true
		l222:
			position, tokenIndex, depth = position222, tokenIndex222, depth222
			return false
		},
		/* 19 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position235, tokenIndex235, depth235 := position, tokenIndex, depth
			{
				position236 := position
				depth++
				{
					position237, tokenIndex237, depth237 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l238
					}
					goto l237
				l238:
					position, tokenIndex, depth = position237, tokenIndex237, depth237
					{
						position239 := position
						depth++
						{
							switch buffer[position] {
							case '(':
								{
									position241 := position
									depth++
									if buffer[position] != rune('(') {
										goto l235
									}
									position++
									if !rules[rulews]() {
										goto l235
									}
									if buffer[position] != rune(')') {
										goto l235
									}
									position++
									if !rules[rulews]() {
										goto l235
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
												goto l235
											}
											position++
											if !rules[rulews]() {
												goto l235
											}
											if buffer[position] != rune(']') {
												goto l235
											}
											position++
											if !rules[rulews]() {
												goto l235
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
													goto l235
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
													goto l235
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
													goto l235
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
													goto l235
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
													goto l235
												}
												position++
											}
										l275:
											if !rules[rulews]() {
												goto l235
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
											goto l235
										}
										position++
									l279:
										{
											position280, tokenIndex280, depth280 := position, tokenIndex, depth
											{
												position281, tokenIndex281, depth281 := position, tokenIndex, depth
												{
													position282, tokenIndex282, depth282 := position, tokenIndex, depth
													if buffer[position] != rune('\'') {
														goto l283
													}
													position++
													goto l282
												l283:
													position, tokenIndex, depth = position282, tokenIndex282, depth282
													if buffer[position] != rune('"') {
														goto l284
													}
													position++
													goto l282
												l284:
													position, tokenIndex, depth = position282, tokenIndex282, depth282
													if buffer[position] != rune('\'') {
														goto l281
													}
													position++
												}
											l282:
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
											goto l235
										}
										position++
										depth--
										add(rulestring, position278)
									}
									{
										position285, tokenIndex285, depth285 := position, tokenIndex, depth
										{
											position287, tokenIndex287, depth287 := position, tokenIndex, depth
											if buffer[position] != rune('@') {
												goto l288
											}
											position++
											{
												position291, tokenIndex291, depth291 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l292
												}
												position++
												goto l291
											l292:
												position, tokenIndex, depth = position291, tokenIndex291, depth291
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l288
												}
												position++
											}
										l291:
										l289:
											{
												position290, tokenIndex290, depth290 := position, tokenIndex, depth
												{
													position293, tokenIndex293, depth293 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('a') || c > rune('z') {
														goto l294
													}
													position++
													goto l293
												l294:
													position, tokenIndex, depth = position293, tokenIndex293, depth293
													if c := buffer[position]; c < rune('A') || c > rune('Z') {
														goto l290
													}
													position++
												}
											l293:
												goto l289
											l290:
												position, tokenIndex, depth = position290, tokenIndex290, depth290
											}
										l295:
											{
												position296, tokenIndex296, depth296 := position, tokenIndex, depth
												if buffer[position] != rune('-') {
													goto l296
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l296
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l296
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l296
														}
														position++
														break
													}
												}

											l297:
												{
													position298, tokenIndex298, depth298 := position, tokenIndex, depth
													{
														switch buffer[position] {
														case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l298
															}
															position++
															break
														case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
															if c := buffer[position]; c < rune('A') || c > rune('Z') {
																goto l298
															}
															position++
															break
														default:
															if c := buffer[position]; c < rune('a') || c > rune('z') {
																goto l298
															}
															position++
															break
														}
													}

													goto l297
												l298:
													position, tokenIndex, depth = position298, tokenIndex298, depth298
												}
												goto l295
											l296:
												position, tokenIndex, depth = position296, tokenIndex296, depth296
											}
											goto l287
										l288:
											position, tokenIndex, depth = position287, tokenIndex287, depth287
											if buffer[position] != rune('^') {
												goto l285
											}
											position++
											if buffer[position] != rune('^') {
												goto l285
											}
											position++
											{
												position301 := position
												depth++
												if !rules[ruleiri]() {
													goto l285
												}
												depth--
												add(rulePegText, position301)
											}
										}
									l287:
										goto l286
									l285:
										position, tokenIndex, depth = position285, tokenIndex285, depth285
									}
								l286:
									if !rules[rulews]() {
										goto l235
									}
									depth--
									add(ruleliteral, position277)
								}
								break
							case '<':
								if !rules[ruleiri]() {
									goto l235
								}
								break
							default:
								{
									position302 := position
									depth++
									{
										position303, tokenIndex303, depth303 := position, tokenIndex, depth
										{
											position305, tokenIndex305, depth305 := position, tokenIndex, depth
											if buffer[position] != rune('+') {
												goto l306
											}
											position++
											goto l305
										l306:
											position, tokenIndex, depth = position305, tokenIndex305, depth305
											if buffer[position] != rune('-') {
												goto l303
											}
											position++
										}
									l305:
										goto l304
									l303:
										position, tokenIndex, depth = position303, tokenIndex303, depth303
									}
								l304:
									{
										position309, tokenIndex309, depth309 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l310
										}
										position++
										goto l309
									l310:
										position, tokenIndex, depth = position309, tokenIndex309, depth309
										if buffer[position] != rune('.') {
											goto l235
										}
										position++
									}
								l309:
								l307:
									{
										position308, tokenIndex308, depth308 := position, tokenIndex, depth
										{
											position311, tokenIndex311, depth311 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l312
											}
											position++
											goto l311
										l312:
											position, tokenIndex, depth = position311, tokenIndex311, depth311
											if buffer[position] != rune('.') {
												goto l308
											}
											position++
										}
									l311:
										goto l307
									l308:
										position, tokenIndex, depth = position308, tokenIndex308, depth308
									}
									if !rules[rulews]() {
										goto l235
									}
									depth--
									add(rulenumericLiteral, position302)
								}
								break
							}
						}

						depth--
						add(rulegraphTerm, position239)
					}
				}
			l237:
				depth--
				add(rulevarOrTerm, position236)
			}
			return true
		l235:
			position, tokenIndex, depth = position235, tokenIndex235, depth235
			return false
		},
		/* 20 graphTerm <- <((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('<') iri) | (&('+' | '-' | '.' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral))> */
		nil,
		/* 21 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 22 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 23 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 24 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			position317, tokenIndex317, depth317 := position, tokenIndex, depth
			{
				position318 := position
				depth++
				{
					position319, tokenIndex319, depth319 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l320
					}
					goto l319
				l320:
					position, tokenIndex, depth = position319, tokenIndex319, depth319
					{
						position321 := position
						depth++
						if !rules[rulepath]() {
							goto l317
						}
						depth--
						add(ruleverbPath, position321)
					}
				}
			l319:
				if !rules[ruleobjectListPath]() {
					goto l317
				}
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					{
						position324 := position
						depth++
						if buffer[position] != rune(';') {
							goto l322
						}
						position++
						if !rules[rulews]() {
							goto l322
						}
						depth--
						add(ruleSEMICOLON, position324)
					}
					if !rules[rulepropertyListPath]() {
						goto l322
					}
					goto l323
				l322:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
				}
			l323:
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
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				{
					position328 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l326
					}
				l329:
					{
						position330, tokenIndex330, depth330 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l330
						}
						if !rules[rulepathSequence]() {
							goto l330
						}
						goto l329
					l330:
						position, tokenIndex, depth = position330, tokenIndex330, depth330
					}
					depth--
					add(rulepathAlternative, position328)
				}
				depth--
				add(rulepath, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 27 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 28 pathSequence <- <(pathElt (SLASH pathElt)*)> */
		func() bool {
			position332, tokenIndex332, depth332 := position, tokenIndex, depth
			{
				position333 := position
				depth++
				if !rules[rulepathElt]() {
					goto l332
				}
			l334:
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					{
						position336 := position
						depth++
						if buffer[position] != rune('/') {
							goto l335
						}
						position++
						if !rules[rulews]() {
							goto l335
						}
						depth--
						add(ruleSLASH, position336)
					}
					if !rules[rulepathElt]() {
						goto l335
					}
					goto l334
				l335:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
				}
				depth--
				add(rulepathSequence, position333)
			}
			return true
		l332:
			position, tokenIndex, depth = position332, tokenIndex332, depth332
			return false
		},
		/* 29 pathElt <- <(INVERSE? pathPrimary &pathMod?)> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if !rules[ruleINVERSE]() {
						goto l339
					}
					goto l340
				l339:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
				}
			l340:
				{
					position341 := position
					depth++
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
								position343 := position
								depth++
								if buffer[position] != rune('!') {
									goto l337
								}
								position++
								if !rules[rulews]() {
									goto l337
								}
								depth--
								add(ruleNOT, position343)
							}
							{
								position344 := position
								depth++
								{
									position345, tokenIndex345, depth345 := position, tokenIndex, depth
									if !rules[rulepathOneInPropertySet]() {
										goto l346
									}
									goto l345
								l346:
									position, tokenIndex, depth = position345, tokenIndex345, depth345
									if !rules[ruleLPAREN]() {
										goto l337
									}
									{
										position347, tokenIndex347, depth347 := position, tokenIndex, depth
										if !rules[rulepathOneInPropertySet]() {
											goto l347
										}
									l349:
										{
											position350, tokenIndex350, depth350 := position, tokenIndex, depth
											if !rules[rulePIPE]() {
												goto l350
											}
											if !rules[rulepathOneInPropertySet]() {
												goto l350
											}
											goto l349
										l350:
											position, tokenIndex, depth = position350, tokenIndex350, depth350
										}
										goto l348
									l347:
										position, tokenIndex, depth = position347, tokenIndex347, depth347
									}
								l348:
									if !rules[ruleRPAREN]() {
										goto l337
									}
								}
							l345:
								depth--
								add(rulepathNegatedPropertySet, position344)
							}
							break
						case 'a':
							if !rules[ruleISA]() {
								goto l337
							}
							break
						default:
							if !rules[ruleiri]() {
								goto l337
							}
							break
						}
					}

					depth--
					add(rulepathPrimary, position341)
				}
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					{
						position352, tokenIndex352, depth352 := position, tokenIndex, depth
						{
							position354 := position
							depth++
							{
								switch buffer[position] {
								case '+':
									{
										position356 := position
										depth++
										if buffer[position] != rune('+') {
											goto l352
										}
										position++
										if !rules[rulews]() {
											goto l352
										}
										depth--
										add(rulePLUS, position356)
									}
									break
								case '?':
									{
										position357 := position
										depth++
										if buffer[position] != rune('?') {
											goto l352
										}
										position++
										if !rules[rulews]() {
											goto l352
										}
										depth--
										add(ruleQUESTION, position357)
									}
									break
								default:
									if !rules[ruleSTAR]() {
										goto l352
									}
									break
								}
							}

							depth--
							add(rulepathMod, position354)
						}
						goto l353
					l352:
						position, tokenIndex, depth = position352, tokenIndex352, depth352
					}
				l353:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
				}
				depth--
				add(rulepathElt, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 30 pathPrimary <- <((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA) | (&('<') iri))> */
		nil,
		/* 31 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 32 pathOneInPropertySet <- <((&('^') (INVERSE (iri / ISA))) | (&('a') ISA) | (&('<') iri))> */
		func() bool {
			position360, tokenIndex360, depth360 := position, tokenIndex, depth
			{
				position361 := position
				depth++
				{
					switch buffer[position] {
					case '^':
						if !rules[ruleINVERSE]() {
							goto l360
						}
						{
							position363, tokenIndex363, depth363 := position, tokenIndex, depth
							if !rules[ruleiri]() {
								goto l364
							}
							goto l363
						l364:
							position, tokenIndex, depth = position363, tokenIndex363, depth363
							if !rules[ruleISA]() {
								goto l360
							}
						}
					l363:
						break
					case 'a':
						if !rules[ruleISA]() {
							goto l360
						}
						break
					default:
						if !rules[ruleiri]() {
							goto l360
						}
						break
					}
				}

				depth--
				add(rulepathOneInPropertySet, position361)
			}
			return true
		l360:
			position, tokenIndex, depth = position360, tokenIndex360, depth360
			return false
		},
		/* 33 pathMod <- <((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR))> */
		nil,
		/* 34 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position366, tokenIndex366, depth366 := position, tokenIndex, depth
			{
				position367 := position
				depth++
				{
					position368 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l366
					}
					depth--
					add(ruleobjectPath, position368)
				}
			l369:
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					{
						position371 := position
						depth++
						if buffer[position] != rune(',') {
							goto l370
						}
						position++
						if !rules[rulews]() {
							goto l370
						}
						depth--
						add(ruleCOMMA, position371)
					}
					if !rules[ruleobjectListPath]() {
						goto l370
					}
					goto l369
				l370:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
				}
				depth--
				add(ruleobjectListPath, position367)
			}
			return true
		l366:
			position, tokenIndex, depth = position366, tokenIndex366, depth366
			return false
		},
		/* 35 objectPath <- <graphNodePath> */
		nil,
		/* 36 graphNodePath <- <varOrTerm> */
		func() bool {
			position373, tokenIndex373, depth373 := position, tokenIndex, depth
			{
				position374 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l373
				}
				depth--
				add(rulegraphNodePath, position374)
			}
			return true
		l373:
			position, tokenIndex, depth = position373, tokenIndex373, depth373
			return false
		},
		/* 37 solutionModifier <- <limitOffsetClauses?> */
		nil,
		/* 38 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 39 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position377, tokenIndex377, depth377 := position, tokenIndex, depth
			{
				position378 := position
				depth++
				{
					position379 := position
					depth++
					{
						position380, tokenIndex380, depth380 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l381
						}
						position++
						goto l380
					l381:
						position, tokenIndex, depth = position380, tokenIndex380, depth380
						if buffer[position] != rune('L') {
							goto l377
						}
						position++
					}
				l380:
					{
						position382, tokenIndex382, depth382 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l383
						}
						position++
						goto l382
					l383:
						position, tokenIndex, depth = position382, tokenIndex382, depth382
						if buffer[position] != rune('I') {
							goto l377
						}
						position++
					}
				l382:
					{
						position384, tokenIndex384, depth384 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l385
						}
						position++
						goto l384
					l385:
						position, tokenIndex, depth = position384, tokenIndex384, depth384
						if buffer[position] != rune('M') {
							goto l377
						}
						position++
					}
				l384:
					{
						position386, tokenIndex386, depth386 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l387
						}
						position++
						goto l386
					l387:
						position, tokenIndex, depth = position386, tokenIndex386, depth386
						if buffer[position] != rune('I') {
							goto l377
						}
						position++
					}
				l386:
					{
						position388, tokenIndex388, depth388 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l389
						}
						position++
						goto l388
					l389:
						position, tokenIndex, depth = position388, tokenIndex388, depth388
						if buffer[position] != rune('T') {
							goto l377
						}
						position++
					}
				l388:
					if !rules[rulews]() {
						goto l377
					}
					depth--
					add(ruleLIMIT, position379)
				}
				if !rules[ruleINTEGER]() {
					goto l377
				}
				depth--
				add(rulelimit, position378)
			}
			return true
		l377:
			position, tokenIndex, depth = position377, tokenIndex377, depth377
			return false
		},
		/* 40 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392 := position
					depth++
					{
						position393, tokenIndex393, depth393 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l394
						}
						position++
						goto l393
					l394:
						position, tokenIndex, depth = position393, tokenIndex393, depth393
						if buffer[position] != rune('O') {
							goto l390
						}
						position++
					}
				l393:
					{
						position395, tokenIndex395, depth395 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l396
						}
						position++
						goto l395
					l396:
						position, tokenIndex, depth = position395, tokenIndex395, depth395
						if buffer[position] != rune('F') {
							goto l390
						}
						position++
					}
				l395:
					{
						position397, tokenIndex397, depth397 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l398
						}
						position++
						goto l397
					l398:
						position, tokenIndex, depth = position397, tokenIndex397, depth397
						if buffer[position] != rune('F') {
							goto l390
						}
						position++
					}
				l397:
					{
						position399, tokenIndex399, depth399 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l400
						}
						position++
						goto l399
					l400:
						position, tokenIndex, depth = position399, tokenIndex399, depth399
						if buffer[position] != rune('S') {
							goto l390
						}
						position++
					}
				l399:
					{
						position401, tokenIndex401, depth401 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l402
						}
						position++
						goto l401
					l402:
						position, tokenIndex, depth = position401, tokenIndex401, depth401
						if buffer[position] != rune('E') {
							goto l390
						}
						position++
					}
				l401:
					{
						position403, tokenIndex403, depth403 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l404
						}
						position++
						goto l403
					l404:
						position, tokenIndex, depth = position403, tokenIndex403, depth403
						if buffer[position] != rune('T') {
							goto l390
						}
						position++
					}
				l403:
					if !rules[rulews]() {
						goto l390
					}
					depth--
					add(ruleOFFSET, position392)
				}
				if !rules[ruleINTEGER]() {
					goto l390
				}
				depth--
				add(ruleoffset, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 41 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position405, tokenIndex405, depth405 := position, tokenIndex, depth
			{
				position406 := position
				depth++
				{
					position407, tokenIndex407, depth407 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l408
					}
					position++
					goto l407
				l408:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
					if buffer[position] != rune('$') {
						goto l405
					}
					position++
				}
			l407:
				{
					position409 := position
					depth++
					{
						position410, tokenIndex410, depth410 := position, tokenIndex, depth
						if !rules[rulePN_CHARS_U]() {
							goto l411
						}
						goto l410
					l411:
						position, tokenIndex, depth = position410, tokenIndex410, depth410
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l405
						}
						position++
					}
				l410:
				l412:
					{
						position413, tokenIndex413, depth413 := position, tokenIndex, depth
						{
							position414 := position
							depth++
							{
								position415, tokenIndex415, depth415 := position, tokenIndex, depth
								if !rules[rulePN_CHARS_U]() {
									goto l416
								}
								goto l415
							l416:
								position, tokenIndex, depth = position415, tokenIndex415, depth415
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l413
								}
								position++
							}
						l415:
							depth--
							add(ruleVAR_CHAR, position414)
						}
						goto l412
					l413:
						position, tokenIndex, depth = position413, tokenIndex413, depth413
					}
					depth--
					add(ruleVARNAME, position409)
				}
				if !rules[rulews]() {
					goto l405
				}
				depth--
				add(rulevar, position406)
			}
			return true
		l405:
			position, tokenIndex, depth = position405, tokenIndex405, depth405
			return false
		},
		/* 42 iri <- <('<' (!('\'' / '>' / '\'') .)* '>' ws)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				if buffer[position] != rune('<') {
					goto l417
				}
				position++
			l419:
				{
					position420, tokenIndex420, depth420 := position, tokenIndex, depth
					{
						position421, tokenIndex421, depth421 := position, tokenIndex, depth
						{
							position422, tokenIndex422, depth422 := position, tokenIndex, depth
							if buffer[position] != rune('\'') {
								goto l423
							}
							position++
							goto l422
						l423:
							position, tokenIndex, depth = position422, tokenIndex422, depth422
							if buffer[position] != rune('>') {
								goto l424
							}
							position++
							goto l422
						l424:
							position, tokenIndex, depth = position422, tokenIndex422, depth422
							if buffer[position] != rune('\'') {
								goto l421
							}
							position++
						}
					l422:
						goto l420
					l421:
						position, tokenIndex, depth = position421, tokenIndex421, depth421
					}
					if !matchDot() {
						goto l420
					}
					goto l419
				l420:
					position, tokenIndex, depth = position420, tokenIndex420, depth420
				}
				if buffer[position] != rune('>') {
					goto l417
				}
				position++
				if !rules[rulews]() {
					goto l417
				}
				depth--
				add(ruleiri, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 43 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' <iri>))? ws)> */
		nil,
		/* 44 string <- <('"' (!('\'' / '"' / '\'') .)* '"')> */
		nil,
		/* 45 numericLiteral <- <(('+' / '-')? ([0-9] / '.')+ ws)> */
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
			position434, tokenIndex434, depth434 := position, tokenIndex, depth
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
								goto l440
							}
							position++
							goto l439
						l440:
							position, tokenIndex, depth = position439, tokenIndex439, depth439
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l437
							}
							position++
						}
					l439:
						depth--
						add(rulePN_CHARS_BASE, position438)
					}
					goto l436
				l437:
					position, tokenIndex, depth = position436, tokenIndex436, depth436
					if buffer[position] != rune('_') {
						goto l434
					}
					position++
				}
			l436:
				depth--
				add(rulePN_CHARS_U, position435)
			}
			return true
		l434:
			position, tokenIndex, depth = position434, tokenIndex434, depth434
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
			position453, tokenIndex453, depth453 := position, tokenIndex, depth
			{
				position454 := position
				depth++
				if buffer[position] != rune('{') {
					goto l453
				}
				position++
				if !rules[rulews]() {
					goto l453
				}
				depth--
				add(ruleLBRACE, position454)
			}
			return true
		l453:
			position, tokenIndex, depth = position453, tokenIndex453, depth453
			return false
		},
		/* 66 RBRACE <- <('}' ws)> */
		func() bool {
			position455, tokenIndex455, depth455 := position, tokenIndex, depth
			{
				position456 := position
				depth++
				if buffer[position] != rune('}') {
					goto l455
				}
				position++
				if !rules[rulews]() {
					goto l455
				}
				depth--
				add(ruleRBRACE, position456)
			}
			return true
		l455:
			position, tokenIndex, depth = position455, tokenIndex455, depth455
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
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if buffer[position] != rune('.') {
					goto l461
				}
				position++
				if !rules[rulews]() {
					goto l461
				}
				depth--
				add(ruleDOT, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 72 COLON <- <(':' ws)> */
		nil,
		/* 73 PIPE <- <('|' ws)> */
		func() bool {
			position464, tokenIndex464, depth464 := position, tokenIndex, depth
			{
				position465 := position
				depth++
				if buffer[position] != rune('|') {
					goto l464
				}
				position++
				if !rules[rulews]() {
					goto l464
				}
				depth--
				add(rulePIPE, position465)
			}
			return true
		l464:
			position, tokenIndex, depth = position464, tokenIndex464, depth464
			return false
		},
		/* 74 SLASH <- <('/' ws)> */
		nil,
		/* 75 INVERSE <- <('^' ws)> */
		func() bool {
			position467, tokenIndex467, depth467 := position, tokenIndex, depth
			{
				position468 := position
				depth++
				if buffer[position] != rune('^') {
					goto l467
				}
				position++
				if !rules[rulews]() {
					goto l467
				}
				depth--
				add(ruleINVERSE, position468)
			}
			return true
		l467:
			position, tokenIndex, depth = position467, tokenIndex467, depth467
			return false
		},
		/* 76 LPAREN <- <('(' ws)> */
		func() bool {
			position469, tokenIndex469, depth469 := position, tokenIndex, depth
			{
				position470 := position
				depth++
				if buffer[position] != rune('(') {
					goto l469
				}
				position++
				if !rules[rulews]() {
					goto l469
				}
				depth--
				add(ruleLPAREN, position470)
			}
			return true
		l469:
			position, tokenIndex, depth = position469, tokenIndex469, depth469
			return false
		},
		/* 77 RPAREN <- <(')' ws)> */
		func() bool {
			position471, tokenIndex471, depth471 := position, tokenIndex, depth
			{
				position472 := position
				depth++
				if buffer[position] != rune(')') {
					goto l471
				}
				position++
				if !rules[rulews]() {
					goto l471
				}
				depth--
				add(ruleRPAREN, position472)
			}
			return true
		l471:
			position, tokenIndex, depth = position471, tokenIndex471, depth471
			return false
		},
		/* 78 ISA <- <('a' ws)> */
		func() bool {
			position473, tokenIndex473, depth473 := position, tokenIndex, depth
			{
				position474 := position
				depth++
				if buffer[position] != rune('a') {
					goto l473
				}
				position++
				if !rules[rulews]() {
					goto l473
				}
				depth--
				add(ruleISA, position474)
			}
			return true
		l473:
			position, tokenIndex, depth = position473, tokenIndex473, depth473
			return false
		},
		/* 79 NOT <- <('!' ws)> */
		nil,
		/* 80 STAR <- <('*' ws)> */
		func() bool {
			position476, tokenIndex476, depth476 := position, tokenIndex, depth
			{
				position477 := position
				depth++
				if buffer[position] != rune('*') {
					goto l476
				}
				position++
				if !rules[rulews]() {
					goto l476
				}
				depth--
				add(ruleSTAR, position477)
			}
			return true
		l476:
			position, tokenIndex, depth = position476, tokenIndex476, depth476
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
			position484, tokenIndex484, depth484 := position, tokenIndex, depth
			{
				position485 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l484
				}
				position++
			l486:
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l487
					}
					position++
					goto l486
				l487:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
				}
				if !rules[rulews]() {
					goto l484
				}
				depth--
				add(ruleINTEGER, position485)
			}
			return true
		l484:
			position, tokenIndex, depth = position484, tokenIndex484, depth484
			return false
		},
		/* 88 ws <- <((&('\f') '\f') | (&('\r') '\r') | (&('\n') '\n') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position489 := position
				depth++
			l490:
				{
					position491, tokenIndex491, depth491 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\f':
							if buffer[position] != rune('\f') {
								goto l491
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l491
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l491
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l491
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l491
							}
							position++
							break
						}
					}

					goto l490
				l491:
					position, tokenIndex, depth = position491, tokenIndex491, depth491
				}
				depth--
				add(rulews, position489)
			}
			return true
		},
		nil,
	}
	p.rules = rules
}
