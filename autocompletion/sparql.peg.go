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
	rulefunctionCall
	rulein
	rulenotin
	ruleargList
	rulebuiltinCall
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
	"functionCall",
	"in",
	"notin",
	"argList",
	"builtinCall",
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

	Buffer string
	buffer []rune
	rules  [197]func() bool
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
			p.addPrefix(p.skipComment(buffer, begin, end))
		case ruleAction1:
			p.setSubject(p.skipComment(buffer, begin, end))
		case ruleAction2:
			p.setSubject(p.skipComment(buffer, begin, end))
		case ruleAction3:
			p.setSubject("?POF")
		case ruleAction4:
			p.setPredicate("?POF")
		case ruleAction5:
			p.setPredicate(p.skipComment(buffer, begin, end))
		case ruleAction6:
			p.setPredicate(p.skipComment(buffer, begin, end))
		case ruleAction7:
			p.setObject("?POF")
			p.addTriplePattern()
		case ruleAction8:
			p.setObject(p.skipComment(buffer, begin, end))
			p.addTriplePattern()
		case ruleAction9:
			p.setObject("?FillVar")
			p.addTriplePattern()
		case ruleAction10:
			p.setPrefix(p.skipComment(buffer, begin, end))
		case ruleAction11:
			p.setPathLength(p.skipComment(buffer, begin, end))
		case ruleAction12:
			p.setKeyword(p.skipComment(buffer, begin, end))
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
						{
							position518 := position
							depth++
							{
								position519, tokenIndex519, depth519 := position, tokenIndex, depth
								{
									position521, tokenIndex521, depth521 := position, tokenIndex, depth
									{
										position523 := position
										depth++
										{
											position524, tokenIndex524, depth524 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l525
											}
											position++
											goto l524
										l525:
											position, tokenIndex, depth = position524, tokenIndex524, depth524
											if buffer[position] != rune('S') {
												goto l522
											}
											position++
										}
									l524:
										{
											position526, tokenIndex526, depth526 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l527
											}
											position++
											goto l526
										l527:
											position, tokenIndex, depth = position526, tokenIndex526, depth526
											if buffer[position] != rune('T') {
												goto l522
											}
											position++
										}
									l526:
										{
											position528, tokenIndex528, depth528 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l529
											}
											position++
											goto l528
										l529:
											position, tokenIndex, depth = position528, tokenIndex528, depth528
											if buffer[position] != rune('R') {
												goto l522
											}
											position++
										}
									l528:
										if !rules[ruleskip]() {
											goto l522
										}
										depth--
										add(ruleSTR, position523)
									}
									goto l521
								l522:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position531 := position
										depth++
										{
											position532, tokenIndex532, depth532 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l533
											}
											position++
											goto l532
										l533:
											position, tokenIndex, depth = position532, tokenIndex532, depth532
											if buffer[position] != rune('L') {
												goto l530
											}
											position++
										}
									l532:
										{
											position534, tokenIndex534, depth534 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l535
											}
											position++
											goto l534
										l535:
											position, tokenIndex, depth = position534, tokenIndex534, depth534
											if buffer[position] != rune('A') {
												goto l530
											}
											position++
										}
									l534:
										{
											position536, tokenIndex536, depth536 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l537
											}
											position++
											goto l536
										l537:
											position, tokenIndex, depth = position536, tokenIndex536, depth536
											if buffer[position] != rune('N') {
												goto l530
											}
											position++
										}
									l536:
										{
											position538, tokenIndex538, depth538 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l539
											}
											position++
											goto l538
										l539:
											position, tokenIndex, depth = position538, tokenIndex538, depth538
											if buffer[position] != rune('G') {
												goto l530
											}
											position++
										}
									l538:
										if !rules[ruleskip]() {
											goto l530
										}
										depth--
										add(ruleLANG, position531)
									}
									goto l521
								l530:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position541 := position
										depth++
										{
											position542, tokenIndex542, depth542 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l543
											}
											position++
											goto l542
										l543:
											position, tokenIndex, depth = position542, tokenIndex542, depth542
											if buffer[position] != rune('D') {
												goto l540
											}
											position++
										}
									l542:
										{
											position544, tokenIndex544, depth544 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l545
											}
											position++
											goto l544
										l545:
											position, tokenIndex, depth = position544, tokenIndex544, depth544
											if buffer[position] != rune('A') {
												goto l540
											}
											position++
										}
									l544:
										{
											position546, tokenIndex546, depth546 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l547
											}
											position++
											goto l546
										l547:
											position, tokenIndex, depth = position546, tokenIndex546, depth546
											if buffer[position] != rune('T') {
												goto l540
											}
											position++
										}
									l546:
										{
											position548, tokenIndex548, depth548 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l549
											}
											position++
											goto l548
										l549:
											position, tokenIndex, depth = position548, tokenIndex548, depth548
											if buffer[position] != rune('A') {
												goto l540
											}
											position++
										}
									l548:
										{
											position550, tokenIndex550, depth550 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l551
											}
											position++
											goto l550
										l551:
											position, tokenIndex, depth = position550, tokenIndex550, depth550
											if buffer[position] != rune('T') {
												goto l540
											}
											position++
										}
									l550:
										{
											position552, tokenIndex552, depth552 := position, tokenIndex, depth
											if buffer[position] != rune('y') {
												goto l553
											}
											position++
											goto l552
										l553:
											position, tokenIndex, depth = position552, tokenIndex552, depth552
											if buffer[position] != rune('Y') {
												goto l540
											}
											position++
										}
									l552:
										{
											position554, tokenIndex554, depth554 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l555
											}
											position++
											goto l554
										l555:
											position, tokenIndex, depth = position554, tokenIndex554, depth554
											if buffer[position] != rune('P') {
												goto l540
											}
											position++
										}
									l554:
										{
											position556, tokenIndex556, depth556 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l557
											}
											position++
											goto l556
										l557:
											position, tokenIndex, depth = position556, tokenIndex556, depth556
											if buffer[position] != rune('E') {
												goto l540
											}
											position++
										}
									l556:
										if !rules[ruleskip]() {
											goto l540
										}
										depth--
										add(ruleDATATYPE, position541)
									}
									goto l521
								l540:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position559 := position
										depth++
										{
											position560, tokenIndex560, depth560 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l561
											}
											position++
											goto l560
										l561:
											position, tokenIndex, depth = position560, tokenIndex560, depth560
											if buffer[position] != rune('I') {
												goto l558
											}
											position++
										}
									l560:
										{
											position562, tokenIndex562, depth562 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l563
											}
											position++
											goto l562
										l563:
											position, tokenIndex, depth = position562, tokenIndex562, depth562
											if buffer[position] != rune('R') {
												goto l558
											}
											position++
										}
									l562:
										{
											position564, tokenIndex564, depth564 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l565
											}
											position++
											goto l564
										l565:
											position, tokenIndex, depth = position564, tokenIndex564, depth564
											if buffer[position] != rune('I') {
												goto l558
											}
											position++
										}
									l564:
										if !rules[ruleskip]() {
											goto l558
										}
										depth--
										add(ruleIRI, position559)
									}
									goto l521
								l558:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position567 := position
										depth++
										{
											position568, tokenIndex568, depth568 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l569
											}
											position++
											goto l568
										l569:
											position, tokenIndex, depth = position568, tokenIndex568, depth568
											if buffer[position] != rune('U') {
												goto l566
											}
											position++
										}
									l568:
										{
											position570, tokenIndex570, depth570 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l571
											}
											position++
											goto l570
										l571:
											position, tokenIndex, depth = position570, tokenIndex570, depth570
											if buffer[position] != rune('R') {
												goto l566
											}
											position++
										}
									l570:
										{
											position572, tokenIndex572, depth572 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l573
											}
											position++
											goto l572
										l573:
											position, tokenIndex, depth = position572, tokenIndex572, depth572
											if buffer[position] != rune('I') {
												goto l566
											}
											position++
										}
									l572:
										if !rules[ruleskip]() {
											goto l566
										}
										depth--
										add(ruleURI, position567)
									}
									goto l521
								l566:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position575 := position
										depth++
										{
											position576, tokenIndex576, depth576 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l577
											}
											position++
											goto l576
										l577:
											position, tokenIndex, depth = position576, tokenIndex576, depth576
											if buffer[position] != rune('S') {
												goto l574
											}
											position++
										}
									l576:
										{
											position578, tokenIndex578, depth578 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l579
											}
											position++
											goto l578
										l579:
											position, tokenIndex, depth = position578, tokenIndex578, depth578
											if buffer[position] != rune('T') {
												goto l574
											}
											position++
										}
									l578:
										{
											position580, tokenIndex580, depth580 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l581
											}
											position++
											goto l580
										l581:
											position, tokenIndex, depth = position580, tokenIndex580, depth580
											if buffer[position] != rune('R') {
												goto l574
											}
											position++
										}
									l580:
										{
											position582, tokenIndex582, depth582 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l583
											}
											position++
											goto l582
										l583:
											position, tokenIndex, depth = position582, tokenIndex582, depth582
											if buffer[position] != rune('L') {
												goto l574
											}
											position++
										}
									l582:
										{
											position584, tokenIndex584, depth584 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l585
											}
											position++
											goto l584
										l585:
											position, tokenIndex, depth = position584, tokenIndex584, depth584
											if buffer[position] != rune('E') {
												goto l574
											}
											position++
										}
									l584:
										{
											position586, tokenIndex586, depth586 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l587
											}
											position++
											goto l586
										l587:
											position, tokenIndex, depth = position586, tokenIndex586, depth586
											if buffer[position] != rune('N') {
												goto l574
											}
											position++
										}
									l586:
										if !rules[ruleskip]() {
											goto l574
										}
										depth--
										add(ruleSTRLEN, position575)
									}
									goto l521
								l574:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position589 := position
										depth++
										{
											position590, tokenIndex590, depth590 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l591
											}
											position++
											goto l590
										l591:
											position, tokenIndex, depth = position590, tokenIndex590, depth590
											if buffer[position] != rune('M') {
												goto l588
											}
											position++
										}
									l590:
										{
											position592, tokenIndex592, depth592 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l593
											}
											position++
											goto l592
										l593:
											position, tokenIndex, depth = position592, tokenIndex592, depth592
											if buffer[position] != rune('O') {
												goto l588
											}
											position++
										}
									l592:
										{
											position594, tokenIndex594, depth594 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l595
											}
											position++
											goto l594
										l595:
											position, tokenIndex, depth = position594, tokenIndex594, depth594
											if buffer[position] != rune('N') {
												goto l588
											}
											position++
										}
									l594:
										{
											position596, tokenIndex596, depth596 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l597
											}
											position++
											goto l596
										l597:
											position, tokenIndex, depth = position596, tokenIndex596, depth596
											if buffer[position] != rune('T') {
												goto l588
											}
											position++
										}
									l596:
										{
											position598, tokenIndex598, depth598 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l599
											}
											position++
											goto l598
										l599:
											position, tokenIndex, depth = position598, tokenIndex598, depth598
											if buffer[position] != rune('H') {
												goto l588
											}
											position++
										}
									l598:
										if !rules[ruleskip]() {
											goto l588
										}
										depth--
										add(ruleMONTH, position589)
									}
									goto l521
								l588:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position601 := position
										depth++
										{
											position602, tokenIndex602, depth602 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l603
											}
											position++
											goto l602
										l603:
											position, tokenIndex, depth = position602, tokenIndex602, depth602
											if buffer[position] != rune('M') {
												goto l600
											}
											position++
										}
									l602:
										{
											position604, tokenIndex604, depth604 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l605
											}
											position++
											goto l604
										l605:
											position, tokenIndex, depth = position604, tokenIndex604, depth604
											if buffer[position] != rune('I') {
												goto l600
											}
											position++
										}
									l604:
										{
											position606, tokenIndex606, depth606 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l607
											}
											position++
											goto l606
										l607:
											position, tokenIndex, depth = position606, tokenIndex606, depth606
											if buffer[position] != rune('N') {
												goto l600
											}
											position++
										}
									l606:
										{
											position608, tokenIndex608, depth608 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l609
											}
											position++
											goto l608
										l609:
											position, tokenIndex, depth = position608, tokenIndex608, depth608
											if buffer[position] != rune('U') {
												goto l600
											}
											position++
										}
									l608:
										{
											position610, tokenIndex610, depth610 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l611
											}
											position++
											goto l610
										l611:
											position, tokenIndex, depth = position610, tokenIndex610, depth610
											if buffer[position] != rune('T') {
												goto l600
											}
											position++
										}
									l610:
										{
											position612, tokenIndex612, depth612 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l613
											}
											position++
											goto l612
										l613:
											position, tokenIndex, depth = position612, tokenIndex612, depth612
											if buffer[position] != rune('E') {
												goto l600
											}
											position++
										}
									l612:
										{
											position614, tokenIndex614, depth614 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l615
											}
											position++
											goto l614
										l615:
											position, tokenIndex, depth = position614, tokenIndex614, depth614
											if buffer[position] != rune('S') {
												goto l600
											}
											position++
										}
									l614:
										if !rules[ruleskip]() {
											goto l600
										}
										depth--
										add(ruleMINUTES, position601)
									}
									goto l521
								l600:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position617 := position
										depth++
										{
											position618, tokenIndex618, depth618 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l619
											}
											position++
											goto l618
										l619:
											position, tokenIndex, depth = position618, tokenIndex618, depth618
											if buffer[position] != rune('S') {
												goto l616
											}
											position++
										}
									l618:
										{
											position620, tokenIndex620, depth620 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l621
											}
											position++
											goto l620
										l621:
											position, tokenIndex, depth = position620, tokenIndex620, depth620
											if buffer[position] != rune('E') {
												goto l616
											}
											position++
										}
									l620:
										{
											position622, tokenIndex622, depth622 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l623
											}
											position++
											goto l622
										l623:
											position, tokenIndex, depth = position622, tokenIndex622, depth622
											if buffer[position] != rune('C') {
												goto l616
											}
											position++
										}
									l622:
										{
											position624, tokenIndex624, depth624 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l625
											}
											position++
											goto l624
										l625:
											position, tokenIndex, depth = position624, tokenIndex624, depth624
											if buffer[position] != rune('O') {
												goto l616
											}
											position++
										}
									l624:
										{
											position626, tokenIndex626, depth626 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l627
											}
											position++
											goto l626
										l627:
											position, tokenIndex, depth = position626, tokenIndex626, depth626
											if buffer[position] != rune('N') {
												goto l616
											}
											position++
										}
									l626:
										{
											position628, tokenIndex628, depth628 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l629
											}
											position++
											goto l628
										l629:
											position, tokenIndex, depth = position628, tokenIndex628, depth628
											if buffer[position] != rune('D') {
												goto l616
											}
											position++
										}
									l628:
										{
											position630, tokenIndex630, depth630 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l631
											}
											position++
											goto l630
										l631:
											position, tokenIndex, depth = position630, tokenIndex630, depth630
											if buffer[position] != rune('S') {
												goto l616
											}
											position++
										}
									l630:
										if !rules[ruleskip]() {
											goto l616
										}
										depth--
										add(ruleSECONDS, position617)
									}
									goto l521
								l616:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position633 := position
										depth++
										{
											position634, tokenIndex634, depth634 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l635
											}
											position++
											goto l634
										l635:
											position, tokenIndex, depth = position634, tokenIndex634, depth634
											if buffer[position] != rune('T') {
												goto l632
											}
											position++
										}
									l634:
										{
											position636, tokenIndex636, depth636 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l637
											}
											position++
											goto l636
										l637:
											position, tokenIndex, depth = position636, tokenIndex636, depth636
											if buffer[position] != rune('I') {
												goto l632
											}
											position++
										}
									l636:
										{
											position638, tokenIndex638, depth638 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l639
											}
											position++
											goto l638
										l639:
											position, tokenIndex, depth = position638, tokenIndex638, depth638
											if buffer[position] != rune('M') {
												goto l632
											}
											position++
										}
									l638:
										{
											position640, tokenIndex640, depth640 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l641
											}
											position++
											goto l640
										l641:
											position, tokenIndex, depth = position640, tokenIndex640, depth640
											if buffer[position] != rune('E') {
												goto l632
											}
											position++
										}
									l640:
										{
											position642, tokenIndex642, depth642 := position, tokenIndex, depth
											if buffer[position] != rune('z') {
												goto l643
											}
											position++
											goto l642
										l643:
											position, tokenIndex, depth = position642, tokenIndex642, depth642
											if buffer[position] != rune('Z') {
												goto l632
											}
											position++
										}
									l642:
										{
											position644, tokenIndex644, depth644 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l645
											}
											position++
											goto l644
										l645:
											position, tokenIndex, depth = position644, tokenIndex644, depth644
											if buffer[position] != rune('O') {
												goto l632
											}
											position++
										}
									l644:
										{
											position646, tokenIndex646, depth646 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l647
											}
											position++
											goto l646
										l647:
											position, tokenIndex, depth = position646, tokenIndex646, depth646
											if buffer[position] != rune('N') {
												goto l632
											}
											position++
										}
									l646:
										{
											position648, tokenIndex648, depth648 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l649
											}
											position++
											goto l648
										l649:
											position, tokenIndex, depth = position648, tokenIndex648, depth648
											if buffer[position] != rune('E') {
												goto l632
											}
											position++
										}
									l648:
										if !rules[ruleskip]() {
											goto l632
										}
										depth--
										add(ruleTIMEZONE, position633)
									}
									goto l521
								l632:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position651 := position
										depth++
										{
											position652, tokenIndex652, depth652 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l653
											}
											position++
											goto l652
										l653:
											position, tokenIndex, depth = position652, tokenIndex652, depth652
											if buffer[position] != rune('S') {
												goto l650
											}
											position++
										}
									l652:
										{
											position654, tokenIndex654, depth654 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l655
											}
											position++
											goto l654
										l655:
											position, tokenIndex, depth = position654, tokenIndex654, depth654
											if buffer[position] != rune('H') {
												goto l650
											}
											position++
										}
									l654:
										{
											position656, tokenIndex656, depth656 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l657
											}
											position++
											goto l656
										l657:
											position, tokenIndex, depth = position656, tokenIndex656, depth656
											if buffer[position] != rune('A') {
												goto l650
											}
											position++
										}
									l656:
										if buffer[position] != rune('1') {
											goto l650
										}
										position++
										if !rules[ruleskip]() {
											goto l650
										}
										depth--
										add(ruleSHA1, position651)
									}
									goto l521
								l650:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position659 := position
										depth++
										{
											position660, tokenIndex660, depth660 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l661
											}
											position++
											goto l660
										l661:
											position, tokenIndex, depth = position660, tokenIndex660, depth660
											if buffer[position] != rune('S') {
												goto l658
											}
											position++
										}
									l660:
										{
											position662, tokenIndex662, depth662 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l663
											}
											position++
											goto l662
										l663:
											position, tokenIndex, depth = position662, tokenIndex662, depth662
											if buffer[position] != rune('H') {
												goto l658
											}
											position++
										}
									l662:
										{
											position664, tokenIndex664, depth664 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l665
											}
											position++
											goto l664
										l665:
											position, tokenIndex, depth = position664, tokenIndex664, depth664
											if buffer[position] != rune('A') {
												goto l658
											}
											position++
										}
									l664:
										if buffer[position] != rune('2') {
											goto l658
										}
										position++
										if buffer[position] != rune('5') {
											goto l658
										}
										position++
										if buffer[position] != rune('6') {
											goto l658
										}
										position++
										if !rules[ruleskip]() {
											goto l658
										}
										depth--
										add(ruleSHA256, position659)
									}
									goto l521
								l658:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position667 := position
										depth++
										{
											position668, tokenIndex668, depth668 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l669
											}
											position++
											goto l668
										l669:
											position, tokenIndex, depth = position668, tokenIndex668, depth668
											if buffer[position] != rune('S') {
												goto l666
											}
											position++
										}
									l668:
										{
											position670, tokenIndex670, depth670 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l671
											}
											position++
											goto l670
										l671:
											position, tokenIndex, depth = position670, tokenIndex670, depth670
											if buffer[position] != rune('H') {
												goto l666
											}
											position++
										}
									l670:
										{
											position672, tokenIndex672, depth672 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l673
											}
											position++
											goto l672
										l673:
											position, tokenIndex, depth = position672, tokenIndex672, depth672
											if buffer[position] != rune('A') {
												goto l666
											}
											position++
										}
									l672:
										if buffer[position] != rune('3') {
											goto l666
										}
										position++
										if buffer[position] != rune('8') {
											goto l666
										}
										position++
										if buffer[position] != rune('4') {
											goto l666
										}
										position++
										if !rules[ruleskip]() {
											goto l666
										}
										depth--
										add(ruleSHA384, position667)
									}
									goto l521
								l666:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position675 := position
										depth++
										{
											position676, tokenIndex676, depth676 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l677
											}
											position++
											goto l676
										l677:
											position, tokenIndex, depth = position676, tokenIndex676, depth676
											if buffer[position] != rune('I') {
												goto l674
											}
											position++
										}
									l676:
										{
											position678, tokenIndex678, depth678 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l679
											}
											position++
											goto l678
										l679:
											position, tokenIndex, depth = position678, tokenIndex678, depth678
											if buffer[position] != rune('S') {
												goto l674
											}
											position++
										}
									l678:
										{
											position680, tokenIndex680, depth680 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l681
											}
											position++
											goto l680
										l681:
											position, tokenIndex, depth = position680, tokenIndex680, depth680
											if buffer[position] != rune('I') {
												goto l674
											}
											position++
										}
									l680:
										{
											position682, tokenIndex682, depth682 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l683
											}
											position++
											goto l682
										l683:
											position, tokenIndex, depth = position682, tokenIndex682, depth682
											if buffer[position] != rune('R') {
												goto l674
											}
											position++
										}
									l682:
										{
											position684, tokenIndex684, depth684 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l685
											}
											position++
											goto l684
										l685:
											position, tokenIndex, depth = position684, tokenIndex684, depth684
											if buffer[position] != rune('I') {
												goto l674
											}
											position++
										}
									l684:
										if !rules[ruleskip]() {
											goto l674
										}
										depth--
										add(ruleISIRI, position675)
									}
									goto l521
								l674:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position687 := position
										depth++
										{
											position688, tokenIndex688, depth688 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l689
											}
											position++
											goto l688
										l689:
											position, tokenIndex, depth = position688, tokenIndex688, depth688
											if buffer[position] != rune('I') {
												goto l686
											}
											position++
										}
									l688:
										{
											position690, tokenIndex690, depth690 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l691
											}
											position++
											goto l690
										l691:
											position, tokenIndex, depth = position690, tokenIndex690, depth690
											if buffer[position] != rune('S') {
												goto l686
											}
											position++
										}
									l690:
										{
											position692, tokenIndex692, depth692 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l693
											}
											position++
											goto l692
										l693:
											position, tokenIndex, depth = position692, tokenIndex692, depth692
											if buffer[position] != rune('U') {
												goto l686
											}
											position++
										}
									l692:
										{
											position694, tokenIndex694, depth694 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l695
											}
											position++
											goto l694
										l695:
											position, tokenIndex, depth = position694, tokenIndex694, depth694
											if buffer[position] != rune('R') {
												goto l686
											}
											position++
										}
									l694:
										{
											position696, tokenIndex696, depth696 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l697
											}
											position++
											goto l696
										l697:
											position, tokenIndex, depth = position696, tokenIndex696, depth696
											if buffer[position] != rune('I') {
												goto l686
											}
											position++
										}
									l696:
										if !rules[ruleskip]() {
											goto l686
										}
										depth--
										add(ruleISURI, position687)
									}
									goto l521
								l686:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position699 := position
										depth++
										{
											position700, tokenIndex700, depth700 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l701
											}
											position++
											goto l700
										l701:
											position, tokenIndex, depth = position700, tokenIndex700, depth700
											if buffer[position] != rune('I') {
												goto l698
											}
											position++
										}
									l700:
										{
											position702, tokenIndex702, depth702 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l703
											}
											position++
											goto l702
										l703:
											position, tokenIndex, depth = position702, tokenIndex702, depth702
											if buffer[position] != rune('S') {
												goto l698
											}
											position++
										}
									l702:
										{
											position704, tokenIndex704, depth704 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l705
											}
											position++
											goto l704
										l705:
											position, tokenIndex, depth = position704, tokenIndex704, depth704
											if buffer[position] != rune('B') {
												goto l698
											}
											position++
										}
									l704:
										{
											position706, tokenIndex706, depth706 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l707
											}
											position++
											goto l706
										l707:
											position, tokenIndex, depth = position706, tokenIndex706, depth706
											if buffer[position] != rune('L') {
												goto l698
											}
											position++
										}
									l706:
										{
											position708, tokenIndex708, depth708 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l709
											}
											position++
											goto l708
										l709:
											position, tokenIndex, depth = position708, tokenIndex708, depth708
											if buffer[position] != rune('A') {
												goto l698
											}
											position++
										}
									l708:
										{
											position710, tokenIndex710, depth710 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l711
											}
											position++
											goto l710
										l711:
											position, tokenIndex, depth = position710, tokenIndex710, depth710
											if buffer[position] != rune('N') {
												goto l698
											}
											position++
										}
									l710:
										{
											position712, tokenIndex712, depth712 := position, tokenIndex, depth
											if buffer[position] != rune('k') {
												goto l713
											}
											position++
											goto l712
										l713:
											position, tokenIndex, depth = position712, tokenIndex712, depth712
											if buffer[position] != rune('K') {
												goto l698
											}
											position++
										}
									l712:
										if !rules[ruleskip]() {
											goto l698
										}
										depth--
										add(ruleISBLANK, position699)
									}
									goto l521
								l698:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										position715 := position
										depth++
										{
											position716, tokenIndex716, depth716 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l717
											}
											position++
											goto l716
										l717:
											position, tokenIndex, depth = position716, tokenIndex716, depth716
											if buffer[position] != rune('I') {
												goto l714
											}
											position++
										}
									l716:
										{
											position718, tokenIndex718, depth718 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l719
											}
											position++
											goto l718
										l719:
											position, tokenIndex, depth = position718, tokenIndex718, depth718
											if buffer[position] != rune('S') {
												goto l714
											}
											position++
										}
									l718:
										{
											position720, tokenIndex720, depth720 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l721
											}
											position++
											goto l720
										l721:
											position, tokenIndex, depth = position720, tokenIndex720, depth720
											if buffer[position] != rune('L') {
												goto l714
											}
											position++
										}
									l720:
										{
											position722, tokenIndex722, depth722 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l723
											}
											position++
											goto l722
										l723:
											position, tokenIndex, depth = position722, tokenIndex722, depth722
											if buffer[position] != rune('I') {
												goto l714
											}
											position++
										}
									l722:
										{
											position724, tokenIndex724, depth724 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l725
											}
											position++
											goto l724
										l725:
											position, tokenIndex, depth = position724, tokenIndex724, depth724
											if buffer[position] != rune('T') {
												goto l714
											}
											position++
										}
									l724:
										{
											position726, tokenIndex726, depth726 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l727
											}
											position++
											goto l726
										l727:
											position, tokenIndex, depth = position726, tokenIndex726, depth726
											if buffer[position] != rune('E') {
												goto l714
											}
											position++
										}
									l726:
										{
											position728, tokenIndex728, depth728 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l729
											}
											position++
											goto l728
										l729:
											position, tokenIndex, depth = position728, tokenIndex728, depth728
											if buffer[position] != rune('R') {
												goto l714
											}
											position++
										}
									l728:
										{
											position730, tokenIndex730, depth730 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l731
											}
											position++
											goto l730
										l731:
											position, tokenIndex, depth = position730, tokenIndex730, depth730
											if buffer[position] != rune('A') {
												goto l714
											}
											position++
										}
									l730:
										{
											position732, tokenIndex732, depth732 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l733
											}
											position++
											goto l732
										l733:
											position, tokenIndex, depth = position732, tokenIndex732, depth732
											if buffer[position] != rune('L') {
												goto l714
											}
											position++
										}
									l732:
										if !rules[ruleskip]() {
											goto l714
										}
										depth--
										add(ruleISLITERAL, position715)
									}
									goto l521
								l714:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
									{
										switch buffer[position] {
										case 'I', 'i':
											{
												position735 := position
												depth++
												{
													position736, tokenIndex736, depth736 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l737
													}
													position++
													goto l736
												l737:
													position, tokenIndex, depth = position736, tokenIndex736, depth736
													if buffer[position] != rune('I') {
														goto l520
													}
													position++
												}
											l736:
												{
													position738, tokenIndex738, depth738 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l739
													}
													position++
													goto l738
												l739:
													position, tokenIndex, depth = position738, tokenIndex738, depth738
													if buffer[position] != rune('S') {
														goto l520
													}
													position++
												}
											l738:
												{
													position740, tokenIndex740, depth740 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l741
													}
													position++
													goto l740
												l741:
													position, tokenIndex, depth = position740, tokenIndex740, depth740
													if buffer[position] != rune('N') {
														goto l520
													}
													position++
												}
											l740:
												{
													position742, tokenIndex742, depth742 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l743
													}
													position++
													goto l742
												l743:
													position, tokenIndex, depth = position742, tokenIndex742, depth742
													if buffer[position] != rune('U') {
														goto l520
													}
													position++
												}
											l742:
												{
													position744, tokenIndex744, depth744 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l745
													}
													position++
													goto l744
												l745:
													position, tokenIndex, depth = position744, tokenIndex744, depth744
													if buffer[position] != rune('M') {
														goto l520
													}
													position++
												}
											l744:
												{
													position746, tokenIndex746, depth746 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l747
													}
													position++
													goto l746
												l747:
													position, tokenIndex, depth = position746, tokenIndex746, depth746
													if buffer[position] != rune('E') {
														goto l520
													}
													position++
												}
											l746:
												{
													position748, tokenIndex748, depth748 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l749
													}
													position++
													goto l748
												l749:
													position, tokenIndex, depth = position748, tokenIndex748, depth748
													if buffer[position] != rune('R') {
														goto l520
													}
													position++
												}
											l748:
												{
													position750, tokenIndex750, depth750 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l751
													}
													position++
													goto l750
												l751:
													position, tokenIndex, depth = position750, tokenIndex750, depth750
													if buffer[position] != rune('I') {
														goto l520
													}
													position++
												}
											l750:
												{
													position752, tokenIndex752, depth752 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l753
													}
													position++
													goto l752
												l753:
													position, tokenIndex, depth = position752, tokenIndex752, depth752
													if buffer[position] != rune('C') {
														goto l520
													}
													position++
												}
											l752:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleISNUMERIC, position735)
											}
											break
										case 'S', 's':
											{
												position754 := position
												depth++
												{
													position755, tokenIndex755, depth755 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l756
													}
													position++
													goto l755
												l756:
													position, tokenIndex, depth = position755, tokenIndex755, depth755
													if buffer[position] != rune('S') {
														goto l520
													}
													position++
												}
											l755:
												{
													position757, tokenIndex757, depth757 := position, tokenIndex, depth
													if buffer[position] != rune('h') {
														goto l758
													}
													position++
													goto l757
												l758:
													position, tokenIndex, depth = position757, tokenIndex757, depth757
													if buffer[position] != rune('H') {
														goto l520
													}
													position++
												}
											l757:
												{
													position759, tokenIndex759, depth759 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l760
													}
													position++
													goto l759
												l760:
													position, tokenIndex, depth = position759, tokenIndex759, depth759
													if buffer[position] != rune('A') {
														goto l520
													}
													position++
												}
											l759:
												if buffer[position] != rune('5') {
													goto l520
												}
												position++
												if buffer[position] != rune('1') {
													goto l520
												}
												position++
												if buffer[position] != rune('2') {
													goto l520
												}
												position++
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleSHA512, position754)
											}
											break
										case 'M', 'm':
											{
												position761 := position
												depth++
												{
													position762, tokenIndex762, depth762 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l763
													}
													position++
													goto l762
												l763:
													position, tokenIndex, depth = position762, tokenIndex762, depth762
													if buffer[position] != rune('M') {
														goto l520
													}
													position++
												}
											l762:
												{
													position764, tokenIndex764, depth764 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l765
													}
													position++
													goto l764
												l765:
													position, tokenIndex, depth = position764, tokenIndex764, depth764
													if buffer[position] != rune('D') {
														goto l520
													}
													position++
												}
											l764:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleMD5, position761)
											}
											break
										case 'T', 't':
											{
												position766 := position
												depth++
												{
													position767, tokenIndex767, depth767 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l768
													}
													position++
													goto l767
												l768:
													position, tokenIndex, depth = position767, tokenIndex767, depth767
													if buffer[position] != rune('T') {
														goto l520
													}
													position++
												}
											l767:
												{
													position769, tokenIndex769, depth769 := position, tokenIndex, depth
													if buffer[position] != rune('z') {
														goto l770
													}
													position++
													goto l769
												l770:
													position, tokenIndex, depth = position769, tokenIndex769, depth769
													if buffer[position] != rune('Z') {
														goto l520
													}
													position++
												}
											l769:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleTZ, position766)
											}
											break
										case 'H', 'h':
											{
												position771 := position
												depth++
												{
													position772, tokenIndex772, depth772 := position, tokenIndex, depth
													if buffer[position] != rune('h') {
														goto l773
													}
													position++
													goto l772
												l773:
													position, tokenIndex, depth = position772, tokenIndex772, depth772
													if buffer[position] != rune('H') {
														goto l520
													}
													position++
												}
											l772:
												{
													position774, tokenIndex774, depth774 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l775
													}
													position++
													goto l774
												l775:
													position, tokenIndex, depth = position774, tokenIndex774, depth774
													if buffer[position] != rune('O') {
														goto l520
													}
													position++
												}
											l774:
												{
													position776, tokenIndex776, depth776 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l777
													}
													position++
													goto l776
												l777:
													position, tokenIndex, depth = position776, tokenIndex776, depth776
													if buffer[position] != rune('U') {
														goto l520
													}
													position++
												}
											l776:
												{
													position778, tokenIndex778, depth778 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l779
													}
													position++
													goto l778
												l779:
													position, tokenIndex, depth = position778, tokenIndex778, depth778
													if buffer[position] != rune('R') {
														goto l520
													}
													position++
												}
											l778:
												{
													position780, tokenIndex780, depth780 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l781
													}
													position++
													goto l780
												l781:
													position, tokenIndex, depth = position780, tokenIndex780, depth780
													if buffer[position] != rune('S') {
														goto l520
													}
													position++
												}
											l780:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleHOURS, position771)
											}
											break
										case 'D', 'd':
											{
												position782 := position
												depth++
												{
													position783, tokenIndex783, depth783 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l784
													}
													position++
													goto l783
												l784:
													position, tokenIndex, depth = position783, tokenIndex783, depth783
													if buffer[position] != rune('D') {
														goto l520
													}
													position++
												}
											l783:
												{
													position785, tokenIndex785, depth785 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l786
													}
													position++
													goto l785
												l786:
													position, tokenIndex, depth = position785, tokenIndex785, depth785
													if buffer[position] != rune('A') {
														goto l520
													}
													position++
												}
											l785:
												{
													position787, tokenIndex787, depth787 := position, tokenIndex, depth
													if buffer[position] != rune('y') {
														goto l788
													}
													position++
													goto l787
												l788:
													position, tokenIndex, depth = position787, tokenIndex787, depth787
													if buffer[position] != rune('Y') {
														goto l520
													}
													position++
												}
											l787:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleDAY, position782)
											}
											break
										case 'Y', 'y':
											{
												position789 := position
												depth++
												{
													position790, tokenIndex790, depth790 := position, tokenIndex, depth
													if buffer[position] != rune('y') {
														goto l791
													}
													position++
													goto l790
												l791:
													position, tokenIndex, depth = position790, tokenIndex790, depth790
													if buffer[position] != rune('Y') {
														goto l520
													}
													position++
												}
											l790:
												{
													position792, tokenIndex792, depth792 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l793
													}
													position++
													goto l792
												l793:
													position, tokenIndex, depth = position792, tokenIndex792, depth792
													if buffer[position] != rune('E') {
														goto l520
													}
													position++
												}
											l792:
												{
													position794, tokenIndex794, depth794 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l795
													}
													position++
													goto l794
												l795:
													position, tokenIndex, depth = position794, tokenIndex794, depth794
													if buffer[position] != rune('A') {
														goto l520
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
														goto l520
													}
													position++
												}
											l796:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleYEAR, position789)
											}
											break
										case 'E', 'e':
											{
												position798 := position
												depth++
												{
													position799, tokenIndex799, depth799 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l800
													}
													position++
													goto l799
												l800:
													position, tokenIndex, depth = position799, tokenIndex799, depth799
													if buffer[position] != rune('E') {
														goto l520
													}
													position++
												}
											l799:
												{
													position801, tokenIndex801, depth801 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l802
													}
													position++
													goto l801
												l802:
													position, tokenIndex, depth = position801, tokenIndex801, depth801
													if buffer[position] != rune('N') {
														goto l520
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
														goto l520
													}
													position++
												}
											l803:
												{
													position805, tokenIndex805, depth805 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l806
													}
													position++
													goto l805
												l806:
													position, tokenIndex, depth = position805, tokenIndex805, depth805
													if buffer[position] != rune('O') {
														goto l520
													}
													position++
												}
											l805:
												{
													position807, tokenIndex807, depth807 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l808
													}
													position++
													goto l807
												l808:
													position, tokenIndex, depth = position807, tokenIndex807, depth807
													if buffer[position] != rune('D') {
														goto l520
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
														goto l520
													}
													position++
												}
											l809:
												if buffer[position] != rune('_') {
													goto l520
												}
												position++
												{
													position811, tokenIndex811, depth811 := position, tokenIndex, depth
													if buffer[position] != rune('f') {
														goto l812
													}
													position++
													goto l811
												l812:
													position, tokenIndex, depth = position811, tokenIndex811, depth811
													if buffer[position] != rune('F') {
														goto l520
													}
													position++
												}
											l811:
												{
													position813, tokenIndex813, depth813 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l814
													}
													position++
													goto l813
												l814:
													position, tokenIndex, depth = position813, tokenIndex813, depth813
													if buffer[position] != rune('O') {
														goto l520
													}
													position++
												}
											l813:
												{
													position815, tokenIndex815, depth815 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l816
													}
													position++
													goto l815
												l816:
													position, tokenIndex, depth = position815, tokenIndex815, depth815
													if buffer[position] != rune('R') {
														goto l520
													}
													position++
												}
											l815:
												if buffer[position] != rune('_') {
													goto l520
												}
												position++
												{
													position817, tokenIndex817, depth817 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l818
													}
													position++
													goto l817
												l818:
													position, tokenIndex, depth = position817, tokenIndex817, depth817
													if buffer[position] != rune('U') {
														goto l520
													}
													position++
												}
											l817:
												{
													position819, tokenIndex819, depth819 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l820
													}
													position++
													goto l819
												l820:
													position, tokenIndex, depth = position819, tokenIndex819, depth819
													if buffer[position] != rune('R') {
														goto l520
													}
													position++
												}
											l819:
												{
													position821, tokenIndex821, depth821 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l822
													}
													position++
													goto l821
												l822:
													position, tokenIndex, depth = position821, tokenIndex821, depth821
													if buffer[position] != rune('I') {
														goto l520
													}
													position++
												}
											l821:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleENCODEFORURI, position798)
											}
											break
										case 'L', 'l':
											{
												position823 := position
												depth++
												{
													position824, tokenIndex824, depth824 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l825
													}
													position++
													goto l824
												l825:
													position, tokenIndex, depth = position824, tokenIndex824, depth824
													if buffer[position] != rune('L') {
														goto l520
													}
													position++
												}
											l824:
												{
													position826, tokenIndex826, depth826 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l827
													}
													position++
													goto l826
												l827:
													position, tokenIndex, depth = position826, tokenIndex826, depth826
													if buffer[position] != rune('C') {
														goto l520
													}
													position++
												}
											l826:
												{
													position828, tokenIndex828, depth828 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l829
													}
													position++
													goto l828
												l829:
													position, tokenIndex, depth = position828, tokenIndex828, depth828
													if buffer[position] != rune('A') {
														goto l520
													}
													position++
												}
											l828:
												{
													position830, tokenIndex830, depth830 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l831
													}
													position++
													goto l830
												l831:
													position, tokenIndex, depth = position830, tokenIndex830, depth830
													if buffer[position] != rune('S') {
														goto l520
													}
													position++
												}
											l830:
												{
													position832, tokenIndex832, depth832 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l833
													}
													position++
													goto l832
												l833:
													position, tokenIndex, depth = position832, tokenIndex832, depth832
													if buffer[position] != rune('E') {
														goto l520
													}
													position++
												}
											l832:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleLCASE, position823)
											}
											break
										case 'U', 'u':
											{
												position834 := position
												depth++
												{
													position835, tokenIndex835, depth835 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l836
													}
													position++
													goto l835
												l836:
													position, tokenIndex, depth = position835, tokenIndex835, depth835
													if buffer[position] != rune('U') {
														goto l520
													}
													position++
												}
											l835:
												{
													position837, tokenIndex837, depth837 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l838
													}
													position++
													goto l837
												l838:
													position, tokenIndex, depth = position837, tokenIndex837, depth837
													if buffer[position] != rune('C') {
														goto l520
													}
													position++
												}
											l837:
												{
													position839, tokenIndex839, depth839 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l840
													}
													position++
													goto l839
												l840:
													position, tokenIndex, depth = position839, tokenIndex839, depth839
													if buffer[position] != rune('A') {
														goto l520
													}
													position++
												}
											l839:
												{
													position841, tokenIndex841, depth841 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l842
													}
													position++
													goto l841
												l842:
													position, tokenIndex, depth = position841, tokenIndex841, depth841
													if buffer[position] != rune('S') {
														goto l520
													}
													position++
												}
											l841:
												{
													position843, tokenIndex843, depth843 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l844
													}
													position++
													goto l843
												l844:
													position, tokenIndex, depth = position843, tokenIndex843, depth843
													if buffer[position] != rune('E') {
														goto l520
													}
													position++
												}
											l843:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleUCASE, position834)
											}
											break
										case 'F', 'f':
											{
												position845 := position
												depth++
												{
													position846, tokenIndex846, depth846 := position, tokenIndex, depth
													if buffer[position] != rune('f') {
														goto l847
													}
													position++
													goto l846
												l847:
													position, tokenIndex, depth = position846, tokenIndex846, depth846
													if buffer[position] != rune('F') {
														goto l520
													}
													position++
												}
											l846:
												{
													position848, tokenIndex848, depth848 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l849
													}
													position++
													goto l848
												l849:
													position, tokenIndex, depth = position848, tokenIndex848, depth848
													if buffer[position] != rune('L') {
														goto l520
													}
													position++
												}
											l848:
												{
													position850, tokenIndex850, depth850 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l851
													}
													position++
													goto l850
												l851:
													position, tokenIndex, depth = position850, tokenIndex850, depth850
													if buffer[position] != rune('O') {
														goto l520
													}
													position++
												}
											l850:
												{
													position852, tokenIndex852, depth852 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l853
													}
													position++
													goto l852
												l853:
													position, tokenIndex, depth = position852, tokenIndex852, depth852
													if buffer[position] != rune('O') {
														goto l520
													}
													position++
												}
											l852:
												{
													position854, tokenIndex854, depth854 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l855
													}
													position++
													goto l854
												l855:
													position, tokenIndex, depth = position854, tokenIndex854, depth854
													if buffer[position] != rune('R') {
														goto l520
													}
													position++
												}
											l854:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleFLOOR, position845)
											}
											break
										case 'R', 'r':
											{
												position856 := position
												depth++
												{
													position857, tokenIndex857, depth857 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l858
													}
													position++
													goto l857
												l858:
													position, tokenIndex, depth = position857, tokenIndex857, depth857
													if buffer[position] != rune('R') {
														goto l520
													}
													position++
												}
											l857:
												{
													position859, tokenIndex859, depth859 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l860
													}
													position++
													goto l859
												l860:
													position, tokenIndex, depth = position859, tokenIndex859, depth859
													if buffer[position] != rune('O') {
														goto l520
													}
													position++
												}
											l859:
												{
													position861, tokenIndex861, depth861 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l862
													}
													position++
													goto l861
												l862:
													position, tokenIndex, depth = position861, tokenIndex861, depth861
													if buffer[position] != rune('U') {
														goto l520
													}
													position++
												}
											l861:
												{
													position863, tokenIndex863, depth863 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l864
													}
													position++
													goto l863
												l864:
													position, tokenIndex, depth = position863, tokenIndex863, depth863
													if buffer[position] != rune('N') {
														goto l520
													}
													position++
												}
											l863:
												{
													position865, tokenIndex865, depth865 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l866
													}
													position++
													goto l865
												l866:
													position, tokenIndex, depth = position865, tokenIndex865, depth865
													if buffer[position] != rune('D') {
														goto l520
													}
													position++
												}
											l865:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleROUND, position856)
											}
											break
										case 'C', 'c':
											{
												position867 := position
												depth++
												{
													position868, tokenIndex868, depth868 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l869
													}
													position++
													goto l868
												l869:
													position, tokenIndex, depth = position868, tokenIndex868, depth868
													if buffer[position] != rune('C') {
														goto l520
													}
													position++
												}
											l868:
												{
													position870, tokenIndex870, depth870 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l871
													}
													position++
													goto l870
												l871:
													position, tokenIndex, depth = position870, tokenIndex870, depth870
													if buffer[position] != rune('E') {
														goto l520
													}
													position++
												}
											l870:
												{
													position872, tokenIndex872, depth872 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l873
													}
													position++
													goto l872
												l873:
													position, tokenIndex, depth = position872, tokenIndex872, depth872
													if buffer[position] != rune('I') {
														goto l520
													}
													position++
												}
											l872:
												{
													position874, tokenIndex874, depth874 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l875
													}
													position++
													goto l874
												l875:
													position, tokenIndex, depth = position874, tokenIndex874, depth874
													if buffer[position] != rune('L') {
														goto l520
													}
													position++
												}
											l874:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleCEIL, position867)
											}
											break
										default:
											{
												position876 := position
												depth++
												{
													position877, tokenIndex877, depth877 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l878
													}
													position++
													goto l877
												l878:
													position, tokenIndex, depth = position877, tokenIndex877, depth877
													if buffer[position] != rune('A') {
														goto l520
													}
													position++
												}
											l877:
												{
													position879, tokenIndex879, depth879 := position, tokenIndex, depth
													if buffer[position] != rune('b') {
														goto l880
													}
													position++
													goto l879
												l880:
													position, tokenIndex, depth = position879, tokenIndex879, depth879
													if buffer[position] != rune('B') {
														goto l520
													}
													position++
												}
											l879:
												{
													position881, tokenIndex881, depth881 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l882
													}
													position++
													goto l881
												l882:
													position, tokenIndex, depth = position881, tokenIndex881, depth881
													if buffer[position] != rune('S') {
														goto l520
													}
													position++
												}
											l881:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleABS, position876)
											}
											break
										}
									}

								}
							l521:
								if !rules[ruleLPAREN]() {
									goto l520
								}
								if !rules[ruleexpression]() {
									goto l520
								}
								if !rules[ruleRPAREN]() {
									goto l520
								}
								goto l519
							l520:
								position, tokenIndex, depth = position519, tokenIndex519, depth519
								{
									position884, tokenIndex884, depth884 := position, tokenIndex, depth
									{
										position886 := position
										depth++
										{
											position887, tokenIndex887, depth887 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l888
											}
											position++
											goto l887
										l888:
											position, tokenIndex, depth = position887, tokenIndex887, depth887
											if buffer[position] != rune('S') {
												goto l885
											}
											position++
										}
									l887:
										{
											position889, tokenIndex889, depth889 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l890
											}
											position++
											goto l889
										l890:
											position, tokenIndex, depth = position889, tokenIndex889, depth889
											if buffer[position] != rune('T') {
												goto l885
											}
											position++
										}
									l889:
										{
											position891, tokenIndex891, depth891 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l892
											}
											position++
											goto l891
										l892:
											position, tokenIndex, depth = position891, tokenIndex891, depth891
											if buffer[position] != rune('R') {
												goto l885
											}
											position++
										}
									l891:
										{
											position893, tokenIndex893, depth893 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l894
											}
											position++
											goto l893
										l894:
											position, tokenIndex, depth = position893, tokenIndex893, depth893
											if buffer[position] != rune('S') {
												goto l885
											}
											position++
										}
									l893:
										{
											position895, tokenIndex895, depth895 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l896
											}
											position++
											goto l895
										l896:
											position, tokenIndex, depth = position895, tokenIndex895, depth895
											if buffer[position] != rune('T') {
												goto l885
											}
											position++
										}
									l895:
										{
											position897, tokenIndex897, depth897 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l898
											}
											position++
											goto l897
										l898:
											position, tokenIndex, depth = position897, tokenIndex897, depth897
											if buffer[position] != rune('A') {
												goto l885
											}
											position++
										}
									l897:
										{
											position899, tokenIndex899, depth899 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l900
											}
											position++
											goto l899
										l900:
											position, tokenIndex, depth = position899, tokenIndex899, depth899
											if buffer[position] != rune('R') {
												goto l885
											}
											position++
										}
									l899:
										{
											position901, tokenIndex901, depth901 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l902
											}
											position++
											goto l901
										l902:
											position, tokenIndex, depth = position901, tokenIndex901, depth901
											if buffer[position] != rune('T') {
												goto l885
											}
											position++
										}
									l901:
										{
											position903, tokenIndex903, depth903 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l904
											}
											position++
											goto l903
										l904:
											position, tokenIndex, depth = position903, tokenIndex903, depth903
											if buffer[position] != rune('S') {
												goto l885
											}
											position++
										}
									l903:
										if !rules[ruleskip]() {
											goto l885
										}
										depth--
										add(ruleSTRSTARTS, position886)
									}
									goto l884
								l885:
									position, tokenIndex, depth = position884, tokenIndex884, depth884
									{
										position906 := position
										depth++
										{
											position907, tokenIndex907, depth907 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l908
											}
											position++
											goto l907
										l908:
											position, tokenIndex, depth = position907, tokenIndex907, depth907
											if buffer[position] != rune('S') {
												goto l905
											}
											position++
										}
									l907:
										{
											position909, tokenIndex909, depth909 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l910
											}
											position++
											goto l909
										l910:
											position, tokenIndex, depth = position909, tokenIndex909, depth909
											if buffer[position] != rune('T') {
												goto l905
											}
											position++
										}
									l909:
										{
											position911, tokenIndex911, depth911 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l912
											}
											position++
											goto l911
										l912:
											position, tokenIndex, depth = position911, tokenIndex911, depth911
											if buffer[position] != rune('R') {
												goto l905
											}
											position++
										}
									l911:
										{
											position913, tokenIndex913, depth913 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l914
											}
											position++
											goto l913
										l914:
											position, tokenIndex, depth = position913, tokenIndex913, depth913
											if buffer[position] != rune('E') {
												goto l905
											}
											position++
										}
									l913:
										{
											position915, tokenIndex915, depth915 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l916
											}
											position++
											goto l915
										l916:
											position, tokenIndex, depth = position915, tokenIndex915, depth915
											if buffer[position] != rune('N') {
												goto l905
											}
											position++
										}
									l915:
										{
											position917, tokenIndex917, depth917 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l918
											}
											position++
											goto l917
										l918:
											position, tokenIndex, depth = position917, tokenIndex917, depth917
											if buffer[position] != rune('D') {
												goto l905
											}
											position++
										}
									l917:
										{
											position919, tokenIndex919, depth919 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l920
											}
											position++
											goto l919
										l920:
											position, tokenIndex, depth = position919, tokenIndex919, depth919
											if buffer[position] != rune('S') {
												goto l905
											}
											position++
										}
									l919:
										if !rules[ruleskip]() {
											goto l905
										}
										depth--
										add(ruleSTRENDS, position906)
									}
									goto l884
								l905:
									position, tokenIndex, depth = position884, tokenIndex884, depth884
									{
										position922 := position
										depth++
										{
											position923, tokenIndex923, depth923 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l924
											}
											position++
											goto l923
										l924:
											position, tokenIndex, depth = position923, tokenIndex923, depth923
											if buffer[position] != rune('S') {
												goto l921
											}
											position++
										}
									l923:
										{
											position925, tokenIndex925, depth925 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l926
											}
											position++
											goto l925
										l926:
											position, tokenIndex, depth = position925, tokenIndex925, depth925
											if buffer[position] != rune('T') {
												goto l921
											}
											position++
										}
									l925:
										{
											position927, tokenIndex927, depth927 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l928
											}
											position++
											goto l927
										l928:
											position, tokenIndex, depth = position927, tokenIndex927, depth927
											if buffer[position] != rune('R') {
												goto l921
											}
											position++
										}
									l927:
										{
											position929, tokenIndex929, depth929 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l930
											}
											position++
											goto l929
										l930:
											position, tokenIndex, depth = position929, tokenIndex929, depth929
											if buffer[position] != rune('B') {
												goto l921
											}
											position++
										}
									l929:
										{
											position931, tokenIndex931, depth931 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l932
											}
											position++
											goto l931
										l932:
											position, tokenIndex, depth = position931, tokenIndex931, depth931
											if buffer[position] != rune('E') {
												goto l921
											}
											position++
										}
									l931:
										{
											position933, tokenIndex933, depth933 := position, tokenIndex, depth
											if buffer[position] != rune('f') {
												goto l934
											}
											position++
											goto l933
										l934:
											position, tokenIndex, depth = position933, tokenIndex933, depth933
											if buffer[position] != rune('F') {
												goto l921
											}
											position++
										}
									l933:
										{
											position935, tokenIndex935, depth935 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l936
											}
											position++
											goto l935
										l936:
											position, tokenIndex, depth = position935, tokenIndex935, depth935
											if buffer[position] != rune('O') {
												goto l921
											}
											position++
										}
									l935:
										{
											position937, tokenIndex937, depth937 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l938
											}
											position++
											goto l937
										l938:
											position, tokenIndex, depth = position937, tokenIndex937, depth937
											if buffer[position] != rune('R') {
												goto l921
											}
											position++
										}
									l937:
										{
											position939, tokenIndex939, depth939 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l940
											}
											position++
											goto l939
										l940:
											position, tokenIndex, depth = position939, tokenIndex939, depth939
											if buffer[position] != rune('E') {
												goto l921
											}
											position++
										}
									l939:
										if !rules[ruleskip]() {
											goto l921
										}
										depth--
										add(ruleSTRBEFORE, position922)
									}
									goto l884
								l921:
									position, tokenIndex, depth = position884, tokenIndex884, depth884
									{
										position942 := position
										depth++
										{
											position943, tokenIndex943, depth943 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l944
											}
											position++
											goto l943
										l944:
											position, tokenIndex, depth = position943, tokenIndex943, depth943
											if buffer[position] != rune('S') {
												goto l941
											}
											position++
										}
									l943:
										{
											position945, tokenIndex945, depth945 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l946
											}
											position++
											goto l945
										l946:
											position, tokenIndex, depth = position945, tokenIndex945, depth945
											if buffer[position] != rune('T') {
												goto l941
											}
											position++
										}
									l945:
										{
											position947, tokenIndex947, depth947 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l948
											}
											position++
											goto l947
										l948:
											position, tokenIndex, depth = position947, tokenIndex947, depth947
											if buffer[position] != rune('R') {
												goto l941
											}
											position++
										}
									l947:
										{
											position949, tokenIndex949, depth949 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l950
											}
											position++
											goto l949
										l950:
											position, tokenIndex, depth = position949, tokenIndex949, depth949
											if buffer[position] != rune('A') {
												goto l941
											}
											position++
										}
									l949:
										{
											position951, tokenIndex951, depth951 := position, tokenIndex, depth
											if buffer[position] != rune('f') {
												goto l952
											}
											position++
											goto l951
										l952:
											position, tokenIndex, depth = position951, tokenIndex951, depth951
											if buffer[position] != rune('F') {
												goto l941
											}
											position++
										}
									l951:
										{
											position953, tokenIndex953, depth953 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l954
											}
											position++
											goto l953
										l954:
											position, tokenIndex, depth = position953, tokenIndex953, depth953
											if buffer[position] != rune('T') {
												goto l941
											}
											position++
										}
									l953:
										{
											position955, tokenIndex955, depth955 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l956
											}
											position++
											goto l955
										l956:
											position, tokenIndex, depth = position955, tokenIndex955, depth955
											if buffer[position] != rune('E') {
												goto l941
											}
											position++
										}
									l955:
										{
											position957, tokenIndex957, depth957 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l958
											}
											position++
											goto l957
										l958:
											position, tokenIndex, depth = position957, tokenIndex957, depth957
											if buffer[position] != rune('R') {
												goto l941
											}
											position++
										}
									l957:
										if !rules[ruleskip]() {
											goto l941
										}
										depth--
										add(ruleSTRAFTER, position942)
									}
									goto l884
								l941:
									position, tokenIndex, depth = position884, tokenIndex884, depth884
									{
										position960 := position
										depth++
										{
											position961, tokenIndex961, depth961 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l962
											}
											position++
											goto l961
										l962:
											position, tokenIndex, depth = position961, tokenIndex961, depth961
											if buffer[position] != rune('S') {
												goto l959
											}
											position++
										}
									l961:
										{
											position963, tokenIndex963, depth963 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l964
											}
											position++
											goto l963
										l964:
											position, tokenIndex, depth = position963, tokenIndex963, depth963
											if buffer[position] != rune('T') {
												goto l959
											}
											position++
										}
									l963:
										{
											position965, tokenIndex965, depth965 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l966
											}
											position++
											goto l965
										l966:
											position, tokenIndex, depth = position965, tokenIndex965, depth965
											if buffer[position] != rune('R') {
												goto l959
											}
											position++
										}
									l965:
										{
											position967, tokenIndex967, depth967 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l968
											}
											position++
											goto l967
										l968:
											position, tokenIndex, depth = position967, tokenIndex967, depth967
											if buffer[position] != rune('L') {
												goto l959
											}
											position++
										}
									l967:
										{
											position969, tokenIndex969, depth969 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l970
											}
											position++
											goto l969
										l970:
											position, tokenIndex, depth = position969, tokenIndex969, depth969
											if buffer[position] != rune('A') {
												goto l959
											}
											position++
										}
									l969:
										{
											position971, tokenIndex971, depth971 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l972
											}
											position++
											goto l971
										l972:
											position, tokenIndex, depth = position971, tokenIndex971, depth971
											if buffer[position] != rune('N') {
												goto l959
											}
											position++
										}
									l971:
										{
											position973, tokenIndex973, depth973 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l974
											}
											position++
											goto l973
										l974:
											position, tokenIndex, depth = position973, tokenIndex973, depth973
											if buffer[position] != rune('G') {
												goto l959
											}
											position++
										}
									l973:
										if !rules[ruleskip]() {
											goto l959
										}
										depth--
										add(ruleSTRLANG, position960)
									}
									goto l884
								l959:
									position, tokenIndex, depth = position884, tokenIndex884, depth884
									{
										position976 := position
										depth++
										{
											position977, tokenIndex977, depth977 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l978
											}
											position++
											goto l977
										l978:
											position, tokenIndex, depth = position977, tokenIndex977, depth977
											if buffer[position] != rune('S') {
												goto l975
											}
											position++
										}
									l977:
										{
											position979, tokenIndex979, depth979 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l980
											}
											position++
											goto l979
										l980:
											position, tokenIndex, depth = position979, tokenIndex979, depth979
											if buffer[position] != rune('T') {
												goto l975
											}
											position++
										}
									l979:
										{
											position981, tokenIndex981, depth981 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l982
											}
											position++
											goto l981
										l982:
											position, tokenIndex, depth = position981, tokenIndex981, depth981
											if buffer[position] != rune('R') {
												goto l975
											}
											position++
										}
									l981:
										{
											position983, tokenIndex983, depth983 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l984
											}
											position++
											goto l983
										l984:
											position, tokenIndex, depth = position983, tokenIndex983, depth983
											if buffer[position] != rune('D') {
												goto l975
											}
											position++
										}
									l983:
										{
											position985, tokenIndex985, depth985 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l986
											}
											position++
											goto l985
										l986:
											position, tokenIndex, depth = position985, tokenIndex985, depth985
											if buffer[position] != rune('T') {
												goto l975
											}
											position++
										}
									l985:
										if !rules[ruleskip]() {
											goto l975
										}
										depth--
										add(ruleSTRDT, position976)
									}
									goto l884
								l975:
									position, tokenIndex, depth = position884, tokenIndex884, depth884
									{
										switch buffer[position] {
										case 'S', 's':
											{
												position988 := position
												depth++
												{
													position989, tokenIndex989, depth989 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l990
													}
													position++
													goto l989
												l990:
													position, tokenIndex, depth = position989, tokenIndex989, depth989
													if buffer[position] != rune('S') {
														goto l883
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
														goto l883
													}
													position++
												}
											l991:
												{
													position993, tokenIndex993, depth993 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l994
													}
													position++
													goto l993
												l994:
													position, tokenIndex, depth = position993, tokenIndex993, depth993
													if buffer[position] != rune('M') {
														goto l883
													}
													position++
												}
											l993:
												{
													position995, tokenIndex995, depth995 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l996
													}
													position++
													goto l995
												l996:
													position, tokenIndex, depth = position995, tokenIndex995, depth995
													if buffer[position] != rune('E') {
														goto l883
													}
													position++
												}
											l995:
												{
													position997, tokenIndex997, depth997 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l998
													}
													position++
													goto l997
												l998:
													position, tokenIndex, depth = position997, tokenIndex997, depth997
													if buffer[position] != rune('T') {
														goto l883
													}
													position++
												}
											l997:
												{
													position999, tokenIndex999, depth999 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1000
													}
													position++
													goto l999
												l1000:
													position, tokenIndex, depth = position999, tokenIndex999, depth999
													if buffer[position] != rune('E') {
														goto l883
													}
													position++
												}
											l999:
												{
													position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l1002
													}
													position++
													goto l1001
												l1002:
													position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
													if buffer[position] != rune('R') {
														goto l883
													}
													position++
												}
											l1001:
												{
													position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l1004
													}
													position++
													goto l1003
												l1004:
													position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
													if buffer[position] != rune('M') {
														goto l883
													}
													position++
												}
											l1003:
												if !rules[ruleskip]() {
													goto l883
												}
												depth--
												add(ruleSAMETERM, position988)
											}
											break
										case 'C', 'c':
											{
												position1005 := position
												depth++
												{
													position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1007
													}
													position++
													goto l1006
												l1007:
													position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
													if buffer[position] != rune('C') {
														goto l883
													}
													position++
												}
											l1006:
												{
													position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l1009
													}
													position++
													goto l1008
												l1009:
													position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
													if buffer[position] != rune('O') {
														goto l883
													}
													position++
												}
											l1008:
												{
													position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1011
													}
													position++
													goto l1010
												l1011:
													position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
													if buffer[position] != rune('N') {
														goto l883
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
														goto l883
													}
													position++
												}
											l1012:
												{
													position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1015
													}
													position++
													goto l1014
												l1015:
													position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
													if buffer[position] != rune('A') {
														goto l883
													}
													position++
												}
											l1014:
												{
													position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l1017
													}
													position++
													goto l1016
												l1017:
													position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
													if buffer[position] != rune('I') {
														goto l883
													}
													position++
												}
											l1016:
												{
													position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1019
													}
													position++
													goto l1018
												l1019:
													position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
													if buffer[position] != rune('N') {
														goto l883
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
														goto l883
													}
													position++
												}
											l1020:
												if !rules[ruleskip]() {
													goto l883
												}
												depth--
												add(ruleCONTAINS, position1005)
											}
											break
										default:
											{
												position1022 := position
												depth++
												{
													position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l1024
													}
													position++
													goto l1023
												l1024:
													position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
													if buffer[position] != rune('L') {
														goto l883
													}
													position++
												}
											l1023:
												{
													position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1026
													}
													position++
													goto l1025
												l1026:
													position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
													if buffer[position] != rune('A') {
														goto l883
													}
													position++
												}
											l1025:
												{
													position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1028
													}
													position++
													goto l1027
												l1028:
													position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
													if buffer[position] != rune('N') {
														goto l883
													}
													position++
												}
											l1027:
												{
													position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
													if buffer[position] != rune('g') {
														goto l1030
													}
													position++
													goto l1029
												l1030:
													position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
													if buffer[position] != rune('G') {
														goto l883
													}
													position++
												}
											l1029:
												{
													position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
													if buffer[position] != rune('m') {
														goto l1032
													}
													position++
													goto l1031
												l1032:
													position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
													if buffer[position] != rune('M') {
														goto l883
													}
													position++
												}
											l1031:
												{
													position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1034
													}
													position++
													goto l1033
												l1034:
													position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
													if buffer[position] != rune('A') {
														goto l883
													}
													position++
												}
											l1033:
												{
													position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1036
													}
													position++
													goto l1035
												l1036:
													position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
													if buffer[position] != rune('T') {
														goto l883
													}
													position++
												}
											l1035:
												{
													position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1038
													}
													position++
													goto l1037
												l1038:
													position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
													if buffer[position] != rune('C') {
														goto l883
													}
													position++
												}
											l1037:
												{
													position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
													if buffer[position] != rune('h') {
														goto l1040
													}
													position++
													goto l1039
												l1040:
													position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
													if buffer[position] != rune('H') {
														goto l883
													}
													position++
												}
											l1039:
												{
													position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1042
													}
													position++
													goto l1041
												l1042:
													position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
													if buffer[position] != rune('E') {
														goto l883
													}
													position++
												}
											l1041:
												{
													position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1044
													}
													position++
													goto l1043
												l1044:
													position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
													if buffer[position] != rune('S') {
														goto l883
													}
													position++
												}
											l1043:
												if !rules[ruleskip]() {
													goto l883
												}
												depth--
												add(ruleLANGMATCHES, position1022)
											}
											break
										}
									}

								}
							l884:
								if !rules[ruleLPAREN]() {
									goto l883
								}
								if !rules[ruleexpression]() {
									goto l883
								}
								if !rules[ruleCOMMA]() {
									goto l883
								}
								if !rules[ruleexpression]() {
									goto l883
								}
								if !rules[ruleRPAREN]() {
									goto l883
								}
								goto l519
							l883:
								position, tokenIndex, depth = position519, tokenIndex519, depth519
								{
									position1046 := position
									depth++
									{
										position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1048
										}
										position++
										goto l1047
									l1048:
										position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
										if buffer[position] != rune('B') {
											goto l1045
										}
										position++
									}
								l1047:
									{
										position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1050
										}
										position++
										goto l1049
									l1050:
										position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
										if buffer[position] != rune('O') {
											goto l1045
										}
										position++
									}
								l1049:
									{
										position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1052
										}
										position++
										goto l1051
									l1052:
										position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
										if buffer[position] != rune('U') {
											goto l1045
										}
										position++
									}
								l1051:
									{
										position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1054
										}
										position++
										goto l1053
									l1054:
										position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
										if buffer[position] != rune('N') {
											goto l1045
										}
										position++
									}
								l1053:
									{
										position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1056
										}
										position++
										goto l1055
									l1056:
										position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
										if buffer[position] != rune('D') {
											goto l1045
										}
										position++
									}
								l1055:
									if !rules[ruleskip]() {
										goto l1045
									}
									depth--
									add(ruleBOUND, position1046)
								}
								if !rules[ruleLPAREN]() {
									goto l1045
								}
								if !rules[rulevar]() {
									goto l1045
								}
								if !rules[ruleRPAREN]() {
									goto l1045
								}
								goto l519
							l1045:
								position, tokenIndex, depth = position519, tokenIndex519, depth519
								{
									switch buffer[position] {
									case 'S', 's':
										{
											position1059 := position
											depth++
											{
												position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l1061
												}
												position++
												goto l1060
											l1061:
												position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
												if buffer[position] != rune('S') {
													goto l1057
												}
												position++
											}
										l1060:
											{
												position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
												if buffer[position] != rune('t') {
													goto l1063
												}
												position++
												goto l1062
											l1063:
												position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
												if buffer[position] != rune('T') {
													goto l1057
												}
												position++
											}
										l1062:
											{
												position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
												if buffer[position] != rune('r') {
													goto l1065
												}
												position++
												goto l1064
											l1065:
												position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
												if buffer[position] != rune('R') {
													goto l1057
												}
												position++
											}
										l1064:
											{
												position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1067
												}
												position++
												goto l1066
											l1067:
												position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
												if buffer[position] != rune('U') {
													goto l1057
												}
												position++
											}
										l1066:
											{
												position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1069
												}
												position++
												goto l1068
											l1069:
												position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
												if buffer[position] != rune('U') {
													goto l1057
												}
												position++
											}
										l1068:
											{
												position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
												if buffer[position] != rune('i') {
													goto l1071
												}
												position++
												goto l1070
											l1071:
												position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
												if buffer[position] != rune('I') {
													goto l1057
												}
												position++
											}
										l1070:
											{
												position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1073
												}
												position++
												goto l1072
											l1073:
												position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
												if buffer[position] != rune('D') {
													goto l1057
												}
												position++
											}
										l1072:
											if !rules[ruleskip]() {
												goto l1057
											}
											depth--
											add(ruleSTRUUID, position1059)
										}
										break
									case 'U', 'u':
										{
											position1074 := position
											depth++
											{
												position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1076
												}
												position++
												goto l1075
											l1076:
												position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
												if buffer[position] != rune('U') {
													goto l1057
												}
												position++
											}
										l1075:
											{
												position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
												if buffer[position] != rune('u') {
													goto l1078
												}
												position++
												goto l1077
											l1078:
												position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
												if buffer[position] != rune('U') {
													goto l1057
												}
												position++
											}
										l1077:
											{
												position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
												if buffer[position] != rune('i') {
													goto l1080
												}
												position++
												goto l1079
											l1080:
												position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
												if buffer[position] != rune('I') {
													goto l1057
												}
												position++
											}
										l1079:
											{
												position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1082
												}
												position++
												goto l1081
											l1082:
												position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
												if buffer[position] != rune('D') {
													goto l1057
												}
												position++
											}
										l1081:
											if !rules[ruleskip]() {
												goto l1057
											}
											depth--
											add(ruleUUID, position1074)
										}
										break
									case 'N', 'n':
										{
											position1083 := position
											depth++
											{
												position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
												if buffer[position] != rune('n') {
													goto l1085
												}
												position++
												goto l1084
											l1085:
												position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
												if buffer[position] != rune('N') {
													goto l1057
												}
												position++
											}
										l1084:
											{
												position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
												if buffer[position] != rune('o') {
													goto l1087
												}
												position++
												goto l1086
											l1087:
												position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
												if buffer[position] != rune('O') {
													goto l1057
												}
												position++
											}
										l1086:
											{
												position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
												if buffer[position] != rune('w') {
													goto l1089
												}
												position++
												goto l1088
											l1089:
												position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
												if buffer[position] != rune('W') {
													goto l1057
												}
												position++
											}
										l1088:
											if !rules[ruleskip]() {
												goto l1057
											}
											depth--
											add(ruleNOW, position1083)
										}
										break
									default:
										{
											position1090 := position
											depth++
											{
												position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
												if buffer[position] != rune('r') {
													goto l1092
												}
												position++
												goto l1091
											l1092:
												position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
												if buffer[position] != rune('R') {
													goto l1057
												}
												position++
											}
										l1091:
											{
												position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l1094
												}
												position++
												goto l1093
											l1094:
												position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
												if buffer[position] != rune('A') {
													goto l1057
												}
												position++
											}
										l1093:
											{
												position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
												if buffer[position] != rune('n') {
													goto l1096
												}
												position++
												goto l1095
											l1096:
												position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
												if buffer[position] != rune('N') {
													goto l1057
												}
												position++
											}
										l1095:
											{
												position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1098
												}
												position++
												goto l1097
											l1098:
												position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
												if buffer[position] != rune('D') {
													goto l1057
												}
												position++
											}
										l1097:
											if !rules[ruleskip]() {
												goto l1057
											}
											depth--
											add(ruleRAND, position1090)
										}
										break
									}
								}

								if !rules[rulenil]() {
									goto l1057
								}
								goto l519
							l1057:
								position, tokenIndex, depth = position519, tokenIndex519, depth519
								{
									switch buffer[position] {
									case 'E', 'N', 'e', 'n':
										{
											position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
											{
												position1102 := position
												depth++
												{
													position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1104
													}
													position++
													goto l1103
												l1104:
													position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
													if buffer[position] != rune('E') {
														goto l1101
													}
													position++
												}
											l1103:
												{
													position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
													if buffer[position] != rune('x') {
														goto l1106
													}
													position++
													goto l1105
												l1106:
													position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
													if buffer[position] != rune('X') {
														goto l1101
													}
													position++
												}
											l1105:
												{
													position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l1108
													}
													position++
													goto l1107
												l1108:
													position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
													if buffer[position] != rune('I') {
														goto l1101
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
														goto l1101
													}
													position++
												}
											l1109:
												{
													position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1112
													}
													position++
													goto l1111
												l1112:
													position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
													if buffer[position] != rune('T') {
														goto l1101
													}
													position++
												}
											l1111:
												{
													position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1114
													}
													position++
													goto l1113
												l1114:
													position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
													if buffer[position] != rune('S') {
														goto l1101
													}
													position++
												}
											l1113:
												if !rules[ruleskip]() {
													goto l1101
												}
												depth--
												add(ruleEXISTS, position1102)
											}
											goto l1100
										l1101:
											position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
											{
												position1115 := position
												depth++
												{
													position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1117
													}
													position++
													goto l1116
												l1117:
													position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
													if buffer[position] != rune('N') {
														goto l517
													}
													position++
												}
											l1116:
												{
													position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l1119
													}
													position++
													goto l1118
												l1119:
													position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
													if buffer[position] != rune('O') {
														goto l517
													}
													position++
												}
											l1118:
												{
													position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1121
													}
													position++
													goto l1120
												l1121:
													position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
													if buffer[position] != rune('T') {
														goto l517
													}
													position++
												}
											l1120:
												if buffer[position] != rune(' ') {
													goto l517
												}
												position++
												{
													position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1123
													}
													position++
													goto l1122
												l1123:
													position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
													if buffer[position] != rune('E') {
														goto l517
													}
													position++
												}
											l1122:
												{
													position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
													if buffer[position] != rune('x') {
														goto l1125
													}
													position++
													goto l1124
												l1125:
													position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
													if buffer[position] != rune('X') {
														goto l517
													}
													position++
												}
											l1124:
												{
													position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
													if buffer[position] != rune('i') {
														goto l1127
													}
													position++
													goto l1126
												l1127:
													position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
													if buffer[position] != rune('I') {
														goto l517
													}
													position++
												}
											l1126:
												{
													position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1129
													}
													position++
													goto l1128
												l1129:
													position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
													if buffer[position] != rune('S') {
														goto l517
													}
													position++
												}
											l1128:
												{
													position1130, tokenIndex1130, depth1130 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1131
													}
													position++
													goto l1130
												l1131:
													position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
													if buffer[position] != rune('T') {
														goto l517
													}
													position++
												}
											l1130:
												{
													position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1133
													}
													position++
													goto l1132
												l1133:
													position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
													if buffer[position] != rune('S') {
														goto l517
													}
													position++
												}
											l1132:
												if !rules[ruleskip]() {
													goto l517
												}
												depth--
												add(ruleNOTEXIST, position1115)
											}
										}
									l1100:
										if !rules[rulegroupGraphPattern]() {
											goto l517
										}
										break
									case 'I', 'i':
										{
											position1134 := position
											depth++
											{
												position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
												if buffer[position] != rune('i') {
													goto l1136
												}
												position++
												goto l1135
											l1136:
												position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
												if buffer[position] != rune('I') {
													goto l517
												}
												position++
											}
										l1135:
											{
												position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
												if buffer[position] != rune('f') {
													goto l1138
												}
												position++
												goto l1137
											l1138:
												position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
												if buffer[position] != rune('F') {
													goto l517
												}
												position++
											}
										l1137:
											if !rules[ruleskip]() {
												goto l517
											}
											depth--
											add(ruleIF, position1134)
										}
										if !rules[ruleLPAREN]() {
											goto l517
										}
										if !rules[ruleexpression]() {
											goto l517
										}
										if !rules[ruleCOMMA]() {
											goto l517
										}
										if !rules[ruleexpression]() {
											goto l517
										}
										if !rules[ruleCOMMA]() {
											goto l517
										}
										if !rules[ruleexpression]() {
											goto l517
										}
										if !rules[ruleRPAREN]() {
											goto l517
										}
										break
									case 'C', 'c':
										{
											position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
											{
												position1141 := position
												depth++
												{
													position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1143
													}
													position++
													goto l1142
												l1143:
													position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
													if buffer[position] != rune('C') {
														goto l1140
													}
													position++
												}
											l1142:
												{
													position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l1145
													}
													position++
													goto l1144
												l1145:
													position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
													if buffer[position] != rune('O') {
														goto l1140
													}
													position++
												}
											l1144:
												{
													position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
													if buffer[position] != rune('n') {
														goto l1147
													}
													position++
													goto l1146
												l1147:
													position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
													if buffer[position] != rune('N') {
														goto l1140
													}
													position++
												}
											l1146:
												{
													position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1149
													}
													position++
													goto l1148
												l1149:
													position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
													if buffer[position] != rune('C') {
														goto l1140
													}
													position++
												}
											l1148:
												{
													position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1151
													}
													position++
													goto l1150
												l1151:
													position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
													if buffer[position] != rune('A') {
														goto l1140
													}
													position++
												}
											l1150:
												{
													position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1153
													}
													position++
													goto l1152
												l1153:
													position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
													if buffer[position] != rune('T') {
														goto l1140
													}
													position++
												}
											l1152:
												if !rules[ruleskip]() {
													goto l1140
												}
												depth--
												add(ruleCONCAT, position1141)
											}
											goto l1139
										l1140:
											position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
											{
												position1154 := position
												depth++
												{
													position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1156
													}
													position++
													goto l1155
												l1156:
													position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
													if buffer[position] != rune('C') {
														goto l517
													}
													position++
												}
											l1155:
												{
													position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
													if buffer[position] != rune('o') {
														goto l1158
													}
													position++
													goto l1157
												l1158:
													position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
													if buffer[position] != rune('O') {
														goto l517
													}
													position++
												}
											l1157:
												{
													position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1160
													}
													position++
													goto l1159
												l1160:
													position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
													if buffer[position] != rune('A') {
														goto l517
													}
													position++
												}
											l1159:
												{
													position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l1162
													}
													position++
													goto l1161
												l1162:
													position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
													if buffer[position] != rune('L') {
														goto l517
													}
													position++
												}
											l1161:
												{
													position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1164
													}
													position++
													goto l1163
												l1164:
													position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
													if buffer[position] != rune('E') {
														goto l517
													}
													position++
												}
											l1163:
												{
													position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1166
													}
													position++
													goto l1165
												l1166:
													position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
													if buffer[position] != rune('S') {
														goto l517
													}
													position++
												}
											l1165:
												{
													position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1168
													}
													position++
													goto l1167
												l1168:
													position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
													if buffer[position] != rune('C') {
														goto l517
													}
													position++
												}
											l1167:
												{
													position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1170
													}
													position++
													goto l1169
												l1170:
													position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
													if buffer[position] != rune('E') {
														goto l517
													}
													position++
												}
											l1169:
												if !rules[ruleskip]() {
													goto l517
												}
												depth--
												add(ruleCOALESCE, position1154)
											}
										}
									l1139:
										if !rules[ruleargList]() {
											goto l517
										}
										break
									case 'B', 'b':
										{
											position1171 := position
											depth++
											{
												position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
												if buffer[position] != rune('b') {
													goto l1173
												}
												position++
												goto l1172
											l1173:
												position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
												if buffer[position] != rune('B') {
													goto l517
												}
												position++
											}
										l1172:
											{
												position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
												if buffer[position] != rune('n') {
													goto l1175
												}
												position++
												goto l1174
											l1175:
												position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
												if buffer[position] != rune('N') {
													goto l517
												}
												position++
											}
										l1174:
											{
												position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
												if buffer[position] != rune('o') {
													goto l1177
												}
												position++
												goto l1176
											l1177:
												position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
												if buffer[position] != rune('O') {
													goto l517
												}
												position++
											}
										l1176:
											{
												position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l1179
												}
												position++
												goto l1178
											l1179:
												position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
												if buffer[position] != rune('D') {
													goto l517
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
													goto l517
												}
												position++
											}
										l1180:
											if !rules[ruleskip]() {
												goto l517
											}
											depth--
											add(ruleBNODE, position1171)
										}
										{
											position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
											if !rules[ruleLPAREN]() {
												goto l1183
											}
											if !rules[ruleexpression]() {
												goto l1183
											}
											if !rules[ruleRPAREN]() {
												goto l1183
											}
											goto l1182
										l1183:
											position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
											if !rules[rulenil]() {
												goto l517
											}
										}
									l1182:
										break
									default:
										{
											position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
											{
												position1186 := position
												depth++
												{
													position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1188
													}
													position++
													goto l1187
												l1188:
													position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
													if buffer[position] != rune('S') {
														goto l1185
													}
													position++
												}
											l1187:
												{
													position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
													if buffer[position] != rune('u') {
														goto l1190
													}
													position++
													goto l1189
												l1190:
													position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
													if buffer[position] != rune('U') {
														goto l1185
													}
													position++
												}
											l1189:
												{
													position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
													if buffer[position] != rune('b') {
														goto l1192
													}
													position++
													goto l1191
												l1192:
													position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
													if buffer[position] != rune('B') {
														goto l1185
													}
													position++
												}
											l1191:
												{
													position1193, tokenIndex1193, depth1193 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l1194
													}
													position++
													goto l1193
												l1194:
													position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
													if buffer[position] != rune('S') {
														goto l1185
													}
													position++
												}
											l1193:
												{
													position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
													if buffer[position] != rune('t') {
														goto l1196
													}
													position++
													goto l1195
												l1196:
													position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
													if buffer[position] != rune('T') {
														goto l1185
													}
													position++
												}
											l1195:
												{
													position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l1198
													}
													position++
													goto l1197
												l1198:
													position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
													if buffer[position] != rune('R') {
														goto l1185
													}
													position++
												}
											l1197:
												if !rules[ruleskip]() {
													goto l1185
												}
												depth--
												add(ruleSUBSTR, position1186)
											}
											goto l1184
										l1185:
											position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
											{
												position1200 := position
												depth++
												{
													position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l1202
													}
													position++
													goto l1201
												l1202:
													position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
													if buffer[position] != rune('R') {
														goto l1199
													}
													position++
												}
											l1201:
												{
													position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1204
													}
													position++
													goto l1203
												l1204:
													position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
													if buffer[position] != rune('E') {
														goto l1199
													}
													position++
												}
											l1203:
												{
													position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
													if buffer[position] != rune('p') {
														goto l1206
													}
													position++
													goto l1205
												l1206:
													position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
													if buffer[position] != rune('P') {
														goto l1199
													}
													position++
												}
											l1205:
												{
													position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
													if buffer[position] != rune('l') {
														goto l1208
													}
													position++
													goto l1207
												l1208:
													position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
													if buffer[position] != rune('L') {
														goto l1199
													}
													position++
												}
											l1207:
												{
													position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l1210
													}
													position++
													goto l1209
												l1210:
													position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
													if buffer[position] != rune('A') {
														goto l1199
													}
													position++
												}
											l1209:
												{
													position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l1212
													}
													position++
													goto l1211
												l1212:
													position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
													if buffer[position] != rune('C') {
														goto l1199
													}
													position++
												}
											l1211:
												{
													position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1214
													}
													position++
													goto l1213
												l1214:
													position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
													if buffer[position] != rune('E') {
														goto l1199
													}
													position++
												}
											l1213:
												if !rules[ruleskip]() {
													goto l1199
												}
												depth--
												add(ruleREPLACE, position1200)
											}
											goto l1184
										l1199:
											position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
											{
												position1215 := position
												depth++
												{
													position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
													if buffer[position] != rune('r') {
														goto l1217
													}
													position++
													goto l1216
												l1217:
													position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
													if buffer[position] != rune('R') {
														goto l517
													}
													position++
												}
											l1216:
												{
													position1218, tokenIndex1218, depth1218 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1219
													}
													position++
													goto l1218
												l1219:
													position, tokenIndex, depth = position1218, tokenIndex1218, depth1218
													if buffer[position] != rune('E') {
														goto l517
													}
													position++
												}
											l1218:
												{
													position1220, tokenIndex1220, depth1220 := position, tokenIndex, depth
													if buffer[position] != rune('g') {
														goto l1221
													}
													position++
													goto l1220
												l1221:
													position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
													if buffer[position] != rune('G') {
														goto l517
													}
													position++
												}
											l1220:
												{
													position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l1223
													}
													position++
													goto l1222
												l1223:
													position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
													if buffer[position] != rune('E') {
														goto l517
													}
													position++
												}
											l1222:
												{
													position1224, tokenIndex1224, depth1224 := position, tokenIndex, depth
													if buffer[position] != rune('x') {
														goto l1225
													}
													position++
													goto l1224
												l1225:
													position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
													if buffer[position] != rune('X') {
														goto l517
													}
													position++
												}
											l1224:
												if !rules[ruleskip]() {
													goto l517
												}
												depth--
												add(ruleREGEX, position1215)
											}
										}
									l1184:
										if !rules[ruleLPAREN]() {
											goto l517
										}
										if !rules[ruleexpression]() {
											goto l517
										}
										if !rules[ruleCOMMA]() {
											goto l517
										}
										if !rules[ruleexpression]() {
											goto l517
										}
										{
											position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
											if !rules[ruleCOMMA]() {
												goto l1226
											}
											if !rules[ruleexpression]() {
												goto l1226
											}
											goto l1227
										l1226:
											position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
										}
									l1227:
										if !rules[ruleRPAREN]() {
											goto l517
										}
										break
									}
								}

							}
						l519:
							depth--
							add(rulebuiltinCall, position518)
						}
						goto l514
					l517:
						position, tokenIndex, depth = position514, tokenIndex514, depth514
						{
							position1229 := position
							depth++
							if !rules[ruleiriref]() {
								goto l1228
							}
							if !rules[ruleargList]() {
								goto l1228
							}
							depth--
							add(rulefunctionCall, position1229)
						}
						goto l514
					l1228:
						position, tokenIndex, depth = position514, tokenIndex514, depth514
						if !rules[ruleiriref]() {
							goto l1230
						}
						goto l514
					l1230:
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
			position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
			{
				position1238 := position
				depth++
				{
					position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l1240
					}
					goto l1239
				l1240:
					position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
					if !rules[ruleLPAREN]() {
						goto l1237
					}
					if !rules[ruleexpression]() {
						goto l1237
					}
				l1241:
					{
						position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l1242
						}
						if !rules[ruleexpression]() {
							goto l1242
						}
						goto l1241
					l1242:
						position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
					}
					if !rules[ruleRPAREN]() {
						goto l1237
					}
				}
			l1239:
				depth--
				add(ruleargList, position1238)
			}
			return true
		l1237:
			position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
			return false
		},
		/* 58 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		nil,
		/* 59 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
			{
				position1245 := position
				depth++
				{
					position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
					{
						position1248 := position
						depth++
					l1249:
						{
							position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
							{
								position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1252
								}
								position++
								goto l1251
							l1252:
								position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1250
								}
								position++
							}
						l1251:
							goto l1249
						l1250:
							position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
						}
						depth--
						add(rulePegText, position1248)
					}
					if buffer[position] != rune(':') {
						goto l1247
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1246
				l1247:
					position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
					{
						position1255 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1254
						}
						position++
					l1256:
						{
							position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1257
							}
							position++
							goto l1256
						l1257:
							position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
						}
						depth--
						add(rulePegText, position1255)
					}
					if buffer[position] != rune('/') {
						goto l1254
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1246
				l1254:
					position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
					{
						position1259 := position
						depth++
					l1260:
						{
							position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1261
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1261
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1261
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1261
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1261
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1261
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1261
									}
									position++
									break
								}
							}

							goto l1260
						l1261:
							position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
						}
						depth--
						add(rulePegText, position1259)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1246:
				if buffer[position] != rune('<') {
					goto l1244
				}
				position++
				if !rules[rulews]() {
					goto l1244
				}
				if !rules[ruleskip]() {
					goto l1244
				}
				depth--
				add(rulepof, position1245)
			}
			return true
		l1244:
			position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
			return false
		},
		/* 60 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
			{
				position1265 := position
				depth++
				{
					position1266, tokenIndex1266, depth1266 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1267
					}
					position++
					goto l1266
				l1267:
					position, tokenIndex, depth = position1266, tokenIndex1266, depth1266
					if buffer[position] != rune('$') {
						goto l1264
					}
					position++
				}
			l1266:
				{
					position1268 := position
					depth++
					{
						position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
						{
							position1273 := position
							depth++
							{
								position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
								{
									position1276 := position
									depth++
									{
										position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1278
										}
										position++
										goto l1277
									l1278:
										position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1275
										}
										position++
									}
								l1277:
									depth--
									add(rulePN_CHARS_BASE, position1276)
								}
								goto l1274
							l1275:
								position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
								if buffer[position] != rune('_') {
									goto l1272
								}
								position++
							}
						l1274:
							depth--
							add(rulePN_CHARS_U, position1273)
						}
						goto l1271
					l1272:
						position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1264
						}
						position++
					}
				l1271:
				l1269:
					{
						position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
						{
							position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
							{
								position1281 := position
								depth++
								{
									position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
									{
										position1284 := position
										depth++
										{
											position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1286
											}
											position++
											goto l1285
										l1286:
											position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1283
											}
											position++
										}
									l1285:
										depth--
										add(rulePN_CHARS_BASE, position1284)
									}
									goto l1282
								l1283:
									position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
									if buffer[position] != rune('_') {
										goto l1280
									}
									position++
								}
							l1282:
								depth--
								add(rulePN_CHARS_U, position1281)
							}
							goto l1279
						l1280:
							position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1270
							}
							position++
						}
					l1279:
						goto l1269
					l1270:
						position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
					}
					depth--
					add(ruleVARNAME, position1268)
				}
				if !rules[ruleskip]() {
					goto l1264
				}
				depth--
				add(rulevar, position1265)
			}
			return true
		l1264:
			position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
			return false
		},
		/* 61 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
			{
				position1288 := position
				depth++
				{
					position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1290
					}
					goto l1289
				l1290:
					position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
					{
						position1291 := position
						depth++
					l1292:
						{
							position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
							{
								position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
								{
									position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1296
									}
									position++
									goto l1295
								l1296:
									position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
									if buffer[position] != rune(' ') {
										goto l1294
									}
									position++
								}
							l1295:
								goto l1293
							l1294:
								position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
							}
							if !matchDot() {
								goto l1293
							}
							goto l1292
						l1293:
							position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
						}
						if buffer[position] != rune(':') {
							goto l1287
						}
						position++
					l1297:
						{
							position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
							{
								position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1300
								}
								position++
								goto l1299
							l1300:
								position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1301
								}
								position++
								goto l1299
							l1301:
								position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1302
								}
								position++
								goto l1299
							l1302:
								position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1298
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1298
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1298
										}
										position++
										break
									}
								}

							}
						l1299:
							goto l1297
						l1298:
							position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
						}
						if !rules[ruleskip]() {
							goto l1287
						}
						depth--
						add(ruleprefixedName, position1291)
					}
				}
			l1289:
				depth--
				add(ruleiriref, position1288)
			}
			return true
		l1287:
			position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
			return false
		},
		/* 62 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
			{
				position1305 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1304
				}
				position++
			l1306:
				{
					position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
					{
						position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1308
						}
						position++
						goto l1307
					l1308:
						position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
					}
					if !matchDot() {
						goto l1307
					}
					goto l1306
				l1307:
					position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
				}
				if buffer[position] != rune('>') {
					goto l1304
				}
				position++
				if !rules[ruleskip]() {
					goto l1304
				}
				depth--
				add(ruleiri, position1305)
			}
			return true
		l1304:
			position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
			return false
		},
		/* 63 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 64 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
			{
				position1311 := position
				depth++
				{
					position1312 := position
					depth++
					if buffer[position] != rune('"') {
						goto l1310
					}
					position++
				l1313:
					{
						position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
						{
							position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l1315
							}
							position++
							goto l1314
						l1315:
							position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
						}
						if !matchDot() {
							goto l1314
						}
						goto l1313
					l1314:
						position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
					}
					if buffer[position] != rune('"') {
						goto l1310
					}
					position++
					depth--
					add(rulestring, position1312)
				}
				{
					position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
					{
						position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1319
						}
						position++
						{
							position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1323
							}
							position++
							goto l1322
						l1323:
							position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1319
							}
							position++
						}
					l1322:
					l1320:
						{
							position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
							{
								position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1325
								}
								position++
								goto l1324
							l1325:
								position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1321
								}
								position++
							}
						l1324:
							goto l1320
						l1321:
							position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
						}
					l1326:
						{
							position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1327
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1327
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1327
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1327
									}
									position++
									break
								}
							}

						l1328:
							{
								position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1329
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1329
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1329
										}
										position++
										break
									}
								}

								goto l1328
							l1329:
								position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
							}
							goto l1326
						l1327:
							position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
						}
						goto l1318
					l1319:
						position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
						if buffer[position] != rune('^') {
							goto l1316
						}
						position++
						if buffer[position] != rune('^') {
							goto l1316
						}
						position++
						if !rules[ruleiriref]() {
							goto l1316
						}
					}
				l1318:
					goto l1317
				l1316:
					position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
				}
			l1317:
				if !rules[ruleskip]() {
					goto l1310
				}
				depth--
				add(ruleliteral, position1311)
			}
			return true
		l1310:
			position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
			return false
		},
		/* 65 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 66 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
			{
				position1334 := position
				depth++
				{
					position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
					{
						position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1338
						}
						position++
						goto l1337
					l1338:
						position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
						if buffer[position] != rune('-') {
							goto l1335
						}
						position++
					}
				l1337:
					goto l1336
				l1335:
					position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
				}
			l1336:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1333
				}
				position++
			l1339:
				{
					position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1340
					}
					position++
					goto l1339
				l1340:
					position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
				}
				{
					position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1341
					}
					position++
				l1343:
					{
						position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1344
						}
						position++
						goto l1343
					l1344:
						position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
					}
					goto l1342
				l1341:
					position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
				}
			l1342:
				if !rules[ruleskip]() {
					goto l1333
				}
				depth--
				add(rulenumericLiteral, position1334)
			}
			return true
		l1333:
			position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
			return false
		},
		/* 67 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 68 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
			{
				position1347 := position
				depth++
				{
					position1348, tokenIndex1348, depth1348 := position, tokenIndex, depth
					{
						position1350 := position
						depth++
						{
							position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1352
							}
							position++
							goto l1351
						l1352:
							position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
							if buffer[position] != rune('T') {
								goto l1349
							}
							position++
						}
					l1351:
						{
							position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1354
							}
							position++
							goto l1353
						l1354:
							position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
							if buffer[position] != rune('R') {
								goto l1349
							}
							position++
						}
					l1353:
						{
							position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1356
							}
							position++
							goto l1355
						l1356:
							position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
							if buffer[position] != rune('U') {
								goto l1349
							}
							position++
						}
					l1355:
						{
							position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1358
							}
							position++
							goto l1357
						l1358:
							position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
							if buffer[position] != rune('E') {
								goto l1349
							}
							position++
						}
					l1357:
						if !rules[ruleskip]() {
							goto l1349
						}
						depth--
						add(ruleTRUE, position1350)
					}
					goto l1348
				l1349:
					position, tokenIndex, depth = position1348, tokenIndex1348, depth1348
					{
						position1359 := position
						depth++
						{
							position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1361
							}
							position++
							goto l1360
						l1361:
							position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
							if buffer[position] != rune('F') {
								goto l1346
							}
							position++
						}
					l1360:
						{
							position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1363
							}
							position++
							goto l1362
						l1363:
							position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
							if buffer[position] != rune('A') {
								goto l1346
							}
							position++
						}
					l1362:
						{
							position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1365
							}
							position++
							goto l1364
						l1365:
							position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
							if buffer[position] != rune('L') {
								goto l1346
							}
							position++
						}
					l1364:
						{
							position1366, tokenIndex1366, depth1366 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1367
							}
							position++
							goto l1366
						l1367:
							position, tokenIndex, depth = position1366, tokenIndex1366, depth1366
							if buffer[position] != rune('S') {
								goto l1346
							}
							position++
						}
					l1366:
						{
							position1368, tokenIndex1368, depth1368 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1369
							}
							position++
							goto l1368
						l1369:
							position, tokenIndex, depth = position1368, tokenIndex1368, depth1368
							if buffer[position] != rune('E') {
								goto l1346
							}
							position++
						}
					l1368:
						if !rules[ruleskip]() {
							goto l1346
						}
						depth--
						add(ruleFALSE, position1359)
					}
				}
			l1348:
				depth--
				add(rulebooleanLiteral, position1347)
			}
			return true
		l1346:
			position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
			return false
		},
		/* 69 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 70 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 71 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 72 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
			{
				position1374 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1373
				}
				position++
			l1375:
				{
					position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1376
					}
					goto l1375
				l1376:
					position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
				}
				if buffer[position] != rune(')') {
					goto l1373
				}
				position++
				if !rules[ruleskip]() {
					goto l1373
				}
				depth--
				add(rulenil, position1374)
			}
			return true
		l1373:
			position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
			return false
		},
		/* 73 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 74 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 75 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 76 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 77 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 78 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 79 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 80 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 81 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 82 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 83 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 84 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 85 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 86 LBRACE <- <('{' skip)> */
		func() bool {
			position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
			{
				position1391 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1390
				}
				position++
				if !rules[ruleskip]() {
					goto l1390
				}
				depth--
				add(ruleLBRACE, position1391)
			}
			return true
		l1390:
			position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
			return false
		},
		/* 87 RBRACE <- <('}' skip)> */
		func() bool {
			position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
			{
				position1393 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1392
				}
				position++
				if !rules[ruleskip]() {
					goto l1392
				}
				depth--
				add(ruleRBRACE, position1393)
			}
			return true
		l1392:
			position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
			return false
		},
		/* 88 LBRACK <- <('[' skip)> */
		nil,
		/* 89 RBRACK <- <(']' skip)> */
		nil,
		/* 90 SEMICOLON <- <(';' skip)> */
		nil,
		/* 91 COMMA <- <(',' skip)> */
		func() bool {
			position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
			{
				position1398 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1397
				}
				position++
				if !rules[ruleskip]() {
					goto l1397
				}
				depth--
				add(ruleCOMMA, position1398)
			}
			return true
		l1397:
			position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
			return false
		},
		/* 92 DOT <- <('.' skip)> */
		func() bool {
			position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
			{
				position1400 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1399
				}
				position++
				if !rules[ruleskip]() {
					goto l1399
				}
				depth--
				add(ruleDOT, position1400)
			}
			return true
		l1399:
			position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
			return false
		},
		/* 93 COLON <- <(':' skip)> */
		nil,
		/* 94 PIPE <- <('|' skip)> */
		func() bool {
			position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
			{
				position1403 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1402
				}
				position++
				if !rules[ruleskip]() {
					goto l1402
				}
				depth--
				add(rulePIPE, position1403)
			}
			return true
		l1402:
			position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
			return false
		},
		/* 95 SLASH <- <('/' skip)> */
		func() bool {
			position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
			{
				position1405 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1404
				}
				position++
				if !rules[ruleskip]() {
					goto l1404
				}
				depth--
				add(ruleSLASH, position1405)
			}
			return true
		l1404:
			position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
			return false
		},
		/* 96 INVERSE <- <('^' skip)> */
		func() bool {
			position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
			{
				position1407 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1406
				}
				position++
				if !rules[ruleskip]() {
					goto l1406
				}
				depth--
				add(ruleINVERSE, position1407)
			}
			return true
		l1406:
			position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
			return false
		},
		/* 97 LPAREN <- <('(' skip)> */
		func() bool {
			position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
			{
				position1409 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1408
				}
				position++
				if !rules[ruleskip]() {
					goto l1408
				}
				depth--
				add(ruleLPAREN, position1409)
			}
			return true
		l1408:
			position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
			return false
		},
		/* 98 RPAREN <- <(')' skip)> */
		func() bool {
			position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
			{
				position1411 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1410
				}
				position++
				if !rules[ruleskip]() {
					goto l1410
				}
				depth--
				add(ruleRPAREN, position1411)
			}
			return true
		l1410:
			position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
			return false
		},
		/* 99 ISA <- <('a' skip)> */
		func() bool {
			position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
			{
				position1413 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1412
				}
				position++
				if !rules[ruleskip]() {
					goto l1412
				}
				depth--
				add(ruleISA, position1413)
			}
			return true
		l1412:
			position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
			return false
		},
		/* 100 NOT <- <('!' skip)> */
		func() bool {
			position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
			{
				position1415 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1414
				}
				position++
				if !rules[ruleskip]() {
					goto l1414
				}
				depth--
				add(ruleNOT, position1415)
			}
			return true
		l1414:
			position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
			return false
		},
		/* 101 STAR <- <('*' skip)> */
		func() bool {
			position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
			{
				position1417 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1416
				}
				position++
				if !rules[ruleskip]() {
					goto l1416
				}
				depth--
				add(ruleSTAR, position1417)
			}
			return true
		l1416:
			position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
			return false
		},
		/* 102 PLUS <- <('+' skip)> */
		func() bool {
			position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
			{
				position1419 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1418
				}
				position++
				if !rules[ruleskip]() {
					goto l1418
				}
				depth--
				add(rulePLUS, position1419)
			}
			return true
		l1418:
			position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
			return false
		},
		/* 103 MINUS <- <('-' skip)> */
		func() bool {
			position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
			{
				position1421 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1420
				}
				position++
				if !rules[ruleskip]() {
					goto l1420
				}
				depth--
				add(ruleMINUS, position1421)
			}
			return true
		l1420:
			position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
			return false
		},
		/* 104 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 105 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 106 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 107 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 108 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1426, tokenIndex1426, depth1426 := position, tokenIndex, depth
			{
				position1427 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1426
				}
				position++
			l1428:
				{
					position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1429
					}
					position++
					goto l1428
				l1429:
					position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
				}
				if !rules[ruleskip]() {
					goto l1426
				}
				depth--
				add(ruleINTEGER, position1427)
			}
			return true
		l1426:
			position, tokenIndex, depth = position1426, tokenIndex1426, depth1426
			return false
		},
		/* 109 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 110 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 111 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 112 OR <- <('|' '|' skip)> */
		nil,
		/* 113 AND <- <('&' '&' skip)> */
		nil,
		/* 114 EQ <- <('=' skip)> */
		nil,
		/* 115 NE <- <('!' '=' skip)> */
		nil,
		/* 116 GT <- <('>' skip)> */
		nil,
		/* 117 LT <- <('<' skip)> */
		nil,
		/* 118 LE <- <('<' '=' skip)> */
		nil,
		/* 119 GE <- <('>' '=' skip)> */
		nil,
		/* 120 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 121 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 122 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		nil,
		/* 123 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 124 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 125 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 126 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 127 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 128 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 129 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 130 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 131 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 132 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 133 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 134 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 135 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 136 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 137 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 138 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 139 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 140 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 141 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 142 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 143 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 144 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 145 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 146 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 147 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 148 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 149 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 150 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 151 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 152 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 153 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 154 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 155 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 156 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 157 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 158 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 159 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 160 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 161 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 162 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 163 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 164 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 165 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 166 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 167 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 168 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 169 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 170 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 171 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 172 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 173 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 174 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 175 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 176 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 177 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1499 := position
				depth++
			l1500:
				{
					position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
					{
						position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1503
						}
						goto l1502
					l1503:
						position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
						{
							position1504 := position
							depth++
							{
								position1505 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1501
								}
								position++
							l1506:
								{
									position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
									{
										position1508, tokenIndex1508, depth1508 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1508
										}
										goto l1507
									l1508:
										position, tokenIndex, depth = position1508, tokenIndex1508, depth1508
									}
									if !matchDot() {
										goto l1507
									}
									goto l1506
								l1507:
									position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
								}
								if !rules[ruleendOfLine]() {
									goto l1501
								}
								depth--
								add(rulePegText, position1505)
							}
							{
								add(ruleAction13, position)
							}
							depth--
							add(rulecomment, position1504)
						}
					}
				l1502:
					goto l1500
				l1501:
					position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
				}
				depth--
				add(ruleskip, position1499)
			}
			return true
		},
		/* 178 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1510, tokenIndex1510, depth1510 := position, tokenIndex, depth
			{
				position1511 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1510
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1510
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1510
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1510
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1510
						}
						break
					}
				}

				depth--
				add(rulews, position1511)
			}
			return true
		l1510:
			position, tokenIndex, depth = position1510, tokenIndex1510, depth1510
			return false
		},
		/* 179 comment <- <(<('#' (!endOfLine .)* endOfLine)> Action13)> */
		nil,
		/* 180 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1514, tokenIndex1514, depth1514 := position, tokenIndex, depth
			{
				position1515 := position
				depth++
				{
					position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1517
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1517
					}
					position++
					goto l1516
				l1517:
					position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
					if buffer[position] != rune('\n') {
						goto l1518
					}
					position++
					goto l1516
				l1518:
					position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
					if buffer[position] != rune('\r') {
						goto l1514
					}
					position++
				}
			l1516:
				depth--
				add(ruleendOfLine, position1515)
			}
			return true
		l1514:
			position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
			return false
		},
		nil,
		/* 183 Action0 <- <{ p.addPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 184 Action1 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 185 Action2 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 186 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 187 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 188 Action5 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 189 Action6 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 190 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 191 Action8 <- <{ p.setObject(p.skipComment(buffer, begin, end)); p.addTriplePattern() }> */
		nil,
		/* 192 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 193 Action10 <- <{ p.setPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 194 Action11 <- <{ p.setPathLength(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 195 Action12 <- <{ p.setKeyword(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 196 Action13 <- <{ p.commentBegin = begin }> */
		nil,
	}
	p.rules = rules
}
