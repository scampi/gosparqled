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
	ruleaggregate
	rulecount
	rulegroupConcat
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
	ruleSUM
	ruleMIN
	ruleMAX
	ruleAVG
	ruleSAMPLE
	ruleCOUNT
	ruleGROUPCONCAT
	ruleSEPARATOR
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
	"aggregate",
	"count",
	"groupConcat",
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
	"SUM",
	"MIN",
	"MAX",
	"AVG",
	"SAMPLE",
	"COUNT",
	"GROUPCONCAT",
	"SEPARATOR",
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
	rules  [212]func() bool
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
						if !rules[ruleDISTINCT]() {
							goto l129
						}
						goto l128
					l129:
						position, tokenIndex, depth = position128, tokenIndex128, depth128
						{
							position130 := position
							depth++
							{
								position131, tokenIndex131, depth131 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l132
								}
								position++
								goto l131
							l132:
								position, tokenIndex, depth = position131, tokenIndex131, depth131
								if buffer[position] != rune('R') {
									goto l126
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
									goto l126
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
									goto l126
								}
								position++
							}
						l135:
							{
								position137, tokenIndex137, depth137 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l138
								}
								position++
								goto l137
							l138:
								position, tokenIndex, depth = position137, tokenIndex137, depth137
								if buffer[position] != rune('U') {
									goto l126
								}
								position++
							}
						l137:
							{
								position139, tokenIndex139, depth139 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l140
								}
								position++
								goto l139
							l140:
								position, tokenIndex, depth = position139, tokenIndex139, depth139
								if buffer[position] != rune('C') {
									goto l126
								}
								position++
							}
						l139:
							{
								position141, tokenIndex141, depth141 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l142
								}
								position++
								goto l141
							l142:
								position, tokenIndex, depth = position141, tokenIndex141, depth141
								if buffer[position] != rune('E') {
									goto l126
								}
								position++
							}
						l141:
							{
								position143, tokenIndex143, depth143 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l144
								}
								position++
								goto l143
							l144:
								position, tokenIndex, depth = position143, tokenIndex143, depth143
								if buffer[position] != rune('D') {
									goto l126
								}
								position++
							}
						l143:
							if !rules[ruleskip]() {
								goto l126
							}
							depth--
							add(ruleREDUCED, position130)
						}
					}
				l128:
					goto l127
				l126:
					position, tokenIndex, depth = position126, tokenIndex126, depth126
				}
			l127:
				{
					position145, tokenIndex145, depth145 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l146
					}
					goto l145
				l146:
					position, tokenIndex, depth = position145, tokenIndex145, depth145
					{
						position149 := position
						depth++
						{
							position150, tokenIndex150, depth150 := position, tokenIndex, depth
							if !rules[rulevar]() {
								goto l151
							}
							goto l150
						l151:
							position, tokenIndex, depth = position150, tokenIndex150, depth150
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
					l150:
						depth--
						add(ruleprojectionElem, position149)
					}
				l147:
					{
						position148, tokenIndex148, depth148 := position, tokenIndex, depth
						{
							position152 := position
							depth++
							{
								position153, tokenIndex153, depth153 := position, tokenIndex, depth
								if !rules[rulevar]() {
									goto l154
								}
								goto l153
							l154:
								position, tokenIndex, depth = position153, tokenIndex153, depth153
								if !rules[ruleLPAREN]() {
									goto l148
								}
								if !rules[ruleexpression]() {
									goto l148
								}
								if !rules[ruleAS]() {
									goto l148
								}
								if !rules[rulevar]() {
									goto l148
								}
								if !rules[ruleRPAREN]() {
									goto l148
								}
							}
						l153:
							depth--
							add(ruleprojectionElem, position152)
						}
						goto l147
					l148:
						position, tokenIndex, depth = position148, tokenIndex148, depth148
					}
				}
			l145:
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
			position155, tokenIndex155, depth155 := position, tokenIndex, depth
			{
				position156 := position
				depth++
				if !rules[ruleselect]() {
					goto l155
				}
				if !rules[rulewhereClause]() {
					goto l155
				}
				depth--
				add(rulesubSelect, position156)
			}
			return true
		l155:
			position, tokenIndex, depth = position155, tokenIndex155, depth155
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
			position163, tokenIndex163, depth163 := position, tokenIndex, depth
			{
				position164 := position
				depth++
				{
					position165 := position
					depth++
					{
						position166, tokenIndex166, depth166 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l167
						}
						position++
						goto l166
					l167:
						position, tokenIndex, depth = position166, tokenIndex166, depth166
						if buffer[position] != rune('F') {
							goto l163
						}
						position++
					}
				l166:
					{
						position168, tokenIndex168, depth168 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l169
						}
						position++
						goto l168
					l169:
						position, tokenIndex, depth = position168, tokenIndex168, depth168
						if buffer[position] != rune('R') {
							goto l163
						}
						position++
					}
				l168:
					{
						position170, tokenIndex170, depth170 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l171
						}
						position++
						goto l170
					l171:
						position, tokenIndex, depth = position170, tokenIndex170, depth170
						if buffer[position] != rune('O') {
							goto l163
						}
						position++
					}
				l170:
					{
						position172, tokenIndex172, depth172 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l173
						}
						position++
						goto l172
					l173:
						position, tokenIndex, depth = position172, tokenIndex172, depth172
						if buffer[position] != rune('M') {
							goto l163
						}
						position++
					}
				l172:
					if !rules[ruleskip]() {
						goto l163
					}
					depth--
					add(ruleFROM, position165)
				}
				{
					position174, tokenIndex174, depth174 := position, tokenIndex, depth
					{
						position176 := position
						depth++
						{
							position177, tokenIndex177, depth177 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l178
							}
							position++
							goto l177
						l178:
							position, tokenIndex, depth = position177, tokenIndex177, depth177
							if buffer[position] != rune('N') {
								goto l174
							}
							position++
						}
					l177:
						{
							position179, tokenIndex179, depth179 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l180
							}
							position++
							goto l179
						l180:
							position, tokenIndex, depth = position179, tokenIndex179, depth179
							if buffer[position] != rune('A') {
								goto l174
							}
							position++
						}
					l179:
						{
							position181, tokenIndex181, depth181 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l182
							}
							position++
							goto l181
						l182:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('M') {
								goto l174
							}
							position++
						}
					l181:
						{
							position183, tokenIndex183, depth183 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l184
							}
							position++
							goto l183
						l184:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('E') {
								goto l174
							}
							position++
						}
					l183:
						{
							position185, tokenIndex185, depth185 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l186
							}
							position++
							goto l185
						l186:
							position, tokenIndex, depth = position185, tokenIndex185, depth185
							if buffer[position] != rune('D') {
								goto l174
							}
							position++
						}
					l185:
						if !rules[ruleskip]() {
							goto l174
						}
						depth--
						add(ruleNAMED, position176)
					}
					goto l175
				l174:
					position, tokenIndex, depth = position174, tokenIndex174, depth174
				}
			l175:
				if !rules[ruleiriref]() {
					goto l163
				}
				depth--
				add(ruledatasetClause, position164)
			}
			return true
		l163:
			position, tokenIndex, depth = position163, tokenIndex163, depth163
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position187, tokenIndex187, depth187 := position, tokenIndex, depth
			{
				position188 := position
				depth++
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					{
						position191 := position
						depth++
						{
							position192, tokenIndex192, depth192 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l193
							}
							position++
							goto l192
						l193:
							position, tokenIndex, depth = position192, tokenIndex192, depth192
							if buffer[position] != rune('W') {
								goto l189
							}
							position++
						}
					l192:
						{
							position194, tokenIndex194, depth194 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l195
							}
							position++
							goto l194
						l195:
							position, tokenIndex, depth = position194, tokenIndex194, depth194
							if buffer[position] != rune('H') {
								goto l189
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
								goto l189
							}
							position++
						}
					l196:
						{
							position198, tokenIndex198, depth198 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l199
							}
							position++
							goto l198
						l199:
							position, tokenIndex, depth = position198, tokenIndex198, depth198
							if buffer[position] != rune('R') {
								goto l189
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
								goto l189
							}
							position++
						}
					l200:
						if !rules[ruleskip]() {
							goto l189
						}
						depth--
						add(ruleWHERE, position191)
					}
					goto l190
				l189:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
				}
			l190:
				if !rules[rulegroupGraphPattern]() {
					goto l187
				}
				depth--
				add(rulewhereClause, position188)
			}
			return true
		l187:
			position, tokenIndex, depth = position187, tokenIndex187, depth187
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l202
				}
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l205
					}
					goto l204
				l205:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
					if !rules[rulegraphPattern]() {
						goto l202
					}
				}
			l204:
				if !rules[ruleRBRACE]() {
					goto l202
				}
				depth--
				add(rulegroupGraphPattern, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position207 := position
				depth++
				{
					position208, tokenIndex208, depth208 := position, tokenIndex, depth
					{
						position210 := position
						depth++
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if !rules[ruletriplesBlock]() {
								goto l212
							}
						l213:
							{
								position214, tokenIndex214, depth214 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l214
								}
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
								{
									position217, tokenIndex217, depth217 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l217
									}
									goto l218
								l217:
									position, tokenIndex, depth = position217, tokenIndex217, depth217
								}
							l218:
								goto l213
							l214:
								position, tokenIndex, depth = position214, tokenIndex214, depth214
							}
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if !rules[rulefilterOrBind]() {
								goto l208
							}
							{
								position221, tokenIndex221, depth221 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l221
								}
								goto l222
							l221:
								position, tokenIndex, depth = position221, tokenIndex221, depth221
							}
						l222:
							{
								position223, tokenIndex223, depth223 := position, tokenIndex, depth
								if !rules[ruletriplesBlock]() {
									goto l223
								}
								goto l224
							l223:
								position, tokenIndex, depth = position223, tokenIndex223, depth223
							}
						l224:
						l219:
							{
								position220, tokenIndex220, depth220 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l220
								}
								{
									position225, tokenIndex225, depth225 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l225
									}
									goto l226
								l225:
									position, tokenIndex, depth = position225, tokenIndex225, depth225
								}
							l226:
								{
									position227, tokenIndex227, depth227 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l227
									}
									goto l228
								l227:
									position, tokenIndex, depth = position227, tokenIndex227, depth227
								}
							l228:
								goto l219
							l220:
								position, tokenIndex, depth = position220, tokenIndex220, depth220
							}
						}
					l211:
						depth--
						add(rulebasicGraphPattern, position210)
					}
					goto l209
				l208:
					position, tokenIndex, depth = position208, tokenIndex208, depth208
				}
			l209:
				{
					position229, tokenIndex229, depth229 := position, tokenIndex, depth
					{
						position231 := position
						depth++
						{
							position232, tokenIndex232, depth232 := position, tokenIndex, depth
							{
								position234 := position
								depth++
								{
									position235 := position
									depth++
									{
										position236, tokenIndex236, depth236 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l237
										}
										position++
										goto l236
									l237:
										position, tokenIndex, depth = position236, tokenIndex236, depth236
										if buffer[position] != rune('O') {
											goto l233
										}
										position++
									}
								l236:
									{
										position238, tokenIndex238, depth238 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l239
										}
										position++
										goto l238
									l239:
										position, tokenIndex, depth = position238, tokenIndex238, depth238
										if buffer[position] != rune('P') {
											goto l233
										}
										position++
									}
								l238:
									{
										position240, tokenIndex240, depth240 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l241
										}
										position++
										goto l240
									l241:
										position, tokenIndex, depth = position240, tokenIndex240, depth240
										if buffer[position] != rune('T') {
											goto l233
										}
										position++
									}
								l240:
									{
										position242, tokenIndex242, depth242 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l243
										}
										position++
										goto l242
									l243:
										position, tokenIndex, depth = position242, tokenIndex242, depth242
										if buffer[position] != rune('I') {
											goto l233
										}
										position++
									}
								l242:
									{
										position244, tokenIndex244, depth244 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l245
										}
										position++
										goto l244
									l245:
										position, tokenIndex, depth = position244, tokenIndex244, depth244
										if buffer[position] != rune('O') {
											goto l233
										}
										position++
									}
								l244:
									{
										position246, tokenIndex246, depth246 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l247
										}
										position++
										goto l246
									l247:
										position, tokenIndex, depth = position246, tokenIndex246, depth246
										if buffer[position] != rune('N') {
											goto l233
										}
										position++
									}
								l246:
									{
										position248, tokenIndex248, depth248 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l249
										}
										position++
										goto l248
									l249:
										position, tokenIndex, depth = position248, tokenIndex248, depth248
										if buffer[position] != rune('A') {
											goto l233
										}
										position++
									}
								l248:
									{
										position250, tokenIndex250, depth250 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l251
										}
										position++
										goto l250
									l251:
										position, tokenIndex, depth = position250, tokenIndex250, depth250
										if buffer[position] != rune('L') {
											goto l233
										}
										position++
									}
								l250:
									if !rules[ruleskip]() {
										goto l233
									}
									depth--
									add(ruleOPTIONAL, position235)
								}
								if !rules[ruleLBRACE]() {
									goto l233
								}
								{
									position252, tokenIndex252, depth252 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l253
									}
									goto l252
								l253:
									position, tokenIndex, depth = position252, tokenIndex252, depth252
									if !rules[rulegraphPattern]() {
										goto l233
									}
								}
							l252:
								if !rules[ruleRBRACE]() {
									goto l233
								}
								depth--
								add(ruleoptionalGraphPattern, position234)
							}
							goto l232
						l233:
							position, tokenIndex, depth = position232, tokenIndex232, depth232
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l229
							}
						}
					l232:
						depth--
						add(rulegraphPatternNotTriples, position231)
					}
					{
						position254, tokenIndex254, depth254 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l254
						}
						goto l255
					l254:
						position, tokenIndex, depth = position254, tokenIndex254, depth254
					}
				l255:
					if !rules[rulegraphPattern]() {
						goto l229
					}
					goto l230
				l229:
					position, tokenIndex, depth = position229, tokenIndex229, depth229
				}
			l230:
				depth--
				add(rulegraphPattern, position207)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position258, tokenIndex258, depth258 := position, tokenIndex, depth
			{
				position259 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l258
				}
				{
					position260, tokenIndex260, depth260 := position, tokenIndex, depth
					{
						position262 := position
						depth++
						{
							position263, tokenIndex263, depth263 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l264
							}
							position++
							goto l263
						l264:
							position, tokenIndex, depth = position263, tokenIndex263, depth263
							if buffer[position] != rune('U') {
								goto l260
							}
							position++
						}
					l263:
						{
							position265, tokenIndex265, depth265 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l266
							}
							position++
							goto l265
						l266:
							position, tokenIndex, depth = position265, tokenIndex265, depth265
							if buffer[position] != rune('N') {
								goto l260
							}
							position++
						}
					l265:
						{
							position267, tokenIndex267, depth267 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l268
							}
							position++
							goto l267
						l268:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
							if buffer[position] != rune('I') {
								goto l260
							}
							position++
						}
					l267:
						{
							position269, tokenIndex269, depth269 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l270
							}
							position++
							goto l269
						l270:
							position, tokenIndex, depth = position269, tokenIndex269, depth269
							if buffer[position] != rune('O') {
								goto l260
							}
							position++
						}
					l269:
						{
							position271, tokenIndex271, depth271 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l272
							}
							position++
							goto l271
						l272:
							position, tokenIndex, depth = position271, tokenIndex271, depth271
							if buffer[position] != rune('N') {
								goto l260
							}
							position++
						}
					l271:
						if !rules[ruleskip]() {
							goto l260
						}
						depth--
						add(ruleUNION, position262)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l260
					}
					goto l261
				l260:
					position, tokenIndex, depth = position260, tokenIndex260, depth260
				}
			l261:
				depth--
				add(rulegroupOrUnionGraphPattern, position259)
			}
			return true
		l258:
			position, tokenIndex, depth = position258, tokenIndex258, depth258
			return false
		},
		/* 21 basicGraphPattern <- <((triplesBlock (filterOrBind DOT? triplesBlock?)*) / (filterOrBind DOT? triplesBlock?)+)> */
		nil,
		/* 22 filterOrBind <- <((FILTER constraint) / (BIND LPAREN expression AS var RPAREN))> */
		func() bool {
			position274, tokenIndex274, depth274 := position, tokenIndex, depth
			{
				position275 := position
				depth++
				{
					position276, tokenIndex276, depth276 := position, tokenIndex, depth
					{
						position278 := position
						depth++
						{
							position279, tokenIndex279, depth279 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l280
							}
							position++
							goto l279
						l280:
							position, tokenIndex, depth = position279, tokenIndex279, depth279
							if buffer[position] != rune('F') {
								goto l277
							}
							position++
						}
					l279:
						{
							position281, tokenIndex281, depth281 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l282
							}
							position++
							goto l281
						l282:
							position, tokenIndex, depth = position281, tokenIndex281, depth281
							if buffer[position] != rune('I') {
								goto l277
							}
							position++
						}
					l281:
						{
							position283, tokenIndex283, depth283 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l284
							}
							position++
							goto l283
						l284:
							position, tokenIndex, depth = position283, tokenIndex283, depth283
							if buffer[position] != rune('L') {
								goto l277
							}
							position++
						}
					l283:
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
								goto l277
							}
							position++
						}
					l285:
						{
							position287, tokenIndex287, depth287 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l288
							}
							position++
							goto l287
						l288:
							position, tokenIndex, depth = position287, tokenIndex287, depth287
							if buffer[position] != rune('E') {
								goto l277
							}
							position++
						}
					l287:
						{
							position289, tokenIndex289, depth289 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l290
							}
							position++
							goto l289
						l290:
							position, tokenIndex, depth = position289, tokenIndex289, depth289
							if buffer[position] != rune('R') {
								goto l277
							}
							position++
						}
					l289:
						if !rules[ruleskip]() {
							goto l277
						}
						depth--
						add(ruleFILTER, position278)
					}
					{
						position291 := position
						depth++
						{
							position292, tokenIndex292, depth292 := position, tokenIndex, depth
							if !rules[rulebrackettedExpression]() {
								goto l293
							}
							goto l292
						l293:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
							if !rules[rulebuiltinCall]() {
								goto l294
							}
							goto l292
						l294:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
							if !rules[rulefunctionCall]() {
								goto l277
							}
						}
					l292:
						depth--
						add(ruleconstraint, position291)
					}
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					{
						position295 := position
						depth++
						{
							position296, tokenIndex296, depth296 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l297
							}
							position++
							goto l296
						l297:
							position, tokenIndex, depth = position296, tokenIndex296, depth296
							if buffer[position] != rune('B') {
								goto l274
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
								goto l274
							}
							position++
						}
					l298:
						{
							position300, tokenIndex300, depth300 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l301
							}
							position++
							goto l300
						l301:
							position, tokenIndex, depth = position300, tokenIndex300, depth300
							if buffer[position] != rune('N') {
								goto l274
							}
							position++
						}
					l300:
						{
							position302, tokenIndex302, depth302 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l303
							}
							position++
							goto l302
						l303:
							position, tokenIndex, depth = position302, tokenIndex302, depth302
							if buffer[position] != rune('D') {
								goto l274
							}
							position++
						}
					l302:
						if !rules[ruleskip]() {
							goto l274
						}
						depth--
						add(ruleBIND, position295)
					}
					if !rules[ruleLPAREN]() {
						goto l274
					}
					if !rules[ruleexpression]() {
						goto l274
					}
					if !rules[ruleAS]() {
						goto l274
					}
					if !rules[rulevar]() {
						goto l274
					}
					if !rules[ruleRPAREN]() {
						goto l274
					}
				}
			l276:
				depth--
				add(rulefilterOrBind, position275)
			}
			return true
		l274:
			position, tokenIndex, depth = position274, tokenIndex274, depth274
			return false
		},
		/* 23 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		nil,
		/* 24 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position305, tokenIndex305, depth305 := position, tokenIndex, depth
			{
				position306 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l305
				}
			l307:
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l308
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l308
					}
					goto l307
				l308:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
				}
				{
					position309, tokenIndex309, depth309 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l309
					}
					goto l310
				l309:
					position, tokenIndex, depth = position309, tokenIndex309, depth309
				}
			l310:
				depth--
				add(ruletriplesBlock, position306)
			}
			return true
		l305:
			position, tokenIndex, depth = position305, tokenIndex305, depth305
			return false
		},
		/* 25 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position311, tokenIndex311, depth311 := position, tokenIndex, depth
			{
				position312 := position
				depth++
				{
					position313, tokenIndex313, depth313 := position, tokenIndex, depth
					{
						position315 := position
						depth++
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							{
								position318 := position
								depth++
								if !rules[rulevar]() {
									goto l317
								}
								depth--
								add(rulePegText, position318)
							}
							{
								add(ruleAction1, position)
							}
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							{
								position321 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l320
								}
								depth--
								add(rulePegText, position321)
							}
							{
								add(ruleAction2, position)
							}
							goto l316
						l320:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if !rules[rulepof]() {
								goto l314
							}
							{
								add(ruleAction3, position)
							}
						}
					l316:
						depth--
						add(rulevarOrTerm, position315)
					}
					if !rules[rulepropertyListPath]() {
						goto l314
					}
					goto l313
				l314:
					position, tokenIndex, depth = position313, tokenIndex313, depth313
					{
						position324 := position
						depth++
						{
							position325, tokenIndex325, depth325 := position, tokenIndex, depth
							{
								position327 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l326
								}
								if !rules[rulegraphNodePath]() {
									goto l326
								}
							l328:
								{
									position329, tokenIndex329, depth329 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l329
									}
									goto l328
								l329:
									position, tokenIndex, depth = position329, tokenIndex329, depth329
								}
								if !rules[ruleRPAREN]() {
									goto l326
								}
								depth--
								add(rulecollectionPath, position327)
							}
							goto l325
						l326:
							position, tokenIndex, depth = position325, tokenIndex325, depth325
							{
								position330 := position
								depth++
								{
									position331 := position
									depth++
									if buffer[position] != rune('[') {
										goto l311
									}
									position++
									if !rules[ruleskip]() {
										goto l311
									}
									depth--
									add(ruleLBRACK, position331)
								}
								if !rules[rulepropertyListPath]() {
									goto l311
								}
								{
									position332 := position
									depth++
									if buffer[position] != rune(']') {
										goto l311
									}
									position++
									if !rules[ruleskip]() {
										goto l311
									}
									depth--
									add(ruleRBRACK, position332)
								}
								depth--
								add(ruleblankNodePropertyListPath, position330)
							}
						}
					l325:
						depth--
						add(ruletriplesNodePath, position324)
					}
					if !rules[rulepropertyListPath]() {
						goto l311
					}
				}
			l313:
				depth--
				add(ruletriplesSameSubjectPath, position312)
			}
			return true
		l311:
			position, tokenIndex, depth = position311, tokenIndex311, depth311
			return false
		},
		/* 26 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 27 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position334, tokenIndex334, depth334 := position, tokenIndex, depth
			{
				position335 := position
				depth++
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l337
					}
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					{
						switch buffer[position] {
						case '(':
							if !rules[rulenil]() {
								goto l334
							}
							break
						case '[', '_':
							{
								position339 := position
								depth++
								{
									position340, tokenIndex340, depth340 := position, tokenIndex, depth
									{
										position342 := position
										depth++
										if buffer[position] != rune('_') {
											goto l341
										}
										position++
										if buffer[position] != rune(':') {
											goto l341
										}
										position++
										{
											switch buffer[position] {
											case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l341
												}
												position++
												break
											case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l341
												}
												position++
												break
											default:
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l341
												}
												position++
												break
											}
										}

										{
											position344, tokenIndex344, depth344 := position, tokenIndex, depth
											{
												position346, tokenIndex346, depth346 := position, tokenIndex, depth
												if c := buffer[position]; c < rune('a') || c > rune('z') {
													goto l347
												}
												position++
												goto l346
											l347:
												position, tokenIndex, depth = position346, tokenIndex346, depth346
												if c := buffer[position]; c < rune('A') || c > rune('Z') {
													goto l348
												}
												position++
												goto l346
											l348:
												position, tokenIndex, depth = position346, tokenIndex346, depth346
												if c := buffer[position]; c < rune('0') || c > rune('9') {
													goto l349
												}
												position++
												goto l346
											l349:
												position, tokenIndex, depth = position346, tokenIndex346, depth346
												if c := buffer[position]; c < rune('.') || c > rune('_') {
													goto l344
												}
												position++
											}
										l346:
											goto l345
										l344:
											position, tokenIndex, depth = position344, tokenIndex344, depth344
										}
									l345:
										if !rules[ruleskip]() {
											goto l341
										}
										depth--
										add(ruleblankNodeLabel, position342)
									}
									goto l340
								l341:
									position, tokenIndex, depth = position340, tokenIndex340, depth340
									{
										position350 := position
										depth++
										if buffer[position] != rune('[') {
											goto l334
										}
										position++
									l351:
										{
											position352, tokenIndex352, depth352 := position, tokenIndex, depth
											if !rules[rulews]() {
												goto l352
											}
											goto l351
										l352:
											position, tokenIndex, depth = position352, tokenIndex352, depth352
										}
										if buffer[position] != rune(']') {
											goto l334
										}
										position++
										if !rules[ruleskip]() {
											goto l334
										}
										depth--
										add(ruleanon, position350)
									}
								}
							l340:
								depth--
								add(ruleblankNode, position339)
							}
							break
						case 'F', 'T', 'f', 't':
							if !rules[rulebooleanLiteral]() {
								goto l334
							}
							break
						case '"':
							if !rules[ruleliteral]() {
								goto l334
							}
							break
						default:
							if !rules[rulenumericLiteral]() {
								goto l334
							}
							break
						}
					}

				}
			l336:
				depth--
				add(rulegraphTerm, position335)
			}
			return true
		l334:
			position, tokenIndex, depth = position334, tokenIndex334, depth334
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
			position356, tokenIndex356, depth356 := position, tokenIndex, depth
			{
				position357 := position
				depth++
				{
					position358, tokenIndex358, depth358 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l359
					}
					{
						add(ruleAction4, position)
					}
					goto l358
				l359:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					{
						position362 := position
						depth++
						if !rules[rulevar]() {
							goto l361
						}
						depth--
						add(rulePegText, position362)
					}
					{
						add(ruleAction5, position)
					}
					goto l358
				l361:
					position, tokenIndex, depth = position358, tokenIndex358, depth358
					{
						position364 := position
						depth++
						if !rules[rulepath]() {
							goto l356
						}
						depth--
						add(ruleverbPath, position364)
					}
				}
			l358:
				if !rules[ruleobjectListPath]() {
					goto l356
				}
				{
					position365, tokenIndex365, depth365 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l365
					}
					if !rules[rulepropertyListPath]() {
						goto l365
					}
					goto l366
				l365:
					position, tokenIndex, depth = position365, tokenIndex365, depth365
				}
			l366:
				depth--
				add(rulepropertyListPath, position357)
			}
			return true
		l356:
			position, tokenIndex, depth = position356, tokenIndex356, depth356
			return false
		},
		/* 32 verbPath <- <path> */
		nil,
		/* 33 path <- <pathAlternative> */
		func() bool {
			position368, tokenIndex368, depth368 := position, tokenIndex, depth
			{
				position369 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l368
				}
				depth--
				add(rulepath, position369)
			}
			return true
		l368:
			position, tokenIndex, depth = position368, tokenIndex368, depth368
			return false
		},
		/* 34 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position370, tokenIndex370, depth370 := position, tokenIndex, depth
			{
				position371 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l370
				}
			l372:
				{
					position373, tokenIndex373, depth373 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l373
					}
					if !rules[rulepathAlternative]() {
						goto l373
					}
					goto l372
				l373:
					position, tokenIndex, depth = position373, tokenIndex373, depth373
				}
				depth--
				add(rulepathAlternative, position371)
			}
			return true
		l370:
			position, tokenIndex, depth = position370, tokenIndex370, depth370
			return false
		},
		/* 35 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				{
					position376 := position
					depth++
					{
						position377 := position
						depth++
						{
							position378, tokenIndex378, depth378 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l378
							}
							goto l379
						l378:
							position, tokenIndex, depth = position378, tokenIndex378, depth378
						}
					l379:
						{
							position380 := position
							depth++
							{
								position381, tokenIndex381, depth381 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l382
								}
								goto l381
							l382:
								position, tokenIndex, depth = position381, tokenIndex381, depth381
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l374
										}
										if !rules[rulepath]() {
											goto l374
										}
										if !rules[ruleRPAREN]() {
											goto l374
										}
										break
									case '!':
										if !rules[ruleNOT]() {
											goto l374
										}
										{
											position384 := position
											depth++
											{
												position385, tokenIndex385, depth385 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l386
												}
												goto l385
											l386:
												position, tokenIndex, depth = position385, tokenIndex385, depth385
												if !rules[ruleLPAREN]() {
													goto l374
												}
												{
													position387, tokenIndex387, depth387 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l387
													}
												l389:
													{
														position390, tokenIndex390, depth390 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l390
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l390
														}
														goto l389
													l390:
														position, tokenIndex, depth = position390, tokenIndex390, depth390
													}
													goto l388
												l387:
													position, tokenIndex, depth = position387, tokenIndex387, depth387
												}
											l388:
												if !rules[ruleRPAREN]() {
													goto l374
												}
											}
										l385:
											depth--
											add(rulepathNegatedPropertySet, position384)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l374
										}
										break
									}
								}

							}
						l381:
							depth--
							add(rulepathPrimary, position380)
						}
						depth--
						add(rulepathElt, position377)
					}
					depth--
					add(rulePegText, position376)
				}
				{
					add(ruleAction6, position)
				}
			l392:
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l393
					}
					if !rules[rulepathSequence]() {
						goto l393
					}
					goto l392
				l393:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
				}
				depth--
				add(rulepathSequence, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
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
			position397, tokenIndex397, depth397 := position, tokenIndex, depth
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
					if !rules[ruleISA]() {
						goto l401
					}
					goto l399
				l401:
					position, tokenIndex, depth = position399, tokenIndex399, depth399
					if !rules[ruleINVERSE]() {
						goto l397
					}
					{
						position402, tokenIndex402, depth402 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l403
						}
						goto l402
					l403:
						position, tokenIndex, depth = position402, tokenIndex402, depth402
						if !rules[ruleISA]() {
							goto l397
						}
					}
				l402:
				}
			l399:
				depth--
				add(rulepathOneInPropertySet, position398)
			}
			return true
		l397:
			position, tokenIndex, depth = position397, tokenIndex397, depth397
			return false
		},
		/* 40 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			{
				position405 := position
				depth++
				{
					position406 := position
					depth++
					{
						position407, tokenIndex407, depth407 := position, tokenIndex, depth
						if !rules[rulepof]() {
							goto l408
						}
						{
							add(ruleAction7, position)
						}
						goto l407
					l408:
						position, tokenIndex, depth = position407, tokenIndex407, depth407
						{
							position411 := position
							depth++
							if !rules[rulegraphNodePath]() {
								goto l410
							}
							depth--
							add(rulePegText, position411)
						}
						{
							add(ruleAction8, position)
						}
						goto l407
					l410:
						position, tokenIndex, depth = position407, tokenIndex407, depth407
						{
							add(ruleAction9, position)
						}
					}
				l407:
					depth--
					add(ruleobjectPath, position406)
				}
			l414:
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if !rules[ruleCOMMA]() {
						goto l415
					}
					if !rules[ruleobjectListPath]() {
						goto l415
					}
					goto l414
				l415:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
				}
				depth--
				add(ruleobjectListPath, position405)
			}
			return true
		},
		/* 41 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		nil,
		/* 42 graphNodePath <- <(var / graphTerm)> */
		func() bool {
			position417, tokenIndex417, depth417 := position, tokenIndex, depth
			{
				position418 := position
				depth++
				{
					position419, tokenIndex419, depth419 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l420
					}
					goto l419
				l420:
					position, tokenIndex, depth = position419, tokenIndex419, depth419
					if !rules[rulegraphTerm]() {
						goto l417
					}
				}
			l419:
				depth--
				add(rulegraphNodePath, position418)
			}
			return true
		l417:
			position, tokenIndex, depth = position417, tokenIndex417, depth417
			return false
		},
		/* 43 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position422 := position
				depth++
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					{
						position425 := position
						depth++
						{
							position426, tokenIndex426, depth426 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l427
							}
							{
								position428, tokenIndex428, depth428 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l428
								}
								goto l429
							l428:
								position, tokenIndex, depth = position428, tokenIndex428, depth428
							}
						l429:
							goto l426
						l427:
							position, tokenIndex, depth = position426, tokenIndex426, depth426
							if !rules[ruleoffset]() {
								goto l423
							}
							{
								position430, tokenIndex430, depth430 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l430
								}
								goto l431
							l430:
								position, tokenIndex, depth = position430, tokenIndex430, depth430
							}
						l431:
						}
					l426:
						depth--
						add(rulelimitOffsetClauses, position425)
					}
					goto l424
				l423:
					position, tokenIndex, depth = position423, tokenIndex423, depth423
				}
			l424:
				depth--
				add(rulesolutionModifier, position422)
			}
			return true
		},
		/* 44 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 45 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position433, tokenIndex433, depth433 := position, tokenIndex, depth
			{
				position434 := position
				depth++
				{
					position435 := position
					depth++
					{
						position436, tokenIndex436, depth436 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l437
						}
						position++
						goto l436
					l437:
						position, tokenIndex, depth = position436, tokenIndex436, depth436
						if buffer[position] != rune('L') {
							goto l433
						}
						position++
					}
				l436:
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
							goto l433
						}
						position++
					}
				l438:
					{
						position440, tokenIndex440, depth440 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l441
						}
						position++
						goto l440
					l441:
						position, tokenIndex, depth = position440, tokenIndex440, depth440
						if buffer[position] != rune('M') {
							goto l433
						}
						position++
					}
				l440:
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
							goto l433
						}
						position++
					}
				l442:
					{
						position444, tokenIndex444, depth444 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l445
						}
						position++
						goto l444
					l445:
						position, tokenIndex, depth = position444, tokenIndex444, depth444
						if buffer[position] != rune('T') {
							goto l433
						}
						position++
					}
				l444:
					if !rules[ruleskip]() {
						goto l433
					}
					depth--
					add(ruleLIMIT, position435)
				}
				if !rules[ruleINTEGER]() {
					goto l433
				}
				depth--
				add(rulelimit, position434)
			}
			return true
		l433:
			position, tokenIndex, depth = position433, tokenIndex433, depth433
			return false
		},
		/* 46 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position446, tokenIndex446, depth446 := position, tokenIndex, depth
			{
				position447 := position
				depth++
				{
					position448 := position
					depth++
					{
						position449, tokenIndex449, depth449 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l450
						}
						position++
						goto l449
					l450:
						position, tokenIndex, depth = position449, tokenIndex449, depth449
						if buffer[position] != rune('O') {
							goto l446
						}
						position++
					}
				l449:
					{
						position451, tokenIndex451, depth451 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l452
						}
						position++
						goto l451
					l452:
						position, tokenIndex, depth = position451, tokenIndex451, depth451
						if buffer[position] != rune('F') {
							goto l446
						}
						position++
					}
				l451:
					{
						position453, tokenIndex453, depth453 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l454
						}
						position++
						goto l453
					l454:
						position, tokenIndex, depth = position453, tokenIndex453, depth453
						if buffer[position] != rune('F') {
							goto l446
						}
						position++
					}
				l453:
					{
						position455, tokenIndex455, depth455 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l456
						}
						position++
						goto l455
					l456:
						position, tokenIndex, depth = position455, tokenIndex455, depth455
						if buffer[position] != rune('S') {
							goto l446
						}
						position++
					}
				l455:
					{
						position457, tokenIndex457, depth457 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l458
						}
						position++
						goto l457
					l458:
						position, tokenIndex, depth = position457, tokenIndex457, depth457
						if buffer[position] != rune('E') {
							goto l446
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
							goto l446
						}
						position++
					}
				l459:
					if !rules[ruleskip]() {
						goto l446
					}
					depth--
					add(ruleOFFSET, position448)
				}
				if !rules[ruleINTEGER]() {
					goto l446
				}
				depth--
				add(ruleoffset, position447)
			}
			return true
		l446:
			position, tokenIndex, depth = position446, tokenIndex446, depth446
			return false
		},
		/* 47 expression <- <conditionalOrExpression> */
		func() bool {
			position461, tokenIndex461, depth461 := position, tokenIndex, depth
			{
				position462 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l461
				}
				depth--
				add(ruleexpression, position462)
			}
			return true
		l461:
			position, tokenIndex, depth = position461, tokenIndex461, depth461
			return false
		},
		/* 48 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l463
				}
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					{
						position467 := position
						depth++
						if buffer[position] != rune('|') {
							goto l465
						}
						position++
						if buffer[position] != rune('|') {
							goto l465
						}
						position++
						if !rules[ruleskip]() {
							goto l465
						}
						depth--
						add(ruleOR, position467)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l465
					}
					goto l466
				l465:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
				}
			l466:
				depth--
				add(ruleconditionalOrExpression, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
			return false
		},
		/* 49 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position468, tokenIndex468, depth468 := position, tokenIndex, depth
			{
				position469 := position
				depth++
				{
					position470 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l468
					}
					{
						position471, tokenIndex471, depth471 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position474 := position
									depth++
									{
										position475 := position
										depth++
										{
											position476, tokenIndex476, depth476 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l477
											}
											position++
											goto l476
										l477:
											position, tokenIndex, depth = position476, tokenIndex476, depth476
											if buffer[position] != rune('N') {
												goto l471
											}
											position++
										}
									l476:
										{
											position478, tokenIndex478, depth478 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l479
											}
											position++
											goto l478
										l479:
											position, tokenIndex, depth = position478, tokenIndex478, depth478
											if buffer[position] != rune('O') {
												goto l471
											}
											position++
										}
									l478:
										{
											position480, tokenIndex480, depth480 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l481
											}
											position++
											goto l480
										l481:
											position, tokenIndex, depth = position480, tokenIndex480, depth480
											if buffer[position] != rune('T') {
												goto l471
											}
											position++
										}
									l480:
										if buffer[position] != rune(' ') {
											goto l471
										}
										position++
										{
											position482, tokenIndex482, depth482 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l483
											}
											position++
											goto l482
										l483:
											position, tokenIndex, depth = position482, tokenIndex482, depth482
											if buffer[position] != rune('I') {
												goto l471
											}
											position++
										}
									l482:
										{
											position484, tokenIndex484, depth484 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l485
											}
											position++
											goto l484
										l485:
											position, tokenIndex, depth = position484, tokenIndex484, depth484
											if buffer[position] != rune('N') {
												goto l471
											}
											position++
										}
									l484:
										if !rules[ruleskip]() {
											goto l471
										}
										depth--
										add(ruleNOTIN, position475)
									}
									if !rules[ruleargList]() {
										goto l471
									}
									depth--
									add(rulenotin, position474)
								}
								break
							case 'I', 'i':
								{
									position486 := position
									depth++
									{
										position487 := position
										depth++
										{
											position488, tokenIndex488, depth488 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l489
											}
											position++
											goto l488
										l489:
											position, tokenIndex, depth = position488, tokenIndex488, depth488
											if buffer[position] != rune('I') {
												goto l471
											}
											position++
										}
									l488:
										{
											position490, tokenIndex490, depth490 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l491
											}
											position++
											goto l490
										l491:
											position, tokenIndex, depth = position490, tokenIndex490, depth490
											if buffer[position] != rune('N') {
												goto l471
											}
											position++
										}
									l490:
										if !rules[ruleskip]() {
											goto l471
										}
										depth--
										add(ruleIN, position487)
									}
									if !rules[ruleargList]() {
										goto l471
									}
									depth--
									add(rulein, position486)
								}
								break
							default:
								{
									position492, tokenIndex492, depth492 := position, tokenIndex, depth
									{
										position494 := position
										depth++
										if buffer[position] != rune('<') {
											goto l493
										}
										position++
										if !rules[ruleskip]() {
											goto l493
										}
										depth--
										add(ruleLT, position494)
									}
									goto l492
								l493:
									position, tokenIndex, depth = position492, tokenIndex492, depth492
									{
										position496 := position
										depth++
										if buffer[position] != rune('>') {
											goto l495
										}
										position++
										if buffer[position] != rune('=') {
											goto l495
										}
										position++
										if !rules[ruleskip]() {
											goto l495
										}
										depth--
										add(ruleGE, position496)
									}
									goto l492
								l495:
									position, tokenIndex, depth = position492, tokenIndex492, depth492
									{
										switch buffer[position] {
										case '>':
											{
												position498 := position
												depth++
												if buffer[position] != rune('>') {
													goto l471
												}
												position++
												if !rules[ruleskip]() {
													goto l471
												}
												depth--
												add(ruleGT, position498)
											}
											break
										case '<':
											{
												position499 := position
												depth++
												if buffer[position] != rune('<') {
													goto l471
												}
												position++
												if buffer[position] != rune('=') {
													goto l471
												}
												position++
												if !rules[ruleskip]() {
													goto l471
												}
												depth--
												add(ruleLE, position499)
											}
											break
										case '!':
											{
												position500 := position
												depth++
												if buffer[position] != rune('!') {
													goto l471
												}
												position++
												if buffer[position] != rune('=') {
													goto l471
												}
												position++
												if !rules[ruleskip]() {
													goto l471
												}
												depth--
												add(ruleNE, position500)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l471
											}
											break
										}
									}

								}
							l492:
								if !rules[rulenumericExpression]() {
									goto l471
								}
								break
							}
						}

						goto l472
					l471:
						position, tokenIndex, depth = position471, tokenIndex471, depth471
					}
				l472:
					depth--
					add(rulevalueLogical, position470)
				}
				{
					position501, tokenIndex501, depth501 := position, tokenIndex, depth
					{
						position503 := position
						depth++
						if buffer[position] != rune('&') {
							goto l501
						}
						position++
						if buffer[position] != rune('&') {
							goto l501
						}
						position++
						if !rules[ruleskip]() {
							goto l501
						}
						depth--
						add(ruleAND, position503)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l501
					}
					goto l502
				l501:
					position, tokenIndex, depth = position501, tokenIndex501, depth501
				}
			l502:
				depth--
				add(ruleconditionalAndExpression, position469)
			}
			return true
		l468:
			position, tokenIndex, depth = position468, tokenIndex468, depth468
			return false
		},
		/* 50 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 51 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position505, tokenIndex505, depth505 := position, tokenIndex, depth
			{
				position506 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l505
				}
			l507:
				{
					position508, tokenIndex508, depth508 := position, tokenIndex, depth
					{
						position509, tokenIndex509, depth509 := position, tokenIndex, depth
						{
							position511, tokenIndex511, depth511 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l512
							}
							goto l511
						l512:
							position, tokenIndex, depth = position511, tokenIndex511, depth511
							if !rules[ruleMINUS]() {
								goto l510
							}
						}
					l511:
						if !rules[rulemultiplicativeExpression]() {
							goto l510
						}
						goto l509
					l510:
						position, tokenIndex, depth = position509, tokenIndex509, depth509
						{
							position513 := position
							depth++
							{
								position514, tokenIndex514, depth514 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l515
								}
								position++
								goto l514
							l515:
								position, tokenIndex, depth = position514, tokenIndex514, depth514
								if buffer[position] != rune('-') {
									goto l508
								}
								position++
							}
						l514:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l508
							}
							position++
						l516:
							{
								position517, tokenIndex517, depth517 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l517
								}
								position++
								goto l516
							l517:
								position, tokenIndex, depth = position517, tokenIndex517, depth517
							}
							{
								position518, tokenIndex518, depth518 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l518
								}
								position++
							l520:
								{
									position521, tokenIndex521, depth521 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l521
									}
									position++
									goto l520
								l521:
									position, tokenIndex, depth = position521, tokenIndex521, depth521
								}
								goto l519
							l518:
								position, tokenIndex, depth = position518, tokenIndex518, depth518
							}
						l519:
							if !rules[ruleskip]() {
								goto l508
							}
							depth--
							add(rulesignedNumericLiteral, position513)
						}
					}
				l509:
					goto l507
				l508:
					position, tokenIndex, depth = position508, tokenIndex508, depth508
				}
				depth--
				add(rulenumericExpression, position506)
			}
			return true
		l505:
			position, tokenIndex, depth = position505, tokenIndex505, depth505
			return false
		},
		/* 52 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position522, tokenIndex522, depth522 := position, tokenIndex, depth
			{
				position523 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l522
				}
			l524:
				{
					position525, tokenIndex525, depth525 := position, tokenIndex, depth
					{
						position526, tokenIndex526, depth526 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l527
						}
						goto l526
					l527:
						position, tokenIndex, depth = position526, tokenIndex526, depth526
						if !rules[ruleSLASH]() {
							goto l525
						}
					}
				l526:
					if !rules[ruleunaryExpression]() {
						goto l525
					}
					goto l524
				l525:
					position, tokenIndex, depth = position525, tokenIndex525, depth525
				}
				depth--
				add(rulemultiplicativeExpression, position523)
			}
			return true
		l522:
			position, tokenIndex, depth = position522, tokenIndex522, depth522
			return false
		},
		/* 53 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position528, tokenIndex528, depth528 := position, tokenIndex, depth
			{
				position529 := position
				depth++
				{
					position530, tokenIndex530, depth530 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l530
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l530
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l530
							}
							break
						}
					}

					goto l531
				l530:
					position, tokenIndex, depth = position530, tokenIndex530, depth530
				}
			l531:
				{
					position533 := position
					depth++
					{
						position534, tokenIndex534, depth534 := position, tokenIndex, depth
						if !rules[rulebrackettedExpression]() {
							goto l535
						}
						goto l534
					l535:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if !rules[rulebuiltinCall]() {
							goto l536
						}
						goto l534
					l536:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if !rules[rulefunctionCall]() {
							goto l537
						}
						goto l534
					l537:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if !rules[ruleiriref]() {
							goto l538
						}
						goto l534
					l538:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position540 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position542 := position
												depth++
												{
													position543 := position
													depth++
													{
														position544, tokenIndex544, depth544 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l545
														}
														position++
														goto l544
													l545:
														position, tokenIndex, depth = position544, tokenIndex544, depth544
														if buffer[position] != rune('G') {
															goto l528
														}
														position++
													}
												l544:
													{
														position546, tokenIndex546, depth546 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l547
														}
														position++
														goto l546
													l547:
														position, tokenIndex, depth = position546, tokenIndex546, depth546
														if buffer[position] != rune('R') {
															goto l528
														}
														position++
													}
												l546:
													{
														position548, tokenIndex548, depth548 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l549
														}
														position++
														goto l548
													l549:
														position, tokenIndex, depth = position548, tokenIndex548, depth548
														if buffer[position] != rune('O') {
															goto l528
														}
														position++
													}
												l548:
													{
														position550, tokenIndex550, depth550 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l551
														}
														position++
														goto l550
													l551:
														position, tokenIndex, depth = position550, tokenIndex550, depth550
														if buffer[position] != rune('U') {
															goto l528
														}
														position++
													}
												l550:
													{
														position552, tokenIndex552, depth552 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l553
														}
														position++
														goto l552
													l553:
														position, tokenIndex, depth = position552, tokenIndex552, depth552
														if buffer[position] != rune('P') {
															goto l528
														}
														position++
													}
												l552:
													if buffer[position] != rune('_') {
														goto l528
													}
													position++
													{
														position554, tokenIndex554, depth554 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l555
														}
														position++
														goto l554
													l555:
														position, tokenIndex, depth = position554, tokenIndex554, depth554
														if buffer[position] != rune('C') {
															goto l528
														}
														position++
													}
												l554:
													{
														position556, tokenIndex556, depth556 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l557
														}
														position++
														goto l556
													l557:
														position, tokenIndex, depth = position556, tokenIndex556, depth556
														if buffer[position] != rune('O') {
															goto l528
														}
														position++
													}
												l556:
													{
														position558, tokenIndex558, depth558 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l559
														}
														position++
														goto l558
													l559:
														position, tokenIndex, depth = position558, tokenIndex558, depth558
														if buffer[position] != rune('N') {
															goto l528
														}
														position++
													}
												l558:
													{
														position560, tokenIndex560, depth560 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l561
														}
														position++
														goto l560
													l561:
														position, tokenIndex, depth = position560, tokenIndex560, depth560
														if buffer[position] != rune('C') {
															goto l528
														}
														position++
													}
												l560:
													{
														position562, tokenIndex562, depth562 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l563
														}
														position++
														goto l562
													l563:
														position, tokenIndex, depth = position562, tokenIndex562, depth562
														if buffer[position] != rune('A') {
															goto l528
														}
														position++
													}
												l562:
													{
														position564, tokenIndex564, depth564 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l565
														}
														position++
														goto l564
													l565:
														position, tokenIndex, depth = position564, tokenIndex564, depth564
														if buffer[position] != rune('T') {
															goto l528
														}
														position++
													}
												l564:
													if !rules[ruleskip]() {
														goto l528
													}
													depth--
													add(ruleGROUPCONCAT, position543)
												}
												if !rules[ruleLPAREN]() {
													goto l528
												}
												{
													position566, tokenIndex566, depth566 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l566
													}
													goto l567
												l566:
													position, tokenIndex, depth = position566, tokenIndex566, depth566
												}
											l567:
												if !rules[ruleexpression]() {
													goto l528
												}
												{
													position568, tokenIndex568, depth568 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l568
													}
													{
														position570 := position
														depth++
														{
															position571, tokenIndex571, depth571 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l572
															}
															position++
															goto l571
														l572:
															position, tokenIndex, depth = position571, tokenIndex571, depth571
															if buffer[position] != rune('S') {
																goto l568
															}
															position++
														}
													l571:
														{
															position573, tokenIndex573, depth573 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l574
															}
															position++
															goto l573
														l574:
															position, tokenIndex, depth = position573, tokenIndex573, depth573
															if buffer[position] != rune('E') {
																goto l568
															}
															position++
														}
													l573:
														{
															position575, tokenIndex575, depth575 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l576
															}
															position++
															goto l575
														l576:
															position, tokenIndex, depth = position575, tokenIndex575, depth575
															if buffer[position] != rune('P') {
																goto l568
															}
															position++
														}
													l575:
														{
															position577, tokenIndex577, depth577 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l578
															}
															position++
															goto l577
														l578:
															position, tokenIndex, depth = position577, tokenIndex577, depth577
															if buffer[position] != rune('A') {
																goto l568
															}
															position++
														}
													l577:
														{
															position579, tokenIndex579, depth579 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l580
															}
															position++
															goto l579
														l580:
															position, tokenIndex, depth = position579, tokenIndex579, depth579
															if buffer[position] != rune('R') {
																goto l568
															}
															position++
														}
													l579:
														{
															position581, tokenIndex581, depth581 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l582
															}
															position++
															goto l581
														l582:
															position, tokenIndex, depth = position581, tokenIndex581, depth581
															if buffer[position] != rune('A') {
																goto l568
															}
															position++
														}
													l581:
														{
															position583, tokenIndex583, depth583 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l584
															}
															position++
															goto l583
														l584:
															position, tokenIndex, depth = position583, tokenIndex583, depth583
															if buffer[position] != rune('T') {
																goto l568
															}
															position++
														}
													l583:
														{
															position585, tokenIndex585, depth585 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l586
															}
															position++
															goto l585
														l586:
															position, tokenIndex, depth = position585, tokenIndex585, depth585
															if buffer[position] != rune('O') {
																goto l568
															}
															position++
														}
													l585:
														{
															position587, tokenIndex587, depth587 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l588
															}
															position++
															goto l587
														l588:
															position, tokenIndex, depth = position587, tokenIndex587, depth587
															if buffer[position] != rune('R') {
																goto l568
															}
															position++
														}
													l587:
														if !rules[ruleskip]() {
															goto l568
														}
														depth--
														add(ruleSEPARATOR, position570)
													}
													if !rules[ruleEQ]() {
														goto l568
													}
													if !rules[rulestring]() {
														goto l568
													}
													goto l569
												l568:
													position, tokenIndex, depth = position568, tokenIndex568, depth568
												}
											l569:
												if !rules[ruleRPAREN]() {
													goto l528
												}
												depth--
												add(rulegroupConcat, position542)
											}
											break
										case 'C', 'c':
											{
												position589 := position
												depth++
												{
													position590 := position
													depth++
													{
														position591, tokenIndex591, depth591 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l592
														}
														position++
														goto l591
													l592:
														position, tokenIndex, depth = position591, tokenIndex591, depth591
														if buffer[position] != rune('C') {
															goto l528
														}
														position++
													}
												l591:
													{
														position593, tokenIndex593, depth593 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l594
														}
														position++
														goto l593
													l594:
														position, tokenIndex, depth = position593, tokenIndex593, depth593
														if buffer[position] != rune('O') {
															goto l528
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
															goto l528
														}
														position++
													}
												l595:
													{
														position597, tokenIndex597, depth597 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l598
														}
														position++
														goto l597
													l598:
														position, tokenIndex, depth = position597, tokenIndex597, depth597
														if buffer[position] != rune('N') {
															goto l528
														}
														position++
													}
												l597:
													{
														position599, tokenIndex599, depth599 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l600
														}
														position++
														goto l599
													l600:
														position, tokenIndex, depth = position599, tokenIndex599, depth599
														if buffer[position] != rune('T') {
															goto l528
														}
														position++
													}
												l599:
													if !rules[ruleskip]() {
														goto l528
													}
													depth--
													add(ruleCOUNT, position590)
												}
												if !rules[ruleLPAREN]() {
													goto l528
												}
												{
													position601, tokenIndex601, depth601 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l601
													}
													goto l602
												l601:
													position, tokenIndex, depth = position601, tokenIndex601, depth601
												}
											l602:
												{
													position603, tokenIndex603, depth603 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l604
													}
													goto l603
												l604:
													position, tokenIndex, depth = position603, tokenIndex603, depth603
													if !rules[ruleexpression]() {
														goto l528
													}
												}
											l603:
												if !rules[ruleRPAREN]() {
													goto l528
												}
												depth--
												add(rulecount, position589)
											}
											break
										default:
											{
												position605, tokenIndex605, depth605 := position, tokenIndex, depth
												{
													position607 := position
													depth++
													{
														position608, tokenIndex608, depth608 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l609
														}
														position++
														goto l608
													l609:
														position, tokenIndex, depth = position608, tokenIndex608, depth608
														if buffer[position] != rune('S') {
															goto l606
														}
														position++
													}
												l608:
													{
														position610, tokenIndex610, depth610 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l611
														}
														position++
														goto l610
													l611:
														position, tokenIndex, depth = position610, tokenIndex610, depth610
														if buffer[position] != rune('U') {
															goto l606
														}
														position++
													}
												l610:
													{
														position612, tokenIndex612, depth612 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l613
														}
														position++
														goto l612
													l613:
														position, tokenIndex, depth = position612, tokenIndex612, depth612
														if buffer[position] != rune('M') {
															goto l606
														}
														position++
													}
												l612:
													if !rules[ruleskip]() {
														goto l606
													}
													depth--
													add(ruleSUM, position607)
												}
												goto l605
											l606:
												position, tokenIndex, depth = position605, tokenIndex605, depth605
												{
													position615 := position
													depth++
													{
														position616, tokenIndex616, depth616 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l617
														}
														position++
														goto l616
													l617:
														position, tokenIndex, depth = position616, tokenIndex616, depth616
														if buffer[position] != rune('M') {
															goto l614
														}
														position++
													}
												l616:
													{
														position618, tokenIndex618, depth618 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l619
														}
														position++
														goto l618
													l619:
														position, tokenIndex, depth = position618, tokenIndex618, depth618
														if buffer[position] != rune('I') {
															goto l614
														}
														position++
													}
												l618:
													{
														position620, tokenIndex620, depth620 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l621
														}
														position++
														goto l620
													l621:
														position, tokenIndex, depth = position620, tokenIndex620, depth620
														if buffer[position] != rune('N') {
															goto l614
														}
														position++
													}
												l620:
													if !rules[ruleskip]() {
														goto l614
													}
													depth--
													add(ruleMIN, position615)
												}
												goto l605
											l614:
												position, tokenIndex, depth = position605, tokenIndex605, depth605
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position623 := position
															depth++
															{
																position624, tokenIndex624, depth624 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l625
																}
																position++
																goto l624
															l625:
																position, tokenIndex, depth = position624, tokenIndex624, depth624
																if buffer[position] != rune('S') {
																	goto l528
																}
																position++
															}
														l624:
															{
																position626, tokenIndex626, depth626 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l627
																}
																position++
																goto l626
															l627:
																position, tokenIndex, depth = position626, tokenIndex626, depth626
																if buffer[position] != rune('A') {
																	goto l528
																}
																position++
															}
														l626:
															{
																position628, tokenIndex628, depth628 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l629
																}
																position++
																goto l628
															l629:
																position, tokenIndex, depth = position628, tokenIndex628, depth628
																if buffer[position] != rune('M') {
																	goto l528
																}
																position++
															}
														l628:
															{
																position630, tokenIndex630, depth630 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l631
																}
																position++
																goto l630
															l631:
																position, tokenIndex, depth = position630, tokenIndex630, depth630
																if buffer[position] != rune('P') {
																	goto l528
																}
																position++
															}
														l630:
															{
																position632, tokenIndex632, depth632 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l633
																}
																position++
																goto l632
															l633:
																position, tokenIndex, depth = position632, tokenIndex632, depth632
																if buffer[position] != rune('L') {
																	goto l528
																}
																position++
															}
														l632:
															{
																position634, tokenIndex634, depth634 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l635
																}
																position++
																goto l634
															l635:
																position, tokenIndex, depth = position634, tokenIndex634, depth634
																if buffer[position] != rune('E') {
																	goto l528
																}
																position++
															}
														l634:
															if !rules[ruleskip]() {
																goto l528
															}
															depth--
															add(ruleSAMPLE, position623)
														}
														break
													case 'A', 'a':
														{
															position636 := position
															depth++
															{
																position637, tokenIndex637, depth637 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l638
																}
																position++
																goto l637
															l638:
																position, tokenIndex, depth = position637, tokenIndex637, depth637
																if buffer[position] != rune('A') {
																	goto l528
																}
																position++
															}
														l637:
															{
																position639, tokenIndex639, depth639 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l640
																}
																position++
																goto l639
															l640:
																position, tokenIndex, depth = position639, tokenIndex639, depth639
																if buffer[position] != rune('V') {
																	goto l528
																}
																position++
															}
														l639:
															{
																position641, tokenIndex641, depth641 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l642
																}
																position++
																goto l641
															l642:
																position, tokenIndex, depth = position641, tokenIndex641, depth641
																if buffer[position] != rune('G') {
																	goto l528
																}
																position++
															}
														l641:
															if !rules[ruleskip]() {
																goto l528
															}
															depth--
															add(ruleAVG, position636)
														}
														break
													default:
														{
															position643 := position
															depth++
															{
																position644, tokenIndex644, depth644 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l645
																}
																position++
																goto l644
															l645:
																position, tokenIndex, depth = position644, tokenIndex644, depth644
																if buffer[position] != rune('M') {
																	goto l528
																}
																position++
															}
														l644:
															{
																position646, tokenIndex646, depth646 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l647
																}
																position++
																goto l646
															l647:
																position, tokenIndex, depth = position646, tokenIndex646, depth646
																if buffer[position] != rune('A') {
																	goto l528
																}
																position++
															}
														l646:
															{
																position648, tokenIndex648, depth648 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l649
																}
																position++
																goto l648
															l649:
																position, tokenIndex, depth = position648, tokenIndex648, depth648
																if buffer[position] != rune('X') {
																	goto l528
																}
																position++
															}
														l648:
															if !rules[ruleskip]() {
																goto l528
															}
															depth--
															add(ruleMAX, position643)
														}
														break
													}
												}

											}
										l605:
											if !rules[ruleLPAREN]() {
												goto l528
											}
											{
												position650, tokenIndex650, depth650 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l650
												}
												goto l651
											l650:
												position, tokenIndex, depth = position650, tokenIndex650, depth650
											}
										l651:
											if !rules[ruleexpression]() {
												goto l528
											}
											if !rules[ruleRPAREN]() {
												goto l528
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position540)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l528
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l528
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l528
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l528
								}
								break
							}
						}

					}
				l534:
					depth--
					add(ruleprimaryExpression, position533)
				}
				depth--
				add(ruleunaryExpression, position529)
			}
			return true
		l528:
			position, tokenIndex, depth = position528, tokenIndex528, depth528
			return false
		},
		/* 54 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 55 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position653, tokenIndex653, depth653 := position, tokenIndex, depth
			{
				position654 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l653
				}
				if !rules[ruleexpression]() {
					goto l653
				}
				if !rules[ruleRPAREN]() {
					goto l653
				}
				depth--
				add(rulebrackettedExpression, position654)
			}
			return true
		l653:
			position, tokenIndex, depth = position653, tokenIndex653, depth653
			return false
		},
		/* 56 functionCall <- <(iriref argList)> */
		func() bool {
			position655, tokenIndex655, depth655 := position, tokenIndex, depth
			{
				position656 := position
				depth++
				if !rules[ruleiriref]() {
					goto l655
				}
				if !rules[ruleargList]() {
					goto l655
				}
				depth--
				add(rulefunctionCall, position656)
			}
			return true
		l655:
			position, tokenIndex, depth = position655, tokenIndex655, depth655
			return false
		},
		/* 57 in <- <(IN argList)> */
		nil,
		/* 58 notin <- <(NOTIN argList)> */
		nil,
		/* 59 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position659, tokenIndex659, depth659 := position, tokenIndex, depth
			{
				position660 := position
				depth++
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l662
					}
					goto l661
				l662:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
					if !rules[ruleLPAREN]() {
						goto l659
					}
					if !rules[ruleexpression]() {
						goto l659
					}
				l663:
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l664
						}
						if !rules[ruleexpression]() {
							goto l664
						}
						goto l663
					l664:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
					}
					if !rules[ruleRPAREN]() {
						goto l659
					}
				}
			l661:
				depth--
				add(ruleargList, position660)
			}
			return true
		l659:
			position, tokenIndex, depth = position659, tokenIndex659, depth659
			return false
		},
		/* 60 aggregate <- <((&('G' | 'g') groupConcat) | (&('C' | 'c') count) | (&('A' | 'M' | 'S' | 'a' | 'm' | 's') ((SUM / MIN / ((&('S' | 's') SAMPLE) | (&('A' | 'a') AVG) | (&('M' | 'm') MAX))) LPAREN DISTINCT? expression RPAREN)))> */
		nil,
		/* 61 count <- <(COUNT LPAREN DISTINCT? (STAR / expression) RPAREN)> */
		nil,
		/* 62 groupConcat <- <(GROUPCONCAT LPAREN DISTINCT? expression (SEMICOLON SEPARATOR EQ string)? RPAREN)> */
		nil,
		/* 63 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		func() bool {
			position668, tokenIndex668, depth668 := position, tokenIndex, depth
			{
				position669 := position
				depth++
				{
					position670, tokenIndex670, depth670 := position, tokenIndex, depth
					{
						position672, tokenIndex672, depth672 := position, tokenIndex, depth
						{
							position674 := position
							depth++
							{
								position675, tokenIndex675, depth675 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l676
								}
								position++
								goto l675
							l676:
								position, tokenIndex, depth = position675, tokenIndex675, depth675
								if buffer[position] != rune('S') {
									goto l673
								}
								position++
							}
						l675:
							{
								position677, tokenIndex677, depth677 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l678
								}
								position++
								goto l677
							l678:
								position, tokenIndex, depth = position677, tokenIndex677, depth677
								if buffer[position] != rune('T') {
									goto l673
								}
								position++
							}
						l677:
							{
								position679, tokenIndex679, depth679 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l680
								}
								position++
								goto l679
							l680:
								position, tokenIndex, depth = position679, tokenIndex679, depth679
								if buffer[position] != rune('R') {
									goto l673
								}
								position++
							}
						l679:
							if !rules[ruleskip]() {
								goto l673
							}
							depth--
							add(ruleSTR, position674)
						}
						goto l672
					l673:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position682 := position
							depth++
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
									goto l681
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
									goto l681
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
									goto l681
								}
								position++
							}
						l687:
							{
								position689, tokenIndex689, depth689 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l690
								}
								position++
								goto l689
							l690:
								position, tokenIndex, depth = position689, tokenIndex689, depth689
								if buffer[position] != rune('G') {
									goto l681
								}
								position++
							}
						l689:
							if !rules[ruleskip]() {
								goto l681
							}
							depth--
							add(ruleLANG, position682)
						}
						goto l672
					l681:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position692 := position
							depth++
							{
								position693, tokenIndex693, depth693 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l694
								}
								position++
								goto l693
							l694:
								position, tokenIndex, depth = position693, tokenIndex693, depth693
								if buffer[position] != rune('D') {
									goto l691
								}
								position++
							}
						l693:
							{
								position695, tokenIndex695, depth695 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l696
								}
								position++
								goto l695
							l696:
								position, tokenIndex, depth = position695, tokenIndex695, depth695
								if buffer[position] != rune('A') {
									goto l691
								}
								position++
							}
						l695:
							{
								position697, tokenIndex697, depth697 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l698
								}
								position++
								goto l697
							l698:
								position, tokenIndex, depth = position697, tokenIndex697, depth697
								if buffer[position] != rune('T') {
									goto l691
								}
								position++
							}
						l697:
							{
								position699, tokenIndex699, depth699 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l700
								}
								position++
								goto l699
							l700:
								position, tokenIndex, depth = position699, tokenIndex699, depth699
								if buffer[position] != rune('A') {
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
								if buffer[position] != rune('y') {
									goto l704
								}
								position++
								goto l703
							l704:
								position, tokenIndex, depth = position703, tokenIndex703, depth703
								if buffer[position] != rune('Y') {
									goto l691
								}
								position++
							}
						l703:
							{
								position705, tokenIndex705, depth705 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l706
								}
								position++
								goto l705
							l706:
								position, tokenIndex, depth = position705, tokenIndex705, depth705
								if buffer[position] != rune('P') {
									goto l691
								}
								position++
							}
						l705:
							{
								position707, tokenIndex707, depth707 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l708
								}
								position++
								goto l707
							l708:
								position, tokenIndex, depth = position707, tokenIndex707, depth707
								if buffer[position] != rune('E') {
									goto l691
								}
								position++
							}
						l707:
							if !rules[ruleskip]() {
								goto l691
							}
							depth--
							add(ruleDATATYPE, position692)
						}
						goto l672
					l691:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position710 := position
							depth++
							{
								position711, tokenIndex711, depth711 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l712
								}
								position++
								goto l711
							l712:
								position, tokenIndex, depth = position711, tokenIndex711, depth711
								if buffer[position] != rune('I') {
									goto l709
								}
								position++
							}
						l711:
							{
								position713, tokenIndex713, depth713 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l714
								}
								position++
								goto l713
							l714:
								position, tokenIndex, depth = position713, tokenIndex713, depth713
								if buffer[position] != rune('R') {
									goto l709
								}
								position++
							}
						l713:
							{
								position715, tokenIndex715, depth715 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l716
								}
								position++
								goto l715
							l716:
								position, tokenIndex, depth = position715, tokenIndex715, depth715
								if buffer[position] != rune('I') {
									goto l709
								}
								position++
							}
						l715:
							if !rules[ruleskip]() {
								goto l709
							}
							depth--
							add(ruleIRI, position710)
						}
						goto l672
					l709:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position718 := position
							depth++
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
									goto l717
								}
								position++
							}
						l719:
							{
								position721, tokenIndex721, depth721 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l722
								}
								position++
								goto l721
							l722:
								position, tokenIndex, depth = position721, tokenIndex721, depth721
								if buffer[position] != rune('R') {
									goto l717
								}
								position++
							}
						l721:
							{
								position723, tokenIndex723, depth723 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l724
								}
								position++
								goto l723
							l724:
								position, tokenIndex, depth = position723, tokenIndex723, depth723
								if buffer[position] != rune('I') {
									goto l717
								}
								position++
							}
						l723:
							if !rules[ruleskip]() {
								goto l717
							}
							depth--
							add(ruleURI, position718)
						}
						goto l672
					l717:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position726 := position
							depth++
							{
								position727, tokenIndex727, depth727 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l728
								}
								position++
								goto l727
							l728:
								position, tokenIndex, depth = position727, tokenIndex727, depth727
								if buffer[position] != rune('S') {
									goto l725
								}
								position++
							}
						l727:
							{
								position729, tokenIndex729, depth729 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l730
								}
								position++
								goto l729
							l730:
								position, tokenIndex, depth = position729, tokenIndex729, depth729
								if buffer[position] != rune('T') {
									goto l725
								}
								position++
							}
						l729:
							{
								position731, tokenIndex731, depth731 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l732
								}
								position++
								goto l731
							l732:
								position, tokenIndex, depth = position731, tokenIndex731, depth731
								if buffer[position] != rune('R') {
									goto l725
								}
								position++
							}
						l731:
							{
								position733, tokenIndex733, depth733 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l734
								}
								position++
								goto l733
							l734:
								position, tokenIndex, depth = position733, tokenIndex733, depth733
								if buffer[position] != rune('L') {
									goto l725
								}
								position++
							}
						l733:
							{
								position735, tokenIndex735, depth735 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l736
								}
								position++
								goto l735
							l736:
								position, tokenIndex, depth = position735, tokenIndex735, depth735
								if buffer[position] != rune('E') {
									goto l725
								}
								position++
							}
						l735:
							{
								position737, tokenIndex737, depth737 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l738
								}
								position++
								goto l737
							l738:
								position, tokenIndex, depth = position737, tokenIndex737, depth737
								if buffer[position] != rune('N') {
									goto l725
								}
								position++
							}
						l737:
							if !rules[ruleskip]() {
								goto l725
							}
							depth--
							add(ruleSTRLEN, position726)
						}
						goto l672
					l725:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position740 := position
							depth++
							{
								position741, tokenIndex741, depth741 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l742
								}
								position++
								goto l741
							l742:
								position, tokenIndex, depth = position741, tokenIndex741, depth741
								if buffer[position] != rune('M') {
									goto l739
								}
								position++
							}
						l741:
							{
								position743, tokenIndex743, depth743 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l744
								}
								position++
								goto l743
							l744:
								position, tokenIndex, depth = position743, tokenIndex743, depth743
								if buffer[position] != rune('O') {
									goto l739
								}
								position++
							}
						l743:
							{
								position745, tokenIndex745, depth745 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l746
								}
								position++
								goto l745
							l746:
								position, tokenIndex, depth = position745, tokenIndex745, depth745
								if buffer[position] != rune('N') {
									goto l739
								}
								position++
							}
						l745:
							{
								position747, tokenIndex747, depth747 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l748
								}
								position++
								goto l747
							l748:
								position, tokenIndex, depth = position747, tokenIndex747, depth747
								if buffer[position] != rune('T') {
									goto l739
								}
								position++
							}
						l747:
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
									goto l739
								}
								position++
							}
						l749:
							if !rules[ruleskip]() {
								goto l739
							}
							depth--
							add(ruleMONTH, position740)
						}
						goto l672
					l739:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position752 := position
							depth++
							{
								position753, tokenIndex753, depth753 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l754
								}
								position++
								goto l753
							l754:
								position, tokenIndex, depth = position753, tokenIndex753, depth753
								if buffer[position] != rune('M') {
									goto l751
								}
								position++
							}
						l753:
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
									goto l751
								}
								position++
							}
						l755:
							{
								position757, tokenIndex757, depth757 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l758
								}
								position++
								goto l757
							l758:
								position, tokenIndex, depth = position757, tokenIndex757, depth757
								if buffer[position] != rune('N') {
									goto l751
								}
								position++
							}
						l757:
							{
								position759, tokenIndex759, depth759 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l760
								}
								position++
								goto l759
							l760:
								position, tokenIndex, depth = position759, tokenIndex759, depth759
								if buffer[position] != rune('U') {
									goto l751
								}
								position++
							}
						l759:
							{
								position761, tokenIndex761, depth761 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l762
								}
								position++
								goto l761
							l762:
								position, tokenIndex, depth = position761, tokenIndex761, depth761
								if buffer[position] != rune('T') {
									goto l751
								}
								position++
							}
						l761:
							{
								position763, tokenIndex763, depth763 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l764
								}
								position++
								goto l763
							l764:
								position, tokenIndex, depth = position763, tokenIndex763, depth763
								if buffer[position] != rune('E') {
									goto l751
								}
								position++
							}
						l763:
							{
								position765, tokenIndex765, depth765 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l766
								}
								position++
								goto l765
							l766:
								position, tokenIndex, depth = position765, tokenIndex765, depth765
								if buffer[position] != rune('S') {
									goto l751
								}
								position++
							}
						l765:
							if !rules[ruleskip]() {
								goto l751
							}
							depth--
							add(ruleMINUTES, position752)
						}
						goto l672
					l751:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position768 := position
							depth++
							{
								position769, tokenIndex769, depth769 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l770
								}
								position++
								goto l769
							l770:
								position, tokenIndex, depth = position769, tokenIndex769, depth769
								if buffer[position] != rune('S') {
									goto l767
								}
								position++
							}
						l769:
							{
								position771, tokenIndex771, depth771 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l772
								}
								position++
								goto l771
							l772:
								position, tokenIndex, depth = position771, tokenIndex771, depth771
								if buffer[position] != rune('E') {
									goto l767
								}
								position++
							}
						l771:
							{
								position773, tokenIndex773, depth773 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l774
								}
								position++
								goto l773
							l774:
								position, tokenIndex, depth = position773, tokenIndex773, depth773
								if buffer[position] != rune('C') {
									goto l767
								}
								position++
							}
						l773:
							{
								position775, tokenIndex775, depth775 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l776
								}
								position++
								goto l775
							l776:
								position, tokenIndex, depth = position775, tokenIndex775, depth775
								if buffer[position] != rune('O') {
									goto l767
								}
								position++
							}
						l775:
							{
								position777, tokenIndex777, depth777 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l778
								}
								position++
								goto l777
							l778:
								position, tokenIndex, depth = position777, tokenIndex777, depth777
								if buffer[position] != rune('N') {
									goto l767
								}
								position++
							}
						l777:
							{
								position779, tokenIndex779, depth779 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l780
								}
								position++
								goto l779
							l780:
								position, tokenIndex, depth = position779, tokenIndex779, depth779
								if buffer[position] != rune('D') {
									goto l767
								}
								position++
							}
						l779:
							{
								position781, tokenIndex781, depth781 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l782
								}
								position++
								goto l781
							l782:
								position, tokenIndex, depth = position781, tokenIndex781, depth781
								if buffer[position] != rune('S') {
									goto l767
								}
								position++
							}
						l781:
							if !rules[ruleskip]() {
								goto l767
							}
							depth--
							add(ruleSECONDS, position768)
						}
						goto l672
					l767:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position784 := position
							depth++
							{
								position785, tokenIndex785, depth785 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l786
								}
								position++
								goto l785
							l786:
								position, tokenIndex, depth = position785, tokenIndex785, depth785
								if buffer[position] != rune('T') {
									goto l783
								}
								position++
							}
						l785:
							{
								position787, tokenIndex787, depth787 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l788
								}
								position++
								goto l787
							l788:
								position, tokenIndex, depth = position787, tokenIndex787, depth787
								if buffer[position] != rune('I') {
									goto l783
								}
								position++
							}
						l787:
							{
								position789, tokenIndex789, depth789 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l790
								}
								position++
								goto l789
							l790:
								position, tokenIndex, depth = position789, tokenIndex789, depth789
								if buffer[position] != rune('M') {
									goto l783
								}
								position++
							}
						l789:
							{
								position791, tokenIndex791, depth791 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l792
								}
								position++
								goto l791
							l792:
								position, tokenIndex, depth = position791, tokenIndex791, depth791
								if buffer[position] != rune('E') {
									goto l783
								}
								position++
							}
						l791:
							{
								position793, tokenIndex793, depth793 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l794
								}
								position++
								goto l793
							l794:
								position, tokenIndex, depth = position793, tokenIndex793, depth793
								if buffer[position] != rune('Z') {
									goto l783
								}
								position++
							}
						l793:
							{
								position795, tokenIndex795, depth795 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l796
								}
								position++
								goto l795
							l796:
								position, tokenIndex, depth = position795, tokenIndex795, depth795
								if buffer[position] != rune('O') {
									goto l783
								}
								position++
							}
						l795:
							{
								position797, tokenIndex797, depth797 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l798
								}
								position++
								goto l797
							l798:
								position, tokenIndex, depth = position797, tokenIndex797, depth797
								if buffer[position] != rune('N') {
									goto l783
								}
								position++
							}
						l797:
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
									goto l783
								}
								position++
							}
						l799:
							if !rules[ruleskip]() {
								goto l783
							}
							depth--
							add(ruleTIMEZONE, position784)
						}
						goto l672
					l783:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position802 := position
							depth++
							{
								position803, tokenIndex803, depth803 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l804
								}
								position++
								goto l803
							l804:
								position, tokenIndex, depth = position803, tokenIndex803, depth803
								if buffer[position] != rune('S') {
									goto l801
								}
								position++
							}
						l803:
							{
								position805, tokenIndex805, depth805 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l806
								}
								position++
								goto l805
							l806:
								position, tokenIndex, depth = position805, tokenIndex805, depth805
								if buffer[position] != rune('H') {
									goto l801
								}
								position++
							}
						l805:
							{
								position807, tokenIndex807, depth807 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l808
								}
								position++
								goto l807
							l808:
								position, tokenIndex, depth = position807, tokenIndex807, depth807
								if buffer[position] != rune('A') {
									goto l801
								}
								position++
							}
						l807:
							if buffer[position] != rune('1') {
								goto l801
							}
							position++
							if !rules[ruleskip]() {
								goto l801
							}
							depth--
							add(ruleSHA1, position802)
						}
						goto l672
					l801:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position810 := position
							depth++
							{
								position811, tokenIndex811, depth811 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l812
								}
								position++
								goto l811
							l812:
								position, tokenIndex, depth = position811, tokenIndex811, depth811
								if buffer[position] != rune('S') {
									goto l809
								}
								position++
							}
						l811:
							{
								position813, tokenIndex813, depth813 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l814
								}
								position++
								goto l813
							l814:
								position, tokenIndex, depth = position813, tokenIndex813, depth813
								if buffer[position] != rune('H') {
									goto l809
								}
								position++
							}
						l813:
							{
								position815, tokenIndex815, depth815 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l816
								}
								position++
								goto l815
							l816:
								position, tokenIndex, depth = position815, tokenIndex815, depth815
								if buffer[position] != rune('A') {
									goto l809
								}
								position++
							}
						l815:
							if buffer[position] != rune('2') {
								goto l809
							}
							position++
							if buffer[position] != rune('5') {
								goto l809
							}
							position++
							if buffer[position] != rune('6') {
								goto l809
							}
							position++
							if !rules[ruleskip]() {
								goto l809
							}
							depth--
							add(ruleSHA256, position810)
						}
						goto l672
					l809:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position818 := position
							depth++
							{
								position819, tokenIndex819, depth819 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l820
								}
								position++
								goto l819
							l820:
								position, tokenIndex, depth = position819, tokenIndex819, depth819
								if buffer[position] != rune('S') {
									goto l817
								}
								position++
							}
						l819:
							{
								position821, tokenIndex821, depth821 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l822
								}
								position++
								goto l821
							l822:
								position, tokenIndex, depth = position821, tokenIndex821, depth821
								if buffer[position] != rune('H') {
									goto l817
								}
								position++
							}
						l821:
							{
								position823, tokenIndex823, depth823 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l824
								}
								position++
								goto l823
							l824:
								position, tokenIndex, depth = position823, tokenIndex823, depth823
								if buffer[position] != rune('A') {
									goto l817
								}
								position++
							}
						l823:
							if buffer[position] != rune('3') {
								goto l817
							}
							position++
							if buffer[position] != rune('8') {
								goto l817
							}
							position++
							if buffer[position] != rune('4') {
								goto l817
							}
							position++
							if !rules[ruleskip]() {
								goto l817
							}
							depth--
							add(ruleSHA384, position818)
						}
						goto l672
					l817:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position826 := position
							depth++
							{
								position827, tokenIndex827, depth827 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l828
								}
								position++
								goto l827
							l828:
								position, tokenIndex, depth = position827, tokenIndex827, depth827
								if buffer[position] != rune('I') {
									goto l825
								}
								position++
							}
						l827:
							{
								position829, tokenIndex829, depth829 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l830
								}
								position++
								goto l829
							l830:
								position, tokenIndex, depth = position829, tokenIndex829, depth829
								if buffer[position] != rune('S') {
									goto l825
								}
								position++
							}
						l829:
							{
								position831, tokenIndex831, depth831 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l832
								}
								position++
								goto l831
							l832:
								position, tokenIndex, depth = position831, tokenIndex831, depth831
								if buffer[position] != rune('I') {
									goto l825
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
									goto l825
								}
								position++
							}
						l833:
							{
								position835, tokenIndex835, depth835 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l836
								}
								position++
								goto l835
							l836:
								position, tokenIndex, depth = position835, tokenIndex835, depth835
								if buffer[position] != rune('I') {
									goto l825
								}
								position++
							}
						l835:
							if !rules[ruleskip]() {
								goto l825
							}
							depth--
							add(ruleISIRI, position826)
						}
						goto l672
					l825:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position838 := position
							depth++
							{
								position839, tokenIndex839, depth839 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l840
								}
								position++
								goto l839
							l840:
								position, tokenIndex, depth = position839, tokenIndex839, depth839
								if buffer[position] != rune('I') {
									goto l837
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
									goto l837
								}
								position++
							}
						l841:
							{
								position843, tokenIndex843, depth843 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l844
								}
								position++
								goto l843
							l844:
								position, tokenIndex, depth = position843, tokenIndex843, depth843
								if buffer[position] != rune('U') {
									goto l837
								}
								position++
							}
						l843:
							{
								position845, tokenIndex845, depth845 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l846
								}
								position++
								goto l845
							l846:
								position, tokenIndex, depth = position845, tokenIndex845, depth845
								if buffer[position] != rune('R') {
									goto l837
								}
								position++
							}
						l845:
							{
								position847, tokenIndex847, depth847 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l848
								}
								position++
								goto l847
							l848:
								position, tokenIndex, depth = position847, tokenIndex847, depth847
								if buffer[position] != rune('I') {
									goto l837
								}
								position++
							}
						l847:
							if !rules[ruleskip]() {
								goto l837
							}
							depth--
							add(ruleISURI, position838)
						}
						goto l672
					l837:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position850 := position
							depth++
							{
								position851, tokenIndex851, depth851 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l852
								}
								position++
								goto l851
							l852:
								position, tokenIndex, depth = position851, tokenIndex851, depth851
								if buffer[position] != rune('I') {
									goto l849
								}
								position++
							}
						l851:
							{
								position853, tokenIndex853, depth853 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l854
								}
								position++
								goto l853
							l854:
								position, tokenIndex, depth = position853, tokenIndex853, depth853
								if buffer[position] != rune('S') {
									goto l849
								}
								position++
							}
						l853:
							{
								position855, tokenIndex855, depth855 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l856
								}
								position++
								goto l855
							l856:
								position, tokenIndex, depth = position855, tokenIndex855, depth855
								if buffer[position] != rune('B') {
									goto l849
								}
								position++
							}
						l855:
							{
								position857, tokenIndex857, depth857 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l858
								}
								position++
								goto l857
							l858:
								position, tokenIndex, depth = position857, tokenIndex857, depth857
								if buffer[position] != rune('L') {
									goto l849
								}
								position++
							}
						l857:
							{
								position859, tokenIndex859, depth859 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l860
								}
								position++
								goto l859
							l860:
								position, tokenIndex, depth = position859, tokenIndex859, depth859
								if buffer[position] != rune('A') {
									goto l849
								}
								position++
							}
						l859:
							{
								position861, tokenIndex861, depth861 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l862
								}
								position++
								goto l861
							l862:
								position, tokenIndex, depth = position861, tokenIndex861, depth861
								if buffer[position] != rune('N') {
									goto l849
								}
								position++
							}
						l861:
							{
								position863, tokenIndex863, depth863 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l864
								}
								position++
								goto l863
							l864:
								position, tokenIndex, depth = position863, tokenIndex863, depth863
								if buffer[position] != rune('K') {
									goto l849
								}
								position++
							}
						l863:
							if !rules[ruleskip]() {
								goto l849
							}
							depth--
							add(ruleISBLANK, position850)
						}
						goto l672
					l849:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							position866 := position
							depth++
							{
								position867, tokenIndex867, depth867 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l868
								}
								position++
								goto l867
							l868:
								position, tokenIndex, depth = position867, tokenIndex867, depth867
								if buffer[position] != rune('I') {
									goto l865
								}
								position++
							}
						l867:
							{
								position869, tokenIndex869, depth869 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l870
								}
								position++
								goto l869
							l870:
								position, tokenIndex, depth = position869, tokenIndex869, depth869
								if buffer[position] != rune('S') {
									goto l865
								}
								position++
							}
						l869:
							{
								position871, tokenIndex871, depth871 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l872
								}
								position++
								goto l871
							l872:
								position, tokenIndex, depth = position871, tokenIndex871, depth871
								if buffer[position] != rune('L') {
									goto l865
								}
								position++
							}
						l871:
							{
								position873, tokenIndex873, depth873 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l874
								}
								position++
								goto l873
							l874:
								position, tokenIndex, depth = position873, tokenIndex873, depth873
								if buffer[position] != rune('I') {
									goto l865
								}
								position++
							}
						l873:
							{
								position875, tokenIndex875, depth875 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l876
								}
								position++
								goto l875
							l876:
								position, tokenIndex, depth = position875, tokenIndex875, depth875
								if buffer[position] != rune('T') {
									goto l865
								}
								position++
							}
						l875:
							{
								position877, tokenIndex877, depth877 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l878
								}
								position++
								goto l877
							l878:
								position, tokenIndex, depth = position877, tokenIndex877, depth877
								if buffer[position] != rune('E') {
									goto l865
								}
								position++
							}
						l877:
							{
								position879, tokenIndex879, depth879 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l880
								}
								position++
								goto l879
							l880:
								position, tokenIndex, depth = position879, tokenIndex879, depth879
								if buffer[position] != rune('R') {
									goto l865
								}
								position++
							}
						l879:
							{
								position881, tokenIndex881, depth881 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l882
								}
								position++
								goto l881
							l882:
								position, tokenIndex, depth = position881, tokenIndex881, depth881
								if buffer[position] != rune('A') {
									goto l865
								}
								position++
							}
						l881:
							{
								position883, tokenIndex883, depth883 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l884
								}
								position++
								goto l883
							l884:
								position, tokenIndex, depth = position883, tokenIndex883, depth883
								if buffer[position] != rune('L') {
									goto l865
								}
								position++
							}
						l883:
							if !rules[ruleskip]() {
								goto l865
							}
							depth--
							add(ruleISLITERAL, position866)
						}
						goto l672
					l865:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position886 := position
									depth++
									{
										position887, tokenIndex887, depth887 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l888
										}
										position++
										goto l887
									l888:
										position, tokenIndex, depth = position887, tokenIndex887, depth887
										if buffer[position] != rune('I') {
											goto l671
										}
										position++
									}
								l887:
									{
										position889, tokenIndex889, depth889 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l890
										}
										position++
										goto l889
									l890:
										position, tokenIndex, depth = position889, tokenIndex889, depth889
										if buffer[position] != rune('S') {
											goto l671
										}
										position++
									}
								l889:
									{
										position891, tokenIndex891, depth891 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l892
										}
										position++
										goto l891
									l892:
										position, tokenIndex, depth = position891, tokenIndex891, depth891
										if buffer[position] != rune('N') {
											goto l671
										}
										position++
									}
								l891:
									{
										position893, tokenIndex893, depth893 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l894
										}
										position++
										goto l893
									l894:
										position, tokenIndex, depth = position893, tokenIndex893, depth893
										if buffer[position] != rune('U') {
											goto l671
										}
										position++
									}
								l893:
									{
										position895, tokenIndex895, depth895 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l896
										}
										position++
										goto l895
									l896:
										position, tokenIndex, depth = position895, tokenIndex895, depth895
										if buffer[position] != rune('M') {
											goto l671
										}
										position++
									}
								l895:
									{
										position897, tokenIndex897, depth897 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l898
										}
										position++
										goto l897
									l898:
										position, tokenIndex, depth = position897, tokenIndex897, depth897
										if buffer[position] != rune('E') {
											goto l671
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
											goto l671
										}
										position++
									}
								l899:
									{
										position901, tokenIndex901, depth901 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l902
										}
										position++
										goto l901
									l902:
										position, tokenIndex, depth = position901, tokenIndex901, depth901
										if buffer[position] != rune('I') {
											goto l671
										}
										position++
									}
								l901:
									{
										position903, tokenIndex903, depth903 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l904
										}
										position++
										goto l903
									l904:
										position, tokenIndex, depth = position903, tokenIndex903, depth903
										if buffer[position] != rune('C') {
											goto l671
										}
										position++
									}
								l903:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleISNUMERIC, position886)
								}
								break
							case 'S', 's':
								{
									position905 := position
									depth++
									{
										position906, tokenIndex906, depth906 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l907
										}
										position++
										goto l906
									l907:
										position, tokenIndex, depth = position906, tokenIndex906, depth906
										if buffer[position] != rune('S') {
											goto l671
										}
										position++
									}
								l906:
									{
										position908, tokenIndex908, depth908 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l909
										}
										position++
										goto l908
									l909:
										position, tokenIndex, depth = position908, tokenIndex908, depth908
										if buffer[position] != rune('H') {
											goto l671
										}
										position++
									}
								l908:
									{
										position910, tokenIndex910, depth910 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l911
										}
										position++
										goto l910
									l911:
										position, tokenIndex, depth = position910, tokenIndex910, depth910
										if buffer[position] != rune('A') {
											goto l671
										}
										position++
									}
								l910:
									if buffer[position] != rune('5') {
										goto l671
									}
									position++
									if buffer[position] != rune('1') {
										goto l671
									}
									position++
									if buffer[position] != rune('2') {
										goto l671
									}
									position++
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleSHA512, position905)
								}
								break
							case 'M', 'm':
								{
									position912 := position
									depth++
									{
										position913, tokenIndex913, depth913 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l914
										}
										position++
										goto l913
									l914:
										position, tokenIndex, depth = position913, tokenIndex913, depth913
										if buffer[position] != rune('M') {
											goto l671
										}
										position++
									}
								l913:
									{
										position915, tokenIndex915, depth915 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l916
										}
										position++
										goto l915
									l916:
										position, tokenIndex, depth = position915, tokenIndex915, depth915
										if buffer[position] != rune('D') {
											goto l671
										}
										position++
									}
								l915:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleMD5, position912)
								}
								break
							case 'T', 't':
								{
									position917 := position
									depth++
									{
										position918, tokenIndex918, depth918 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l919
										}
										position++
										goto l918
									l919:
										position, tokenIndex, depth = position918, tokenIndex918, depth918
										if buffer[position] != rune('T') {
											goto l671
										}
										position++
									}
								l918:
									{
										position920, tokenIndex920, depth920 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l921
										}
										position++
										goto l920
									l921:
										position, tokenIndex, depth = position920, tokenIndex920, depth920
										if buffer[position] != rune('Z') {
											goto l671
										}
										position++
									}
								l920:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleTZ, position917)
								}
								break
							case 'H', 'h':
								{
									position922 := position
									depth++
									{
										position923, tokenIndex923, depth923 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l924
										}
										position++
										goto l923
									l924:
										position, tokenIndex, depth = position923, tokenIndex923, depth923
										if buffer[position] != rune('H') {
											goto l671
										}
										position++
									}
								l923:
									{
										position925, tokenIndex925, depth925 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l926
										}
										position++
										goto l925
									l926:
										position, tokenIndex, depth = position925, tokenIndex925, depth925
										if buffer[position] != rune('O') {
											goto l671
										}
										position++
									}
								l925:
									{
										position927, tokenIndex927, depth927 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l928
										}
										position++
										goto l927
									l928:
										position, tokenIndex, depth = position927, tokenIndex927, depth927
										if buffer[position] != rune('U') {
											goto l671
										}
										position++
									}
								l927:
									{
										position929, tokenIndex929, depth929 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l930
										}
										position++
										goto l929
									l930:
										position, tokenIndex, depth = position929, tokenIndex929, depth929
										if buffer[position] != rune('R') {
											goto l671
										}
										position++
									}
								l929:
									{
										position931, tokenIndex931, depth931 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l932
										}
										position++
										goto l931
									l932:
										position, tokenIndex, depth = position931, tokenIndex931, depth931
										if buffer[position] != rune('S') {
											goto l671
										}
										position++
									}
								l931:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleHOURS, position922)
								}
								break
							case 'D', 'd':
								{
									position933 := position
									depth++
									{
										position934, tokenIndex934, depth934 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l935
										}
										position++
										goto l934
									l935:
										position, tokenIndex, depth = position934, tokenIndex934, depth934
										if buffer[position] != rune('D') {
											goto l671
										}
										position++
									}
								l934:
									{
										position936, tokenIndex936, depth936 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l937
										}
										position++
										goto l936
									l937:
										position, tokenIndex, depth = position936, tokenIndex936, depth936
										if buffer[position] != rune('A') {
											goto l671
										}
										position++
									}
								l936:
									{
										position938, tokenIndex938, depth938 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l939
										}
										position++
										goto l938
									l939:
										position, tokenIndex, depth = position938, tokenIndex938, depth938
										if buffer[position] != rune('Y') {
											goto l671
										}
										position++
									}
								l938:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleDAY, position933)
								}
								break
							case 'Y', 'y':
								{
									position940 := position
									depth++
									{
										position941, tokenIndex941, depth941 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l942
										}
										position++
										goto l941
									l942:
										position, tokenIndex, depth = position941, tokenIndex941, depth941
										if buffer[position] != rune('Y') {
											goto l671
										}
										position++
									}
								l941:
									{
										position943, tokenIndex943, depth943 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l944
										}
										position++
										goto l943
									l944:
										position, tokenIndex, depth = position943, tokenIndex943, depth943
										if buffer[position] != rune('E') {
											goto l671
										}
										position++
									}
								l943:
									{
										position945, tokenIndex945, depth945 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l946
										}
										position++
										goto l945
									l946:
										position, tokenIndex, depth = position945, tokenIndex945, depth945
										if buffer[position] != rune('A') {
											goto l671
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
											goto l671
										}
										position++
									}
								l947:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleYEAR, position940)
								}
								break
							case 'E', 'e':
								{
									position949 := position
									depth++
									{
										position950, tokenIndex950, depth950 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l951
										}
										position++
										goto l950
									l951:
										position, tokenIndex, depth = position950, tokenIndex950, depth950
										if buffer[position] != rune('E') {
											goto l671
										}
										position++
									}
								l950:
									{
										position952, tokenIndex952, depth952 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l953
										}
										position++
										goto l952
									l953:
										position, tokenIndex, depth = position952, tokenIndex952, depth952
										if buffer[position] != rune('N') {
											goto l671
										}
										position++
									}
								l952:
									{
										position954, tokenIndex954, depth954 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l955
										}
										position++
										goto l954
									l955:
										position, tokenIndex, depth = position954, tokenIndex954, depth954
										if buffer[position] != rune('C') {
											goto l671
										}
										position++
									}
								l954:
									{
										position956, tokenIndex956, depth956 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l957
										}
										position++
										goto l956
									l957:
										position, tokenIndex, depth = position956, tokenIndex956, depth956
										if buffer[position] != rune('O') {
											goto l671
										}
										position++
									}
								l956:
									{
										position958, tokenIndex958, depth958 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l959
										}
										position++
										goto l958
									l959:
										position, tokenIndex, depth = position958, tokenIndex958, depth958
										if buffer[position] != rune('D') {
											goto l671
										}
										position++
									}
								l958:
									{
										position960, tokenIndex960, depth960 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l961
										}
										position++
										goto l960
									l961:
										position, tokenIndex, depth = position960, tokenIndex960, depth960
										if buffer[position] != rune('E') {
											goto l671
										}
										position++
									}
								l960:
									if buffer[position] != rune('_') {
										goto l671
									}
									position++
									{
										position962, tokenIndex962, depth962 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l963
										}
										position++
										goto l962
									l963:
										position, tokenIndex, depth = position962, tokenIndex962, depth962
										if buffer[position] != rune('F') {
											goto l671
										}
										position++
									}
								l962:
									{
										position964, tokenIndex964, depth964 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l965
										}
										position++
										goto l964
									l965:
										position, tokenIndex, depth = position964, tokenIndex964, depth964
										if buffer[position] != rune('O') {
											goto l671
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
											goto l671
										}
										position++
									}
								l966:
									if buffer[position] != rune('_') {
										goto l671
									}
									position++
									{
										position968, tokenIndex968, depth968 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l969
										}
										position++
										goto l968
									l969:
										position, tokenIndex, depth = position968, tokenIndex968, depth968
										if buffer[position] != rune('U') {
											goto l671
										}
										position++
									}
								l968:
									{
										position970, tokenIndex970, depth970 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l971
										}
										position++
										goto l970
									l971:
										position, tokenIndex, depth = position970, tokenIndex970, depth970
										if buffer[position] != rune('R') {
											goto l671
										}
										position++
									}
								l970:
									{
										position972, tokenIndex972, depth972 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l973
										}
										position++
										goto l972
									l973:
										position, tokenIndex, depth = position972, tokenIndex972, depth972
										if buffer[position] != rune('I') {
											goto l671
										}
										position++
									}
								l972:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleENCODEFORURI, position949)
								}
								break
							case 'L', 'l':
								{
									position974 := position
									depth++
									{
										position975, tokenIndex975, depth975 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l976
										}
										position++
										goto l975
									l976:
										position, tokenIndex, depth = position975, tokenIndex975, depth975
										if buffer[position] != rune('L') {
											goto l671
										}
										position++
									}
								l975:
									{
										position977, tokenIndex977, depth977 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l978
										}
										position++
										goto l977
									l978:
										position, tokenIndex, depth = position977, tokenIndex977, depth977
										if buffer[position] != rune('C') {
											goto l671
										}
										position++
									}
								l977:
									{
										position979, tokenIndex979, depth979 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l980
										}
										position++
										goto l979
									l980:
										position, tokenIndex, depth = position979, tokenIndex979, depth979
										if buffer[position] != rune('A') {
											goto l671
										}
										position++
									}
								l979:
									{
										position981, tokenIndex981, depth981 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l982
										}
										position++
										goto l981
									l982:
										position, tokenIndex, depth = position981, tokenIndex981, depth981
										if buffer[position] != rune('S') {
											goto l671
										}
										position++
									}
								l981:
									{
										position983, tokenIndex983, depth983 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l984
										}
										position++
										goto l983
									l984:
										position, tokenIndex, depth = position983, tokenIndex983, depth983
										if buffer[position] != rune('E') {
											goto l671
										}
										position++
									}
								l983:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleLCASE, position974)
								}
								break
							case 'U', 'u':
								{
									position985 := position
									depth++
									{
										position986, tokenIndex986, depth986 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l987
										}
										position++
										goto l986
									l987:
										position, tokenIndex, depth = position986, tokenIndex986, depth986
										if buffer[position] != rune('U') {
											goto l671
										}
										position++
									}
								l986:
									{
										position988, tokenIndex988, depth988 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l989
										}
										position++
										goto l988
									l989:
										position, tokenIndex, depth = position988, tokenIndex988, depth988
										if buffer[position] != rune('C') {
											goto l671
										}
										position++
									}
								l988:
									{
										position990, tokenIndex990, depth990 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l991
										}
										position++
										goto l990
									l991:
										position, tokenIndex, depth = position990, tokenIndex990, depth990
										if buffer[position] != rune('A') {
											goto l671
										}
										position++
									}
								l990:
									{
										position992, tokenIndex992, depth992 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l993
										}
										position++
										goto l992
									l993:
										position, tokenIndex, depth = position992, tokenIndex992, depth992
										if buffer[position] != rune('S') {
											goto l671
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
											goto l671
										}
										position++
									}
								l994:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleUCASE, position985)
								}
								break
							case 'F', 'f':
								{
									position996 := position
									depth++
									{
										position997, tokenIndex997, depth997 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l998
										}
										position++
										goto l997
									l998:
										position, tokenIndex, depth = position997, tokenIndex997, depth997
										if buffer[position] != rune('F') {
											goto l671
										}
										position++
									}
								l997:
									{
										position999, tokenIndex999, depth999 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1000
										}
										position++
										goto l999
									l1000:
										position, tokenIndex, depth = position999, tokenIndex999, depth999
										if buffer[position] != rune('L') {
											goto l671
										}
										position++
									}
								l999:
									{
										position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1002
										}
										position++
										goto l1001
									l1002:
										position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
										if buffer[position] != rune('O') {
											goto l671
										}
										position++
									}
								l1001:
									{
										position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1004
										}
										position++
										goto l1003
									l1004:
										position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
										if buffer[position] != rune('O') {
											goto l671
										}
										position++
									}
								l1003:
									{
										position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1006
										}
										position++
										goto l1005
									l1006:
										position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
										if buffer[position] != rune('R') {
											goto l671
										}
										position++
									}
								l1005:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleFLOOR, position996)
								}
								break
							case 'R', 'r':
								{
									position1007 := position
									depth++
									{
										position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1009
										}
										position++
										goto l1008
									l1009:
										position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
										if buffer[position] != rune('R') {
											goto l671
										}
										position++
									}
								l1008:
									{
										position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1011
										}
										position++
										goto l1010
									l1011:
										position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
										if buffer[position] != rune('O') {
											goto l671
										}
										position++
									}
								l1010:
									{
										position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1013
										}
										position++
										goto l1012
									l1013:
										position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
										if buffer[position] != rune('U') {
											goto l671
										}
										position++
									}
								l1012:
									{
										position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1015
										}
										position++
										goto l1014
									l1015:
										position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
										if buffer[position] != rune('N') {
											goto l671
										}
										position++
									}
								l1014:
									{
										position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1017
										}
										position++
										goto l1016
									l1017:
										position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
										if buffer[position] != rune('D') {
											goto l671
										}
										position++
									}
								l1016:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleROUND, position1007)
								}
								break
							case 'C', 'c':
								{
									position1018 := position
									depth++
									{
										position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1020
										}
										position++
										goto l1019
									l1020:
										position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
										if buffer[position] != rune('C') {
											goto l671
										}
										position++
									}
								l1019:
									{
										position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1022
										}
										position++
										goto l1021
									l1022:
										position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
										if buffer[position] != rune('E') {
											goto l671
										}
										position++
									}
								l1021:
									{
										position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1024
										}
										position++
										goto l1023
									l1024:
										position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
										if buffer[position] != rune('I') {
											goto l671
										}
										position++
									}
								l1023:
									{
										position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1026
										}
										position++
										goto l1025
									l1026:
										position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
										if buffer[position] != rune('L') {
											goto l671
										}
										position++
									}
								l1025:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleCEIL, position1018)
								}
								break
							default:
								{
									position1027 := position
									depth++
									{
										position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1029
										}
										position++
										goto l1028
									l1029:
										position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
										if buffer[position] != rune('A') {
											goto l671
										}
										position++
									}
								l1028:
									{
										position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1031
										}
										position++
										goto l1030
									l1031:
										position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
										if buffer[position] != rune('B') {
											goto l671
										}
										position++
									}
								l1030:
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
											goto l671
										}
										position++
									}
								l1032:
									if !rules[ruleskip]() {
										goto l671
									}
									depth--
									add(ruleABS, position1027)
								}
								break
							}
						}

					}
				l672:
					if !rules[ruleLPAREN]() {
						goto l671
					}
					if !rules[ruleexpression]() {
						goto l671
					}
					if !rules[ruleRPAREN]() {
						goto l671
					}
					goto l670
				l671:
					position, tokenIndex, depth = position670, tokenIndex670, depth670
					{
						position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
						{
							position1037 := position
							depth++
							{
								position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1039
								}
								position++
								goto l1038
							l1039:
								position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
								if buffer[position] != rune('S') {
									goto l1036
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
									goto l1036
								}
								position++
							}
						l1040:
							{
								position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1043
								}
								position++
								goto l1042
							l1043:
								position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
								if buffer[position] != rune('R') {
									goto l1036
								}
								position++
							}
						l1042:
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
									goto l1036
								}
								position++
							}
						l1044:
							{
								position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1047
								}
								position++
								goto l1046
							l1047:
								position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
								if buffer[position] != rune('T') {
									goto l1036
								}
								position++
							}
						l1046:
							{
								position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1049
								}
								position++
								goto l1048
							l1049:
								position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
								if buffer[position] != rune('A') {
									goto l1036
								}
								position++
							}
						l1048:
							{
								position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1051
								}
								position++
								goto l1050
							l1051:
								position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
								if buffer[position] != rune('R') {
									goto l1036
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
									goto l1036
								}
								position++
							}
						l1052:
							{
								position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1055
								}
								position++
								goto l1054
							l1055:
								position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
								if buffer[position] != rune('S') {
									goto l1036
								}
								position++
							}
						l1054:
							if !rules[ruleskip]() {
								goto l1036
							}
							depth--
							add(ruleSTRSTARTS, position1037)
						}
						goto l1035
					l1036:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						{
							position1057 := position
							depth++
							{
								position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1059
								}
								position++
								goto l1058
							l1059:
								position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
								if buffer[position] != rune('S') {
									goto l1056
								}
								position++
							}
						l1058:
							{
								position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1061
								}
								position++
								goto l1060
							l1061:
								position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
								if buffer[position] != rune('T') {
									goto l1056
								}
								position++
							}
						l1060:
							{
								position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1063
								}
								position++
								goto l1062
							l1063:
								position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
								if buffer[position] != rune('R') {
									goto l1056
								}
								position++
							}
						l1062:
							{
								position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1065
								}
								position++
								goto l1064
							l1065:
								position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
								if buffer[position] != rune('E') {
									goto l1056
								}
								position++
							}
						l1064:
							{
								position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1067
								}
								position++
								goto l1066
							l1067:
								position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
								if buffer[position] != rune('N') {
									goto l1056
								}
								position++
							}
						l1066:
							{
								position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1069
								}
								position++
								goto l1068
							l1069:
								position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
								if buffer[position] != rune('D') {
									goto l1056
								}
								position++
							}
						l1068:
							{
								position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1071
								}
								position++
								goto l1070
							l1071:
								position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
								if buffer[position] != rune('S') {
									goto l1056
								}
								position++
							}
						l1070:
							if !rules[ruleskip]() {
								goto l1056
							}
							depth--
							add(ruleSTRENDS, position1057)
						}
						goto l1035
					l1056:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						{
							position1073 := position
							depth++
							{
								position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1075
								}
								position++
								goto l1074
							l1075:
								position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
								if buffer[position] != rune('S') {
									goto l1072
								}
								position++
							}
						l1074:
							{
								position1076, tokenIndex1076, depth1076 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1077
								}
								position++
								goto l1076
							l1077:
								position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
								if buffer[position] != rune('T') {
									goto l1072
								}
								position++
							}
						l1076:
							{
								position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1079
								}
								position++
								goto l1078
							l1079:
								position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
								if buffer[position] != rune('R') {
									goto l1072
								}
								position++
							}
						l1078:
							{
								position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1081
								}
								position++
								goto l1080
							l1081:
								position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
								if buffer[position] != rune('B') {
									goto l1072
								}
								position++
							}
						l1080:
							{
								position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1083
								}
								position++
								goto l1082
							l1083:
								position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
								if buffer[position] != rune('E') {
									goto l1072
								}
								position++
							}
						l1082:
							{
								position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1085
								}
								position++
								goto l1084
							l1085:
								position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
								if buffer[position] != rune('F') {
									goto l1072
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
									goto l1072
								}
								position++
							}
						l1086:
							{
								position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1089
								}
								position++
								goto l1088
							l1089:
								position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
								if buffer[position] != rune('R') {
									goto l1072
								}
								position++
							}
						l1088:
							{
								position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1091
								}
								position++
								goto l1090
							l1091:
								position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
								if buffer[position] != rune('E') {
									goto l1072
								}
								position++
							}
						l1090:
							if !rules[ruleskip]() {
								goto l1072
							}
							depth--
							add(ruleSTRBEFORE, position1073)
						}
						goto l1035
					l1072:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						{
							position1093 := position
							depth++
							{
								position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1095
								}
								position++
								goto l1094
							l1095:
								position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
								if buffer[position] != rune('S') {
									goto l1092
								}
								position++
							}
						l1094:
							{
								position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1097
								}
								position++
								goto l1096
							l1097:
								position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
								if buffer[position] != rune('T') {
									goto l1092
								}
								position++
							}
						l1096:
							{
								position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1099
								}
								position++
								goto l1098
							l1099:
								position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
								if buffer[position] != rune('R') {
									goto l1092
								}
								position++
							}
						l1098:
							{
								position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1101
								}
								position++
								goto l1100
							l1101:
								position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
								if buffer[position] != rune('A') {
									goto l1092
								}
								position++
							}
						l1100:
							{
								position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1103
								}
								position++
								goto l1102
							l1103:
								position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
								if buffer[position] != rune('F') {
									goto l1092
								}
								position++
							}
						l1102:
							{
								position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1105
								}
								position++
								goto l1104
							l1105:
								position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
								if buffer[position] != rune('T') {
									goto l1092
								}
								position++
							}
						l1104:
							{
								position1106, tokenIndex1106, depth1106 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1107
								}
								position++
								goto l1106
							l1107:
								position, tokenIndex, depth = position1106, tokenIndex1106, depth1106
								if buffer[position] != rune('E') {
									goto l1092
								}
								position++
							}
						l1106:
							{
								position1108, tokenIndex1108, depth1108 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1109
								}
								position++
								goto l1108
							l1109:
								position, tokenIndex, depth = position1108, tokenIndex1108, depth1108
								if buffer[position] != rune('R') {
									goto l1092
								}
								position++
							}
						l1108:
							if !rules[ruleskip]() {
								goto l1092
							}
							depth--
							add(ruleSTRAFTER, position1093)
						}
						goto l1035
					l1092:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						{
							position1111 := position
							depth++
							{
								position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1113
								}
								position++
								goto l1112
							l1113:
								position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
								if buffer[position] != rune('S') {
									goto l1110
								}
								position++
							}
						l1112:
							{
								position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1115
								}
								position++
								goto l1114
							l1115:
								position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
								if buffer[position] != rune('T') {
									goto l1110
								}
								position++
							}
						l1114:
							{
								position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1117
								}
								position++
								goto l1116
							l1117:
								position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
								if buffer[position] != rune('R') {
									goto l1110
								}
								position++
							}
						l1116:
							{
								position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1119
								}
								position++
								goto l1118
							l1119:
								position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
								if buffer[position] != rune('L') {
									goto l1110
								}
								position++
							}
						l1118:
							{
								position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1121
								}
								position++
								goto l1120
							l1121:
								position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
								if buffer[position] != rune('A') {
									goto l1110
								}
								position++
							}
						l1120:
							{
								position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1123
								}
								position++
								goto l1122
							l1123:
								position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
								if buffer[position] != rune('N') {
									goto l1110
								}
								position++
							}
						l1122:
							{
								position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1125
								}
								position++
								goto l1124
							l1125:
								position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
								if buffer[position] != rune('G') {
									goto l1110
								}
								position++
							}
						l1124:
							if !rules[ruleskip]() {
								goto l1110
							}
							depth--
							add(ruleSTRLANG, position1111)
						}
						goto l1035
					l1110:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						{
							position1127 := position
							depth++
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
									goto l1126
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
									goto l1126
								}
								position++
							}
						l1130:
							{
								position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1133
								}
								position++
								goto l1132
							l1133:
								position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
								if buffer[position] != rune('R') {
									goto l1126
								}
								position++
							}
						l1132:
							{
								position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1135
								}
								position++
								goto l1134
							l1135:
								position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
								if buffer[position] != rune('D') {
									goto l1126
								}
								position++
							}
						l1134:
							{
								position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1137
								}
								position++
								goto l1136
							l1137:
								position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
								if buffer[position] != rune('T') {
									goto l1126
								}
								position++
							}
						l1136:
							if !rules[ruleskip]() {
								goto l1126
							}
							depth--
							add(ruleSTRDT, position1127)
						}
						goto l1035
					l1126:
						position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1139 := position
									depth++
									{
										position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1141
										}
										position++
										goto l1140
									l1141:
										position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
										if buffer[position] != rune('S') {
											goto l1034
										}
										position++
									}
								l1140:
									{
										position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1143
										}
										position++
										goto l1142
									l1143:
										position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
										if buffer[position] != rune('A') {
											goto l1034
										}
										position++
									}
								l1142:
									{
										position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1145
										}
										position++
										goto l1144
									l1145:
										position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
										if buffer[position] != rune('M') {
											goto l1034
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
											goto l1034
										}
										position++
									}
								l1146:
									{
										position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1149
										}
										position++
										goto l1148
									l1149:
										position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
										if buffer[position] != rune('T') {
											goto l1034
										}
										position++
									}
								l1148:
									{
										position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1151
										}
										position++
										goto l1150
									l1151:
										position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
										if buffer[position] != rune('E') {
											goto l1034
										}
										position++
									}
								l1150:
									{
										position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1153
										}
										position++
										goto l1152
									l1153:
										position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
										if buffer[position] != rune('R') {
											goto l1034
										}
										position++
									}
								l1152:
									{
										position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1155
										}
										position++
										goto l1154
									l1155:
										position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
										if buffer[position] != rune('M') {
											goto l1034
										}
										position++
									}
								l1154:
									if !rules[ruleskip]() {
										goto l1034
									}
									depth--
									add(ruleSAMETERM, position1139)
								}
								break
							case 'C', 'c':
								{
									position1156 := position
									depth++
									{
										position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1158
										}
										position++
										goto l1157
									l1158:
										position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
										if buffer[position] != rune('C') {
											goto l1034
										}
										position++
									}
								l1157:
									{
										position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1160
										}
										position++
										goto l1159
									l1160:
										position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
										if buffer[position] != rune('O') {
											goto l1034
										}
										position++
									}
								l1159:
									{
										position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1162
										}
										position++
										goto l1161
									l1162:
										position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
										if buffer[position] != rune('N') {
											goto l1034
										}
										position++
									}
								l1161:
									{
										position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1164
										}
										position++
										goto l1163
									l1164:
										position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
										if buffer[position] != rune('T') {
											goto l1034
										}
										position++
									}
								l1163:
									{
										position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1166
										}
										position++
										goto l1165
									l1166:
										position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
										if buffer[position] != rune('A') {
											goto l1034
										}
										position++
									}
								l1165:
									{
										position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1168
										}
										position++
										goto l1167
									l1168:
										position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
										if buffer[position] != rune('I') {
											goto l1034
										}
										position++
									}
								l1167:
									{
										position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1170
										}
										position++
										goto l1169
									l1170:
										position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
										if buffer[position] != rune('N') {
											goto l1034
										}
										position++
									}
								l1169:
									{
										position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1172
										}
										position++
										goto l1171
									l1172:
										position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
										if buffer[position] != rune('S') {
											goto l1034
										}
										position++
									}
								l1171:
									if !rules[ruleskip]() {
										goto l1034
									}
									depth--
									add(ruleCONTAINS, position1156)
								}
								break
							default:
								{
									position1173 := position
									depth++
									{
										position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1175
										}
										position++
										goto l1174
									l1175:
										position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
										if buffer[position] != rune('L') {
											goto l1034
										}
										position++
									}
								l1174:
									{
										position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1177
										}
										position++
										goto l1176
									l1177:
										position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
										if buffer[position] != rune('A') {
											goto l1034
										}
										position++
									}
								l1176:
									{
										position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1179
										}
										position++
										goto l1178
									l1179:
										position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
										if buffer[position] != rune('N') {
											goto l1034
										}
										position++
									}
								l1178:
									{
										position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1181
										}
										position++
										goto l1180
									l1181:
										position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
										if buffer[position] != rune('G') {
											goto l1034
										}
										position++
									}
								l1180:
									{
										position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1183
										}
										position++
										goto l1182
									l1183:
										position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
										if buffer[position] != rune('M') {
											goto l1034
										}
										position++
									}
								l1182:
									{
										position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1185
										}
										position++
										goto l1184
									l1185:
										position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
										if buffer[position] != rune('A') {
											goto l1034
										}
										position++
									}
								l1184:
									{
										position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1187
										}
										position++
										goto l1186
									l1187:
										position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
										if buffer[position] != rune('T') {
											goto l1034
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
											goto l1034
										}
										position++
									}
								l1188:
									{
										position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1191
										}
										position++
										goto l1190
									l1191:
										position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
										if buffer[position] != rune('H') {
											goto l1034
										}
										position++
									}
								l1190:
									{
										position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1193
										}
										position++
										goto l1192
									l1193:
										position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
										if buffer[position] != rune('E') {
											goto l1034
										}
										position++
									}
								l1192:
									{
										position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1195
										}
										position++
										goto l1194
									l1195:
										position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
										if buffer[position] != rune('S') {
											goto l1034
										}
										position++
									}
								l1194:
									if !rules[ruleskip]() {
										goto l1034
									}
									depth--
									add(ruleLANGMATCHES, position1173)
								}
								break
							}
						}

					}
				l1035:
					if !rules[ruleLPAREN]() {
						goto l1034
					}
					if !rules[ruleexpression]() {
						goto l1034
					}
					if !rules[ruleCOMMA]() {
						goto l1034
					}
					if !rules[ruleexpression]() {
						goto l1034
					}
					if !rules[ruleRPAREN]() {
						goto l1034
					}
					goto l670
				l1034:
					position, tokenIndex, depth = position670, tokenIndex670, depth670
					{
						position1197 := position
						depth++
						{
							position1198, tokenIndex1198, depth1198 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1199
							}
							position++
							goto l1198
						l1199:
							position, tokenIndex, depth = position1198, tokenIndex1198, depth1198
							if buffer[position] != rune('B') {
								goto l1196
							}
							position++
						}
					l1198:
						{
							position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1201
							}
							position++
							goto l1200
						l1201:
							position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
							if buffer[position] != rune('O') {
								goto l1196
							}
							position++
						}
					l1200:
						{
							position1202, tokenIndex1202, depth1202 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1203
							}
							position++
							goto l1202
						l1203:
							position, tokenIndex, depth = position1202, tokenIndex1202, depth1202
							if buffer[position] != rune('U') {
								goto l1196
							}
							position++
						}
					l1202:
						{
							position1204, tokenIndex1204, depth1204 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1205
							}
							position++
							goto l1204
						l1205:
							position, tokenIndex, depth = position1204, tokenIndex1204, depth1204
							if buffer[position] != rune('N') {
								goto l1196
							}
							position++
						}
					l1204:
						{
							position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1207
							}
							position++
							goto l1206
						l1207:
							position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
							if buffer[position] != rune('D') {
								goto l1196
							}
							position++
						}
					l1206:
						if !rules[ruleskip]() {
							goto l1196
						}
						depth--
						add(ruleBOUND, position1197)
					}
					if !rules[ruleLPAREN]() {
						goto l1196
					}
					if !rules[rulevar]() {
						goto l1196
					}
					if !rules[ruleRPAREN]() {
						goto l1196
					}
					goto l670
				l1196:
					position, tokenIndex, depth = position670, tokenIndex670, depth670
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1210 := position
								depth++
								{
									position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1212
									}
									position++
									goto l1211
								l1212:
									position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
									if buffer[position] != rune('S') {
										goto l1208
									}
									position++
								}
							l1211:
								{
									position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1214
									}
									position++
									goto l1213
								l1214:
									position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
									if buffer[position] != rune('T') {
										goto l1208
									}
									position++
								}
							l1213:
								{
									position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1216
									}
									position++
									goto l1215
								l1216:
									position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
									if buffer[position] != rune('R') {
										goto l1208
									}
									position++
								}
							l1215:
								{
									position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1218
									}
									position++
									goto l1217
								l1218:
									position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
									if buffer[position] != rune('U') {
										goto l1208
									}
									position++
								}
							l1217:
								{
									position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1220
									}
									position++
									goto l1219
								l1220:
									position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
									if buffer[position] != rune('U') {
										goto l1208
									}
									position++
								}
							l1219:
								{
									position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1222
									}
									position++
									goto l1221
								l1222:
									position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
									if buffer[position] != rune('I') {
										goto l1208
									}
									position++
								}
							l1221:
								{
									position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1224
									}
									position++
									goto l1223
								l1224:
									position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
									if buffer[position] != rune('D') {
										goto l1208
									}
									position++
								}
							l1223:
								if !rules[ruleskip]() {
									goto l1208
								}
								depth--
								add(ruleSTRUUID, position1210)
							}
							break
						case 'U', 'u':
							{
								position1225 := position
								depth++
								{
									position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1227
									}
									position++
									goto l1226
								l1227:
									position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
									if buffer[position] != rune('U') {
										goto l1208
									}
									position++
								}
							l1226:
								{
									position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1229
									}
									position++
									goto l1228
								l1229:
									position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
									if buffer[position] != rune('U') {
										goto l1208
									}
									position++
								}
							l1228:
								{
									position1230, tokenIndex1230, depth1230 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1231
									}
									position++
									goto l1230
								l1231:
									position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
									if buffer[position] != rune('I') {
										goto l1208
									}
									position++
								}
							l1230:
								{
									position1232, tokenIndex1232, depth1232 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1233
									}
									position++
									goto l1232
								l1233:
									position, tokenIndex, depth = position1232, tokenIndex1232, depth1232
									if buffer[position] != rune('D') {
										goto l1208
									}
									position++
								}
							l1232:
								if !rules[ruleskip]() {
									goto l1208
								}
								depth--
								add(ruleUUID, position1225)
							}
							break
						case 'N', 'n':
							{
								position1234 := position
								depth++
								{
									position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1236
									}
									position++
									goto l1235
								l1236:
									position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
									if buffer[position] != rune('N') {
										goto l1208
									}
									position++
								}
							l1235:
								{
									position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1238
									}
									position++
									goto l1237
								l1238:
									position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
									if buffer[position] != rune('O') {
										goto l1208
									}
									position++
								}
							l1237:
								{
									position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1240
									}
									position++
									goto l1239
								l1240:
									position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
									if buffer[position] != rune('W') {
										goto l1208
									}
									position++
								}
							l1239:
								if !rules[ruleskip]() {
									goto l1208
								}
								depth--
								add(ruleNOW, position1234)
							}
							break
						default:
							{
								position1241 := position
								depth++
								{
									position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1243
									}
									position++
									goto l1242
								l1243:
									position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
									if buffer[position] != rune('R') {
										goto l1208
									}
									position++
								}
							l1242:
								{
									position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1245
									}
									position++
									goto l1244
								l1245:
									position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
									if buffer[position] != rune('A') {
										goto l1208
									}
									position++
								}
							l1244:
								{
									position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1247
									}
									position++
									goto l1246
								l1247:
									position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
									if buffer[position] != rune('N') {
										goto l1208
									}
									position++
								}
							l1246:
								{
									position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1249
									}
									position++
									goto l1248
								l1249:
									position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
									if buffer[position] != rune('D') {
										goto l1208
									}
									position++
								}
							l1248:
								if !rules[ruleskip]() {
									goto l1208
								}
								depth--
								add(ruleRAND, position1241)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1208
					}
					goto l670
				l1208:
					position, tokenIndex, depth = position670, tokenIndex670, depth670
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
								{
									position1253 := position
									depth++
									{
										position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1255
										}
										position++
										goto l1254
									l1255:
										position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
										if buffer[position] != rune('E') {
											goto l1252
										}
										position++
									}
								l1254:
									{
										position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1257
										}
										position++
										goto l1256
									l1257:
										position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
										if buffer[position] != rune('X') {
											goto l1252
										}
										position++
									}
								l1256:
									{
										position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1259
										}
										position++
										goto l1258
									l1259:
										position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
										if buffer[position] != rune('I') {
											goto l1252
										}
										position++
									}
								l1258:
									{
										position1260, tokenIndex1260, depth1260 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1261
										}
										position++
										goto l1260
									l1261:
										position, tokenIndex, depth = position1260, tokenIndex1260, depth1260
										if buffer[position] != rune('S') {
											goto l1252
										}
										position++
									}
								l1260:
									{
										position1262, tokenIndex1262, depth1262 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1263
										}
										position++
										goto l1262
									l1263:
										position, tokenIndex, depth = position1262, tokenIndex1262, depth1262
										if buffer[position] != rune('T') {
											goto l1252
										}
										position++
									}
								l1262:
									{
										position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1265
										}
										position++
										goto l1264
									l1265:
										position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
										if buffer[position] != rune('S') {
											goto l1252
										}
										position++
									}
								l1264:
									if !rules[ruleskip]() {
										goto l1252
									}
									depth--
									add(ruleEXISTS, position1253)
								}
								goto l1251
							l1252:
								position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
								{
									position1266 := position
									depth++
									{
										position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1268
										}
										position++
										goto l1267
									l1268:
										position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
										if buffer[position] != rune('N') {
											goto l668
										}
										position++
									}
								l1267:
									{
										position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1270
										}
										position++
										goto l1269
									l1270:
										position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
										if buffer[position] != rune('O') {
											goto l668
										}
										position++
									}
								l1269:
									{
										position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1272
										}
										position++
										goto l1271
									l1272:
										position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
										if buffer[position] != rune('T') {
											goto l668
										}
										position++
									}
								l1271:
									if buffer[position] != rune(' ') {
										goto l668
									}
									position++
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
											goto l668
										}
										position++
									}
								l1273:
									{
										position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1276
										}
										position++
										goto l1275
									l1276:
										position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
										if buffer[position] != rune('X') {
											goto l668
										}
										position++
									}
								l1275:
									{
										position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1278
										}
										position++
										goto l1277
									l1278:
										position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
										if buffer[position] != rune('I') {
											goto l668
										}
										position++
									}
								l1277:
									{
										position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1280
										}
										position++
										goto l1279
									l1280:
										position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
										if buffer[position] != rune('S') {
											goto l668
										}
										position++
									}
								l1279:
									{
										position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1282
										}
										position++
										goto l1281
									l1282:
										position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
										if buffer[position] != rune('T') {
											goto l668
										}
										position++
									}
								l1281:
									{
										position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1284
										}
										position++
										goto l1283
									l1284:
										position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
										if buffer[position] != rune('S') {
											goto l668
										}
										position++
									}
								l1283:
									if !rules[ruleskip]() {
										goto l668
									}
									depth--
									add(ruleNOTEXIST, position1266)
								}
							}
						l1251:
							if !rules[rulegroupGraphPattern]() {
								goto l668
							}
							break
						case 'I', 'i':
							{
								position1285 := position
								depth++
								{
									position1286, tokenIndex1286, depth1286 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1287
									}
									position++
									goto l1286
								l1287:
									position, tokenIndex, depth = position1286, tokenIndex1286, depth1286
									if buffer[position] != rune('I') {
										goto l668
									}
									position++
								}
							l1286:
								{
									position1288, tokenIndex1288, depth1288 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1289
									}
									position++
									goto l1288
								l1289:
									position, tokenIndex, depth = position1288, tokenIndex1288, depth1288
									if buffer[position] != rune('F') {
										goto l668
									}
									position++
								}
							l1288:
								if !rules[ruleskip]() {
									goto l668
								}
								depth--
								add(ruleIF, position1285)
							}
							if !rules[ruleLPAREN]() {
								goto l668
							}
							if !rules[ruleexpression]() {
								goto l668
							}
							if !rules[ruleCOMMA]() {
								goto l668
							}
							if !rules[ruleexpression]() {
								goto l668
							}
							if !rules[ruleCOMMA]() {
								goto l668
							}
							if !rules[ruleexpression]() {
								goto l668
							}
							if !rules[ruleRPAREN]() {
								goto l668
							}
							break
						case 'C', 'c':
							{
								position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
								{
									position1292 := position
									depth++
									{
										position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1294
										}
										position++
										goto l1293
									l1294:
										position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
										if buffer[position] != rune('C') {
											goto l1291
										}
										position++
									}
								l1293:
									{
										position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1296
										}
										position++
										goto l1295
									l1296:
										position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
										if buffer[position] != rune('O') {
											goto l1291
										}
										position++
									}
								l1295:
									{
										position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1298
										}
										position++
										goto l1297
									l1298:
										position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
										if buffer[position] != rune('N') {
											goto l1291
										}
										position++
									}
								l1297:
									{
										position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1300
										}
										position++
										goto l1299
									l1300:
										position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
										if buffer[position] != rune('C') {
											goto l1291
										}
										position++
									}
								l1299:
									{
										position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1302
										}
										position++
										goto l1301
									l1302:
										position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
										if buffer[position] != rune('A') {
											goto l1291
										}
										position++
									}
								l1301:
									{
										position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1304
										}
										position++
										goto l1303
									l1304:
										position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
										if buffer[position] != rune('T') {
											goto l1291
										}
										position++
									}
								l1303:
									if !rules[ruleskip]() {
										goto l1291
									}
									depth--
									add(ruleCONCAT, position1292)
								}
								goto l1290
							l1291:
								position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
								{
									position1305 := position
									depth++
									{
										position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1307
										}
										position++
										goto l1306
									l1307:
										position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
										if buffer[position] != rune('C') {
											goto l668
										}
										position++
									}
								l1306:
									{
										position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1309
										}
										position++
										goto l1308
									l1309:
										position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
										if buffer[position] != rune('O') {
											goto l668
										}
										position++
									}
								l1308:
									{
										position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1311
										}
										position++
										goto l1310
									l1311:
										position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
										if buffer[position] != rune('A') {
											goto l668
										}
										position++
									}
								l1310:
									{
										position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1313
										}
										position++
										goto l1312
									l1313:
										position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
										if buffer[position] != rune('L') {
											goto l668
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
											goto l668
										}
										position++
									}
								l1314:
									{
										position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1317
										}
										position++
										goto l1316
									l1317:
										position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
										if buffer[position] != rune('S') {
											goto l668
										}
										position++
									}
								l1316:
									{
										position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1319
										}
										position++
										goto l1318
									l1319:
										position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
										if buffer[position] != rune('C') {
											goto l668
										}
										position++
									}
								l1318:
									{
										position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1321
										}
										position++
										goto l1320
									l1321:
										position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
										if buffer[position] != rune('E') {
											goto l668
										}
										position++
									}
								l1320:
									if !rules[ruleskip]() {
										goto l668
									}
									depth--
									add(ruleCOALESCE, position1305)
								}
							}
						l1290:
							if !rules[ruleargList]() {
								goto l668
							}
							break
						case 'B', 'b':
							{
								position1322 := position
								depth++
								{
									position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1324
									}
									position++
									goto l1323
								l1324:
									position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
									if buffer[position] != rune('B') {
										goto l668
									}
									position++
								}
							l1323:
								{
									position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1326
									}
									position++
									goto l1325
								l1326:
									position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
									if buffer[position] != rune('N') {
										goto l668
									}
									position++
								}
							l1325:
								{
									position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1328
									}
									position++
									goto l1327
								l1328:
									position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
									if buffer[position] != rune('O') {
										goto l668
									}
									position++
								}
							l1327:
								{
									position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1330
									}
									position++
									goto l1329
								l1330:
									position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
									if buffer[position] != rune('D') {
										goto l668
									}
									position++
								}
							l1329:
								{
									position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1332
									}
									position++
									goto l1331
								l1332:
									position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
									if buffer[position] != rune('E') {
										goto l668
									}
									position++
								}
							l1331:
								if !rules[ruleskip]() {
									goto l668
								}
								depth--
								add(ruleBNODE, position1322)
							}
							{
								position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1334
								}
								if !rules[ruleexpression]() {
									goto l1334
								}
								if !rules[ruleRPAREN]() {
									goto l1334
								}
								goto l1333
							l1334:
								position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
								if !rules[rulenil]() {
									goto l668
								}
							}
						l1333:
							break
						default:
							{
								position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
								{
									position1337 := position
									depth++
									{
										position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1339
										}
										position++
										goto l1338
									l1339:
										position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
										if buffer[position] != rune('S') {
											goto l1336
										}
										position++
									}
								l1338:
									{
										position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1341
										}
										position++
										goto l1340
									l1341:
										position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
										if buffer[position] != rune('U') {
											goto l1336
										}
										position++
									}
								l1340:
									{
										position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1343
										}
										position++
										goto l1342
									l1343:
										position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
										if buffer[position] != rune('B') {
											goto l1336
										}
										position++
									}
								l1342:
									{
										position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1345
										}
										position++
										goto l1344
									l1345:
										position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
										if buffer[position] != rune('S') {
											goto l1336
										}
										position++
									}
								l1344:
									{
										position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1347
										}
										position++
										goto l1346
									l1347:
										position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
										if buffer[position] != rune('T') {
											goto l1336
										}
										position++
									}
								l1346:
									{
										position1348, tokenIndex1348, depth1348 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1349
										}
										position++
										goto l1348
									l1349:
										position, tokenIndex, depth = position1348, tokenIndex1348, depth1348
										if buffer[position] != rune('R') {
											goto l1336
										}
										position++
									}
								l1348:
									if !rules[ruleskip]() {
										goto l1336
									}
									depth--
									add(ruleSUBSTR, position1337)
								}
								goto l1335
							l1336:
								position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
								{
									position1351 := position
									depth++
									{
										position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1353
										}
										position++
										goto l1352
									l1353:
										position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
										if buffer[position] != rune('R') {
											goto l1350
										}
										position++
									}
								l1352:
									{
										position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1355
										}
										position++
										goto l1354
									l1355:
										position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
										if buffer[position] != rune('E') {
											goto l1350
										}
										position++
									}
								l1354:
									{
										position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1357
										}
										position++
										goto l1356
									l1357:
										position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
										if buffer[position] != rune('P') {
											goto l1350
										}
										position++
									}
								l1356:
									{
										position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1359
										}
										position++
										goto l1358
									l1359:
										position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
										if buffer[position] != rune('L') {
											goto l1350
										}
										position++
									}
								l1358:
									{
										position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1361
										}
										position++
										goto l1360
									l1361:
										position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
										if buffer[position] != rune('A') {
											goto l1350
										}
										position++
									}
								l1360:
									{
										position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1363
										}
										position++
										goto l1362
									l1363:
										position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
										if buffer[position] != rune('C') {
											goto l1350
										}
										position++
									}
								l1362:
									{
										position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1365
										}
										position++
										goto l1364
									l1365:
										position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
										if buffer[position] != rune('E') {
											goto l1350
										}
										position++
									}
								l1364:
									if !rules[ruleskip]() {
										goto l1350
									}
									depth--
									add(ruleREPLACE, position1351)
								}
								goto l1335
							l1350:
								position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
								{
									position1366 := position
									depth++
									{
										position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1368
										}
										position++
										goto l1367
									l1368:
										position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
										if buffer[position] != rune('R') {
											goto l668
										}
										position++
									}
								l1367:
									{
										position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1370
										}
										position++
										goto l1369
									l1370:
										position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
										if buffer[position] != rune('E') {
											goto l668
										}
										position++
									}
								l1369:
									{
										position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1372
										}
										position++
										goto l1371
									l1372:
										position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
										if buffer[position] != rune('G') {
											goto l668
										}
										position++
									}
								l1371:
									{
										position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1374
										}
										position++
										goto l1373
									l1374:
										position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
										if buffer[position] != rune('E') {
											goto l668
										}
										position++
									}
								l1373:
									{
										position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1376
										}
										position++
										goto l1375
									l1376:
										position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
										if buffer[position] != rune('X') {
											goto l668
										}
										position++
									}
								l1375:
									if !rules[ruleskip]() {
										goto l668
									}
									depth--
									add(ruleREGEX, position1366)
								}
							}
						l1335:
							if !rules[ruleLPAREN]() {
								goto l668
							}
							if !rules[ruleexpression]() {
								goto l668
							}
							if !rules[ruleCOMMA]() {
								goto l668
							}
							if !rules[ruleexpression]() {
								goto l668
							}
							{
								position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1377
								}
								if !rules[ruleexpression]() {
									goto l1377
								}
								goto l1378
							l1377:
								position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
							}
						l1378:
							if !rules[ruleRPAREN]() {
								goto l668
							}
							break
						}
					}

				}
			l670:
				depth--
				add(rulebuiltinCall, position669)
			}
			return true
		l668:
			position, tokenIndex, depth = position668, tokenIndex668, depth668
			return false
		},
		/* 64 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
			{
				position1380 := position
				depth++
				{
					position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
					{
						position1383 := position
						depth++
					l1384:
						{
							position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
							{
								position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1387
								}
								position++
								goto l1386
							l1387:
								position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1385
								}
								position++
							}
						l1386:
							goto l1384
						l1385:
							position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
						}
						depth--
						add(rulePegText, position1383)
					}
					if buffer[position] != rune(':') {
						goto l1382
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1381
				l1382:
					position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
					{
						position1390 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1389
						}
						position++
					l1391:
						{
							position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1392
							}
							position++
							goto l1391
						l1392:
							position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
						}
						depth--
						add(rulePegText, position1390)
					}
					if buffer[position] != rune('/') {
						goto l1389
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1381
				l1389:
					position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
					{
						position1394 := position
						depth++
					l1395:
						{
							position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1396
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1396
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1396
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1396
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1396
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1396
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1396
									}
									position++
									break
								}
							}

							goto l1395
						l1396:
							position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
						}
						depth--
						add(rulePegText, position1394)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1381:
				if buffer[position] != rune('<') {
					goto l1379
				}
				position++
				if !rules[rulews]() {
					goto l1379
				}
				if !rules[ruleskip]() {
					goto l1379
				}
				depth--
				add(rulepof, position1380)
			}
			return true
		l1379:
			position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
			return false
		},
		/* 65 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
			{
				position1400 := position
				depth++
				{
					position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1402
					}
					position++
					goto l1401
				l1402:
					position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
					if buffer[position] != rune('$') {
						goto l1399
					}
					position++
				}
			l1401:
				{
					position1403 := position
					depth++
					{
						position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
						{
							position1408 := position
							depth++
							{
								position1409, tokenIndex1409, depth1409 := position, tokenIndex, depth
								{
									position1411 := position
									depth++
									{
										position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1413
										}
										position++
										goto l1412
									l1413:
										position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1410
										}
										position++
									}
								l1412:
									depth--
									add(rulePN_CHARS_BASE, position1411)
								}
								goto l1409
							l1410:
								position, tokenIndex, depth = position1409, tokenIndex1409, depth1409
								if buffer[position] != rune('_') {
									goto l1407
								}
								position++
							}
						l1409:
							depth--
							add(rulePN_CHARS_U, position1408)
						}
						goto l1406
					l1407:
						position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1399
						}
						position++
					}
				l1406:
				l1404:
					{
						position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
						{
							position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
							{
								position1416 := position
								depth++
								{
									position1417, tokenIndex1417, depth1417 := position, tokenIndex, depth
									{
										position1419 := position
										depth++
										{
											position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1421
											}
											position++
											goto l1420
										l1421:
											position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1418
											}
											position++
										}
									l1420:
										depth--
										add(rulePN_CHARS_BASE, position1419)
									}
									goto l1417
								l1418:
									position, tokenIndex, depth = position1417, tokenIndex1417, depth1417
									if buffer[position] != rune('_') {
										goto l1415
									}
									position++
								}
							l1417:
								depth--
								add(rulePN_CHARS_U, position1416)
							}
							goto l1414
						l1415:
							position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1405
							}
							position++
						}
					l1414:
						goto l1404
					l1405:
						position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
					}
					depth--
					add(ruleVARNAME, position1403)
				}
				if !rules[ruleskip]() {
					goto l1399
				}
				depth--
				add(rulevar, position1400)
			}
			return true
		l1399:
			position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
			return false
		},
		/* 66 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1422, tokenIndex1422, depth1422 := position, tokenIndex, depth
			{
				position1423 := position
				depth++
				{
					position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1425
					}
					goto l1424
				l1425:
					position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
					{
						position1426 := position
						depth++
					l1427:
						{
							position1428, tokenIndex1428, depth1428 := position, tokenIndex, depth
							{
								position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
								{
									position1430, tokenIndex1430, depth1430 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1431
									}
									position++
									goto l1430
								l1431:
									position, tokenIndex, depth = position1430, tokenIndex1430, depth1430
									if buffer[position] != rune(' ') {
										goto l1429
									}
									position++
								}
							l1430:
								goto l1428
							l1429:
								position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
							}
							if !matchDot() {
								goto l1428
							}
							goto l1427
						l1428:
							position, tokenIndex, depth = position1428, tokenIndex1428, depth1428
						}
						if buffer[position] != rune(':') {
							goto l1422
						}
						position++
					l1432:
						{
							position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
							{
								position1434, tokenIndex1434, depth1434 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1435
								}
								position++
								goto l1434
							l1435:
								position, tokenIndex, depth = position1434, tokenIndex1434, depth1434
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1436
								}
								position++
								goto l1434
							l1436:
								position, tokenIndex, depth = position1434, tokenIndex1434, depth1434
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1437
								}
								position++
								goto l1434
							l1437:
								position, tokenIndex, depth = position1434, tokenIndex1434, depth1434
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1433
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1433
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1433
										}
										position++
										break
									}
								}

							}
						l1434:
							goto l1432
						l1433:
							position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
						}
						if !rules[ruleskip]() {
							goto l1422
						}
						depth--
						add(ruleprefixedName, position1426)
					}
				}
			l1424:
				depth--
				add(ruleiriref, position1423)
			}
			return true
		l1422:
			position, tokenIndex, depth = position1422, tokenIndex1422, depth1422
			return false
		},
		/* 67 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
			{
				position1440 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1439
				}
				position++
			l1441:
				{
					position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
					{
						position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1443
						}
						position++
						goto l1442
					l1443:
						position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
					}
					if !matchDot() {
						goto l1442
					}
					goto l1441
				l1442:
					position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
				}
				if buffer[position] != rune('>') {
					goto l1439
				}
				position++
				if !rules[ruleskip]() {
					goto l1439
				}
				depth--
				add(ruleiri, position1440)
			}
			return true
		l1439:
			position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
			return false
		},
		/* 68 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 69 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1445, tokenIndex1445, depth1445 := position, tokenIndex, depth
			{
				position1446 := position
				depth++
				if !rules[rulestring]() {
					goto l1445
				}
				{
					position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
					{
						position1449, tokenIndex1449, depth1449 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1450
						}
						position++
						{
							position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1454
							}
							position++
							goto l1453
						l1454:
							position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1450
							}
							position++
						}
					l1453:
					l1451:
						{
							position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
							{
								position1455, tokenIndex1455, depth1455 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1456
								}
								position++
								goto l1455
							l1456:
								position, tokenIndex, depth = position1455, tokenIndex1455, depth1455
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1452
								}
								position++
							}
						l1455:
							goto l1451
						l1452:
							position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
						}
					l1457:
						{
							position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1458
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1458
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1458
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1458
									}
									position++
									break
								}
							}

						l1459:
							{
								position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1460
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1460
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1460
										}
										position++
										break
									}
								}

								goto l1459
							l1460:
								position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
							}
							goto l1457
						l1458:
							position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
						}
						goto l1449
					l1450:
						position, tokenIndex, depth = position1449, tokenIndex1449, depth1449
						if buffer[position] != rune('^') {
							goto l1447
						}
						position++
						if buffer[position] != rune('^') {
							goto l1447
						}
						position++
						if !rules[ruleiriref]() {
							goto l1447
						}
					}
				l1449:
					goto l1448
				l1447:
					position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
				}
			l1448:
				if !rules[ruleskip]() {
					goto l1445
				}
				depth--
				add(ruleliteral, position1446)
			}
			return true
		l1445:
			position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
			return false
		},
		/* 70 string <- <('"' (!'"' .)* '"')> */
		func() bool {
			position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
			{
				position1464 := position
				depth++
				if buffer[position] != rune('"') {
					goto l1463
				}
				position++
			l1465:
				{
					position1466, tokenIndex1466, depth1466 := position, tokenIndex, depth
					{
						position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
						if buffer[position] != rune('"') {
							goto l1467
						}
						position++
						goto l1466
					l1467:
						position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
					}
					if !matchDot() {
						goto l1466
					}
					goto l1465
				l1466:
					position, tokenIndex, depth = position1466, tokenIndex1466, depth1466
				}
				if buffer[position] != rune('"') {
					goto l1463
				}
				position++
				depth--
				add(rulestring, position1464)
			}
			return true
		l1463:
			position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
			return false
		},
		/* 71 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
			{
				position1469 := position
				depth++
				{
					position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
					{
						position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1473
						}
						position++
						goto l1472
					l1473:
						position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
						if buffer[position] != rune('-') {
							goto l1470
						}
						position++
					}
				l1472:
					goto l1471
				l1470:
					position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
				}
			l1471:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1468
				}
				position++
			l1474:
				{
					position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1475
					}
					position++
					goto l1474
				l1475:
					position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
				}
				{
					position1476, tokenIndex1476, depth1476 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1476
					}
					position++
				l1478:
					{
						position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1479
						}
						position++
						goto l1478
					l1479:
						position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
					}
					goto l1477
				l1476:
					position, tokenIndex, depth = position1476, tokenIndex1476, depth1476
				}
			l1477:
				if !rules[ruleskip]() {
					goto l1468
				}
				depth--
				add(rulenumericLiteral, position1469)
			}
			return true
		l1468:
			position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
			return false
		},
		/* 72 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 73 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
			{
				position1482 := position
				depth++
				{
					position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
					{
						position1485 := position
						depth++
						{
							position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1487
							}
							position++
							goto l1486
						l1487:
							position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
							if buffer[position] != rune('T') {
								goto l1484
							}
							position++
						}
					l1486:
						{
							position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1489
							}
							position++
							goto l1488
						l1489:
							position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
							if buffer[position] != rune('R') {
								goto l1484
							}
							position++
						}
					l1488:
						{
							position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1491
							}
							position++
							goto l1490
						l1491:
							position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
							if buffer[position] != rune('U') {
								goto l1484
							}
							position++
						}
					l1490:
						{
							position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1493
							}
							position++
							goto l1492
						l1493:
							position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
							if buffer[position] != rune('E') {
								goto l1484
							}
							position++
						}
					l1492:
						if !rules[ruleskip]() {
							goto l1484
						}
						depth--
						add(ruleTRUE, position1485)
					}
					goto l1483
				l1484:
					position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
					{
						position1494 := position
						depth++
						{
							position1495, tokenIndex1495, depth1495 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1496
							}
							position++
							goto l1495
						l1496:
							position, tokenIndex, depth = position1495, tokenIndex1495, depth1495
							if buffer[position] != rune('F') {
								goto l1481
							}
							position++
						}
					l1495:
						{
							position1497, tokenIndex1497, depth1497 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1498
							}
							position++
							goto l1497
						l1498:
							position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
							if buffer[position] != rune('A') {
								goto l1481
							}
							position++
						}
					l1497:
						{
							position1499, tokenIndex1499, depth1499 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1500
							}
							position++
							goto l1499
						l1500:
							position, tokenIndex, depth = position1499, tokenIndex1499, depth1499
							if buffer[position] != rune('L') {
								goto l1481
							}
							position++
						}
					l1499:
						{
							position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1502
							}
							position++
							goto l1501
						l1502:
							position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
							if buffer[position] != rune('S') {
								goto l1481
							}
							position++
						}
					l1501:
						{
							position1503, tokenIndex1503, depth1503 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1504
							}
							position++
							goto l1503
						l1504:
							position, tokenIndex, depth = position1503, tokenIndex1503, depth1503
							if buffer[position] != rune('E') {
								goto l1481
							}
							position++
						}
					l1503:
						if !rules[ruleskip]() {
							goto l1481
						}
						depth--
						add(ruleFALSE, position1494)
					}
				}
			l1483:
				depth--
				add(rulebooleanLiteral, position1482)
			}
			return true
		l1481:
			position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
			return false
		},
		/* 74 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 75 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 76 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 77 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1508, tokenIndex1508, depth1508 := position, tokenIndex, depth
			{
				position1509 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1508
				}
				position++
			l1510:
				{
					position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1511
					}
					goto l1510
				l1511:
					position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
				}
				if buffer[position] != rune(')') {
					goto l1508
				}
				position++
				if !rules[ruleskip]() {
					goto l1508
				}
				depth--
				add(rulenil, position1509)
			}
			return true
		l1508:
			position, tokenIndex, depth = position1508, tokenIndex1508, depth1508
			return false
		},
		/* 78 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 79 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 80 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 81 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 82 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 83 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 84 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 85 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 86 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 87 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1521, tokenIndex1521, depth1521 := position, tokenIndex, depth
			{
				position1522 := position
				depth++
				{
					position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1524
					}
					position++
					goto l1523
				l1524:
					position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
					if buffer[position] != rune('D') {
						goto l1521
					}
					position++
				}
			l1523:
				{
					position1525, tokenIndex1525, depth1525 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1526
					}
					position++
					goto l1525
				l1526:
					position, tokenIndex, depth = position1525, tokenIndex1525, depth1525
					if buffer[position] != rune('I') {
						goto l1521
					}
					position++
				}
			l1525:
				{
					position1527, tokenIndex1527, depth1527 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1528
					}
					position++
					goto l1527
				l1528:
					position, tokenIndex, depth = position1527, tokenIndex1527, depth1527
					if buffer[position] != rune('S') {
						goto l1521
					}
					position++
				}
			l1527:
				{
					position1529, tokenIndex1529, depth1529 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1530
					}
					position++
					goto l1529
				l1530:
					position, tokenIndex, depth = position1529, tokenIndex1529, depth1529
					if buffer[position] != rune('T') {
						goto l1521
					}
					position++
				}
			l1529:
				{
					position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1532
					}
					position++
					goto l1531
				l1532:
					position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
					if buffer[position] != rune('I') {
						goto l1521
					}
					position++
				}
			l1531:
				{
					position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1534
					}
					position++
					goto l1533
				l1534:
					position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
					if buffer[position] != rune('N') {
						goto l1521
					}
					position++
				}
			l1533:
				{
					position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1536
					}
					position++
					goto l1535
				l1536:
					position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
					if buffer[position] != rune('C') {
						goto l1521
					}
					position++
				}
			l1535:
				{
					position1537, tokenIndex1537, depth1537 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1538
					}
					position++
					goto l1537
				l1538:
					position, tokenIndex, depth = position1537, tokenIndex1537, depth1537
					if buffer[position] != rune('T') {
						goto l1521
					}
					position++
				}
			l1537:
				if !rules[ruleskip]() {
					goto l1521
				}
				depth--
				add(ruleDISTINCT, position1522)
			}
			return true
		l1521:
			position, tokenIndex, depth = position1521, tokenIndex1521, depth1521
			return false
		},
		/* 88 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 89 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 90 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 91 LBRACE <- <('{' skip)> */
		func() bool {
			position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
			{
				position1543 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1542
				}
				position++
				if !rules[ruleskip]() {
					goto l1542
				}
				depth--
				add(ruleLBRACE, position1543)
			}
			return true
		l1542:
			position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
			return false
		},
		/* 92 RBRACE <- <('}' skip)> */
		func() bool {
			position1544, tokenIndex1544, depth1544 := position, tokenIndex, depth
			{
				position1545 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1544
				}
				position++
				if !rules[ruleskip]() {
					goto l1544
				}
				depth--
				add(ruleRBRACE, position1545)
			}
			return true
		l1544:
			position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
			return false
		},
		/* 93 LBRACK <- <('[' skip)> */
		nil,
		/* 94 RBRACK <- <(']' skip)> */
		nil,
		/* 95 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1548, tokenIndex1548, depth1548 := position, tokenIndex, depth
			{
				position1549 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1548
				}
				position++
				if !rules[ruleskip]() {
					goto l1548
				}
				depth--
				add(ruleSEMICOLON, position1549)
			}
			return true
		l1548:
			position, tokenIndex, depth = position1548, tokenIndex1548, depth1548
			return false
		},
		/* 96 COMMA <- <(',' skip)> */
		func() bool {
			position1550, tokenIndex1550, depth1550 := position, tokenIndex, depth
			{
				position1551 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1550
				}
				position++
				if !rules[ruleskip]() {
					goto l1550
				}
				depth--
				add(ruleCOMMA, position1551)
			}
			return true
		l1550:
			position, tokenIndex, depth = position1550, tokenIndex1550, depth1550
			return false
		},
		/* 97 DOT <- <('.' skip)> */
		func() bool {
			position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
			{
				position1553 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1552
				}
				position++
				if !rules[ruleskip]() {
					goto l1552
				}
				depth--
				add(ruleDOT, position1553)
			}
			return true
		l1552:
			position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
			return false
		},
		/* 98 COLON <- <(':' skip)> */
		nil,
		/* 99 PIPE <- <('|' skip)> */
		func() bool {
			position1555, tokenIndex1555, depth1555 := position, tokenIndex, depth
			{
				position1556 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1555
				}
				position++
				if !rules[ruleskip]() {
					goto l1555
				}
				depth--
				add(rulePIPE, position1556)
			}
			return true
		l1555:
			position, tokenIndex, depth = position1555, tokenIndex1555, depth1555
			return false
		},
		/* 100 SLASH <- <('/' skip)> */
		func() bool {
			position1557, tokenIndex1557, depth1557 := position, tokenIndex, depth
			{
				position1558 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1557
				}
				position++
				if !rules[ruleskip]() {
					goto l1557
				}
				depth--
				add(ruleSLASH, position1558)
			}
			return true
		l1557:
			position, tokenIndex, depth = position1557, tokenIndex1557, depth1557
			return false
		},
		/* 101 INVERSE <- <('^' skip)> */
		func() bool {
			position1559, tokenIndex1559, depth1559 := position, tokenIndex, depth
			{
				position1560 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1559
				}
				position++
				if !rules[ruleskip]() {
					goto l1559
				}
				depth--
				add(ruleINVERSE, position1560)
			}
			return true
		l1559:
			position, tokenIndex, depth = position1559, tokenIndex1559, depth1559
			return false
		},
		/* 102 LPAREN <- <('(' skip)> */
		func() bool {
			position1561, tokenIndex1561, depth1561 := position, tokenIndex, depth
			{
				position1562 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1561
				}
				position++
				if !rules[ruleskip]() {
					goto l1561
				}
				depth--
				add(ruleLPAREN, position1562)
			}
			return true
		l1561:
			position, tokenIndex, depth = position1561, tokenIndex1561, depth1561
			return false
		},
		/* 103 RPAREN <- <(')' skip)> */
		func() bool {
			position1563, tokenIndex1563, depth1563 := position, tokenIndex, depth
			{
				position1564 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1563
				}
				position++
				if !rules[ruleskip]() {
					goto l1563
				}
				depth--
				add(ruleRPAREN, position1564)
			}
			return true
		l1563:
			position, tokenIndex, depth = position1563, tokenIndex1563, depth1563
			return false
		},
		/* 104 ISA <- <('a' skip)> */
		func() bool {
			position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
			{
				position1566 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1565
				}
				position++
				if !rules[ruleskip]() {
					goto l1565
				}
				depth--
				add(ruleISA, position1566)
			}
			return true
		l1565:
			position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
			return false
		},
		/* 105 NOT <- <('!' skip)> */
		func() bool {
			position1567, tokenIndex1567, depth1567 := position, tokenIndex, depth
			{
				position1568 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1567
				}
				position++
				if !rules[ruleskip]() {
					goto l1567
				}
				depth--
				add(ruleNOT, position1568)
			}
			return true
		l1567:
			position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
			return false
		},
		/* 106 STAR <- <('*' skip)> */
		func() bool {
			position1569, tokenIndex1569, depth1569 := position, tokenIndex, depth
			{
				position1570 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1569
				}
				position++
				if !rules[ruleskip]() {
					goto l1569
				}
				depth--
				add(ruleSTAR, position1570)
			}
			return true
		l1569:
			position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
			return false
		},
		/* 107 PLUS <- <('+' skip)> */
		func() bool {
			position1571, tokenIndex1571, depth1571 := position, tokenIndex, depth
			{
				position1572 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1571
				}
				position++
				if !rules[ruleskip]() {
					goto l1571
				}
				depth--
				add(rulePLUS, position1572)
			}
			return true
		l1571:
			position, tokenIndex, depth = position1571, tokenIndex1571, depth1571
			return false
		},
		/* 108 MINUS <- <('-' skip)> */
		func() bool {
			position1573, tokenIndex1573, depth1573 := position, tokenIndex, depth
			{
				position1574 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1573
				}
				position++
				if !rules[ruleskip]() {
					goto l1573
				}
				depth--
				add(ruleMINUS, position1574)
			}
			return true
		l1573:
			position, tokenIndex, depth = position1573, tokenIndex1573, depth1573
			return false
		},
		/* 109 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 110 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 111 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 112 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 113 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1579, tokenIndex1579, depth1579 := position, tokenIndex, depth
			{
				position1580 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1579
				}
				position++
			l1581:
				{
					position1582, tokenIndex1582, depth1582 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1582
					}
					position++
					goto l1581
				l1582:
					position, tokenIndex, depth = position1582, tokenIndex1582, depth1582
				}
				if !rules[ruleskip]() {
					goto l1579
				}
				depth--
				add(ruleINTEGER, position1580)
			}
			return true
		l1579:
			position, tokenIndex, depth = position1579, tokenIndex1579, depth1579
			return false
		},
		/* 114 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 115 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 116 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 117 OR <- <('|' '|' skip)> */
		nil,
		/* 118 AND <- <('&' '&' skip)> */
		nil,
		/* 119 EQ <- <('=' skip)> */
		func() bool {
			position1588, tokenIndex1588, depth1588 := position, tokenIndex, depth
			{
				position1589 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1588
				}
				position++
				if !rules[ruleskip]() {
					goto l1588
				}
				depth--
				add(ruleEQ, position1589)
			}
			return true
		l1588:
			position, tokenIndex, depth = position1588, tokenIndex1588, depth1588
			return false
		},
		/* 120 NE <- <('!' '=' skip)> */
		nil,
		/* 121 GT <- <('>' skip)> */
		nil,
		/* 122 LT <- <('<' skip)> */
		nil,
		/* 123 LE <- <('<' '=' skip)> */
		nil,
		/* 124 GE <- <('>' '=' skip)> */
		nil,
		/* 125 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 126 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 127 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1597, tokenIndex1597, depth1597 := position, tokenIndex, depth
			{
				position1598 := position
				depth++
				{
					position1599, tokenIndex1599, depth1599 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1600
					}
					position++
					goto l1599
				l1600:
					position, tokenIndex, depth = position1599, tokenIndex1599, depth1599
					if buffer[position] != rune('A') {
						goto l1597
					}
					position++
				}
			l1599:
				{
					position1601, tokenIndex1601, depth1601 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1602
					}
					position++
					goto l1601
				l1602:
					position, tokenIndex, depth = position1601, tokenIndex1601, depth1601
					if buffer[position] != rune('S') {
						goto l1597
					}
					position++
				}
			l1601:
				if !rules[ruleskip]() {
					goto l1597
				}
				depth--
				add(ruleAS, position1598)
			}
			return true
		l1597:
			position, tokenIndex, depth = position1597, tokenIndex1597, depth1597
			return false
		},
		/* 128 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 129 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 130 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 131 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 132 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 133 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 134 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 135 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 136 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 137 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 138 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 139 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 140 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 141 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 142 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 143 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 144 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 145 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 146 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 147 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 148 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 149 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 150 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 151 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 152 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 153 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 154 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 155 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 156 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 157 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 158 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 159 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 160 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 161 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 162 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 163 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 164 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 165 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 166 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 167 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 168 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 169 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 170 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 171 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 172 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 173 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 174 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 175 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 176 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 177 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 178 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 179 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 180 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 181 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 182 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 183 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 184 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 185 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 186 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 187 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 188 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 189 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 190 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 191 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 192 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1668 := position
				depth++
			l1669:
				{
					position1670, tokenIndex1670, depth1670 := position, tokenIndex, depth
					{
						position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1672
						}
						goto l1671
					l1672:
						position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
						{
							position1673 := position
							depth++
							{
								position1674 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1670
								}
								position++
							l1675:
								{
									position1676, tokenIndex1676, depth1676 := position, tokenIndex, depth
									{
										position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1677
										}
										goto l1676
									l1677:
										position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
									}
									if !matchDot() {
										goto l1676
									}
									goto l1675
								l1676:
									position, tokenIndex, depth = position1676, tokenIndex1676, depth1676
								}
								if !rules[ruleendOfLine]() {
									goto l1670
								}
								depth--
								add(rulePegText, position1674)
							}
							{
								add(ruleAction13, position)
							}
							depth--
							add(rulecomment, position1673)
						}
					}
				l1671:
					goto l1669
				l1670:
					position, tokenIndex, depth = position1670, tokenIndex1670, depth1670
				}
				depth--
				add(ruleskip, position1668)
			}
			return true
		},
		/* 193 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1679, tokenIndex1679, depth1679 := position, tokenIndex, depth
			{
				position1680 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1679
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1679
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1679
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1679
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1679
						}
						break
					}
				}

				depth--
				add(rulews, position1680)
			}
			return true
		l1679:
			position, tokenIndex, depth = position1679, tokenIndex1679, depth1679
			return false
		},
		/* 194 comment <- <(<('#' (!endOfLine .)* endOfLine)> Action13)> */
		nil,
		/* 195 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1683, tokenIndex1683, depth1683 := position, tokenIndex, depth
			{
				position1684 := position
				depth++
				{
					position1685, tokenIndex1685, depth1685 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1686
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1686
					}
					position++
					goto l1685
				l1686:
					position, tokenIndex, depth = position1685, tokenIndex1685, depth1685
					if buffer[position] != rune('\n') {
						goto l1687
					}
					position++
					goto l1685
				l1687:
					position, tokenIndex, depth = position1685, tokenIndex1685, depth1685
					if buffer[position] != rune('\r') {
						goto l1683
					}
					position++
				}
			l1685:
				depth--
				add(ruleendOfLine, position1684)
			}
			return true
		l1683:
			position, tokenIndex, depth = position1683, tokenIndex1683, depth1683
			return false
		},
		nil,
		/* 198 Action0 <- <{ p.addPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 199 Action1 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 200 Action2 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 201 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 202 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 203 Action5 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 204 Action6 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 205 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 206 Action8 <- <{ p.setObject(p.skipComment(buffer, begin, end)); p.addTriplePattern() }> */
		nil,
		/* 207 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 208 Action10 <- <{ p.setPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 209 Action11 <- <{ p.setPathLength(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 210 Action12 <- <{ p.setKeyword(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 211 Action13 <- <{ p.commentBegin = begin }> */
		nil,
	}
	p.rules = rules
}
