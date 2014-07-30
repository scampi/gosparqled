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
	rulefunctionCall
	rulein
	rulenotin
	ruleargList
	rulebuiltinCall
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
	ruleSTR
	ruleLANG
	ruleDATATYPE
	ruleIRI
	ruleURI
	ruleABS
	ruleCEIL
	ruleROUND
	ruleFLOOR
	ruleSTRLEN
	ruleUCASE
	ruleLCASE
	ruleENCODEFORURI
	ruleYEAR
	ruleMONTH
	ruleDAY
	ruleHOURS
	ruleMINUTES
	ruleSECONDS
	ruleTIMEZONE
	ruleTZ
	ruleMD5
	ruleSHA1
	ruleSHA256
	ruleSHA384
	ruleSHA512
	ruleISIRI
	ruleISURI
	ruleISBLANK
	ruleISLITERAL
	ruleISNUMERIC
	ruleLANGMATCHES
	ruleCONTAINS
	ruleSTRSTARTS
	ruleSTRENDS
	ruleSTRBEFORE
	ruleSTRAFTER
	ruleSTRLANG
	ruleSTRDT
	ruleSAMETERM
	ruleBOUND
	ruleBNODE
	ruleRAND
	ruleNOW
	ruleUUID
	ruleSTRUUID
	ruleCONCAT
	ruleSUBSTR
	ruleREPLACE
	ruleREGEX
	ruleIF
	ruleEXISTS
	ruleNOTEXIST
	ruleCOALESCE
	ruleskip
	rulews
	rulecomment
	ruleendOfLine

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
	"functionCall",
	"in",
	"notin",
	"argList",
	"builtinCall",
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
	"STR",
	"LANG",
	"DATATYPE",
	"IRI",
	"URI",
	"ABS",
	"CEIL",
	"ROUND",
	"FLOOR",
	"STRLEN",
	"UCASE",
	"LCASE",
	"ENCODEFORURI",
	"YEAR",
	"MONTH",
	"DAY",
	"HOURS",
	"MINUTES",
	"SECONDS",
	"TIMEZONE",
	"TZ",
	"MD5",
	"SHA1",
	"SHA256",
	"SHA384",
	"SHA512",
	"ISIRI",
	"ISURI",
	"ISBLANK",
	"ISLITERAL",
	"ISNUMERIC",
	"LANGMATCHES",
	"CONTAINS",
	"STRSTARTS",
	"STRENDS",
	"STRBEFORE",
	"STRAFTER",
	"STRLANG",
	"STRDT",
	"SAMETERM",
	"BOUND",
	"BNODE",
	"RAND",
	"NOW",
	"UUID",
	"STRUUID",
	"CONCAT",
	"SUBSTR",
	"REPLACE",
	"REGEX",
	"IF",
	"EXISTS",
	"NOTEXIST",
	"COALESCE",
	"skip",
	"ws",
	"comment",
	"endOfLine",

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
	rules  [181]func() bool
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
									if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
										if !rules[ruleskip]() {
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
										if !rules[ruleskip]() {
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
					if !rules[ruleskip]() {
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
							if !rules[ruleskip]() {
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
							if !rules[ruleskip]() {
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
								if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
					if !rules[ruleskip]() {
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
						if !rules[ruleskip]() {
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
						if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
						if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
									if !rules[ruleskip]() {
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
												if !rules[ruleskip]() {
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
											l320:
												{
													position321, tokenIndex321, depth321 := position, tokenIndex, depth
													if !rules[rulews]() {
														goto l321
													}
													goto l320
												l321:
													position, tokenIndex, depth = position321, tokenIndex321, depth321
												}
												if buffer[position] != rune(']') {
													goto l300
												}
												position++
												if !rules[ruleskip]() {
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
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				{
					position328, tokenIndex328, depth328 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l329
					}
					goto l328
				l329:
					position, tokenIndex, depth = position328, tokenIndex328, depth328
					{
						position330 := position
						depth++
						if !rules[rulepath]() {
							goto l326
						}
						depth--
						add(ruleverbPath, position330)
					}
				}
			l328:
				if !rules[ruleobjectListPath]() {
					goto l326
				}
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					{
						position333 := position
						depth++
						if buffer[position] != rune(';') {
							goto l331
						}
						position++
						if !rules[ruleskip]() {
							goto l331
						}
						depth--
						add(ruleSEMICOLON, position333)
					}
					if !rules[rulepropertyListPath]() {
						goto l331
					}
					goto l332
				l331:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
				}
			l332:
				depth--
				add(rulepropertyListPath, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 30 verbPath <- <path> */
		nil,
		/* 31 path <- <pathAlternative> */
		func() bool {
			position335, tokenIndex335, depth335 := position, tokenIndex, depth
			{
				position336 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l335
				}
				depth--
				add(rulepath, position336)
			}
			return true
		l335:
			position, tokenIndex, depth = position335, tokenIndex335, depth335
			return false
		},
		/* 32 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l337
				}
			l339:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l340
					}
					if !rules[rulepathAlternative]() {
						goto l340
					}
					goto l339
				l340:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
				}
				depth--
				add(rulepathAlternative, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 33 pathSequence <- <(pathElt (SLASH pathSequence)*)> */
		func() bool {
			position341, tokenIndex341, depth341 := position, tokenIndex, depth
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
							position347, tokenIndex347, depth347 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l348
							}
							goto l347
						l348:
							position, tokenIndex, depth = position347, tokenIndex347, depth347
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l341
									}
									if !rules[rulepath]() {
										goto l341
									}
									if !rules[ruleRPAREN]() {
										goto l341
									}
									break
								case '!':
									if !rules[ruleNOT]() {
										goto l341
									}
									{
										position350 := position
										depth++
										{
											position351, tokenIndex351, depth351 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l352
											}
											goto l351
										l352:
											position, tokenIndex, depth = position351, tokenIndex351, depth351
											if !rules[ruleLPAREN]() {
												goto l341
											}
											{
												position353, tokenIndex353, depth353 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l353
												}
											l355:
												{
													position356, tokenIndex356, depth356 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l356
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l356
													}
													goto l355
												l356:
													position, tokenIndex, depth = position356, tokenIndex356, depth356
												}
												goto l354
											l353:
												position, tokenIndex, depth = position353, tokenIndex353, depth353
											}
										l354:
											if !rules[ruleRPAREN]() {
												goto l341
											}
										}
									l351:
										depth--
										add(rulepathNegatedPropertySet, position350)
									}
									break
								default:
									if !rules[ruleISA]() {
										goto l341
									}
									break
								}
							}

						}
					l347:
						depth--
						add(rulepathPrimary, position346)
					}
					depth--
					add(rulepathElt, position343)
				}
			l357:
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l358
					}
					if !rules[rulepathSequence]() {
						goto l358
					}
					goto l357
				l358:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
				}
				depth--
				add(rulepathSequence, position342)
			}
			return true
		l341:
			position, tokenIndex, depth = position341, tokenIndex341, depth341
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
		/* 38 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position369, tokenIndex369, depth369 := position, tokenIndex, depth
			{
				position370 := position
				depth++
				{
					position371 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l369
					}
					depth--
					add(ruleobjectPath, position371)
				}
			l372:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if !rules[ruleCOMMA]() {
						goto l373
					}
					if !rules[ruleobjectListPath]() {
						goto l373
					}
					goto l372
				l373:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
				}
				depth--
				add(ruleobjectListPath, position370)
			}
			return true
		l369:
			position, tokenIndex, depth = position369, tokenIndex369, depth369
			return false
		},
		/* 39 objectPath <- <graphNodePath> */
		nil,
		/* 40 graphNodePath <- <varOrTerm> */
		func() bool {
			position375, tokenIndex375, depth375 := position, tokenIndex, depth
			{
				position376 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l375
				}
				depth--
				add(rulegraphNodePath, position376)
			}
			return true
		l375:
			position, tokenIndex, depth = position375, tokenIndex375, depth375
			return false
		},
		/* 41 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position378 := position
				depth++
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					{
						position381 := position
						depth++
						{
							position382, tokenIndex382, depth382 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l383
							}
							{
								position384, tokenIndex384, depth384 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l384
								}
								goto l385
							l384:
								position, tokenIndex, depth = position384, tokenIndex384, depth384
							}
						l385:
							goto l382
						l383:
							position, tokenIndex, depth = position382, tokenIndex382, depth382
							if !rules[ruleoffset]() {
								goto l379
							}
							{
								position386, tokenIndex386, depth386 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l386
								}
								goto l387
							l386:
								position, tokenIndex, depth = position386, tokenIndex386, depth386
							}
						l387:
						}
					l382:
						depth--
						add(rulelimitOffsetClauses, position381)
					}
					goto l380
				l379:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
				}
			l380:
				depth--
				add(rulesolutionModifier, position378)
			}
			return true
		},
		/* 42 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 43 limit <- <(LIMIT INTEGER)> */
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
					if !rules[ruleskip]() {
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
		/* 44 offset <- <(OFFSET INTEGER)> */
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
					if !rules[ruleskip]() {
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
		/* 45 expression <- <conditionalOrExpression> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l417
				}
				depth--
				add(ruleexpression, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 46 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position419, tokenIndex419, depth419 := position, tokenIndex, depth
			{
				position420 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l419
				}
				{
					position421, tokenIndex421, depth421 := position, tokenIndex, depth
					{
						position423 := position
						depth++
						if buffer[position] != rune('|') {
							goto l421
						}
						position++
						if buffer[position] != rune('|') {
							goto l421
						}
						position++
						if !rules[ruleskip]() {
							goto l421
						}
						depth--
						add(ruleOR, position423)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l421
					}
					goto l422
				l421:
					position, tokenIndex, depth = position421, tokenIndex421, depth421
				}
			l422:
				depth--
				add(ruleconditionalOrExpression, position420)
			}
			return true
		l419:
			position, tokenIndex, depth = position419, tokenIndex419, depth419
			return false
		},
		/* 47 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position424, tokenIndex424, depth424 := position, tokenIndex, depth
			{
				position425 := position
				depth++
				{
					position426 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l424
					}
					{
						position427, tokenIndex427, depth427 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position430 := position
									depth++
									{
										position431 := position
										depth++
										{
											position432, tokenIndex432, depth432 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l433
											}
											position++
											goto l432
										l433:
											position, tokenIndex, depth = position432, tokenIndex432, depth432
											if buffer[position] != rune('N') {
												goto l427
											}
											position++
										}
									l432:
										{
											position434, tokenIndex434, depth434 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l435
											}
											position++
											goto l434
										l435:
											position, tokenIndex, depth = position434, tokenIndex434, depth434
											if buffer[position] != rune('O') {
												goto l427
											}
											position++
										}
									l434:
										{
											position436, tokenIndex436, depth436 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l437
											}
											position++
											goto l436
										l437:
											position, tokenIndex, depth = position436, tokenIndex436, depth436
											if buffer[position] != rune('T') {
												goto l427
											}
											position++
										}
									l436:
										if buffer[position] != rune(' ') {
											goto l427
										}
										position++
										{
											position438, tokenIndex438, depth438 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l439
											}
											position++
											goto l438
										l439:
											position, tokenIndex, depth = position438, tokenIndex438, depth438
											if buffer[position] != rune('I') {
												goto l427
											}
											position++
										}
									l438:
										{
											position440, tokenIndex440, depth440 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l441
											}
											position++
											goto l440
										l441:
											position, tokenIndex, depth = position440, tokenIndex440, depth440
											if buffer[position] != rune('N') {
												goto l427
											}
											position++
										}
									l440:
										if !rules[ruleskip]() {
											goto l427
										}
										depth--
										add(ruleNOTIN, position431)
									}
									if !rules[ruleargList]() {
										goto l427
									}
									depth--
									add(rulenotin, position430)
								}
								break
							case 'I', 'i':
								{
									position442 := position
									depth++
									{
										position443 := position
										depth++
										{
											position444, tokenIndex444, depth444 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l445
											}
											position++
											goto l444
										l445:
											position, tokenIndex, depth = position444, tokenIndex444, depth444
											if buffer[position] != rune('I') {
												goto l427
											}
											position++
										}
									l444:
										{
											position446, tokenIndex446, depth446 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l447
											}
											position++
											goto l446
										l447:
											position, tokenIndex, depth = position446, tokenIndex446, depth446
											if buffer[position] != rune('N') {
												goto l427
											}
											position++
										}
									l446:
										if !rules[ruleskip]() {
											goto l427
										}
										depth--
										add(ruleIN, position443)
									}
									if !rules[ruleargList]() {
										goto l427
									}
									depth--
									add(rulein, position442)
								}
								break
							default:
								{
									position448, tokenIndex448, depth448 := position, tokenIndex, depth
									{
										position450 := position
										depth++
										if buffer[position] != rune('<') {
											goto l449
										}
										position++
										if !rules[ruleskip]() {
											goto l449
										}
										depth--
										add(ruleLT, position450)
									}
									goto l448
								l449:
									position, tokenIndex, depth = position448, tokenIndex448, depth448
									{
										position452 := position
										depth++
										if buffer[position] != rune('>') {
											goto l451
										}
										position++
										if buffer[position] != rune('=') {
											goto l451
										}
										position++
										if !rules[ruleskip]() {
											goto l451
										}
										depth--
										add(ruleGE, position452)
									}
									goto l448
								l451:
									position, tokenIndex, depth = position448, tokenIndex448, depth448
									{
										switch buffer[position] {
										case '>':
											{
												position454 := position
												depth++
												if buffer[position] != rune('>') {
													goto l427
												}
												position++
												if !rules[ruleskip]() {
													goto l427
												}
												depth--
												add(ruleGT, position454)
											}
											break
										case '<':
											{
												position455 := position
												depth++
												if buffer[position] != rune('<') {
													goto l427
												}
												position++
												if buffer[position] != rune('=') {
													goto l427
												}
												position++
												if !rules[ruleskip]() {
													goto l427
												}
												depth--
												add(ruleLE, position455)
											}
											break
										case '!':
											{
												position456 := position
												depth++
												if buffer[position] != rune('!') {
													goto l427
												}
												position++
												if buffer[position] != rune('=') {
													goto l427
												}
												position++
												if !rules[ruleskip]() {
													goto l427
												}
												depth--
												add(ruleNE, position456)
											}
											break
										default:
											{
												position457 := position
												depth++
												if buffer[position] != rune('=') {
													goto l427
												}
												position++
												if !rules[ruleskip]() {
													goto l427
												}
												depth--
												add(ruleEQ, position457)
											}
											break
										}
									}

								}
							l448:
								if !rules[rulenumericExpression]() {
									goto l427
								}
								break
							}
						}

						goto l428
					l427:
						position, tokenIndex, depth = position427, tokenIndex427, depth427
					}
				l428:
					depth--
					add(rulevalueLogical, position426)
				}
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					{
						position460 := position
						depth++
						if buffer[position] != rune('&') {
							goto l458
						}
						position++
						if buffer[position] != rune('&') {
							goto l458
						}
						position++
						if !rules[ruleskip]() {
							goto l458
						}
						depth--
						add(ruleAND, position460)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l458
					}
					goto l459
				l458:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
				}
			l459:
				depth--
				add(ruleconditionalAndExpression, position425)
			}
			return true
		l424:
			position, tokenIndex, depth = position424, tokenIndex424, depth424
			return false
		},
		/* 48 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 49 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position462, tokenIndex462, depth462 := position, tokenIndex, depth
			{
				position463 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l462
				}
			l464:
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						{
							position468, tokenIndex468, depth468 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l469
							}
							goto l468
						l469:
							position, tokenIndex, depth = position468, tokenIndex468, depth468
							if !rules[ruleMINUS]() {
								goto l467
							}
						}
					l468:
						if !rules[rulemultiplicativeExpression]() {
							goto l467
						}
						goto l466
					l467:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
						{
							position470 := position
							depth++
							{
								position471, tokenIndex471, depth471 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l472
								}
								position++
								goto l471
							l472:
								position, tokenIndex, depth = position471, tokenIndex471, depth471
								if buffer[position] != rune('-') {
									goto l465
								}
								position++
							}
						l471:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l465
							}
							position++
						l473:
							{
								position474, tokenIndex474, depth474 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l474
								}
								position++
								goto l473
							l474:
								position, tokenIndex, depth = position474, tokenIndex474, depth474
							}
							{
								position475, tokenIndex475, depth475 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l475
								}
								position++
							l477:
								{
									position478, tokenIndex478, depth478 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l478
									}
									position++
									goto l477
								l478:
									position, tokenIndex, depth = position478, tokenIndex478, depth478
								}
								goto l476
							l475:
								position, tokenIndex, depth = position475, tokenIndex475, depth475
							}
						l476:
							if !rules[ruleskip]() {
								goto l465
							}
							depth--
							add(rulesignedNumericLiteral, position470)
						}
					}
				l466:
					goto l464
				l465:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
				}
				depth--
				add(rulenumericExpression, position463)
			}
			return true
		l462:
			position, tokenIndex, depth = position462, tokenIndex462, depth462
			return false
		},
		/* 50 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l479
				}
			l481:
				{
					position482, tokenIndex482, depth482 := position, tokenIndex, depth
					{
						position483, tokenIndex483, depth483 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l484
						}
						goto l483
					l484:
						position, tokenIndex, depth = position483, tokenIndex483, depth483
						if !rules[ruleSLASH]() {
							goto l482
						}
					}
				l483:
					if !rules[ruleunaryExpression]() {
						goto l482
					}
					goto l481
				l482:
					position, tokenIndex, depth = position482, tokenIndex482, depth482
				}
				depth--
				add(rulemultiplicativeExpression, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 51 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position485, tokenIndex485, depth485 := position, tokenIndex, depth
			{
				position486 := position
				depth++
				{
					position487, tokenIndex487, depth487 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l487
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l487
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l487
							}
							break
						}
					}

					goto l488
				l487:
					position, tokenIndex, depth = position487, tokenIndex487, depth487
				}
			l488:
				{
					position490 := position
					depth++
					{
						position491, tokenIndex491, depth491 := position, tokenIndex, depth
						{
							position493 := position
							depth++
							if !rules[ruleLPAREN]() {
								goto l492
							}
							if !rules[ruleexpression]() {
								goto l492
							}
							if !rules[ruleRPAREN]() {
								goto l492
							}
							depth--
							add(rulebrackettedExpression, position493)
						}
						goto l491
					l492:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
						{
							position495 := position
							depth++
							{
								position496, tokenIndex496, depth496 := position, tokenIndex, depth
								{
									position498, tokenIndex498, depth498 := position, tokenIndex, depth
									{
										position500 := position
										depth++
										{
											position501, tokenIndex501, depth501 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l502
											}
											position++
											goto l501
										l502:
											position, tokenIndex, depth = position501, tokenIndex501, depth501
											if buffer[position] != rune('S') {
												goto l499
											}
											position++
										}
									l501:
										{
											position503, tokenIndex503, depth503 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l504
											}
											position++
											goto l503
										l504:
											position, tokenIndex, depth = position503, tokenIndex503, depth503
											if buffer[position] != rune('T') {
												goto l499
											}
											position++
										}
									l503:
										{
											position505, tokenIndex505, depth505 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l506
											}
											position++
											goto l505
										l506:
											position, tokenIndex, depth = position505, tokenIndex505, depth505
											if buffer[position] != rune('R') {
												goto l499
											}
											position++
										}
									l505:
										if !rules[ruleskip]() {
											goto l499
										}
										depth--
										add(ruleSTR, position500)
									}
									goto l498
								l499:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position508 := position
										depth++
										{
											position509, tokenIndex509, depth509 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l510
											}
											position++
											goto l509
										l510:
											position, tokenIndex, depth = position509, tokenIndex509, depth509
											if buffer[position] != rune('L') {
												goto l507
											}
											position++
										}
									l509:
										{
											position511, tokenIndex511, depth511 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l512
											}
											position++
											goto l511
										l512:
											position, tokenIndex, depth = position511, tokenIndex511, depth511
											if buffer[position] != rune('A') {
												goto l507
											}
											position++
										}
									l511:
										{
											position513, tokenIndex513, depth513 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l514
											}
											position++
											goto l513
										l514:
											position, tokenIndex, depth = position513, tokenIndex513, depth513
											if buffer[position] != rune('N') {
												goto l507
											}
											position++
										}
									l513:
										{
											position515, tokenIndex515, depth515 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l516
											}
											position++
											goto l515
										l516:
											position, tokenIndex, depth = position515, tokenIndex515, depth515
											if buffer[position] != rune('G') {
												goto l507
											}
											position++
										}
									l515:
										if !rules[ruleskip]() {
											goto l507
										}
										depth--
										add(ruleLANG, position508)
									}
									goto l498
								l507:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position518 := position
										depth++
										{
											position519, tokenIndex519, depth519 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l520
											}
											position++
											goto l519
										l520:
											position, tokenIndex, depth = position519, tokenIndex519, depth519
											if buffer[position] != rune('D') {
												goto l517
											}
											position++
										}
									l519:
										{
											position521, tokenIndex521, depth521 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l522
											}
											position++
											goto l521
										l522:
											position, tokenIndex, depth = position521, tokenIndex521, depth521
											if buffer[position] != rune('A') {
												goto l517
											}
											position++
										}
									l521:
										{
											position523, tokenIndex523, depth523 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l524
											}
											position++
											goto l523
										l524:
											position, tokenIndex, depth = position523, tokenIndex523, depth523
											if buffer[position] != rune('T') {
												goto l517
											}
											position++
										}
									l523:
										{
											position525, tokenIndex525, depth525 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l526
											}
											position++
											goto l525
										l526:
											position, tokenIndex, depth = position525, tokenIndex525, depth525
											if buffer[position] != rune('A') {
												goto l517
											}
											position++
										}
									l525:
										{
											position527, tokenIndex527, depth527 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l528
											}
											position++
											goto l527
										l528:
											position, tokenIndex, depth = position527, tokenIndex527, depth527
											if buffer[position] != rune('T') {
												goto l517
											}
											position++
										}
									l527:
										{
											position529, tokenIndex529, depth529 := position, tokenIndex, depth
											if buffer[position] != rune('y') {
												goto l530
											}
											position++
											goto l529
										l530:
											position, tokenIndex, depth = position529, tokenIndex529, depth529
											if buffer[position] != rune('Y') {
												goto l517
											}
											position++
										}
									l529:
										{
											position531, tokenIndex531, depth531 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l532
											}
											position++
											goto l531
										l532:
											position, tokenIndex, depth = position531, tokenIndex531, depth531
											if buffer[position] != rune('P') {
												goto l517
											}
											position++
										}
									l531:
										{
											position533, tokenIndex533, depth533 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l534
											}
											position++
											goto l533
										l534:
											position, tokenIndex, depth = position533, tokenIndex533, depth533
											if buffer[position] != rune('E') {
												goto l517
											}
											position++
										}
									l533:
										if !rules[ruleskip]() {
											goto l517
										}
										depth--
										add(ruleDATATYPE, position518)
									}
									goto l498
								l517:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position536 := position
										depth++
										{
											position537, tokenIndex537, depth537 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l538
											}
											position++
											goto l537
										l538:
											position, tokenIndex, depth = position537, tokenIndex537, depth537
											if buffer[position] != rune('I') {
												goto l535
											}
											position++
										}
									l537:
										{
											position539, tokenIndex539, depth539 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l540
											}
											position++
											goto l539
										l540:
											position, tokenIndex, depth = position539, tokenIndex539, depth539
											if buffer[position] != rune('R') {
												goto l535
											}
											position++
										}
									l539:
										{
											position541, tokenIndex541, depth541 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l542
											}
											position++
											goto l541
										l542:
											position, tokenIndex, depth = position541, tokenIndex541, depth541
											if buffer[position] != rune('I') {
												goto l535
											}
											position++
										}
									l541:
										if !rules[ruleskip]() {
											goto l535
										}
										depth--
										add(ruleIRI, position536)
									}
									goto l498
								l535:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position544 := position
										depth++
										{
											position545, tokenIndex545, depth545 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l546
											}
											position++
											goto l545
										l546:
											position, tokenIndex, depth = position545, tokenIndex545, depth545
											if buffer[position] != rune('U') {
												goto l543
											}
											position++
										}
									l545:
										{
											position547, tokenIndex547, depth547 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l548
											}
											position++
											goto l547
										l548:
											position, tokenIndex, depth = position547, tokenIndex547, depth547
											if buffer[position] != rune('R') {
												goto l543
											}
											position++
										}
									l547:
										{
											position549, tokenIndex549, depth549 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l550
											}
											position++
											goto l549
										l550:
											position, tokenIndex, depth = position549, tokenIndex549, depth549
											if buffer[position] != rune('I') {
												goto l543
											}
											position++
										}
									l549:
										if !rules[ruleskip]() {
											goto l543
										}
										depth--
										add(ruleURI, position544)
									}
									goto l498
								l543:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position552 := position
										depth++
										{
											position553, tokenIndex553, depth553 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l554
											}
											position++
											goto l553
										l554:
											position, tokenIndex, depth = position553, tokenIndex553, depth553
											if buffer[position] != rune('S') {
												goto l551
											}
											position++
										}
									l553:
										{
											position555, tokenIndex555, depth555 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l556
											}
											position++
											goto l555
										l556:
											position, tokenIndex, depth = position555, tokenIndex555, depth555
											if buffer[position] != rune('T') {
												goto l551
											}
											position++
										}
									l555:
										{
											position557, tokenIndex557, depth557 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l558
											}
											position++
											goto l557
										l558:
											position, tokenIndex, depth = position557, tokenIndex557, depth557
											if buffer[position] != rune('R') {
												goto l551
											}
											position++
										}
									l557:
										{
											position559, tokenIndex559, depth559 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l560
											}
											position++
											goto l559
										l560:
											position, tokenIndex, depth = position559, tokenIndex559, depth559
											if buffer[position] != rune('L') {
												goto l551
											}
											position++
										}
									l559:
										{
											position561, tokenIndex561, depth561 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l562
											}
											position++
											goto l561
										l562:
											position, tokenIndex, depth = position561, tokenIndex561, depth561
											if buffer[position] != rune('E') {
												goto l551
											}
											position++
										}
									l561:
										{
											position563, tokenIndex563, depth563 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l564
											}
											position++
											goto l563
										l564:
											position, tokenIndex, depth = position563, tokenIndex563, depth563
											if buffer[position] != rune('N') {
												goto l551
											}
											position++
										}
									l563:
										if !rules[ruleskip]() {
											goto l551
										}
										depth--
										add(ruleSTRLEN, position552)
									}
									goto l498
								l551:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position566 := position
										depth++
										{
											position567, tokenIndex567, depth567 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l568
											}
											position++
											goto l567
										l568:
											position, tokenIndex, depth = position567, tokenIndex567, depth567
											if buffer[position] != rune('M') {
												goto l565
											}
											position++
										}
									l567:
										{
											position569, tokenIndex569, depth569 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l570
											}
											position++
											goto l569
										l570:
											position, tokenIndex, depth = position569, tokenIndex569, depth569
											if buffer[position] != rune('O') {
												goto l565
											}
											position++
										}
									l569:
										{
											position571, tokenIndex571, depth571 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l572
											}
											position++
											goto l571
										l572:
											position, tokenIndex, depth = position571, tokenIndex571, depth571
											if buffer[position] != rune('N') {
												goto l565
											}
											position++
										}
									l571:
										{
											position573, tokenIndex573, depth573 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l574
											}
											position++
											goto l573
										l574:
											position, tokenIndex, depth = position573, tokenIndex573, depth573
											if buffer[position] != rune('T') {
												goto l565
											}
											position++
										}
									l573:
										{
											position575, tokenIndex575, depth575 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l576
											}
											position++
											goto l575
										l576:
											position, tokenIndex, depth = position575, tokenIndex575, depth575
											if buffer[position] != rune('H') {
												goto l565
											}
											position++
										}
									l575:
										if !rules[ruleskip]() {
											goto l565
										}
										depth--
										add(ruleMONTH, position566)
									}
									goto l498
								l565:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position578 := position
										depth++
										{
											position579, tokenIndex579, depth579 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l580
											}
											position++
											goto l579
										l580:
											position, tokenIndex, depth = position579, tokenIndex579, depth579
											if buffer[position] != rune('M') {
												goto l577
											}
											position++
										}
									l579:
										{
											position581, tokenIndex581, depth581 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l582
											}
											position++
											goto l581
										l582:
											position, tokenIndex, depth = position581, tokenIndex581, depth581
											if buffer[position] != rune('I') {
												goto l577
											}
											position++
										}
									l581:
										{
											position583, tokenIndex583, depth583 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l584
											}
											position++
											goto l583
										l584:
											position, tokenIndex, depth = position583, tokenIndex583, depth583
											if buffer[position] != rune('N') {
												goto l577
											}
											position++
										}
									l583:
										{
											position585, tokenIndex585, depth585 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l586
											}
											position++
											goto l585
										l586:
											position, tokenIndex, depth = position585, tokenIndex585, depth585
											if buffer[position] != rune('U') {
												goto l577
											}
											position++
										}
									l585:
										{
											position587, tokenIndex587, depth587 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l588
											}
											position++
											goto l587
										l588:
											position, tokenIndex, depth = position587, tokenIndex587, depth587
											if buffer[position] != rune('T') {
												goto l577
											}
											position++
										}
									l587:
										{
											position589, tokenIndex589, depth589 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l590
											}
											position++
											goto l589
										l590:
											position, tokenIndex, depth = position589, tokenIndex589, depth589
											if buffer[position] != rune('E') {
												goto l577
											}
											position++
										}
									l589:
										{
											position591, tokenIndex591, depth591 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l592
											}
											position++
											goto l591
										l592:
											position, tokenIndex, depth = position591, tokenIndex591, depth591
											if buffer[position] != rune('S') {
												goto l577
											}
											position++
										}
									l591:
										if !rules[ruleskip]() {
											goto l577
										}
										depth--
										add(ruleMINUTES, position578)
									}
									goto l498
								l577:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position594 := position
										depth++
										{
											position595, tokenIndex595, depth595 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l596
											}
											position++
											goto l595
										l596:
											position, tokenIndex, depth = position595, tokenIndex595, depth595
											if buffer[position] != rune('S') {
												goto l593
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
												goto l593
											}
											position++
										}
									l597:
										{
											position599, tokenIndex599, depth599 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l600
											}
											position++
											goto l599
										l600:
											position, tokenIndex, depth = position599, tokenIndex599, depth599
											if buffer[position] != rune('C') {
												goto l593
											}
											position++
										}
									l599:
										{
											position601, tokenIndex601, depth601 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l602
											}
											position++
											goto l601
										l602:
											position, tokenIndex, depth = position601, tokenIndex601, depth601
											if buffer[position] != rune('O') {
												goto l593
											}
											position++
										}
									l601:
										{
											position603, tokenIndex603, depth603 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l604
											}
											position++
											goto l603
										l604:
											position, tokenIndex, depth = position603, tokenIndex603, depth603
											if buffer[position] != rune('N') {
												goto l593
											}
											position++
										}
									l603:
										{
											position605, tokenIndex605, depth605 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l606
											}
											position++
											goto l605
										l606:
											position, tokenIndex, depth = position605, tokenIndex605, depth605
											if buffer[position] != rune('D') {
												goto l593
											}
											position++
										}
									l605:
										{
											position607, tokenIndex607, depth607 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l608
											}
											position++
											goto l607
										l608:
											position, tokenIndex, depth = position607, tokenIndex607, depth607
											if buffer[position] != rune('S') {
												goto l593
											}
											position++
										}
									l607:
										if !rules[ruleskip]() {
											goto l593
										}
										depth--
										add(ruleSECONDS, position594)
									}
									goto l498
								l593:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position610 := position
										depth++
										{
											position611, tokenIndex611, depth611 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l612
											}
											position++
											goto l611
										l612:
											position, tokenIndex, depth = position611, tokenIndex611, depth611
											if buffer[position] != rune('T') {
												goto l609
											}
											position++
										}
									l611:
										{
											position613, tokenIndex613, depth613 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l614
											}
											position++
											goto l613
										l614:
											position, tokenIndex, depth = position613, tokenIndex613, depth613
											if buffer[position] != rune('I') {
												goto l609
											}
											position++
										}
									l613:
										{
											position615, tokenIndex615, depth615 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l616
											}
											position++
											goto l615
										l616:
											position, tokenIndex, depth = position615, tokenIndex615, depth615
											if buffer[position] != rune('M') {
												goto l609
											}
											position++
										}
									l615:
										{
											position617, tokenIndex617, depth617 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l618
											}
											position++
											goto l617
										l618:
											position, tokenIndex, depth = position617, tokenIndex617, depth617
											if buffer[position] != rune('E') {
												goto l609
											}
											position++
										}
									l617:
										{
											position619, tokenIndex619, depth619 := position, tokenIndex, depth
											if buffer[position] != rune('z') {
												goto l620
											}
											position++
											goto l619
										l620:
											position, tokenIndex, depth = position619, tokenIndex619, depth619
											if buffer[position] != rune('Z') {
												goto l609
											}
											position++
										}
									l619:
										{
											position621, tokenIndex621, depth621 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l622
											}
											position++
											goto l621
										l622:
											position, tokenIndex, depth = position621, tokenIndex621, depth621
											if buffer[position] != rune('O') {
												goto l609
											}
											position++
										}
									l621:
										{
											position623, tokenIndex623, depth623 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l624
											}
											position++
											goto l623
										l624:
											position, tokenIndex, depth = position623, tokenIndex623, depth623
											if buffer[position] != rune('N') {
												goto l609
											}
											position++
										}
									l623:
										{
											position625, tokenIndex625, depth625 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l626
											}
											position++
											goto l625
										l626:
											position, tokenIndex, depth = position625, tokenIndex625, depth625
											if buffer[position] != rune('E') {
												goto l609
											}
											position++
										}
									l625:
										if !rules[ruleskip]() {
											goto l609
										}
										depth--
										add(ruleTIMEZONE, position610)
									}
									goto l498
								l609:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position628 := position
										depth++
										{
											position629, tokenIndex629, depth629 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l630
											}
											position++
											goto l629
										l630:
											position, tokenIndex, depth = position629, tokenIndex629, depth629
											if buffer[position] != rune('S') {
												goto l627
											}
											position++
										}
									l629:
										{
											position631, tokenIndex631, depth631 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l632
											}
											position++
											goto l631
										l632:
											position, tokenIndex, depth = position631, tokenIndex631, depth631
											if buffer[position] != rune('H') {
												goto l627
											}
											position++
										}
									l631:
										{
											position633, tokenIndex633, depth633 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l634
											}
											position++
											goto l633
										l634:
											position, tokenIndex, depth = position633, tokenIndex633, depth633
											if buffer[position] != rune('A') {
												goto l627
											}
											position++
										}
									l633:
										if buffer[position] != rune('1') {
											goto l627
										}
										position++
										if !rules[ruleskip]() {
											goto l627
										}
										depth--
										add(ruleSHA1, position628)
									}
									goto l498
								l627:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position636 := position
										depth++
										{
											position637, tokenIndex637, depth637 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l638
											}
											position++
											goto l637
										l638:
											position, tokenIndex, depth = position637, tokenIndex637, depth637
											if buffer[position] != rune('S') {
												goto l635
											}
											position++
										}
									l637:
										{
											position639, tokenIndex639, depth639 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l640
											}
											position++
											goto l639
										l640:
											position, tokenIndex, depth = position639, tokenIndex639, depth639
											if buffer[position] != rune('H') {
												goto l635
											}
											position++
										}
									l639:
										{
											position641, tokenIndex641, depth641 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l642
											}
											position++
											goto l641
										l642:
											position, tokenIndex, depth = position641, tokenIndex641, depth641
											if buffer[position] != rune('A') {
												goto l635
											}
											position++
										}
									l641:
										if buffer[position] != rune('2') {
											goto l635
										}
										position++
										if buffer[position] != rune('5') {
											goto l635
										}
										position++
										if buffer[position] != rune('6') {
											goto l635
										}
										position++
										if !rules[ruleskip]() {
											goto l635
										}
										depth--
										add(ruleSHA256, position636)
									}
									goto l498
								l635:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position644 := position
										depth++
										{
											position645, tokenIndex645, depth645 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l646
											}
											position++
											goto l645
										l646:
											position, tokenIndex, depth = position645, tokenIndex645, depth645
											if buffer[position] != rune('S') {
												goto l643
											}
											position++
										}
									l645:
										{
											position647, tokenIndex647, depth647 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l648
											}
											position++
											goto l647
										l648:
											position, tokenIndex, depth = position647, tokenIndex647, depth647
											if buffer[position] != rune('H') {
												goto l643
											}
											position++
										}
									l647:
										{
											position649, tokenIndex649, depth649 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l650
											}
											position++
											goto l649
										l650:
											position, tokenIndex, depth = position649, tokenIndex649, depth649
											if buffer[position] != rune('A') {
												goto l643
											}
											position++
										}
									l649:
										if buffer[position] != rune('3') {
											goto l643
										}
										position++
										if buffer[position] != rune('8') {
											goto l643
										}
										position++
										if buffer[position] != rune('4') {
											goto l643
										}
										position++
										if !rules[ruleskip]() {
											goto l643
										}
										depth--
										add(ruleSHA384, position644)
									}
									goto l498
								l643:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position652 := position
										depth++
										{
											position653, tokenIndex653, depth653 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l654
											}
											position++
											goto l653
										l654:
											position, tokenIndex, depth = position653, tokenIndex653, depth653
											if buffer[position] != rune('I') {
												goto l651
											}
											position++
										}
									l653:
										{
											position655, tokenIndex655, depth655 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l656
											}
											position++
											goto l655
										l656:
											position, tokenIndex, depth = position655, tokenIndex655, depth655
											if buffer[position] != rune('S') {
												goto l651
											}
											position++
										}
									l655:
										{
											position657, tokenIndex657, depth657 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l658
											}
											position++
											goto l657
										l658:
											position, tokenIndex, depth = position657, tokenIndex657, depth657
											if buffer[position] != rune('I') {
												goto l651
											}
											position++
										}
									l657:
										{
											position659, tokenIndex659, depth659 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l660
											}
											position++
											goto l659
										l660:
											position, tokenIndex, depth = position659, tokenIndex659, depth659
											if buffer[position] != rune('R') {
												goto l651
											}
											position++
										}
									l659:
										{
											position661, tokenIndex661, depth661 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l662
											}
											position++
											goto l661
										l662:
											position, tokenIndex, depth = position661, tokenIndex661, depth661
											if buffer[position] != rune('I') {
												goto l651
											}
											position++
										}
									l661:
										if !rules[ruleskip]() {
											goto l651
										}
										depth--
										add(ruleISIRI, position652)
									}
									goto l498
								l651:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position664 := position
										depth++
										{
											position665, tokenIndex665, depth665 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l666
											}
											position++
											goto l665
										l666:
											position, tokenIndex, depth = position665, tokenIndex665, depth665
											if buffer[position] != rune('I') {
												goto l663
											}
											position++
										}
									l665:
										{
											position667, tokenIndex667, depth667 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l668
											}
											position++
											goto l667
										l668:
											position, tokenIndex, depth = position667, tokenIndex667, depth667
											if buffer[position] != rune('S') {
												goto l663
											}
											position++
										}
									l667:
										{
											position669, tokenIndex669, depth669 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l670
											}
											position++
											goto l669
										l670:
											position, tokenIndex, depth = position669, tokenIndex669, depth669
											if buffer[position] != rune('U') {
												goto l663
											}
											position++
										}
									l669:
										{
											position671, tokenIndex671, depth671 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l672
											}
											position++
											goto l671
										l672:
											position, tokenIndex, depth = position671, tokenIndex671, depth671
											if buffer[position] != rune('R') {
												goto l663
											}
											position++
										}
									l671:
										{
											position673, tokenIndex673, depth673 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l674
											}
											position++
											goto l673
										l674:
											position, tokenIndex, depth = position673, tokenIndex673, depth673
											if buffer[position] != rune('I') {
												goto l663
											}
											position++
										}
									l673:
										if !rules[ruleskip]() {
											goto l663
										}
										depth--
										add(ruleISURI, position664)
									}
									goto l498
								l663:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position676 := position
										depth++
										{
											position677, tokenIndex677, depth677 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l678
											}
											position++
											goto l677
										l678:
											position, tokenIndex, depth = position677, tokenIndex677, depth677
											if buffer[position] != rune('I') {
												goto l675
											}
											position++
										}
									l677:
										{
											position679, tokenIndex679, depth679 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l680
											}
											position++
											goto l679
										l680:
											position, tokenIndex, depth = position679, tokenIndex679, depth679
											if buffer[position] != rune('S') {
												goto l675
											}
											position++
										}
									l679:
										{
											position681, tokenIndex681, depth681 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l682
											}
											position++
											goto l681
										l682:
											position, tokenIndex, depth = position681, tokenIndex681, depth681
											if buffer[position] != rune('B') {
												goto l675
											}
											position++
										}
									l681:
										{
											position683, tokenIndex683, depth683 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l684
											}
											position++
											goto l683
										l684:
											position, tokenIndex, depth = position683, tokenIndex683, depth683
											if buffer[position] != rune('L') {
												goto l675
											}
											position++
										}
									l683:
										{
											position685, tokenIndex685, depth685 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l686
											}
											position++
											goto l685
										l686:
											position, tokenIndex, depth = position685, tokenIndex685, depth685
											if buffer[position] != rune('A') {
												goto l675
											}
											position++
										}
									l685:
										{
											position687, tokenIndex687, depth687 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l688
											}
											position++
											goto l687
										l688:
											position, tokenIndex, depth = position687, tokenIndex687, depth687
											if buffer[position] != rune('N') {
												goto l675
											}
											position++
										}
									l687:
										{
											position689, tokenIndex689, depth689 := position, tokenIndex, depth
											if buffer[position] != rune('k') {
												goto l690
											}
											position++
											goto l689
										l690:
											position, tokenIndex, depth = position689, tokenIndex689, depth689
											if buffer[position] != rune('K') {
												goto l675
											}
											position++
										}
									l689:
										if !rules[ruleskip]() {
											goto l675
										}
										depth--
										add(ruleISBLANK, position676)
									}
									goto l498
								l675:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										position692 := position
										depth++
										{
											position693, tokenIndex693, depth693 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l694
											}
											position++
											goto l693
										l694:
											position, tokenIndex, depth = position693, tokenIndex693, depth693
											if buffer[position] != rune('I') {
												goto l691
											}
											position++
										}
									l693:
										{
											position695, tokenIndex695, depth695 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l696
											}
											position++
											goto l695
										l696:
											position, tokenIndex, depth = position695, tokenIndex695, depth695
											if buffer[position] != rune('S') {
												goto l691
											}
											position++
										}
									l695:
										{
											position697, tokenIndex697, depth697 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l698
											}
											position++
											goto l697
										l698:
											position, tokenIndex, depth = position697, tokenIndex697, depth697
											if buffer[position] != rune('L') {
												goto l691
											}
											position++
										}
									l697:
										{
											position699, tokenIndex699, depth699 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l700
											}
											position++
											goto l699
										l700:
											position, tokenIndex, depth = position699, tokenIndex699, depth699
											if buffer[position] != rune('I') {
												goto l691
											}
											position++
										}
									l699:
										{
											position701, tokenIndex701, depth701 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l702
											}
											position++
											goto l701
										l702:
											position, tokenIndex, depth = position701, tokenIndex701, depth701
											if buffer[position] != rune('T') {
												goto l691
											}
											position++
										}
									l701:
										{
											position703, tokenIndex703, depth703 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l704
											}
											position++
											goto l703
										l704:
											position, tokenIndex, depth = position703, tokenIndex703, depth703
											if buffer[position] != rune('E') {
												goto l691
											}
											position++
										}
									l703:
										{
											position705, tokenIndex705, depth705 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l706
											}
											position++
											goto l705
										l706:
											position, tokenIndex, depth = position705, tokenIndex705, depth705
											if buffer[position] != rune('R') {
												goto l691
											}
											position++
										}
									l705:
										{
											position707, tokenIndex707, depth707 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l708
											}
											position++
											goto l707
										l708:
											position, tokenIndex, depth = position707, tokenIndex707, depth707
											if buffer[position] != rune('A') {
												goto l691
											}
											position++
										}
									l707:
										{
											position709, tokenIndex709, depth709 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l710
											}
											position++
											goto l709
										l710:
											position, tokenIndex, depth = position709, tokenIndex709, depth709
											if buffer[position] != rune('L') {
												goto l691
											}
											position++
										}
									l709:
										if !rules[ruleskip]() {
											goto l691
										}
										depth--
										add(ruleISLITERAL, position692)
									}
									goto l498
								l691:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
									{
										switch buffer[position] {
										case 'I', 'i':
											{
												position712 := position
												depth++
												{
													position713, tokenIndex713, depth713 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l714
													}
													position++
													goto l713
												l714:
													position, tokenIndex, depth = position713, tokenIndex713, depth713
													if buffer[position] != rune('I') {
														goto l497
													}
													position++
												}
											l713:
												{
													position715, tokenIndex715, depth715 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l716
													}
													position++
													goto l715
												l716:
													position, tokenIndex, depth = position715, tokenIndex715, depth715
													if buffer[position] != rune('S') {
														goto l497
													}
													position++
												}
											l715:
												{
													position717, tokenIndex717, depth717 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l718
													}
													position++
													goto l717
												l718:
													position, tokenIndex, depth = position717, tokenIndex717, depth717
													if buffer[position] != rune('N') {
														goto l497
													}
													position++
												}
											l717:
												{
													position719, tokenIndex719, depth719 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l720
													}
													position++
													goto l719
												l720:
													position, tokenIndex, depth = position719, tokenIndex719, depth719
													if buffer[position] != rune('U') {
														goto l497
													}
													position++
												}
											l719:
												{
													position721, tokenIndex721, depth721 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l722
													}
													position++
													goto l721
												l722:
													position, tokenIndex, depth = position721, tokenIndex721, depth721
													if buffer[position] != rune('M') {
														goto l497
													}
													position++
												}
											l721:
												{
													position723, tokenIndex723, depth723 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l724
													}
													position++
													goto l723
												l724:
													position, tokenIndex, depth = position723, tokenIndex723, depth723
													if buffer[position] != rune('E') {
														goto l497
													}
													position++
												}
											l723:
												{
													position725, tokenIndex725, depth725 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l726
													}
													position++
													goto l725
												l726:
													position, tokenIndex, depth = position725, tokenIndex725, depth725
													if buffer[position] != rune('R') {
														goto l497
													}
													position++
												}
											l725:
												{
													position727, tokenIndex727, depth727 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l728
													}
													position++
													goto l727
												l728:
													position, tokenIndex, depth = position727, tokenIndex727, depth727
													if buffer[position] != rune('I') {
														goto l497
													}
													position++
												}
											l727:
												{
													position729, tokenIndex729, depth729 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l730
													}
													position++
													goto l729
												l730:
													position, tokenIndex, depth = position729, tokenIndex729, depth729
													if buffer[position] != rune('C') {
														goto l497
													}
													position++
												}
											l729:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleISNUMERIC, position712)
											}
											break
										case 'S', 's':
											{
												position731 := position
												depth++
												{
													position732, tokenIndex732, depth732 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l733
													}
													position++
													goto l732
												l733:
													position, tokenIndex, depth = position732, tokenIndex732, depth732
													if buffer[position] != rune('S') {
														goto l497
													}
													position++
												}
											l732:
												{
													position734, tokenIndex734, depth734 := position, tokenIndex, depth
													if buffer[position] != rune('h') {
														goto l735
													}
													position++
													goto l734
												l735:
													position, tokenIndex, depth = position734, tokenIndex734, depth734
													if buffer[position] != rune('H') {
														goto l497
													}
													position++
												}
											l734:
												{
													position736, tokenIndex736, depth736 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l737
													}
													position++
													goto l736
												l737:
													position, tokenIndex, depth = position736, tokenIndex736, depth736
													if buffer[position] != rune('A') {
														goto l497
													}
													position++
												}
											l736:
												if buffer[position] != rune('5') {
													goto l497
												}
												position++
												if buffer[position] != rune('1') {
													goto l497
												}
												position++
												if buffer[position] != rune('2') {
													goto l497
												}
												position++
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleSHA512, position731)
											}
											break
										case 'M', 'm':
											{
												position738 := position
												depth++
												{
													position739, tokenIndex739, depth739 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l740
													}
													position++
													goto l739
												l740:
													position, tokenIndex, depth = position739, tokenIndex739, depth739
													if buffer[position] != rune('M') {
														goto l497
													}
													position++
												}
											l739:
												{
													position741, tokenIndex741, depth741 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l742
													}
													position++
													goto l741
												l742:
													position, tokenIndex, depth = position741, tokenIndex741, depth741
													if buffer[position] != rune('D') {
														goto l497
													}
													position++
												}
											l741:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleMD5, position738)
											}
											break
										case 'T', 't':
											{
												position743 := position
												depth++
												{
													position744, tokenIndex744, depth744 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l745
													}
													position++
													goto l744
												l745:
													position, tokenIndex, depth = position744, tokenIndex744, depth744
													if buffer[position] != rune('T') {
														goto l497
													}
													position++
												}
											l744:
												{
													position746, tokenIndex746, depth746 := position, tokenIndex, depth
													if buffer[position] != rune('z') {
														goto l747
													}
													position++
													goto l746
												l747:
													position, tokenIndex, depth = position746, tokenIndex746, depth746
													if buffer[position] != rune('Z') {
														goto l497
													}
													position++
												}
											l746:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleTZ, position743)
											}
											break
										case 'H', 'h':
											{
												position748 := position
												depth++
												{
													position749, tokenIndex749, depth749 := position, tokenIndex, depth
													if buffer[position] != rune('h') {
														goto l750
													}
													position++
													goto l749
												l750:
													position, tokenIndex, depth = position749, tokenIndex749, depth749
													if buffer[position] != rune('H') {
														goto l497
													}
													position++
												}
											l749:
												{
													position751, tokenIndex751, depth751 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l752
													}
													position++
													goto l751
												l752:
													position, tokenIndex, depth = position751, tokenIndex751, depth751
													if buffer[position] != rune('O') {
														goto l497
													}
													position++
												}
											l751:
												{
													position753, tokenIndex753, depth753 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l754
													}
													position++
													goto l753
												l754:
													position, tokenIndex, depth = position753, tokenIndex753, depth753
													if buffer[position] != rune('U') {
														goto l497
													}
													position++
												}
											l753:
												{
													position755, tokenIndex755, depth755 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l756
													}
													position++
													goto l755
												l756:
													position, tokenIndex, depth = position755, tokenIndex755, depth755
													if buffer[position] != rune('R') {
														goto l497
													}
													position++
												}
											l755:
												{
													position757, tokenIndex757, depth757 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l758
													}
													position++
													goto l757
												l758:
													position, tokenIndex, depth = position757, tokenIndex757, depth757
													if buffer[position] != rune('S') {
														goto l497
													}
													position++
												}
											l757:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleHOURS, position748)
											}
											break
										case 'D', 'd':
											{
												position759 := position
												depth++
												{
													position760, tokenIndex760, depth760 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l761
													}
													position++
													goto l760
												l761:
													position, tokenIndex, depth = position760, tokenIndex760, depth760
													if buffer[position] != rune('D') {
														goto l497
													}
													position++
												}
											l760:
												{
													position762, tokenIndex762, depth762 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l763
													}
													position++
													goto l762
												l763:
													position, tokenIndex, depth = position762, tokenIndex762, depth762
													if buffer[position] != rune('A') {
														goto l497
													}
													position++
												}
											l762:
												{
													position764, tokenIndex764, depth764 := position, tokenIndex, depth
													if buffer[position] != rune('y') {
														goto l765
													}
													position++
													goto l764
												l765:
													position, tokenIndex, depth = position764, tokenIndex764, depth764
													if buffer[position] != rune('Y') {
														goto l497
													}
													position++
												}
											l764:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleDAY, position759)
											}
											break
										case 'Y', 'y':
											{
												position766 := position
												depth++
												{
													position767, tokenIndex767, depth767 := position, tokenIndex, depth
													if buffer[position] != rune('y') {
														goto l768
													}
													position++
													goto l767
												l768:
													position, tokenIndex, depth = position767, tokenIndex767, depth767
													if buffer[position] != rune('Y') {
														goto l497
													}
													position++
												}
											l767:
												{
													position769, tokenIndex769, depth769 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l770
													}
													position++
													goto l769
												l770:
													position, tokenIndex, depth = position769, tokenIndex769, depth769
													if buffer[position] != rune('E') {
														goto l497
													}
													position++
												}
											l769:
												{
													position771, tokenIndex771, depth771 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l772
													}
													position++
													goto l771
												l772:
													position, tokenIndex, depth = position771, tokenIndex771, depth771
													if buffer[position] != rune('A') {
														goto l497
													}
													position++
												}
											l771:
												{
													position773, tokenIndex773, depth773 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l774
													}
													position++
													goto l773
												l774:
													position, tokenIndex, depth = position773, tokenIndex773, depth773
													if buffer[position] != rune('R') {
														goto l497
													}
													position++
												}
											l773:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleYEAR, position766)
											}
											break
										case 'E', 'e':
											{
												position775 := position
												depth++
												{
													position776, tokenIndex776, depth776 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l777
													}
													position++
													goto l776
												l777:
													position, tokenIndex, depth = position776, tokenIndex776, depth776
													if buffer[position] != rune('E') {
														goto l497
													}
													position++
												}
											l776:
												{
													position778, tokenIndex778, depth778 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l779
													}
													position++
													goto l778
												l779:
													position, tokenIndex, depth = position778, tokenIndex778, depth778
													if buffer[position] != rune('N') {
														goto l497
													}
													position++
												}
											l778:
												{
													position780, tokenIndex780, depth780 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l781
													}
													position++
													goto l780
												l781:
													position, tokenIndex, depth = position780, tokenIndex780, depth780
													if buffer[position] != rune('C') {
														goto l497
													}
													position++
												}
											l780:
												{
													position782, tokenIndex782, depth782 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l783
													}
													position++
													goto l782
												l783:
													position, tokenIndex, depth = position782, tokenIndex782, depth782
													if buffer[position] != rune('O') {
														goto l497
													}
													position++
												}
											l782:
												{
													position784, tokenIndex784, depth784 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l785
													}
													position++
													goto l784
												l785:
													position, tokenIndex, depth = position784, tokenIndex784, depth784
													if buffer[position] != rune('D') {
														goto l497
													}
													position++
												}
											l784:
												{
													position786, tokenIndex786, depth786 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l787
													}
													position++
													goto l786
												l787:
													position, tokenIndex, depth = position786, tokenIndex786, depth786
													if buffer[position] != rune('E') {
														goto l497
													}
													position++
												}
											l786:
												if buffer[position] != rune('_') {
													goto l497
												}
												position++
												{
													position788, tokenIndex788, depth788 := position, tokenIndex, depth
													if buffer[position] != rune('f') {
														goto l789
													}
													position++
													goto l788
												l789:
													position, tokenIndex, depth = position788, tokenIndex788, depth788
													if buffer[position] != rune('F') {
														goto l497
													}
													position++
												}
											l788:
												{
													position790, tokenIndex790, depth790 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l791
													}
													position++
													goto l790
												l791:
													position, tokenIndex, depth = position790, tokenIndex790, depth790
													if buffer[position] != rune('O') {
														goto l497
													}
													position++
												}
											l790:
												{
													position792, tokenIndex792, depth792 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l793
													}
													position++
													goto l792
												l793:
													position, tokenIndex, depth = position792, tokenIndex792, depth792
													if buffer[position] != rune('R') {
														goto l497
													}
													position++
												}
											l792:
												if buffer[position] != rune('_') {
													goto l497
												}
												position++
												{
													position794, tokenIndex794, depth794 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l795
													}
													position++
													goto l794
												l795:
													position, tokenIndex, depth = position794, tokenIndex794, depth794
													if buffer[position] != rune('U') {
														goto l497
													}
													position++
												}
											l794:
												{
													position796, tokenIndex796, depth796 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l797
													}
													position++
													goto l796
												l797:
													position, tokenIndex, depth = position796, tokenIndex796, depth796
													if buffer[position] != rune('R') {
														goto l497
													}
													position++
												}
											l796:
												{
													position798, tokenIndex798, depth798 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l799
													}
													position++
													goto l798
												l799:
													position, tokenIndex, depth = position798, tokenIndex798, depth798
													if buffer[position] != rune('I') {
														goto l497
													}
													position++
												}
											l798:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleENCODEFORURI, position775)
											}
											break
										case 'L', 'l':
											{
												position800 := position
												depth++
												{
													position801, tokenIndex801, depth801 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l802
													}
													position++
													goto l801
												l802:
													position, tokenIndex, depth = position801, tokenIndex801, depth801
													if buffer[position] != rune('L') {
														goto l497
													}
													position++
												}
											l801:
												{
													position803, tokenIndex803, depth803 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l804
													}
													position++
													goto l803
												l804:
													position, tokenIndex, depth = position803, tokenIndex803, depth803
													if buffer[position] != rune('C') {
														goto l497
													}
													position++
												}
											l803:
												{
													position805, tokenIndex805, depth805 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l806
													}
													position++
													goto l805
												l806:
													position, tokenIndex, depth = position805, tokenIndex805, depth805
													if buffer[position] != rune('A') {
														goto l497
													}
													position++
												}
											l805:
												{
													position807, tokenIndex807, depth807 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l808
													}
													position++
													goto l807
												l808:
													position, tokenIndex, depth = position807, tokenIndex807, depth807
													if buffer[position] != rune('S') {
														goto l497
													}
													position++
												}
											l807:
												{
													position809, tokenIndex809, depth809 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l810
													}
													position++
													goto l809
												l810:
													position, tokenIndex, depth = position809, tokenIndex809, depth809
													if buffer[position] != rune('E') {
														goto l497
													}
													position++
												}
											l809:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleLCASE, position800)
											}
											break
										case 'U', 'u':
											{
												position811 := position
												depth++
												{
													position812, tokenIndex812, depth812 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l813
													}
													position++
													goto l812
												l813:
													position, tokenIndex, depth = position812, tokenIndex812, depth812
													if buffer[position] != rune('U') {
														goto l497
													}
													position++
												}
											l812:
												{
													position814, tokenIndex814, depth814 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l815
													}
													position++
													goto l814
												l815:
													position, tokenIndex, depth = position814, tokenIndex814, depth814
													if buffer[position] != rune('C') {
														goto l497
													}
													position++
												}
											l814:
												{
													position816, tokenIndex816, depth816 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l817
													}
													position++
													goto l816
												l817:
													position, tokenIndex, depth = position816, tokenIndex816, depth816
													if buffer[position] != rune('A') {
														goto l497
													}
													position++
												}
											l816:
												{
													position818, tokenIndex818, depth818 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l819
													}
													position++
													goto l818
												l819:
													position, tokenIndex, depth = position818, tokenIndex818, depth818
													if buffer[position] != rune('S') {
														goto l497
													}
													position++
												}
											l818:
												{
													position820, tokenIndex820, depth820 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l821
													}
													position++
													goto l820
												l821:
													position, tokenIndex, depth = position820, tokenIndex820, depth820
													if buffer[position] != rune('E') {
														goto l497
													}
													position++
												}
											l820:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleUCASE, position811)
											}
											break
										case 'F', 'f':
											{
												position822 := position
												depth++
												{
													position823, tokenIndex823, depth823 := position, tokenIndex, depth
													if buffer[position] != rune('f') {
														goto l824
													}
													position++
													goto l823
												l824:
													position, tokenIndex, depth = position823, tokenIndex823, depth823
													if buffer[position] != rune('F') {
														goto l497
													}
													position++
												}
											l823:
												{
													position825, tokenIndex825, depth825 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l826
													}
													position++
													goto l825
												l826:
													position, tokenIndex, depth = position825, tokenIndex825, depth825
													if buffer[position] != rune('L') {
														goto l497
													}
													position++
												}
											l825:
												{
													position827, tokenIndex827, depth827 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l828
													}
													position++
													goto l827
												l828:
													position, tokenIndex, depth = position827, tokenIndex827, depth827
													if buffer[position] != rune('O') {
														goto l497
													}
													position++
												}
											l827:
												{
													position829, tokenIndex829, depth829 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l830
													}
													position++
													goto l829
												l830:
													position, tokenIndex, depth = position829, tokenIndex829, depth829
													if buffer[position] != rune('O') {
														goto l497
													}
													position++
												}
											l829:
												{
													position831, tokenIndex831, depth831 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l832
													}
													position++
													goto l831
												l832:
													position, tokenIndex, depth = position831, tokenIndex831, depth831
													if buffer[position] != rune('R') {
														goto l497
													}
													position++
												}
											l831:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleFLOOR, position822)
											}
											break
										case 'R', 'r':
											{
												position833 := position
												depth++
												{
													position834, tokenIndex834, depth834 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l835
													}
													position++
													goto l834
												l835:
													position, tokenIndex, depth = position834, tokenIndex834, depth834
													if buffer[position] != rune('R') {
														goto l497
													}
													position++
												}
											l834:
												{
													position836, tokenIndex836, depth836 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l837
													}
													position++
													goto l836
												l837:
													position, tokenIndex, depth = position836, tokenIndex836, depth836
													if buffer[position] != rune('O') {
														goto l497
													}
													position++
												}
											l836:
												{
													position838, tokenIndex838, depth838 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l839
													}
													position++
													goto l838
												l839:
													position, tokenIndex, depth = position838, tokenIndex838, depth838
													if buffer[position] != rune('U') {
														goto l497
													}
													position++
												}
											l838:
												{
													position840, tokenIndex840, depth840 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l841
													}
													position++
													goto l840
												l841:
													position, tokenIndex, depth = position840, tokenIndex840, depth840
													if buffer[position] != rune('N') {
														goto l497
													}
													position++
												}
											l840:
												{
													position842, tokenIndex842, depth842 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l843
													}
													position++
													goto l842
												l843:
													position, tokenIndex, depth = position842, tokenIndex842, depth842
													if buffer[position] != rune('D') {
														goto l497
													}
													position++
												}
											l842:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleROUND, position833)
											}
											break
										case 'C', 'c':
											{
												position844 := position
												depth++
												{
													position845, tokenIndex845, depth845 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l846
													}
													position++
													goto l845
												l846:
													position, tokenIndex, depth = position845, tokenIndex845, depth845
													if buffer[position] != rune('C') {
														goto l497
													}
													position++
												}
											l845:
												{
													position847, tokenIndex847, depth847 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l848
													}
													position++
													goto l847
												l848:
													position, tokenIndex, depth = position847, tokenIndex847, depth847
													if buffer[position] != rune('E') {
														goto l497
													}
													position++
												}
											l847:
												{
													position849, tokenIndex849, depth849 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l850
													}
													position++
													goto l849
												l850:
													position, tokenIndex, depth = position849, tokenIndex849, depth849
													if buffer[position] != rune('I') {
														goto l497
													}
													position++
												}
											l849:
												{
													position851, tokenIndex851, depth851 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l852
													}
													position++
													goto l851
												l852:
													position, tokenIndex, depth = position851, tokenIndex851, depth851
													if buffer[position] != rune('L') {
														goto l497
													}
													position++
												}
											l851:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleCEIL, position844)
											}
											break
										default:
											{
												position853 := position
												depth++
												{
													position854, tokenIndex854, depth854 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l855
													}
													position++
													goto l854
												l855:
													position, tokenIndex, depth = position854, tokenIndex854, depth854
													if buffer[position] != rune('A') {
														goto l497
													}
													position++
												}
											l854:
												{
													position856, tokenIndex856, depth856 := position, tokenIndex, depth
													if buffer[position] != rune('b') {
														goto l857
													}
													position++
													goto l856
												l857:
													position, tokenIndex, depth = position856, tokenIndex856, depth856
													if buffer[position] != rune('B') {
														goto l497
													}
													position++
												}
											l856:
												{
													position858, tokenIndex858, depth858 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l859
													}
													position++
													goto l858
												l859:
													position, tokenIndex, depth = position858, tokenIndex858, depth858
													if buffer[position] != rune('S') {
														goto l497
													}
													position++
												}
											l858:
												if !rules[ruleskip]() {
													goto l497
												}
												depth--
												add(ruleABS, position853)
											}
											break
										}
									}

								}
							l498:
								if !rules[ruleLPAREN]() {
									goto l497
								}
								if !rules[ruleexpression]() {
									goto l497
								}
								if !rules[ruleRPAREN]() {
									goto l497
								}
								goto l496
							l497:
								position, tokenIndex, depth = position496, tokenIndex496, depth496
								{
									position861, tokenIndex861, depth861 := position, tokenIndex, depth
									{
										position863 := position
										depth++
										{
											position864, tokenIndex864, depth864 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l865
											}
											position++
											goto l864
										l865:
											position, tokenIndex, depth = position864, tokenIndex864, depth864
											if buffer[position] != rune('S') {
												goto l862
											}
											position++
										}
									l864:
										{
											position866, tokenIndex866, depth866 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l867
											}
											position++
											goto l866
										l867:
											position, tokenIndex, depth = position866, tokenIndex866, depth866
											if buffer[position] != rune('T') {
												goto l862
											}
											position++
										}
									l866:
										{
											position868, tokenIndex868, depth868 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l869
											}
											position++
											goto l868
										l869:
											position, tokenIndex, depth = position868, tokenIndex868, depth868
											if buffer[position] != rune('R') {
												goto l862
											}
											position++
										}
									l868:
										{
											position870, tokenIndex870, depth870 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l871
											}
											position++
											goto l870
										l871:
											position, tokenIndex, depth = position870, tokenIndex870, depth870
											if buffer[position] != rune('S') {
												goto l862
											}
											position++
										}
									l870:
										{
											position872, tokenIndex872, depth872 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l873
											}
											position++
											goto l872
										l873:
											position, tokenIndex, depth = position872, tokenIndex872, depth872
											if buffer[position] != rune('T') {
												goto l862
											}
											position++
										}
									l872:
										{
											position874, tokenIndex874, depth874 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l875
											}
											position++
											goto l874
										l875:
											position, tokenIndex, depth = position874, tokenIndex874, depth874
											if buffer[position] != rune('A') {
												goto l862
											}
											position++
										}
									l874:
										{
											position876, tokenIndex876, depth876 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l877
											}
											position++
											goto l876
										l877:
											position, tokenIndex, depth = position876, tokenIndex876, depth876
											if buffer[position] != rune('R') {
												goto l862
											}
											position++
										}
									l876:
										{
											position878, tokenIndex878, depth878 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l879
											}
											position++
											goto l878
										l879:
											position, tokenIndex, depth = position878, tokenIndex878, depth878
											if buffer[position] != rune('T') {
												goto l862
											}
											position++
										}
									l878:
										{
											position880, tokenIndex880, depth880 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l881
											}
											position++
											goto l880
										l881:
											position, tokenIndex, depth = position880, tokenIndex880, depth880
											if buffer[position] != rune('S') {
												goto l862
											}
											position++
										}
									l880:
										if !rules[ruleskip]() {
											goto l862
										}
										depth--
										add(ruleSTRSTARTS, position863)
									}
									goto l861
								l862:
									position, tokenIndex, depth = position861, tokenIndex861, depth861
									{
										position883 := position
										depth++
										{
											position884, tokenIndex884, depth884 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l885
											}
											position++
											goto l884
										l885:
											position, tokenIndex, depth = position884, tokenIndex884, depth884
											if buffer[position] != rune('S') {
												goto l882
											}
											position++
										}
									l884:
										{
											position886, tokenIndex886, depth886 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l887
											}
											position++
											goto l886
										l887:
											position, tokenIndex, depth = position886, tokenIndex886, depth886
											if buffer[position] != rune('T') {
												goto l882
											}
											position++
										}
									l886:
										{
											position888, tokenIndex888, depth888 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l889
											}
											position++
											goto l888
										l889:
											position, tokenIndex, depth = position888, tokenIndex888, depth888
											if buffer[position] != rune('R') {
												goto l882
											}
											position++
										}
									l888:
										{
											position890, tokenIndex890, depth890 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l891
											}
											position++
											goto l890
										l891:
											position, tokenIndex, depth = position890, tokenIndex890, depth890
											if buffer[position] != rune('E') {
												goto l882
											}
											position++
										}
									l890:
										{
											position892, tokenIndex892, depth892 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l893
											}
											position++
											goto l892
										l893:
											position, tokenIndex, depth = position892, tokenIndex892, depth892
											if buffer[position] != rune('N') {
												goto l882
											}
											position++
										}
									l892:
										{
											position894, tokenIndex894, depth894 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l895
											}
											position++
											goto l894
										l895:
											position, tokenIndex, depth = position894, tokenIndex894, depth894
											if buffer[position] != rune('D') {
												goto l882
											}
											position++
										}
									l894:
										{
											position896, tokenIndex896, depth896 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l897
											}
											position++
											goto l896
										l897:
											position, tokenIndex, depth = position896, tokenIndex896, depth896
											if buffer[position] != rune('S') {
												goto l882
											}
											position++
										}
									l896:
										if !rules[ruleskip]() {
											goto l882
										}
										depth--
										add(ruleSTRENDS, position883)
									}
									goto l861
								l882:
									position, tokenIndex, depth = position861, tokenIndex861, depth861
									{
										position899 := position
										depth++
										{
											position900, tokenIndex900, depth900 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l901
											}
											position++
											goto l900
										l901:
											position, tokenIndex, depth = position900, tokenIndex900, depth900
											if buffer[position] != rune('S') {
												goto l898
											}
											position++
										}
									l900:
										{
											position902, tokenIndex902, depth902 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l903
											}
											position++
											goto l902
										l903:
											position, tokenIndex, depth = position902, tokenIndex902, depth902
											if buffer[position] != rune('T') {
												goto l898
											}
											position++
										}
									l902:
										{
											position904, tokenIndex904, depth904 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l905
											}
											position++
											goto l904
										l905:
											position, tokenIndex, depth = position904, tokenIndex904, depth904
											if buffer[position] != rune('R') {
												goto l898
											}
											position++
										}
									l904:
										{
											position906, tokenIndex906, depth906 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l907
											}
											position++
											goto l906
										l907:
											position, tokenIndex, depth = position906, tokenIndex906, depth906
											if buffer[position] != rune('B') {
												goto l898
											}
											position++
										}
									l906:
										{
											position908, tokenIndex908, depth908 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l909
											}
											position++
											goto l908
										l909:
											position, tokenIndex, depth = position908, tokenIndex908, depth908
											if buffer[position] != rune('E') {
												goto l898
											}
											position++
										}
									l908:
										{
											position910, tokenIndex910, depth910 := position, tokenIndex, depth
											if buffer[position] != rune('f') {
												goto l911
											}
											position++
											goto l910
										l911:
											position, tokenIndex, depth = position910, tokenIndex910, depth910
											if buffer[position] != rune('F') {
												goto l898
											}
											position++
										}
									l910:
										{
											position912, tokenIndex912, depth912 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l913
											}
											position++
											goto l912
										l913:
											position, tokenIndex, depth = position912, tokenIndex912, depth912
											if buffer[position] != rune('O') {
												goto l898
											}
											position++
										}
									l912:
										{
											position914, tokenIndex914, depth914 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l915
											}
											position++
											goto l914
										l915:
											position, tokenIndex, depth = position914, tokenIndex914, depth914
											if buffer[position] != rune('R') {
												goto l898
											}
											position++
										}
									l914:
										{
											position916, tokenIndex916, depth916 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l917
											}
											position++
											goto l916
										l917:
											position, tokenIndex, depth = position916, tokenIndex916, depth916
											if buffer[position] != rune('E') {
												goto l898
											}
											position++
										}
									l916:
										if !rules[ruleskip]() {
											goto l898
										}
										depth--
										add(ruleSTRBEFORE, position899)
									}
									goto l861
								l898:
									position, tokenIndex, depth = position861, tokenIndex861, depth861
									{
										position919 := position
										depth++
										{
											position920, tokenIndex920, depth920 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l921
											}
											position++
											goto l920
										l921:
											position, tokenIndex, depth = position920, tokenIndex920, depth920
											if buffer[position] != rune('S') {
												goto l918
											}
											position++
										}
									l920:
										{
											position922, tokenIndex922, depth922 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l923
											}
											position++
											goto l922
										l923:
											position, tokenIndex, depth = position922, tokenIndex922, depth922
											if buffer[position] != rune('T') {
												goto l918
											}
											position++
										}
									l922:
										{
											position924, tokenIndex924, depth924 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l925
											}
											position++
											goto l924
										l925:
											position, tokenIndex, depth = position924, tokenIndex924, depth924
											if buffer[position] != rune('R') {
												goto l918
											}
											position++
										}
									l924:
										{
											position926, tokenIndex926, depth926 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l927
											}
											position++
											goto l926
										l927:
											position, tokenIndex, depth = position926, tokenIndex926, depth926
											if buffer[position] != rune('A') {
												goto l918
											}
											position++
										}
									l926:
										{
											position928, tokenIndex928, depth928 := position, tokenIndex, depth
											if buffer[position] != rune('f') {
												goto l929
											}
											position++
											goto l928
										l929:
											position, tokenIndex, depth = position928, tokenIndex928, depth928
											if buffer[position] != rune('F') {
												goto l918
											}
											position++
										}
									l928:
										{
											position930, tokenIndex930, depth930 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l931
											}
											position++
											goto l930
										l931:
											position, tokenIndex, depth = position930, tokenIndex930, depth930
											if buffer[position] != rune('T') {
												goto l918
											}
											position++
										}
									l930:
										{
											position932, tokenIndex932, depth932 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l933
											}
											position++
											goto l932
										l933:
											position, tokenIndex, depth = position932, tokenIndex932, depth932
											if buffer[position] != rune('E') {
												goto l918
											}
											position++
										}
									l932:
										{
											position934, tokenIndex934, depth934 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l935
											}
											position++
											goto l934
										l935:
											position, tokenIndex, depth = position934, tokenIndex934, depth934
											if buffer[position] != rune('R') {
												goto l918
											}
											position++
										}
									l934:
										if !rules[ruleskip]() {
											goto l918
										}
										depth--
										add(ruleSTRAFTER, position919)
									}
									goto l861
								l918:
									position, tokenIndex, depth = position861, tokenIndex861, depth861
									{
										position937 := position
										depth++
										{
											position938, tokenIndex938, depth938 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l939
											}
											position++
											goto l938
										l939:
											position, tokenIndex, depth = position938, tokenIndex938, depth938
											if buffer[position] != rune('S') {
												goto l936
											}
											position++
										}
									l938:
										{
											position940, tokenIndex940, depth940 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l941
											}
											position++
											goto l940
										l941:
											position, tokenIndex, depth = position940, tokenIndex940, depth940
											if buffer[position] != rune('T') {
												goto l936
											}
											position++
										}
									l940:
										{
											position942, tokenIndex942, depth942 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l943
											}
											position++
											goto l942
										l943:
											position, tokenIndex, depth = position942, tokenIndex942, depth942
											if buffer[position] != rune('R') {
												goto l936
											}
											position++
										}
									l942:
										{
											position944, tokenIndex944, depth944 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l945
											}
											position++
											goto l944
										l945:
											position, tokenIndex, depth = position944, tokenIndex944, depth944
											if buffer[position] != rune('L') {
												goto l936
											}
											position++
										}
									l944:
										{
											position946, tokenIndex946, depth946 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l947
											}
											position++
											goto l946
										l947:
											position, tokenIndex, depth = position946, tokenIndex946, depth946
											if buffer[position] != rune('A') {
												goto l936
											}
											position++
										}
									l946:
										{
											position948, tokenIndex948, depth948 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l949
											}
											position++
											goto l948
										l949:
											position, tokenIndex, depth = position948, tokenIndex948, depth948
											if buffer[position] != rune('N') {
												goto l936
											}
											position++
										}
									l948:
										{
											position950, tokenIndex950, depth950 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l951
											}
											position++
											goto l950
										l951:
											position, tokenIndex, depth = position950, tokenIndex950, depth950
											if buffer[position] != rune('G') {
												goto l936
											}
											position++
										}
									l950:
										if !rules[ruleskip]() {
											goto l936
										}
										depth--
										add(ruleSTRLANG, position937)
									}
									goto l861
								l936:
									position, tokenIndex, depth = position861, tokenIndex861, depth861
									{
										position953 := position
										depth++
										{
											position954, tokenIndex954, depth954 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l955
											}
											position++
											goto l954
										l955:
											position, tokenIndex, depth = position954, tokenIndex954, depth954
											if buffer[position] != rune('S') {
												goto l952
											}
											position++
										}
									l954:
										{
											position956, tokenIndex956, depth956 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l957
											}
											position++
											goto l956
										l957:
											position, tokenIndex, depth = position956, tokenIndex956, depth956
											if buffer[position] != rune('T') {
												goto l952
											}
											position++
										}
									l956:
										{
											position958, tokenIndex958, depth958 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l959
											}
											position++
											goto l958
										l959:
											position, tokenIndex, depth = position958, tokenIndex958, depth958
											if buffer[position] != rune('R') {
												goto l952
											}
											position++
										}
									l958:
										{
											position960, tokenIndex960, depth960 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l961
											}
											position++
											goto l960
										l961:
											position, tokenIndex, depth = position960, tokenIndex960, depth960
											if buffer[position] != rune('D') {
												goto l952
											}
											position++
										}
									l960:
										{
											position962, tokenIndex962, depth962 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l963
											}
											position++
											goto l962
										l963:
											position, tokenIndex, depth = position962, tokenIndex962, depth962
											if buffer[position] != rune('T') {
												goto l952
											}
											position++
										}
									l962:
										if !rules[ruleskip]() {
											goto l952
										}
										depth--
										add(ruleSTRDT, position953)
									}
									goto l861
								l952:
									position, tokenIndex, depth = position861, tokenIndex861, depth861
									{
										switch buffer[position] {
										case 'S', 's':
											{
												position965 := position
												depth++
												{
													position966, tokenIndex966, depth966 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l967
													}
													position++
													goto l966
												l967:
													position, tokenIndex, depth = position966, tokenIndex966, depth966
													if buffer[position] != rune('S') {
														goto l860
													}
													position++
												}
											l966:
												{
													position968, tokenIndex968, depth968 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l969
													}
													position++
													goto l968
												l969:
													position, tokenIndex, depth = position968, tokenIndex968, depth968
													if buffer[position] != rune('A') {
														goto l860
													}
													position++
												}
											l968:
												{
													position970, tokenIndex970, depth970 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l971
													}
													position++
													goto l970
												l971:
													position, tokenIndex, depth = position970, tokenIndex970, depth970
													if buffer[position] != rune('M') {
														goto l860
													}
													position++
												}
											l970:
												{
													position972, tokenIndex972, depth972 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l973
													}
													position++
													goto l972
												l973:
													position, tokenIndex, depth = position972, tokenIndex972, depth972
													if buffer[position] != rune('E') {
														goto l860
													}
													position++
												}
											l972:
												{
													position974, tokenIndex974, depth974 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l975
													}
													position++
													goto l974
												l975:
													position, tokenIndex, depth = position974, tokenIndex974, depth974
													if buffer[position] != rune('T') {
														goto l860
													}
													position++
												}
											l974:
												{
													position976, tokenIndex976, depth976 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l977
													}
													position++
													goto l976
												l977:
													position, tokenIndex, depth = position976, tokenIndex976, depth976
													if buffer[position] != rune('E') {
														goto l860
													}
													position++
												}
											l976:
												{
													position978, tokenIndex978, depth978 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l979
													}
													position++
													goto l978
												l979:
													position, tokenIndex, depth = position978, tokenIndex978, depth978
													if buffer[position] != rune('R') {
														goto l860
													}
													position++
												}
											l978:
												{
													position980, tokenIndex980, depth980 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l981
													}
													position++
													goto l980
												l981:
													position, tokenIndex, depth = position980, tokenIndex980, depth980
													if buffer[position] != rune('M') {
														goto l860
													}
													position++
												}
											l980:
												if !rules[ruleskip]() {
													goto l860
												}
												depth--
												add(ruleSAMETERM, position965)
											}
											break
										case 'C', 'c':
											{
												position982 := position
												depth++
												{
													position983, tokenIndex983, depth983 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l984
													}
													position++
													goto l983
												l984:
													position, tokenIndex, depth = position983, tokenIndex983, depth983
													if buffer[position] != rune('C') {
														goto l860
													}
													position++
												}
											l983:
												{
													position985, tokenIndex985, depth985 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l986
													}
													position++
													goto l985
												l986:
													position, tokenIndex, depth = position985, tokenIndex985, depth985
													if buffer[position] != rune('O') {
														goto l860
													}
													position++
												}
											l985:
												{
													position987, tokenIndex987, depth987 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l988
													}
													position++
													goto l987
												l988:
													position, tokenIndex, depth = position987, tokenIndex987, depth987
													if buffer[position] != rune('N') {
														goto l860
													}
													position++
												}
											l987:
												{
													position989, tokenIndex989, depth989 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l990
													}
													position++
													goto l989
												l990:
													position, tokenIndex, depth = position989, tokenIndex989, depth989
													if buffer[position] != rune('T') {
														goto l860
													}
													position++
												}
											l989:
												{
													position991, tokenIndex991, depth991 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l992
													}
													position++
													goto l991
												l992:
													position, tokenIndex, depth = position991, tokenIndex991, depth991
													if buffer[position] != rune('A') {
														goto l860
													}
													position++
												}
											l991:
												{
													position993, tokenIndex993, depth993 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l994
													}
													position++
													goto l993
												l994:
													position, tokenIndex, depth = position993, tokenIndex993, depth993
													if buffer[position] != rune('I') {
														goto l860
													}
													position++
												}
											l993:
												{
													position995, tokenIndex995, depth995 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l996
													}
													position++
													goto l995
												l996:
													position, tokenIndex, depth = position995, tokenIndex995, depth995
													if buffer[position] != rune('N') {
														goto l860
													}
													position++
												}
											l995:
												{
													position997, tokenIndex997, depth997 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l998
													}
													position++
													goto l997
												l998:
													position, tokenIndex, depth = position997, tokenIndex997, depth997
													if buffer[position] != rune('S') {
														goto l860
													}
													position++
												}
											l997:
												if !rules[ruleskip]() {
													goto l860
												}
												depth--
												add(ruleCONTAINS, position982)
											}
											break
										default:
											{
												position999 := position
												depth++
												{
													position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l1001
													}
													position++
													goto l1000
												l1001:
													position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
													if buffer[position] != rune('L') {
														goto l860
													}
													position++
												}
											l1000:
												{
													position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1003
													}
													position++
													goto l1002
												l1003:
													position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
													if buffer[position] != rune('A') {
														goto l860
													}
													position++
												}
											l1002:
												{
													position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1005
													}
													position++
													goto l1004
												l1005:
													position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
													if buffer[position] != rune('N') {
														goto l860
													}
													position++
												}
											l1004:
												{
													position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
													if buffer[position] != rune('g') {
														goto l1007
													}
													position++
													goto l1006
												l1007:
													position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
													if buffer[position] != rune('G') {
														goto l860
													}
													position++
												}
											l1006:
												{
													position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l1009
													}
													position++
													goto l1008
												l1009:
													position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
													if buffer[position] != rune('M') {
														goto l860
													}
													position++
												}
											l1008:
												{
													position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1011
													}
													position++
													goto l1010
												l1011:
													position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
													if buffer[position] != rune('A') {
														goto l860
													}
													position++
												}
											l1010:
												{
													position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1013
													}
													position++
													goto l1012
												l1013:
													position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
													if buffer[position] != rune('T') {
														goto l860
													}
													position++
												}
											l1012:
												{
													position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1015
													}
													position++
													goto l1014
												l1015:
													position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
													if buffer[position] != rune('C') {
														goto l860
													}
													position++
												}
											l1014:
												{
													position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
													if buffer[position] != rune('h') {
														goto l1017
													}
													position++
													goto l1016
												l1017:
													position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
													if buffer[position] != rune('H') {
														goto l860
													}
													position++
												}
											l1016:
												{
													position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1019
													}
													position++
													goto l1018
												l1019:
													position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
													if buffer[position] != rune('E') {
														goto l860
													}
													position++
												}
											l1018:
												{
													position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1021
													}
													position++
													goto l1020
												l1021:
													position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
													if buffer[position] != rune('S') {
														goto l860
													}
													position++
												}
											l1020:
												if !rules[ruleskip]() {
													goto l860
												}
												depth--
												add(ruleLANGMATCHES, position999)
											}
											break
										}
									}

								}
							l861:
								if !rules[ruleLPAREN]() {
									goto l860
								}
								if !rules[ruleexpression]() {
									goto l860
								}
								if !rules[ruleCOMMA]() {
									goto l860
								}
								if !rules[ruleexpression]() {
									goto l860
								}
								if !rules[ruleRPAREN]() {
									goto l860
								}
								goto l496
							l860:
								position, tokenIndex, depth = position496, tokenIndex496, depth496
								{
									position1023 := position
									depth++
									{
										position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1025
										}
										position++
										goto l1024
									l1025:
										position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
										if buffer[position] != rune('B') {
											goto l1022
										}
										position++
									}
								l1024:
									{
										position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1027
										}
										position++
										goto l1026
									l1027:
										position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
										if buffer[position] != rune('O') {
											goto l1022
										}
										position++
									}
								l1026:
									{
										position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1029
										}
										position++
										goto l1028
									l1029:
										position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
										if buffer[position] != rune('U') {
											goto l1022
										}
										position++
									}
								l1028:
									{
										position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1031
										}
										position++
										goto l1030
									l1031:
										position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
										if buffer[position] != rune('N') {
											goto l1022
										}
										position++
									}
								l1030:
									{
										position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1033
										}
										position++
										goto l1032
									l1033:
										position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
										if buffer[position] != rune('D') {
											goto l1022
										}
										position++
									}
								l1032:
									if !rules[ruleskip]() {
										goto l1022
									}
									depth--
									add(ruleBOUND, position1023)
								}
								if !rules[ruleLPAREN]() {
									goto l1022
								}
								if !rules[rulevar]() {
									goto l1022
								}
								if !rules[ruleRPAREN]() {
									goto l1022
								}
								goto l496
							l1022:
								position, tokenIndex, depth = position496, tokenIndex496, depth496
								{
									switch buffer[position] {
									case 'S', 's':
										{
											position1036 := position
											depth++
											{
												position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l1038
												}
												position++
												goto l1037
											l1038:
												position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
												if buffer[position] != rune('S') {
													goto l1034
												}
												position++
											}
										l1037:
											{
												position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
												if buffer[position] != rune('t') {
													goto l1040
												}
												position++
												goto l1039
											l1040:
												position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
												if buffer[position] != rune('T') {
													goto l1034
												}
												position++
											}
										l1039:
											{
												position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
												if buffer[position] != rune('r') {
													goto l1042
												}
												position++
												goto l1041
											l1042:
												position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
												if buffer[position] != rune('R') {
													goto l1034
												}
												position++
											}
										l1041:
											{
												position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1044
												}
												position++
												goto l1043
											l1044:
												position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
												if buffer[position] != rune('U') {
													goto l1034
												}
												position++
											}
										l1043:
											{
												position1045, tokenIndex1045, depth1045 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1046
												}
												position++
												goto l1045
											l1046:
												position, tokenIndex, depth = position1045, tokenIndex1045, depth1045
												if buffer[position] != rune('U') {
													goto l1034
												}
												position++
											}
										l1045:
											{
												position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
												if buffer[position] != rune('i') {
													goto l1048
												}
												position++
												goto l1047
											l1048:
												position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
												if buffer[position] != rune('I') {
													goto l1034
												}
												position++
											}
										l1047:
											{
												position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1050
												}
												position++
												goto l1049
											l1050:
												position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
												if buffer[position] != rune('D') {
													goto l1034
												}
												position++
											}
										l1049:
											if !rules[ruleskip]() {
												goto l1034
											}
											depth--
											add(ruleSTRUUID, position1036)
										}
										break
									case 'U', 'u':
										{
											position1051 := position
											depth++
											{
												position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1053
												}
												position++
												goto l1052
											l1053:
												position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
												if buffer[position] != rune('U') {
													goto l1034
												}
												position++
											}
										l1052:
											{
												position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1055
												}
												position++
												goto l1054
											l1055:
												position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
												if buffer[position] != rune('U') {
													goto l1034
												}
												position++
											}
										l1054:
											{
												position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
												if buffer[position] != rune('i') {
													goto l1057
												}
												position++
												goto l1056
											l1057:
												position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
												if buffer[position] != rune('I') {
													goto l1034
												}
												position++
											}
										l1056:
											{
												position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1059
												}
												position++
												goto l1058
											l1059:
												position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
												if buffer[position] != rune('D') {
													goto l1034
												}
												position++
											}
										l1058:
											if !rules[ruleskip]() {
												goto l1034
											}
											depth--
											add(ruleUUID, position1051)
										}
										break
									case 'N', 'n':
										{
											position1060 := position
											depth++
											{
												position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
												if buffer[position] != rune('n') {
													goto l1062
												}
												position++
												goto l1061
											l1062:
												position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
												if buffer[position] != rune('N') {
													goto l1034
												}
												position++
											}
										l1061:
											{
												position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
												if buffer[position] != rune('o') {
													goto l1064
												}
												position++
												goto l1063
											l1064:
												position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
												if buffer[position] != rune('O') {
													goto l1034
												}
												position++
											}
										l1063:
											{
												position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
												if buffer[position] != rune('w') {
													goto l1066
												}
												position++
												goto l1065
											l1066:
												position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
												if buffer[position] != rune('W') {
													goto l1034
												}
												position++
											}
										l1065:
											if !rules[ruleskip]() {
												goto l1034
											}
											depth--
											add(ruleNOW, position1060)
										}
										break
									default:
										{
											position1067 := position
											depth++
											{
												position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
												if buffer[position] != rune('r') {
													goto l1069
												}
												position++
												goto l1068
											l1069:
												position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
												if buffer[position] != rune('R') {
													goto l1034
												}
												position++
											}
										l1068:
											{
												position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l1071
												}
												position++
												goto l1070
											l1071:
												position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
												if buffer[position] != rune('A') {
													goto l1034
												}
												position++
											}
										l1070:
											{
												position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
												if buffer[position] != rune('n') {
													goto l1073
												}
												position++
												goto l1072
											l1073:
												position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
												if buffer[position] != rune('N') {
													goto l1034
												}
												position++
											}
										l1072:
											{
												position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1075
												}
												position++
												goto l1074
											l1075:
												position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
												if buffer[position] != rune('D') {
													goto l1034
												}
												position++
											}
										l1074:
											if !rules[ruleskip]() {
												goto l1034
											}
											depth--
											add(ruleRAND, position1067)
										}
										break
									}
								}

								if !rules[rulenil]() {
									goto l1034
								}
								goto l496
							l1034:
								position, tokenIndex, depth = position496, tokenIndex496, depth496
								{
									switch buffer[position] {
									case 'E', 'N', 'e', 'n':
										{
											position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
											{
												position1079 := position
												depth++
												{
													position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1081
													}
													position++
													goto l1080
												l1081:
													position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
													if buffer[position] != rune('E') {
														goto l1078
													}
													position++
												}
											l1080:
												{
													position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
													if buffer[position] != rune('x') {
														goto l1083
													}
													position++
													goto l1082
												l1083:
													position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
													if buffer[position] != rune('X') {
														goto l1078
													}
													position++
												}
											l1082:
												{
													position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l1085
													}
													position++
													goto l1084
												l1085:
													position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
													if buffer[position] != rune('I') {
														goto l1078
													}
													position++
												}
											l1084:
												{
													position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1087
													}
													position++
													goto l1086
												l1087:
													position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
													if buffer[position] != rune('S') {
														goto l1078
													}
													position++
												}
											l1086:
												{
													position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1089
													}
													position++
													goto l1088
												l1089:
													position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
													if buffer[position] != rune('T') {
														goto l1078
													}
													position++
												}
											l1088:
												{
													position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1091
													}
													position++
													goto l1090
												l1091:
													position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
													if buffer[position] != rune('S') {
														goto l1078
													}
													position++
												}
											l1090:
												if !rules[ruleskip]() {
													goto l1078
												}
												depth--
												add(ruleEXISTS, position1079)
											}
											goto l1077
										l1078:
											position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
											{
												position1092 := position
												depth++
												{
													position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1094
													}
													position++
													goto l1093
												l1094:
													position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
													if buffer[position] != rune('N') {
														goto l494
													}
													position++
												}
											l1093:
												{
													position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l1096
													}
													position++
													goto l1095
												l1096:
													position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
													if buffer[position] != rune('O') {
														goto l494
													}
													position++
												}
											l1095:
												{
													position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1098
													}
													position++
													goto l1097
												l1098:
													position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
													if buffer[position] != rune('T') {
														goto l494
													}
													position++
												}
											l1097:
												if buffer[position] != rune(' ') {
													goto l494
												}
												position++
												{
													position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1100
													}
													position++
													goto l1099
												l1100:
													position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
													if buffer[position] != rune('E') {
														goto l494
													}
													position++
												}
											l1099:
												{
													position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
													if buffer[position] != rune('x') {
														goto l1102
													}
													position++
													goto l1101
												l1102:
													position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
													if buffer[position] != rune('X') {
														goto l494
													}
													position++
												}
											l1101:
												{
													position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l1104
													}
													position++
													goto l1103
												l1104:
													position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
													if buffer[position] != rune('I') {
														goto l494
													}
													position++
												}
											l1103:
												{
													position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1106
													}
													position++
													goto l1105
												l1106:
													position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
													if buffer[position] != rune('S') {
														goto l494
													}
													position++
												}
											l1105:
												{
													position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1108
													}
													position++
													goto l1107
												l1108:
													position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
													if buffer[position] != rune('T') {
														goto l494
													}
													position++
												}
											l1107:
												{
													position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1110
													}
													position++
													goto l1109
												l1110:
													position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
													if buffer[position] != rune('S') {
														goto l494
													}
													position++
												}
											l1109:
												if !rules[ruleskip]() {
													goto l494
												}
												depth--
												add(ruleNOTEXIST, position1092)
											}
										}
									l1077:
										if !rules[rulegroupGraphPattern]() {
											goto l494
										}
										break
									case 'I', 'i':
										{
											position1111 := position
											depth++
											{
												position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
												if buffer[position] != rune('i') {
													goto l1113
												}
												position++
												goto l1112
											l1113:
												position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
												if buffer[position] != rune('I') {
													goto l494
												}
												position++
											}
										l1112:
											{
												position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
												if buffer[position] != rune('f') {
													goto l1115
												}
												position++
												goto l1114
											l1115:
												position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
												if buffer[position] != rune('F') {
													goto l494
												}
												position++
											}
										l1114:
											if !rules[ruleskip]() {
												goto l494
											}
											depth--
											add(ruleIF, position1111)
										}
										if !rules[ruleLPAREN]() {
											goto l494
										}
										if !rules[ruleexpression]() {
											goto l494
										}
										if !rules[ruleCOMMA]() {
											goto l494
										}
										if !rules[ruleexpression]() {
											goto l494
										}
										if !rules[ruleCOMMA]() {
											goto l494
										}
										if !rules[ruleexpression]() {
											goto l494
										}
										if !rules[ruleRPAREN]() {
											goto l494
										}
										break
									case 'C', 'c':
										{
											position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
											{
												position1118 := position
												depth++
												{
													position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1120
													}
													position++
													goto l1119
												l1120:
													position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
													if buffer[position] != rune('C') {
														goto l1117
													}
													position++
												}
											l1119:
												{
													position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l1122
													}
													position++
													goto l1121
												l1122:
													position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
													if buffer[position] != rune('O') {
														goto l1117
													}
													position++
												}
											l1121:
												{
													position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1124
													}
													position++
													goto l1123
												l1124:
													position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
													if buffer[position] != rune('N') {
														goto l1117
													}
													position++
												}
											l1123:
												{
													position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1126
													}
													position++
													goto l1125
												l1126:
													position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
													if buffer[position] != rune('C') {
														goto l1117
													}
													position++
												}
											l1125:
												{
													position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1128
													}
													position++
													goto l1127
												l1128:
													position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
													if buffer[position] != rune('A') {
														goto l1117
													}
													position++
												}
											l1127:
												{
													position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1130
													}
													position++
													goto l1129
												l1130:
													position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
													if buffer[position] != rune('T') {
														goto l1117
													}
													position++
												}
											l1129:
												if !rules[ruleskip]() {
													goto l1117
												}
												depth--
												add(ruleCONCAT, position1118)
											}
											goto l1116
										l1117:
											position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
											{
												position1131 := position
												depth++
												{
													position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1133
													}
													position++
													goto l1132
												l1133:
													position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
													if buffer[position] != rune('C') {
														goto l494
													}
													position++
												}
											l1132:
												{
													position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l1135
													}
													position++
													goto l1134
												l1135:
													position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
													if buffer[position] != rune('O') {
														goto l494
													}
													position++
												}
											l1134:
												{
													position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1137
													}
													position++
													goto l1136
												l1137:
													position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
													if buffer[position] != rune('A') {
														goto l494
													}
													position++
												}
											l1136:
												{
													position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l1139
													}
													position++
													goto l1138
												l1139:
													position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
													if buffer[position] != rune('L') {
														goto l494
													}
													position++
												}
											l1138:
												{
													position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1141
													}
													position++
													goto l1140
												l1141:
													position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
													if buffer[position] != rune('E') {
														goto l494
													}
													position++
												}
											l1140:
												{
													position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1143
													}
													position++
													goto l1142
												l1143:
													position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
													if buffer[position] != rune('S') {
														goto l494
													}
													position++
												}
											l1142:
												{
													position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1145
													}
													position++
													goto l1144
												l1145:
													position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
													if buffer[position] != rune('C') {
														goto l494
													}
													position++
												}
											l1144:
												{
													position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1147
													}
													position++
													goto l1146
												l1147:
													position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
													if buffer[position] != rune('E') {
														goto l494
													}
													position++
												}
											l1146:
												if !rules[ruleskip]() {
													goto l494
												}
												depth--
												add(ruleCOALESCE, position1131)
											}
										}
									l1116:
										if !rules[ruleargList]() {
											goto l494
										}
										break
									case 'B', 'b':
										{
											position1148 := position
											depth++
											{
												position1149, tokenIndex1149, depth1149 := position, tokenIndex, depth
												if buffer[position] != rune('b') {
													goto l1150
												}
												position++
												goto l1149
											l1150:
												position, tokenIndex, depth = position1149, tokenIndex1149, depth1149
												if buffer[position] != rune('B') {
													goto l494
												}
												position++
											}
										l1149:
											{
												position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
												if buffer[position] != rune('n') {
													goto l1152
												}
												position++
												goto l1151
											l1152:
												position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
												if buffer[position] != rune('N') {
													goto l494
												}
												position++
											}
										l1151:
											{
												position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
												if buffer[position] != rune('o') {
													goto l1154
												}
												position++
												goto l1153
											l1154:
												position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
												if buffer[position] != rune('O') {
													goto l494
												}
												position++
											}
										l1153:
											{
												position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1156
												}
												position++
												goto l1155
											l1156:
												position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
												if buffer[position] != rune('D') {
													goto l494
												}
												position++
											}
										l1155:
											{
												position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l1158
												}
												position++
												goto l1157
											l1158:
												position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
												if buffer[position] != rune('E') {
													goto l494
												}
												position++
											}
										l1157:
											if !rules[ruleskip]() {
												goto l494
											}
											depth--
											add(ruleBNODE, position1148)
										}
										{
											position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
											if !rules[ruleLPAREN]() {
												goto l1160
											}
											if !rules[ruleexpression]() {
												goto l1160
											}
											if !rules[ruleRPAREN]() {
												goto l1160
											}
											goto l1159
										l1160:
											position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
											if !rules[rulenil]() {
												goto l494
											}
										}
									l1159:
										break
									default:
										{
											position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
											{
												position1163 := position
												depth++
												{
													position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1165
													}
													position++
													goto l1164
												l1165:
													position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
													if buffer[position] != rune('S') {
														goto l1162
													}
													position++
												}
											l1164:
												{
													position1166, tokenIndex1166, depth1166 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l1167
													}
													position++
													goto l1166
												l1167:
													position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
													if buffer[position] != rune('U') {
														goto l1162
													}
													position++
												}
											l1166:
												{
													position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
													if buffer[position] != rune('b') {
														goto l1169
													}
													position++
													goto l1168
												l1169:
													position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
													if buffer[position] != rune('B') {
														goto l1162
													}
													position++
												}
											l1168:
												{
													position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1171
													}
													position++
													goto l1170
												l1171:
													position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
													if buffer[position] != rune('S') {
														goto l1162
													}
													position++
												}
											l1170:
												{
													position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1173
													}
													position++
													goto l1172
												l1173:
													position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
													if buffer[position] != rune('T') {
														goto l1162
													}
													position++
												}
											l1172:
												{
													position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l1175
													}
													position++
													goto l1174
												l1175:
													position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
													if buffer[position] != rune('R') {
														goto l1162
													}
													position++
												}
											l1174:
												if !rules[ruleskip]() {
													goto l1162
												}
												depth--
												add(ruleSUBSTR, position1163)
											}
											goto l1161
										l1162:
											position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
											{
												position1177 := position
												depth++
												{
													position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l1179
													}
													position++
													goto l1178
												l1179:
													position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
													if buffer[position] != rune('R') {
														goto l1176
													}
													position++
												}
											l1178:
												{
													position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1181
													}
													position++
													goto l1180
												l1181:
													position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
													if buffer[position] != rune('E') {
														goto l1176
													}
													position++
												}
											l1180:
												{
													position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
													if buffer[position] != rune('p') {
														goto l1183
													}
													position++
													goto l1182
												l1183:
													position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
													if buffer[position] != rune('P') {
														goto l1176
													}
													position++
												}
											l1182:
												{
													position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l1185
													}
													position++
													goto l1184
												l1185:
													position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
													if buffer[position] != rune('L') {
														goto l1176
													}
													position++
												}
											l1184:
												{
													position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1187
													}
													position++
													goto l1186
												l1187:
													position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
													if buffer[position] != rune('A') {
														goto l1176
													}
													position++
												}
											l1186:
												{
													position1188, tokenIndex1188, depth1188 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1189
													}
													position++
													goto l1188
												l1189:
													position, tokenIndex, depth = position1188, tokenIndex1188, depth1188
													if buffer[position] != rune('C') {
														goto l1176
													}
													position++
												}
											l1188:
												{
													position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1191
													}
													position++
													goto l1190
												l1191:
													position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
													if buffer[position] != rune('E') {
														goto l1176
													}
													position++
												}
											l1190:
												if !rules[ruleskip]() {
													goto l1176
												}
												depth--
												add(ruleREPLACE, position1177)
											}
											goto l1161
										l1176:
											position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
											{
												position1192 := position
												depth++
												{
													position1193, tokenIndex1193, depth1193 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l1194
													}
													position++
													goto l1193
												l1194:
													position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
													if buffer[position] != rune('R') {
														goto l494
													}
													position++
												}
											l1193:
												{
													position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1196
													}
													position++
													goto l1195
												l1196:
													position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
													if buffer[position] != rune('E') {
														goto l494
													}
													position++
												}
											l1195:
												{
													position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
													if buffer[position] != rune('g') {
														goto l1198
													}
													position++
													goto l1197
												l1198:
													position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
													if buffer[position] != rune('G') {
														goto l494
													}
													position++
												}
											l1197:
												{
													position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1200
													}
													position++
													goto l1199
												l1200:
													position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
													if buffer[position] != rune('E') {
														goto l494
													}
													position++
												}
											l1199:
												{
													position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
													if buffer[position] != rune('x') {
														goto l1202
													}
													position++
													goto l1201
												l1202:
													position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
													if buffer[position] != rune('X') {
														goto l494
													}
													position++
												}
											l1201:
												if !rules[ruleskip]() {
													goto l494
												}
												depth--
												add(ruleREGEX, position1192)
											}
										}
									l1161:
										if !rules[ruleLPAREN]() {
											goto l494
										}
										if !rules[ruleexpression]() {
											goto l494
										}
										if !rules[ruleCOMMA]() {
											goto l494
										}
										if !rules[ruleexpression]() {
											goto l494
										}
										{
											position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
											if !rules[ruleCOMMA]() {
												goto l1203
											}
											if !rules[ruleexpression]() {
												goto l1203
											}
											goto l1204
										l1203:
											position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
										}
									l1204:
										if !rules[ruleRPAREN]() {
											goto l494
										}
										break
									}
								}

							}
						l496:
							depth--
							add(rulebuiltinCall, position495)
						}
						goto l491
					l494:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
						{
							position1206 := position
							depth++
							if !rules[ruleiriref]() {
								goto l1205
							}
							if !rules[ruleargList]() {
								goto l1205
							}
							depth--
							add(rulefunctionCall, position1206)
						}
						goto l491
					l1205:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
						if !rules[ruleiriref]() {
							goto l1207
						}
						goto l491
					l1207:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
						{
							switch buffer[position] {
							case '$', '?':
								if !rules[rulevar]() {
									goto l485
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l485
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l485
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l485
								}
								break
							}
						}

					}
				l491:
					depth--
					add(ruleprimaryExpression, position490)
				}
				depth--
				add(ruleunaryExpression, position486)
			}
			return true
		l485:
			position, tokenIndex, depth = position485, tokenIndex485, depth485
			return false
		},
		/* 52 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 53 brackettedExpression <- <(LPAREN expression RPAREN)> */
		nil,
		/* 54 functionCall <- <(iriref argList)> */
		nil,
		/* 55 in <- <(IN argList)> */
		nil,
		/* 56 notin <- <(NOTIN argList)> */
		nil,
		/* 57 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
			{
				position1215 := position
				depth++
				{
					position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l1217
					}
					goto l1216
				l1217:
					position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
					if !rules[ruleLPAREN]() {
						goto l1214
					}
					if !rules[ruleexpression]() {
						goto l1214
					}
				l1218:
					{
						position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l1219
						}
						if !rules[ruleexpression]() {
							goto l1219
						}
						goto l1218
					l1219:
						position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
					}
					if !rules[ruleRPAREN]() {
						goto l1214
					}
				}
			l1216:
				depth--
				add(ruleargList, position1215)
			}
			return true
		l1214:
			position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
			return false
		},
		/* 58 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		nil,
		/* 59 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
			{
				position1222 := position
				depth++
				{
					position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1224
					}
					position++
					goto l1223
				l1224:
					position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
					if buffer[position] != rune('$') {
						goto l1221
					}
					position++
				}
			l1223:
				{
					position1225 := position
					depth++
					{
						position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
						{
							position1230 := position
							depth++
							{
								position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
								{
									position1233 := position
									depth++
									{
										position1234, tokenIndex1234, depth1234 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1235
										}
										position++
										goto l1234
									l1235:
										position, tokenIndex, depth = position1234, tokenIndex1234, depth1234
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1232
										}
										position++
									}
								l1234:
									depth--
									add(rulePN_CHARS_BASE, position1233)
								}
								goto l1231
							l1232:
								position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
								if buffer[position] != rune('_') {
									goto l1229
								}
								position++
							}
						l1231:
							depth--
							add(rulePN_CHARS_U, position1230)
						}
						goto l1228
					l1229:
						position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1221
						}
						position++
					}
				l1228:
				l1226:
					{
						position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
						{
							position1236, tokenIndex1236, depth1236 := position, tokenIndex, depth
							{
								position1238 := position
								depth++
								{
									position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
									{
										position1241 := position
										depth++
										{
											position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1243
											}
											position++
											goto l1242
										l1243:
											position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1240
											}
											position++
										}
									l1242:
										depth--
										add(rulePN_CHARS_BASE, position1241)
									}
									goto l1239
								l1240:
									position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
									if buffer[position] != rune('_') {
										goto l1237
									}
									position++
								}
							l1239:
								depth--
								add(rulePN_CHARS_U, position1238)
							}
							goto l1236
						l1237:
							position, tokenIndex, depth = position1236, tokenIndex1236, depth1236
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1227
							}
							position++
						}
					l1236:
						goto l1226
					l1227:
						position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
					}
					depth--
					add(ruleVARNAME, position1225)
				}
				if !rules[ruleskip]() {
					goto l1221
				}
				depth--
				add(rulevar, position1222)
			}
			return true
		l1221:
			position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
			return false
		},
		/* 60 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
			{
				position1245 := position
				depth++
				{
					position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1247
					}
					goto l1246
				l1247:
					position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
					{
						position1248 := position
						depth++
					l1249:
						{
							position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
							{
								position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
								{
									position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1253
									}
									position++
									goto l1252
								l1253:
									position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
									if buffer[position] != rune(' ') {
										goto l1251
									}
									position++
								}
							l1252:
								goto l1250
							l1251:
								position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
							}
							if !matchDot() {
								goto l1250
							}
							goto l1249
						l1250:
							position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
						}
						if buffer[position] != rune(':') {
							goto l1244
						}
						position++
					l1254:
						{
							position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
							{
								position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1257
								}
								position++
								goto l1256
							l1257:
								position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1258
								}
								position++
								goto l1256
							l1258:
								position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1259
								}
								position++
								goto l1256
							l1259:
								position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1255
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1255
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1255
										}
										position++
										break
									}
								}

							}
						l1256:
							goto l1254
						l1255:
							position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
						}
						if !rules[ruleskip]() {
							goto l1244
						}
						depth--
						add(ruleprefixedName, position1248)
					}
				}
			l1246:
				depth--
				add(ruleiriref, position1245)
			}
			return true
		l1244:
			position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
			return false
		},
		/* 61 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
			{
				position1262 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1261
				}
				position++
			l1263:
				{
					position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
					{
						position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1265
						}
						position++
						goto l1264
					l1265:
						position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
					}
					if !matchDot() {
						goto l1264
					}
					goto l1263
				l1264:
					position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
				}
				if buffer[position] != rune('>') {
					goto l1261
				}
				position++
				if !rules[ruleskip]() {
					goto l1261
				}
				depth--
				add(ruleiri, position1262)
			}
			return true
		l1261:
			position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
			return false
		},
		/* 62 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 63 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
			{
				position1268 := position
				depth++
				{
					position1269 := position
					depth++
					if buffer[position] != rune('"') {
						goto l1267
					}
					position++
				l1270:
					{
						position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
						{
							position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l1272
							}
							position++
							goto l1271
						l1272:
							position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
						}
						if !matchDot() {
							goto l1271
						}
						goto l1270
					l1271:
						position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
					}
					if buffer[position] != rune('"') {
						goto l1267
					}
					position++
					depth--
					add(rulestring, position1269)
				}
				{
					position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
					{
						position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1276
						}
						position++
						{
							position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1280
							}
							position++
							goto l1279
						l1280:
							position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1276
							}
							position++
						}
					l1279:
					l1277:
						{
							position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
							{
								position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1282
								}
								position++
								goto l1281
							l1282:
								position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1278
								}
								position++
							}
						l1281:
							goto l1277
						l1278:
							position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
						}
					l1283:
						{
							position1284, tokenIndex1284, depth1284 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1284
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1284
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1284
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1284
									}
									position++
									break
								}
							}

						l1285:
							{
								position1286, tokenIndex1286, depth1286 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1286
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1286
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1286
										}
										position++
										break
									}
								}

								goto l1285
							l1286:
								position, tokenIndex, depth = position1286, tokenIndex1286, depth1286
							}
							goto l1283
						l1284:
							position, tokenIndex, depth = position1284, tokenIndex1284, depth1284
						}
						goto l1275
					l1276:
						position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
						if buffer[position] != rune('^') {
							goto l1273
						}
						position++
						if buffer[position] != rune('^') {
							goto l1273
						}
						position++
						if !rules[ruleiriref]() {
							goto l1273
						}
					}
				l1275:
					goto l1274
				l1273:
					position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
				}
			l1274:
				if !rules[ruleskip]() {
					goto l1267
				}
				depth--
				add(ruleliteral, position1268)
			}
			return true
		l1267:
			position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
			return false
		},
		/* 64 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 65 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
			{
				position1291 := position
				depth++
				{
					position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
					{
						position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1295
						}
						position++
						goto l1294
					l1295:
						position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
						if buffer[position] != rune('-') {
							goto l1292
						}
						position++
					}
				l1294:
					goto l1293
				l1292:
					position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
				}
			l1293:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1290
				}
				position++
			l1296:
				{
					position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1297
					}
					position++
					goto l1296
				l1297:
					position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
				}
				{
					position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1298
					}
					position++
				l1300:
					{
						position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1301
						}
						position++
						goto l1300
					l1301:
						position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
					}
					goto l1299
				l1298:
					position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
				}
			l1299:
				if !rules[ruleskip]() {
					goto l1290
				}
				depth--
				add(rulenumericLiteral, position1291)
			}
			return true
		l1290:
			position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
			return false
		},
		/* 66 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 67 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
			{
				position1304 := position
				depth++
				{
					position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
					{
						position1307 := position
						depth++
						{
							position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1309
							}
							position++
							goto l1308
						l1309:
							position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
							if buffer[position] != rune('T') {
								goto l1306
							}
							position++
						}
					l1308:
						{
							position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1311
							}
							position++
							goto l1310
						l1311:
							position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
							if buffer[position] != rune('R') {
								goto l1306
							}
							position++
						}
					l1310:
						{
							position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1313
							}
							position++
							goto l1312
						l1313:
							position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
							if buffer[position] != rune('U') {
								goto l1306
							}
							position++
						}
					l1312:
						{
							position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1315
							}
							position++
							goto l1314
						l1315:
							position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
							if buffer[position] != rune('E') {
								goto l1306
							}
							position++
						}
					l1314:
						if !rules[ruleskip]() {
							goto l1306
						}
						depth--
						add(ruleTRUE, position1307)
					}
					goto l1305
				l1306:
					position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
					{
						position1316 := position
						depth++
						{
							position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1318
							}
							position++
							goto l1317
						l1318:
							position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
							if buffer[position] != rune('F') {
								goto l1303
							}
							position++
						}
					l1317:
						{
							position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1320
							}
							position++
							goto l1319
						l1320:
							position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
							if buffer[position] != rune('A') {
								goto l1303
							}
							position++
						}
					l1319:
						{
							position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1322
							}
							position++
							goto l1321
						l1322:
							position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
							if buffer[position] != rune('L') {
								goto l1303
							}
							position++
						}
					l1321:
						{
							position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1324
							}
							position++
							goto l1323
						l1324:
							position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
							if buffer[position] != rune('S') {
								goto l1303
							}
							position++
						}
					l1323:
						{
							position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1326
							}
							position++
							goto l1325
						l1326:
							position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
							if buffer[position] != rune('E') {
								goto l1303
							}
							position++
						}
					l1325:
						if !rules[ruleskip]() {
							goto l1303
						}
						depth--
						add(ruleFALSE, position1316)
					}
				}
			l1305:
				depth--
				add(rulebooleanLiteral, position1304)
			}
			return true
		l1303:
			position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
			return false
		},
		/* 68 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 69 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 70 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 71 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1330, tokenIndex1330, depth1330 := position, tokenIndex, depth
			{
				position1331 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1330
				}
				position++
			l1332:
				{
					position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1333
					}
					goto l1332
				l1333:
					position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
				}
				if buffer[position] != rune(')') {
					goto l1330
				}
				position++
				if !rules[ruleskip]() {
					goto l1330
				}
				depth--
				add(rulenil, position1331)
			}
			return true
		l1330:
			position, tokenIndex, depth = position1330, tokenIndex1330, depth1330
			return false
		},
		/* 72 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 73 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 74 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 75 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 76 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 77 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 78 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 79 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 80 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 81 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 82 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 83 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 84 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 85 LBRACE <- <('{' skip)> */
		func() bool {
			position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
			{
				position1348 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1347
				}
				position++
				if !rules[ruleskip]() {
					goto l1347
				}
				depth--
				add(ruleLBRACE, position1348)
			}
			return true
		l1347:
			position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
			return false
		},
		/* 86 RBRACE <- <('}' skip)> */
		func() bool {
			position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
			{
				position1350 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1349
				}
				position++
				if !rules[ruleskip]() {
					goto l1349
				}
				depth--
				add(ruleRBRACE, position1350)
			}
			return true
		l1349:
			position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
			return false
		},
		/* 87 LBRACK <- <('[' skip)> */
		nil,
		/* 88 RBRACK <- <(']' skip)> */
		nil,
		/* 89 SEMICOLON <- <(';' skip)> */
		nil,
		/* 90 COMMA <- <(',' skip)> */
		func() bool {
			position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
			{
				position1355 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1354
				}
				position++
				if !rules[ruleskip]() {
					goto l1354
				}
				depth--
				add(ruleCOMMA, position1355)
			}
			return true
		l1354:
			position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
			return false
		},
		/* 91 DOT <- <('.' skip)> */
		func() bool {
			position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
			{
				position1357 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1356
				}
				position++
				if !rules[ruleskip]() {
					goto l1356
				}
				depth--
				add(ruleDOT, position1357)
			}
			return true
		l1356:
			position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
			return false
		},
		/* 92 COLON <- <(':' skip)> */
		nil,
		/* 93 PIPE <- <('|' skip)> */
		func() bool {
			position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
			{
				position1360 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1359
				}
				position++
				if !rules[ruleskip]() {
					goto l1359
				}
				depth--
				add(rulePIPE, position1360)
			}
			return true
		l1359:
			position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
			return false
		},
		/* 94 SLASH <- <('/' skip)> */
		func() bool {
			position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
			{
				position1362 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1361
				}
				position++
				if !rules[ruleskip]() {
					goto l1361
				}
				depth--
				add(ruleSLASH, position1362)
			}
			return true
		l1361:
			position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
			return false
		},
		/* 95 INVERSE <- <('^' skip)> */
		func() bool {
			position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
			{
				position1364 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1363
				}
				position++
				if !rules[ruleskip]() {
					goto l1363
				}
				depth--
				add(ruleINVERSE, position1364)
			}
			return true
		l1363:
			position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
			return false
		},
		/* 96 LPAREN <- <('(' skip)> */
		func() bool {
			position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
			{
				position1366 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1365
				}
				position++
				if !rules[ruleskip]() {
					goto l1365
				}
				depth--
				add(ruleLPAREN, position1366)
			}
			return true
		l1365:
			position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
			return false
		},
		/* 97 RPAREN <- <(')' skip)> */
		func() bool {
			position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
			{
				position1368 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1367
				}
				position++
				if !rules[ruleskip]() {
					goto l1367
				}
				depth--
				add(ruleRPAREN, position1368)
			}
			return true
		l1367:
			position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
			return false
		},
		/* 98 ISA <- <('a' skip)> */
		func() bool {
			position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
			{
				position1370 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1369
				}
				position++
				if !rules[ruleskip]() {
					goto l1369
				}
				depth--
				add(ruleISA, position1370)
			}
			return true
		l1369:
			position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
			return false
		},
		/* 99 NOT <- <('!' skip)> */
		func() bool {
			position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
			{
				position1372 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1371
				}
				position++
				if !rules[ruleskip]() {
					goto l1371
				}
				depth--
				add(ruleNOT, position1372)
			}
			return true
		l1371:
			position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
			return false
		},
		/* 100 STAR <- <('*' skip)> */
		func() bool {
			position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
			{
				position1374 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1373
				}
				position++
				if !rules[ruleskip]() {
					goto l1373
				}
				depth--
				add(ruleSTAR, position1374)
			}
			return true
		l1373:
			position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
			return false
		},
		/* 101 PLUS <- <('+' skip)> */
		func() bool {
			position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
			{
				position1376 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1375
				}
				position++
				if !rules[ruleskip]() {
					goto l1375
				}
				depth--
				add(rulePLUS, position1376)
			}
			return true
		l1375:
			position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
			return false
		},
		/* 102 MINUS <- <('-' skip)> */
		func() bool {
			position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
			{
				position1378 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1377
				}
				position++
				if !rules[ruleskip]() {
					goto l1377
				}
				depth--
				add(ruleMINUS, position1378)
			}
			return true
		l1377:
			position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
			return false
		},
		/* 103 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 104 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 105 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 106 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 107 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
			{
				position1384 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1383
				}
				position++
			l1385:
				{
					position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1386
					}
					position++
					goto l1385
				l1386:
					position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
				}
				if !rules[ruleskip]() {
					goto l1383
				}
				depth--
				add(ruleINTEGER, position1384)
			}
			return true
		l1383:
			position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
			return false
		},
		/* 108 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 109 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 110 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 111 OR <- <('|' '|' skip)> */
		nil,
		/* 112 AND <- <('&' '&' skip)> */
		nil,
		/* 113 EQ <- <('=' skip)> */
		nil,
		/* 114 NE <- <('!' '=' skip)> */
		nil,
		/* 115 GT <- <('>' skip)> */
		nil,
		/* 116 LT <- <('<' skip)> */
		nil,
		/* 117 LE <- <('<' '=' skip)> */
		nil,
		/* 118 GE <- <('>' '=' skip)> */
		nil,
		/* 119 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 120 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 121 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		nil,
		/* 122 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 123 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 124 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 125 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 126 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 127 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 128 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 129 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 130 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 131 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 132 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 133 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 134 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 135 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 136 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 137 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 138 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 139 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 140 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 141 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 142 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 143 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 144 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 145 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 146 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 147 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 148 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 149 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 150 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 151 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 152 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 153 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 154 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 155 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 156 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 157 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 158 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 159 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 160 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 161 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 162 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 163 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 164 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 165 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 166 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 167 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 168 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 169 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 170 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 171 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 172 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 173 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 174 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 175 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 176 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1456 := position
				depth++
			l1457:
				{
					position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
					{
						position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1460
						}
						goto l1459
					l1460:
						position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
						{
							position1461 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1458
							}
							position++
						l1462:
							{
								position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
								{
									position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1464
									}
									goto l1463
								l1464:
									position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
								}
								if !matchDot() {
									goto l1463
								}
								goto l1462
							l1463:
								position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
							}
							if !rules[ruleendOfLine]() {
								goto l1458
							}
							depth--
							add(rulecomment, position1461)
						}
					}
				l1459:
					goto l1457
				l1458:
					position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
				}
				depth--
				add(ruleskip, position1456)
			}
			return true
		},
		/* 177 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
			{
				position1466 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1465
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1465
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1465
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1465
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1465
						}
						break
					}
				}

				depth--
				add(rulews, position1466)
			}
			return true
		l1465:
			position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
			return false
		},
		/* 178 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 179 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
			{
				position1470 := position
				depth++
				{
					position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1472
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1472
					}
					position++
					goto l1471
				l1472:
					position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
					if buffer[position] != rune('\n') {
						goto l1473
					}
					position++
					goto l1471
				l1473:
					position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
					if buffer[position] != rune('\r') {
						goto l1469
					}
					position++
				}
			l1471:
				depth--
				add(ruleendOfLine, position1470)
			}
			return true
		l1469:
			position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
			return false
		},
	}
	p.rules = rules
}
