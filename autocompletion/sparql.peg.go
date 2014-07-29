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
	rulepof
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
	ruleskip
	rulews
	rulecomment
	ruleendOfLine
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
	ruleAction13

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
	"pof",
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
	"skip",
	"ws",
	"comment",
	"endOfLine",
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
	"Action13",

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
	*Skip

	Buffer string
	buffer []rune
	rules  [141]func() bool
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
			p.addPrefix(p.skip(buffer, begin, end))
		case ruleAction1:
			p.setSubject(p.skip(buffer, begin, end))
		case ruleAction2:
			p.setSubject(p.skip(buffer, begin, end))
		case ruleAction3:
			p.setSubject("?POF")
		case ruleAction4:
			p.setPredicate("?POF")
		case ruleAction5:
			p.setPredicate(p.skip(buffer, begin, end))
		case ruleAction6:
			p.setPredicate(p.skip(buffer, begin, end))
		case ruleAction7:
			p.setObject("?POF")
			p.addTriplePattern()
		case ruleAction8:
			p.setObject(p.skip(buffer, begin, end))
			p.addTriplePattern()
		case ruleAction9:
			p.setObject("?FillVar")
			p.addTriplePattern()
		case ruleAction10:
			p.setPrefix(p.skip(buffer, begin, end))
		case ruleAction11:
			p.setPathLength(p.skip(buffer, begin, end))
		case ruleAction12:
			p.setKeyword(p.skip(buffer, begin, end))
		case ruleAction13:
			p.commentBegin = begin

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
		/* 0 queryContainer <- <(skip prolog query !.)> */
		func() bool {
			position0, tokenIndex0, depth0 := position, tokenIndex, depth
			{
				position1 := position
				depth++
				if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
										if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
										if !rules[ruleskip]() {
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
										if !rules[ruleskip]() {
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
					if !rules[ruleskip]() {
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
							if !rules[ruleskip]() {
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
							if !rules[ruleskip]() {
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
						{
							position167, tokenIndex167, depth167 := position, tokenIndex, depth
							if !rules[rulevar]() {
								goto l168
							}
							goto l167
						l168:
							position, tokenIndex, depth = position167, tokenIndex167, depth167
							if !rules[ruleLPAREN]() {
								goto l111
							}
							if !rules[ruleexpression]() {
								goto l111
							}
							{
								position169 := position
								depth++
								{
									position170, tokenIndex170, depth170 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l171
									}
									position++
									goto l170
								l171:
									position, tokenIndex, depth = position170, tokenIndex170, depth170
									if buffer[position] != rune('A') {
										goto l111
									}
									position++
								}
							l170:
								{
									position172, tokenIndex172, depth172 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l173
									}
									position++
									goto l172
								l173:
									position, tokenIndex, depth = position172, tokenIndex172, depth172
									if buffer[position] != rune('S') {
										goto l111
									}
									position++
								}
							l172:
								if !rules[ruleskip]() {
									goto l111
								}
								depth--
								add(ruleAS, position169)
							}
							if !rules[rulevar]() {
								goto l111
							}
							if !rules[ruleRPAREN]() {
								goto l111
							}
						}
					l167:
						depth--
						add(ruleprojectionElem, position166)
					}
				l164:
					{
						position165, tokenIndex165, depth165 := position, tokenIndex, depth
						{
							position174 := position
							depth++
							{
								position175, tokenIndex175, depth175 := position, tokenIndex, depth
								if !rules[rulevar]() {
									goto l176
								}
								goto l175
							l176:
								position, tokenIndex, depth = position175, tokenIndex175, depth175
								if !rules[ruleLPAREN]() {
									goto l165
								}
								if !rules[ruleexpression]() {
									goto l165
								}
								{
									position177 := position
									depth++
									{
										position178, tokenIndex178, depth178 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l179
										}
										position++
										goto l178
									l179:
										position, tokenIndex, depth = position178, tokenIndex178, depth178
										if buffer[position] != rune('A') {
											goto l165
										}
										position++
									}
								l178:
									{
										position180, tokenIndex180, depth180 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l181
										}
										position++
										goto l180
									l181:
										position, tokenIndex, depth = position180, tokenIndex180, depth180
										if buffer[position] != rune('S') {
											goto l165
										}
										position++
									}
								l180:
									if !rules[ruleskip]() {
										goto l165
									}
									depth--
									add(ruleAS, position177)
								}
								if !rules[rulevar]() {
									goto l165
								}
								if !rules[ruleRPAREN]() {
									goto l165
								}
							}
						l175:
							depth--
							add(ruleprojectionElem, position174)
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
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				if !rules[ruleselect]() {
					goto l182
				}
				if !rules[rulewhereClause]() {
					goto l182
				}
				depth--
				add(rulesubSelect, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
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
			position190, tokenIndex190, depth190 := position, tokenIndex, depth
			{
				position191 := position
				depth++
				{
					position192 := position
					depth++
					{
						position193, tokenIndex193, depth193 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l194
						}
						position++
						goto l193
					l194:
						position, tokenIndex, depth = position193, tokenIndex193, depth193
						if buffer[position] != rune('F') {
							goto l190
						}
						position++
					}
				l193:
					{
						position195, tokenIndex195, depth195 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l196
						}
						position++
						goto l195
					l196:
						position, tokenIndex, depth = position195, tokenIndex195, depth195
						if buffer[position] != rune('R') {
							goto l190
						}
						position++
					}
				l195:
					{
						position197, tokenIndex197, depth197 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l198
						}
						position++
						goto l197
					l198:
						position, tokenIndex, depth = position197, tokenIndex197, depth197
						if buffer[position] != rune('O') {
							goto l190
						}
						position++
					}
				l197:
					{
						position199, tokenIndex199, depth199 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l200
						}
						position++
						goto l199
					l200:
						position, tokenIndex, depth = position199, tokenIndex199, depth199
						if buffer[position] != rune('M') {
							goto l190
						}
						position++
					}
				l199:
					if !rules[ruleskip]() {
						goto l190
					}
					depth--
					add(ruleFROM, position192)
				}
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					{
						position203 := position
						depth++
						{
							position204, tokenIndex204, depth204 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l205
							}
							position++
							goto l204
						l205:
							position, tokenIndex, depth = position204, tokenIndex204, depth204
							if buffer[position] != rune('N') {
								goto l201
							}
							position++
						}
					l204:
						{
							position206, tokenIndex206, depth206 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l207
							}
							position++
							goto l206
						l207:
							position, tokenIndex, depth = position206, tokenIndex206, depth206
							if buffer[position] != rune('A') {
								goto l201
							}
							position++
						}
					l206:
						{
							position208, tokenIndex208, depth208 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l209
							}
							position++
							goto l208
						l209:
							position, tokenIndex, depth = position208, tokenIndex208, depth208
							if buffer[position] != rune('M') {
								goto l201
							}
							position++
						}
					l208:
						{
							position210, tokenIndex210, depth210 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l211
							}
							position++
							goto l210
						l211:
							position, tokenIndex, depth = position210, tokenIndex210, depth210
							if buffer[position] != rune('E') {
								goto l201
							}
							position++
						}
					l210:
						{
							position212, tokenIndex212, depth212 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l213
							}
							position++
							goto l212
						l213:
							position, tokenIndex, depth = position212, tokenIndex212, depth212
							if buffer[position] != rune('D') {
								goto l201
							}
							position++
						}
					l212:
						if !rules[ruleskip]() {
							goto l201
						}
						depth--
						add(ruleNAMED, position203)
					}
					goto l202
				l201:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
				}
			l202:
				if !rules[ruleiriref]() {
					goto l190
				}
				depth--
				add(ruledatasetClause, position191)
			}
			return true
		l190:
			position, tokenIndex, depth = position190, tokenIndex190, depth190
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position214, tokenIndex214, depth214 := position, tokenIndex, depth
			{
				position215 := position
				depth++
				{
					position216, tokenIndex216, depth216 := position, tokenIndex, depth
					{
						position218 := position
						depth++
						{
							position219, tokenIndex219, depth219 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l220
							}
							position++
							goto l219
						l220:
							position, tokenIndex, depth = position219, tokenIndex219, depth219
							if buffer[position] != rune('W') {
								goto l216
							}
							position++
						}
					l219:
						{
							position221, tokenIndex221, depth221 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l222
							}
							position++
							goto l221
						l222:
							position, tokenIndex, depth = position221, tokenIndex221, depth221
							if buffer[position] != rune('H') {
								goto l216
							}
							position++
						}
					l221:
						{
							position223, tokenIndex223, depth223 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l224
							}
							position++
							goto l223
						l224:
							position, tokenIndex, depth = position223, tokenIndex223, depth223
							if buffer[position] != rune('E') {
								goto l216
							}
							position++
						}
					l223:
						{
							position225, tokenIndex225, depth225 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l226
							}
							position++
							goto l225
						l226:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if buffer[position] != rune('R') {
								goto l216
							}
							position++
						}
					l225:
						{
							position227, tokenIndex227, depth227 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l228
							}
							position++
							goto l227
						l228:
							position, tokenIndex, depth = position227, tokenIndex227, depth227
							if buffer[position] != rune('E') {
								goto l216
							}
							position++
						}
					l227:
						if !rules[ruleskip]() {
							goto l216
						}
						depth--
						add(ruleWHERE, position218)
					}
					goto l217
				l216:
					position, tokenIndex, depth = position216, tokenIndex216, depth216
				}
			l217:
				if !rules[rulegroupGraphPattern]() {
					goto l214
				}
				depth--
				add(rulewhereClause, position215)
			}
			return true
		l214:
			position, tokenIndex, depth = position214, tokenIndex214, depth214
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position229, tokenIndex229, depth229 := position, tokenIndex, depth
			{
				position230 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l229
				}
				{
					position231, tokenIndex231, depth231 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l232
					}
					goto l231
				l232:
					position, tokenIndex, depth = position231, tokenIndex231, depth231
					if !rules[rulegraphPattern]() {
						goto l229
					}
				}
			l231:
				if !rules[ruleRBRACE]() {
					goto l229
				}
				depth--
				add(rulegroupGraphPattern, position230)
			}
			return true
		l229:
			position, tokenIndex, depth = position229, tokenIndex229, depth229
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position234 := position
				depth++
				{
					position235, tokenIndex235, depth235 := position, tokenIndex, depth
					{
						position237 := position
						depth++
						if !rules[ruletriplesBlock]() {
							goto l235
						}
						depth--
						add(rulebasicGraphPattern, position237)
					}
					goto l236
				l235:
					position, tokenIndex, depth = position235, tokenIndex235, depth235
				}
			l236:
				{
					position238, tokenIndex238, depth238 := position, tokenIndex, depth
					{
						position240 := position
						depth++
						{
							position241, tokenIndex241, depth241 := position, tokenIndex, depth
							{
								position243 := position
								depth++
								{
									position244 := position
									depth++
									{
										position245, tokenIndex245, depth245 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l246
										}
										position++
										goto l245
									l246:
										position, tokenIndex, depth = position245, tokenIndex245, depth245
										if buffer[position] != rune('O') {
											goto l242
										}
										position++
									}
								l245:
									{
										position247, tokenIndex247, depth247 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l248
										}
										position++
										goto l247
									l248:
										position, tokenIndex, depth = position247, tokenIndex247, depth247
										if buffer[position] != rune('P') {
											goto l242
										}
										position++
									}
								l247:
									{
										position249, tokenIndex249, depth249 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l250
										}
										position++
										goto l249
									l250:
										position, tokenIndex, depth = position249, tokenIndex249, depth249
										if buffer[position] != rune('T') {
											goto l242
										}
										position++
									}
								l249:
									{
										position251, tokenIndex251, depth251 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l252
										}
										position++
										goto l251
									l252:
										position, tokenIndex, depth = position251, tokenIndex251, depth251
										if buffer[position] != rune('I') {
											goto l242
										}
										position++
									}
								l251:
									{
										position253, tokenIndex253, depth253 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l254
										}
										position++
										goto l253
									l254:
										position, tokenIndex, depth = position253, tokenIndex253, depth253
										if buffer[position] != rune('O') {
											goto l242
										}
										position++
									}
								l253:
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l256
										}
										position++
										goto l255
									l256:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
										if buffer[position] != rune('N') {
											goto l242
										}
										position++
									}
								l255:
									{
										position257, tokenIndex257, depth257 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l258
										}
										position++
										goto l257
									l258:
										position, tokenIndex, depth = position257, tokenIndex257, depth257
										if buffer[position] != rune('A') {
											goto l242
										}
										position++
									}
								l257:
									{
										position259, tokenIndex259, depth259 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l260
										}
										position++
										goto l259
									l260:
										position, tokenIndex, depth = position259, tokenIndex259, depth259
										if buffer[position] != rune('L') {
											goto l242
										}
										position++
									}
								l259:
									if !rules[ruleskip]() {
										goto l242
									}
									depth--
									add(ruleOPTIONAL, position244)
								}
								if !rules[ruleLBRACE]() {
									goto l242
								}
								{
									position261, tokenIndex261, depth261 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l262
									}
									goto l261
								l262:
									position, tokenIndex, depth = position261, tokenIndex261, depth261
									if !rules[rulegraphPattern]() {
										goto l242
									}
								}
							l261:
								if !rules[ruleRBRACE]() {
									goto l242
								}
								depth--
								add(ruleoptionalGraphPattern, position243)
							}
							goto l241
						l242:
							position, tokenIndex, depth = position241, tokenIndex241, depth241
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l238
							}
						}
					l241:
						depth--
						add(rulegraphPatternNotTriples, position240)
					}
					{
						position263, tokenIndex263, depth263 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l263
						}
						goto l264
					l263:
						position, tokenIndex, depth = position263, tokenIndex263, depth263
					}
				l264:
					if !rules[rulegraphPattern]() {
						goto l238
					}
					goto l239
				l238:
					position, tokenIndex, depth = position238, tokenIndex238, depth238
				}
			l239:
				depth--
				add(rulegraphPattern, position234)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position267, tokenIndex267, depth267 := position, tokenIndex, depth
			{
				position268 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l267
				}
				{
					position269, tokenIndex269, depth269 := position, tokenIndex, depth
					{
						position271 := position
						depth++
						{
							position272, tokenIndex272, depth272 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l273
							}
							position++
							goto l272
						l273:
							position, tokenIndex, depth = position272, tokenIndex272, depth272
							if buffer[position] != rune('U') {
								goto l269
							}
							position++
						}
					l272:
						{
							position274, tokenIndex274, depth274 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l275
							}
							position++
							goto l274
						l275:
							position, tokenIndex, depth = position274, tokenIndex274, depth274
							if buffer[position] != rune('N') {
								goto l269
							}
							position++
						}
					l274:
						{
							position276, tokenIndex276, depth276 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l277
							}
							position++
							goto l276
						l277:
							position, tokenIndex, depth = position276, tokenIndex276, depth276
							if buffer[position] != rune('I') {
								goto l269
							}
							position++
						}
					l276:
						{
							position278, tokenIndex278, depth278 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l279
							}
							position++
							goto l278
						l279:
							position, tokenIndex, depth = position278, tokenIndex278, depth278
							if buffer[position] != rune('O') {
								goto l269
							}
							position++
						}
					l278:
						{
							position280, tokenIndex280, depth280 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l281
							}
							position++
							goto l280
						l281:
							position, tokenIndex, depth = position280, tokenIndex280, depth280
							if buffer[position] != rune('N') {
								goto l269
							}
							position++
						}
					l280:
						if !rules[ruleskip]() {
							goto l269
						}
						depth--
						add(ruleUNION, position271)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l269
					}
					goto l270
				l269:
					position, tokenIndex, depth = position269, tokenIndex269, depth269
				}
			l270:
				depth--
				add(rulegroupOrUnionGraphPattern, position268)
			}
			return true
		l267:
			position, tokenIndex, depth = position267, tokenIndex267, depth267
			return false
		},
		/* 21 basicGraphPattern <- <triplesBlock> */
		nil,
		/* 22 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position283, tokenIndex283, depth283 := position, tokenIndex, depth
			{
				position284 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l283
				}
			l285:
				{
					position286, tokenIndex286, depth286 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l286
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l286
					}
					goto l285
				l286:
					position, tokenIndex, depth = position286, tokenIndex286, depth286
				}
				{
					position287, tokenIndex287, depth287 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l287
					}
					goto l288
				l287:
					position, tokenIndex, depth = position287, tokenIndex287, depth287
				}
			l288:
				depth--
				add(ruletriplesBlock, position284)
			}
			return true
		l283:
			position, tokenIndex, depth = position283, tokenIndex283, depth283
			return false
		},
		/* 23 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position289, tokenIndex289, depth289 := position, tokenIndex, depth
			{
				position290 := position
				depth++
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					{
						position293 := position
						depth++
						{
							position294, tokenIndex294, depth294 := position, tokenIndex, depth
							{
								position296 := position
								depth++
								if !rules[rulevar]() {
									goto l295
								}
								depth--
								add(rulePegText, position296)
							}
							{
								add(ruleAction1, position)
							}
							goto l294
						l295:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							{
								position299 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l298
								}
								depth--
								add(rulePegText, position299)
							}
							{
								add(ruleAction2, position)
							}
							goto l294
						l298:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							if !rules[rulepof]() {
								goto l292
							}
							{
								add(ruleAction3, position)
							}
						}
					l294:
						depth--
						add(rulevarOrTerm, position293)
					}
					if !rules[rulepropertyListPath]() {
						goto l292
					}
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					{
						position302 := position
						depth++
						{
							position303, tokenIndex303, depth303 := position, tokenIndex, depth
							{
								position305 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l304
								}
								if !rules[rulegraphNodePath]() {
									goto l304
								}
							l306:
								{
									position307, tokenIndex307, depth307 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l307
									}
									goto l306
								l307:
									position, tokenIndex, depth = position307, tokenIndex307, depth307
								}
								if !rules[ruleRPAREN]() {
									goto l304
								}
								depth--
								add(rulecollectionPath, position305)
							}
							goto l303
						l304:
							position, tokenIndex, depth = position303, tokenIndex303, depth303
							{
								position308 := position
								depth++
								{
									position309 := position
									depth++
									if buffer[position] != rune('[') {
										goto l289
									}
									position++
									if !rules[ruleskip]() {
										goto l289
									}
									depth--
									add(ruleLBRACK, position309)
								}
								if !rules[rulepropertyListPath]() {
									goto l289
								}
								{
									position310 := position
									depth++
									if buffer[position] != rune(']') {
										goto l289
									}
									position++
									if !rules[ruleskip]() {
										goto l289
									}
									depth--
									add(ruleRBRACK, position310)
								}
								depth--
								add(ruleblankNodePropertyListPath, position308)
							}
						}
					l303:
						depth--
						add(ruletriplesNodePath, position302)
					}
					if !rules[rulepropertyListPath]() {
						goto l289
					}
				}
			l291:
				depth--
				add(ruletriplesSameSubjectPath, position290)
			}
			return true
		l289:
			position, tokenIndex, depth = position289, tokenIndex289, depth289
			return false
		},
		/* 24 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 25 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position312, tokenIndex312, depth312 := position, tokenIndex, depth
			{
				position313 := position
				depth++
				{
					position314, tokenIndex314, depth314 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l315
					}
					goto l314
				l315:
					position, tokenIndex, depth = position314, tokenIndex314, depth314
					{
						switch buffer[position] {
						case '(':
							if !rules[rulenil]() {
								goto l312
							}
							break
						case '[', '_':
							{
								position317 := position
								depth++
								{
									position318, tokenIndex318, depth318 := position, tokenIndex, depth
									{
										position320 := position
										depth++
										if buffer[position] != rune('_') {
											goto l319
										}
										position++
										if buffer[position] != rune(':') {
											goto l319
										}
										position++
										{
											switch buffer[position] {
											case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l319
												}
												position++
												break
											case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l319
												}
												position++
												break
											default:
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l319
												}
												position++
												break
											}
										}

										{
											position322, tokenIndex322, depth322 := position, tokenIndex, depth
											{
												position324, tokenIndex324, depth324 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l325
												}
												position++
												goto l324
											l325:
												position, tokenIndex, depth = position324, tokenIndex324, depth324
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l326
												}
												position++
												goto l324
											l326:
												position, tokenIndex, depth = position324, tokenIndex324, depth324
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l327
												}
												position++
												goto l324
											l327:
												position, tokenIndex, depth = position324, tokenIndex324, depth324
												if c := buffer[position]; c < rune('.') || c > rune('_') {
													goto l322
												}
												position++
											}
										l324:
											goto l323
										l322:
											position, tokenIndex, depth = position322, tokenIndex322, depth322
										}
									l323:
										if !rules[ruleskip]() {
											goto l319
										}
										depth--
										add(ruleblankNodeLabel, position320)
									}
									goto l318
								l319:
									position, tokenIndex, depth = position318, tokenIndex318, depth318
									{
										position328 := position
										depth++
										if buffer[position] != rune('[') {
											goto l312
										}
										position++
									l329:
										{
											position330, tokenIndex330, depth330 := position, tokenIndex, depth
											if !rules[rulews]() {
												goto l330
											}
											goto l329
										l330:
											position, tokenIndex, depth = position330, tokenIndex330, depth330
										}
										if buffer[position] != rune(']') {
											goto l312
										}
										position++
										if !rules[ruleskip]() {
											goto l312
										}
										depth--
										add(ruleanon, position328)
									}
								}
							l318:
								depth--
								add(ruleblankNode, position317)
							}
							break
						case 'F', 'T', 'f', 't':
							if !rules[rulebooleanLiteral]() {
								goto l312
							}
							break
						case '"':
							if !rules[ruleliteral]() {
								goto l312
							}
							break
						default:
							if !rules[rulenumericLiteral]() {
								goto l312
							}
							break
						}
					}

				}
			l314:
				depth--
				add(rulegraphTerm, position313)
			}
			return true
		l312:
			position, tokenIndex, depth = position312, tokenIndex312, depth312
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
			position334, tokenIndex334, depth334 := position, tokenIndex, depth
			{
				position335 := position
				depth++
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l337
					}
					{
						add(ruleAction4, position)
					}
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					{
						position340 := position
						depth++
						if !rules[rulevar]() {
							goto l339
						}
						depth--
						add(rulePegText, position340)
					}
					{
						add(ruleAction5, position)
					}
					goto l336
				l339:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					{
						position342 := position
						depth++
						if !rules[rulepath]() {
							goto l334
						}
						depth--
						add(ruleverbPath, position342)
					}
				}
			l336:
				if !rules[ruleobjectListPath]() {
					goto l334
				}
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					{
						position345 := position
						depth++
						if buffer[position] != rune(';') {
							goto l343
						}
						position++
						if !rules[ruleskip]() {
							goto l343
						}
						depth--
						add(ruleSEMICOLON, position345)
					}
					if !rules[rulepropertyListPath]() {
						goto l343
					}
					goto l344
				l343:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
				}
			l344:
				depth--
				add(rulepropertyListPath, position335)
			}
			return true
		l334:
			position, tokenIndex, depth = position334, tokenIndex334, depth334
			return false
		},
		/* 30 verbPath <- <path> */
		nil,
		/* 31 path <- <pathAlternative> */
		func() bool {
			position347, tokenIndex347, depth347 := position, tokenIndex, depth
			{
				position348 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l347
				}
				depth--
				add(rulepath, position348)
			}
			return true
		l347:
			position, tokenIndex, depth = position347, tokenIndex347, depth347
			return false
		},
		/* 32 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l349
				}
			l351:
				{
					position352, tokenIndex352, depth352 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l352
					}
					if !rules[rulepathAlternative]() {
						goto l352
					}
					goto l351
				l352:
					position, tokenIndex, depth = position352, tokenIndex352, depth352
				}
				depth--
				add(rulepathAlternative, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 33 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position353, tokenIndex353, depth353 := position, tokenIndex, depth
			{
				position354 := position
				depth++
				{
					position355 := position
					depth++
					{
						position356 := position
						depth++
						{
							position357, tokenIndex357, depth357 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l357
							}
							goto l358
						l357:
							position, tokenIndex, depth = position357, tokenIndex357, depth357
						}
					l358:
						{
							position359 := position
							depth++
							{
								position360, tokenIndex360, depth360 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l361
								}
								goto l360
							l361:
								position, tokenIndex, depth = position360, tokenIndex360, depth360
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l353
										}
										if !rules[rulepath]() {
											goto l353
										}
										if !rules[ruleRPAREN]() {
											goto l353
										}
										break
									case '!':
										if !rules[ruleNOT]() {
											goto l353
										}
										{
											position363 := position
											depth++
											{
												position364, tokenIndex364, depth364 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l365
												}
												goto l364
											l365:
												position, tokenIndex, depth = position364, tokenIndex364, depth364
												if !rules[ruleLPAREN]() {
													goto l353
												}
												{
													position366, tokenIndex366, depth366 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l366
													}
												l368:
													{
														position369, tokenIndex369, depth369 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l369
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l369
														}
														goto l368
													l369:
														position, tokenIndex, depth = position369, tokenIndex369, depth369
													}
													goto l367
												l366:
													position, tokenIndex, depth = position366, tokenIndex366, depth366
												}
											l367:
												if !rules[ruleRPAREN]() {
													goto l353
												}
											}
										l364:
											depth--
											add(rulepathNegatedPropertySet, position363)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l353
										}
										break
									}
								}

							}
						l360:
							depth--
							add(rulepathPrimary, position359)
						}
						depth--
						add(rulepathElt, position356)
					}
					depth--
					add(rulePegText, position355)
				}
				{
					add(ruleAction6, position)
				}
			l371:
				{
					position372, tokenIndex372, depth372 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l372
					}
					if !rules[rulepathSequence]() {
						goto l372
					}
					goto l371
				l372:
					position, tokenIndex, depth = position372, tokenIndex372, depth372
				}
				depth--
				add(rulepathSequence, position354)
			}
			return true
		l353:
			position, tokenIndex, depth = position353, tokenIndex353, depth353
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
			position376, tokenIndex376, depth376 := position, tokenIndex, depth
			{
				position377 := position
				depth++
				{
					position378, tokenIndex378, depth378 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l379
					}
					goto l378
				l379:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
					if !rules[ruleISA]() {
						goto l380
					}
					goto l378
				l380:
					position, tokenIndex, depth = position378, tokenIndex378, depth378
					if !rules[ruleINVERSE]() {
						goto l376
					}
					{
						position381, tokenIndex381, depth381 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l382
						}
						goto l381
					l382:
						position, tokenIndex, depth = position381, tokenIndex381, depth381
						if !rules[ruleISA]() {
							goto l376
						}
					}
				l381:
				}
			l378:
				depth--
				add(rulepathOneInPropertySet, position377)
			}
			return true
		l376:
			position, tokenIndex, depth = position376, tokenIndex376, depth376
			return false
		},
		/* 38 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			{
				position384 := position
				depth++
				{
					position385 := position
					depth++
					{
						position386, tokenIndex386, depth386 := position, tokenIndex, depth
						if !rules[rulepof]() {
							goto l387
						}
						{
							add(ruleAction7, position)
						}
						goto l386
					l387:
						position, tokenIndex, depth = position386, tokenIndex386, depth386
						{
							position390 := position
							depth++
							if !rules[rulegraphNodePath]() {
								goto l389
							}
							depth--
							add(rulePegText, position390)
						}
						{
							add(ruleAction8, position)
						}
						goto l386
					l389:
						position, tokenIndex, depth = position386, tokenIndex386, depth386
						{
							add(ruleAction9, position)
						}
					}
				l386:
					depth--
					add(ruleobjectPath, position385)
				}
			l393:
				{
					position394, tokenIndex394, depth394 := position, tokenIndex, depth
					if !rules[ruleCOMMA]() {
						goto l394
					}
					if !rules[ruleobjectListPath]() {
						goto l394
					}
					goto l393
				l394:
					position, tokenIndex, depth = position394, tokenIndex394, depth394
				}
				depth--
				add(ruleobjectListPath, position384)
			}
			return true
		},
		/* 39 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		nil,
		/* 40 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			position396, tokenIndex396, depth396 := position, tokenIndex, depth
			{
				position397 := position
				depth++
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l399
					}
					goto l398
				l399:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
					if !rules[rulegraphTerm]() {
						goto l396
					}
				}
			l398:
				depth--
				add(rulegraphNodePath, position397)
			}
			return true
		l396:
			position, tokenIndex, depth = position396, tokenIndex396, depth396
			return false
		},
		/* 41 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position401 := position
				depth++
				{
					position402, tokenIndex402, depth402 := position, tokenIndex, depth
					{
						position404 := position
						depth++
						{
							position405, tokenIndex405, depth405 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l406
							}
							{
								position407, tokenIndex407, depth407 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l407
								}
								goto l408
							l407:
								position, tokenIndex, depth = position407, tokenIndex407, depth407
							}
						l408:
							goto l405
						l406:
							position, tokenIndex, depth = position405, tokenIndex405, depth405
							if !rules[ruleoffset]() {
								goto l402
							}
							{
								position409, tokenIndex409, depth409 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l409
								}
								goto l410
							l409:
								position, tokenIndex, depth = position409, tokenIndex409, depth409
							}
						l410:
						}
					l405:
						depth--
						add(rulelimitOffsetClauses, position404)
					}
					goto l403
				l402:
					position, tokenIndex, depth = position402, tokenIndex402, depth402
				}
			l403:
				depth--
				add(rulesolutionModifier, position401)
			}
			return true
		},
		/* 42 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 43 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position412, tokenIndex412, depth412 := position, tokenIndex, depth
			{
				position413 := position
				depth++
				{
					position414 := position
					depth++
					{
						position415, tokenIndex415, depth415 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l416
						}
						position++
						goto l415
					l416:
						position, tokenIndex, depth = position415, tokenIndex415, depth415
						if buffer[position] != rune('L') {
							goto l412
						}
						position++
					}
				l415:
					{
						position417, tokenIndex417, depth417 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l418
						}
						position++
						goto l417
					l418:
						position, tokenIndex, depth = position417, tokenIndex417, depth417
						if buffer[position] != rune('I') {
							goto l412
						}
						position++
					}
				l417:
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l420
						}
						position++
						goto l419
					l420:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
						if buffer[position] != rune('M') {
							goto l412
						}
						position++
					}
				l419:
					{
						position421, tokenIndex421, depth421 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l422
						}
						position++
						goto l421
					l422:
						position, tokenIndex, depth = position421, tokenIndex421, depth421
						if buffer[position] != rune('I') {
							goto l412
						}
						position++
					}
				l421:
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l424
						}
						position++
						goto l423
					l424:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
						if buffer[position] != rune('T') {
							goto l412
						}
						position++
					}
				l423:
					if !rules[ruleskip]() {
						goto l412
					}
					depth--
					add(ruleLIMIT, position414)
				}
				if !rules[ruleINTEGER]() {
					goto l412
				}
				depth--
				add(rulelimit, position413)
			}
			return true
		l412:
			position, tokenIndex, depth = position412, tokenIndex412, depth412
			return false
		},
		/* 44 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position425, tokenIndex425, depth425 := position, tokenIndex, depth
			{
				position426 := position
				depth++
				{
					position427 := position
					depth++
					{
						position428, tokenIndex428, depth428 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l429
						}
						position++
						goto l428
					l429:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
						if buffer[position] != rune('O') {
							goto l425
						}
						position++
					}
				l428:
					{
						position430, tokenIndex430, depth430 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l431
						}
						position++
						goto l430
					l431:
						position, tokenIndex, depth = position430, tokenIndex430, depth430
						if buffer[position] != rune('F') {
							goto l425
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
							goto l425
						}
						position++
					}
				l432:
					{
						position434, tokenIndex434, depth434 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l435
						}
						position++
						goto l434
					l435:
						position, tokenIndex, depth = position434, tokenIndex434, depth434
						if buffer[position] != rune('S') {
							goto l425
						}
						position++
					}
				l434:
					{
						position436, tokenIndex436, depth436 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l437
						}
						position++
						goto l436
					l437:
						position, tokenIndex, depth = position436, tokenIndex436, depth436
						if buffer[position] != rune('E') {
							goto l425
						}
						position++
					}
				l436:
					{
						position438, tokenIndex438, depth438 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l439
						}
						position++
						goto l438
					l439:
						position, tokenIndex, depth = position438, tokenIndex438, depth438
						if buffer[position] != rune('T') {
							goto l425
						}
						position++
					}
				l438:
					if !rules[ruleskip]() {
						goto l425
					}
					depth--
					add(ruleOFFSET, position427)
				}
				if !rules[ruleINTEGER]() {
					goto l425
				}
				depth--
				add(ruleoffset, position426)
			}
			return true
		l425:
			position, tokenIndex, depth = position425, tokenIndex425, depth425
			return false
		},
		/* 45 expression <- <conditionalOrExpression> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l440
				}
				depth--
				add(ruleexpression, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 46 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position442, tokenIndex442, depth442 := position, tokenIndex, depth
			{
				position443 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l442
				}
				{
					position444, tokenIndex444, depth444 := position, tokenIndex, depth
					{
						position446 := position
						depth++
						if buffer[position] != rune('|') {
							goto l444
						}
						position++
						if buffer[position] != rune('|') {
							goto l444
						}
						position++
						if !rules[ruleskip]() {
							goto l444
						}
						depth--
						add(ruleOR, position446)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l444
					}
					goto l445
				l444:
					position, tokenIndex, depth = position444, tokenIndex444, depth444
				}
			l445:
				depth--
				add(ruleconditionalOrExpression, position443)
			}
			return true
		l442:
			position, tokenIndex, depth = position442, tokenIndex442, depth442
			return false
		},
		/* 47 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position447, tokenIndex447, depth447 := position, tokenIndex, depth
			{
				position448 := position
				depth++
				{
					position449 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l447
					}
					{
						position450, tokenIndex450, depth450 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position453 := position
									depth++
									{
										position454 := position
										depth++
										{
											position455, tokenIndex455, depth455 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l456
											}
											position++
											goto l455
										l456:
											position, tokenIndex, depth = position455, tokenIndex455, depth455
											if buffer[position] != rune('N') {
												goto l450
											}
											position++
										}
									l455:
										{
											position457, tokenIndex457, depth457 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l458
											}
											position++
											goto l457
										l458:
											position, tokenIndex, depth = position457, tokenIndex457, depth457
											if buffer[position] != rune('O') {
												goto l450
											}
											position++
										}
									l457:
										{
											position459, tokenIndex459, depth459 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l460
											}
											position++
											goto l459
										l460:
											position, tokenIndex, depth = position459, tokenIndex459, depth459
											if buffer[position] != rune('T') {
												goto l450
											}
											position++
										}
									l459:
										if buffer[position] != rune(' ') {
											goto l450
										}
										position++
										{
											position461, tokenIndex461, depth461 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l462
											}
											position++
											goto l461
										l462:
											position, tokenIndex, depth = position461, tokenIndex461, depth461
											if buffer[position] != rune('I') {
												goto l450
											}
											position++
										}
									l461:
										{
											position463, tokenIndex463, depth463 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l464
											}
											position++
											goto l463
										l464:
											position, tokenIndex, depth = position463, tokenIndex463, depth463
											if buffer[position] != rune('N') {
												goto l450
											}
											position++
										}
									l463:
										if !rules[ruleskip]() {
											goto l450
										}
										depth--
										add(ruleNOTIN, position454)
									}
									if !rules[ruleargList]() {
										goto l450
									}
									depth--
									add(rulenotin, position453)
								}
								break
							case 'I', 'i':
								{
									position465 := position
									depth++
									{
										position466 := position
										depth++
										{
											position467, tokenIndex467, depth467 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l468
											}
											position++
											goto l467
										l468:
											position, tokenIndex, depth = position467, tokenIndex467, depth467
											if buffer[position] != rune('I') {
												goto l450
											}
											position++
										}
									l467:
										{
											position469, tokenIndex469, depth469 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l470
											}
											position++
											goto l469
										l470:
											position, tokenIndex, depth = position469, tokenIndex469, depth469
											if buffer[position] != rune('N') {
												goto l450
											}
											position++
										}
									l469:
										if !rules[ruleskip]() {
											goto l450
										}
										depth--
										add(ruleIN, position466)
									}
									if !rules[ruleargList]() {
										goto l450
									}
									depth--
									add(rulein, position465)
								}
								break
							default:
								{
									position471, tokenIndex471, depth471 := position, tokenIndex, depth
									{
										position473 := position
										depth++
										if buffer[position] != rune('<') {
											goto l472
										}
										position++
										if !rules[ruleskip]() {
											goto l472
										}
										depth--
										add(ruleLT, position473)
									}
									goto l471
								l472:
									position, tokenIndex, depth = position471, tokenIndex471, depth471
									{
										position475 := position
										depth++
										if buffer[position] != rune('>') {
											goto l474
										}
										position++
										if buffer[position] != rune('=') {
											goto l474
										}
										position++
										if !rules[ruleskip]() {
											goto l474
										}
										depth--
										add(ruleGE, position475)
									}
									goto l471
								l474:
									position, tokenIndex, depth = position471, tokenIndex471, depth471
									{
										switch buffer[position] {
										case '>':
											{
												position477 := position
												depth++
												if buffer[position] != rune('>') {
													goto l450
												}
												position++
												if !rules[ruleskip]() {
													goto l450
												}
												depth--
												add(ruleGT, position477)
											}
											break
										case '<':
											{
												position478 := position
												depth++
												if buffer[position] != rune('<') {
													goto l450
												}
												position++
												if buffer[position] != rune('=') {
													goto l450
												}
												position++
												if !rules[ruleskip]() {
													goto l450
												}
												depth--
												add(ruleLE, position478)
											}
											break
										case '!':
											{
												position479 := position
												depth++
												if buffer[position] != rune('!') {
													goto l450
												}
												position++
												if buffer[position] != rune('=') {
													goto l450
												}
												position++
												if !rules[ruleskip]() {
													goto l450
												}
												depth--
												add(ruleNE, position479)
											}
											break
										default:
											{
												position480 := position
												depth++
												if buffer[position] != rune('=') {
													goto l450
												}
												position++
												if !rules[ruleskip]() {
													goto l450
												}
												depth--
												add(ruleEQ, position480)
											}
											break
										}
									}

								}
							l471:
								if !rules[rulenumericExpression]() {
									goto l450
								}
								break
							}
						}

						goto l451
					l450:
						position, tokenIndex, depth = position450, tokenIndex450, depth450
					}
				l451:
					depth--
					add(rulevalueLogical, position449)
				}
				{
					position481, tokenIndex481, depth481 := position, tokenIndex, depth
					{
						position483 := position
						depth++
						if buffer[position] != rune('&') {
							goto l481
						}
						position++
						if buffer[position] != rune('&') {
							goto l481
						}
						position++
						if !rules[ruleskip]() {
							goto l481
						}
						depth--
						add(ruleAND, position483)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l481
					}
					goto l482
				l481:
					position, tokenIndex, depth = position481, tokenIndex481, depth481
				}
			l482:
				depth--
				add(ruleconditionalAndExpression, position448)
			}
			return true
		l447:
			position, tokenIndex, depth = position447, tokenIndex447, depth447
			return false
		},
		/* 48 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 49 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l485
				}
			l487:
				{
					position488, tokenIndex488, depth488 := position, tokenIndex, depth
					{
						position489, tokenIndex489, depth489 := position, tokenIndex, depth
						{
							position491, tokenIndex491, depth491 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l492
							}
							goto l491
						l492:
							position, tokenIndex, depth = position491, tokenIndex491, depth491
							if !rules[ruleMINUS]() {
								goto l490
							}
						}
					l491:
						if !rules[rulemultiplicativeExpression]() {
							goto l490
						}
						goto l489
					l490:
						position, tokenIndex, depth = position489, tokenIndex489, depth489
						{
							position493 := position
							depth++
							{
								position494, tokenIndex494, depth494 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l495
								}
								position++
								goto l494
							l495:
								position, tokenIndex, depth = position494, tokenIndex494, depth494
								if buffer[position] != rune('-') {
									goto l488
								}
								position++
							}
						l494:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l488
							}
							position++
						l496:
							{
								position497, tokenIndex497, depth497 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l497
								}
								position++
								goto l496
							l497:
								position, tokenIndex, depth = position497, tokenIndex497, depth497
							}
							{
								position498, tokenIndex498, depth498 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l498
								}
								position++
							l500:
								{
									position501, tokenIndex501, depth501 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l501
									}
									position++
									goto l500
								l501:
									position, tokenIndex, depth = position501, tokenIndex501, depth501
								}
								goto l499
							l498:
								position, tokenIndex, depth = position498, tokenIndex498, depth498
							}
						l499:
							if !rules[ruleskip]() {
								goto l488
							}
							depth--
							add(rulesignedNumericLiteral, position493)
						}
					}
				l489:
					goto l487
				l488:
					position, tokenIndex, depth = position488, tokenIndex488, depth488
				}
				depth--
				add(rulenumericExpression, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 50 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position502, tokenIndex502, depth502 := position, tokenIndex, depth
			{
				position503 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l502
				}
			l504:
				{
					position505, tokenIndex505, depth505 := position, tokenIndex, depth
					{
						position506, tokenIndex506, depth506 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l507
						}
						goto l506
					l507:
						position, tokenIndex, depth = position506, tokenIndex506, depth506
						if !rules[ruleSLASH]() {
							goto l505
						}
					}
				l506:
					if !rules[ruleunaryExpression]() {
						goto l505
					}
					goto l504
				l505:
					position, tokenIndex, depth = position505, tokenIndex505, depth505
				}
				depth--
				add(rulemultiplicativeExpression, position503)
			}
			return true
		l502:
			position, tokenIndex, depth = position502, tokenIndex502, depth502
			return false
		},
		/* 51 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position508, tokenIndex508, depth508 := position, tokenIndex, depth
			{
				position509 := position
				depth++
				{
					position510, tokenIndex510, depth510 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l510
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l510
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l510
							}
							break
						}
					}

					goto l511
				l510:
					position, tokenIndex, depth = position510, tokenIndex510, depth510
				}
			l511:
				{
					position513 := position
					depth++
					{
						position514, tokenIndex514, depth514 := position, tokenIndex, depth
						{
							position516 := position
							depth++
							if !rules[ruleLPAREN]() {
								goto l515
							}
							if !rules[ruleexpression]() {
								goto l515
							}
							if !rules[ruleRPAREN]() {
								goto l515
							}
							depth--
							add(rulebrackettedExpression, position516)
						}
						goto l514
					l515:
						position, tokenIndex, depth = position514, tokenIndex514, depth514
						if !rules[ruleiriref]() {
							goto l517
						}
						goto l514
					l517:
						position, tokenIndex, depth = position514, tokenIndex514, depth514
						{
							switch buffer[position] {
							case '$', '?':
								if !rules[rulevar]() {
									goto l508
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l508
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l508
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l508
								}
								break
							}
						}

					}
				l514:
					depth--
					add(ruleprimaryExpression, position513)
				}
				depth--
				add(ruleunaryExpression, position509)
			}
			return true
		l508:
			position, tokenIndex, depth = position508, tokenIndex508, depth508
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
			position523, tokenIndex523, depth523 := position, tokenIndex, depth
			{
				position524 := position
				depth++
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l526
					}
					goto l525
				l526:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
					if !rules[ruleLPAREN]() {
						goto l523
					}
					if !rules[ruleexpression]() {
						goto l523
					}
				l527:
					{
						position528, tokenIndex528, depth528 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l528
						}
						if !rules[ruleexpression]() {
							goto l528
						}
						goto l527
					l528:
						position, tokenIndex, depth = position528, tokenIndex528, depth528
					}
					if !rules[ruleRPAREN]() {
						goto l523
					}
				}
			l525:
				depth--
				add(ruleargList, position524)
			}
			return true
		l523:
			position, tokenIndex, depth = position523, tokenIndex523, depth523
			return false
		},
		/* 57 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531, tokenIndex531, depth531 := position, tokenIndex, depth
					{
						position533 := position
						depth++
					l534:
						{
							position535, tokenIndex535, depth535 := position, tokenIndex, depth
							{
								position536, tokenIndex536, depth536 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l537
								}
								position++
								goto l536
							l537:
								position, tokenIndex, depth = position536, tokenIndex536, depth536
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l535
								}
								position++
							}
						l536:
							goto l534
						l535:
							position, tokenIndex, depth = position535, tokenIndex535, depth535
						}
						depth--
						add(rulePegText, position533)
					}
					if buffer[position] != rune(':') {
						goto l532
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l531
				l532:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					{
						position540 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l539
						}
						position++
					l541:
						{
							position542, tokenIndex542, depth542 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l542
							}
							position++
							goto l541
						l542:
							position, tokenIndex, depth = position542, tokenIndex542, depth542
						}
						depth--
						add(rulePegText, position540)
					}
					if buffer[position] != rune('/') {
						goto l539
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l531
				l539:
					position, tokenIndex, depth = position531, tokenIndex531, depth531
					{
						position544 := position
						depth++
					l545:
						{
							position546, tokenIndex546, depth546 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l546
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l546
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l546
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l546
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l546
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l546
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l546
									}
									position++
									break
								}
							}

							goto l545
						l546:
							position, tokenIndex, depth = position546, tokenIndex546, depth546
						}
						depth--
						add(rulePegText, position544)
					}
					{
						add(ruleAction12, position)
					}
				}
			l531:
				if buffer[position] != rune('<') {
					goto l529
				}
				position++
				if !rules[rulews]() {
					goto l529
				}
				if !rules[ruleskip]() {
					goto l529
				}
				depth--
				add(rulepof, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 58 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position549, tokenIndex549, depth549 := position, tokenIndex, depth
			{
				position550 := position
				depth++
				{
					position551, tokenIndex551, depth551 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l552
					}
					position++
					goto l551
				l552:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					if buffer[position] != rune('$') {
						goto l549
					}
					position++
				}
			l551:
				{
					position553 := position
					depth++
					{
						position556, tokenIndex556, depth556 := position, tokenIndex, depth
						{
							position558 := position
							depth++
							{
								position559, tokenIndex559, depth559 := position, tokenIndex, depth
								{
									position561 := position
									depth++
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
											goto l560
										}
										position++
									}
								l562:
									depth--
									add(rulePN_CHARS_BASE, position561)
								}
								goto l559
							l560:
								position, tokenIndex, depth = position559, tokenIndex559, depth559
								if buffer[position] != rune('_') {
									goto l557
								}
								position++
							}
						l559:
							depth--
							add(rulePN_CHARS_U, position558)
						}
						goto l556
					l557:
						position, tokenIndex, depth = position556, tokenIndex556, depth556
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l549
						}
						position++
					}
				l556:
				l554:
					{
						position555, tokenIndex555, depth555 := position, tokenIndex, depth
						{
							position564, tokenIndex564, depth564 := position, tokenIndex, depth
							{
								position566 := position
								depth++
								{
									position567, tokenIndex567, depth567 := position, tokenIndex, depth
									{
										position569 := position
										depth++
										{
											position570, tokenIndex570, depth570 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l571
											}
											position++
											goto l570
										l571:
											position, tokenIndex, depth = position570, tokenIndex570, depth570
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l568
											}
											position++
										}
									l570:
										depth--
										add(rulePN_CHARS_BASE, position569)
									}
									goto l567
								l568:
									position, tokenIndex, depth = position567, tokenIndex567, depth567
									if buffer[position] != rune('_') {
										goto l565
									}
									position++
								}
							l567:
								depth--
								add(rulePN_CHARS_U, position566)
							}
							goto l564
						l565:
							position, tokenIndex, depth = position564, tokenIndex564, depth564
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l555
							}
							position++
						}
					l564:
						goto l554
					l555:
						position, tokenIndex, depth = position555, tokenIndex555, depth555
					}
					depth--
					add(ruleVARNAME, position553)
				}
				if !rules[ruleskip]() {
					goto l549
				}
				depth--
				add(rulevar, position550)
			}
			return true
		l549:
			position, tokenIndex, depth = position549, tokenIndex549, depth549
			return false
		},
		/* 59 iriref <- <(iri / prefixedName)> */
		func() bool {
			position572, tokenIndex572, depth572 := position, tokenIndex, depth
			{
				position573 := position
				depth++
				{
					position574, tokenIndex574, depth574 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l575
					}
					goto l574
				l575:
					position, tokenIndex, depth = position574, tokenIndex574, depth574
					{
						position576 := position
						depth++
					l577:
						{
							position578, tokenIndex578, depth578 := position, tokenIndex, depth
							{
								position579, tokenIndex579, depth579 := position, tokenIndex, depth
								{
									position580, tokenIndex580, depth580 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l581
									}
									position++
									goto l580
								l581:
									position, tokenIndex, depth = position580, tokenIndex580, depth580
									if buffer[position] != rune(' ') {
										goto l579
									}
									position++
								}
							l580:
								goto l578
							l579:
								position, tokenIndex, depth = position579, tokenIndex579, depth579
							}
							if !matchDot() {
								goto l578
							}
							goto l577
						l578:
							position, tokenIndex, depth = position578, tokenIndex578, depth578
						}
						if buffer[position] != rune(':') {
							goto l572
						}
						position++
					l582:
						{
							position583, tokenIndex583, depth583 := position, tokenIndex, depth
							{
								position584, tokenIndex584, depth584 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l585
								}
								position++
								goto l584
							l585:
								position, tokenIndex, depth = position584, tokenIndex584, depth584
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l586
								}
								position++
								goto l584
							l586:
								position, tokenIndex, depth = position584, tokenIndex584, depth584
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l587
								}
								position++
								goto l584
							l587:
								position, tokenIndex, depth = position584, tokenIndex584, depth584
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l583
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l583
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l583
										}
										position++
										break
									}
								}

							}
						l584:
							goto l582
						l583:
							position, tokenIndex, depth = position583, tokenIndex583, depth583
						}
						if !rules[ruleskip]() {
							goto l572
						}
						depth--
						add(ruleprefixedName, position576)
					}
				}
			l574:
				depth--
				add(ruleiriref, position573)
			}
			return true
		l572:
			position, tokenIndex, depth = position572, tokenIndex572, depth572
			return false
		},
		/* 60 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position589, tokenIndex589, depth589 := position, tokenIndex, depth
			{
				position590 := position
				depth++
				if buffer[position] != rune('<') {
					goto l589
				}
				position++
			l591:
				{
					position592, tokenIndex592, depth592 := position, tokenIndex, depth
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l593
						}
						position++
						goto l592
					l593:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
					}
					if !matchDot() {
						goto l592
					}
					goto l591
				l592:
					position, tokenIndex, depth = position592, tokenIndex592, depth592
				}
				if buffer[position] != rune('>') {
					goto l589
				}
				position++
				if !rules[ruleskip]() {
					goto l589
				}
				depth--
				add(ruleiri, position590)
			}
			return true
		l589:
			position, tokenIndex, depth = position589, tokenIndex589, depth589
			return false
		},
		/* 61 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 62 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position595, tokenIndex595, depth595 := position, tokenIndex, depth
			{
				position596 := position
				depth++
				{
					position597 := position
					depth++
					if buffer[position] != rune('"') {
						goto l595
					}
					position++
				l598:
					{
						position599, tokenIndex599, depth599 := position, tokenIndex, depth
						{
							position600, tokenIndex600, depth600 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l600
							}
							position++
							goto l599
						l600:
							position, tokenIndex, depth = position600, tokenIndex600, depth600
						}
						if !matchDot() {
							goto l599
						}
						goto l598
					l599:
						position, tokenIndex, depth = position599, tokenIndex599, depth599
					}
					if buffer[position] != rune('"') {
						goto l595
					}
					position++
					depth--
					add(rulestring, position597)
				}
				{
					position601, tokenIndex601, depth601 := position, tokenIndex, depth
					{
						position603, tokenIndex603, depth603 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l604
						}
						position++
						{
							position607, tokenIndex607, depth607 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l608
							}
							position++
							goto l607
						l608:
							position, tokenIndex, depth = position607, tokenIndex607, depth607
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l604
							}
							position++
						}
					l607:
					l605:
						{
							position606, tokenIndex606, depth606 := position, tokenIndex, depth
							{
								position609, tokenIndex609, depth609 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l610
								}
								position++
								goto l609
							l610:
								position, tokenIndex, depth = position609, tokenIndex609, depth609
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l606
								}
								position++
							}
						l609:
							goto l605
						l606:
							position, tokenIndex, depth = position606, tokenIndex606, depth606
						}
					l611:
						{
							position612, tokenIndex612, depth612 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l612
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l612
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l612
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l612
									}
									position++
									break
								}
							}

						l613:
							{
								position614, tokenIndex614, depth614 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l614
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l614
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l614
										}
										position++
										break
									}
								}

								goto l613
							l614:
								position, tokenIndex, depth = position614, tokenIndex614, depth614
							}
							goto l611
						l612:
							position, tokenIndex, depth = position612, tokenIndex612, depth612
						}
						goto l603
					l604:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
						if buffer[position] != rune('^') {
							goto l601
						}
						position++
						if buffer[position] != rune('^') {
							goto l601
						}
						position++
						if !rules[ruleiriref]() {
							goto l601
						}
					}
				l603:
					goto l602
				l601:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
				}
			l602:
				if !rules[ruleskip]() {
					goto l595
				}
				depth--
				add(ruleliteral, position596)
			}
			return true
		l595:
			position, tokenIndex, depth = position595, tokenIndex595, depth595
			return false
		},
		/* 63 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 64 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position618, tokenIndex618, depth618 := position, tokenIndex, depth
			{
				position619 := position
				depth++
				{
					position620, tokenIndex620, depth620 := position, tokenIndex, depth
					{
						position622, tokenIndex622, depth622 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l623
						}
						position++
						goto l622
					l623:
						position, tokenIndex, depth = position622, tokenIndex622, depth622
						if buffer[position] != rune('-') {
							goto l620
						}
						position++
					}
				l622:
					goto l621
				l620:
					position, tokenIndex, depth = position620, tokenIndex620, depth620
				}
			l621:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l618
				}
				position++
			l624:
				{
					position625, tokenIndex625, depth625 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l625
					}
					position++
					goto l624
				l625:
					position, tokenIndex, depth = position625, tokenIndex625, depth625
				}
				{
					position626, tokenIndex626, depth626 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l626
					}
					position++
				l628:
					{
						position629, tokenIndex629, depth629 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l629
						}
						position++
						goto l628
					l629:
						position, tokenIndex, depth = position629, tokenIndex629, depth629
					}
					goto l627
				l626:
					position, tokenIndex, depth = position626, tokenIndex626, depth626
				}
			l627:
				if !rules[ruleskip]() {
					goto l618
				}
				depth--
				add(rulenumericLiteral, position619)
			}
			return true
		l618:
			position, tokenIndex, depth = position618, tokenIndex618, depth618
			return false
		},
		/* 65 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 66 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position631, tokenIndex631, depth631 := position, tokenIndex, depth
			{
				position632 := position
				depth++
				{
					position633, tokenIndex633, depth633 := position, tokenIndex, depth
					{
						position635 := position
						depth++
						{
							position636, tokenIndex636, depth636 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l637
							}
							position++
							goto l636
						l637:
							position, tokenIndex, depth = position636, tokenIndex636, depth636
							if buffer[position] != rune('T') {
								goto l634
							}
							position++
						}
					l636:
						{
							position638, tokenIndex638, depth638 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l639
							}
							position++
							goto l638
						l639:
							position, tokenIndex, depth = position638, tokenIndex638, depth638
							if buffer[position] != rune('R') {
								goto l634
							}
							position++
						}
					l638:
						{
							position640, tokenIndex640, depth640 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l641
							}
							position++
							goto l640
						l641:
							position, tokenIndex, depth = position640, tokenIndex640, depth640
							if buffer[position] != rune('U') {
								goto l634
							}
							position++
						}
					l640:
						{
							position642, tokenIndex642, depth642 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l643
							}
							position++
							goto l642
						l643:
							position, tokenIndex, depth = position642, tokenIndex642, depth642
							if buffer[position] != rune('E') {
								goto l634
							}
							position++
						}
					l642:
						if !rules[ruleskip]() {
							goto l634
						}
						depth--
						add(ruleTRUE, position635)
					}
					goto l633
				l634:
					position, tokenIndex, depth = position633, tokenIndex633, depth633
					{
						position644 := position
						depth++
						{
							position645, tokenIndex645, depth645 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l646
							}
							position++
							goto l645
						l646:
							position, tokenIndex, depth = position645, tokenIndex645, depth645
							if buffer[position] != rune('F') {
								goto l631
							}
							position++
						}
					l645:
						{
							position647, tokenIndex647, depth647 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l648
							}
							position++
							goto l647
						l648:
							position, tokenIndex, depth = position647, tokenIndex647, depth647
							if buffer[position] != rune('A') {
								goto l631
							}
							position++
						}
					l647:
						{
							position649, tokenIndex649, depth649 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l650
							}
							position++
							goto l649
						l650:
							position, tokenIndex, depth = position649, tokenIndex649, depth649
							if buffer[position] != rune('L') {
								goto l631
							}
							position++
						}
					l649:
						{
							position651, tokenIndex651, depth651 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l652
							}
							position++
							goto l651
						l652:
							position, tokenIndex, depth = position651, tokenIndex651, depth651
							if buffer[position] != rune('S') {
								goto l631
							}
							position++
						}
					l651:
						{
							position653, tokenIndex653, depth653 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l654
							}
							position++
							goto l653
						l654:
							position, tokenIndex, depth = position653, tokenIndex653, depth653
							if buffer[position] != rune('E') {
								goto l631
							}
							position++
						}
					l653:
						if !rules[ruleskip]() {
							goto l631
						}
						depth--
						add(ruleFALSE, position644)
					}
				}
			l633:
				depth--
				add(rulebooleanLiteral, position632)
			}
			return true
		l631:
			position, tokenIndex, depth = position631, tokenIndex631, depth631
			return false
		},
		/* 67 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 68 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 69 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 70 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				if buffer[position] != rune('(') {
					goto l658
				}
				position++
			l660:
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l661
					}
					goto l660
				l661:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
				}
				if buffer[position] != rune(')') {
					goto l658
				}
				position++
				if !rules[ruleskip]() {
					goto l658
				}
				depth--
				add(rulenil, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 71 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 72 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 73 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 74 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 75 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 76 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 77 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 78 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 79 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 80 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 81 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 82 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 83 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 84 LBRACE <- <('{' skip)> */
		func() bool {
			position675, tokenIndex675, depth675 := position, tokenIndex, depth
			{
				position676 := position
				depth++
				if buffer[position] != rune('{') {
					goto l675
				}
				position++
				if !rules[ruleskip]() {
					goto l675
				}
				depth--
				add(ruleLBRACE, position676)
			}
			return true
		l675:
			position, tokenIndex, depth = position675, tokenIndex675, depth675
			return false
		},
		/* 85 RBRACE <- <('}' skip)> */
		func() bool {
			position677, tokenIndex677, depth677 := position, tokenIndex, depth
			{
				position678 := position
				depth++
				if buffer[position] != rune('}') {
					goto l677
				}
				position++
				if !rules[ruleskip]() {
					goto l677
				}
				depth--
				add(ruleRBRACE, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 86 LBRACK <- <('[' skip)> */
		nil,
		/* 87 RBRACK <- <(']' skip)> */
		nil,
		/* 88 SEMICOLON <- <(';' skip)> */
		nil,
		/* 89 COMMA <- <(',' skip)> */
		func() bool {
			position682, tokenIndex682, depth682 := position, tokenIndex, depth
			{
				position683 := position
				depth++
				if buffer[position] != rune(',') {
					goto l682
				}
				position++
				if !rules[ruleskip]() {
					goto l682
				}
				depth--
				add(ruleCOMMA, position683)
			}
			return true
		l682:
			position, tokenIndex, depth = position682, tokenIndex682, depth682
			return false
		},
		/* 90 DOT <- <('.' skip)> */
		func() bool {
			position684, tokenIndex684, depth684 := position, tokenIndex, depth
			{
				position685 := position
				depth++
				if buffer[position] != rune('.') {
					goto l684
				}
				position++
				if !rules[ruleskip]() {
					goto l684
				}
				depth--
				add(ruleDOT, position685)
			}
			return true
		l684:
			position, tokenIndex, depth = position684, tokenIndex684, depth684
			return false
		},
		/* 91 COLON <- <(':' skip)> */
		nil,
		/* 92 PIPE <- <('|' skip)> */
		func() bool {
			position687, tokenIndex687, depth687 := position, tokenIndex, depth
			{
				position688 := position
				depth++
				if buffer[position] != rune('|') {
					goto l687
				}
				position++
				if !rules[ruleskip]() {
					goto l687
				}
				depth--
				add(rulePIPE, position688)
			}
			return true
		l687:
			position, tokenIndex, depth = position687, tokenIndex687, depth687
			return false
		},
		/* 93 SLASH <- <('/' skip)> */
		func() bool {
			position689, tokenIndex689, depth689 := position, tokenIndex, depth
			{
				position690 := position
				depth++
				if buffer[position] != rune('/') {
					goto l689
				}
				position++
				if !rules[ruleskip]() {
					goto l689
				}
				depth--
				add(ruleSLASH, position690)
			}
			return true
		l689:
			position, tokenIndex, depth = position689, tokenIndex689, depth689
			return false
		},
		/* 94 INVERSE <- <('^' skip)> */
		func() bool {
			position691, tokenIndex691, depth691 := position, tokenIndex, depth
			{
				position692 := position
				depth++
				if buffer[position] != rune('^') {
					goto l691
				}
				position++
				if !rules[ruleskip]() {
					goto l691
				}
				depth--
				add(ruleINVERSE, position692)
			}
			return true
		l691:
			position, tokenIndex, depth = position691, tokenIndex691, depth691
			return false
		},
		/* 95 LPAREN <- <('(' skip)> */
		func() bool {
			position693, tokenIndex693, depth693 := position, tokenIndex, depth
			{
				position694 := position
				depth++
				if buffer[position] != rune('(') {
					goto l693
				}
				position++
				if !rules[ruleskip]() {
					goto l693
				}
				depth--
				add(ruleLPAREN, position694)
			}
			return true
		l693:
			position, tokenIndex, depth = position693, tokenIndex693, depth693
			return false
		},
		/* 96 RPAREN <- <(')' skip)> */
		func() bool {
			position695, tokenIndex695, depth695 := position, tokenIndex, depth
			{
				position696 := position
				depth++
				if buffer[position] != rune(')') {
					goto l695
				}
				position++
				if !rules[ruleskip]() {
					goto l695
				}
				depth--
				add(ruleRPAREN, position696)
			}
			return true
		l695:
			position, tokenIndex, depth = position695, tokenIndex695, depth695
			return false
		},
		/* 97 ISA <- <('a' skip)> */
		func() bool {
			position697, tokenIndex697, depth697 := position, tokenIndex, depth
			{
				position698 := position
				depth++
				if buffer[position] != rune('a') {
					goto l697
				}
				position++
				if !rules[ruleskip]() {
					goto l697
				}
				depth--
				add(ruleISA, position698)
			}
			return true
		l697:
			position, tokenIndex, depth = position697, tokenIndex697, depth697
			return false
		},
		/* 98 NOT <- <('!' skip)> */
		func() bool {
			position699, tokenIndex699, depth699 := position, tokenIndex, depth
			{
				position700 := position
				depth++
				if buffer[position] != rune('!') {
					goto l699
				}
				position++
				if !rules[ruleskip]() {
					goto l699
				}
				depth--
				add(ruleNOT, position700)
			}
			return true
		l699:
			position, tokenIndex, depth = position699, tokenIndex699, depth699
			return false
		},
		/* 99 STAR <- <('*' skip)> */
		func() bool {
			position701, tokenIndex701, depth701 := position, tokenIndex, depth
			{
				position702 := position
				depth++
				if buffer[position] != rune('*') {
					goto l701
				}
				position++
				if !rules[ruleskip]() {
					goto l701
				}
				depth--
				add(ruleSTAR, position702)
			}
			return true
		l701:
			position, tokenIndex, depth = position701, tokenIndex701, depth701
			return false
		},
		/* 100 PLUS <- <('+' skip)> */
		func() bool {
			position703, tokenIndex703, depth703 := position, tokenIndex, depth
			{
				position704 := position
				depth++
				if buffer[position] != rune('+') {
					goto l703
				}
				position++
				if !rules[ruleskip]() {
					goto l703
				}
				depth--
				add(rulePLUS, position704)
			}
			return true
		l703:
			position, tokenIndex, depth = position703, tokenIndex703, depth703
			return false
		},
		/* 101 MINUS <- <('-' skip)> */
		func() bool {
			position705, tokenIndex705, depth705 := position, tokenIndex, depth
			{
				position706 := position
				depth++
				if buffer[position] != rune('-') {
					goto l705
				}
				position++
				if !rules[ruleskip]() {
					goto l705
				}
				depth--
				add(ruleMINUS, position706)
			}
			return true
		l705:
			position, tokenIndex, depth = position705, tokenIndex705, depth705
			return false
		},
		/* 102 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 103 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 104 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 105 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 106 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position711, tokenIndex711, depth711 := position, tokenIndex, depth
			{
				position712 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l711
				}
				position++
			l713:
				{
					position714, tokenIndex714, depth714 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l714
					}
					position++
					goto l713
				l714:
					position, tokenIndex, depth = position714, tokenIndex714, depth714
				}
				if !rules[ruleskip]() {
					goto l711
				}
				depth--
				add(ruleINTEGER, position712)
			}
			return true
		l711:
			position, tokenIndex, depth = position711, tokenIndex711, depth711
			return false
		},
		/* 107 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 108 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 109 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 110 OR <- <('|' '|' skip)> */
		nil,
		/* 111 AND <- <('&' '&' skip)> */
		nil,
		/* 112 EQ <- <('=' skip)> */
		nil,
		/* 113 NE <- <('!' '=' skip)> */
		nil,
		/* 114 GT <- <('>' skip)> */
		nil,
		/* 115 LT <- <('<' skip)> */
		nil,
		/* 116 LE <- <('<' '=' skip)> */
		nil,
		/* 117 GE <- <('>' '=' skip)> */
		nil,
		/* 118 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 119 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 120 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		nil,
		/* 121 skip <- <(ws / comment)*> */
		func() bool {
			{
				position730 := position
				depth++
			l731:
				{
					position732, tokenIndex732, depth732 := position, tokenIndex, depth
					{
						position733, tokenIndex733, depth733 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l734
						}
						goto l733
					l734:
						position, tokenIndex, depth = position733, tokenIndex733, depth733
						{
							position735 := position
							depth++
							{
								position736 := position
								depth++
								if buffer[position] != rune('#') {
									goto l732
								}
								position++
							l737:
								{
									position738, tokenIndex738, depth738 := position, tokenIndex, depth
									{
										position739, tokenIndex739, depth739 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l739
										}
										goto l738
									l739:
										position, tokenIndex, depth = position739, tokenIndex739, depth739
									}
									if !matchDot() {
										goto l738
									}
									goto l737
								l738:
									position, tokenIndex, depth = position738, tokenIndex738, depth738
								}
								if !rules[ruleendOfLine]() {
									goto l732
								}
								depth--
								add(rulePegText, position736)
							}
							{
								add(ruleAction13, position)
							}
							depth--
							add(rulecomment, position735)
						}
					}
				l733:
					goto l731
				l732:
					position, tokenIndex, depth = position732, tokenIndex732, depth732
				}
				depth--
				add(ruleskip, position730)
			}
			return true
		},
		/* 122 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position741, tokenIndex741, depth741 := position, tokenIndex, depth
			{
				position742 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l741
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l741
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l741
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l741
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l741
						}
						break
					}
				}

				depth--
				add(rulews, position742)
			}
			return true
		l741:
			position, tokenIndex, depth = position741, tokenIndex741, depth741
			return false
		},
		/* 123 comment <- <(<('#' (!endOfLine .)* endOfLine)> Action13)> */
		nil,
		/* 124 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position745, tokenIndex745, depth745 := position, tokenIndex, depth
			{
				position746 := position
				depth++
				{
					position747, tokenIndex747, depth747 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l748
					}
					position++
					if buffer[position] != rune('\n') {
						goto l748
					}
					position++
					goto l747
				l748:
					position, tokenIndex, depth = position747, tokenIndex747, depth747
					if buffer[position] != rune('\n') {
						goto l749
					}
					position++
					goto l747
				l749:
					position, tokenIndex, depth = position747, tokenIndex747, depth747
					if buffer[position] != rune('\r') {
						goto l745
					}
					position++
				}
			l747:
				depth--
				add(ruleendOfLine, position746)
			}
			return true
		l745:
			position, tokenIndex, depth = position745, tokenIndex745, depth745
			return false
		},
		nil,
		/* 127 Action0 <- <{ p.addPrefix(p.skip(buffer, begin, end)) }> */
		nil,
		/* 128 Action1 <- <{ p.setSubject(p.skip(buffer, begin, end)) }> */
		nil,
		/* 129 Action2 <- <{ p.setSubject(p.skip(buffer, begin, end)) }> */
		nil,
		/* 130 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 131 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 132 Action5 <- <{ p.setPredicate(p.skip(buffer, begin, end)) }> */
		nil,
		/* 133 Action6 <- <{ p.setPredicate(p.skip(buffer, begin, end)) }> */
		nil,
		/* 134 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 135 Action8 <- <{ p.setObject(p.skip(buffer, begin, end)); p.addTriplePattern() }> */
		nil,
		/* 136 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 137 Action10 <- <{ p.setPrefix(p.skip(buffer, begin, end)) }> */
		nil,
		/* 138 Action11 <- <{ p.setPathLength(p.skip(buffer, begin, end)) }> */
		nil,
		/* 139 Action12 <- <{ p.setKeyword(p.skip(buffer, begin, end)) }> */
		nil,
		/* 140 Action13 <- <{ p.commentBegin = begin }> */
		nil,
	}
	p.rules = rules
}
