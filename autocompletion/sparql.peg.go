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
	rulegroupCondition
	ruleorderCondition
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
	ruleASC
	ruleDESC
	ruleORDER
	ruleGROUP
	ruleBY
	ruleHAVING
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
	"groupCondition",
	"orderCondition",
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
	"ASC",
	"DESC",
	"ORDER",
	"GROUP",
	"BY",
	"HAVING",
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
	rules  [220]func() bool
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
		/* 7 subSelect <- <(select whereClause solutionModifier)> */
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
				if !rules[rulesolutionModifier]() {
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
					if !rules[ruleconstraint]() {
						goto l277
					}
					goto l276
				l277:
					position, tokenIndex, depth = position276, tokenIndex276, depth276
					{
						position291 := position
						depth++
						{
							position292, tokenIndex292, depth292 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l293
							}
							position++
							goto l292
						l293:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
							if buffer[position] != rune('B') {
								goto l274
							}
							position++
						}
					l292:
						{
							position294, tokenIndex294, depth294 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l295
							}
							position++
							goto l294
						l295:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							if buffer[position] != rune('I') {
								goto l274
							}
							position++
						}
					l294:
						{
							position296, tokenIndex296, depth296 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l297
							}
							position++
							goto l296
						l297:
							position, tokenIndex, depth = position296, tokenIndex296, depth296
							if buffer[position] != rune('N') {
								goto l274
							}
							position++
						}
					l296:
						{
							position298, tokenIndex298, depth298 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l299
							}
							position++
							goto l298
						l299:
							position, tokenIndex, depth = position298, tokenIndex298, depth298
							if buffer[position] != rune('D') {
								goto l274
							}
							position++
						}
					l298:
						if !rules[ruleskip]() {
							goto l274
						}
						depth--
						add(ruleBIND, position291)
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
		func() bool {
			position300, tokenIndex300, depth300 := position, tokenIndex, depth
			{
				position301 := position
				depth++
				{
					position302, tokenIndex302, depth302 := position, tokenIndex, depth
					if !rules[rulebrackettedExpression]() {
						goto l303
					}
					goto l302
				l303:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					if !rules[rulebuiltinCall]() {
						goto l304
					}
					goto l302
				l304:
					position, tokenIndex, depth = position302, tokenIndex302, depth302
					if !rules[rulefunctionCall]() {
						goto l300
					}
				}
			l302:
				depth--
				add(ruleconstraint, position301)
			}
			return true
		l300:
			position, tokenIndex, depth = position300, tokenIndex300, depth300
			return false
		},
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
		/* 43 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position422 := position
				depth++
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					{
						position425, tokenIndex425, depth425 := position, tokenIndex, depth
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
									goto l426
								}
								position++
							}
						l428:
							{
								position430, tokenIndex430, depth430 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l431
								}
								position++
								goto l430
							l431:
								position, tokenIndex, depth = position430, tokenIndex430, depth430
								if buffer[position] != rune('R') {
									goto l426
								}
								position++
							}
						l430:
							{
								position432, tokenIndex432, depth432 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l433
								}
								position++
								goto l432
							l433:
								position, tokenIndex, depth = position432, tokenIndex432, depth432
								if buffer[position] != rune('D') {
									goto l426
								}
								position++
							}
						l432:
							{
								position434, tokenIndex434, depth434 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l435
								}
								position++
								goto l434
							l435:
								position, tokenIndex, depth = position434, tokenIndex434, depth434
								if buffer[position] != rune('E') {
									goto l426
								}
								position++
							}
						l434:
							{
								position436, tokenIndex436, depth436 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l437
								}
								position++
								goto l436
							l437:
								position, tokenIndex, depth = position436, tokenIndex436, depth436
								if buffer[position] != rune('R') {
									goto l426
								}
								position++
							}
						l436:
							if !rules[ruleskip]() {
								goto l426
							}
							depth--
							add(ruleORDER, position427)
						}
						if !rules[ruleBY]() {
							goto l426
						}
						{
							position440 := position
							depth++
							{
								position441, tokenIndex441, depth441 := position, tokenIndex, depth
								{
									position443, tokenIndex443, depth443 := position, tokenIndex, depth
									{
										position445, tokenIndex445, depth445 := position, tokenIndex, depth
										{
											position447 := position
											depth++
											{
												position448, tokenIndex448, depth448 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l449
												}
												position++
												goto l448
											l449:
												position, tokenIndex, depth = position448, tokenIndex448, depth448
												if buffer[position] != rune('A') {
													goto l446
												}
												position++
											}
										l448:
											{
												position450, tokenIndex450, depth450 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l451
												}
												position++
												goto l450
											l451:
												position, tokenIndex, depth = position450, tokenIndex450, depth450
												if buffer[position] != rune('S') {
													goto l446
												}
												position++
											}
										l450:
											{
												position452, tokenIndex452, depth452 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l453
												}
												position++
												goto l452
											l453:
												position, tokenIndex, depth = position452, tokenIndex452, depth452
												if buffer[position] != rune('C') {
													goto l446
												}
												position++
											}
										l452:
											if !rules[ruleskip]() {
												goto l446
											}
											depth--
											add(ruleASC, position447)
										}
										goto l445
									l446:
										position, tokenIndex, depth = position445, tokenIndex445, depth445
										{
											position454 := position
											depth++
											{
												position455, tokenIndex455, depth455 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l456
												}
												position++
												goto l455
											l456:
												position, tokenIndex, depth = position455, tokenIndex455, depth455
												if buffer[position] != rune('D') {
													goto l443
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
													goto l443
												}
												position++
											}
										l457:
											{
												position459, tokenIndex459, depth459 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l460
												}
												position++
												goto l459
											l460:
												position, tokenIndex, depth = position459, tokenIndex459, depth459
												if buffer[position] != rune('S') {
													goto l443
												}
												position++
											}
										l459:
											{
												position461, tokenIndex461, depth461 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l462
												}
												position++
												goto l461
											l462:
												position, tokenIndex, depth = position461, tokenIndex461, depth461
												if buffer[position] != rune('C') {
													goto l443
												}
												position++
											}
										l461:
											if !rules[ruleskip]() {
												goto l443
											}
											depth--
											add(ruleDESC, position454)
										}
									}
								l445:
									goto l444
								l443:
									position, tokenIndex, depth = position443, tokenIndex443, depth443
								}
							l444:
								if !rules[rulebrackettedExpression]() {
									goto l442
								}
								goto l441
							l442:
								position, tokenIndex, depth = position441, tokenIndex441, depth441
								if !rules[rulefunctionCall]() {
									goto l463
								}
								goto l441
							l463:
								position, tokenIndex, depth = position441, tokenIndex441, depth441
								if !rules[rulebuiltinCall]() {
									goto l464
								}
								goto l441
							l464:
								position, tokenIndex, depth = position441, tokenIndex441, depth441
								if !rules[rulevar]() {
									goto l426
								}
							}
						l441:
							depth--
							add(ruleorderCondition, position440)
						}
					l438:
						{
							position439, tokenIndex439, depth439 := position, tokenIndex, depth
							{
								position465 := position
								depth++
								{
									position466, tokenIndex466, depth466 := position, tokenIndex, depth
									{
										position468, tokenIndex468, depth468 := position, tokenIndex, depth
										{
											position470, tokenIndex470, depth470 := position, tokenIndex, depth
											{
												position472 := position
												depth++
												{
													position473, tokenIndex473, depth473 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l474
													}
													position++
													goto l473
												l474:
													position, tokenIndex, depth = position473, tokenIndex473, depth473
													if buffer[position] != rune('A') {
														goto l471
													}
													position++
												}
											l473:
												{
													position475, tokenIndex475, depth475 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l476
													}
													position++
													goto l475
												l476:
													position, tokenIndex, depth = position475, tokenIndex475, depth475
													if buffer[position] != rune('S') {
														goto l471
													}
													position++
												}
											l475:
												{
													position477, tokenIndex477, depth477 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l478
													}
													position++
													goto l477
												l478:
													position, tokenIndex, depth = position477, tokenIndex477, depth477
													if buffer[position] != rune('C') {
														goto l471
													}
													position++
												}
											l477:
												if !rules[ruleskip]() {
													goto l471
												}
												depth--
												add(ruleASC, position472)
											}
											goto l470
										l471:
											position, tokenIndex, depth = position470, tokenIndex470, depth470
											{
												position479 := position
												depth++
												{
													position480, tokenIndex480, depth480 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l481
													}
													position++
													goto l480
												l481:
													position, tokenIndex, depth = position480, tokenIndex480, depth480
													if buffer[position] != rune('D') {
														goto l468
													}
													position++
												}
											l480:
												{
													position482, tokenIndex482, depth482 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l483
													}
													position++
													goto l482
												l483:
													position, tokenIndex, depth = position482, tokenIndex482, depth482
													if buffer[position] != rune('E') {
														goto l468
													}
													position++
												}
											l482:
												{
													position484, tokenIndex484, depth484 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l485
													}
													position++
													goto l484
												l485:
													position, tokenIndex, depth = position484, tokenIndex484, depth484
													if buffer[position] != rune('S') {
														goto l468
													}
													position++
												}
											l484:
												{
													position486, tokenIndex486, depth486 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l487
													}
													position++
													goto l486
												l487:
													position, tokenIndex, depth = position486, tokenIndex486, depth486
													if buffer[position] != rune('C') {
														goto l468
													}
													position++
												}
											l486:
												if !rules[ruleskip]() {
													goto l468
												}
												depth--
												add(ruleDESC, position479)
											}
										}
									l470:
										goto l469
									l468:
										position, tokenIndex, depth = position468, tokenIndex468, depth468
									}
								l469:
									if !rules[rulebrackettedExpression]() {
										goto l467
									}
									goto l466
								l467:
									position, tokenIndex, depth = position466, tokenIndex466, depth466
									if !rules[rulefunctionCall]() {
										goto l488
									}
									goto l466
								l488:
									position, tokenIndex, depth = position466, tokenIndex466, depth466
									if !rules[rulebuiltinCall]() {
										goto l489
									}
									goto l466
								l489:
									position, tokenIndex, depth = position466, tokenIndex466, depth466
									if !rules[rulevar]() {
										goto l439
									}
								}
							l466:
								depth--
								add(ruleorderCondition, position465)
							}
							goto l438
						l439:
							position, tokenIndex, depth = position439, tokenIndex439, depth439
						}
						goto l425
					l426:
						position, tokenIndex, depth = position425, tokenIndex425, depth425
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position491 := position
									depth++
									{
										position492, tokenIndex492, depth492 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l493
										}
										position++
										goto l492
									l493:
										position, tokenIndex, depth = position492, tokenIndex492, depth492
										if buffer[position] != rune('H') {
											goto l423
										}
										position++
									}
								l492:
									{
										position494, tokenIndex494, depth494 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l495
										}
										position++
										goto l494
									l495:
										position, tokenIndex, depth = position494, tokenIndex494, depth494
										if buffer[position] != rune('A') {
											goto l423
										}
										position++
									}
								l494:
									{
										position496, tokenIndex496, depth496 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l497
										}
										position++
										goto l496
									l497:
										position, tokenIndex, depth = position496, tokenIndex496, depth496
										if buffer[position] != rune('V') {
											goto l423
										}
										position++
									}
								l496:
									{
										position498, tokenIndex498, depth498 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l499
										}
										position++
										goto l498
									l499:
										position, tokenIndex, depth = position498, tokenIndex498, depth498
										if buffer[position] != rune('I') {
											goto l423
										}
										position++
									}
								l498:
									{
										position500, tokenIndex500, depth500 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l501
										}
										position++
										goto l500
									l501:
										position, tokenIndex, depth = position500, tokenIndex500, depth500
										if buffer[position] != rune('N') {
											goto l423
										}
										position++
									}
								l500:
									{
										position502, tokenIndex502, depth502 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l503
										}
										position++
										goto l502
									l503:
										position, tokenIndex, depth = position502, tokenIndex502, depth502
										if buffer[position] != rune('G') {
											goto l423
										}
										position++
									}
								l502:
									if !rules[ruleskip]() {
										goto l423
									}
									depth--
									add(ruleHAVING, position491)
								}
								if !rules[ruleconstraint]() {
									goto l423
								}
								break
							case 'G', 'g':
								{
									position504 := position
									depth++
									{
										position505, tokenIndex505, depth505 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l506
										}
										position++
										goto l505
									l506:
										position, tokenIndex, depth = position505, tokenIndex505, depth505
										if buffer[position] != rune('G') {
											goto l423
										}
										position++
									}
								l505:
									{
										position507, tokenIndex507, depth507 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l508
										}
										position++
										goto l507
									l508:
										position, tokenIndex, depth = position507, tokenIndex507, depth507
										if buffer[position] != rune('R') {
											goto l423
										}
										position++
									}
								l507:
									{
										position509, tokenIndex509, depth509 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l510
										}
										position++
										goto l509
									l510:
										position, tokenIndex, depth = position509, tokenIndex509, depth509
										if buffer[position] != rune('O') {
											goto l423
										}
										position++
									}
								l509:
									{
										position511, tokenIndex511, depth511 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l512
										}
										position++
										goto l511
									l512:
										position, tokenIndex, depth = position511, tokenIndex511, depth511
										if buffer[position] != rune('U') {
											goto l423
										}
										position++
									}
								l511:
									{
										position513, tokenIndex513, depth513 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l514
										}
										position++
										goto l513
									l514:
										position, tokenIndex, depth = position513, tokenIndex513, depth513
										if buffer[position] != rune('P') {
											goto l423
										}
										position++
									}
								l513:
									if !rules[ruleskip]() {
										goto l423
									}
									depth--
									add(ruleGROUP, position504)
								}
								if !rules[ruleBY]() {
									goto l423
								}
								{
									position517 := position
									depth++
									{
										position518, tokenIndex518, depth518 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l519
										}
										goto l518
									l519:
										position, tokenIndex, depth = position518, tokenIndex518, depth518
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l423
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l423
												}
												if !rules[ruleexpression]() {
													goto l423
												}
												{
													position521, tokenIndex521, depth521 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l521
													}
													if !rules[rulevar]() {
														goto l521
													}
													goto l522
												l521:
													position, tokenIndex, depth = position521, tokenIndex521, depth521
												}
											l522:
												if !rules[ruleRPAREN]() {
													goto l423
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l423
												}
												break
											}
										}

									}
								l518:
									depth--
									add(rulegroupCondition, position517)
								}
							l515:
								{
									position516, tokenIndex516, depth516 := position, tokenIndex, depth
									{
										position523 := position
										depth++
										{
											position524, tokenIndex524, depth524 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l525
											}
											goto l524
										l525:
											position, tokenIndex, depth = position524, tokenIndex524, depth524
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l516
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l516
													}
													if !rules[ruleexpression]() {
														goto l516
													}
													{
														position527, tokenIndex527, depth527 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l527
														}
														if !rules[rulevar]() {
															goto l527
														}
														goto l528
													l527:
														position, tokenIndex, depth = position527, tokenIndex527, depth527
													}
												l528:
													if !rules[ruleRPAREN]() {
														goto l516
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l516
													}
													break
												}
											}

										}
									l524:
										depth--
										add(rulegroupCondition, position523)
									}
									goto l515
								l516:
									position, tokenIndex, depth = position516, tokenIndex516, depth516
								}
								break
							default:
								{
									position529 := position
									depth++
									{
										position530, tokenIndex530, depth530 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l531
										}
										{
											position532, tokenIndex532, depth532 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l532
											}
											goto l533
										l532:
											position, tokenIndex, depth = position532, tokenIndex532, depth532
										}
									l533:
										goto l530
									l531:
										position, tokenIndex, depth = position530, tokenIndex530, depth530
										if !rules[ruleoffset]() {
											goto l423
										}
										{
											position534, tokenIndex534, depth534 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l534
											}
											goto l535
										l534:
											position, tokenIndex, depth = position534, tokenIndex534, depth534
										}
									l535:
									}
								l530:
									depth--
									add(rulelimitOffsetClauses, position529)
								}
								break
							}
						}

					}
				l425:
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
		/* 44 groupCondition <- <(functionCall / ((&('$' | '?') var) | (&('(') (LPAREN expression (AS var)? RPAREN)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'H' | 'I' | 'L' | 'M' | 'N' | 'R' | 'S' | 'T' | 'U' | 'Y' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'h' | 'i' | 'l' | 'm' | 'n' | 'r' | 's' | 't' | 'u' | 'y') builtinCall)))> */
		nil,
		/* 45 orderCondition <- <(((ASC / DESC)? brackettedExpression) / functionCall / builtinCall / var)> */
		nil,
		/* 46 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 47 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position539, tokenIndex539, depth539 := position, tokenIndex, depth
			{
				position540 := position
				depth++
				{
					position541 := position
					depth++
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l543
						}
						position++
						goto l542
					l543:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
						if buffer[position] != rune('L') {
							goto l539
						}
						position++
					}
				l542:
					{
						position544, tokenIndex544, depth544 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l545
						}
						position++
						goto l544
					l545:
						position, tokenIndex, depth = position544, tokenIndex544, depth544
						if buffer[position] != rune('I') {
							goto l539
						}
						position++
					}
				l544:
					{
						position546, tokenIndex546, depth546 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l547
						}
						position++
						goto l546
					l547:
						position, tokenIndex, depth = position546, tokenIndex546, depth546
						if buffer[position] != rune('M') {
							goto l539
						}
						position++
					}
				l546:
					{
						position548, tokenIndex548, depth548 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l549
						}
						position++
						goto l548
					l549:
						position, tokenIndex, depth = position548, tokenIndex548, depth548
						if buffer[position] != rune('I') {
							goto l539
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
							goto l539
						}
						position++
					}
				l550:
					if !rules[ruleskip]() {
						goto l539
					}
					depth--
					add(ruleLIMIT, position541)
				}
				if !rules[ruleINTEGER]() {
					goto l539
				}
				depth--
				add(rulelimit, position540)
			}
			return true
		l539:
			position, tokenIndex, depth = position539, tokenIndex539, depth539
			return false
		},
		/* 48 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position552, tokenIndex552, depth552 := position, tokenIndex, depth
			{
				position553 := position
				depth++
				{
					position554 := position
					depth++
					{
						position555, tokenIndex555, depth555 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l556
						}
						position++
						goto l555
					l556:
						position, tokenIndex, depth = position555, tokenIndex555, depth555
						if buffer[position] != rune('O') {
							goto l552
						}
						position++
					}
				l555:
					{
						position557, tokenIndex557, depth557 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l558
						}
						position++
						goto l557
					l558:
						position, tokenIndex, depth = position557, tokenIndex557, depth557
						if buffer[position] != rune('F') {
							goto l552
						}
						position++
					}
				l557:
					{
						position559, tokenIndex559, depth559 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l560
						}
						position++
						goto l559
					l560:
						position, tokenIndex, depth = position559, tokenIndex559, depth559
						if buffer[position] != rune('F') {
							goto l552
						}
						position++
					}
				l559:
					{
						position561, tokenIndex561, depth561 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l562
						}
						position++
						goto l561
					l562:
						position, tokenIndex, depth = position561, tokenIndex561, depth561
						if buffer[position] != rune('S') {
							goto l552
						}
						position++
					}
				l561:
					{
						position563, tokenIndex563, depth563 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l564
						}
						position++
						goto l563
					l564:
						position, tokenIndex, depth = position563, tokenIndex563, depth563
						if buffer[position] != rune('E') {
							goto l552
						}
						position++
					}
				l563:
					{
						position565, tokenIndex565, depth565 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l566
						}
						position++
						goto l565
					l566:
						position, tokenIndex, depth = position565, tokenIndex565, depth565
						if buffer[position] != rune('T') {
							goto l552
						}
						position++
					}
				l565:
					if !rules[ruleskip]() {
						goto l552
					}
					depth--
					add(ruleOFFSET, position554)
				}
				if !rules[ruleINTEGER]() {
					goto l552
				}
				depth--
				add(ruleoffset, position553)
			}
			return true
		l552:
			position, tokenIndex, depth = position552, tokenIndex552, depth552
			return false
		},
		/* 49 expression <- <conditionalOrExpression> */
		func() bool {
			position567, tokenIndex567, depth567 := position, tokenIndex, depth
			{
				position568 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l567
				}
				depth--
				add(ruleexpression, position568)
			}
			return true
		l567:
			position, tokenIndex, depth = position567, tokenIndex567, depth567
			return false
		},
		/* 50 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position569, tokenIndex569, depth569 := position, tokenIndex, depth
			{
				position570 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l569
				}
				{
					position571, tokenIndex571, depth571 := position, tokenIndex, depth
					{
						position573 := position
						depth++
						if buffer[position] != rune('|') {
							goto l571
						}
						position++
						if buffer[position] != rune('|') {
							goto l571
						}
						position++
						if !rules[ruleskip]() {
							goto l571
						}
						depth--
						add(ruleOR, position573)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l571
					}
					goto l572
				l571:
					position, tokenIndex, depth = position571, tokenIndex571, depth571
				}
			l572:
				depth--
				add(ruleconditionalOrExpression, position570)
			}
			return true
		l569:
			position, tokenIndex, depth = position569, tokenIndex569, depth569
			return false
		},
		/* 51 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position574, tokenIndex574, depth574 := position, tokenIndex, depth
			{
				position575 := position
				depth++
				{
					position576 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l574
					}
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position580 := position
									depth++
									{
										position581 := position
										depth++
										{
											position582, tokenIndex582, depth582 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l583
											}
											position++
											goto l582
										l583:
											position, tokenIndex, depth = position582, tokenIndex582, depth582
											if buffer[position] != rune('N') {
												goto l577
											}
											position++
										}
									l582:
										{
											position584, tokenIndex584, depth584 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l585
											}
											position++
											goto l584
										l585:
											position, tokenIndex, depth = position584, tokenIndex584, depth584
											if buffer[position] != rune('O') {
												goto l577
											}
											position++
										}
									l584:
										{
											position586, tokenIndex586, depth586 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l587
											}
											position++
											goto l586
										l587:
											position, tokenIndex, depth = position586, tokenIndex586, depth586
											if buffer[position] != rune('T') {
												goto l577
											}
											position++
										}
									l586:
										if buffer[position] != rune(' ') {
											goto l577
										}
										position++
										{
											position588, tokenIndex588, depth588 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l589
											}
											position++
											goto l588
										l589:
											position, tokenIndex, depth = position588, tokenIndex588, depth588
											if buffer[position] != rune('I') {
												goto l577
											}
											position++
										}
									l588:
										{
											position590, tokenIndex590, depth590 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l591
											}
											position++
											goto l590
										l591:
											position, tokenIndex, depth = position590, tokenIndex590, depth590
											if buffer[position] != rune('N') {
												goto l577
											}
											position++
										}
									l590:
										if !rules[ruleskip]() {
											goto l577
										}
										depth--
										add(ruleNOTIN, position581)
									}
									if !rules[ruleargList]() {
										goto l577
									}
									depth--
									add(rulenotin, position580)
								}
								break
							case 'I', 'i':
								{
									position592 := position
									depth++
									{
										position593 := position
										depth++
										{
											position594, tokenIndex594, depth594 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l595
											}
											position++
											goto l594
										l595:
											position, tokenIndex, depth = position594, tokenIndex594, depth594
											if buffer[position] != rune('I') {
												goto l577
											}
											position++
										}
									l594:
										{
											position596, tokenIndex596, depth596 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l597
											}
											position++
											goto l596
										l597:
											position, tokenIndex, depth = position596, tokenIndex596, depth596
											if buffer[position] != rune('N') {
												goto l577
											}
											position++
										}
									l596:
										if !rules[ruleskip]() {
											goto l577
										}
										depth--
										add(ruleIN, position593)
									}
									if !rules[ruleargList]() {
										goto l577
									}
									depth--
									add(rulein, position592)
								}
								break
							default:
								{
									position598, tokenIndex598, depth598 := position, tokenIndex, depth
									{
										position600 := position
										depth++
										if buffer[position] != rune('<') {
											goto l599
										}
										position++
										if !rules[ruleskip]() {
											goto l599
										}
										depth--
										add(ruleLT, position600)
									}
									goto l598
								l599:
									position, tokenIndex, depth = position598, tokenIndex598, depth598
									{
										position602 := position
										depth++
										if buffer[position] != rune('>') {
											goto l601
										}
										position++
										if buffer[position] != rune('=') {
											goto l601
										}
										position++
										if !rules[ruleskip]() {
											goto l601
										}
										depth--
										add(ruleGE, position602)
									}
									goto l598
								l601:
									position, tokenIndex, depth = position598, tokenIndex598, depth598
									{
										switch buffer[position] {
										case '>':
											{
												position604 := position
												depth++
												if buffer[position] != rune('>') {
													goto l577
												}
												position++
												if !rules[ruleskip]() {
													goto l577
												}
												depth--
												add(ruleGT, position604)
											}
											break
										case '<':
											{
												position605 := position
												depth++
												if buffer[position] != rune('<') {
													goto l577
												}
												position++
												if buffer[position] != rune('=') {
													goto l577
												}
												position++
												if !rules[ruleskip]() {
													goto l577
												}
												depth--
												add(ruleLE, position605)
											}
											break
										case '!':
											{
												position606 := position
												depth++
												if buffer[position] != rune('!') {
													goto l577
												}
												position++
												if buffer[position] != rune('=') {
													goto l577
												}
												position++
												if !rules[ruleskip]() {
													goto l577
												}
												depth--
												add(ruleNE, position606)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l577
											}
											break
										}
									}

								}
							l598:
								if !rules[rulenumericExpression]() {
									goto l577
								}
								break
							}
						}

						goto l578
					l577:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
					}
				l578:
					depth--
					add(rulevalueLogical, position576)
				}
				{
					position607, tokenIndex607, depth607 := position, tokenIndex, depth
					{
						position609 := position
						depth++
						if buffer[position] != rune('&') {
							goto l607
						}
						position++
						if buffer[position] != rune('&') {
							goto l607
						}
						position++
						if !rules[ruleskip]() {
							goto l607
						}
						depth--
						add(ruleAND, position609)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l607
					}
					goto l608
				l607:
					position, tokenIndex, depth = position607, tokenIndex607, depth607
				}
			l608:
				depth--
				add(ruleconditionalAndExpression, position575)
			}
			return true
		l574:
			position, tokenIndex, depth = position574, tokenIndex574, depth574
			return false
		},
		/* 52 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 53 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position611, tokenIndex611, depth611 := position, tokenIndex, depth
			{
				position612 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l611
				}
			l613:
				{
					position614, tokenIndex614, depth614 := position, tokenIndex, depth
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						{
							position617, tokenIndex617, depth617 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l618
							}
							goto l617
						l618:
							position, tokenIndex, depth = position617, tokenIndex617, depth617
							if !rules[ruleMINUS]() {
								goto l616
							}
						}
					l617:
						if !rules[rulemultiplicativeExpression]() {
							goto l616
						}
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						{
							position619 := position
							depth++
							{
								position620, tokenIndex620, depth620 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l621
								}
								position++
								goto l620
							l621:
								position, tokenIndex, depth = position620, tokenIndex620, depth620
								if buffer[position] != rune('-') {
									goto l614
								}
								position++
							}
						l620:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l614
							}
							position++
						l622:
							{
								position623, tokenIndex623, depth623 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l623
								}
								position++
								goto l622
							l623:
								position, tokenIndex, depth = position623, tokenIndex623, depth623
							}
							{
								position624, tokenIndex624, depth624 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l624
								}
								position++
							l626:
								{
									position627, tokenIndex627, depth627 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l627
									}
									position++
									goto l626
								l627:
									position, tokenIndex, depth = position627, tokenIndex627, depth627
								}
								goto l625
							l624:
								position, tokenIndex, depth = position624, tokenIndex624, depth624
							}
						l625:
							if !rules[ruleskip]() {
								goto l614
							}
							depth--
							add(rulesignedNumericLiteral, position619)
						}
					}
				l615:
					goto l613
				l614:
					position, tokenIndex, depth = position614, tokenIndex614, depth614
				}
				depth--
				add(rulenumericExpression, position612)
			}
			return true
		l611:
			position, tokenIndex, depth = position611, tokenIndex611, depth611
			return false
		},
		/* 54 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position628, tokenIndex628, depth628 := position, tokenIndex, depth
			{
				position629 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l628
				}
			l630:
				{
					position631, tokenIndex631, depth631 := position, tokenIndex, depth
					{
						position632, tokenIndex632, depth632 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l633
						}
						goto l632
					l633:
						position, tokenIndex, depth = position632, tokenIndex632, depth632
						if !rules[ruleSLASH]() {
							goto l631
						}
					}
				l632:
					if !rules[ruleunaryExpression]() {
						goto l631
					}
					goto l630
				l631:
					position, tokenIndex, depth = position631, tokenIndex631, depth631
				}
				depth--
				add(rulemultiplicativeExpression, position629)
			}
			return true
		l628:
			position, tokenIndex, depth = position628, tokenIndex628, depth628
			return false
		},
		/* 55 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position634, tokenIndex634, depth634 := position, tokenIndex, depth
			{
				position635 := position
				depth++
				{
					position636, tokenIndex636, depth636 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l636
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l636
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l636
							}
							break
						}
					}

					goto l637
				l636:
					position, tokenIndex, depth = position636, tokenIndex636, depth636
				}
			l637:
				{
					position639 := position
					depth++
					{
						position640, tokenIndex640, depth640 := position, tokenIndex, depth
						if !rules[rulebrackettedExpression]() {
							goto l641
						}
						goto l640
					l641:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
						if !rules[rulebuiltinCall]() {
							goto l642
						}
						goto l640
					l642:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
						if !rules[rulefunctionCall]() {
							goto l643
						}
						goto l640
					l643:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
						if !rules[ruleiriref]() {
							goto l644
						}
						goto l640
					l644:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position646 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position648 := position
												depth++
												{
													position649 := position
													depth++
													{
														position650, tokenIndex650, depth650 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l651
														}
														position++
														goto l650
													l651:
														position, tokenIndex, depth = position650, tokenIndex650, depth650
														if buffer[position] != rune('G') {
															goto l634
														}
														position++
													}
												l650:
													{
														position652, tokenIndex652, depth652 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l653
														}
														position++
														goto l652
													l653:
														position, tokenIndex, depth = position652, tokenIndex652, depth652
														if buffer[position] != rune('R') {
															goto l634
														}
														position++
													}
												l652:
													{
														position654, tokenIndex654, depth654 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l655
														}
														position++
														goto l654
													l655:
														position, tokenIndex, depth = position654, tokenIndex654, depth654
														if buffer[position] != rune('O') {
															goto l634
														}
														position++
													}
												l654:
													{
														position656, tokenIndex656, depth656 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l657
														}
														position++
														goto l656
													l657:
														position, tokenIndex, depth = position656, tokenIndex656, depth656
														if buffer[position] != rune('U') {
															goto l634
														}
														position++
													}
												l656:
													{
														position658, tokenIndex658, depth658 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l659
														}
														position++
														goto l658
													l659:
														position, tokenIndex, depth = position658, tokenIndex658, depth658
														if buffer[position] != rune('P') {
															goto l634
														}
														position++
													}
												l658:
													if buffer[position] != rune('_') {
														goto l634
													}
													position++
													{
														position660, tokenIndex660, depth660 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l661
														}
														position++
														goto l660
													l661:
														position, tokenIndex, depth = position660, tokenIndex660, depth660
														if buffer[position] != rune('C') {
															goto l634
														}
														position++
													}
												l660:
													{
														position662, tokenIndex662, depth662 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l663
														}
														position++
														goto l662
													l663:
														position, tokenIndex, depth = position662, tokenIndex662, depth662
														if buffer[position] != rune('O') {
															goto l634
														}
														position++
													}
												l662:
													{
														position664, tokenIndex664, depth664 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l665
														}
														position++
														goto l664
													l665:
														position, tokenIndex, depth = position664, tokenIndex664, depth664
														if buffer[position] != rune('N') {
															goto l634
														}
														position++
													}
												l664:
													{
														position666, tokenIndex666, depth666 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l667
														}
														position++
														goto l666
													l667:
														position, tokenIndex, depth = position666, tokenIndex666, depth666
														if buffer[position] != rune('C') {
															goto l634
														}
														position++
													}
												l666:
													{
														position668, tokenIndex668, depth668 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l669
														}
														position++
														goto l668
													l669:
														position, tokenIndex, depth = position668, tokenIndex668, depth668
														if buffer[position] != rune('A') {
															goto l634
														}
														position++
													}
												l668:
													{
														position670, tokenIndex670, depth670 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l671
														}
														position++
														goto l670
													l671:
														position, tokenIndex, depth = position670, tokenIndex670, depth670
														if buffer[position] != rune('T') {
															goto l634
														}
														position++
													}
												l670:
													if !rules[ruleskip]() {
														goto l634
													}
													depth--
													add(ruleGROUPCONCAT, position649)
												}
												if !rules[ruleLPAREN]() {
													goto l634
												}
												{
													position672, tokenIndex672, depth672 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l672
													}
													goto l673
												l672:
													position, tokenIndex, depth = position672, tokenIndex672, depth672
												}
											l673:
												if !rules[ruleexpression]() {
													goto l634
												}
												{
													position674, tokenIndex674, depth674 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l674
													}
													{
														position676 := position
														depth++
														{
															position677, tokenIndex677, depth677 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l678
															}
															position++
															goto l677
														l678:
															position, tokenIndex, depth = position677, tokenIndex677, depth677
															if buffer[position] != rune('S') {
																goto l674
															}
															position++
														}
													l677:
														{
															position679, tokenIndex679, depth679 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l680
															}
															position++
															goto l679
														l680:
															position, tokenIndex, depth = position679, tokenIndex679, depth679
															if buffer[position] != rune('E') {
																goto l674
															}
															position++
														}
													l679:
														{
															position681, tokenIndex681, depth681 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l682
															}
															position++
															goto l681
														l682:
															position, tokenIndex, depth = position681, tokenIndex681, depth681
															if buffer[position] != rune('P') {
																goto l674
															}
															position++
														}
													l681:
														{
															position683, tokenIndex683, depth683 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l684
															}
															position++
															goto l683
														l684:
															position, tokenIndex, depth = position683, tokenIndex683, depth683
															if buffer[position] != rune('A') {
																goto l674
															}
															position++
														}
													l683:
														{
															position685, tokenIndex685, depth685 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l686
															}
															position++
															goto l685
														l686:
															position, tokenIndex, depth = position685, tokenIndex685, depth685
															if buffer[position] != rune('R') {
																goto l674
															}
															position++
														}
													l685:
														{
															position687, tokenIndex687, depth687 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l688
															}
															position++
															goto l687
														l688:
															position, tokenIndex, depth = position687, tokenIndex687, depth687
															if buffer[position] != rune('A') {
																goto l674
															}
															position++
														}
													l687:
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
																goto l674
															}
															position++
														}
													l689:
														{
															position691, tokenIndex691, depth691 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l692
															}
															position++
															goto l691
														l692:
															position, tokenIndex, depth = position691, tokenIndex691, depth691
															if buffer[position] != rune('O') {
																goto l674
															}
															position++
														}
													l691:
														{
															position693, tokenIndex693, depth693 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l694
															}
															position++
															goto l693
														l694:
															position, tokenIndex, depth = position693, tokenIndex693, depth693
															if buffer[position] != rune('R') {
																goto l674
															}
															position++
														}
													l693:
														if !rules[ruleskip]() {
															goto l674
														}
														depth--
														add(ruleSEPARATOR, position676)
													}
													if !rules[ruleEQ]() {
														goto l674
													}
													if !rules[rulestring]() {
														goto l674
													}
													goto l675
												l674:
													position, tokenIndex, depth = position674, tokenIndex674, depth674
												}
											l675:
												if !rules[ruleRPAREN]() {
													goto l634
												}
												depth--
												add(rulegroupConcat, position648)
											}
											break
										case 'C', 'c':
											{
												position695 := position
												depth++
												{
													position696 := position
													depth++
													{
														position697, tokenIndex697, depth697 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l698
														}
														position++
														goto l697
													l698:
														position, tokenIndex, depth = position697, tokenIndex697, depth697
														if buffer[position] != rune('C') {
															goto l634
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
															goto l634
														}
														position++
													}
												l699:
													{
														position701, tokenIndex701, depth701 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l702
														}
														position++
														goto l701
													l702:
														position, tokenIndex, depth = position701, tokenIndex701, depth701
														if buffer[position] != rune('U') {
															goto l634
														}
														position++
													}
												l701:
													{
														position703, tokenIndex703, depth703 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l704
														}
														position++
														goto l703
													l704:
														position, tokenIndex, depth = position703, tokenIndex703, depth703
														if buffer[position] != rune('N') {
															goto l634
														}
														position++
													}
												l703:
													{
														position705, tokenIndex705, depth705 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l706
														}
														position++
														goto l705
													l706:
														position, tokenIndex, depth = position705, tokenIndex705, depth705
														if buffer[position] != rune('T') {
															goto l634
														}
														position++
													}
												l705:
													if !rules[ruleskip]() {
														goto l634
													}
													depth--
													add(ruleCOUNT, position696)
												}
												if !rules[ruleLPAREN]() {
													goto l634
												}
												{
													position707, tokenIndex707, depth707 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l707
													}
													goto l708
												l707:
													position, tokenIndex, depth = position707, tokenIndex707, depth707
												}
											l708:
												{
													position709, tokenIndex709, depth709 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l710
													}
													goto l709
												l710:
													position, tokenIndex, depth = position709, tokenIndex709, depth709
													if !rules[ruleexpression]() {
														goto l634
													}
												}
											l709:
												if !rules[ruleRPAREN]() {
													goto l634
												}
												depth--
												add(rulecount, position695)
											}
											break
										default:
											{
												position711, tokenIndex711, depth711 := position, tokenIndex, depth
												{
													position713 := position
													depth++
													{
														position714, tokenIndex714, depth714 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l715
														}
														position++
														goto l714
													l715:
														position, tokenIndex, depth = position714, tokenIndex714, depth714
														if buffer[position] != rune('S') {
															goto l712
														}
														position++
													}
												l714:
													{
														position716, tokenIndex716, depth716 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l717
														}
														position++
														goto l716
													l717:
														position, tokenIndex, depth = position716, tokenIndex716, depth716
														if buffer[position] != rune('U') {
															goto l712
														}
														position++
													}
												l716:
													{
														position718, tokenIndex718, depth718 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l719
														}
														position++
														goto l718
													l719:
														position, tokenIndex, depth = position718, tokenIndex718, depth718
														if buffer[position] != rune('M') {
															goto l712
														}
														position++
													}
												l718:
													if !rules[ruleskip]() {
														goto l712
													}
													depth--
													add(ruleSUM, position713)
												}
												goto l711
											l712:
												position, tokenIndex, depth = position711, tokenIndex711, depth711
												{
													position721 := position
													depth++
													{
														position722, tokenIndex722, depth722 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l723
														}
														position++
														goto l722
													l723:
														position, tokenIndex, depth = position722, tokenIndex722, depth722
														if buffer[position] != rune('M') {
															goto l720
														}
														position++
													}
												l722:
													{
														position724, tokenIndex724, depth724 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l725
														}
														position++
														goto l724
													l725:
														position, tokenIndex, depth = position724, tokenIndex724, depth724
														if buffer[position] != rune('I') {
															goto l720
														}
														position++
													}
												l724:
													{
														position726, tokenIndex726, depth726 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l727
														}
														position++
														goto l726
													l727:
														position, tokenIndex, depth = position726, tokenIndex726, depth726
														if buffer[position] != rune('N') {
															goto l720
														}
														position++
													}
												l726:
													if !rules[ruleskip]() {
														goto l720
													}
													depth--
													add(ruleMIN, position721)
												}
												goto l711
											l720:
												position, tokenIndex, depth = position711, tokenIndex711, depth711
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position729 := position
															depth++
															{
																position730, tokenIndex730, depth730 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l731
																}
																position++
																goto l730
															l731:
																position, tokenIndex, depth = position730, tokenIndex730, depth730
																if buffer[position] != rune('S') {
																	goto l634
																}
																position++
															}
														l730:
															{
																position732, tokenIndex732, depth732 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l733
																}
																position++
																goto l732
															l733:
																position, tokenIndex, depth = position732, tokenIndex732, depth732
																if buffer[position] != rune('A') {
																	goto l634
																}
																position++
															}
														l732:
															{
																position734, tokenIndex734, depth734 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l735
																}
																position++
																goto l734
															l735:
																position, tokenIndex, depth = position734, tokenIndex734, depth734
																if buffer[position] != rune('M') {
																	goto l634
																}
																position++
															}
														l734:
															{
																position736, tokenIndex736, depth736 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l737
																}
																position++
																goto l736
															l737:
																position, tokenIndex, depth = position736, tokenIndex736, depth736
																if buffer[position] != rune('P') {
																	goto l634
																}
																position++
															}
														l736:
															{
																position738, tokenIndex738, depth738 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l739
																}
																position++
																goto l738
															l739:
																position, tokenIndex, depth = position738, tokenIndex738, depth738
																if buffer[position] != rune('L') {
																	goto l634
																}
																position++
															}
														l738:
															{
																position740, tokenIndex740, depth740 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l741
																}
																position++
																goto l740
															l741:
																position, tokenIndex, depth = position740, tokenIndex740, depth740
																if buffer[position] != rune('E') {
																	goto l634
																}
																position++
															}
														l740:
															if !rules[ruleskip]() {
																goto l634
															}
															depth--
															add(ruleSAMPLE, position729)
														}
														break
													case 'A', 'a':
														{
															position742 := position
															depth++
															{
																position743, tokenIndex743, depth743 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l744
																}
																position++
																goto l743
															l744:
																position, tokenIndex, depth = position743, tokenIndex743, depth743
																if buffer[position] != rune('A') {
																	goto l634
																}
																position++
															}
														l743:
															{
																position745, tokenIndex745, depth745 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l746
																}
																position++
																goto l745
															l746:
																position, tokenIndex, depth = position745, tokenIndex745, depth745
																if buffer[position] != rune('V') {
																	goto l634
																}
																position++
															}
														l745:
															{
																position747, tokenIndex747, depth747 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l748
																}
																position++
																goto l747
															l748:
																position, tokenIndex, depth = position747, tokenIndex747, depth747
																if buffer[position] != rune('G') {
																	goto l634
																}
																position++
															}
														l747:
															if !rules[ruleskip]() {
																goto l634
															}
															depth--
															add(ruleAVG, position742)
														}
														break
													default:
														{
															position749 := position
															depth++
															{
																position750, tokenIndex750, depth750 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l751
																}
																position++
																goto l750
															l751:
																position, tokenIndex, depth = position750, tokenIndex750, depth750
																if buffer[position] != rune('M') {
																	goto l634
																}
																position++
															}
														l750:
															{
																position752, tokenIndex752, depth752 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l753
																}
																position++
																goto l752
															l753:
																position, tokenIndex, depth = position752, tokenIndex752, depth752
																if buffer[position] != rune('A') {
																	goto l634
																}
																position++
															}
														l752:
															{
																position754, tokenIndex754, depth754 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l755
																}
																position++
																goto l754
															l755:
																position, tokenIndex, depth = position754, tokenIndex754, depth754
																if buffer[position] != rune('X') {
																	goto l634
																}
																position++
															}
														l754:
															if !rules[ruleskip]() {
																goto l634
															}
															depth--
															add(ruleMAX, position749)
														}
														break
													}
												}

											}
										l711:
											if !rules[ruleLPAREN]() {
												goto l634
											}
											{
												position756, tokenIndex756, depth756 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l756
												}
												goto l757
											l756:
												position, tokenIndex, depth = position756, tokenIndex756, depth756
											}
										l757:
											if !rules[ruleexpression]() {
												goto l634
											}
											if !rules[ruleRPAREN]() {
												goto l634
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position646)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l634
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l634
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l634
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l634
								}
								break
							}
						}

					}
				l640:
					depth--
					add(ruleprimaryExpression, position639)
				}
				depth--
				add(ruleunaryExpression, position635)
			}
			return true
		l634:
			position, tokenIndex, depth = position634, tokenIndex634, depth634
			return false
		},
		/* 56 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 57 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position759, tokenIndex759, depth759 := position, tokenIndex, depth
			{
				position760 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l759
				}
				if !rules[ruleexpression]() {
					goto l759
				}
				if !rules[ruleRPAREN]() {
					goto l759
				}
				depth--
				add(rulebrackettedExpression, position760)
			}
			return true
		l759:
			position, tokenIndex, depth = position759, tokenIndex759, depth759
			return false
		},
		/* 58 functionCall <- <(iriref argList)> */
		func() bool {
			position761, tokenIndex761, depth761 := position, tokenIndex, depth
			{
				position762 := position
				depth++
				if !rules[ruleiriref]() {
					goto l761
				}
				if !rules[ruleargList]() {
					goto l761
				}
				depth--
				add(rulefunctionCall, position762)
			}
			return true
		l761:
			position, tokenIndex, depth = position761, tokenIndex761, depth761
			return false
		},
		/* 59 in <- <(IN argList)> */
		nil,
		/* 60 notin <- <(NOTIN argList)> */
		nil,
		/* 61 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position765, tokenIndex765, depth765 := position, tokenIndex, depth
			{
				position766 := position
				depth++
				{
					position767, tokenIndex767, depth767 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l768
					}
					goto l767
				l768:
					position, tokenIndex, depth = position767, tokenIndex767, depth767
					if !rules[ruleLPAREN]() {
						goto l765
					}
					if !rules[ruleexpression]() {
						goto l765
					}
				l769:
					{
						position770, tokenIndex770, depth770 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l770
						}
						if !rules[ruleexpression]() {
							goto l770
						}
						goto l769
					l770:
						position, tokenIndex, depth = position770, tokenIndex770, depth770
					}
					if !rules[ruleRPAREN]() {
						goto l765
					}
				}
			l767:
				depth--
				add(ruleargList, position766)
			}
			return true
		l765:
			position, tokenIndex, depth = position765, tokenIndex765, depth765
			return false
		},
		/* 62 aggregate <- <((&('G' | 'g') groupConcat) | (&('C' | 'c') count) | (&('A' | 'M' | 'S' | 'a' | 'm' | 's') ((SUM / MIN / ((&('S' | 's') SAMPLE) | (&('A' | 'a') AVG) | (&('M' | 'm') MAX))) LPAREN DISTINCT? expression RPAREN)))> */
		nil,
		/* 63 count <- <(COUNT LPAREN DISTINCT? (STAR / expression) RPAREN)> */
		nil,
		/* 64 groupConcat <- <(GROUPCONCAT LPAREN DISTINCT? expression (SEMICOLON SEPARATOR EQ string)? RPAREN)> */
		nil,
		/* 65 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		func() bool {
			position774, tokenIndex774, depth774 := position, tokenIndex, depth
			{
				position775 := position
				depth++
				{
					position776, tokenIndex776, depth776 := position, tokenIndex, depth
					{
						position778, tokenIndex778, depth778 := position, tokenIndex, depth
						{
							position780 := position
							depth++
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
									goto l779
								}
								position++
							}
						l781:
							{
								position783, tokenIndex783, depth783 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l784
								}
								position++
								goto l783
							l784:
								position, tokenIndex, depth = position783, tokenIndex783, depth783
								if buffer[position] != rune('T') {
									goto l779
								}
								position++
							}
						l783:
							{
								position785, tokenIndex785, depth785 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l786
								}
								position++
								goto l785
							l786:
								position, tokenIndex, depth = position785, tokenIndex785, depth785
								if buffer[position] != rune('R') {
									goto l779
								}
								position++
							}
						l785:
							if !rules[ruleskip]() {
								goto l779
							}
							depth--
							add(ruleSTR, position780)
						}
						goto l778
					l779:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position788 := position
							depth++
							{
								position789, tokenIndex789, depth789 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l790
								}
								position++
								goto l789
							l790:
								position, tokenIndex, depth = position789, tokenIndex789, depth789
								if buffer[position] != rune('L') {
									goto l787
								}
								position++
							}
						l789:
							{
								position791, tokenIndex791, depth791 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l792
								}
								position++
								goto l791
							l792:
								position, tokenIndex, depth = position791, tokenIndex791, depth791
								if buffer[position] != rune('A') {
									goto l787
								}
								position++
							}
						l791:
							{
								position793, tokenIndex793, depth793 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l794
								}
								position++
								goto l793
							l794:
								position, tokenIndex, depth = position793, tokenIndex793, depth793
								if buffer[position] != rune('N') {
									goto l787
								}
								position++
							}
						l793:
							{
								position795, tokenIndex795, depth795 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l796
								}
								position++
								goto l795
							l796:
								position, tokenIndex, depth = position795, tokenIndex795, depth795
								if buffer[position] != rune('G') {
									goto l787
								}
								position++
							}
						l795:
							if !rules[ruleskip]() {
								goto l787
							}
							depth--
							add(ruleLANG, position788)
						}
						goto l778
					l787:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position798 := position
							depth++
							{
								position799, tokenIndex799, depth799 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l800
								}
								position++
								goto l799
							l800:
								position, tokenIndex, depth = position799, tokenIndex799, depth799
								if buffer[position] != rune('D') {
									goto l797
								}
								position++
							}
						l799:
							{
								position801, tokenIndex801, depth801 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l802
								}
								position++
								goto l801
							l802:
								position, tokenIndex, depth = position801, tokenIndex801, depth801
								if buffer[position] != rune('A') {
									goto l797
								}
								position++
							}
						l801:
							{
								position803, tokenIndex803, depth803 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l804
								}
								position++
								goto l803
							l804:
								position, tokenIndex, depth = position803, tokenIndex803, depth803
								if buffer[position] != rune('T') {
									goto l797
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
									goto l797
								}
								position++
							}
						l805:
							{
								position807, tokenIndex807, depth807 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l808
								}
								position++
								goto l807
							l808:
								position, tokenIndex, depth = position807, tokenIndex807, depth807
								if buffer[position] != rune('T') {
									goto l797
								}
								position++
							}
						l807:
							{
								position809, tokenIndex809, depth809 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l810
								}
								position++
								goto l809
							l810:
								position, tokenIndex, depth = position809, tokenIndex809, depth809
								if buffer[position] != rune('Y') {
									goto l797
								}
								position++
							}
						l809:
							{
								position811, tokenIndex811, depth811 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l812
								}
								position++
								goto l811
							l812:
								position, tokenIndex, depth = position811, tokenIndex811, depth811
								if buffer[position] != rune('P') {
									goto l797
								}
								position++
							}
						l811:
							{
								position813, tokenIndex813, depth813 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l814
								}
								position++
								goto l813
							l814:
								position, tokenIndex, depth = position813, tokenIndex813, depth813
								if buffer[position] != rune('E') {
									goto l797
								}
								position++
							}
						l813:
							if !rules[ruleskip]() {
								goto l797
							}
							depth--
							add(ruleDATATYPE, position798)
						}
						goto l778
					l797:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position816 := position
							depth++
							{
								position817, tokenIndex817, depth817 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l818
								}
								position++
								goto l817
							l818:
								position, tokenIndex, depth = position817, tokenIndex817, depth817
								if buffer[position] != rune('I') {
									goto l815
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
									goto l815
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
									goto l815
								}
								position++
							}
						l821:
							if !rules[ruleskip]() {
								goto l815
							}
							depth--
							add(ruleIRI, position816)
						}
						goto l778
					l815:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position824 := position
							depth++
							{
								position825, tokenIndex825, depth825 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l826
								}
								position++
								goto l825
							l826:
								position, tokenIndex, depth = position825, tokenIndex825, depth825
								if buffer[position] != rune('U') {
									goto l823
								}
								position++
							}
						l825:
							{
								position827, tokenIndex827, depth827 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l828
								}
								position++
								goto l827
							l828:
								position, tokenIndex, depth = position827, tokenIndex827, depth827
								if buffer[position] != rune('R') {
									goto l823
								}
								position++
							}
						l827:
							{
								position829, tokenIndex829, depth829 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l830
								}
								position++
								goto l829
							l830:
								position, tokenIndex, depth = position829, tokenIndex829, depth829
								if buffer[position] != rune('I') {
									goto l823
								}
								position++
							}
						l829:
							if !rules[ruleskip]() {
								goto l823
							}
							depth--
							add(ruleURI, position824)
						}
						goto l778
					l823:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position832 := position
							depth++
							{
								position833, tokenIndex833, depth833 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l834
								}
								position++
								goto l833
							l834:
								position, tokenIndex, depth = position833, tokenIndex833, depth833
								if buffer[position] != rune('S') {
									goto l831
								}
								position++
							}
						l833:
							{
								position835, tokenIndex835, depth835 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l836
								}
								position++
								goto l835
							l836:
								position, tokenIndex, depth = position835, tokenIndex835, depth835
								if buffer[position] != rune('T') {
									goto l831
								}
								position++
							}
						l835:
							{
								position837, tokenIndex837, depth837 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l838
								}
								position++
								goto l837
							l838:
								position, tokenIndex, depth = position837, tokenIndex837, depth837
								if buffer[position] != rune('R') {
									goto l831
								}
								position++
							}
						l837:
							{
								position839, tokenIndex839, depth839 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l840
								}
								position++
								goto l839
							l840:
								position, tokenIndex, depth = position839, tokenIndex839, depth839
								if buffer[position] != rune('L') {
									goto l831
								}
								position++
							}
						l839:
							{
								position841, tokenIndex841, depth841 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l842
								}
								position++
								goto l841
							l842:
								position, tokenIndex, depth = position841, tokenIndex841, depth841
								if buffer[position] != rune('E') {
									goto l831
								}
								position++
							}
						l841:
							{
								position843, tokenIndex843, depth843 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l844
								}
								position++
								goto l843
							l844:
								position, tokenIndex, depth = position843, tokenIndex843, depth843
								if buffer[position] != rune('N') {
									goto l831
								}
								position++
							}
						l843:
							if !rules[ruleskip]() {
								goto l831
							}
							depth--
							add(ruleSTRLEN, position832)
						}
						goto l778
					l831:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position846 := position
							depth++
							{
								position847, tokenIndex847, depth847 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l848
								}
								position++
								goto l847
							l848:
								position, tokenIndex, depth = position847, tokenIndex847, depth847
								if buffer[position] != rune('M') {
									goto l845
								}
								position++
							}
						l847:
							{
								position849, tokenIndex849, depth849 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l850
								}
								position++
								goto l849
							l850:
								position, tokenIndex, depth = position849, tokenIndex849, depth849
								if buffer[position] != rune('O') {
									goto l845
								}
								position++
							}
						l849:
							{
								position851, tokenIndex851, depth851 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l852
								}
								position++
								goto l851
							l852:
								position, tokenIndex, depth = position851, tokenIndex851, depth851
								if buffer[position] != rune('N') {
									goto l845
								}
								position++
							}
						l851:
							{
								position853, tokenIndex853, depth853 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l854
								}
								position++
								goto l853
							l854:
								position, tokenIndex, depth = position853, tokenIndex853, depth853
								if buffer[position] != rune('T') {
									goto l845
								}
								position++
							}
						l853:
							{
								position855, tokenIndex855, depth855 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l856
								}
								position++
								goto l855
							l856:
								position, tokenIndex, depth = position855, tokenIndex855, depth855
								if buffer[position] != rune('H') {
									goto l845
								}
								position++
							}
						l855:
							if !rules[ruleskip]() {
								goto l845
							}
							depth--
							add(ruleMONTH, position846)
						}
						goto l778
					l845:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position858 := position
							depth++
							{
								position859, tokenIndex859, depth859 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l860
								}
								position++
								goto l859
							l860:
								position, tokenIndex, depth = position859, tokenIndex859, depth859
								if buffer[position] != rune('M') {
									goto l857
								}
								position++
							}
						l859:
							{
								position861, tokenIndex861, depth861 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l862
								}
								position++
								goto l861
							l862:
								position, tokenIndex, depth = position861, tokenIndex861, depth861
								if buffer[position] != rune('I') {
									goto l857
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
									goto l857
								}
								position++
							}
						l863:
							{
								position865, tokenIndex865, depth865 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l866
								}
								position++
								goto l865
							l866:
								position, tokenIndex, depth = position865, tokenIndex865, depth865
								if buffer[position] != rune('U') {
									goto l857
								}
								position++
							}
						l865:
							{
								position867, tokenIndex867, depth867 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l868
								}
								position++
								goto l867
							l868:
								position, tokenIndex, depth = position867, tokenIndex867, depth867
								if buffer[position] != rune('T') {
									goto l857
								}
								position++
							}
						l867:
							{
								position869, tokenIndex869, depth869 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l870
								}
								position++
								goto l869
							l870:
								position, tokenIndex, depth = position869, tokenIndex869, depth869
								if buffer[position] != rune('E') {
									goto l857
								}
								position++
							}
						l869:
							{
								position871, tokenIndex871, depth871 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l872
								}
								position++
								goto l871
							l872:
								position, tokenIndex, depth = position871, tokenIndex871, depth871
								if buffer[position] != rune('S') {
									goto l857
								}
								position++
							}
						l871:
							if !rules[ruleskip]() {
								goto l857
							}
							depth--
							add(ruleMINUTES, position858)
						}
						goto l778
					l857:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position874 := position
							depth++
							{
								position875, tokenIndex875, depth875 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l876
								}
								position++
								goto l875
							l876:
								position, tokenIndex, depth = position875, tokenIndex875, depth875
								if buffer[position] != rune('S') {
									goto l873
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
									goto l873
								}
								position++
							}
						l877:
							{
								position879, tokenIndex879, depth879 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l880
								}
								position++
								goto l879
							l880:
								position, tokenIndex, depth = position879, tokenIndex879, depth879
								if buffer[position] != rune('C') {
									goto l873
								}
								position++
							}
						l879:
							{
								position881, tokenIndex881, depth881 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l882
								}
								position++
								goto l881
							l882:
								position, tokenIndex, depth = position881, tokenIndex881, depth881
								if buffer[position] != rune('O') {
									goto l873
								}
								position++
							}
						l881:
							{
								position883, tokenIndex883, depth883 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l884
								}
								position++
								goto l883
							l884:
								position, tokenIndex, depth = position883, tokenIndex883, depth883
								if buffer[position] != rune('N') {
									goto l873
								}
								position++
							}
						l883:
							{
								position885, tokenIndex885, depth885 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l886
								}
								position++
								goto l885
							l886:
								position, tokenIndex, depth = position885, tokenIndex885, depth885
								if buffer[position] != rune('D') {
									goto l873
								}
								position++
							}
						l885:
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
									goto l873
								}
								position++
							}
						l887:
							if !rules[ruleskip]() {
								goto l873
							}
							depth--
							add(ruleSECONDS, position874)
						}
						goto l778
					l873:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position890 := position
							depth++
							{
								position891, tokenIndex891, depth891 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l892
								}
								position++
								goto l891
							l892:
								position, tokenIndex, depth = position891, tokenIndex891, depth891
								if buffer[position] != rune('T') {
									goto l889
								}
								position++
							}
						l891:
							{
								position893, tokenIndex893, depth893 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l894
								}
								position++
								goto l893
							l894:
								position, tokenIndex, depth = position893, tokenIndex893, depth893
								if buffer[position] != rune('I') {
									goto l889
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
									goto l889
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
									goto l889
								}
								position++
							}
						l897:
							{
								position899, tokenIndex899, depth899 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l900
								}
								position++
								goto l899
							l900:
								position, tokenIndex, depth = position899, tokenIndex899, depth899
								if buffer[position] != rune('Z') {
									goto l889
								}
								position++
							}
						l899:
							{
								position901, tokenIndex901, depth901 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l902
								}
								position++
								goto l901
							l902:
								position, tokenIndex, depth = position901, tokenIndex901, depth901
								if buffer[position] != rune('O') {
									goto l889
								}
								position++
							}
						l901:
							{
								position903, tokenIndex903, depth903 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l904
								}
								position++
								goto l903
							l904:
								position, tokenIndex, depth = position903, tokenIndex903, depth903
								if buffer[position] != rune('N') {
									goto l889
								}
								position++
							}
						l903:
							{
								position905, tokenIndex905, depth905 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l906
								}
								position++
								goto l905
							l906:
								position, tokenIndex, depth = position905, tokenIndex905, depth905
								if buffer[position] != rune('E') {
									goto l889
								}
								position++
							}
						l905:
							if !rules[ruleskip]() {
								goto l889
							}
							depth--
							add(ruleTIMEZONE, position890)
						}
						goto l778
					l889:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position908 := position
							depth++
							{
								position909, tokenIndex909, depth909 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l910
								}
								position++
								goto l909
							l910:
								position, tokenIndex, depth = position909, tokenIndex909, depth909
								if buffer[position] != rune('S') {
									goto l907
								}
								position++
							}
						l909:
							{
								position911, tokenIndex911, depth911 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l912
								}
								position++
								goto l911
							l912:
								position, tokenIndex, depth = position911, tokenIndex911, depth911
								if buffer[position] != rune('H') {
									goto l907
								}
								position++
							}
						l911:
							{
								position913, tokenIndex913, depth913 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l914
								}
								position++
								goto l913
							l914:
								position, tokenIndex, depth = position913, tokenIndex913, depth913
								if buffer[position] != rune('A') {
									goto l907
								}
								position++
							}
						l913:
							if buffer[position] != rune('1') {
								goto l907
							}
							position++
							if !rules[ruleskip]() {
								goto l907
							}
							depth--
							add(ruleSHA1, position908)
						}
						goto l778
					l907:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position916 := position
							depth++
							{
								position917, tokenIndex917, depth917 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l918
								}
								position++
								goto l917
							l918:
								position, tokenIndex, depth = position917, tokenIndex917, depth917
								if buffer[position] != rune('S') {
									goto l915
								}
								position++
							}
						l917:
							{
								position919, tokenIndex919, depth919 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l920
								}
								position++
								goto l919
							l920:
								position, tokenIndex, depth = position919, tokenIndex919, depth919
								if buffer[position] != rune('H') {
									goto l915
								}
								position++
							}
						l919:
							{
								position921, tokenIndex921, depth921 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l922
								}
								position++
								goto l921
							l922:
								position, tokenIndex, depth = position921, tokenIndex921, depth921
								if buffer[position] != rune('A') {
									goto l915
								}
								position++
							}
						l921:
							if buffer[position] != rune('2') {
								goto l915
							}
							position++
							if buffer[position] != rune('5') {
								goto l915
							}
							position++
							if buffer[position] != rune('6') {
								goto l915
							}
							position++
							if !rules[ruleskip]() {
								goto l915
							}
							depth--
							add(ruleSHA256, position916)
						}
						goto l778
					l915:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position924 := position
							depth++
							{
								position925, tokenIndex925, depth925 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l926
								}
								position++
								goto l925
							l926:
								position, tokenIndex, depth = position925, tokenIndex925, depth925
								if buffer[position] != rune('S') {
									goto l923
								}
								position++
							}
						l925:
							{
								position927, tokenIndex927, depth927 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l928
								}
								position++
								goto l927
							l928:
								position, tokenIndex, depth = position927, tokenIndex927, depth927
								if buffer[position] != rune('H') {
									goto l923
								}
								position++
							}
						l927:
							{
								position929, tokenIndex929, depth929 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l930
								}
								position++
								goto l929
							l930:
								position, tokenIndex, depth = position929, tokenIndex929, depth929
								if buffer[position] != rune('A') {
									goto l923
								}
								position++
							}
						l929:
							if buffer[position] != rune('3') {
								goto l923
							}
							position++
							if buffer[position] != rune('8') {
								goto l923
							}
							position++
							if buffer[position] != rune('4') {
								goto l923
							}
							position++
							if !rules[ruleskip]() {
								goto l923
							}
							depth--
							add(ruleSHA384, position924)
						}
						goto l778
					l923:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position932 := position
							depth++
							{
								position933, tokenIndex933, depth933 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l934
								}
								position++
								goto l933
							l934:
								position, tokenIndex, depth = position933, tokenIndex933, depth933
								if buffer[position] != rune('I') {
									goto l931
								}
								position++
							}
						l933:
							{
								position935, tokenIndex935, depth935 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l936
								}
								position++
								goto l935
							l936:
								position, tokenIndex, depth = position935, tokenIndex935, depth935
								if buffer[position] != rune('S') {
									goto l931
								}
								position++
							}
						l935:
							{
								position937, tokenIndex937, depth937 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l938
								}
								position++
								goto l937
							l938:
								position, tokenIndex, depth = position937, tokenIndex937, depth937
								if buffer[position] != rune('I') {
									goto l931
								}
								position++
							}
						l937:
							{
								position939, tokenIndex939, depth939 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l940
								}
								position++
								goto l939
							l940:
								position, tokenIndex, depth = position939, tokenIndex939, depth939
								if buffer[position] != rune('R') {
									goto l931
								}
								position++
							}
						l939:
							{
								position941, tokenIndex941, depth941 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l942
								}
								position++
								goto l941
							l942:
								position, tokenIndex, depth = position941, tokenIndex941, depth941
								if buffer[position] != rune('I') {
									goto l931
								}
								position++
							}
						l941:
							if !rules[ruleskip]() {
								goto l931
							}
							depth--
							add(ruleISIRI, position932)
						}
						goto l778
					l931:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position944 := position
							depth++
							{
								position945, tokenIndex945, depth945 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l946
								}
								position++
								goto l945
							l946:
								position, tokenIndex, depth = position945, tokenIndex945, depth945
								if buffer[position] != rune('I') {
									goto l943
								}
								position++
							}
						l945:
							{
								position947, tokenIndex947, depth947 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l948
								}
								position++
								goto l947
							l948:
								position, tokenIndex, depth = position947, tokenIndex947, depth947
								if buffer[position] != rune('S') {
									goto l943
								}
								position++
							}
						l947:
							{
								position949, tokenIndex949, depth949 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l950
								}
								position++
								goto l949
							l950:
								position, tokenIndex, depth = position949, tokenIndex949, depth949
								if buffer[position] != rune('U') {
									goto l943
								}
								position++
							}
						l949:
							{
								position951, tokenIndex951, depth951 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l952
								}
								position++
								goto l951
							l952:
								position, tokenIndex, depth = position951, tokenIndex951, depth951
								if buffer[position] != rune('R') {
									goto l943
								}
								position++
							}
						l951:
							{
								position953, tokenIndex953, depth953 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l954
								}
								position++
								goto l953
							l954:
								position, tokenIndex, depth = position953, tokenIndex953, depth953
								if buffer[position] != rune('I') {
									goto l943
								}
								position++
							}
						l953:
							if !rules[ruleskip]() {
								goto l943
							}
							depth--
							add(ruleISURI, position944)
						}
						goto l778
					l943:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position956 := position
							depth++
							{
								position957, tokenIndex957, depth957 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l958
								}
								position++
								goto l957
							l958:
								position, tokenIndex, depth = position957, tokenIndex957, depth957
								if buffer[position] != rune('I') {
									goto l955
								}
								position++
							}
						l957:
							{
								position959, tokenIndex959, depth959 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l960
								}
								position++
								goto l959
							l960:
								position, tokenIndex, depth = position959, tokenIndex959, depth959
								if buffer[position] != rune('S') {
									goto l955
								}
								position++
							}
						l959:
							{
								position961, tokenIndex961, depth961 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l962
								}
								position++
								goto l961
							l962:
								position, tokenIndex, depth = position961, tokenIndex961, depth961
								if buffer[position] != rune('B') {
									goto l955
								}
								position++
							}
						l961:
							{
								position963, tokenIndex963, depth963 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l964
								}
								position++
								goto l963
							l964:
								position, tokenIndex, depth = position963, tokenIndex963, depth963
								if buffer[position] != rune('L') {
									goto l955
								}
								position++
							}
						l963:
							{
								position965, tokenIndex965, depth965 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l966
								}
								position++
								goto l965
							l966:
								position, tokenIndex, depth = position965, tokenIndex965, depth965
								if buffer[position] != rune('A') {
									goto l955
								}
								position++
							}
						l965:
							{
								position967, tokenIndex967, depth967 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l968
								}
								position++
								goto l967
							l968:
								position, tokenIndex, depth = position967, tokenIndex967, depth967
								if buffer[position] != rune('N') {
									goto l955
								}
								position++
							}
						l967:
							{
								position969, tokenIndex969, depth969 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l970
								}
								position++
								goto l969
							l970:
								position, tokenIndex, depth = position969, tokenIndex969, depth969
								if buffer[position] != rune('K') {
									goto l955
								}
								position++
							}
						l969:
							if !rules[ruleskip]() {
								goto l955
							}
							depth--
							add(ruleISBLANK, position956)
						}
						goto l778
					l955:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							position972 := position
							depth++
							{
								position973, tokenIndex973, depth973 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l974
								}
								position++
								goto l973
							l974:
								position, tokenIndex, depth = position973, tokenIndex973, depth973
								if buffer[position] != rune('I') {
									goto l971
								}
								position++
							}
						l973:
							{
								position975, tokenIndex975, depth975 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l976
								}
								position++
								goto l975
							l976:
								position, tokenIndex, depth = position975, tokenIndex975, depth975
								if buffer[position] != rune('S') {
									goto l971
								}
								position++
							}
						l975:
							{
								position977, tokenIndex977, depth977 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l978
								}
								position++
								goto l977
							l978:
								position, tokenIndex, depth = position977, tokenIndex977, depth977
								if buffer[position] != rune('L') {
									goto l971
								}
								position++
							}
						l977:
							{
								position979, tokenIndex979, depth979 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l980
								}
								position++
								goto l979
							l980:
								position, tokenIndex, depth = position979, tokenIndex979, depth979
								if buffer[position] != rune('I') {
									goto l971
								}
								position++
							}
						l979:
							{
								position981, tokenIndex981, depth981 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l982
								}
								position++
								goto l981
							l982:
								position, tokenIndex, depth = position981, tokenIndex981, depth981
								if buffer[position] != rune('T') {
									goto l971
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
									goto l971
								}
								position++
							}
						l983:
							{
								position985, tokenIndex985, depth985 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l986
								}
								position++
								goto l985
							l986:
								position, tokenIndex, depth = position985, tokenIndex985, depth985
								if buffer[position] != rune('R') {
									goto l971
								}
								position++
							}
						l985:
							{
								position987, tokenIndex987, depth987 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l988
								}
								position++
								goto l987
							l988:
								position, tokenIndex, depth = position987, tokenIndex987, depth987
								if buffer[position] != rune('A') {
									goto l971
								}
								position++
							}
						l987:
							{
								position989, tokenIndex989, depth989 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l990
								}
								position++
								goto l989
							l990:
								position, tokenIndex, depth = position989, tokenIndex989, depth989
								if buffer[position] != rune('L') {
									goto l971
								}
								position++
							}
						l989:
							if !rules[ruleskip]() {
								goto l971
							}
							depth--
							add(ruleISLITERAL, position972)
						}
						goto l778
					l971:
						position, tokenIndex, depth = position778, tokenIndex778, depth778
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position992 := position
									depth++
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
											goto l777
										}
										position++
									}
								l993:
									{
										position995, tokenIndex995, depth995 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l996
										}
										position++
										goto l995
									l996:
										position, tokenIndex, depth = position995, tokenIndex995, depth995
										if buffer[position] != rune('S') {
											goto l777
										}
										position++
									}
								l995:
									{
										position997, tokenIndex997, depth997 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l998
										}
										position++
										goto l997
									l998:
										position, tokenIndex, depth = position997, tokenIndex997, depth997
										if buffer[position] != rune('N') {
											goto l777
										}
										position++
									}
								l997:
									{
										position999, tokenIndex999, depth999 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1000
										}
										position++
										goto l999
									l1000:
										position, tokenIndex, depth = position999, tokenIndex999, depth999
										if buffer[position] != rune('U') {
											goto l777
										}
										position++
									}
								l999:
									{
										position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1002
										}
										position++
										goto l1001
									l1002:
										position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
										if buffer[position] != rune('M') {
											goto l777
										}
										position++
									}
								l1001:
									{
										position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1004
										}
										position++
										goto l1003
									l1004:
										position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
										if buffer[position] != rune('E') {
											goto l777
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
											goto l777
										}
										position++
									}
								l1005:
									{
										position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1008
										}
										position++
										goto l1007
									l1008:
										position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
										if buffer[position] != rune('I') {
											goto l777
										}
										position++
									}
								l1007:
									{
										position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1010
										}
										position++
										goto l1009
									l1010:
										position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
										if buffer[position] != rune('C') {
											goto l777
										}
										position++
									}
								l1009:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleISNUMERIC, position992)
								}
								break
							case 'S', 's':
								{
									position1011 := position
									depth++
									{
										position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1013
										}
										position++
										goto l1012
									l1013:
										position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
										if buffer[position] != rune('S') {
											goto l777
										}
										position++
									}
								l1012:
									{
										position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1015
										}
										position++
										goto l1014
									l1015:
										position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
										if buffer[position] != rune('H') {
											goto l777
										}
										position++
									}
								l1014:
									{
										position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1017
										}
										position++
										goto l1016
									l1017:
										position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
										if buffer[position] != rune('A') {
											goto l777
										}
										position++
									}
								l1016:
									if buffer[position] != rune('5') {
										goto l777
									}
									position++
									if buffer[position] != rune('1') {
										goto l777
									}
									position++
									if buffer[position] != rune('2') {
										goto l777
									}
									position++
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleSHA512, position1011)
								}
								break
							case 'M', 'm':
								{
									position1018 := position
									depth++
									{
										position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1020
										}
										position++
										goto l1019
									l1020:
										position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
										if buffer[position] != rune('M') {
											goto l777
										}
										position++
									}
								l1019:
									{
										position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1022
										}
										position++
										goto l1021
									l1022:
										position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
										if buffer[position] != rune('D') {
											goto l777
										}
										position++
									}
								l1021:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleMD5, position1018)
								}
								break
							case 'T', 't':
								{
									position1023 := position
									depth++
									{
										position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1025
										}
										position++
										goto l1024
									l1025:
										position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
										if buffer[position] != rune('T') {
											goto l777
										}
										position++
									}
								l1024:
									{
										position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1027
										}
										position++
										goto l1026
									l1027:
										position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
										if buffer[position] != rune('Z') {
											goto l777
										}
										position++
									}
								l1026:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleTZ, position1023)
								}
								break
							case 'H', 'h':
								{
									position1028 := position
									depth++
									{
										position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1030
										}
										position++
										goto l1029
									l1030:
										position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
										if buffer[position] != rune('H') {
											goto l777
										}
										position++
									}
								l1029:
									{
										position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1032
										}
										position++
										goto l1031
									l1032:
										position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
										if buffer[position] != rune('O') {
											goto l777
										}
										position++
									}
								l1031:
									{
										position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1034
										}
										position++
										goto l1033
									l1034:
										position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
										if buffer[position] != rune('U') {
											goto l777
										}
										position++
									}
								l1033:
									{
										position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1036
										}
										position++
										goto l1035
									l1036:
										position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
										if buffer[position] != rune('R') {
											goto l777
										}
										position++
									}
								l1035:
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
											goto l777
										}
										position++
									}
								l1037:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleHOURS, position1028)
								}
								break
							case 'D', 'd':
								{
									position1039 := position
									depth++
									{
										position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1041
										}
										position++
										goto l1040
									l1041:
										position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
										if buffer[position] != rune('D') {
											goto l777
										}
										position++
									}
								l1040:
									{
										position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1043
										}
										position++
										goto l1042
									l1043:
										position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
										if buffer[position] != rune('A') {
											goto l777
										}
										position++
									}
								l1042:
									{
										position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1045
										}
										position++
										goto l1044
									l1045:
										position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
										if buffer[position] != rune('Y') {
											goto l777
										}
										position++
									}
								l1044:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleDAY, position1039)
								}
								break
							case 'Y', 'y':
								{
									position1046 := position
									depth++
									{
										position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1048
										}
										position++
										goto l1047
									l1048:
										position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
										if buffer[position] != rune('Y') {
											goto l777
										}
										position++
									}
								l1047:
									{
										position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1050
										}
										position++
										goto l1049
									l1050:
										position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
										if buffer[position] != rune('E') {
											goto l777
										}
										position++
									}
								l1049:
									{
										position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1052
										}
										position++
										goto l1051
									l1052:
										position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
										if buffer[position] != rune('A') {
											goto l777
										}
										position++
									}
								l1051:
									{
										position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1054
										}
										position++
										goto l1053
									l1054:
										position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
										if buffer[position] != rune('R') {
											goto l777
										}
										position++
									}
								l1053:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleYEAR, position1046)
								}
								break
							case 'E', 'e':
								{
									position1055 := position
									depth++
									{
										position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1057
										}
										position++
										goto l1056
									l1057:
										position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
										if buffer[position] != rune('E') {
											goto l777
										}
										position++
									}
								l1056:
									{
										position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1059
										}
										position++
										goto l1058
									l1059:
										position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
										if buffer[position] != rune('N') {
											goto l777
										}
										position++
									}
								l1058:
									{
										position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1061
										}
										position++
										goto l1060
									l1061:
										position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
										if buffer[position] != rune('C') {
											goto l777
										}
										position++
									}
								l1060:
									{
										position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1063
										}
										position++
										goto l1062
									l1063:
										position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
										if buffer[position] != rune('O') {
											goto l777
										}
										position++
									}
								l1062:
									{
										position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1065
										}
										position++
										goto l1064
									l1065:
										position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
										if buffer[position] != rune('D') {
											goto l777
										}
										position++
									}
								l1064:
									{
										position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1067
										}
										position++
										goto l1066
									l1067:
										position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
										if buffer[position] != rune('E') {
											goto l777
										}
										position++
									}
								l1066:
									if buffer[position] != rune('_') {
										goto l777
									}
									position++
									{
										position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1069
										}
										position++
										goto l1068
									l1069:
										position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
										if buffer[position] != rune('F') {
											goto l777
										}
										position++
									}
								l1068:
									{
										position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1071
										}
										position++
										goto l1070
									l1071:
										position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
										if buffer[position] != rune('O') {
											goto l777
										}
										position++
									}
								l1070:
									{
										position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1073
										}
										position++
										goto l1072
									l1073:
										position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
										if buffer[position] != rune('R') {
											goto l777
										}
										position++
									}
								l1072:
									if buffer[position] != rune('_') {
										goto l777
									}
									position++
									{
										position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1075
										}
										position++
										goto l1074
									l1075:
										position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
										if buffer[position] != rune('U') {
											goto l777
										}
										position++
									}
								l1074:
									{
										position1076, tokenIndex1076, depth1076 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1077
										}
										position++
										goto l1076
									l1077:
										position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
										if buffer[position] != rune('R') {
											goto l777
										}
										position++
									}
								l1076:
									{
										position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1079
										}
										position++
										goto l1078
									l1079:
										position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
										if buffer[position] != rune('I') {
											goto l777
										}
										position++
									}
								l1078:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleENCODEFORURI, position1055)
								}
								break
							case 'L', 'l':
								{
									position1080 := position
									depth++
									{
										position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1082
										}
										position++
										goto l1081
									l1082:
										position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
										if buffer[position] != rune('L') {
											goto l777
										}
										position++
									}
								l1081:
									{
										position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1084
										}
										position++
										goto l1083
									l1084:
										position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
										if buffer[position] != rune('C') {
											goto l777
										}
										position++
									}
								l1083:
									{
										position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1086
										}
										position++
										goto l1085
									l1086:
										position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
										if buffer[position] != rune('A') {
											goto l777
										}
										position++
									}
								l1085:
									{
										position1087, tokenIndex1087, depth1087 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1088
										}
										position++
										goto l1087
									l1088:
										position, tokenIndex, depth = position1087, tokenIndex1087, depth1087
										if buffer[position] != rune('S') {
											goto l777
										}
										position++
									}
								l1087:
									{
										position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1090
										}
										position++
										goto l1089
									l1090:
										position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
										if buffer[position] != rune('E') {
											goto l777
										}
										position++
									}
								l1089:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleLCASE, position1080)
								}
								break
							case 'U', 'u':
								{
									position1091 := position
									depth++
									{
										position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1093
										}
										position++
										goto l1092
									l1093:
										position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
										if buffer[position] != rune('U') {
											goto l777
										}
										position++
									}
								l1092:
									{
										position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1095
										}
										position++
										goto l1094
									l1095:
										position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
										if buffer[position] != rune('C') {
											goto l777
										}
										position++
									}
								l1094:
									{
										position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1097
										}
										position++
										goto l1096
									l1097:
										position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
										if buffer[position] != rune('A') {
											goto l777
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
											goto l777
										}
										position++
									}
								l1098:
									{
										position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1101
										}
										position++
										goto l1100
									l1101:
										position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
										if buffer[position] != rune('E') {
											goto l777
										}
										position++
									}
								l1100:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleUCASE, position1091)
								}
								break
							case 'F', 'f':
								{
									position1102 := position
									depth++
									{
										position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1104
										}
										position++
										goto l1103
									l1104:
										position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
										if buffer[position] != rune('F') {
											goto l777
										}
										position++
									}
								l1103:
									{
										position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1106
										}
										position++
										goto l1105
									l1106:
										position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
										if buffer[position] != rune('L') {
											goto l777
										}
										position++
									}
								l1105:
									{
										position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1108
										}
										position++
										goto l1107
									l1108:
										position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
										if buffer[position] != rune('O') {
											goto l777
										}
										position++
									}
								l1107:
									{
										position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1110
										}
										position++
										goto l1109
									l1110:
										position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
										if buffer[position] != rune('O') {
											goto l777
										}
										position++
									}
								l1109:
									{
										position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1112
										}
										position++
										goto l1111
									l1112:
										position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
										if buffer[position] != rune('R') {
											goto l777
										}
										position++
									}
								l1111:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleFLOOR, position1102)
								}
								break
							case 'R', 'r':
								{
									position1113 := position
									depth++
									{
										position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1115
										}
										position++
										goto l1114
									l1115:
										position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
										if buffer[position] != rune('R') {
											goto l777
										}
										position++
									}
								l1114:
									{
										position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1117
										}
										position++
										goto l1116
									l1117:
										position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
										if buffer[position] != rune('O') {
											goto l777
										}
										position++
									}
								l1116:
									{
										position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1119
										}
										position++
										goto l1118
									l1119:
										position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
										if buffer[position] != rune('U') {
											goto l777
										}
										position++
									}
								l1118:
									{
										position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1121
										}
										position++
										goto l1120
									l1121:
										position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
										if buffer[position] != rune('N') {
											goto l777
										}
										position++
									}
								l1120:
									{
										position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1123
										}
										position++
										goto l1122
									l1123:
										position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
										if buffer[position] != rune('D') {
											goto l777
										}
										position++
									}
								l1122:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleROUND, position1113)
								}
								break
							case 'C', 'c':
								{
									position1124 := position
									depth++
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
											goto l777
										}
										position++
									}
								l1125:
									{
										position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1128
										}
										position++
										goto l1127
									l1128:
										position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
										if buffer[position] != rune('E') {
											goto l777
										}
										position++
									}
								l1127:
									{
										position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1130
										}
										position++
										goto l1129
									l1130:
										position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
										if buffer[position] != rune('I') {
											goto l777
										}
										position++
									}
								l1129:
									{
										position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1132
										}
										position++
										goto l1131
									l1132:
										position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
										if buffer[position] != rune('L') {
											goto l777
										}
										position++
									}
								l1131:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleCEIL, position1124)
								}
								break
							default:
								{
									position1133 := position
									depth++
									{
										position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1135
										}
										position++
										goto l1134
									l1135:
										position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
										if buffer[position] != rune('A') {
											goto l777
										}
										position++
									}
								l1134:
									{
										position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1137
										}
										position++
										goto l1136
									l1137:
										position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
										if buffer[position] != rune('B') {
											goto l777
										}
										position++
									}
								l1136:
									{
										position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1139
										}
										position++
										goto l1138
									l1139:
										position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
										if buffer[position] != rune('S') {
											goto l777
										}
										position++
									}
								l1138:
									if !rules[ruleskip]() {
										goto l777
									}
									depth--
									add(ruleABS, position1133)
								}
								break
							}
						}

					}
				l778:
					if !rules[ruleLPAREN]() {
						goto l777
					}
					if !rules[ruleexpression]() {
						goto l777
					}
					if !rules[ruleRPAREN]() {
						goto l777
					}
					goto l776
				l777:
					position, tokenIndex, depth = position776, tokenIndex776, depth776
					{
						position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
						{
							position1143 := position
							depth++
							{
								position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1145
								}
								position++
								goto l1144
							l1145:
								position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
								if buffer[position] != rune('S') {
									goto l1142
								}
								position++
							}
						l1144:
							{
								position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1147
								}
								position++
								goto l1146
							l1147:
								position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
								if buffer[position] != rune('T') {
									goto l1142
								}
								position++
							}
						l1146:
							{
								position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1149
								}
								position++
								goto l1148
							l1149:
								position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
								if buffer[position] != rune('R') {
									goto l1142
								}
								position++
							}
						l1148:
							{
								position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1151
								}
								position++
								goto l1150
							l1151:
								position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
								if buffer[position] != rune('S') {
									goto l1142
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
									goto l1142
								}
								position++
							}
						l1152:
							{
								position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1155
								}
								position++
								goto l1154
							l1155:
								position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
								if buffer[position] != rune('A') {
									goto l1142
								}
								position++
							}
						l1154:
							{
								position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1157
								}
								position++
								goto l1156
							l1157:
								position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
								if buffer[position] != rune('R') {
									goto l1142
								}
								position++
							}
						l1156:
							{
								position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1159
								}
								position++
								goto l1158
							l1159:
								position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
								if buffer[position] != rune('T') {
									goto l1142
								}
								position++
							}
						l1158:
							{
								position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1161
								}
								position++
								goto l1160
							l1161:
								position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
								if buffer[position] != rune('S') {
									goto l1142
								}
								position++
							}
						l1160:
							if !rules[ruleskip]() {
								goto l1142
							}
							depth--
							add(ruleSTRSTARTS, position1143)
						}
						goto l1141
					l1142:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
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
								if buffer[position] != rune('t') {
									goto l1167
								}
								position++
								goto l1166
							l1167:
								position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
								if buffer[position] != rune('T') {
									goto l1162
								}
								position++
							}
						l1166:
							{
								position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1169
								}
								position++
								goto l1168
							l1169:
								position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
								if buffer[position] != rune('R') {
									goto l1162
								}
								position++
							}
						l1168:
							{
								position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1171
								}
								position++
								goto l1170
							l1171:
								position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
								if buffer[position] != rune('E') {
									goto l1162
								}
								position++
							}
						l1170:
							{
								position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1173
								}
								position++
								goto l1172
							l1173:
								position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
								if buffer[position] != rune('N') {
									goto l1162
								}
								position++
							}
						l1172:
							{
								position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1175
								}
								position++
								goto l1174
							l1175:
								position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
								if buffer[position] != rune('D') {
									goto l1162
								}
								position++
							}
						l1174:
							{
								position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1177
								}
								position++
								goto l1176
							l1177:
								position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
								if buffer[position] != rune('S') {
									goto l1162
								}
								position++
							}
						l1176:
							if !rules[ruleskip]() {
								goto l1162
							}
							depth--
							add(ruleSTRENDS, position1163)
						}
						goto l1141
					l1162:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
						{
							position1179 := position
							depth++
							{
								position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1181
								}
								position++
								goto l1180
							l1181:
								position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
								if buffer[position] != rune('S') {
									goto l1178
								}
								position++
							}
						l1180:
							{
								position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1183
								}
								position++
								goto l1182
							l1183:
								position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
								if buffer[position] != rune('T') {
									goto l1178
								}
								position++
							}
						l1182:
							{
								position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1185
								}
								position++
								goto l1184
							l1185:
								position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
								if buffer[position] != rune('R') {
									goto l1178
								}
								position++
							}
						l1184:
							{
								position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1187
								}
								position++
								goto l1186
							l1187:
								position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
								if buffer[position] != rune('B') {
									goto l1178
								}
								position++
							}
						l1186:
							{
								position1188, tokenIndex1188, depth1188 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1189
								}
								position++
								goto l1188
							l1189:
								position, tokenIndex, depth = position1188, tokenIndex1188, depth1188
								if buffer[position] != rune('E') {
									goto l1178
								}
								position++
							}
						l1188:
							{
								position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1191
								}
								position++
								goto l1190
							l1191:
								position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
								if buffer[position] != rune('F') {
									goto l1178
								}
								position++
							}
						l1190:
							{
								position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1193
								}
								position++
								goto l1192
							l1193:
								position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
								if buffer[position] != rune('O') {
									goto l1178
								}
								position++
							}
						l1192:
							{
								position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1195
								}
								position++
								goto l1194
							l1195:
								position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
								if buffer[position] != rune('R') {
									goto l1178
								}
								position++
							}
						l1194:
							{
								position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1197
								}
								position++
								goto l1196
							l1197:
								position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
								if buffer[position] != rune('E') {
									goto l1178
								}
								position++
							}
						l1196:
							if !rules[ruleskip]() {
								goto l1178
							}
							depth--
							add(ruleSTRBEFORE, position1179)
						}
						goto l1141
					l1178:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
						{
							position1199 := position
							depth++
							{
								position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1201
								}
								position++
								goto l1200
							l1201:
								position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
								if buffer[position] != rune('S') {
									goto l1198
								}
								position++
							}
						l1200:
							{
								position1202, tokenIndex1202, depth1202 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1203
								}
								position++
								goto l1202
							l1203:
								position, tokenIndex, depth = position1202, tokenIndex1202, depth1202
								if buffer[position] != rune('T') {
									goto l1198
								}
								position++
							}
						l1202:
							{
								position1204, tokenIndex1204, depth1204 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1205
								}
								position++
								goto l1204
							l1205:
								position, tokenIndex, depth = position1204, tokenIndex1204, depth1204
								if buffer[position] != rune('R') {
									goto l1198
								}
								position++
							}
						l1204:
							{
								position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1207
								}
								position++
								goto l1206
							l1207:
								position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
								if buffer[position] != rune('A') {
									goto l1198
								}
								position++
							}
						l1206:
							{
								position1208, tokenIndex1208, depth1208 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1209
								}
								position++
								goto l1208
							l1209:
								position, tokenIndex, depth = position1208, tokenIndex1208, depth1208
								if buffer[position] != rune('F') {
									goto l1198
								}
								position++
							}
						l1208:
							{
								position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1211
								}
								position++
								goto l1210
							l1211:
								position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
								if buffer[position] != rune('T') {
									goto l1198
								}
								position++
							}
						l1210:
							{
								position1212, tokenIndex1212, depth1212 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1213
								}
								position++
								goto l1212
							l1213:
								position, tokenIndex, depth = position1212, tokenIndex1212, depth1212
								if buffer[position] != rune('E') {
									goto l1198
								}
								position++
							}
						l1212:
							{
								position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1215
								}
								position++
								goto l1214
							l1215:
								position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
								if buffer[position] != rune('R') {
									goto l1198
								}
								position++
							}
						l1214:
							if !rules[ruleskip]() {
								goto l1198
							}
							depth--
							add(ruleSTRAFTER, position1199)
						}
						goto l1141
					l1198:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
						{
							position1217 := position
							depth++
							{
								position1218, tokenIndex1218, depth1218 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1219
								}
								position++
								goto l1218
							l1219:
								position, tokenIndex, depth = position1218, tokenIndex1218, depth1218
								if buffer[position] != rune('S') {
									goto l1216
								}
								position++
							}
						l1218:
							{
								position1220, tokenIndex1220, depth1220 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1221
								}
								position++
								goto l1220
							l1221:
								position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
								if buffer[position] != rune('T') {
									goto l1216
								}
								position++
							}
						l1220:
							{
								position1222, tokenIndex1222, depth1222 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1223
								}
								position++
								goto l1222
							l1223:
								position, tokenIndex, depth = position1222, tokenIndex1222, depth1222
								if buffer[position] != rune('R') {
									goto l1216
								}
								position++
							}
						l1222:
							{
								position1224, tokenIndex1224, depth1224 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1225
								}
								position++
								goto l1224
							l1225:
								position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
								if buffer[position] != rune('L') {
									goto l1216
								}
								position++
							}
						l1224:
							{
								position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1227
								}
								position++
								goto l1226
							l1227:
								position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
								if buffer[position] != rune('A') {
									goto l1216
								}
								position++
							}
						l1226:
							{
								position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1229
								}
								position++
								goto l1228
							l1229:
								position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
								if buffer[position] != rune('N') {
									goto l1216
								}
								position++
							}
						l1228:
							{
								position1230, tokenIndex1230, depth1230 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1231
								}
								position++
								goto l1230
							l1231:
								position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
								if buffer[position] != rune('G') {
									goto l1216
								}
								position++
							}
						l1230:
							if !rules[ruleskip]() {
								goto l1216
							}
							depth--
							add(ruleSTRLANG, position1217)
						}
						goto l1141
					l1216:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
						{
							position1233 := position
							depth++
							{
								position1234, tokenIndex1234, depth1234 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1235
								}
								position++
								goto l1234
							l1235:
								position, tokenIndex, depth = position1234, tokenIndex1234, depth1234
								if buffer[position] != rune('S') {
									goto l1232
								}
								position++
							}
						l1234:
							{
								position1236, tokenIndex1236, depth1236 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1237
								}
								position++
								goto l1236
							l1237:
								position, tokenIndex, depth = position1236, tokenIndex1236, depth1236
								if buffer[position] != rune('T') {
									goto l1232
								}
								position++
							}
						l1236:
							{
								position1238, tokenIndex1238, depth1238 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1239
								}
								position++
								goto l1238
							l1239:
								position, tokenIndex, depth = position1238, tokenIndex1238, depth1238
								if buffer[position] != rune('R') {
									goto l1232
								}
								position++
							}
						l1238:
							{
								position1240, tokenIndex1240, depth1240 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1241
								}
								position++
								goto l1240
							l1241:
								position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
								if buffer[position] != rune('D') {
									goto l1232
								}
								position++
							}
						l1240:
							{
								position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1243
								}
								position++
								goto l1242
							l1243:
								position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
								if buffer[position] != rune('T') {
									goto l1232
								}
								position++
							}
						l1242:
							if !rules[ruleskip]() {
								goto l1232
							}
							depth--
							add(ruleSTRDT, position1233)
						}
						goto l1141
					l1232:
						position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1245 := position
									depth++
									{
										position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1247
										}
										position++
										goto l1246
									l1247:
										position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
										if buffer[position] != rune('S') {
											goto l1140
										}
										position++
									}
								l1246:
									{
										position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1249
										}
										position++
										goto l1248
									l1249:
										position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
										if buffer[position] != rune('A') {
											goto l1140
										}
										position++
									}
								l1248:
									{
										position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1251
										}
										position++
										goto l1250
									l1251:
										position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
										if buffer[position] != rune('M') {
											goto l1140
										}
										position++
									}
								l1250:
									{
										position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1253
										}
										position++
										goto l1252
									l1253:
										position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
										if buffer[position] != rune('E') {
											goto l1140
										}
										position++
									}
								l1252:
									{
										position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1255
										}
										position++
										goto l1254
									l1255:
										position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
										if buffer[position] != rune('T') {
											goto l1140
										}
										position++
									}
								l1254:
									{
										position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1257
										}
										position++
										goto l1256
									l1257:
										position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
										if buffer[position] != rune('E') {
											goto l1140
										}
										position++
									}
								l1256:
									{
										position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1259
										}
										position++
										goto l1258
									l1259:
										position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
										if buffer[position] != rune('R') {
											goto l1140
										}
										position++
									}
								l1258:
									{
										position1260, tokenIndex1260, depth1260 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1261
										}
										position++
										goto l1260
									l1261:
										position, tokenIndex, depth = position1260, tokenIndex1260, depth1260
										if buffer[position] != rune('M') {
											goto l1140
										}
										position++
									}
								l1260:
									if !rules[ruleskip]() {
										goto l1140
									}
									depth--
									add(ruleSAMETERM, position1245)
								}
								break
							case 'C', 'c':
								{
									position1262 := position
									depth++
									{
										position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1264
										}
										position++
										goto l1263
									l1264:
										position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
										if buffer[position] != rune('C') {
											goto l1140
										}
										position++
									}
								l1263:
									{
										position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1266
										}
										position++
										goto l1265
									l1266:
										position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
										if buffer[position] != rune('O') {
											goto l1140
										}
										position++
									}
								l1265:
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
											goto l1140
										}
										position++
									}
								l1267:
									{
										position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1270
										}
										position++
										goto l1269
									l1270:
										position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
										if buffer[position] != rune('T') {
											goto l1140
										}
										position++
									}
								l1269:
									{
										position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1272
										}
										position++
										goto l1271
									l1272:
										position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
										if buffer[position] != rune('A') {
											goto l1140
										}
										position++
									}
								l1271:
									{
										position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1274
										}
										position++
										goto l1273
									l1274:
										position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
										if buffer[position] != rune('I') {
											goto l1140
										}
										position++
									}
								l1273:
									{
										position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1276
										}
										position++
										goto l1275
									l1276:
										position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
										if buffer[position] != rune('N') {
											goto l1140
										}
										position++
									}
								l1275:
									{
										position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1278
										}
										position++
										goto l1277
									l1278:
										position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
										if buffer[position] != rune('S') {
											goto l1140
										}
										position++
									}
								l1277:
									if !rules[ruleskip]() {
										goto l1140
									}
									depth--
									add(ruleCONTAINS, position1262)
								}
								break
							default:
								{
									position1279 := position
									depth++
									{
										position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1281
										}
										position++
										goto l1280
									l1281:
										position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
										if buffer[position] != rune('L') {
											goto l1140
										}
										position++
									}
								l1280:
									{
										position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1283
										}
										position++
										goto l1282
									l1283:
										position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
										if buffer[position] != rune('A') {
											goto l1140
										}
										position++
									}
								l1282:
									{
										position1284, tokenIndex1284, depth1284 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1285
										}
										position++
										goto l1284
									l1285:
										position, tokenIndex, depth = position1284, tokenIndex1284, depth1284
										if buffer[position] != rune('N') {
											goto l1140
										}
										position++
									}
								l1284:
									{
										position1286, tokenIndex1286, depth1286 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1287
										}
										position++
										goto l1286
									l1287:
										position, tokenIndex, depth = position1286, tokenIndex1286, depth1286
										if buffer[position] != rune('G') {
											goto l1140
										}
										position++
									}
								l1286:
									{
										position1288, tokenIndex1288, depth1288 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1289
										}
										position++
										goto l1288
									l1289:
										position, tokenIndex, depth = position1288, tokenIndex1288, depth1288
										if buffer[position] != rune('M') {
											goto l1140
										}
										position++
									}
								l1288:
									{
										position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1291
										}
										position++
										goto l1290
									l1291:
										position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
										if buffer[position] != rune('A') {
											goto l1140
										}
										position++
									}
								l1290:
									{
										position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1293
										}
										position++
										goto l1292
									l1293:
										position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
										if buffer[position] != rune('T') {
											goto l1140
										}
										position++
									}
								l1292:
									{
										position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1295
										}
										position++
										goto l1294
									l1295:
										position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
										if buffer[position] != rune('C') {
											goto l1140
										}
										position++
									}
								l1294:
									{
										position1296, tokenIndex1296, depth1296 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1297
										}
										position++
										goto l1296
									l1297:
										position, tokenIndex, depth = position1296, tokenIndex1296, depth1296
										if buffer[position] != rune('H') {
											goto l1140
										}
										position++
									}
								l1296:
									{
										position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1299
										}
										position++
										goto l1298
									l1299:
										position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
										if buffer[position] != rune('E') {
											goto l1140
										}
										position++
									}
								l1298:
									{
										position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1301
										}
										position++
										goto l1300
									l1301:
										position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
										if buffer[position] != rune('S') {
											goto l1140
										}
										position++
									}
								l1300:
									if !rules[ruleskip]() {
										goto l1140
									}
									depth--
									add(ruleLANGMATCHES, position1279)
								}
								break
							}
						}

					}
				l1141:
					if !rules[ruleLPAREN]() {
						goto l1140
					}
					if !rules[ruleexpression]() {
						goto l1140
					}
					if !rules[ruleCOMMA]() {
						goto l1140
					}
					if !rules[ruleexpression]() {
						goto l1140
					}
					if !rules[ruleRPAREN]() {
						goto l1140
					}
					goto l776
				l1140:
					position, tokenIndex, depth = position776, tokenIndex776, depth776
					{
						position1303 := position
						depth++
						{
							position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1305
							}
							position++
							goto l1304
						l1305:
							position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
							if buffer[position] != rune('B') {
								goto l1302
							}
							position++
						}
					l1304:
						{
							position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1307
							}
							position++
							goto l1306
						l1307:
							position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
							if buffer[position] != rune('O') {
								goto l1302
							}
							position++
						}
					l1306:
						{
							position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1309
							}
							position++
							goto l1308
						l1309:
							position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
							if buffer[position] != rune('U') {
								goto l1302
							}
							position++
						}
					l1308:
						{
							position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1311
							}
							position++
							goto l1310
						l1311:
							position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
							if buffer[position] != rune('N') {
								goto l1302
							}
							position++
						}
					l1310:
						{
							position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1313
							}
							position++
							goto l1312
						l1313:
							position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
							if buffer[position] != rune('D') {
								goto l1302
							}
							position++
						}
					l1312:
						if !rules[ruleskip]() {
							goto l1302
						}
						depth--
						add(ruleBOUND, position1303)
					}
					if !rules[ruleLPAREN]() {
						goto l1302
					}
					if !rules[rulevar]() {
						goto l1302
					}
					if !rules[ruleRPAREN]() {
						goto l1302
					}
					goto l776
				l1302:
					position, tokenIndex, depth = position776, tokenIndex776, depth776
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1316 := position
								depth++
								{
									position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1318
									}
									position++
									goto l1317
								l1318:
									position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
									if buffer[position] != rune('S') {
										goto l1314
									}
									position++
								}
							l1317:
								{
									position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1320
									}
									position++
									goto l1319
								l1320:
									position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
									if buffer[position] != rune('T') {
										goto l1314
									}
									position++
								}
							l1319:
								{
									position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1322
									}
									position++
									goto l1321
								l1322:
									position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
									if buffer[position] != rune('R') {
										goto l1314
									}
									position++
								}
							l1321:
								{
									position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1324
									}
									position++
									goto l1323
								l1324:
									position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
									if buffer[position] != rune('U') {
										goto l1314
									}
									position++
								}
							l1323:
								{
									position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1326
									}
									position++
									goto l1325
								l1326:
									position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
									if buffer[position] != rune('U') {
										goto l1314
									}
									position++
								}
							l1325:
								{
									position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1328
									}
									position++
									goto l1327
								l1328:
									position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
									if buffer[position] != rune('I') {
										goto l1314
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
										goto l1314
									}
									position++
								}
							l1329:
								if !rules[ruleskip]() {
									goto l1314
								}
								depth--
								add(ruleSTRUUID, position1316)
							}
							break
						case 'U', 'u':
							{
								position1331 := position
								depth++
								{
									position1332, tokenIndex1332, depth1332 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1333
									}
									position++
									goto l1332
								l1333:
									position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
									if buffer[position] != rune('U') {
										goto l1314
									}
									position++
								}
							l1332:
								{
									position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1335
									}
									position++
									goto l1334
								l1335:
									position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
									if buffer[position] != rune('U') {
										goto l1314
									}
									position++
								}
							l1334:
								{
									position1336, tokenIndex1336, depth1336 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1337
									}
									position++
									goto l1336
								l1337:
									position, tokenIndex, depth = position1336, tokenIndex1336, depth1336
									if buffer[position] != rune('I') {
										goto l1314
									}
									position++
								}
							l1336:
								{
									position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1339
									}
									position++
									goto l1338
								l1339:
									position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
									if buffer[position] != rune('D') {
										goto l1314
									}
									position++
								}
							l1338:
								if !rules[ruleskip]() {
									goto l1314
								}
								depth--
								add(ruleUUID, position1331)
							}
							break
						case 'N', 'n':
							{
								position1340 := position
								depth++
								{
									position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1342
									}
									position++
									goto l1341
								l1342:
									position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
									if buffer[position] != rune('N') {
										goto l1314
									}
									position++
								}
							l1341:
								{
									position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1344
									}
									position++
									goto l1343
								l1344:
									position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
									if buffer[position] != rune('O') {
										goto l1314
									}
									position++
								}
							l1343:
								{
									position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1346
									}
									position++
									goto l1345
								l1346:
									position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
									if buffer[position] != rune('W') {
										goto l1314
									}
									position++
								}
							l1345:
								if !rules[ruleskip]() {
									goto l1314
								}
								depth--
								add(ruleNOW, position1340)
							}
							break
						default:
							{
								position1347 := position
								depth++
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
										goto l1314
									}
									position++
								}
							l1348:
								{
									position1350, tokenIndex1350, depth1350 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1351
									}
									position++
									goto l1350
								l1351:
									position, tokenIndex, depth = position1350, tokenIndex1350, depth1350
									if buffer[position] != rune('A') {
										goto l1314
									}
									position++
								}
							l1350:
								{
									position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1353
									}
									position++
									goto l1352
								l1353:
									position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
									if buffer[position] != rune('N') {
										goto l1314
									}
									position++
								}
							l1352:
								{
									position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1355
									}
									position++
									goto l1354
								l1355:
									position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
									if buffer[position] != rune('D') {
										goto l1314
									}
									position++
								}
							l1354:
								if !rules[ruleskip]() {
									goto l1314
								}
								depth--
								add(ruleRAND, position1347)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1314
					}
					goto l776
				l1314:
					position, tokenIndex, depth = position776, tokenIndex776, depth776
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
								{
									position1359 := position
									depth++
									{
										position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1361
										}
										position++
										goto l1360
									l1361:
										position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
										if buffer[position] != rune('E') {
											goto l1358
										}
										position++
									}
								l1360:
									{
										position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1363
										}
										position++
										goto l1362
									l1363:
										position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
										if buffer[position] != rune('X') {
											goto l1358
										}
										position++
									}
								l1362:
									{
										position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1365
										}
										position++
										goto l1364
									l1365:
										position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
										if buffer[position] != rune('I') {
											goto l1358
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
											goto l1358
										}
										position++
									}
								l1366:
									{
										position1368, tokenIndex1368, depth1368 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1369
										}
										position++
										goto l1368
									l1369:
										position, tokenIndex, depth = position1368, tokenIndex1368, depth1368
										if buffer[position] != rune('T') {
											goto l1358
										}
										position++
									}
								l1368:
									{
										position1370, tokenIndex1370, depth1370 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1371
										}
										position++
										goto l1370
									l1371:
										position, tokenIndex, depth = position1370, tokenIndex1370, depth1370
										if buffer[position] != rune('S') {
											goto l1358
										}
										position++
									}
								l1370:
									if !rules[ruleskip]() {
										goto l1358
									}
									depth--
									add(ruleEXISTS, position1359)
								}
								goto l1357
							l1358:
								position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
								{
									position1372 := position
									depth++
									{
										position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1374
										}
										position++
										goto l1373
									l1374:
										position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
										if buffer[position] != rune('N') {
											goto l774
										}
										position++
									}
								l1373:
									{
										position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1376
										}
										position++
										goto l1375
									l1376:
										position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
										if buffer[position] != rune('O') {
											goto l774
										}
										position++
									}
								l1375:
									{
										position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1378
										}
										position++
										goto l1377
									l1378:
										position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
										if buffer[position] != rune('T') {
											goto l774
										}
										position++
									}
								l1377:
									if buffer[position] != rune(' ') {
										goto l774
									}
									position++
									{
										position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1380
										}
										position++
										goto l1379
									l1380:
										position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
										if buffer[position] != rune('E') {
											goto l774
										}
										position++
									}
								l1379:
									{
										position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1382
										}
										position++
										goto l1381
									l1382:
										position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
										if buffer[position] != rune('X') {
											goto l774
										}
										position++
									}
								l1381:
									{
										position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1384
										}
										position++
										goto l1383
									l1384:
										position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
										if buffer[position] != rune('I') {
											goto l774
										}
										position++
									}
								l1383:
									{
										position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1386
										}
										position++
										goto l1385
									l1386:
										position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
										if buffer[position] != rune('S') {
											goto l774
										}
										position++
									}
								l1385:
									{
										position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1388
										}
										position++
										goto l1387
									l1388:
										position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
										if buffer[position] != rune('T') {
											goto l774
										}
										position++
									}
								l1387:
									{
										position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1390
										}
										position++
										goto l1389
									l1390:
										position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
										if buffer[position] != rune('S') {
											goto l774
										}
										position++
									}
								l1389:
									if !rules[ruleskip]() {
										goto l774
									}
									depth--
									add(ruleNOTEXIST, position1372)
								}
							}
						l1357:
							if !rules[rulegroupGraphPattern]() {
								goto l774
							}
							break
						case 'I', 'i':
							{
								position1391 := position
								depth++
								{
									position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1393
									}
									position++
									goto l1392
								l1393:
									position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
									if buffer[position] != rune('I') {
										goto l774
									}
									position++
								}
							l1392:
								{
									position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1395
									}
									position++
									goto l1394
								l1395:
									position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
									if buffer[position] != rune('F') {
										goto l774
									}
									position++
								}
							l1394:
								if !rules[ruleskip]() {
									goto l774
								}
								depth--
								add(ruleIF, position1391)
							}
							if !rules[ruleLPAREN]() {
								goto l774
							}
							if !rules[ruleexpression]() {
								goto l774
							}
							if !rules[ruleCOMMA]() {
								goto l774
							}
							if !rules[ruleexpression]() {
								goto l774
							}
							if !rules[ruleCOMMA]() {
								goto l774
							}
							if !rules[ruleexpression]() {
								goto l774
							}
							if !rules[ruleRPAREN]() {
								goto l774
							}
							break
						case 'C', 'c':
							{
								position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
								{
									position1398 := position
									depth++
									{
										position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1400
										}
										position++
										goto l1399
									l1400:
										position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
										if buffer[position] != rune('C') {
											goto l1397
										}
										position++
									}
								l1399:
									{
										position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1402
										}
										position++
										goto l1401
									l1402:
										position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
										if buffer[position] != rune('O') {
											goto l1397
										}
										position++
									}
								l1401:
									{
										position1403, tokenIndex1403, depth1403 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1404
										}
										position++
										goto l1403
									l1404:
										position, tokenIndex, depth = position1403, tokenIndex1403, depth1403
										if buffer[position] != rune('N') {
											goto l1397
										}
										position++
									}
								l1403:
									{
										position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1406
										}
										position++
										goto l1405
									l1406:
										position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
										if buffer[position] != rune('C') {
											goto l1397
										}
										position++
									}
								l1405:
									{
										position1407, tokenIndex1407, depth1407 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1408
										}
										position++
										goto l1407
									l1408:
										position, tokenIndex, depth = position1407, tokenIndex1407, depth1407
										if buffer[position] != rune('A') {
											goto l1397
										}
										position++
									}
								l1407:
									{
										position1409, tokenIndex1409, depth1409 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1410
										}
										position++
										goto l1409
									l1410:
										position, tokenIndex, depth = position1409, tokenIndex1409, depth1409
										if buffer[position] != rune('T') {
											goto l1397
										}
										position++
									}
								l1409:
									if !rules[ruleskip]() {
										goto l1397
									}
									depth--
									add(ruleCONCAT, position1398)
								}
								goto l1396
							l1397:
								position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
								{
									position1411 := position
									depth++
									{
										position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1413
										}
										position++
										goto l1412
									l1413:
										position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
										if buffer[position] != rune('C') {
											goto l774
										}
										position++
									}
								l1412:
									{
										position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1415
										}
										position++
										goto l1414
									l1415:
										position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
										if buffer[position] != rune('O') {
											goto l774
										}
										position++
									}
								l1414:
									{
										position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1417
										}
										position++
										goto l1416
									l1417:
										position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
										if buffer[position] != rune('A') {
											goto l774
										}
										position++
									}
								l1416:
									{
										position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1419
										}
										position++
										goto l1418
									l1419:
										position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
										if buffer[position] != rune('L') {
											goto l774
										}
										position++
									}
								l1418:
									{
										position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1421
										}
										position++
										goto l1420
									l1421:
										position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
										if buffer[position] != rune('E') {
											goto l774
										}
										position++
									}
								l1420:
									{
										position1422, tokenIndex1422, depth1422 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1423
										}
										position++
										goto l1422
									l1423:
										position, tokenIndex, depth = position1422, tokenIndex1422, depth1422
										if buffer[position] != rune('S') {
											goto l774
										}
										position++
									}
								l1422:
									{
										position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1425
										}
										position++
										goto l1424
									l1425:
										position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
										if buffer[position] != rune('C') {
											goto l774
										}
										position++
									}
								l1424:
									{
										position1426, tokenIndex1426, depth1426 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1427
										}
										position++
										goto l1426
									l1427:
										position, tokenIndex, depth = position1426, tokenIndex1426, depth1426
										if buffer[position] != rune('E') {
											goto l774
										}
										position++
									}
								l1426:
									if !rules[ruleskip]() {
										goto l774
									}
									depth--
									add(ruleCOALESCE, position1411)
								}
							}
						l1396:
							if !rules[ruleargList]() {
								goto l774
							}
							break
						case 'B', 'b':
							{
								position1428 := position
								depth++
								{
									position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1430
									}
									position++
									goto l1429
								l1430:
									position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
									if buffer[position] != rune('B') {
										goto l774
									}
									position++
								}
							l1429:
								{
									position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1432
									}
									position++
									goto l1431
								l1432:
									position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
									if buffer[position] != rune('N') {
										goto l774
									}
									position++
								}
							l1431:
								{
									position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1434
									}
									position++
									goto l1433
								l1434:
									position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
									if buffer[position] != rune('O') {
										goto l774
									}
									position++
								}
							l1433:
								{
									position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1436
									}
									position++
									goto l1435
								l1436:
									position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
									if buffer[position] != rune('D') {
										goto l774
									}
									position++
								}
							l1435:
								{
									position1437, tokenIndex1437, depth1437 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1438
									}
									position++
									goto l1437
								l1438:
									position, tokenIndex, depth = position1437, tokenIndex1437, depth1437
									if buffer[position] != rune('E') {
										goto l774
									}
									position++
								}
							l1437:
								if !rules[ruleskip]() {
									goto l774
								}
								depth--
								add(ruleBNODE, position1428)
							}
							{
								position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1440
								}
								if !rules[ruleexpression]() {
									goto l1440
								}
								if !rules[ruleRPAREN]() {
									goto l1440
								}
								goto l1439
							l1440:
								position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
								if !rules[rulenil]() {
									goto l774
								}
							}
						l1439:
							break
						default:
							{
								position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
								{
									position1443 := position
									depth++
									{
										position1444, tokenIndex1444, depth1444 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1445
										}
										position++
										goto l1444
									l1445:
										position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
										if buffer[position] != rune('S') {
											goto l1442
										}
										position++
									}
								l1444:
									{
										position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1447
										}
										position++
										goto l1446
									l1447:
										position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
										if buffer[position] != rune('U') {
											goto l1442
										}
										position++
									}
								l1446:
									{
										position1448, tokenIndex1448, depth1448 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1449
										}
										position++
										goto l1448
									l1449:
										position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
										if buffer[position] != rune('B') {
											goto l1442
										}
										position++
									}
								l1448:
									{
										position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1451
										}
										position++
										goto l1450
									l1451:
										position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
										if buffer[position] != rune('S') {
											goto l1442
										}
										position++
									}
								l1450:
									{
										position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1453
										}
										position++
										goto l1452
									l1453:
										position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
										if buffer[position] != rune('T') {
											goto l1442
										}
										position++
									}
								l1452:
									{
										position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1455
										}
										position++
										goto l1454
									l1455:
										position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
										if buffer[position] != rune('R') {
											goto l1442
										}
										position++
									}
								l1454:
									if !rules[ruleskip]() {
										goto l1442
									}
									depth--
									add(ruleSUBSTR, position1443)
								}
								goto l1441
							l1442:
								position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
								{
									position1457 := position
									depth++
									{
										position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1459
										}
										position++
										goto l1458
									l1459:
										position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
										if buffer[position] != rune('R') {
											goto l1456
										}
										position++
									}
								l1458:
									{
										position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1461
										}
										position++
										goto l1460
									l1461:
										position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
										if buffer[position] != rune('E') {
											goto l1456
										}
										position++
									}
								l1460:
									{
										position1462, tokenIndex1462, depth1462 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1463
										}
										position++
										goto l1462
									l1463:
										position, tokenIndex, depth = position1462, tokenIndex1462, depth1462
										if buffer[position] != rune('P') {
											goto l1456
										}
										position++
									}
								l1462:
									{
										position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1465
										}
										position++
										goto l1464
									l1465:
										position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
										if buffer[position] != rune('L') {
											goto l1456
										}
										position++
									}
								l1464:
									{
										position1466, tokenIndex1466, depth1466 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1467
										}
										position++
										goto l1466
									l1467:
										position, tokenIndex, depth = position1466, tokenIndex1466, depth1466
										if buffer[position] != rune('A') {
											goto l1456
										}
										position++
									}
								l1466:
									{
										position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1469
										}
										position++
										goto l1468
									l1469:
										position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
										if buffer[position] != rune('C') {
											goto l1456
										}
										position++
									}
								l1468:
									{
										position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1471
										}
										position++
										goto l1470
									l1471:
										position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
										if buffer[position] != rune('E') {
											goto l1456
										}
										position++
									}
								l1470:
									if !rules[ruleskip]() {
										goto l1456
									}
									depth--
									add(ruleREPLACE, position1457)
								}
								goto l1441
							l1456:
								position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
								{
									position1472 := position
									depth++
									{
										position1473, tokenIndex1473, depth1473 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1474
										}
										position++
										goto l1473
									l1474:
										position, tokenIndex, depth = position1473, tokenIndex1473, depth1473
										if buffer[position] != rune('R') {
											goto l774
										}
										position++
									}
								l1473:
									{
										position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1476
										}
										position++
										goto l1475
									l1476:
										position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
										if buffer[position] != rune('E') {
											goto l774
										}
										position++
									}
								l1475:
									{
										position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1478
										}
										position++
										goto l1477
									l1478:
										position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
										if buffer[position] != rune('G') {
											goto l774
										}
										position++
									}
								l1477:
									{
										position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1480
										}
										position++
										goto l1479
									l1480:
										position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
										if buffer[position] != rune('E') {
											goto l774
										}
										position++
									}
								l1479:
									{
										position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1482
										}
										position++
										goto l1481
									l1482:
										position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
										if buffer[position] != rune('X') {
											goto l774
										}
										position++
									}
								l1481:
									if !rules[ruleskip]() {
										goto l774
									}
									depth--
									add(ruleREGEX, position1472)
								}
							}
						l1441:
							if !rules[ruleLPAREN]() {
								goto l774
							}
							if !rules[ruleexpression]() {
								goto l774
							}
							if !rules[ruleCOMMA]() {
								goto l774
							}
							if !rules[ruleexpression]() {
								goto l774
							}
							{
								position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1483
								}
								if !rules[ruleexpression]() {
									goto l1483
								}
								goto l1484
							l1483:
								position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
							}
						l1484:
							if !rules[ruleRPAREN]() {
								goto l774
							}
							break
						}
					}

				}
			l776:
				depth--
				add(rulebuiltinCall, position775)
			}
			return true
		l774:
			position, tokenIndex, depth = position774, tokenIndex774, depth774
			return false
		},
		/* 66 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1485, tokenIndex1485, depth1485 := position, tokenIndex, depth
			{
				position1486 := position
				depth++
				{
					position1487, tokenIndex1487, depth1487 := position, tokenIndex, depth
					{
						position1489 := position
						depth++
					l1490:
						{
							position1491, tokenIndex1491, depth1491 := position, tokenIndex, depth
							{
								position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1493
								}
								position++
								goto l1492
							l1493:
								position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1491
								}
								position++
							}
						l1492:
							goto l1490
						l1491:
							position, tokenIndex, depth = position1491, tokenIndex1491, depth1491
						}
						depth--
						add(rulePegText, position1489)
					}
					if buffer[position] != rune(':') {
						goto l1488
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1487
				l1488:
					position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
					{
						position1496 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1495
						}
						position++
					l1497:
						{
							position1498, tokenIndex1498, depth1498 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1498
							}
							position++
							goto l1497
						l1498:
							position, tokenIndex, depth = position1498, tokenIndex1498, depth1498
						}
						depth--
						add(rulePegText, position1496)
					}
					if buffer[position] != rune('/') {
						goto l1495
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1487
				l1495:
					position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
					{
						position1500 := position
						depth++
					l1501:
						{
							position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1502
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1502
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1502
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1502
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1502
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1502
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1502
									}
									position++
									break
								}
							}

							goto l1501
						l1502:
							position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
						}
						depth--
						add(rulePegText, position1500)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1487:
				if buffer[position] != rune('<') {
					goto l1485
				}
				position++
				if !rules[rulews]() {
					goto l1485
				}
				if !rules[ruleskip]() {
					goto l1485
				}
				depth--
				add(rulepof, position1486)
			}
			return true
		l1485:
			position, tokenIndex, depth = position1485, tokenIndex1485, depth1485
			return false
		},
		/* 67 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
			{
				position1506 := position
				depth++
				{
					position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1508
					}
					position++
					goto l1507
				l1508:
					position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
					if buffer[position] != rune('$') {
						goto l1505
					}
					position++
				}
			l1507:
				{
					position1509 := position
					depth++
					{
						position1512, tokenIndex1512, depth1512 := position, tokenIndex, depth
						{
							position1514 := position
							depth++
							{
								position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
								{
									position1517 := position
									depth++
									{
										position1518, tokenIndex1518, depth1518 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1519
										}
										position++
										goto l1518
									l1519:
										position, tokenIndex, depth = position1518, tokenIndex1518, depth1518
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1516
										}
										position++
									}
								l1518:
									depth--
									add(rulePN_CHARS_BASE, position1517)
								}
								goto l1515
							l1516:
								position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
								if buffer[position] != rune('_') {
									goto l1513
								}
								position++
							}
						l1515:
							depth--
							add(rulePN_CHARS_U, position1514)
						}
						goto l1512
					l1513:
						position, tokenIndex, depth = position1512, tokenIndex1512, depth1512
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1505
						}
						position++
					}
				l1512:
				l1510:
					{
						position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
						{
							position1520, tokenIndex1520, depth1520 := position, tokenIndex, depth
							{
								position1522 := position
								depth++
								{
									position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
									{
										position1525 := position
										depth++
										{
											position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1527
											}
											position++
											goto l1526
										l1527:
											position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1524
											}
											position++
										}
									l1526:
										depth--
										add(rulePN_CHARS_BASE, position1525)
									}
									goto l1523
								l1524:
									position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
									if buffer[position] != rune('_') {
										goto l1521
									}
									position++
								}
							l1523:
								depth--
								add(rulePN_CHARS_U, position1522)
							}
							goto l1520
						l1521:
							position, tokenIndex, depth = position1520, tokenIndex1520, depth1520
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1511
							}
							position++
						}
					l1520:
						goto l1510
					l1511:
						position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
					}
					depth--
					add(ruleVARNAME, position1509)
				}
				if !rules[ruleskip]() {
					goto l1505
				}
				depth--
				add(rulevar, position1506)
			}
			return true
		l1505:
			position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
			return false
		},
		/* 68 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1528, tokenIndex1528, depth1528 := position, tokenIndex, depth
			{
				position1529 := position
				depth++
				{
					position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1531
					}
					goto l1530
				l1531:
					position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
					{
						position1532 := position
						depth++
					l1533:
						{
							position1534, tokenIndex1534, depth1534 := position, tokenIndex, depth
							{
								position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
								{
									position1536, tokenIndex1536, depth1536 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1537
									}
									position++
									goto l1536
								l1537:
									position, tokenIndex, depth = position1536, tokenIndex1536, depth1536
									if buffer[position] != rune(' ') {
										goto l1535
									}
									position++
								}
							l1536:
								goto l1534
							l1535:
								position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
							}
							if !matchDot() {
								goto l1534
							}
							goto l1533
						l1534:
							position, tokenIndex, depth = position1534, tokenIndex1534, depth1534
						}
						if buffer[position] != rune(':') {
							goto l1528
						}
						position++
					l1538:
						{
							position1539, tokenIndex1539, depth1539 := position, tokenIndex, depth
							{
								position1540, tokenIndex1540, depth1540 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1541
								}
								position++
								goto l1540
							l1541:
								position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1542
								}
								position++
								goto l1540
							l1542:
								position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1543
								}
								position++
								goto l1540
							l1543:
								position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1539
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1539
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1539
										}
										position++
										break
									}
								}

							}
						l1540:
							goto l1538
						l1539:
							position, tokenIndex, depth = position1539, tokenIndex1539, depth1539
						}
						if !rules[ruleskip]() {
							goto l1528
						}
						depth--
						add(ruleprefixedName, position1532)
					}
				}
			l1530:
				depth--
				add(ruleiriref, position1529)
			}
			return true
		l1528:
			position, tokenIndex, depth = position1528, tokenIndex1528, depth1528
			return false
		},
		/* 69 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1545, tokenIndex1545, depth1545 := position, tokenIndex, depth
			{
				position1546 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1545
				}
				position++
			l1547:
				{
					position1548, tokenIndex1548, depth1548 := position, tokenIndex, depth
					{
						position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1549
						}
						position++
						goto l1548
					l1549:
						position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
					}
					if !matchDot() {
						goto l1548
					}
					goto l1547
				l1548:
					position, tokenIndex, depth = position1548, tokenIndex1548, depth1548
				}
				if buffer[position] != rune('>') {
					goto l1545
				}
				position++
				if !rules[ruleskip]() {
					goto l1545
				}
				depth--
				add(ruleiri, position1546)
			}
			return true
		l1545:
			position, tokenIndex, depth = position1545, tokenIndex1545, depth1545
			return false
		},
		/* 70 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 71 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1551, tokenIndex1551, depth1551 := position, tokenIndex, depth
			{
				position1552 := position
				depth++
				if !rules[rulestring]() {
					goto l1551
				}
				{
					position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
					{
						position1555, tokenIndex1555, depth1555 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1556
						}
						position++
						{
							position1559, tokenIndex1559, depth1559 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1560
							}
							position++
							goto l1559
						l1560:
							position, tokenIndex, depth = position1559, tokenIndex1559, depth1559
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1556
							}
							position++
						}
					l1559:
					l1557:
						{
							position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
							{
								position1561, tokenIndex1561, depth1561 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1562
								}
								position++
								goto l1561
							l1562:
								position, tokenIndex, depth = position1561, tokenIndex1561, depth1561
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1558
								}
								position++
							}
						l1561:
							goto l1557
						l1558:
							position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
						}
					l1563:
						{
							position1564, tokenIndex1564, depth1564 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1564
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1564
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1564
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1564
									}
									position++
									break
								}
							}

						l1565:
							{
								position1566, tokenIndex1566, depth1566 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1566
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1566
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1566
										}
										position++
										break
									}
								}

								goto l1565
							l1566:
								position, tokenIndex, depth = position1566, tokenIndex1566, depth1566
							}
							goto l1563
						l1564:
							position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
						}
						goto l1555
					l1556:
						position, tokenIndex, depth = position1555, tokenIndex1555, depth1555
						if buffer[position] != rune('^') {
							goto l1553
						}
						position++
						if buffer[position] != rune('^') {
							goto l1553
						}
						position++
						if !rules[ruleiriref]() {
							goto l1553
						}
					}
				l1555:
					goto l1554
				l1553:
					position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
				}
			l1554:
				if !rules[ruleskip]() {
					goto l1551
				}
				depth--
				add(ruleliteral, position1552)
			}
			return true
		l1551:
			position, tokenIndex, depth = position1551, tokenIndex1551, depth1551
			return false
		},
		/* 72 string <- <('"' (!'"' .)* '"')> */
		func() bool {
			position1569, tokenIndex1569, depth1569 := position, tokenIndex, depth
			{
				position1570 := position
				depth++
				if buffer[position] != rune('"') {
					goto l1569
				}
				position++
			l1571:
				{
					position1572, tokenIndex1572, depth1572 := position, tokenIndex, depth
					{
						position1573, tokenIndex1573, depth1573 := position, tokenIndex, depth
						if buffer[position] != rune('"') {
							goto l1573
						}
						position++
						goto l1572
					l1573:
						position, tokenIndex, depth = position1573, tokenIndex1573, depth1573
					}
					if !matchDot() {
						goto l1572
					}
					goto l1571
				l1572:
					position, tokenIndex, depth = position1572, tokenIndex1572, depth1572
				}
				if buffer[position] != rune('"') {
					goto l1569
				}
				position++
				depth--
				add(rulestring, position1570)
			}
			return true
		l1569:
			position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
			return false
		},
		/* 73 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1574, tokenIndex1574, depth1574 := position, tokenIndex, depth
			{
				position1575 := position
				depth++
				{
					position1576, tokenIndex1576, depth1576 := position, tokenIndex, depth
					{
						position1578, tokenIndex1578, depth1578 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1579
						}
						position++
						goto l1578
					l1579:
						position, tokenIndex, depth = position1578, tokenIndex1578, depth1578
						if buffer[position] != rune('-') {
							goto l1576
						}
						position++
					}
				l1578:
					goto l1577
				l1576:
					position, tokenIndex, depth = position1576, tokenIndex1576, depth1576
				}
			l1577:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1574
				}
				position++
			l1580:
				{
					position1581, tokenIndex1581, depth1581 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1581
					}
					position++
					goto l1580
				l1581:
					position, tokenIndex, depth = position1581, tokenIndex1581, depth1581
				}
				{
					position1582, tokenIndex1582, depth1582 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1582
					}
					position++
				l1584:
					{
						position1585, tokenIndex1585, depth1585 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1585
						}
						position++
						goto l1584
					l1585:
						position, tokenIndex, depth = position1585, tokenIndex1585, depth1585
					}
					goto l1583
				l1582:
					position, tokenIndex, depth = position1582, tokenIndex1582, depth1582
				}
			l1583:
				if !rules[ruleskip]() {
					goto l1574
				}
				depth--
				add(rulenumericLiteral, position1575)
			}
			return true
		l1574:
			position, tokenIndex, depth = position1574, tokenIndex1574, depth1574
			return false
		},
		/* 74 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 75 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1587, tokenIndex1587, depth1587 := position, tokenIndex, depth
			{
				position1588 := position
				depth++
				{
					position1589, tokenIndex1589, depth1589 := position, tokenIndex, depth
					{
						position1591 := position
						depth++
						{
							position1592, tokenIndex1592, depth1592 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1593
							}
							position++
							goto l1592
						l1593:
							position, tokenIndex, depth = position1592, tokenIndex1592, depth1592
							if buffer[position] != rune('T') {
								goto l1590
							}
							position++
						}
					l1592:
						{
							position1594, tokenIndex1594, depth1594 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1595
							}
							position++
							goto l1594
						l1595:
							position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
							if buffer[position] != rune('R') {
								goto l1590
							}
							position++
						}
					l1594:
						{
							position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1597
							}
							position++
							goto l1596
						l1597:
							position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
							if buffer[position] != rune('U') {
								goto l1590
							}
							position++
						}
					l1596:
						{
							position1598, tokenIndex1598, depth1598 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1599
							}
							position++
							goto l1598
						l1599:
							position, tokenIndex, depth = position1598, tokenIndex1598, depth1598
							if buffer[position] != rune('E') {
								goto l1590
							}
							position++
						}
					l1598:
						if !rules[ruleskip]() {
							goto l1590
						}
						depth--
						add(ruleTRUE, position1591)
					}
					goto l1589
				l1590:
					position, tokenIndex, depth = position1589, tokenIndex1589, depth1589
					{
						position1600 := position
						depth++
						{
							position1601, tokenIndex1601, depth1601 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1602
							}
							position++
							goto l1601
						l1602:
							position, tokenIndex, depth = position1601, tokenIndex1601, depth1601
							if buffer[position] != rune('F') {
								goto l1587
							}
							position++
						}
					l1601:
						{
							position1603, tokenIndex1603, depth1603 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1604
							}
							position++
							goto l1603
						l1604:
							position, tokenIndex, depth = position1603, tokenIndex1603, depth1603
							if buffer[position] != rune('A') {
								goto l1587
							}
							position++
						}
					l1603:
						{
							position1605, tokenIndex1605, depth1605 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1606
							}
							position++
							goto l1605
						l1606:
							position, tokenIndex, depth = position1605, tokenIndex1605, depth1605
							if buffer[position] != rune('L') {
								goto l1587
							}
							position++
						}
					l1605:
						{
							position1607, tokenIndex1607, depth1607 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1608
							}
							position++
							goto l1607
						l1608:
							position, tokenIndex, depth = position1607, tokenIndex1607, depth1607
							if buffer[position] != rune('S') {
								goto l1587
							}
							position++
						}
					l1607:
						{
							position1609, tokenIndex1609, depth1609 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1610
							}
							position++
							goto l1609
						l1610:
							position, tokenIndex, depth = position1609, tokenIndex1609, depth1609
							if buffer[position] != rune('E') {
								goto l1587
							}
							position++
						}
					l1609:
						if !rules[ruleskip]() {
							goto l1587
						}
						depth--
						add(ruleFALSE, position1600)
					}
				}
			l1589:
				depth--
				add(rulebooleanLiteral, position1588)
			}
			return true
		l1587:
			position, tokenIndex, depth = position1587, tokenIndex1587, depth1587
			return false
		},
		/* 76 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 77 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 78 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 79 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1614, tokenIndex1614, depth1614 := position, tokenIndex, depth
			{
				position1615 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1614
				}
				position++
			l1616:
				{
					position1617, tokenIndex1617, depth1617 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1617
					}
					goto l1616
				l1617:
					position, tokenIndex, depth = position1617, tokenIndex1617, depth1617
				}
				if buffer[position] != rune(')') {
					goto l1614
				}
				position++
				if !rules[ruleskip]() {
					goto l1614
				}
				depth--
				add(rulenil, position1615)
			}
			return true
		l1614:
			position, tokenIndex, depth = position1614, tokenIndex1614, depth1614
			return false
		},
		/* 80 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 81 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 82 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 83 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 84 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 85 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 86 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 87 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 88 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 89 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
			{
				position1628 := position
				depth++
				{
					position1629, tokenIndex1629, depth1629 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1630
					}
					position++
					goto l1629
				l1630:
					position, tokenIndex, depth = position1629, tokenIndex1629, depth1629
					if buffer[position] != rune('D') {
						goto l1627
					}
					position++
				}
			l1629:
				{
					position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1632
					}
					position++
					goto l1631
				l1632:
					position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
					if buffer[position] != rune('I') {
						goto l1627
					}
					position++
				}
			l1631:
				{
					position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1634
					}
					position++
					goto l1633
				l1634:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					if buffer[position] != rune('S') {
						goto l1627
					}
					position++
				}
			l1633:
				{
					position1635, tokenIndex1635, depth1635 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1636
					}
					position++
					goto l1635
				l1636:
					position, tokenIndex, depth = position1635, tokenIndex1635, depth1635
					if buffer[position] != rune('T') {
						goto l1627
					}
					position++
				}
			l1635:
				{
					position1637, tokenIndex1637, depth1637 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1638
					}
					position++
					goto l1637
				l1638:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if buffer[position] != rune('I') {
						goto l1627
					}
					position++
				}
			l1637:
				{
					position1639, tokenIndex1639, depth1639 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1640
					}
					position++
					goto l1639
				l1640:
					position, tokenIndex, depth = position1639, tokenIndex1639, depth1639
					if buffer[position] != rune('N') {
						goto l1627
					}
					position++
				}
			l1639:
				{
					position1641, tokenIndex1641, depth1641 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1642
					}
					position++
					goto l1641
				l1642:
					position, tokenIndex, depth = position1641, tokenIndex1641, depth1641
					if buffer[position] != rune('C') {
						goto l1627
					}
					position++
				}
			l1641:
				{
					position1643, tokenIndex1643, depth1643 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1644
					}
					position++
					goto l1643
				l1644:
					position, tokenIndex, depth = position1643, tokenIndex1643, depth1643
					if buffer[position] != rune('T') {
						goto l1627
					}
					position++
				}
			l1643:
				if !rules[ruleskip]() {
					goto l1627
				}
				depth--
				add(ruleDISTINCT, position1628)
			}
			return true
		l1627:
			position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
			return false
		},
		/* 90 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 91 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 92 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 93 LBRACE <- <('{' skip)> */
		func() bool {
			position1648, tokenIndex1648, depth1648 := position, tokenIndex, depth
			{
				position1649 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1648
				}
				position++
				if !rules[ruleskip]() {
					goto l1648
				}
				depth--
				add(ruleLBRACE, position1649)
			}
			return true
		l1648:
			position, tokenIndex, depth = position1648, tokenIndex1648, depth1648
			return false
		},
		/* 94 RBRACE <- <('}' skip)> */
		func() bool {
			position1650, tokenIndex1650, depth1650 := position, tokenIndex, depth
			{
				position1651 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1650
				}
				position++
				if !rules[ruleskip]() {
					goto l1650
				}
				depth--
				add(ruleRBRACE, position1651)
			}
			return true
		l1650:
			position, tokenIndex, depth = position1650, tokenIndex1650, depth1650
			return false
		},
		/* 95 LBRACK <- <('[' skip)> */
		nil,
		/* 96 RBRACK <- <(']' skip)> */
		nil,
		/* 97 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1654, tokenIndex1654, depth1654 := position, tokenIndex, depth
			{
				position1655 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1654
				}
				position++
				if !rules[ruleskip]() {
					goto l1654
				}
				depth--
				add(ruleSEMICOLON, position1655)
			}
			return true
		l1654:
			position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
			return false
		},
		/* 98 COMMA <- <(',' skip)> */
		func() bool {
			position1656, tokenIndex1656, depth1656 := position, tokenIndex, depth
			{
				position1657 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1656
				}
				position++
				if !rules[ruleskip]() {
					goto l1656
				}
				depth--
				add(ruleCOMMA, position1657)
			}
			return true
		l1656:
			position, tokenIndex, depth = position1656, tokenIndex1656, depth1656
			return false
		},
		/* 99 DOT <- <('.' skip)> */
		func() bool {
			position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
			{
				position1659 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1658
				}
				position++
				if !rules[ruleskip]() {
					goto l1658
				}
				depth--
				add(ruleDOT, position1659)
			}
			return true
		l1658:
			position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
			return false
		},
		/* 100 COLON <- <(':' skip)> */
		nil,
		/* 101 PIPE <- <('|' skip)> */
		func() bool {
			position1661, tokenIndex1661, depth1661 := position, tokenIndex, depth
			{
				position1662 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1661
				}
				position++
				if !rules[ruleskip]() {
					goto l1661
				}
				depth--
				add(rulePIPE, position1662)
			}
			return true
		l1661:
			position, tokenIndex, depth = position1661, tokenIndex1661, depth1661
			return false
		},
		/* 102 SLASH <- <('/' skip)> */
		func() bool {
			position1663, tokenIndex1663, depth1663 := position, tokenIndex, depth
			{
				position1664 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1663
				}
				position++
				if !rules[ruleskip]() {
					goto l1663
				}
				depth--
				add(ruleSLASH, position1664)
			}
			return true
		l1663:
			position, tokenIndex, depth = position1663, tokenIndex1663, depth1663
			return false
		},
		/* 103 INVERSE <- <('^' skip)> */
		func() bool {
			position1665, tokenIndex1665, depth1665 := position, tokenIndex, depth
			{
				position1666 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1665
				}
				position++
				if !rules[ruleskip]() {
					goto l1665
				}
				depth--
				add(ruleINVERSE, position1666)
			}
			return true
		l1665:
			position, tokenIndex, depth = position1665, tokenIndex1665, depth1665
			return false
		},
		/* 104 LPAREN <- <('(' skip)> */
		func() bool {
			position1667, tokenIndex1667, depth1667 := position, tokenIndex, depth
			{
				position1668 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1667
				}
				position++
				if !rules[ruleskip]() {
					goto l1667
				}
				depth--
				add(ruleLPAREN, position1668)
			}
			return true
		l1667:
			position, tokenIndex, depth = position1667, tokenIndex1667, depth1667
			return false
		},
		/* 105 RPAREN <- <(')' skip)> */
		func() bool {
			position1669, tokenIndex1669, depth1669 := position, tokenIndex, depth
			{
				position1670 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1669
				}
				position++
				if !rules[ruleskip]() {
					goto l1669
				}
				depth--
				add(ruleRPAREN, position1670)
			}
			return true
		l1669:
			position, tokenIndex, depth = position1669, tokenIndex1669, depth1669
			return false
		},
		/* 106 ISA <- <('a' skip)> */
		func() bool {
			position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
			{
				position1672 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1671
				}
				position++
				if !rules[ruleskip]() {
					goto l1671
				}
				depth--
				add(ruleISA, position1672)
			}
			return true
		l1671:
			position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
			return false
		},
		/* 107 NOT <- <('!' skip)> */
		func() bool {
			position1673, tokenIndex1673, depth1673 := position, tokenIndex, depth
			{
				position1674 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1673
				}
				position++
				if !rules[ruleskip]() {
					goto l1673
				}
				depth--
				add(ruleNOT, position1674)
			}
			return true
		l1673:
			position, tokenIndex, depth = position1673, tokenIndex1673, depth1673
			return false
		},
		/* 108 STAR <- <('*' skip)> */
		func() bool {
			position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
			{
				position1676 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1675
				}
				position++
				if !rules[ruleskip]() {
					goto l1675
				}
				depth--
				add(ruleSTAR, position1676)
			}
			return true
		l1675:
			position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
			return false
		},
		/* 109 PLUS <- <('+' skip)> */
		func() bool {
			position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
			{
				position1678 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1677
				}
				position++
				if !rules[ruleskip]() {
					goto l1677
				}
				depth--
				add(rulePLUS, position1678)
			}
			return true
		l1677:
			position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
			return false
		},
		/* 110 MINUS <- <('-' skip)> */
		func() bool {
			position1679, tokenIndex1679, depth1679 := position, tokenIndex, depth
			{
				position1680 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1679
				}
				position++
				if !rules[ruleskip]() {
					goto l1679
				}
				depth--
				add(ruleMINUS, position1680)
			}
			return true
		l1679:
			position, tokenIndex, depth = position1679, tokenIndex1679, depth1679
			return false
		},
		/* 111 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 112 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 113 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 114 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 115 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1685, tokenIndex1685, depth1685 := position, tokenIndex, depth
			{
				position1686 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1685
				}
				position++
			l1687:
				{
					position1688, tokenIndex1688, depth1688 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1688
					}
					position++
					goto l1687
				l1688:
					position, tokenIndex, depth = position1688, tokenIndex1688, depth1688
				}
				if !rules[ruleskip]() {
					goto l1685
				}
				depth--
				add(ruleINTEGER, position1686)
			}
			return true
		l1685:
			position, tokenIndex, depth = position1685, tokenIndex1685, depth1685
			return false
		},
		/* 116 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 117 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 118 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 119 OR <- <('|' '|' skip)> */
		nil,
		/* 120 AND <- <('&' '&' skip)> */
		nil,
		/* 121 EQ <- <('=' skip)> */
		func() bool {
			position1694, tokenIndex1694, depth1694 := position, tokenIndex, depth
			{
				position1695 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1694
				}
				position++
				if !rules[ruleskip]() {
					goto l1694
				}
				depth--
				add(ruleEQ, position1695)
			}
			return true
		l1694:
			position, tokenIndex, depth = position1694, tokenIndex1694, depth1694
			return false
		},
		/* 122 NE <- <('!' '=' skip)> */
		nil,
		/* 123 GT <- <('>' skip)> */
		nil,
		/* 124 LT <- <('<' skip)> */
		nil,
		/* 125 LE <- <('<' '=' skip)> */
		nil,
		/* 126 GE <- <('>' '=' skip)> */
		nil,
		/* 127 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 128 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 129 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1703, tokenIndex1703, depth1703 := position, tokenIndex, depth
			{
				position1704 := position
				depth++
				{
					position1705, tokenIndex1705, depth1705 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1706
					}
					position++
					goto l1705
				l1706:
					position, tokenIndex, depth = position1705, tokenIndex1705, depth1705
					if buffer[position] != rune('A') {
						goto l1703
					}
					position++
				}
			l1705:
				{
					position1707, tokenIndex1707, depth1707 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1708
					}
					position++
					goto l1707
				l1708:
					position, tokenIndex, depth = position1707, tokenIndex1707, depth1707
					if buffer[position] != rune('S') {
						goto l1703
					}
					position++
				}
			l1707:
				if !rules[ruleskip]() {
					goto l1703
				}
				depth--
				add(ruleAS, position1704)
			}
			return true
		l1703:
			position, tokenIndex, depth = position1703, tokenIndex1703, depth1703
			return false
		},
		/* 130 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 131 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 132 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 133 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 134 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 135 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 136 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 137 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 138 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 139 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 140 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 141 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 142 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 143 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 144 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 145 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 146 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 147 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 148 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 149 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 150 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 151 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 152 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 153 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 154 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 155 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 156 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 157 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 158 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 159 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 160 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 161 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 162 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 163 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 164 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 165 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 166 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 167 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 168 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 169 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 170 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 171 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 172 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 173 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 174 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 175 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 176 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 177 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 178 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 179 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 180 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 181 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 182 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 183 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 184 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 185 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 186 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 187 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 188 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 189 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 190 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 191 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 192 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 193 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 194 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 195 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 196 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 197 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 198 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1777, tokenIndex1777, depth1777 := position, tokenIndex, depth
			{
				position1778 := position
				depth++
				{
					position1779, tokenIndex1779, depth1779 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1780
					}
					position++
					goto l1779
				l1780:
					position, tokenIndex, depth = position1779, tokenIndex1779, depth1779
					if buffer[position] != rune('B') {
						goto l1777
					}
					position++
				}
			l1779:
				{
					position1781, tokenIndex1781, depth1781 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1782
					}
					position++
					goto l1781
				l1782:
					position, tokenIndex, depth = position1781, tokenIndex1781, depth1781
					if buffer[position] != rune('Y') {
						goto l1777
					}
					position++
				}
			l1781:
				if !rules[ruleskip]() {
					goto l1777
				}
				depth--
				add(ruleBY, position1778)
			}
			return true
		l1777:
			position, tokenIndex, depth = position1777, tokenIndex1777, depth1777
			return false
		},
		/* 199 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 200 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1785 := position
				depth++
			l1786:
				{
					position1787, tokenIndex1787, depth1787 := position, tokenIndex, depth
					{
						position1788, tokenIndex1788, depth1788 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1789
						}
						goto l1788
					l1789:
						position, tokenIndex, depth = position1788, tokenIndex1788, depth1788
						{
							position1790 := position
							depth++
							{
								position1791 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1787
								}
								position++
							l1792:
								{
									position1793, tokenIndex1793, depth1793 := position, tokenIndex, depth
									{
										position1794, tokenIndex1794, depth1794 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1794
										}
										goto l1793
									l1794:
										position, tokenIndex, depth = position1794, tokenIndex1794, depth1794
									}
									if !matchDot() {
										goto l1793
									}
									goto l1792
								l1793:
									position, tokenIndex, depth = position1793, tokenIndex1793, depth1793
								}
								if !rules[ruleendOfLine]() {
									goto l1787
								}
								depth--
								add(rulePegText, position1791)
							}
							{
								add(ruleAction13, position)
							}
							depth--
							add(rulecomment, position1790)
						}
					}
				l1788:
					goto l1786
				l1787:
					position, tokenIndex, depth = position1787, tokenIndex1787, depth1787
				}
				depth--
				add(ruleskip, position1785)
			}
			return true
		},
		/* 201 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1796, tokenIndex1796, depth1796 := position, tokenIndex, depth
			{
				position1797 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1796
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1796
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1796
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1796
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1796
						}
						break
					}
				}

				depth--
				add(rulews, position1797)
			}
			return true
		l1796:
			position, tokenIndex, depth = position1796, tokenIndex1796, depth1796
			return false
		},
		/* 202 comment <- <(<('#' (!endOfLine .)* endOfLine)> Action13)> */
		nil,
		/* 203 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1800, tokenIndex1800, depth1800 := position, tokenIndex, depth
			{
				position1801 := position
				depth++
				{
					position1802, tokenIndex1802, depth1802 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1803
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1803
					}
					position++
					goto l1802
				l1803:
					position, tokenIndex, depth = position1802, tokenIndex1802, depth1802
					if buffer[position] != rune('\n') {
						goto l1804
					}
					position++
					goto l1802
				l1804:
					position, tokenIndex, depth = position1802, tokenIndex1802, depth1802
					if buffer[position] != rune('\r') {
						goto l1800
					}
					position++
				}
			l1802:
				depth--
				add(ruleendOfLine, position1801)
			}
			return true
		l1800:
			position, tokenIndex, depth = position1800, tokenIndex1800, depth1800
			return false
		},
		nil,
		/* 206 Action0 <- <{ p.addPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 207 Action1 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 208 Action2 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 209 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 210 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 211 Action5 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 212 Action6 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 213 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 214 Action8 <- <{ p.setObject(p.skipComment(buffer, begin, end)); p.addTriplePattern() }> */
		nil,
		/* 215 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 216 Action10 <- <{ p.setPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 217 Action11 <- <{ p.setPathLength(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 218 Action12 <- <{ p.setKeyword(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 219 Action13 <- <{ p.commentBegin = begin }> */
		nil,
	}
	p.rules = rules
}
