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
	rules  [196]func() bool
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
		/* 7 subSelect <- <(select whereClause)> */
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
					{
						position289 := position
						depth++
						{
							position290, tokenIndex290, depth290 := position, tokenIndex, depth
							if !rules[rulebrackettedExpression]() {
								goto l291
							}
							goto l290
						l291:
							position, tokenIndex, depth = position290, tokenIndex290, depth290
							if !rules[rulebuiltinCall]() {
								goto l292
							}
							goto l290
						l292:
							position, tokenIndex, depth = position290, tokenIndex290, depth290
							if !rules[rulefunctionCall]() {
								goto l275
							}
						}
					l290:
						depth--
						add(ruleconstraint, position289)
					}
					goto l274
				l275:
					position, tokenIndex, depth = position274, tokenIndex274, depth274
					{
						position293 := position
						depth++
						{
							position294, tokenIndex294, depth294 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l295
							}
							position++
							goto l294
						l295:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							if buffer[position] != rune('B') {
								goto l272
							}
							position++
						}
					l294:
						{
							position296, tokenIndex296, depth296 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l297
							}
							position++
							goto l296
						l297:
							position, tokenIndex, depth = position296, tokenIndex296, depth296
							if buffer[position] != rune('I') {
								goto l272
							}
							position++
						}
					l296:
						{
							position298, tokenIndex298, depth298 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l299
							}
							position++
							goto l298
						l299:
							position, tokenIndex, depth = position298, tokenIndex298, depth298
							if buffer[position] != rune('N') {
								goto l272
							}
							position++
						}
					l298:
						{
							position300, tokenIndex300, depth300 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l301
							}
							position++
							goto l300
						l301:
							position, tokenIndex, depth = position300, tokenIndex300, depth300
							if buffer[position] != rune('D') {
								goto l272
							}
							position++
						}
					l300:
						if !rules[ruleskip]() {
							goto l272
						}
						depth--
						add(ruleBIND, position293)
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
		nil,
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
		/* 43 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position399 := position
				depth++
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					{
						position402 := position
						depth++
						{
							position403, tokenIndex403, depth403 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l404
							}
							{
								position405, tokenIndex405, depth405 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l405
								}
								goto l406
							l405:
								position, tokenIndex, depth = position405, tokenIndex405, depth405
							}
						l406:
							goto l403
						l404:
							position, tokenIndex, depth = position403, tokenIndex403, depth403
							if !rules[ruleoffset]() {
								goto l400
							}
							{
								position407, tokenIndex407, depth407 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l407
								}
								goto l408
							l407:
								position, tokenIndex, depth = position407, tokenIndex407, depth407
							}
						l408:
						}
					l403:
						depth--
						add(rulelimitOffsetClauses, position402)
					}
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
		/* 44 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 45 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position410, tokenIndex410, depth410 := position, tokenIndex, depth
			{
				position411 := position
				depth++
				{
					position412 := position
					depth++
					{
						position413, tokenIndex413, depth413 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l414
						}
						position++
						goto l413
					l414:
						position, tokenIndex, depth = position413, tokenIndex413, depth413
						if buffer[position] != rune('L') {
							goto l410
						}
						position++
					}
				l413:
					{
						position415, tokenIndex415, depth415 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l416
						}
						position++
						goto l415
					l416:
						position, tokenIndex, depth = position415, tokenIndex415, depth415
						if buffer[position] != rune('I') {
							goto l410
						}
						position++
					}
				l415:
					{
						position417, tokenIndex417, depth417 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l418
						}
						position++
						goto l417
					l418:
						position, tokenIndex, depth = position417, tokenIndex417, depth417
						if buffer[position] != rune('M') {
							goto l410
						}
						position++
					}
				l417:
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l420
						}
						position++
						goto l419
					l420:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
						if buffer[position] != rune('I') {
							goto l410
						}
						position++
					}
				l419:
					{
						position421, tokenIndex421, depth421 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l422
						}
						position++
						goto l421
					l422:
						position, tokenIndex, depth = position421, tokenIndex421, depth421
						if buffer[position] != rune('T') {
							goto l410
						}
						position++
					}
				l421:
					if !rules[ruleskip]() {
						goto l410
					}
					depth--
					add(ruleLIMIT, position412)
				}
				if !rules[ruleINTEGER]() {
					goto l410
				}
				depth--
				add(rulelimit, position411)
			}
			return true
		l410:
			position, tokenIndex, depth = position410, tokenIndex410, depth410
			return false
		},
		/* 46 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position423, tokenIndex423, depth423 := position, tokenIndex, depth
			{
				position424 := position
				depth++
				{
					position425 := position
					depth++
					{
						position426, tokenIndex426, depth426 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l427
						}
						position++
						goto l426
					l427:
						position, tokenIndex, depth = position426, tokenIndex426, depth426
						if buffer[position] != rune('O') {
							goto l423
						}
						position++
					}
				l426:
					{
						position428, tokenIndex428, depth428 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l429
						}
						position++
						goto l428
					l429:
						position, tokenIndex, depth = position428, tokenIndex428, depth428
						if buffer[position] != rune('F') {
							goto l423
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
							goto l423
						}
						position++
					}
				l430:
					{
						position432, tokenIndex432, depth432 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l433
						}
						position++
						goto l432
					l433:
						position, tokenIndex, depth = position432, tokenIndex432, depth432
						if buffer[position] != rune('S') {
							goto l423
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
							goto l423
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
							goto l423
						}
						position++
					}
				l436:
					if !rules[ruleskip]() {
						goto l423
					}
					depth--
					add(ruleOFFSET, position425)
				}
				if !rules[ruleINTEGER]() {
					goto l423
				}
				depth--
				add(ruleoffset, position424)
			}
			return true
		l423:
			position, tokenIndex, depth = position423, tokenIndex423, depth423
			return false
		},
		/* 47 expression <- <conditionalOrExpression> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l438
				}
				depth--
				add(ruleexpression, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 48 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l440
				}
				{
					position442, tokenIndex442, depth442 := position, tokenIndex, depth
					{
						position444 := position
						depth++
						if buffer[position] != rune('|') {
							goto l442
						}
						position++
						if buffer[position] != rune('|') {
							goto l442
						}
						position++
						if !rules[ruleskip]() {
							goto l442
						}
						depth--
						add(ruleOR, position444)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l442
					}
					goto l443
				l442:
					position, tokenIndex, depth = position442, tokenIndex442, depth442
				}
			l443:
				depth--
				add(ruleconditionalOrExpression, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 49 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				{
					position447 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l445
					}
					{
						position448, tokenIndex448, depth448 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position451 := position
									depth++
									{
										position452 := position
										depth++
										{
											position453, tokenIndex453, depth453 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l454
											}
											position++
											goto l453
										l454:
											position, tokenIndex, depth = position453, tokenIndex453, depth453
											if buffer[position] != rune('N') {
												goto l448
											}
											position++
										}
									l453:
										{
											position455, tokenIndex455, depth455 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l456
											}
											position++
											goto l455
										l456:
											position, tokenIndex, depth = position455, tokenIndex455, depth455
											if buffer[position] != rune('O') {
												goto l448
											}
											position++
										}
									l455:
										{
											position457, tokenIndex457, depth457 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l458
											}
											position++
											goto l457
										l458:
											position, tokenIndex, depth = position457, tokenIndex457, depth457
											if buffer[position] != rune('T') {
												goto l448
											}
											position++
										}
									l457:
										if buffer[position] != rune(' ') {
											goto l448
										}
										position++
										{
											position459, tokenIndex459, depth459 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l460
											}
											position++
											goto l459
										l460:
											position, tokenIndex, depth = position459, tokenIndex459, depth459
											if buffer[position] != rune('I') {
												goto l448
											}
											position++
										}
									l459:
										{
											position461, tokenIndex461, depth461 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l462
											}
											position++
											goto l461
										l462:
											position, tokenIndex, depth = position461, tokenIndex461, depth461
											if buffer[position] != rune('N') {
												goto l448
											}
											position++
										}
									l461:
										if !rules[ruleskip]() {
											goto l448
										}
										depth--
										add(ruleNOTIN, position452)
									}
									if !rules[ruleargList]() {
										goto l448
									}
									depth--
									add(rulenotin, position451)
								}
								break
							case 'I', 'i':
								{
									position463 := position
									depth++
									{
										position464 := position
										depth++
										{
											position465, tokenIndex465, depth465 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l466
											}
											position++
											goto l465
										l466:
											position, tokenIndex, depth = position465, tokenIndex465, depth465
											if buffer[position] != rune('I') {
												goto l448
											}
											position++
										}
									l465:
										{
											position467, tokenIndex467, depth467 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l468
											}
											position++
											goto l467
										l468:
											position, tokenIndex, depth = position467, tokenIndex467, depth467
											if buffer[position] != rune('N') {
												goto l448
											}
											position++
										}
									l467:
										if !rules[ruleskip]() {
											goto l448
										}
										depth--
										add(ruleIN, position464)
									}
									if !rules[ruleargList]() {
										goto l448
									}
									depth--
									add(rulein, position463)
								}
								break
							default:
								{
									position469, tokenIndex469, depth469 := position, tokenIndex, depth
									{
										position471 := position
										depth++
										if buffer[position] != rune('<') {
											goto l470
										}
										position++
										if !rules[ruleskip]() {
											goto l470
										}
										depth--
										add(ruleLT, position471)
									}
									goto l469
								l470:
									position, tokenIndex, depth = position469, tokenIndex469, depth469
									{
										position473 := position
										depth++
										if buffer[position] != rune('>') {
											goto l472
										}
										position++
										if buffer[position] != rune('=') {
											goto l472
										}
										position++
										if !rules[ruleskip]() {
											goto l472
										}
										depth--
										add(ruleGE, position473)
									}
									goto l469
								l472:
									position, tokenIndex, depth = position469, tokenIndex469, depth469
									{
										switch buffer[position] {
										case '>':
											{
												position475 := position
												depth++
												if buffer[position] != rune('>') {
													goto l448
												}
												position++
												if !rules[ruleskip]() {
													goto l448
												}
												depth--
												add(ruleGT, position475)
											}
											break
										case '<':
											{
												position476 := position
												depth++
												if buffer[position] != rune('<') {
													goto l448
												}
												position++
												if buffer[position] != rune('=') {
													goto l448
												}
												position++
												if !rules[ruleskip]() {
													goto l448
												}
												depth--
												add(ruleLE, position476)
											}
											break
										case '!':
											{
												position477 := position
												depth++
												if buffer[position] != rune('!') {
													goto l448
												}
												position++
												if buffer[position] != rune('=') {
													goto l448
												}
												position++
												if !rules[ruleskip]() {
													goto l448
												}
												depth--
												add(ruleNE, position477)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l448
											}
											break
										}
									}

								}
							l469:
								if !rules[rulenumericExpression]() {
									goto l448
								}
								break
							}
						}

						goto l449
					l448:
						position, tokenIndex, depth = position448, tokenIndex448, depth448
					}
				l449:
					depth--
					add(rulevalueLogical, position447)
				}
				{
					position478, tokenIndex478, depth478 := position, tokenIndex, depth
					{
						position480 := position
						depth++
						if buffer[position] != rune('&') {
							goto l478
						}
						position++
						if buffer[position] != rune('&') {
							goto l478
						}
						position++
						if !rules[ruleskip]() {
							goto l478
						}
						depth--
						add(ruleAND, position480)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l478
					}
					goto l479
				l478:
					position, tokenIndex, depth = position478, tokenIndex478, depth478
				}
			l479:
				depth--
				add(ruleconditionalAndExpression, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 50 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 51 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position482, tokenIndex482, depth482 := position, tokenIndex, depth
			{
				position483 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l482
				}
			l484:
				{
					position485, tokenIndex485, depth485 := position, tokenIndex, depth
					{
						position486, tokenIndex486, depth486 := position, tokenIndex, depth
						{
							position488, tokenIndex488, depth488 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l489
							}
							goto l488
						l489:
							position, tokenIndex, depth = position488, tokenIndex488, depth488
							if !rules[ruleMINUS]() {
								goto l487
							}
						}
					l488:
						if !rules[rulemultiplicativeExpression]() {
							goto l487
						}
						goto l486
					l487:
						position, tokenIndex, depth = position486, tokenIndex486, depth486
						{
							position490 := position
							depth++
							{
								position491, tokenIndex491, depth491 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l492
								}
								position++
								goto l491
							l492:
								position, tokenIndex, depth = position491, tokenIndex491, depth491
								if buffer[position] != rune('-') {
									goto l485
								}
								position++
							}
						l491:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l485
							}
							position++
						l493:
							{
								position494, tokenIndex494, depth494 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l494
								}
								position++
								goto l493
							l494:
								position, tokenIndex, depth = position494, tokenIndex494, depth494
							}
							{
								position495, tokenIndex495, depth495 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l495
								}
								position++
							l497:
								{
									position498, tokenIndex498, depth498 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l498
									}
									position++
									goto l497
								l498:
									position, tokenIndex, depth = position498, tokenIndex498, depth498
								}
								goto l496
							l495:
								position, tokenIndex, depth = position495, tokenIndex495, depth495
							}
						l496:
							if !rules[ruleskip]() {
								goto l485
							}
							depth--
							add(rulesignedNumericLiteral, position490)
						}
					}
				l486:
					goto l484
				l485:
					position, tokenIndex, depth = position485, tokenIndex485, depth485
				}
				depth--
				add(rulenumericExpression, position483)
			}
			return true
		l482:
			position, tokenIndex, depth = position482, tokenIndex482, depth482
			return false
		},
		/* 52 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position499, tokenIndex499, depth499 := position, tokenIndex, depth
			{
				position500 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l499
				}
			l501:
				{
					position502, tokenIndex502, depth502 := position, tokenIndex, depth
					{
						position503, tokenIndex503, depth503 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l504
						}
						goto l503
					l504:
						position, tokenIndex, depth = position503, tokenIndex503, depth503
						if !rules[ruleSLASH]() {
							goto l502
						}
					}
				l503:
					if !rules[ruleunaryExpression]() {
						goto l502
					}
					goto l501
				l502:
					position, tokenIndex, depth = position502, tokenIndex502, depth502
				}
				depth--
				add(rulemultiplicativeExpression, position500)
			}
			return true
		l499:
			position, tokenIndex, depth = position499, tokenIndex499, depth499
			return false
		},
		/* 53 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position505, tokenIndex505, depth505 := position, tokenIndex, depth
			{
				position506 := position
				depth++
				{
					position507, tokenIndex507, depth507 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l507
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l507
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l507
							}
							break
						}
					}

					goto l508
				l507:
					position, tokenIndex, depth = position507, tokenIndex507, depth507
				}
			l508:
				{
					position510 := position
					depth++
					{
						position511, tokenIndex511, depth511 := position, tokenIndex, depth
						if !rules[rulebrackettedExpression]() {
							goto l512
						}
						goto l511
					l512:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
						if !rules[rulebuiltinCall]() {
							goto l513
						}
						goto l511
					l513:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
						if !rules[rulefunctionCall]() {
							goto l514
						}
						goto l511
					l514:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
						if !rules[ruleiriref]() {
							goto l515
						}
						goto l511
					l515:
						position, tokenIndex, depth = position511, tokenIndex511, depth511
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position517 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position519 := position
												depth++
												{
													position520 := position
													depth++
													{
														position521, tokenIndex521, depth521 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l522
														}
														position++
														goto l521
													l522:
														position, tokenIndex, depth = position521, tokenIndex521, depth521
														if buffer[position] != rune('G') {
															goto l505
														}
														position++
													}
												l521:
													{
														position523, tokenIndex523, depth523 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l524
														}
														position++
														goto l523
													l524:
														position, tokenIndex, depth = position523, tokenIndex523, depth523
														if buffer[position] != rune('R') {
															goto l505
														}
														position++
													}
												l523:
													{
														position525, tokenIndex525, depth525 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l526
														}
														position++
														goto l525
													l526:
														position, tokenIndex, depth = position525, tokenIndex525, depth525
														if buffer[position] != rune('O') {
															goto l505
														}
														position++
													}
												l525:
													{
														position527, tokenIndex527, depth527 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l528
														}
														position++
														goto l527
													l528:
														position, tokenIndex, depth = position527, tokenIndex527, depth527
														if buffer[position] != rune('U') {
															goto l505
														}
														position++
													}
												l527:
													{
														position529, tokenIndex529, depth529 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l530
														}
														position++
														goto l529
													l530:
														position, tokenIndex, depth = position529, tokenIndex529, depth529
														if buffer[position] != rune('P') {
															goto l505
														}
														position++
													}
												l529:
													if buffer[position] != rune('_') {
														goto l505
													}
													position++
													{
														position531, tokenIndex531, depth531 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l532
														}
														position++
														goto l531
													l532:
														position, tokenIndex, depth = position531, tokenIndex531, depth531
														if buffer[position] != rune('C') {
															goto l505
														}
														position++
													}
												l531:
													{
														position533, tokenIndex533, depth533 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l534
														}
														position++
														goto l533
													l534:
														position, tokenIndex, depth = position533, tokenIndex533, depth533
														if buffer[position] != rune('O') {
															goto l505
														}
														position++
													}
												l533:
													{
														position535, tokenIndex535, depth535 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l536
														}
														position++
														goto l535
													l536:
														position, tokenIndex, depth = position535, tokenIndex535, depth535
														if buffer[position] != rune('N') {
															goto l505
														}
														position++
													}
												l535:
													{
														position537, tokenIndex537, depth537 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l538
														}
														position++
														goto l537
													l538:
														position, tokenIndex, depth = position537, tokenIndex537, depth537
														if buffer[position] != rune('C') {
															goto l505
														}
														position++
													}
												l537:
													{
														position539, tokenIndex539, depth539 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l540
														}
														position++
														goto l539
													l540:
														position, tokenIndex, depth = position539, tokenIndex539, depth539
														if buffer[position] != rune('A') {
															goto l505
														}
														position++
													}
												l539:
													{
														position541, tokenIndex541, depth541 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l542
														}
														position++
														goto l541
													l542:
														position, tokenIndex, depth = position541, tokenIndex541, depth541
														if buffer[position] != rune('T') {
															goto l505
														}
														position++
													}
												l541:
													if !rules[ruleskip]() {
														goto l505
													}
													depth--
													add(ruleGROUPCONCAT, position520)
												}
												if !rules[ruleLPAREN]() {
													goto l505
												}
												{
													position543, tokenIndex543, depth543 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l543
													}
													goto l544
												l543:
													position, tokenIndex, depth = position543, tokenIndex543, depth543
												}
											l544:
												if !rules[ruleexpression]() {
													goto l505
												}
												{
													position545, tokenIndex545, depth545 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l545
													}
													{
														position547 := position
														depth++
														{
															position548, tokenIndex548, depth548 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l549
															}
															position++
															goto l548
														l549:
															position, tokenIndex, depth = position548, tokenIndex548, depth548
															if buffer[position] != rune('S') {
																goto l545
															}
															position++
														}
													l548:
														{
															position550, tokenIndex550, depth550 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l551
															}
															position++
															goto l550
														l551:
															position, tokenIndex, depth = position550, tokenIndex550, depth550
															if buffer[position] != rune('E') {
																goto l545
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
																goto l545
															}
															position++
														}
													l552:
														{
															position554, tokenIndex554, depth554 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l555
															}
															position++
															goto l554
														l555:
															position, tokenIndex, depth = position554, tokenIndex554, depth554
															if buffer[position] != rune('A') {
																goto l545
															}
															position++
														}
													l554:
														{
															position556, tokenIndex556, depth556 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l557
															}
															position++
															goto l556
														l557:
															position, tokenIndex, depth = position556, tokenIndex556, depth556
															if buffer[position] != rune('R') {
																goto l545
															}
															position++
														}
													l556:
														{
															position558, tokenIndex558, depth558 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l559
															}
															position++
															goto l558
														l559:
															position, tokenIndex, depth = position558, tokenIndex558, depth558
															if buffer[position] != rune('A') {
																goto l545
															}
															position++
														}
													l558:
														{
															position560, tokenIndex560, depth560 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l561
															}
															position++
															goto l560
														l561:
															position, tokenIndex, depth = position560, tokenIndex560, depth560
															if buffer[position] != rune('T') {
																goto l545
															}
															position++
														}
													l560:
														{
															position562, tokenIndex562, depth562 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l563
															}
															position++
															goto l562
														l563:
															position, tokenIndex, depth = position562, tokenIndex562, depth562
															if buffer[position] != rune('O') {
																goto l545
															}
															position++
														}
													l562:
														{
															position564, tokenIndex564, depth564 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l565
															}
															position++
															goto l564
														l565:
															position, tokenIndex, depth = position564, tokenIndex564, depth564
															if buffer[position] != rune('R') {
																goto l545
															}
															position++
														}
													l564:
														if !rules[ruleskip]() {
															goto l545
														}
														depth--
														add(ruleSEPARATOR, position547)
													}
													if !rules[ruleEQ]() {
														goto l545
													}
													if !rules[rulestring]() {
														goto l545
													}
													goto l546
												l545:
													position, tokenIndex, depth = position545, tokenIndex545, depth545
												}
											l546:
												if !rules[ruleRPAREN]() {
													goto l505
												}
												depth--
												add(rulegroupConcat, position519)
											}
											break
										case 'C', 'c':
											{
												position566 := position
												depth++
												{
													position567 := position
													depth++
													{
														position568, tokenIndex568, depth568 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l569
														}
														position++
														goto l568
													l569:
														position, tokenIndex, depth = position568, tokenIndex568, depth568
														if buffer[position] != rune('C') {
															goto l505
														}
														position++
													}
												l568:
													{
														position570, tokenIndex570, depth570 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l571
														}
														position++
														goto l570
													l571:
														position, tokenIndex, depth = position570, tokenIndex570, depth570
														if buffer[position] != rune('O') {
															goto l505
														}
														position++
													}
												l570:
													{
														position572, tokenIndex572, depth572 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l573
														}
														position++
														goto l572
													l573:
														position, tokenIndex, depth = position572, tokenIndex572, depth572
														if buffer[position] != rune('U') {
															goto l505
														}
														position++
													}
												l572:
													{
														position574, tokenIndex574, depth574 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l575
														}
														position++
														goto l574
													l575:
														position, tokenIndex, depth = position574, tokenIndex574, depth574
														if buffer[position] != rune('N') {
															goto l505
														}
														position++
													}
												l574:
													{
														position576, tokenIndex576, depth576 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l577
														}
														position++
														goto l576
													l577:
														position, tokenIndex, depth = position576, tokenIndex576, depth576
														if buffer[position] != rune('T') {
															goto l505
														}
														position++
													}
												l576:
													if !rules[ruleskip]() {
														goto l505
													}
													depth--
													add(ruleCOUNT, position567)
												}
												if !rules[ruleLPAREN]() {
													goto l505
												}
												{
													position578, tokenIndex578, depth578 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l578
													}
													goto l579
												l578:
													position, tokenIndex, depth = position578, tokenIndex578, depth578
												}
											l579:
												{
													position580, tokenIndex580, depth580 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l581
													}
													goto l580
												l581:
													position, tokenIndex, depth = position580, tokenIndex580, depth580
													if !rules[ruleexpression]() {
														goto l505
													}
												}
											l580:
												if !rules[ruleRPAREN]() {
													goto l505
												}
												depth--
												add(rulecount, position566)
											}
											break
										default:
											{
												position582, tokenIndex582, depth582 := position, tokenIndex, depth
												{
													position584 := position
													depth++
													{
														position585, tokenIndex585, depth585 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l586
														}
														position++
														goto l585
													l586:
														position, tokenIndex, depth = position585, tokenIndex585, depth585
														if buffer[position] != rune('S') {
															goto l583
														}
														position++
													}
												l585:
													{
														position587, tokenIndex587, depth587 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l588
														}
														position++
														goto l587
													l588:
														position, tokenIndex, depth = position587, tokenIndex587, depth587
														if buffer[position] != rune('U') {
															goto l583
														}
														position++
													}
												l587:
													{
														position589, tokenIndex589, depth589 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l590
														}
														position++
														goto l589
													l590:
														position, tokenIndex, depth = position589, tokenIndex589, depth589
														if buffer[position] != rune('M') {
															goto l583
														}
														position++
													}
												l589:
													if !rules[ruleskip]() {
														goto l583
													}
													depth--
													add(ruleSUM, position584)
												}
												goto l582
											l583:
												position, tokenIndex, depth = position582, tokenIndex582, depth582
												{
													position592 := position
													depth++
													{
														position593, tokenIndex593, depth593 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l594
														}
														position++
														goto l593
													l594:
														position, tokenIndex, depth = position593, tokenIndex593, depth593
														if buffer[position] != rune('M') {
															goto l591
														}
														position++
													}
												l593:
													{
														position595, tokenIndex595, depth595 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l596
														}
														position++
														goto l595
													l596:
														position, tokenIndex, depth = position595, tokenIndex595, depth595
														if buffer[position] != rune('I') {
															goto l591
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
															goto l591
														}
														position++
													}
												l597:
													if !rules[ruleskip]() {
														goto l591
													}
													depth--
													add(ruleMIN, position592)
												}
												goto l582
											l591:
												position, tokenIndex, depth = position582, tokenIndex582, depth582
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position600 := position
															depth++
															{
																position601, tokenIndex601, depth601 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l602
																}
																position++
																goto l601
															l602:
																position, tokenIndex, depth = position601, tokenIndex601, depth601
																if buffer[position] != rune('S') {
																	goto l505
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
																	goto l505
																}
																position++
															}
														l603:
															{
																position605, tokenIndex605, depth605 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l606
																}
																position++
																goto l605
															l606:
																position, tokenIndex, depth = position605, tokenIndex605, depth605
																if buffer[position] != rune('M') {
																	goto l505
																}
																position++
															}
														l605:
															{
																position607, tokenIndex607, depth607 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l608
																}
																position++
																goto l607
															l608:
																position, tokenIndex, depth = position607, tokenIndex607, depth607
																if buffer[position] != rune('P') {
																	goto l505
																}
																position++
															}
														l607:
															{
																position609, tokenIndex609, depth609 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l610
																}
																position++
																goto l609
															l610:
																position, tokenIndex, depth = position609, tokenIndex609, depth609
																if buffer[position] != rune('L') {
																	goto l505
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
																	goto l505
																}
																position++
															}
														l611:
															if !rules[ruleskip]() {
																goto l505
															}
															depth--
															add(ruleSAMPLE, position600)
														}
														break
													case 'A', 'a':
														{
															position613 := position
															depth++
															{
																position614, tokenIndex614, depth614 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l615
																}
																position++
																goto l614
															l615:
																position, tokenIndex, depth = position614, tokenIndex614, depth614
																if buffer[position] != rune('A') {
																	goto l505
																}
																position++
															}
														l614:
															{
																position616, tokenIndex616, depth616 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l617
																}
																position++
																goto l616
															l617:
																position, tokenIndex, depth = position616, tokenIndex616, depth616
																if buffer[position] != rune('V') {
																	goto l505
																}
																position++
															}
														l616:
															{
																position618, tokenIndex618, depth618 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l619
																}
																position++
																goto l618
															l619:
																position, tokenIndex, depth = position618, tokenIndex618, depth618
																if buffer[position] != rune('G') {
																	goto l505
																}
																position++
															}
														l618:
															if !rules[ruleskip]() {
																goto l505
															}
															depth--
															add(ruleAVG, position613)
														}
														break
													default:
														{
															position620 := position
															depth++
															{
																position621, tokenIndex621, depth621 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l622
																}
																position++
																goto l621
															l622:
																position, tokenIndex, depth = position621, tokenIndex621, depth621
																if buffer[position] != rune('M') {
																	goto l505
																}
																position++
															}
														l621:
															{
																position623, tokenIndex623, depth623 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l624
																}
																position++
																goto l623
															l624:
																position, tokenIndex, depth = position623, tokenIndex623, depth623
																if buffer[position] != rune('A') {
																	goto l505
																}
																position++
															}
														l623:
															{
																position625, tokenIndex625, depth625 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l626
																}
																position++
																goto l625
															l626:
																position, tokenIndex, depth = position625, tokenIndex625, depth625
																if buffer[position] != rune('X') {
																	goto l505
																}
																position++
															}
														l625:
															if !rules[ruleskip]() {
																goto l505
															}
															depth--
															add(ruleMAX, position620)
														}
														break
													}
												}

											}
										l582:
											if !rules[ruleLPAREN]() {
												goto l505
											}
											{
												position627, tokenIndex627, depth627 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l627
												}
												goto l628
											l627:
												position, tokenIndex, depth = position627, tokenIndex627, depth627
											}
										l628:
											if !rules[ruleexpression]() {
												goto l505
											}
											if !rules[ruleRPAREN]() {
												goto l505
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position517)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l505
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l505
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l505
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l505
								}
								break
							}
						}

					}
				l511:
					depth--
					add(ruleprimaryExpression, position510)
				}
				depth--
				add(ruleunaryExpression, position506)
			}
			return true
		l505:
			position, tokenIndex, depth = position505, tokenIndex505, depth505
			return false
		},
		/* 54 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 55 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position630, tokenIndex630, depth630 := position, tokenIndex, depth
			{
				position631 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l630
				}
				if !rules[ruleexpression]() {
					goto l630
				}
				if !rules[ruleRPAREN]() {
					goto l630
				}
				depth--
				add(rulebrackettedExpression, position631)
			}
			return true
		l630:
			position, tokenIndex, depth = position630, tokenIndex630, depth630
			return false
		},
		/* 56 functionCall <- <(iriref argList)> */
		func() bool {
			position632, tokenIndex632, depth632 := position, tokenIndex, depth
			{
				position633 := position
				depth++
				if !rules[ruleiriref]() {
					goto l632
				}
				if !rules[ruleargList]() {
					goto l632
				}
				depth--
				add(rulefunctionCall, position633)
			}
			return true
		l632:
			position, tokenIndex, depth = position632, tokenIndex632, depth632
			return false
		},
		/* 57 in <- <(IN argList)> */
		nil,
		/* 58 notin <- <(NOTIN argList)> */
		nil,
		/* 59 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position636, tokenIndex636, depth636 := position, tokenIndex, depth
			{
				position637 := position
				depth++
				{
					position638, tokenIndex638, depth638 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l639
					}
					goto l638
				l639:
					position, tokenIndex, depth = position638, tokenIndex638, depth638
					if !rules[ruleLPAREN]() {
						goto l636
					}
					if !rules[ruleexpression]() {
						goto l636
					}
				l640:
					{
						position641, tokenIndex641, depth641 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l641
						}
						if !rules[ruleexpression]() {
							goto l641
						}
						goto l640
					l641:
						position, tokenIndex, depth = position641, tokenIndex641, depth641
					}
					if !rules[ruleRPAREN]() {
						goto l636
					}
				}
			l638:
				depth--
				add(ruleargList, position637)
			}
			return true
		l636:
			position, tokenIndex, depth = position636, tokenIndex636, depth636
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
			position645, tokenIndex645, depth645 := position, tokenIndex, depth
			{
				position646 := position
				depth++
				{
					position647, tokenIndex647, depth647 := position, tokenIndex, depth
					{
						position649, tokenIndex649, depth649 := position, tokenIndex, depth
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
								if buffer[position] != rune('t') {
									goto l655
								}
								position++
								goto l654
							l655:
								position, tokenIndex, depth = position654, tokenIndex654, depth654
								if buffer[position] != rune('T') {
									goto l650
								}
								position++
							}
						l654:
							{
								position656, tokenIndex656, depth656 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l657
								}
								position++
								goto l656
							l657:
								position, tokenIndex, depth = position656, tokenIndex656, depth656
								if buffer[position] != rune('R') {
									goto l650
								}
								position++
							}
						l656:
							if !rules[ruleskip]() {
								goto l650
							}
							depth--
							add(ruleSTR, position651)
						}
						goto l649
					l650:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position659 := position
							depth++
							{
								position660, tokenIndex660, depth660 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l661
								}
								position++
								goto l660
							l661:
								position, tokenIndex, depth = position660, tokenIndex660, depth660
								if buffer[position] != rune('L') {
									goto l658
								}
								position++
							}
						l660:
							{
								position662, tokenIndex662, depth662 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l663
								}
								position++
								goto l662
							l663:
								position, tokenIndex, depth = position662, tokenIndex662, depth662
								if buffer[position] != rune('A') {
									goto l658
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
									goto l658
								}
								position++
							}
						l664:
							{
								position666, tokenIndex666, depth666 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l667
								}
								position++
								goto l666
							l667:
								position, tokenIndex, depth = position666, tokenIndex666, depth666
								if buffer[position] != rune('G') {
									goto l658
								}
								position++
							}
						l666:
							if !rules[ruleskip]() {
								goto l658
							}
							depth--
							add(ruleLANG, position659)
						}
						goto l649
					l658:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position669 := position
							depth++
							{
								position670, tokenIndex670, depth670 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l671
								}
								position++
								goto l670
							l671:
								position, tokenIndex, depth = position670, tokenIndex670, depth670
								if buffer[position] != rune('D') {
									goto l668
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
									goto l668
								}
								position++
							}
						l672:
							{
								position674, tokenIndex674, depth674 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l675
								}
								position++
								goto l674
							l675:
								position, tokenIndex, depth = position674, tokenIndex674, depth674
								if buffer[position] != rune('T') {
									goto l668
								}
								position++
							}
						l674:
							{
								position676, tokenIndex676, depth676 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l677
								}
								position++
								goto l676
							l677:
								position, tokenIndex, depth = position676, tokenIndex676, depth676
								if buffer[position] != rune('A') {
									goto l668
								}
								position++
							}
						l676:
							{
								position678, tokenIndex678, depth678 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l679
								}
								position++
								goto l678
							l679:
								position, tokenIndex, depth = position678, tokenIndex678, depth678
								if buffer[position] != rune('T') {
									goto l668
								}
								position++
							}
						l678:
							{
								position680, tokenIndex680, depth680 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l681
								}
								position++
								goto l680
							l681:
								position, tokenIndex, depth = position680, tokenIndex680, depth680
								if buffer[position] != rune('Y') {
									goto l668
								}
								position++
							}
						l680:
							{
								position682, tokenIndex682, depth682 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l683
								}
								position++
								goto l682
							l683:
								position, tokenIndex, depth = position682, tokenIndex682, depth682
								if buffer[position] != rune('P') {
									goto l668
								}
								position++
							}
						l682:
							{
								position684, tokenIndex684, depth684 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l685
								}
								position++
								goto l684
							l685:
								position, tokenIndex, depth = position684, tokenIndex684, depth684
								if buffer[position] != rune('E') {
									goto l668
								}
								position++
							}
						l684:
							if !rules[ruleskip]() {
								goto l668
							}
							depth--
							add(ruleDATATYPE, position669)
						}
						goto l649
					l668:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
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
								if buffer[position] != rune('r') {
									goto l691
								}
								position++
								goto l690
							l691:
								position, tokenIndex, depth = position690, tokenIndex690, depth690
								if buffer[position] != rune('R') {
									goto l686
								}
								position++
							}
						l690:
							{
								position692, tokenIndex692, depth692 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l693
								}
								position++
								goto l692
							l693:
								position, tokenIndex, depth = position692, tokenIndex692, depth692
								if buffer[position] != rune('I') {
									goto l686
								}
								position++
							}
						l692:
							if !rules[ruleskip]() {
								goto l686
							}
							depth--
							add(ruleIRI, position687)
						}
						goto l649
					l686:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position695 := position
							depth++
							{
								position696, tokenIndex696, depth696 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l697
								}
								position++
								goto l696
							l697:
								position, tokenIndex, depth = position696, tokenIndex696, depth696
								if buffer[position] != rune('U') {
									goto l694
								}
								position++
							}
						l696:
							{
								position698, tokenIndex698, depth698 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l699
								}
								position++
								goto l698
							l699:
								position, tokenIndex, depth = position698, tokenIndex698, depth698
								if buffer[position] != rune('R') {
									goto l694
								}
								position++
							}
						l698:
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
									goto l694
								}
								position++
							}
						l700:
							if !rules[ruleskip]() {
								goto l694
							}
							depth--
							add(ruleURI, position695)
						}
						goto l649
					l694:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position703 := position
							depth++
							{
								position704, tokenIndex704, depth704 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l705
								}
								position++
								goto l704
							l705:
								position, tokenIndex, depth = position704, tokenIndex704, depth704
								if buffer[position] != rune('S') {
									goto l702
								}
								position++
							}
						l704:
							{
								position706, tokenIndex706, depth706 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l707
								}
								position++
								goto l706
							l707:
								position, tokenIndex, depth = position706, tokenIndex706, depth706
								if buffer[position] != rune('T') {
									goto l702
								}
								position++
							}
						l706:
							{
								position708, tokenIndex708, depth708 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l709
								}
								position++
								goto l708
							l709:
								position, tokenIndex, depth = position708, tokenIndex708, depth708
								if buffer[position] != rune('R') {
									goto l702
								}
								position++
							}
						l708:
							{
								position710, tokenIndex710, depth710 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l711
								}
								position++
								goto l710
							l711:
								position, tokenIndex, depth = position710, tokenIndex710, depth710
								if buffer[position] != rune('L') {
									goto l702
								}
								position++
							}
						l710:
							{
								position712, tokenIndex712, depth712 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l713
								}
								position++
								goto l712
							l713:
								position, tokenIndex, depth = position712, tokenIndex712, depth712
								if buffer[position] != rune('E') {
									goto l702
								}
								position++
							}
						l712:
							{
								position714, tokenIndex714, depth714 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l715
								}
								position++
								goto l714
							l715:
								position, tokenIndex, depth = position714, tokenIndex714, depth714
								if buffer[position] != rune('N') {
									goto l702
								}
								position++
							}
						l714:
							if !rules[ruleskip]() {
								goto l702
							}
							depth--
							add(ruleSTRLEN, position703)
						}
						goto l649
					l702:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position717 := position
							depth++
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
									goto l716
								}
								position++
							}
						l718:
							{
								position720, tokenIndex720, depth720 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l721
								}
								position++
								goto l720
							l721:
								position, tokenIndex, depth = position720, tokenIndex720, depth720
								if buffer[position] != rune('O') {
									goto l716
								}
								position++
							}
						l720:
							{
								position722, tokenIndex722, depth722 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l723
								}
								position++
								goto l722
							l723:
								position, tokenIndex, depth = position722, tokenIndex722, depth722
								if buffer[position] != rune('N') {
									goto l716
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
									goto l716
								}
								position++
							}
						l724:
							{
								position726, tokenIndex726, depth726 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l727
								}
								position++
								goto l726
							l727:
								position, tokenIndex, depth = position726, tokenIndex726, depth726
								if buffer[position] != rune('H') {
									goto l716
								}
								position++
							}
						l726:
							if !rules[ruleskip]() {
								goto l716
							}
							depth--
							add(ruleMONTH, position717)
						}
						goto l649
					l716:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position729 := position
							depth++
							{
								position730, tokenIndex730, depth730 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l731
								}
								position++
								goto l730
							l731:
								position, tokenIndex, depth = position730, tokenIndex730, depth730
								if buffer[position] != rune('M') {
									goto l728
								}
								position++
							}
						l730:
							{
								position732, tokenIndex732, depth732 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l733
								}
								position++
								goto l732
							l733:
								position, tokenIndex, depth = position732, tokenIndex732, depth732
								if buffer[position] != rune('I') {
									goto l728
								}
								position++
							}
						l732:
							{
								position734, tokenIndex734, depth734 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l735
								}
								position++
								goto l734
							l735:
								position, tokenIndex, depth = position734, tokenIndex734, depth734
								if buffer[position] != rune('N') {
									goto l728
								}
								position++
							}
						l734:
							{
								position736, tokenIndex736, depth736 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l737
								}
								position++
								goto l736
							l737:
								position, tokenIndex, depth = position736, tokenIndex736, depth736
								if buffer[position] != rune('U') {
									goto l728
								}
								position++
							}
						l736:
							{
								position738, tokenIndex738, depth738 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l739
								}
								position++
								goto l738
							l739:
								position, tokenIndex, depth = position738, tokenIndex738, depth738
								if buffer[position] != rune('T') {
									goto l728
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
									goto l728
								}
								position++
							}
						l740:
							{
								position742, tokenIndex742, depth742 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l743
								}
								position++
								goto l742
							l743:
								position, tokenIndex, depth = position742, tokenIndex742, depth742
								if buffer[position] != rune('S') {
									goto l728
								}
								position++
							}
						l742:
							if !rules[ruleskip]() {
								goto l728
							}
							depth--
							add(ruleMINUTES, position729)
						}
						goto l649
					l728:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position745 := position
							depth++
							{
								position746, tokenIndex746, depth746 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l747
								}
								position++
								goto l746
							l747:
								position, tokenIndex, depth = position746, tokenIndex746, depth746
								if buffer[position] != rune('S') {
									goto l744
								}
								position++
							}
						l746:
							{
								position748, tokenIndex748, depth748 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l749
								}
								position++
								goto l748
							l749:
								position, tokenIndex, depth = position748, tokenIndex748, depth748
								if buffer[position] != rune('E') {
									goto l744
								}
								position++
							}
						l748:
							{
								position750, tokenIndex750, depth750 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l751
								}
								position++
								goto l750
							l751:
								position, tokenIndex, depth = position750, tokenIndex750, depth750
								if buffer[position] != rune('C') {
									goto l744
								}
								position++
							}
						l750:
							{
								position752, tokenIndex752, depth752 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l753
								}
								position++
								goto l752
							l753:
								position, tokenIndex, depth = position752, tokenIndex752, depth752
								if buffer[position] != rune('O') {
									goto l744
								}
								position++
							}
						l752:
							{
								position754, tokenIndex754, depth754 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l755
								}
								position++
								goto l754
							l755:
								position, tokenIndex, depth = position754, tokenIndex754, depth754
								if buffer[position] != rune('N') {
									goto l744
								}
								position++
							}
						l754:
							{
								position756, tokenIndex756, depth756 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l757
								}
								position++
								goto l756
							l757:
								position, tokenIndex, depth = position756, tokenIndex756, depth756
								if buffer[position] != rune('D') {
									goto l744
								}
								position++
							}
						l756:
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
									goto l744
								}
								position++
							}
						l758:
							if !rules[ruleskip]() {
								goto l744
							}
							depth--
							add(ruleSECONDS, position745)
						}
						goto l649
					l744:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position761 := position
							depth++
							{
								position762, tokenIndex762, depth762 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l763
								}
								position++
								goto l762
							l763:
								position, tokenIndex, depth = position762, tokenIndex762, depth762
								if buffer[position] != rune('T') {
									goto l760
								}
								position++
							}
						l762:
							{
								position764, tokenIndex764, depth764 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l765
								}
								position++
								goto l764
							l765:
								position, tokenIndex, depth = position764, tokenIndex764, depth764
								if buffer[position] != rune('I') {
									goto l760
								}
								position++
							}
						l764:
							{
								position766, tokenIndex766, depth766 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l767
								}
								position++
								goto l766
							l767:
								position, tokenIndex, depth = position766, tokenIndex766, depth766
								if buffer[position] != rune('M') {
									goto l760
								}
								position++
							}
						l766:
							{
								position768, tokenIndex768, depth768 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l769
								}
								position++
								goto l768
							l769:
								position, tokenIndex, depth = position768, tokenIndex768, depth768
								if buffer[position] != rune('E') {
									goto l760
								}
								position++
							}
						l768:
							{
								position770, tokenIndex770, depth770 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l771
								}
								position++
								goto l770
							l771:
								position, tokenIndex, depth = position770, tokenIndex770, depth770
								if buffer[position] != rune('Z') {
									goto l760
								}
								position++
							}
						l770:
							{
								position772, tokenIndex772, depth772 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l773
								}
								position++
								goto l772
							l773:
								position, tokenIndex, depth = position772, tokenIndex772, depth772
								if buffer[position] != rune('O') {
									goto l760
								}
								position++
							}
						l772:
							{
								position774, tokenIndex774, depth774 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l775
								}
								position++
								goto l774
							l775:
								position, tokenIndex, depth = position774, tokenIndex774, depth774
								if buffer[position] != rune('N') {
									goto l760
								}
								position++
							}
						l774:
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
									goto l760
								}
								position++
							}
						l776:
							if !rules[ruleskip]() {
								goto l760
							}
							depth--
							add(ruleTIMEZONE, position761)
						}
						goto l649
					l760:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position779 := position
							depth++
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
									goto l778
								}
								position++
							}
						l780:
							{
								position782, tokenIndex782, depth782 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l783
								}
								position++
								goto l782
							l783:
								position, tokenIndex, depth = position782, tokenIndex782, depth782
								if buffer[position] != rune('H') {
									goto l778
								}
								position++
							}
						l782:
							{
								position784, tokenIndex784, depth784 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l785
								}
								position++
								goto l784
							l785:
								position, tokenIndex, depth = position784, tokenIndex784, depth784
								if buffer[position] != rune('A') {
									goto l778
								}
								position++
							}
						l784:
							if buffer[position] != rune('1') {
								goto l778
							}
							position++
							if !rules[ruleskip]() {
								goto l778
							}
							depth--
							add(ruleSHA1, position779)
						}
						goto l649
					l778:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position787 := position
							depth++
							{
								position788, tokenIndex788, depth788 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l789
								}
								position++
								goto l788
							l789:
								position, tokenIndex, depth = position788, tokenIndex788, depth788
								if buffer[position] != rune('S') {
									goto l786
								}
								position++
							}
						l788:
							{
								position790, tokenIndex790, depth790 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l791
								}
								position++
								goto l790
							l791:
								position, tokenIndex, depth = position790, tokenIndex790, depth790
								if buffer[position] != rune('H') {
									goto l786
								}
								position++
							}
						l790:
							{
								position792, tokenIndex792, depth792 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l793
								}
								position++
								goto l792
							l793:
								position, tokenIndex, depth = position792, tokenIndex792, depth792
								if buffer[position] != rune('A') {
									goto l786
								}
								position++
							}
						l792:
							if buffer[position] != rune('2') {
								goto l786
							}
							position++
							if buffer[position] != rune('5') {
								goto l786
							}
							position++
							if buffer[position] != rune('6') {
								goto l786
							}
							position++
							if !rules[ruleskip]() {
								goto l786
							}
							depth--
							add(ruleSHA256, position787)
						}
						goto l649
					l786:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position795 := position
							depth++
							{
								position796, tokenIndex796, depth796 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l797
								}
								position++
								goto l796
							l797:
								position, tokenIndex, depth = position796, tokenIndex796, depth796
								if buffer[position] != rune('S') {
									goto l794
								}
								position++
							}
						l796:
							{
								position798, tokenIndex798, depth798 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l799
								}
								position++
								goto l798
							l799:
								position, tokenIndex, depth = position798, tokenIndex798, depth798
								if buffer[position] != rune('H') {
									goto l794
								}
								position++
							}
						l798:
							{
								position800, tokenIndex800, depth800 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l801
								}
								position++
								goto l800
							l801:
								position, tokenIndex, depth = position800, tokenIndex800, depth800
								if buffer[position] != rune('A') {
									goto l794
								}
								position++
							}
						l800:
							if buffer[position] != rune('3') {
								goto l794
							}
							position++
							if buffer[position] != rune('8') {
								goto l794
							}
							position++
							if buffer[position] != rune('4') {
								goto l794
							}
							position++
							if !rules[ruleskip]() {
								goto l794
							}
							depth--
							add(ruleSHA384, position795)
						}
						goto l649
					l794:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position803 := position
							depth++
							{
								position804, tokenIndex804, depth804 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l805
								}
								position++
								goto l804
							l805:
								position, tokenIndex, depth = position804, tokenIndex804, depth804
								if buffer[position] != rune('I') {
									goto l802
								}
								position++
							}
						l804:
							{
								position806, tokenIndex806, depth806 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l807
								}
								position++
								goto l806
							l807:
								position, tokenIndex, depth = position806, tokenIndex806, depth806
								if buffer[position] != rune('S') {
									goto l802
								}
								position++
							}
						l806:
							{
								position808, tokenIndex808, depth808 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l809
								}
								position++
								goto l808
							l809:
								position, tokenIndex, depth = position808, tokenIndex808, depth808
								if buffer[position] != rune('I') {
									goto l802
								}
								position++
							}
						l808:
							{
								position810, tokenIndex810, depth810 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l811
								}
								position++
								goto l810
							l811:
								position, tokenIndex, depth = position810, tokenIndex810, depth810
								if buffer[position] != rune('R') {
									goto l802
								}
								position++
							}
						l810:
							{
								position812, tokenIndex812, depth812 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l813
								}
								position++
								goto l812
							l813:
								position, tokenIndex, depth = position812, tokenIndex812, depth812
								if buffer[position] != rune('I') {
									goto l802
								}
								position++
							}
						l812:
							if !rules[ruleskip]() {
								goto l802
							}
							depth--
							add(ruleISIRI, position803)
						}
						goto l649
					l802:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position815 := position
							depth++
							{
								position816, tokenIndex816, depth816 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l817
								}
								position++
								goto l816
							l817:
								position, tokenIndex, depth = position816, tokenIndex816, depth816
								if buffer[position] != rune('I') {
									goto l814
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
									goto l814
								}
								position++
							}
						l818:
							{
								position820, tokenIndex820, depth820 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l821
								}
								position++
								goto l820
							l821:
								position, tokenIndex, depth = position820, tokenIndex820, depth820
								if buffer[position] != rune('U') {
									goto l814
								}
								position++
							}
						l820:
							{
								position822, tokenIndex822, depth822 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l823
								}
								position++
								goto l822
							l823:
								position, tokenIndex, depth = position822, tokenIndex822, depth822
								if buffer[position] != rune('R') {
									goto l814
								}
								position++
							}
						l822:
							{
								position824, tokenIndex824, depth824 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l825
								}
								position++
								goto l824
							l825:
								position, tokenIndex, depth = position824, tokenIndex824, depth824
								if buffer[position] != rune('I') {
									goto l814
								}
								position++
							}
						l824:
							if !rules[ruleskip]() {
								goto l814
							}
							depth--
							add(ruleISURI, position815)
						}
						goto l649
					l814:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position827 := position
							depth++
							{
								position828, tokenIndex828, depth828 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l829
								}
								position++
								goto l828
							l829:
								position, tokenIndex, depth = position828, tokenIndex828, depth828
								if buffer[position] != rune('I') {
									goto l826
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
									goto l826
								}
								position++
							}
						l830:
							{
								position832, tokenIndex832, depth832 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l833
								}
								position++
								goto l832
							l833:
								position, tokenIndex, depth = position832, tokenIndex832, depth832
								if buffer[position] != rune('B') {
									goto l826
								}
								position++
							}
						l832:
							{
								position834, tokenIndex834, depth834 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l835
								}
								position++
								goto l834
							l835:
								position, tokenIndex, depth = position834, tokenIndex834, depth834
								if buffer[position] != rune('L') {
									goto l826
								}
								position++
							}
						l834:
							{
								position836, tokenIndex836, depth836 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l837
								}
								position++
								goto l836
							l837:
								position, tokenIndex, depth = position836, tokenIndex836, depth836
								if buffer[position] != rune('A') {
									goto l826
								}
								position++
							}
						l836:
							{
								position838, tokenIndex838, depth838 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l839
								}
								position++
								goto l838
							l839:
								position, tokenIndex, depth = position838, tokenIndex838, depth838
								if buffer[position] != rune('N') {
									goto l826
								}
								position++
							}
						l838:
							{
								position840, tokenIndex840, depth840 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l841
								}
								position++
								goto l840
							l841:
								position, tokenIndex, depth = position840, tokenIndex840, depth840
								if buffer[position] != rune('K') {
									goto l826
								}
								position++
							}
						l840:
							if !rules[ruleskip]() {
								goto l826
							}
							depth--
							add(ruleISBLANK, position827)
						}
						goto l649
					l826:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							position843 := position
							depth++
							{
								position844, tokenIndex844, depth844 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l845
								}
								position++
								goto l844
							l845:
								position, tokenIndex, depth = position844, tokenIndex844, depth844
								if buffer[position] != rune('I') {
									goto l842
								}
								position++
							}
						l844:
							{
								position846, tokenIndex846, depth846 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l847
								}
								position++
								goto l846
							l847:
								position, tokenIndex, depth = position846, tokenIndex846, depth846
								if buffer[position] != rune('S') {
									goto l842
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
									goto l842
								}
								position++
							}
						l848:
							{
								position850, tokenIndex850, depth850 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l851
								}
								position++
								goto l850
							l851:
								position, tokenIndex, depth = position850, tokenIndex850, depth850
								if buffer[position] != rune('I') {
									goto l842
								}
								position++
							}
						l850:
							{
								position852, tokenIndex852, depth852 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l853
								}
								position++
								goto l852
							l853:
								position, tokenIndex, depth = position852, tokenIndex852, depth852
								if buffer[position] != rune('T') {
									goto l842
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
									goto l842
								}
								position++
							}
						l854:
							{
								position856, tokenIndex856, depth856 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l857
								}
								position++
								goto l856
							l857:
								position, tokenIndex, depth = position856, tokenIndex856, depth856
								if buffer[position] != rune('R') {
									goto l842
								}
								position++
							}
						l856:
							{
								position858, tokenIndex858, depth858 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('A') {
									goto l842
								}
								position++
							}
						l858:
							{
								position860, tokenIndex860, depth860 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l861
								}
								position++
								goto l860
							l861:
								position, tokenIndex, depth = position860, tokenIndex860, depth860
								if buffer[position] != rune('L') {
									goto l842
								}
								position++
							}
						l860:
							if !rules[ruleskip]() {
								goto l842
							}
							depth--
							add(ruleISLITERAL, position843)
						}
						goto l649
					l842:
						position, tokenIndex, depth = position649, tokenIndex649, depth649
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position863 := position
									depth++
									{
										position864, tokenIndex864, depth864 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l865
										}
										position++
										goto l864
									l865:
										position, tokenIndex, depth = position864, tokenIndex864, depth864
										if buffer[position] != rune('I') {
											goto l648
										}
										position++
									}
								l864:
									{
										position866, tokenIndex866, depth866 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l867
										}
										position++
										goto l866
									l867:
										position, tokenIndex, depth = position866, tokenIndex866, depth866
										if buffer[position] != rune('S') {
											goto l648
										}
										position++
									}
								l866:
									{
										position868, tokenIndex868, depth868 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l869
										}
										position++
										goto l868
									l869:
										position, tokenIndex, depth = position868, tokenIndex868, depth868
										if buffer[position] != rune('N') {
											goto l648
										}
										position++
									}
								l868:
									{
										position870, tokenIndex870, depth870 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l871
										}
										position++
										goto l870
									l871:
										position, tokenIndex, depth = position870, tokenIndex870, depth870
										if buffer[position] != rune('U') {
											goto l648
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
											goto l648
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
											goto l648
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
											goto l648
										}
										position++
									}
								l876:
									{
										position878, tokenIndex878, depth878 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l879
										}
										position++
										goto l878
									l879:
										position, tokenIndex, depth = position878, tokenIndex878, depth878
										if buffer[position] != rune('I') {
											goto l648
										}
										position++
									}
								l878:
									{
										position880, tokenIndex880, depth880 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l881
										}
										position++
										goto l880
									l881:
										position, tokenIndex, depth = position880, tokenIndex880, depth880
										if buffer[position] != rune('C') {
											goto l648
										}
										position++
									}
								l880:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleISNUMERIC, position863)
								}
								break
							case 'S', 's':
								{
									position882 := position
									depth++
									{
										position883, tokenIndex883, depth883 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l884
										}
										position++
										goto l883
									l884:
										position, tokenIndex, depth = position883, tokenIndex883, depth883
										if buffer[position] != rune('S') {
											goto l648
										}
										position++
									}
								l883:
									{
										position885, tokenIndex885, depth885 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l886
										}
										position++
										goto l885
									l886:
										position, tokenIndex, depth = position885, tokenIndex885, depth885
										if buffer[position] != rune('H') {
											goto l648
										}
										position++
									}
								l885:
									{
										position887, tokenIndex887, depth887 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l888
										}
										position++
										goto l887
									l888:
										position, tokenIndex, depth = position887, tokenIndex887, depth887
										if buffer[position] != rune('A') {
											goto l648
										}
										position++
									}
								l887:
									if buffer[position] != rune('5') {
										goto l648
									}
									position++
									if buffer[position] != rune('1') {
										goto l648
									}
									position++
									if buffer[position] != rune('2') {
										goto l648
									}
									position++
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleSHA512, position882)
								}
								break
							case 'M', 'm':
								{
									position889 := position
									depth++
									{
										position890, tokenIndex890, depth890 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l891
										}
										position++
										goto l890
									l891:
										position, tokenIndex, depth = position890, tokenIndex890, depth890
										if buffer[position] != rune('M') {
											goto l648
										}
										position++
									}
								l890:
									{
										position892, tokenIndex892, depth892 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l893
										}
										position++
										goto l892
									l893:
										position, tokenIndex, depth = position892, tokenIndex892, depth892
										if buffer[position] != rune('D') {
											goto l648
										}
										position++
									}
								l892:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleMD5, position889)
								}
								break
							case 'T', 't':
								{
									position894 := position
									depth++
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
											goto l648
										}
										position++
									}
								l895:
									{
										position897, tokenIndex897, depth897 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l898
										}
										position++
										goto l897
									l898:
										position, tokenIndex, depth = position897, tokenIndex897, depth897
										if buffer[position] != rune('Z') {
											goto l648
										}
										position++
									}
								l897:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleTZ, position894)
								}
								break
							case 'H', 'h':
								{
									position899 := position
									depth++
									{
										position900, tokenIndex900, depth900 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l901
										}
										position++
										goto l900
									l901:
										position, tokenIndex, depth = position900, tokenIndex900, depth900
										if buffer[position] != rune('H') {
											goto l648
										}
										position++
									}
								l900:
									{
										position902, tokenIndex902, depth902 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l903
										}
										position++
										goto l902
									l903:
										position, tokenIndex, depth = position902, tokenIndex902, depth902
										if buffer[position] != rune('O') {
											goto l648
										}
										position++
									}
								l902:
									{
										position904, tokenIndex904, depth904 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l905
										}
										position++
										goto l904
									l905:
										position, tokenIndex, depth = position904, tokenIndex904, depth904
										if buffer[position] != rune('U') {
											goto l648
										}
										position++
									}
								l904:
									{
										position906, tokenIndex906, depth906 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l907
										}
										position++
										goto l906
									l907:
										position, tokenIndex, depth = position906, tokenIndex906, depth906
										if buffer[position] != rune('R') {
											goto l648
										}
										position++
									}
								l906:
									{
										position908, tokenIndex908, depth908 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l909
										}
										position++
										goto l908
									l909:
										position, tokenIndex, depth = position908, tokenIndex908, depth908
										if buffer[position] != rune('S') {
											goto l648
										}
										position++
									}
								l908:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleHOURS, position899)
								}
								break
							case 'D', 'd':
								{
									position910 := position
									depth++
									{
										position911, tokenIndex911, depth911 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l912
										}
										position++
										goto l911
									l912:
										position, tokenIndex, depth = position911, tokenIndex911, depth911
										if buffer[position] != rune('D') {
											goto l648
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
											goto l648
										}
										position++
									}
								l913:
									{
										position915, tokenIndex915, depth915 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l916
										}
										position++
										goto l915
									l916:
										position, tokenIndex, depth = position915, tokenIndex915, depth915
										if buffer[position] != rune('Y') {
											goto l648
										}
										position++
									}
								l915:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleDAY, position910)
								}
								break
							case 'Y', 'y':
								{
									position917 := position
									depth++
									{
										position918, tokenIndex918, depth918 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l919
										}
										position++
										goto l918
									l919:
										position, tokenIndex, depth = position918, tokenIndex918, depth918
										if buffer[position] != rune('Y') {
											goto l648
										}
										position++
									}
								l918:
									{
										position920, tokenIndex920, depth920 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l921
										}
										position++
										goto l920
									l921:
										position, tokenIndex, depth = position920, tokenIndex920, depth920
										if buffer[position] != rune('E') {
											goto l648
										}
										position++
									}
								l920:
									{
										position922, tokenIndex922, depth922 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l923
										}
										position++
										goto l922
									l923:
										position, tokenIndex, depth = position922, tokenIndex922, depth922
										if buffer[position] != rune('A') {
											goto l648
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
											goto l648
										}
										position++
									}
								l924:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleYEAR, position917)
								}
								break
							case 'E', 'e':
								{
									position926 := position
									depth++
									{
										position927, tokenIndex927, depth927 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l928
										}
										position++
										goto l927
									l928:
										position, tokenIndex, depth = position927, tokenIndex927, depth927
										if buffer[position] != rune('E') {
											goto l648
										}
										position++
									}
								l927:
									{
										position929, tokenIndex929, depth929 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l930
										}
										position++
										goto l929
									l930:
										position, tokenIndex, depth = position929, tokenIndex929, depth929
										if buffer[position] != rune('N') {
											goto l648
										}
										position++
									}
								l929:
									{
										position931, tokenIndex931, depth931 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l932
										}
										position++
										goto l931
									l932:
										position, tokenIndex, depth = position931, tokenIndex931, depth931
										if buffer[position] != rune('C') {
											goto l648
										}
										position++
									}
								l931:
									{
										position933, tokenIndex933, depth933 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l934
										}
										position++
										goto l933
									l934:
										position, tokenIndex, depth = position933, tokenIndex933, depth933
										if buffer[position] != rune('O') {
											goto l648
										}
										position++
									}
								l933:
									{
										position935, tokenIndex935, depth935 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l936
										}
										position++
										goto l935
									l936:
										position, tokenIndex, depth = position935, tokenIndex935, depth935
										if buffer[position] != rune('D') {
											goto l648
										}
										position++
									}
								l935:
									{
										position937, tokenIndex937, depth937 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l938
										}
										position++
										goto l937
									l938:
										position, tokenIndex, depth = position937, tokenIndex937, depth937
										if buffer[position] != rune('E') {
											goto l648
										}
										position++
									}
								l937:
									if buffer[position] != rune('_') {
										goto l648
									}
									position++
									{
										position939, tokenIndex939, depth939 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l940
										}
										position++
										goto l939
									l940:
										position, tokenIndex, depth = position939, tokenIndex939, depth939
										if buffer[position] != rune('F') {
											goto l648
										}
										position++
									}
								l939:
									{
										position941, tokenIndex941, depth941 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l942
										}
										position++
										goto l941
									l942:
										position, tokenIndex, depth = position941, tokenIndex941, depth941
										if buffer[position] != rune('O') {
											goto l648
										}
										position++
									}
								l941:
									{
										position943, tokenIndex943, depth943 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l944
										}
										position++
										goto l943
									l944:
										position, tokenIndex, depth = position943, tokenIndex943, depth943
										if buffer[position] != rune('R') {
											goto l648
										}
										position++
									}
								l943:
									if buffer[position] != rune('_') {
										goto l648
									}
									position++
									{
										position945, tokenIndex945, depth945 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l946
										}
										position++
										goto l945
									l946:
										position, tokenIndex, depth = position945, tokenIndex945, depth945
										if buffer[position] != rune('U') {
											goto l648
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
											goto l648
										}
										position++
									}
								l947:
									{
										position949, tokenIndex949, depth949 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l950
										}
										position++
										goto l949
									l950:
										position, tokenIndex, depth = position949, tokenIndex949, depth949
										if buffer[position] != rune('I') {
											goto l648
										}
										position++
									}
								l949:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleENCODEFORURI, position926)
								}
								break
							case 'L', 'l':
								{
									position951 := position
									depth++
									{
										position952, tokenIndex952, depth952 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l953
										}
										position++
										goto l952
									l953:
										position, tokenIndex, depth = position952, tokenIndex952, depth952
										if buffer[position] != rune('L') {
											goto l648
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
											goto l648
										}
										position++
									}
								l954:
									{
										position956, tokenIndex956, depth956 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l957
										}
										position++
										goto l956
									l957:
										position, tokenIndex, depth = position956, tokenIndex956, depth956
										if buffer[position] != rune('A') {
											goto l648
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
											goto l648
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
											goto l648
										}
										position++
									}
								l960:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleLCASE, position951)
								}
								break
							case 'U', 'u':
								{
									position962 := position
									depth++
									{
										position963, tokenIndex963, depth963 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l964
										}
										position++
										goto l963
									l964:
										position, tokenIndex, depth = position963, tokenIndex963, depth963
										if buffer[position] != rune('U') {
											goto l648
										}
										position++
									}
								l963:
									{
										position965, tokenIndex965, depth965 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l966
										}
										position++
										goto l965
									l966:
										position, tokenIndex, depth = position965, tokenIndex965, depth965
										if buffer[position] != rune('C') {
											goto l648
										}
										position++
									}
								l965:
									{
										position967, tokenIndex967, depth967 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l968
										}
										position++
										goto l967
									l968:
										position, tokenIndex, depth = position967, tokenIndex967, depth967
										if buffer[position] != rune('A') {
											goto l648
										}
										position++
									}
								l967:
									{
										position969, tokenIndex969, depth969 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l970
										}
										position++
										goto l969
									l970:
										position, tokenIndex, depth = position969, tokenIndex969, depth969
										if buffer[position] != rune('S') {
											goto l648
										}
										position++
									}
								l969:
									{
										position971, tokenIndex971, depth971 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l972
										}
										position++
										goto l971
									l972:
										position, tokenIndex, depth = position971, tokenIndex971, depth971
										if buffer[position] != rune('E') {
											goto l648
										}
										position++
									}
								l971:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleUCASE, position962)
								}
								break
							case 'F', 'f':
								{
									position973 := position
									depth++
									{
										position974, tokenIndex974, depth974 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l975
										}
										position++
										goto l974
									l975:
										position, tokenIndex, depth = position974, tokenIndex974, depth974
										if buffer[position] != rune('F') {
											goto l648
										}
										position++
									}
								l974:
									{
										position976, tokenIndex976, depth976 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l977
										}
										position++
										goto l976
									l977:
										position, tokenIndex, depth = position976, tokenIndex976, depth976
										if buffer[position] != rune('L') {
											goto l648
										}
										position++
									}
								l976:
									{
										position978, tokenIndex978, depth978 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l979
										}
										position++
										goto l978
									l979:
										position, tokenIndex, depth = position978, tokenIndex978, depth978
										if buffer[position] != rune('O') {
											goto l648
										}
										position++
									}
								l978:
									{
										position980, tokenIndex980, depth980 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l981
										}
										position++
										goto l980
									l981:
										position, tokenIndex, depth = position980, tokenIndex980, depth980
										if buffer[position] != rune('O') {
											goto l648
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
											goto l648
										}
										position++
									}
								l982:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleFLOOR, position973)
								}
								break
							case 'R', 'r':
								{
									position984 := position
									depth++
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
											goto l648
										}
										position++
									}
								l985:
									{
										position987, tokenIndex987, depth987 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l988
										}
										position++
										goto l987
									l988:
										position, tokenIndex, depth = position987, tokenIndex987, depth987
										if buffer[position] != rune('O') {
											goto l648
										}
										position++
									}
								l987:
									{
										position989, tokenIndex989, depth989 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l990
										}
										position++
										goto l989
									l990:
										position, tokenIndex, depth = position989, tokenIndex989, depth989
										if buffer[position] != rune('U') {
											goto l648
										}
										position++
									}
								l989:
									{
										position991, tokenIndex991, depth991 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l992
										}
										position++
										goto l991
									l992:
										position, tokenIndex, depth = position991, tokenIndex991, depth991
										if buffer[position] != rune('N') {
											goto l648
										}
										position++
									}
								l991:
									{
										position993, tokenIndex993, depth993 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l994
										}
										position++
										goto l993
									l994:
										position, tokenIndex, depth = position993, tokenIndex993, depth993
										if buffer[position] != rune('D') {
											goto l648
										}
										position++
									}
								l993:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleROUND, position984)
								}
								break
							case 'C', 'c':
								{
									position995 := position
									depth++
									{
										position996, tokenIndex996, depth996 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l997
										}
										position++
										goto l996
									l997:
										position, tokenIndex, depth = position996, tokenIndex996, depth996
										if buffer[position] != rune('C') {
											goto l648
										}
										position++
									}
								l996:
									{
										position998, tokenIndex998, depth998 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l999
										}
										position++
										goto l998
									l999:
										position, tokenIndex, depth = position998, tokenIndex998, depth998
										if buffer[position] != rune('E') {
											goto l648
										}
										position++
									}
								l998:
									{
										position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1001
										}
										position++
										goto l1000
									l1001:
										position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
										if buffer[position] != rune('I') {
											goto l648
										}
										position++
									}
								l1000:
									{
										position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1003
										}
										position++
										goto l1002
									l1003:
										position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
										if buffer[position] != rune('L') {
											goto l648
										}
										position++
									}
								l1002:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleCEIL, position995)
								}
								break
							default:
								{
									position1004 := position
									depth++
									{
										position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1006
										}
										position++
										goto l1005
									l1006:
										position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
										if buffer[position] != rune('A') {
											goto l648
										}
										position++
									}
								l1005:
									{
										position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1008
										}
										position++
										goto l1007
									l1008:
										position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
										if buffer[position] != rune('B') {
											goto l648
										}
										position++
									}
								l1007:
									{
										position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1010
										}
										position++
										goto l1009
									l1010:
										position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
										if buffer[position] != rune('S') {
											goto l648
										}
										position++
									}
								l1009:
									if !rules[ruleskip]() {
										goto l648
									}
									depth--
									add(ruleABS, position1004)
								}
								break
							}
						}

					}
				l649:
					if !rules[ruleLPAREN]() {
						goto l648
					}
					if !rules[ruleexpression]() {
						goto l648
					}
					if !rules[ruleRPAREN]() {
						goto l648
					}
					goto l647
				l648:
					position, tokenIndex, depth = position647, tokenIndex647, depth647
					{
						position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
						{
							position1014 := position
							depth++
							{
								position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1016
								}
								position++
								goto l1015
							l1016:
								position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
								if buffer[position] != rune('S') {
									goto l1013
								}
								position++
							}
						l1015:
							{
								position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1018
								}
								position++
								goto l1017
							l1018:
								position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
								if buffer[position] != rune('T') {
									goto l1013
								}
								position++
							}
						l1017:
							{
								position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1020
								}
								position++
								goto l1019
							l1020:
								position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
								if buffer[position] != rune('R') {
									goto l1013
								}
								position++
							}
						l1019:
							{
								position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1022
								}
								position++
								goto l1021
							l1022:
								position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
								if buffer[position] != rune('S') {
									goto l1013
								}
								position++
							}
						l1021:
							{
								position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1024
								}
								position++
								goto l1023
							l1024:
								position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
								if buffer[position] != rune('T') {
									goto l1013
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
									goto l1013
								}
								position++
							}
						l1025:
							{
								position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1028
								}
								position++
								goto l1027
							l1028:
								position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
								if buffer[position] != rune('R') {
									goto l1013
								}
								position++
							}
						l1027:
							{
								position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1030
								}
								position++
								goto l1029
							l1030:
								position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
								if buffer[position] != rune('T') {
									goto l1013
								}
								position++
							}
						l1029:
							{
								position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1032
								}
								position++
								goto l1031
							l1032:
								position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
								if buffer[position] != rune('S') {
									goto l1013
								}
								position++
							}
						l1031:
							if !rules[ruleskip]() {
								goto l1013
							}
							depth--
							add(ruleSTRSTARTS, position1014)
						}
						goto l1012
					l1013:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						{
							position1034 := position
							depth++
							{
								position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1036
								}
								position++
								goto l1035
							l1036:
								position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
								if buffer[position] != rune('S') {
									goto l1033
								}
								position++
							}
						l1035:
							{
								position1037, tokenIndex1037, depth1037 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1038
								}
								position++
								goto l1037
							l1038:
								position, tokenIndex, depth = position1037, tokenIndex1037, depth1037
								if buffer[position] != rune('T') {
									goto l1033
								}
								position++
							}
						l1037:
							{
								position1039, tokenIndex1039, depth1039 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1040
								}
								position++
								goto l1039
							l1040:
								position, tokenIndex, depth = position1039, tokenIndex1039, depth1039
								if buffer[position] != rune('R') {
									goto l1033
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
									goto l1033
								}
								position++
							}
						l1041:
							{
								position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1044
								}
								position++
								goto l1043
							l1044:
								position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
								if buffer[position] != rune('N') {
									goto l1033
								}
								position++
							}
						l1043:
							{
								position1045, tokenIndex1045, depth1045 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1046
								}
								position++
								goto l1045
							l1046:
								position, tokenIndex, depth = position1045, tokenIndex1045, depth1045
								if buffer[position] != rune('D') {
									goto l1033
								}
								position++
							}
						l1045:
							{
								position1047, tokenIndex1047, depth1047 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1048
								}
								position++
								goto l1047
							l1048:
								position, tokenIndex, depth = position1047, tokenIndex1047, depth1047
								if buffer[position] != rune('S') {
									goto l1033
								}
								position++
							}
						l1047:
							if !rules[ruleskip]() {
								goto l1033
							}
							depth--
							add(ruleSTRENDS, position1034)
						}
						goto l1012
					l1033:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						{
							position1050 := position
							depth++
							{
								position1051, tokenIndex1051, depth1051 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1052
								}
								position++
								goto l1051
							l1052:
								position, tokenIndex, depth = position1051, tokenIndex1051, depth1051
								if buffer[position] != rune('S') {
									goto l1049
								}
								position++
							}
						l1051:
							{
								position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1054
								}
								position++
								goto l1053
							l1054:
								position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
								if buffer[position] != rune('T') {
									goto l1049
								}
								position++
							}
						l1053:
							{
								position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1056
								}
								position++
								goto l1055
							l1056:
								position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
								if buffer[position] != rune('R') {
									goto l1049
								}
								position++
							}
						l1055:
							{
								position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1058
								}
								position++
								goto l1057
							l1058:
								position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
								if buffer[position] != rune('B') {
									goto l1049
								}
								position++
							}
						l1057:
							{
								position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1060
								}
								position++
								goto l1059
							l1060:
								position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
								if buffer[position] != rune('E') {
									goto l1049
								}
								position++
							}
						l1059:
							{
								position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1062
								}
								position++
								goto l1061
							l1062:
								position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
								if buffer[position] != rune('F') {
									goto l1049
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
									goto l1049
								}
								position++
							}
						l1063:
							{
								position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1066
								}
								position++
								goto l1065
							l1066:
								position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
								if buffer[position] != rune('R') {
									goto l1049
								}
								position++
							}
						l1065:
							{
								position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1068
								}
								position++
								goto l1067
							l1068:
								position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
								if buffer[position] != rune('E') {
									goto l1049
								}
								position++
							}
						l1067:
							if !rules[ruleskip]() {
								goto l1049
							}
							depth--
							add(ruleSTRBEFORE, position1050)
						}
						goto l1012
					l1049:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						{
							position1070 := position
							depth++
							{
								position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1072
								}
								position++
								goto l1071
							l1072:
								position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
								if buffer[position] != rune('S') {
									goto l1069
								}
								position++
							}
						l1071:
							{
								position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1074
								}
								position++
								goto l1073
							l1074:
								position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
								if buffer[position] != rune('T') {
									goto l1069
								}
								position++
							}
						l1073:
							{
								position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1076
								}
								position++
								goto l1075
							l1076:
								position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
								if buffer[position] != rune('R') {
									goto l1069
								}
								position++
							}
						l1075:
							{
								position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1078
								}
								position++
								goto l1077
							l1078:
								position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
								if buffer[position] != rune('A') {
									goto l1069
								}
								position++
							}
						l1077:
							{
								position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1080
								}
								position++
								goto l1079
							l1080:
								position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
								if buffer[position] != rune('F') {
									goto l1069
								}
								position++
							}
						l1079:
							{
								position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1082
								}
								position++
								goto l1081
							l1082:
								position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
								if buffer[position] != rune('T') {
									goto l1069
								}
								position++
							}
						l1081:
							{
								position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1084
								}
								position++
								goto l1083
							l1084:
								position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
								if buffer[position] != rune('E') {
									goto l1069
								}
								position++
							}
						l1083:
							{
								position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1086
								}
								position++
								goto l1085
							l1086:
								position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
								if buffer[position] != rune('R') {
									goto l1069
								}
								position++
							}
						l1085:
							if !rules[ruleskip]() {
								goto l1069
							}
							depth--
							add(ruleSTRAFTER, position1070)
						}
						goto l1012
					l1069:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						{
							position1088 := position
							depth++
							{
								position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1090
								}
								position++
								goto l1089
							l1090:
								position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
								if buffer[position] != rune('S') {
									goto l1087
								}
								position++
							}
						l1089:
							{
								position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1092
								}
								position++
								goto l1091
							l1092:
								position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
								if buffer[position] != rune('T') {
									goto l1087
								}
								position++
							}
						l1091:
							{
								position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1094
								}
								position++
								goto l1093
							l1094:
								position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
								if buffer[position] != rune('R') {
									goto l1087
								}
								position++
							}
						l1093:
							{
								position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1096
								}
								position++
								goto l1095
							l1096:
								position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
								if buffer[position] != rune('L') {
									goto l1087
								}
								position++
							}
						l1095:
							{
								position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1098
								}
								position++
								goto l1097
							l1098:
								position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
								if buffer[position] != rune('A') {
									goto l1087
								}
								position++
							}
						l1097:
							{
								position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1100
								}
								position++
								goto l1099
							l1100:
								position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
								if buffer[position] != rune('N') {
									goto l1087
								}
								position++
							}
						l1099:
							{
								position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1102
								}
								position++
								goto l1101
							l1102:
								position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
								if buffer[position] != rune('G') {
									goto l1087
								}
								position++
							}
						l1101:
							if !rules[ruleskip]() {
								goto l1087
							}
							depth--
							add(ruleSTRLANG, position1088)
						}
						goto l1012
					l1087:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						{
							position1104 := position
							depth++
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
									goto l1103
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
									goto l1103
								}
								position++
							}
						l1107:
							{
								position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1110
								}
								position++
								goto l1109
							l1110:
								position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
								if buffer[position] != rune('R') {
									goto l1103
								}
								position++
							}
						l1109:
							{
								position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1112
								}
								position++
								goto l1111
							l1112:
								position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
								if buffer[position] != rune('D') {
									goto l1103
								}
								position++
							}
						l1111:
							{
								position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1114
								}
								position++
								goto l1113
							l1114:
								position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
								if buffer[position] != rune('T') {
									goto l1103
								}
								position++
							}
						l1113:
							if !rules[ruleskip]() {
								goto l1103
							}
							depth--
							add(ruleSTRDT, position1104)
						}
						goto l1012
					l1103:
						position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1116 := position
									depth++
									{
										position1117, tokenIndex1117, depth1117 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1118
										}
										position++
										goto l1117
									l1118:
										position, tokenIndex, depth = position1117, tokenIndex1117, depth1117
										if buffer[position] != rune('S') {
											goto l1011
										}
										position++
									}
								l1117:
									{
										position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1120
										}
										position++
										goto l1119
									l1120:
										position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
										if buffer[position] != rune('A') {
											goto l1011
										}
										position++
									}
								l1119:
									{
										position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1122
										}
										position++
										goto l1121
									l1122:
										position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
										if buffer[position] != rune('M') {
											goto l1011
										}
										position++
									}
								l1121:
									{
										position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1124
										}
										position++
										goto l1123
									l1124:
										position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
										if buffer[position] != rune('E') {
											goto l1011
										}
										position++
									}
								l1123:
									{
										position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1126
										}
										position++
										goto l1125
									l1126:
										position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
										if buffer[position] != rune('T') {
											goto l1011
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
											goto l1011
										}
										position++
									}
								l1127:
									{
										position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1130
										}
										position++
										goto l1129
									l1130:
										position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
										if buffer[position] != rune('R') {
											goto l1011
										}
										position++
									}
								l1129:
									{
										position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1132
										}
										position++
										goto l1131
									l1132:
										position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
										if buffer[position] != rune('M') {
											goto l1011
										}
										position++
									}
								l1131:
									if !rules[ruleskip]() {
										goto l1011
									}
									depth--
									add(ruleSAMETERM, position1116)
								}
								break
							case 'C', 'c':
								{
									position1133 := position
									depth++
									{
										position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1135
										}
										position++
										goto l1134
									l1135:
										position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
										if buffer[position] != rune('C') {
											goto l1011
										}
										position++
									}
								l1134:
									{
										position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1137
										}
										position++
										goto l1136
									l1137:
										position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
										if buffer[position] != rune('O') {
											goto l1011
										}
										position++
									}
								l1136:
									{
										position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1139
										}
										position++
										goto l1138
									l1139:
										position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
										if buffer[position] != rune('N') {
											goto l1011
										}
										position++
									}
								l1138:
									{
										position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1141
										}
										position++
										goto l1140
									l1141:
										position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
										if buffer[position] != rune('T') {
											goto l1011
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
											goto l1011
										}
										position++
									}
								l1142:
									{
										position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1145
										}
										position++
										goto l1144
									l1145:
										position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
										if buffer[position] != rune('I') {
											goto l1011
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
											goto l1011
										}
										position++
									}
								l1146:
									{
										position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1149
										}
										position++
										goto l1148
									l1149:
										position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
										if buffer[position] != rune('S') {
											goto l1011
										}
										position++
									}
								l1148:
									if !rules[ruleskip]() {
										goto l1011
									}
									depth--
									add(ruleCONTAINS, position1133)
								}
								break
							default:
								{
									position1150 := position
									depth++
									{
										position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1152
										}
										position++
										goto l1151
									l1152:
										position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
										if buffer[position] != rune('L') {
											goto l1011
										}
										position++
									}
								l1151:
									{
										position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1154
										}
										position++
										goto l1153
									l1154:
										position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
										if buffer[position] != rune('A') {
											goto l1011
										}
										position++
									}
								l1153:
									{
										position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1156
										}
										position++
										goto l1155
									l1156:
										position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
										if buffer[position] != rune('N') {
											goto l1011
										}
										position++
									}
								l1155:
									{
										position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1158
										}
										position++
										goto l1157
									l1158:
										position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
										if buffer[position] != rune('G') {
											goto l1011
										}
										position++
									}
								l1157:
									{
										position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1160
										}
										position++
										goto l1159
									l1160:
										position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
										if buffer[position] != rune('M') {
											goto l1011
										}
										position++
									}
								l1159:
									{
										position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1162
										}
										position++
										goto l1161
									l1162:
										position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
										if buffer[position] != rune('A') {
											goto l1011
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
											goto l1011
										}
										position++
									}
								l1163:
									{
										position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1166
										}
										position++
										goto l1165
									l1166:
										position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
										if buffer[position] != rune('C') {
											goto l1011
										}
										position++
									}
								l1165:
									{
										position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1168
										}
										position++
										goto l1167
									l1168:
										position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
										if buffer[position] != rune('H') {
											goto l1011
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
											goto l1011
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
											goto l1011
										}
										position++
									}
								l1171:
									if !rules[ruleskip]() {
										goto l1011
									}
									depth--
									add(ruleLANGMATCHES, position1150)
								}
								break
							}
						}

					}
				l1012:
					if !rules[ruleLPAREN]() {
						goto l1011
					}
					if !rules[ruleexpression]() {
						goto l1011
					}
					if !rules[ruleCOMMA]() {
						goto l1011
					}
					if !rules[ruleexpression]() {
						goto l1011
					}
					if !rules[ruleRPAREN]() {
						goto l1011
					}
					goto l647
				l1011:
					position, tokenIndex, depth = position647, tokenIndex647, depth647
					{
						position1174 := position
						depth++
						{
							position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1176
							}
							position++
							goto l1175
						l1176:
							position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
							if buffer[position] != rune('B') {
								goto l1173
							}
							position++
						}
					l1175:
						{
							position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1178
							}
							position++
							goto l1177
						l1178:
							position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
							if buffer[position] != rune('O') {
								goto l1173
							}
							position++
						}
					l1177:
						{
							position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1180
							}
							position++
							goto l1179
						l1180:
							position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
							if buffer[position] != rune('U') {
								goto l1173
							}
							position++
						}
					l1179:
						{
							position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1182
							}
							position++
							goto l1181
						l1182:
							position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
							if buffer[position] != rune('N') {
								goto l1173
							}
							position++
						}
					l1181:
						{
							position1183, tokenIndex1183, depth1183 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1184
							}
							position++
							goto l1183
						l1184:
							position, tokenIndex, depth = position1183, tokenIndex1183, depth1183
							if buffer[position] != rune('D') {
								goto l1173
							}
							position++
						}
					l1183:
						if !rules[ruleskip]() {
							goto l1173
						}
						depth--
						add(ruleBOUND, position1174)
					}
					if !rules[ruleLPAREN]() {
						goto l1173
					}
					if !rules[rulevar]() {
						goto l1173
					}
					if !rules[ruleRPAREN]() {
						goto l1173
					}
					goto l647
				l1173:
					position, tokenIndex, depth = position647, tokenIndex647, depth647
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1187 := position
								depth++
								{
									position1188, tokenIndex1188, depth1188 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1189
									}
									position++
									goto l1188
								l1189:
									position, tokenIndex, depth = position1188, tokenIndex1188, depth1188
									if buffer[position] != rune('S') {
										goto l1185
									}
									position++
								}
							l1188:
								{
									position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1191
									}
									position++
									goto l1190
								l1191:
									position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
									if buffer[position] != rune('T') {
										goto l1185
									}
									position++
								}
							l1190:
								{
									position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1193
									}
									position++
									goto l1192
								l1193:
									position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
									if buffer[position] != rune('R') {
										goto l1185
									}
									position++
								}
							l1192:
								{
									position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1195
									}
									position++
									goto l1194
								l1195:
									position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
									if buffer[position] != rune('U') {
										goto l1185
									}
									position++
								}
							l1194:
								{
									position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1197
									}
									position++
									goto l1196
								l1197:
									position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
									if buffer[position] != rune('U') {
										goto l1185
									}
									position++
								}
							l1196:
								{
									position1198, tokenIndex1198, depth1198 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1199
									}
									position++
									goto l1198
								l1199:
									position, tokenIndex, depth = position1198, tokenIndex1198, depth1198
									if buffer[position] != rune('I') {
										goto l1185
									}
									position++
								}
							l1198:
								{
									position1200, tokenIndex1200, depth1200 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1201
									}
									position++
									goto l1200
								l1201:
									position, tokenIndex, depth = position1200, tokenIndex1200, depth1200
									if buffer[position] != rune('D') {
										goto l1185
									}
									position++
								}
							l1200:
								if !rules[ruleskip]() {
									goto l1185
								}
								depth--
								add(ruleSTRUUID, position1187)
							}
							break
						case 'U', 'u':
							{
								position1202 := position
								depth++
								{
									position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1204
									}
									position++
									goto l1203
								l1204:
									position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
									if buffer[position] != rune('U') {
										goto l1185
									}
									position++
								}
							l1203:
								{
									position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1206
									}
									position++
									goto l1205
								l1206:
									position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
									if buffer[position] != rune('U') {
										goto l1185
									}
									position++
								}
							l1205:
								{
									position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1208
									}
									position++
									goto l1207
								l1208:
									position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
									if buffer[position] != rune('I') {
										goto l1185
									}
									position++
								}
							l1207:
								{
									position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1210
									}
									position++
									goto l1209
								l1210:
									position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
									if buffer[position] != rune('D') {
										goto l1185
									}
									position++
								}
							l1209:
								if !rules[ruleskip]() {
									goto l1185
								}
								depth--
								add(ruleUUID, position1202)
							}
							break
						case 'N', 'n':
							{
								position1211 := position
								depth++
								{
									position1212, tokenIndex1212, depth1212 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1213
									}
									position++
									goto l1212
								l1213:
									position, tokenIndex, depth = position1212, tokenIndex1212, depth1212
									if buffer[position] != rune('N') {
										goto l1185
									}
									position++
								}
							l1212:
								{
									position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1215
									}
									position++
									goto l1214
								l1215:
									position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
									if buffer[position] != rune('O') {
										goto l1185
									}
									position++
								}
							l1214:
								{
									position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1217
									}
									position++
									goto l1216
								l1217:
									position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
									if buffer[position] != rune('W') {
										goto l1185
									}
									position++
								}
							l1216:
								if !rules[ruleskip]() {
									goto l1185
								}
								depth--
								add(ruleNOW, position1211)
							}
							break
						default:
							{
								position1218 := position
								depth++
								{
									position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1220
									}
									position++
									goto l1219
								l1220:
									position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
									if buffer[position] != rune('R') {
										goto l1185
									}
									position++
								}
							l1219:
								{
									position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1222
									}
									position++
									goto l1221
								l1222:
									position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
									if buffer[position] != rune('A') {
										goto l1185
									}
									position++
								}
							l1221:
								{
									position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1224
									}
									position++
									goto l1223
								l1224:
									position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
									if buffer[position] != rune('N') {
										goto l1185
									}
									position++
								}
							l1223:
								{
									position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1226
									}
									position++
									goto l1225
								l1226:
									position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
									if buffer[position] != rune('D') {
										goto l1185
									}
									position++
								}
							l1225:
								if !rules[ruleskip]() {
									goto l1185
								}
								depth--
								add(ruleRAND, position1218)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1185
					}
					goto l647
				l1185:
					position, tokenIndex, depth = position647, tokenIndex647, depth647
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
								{
									position1230 := position
									depth++
									{
										position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1232
										}
										position++
										goto l1231
									l1232:
										position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
										if buffer[position] != rune('E') {
											goto l1229
										}
										position++
									}
								l1231:
									{
										position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1234
										}
										position++
										goto l1233
									l1234:
										position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
										if buffer[position] != rune('X') {
											goto l1229
										}
										position++
									}
								l1233:
									{
										position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1236
										}
										position++
										goto l1235
									l1236:
										position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
										if buffer[position] != rune('I') {
											goto l1229
										}
										position++
									}
								l1235:
									{
										position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1238
										}
										position++
										goto l1237
									l1238:
										position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
										if buffer[position] != rune('S') {
											goto l1229
										}
										position++
									}
								l1237:
									{
										position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1240
										}
										position++
										goto l1239
									l1240:
										position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
										if buffer[position] != rune('T') {
											goto l1229
										}
										position++
									}
								l1239:
									{
										position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1242
										}
										position++
										goto l1241
									l1242:
										position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
										if buffer[position] != rune('S') {
											goto l1229
										}
										position++
									}
								l1241:
									if !rules[ruleskip]() {
										goto l1229
									}
									depth--
									add(ruleEXISTS, position1230)
								}
								goto l1228
							l1229:
								position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
								{
									position1243 := position
									depth++
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
											goto l645
										}
										position++
									}
								l1244:
									{
										position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1247
										}
										position++
										goto l1246
									l1247:
										position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
										if buffer[position] != rune('O') {
											goto l645
										}
										position++
									}
								l1246:
									{
										position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1249
										}
										position++
										goto l1248
									l1249:
										position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
										if buffer[position] != rune('T') {
											goto l645
										}
										position++
									}
								l1248:
									if buffer[position] != rune(' ') {
										goto l645
									}
									position++
									{
										position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1251
										}
										position++
										goto l1250
									l1251:
										position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
										if buffer[position] != rune('E') {
											goto l645
										}
										position++
									}
								l1250:
									{
										position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1253
										}
										position++
										goto l1252
									l1253:
										position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
										if buffer[position] != rune('X') {
											goto l645
										}
										position++
									}
								l1252:
									{
										position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1255
										}
										position++
										goto l1254
									l1255:
										position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
										if buffer[position] != rune('I') {
											goto l645
										}
										position++
									}
								l1254:
									{
										position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1257
										}
										position++
										goto l1256
									l1257:
										position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
										if buffer[position] != rune('S') {
											goto l645
										}
										position++
									}
								l1256:
									{
										position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1259
										}
										position++
										goto l1258
									l1259:
										position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
										if buffer[position] != rune('T') {
											goto l645
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
											goto l645
										}
										position++
									}
								l1260:
									if !rules[ruleskip]() {
										goto l645
									}
									depth--
									add(ruleNOTEXIST, position1243)
								}
							}
						l1228:
							if !rules[rulegroupGraphPattern]() {
								goto l645
							}
							break
						case 'I', 'i':
							{
								position1262 := position
								depth++
								{
									position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1264
									}
									position++
									goto l1263
								l1264:
									position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
									if buffer[position] != rune('I') {
										goto l645
									}
									position++
								}
							l1263:
								{
									position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1266
									}
									position++
									goto l1265
								l1266:
									position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
									if buffer[position] != rune('F') {
										goto l645
									}
									position++
								}
							l1265:
								if !rules[ruleskip]() {
									goto l645
								}
								depth--
								add(ruleIF, position1262)
							}
							if !rules[ruleLPAREN]() {
								goto l645
							}
							if !rules[ruleexpression]() {
								goto l645
							}
							if !rules[ruleCOMMA]() {
								goto l645
							}
							if !rules[ruleexpression]() {
								goto l645
							}
							if !rules[ruleCOMMA]() {
								goto l645
							}
							if !rules[ruleexpression]() {
								goto l645
							}
							if !rules[ruleRPAREN]() {
								goto l645
							}
							break
						case 'C', 'c':
							{
								position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
								{
									position1269 := position
									depth++
									{
										position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1271
										}
										position++
										goto l1270
									l1271:
										position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
										if buffer[position] != rune('C') {
											goto l1268
										}
										position++
									}
								l1270:
									{
										position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1273
										}
										position++
										goto l1272
									l1273:
										position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
										if buffer[position] != rune('O') {
											goto l1268
										}
										position++
									}
								l1272:
									{
										position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1275
										}
										position++
										goto l1274
									l1275:
										position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
										if buffer[position] != rune('N') {
											goto l1268
										}
										position++
									}
								l1274:
									{
										position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1277
										}
										position++
										goto l1276
									l1277:
										position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
										if buffer[position] != rune('C') {
											goto l1268
										}
										position++
									}
								l1276:
									{
										position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1279
										}
										position++
										goto l1278
									l1279:
										position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
										if buffer[position] != rune('A') {
											goto l1268
										}
										position++
									}
								l1278:
									{
										position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1281
										}
										position++
										goto l1280
									l1281:
										position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
										if buffer[position] != rune('T') {
											goto l1268
										}
										position++
									}
								l1280:
									if !rules[ruleskip]() {
										goto l1268
									}
									depth--
									add(ruleCONCAT, position1269)
								}
								goto l1267
							l1268:
								position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
								{
									position1282 := position
									depth++
									{
										position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1284
										}
										position++
										goto l1283
									l1284:
										position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
										if buffer[position] != rune('C') {
											goto l645
										}
										position++
									}
								l1283:
									{
										position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1286
										}
										position++
										goto l1285
									l1286:
										position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
										if buffer[position] != rune('O') {
											goto l645
										}
										position++
									}
								l1285:
									{
										position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1288
										}
										position++
										goto l1287
									l1288:
										position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
										if buffer[position] != rune('A') {
											goto l645
										}
										position++
									}
								l1287:
									{
										position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1290
										}
										position++
										goto l1289
									l1290:
										position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
										if buffer[position] != rune('L') {
											goto l645
										}
										position++
									}
								l1289:
									{
										position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1292
										}
										position++
										goto l1291
									l1292:
										position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
										if buffer[position] != rune('E') {
											goto l645
										}
										position++
									}
								l1291:
									{
										position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1294
										}
										position++
										goto l1293
									l1294:
										position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
										if buffer[position] != rune('S') {
											goto l645
										}
										position++
									}
								l1293:
									{
										position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1296
										}
										position++
										goto l1295
									l1296:
										position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
										if buffer[position] != rune('C') {
											goto l645
										}
										position++
									}
								l1295:
									{
										position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1298
										}
										position++
										goto l1297
									l1298:
										position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
										if buffer[position] != rune('E') {
											goto l645
										}
										position++
									}
								l1297:
									if !rules[ruleskip]() {
										goto l645
									}
									depth--
									add(ruleCOALESCE, position1282)
								}
							}
						l1267:
							if !rules[ruleargList]() {
								goto l645
							}
							break
						case 'B', 'b':
							{
								position1299 := position
								depth++
								{
									position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1301
									}
									position++
									goto l1300
								l1301:
									position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
									if buffer[position] != rune('B') {
										goto l645
									}
									position++
								}
							l1300:
								{
									position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1303
									}
									position++
									goto l1302
								l1303:
									position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
									if buffer[position] != rune('N') {
										goto l645
									}
									position++
								}
							l1302:
								{
									position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1305
									}
									position++
									goto l1304
								l1305:
									position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
									if buffer[position] != rune('O') {
										goto l645
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
										goto l645
									}
									position++
								}
							l1306:
								{
									position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1309
									}
									position++
									goto l1308
								l1309:
									position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
									if buffer[position] != rune('E') {
										goto l645
									}
									position++
								}
							l1308:
								if !rules[ruleskip]() {
									goto l645
								}
								depth--
								add(ruleBNODE, position1299)
							}
							{
								position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1311
								}
								if !rules[ruleexpression]() {
									goto l1311
								}
								if !rules[ruleRPAREN]() {
									goto l1311
								}
								goto l1310
							l1311:
								position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
								if !rules[rulenil]() {
									goto l645
								}
							}
						l1310:
							break
						default:
							{
								position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
								{
									position1314 := position
									depth++
									{
										position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1316
										}
										position++
										goto l1315
									l1316:
										position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
										if buffer[position] != rune('S') {
											goto l1313
										}
										position++
									}
								l1315:
									{
										position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1318
										}
										position++
										goto l1317
									l1318:
										position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
										if buffer[position] != rune('U') {
											goto l1313
										}
										position++
									}
								l1317:
									{
										position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1320
										}
										position++
										goto l1319
									l1320:
										position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
										if buffer[position] != rune('B') {
											goto l1313
										}
										position++
									}
								l1319:
									{
										position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1322
										}
										position++
										goto l1321
									l1322:
										position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
										if buffer[position] != rune('S') {
											goto l1313
										}
										position++
									}
								l1321:
									{
										position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1324
										}
										position++
										goto l1323
									l1324:
										position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
										if buffer[position] != rune('T') {
											goto l1313
										}
										position++
									}
								l1323:
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
											goto l1313
										}
										position++
									}
								l1325:
									if !rules[ruleskip]() {
										goto l1313
									}
									depth--
									add(ruleSUBSTR, position1314)
								}
								goto l1312
							l1313:
								position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
								{
									position1328 := position
									depth++
									{
										position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1330
										}
										position++
										goto l1329
									l1330:
										position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
										if buffer[position] != rune('R') {
											goto l1327
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
											goto l1327
										}
										position++
									}
								l1331:
									{
										position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1334
										}
										position++
										goto l1333
									l1334:
										position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
										if buffer[position] != rune('P') {
											goto l1327
										}
										position++
									}
								l1333:
									{
										position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1336
										}
										position++
										goto l1335
									l1336:
										position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
										if buffer[position] != rune('L') {
											goto l1327
										}
										position++
									}
								l1335:
									{
										position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1338
										}
										position++
										goto l1337
									l1338:
										position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
										if buffer[position] != rune('A') {
											goto l1327
										}
										position++
									}
								l1337:
									{
										position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1340
										}
										position++
										goto l1339
									l1340:
										position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
										if buffer[position] != rune('C') {
											goto l1327
										}
										position++
									}
								l1339:
									{
										position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1342
										}
										position++
										goto l1341
									l1342:
										position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
										if buffer[position] != rune('E') {
											goto l1327
										}
										position++
									}
								l1341:
									if !rules[ruleskip]() {
										goto l1327
									}
									depth--
									add(ruleREPLACE, position1328)
								}
								goto l1312
							l1327:
								position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
								{
									position1343 := position
									depth++
									{
										position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1345
										}
										position++
										goto l1344
									l1345:
										position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
										if buffer[position] != rune('R') {
											goto l645
										}
										position++
									}
								l1344:
									{
										position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1347
										}
										position++
										goto l1346
									l1347:
										position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
										if buffer[position] != rune('E') {
											goto l645
										}
										position++
									}
								l1346:
									{
										position1348, tokenIndex1348, depth1348 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1349
										}
										position++
										goto l1348
									l1349:
										position, tokenIndex, depth = position1348, tokenIndex1348, depth1348
										if buffer[position] != rune('G') {
											goto l645
										}
										position++
									}
								l1348:
									{
										position1350, tokenIndex1350, depth1350 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1351
										}
										position++
										goto l1350
									l1351:
										position, tokenIndex, depth = position1350, tokenIndex1350, depth1350
										if buffer[position] != rune('E') {
											goto l645
										}
										position++
									}
								l1350:
									{
										position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1353
										}
										position++
										goto l1352
									l1353:
										position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
										if buffer[position] != rune('X') {
											goto l645
										}
										position++
									}
								l1352:
									if !rules[ruleskip]() {
										goto l645
									}
									depth--
									add(ruleREGEX, position1343)
								}
							}
						l1312:
							if !rules[ruleLPAREN]() {
								goto l645
							}
							if !rules[ruleexpression]() {
								goto l645
							}
							if !rules[ruleCOMMA]() {
								goto l645
							}
							if !rules[ruleexpression]() {
								goto l645
							}
							{
								position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1354
								}
								if !rules[ruleexpression]() {
									goto l1354
								}
								goto l1355
							l1354:
								position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
							}
						l1355:
							if !rules[ruleRPAREN]() {
								goto l645
							}
							break
						}
					}

				}
			l647:
				depth--
				add(rulebuiltinCall, position646)
			}
			return true
		l645:
			position, tokenIndex, depth = position645, tokenIndex645, depth645
			return false
		},
		/* 64 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
			{
				position1357 := position
				depth++
				{
					position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1359
					}
					position++
					goto l1358
				l1359:
					position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
					if buffer[position] != rune('$') {
						goto l1356
					}
					position++
				}
			l1358:
				{
					position1360 := position
					depth++
					{
						position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
						{
							position1365 := position
							depth++
							{
								position1366, tokenIndex1366, depth1366 := position, tokenIndex, depth
								{
									position1368 := position
									depth++
									{
										position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1370
										}
										position++
										goto l1369
									l1370:
										position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1367
										}
										position++
									}
								l1369:
									depth--
									add(rulePN_CHARS_BASE, position1368)
								}
								goto l1366
							l1367:
								position, tokenIndex, depth = position1366, tokenIndex1366, depth1366
								if buffer[position] != rune('_') {
									goto l1364
								}
								position++
							}
						l1366:
							depth--
							add(rulePN_CHARS_U, position1365)
						}
						goto l1363
					l1364:
						position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1356
						}
						position++
					}
				l1363:
				l1361:
					{
						position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
						{
							position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
							{
								position1373 := position
								depth++
								{
									position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
									{
										position1376 := position
										depth++
										{
											position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1378
											}
											position++
											goto l1377
										l1378:
											position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1375
											}
											position++
										}
									l1377:
										depth--
										add(rulePN_CHARS_BASE, position1376)
									}
									goto l1374
								l1375:
									position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
									if buffer[position] != rune('_') {
										goto l1372
									}
									position++
								}
							l1374:
								depth--
								add(rulePN_CHARS_U, position1373)
							}
							goto l1371
						l1372:
							position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1362
							}
							position++
						}
					l1371:
						goto l1361
					l1362:
						position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
					}
					depth--
					add(ruleVARNAME, position1360)
				}
				if !rules[ruleskip]() {
					goto l1356
				}
				depth--
				add(rulevar, position1357)
			}
			return true
		l1356:
			position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
			return false
		},
		/* 65 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
			{
				position1380 := position
				depth++
				{
					position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1382
					}
					goto l1381
				l1382:
					position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
					{
						position1383 := position
						depth++
					l1384:
						{
							position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
							{
								position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
								{
									position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1388
									}
									position++
									goto l1387
								l1388:
									position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
									if buffer[position] != rune(' ') {
										goto l1386
									}
									position++
								}
							l1387:
								goto l1385
							l1386:
								position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
							}
							if !matchDot() {
								goto l1385
							}
							goto l1384
						l1385:
							position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
						}
						if buffer[position] != rune(':') {
							goto l1379
						}
						position++
					l1389:
						{
							position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
							{
								position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1392
								}
								position++
								goto l1391
							l1392:
								position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1393
								}
								position++
								goto l1391
							l1393:
								position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1394
								}
								position++
								goto l1391
							l1394:
								position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1390
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1390
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1390
										}
										position++
										break
									}
								}

							}
						l1391:
							goto l1389
						l1390:
							position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
						}
						if !rules[ruleskip]() {
							goto l1379
						}
						depth--
						add(ruleprefixedName, position1383)
					}
				}
			l1381:
				depth--
				add(ruleiriref, position1380)
			}
			return true
		l1379:
			position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
			return false
		},
		/* 66 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
			{
				position1397 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1396
				}
				position++
			l1398:
				{
					position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
					{
						position1400, tokenIndex1400, depth1400 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1400
						}
						position++
						goto l1399
					l1400:
						position, tokenIndex, depth = position1400, tokenIndex1400, depth1400
					}
					if !matchDot() {
						goto l1399
					}
					goto l1398
				l1399:
					position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
				}
				if buffer[position] != rune('>') {
					goto l1396
				}
				position++
				if !rules[ruleskip]() {
					goto l1396
				}
				depth--
				add(ruleiri, position1397)
			}
			return true
		l1396:
			position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
			return false
		},
		/* 67 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 68 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
			{
				position1403 := position
				depth++
				if !rules[rulestring]() {
					goto l1402
				}
				{
					position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
					{
						position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1407
						}
						position++
						{
							position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1411
							}
							position++
							goto l1410
						l1411:
							position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1407
							}
							position++
						}
					l1410:
					l1408:
						{
							position1409, tokenIndex1409, depth1409 := position, tokenIndex, depth
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
									goto l1409
								}
								position++
							}
						l1412:
							goto l1408
						l1409:
							position, tokenIndex, depth = position1409, tokenIndex1409, depth1409
						}
					l1414:
						{
							position1415, tokenIndex1415, depth1415 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1415
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1415
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1415
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1415
									}
									position++
									break
								}
							}

						l1416:
							{
								position1417, tokenIndex1417, depth1417 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1417
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1417
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1417
										}
										position++
										break
									}
								}

								goto l1416
							l1417:
								position, tokenIndex, depth = position1417, tokenIndex1417, depth1417
							}
							goto l1414
						l1415:
							position, tokenIndex, depth = position1415, tokenIndex1415, depth1415
						}
						goto l1406
					l1407:
						position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
						if buffer[position] != rune('^') {
							goto l1404
						}
						position++
						if buffer[position] != rune('^') {
							goto l1404
						}
						position++
						if !rules[ruleiriref]() {
							goto l1404
						}
					}
				l1406:
					goto l1405
				l1404:
					position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
				}
			l1405:
				if !rules[ruleskip]() {
					goto l1402
				}
				depth--
				add(ruleliteral, position1403)
			}
			return true
		l1402:
			position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
			return false
		},
		/* 69 string <- <('"' (!'"' .)* '"')> */
		func() bool {
			position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
			{
				position1421 := position
				depth++
				if buffer[position] != rune('"') {
					goto l1420
				}
				position++
			l1422:
				{
					position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
					{
						position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
						if buffer[position] != rune('"') {
							goto l1424
						}
						position++
						goto l1423
					l1424:
						position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
					}
					if !matchDot() {
						goto l1423
					}
					goto l1422
				l1423:
					position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
				}
				if buffer[position] != rune('"') {
					goto l1420
				}
				position++
				depth--
				add(rulestring, position1421)
			}
			return true
		l1420:
			position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
			return false
		},
		/* 70 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
			{
				position1426 := position
				depth++
				{
					position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
					{
						position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1430
						}
						position++
						goto l1429
					l1430:
						position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
						if buffer[position] != rune('-') {
							goto l1427
						}
						position++
					}
				l1429:
					goto l1428
				l1427:
					position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
				}
			l1428:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1425
				}
				position++
			l1431:
				{
					position1432, tokenIndex1432, depth1432 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1432
					}
					position++
					goto l1431
				l1432:
					position, tokenIndex, depth = position1432, tokenIndex1432, depth1432
				}
				{
					position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1433
					}
					position++
				l1435:
					{
						position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1436
						}
						position++
						goto l1435
					l1436:
						position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
					}
					goto l1434
				l1433:
					position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
				}
			l1434:
				if !rules[ruleskip]() {
					goto l1425
				}
				depth--
				add(rulenumericLiteral, position1426)
			}
			return true
		l1425:
			position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
			return false
		},
		/* 71 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 72 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
			{
				position1439 := position
				depth++
				{
					position1440, tokenIndex1440, depth1440 := position, tokenIndex, depth
					{
						position1442 := position
						depth++
						{
							position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1444
							}
							position++
							goto l1443
						l1444:
							position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
							if buffer[position] != rune('T') {
								goto l1441
							}
							position++
						}
					l1443:
						{
							position1445, tokenIndex1445, depth1445 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1446
							}
							position++
							goto l1445
						l1446:
							position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
							if buffer[position] != rune('R') {
								goto l1441
							}
							position++
						}
					l1445:
						{
							position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1448
							}
							position++
							goto l1447
						l1448:
							position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
							if buffer[position] != rune('U') {
								goto l1441
							}
							position++
						}
					l1447:
						{
							position1449, tokenIndex1449, depth1449 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1450
							}
							position++
							goto l1449
						l1450:
							position, tokenIndex, depth = position1449, tokenIndex1449, depth1449
							if buffer[position] != rune('E') {
								goto l1441
							}
							position++
						}
					l1449:
						if !rules[ruleskip]() {
							goto l1441
						}
						depth--
						add(ruleTRUE, position1442)
					}
					goto l1440
				l1441:
					position, tokenIndex, depth = position1440, tokenIndex1440, depth1440
					{
						position1451 := position
						depth++
						{
							position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1453
							}
							position++
							goto l1452
						l1453:
							position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
							if buffer[position] != rune('F') {
								goto l1438
							}
							position++
						}
					l1452:
						{
							position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1455
							}
							position++
							goto l1454
						l1455:
							position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
							if buffer[position] != rune('A') {
								goto l1438
							}
							position++
						}
					l1454:
						{
							position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1457
							}
							position++
							goto l1456
						l1457:
							position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
							if buffer[position] != rune('L') {
								goto l1438
							}
							position++
						}
					l1456:
						{
							position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1459
							}
							position++
							goto l1458
						l1459:
							position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
							if buffer[position] != rune('S') {
								goto l1438
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
								goto l1438
							}
							position++
						}
					l1460:
						if !rules[ruleskip]() {
							goto l1438
						}
						depth--
						add(ruleFALSE, position1451)
					}
				}
			l1440:
				depth--
				add(rulebooleanLiteral, position1439)
			}
			return true
		l1438:
			position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
			return false
		},
		/* 73 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 74 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 75 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 76 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
			{
				position1466 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1465
				}
				position++
			l1467:
				{
					position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1468
					}
					goto l1467
				l1468:
					position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
				}
				if buffer[position] != rune(')') {
					goto l1465
				}
				position++
				if !rules[ruleskip]() {
					goto l1465
				}
				depth--
				add(rulenil, position1466)
			}
			return true
		l1465:
			position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
			return false
		},
		/* 77 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 78 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 79 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 80 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 81 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 82 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 83 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 84 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 85 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 86 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1478, tokenIndex1478, depth1478 := position, tokenIndex, depth
			{
				position1479 := position
				depth++
				{
					position1480, tokenIndex1480, depth1480 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1481
					}
					position++
					goto l1480
				l1481:
					position, tokenIndex, depth = position1480, tokenIndex1480, depth1480
					if buffer[position] != rune('D') {
						goto l1478
					}
					position++
				}
			l1480:
				{
					position1482, tokenIndex1482, depth1482 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1483
					}
					position++
					goto l1482
				l1483:
					position, tokenIndex, depth = position1482, tokenIndex1482, depth1482
					if buffer[position] != rune('I') {
						goto l1478
					}
					position++
				}
			l1482:
				{
					position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1485
					}
					position++
					goto l1484
				l1485:
					position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
					if buffer[position] != rune('S') {
						goto l1478
					}
					position++
				}
			l1484:
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
						goto l1478
					}
					position++
				}
			l1486:
				{
					position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1489
					}
					position++
					goto l1488
				l1489:
					position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
					if buffer[position] != rune('I') {
						goto l1478
					}
					position++
				}
			l1488:
				{
					position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1491
					}
					position++
					goto l1490
				l1491:
					position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
					if buffer[position] != rune('N') {
						goto l1478
					}
					position++
				}
			l1490:
				{
					position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1493
					}
					position++
					goto l1492
				l1493:
					position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
					if buffer[position] != rune('C') {
						goto l1478
					}
					position++
				}
			l1492:
				{
					position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1495
					}
					position++
					goto l1494
				l1495:
					position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
					if buffer[position] != rune('T') {
						goto l1478
					}
					position++
				}
			l1494:
				if !rules[ruleskip]() {
					goto l1478
				}
				depth--
				add(ruleDISTINCT, position1479)
			}
			return true
		l1478:
			position, tokenIndex, depth = position1478, tokenIndex1478, depth1478
			return false
		},
		/* 87 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 88 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 89 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 90 LBRACE <- <('{' skip)> */
		func() bool {
			position1499, tokenIndex1499, depth1499 := position, tokenIndex, depth
			{
				position1500 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1499
				}
				position++
				if !rules[ruleskip]() {
					goto l1499
				}
				depth--
				add(ruleLBRACE, position1500)
			}
			return true
		l1499:
			position, tokenIndex, depth = position1499, tokenIndex1499, depth1499
			return false
		},
		/* 91 RBRACE <- <('}' skip)> */
		func() bool {
			position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
			{
				position1502 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1501
				}
				position++
				if !rules[ruleskip]() {
					goto l1501
				}
				depth--
				add(ruleRBRACE, position1502)
			}
			return true
		l1501:
			position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
			return false
		},
		/* 92 LBRACK <- <('[' skip)> */
		nil,
		/* 93 RBRACK <- <(']' skip)> */
		nil,
		/* 94 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
			{
				position1506 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1505
				}
				position++
				if !rules[ruleskip]() {
					goto l1505
				}
				depth--
				add(ruleSEMICOLON, position1506)
			}
			return true
		l1505:
			position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
			return false
		},
		/* 95 COMMA <- <(',' skip)> */
		func() bool {
			position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
			{
				position1508 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1507
				}
				position++
				if !rules[ruleskip]() {
					goto l1507
				}
				depth--
				add(ruleCOMMA, position1508)
			}
			return true
		l1507:
			position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
			return false
		},
		/* 96 DOT <- <('.' skip)> */
		func() bool {
			position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
			{
				position1510 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1509
				}
				position++
				if !rules[ruleskip]() {
					goto l1509
				}
				depth--
				add(ruleDOT, position1510)
			}
			return true
		l1509:
			position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
			return false
		},
		/* 97 COLON <- <(':' skip)> */
		nil,
		/* 98 PIPE <- <('|' skip)> */
		func() bool {
			position1512, tokenIndex1512, depth1512 := position, tokenIndex, depth
			{
				position1513 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1512
				}
				position++
				if !rules[ruleskip]() {
					goto l1512
				}
				depth--
				add(rulePIPE, position1513)
			}
			return true
		l1512:
			position, tokenIndex, depth = position1512, tokenIndex1512, depth1512
			return false
		},
		/* 99 SLASH <- <('/' skip)> */
		func() bool {
			position1514, tokenIndex1514, depth1514 := position, tokenIndex, depth
			{
				position1515 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1514
				}
				position++
				if !rules[ruleskip]() {
					goto l1514
				}
				depth--
				add(ruleSLASH, position1515)
			}
			return true
		l1514:
			position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
			return false
		},
		/* 100 INVERSE <- <('^' skip)> */
		func() bool {
			position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
			{
				position1517 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1516
				}
				position++
				if !rules[ruleskip]() {
					goto l1516
				}
				depth--
				add(ruleINVERSE, position1517)
			}
			return true
		l1516:
			position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
			return false
		},
		/* 101 LPAREN <- <('(' skip)> */
		func() bool {
			position1518, tokenIndex1518, depth1518 := position, tokenIndex, depth
			{
				position1519 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1518
				}
				position++
				if !rules[ruleskip]() {
					goto l1518
				}
				depth--
				add(ruleLPAREN, position1519)
			}
			return true
		l1518:
			position, tokenIndex, depth = position1518, tokenIndex1518, depth1518
			return false
		},
		/* 102 RPAREN <- <(')' skip)> */
		func() bool {
			position1520, tokenIndex1520, depth1520 := position, tokenIndex, depth
			{
				position1521 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1520
				}
				position++
				if !rules[ruleskip]() {
					goto l1520
				}
				depth--
				add(ruleRPAREN, position1521)
			}
			return true
		l1520:
			position, tokenIndex, depth = position1520, tokenIndex1520, depth1520
			return false
		},
		/* 103 ISA <- <('a' skip)> */
		func() bool {
			position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
			{
				position1523 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1522
				}
				position++
				if !rules[ruleskip]() {
					goto l1522
				}
				depth--
				add(ruleISA, position1523)
			}
			return true
		l1522:
			position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
			return false
		},
		/* 104 NOT <- <('!' skip)> */
		func() bool {
			position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
			{
				position1525 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1524
				}
				position++
				if !rules[ruleskip]() {
					goto l1524
				}
				depth--
				add(ruleNOT, position1525)
			}
			return true
		l1524:
			position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
			return false
		},
		/* 105 STAR <- <('*' skip)> */
		func() bool {
			position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
			{
				position1527 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1526
				}
				position++
				if !rules[ruleskip]() {
					goto l1526
				}
				depth--
				add(ruleSTAR, position1527)
			}
			return true
		l1526:
			position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
			return false
		},
		/* 106 PLUS <- <('+' skip)> */
		func() bool {
			position1528, tokenIndex1528, depth1528 := position, tokenIndex, depth
			{
				position1529 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1528
				}
				position++
				if !rules[ruleskip]() {
					goto l1528
				}
				depth--
				add(rulePLUS, position1529)
			}
			return true
		l1528:
			position, tokenIndex, depth = position1528, tokenIndex1528, depth1528
			return false
		},
		/* 107 MINUS <- <('-' skip)> */
		func() bool {
			position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
			{
				position1531 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1530
				}
				position++
				if !rules[ruleskip]() {
					goto l1530
				}
				depth--
				add(ruleMINUS, position1531)
			}
			return true
		l1530:
			position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
			return false
		},
		/* 108 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 109 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 110 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 111 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 112 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1536, tokenIndex1536, depth1536 := position, tokenIndex, depth
			{
				position1537 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1536
				}
				position++
			l1538:
				{
					position1539, tokenIndex1539, depth1539 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1539
					}
					position++
					goto l1538
				l1539:
					position, tokenIndex, depth = position1539, tokenIndex1539, depth1539
				}
				if !rules[ruleskip]() {
					goto l1536
				}
				depth--
				add(ruleINTEGER, position1537)
			}
			return true
		l1536:
			position, tokenIndex, depth = position1536, tokenIndex1536, depth1536
			return false
		},
		/* 113 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 114 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 115 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 116 OR <- <('|' '|' skip)> */
		nil,
		/* 117 AND <- <('&' '&' skip)> */
		nil,
		/* 118 EQ <- <('=' skip)> */
		func() bool {
			position1545, tokenIndex1545, depth1545 := position, tokenIndex, depth
			{
				position1546 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1545
				}
				position++
				if !rules[ruleskip]() {
					goto l1545
				}
				depth--
				add(ruleEQ, position1546)
			}
			return true
		l1545:
			position, tokenIndex, depth = position1545, tokenIndex1545, depth1545
			return false
		},
		/* 119 NE <- <('!' '=' skip)> */
		nil,
		/* 120 GT <- <('>' skip)> */
		nil,
		/* 121 LT <- <('<' skip)> */
		nil,
		/* 122 LE <- <('<' '=' skip)> */
		nil,
		/* 123 GE <- <('>' '=' skip)> */
		nil,
		/* 124 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 125 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 126 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1554, tokenIndex1554, depth1554 := position, tokenIndex, depth
			{
				position1555 := position
				depth++
				{
					position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1557
					}
					position++
					goto l1556
				l1557:
					position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
					if buffer[position] != rune('A') {
						goto l1554
					}
					position++
				}
			l1556:
				{
					position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1559
					}
					position++
					goto l1558
				l1559:
					position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
					if buffer[position] != rune('S') {
						goto l1554
					}
					position++
				}
			l1558:
				if !rules[ruleskip]() {
					goto l1554
				}
				depth--
				add(ruleAS, position1555)
			}
			return true
		l1554:
			position, tokenIndex, depth = position1554, tokenIndex1554, depth1554
			return false
		},
		/* 127 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 128 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 129 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 130 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 131 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 132 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 133 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 134 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 135 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 136 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 137 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 138 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 139 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 140 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 141 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 142 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 143 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 144 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 145 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 146 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 147 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 148 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 149 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 150 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 151 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 152 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 153 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 154 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 155 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 156 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 157 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 158 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 159 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 160 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 161 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 162 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 163 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 164 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 165 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 166 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 167 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 168 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 169 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 170 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 171 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 172 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 173 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 174 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 175 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 176 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 177 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 178 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 179 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 180 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 181 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 182 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 183 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 184 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 185 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 186 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 187 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 188 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 189 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 190 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 191 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1625 := position
				depth++
			l1626:
				{
					position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
					{
						position1628, tokenIndex1628, depth1628 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1629
						}
						goto l1628
					l1629:
						position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
						{
							position1630 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1627
							}
							position++
						l1631:
							{
								position1632, tokenIndex1632, depth1632 := position, tokenIndex, depth
								{
									position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1633
									}
									goto l1632
								l1633:
									position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
								}
								if !matchDot() {
									goto l1632
								}
								goto l1631
							l1632:
								position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
							}
							if !rules[ruleendOfLine]() {
								goto l1627
							}
							depth--
							add(rulecomment, position1630)
						}
					}
				l1628:
					goto l1626
				l1627:
					position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
				}
				depth--
				add(ruleskip, position1625)
			}
			return true
		},
		/* 192 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1634, tokenIndex1634, depth1634 := position, tokenIndex, depth
			{
				position1635 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1634
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1634
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1634
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1634
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1634
						}
						break
					}
				}

				depth--
				add(rulews, position1635)
			}
			return true
		l1634:
			position, tokenIndex, depth = position1634, tokenIndex1634, depth1634
			return false
		},
		/* 193 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 194 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1638, tokenIndex1638, depth1638 := position, tokenIndex, depth
			{
				position1639 := position
				depth++
				{
					position1640, tokenIndex1640, depth1640 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1641
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1641
					}
					position++
					goto l1640
				l1641:
					position, tokenIndex, depth = position1640, tokenIndex1640, depth1640
					if buffer[position] != rune('\n') {
						goto l1642
					}
					position++
					goto l1640
				l1642:
					position, tokenIndex, depth = position1640, tokenIndex1640, depth1640
					if buffer[position] != rune('\r') {
						goto l1638
					}
					position++
				}
			l1640:
				depth--
				add(ruleendOfLine, position1639)
			}
			return true
		l1638:
			position, tokenIndex, depth = position1638, tokenIndex1638, depth1638
			return false
		},
	}
	p.rules = rules
}
