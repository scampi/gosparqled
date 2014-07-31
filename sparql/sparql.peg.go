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
	rulegraphGraphPattern
	ruleminusGraphPattern
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
	rulepathMod
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
	ruleQUESTION
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
	ruleGRAPH
	ruleMINUSSETOPER
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
	"graphGraphPattern",
	"minusGraphPattern",
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
	"pathMod",
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
	"QUESTION",
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
	"GRAPH",
	"MINUSSETOPER",
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
	rules  [210]func() bool
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
							l47:
								{
									position48, tokenIndex48, depth48 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l48
									}
									goto l47
								l48:
									position, tokenIndex, depth = position48, tokenIndex48, depth48
								}
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
							l71:
								{
									position72, tokenIndex72, depth72 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l72
									}
									goto l71
								l72:
									position, tokenIndex, depth = position72, tokenIndex72, depth72
								}
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
							l98:
								{
									position99, tokenIndex99, depth99 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l99
									}
									goto l98
								l99:
									position, tokenIndex, depth = position99, tokenIndex99, depth99
								}
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
							l101:
								{
									position102, tokenIndex102, depth102 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l102
									}
									goto l101
								l102:
									position, tokenIndex, depth = position102, tokenIndex102, depth102
								}
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
		/* 5 selectQuery <- <(select datasetClause* whereClause solutionModifier)> */
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
		/* 8 constructQuery <- <(construct datasetClause* whereClause solutionModifier)> */
		nil,
		/* 9 construct <- <(CONSTRUCT LBRACE triplesBlock? RBRACE)> */
		nil,
		/* 10 describeQuery <- <(describe datasetClause* whereClause? solutionModifier)> */
		nil,
		/* 11 describe <- <(DESCRIBE (STAR / var / iriref))> */
		nil,
		/* 12 askQuery <- <(ASK datasetClause* whereClause)> */
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
							switch buffer[position] {
							case 'M', 'm':
								{
									position231 := position
									depth++
									{
										position232 := position
										depth++
										{
											position233, tokenIndex233, depth233 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l234
											}
											position++
											goto l233
										l234:
											position, tokenIndex, depth = position233, tokenIndex233, depth233
											if buffer[position] != rune('M') {
												goto l227
											}
											position++
										}
									l233:
										{
											position235, tokenIndex235, depth235 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l236
											}
											position++
											goto l235
										l236:
											position, tokenIndex, depth = position235, tokenIndex235, depth235
											if buffer[position] != rune('I') {
												goto l227
											}
											position++
										}
									l235:
										{
											position237, tokenIndex237, depth237 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l238
											}
											position++
											goto l237
										l238:
											position, tokenIndex, depth = position237, tokenIndex237, depth237
											if buffer[position] != rune('N') {
												goto l227
											}
											position++
										}
									l237:
										{
											position239, tokenIndex239, depth239 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l240
											}
											position++
											goto l239
										l240:
											position, tokenIndex, depth = position239, tokenIndex239, depth239
											if buffer[position] != rune('U') {
												goto l227
											}
											position++
										}
									l239:
										{
											position241, tokenIndex241, depth241 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l242
											}
											position++
											goto l241
										l242:
											position, tokenIndex, depth = position241, tokenIndex241, depth241
											if buffer[position] != rune('S') {
												goto l227
											}
											position++
										}
									l241:
										if !rules[ruleskip]() {
											goto l227
										}
										depth--
										add(ruleMINUSSETOPER, position232)
									}
									if !rules[rulegroupGraphPattern]() {
										goto l227
									}
									depth--
									add(ruleminusGraphPattern, position231)
								}
								break
							case 'G', 'g':
								{
									position243 := position
									depth++
									{
										position244 := position
										depth++
										{
											position245, tokenIndex245, depth245 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l246
											}
											position++
											goto l245
										l246:
											position, tokenIndex, depth = position245, tokenIndex245, depth245
											if buffer[position] != rune('G') {
												goto l227
											}
											position++
										}
									l245:
										{
											position247, tokenIndex247, depth247 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l248
											}
											position++
											goto l247
										l248:
											position, tokenIndex, depth = position247, tokenIndex247, depth247
											if buffer[position] != rune('R') {
												goto l227
											}
											position++
										}
									l247:
										{
											position249, tokenIndex249, depth249 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l250
											}
											position++
											goto l249
										l250:
											position, tokenIndex, depth = position249, tokenIndex249, depth249
											if buffer[position] != rune('A') {
												goto l227
											}
											position++
										}
									l249:
										{
											position251, tokenIndex251, depth251 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l252
											}
											position++
											goto l251
										l252:
											position, tokenIndex, depth = position251, tokenIndex251, depth251
											if buffer[position] != rune('P') {
												goto l227
											}
											position++
										}
									l251:
										{
											position253, tokenIndex253, depth253 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l254
											}
											position++
											goto l253
										l254:
											position, tokenIndex, depth = position253, tokenIndex253, depth253
											if buffer[position] != rune('H') {
												goto l227
											}
											position++
										}
									l253:
										if !rules[ruleskip]() {
											goto l227
										}
										depth--
										add(ruleGRAPH, position244)
									}
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										if !rules[rulevar]() {
											goto l256
										}
										goto l255
									l256:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
										if !rules[ruleiriref]() {
											goto l227
										}
									}
								l255:
									if !rules[rulegroupGraphPattern]() {
										goto l227
									}
									depth--
									add(rulegraphGraphPattern, position243)
								}
								break
							case '{':
								if !rules[rulegroupOrUnionGraphPattern]() {
									goto l227
								}
								break
							default:
								{
									position257 := position
									depth++
									{
										position258 := position
										depth++
										{
											position259, tokenIndex259, depth259 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l260
											}
											position++
											goto l259
										l260:
											position, tokenIndex, depth = position259, tokenIndex259, depth259
											if buffer[position] != rune('O') {
												goto l227
											}
											position++
										}
									l259:
										{
											position261, tokenIndex261, depth261 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l262
											}
											position++
											goto l261
										l262:
											position, tokenIndex, depth = position261, tokenIndex261, depth261
											if buffer[position] != rune('P') {
												goto l227
											}
											position++
										}
									l261:
										{
											position263, tokenIndex263, depth263 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l264
											}
											position++
											goto l263
										l264:
											position, tokenIndex, depth = position263, tokenIndex263, depth263
											if buffer[position] != rune('T') {
												goto l227
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
												goto l227
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
												goto l227
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
												goto l227
											}
											position++
										}
									l269:
										{
											position271, tokenIndex271, depth271 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l272
											}
											position++
											goto l271
										l272:
											position, tokenIndex, depth = position271, tokenIndex271, depth271
											if buffer[position] != rune('A') {
												goto l227
											}
											position++
										}
									l271:
										{
											position273, tokenIndex273, depth273 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l274
											}
											position++
											goto l273
										l274:
											position, tokenIndex, depth = position273, tokenIndex273, depth273
											if buffer[position] != rune('L') {
												goto l227
											}
											position++
										}
									l273:
										if !rules[ruleskip]() {
											goto l227
										}
										depth--
										add(ruleOPTIONAL, position258)
									}
									if !rules[ruleLBRACE]() {
										goto l227
									}
									{
										position275, tokenIndex275, depth275 := position, tokenIndex, depth
										if !rules[rulesubSelect]() {
											goto l276
										}
										goto l275
									l276:
										position, tokenIndex, depth = position275, tokenIndex275, depth275
										if !rules[rulegraphPattern]() {
											goto l227
										}
									}
								l275:
									if !rules[ruleRBRACE]() {
										goto l227
									}
									depth--
									add(ruleoptionalGraphPattern, position257)
								}
								break
							}
						}

						depth--
						add(rulegraphPatternNotTriples, position229)
					}
					{
						position277, tokenIndex277, depth277 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l277
						}
						goto l278
					l277:
						position, tokenIndex, depth = position277, tokenIndex277, depth277
					}
				l278:
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
		/* 18 graphPatternNotTriples <- <((&('M' | 'm') minusGraphPattern) | (&('G' | 'g') graphGraphPattern) | (&('{') groupOrUnionGraphPattern) | (&('O' | 'o') optionalGraphPattern))> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position281, tokenIndex281, depth281 := position, tokenIndex, depth
			{
				position282 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l281
				}
				{
					position283, tokenIndex283, depth283 := position, tokenIndex, depth
					{
						position285 := position
						depth++
						{
							position286, tokenIndex286, depth286 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l287
							}
							position++
							goto l286
						l287:
							position, tokenIndex, depth = position286, tokenIndex286, depth286
							if buffer[position] != rune('U') {
								goto l283
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
								goto l283
							}
							position++
						}
					l288:
						{
							position290, tokenIndex290, depth290 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l291
							}
							position++
							goto l290
						l291:
							position, tokenIndex, depth = position290, tokenIndex290, depth290
							if buffer[position] != rune('I') {
								goto l283
							}
							position++
						}
					l290:
						{
							position292, tokenIndex292, depth292 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l293
							}
							position++
							goto l292
						l293:
							position, tokenIndex, depth = position292, tokenIndex292, depth292
							if buffer[position] != rune('O') {
								goto l283
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
								goto l283
							}
							position++
						}
					l294:
						if !rules[ruleskip]() {
							goto l283
						}
						depth--
						add(ruleUNION, position285)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l283
					}
					goto l284
				l283:
					position, tokenIndex, depth = position283, tokenIndex283, depth283
				}
			l284:
				depth--
				add(rulegroupOrUnionGraphPattern, position282)
			}
			return true
		l281:
			position, tokenIndex, depth = position281, tokenIndex281, depth281
			return false
		},
		/* 21 graphGraphPattern <- <(GRAPH (var / iriref) groupGraphPattern)> */
		nil,
		/* 22 minusGraphPattern <- <(MINUSSETOPER groupGraphPattern)> */
		nil,
		/* 23 basicGraphPattern <- <((triplesBlock (filterOrBind DOT? triplesBlock?)*) / (filterOrBind DOT? triplesBlock?)+)> */
		nil,
		/* 24 filterOrBind <- <((FILTER constraint) / (BIND LPAREN expression AS var RPAREN))> */
		func() bool {
			position299, tokenIndex299, depth299 := position, tokenIndex, depth
			{
				position300 := position
				depth++
				{
					position301, tokenIndex301, depth301 := position, tokenIndex, depth
					{
						position303 := position
						depth++
						{
							position304, tokenIndex304, depth304 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l305
							}
							position++
							goto l304
						l305:
							position, tokenIndex, depth = position304, tokenIndex304, depth304
							if buffer[position] != rune('F') {
								goto l302
							}
							position++
						}
					l304:
						{
							position306, tokenIndex306, depth306 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l307
							}
							position++
							goto l306
						l307:
							position, tokenIndex, depth = position306, tokenIndex306, depth306
							if buffer[position] != rune('I') {
								goto l302
							}
							position++
						}
					l306:
						{
							position308, tokenIndex308, depth308 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l309
							}
							position++
							goto l308
						l309:
							position, tokenIndex, depth = position308, tokenIndex308, depth308
							if buffer[position] != rune('L') {
								goto l302
							}
							position++
						}
					l308:
						{
							position310, tokenIndex310, depth310 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l311
							}
							position++
							goto l310
						l311:
							position, tokenIndex, depth = position310, tokenIndex310, depth310
							if buffer[position] != rune('T') {
								goto l302
							}
							position++
						}
					l310:
						{
							position312, tokenIndex312, depth312 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l313
							}
							position++
							goto l312
						l313:
							position, tokenIndex, depth = position312, tokenIndex312, depth312
							if buffer[position] != rune('E') {
								goto l302
							}
							position++
						}
					l312:
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('R') {
								goto l302
							}
							position++
						}
					l314:
						if !rules[ruleskip]() {
							goto l302
						}
						depth--
						add(ruleFILTER, position303)
					}
					if !rules[ruleconstraint]() {
						goto l302
					}
					goto l301
				l302:
					position, tokenIndex, depth = position301, tokenIndex301, depth301
					{
						position316 := position
						depth++
						{
							position317, tokenIndex317, depth317 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l318
							}
							position++
							goto l317
						l318:
							position, tokenIndex, depth = position317, tokenIndex317, depth317
							if buffer[position] != rune('B') {
								goto l299
							}
							position++
						}
					l317:
						{
							position319, tokenIndex319, depth319 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l320
							}
							position++
							goto l319
						l320:
							position, tokenIndex, depth = position319, tokenIndex319, depth319
							if buffer[position] != rune('I') {
								goto l299
							}
							position++
						}
					l319:
						{
							position321, tokenIndex321, depth321 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l322
							}
							position++
							goto l321
						l322:
							position, tokenIndex, depth = position321, tokenIndex321, depth321
							if buffer[position] != rune('N') {
								goto l299
							}
							position++
						}
					l321:
						{
							position323, tokenIndex323, depth323 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l324
							}
							position++
							goto l323
						l324:
							position, tokenIndex, depth = position323, tokenIndex323, depth323
							if buffer[position] != rune('D') {
								goto l299
							}
							position++
						}
					l323:
						if !rules[ruleskip]() {
							goto l299
						}
						depth--
						add(ruleBIND, position316)
					}
					if !rules[ruleLPAREN]() {
						goto l299
					}
					if !rules[ruleexpression]() {
						goto l299
					}
					if !rules[ruleAS]() {
						goto l299
					}
					if !rules[rulevar]() {
						goto l299
					}
					if !rules[ruleRPAREN]() {
						goto l299
					}
				}
			l301:
				depth--
				add(rulefilterOrBind, position300)
			}
			return true
		l299:
			position, tokenIndex, depth = position299, tokenIndex299, depth299
			return false
		},
		/* 25 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				{
					position327, tokenIndex327, depth327 := position, tokenIndex, depth
					if !rules[rulebrackettedExpression]() {
						goto l328
					}
					goto l327
				l328:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if !rules[rulebuiltinCall]() {
						goto l329
					}
					goto l327
				l329:
					position, tokenIndex, depth = position327, tokenIndex327, depth327
					if !rules[rulefunctionCall]() {
						goto l325
					}
				}
			l327:
				depth--
				add(ruleconstraint, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
			return false
		},
		/* 26 triplesBlock <- <triplesSameSubjectPath+> */
		func() bool {
			position330, tokenIndex330, depth330 := position, tokenIndex, depth
			{
				position331 := position
				depth++
				{
					position334 := position
					depth++
					{
						position335, tokenIndex335, depth335 := position, tokenIndex, depth
						if !rules[rulevarOrTerm]() {
							goto l336
						}
						if !rules[rulepropertyListPath]() {
							goto l336
						}
						goto l335
					l336:
						position, tokenIndex, depth = position335, tokenIndex335, depth335
						if !rules[ruletriplesNodePath]() {
							goto l330
						}
						{
							position337, tokenIndex337, depth337 := position, tokenIndex, depth
							if !rules[rulepropertyListPath]() {
								goto l337
							}
							goto l338
						l337:
							position, tokenIndex, depth = position337, tokenIndex337, depth337
						}
					l338:
					}
				l335:
					{
						position339, tokenIndex339, depth339 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l339
						}
						goto l340
					l339:
						position, tokenIndex, depth = position339, tokenIndex339, depth339
					}
				l340:
					depth--
					add(ruletriplesSameSubjectPath, position334)
				}
			l332:
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					{
						position341 := position
						depth++
						{
							position342, tokenIndex342, depth342 := position, tokenIndex, depth
							if !rules[rulevarOrTerm]() {
								goto l343
							}
							if !rules[rulepropertyListPath]() {
								goto l343
							}
							goto l342
						l343:
							position, tokenIndex, depth = position342, tokenIndex342, depth342
							if !rules[ruletriplesNodePath]() {
								goto l333
							}
							{
								position344, tokenIndex344, depth344 := position, tokenIndex, depth
								if !rules[rulepropertyListPath]() {
									goto l344
								}
								goto l345
							l344:
								position, tokenIndex, depth = position344, tokenIndex344, depth344
							}
						l345:
						}
					l342:
						{
							position346, tokenIndex346, depth346 := position, tokenIndex, depth
							if !rules[ruleDOT]() {
								goto l346
							}
							goto l347
						l346:
							position, tokenIndex, depth = position346, tokenIndex346, depth346
						}
					l347:
						depth--
						add(ruletriplesSameSubjectPath, position341)
					}
					goto l332
				l333:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
				}
				depth--
				add(ruletriplesBlock, position331)
			}
			return true
		l330:
			position, tokenIndex, depth = position330, tokenIndex330, depth330
			return false
		},
		/* 27 triplesSameSubjectPath <- <(((varOrTerm propertyListPath) / (triplesNodePath propertyListPath?)) DOT?)> */
		nil,
		/* 28 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l352
					}
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					{
						position353 := position
						depth++
						{
							position354, tokenIndex354, depth354 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l355
							}
							goto l354
						l355:
							position, tokenIndex, depth = position354, tokenIndex354, depth354
							{
								switch buffer[position] {
								case '(':
									if !rules[rulenil]() {
										goto l349
									}
									break
								case '[', '_':
									{
										position357 := position
										depth++
										{
											position358, tokenIndex358, depth358 := position, tokenIndex, depth
											{
												position360 := position
												depth++
												if buffer[position] != rune('_') {
													goto l359
												}
												position++
												if buffer[position] != rune(':') {
													goto l359
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l359
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l359
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l359
														}
														position++
														break
													}
												}

												{
													position362, tokenIndex362, depth362 := position, tokenIndex, depth
													{
														position364, tokenIndex364, depth364 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l365
														}
														position++
														goto l364
													l365:
														position, tokenIndex, depth = position364, tokenIndex364, depth364
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l366
														}
														position++
														goto l364
													l366:
														position, tokenIndex, depth = position364, tokenIndex364, depth364
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l367
														}
														position++
														goto l364
													l367:
														position, tokenIndex, depth = position364, tokenIndex364, depth364
														if c := buffer[position]; c < rune('.') || c > rune('_') {
															goto l362
														}
														position++
													}
												l364:
													goto l363
												l362:
													position, tokenIndex, depth = position362, tokenIndex362, depth362
												}
											l363:
												if !rules[ruleskip]() {
													goto l359
												}
												depth--
												add(ruleblankNodeLabel, position360)
											}
											goto l358
										l359:
											position, tokenIndex, depth = position358, tokenIndex358, depth358
											{
												position368 := position
												depth++
												if buffer[position] != rune('[') {
													goto l349
												}
												position++
											l369:
												{
													position370, tokenIndex370, depth370 := position, tokenIndex, depth
													if !rules[rulews]() {
														goto l370
													}
													goto l369
												l370:
													position, tokenIndex, depth = position370, tokenIndex370, depth370
												}
												if buffer[position] != rune(']') {
													goto l349
												}
												position++
												if !rules[ruleskip]() {
													goto l349
												}
												depth--
												add(ruleanon, position368)
											}
										}
									l358:
										depth--
										add(ruleblankNode, position357)
									}
									break
								case 'F', 'T', 'f', 't':
									if !rules[rulebooleanLiteral]() {
										goto l349
									}
									break
								case '"':
									if !rules[ruleliteral]() {
										goto l349
									}
									break
								default:
									if !rules[rulenumericLiteral]() {
										goto l349
									}
									break
								}
							}

						}
					l354:
						depth--
						add(rulegraphTerm, position353)
					}
				}
			l351:
				depth--
				add(rulevarOrTerm, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 29 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 30 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		func() bool {
			position372, tokenIndex372, depth372 := position, tokenIndex, depth
			{
				position373 := position
				depth++
				{
					position374, tokenIndex374, depth374 := position, tokenIndex, depth
					{
						position376 := position
						depth++
						if !rules[ruleLPAREN]() {
							goto l375
						}
						if !rules[rulegraphNodePath]() {
							goto l375
						}
					l377:
						{
							position378, tokenIndex378, depth378 := position, tokenIndex, depth
							if !rules[rulegraphNodePath]() {
								goto l378
							}
							goto l377
						l378:
							position, tokenIndex, depth = position378, tokenIndex378, depth378
						}
						if !rules[ruleRPAREN]() {
							goto l375
						}
						depth--
						add(rulecollectionPath, position376)
					}
					goto l374
				l375:
					position, tokenIndex, depth = position374, tokenIndex374, depth374
					{
						position379 := position
						depth++
						{
							position380 := position
							depth++
							if buffer[position] != rune('[') {
								goto l372
							}
							position++
							if !rules[ruleskip]() {
								goto l372
							}
							depth--
							add(ruleLBRACK, position380)
						}
						if !rules[rulepropertyListPath]() {
							goto l372
						}
						{
							position381 := position
							depth++
							if buffer[position] != rune(']') {
								goto l372
							}
							position++
							if !rules[ruleskip]() {
								goto l372
							}
							depth--
							add(ruleRBRACK, position381)
						}
						depth--
						add(ruleblankNodePropertyListPath, position379)
					}
				}
			l374:
				depth--
				add(ruletriplesNodePath, position373)
			}
			return true
		l372:
			position, tokenIndex, depth = position372, tokenIndex372, depth372
			return false
		},
		/* 31 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 32 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 33 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position384, tokenIndex384, depth384 := position, tokenIndex, depth
			{
				position385 := position
				depth++
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l387
					}
					goto l386
				l387:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
					{
						position388 := position
						depth++
						if !rules[rulepath]() {
							goto l384
						}
						depth--
						add(ruleverbPath, position388)
					}
				}
			l386:
				{
					position389 := position
					depth++
					if !rules[ruleobjectPath]() {
						goto l384
					}
				l390:
					{
						position391, tokenIndex391, depth391 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l391
						}
						if !rules[ruleobjectPath]() {
							goto l391
						}
						goto l390
					l391:
						position, tokenIndex, depth = position391, tokenIndex391, depth391
					}
					depth--
					add(ruleobjectListPath, position389)
				}
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l392
					}
					{
						position394, tokenIndex394, depth394 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l394
						}
						goto l395
					l394:
						position, tokenIndex, depth = position394, tokenIndex394, depth394
					}
				l395:
					goto l393
				l392:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
				}
			l393:
				depth--
				add(rulepropertyListPath, position385)
			}
			return true
		l384:
			position, tokenIndex, depth = position384, tokenIndex384, depth384
			return false
		},
		/* 34 verbPath <- <path> */
		nil,
		/* 35 path <- <pathAlternative> */
		func() bool {
			position397, tokenIndex397, depth397 := position, tokenIndex, depth
			{
				position398 := position
				depth++
				{
					position399 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l397
					}
				l400:
					{
						position401, tokenIndex401, depth401 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l401
						}
						if !rules[rulepathSequence]() {
							goto l401
						}
						goto l400
					l401:
						position, tokenIndex, depth = position401, tokenIndex401, depth401
					}
					depth--
					add(rulepathAlternative, position399)
				}
				depth--
				add(rulepath, position398)
			}
			return true
		l397:
			position, tokenIndex, depth = position397, tokenIndex397, depth397
			return false
		},
		/* 36 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 37 pathSequence <- <(pathElt (SLASH pathElt)*)> */
		func() bool {
			position403, tokenIndex403, depth403 := position, tokenIndex, depth
			{
				position404 := position
				depth++
				if !rules[rulepathElt]() {
					goto l403
				}
			l405:
				{
					position406, tokenIndex406, depth406 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l406
					}
					if !rules[rulepathElt]() {
						goto l406
					}
					goto l405
				l406:
					position, tokenIndex, depth = position406, tokenIndex406, depth406
				}
				depth--
				add(rulepathSequence, position404)
			}
			return true
		l403:
			position, tokenIndex, depth = position403, tokenIndex403, depth403
			return false
		},
		/* 38 pathElt <- <(INVERSE? pathPrimary pathMod?)> */
		func() bool {
			position407, tokenIndex407, depth407 := position, tokenIndex, depth
			{
				position408 := position
				depth++
				{
					position409, tokenIndex409, depth409 := position, tokenIndex, depth
					if !rules[ruleINVERSE]() {
						goto l409
					}
					goto l410
				l409:
					position, tokenIndex, depth = position409, tokenIndex409, depth409
				}
			l410:
				{
					position411 := position
					depth++
					{
						position412, tokenIndex412, depth412 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l413
						}
						goto l412
					l413:
						position, tokenIndex, depth = position412, tokenIndex412, depth412
						{
							switch buffer[position] {
							case '(':
								if !rules[ruleLPAREN]() {
									goto l407
								}
								if !rules[rulepath]() {
									goto l407
								}
								if !rules[ruleRPAREN]() {
									goto l407
								}
								break
							case '!':
								if !rules[ruleNOT]() {
									goto l407
								}
								{
									position415 := position
									depth++
									{
										position416, tokenIndex416, depth416 := position, tokenIndex, depth
										if !rules[rulepathOneInPropertySet]() {
											goto l417
										}
										goto l416
									l417:
										position, tokenIndex, depth = position416, tokenIndex416, depth416
										if !rules[ruleLPAREN]() {
											goto l407
										}
										{
											position418, tokenIndex418, depth418 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l418
											}
										l420:
											{
												position421, tokenIndex421, depth421 := position, tokenIndex, depth
												if !rules[rulePIPE]() {
													goto l421
												}
												if !rules[rulepathOneInPropertySet]() {
													goto l421
												}
												goto l420
											l421:
												position, tokenIndex, depth = position421, tokenIndex421, depth421
											}
											goto l419
										l418:
											position, tokenIndex, depth = position418, tokenIndex418, depth418
										}
									l419:
										if !rules[ruleRPAREN]() {
											goto l407
										}
									}
								l416:
									depth--
									add(rulepathNegatedPropertySet, position415)
								}
								break
							default:
								if !rules[ruleISA]() {
									goto l407
								}
								break
							}
						}

					}
				l412:
					depth--
					add(rulepathPrimary, position411)
				}
				{
					position422, tokenIndex422, depth422 := position, tokenIndex, depth
					{
						position424 := position
						depth++
						{
							switch buffer[position] {
							case '+':
								if !rules[rulePLUS]() {
									goto l422
								}
								break
							case '?':
								{
									position426 := position
									depth++
									if buffer[position] != rune('?') {
										goto l422
									}
									position++
									if !rules[ruleskip]() {
										goto l422
									}
									depth--
									add(ruleQUESTION, position426)
								}
								break
							default:
								if !rules[ruleSTAR]() {
									goto l422
								}
								break
							}
						}

						{
							position427, tokenIndex427, depth427 := position, tokenIndex, depth
							if !rules[ruleskip]() {
								goto l427
							}
							goto l422
						l427:
							position, tokenIndex, depth = position427, tokenIndex427, depth427
						}
						depth--
						add(rulepathMod, position424)
					}
					goto l423
				l422:
					position, tokenIndex, depth = position422, tokenIndex422, depth422
				}
			l423:
				depth--
				add(rulepathElt, position408)
			}
			return true
		l407:
			position, tokenIndex, depth = position407, tokenIndex407, depth407
			return false
		},
		/* 39 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 40 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 41 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position430, tokenIndex430, depth430 := position, tokenIndex, depth
			{
				position431 := position
				depth++
				{
					position432, tokenIndex432, depth432 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l433
					}
					goto l432
				l433:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if !rules[ruleISA]() {
						goto l434
					}
					goto l432
				l434:
					position, tokenIndex, depth = position432, tokenIndex432, depth432
					if !rules[ruleINVERSE]() {
						goto l430
					}
					{
						position435, tokenIndex435, depth435 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l436
						}
						goto l435
					l436:
						position, tokenIndex, depth = position435, tokenIndex435, depth435
						if !rules[ruleISA]() {
							goto l430
						}
					}
				l435:
				}
			l432:
				depth--
				add(rulepathOneInPropertySet, position431)
			}
			return true
		l430:
			position, tokenIndex, depth = position430, tokenIndex430, depth430
			return false
		},
		/* 42 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !skip)> */
		nil,
		/* 43 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 44 objectPath <- <graphNodePath> */
		func() bool {
			position439, tokenIndex439, depth439 := position, tokenIndex, depth
			{
				position440 := position
				depth++
				if !rules[rulegraphNodePath]() {
					goto l439
				}
				depth--
				add(ruleobjectPath, position440)
			}
			return true
		l439:
			position, tokenIndex, depth = position439, tokenIndex439, depth439
			return false
		},
		/* 45 graphNodePath <- <(varOrTerm / triplesNodePath)> */
		func() bool {
			position441, tokenIndex441, depth441 := position, tokenIndex, depth
			{
				position442 := position
				depth++
				{
					position443, tokenIndex443, depth443 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l444
					}
					goto l443
				l444:
					position, tokenIndex, depth = position443, tokenIndex443, depth443
					if !rules[ruletriplesNodePath]() {
						goto l441
					}
				}
			l443:
				depth--
				add(rulegraphNodePath, position442)
			}
			return true
		l441:
			position, tokenIndex, depth = position441, tokenIndex441, depth441
			return false
		},
		/* 46 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position446 := position
				depth++
				{
					position447, tokenIndex447, depth447 := position, tokenIndex, depth
					{
						position449, tokenIndex449, depth449 := position, tokenIndex, depth
						{
							position451 := position
							depth++
							{
								position452, tokenIndex452, depth452 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l453
								}
								position++
								goto l452
							l453:
								position, tokenIndex, depth = position452, tokenIndex452, depth452
								if buffer[position] != rune('O') {
									goto l450
								}
								position++
							}
						l452:
							{
								position454, tokenIndex454, depth454 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l455
								}
								position++
								goto l454
							l455:
								position, tokenIndex, depth = position454, tokenIndex454, depth454
								if buffer[position] != rune('R') {
									goto l450
								}
								position++
							}
						l454:
							{
								position456, tokenIndex456, depth456 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l457
								}
								position++
								goto l456
							l457:
								position, tokenIndex, depth = position456, tokenIndex456, depth456
								if buffer[position] != rune('D') {
									goto l450
								}
								position++
							}
						l456:
							{
								position458, tokenIndex458, depth458 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l459
								}
								position++
								goto l458
							l459:
								position, tokenIndex, depth = position458, tokenIndex458, depth458
								if buffer[position] != rune('E') {
									goto l450
								}
								position++
							}
						l458:
							{
								position460, tokenIndex460, depth460 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l461
								}
								position++
								goto l460
							l461:
								position, tokenIndex, depth = position460, tokenIndex460, depth460
								if buffer[position] != rune('R') {
									goto l450
								}
								position++
							}
						l460:
							if !rules[ruleskip]() {
								goto l450
							}
							depth--
							add(ruleORDER, position451)
						}
						if !rules[ruleBY]() {
							goto l450
						}
						{
							position464 := position
							depth++
							{
								position465, tokenIndex465, depth465 := position, tokenIndex, depth
								{
									position467, tokenIndex467, depth467 := position, tokenIndex, depth
									{
										position469, tokenIndex469, depth469 := position, tokenIndex, depth
										{
											position471 := position
											depth++
											{
												position472, tokenIndex472, depth472 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l473
												}
												position++
												goto l472
											l473:
												position, tokenIndex, depth = position472, tokenIndex472, depth472
												if buffer[position] != rune('A') {
													goto l470
												}
												position++
											}
										l472:
											{
												position474, tokenIndex474, depth474 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l475
												}
												position++
												goto l474
											l475:
												position, tokenIndex, depth = position474, tokenIndex474, depth474
												if buffer[position] != rune('S') {
													goto l470
												}
												position++
											}
										l474:
											{
												position476, tokenIndex476, depth476 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l477
												}
												position++
												goto l476
											l477:
												position, tokenIndex, depth = position476, tokenIndex476, depth476
												if buffer[position] != rune('C') {
													goto l470
												}
												position++
											}
										l476:
											if !rules[ruleskip]() {
												goto l470
											}
											depth--
											add(ruleASC, position471)
										}
										goto l469
									l470:
										position, tokenIndex, depth = position469, tokenIndex469, depth469
										{
											position478 := position
											depth++
											{
												position479, tokenIndex479, depth479 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l480
												}
												position++
												goto l479
											l480:
												position, tokenIndex, depth = position479, tokenIndex479, depth479
												if buffer[position] != rune('D') {
													goto l467
												}
												position++
											}
										l479:
											{
												position481, tokenIndex481, depth481 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l482
												}
												position++
												goto l481
											l482:
												position, tokenIndex, depth = position481, tokenIndex481, depth481
												if buffer[position] != rune('E') {
													goto l467
												}
												position++
											}
										l481:
											{
												position483, tokenIndex483, depth483 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l484
												}
												position++
												goto l483
											l484:
												position, tokenIndex, depth = position483, tokenIndex483, depth483
												if buffer[position] != rune('S') {
													goto l467
												}
												position++
											}
										l483:
											{
												position485, tokenIndex485, depth485 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l486
												}
												position++
												goto l485
											l486:
												position, tokenIndex, depth = position485, tokenIndex485, depth485
												if buffer[position] != rune('C') {
													goto l467
												}
												position++
											}
										l485:
											if !rules[ruleskip]() {
												goto l467
											}
											depth--
											add(ruleDESC, position478)
										}
									}
								l469:
									goto l468
								l467:
									position, tokenIndex, depth = position467, tokenIndex467, depth467
								}
							l468:
								if !rules[rulebrackettedExpression]() {
									goto l466
								}
								goto l465
							l466:
								position, tokenIndex, depth = position465, tokenIndex465, depth465
								if !rules[rulefunctionCall]() {
									goto l487
								}
								goto l465
							l487:
								position, tokenIndex, depth = position465, tokenIndex465, depth465
								if !rules[rulebuiltinCall]() {
									goto l488
								}
								goto l465
							l488:
								position, tokenIndex, depth = position465, tokenIndex465, depth465
								if !rules[rulevar]() {
									goto l450
								}
							}
						l465:
							depth--
							add(ruleorderCondition, position464)
						}
					l462:
						{
							position463, tokenIndex463, depth463 := position, tokenIndex, depth
							{
								position489 := position
								depth++
								{
									position490, tokenIndex490, depth490 := position, tokenIndex, depth
									{
										position492, tokenIndex492, depth492 := position, tokenIndex, depth
										{
											position494, tokenIndex494, depth494 := position, tokenIndex, depth
											{
												position496 := position
												depth++
												{
													position497, tokenIndex497, depth497 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l498
													}
													position++
													goto l497
												l498:
													position, tokenIndex, depth = position497, tokenIndex497, depth497
													if buffer[position] != rune('A') {
														goto l495
													}
													position++
												}
											l497:
												{
													position499, tokenIndex499, depth499 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l500
													}
													position++
													goto l499
												l500:
													position, tokenIndex, depth = position499, tokenIndex499, depth499
													if buffer[position] != rune('S') {
														goto l495
													}
													position++
												}
											l499:
												{
													position501, tokenIndex501, depth501 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l502
													}
													position++
													goto l501
												l502:
													position, tokenIndex, depth = position501, tokenIndex501, depth501
													if buffer[position] != rune('C') {
														goto l495
													}
													position++
												}
											l501:
												if !rules[ruleskip]() {
													goto l495
												}
												depth--
												add(ruleASC, position496)
											}
											goto l494
										l495:
											position, tokenIndex, depth = position494, tokenIndex494, depth494
											{
												position503 := position
												depth++
												{
													position504, tokenIndex504, depth504 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l505
													}
													position++
													goto l504
												l505:
													position, tokenIndex, depth = position504, tokenIndex504, depth504
													if buffer[position] != rune('D') {
														goto l492
													}
													position++
												}
											l504:
												{
													position506, tokenIndex506, depth506 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l507
													}
													position++
													goto l506
												l507:
													position, tokenIndex, depth = position506, tokenIndex506, depth506
													if buffer[position] != rune('E') {
														goto l492
													}
													position++
												}
											l506:
												{
													position508, tokenIndex508, depth508 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l509
													}
													position++
													goto l508
												l509:
													position, tokenIndex, depth = position508, tokenIndex508, depth508
													if buffer[position] != rune('S') {
														goto l492
													}
													position++
												}
											l508:
												{
													position510, tokenIndex510, depth510 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l511
													}
													position++
													goto l510
												l511:
													position, tokenIndex, depth = position510, tokenIndex510, depth510
													if buffer[position] != rune('C') {
														goto l492
													}
													position++
												}
											l510:
												if !rules[ruleskip]() {
													goto l492
												}
												depth--
												add(ruleDESC, position503)
											}
										}
									l494:
										goto l493
									l492:
										position, tokenIndex, depth = position492, tokenIndex492, depth492
									}
								l493:
									if !rules[rulebrackettedExpression]() {
										goto l491
									}
									goto l490
								l491:
									position, tokenIndex, depth = position490, tokenIndex490, depth490
									if !rules[rulefunctionCall]() {
										goto l512
									}
									goto l490
								l512:
									position, tokenIndex, depth = position490, tokenIndex490, depth490
									if !rules[rulebuiltinCall]() {
										goto l513
									}
									goto l490
								l513:
									position, tokenIndex, depth = position490, tokenIndex490, depth490
									if !rules[rulevar]() {
										goto l463
									}
								}
							l490:
								depth--
								add(ruleorderCondition, position489)
							}
							goto l462
						l463:
							position, tokenIndex, depth = position463, tokenIndex463, depth463
						}
						goto l449
					l450:
						position, tokenIndex, depth = position449, tokenIndex449, depth449
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position515 := position
									depth++
									{
										position516, tokenIndex516, depth516 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l517
										}
										position++
										goto l516
									l517:
										position, tokenIndex, depth = position516, tokenIndex516, depth516
										if buffer[position] != rune('H') {
											goto l447
										}
										position++
									}
								l516:
									{
										position518, tokenIndex518, depth518 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l519
										}
										position++
										goto l518
									l519:
										position, tokenIndex, depth = position518, tokenIndex518, depth518
										if buffer[position] != rune('A') {
											goto l447
										}
										position++
									}
								l518:
									{
										position520, tokenIndex520, depth520 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l521
										}
										position++
										goto l520
									l521:
										position, tokenIndex, depth = position520, tokenIndex520, depth520
										if buffer[position] != rune('V') {
											goto l447
										}
										position++
									}
								l520:
									{
										position522, tokenIndex522, depth522 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l523
										}
										position++
										goto l522
									l523:
										position, tokenIndex, depth = position522, tokenIndex522, depth522
										if buffer[position] != rune('I') {
											goto l447
										}
										position++
									}
								l522:
									{
										position524, tokenIndex524, depth524 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l525
										}
										position++
										goto l524
									l525:
										position, tokenIndex, depth = position524, tokenIndex524, depth524
										if buffer[position] != rune('N') {
											goto l447
										}
										position++
									}
								l524:
									{
										position526, tokenIndex526, depth526 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l527
										}
										position++
										goto l526
									l527:
										position, tokenIndex, depth = position526, tokenIndex526, depth526
										if buffer[position] != rune('G') {
											goto l447
										}
										position++
									}
								l526:
									if !rules[ruleskip]() {
										goto l447
									}
									depth--
									add(ruleHAVING, position515)
								}
								if !rules[ruleconstraint]() {
									goto l447
								}
								break
							case 'G', 'g':
								{
									position528 := position
									depth++
									{
										position529, tokenIndex529, depth529 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l530
										}
										position++
										goto l529
									l530:
										position, tokenIndex, depth = position529, tokenIndex529, depth529
										if buffer[position] != rune('G') {
											goto l447
										}
										position++
									}
								l529:
									{
										position531, tokenIndex531, depth531 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l532
										}
										position++
										goto l531
									l532:
										position, tokenIndex, depth = position531, tokenIndex531, depth531
										if buffer[position] != rune('R') {
											goto l447
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
											goto l447
										}
										position++
									}
								l533:
									{
										position535, tokenIndex535, depth535 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l536
										}
										position++
										goto l535
									l536:
										position, tokenIndex, depth = position535, tokenIndex535, depth535
										if buffer[position] != rune('U') {
											goto l447
										}
										position++
									}
								l535:
									{
										position537, tokenIndex537, depth537 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l538
										}
										position++
										goto l537
									l538:
										position, tokenIndex, depth = position537, tokenIndex537, depth537
										if buffer[position] != rune('P') {
											goto l447
										}
										position++
									}
								l537:
									if !rules[ruleskip]() {
										goto l447
									}
									depth--
									add(ruleGROUP, position528)
								}
								if !rules[ruleBY]() {
									goto l447
								}
								{
									position541 := position
									depth++
									{
										position542, tokenIndex542, depth542 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l543
										}
										goto l542
									l543:
										position, tokenIndex, depth = position542, tokenIndex542, depth542
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l447
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l447
												}
												if !rules[ruleexpression]() {
													goto l447
												}
												{
													position545, tokenIndex545, depth545 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l545
													}
													if !rules[rulevar]() {
														goto l545
													}
													goto l546
												l545:
													position, tokenIndex, depth = position545, tokenIndex545, depth545
												}
											l546:
												if !rules[ruleRPAREN]() {
													goto l447
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l447
												}
												break
											}
										}

									}
								l542:
									depth--
									add(rulegroupCondition, position541)
								}
							l539:
								{
									position540, tokenIndex540, depth540 := position, tokenIndex, depth
									{
										position547 := position
										depth++
										{
											position548, tokenIndex548, depth548 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l549
											}
											goto l548
										l549:
											position, tokenIndex, depth = position548, tokenIndex548, depth548
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l540
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l540
													}
													if !rules[ruleexpression]() {
														goto l540
													}
													{
														position551, tokenIndex551, depth551 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l551
														}
														if !rules[rulevar]() {
															goto l551
														}
														goto l552
													l551:
														position, tokenIndex, depth = position551, tokenIndex551, depth551
													}
												l552:
													if !rules[ruleRPAREN]() {
														goto l540
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l540
													}
													break
												}
											}

										}
									l548:
										depth--
										add(rulegroupCondition, position547)
									}
									goto l539
								l540:
									position, tokenIndex, depth = position540, tokenIndex540, depth540
								}
								break
							default:
								{
									position553 := position
									depth++
									{
										position554, tokenIndex554, depth554 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l555
										}
										{
											position556, tokenIndex556, depth556 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l556
											}
											goto l557
										l556:
											position, tokenIndex, depth = position556, tokenIndex556, depth556
										}
									l557:
										goto l554
									l555:
										position, tokenIndex, depth = position554, tokenIndex554, depth554
										if !rules[ruleoffset]() {
											goto l447
										}
										{
											position558, tokenIndex558, depth558 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l558
											}
											goto l559
										l558:
											position, tokenIndex, depth = position558, tokenIndex558, depth558
										}
									l559:
									}
								l554:
									depth--
									add(rulelimitOffsetClauses, position553)
								}
								break
							}
						}

					}
				l449:
					goto l448
				l447:
					position, tokenIndex, depth = position447, tokenIndex447, depth447
				}
			l448:
				depth--
				add(rulesolutionModifier, position446)
			}
			return true
		},
		/* 47 groupCondition <- <(functionCall / ((&('$' | '?') var) | (&('(') (LPAREN expression (AS var)? RPAREN)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'H' | 'I' | 'L' | 'M' | 'N' | 'R' | 'S' | 'T' | 'U' | 'Y' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'h' | 'i' | 'l' | 'm' | 'n' | 'r' | 's' | 't' | 'u' | 'y') builtinCall)))> */
		nil,
		/* 48 orderCondition <- <(((ASC / DESC)? brackettedExpression) / functionCall / builtinCall / var)> */
		nil,
		/* 49 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 50 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position563, tokenIndex563, depth563 := position, tokenIndex, depth
			{
				position564 := position
				depth++
				{
					position565 := position
					depth++
					{
						position566, tokenIndex566, depth566 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l567
						}
						position++
						goto l566
					l567:
						position, tokenIndex, depth = position566, tokenIndex566, depth566
						if buffer[position] != rune('L') {
							goto l563
						}
						position++
					}
				l566:
					{
						position568, tokenIndex568, depth568 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l569
						}
						position++
						goto l568
					l569:
						position, tokenIndex, depth = position568, tokenIndex568, depth568
						if buffer[position] != rune('I') {
							goto l563
						}
						position++
					}
				l568:
					{
						position570, tokenIndex570, depth570 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l571
						}
						position++
						goto l570
					l571:
						position, tokenIndex, depth = position570, tokenIndex570, depth570
						if buffer[position] != rune('M') {
							goto l563
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
							goto l563
						}
						position++
					}
				l572:
					{
						position574, tokenIndex574, depth574 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l575
						}
						position++
						goto l574
					l575:
						position, tokenIndex, depth = position574, tokenIndex574, depth574
						if buffer[position] != rune('T') {
							goto l563
						}
						position++
					}
				l574:
					if !rules[ruleskip]() {
						goto l563
					}
					depth--
					add(ruleLIMIT, position565)
				}
				if !rules[ruleINTEGER]() {
					goto l563
				}
				depth--
				add(rulelimit, position564)
			}
			return true
		l563:
			position, tokenIndex, depth = position563, tokenIndex563, depth563
			return false
		},
		/* 51 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position576, tokenIndex576, depth576 := position, tokenIndex, depth
			{
				position577 := position
				depth++
				{
					position578 := position
					depth++
					{
						position579, tokenIndex579, depth579 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l580
						}
						position++
						goto l579
					l580:
						position, tokenIndex, depth = position579, tokenIndex579, depth579
						if buffer[position] != rune('O') {
							goto l576
						}
						position++
					}
				l579:
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l582
						}
						position++
						goto l581
					l582:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
						if buffer[position] != rune('F') {
							goto l576
						}
						position++
					}
				l581:
					{
						position583, tokenIndex583, depth583 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l584
						}
						position++
						goto l583
					l584:
						position, tokenIndex, depth = position583, tokenIndex583, depth583
						if buffer[position] != rune('F') {
							goto l576
						}
						position++
					}
				l583:
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
							goto l576
						}
						position++
					}
				l585:
					{
						position587, tokenIndex587, depth587 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l588
						}
						position++
						goto l587
					l588:
						position, tokenIndex, depth = position587, tokenIndex587, depth587
						if buffer[position] != rune('E') {
							goto l576
						}
						position++
					}
				l587:
					{
						position589, tokenIndex589, depth589 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l590
						}
						position++
						goto l589
					l590:
						position, tokenIndex, depth = position589, tokenIndex589, depth589
						if buffer[position] != rune('T') {
							goto l576
						}
						position++
					}
				l589:
					if !rules[ruleskip]() {
						goto l576
					}
					depth--
					add(ruleOFFSET, position578)
				}
				if !rules[ruleINTEGER]() {
					goto l576
				}
				depth--
				add(ruleoffset, position577)
			}
			return true
		l576:
			position, tokenIndex, depth = position576, tokenIndex576, depth576
			return false
		},
		/* 52 expression <- <conditionalOrExpression> */
		func() bool {
			position591, tokenIndex591, depth591 := position, tokenIndex, depth
			{
				position592 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l591
				}
				depth--
				add(ruleexpression, position592)
			}
			return true
		l591:
			position, tokenIndex, depth = position591, tokenIndex591, depth591
			return false
		},
		/* 53 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position593, tokenIndex593, depth593 := position, tokenIndex, depth
			{
				position594 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l593
				}
				{
					position595, tokenIndex595, depth595 := position, tokenIndex, depth
					{
						position597 := position
						depth++
						if buffer[position] != rune('|') {
							goto l595
						}
						position++
						if buffer[position] != rune('|') {
							goto l595
						}
						position++
						if !rules[ruleskip]() {
							goto l595
						}
						depth--
						add(ruleOR, position597)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l595
					}
					goto l596
				l595:
					position, tokenIndex, depth = position595, tokenIndex595, depth595
				}
			l596:
				depth--
				add(ruleconditionalOrExpression, position594)
			}
			return true
		l593:
			position, tokenIndex, depth = position593, tokenIndex593, depth593
			return false
		},
		/* 54 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position598, tokenIndex598, depth598 := position, tokenIndex, depth
			{
				position599 := position
				depth++
				{
					position600 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l598
					}
					{
						position601, tokenIndex601, depth601 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position604 := position
									depth++
									{
										position605 := position
										depth++
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
												goto l601
											}
											position++
										}
									l606:
										{
											position608, tokenIndex608, depth608 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l609
											}
											position++
											goto l608
										l609:
											position, tokenIndex, depth = position608, tokenIndex608, depth608
											if buffer[position] != rune('O') {
												goto l601
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
												goto l601
											}
											position++
										}
									l610:
										if buffer[position] != rune(' ') {
											goto l601
										}
										position++
										{
											position612, tokenIndex612, depth612 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l613
											}
											position++
											goto l612
										l613:
											position, tokenIndex, depth = position612, tokenIndex612, depth612
											if buffer[position] != rune('I') {
												goto l601
											}
											position++
										}
									l612:
										{
											position614, tokenIndex614, depth614 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l615
											}
											position++
											goto l614
										l615:
											position, tokenIndex, depth = position614, tokenIndex614, depth614
											if buffer[position] != rune('N') {
												goto l601
											}
											position++
										}
									l614:
										if !rules[ruleskip]() {
											goto l601
										}
										depth--
										add(ruleNOTIN, position605)
									}
									if !rules[ruleargList]() {
										goto l601
									}
									depth--
									add(rulenotin, position604)
								}
								break
							case 'I', 'i':
								{
									position616 := position
									depth++
									{
										position617 := position
										depth++
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
												goto l601
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
												goto l601
											}
											position++
										}
									l620:
										if !rules[ruleskip]() {
											goto l601
										}
										depth--
										add(ruleIN, position617)
									}
									if !rules[ruleargList]() {
										goto l601
									}
									depth--
									add(rulein, position616)
								}
								break
							default:
								{
									position622, tokenIndex622, depth622 := position, tokenIndex, depth
									{
										position624 := position
										depth++
										if buffer[position] != rune('<') {
											goto l623
										}
										position++
										if !rules[ruleskip]() {
											goto l623
										}
										depth--
										add(ruleLT, position624)
									}
									goto l622
								l623:
									position, tokenIndex, depth = position622, tokenIndex622, depth622
									{
										position626 := position
										depth++
										if buffer[position] != rune('>') {
											goto l625
										}
										position++
										if buffer[position] != rune('=') {
											goto l625
										}
										position++
										if !rules[ruleskip]() {
											goto l625
										}
										depth--
										add(ruleGE, position626)
									}
									goto l622
								l625:
									position, tokenIndex, depth = position622, tokenIndex622, depth622
									{
										switch buffer[position] {
										case '>':
											{
												position628 := position
												depth++
												if buffer[position] != rune('>') {
													goto l601
												}
												position++
												if !rules[ruleskip]() {
													goto l601
												}
												depth--
												add(ruleGT, position628)
											}
											break
										case '<':
											{
												position629 := position
												depth++
												if buffer[position] != rune('<') {
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
												add(ruleLE, position629)
											}
											break
										case '!':
											{
												position630 := position
												depth++
												if buffer[position] != rune('!') {
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
												add(ruleNE, position630)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l601
											}
											break
										}
									}

								}
							l622:
								if !rules[rulenumericExpression]() {
									goto l601
								}
								break
							}
						}

						goto l602
					l601:
						position, tokenIndex, depth = position601, tokenIndex601, depth601
					}
				l602:
					depth--
					add(rulevalueLogical, position600)
				}
				{
					position631, tokenIndex631, depth631 := position, tokenIndex, depth
					{
						position633 := position
						depth++
						if buffer[position] != rune('&') {
							goto l631
						}
						position++
						if buffer[position] != rune('&') {
							goto l631
						}
						position++
						if !rules[ruleskip]() {
							goto l631
						}
						depth--
						add(ruleAND, position633)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l631
					}
					goto l632
				l631:
					position, tokenIndex, depth = position631, tokenIndex631, depth631
				}
			l632:
				depth--
				add(ruleconditionalAndExpression, position599)
			}
			return true
		l598:
			position, tokenIndex, depth = position598, tokenIndex598, depth598
			return false
		},
		/* 55 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 56 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l635
				}
			l637:
				{
					position638, tokenIndex638, depth638 := position, tokenIndex, depth
					{
						position639, tokenIndex639, depth639 := position, tokenIndex, depth
						{
							position641, tokenIndex641, depth641 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l642
							}
							goto l641
						l642:
							position, tokenIndex, depth = position641, tokenIndex641, depth641
							if !rules[ruleMINUS]() {
								goto l640
							}
						}
					l641:
						if !rules[rulemultiplicativeExpression]() {
							goto l640
						}
						goto l639
					l640:
						position, tokenIndex, depth = position639, tokenIndex639, depth639
						{
							position643 := position
							depth++
							{
								position644, tokenIndex644, depth644 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l645
								}
								position++
								goto l644
							l645:
								position, tokenIndex, depth = position644, tokenIndex644, depth644
								if buffer[position] != rune('-') {
									goto l638
								}
								position++
							}
						l644:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l638
							}
							position++
						l646:
							{
								position647, tokenIndex647, depth647 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l647
								}
								position++
								goto l646
							l647:
								position, tokenIndex, depth = position647, tokenIndex647, depth647
							}
							{
								position648, tokenIndex648, depth648 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l648
								}
								position++
							l650:
								{
									position651, tokenIndex651, depth651 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l651
									}
									position++
									goto l650
								l651:
									position, tokenIndex, depth = position651, tokenIndex651, depth651
								}
								goto l649
							l648:
								position, tokenIndex, depth = position648, tokenIndex648, depth648
							}
						l649:
							if !rules[ruleskip]() {
								goto l638
							}
							depth--
							add(rulesignedNumericLiteral, position643)
						}
					}
				l639:
					goto l637
				l638:
					position, tokenIndex, depth = position638, tokenIndex638, depth638
				}
				depth--
				add(rulenumericExpression, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 57 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l652
				}
			l654:
				{
					position655, tokenIndex655, depth655 := position, tokenIndex, depth
					{
						position656, tokenIndex656, depth656 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l657
						}
						goto l656
					l657:
						position, tokenIndex, depth = position656, tokenIndex656, depth656
						if !rules[ruleSLASH]() {
							goto l655
						}
					}
				l656:
					if !rules[ruleunaryExpression]() {
						goto l655
					}
					goto l654
				l655:
					position, tokenIndex, depth = position655, tokenIndex655, depth655
				}
				depth--
				add(rulemultiplicativeExpression, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 58 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				{
					position660, tokenIndex660, depth660 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l660
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l660
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l660
							}
							break
						}
					}

					goto l661
				l660:
					position, tokenIndex, depth = position660, tokenIndex660, depth660
				}
			l661:
				{
					position663 := position
					depth++
					{
						position664, tokenIndex664, depth664 := position, tokenIndex, depth
						if !rules[rulebrackettedExpression]() {
							goto l665
						}
						goto l664
					l665:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if !rules[rulebuiltinCall]() {
							goto l666
						}
						goto l664
					l666:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if !rules[rulefunctionCall]() {
							goto l667
						}
						goto l664
					l667:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						if !rules[ruleiriref]() {
							goto l668
						}
						goto l664
					l668:
						position, tokenIndex, depth = position664, tokenIndex664, depth664
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position670 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position672 := position
												depth++
												{
													position673 := position
													depth++
													{
														position674, tokenIndex674, depth674 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l675
														}
														position++
														goto l674
													l675:
														position, tokenIndex, depth = position674, tokenIndex674, depth674
														if buffer[position] != rune('G') {
															goto l658
														}
														position++
													}
												l674:
													{
														position676, tokenIndex676, depth676 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l677
														}
														position++
														goto l676
													l677:
														position, tokenIndex, depth = position676, tokenIndex676, depth676
														if buffer[position] != rune('R') {
															goto l658
														}
														position++
													}
												l676:
													{
														position678, tokenIndex678, depth678 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l679
														}
														position++
														goto l678
													l679:
														position, tokenIndex, depth = position678, tokenIndex678, depth678
														if buffer[position] != rune('O') {
															goto l658
														}
														position++
													}
												l678:
													{
														position680, tokenIndex680, depth680 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l681
														}
														position++
														goto l680
													l681:
														position, tokenIndex, depth = position680, tokenIndex680, depth680
														if buffer[position] != rune('U') {
															goto l658
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
															goto l658
														}
														position++
													}
												l682:
													if buffer[position] != rune('_') {
														goto l658
													}
													position++
													{
														position684, tokenIndex684, depth684 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l685
														}
														position++
														goto l684
													l685:
														position, tokenIndex, depth = position684, tokenIndex684, depth684
														if buffer[position] != rune('C') {
															goto l658
														}
														position++
													}
												l684:
													{
														position686, tokenIndex686, depth686 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l687
														}
														position++
														goto l686
													l687:
														position, tokenIndex, depth = position686, tokenIndex686, depth686
														if buffer[position] != rune('O') {
															goto l658
														}
														position++
													}
												l686:
													{
														position688, tokenIndex688, depth688 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l689
														}
														position++
														goto l688
													l689:
														position, tokenIndex, depth = position688, tokenIndex688, depth688
														if buffer[position] != rune('N') {
															goto l658
														}
														position++
													}
												l688:
													{
														position690, tokenIndex690, depth690 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l691
														}
														position++
														goto l690
													l691:
														position, tokenIndex, depth = position690, tokenIndex690, depth690
														if buffer[position] != rune('C') {
															goto l658
														}
														position++
													}
												l690:
													{
														position692, tokenIndex692, depth692 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l693
														}
														position++
														goto l692
													l693:
														position, tokenIndex, depth = position692, tokenIndex692, depth692
														if buffer[position] != rune('A') {
															goto l658
														}
														position++
													}
												l692:
													{
														position694, tokenIndex694, depth694 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l695
														}
														position++
														goto l694
													l695:
														position, tokenIndex, depth = position694, tokenIndex694, depth694
														if buffer[position] != rune('T') {
															goto l658
														}
														position++
													}
												l694:
													if !rules[ruleskip]() {
														goto l658
													}
													depth--
													add(ruleGROUPCONCAT, position673)
												}
												if !rules[ruleLPAREN]() {
													goto l658
												}
												{
													position696, tokenIndex696, depth696 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l696
													}
													goto l697
												l696:
													position, tokenIndex, depth = position696, tokenIndex696, depth696
												}
											l697:
												if !rules[ruleexpression]() {
													goto l658
												}
												{
													position698, tokenIndex698, depth698 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l698
													}
													{
														position700 := position
														depth++
														{
															position701, tokenIndex701, depth701 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l702
															}
															position++
															goto l701
														l702:
															position, tokenIndex, depth = position701, tokenIndex701, depth701
															if buffer[position] != rune('S') {
																goto l698
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
																goto l698
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
																goto l698
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
																goto l698
															}
															position++
														}
													l707:
														{
															position709, tokenIndex709, depth709 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l710
															}
															position++
															goto l709
														l710:
															position, tokenIndex, depth = position709, tokenIndex709, depth709
															if buffer[position] != rune('R') {
																goto l698
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
																goto l698
															}
															position++
														}
													l711:
														{
															position713, tokenIndex713, depth713 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l714
															}
															position++
															goto l713
														l714:
															position, tokenIndex, depth = position713, tokenIndex713, depth713
															if buffer[position] != rune('T') {
																goto l698
															}
															position++
														}
													l713:
														{
															position715, tokenIndex715, depth715 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l716
															}
															position++
															goto l715
														l716:
															position, tokenIndex, depth = position715, tokenIndex715, depth715
															if buffer[position] != rune('O') {
																goto l698
															}
															position++
														}
													l715:
														{
															position717, tokenIndex717, depth717 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l718
															}
															position++
															goto l717
														l718:
															position, tokenIndex, depth = position717, tokenIndex717, depth717
															if buffer[position] != rune('R') {
																goto l698
															}
															position++
														}
													l717:
														if !rules[ruleskip]() {
															goto l698
														}
														depth--
														add(ruleSEPARATOR, position700)
													}
													if !rules[ruleEQ]() {
														goto l698
													}
													if !rules[rulestring]() {
														goto l698
													}
													goto l699
												l698:
													position, tokenIndex, depth = position698, tokenIndex698, depth698
												}
											l699:
												if !rules[ruleRPAREN]() {
													goto l658
												}
												depth--
												add(rulegroupConcat, position672)
											}
											break
										case 'C', 'c':
											{
												position719 := position
												depth++
												{
													position720 := position
													depth++
													{
														position721, tokenIndex721, depth721 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l722
														}
														position++
														goto l721
													l722:
														position, tokenIndex, depth = position721, tokenIndex721, depth721
														if buffer[position] != rune('C') {
															goto l658
														}
														position++
													}
												l721:
													{
														position723, tokenIndex723, depth723 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l724
														}
														position++
														goto l723
													l724:
														position, tokenIndex, depth = position723, tokenIndex723, depth723
														if buffer[position] != rune('O') {
															goto l658
														}
														position++
													}
												l723:
													{
														position725, tokenIndex725, depth725 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l726
														}
														position++
														goto l725
													l726:
														position, tokenIndex, depth = position725, tokenIndex725, depth725
														if buffer[position] != rune('U') {
															goto l658
														}
														position++
													}
												l725:
													{
														position727, tokenIndex727, depth727 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l728
														}
														position++
														goto l727
													l728:
														position, tokenIndex, depth = position727, tokenIndex727, depth727
														if buffer[position] != rune('N') {
															goto l658
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
															goto l658
														}
														position++
													}
												l729:
													if !rules[ruleskip]() {
														goto l658
													}
													depth--
													add(ruleCOUNT, position720)
												}
												if !rules[ruleLPAREN]() {
													goto l658
												}
												{
													position731, tokenIndex731, depth731 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l731
													}
													goto l732
												l731:
													position, tokenIndex, depth = position731, tokenIndex731, depth731
												}
											l732:
												{
													position733, tokenIndex733, depth733 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l734
													}
													goto l733
												l734:
													position, tokenIndex, depth = position733, tokenIndex733, depth733
													if !rules[ruleexpression]() {
														goto l658
													}
												}
											l733:
												if !rules[ruleRPAREN]() {
													goto l658
												}
												depth--
												add(rulecount, position719)
											}
											break
										default:
											{
												position735, tokenIndex735, depth735 := position, tokenIndex, depth
												{
													position737 := position
													depth++
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
															goto l736
														}
														position++
													}
												l738:
													{
														position740, tokenIndex740, depth740 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l741
														}
														position++
														goto l740
													l741:
														position, tokenIndex, depth = position740, tokenIndex740, depth740
														if buffer[position] != rune('U') {
															goto l736
														}
														position++
													}
												l740:
													{
														position742, tokenIndex742, depth742 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l743
														}
														position++
														goto l742
													l743:
														position, tokenIndex, depth = position742, tokenIndex742, depth742
														if buffer[position] != rune('M') {
															goto l736
														}
														position++
													}
												l742:
													if !rules[ruleskip]() {
														goto l736
													}
													depth--
													add(ruleSUM, position737)
												}
												goto l735
											l736:
												position, tokenIndex, depth = position735, tokenIndex735, depth735
												{
													position745 := position
													depth++
													{
														position746, tokenIndex746, depth746 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l747
														}
														position++
														goto l746
													l747:
														position, tokenIndex, depth = position746, tokenIndex746, depth746
														if buffer[position] != rune('M') {
															goto l744
														}
														position++
													}
												l746:
													{
														position748, tokenIndex748, depth748 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l749
														}
														position++
														goto l748
													l749:
														position, tokenIndex, depth = position748, tokenIndex748, depth748
														if buffer[position] != rune('I') {
															goto l744
														}
														position++
													}
												l748:
													{
														position750, tokenIndex750, depth750 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l751
														}
														position++
														goto l750
													l751:
														position, tokenIndex, depth = position750, tokenIndex750, depth750
														if buffer[position] != rune('N') {
															goto l744
														}
														position++
													}
												l750:
													if !rules[ruleskip]() {
														goto l744
													}
													depth--
													add(ruleMIN, position745)
												}
												goto l735
											l744:
												position, tokenIndex, depth = position735, tokenIndex735, depth735
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position753 := position
															depth++
															{
																position754, tokenIndex754, depth754 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l755
																}
																position++
																goto l754
															l755:
																position, tokenIndex, depth = position754, tokenIndex754, depth754
																if buffer[position] != rune('S') {
																	goto l658
																}
																position++
															}
														l754:
															{
																position756, tokenIndex756, depth756 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l757
																}
																position++
																goto l756
															l757:
																position, tokenIndex, depth = position756, tokenIndex756, depth756
																if buffer[position] != rune('A') {
																	goto l658
																}
																position++
															}
														l756:
															{
																position758, tokenIndex758, depth758 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l759
																}
																position++
																goto l758
															l759:
																position, tokenIndex, depth = position758, tokenIndex758, depth758
																if buffer[position] != rune('M') {
																	goto l658
																}
																position++
															}
														l758:
															{
																position760, tokenIndex760, depth760 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l761
																}
																position++
																goto l760
															l761:
																position, tokenIndex, depth = position760, tokenIndex760, depth760
																if buffer[position] != rune('P') {
																	goto l658
																}
																position++
															}
														l760:
															{
																position762, tokenIndex762, depth762 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l763
																}
																position++
																goto l762
															l763:
																position, tokenIndex, depth = position762, tokenIndex762, depth762
																if buffer[position] != rune('L') {
																	goto l658
																}
																position++
															}
														l762:
															{
																position764, tokenIndex764, depth764 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l765
																}
																position++
																goto l764
															l765:
																position, tokenIndex, depth = position764, tokenIndex764, depth764
																if buffer[position] != rune('E') {
																	goto l658
																}
																position++
															}
														l764:
															if !rules[ruleskip]() {
																goto l658
															}
															depth--
															add(ruleSAMPLE, position753)
														}
														break
													case 'A', 'a':
														{
															position766 := position
															depth++
															{
																position767, tokenIndex767, depth767 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l768
																}
																position++
																goto l767
															l768:
																position, tokenIndex, depth = position767, tokenIndex767, depth767
																if buffer[position] != rune('A') {
																	goto l658
																}
																position++
															}
														l767:
															{
																position769, tokenIndex769, depth769 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l770
																}
																position++
																goto l769
															l770:
																position, tokenIndex, depth = position769, tokenIndex769, depth769
																if buffer[position] != rune('V') {
																	goto l658
																}
																position++
															}
														l769:
															{
																position771, tokenIndex771, depth771 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l772
																}
																position++
																goto l771
															l772:
																position, tokenIndex, depth = position771, tokenIndex771, depth771
																if buffer[position] != rune('G') {
																	goto l658
																}
																position++
															}
														l771:
															if !rules[ruleskip]() {
																goto l658
															}
															depth--
															add(ruleAVG, position766)
														}
														break
													default:
														{
															position773 := position
															depth++
															{
																position774, tokenIndex774, depth774 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l775
																}
																position++
																goto l774
															l775:
																position, tokenIndex, depth = position774, tokenIndex774, depth774
																if buffer[position] != rune('M') {
																	goto l658
																}
																position++
															}
														l774:
															{
																position776, tokenIndex776, depth776 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l777
																}
																position++
																goto l776
															l777:
																position, tokenIndex, depth = position776, tokenIndex776, depth776
																if buffer[position] != rune('A') {
																	goto l658
																}
																position++
															}
														l776:
															{
																position778, tokenIndex778, depth778 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l779
																}
																position++
																goto l778
															l779:
																position, tokenIndex, depth = position778, tokenIndex778, depth778
																if buffer[position] != rune('X') {
																	goto l658
																}
																position++
															}
														l778:
															if !rules[ruleskip]() {
																goto l658
															}
															depth--
															add(ruleMAX, position773)
														}
														break
													}
												}

											}
										l735:
											if !rules[ruleLPAREN]() {
												goto l658
											}
											{
												position780, tokenIndex780, depth780 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l780
												}
												goto l781
											l780:
												position, tokenIndex, depth = position780, tokenIndex780, depth780
											}
										l781:
											if !rules[ruleexpression]() {
												goto l658
											}
											if !rules[ruleRPAREN]() {
												goto l658
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position670)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l658
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l658
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l658
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l658
								}
								break
							}
						}

					}
				l664:
					depth--
					add(ruleprimaryExpression, position663)
				}
				depth--
				add(ruleunaryExpression, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 59 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 60 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position783, tokenIndex783, depth783 := position, tokenIndex, depth
			{
				position784 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l783
				}
				if !rules[ruleexpression]() {
					goto l783
				}
				if !rules[ruleRPAREN]() {
					goto l783
				}
				depth--
				add(rulebrackettedExpression, position784)
			}
			return true
		l783:
			position, tokenIndex, depth = position783, tokenIndex783, depth783
			return false
		},
		/* 61 functionCall <- <(iriref argList)> */
		func() bool {
			position785, tokenIndex785, depth785 := position, tokenIndex, depth
			{
				position786 := position
				depth++
				if !rules[ruleiriref]() {
					goto l785
				}
				if !rules[ruleargList]() {
					goto l785
				}
				depth--
				add(rulefunctionCall, position786)
			}
			return true
		l785:
			position, tokenIndex, depth = position785, tokenIndex785, depth785
			return false
		},
		/* 62 in <- <(IN argList)> */
		nil,
		/* 63 notin <- <(NOTIN argList)> */
		nil,
		/* 64 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position789, tokenIndex789, depth789 := position, tokenIndex, depth
			{
				position790 := position
				depth++
				{
					position791, tokenIndex791, depth791 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l792
					}
					goto l791
				l792:
					position, tokenIndex, depth = position791, tokenIndex791, depth791
					if !rules[ruleLPAREN]() {
						goto l789
					}
					if !rules[ruleexpression]() {
						goto l789
					}
				l793:
					{
						position794, tokenIndex794, depth794 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l794
						}
						if !rules[ruleexpression]() {
							goto l794
						}
						goto l793
					l794:
						position, tokenIndex, depth = position794, tokenIndex794, depth794
					}
					if !rules[ruleRPAREN]() {
						goto l789
					}
				}
			l791:
				depth--
				add(ruleargList, position790)
			}
			return true
		l789:
			position, tokenIndex, depth = position789, tokenIndex789, depth789
			return false
		},
		/* 65 aggregate <- <((&('G' | 'g') groupConcat) | (&('C' | 'c') count) | (&('A' | 'M' | 'S' | 'a' | 'm' | 's') ((SUM / MIN / ((&('S' | 's') SAMPLE) | (&('A' | 'a') AVG) | (&('M' | 'm') MAX))) LPAREN DISTINCT? expression RPAREN)))> */
		nil,
		/* 66 count <- <(COUNT LPAREN DISTINCT? (STAR / expression) RPAREN)> */
		nil,
		/* 67 groupConcat <- <(GROUPCONCAT LPAREN DISTINCT? expression (SEMICOLON SEPARATOR EQ string)? RPAREN)> */
		nil,
		/* 68 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		func() bool {
			position798, tokenIndex798, depth798 := position, tokenIndex, depth
			{
				position799 := position
				depth++
				{
					position800, tokenIndex800, depth800 := position, tokenIndex, depth
					{
						position802, tokenIndex802, depth802 := position, tokenIndex, depth
						{
							position804 := position
							depth++
							{
								position805, tokenIndex805, depth805 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l806
								}
								position++
								goto l805
							l806:
								position, tokenIndex, depth = position805, tokenIndex805, depth805
								if buffer[position] != rune('S') {
									goto l803
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
									goto l803
								}
								position++
							}
						l807:
							{
								position809, tokenIndex809, depth809 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l810
								}
								position++
								goto l809
							l810:
								position, tokenIndex, depth = position809, tokenIndex809, depth809
								if buffer[position] != rune('R') {
									goto l803
								}
								position++
							}
						l809:
							if !rules[ruleskip]() {
								goto l803
							}
							depth--
							add(ruleSTR, position804)
						}
						goto l802
					l803:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position812 := position
							depth++
							{
								position813, tokenIndex813, depth813 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l814
								}
								position++
								goto l813
							l814:
								position, tokenIndex, depth = position813, tokenIndex813, depth813
								if buffer[position] != rune('L') {
									goto l811
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
									goto l811
								}
								position++
							}
						l815:
							{
								position817, tokenIndex817, depth817 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l818
								}
								position++
								goto l817
							l818:
								position, tokenIndex, depth = position817, tokenIndex817, depth817
								if buffer[position] != rune('N') {
									goto l811
								}
								position++
							}
						l817:
							{
								position819, tokenIndex819, depth819 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l820
								}
								position++
								goto l819
							l820:
								position, tokenIndex, depth = position819, tokenIndex819, depth819
								if buffer[position] != rune('G') {
									goto l811
								}
								position++
							}
						l819:
							if !rules[ruleskip]() {
								goto l811
							}
							depth--
							add(ruleLANG, position812)
						}
						goto l802
					l811:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position822 := position
							depth++
							{
								position823, tokenIndex823, depth823 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l824
								}
								position++
								goto l823
							l824:
								position, tokenIndex, depth = position823, tokenIndex823, depth823
								if buffer[position] != rune('D') {
									goto l821
								}
								position++
							}
						l823:
							{
								position825, tokenIndex825, depth825 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l826
								}
								position++
								goto l825
							l826:
								position, tokenIndex, depth = position825, tokenIndex825, depth825
								if buffer[position] != rune('A') {
									goto l821
								}
								position++
							}
						l825:
							{
								position827, tokenIndex827, depth827 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l828
								}
								position++
								goto l827
							l828:
								position, tokenIndex, depth = position827, tokenIndex827, depth827
								if buffer[position] != rune('T') {
									goto l821
								}
								position++
							}
						l827:
							{
								position829, tokenIndex829, depth829 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l830
								}
								position++
								goto l829
							l830:
								position, tokenIndex, depth = position829, tokenIndex829, depth829
								if buffer[position] != rune('A') {
									goto l821
								}
								position++
							}
						l829:
							{
								position831, tokenIndex831, depth831 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l832
								}
								position++
								goto l831
							l832:
								position, tokenIndex, depth = position831, tokenIndex831, depth831
								if buffer[position] != rune('T') {
									goto l821
								}
								position++
							}
						l831:
							{
								position833, tokenIndex833, depth833 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l834
								}
								position++
								goto l833
							l834:
								position, tokenIndex, depth = position833, tokenIndex833, depth833
								if buffer[position] != rune('Y') {
									goto l821
								}
								position++
							}
						l833:
							{
								position835, tokenIndex835, depth835 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l836
								}
								position++
								goto l835
							l836:
								position, tokenIndex, depth = position835, tokenIndex835, depth835
								if buffer[position] != rune('P') {
									goto l821
								}
								position++
							}
						l835:
							{
								position837, tokenIndex837, depth837 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l838
								}
								position++
								goto l837
							l838:
								position, tokenIndex, depth = position837, tokenIndex837, depth837
								if buffer[position] != rune('E') {
									goto l821
								}
								position++
							}
						l837:
							if !rules[ruleskip]() {
								goto l821
							}
							depth--
							add(ruleDATATYPE, position822)
						}
						goto l802
					l821:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position840 := position
							depth++
							{
								position841, tokenIndex841, depth841 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l842
								}
								position++
								goto l841
							l842:
								position, tokenIndex, depth = position841, tokenIndex841, depth841
								if buffer[position] != rune('I') {
									goto l839
								}
								position++
							}
						l841:
							{
								position843, tokenIndex843, depth843 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l844
								}
								position++
								goto l843
							l844:
								position, tokenIndex, depth = position843, tokenIndex843, depth843
								if buffer[position] != rune('R') {
									goto l839
								}
								position++
							}
						l843:
							{
								position845, tokenIndex845, depth845 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l846
								}
								position++
								goto l845
							l846:
								position, tokenIndex, depth = position845, tokenIndex845, depth845
								if buffer[position] != rune('I') {
									goto l839
								}
								position++
							}
						l845:
							if !rules[ruleskip]() {
								goto l839
							}
							depth--
							add(ruleIRI, position840)
						}
						goto l802
					l839:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position848 := position
							depth++
							{
								position849, tokenIndex849, depth849 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l850
								}
								position++
								goto l849
							l850:
								position, tokenIndex, depth = position849, tokenIndex849, depth849
								if buffer[position] != rune('U') {
									goto l847
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
									goto l847
								}
								position++
							}
						l851:
							{
								position853, tokenIndex853, depth853 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l854
								}
								position++
								goto l853
							l854:
								position, tokenIndex, depth = position853, tokenIndex853, depth853
								if buffer[position] != rune('I') {
									goto l847
								}
								position++
							}
						l853:
							if !rules[ruleskip]() {
								goto l847
							}
							depth--
							add(ruleURI, position848)
						}
						goto l802
					l847:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position856 := position
							depth++
							{
								position857, tokenIndex857, depth857 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l858
								}
								position++
								goto l857
							l858:
								position, tokenIndex, depth = position857, tokenIndex857, depth857
								if buffer[position] != rune('S') {
									goto l855
								}
								position++
							}
						l857:
							{
								position859, tokenIndex859, depth859 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l860
								}
								position++
								goto l859
							l860:
								position, tokenIndex, depth = position859, tokenIndex859, depth859
								if buffer[position] != rune('T') {
									goto l855
								}
								position++
							}
						l859:
							{
								position861, tokenIndex861, depth861 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l862
								}
								position++
								goto l861
							l862:
								position, tokenIndex, depth = position861, tokenIndex861, depth861
								if buffer[position] != rune('R') {
									goto l855
								}
								position++
							}
						l861:
							{
								position863, tokenIndex863, depth863 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l864
								}
								position++
								goto l863
							l864:
								position, tokenIndex, depth = position863, tokenIndex863, depth863
								if buffer[position] != rune('L') {
									goto l855
								}
								position++
							}
						l863:
							{
								position865, tokenIndex865, depth865 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l866
								}
								position++
								goto l865
							l866:
								position, tokenIndex, depth = position865, tokenIndex865, depth865
								if buffer[position] != rune('E') {
									goto l855
								}
								position++
							}
						l865:
							{
								position867, tokenIndex867, depth867 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l868
								}
								position++
								goto l867
							l868:
								position, tokenIndex, depth = position867, tokenIndex867, depth867
								if buffer[position] != rune('N') {
									goto l855
								}
								position++
							}
						l867:
							if !rules[ruleskip]() {
								goto l855
							}
							depth--
							add(ruleSTRLEN, position856)
						}
						goto l802
					l855:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position870 := position
							depth++
							{
								position871, tokenIndex871, depth871 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l872
								}
								position++
								goto l871
							l872:
								position, tokenIndex, depth = position871, tokenIndex871, depth871
								if buffer[position] != rune('M') {
									goto l869
								}
								position++
							}
						l871:
							{
								position873, tokenIndex873, depth873 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l874
								}
								position++
								goto l873
							l874:
								position, tokenIndex, depth = position873, tokenIndex873, depth873
								if buffer[position] != rune('O') {
									goto l869
								}
								position++
							}
						l873:
							{
								position875, tokenIndex875, depth875 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l876
								}
								position++
								goto l875
							l876:
								position, tokenIndex, depth = position875, tokenIndex875, depth875
								if buffer[position] != rune('N') {
									goto l869
								}
								position++
							}
						l875:
							{
								position877, tokenIndex877, depth877 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l878
								}
								position++
								goto l877
							l878:
								position, tokenIndex, depth = position877, tokenIndex877, depth877
								if buffer[position] != rune('T') {
									goto l869
								}
								position++
							}
						l877:
							{
								position879, tokenIndex879, depth879 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l880
								}
								position++
								goto l879
							l880:
								position, tokenIndex, depth = position879, tokenIndex879, depth879
								if buffer[position] != rune('H') {
									goto l869
								}
								position++
							}
						l879:
							if !rules[ruleskip]() {
								goto l869
							}
							depth--
							add(ruleMONTH, position870)
						}
						goto l802
					l869:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position882 := position
							depth++
							{
								position883, tokenIndex883, depth883 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l884
								}
								position++
								goto l883
							l884:
								position, tokenIndex, depth = position883, tokenIndex883, depth883
								if buffer[position] != rune('M') {
									goto l881
								}
								position++
							}
						l883:
							{
								position885, tokenIndex885, depth885 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l886
								}
								position++
								goto l885
							l886:
								position, tokenIndex, depth = position885, tokenIndex885, depth885
								if buffer[position] != rune('I') {
									goto l881
								}
								position++
							}
						l885:
							{
								position887, tokenIndex887, depth887 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l888
								}
								position++
								goto l887
							l888:
								position, tokenIndex, depth = position887, tokenIndex887, depth887
								if buffer[position] != rune('N') {
									goto l881
								}
								position++
							}
						l887:
							{
								position889, tokenIndex889, depth889 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l890
								}
								position++
								goto l889
							l890:
								position, tokenIndex, depth = position889, tokenIndex889, depth889
								if buffer[position] != rune('U') {
									goto l881
								}
								position++
							}
						l889:
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
									goto l881
								}
								position++
							}
						l891:
							{
								position893, tokenIndex893, depth893 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l894
								}
								position++
								goto l893
							l894:
								position, tokenIndex, depth = position893, tokenIndex893, depth893
								if buffer[position] != rune('E') {
									goto l881
								}
								position++
							}
						l893:
							{
								position895, tokenIndex895, depth895 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l896
								}
								position++
								goto l895
							l896:
								position, tokenIndex, depth = position895, tokenIndex895, depth895
								if buffer[position] != rune('S') {
									goto l881
								}
								position++
							}
						l895:
							if !rules[ruleskip]() {
								goto l881
							}
							depth--
							add(ruleMINUTES, position882)
						}
						goto l802
					l881:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position898 := position
							depth++
							{
								position899, tokenIndex899, depth899 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l900
								}
								position++
								goto l899
							l900:
								position, tokenIndex, depth = position899, tokenIndex899, depth899
								if buffer[position] != rune('S') {
									goto l897
								}
								position++
							}
						l899:
							{
								position901, tokenIndex901, depth901 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l902
								}
								position++
								goto l901
							l902:
								position, tokenIndex, depth = position901, tokenIndex901, depth901
								if buffer[position] != rune('E') {
									goto l897
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
									goto l897
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
									goto l897
								}
								position++
							}
						l905:
							{
								position907, tokenIndex907, depth907 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l908
								}
								position++
								goto l907
							l908:
								position, tokenIndex, depth = position907, tokenIndex907, depth907
								if buffer[position] != rune('N') {
									goto l897
								}
								position++
							}
						l907:
							{
								position909, tokenIndex909, depth909 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l910
								}
								position++
								goto l909
							l910:
								position, tokenIndex, depth = position909, tokenIndex909, depth909
								if buffer[position] != rune('D') {
									goto l897
								}
								position++
							}
						l909:
							{
								position911, tokenIndex911, depth911 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l912
								}
								position++
								goto l911
							l912:
								position, tokenIndex, depth = position911, tokenIndex911, depth911
								if buffer[position] != rune('S') {
									goto l897
								}
								position++
							}
						l911:
							if !rules[ruleskip]() {
								goto l897
							}
							depth--
							add(ruleSECONDS, position898)
						}
						goto l802
					l897:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position914 := position
							depth++
							{
								position915, tokenIndex915, depth915 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l916
								}
								position++
								goto l915
							l916:
								position, tokenIndex, depth = position915, tokenIndex915, depth915
								if buffer[position] != rune('T') {
									goto l913
								}
								position++
							}
						l915:
							{
								position917, tokenIndex917, depth917 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l918
								}
								position++
								goto l917
							l918:
								position, tokenIndex, depth = position917, tokenIndex917, depth917
								if buffer[position] != rune('I') {
									goto l913
								}
								position++
							}
						l917:
							{
								position919, tokenIndex919, depth919 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l920
								}
								position++
								goto l919
							l920:
								position, tokenIndex, depth = position919, tokenIndex919, depth919
								if buffer[position] != rune('M') {
									goto l913
								}
								position++
							}
						l919:
							{
								position921, tokenIndex921, depth921 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l922
								}
								position++
								goto l921
							l922:
								position, tokenIndex, depth = position921, tokenIndex921, depth921
								if buffer[position] != rune('E') {
									goto l913
								}
								position++
							}
						l921:
							{
								position923, tokenIndex923, depth923 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l924
								}
								position++
								goto l923
							l924:
								position, tokenIndex, depth = position923, tokenIndex923, depth923
								if buffer[position] != rune('Z') {
									goto l913
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
									goto l913
								}
								position++
							}
						l925:
							{
								position927, tokenIndex927, depth927 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l928
								}
								position++
								goto l927
							l928:
								position, tokenIndex, depth = position927, tokenIndex927, depth927
								if buffer[position] != rune('N') {
									goto l913
								}
								position++
							}
						l927:
							{
								position929, tokenIndex929, depth929 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l930
								}
								position++
								goto l929
							l930:
								position, tokenIndex, depth = position929, tokenIndex929, depth929
								if buffer[position] != rune('E') {
									goto l913
								}
								position++
							}
						l929:
							if !rules[ruleskip]() {
								goto l913
							}
							depth--
							add(ruleTIMEZONE, position914)
						}
						goto l802
					l913:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position932 := position
							depth++
							{
								position933, tokenIndex933, depth933 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l934
								}
								position++
								goto l933
							l934:
								position, tokenIndex, depth = position933, tokenIndex933, depth933
								if buffer[position] != rune('S') {
									goto l931
								}
								position++
							}
						l933:
							{
								position935, tokenIndex935, depth935 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l936
								}
								position++
								goto l935
							l936:
								position, tokenIndex, depth = position935, tokenIndex935, depth935
								if buffer[position] != rune('H') {
									goto l931
								}
								position++
							}
						l935:
							{
								position937, tokenIndex937, depth937 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l938
								}
								position++
								goto l937
							l938:
								position, tokenIndex, depth = position937, tokenIndex937, depth937
								if buffer[position] != rune('A') {
									goto l931
								}
								position++
							}
						l937:
							if buffer[position] != rune('1') {
								goto l931
							}
							position++
							if !rules[ruleskip]() {
								goto l931
							}
							depth--
							add(ruleSHA1, position932)
						}
						goto l802
					l931:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position940 := position
							depth++
							{
								position941, tokenIndex941, depth941 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l942
								}
								position++
								goto l941
							l942:
								position, tokenIndex, depth = position941, tokenIndex941, depth941
								if buffer[position] != rune('S') {
									goto l939
								}
								position++
							}
						l941:
							{
								position943, tokenIndex943, depth943 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l944
								}
								position++
								goto l943
							l944:
								position, tokenIndex, depth = position943, tokenIndex943, depth943
								if buffer[position] != rune('H') {
									goto l939
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
									goto l939
								}
								position++
							}
						l945:
							if buffer[position] != rune('2') {
								goto l939
							}
							position++
							if buffer[position] != rune('5') {
								goto l939
							}
							position++
							if buffer[position] != rune('6') {
								goto l939
							}
							position++
							if !rules[ruleskip]() {
								goto l939
							}
							depth--
							add(ruleSHA256, position940)
						}
						goto l802
					l939:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position948 := position
							depth++
							{
								position949, tokenIndex949, depth949 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l950
								}
								position++
								goto l949
							l950:
								position, tokenIndex, depth = position949, tokenIndex949, depth949
								if buffer[position] != rune('S') {
									goto l947
								}
								position++
							}
						l949:
							{
								position951, tokenIndex951, depth951 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l952
								}
								position++
								goto l951
							l952:
								position, tokenIndex, depth = position951, tokenIndex951, depth951
								if buffer[position] != rune('H') {
									goto l947
								}
								position++
							}
						l951:
							{
								position953, tokenIndex953, depth953 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l954
								}
								position++
								goto l953
							l954:
								position, tokenIndex, depth = position953, tokenIndex953, depth953
								if buffer[position] != rune('A') {
									goto l947
								}
								position++
							}
						l953:
							if buffer[position] != rune('3') {
								goto l947
							}
							position++
							if buffer[position] != rune('8') {
								goto l947
							}
							position++
							if buffer[position] != rune('4') {
								goto l947
							}
							position++
							if !rules[ruleskip]() {
								goto l947
							}
							depth--
							add(ruleSHA384, position948)
						}
						goto l802
					l947:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
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
								if buffer[position] != rune('i') {
									goto l962
								}
								position++
								goto l961
							l962:
								position, tokenIndex, depth = position961, tokenIndex961, depth961
								if buffer[position] != rune('I') {
									goto l955
								}
								position++
							}
						l961:
							{
								position963, tokenIndex963, depth963 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l964
								}
								position++
								goto l963
							l964:
								position, tokenIndex, depth = position963, tokenIndex963, depth963
								if buffer[position] != rune('R') {
									goto l955
								}
								position++
							}
						l963:
							{
								position965, tokenIndex965, depth965 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l966
								}
								position++
								goto l965
							l966:
								position, tokenIndex, depth = position965, tokenIndex965, depth965
								if buffer[position] != rune('I') {
									goto l955
								}
								position++
							}
						l965:
							if !rules[ruleskip]() {
								goto l955
							}
							depth--
							add(ruleISIRI, position956)
						}
						goto l802
					l955:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position968 := position
							depth++
							{
								position969, tokenIndex969, depth969 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l970
								}
								position++
								goto l969
							l970:
								position, tokenIndex, depth = position969, tokenIndex969, depth969
								if buffer[position] != rune('I') {
									goto l967
								}
								position++
							}
						l969:
							{
								position971, tokenIndex971, depth971 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l972
								}
								position++
								goto l971
							l972:
								position, tokenIndex, depth = position971, tokenIndex971, depth971
								if buffer[position] != rune('S') {
									goto l967
								}
								position++
							}
						l971:
							{
								position973, tokenIndex973, depth973 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l974
								}
								position++
								goto l973
							l974:
								position, tokenIndex, depth = position973, tokenIndex973, depth973
								if buffer[position] != rune('U') {
									goto l967
								}
								position++
							}
						l973:
							{
								position975, tokenIndex975, depth975 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l976
								}
								position++
								goto l975
							l976:
								position, tokenIndex, depth = position975, tokenIndex975, depth975
								if buffer[position] != rune('R') {
									goto l967
								}
								position++
							}
						l975:
							{
								position977, tokenIndex977, depth977 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l978
								}
								position++
								goto l977
							l978:
								position, tokenIndex, depth = position977, tokenIndex977, depth977
								if buffer[position] != rune('I') {
									goto l967
								}
								position++
							}
						l977:
							if !rules[ruleskip]() {
								goto l967
							}
							depth--
							add(ruleISURI, position968)
						}
						goto l802
					l967:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position980 := position
							depth++
							{
								position981, tokenIndex981, depth981 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l982
								}
								position++
								goto l981
							l982:
								position, tokenIndex, depth = position981, tokenIndex981, depth981
								if buffer[position] != rune('I') {
									goto l979
								}
								position++
							}
						l981:
							{
								position983, tokenIndex983, depth983 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l984
								}
								position++
								goto l983
							l984:
								position, tokenIndex, depth = position983, tokenIndex983, depth983
								if buffer[position] != rune('S') {
									goto l979
								}
								position++
							}
						l983:
							{
								position985, tokenIndex985, depth985 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l986
								}
								position++
								goto l985
							l986:
								position, tokenIndex, depth = position985, tokenIndex985, depth985
								if buffer[position] != rune('B') {
									goto l979
								}
								position++
							}
						l985:
							{
								position987, tokenIndex987, depth987 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l988
								}
								position++
								goto l987
							l988:
								position, tokenIndex, depth = position987, tokenIndex987, depth987
								if buffer[position] != rune('L') {
									goto l979
								}
								position++
							}
						l987:
							{
								position989, tokenIndex989, depth989 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l990
								}
								position++
								goto l989
							l990:
								position, tokenIndex, depth = position989, tokenIndex989, depth989
								if buffer[position] != rune('A') {
									goto l979
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
									goto l979
								}
								position++
							}
						l991:
							{
								position993, tokenIndex993, depth993 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l994
								}
								position++
								goto l993
							l994:
								position, tokenIndex, depth = position993, tokenIndex993, depth993
								if buffer[position] != rune('K') {
									goto l979
								}
								position++
							}
						l993:
							if !rules[ruleskip]() {
								goto l979
							}
							depth--
							add(ruleISBLANK, position980)
						}
						goto l802
					l979:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							position996 := position
							depth++
							{
								position997, tokenIndex997, depth997 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l998
								}
								position++
								goto l997
							l998:
								position, tokenIndex, depth = position997, tokenIndex997, depth997
								if buffer[position] != rune('I') {
									goto l995
								}
								position++
							}
						l997:
							{
								position999, tokenIndex999, depth999 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1000
								}
								position++
								goto l999
							l1000:
								position, tokenIndex, depth = position999, tokenIndex999, depth999
								if buffer[position] != rune('S') {
									goto l995
								}
								position++
							}
						l999:
							{
								position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1002
								}
								position++
								goto l1001
							l1002:
								position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
								if buffer[position] != rune('L') {
									goto l995
								}
								position++
							}
						l1001:
							{
								position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1004
								}
								position++
								goto l1003
							l1004:
								position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
								if buffer[position] != rune('I') {
									goto l995
								}
								position++
							}
						l1003:
							{
								position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1006
								}
								position++
								goto l1005
							l1006:
								position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
								if buffer[position] != rune('T') {
									goto l995
								}
								position++
							}
						l1005:
							{
								position1007, tokenIndex1007, depth1007 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1008
								}
								position++
								goto l1007
							l1008:
								position, tokenIndex, depth = position1007, tokenIndex1007, depth1007
								if buffer[position] != rune('E') {
									goto l995
								}
								position++
							}
						l1007:
							{
								position1009, tokenIndex1009, depth1009 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1010
								}
								position++
								goto l1009
							l1010:
								position, tokenIndex, depth = position1009, tokenIndex1009, depth1009
								if buffer[position] != rune('R') {
									goto l995
								}
								position++
							}
						l1009:
							{
								position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1012
								}
								position++
								goto l1011
							l1012:
								position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
								if buffer[position] != rune('A') {
									goto l995
								}
								position++
							}
						l1011:
							{
								position1013, tokenIndex1013, depth1013 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1014
								}
								position++
								goto l1013
							l1014:
								position, tokenIndex, depth = position1013, tokenIndex1013, depth1013
								if buffer[position] != rune('L') {
									goto l995
								}
								position++
							}
						l1013:
							if !rules[ruleskip]() {
								goto l995
							}
							depth--
							add(ruleISLITERAL, position996)
						}
						goto l802
					l995:
						position, tokenIndex, depth = position802, tokenIndex802, depth802
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position1016 := position
									depth++
									{
										position1017, tokenIndex1017, depth1017 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1018
										}
										position++
										goto l1017
									l1018:
										position, tokenIndex, depth = position1017, tokenIndex1017, depth1017
										if buffer[position] != rune('I') {
											goto l801
										}
										position++
									}
								l1017:
									{
										position1019, tokenIndex1019, depth1019 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1020
										}
										position++
										goto l1019
									l1020:
										position, tokenIndex, depth = position1019, tokenIndex1019, depth1019
										if buffer[position] != rune('S') {
											goto l801
										}
										position++
									}
								l1019:
									{
										position1021, tokenIndex1021, depth1021 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1022
										}
										position++
										goto l1021
									l1022:
										position, tokenIndex, depth = position1021, tokenIndex1021, depth1021
										if buffer[position] != rune('N') {
											goto l801
										}
										position++
									}
								l1021:
									{
										position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1024
										}
										position++
										goto l1023
									l1024:
										position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
										if buffer[position] != rune('U') {
											goto l801
										}
										position++
									}
								l1023:
									{
										position1025, tokenIndex1025, depth1025 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1026
										}
										position++
										goto l1025
									l1026:
										position, tokenIndex, depth = position1025, tokenIndex1025, depth1025
										if buffer[position] != rune('M') {
											goto l801
										}
										position++
									}
								l1025:
									{
										position1027, tokenIndex1027, depth1027 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1028
										}
										position++
										goto l1027
									l1028:
										position, tokenIndex, depth = position1027, tokenIndex1027, depth1027
										if buffer[position] != rune('E') {
											goto l801
										}
										position++
									}
								l1027:
									{
										position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1030
										}
										position++
										goto l1029
									l1030:
										position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
										if buffer[position] != rune('R') {
											goto l801
										}
										position++
									}
								l1029:
									{
										position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1032
										}
										position++
										goto l1031
									l1032:
										position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
										if buffer[position] != rune('I') {
											goto l801
										}
										position++
									}
								l1031:
									{
										position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1034
										}
										position++
										goto l1033
									l1034:
										position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
										if buffer[position] != rune('C') {
											goto l801
										}
										position++
									}
								l1033:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleISNUMERIC, position1016)
								}
								break
							case 'S', 's':
								{
									position1035 := position
									depth++
									{
										position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1037
										}
										position++
										goto l1036
									l1037:
										position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
										if buffer[position] != rune('S') {
											goto l801
										}
										position++
									}
								l1036:
									{
										position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1039
										}
										position++
										goto l1038
									l1039:
										position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
										if buffer[position] != rune('H') {
											goto l801
										}
										position++
									}
								l1038:
									{
										position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1041
										}
										position++
										goto l1040
									l1041:
										position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
										if buffer[position] != rune('A') {
											goto l801
										}
										position++
									}
								l1040:
									if buffer[position] != rune('5') {
										goto l801
									}
									position++
									if buffer[position] != rune('1') {
										goto l801
									}
									position++
									if buffer[position] != rune('2') {
										goto l801
									}
									position++
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleSHA512, position1035)
								}
								break
							case 'M', 'm':
								{
									position1042 := position
									depth++
									{
										position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1044
										}
										position++
										goto l1043
									l1044:
										position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
										if buffer[position] != rune('M') {
											goto l801
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
											goto l801
										}
										position++
									}
								l1045:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleMD5, position1042)
								}
								break
							case 'T', 't':
								{
									position1047 := position
									depth++
									{
										position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1049
										}
										position++
										goto l1048
									l1049:
										position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
										if buffer[position] != rune('T') {
											goto l801
										}
										position++
									}
								l1048:
									{
										position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1051
										}
										position++
										goto l1050
									l1051:
										position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
										if buffer[position] != rune('Z') {
											goto l801
										}
										position++
									}
								l1050:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleTZ, position1047)
								}
								break
							case 'H', 'h':
								{
									position1052 := position
									depth++
									{
										position1053, tokenIndex1053, depth1053 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1054
										}
										position++
										goto l1053
									l1054:
										position, tokenIndex, depth = position1053, tokenIndex1053, depth1053
										if buffer[position] != rune('H') {
											goto l801
										}
										position++
									}
								l1053:
									{
										position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1056
										}
										position++
										goto l1055
									l1056:
										position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
										if buffer[position] != rune('O') {
											goto l801
										}
										position++
									}
								l1055:
									{
										position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1058
										}
										position++
										goto l1057
									l1058:
										position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
										if buffer[position] != rune('U') {
											goto l801
										}
										position++
									}
								l1057:
									{
										position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1060
										}
										position++
										goto l1059
									l1060:
										position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
										if buffer[position] != rune('R') {
											goto l801
										}
										position++
									}
								l1059:
									{
										position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1062
										}
										position++
										goto l1061
									l1062:
										position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
										if buffer[position] != rune('S') {
											goto l801
										}
										position++
									}
								l1061:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleHOURS, position1052)
								}
								break
							case 'D', 'd':
								{
									position1063 := position
									depth++
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
											goto l801
										}
										position++
									}
								l1064:
									{
										position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1067
										}
										position++
										goto l1066
									l1067:
										position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
										if buffer[position] != rune('A') {
											goto l801
										}
										position++
									}
								l1066:
									{
										position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1069
										}
										position++
										goto l1068
									l1069:
										position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
										if buffer[position] != rune('Y') {
											goto l801
										}
										position++
									}
								l1068:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleDAY, position1063)
								}
								break
							case 'Y', 'y':
								{
									position1070 := position
									depth++
									{
										position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1072
										}
										position++
										goto l1071
									l1072:
										position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
										if buffer[position] != rune('Y') {
											goto l801
										}
										position++
									}
								l1071:
									{
										position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1074
										}
										position++
										goto l1073
									l1074:
										position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
										if buffer[position] != rune('E') {
											goto l801
										}
										position++
									}
								l1073:
									{
										position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1076
										}
										position++
										goto l1075
									l1076:
										position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
										if buffer[position] != rune('A') {
											goto l801
										}
										position++
									}
								l1075:
									{
										position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1078
										}
										position++
										goto l1077
									l1078:
										position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
										if buffer[position] != rune('R') {
											goto l801
										}
										position++
									}
								l1077:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleYEAR, position1070)
								}
								break
							case 'E', 'e':
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
											goto l801
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
											goto l801
										}
										position++
									}
								l1082:
									{
										position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1085
										}
										position++
										goto l1084
									l1085:
										position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
										if buffer[position] != rune('C') {
											goto l801
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
											goto l801
										}
										position++
									}
								l1086:
									{
										position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1089
										}
										position++
										goto l1088
									l1089:
										position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
										if buffer[position] != rune('D') {
											goto l801
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
											goto l801
										}
										position++
									}
								l1090:
									if buffer[position] != rune('_') {
										goto l801
									}
									position++
									{
										position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1093
										}
										position++
										goto l1092
									l1093:
										position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
										if buffer[position] != rune('F') {
											goto l801
										}
										position++
									}
								l1092:
									{
										position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1095
										}
										position++
										goto l1094
									l1095:
										position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
										if buffer[position] != rune('O') {
											goto l801
										}
										position++
									}
								l1094:
									{
										position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1097
										}
										position++
										goto l1096
									l1097:
										position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
										if buffer[position] != rune('R') {
											goto l801
										}
										position++
									}
								l1096:
									if buffer[position] != rune('_') {
										goto l801
									}
									position++
									{
										position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1099
										}
										position++
										goto l1098
									l1099:
										position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
										if buffer[position] != rune('U') {
											goto l801
										}
										position++
									}
								l1098:
									{
										position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1101
										}
										position++
										goto l1100
									l1101:
										position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
										if buffer[position] != rune('R') {
											goto l801
										}
										position++
									}
								l1100:
									{
										position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1103
										}
										position++
										goto l1102
									l1103:
										position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
										if buffer[position] != rune('I') {
											goto l801
										}
										position++
									}
								l1102:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleENCODEFORURI, position1079)
								}
								break
							case 'L', 'l':
								{
									position1104 := position
									depth++
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
											goto l801
										}
										position++
									}
								l1105:
									{
										position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1108
										}
										position++
										goto l1107
									l1108:
										position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
										if buffer[position] != rune('C') {
											goto l801
										}
										position++
									}
								l1107:
									{
										position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1110
										}
										position++
										goto l1109
									l1110:
										position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
										if buffer[position] != rune('A') {
											goto l801
										}
										position++
									}
								l1109:
									{
										position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1112
										}
										position++
										goto l1111
									l1112:
										position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
										if buffer[position] != rune('S') {
											goto l801
										}
										position++
									}
								l1111:
									{
										position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1114
										}
										position++
										goto l1113
									l1114:
										position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
										if buffer[position] != rune('E') {
											goto l801
										}
										position++
									}
								l1113:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleLCASE, position1104)
								}
								break
							case 'U', 'u':
								{
									position1115 := position
									depth++
									{
										position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1117
										}
										position++
										goto l1116
									l1117:
										position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
										if buffer[position] != rune('U') {
											goto l801
										}
										position++
									}
								l1116:
									{
										position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1119
										}
										position++
										goto l1118
									l1119:
										position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
										if buffer[position] != rune('C') {
											goto l801
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
											goto l801
										}
										position++
									}
								l1120:
									{
										position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1123
										}
										position++
										goto l1122
									l1123:
										position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
										if buffer[position] != rune('S') {
											goto l801
										}
										position++
									}
								l1122:
									{
										position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1125
										}
										position++
										goto l1124
									l1125:
										position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
										if buffer[position] != rune('E') {
											goto l801
										}
										position++
									}
								l1124:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleUCASE, position1115)
								}
								break
							case 'F', 'f':
								{
									position1126 := position
									depth++
									{
										position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1128
										}
										position++
										goto l1127
									l1128:
										position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
										if buffer[position] != rune('F') {
											goto l801
										}
										position++
									}
								l1127:
									{
										position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1130
										}
										position++
										goto l1129
									l1130:
										position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
										if buffer[position] != rune('L') {
											goto l801
										}
										position++
									}
								l1129:
									{
										position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1132
										}
										position++
										goto l1131
									l1132:
										position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
										if buffer[position] != rune('O') {
											goto l801
										}
										position++
									}
								l1131:
									{
										position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1134
										}
										position++
										goto l1133
									l1134:
										position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
										if buffer[position] != rune('O') {
											goto l801
										}
										position++
									}
								l1133:
									{
										position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1136
										}
										position++
										goto l1135
									l1136:
										position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
										if buffer[position] != rune('R') {
											goto l801
										}
										position++
									}
								l1135:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleFLOOR, position1126)
								}
								break
							case 'R', 'r':
								{
									position1137 := position
									depth++
									{
										position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1139
										}
										position++
										goto l1138
									l1139:
										position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
										if buffer[position] != rune('R') {
											goto l801
										}
										position++
									}
								l1138:
									{
										position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1141
										}
										position++
										goto l1140
									l1141:
										position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
										if buffer[position] != rune('O') {
											goto l801
										}
										position++
									}
								l1140:
									{
										position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1143
										}
										position++
										goto l1142
									l1143:
										position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
										if buffer[position] != rune('U') {
											goto l801
										}
										position++
									}
								l1142:
									{
										position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1145
										}
										position++
										goto l1144
									l1145:
										position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
										if buffer[position] != rune('N') {
											goto l801
										}
										position++
									}
								l1144:
									{
										position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1147
										}
										position++
										goto l1146
									l1147:
										position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
										if buffer[position] != rune('D') {
											goto l801
										}
										position++
									}
								l1146:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleROUND, position1137)
								}
								break
							case 'C', 'c':
								{
									position1148 := position
									depth++
									{
										position1149, tokenIndex1149, depth1149 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1150
										}
										position++
										goto l1149
									l1150:
										position, tokenIndex, depth = position1149, tokenIndex1149, depth1149
										if buffer[position] != rune('C') {
											goto l801
										}
										position++
									}
								l1149:
									{
										position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1152
										}
										position++
										goto l1151
									l1152:
										position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
										if buffer[position] != rune('E') {
											goto l801
										}
										position++
									}
								l1151:
									{
										position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1154
										}
										position++
										goto l1153
									l1154:
										position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
										if buffer[position] != rune('I') {
											goto l801
										}
										position++
									}
								l1153:
									{
										position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1156
										}
										position++
										goto l1155
									l1156:
										position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
										if buffer[position] != rune('L') {
											goto l801
										}
										position++
									}
								l1155:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleCEIL, position1148)
								}
								break
							default:
								{
									position1157 := position
									depth++
									{
										position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1159
										}
										position++
										goto l1158
									l1159:
										position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
										if buffer[position] != rune('A') {
											goto l801
										}
										position++
									}
								l1158:
									{
										position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1161
										}
										position++
										goto l1160
									l1161:
										position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
										if buffer[position] != rune('B') {
											goto l801
										}
										position++
									}
								l1160:
									{
										position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1163
										}
										position++
										goto l1162
									l1163:
										position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
										if buffer[position] != rune('S') {
											goto l801
										}
										position++
									}
								l1162:
									if !rules[ruleskip]() {
										goto l801
									}
									depth--
									add(ruleABS, position1157)
								}
								break
							}
						}

					}
				l802:
					if !rules[ruleLPAREN]() {
						goto l801
					}
					if !rules[ruleexpression]() {
						goto l801
					}
					if !rules[ruleRPAREN]() {
						goto l801
					}
					goto l800
				l801:
					position, tokenIndex, depth = position800, tokenIndex800, depth800
					{
						position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
						{
							position1167 := position
							depth++
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
									goto l1166
								}
								position++
							}
						l1168:
							{
								position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1171
								}
								position++
								goto l1170
							l1171:
								position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
								if buffer[position] != rune('T') {
									goto l1166
								}
								position++
							}
						l1170:
							{
								position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1173
								}
								position++
								goto l1172
							l1173:
								position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
								if buffer[position] != rune('R') {
									goto l1166
								}
								position++
							}
						l1172:
							{
								position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1175
								}
								position++
								goto l1174
							l1175:
								position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
								if buffer[position] != rune('S') {
									goto l1166
								}
								position++
							}
						l1174:
							{
								position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1177
								}
								position++
								goto l1176
							l1177:
								position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
								if buffer[position] != rune('T') {
									goto l1166
								}
								position++
							}
						l1176:
							{
								position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1179
								}
								position++
								goto l1178
							l1179:
								position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
								if buffer[position] != rune('A') {
									goto l1166
								}
								position++
							}
						l1178:
							{
								position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1181
								}
								position++
								goto l1180
							l1181:
								position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
								if buffer[position] != rune('R') {
									goto l1166
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
									goto l1166
								}
								position++
							}
						l1182:
							{
								position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1185
								}
								position++
								goto l1184
							l1185:
								position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
								if buffer[position] != rune('S') {
									goto l1166
								}
								position++
							}
						l1184:
							if !rules[ruleskip]() {
								goto l1166
							}
							depth--
							add(ruleSTRSTARTS, position1167)
						}
						goto l1165
					l1166:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
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
									goto l1186
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
									goto l1186
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
									goto l1186
								}
								position++
							}
						l1192:
							{
								position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1195
								}
								position++
								goto l1194
							l1195:
								position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
								if buffer[position] != rune('E') {
									goto l1186
								}
								position++
							}
						l1194:
							{
								position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1197
								}
								position++
								goto l1196
							l1197:
								position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
								if buffer[position] != rune('N') {
									goto l1186
								}
								position++
							}
						l1196:
							{
								position1198, tokenIndex1198, depth1198 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1199
								}
								position++
								goto l1198
							l1199:
								position, tokenIndex, depth = position1198, tokenIndex1198, depth1198
								if buffer[position] != rune('D') {
									goto l1186
								}
								position++
							}
						l1198:
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
									goto l1186
								}
								position++
							}
						l1200:
							if !rules[ruleskip]() {
								goto l1186
							}
							depth--
							add(ruleSTRENDS, position1187)
						}
						goto l1165
					l1186:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						{
							position1203 := position
							depth++
							{
								position1204, tokenIndex1204, depth1204 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1205
								}
								position++
								goto l1204
							l1205:
								position, tokenIndex, depth = position1204, tokenIndex1204, depth1204
								if buffer[position] != rune('S') {
									goto l1202
								}
								position++
							}
						l1204:
							{
								position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1207
								}
								position++
								goto l1206
							l1207:
								position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
								if buffer[position] != rune('T') {
									goto l1202
								}
								position++
							}
						l1206:
							{
								position1208, tokenIndex1208, depth1208 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1209
								}
								position++
								goto l1208
							l1209:
								position, tokenIndex, depth = position1208, tokenIndex1208, depth1208
								if buffer[position] != rune('R') {
									goto l1202
								}
								position++
							}
						l1208:
							{
								position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1211
								}
								position++
								goto l1210
							l1211:
								position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
								if buffer[position] != rune('B') {
									goto l1202
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
									goto l1202
								}
								position++
							}
						l1212:
							{
								position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1215
								}
								position++
								goto l1214
							l1215:
								position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
								if buffer[position] != rune('F') {
									goto l1202
								}
								position++
							}
						l1214:
							{
								position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1217
								}
								position++
								goto l1216
							l1217:
								position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
								if buffer[position] != rune('O') {
									goto l1202
								}
								position++
							}
						l1216:
							{
								position1218, tokenIndex1218, depth1218 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1219
								}
								position++
								goto l1218
							l1219:
								position, tokenIndex, depth = position1218, tokenIndex1218, depth1218
								if buffer[position] != rune('R') {
									goto l1202
								}
								position++
							}
						l1218:
							{
								position1220, tokenIndex1220, depth1220 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1221
								}
								position++
								goto l1220
							l1221:
								position, tokenIndex, depth = position1220, tokenIndex1220, depth1220
								if buffer[position] != rune('E') {
									goto l1202
								}
								position++
							}
						l1220:
							if !rules[ruleskip]() {
								goto l1202
							}
							depth--
							add(ruleSTRBEFORE, position1203)
						}
						goto l1165
					l1202:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						{
							position1223 := position
							depth++
							{
								position1224, tokenIndex1224, depth1224 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1225
								}
								position++
								goto l1224
							l1225:
								position, tokenIndex, depth = position1224, tokenIndex1224, depth1224
								if buffer[position] != rune('S') {
									goto l1222
								}
								position++
							}
						l1224:
							{
								position1226, tokenIndex1226, depth1226 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1227
								}
								position++
								goto l1226
							l1227:
								position, tokenIndex, depth = position1226, tokenIndex1226, depth1226
								if buffer[position] != rune('T') {
									goto l1222
								}
								position++
							}
						l1226:
							{
								position1228, tokenIndex1228, depth1228 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1229
								}
								position++
								goto l1228
							l1229:
								position, tokenIndex, depth = position1228, tokenIndex1228, depth1228
								if buffer[position] != rune('R') {
									goto l1222
								}
								position++
							}
						l1228:
							{
								position1230, tokenIndex1230, depth1230 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1231
								}
								position++
								goto l1230
							l1231:
								position, tokenIndex, depth = position1230, tokenIndex1230, depth1230
								if buffer[position] != rune('A') {
									goto l1222
								}
								position++
							}
						l1230:
							{
								position1232, tokenIndex1232, depth1232 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1233
								}
								position++
								goto l1232
							l1233:
								position, tokenIndex, depth = position1232, tokenIndex1232, depth1232
								if buffer[position] != rune('F') {
									goto l1222
								}
								position++
							}
						l1232:
							{
								position1234, tokenIndex1234, depth1234 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1235
								}
								position++
								goto l1234
							l1235:
								position, tokenIndex, depth = position1234, tokenIndex1234, depth1234
								if buffer[position] != rune('T') {
									goto l1222
								}
								position++
							}
						l1234:
							{
								position1236, tokenIndex1236, depth1236 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1237
								}
								position++
								goto l1236
							l1237:
								position, tokenIndex, depth = position1236, tokenIndex1236, depth1236
								if buffer[position] != rune('E') {
									goto l1222
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
									goto l1222
								}
								position++
							}
						l1238:
							if !rules[ruleskip]() {
								goto l1222
							}
							depth--
							add(ruleSTRAFTER, position1223)
						}
						goto l1165
					l1222:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
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
								if buffer[position] != rune('t') {
									goto l1245
								}
								position++
								goto l1244
							l1245:
								position, tokenIndex, depth = position1244, tokenIndex1244, depth1244
								if buffer[position] != rune('T') {
									goto l1240
								}
								position++
							}
						l1244:
							{
								position1246, tokenIndex1246, depth1246 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1247
								}
								position++
								goto l1246
							l1247:
								position, tokenIndex, depth = position1246, tokenIndex1246, depth1246
								if buffer[position] != rune('R') {
									goto l1240
								}
								position++
							}
						l1246:
							{
								position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1249
								}
								position++
								goto l1248
							l1249:
								position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
								if buffer[position] != rune('L') {
									goto l1240
								}
								position++
							}
						l1248:
							{
								position1250, tokenIndex1250, depth1250 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1251
								}
								position++
								goto l1250
							l1251:
								position, tokenIndex, depth = position1250, tokenIndex1250, depth1250
								if buffer[position] != rune('A') {
									goto l1240
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
									goto l1240
								}
								position++
							}
						l1252:
							{
								position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1255
								}
								position++
								goto l1254
							l1255:
								position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
								if buffer[position] != rune('G') {
									goto l1240
								}
								position++
							}
						l1254:
							if !rules[ruleskip]() {
								goto l1240
							}
							depth--
							add(ruleSTRLANG, position1241)
						}
						goto l1165
					l1240:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						{
							position1257 := position
							depth++
							{
								position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1259
								}
								position++
								goto l1258
							l1259:
								position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
								if buffer[position] != rune('S') {
									goto l1256
								}
								position++
							}
						l1258:
							{
								position1260, tokenIndex1260, depth1260 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1261
								}
								position++
								goto l1260
							l1261:
								position, tokenIndex, depth = position1260, tokenIndex1260, depth1260
								if buffer[position] != rune('T') {
									goto l1256
								}
								position++
							}
						l1260:
							{
								position1262, tokenIndex1262, depth1262 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1263
								}
								position++
								goto l1262
							l1263:
								position, tokenIndex, depth = position1262, tokenIndex1262, depth1262
								if buffer[position] != rune('R') {
									goto l1256
								}
								position++
							}
						l1262:
							{
								position1264, tokenIndex1264, depth1264 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1265
								}
								position++
								goto l1264
							l1265:
								position, tokenIndex, depth = position1264, tokenIndex1264, depth1264
								if buffer[position] != rune('D') {
									goto l1256
								}
								position++
							}
						l1264:
							{
								position1266, tokenIndex1266, depth1266 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1267
								}
								position++
								goto l1266
							l1267:
								position, tokenIndex, depth = position1266, tokenIndex1266, depth1266
								if buffer[position] != rune('T') {
									goto l1256
								}
								position++
							}
						l1266:
							if !rules[ruleskip]() {
								goto l1256
							}
							depth--
							add(ruleSTRDT, position1257)
						}
						goto l1165
					l1256:
						position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1269 := position
									depth++
									{
										position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1271
										}
										position++
										goto l1270
									l1271:
										position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
										if buffer[position] != rune('S') {
											goto l1164
										}
										position++
									}
								l1270:
									{
										position1272, tokenIndex1272, depth1272 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1273
										}
										position++
										goto l1272
									l1273:
										position, tokenIndex, depth = position1272, tokenIndex1272, depth1272
										if buffer[position] != rune('A') {
											goto l1164
										}
										position++
									}
								l1272:
									{
										position1274, tokenIndex1274, depth1274 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1275
										}
										position++
										goto l1274
									l1275:
										position, tokenIndex, depth = position1274, tokenIndex1274, depth1274
										if buffer[position] != rune('M') {
											goto l1164
										}
										position++
									}
								l1274:
									{
										position1276, tokenIndex1276, depth1276 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1277
										}
										position++
										goto l1276
									l1277:
										position, tokenIndex, depth = position1276, tokenIndex1276, depth1276
										if buffer[position] != rune('E') {
											goto l1164
										}
										position++
									}
								l1276:
									{
										position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1279
										}
										position++
										goto l1278
									l1279:
										position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
										if buffer[position] != rune('T') {
											goto l1164
										}
										position++
									}
								l1278:
									{
										position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1281
										}
										position++
										goto l1280
									l1281:
										position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
										if buffer[position] != rune('E') {
											goto l1164
										}
										position++
									}
								l1280:
									{
										position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1283
										}
										position++
										goto l1282
									l1283:
										position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
										if buffer[position] != rune('R') {
											goto l1164
										}
										position++
									}
								l1282:
									{
										position1284, tokenIndex1284, depth1284 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1285
										}
										position++
										goto l1284
									l1285:
										position, tokenIndex, depth = position1284, tokenIndex1284, depth1284
										if buffer[position] != rune('M') {
											goto l1164
										}
										position++
									}
								l1284:
									if !rules[ruleskip]() {
										goto l1164
									}
									depth--
									add(ruleSAMETERM, position1269)
								}
								break
							case 'C', 'c':
								{
									position1286 := position
									depth++
									{
										position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1288
										}
										position++
										goto l1287
									l1288:
										position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
										if buffer[position] != rune('C') {
											goto l1164
										}
										position++
									}
								l1287:
									{
										position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1290
										}
										position++
										goto l1289
									l1290:
										position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
										if buffer[position] != rune('O') {
											goto l1164
										}
										position++
									}
								l1289:
									{
										position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1292
										}
										position++
										goto l1291
									l1292:
										position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
										if buffer[position] != rune('N') {
											goto l1164
										}
										position++
									}
								l1291:
									{
										position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1294
										}
										position++
										goto l1293
									l1294:
										position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
										if buffer[position] != rune('T') {
											goto l1164
										}
										position++
									}
								l1293:
									{
										position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1296
										}
										position++
										goto l1295
									l1296:
										position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
										if buffer[position] != rune('A') {
											goto l1164
										}
										position++
									}
								l1295:
									{
										position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1298
										}
										position++
										goto l1297
									l1298:
										position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
										if buffer[position] != rune('I') {
											goto l1164
										}
										position++
									}
								l1297:
									{
										position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1300
										}
										position++
										goto l1299
									l1300:
										position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
										if buffer[position] != rune('N') {
											goto l1164
										}
										position++
									}
								l1299:
									{
										position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1302
										}
										position++
										goto l1301
									l1302:
										position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
										if buffer[position] != rune('S') {
											goto l1164
										}
										position++
									}
								l1301:
									if !rules[ruleskip]() {
										goto l1164
									}
									depth--
									add(ruleCONTAINS, position1286)
								}
								break
							default:
								{
									position1303 := position
									depth++
									{
										position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1305
										}
										position++
										goto l1304
									l1305:
										position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
										if buffer[position] != rune('L') {
											goto l1164
										}
										position++
									}
								l1304:
									{
										position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1307
										}
										position++
										goto l1306
									l1307:
										position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
										if buffer[position] != rune('A') {
											goto l1164
										}
										position++
									}
								l1306:
									{
										position1308, tokenIndex1308, depth1308 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1309
										}
										position++
										goto l1308
									l1309:
										position, tokenIndex, depth = position1308, tokenIndex1308, depth1308
										if buffer[position] != rune('N') {
											goto l1164
										}
										position++
									}
								l1308:
									{
										position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1311
										}
										position++
										goto l1310
									l1311:
										position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
										if buffer[position] != rune('G') {
											goto l1164
										}
										position++
									}
								l1310:
									{
										position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1313
										}
										position++
										goto l1312
									l1313:
										position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
										if buffer[position] != rune('M') {
											goto l1164
										}
										position++
									}
								l1312:
									{
										position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1315
										}
										position++
										goto l1314
									l1315:
										position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
										if buffer[position] != rune('A') {
											goto l1164
										}
										position++
									}
								l1314:
									{
										position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1317
										}
										position++
										goto l1316
									l1317:
										position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
										if buffer[position] != rune('T') {
											goto l1164
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
											goto l1164
										}
										position++
									}
								l1318:
									{
										position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1321
										}
										position++
										goto l1320
									l1321:
										position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
										if buffer[position] != rune('H') {
											goto l1164
										}
										position++
									}
								l1320:
									{
										position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1323
										}
										position++
										goto l1322
									l1323:
										position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
										if buffer[position] != rune('E') {
											goto l1164
										}
										position++
									}
								l1322:
									{
										position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1325
										}
										position++
										goto l1324
									l1325:
										position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
										if buffer[position] != rune('S') {
											goto l1164
										}
										position++
									}
								l1324:
									if !rules[ruleskip]() {
										goto l1164
									}
									depth--
									add(ruleLANGMATCHES, position1303)
								}
								break
							}
						}

					}
				l1165:
					if !rules[ruleLPAREN]() {
						goto l1164
					}
					if !rules[ruleexpression]() {
						goto l1164
					}
					if !rules[ruleCOMMA]() {
						goto l1164
					}
					if !rules[ruleexpression]() {
						goto l1164
					}
					if !rules[ruleRPAREN]() {
						goto l1164
					}
					goto l800
				l1164:
					position, tokenIndex, depth = position800, tokenIndex800, depth800
					{
						position1327 := position
						depth++
						{
							position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1329
							}
							position++
							goto l1328
						l1329:
							position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
							if buffer[position] != rune('B') {
								goto l1326
							}
							position++
						}
					l1328:
						{
							position1330, tokenIndex1330, depth1330 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1331
							}
							position++
							goto l1330
						l1331:
							position, tokenIndex, depth = position1330, tokenIndex1330, depth1330
							if buffer[position] != rune('O') {
								goto l1326
							}
							position++
						}
					l1330:
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
								goto l1326
							}
							position++
						}
					l1332:
						{
							position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1335
							}
							position++
							goto l1334
						l1335:
							position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
							if buffer[position] != rune('N') {
								goto l1326
							}
							position++
						}
					l1334:
						{
							position1336, tokenIndex1336, depth1336 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1337
							}
							position++
							goto l1336
						l1337:
							position, tokenIndex, depth = position1336, tokenIndex1336, depth1336
							if buffer[position] != rune('D') {
								goto l1326
							}
							position++
						}
					l1336:
						if !rules[ruleskip]() {
							goto l1326
						}
						depth--
						add(ruleBOUND, position1327)
					}
					if !rules[ruleLPAREN]() {
						goto l1326
					}
					if !rules[rulevar]() {
						goto l1326
					}
					if !rules[ruleRPAREN]() {
						goto l1326
					}
					goto l800
				l1326:
					position, tokenIndex, depth = position800, tokenIndex800, depth800
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1340 := position
								depth++
								{
									position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1342
									}
									position++
									goto l1341
								l1342:
									position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
									if buffer[position] != rune('S') {
										goto l1338
									}
									position++
								}
							l1341:
								{
									position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1344
									}
									position++
									goto l1343
								l1344:
									position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
									if buffer[position] != rune('T') {
										goto l1338
									}
									position++
								}
							l1343:
								{
									position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1346
									}
									position++
									goto l1345
								l1346:
									position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
									if buffer[position] != rune('R') {
										goto l1338
									}
									position++
								}
							l1345:
								{
									position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1348
									}
									position++
									goto l1347
								l1348:
									position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
									if buffer[position] != rune('U') {
										goto l1338
									}
									position++
								}
							l1347:
								{
									position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1350
									}
									position++
									goto l1349
								l1350:
									position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
									if buffer[position] != rune('U') {
										goto l1338
									}
									position++
								}
							l1349:
								{
									position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1352
									}
									position++
									goto l1351
								l1352:
									position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
									if buffer[position] != rune('I') {
										goto l1338
									}
									position++
								}
							l1351:
								{
									position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1354
									}
									position++
									goto l1353
								l1354:
									position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
									if buffer[position] != rune('D') {
										goto l1338
									}
									position++
								}
							l1353:
								if !rules[ruleskip]() {
									goto l1338
								}
								depth--
								add(ruleSTRUUID, position1340)
							}
							break
						case 'U', 'u':
							{
								position1355 := position
								depth++
								{
									position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1357
									}
									position++
									goto l1356
								l1357:
									position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
									if buffer[position] != rune('U') {
										goto l1338
									}
									position++
								}
							l1356:
								{
									position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1359
									}
									position++
									goto l1358
								l1359:
									position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
									if buffer[position] != rune('U') {
										goto l1338
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
										goto l1338
									}
									position++
								}
							l1360:
								{
									position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1363
									}
									position++
									goto l1362
								l1363:
									position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
									if buffer[position] != rune('D') {
										goto l1338
									}
									position++
								}
							l1362:
								if !rules[ruleskip]() {
									goto l1338
								}
								depth--
								add(ruleUUID, position1355)
							}
							break
						case 'N', 'n':
							{
								position1364 := position
								depth++
								{
									position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1366
									}
									position++
									goto l1365
								l1366:
									position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
									if buffer[position] != rune('N') {
										goto l1338
									}
									position++
								}
							l1365:
								{
									position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1368
									}
									position++
									goto l1367
								l1368:
									position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
									if buffer[position] != rune('O') {
										goto l1338
									}
									position++
								}
							l1367:
								{
									position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1370
									}
									position++
									goto l1369
								l1370:
									position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
									if buffer[position] != rune('W') {
										goto l1338
									}
									position++
								}
							l1369:
								if !rules[ruleskip]() {
									goto l1338
								}
								depth--
								add(ruleNOW, position1364)
							}
							break
						default:
							{
								position1371 := position
								depth++
								{
									position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1373
									}
									position++
									goto l1372
								l1373:
									position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
									if buffer[position] != rune('R') {
										goto l1338
									}
									position++
								}
							l1372:
								{
									position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1375
									}
									position++
									goto l1374
								l1375:
									position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
									if buffer[position] != rune('A') {
										goto l1338
									}
									position++
								}
							l1374:
								{
									position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1377
									}
									position++
									goto l1376
								l1377:
									position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
									if buffer[position] != rune('N') {
										goto l1338
									}
									position++
								}
							l1376:
								{
									position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1379
									}
									position++
									goto l1378
								l1379:
									position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
									if buffer[position] != rune('D') {
										goto l1338
									}
									position++
								}
							l1378:
								if !rules[ruleskip]() {
									goto l1338
								}
								depth--
								add(ruleRAND, position1371)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1338
					}
					goto l800
				l1338:
					position, tokenIndex, depth = position800, tokenIndex800, depth800
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
								{
									position1383 := position
									depth++
									{
										position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1385
										}
										position++
										goto l1384
									l1385:
										position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
										if buffer[position] != rune('E') {
											goto l1382
										}
										position++
									}
								l1384:
									{
										position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1387
										}
										position++
										goto l1386
									l1387:
										position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
										if buffer[position] != rune('X') {
											goto l1382
										}
										position++
									}
								l1386:
									{
										position1388, tokenIndex1388, depth1388 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1389
										}
										position++
										goto l1388
									l1389:
										position, tokenIndex, depth = position1388, tokenIndex1388, depth1388
										if buffer[position] != rune('I') {
											goto l1382
										}
										position++
									}
								l1388:
									{
										position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1391
										}
										position++
										goto l1390
									l1391:
										position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
										if buffer[position] != rune('S') {
											goto l1382
										}
										position++
									}
								l1390:
									{
										position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1393
										}
										position++
										goto l1392
									l1393:
										position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
										if buffer[position] != rune('T') {
											goto l1382
										}
										position++
									}
								l1392:
									{
										position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1395
										}
										position++
										goto l1394
									l1395:
										position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
										if buffer[position] != rune('S') {
											goto l1382
										}
										position++
									}
								l1394:
									if !rules[ruleskip]() {
										goto l1382
									}
									depth--
									add(ruleEXISTS, position1383)
								}
								goto l1381
							l1382:
								position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
								{
									position1396 := position
									depth++
									{
										position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1398
										}
										position++
										goto l1397
									l1398:
										position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
										if buffer[position] != rune('N') {
											goto l798
										}
										position++
									}
								l1397:
									{
										position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1400
										}
										position++
										goto l1399
									l1400:
										position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
										if buffer[position] != rune('O') {
											goto l798
										}
										position++
									}
								l1399:
									{
										position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1402
										}
										position++
										goto l1401
									l1402:
										position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
										if buffer[position] != rune('T') {
											goto l798
										}
										position++
									}
								l1401:
									if buffer[position] != rune(' ') {
										goto l798
									}
									position++
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
											goto l798
										}
										position++
									}
								l1403:
									{
										position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1406
										}
										position++
										goto l1405
									l1406:
										position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
										if buffer[position] != rune('X') {
											goto l798
										}
										position++
									}
								l1405:
									{
										position1407, tokenIndex1407, depth1407 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1408
										}
										position++
										goto l1407
									l1408:
										position, tokenIndex, depth = position1407, tokenIndex1407, depth1407
										if buffer[position] != rune('I') {
											goto l798
										}
										position++
									}
								l1407:
									{
										position1409, tokenIndex1409, depth1409 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1410
										}
										position++
										goto l1409
									l1410:
										position, tokenIndex, depth = position1409, tokenIndex1409, depth1409
										if buffer[position] != rune('S') {
											goto l798
										}
										position++
									}
								l1409:
									{
										position1411, tokenIndex1411, depth1411 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1412
										}
										position++
										goto l1411
									l1412:
										position, tokenIndex, depth = position1411, tokenIndex1411, depth1411
										if buffer[position] != rune('T') {
											goto l798
										}
										position++
									}
								l1411:
									{
										position1413, tokenIndex1413, depth1413 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1414
										}
										position++
										goto l1413
									l1414:
										position, tokenIndex, depth = position1413, tokenIndex1413, depth1413
										if buffer[position] != rune('S') {
											goto l798
										}
										position++
									}
								l1413:
									if !rules[ruleskip]() {
										goto l798
									}
									depth--
									add(ruleNOTEXIST, position1396)
								}
							}
						l1381:
							if !rules[rulegroupGraphPattern]() {
								goto l798
							}
							break
						case 'I', 'i':
							{
								position1415 := position
								depth++
								{
									position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1417
									}
									position++
									goto l1416
								l1417:
									position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
									if buffer[position] != rune('I') {
										goto l798
									}
									position++
								}
							l1416:
								{
									position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1419
									}
									position++
									goto l1418
								l1419:
									position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
									if buffer[position] != rune('F') {
										goto l798
									}
									position++
								}
							l1418:
								if !rules[ruleskip]() {
									goto l798
								}
								depth--
								add(ruleIF, position1415)
							}
							if !rules[ruleLPAREN]() {
								goto l798
							}
							if !rules[ruleexpression]() {
								goto l798
							}
							if !rules[ruleCOMMA]() {
								goto l798
							}
							if !rules[ruleexpression]() {
								goto l798
							}
							if !rules[ruleCOMMA]() {
								goto l798
							}
							if !rules[ruleexpression]() {
								goto l798
							}
							if !rules[ruleRPAREN]() {
								goto l798
							}
							break
						case 'C', 'c':
							{
								position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
								{
									position1422 := position
									depth++
									{
										position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1424
										}
										position++
										goto l1423
									l1424:
										position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
										if buffer[position] != rune('C') {
											goto l1421
										}
										position++
									}
								l1423:
									{
										position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1426
										}
										position++
										goto l1425
									l1426:
										position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
										if buffer[position] != rune('O') {
											goto l1421
										}
										position++
									}
								l1425:
									{
										position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1428
										}
										position++
										goto l1427
									l1428:
										position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
										if buffer[position] != rune('N') {
											goto l1421
										}
										position++
									}
								l1427:
									{
										position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1430
										}
										position++
										goto l1429
									l1430:
										position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
										if buffer[position] != rune('C') {
											goto l1421
										}
										position++
									}
								l1429:
									{
										position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1432
										}
										position++
										goto l1431
									l1432:
										position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
										if buffer[position] != rune('A') {
											goto l1421
										}
										position++
									}
								l1431:
									{
										position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1434
										}
										position++
										goto l1433
									l1434:
										position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
										if buffer[position] != rune('T') {
											goto l1421
										}
										position++
									}
								l1433:
									if !rules[ruleskip]() {
										goto l1421
									}
									depth--
									add(ruleCONCAT, position1422)
								}
								goto l1420
							l1421:
								position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
								{
									position1435 := position
									depth++
									{
										position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1437
										}
										position++
										goto l1436
									l1437:
										position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
										if buffer[position] != rune('C') {
											goto l798
										}
										position++
									}
								l1436:
									{
										position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1439
										}
										position++
										goto l1438
									l1439:
										position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
										if buffer[position] != rune('O') {
											goto l798
										}
										position++
									}
								l1438:
									{
										position1440, tokenIndex1440, depth1440 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1441
										}
										position++
										goto l1440
									l1441:
										position, tokenIndex, depth = position1440, tokenIndex1440, depth1440
										if buffer[position] != rune('A') {
											goto l798
										}
										position++
									}
								l1440:
									{
										position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1443
										}
										position++
										goto l1442
									l1443:
										position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
										if buffer[position] != rune('L') {
											goto l798
										}
										position++
									}
								l1442:
									{
										position1444, tokenIndex1444, depth1444 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1445
										}
										position++
										goto l1444
									l1445:
										position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
										if buffer[position] != rune('E') {
											goto l798
										}
										position++
									}
								l1444:
									{
										position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1447
										}
										position++
										goto l1446
									l1447:
										position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
										if buffer[position] != rune('S') {
											goto l798
										}
										position++
									}
								l1446:
									{
										position1448, tokenIndex1448, depth1448 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1449
										}
										position++
										goto l1448
									l1449:
										position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
										if buffer[position] != rune('C') {
											goto l798
										}
										position++
									}
								l1448:
									{
										position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1451
										}
										position++
										goto l1450
									l1451:
										position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
										if buffer[position] != rune('E') {
											goto l798
										}
										position++
									}
								l1450:
									if !rules[ruleskip]() {
										goto l798
									}
									depth--
									add(ruleCOALESCE, position1435)
								}
							}
						l1420:
							if !rules[ruleargList]() {
								goto l798
							}
							break
						case 'B', 'b':
							{
								position1452 := position
								depth++
								{
									position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1454
									}
									position++
									goto l1453
								l1454:
									position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
									if buffer[position] != rune('B') {
										goto l798
									}
									position++
								}
							l1453:
								{
									position1455, tokenIndex1455, depth1455 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1456
									}
									position++
									goto l1455
								l1456:
									position, tokenIndex, depth = position1455, tokenIndex1455, depth1455
									if buffer[position] != rune('N') {
										goto l798
									}
									position++
								}
							l1455:
								{
									position1457, tokenIndex1457, depth1457 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1458
									}
									position++
									goto l1457
								l1458:
									position, tokenIndex, depth = position1457, tokenIndex1457, depth1457
									if buffer[position] != rune('O') {
										goto l798
									}
									position++
								}
							l1457:
								{
									position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1460
									}
									position++
									goto l1459
								l1460:
									position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
									if buffer[position] != rune('D') {
										goto l798
									}
									position++
								}
							l1459:
								{
									position1461, tokenIndex1461, depth1461 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1462
									}
									position++
									goto l1461
								l1462:
									position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
									if buffer[position] != rune('E') {
										goto l798
									}
									position++
								}
							l1461:
								if !rules[ruleskip]() {
									goto l798
								}
								depth--
								add(ruleBNODE, position1452)
							}
							{
								position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1464
								}
								if !rules[ruleexpression]() {
									goto l1464
								}
								if !rules[ruleRPAREN]() {
									goto l1464
								}
								goto l1463
							l1464:
								position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
								if !rules[rulenil]() {
									goto l798
								}
							}
						l1463:
							break
						default:
							{
								position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
								{
									position1467 := position
									depth++
									{
										position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1469
										}
										position++
										goto l1468
									l1469:
										position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
										if buffer[position] != rune('S') {
											goto l1466
										}
										position++
									}
								l1468:
									{
										position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1471
										}
										position++
										goto l1470
									l1471:
										position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
										if buffer[position] != rune('U') {
											goto l1466
										}
										position++
									}
								l1470:
									{
										position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1473
										}
										position++
										goto l1472
									l1473:
										position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
										if buffer[position] != rune('B') {
											goto l1466
										}
										position++
									}
								l1472:
									{
										position1474, tokenIndex1474, depth1474 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1475
										}
										position++
										goto l1474
									l1475:
										position, tokenIndex, depth = position1474, tokenIndex1474, depth1474
										if buffer[position] != rune('S') {
											goto l1466
										}
										position++
									}
								l1474:
									{
										position1476, tokenIndex1476, depth1476 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1477
										}
										position++
										goto l1476
									l1477:
										position, tokenIndex, depth = position1476, tokenIndex1476, depth1476
										if buffer[position] != rune('T') {
											goto l1466
										}
										position++
									}
								l1476:
									{
										position1478, tokenIndex1478, depth1478 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1479
										}
										position++
										goto l1478
									l1479:
										position, tokenIndex, depth = position1478, tokenIndex1478, depth1478
										if buffer[position] != rune('R') {
											goto l1466
										}
										position++
									}
								l1478:
									if !rules[ruleskip]() {
										goto l1466
									}
									depth--
									add(ruleSUBSTR, position1467)
								}
								goto l1465
							l1466:
								position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
								{
									position1481 := position
									depth++
									{
										position1482, tokenIndex1482, depth1482 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1483
										}
										position++
										goto l1482
									l1483:
										position, tokenIndex, depth = position1482, tokenIndex1482, depth1482
										if buffer[position] != rune('R') {
											goto l1480
										}
										position++
									}
								l1482:
									{
										position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1485
										}
										position++
										goto l1484
									l1485:
										position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
										if buffer[position] != rune('E') {
											goto l1480
										}
										position++
									}
								l1484:
									{
										position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1487
										}
										position++
										goto l1486
									l1487:
										position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
										if buffer[position] != rune('P') {
											goto l1480
										}
										position++
									}
								l1486:
									{
										position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1489
										}
										position++
										goto l1488
									l1489:
										position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
										if buffer[position] != rune('L') {
											goto l1480
										}
										position++
									}
								l1488:
									{
										position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1491
										}
										position++
										goto l1490
									l1491:
										position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
										if buffer[position] != rune('A') {
											goto l1480
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
											goto l1480
										}
										position++
									}
								l1492:
									{
										position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1495
										}
										position++
										goto l1494
									l1495:
										position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
										if buffer[position] != rune('E') {
											goto l1480
										}
										position++
									}
								l1494:
									if !rules[ruleskip]() {
										goto l1480
									}
									depth--
									add(ruleREPLACE, position1481)
								}
								goto l1465
							l1480:
								position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
								{
									position1496 := position
									depth++
									{
										position1497, tokenIndex1497, depth1497 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1498
										}
										position++
										goto l1497
									l1498:
										position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
										if buffer[position] != rune('R') {
											goto l798
										}
										position++
									}
								l1497:
									{
										position1499, tokenIndex1499, depth1499 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1500
										}
										position++
										goto l1499
									l1500:
										position, tokenIndex, depth = position1499, tokenIndex1499, depth1499
										if buffer[position] != rune('E') {
											goto l798
										}
										position++
									}
								l1499:
									{
										position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1502
										}
										position++
										goto l1501
									l1502:
										position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
										if buffer[position] != rune('G') {
											goto l798
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
											goto l798
										}
										position++
									}
								l1503:
									{
										position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1506
										}
										position++
										goto l1505
									l1506:
										position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
										if buffer[position] != rune('X') {
											goto l798
										}
										position++
									}
								l1505:
									if !rules[ruleskip]() {
										goto l798
									}
									depth--
									add(ruleREGEX, position1496)
								}
							}
						l1465:
							if !rules[ruleLPAREN]() {
								goto l798
							}
							if !rules[ruleexpression]() {
								goto l798
							}
							if !rules[ruleCOMMA]() {
								goto l798
							}
							if !rules[ruleexpression]() {
								goto l798
							}
							{
								position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1507
								}
								if !rules[ruleexpression]() {
									goto l1507
								}
								goto l1508
							l1507:
								position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
							}
						l1508:
							if !rules[ruleRPAREN]() {
								goto l798
							}
							break
						}
					}

				}
			l800:
				depth--
				add(rulebuiltinCall, position799)
			}
			return true
		l798:
			position, tokenIndex, depth = position798, tokenIndex798, depth798
			return false
		},
		/* 69 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
			{
				position1510 := position
				depth++
				{
					position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1512
					}
					position++
					goto l1511
				l1512:
					position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
					if buffer[position] != rune('$') {
						goto l1509
					}
					position++
				}
			l1511:
				{
					position1513 := position
					depth++
					{
						position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
						{
							position1518 := position
							depth++
							{
								position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
								{
									position1521 := position
									depth++
									{
										position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1523
										}
										position++
										goto l1522
									l1523:
										position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1520
										}
										position++
									}
								l1522:
									depth--
									add(rulePN_CHARS_BASE, position1521)
								}
								goto l1519
							l1520:
								position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
								if buffer[position] != rune('_') {
									goto l1517
								}
								position++
							}
						l1519:
							depth--
							add(rulePN_CHARS_U, position1518)
						}
						goto l1516
					l1517:
						position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1509
						}
						position++
					}
				l1516:
				l1514:
					{
						position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
						{
							position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
							{
								position1526 := position
								depth++
								{
									position1527, tokenIndex1527, depth1527 := position, tokenIndex, depth
									{
										position1529 := position
										depth++
										{
											position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1531
											}
											position++
											goto l1530
										l1531:
											position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1528
											}
											position++
										}
									l1530:
										depth--
										add(rulePN_CHARS_BASE, position1529)
									}
									goto l1527
								l1528:
									position, tokenIndex, depth = position1527, tokenIndex1527, depth1527
									if buffer[position] != rune('_') {
										goto l1525
									}
									position++
								}
							l1527:
								depth--
								add(rulePN_CHARS_U, position1526)
							}
							goto l1524
						l1525:
							position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1515
							}
							position++
						}
					l1524:
						goto l1514
					l1515:
						position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
					}
					depth--
					add(ruleVARNAME, position1513)
				}
				if !rules[ruleskip]() {
					goto l1509
				}
				depth--
				add(rulevar, position1510)
			}
			return true
		l1509:
			position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
			return false
		},
		/* 70 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1532, tokenIndex1532, depth1532 := position, tokenIndex, depth
			{
				position1533 := position
				depth++
				{
					position1534, tokenIndex1534, depth1534 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1535
					}
					goto l1534
				l1535:
					position, tokenIndex, depth = position1534, tokenIndex1534, depth1534
					{
						position1536 := position
						depth++
					l1537:
						{
							position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
							{
								position1539, tokenIndex1539, depth1539 := position, tokenIndex, depth
								{
									position1540, tokenIndex1540, depth1540 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1541
									}
									position++
									goto l1540
								l1541:
									position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
									if buffer[position] != rune(' ') {
										goto l1539
									}
									position++
								}
							l1540:
								goto l1538
							l1539:
								position, tokenIndex, depth = position1539, tokenIndex1539, depth1539
							}
							if !matchDot() {
								goto l1538
							}
							goto l1537
						l1538:
							position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
						}
						if buffer[position] != rune(':') {
							goto l1532
						}
						position++
					l1542:
						{
							position1543, tokenIndex1543, depth1543 := position, tokenIndex, depth
							{
								position1544, tokenIndex1544, depth1544 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1545
								}
								position++
								goto l1544
							l1545:
								position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1546
								}
								position++
								goto l1544
							l1546:
								position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1547
								}
								position++
								goto l1544
							l1547:
								position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1543
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1543
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1543
										}
										position++
										break
									}
								}

							}
						l1544:
							goto l1542
						l1543:
							position, tokenIndex, depth = position1543, tokenIndex1543, depth1543
						}
						if !rules[ruleskip]() {
							goto l1532
						}
						depth--
						add(ruleprefixedName, position1536)
					}
				}
			l1534:
				depth--
				add(ruleiriref, position1533)
			}
			return true
		l1532:
			position, tokenIndex, depth = position1532, tokenIndex1532, depth1532
			return false
		},
		/* 71 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
			{
				position1550 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1549
				}
				position++
			l1551:
				{
					position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
					{
						position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1553
						}
						position++
						goto l1552
					l1553:
						position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
					}
					if !matchDot() {
						goto l1552
					}
					goto l1551
				l1552:
					position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
				}
				if buffer[position] != rune('>') {
					goto l1549
				}
				position++
				if !rules[ruleskip]() {
					goto l1549
				}
				depth--
				add(ruleiri, position1550)
			}
			return true
		l1549:
			position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
			return false
		},
		/* 72 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 73 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1555, tokenIndex1555, depth1555 := position, tokenIndex, depth
			{
				position1556 := position
				depth++
				if !rules[rulestring]() {
					goto l1555
				}
				{
					position1557, tokenIndex1557, depth1557 := position, tokenIndex, depth
					{
						position1559, tokenIndex1559, depth1559 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1560
						}
						position++
						{
							position1563, tokenIndex1563, depth1563 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1564
							}
							position++
							goto l1563
						l1564:
							position, tokenIndex, depth = position1563, tokenIndex1563, depth1563
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1560
							}
							position++
						}
					l1563:
					l1561:
						{
							position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
							{
								position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1566
								}
								position++
								goto l1565
							l1566:
								position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1562
								}
								position++
							}
						l1565:
							goto l1561
						l1562:
							position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
						}
					l1567:
						{
							position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1568
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1568
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1568
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1568
									}
									position++
									break
								}
							}

						l1569:
							{
								position1570, tokenIndex1570, depth1570 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1570
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1570
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1570
										}
										position++
										break
									}
								}

								goto l1569
							l1570:
								position, tokenIndex, depth = position1570, tokenIndex1570, depth1570
							}
							goto l1567
						l1568:
							position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
						}
						goto l1559
					l1560:
						position, tokenIndex, depth = position1559, tokenIndex1559, depth1559
						if buffer[position] != rune('^') {
							goto l1557
						}
						position++
						if buffer[position] != rune('^') {
							goto l1557
						}
						position++
						if !rules[ruleiriref]() {
							goto l1557
						}
					}
				l1559:
					goto l1558
				l1557:
					position, tokenIndex, depth = position1557, tokenIndex1557, depth1557
				}
			l1558:
				if !rules[ruleskip]() {
					goto l1555
				}
				depth--
				add(ruleliteral, position1556)
			}
			return true
		l1555:
			position, tokenIndex, depth = position1555, tokenIndex1555, depth1555
			return false
		},
		/* 74 string <- <('"' (!'"' .)* '"')> */
		func() bool {
			position1573, tokenIndex1573, depth1573 := position, tokenIndex, depth
			{
				position1574 := position
				depth++
				if buffer[position] != rune('"') {
					goto l1573
				}
				position++
			l1575:
				{
					position1576, tokenIndex1576, depth1576 := position, tokenIndex, depth
					{
						position1577, tokenIndex1577, depth1577 := position, tokenIndex, depth
						if buffer[position] != rune('"') {
							goto l1577
						}
						position++
						goto l1576
					l1577:
						position, tokenIndex, depth = position1577, tokenIndex1577, depth1577
					}
					if !matchDot() {
						goto l1576
					}
					goto l1575
				l1576:
					position, tokenIndex, depth = position1576, tokenIndex1576, depth1576
				}
				if buffer[position] != rune('"') {
					goto l1573
				}
				position++
				depth--
				add(rulestring, position1574)
			}
			return true
		l1573:
			position, tokenIndex, depth = position1573, tokenIndex1573, depth1573
			return false
		},
		/* 75 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1578, tokenIndex1578, depth1578 := position, tokenIndex, depth
			{
				position1579 := position
				depth++
				{
					position1580, tokenIndex1580, depth1580 := position, tokenIndex, depth
					{
						position1582, tokenIndex1582, depth1582 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1583
						}
						position++
						goto l1582
					l1583:
						position, tokenIndex, depth = position1582, tokenIndex1582, depth1582
						if buffer[position] != rune('-') {
							goto l1580
						}
						position++
					}
				l1582:
					goto l1581
				l1580:
					position, tokenIndex, depth = position1580, tokenIndex1580, depth1580
				}
			l1581:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1578
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
				{
					position1586, tokenIndex1586, depth1586 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1586
					}
					position++
				l1588:
					{
						position1589, tokenIndex1589, depth1589 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1589
						}
						position++
						goto l1588
					l1589:
						position, tokenIndex, depth = position1589, tokenIndex1589, depth1589
					}
					goto l1587
				l1586:
					position, tokenIndex, depth = position1586, tokenIndex1586, depth1586
				}
			l1587:
				if !rules[ruleskip]() {
					goto l1578
				}
				depth--
				add(rulenumericLiteral, position1579)
			}
			return true
		l1578:
			position, tokenIndex, depth = position1578, tokenIndex1578, depth1578
			return false
		},
		/* 76 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 77 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1591, tokenIndex1591, depth1591 := position, tokenIndex, depth
			{
				position1592 := position
				depth++
				{
					position1593, tokenIndex1593, depth1593 := position, tokenIndex, depth
					{
						position1595 := position
						depth++
						{
							position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1597
							}
							position++
							goto l1596
						l1597:
							position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
							if buffer[position] != rune('T') {
								goto l1594
							}
							position++
						}
					l1596:
						{
							position1598, tokenIndex1598, depth1598 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1599
							}
							position++
							goto l1598
						l1599:
							position, tokenIndex, depth = position1598, tokenIndex1598, depth1598
							if buffer[position] != rune('R') {
								goto l1594
							}
							position++
						}
					l1598:
						{
							position1600, tokenIndex1600, depth1600 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1601
							}
							position++
							goto l1600
						l1601:
							position, tokenIndex, depth = position1600, tokenIndex1600, depth1600
							if buffer[position] != rune('U') {
								goto l1594
							}
							position++
						}
					l1600:
						{
							position1602, tokenIndex1602, depth1602 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1603
							}
							position++
							goto l1602
						l1603:
							position, tokenIndex, depth = position1602, tokenIndex1602, depth1602
							if buffer[position] != rune('E') {
								goto l1594
							}
							position++
						}
					l1602:
						if !rules[ruleskip]() {
							goto l1594
						}
						depth--
						add(ruleTRUE, position1595)
					}
					goto l1593
				l1594:
					position, tokenIndex, depth = position1593, tokenIndex1593, depth1593
					{
						position1604 := position
						depth++
						{
							position1605, tokenIndex1605, depth1605 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1606
							}
							position++
							goto l1605
						l1606:
							position, tokenIndex, depth = position1605, tokenIndex1605, depth1605
							if buffer[position] != rune('F') {
								goto l1591
							}
							position++
						}
					l1605:
						{
							position1607, tokenIndex1607, depth1607 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1608
							}
							position++
							goto l1607
						l1608:
							position, tokenIndex, depth = position1607, tokenIndex1607, depth1607
							if buffer[position] != rune('A') {
								goto l1591
							}
							position++
						}
					l1607:
						{
							position1609, tokenIndex1609, depth1609 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1610
							}
							position++
							goto l1609
						l1610:
							position, tokenIndex, depth = position1609, tokenIndex1609, depth1609
							if buffer[position] != rune('L') {
								goto l1591
							}
							position++
						}
					l1609:
						{
							position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1612
							}
							position++
							goto l1611
						l1612:
							position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
							if buffer[position] != rune('S') {
								goto l1591
							}
							position++
						}
					l1611:
						{
							position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1614
							}
							position++
							goto l1613
						l1614:
							position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
							if buffer[position] != rune('E') {
								goto l1591
							}
							position++
						}
					l1613:
						if !rules[ruleskip]() {
							goto l1591
						}
						depth--
						add(ruleFALSE, position1604)
					}
				}
			l1593:
				depth--
				add(rulebooleanLiteral, position1592)
			}
			return true
		l1591:
			position, tokenIndex, depth = position1591, tokenIndex1591, depth1591
			return false
		},
		/* 78 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 79 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 80 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 81 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1618, tokenIndex1618, depth1618 := position, tokenIndex, depth
			{
				position1619 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1618
				}
				position++
			l1620:
				{
					position1621, tokenIndex1621, depth1621 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1621
					}
					goto l1620
				l1621:
					position, tokenIndex, depth = position1621, tokenIndex1621, depth1621
				}
				if buffer[position] != rune(')') {
					goto l1618
				}
				position++
				if !rules[ruleskip]() {
					goto l1618
				}
				depth--
				add(rulenil, position1619)
			}
			return true
		l1618:
			position, tokenIndex, depth = position1618, tokenIndex1618, depth1618
			return false
		},
		/* 82 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 83 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 84 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 85 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 86 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 87 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 88 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 89 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 90 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 91 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
			{
				position1632 := position
				depth++
				{
					position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1634
					}
					position++
					goto l1633
				l1634:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					if buffer[position] != rune('D') {
						goto l1631
					}
					position++
				}
			l1633:
				{
					position1635, tokenIndex1635, depth1635 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1636
					}
					position++
					goto l1635
				l1636:
					position, tokenIndex, depth = position1635, tokenIndex1635, depth1635
					if buffer[position] != rune('I') {
						goto l1631
					}
					position++
				}
			l1635:
				{
					position1637, tokenIndex1637, depth1637 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1638
					}
					position++
					goto l1637
				l1638:
					position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
					if buffer[position] != rune('S') {
						goto l1631
					}
					position++
				}
			l1637:
				{
					position1639, tokenIndex1639, depth1639 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1640
					}
					position++
					goto l1639
				l1640:
					position, tokenIndex, depth = position1639, tokenIndex1639, depth1639
					if buffer[position] != rune('T') {
						goto l1631
					}
					position++
				}
			l1639:
				{
					position1641, tokenIndex1641, depth1641 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1642
					}
					position++
					goto l1641
				l1642:
					position, tokenIndex, depth = position1641, tokenIndex1641, depth1641
					if buffer[position] != rune('I') {
						goto l1631
					}
					position++
				}
			l1641:
				{
					position1643, tokenIndex1643, depth1643 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1644
					}
					position++
					goto l1643
				l1644:
					position, tokenIndex, depth = position1643, tokenIndex1643, depth1643
					if buffer[position] != rune('N') {
						goto l1631
					}
					position++
				}
			l1643:
				{
					position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1646
					}
					position++
					goto l1645
				l1646:
					position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
					if buffer[position] != rune('C') {
						goto l1631
					}
					position++
				}
			l1645:
				{
					position1647, tokenIndex1647, depth1647 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1648
					}
					position++
					goto l1647
				l1648:
					position, tokenIndex, depth = position1647, tokenIndex1647, depth1647
					if buffer[position] != rune('T') {
						goto l1631
					}
					position++
				}
			l1647:
				if !rules[ruleskip]() {
					goto l1631
				}
				depth--
				add(ruleDISTINCT, position1632)
			}
			return true
		l1631:
			position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
			return false
		},
		/* 92 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 93 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 94 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 95 LBRACE <- <('{' skip)> */
		func() bool {
			position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
			{
				position1653 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1652
				}
				position++
				if !rules[ruleskip]() {
					goto l1652
				}
				depth--
				add(ruleLBRACE, position1653)
			}
			return true
		l1652:
			position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
			return false
		},
		/* 96 RBRACE <- <('}' skip)> */
		func() bool {
			position1654, tokenIndex1654, depth1654 := position, tokenIndex, depth
			{
				position1655 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1654
				}
				position++
				if !rules[ruleskip]() {
					goto l1654
				}
				depth--
				add(ruleRBRACE, position1655)
			}
			return true
		l1654:
			position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
			return false
		},
		/* 97 LBRACK <- <('[' skip)> */
		nil,
		/* 98 RBRACK <- <(']' skip)> */
		nil,
		/* 99 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
			{
				position1659 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1658
				}
				position++
				if !rules[ruleskip]() {
					goto l1658
				}
				depth--
				add(ruleSEMICOLON, position1659)
			}
			return true
		l1658:
			position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
			return false
		},
		/* 100 COMMA <- <(',' skip)> */
		func() bool {
			position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
			{
				position1661 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1660
				}
				position++
				if !rules[ruleskip]() {
					goto l1660
				}
				depth--
				add(ruleCOMMA, position1661)
			}
			return true
		l1660:
			position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
			return false
		},
		/* 101 DOT <- <('.' skip)> */
		func() bool {
			position1662, tokenIndex1662, depth1662 := position, tokenIndex, depth
			{
				position1663 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1662
				}
				position++
				if !rules[ruleskip]() {
					goto l1662
				}
				depth--
				add(ruleDOT, position1663)
			}
			return true
		l1662:
			position, tokenIndex, depth = position1662, tokenIndex1662, depth1662
			return false
		},
		/* 102 COLON <- <(':' skip)> */
		nil,
		/* 103 PIPE <- <('|' skip)> */
		func() bool {
			position1665, tokenIndex1665, depth1665 := position, tokenIndex, depth
			{
				position1666 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1665
				}
				position++
				if !rules[ruleskip]() {
					goto l1665
				}
				depth--
				add(rulePIPE, position1666)
			}
			return true
		l1665:
			position, tokenIndex, depth = position1665, tokenIndex1665, depth1665
			return false
		},
		/* 104 SLASH <- <('/' skip)> */
		func() bool {
			position1667, tokenIndex1667, depth1667 := position, tokenIndex, depth
			{
				position1668 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1667
				}
				position++
				if !rules[ruleskip]() {
					goto l1667
				}
				depth--
				add(ruleSLASH, position1668)
			}
			return true
		l1667:
			position, tokenIndex, depth = position1667, tokenIndex1667, depth1667
			return false
		},
		/* 105 INVERSE <- <('^' skip)> */
		func() bool {
			position1669, tokenIndex1669, depth1669 := position, tokenIndex, depth
			{
				position1670 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1669
				}
				position++
				if !rules[ruleskip]() {
					goto l1669
				}
				depth--
				add(ruleINVERSE, position1670)
			}
			return true
		l1669:
			position, tokenIndex, depth = position1669, tokenIndex1669, depth1669
			return false
		},
		/* 106 LPAREN <- <('(' skip)> */
		func() bool {
			position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
			{
				position1672 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1671
				}
				position++
				if !rules[ruleskip]() {
					goto l1671
				}
				depth--
				add(ruleLPAREN, position1672)
			}
			return true
		l1671:
			position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
			return false
		},
		/* 107 RPAREN <- <(')' skip)> */
		func() bool {
			position1673, tokenIndex1673, depth1673 := position, tokenIndex, depth
			{
				position1674 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1673
				}
				position++
				if !rules[ruleskip]() {
					goto l1673
				}
				depth--
				add(ruleRPAREN, position1674)
			}
			return true
		l1673:
			position, tokenIndex, depth = position1673, tokenIndex1673, depth1673
			return false
		},
		/* 108 ISA <- <('a' skip)> */
		func() bool {
			position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
			{
				position1676 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1675
				}
				position++
				if !rules[ruleskip]() {
					goto l1675
				}
				depth--
				add(ruleISA, position1676)
			}
			return true
		l1675:
			position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
			return false
		},
		/* 109 NOT <- <('!' skip)> */
		func() bool {
			position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
			{
				position1678 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1677
				}
				position++
				if !rules[ruleskip]() {
					goto l1677
				}
				depth--
				add(ruleNOT, position1678)
			}
			return true
		l1677:
			position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
			return false
		},
		/* 110 STAR <- <('*' skip)> */
		func() bool {
			position1679, tokenIndex1679, depth1679 := position, tokenIndex, depth
			{
				position1680 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1679
				}
				position++
				if !rules[ruleskip]() {
					goto l1679
				}
				depth--
				add(ruleSTAR, position1680)
			}
			return true
		l1679:
			position, tokenIndex, depth = position1679, tokenIndex1679, depth1679
			return false
		},
		/* 111 QUESTION <- <('?' skip)> */
		nil,
		/* 112 PLUS <- <('+' skip)> */
		func() bool {
			position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
			{
				position1683 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1682
				}
				position++
				if !rules[ruleskip]() {
					goto l1682
				}
				depth--
				add(rulePLUS, position1683)
			}
			return true
		l1682:
			position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
			return false
		},
		/* 113 MINUS <- <('-' skip)> */
		func() bool {
			position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
			{
				position1685 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1684
				}
				position++
				if !rules[ruleskip]() {
					goto l1684
				}
				depth--
				add(ruleMINUS, position1685)
			}
			return true
		l1684:
			position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
			return false
		},
		/* 114 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 115 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 116 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 117 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 118 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1690, tokenIndex1690, depth1690 := position, tokenIndex, depth
			{
				position1691 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1690
				}
				position++
			l1692:
				{
					position1693, tokenIndex1693, depth1693 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1693
					}
					position++
					goto l1692
				l1693:
					position, tokenIndex, depth = position1693, tokenIndex1693, depth1693
				}
				if !rules[ruleskip]() {
					goto l1690
				}
				depth--
				add(ruleINTEGER, position1691)
			}
			return true
		l1690:
			position, tokenIndex, depth = position1690, tokenIndex1690, depth1690
			return false
		},
		/* 119 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 120 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 121 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 122 OR <- <('|' '|' skip)> */
		nil,
		/* 123 AND <- <('&' '&' skip)> */
		nil,
		/* 124 EQ <- <('=' skip)> */
		func() bool {
			position1699, tokenIndex1699, depth1699 := position, tokenIndex, depth
			{
				position1700 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1699
				}
				position++
				if !rules[ruleskip]() {
					goto l1699
				}
				depth--
				add(ruleEQ, position1700)
			}
			return true
		l1699:
			position, tokenIndex, depth = position1699, tokenIndex1699, depth1699
			return false
		},
		/* 125 NE <- <('!' '=' skip)> */
		nil,
		/* 126 GT <- <('>' skip)> */
		nil,
		/* 127 LT <- <('<' skip)> */
		nil,
		/* 128 LE <- <('<' '=' skip)> */
		nil,
		/* 129 GE <- <('>' '=' skip)> */
		nil,
		/* 130 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 131 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 132 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
			{
				position1709 := position
				depth++
				{
					position1710, tokenIndex1710, depth1710 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1711
					}
					position++
					goto l1710
				l1711:
					position, tokenIndex, depth = position1710, tokenIndex1710, depth1710
					if buffer[position] != rune('A') {
						goto l1708
					}
					position++
				}
			l1710:
				{
					position1712, tokenIndex1712, depth1712 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1713
					}
					position++
					goto l1712
				l1713:
					position, tokenIndex, depth = position1712, tokenIndex1712, depth1712
					if buffer[position] != rune('S') {
						goto l1708
					}
					position++
				}
			l1712:
				if !rules[ruleskip]() {
					goto l1708
				}
				depth--
				add(ruleAS, position1709)
			}
			return true
		l1708:
			position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
			return false
		},
		/* 133 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 134 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 135 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 136 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 137 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 138 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 139 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 140 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 141 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 142 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 143 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 144 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 145 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 146 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 147 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 148 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 149 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 150 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 151 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 152 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 153 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 154 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 155 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 156 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 157 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 158 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 159 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 160 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 161 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 162 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 163 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 164 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 165 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 166 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 167 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 168 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 169 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 170 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 171 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 172 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 173 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 174 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 175 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 176 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 177 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 178 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 179 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 180 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 181 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 182 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 183 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 184 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 185 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 186 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 187 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 188 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 189 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 190 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 191 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 192 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 193 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 194 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 195 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 196 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 197 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 198 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 199 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 200 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 201 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1782, tokenIndex1782, depth1782 := position, tokenIndex, depth
			{
				position1783 := position
				depth++
				{
					position1784, tokenIndex1784, depth1784 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1785
					}
					position++
					goto l1784
				l1785:
					position, tokenIndex, depth = position1784, tokenIndex1784, depth1784
					if buffer[position] != rune('B') {
						goto l1782
					}
					position++
				}
			l1784:
				{
					position1786, tokenIndex1786, depth1786 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1787
					}
					position++
					goto l1786
				l1787:
					position, tokenIndex, depth = position1786, tokenIndex1786, depth1786
					if buffer[position] != rune('Y') {
						goto l1782
					}
					position++
				}
			l1786:
				if !rules[ruleskip]() {
					goto l1782
				}
				depth--
				add(ruleBY, position1783)
			}
			return true
		l1782:
			position, tokenIndex, depth = position1782, tokenIndex1782, depth1782
			return false
		},
		/* 202 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 203 GRAPH <- <(('g' / 'G') ('r' / 'R') ('a' / 'A') ('p' / 'P') ('h' / 'H') skip)> */
		nil,
		/* 204 MINUSSETOPER <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('s' / 'S') skip)> */
		nil,
		/* 205 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1792 := position
				depth++
			l1793:
				{
					position1794, tokenIndex1794, depth1794 := position, tokenIndex, depth
					{
						position1795, tokenIndex1795, depth1795 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1796
						}
						goto l1795
					l1796:
						position, tokenIndex, depth = position1795, tokenIndex1795, depth1795
						{
							position1797 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1794
							}
							position++
						l1798:
							{
								position1799, tokenIndex1799, depth1799 := position, tokenIndex, depth
								{
									position1800, tokenIndex1800, depth1800 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1800
									}
									goto l1799
								l1800:
									position, tokenIndex, depth = position1800, tokenIndex1800, depth1800
								}
								if !matchDot() {
									goto l1799
								}
								goto l1798
							l1799:
								position, tokenIndex, depth = position1799, tokenIndex1799, depth1799
							}
							if !rules[ruleendOfLine]() {
								goto l1794
							}
							depth--
							add(rulecomment, position1797)
						}
					}
				l1795:
					goto l1793
				l1794:
					position, tokenIndex, depth = position1794, tokenIndex1794, depth1794
				}
				depth--
				add(ruleskip, position1792)
			}
			return true
		},
		/* 206 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1801, tokenIndex1801, depth1801 := position, tokenIndex, depth
			{
				position1802 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1801
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1801
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1801
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1801
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1801
						}
						break
					}
				}

				depth--
				add(rulews, position1802)
			}
			return true
		l1801:
			position, tokenIndex, depth = position1801, tokenIndex1801, depth1801
			return false
		},
		/* 207 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 208 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1805, tokenIndex1805, depth1805 := position, tokenIndex, depth
			{
				position1806 := position
				depth++
				{
					position1807, tokenIndex1807, depth1807 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1808
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1808
					}
					position++
					goto l1807
				l1808:
					position, tokenIndex, depth = position1807, tokenIndex1807, depth1807
					if buffer[position] != rune('\n') {
						goto l1809
					}
					position++
					goto l1807
				l1809:
					position, tokenIndex, depth = position1807, tokenIndex1807, depth1807
					if buffer[position] != rune('\r') {
						goto l1805
					}
					position++
				}
			l1807:
				depth--
				add(ruleendOfLine, position1806)
			}
			return true
		l1805:
			position, tokenIndex, depth = position1805, tokenIndex1805, depth1805
			return false
		},
	}
	p.rules = rules
}
