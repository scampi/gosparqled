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
	rulefilterOrBind
	ruleconstraint
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
	ruleFILTER
	ruleBIND
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
	"filterOrBind",
	"constraint",
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
	"FILTER",
	"BIND",
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
	rules  [201]func() bool
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
							if !rules[ruleAS]() {
								goto l111
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
							position169 := position
							depth++
							{
								position170, tokenIndex170, depth170 := position, tokenIndex, depth
								if !rules[rulevar]() {
									goto l171
								}
								goto l170
							l171:
								position, tokenIndex, depth = position170, tokenIndex170, depth170
								if !rules[ruleLPAREN]() {
									goto l165
								}
								if !rules[ruleexpression]() {
									goto l165
								}
								if !rules[ruleAS]() {
									goto l165
								}
								if !rules[rulevar]() {
									goto l165
								}
								if !rules[ruleRPAREN]() {
									goto l165
								}
							}
						l170:
							depth--
							add(ruleprojectionElem, position169)
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
			position172, tokenIndex172, depth172 := position, tokenIndex, depth
			{
				position173 := position
				depth++
				if !rules[ruleselect]() {
					goto l172
				}
				if !rules[rulewhereClause]() {
					goto l172
				}
				depth--
				add(rulesubSelect, position173)
			}
			return true
		l172:
			position, tokenIndex, depth = position172, tokenIndex172, depth172
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
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				{
					position182 := position
					depth++
					{
						position183, tokenIndex183, depth183 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l184
						}
						position++
						goto l183
					l184:
						position, tokenIndex, depth = position183, tokenIndex183, depth183
						if buffer[position] != rune('F') {
							goto l180
						}
						position++
					}
				l183:
					{
						position185, tokenIndex185, depth185 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l186
						}
						position++
						goto l185
					l186:
						position, tokenIndex, depth = position185, tokenIndex185, depth185
						if buffer[position] != rune('R') {
							goto l180
						}
						position++
					}
				l185:
					{
						position187, tokenIndex187, depth187 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l188
						}
						position++
						goto l187
					l188:
						position, tokenIndex, depth = position187, tokenIndex187, depth187
						if buffer[position] != rune('O') {
							goto l180
						}
						position++
					}
				l187:
					{
						position189, tokenIndex189, depth189 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l190
						}
						position++
						goto l189
					l190:
						position, tokenIndex, depth = position189, tokenIndex189, depth189
						if buffer[position] != rune('M') {
							goto l180
						}
						position++
					}
				l189:
					if !rules[ruleskip]() {
						goto l180
					}
					depth--
					add(ruleFROM, position182)
				}
				{
					position191, tokenIndex191, depth191 := position, tokenIndex, depth
					{
						position193 := position
						depth++
						{
							position194, tokenIndex194, depth194 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l195
							}
							position++
							goto l194
						l195:
							position, tokenIndex, depth = position194, tokenIndex194, depth194
							if buffer[position] != rune('N') {
								goto l191
							}
							position++
						}
					l194:
						{
							position196, tokenIndex196, depth196 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l197
							}
							position++
							goto l196
						l197:
							position, tokenIndex, depth = position196, tokenIndex196, depth196
							if buffer[position] != rune('A') {
								goto l191
							}
							position++
						}
					l196:
						{
							position198, tokenIndex198, depth198 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l199
							}
							position++
							goto l198
						l199:
							position, tokenIndex, depth = position198, tokenIndex198, depth198
							if buffer[position] != rune('M') {
								goto l191
							}
							position++
						}
					l198:
						{
							position200, tokenIndex200, depth200 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l201
							}
							position++
							goto l200
						l201:
							position, tokenIndex, depth = position200, tokenIndex200, depth200
							if buffer[position] != rune('E') {
								goto l191
							}
							position++
						}
					l200:
						{
							position202, tokenIndex202, depth202 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l203
							}
							position++
							goto l202
						l203:
							position, tokenIndex, depth = position202, tokenIndex202, depth202
							if buffer[position] != rune('D') {
								goto l191
							}
							position++
						}
					l202:
						if !rules[ruleskip]() {
							goto l191
						}
						depth--
						add(ruleNAMED, position193)
					}
					goto l192
				l191:
					position, tokenIndex, depth = position191, tokenIndex191, depth191
				}
			l192:
				if !rules[ruleiriref]() {
					goto l180
				}
				depth--
				add(ruledatasetClause, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position204, tokenIndex204, depth204 := position, tokenIndex, depth
			{
				position205 := position
				depth++
				{
					position206, tokenIndex206, depth206 := position, tokenIndex, depth
					{
						position208 := position
						depth++
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('W') {
								goto l206
							}
							position++
						}
					l209:
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l212
							}
							position++
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if buffer[position] != rune('H') {
								goto l206
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
								goto l206
							}
							position++
						}
					l213:
						{
							position215, tokenIndex215, depth215 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l216
							}
							position++
							goto l215
						l216:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune('R') {
								goto l206
							}
							position++
						}
					l215:
						{
							position217, tokenIndex217, depth217 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l218
							}
							position++
							goto l217
						l218:
							position, tokenIndex, depth = position217, tokenIndex217, depth217
							if buffer[position] != rune('E') {
								goto l206
							}
							position++
						}
					l217:
						if !rules[ruleskip]() {
							goto l206
						}
						depth--
						add(ruleWHERE, position208)
					}
					goto l207
				l206:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
				}
			l207:
				if !rules[rulegroupGraphPattern]() {
					goto l204
				}
				depth--
				add(rulewhereClause, position205)
			}
			return true
		l204:
			position, tokenIndex, depth = position204, tokenIndex204, depth204
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position219, tokenIndex219, depth219 := position, tokenIndex, depth
			{
				position220 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l219
				}
				{
					position221, tokenIndex221, depth221 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l222
					}
					goto l221
				l222:
					position, tokenIndex, depth = position221, tokenIndex221, depth221
					if !rules[rulegraphPattern]() {
						goto l219
					}
				}
			l221:
				if !rules[ruleRBRACE]() {
					goto l219
				}
				depth--
				add(rulegroupGraphPattern, position220)
			}
			return true
		l219:
			position, tokenIndex, depth = position219, tokenIndex219, depth219
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position224 := position
				depth++
				{
					position225, tokenIndex225, depth225 := position, tokenIndex, depth
					{
						position227 := position
						depth++
						{
							position228, tokenIndex228, depth228 := position, tokenIndex, depth
							if !rules[ruletriplesBlock]() {
								goto l229
							}
						l230:
							{
								position231, tokenIndex231, depth231 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l231
								}
								{
									position232, tokenIndex232, depth232 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l232
									}
									goto l233
								l232:
									position, tokenIndex, depth = position232, tokenIndex232, depth232
								}
							l233:
								{
									position234, tokenIndex234, depth234 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l234
									}
									goto l235
								l234:
									position, tokenIndex, depth = position234, tokenIndex234, depth234
								}
							l235:
								goto l230
							l231:
								position, tokenIndex, depth = position231, tokenIndex231, depth231
							}
							goto l228
						l229:
							position, tokenIndex, depth = position228, tokenIndex228, depth228
							if !rules[rulefilterOrBind]() {
								goto l225
							}
							{
								position238, tokenIndex238, depth238 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l238
								}
								goto l239
							l238:
								position, tokenIndex, depth = position238, tokenIndex238, depth238
							}
						l239:
							{
								position240, tokenIndex240, depth240 := position, tokenIndex, depth
								if !rules[ruletriplesBlock]() {
									goto l240
								}
								goto l241
							l240:
								position, tokenIndex, depth = position240, tokenIndex240, depth240
							}
						l241:
						l236:
							{
								position237, tokenIndex237, depth237 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l237
								}
								{
									position242, tokenIndex242, depth242 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l242
									}
									goto l243
								l242:
									position, tokenIndex, depth = position242, tokenIndex242, depth242
								}
							l243:
								{
									position244, tokenIndex244, depth244 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l244
									}
									goto l245
								l244:
									position, tokenIndex, depth = position244, tokenIndex244, depth244
								}
							l245:
								goto l236
							l237:
								position, tokenIndex, depth = position237, tokenIndex237, depth237
							}
						}
					l228:
						depth--
						add(rulebasicGraphPattern, position227)
					}
					goto l226
				l225:
					position, tokenIndex, depth = position225, tokenIndex225, depth225
				}
			l226:
				{
					position246, tokenIndex246, depth246 := position, tokenIndex, depth
					{
						position248 := position
						depth++
						{
							position249, tokenIndex249, depth249 := position, tokenIndex, depth
							{
								position251 := position
								depth++
								{
									position252 := position
									depth++
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
											goto l250
										}
										position++
									}
								l253:
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l256
										}
										position++
										goto l255
									l256:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
										if buffer[position] != rune('P') {
											goto l250
										}
										position++
									}
								l255:
									{
										position257, tokenIndex257, depth257 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l258
										}
										position++
										goto l257
									l258:
										position, tokenIndex, depth = position257, tokenIndex257, depth257
										if buffer[position] != rune('T') {
											goto l250
										}
										position++
									}
								l257:
									{
										position259, tokenIndex259, depth259 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l260
										}
										position++
										goto l259
									l260:
										position, tokenIndex, depth = position259, tokenIndex259, depth259
										if buffer[position] != rune('I') {
											goto l250
										}
										position++
									}
								l259:
									{
										position261, tokenIndex261, depth261 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l262
										}
										position++
										goto l261
									l262:
										position, tokenIndex, depth = position261, tokenIndex261, depth261
										if buffer[position] != rune('O') {
											goto l250
										}
										position++
									}
								l261:
									{
										position263, tokenIndex263, depth263 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l264
										}
										position++
										goto l263
									l264:
										position, tokenIndex, depth = position263, tokenIndex263, depth263
										if buffer[position] != rune('N') {
											goto l250
										}
										position++
									}
								l263:
									{
										position265, tokenIndex265, depth265 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l266
										}
										position++
										goto l265
									l266:
										position, tokenIndex, depth = position265, tokenIndex265, depth265
										if buffer[position] != rune('A') {
											goto l250
										}
										position++
									}
								l265:
									{
										position267, tokenIndex267, depth267 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l268
										}
										position++
										goto l267
									l268:
										position, tokenIndex, depth = position267, tokenIndex267, depth267
										if buffer[position] != rune('L') {
											goto l250
										}
										position++
									}
								l267:
									if !rules[ruleskip]() {
										goto l250
									}
									depth--
									add(ruleOPTIONAL, position252)
								}
								if !rules[ruleLBRACE]() {
									goto l250
								}
								{
									position269, tokenIndex269, depth269 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l270
									}
									goto l269
								l270:
									position, tokenIndex, depth = position269, tokenIndex269, depth269
									if !rules[rulegraphPattern]() {
										goto l250
									}
								}
							l269:
								if !rules[ruleRBRACE]() {
									goto l250
								}
								depth--
								add(ruleoptionalGraphPattern, position251)
							}
							goto l249
						l250:
							position, tokenIndex, depth = position249, tokenIndex249, depth249
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l246
							}
						}
					l249:
						depth--
						add(rulegraphPatternNotTriples, position248)
					}
					{
						position271, tokenIndex271, depth271 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l271
						}
						goto l272
					l271:
						position, tokenIndex, depth = position271, tokenIndex271, depth271
					}
				l272:
					if !rules[rulegraphPattern]() {
						goto l246
					}
					goto l247
				l246:
					position, tokenIndex, depth = position246, tokenIndex246, depth246
				}
			l247:
				depth--
				add(rulegraphPattern, position224)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position275, tokenIndex275, depth275 := position, tokenIndex, depth
			{
				position276 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l275
				}
				{
					position277, tokenIndex277, depth277 := position, tokenIndex, depth
					{
						position279 := position
						depth++
						{
							position280, tokenIndex280, depth280 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l281
							}
							position++
							goto l280
						l281:
							position, tokenIndex, depth = position280, tokenIndex280, depth280
							if buffer[position] != rune('U') {
								goto l277
							}
							position++
						}
					l280:
						{
							position282, tokenIndex282, depth282 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l283
							}
							position++
							goto l282
						l283:
							position, tokenIndex, depth = position282, tokenIndex282, depth282
							if buffer[position] != rune('N') {
								goto l277
							}
							position++
						}
					l282:
						{
							position284, tokenIndex284, depth284 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l285
							}
							position++
							goto l284
						l285:
							position, tokenIndex, depth = position284, tokenIndex284, depth284
							if buffer[position] != rune('I') {
								goto l277
							}
							position++
						}
					l284:
						{
							position286, tokenIndex286, depth286 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l287
							}
							position++
							goto l286
						l287:
							position, tokenIndex, depth = position286, tokenIndex286, depth286
							if buffer[position] != rune('O') {
								goto l277
							}
							position++
						}
					l286:
						{
							position288, tokenIndex288, depth288 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l289
							}
							position++
							goto l288
						l289:
							position, tokenIndex, depth = position288, tokenIndex288, depth288
							if buffer[position] != rune('N') {
								goto l277
							}
							position++
						}
					l288:
						if !rules[ruleskip]() {
							goto l277
						}
						depth--
						add(ruleUNION, position279)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l277
					}
					goto l278
				l277:
					position, tokenIndex, depth = position277, tokenIndex277, depth277
				}
			l278:
				depth--
				add(rulegroupOrUnionGraphPattern, position276)
			}
			return true
		l275:
			position, tokenIndex, depth = position275, tokenIndex275, depth275
			return false
		},
		/* 21 basicGraphPattern <- <((triplesBlock (filterOrBind DOT? triplesBlock?)*) / (filterOrBind DOT? triplesBlock?)+)> */
		nil,
		/* 22 filterOrBind <- <((FILTER constraint) / (BIND LPAREN expression AS var RPAREN))> */
		func() bool {
			position291, tokenIndex291, depth291 := position, tokenIndex, depth
			{
				position292 := position
				depth++
				{
					position293, tokenIndex293, depth293 := position, tokenIndex, depth
					{
						position295 := position
						depth++
						{
							position296, tokenIndex296, depth296 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l297
							}
							position++
							goto l296
						l297:
							position, tokenIndex, depth = position296, tokenIndex296, depth296
							if buffer[position] != rune('F') {
								goto l294
							}
							position++
						}
					l296:
						{
							position298, tokenIndex298, depth298 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l299
							}
							position++
							goto l298
						l299:
							position, tokenIndex, depth = position298, tokenIndex298, depth298
							if buffer[position] != rune('I') {
								goto l294
							}
							position++
						}
					l298:
						{
							position300, tokenIndex300, depth300 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l301
							}
							position++
							goto l300
						l301:
							position, tokenIndex, depth = position300, tokenIndex300, depth300
							if buffer[position] != rune('L') {
								goto l294
							}
							position++
						}
					l300:
						{
							position302, tokenIndex302, depth302 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l303
							}
							position++
							goto l302
						l303:
							position, tokenIndex, depth = position302, tokenIndex302, depth302
							if buffer[position] != rune('T') {
								goto l294
							}
							position++
						}
					l302:
						{
							position304, tokenIndex304, depth304 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l305
							}
							position++
							goto l304
						l305:
							position, tokenIndex, depth = position304, tokenIndex304, depth304
							if buffer[position] != rune('E') {
								goto l294
							}
							position++
						}
					l304:
						{
							position306, tokenIndex306, depth306 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l307
							}
							position++
							goto l306
						l307:
							position, tokenIndex, depth = position306, tokenIndex306, depth306
							if buffer[position] != rune('R') {
								goto l294
							}
							position++
						}
					l306:
						if !rules[ruleskip]() {
							goto l294
						}
						depth--
						add(ruleFILTER, position295)
					}
					{
						position308 := position
						depth++
						{
							position309, tokenIndex309, depth309 := position, tokenIndex, depth
							if !rules[rulebrackettedExpression]() {
								goto l310
							}
							goto l309
						l310:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if !rules[rulebuiltinCall]() {
								goto l311
							}
							goto l309
						l311:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if !rules[rulefunctionCall]() {
								goto l294
							}
						}
					l309:
						depth--
						add(ruleconstraint, position308)
					}
					goto l293
				l294:
					position, tokenIndex, depth = position293, tokenIndex293, depth293
					{
						position312 := position
						depth++
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l314
							}
							position++
							goto l313
						l314:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
							if buffer[position] != rune('B') {
								goto l291
							}
							position++
						}
					l313:
						{
							position315, tokenIndex315, depth315 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l316
							}
							position++
							goto l315
						l316:
							position, tokenIndex, depth = position315, tokenIndex315, depth315
							if buffer[position] != rune('I') {
								goto l291
							}
							position++
						}
					l315:
						{
							position317, tokenIndex317, depth317 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l318
							}
							position++
							goto l317
						l318:
							position, tokenIndex, depth = position317, tokenIndex317, depth317
							if buffer[position] != rune('N') {
								goto l291
							}
							position++
						}
					l317:
						{
							position319, tokenIndex319, depth319 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l320
							}
							position++
							goto l319
						l320:
							position, tokenIndex, depth = position319, tokenIndex319, depth319
							if buffer[position] != rune('D') {
								goto l291
							}
							position++
						}
					l319:
						if !rules[ruleskip]() {
							goto l291
						}
						depth--
						add(ruleBIND, position312)
					}
					if !rules[ruleLPAREN]() {
						goto l291
					}
					if !rules[ruleexpression]() {
						goto l291
					}
					if !rules[ruleAS]() {
						goto l291
					}
					if !rules[rulevar]() {
						goto l291
					}
					if !rules[ruleRPAREN]() {
						goto l291
					}
				}
			l293:
				depth--
				add(rulefilterOrBind, position292)
			}
			return true
		l291:
			position, tokenIndex, depth = position291, tokenIndex291, depth291
			return false
		},
		/* 23 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		nil,
		/* 24 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l322
				}
			l324:
				{
					position325, tokenIndex325, depth325 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l325
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l325
					}
					goto l324
				l325:
					position, tokenIndex, depth = position325, tokenIndex325, depth325
				}
				{
					position326, tokenIndex326, depth326 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l326
					}
					goto l327
				l326:
					position, tokenIndex, depth = position326, tokenIndex326, depth326
				}
			l327:
				depth--
				add(ruletriplesBlock, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 25 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position328, tokenIndex328, depth328 := position, tokenIndex, depth
			{
				position329 := position
				depth++
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					{
						position332 := position
						depth++
						{
							position333, tokenIndex333, depth333 := position, tokenIndex, depth
							{
								position335 := position
								depth++
								if !rules[rulevar]() {
									goto l334
								}
								depth--
								add(rulePegText, position335)
							}
							{
								add(ruleAction1, position)
							}
							goto l333
						l334:
							position, tokenIndex, depth = position333, tokenIndex333, depth333
							{
								position338 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l337
								}
								depth--
								add(rulePegText, position338)
							}
							{
								add(ruleAction2, position)
							}
							goto l333
						l337:
							position, tokenIndex, depth = position333, tokenIndex333, depth333
							if !rules[rulepof]() {
								goto l331
							}
							{
								add(ruleAction3, position)
							}
						}
					l333:
						depth--
						add(rulevarOrTerm, position332)
					}
					if !rules[rulepropertyListPath]() {
						goto l331
					}
					goto l330
				l331:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
					{
						position341 := position
						depth++
						{
							position342, tokenIndex342, depth342 := position, tokenIndex, depth
							{
								position344 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l343
								}
								if !rules[rulegraphNodePath]() {
									goto l343
								}
							l345:
								{
									position346, tokenIndex346, depth346 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l346
									}
									goto l345
								l346:
									position, tokenIndex, depth = position346, tokenIndex346, depth346
								}
								if !rules[ruleRPAREN]() {
									goto l343
								}
								depth--
								add(rulecollectionPath, position344)
							}
							goto l342
						l343:
							position, tokenIndex, depth = position342, tokenIndex342, depth342
							{
								position347 := position
								depth++
								{
									position348 := position
									depth++
									if buffer[position] != rune('[') {
										goto l328
									}
									position++
									if !rules[ruleskip]() {
										goto l328
									}
									depth--
									add(ruleLBRACK, position348)
								}
								if !rules[rulepropertyListPath]() {
									goto l328
								}
								{
									position349 := position
									depth++
									if buffer[position] != rune(']') {
										goto l328
									}
									position++
									if !rules[ruleskip]() {
										goto l328
									}
									depth--
									add(ruleRBRACK, position349)
								}
								depth--
								add(ruleblankNodePropertyListPath, position347)
							}
						}
					l342:
						depth--
						add(ruletriplesNodePath, position341)
					}
					if !rules[rulepropertyListPath]() {
						goto l328
					}
				}
			l330:
				depth--
				add(ruletriplesSameSubjectPath, position329)
			}
			return true
		l328:
			position, tokenIndex, depth = position328, tokenIndex328, depth328
			return false
		},
		/* 26 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 27 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position351, tokenIndex351, depth351 := position, tokenIndex, depth
			{
				position352 := position
				depth++
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l354
					}
					goto l353
				l354:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
					{
						switch buffer[position] {
						case '(':
							if !rules[rulenil]() {
								goto l351
							}
							break
						case '[', '_':
							{
								position356 := position
								depth++
								{
									position357, tokenIndex357, depth357 := position, tokenIndex, depth
									{
										position359 := position
										depth++
										if buffer[position] != rune('_') {
											goto l358
										}
										position++
										if buffer[position] != rune(':') {
											goto l358
										}
										position++
										{
											switch buffer[position] {
											case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l358
												}
												position++
												break
											case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l358
												}
												position++
												break
											default:
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l358
												}
												position++
												break
											}
										}

										{
											position361, tokenIndex361, depth361 := position, tokenIndex, depth
											{
												position363, tokenIndex363, depth363 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l364
												}
												position++
												goto l363
											l364:
												position, tokenIndex, depth = position363, tokenIndex363, depth363
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l365
												}
												position++
												goto l363
											l365:
												position, tokenIndex, depth = position363, tokenIndex363, depth363
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l366
												}
												position++
												goto l363
											l366:
												position, tokenIndex, depth = position363, tokenIndex363, depth363
												if c := buffer[position]; c < rune('.') || c > rune('_') {
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
										if !rules[ruleskip]() {
											goto l358
										}
										depth--
										add(ruleblankNodeLabel, position359)
									}
									goto l357
								l358:
									position, tokenIndex, depth = position357, tokenIndex357, depth357
									{
										position367 := position
										depth++
										if buffer[position] != rune('[') {
											goto l351
										}
										position++
									l368:
										{
											position369, tokenIndex369, depth369 := position, tokenIndex, depth
											if !rules[rulews]() {
												goto l369
											}
											goto l368
										l369:
											position, tokenIndex, depth = position369, tokenIndex369, depth369
										}
										if buffer[position] != rune(']') {
											goto l351
										}
										position++
										if !rules[ruleskip]() {
											goto l351
										}
										depth--
										add(ruleanon, position367)
									}
								}
							l357:
								depth--
								add(ruleblankNode, position356)
							}
							break
						case 'F', 'T', 'f', 't':
							if !rules[rulebooleanLiteral]() {
								goto l351
							}
							break
						case '"':
							if !rules[ruleliteral]() {
								goto l351
							}
							break
						default:
							if !rules[rulenumericLiteral]() {
								goto l351
							}
							break
						}
					}

				}
			l353:
				depth--
				add(rulegraphTerm, position352)
			}
			return true
		l351:
			position, tokenIndex, depth = position351, tokenIndex351, depth351
			return false
		},
		/* 28 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 29 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 30 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 31 propertyListPath <- <(((pof Action4) / (<var> Action5) / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			position373, tokenIndex373, depth373 := position, tokenIndex, depth
			{
				position374 := position
				depth++
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l376
					}
					{
						add(ruleAction4, position)
					}
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					{
						position379 := position
						depth++
						if !rules[rulevar]() {
							goto l378
						}
						depth--
						add(rulePegText, position379)
					}
					{
						add(ruleAction5, position)
					}
					goto l375
				l378:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					{
						position381 := position
						depth++
						if !rules[rulepath]() {
							goto l373
						}
						depth--
						add(ruleverbPath, position381)
					}
				}
			l375:
				if !rules[ruleobjectListPath]() {
					goto l373
				}
				{
					position382, tokenIndex382, depth382 := position, tokenIndex, depth
					{
						position384 := position
						depth++
						if buffer[position] != rune(';') {
							goto l382
						}
						position++
						if !rules[ruleskip]() {
							goto l382
						}
						depth--
						add(ruleSEMICOLON, position384)
					}
					if !rules[rulepropertyListPath]() {
						goto l382
					}
					goto l383
				l382:
					position, tokenIndex, depth = position382, tokenIndex382, depth382
				}
			l383:
				depth--
				add(rulepropertyListPath, position374)
			}
			return true
		l373:
			position, tokenIndex, depth = position373, tokenIndex373, depth373
			return false
		},
		/* 32 verbPath <- <path> */
		nil,
		/* 33 path <- <pathAlternative> */
		func() bool {
			position386, tokenIndex386, depth386 := position, tokenIndex, depth
			{
				position387 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l386
				}
				depth--
				add(rulepath, position387)
			}
			return true
		l386:
			position, tokenIndex, depth = position386, tokenIndex386, depth386
			return false
		},
		/* 34 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position388, tokenIndex388, depth388 := position, tokenIndex, depth
			{
				position389 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l388
				}
			l390:
				{
					position391, tokenIndex391, depth391 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l391
					}
					if !rules[rulepathAlternative]() {
						goto l391
					}
					goto l390
				l391:
					position, tokenIndex, depth = position391, tokenIndex391, depth391
				}
				depth--
				add(rulepathAlternative, position389)
			}
			return true
		l388:
			position, tokenIndex, depth = position388, tokenIndex388, depth388
			return false
		},
		/* 35 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position392, tokenIndex392, depth392 := position, tokenIndex, depth
			{
				position393 := position
				depth++
				{
					position394 := position
					depth++
					{
						position395 := position
						depth++
						{
							position396, tokenIndex396, depth396 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l396
							}
							goto l397
						l396:
							position, tokenIndex, depth = position396, tokenIndex396, depth396
						}
					l397:
						{
							position398 := position
							depth++
							{
								position399, tokenIndex399, depth399 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l400
								}
								goto l399
							l400:
								position, tokenIndex, depth = position399, tokenIndex399, depth399
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l392
										}
										if !rules[rulepath]() {
											goto l392
										}
										if !rules[ruleRPAREN]() {
											goto l392
										}
										break
									case '!':
										if !rules[ruleNOT]() {
											goto l392
										}
										{
											position402 := position
											depth++
											{
												position403, tokenIndex403, depth403 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l404
												}
												goto l403
											l404:
												position, tokenIndex, depth = position403, tokenIndex403, depth403
												if !rules[ruleLPAREN]() {
													goto l392
												}
												{
													position405, tokenIndex405, depth405 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l405
													}
												l407:
													{
														position408, tokenIndex408, depth408 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l408
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l408
														}
														goto l407
													l408:
														position, tokenIndex, depth = position408, tokenIndex408, depth408
													}
													goto l406
												l405:
													position, tokenIndex, depth = position405, tokenIndex405, depth405
												}
											l406:
												if !rules[ruleRPAREN]() {
													goto l392
												}
											}
										l403:
											depth--
											add(rulepathNegatedPropertySet, position402)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l392
										}
										break
									}
								}

							}
						l399:
							depth--
							add(rulepathPrimary, position398)
						}
						depth--
						add(rulepathElt, position395)
					}
					depth--
					add(rulePegText, position394)
				}
				{
					add(ruleAction6, position)
				}
			l410:
				{
					position411, tokenIndex411, depth411 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l411
					}
					if !rules[rulepathSequence]() {
						goto l411
					}
					goto l410
				l411:
					position, tokenIndex, depth = position411, tokenIndex411, depth411
				}
				depth--
				add(rulepathSequence, position393)
			}
			return true
		l392:
			position, tokenIndex, depth = position392, tokenIndex392, depth392
			return false
		},
		/* 36 pathElt <- <(INVERSE? pathPrimary)> */
		nil,
		/* 37 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 38 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 39 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position415, tokenIndex415, depth415 := position, tokenIndex, depth
			{
				position416 := position
				depth++
				{
					position417, tokenIndex417, depth417 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l418
					}
					goto l417
				l418:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if !rules[ruleISA]() {
						goto l419
					}
					goto l417
				l419:
					position, tokenIndex, depth = position417, tokenIndex417, depth417
					if !rules[ruleINVERSE]() {
						goto l415
					}
					{
						position420, tokenIndex420, depth420 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l421
						}
						goto l420
					l421:
						position, tokenIndex, depth = position420, tokenIndex420, depth420
						if !rules[ruleISA]() {
							goto l415
						}
					}
				l420:
				}
			l417:
				depth--
				add(rulepathOneInPropertySet, position416)
			}
			return true
		l415:
			position, tokenIndex, depth = position415, tokenIndex415, depth415
			return false
		},
		/* 40 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			{
				position423 := position
				depth++
				{
					position424 := position
					depth++
					{
						position425, tokenIndex425, depth425 := position, tokenIndex, depth
						if !rules[rulepof]() {
							goto l426
						}
						{
							add(ruleAction7, position)
						}
						goto l425
					l426:
						position, tokenIndex, depth = position425, tokenIndex425, depth425
						{
							position429 := position
							depth++
							if !rules[rulegraphNodePath]() {
								goto l428
							}
							depth--
							add(rulePegText, position429)
						}
						{
							add(ruleAction8, position)
						}
						goto l425
					l428:
						position, tokenIndex, depth = position425, tokenIndex425, depth425
						{
							add(ruleAction9, position)
						}
					}
				l425:
					depth--
					add(ruleobjectPath, position424)
				}
			l432:
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if !rules[ruleCOMMA]() {
						goto l433
					}
					if !rules[ruleobjectListPath]() {
						goto l433
					}
					goto l432
				l433:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
				}
				depth--
				add(ruleobjectListPath, position423)
			}
			return true
		},
		/* 41 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		nil,
		/* 42 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l438
					}
					goto l437
				l438:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
					if !rules[rulegraphTerm]() {
						goto l435
					}
				}
			l437:
				depth--
				add(rulegraphNodePath, position436)
			}
			return true
		l435:
			position, tokenIndex, depth = position435, tokenIndex435, depth435
			return false
		},
		/* 43 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position440 := position
				depth++
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					{
						position443 := position
						depth++
						{
							position444, tokenIndex444, depth444 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l445
							}
							{
								position446, tokenIndex446, depth446 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l446
								}
								goto l447
							l446:
								position, tokenIndex, depth = position446, tokenIndex446, depth446
							}
						l447:
							goto l444
						l445:
							position, tokenIndex, depth = position444, tokenIndex444, depth444
							if !rules[ruleoffset]() {
								goto l441
							}
							{
								position448, tokenIndex448, depth448 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l448
								}
								goto l449
							l448:
								position, tokenIndex, depth = position448, tokenIndex448, depth448
							}
						l449:
						}
					l444:
						depth--
						add(rulelimitOffsetClauses, position443)
					}
					goto l442
				l441:
					position, tokenIndex, depth = position441, tokenIndex441, depth441
				}
			l442:
				depth--
				add(rulesolutionModifier, position440)
			}
			return true
		},
		/* 44 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 45 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position451, tokenIndex451, depth451 := position, tokenIndex, depth
			{
				position452 := position
				depth++
				{
					position453 := position
					depth++
					{
						position454, tokenIndex454, depth454 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l455
						}
						position++
						goto l454
					l455:
						position, tokenIndex, depth = position454, tokenIndex454, depth454
						if buffer[position] != rune('L') {
							goto l451
						}
						position++
					}
				l454:
					{
						position456, tokenIndex456, depth456 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l457
						}
						position++
						goto l456
					l457:
						position, tokenIndex, depth = position456, tokenIndex456, depth456
						if buffer[position] != rune('I') {
							goto l451
						}
						position++
					}
				l456:
					{
						position458, tokenIndex458, depth458 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l459
						}
						position++
						goto l458
					l459:
						position, tokenIndex, depth = position458, tokenIndex458, depth458
						if buffer[position] != rune('M') {
							goto l451
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
							goto l451
						}
						position++
					}
				l460:
					{
						position462, tokenIndex462, depth462 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l463
						}
						position++
						goto l462
					l463:
						position, tokenIndex, depth = position462, tokenIndex462, depth462
						if buffer[position] != rune('T') {
							goto l451
						}
						position++
					}
				l462:
					if !rules[ruleskip]() {
						goto l451
					}
					depth--
					add(ruleLIMIT, position453)
				}
				if !rules[ruleINTEGER]() {
					goto l451
				}
				depth--
				add(rulelimit, position452)
			}
			return true
		l451:
			position, tokenIndex, depth = position451, tokenIndex451, depth451
			return false
		},
		/* 46 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position464, tokenIndex464, depth464 := position, tokenIndex, depth
			{
				position465 := position
				depth++
				{
					position466 := position
					depth++
					{
						position467, tokenIndex467, depth467 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l468
						}
						position++
						goto l467
					l468:
						position, tokenIndex, depth = position467, tokenIndex467, depth467
						if buffer[position] != rune('O') {
							goto l464
						}
						position++
					}
				l467:
					{
						position469, tokenIndex469, depth469 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l470
						}
						position++
						goto l469
					l470:
						position, tokenIndex, depth = position469, tokenIndex469, depth469
						if buffer[position] != rune('F') {
							goto l464
						}
						position++
					}
				l469:
					{
						position471, tokenIndex471, depth471 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l472
						}
						position++
						goto l471
					l472:
						position, tokenIndex, depth = position471, tokenIndex471, depth471
						if buffer[position] != rune('F') {
							goto l464
						}
						position++
					}
				l471:
					{
						position473, tokenIndex473, depth473 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l474
						}
						position++
						goto l473
					l474:
						position, tokenIndex, depth = position473, tokenIndex473, depth473
						if buffer[position] != rune('S') {
							goto l464
						}
						position++
					}
				l473:
					{
						position475, tokenIndex475, depth475 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l476
						}
						position++
						goto l475
					l476:
						position, tokenIndex, depth = position475, tokenIndex475, depth475
						if buffer[position] != rune('E') {
							goto l464
						}
						position++
					}
				l475:
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l478
						}
						position++
						goto l477
					l478:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
						if buffer[position] != rune('T') {
							goto l464
						}
						position++
					}
				l477:
					if !rules[ruleskip]() {
						goto l464
					}
					depth--
					add(ruleOFFSET, position466)
				}
				if !rules[ruleINTEGER]() {
					goto l464
				}
				depth--
				add(ruleoffset, position465)
			}
			return true
		l464:
			position, tokenIndex, depth = position464, tokenIndex464, depth464
			return false
		},
		/* 47 expression <- <conditionalOrExpression> */
		func() bool {
			position479, tokenIndex479, depth479 := position, tokenIndex, depth
			{
				position480 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l479
				}
				depth--
				add(ruleexpression, position480)
			}
			return true
		l479:
			position, tokenIndex, depth = position479, tokenIndex479, depth479
			return false
		},
		/* 48 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position481, tokenIndex481, depth481 := position, tokenIndex, depth
			{
				position482 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l481
				}
				{
					position483, tokenIndex483, depth483 := position, tokenIndex, depth
					{
						position485 := position
						depth++
						if buffer[position] != rune('|') {
							goto l483
						}
						position++
						if buffer[position] != rune('|') {
							goto l483
						}
						position++
						if !rules[ruleskip]() {
							goto l483
						}
						depth--
						add(ruleOR, position485)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l483
					}
					goto l484
				l483:
					position, tokenIndex, depth = position483, tokenIndex483, depth483
				}
			l484:
				depth--
				add(ruleconditionalOrExpression, position482)
			}
			return true
		l481:
			position, tokenIndex, depth = position481, tokenIndex481, depth481
			return false
		},
		/* 49 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position486, tokenIndex486, depth486 := position, tokenIndex, depth
			{
				position487 := position
				depth++
				{
					position488 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l486
					}
					{
						position489, tokenIndex489, depth489 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position492 := position
									depth++
									{
										position493 := position
										depth++
										{
											position494, tokenIndex494, depth494 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l495
											}
											position++
											goto l494
										l495:
											position, tokenIndex, depth = position494, tokenIndex494, depth494
											if buffer[position] != rune('N') {
												goto l489
											}
											position++
										}
									l494:
										{
											position496, tokenIndex496, depth496 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l497
											}
											position++
											goto l496
										l497:
											position, tokenIndex, depth = position496, tokenIndex496, depth496
											if buffer[position] != rune('O') {
												goto l489
											}
											position++
										}
									l496:
										{
											position498, tokenIndex498, depth498 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l499
											}
											position++
											goto l498
										l499:
											position, tokenIndex, depth = position498, tokenIndex498, depth498
											if buffer[position] != rune('T') {
												goto l489
											}
											position++
										}
									l498:
										if buffer[position] != rune(' ') {
											goto l489
										}
										position++
										{
											position500, tokenIndex500, depth500 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l501
											}
											position++
											goto l500
										l501:
											position, tokenIndex, depth = position500, tokenIndex500, depth500
											if buffer[position] != rune('I') {
												goto l489
											}
											position++
										}
									l500:
										{
											position502, tokenIndex502, depth502 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l503
											}
											position++
											goto l502
										l503:
											position, tokenIndex, depth = position502, tokenIndex502, depth502
											if buffer[position] != rune('N') {
												goto l489
											}
											position++
										}
									l502:
										if !rules[ruleskip]() {
											goto l489
										}
										depth--
										add(ruleNOTIN, position493)
									}
									if !rules[ruleargList]() {
										goto l489
									}
									depth--
									add(rulenotin, position492)
								}
								break
							case 'I', 'i':
								{
									position504 := position
									depth++
									{
										position505 := position
										depth++
										{
											position506, tokenIndex506, depth506 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l507
											}
											position++
											goto l506
										l507:
											position, tokenIndex, depth = position506, tokenIndex506, depth506
											if buffer[position] != rune('I') {
												goto l489
											}
											position++
										}
									l506:
										{
											position508, tokenIndex508, depth508 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l509
											}
											position++
											goto l508
										l509:
											position, tokenIndex, depth = position508, tokenIndex508, depth508
											if buffer[position] != rune('N') {
												goto l489
											}
											position++
										}
									l508:
										if !rules[ruleskip]() {
											goto l489
										}
										depth--
										add(ruleIN, position505)
									}
									if !rules[ruleargList]() {
										goto l489
									}
									depth--
									add(rulein, position504)
								}
								break
							default:
								{
									position510, tokenIndex510, depth510 := position, tokenIndex, depth
									{
										position512 := position
										depth++
										if buffer[position] != rune('<') {
											goto l511
										}
										position++
										if !rules[ruleskip]() {
											goto l511
										}
										depth--
										add(ruleLT, position512)
									}
									goto l510
								l511:
									position, tokenIndex, depth = position510, tokenIndex510, depth510
									{
										position514 := position
										depth++
										if buffer[position] != rune('>') {
											goto l513
										}
										position++
										if buffer[position] != rune('=') {
											goto l513
										}
										position++
										if !rules[ruleskip]() {
											goto l513
										}
										depth--
										add(ruleGE, position514)
									}
									goto l510
								l513:
									position, tokenIndex, depth = position510, tokenIndex510, depth510
									{
										switch buffer[position] {
										case '>':
											{
												position516 := position
												depth++
												if buffer[position] != rune('>') {
													goto l489
												}
												position++
												if !rules[ruleskip]() {
													goto l489
												}
												depth--
												add(ruleGT, position516)
											}
											break
										case '<':
											{
												position517 := position
												depth++
												if buffer[position] != rune('<') {
													goto l489
												}
												position++
												if buffer[position] != rune('=') {
													goto l489
												}
												position++
												if !rules[ruleskip]() {
													goto l489
												}
												depth--
												add(ruleLE, position517)
											}
											break
										case '!':
											{
												position518 := position
												depth++
												if buffer[position] != rune('!') {
													goto l489
												}
												position++
												if buffer[position] != rune('=') {
													goto l489
												}
												position++
												if !rules[ruleskip]() {
													goto l489
												}
												depth--
												add(ruleNE, position518)
											}
											break
										default:
											{
												position519 := position
												depth++
												if buffer[position] != rune('=') {
													goto l489
												}
												position++
												if !rules[ruleskip]() {
													goto l489
												}
												depth--
												add(ruleEQ, position519)
											}
											break
										}
									}

								}
							l510:
								if !rules[rulenumericExpression]() {
									goto l489
								}
								break
							}
						}

						goto l490
					l489:
						position, tokenIndex, depth = position489, tokenIndex489, depth489
					}
				l490:
					depth--
					add(rulevalueLogical, position488)
				}
				{
					position520, tokenIndex520, depth520 := position, tokenIndex, depth
					{
						position522 := position
						depth++
						if buffer[position] != rune('&') {
							goto l520
						}
						position++
						if buffer[position] != rune('&') {
							goto l520
						}
						position++
						if !rules[ruleskip]() {
							goto l520
						}
						depth--
						add(ruleAND, position522)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l520
					}
					goto l521
				l520:
					position, tokenIndex, depth = position520, tokenIndex520, depth520
				}
			l521:
				depth--
				add(ruleconditionalAndExpression, position487)
			}
			return true
		l486:
			position, tokenIndex, depth = position486, tokenIndex486, depth486
			return false
		},
		/* 50 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 51 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l524
				}
			l526:
				{
					position527, tokenIndex527, depth527 := position, tokenIndex, depth
					{
						position528, tokenIndex528, depth528 := position, tokenIndex, depth
						{
							position530, tokenIndex530, depth530 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l531
							}
							goto l530
						l531:
							position, tokenIndex, depth = position530, tokenIndex530, depth530
							if !rules[ruleMINUS]() {
								goto l529
							}
						}
					l530:
						if !rules[rulemultiplicativeExpression]() {
							goto l529
						}
						goto l528
					l529:
						position, tokenIndex, depth = position528, tokenIndex528, depth528
						{
							position532 := position
							depth++
							{
								position533, tokenIndex533, depth533 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l534
								}
								position++
								goto l533
							l534:
								position, tokenIndex, depth = position533, tokenIndex533, depth533
								if buffer[position] != rune('-') {
									goto l527
								}
								position++
							}
						l533:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l527
							}
							position++
						l535:
							{
								position536, tokenIndex536, depth536 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l536
								}
								position++
								goto l535
							l536:
								position, tokenIndex, depth = position536, tokenIndex536, depth536
							}
							{
								position537, tokenIndex537, depth537 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l537
								}
								position++
							l539:
								{
									position540, tokenIndex540, depth540 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l540
									}
									position++
									goto l539
								l540:
									position, tokenIndex, depth = position540, tokenIndex540, depth540
								}
								goto l538
							l537:
								position, tokenIndex, depth = position537, tokenIndex537, depth537
							}
						l538:
							if !rules[ruleskip]() {
								goto l527
							}
							depth--
							add(rulesignedNumericLiteral, position532)
						}
					}
				l528:
					goto l526
				l527:
					position, tokenIndex, depth = position527, tokenIndex527, depth527
				}
				depth--
				add(rulenumericExpression, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 52 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position541, tokenIndex541, depth541 := position, tokenIndex, depth
			{
				position542 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l541
				}
			l543:
				{
					position544, tokenIndex544, depth544 := position, tokenIndex, depth
					{
						position545, tokenIndex545, depth545 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l546
						}
						goto l545
					l546:
						position, tokenIndex, depth = position545, tokenIndex545, depth545
						if !rules[ruleSLASH]() {
							goto l544
						}
					}
				l545:
					if !rules[ruleunaryExpression]() {
						goto l544
					}
					goto l543
				l544:
					position, tokenIndex, depth = position544, tokenIndex544, depth544
				}
				depth--
				add(rulemultiplicativeExpression, position542)
			}
			return true
		l541:
			position, tokenIndex, depth = position541, tokenIndex541, depth541
			return false
		},
		/* 53 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position547, tokenIndex547, depth547 := position, tokenIndex, depth
			{
				position548 := position
				depth++
				{
					position549, tokenIndex549, depth549 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l549
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l549
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l549
							}
							break
						}
					}

					goto l550
				l549:
					position, tokenIndex, depth = position549, tokenIndex549, depth549
				}
			l550:
				{
					position552 := position
					depth++
					{
						position553, tokenIndex553, depth553 := position, tokenIndex, depth
						if !rules[rulebrackettedExpression]() {
							goto l554
						}
						goto l553
					l554:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						if !rules[rulebuiltinCall]() {
							goto l555
						}
						goto l553
					l555:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						if !rules[rulefunctionCall]() {
							goto l556
						}
						goto l553
					l556:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						if !rules[ruleiriref]() {
							goto l557
						}
						goto l553
					l557:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							switch buffer[position] {
							case '$', '?':
								if !rules[rulevar]() {
									goto l547
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l547
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l547
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l547
								}
								break
							}
						}

					}
				l553:
					depth--
					add(ruleprimaryExpression, position552)
				}
				depth--
				add(ruleunaryExpression, position548)
			}
			return true
		l547:
			position, tokenIndex, depth = position547, tokenIndex547, depth547
			return false
		},
		/* 54 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 55 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position560, tokenIndex560, depth560 := position, tokenIndex, depth
			{
				position561 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l560
				}
				if !rules[ruleexpression]() {
					goto l560
				}
				if !rules[ruleRPAREN]() {
					goto l560
				}
				depth--
				add(rulebrackettedExpression, position561)
			}
			return true
		l560:
			position, tokenIndex, depth = position560, tokenIndex560, depth560
			return false
		},
		/* 56 functionCall <- <(iriref argList)> */
		func() bool {
			position562, tokenIndex562, depth562 := position, tokenIndex, depth
			{
				position563 := position
				depth++
				if !rules[ruleiriref]() {
					goto l562
				}
				if !rules[ruleargList]() {
					goto l562
				}
				depth--
				add(rulefunctionCall, position563)
			}
			return true
		l562:
			position, tokenIndex, depth = position562, tokenIndex562, depth562
			return false
		},
		/* 57 in <- <(IN argList)> */
		nil,
		/* 58 notin <- <(NOTIN argList)> */
		nil,
		/* 59 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position566, tokenIndex566, depth566 := position, tokenIndex, depth
			{
				position567 := position
				depth++
				{
					position568, tokenIndex568, depth568 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l569
					}
					goto l568
				l569:
					position, tokenIndex, depth = position568, tokenIndex568, depth568
					if !rules[ruleLPAREN]() {
						goto l566
					}
					if !rules[ruleexpression]() {
						goto l566
					}
				l570:
					{
						position571, tokenIndex571, depth571 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l571
						}
						if !rules[ruleexpression]() {
							goto l571
						}
						goto l570
					l571:
						position, tokenIndex, depth = position571, tokenIndex571, depth571
					}
					if !rules[ruleRPAREN]() {
						goto l566
					}
				}
			l568:
				depth--
				add(ruleargList, position567)
			}
			return true
		l566:
			position, tokenIndex, depth = position566, tokenIndex566, depth566
			return false
		},
		/* 60 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		func() bool {
			position572, tokenIndex572, depth572 := position, tokenIndex, depth
			{
				position573 := position
				depth++
				{
					position574, tokenIndex574, depth574 := position, tokenIndex, depth
					{
						position576, tokenIndex576, depth576 := position, tokenIndex, depth
						{
							position578 := position
							depth++
							{
								position579, tokenIndex579, depth579 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l580
								}
								position++
								goto l579
							l580:
								position, tokenIndex, depth = position579, tokenIndex579, depth579
								if buffer[position] != rune('S') {
									goto l577
								}
								position++
							}
						l579:
							{
								position581, tokenIndex581, depth581 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l582
								}
								position++
								goto l581
							l582:
								position, tokenIndex, depth = position581, tokenIndex581, depth581
								if buffer[position] != rune('T') {
									goto l577
								}
								position++
							}
						l581:
							{
								position583, tokenIndex583, depth583 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l584
								}
								position++
								goto l583
							l584:
								position, tokenIndex, depth = position583, tokenIndex583, depth583
								if buffer[position] != rune('R') {
									goto l577
								}
								position++
							}
						l583:
							if !rules[ruleskip]() {
								goto l577
							}
							depth--
							add(ruleSTR, position578)
						}
						goto l576
					l577:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position586 := position
							depth++
							{
								position587, tokenIndex587, depth587 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l588
								}
								position++
								goto l587
							l588:
								position, tokenIndex, depth = position587, tokenIndex587, depth587
								if buffer[position] != rune('L') {
									goto l585
								}
								position++
							}
						l587:
							{
								position589, tokenIndex589, depth589 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l590
								}
								position++
								goto l589
							l590:
								position, tokenIndex, depth = position589, tokenIndex589, depth589
								if buffer[position] != rune('A') {
									goto l585
								}
								position++
							}
						l589:
							{
								position591, tokenIndex591, depth591 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l592
								}
								position++
								goto l591
							l592:
								position, tokenIndex, depth = position591, tokenIndex591, depth591
								if buffer[position] != rune('N') {
									goto l585
								}
								position++
							}
						l591:
							{
								position593, tokenIndex593, depth593 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l594
								}
								position++
								goto l593
							l594:
								position, tokenIndex, depth = position593, tokenIndex593, depth593
								if buffer[position] != rune('G') {
									goto l585
								}
								position++
							}
						l593:
							if !rules[ruleskip]() {
								goto l585
							}
							depth--
							add(ruleLANG, position586)
						}
						goto l576
					l585:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position596 := position
							depth++
							{
								position597, tokenIndex597, depth597 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l598
								}
								position++
								goto l597
							l598:
								position, tokenIndex, depth = position597, tokenIndex597, depth597
								if buffer[position] != rune('D') {
									goto l595
								}
								position++
							}
						l597:
							{
								position599, tokenIndex599, depth599 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l600
								}
								position++
								goto l599
							l600:
								position, tokenIndex, depth = position599, tokenIndex599, depth599
								if buffer[position] != rune('A') {
									goto l595
								}
								position++
							}
						l599:
							{
								position601, tokenIndex601, depth601 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l602
								}
								position++
								goto l601
							l602:
								position, tokenIndex, depth = position601, tokenIndex601, depth601
								if buffer[position] != rune('T') {
									goto l595
								}
								position++
							}
						l601:
							{
								position603, tokenIndex603, depth603 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l604
								}
								position++
								goto l603
							l604:
								position, tokenIndex, depth = position603, tokenIndex603, depth603
								if buffer[position] != rune('A') {
									goto l595
								}
								position++
							}
						l603:
							{
								position605, tokenIndex605, depth605 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l606
								}
								position++
								goto l605
							l606:
								position, tokenIndex, depth = position605, tokenIndex605, depth605
								if buffer[position] != rune('T') {
									goto l595
								}
								position++
							}
						l605:
							{
								position607, tokenIndex607, depth607 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l608
								}
								position++
								goto l607
							l608:
								position, tokenIndex, depth = position607, tokenIndex607, depth607
								if buffer[position] != rune('Y') {
									goto l595
								}
								position++
							}
						l607:
							{
								position609, tokenIndex609, depth609 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l610
								}
								position++
								goto l609
							l610:
								position, tokenIndex, depth = position609, tokenIndex609, depth609
								if buffer[position] != rune('P') {
									goto l595
								}
								position++
							}
						l609:
							{
								position611, tokenIndex611, depth611 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l612
								}
								position++
								goto l611
							l612:
								position, tokenIndex, depth = position611, tokenIndex611, depth611
								if buffer[position] != rune('E') {
									goto l595
								}
								position++
							}
						l611:
							if !rules[ruleskip]() {
								goto l595
							}
							depth--
							add(ruleDATATYPE, position596)
						}
						goto l576
					l595:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position614 := position
							depth++
							{
								position615, tokenIndex615, depth615 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l616
								}
								position++
								goto l615
							l616:
								position, tokenIndex, depth = position615, tokenIndex615, depth615
								if buffer[position] != rune('I') {
									goto l613
								}
								position++
							}
						l615:
							{
								position617, tokenIndex617, depth617 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l618
								}
								position++
								goto l617
							l618:
								position, tokenIndex, depth = position617, tokenIndex617, depth617
								if buffer[position] != rune('R') {
									goto l613
								}
								position++
							}
						l617:
							{
								position619, tokenIndex619, depth619 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l620
								}
								position++
								goto l619
							l620:
								position, tokenIndex, depth = position619, tokenIndex619, depth619
								if buffer[position] != rune('I') {
									goto l613
								}
								position++
							}
						l619:
							if !rules[ruleskip]() {
								goto l613
							}
							depth--
							add(ruleIRI, position614)
						}
						goto l576
					l613:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position622 := position
							depth++
							{
								position623, tokenIndex623, depth623 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l624
								}
								position++
								goto l623
							l624:
								position, tokenIndex, depth = position623, tokenIndex623, depth623
								if buffer[position] != rune('U') {
									goto l621
								}
								position++
							}
						l623:
							{
								position625, tokenIndex625, depth625 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l626
								}
								position++
								goto l625
							l626:
								position, tokenIndex, depth = position625, tokenIndex625, depth625
								if buffer[position] != rune('R') {
									goto l621
								}
								position++
							}
						l625:
							{
								position627, tokenIndex627, depth627 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l628
								}
								position++
								goto l627
							l628:
								position, tokenIndex, depth = position627, tokenIndex627, depth627
								if buffer[position] != rune('I') {
									goto l621
								}
								position++
							}
						l627:
							if !rules[ruleskip]() {
								goto l621
							}
							depth--
							add(ruleURI, position622)
						}
						goto l576
					l621:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position630 := position
							depth++
							{
								position631, tokenIndex631, depth631 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l632
								}
								position++
								goto l631
							l632:
								position, tokenIndex, depth = position631, tokenIndex631, depth631
								if buffer[position] != rune('S') {
									goto l629
								}
								position++
							}
						l631:
							{
								position633, tokenIndex633, depth633 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l634
								}
								position++
								goto l633
							l634:
								position, tokenIndex, depth = position633, tokenIndex633, depth633
								if buffer[position] != rune('T') {
									goto l629
								}
								position++
							}
						l633:
							{
								position635, tokenIndex635, depth635 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l636
								}
								position++
								goto l635
							l636:
								position, tokenIndex, depth = position635, tokenIndex635, depth635
								if buffer[position] != rune('R') {
									goto l629
								}
								position++
							}
						l635:
							{
								position637, tokenIndex637, depth637 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l638
								}
								position++
								goto l637
							l638:
								position, tokenIndex, depth = position637, tokenIndex637, depth637
								if buffer[position] != rune('L') {
									goto l629
								}
								position++
							}
						l637:
							{
								position639, tokenIndex639, depth639 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l640
								}
								position++
								goto l639
							l640:
								position, tokenIndex, depth = position639, tokenIndex639, depth639
								if buffer[position] != rune('E') {
									goto l629
								}
								position++
							}
						l639:
							{
								position641, tokenIndex641, depth641 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l642
								}
								position++
								goto l641
							l642:
								position, tokenIndex, depth = position641, tokenIndex641, depth641
								if buffer[position] != rune('N') {
									goto l629
								}
								position++
							}
						l641:
							if !rules[ruleskip]() {
								goto l629
							}
							depth--
							add(ruleSTRLEN, position630)
						}
						goto l576
					l629:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position644 := position
							depth++
							{
								position645, tokenIndex645, depth645 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l646
								}
								position++
								goto l645
							l646:
								position, tokenIndex, depth = position645, tokenIndex645, depth645
								if buffer[position] != rune('M') {
									goto l643
								}
								position++
							}
						l645:
							{
								position647, tokenIndex647, depth647 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l648
								}
								position++
								goto l647
							l648:
								position, tokenIndex, depth = position647, tokenIndex647, depth647
								if buffer[position] != rune('O') {
									goto l643
								}
								position++
							}
						l647:
							{
								position649, tokenIndex649, depth649 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l650
								}
								position++
								goto l649
							l650:
								position, tokenIndex, depth = position649, tokenIndex649, depth649
								if buffer[position] != rune('N') {
									goto l643
								}
								position++
							}
						l649:
							{
								position651, tokenIndex651, depth651 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l652
								}
								position++
								goto l651
							l652:
								position, tokenIndex, depth = position651, tokenIndex651, depth651
								if buffer[position] != rune('T') {
									goto l643
								}
								position++
							}
						l651:
							{
								position653, tokenIndex653, depth653 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l654
								}
								position++
								goto l653
							l654:
								position, tokenIndex, depth = position653, tokenIndex653, depth653
								if buffer[position] != rune('H') {
									goto l643
								}
								position++
							}
						l653:
							if !rules[ruleskip]() {
								goto l643
							}
							depth--
							add(ruleMONTH, position644)
						}
						goto l576
					l643:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position656 := position
							depth++
							{
								position657, tokenIndex657, depth657 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l658
								}
								position++
								goto l657
							l658:
								position, tokenIndex, depth = position657, tokenIndex657, depth657
								if buffer[position] != rune('M') {
									goto l655
								}
								position++
							}
						l657:
							{
								position659, tokenIndex659, depth659 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l660
								}
								position++
								goto l659
							l660:
								position, tokenIndex, depth = position659, tokenIndex659, depth659
								if buffer[position] != rune('I') {
									goto l655
								}
								position++
							}
						l659:
							{
								position661, tokenIndex661, depth661 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l662
								}
								position++
								goto l661
							l662:
								position, tokenIndex, depth = position661, tokenIndex661, depth661
								if buffer[position] != rune('N') {
									goto l655
								}
								position++
							}
						l661:
							{
								position663, tokenIndex663, depth663 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l664
								}
								position++
								goto l663
							l664:
								position, tokenIndex, depth = position663, tokenIndex663, depth663
								if buffer[position] != rune('U') {
									goto l655
								}
								position++
							}
						l663:
							{
								position665, tokenIndex665, depth665 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l666
								}
								position++
								goto l665
							l666:
								position, tokenIndex, depth = position665, tokenIndex665, depth665
								if buffer[position] != rune('T') {
									goto l655
								}
								position++
							}
						l665:
							{
								position667, tokenIndex667, depth667 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l668
								}
								position++
								goto l667
							l668:
								position, tokenIndex, depth = position667, tokenIndex667, depth667
								if buffer[position] != rune('E') {
									goto l655
								}
								position++
							}
						l667:
							{
								position669, tokenIndex669, depth669 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l670
								}
								position++
								goto l669
							l670:
								position, tokenIndex, depth = position669, tokenIndex669, depth669
								if buffer[position] != rune('S') {
									goto l655
								}
								position++
							}
						l669:
							if !rules[ruleskip]() {
								goto l655
							}
							depth--
							add(ruleMINUTES, position656)
						}
						goto l576
					l655:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position672 := position
							depth++
							{
								position673, tokenIndex673, depth673 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l674
								}
								position++
								goto l673
							l674:
								position, tokenIndex, depth = position673, tokenIndex673, depth673
								if buffer[position] != rune('S') {
									goto l671
								}
								position++
							}
						l673:
							{
								position675, tokenIndex675, depth675 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l676
								}
								position++
								goto l675
							l676:
								position, tokenIndex, depth = position675, tokenIndex675, depth675
								if buffer[position] != rune('E') {
									goto l671
								}
								position++
							}
						l675:
							{
								position677, tokenIndex677, depth677 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l678
								}
								position++
								goto l677
							l678:
								position, tokenIndex, depth = position677, tokenIndex677, depth677
								if buffer[position] != rune('C') {
									goto l671
								}
								position++
							}
						l677:
							{
								position679, tokenIndex679, depth679 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l680
								}
								position++
								goto l679
							l680:
								position, tokenIndex, depth = position679, tokenIndex679, depth679
								if buffer[position] != rune('O') {
									goto l671
								}
								position++
							}
						l679:
							{
								position681, tokenIndex681, depth681 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l682
								}
								position++
								goto l681
							l682:
								position, tokenIndex, depth = position681, tokenIndex681, depth681
								if buffer[position] != rune('N') {
									goto l671
								}
								position++
							}
						l681:
							{
								position683, tokenIndex683, depth683 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l684
								}
								position++
								goto l683
							l684:
								position, tokenIndex, depth = position683, tokenIndex683, depth683
								if buffer[position] != rune('D') {
									goto l671
								}
								position++
							}
						l683:
							{
								position685, tokenIndex685, depth685 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l686
								}
								position++
								goto l685
							l686:
								position, tokenIndex, depth = position685, tokenIndex685, depth685
								if buffer[position] != rune('S') {
									goto l671
								}
								position++
							}
						l685:
							if !rules[ruleskip]() {
								goto l671
							}
							depth--
							add(ruleSECONDS, position672)
						}
						goto l576
					l671:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position688 := position
							depth++
							{
								position689, tokenIndex689, depth689 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l690
								}
								position++
								goto l689
							l690:
								position, tokenIndex, depth = position689, tokenIndex689, depth689
								if buffer[position] != rune('T') {
									goto l687
								}
								position++
							}
						l689:
							{
								position691, tokenIndex691, depth691 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l692
								}
								position++
								goto l691
							l692:
								position, tokenIndex, depth = position691, tokenIndex691, depth691
								if buffer[position] != rune('I') {
									goto l687
								}
								position++
							}
						l691:
							{
								position693, tokenIndex693, depth693 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l694
								}
								position++
								goto l693
							l694:
								position, tokenIndex, depth = position693, tokenIndex693, depth693
								if buffer[position] != rune('M') {
									goto l687
								}
								position++
							}
						l693:
							{
								position695, tokenIndex695, depth695 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l696
								}
								position++
								goto l695
							l696:
								position, tokenIndex, depth = position695, tokenIndex695, depth695
								if buffer[position] != rune('E') {
									goto l687
								}
								position++
							}
						l695:
							{
								position697, tokenIndex697, depth697 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l698
								}
								position++
								goto l697
							l698:
								position, tokenIndex, depth = position697, tokenIndex697, depth697
								if buffer[position] != rune('Z') {
									goto l687
								}
								position++
							}
						l697:
							{
								position699, tokenIndex699, depth699 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l700
								}
								position++
								goto l699
							l700:
								position, tokenIndex, depth = position699, tokenIndex699, depth699
								if buffer[position] != rune('O') {
									goto l687
								}
								position++
							}
						l699:
							{
								position701, tokenIndex701, depth701 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l702
								}
								position++
								goto l701
							l702:
								position, tokenIndex, depth = position701, tokenIndex701, depth701
								if buffer[position] != rune('N') {
									goto l687
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
									goto l687
								}
								position++
							}
						l703:
							if !rules[ruleskip]() {
								goto l687
							}
							depth--
							add(ruleTIMEZONE, position688)
						}
						goto l576
					l687:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position706 := position
							depth++
							{
								position707, tokenIndex707, depth707 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l708
								}
								position++
								goto l707
							l708:
								position, tokenIndex, depth = position707, tokenIndex707, depth707
								if buffer[position] != rune('S') {
									goto l705
								}
								position++
							}
						l707:
							{
								position709, tokenIndex709, depth709 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l710
								}
								position++
								goto l709
							l710:
								position, tokenIndex, depth = position709, tokenIndex709, depth709
								if buffer[position] != rune('H') {
									goto l705
								}
								position++
							}
						l709:
							{
								position711, tokenIndex711, depth711 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l712
								}
								position++
								goto l711
							l712:
								position, tokenIndex, depth = position711, tokenIndex711, depth711
								if buffer[position] != rune('A') {
									goto l705
								}
								position++
							}
						l711:
							if buffer[position] != rune('1') {
								goto l705
							}
							position++
							if !rules[ruleskip]() {
								goto l705
							}
							depth--
							add(ruleSHA1, position706)
						}
						goto l576
					l705:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position714 := position
							depth++
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
									goto l713
								}
								position++
							}
						l715:
							{
								position717, tokenIndex717, depth717 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l718
								}
								position++
								goto l717
							l718:
								position, tokenIndex, depth = position717, tokenIndex717, depth717
								if buffer[position] != rune('H') {
									goto l713
								}
								position++
							}
						l717:
							{
								position719, tokenIndex719, depth719 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l720
								}
								position++
								goto l719
							l720:
								position, tokenIndex, depth = position719, tokenIndex719, depth719
								if buffer[position] != rune('A') {
									goto l713
								}
								position++
							}
						l719:
							if buffer[position] != rune('2') {
								goto l713
							}
							position++
							if buffer[position] != rune('5') {
								goto l713
							}
							position++
							if buffer[position] != rune('6') {
								goto l713
							}
							position++
							if !rules[ruleskip]() {
								goto l713
							}
							depth--
							add(ruleSHA256, position714)
						}
						goto l576
					l713:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position722 := position
							depth++
							{
								position723, tokenIndex723, depth723 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l724
								}
								position++
								goto l723
							l724:
								position, tokenIndex, depth = position723, tokenIndex723, depth723
								if buffer[position] != rune('S') {
									goto l721
								}
								position++
							}
						l723:
							{
								position725, tokenIndex725, depth725 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l726
								}
								position++
								goto l725
							l726:
								position, tokenIndex, depth = position725, tokenIndex725, depth725
								if buffer[position] != rune('H') {
									goto l721
								}
								position++
							}
						l725:
							{
								position727, tokenIndex727, depth727 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l728
								}
								position++
								goto l727
							l728:
								position, tokenIndex, depth = position727, tokenIndex727, depth727
								if buffer[position] != rune('A') {
									goto l721
								}
								position++
							}
						l727:
							if buffer[position] != rune('3') {
								goto l721
							}
							position++
							if buffer[position] != rune('8') {
								goto l721
							}
							position++
							if buffer[position] != rune('4') {
								goto l721
							}
							position++
							if !rules[ruleskip]() {
								goto l721
							}
							depth--
							add(ruleSHA384, position722)
						}
						goto l576
					l721:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position730 := position
							depth++
							{
								position731, tokenIndex731, depth731 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l732
								}
								position++
								goto l731
							l732:
								position, tokenIndex, depth = position731, tokenIndex731, depth731
								if buffer[position] != rune('I') {
									goto l729
								}
								position++
							}
						l731:
							{
								position733, tokenIndex733, depth733 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l734
								}
								position++
								goto l733
							l734:
								position, tokenIndex, depth = position733, tokenIndex733, depth733
								if buffer[position] != rune('S') {
									goto l729
								}
								position++
							}
						l733:
							{
								position735, tokenIndex735, depth735 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l736
								}
								position++
								goto l735
							l736:
								position, tokenIndex, depth = position735, tokenIndex735, depth735
								if buffer[position] != rune('I') {
									goto l729
								}
								position++
							}
						l735:
							{
								position737, tokenIndex737, depth737 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l738
								}
								position++
								goto l737
							l738:
								position, tokenIndex, depth = position737, tokenIndex737, depth737
								if buffer[position] != rune('R') {
									goto l729
								}
								position++
							}
						l737:
							{
								position739, tokenIndex739, depth739 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l740
								}
								position++
								goto l739
							l740:
								position, tokenIndex, depth = position739, tokenIndex739, depth739
								if buffer[position] != rune('I') {
									goto l729
								}
								position++
							}
						l739:
							if !rules[ruleskip]() {
								goto l729
							}
							depth--
							add(ruleISIRI, position730)
						}
						goto l576
					l729:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position742 := position
							depth++
							{
								position743, tokenIndex743, depth743 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l744
								}
								position++
								goto l743
							l744:
								position, tokenIndex, depth = position743, tokenIndex743, depth743
								if buffer[position] != rune('I') {
									goto l741
								}
								position++
							}
						l743:
							{
								position745, tokenIndex745, depth745 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l746
								}
								position++
								goto l745
							l746:
								position, tokenIndex, depth = position745, tokenIndex745, depth745
								if buffer[position] != rune('S') {
									goto l741
								}
								position++
							}
						l745:
							{
								position747, tokenIndex747, depth747 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l748
								}
								position++
								goto l747
							l748:
								position, tokenIndex, depth = position747, tokenIndex747, depth747
								if buffer[position] != rune('U') {
									goto l741
								}
								position++
							}
						l747:
							{
								position749, tokenIndex749, depth749 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l750
								}
								position++
								goto l749
							l750:
								position, tokenIndex, depth = position749, tokenIndex749, depth749
								if buffer[position] != rune('R') {
									goto l741
								}
								position++
							}
						l749:
							{
								position751, tokenIndex751, depth751 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l752
								}
								position++
								goto l751
							l752:
								position, tokenIndex, depth = position751, tokenIndex751, depth751
								if buffer[position] != rune('I') {
									goto l741
								}
								position++
							}
						l751:
							if !rules[ruleskip]() {
								goto l741
							}
							depth--
							add(ruleISURI, position742)
						}
						goto l576
					l741:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position754 := position
							depth++
							{
								position755, tokenIndex755, depth755 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l756
								}
								position++
								goto l755
							l756:
								position, tokenIndex, depth = position755, tokenIndex755, depth755
								if buffer[position] != rune('I') {
									goto l753
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
									goto l753
								}
								position++
							}
						l757:
							{
								position759, tokenIndex759, depth759 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l760
								}
								position++
								goto l759
							l760:
								position, tokenIndex, depth = position759, tokenIndex759, depth759
								if buffer[position] != rune('B') {
									goto l753
								}
								position++
							}
						l759:
							{
								position761, tokenIndex761, depth761 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l762
								}
								position++
								goto l761
							l762:
								position, tokenIndex, depth = position761, tokenIndex761, depth761
								if buffer[position] != rune('L') {
									goto l753
								}
								position++
							}
						l761:
							{
								position763, tokenIndex763, depth763 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l764
								}
								position++
								goto l763
							l764:
								position, tokenIndex, depth = position763, tokenIndex763, depth763
								if buffer[position] != rune('A') {
									goto l753
								}
								position++
							}
						l763:
							{
								position765, tokenIndex765, depth765 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l766
								}
								position++
								goto l765
							l766:
								position, tokenIndex, depth = position765, tokenIndex765, depth765
								if buffer[position] != rune('N') {
									goto l753
								}
								position++
							}
						l765:
							{
								position767, tokenIndex767, depth767 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l768
								}
								position++
								goto l767
							l768:
								position, tokenIndex, depth = position767, tokenIndex767, depth767
								if buffer[position] != rune('K') {
									goto l753
								}
								position++
							}
						l767:
							if !rules[ruleskip]() {
								goto l753
							}
							depth--
							add(ruleISBLANK, position754)
						}
						goto l576
					l753:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							position770 := position
							depth++
							{
								position771, tokenIndex771, depth771 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l772
								}
								position++
								goto l771
							l772:
								position, tokenIndex, depth = position771, tokenIndex771, depth771
								if buffer[position] != rune('I') {
									goto l769
								}
								position++
							}
						l771:
							{
								position773, tokenIndex773, depth773 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l774
								}
								position++
								goto l773
							l774:
								position, tokenIndex, depth = position773, tokenIndex773, depth773
								if buffer[position] != rune('S') {
									goto l769
								}
								position++
							}
						l773:
							{
								position775, tokenIndex775, depth775 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l776
								}
								position++
								goto l775
							l776:
								position, tokenIndex, depth = position775, tokenIndex775, depth775
								if buffer[position] != rune('L') {
									goto l769
								}
								position++
							}
						l775:
							{
								position777, tokenIndex777, depth777 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l778
								}
								position++
								goto l777
							l778:
								position, tokenIndex, depth = position777, tokenIndex777, depth777
								if buffer[position] != rune('I') {
									goto l769
								}
								position++
							}
						l777:
							{
								position779, tokenIndex779, depth779 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l780
								}
								position++
								goto l779
							l780:
								position, tokenIndex, depth = position779, tokenIndex779, depth779
								if buffer[position] != rune('T') {
									goto l769
								}
								position++
							}
						l779:
							{
								position781, tokenIndex781, depth781 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l782
								}
								position++
								goto l781
							l782:
								position, tokenIndex, depth = position781, tokenIndex781, depth781
								if buffer[position] != rune('E') {
									goto l769
								}
								position++
							}
						l781:
							{
								position783, tokenIndex783, depth783 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l784
								}
								position++
								goto l783
							l784:
								position, tokenIndex, depth = position783, tokenIndex783, depth783
								if buffer[position] != rune('R') {
									goto l769
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
									goto l769
								}
								position++
							}
						l785:
							{
								position787, tokenIndex787, depth787 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l788
								}
								position++
								goto l787
							l788:
								position, tokenIndex, depth = position787, tokenIndex787, depth787
								if buffer[position] != rune('L') {
									goto l769
								}
								position++
							}
						l787:
							if !rules[ruleskip]() {
								goto l769
							}
							depth--
							add(ruleISLITERAL, position770)
						}
						goto l576
					l769:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position790 := position
									depth++
									{
										position791, tokenIndex791, depth791 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l792
										}
										position++
										goto l791
									l792:
										position, tokenIndex, depth = position791, tokenIndex791, depth791
										if buffer[position] != rune('I') {
											goto l575
										}
										position++
									}
								l791:
									{
										position793, tokenIndex793, depth793 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l794
										}
										position++
										goto l793
									l794:
										position, tokenIndex, depth = position793, tokenIndex793, depth793
										if buffer[position] != rune('S') {
											goto l575
										}
										position++
									}
								l793:
									{
										position795, tokenIndex795, depth795 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l796
										}
										position++
										goto l795
									l796:
										position, tokenIndex, depth = position795, tokenIndex795, depth795
										if buffer[position] != rune('N') {
											goto l575
										}
										position++
									}
								l795:
									{
										position797, tokenIndex797, depth797 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l798
										}
										position++
										goto l797
									l798:
										position, tokenIndex, depth = position797, tokenIndex797, depth797
										if buffer[position] != rune('U') {
											goto l575
										}
										position++
									}
								l797:
									{
										position799, tokenIndex799, depth799 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l800
										}
										position++
										goto l799
									l800:
										position, tokenIndex, depth = position799, tokenIndex799, depth799
										if buffer[position] != rune('M') {
											goto l575
										}
										position++
									}
								l799:
									{
										position801, tokenIndex801, depth801 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l802
										}
										position++
										goto l801
									l802:
										position, tokenIndex, depth = position801, tokenIndex801, depth801
										if buffer[position] != rune('E') {
											goto l575
										}
										position++
									}
								l801:
									{
										position803, tokenIndex803, depth803 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l804
										}
										position++
										goto l803
									l804:
										position, tokenIndex, depth = position803, tokenIndex803, depth803
										if buffer[position] != rune('R') {
											goto l575
										}
										position++
									}
								l803:
									{
										position805, tokenIndex805, depth805 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l806
										}
										position++
										goto l805
									l806:
										position, tokenIndex, depth = position805, tokenIndex805, depth805
										if buffer[position] != rune('I') {
											goto l575
										}
										position++
									}
								l805:
									{
										position807, tokenIndex807, depth807 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l808
										}
										position++
										goto l807
									l808:
										position, tokenIndex, depth = position807, tokenIndex807, depth807
										if buffer[position] != rune('C') {
											goto l575
										}
										position++
									}
								l807:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleISNUMERIC, position790)
								}
								break
							case 'S', 's':
								{
									position809 := position
									depth++
									{
										position810, tokenIndex810, depth810 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l811
										}
										position++
										goto l810
									l811:
										position, tokenIndex, depth = position810, tokenIndex810, depth810
										if buffer[position] != rune('S') {
											goto l575
										}
										position++
									}
								l810:
									{
										position812, tokenIndex812, depth812 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l813
										}
										position++
										goto l812
									l813:
										position, tokenIndex, depth = position812, tokenIndex812, depth812
										if buffer[position] != rune('H') {
											goto l575
										}
										position++
									}
								l812:
									{
										position814, tokenIndex814, depth814 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l815
										}
										position++
										goto l814
									l815:
										position, tokenIndex, depth = position814, tokenIndex814, depth814
										if buffer[position] != rune('A') {
											goto l575
										}
										position++
									}
								l814:
									if buffer[position] != rune('5') {
										goto l575
									}
									position++
									if buffer[position] != rune('1') {
										goto l575
									}
									position++
									if buffer[position] != rune('2') {
										goto l575
									}
									position++
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleSHA512, position809)
								}
								break
							case 'M', 'm':
								{
									position816 := position
									depth++
									{
										position817, tokenIndex817, depth817 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l818
										}
										position++
										goto l817
									l818:
										position, tokenIndex, depth = position817, tokenIndex817, depth817
										if buffer[position] != rune('M') {
											goto l575
										}
										position++
									}
								l817:
									{
										position819, tokenIndex819, depth819 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l820
										}
										position++
										goto l819
									l820:
										position, tokenIndex, depth = position819, tokenIndex819, depth819
										if buffer[position] != rune('D') {
											goto l575
										}
										position++
									}
								l819:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleMD5, position816)
								}
								break
							case 'T', 't':
								{
									position821 := position
									depth++
									{
										position822, tokenIndex822, depth822 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l823
										}
										position++
										goto l822
									l823:
										position, tokenIndex, depth = position822, tokenIndex822, depth822
										if buffer[position] != rune('T') {
											goto l575
										}
										position++
									}
								l822:
									{
										position824, tokenIndex824, depth824 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l825
										}
										position++
										goto l824
									l825:
										position, tokenIndex, depth = position824, tokenIndex824, depth824
										if buffer[position] != rune('Z') {
											goto l575
										}
										position++
									}
								l824:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleTZ, position821)
								}
								break
							case 'H', 'h':
								{
									position826 := position
									depth++
									{
										position827, tokenIndex827, depth827 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l828
										}
										position++
										goto l827
									l828:
										position, tokenIndex, depth = position827, tokenIndex827, depth827
										if buffer[position] != rune('H') {
											goto l575
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
											goto l575
										}
										position++
									}
								l829:
									{
										position831, tokenIndex831, depth831 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l832
										}
										position++
										goto l831
									l832:
										position, tokenIndex, depth = position831, tokenIndex831, depth831
										if buffer[position] != rune('U') {
											goto l575
										}
										position++
									}
								l831:
									{
										position833, tokenIndex833, depth833 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l834
										}
										position++
										goto l833
									l834:
										position, tokenIndex, depth = position833, tokenIndex833, depth833
										if buffer[position] != rune('R') {
											goto l575
										}
										position++
									}
								l833:
									{
										position835, tokenIndex835, depth835 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l836
										}
										position++
										goto l835
									l836:
										position, tokenIndex, depth = position835, tokenIndex835, depth835
										if buffer[position] != rune('S') {
											goto l575
										}
										position++
									}
								l835:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleHOURS, position826)
								}
								break
							case 'D', 'd':
								{
									position837 := position
									depth++
									{
										position838, tokenIndex838, depth838 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l839
										}
										position++
										goto l838
									l839:
										position, tokenIndex, depth = position838, tokenIndex838, depth838
										if buffer[position] != rune('D') {
											goto l575
										}
										position++
									}
								l838:
									{
										position840, tokenIndex840, depth840 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l841
										}
										position++
										goto l840
									l841:
										position, tokenIndex, depth = position840, tokenIndex840, depth840
										if buffer[position] != rune('A') {
											goto l575
										}
										position++
									}
								l840:
									{
										position842, tokenIndex842, depth842 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l843
										}
										position++
										goto l842
									l843:
										position, tokenIndex, depth = position842, tokenIndex842, depth842
										if buffer[position] != rune('Y') {
											goto l575
										}
										position++
									}
								l842:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleDAY, position837)
								}
								break
							case 'Y', 'y':
								{
									position844 := position
									depth++
									{
										position845, tokenIndex845, depth845 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l846
										}
										position++
										goto l845
									l846:
										position, tokenIndex, depth = position845, tokenIndex845, depth845
										if buffer[position] != rune('Y') {
											goto l575
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
											goto l575
										}
										position++
									}
								l847:
									{
										position849, tokenIndex849, depth849 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l850
										}
										position++
										goto l849
									l850:
										position, tokenIndex, depth = position849, tokenIndex849, depth849
										if buffer[position] != rune('A') {
											goto l575
										}
										position++
									}
								l849:
									{
										position851, tokenIndex851, depth851 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l852
										}
										position++
										goto l851
									l852:
										position, tokenIndex, depth = position851, tokenIndex851, depth851
										if buffer[position] != rune('R') {
											goto l575
										}
										position++
									}
								l851:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleYEAR, position844)
								}
								break
							case 'E', 'e':
								{
									position853 := position
									depth++
									{
										position854, tokenIndex854, depth854 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l855
										}
										position++
										goto l854
									l855:
										position, tokenIndex, depth = position854, tokenIndex854, depth854
										if buffer[position] != rune('E') {
											goto l575
										}
										position++
									}
								l854:
									{
										position856, tokenIndex856, depth856 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l857
										}
										position++
										goto l856
									l857:
										position, tokenIndex, depth = position856, tokenIndex856, depth856
										if buffer[position] != rune('N') {
											goto l575
										}
										position++
									}
								l856:
									{
										position858, tokenIndex858, depth858 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l859
										}
										position++
										goto l858
									l859:
										position, tokenIndex, depth = position858, tokenIndex858, depth858
										if buffer[position] != rune('C') {
											goto l575
										}
										position++
									}
								l858:
									{
										position860, tokenIndex860, depth860 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l861
										}
										position++
										goto l860
									l861:
										position, tokenIndex, depth = position860, tokenIndex860, depth860
										if buffer[position] != rune('O') {
											goto l575
										}
										position++
									}
								l860:
									{
										position862, tokenIndex862, depth862 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l863
										}
										position++
										goto l862
									l863:
										position, tokenIndex, depth = position862, tokenIndex862, depth862
										if buffer[position] != rune('D') {
											goto l575
										}
										position++
									}
								l862:
									{
										position864, tokenIndex864, depth864 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l865
										}
										position++
										goto l864
									l865:
										position, tokenIndex, depth = position864, tokenIndex864, depth864
										if buffer[position] != rune('E') {
											goto l575
										}
										position++
									}
								l864:
									if buffer[position] != rune('_') {
										goto l575
									}
									position++
									{
										position866, tokenIndex866, depth866 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l867
										}
										position++
										goto l866
									l867:
										position, tokenIndex, depth = position866, tokenIndex866, depth866
										if buffer[position] != rune('F') {
											goto l575
										}
										position++
									}
								l866:
									{
										position868, tokenIndex868, depth868 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l869
										}
										position++
										goto l868
									l869:
										position, tokenIndex, depth = position868, tokenIndex868, depth868
										if buffer[position] != rune('O') {
											goto l575
										}
										position++
									}
								l868:
									{
										position870, tokenIndex870, depth870 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l871
										}
										position++
										goto l870
									l871:
										position, tokenIndex, depth = position870, tokenIndex870, depth870
										if buffer[position] != rune('R') {
											goto l575
										}
										position++
									}
								l870:
									if buffer[position] != rune('_') {
										goto l575
									}
									position++
									{
										position872, tokenIndex872, depth872 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l873
										}
										position++
										goto l872
									l873:
										position, tokenIndex, depth = position872, tokenIndex872, depth872
										if buffer[position] != rune('U') {
											goto l575
										}
										position++
									}
								l872:
									{
										position874, tokenIndex874, depth874 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l875
										}
										position++
										goto l874
									l875:
										position, tokenIndex, depth = position874, tokenIndex874, depth874
										if buffer[position] != rune('R') {
											goto l575
										}
										position++
									}
								l874:
									{
										position876, tokenIndex876, depth876 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l877
										}
										position++
										goto l876
									l877:
										position, tokenIndex, depth = position876, tokenIndex876, depth876
										if buffer[position] != rune('I') {
											goto l575
										}
										position++
									}
								l876:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleENCODEFORURI, position853)
								}
								break
							case 'L', 'l':
								{
									position878 := position
									depth++
									{
										position879, tokenIndex879, depth879 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l880
										}
										position++
										goto l879
									l880:
										position, tokenIndex, depth = position879, tokenIndex879, depth879
										if buffer[position] != rune('L') {
											goto l575
										}
										position++
									}
								l879:
									{
										position881, tokenIndex881, depth881 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l882
										}
										position++
										goto l881
									l882:
										position, tokenIndex, depth = position881, tokenIndex881, depth881
										if buffer[position] != rune('C') {
											goto l575
										}
										position++
									}
								l881:
									{
										position883, tokenIndex883, depth883 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l884
										}
										position++
										goto l883
									l884:
										position, tokenIndex, depth = position883, tokenIndex883, depth883
										if buffer[position] != rune('A') {
											goto l575
										}
										position++
									}
								l883:
									{
										position885, tokenIndex885, depth885 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l886
										}
										position++
										goto l885
									l886:
										position, tokenIndex, depth = position885, tokenIndex885, depth885
										if buffer[position] != rune('S') {
											goto l575
										}
										position++
									}
								l885:
									{
										position887, tokenIndex887, depth887 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l888
										}
										position++
										goto l887
									l888:
										position, tokenIndex, depth = position887, tokenIndex887, depth887
										if buffer[position] != rune('E') {
											goto l575
										}
										position++
									}
								l887:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleLCASE, position878)
								}
								break
							case 'U', 'u':
								{
									position889 := position
									depth++
									{
										position890, tokenIndex890, depth890 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l891
										}
										position++
										goto l890
									l891:
										position, tokenIndex, depth = position890, tokenIndex890, depth890
										if buffer[position] != rune('U') {
											goto l575
										}
										position++
									}
								l890:
									{
										position892, tokenIndex892, depth892 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l893
										}
										position++
										goto l892
									l893:
										position, tokenIndex, depth = position892, tokenIndex892, depth892
										if buffer[position] != rune('C') {
											goto l575
										}
										position++
									}
								l892:
									{
										position894, tokenIndex894, depth894 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l895
										}
										position++
										goto l894
									l895:
										position, tokenIndex, depth = position894, tokenIndex894, depth894
										if buffer[position] != rune('A') {
											goto l575
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
											goto l575
										}
										position++
									}
								l896:
									{
										position898, tokenIndex898, depth898 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l899
										}
										position++
										goto l898
									l899:
										position, tokenIndex, depth = position898, tokenIndex898, depth898
										if buffer[position] != rune('E') {
											goto l575
										}
										position++
									}
								l898:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleUCASE, position889)
								}
								break
							case 'F', 'f':
								{
									position900 := position
									depth++
									{
										position901, tokenIndex901, depth901 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l902
										}
										position++
										goto l901
									l902:
										position, tokenIndex, depth = position901, tokenIndex901, depth901
										if buffer[position] != rune('F') {
											goto l575
										}
										position++
									}
								l901:
									{
										position903, tokenIndex903, depth903 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l904
										}
										position++
										goto l903
									l904:
										position, tokenIndex, depth = position903, tokenIndex903, depth903
										if buffer[position] != rune('L') {
											goto l575
										}
										position++
									}
								l903:
									{
										position905, tokenIndex905, depth905 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l906
										}
										position++
										goto l905
									l906:
										position, tokenIndex, depth = position905, tokenIndex905, depth905
										if buffer[position] != rune('O') {
											goto l575
										}
										position++
									}
								l905:
									{
										position907, tokenIndex907, depth907 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l908
										}
										position++
										goto l907
									l908:
										position, tokenIndex, depth = position907, tokenIndex907, depth907
										if buffer[position] != rune('O') {
											goto l575
										}
										position++
									}
								l907:
									{
										position909, tokenIndex909, depth909 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l910
										}
										position++
										goto l909
									l910:
										position, tokenIndex, depth = position909, tokenIndex909, depth909
										if buffer[position] != rune('R') {
											goto l575
										}
										position++
									}
								l909:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleFLOOR, position900)
								}
								break
							case 'R', 'r':
								{
									position911 := position
									depth++
									{
										position912, tokenIndex912, depth912 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l913
										}
										position++
										goto l912
									l913:
										position, tokenIndex, depth = position912, tokenIndex912, depth912
										if buffer[position] != rune('R') {
											goto l575
										}
										position++
									}
								l912:
									{
										position914, tokenIndex914, depth914 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l915
										}
										position++
										goto l914
									l915:
										position, tokenIndex, depth = position914, tokenIndex914, depth914
										if buffer[position] != rune('O') {
											goto l575
										}
										position++
									}
								l914:
									{
										position916, tokenIndex916, depth916 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l917
										}
										position++
										goto l916
									l917:
										position, tokenIndex, depth = position916, tokenIndex916, depth916
										if buffer[position] != rune('U') {
											goto l575
										}
										position++
									}
								l916:
									{
										position918, tokenIndex918, depth918 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l919
										}
										position++
										goto l918
									l919:
										position, tokenIndex, depth = position918, tokenIndex918, depth918
										if buffer[position] != rune('N') {
											goto l575
										}
										position++
									}
								l918:
									{
										position920, tokenIndex920, depth920 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l921
										}
										position++
										goto l920
									l921:
										position, tokenIndex, depth = position920, tokenIndex920, depth920
										if buffer[position] != rune('D') {
											goto l575
										}
										position++
									}
								l920:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleROUND, position911)
								}
								break
							case 'C', 'c':
								{
									position922 := position
									depth++
									{
										position923, tokenIndex923, depth923 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l924
										}
										position++
										goto l923
									l924:
										position, tokenIndex, depth = position923, tokenIndex923, depth923
										if buffer[position] != rune('C') {
											goto l575
										}
										position++
									}
								l923:
									{
										position925, tokenIndex925, depth925 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l926
										}
										position++
										goto l925
									l926:
										position, tokenIndex, depth = position925, tokenIndex925, depth925
										if buffer[position] != rune('E') {
											goto l575
										}
										position++
									}
								l925:
									{
										position927, tokenIndex927, depth927 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l928
										}
										position++
										goto l927
									l928:
										position, tokenIndex, depth = position927, tokenIndex927, depth927
										if buffer[position] != rune('I') {
											goto l575
										}
										position++
									}
								l927:
									{
										position929, tokenIndex929, depth929 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l930
										}
										position++
										goto l929
									l930:
										position, tokenIndex, depth = position929, tokenIndex929, depth929
										if buffer[position] != rune('L') {
											goto l575
										}
										position++
									}
								l929:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleCEIL, position922)
								}
								break
							default:
								{
									position931 := position
									depth++
									{
										position932, tokenIndex932, depth932 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l933
										}
										position++
										goto l932
									l933:
										position, tokenIndex, depth = position932, tokenIndex932, depth932
										if buffer[position] != rune('A') {
											goto l575
										}
										position++
									}
								l932:
									{
										position934, tokenIndex934, depth934 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l935
										}
										position++
										goto l934
									l935:
										position, tokenIndex, depth = position934, tokenIndex934, depth934
										if buffer[position] != rune('B') {
											goto l575
										}
										position++
									}
								l934:
									{
										position936, tokenIndex936, depth936 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l937
										}
										position++
										goto l936
									l937:
										position, tokenIndex, depth = position936, tokenIndex936, depth936
										if buffer[position] != rune('S') {
											goto l575
										}
										position++
									}
								l936:
									if !rules[ruleskip]() {
										goto l575
									}
									depth--
									add(ruleABS, position931)
								}
								break
							}
						}

					}
				l576:
					if !rules[ruleLPAREN]() {
						goto l575
					}
					if !rules[ruleexpression]() {
						goto l575
					}
					if !rules[ruleRPAREN]() {
						goto l575
					}
					goto l574
				l575:
					position, tokenIndex, depth = position574, tokenIndex574, depth574
					{
						position939, tokenIndex939, depth939 := position, tokenIndex, depth
						{
							position941 := position
							depth++
							{
								position942, tokenIndex942, depth942 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l943
								}
								position++
								goto l942
							l943:
								position, tokenIndex, depth = position942, tokenIndex942, depth942
								if buffer[position] != rune('S') {
									goto l940
								}
								position++
							}
						l942:
							{
								position944, tokenIndex944, depth944 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l945
								}
								position++
								goto l944
							l945:
								position, tokenIndex, depth = position944, tokenIndex944, depth944
								if buffer[position] != rune('T') {
									goto l940
								}
								position++
							}
						l944:
							{
								position946, tokenIndex946, depth946 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l947
								}
								position++
								goto l946
							l947:
								position, tokenIndex, depth = position946, tokenIndex946, depth946
								if buffer[position] != rune('R') {
									goto l940
								}
								position++
							}
						l946:
							{
								position948, tokenIndex948, depth948 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l949
								}
								position++
								goto l948
							l949:
								position, tokenIndex, depth = position948, tokenIndex948, depth948
								if buffer[position] != rune('S') {
									goto l940
								}
								position++
							}
						l948:
							{
								position950, tokenIndex950, depth950 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l951
								}
								position++
								goto l950
							l951:
								position, tokenIndex, depth = position950, tokenIndex950, depth950
								if buffer[position] != rune('T') {
									goto l940
								}
								position++
							}
						l950:
							{
								position952, tokenIndex952, depth952 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l953
								}
								position++
								goto l952
							l953:
								position, tokenIndex, depth = position952, tokenIndex952, depth952
								if buffer[position] != rune('A') {
									goto l940
								}
								position++
							}
						l952:
							{
								position954, tokenIndex954, depth954 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l955
								}
								position++
								goto l954
							l955:
								position, tokenIndex, depth = position954, tokenIndex954, depth954
								if buffer[position] != rune('R') {
									goto l940
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
									goto l940
								}
								position++
							}
						l956:
							{
								position958, tokenIndex958, depth958 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l959
								}
								position++
								goto l958
							l959:
								position, tokenIndex, depth = position958, tokenIndex958, depth958
								if buffer[position] != rune('S') {
									goto l940
								}
								position++
							}
						l958:
							if !rules[ruleskip]() {
								goto l940
							}
							depth--
							add(ruleSTRSTARTS, position941)
						}
						goto l939
					l940:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						{
							position961 := position
							depth++
							{
								position962, tokenIndex962, depth962 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l963
								}
								position++
								goto l962
							l963:
								position, tokenIndex, depth = position962, tokenIndex962, depth962
								if buffer[position] != rune('S') {
									goto l960
								}
								position++
							}
						l962:
							{
								position964, tokenIndex964, depth964 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l965
								}
								position++
								goto l964
							l965:
								position, tokenIndex, depth = position964, tokenIndex964, depth964
								if buffer[position] != rune('T') {
									goto l960
								}
								position++
							}
						l964:
							{
								position966, tokenIndex966, depth966 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l967
								}
								position++
								goto l966
							l967:
								position, tokenIndex, depth = position966, tokenIndex966, depth966
								if buffer[position] != rune('R') {
									goto l960
								}
								position++
							}
						l966:
							{
								position968, tokenIndex968, depth968 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l969
								}
								position++
								goto l968
							l969:
								position, tokenIndex, depth = position968, tokenIndex968, depth968
								if buffer[position] != rune('E') {
									goto l960
								}
								position++
							}
						l968:
							{
								position970, tokenIndex970, depth970 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l971
								}
								position++
								goto l970
							l971:
								position, tokenIndex, depth = position970, tokenIndex970, depth970
								if buffer[position] != rune('N') {
									goto l960
								}
								position++
							}
						l970:
							{
								position972, tokenIndex972, depth972 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l973
								}
								position++
								goto l972
							l973:
								position, tokenIndex, depth = position972, tokenIndex972, depth972
								if buffer[position] != rune('D') {
									goto l960
								}
								position++
							}
						l972:
							{
								position974, tokenIndex974, depth974 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l975
								}
								position++
								goto l974
							l975:
								position, tokenIndex, depth = position974, tokenIndex974, depth974
								if buffer[position] != rune('S') {
									goto l960
								}
								position++
							}
						l974:
							if !rules[ruleskip]() {
								goto l960
							}
							depth--
							add(ruleSTRENDS, position961)
						}
						goto l939
					l960:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						{
							position977 := position
							depth++
							{
								position978, tokenIndex978, depth978 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l979
								}
								position++
								goto l978
							l979:
								position, tokenIndex, depth = position978, tokenIndex978, depth978
								if buffer[position] != rune('S') {
									goto l976
								}
								position++
							}
						l978:
							{
								position980, tokenIndex980, depth980 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l981
								}
								position++
								goto l980
							l981:
								position, tokenIndex, depth = position980, tokenIndex980, depth980
								if buffer[position] != rune('T') {
									goto l976
								}
								position++
							}
						l980:
							{
								position982, tokenIndex982, depth982 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l983
								}
								position++
								goto l982
							l983:
								position, tokenIndex, depth = position982, tokenIndex982, depth982
								if buffer[position] != rune('R') {
									goto l976
								}
								position++
							}
						l982:
							{
								position984, tokenIndex984, depth984 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l985
								}
								position++
								goto l984
							l985:
								position, tokenIndex, depth = position984, tokenIndex984, depth984
								if buffer[position] != rune('B') {
									goto l976
								}
								position++
							}
						l984:
							{
								position986, tokenIndex986, depth986 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l987
								}
								position++
								goto l986
							l987:
								position, tokenIndex, depth = position986, tokenIndex986, depth986
								if buffer[position] != rune('E') {
									goto l976
								}
								position++
							}
						l986:
							{
								position988, tokenIndex988, depth988 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l989
								}
								position++
								goto l988
							l989:
								position, tokenIndex, depth = position988, tokenIndex988, depth988
								if buffer[position] != rune('F') {
									goto l976
								}
								position++
							}
						l988:
							{
								position990, tokenIndex990, depth990 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l991
								}
								position++
								goto l990
							l991:
								position, tokenIndex, depth = position990, tokenIndex990, depth990
								if buffer[position] != rune('O') {
									goto l976
								}
								position++
							}
						l990:
							{
								position992, tokenIndex992, depth992 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l993
								}
								position++
								goto l992
							l993:
								position, tokenIndex, depth = position992, tokenIndex992, depth992
								if buffer[position] != rune('R') {
									goto l976
								}
								position++
							}
						l992:
							{
								position994, tokenIndex994, depth994 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l995
								}
								position++
								goto l994
							l995:
								position, tokenIndex, depth = position994, tokenIndex994, depth994
								if buffer[position] != rune('E') {
									goto l976
								}
								position++
							}
						l994:
							if !rules[ruleskip]() {
								goto l976
							}
							depth--
							add(ruleSTRBEFORE, position977)
						}
						goto l939
					l976:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						{
							position997 := position
							depth++
							{
								position998, tokenIndex998, depth998 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l999
								}
								position++
								goto l998
							l999:
								position, tokenIndex, depth = position998, tokenIndex998, depth998
								if buffer[position] != rune('S') {
									goto l996
								}
								position++
							}
						l998:
							{
								position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1001
								}
								position++
								goto l1000
							l1001:
								position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
								if buffer[position] != rune('T') {
									goto l996
								}
								position++
							}
						l1000:
							{
								position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1003
								}
								position++
								goto l1002
							l1003:
								position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
								if buffer[position] != rune('R') {
									goto l996
								}
								position++
							}
						l1002:
							{
								position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1005
								}
								position++
								goto l1004
							l1005:
								position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
								if buffer[position] != rune('A') {
									goto l996
								}
								position++
							}
						l1004:
							{
								position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1007
								}
								position++
								goto l1006
							l1007:
								position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
								if buffer[position] != rune('F') {
									goto l996
								}
								position++
							}
						l1006:
							{
								position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1009
								}
								position++
								goto l1008
							l1009:
								position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
								if buffer[position] != rune('T') {
									goto l996
								}
								position++
							}
						l1008:
							{
								position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1011
								}
								position++
								goto l1010
							l1011:
								position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
								if buffer[position] != rune('E') {
									goto l996
								}
								position++
							}
						l1010:
							{
								position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1013
								}
								position++
								goto l1012
							l1013:
								position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
								if buffer[position] != rune('R') {
									goto l996
								}
								position++
							}
						l1012:
							if !rules[ruleskip]() {
								goto l996
							}
							depth--
							add(ruleSTRAFTER, position997)
						}
						goto l939
					l996:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						{
							position1015 := position
							depth++
							{
								position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1017
								}
								position++
								goto l1016
							l1017:
								position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
								if buffer[position] != rune('S') {
									goto l1014
								}
								position++
							}
						l1016:
							{
								position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1019
								}
								position++
								goto l1018
							l1019:
								position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
								if buffer[position] != rune('T') {
									goto l1014
								}
								position++
							}
						l1018:
							{
								position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1021
								}
								position++
								goto l1020
							l1021:
								position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
								if buffer[position] != rune('R') {
									goto l1014
								}
								position++
							}
						l1020:
							{
								position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1023
								}
								position++
								goto l1022
							l1023:
								position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
								if buffer[position] != rune('L') {
									goto l1014
								}
								position++
							}
						l1022:
							{
								position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1025
								}
								position++
								goto l1024
							l1025:
								position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
								if buffer[position] != rune('A') {
									goto l1014
								}
								position++
							}
						l1024:
							{
								position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1027
								}
								position++
								goto l1026
							l1027:
								position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
								if buffer[position] != rune('N') {
									goto l1014
								}
								position++
							}
						l1026:
							{
								position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1029
								}
								position++
								goto l1028
							l1029:
								position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
								if buffer[position] != rune('G') {
									goto l1014
								}
								position++
							}
						l1028:
							if !rules[ruleskip]() {
								goto l1014
							}
							depth--
							add(ruleSTRLANG, position1015)
						}
						goto l939
					l1014:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						{
							position1031 := position
							depth++
							{
								position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1033
								}
								position++
								goto l1032
							l1033:
								position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
								if buffer[position] != rune('S') {
									goto l1030
								}
								position++
							}
						l1032:
							{
								position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1035
								}
								position++
								goto l1034
							l1035:
								position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
								if buffer[position] != rune('T') {
									goto l1030
								}
								position++
							}
						l1034:
							{
								position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1037
								}
								position++
								goto l1036
							l1037:
								position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
								if buffer[position] != rune('R') {
									goto l1030
								}
								position++
							}
						l1036:
							{
								position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1039
								}
								position++
								goto l1038
							l1039:
								position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
								if buffer[position] != rune('D') {
									goto l1030
								}
								position++
							}
						l1038:
							{
								position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1041
								}
								position++
								goto l1040
							l1041:
								position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
								if buffer[position] != rune('T') {
									goto l1030
								}
								position++
							}
						l1040:
							if !rules[ruleskip]() {
								goto l1030
							}
							depth--
							add(ruleSTRDT, position1031)
						}
						goto l939
					l1030:
						position, tokenIndex, depth = position939, tokenIndex939, depth939
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1043 := position
									depth++
									{
										position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1045
										}
										position++
										goto l1044
									l1045:
										position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
										if buffer[position] != rune('S') {
											goto l938
										}
										position++
									}
								l1044:
									{
										position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1047
										}
										position++
										goto l1046
									l1047:
										position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
										if buffer[position] != rune('A') {
											goto l938
										}
										position++
									}
								l1046:
									{
										position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1049
										}
										position++
										goto l1048
									l1049:
										position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
										if buffer[position] != rune('M') {
											goto l938
										}
										position++
									}
								l1048:
									{
										position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1051
										}
										position++
										goto l1050
									l1051:
										position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
										if buffer[position] != rune('E') {
											goto l938
										}
										position++
									}
								l1050:
									{
										position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1053
										}
										position++
										goto l1052
									l1053:
										position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
										if buffer[position] != rune('T') {
											goto l938
										}
										position++
									}
								l1052:
									{
										position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1055
										}
										position++
										goto l1054
									l1055:
										position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
										if buffer[position] != rune('E') {
											goto l938
										}
										position++
									}
								l1054:
									{
										position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1057
										}
										position++
										goto l1056
									l1057:
										position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
										if buffer[position] != rune('R') {
											goto l938
										}
										position++
									}
								l1056:
									{
										position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1059
										}
										position++
										goto l1058
									l1059:
										position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
										if buffer[position] != rune('M') {
											goto l938
										}
										position++
									}
								l1058:
									if !rules[ruleskip]() {
										goto l938
									}
									depth--
									add(ruleSAMETERM, position1043)
								}
								break
							case 'C', 'c':
								{
									position1060 := position
									depth++
									{
										position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1062
										}
										position++
										goto l1061
									l1062:
										position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
										if buffer[position] != rune('C') {
											goto l938
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
											goto l938
										}
										position++
									}
								l1063:
									{
										position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1066
										}
										position++
										goto l1065
									l1066:
										position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
										if buffer[position] != rune('N') {
											goto l938
										}
										position++
									}
								l1065:
									{
										position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1068
										}
										position++
										goto l1067
									l1068:
										position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
										if buffer[position] != rune('T') {
											goto l938
										}
										position++
									}
								l1067:
									{
										position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1070
										}
										position++
										goto l1069
									l1070:
										position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
										if buffer[position] != rune('A') {
											goto l938
										}
										position++
									}
								l1069:
									{
										position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1072
										}
										position++
										goto l1071
									l1072:
										position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
										if buffer[position] != rune('I') {
											goto l938
										}
										position++
									}
								l1071:
									{
										position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1074
										}
										position++
										goto l1073
									l1074:
										position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
										if buffer[position] != rune('N') {
											goto l938
										}
										position++
									}
								l1073:
									{
										position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1076
										}
										position++
										goto l1075
									l1076:
										position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
										if buffer[position] != rune('S') {
											goto l938
										}
										position++
									}
								l1075:
									if !rules[ruleskip]() {
										goto l938
									}
									depth--
									add(ruleCONTAINS, position1060)
								}
								break
							default:
								{
									position1077 := position
									depth++
									{
										position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1079
										}
										position++
										goto l1078
									l1079:
										position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
										if buffer[position] != rune('L') {
											goto l938
										}
										position++
									}
								l1078:
									{
										position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1081
										}
										position++
										goto l1080
									l1081:
										position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
										if buffer[position] != rune('A') {
											goto l938
										}
										position++
									}
								l1080:
									{
										position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1083
										}
										position++
										goto l1082
									l1083:
										position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
										if buffer[position] != rune('N') {
											goto l938
										}
										position++
									}
								l1082:
									{
										position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1085
										}
										position++
										goto l1084
									l1085:
										position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
										if buffer[position] != rune('G') {
											goto l938
										}
										position++
									}
								l1084:
									{
										position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1087
										}
										position++
										goto l1086
									l1087:
										position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
										if buffer[position] != rune('M') {
											goto l938
										}
										position++
									}
								l1086:
									{
										position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1089
										}
										position++
										goto l1088
									l1089:
										position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
										if buffer[position] != rune('A') {
											goto l938
										}
										position++
									}
								l1088:
									{
										position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1091
										}
										position++
										goto l1090
									l1091:
										position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
										if buffer[position] != rune('T') {
											goto l938
										}
										position++
									}
								l1090:
									{
										position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1093
										}
										position++
										goto l1092
									l1093:
										position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
										if buffer[position] != rune('C') {
											goto l938
										}
										position++
									}
								l1092:
									{
										position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1095
										}
										position++
										goto l1094
									l1095:
										position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
										if buffer[position] != rune('H') {
											goto l938
										}
										position++
									}
								l1094:
									{
										position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1097
										}
										position++
										goto l1096
									l1097:
										position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
										if buffer[position] != rune('E') {
											goto l938
										}
										position++
									}
								l1096:
									{
										position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1099
										}
										position++
										goto l1098
									l1099:
										position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
										if buffer[position] != rune('S') {
											goto l938
										}
										position++
									}
								l1098:
									if !rules[ruleskip]() {
										goto l938
									}
									depth--
									add(ruleLANGMATCHES, position1077)
								}
								break
							}
						}

					}
				l939:
					if !rules[ruleLPAREN]() {
						goto l938
					}
					if !rules[ruleexpression]() {
						goto l938
					}
					if !rules[ruleCOMMA]() {
						goto l938
					}
					if !rules[ruleexpression]() {
						goto l938
					}
					if !rules[ruleRPAREN]() {
						goto l938
					}
					goto l574
				l938:
					position, tokenIndex, depth = position574, tokenIndex574, depth574
					{
						position1101 := position
						depth++
						{
							position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1103
							}
							position++
							goto l1102
						l1103:
							position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
							if buffer[position] != rune('B') {
								goto l1100
							}
							position++
						}
					l1102:
						{
							position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1105
							}
							position++
							goto l1104
						l1105:
							position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
							if buffer[position] != rune('O') {
								goto l1100
							}
							position++
						}
					l1104:
						{
							position1106, tokenIndex1106, depth1106 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1107
							}
							position++
							goto l1106
						l1107:
							position, tokenIndex, depth = position1106, tokenIndex1106, depth1106
							if buffer[position] != rune('U') {
								goto l1100
							}
							position++
						}
					l1106:
						{
							position1108, tokenIndex1108, depth1108 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1109
							}
							position++
							goto l1108
						l1109:
							position, tokenIndex, depth = position1108, tokenIndex1108, depth1108
							if buffer[position] != rune('N') {
								goto l1100
							}
							position++
						}
					l1108:
						{
							position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1111
							}
							position++
							goto l1110
						l1111:
							position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
							if buffer[position] != rune('D') {
								goto l1100
							}
							position++
						}
					l1110:
						if !rules[ruleskip]() {
							goto l1100
						}
						depth--
						add(ruleBOUND, position1101)
					}
					if !rules[ruleLPAREN]() {
						goto l1100
					}
					if !rules[rulevar]() {
						goto l1100
					}
					if !rules[ruleRPAREN]() {
						goto l1100
					}
					goto l574
				l1100:
					position, tokenIndex, depth = position574, tokenIndex574, depth574
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1114 := position
								depth++
								{
									position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1116
									}
									position++
									goto l1115
								l1116:
									position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
									if buffer[position] != rune('S') {
										goto l1112
									}
									position++
								}
							l1115:
								{
									position1117, tokenIndex1117, depth1117 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1118
									}
									position++
									goto l1117
								l1118:
									position, tokenIndex, depth = position1117, tokenIndex1117, depth1117
									if buffer[position] != rune('T') {
										goto l1112
									}
									position++
								}
							l1117:
								{
									position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1120
									}
									position++
									goto l1119
								l1120:
									position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
									if buffer[position] != rune('R') {
										goto l1112
									}
									position++
								}
							l1119:
								{
									position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1122
									}
									position++
									goto l1121
								l1122:
									position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
									if buffer[position] != rune('U') {
										goto l1112
									}
									position++
								}
							l1121:
								{
									position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1124
									}
									position++
									goto l1123
								l1124:
									position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
									if buffer[position] != rune('U') {
										goto l1112
									}
									position++
								}
							l1123:
								{
									position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1126
									}
									position++
									goto l1125
								l1126:
									position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
									if buffer[position] != rune('I') {
										goto l1112
									}
									position++
								}
							l1125:
								{
									position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1128
									}
									position++
									goto l1127
								l1128:
									position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
									if buffer[position] != rune('D') {
										goto l1112
									}
									position++
								}
							l1127:
								if !rules[ruleskip]() {
									goto l1112
								}
								depth--
								add(ruleSTRUUID, position1114)
							}
							break
						case 'U', 'u':
							{
								position1129 := position
								depth++
								{
									position1130, tokenIndex1130, depth1130 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1131
									}
									position++
									goto l1130
								l1131:
									position, tokenIndex, depth = position1130, tokenIndex1130, depth1130
									if buffer[position] != rune('U') {
										goto l1112
									}
									position++
								}
							l1130:
								{
									position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1133
									}
									position++
									goto l1132
								l1133:
									position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
									if buffer[position] != rune('U') {
										goto l1112
									}
									position++
								}
							l1132:
								{
									position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1135
									}
									position++
									goto l1134
								l1135:
									position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
									if buffer[position] != rune('I') {
										goto l1112
									}
									position++
								}
							l1134:
								{
									position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1137
									}
									position++
									goto l1136
								l1137:
									position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
									if buffer[position] != rune('D') {
										goto l1112
									}
									position++
								}
							l1136:
								if !rules[ruleskip]() {
									goto l1112
								}
								depth--
								add(ruleUUID, position1129)
							}
							break
						case 'N', 'n':
							{
								position1138 := position
								depth++
								{
									position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1140
									}
									position++
									goto l1139
								l1140:
									position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
									if buffer[position] != rune('N') {
										goto l1112
									}
									position++
								}
							l1139:
								{
									position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1142
									}
									position++
									goto l1141
								l1142:
									position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
									if buffer[position] != rune('O') {
										goto l1112
									}
									position++
								}
							l1141:
								{
									position1143, tokenIndex1143, depth1143 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1144
									}
									position++
									goto l1143
								l1144:
									position, tokenIndex, depth = position1143, tokenIndex1143, depth1143
									if buffer[position] != rune('W') {
										goto l1112
									}
									position++
								}
							l1143:
								if !rules[ruleskip]() {
									goto l1112
								}
								depth--
								add(ruleNOW, position1138)
							}
							break
						default:
							{
								position1145 := position
								depth++
								{
									position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1147
									}
									position++
									goto l1146
								l1147:
									position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
									if buffer[position] != rune('R') {
										goto l1112
									}
									position++
								}
							l1146:
								{
									position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1149
									}
									position++
									goto l1148
								l1149:
									position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
									if buffer[position] != rune('A') {
										goto l1112
									}
									position++
								}
							l1148:
								{
									position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1151
									}
									position++
									goto l1150
								l1151:
									position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
									if buffer[position] != rune('N') {
										goto l1112
									}
									position++
								}
							l1150:
								{
									position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1153
									}
									position++
									goto l1152
								l1153:
									position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
									if buffer[position] != rune('D') {
										goto l1112
									}
									position++
								}
							l1152:
								if !rules[ruleskip]() {
									goto l1112
								}
								depth--
								add(ruleRAND, position1145)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1112
					}
					goto l574
				l1112:
					position, tokenIndex, depth = position574, tokenIndex574, depth574
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
								{
									position1157 := position
									depth++
									{
										position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1159
										}
										position++
										goto l1158
									l1159:
										position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
										if buffer[position] != rune('E') {
											goto l1156
										}
										position++
									}
								l1158:
									{
										position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1161
										}
										position++
										goto l1160
									l1161:
										position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
										if buffer[position] != rune('X') {
											goto l1156
										}
										position++
									}
								l1160:
									{
										position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1163
										}
										position++
										goto l1162
									l1163:
										position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
										if buffer[position] != rune('I') {
											goto l1156
										}
										position++
									}
								l1162:
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
											goto l1156
										}
										position++
									}
								l1164:
									{
										position1166, tokenIndex1166, depth1166 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1167
										}
										position++
										goto l1166
									l1167:
										position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
										if buffer[position] != rune('T') {
											goto l1156
										}
										position++
									}
								l1166:
									{
										position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1169
										}
										position++
										goto l1168
									l1169:
										position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
										if buffer[position] != rune('S') {
											goto l1156
										}
										position++
									}
								l1168:
									if !rules[ruleskip]() {
										goto l1156
									}
									depth--
									add(ruleEXISTS, position1157)
								}
								goto l1155
							l1156:
								position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
								{
									position1170 := position
									depth++
									{
										position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1172
										}
										position++
										goto l1171
									l1172:
										position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
										if buffer[position] != rune('N') {
											goto l572
										}
										position++
									}
								l1171:
									{
										position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1174
										}
										position++
										goto l1173
									l1174:
										position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
										if buffer[position] != rune('O') {
											goto l572
										}
										position++
									}
								l1173:
									{
										position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1176
										}
										position++
										goto l1175
									l1176:
										position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
										if buffer[position] != rune('T') {
											goto l572
										}
										position++
									}
								l1175:
									if buffer[position] != rune(' ') {
										goto l572
									}
									position++
									{
										position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1178
										}
										position++
										goto l1177
									l1178:
										position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
										if buffer[position] != rune('E') {
											goto l572
										}
										position++
									}
								l1177:
									{
										position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1180
										}
										position++
										goto l1179
									l1180:
										position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
										if buffer[position] != rune('X') {
											goto l572
										}
										position++
									}
								l1179:
									{
										position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1182
										}
										position++
										goto l1181
									l1182:
										position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
										if buffer[position] != rune('I') {
											goto l572
										}
										position++
									}
								l1181:
									{
										position1183, tokenIndex1183, depth1183 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1184
										}
										position++
										goto l1183
									l1184:
										position, tokenIndex, depth = position1183, tokenIndex1183, depth1183
										if buffer[position] != rune('S') {
											goto l572
										}
										position++
									}
								l1183:
									{
										position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1186
										}
										position++
										goto l1185
									l1186:
										position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
										if buffer[position] != rune('T') {
											goto l572
										}
										position++
									}
								l1185:
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
											goto l572
										}
										position++
									}
								l1187:
									if !rules[ruleskip]() {
										goto l572
									}
									depth--
									add(ruleNOTEXIST, position1170)
								}
							}
						l1155:
							if !rules[rulegroupGraphPattern]() {
								goto l572
							}
							break
						case 'I', 'i':
							{
								position1189 := position
								depth++
								{
									position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1191
									}
									position++
									goto l1190
								l1191:
									position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
									if buffer[position] != rune('I') {
										goto l572
									}
									position++
								}
							l1190:
								{
									position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1193
									}
									position++
									goto l1192
								l1193:
									position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
									if buffer[position] != rune('F') {
										goto l572
									}
									position++
								}
							l1192:
								if !rules[ruleskip]() {
									goto l572
								}
								depth--
								add(ruleIF, position1189)
							}
							if !rules[ruleLPAREN]() {
								goto l572
							}
							if !rules[ruleexpression]() {
								goto l572
							}
							if !rules[ruleCOMMA]() {
								goto l572
							}
							if !rules[ruleexpression]() {
								goto l572
							}
							if !rules[ruleCOMMA]() {
								goto l572
							}
							if !rules[ruleexpression]() {
								goto l572
							}
							if !rules[ruleRPAREN]() {
								goto l572
							}
							break
						case 'C', 'c':
							{
								position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
								{
									position1196 := position
									depth++
									{
										position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1198
										}
										position++
										goto l1197
									l1198:
										position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
										if buffer[position] != rune('C') {
											goto l1195
										}
										position++
									}
								l1197:
									{
										position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1200
										}
										position++
										goto l1199
									l1200:
										position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
										if buffer[position] != rune('O') {
											goto l1195
										}
										position++
									}
								l1199:
									{
										position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1202
										}
										position++
										goto l1201
									l1202:
										position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
										if buffer[position] != rune('N') {
											goto l1195
										}
										position++
									}
								l1201:
									{
										position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1204
										}
										position++
										goto l1203
									l1204:
										position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
										if buffer[position] != rune('C') {
											goto l1195
										}
										position++
									}
								l1203:
									{
										position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1206
										}
										position++
										goto l1205
									l1206:
										position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
										if buffer[position] != rune('A') {
											goto l1195
										}
										position++
									}
								l1205:
									{
										position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1208
										}
										position++
										goto l1207
									l1208:
										position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
										if buffer[position] != rune('T') {
											goto l1195
										}
										position++
									}
								l1207:
									if !rules[ruleskip]() {
										goto l1195
									}
									depth--
									add(ruleCONCAT, position1196)
								}
								goto l1194
							l1195:
								position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
								{
									position1209 := position
									depth++
									{
										position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1211
										}
										position++
										goto l1210
									l1211:
										position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
										if buffer[position] != rune('C') {
											goto l572
										}
										position++
									}
								l1210:
									{
										position1212, tokenIndex1212, depth1212 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1213
										}
										position++
										goto l1212
									l1213:
										position, tokenIndex, depth = position1212, tokenIndex1212, depth1212
										if buffer[position] != rune('O') {
											goto l572
										}
										position++
									}
								l1212:
									{
										position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1215
										}
										position++
										goto l1214
									l1215:
										position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
										if buffer[position] != rune('A') {
											goto l572
										}
										position++
									}
								l1214:
									{
										position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1217
										}
										position++
										goto l1216
									l1217:
										position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
										if buffer[position] != rune('L') {
											goto l572
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
											goto l572
										}
										position++
									}
								l1218:
									{
										position1220, tokenIndex1220, depth1220 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1221
										}
										position++
										goto l1220
									l1221:
										position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
										if buffer[position] != rune('S') {
											goto l572
										}
										position++
									}
								l1220:
									{
										position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1223
										}
										position++
										goto l1222
									l1223:
										position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
										if buffer[position] != rune('C') {
											goto l572
										}
										position++
									}
								l1222:
									{
										position1224, tokenIndex1224, depth1224 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1225
										}
										position++
										goto l1224
									l1225:
										position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
										if buffer[position] != rune('E') {
											goto l572
										}
										position++
									}
								l1224:
									if !rules[ruleskip]() {
										goto l572
									}
									depth--
									add(ruleCOALESCE, position1209)
								}
							}
						l1194:
							if !rules[ruleargList]() {
								goto l572
							}
							break
						case 'B', 'b':
							{
								position1226 := position
								depth++
								{
									position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1228
									}
									position++
									goto l1227
								l1228:
									position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
									if buffer[position] != rune('B') {
										goto l572
									}
									position++
								}
							l1227:
								{
									position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1230
									}
									position++
									goto l1229
								l1230:
									position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
									if buffer[position] != rune('N') {
										goto l572
									}
									position++
								}
							l1229:
								{
									position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1232
									}
									position++
									goto l1231
								l1232:
									position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
									if buffer[position] != rune('O') {
										goto l572
									}
									position++
								}
							l1231:
								{
									position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1234
									}
									position++
									goto l1233
								l1234:
									position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
									if buffer[position] != rune('D') {
										goto l572
									}
									position++
								}
							l1233:
								{
									position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1236
									}
									position++
									goto l1235
								l1236:
									position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
									if buffer[position] != rune('E') {
										goto l572
									}
									position++
								}
							l1235:
								if !rules[ruleskip]() {
									goto l572
								}
								depth--
								add(ruleBNODE, position1226)
							}
							{
								position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1238
								}
								if !rules[ruleexpression]() {
									goto l1238
								}
								if !rules[ruleRPAREN]() {
									goto l1238
								}
								goto l1237
							l1238:
								position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
								if !rules[rulenil]() {
									goto l572
								}
							}
						l1237:
							break
						default:
							{
								position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
								{
									position1241 := position
									depth++
									{
										position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1243
										}
										position++
										goto l1242
									l1243:
										position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
										if buffer[position] != rune('S') {
											goto l1240
										}
										position++
									}
								l1242:
									{
										position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1245
										}
										position++
										goto l1244
									l1245:
										position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
										if buffer[position] != rune('U') {
											goto l1240
										}
										position++
									}
								l1244:
									{
										position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1247
										}
										position++
										goto l1246
									l1247:
										position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
										if buffer[position] != rune('B') {
											goto l1240
										}
										position++
									}
								l1246:
									{
										position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1249
										}
										position++
										goto l1248
									l1249:
										position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
										if buffer[position] != rune('S') {
											goto l1240
										}
										position++
									}
								l1248:
									{
										position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1251
										}
										position++
										goto l1250
									l1251:
										position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
										if buffer[position] != rune('T') {
											goto l1240
										}
										position++
									}
								l1250:
									{
										position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1253
										}
										position++
										goto l1252
									l1253:
										position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
										if buffer[position] != rune('R') {
											goto l1240
										}
										position++
									}
								l1252:
									if !rules[ruleskip]() {
										goto l1240
									}
									depth--
									add(ruleSUBSTR, position1241)
								}
								goto l1239
							l1240:
								position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
								{
									position1255 := position
									depth++
									{
										position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1257
										}
										position++
										goto l1256
									l1257:
										position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
										if buffer[position] != rune('R') {
											goto l1254
										}
										position++
									}
								l1256:
									{
										position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1259
										}
										position++
										goto l1258
									l1259:
										position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
										if buffer[position] != rune('E') {
											goto l1254
										}
										position++
									}
								l1258:
									{
										position1260, tokenIndex1260, depth1260 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1261
										}
										position++
										goto l1260
									l1261:
										position, tokenIndex, depth = position1260, tokenIndex1260, depth1260
										if buffer[position] != rune('P') {
											goto l1254
										}
										position++
									}
								l1260:
									{
										position1262, tokenIndex1262, depth1262 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1263
										}
										position++
										goto l1262
									l1263:
										position, tokenIndex, depth = position1262, tokenIndex1262, depth1262
										if buffer[position] != rune('L') {
											goto l1254
										}
										position++
									}
								l1262:
									{
										position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1265
										}
										position++
										goto l1264
									l1265:
										position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
										if buffer[position] != rune('A') {
											goto l1254
										}
										position++
									}
								l1264:
									{
										position1266, tokenIndex1266, depth1266 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1267
										}
										position++
										goto l1266
									l1267:
										position, tokenIndex, depth = position1266, tokenIndex1266, depth1266
										if buffer[position] != rune('C') {
											goto l1254
										}
										position++
									}
								l1266:
									{
										position1268, tokenIndex1268, depth1268 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1269
										}
										position++
										goto l1268
									l1269:
										position, tokenIndex, depth = position1268, tokenIndex1268, depth1268
										if buffer[position] != rune('E') {
											goto l1254
										}
										position++
									}
								l1268:
									if !rules[ruleskip]() {
										goto l1254
									}
									depth--
									add(ruleREPLACE, position1255)
								}
								goto l1239
							l1254:
								position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
								{
									position1270 := position
									depth++
									{
										position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1272
										}
										position++
										goto l1271
									l1272:
										position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
										if buffer[position] != rune('R') {
											goto l572
										}
										position++
									}
								l1271:
									{
										position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1274
										}
										position++
										goto l1273
									l1274:
										position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
										if buffer[position] != rune('E') {
											goto l572
										}
										position++
									}
								l1273:
									{
										position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1276
										}
										position++
										goto l1275
									l1276:
										position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
										if buffer[position] != rune('G') {
											goto l572
										}
										position++
									}
								l1275:
									{
										position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1278
										}
										position++
										goto l1277
									l1278:
										position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
										if buffer[position] != rune('E') {
											goto l572
										}
										position++
									}
								l1277:
									{
										position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1280
										}
										position++
										goto l1279
									l1280:
										position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
										if buffer[position] != rune('X') {
											goto l572
										}
										position++
									}
								l1279:
									if !rules[ruleskip]() {
										goto l572
									}
									depth--
									add(ruleREGEX, position1270)
								}
							}
						l1239:
							if !rules[ruleLPAREN]() {
								goto l572
							}
							if !rules[ruleexpression]() {
								goto l572
							}
							if !rules[ruleCOMMA]() {
								goto l572
							}
							if !rules[ruleexpression]() {
								goto l572
							}
							{
								position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1281
								}
								if !rules[ruleexpression]() {
									goto l1281
								}
								goto l1282
							l1281:
								position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
							}
						l1282:
							if !rules[ruleRPAREN]() {
								goto l572
							}
							break
						}
					}

				}
			l574:
				depth--
				add(rulebuiltinCall, position573)
			}
			return true
		l572:
			position, tokenIndex, depth = position572, tokenIndex572, depth572
			return false
		},
		/* 61 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
			{
				position1284 := position
				depth++
				{
					position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
					{
						position1287 := position
						depth++
					l1288:
						{
							position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
							{
								position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1291
								}
								position++
								goto l1290
							l1291:
								position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1289
								}
								position++
							}
						l1290:
							goto l1288
						l1289:
							position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
						}
						depth--
						add(rulePegText, position1287)
					}
					if buffer[position] != rune(':') {
						goto l1286
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1285
				l1286:
					position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
					{
						position1294 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1293
						}
						position++
					l1295:
						{
							position1296, tokenIndex1296, depth1296 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1296
							}
							position++
							goto l1295
						l1296:
							position, tokenIndex, depth = position1296, tokenIndex1296, depth1296
						}
						depth--
						add(rulePegText, position1294)
					}
					if buffer[position] != rune('/') {
						goto l1293
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1285
				l1293:
					position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
					{
						position1298 := position
						depth++
					l1299:
						{
							position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1300
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1300
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1300
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1300
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1300
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1300
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1300
									}
									position++
									break
								}
							}

							goto l1299
						l1300:
							position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
						}
						depth--
						add(rulePegText, position1298)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1285:
				if buffer[position] != rune('<') {
					goto l1283
				}
				position++
				if !rules[rulews]() {
					goto l1283
				}
				if !rules[ruleskip]() {
					goto l1283
				}
				depth--
				add(rulepof, position1284)
			}
			return true
		l1283:
			position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
			return false
		},
		/* 62 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
			{
				position1304 := position
				depth++
				{
					position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1306
					}
					position++
					goto l1305
				l1306:
					position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
					if buffer[position] != rune('$') {
						goto l1303
					}
					position++
				}
			l1305:
				{
					position1307 := position
					depth++
					{
						position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
						{
							position1312 := position
							depth++
							{
								position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
								{
									position1315 := position
									depth++
									{
										position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1317
										}
										position++
										goto l1316
									l1317:
										position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1314
										}
										position++
									}
								l1316:
									depth--
									add(rulePN_CHARS_BASE, position1315)
								}
								goto l1313
							l1314:
								position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
								if buffer[position] != rune('_') {
									goto l1311
								}
								position++
							}
						l1313:
							depth--
							add(rulePN_CHARS_U, position1312)
						}
						goto l1310
					l1311:
						position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1303
						}
						position++
					}
				l1310:
				l1308:
					{
						position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
						{
							position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
							{
								position1320 := position
								depth++
								{
									position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
									{
										position1323 := position
										depth++
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
												goto l1322
											}
											position++
										}
									l1324:
										depth--
										add(rulePN_CHARS_BASE, position1323)
									}
									goto l1321
								l1322:
									position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
									if buffer[position] != rune('_') {
										goto l1319
									}
									position++
								}
							l1321:
								depth--
								add(rulePN_CHARS_U, position1320)
							}
							goto l1318
						l1319:
							position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1309
							}
							position++
						}
					l1318:
						goto l1308
					l1309:
						position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
					}
					depth--
					add(ruleVARNAME, position1307)
				}
				if !rules[ruleskip]() {
					goto l1303
				}
				depth--
				add(rulevar, position1304)
			}
			return true
		l1303:
			position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
			return false
		},
		/* 63 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
			{
				position1327 := position
				depth++
				{
					position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1329
					}
					goto l1328
				l1329:
					position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
					{
						position1330 := position
						depth++
					l1331:
						{
							position1332, tokenIndex1332, depth1332 := position, tokenIndex, depth
							{
								position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
								{
									position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1335
									}
									position++
									goto l1334
								l1335:
									position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
									if buffer[position] != rune(' ') {
										goto l1333
									}
									position++
								}
							l1334:
								goto l1332
							l1333:
								position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
							}
							if !matchDot() {
								goto l1332
							}
							goto l1331
						l1332:
							position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
						}
						if buffer[position] != rune(':') {
							goto l1326
						}
						position++
					l1336:
						{
							position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
							{
								position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1339
								}
								position++
								goto l1338
							l1339:
								position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1340
								}
								position++
								goto l1338
							l1340:
								position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1341
								}
								position++
								goto l1338
							l1341:
								position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1337
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1337
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1337
										}
										position++
										break
									}
								}

							}
						l1338:
							goto l1336
						l1337:
							position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
						}
						if !rules[ruleskip]() {
							goto l1326
						}
						depth--
						add(ruleprefixedName, position1330)
					}
				}
			l1328:
				depth--
				add(ruleiriref, position1327)
			}
			return true
		l1326:
			position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
			return false
		},
		/* 64 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
			{
				position1344 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1343
				}
				position++
			l1345:
				{
					position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
					{
						position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1347
						}
						position++
						goto l1346
					l1347:
						position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
					}
					if !matchDot() {
						goto l1346
					}
					goto l1345
				l1346:
					position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
				}
				if buffer[position] != rune('>') {
					goto l1343
				}
				position++
				if !rules[ruleskip]() {
					goto l1343
				}
				depth--
				add(ruleiri, position1344)
			}
			return true
		l1343:
			position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
			return false
		},
		/* 65 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 66 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
			{
				position1350 := position
				depth++
				{
					position1351 := position
					depth++
					if buffer[position] != rune('"') {
						goto l1349
					}
					position++
				l1352:
					{
						position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
						{
							position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l1354
							}
							position++
							goto l1353
						l1354:
							position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
						}
						if !matchDot() {
							goto l1353
						}
						goto l1352
					l1353:
						position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
					}
					if buffer[position] != rune('"') {
						goto l1349
					}
					position++
					depth--
					add(rulestring, position1351)
				}
				{
					position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
					{
						position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1358
						}
						position++
						{
							position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1362
							}
							position++
							goto l1361
						l1362:
							position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1358
							}
							position++
						}
					l1361:
					l1359:
						{
							position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
							{
								position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1364
								}
								position++
								goto l1363
							l1364:
								position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1360
								}
								position++
							}
						l1363:
							goto l1359
						l1360:
							position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
						}
					l1365:
						{
							position1366, tokenIndex1366, depth1366 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1366
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1366
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1366
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1366
									}
									position++
									break
								}
							}

						l1367:
							{
								position1368, tokenIndex1368, depth1368 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1368
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1368
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1368
										}
										position++
										break
									}
								}

								goto l1367
							l1368:
								position, tokenIndex, depth = position1368, tokenIndex1368, depth1368
							}
							goto l1365
						l1366:
							position, tokenIndex, depth = position1366, tokenIndex1366, depth1366
						}
						goto l1357
					l1358:
						position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
						if buffer[position] != rune('^') {
							goto l1355
						}
						position++
						if buffer[position] != rune('^') {
							goto l1355
						}
						position++
						if !rules[ruleiriref]() {
							goto l1355
						}
					}
				l1357:
					goto l1356
				l1355:
					position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
				}
			l1356:
				if !rules[ruleskip]() {
					goto l1349
				}
				depth--
				add(ruleliteral, position1350)
			}
			return true
		l1349:
			position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
			return false
		},
		/* 67 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 68 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
			{
				position1373 := position
				depth++
				{
					position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
					{
						position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1377
						}
						position++
						goto l1376
					l1377:
						position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
						if buffer[position] != rune('-') {
							goto l1374
						}
						position++
					}
				l1376:
					goto l1375
				l1374:
					position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
				}
			l1375:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1372
				}
				position++
			l1378:
				{
					position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1379
					}
					position++
					goto l1378
				l1379:
					position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
				}
				{
					position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1380
					}
					position++
				l1382:
					{
						position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1383
						}
						position++
						goto l1382
					l1383:
						position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
					}
					goto l1381
				l1380:
					position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
				}
			l1381:
				if !rules[ruleskip]() {
					goto l1372
				}
				depth--
				add(rulenumericLiteral, position1373)
			}
			return true
		l1372:
			position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
			return false
		},
		/* 69 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 70 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
			{
				position1386 := position
				depth++
				{
					position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
					{
						position1389 := position
						depth++
						{
							position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1391
							}
							position++
							goto l1390
						l1391:
							position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
							if buffer[position] != rune('T') {
								goto l1388
							}
							position++
						}
					l1390:
						{
							position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1393
							}
							position++
							goto l1392
						l1393:
							position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
							if buffer[position] != rune('R') {
								goto l1388
							}
							position++
						}
					l1392:
						{
							position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1395
							}
							position++
							goto l1394
						l1395:
							position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
							if buffer[position] != rune('U') {
								goto l1388
							}
							position++
						}
					l1394:
						{
							position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1397
							}
							position++
							goto l1396
						l1397:
							position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
							if buffer[position] != rune('E') {
								goto l1388
							}
							position++
						}
					l1396:
						if !rules[ruleskip]() {
							goto l1388
						}
						depth--
						add(ruleTRUE, position1389)
					}
					goto l1387
				l1388:
					position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
					{
						position1398 := position
						depth++
						{
							position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1400
							}
							position++
							goto l1399
						l1400:
							position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
							if buffer[position] != rune('F') {
								goto l1385
							}
							position++
						}
					l1399:
						{
							position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1402
							}
							position++
							goto l1401
						l1402:
							position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
							if buffer[position] != rune('A') {
								goto l1385
							}
							position++
						}
					l1401:
						{
							position1403, tokenIndex1403, depth1403 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1404
							}
							position++
							goto l1403
						l1404:
							position, tokenIndex, depth = position1403, tokenIndex1403, depth1403
							if buffer[position] != rune('L') {
								goto l1385
							}
							position++
						}
					l1403:
						{
							position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1406
							}
							position++
							goto l1405
						l1406:
							position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
							if buffer[position] != rune('S') {
								goto l1385
							}
							position++
						}
					l1405:
						{
							position1407, tokenIndex1407, depth1407 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1408
							}
							position++
							goto l1407
						l1408:
							position, tokenIndex, depth = position1407, tokenIndex1407, depth1407
							if buffer[position] != rune('E') {
								goto l1385
							}
							position++
						}
					l1407:
						if !rules[ruleskip]() {
							goto l1385
						}
						depth--
						add(ruleFALSE, position1398)
					}
				}
			l1387:
				depth--
				add(rulebooleanLiteral, position1386)
			}
			return true
		l1385:
			position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
			return false
		},
		/* 71 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 72 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 73 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 74 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
			{
				position1413 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1412
				}
				position++
			l1414:
				{
					position1415, tokenIndex1415, depth1415 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1415
					}
					goto l1414
				l1415:
					position, tokenIndex, depth = position1415, tokenIndex1415, depth1415
				}
				if buffer[position] != rune(')') {
					goto l1412
				}
				position++
				if !rules[ruleskip]() {
					goto l1412
				}
				depth--
				add(rulenil, position1413)
			}
			return true
		l1412:
			position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
			return false
		},
		/* 75 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 76 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 77 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 78 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 79 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 80 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 81 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 82 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 83 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 84 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 85 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 86 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 87 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 88 LBRACE <- <('{' skip)> */
		func() bool {
			position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
			{
				position1430 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1429
				}
				position++
				if !rules[ruleskip]() {
					goto l1429
				}
				depth--
				add(ruleLBRACE, position1430)
			}
			return true
		l1429:
			position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
			return false
		},
		/* 89 RBRACE <- <('}' skip)> */
		func() bool {
			position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
			{
				position1432 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1431
				}
				position++
				if !rules[ruleskip]() {
					goto l1431
				}
				depth--
				add(ruleRBRACE, position1432)
			}
			return true
		l1431:
			position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
			return false
		},
		/* 90 LBRACK <- <('[' skip)> */
		nil,
		/* 91 RBRACK <- <(']' skip)> */
		nil,
		/* 92 SEMICOLON <- <(';' skip)> */
		nil,
		/* 93 COMMA <- <(',' skip)> */
		func() bool {
			position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
			{
				position1437 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1436
				}
				position++
				if !rules[ruleskip]() {
					goto l1436
				}
				depth--
				add(ruleCOMMA, position1437)
			}
			return true
		l1436:
			position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
			return false
		},
		/* 94 DOT <- <('.' skip)> */
		func() bool {
			position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
			{
				position1439 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1438
				}
				position++
				if !rules[ruleskip]() {
					goto l1438
				}
				depth--
				add(ruleDOT, position1439)
			}
			return true
		l1438:
			position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
			return false
		},
		/* 95 COLON <- <(':' skip)> */
		nil,
		/* 96 PIPE <- <('|' skip)> */
		func() bool {
			position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
			{
				position1442 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1441
				}
				position++
				if !rules[ruleskip]() {
					goto l1441
				}
				depth--
				add(rulePIPE, position1442)
			}
			return true
		l1441:
			position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
			return false
		},
		/* 97 SLASH <- <('/' skip)> */
		func() bool {
			position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
			{
				position1444 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1443
				}
				position++
				if !rules[ruleskip]() {
					goto l1443
				}
				depth--
				add(ruleSLASH, position1444)
			}
			return true
		l1443:
			position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
			return false
		},
		/* 98 INVERSE <- <('^' skip)> */
		func() bool {
			position1445, tokenIndex1445, depth1445 := position, tokenIndex, depth
			{
				position1446 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1445
				}
				position++
				if !rules[ruleskip]() {
					goto l1445
				}
				depth--
				add(ruleINVERSE, position1446)
			}
			return true
		l1445:
			position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
			return false
		},
		/* 99 LPAREN <- <('(' skip)> */
		func() bool {
			position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
			{
				position1448 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1447
				}
				position++
				if !rules[ruleskip]() {
					goto l1447
				}
				depth--
				add(ruleLPAREN, position1448)
			}
			return true
		l1447:
			position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
			return false
		},
		/* 100 RPAREN <- <(')' skip)> */
		func() bool {
			position1449, tokenIndex1449, depth1449 := position, tokenIndex, depth
			{
				position1450 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1449
				}
				position++
				if !rules[ruleskip]() {
					goto l1449
				}
				depth--
				add(ruleRPAREN, position1450)
			}
			return true
		l1449:
			position, tokenIndex, depth = position1449, tokenIndex1449, depth1449
			return false
		},
		/* 101 ISA <- <('a' skip)> */
		func() bool {
			position1451, tokenIndex1451, depth1451 := position, tokenIndex, depth
			{
				position1452 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1451
				}
				position++
				if !rules[ruleskip]() {
					goto l1451
				}
				depth--
				add(ruleISA, position1452)
			}
			return true
		l1451:
			position, tokenIndex, depth = position1451, tokenIndex1451, depth1451
			return false
		},
		/* 102 NOT <- <('!' skip)> */
		func() bool {
			position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
			{
				position1454 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1453
				}
				position++
				if !rules[ruleskip]() {
					goto l1453
				}
				depth--
				add(ruleNOT, position1454)
			}
			return true
		l1453:
			position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
			return false
		},
		/* 103 STAR <- <('*' skip)> */
		func() bool {
			position1455, tokenIndex1455, depth1455 := position, tokenIndex, depth
			{
				position1456 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1455
				}
				position++
				if !rules[ruleskip]() {
					goto l1455
				}
				depth--
				add(ruleSTAR, position1456)
			}
			return true
		l1455:
			position, tokenIndex, depth = position1455, tokenIndex1455, depth1455
			return false
		},
		/* 104 PLUS <- <('+' skip)> */
		func() bool {
			position1457, tokenIndex1457, depth1457 := position, tokenIndex, depth
			{
				position1458 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1457
				}
				position++
				if !rules[ruleskip]() {
					goto l1457
				}
				depth--
				add(rulePLUS, position1458)
			}
			return true
		l1457:
			position, tokenIndex, depth = position1457, tokenIndex1457, depth1457
			return false
		},
		/* 105 MINUS <- <('-' skip)> */
		func() bool {
			position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
			{
				position1460 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1459
				}
				position++
				if !rules[ruleskip]() {
					goto l1459
				}
				depth--
				add(ruleMINUS, position1460)
			}
			return true
		l1459:
			position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
			return false
		},
		/* 106 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 107 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 108 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 109 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 110 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
			{
				position1466 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1465
				}
				position++
			l1467:
				{
					position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1468
					}
					position++
					goto l1467
				l1468:
					position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
				}
				if !rules[ruleskip]() {
					goto l1465
				}
				depth--
				add(ruleINTEGER, position1466)
			}
			return true
		l1465:
			position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
			return false
		},
		/* 111 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 112 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 113 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 114 OR <- <('|' '|' skip)> */
		nil,
		/* 115 AND <- <('&' '&' skip)> */
		nil,
		/* 116 EQ <- <('=' skip)> */
		nil,
		/* 117 NE <- <('!' '=' skip)> */
		nil,
		/* 118 GT <- <('>' skip)> */
		nil,
		/* 119 LT <- <('<' skip)> */
		nil,
		/* 120 LE <- <('<' '=' skip)> */
		nil,
		/* 121 GE <- <('>' '=' skip)> */
		nil,
		/* 122 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 123 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 124 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1482, tokenIndex1482, depth1482 := position, tokenIndex, depth
			{
				position1483 := position
				depth++
				{
					position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1485
					}
					position++
					goto l1484
				l1485:
					position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
					if buffer[position] != rune('A') {
						goto l1482
					}
					position++
				}
			l1484:
				{
					position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1487
					}
					position++
					goto l1486
				l1487:
					position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
					if buffer[position] != rune('S') {
						goto l1482
					}
					position++
				}
			l1486:
				if !rules[ruleskip]() {
					goto l1482
				}
				depth--
				add(ruleAS, position1483)
			}
			return true
		l1482:
			position, tokenIndex, depth = position1482, tokenIndex1482, depth1482
			return false
		},
		/* 125 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 126 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 127 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 128 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 129 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 130 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 131 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 132 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 133 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 134 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 135 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 136 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 137 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 138 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 139 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 140 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 141 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 142 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 143 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 144 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 145 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 146 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 147 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 148 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 149 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 150 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 151 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 152 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 153 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 154 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 155 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 156 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 157 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 158 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 159 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 160 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 161 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 162 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 163 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 164 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 165 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 166 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 167 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 168 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 169 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 170 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 171 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 172 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 173 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 174 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 175 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 176 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 177 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 178 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 179 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 180 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 181 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1545 := position
				depth++
			l1546:
				{
					position1547, tokenIndex1547, depth1547 := position, tokenIndex, depth
					{
						position1548, tokenIndex1548, depth1548 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1549
						}
						goto l1548
					l1549:
						position, tokenIndex, depth = position1548, tokenIndex1548, depth1548
						{
							position1550 := position
							depth++
							{
								position1551 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1547
								}
								position++
							l1552:
								{
									position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
									{
										position1554, tokenIndex1554, depth1554 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1554
										}
										goto l1553
									l1554:
										position, tokenIndex, depth = position1554, tokenIndex1554, depth1554
									}
									if !matchDot() {
										goto l1553
									}
									goto l1552
								l1553:
									position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
								}
								if !rules[ruleendOfLine]() {
									goto l1547
								}
								depth--
								add(rulePegText, position1551)
							}
							{
								add(ruleAction13, position)
							}
							depth--
							add(rulecomment, position1550)
						}
					}
				l1548:
					goto l1546
				l1547:
					position, tokenIndex, depth = position1547, tokenIndex1547, depth1547
				}
				depth--
				add(ruleskip, position1545)
			}
			return true
		},
		/* 182 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
			{
				position1557 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1556
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1556
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1556
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1556
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1556
						}
						break
					}
				}

				depth--
				add(rulews, position1557)
			}
			return true
		l1556:
			position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
			return false
		},
		/* 183 comment <- <(<('#' (!endOfLine .)* endOfLine)> Action13)> */
		nil,
		/* 184 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1560, tokenIndex1560, depth1560 := position, tokenIndex, depth
			{
				position1561 := position
				depth++
				{
					position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1563
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1563
					}
					position++
					goto l1562
				l1563:
					position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
					if buffer[position] != rune('\n') {
						goto l1564
					}
					position++
					goto l1562
				l1564:
					position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
					if buffer[position] != rune('\r') {
						goto l1560
					}
					position++
				}
			l1562:
				depth--
				add(ruleendOfLine, position1561)
			}
			return true
		l1560:
			position, tokenIndex, depth = position1560, tokenIndex1560, depth1560
			return false
		},
		nil,
		/* 187 Action0 <- <{ p.addPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 188 Action1 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 189 Action2 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 190 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 191 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 192 Action5 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 193 Action6 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 194 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 195 Action8 <- <{ p.setObject(p.skipComment(buffer, begin, end)); p.addTriplePattern() }> */
		nil,
		/* 196 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 197 Action10 <- <{ p.setPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 198 Action11 <- <{ p.setPathLength(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 199 Action12 <- <{ p.setKeyword(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 200 Action13 <- <{ p.commentBegin = begin }> */
		nil,
	}
	p.rules = rules
}
