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
	ruleexpression
	ruleconditionalOrExpression
	ruleconditionalAndExpression
	rulevalueLogical
	rulenumericExpression
	rulemultiplicativeExpression
	ruleunaryExpression
	ruleprimaryExpression
	rulebrackettedExpression
	rulein
	rulenotin
	ruleargList
	rulevar
	ruleiriref
	ruleiri
	ruleprefixedName
	ruleliteral
	rulestring
	rulenumericLiteral
	rulesignedNumericLiteral
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
	rulePLUS
	ruleMINUS
	ruleOPTIONAL
	ruleUNION
	ruleLIMIT
	ruleOFFSET
	ruleINTEGER
	ruleCONSTRUCT
	ruleDESCRIBE
	ruleASK
	ruleOR
	ruleAND
	ruleEQ
	ruleNE
	ruleGT
	ruleLT
	ruleLE
	ruleGE
	ruleIN
	ruleNOTIN
	ruleAS
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
	"expression",
	"conditionalOrExpression",
	"conditionalAndExpression",
	"valueLogical",
	"numericExpression",
	"multiplicativeExpression",
	"unaryExpression",
	"primaryExpression",
	"brackettedExpression",
	"in",
	"notin",
	"argList",
	"var",
	"iriref",
	"iri",
	"prefixedName",
	"literal",
	"string",
	"numericLiteral",
	"signedNumericLiteral",
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
	"PLUS",
	"MINUS",
	"OPTIONAL",
	"UNION",
	"LIMIT",
	"OFFSET",
	"INTEGER",
	"CONSTRUCT",
	"DESCRIBE",
	"ASK",
	"OR",
	"AND",
	"EQ",
	"NE",
	"GT",
	"LT",
	"LE",
	"GE",
	"IN",
	"NOTIN",
	"AS",
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
	rules  [122]func() bool
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
						{
							position165, tokenIndex165, depth165 := position, tokenIndex, depth
							if !rules[rulevar]() {
								goto l166
							}
							goto l165
						l166:
							position, tokenIndex, depth = position165, tokenIndex165, depth165
							if !rules[ruleLPAREN]() {
								goto l109
							}
							if !rules[ruleexpression]() {
								goto l109
							}
							{
								position167 := position
								depth++
								{
									position168, tokenIndex168, depth168 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l169
									}
									position++
									goto l168
								l169:
									position, tokenIndex, depth = position168, tokenIndex168, depth168
									if buffer[position] != rune('A') {
										goto l109
									}
									position++
								}
							l168:
								{
									position170, tokenIndex170, depth170 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l171
									}
									position++
									goto l170
								l171:
									position, tokenIndex, depth = position170, tokenIndex170, depth170
									if buffer[position] != rune('S') {
										goto l109
									}
									position++
								}
							l170:
								if !rules[rulews]() {
									goto l109
								}
								depth--
								add(ruleAS, position167)
							}
							if !rules[rulevar]() {
								goto l109
							}
							if !rules[ruleRPAREN]() {
								goto l109
							}
						}
					l165:
						depth--
						add(ruleprojectionElem, position164)
					}
				l162:
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						{
							position172 := position
							depth++
							{
								position173, tokenIndex173, depth173 := position, tokenIndex, depth
								if !rules[rulevar]() {
									goto l174
								}
								goto l173
							l174:
								position, tokenIndex, depth = position173, tokenIndex173, depth173
								if !rules[ruleLPAREN]() {
									goto l163
								}
								if !rules[ruleexpression]() {
									goto l163
								}
								{
									position175 := position
									depth++
									{
										position176, tokenIndex176, depth176 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l177
										}
										position++
										goto l176
									l177:
										position, tokenIndex, depth = position176, tokenIndex176, depth176
										if buffer[position] != rune('A') {
											goto l163
										}
										position++
									}
								l176:
									{
										position178, tokenIndex178, depth178 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l179
										}
										position++
										goto l178
									l179:
										position, tokenIndex, depth = position178, tokenIndex178, depth178
										if buffer[position] != rune('S') {
											goto l163
										}
										position++
									}
								l178:
									if !rules[rulews]() {
										goto l163
									}
									depth--
									add(ruleAS, position175)
								}
								if !rules[rulevar]() {
									goto l163
								}
								if !rules[ruleRPAREN]() {
									goto l163
								}
							}
						l173:
							depth--
							add(ruleprojectionElem, position172)
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
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				if !rules[ruleselect]() {
					goto l180
				}
				if !rules[rulewhereClause]() {
					goto l180
				}
				depth--
				add(rulesubSelect, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
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
		/* 13 projectionElem <- <(var / (LPAREN expression AS var RPAREN))> */
		nil,
		/* 14 datasetClause <- <(FROM NAMED? iriref)> */
		func() bool {
			position188, tokenIndex188, depth188 := position, tokenIndex, depth
			{
				position189 := position
				depth++
				{
					position190 := position
					depth++
					{
						position191, tokenIndex191, depth191 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l192
						}
						position++
						goto l191
					l192:
						position, tokenIndex, depth = position191, tokenIndex191, depth191
						if buffer[position] != rune('F') {
							goto l188
						}
						position++
					}
				l191:
					{
						position193, tokenIndex193, depth193 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l194
						}
						position++
						goto l193
					l194:
						position, tokenIndex, depth = position193, tokenIndex193, depth193
						if buffer[position] != rune('R') {
							goto l188
						}
						position++
					}
				l193:
					{
						position195, tokenIndex195, depth195 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l196
						}
						position++
						goto l195
					l196:
						position, tokenIndex, depth = position195, tokenIndex195, depth195
						if buffer[position] != rune('O') {
							goto l188
						}
						position++
					}
				l195:
					{
						position197, tokenIndex197, depth197 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l198
						}
						position++
						goto l197
					l198:
						position, tokenIndex, depth = position197, tokenIndex197, depth197
						if buffer[position] != rune('M') {
							goto l188
						}
						position++
					}
				l197:
					if !rules[rulews]() {
						goto l188
					}
					depth--
					add(ruleFROM, position190)
				}
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					{
						position201 := position
						depth++
						{
							position202, tokenIndex202, depth202 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l203
							}
							position++
							goto l202
						l203:
							position, tokenIndex, depth = position202, tokenIndex202, depth202
							if buffer[position] != rune('N') {
								goto l199
							}
							position++
						}
					l202:
						{
							position204, tokenIndex204, depth204 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l205
							}
							position++
							goto l204
						l205:
							position, tokenIndex, depth = position204, tokenIndex204, depth204
							if buffer[position] != rune('A') {
								goto l199
							}
							position++
						}
					l204:
						{
							position206, tokenIndex206, depth206 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l207
							}
							position++
							goto l206
						l207:
							position, tokenIndex, depth = position206, tokenIndex206, depth206
							if buffer[position] != rune('M') {
								goto l199
							}
							position++
						}
					l206:
						{
							position208, tokenIndex208, depth208 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l209
							}
							position++
							goto l208
						l209:
							position, tokenIndex, depth = position208, tokenIndex208, depth208
							if buffer[position] != rune('E') {
								goto l199
							}
							position++
						}
					l208:
						{
							position210, tokenIndex210, depth210 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l211
							}
							position++
							goto l210
						l211:
							position, tokenIndex, depth = position210, tokenIndex210, depth210
							if buffer[position] != rune('D') {
								goto l199
							}
							position++
						}
					l210:
						if !rules[rulews]() {
							goto l199
						}
						depth--
						add(ruleNAMED, position201)
					}
					goto l200
				l199:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
				}
			l200:
				if !rules[ruleiriref]() {
					goto l188
				}
				depth--
				add(ruledatasetClause, position189)
			}
			return true
		l188:
			position, tokenIndex, depth = position188, tokenIndex188, depth188
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position212, tokenIndex212, depth212 := position, tokenIndex, depth
			{
				position213 := position
				depth++
				{
					position214, tokenIndex214, depth214 := position, tokenIndex, depth
					{
						position216 := position
						depth++
						{
							position217, tokenIndex217, depth217 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l218
							}
							position++
							goto l217
						l218:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if buffer[position] != rune('W') {
								goto l214
							}
							position++
						}
					l217:
						{
							position219, tokenIndex219, depth219 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l220
							}
							position++
							goto l219
						l220:
							position, tokenIndex, depth = position219, tokenIndex219, depth219
							if buffer[position] != rune('H') {
								goto l214
							}
							position++
						}
					l219:
						{
							position221, tokenIndex221, depth221 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l222
							}
							position++
							goto l221
						l222:
							position, tokenIndex, depth = position221, tokenIndex221, depth221
							if buffer[position] != rune('E') {
								goto l214
							}
							position++
						}
					l221:
						{
							position223, tokenIndex223, depth223 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l224
							}
							position++
							goto l223
						l224:
							position, tokenIndex, depth = position223, tokenIndex223, depth223
							if buffer[position] != rune('R') {
								goto l214
							}
							position++
						}
					l223:
						{
							position225, tokenIndex225, depth225 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l226
							}
							position++
							goto l225
						l226:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if buffer[position] != rune('E') {
								goto l214
							}
							position++
						}
					l225:
						if !rules[rulews]() {
							goto l214
						}
						depth--
						add(ruleWHERE, position216)
					}
					goto l215
				l214:
					position, tokenIndex, depth = position214, tokenIndex214, depth214
				}
			l215:
				if !rules[rulegroupGraphPattern]() {
					goto l212
				}
				depth--
				add(rulewhereClause, position213)
			}
			return true
		l212:
			position, tokenIndex, depth = position212, tokenIndex212, depth212
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position227, tokenIndex227, depth227 := position, tokenIndex, depth
			{
				position228 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l227
				}
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l230
					}
					goto l229
				l230:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
					if !rules[rulegraphPattern]() {
						goto l227
					}
				}
			l229:
				if !rules[ruleRBRACE]() {
					goto l227
				}
				depth--
				add(rulegroupGraphPattern, position228)
			}
			return true
		l227:
			position, tokenIndex, depth = position227, tokenIndex227, depth227
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position232 := position
				depth++
				{
					position233, tokenIndex233, depth233 := position, tokenIndex, depth
					{
						position235 := position
						depth++
						if !rules[ruletriplesBlock]() {
							goto l233
						}
						depth--
						add(rulebasicGraphPattern, position235)
					}
					goto l234
				l233:
					position, tokenIndex, depth = position233, tokenIndex233, depth233
				}
			l234:
				{
					position236, tokenIndex236, depth236 := position, tokenIndex, depth
					{
						position238 := position
						depth++
						{
							position239, tokenIndex239, depth239 := position, tokenIndex, depth
							{
								position241 := position
								depth++
								{
									position242 := position
									depth++
									{
										position243, tokenIndex243, depth243 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l244
										}
										position++
										goto l243
									l244:
										position, tokenIndex, depth = position243, tokenIndex243, depth243
										if buffer[position] != rune('O') {
											goto l240
										}
										position++
									}
								l243:
									{
										position245, tokenIndex245, depth245 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l246
										}
										position++
										goto l245
									l246:
										position, tokenIndex, depth = position245, tokenIndex245, depth245
										if buffer[position] != rune('P') {
											goto l240
										}
										position++
									}
								l245:
									{
										position247, tokenIndex247, depth247 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l248
										}
										position++
										goto l247
									l248:
										position, tokenIndex, depth = position247, tokenIndex247, depth247
										if buffer[position] != rune('T') {
											goto l240
										}
										position++
									}
								l247:
									{
										position249, tokenIndex249, depth249 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l250
										}
										position++
										goto l249
									l250:
										position, tokenIndex, depth = position249, tokenIndex249, depth249
										if buffer[position] != rune('I') {
											goto l240
										}
										position++
									}
								l249:
									{
										position251, tokenIndex251, depth251 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l252
										}
										position++
										goto l251
									l252:
										position, tokenIndex, depth = position251, tokenIndex251, depth251
										if buffer[position] != rune('O') {
											goto l240
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
											goto l240
										}
										position++
									}
								l253:
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l256
										}
										position++
										goto l255
									l256:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
										if buffer[position] != rune('A') {
											goto l240
										}
										position++
									}
								l255:
									{
										position257, tokenIndex257, depth257 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l258
										}
										position++
										goto l257
									l258:
										position, tokenIndex, depth = position257, tokenIndex257, depth257
										if buffer[position] != rune('L') {
											goto l240
										}
										position++
									}
								l257:
									if !rules[rulews]() {
										goto l240
									}
									depth--
									add(ruleOPTIONAL, position242)
								}
								if !rules[ruleLBRACE]() {
									goto l240
								}
								{
									position259, tokenIndex259, depth259 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l260
									}
									goto l259
								l260:
									position, tokenIndex, depth = position259, tokenIndex259, depth259
									if !rules[rulegraphPattern]() {
										goto l240
									}
								}
							l259:
								if !rules[ruleRBRACE]() {
									goto l240
								}
								depth--
								add(ruleoptionalGraphPattern, position241)
							}
							goto l239
						l240:
							position, tokenIndex, depth = position239, tokenIndex239, depth239
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l236
							}
						}
					l239:
						depth--
						add(rulegraphPatternNotTriples, position238)
					}
					{
						position261, tokenIndex261, depth261 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l261
						}
						goto l262
					l261:
						position, tokenIndex, depth = position261, tokenIndex261, depth261
					}
				l262:
					if !rules[rulegraphPattern]() {
						goto l236
					}
					goto l237
				l236:
					position, tokenIndex, depth = position236, tokenIndex236, depth236
				}
			l237:
				depth--
				add(rulegraphPattern, position232)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position265, tokenIndex265, depth265 := position, tokenIndex, depth
			{
				position266 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l265
				}
				{
					position267, tokenIndex267, depth267 := position, tokenIndex, depth
					{
						position269 := position
						depth++
						{
							position270, tokenIndex270, depth270 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l271
							}
							position++
							goto l270
						l271:
							position, tokenIndex, depth = position270, tokenIndex270, depth270
							if buffer[position] != rune('U') {
								goto l267
							}
							position++
						}
					l270:
						{
							position272, tokenIndex272, depth272 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l273
							}
							position++
							goto l272
						l273:
							position, tokenIndex, depth = position272, tokenIndex272, depth272
							if buffer[position] != rune('N') {
								goto l267
							}
							position++
						}
					l272:
						{
							position274, tokenIndex274, depth274 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l275
							}
							position++
							goto l274
						l275:
							position, tokenIndex, depth = position274, tokenIndex274, depth274
							if buffer[position] != rune('I') {
								goto l267
							}
							position++
						}
					l274:
						{
							position276, tokenIndex276, depth276 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l277
							}
							position++
							goto l276
						l277:
							position, tokenIndex, depth = position276, tokenIndex276, depth276
							if buffer[position] != rune('O') {
								goto l267
							}
							position++
						}
					l276:
						{
							position278, tokenIndex278, depth278 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l279
							}
							position++
							goto l278
						l279:
							position, tokenIndex, depth = position278, tokenIndex278, depth278
							if buffer[position] != rune('N') {
								goto l267
							}
							position++
						}
					l278:
						if !rules[rulews]() {
							goto l267
						}
						depth--
						add(ruleUNION, position269)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l267
					}
					goto l268
				l267:
					position, tokenIndex, depth = position267, tokenIndex267, depth267
				}
			l268:
				depth--
				add(rulegroupOrUnionGraphPattern, position266)
			}
			return true
		l265:
			position, tokenIndex, depth = position265, tokenIndex265, depth265
			return false
		},
		/* 21 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 22 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position281, tokenIndex281, depth281 := position, tokenIndex, depth
			{
				position282 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l281
				}
			l283:
				{
					position284, tokenIndex284, depth284 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l284
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l284
					}
					goto l283
				l284:
					position, tokenIndex, depth = position284, tokenIndex284, depth284
				}
				{
					position285, tokenIndex285, depth285 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l285
					}
					goto l286
				l285:
					position, tokenIndex, depth = position285, tokenIndex285, depth285
				}
			l286:
				depth--
				add(ruletriplesBlock, position282)
			}
			return true
		l281:
			position, tokenIndex, depth = position281, tokenIndex281, depth281
			return false
		},
		/* 23 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position287, tokenIndex287, depth287 := position, tokenIndex, depth
			{
				position288 := position
				depth++
				{
					position289, tokenIndex289, depth289 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l290
					}
					if !rules[rulepropertyListPath]() {
						goto l290
					}
					goto l289
				l290:
					position, tokenIndex, depth = position289, tokenIndex289, depth289
					{
						position291 := position
						depth++
						{
							position292, tokenIndex292, depth292 := position, tokenIndex, depth
							{
								position294 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l293
								}
								if !rules[rulegraphNodePath]() {
									goto l293
								}
							l295:
								{
									position296, tokenIndex296, depth296 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l296
									}
									goto l295
								l296:
									position, tokenIndex, depth = position296, tokenIndex296, depth296
								}
								if !rules[ruleRPAREN]() {
									goto l293
								}
								depth--
								add(rulecollectionPath, position294)
							}
							goto l292
						l293:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
							{
								position297 := position
								depth++
								{
									position298 := position
									depth++
									if buffer[position] != rune('[') {
										goto l287
									}
									position++
									if !rules[rulews]() {
										goto l287
									}
									depth--
									add(ruleLBRACK, position298)
								}
								if !rules[rulepropertyListPath]() {
									goto l287
								}
								{
									position299 := position
									depth++
									if buffer[position] != rune(']') {
										goto l287
									}
									position++
									if !rules[rulews]() {
										goto l287
									}
									depth--
									add(ruleRBRACK, position299)
								}
								depth--
								add(ruleblankNodePropertyListPath, position297)
							}
						}
					l292:
						depth--
						add(ruletriplesNodePath, position291)
					}
					if !rules[rulepropertyListPath]() {
						goto l287
					}
				}
			l289:
				depth--
				add(ruletriplesSameSubjectPath, position288)
			}
			return true
		l287:
			position, tokenIndex, depth = position287, tokenIndex287, depth287
			return false
		},
		/* 24 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position300, tokenIndex300, depth300 := position, tokenIndex, depth
			{
				position301 := position
				depth++
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l303
					}
					goto l302
				l303:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					{
						position304 := position
						depth++
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l306
							}
							goto l305
						l306:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
							{
								switch buffer[position] {
								case '(':
									if !rules[rulenil]() {
										goto l300
									}
									break
								case '[', '_':
									{
										position308 := position
										depth++
										{
											position309, tokenIndex309, depth309 := position, tokenIndex, depth
											{
												position311 := position
												depth++
												if buffer[position] != rune('_') {
													goto l310
												}
												position++
												if buffer[position] != rune(':') {
													goto l310
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l310
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l310
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l310
														}
														position++
														break
													}
												}

												{
													position313, tokenIndex313, depth313 := position, tokenIndex, depth
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
															goto l317
														}
														position++
														goto l315
													l317:
														position, tokenIndex, depth = position315, tokenIndex315, depth315
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l318
														}
														position++
														goto l315
													l318:
														position, tokenIndex, depth = position315, tokenIndex315, depth315
														if c := buffer[position]; c < rune('.') || c > rune('_') {
															goto l313
														}
														position++
													}
												l315:
													goto l314
												l313:
													position, tokenIndex, depth = position313, tokenIndex313, depth313
												}
											l314:
												if !rules[rulews]() {
													goto l310
												}
												depth--
												add(ruleblankNodeLabel, position311)
											}
											goto l309
										l310:
											position, tokenIndex, depth = position309, tokenIndex309, depth309
											{
												position319 := position
												depth++
												if buffer[position] != rune('[') {
													goto l300
												}
												position++
												if !rules[rulews]() {
													goto l300
												}
												if buffer[position] != rune(']') {
													goto l300
												}
												position++
												if !rules[rulews]() {
													goto l300
												}
												depth--
												add(ruleanon, position319)
											}
										}
									l309:
										depth--
										add(ruleblankNode, position308)
									}
									break
								case 'F', 'T', 'f', 't':
									if !rules[rulebooleanLiteral]() {
										goto l300
									}
									break
								case '"':
									if !rules[ruleliteral]() {
										goto l300
									}
									break
								default:
									if !rules[rulenumericLiteral]() {
										goto l300
									}
									break
								}
							}

						}
					l305:
						depth--
						add(rulegraphTerm, position304)
					}
				}
			l302:
				depth--
				add(rulevarOrTerm, position301)
			}
			return true
		l300:
			position, tokenIndex, depth = position300, tokenIndex300, depth300
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
			position324, tokenIndex324, depth324 := position, tokenIndex, depth
			{
				position325 := position
				depth++
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l327
					}
					goto l326
				l327:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
					{
						position328 := position
						depth++
						if !rules[rulepath]() {
							goto l324
						}
						depth--
						add(ruleverbPath, position328)
					}
				}
			l326:
				if !rules[ruleobjectListPath]() {
					goto l324
				}
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					{
						position331 := position
						depth++
						if buffer[position] != rune(';') {
							goto l329
						}
						position++
						if !rules[rulews]() {
							goto l329
						}
						depth--
						add(ruleSEMICOLON, position331)
					}
					if !rules[rulepropertyListPath]() {
						goto l329
					}
					goto l330
				l329:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
				}
			l330:
				depth--
				add(rulepropertyListPath, position325)
			}
			return true
		l324:
			position, tokenIndex, depth = position324, tokenIndex324, depth324
			return false
		},
		/* 30 verbPath <- <path> */
		nil,
		/* 31 path <- <pathAlternative> */
		func() bool {
			position333, tokenIndex333, depth333 := position, tokenIndex, depth
			{
				position334 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l333
				}
				depth--
				add(rulepath, position334)
			}
			return true
		l333:
			position, tokenIndex, depth = position333, tokenIndex333, depth333
			return false
		},
		/* 32 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position335, tokenIndex335, depth335 := position, tokenIndex, depth
			{
				position336 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l335
				}
			l337:
				{
					position338, tokenIndex338, depth338 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l338
					}
					if !rules[rulepathAlternative]() {
						goto l338
					}
					goto l337
				l338:
					position, tokenIndex, depth = position338, tokenIndex338, depth338
				}
				depth--
				add(rulepathAlternative, position336)
			}
			return true
		l335:
			position, tokenIndex, depth = position335, tokenIndex335, depth335
			return false
		},
		/* 33 pathSequence <- <(pathElt (SLASH pathSequence)*)> */
		func() bool {
			position339, tokenIndex339, depth339 := position, tokenIndex, depth
			{
				position340 := position
				depth++
				{
					position341 := position
					depth++
					{
						position342, tokenIndex342, depth342 := position, tokenIndex, depth
						if !rules[ruleINVERSE]() {
							goto l342
						}
						goto l343
					l342:
						position, tokenIndex, depth = position342, tokenIndex342, depth342
					}
				l343:
					{
						position344 := position
						depth++
						{
							position345, tokenIndex345, depth345 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l346
							}
							goto l345
						l346:
							position, tokenIndex, depth = position345, tokenIndex345, depth345
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l339
									}
									if !rules[rulepath]() {
										goto l339
									}
									if !rules[ruleRPAREN]() {
										goto l339
									}
									break
								case '!':
									if !rules[ruleNOT]() {
										goto l339
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
												goto l339
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
												goto l339
											}
										}
									l349:
										depth--
										add(rulepathNegatedPropertySet, position348)
									}
									break
								default:
									if !rules[ruleISA]() {
										goto l339
									}
									break
								}
							}

						}
					l345:
						depth--
						add(rulepathPrimary, position344)
					}
					depth--
					add(rulepathElt, position341)
				}
			l355:
				{
					position356, tokenIndex356, depth356 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l356
					}
					if !rules[rulepathSequence]() {
						goto l356
					}
					goto l355
				l356:
					position, tokenIndex, depth = position356, tokenIndex356, depth356
				}
				depth--
				add(rulepathSequence, position340)
			}
			return true
		l339:
			position, tokenIndex, depth = position339, tokenIndex339, depth339
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
			position360, tokenIndex360, depth360 := position, tokenIndex, depth
			{
				position361 := position
				depth++
				{
					position362, tokenIndex362, depth362 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l363
					}
					goto l362
				l363:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
					if !rules[ruleISA]() {
						goto l364
					}
					goto l362
				l364:
					position, tokenIndex, depth = position362, tokenIndex362, depth362
					if !rules[ruleINVERSE]() {
						goto l360
					}
					{
						position365, tokenIndex365, depth365 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l366
						}
						goto l365
					l366:
						position, tokenIndex, depth = position365, tokenIndex365, depth365
						if !rules[ruleISA]() {
							goto l360
						}
					}
				l365:
				}
			l362:
				depth--
				add(rulepathOneInPropertySet, position361)
			}
			return true
		l360:
			position, tokenIndex, depth = position360, tokenIndex360, depth360
			return false
		},
		/* 38 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position367, tokenIndex367, depth367 := position, tokenIndex, depth
			{
				position368 := position
				depth++
				{
					position369 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l367
					}
					depth--
					add(ruleobjectPath, position369)
				}
			l370:
				{
					position371, tokenIndex371, depth371 := position, tokenIndex, depth
					if !rules[ruleCOMMA]() {
						goto l371
					}
					if !rules[ruleobjectListPath]() {
						goto l371
					}
					goto l370
				l371:
					position, tokenIndex, depth = position371, tokenIndex371, depth371
				}
				depth--
				add(ruleobjectListPath, position368)
			}
			return true
		l367:
			position, tokenIndex, depth = position367, tokenIndex367, depth367
			return false
		},
		/* 39 objectPath <- <graphNodePath> */
		nil,
		/* 40 graphNodePath <- <varOrTerm> */
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
		/* 41 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position376 := position
				depth++
				{
					position377, tokenIndex377, depth377 := position, tokenIndex, depth
					{
						position379 := position
						depth++
						{
							position380, tokenIndex380, depth380 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l381
							}
							{
								position382, tokenIndex382, depth382 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l382
								}
								goto l383
							l382:
								position, tokenIndex, depth = position382, tokenIndex382, depth382
							}
						l383:
							goto l380
						l381:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
							if !rules[ruleoffset]() {
								goto l377
							}
							{
								position384, tokenIndex384, depth384 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l384
								}
								goto l385
							l384:
								position, tokenIndex, depth = position384, tokenIndex384, depth384
							}
						l385:
						}
					l380:
						depth--
						add(rulelimitOffsetClauses, position379)
					}
					goto l378
				l377:
					position, tokenIndex, depth = position377, tokenIndex377, depth377
				}
			l378:
				depth--
				add(rulesolutionModifier, position376)
			}
			return true
		},
		/* 42 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 43 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position387, tokenIndex387, depth387 := position, tokenIndex, depth
			{
				position388 := position
				depth++
				{
					position389 := position
					depth++
					{
						position390, tokenIndex390, depth390 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l391
						}
						position++
						goto l390
					l391:
						position, tokenIndex, depth = position390, tokenIndex390, depth390
						if buffer[position] != rune('L') {
							goto l387
						}
						position++
					}
				l390:
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l393
						}
						position++
						goto l392
					l393:
						position, tokenIndex, depth = position392, tokenIndex392, depth392
						if buffer[position] != rune('I') {
							goto l387
						}
						position++
					}
				l392:
					{
						position394, tokenIndex394, depth394 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l395
						}
						position++
						goto l394
					l395:
						position, tokenIndex, depth = position394, tokenIndex394, depth394
						if buffer[position] != rune('M') {
							goto l387
						}
						position++
					}
				l394:
					{
						position396, tokenIndex396, depth396 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l397
						}
						position++
						goto l396
					l397:
						position, tokenIndex, depth = position396, tokenIndex396, depth396
						if buffer[position] != rune('I') {
							goto l387
						}
						position++
					}
				l396:
					{
						position398, tokenIndex398, depth398 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l399
						}
						position++
						goto l398
					l399:
						position, tokenIndex, depth = position398, tokenIndex398, depth398
						if buffer[position] != rune('T') {
							goto l387
						}
						position++
					}
				l398:
					if !rules[rulews]() {
						goto l387
					}
					depth--
					add(ruleLIMIT, position389)
				}
				if !rules[ruleINTEGER]() {
					goto l387
				}
				depth--
				add(rulelimit, position388)
			}
			return true
		l387:
			position, tokenIndex, depth = position387, tokenIndex387, depth387
			return false
		},
		/* 44 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position400, tokenIndex400, depth400 := position, tokenIndex, depth
			{
				position401 := position
				depth++
				{
					position402 := position
					depth++
					{
						position403, tokenIndex403, depth403 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l404
						}
						position++
						goto l403
					l404:
						position, tokenIndex, depth = position403, tokenIndex403, depth403
						if buffer[position] != rune('O') {
							goto l400
						}
						position++
					}
				l403:
					{
						position405, tokenIndex405, depth405 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l406
						}
						position++
						goto l405
					l406:
						position, tokenIndex, depth = position405, tokenIndex405, depth405
						if buffer[position] != rune('F') {
							goto l400
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
							goto l400
						}
						position++
					}
				l407:
					{
						position409, tokenIndex409, depth409 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l410
						}
						position++
						goto l409
					l410:
						position, tokenIndex, depth = position409, tokenIndex409, depth409
						if buffer[position] != rune('S') {
							goto l400
						}
						position++
					}
				l409:
					{
						position411, tokenIndex411, depth411 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l412
						}
						position++
						goto l411
					l412:
						position, tokenIndex, depth = position411, tokenIndex411, depth411
						if buffer[position] != rune('E') {
							goto l400
						}
						position++
					}
				l411:
					{
						position413, tokenIndex413, depth413 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l414
						}
						position++
						goto l413
					l414:
						position, tokenIndex, depth = position413, tokenIndex413, depth413
						if buffer[position] != rune('T') {
							goto l400
						}
						position++
					}
				l413:
					if !rules[rulews]() {
						goto l400
					}
					depth--
					add(ruleOFFSET, position402)
				}
				if !rules[ruleINTEGER]() {
					goto l400
				}
				depth--
				add(ruleoffset, position401)
			}
			return true
		l400:
			position, tokenIndex, depth = position400, tokenIndex400, depth400
			return false
		},
		/* 45 expression <- <conditionalOrExpression> */
		func() bool {
			position415, tokenIndex415, depth415 := position, tokenIndex, depth
			{
				position416 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l415
				}
				depth--
				add(ruleexpression, position416)
			}
			return true
		l415:
			position, tokenIndex, depth = position415, tokenIndex415, depth415
			return false
		},
		/* 46 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l417
				}
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					{
						position421 := position
						depth++
						if buffer[position] != rune('|') {
							goto l419
						}
						position++
						if buffer[position] != rune('|') {
							goto l419
						}
						position++
						if !rules[rulews]() {
							goto l419
						}
						depth--
						add(ruleOR, position421)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l419
					}
					goto l420
				l419:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
				}
			l420:
				depth--
				add(ruleconditionalOrExpression, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 47 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position422, tokenIndex422, depth422 := position, tokenIndex, depth
			{
				position423 := position
				depth++
				{
					position424 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l422
					}
					{
						position425, tokenIndex425, depth425 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position428 := position
									depth++
									{
										position429 := position
										depth++
										{
											position430, tokenIndex430, depth430 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l431
											}
											position++
											goto l430
										l431:
											position, tokenIndex, depth = position430, tokenIndex430, depth430
											if buffer[position] != rune('N') {
												goto l425
											}
											position++
										}
									l430:
										{
											position432, tokenIndex432, depth432 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l433
											}
											position++
											goto l432
										l433:
											position, tokenIndex, depth = position432, tokenIndex432, depth432
											if buffer[position] != rune('O') {
												goto l425
											}
											position++
										}
									l432:
										{
											position434, tokenIndex434, depth434 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l435
											}
											position++
											goto l434
										l435:
											position, tokenIndex, depth = position434, tokenIndex434, depth434
											if buffer[position] != rune('T') {
												goto l425
											}
											position++
										}
									l434:
										if buffer[position] != rune(' ') {
											goto l425
										}
										position++
										{
											position436, tokenIndex436, depth436 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l437
											}
											position++
											goto l436
										l437:
											position, tokenIndex, depth = position436, tokenIndex436, depth436
											if buffer[position] != rune('I') {
												goto l425
											}
											position++
										}
									l436:
										{
											position438, tokenIndex438, depth438 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l439
											}
											position++
											goto l438
										l439:
											position, tokenIndex, depth = position438, tokenIndex438, depth438
											if buffer[position] != rune('N') {
												goto l425
											}
											position++
										}
									l438:
										if !rules[rulews]() {
											goto l425
										}
										depth--
										add(ruleNOTIN, position429)
									}
									if !rules[ruleargList]() {
										goto l425
									}
									depth--
									add(rulenotin, position428)
								}
								break
							case 'I', 'i':
								{
									position440 := position
									depth++
									{
										position441 := position
										depth++
										{
											position442, tokenIndex442, depth442 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l443
											}
											position++
											goto l442
										l443:
											position, tokenIndex, depth = position442, tokenIndex442, depth442
											if buffer[position] != rune('I') {
												goto l425
											}
											position++
										}
									l442:
										{
											position444, tokenIndex444, depth444 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l445
											}
											position++
											goto l444
										l445:
											position, tokenIndex, depth = position444, tokenIndex444, depth444
											if buffer[position] != rune('N') {
												goto l425
											}
											position++
										}
									l444:
										if !rules[rulews]() {
											goto l425
										}
										depth--
										add(ruleIN, position441)
									}
									if !rules[ruleargList]() {
										goto l425
									}
									depth--
									add(rulein, position440)
								}
								break
							default:
								{
									position446, tokenIndex446, depth446 := position, tokenIndex, depth
									{
										position448 := position
										depth++
										if buffer[position] != rune('<') {
											goto l447
										}
										position++
										if !rules[rulews]() {
											goto l447
										}
										depth--
										add(ruleLT, position448)
									}
									goto l446
								l447:
									position, tokenIndex, depth = position446, tokenIndex446, depth446
									{
										position450 := position
										depth++
										if buffer[position] != rune('>') {
											goto l449
										}
										position++
										if buffer[position] != rune('=') {
											goto l449
										}
										position++
										if !rules[rulews]() {
											goto l449
										}
										depth--
										add(ruleGE, position450)
									}
									goto l446
								l449:
									position, tokenIndex, depth = position446, tokenIndex446, depth446
									{
										switch buffer[position] {
										case '>':
											{
												position452 := position
												depth++
												if buffer[position] != rune('>') {
													goto l425
												}
												position++
												if !rules[rulews]() {
													goto l425
												}
												depth--
												add(ruleGT, position452)
											}
											break
										case '<':
											{
												position453 := position
												depth++
												if buffer[position] != rune('<') {
													goto l425
												}
												position++
												if buffer[position] != rune('=') {
													goto l425
												}
												position++
												if !rules[rulews]() {
													goto l425
												}
												depth--
												add(ruleLE, position453)
											}
											break
										case '!':
											{
												position454 := position
												depth++
												if buffer[position] != rune('!') {
													goto l425
												}
												position++
												if buffer[position] != rune('=') {
													goto l425
												}
												position++
												if !rules[rulews]() {
													goto l425
												}
												depth--
												add(ruleNE, position454)
											}
											break
										default:
											{
												position455 := position
												depth++
												if buffer[position] != rune('=') {
													goto l425
												}
												position++
												if !rules[rulews]() {
													goto l425
												}
												depth--
												add(ruleEQ, position455)
											}
											break
										}
									}

								}
							l446:
								if !rules[rulenumericExpression]() {
									goto l425
								}
								break
							}
						}

						goto l426
					l425:
						position, tokenIndex, depth = position425, tokenIndex425, depth425
					}
				l426:
					depth--
					add(rulevalueLogical, position424)
				}
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					{
						position458 := position
						depth++
						if buffer[position] != rune('&') {
							goto l456
						}
						position++
						if buffer[position] != rune('&') {
							goto l456
						}
						position++
						if !rules[rulews]() {
							goto l456
						}
						depth--
						add(ruleAND, position458)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l456
					}
					goto l457
				l456:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
				}
			l457:
				depth--
				add(ruleconditionalAndExpression, position423)
			}
			return true
		l422:
			position, tokenIndex, depth = position422, tokenIndex422, depth422
			return false
		},
		/* 48 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 49 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position460, tokenIndex460, depth460 := position, tokenIndex, depth
			{
				position461 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l460
				}
			l462:
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					{
						position464, tokenIndex464, depth464 := position, tokenIndex, depth
						{
							position466, tokenIndex466, depth466 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l467
							}
							goto l466
						l467:
							position, tokenIndex, depth = position466, tokenIndex466, depth466
							if !rules[ruleMINUS]() {
								goto l465
							}
						}
					l466:
						if !rules[rulemultiplicativeExpression]() {
							goto l465
						}
						goto l464
					l465:
						position, tokenIndex, depth = position464, tokenIndex464, depth464
						{
							position468 := position
							depth++
							{
								position469, tokenIndex469, depth469 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l470
								}
								position++
								goto l469
							l470:
								position, tokenIndex, depth = position469, tokenIndex469, depth469
								if buffer[position] != rune('-') {
									goto l463
								}
								position++
							}
						l469:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l463
							}
							position++
						l471:
							{
								position472, tokenIndex472, depth472 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l472
								}
								position++
								goto l471
							l472:
								position, tokenIndex, depth = position472, tokenIndex472, depth472
							}
							{
								position473, tokenIndex473, depth473 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l473
								}
								position++
							l475:
								{
									position476, tokenIndex476, depth476 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l476
									}
									position++
									goto l475
								l476:
									position, tokenIndex, depth = position476, tokenIndex476, depth476
								}
								goto l474
							l473:
								position, tokenIndex, depth = position473, tokenIndex473, depth473
							}
						l474:
							if !rules[rulews]() {
								goto l463
							}
							depth--
							add(rulesignedNumericLiteral, position468)
						}
					}
				l464:
					goto l462
				l463:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
				}
				depth--
				add(rulenumericExpression, position461)
			}
			return true
		l460:
			position, tokenIndex, depth = position460, tokenIndex460, depth460
			return false
		},
		/* 50 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position477, tokenIndex477, depth477 := position, tokenIndex, depth
			{
				position478 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l477
				}
			l479:
				{
					position480, tokenIndex480, depth480 := position, tokenIndex, depth
					{
						position481, tokenIndex481, depth481 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l482
						}
						goto l481
					l482:
						position, tokenIndex, depth = position481, tokenIndex481, depth481
						if !rules[ruleSLASH]() {
							goto l480
						}
					}
				l481:
					if !rules[ruleunaryExpression]() {
						goto l480
					}
					goto l479
				l480:
					position, tokenIndex, depth = position480, tokenIndex480, depth480
				}
				depth--
				add(rulemultiplicativeExpression, position478)
			}
			return true
		l477:
			position, tokenIndex, depth = position477, tokenIndex477, depth477
			return false
		},
		/* 51 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position483, tokenIndex483, depth483 := position, tokenIndex, depth
			{
				position484 := position
				depth++
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l485
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l485
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l485
							}
							break
						}
					}

					goto l486
				l485:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
				}
			l486:
				{
					position488 := position
					depth++
					{
						position489, tokenIndex489, depth489 := position, tokenIndex, depth
						{
							position491 := position
							depth++
							if !rules[ruleLPAREN]() {
								goto l490
							}
							if !rules[ruleexpression]() {
								goto l490
							}
							if !rules[ruleRPAREN]() {
								goto l490
							}
							depth--
							add(rulebrackettedExpression, position491)
						}
						goto l489
					l490:
						position, tokenIndex, depth = position489, tokenIndex489, depth489
						if !rules[ruleiriref]() {
							goto l492
						}
						goto l489
					l492:
						position, tokenIndex, depth = position489, tokenIndex489, depth489
						{
							switch buffer[position] {
							case '$', '?':
								if !rules[rulevar]() {
									goto l483
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l483
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l483
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l483
								}
								break
							}
						}

					}
				l489:
					depth--
					add(ruleprimaryExpression, position488)
				}
				depth--
				add(ruleunaryExpression, position484)
			}
			return true
		l483:
			position, tokenIndex, depth = position483, tokenIndex483, depth483
			return false
		},
		/* 52 primaryExpression <- <(brackettedExpression / iriref / ((&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 53 brackettedExpression <- <(LPAREN expression RPAREN)> */
		nil,
		/* 54 in <- <(IN argList)> */
		nil,
		/* 55 notin <- <(NOTIN argList)> */
		nil,
		/* 56 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position498, tokenIndex498, depth498 := position, tokenIndex, depth
			{
				position499 := position
				depth++
				{
					position500, tokenIndex500, depth500 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l501
					}
					goto l500
				l501:
					position, tokenIndex, depth = position500, tokenIndex500, depth500
					if !rules[ruleLPAREN]() {
						goto l498
					}
					if !rules[ruleexpression]() {
						goto l498
					}
				l502:
					{
						position503, tokenIndex503, depth503 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l503
						}
						if !rules[ruleexpression]() {
							goto l503
						}
						goto l502
					l503:
						position, tokenIndex, depth = position503, tokenIndex503, depth503
					}
					if !rules[ruleRPAREN]() {
						goto l498
					}
				}
			l500:
				depth--
				add(ruleargList, position499)
			}
			return true
		l498:
			position, tokenIndex, depth = position498, tokenIndex498, depth498
			return false
		},
		/* 57 var <- <(('?' / '$') VARNAME ws)> */
		func() bool {
			position504, tokenIndex504, depth504 := position, tokenIndex, depth
			{
				position505 := position
				depth++
				{
					position506, tokenIndex506, depth506 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l507
					}
					position++
					goto l506
				l507:
					position, tokenIndex, depth = position506, tokenIndex506, depth506
					if buffer[position] != rune('$') {
						goto l504
					}
					position++
				}
			l506:
				{
					position508 := position
					depth++
					{
						position511, tokenIndex511, depth511 := position, tokenIndex, depth
						{
							position513 := position
							depth++
							{
								position514, tokenIndex514, depth514 := position, tokenIndex, depth
								{
									position516 := position
									depth++
									{
										position517, tokenIndex517, depth517 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l518
										}
										position++
										goto l517
									l518:
										position, tokenIndex, depth = position517, tokenIndex517, depth517
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l515
										}
										position++
									}
								l517:
									depth--
									add(rulePN_CHARS_BASE, position516)
								}
								goto l514
							l515:
								position, tokenIndex, depth = position514, tokenIndex514, depth514
								if buffer[position] != rune('_') {
									goto l512
								}
								position++
							}
						l514:
							depth--
							add(rulePN_CHARS_U, position513)
						}
						goto l511
					l512:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l504
						}
						position++
					}
				l511:
				l509:
					{
						position510, tokenIndex510, depth510 := position, tokenIndex, depth
						{
							position519, tokenIndex519, depth519 := position, tokenIndex, depth
							{
								position521 := position
								depth++
								{
									position522, tokenIndex522, depth522 := position, tokenIndex, depth
									{
										position524 := position
										depth++
										{
											position525, tokenIndex525, depth525 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l526
											}
											position++
											goto l525
										l526:
											position, tokenIndex, depth = position525, tokenIndex525, depth525
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l523
											}
											position++
										}
									l525:
										depth--
										add(rulePN_CHARS_BASE, position524)
									}
									goto l522
								l523:
									position, tokenIndex, depth = position522, tokenIndex522, depth522
									if buffer[position] != rune('_') {
										goto l520
									}
									position++
								}
							l522:
								depth--
								add(rulePN_CHARS_U, position521)
							}
							goto l519
						l520:
							position, tokenIndex, depth = position519, tokenIndex519, depth519
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l510
							}
							position++
						}
					l519:
						goto l509
					l510:
						position, tokenIndex, depth = position510, tokenIndex510, depth510
					}
					depth--
					add(ruleVARNAME, position508)
				}
				if !rules[rulews]() {
					goto l504
				}
				depth--
				add(rulevar, position505)
			}
			return true
		l504:
			position, tokenIndex, depth = position504, tokenIndex504, depth504
			return false
		},
		/* 58 iriref <- <(iri / prefixedName)> */
		func() bool {
			position527, tokenIndex527, depth527 := position, tokenIndex, depth
			{
				position528 := position
				depth++
				{
					position529, tokenIndex529, depth529 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l530
					}
					goto l529
				l530:
					position, tokenIndex, depth = position529, tokenIndex529, depth529
					{
						position531 := position
						depth++
					l532:
						{
							position533, tokenIndex533, depth533 := position, tokenIndex, depth
							{
								position534, tokenIndex534, depth534 := position, tokenIndex, depth
								{
									position535, tokenIndex535, depth535 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l536
									}
									position++
									goto l535
								l536:
									position, tokenIndex, depth = position535, tokenIndex535, depth535
									if buffer[position] != rune(' ') {
										goto l534
									}
									position++
								}
							l535:
								goto l533
							l534:
								position, tokenIndex, depth = position534, tokenIndex534, depth534
							}
							if !matchDot() {
								goto l533
							}
							goto l532
						l533:
							position, tokenIndex, depth = position533, tokenIndex533, depth533
						}
						if buffer[position] != rune(':') {
							goto l527
						}
						position++
					l537:
						{
							position538, tokenIndex538, depth538 := position, tokenIndex, depth
							{
								position539, tokenIndex539, depth539 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l540
								}
								position++
								goto l539
							l540:
								position, tokenIndex, depth = position539, tokenIndex539, depth539
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l541
								}
								position++
								goto l539
							l541:
								position, tokenIndex, depth = position539, tokenIndex539, depth539
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l542
								}
								position++
								goto l539
							l542:
								position, tokenIndex, depth = position539, tokenIndex539, depth539
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l538
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l538
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l538
										}
										position++
										break
									}
								}

							}
						l539:
							goto l537
						l538:
							position, tokenIndex, depth = position538, tokenIndex538, depth538
						}
						if !rules[rulews]() {
							goto l527
						}
						depth--
						add(ruleprefixedName, position531)
					}
				}
			l529:
				depth--
				add(ruleiriref, position528)
			}
			return true
		l527:
			position, tokenIndex, depth = position527, tokenIndex527, depth527
			return false
		},
		/* 59 iri <- <('<' (!'>' .)* '>' ws)> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				if buffer[position] != rune('<') {
					goto l544
				}
				position++
			l546:
				{
					position547, tokenIndex547, depth547 := position, tokenIndex, depth
					{
						position548, tokenIndex548, depth548 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l548
						}
						position++
						goto l547
					l548:
						position, tokenIndex, depth = position548, tokenIndex548, depth548
					}
					if !matchDot() {
						goto l547
					}
					goto l546
				l547:
					position, tokenIndex, depth = position547, tokenIndex547, depth547
				}
				if buffer[position] != rune('>') {
					goto l544
				}
				position++
				if !rules[rulews]() {
					goto l544
				}
				depth--
				add(ruleiri, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 60 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* ws)> */
		nil,
		/* 61 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? ws)> */
		func() bool {
			position550, tokenIndex550, depth550 := position, tokenIndex, depth
			{
				position551 := position
				depth++
				{
					position552 := position
					depth++
					if buffer[position] != rune('"') {
						goto l550
					}
					position++
				l553:
					{
						position554, tokenIndex554, depth554 := position, tokenIndex, depth
						{
							position555, tokenIndex555, depth555 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l555
							}
							position++
							goto l554
						l555:
							position, tokenIndex, depth = position555, tokenIndex555, depth555
						}
						if !matchDot() {
							goto l554
						}
						goto l553
					l554:
						position, tokenIndex, depth = position554, tokenIndex554, depth554
					}
					if buffer[position] != rune('"') {
						goto l550
					}
					position++
					depth--
					add(rulestring, position552)
				}
				{
					position556, tokenIndex556, depth556 := position, tokenIndex, depth
					{
						position558, tokenIndex558, depth558 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l559
						}
						position++
						{
							position562, tokenIndex562, depth562 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l563
							}
							position++
							goto l562
						l563:
							position, tokenIndex, depth = position562, tokenIndex562, depth562
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l559
							}
							position++
						}
					l562:
					l560:
						{
							position561, tokenIndex561, depth561 := position, tokenIndex, depth
							{
								position564, tokenIndex564, depth564 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l565
								}
								position++
								goto l564
							l565:
								position, tokenIndex, depth = position564, tokenIndex564, depth564
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l561
								}
								position++
							}
						l564:
							goto l560
						l561:
							position, tokenIndex, depth = position561, tokenIndex561, depth561
						}
					l566:
						{
							position567, tokenIndex567, depth567 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l567
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l567
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l567
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l567
									}
									position++
									break
								}
							}

						l568:
							{
								position569, tokenIndex569, depth569 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l569
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l569
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l569
										}
										position++
										break
									}
								}

								goto l568
							l569:
								position, tokenIndex, depth = position569, tokenIndex569, depth569
							}
							goto l566
						l567:
							position, tokenIndex, depth = position567, tokenIndex567, depth567
						}
						goto l558
					l559:
						position, tokenIndex, depth = position558, tokenIndex558, depth558
						if buffer[position] != rune('^') {
							goto l556
						}
						position++
						if buffer[position] != rune('^') {
							goto l556
						}
						position++
						if !rules[ruleiriref]() {
							goto l556
						}
					}
				l558:
					goto l557
				l556:
					position, tokenIndex, depth = position556, tokenIndex556, depth556
				}
			l557:
				if !rules[rulews]() {
					goto l550
				}
				depth--
				add(ruleliteral, position551)
			}
			return true
		l550:
			position, tokenIndex, depth = position550, tokenIndex550, depth550
			return false
		},
		/* 62 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 63 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? ws)> */
		func() bool {
			position573, tokenIndex573, depth573 := position, tokenIndex, depth
			{
				position574 := position
				depth++
				{
					position575, tokenIndex575, depth575 := position, tokenIndex, depth
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('-') {
							goto l575
						}
						position++
					}
				l577:
					goto l576
				l575:
					position, tokenIndex, depth = position575, tokenIndex575, depth575
				}
			l576:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l573
				}
				position++
			l579:
				{
					position580, tokenIndex580, depth580 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l580
					}
					position++
					goto l579
				l580:
					position, tokenIndex, depth = position580, tokenIndex580, depth580
				}
				{
					position581, tokenIndex581, depth581 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l581
					}
					position++
				l583:
					{
						position584, tokenIndex584, depth584 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l584
						}
						position++
						goto l583
					l584:
						position, tokenIndex, depth = position584, tokenIndex584, depth584
					}
					goto l582
				l581:
					position, tokenIndex, depth = position581, tokenIndex581, depth581
				}
			l582:
				if !rules[rulews]() {
					goto l573
				}
				depth--
				add(rulenumericLiteral, position574)
			}
			return true
		l573:
			position, tokenIndex, depth = position573, tokenIndex573, depth573
			return false
		},
		/* 64 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? ws)> */
		nil,
		/* 65 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position586, tokenIndex586, depth586 := position, tokenIndex, depth
			{
				position587 := position
				depth++
				{
					position588, tokenIndex588, depth588 := position, tokenIndex, depth
					{
						position590 := position
						depth++
						{
							position591, tokenIndex591, depth591 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l592
							}
							position++
							goto l591
						l592:
							position, tokenIndex, depth = position591, tokenIndex591, depth591
							if buffer[position] != rune('T') {
								goto l589
							}
							position++
						}
					l591:
						{
							position593, tokenIndex593, depth593 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l594
							}
							position++
							goto l593
						l594:
							position, tokenIndex, depth = position593, tokenIndex593, depth593
							if buffer[position] != rune('R') {
								goto l589
							}
							position++
						}
					l593:
						{
							position595, tokenIndex595, depth595 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l596
							}
							position++
							goto l595
						l596:
							position, tokenIndex, depth = position595, tokenIndex595, depth595
							if buffer[position] != rune('U') {
								goto l589
							}
							position++
						}
					l595:
						{
							position597, tokenIndex597, depth597 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l598
							}
							position++
							goto l597
						l598:
							position, tokenIndex, depth = position597, tokenIndex597, depth597
							if buffer[position] != rune('E') {
								goto l589
							}
							position++
						}
					l597:
						if !rules[rulews]() {
							goto l589
						}
						depth--
						add(ruleTRUE, position590)
					}
					goto l588
				l589:
					position, tokenIndex, depth = position588, tokenIndex588, depth588
					{
						position599 := position
						depth++
						{
							position600, tokenIndex600, depth600 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l601
							}
							position++
							goto l600
						l601:
							position, tokenIndex, depth = position600, tokenIndex600, depth600
							if buffer[position] != rune('F') {
								goto l586
							}
							position++
						}
					l600:
						{
							position602, tokenIndex602, depth602 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l603
							}
							position++
							goto l602
						l603:
							position, tokenIndex, depth = position602, tokenIndex602, depth602
							if buffer[position] != rune('A') {
								goto l586
							}
							position++
						}
					l602:
						{
							position604, tokenIndex604, depth604 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l605
							}
							position++
							goto l604
						l605:
							position, tokenIndex, depth = position604, tokenIndex604, depth604
							if buffer[position] != rune('L') {
								goto l586
							}
							position++
						}
					l604:
						{
							position606, tokenIndex606, depth606 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l607
							}
							position++
							goto l606
						l607:
							position, tokenIndex, depth = position606, tokenIndex606, depth606
							if buffer[position] != rune('S') {
								goto l586
							}
							position++
						}
					l606:
						{
							position608, tokenIndex608, depth608 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l609
							}
							position++
							goto l608
						l609:
							position, tokenIndex, depth = position608, tokenIndex608, depth608
							if buffer[position] != rune('E') {
								goto l586
							}
							position++
						}
					l608:
						if !rules[rulews]() {
							goto l586
						}
						depth--
						add(ruleFALSE, position599)
					}
				}
			l588:
				depth--
				add(rulebooleanLiteral, position587)
			}
			return true
		l586:
			position, tokenIndex, depth = position586, tokenIndex586, depth586
			return false
		},
		/* 66 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 67 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? ws)> */
		nil,
		/* 68 anon <- <('[' ws ']' ws)> */
		nil,
		/* 69 nil <- <('(' ws ')' ws)> */
		func() bool {
			position613, tokenIndex613, depth613 := position, tokenIndex, depth
			{
				position614 := position
				depth++
				if buffer[position] != rune('(') {
					goto l613
				}
				position++
				if !rules[rulews]() {
					goto l613
				}
				if buffer[position] != rune(')') {
					goto l613
				}
				position++
				if !rules[rulews]() {
					goto l613
				}
				depth--
				add(rulenil, position614)
			}
			return true
		l613:
			position, tokenIndex, depth = position613, tokenIndex613, depth613
			return false
		},
		/* 70 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 71 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 72 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 73 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') ws)> */
		nil,
		/* 74 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') ws)> */
		nil,
		/* 75 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 76 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') ws)> */
		nil,
		/* 77 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 78 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 79 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 80 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') ws)> */
		nil,
		/* 81 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') ws)> */
		nil,
		/* 82 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') ws)> */
		nil,
		/* 83 LBRACE <- <('{' ws)> */
		func() bool {
			position628, tokenIndex628, depth628 := position, tokenIndex, depth
			{
				position629 := position
				depth++
				if buffer[position] != rune('{') {
					goto l628
				}
				position++
				if !rules[rulews]() {
					goto l628
				}
				depth--
				add(ruleLBRACE, position629)
			}
			return true
		l628:
			position, tokenIndex, depth = position628, tokenIndex628, depth628
			return false
		},
		/* 84 RBRACE <- <('}' ws)> */
		func() bool {
			position630, tokenIndex630, depth630 := position, tokenIndex, depth
			{
				position631 := position
				depth++
				if buffer[position] != rune('}') {
					goto l630
				}
				position++
				if !rules[rulews]() {
					goto l630
				}
				depth--
				add(ruleRBRACE, position631)
			}
			return true
		l630:
			position, tokenIndex, depth = position630, tokenIndex630, depth630
			return false
		},
		/* 85 LBRACK <- <('[' ws)> */
		nil,
		/* 86 RBRACK <- <(']' ws)> */
		nil,
		/* 87 SEMICOLON <- <(';' ws)> */
		nil,
		/* 88 COMMA <- <(',' ws)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				if buffer[position] != rune(',') {
					goto l635
				}
				position++
				if !rules[rulews]() {
					goto l635
				}
				depth--
				add(ruleCOMMA, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 89 DOT <- <('.' ws)> */
		func() bool {
			position637, tokenIndex637, depth637 := position, tokenIndex, depth
			{
				position638 := position
				depth++
				if buffer[position] != rune('.') {
					goto l637
				}
				position++
				if !rules[rulews]() {
					goto l637
				}
				depth--
				add(ruleDOT, position638)
			}
			return true
		l637:
			position, tokenIndex, depth = position637, tokenIndex637, depth637
			return false
		},
		/* 90 COLON <- <(':' ws)> */
		nil,
		/* 91 PIPE <- <('|' ws)> */
		func() bool {
			position640, tokenIndex640, depth640 := position, tokenIndex, depth
			{
				position641 := position
				depth++
				if buffer[position] != rune('|') {
					goto l640
				}
				position++
				if !rules[rulews]() {
					goto l640
				}
				depth--
				add(rulePIPE, position641)
			}
			return true
		l640:
			position, tokenIndex, depth = position640, tokenIndex640, depth640
			return false
		},
		/* 92 SLASH <- <('/' ws)> */
		func() bool {
			position642, tokenIndex642, depth642 := position, tokenIndex, depth
			{
				position643 := position
				depth++
				if buffer[position] != rune('/') {
					goto l642
				}
				position++
				if !rules[rulews]() {
					goto l642
				}
				depth--
				add(ruleSLASH, position643)
			}
			return true
		l642:
			position, tokenIndex, depth = position642, tokenIndex642, depth642
			return false
		},
		/* 93 INVERSE <- <('^' ws)> */
		func() bool {
			position644, tokenIndex644, depth644 := position, tokenIndex, depth
			{
				position645 := position
				depth++
				if buffer[position] != rune('^') {
					goto l644
				}
				position++
				if !rules[rulews]() {
					goto l644
				}
				depth--
				add(ruleINVERSE, position645)
			}
			return true
		l644:
			position, tokenIndex, depth = position644, tokenIndex644, depth644
			return false
		},
		/* 94 LPAREN <- <('(' ws)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				if buffer[position] != rune('(') {
					goto l646
				}
				position++
				if !rules[rulews]() {
					goto l646
				}
				depth--
				add(ruleLPAREN, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 95 RPAREN <- <(')' ws)> */
		func() bool {
			position648, tokenIndex648, depth648 := position, tokenIndex, depth
			{
				position649 := position
				depth++
				if buffer[position] != rune(')') {
					goto l648
				}
				position++
				if !rules[rulews]() {
					goto l648
				}
				depth--
				add(ruleRPAREN, position649)
			}
			return true
		l648:
			position, tokenIndex, depth = position648, tokenIndex648, depth648
			return false
		},
		/* 96 ISA <- <('a' ws)> */
		func() bool {
			position650, tokenIndex650, depth650 := position, tokenIndex, depth
			{
				position651 := position
				depth++
				if buffer[position] != rune('a') {
					goto l650
				}
				position++
				if !rules[rulews]() {
					goto l650
				}
				depth--
				add(ruleISA, position651)
			}
			return true
		l650:
			position, tokenIndex, depth = position650, tokenIndex650, depth650
			return false
		},
		/* 97 NOT <- <('!' ws)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				if buffer[position] != rune('!') {
					goto l652
				}
				position++
				if !rules[rulews]() {
					goto l652
				}
				depth--
				add(ruleNOT, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 98 STAR <- <('*' ws)> */
		func() bool {
			position654, tokenIndex654, depth654 := position, tokenIndex, depth
			{
				position655 := position
				depth++
				if buffer[position] != rune('*') {
					goto l654
				}
				position++
				if !rules[rulews]() {
					goto l654
				}
				depth--
				add(ruleSTAR, position655)
			}
			return true
		l654:
			position, tokenIndex, depth = position654, tokenIndex654, depth654
			return false
		},
		/* 99 PLUS <- <('+' ws)> */
		func() bool {
			position656, tokenIndex656, depth656 := position, tokenIndex, depth
			{
				position657 := position
				depth++
				if buffer[position] != rune('+') {
					goto l656
				}
				position++
				if !rules[rulews]() {
					goto l656
				}
				depth--
				add(rulePLUS, position657)
			}
			return true
		l656:
			position, tokenIndex, depth = position656, tokenIndex656, depth656
			return false
		},
		/* 100 MINUS <- <('-' ws)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				if buffer[position] != rune('-') {
					goto l658
				}
				position++
				if !rules[rulews]() {
					goto l658
				}
				depth--
				add(ruleMINUS, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 101 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') ws)> */
		nil,
		/* 102 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') ws)> */
		nil,
		/* 103 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') ws)> */
		nil,
		/* 104 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') ws)> */
		nil,
		/* 105 INTEGER <- <([0-9]+ ws)> */
		func() bool {
			position664, tokenIndex664, depth664 := position, tokenIndex, depth
			{
				position665 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l664
				}
				position++
			l666:
				{
					position667, tokenIndex667, depth667 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l667
					}
					position++
					goto l666
				l667:
					position, tokenIndex, depth = position667, tokenIndex667, depth667
				}
				if !rules[rulews]() {
					goto l664
				}
				depth--
				add(ruleINTEGER, position665)
			}
			return true
		l664:
			position, tokenIndex, depth = position664, tokenIndex664, depth664
			return false
		},
		/* 106 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') ws)> */
		nil,
		/* 107 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') ws)> */
		nil,
		/* 108 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') ws)> */
		nil,
		/* 109 OR <- <('|' '|' ws)> */
		nil,
		/* 110 AND <- <('&' '&' ws)> */
		nil,
		/* 111 EQ <- <('=' ws)> */
		nil,
		/* 112 NE <- <('!' '=' ws)> */
		nil,
		/* 113 GT <- <('>' ws)> */
		nil,
		/* 114 LT <- <('<' ws)> */
		nil,
		/* 115 LE <- <('<' '=' ws)> */
		nil,
		/* 116 GE <- <('>' '=' ws)> */
		nil,
		/* 117 IN <- <(('i' / 'I') ('n' / 'N') ws)> */
		nil,
		/* 118 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') ws)> */
		nil,
		/* 119 AS <- <(('a' / 'A') ('s' / 'S') ws)> */
		nil,
		/* 120 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\n') '\n') | (&('\r') '\r') | (&('\t') '\t') | (&(' ') ' '))*> */
		func() bool {
			{
				position683 := position
				depth++
			l684:
				{
					position685, tokenIndex685, depth685 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '\v':
							if buffer[position] != rune('\v') {
								goto l685
							}
							position++
							break
						case '\f':
							if buffer[position] != rune('\f') {
								goto l685
							}
							position++
							break
						case '\n':
							if buffer[position] != rune('\n') {
								goto l685
							}
							position++
							break
						case '\r':
							if buffer[position] != rune('\r') {
								goto l685
							}
							position++
							break
						case '\t':
							if buffer[position] != rune('\t') {
								goto l685
							}
							position++
							break
						default:
							if buffer[position] != rune(' ') {
								goto l685
							}
							position++
							break
						}
					}

					goto l684
				l685:
					position, tokenIndex, depth = position685, tokenIndex685, depth685
				}
				depth--
				add(rulews, position683)
			}
			return true
		},
	}
	p.rules = rules
}
