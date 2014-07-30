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
	rules  [204]func() bool
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
						if !rules[ruleDISTINCT]() {
							goto l127
						}
						goto l126
					l127:
						position, tokenIndex, depth = position126, tokenIndex126, depth126
						{
							position128 := position
							depth++
							{
								position129, tokenIndex129, depth129 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l130
								}
								position++
								goto l129
							l130:
								position, tokenIndex, depth = position129, tokenIndex129, depth129
								if buffer[position] != rune('R') {
									goto l124
								}
								position++
							}
						l129:
							{
								position131, tokenIndex131, depth131 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l132
								}
								position++
								goto l131
							l132:
								position, tokenIndex, depth = position131, tokenIndex131, depth131
								if buffer[position] != rune('E') {
									goto l124
								}
								position++
							}
						l131:
							{
								position133, tokenIndex133, depth133 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l134
								}
								position++
								goto l133
							l134:
								position, tokenIndex, depth = position133, tokenIndex133, depth133
								if buffer[position] != rune('D') {
									goto l124
								}
								position++
							}
						l133:
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l136
								}
								position++
								goto l135
							l136:
								position, tokenIndex, depth = position135, tokenIndex135, depth135
								if buffer[position] != rune('U') {
									goto l124
								}
								position++
							}
						l135:
							{
								position137, tokenIndex137, depth137 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l138
								}
								position++
								goto l137
							l138:
								position, tokenIndex, depth = position137, tokenIndex137, depth137
								if buffer[position] != rune('C') {
									goto l124
								}
								position++
							}
						l137:
							{
								position139, tokenIndex139, depth139 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l140
								}
								position++
								goto l139
							l140:
								position, tokenIndex, depth = position139, tokenIndex139, depth139
								if buffer[position] != rune('E') {
									goto l124
								}
								position++
							}
						l139:
							{
								position141, tokenIndex141, depth141 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l142
								}
								position++
								goto l141
							l142:
								position, tokenIndex, depth = position141, tokenIndex141, depth141
								if buffer[position] != rune('D') {
									goto l124
								}
								position++
							}
						l141:
							if !rules[ruleskip]() {
								goto l124
							}
							depth--
							add(ruleREDUCED, position128)
						}
					}
				l126:
					goto l125
				l124:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
				}
			l125:
				{
					position143, tokenIndex143, depth143 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l144
					}
					goto l143
				l144:
					position, tokenIndex, depth = position143, tokenIndex143, depth143
					{
						position147 := position
						depth++
						{
							position148, tokenIndex148, depth148 := position, tokenIndex, depth
							if !rules[rulevar]() {
								goto l149
							}
							goto l148
						l149:
							position, tokenIndex, depth = position148, tokenIndex148, depth148
							if !rules[ruleLPAREN]() {
								goto l109
							}
							if !rules[ruleexpression]() {
								goto l109
							}
							if !rules[ruleAS]() {
								goto l109
							}
							if !rules[rulevar]() {
								goto l109
							}
							if !rules[ruleRPAREN]() {
								goto l109
							}
						}
					l148:
						depth--
						add(ruleprojectionElem, position147)
					}
				l145:
					{
						position146, tokenIndex146, depth146 := position, tokenIndex, depth
						{
							position150 := position
							depth++
							{
								position151, tokenIndex151, depth151 := position, tokenIndex, depth
								if !rules[rulevar]() {
									goto l152
								}
								goto l151
							l152:
								position, tokenIndex, depth = position151, tokenIndex151, depth151
								if !rules[ruleLPAREN]() {
									goto l146
								}
								if !rules[ruleexpression]() {
									goto l146
								}
								if !rules[ruleAS]() {
									goto l146
								}
								if !rules[rulevar]() {
									goto l146
								}
								if !rules[ruleRPAREN]() {
									goto l146
								}
							}
						l151:
							depth--
							add(ruleprojectionElem, position150)
						}
						goto l145
					l146:
						position, tokenIndex, depth = position146, tokenIndex146, depth146
					}
				}
			l143:
				depth--
				add(ruleselect, position110)
			}
			return true
		l109:
			position, tokenIndex, depth = position109, tokenIndex109, depth109
			return false
		},
		/* 7 subSelect <- <(select whereClause solutionModifier)> */
		func() bool {
			position153, tokenIndex153, depth153 := position, tokenIndex, depth
			{
				position154 := position
				depth++
				if !rules[ruleselect]() {
					goto l153
				}
				if !rules[rulewhereClause]() {
					goto l153
				}
				if !rules[rulesolutionModifier]() {
					goto l153
				}
				depth--
				add(rulesubSelect, position154)
			}
			return true
		l153:
			position, tokenIndex, depth = position153, tokenIndex153, depth153
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
			position161, tokenIndex161, depth161 := position, tokenIndex, depth
			{
				position162 := position
				depth++
				{
					position163 := position
					depth++
					{
						position164, tokenIndex164, depth164 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l165
						}
						position++
						goto l164
					l165:
						position, tokenIndex, depth = position164, tokenIndex164, depth164
						if buffer[position] != rune('F') {
							goto l161
						}
						position++
					}
				l164:
					{
						position166, tokenIndex166, depth166 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l167
						}
						position++
						goto l166
					l167:
						position, tokenIndex, depth = position166, tokenIndex166, depth166
						if buffer[position] != rune('R') {
							goto l161
						}
						position++
					}
				l166:
					{
						position168, tokenIndex168, depth168 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l169
						}
						position++
						goto l168
					l169:
						position, tokenIndex, depth = position168, tokenIndex168, depth168
						if buffer[position] != rune('O') {
							goto l161
						}
						position++
					}
				l168:
					{
						position170, tokenIndex170, depth170 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l171
						}
						position++
						goto l170
					l171:
						position, tokenIndex, depth = position170, tokenIndex170, depth170
						if buffer[position] != rune('M') {
							goto l161
						}
						position++
					}
				l170:
					if !rules[ruleskip]() {
						goto l161
					}
					depth--
					add(ruleFROM, position163)
				}
				{
					position172, tokenIndex172, depth172 := position, tokenIndex, depth
					{
						position174 := position
						depth++
						{
							position175, tokenIndex175, depth175 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l176
							}
							position++
							goto l175
						l176:
							position, tokenIndex, depth = position175, tokenIndex175, depth175
							if buffer[position] != rune('N') {
								goto l172
							}
							position++
						}
					l175:
						{
							position177, tokenIndex177, depth177 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l178
							}
							position++
							goto l177
						l178:
							position, tokenIndex, depth = position177, tokenIndex177, depth177
							if buffer[position] != rune('A') {
								goto l172
							}
							position++
						}
					l177:
						{
							position179, tokenIndex179, depth179 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l180
							}
							position++
							goto l179
						l180:
							position, tokenIndex, depth = position179, tokenIndex179, depth179
							if buffer[position] != rune('M') {
								goto l172
							}
							position++
						}
					l179:
						{
							position181, tokenIndex181, depth181 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l182
							}
							position++
							goto l181
						l182:
							position, tokenIndex, depth = position181, tokenIndex181, depth181
							if buffer[position] != rune('E') {
								goto l172
							}
							position++
						}
					l181:
						{
							position183, tokenIndex183, depth183 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l184
							}
							position++
							goto l183
						l184:
							position, tokenIndex, depth = position183, tokenIndex183, depth183
							if buffer[position] != rune('D') {
								goto l172
							}
							position++
						}
					l183:
						if !rules[ruleskip]() {
							goto l172
						}
						depth--
						add(ruleNAMED, position174)
					}
					goto l173
				l172:
					position, tokenIndex, depth = position172, tokenIndex172, depth172
				}
			l173:
				if !rules[ruleiriref]() {
					goto l161
				}
				depth--
				add(ruledatasetClause, position162)
			}
			return true
		l161:
			position, tokenIndex, depth = position161, tokenIndex161, depth161
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position185, tokenIndex185, depth185 := position, tokenIndex, depth
			{
				position186 := position
				depth++
				{
					position187, tokenIndex187, depth187 := position, tokenIndex, depth
					{
						position189 := position
						depth++
						{
							position190, tokenIndex190, depth190 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l191
							}
							position++
							goto l190
						l191:
							position, tokenIndex, depth = position190, tokenIndex190, depth190
							if buffer[position] != rune('W') {
								goto l187
							}
							position++
						}
					l190:
						{
							position192, tokenIndex192, depth192 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l193
							}
							position++
							goto l192
						l193:
							position, tokenIndex, depth = position192, tokenIndex192, depth192
							if buffer[position] != rune('H') {
								goto l187
							}
							position++
						}
					l192:
						{
							position194, tokenIndex194, depth194 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l195
							}
							position++
							goto l194
						l195:
							position, tokenIndex, depth = position194, tokenIndex194, depth194
							if buffer[position] != rune('E') {
								goto l187
							}
							position++
						}
					l194:
						{
							position196, tokenIndex196, depth196 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l197
							}
							position++
							goto l196
						l197:
							position, tokenIndex, depth = position196, tokenIndex196, depth196
							if buffer[position] != rune('R') {
								goto l187
							}
							position++
						}
					l196:
						{
							position198, tokenIndex198, depth198 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l199
							}
							position++
							goto l198
						l199:
							position, tokenIndex, depth = position198, tokenIndex198, depth198
							if buffer[position] != rune('E') {
								goto l187
							}
							position++
						}
					l198:
						if !rules[ruleskip]() {
							goto l187
						}
						depth--
						add(ruleWHERE, position189)
					}
					goto l188
				l187:
					position, tokenIndex, depth = position187, tokenIndex187, depth187
				}
			l188:
				if !rules[rulegroupGraphPattern]() {
					goto l185
				}
				depth--
				add(rulewhereClause, position186)
			}
			return true
		l185:
			position, tokenIndex, depth = position185, tokenIndex185, depth185
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position200, tokenIndex200, depth200 := position, tokenIndex, depth
			{
				position201 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l200
				}
				{
					position202, tokenIndex202, depth202 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l203
					}
					goto l202
				l203:
					position, tokenIndex, depth = position202, tokenIndex202, depth202
					if !rules[rulegraphPattern]() {
						goto l200
					}
				}
			l202:
				if !rules[ruleRBRACE]() {
					goto l200
				}
				depth--
				add(rulegroupGraphPattern, position201)
			}
			return true
		l200:
			position, tokenIndex, depth = position200, tokenIndex200, depth200
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
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
							if !rules[ruletriplesBlock]() {
								goto l210
							}
						l211:
							{
								position212, tokenIndex212, depth212 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l212
								}
								{
									position213, tokenIndex213, depth213 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l213
									}
									goto l214
								l213:
									position, tokenIndex, depth = position213, tokenIndex213, depth213
								}
							l214:
								{
									position215, tokenIndex215, depth215 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l215
									}
									goto l216
								l215:
									position, tokenIndex, depth = position215, tokenIndex215, depth215
								}
							l216:
								goto l211
							l212:
								position, tokenIndex, depth = position212, tokenIndex212, depth212
							}
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if !rules[rulefilterOrBind]() {
								goto l206
							}
							{
								position219, tokenIndex219, depth219 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l219
								}
								goto l220
							l219:
								position, tokenIndex, depth = position219, tokenIndex219, depth219
							}
						l220:
							{
								position221, tokenIndex221, depth221 := position, tokenIndex, depth
								if !rules[ruletriplesBlock]() {
									goto l221
								}
								goto l222
							l221:
								position, tokenIndex, depth = position221, tokenIndex221, depth221
							}
						l222:
						l217:
							{
								position218, tokenIndex218, depth218 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l218
								}
								{
									position223, tokenIndex223, depth223 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l223
									}
									goto l224
								l223:
									position, tokenIndex, depth = position223, tokenIndex223, depth223
								}
							l224:
								{
									position225, tokenIndex225, depth225 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l225
									}
									goto l226
								l225:
									position, tokenIndex, depth = position225, tokenIndex225, depth225
								}
							l226:
								goto l217
							l218:
								position, tokenIndex, depth = position218, tokenIndex218, depth218
							}
						}
					l209:
						depth--
						add(rulebasicGraphPattern, position208)
					}
					goto l207
				l206:
					position, tokenIndex, depth = position206, tokenIndex206, depth206
				}
			l207:
				{
					position227, tokenIndex227, depth227 := position, tokenIndex, depth
					{
						position229 := position
						depth++
						{
							position230, tokenIndex230, depth230 := position, tokenIndex, depth
							{
								position232 := position
								depth++
								{
									position233 := position
									depth++
									{
										position234, tokenIndex234, depth234 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l235
										}
										position++
										goto l234
									l235:
										position, tokenIndex, depth = position234, tokenIndex234, depth234
										if buffer[position] != rune('O') {
											goto l231
										}
										position++
									}
								l234:
									{
										position236, tokenIndex236, depth236 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l237
										}
										position++
										goto l236
									l237:
										position, tokenIndex, depth = position236, tokenIndex236, depth236
										if buffer[position] != rune('P') {
											goto l231
										}
										position++
									}
								l236:
									{
										position238, tokenIndex238, depth238 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l239
										}
										position++
										goto l238
									l239:
										position, tokenIndex, depth = position238, tokenIndex238, depth238
										if buffer[position] != rune('T') {
											goto l231
										}
										position++
									}
								l238:
									{
										position240, tokenIndex240, depth240 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l241
										}
										position++
										goto l240
									l241:
										position, tokenIndex, depth = position240, tokenIndex240, depth240
										if buffer[position] != rune('I') {
											goto l231
										}
										position++
									}
								l240:
									{
										position242, tokenIndex242, depth242 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l243
										}
										position++
										goto l242
									l243:
										position, tokenIndex, depth = position242, tokenIndex242, depth242
										if buffer[position] != rune('O') {
											goto l231
										}
										position++
									}
								l242:
									{
										position244, tokenIndex244, depth244 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l245
										}
										position++
										goto l244
									l245:
										position, tokenIndex, depth = position244, tokenIndex244, depth244
										if buffer[position] != rune('N') {
											goto l231
										}
										position++
									}
								l244:
									{
										position246, tokenIndex246, depth246 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l247
										}
										position++
										goto l246
									l247:
										position, tokenIndex, depth = position246, tokenIndex246, depth246
										if buffer[position] != rune('A') {
											goto l231
										}
										position++
									}
								l246:
									{
										position248, tokenIndex248, depth248 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l249
										}
										position++
										goto l248
									l249:
										position, tokenIndex, depth = position248, tokenIndex248, depth248
										if buffer[position] != rune('L') {
											goto l231
										}
										position++
									}
								l248:
									if !rules[ruleskip]() {
										goto l231
									}
									depth--
									add(ruleOPTIONAL, position233)
								}
								if !rules[ruleLBRACE]() {
									goto l231
								}
								{
									position250, tokenIndex250, depth250 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l251
									}
									goto l250
								l251:
									position, tokenIndex, depth = position250, tokenIndex250, depth250
									if !rules[rulegraphPattern]() {
										goto l231
									}
								}
							l250:
								if !rules[ruleRBRACE]() {
									goto l231
								}
								depth--
								add(ruleoptionalGraphPattern, position232)
							}
							goto l230
						l231:
							position, tokenIndex, depth = position230, tokenIndex230, depth230
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l227
							}
						}
					l230:
						depth--
						add(rulegraphPatternNotTriples, position229)
					}
					{
						position252, tokenIndex252, depth252 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l252
						}
						goto l253
					l252:
						position, tokenIndex, depth = position252, tokenIndex252, depth252
					}
				l253:
					if !rules[rulegraphPattern]() {
						goto l227
					}
					goto l228
				l227:
					position, tokenIndex, depth = position227, tokenIndex227, depth227
				}
			l228:
				depth--
				add(rulegraphPattern, position205)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position256, tokenIndex256, depth256 := position, tokenIndex, depth
			{
				position257 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l256
				}
				{
					position258, tokenIndex258, depth258 := position, tokenIndex, depth
					{
						position260 := position
						depth++
						{
							position261, tokenIndex261, depth261 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l262
							}
							position++
							goto l261
						l262:
							position, tokenIndex, depth = position261, tokenIndex261, depth261
							if buffer[position] != rune('U') {
								goto l258
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
								goto l258
							}
							position++
						}
					l263:
						{
							position265, tokenIndex265, depth265 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l266
							}
							position++
							goto l265
						l266:
							position, tokenIndex, depth = position265, tokenIndex265, depth265
							if buffer[position] != rune('I') {
								goto l258
							}
							position++
						}
					l265:
						{
							position267, tokenIndex267, depth267 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l268
							}
							position++
							goto l267
						l268:
							position, tokenIndex, depth = position267, tokenIndex267, depth267
							if buffer[position] != rune('O') {
								goto l258
							}
							position++
						}
					l267:
						{
							position269, tokenIndex269, depth269 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l270
							}
							position++
							goto l269
						l270:
							position, tokenIndex, depth = position269, tokenIndex269, depth269
							if buffer[position] != rune('N') {
								goto l258
							}
							position++
						}
					l269:
						if !rules[ruleskip]() {
							goto l258
						}
						depth--
						add(ruleUNION, position260)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l258
					}
					goto l259
				l258:
					position, tokenIndex, depth = position258, tokenIndex258, depth258
				}
			l259:
				depth--
				add(rulegroupOrUnionGraphPattern, position257)
			}
			return true
		l256:
			position, tokenIndex, depth = position256, tokenIndex256, depth256
			return false
		},
		/* 21 basicGraphPattern <- <((triplesBlock (filterOrBind DOT? triplesBlock?)*) / (filterOrBind DOT? triplesBlock?)+)> */
		nil,
		/* 22 filterOrBind <- <((FILTER constraint) / (BIND LPAREN expression AS var RPAREN))> */
		func() bool {
			position272, tokenIndex272, depth272 := position, tokenIndex, depth
			{
				position273 := position
				depth++
				{
					position274, tokenIndex274, depth274 := position, tokenIndex, depth
					{
						position276 := position
						depth++
						{
							position277, tokenIndex277, depth277 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l278
							}
							position++
							goto l277
						l278:
							position, tokenIndex, depth = position277, tokenIndex277, depth277
							if buffer[position] != rune('F') {
								goto l275
							}
							position++
						}
					l277:
						{
							position279, tokenIndex279, depth279 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l280
							}
							position++
							goto l279
						l280:
							position, tokenIndex, depth = position279, tokenIndex279, depth279
							if buffer[position] != rune('I') {
								goto l275
							}
							position++
						}
					l279:
						{
							position281, tokenIndex281, depth281 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l282
							}
							position++
							goto l281
						l282:
							position, tokenIndex, depth = position281, tokenIndex281, depth281
							if buffer[position] != rune('L') {
								goto l275
							}
							position++
						}
					l281:
						{
							position283, tokenIndex283, depth283 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l284
							}
							position++
							goto l283
						l284:
							position, tokenIndex, depth = position283, tokenIndex283, depth283
							if buffer[position] != rune('T') {
								goto l275
							}
							position++
						}
					l283:
						{
							position285, tokenIndex285, depth285 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l286
							}
							position++
							goto l285
						l286:
							position, tokenIndex, depth = position285, tokenIndex285, depth285
							if buffer[position] != rune('E') {
								goto l275
							}
							position++
						}
					l285:
						{
							position287, tokenIndex287, depth287 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l288
							}
							position++
							goto l287
						l288:
							position, tokenIndex, depth = position287, tokenIndex287, depth287
							if buffer[position] != rune('R') {
								goto l275
							}
							position++
						}
					l287:
						if !rules[ruleskip]() {
							goto l275
						}
						depth--
						add(ruleFILTER, position276)
					}
					if !rules[ruleconstraint]() {
						goto l275
					}
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					{
						position289 := position
						depth++
						{
							position290, tokenIndex290, depth290 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l291
							}
							position++
							goto l290
						l291:
							position, tokenIndex, depth = position290, tokenIndex290, depth290
							if buffer[position] != rune('B') {
								goto l272
							}
							position++
						}
					l290:
						{
							position292, tokenIndex292, depth292 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l293
							}
							position++
							goto l292
						l293:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
							if buffer[position] != rune('I') {
								goto l272
							}
							position++
						}
					l292:
						{
							position294, tokenIndex294, depth294 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l295
							}
							position++
							goto l294
						l295:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							if buffer[position] != rune('N') {
								goto l272
							}
							position++
						}
					l294:
						{
							position296, tokenIndex296, depth296 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l297
							}
							position++
							goto l296
						l297:
							position, tokenIndex, depth = position296, tokenIndex296, depth296
							if buffer[position] != rune('D') {
								goto l272
							}
							position++
						}
					l296:
						if !rules[ruleskip]() {
							goto l272
						}
						depth--
						add(ruleBIND, position289)
					}
					if !rules[ruleLPAREN]() {
						goto l272
					}
					if !rules[ruleexpression]() {
						goto l272
					}
					if !rules[ruleAS]() {
						goto l272
					}
					if !rules[rulevar]() {
						goto l272
					}
					if !rules[ruleRPAREN]() {
						goto l272
					}
				}
			l274:
				depth--
				add(rulefilterOrBind, position273)
			}
			return true
		l272:
			position, tokenIndex, depth = position272, tokenIndex272, depth272
			return false
		},
		/* 23 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		func() bool {
			position298, tokenIndex298, depth298 := position, tokenIndex, depth
			{
				position299 := position
				depth++
				{
					position300, tokenIndex300, depth300 := position, tokenIndex, depth
					if !rules[rulebrackettedExpression]() {
						goto l301
					}
					goto l300
				l301:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if !rules[rulebuiltinCall]() {
						goto l302
					}
					goto l300
				l302:
					position, tokenIndex, depth = position300, tokenIndex300, depth300
					if !rules[rulefunctionCall]() {
						goto l298
					}
				}
			l300:
				depth--
				add(ruleconstraint, position299)
			}
			return true
		l298:
			position, tokenIndex, depth = position298, tokenIndex298, depth298
			return false
		},
		/* 24 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position303, tokenIndex303, depth303 := position, tokenIndex, depth
			{
				position304 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l303
				}
			l305:
				{
					position306, tokenIndex306, depth306 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l306
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l306
					}
					goto l305
				l306:
					position, tokenIndex, depth = position306, tokenIndex306, depth306
				}
				{
					position307, tokenIndex307, depth307 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l307
					}
					goto l308
				l307:
					position, tokenIndex, depth = position307, tokenIndex307, depth307
				}
			l308:
				depth--
				add(ruletriplesBlock, position304)
			}
			return true
		l303:
			position, tokenIndex, depth = position303, tokenIndex303, depth303
			return false
		},
		/* 25 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position309, tokenIndex309, depth309 := position, tokenIndex, depth
			{
				position310 := position
				depth++
				{
					position311, tokenIndex311, depth311 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l312
					}
					if !rules[rulepropertyListPath]() {
						goto l312
					}
					goto l311
				l312:
					position, tokenIndex, depth = position311, tokenIndex311, depth311
					{
						position313 := position
						depth++
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							{
								position316 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l315
								}
								if !rules[rulegraphNodePath]() {
									goto l315
								}
							l317:
								{
									position318, tokenIndex318, depth318 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l318
									}
									goto l317
								l318:
									position, tokenIndex, depth = position318, tokenIndex318, depth318
								}
								if !rules[ruleRPAREN]() {
									goto l315
								}
								depth--
								add(rulecollectionPath, position316)
							}
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							{
								position319 := position
								depth++
								{
									position320 := position
									depth++
									if buffer[position] != rune('[') {
										goto l309
									}
									position++
									if !rules[ruleskip]() {
										goto l309
									}
									depth--
									add(ruleLBRACK, position320)
								}
								if !rules[rulepropertyListPath]() {
									goto l309
								}
								{
									position321 := position
									depth++
									if buffer[position] != rune(']') {
										goto l309
									}
									position++
									if !rules[ruleskip]() {
										goto l309
									}
									depth--
									add(ruleRBRACK, position321)
								}
								depth--
								add(ruleblankNodePropertyListPath, position319)
							}
						}
					l314:
						depth--
						add(ruletriplesNodePath, position313)
					}
					if !rules[rulepropertyListPath]() {
						goto l309
					}
				}
			l311:
				depth--
				add(ruletriplesSameSubjectPath, position310)
			}
			return true
		l309:
			position, tokenIndex, depth = position309, tokenIndex309, depth309
			return false
		},
		/* 26 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l325
					}
					goto l324
				l325:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					{
						position326 := position
						depth++
						{
							position327, tokenIndex327, depth327 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l328
							}
							goto l327
						l328:
							position, tokenIndex, depth = position327, tokenIndex327, depth327
							{
								switch buffer[position] {
								case '(':
									if !rules[rulenil]() {
										goto l322
									}
									break
								case '[', '_':
									{
										position330 := position
										depth++
										{
											position331, tokenIndex331, depth331 := position, tokenIndex, depth
											{
												position333 := position
												depth++
												if buffer[position] != rune('_') {
													goto l332
												}
												position++
												if buffer[position] != rune(':') {
													goto l332
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l332
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l332
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l332
														}
														position++
														break
													}
												}

												{
													position335, tokenIndex335, depth335 := position, tokenIndex, depth
													{
														position337, tokenIndex337, depth337 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l338
														}
														position++
														goto l337
													l338:
														position, tokenIndex, depth = position337, tokenIndex337, depth337
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l339
														}
														position++
														goto l337
													l339:
														position, tokenIndex, depth = position337, tokenIndex337, depth337
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l340
														}
														position++
														goto l337
													l340:
														position, tokenIndex, depth = position337, tokenIndex337, depth337
														if c := buffer[position]; c < rune('.') || c > rune('_') {
															goto l335
														}
														position++
													}
												l337:
													goto l336
												l335:
													position, tokenIndex, depth = position335, tokenIndex335, depth335
												}
											l336:
												if !rules[ruleskip]() {
													goto l332
												}
												depth--
												add(ruleblankNodeLabel, position333)
											}
											goto l331
										l332:
											position, tokenIndex, depth = position331, tokenIndex331, depth331
											{
												position341 := position
												depth++
												if buffer[position] != rune('[') {
													goto l322
												}
												position++
											l342:
												{
													position343, tokenIndex343, depth343 := position, tokenIndex, depth
													if !rules[rulews]() {
														goto l343
													}
													goto l342
												l343:
													position, tokenIndex, depth = position343, tokenIndex343, depth343
												}
												if buffer[position] != rune(']') {
													goto l322
												}
												position++
												if !rules[ruleskip]() {
													goto l322
												}
												depth--
												add(ruleanon, position341)
											}
										}
									l331:
										depth--
										add(ruleblankNode, position330)
									}
									break
								case 'F', 'T', 'f', 't':
									if !rules[rulebooleanLiteral]() {
										goto l322
									}
									break
								case '"':
									if !rules[ruleliteral]() {
										goto l322
									}
									break
								default:
									if !rules[rulenumericLiteral]() {
										goto l322
									}
									break
								}
							}

						}
					l327:
						depth--
						add(rulegraphTerm, position326)
					}
				}
			l324:
				depth--
				add(rulevarOrTerm, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 27 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 28 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		nil,
		/* 29 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 30 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 31 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath)?)> */
		func() bool {
			position348, tokenIndex348, depth348 := position, tokenIndex, depth
			{
				position349 := position
				depth++
				{
					position350, tokenIndex350, depth350 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l351
					}
					goto l350
				l351:
					position, tokenIndex, depth = position350, tokenIndex350, depth350
					{
						position352 := position
						depth++
						if !rules[rulepath]() {
							goto l348
						}
						depth--
						add(ruleverbPath, position352)
					}
				}
			l350:
				if !rules[ruleobjectListPath]() {
					goto l348
				}
				{
					position353, tokenIndex353, depth353 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l353
					}
					if !rules[rulepropertyListPath]() {
						goto l353
					}
					goto l354
				l353:
					position, tokenIndex, depth = position353, tokenIndex353, depth353
				}
			l354:
				depth--
				add(rulepropertyListPath, position349)
			}
			return true
		l348:
			position, tokenIndex, depth = position348, tokenIndex348, depth348
			return false
		},
		/* 32 verbPath <- <path> */
		nil,
		/* 33 path <- <pathAlternative> */
		func() bool {
			position356, tokenIndex356, depth356 := position, tokenIndex, depth
			{
				position357 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l356
				}
				depth--
				add(rulepath, position357)
			}
			return true
		l356:
			position, tokenIndex, depth = position356, tokenIndex356, depth356
			return false
		},
		/* 34 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position358, tokenIndex358, depth358 := position, tokenIndex, depth
			{
				position359 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l358
				}
			l360:
				{
					position361, tokenIndex361, depth361 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l361
					}
					if !rules[rulepathAlternative]() {
						goto l361
					}
					goto l360
				l361:
					position, tokenIndex, depth = position361, tokenIndex361, depth361
				}
				depth--
				add(rulepathAlternative, position359)
			}
			return true
		l358:
			position, tokenIndex, depth = position358, tokenIndex358, depth358
			return false
		},
		/* 35 pathSequence <- <(pathElt (SLASH pathSequence)*)> */
		func() bool {
			position362, tokenIndex362, depth362 := position, tokenIndex, depth
			{
				position363 := position
				depth++
				{
					position364 := position
					depth++
					{
						position365, tokenIndex365, depth365 := position, tokenIndex, depth
						if !rules[ruleINVERSE]() {
							goto l365
						}
						goto l366
					l365:
						position, tokenIndex, depth = position365, tokenIndex365, depth365
					}
				l366:
					{
						position367 := position
						depth++
						{
							position368, tokenIndex368, depth368 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l369
							}
							goto l368
						l369:
							position, tokenIndex, depth = position368, tokenIndex368, depth368
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l362
									}
									if !rules[rulepath]() {
										goto l362
									}
									if !rules[ruleRPAREN]() {
										goto l362
									}
									break
								case '!':
									if !rules[ruleNOT]() {
										goto l362
									}
									{
										position371 := position
										depth++
										{
											position372, tokenIndex372, depth372 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l373
											}
											goto l372
										l373:
											position, tokenIndex, depth = position372, tokenIndex372, depth372
											if !rules[ruleLPAREN]() {
												goto l362
											}
											{
												position374, tokenIndex374, depth374 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l374
												}
											l376:
												{
													position377, tokenIndex377, depth377 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l377
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l377
													}
													goto l376
												l377:
													position, tokenIndex, depth = position377, tokenIndex377, depth377
												}
												goto l375
											l374:
												position, tokenIndex, depth = position374, tokenIndex374, depth374
											}
										l375:
											if !rules[ruleRPAREN]() {
												goto l362
											}
										}
									l372:
										depth--
										add(rulepathNegatedPropertySet, position371)
									}
									break
								default:
									if !rules[ruleISA]() {
										goto l362
									}
									break
								}
							}

						}
					l368:
						depth--
						add(rulepathPrimary, position367)
					}
					depth--
					add(rulepathElt, position364)
				}
			l378:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l379
					}
					if !rules[rulepathSequence]() {
						goto l379
					}
					goto l378
				l379:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
				}
				depth--
				add(rulepathSequence, position363)
			}
			return true
		l362:
			position, tokenIndex, depth = position362, tokenIndex362, depth362
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
			position383, tokenIndex383, depth383 := position, tokenIndex, depth
			{
				position384 := position
				depth++
				{
					position385, tokenIndex385, depth385 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l386
					}
					goto l385
				l386:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if !rules[ruleISA]() {
						goto l387
					}
					goto l385
				l387:
					position, tokenIndex, depth = position385, tokenIndex385, depth385
					if !rules[ruleINVERSE]() {
						goto l383
					}
					{
						position388, tokenIndex388, depth388 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l389
						}
						goto l388
					l389:
						position, tokenIndex, depth = position388, tokenIndex388, depth388
						if !rules[ruleISA]() {
							goto l383
						}
					}
				l388:
				}
			l385:
				depth--
				add(rulepathOneInPropertySet, position384)
			}
			return true
		l383:
			position, tokenIndex, depth = position383, tokenIndex383, depth383
			return false
		},
		/* 40 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l390
					}
					depth--
					add(ruleobjectPath, position392)
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
				add(ruleobjectListPath, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 41 objectPath <- <graphNodePath> */
		nil,
		/* 42 graphNodePath <- <varOrTerm> */
		func() bool {
			position396, tokenIndex396, depth396 := position, tokenIndex, depth
			{
				position397 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l396
				}
				depth--
				add(rulegraphNodePath, position397)
			}
			return true
		l396:
			position, tokenIndex, depth = position396, tokenIndex396, depth396
			return false
		},
		/* 43 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position399 := position
				depth++
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					{
						position402, tokenIndex402, depth402 := position, tokenIndex, depth
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
									goto l403
								}
								position++
							}
						l405:
							{
								position407, tokenIndex407, depth407 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l408
								}
								position++
								goto l407
							l408:
								position, tokenIndex, depth = position407, tokenIndex407, depth407
								if buffer[position] != rune('R') {
									goto l403
								}
								position++
							}
						l407:
							{
								position409, tokenIndex409, depth409 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l410
								}
								position++
								goto l409
							l410:
								position, tokenIndex, depth = position409, tokenIndex409, depth409
								if buffer[position] != rune('D') {
									goto l403
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
									goto l403
								}
								position++
							}
						l411:
							{
								position413, tokenIndex413, depth413 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l414
								}
								position++
								goto l413
							l414:
								position, tokenIndex, depth = position413, tokenIndex413, depth413
								if buffer[position] != rune('R') {
									goto l403
								}
								position++
							}
						l413:
							if !rules[ruleskip]() {
								goto l403
							}
							depth--
							add(ruleORDER, position404)
						}
						if !rules[ruleBY]() {
							goto l403
						}
						{
							position417 := position
							depth++
							{
								position418, tokenIndex418, depth418 := position, tokenIndex, depth
								{
									position420, tokenIndex420, depth420 := position, tokenIndex, depth
									{
										position422, tokenIndex422, depth422 := position, tokenIndex, depth
										{
											position424 := position
											depth++
											{
												position425, tokenIndex425, depth425 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l426
												}
												position++
												goto l425
											l426:
												position, tokenIndex, depth = position425, tokenIndex425, depth425
												if buffer[position] != rune('A') {
													goto l423
												}
												position++
											}
										l425:
											{
												position427, tokenIndex427, depth427 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l428
												}
												position++
												goto l427
											l428:
												position, tokenIndex, depth = position427, tokenIndex427, depth427
												if buffer[position] != rune('S') {
													goto l423
												}
												position++
											}
										l427:
											{
												position429, tokenIndex429, depth429 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l430
												}
												position++
												goto l429
											l430:
												position, tokenIndex, depth = position429, tokenIndex429, depth429
												if buffer[position] != rune('C') {
													goto l423
												}
												position++
											}
										l429:
											if !rules[ruleskip]() {
												goto l423
											}
											depth--
											add(ruleASC, position424)
										}
										goto l422
									l423:
										position, tokenIndex, depth = position422, tokenIndex422, depth422
										{
											position431 := position
											depth++
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
													goto l420
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
													goto l420
												}
												position++
											}
										l434:
											{
												position436, tokenIndex436, depth436 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l437
												}
												position++
												goto l436
											l437:
												position, tokenIndex, depth = position436, tokenIndex436, depth436
												if buffer[position] != rune('S') {
													goto l420
												}
												position++
											}
										l436:
											{
												position438, tokenIndex438, depth438 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l439
												}
												position++
												goto l438
											l439:
												position, tokenIndex, depth = position438, tokenIndex438, depth438
												if buffer[position] != rune('C') {
													goto l420
												}
												position++
											}
										l438:
											if !rules[ruleskip]() {
												goto l420
											}
											depth--
											add(ruleDESC, position431)
										}
									}
								l422:
									goto l421
								l420:
									position, tokenIndex, depth = position420, tokenIndex420, depth420
								}
							l421:
								if !rules[rulebrackettedExpression]() {
									goto l419
								}
								goto l418
							l419:
								position, tokenIndex, depth = position418, tokenIndex418, depth418
								if !rules[rulefunctionCall]() {
									goto l440
								}
								goto l418
							l440:
								position, tokenIndex, depth = position418, tokenIndex418, depth418
								if !rules[rulebuiltinCall]() {
									goto l441
								}
								goto l418
							l441:
								position, tokenIndex, depth = position418, tokenIndex418, depth418
								if !rules[rulevar]() {
									goto l403
								}
							}
						l418:
							depth--
							add(ruleorderCondition, position417)
						}
					l415:
						{
							position416, tokenIndex416, depth416 := position, tokenIndex, depth
							{
								position442 := position
								depth++
								{
									position443, tokenIndex443, depth443 := position, tokenIndex, depth
									{
										position445, tokenIndex445, depth445 := position, tokenIndex, depth
										{
											position447, tokenIndex447, depth447 := position, tokenIndex, depth
											{
												position449 := position
												depth++
												{
													position450, tokenIndex450, depth450 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l451
													}
													position++
													goto l450
												l451:
													position, tokenIndex, depth = position450, tokenIndex450, depth450
													if buffer[position] != rune('A') {
														goto l448
													}
													position++
												}
											l450:
												{
													position452, tokenIndex452, depth452 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l453
													}
													position++
													goto l452
												l453:
													position, tokenIndex, depth = position452, tokenIndex452, depth452
													if buffer[position] != rune('S') {
														goto l448
													}
													position++
												}
											l452:
												{
													position454, tokenIndex454, depth454 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l455
													}
													position++
													goto l454
												l455:
													position, tokenIndex, depth = position454, tokenIndex454, depth454
													if buffer[position] != rune('C') {
														goto l448
													}
													position++
												}
											l454:
												if !rules[ruleskip]() {
													goto l448
												}
												depth--
												add(ruleASC, position449)
											}
											goto l447
										l448:
											position, tokenIndex, depth = position447, tokenIndex447, depth447
											{
												position456 := position
												depth++
												{
													position457, tokenIndex457, depth457 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l458
													}
													position++
													goto l457
												l458:
													position, tokenIndex, depth = position457, tokenIndex457, depth457
													if buffer[position] != rune('D') {
														goto l445
													}
													position++
												}
											l457:
												{
													position459, tokenIndex459, depth459 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l460
													}
													position++
													goto l459
												l460:
													position, tokenIndex, depth = position459, tokenIndex459, depth459
													if buffer[position] != rune('E') {
														goto l445
													}
													position++
												}
											l459:
												{
													position461, tokenIndex461, depth461 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l462
													}
													position++
													goto l461
												l462:
													position, tokenIndex, depth = position461, tokenIndex461, depth461
													if buffer[position] != rune('S') {
														goto l445
													}
													position++
												}
											l461:
												{
													position463, tokenIndex463, depth463 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l464
													}
													position++
													goto l463
												l464:
													position, tokenIndex, depth = position463, tokenIndex463, depth463
													if buffer[position] != rune('C') {
														goto l445
													}
													position++
												}
											l463:
												if !rules[ruleskip]() {
													goto l445
												}
												depth--
												add(ruleDESC, position456)
											}
										}
									l447:
										goto l446
									l445:
										position, tokenIndex, depth = position445, tokenIndex445, depth445
									}
								l446:
									if !rules[rulebrackettedExpression]() {
										goto l444
									}
									goto l443
								l444:
									position, tokenIndex, depth = position443, tokenIndex443, depth443
									if !rules[rulefunctionCall]() {
										goto l465
									}
									goto l443
								l465:
									position, tokenIndex, depth = position443, tokenIndex443, depth443
									if !rules[rulebuiltinCall]() {
										goto l466
									}
									goto l443
								l466:
									position, tokenIndex, depth = position443, tokenIndex443, depth443
									if !rules[rulevar]() {
										goto l416
									}
								}
							l443:
								depth--
								add(ruleorderCondition, position442)
							}
							goto l415
						l416:
							position, tokenIndex, depth = position416, tokenIndex416, depth416
						}
						goto l402
					l403:
						position, tokenIndex, depth = position402, tokenIndex402, depth402
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position468 := position
									depth++
									{
										position469, tokenIndex469, depth469 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l470
										}
										position++
										goto l469
									l470:
										position, tokenIndex, depth = position469, tokenIndex469, depth469
										if buffer[position] != rune('H') {
											goto l400
										}
										position++
									}
								l469:
									{
										position471, tokenIndex471, depth471 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l472
										}
										position++
										goto l471
									l472:
										position, tokenIndex, depth = position471, tokenIndex471, depth471
										if buffer[position] != rune('A') {
											goto l400
										}
										position++
									}
								l471:
									{
										position473, tokenIndex473, depth473 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l474
										}
										position++
										goto l473
									l474:
										position, tokenIndex, depth = position473, tokenIndex473, depth473
										if buffer[position] != rune('V') {
											goto l400
										}
										position++
									}
								l473:
									{
										position475, tokenIndex475, depth475 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l476
										}
										position++
										goto l475
									l476:
										position, tokenIndex, depth = position475, tokenIndex475, depth475
										if buffer[position] != rune('I') {
											goto l400
										}
										position++
									}
								l475:
									{
										position477, tokenIndex477, depth477 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l478
										}
										position++
										goto l477
									l478:
										position, tokenIndex, depth = position477, tokenIndex477, depth477
										if buffer[position] != rune('N') {
											goto l400
										}
										position++
									}
								l477:
									{
										position479, tokenIndex479, depth479 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l480
										}
										position++
										goto l479
									l480:
										position, tokenIndex, depth = position479, tokenIndex479, depth479
										if buffer[position] != rune('G') {
											goto l400
										}
										position++
									}
								l479:
									if !rules[ruleskip]() {
										goto l400
									}
									depth--
									add(ruleHAVING, position468)
								}
								if !rules[ruleconstraint]() {
									goto l400
								}
								break
							case 'G', 'g':
								{
									position481 := position
									depth++
									{
										position482, tokenIndex482, depth482 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l483
										}
										position++
										goto l482
									l483:
										position, tokenIndex, depth = position482, tokenIndex482, depth482
										if buffer[position] != rune('G') {
											goto l400
										}
										position++
									}
								l482:
									{
										position484, tokenIndex484, depth484 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l485
										}
										position++
										goto l484
									l485:
										position, tokenIndex, depth = position484, tokenIndex484, depth484
										if buffer[position] != rune('R') {
											goto l400
										}
										position++
									}
								l484:
									{
										position486, tokenIndex486, depth486 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l487
										}
										position++
										goto l486
									l487:
										position, tokenIndex, depth = position486, tokenIndex486, depth486
										if buffer[position] != rune('O') {
											goto l400
										}
										position++
									}
								l486:
									{
										position488, tokenIndex488, depth488 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l489
										}
										position++
										goto l488
									l489:
										position, tokenIndex, depth = position488, tokenIndex488, depth488
										if buffer[position] != rune('U') {
											goto l400
										}
										position++
									}
								l488:
									{
										position490, tokenIndex490, depth490 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l491
										}
										position++
										goto l490
									l491:
										position, tokenIndex, depth = position490, tokenIndex490, depth490
										if buffer[position] != rune('P') {
											goto l400
										}
										position++
									}
								l490:
									if !rules[ruleskip]() {
										goto l400
									}
									depth--
									add(ruleGROUP, position481)
								}
								if !rules[ruleBY]() {
									goto l400
								}
								{
									position494 := position
									depth++
									{
										position495, tokenIndex495, depth495 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l496
										}
										goto l495
									l496:
										position, tokenIndex, depth = position495, tokenIndex495, depth495
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l400
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l400
												}
												if !rules[ruleexpression]() {
													goto l400
												}
												{
													position498, tokenIndex498, depth498 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l498
													}
													if !rules[rulevar]() {
														goto l498
													}
													goto l499
												l498:
													position, tokenIndex, depth = position498, tokenIndex498, depth498
												}
											l499:
												if !rules[ruleRPAREN]() {
													goto l400
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l400
												}
												break
											}
										}

									}
								l495:
									depth--
									add(rulegroupCondition, position494)
								}
							l492:
								{
									position493, tokenIndex493, depth493 := position, tokenIndex, depth
									{
										position500 := position
										depth++
										{
											position501, tokenIndex501, depth501 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l502
											}
											goto l501
										l502:
											position, tokenIndex, depth = position501, tokenIndex501, depth501
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l493
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l493
													}
													if !rules[ruleexpression]() {
														goto l493
													}
													{
														position504, tokenIndex504, depth504 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l504
														}
														if !rules[rulevar]() {
															goto l504
														}
														goto l505
													l504:
														position, tokenIndex, depth = position504, tokenIndex504, depth504
													}
												l505:
													if !rules[ruleRPAREN]() {
														goto l493
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l493
													}
													break
												}
											}

										}
									l501:
										depth--
										add(rulegroupCondition, position500)
									}
									goto l492
								l493:
									position, tokenIndex, depth = position493, tokenIndex493, depth493
								}
								break
							default:
								{
									position506 := position
									depth++
									{
										position507, tokenIndex507, depth507 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l508
										}
										{
											position509, tokenIndex509, depth509 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l509
											}
											goto l510
										l509:
											position, tokenIndex, depth = position509, tokenIndex509, depth509
										}
									l510:
										goto l507
									l508:
										position, tokenIndex, depth = position507, tokenIndex507, depth507
										if !rules[ruleoffset]() {
											goto l400
										}
										{
											position511, tokenIndex511, depth511 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l511
											}
											goto l512
										l511:
											position, tokenIndex, depth = position511, tokenIndex511, depth511
										}
									l512:
									}
								l507:
									depth--
									add(rulelimitOffsetClauses, position506)
								}
								break
							}
						}

					}
				l402:
					goto l401
				l400:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
				}
			l401:
				depth--
				add(rulesolutionModifier, position399)
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
			position516, tokenIndex516, depth516 := position, tokenIndex, depth
			{
				position517 := position
				depth++
				{
					position518 := position
					depth++
					{
						position519, tokenIndex519, depth519 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l520
						}
						position++
						goto l519
					l520:
						position, tokenIndex, depth = position519, tokenIndex519, depth519
						if buffer[position] != rune('L') {
							goto l516
						}
						position++
					}
				l519:
					{
						position521, tokenIndex521, depth521 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l522
						}
						position++
						goto l521
					l522:
						position, tokenIndex, depth = position521, tokenIndex521, depth521
						if buffer[position] != rune('I') {
							goto l516
						}
						position++
					}
				l521:
					{
						position523, tokenIndex523, depth523 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l524
						}
						position++
						goto l523
					l524:
						position, tokenIndex, depth = position523, tokenIndex523, depth523
						if buffer[position] != rune('M') {
							goto l516
						}
						position++
					}
				l523:
					{
						position525, tokenIndex525, depth525 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l526
						}
						position++
						goto l525
					l526:
						position, tokenIndex, depth = position525, tokenIndex525, depth525
						if buffer[position] != rune('I') {
							goto l516
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
							goto l516
						}
						position++
					}
				l527:
					if !rules[ruleskip]() {
						goto l516
					}
					depth--
					add(ruleLIMIT, position518)
				}
				if !rules[ruleINTEGER]() {
					goto l516
				}
				depth--
				add(rulelimit, position517)
			}
			return true
		l516:
			position, tokenIndex, depth = position516, tokenIndex516, depth516
			return false
		},
		/* 48 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position529, tokenIndex529, depth529 := position, tokenIndex, depth
			{
				position530 := position
				depth++
				{
					position531 := position
					depth++
					{
						position532, tokenIndex532, depth532 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l533
						}
						position++
						goto l532
					l533:
						position, tokenIndex, depth = position532, tokenIndex532, depth532
						if buffer[position] != rune('O') {
							goto l529
						}
						position++
					}
				l532:
					{
						position534, tokenIndex534, depth534 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l535
						}
						position++
						goto l534
					l535:
						position, tokenIndex, depth = position534, tokenIndex534, depth534
						if buffer[position] != rune('F') {
							goto l529
						}
						position++
					}
				l534:
					{
						position536, tokenIndex536, depth536 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l537
						}
						position++
						goto l536
					l537:
						position, tokenIndex, depth = position536, tokenIndex536, depth536
						if buffer[position] != rune('F') {
							goto l529
						}
						position++
					}
				l536:
					{
						position538, tokenIndex538, depth538 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l539
						}
						position++
						goto l538
					l539:
						position, tokenIndex, depth = position538, tokenIndex538, depth538
						if buffer[position] != rune('S') {
							goto l529
						}
						position++
					}
				l538:
					{
						position540, tokenIndex540, depth540 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l541
						}
						position++
						goto l540
					l541:
						position, tokenIndex, depth = position540, tokenIndex540, depth540
						if buffer[position] != rune('E') {
							goto l529
						}
						position++
					}
				l540:
					{
						position542, tokenIndex542, depth542 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l543
						}
						position++
						goto l542
					l543:
						position, tokenIndex, depth = position542, tokenIndex542, depth542
						if buffer[position] != rune('T') {
							goto l529
						}
						position++
					}
				l542:
					if !rules[ruleskip]() {
						goto l529
					}
					depth--
					add(ruleOFFSET, position531)
				}
				if !rules[ruleINTEGER]() {
					goto l529
				}
				depth--
				add(ruleoffset, position530)
			}
			return true
		l529:
			position, tokenIndex, depth = position529, tokenIndex529, depth529
			return false
		},
		/* 49 expression <- <conditionalOrExpression> */
		func() bool {
			position544, tokenIndex544, depth544 := position, tokenIndex, depth
			{
				position545 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l544
				}
				depth--
				add(ruleexpression, position545)
			}
			return true
		l544:
			position, tokenIndex, depth = position544, tokenIndex544, depth544
			return false
		},
		/* 50 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position546, tokenIndex546, depth546 := position, tokenIndex, depth
			{
				position547 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l546
				}
				{
					position548, tokenIndex548, depth548 := position, tokenIndex, depth
					{
						position550 := position
						depth++
						if buffer[position] != rune('|') {
							goto l548
						}
						position++
						if buffer[position] != rune('|') {
							goto l548
						}
						position++
						if !rules[ruleskip]() {
							goto l548
						}
						depth--
						add(ruleOR, position550)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l548
					}
					goto l549
				l548:
					position, tokenIndex, depth = position548, tokenIndex548, depth548
				}
			l549:
				depth--
				add(ruleconditionalOrExpression, position547)
			}
			return true
		l546:
			position, tokenIndex, depth = position546, tokenIndex546, depth546
			return false
		},
		/* 51 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position551, tokenIndex551, depth551 := position, tokenIndex, depth
			{
				position552 := position
				depth++
				{
					position553 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l551
					}
					{
						position554, tokenIndex554, depth554 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position557 := position
									depth++
									{
										position558 := position
										depth++
										{
											position559, tokenIndex559, depth559 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l560
											}
											position++
											goto l559
										l560:
											position, tokenIndex, depth = position559, tokenIndex559, depth559
											if buffer[position] != rune('N') {
												goto l554
											}
											position++
										}
									l559:
										{
											position561, tokenIndex561, depth561 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l562
											}
											position++
											goto l561
										l562:
											position, tokenIndex, depth = position561, tokenIndex561, depth561
											if buffer[position] != rune('O') {
												goto l554
											}
											position++
										}
									l561:
										{
											position563, tokenIndex563, depth563 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l564
											}
											position++
											goto l563
										l564:
											position, tokenIndex, depth = position563, tokenIndex563, depth563
											if buffer[position] != rune('T') {
												goto l554
											}
											position++
										}
									l563:
										if buffer[position] != rune(' ') {
											goto l554
										}
										position++
										{
											position565, tokenIndex565, depth565 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l566
											}
											position++
											goto l565
										l566:
											position, tokenIndex, depth = position565, tokenIndex565, depth565
											if buffer[position] != rune('I') {
												goto l554
											}
											position++
										}
									l565:
										{
											position567, tokenIndex567, depth567 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l568
											}
											position++
											goto l567
										l568:
											position, tokenIndex, depth = position567, tokenIndex567, depth567
											if buffer[position] != rune('N') {
												goto l554
											}
											position++
										}
									l567:
										if !rules[ruleskip]() {
											goto l554
										}
										depth--
										add(ruleNOTIN, position558)
									}
									if !rules[ruleargList]() {
										goto l554
									}
									depth--
									add(rulenotin, position557)
								}
								break
							case 'I', 'i':
								{
									position569 := position
									depth++
									{
										position570 := position
										depth++
										{
											position571, tokenIndex571, depth571 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l572
											}
											position++
											goto l571
										l572:
											position, tokenIndex, depth = position571, tokenIndex571, depth571
											if buffer[position] != rune('I') {
												goto l554
											}
											position++
										}
									l571:
										{
											position573, tokenIndex573, depth573 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l574
											}
											position++
											goto l573
										l574:
											position, tokenIndex, depth = position573, tokenIndex573, depth573
											if buffer[position] != rune('N') {
												goto l554
											}
											position++
										}
									l573:
										if !rules[ruleskip]() {
											goto l554
										}
										depth--
										add(ruleIN, position570)
									}
									if !rules[ruleargList]() {
										goto l554
									}
									depth--
									add(rulein, position569)
								}
								break
							default:
								{
									position575, tokenIndex575, depth575 := position, tokenIndex, depth
									{
										position577 := position
										depth++
										if buffer[position] != rune('<') {
											goto l576
										}
										position++
										if !rules[ruleskip]() {
											goto l576
										}
										depth--
										add(ruleLT, position577)
									}
									goto l575
								l576:
									position, tokenIndex, depth = position575, tokenIndex575, depth575
									{
										position579 := position
										depth++
										if buffer[position] != rune('>') {
											goto l578
										}
										position++
										if buffer[position] != rune('=') {
											goto l578
										}
										position++
										if !rules[ruleskip]() {
											goto l578
										}
										depth--
										add(ruleGE, position579)
									}
									goto l575
								l578:
									position, tokenIndex, depth = position575, tokenIndex575, depth575
									{
										switch buffer[position] {
										case '>':
											{
												position581 := position
												depth++
												if buffer[position] != rune('>') {
													goto l554
												}
												position++
												if !rules[ruleskip]() {
													goto l554
												}
												depth--
												add(ruleGT, position581)
											}
											break
										case '<':
											{
												position582 := position
												depth++
												if buffer[position] != rune('<') {
													goto l554
												}
												position++
												if buffer[position] != rune('=') {
													goto l554
												}
												position++
												if !rules[ruleskip]() {
													goto l554
												}
												depth--
												add(ruleLE, position582)
											}
											break
										case '!':
											{
												position583 := position
												depth++
												if buffer[position] != rune('!') {
													goto l554
												}
												position++
												if buffer[position] != rune('=') {
													goto l554
												}
												position++
												if !rules[ruleskip]() {
													goto l554
												}
												depth--
												add(ruleNE, position583)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l554
											}
											break
										}
									}

								}
							l575:
								if !rules[rulenumericExpression]() {
									goto l554
								}
								break
							}
						}

						goto l555
					l554:
						position, tokenIndex, depth = position554, tokenIndex554, depth554
					}
				l555:
					depth--
					add(rulevalueLogical, position553)
				}
				{
					position584, tokenIndex584, depth584 := position, tokenIndex, depth
					{
						position586 := position
						depth++
						if buffer[position] != rune('&') {
							goto l584
						}
						position++
						if buffer[position] != rune('&') {
							goto l584
						}
						position++
						if !rules[ruleskip]() {
							goto l584
						}
						depth--
						add(ruleAND, position586)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l584
					}
					goto l585
				l584:
					position, tokenIndex, depth = position584, tokenIndex584, depth584
				}
			l585:
				depth--
				add(ruleconditionalAndExpression, position552)
			}
			return true
		l551:
			position, tokenIndex, depth = position551, tokenIndex551, depth551
			return false
		},
		/* 52 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 53 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position588, tokenIndex588, depth588 := position, tokenIndex, depth
			{
				position589 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l588
				}
			l590:
				{
					position591, tokenIndex591, depth591 := position, tokenIndex, depth
					{
						position592, tokenIndex592, depth592 := position, tokenIndex, depth
						{
							position594, tokenIndex594, depth594 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l595
							}
							goto l594
						l595:
							position, tokenIndex, depth = position594, tokenIndex594, depth594
							if !rules[ruleMINUS]() {
								goto l593
							}
						}
					l594:
						if !rules[rulemultiplicativeExpression]() {
							goto l593
						}
						goto l592
					l593:
						position, tokenIndex, depth = position592, tokenIndex592, depth592
						{
							position596 := position
							depth++
							{
								position597, tokenIndex597, depth597 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l598
								}
								position++
								goto l597
							l598:
								position, tokenIndex, depth = position597, tokenIndex597, depth597
								if buffer[position] != rune('-') {
									goto l591
								}
								position++
							}
						l597:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l591
							}
							position++
						l599:
							{
								position600, tokenIndex600, depth600 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l600
								}
								position++
								goto l599
							l600:
								position, tokenIndex, depth = position600, tokenIndex600, depth600
							}
							{
								position601, tokenIndex601, depth601 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
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
								goto l602
							l601:
								position, tokenIndex, depth = position601, tokenIndex601, depth601
							}
						l602:
							if !rules[ruleskip]() {
								goto l591
							}
							depth--
							add(rulesignedNumericLiteral, position596)
						}
					}
				l592:
					goto l590
				l591:
					position, tokenIndex, depth = position591, tokenIndex591, depth591
				}
				depth--
				add(rulenumericExpression, position589)
			}
			return true
		l588:
			position, tokenIndex, depth = position588, tokenIndex588, depth588
			return false
		},
		/* 54 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l605
				}
			l607:
				{
					position608, tokenIndex608, depth608 := position, tokenIndex, depth
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l610
						}
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if !rules[ruleSLASH]() {
							goto l608
						}
					}
				l609:
					if !rules[ruleunaryExpression]() {
						goto l608
					}
					goto l607
				l608:
					position, tokenIndex, depth = position608, tokenIndex608, depth608
				}
				depth--
				add(rulemultiplicativeExpression, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 55 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position611, tokenIndex611, depth611 := position, tokenIndex, depth
			{
				position612 := position
				depth++
				{
					position613, tokenIndex613, depth613 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l613
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l613
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l613
							}
							break
						}
					}

					goto l614
				l613:
					position, tokenIndex, depth = position613, tokenIndex613, depth613
				}
			l614:
				{
					position616 := position
					depth++
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if !rules[rulebrackettedExpression]() {
							goto l618
						}
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if !rules[rulebuiltinCall]() {
							goto l619
						}
						goto l617
					l619:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if !rules[rulefunctionCall]() {
							goto l620
						}
						goto l617
					l620:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if !rules[ruleiriref]() {
							goto l621
						}
						goto l617
					l621:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position623 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position625 := position
												depth++
												{
													position626 := position
													depth++
													{
														position627, tokenIndex627, depth627 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l628
														}
														position++
														goto l627
													l628:
														position, tokenIndex, depth = position627, tokenIndex627, depth627
														if buffer[position] != rune('G') {
															goto l611
														}
														position++
													}
												l627:
													{
														position629, tokenIndex629, depth629 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l630
														}
														position++
														goto l629
													l630:
														position, tokenIndex, depth = position629, tokenIndex629, depth629
														if buffer[position] != rune('R') {
															goto l611
														}
														position++
													}
												l629:
													{
														position631, tokenIndex631, depth631 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l632
														}
														position++
														goto l631
													l632:
														position, tokenIndex, depth = position631, tokenIndex631, depth631
														if buffer[position] != rune('O') {
															goto l611
														}
														position++
													}
												l631:
													{
														position633, tokenIndex633, depth633 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l634
														}
														position++
														goto l633
													l634:
														position, tokenIndex, depth = position633, tokenIndex633, depth633
														if buffer[position] != rune('U') {
															goto l611
														}
														position++
													}
												l633:
													{
														position635, tokenIndex635, depth635 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l636
														}
														position++
														goto l635
													l636:
														position, tokenIndex, depth = position635, tokenIndex635, depth635
														if buffer[position] != rune('P') {
															goto l611
														}
														position++
													}
												l635:
													if buffer[position] != rune('_') {
														goto l611
													}
													position++
													{
														position637, tokenIndex637, depth637 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l638
														}
														position++
														goto l637
													l638:
														position, tokenIndex, depth = position637, tokenIndex637, depth637
														if buffer[position] != rune('C') {
															goto l611
														}
														position++
													}
												l637:
													{
														position639, tokenIndex639, depth639 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l640
														}
														position++
														goto l639
													l640:
														position, tokenIndex, depth = position639, tokenIndex639, depth639
														if buffer[position] != rune('O') {
															goto l611
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
															goto l611
														}
														position++
													}
												l641:
													{
														position643, tokenIndex643, depth643 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l644
														}
														position++
														goto l643
													l644:
														position, tokenIndex, depth = position643, tokenIndex643, depth643
														if buffer[position] != rune('C') {
															goto l611
														}
														position++
													}
												l643:
													{
														position645, tokenIndex645, depth645 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l646
														}
														position++
														goto l645
													l646:
														position, tokenIndex, depth = position645, tokenIndex645, depth645
														if buffer[position] != rune('A') {
															goto l611
														}
														position++
													}
												l645:
													{
														position647, tokenIndex647, depth647 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l648
														}
														position++
														goto l647
													l648:
														position, tokenIndex, depth = position647, tokenIndex647, depth647
														if buffer[position] != rune('T') {
															goto l611
														}
														position++
													}
												l647:
													if !rules[ruleskip]() {
														goto l611
													}
													depth--
													add(ruleGROUPCONCAT, position626)
												}
												if !rules[ruleLPAREN]() {
													goto l611
												}
												{
													position649, tokenIndex649, depth649 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l649
													}
													goto l650
												l649:
													position, tokenIndex, depth = position649, tokenIndex649, depth649
												}
											l650:
												if !rules[ruleexpression]() {
													goto l611
												}
												{
													position651, tokenIndex651, depth651 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l651
													}
													{
														position653 := position
														depth++
														{
															position654, tokenIndex654, depth654 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l655
															}
															position++
															goto l654
														l655:
															position, tokenIndex, depth = position654, tokenIndex654, depth654
															if buffer[position] != rune('S') {
																goto l651
															}
															position++
														}
													l654:
														{
															position656, tokenIndex656, depth656 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l657
															}
															position++
															goto l656
														l657:
															position, tokenIndex, depth = position656, tokenIndex656, depth656
															if buffer[position] != rune('E') {
																goto l651
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
																goto l651
															}
															position++
														}
													l658:
														{
															position660, tokenIndex660, depth660 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l661
															}
															position++
															goto l660
														l661:
															position, tokenIndex, depth = position660, tokenIndex660, depth660
															if buffer[position] != rune('A') {
																goto l651
															}
															position++
														}
													l660:
														{
															position662, tokenIndex662, depth662 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l663
															}
															position++
															goto l662
														l663:
															position, tokenIndex, depth = position662, tokenIndex662, depth662
															if buffer[position] != rune('R') {
																goto l651
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
																goto l651
															}
															position++
														}
													l664:
														{
															position666, tokenIndex666, depth666 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l667
															}
															position++
															goto l666
														l667:
															position, tokenIndex, depth = position666, tokenIndex666, depth666
															if buffer[position] != rune('T') {
																goto l651
															}
															position++
														}
													l666:
														{
															position668, tokenIndex668, depth668 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l669
															}
															position++
															goto l668
														l669:
															position, tokenIndex, depth = position668, tokenIndex668, depth668
															if buffer[position] != rune('O') {
																goto l651
															}
															position++
														}
													l668:
														{
															position670, tokenIndex670, depth670 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l671
															}
															position++
															goto l670
														l671:
															position, tokenIndex, depth = position670, tokenIndex670, depth670
															if buffer[position] != rune('R') {
																goto l651
															}
															position++
														}
													l670:
														if !rules[ruleskip]() {
															goto l651
														}
														depth--
														add(ruleSEPARATOR, position653)
													}
													if !rules[ruleEQ]() {
														goto l651
													}
													if !rules[rulestring]() {
														goto l651
													}
													goto l652
												l651:
													position, tokenIndex, depth = position651, tokenIndex651, depth651
												}
											l652:
												if !rules[ruleRPAREN]() {
													goto l611
												}
												depth--
												add(rulegroupConcat, position625)
											}
											break
										case 'C', 'c':
											{
												position672 := position
												depth++
												{
													position673 := position
													depth++
													{
														position674, tokenIndex674, depth674 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l675
														}
														position++
														goto l674
													l675:
														position, tokenIndex, depth = position674, tokenIndex674, depth674
														if buffer[position] != rune('C') {
															goto l611
														}
														position++
													}
												l674:
													{
														position676, tokenIndex676, depth676 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l677
														}
														position++
														goto l676
													l677:
														position, tokenIndex, depth = position676, tokenIndex676, depth676
														if buffer[position] != rune('O') {
															goto l611
														}
														position++
													}
												l676:
													{
														position678, tokenIndex678, depth678 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l679
														}
														position++
														goto l678
													l679:
														position, tokenIndex, depth = position678, tokenIndex678, depth678
														if buffer[position] != rune('U') {
															goto l611
														}
														position++
													}
												l678:
													{
														position680, tokenIndex680, depth680 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l681
														}
														position++
														goto l680
													l681:
														position, tokenIndex, depth = position680, tokenIndex680, depth680
														if buffer[position] != rune('N') {
															goto l611
														}
														position++
													}
												l680:
													{
														position682, tokenIndex682, depth682 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l683
														}
														position++
														goto l682
													l683:
														position, tokenIndex, depth = position682, tokenIndex682, depth682
														if buffer[position] != rune('T') {
															goto l611
														}
														position++
													}
												l682:
													if !rules[ruleskip]() {
														goto l611
													}
													depth--
													add(ruleCOUNT, position673)
												}
												if !rules[ruleLPAREN]() {
													goto l611
												}
												{
													position684, tokenIndex684, depth684 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l684
													}
													goto l685
												l684:
													position, tokenIndex, depth = position684, tokenIndex684, depth684
												}
											l685:
												{
													position686, tokenIndex686, depth686 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l687
													}
													goto l686
												l687:
													position, tokenIndex, depth = position686, tokenIndex686, depth686
													if !rules[ruleexpression]() {
														goto l611
													}
												}
											l686:
												if !rules[ruleRPAREN]() {
													goto l611
												}
												depth--
												add(rulecount, position672)
											}
											break
										default:
											{
												position688, tokenIndex688, depth688 := position, tokenIndex, depth
												{
													position690 := position
													depth++
													{
														position691, tokenIndex691, depth691 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l692
														}
														position++
														goto l691
													l692:
														position, tokenIndex, depth = position691, tokenIndex691, depth691
														if buffer[position] != rune('S') {
															goto l689
														}
														position++
													}
												l691:
													{
														position693, tokenIndex693, depth693 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l694
														}
														position++
														goto l693
													l694:
														position, tokenIndex, depth = position693, tokenIndex693, depth693
														if buffer[position] != rune('U') {
															goto l689
														}
														position++
													}
												l693:
													{
														position695, tokenIndex695, depth695 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l696
														}
														position++
														goto l695
													l696:
														position, tokenIndex, depth = position695, tokenIndex695, depth695
														if buffer[position] != rune('M') {
															goto l689
														}
														position++
													}
												l695:
													if !rules[ruleskip]() {
														goto l689
													}
													depth--
													add(ruleSUM, position690)
												}
												goto l688
											l689:
												position, tokenIndex, depth = position688, tokenIndex688, depth688
												{
													position698 := position
													depth++
													{
														position699, tokenIndex699, depth699 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l700
														}
														position++
														goto l699
													l700:
														position, tokenIndex, depth = position699, tokenIndex699, depth699
														if buffer[position] != rune('M') {
															goto l697
														}
														position++
													}
												l699:
													{
														position701, tokenIndex701, depth701 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l702
														}
														position++
														goto l701
													l702:
														position, tokenIndex, depth = position701, tokenIndex701, depth701
														if buffer[position] != rune('I') {
															goto l697
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
															goto l697
														}
														position++
													}
												l703:
													if !rules[ruleskip]() {
														goto l697
													}
													depth--
													add(ruleMIN, position698)
												}
												goto l688
											l697:
												position, tokenIndex, depth = position688, tokenIndex688, depth688
												{
													switch buffer[position] {
													case 'S', 's':
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
																	goto l611
																}
																position++
															}
														l707:
															{
																position709, tokenIndex709, depth709 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l710
																}
																position++
																goto l709
															l710:
																position, tokenIndex, depth = position709, tokenIndex709, depth709
																if buffer[position] != rune('A') {
																	goto l611
																}
																position++
															}
														l709:
															{
																position711, tokenIndex711, depth711 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l712
																}
																position++
																goto l711
															l712:
																position, tokenIndex, depth = position711, tokenIndex711, depth711
																if buffer[position] != rune('M') {
																	goto l611
																}
																position++
															}
														l711:
															{
																position713, tokenIndex713, depth713 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l714
																}
																position++
																goto l713
															l714:
																position, tokenIndex, depth = position713, tokenIndex713, depth713
																if buffer[position] != rune('P') {
																	goto l611
																}
																position++
															}
														l713:
															{
																position715, tokenIndex715, depth715 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l716
																}
																position++
																goto l715
															l716:
																position, tokenIndex, depth = position715, tokenIndex715, depth715
																if buffer[position] != rune('L') {
																	goto l611
																}
																position++
															}
														l715:
															{
																position717, tokenIndex717, depth717 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l718
																}
																position++
																goto l717
															l718:
																position, tokenIndex, depth = position717, tokenIndex717, depth717
																if buffer[position] != rune('E') {
																	goto l611
																}
																position++
															}
														l717:
															if !rules[ruleskip]() {
																goto l611
															}
															depth--
															add(ruleSAMPLE, position706)
														}
														break
													case 'A', 'a':
														{
															position719 := position
															depth++
															{
																position720, tokenIndex720, depth720 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l721
																}
																position++
																goto l720
															l721:
																position, tokenIndex, depth = position720, tokenIndex720, depth720
																if buffer[position] != rune('A') {
																	goto l611
																}
																position++
															}
														l720:
															{
																position722, tokenIndex722, depth722 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l723
																}
																position++
																goto l722
															l723:
																position, tokenIndex, depth = position722, tokenIndex722, depth722
																if buffer[position] != rune('V') {
																	goto l611
																}
																position++
															}
														l722:
															{
																position724, tokenIndex724, depth724 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l725
																}
																position++
																goto l724
															l725:
																position, tokenIndex, depth = position724, tokenIndex724, depth724
																if buffer[position] != rune('G') {
																	goto l611
																}
																position++
															}
														l724:
															if !rules[ruleskip]() {
																goto l611
															}
															depth--
															add(ruleAVG, position719)
														}
														break
													default:
														{
															position726 := position
															depth++
															{
																position727, tokenIndex727, depth727 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l728
																}
																position++
																goto l727
															l728:
																position, tokenIndex, depth = position727, tokenIndex727, depth727
																if buffer[position] != rune('M') {
																	goto l611
																}
																position++
															}
														l727:
															{
																position729, tokenIndex729, depth729 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l730
																}
																position++
																goto l729
															l730:
																position, tokenIndex, depth = position729, tokenIndex729, depth729
																if buffer[position] != rune('A') {
																	goto l611
																}
																position++
															}
														l729:
															{
																position731, tokenIndex731, depth731 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l732
																}
																position++
																goto l731
															l732:
																position, tokenIndex, depth = position731, tokenIndex731, depth731
																if buffer[position] != rune('X') {
																	goto l611
																}
																position++
															}
														l731:
															if !rules[ruleskip]() {
																goto l611
															}
															depth--
															add(ruleMAX, position726)
														}
														break
													}
												}

											}
										l688:
											if !rules[ruleLPAREN]() {
												goto l611
											}
											{
												position733, tokenIndex733, depth733 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l733
												}
												goto l734
											l733:
												position, tokenIndex, depth = position733, tokenIndex733, depth733
											}
										l734:
											if !rules[ruleexpression]() {
												goto l611
											}
											if !rules[ruleRPAREN]() {
												goto l611
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position623)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l611
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l611
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l611
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l611
								}
								break
							}
						}

					}
				l617:
					depth--
					add(ruleprimaryExpression, position616)
				}
				depth--
				add(ruleunaryExpression, position612)
			}
			return true
		l611:
			position, tokenIndex, depth = position611, tokenIndex611, depth611
			return false
		},
		/* 56 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 57 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position736, tokenIndex736, depth736 := position, tokenIndex, depth
			{
				position737 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l736
				}
				if !rules[ruleexpression]() {
					goto l736
				}
				if !rules[ruleRPAREN]() {
					goto l736
				}
				depth--
				add(rulebrackettedExpression, position737)
			}
			return true
		l736:
			position, tokenIndex, depth = position736, tokenIndex736, depth736
			return false
		},
		/* 58 functionCall <- <(iriref argList)> */
		func() bool {
			position738, tokenIndex738, depth738 := position, tokenIndex, depth
			{
				position739 := position
				depth++
				if !rules[ruleiriref]() {
					goto l738
				}
				if !rules[ruleargList]() {
					goto l738
				}
				depth--
				add(rulefunctionCall, position739)
			}
			return true
		l738:
			position, tokenIndex, depth = position738, tokenIndex738, depth738
			return false
		},
		/* 59 in <- <(IN argList)> */
		nil,
		/* 60 notin <- <(NOTIN argList)> */
		nil,
		/* 61 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position742, tokenIndex742, depth742 := position, tokenIndex, depth
			{
				position743 := position
				depth++
				{
					position744, tokenIndex744, depth744 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l745
					}
					goto l744
				l745:
					position, tokenIndex, depth = position744, tokenIndex744, depth744
					if !rules[ruleLPAREN]() {
						goto l742
					}
					if !rules[ruleexpression]() {
						goto l742
					}
				l746:
					{
						position747, tokenIndex747, depth747 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l747
						}
						if !rules[ruleexpression]() {
							goto l747
						}
						goto l746
					l747:
						position, tokenIndex, depth = position747, tokenIndex747, depth747
					}
					if !rules[ruleRPAREN]() {
						goto l742
					}
				}
			l744:
				depth--
				add(ruleargList, position743)
			}
			return true
		l742:
			position, tokenIndex, depth = position742, tokenIndex742, depth742
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
			position751, tokenIndex751, depth751 := position, tokenIndex, depth
			{
				position752 := position
				depth++
				{
					position753, tokenIndex753, depth753 := position, tokenIndex, depth
					{
						position755, tokenIndex755, depth755 := position, tokenIndex, depth
						{
							position757 := position
							depth++
							{
								position758, tokenIndex758, depth758 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l759
								}
								position++
								goto l758
							l759:
								position, tokenIndex, depth = position758, tokenIndex758, depth758
								if buffer[position] != rune('S') {
									goto l756
								}
								position++
							}
						l758:
							{
								position760, tokenIndex760, depth760 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l761
								}
								position++
								goto l760
							l761:
								position, tokenIndex, depth = position760, tokenIndex760, depth760
								if buffer[position] != rune('T') {
									goto l756
								}
								position++
							}
						l760:
							{
								position762, tokenIndex762, depth762 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l763
								}
								position++
								goto l762
							l763:
								position, tokenIndex, depth = position762, tokenIndex762, depth762
								if buffer[position] != rune('R') {
									goto l756
								}
								position++
							}
						l762:
							if !rules[ruleskip]() {
								goto l756
							}
							depth--
							add(ruleSTR, position757)
						}
						goto l755
					l756:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position765 := position
							depth++
							{
								position766, tokenIndex766, depth766 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l767
								}
								position++
								goto l766
							l767:
								position, tokenIndex, depth = position766, tokenIndex766, depth766
								if buffer[position] != rune('L') {
									goto l764
								}
								position++
							}
						l766:
							{
								position768, tokenIndex768, depth768 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l769
								}
								position++
								goto l768
							l769:
								position, tokenIndex, depth = position768, tokenIndex768, depth768
								if buffer[position] != rune('A') {
									goto l764
								}
								position++
							}
						l768:
							{
								position770, tokenIndex770, depth770 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l771
								}
								position++
								goto l770
							l771:
								position, tokenIndex, depth = position770, tokenIndex770, depth770
								if buffer[position] != rune('N') {
									goto l764
								}
								position++
							}
						l770:
							{
								position772, tokenIndex772, depth772 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l773
								}
								position++
								goto l772
							l773:
								position, tokenIndex, depth = position772, tokenIndex772, depth772
								if buffer[position] != rune('G') {
									goto l764
								}
								position++
							}
						l772:
							if !rules[ruleskip]() {
								goto l764
							}
							depth--
							add(ruleLANG, position765)
						}
						goto l755
					l764:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position775 := position
							depth++
							{
								position776, tokenIndex776, depth776 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l777
								}
								position++
								goto l776
							l777:
								position, tokenIndex, depth = position776, tokenIndex776, depth776
								if buffer[position] != rune('D') {
									goto l774
								}
								position++
							}
						l776:
							{
								position778, tokenIndex778, depth778 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l779
								}
								position++
								goto l778
							l779:
								position, tokenIndex, depth = position778, tokenIndex778, depth778
								if buffer[position] != rune('A') {
									goto l774
								}
								position++
							}
						l778:
							{
								position780, tokenIndex780, depth780 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l781
								}
								position++
								goto l780
							l781:
								position, tokenIndex, depth = position780, tokenIndex780, depth780
								if buffer[position] != rune('T') {
									goto l774
								}
								position++
							}
						l780:
							{
								position782, tokenIndex782, depth782 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l783
								}
								position++
								goto l782
							l783:
								position, tokenIndex, depth = position782, tokenIndex782, depth782
								if buffer[position] != rune('A') {
									goto l774
								}
								position++
							}
						l782:
							{
								position784, tokenIndex784, depth784 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l785
								}
								position++
								goto l784
							l785:
								position, tokenIndex, depth = position784, tokenIndex784, depth784
								if buffer[position] != rune('T') {
									goto l774
								}
								position++
							}
						l784:
							{
								position786, tokenIndex786, depth786 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l787
								}
								position++
								goto l786
							l787:
								position, tokenIndex, depth = position786, tokenIndex786, depth786
								if buffer[position] != rune('Y') {
									goto l774
								}
								position++
							}
						l786:
							{
								position788, tokenIndex788, depth788 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l789
								}
								position++
								goto l788
							l789:
								position, tokenIndex, depth = position788, tokenIndex788, depth788
								if buffer[position] != rune('P') {
									goto l774
								}
								position++
							}
						l788:
							{
								position790, tokenIndex790, depth790 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l791
								}
								position++
								goto l790
							l791:
								position, tokenIndex, depth = position790, tokenIndex790, depth790
								if buffer[position] != rune('E') {
									goto l774
								}
								position++
							}
						l790:
							if !rules[ruleskip]() {
								goto l774
							}
							depth--
							add(ruleDATATYPE, position775)
						}
						goto l755
					l774:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position793 := position
							depth++
							{
								position794, tokenIndex794, depth794 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l795
								}
								position++
								goto l794
							l795:
								position, tokenIndex, depth = position794, tokenIndex794, depth794
								if buffer[position] != rune('I') {
									goto l792
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
									goto l792
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
									goto l792
								}
								position++
							}
						l798:
							if !rules[ruleskip]() {
								goto l792
							}
							depth--
							add(ruleIRI, position793)
						}
						goto l755
					l792:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position801 := position
							depth++
							{
								position802, tokenIndex802, depth802 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l803
								}
								position++
								goto l802
							l803:
								position, tokenIndex, depth = position802, tokenIndex802, depth802
								if buffer[position] != rune('U') {
									goto l800
								}
								position++
							}
						l802:
							{
								position804, tokenIndex804, depth804 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l805
								}
								position++
								goto l804
							l805:
								position, tokenIndex, depth = position804, tokenIndex804, depth804
								if buffer[position] != rune('R') {
									goto l800
								}
								position++
							}
						l804:
							{
								position806, tokenIndex806, depth806 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l807
								}
								position++
								goto l806
							l807:
								position, tokenIndex, depth = position806, tokenIndex806, depth806
								if buffer[position] != rune('I') {
									goto l800
								}
								position++
							}
						l806:
							if !rules[ruleskip]() {
								goto l800
							}
							depth--
							add(ruleURI, position801)
						}
						goto l755
					l800:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
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
									goto l808
								}
								position++
							}
						l810:
							{
								position812, tokenIndex812, depth812 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l813
								}
								position++
								goto l812
							l813:
								position, tokenIndex, depth = position812, tokenIndex812, depth812
								if buffer[position] != rune('T') {
									goto l808
								}
								position++
							}
						l812:
							{
								position814, tokenIndex814, depth814 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l815
								}
								position++
								goto l814
							l815:
								position, tokenIndex, depth = position814, tokenIndex814, depth814
								if buffer[position] != rune('R') {
									goto l808
								}
								position++
							}
						l814:
							{
								position816, tokenIndex816, depth816 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l817
								}
								position++
								goto l816
							l817:
								position, tokenIndex, depth = position816, tokenIndex816, depth816
								if buffer[position] != rune('L') {
									goto l808
								}
								position++
							}
						l816:
							{
								position818, tokenIndex818, depth818 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l819
								}
								position++
								goto l818
							l819:
								position, tokenIndex, depth = position818, tokenIndex818, depth818
								if buffer[position] != rune('E') {
									goto l808
								}
								position++
							}
						l818:
							{
								position820, tokenIndex820, depth820 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l821
								}
								position++
								goto l820
							l821:
								position, tokenIndex, depth = position820, tokenIndex820, depth820
								if buffer[position] != rune('N') {
									goto l808
								}
								position++
							}
						l820:
							if !rules[ruleskip]() {
								goto l808
							}
							depth--
							add(ruleSTRLEN, position809)
						}
						goto l755
					l808:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position823 := position
							depth++
							{
								position824, tokenIndex824, depth824 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l825
								}
								position++
								goto l824
							l825:
								position, tokenIndex, depth = position824, tokenIndex824, depth824
								if buffer[position] != rune('M') {
									goto l822
								}
								position++
							}
						l824:
							{
								position826, tokenIndex826, depth826 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l827
								}
								position++
								goto l826
							l827:
								position, tokenIndex, depth = position826, tokenIndex826, depth826
								if buffer[position] != rune('O') {
									goto l822
								}
								position++
							}
						l826:
							{
								position828, tokenIndex828, depth828 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l829
								}
								position++
								goto l828
							l829:
								position, tokenIndex, depth = position828, tokenIndex828, depth828
								if buffer[position] != rune('N') {
									goto l822
								}
								position++
							}
						l828:
							{
								position830, tokenIndex830, depth830 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l831
								}
								position++
								goto l830
							l831:
								position, tokenIndex, depth = position830, tokenIndex830, depth830
								if buffer[position] != rune('T') {
									goto l822
								}
								position++
							}
						l830:
							{
								position832, tokenIndex832, depth832 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l833
								}
								position++
								goto l832
							l833:
								position, tokenIndex, depth = position832, tokenIndex832, depth832
								if buffer[position] != rune('H') {
									goto l822
								}
								position++
							}
						l832:
							if !rules[ruleskip]() {
								goto l822
							}
							depth--
							add(ruleMONTH, position823)
						}
						goto l755
					l822:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position835 := position
							depth++
							{
								position836, tokenIndex836, depth836 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l837
								}
								position++
								goto l836
							l837:
								position, tokenIndex, depth = position836, tokenIndex836, depth836
								if buffer[position] != rune('M') {
									goto l834
								}
								position++
							}
						l836:
							{
								position838, tokenIndex838, depth838 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l839
								}
								position++
								goto l838
							l839:
								position, tokenIndex, depth = position838, tokenIndex838, depth838
								if buffer[position] != rune('I') {
									goto l834
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
									goto l834
								}
								position++
							}
						l840:
							{
								position842, tokenIndex842, depth842 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l843
								}
								position++
								goto l842
							l843:
								position, tokenIndex, depth = position842, tokenIndex842, depth842
								if buffer[position] != rune('U') {
									goto l834
								}
								position++
							}
						l842:
							{
								position844, tokenIndex844, depth844 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l845
								}
								position++
								goto l844
							l845:
								position, tokenIndex, depth = position844, tokenIndex844, depth844
								if buffer[position] != rune('T') {
									goto l834
								}
								position++
							}
						l844:
							{
								position846, tokenIndex846, depth846 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l847
								}
								position++
								goto l846
							l847:
								position, tokenIndex, depth = position846, tokenIndex846, depth846
								if buffer[position] != rune('E') {
									goto l834
								}
								position++
							}
						l846:
							{
								position848, tokenIndex848, depth848 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l849
								}
								position++
								goto l848
							l849:
								position, tokenIndex, depth = position848, tokenIndex848, depth848
								if buffer[position] != rune('S') {
									goto l834
								}
								position++
							}
						l848:
							if !rules[ruleskip]() {
								goto l834
							}
							depth--
							add(ruleMINUTES, position835)
						}
						goto l755
					l834:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position851 := position
							depth++
							{
								position852, tokenIndex852, depth852 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l853
								}
								position++
								goto l852
							l853:
								position, tokenIndex, depth = position852, tokenIndex852, depth852
								if buffer[position] != rune('S') {
									goto l850
								}
								position++
							}
						l852:
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
									goto l850
								}
								position++
							}
						l854:
							{
								position856, tokenIndex856, depth856 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l857
								}
								position++
								goto l856
							l857:
								position, tokenIndex, depth = position856, tokenIndex856, depth856
								if buffer[position] != rune('C') {
									goto l850
								}
								position++
							}
						l856:
							{
								position858, tokenIndex858, depth858 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('O') {
									goto l850
								}
								position++
							}
						l858:
							{
								position860, tokenIndex860, depth860 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l861
								}
								position++
								goto l860
							l861:
								position, tokenIndex, depth = position860, tokenIndex860, depth860
								if buffer[position] != rune('N') {
									goto l850
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
									goto l850
								}
								position++
							}
						l862:
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
									goto l850
								}
								position++
							}
						l864:
							if !rules[ruleskip]() {
								goto l850
							}
							depth--
							add(ruleSECONDS, position851)
						}
						goto l755
					l850:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position867 := position
							depth++
							{
								position868, tokenIndex868, depth868 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l869
								}
								position++
								goto l868
							l869:
								position, tokenIndex, depth = position868, tokenIndex868, depth868
								if buffer[position] != rune('T') {
									goto l866
								}
								position++
							}
						l868:
							{
								position870, tokenIndex870, depth870 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l871
								}
								position++
								goto l870
							l871:
								position, tokenIndex, depth = position870, tokenIndex870, depth870
								if buffer[position] != rune('I') {
									goto l866
								}
								position++
							}
						l870:
							{
								position872, tokenIndex872, depth872 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l873
								}
								position++
								goto l872
							l873:
								position, tokenIndex, depth = position872, tokenIndex872, depth872
								if buffer[position] != rune('M') {
									goto l866
								}
								position++
							}
						l872:
							{
								position874, tokenIndex874, depth874 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l875
								}
								position++
								goto l874
							l875:
								position, tokenIndex, depth = position874, tokenIndex874, depth874
								if buffer[position] != rune('E') {
									goto l866
								}
								position++
							}
						l874:
							{
								position876, tokenIndex876, depth876 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l877
								}
								position++
								goto l876
							l877:
								position, tokenIndex, depth = position876, tokenIndex876, depth876
								if buffer[position] != rune('Z') {
									goto l866
								}
								position++
							}
						l876:
							{
								position878, tokenIndex878, depth878 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l879
								}
								position++
								goto l878
							l879:
								position, tokenIndex, depth = position878, tokenIndex878, depth878
								if buffer[position] != rune('O') {
									goto l866
								}
								position++
							}
						l878:
							{
								position880, tokenIndex880, depth880 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l881
								}
								position++
								goto l880
							l881:
								position, tokenIndex, depth = position880, tokenIndex880, depth880
								if buffer[position] != rune('N') {
									goto l866
								}
								position++
							}
						l880:
							{
								position882, tokenIndex882, depth882 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l883
								}
								position++
								goto l882
							l883:
								position, tokenIndex, depth = position882, tokenIndex882, depth882
								if buffer[position] != rune('E') {
									goto l866
								}
								position++
							}
						l882:
							if !rules[ruleskip]() {
								goto l866
							}
							depth--
							add(ruleTIMEZONE, position867)
						}
						goto l755
					l866:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position885 := position
							depth++
							{
								position886, tokenIndex886, depth886 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l887
								}
								position++
								goto l886
							l887:
								position, tokenIndex, depth = position886, tokenIndex886, depth886
								if buffer[position] != rune('S') {
									goto l884
								}
								position++
							}
						l886:
							{
								position888, tokenIndex888, depth888 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l889
								}
								position++
								goto l888
							l889:
								position, tokenIndex, depth = position888, tokenIndex888, depth888
								if buffer[position] != rune('H') {
									goto l884
								}
								position++
							}
						l888:
							{
								position890, tokenIndex890, depth890 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l891
								}
								position++
								goto l890
							l891:
								position, tokenIndex, depth = position890, tokenIndex890, depth890
								if buffer[position] != rune('A') {
									goto l884
								}
								position++
							}
						l890:
							if buffer[position] != rune('1') {
								goto l884
							}
							position++
							if !rules[ruleskip]() {
								goto l884
							}
							depth--
							add(ruleSHA1, position885)
						}
						goto l755
					l884:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position893 := position
							depth++
							{
								position894, tokenIndex894, depth894 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l895
								}
								position++
								goto l894
							l895:
								position, tokenIndex, depth = position894, tokenIndex894, depth894
								if buffer[position] != rune('S') {
									goto l892
								}
								position++
							}
						l894:
							{
								position896, tokenIndex896, depth896 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l897
								}
								position++
								goto l896
							l897:
								position, tokenIndex, depth = position896, tokenIndex896, depth896
								if buffer[position] != rune('H') {
									goto l892
								}
								position++
							}
						l896:
							{
								position898, tokenIndex898, depth898 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l899
								}
								position++
								goto l898
							l899:
								position, tokenIndex, depth = position898, tokenIndex898, depth898
								if buffer[position] != rune('A') {
									goto l892
								}
								position++
							}
						l898:
							if buffer[position] != rune('2') {
								goto l892
							}
							position++
							if buffer[position] != rune('5') {
								goto l892
							}
							position++
							if buffer[position] != rune('6') {
								goto l892
							}
							position++
							if !rules[ruleskip]() {
								goto l892
							}
							depth--
							add(ruleSHA256, position893)
						}
						goto l755
					l892:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position901 := position
							depth++
							{
								position902, tokenIndex902, depth902 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l903
								}
								position++
								goto l902
							l903:
								position, tokenIndex, depth = position902, tokenIndex902, depth902
								if buffer[position] != rune('S') {
									goto l900
								}
								position++
							}
						l902:
							{
								position904, tokenIndex904, depth904 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l905
								}
								position++
								goto l904
							l905:
								position, tokenIndex, depth = position904, tokenIndex904, depth904
								if buffer[position] != rune('H') {
									goto l900
								}
								position++
							}
						l904:
							{
								position906, tokenIndex906, depth906 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l907
								}
								position++
								goto l906
							l907:
								position, tokenIndex, depth = position906, tokenIndex906, depth906
								if buffer[position] != rune('A') {
									goto l900
								}
								position++
							}
						l906:
							if buffer[position] != rune('3') {
								goto l900
							}
							position++
							if buffer[position] != rune('8') {
								goto l900
							}
							position++
							if buffer[position] != rune('4') {
								goto l900
							}
							position++
							if !rules[ruleskip]() {
								goto l900
							}
							depth--
							add(ruleSHA384, position901)
						}
						goto l755
					l900:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position909 := position
							depth++
							{
								position910, tokenIndex910, depth910 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l911
								}
								position++
								goto l910
							l911:
								position, tokenIndex, depth = position910, tokenIndex910, depth910
								if buffer[position] != rune('I') {
									goto l908
								}
								position++
							}
						l910:
							{
								position912, tokenIndex912, depth912 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l913
								}
								position++
								goto l912
							l913:
								position, tokenIndex, depth = position912, tokenIndex912, depth912
								if buffer[position] != rune('S') {
									goto l908
								}
								position++
							}
						l912:
							{
								position914, tokenIndex914, depth914 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l915
								}
								position++
								goto l914
							l915:
								position, tokenIndex, depth = position914, tokenIndex914, depth914
								if buffer[position] != rune('I') {
									goto l908
								}
								position++
							}
						l914:
							{
								position916, tokenIndex916, depth916 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l917
								}
								position++
								goto l916
							l917:
								position, tokenIndex, depth = position916, tokenIndex916, depth916
								if buffer[position] != rune('R') {
									goto l908
								}
								position++
							}
						l916:
							{
								position918, tokenIndex918, depth918 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l919
								}
								position++
								goto l918
							l919:
								position, tokenIndex, depth = position918, tokenIndex918, depth918
								if buffer[position] != rune('I') {
									goto l908
								}
								position++
							}
						l918:
							if !rules[ruleskip]() {
								goto l908
							}
							depth--
							add(ruleISIRI, position909)
						}
						goto l755
					l908:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position921 := position
							depth++
							{
								position922, tokenIndex922, depth922 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l923
								}
								position++
								goto l922
							l923:
								position, tokenIndex, depth = position922, tokenIndex922, depth922
								if buffer[position] != rune('I') {
									goto l920
								}
								position++
							}
						l922:
							{
								position924, tokenIndex924, depth924 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l925
								}
								position++
								goto l924
							l925:
								position, tokenIndex, depth = position924, tokenIndex924, depth924
								if buffer[position] != rune('S') {
									goto l920
								}
								position++
							}
						l924:
							{
								position926, tokenIndex926, depth926 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l927
								}
								position++
								goto l926
							l927:
								position, tokenIndex, depth = position926, tokenIndex926, depth926
								if buffer[position] != rune('U') {
									goto l920
								}
								position++
							}
						l926:
							{
								position928, tokenIndex928, depth928 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l929
								}
								position++
								goto l928
							l929:
								position, tokenIndex, depth = position928, tokenIndex928, depth928
								if buffer[position] != rune('R') {
									goto l920
								}
								position++
							}
						l928:
							{
								position930, tokenIndex930, depth930 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l931
								}
								position++
								goto l930
							l931:
								position, tokenIndex, depth = position930, tokenIndex930, depth930
								if buffer[position] != rune('I') {
									goto l920
								}
								position++
							}
						l930:
							if !rules[ruleskip]() {
								goto l920
							}
							depth--
							add(ruleISURI, position921)
						}
						goto l755
					l920:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position933 := position
							depth++
							{
								position934, tokenIndex934, depth934 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l935
								}
								position++
								goto l934
							l935:
								position, tokenIndex, depth = position934, tokenIndex934, depth934
								if buffer[position] != rune('I') {
									goto l932
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
									goto l932
								}
								position++
							}
						l936:
							{
								position938, tokenIndex938, depth938 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l939
								}
								position++
								goto l938
							l939:
								position, tokenIndex, depth = position938, tokenIndex938, depth938
								if buffer[position] != rune('B') {
									goto l932
								}
								position++
							}
						l938:
							{
								position940, tokenIndex940, depth940 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l941
								}
								position++
								goto l940
							l941:
								position, tokenIndex, depth = position940, tokenIndex940, depth940
								if buffer[position] != rune('L') {
									goto l932
								}
								position++
							}
						l940:
							{
								position942, tokenIndex942, depth942 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l943
								}
								position++
								goto l942
							l943:
								position, tokenIndex, depth = position942, tokenIndex942, depth942
								if buffer[position] != rune('A') {
									goto l932
								}
								position++
							}
						l942:
							{
								position944, tokenIndex944, depth944 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l945
								}
								position++
								goto l944
							l945:
								position, tokenIndex, depth = position944, tokenIndex944, depth944
								if buffer[position] != rune('N') {
									goto l932
								}
								position++
							}
						l944:
							{
								position946, tokenIndex946, depth946 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l947
								}
								position++
								goto l946
							l947:
								position, tokenIndex, depth = position946, tokenIndex946, depth946
								if buffer[position] != rune('K') {
									goto l932
								}
								position++
							}
						l946:
							if !rules[ruleskip]() {
								goto l932
							}
							depth--
							add(ruleISBLANK, position933)
						}
						goto l755
					l932:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							position949 := position
							depth++
							{
								position950, tokenIndex950, depth950 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l951
								}
								position++
								goto l950
							l951:
								position, tokenIndex, depth = position950, tokenIndex950, depth950
								if buffer[position] != rune('I') {
									goto l948
								}
								position++
							}
						l950:
							{
								position952, tokenIndex952, depth952 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l953
								}
								position++
								goto l952
							l953:
								position, tokenIndex, depth = position952, tokenIndex952, depth952
								if buffer[position] != rune('S') {
									goto l948
								}
								position++
							}
						l952:
							{
								position954, tokenIndex954, depth954 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l955
								}
								position++
								goto l954
							l955:
								position, tokenIndex, depth = position954, tokenIndex954, depth954
								if buffer[position] != rune('L') {
									goto l948
								}
								position++
							}
						l954:
							{
								position956, tokenIndex956, depth956 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l957
								}
								position++
								goto l956
							l957:
								position, tokenIndex, depth = position956, tokenIndex956, depth956
								if buffer[position] != rune('I') {
									goto l948
								}
								position++
							}
						l956:
							{
								position958, tokenIndex958, depth958 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l959
								}
								position++
								goto l958
							l959:
								position, tokenIndex, depth = position958, tokenIndex958, depth958
								if buffer[position] != rune('T') {
									goto l948
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
									goto l948
								}
								position++
							}
						l960:
							{
								position962, tokenIndex962, depth962 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l963
								}
								position++
								goto l962
							l963:
								position, tokenIndex, depth = position962, tokenIndex962, depth962
								if buffer[position] != rune('R') {
									goto l948
								}
								position++
							}
						l962:
							{
								position964, tokenIndex964, depth964 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l965
								}
								position++
								goto l964
							l965:
								position, tokenIndex, depth = position964, tokenIndex964, depth964
								if buffer[position] != rune('A') {
									goto l948
								}
								position++
							}
						l964:
							{
								position966, tokenIndex966, depth966 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l967
								}
								position++
								goto l966
							l967:
								position, tokenIndex, depth = position966, tokenIndex966, depth966
								if buffer[position] != rune('L') {
									goto l948
								}
								position++
							}
						l966:
							if !rules[ruleskip]() {
								goto l948
							}
							depth--
							add(ruleISLITERAL, position949)
						}
						goto l755
					l948:
						position, tokenIndex, depth = position755, tokenIndex755, depth755
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position969 := position
									depth++
									{
										position970, tokenIndex970, depth970 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l971
										}
										position++
										goto l970
									l971:
										position, tokenIndex, depth = position970, tokenIndex970, depth970
										if buffer[position] != rune('I') {
											goto l754
										}
										position++
									}
								l970:
									{
										position972, tokenIndex972, depth972 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l973
										}
										position++
										goto l972
									l973:
										position, tokenIndex, depth = position972, tokenIndex972, depth972
										if buffer[position] != rune('S') {
											goto l754
										}
										position++
									}
								l972:
									{
										position974, tokenIndex974, depth974 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l975
										}
										position++
										goto l974
									l975:
										position, tokenIndex, depth = position974, tokenIndex974, depth974
										if buffer[position] != rune('N') {
											goto l754
										}
										position++
									}
								l974:
									{
										position976, tokenIndex976, depth976 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l977
										}
										position++
										goto l976
									l977:
										position, tokenIndex, depth = position976, tokenIndex976, depth976
										if buffer[position] != rune('U') {
											goto l754
										}
										position++
									}
								l976:
									{
										position978, tokenIndex978, depth978 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l979
										}
										position++
										goto l978
									l979:
										position, tokenIndex, depth = position978, tokenIndex978, depth978
										if buffer[position] != rune('M') {
											goto l754
										}
										position++
									}
								l978:
									{
										position980, tokenIndex980, depth980 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l981
										}
										position++
										goto l980
									l981:
										position, tokenIndex, depth = position980, tokenIndex980, depth980
										if buffer[position] != rune('E') {
											goto l754
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
											goto l754
										}
										position++
									}
								l982:
									{
										position984, tokenIndex984, depth984 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l985
										}
										position++
										goto l984
									l985:
										position, tokenIndex, depth = position984, tokenIndex984, depth984
										if buffer[position] != rune('I') {
											goto l754
										}
										position++
									}
								l984:
									{
										position986, tokenIndex986, depth986 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l987
										}
										position++
										goto l986
									l987:
										position, tokenIndex, depth = position986, tokenIndex986, depth986
										if buffer[position] != rune('C') {
											goto l754
										}
										position++
									}
								l986:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleISNUMERIC, position969)
								}
								break
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
											goto l754
										}
										position++
									}
								l989:
									{
										position991, tokenIndex991, depth991 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l992
										}
										position++
										goto l991
									l992:
										position, tokenIndex, depth = position991, tokenIndex991, depth991
										if buffer[position] != rune('H') {
											goto l754
										}
										position++
									}
								l991:
									{
										position993, tokenIndex993, depth993 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l994
										}
										position++
										goto l993
									l994:
										position, tokenIndex, depth = position993, tokenIndex993, depth993
										if buffer[position] != rune('A') {
											goto l754
										}
										position++
									}
								l993:
									if buffer[position] != rune('5') {
										goto l754
									}
									position++
									if buffer[position] != rune('1') {
										goto l754
									}
									position++
									if buffer[position] != rune('2') {
										goto l754
									}
									position++
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleSHA512, position988)
								}
								break
							case 'M', 'm':
								{
									position995 := position
									depth++
									{
										position996, tokenIndex996, depth996 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l997
										}
										position++
										goto l996
									l997:
										position, tokenIndex, depth = position996, tokenIndex996, depth996
										if buffer[position] != rune('M') {
											goto l754
										}
										position++
									}
								l996:
									{
										position998, tokenIndex998, depth998 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l999
										}
										position++
										goto l998
									l999:
										position, tokenIndex, depth = position998, tokenIndex998, depth998
										if buffer[position] != rune('D') {
											goto l754
										}
										position++
									}
								l998:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleMD5, position995)
								}
								break
							case 'T', 't':
								{
									position1000 := position
									depth++
									{
										position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1002
										}
										position++
										goto l1001
									l1002:
										position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
										if buffer[position] != rune('T') {
											goto l754
										}
										position++
									}
								l1001:
									{
										position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1004
										}
										position++
										goto l1003
									l1004:
										position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
										if buffer[position] != rune('Z') {
											goto l754
										}
										position++
									}
								l1003:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleTZ, position1000)
								}
								break
							case 'H', 'h':
								{
									position1005 := position
									depth++
									{
										position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1007
										}
										position++
										goto l1006
									l1007:
										position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
										if buffer[position] != rune('H') {
											goto l754
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
											goto l754
										}
										position++
									}
								l1008:
									{
										position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1011
										}
										position++
										goto l1010
									l1011:
										position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
										if buffer[position] != rune('U') {
											goto l754
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
											goto l754
										}
										position++
									}
								l1012:
									{
										position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1015
										}
										position++
										goto l1014
									l1015:
										position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
										if buffer[position] != rune('S') {
											goto l754
										}
										position++
									}
								l1014:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleHOURS, position1005)
								}
								break
							case 'D', 'd':
								{
									position1016 := position
									depth++
									{
										position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1018
										}
										position++
										goto l1017
									l1018:
										position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
										if buffer[position] != rune('D') {
											goto l754
										}
										position++
									}
								l1017:
									{
										position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1020
										}
										position++
										goto l1019
									l1020:
										position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
										if buffer[position] != rune('A') {
											goto l754
										}
										position++
									}
								l1019:
									{
										position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1022
										}
										position++
										goto l1021
									l1022:
										position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
										if buffer[position] != rune('Y') {
											goto l754
										}
										position++
									}
								l1021:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleDAY, position1016)
								}
								break
							case 'Y', 'y':
								{
									position1023 := position
									depth++
									{
										position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1025
										}
										position++
										goto l1024
									l1025:
										position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
										if buffer[position] != rune('Y') {
											goto l754
										}
										position++
									}
								l1024:
									{
										position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1027
										}
										position++
										goto l1026
									l1027:
										position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
										if buffer[position] != rune('E') {
											goto l754
										}
										position++
									}
								l1026:
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
											goto l754
										}
										position++
									}
								l1028:
									{
										position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1031
										}
										position++
										goto l1030
									l1031:
										position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
										if buffer[position] != rune('R') {
											goto l754
										}
										position++
									}
								l1030:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleYEAR, position1023)
								}
								break
							case 'E', 'e':
								{
									position1032 := position
									depth++
									{
										position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1034
										}
										position++
										goto l1033
									l1034:
										position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
										if buffer[position] != rune('E') {
											goto l754
										}
										position++
									}
								l1033:
									{
										position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1036
										}
										position++
										goto l1035
									l1036:
										position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
										if buffer[position] != rune('N') {
											goto l754
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
											goto l754
										}
										position++
									}
								l1037:
									{
										position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1040
										}
										position++
										goto l1039
									l1040:
										position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
										if buffer[position] != rune('O') {
											goto l754
										}
										position++
									}
								l1039:
									{
										position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1042
										}
										position++
										goto l1041
									l1042:
										position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
										if buffer[position] != rune('D') {
											goto l754
										}
										position++
									}
								l1041:
									{
										position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1044
										}
										position++
										goto l1043
									l1044:
										position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
										if buffer[position] != rune('E') {
											goto l754
										}
										position++
									}
								l1043:
									if buffer[position] != rune('_') {
										goto l754
									}
									position++
									{
										position1045, tokenIndex1045, depth1045 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1046
										}
										position++
										goto l1045
									l1046:
										position, tokenIndex, depth = position1045, tokenIndex1045, depth1045
										if buffer[position] != rune('F') {
											goto l754
										}
										position++
									}
								l1045:
									{
										position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1048
										}
										position++
										goto l1047
									l1048:
										position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
										if buffer[position] != rune('O') {
											goto l754
										}
										position++
									}
								l1047:
									{
										position1049, tokenIndex1049, depth1049 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1050
										}
										position++
										goto l1049
									l1050:
										position, tokenIndex, depth = position1049, tokenIndex1049, depth1049
										if buffer[position] != rune('R') {
											goto l754
										}
										position++
									}
								l1049:
									if buffer[position] != rune('_') {
										goto l754
									}
									position++
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
											goto l754
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
											goto l754
										}
										position++
									}
								l1053:
									{
										position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1056
										}
										position++
										goto l1055
									l1056:
										position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
										if buffer[position] != rune('I') {
											goto l754
										}
										position++
									}
								l1055:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleENCODEFORURI, position1032)
								}
								break
							case 'L', 'l':
								{
									position1057 := position
									depth++
									{
										position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1059
										}
										position++
										goto l1058
									l1059:
										position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
										if buffer[position] != rune('L') {
											goto l754
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
											goto l754
										}
										position++
									}
								l1060:
									{
										position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1063
										}
										position++
										goto l1062
									l1063:
										position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
										if buffer[position] != rune('A') {
											goto l754
										}
										position++
									}
								l1062:
									{
										position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1065
										}
										position++
										goto l1064
									l1065:
										position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
										if buffer[position] != rune('S') {
											goto l754
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
											goto l754
										}
										position++
									}
								l1066:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleLCASE, position1057)
								}
								break
							case 'U', 'u':
								{
									position1068 := position
									depth++
									{
										position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1070
										}
										position++
										goto l1069
									l1070:
										position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
										if buffer[position] != rune('U') {
											goto l754
										}
										position++
									}
								l1069:
									{
										position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1072
										}
										position++
										goto l1071
									l1072:
										position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
										if buffer[position] != rune('C') {
											goto l754
										}
										position++
									}
								l1071:
									{
										position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1074
										}
										position++
										goto l1073
									l1074:
										position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
										if buffer[position] != rune('A') {
											goto l754
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
											goto l754
										}
										position++
									}
								l1075:
									{
										position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1078
										}
										position++
										goto l1077
									l1078:
										position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
										if buffer[position] != rune('E') {
											goto l754
										}
										position++
									}
								l1077:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleUCASE, position1068)
								}
								break
							case 'F', 'f':
								{
									position1079 := position
									depth++
									{
										position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1081
										}
										position++
										goto l1080
									l1081:
										position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
										if buffer[position] != rune('F') {
											goto l754
										}
										position++
									}
								l1080:
									{
										position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1083
										}
										position++
										goto l1082
									l1083:
										position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
										if buffer[position] != rune('L') {
											goto l754
										}
										position++
									}
								l1082:
									{
										position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1085
										}
										position++
										goto l1084
									l1085:
										position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
										if buffer[position] != rune('O') {
											goto l754
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
											goto l754
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
											goto l754
										}
										position++
									}
								l1088:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleFLOOR, position1079)
								}
								break
							case 'R', 'r':
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
											goto l754
										}
										position++
									}
								l1091:
									{
										position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1094
										}
										position++
										goto l1093
									l1094:
										position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
										if buffer[position] != rune('O') {
											goto l754
										}
										position++
									}
								l1093:
									{
										position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1096
										}
										position++
										goto l1095
									l1096:
										position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
										if buffer[position] != rune('U') {
											goto l754
										}
										position++
									}
								l1095:
									{
										position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1098
										}
										position++
										goto l1097
									l1098:
										position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
										if buffer[position] != rune('N') {
											goto l754
										}
										position++
									}
								l1097:
									{
										position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1100
										}
										position++
										goto l1099
									l1100:
										position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
										if buffer[position] != rune('D') {
											goto l754
										}
										position++
									}
								l1099:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleROUND, position1090)
								}
								break
							case 'C', 'c':
								{
									position1101 := position
									depth++
									{
										position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1103
										}
										position++
										goto l1102
									l1103:
										position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
										if buffer[position] != rune('C') {
											goto l754
										}
										position++
									}
								l1102:
									{
										position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1105
										}
										position++
										goto l1104
									l1105:
										position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
										if buffer[position] != rune('E') {
											goto l754
										}
										position++
									}
								l1104:
									{
										position1106, tokenIndex1106, depth1106 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1107
										}
										position++
										goto l1106
									l1107:
										position, tokenIndex, depth = position1106, tokenIndex1106, depth1106
										if buffer[position] != rune('I') {
											goto l754
										}
										position++
									}
								l1106:
									{
										position1108, tokenIndex1108, depth1108 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1109
										}
										position++
										goto l1108
									l1109:
										position, tokenIndex, depth = position1108, tokenIndex1108, depth1108
										if buffer[position] != rune('L') {
											goto l754
										}
										position++
									}
								l1108:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleCEIL, position1101)
								}
								break
							default:
								{
									position1110 := position
									depth++
									{
										position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1112
										}
										position++
										goto l1111
									l1112:
										position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
										if buffer[position] != rune('A') {
											goto l754
										}
										position++
									}
								l1111:
									{
										position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1114
										}
										position++
										goto l1113
									l1114:
										position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
										if buffer[position] != rune('B') {
											goto l754
										}
										position++
									}
								l1113:
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
											goto l754
										}
										position++
									}
								l1115:
									if !rules[ruleskip]() {
										goto l754
									}
									depth--
									add(ruleABS, position1110)
								}
								break
							}
						}

					}
				l755:
					if !rules[ruleLPAREN]() {
						goto l754
					}
					if !rules[ruleexpression]() {
						goto l754
					}
					if !rules[ruleRPAREN]() {
						goto l754
					}
					goto l753
				l754:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					{
						position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
						{
							position1120 := position
							depth++
							{
								position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1122
								}
								position++
								goto l1121
							l1122:
								position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
								if buffer[position] != rune('S') {
									goto l1119
								}
								position++
							}
						l1121:
							{
								position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1124
								}
								position++
								goto l1123
							l1124:
								position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
								if buffer[position] != rune('T') {
									goto l1119
								}
								position++
							}
						l1123:
							{
								position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1126
								}
								position++
								goto l1125
							l1126:
								position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
								if buffer[position] != rune('R') {
									goto l1119
								}
								position++
							}
						l1125:
							{
								position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1128
								}
								position++
								goto l1127
							l1128:
								position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
								if buffer[position] != rune('S') {
									goto l1119
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
									goto l1119
								}
								position++
							}
						l1129:
							{
								position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1132
								}
								position++
								goto l1131
							l1132:
								position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
								if buffer[position] != rune('A') {
									goto l1119
								}
								position++
							}
						l1131:
							{
								position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1134
								}
								position++
								goto l1133
							l1134:
								position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
								if buffer[position] != rune('R') {
									goto l1119
								}
								position++
							}
						l1133:
							{
								position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1136
								}
								position++
								goto l1135
							l1136:
								position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
								if buffer[position] != rune('T') {
									goto l1119
								}
								position++
							}
						l1135:
							{
								position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1138
								}
								position++
								goto l1137
							l1138:
								position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
								if buffer[position] != rune('S') {
									goto l1119
								}
								position++
							}
						l1137:
							if !rules[ruleskip]() {
								goto l1119
							}
							depth--
							add(ruleSTRSTARTS, position1120)
						}
						goto l1118
					l1119:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						{
							position1140 := position
							depth++
							{
								position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1142
								}
								position++
								goto l1141
							l1142:
								position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
								if buffer[position] != rune('S') {
									goto l1139
								}
								position++
							}
						l1141:
							{
								position1143, tokenIndex1143, depth1143 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1144
								}
								position++
								goto l1143
							l1144:
								position, tokenIndex, depth = position1143, tokenIndex1143, depth1143
								if buffer[position] != rune('T') {
									goto l1139
								}
								position++
							}
						l1143:
							{
								position1145, tokenIndex1145, depth1145 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1146
								}
								position++
								goto l1145
							l1146:
								position, tokenIndex, depth = position1145, tokenIndex1145, depth1145
								if buffer[position] != rune('R') {
									goto l1139
								}
								position++
							}
						l1145:
							{
								position1147, tokenIndex1147, depth1147 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1148
								}
								position++
								goto l1147
							l1148:
								position, tokenIndex, depth = position1147, tokenIndex1147, depth1147
								if buffer[position] != rune('E') {
									goto l1139
								}
								position++
							}
						l1147:
							{
								position1149, tokenIndex1149, depth1149 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1150
								}
								position++
								goto l1149
							l1150:
								position, tokenIndex, depth = position1149, tokenIndex1149, depth1149
								if buffer[position] != rune('N') {
									goto l1139
								}
								position++
							}
						l1149:
							{
								position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1152
								}
								position++
								goto l1151
							l1152:
								position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
								if buffer[position] != rune('D') {
									goto l1139
								}
								position++
							}
						l1151:
							{
								position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1154
								}
								position++
								goto l1153
							l1154:
								position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
								if buffer[position] != rune('S') {
									goto l1139
								}
								position++
							}
						l1153:
							if !rules[ruleskip]() {
								goto l1139
							}
							depth--
							add(ruleSTRENDS, position1140)
						}
						goto l1118
					l1139:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						{
							position1156 := position
							depth++
							{
								position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1158
								}
								position++
								goto l1157
							l1158:
								position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
								if buffer[position] != rune('S') {
									goto l1155
								}
								position++
							}
						l1157:
							{
								position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1160
								}
								position++
								goto l1159
							l1160:
								position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
								if buffer[position] != rune('T') {
									goto l1155
								}
								position++
							}
						l1159:
							{
								position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1162
								}
								position++
								goto l1161
							l1162:
								position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
								if buffer[position] != rune('R') {
									goto l1155
								}
								position++
							}
						l1161:
							{
								position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1164
								}
								position++
								goto l1163
							l1164:
								position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
								if buffer[position] != rune('B') {
									goto l1155
								}
								position++
							}
						l1163:
							{
								position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1166
								}
								position++
								goto l1165
							l1166:
								position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
								if buffer[position] != rune('E') {
									goto l1155
								}
								position++
							}
						l1165:
							{
								position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1168
								}
								position++
								goto l1167
							l1168:
								position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
								if buffer[position] != rune('F') {
									goto l1155
								}
								position++
							}
						l1167:
							{
								position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1170
								}
								position++
								goto l1169
							l1170:
								position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
								if buffer[position] != rune('O') {
									goto l1155
								}
								position++
							}
						l1169:
							{
								position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1172
								}
								position++
								goto l1171
							l1172:
								position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
								if buffer[position] != rune('R') {
									goto l1155
								}
								position++
							}
						l1171:
							{
								position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1174
								}
								position++
								goto l1173
							l1174:
								position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
								if buffer[position] != rune('E') {
									goto l1155
								}
								position++
							}
						l1173:
							if !rules[ruleskip]() {
								goto l1155
							}
							depth--
							add(ruleSTRBEFORE, position1156)
						}
						goto l1118
					l1155:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						{
							position1176 := position
							depth++
							{
								position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1178
								}
								position++
								goto l1177
							l1178:
								position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
								if buffer[position] != rune('S') {
									goto l1175
								}
								position++
							}
						l1177:
							{
								position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1180
								}
								position++
								goto l1179
							l1180:
								position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
								if buffer[position] != rune('T') {
									goto l1175
								}
								position++
							}
						l1179:
							{
								position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1182
								}
								position++
								goto l1181
							l1182:
								position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
								if buffer[position] != rune('R') {
									goto l1175
								}
								position++
							}
						l1181:
							{
								position1183, tokenIndex1183, depth1183 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1184
								}
								position++
								goto l1183
							l1184:
								position, tokenIndex, depth = position1183, tokenIndex1183, depth1183
								if buffer[position] != rune('A') {
									goto l1175
								}
								position++
							}
						l1183:
							{
								position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1186
								}
								position++
								goto l1185
							l1186:
								position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
								if buffer[position] != rune('F') {
									goto l1175
								}
								position++
							}
						l1185:
							{
								position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1188
								}
								position++
								goto l1187
							l1188:
								position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
								if buffer[position] != rune('T') {
									goto l1175
								}
								position++
							}
						l1187:
							{
								position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1190
								}
								position++
								goto l1189
							l1190:
								position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
								if buffer[position] != rune('E') {
									goto l1175
								}
								position++
							}
						l1189:
							{
								position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1192
								}
								position++
								goto l1191
							l1192:
								position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
								if buffer[position] != rune('R') {
									goto l1175
								}
								position++
							}
						l1191:
							if !rules[ruleskip]() {
								goto l1175
							}
							depth--
							add(ruleSTRAFTER, position1176)
						}
						goto l1118
					l1175:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						{
							position1194 := position
							depth++
							{
								position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1196
								}
								position++
								goto l1195
							l1196:
								position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
								if buffer[position] != rune('S') {
									goto l1193
								}
								position++
							}
						l1195:
							{
								position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1198
								}
								position++
								goto l1197
							l1198:
								position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
								if buffer[position] != rune('T') {
									goto l1193
								}
								position++
							}
						l1197:
							{
								position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1200
								}
								position++
								goto l1199
							l1200:
								position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
								if buffer[position] != rune('R') {
									goto l1193
								}
								position++
							}
						l1199:
							{
								position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1202
								}
								position++
								goto l1201
							l1202:
								position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
								if buffer[position] != rune('L') {
									goto l1193
								}
								position++
							}
						l1201:
							{
								position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1204
								}
								position++
								goto l1203
							l1204:
								position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
								if buffer[position] != rune('A') {
									goto l1193
								}
								position++
							}
						l1203:
							{
								position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1206
								}
								position++
								goto l1205
							l1206:
								position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
								if buffer[position] != rune('N') {
									goto l1193
								}
								position++
							}
						l1205:
							{
								position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1208
								}
								position++
								goto l1207
							l1208:
								position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
								if buffer[position] != rune('G') {
									goto l1193
								}
								position++
							}
						l1207:
							if !rules[ruleskip]() {
								goto l1193
							}
							depth--
							add(ruleSTRLANG, position1194)
						}
						goto l1118
					l1193:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
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
									goto l1209
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
									goto l1209
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
									goto l1209
								}
								position++
							}
						l1215:
							{
								position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1218
								}
								position++
								goto l1217
							l1218:
								position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
								if buffer[position] != rune('D') {
									goto l1209
								}
								position++
							}
						l1217:
							{
								position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1220
								}
								position++
								goto l1219
							l1220:
								position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
								if buffer[position] != rune('T') {
									goto l1209
								}
								position++
							}
						l1219:
							if !rules[ruleskip]() {
								goto l1209
							}
							depth--
							add(ruleSTRDT, position1210)
						}
						goto l1118
					l1209:
						position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1222 := position
									depth++
									{
										position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1224
										}
										position++
										goto l1223
									l1224:
										position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
										if buffer[position] != rune('S') {
											goto l1117
										}
										position++
									}
								l1223:
									{
										position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1226
										}
										position++
										goto l1225
									l1226:
										position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
										if buffer[position] != rune('A') {
											goto l1117
										}
										position++
									}
								l1225:
									{
										position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1228
										}
										position++
										goto l1227
									l1228:
										position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
										if buffer[position] != rune('M') {
											goto l1117
										}
										position++
									}
								l1227:
									{
										position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1230
										}
										position++
										goto l1229
									l1230:
										position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
										if buffer[position] != rune('E') {
											goto l1117
										}
										position++
									}
								l1229:
									{
										position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1232
										}
										position++
										goto l1231
									l1232:
										position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
										if buffer[position] != rune('T') {
											goto l1117
										}
										position++
									}
								l1231:
									{
										position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1234
										}
										position++
										goto l1233
									l1234:
										position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
										if buffer[position] != rune('E') {
											goto l1117
										}
										position++
									}
								l1233:
									{
										position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1236
										}
										position++
										goto l1235
									l1236:
										position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
										if buffer[position] != rune('R') {
											goto l1117
										}
										position++
									}
								l1235:
									{
										position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1238
										}
										position++
										goto l1237
									l1238:
										position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
										if buffer[position] != rune('M') {
											goto l1117
										}
										position++
									}
								l1237:
									if !rules[ruleskip]() {
										goto l1117
									}
									depth--
									add(ruleSAMETERM, position1222)
								}
								break
							case 'C', 'c':
								{
									position1239 := position
									depth++
									{
										position1240, tokenIndex1240, depth1240 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1241
										}
										position++
										goto l1240
									l1241:
										position, tokenIndex, depth = position1240, tokenIndex1240, depth1240
										if buffer[position] != rune('C') {
											goto l1117
										}
										position++
									}
								l1240:
									{
										position1242, tokenIndex1242, depth1242 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1243
										}
										position++
										goto l1242
									l1243:
										position, tokenIndex, depth = position1242, tokenIndex1242, depth1242
										if buffer[position] != rune('O') {
											goto l1117
										}
										position++
									}
								l1242:
									{
										position1244, tokenIndex1244, depth1244 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1245
										}
										position++
										goto l1244
									l1245:
										position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
										if buffer[position] != rune('N') {
											goto l1117
										}
										position++
									}
								l1244:
									{
										position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1247
										}
										position++
										goto l1246
									l1247:
										position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
										if buffer[position] != rune('T') {
											goto l1117
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
											goto l1117
										}
										position++
									}
								l1248:
									{
										position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1251
										}
										position++
										goto l1250
									l1251:
										position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
										if buffer[position] != rune('I') {
											goto l1117
										}
										position++
									}
								l1250:
									{
										position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1253
										}
										position++
										goto l1252
									l1253:
										position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
										if buffer[position] != rune('N') {
											goto l1117
										}
										position++
									}
								l1252:
									{
										position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1255
										}
										position++
										goto l1254
									l1255:
										position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
										if buffer[position] != rune('S') {
											goto l1117
										}
										position++
									}
								l1254:
									if !rules[ruleskip]() {
										goto l1117
									}
									depth--
									add(ruleCONTAINS, position1239)
								}
								break
							default:
								{
									position1256 := position
									depth++
									{
										position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1258
										}
										position++
										goto l1257
									l1258:
										position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
										if buffer[position] != rune('L') {
											goto l1117
										}
										position++
									}
								l1257:
									{
										position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1260
										}
										position++
										goto l1259
									l1260:
										position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
										if buffer[position] != rune('A') {
											goto l1117
										}
										position++
									}
								l1259:
									{
										position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1262
										}
										position++
										goto l1261
									l1262:
										position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
										if buffer[position] != rune('N') {
											goto l1117
										}
										position++
									}
								l1261:
									{
										position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1264
										}
										position++
										goto l1263
									l1264:
										position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
										if buffer[position] != rune('G') {
											goto l1117
										}
										position++
									}
								l1263:
									{
										position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1266
										}
										position++
										goto l1265
									l1266:
										position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
										if buffer[position] != rune('M') {
											goto l1117
										}
										position++
									}
								l1265:
									{
										position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1268
										}
										position++
										goto l1267
									l1268:
										position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
										if buffer[position] != rune('A') {
											goto l1117
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
											goto l1117
										}
										position++
									}
								l1269:
									{
										position1271, tokenIndex1271, depth1271 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1272
										}
										position++
										goto l1271
									l1272:
										position, tokenIndex, depth = position1271, tokenIndex1271, depth1271
										if buffer[position] != rune('C') {
											goto l1117
										}
										position++
									}
								l1271:
									{
										position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1274
										}
										position++
										goto l1273
									l1274:
										position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
										if buffer[position] != rune('H') {
											goto l1117
										}
										position++
									}
								l1273:
									{
										position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1276
										}
										position++
										goto l1275
									l1276:
										position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
										if buffer[position] != rune('E') {
											goto l1117
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
											goto l1117
										}
										position++
									}
								l1277:
									if !rules[ruleskip]() {
										goto l1117
									}
									depth--
									add(ruleLANGMATCHES, position1256)
								}
								break
							}
						}

					}
				l1118:
					if !rules[ruleLPAREN]() {
						goto l1117
					}
					if !rules[ruleexpression]() {
						goto l1117
					}
					if !rules[ruleCOMMA]() {
						goto l1117
					}
					if !rules[ruleexpression]() {
						goto l1117
					}
					if !rules[ruleRPAREN]() {
						goto l1117
					}
					goto l753
				l1117:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					{
						position1280 := position
						depth++
						{
							position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1282
							}
							position++
							goto l1281
						l1282:
							position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
							if buffer[position] != rune('B') {
								goto l1279
							}
							position++
						}
					l1281:
						{
							position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1284
							}
							position++
							goto l1283
						l1284:
							position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
							if buffer[position] != rune('O') {
								goto l1279
							}
							position++
						}
					l1283:
						{
							position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1286
							}
							position++
							goto l1285
						l1286:
							position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
							if buffer[position] != rune('U') {
								goto l1279
							}
							position++
						}
					l1285:
						{
							position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1288
							}
							position++
							goto l1287
						l1288:
							position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
							if buffer[position] != rune('N') {
								goto l1279
							}
							position++
						}
					l1287:
						{
							position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1290
							}
							position++
							goto l1289
						l1290:
							position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
							if buffer[position] != rune('D') {
								goto l1279
							}
							position++
						}
					l1289:
						if !rules[ruleskip]() {
							goto l1279
						}
						depth--
						add(ruleBOUND, position1280)
					}
					if !rules[ruleLPAREN]() {
						goto l1279
					}
					if !rules[rulevar]() {
						goto l1279
					}
					if !rules[ruleRPAREN]() {
						goto l1279
					}
					goto l753
				l1279:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1293 := position
								depth++
								{
									position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1295
									}
									position++
									goto l1294
								l1295:
									position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
									if buffer[position] != rune('S') {
										goto l1291
									}
									position++
								}
							l1294:
								{
									position1296, tokenIndex1296, depth1296 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1297
									}
									position++
									goto l1296
								l1297:
									position, tokenIndex, depth = position1296, tokenIndex1296, depth1296
									if buffer[position] != rune('T') {
										goto l1291
									}
									position++
								}
							l1296:
								{
									position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1299
									}
									position++
									goto l1298
								l1299:
									position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
									if buffer[position] != rune('R') {
										goto l1291
									}
									position++
								}
							l1298:
								{
									position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1301
									}
									position++
									goto l1300
								l1301:
									position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
									if buffer[position] != rune('U') {
										goto l1291
									}
									position++
								}
							l1300:
								{
									position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1303
									}
									position++
									goto l1302
								l1303:
									position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
									if buffer[position] != rune('U') {
										goto l1291
									}
									position++
								}
							l1302:
								{
									position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1305
									}
									position++
									goto l1304
								l1305:
									position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
									if buffer[position] != rune('I') {
										goto l1291
									}
									position++
								}
							l1304:
								{
									position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1307
									}
									position++
									goto l1306
								l1307:
									position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
									if buffer[position] != rune('D') {
										goto l1291
									}
									position++
								}
							l1306:
								if !rules[ruleskip]() {
									goto l1291
								}
								depth--
								add(ruleSTRUUID, position1293)
							}
							break
						case 'U', 'u':
							{
								position1308 := position
								depth++
								{
									position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1310
									}
									position++
									goto l1309
								l1310:
									position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
									if buffer[position] != rune('U') {
										goto l1291
									}
									position++
								}
							l1309:
								{
									position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1312
									}
									position++
									goto l1311
								l1312:
									position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
									if buffer[position] != rune('U') {
										goto l1291
									}
									position++
								}
							l1311:
								{
									position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1314
									}
									position++
									goto l1313
								l1314:
									position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
									if buffer[position] != rune('I') {
										goto l1291
									}
									position++
								}
							l1313:
								{
									position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1316
									}
									position++
									goto l1315
								l1316:
									position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
									if buffer[position] != rune('D') {
										goto l1291
									}
									position++
								}
							l1315:
								if !rules[ruleskip]() {
									goto l1291
								}
								depth--
								add(ruleUUID, position1308)
							}
							break
						case 'N', 'n':
							{
								position1317 := position
								depth++
								{
									position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1319
									}
									position++
									goto l1318
								l1319:
									position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
									if buffer[position] != rune('N') {
										goto l1291
									}
									position++
								}
							l1318:
								{
									position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1321
									}
									position++
									goto l1320
								l1321:
									position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
									if buffer[position] != rune('O') {
										goto l1291
									}
									position++
								}
							l1320:
								{
									position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1323
									}
									position++
									goto l1322
								l1323:
									position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
									if buffer[position] != rune('W') {
										goto l1291
									}
									position++
								}
							l1322:
								if !rules[ruleskip]() {
									goto l1291
								}
								depth--
								add(ruleNOW, position1317)
							}
							break
						default:
							{
								position1324 := position
								depth++
								{
									position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1326
									}
									position++
									goto l1325
								l1326:
									position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
									if buffer[position] != rune('R') {
										goto l1291
									}
									position++
								}
							l1325:
								{
									position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1328
									}
									position++
									goto l1327
								l1328:
									position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
									if buffer[position] != rune('A') {
										goto l1291
									}
									position++
								}
							l1327:
								{
									position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1330
									}
									position++
									goto l1329
								l1330:
									position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
									if buffer[position] != rune('N') {
										goto l1291
									}
									position++
								}
							l1329:
								{
									position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1332
									}
									position++
									goto l1331
								l1332:
									position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
									if buffer[position] != rune('D') {
										goto l1291
									}
									position++
								}
							l1331:
								if !rules[ruleskip]() {
									goto l1291
								}
								depth--
								add(ruleRAND, position1324)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1291
					}
					goto l753
				l1291:
					position, tokenIndex, depth = position753, tokenIndex753, depth753
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
								{
									position1336 := position
									depth++
									{
										position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1338
										}
										position++
										goto l1337
									l1338:
										position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
										if buffer[position] != rune('E') {
											goto l1335
										}
										position++
									}
								l1337:
									{
										position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1340
										}
										position++
										goto l1339
									l1340:
										position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
										if buffer[position] != rune('X') {
											goto l1335
										}
										position++
									}
								l1339:
									{
										position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1342
										}
										position++
										goto l1341
									l1342:
										position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
										if buffer[position] != rune('I') {
											goto l1335
										}
										position++
									}
								l1341:
									{
										position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1344
										}
										position++
										goto l1343
									l1344:
										position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
										if buffer[position] != rune('S') {
											goto l1335
										}
										position++
									}
								l1343:
									{
										position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1346
										}
										position++
										goto l1345
									l1346:
										position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
										if buffer[position] != rune('T') {
											goto l1335
										}
										position++
									}
								l1345:
									{
										position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1348
										}
										position++
										goto l1347
									l1348:
										position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
										if buffer[position] != rune('S') {
											goto l1335
										}
										position++
									}
								l1347:
									if !rules[ruleskip]() {
										goto l1335
									}
									depth--
									add(ruleEXISTS, position1336)
								}
								goto l1334
							l1335:
								position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
								{
									position1349 := position
									depth++
									{
										position1350, tokenIndex1350, depth1350 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1351
										}
										position++
										goto l1350
									l1351:
										position, tokenIndex, depth = position1350, tokenIndex1350, depth1350
										if buffer[position] != rune('N') {
											goto l751
										}
										position++
									}
								l1350:
									{
										position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1353
										}
										position++
										goto l1352
									l1353:
										position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
										if buffer[position] != rune('O') {
											goto l751
										}
										position++
									}
								l1352:
									{
										position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1355
										}
										position++
										goto l1354
									l1355:
										position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
										if buffer[position] != rune('T') {
											goto l751
										}
										position++
									}
								l1354:
									if buffer[position] != rune(' ') {
										goto l751
									}
									position++
									{
										position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1357
										}
										position++
										goto l1356
									l1357:
										position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
										if buffer[position] != rune('E') {
											goto l751
										}
										position++
									}
								l1356:
									{
										position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1359
										}
										position++
										goto l1358
									l1359:
										position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
										if buffer[position] != rune('X') {
											goto l751
										}
										position++
									}
								l1358:
									{
										position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1361
										}
										position++
										goto l1360
									l1361:
										position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
										if buffer[position] != rune('I') {
											goto l751
										}
										position++
									}
								l1360:
									{
										position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1363
										}
										position++
										goto l1362
									l1363:
										position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
										if buffer[position] != rune('S') {
											goto l751
										}
										position++
									}
								l1362:
									{
										position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1365
										}
										position++
										goto l1364
									l1365:
										position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
										if buffer[position] != rune('T') {
											goto l751
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
											goto l751
										}
										position++
									}
								l1366:
									if !rules[ruleskip]() {
										goto l751
									}
									depth--
									add(ruleNOTEXIST, position1349)
								}
							}
						l1334:
							if !rules[rulegroupGraphPattern]() {
								goto l751
							}
							break
						case 'I', 'i':
							{
								position1368 := position
								depth++
								{
									position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1370
									}
									position++
									goto l1369
								l1370:
									position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
									if buffer[position] != rune('I') {
										goto l751
									}
									position++
								}
							l1369:
								{
									position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1372
									}
									position++
									goto l1371
								l1372:
									position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
									if buffer[position] != rune('F') {
										goto l751
									}
									position++
								}
							l1371:
								if !rules[ruleskip]() {
									goto l751
								}
								depth--
								add(ruleIF, position1368)
							}
							if !rules[ruleLPAREN]() {
								goto l751
							}
							if !rules[ruleexpression]() {
								goto l751
							}
							if !rules[ruleCOMMA]() {
								goto l751
							}
							if !rules[ruleexpression]() {
								goto l751
							}
							if !rules[ruleCOMMA]() {
								goto l751
							}
							if !rules[ruleexpression]() {
								goto l751
							}
							if !rules[ruleRPAREN]() {
								goto l751
							}
							break
						case 'C', 'c':
							{
								position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
								{
									position1375 := position
									depth++
									{
										position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1377
										}
										position++
										goto l1376
									l1377:
										position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
										if buffer[position] != rune('C') {
											goto l1374
										}
										position++
									}
								l1376:
									{
										position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1379
										}
										position++
										goto l1378
									l1379:
										position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
										if buffer[position] != rune('O') {
											goto l1374
										}
										position++
									}
								l1378:
									{
										position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1381
										}
										position++
										goto l1380
									l1381:
										position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
										if buffer[position] != rune('N') {
											goto l1374
										}
										position++
									}
								l1380:
									{
										position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1383
										}
										position++
										goto l1382
									l1383:
										position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
										if buffer[position] != rune('C') {
											goto l1374
										}
										position++
									}
								l1382:
									{
										position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1385
										}
										position++
										goto l1384
									l1385:
										position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
										if buffer[position] != rune('A') {
											goto l1374
										}
										position++
									}
								l1384:
									{
										position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1387
										}
										position++
										goto l1386
									l1387:
										position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
										if buffer[position] != rune('T') {
											goto l1374
										}
										position++
									}
								l1386:
									if !rules[ruleskip]() {
										goto l1374
									}
									depth--
									add(ruleCONCAT, position1375)
								}
								goto l1373
							l1374:
								position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
								{
									position1388 := position
									depth++
									{
										position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1390
										}
										position++
										goto l1389
									l1390:
										position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
										if buffer[position] != rune('C') {
											goto l751
										}
										position++
									}
								l1389:
									{
										position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1392
										}
										position++
										goto l1391
									l1392:
										position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
										if buffer[position] != rune('O') {
											goto l751
										}
										position++
									}
								l1391:
									{
										position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1394
										}
										position++
										goto l1393
									l1394:
										position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
										if buffer[position] != rune('A') {
											goto l751
										}
										position++
									}
								l1393:
									{
										position1395, tokenIndex1395, depth1395 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1396
										}
										position++
										goto l1395
									l1396:
										position, tokenIndex, depth = position1395, tokenIndex1395, depth1395
										if buffer[position] != rune('L') {
											goto l751
										}
										position++
									}
								l1395:
									{
										position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1398
										}
										position++
										goto l1397
									l1398:
										position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
										if buffer[position] != rune('E') {
											goto l751
										}
										position++
									}
								l1397:
									{
										position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1400
										}
										position++
										goto l1399
									l1400:
										position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
										if buffer[position] != rune('S') {
											goto l751
										}
										position++
									}
								l1399:
									{
										position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1402
										}
										position++
										goto l1401
									l1402:
										position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
										if buffer[position] != rune('C') {
											goto l751
										}
										position++
									}
								l1401:
									{
										position1403, tokenIndex1403, depth1403 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1404
										}
										position++
										goto l1403
									l1404:
										position, tokenIndex, depth = position1403, tokenIndex1403, depth1403
										if buffer[position] != rune('E') {
											goto l751
										}
										position++
									}
								l1403:
									if !rules[ruleskip]() {
										goto l751
									}
									depth--
									add(ruleCOALESCE, position1388)
								}
							}
						l1373:
							if !rules[ruleargList]() {
								goto l751
							}
							break
						case 'B', 'b':
							{
								position1405 := position
								depth++
								{
									position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1407
									}
									position++
									goto l1406
								l1407:
									position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
									if buffer[position] != rune('B') {
										goto l751
									}
									position++
								}
							l1406:
								{
									position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1409
									}
									position++
									goto l1408
								l1409:
									position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
									if buffer[position] != rune('N') {
										goto l751
									}
									position++
								}
							l1408:
								{
									position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1411
									}
									position++
									goto l1410
								l1411:
									position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
									if buffer[position] != rune('O') {
										goto l751
									}
									position++
								}
							l1410:
								{
									position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1413
									}
									position++
									goto l1412
								l1413:
									position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
									if buffer[position] != rune('D') {
										goto l751
									}
									position++
								}
							l1412:
								{
									position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1415
									}
									position++
									goto l1414
								l1415:
									position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
									if buffer[position] != rune('E') {
										goto l751
									}
									position++
								}
							l1414:
								if !rules[ruleskip]() {
									goto l751
								}
								depth--
								add(ruleBNODE, position1405)
							}
							{
								position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1417
								}
								if !rules[ruleexpression]() {
									goto l1417
								}
								if !rules[ruleRPAREN]() {
									goto l1417
								}
								goto l1416
							l1417:
								position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
								if !rules[rulenil]() {
									goto l751
								}
							}
						l1416:
							break
						default:
							{
								position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
								{
									position1420 := position
									depth++
									{
										position1421, tokenIndex1421, depth1421 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1422
										}
										position++
										goto l1421
									l1422:
										position, tokenIndex, depth = position1421, tokenIndex1421, depth1421
										if buffer[position] != rune('S') {
											goto l1419
										}
										position++
									}
								l1421:
									{
										position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1424
										}
										position++
										goto l1423
									l1424:
										position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
										if buffer[position] != rune('U') {
											goto l1419
										}
										position++
									}
								l1423:
									{
										position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1426
										}
										position++
										goto l1425
									l1426:
										position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
										if buffer[position] != rune('B') {
											goto l1419
										}
										position++
									}
								l1425:
									{
										position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1428
										}
										position++
										goto l1427
									l1428:
										position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
										if buffer[position] != rune('S') {
											goto l1419
										}
										position++
									}
								l1427:
									{
										position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1430
										}
										position++
										goto l1429
									l1430:
										position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
										if buffer[position] != rune('T') {
											goto l1419
										}
										position++
									}
								l1429:
									{
										position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1432
										}
										position++
										goto l1431
									l1432:
										position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
										if buffer[position] != rune('R') {
											goto l1419
										}
										position++
									}
								l1431:
									if !rules[ruleskip]() {
										goto l1419
									}
									depth--
									add(ruleSUBSTR, position1420)
								}
								goto l1418
							l1419:
								position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
								{
									position1434 := position
									depth++
									{
										position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1436
										}
										position++
										goto l1435
									l1436:
										position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
										if buffer[position] != rune('R') {
											goto l1433
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
											goto l1433
										}
										position++
									}
								l1437:
									{
										position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1440
										}
										position++
										goto l1439
									l1440:
										position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
										if buffer[position] != rune('P') {
											goto l1433
										}
										position++
									}
								l1439:
									{
										position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1442
										}
										position++
										goto l1441
									l1442:
										position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
										if buffer[position] != rune('L') {
											goto l1433
										}
										position++
									}
								l1441:
									{
										position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1444
										}
										position++
										goto l1443
									l1444:
										position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
										if buffer[position] != rune('A') {
											goto l1433
										}
										position++
									}
								l1443:
									{
										position1445, tokenIndex1445, depth1445 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1446
										}
										position++
										goto l1445
									l1446:
										position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
										if buffer[position] != rune('C') {
											goto l1433
										}
										position++
									}
								l1445:
									{
										position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1448
										}
										position++
										goto l1447
									l1448:
										position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
										if buffer[position] != rune('E') {
											goto l1433
										}
										position++
									}
								l1447:
									if !rules[ruleskip]() {
										goto l1433
									}
									depth--
									add(ruleREPLACE, position1434)
								}
								goto l1418
							l1433:
								position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
								{
									position1449 := position
									depth++
									{
										position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1451
										}
										position++
										goto l1450
									l1451:
										position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
										if buffer[position] != rune('R') {
											goto l751
										}
										position++
									}
								l1450:
									{
										position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1453
										}
										position++
										goto l1452
									l1453:
										position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
										if buffer[position] != rune('E') {
											goto l751
										}
										position++
									}
								l1452:
									{
										position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1455
										}
										position++
										goto l1454
									l1455:
										position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
										if buffer[position] != rune('G') {
											goto l751
										}
										position++
									}
								l1454:
									{
										position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1457
										}
										position++
										goto l1456
									l1457:
										position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
										if buffer[position] != rune('E') {
											goto l751
										}
										position++
									}
								l1456:
									{
										position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1459
										}
										position++
										goto l1458
									l1459:
										position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
										if buffer[position] != rune('X') {
											goto l751
										}
										position++
									}
								l1458:
									if !rules[ruleskip]() {
										goto l751
									}
									depth--
									add(ruleREGEX, position1449)
								}
							}
						l1418:
							if !rules[ruleLPAREN]() {
								goto l751
							}
							if !rules[ruleexpression]() {
								goto l751
							}
							if !rules[ruleCOMMA]() {
								goto l751
							}
							if !rules[ruleexpression]() {
								goto l751
							}
							{
								position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1460
								}
								if !rules[ruleexpression]() {
									goto l1460
								}
								goto l1461
							l1460:
								position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
							}
						l1461:
							if !rules[ruleRPAREN]() {
								goto l751
							}
							break
						}
					}

				}
			l753:
				depth--
				add(rulebuiltinCall, position752)
			}
			return true
		l751:
			position, tokenIndex, depth = position751, tokenIndex751, depth751
			return false
		},
		/* 66 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1462, tokenIndex1462, depth1462 := position, tokenIndex, depth
			{
				position1463 := position
				depth++
				{
					position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1465
					}
					position++
					goto l1464
				l1465:
					position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
					if buffer[position] != rune('$') {
						goto l1462
					}
					position++
				}
			l1464:
				{
					position1466 := position
					depth++
					{
						position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
						{
							position1471 := position
							depth++
							{
								position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
								{
									position1474 := position
									depth++
									{
										position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1476
										}
										position++
										goto l1475
									l1476:
										position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1473
										}
										position++
									}
								l1475:
									depth--
									add(rulePN_CHARS_BASE, position1474)
								}
								goto l1472
							l1473:
								position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
								if buffer[position] != rune('_') {
									goto l1470
								}
								position++
							}
						l1472:
							depth--
							add(rulePN_CHARS_U, position1471)
						}
						goto l1469
					l1470:
						position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1462
						}
						position++
					}
				l1469:
				l1467:
					{
						position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
						{
							position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
							{
								position1479 := position
								depth++
								{
									position1480, tokenIndex1480, depth1480 := position, tokenIndex, depth
									{
										position1482 := position
										depth++
										{
											position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1484
											}
											position++
											goto l1483
										l1484:
											position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1481
											}
											position++
										}
									l1483:
										depth--
										add(rulePN_CHARS_BASE, position1482)
									}
									goto l1480
								l1481:
									position, tokenIndex, depth = position1480, tokenIndex1480, depth1480
									if buffer[position] != rune('_') {
										goto l1478
									}
									position++
								}
							l1480:
								depth--
								add(rulePN_CHARS_U, position1479)
							}
							goto l1477
						l1478:
							position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1468
							}
							position++
						}
					l1477:
						goto l1467
					l1468:
						position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
					}
					depth--
					add(ruleVARNAME, position1466)
				}
				if !rules[ruleskip]() {
					goto l1462
				}
				depth--
				add(rulevar, position1463)
			}
			return true
		l1462:
			position, tokenIndex, depth = position1462, tokenIndex1462, depth1462
			return false
		},
		/* 67 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1485, tokenIndex1485, depth1485 := position, tokenIndex, depth
			{
				position1486 := position
				depth++
				{
					position1487, tokenIndex1487, depth1487 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1488
					}
					goto l1487
				l1488:
					position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
					{
						position1489 := position
						depth++
					l1490:
						{
							position1491, tokenIndex1491, depth1491 := position, tokenIndex, depth
							{
								position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
								{
									position1493, tokenIndex1493, depth1493 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1494
									}
									position++
									goto l1493
								l1494:
									position, tokenIndex, depth = position1493, tokenIndex1493, depth1493
									if buffer[position] != rune(' ') {
										goto l1492
									}
									position++
								}
							l1493:
								goto l1491
							l1492:
								position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
							}
							if !matchDot() {
								goto l1491
							}
							goto l1490
						l1491:
							position, tokenIndex, depth = position1491, tokenIndex1491, depth1491
						}
						if buffer[position] != rune(':') {
							goto l1485
						}
						position++
					l1495:
						{
							position1496, tokenIndex1496, depth1496 := position, tokenIndex, depth
							{
								position1497, tokenIndex1497, depth1497 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1498
								}
								position++
								goto l1497
							l1498:
								position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1499
								}
								position++
								goto l1497
							l1499:
								position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1500
								}
								position++
								goto l1497
							l1500:
								position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1496
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1496
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1496
										}
										position++
										break
									}
								}

							}
						l1497:
							goto l1495
						l1496:
							position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
						}
						if !rules[ruleskip]() {
							goto l1485
						}
						depth--
						add(ruleprefixedName, position1489)
					}
				}
			l1487:
				depth--
				add(ruleiriref, position1486)
			}
			return true
		l1485:
			position, tokenIndex, depth = position1485, tokenIndex1485, depth1485
			return false
		},
		/* 68 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
			{
				position1503 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1502
				}
				position++
			l1504:
				{
					position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
					{
						position1506, tokenIndex1506, depth1506 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1506
						}
						position++
						goto l1505
					l1506:
						position, tokenIndex, depth = position1506, tokenIndex1506, depth1506
					}
					if !matchDot() {
						goto l1505
					}
					goto l1504
				l1505:
					position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
				}
				if buffer[position] != rune('>') {
					goto l1502
				}
				position++
				if !rules[ruleskip]() {
					goto l1502
				}
				depth--
				add(ruleiri, position1503)
			}
			return true
		l1502:
			position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
			return false
		},
		/* 69 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 70 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1508, tokenIndex1508, depth1508 := position, tokenIndex, depth
			{
				position1509 := position
				depth++
				if !rules[rulestring]() {
					goto l1508
				}
				{
					position1510, tokenIndex1510, depth1510 := position, tokenIndex, depth
					{
						position1512, tokenIndex1512, depth1512 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1513
						}
						position++
						{
							position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1517
							}
							position++
							goto l1516
						l1517:
							position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1513
							}
							position++
						}
					l1516:
					l1514:
						{
							position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
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
									goto l1515
								}
								position++
							}
						l1518:
							goto l1514
						l1515:
							position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
						}
					l1520:
						{
							position1521, tokenIndex1521, depth1521 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1521
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1521
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1521
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1521
									}
									position++
									break
								}
							}

						l1522:
							{
								position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1523
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1523
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1523
										}
										position++
										break
									}
								}

								goto l1522
							l1523:
								position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
							}
							goto l1520
						l1521:
							position, tokenIndex, depth = position1521, tokenIndex1521, depth1521
						}
						goto l1512
					l1513:
						position, tokenIndex, depth = position1512, tokenIndex1512, depth1512
						if buffer[position] != rune('^') {
							goto l1510
						}
						position++
						if buffer[position] != rune('^') {
							goto l1510
						}
						position++
						if !rules[ruleiriref]() {
							goto l1510
						}
					}
				l1512:
					goto l1511
				l1510:
					position, tokenIndex, depth = position1510, tokenIndex1510, depth1510
				}
			l1511:
				if !rules[ruleskip]() {
					goto l1508
				}
				depth--
				add(ruleliteral, position1509)
			}
			return true
		l1508:
			position, tokenIndex, depth = position1508, tokenIndex1508, depth1508
			return false
		},
		/* 71 string <- <('"' (!'"' .)* '"')> */
		func() bool {
			position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
			{
				position1527 := position
				depth++
				if buffer[position] != rune('"') {
					goto l1526
				}
				position++
			l1528:
				{
					position1529, tokenIndex1529, depth1529 := position, tokenIndex, depth
					{
						position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
						if buffer[position] != rune('"') {
							goto l1530
						}
						position++
						goto l1529
					l1530:
						position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
					}
					if !matchDot() {
						goto l1529
					}
					goto l1528
				l1529:
					position, tokenIndex, depth = position1529, tokenIndex1529, depth1529
				}
				if buffer[position] != rune('"') {
					goto l1526
				}
				position++
				depth--
				add(rulestring, position1527)
			}
			return true
		l1526:
			position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
			return false
		},
		/* 72 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
			{
				position1532 := position
				depth++
				{
					position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
					{
						position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1536
						}
						position++
						goto l1535
					l1536:
						position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
						if buffer[position] != rune('-') {
							goto l1533
						}
						position++
					}
				l1535:
					goto l1534
				l1533:
					position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
				}
			l1534:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1531
				}
				position++
			l1537:
				{
					position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1538
					}
					position++
					goto l1537
				l1538:
					position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
				}
				{
					position1539, tokenIndex1539, depth1539 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1539
					}
					position++
				l1541:
					{
						position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1542
						}
						position++
						goto l1541
					l1542:
						position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
					}
					goto l1540
				l1539:
					position, tokenIndex, depth = position1539, tokenIndex1539, depth1539
				}
			l1540:
				if !rules[ruleskip]() {
					goto l1531
				}
				depth--
				add(rulenumericLiteral, position1532)
			}
			return true
		l1531:
			position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
			return false
		},
		/* 73 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 74 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1544, tokenIndex1544, depth1544 := position, tokenIndex, depth
			{
				position1545 := position
				depth++
				{
					position1546, tokenIndex1546, depth1546 := position, tokenIndex, depth
					{
						position1548 := position
						depth++
						{
							position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1550
							}
							position++
							goto l1549
						l1550:
							position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
							if buffer[position] != rune('T') {
								goto l1547
							}
							position++
						}
					l1549:
						{
							position1551, tokenIndex1551, depth1551 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1552
							}
							position++
							goto l1551
						l1552:
							position, tokenIndex, depth = position1551, tokenIndex1551, depth1551
							if buffer[position] != rune('R') {
								goto l1547
							}
							position++
						}
					l1551:
						{
							position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1554
							}
							position++
							goto l1553
						l1554:
							position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
							if buffer[position] != rune('U') {
								goto l1547
							}
							position++
						}
					l1553:
						{
							position1555, tokenIndex1555, depth1555 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1556
							}
							position++
							goto l1555
						l1556:
							position, tokenIndex, depth = position1555, tokenIndex1555, depth1555
							if buffer[position] != rune('E') {
								goto l1547
							}
							position++
						}
					l1555:
						if !rules[ruleskip]() {
							goto l1547
						}
						depth--
						add(ruleTRUE, position1548)
					}
					goto l1546
				l1547:
					position, tokenIndex, depth = position1546, tokenIndex1546, depth1546
					{
						position1557 := position
						depth++
						{
							position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1559
							}
							position++
							goto l1558
						l1559:
							position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
							if buffer[position] != rune('F') {
								goto l1544
							}
							position++
						}
					l1558:
						{
							position1560, tokenIndex1560, depth1560 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1561
							}
							position++
							goto l1560
						l1561:
							position, tokenIndex, depth = position1560, tokenIndex1560, depth1560
							if buffer[position] != rune('A') {
								goto l1544
							}
							position++
						}
					l1560:
						{
							position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1563
							}
							position++
							goto l1562
						l1563:
							position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
							if buffer[position] != rune('L') {
								goto l1544
							}
							position++
						}
					l1562:
						{
							position1564, tokenIndex1564, depth1564 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1565
							}
							position++
							goto l1564
						l1565:
							position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
							if buffer[position] != rune('S') {
								goto l1544
							}
							position++
						}
					l1564:
						{
							position1566, tokenIndex1566, depth1566 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1567
							}
							position++
							goto l1566
						l1567:
							position, tokenIndex, depth = position1566, tokenIndex1566, depth1566
							if buffer[position] != rune('E') {
								goto l1544
							}
							position++
						}
					l1566:
						if !rules[ruleskip]() {
							goto l1544
						}
						depth--
						add(ruleFALSE, position1557)
					}
				}
			l1546:
				depth--
				add(rulebooleanLiteral, position1545)
			}
			return true
		l1544:
			position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
			return false
		},
		/* 75 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 76 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 77 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 78 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1571, tokenIndex1571, depth1571 := position, tokenIndex, depth
			{
				position1572 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1571
				}
				position++
			l1573:
				{
					position1574, tokenIndex1574, depth1574 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1574
					}
					goto l1573
				l1574:
					position, tokenIndex, depth = position1574, tokenIndex1574, depth1574
				}
				if buffer[position] != rune(')') {
					goto l1571
				}
				position++
				if !rules[ruleskip]() {
					goto l1571
				}
				depth--
				add(rulenil, position1572)
			}
			return true
		l1571:
			position, tokenIndex, depth = position1571, tokenIndex1571, depth1571
			return false
		},
		/* 79 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 80 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 81 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 82 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 83 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 84 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 85 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 86 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 87 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 88 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1584, tokenIndex1584, depth1584 := position, tokenIndex, depth
			{
				position1585 := position
				depth++
				{
					position1586, tokenIndex1586, depth1586 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1587
					}
					position++
					goto l1586
				l1587:
					position, tokenIndex, depth = position1586, tokenIndex1586, depth1586
					if buffer[position] != rune('D') {
						goto l1584
					}
					position++
				}
			l1586:
				{
					position1588, tokenIndex1588, depth1588 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1589
					}
					position++
					goto l1588
				l1589:
					position, tokenIndex, depth = position1588, tokenIndex1588, depth1588
					if buffer[position] != rune('I') {
						goto l1584
					}
					position++
				}
			l1588:
				{
					position1590, tokenIndex1590, depth1590 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1591
					}
					position++
					goto l1590
				l1591:
					position, tokenIndex, depth = position1590, tokenIndex1590, depth1590
					if buffer[position] != rune('S') {
						goto l1584
					}
					position++
				}
			l1590:
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
						goto l1584
					}
					position++
				}
			l1592:
				{
					position1594, tokenIndex1594, depth1594 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1595
					}
					position++
					goto l1594
				l1595:
					position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
					if buffer[position] != rune('I') {
						goto l1584
					}
					position++
				}
			l1594:
				{
					position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1597
					}
					position++
					goto l1596
				l1597:
					position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
					if buffer[position] != rune('N') {
						goto l1584
					}
					position++
				}
			l1596:
				{
					position1598, tokenIndex1598, depth1598 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1599
					}
					position++
					goto l1598
				l1599:
					position, tokenIndex, depth = position1598, tokenIndex1598, depth1598
					if buffer[position] != rune('C') {
						goto l1584
					}
					position++
				}
			l1598:
				{
					position1600, tokenIndex1600, depth1600 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1601
					}
					position++
					goto l1600
				l1601:
					position, tokenIndex, depth = position1600, tokenIndex1600, depth1600
					if buffer[position] != rune('T') {
						goto l1584
					}
					position++
				}
			l1600:
				if !rules[ruleskip]() {
					goto l1584
				}
				depth--
				add(ruleDISTINCT, position1585)
			}
			return true
		l1584:
			position, tokenIndex, depth = position1584, tokenIndex1584, depth1584
			return false
		},
		/* 89 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 90 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 91 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 92 LBRACE <- <('{' skip)> */
		func() bool {
			position1605, tokenIndex1605, depth1605 := position, tokenIndex, depth
			{
				position1606 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1605
				}
				position++
				if !rules[ruleskip]() {
					goto l1605
				}
				depth--
				add(ruleLBRACE, position1606)
			}
			return true
		l1605:
			position, tokenIndex, depth = position1605, tokenIndex1605, depth1605
			return false
		},
		/* 93 RBRACE <- <('}' skip)> */
		func() bool {
			position1607, tokenIndex1607, depth1607 := position, tokenIndex, depth
			{
				position1608 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1607
				}
				position++
				if !rules[ruleskip]() {
					goto l1607
				}
				depth--
				add(ruleRBRACE, position1608)
			}
			return true
		l1607:
			position, tokenIndex, depth = position1607, tokenIndex1607, depth1607
			return false
		},
		/* 94 LBRACK <- <('[' skip)> */
		nil,
		/* 95 RBRACK <- <(']' skip)> */
		nil,
		/* 96 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
			{
				position1612 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1611
				}
				position++
				if !rules[ruleskip]() {
					goto l1611
				}
				depth--
				add(ruleSEMICOLON, position1612)
			}
			return true
		l1611:
			position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
			return false
		},
		/* 97 COMMA <- <(',' skip)> */
		func() bool {
			position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
			{
				position1614 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1613
				}
				position++
				if !rules[ruleskip]() {
					goto l1613
				}
				depth--
				add(ruleCOMMA, position1614)
			}
			return true
		l1613:
			position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
			return false
		},
		/* 98 DOT <- <('.' skip)> */
		func() bool {
			position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
			{
				position1616 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1615
				}
				position++
				if !rules[ruleskip]() {
					goto l1615
				}
				depth--
				add(ruleDOT, position1616)
			}
			return true
		l1615:
			position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
			return false
		},
		/* 99 COLON <- <(':' skip)> */
		nil,
		/* 100 PIPE <- <('|' skip)> */
		func() bool {
			position1618, tokenIndex1618, depth1618 := position, tokenIndex, depth
			{
				position1619 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1618
				}
				position++
				if !rules[ruleskip]() {
					goto l1618
				}
				depth--
				add(rulePIPE, position1619)
			}
			return true
		l1618:
			position, tokenIndex, depth = position1618, tokenIndex1618, depth1618
			return false
		},
		/* 101 SLASH <- <('/' skip)> */
		func() bool {
			position1620, tokenIndex1620, depth1620 := position, tokenIndex, depth
			{
				position1621 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1620
				}
				position++
				if !rules[ruleskip]() {
					goto l1620
				}
				depth--
				add(ruleSLASH, position1621)
			}
			return true
		l1620:
			position, tokenIndex, depth = position1620, tokenIndex1620, depth1620
			return false
		},
		/* 102 INVERSE <- <('^' skip)> */
		func() bool {
			position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
			{
				position1623 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1622
				}
				position++
				if !rules[ruleskip]() {
					goto l1622
				}
				depth--
				add(ruleINVERSE, position1623)
			}
			return true
		l1622:
			position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
			return false
		},
		/* 103 LPAREN <- <('(' skip)> */
		func() bool {
			position1624, tokenIndex1624, depth1624 := position, tokenIndex, depth
			{
				position1625 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1624
				}
				position++
				if !rules[ruleskip]() {
					goto l1624
				}
				depth--
				add(ruleLPAREN, position1625)
			}
			return true
		l1624:
			position, tokenIndex, depth = position1624, tokenIndex1624, depth1624
			return false
		},
		/* 104 RPAREN <- <(')' skip)> */
		func() bool {
			position1626, tokenIndex1626, depth1626 := position, tokenIndex, depth
			{
				position1627 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1626
				}
				position++
				if !rules[ruleskip]() {
					goto l1626
				}
				depth--
				add(ruleRPAREN, position1627)
			}
			return true
		l1626:
			position, tokenIndex, depth = position1626, tokenIndex1626, depth1626
			return false
		},
		/* 105 ISA <- <('a' skip)> */
		func() bool {
			position1628, tokenIndex1628, depth1628 := position, tokenIndex, depth
			{
				position1629 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1628
				}
				position++
				if !rules[ruleskip]() {
					goto l1628
				}
				depth--
				add(ruleISA, position1629)
			}
			return true
		l1628:
			position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
			return false
		},
		/* 106 NOT <- <('!' skip)> */
		func() bool {
			position1630, tokenIndex1630, depth1630 := position, tokenIndex, depth
			{
				position1631 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1630
				}
				position++
				if !rules[ruleskip]() {
					goto l1630
				}
				depth--
				add(ruleNOT, position1631)
			}
			return true
		l1630:
			position, tokenIndex, depth = position1630, tokenIndex1630, depth1630
			return false
		},
		/* 107 STAR <- <('*' skip)> */
		func() bool {
			position1632, tokenIndex1632, depth1632 := position, tokenIndex, depth
			{
				position1633 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1632
				}
				position++
				if !rules[ruleskip]() {
					goto l1632
				}
				depth--
				add(ruleSTAR, position1633)
			}
			return true
		l1632:
			position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
			return false
		},
		/* 108 PLUS <- <('+' skip)> */
		func() bool {
			position1634, tokenIndex1634, depth1634 := position, tokenIndex, depth
			{
				position1635 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1634
				}
				position++
				if !rules[ruleskip]() {
					goto l1634
				}
				depth--
				add(rulePLUS, position1635)
			}
			return true
		l1634:
			position, tokenIndex, depth = position1634, tokenIndex1634, depth1634
			return false
		},
		/* 109 MINUS <- <('-' skip)> */
		func() bool {
			position1636, tokenIndex1636, depth1636 := position, tokenIndex, depth
			{
				position1637 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1636
				}
				position++
				if !rules[ruleskip]() {
					goto l1636
				}
				depth--
				add(ruleMINUS, position1637)
			}
			return true
		l1636:
			position, tokenIndex, depth = position1636, tokenIndex1636, depth1636
			return false
		},
		/* 110 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 111 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 112 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 113 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 114 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1642, tokenIndex1642, depth1642 := position, tokenIndex, depth
			{
				position1643 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1642
				}
				position++
			l1644:
				{
					position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1645
					}
					position++
					goto l1644
				l1645:
					position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
				}
				if !rules[ruleskip]() {
					goto l1642
				}
				depth--
				add(ruleINTEGER, position1643)
			}
			return true
		l1642:
			position, tokenIndex, depth = position1642, tokenIndex1642, depth1642
			return false
		},
		/* 115 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 116 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 117 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 118 OR <- <('|' '|' skip)> */
		nil,
		/* 119 AND <- <('&' '&' skip)> */
		nil,
		/* 120 EQ <- <('=' skip)> */
		func() bool {
			position1651, tokenIndex1651, depth1651 := position, tokenIndex, depth
			{
				position1652 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1651
				}
				position++
				if !rules[ruleskip]() {
					goto l1651
				}
				depth--
				add(ruleEQ, position1652)
			}
			return true
		l1651:
			position, tokenIndex, depth = position1651, tokenIndex1651, depth1651
			return false
		},
		/* 121 NE <- <('!' '=' skip)> */
		nil,
		/* 122 GT <- <('>' skip)> */
		nil,
		/* 123 LT <- <('<' skip)> */
		nil,
		/* 124 LE <- <('<' '=' skip)> */
		nil,
		/* 125 GE <- <('>' '=' skip)> */
		nil,
		/* 126 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 127 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 128 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
			{
				position1661 := position
				depth++
				{
					position1662, tokenIndex1662, depth1662 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1663
					}
					position++
					goto l1662
				l1663:
					position, tokenIndex, depth = position1662, tokenIndex1662, depth1662
					if buffer[position] != rune('A') {
						goto l1660
					}
					position++
				}
			l1662:
				{
					position1664, tokenIndex1664, depth1664 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1665
					}
					position++
					goto l1664
				l1665:
					position, tokenIndex, depth = position1664, tokenIndex1664, depth1664
					if buffer[position] != rune('S') {
						goto l1660
					}
					position++
				}
			l1664:
				if !rules[ruleskip]() {
					goto l1660
				}
				depth--
				add(ruleAS, position1661)
			}
			return true
		l1660:
			position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
			return false
		},
		/* 129 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 130 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 131 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 132 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 133 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 134 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 135 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 136 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 137 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 138 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 139 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 140 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 141 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 142 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 143 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 144 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 145 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 146 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 147 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 148 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 149 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 150 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 151 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 152 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 153 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 154 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 155 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 156 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 157 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 158 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 159 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 160 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 161 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 162 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 163 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 164 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 165 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 166 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 167 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 168 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 169 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 170 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 171 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 172 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 173 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 174 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 175 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 176 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 177 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 178 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 179 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 180 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 181 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 182 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 183 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 184 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 185 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 186 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 187 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 188 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 189 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 190 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 191 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 192 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 193 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 194 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 195 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 196 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 197 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1734, tokenIndex1734, depth1734 := position, tokenIndex, depth
			{
				position1735 := position
				depth++
				{
					position1736, tokenIndex1736, depth1736 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1737
					}
					position++
					goto l1736
				l1737:
					position, tokenIndex, depth = position1736, tokenIndex1736, depth1736
					if buffer[position] != rune('B') {
						goto l1734
					}
					position++
				}
			l1736:
				{
					position1738, tokenIndex1738, depth1738 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1739
					}
					position++
					goto l1738
				l1739:
					position, tokenIndex, depth = position1738, tokenIndex1738, depth1738
					if buffer[position] != rune('Y') {
						goto l1734
					}
					position++
				}
			l1738:
				if !rules[ruleskip]() {
					goto l1734
				}
				depth--
				add(ruleBY, position1735)
			}
			return true
		l1734:
			position, tokenIndex, depth = position1734, tokenIndex1734, depth1734
			return false
		},
		/* 198 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 199 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1742 := position
				depth++
			l1743:
				{
					position1744, tokenIndex1744, depth1744 := position, tokenIndex, depth
					{
						position1745, tokenIndex1745, depth1745 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1746
						}
						goto l1745
					l1746:
						position, tokenIndex, depth = position1745, tokenIndex1745, depth1745
						{
							position1747 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1744
							}
							position++
						l1748:
							{
								position1749, tokenIndex1749, depth1749 := position, tokenIndex, depth
								{
									position1750, tokenIndex1750, depth1750 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1750
									}
									goto l1749
								l1750:
									position, tokenIndex, depth = position1750, tokenIndex1750, depth1750
								}
								if !matchDot() {
									goto l1749
								}
								goto l1748
							l1749:
								position, tokenIndex, depth = position1749, tokenIndex1749, depth1749
							}
							if !rules[ruleendOfLine]() {
								goto l1744
							}
							depth--
							add(rulecomment, position1747)
						}
					}
				l1745:
					goto l1743
				l1744:
					position, tokenIndex, depth = position1744, tokenIndex1744, depth1744
				}
				depth--
				add(ruleskip, position1742)
			}
			return true
		},
		/* 200 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1751, tokenIndex1751, depth1751 := position, tokenIndex, depth
			{
				position1752 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1751
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1751
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1751
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1751
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1751
						}
						break
					}
				}

				depth--
				add(rulews, position1752)
			}
			return true
		l1751:
			position, tokenIndex, depth = position1751, tokenIndex1751, depth1751
			return false
		},
		/* 201 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 202 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1755, tokenIndex1755, depth1755 := position, tokenIndex, depth
			{
				position1756 := position
				depth++
				{
					position1757, tokenIndex1757, depth1757 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1758
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1758
					}
					position++
					goto l1757
				l1758:
					position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
					if buffer[position] != rune('\n') {
						goto l1759
					}
					position++
					goto l1757
				l1759:
					position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
					if buffer[position] != rune('\r') {
						goto l1755
					}
					position++
				}
			l1757:
				depth--
				add(ruleendOfLine, position1756)
			}
			return true
		l1755:
			position, tokenIndex, depth = position1755, tokenIndex1755, depth1755
			return false
		},
	}
	p.rules = rules
}
