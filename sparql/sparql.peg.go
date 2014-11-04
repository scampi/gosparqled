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
	ruleserviceGraphPattern
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
	rulestringLiteralA
	rulestringLiteralB
	rulestringLiteralLongA
	rulestringLiteralLongB
	ruleechar
	rulenumericLiteral
	rulesignedNumericLiteral
	rulebooleanLiteral
	ruleblankNode
	ruleblankNodeLabel
	ruleanon
	rulenil
	ruleVARNAME
	rulepnPrefix
	rulepnLocal
	rulepnChars
	rulepnCharsU
	rulepnCharsBase
	ruleplx
	rulepercent
	rulehex
	rulepnLocalEsc
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
	ruleSERVICE
	ruleSILENT

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
	"serviceGraphPattern",
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
	"stringLiteralA",
	"stringLiteralB",
	"stringLiteralLongA",
	"stringLiteralLongB",
	"echar",
	"numericLiteral",
	"signedNumericLiteral",
	"booleanLiteral",
	"blankNode",
	"blankNodeLabel",
	"anon",
	"nil",
	"VARNAME",
	"pnPrefix",
	"pnLocal",
	"pnChars",
	"pnCharsU",
	"pnCharsBase",
	"plx",
	"percent",
	"hex",
	"pnLocalEsc",
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
	"SERVICE",
	"SILENT",

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
	rules  [225]func() bool
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
								{
									position21, tokenIndex21, depth21 := position, tokenIndex, depth
									if !rules[rulepnPrefix]() {
										goto l21
									}
									goto l22
								l21:
									position, tokenIndex, depth = position21, tokenIndex21, depth21
								}
							l22:
								{
									position23 := position
									depth++
									if buffer[position] != rune(':') {
										goto l6
									}
									position++
									if !rules[ruleskip]() {
										goto l6
									}
									depth--
									add(ruleCOLON, position23)
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
								position24 := position
								depth++
								{
									position25 := position
									depth++
									{
										position26, tokenIndex26, depth26 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l27
										}
										position++
										goto l26
									l27:
										position, tokenIndex, depth = position26, tokenIndex26, depth26
										if buffer[position] != rune('B') {
											goto l4
										}
										position++
									}
								l26:
									{
										position28, tokenIndex28, depth28 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l29
										}
										position++
										goto l28
									l29:
										position, tokenIndex, depth = position28, tokenIndex28, depth28
										if buffer[position] != rune('A') {
											goto l4
										}
										position++
									}
								l28:
									{
										position30, tokenIndex30, depth30 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l31
										}
										position++
										goto l30
									l31:
										position, tokenIndex, depth = position30, tokenIndex30, depth30
										if buffer[position] != rune('S') {
											goto l4
										}
										position++
									}
								l30:
									{
										position32, tokenIndex32, depth32 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l33
										}
										position++
										goto l32
									l33:
										position, tokenIndex, depth = position32, tokenIndex32, depth32
										if buffer[position] != rune('E') {
											goto l4
										}
										position++
									}
								l32:
									if !rules[ruleskip]() {
										goto l4
									}
									depth--
									add(ruleBASE, position25)
								}
								if !rules[ruleiri]() {
									goto l4
								}
								depth--
								add(rulebaseDecl, position24)
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
					position34 := position
					depth++
					{
						switch buffer[position] {
						case 'A', 'a':
							{
								position36 := position
								depth++
								{
									position37 := position
									depth++
									{
										position38, tokenIndex38, depth38 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l39
										}
										position++
										goto l38
									l39:
										position, tokenIndex, depth = position38, tokenIndex38, depth38
										if buffer[position] != rune('A') {
											goto l0
										}
										position++
									}
								l38:
									{
										position40, tokenIndex40, depth40 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l41
										}
										position++
										goto l40
									l41:
										position, tokenIndex, depth = position40, tokenIndex40, depth40
										if buffer[position] != rune('S') {
											goto l0
										}
										position++
									}
								l40:
									{
										position42, tokenIndex42, depth42 := position, tokenIndex, depth
										if buffer[position] != rune('k') {
											goto l43
										}
										position++
										goto l42
									l43:
										position, tokenIndex, depth = position42, tokenIndex42, depth42
										if buffer[position] != rune('K') {
											goto l0
										}
										position++
									}
								l42:
									if !rules[ruleskip]() {
										goto l0
									}
									depth--
									add(ruleASK, position37)
								}
							l44:
								{
									position45, tokenIndex45, depth45 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l45
									}
									goto l44
								l45:
									position, tokenIndex, depth = position45, tokenIndex45, depth45
								}
								if !rules[rulewhereClause]() {
									goto l0
								}
								depth--
								add(ruleaskQuery, position36)
							}
							break
						case 'D', 'd':
							{
								position46 := position
								depth++
								{
									position47 := position
									depth++
									{
										position48 := position
										depth++
										{
											position49, tokenIndex49, depth49 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l50
											}
											position++
											goto l49
										l50:
											position, tokenIndex, depth = position49, tokenIndex49, depth49
											if buffer[position] != rune('D') {
												goto l0
											}
											position++
										}
									l49:
										{
											position51, tokenIndex51, depth51 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l52
											}
											position++
											goto l51
										l52:
											position, tokenIndex, depth = position51, tokenIndex51, depth51
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l51:
										{
											position53, tokenIndex53, depth53 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l54
											}
											position++
											goto l53
										l54:
											position, tokenIndex, depth = position53, tokenIndex53, depth53
											if buffer[position] != rune('S') {
												goto l0
											}
											position++
										}
									l53:
										{
											position55, tokenIndex55, depth55 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l56
											}
											position++
											goto l55
										l56:
											position, tokenIndex, depth = position55, tokenIndex55, depth55
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l55:
										{
											position57, tokenIndex57, depth57 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l58
											}
											position++
											goto l57
										l58:
											position, tokenIndex, depth = position57, tokenIndex57, depth57
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l57:
										{
											position59, tokenIndex59, depth59 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l60
											}
											position++
											goto l59
										l60:
											position, tokenIndex, depth = position59, tokenIndex59, depth59
											if buffer[position] != rune('I') {
												goto l0
											}
											position++
										}
									l59:
										{
											position61, tokenIndex61, depth61 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l62
											}
											position++
											goto l61
										l62:
											position, tokenIndex, depth = position61, tokenIndex61, depth61
											if buffer[position] != rune('B') {
												goto l0
											}
											position++
										}
									l61:
										{
											position63, tokenIndex63, depth63 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l64
											}
											position++
											goto l63
										l64:
											position, tokenIndex, depth = position63, tokenIndex63, depth63
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l63:
										if !rules[ruleskip]() {
											goto l0
										}
										depth--
										add(ruleDESCRIBE, position48)
									}
									{
										switch buffer[position] {
										case '*':
											if !rules[ruleSTAR]() {
												goto l0
											}
											break
										case '$', '?':
											if !rules[rulevar]() {
												goto l0
											}
											break
										default:
											if !rules[ruleiriref]() {
												goto l0
											}
											break
										}
									}

									depth--
									add(ruledescribe, position47)
								}
							l66:
								{
									position67, tokenIndex67, depth67 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l67
									}
									goto l66
								l67:
									position, tokenIndex, depth = position67, tokenIndex67, depth67
								}
								{
									position68, tokenIndex68, depth68 := position, tokenIndex, depth
									if !rules[rulewhereClause]() {
										goto l68
									}
									goto l69
								l68:
									position, tokenIndex, depth = position68, tokenIndex68, depth68
								}
							l69:
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruledescribeQuery, position46)
							}
							break
						case 'C', 'c':
							{
								position70 := position
								depth++
								{
									position71 := position
									depth++
									{
										position72 := position
										depth++
										{
											position73, tokenIndex73, depth73 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l74
											}
											position++
											goto l73
										l74:
											position, tokenIndex, depth = position73, tokenIndex73, depth73
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l73:
										{
											position75, tokenIndex75, depth75 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l76
											}
											position++
											goto l75
										l76:
											position, tokenIndex, depth = position75, tokenIndex75, depth75
											if buffer[position] != rune('O') {
												goto l0
											}
											position++
										}
									l75:
										{
											position77, tokenIndex77, depth77 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l78
											}
											position++
											goto l77
										l78:
											position, tokenIndex, depth = position77, tokenIndex77, depth77
											if buffer[position] != rune('N') {
												goto l0
											}
											position++
										}
									l77:
										{
											position79, tokenIndex79, depth79 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l80
											}
											position++
											goto l79
										l80:
											position, tokenIndex, depth = position79, tokenIndex79, depth79
											if buffer[position] != rune('S') {
												goto l0
											}
											position++
										}
									l79:
										{
											position81, tokenIndex81, depth81 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l82
											}
											position++
											goto l81
										l82:
											position, tokenIndex, depth = position81, tokenIndex81, depth81
											if buffer[position] != rune('T') {
												goto l0
											}
											position++
										}
									l81:
										{
											position83, tokenIndex83, depth83 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l84
											}
											position++
											goto l83
										l84:
											position, tokenIndex, depth = position83, tokenIndex83, depth83
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l83:
										{
											position85, tokenIndex85, depth85 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l86
											}
											position++
											goto l85
										l86:
											position, tokenIndex, depth = position85, tokenIndex85, depth85
											if buffer[position] != rune('U') {
												goto l0
											}
											position++
										}
									l85:
										{
											position87, tokenIndex87, depth87 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l88
											}
											position++
											goto l87
										l88:
											position, tokenIndex, depth = position87, tokenIndex87, depth87
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l87:
										{
											position89, tokenIndex89, depth89 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l90
											}
											position++
											goto l89
										l90:
											position, tokenIndex, depth = position89, tokenIndex89, depth89
											if buffer[position] != rune('T') {
												goto l0
											}
											position++
										}
									l89:
										if !rules[ruleskip]() {
											goto l0
										}
										depth--
										add(ruleCONSTRUCT, position72)
									}
									if !rules[ruleLBRACE]() {
										goto l0
									}
									{
										position91, tokenIndex91, depth91 := position, tokenIndex, depth
										if !rules[ruletriplesBlock]() {
											goto l91
										}
										goto l92
									l91:
										position, tokenIndex, depth = position91, tokenIndex91, depth91
									}
								l92:
									if !rules[ruleRBRACE]() {
										goto l0
									}
									depth--
									add(ruleconstruct, position71)
								}
							l93:
								{
									position94, tokenIndex94, depth94 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l94
									}
									goto l93
								l94:
									position, tokenIndex, depth = position94, tokenIndex94, depth94
								}
								if !rules[rulewhereClause]() {
									goto l0
								}
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruleconstructQuery, position70)
							}
							break
						default:
							{
								position95 := position
								depth++
								if !rules[ruleselect]() {
									goto l0
								}
							l96:
								{
									position97, tokenIndex97, depth97 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l97
									}
									goto l96
								l97:
									position, tokenIndex, depth = position97, tokenIndex97, depth97
								}
								if !rules[rulewhereClause]() {
									goto l0
								}
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruleselectQuery, position95)
							}
							break
						}
					}

					depth--
					add(rulequery, position34)
				}
				{
					position98, tokenIndex98, depth98 := position, tokenIndex, depth
					if !matchDot() {
						goto l98
					}
					goto l0
				l98:
					position, tokenIndex, depth = position98, tokenIndex98, depth98
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
		/* 2 prefixDecl <- <(PREFIX pnPrefix? COLON iri)> */
		nil,
		/* 3 baseDecl <- <(BASE iri)> */
		nil,
		/* 4 query <- <((&('A' | 'a') askQuery) | (&('D' | 'd') describeQuery) | (&('C' | 'c') constructQuery) | (&('S' | 's') selectQuery))> */
		nil,
		/* 5 selectQuery <- <(select datasetClause* whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position104, tokenIndex104, depth104 := position, tokenIndex, depth
			{
				position105 := position
				depth++
				{
					position106 := position
					depth++
					{
						position107, tokenIndex107, depth107 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l108
						}
						position++
						goto l107
					l108:
						position, tokenIndex, depth = position107, tokenIndex107, depth107
						if buffer[position] != rune('S') {
							goto l104
						}
						position++
					}
				l107:
					{
						position109, tokenIndex109, depth109 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l110
						}
						position++
						goto l109
					l110:
						position, tokenIndex, depth = position109, tokenIndex109, depth109
						if buffer[position] != rune('E') {
							goto l104
						}
						position++
					}
				l109:
					{
						position111, tokenIndex111, depth111 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l112
						}
						position++
						goto l111
					l112:
						position, tokenIndex, depth = position111, tokenIndex111, depth111
						if buffer[position] != rune('L') {
							goto l104
						}
						position++
					}
				l111:
					{
						position113, tokenIndex113, depth113 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l114
						}
						position++
						goto l113
					l114:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
						if buffer[position] != rune('E') {
							goto l104
						}
						position++
					}
				l113:
					{
						position115, tokenIndex115, depth115 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l116
						}
						position++
						goto l115
					l116:
						position, tokenIndex, depth = position115, tokenIndex115, depth115
						if buffer[position] != rune('C') {
							goto l104
						}
						position++
					}
				l115:
					{
						position117, tokenIndex117, depth117 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l118
						}
						position++
						goto l117
					l118:
						position, tokenIndex, depth = position117, tokenIndex117, depth117
						if buffer[position] != rune('T') {
							goto l104
						}
						position++
					}
				l117:
					if !rules[ruleskip]() {
						goto l104
					}
					depth--
					add(ruleSELECT, position106)
				}
				{
					position119, tokenIndex119, depth119 := position, tokenIndex, depth
					{
						position121, tokenIndex121, depth121 := position, tokenIndex, depth
						if !rules[ruleDISTINCT]() {
							goto l122
						}
						goto l121
					l122:
						position, tokenIndex, depth = position121, tokenIndex121, depth121
						{
							position123 := position
							depth++
							{
								position124, tokenIndex124, depth124 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l125
								}
								position++
								goto l124
							l125:
								position, tokenIndex, depth = position124, tokenIndex124, depth124
								if buffer[position] != rune('R') {
									goto l119
								}
								position++
							}
						l124:
							{
								position126, tokenIndex126, depth126 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l127
								}
								position++
								goto l126
							l127:
								position, tokenIndex, depth = position126, tokenIndex126, depth126
								if buffer[position] != rune('E') {
									goto l119
								}
								position++
							}
						l126:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l129
								}
								position++
								goto l128
							l129:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
								if buffer[position] != rune('D') {
									goto l119
								}
								position++
							}
						l128:
							{
								position130, tokenIndex130, depth130 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l131
								}
								position++
								goto l130
							l131:
								position, tokenIndex, depth = position130, tokenIndex130, depth130
								if buffer[position] != rune('U') {
									goto l119
								}
								position++
							}
						l130:
							{
								position132, tokenIndex132, depth132 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l133
								}
								position++
								goto l132
							l133:
								position, tokenIndex, depth = position132, tokenIndex132, depth132
								if buffer[position] != rune('C') {
									goto l119
								}
								position++
							}
						l132:
							{
								position134, tokenIndex134, depth134 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l135
								}
								position++
								goto l134
							l135:
								position, tokenIndex, depth = position134, tokenIndex134, depth134
								if buffer[position] != rune('E') {
									goto l119
								}
								position++
							}
						l134:
							{
								position136, tokenIndex136, depth136 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l137
								}
								position++
								goto l136
							l137:
								position, tokenIndex, depth = position136, tokenIndex136, depth136
								if buffer[position] != rune('D') {
									goto l119
								}
								position++
							}
						l136:
							if !rules[ruleskip]() {
								goto l119
							}
							depth--
							add(ruleREDUCED, position123)
						}
					}
				l121:
					goto l120
				l119:
					position, tokenIndex, depth = position119, tokenIndex119, depth119
				}
			l120:
				{
					position138, tokenIndex138, depth138 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l139
					}
					goto l138
				l139:
					position, tokenIndex, depth = position138, tokenIndex138, depth138
					{
						position142 := position
						depth++
						{
							position143, tokenIndex143, depth143 := position, tokenIndex, depth
							if !rules[rulevar]() {
								goto l144
							}
							goto l143
						l144:
							position, tokenIndex, depth = position143, tokenIndex143, depth143
							if !rules[ruleLPAREN]() {
								goto l104
							}
							if !rules[ruleexpression]() {
								goto l104
							}
							if !rules[ruleAS]() {
								goto l104
							}
							if !rules[rulevar]() {
								goto l104
							}
							if !rules[ruleRPAREN]() {
								goto l104
							}
						}
					l143:
						depth--
						add(ruleprojectionElem, position142)
					}
				l140:
					{
						position141, tokenIndex141, depth141 := position, tokenIndex, depth
						{
							position145 := position
							depth++
							{
								position146, tokenIndex146, depth146 := position, tokenIndex, depth
								if !rules[rulevar]() {
									goto l147
								}
								goto l146
							l147:
								position, tokenIndex, depth = position146, tokenIndex146, depth146
								if !rules[ruleLPAREN]() {
									goto l141
								}
								if !rules[ruleexpression]() {
									goto l141
								}
								if !rules[ruleAS]() {
									goto l141
								}
								if !rules[rulevar]() {
									goto l141
								}
								if !rules[ruleRPAREN]() {
									goto l141
								}
							}
						l146:
							depth--
							add(ruleprojectionElem, position145)
						}
						goto l140
					l141:
						position, tokenIndex, depth = position141, tokenIndex141, depth141
					}
				}
			l138:
				depth--
				add(ruleselect, position105)
			}
			return true
		l104:
			position, tokenIndex, depth = position104, tokenIndex104, depth104
			return false
		},
		/* 7 subSelect <- <(select whereClause solutionModifier)> */
		func() bool {
			position148, tokenIndex148, depth148 := position, tokenIndex, depth
			{
				position149 := position
				depth++
				if !rules[ruleselect]() {
					goto l148
				}
				if !rules[rulewhereClause]() {
					goto l148
				}
				if !rules[rulesolutionModifier]() {
					goto l148
				}
				depth--
				add(rulesubSelect, position149)
			}
			return true
		l148:
			position, tokenIndex, depth = position148, tokenIndex148, depth148
			return false
		},
		/* 8 constructQuery <- <(construct datasetClause* whereClause solutionModifier)> */
		nil,
		/* 9 construct <- <(CONSTRUCT LBRACE triplesBlock? RBRACE)> */
		nil,
		/* 10 describeQuery <- <(describe datasetClause* whereClause? solutionModifier)> */
		nil,
		/* 11 describe <- <(DESCRIBE ((&('*') STAR) | (&('$' | '?') var) | (&(':' | '<' | 'A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') iriref)))> */
		nil,
		/* 12 askQuery <- <(ASK datasetClause* whereClause)> */
		nil,
		/* 13 projectionElem <- <(var / (LPAREN expression AS var RPAREN))> */
		nil,
		/* 14 datasetClause <- <(FROM NAMED? iriref)> */
		func() bool {
			position156, tokenIndex156, depth156 := position, tokenIndex, depth
			{
				position157 := position
				depth++
				{
					position158 := position
					depth++
					{
						position159, tokenIndex159, depth159 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l160
						}
						position++
						goto l159
					l160:
						position, tokenIndex, depth = position159, tokenIndex159, depth159
						if buffer[position] != rune('F') {
							goto l156
						}
						position++
					}
				l159:
					{
						position161, tokenIndex161, depth161 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l162
						}
						position++
						goto l161
					l162:
						position, tokenIndex, depth = position161, tokenIndex161, depth161
						if buffer[position] != rune('R') {
							goto l156
						}
						position++
					}
				l161:
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l164
						}
						position++
						goto l163
					l164:
						position, tokenIndex, depth = position163, tokenIndex163, depth163
						if buffer[position] != rune('O') {
							goto l156
						}
						position++
					}
				l163:
					{
						position165, tokenIndex165, depth165 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l166
						}
						position++
						goto l165
					l166:
						position, tokenIndex, depth = position165, tokenIndex165, depth165
						if buffer[position] != rune('M') {
							goto l156
						}
						position++
					}
				l165:
					if !rules[ruleskip]() {
						goto l156
					}
					depth--
					add(ruleFROM, position158)
				}
				{
					position167, tokenIndex167, depth167 := position, tokenIndex, depth
					{
						position169 := position
						depth++
						{
							position170, tokenIndex170, depth170 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l171
							}
							position++
							goto l170
						l171:
							position, tokenIndex, depth = position170, tokenIndex170, depth170
							if buffer[position] != rune('N') {
								goto l167
							}
							position++
						}
					l170:
						{
							position172, tokenIndex172, depth172 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l173
							}
							position++
							goto l172
						l173:
							position, tokenIndex, depth = position172, tokenIndex172, depth172
							if buffer[position] != rune('A') {
								goto l167
							}
							position++
						}
					l172:
						{
							position174, tokenIndex174, depth174 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l175
							}
							position++
							goto l174
						l175:
							position, tokenIndex, depth = position174, tokenIndex174, depth174
							if buffer[position] != rune('M') {
								goto l167
							}
							position++
						}
					l174:
						{
							position176, tokenIndex176, depth176 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l177
							}
							position++
							goto l176
						l177:
							position, tokenIndex, depth = position176, tokenIndex176, depth176
							if buffer[position] != rune('E') {
								goto l167
							}
							position++
						}
					l176:
						{
							position178, tokenIndex178, depth178 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l179
							}
							position++
							goto l178
						l179:
							position, tokenIndex, depth = position178, tokenIndex178, depth178
							if buffer[position] != rune('D') {
								goto l167
							}
							position++
						}
					l178:
						if !rules[ruleskip]() {
							goto l167
						}
						depth--
						add(ruleNAMED, position169)
					}
					goto l168
				l167:
					position, tokenIndex, depth = position167, tokenIndex167, depth167
				}
			l168:
				if !rules[ruleiriref]() {
					goto l156
				}
				depth--
				add(ruledatasetClause, position157)
			}
			return true
		l156:
			position, tokenIndex, depth = position156, tokenIndex156, depth156
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position180, tokenIndex180, depth180 := position, tokenIndex, depth
			{
				position181 := position
				depth++
				{
					position182, tokenIndex182, depth182 := position, tokenIndex, depth
					{
						position184 := position
						depth++
						{
							position185, tokenIndex185, depth185 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l186
							}
							position++
							goto l185
						l186:
							position, tokenIndex, depth = position185, tokenIndex185, depth185
							if buffer[position] != rune('W') {
								goto l182
							}
							position++
						}
					l185:
						{
							position187, tokenIndex187, depth187 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l188
							}
							position++
							goto l187
						l188:
							position, tokenIndex, depth = position187, tokenIndex187, depth187
							if buffer[position] != rune('H') {
								goto l182
							}
							position++
						}
					l187:
						{
							position189, tokenIndex189, depth189 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l190
							}
							position++
							goto l189
						l190:
							position, tokenIndex, depth = position189, tokenIndex189, depth189
							if buffer[position] != rune('E') {
								goto l182
							}
							position++
						}
					l189:
						{
							position191, tokenIndex191, depth191 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l192
							}
							position++
							goto l191
						l192:
							position, tokenIndex, depth = position191, tokenIndex191, depth191
							if buffer[position] != rune('R') {
								goto l182
							}
							position++
						}
					l191:
						{
							position193, tokenIndex193, depth193 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l194
							}
							position++
							goto l193
						l194:
							position, tokenIndex, depth = position193, tokenIndex193, depth193
							if buffer[position] != rune('E') {
								goto l182
							}
							position++
						}
					l193:
						if !rules[ruleskip]() {
							goto l182
						}
						depth--
						add(ruleWHERE, position184)
					}
					goto l183
				l182:
					position, tokenIndex, depth = position182, tokenIndex182, depth182
				}
			l183:
				if !rules[rulegroupGraphPattern]() {
					goto l180
				}
				depth--
				add(rulewhereClause, position181)
			}
			return true
		l180:
			position, tokenIndex, depth = position180, tokenIndex180, depth180
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position195, tokenIndex195, depth195 := position, tokenIndex, depth
			{
				position196 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l195
				}
				{
					position197, tokenIndex197, depth197 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l198
					}
					goto l197
				l198:
					position, tokenIndex, depth = position197, tokenIndex197, depth197
					if !rules[rulegraphPattern]() {
						goto l195
					}
				}
			l197:
				if !rules[ruleRBRACE]() {
					goto l195
				}
				depth--
				add(rulegroupGraphPattern, position196)
			}
			return true
		l195:
			position, tokenIndex, depth = position195, tokenIndex195, depth195
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position200 := position
				depth++
				{
					position201, tokenIndex201, depth201 := position, tokenIndex, depth
					{
						position203 := position
						depth++
						{
							position204, tokenIndex204, depth204 := position, tokenIndex, depth
							if !rules[ruletriplesBlock]() {
								goto l205
							}
						l206:
							{
								position207, tokenIndex207, depth207 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l207
								}
								{
									position208, tokenIndex208, depth208 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l208
									}
									goto l209
								l208:
									position, tokenIndex, depth = position208, tokenIndex208, depth208
								}
							l209:
								{
									position210, tokenIndex210, depth210 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l210
									}
									goto l211
								l210:
									position, tokenIndex, depth = position210, tokenIndex210, depth210
								}
							l211:
								goto l206
							l207:
								position, tokenIndex, depth = position207, tokenIndex207, depth207
							}
							goto l204
						l205:
							position, tokenIndex, depth = position204, tokenIndex204, depth204
							if !rules[rulefilterOrBind]() {
								goto l201
							}
							{
								position214, tokenIndex214, depth214 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l214
								}
								goto l215
							l214:
								position, tokenIndex, depth = position214, tokenIndex214, depth214
							}
						l215:
							{
								position216, tokenIndex216, depth216 := position, tokenIndex, depth
								if !rules[ruletriplesBlock]() {
									goto l216
								}
								goto l217
							l216:
								position, tokenIndex, depth = position216, tokenIndex216, depth216
							}
						l217:
						l212:
							{
								position213, tokenIndex213, depth213 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l213
								}
								{
									position218, tokenIndex218, depth218 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l218
									}
									goto l219
								l218:
									position, tokenIndex, depth = position218, tokenIndex218, depth218
								}
							l219:
								{
									position220, tokenIndex220, depth220 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l220
									}
									goto l221
								l220:
									position, tokenIndex, depth = position220, tokenIndex220, depth220
								}
							l221:
								goto l212
							l213:
								position, tokenIndex, depth = position213, tokenIndex213, depth213
							}
						}
					l204:
						depth--
						add(rulebasicGraphPattern, position203)
					}
					goto l202
				l201:
					position, tokenIndex, depth = position201, tokenIndex201, depth201
				}
			l202:
				{
					position222, tokenIndex222, depth222 := position, tokenIndex, depth
					{
						position224 := position
						depth++
						{
							position225, tokenIndex225, depth225 := position, tokenIndex, depth
							{
								position227 := position
								depth++
								{
									position228 := position
									depth++
									{
										position229, tokenIndex229, depth229 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l230
										}
										position++
										goto l229
									l230:
										position, tokenIndex, depth = position229, tokenIndex229, depth229
										if buffer[position] != rune('O') {
											goto l226
										}
										position++
									}
								l229:
									{
										position231, tokenIndex231, depth231 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l232
										}
										position++
										goto l231
									l232:
										position, tokenIndex, depth = position231, tokenIndex231, depth231
										if buffer[position] != rune('P') {
											goto l226
										}
										position++
									}
								l231:
									{
										position233, tokenIndex233, depth233 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l234
										}
										position++
										goto l233
									l234:
										position, tokenIndex, depth = position233, tokenIndex233, depth233
										if buffer[position] != rune('T') {
											goto l226
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
											goto l226
										}
										position++
									}
								l235:
									{
										position237, tokenIndex237, depth237 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l238
										}
										position++
										goto l237
									l238:
										position, tokenIndex, depth = position237, tokenIndex237, depth237
										if buffer[position] != rune('O') {
											goto l226
										}
										position++
									}
								l237:
									{
										position239, tokenIndex239, depth239 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l240
										}
										position++
										goto l239
									l240:
										position, tokenIndex, depth = position239, tokenIndex239, depth239
										if buffer[position] != rune('N') {
											goto l226
										}
										position++
									}
								l239:
									{
										position241, tokenIndex241, depth241 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l242
										}
										position++
										goto l241
									l242:
										position, tokenIndex, depth = position241, tokenIndex241, depth241
										if buffer[position] != rune('A') {
											goto l226
										}
										position++
									}
								l241:
									{
										position243, tokenIndex243, depth243 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l244
										}
										position++
										goto l243
									l244:
										position, tokenIndex, depth = position243, tokenIndex243, depth243
										if buffer[position] != rune('L') {
											goto l226
										}
										position++
									}
								l243:
									if !rules[ruleskip]() {
										goto l226
									}
									depth--
									add(ruleOPTIONAL, position228)
								}
								if !rules[ruleLBRACE]() {
									goto l226
								}
								{
									position245, tokenIndex245, depth245 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l246
									}
									goto l245
								l246:
									position, tokenIndex, depth = position245, tokenIndex245, depth245
									if !rules[rulegraphPattern]() {
										goto l226
									}
								}
							l245:
								if !rules[ruleRBRACE]() {
									goto l226
								}
								depth--
								add(ruleoptionalGraphPattern, position227)
							}
							goto l225
						l226:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l247
							}
							goto l225
						l247:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							{
								position249 := position
								depth++
								{
									position250 := position
									depth++
									{
										position251, tokenIndex251, depth251 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l252
										}
										position++
										goto l251
									l252:
										position, tokenIndex, depth = position251, tokenIndex251, depth251
										if buffer[position] != rune('G') {
											goto l248
										}
										position++
									}
								l251:
									{
										position253, tokenIndex253, depth253 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l254
										}
										position++
										goto l253
									l254:
										position, tokenIndex, depth = position253, tokenIndex253, depth253
										if buffer[position] != rune('R') {
											goto l248
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
											goto l248
										}
										position++
									}
								l255:
									{
										position257, tokenIndex257, depth257 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l258
										}
										position++
										goto l257
									l258:
										position, tokenIndex, depth = position257, tokenIndex257, depth257
										if buffer[position] != rune('P') {
											goto l248
										}
										position++
									}
								l257:
									{
										position259, tokenIndex259, depth259 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l260
										}
										position++
										goto l259
									l260:
										position, tokenIndex, depth = position259, tokenIndex259, depth259
										if buffer[position] != rune('H') {
											goto l248
										}
										position++
									}
								l259:
									if !rules[ruleskip]() {
										goto l248
									}
									depth--
									add(ruleGRAPH, position250)
								}
								{
									position261, tokenIndex261, depth261 := position, tokenIndex, depth
									if !rules[rulevar]() {
										goto l262
									}
									goto l261
								l262:
									position, tokenIndex, depth = position261, tokenIndex261, depth261
									if !rules[ruleiriref]() {
										goto l248
									}
								}
							l261:
								if !rules[rulegroupGraphPattern]() {
									goto l248
								}
								depth--
								add(rulegraphGraphPattern, position249)
							}
							goto l225
						l248:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							{
								position264 := position
								depth++
								{
									position265 := position
									depth++
									{
										position266, tokenIndex266, depth266 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l267
										}
										position++
										goto l266
									l267:
										position, tokenIndex, depth = position266, tokenIndex266, depth266
										if buffer[position] != rune('M') {
											goto l263
										}
										position++
									}
								l266:
									{
										position268, tokenIndex268, depth268 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l269
										}
										position++
										goto l268
									l269:
										position, tokenIndex, depth = position268, tokenIndex268, depth268
										if buffer[position] != rune('I') {
											goto l263
										}
										position++
									}
								l268:
									{
										position270, tokenIndex270, depth270 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l271
										}
										position++
										goto l270
									l271:
										position, tokenIndex, depth = position270, tokenIndex270, depth270
										if buffer[position] != rune('N') {
											goto l263
										}
										position++
									}
								l270:
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
											goto l263
										}
										position++
									}
								l272:
									{
										position274, tokenIndex274, depth274 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l275
										}
										position++
										goto l274
									l275:
										position, tokenIndex, depth = position274, tokenIndex274, depth274
										if buffer[position] != rune('S') {
											goto l263
										}
										position++
									}
								l274:
									if !rules[ruleskip]() {
										goto l263
									}
									depth--
									add(ruleMINUSSETOPER, position265)
								}
								if !rules[rulegroupGraphPattern]() {
									goto l263
								}
								depth--
								add(ruleminusGraphPattern, position264)
							}
							goto l225
						l263:
							position, tokenIndex, depth = position225, tokenIndex225, depth225
							{
								position276 := position
								depth++
								{
									position277 := position
									depth++
									depth--
									add(ruleSERVICE, position277)
								}
								{
									position278, tokenIndex278, depth278 := position, tokenIndex, depth
									{
										position280 := position
										depth++
										depth--
										add(ruleSILENT, position280)
									}
									goto l279

									position, tokenIndex, depth = position278, tokenIndex278, depth278
								}
							l279:
								{
									position281, tokenIndex281, depth281 := position, tokenIndex, depth
									if !rules[rulevar]() {
										goto l282
									}
									goto l281
								l282:
									position, tokenIndex, depth = position281, tokenIndex281, depth281
									if !rules[ruleiriref]() {
										goto l222
									}
								}
							l281:
								if !rules[rulegroupGraphPattern]() {
									goto l222
								}
								depth--
								add(ruleserviceGraphPattern, position276)
							}
						}
					l225:
						depth--
						add(rulegraphPatternNotTriples, position224)
					}
					{
						position283, tokenIndex283, depth283 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l283
						}
						goto l284
					l283:
						position, tokenIndex, depth = position283, tokenIndex283, depth283
					}
				l284:
					if !rules[rulegraphPattern]() {
						goto l222
					}
					goto l223
				l222:
					position, tokenIndex, depth = position222, tokenIndex222, depth222
				}
			l223:
				depth--
				add(rulegraphPattern, position200)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern / graphGraphPattern / minusGraphPattern / serviceGraphPattern)> */
		nil,
		/* 19 serviceGraphPattern <- <(SERVICE SILENT? (var / iriref) groupGraphPattern)> */
		nil,
		/* 20 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 21 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position288, tokenIndex288, depth288 := position, tokenIndex, depth
			{
				position289 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l288
				}
				{
					position290, tokenIndex290, depth290 := position, tokenIndex, depth
					{
						position292 := position
						depth++
						{
							position293, tokenIndex293, depth293 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l294
							}
							position++
							goto l293
						l294:
							position, tokenIndex, depth = position293, tokenIndex293, depth293
							if buffer[position] != rune('U') {
								goto l290
							}
							position++
						}
					l293:
						{
							position295, tokenIndex295, depth295 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l296
							}
							position++
							goto l295
						l296:
							position, tokenIndex, depth = position295, tokenIndex295, depth295
							if buffer[position] != rune('N') {
								goto l290
							}
							position++
						}
					l295:
						{
							position297, tokenIndex297, depth297 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l298
							}
							position++
							goto l297
						l298:
							position, tokenIndex, depth = position297, tokenIndex297, depth297
							if buffer[position] != rune('I') {
								goto l290
							}
							position++
						}
					l297:
						{
							position299, tokenIndex299, depth299 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l300
							}
							position++
							goto l299
						l300:
							position, tokenIndex, depth = position299, tokenIndex299, depth299
							if buffer[position] != rune('O') {
								goto l290
							}
							position++
						}
					l299:
						{
							position301, tokenIndex301, depth301 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l302
							}
							position++
							goto l301
						l302:
							position, tokenIndex, depth = position301, tokenIndex301, depth301
							if buffer[position] != rune('N') {
								goto l290
							}
							position++
						}
					l301:
						if !rules[ruleskip]() {
							goto l290
						}
						depth--
						add(ruleUNION, position292)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l290
					}
					goto l291
				l290:
					position, tokenIndex, depth = position290, tokenIndex290, depth290
				}
			l291:
				depth--
				add(rulegroupOrUnionGraphPattern, position289)
			}
			return true
		l288:
			position, tokenIndex, depth = position288, tokenIndex288, depth288
			return false
		},
		/* 22 graphGraphPattern <- <(GRAPH (var / iriref) groupGraphPattern)> */
		nil,
		/* 23 minusGraphPattern <- <(MINUSSETOPER groupGraphPattern)> */
		nil,
		/* 24 basicGraphPattern <- <((triplesBlock (filterOrBind DOT? triplesBlock?)*) / (filterOrBind DOT? triplesBlock?)+)> */
		nil,
		/* 25 filterOrBind <- <((FILTER constraint) / (BIND LPAREN expression AS var RPAREN))> */
		func() bool {
			position306, tokenIndex306, depth306 := position, tokenIndex, depth
			{
				position307 := position
				depth++
				{
					position308, tokenIndex308, depth308 := position, tokenIndex, depth
					{
						position310 := position
						depth++
						{
							position311, tokenIndex311, depth311 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l312
							}
							position++
							goto l311
						l312:
							position, tokenIndex, depth = position311, tokenIndex311, depth311
							if buffer[position] != rune('F') {
								goto l309
							}
							position++
						}
					l311:
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l314
							}
							position++
							goto l313
						l314:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
							if buffer[position] != rune('I') {
								goto l309
							}
							position++
						}
					l313:
						{
							position315, tokenIndex315, depth315 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l316
							}
							position++
							goto l315
						l316:
							position, tokenIndex, depth = position315, tokenIndex315, depth315
							if buffer[position] != rune('L') {
								goto l309
							}
							position++
						}
					l315:
						{
							position317, tokenIndex317, depth317 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l318
							}
							position++
							goto l317
						l318:
							position, tokenIndex, depth = position317, tokenIndex317, depth317
							if buffer[position] != rune('T') {
								goto l309
							}
							position++
						}
					l317:
						{
							position319, tokenIndex319, depth319 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l320
							}
							position++
							goto l319
						l320:
							position, tokenIndex, depth = position319, tokenIndex319, depth319
							if buffer[position] != rune('E') {
								goto l309
							}
							position++
						}
					l319:
						{
							position321, tokenIndex321, depth321 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l322
							}
							position++
							goto l321
						l322:
							position, tokenIndex, depth = position321, tokenIndex321, depth321
							if buffer[position] != rune('R') {
								goto l309
							}
							position++
						}
					l321:
						if !rules[ruleskip]() {
							goto l309
						}
						depth--
						add(ruleFILTER, position310)
					}
					if !rules[ruleconstraint]() {
						goto l309
					}
					goto l308
				l309:
					position, tokenIndex, depth = position308, tokenIndex308, depth308
					{
						position323 := position
						depth++
						{
							position324, tokenIndex324, depth324 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l325
							}
							position++
							goto l324
						l325:
							position, tokenIndex, depth = position324, tokenIndex324, depth324
							if buffer[position] != rune('B') {
								goto l306
							}
							position++
						}
					l324:
						{
							position326, tokenIndex326, depth326 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l327
							}
							position++
							goto l326
						l327:
							position, tokenIndex, depth = position326, tokenIndex326, depth326
							if buffer[position] != rune('I') {
								goto l306
							}
							position++
						}
					l326:
						{
							position328, tokenIndex328, depth328 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l329
							}
							position++
							goto l328
						l329:
							position, tokenIndex, depth = position328, tokenIndex328, depth328
							if buffer[position] != rune('N') {
								goto l306
							}
							position++
						}
					l328:
						{
							position330, tokenIndex330, depth330 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l331
							}
							position++
							goto l330
						l331:
							position, tokenIndex, depth = position330, tokenIndex330, depth330
							if buffer[position] != rune('D') {
								goto l306
							}
							position++
						}
					l330:
						if !rules[ruleskip]() {
							goto l306
						}
						depth--
						add(ruleBIND, position323)
					}
					if !rules[ruleLPAREN]() {
						goto l306
					}
					if !rules[ruleexpression]() {
						goto l306
					}
					if !rules[ruleAS]() {
						goto l306
					}
					if !rules[rulevar]() {
						goto l306
					}
					if !rules[ruleRPAREN]() {
						goto l306
					}
				}
			l308:
				depth--
				add(rulefilterOrBind, position307)
			}
			return true
		l306:
			position, tokenIndex, depth = position306, tokenIndex306, depth306
			return false
		},
		/* 26 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		func() bool {
			position332, tokenIndex332, depth332 := position, tokenIndex, depth
			{
				position333 := position
				depth++
				{
					position334, tokenIndex334, depth334 := position, tokenIndex, depth
					if !rules[rulebrackettedExpression]() {
						goto l335
					}
					goto l334
				l335:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if !rules[rulebuiltinCall]() {
						goto l336
					}
					goto l334
				l336:
					position, tokenIndex, depth = position334, tokenIndex334, depth334
					if !rules[rulefunctionCall]() {
						goto l332
					}
				}
			l334:
				depth--
				add(ruleconstraint, position333)
			}
			return true
		l332:
			position, tokenIndex, depth = position332, tokenIndex332, depth332
			return false
		},
		/* 27 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l337
				}
			l339:
				{
					position340, tokenIndex340, depth340 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l340
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l340
					}
					goto l339
				l340:
					position, tokenIndex, depth = position340, tokenIndex340, depth340
				}
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l341
					}
					goto l342
				l341:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
				}
			l342:
				depth--
				add(ruletriplesBlock, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 28 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath?))> */
		func() bool {
			position343, tokenIndex343, depth343 := position, tokenIndex, depth
			{
				position344 := position
				depth++
				{
					position345, tokenIndex345, depth345 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l346
					}
					if !rules[rulepropertyListPath]() {
						goto l346
					}
					goto l345
				l346:
					position, tokenIndex, depth = position345, tokenIndex345, depth345
					if !rules[ruletriplesNodePath]() {
						goto l343
					}
					{
						position347, tokenIndex347, depth347 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l347
						}
						goto l348
					l347:
						position, tokenIndex, depth = position347, tokenIndex347, depth347
					}
				l348:
				}
			l345:
				depth--
				add(ruletriplesSameSubjectPath, position344)
			}
			return true
		l343:
			position, tokenIndex, depth = position343, tokenIndex343, depth343
			return false
		},
		/* 29 varOrTerm <- <(var / graphTerm)> */
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
													position361, tokenIndex361, depth361 := position, tokenIndex, depth
													if !rules[rulepnCharsU]() {
														goto l362
													}
													goto l361
												l362:
													position, tokenIndex, depth = position361, tokenIndex361, depth361
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l359
													}
													position++
												}
											l361:
												{
													position363, tokenIndex363, depth363 := position, tokenIndex, depth
													{
														position365, tokenIndex365, depth365 := position, tokenIndex, depth
													l367:
														{
															position368, tokenIndex368, depth368 := position, tokenIndex, depth
															{
																position369, tokenIndex369, depth369 := position, tokenIndex, depth
																if !rules[rulepnCharsU]() {
																	goto l370
																}
																goto l369
															l370:
																position, tokenIndex, depth = position369, tokenIndex369, depth369
																{
																	switch buffer[position] {
																	case '.':
																		if buffer[position] != rune('.') {
																			goto l368
																		}
																		position++
																		break
																	case '-':
																		if buffer[position] != rune('-') {
																			goto l368
																		}
																		position++
																		break
																	default:
																		if c := buffer[position]; c < rune('0') || c > rune('9') {
																			goto l368
																		}
																		position++
																		break
																	}
																}

															}
														l369:
															goto l367
														l368:
															position, tokenIndex, depth = position368, tokenIndex368, depth368
														}
														if !rules[rulepnCharsU]() {
															goto l366
														}
														goto l365
													l366:
														position, tokenIndex, depth = position365, tokenIndex365, depth365
														{
															position372, tokenIndex372, depth372 := position, tokenIndex, depth
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l373
															}
															position++
															goto l372
														l373:
															position, tokenIndex, depth = position372, tokenIndex372, depth372
															if buffer[position] != rune('-') {
																goto l363
															}
															position++
														}
													l372:
													}
												l365:
													goto l364
												l363:
													position, tokenIndex, depth = position363, tokenIndex363, depth363
												}
											l364:
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
												position374 := position
												depth++
												if buffer[position] != rune('[') {
													goto l349
												}
												position++
											l375:
												{
													position376, tokenIndex376, depth376 := position, tokenIndex, depth
													if !rules[rulews]() {
														goto l376
													}
													goto l375
												l376:
													position, tokenIndex, depth = position376, tokenIndex376, depth376
												}
												if buffer[position] != rune(']') {
													goto l349
												}
												position++
												if !rules[ruleskip]() {
													goto l349
												}
												depth--
												add(ruleanon, position374)
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
								case '"', '\'':
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
		/* 30 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 31 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					{
						position382 := position
						depth++
						if !rules[ruleLPAREN]() {
							goto l381
						}
						if !rules[rulegraphNodePath]() {
							goto l381
						}
					l383:
						{
							position384, tokenIndex384, depth384 := position, tokenIndex, depth
							if !rules[rulegraphNodePath]() {
								goto l384
							}
							goto l383
						l384:
							position, tokenIndex, depth = position384, tokenIndex384, depth384
						}
						if !rules[ruleRPAREN]() {
							goto l381
						}
						depth--
						add(rulecollectionPath, position382)
					}
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					{
						position385 := position
						depth++
						{
							position386 := position
							depth++
							if buffer[position] != rune('[') {
								goto l378
							}
							position++
							if !rules[ruleskip]() {
								goto l378
							}
							depth--
							add(ruleLBRACK, position386)
						}
						if !rules[rulepropertyListPath]() {
							goto l378
						}
						{
							position387 := position
							depth++
							if buffer[position] != rune(']') {
								goto l378
							}
							position++
							if !rules[ruleskip]() {
								goto l378
							}
							depth--
							add(ruleRBRACK, position387)
						}
						depth--
						add(ruleblankNodePropertyListPath, position385)
					}
				}
			l380:
				depth--
				add(ruletriplesNodePath, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 32 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 33 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 34 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l393
					}
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					{
						position394 := position
						depth++
						if !rules[rulepath]() {
							goto l390
						}
						depth--
						add(ruleverbPath, position394)
					}
				}
			l392:
				{
					position395 := position
					depth++
					if !rules[ruleobjectPath]() {
						goto l390
					}
				l396:
					{
						position397, tokenIndex397, depth397 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l397
						}
						if !rules[ruleobjectPath]() {
							goto l397
						}
						goto l396
					l397:
						position, tokenIndex, depth = position397, tokenIndex397, depth397
					}
					depth--
					add(ruleobjectListPath, position395)
				}
				{
					position398, tokenIndex398, depth398 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l398
					}
					{
						position400, tokenIndex400, depth400 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l400
						}
						goto l401
					l400:
						position, tokenIndex, depth = position400, tokenIndex400, depth400
					}
				l401:
					goto l399
				l398:
					position, tokenIndex, depth = position398, tokenIndex398, depth398
				}
			l399:
				depth--
				add(rulepropertyListPath, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 35 verbPath <- <path> */
		nil,
		/* 36 path <- <pathAlternative> */
		func() bool {
			position403, tokenIndex403, depth403 := position, tokenIndex, depth
			{
				position404 := position
				depth++
				{
					position405 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l403
					}
				l406:
					{
						position407, tokenIndex407, depth407 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l407
						}
						if !rules[rulepathSequence]() {
							goto l407
						}
						goto l406
					l407:
						position, tokenIndex, depth = position407, tokenIndex407, depth407
					}
					depth--
					add(rulepathAlternative, position405)
				}
				depth--
				add(rulepath, position404)
			}
			return true
		l403:
			position, tokenIndex, depth = position403, tokenIndex403, depth403
			return false
		},
		/* 37 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 38 pathSequence <- <(pathElt (SLASH pathElt)*)> */
		func() bool {
			position409, tokenIndex409, depth409 := position, tokenIndex, depth
			{
				position410 := position
				depth++
				if !rules[rulepathElt]() {
					goto l409
				}
			l411:
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l412
					}
					if !rules[rulepathElt]() {
						goto l412
					}
					goto l411
				l412:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
				}
				depth--
				add(rulepathSequence, position410)
			}
			return true
		l409:
			position, tokenIndex, depth = position409, tokenIndex409, depth409
			return false
		},
		/* 39 pathElt <- <(INVERSE? pathPrimary pathMod?)> */
		func() bool {
			position413, tokenIndex413, depth413 := position, tokenIndex, depth
			{
				position414 := position
				depth++
				{
					position415, tokenIndex415, depth415 := position, tokenIndex, depth
					if !rules[ruleINVERSE]() {
						goto l415
					}
					goto l416
				l415:
					position, tokenIndex, depth = position415, tokenIndex415, depth415
				}
			l416:
				{
					position417 := position
					depth++
					{
						position418, tokenIndex418, depth418 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l419
						}
						goto l418
					l419:
						position, tokenIndex, depth = position418, tokenIndex418, depth418
						{
							switch buffer[position] {
							case '(':
								if !rules[ruleLPAREN]() {
									goto l413
								}
								if !rules[rulepath]() {
									goto l413
								}
								if !rules[ruleRPAREN]() {
									goto l413
								}
								break
							case '!':
								if !rules[ruleNOT]() {
									goto l413
								}
								{
									position421 := position
									depth++
									{
										position422, tokenIndex422, depth422 := position, tokenIndex, depth
										if !rules[rulepathOneInPropertySet]() {
											goto l423
										}
										goto l422
									l423:
										position, tokenIndex, depth = position422, tokenIndex422, depth422
										if !rules[ruleLPAREN]() {
											goto l413
										}
										{
											position424, tokenIndex424, depth424 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l424
											}
										l426:
											{
												position427, tokenIndex427, depth427 := position, tokenIndex, depth
												if !rules[rulePIPE]() {
													goto l427
												}
												if !rules[rulepathOneInPropertySet]() {
													goto l427
												}
												goto l426
											l427:
												position, tokenIndex, depth = position427, tokenIndex427, depth427
											}
											goto l425
										l424:
											position, tokenIndex, depth = position424, tokenIndex424, depth424
										}
									l425:
										if !rules[ruleRPAREN]() {
											goto l413
										}
									}
								l422:
									depth--
									add(rulepathNegatedPropertySet, position421)
								}
								break
							default:
								if !rules[ruleISA]() {
									goto l413
								}
								break
							}
						}

					}
				l418:
					depth--
					add(rulepathPrimary, position417)
				}
				{
					position428, tokenIndex428, depth428 := position, tokenIndex, depth
					{
						position430 := position
						depth++
						{
							switch buffer[position] {
							case '+':
								if !rules[rulePLUS]() {
									goto l428
								}
								break
							case '?':
								{
									position432 := position
									depth++
									if buffer[position] != rune('?') {
										goto l428
									}
									position++
									if !rules[ruleskip]() {
										goto l428
									}
									depth--
									add(ruleQUESTION, position432)
								}
								break
							default:
								if !rules[ruleSTAR]() {
									goto l428
								}
								break
							}
						}

						{
							position433, tokenIndex433, depth433 := position, tokenIndex, depth
							if !matchDot() {
								goto l433
							}
							goto l428
						l433:
							position, tokenIndex, depth = position433, tokenIndex433, depth433
						}
						depth--
						add(rulepathMod, position430)
					}
					goto l429
				l428:
					position, tokenIndex, depth = position428, tokenIndex428, depth428
				}
			l429:
				depth--
				add(rulepathElt, position414)
			}
			return true
		l413:
			position, tokenIndex, depth = position413, tokenIndex413, depth413
			return false
		},
		/* 40 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 41 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 42 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position436, tokenIndex436, depth436 := position, tokenIndex, depth
			{
				position437 := position
				depth++
				{
					position438, tokenIndex438, depth438 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l439
					}
					goto l438
				l439:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if !rules[ruleISA]() {
						goto l440
					}
					goto l438
				l440:
					position, tokenIndex, depth = position438, tokenIndex438, depth438
					if !rules[ruleINVERSE]() {
						goto l436
					}
					{
						position441, tokenIndex441, depth441 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l442
						}
						goto l441
					l442:
						position, tokenIndex, depth = position441, tokenIndex441, depth441
						if !rules[ruleISA]() {
							goto l436
						}
					}
				l441:
				}
			l438:
				depth--
				add(rulepathOneInPropertySet, position437)
			}
			return true
		l436:
			position, tokenIndex, depth = position436, tokenIndex436, depth436
			return false
		},
		/* 43 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !.)> */
		nil,
		/* 44 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 45 objectPath <- <graphNodePath> */
		func() bool {
			position445, tokenIndex445, depth445 := position, tokenIndex, depth
			{
				position446 := position
				depth++
				if !rules[rulegraphNodePath]() {
					goto l445
				}
				depth--
				add(ruleobjectPath, position446)
			}
			return true
		l445:
			position, tokenIndex, depth = position445, tokenIndex445, depth445
			return false
		},
		/* 46 graphNodePath <- <(varOrTerm / triplesNodePath)> */
		func() bool {
			position447, tokenIndex447, depth447 := position, tokenIndex, depth
			{
				position448 := position
				depth++
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l450
					}
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					if !rules[ruletriplesNodePath]() {
						goto l447
					}
				}
			l449:
				depth--
				add(rulegraphNodePath, position448)
			}
			return true
		l447:
			position, tokenIndex, depth = position447, tokenIndex447, depth447
			return false
		},
		/* 47 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position452 := position
				depth++
				{
					position453, tokenIndex453, depth453 := position, tokenIndex, depth
					{
						position455, tokenIndex455, depth455 := position, tokenIndex, depth
						{
							position457 := position
							depth++
							{
								position458, tokenIndex458, depth458 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l459
								}
								position++
								goto l458
							l459:
								position, tokenIndex, depth = position458, tokenIndex458, depth458
								if buffer[position] != rune('O') {
									goto l456
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
									goto l456
								}
								position++
							}
						l460:
							{
								position462, tokenIndex462, depth462 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l463
								}
								position++
								goto l462
							l463:
								position, tokenIndex, depth = position462, tokenIndex462, depth462
								if buffer[position] != rune('D') {
									goto l456
								}
								position++
							}
						l462:
							{
								position464, tokenIndex464, depth464 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l465
								}
								position++
								goto l464
							l465:
								position, tokenIndex, depth = position464, tokenIndex464, depth464
								if buffer[position] != rune('E') {
									goto l456
								}
								position++
							}
						l464:
							{
								position466, tokenIndex466, depth466 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l467
								}
								position++
								goto l466
							l467:
								position, tokenIndex, depth = position466, tokenIndex466, depth466
								if buffer[position] != rune('R') {
									goto l456
								}
								position++
							}
						l466:
							if !rules[ruleskip]() {
								goto l456
							}
							depth--
							add(ruleORDER, position457)
						}
						if !rules[ruleBY]() {
							goto l456
						}
						{
							position470 := position
							depth++
							{
								position471, tokenIndex471, depth471 := position, tokenIndex, depth
								{
									position473, tokenIndex473, depth473 := position, tokenIndex, depth
									{
										position475, tokenIndex475, depth475 := position, tokenIndex, depth
										{
											position477 := position
											depth++
											{
												position478, tokenIndex478, depth478 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l479
												}
												position++
												goto l478
											l479:
												position, tokenIndex, depth = position478, tokenIndex478, depth478
												if buffer[position] != rune('A') {
													goto l476
												}
												position++
											}
										l478:
											{
												position480, tokenIndex480, depth480 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l481
												}
												position++
												goto l480
											l481:
												position, tokenIndex, depth = position480, tokenIndex480, depth480
												if buffer[position] != rune('S') {
													goto l476
												}
												position++
											}
										l480:
											{
												position482, tokenIndex482, depth482 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l483
												}
												position++
												goto l482
											l483:
												position, tokenIndex, depth = position482, tokenIndex482, depth482
												if buffer[position] != rune('C') {
													goto l476
												}
												position++
											}
										l482:
											if !rules[ruleskip]() {
												goto l476
											}
											depth--
											add(ruleASC, position477)
										}
										goto l475
									l476:
										position, tokenIndex, depth = position475, tokenIndex475, depth475
										{
											position484 := position
											depth++
											{
												position485, tokenIndex485, depth485 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l486
												}
												position++
												goto l485
											l486:
												position, tokenIndex, depth = position485, tokenIndex485, depth485
												if buffer[position] != rune('D') {
													goto l473
												}
												position++
											}
										l485:
											{
												position487, tokenIndex487, depth487 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l488
												}
												position++
												goto l487
											l488:
												position, tokenIndex, depth = position487, tokenIndex487, depth487
												if buffer[position] != rune('E') {
													goto l473
												}
												position++
											}
										l487:
											{
												position489, tokenIndex489, depth489 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l490
												}
												position++
												goto l489
											l490:
												position, tokenIndex, depth = position489, tokenIndex489, depth489
												if buffer[position] != rune('S') {
													goto l473
												}
												position++
											}
										l489:
											{
												position491, tokenIndex491, depth491 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l492
												}
												position++
												goto l491
											l492:
												position, tokenIndex, depth = position491, tokenIndex491, depth491
												if buffer[position] != rune('C') {
													goto l473
												}
												position++
											}
										l491:
											if !rules[ruleskip]() {
												goto l473
											}
											depth--
											add(ruleDESC, position484)
										}
									}
								l475:
									goto l474
								l473:
									position, tokenIndex, depth = position473, tokenIndex473, depth473
								}
							l474:
								if !rules[rulebrackettedExpression]() {
									goto l472
								}
								goto l471
							l472:
								position, tokenIndex, depth = position471, tokenIndex471, depth471
								if !rules[rulefunctionCall]() {
									goto l493
								}
								goto l471
							l493:
								position, tokenIndex, depth = position471, tokenIndex471, depth471
								if !rules[rulebuiltinCall]() {
									goto l494
								}
								goto l471
							l494:
								position, tokenIndex, depth = position471, tokenIndex471, depth471
								if !rules[rulevar]() {
									goto l456
								}
							}
						l471:
							depth--
							add(ruleorderCondition, position470)
						}
					l468:
						{
							position469, tokenIndex469, depth469 := position, tokenIndex, depth
							{
								position495 := position
								depth++
								{
									position496, tokenIndex496, depth496 := position, tokenIndex, depth
									{
										position498, tokenIndex498, depth498 := position, tokenIndex, depth
										{
											position500, tokenIndex500, depth500 := position, tokenIndex, depth
											{
												position502 := position
												depth++
												{
													position503, tokenIndex503, depth503 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l504
													}
													position++
													goto l503
												l504:
													position, tokenIndex, depth = position503, tokenIndex503, depth503
													if buffer[position] != rune('A') {
														goto l501
													}
													position++
												}
											l503:
												{
													position505, tokenIndex505, depth505 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l506
													}
													position++
													goto l505
												l506:
													position, tokenIndex, depth = position505, tokenIndex505, depth505
													if buffer[position] != rune('S') {
														goto l501
													}
													position++
												}
											l505:
												{
													position507, tokenIndex507, depth507 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l508
													}
													position++
													goto l507
												l508:
													position, tokenIndex, depth = position507, tokenIndex507, depth507
													if buffer[position] != rune('C') {
														goto l501
													}
													position++
												}
											l507:
												if !rules[ruleskip]() {
													goto l501
												}
												depth--
												add(ruleASC, position502)
											}
											goto l500
										l501:
											position, tokenIndex, depth = position500, tokenIndex500, depth500
											{
												position509 := position
												depth++
												{
													position510, tokenIndex510, depth510 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l511
													}
													position++
													goto l510
												l511:
													position, tokenIndex, depth = position510, tokenIndex510, depth510
													if buffer[position] != rune('D') {
														goto l498
													}
													position++
												}
											l510:
												{
													position512, tokenIndex512, depth512 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l513
													}
													position++
													goto l512
												l513:
													position, tokenIndex, depth = position512, tokenIndex512, depth512
													if buffer[position] != rune('E') {
														goto l498
													}
													position++
												}
											l512:
												{
													position514, tokenIndex514, depth514 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l515
													}
													position++
													goto l514
												l515:
													position, tokenIndex, depth = position514, tokenIndex514, depth514
													if buffer[position] != rune('S') {
														goto l498
													}
													position++
												}
											l514:
												{
													position516, tokenIndex516, depth516 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l517
													}
													position++
													goto l516
												l517:
													position, tokenIndex, depth = position516, tokenIndex516, depth516
													if buffer[position] != rune('C') {
														goto l498
													}
													position++
												}
											l516:
												if !rules[ruleskip]() {
													goto l498
												}
												depth--
												add(ruleDESC, position509)
											}
										}
									l500:
										goto l499
									l498:
										position, tokenIndex, depth = position498, tokenIndex498, depth498
									}
								l499:
									if !rules[rulebrackettedExpression]() {
										goto l497
									}
									goto l496
								l497:
									position, tokenIndex, depth = position496, tokenIndex496, depth496
									if !rules[rulefunctionCall]() {
										goto l518
									}
									goto l496
								l518:
									position, tokenIndex, depth = position496, tokenIndex496, depth496
									if !rules[rulebuiltinCall]() {
										goto l519
									}
									goto l496
								l519:
									position, tokenIndex, depth = position496, tokenIndex496, depth496
									if !rules[rulevar]() {
										goto l469
									}
								}
							l496:
								depth--
								add(ruleorderCondition, position495)
							}
							goto l468
						l469:
							position, tokenIndex, depth = position469, tokenIndex469, depth469
						}
						goto l455
					l456:
						position, tokenIndex, depth = position455, tokenIndex455, depth455
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position521 := position
									depth++
									{
										position522, tokenIndex522, depth522 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l523
										}
										position++
										goto l522
									l523:
										position, tokenIndex, depth = position522, tokenIndex522, depth522
										if buffer[position] != rune('H') {
											goto l453
										}
										position++
									}
								l522:
									{
										position524, tokenIndex524, depth524 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l525
										}
										position++
										goto l524
									l525:
										position, tokenIndex, depth = position524, tokenIndex524, depth524
										if buffer[position] != rune('A') {
											goto l453
										}
										position++
									}
								l524:
									{
										position526, tokenIndex526, depth526 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l527
										}
										position++
										goto l526
									l527:
										position, tokenIndex, depth = position526, tokenIndex526, depth526
										if buffer[position] != rune('V') {
											goto l453
										}
										position++
									}
								l526:
									{
										position528, tokenIndex528, depth528 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l529
										}
										position++
										goto l528
									l529:
										position, tokenIndex, depth = position528, tokenIndex528, depth528
										if buffer[position] != rune('I') {
											goto l453
										}
										position++
									}
								l528:
									{
										position530, tokenIndex530, depth530 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l531
										}
										position++
										goto l530
									l531:
										position, tokenIndex, depth = position530, tokenIndex530, depth530
										if buffer[position] != rune('N') {
											goto l453
										}
										position++
									}
								l530:
									{
										position532, tokenIndex532, depth532 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l533
										}
										position++
										goto l532
									l533:
										position, tokenIndex, depth = position532, tokenIndex532, depth532
										if buffer[position] != rune('G') {
											goto l453
										}
										position++
									}
								l532:
									if !rules[ruleskip]() {
										goto l453
									}
									depth--
									add(ruleHAVING, position521)
								}
								if !rules[ruleconstraint]() {
									goto l453
								}
								break
							case 'G', 'g':
								{
									position534 := position
									depth++
									{
										position535, tokenIndex535, depth535 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l536
										}
										position++
										goto l535
									l536:
										position, tokenIndex, depth = position535, tokenIndex535, depth535
										if buffer[position] != rune('G') {
											goto l453
										}
										position++
									}
								l535:
									{
										position537, tokenIndex537, depth537 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l538
										}
										position++
										goto l537
									l538:
										position, tokenIndex, depth = position537, tokenIndex537, depth537
										if buffer[position] != rune('R') {
											goto l453
										}
										position++
									}
								l537:
									{
										position539, tokenIndex539, depth539 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l540
										}
										position++
										goto l539
									l540:
										position, tokenIndex, depth = position539, tokenIndex539, depth539
										if buffer[position] != rune('O') {
											goto l453
										}
										position++
									}
								l539:
									{
										position541, tokenIndex541, depth541 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l542
										}
										position++
										goto l541
									l542:
										position, tokenIndex, depth = position541, tokenIndex541, depth541
										if buffer[position] != rune('U') {
											goto l453
										}
										position++
									}
								l541:
									{
										position543, tokenIndex543, depth543 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l544
										}
										position++
										goto l543
									l544:
										position, tokenIndex, depth = position543, tokenIndex543, depth543
										if buffer[position] != rune('P') {
											goto l453
										}
										position++
									}
								l543:
									if !rules[ruleskip]() {
										goto l453
									}
									depth--
									add(ruleGROUP, position534)
								}
								if !rules[ruleBY]() {
									goto l453
								}
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
													goto l453
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l453
												}
												if !rules[ruleexpression]() {
													goto l453
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
													goto l453
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l453
												}
												break
											}
										}

									}
								l548:
									depth--
									add(rulegroupCondition, position547)
								}
							l545:
								{
									position546, tokenIndex546, depth546 := position, tokenIndex, depth
									{
										position553 := position
										depth++
										{
											position554, tokenIndex554, depth554 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l555
											}
											goto l554
										l555:
											position, tokenIndex, depth = position554, tokenIndex554, depth554
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l546
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l546
													}
													if !rules[ruleexpression]() {
														goto l546
													}
													{
														position557, tokenIndex557, depth557 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l557
														}
														if !rules[rulevar]() {
															goto l557
														}
														goto l558
													l557:
														position, tokenIndex, depth = position557, tokenIndex557, depth557
													}
												l558:
													if !rules[ruleRPAREN]() {
														goto l546
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l546
													}
													break
												}
											}

										}
									l554:
										depth--
										add(rulegroupCondition, position553)
									}
									goto l545
								l546:
									position, tokenIndex, depth = position546, tokenIndex546, depth546
								}
								break
							default:
								{
									position559 := position
									depth++
									{
										position560, tokenIndex560, depth560 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l561
										}
										{
											position562, tokenIndex562, depth562 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l562
											}
											goto l563
										l562:
											position, tokenIndex, depth = position562, tokenIndex562, depth562
										}
									l563:
										goto l560
									l561:
										position, tokenIndex, depth = position560, tokenIndex560, depth560
										if !rules[ruleoffset]() {
											goto l453
										}
										{
											position564, tokenIndex564, depth564 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l564
											}
											goto l565
										l564:
											position, tokenIndex, depth = position564, tokenIndex564, depth564
										}
									l565:
									}
								l560:
									depth--
									add(rulelimitOffsetClauses, position559)
								}
								break
							}
						}

					}
				l455:
					goto l454
				l453:
					position, tokenIndex, depth = position453, tokenIndex453, depth453
				}
			l454:
				depth--
				add(rulesolutionModifier, position452)
			}
			return true
		},
		/* 48 groupCondition <- <(functionCall / ((&('$' | '?') var) | (&('(') (LPAREN expression (AS var)? RPAREN)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'H' | 'I' | 'L' | 'M' | 'N' | 'R' | 'S' | 'T' | 'U' | 'Y' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'h' | 'i' | 'l' | 'm' | 'n' | 'r' | 's' | 't' | 'u' | 'y') builtinCall)))> */
		nil,
		/* 49 orderCondition <- <(((ASC / DESC)? brackettedExpression) / functionCall / builtinCall / var)> */
		nil,
		/* 50 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 51 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position569, tokenIndex569, depth569 := position, tokenIndex, depth
			{
				position570 := position
				depth++
				{
					position571 := position
					depth++
					{
						position572, tokenIndex572, depth572 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l573
						}
						position++
						goto l572
					l573:
						position, tokenIndex, depth = position572, tokenIndex572, depth572
						if buffer[position] != rune('L') {
							goto l569
						}
						position++
					}
				l572:
					{
						position574, tokenIndex574, depth574 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l575
						}
						position++
						goto l574
					l575:
						position, tokenIndex, depth = position574, tokenIndex574, depth574
						if buffer[position] != rune('I') {
							goto l569
						}
						position++
					}
				l574:
					{
						position576, tokenIndex576, depth576 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l577
						}
						position++
						goto l576
					l577:
						position, tokenIndex, depth = position576, tokenIndex576, depth576
						if buffer[position] != rune('M') {
							goto l569
						}
						position++
					}
				l576:
					{
						position578, tokenIndex578, depth578 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l579
						}
						position++
						goto l578
					l579:
						position, tokenIndex, depth = position578, tokenIndex578, depth578
						if buffer[position] != rune('I') {
							goto l569
						}
						position++
					}
				l578:
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l581
						}
						position++
						goto l580
					l581:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
						if buffer[position] != rune('T') {
							goto l569
						}
						position++
					}
				l580:
					if !rules[ruleskip]() {
						goto l569
					}
					depth--
					add(ruleLIMIT, position571)
				}
				if !rules[ruleINTEGER]() {
					goto l569
				}
				depth--
				add(rulelimit, position570)
			}
			return true
		l569:
			position, tokenIndex, depth = position569, tokenIndex569, depth569
			return false
		},
		/* 52 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position582, tokenIndex582, depth582 := position, tokenIndex, depth
			{
				position583 := position
				depth++
				{
					position584 := position
					depth++
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
							goto l582
						}
						position++
					}
				l585:
					{
						position587, tokenIndex587, depth587 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l588
						}
						position++
						goto l587
					l588:
						position, tokenIndex, depth = position587, tokenIndex587, depth587
						if buffer[position] != rune('F') {
							goto l582
						}
						position++
					}
				l587:
					{
						position589, tokenIndex589, depth589 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l590
						}
						position++
						goto l589
					l590:
						position, tokenIndex, depth = position589, tokenIndex589, depth589
						if buffer[position] != rune('F') {
							goto l582
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
							goto l582
						}
						position++
					}
				l591:
					{
						position593, tokenIndex593, depth593 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l594
						}
						position++
						goto l593
					l594:
						position, tokenIndex, depth = position593, tokenIndex593, depth593
						if buffer[position] != rune('E') {
							goto l582
						}
						position++
					}
				l593:
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l596
						}
						position++
						goto l595
					l596:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
						if buffer[position] != rune('T') {
							goto l582
						}
						position++
					}
				l595:
					if !rules[ruleskip]() {
						goto l582
					}
					depth--
					add(ruleOFFSET, position584)
				}
				if !rules[ruleINTEGER]() {
					goto l582
				}
				depth--
				add(ruleoffset, position583)
			}
			return true
		l582:
			position, tokenIndex, depth = position582, tokenIndex582, depth582
			return false
		},
		/* 53 expression <- <conditionalOrExpression> */
		func() bool {
			position597, tokenIndex597, depth597 := position, tokenIndex, depth
			{
				position598 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l597
				}
				depth--
				add(ruleexpression, position598)
			}
			return true
		l597:
			position, tokenIndex, depth = position597, tokenIndex597, depth597
			return false
		},
		/* 54 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l599
				}
				{
					position601, tokenIndex601, depth601 := position, tokenIndex, depth
					{
						position603 := position
						depth++
						if buffer[position] != rune('|') {
							goto l601
						}
						position++
						if buffer[position] != rune('|') {
							goto l601
						}
						position++
						if !rules[ruleskip]() {
							goto l601
						}
						depth--
						add(ruleOR, position603)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l601
					}
					goto l602
				l601:
					position, tokenIndex, depth = position601, tokenIndex601, depth601
				}
			l602:
				depth--
				add(ruleconditionalOrExpression, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 55 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position604, tokenIndex604, depth604 := position, tokenIndex, depth
			{
				position605 := position
				depth++
				{
					position606 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l604
					}
					{
						position607, tokenIndex607, depth607 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position610 := position
									depth++
									{
										position611 := position
										depth++
										{
											position612, tokenIndex612, depth612 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l613
											}
											position++
											goto l612
										l613:
											position, tokenIndex, depth = position612, tokenIndex612, depth612
											if buffer[position] != rune('N') {
												goto l607
											}
											position++
										}
									l612:
										{
											position614, tokenIndex614, depth614 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l615
											}
											position++
											goto l614
										l615:
											position, tokenIndex, depth = position614, tokenIndex614, depth614
											if buffer[position] != rune('O') {
												goto l607
											}
											position++
										}
									l614:
										{
											position616, tokenIndex616, depth616 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l617
											}
											position++
											goto l616
										l617:
											position, tokenIndex, depth = position616, tokenIndex616, depth616
											if buffer[position] != rune('T') {
												goto l607
											}
											position++
										}
									l616:
										if buffer[position] != rune(' ') {
											goto l607
										}
										position++
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
												goto l607
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
												goto l607
											}
											position++
										}
									l620:
										if !rules[ruleskip]() {
											goto l607
										}
										depth--
										add(ruleNOTIN, position611)
									}
									if !rules[ruleargList]() {
										goto l607
									}
									depth--
									add(rulenotin, position610)
								}
								break
							case 'I', 'i':
								{
									position622 := position
									depth++
									{
										position623 := position
										depth++
										{
											position624, tokenIndex624, depth624 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l625
											}
											position++
											goto l624
										l625:
											position, tokenIndex, depth = position624, tokenIndex624, depth624
											if buffer[position] != rune('I') {
												goto l607
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
												goto l607
											}
											position++
										}
									l626:
										if !rules[ruleskip]() {
											goto l607
										}
										depth--
										add(ruleIN, position623)
									}
									if !rules[ruleargList]() {
										goto l607
									}
									depth--
									add(rulein, position622)
								}
								break
							default:
								{
									position628, tokenIndex628, depth628 := position, tokenIndex, depth
									{
										position630 := position
										depth++
										if buffer[position] != rune('<') {
											goto l629
										}
										position++
										if !rules[ruleskip]() {
											goto l629
										}
										depth--
										add(ruleLT, position630)
									}
									goto l628
								l629:
									position, tokenIndex, depth = position628, tokenIndex628, depth628
									{
										position632 := position
										depth++
										if buffer[position] != rune('>') {
											goto l631
										}
										position++
										if buffer[position] != rune('=') {
											goto l631
										}
										position++
										if !rules[ruleskip]() {
											goto l631
										}
										depth--
										add(ruleGE, position632)
									}
									goto l628
								l631:
									position, tokenIndex, depth = position628, tokenIndex628, depth628
									{
										switch buffer[position] {
										case '>':
											{
												position634 := position
												depth++
												if buffer[position] != rune('>') {
													goto l607
												}
												position++
												if !rules[ruleskip]() {
													goto l607
												}
												depth--
												add(ruleGT, position634)
											}
											break
										case '<':
											{
												position635 := position
												depth++
												if buffer[position] != rune('<') {
													goto l607
												}
												position++
												if buffer[position] != rune('=') {
													goto l607
												}
												position++
												if !rules[ruleskip]() {
													goto l607
												}
												depth--
												add(ruleLE, position635)
											}
											break
										case '!':
											{
												position636 := position
												depth++
												if buffer[position] != rune('!') {
													goto l607
												}
												position++
												if buffer[position] != rune('=') {
													goto l607
												}
												position++
												if !rules[ruleskip]() {
													goto l607
												}
												depth--
												add(ruleNE, position636)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l607
											}
											break
										}
									}

								}
							l628:
								if !rules[rulenumericExpression]() {
									goto l607
								}
								break
							}
						}

						goto l608
					l607:
						position, tokenIndex, depth = position607, tokenIndex607, depth607
					}
				l608:
					depth--
					add(rulevalueLogical, position606)
				}
				{
					position637, tokenIndex637, depth637 := position, tokenIndex, depth
					{
						position639 := position
						depth++
						if buffer[position] != rune('&') {
							goto l637
						}
						position++
						if buffer[position] != rune('&') {
							goto l637
						}
						position++
						if !rules[ruleskip]() {
							goto l637
						}
						depth--
						add(ruleAND, position639)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l637
					}
					goto l638
				l637:
					position, tokenIndex, depth = position637, tokenIndex637, depth637
				}
			l638:
				depth--
				add(ruleconditionalAndExpression, position605)
			}
			return true
		l604:
			position, tokenIndex, depth = position604, tokenIndex604, depth604
			return false
		},
		/* 56 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 57 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position641, tokenIndex641, depth641 := position, tokenIndex, depth
			{
				position642 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l641
				}
			l643:
				{
					position644, tokenIndex644, depth644 := position, tokenIndex, depth
					{
						position645, tokenIndex645, depth645 := position, tokenIndex, depth
						{
							position647, tokenIndex647, depth647 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l648
							}
							goto l647
						l648:
							position, tokenIndex, depth = position647, tokenIndex647, depth647
							if !rules[ruleMINUS]() {
								goto l646
							}
						}
					l647:
						if !rules[rulemultiplicativeExpression]() {
							goto l646
						}
						goto l645
					l646:
						position, tokenIndex, depth = position645, tokenIndex645, depth645
						{
							position649 := position
							depth++
							{
								position650, tokenIndex650, depth650 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l651
								}
								position++
								goto l650
							l651:
								position, tokenIndex, depth = position650, tokenIndex650, depth650
								if buffer[position] != rune('-') {
									goto l644
								}
								position++
							}
						l650:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l644
							}
							position++
						l652:
							{
								position653, tokenIndex653, depth653 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l653
								}
								position++
								goto l652
							l653:
								position, tokenIndex, depth = position653, tokenIndex653, depth653
							}
							{
								position654, tokenIndex654, depth654 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l654
								}
								position++
							l656:
								{
									position657, tokenIndex657, depth657 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l657
									}
									position++
									goto l656
								l657:
									position, tokenIndex, depth = position657, tokenIndex657, depth657
								}
								goto l655
							l654:
								position, tokenIndex, depth = position654, tokenIndex654, depth654
							}
						l655:
							if !rules[ruleskip]() {
								goto l644
							}
							depth--
							add(rulesignedNumericLiteral, position649)
						}
					}
				l645:
					goto l643
				l644:
					position, tokenIndex, depth = position644, tokenIndex644, depth644
				}
				depth--
				add(rulenumericExpression, position642)
			}
			return true
		l641:
			position, tokenIndex, depth = position641, tokenIndex641, depth641
			return false
		},
		/* 58 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position658, tokenIndex658, depth658 := position, tokenIndex, depth
			{
				position659 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l658
				}
			l660:
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					{
						position662, tokenIndex662, depth662 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l663
						}
						goto l662
					l663:
						position, tokenIndex, depth = position662, tokenIndex662, depth662
						if !rules[ruleSLASH]() {
							goto l661
						}
					}
				l662:
					if !rules[ruleunaryExpression]() {
						goto l661
					}
					goto l660
				l661:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
				}
				depth--
				add(rulemultiplicativeExpression, position659)
			}
			return true
		l658:
			position, tokenIndex, depth = position658, tokenIndex658, depth658
			return false
		},
		/* 59 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position664, tokenIndex664, depth664 := position, tokenIndex, depth
			{
				position665 := position
				depth++
				{
					position666, tokenIndex666, depth666 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l666
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l666
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l666
							}
							break
						}
					}

					goto l667
				l666:
					position, tokenIndex, depth = position666, tokenIndex666, depth666
				}
			l667:
				{
					position669 := position
					depth++
					{
						position670, tokenIndex670, depth670 := position, tokenIndex, depth
						if !rules[rulebuiltinCall]() {
							goto l671
						}
						goto l670
					l671:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
						if !rules[rulefunctionCall]() {
							goto l672
						}
						goto l670
					l672:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
						if !rules[ruleiriref]() {
							goto l673
						}
						goto l670
					l673:
						position, tokenIndex, depth = position670, tokenIndex670, depth670
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position675 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position677 := position
												depth++
												{
													position678 := position
													depth++
													{
														position679, tokenIndex679, depth679 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l680
														}
														position++
														goto l679
													l680:
														position, tokenIndex, depth = position679, tokenIndex679, depth679
														if buffer[position] != rune('G') {
															goto l664
														}
														position++
													}
												l679:
													{
														position681, tokenIndex681, depth681 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l682
														}
														position++
														goto l681
													l682:
														position, tokenIndex, depth = position681, tokenIndex681, depth681
														if buffer[position] != rune('R') {
															goto l664
														}
														position++
													}
												l681:
													{
														position683, tokenIndex683, depth683 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l684
														}
														position++
														goto l683
													l684:
														position, tokenIndex, depth = position683, tokenIndex683, depth683
														if buffer[position] != rune('O') {
															goto l664
														}
														position++
													}
												l683:
													{
														position685, tokenIndex685, depth685 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l686
														}
														position++
														goto l685
													l686:
														position, tokenIndex, depth = position685, tokenIndex685, depth685
														if buffer[position] != rune('U') {
															goto l664
														}
														position++
													}
												l685:
													{
														position687, tokenIndex687, depth687 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l688
														}
														position++
														goto l687
													l688:
														position, tokenIndex, depth = position687, tokenIndex687, depth687
														if buffer[position] != rune('P') {
															goto l664
														}
														position++
													}
												l687:
													if buffer[position] != rune('_') {
														goto l664
													}
													position++
													{
														position689, tokenIndex689, depth689 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l690
														}
														position++
														goto l689
													l690:
														position, tokenIndex, depth = position689, tokenIndex689, depth689
														if buffer[position] != rune('C') {
															goto l664
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
															goto l664
														}
														position++
													}
												l691:
													{
														position693, tokenIndex693, depth693 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l694
														}
														position++
														goto l693
													l694:
														position, tokenIndex, depth = position693, tokenIndex693, depth693
														if buffer[position] != rune('N') {
															goto l664
														}
														position++
													}
												l693:
													{
														position695, tokenIndex695, depth695 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l696
														}
														position++
														goto l695
													l696:
														position, tokenIndex, depth = position695, tokenIndex695, depth695
														if buffer[position] != rune('C') {
															goto l664
														}
														position++
													}
												l695:
													{
														position697, tokenIndex697, depth697 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l698
														}
														position++
														goto l697
													l698:
														position, tokenIndex, depth = position697, tokenIndex697, depth697
														if buffer[position] != rune('A') {
															goto l664
														}
														position++
													}
												l697:
													{
														position699, tokenIndex699, depth699 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l700
														}
														position++
														goto l699
													l700:
														position, tokenIndex, depth = position699, tokenIndex699, depth699
														if buffer[position] != rune('T') {
															goto l664
														}
														position++
													}
												l699:
													if !rules[ruleskip]() {
														goto l664
													}
													depth--
													add(ruleGROUPCONCAT, position678)
												}
												if !rules[ruleLPAREN]() {
													goto l664
												}
												{
													position701, tokenIndex701, depth701 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l701
													}
													goto l702
												l701:
													position, tokenIndex, depth = position701, tokenIndex701, depth701
												}
											l702:
												if !rules[ruleexpression]() {
													goto l664
												}
												{
													position703, tokenIndex703, depth703 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l703
													}
													{
														position705 := position
														depth++
														{
															position706, tokenIndex706, depth706 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l707
															}
															position++
															goto l706
														l707:
															position, tokenIndex, depth = position706, tokenIndex706, depth706
															if buffer[position] != rune('S') {
																goto l703
															}
															position++
														}
													l706:
														{
															position708, tokenIndex708, depth708 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l709
															}
															position++
															goto l708
														l709:
															position, tokenIndex, depth = position708, tokenIndex708, depth708
															if buffer[position] != rune('E') {
																goto l703
															}
															position++
														}
													l708:
														{
															position710, tokenIndex710, depth710 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l711
															}
															position++
															goto l710
														l711:
															position, tokenIndex, depth = position710, tokenIndex710, depth710
															if buffer[position] != rune('P') {
																goto l703
															}
															position++
														}
													l710:
														{
															position712, tokenIndex712, depth712 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l713
															}
															position++
															goto l712
														l713:
															position, tokenIndex, depth = position712, tokenIndex712, depth712
															if buffer[position] != rune('A') {
																goto l703
															}
															position++
														}
													l712:
														{
															position714, tokenIndex714, depth714 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l715
															}
															position++
															goto l714
														l715:
															position, tokenIndex, depth = position714, tokenIndex714, depth714
															if buffer[position] != rune('R') {
																goto l703
															}
															position++
														}
													l714:
														{
															position716, tokenIndex716, depth716 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l717
															}
															position++
															goto l716
														l717:
															position, tokenIndex, depth = position716, tokenIndex716, depth716
															if buffer[position] != rune('A') {
																goto l703
															}
															position++
														}
													l716:
														{
															position718, tokenIndex718, depth718 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l719
															}
															position++
															goto l718
														l719:
															position, tokenIndex, depth = position718, tokenIndex718, depth718
															if buffer[position] != rune('T') {
																goto l703
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
																goto l703
															}
															position++
														}
													l720:
														{
															position722, tokenIndex722, depth722 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l723
															}
															position++
															goto l722
														l723:
															position, tokenIndex, depth = position722, tokenIndex722, depth722
															if buffer[position] != rune('R') {
																goto l703
															}
															position++
														}
													l722:
														if !rules[ruleskip]() {
															goto l703
														}
														depth--
														add(ruleSEPARATOR, position705)
													}
													if !rules[ruleEQ]() {
														goto l703
													}
													if !rules[rulestring]() {
														goto l703
													}
													goto l704
												l703:
													position, tokenIndex, depth = position703, tokenIndex703, depth703
												}
											l704:
												if !rules[ruleRPAREN]() {
													goto l664
												}
												depth--
												add(rulegroupConcat, position677)
											}
											break
										case 'C', 'c':
											{
												position724 := position
												depth++
												{
													position725 := position
													depth++
													{
														position726, tokenIndex726, depth726 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l727
														}
														position++
														goto l726
													l727:
														position, tokenIndex, depth = position726, tokenIndex726, depth726
														if buffer[position] != rune('C') {
															goto l664
														}
														position++
													}
												l726:
													{
														position728, tokenIndex728, depth728 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l729
														}
														position++
														goto l728
													l729:
														position, tokenIndex, depth = position728, tokenIndex728, depth728
														if buffer[position] != rune('O') {
															goto l664
														}
														position++
													}
												l728:
													{
														position730, tokenIndex730, depth730 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l731
														}
														position++
														goto l730
													l731:
														position, tokenIndex, depth = position730, tokenIndex730, depth730
														if buffer[position] != rune('U') {
															goto l664
														}
														position++
													}
												l730:
													{
														position732, tokenIndex732, depth732 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l733
														}
														position++
														goto l732
													l733:
														position, tokenIndex, depth = position732, tokenIndex732, depth732
														if buffer[position] != rune('N') {
															goto l664
														}
														position++
													}
												l732:
													{
														position734, tokenIndex734, depth734 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l735
														}
														position++
														goto l734
													l735:
														position, tokenIndex, depth = position734, tokenIndex734, depth734
														if buffer[position] != rune('T') {
															goto l664
														}
														position++
													}
												l734:
													if !rules[ruleskip]() {
														goto l664
													}
													depth--
													add(ruleCOUNT, position725)
												}
												if !rules[ruleLPAREN]() {
													goto l664
												}
												{
													position736, tokenIndex736, depth736 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l736
													}
													goto l737
												l736:
													position, tokenIndex, depth = position736, tokenIndex736, depth736
												}
											l737:
												{
													position738, tokenIndex738, depth738 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l739
													}
													goto l738
												l739:
													position, tokenIndex, depth = position738, tokenIndex738, depth738
													if !rules[ruleexpression]() {
														goto l664
													}
												}
											l738:
												if !rules[ruleRPAREN]() {
													goto l664
												}
												depth--
												add(rulecount, position724)
											}
											break
										default:
											{
												position740, tokenIndex740, depth740 := position, tokenIndex, depth
												{
													position742 := position
													depth++
													{
														position743, tokenIndex743, depth743 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l744
														}
														position++
														goto l743
													l744:
														position, tokenIndex, depth = position743, tokenIndex743, depth743
														if buffer[position] != rune('S') {
															goto l741
														}
														position++
													}
												l743:
													{
														position745, tokenIndex745, depth745 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l746
														}
														position++
														goto l745
													l746:
														position, tokenIndex, depth = position745, tokenIndex745, depth745
														if buffer[position] != rune('U') {
															goto l741
														}
														position++
													}
												l745:
													{
														position747, tokenIndex747, depth747 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l748
														}
														position++
														goto l747
													l748:
														position, tokenIndex, depth = position747, tokenIndex747, depth747
														if buffer[position] != rune('M') {
															goto l741
														}
														position++
													}
												l747:
													if !rules[ruleskip]() {
														goto l741
													}
													depth--
													add(ruleSUM, position742)
												}
												goto l740
											l741:
												position, tokenIndex, depth = position740, tokenIndex740, depth740
												{
													position750 := position
													depth++
													{
														position751, tokenIndex751, depth751 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l752
														}
														position++
														goto l751
													l752:
														position, tokenIndex, depth = position751, tokenIndex751, depth751
														if buffer[position] != rune('M') {
															goto l749
														}
														position++
													}
												l751:
													{
														position753, tokenIndex753, depth753 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l754
														}
														position++
														goto l753
													l754:
														position, tokenIndex, depth = position753, tokenIndex753, depth753
														if buffer[position] != rune('I') {
															goto l749
														}
														position++
													}
												l753:
													{
														position755, tokenIndex755, depth755 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l756
														}
														position++
														goto l755
													l756:
														position, tokenIndex, depth = position755, tokenIndex755, depth755
														if buffer[position] != rune('N') {
															goto l749
														}
														position++
													}
												l755:
													if !rules[ruleskip]() {
														goto l749
													}
													depth--
													add(ruleMIN, position750)
												}
												goto l740
											l749:
												position, tokenIndex, depth = position740, tokenIndex740, depth740
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position758 := position
															depth++
															{
																position759, tokenIndex759, depth759 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l760
																}
																position++
																goto l759
															l760:
																position, tokenIndex, depth = position759, tokenIndex759, depth759
																if buffer[position] != rune('S') {
																	goto l664
																}
																position++
															}
														l759:
															{
																position761, tokenIndex761, depth761 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l762
																}
																position++
																goto l761
															l762:
																position, tokenIndex, depth = position761, tokenIndex761, depth761
																if buffer[position] != rune('A') {
																	goto l664
																}
																position++
															}
														l761:
															{
																position763, tokenIndex763, depth763 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l764
																}
																position++
																goto l763
															l764:
																position, tokenIndex, depth = position763, tokenIndex763, depth763
																if buffer[position] != rune('M') {
																	goto l664
																}
																position++
															}
														l763:
															{
																position765, tokenIndex765, depth765 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l766
																}
																position++
																goto l765
															l766:
																position, tokenIndex, depth = position765, tokenIndex765, depth765
																if buffer[position] != rune('P') {
																	goto l664
																}
																position++
															}
														l765:
															{
																position767, tokenIndex767, depth767 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l768
																}
																position++
																goto l767
															l768:
																position, tokenIndex, depth = position767, tokenIndex767, depth767
																if buffer[position] != rune('L') {
																	goto l664
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
																	goto l664
																}
																position++
															}
														l769:
															if !rules[ruleskip]() {
																goto l664
															}
															depth--
															add(ruleSAMPLE, position758)
														}
														break
													case 'A', 'a':
														{
															position771 := position
															depth++
															{
																position772, tokenIndex772, depth772 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l773
																}
																position++
																goto l772
															l773:
																position, tokenIndex, depth = position772, tokenIndex772, depth772
																if buffer[position] != rune('A') {
																	goto l664
																}
																position++
															}
														l772:
															{
																position774, tokenIndex774, depth774 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l775
																}
																position++
																goto l774
															l775:
																position, tokenIndex, depth = position774, tokenIndex774, depth774
																if buffer[position] != rune('V') {
																	goto l664
																}
																position++
															}
														l774:
															{
																position776, tokenIndex776, depth776 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l777
																}
																position++
																goto l776
															l777:
																position, tokenIndex, depth = position776, tokenIndex776, depth776
																if buffer[position] != rune('G') {
																	goto l664
																}
																position++
															}
														l776:
															if !rules[ruleskip]() {
																goto l664
															}
															depth--
															add(ruleAVG, position771)
														}
														break
													default:
														{
															position778 := position
															depth++
															{
																position779, tokenIndex779, depth779 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l780
																}
																position++
																goto l779
															l780:
																position, tokenIndex, depth = position779, tokenIndex779, depth779
																if buffer[position] != rune('M') {
																	goto l664
																}
																position++
															}
														l779:
															{
																position781, tokenIndex781, depth781 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l782
																}
																position++
																goto l781
															l782:
																position, tokenIndex, depth = position781, tokenIndex781, depth781
																if buffer[position] != rune('A') {
																	goto l664
																}
																position++
															}
														l781:
															{
																position783, tokenIndex783, depth783 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l784
																}
																position++
																goto l783
															l784:
																position, tokenIndex, depth = position783, tokenIndex783, depth783
																if buffer[position] != rune('X') {
																	goto l664
																}
																position++
															}
														l783:
															if !rules[ruleskip]() {
																goto l664
															}
															depth--
															add(ruleMAX, position778)
														}
														break
													}
												}

											}
										l740:
											if !rules[ruleLPAREN]() {
												goto l664
											}
											{
												position785, tokenIndex785, depth785 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l785
												}
												goto l786
											l785:
												position, tokenIndex, depth = position785, tokenIndex785, depth785
											}
										l786:
											if !rules[ruleexpression]() {
												goto l664
											}
											if !rules[ruleRPAREN]() {
												goto l664
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position675)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l664
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l664
								}
								break
							case '(':
								if !rules[rulebrackettedExpression]() {
									goto l664
								}
								break
							case '"', '\'':
								if !rules[ruleliteral]() {
									goto l664
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l664
								}
								break
							}
						}

					}
				l670:
					depth--
					add(ruleprimaryExpression, position669)
				}
				depth--
				add(ruleunaryExpression, position665)
			}
			return true
		l664:
			position, tokenIndex, depth = position664, tokenIndex664, depth664
			return false
		},
		/* 60 primaryExpression <- <(builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('(') brackettedExpression) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 61 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position788, tokenIndex788, depth788 := position, tokenIndex, depth
			{
				position789 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l788
				}
				if !rules[ruleexpression]() {
					goto l788
				}
				if !rules[ruleRPAREN]() {
					goto l788
				}
				depth--
				add(rulebrackettedExpression, position789)
			}
			return true
		l788:
			position, tokenIndex, depth = position788, tokenIndex788, depth788
			return false
		},
		/* 62 functionCall <- <(iriref argList)> */
		func() bool {
			position790, tokenIndex790, depth790 := position, tokenIndex, depth
			{
				position791 := position
				depth++
				if !rules[ruleiriref]() {
					goto l790
				}
				if !rules[ruleargList]() {
					goto l790
				}
				depth--
				add(rulefunctionCall, position791)
			}
			return true
		l790:
			position, tokenIndex, depth = position790, tokenIndex790, depth790
			return false
		},
		/* 63 in <- <(IN argList)> */
		nil,
		/* 64 notin <- <(NOTIN argList)> */
		nil,
		/* 65 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position794, tokenIndex794, depth794 := position, tokenIndex, depth
			{
				position795 := position
				depth++
				{
					position796, tokenIndex796, depth796 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l797
					}
					goto l796
				l797:
					position, tokenIndex, depth = position796, tokenIndex796, depth796
					if !rules[ruleLPAREN]() {
						goto l794
					}
					if !rules[ruleexpression]() {
						goto l794
					}
				l798:
					{
						position799, tokenIndex799, depth799 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l799
						}
						if !rules[ruleexpression]() {
							goto l799
						}
						goto l798
					l799:
						position, tokenIndex, depth = position799, tokenIndex799, depth799
					}
					if !rules[ruleRPAREN]() {
						goto l794
					}
				}
			l796:
				depth--
				add(ruleargList, position795)
			}
			return true
		l794:
			position, tokenIndex, depth = position794, tokenIndex794, depth794
			return false
		},
		/* 66 aggregate <- <((&('G' | 'g') groupConcat) | (&('C' | 'c') count) | (&('A' | 'M' | 'S' | 'a' | 'm' | 's') ((SUM / MIN / ((&('S' | 's') SAMPLE) | (&('A' | 'a') AVG) | (&('M' | 'm') MAX))) LPAREN DISTINCT? expression RPAREN)))> */
		nil,
		/* 67 count <- <(COUNT LPAREN DISTINCT? (STAR / expression) RPAREN)> */
		nil,
		/* 68 groupConcat <- <(GROUPCONCAT LPAREN DISTINCT? expression (SEMICOLON SEPARATOR EQ string)? RPAREN)> */
		nil,
		/* 69 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		func() bool {
			position803, tokenIndex803, depth803 := position, tokenIndex, depth
			{
				position804 := position
				depth++
				{
					position805, tokenIndex805, depth805 := position, tokenIndex, depth
					{
						position807, tokenIndex807, depth807 := position, tokenIndex, depth
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
							if !rules[ruleskip]() {
								goto l808
							}
							depth--
							add(ruleSTR, position809)
						}
						goto l807
					l808:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position817 := position
							depth++
							{
								position818, tokenIndex818, depth818 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l819
								}
								position++
								goto l818
							l819:
								position, tokenIndex, depth = position818, tokenIndex818, depth818
								if buffer[position] != rune('L') {
									goto l816
								}
								position++
							}
						l818:
							{
								position820, tokenIndex820, depth820 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l821
								}
								position++
								goto l820
							l821:
								position, tokenIndex, depth = position820, tokenIndex820, depth820
								if buffer[position] != rune('A') {
									goto l816
								}
								position++
							}
						l820:
							{
								position822, tokenIndex822, depth822 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l823
								}
								position++
								goto l822
							l823:
								position, tokenIndex, depth = position822, tokenIndex822, depth822
								if buffer[position] != rune('N') {
									goto l816
								}
								position++
							}
						l822:
							{
								position824, tokenIndex824, depth824 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l825
								}
								position++
								goto l824
							l825:
								position, tokenIndex, depth = position824, tokenIndex824, depth824
								if buffer[position] != rune('G') {
									goto l816
								}
								position++
							}
						l824:
							if !rules[ruleskip]() {
								goto l816
							}
							depth--
							add(ruleLANG, position817)
						}
						goto l807
					l816:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position827 := position
							depth++
							{
								position828, tokenIndex828, depth828 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l829
								}
								position++
								goto l828
							l829:
								position, tokenIndex, depth = position828, tokenIndex828, depth828
								if buffer[position] != rune('D') {
									goto l826
								}
								position++
							}
						l828:
							{
								position830, tokenIndex830, depth830 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l831
								}
								position++
								goto l830
							l831:
								position, tokenIndex, depth = position830, tokenIndex830, depth830
								if buffer[position] != rune('A') {
									goto l826
								}
								position++
							}
						l830:
							{
								position832, tokenIndex832, depth832 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l833
								}
								position++
								goto l832
							l833:
								position, tokenIndex, depth = position832, tokenIndex832, depth832
								if buffer[position] != rune('T') {
									goto l826
								}
								position++
							}
						l832:
							{
								position834, tokenIndex834, depth834 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l835
								}
								position++
								goto l834
							l835:
								position, tokenIndex, depth = position834, tokenIndex834, depth834
								if buffer[position] != rune('A') {
									goto l826
								}
								position++
							}
						l834:
							{
								position836, tokenIndex836, depth836 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l837
								}
								position++
								goto l836
							l837:
								position, tokenIndex, depth = position836, tokenIndex836, depth836
								if buffer[position] != rune('T') {
									goto l826
								}
								position++
							}
						l836:
							{
								position838, tokenIndex838, depth838 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l839
								}
								position++
								goto l838
							l839:
								position, tokenIndex, depth = position838, tokenIndex838, depth838
								if buffer[position] != rune('Y') {
									goto l826
								}
								position++
							}
						l838:
							{
								position840, tokenIndex840, depth840 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l841
								}
								position++
								goto l840
							l841:
								position, tokenIndex, depth = position840, tokenIndex840, depth840
								if buffer[position] != rune('P') {
									goto l826
								}
								position++
							}
						l840:
							{
								position842, tokenIndex842, depth842 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l843
								}
								position++
								goto l842
							l843:
								position, tokenIndex, depth = position842, tokenIndex842, depth842
								if buffer[position] != rune('E') {
									goto l826
								}
								position++
							}
						l842:
							if !rules[ruleskip]() {
								goto l826
							}
							depth--
							add(ruleDATATYPE, position827)
						}
						goto l807
					l826:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position845 := position
							depth++
							{
								position846, tokenIndex846, depth846 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l847
								}
								position++
								goto l846
							l847:
								position, tokenIndex, depth = position846, tokenIndex846, depth846
								if buffer[position] != rune('I') {
									goto l844
								}
								position++
							}
						l846:
							{
								position848, tokenIndex848, depth848 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l849
								}
								position++
								goto l848
							l849:
								position, tokenIndex, depth = position848, tokenIndex848, depth848
								if buffer[position] != rune('R') {
									goto l844
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
									goto l844
								}
								position++
							}
						l850:
							if !rules[ruleskip]() {
								goto l844
							}
							depth--
							add(ruleIRI, position845)
						}
						goto l807
					l844:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position853 := position
							depth++
							{
								position854, tokenIndex854, depth854 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l855
								}
								position++
								goto l854
							l855:
								position, tokenIndex, depth = position854, tokenIndex854, depth854
								if buffer[position] != rune('U') {
									goto l852
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
									goto l852
								}
								position++
							}
						l856:
							{
								position858, tokenIndex858, depth858 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('I') {
									goto l852
								}
								position++
							}
						l858:
							if !rules[ruleskip]() {
								goto l852
							}
							depth--
							add(ruleURI, position853)
						}
						goto l807
					l852:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position861 := position
							depth++
							{
								position862, tokenIndex862, depth862 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l863
								}
								position++
								goto l862
							l863:
								position, tokenIndex, depth = position862, tokenIndex862, depth862
								if buffer[position] != rune('S') {
									goto l860
								}
								position++
							}
						l862:
							{
								position864, tokenIndex864, depth864 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l865
								}
								position++
								goto l864
							l865:
								position, tokenIndex, depth = position864, tokenIndex864, depth864
								if buffer[position] != rune('T') {
									goto l860
								}
								position++
							}
						l864:
							{
								position866, tokenIndex866, depth866 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l867
								}
								position++
								goto l866
							l867:
								position, tokenIndex, depth = position866, tokenIndex866, depth866
								if buffer[position] != rune('R') {
									goto l860
								}
								position++
							}
						l866:
							{
								position868, tokenIndex868, depth868 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l869
								}
								position++
								goto l868
							l869:
								position, tokenIndex, depth = position868, tokenIndex868, depth868
								if buffer[position] != rune('L') {
									goto l860
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
									goto l860
								}
								position++
							}
						l870:
							{
								position872, tokenIndex872, depth872 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l873
								}
								position++
								goto l872
							l873:
								position, tokenIndex, depth = position872, tokenIndex872, depth872
								if buffer[position] != rune('N') {
									goto l860
								}
								position++
							}
						l872:
							if !rules[ruleskip]() {
								goto l860
							}
							depth--
							add(ruleSTRLEN, position861)
						}
						goto l807
					l860:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position875 := position
							depth++
							{
								position876, tokenIndex876, depth876 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l877
								}
								position++
								goto l876
							l877:
								position, tokenIndex, depth = position876, tokenIndex876, depth876
								if buffer[position] != rune('M') {
									goto l874
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
									goto l874
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
									goto l874
								}
								position++
							}
						l880:
							{
								position882, tokenIndex882, depth882 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l883
								}
								position++
								goto l882
							l883:
								position, tokenIndex, depth = position882, tokenIndex882, depth882
								if buffer[position] != rune('T') {
									goto l874
								}
								position++
							}
						l882:
							{
								position884, tokenIndex884, depth884 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l885
								}
								position++
								goto l884
							l885:
								position, tokenIndex, depth = position884, tokenIndex884, depth884
								if buffer[position] != rune('H') {
									goto l874
								}
								position++
							}
						l884:
							if !rules[ruleskip]() {
								goto l874
							}
							depth--
							add(ruleMONTH, position875)
						}
						goto l807
					l874:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position887 := position
							depth++
							{
								position888, tokenIndex888, depth888 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l889
								}
								position++
								goto l888
							l889:
								position, tokenIndex, depth = position888, tokenIndex888, depth888
								if buffer[position] != rune('M') {
									goto l886
								}
								position++
							}
						l888:
							{
								position890, tokenIndex890, depth890 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l891
								}
								position++
								goto l890
							l891:
								position, tokenIndex, depth = position890, tokenIndex890, depth890
								if buffer[position] != rune('I') {
									goto l886
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
									goto l886
								}
								position++
							}
						l892:
							{
								position894, tokenIndex894, depth894 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l895
								}
								position++
								goto l894
							l895:
								position, tokenIndex, depth = position894, tokenIndex894, depth894
								if buffer[position] != rune('U') {
									goto l886
								}
								position++
							}
						l894:
							{
								position896, tokenIndex896, depth896 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l897
								}
								position++
								goto l896
							l897:
								position, tokenIndex, depth = position896, tokenIndex896, depth896
								if buffer[position] != rune('T') {
									goto l886
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
									goto l886
								}
								position++
							}
						l898:
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
									goto l886
								}
								position++
							}
						l900:
							if !rules[ruleskip]() {
								goto l886
							}
							depth--
							add(ruleMINUTES, position887)
						}
						goto l807
					l886:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position903 := position
							depth++
							{
								position904, tokenIndex904, depth904 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l905
								}
								position++
								goto l904
							l905:
								position, tokenIndex, depth = position904, tokenIndex904, depth904
								if buffer[position] != rune('S') {
									goto l902
								}
								position++
							}
						l904:
							{
								position906, tokenIndex906, depth906 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l907
								}
								position++
								goto l906
							l907:
								position, tokenIndex, depth = position906, tokenIndex906, depth906
								if buffer[position] != rune('E') {
									goto l902
								}
								position++
							}
						l906:
							{
								position908, tokenIndex908, depth908 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l909
								}
								position++
								goto l908
							l909:
								position, tokenIndex, depth = position908, tokenIndex908, depth908
								if buffer[position] != rune('C') {
									goto l902
								}
								position++
							}
						l908:
							{
								position910, tokenIndex910, depth910 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l911
								}
								position++
								goto l910
							l911:
								position, tokenIndex, depth = position910, tokenIndex910, depth910
								if buffer[position] != rune('O') {
									goto l902
								}
								position++
							}
						l910:
							{
								position912, tokenIndex912, depth912 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l913
								}
								position++
								goto l912
							l913:
								position, tokenIndex, depth = position912, tokenIndex912, depth912
								if buffer[position] != rune('N') {
									goto l902
								}
								position++
							}
						l912:
							{
								position914, tokenIndex914, depth914 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l915
								}
								position++
								goto l914
							l915:
								position, tokenIndex, depth = position914, tokenIndex914, depth914
								if buffer[position] != rune('D') {
									goto l902
								}
								position++
							}
						l914:
							{
								position916, tokenIndex916, depth916 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l917
								}
								position++
								goto l916
							l917:
								position, tokenIndex, depth = position916, tokenIndex916, depth916
								if buffer[position] != rune('S') {
									goto l902
								}
								position++
							}
						l916:
							if !rules[ruleskip]() {
								goto l902
							}
							depth--
							add(ruleSECONDS, position903)
						}
						goto l807
					l902:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position919 := position
							depth++
							{
								position920, tokenIndex920, depth920 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l921
								}
								position++
								goto l920
							l921:
								position, tokenIndex, depth = position920, tokenIndex920, depth920
								if buffer[position] != rune('T') {
									goto l918
								}
								position++
							}
						l920:
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
									goto l918
								}
								position++
							}
						l922:
							{
								position924, tokenIndex924, depth924 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l925
								}
								position++
								goto l924
							l925:
								position, tokenIndex, depth = position924, tokenIndex924, depth924
								if buffer[position] != rune('M') {
									goto l918
								}
								position++
							}
						l924:
							{
								position926, tokenIndex926, depth926 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l927
								}
								position++
								goto l926
							l927:
								position, tokenIndex, depth = position926, tokenIndex926, depth926
								if buffer[position] != rune('E') {
									goto l918
								}
								position++
							}
						l926:
							{
								position928, tokenIndex928, depth928 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l929
								}
								position++
								goto l928
							l929:
								position, tokenIndex, depth = position928, tokenIndex928, depth928
								if buffer[position] != rune('Z') {
									goto l918
								}
								position++
							}
						l928:
							{
								position930, tokenIndex930, depth930 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l931
								}
								position++
								goto l930
							l931:
								position, tokenIndex, depth = position930, tokenIndex930, depth930
								if buffer[position] != rune('O') {
									goto l918
								}
								position++
							}
						l930:
							{
								position932, tokenIndex932, depth932 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l933
								}
								position++
								goto l932
							l933:
								position, tokenIndex, depth = position932, tokenIndex932, depth932
								if buffer[position] != rune('N') {
									goto l918
								}
								position++
							}
						l932:
							{
								position934, tokenIndex934, depth934 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l935
								}
								position++
								goto l934
							l935:
								position, tokenIndex, depth = position934, tokenIndex934, depth934
								if buffer[position] != rune('E') {
									goto l918
								}
								position++
							}
						l934:
							if !rules[ruleskip]() {
								goto l918
							}
							depth--
							add(ruleTIMEZONE, position919)
						}
						goto l807
					l918:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
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
								if buffer[position] != rune('h') {
									goto l941
								}
								position++
								goto l940
							l941:
								position, tokenIndex, depth = position940, tokenIndex940, depth940
								if buffer[position] != rune('H') {
									goto l936
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
									goto l936
								}
								position++
							}
						l942:
							if buffer[position] != rune('1') {
								goto l936
							}
							position++
							if !rules[ruleskip]() {
								goto l936
							}
							depth--
							add(ruleSHA1, position937)
						}
						goto l807
					l936:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position945 := position
							depth++
							{
								position946, tokenIndex946, depth946 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l947
								}
								position++
								goto l946
							l947:
								position, tokenIndex, depth = position946, tokenIndex946, depth946
								if buffer[position] != rune('S') {
									goto l944
								}
								position++
							}
						l946:
							{
								position948, tokenIndex948, depth948 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l949
								}
								position++
								goto l948
							l949:
								position, tokenIndex, depth = position948, tokenIndex948, depth948
								if buffer[position] != rune('H') {
									goto l944
								}
								position++
							}
						l948:
							{
								position950, tokenIndex950, depth950 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l951
								}
								position++
								goto l950
							l951:
								position, tokenIndex, depth = position950, tokenIndex950, depth950
								if buffer[position] != rune('A') {
									goto l944
								}
								position++
							}
						l950:
							if buffer[position] != rune('2') {
								goto l944
							}
							position++
							if buffer[position] != rune('5') {
								goto l944
							}
							position++
							if buffer[position] != rune('6') {
								goto l944
							}
							position++
							if !rules[ruleskip]() {
								goto l944
							}
							depth--
							add(ruleSHA256, position945)
						}
						goto l807
					l944:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
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
								if buffer[position] != rune('h') {
									goto l957
								}
								position++
								goto l956
							l957:
								position, tokenIndex, depth = position956, tokenIndex956, depth956
								if buffer[position] != rune('H') {
									goto l952
								}
								position++
							}
						l956:
							{
								position958, tokenIndex958, depth958 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l959
								}
								position++
								goto l958
							l959:
								position, tokenIndex, depth = position958, tokenIndex958, depth958
								if buffer[position] != rune('A') {
									goto l952
								}
								position++
							}
						l958:
							if buffer[position] != rune('3') {
								goto l952
							}
							position++
							if buffer[position] != rune('8') {
								goto l952
							}
							position++
							if buffer[position] != rune('4') {
								goto l952
							}
							position++
							if !rules[ruleskip]() {
								goto l952
							}
							depth--
							add(ruleSHA384, position953)
						}
						goto l807
					l952:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position961 := position
							depth++
							{
								position962, tokenIndex962, depth962 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l963
								}
								position++
								goto l962
							l963:
								position, tokenIndex, depth = position962, tokenIndex962, depth962
								if buffer[position] != rune('I') {
									goto l960
								}
								position++
							}
						l962:
							{
								position964, tokenIndex964, depth964 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l965
								}
								position++
								goto l964
							l965:
								position, tokenIndex, depth = position964, tokenIndex964, depth964
								if buffer[position] != rune('S') {
									goto l960
								}
								position++
							}
						l964:
							{
								position966, tokenIndex966, depth966 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l967
								}
								position++
								goto l966
							l967:
								position, tokenIndex, depth = position966, tokenIndex966, depth966
								if buffer[position] != rune('I') {
									goto l960
								}
								position++
							}
						l966:
							{
								position968, tokenIndex968, depth968 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l969
								}
								position++
								goto l968
							l969:
								position, tokenIndex, depth = position968, tokenIndex968, depth968
								if buffer[position] != rune('R') {
									goto l960
								}
								position++
							}
						l968:
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
									goto l960
								}
								position++
							}
						l970:
							if !rules[ruleskip]() {
								goto l960
							}
							depth--
							add(ruleISIRI, position961)
						}
						goto l807
					l960:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position973 := position
							depth++
							{
								position974, tokenIndex974, depth974 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l975
								}
								position++
								goto l974
							l975:
								position, tokenIndex, depth = position974, tokenIndex974, depth974
								if buffer[position] != rune('I') {
									goto l972
								}
								position++
							}
						l974:
							{
								position976, tokenIndex976, depth976 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l977
								}
								position++
								goto l976
							l977:
								position, tokenIndex, depth = position976, tokenIndex976, depth976
								if buffer[position] != rune('S') {
									goto l972
								}
								position++
							}
						l976:
							{
								position978, tokenIndex978, depth978 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l979
								}
								position++
								goto l978
							l979:
								position, tokenIndex, depth = position978, tokenIndex978, depth978
								if buffer[position] != rune('U') {
									goto l972
								}
								position++
							}
						l978:
							{
								position980, tokenIndex980, depth980 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l981
								}
								position++
								goto l980
							l981:
								position, tokenIndex, depth = position980, tokenIndex980, depth980
								if buffer[position] != rune('R') {
									goto l972
								}
								position++
							}
						l980:
							{
								position982, tokenIndex982, depth982 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l983
								}
								position++
								goto l982
							l983:
								position, tokenIndex, depth = position982, tokenIndex982, depth982
								if buffer[position] != rune('I') {
									goto l972
								}
								position++
							}
						l982:
							if !rules[ruleskip]() {
								goto l972
							}
							depth--
							add(ruleISURI, position973)
						}
						goto l807
					l972:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position985 := position
							depth++
							{
								position986, tokenIndex986, depth986 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l987
								}
								position++
								goto l986
							l987:
								position, tokenIndex, depth = position986, tokenIndex986, depth986
								if buffer[position] != rune('I') {
									goto l984
								}
								position++
							}
						l986:
							{
								position988, tokenIndex988, depth988 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l989
								}
								position++
								goto l988
							l989:
								position, tokenIndex, depth = position988, tokenIndex988, depth988
								if buffer[position] != rune('S') {
									goto l984
								}
								position++
							}
						l988:
							{
								position990, tokenIndex990, depth990 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l991
								}
								position++
								goto l990
							l991:
								position, tokenIndex, depth = position990, tokenIndex990, depth990
								if buffer[position] != rune('B') {
									goto l984
								}
								position++
							}
						l990:
							{
								position992, tokenIndex992, depth992 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l993
								}
								position++
								goto l992
							l993:
								position, tokenIndex, depth = position992, tokenIndex992, depth992
								if buffer[position] != rune('L') {
									goto l984
								}
								position++
							}
						l992:
							{
								position994, tokenIndex994, depth994 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l995
								}
								position++
								goto l994
							l995:
								position, tokenIndex, depth = position994, tokenIndex994, depth994
								if buffer[position] != rune('A') {
									goto l984
								}
								position++
							}
						l994:
							{
								position996, tokenIndex996, depth996 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l997
								}
								position++
								goto l996
							l997:
								position, tokenIndex, depth = position996, tokenIndex996, depth996
								if buffer[position] != rune('N') {
									goto l984
								}
								position++
							}
						l996:
							{
								position998, tokenIndex998, depth998 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l999
								}
								position++
								goto l998
							l999:
								position, tokenIndex, depth = position998, tokenIndex998, depth998
								if buffer[position] != rune('K') {
									goto l984
								}
								position++
							}
						l998:
							if !rules[ruleskip]() {
								goto l984
							}
							depth--
							add(ruleISBLANK, position985)
						}
						goto l807
					l984:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							position1001 := position
							depth++
							{
								position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1003
								}
								position++
								goto l1002
							l1003:
								position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
								if buffer[position] != rune('I') {
									goto l1000
								}
								position++
							}
						l1002:
							{
								position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1005
								}
								position++
								goto l1004
							l1005:
								position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
								if buffer[position] != rune('S') {
									goto l1000
								}
								position++
							}
						l1004:
							{
								position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1007
								}
								position++
								goto l1006
							l1007:
								position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
								if buffer[position] != rune('L') {
									goto l1000
								}
								position++
							}
						l1006:
							{
								position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1009
								}
								position++
								goto l1008
							l1009:
								position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
								if buffer[position] != rune('I') {
									goto l1000
								}
								position++
							}
						l1008:
							{
								position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1011
								}
								position++
								goto l1010
							l1011:
								position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
								if buffer[position] != rune('T') {
									goto l1000
								}
								position++
							}
						l1010:
							{
								position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1013
								}
								position++
								goto l1012
							l1013:
								position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
								if buffer[position] != rune('E') {
									goto l1000
								}
								position++
							}
						l1012:
							{
								position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1015
								}
								position++
								goto l1014
							l1015:
								position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
								if buffer[position] != rune('R') {
									goto l1000
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
									goto l1000
								}
								position++
							}
						l1016:
							{
								position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1019
								}
								position++
								goto l1018
							l1019:
								position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
								if buffer[position] != rune('L') {
									goto l1000
								}
								position++
							}
						l1018:
							if !rules[ruleskip]() {
								goto l1000
							}
							depth--
							add(ruleISLITERAL, position1001)
						}
						goto l807
					l1000:
						position, tokenIndex, depth = position807, tokenIndex807, depth807
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position1021 := position
									depth++
									{
										position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1023
										}
										position++
										goto l1022
									l1023:
										position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
										if buffer[position] != rune('I') {
											goto l806
										}
										position++
									}
								l1022:
									{
										position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1025
										}
										position++
										goto l1024
									l1025:
										position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
										if buffer[position] != rune('S') {
											goto l806
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
											goto l806
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
											goto l806
										}
										position++
									}
								l1028:
									{
										position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1031
										}
										position++
										goto l1030
									l1031:
										position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
										if buffer[position] != rune('M') {
											goto l806
										}
										position++
									}
								l1030:
									{
										position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1033
										}
										position++
										goto l1032
									l1033:
										position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
										if buffer[position] != rune('E') {
											goto l806
										}
										position++
									}
								l1032:
									{
										position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1035
										}
										position++
										goto l1034
									l1035:
										position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
										if buffer[position] != rune('R') {
											goto l806
										}
										position++
									}
								l1034:
									{
										position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1037
										}
										position++
										goto l1036
									l1037:
										position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
										if buffer[position] != rune('I') {
											goto l806
										}
										position++
									}
								l1036:
									{
										position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1039
										}
										position++
										goto l1038
									l1039:
										position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
										if buffer[position] != rune('C') {
											goto l806
										}
										position++
									}
								l1038:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleISNUMERIC, position1021)
								}
								break
							case 'S', 's':
								{
									position1040 := position
									depth++
									{
										position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1042
										}
										position++
										goto l1041
									l1042:
										position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
										if buffer[position] != rune('S') {
											goto l806
										}
										position++
									}
								l1041:
									{
										position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1044
										}
										position++
										goto l1043
									l1044:
										position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
										if buffer[position] != rune('H') {
											goto l806
										}
										position++
									}
								l1043:
									{
										position1045, tokenIndex1045, depth1045 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1046
										}
										position++
										goto l1045
									l1046:
										position, tokenIndex, depth = position1045, tokenIndex1045, depth1045
										if buffer[position] != rune('A') {
											goto l806
										}
										position++
									}
								l1045:
									if buffer[position] != rune('5') {
										goto l806
									}
									position++
									if buffer[position] != rune('1') {
										goto l806
									}
									position++
									if buffer[position] != rune('2') {
										goto l806
									}
									position++
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleSHA512, position1040)
								}
								break
							case 'M', 'm':
								{
									position1047 := position
									depth++
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
											goto l806
										}
										position++
									}
								l1048:
									{
										position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1051
										}
										position++
										goto l1050
									l1051:
										position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
										if buffer[position] != rune('D') {
											goto l806
										}
										position++
									}
								l1050:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleMD5, position1047)
								}
								break
							case 'T', 't':
								{
									position1052 := position
									depth++
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
											goto l806
										}
										position++
									}
								l1053:
									{
										position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1056
										}
										position++
										goto l1055
									l1056:
										position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
										if buffer[position] != rune('Z') {
											goto l806
										}
										position++
									}
								l1055:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleTZ, position1052)
								}
								break
							case 'H', 'h':
								{
									position1057 := position
									depth++
									{
										position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1059
										}
										position++
										goto l1058
									l1059:
										position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
										if buffer[position] != rune('H') {
											goto l806
										}
										position++
									}
								l1058:
									{
										position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1061
										}
										position++
										goto l1060
									l1061:
										position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
										if buffer[position] != rune('O') {
											goto l806
										}
										position++
									}
								l1060:
									{
										position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1063
										}
										position++
										goto l1062
									l1063:
										position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
										if buffer[position] != rune('U') {
											goto l806
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
											goto l806
										}
										position++
									}
								l1064:
									{
										position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1067
										}
										position++
										goto l1066
									l1067:
										position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
										if buffer[position] != rune('S') {
											goto l806
										}
										position++
									}
								l1066:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleHOURS, position1057)
								}
								break
							case 'D', 'd':
								{
									position1068 := position
									depth++
									{
										position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1070
										}
										position++
										goto l1069
									l1070:
										position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
										if buffer[position] != rune('D') {
											goto l806
										}
										position++
									}
								l1069:
									{
										position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1072
										}
										position++
										goto l1071
									l1072:
										position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
										if buffer[position] != rune('A') {
											goto l806
										}
										position++
									}
								l1071:
									{
										position1073, tokenIndex1073, depth1073 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1074
										}
										position++
										goto l1073
									l1074:
										position, tokenIndex, depth = position1073, tokenIndex1073, depth1073
										if buffer[position] != rune('Y') {
											goto l806
										}
										position++
									}
								l1073:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleDAY, position1068)
								}
								break
							case 'Y', 'y':
								{
									position1075 := position
									depth++
									{
										position1076, tokenIndex1076, depth1076 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1077
										}
										position++
										goto l1076
									l1077:
										position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
										if buffer[position] != rune('Y') {
											goto l806
										}
										position++
									}
								l1076:
									{
										position1078, tokenIndex1078, depth1078 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1079
										}
										position++
										goto l1078
									l1079:
										position, tokenIndex, depth = position1078, tokenIndex1078, depth1078
										if buffer[position] != rune('E') {
											goto l806
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
											goto l806
										}
										position++
									}
								l1080:
									{
										position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1083
										}
										position++
										goto l1082
									l1083:
										position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
										if buffer[position] != rune('R') {
											goto l806
										}
										position++
									}
								l1082:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleYEAR, position1075)
								}
								break
							case 'E', 'e':
								{
									position1084 := position
									depth++
									{
										position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1086
										}
										position++
										goto l1085
									l1086:
										position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
										if buffer[position] != rune('E') {
											goto l806
										}
										position++
									}
								l1085:
									{
										position1087, tokenIndex1087, depth1087 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1088
										}
										position++
										goto l1087
									l1088:
										position, tokenIndex, depth = position1087, tokenIndex1087, depth1087
										if buffer[position] != rune('N') {
											goto l806
										}
										position++
									}
								l1087:
									{
										position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1090
										}
										position++
										goto l1089
									l1090:
										position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
										if buffer[position] != rune('C') {
											goto l806
										}
										position++
									}
								l1089:
									{
										position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1092
										}
										position++
										goto l1091
									l1092:
										position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
										if buffer[position] != rune('O') {
											goto l806
										}
										position++
									}
								l1091:
									{
										position1093, tokenIndex1093, depth1093 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1094
										}
										position++
										goto l1093
									l1094:
										position, tokenIndex, depth = position1093, tokenIndex1093, depth1093
										if buffer[position] != rune('D') {
											goto l806
										}
										position++
									}
								l1093:
									{
										position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1096
										}
										position++
										goto l1095
									l1096:
										position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
										if buffer[position] != rune('E') {
											goto l806
										}
										position++
									}
								l1095:
									if buffer[position] != rune('_') {
										goto l806
									}
									position++
									{
										position1097, tokenIndex1097, depth1097 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1098
										}
										position++
										goto l1097
									l1098:
										position, tokenIndex, depth = position1097, tokenIndex1097, depth1097
										if buffer[position] != rune('F') {
											goto l806
										}
										position++
									}
								l1097:
									{
										position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1100
										}
										position++
										goto l1099
									l1100:
										position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
										if buffer[position] != rune('O') {
											goto l806
										}
										position++
									}
								l1099:
									{
										position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1102
										}
										position++
										goto l1101
									l1102:
										position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
										if buffer[position] != rune('R') {
											goto l806
										}
										position++
									}
								l1101:
									if buffer[position] != rune('_') {
										goto l806
									}
									position++
									{
										position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1104
										}
										position++
										goto l1103
									l1104:
										position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
										if buffer[position] != rune('U') {
											goto l806
										}
										position++
									}
								l1103:
									{
										position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1106
										}
										position++
										goto l1105
									l1106:
										position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
										if buffer[position] != rune('R') {
											goto l806
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
											goto l806
										}
										position++
									}
								l1107:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleENCODEFORURI, position1084)
								}
								break
							case 'L', 'l':
								{
									position1109 := position
									depth++
									{
										position1110, tokenIndex1110, depth1110 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1111
										}
										position++
										goto l1110
									l1111:
										position, tokenIndex, depth = position1110, tokenIndex1110, depth1110
										if buffer[position] != rune('L') {
											goto l806
										}
										position++
									}
								l1110:
									{
										position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1113
										}
										position++
										goto l1112
									l1113:
										position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
										if buffer[position] != rune('C') {
											goto l806
										}
										position++
									}
								l1112:
									{
										position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1115
										}
										position++
										goto l1114
									l1115:
										position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
										if buffer[position] != rune('A') {
											goto l806
										}
										position++
									}
								l1114:
									{
										position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1117
										}
										position++
										goto l1116
									l1117:
										position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
										if buffer[position] != rune('S') {
											goto l806
										}
										position++
									}
								l1116:
									{
										position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1119
										}
										position++
										goto l1118
									l1119:
										position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
										if buffer[position] != rune('E') {
											goto l806
										}
										position++
									}
								l1118:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleLCASE, position1109)
								}
								break
							case 'U', 'u':
								{
									position1120 := position
									depth++
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
											goto l806
										}
										position++
									}
								l1121:
									{
										position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1124
										}
										position++
										goto l1123
									l1124:
										position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
										if buffer[position] != rune('C') {
											goto l806
										}
										position++
									}
								l1123:
									{
										position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1126
										}
										position++
										goto l1125
									l1126:
										position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
										if buffer[position] != rune('A') {
											goto l806
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
											goto l806
										}
										position++
									}
								l1127:
									{
										position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1130
										}
										position++
										goto l1129
									l1130:
										position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
										if buffer[position] != rune('E') {
											goto l806
										}
										position++
									}
								l1129:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleUCASE, position1120)
								}
								break
							case 'F', 'f':
								{
									position1131 := position
									depth++
									{
										position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1133
										}
										position++
										goto l1132
									l1133:
										position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
										if buffer[position] != rune('F') {
											goto l806
										}
										position++
									}
								l1132:
									{
										position1134, tokenIndex1134, depth1134 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1135
										}
										position++
										goto l1134
									l1135:
										position, tokenIndex, depth = position1134, tokenIndex1134, depth1134
										if buffer[position] != rune('L') {
											goto l806
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
											goto l806
										}
										position++
									}
								l1136:
									{
										position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1139
										}
										position++
										goto l1138
									l1139:
										position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
										if buffer[position] != rune('O') {
											goto l806
										}
										position++
									}
								l1138:
									{
										position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1141
										}
										position++
										goto l1140
									l1141:
										position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
										if buffer[position] != rune('R') {
											goto l806
										}
										position++
									}
								l1140:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleFLOOR, position1131)
								}
								break
							case 'R', 'r':
								{
									position1142 := position
									depth++
									{
										position1143, tokenIndex1143, depth1143 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1144
										}
										position++
										goto l1143
									l1144:
										position, tokenIndex, depth = position1143, tokenIndex1143, depth1143
										if buffer[position] != rune('R') {
											goto l806
										}
										position++
									}
								l1143:
									{
										position1145, tokenIndex1145, depth1145 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1146
										}
										position++
										goto l1145
									l1146:
										position, tokenIndex, depth = position1145, tokenIndex1145, depth1145
										if buffer[position] != rune('O') {
											goto l806
										}
										position++
									}
								l1145:
									{
										position1147, tokenIndex1147, depth1147 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1148
										}
										position++
										goto l1147
									l1148:
										position, tokenIndex, depth = position1147, tokenIndex1147, depth1147
										if buffer[position] != rune('U') {
											goto l806
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
											goto l806
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
											goto l806
										}
										position++
									}
								l1151:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleROUND, position1142)
								}
								break
							case 'C', 'c':
								{
									position1153 := position
									depth++
									{
										position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1155
										}
										position++
										goto l1154
									l1155:
										position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
										if buffer[position] != rune('C') {
											goto l806
										}
										position++
									}
								l1154:
									{
										position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1157
										}
										position++
										goto l1156
									l1157:
										position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
										if buffer[position] != rune('E') {
											goto l806
										}
										position++
									}
								l1156:
									{
										position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1159
										}
										position++
										goto l1158
									l1159:
										position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
										if buffer[position] != rune('I') {
											goto l806
										}
										position++
									}
								l1158:
									{
										position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1161
										}
										position++
										goto l1160
									l1161:
										position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
										if buffer[position] != rune('L') {
											goto l806
										}
										position++
									}
								l1160:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleCEIL, position1153)
								}
								break
							default:
								{
									position1162 := position
									depth++
									{
										position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1164
										}
										position++
										goto l1163
									l1164:
										position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
										if buffer[position] != rune('A') {
											goto l806
										}
										position++
									}
								l1163:
									{
										position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1166
										}
										position++
										goto l1165
									l1166:
										position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
										if buffer[position] != rune('B') {
											goto l806
										}
										position++
									}
								l1165:
									{
										position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1168
										}
										position++
										goto l1167
									l1168:
										position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
										if buffer[position] != rune('S') {
											goto l806
										}
										position++
									}
								l1167:
									if !rules[ruleskip]() {
										goto l806
									}
									depth--
									add(ruleABS, position1162)
								}
								break
							}
						}

					}
				l807:
					if !rules[ruleLPAREN]() {
						goto l806
					}
					if !rules[ruleexpression]() {
						goto l806
					}
					if !rules[ruleRPAREN]() {
						goto l806
					}
					goto l805
				l806:
					position, tokenIndex, depth = position805, tokenIndex805, depth805
					{
						position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
						{
							position1172 := position
							depth++
							{
								position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1174
								}
								position++
								goto l1173
							l1174:
								position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
								if buffer[position] != rune('S') {
									goto l1171
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
									goto l1171
								}
								position++
							}
						l1175:
							{
								position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1178
								}
								position++
								goto l1177
							l1178:
								position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
								if buffer[position] != rune('R') {
									goto l1171
								}
								position++
							}
						l1177:
							{
								position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1180
								}
								position++
								goto l1179
							l1180:
								position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
								if buffer[position] != rune('S') {
									goto l1171
								}
								position++
							}
						l1179:
							{
								position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1182
								}
								position++
								goto l1181
							l1182:
								position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
								if buffer[position] != rune('T') {
									goto l1171
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
									goto l1171
								}
								position++
							}
						l1183:
							{
								position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1186
								}
								position++
								goto l1185
							l1186:
								position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
								if buffer[position] != rune('R') {
									goto l1171
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
									goto l1171
								}
								position++
							}
						l1187:
							{
								position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1190
								}
								position++
								goto l1189
							l1190:
								position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
								if buffer[position] != rune('S') {
									goto l1171
								}
								position++
							}
						l1189:
							if !rules[ruleskip]() {
								goto l1171
							}
							depth--
							add(ruleSTRSTARTS, position1172)
						}
						goto l1170
					l1171:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						{
							position1192 := position
							depth++
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
									goto l1191
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
									goto l1191
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
									goto l1191
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
									goto l1191
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
									goto l1191
								}
								position++
							}
						l1201:
							{
								position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1204
								}
								position++
								goto l1203
							l1204:
								position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
								if buffer[position] != rune('D') {
									goto l1191
								}
								position++
							}
						l1203:
							{
								position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1206
								}
								position++
								goto l1205
							l1206:
								position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
								if buffer[position] != rune('S') {
									goto l1191
								}
								position++
							}
						l1205:
							if !rules[ruleskip]() {
								goto l1191
							}
							depth--
							add(ruleSTRENDS, position1192)
						}
						goto l1170
					l1191:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						{
							position1208 := position
							depth++
							{
								position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1210
								}
								position++
								goto l1209
							l1210:
								position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
								if buffer[position] != rune('S') {
									goto l1207
								}
								position++
							}
						l1209:
							{
								position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1212
								}
								position++
								goto l1211
							l1212:
								position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
								if buffer[position] != rune('T') {
									goto l1207
								}
								position++
							}
						l1211:
							{
								position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1214
								}
								position++
								goto l1213
							l1214:
								position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
								if buffer[position] != rune('R') {
									goto l1207
								}
								position++
							}
						l1213:
							{
								position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1216
								}
								position++
								goto l1215
							l1216:
								position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
								if buffer[position] != rune('B') {
									goto l1207
								}
								position++
							}
						l1215:
							{
								position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1218
								}
								position++
								goto l1217
							l1218:
								position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
								if buffer[position] != rune('E') {
									goto l1207
								}
								position++
							}
						l1217:
							{
								position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1220
								}
								position++
								goto l1219
							l1220:
								position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
								if buffer[position] != rune('F') {
									goto l1207
								}
								position++
							}
						l1219:
							{
								position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1222
								}
								position++
								goto l1221
							l1222:
								position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
								if buffer[position] != rune('O') {
									goto l1207
								}
								position++
							}
						l1221:
							{
								position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1224
								}
								position++
								goto l1223
							l1224:
								position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
								if buffer[position] != rune('R') {
									goto l1207
								}
								position++
							}
						l1223:
							{
								position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1226
								}
								position++
								goto l1225
							l1226:
								position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
								if buffer[position] != rune('E') {
									goto l1207
								}
								position++
							}
						l1225:
							if !rules[ruleskip]() {
								goto l1207
							}
							depth--
							add(ruleSTRBEFORE, position1208)
						}
						goto l1170
					l1207:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						{
							position1228 := position
							depth++
							{
								position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1230
								}
								position++
								goto l1229
							l1230:
								position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
								if buffer[position] != rune('S') {
									goto l1227
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
									goto l1227
								}
								position++
							}
						l1231:
							{
								position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1234
								}
								position++
								goto l1233
							l1234:
								position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
								if buffer[position] != rune('R') {
									goto l1227
								}
								position++
							}
						l1233:
							{
								position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1236
								}
								position++
								goto l1235
							l1236:
								position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
								if buffer[position] != rune('A') {
									goto l1227
								}
								position++
							}
						l1235:
							{
								position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1238
								}
								position++
								goto l1237
							l1238:
								position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
								if buffer[position] != rune('F') {
									goto l1227
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
									goto l1227
								}
								position++
							}
						l1239:
							{
								position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1242
								}
								position++
								goto l1241
							l1242:
								position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
								if buffer[position] != rune('E') {
									goto l1227
								}
								position++
							}
						l1241:
							{
								position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1244
								}
								position++
								goto l1243
							l1244:
								position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
								if buffer[position] != rune('R') {
									goto l1227
								}
								position++
							}
						l1243:
							if !rules[ruleskip]() {
								goto l1227
							}
							depth--
							add(ruleSTRAFTER, position1228)
						}
						goto l1170
					l1227:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						{
							position1246 := position
							depth++
							{
								position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1248
								}
								position++
								goto l1247
							l1248:
								position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
								if buffer[position] != rune('S') {
									goto l1245
								}
								position++
							}
						l1247:
							{
								position1249, tokenIndex1249, depth1249 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1250
								}
								position++
								goto l1249
							l1250:
								position, tokenIndex, depth = position1249, tokenIndex1249, depth1249
								if buffer[position] != rune('T') {
									goto l1245
								}
								position++
							}
						l1249:
							{
								position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1252
								}
								position++
								goto l1251
							l1252:
								position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
								if buffer[position] != rune('R') {
									goto l1245
								}
								position++
							}
						l1251:
							{
								position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1254
								}
								position++
								goto l1253
							l1254:
								position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
								if buffer[position] != rune('L') {
									goto l1245
								}
								position++
							}
						l1253:
							{
								position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1256
								}
								position++
								goto l1255
							l1256:
								position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
								if buffer[position] != rune('A') {
									goto l1245
								}
								position++
							}
						l1255:
							{
								position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1258
								}
								position++
								goto l1257
							l1258:
								position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
								if buffer[position] != rune('N') {
									goto l1245
								}
								position++
							}
						l1257:
							{
								position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1260
								}
								position++
								goto l1259
							l1260:
								position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
								if buffer[position] != rune('G') {
									goto l1245
								}
								position++
							}
						l1259:
							if !rules[ruleskip]() {
								goto l1245
							}
							depth--
							add(ruleSTRLANG, position1246)
						}
						goto l1170
					l1245:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						{
							position1262 := position
							depth++
							{
								position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1264
								}
								position++
								goto l1263
							l1264:
								position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
								if buffer[position] != rune('S') {
									goto l1261
								}
								position++
							}
						l1263:
							{
								position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1266
								}
								position++
								goto l1265
							l1266:
								position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
								if buffer[position] != rune('T') {
									goto l1261
								}
								position++
							}
						l1265:
							{
								position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1268
								}
								position++
								goto l1267
							l1268:
								position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
								if buffer[position] != rune('R') {
									goto l1261
								}
								position++
							}
						l1267:
							{
								position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1270
								}
								position++
								goto l1269
							l1270:
								position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
								if buffer[position] != rune('D') {
									goto l1261
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
									goto l1261
								}
								position++
							}
						l1271:
							if !rules[ruleskip]() {
								goto l1261
							}
							depth--
							add(ruleSTRDT, position1262)
						}
						goto l1170
					l1261:
						position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1274 := position
									depth++
									{
										position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1276
										}
										position++
										goto l1275
									l1276:
										position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
										if buffer[position] != rune('S') {
											goto l1169
										}
										position++
									}
								l1275:
									{
										position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1278
										}
										position++
										goto l1277
									l1278:
										position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
										if buffer[position] != rune('A') {
											goto l1169
										}
										position++
									}
								l1277:
									{
										position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1280
										}
										position++
										goto l1279
									l1280:
										position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
										if buffer[position] != rune('M') {
											goto l1169
										}
										position++
									}
								l1279:
									{
										position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1282
										}
										position++
										goto l1281
									l1282:
										position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
										if buffer[position] != rune('E') {
											goto l1169
										}
										position++
									}
								l1281:
									{
										position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1284
										}
										position++
										goto l1283
									l1284:
										position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
										if buffer[position] != rune('T') {
											goto l1169
										}
										position++
									}
								l1283:
									{
										position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1286
										}
										position++
										goto l1285
									l1286:
										position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
										if buffer[position] != rune('E') {
											goto l1169
										}
										position++
									}
								l1285:
									{
										position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1288
										}
										position++
										goto l1287
									l1288:
										position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
										if buffer[position] != rune('R') {
											goto l1169
										}
										position++
									}
								l1287:
									{
										position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1290
										}
										position++
										goto l1289
									l1290:
										position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
										if buffer[position] != rune('M') {
											goto l1169
										}
										position++
									}
								l1289:
									if !rules[ruleskip]() {
										goto l1169
									}
									depth--
									add(ruleSAMETERM, position1274)
								}
								break
							case 'C', 'c':
								{
									position1291 := position
									depth++
									{
										position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1293
										}
										position++
										goto l1292
									l1293:
										position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
										if buffer[position] != rune('C') {
											goto l1169
										}
										position++
									}
								l1292:
									{
										position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1295
										}
										position++
										goto l1294
									l1295:
										position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
										if buffer[position] != rune('O') {
											goto l1169
										}
										position++
									}
								l1294:
									{
										position1296, tokenIndex1296, depth1296 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1297
										}
										position++
										goto l1296
									l1297:
										position, tokenIndex, depth = position1296, tokenIndex1296, depth1296
										if buffer[position] != rune('N') {
											goto l1169
										}
										position++
									}
								l1296:
									{
										position1298, tokenIndex1298, depth1298 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1299
										}
										position++
										goto l1298
									l1299:
										position, tokenIndex, depth = position1298, tokenIndex1298, depth1298
										if buffer[position] != rune('T') {
											goto l1169
										}
										position++
									}
								l1298:
									{
										position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1301
										}
										position++
										goto l1300
									l1301:
										position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
										if buffer[position] != rune('A') {
											goto l1169
										}
										position++
									}
								l1300:
									{
										position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1303
										}
										position++
										goto l1302
									l1303:
										position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
										if buffer[position] != rune('I') {
											goto l1169
										}
										position++
									}
								l1302:
									{
										position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1305
										}
										position++
										goto l1304
									l1305:
										position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
										if buffer[position] != rune('N') {
											goto l1169
										}
										position++
									}
								l1304:
									{
										position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1307
										}
										position++
										goto l1306
									l1307:
										position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
										if buffer[position] != rune('S') {
											goto l1169
										}
										position++
									}
								l1306:
									if !rules[ruleskip]() {
										goto l1169
									}
									depth--
									add(ruleCONTAINS, position1291)
								}
								break
							default:
								{
									position1308 := position
									depth++
									{
										position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1310
										}
										position++
										goto l1309
									l1310:
										position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
										if buffer[position] != rune('L') {
											goto l1169
										}
										position++
									}
								l1309:
									{
										position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1312
										}
										position++
										goto l1311
									l1312:
										position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
										if buffer[position] != rune('A') {
											goto l1169
										}
										position++
									}
								l1311:
									{
										position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1314
										}
										position++
										goto l1313
									l1314:
										position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
										if buffer[position] != rune('N') {
											goto l1169
										}
										position++
									}
								l1313:
									{
										position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1316
										}
										position++
										goto l1315
									l1316:
										position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
										if buffer[position] != rune('G') {
											goto l1169
										}
										position++
									}
								l1315:
									{
										position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1318
										}
										position++
										goto l1317
									l1318:
										position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
										if buffer[position] != rune('M') {
											goto l1169
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
											goto l1169
										}
										position++
									}
								l1319:
									{
										position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1322
										}
										position++
										goto l1321
									l1322:
										position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
										if buffer[position] != rune('T') {
											goto l1169
										}
										position++
									}
								l1321:
									{
										position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1324
										}
										position++
										goto l1323
									l1324:
										position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
										if buffer[position] != rune('C') {
											goto l1169
										}
										position++
									}
								l1323:
									{
										position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1326
										}
										position++
										goto l1325
									l1326:
										position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
										if buffer[position] != rune('H') {
											goto l1169
										}
										position++
									}
								l1325:
									{
										position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1328
										}
										position++
										goto l1327
									l1328:
										position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
										if buffer[position] != rune('E') {
											goto l1169
										}
										position++
									}
								l1327:
									{
										position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1330
										}
										position++
										goto l1329
									l1330:
										position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
										if buffer[position] != rune('S') {
											goto l1169
										}
										position++
									}
								l1329:
									if !rules[ruleskip]() {
										goto l1169
									}
									depth--
									add(ruleLANGMATCHES, position1308)
								}
								break
							}
						}

					}
				l1170:
					if !rules[ruleLPAREN]() {
						goto l1169
					}
					if !rules[ruleexpression]() {
						goto l1169
					}
					if !rules[ruleCOMMA]() {
						goto l1169
					}
					if !rules[ruleexpression]() {
						goto l1169
					}
					if !rules[ruleRPAREN]() {
						goto l1169
					}
					goto l805
				l1169:
					position, tokenIndex, depth = position805, tokenIndex805, depth805
					{
						position1332 := position
						depth++
						{
							position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1334
							}
							position++
							goto l1333
						l1334:
							position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
							if buffer[position] != rune('B') {
								goto l1331
							}
							position++
						}
					l1333:
						{
							position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1336
							}
							position++
							goto l1335
						l1336:
							position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
							if buffer[position] != rune('O') {
								goto l1331
							}
							position++
						}
					l1335:
						{
							position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1338
							}
							position++
							goto l1337
						l1338:
							position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
							if buffer[position] != rune('U') {
								goto l1331
							}
							position++
						}
					l1337:
						{
							position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1340
							}
							position++
							goto l1339
						l1340:
							position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
							if buffer[position] != rune('N') {
								goto l1331
							}
							position++
						}
					l1339:
						{
							position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1342
							}
							position++
							goto l1341
						l1342:
							position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
							if buffer[position] != rune('D') {
								goto l1331
							}
							position++
						}
					l1341:
						if !rules[ruleskip]() {
							goto l1331
						}
						depth--
						add(ruleBOUND, position1332)
					}
					if !rules[ruleLPAREN]() {
						goto l1331
					}
					if !rules[rulevar]() {
						goto l1331
					}
					if !rules[ruleRPAREN]() {
						goto l1331
					}
					goto l805
				l1331:
					position, tokenIndex, depth = position805, tokenIndex805, depth805
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1345 := position
								depth++
								{
									position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1347
									}
									position++
									goto l1346
								l1347:
									position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
									if buffer[position] != rune('S') {
										goto l1343
									}
									position++
								}
							l1346:
								{
									position1348, tokenIndex1348, depth1348 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1349
									}
									position++
									goto l1348
								l1349:
									position, tokenIndex, depth = position1348, tokenIndex1348, depth1348
									if buffer[position] != rune('T') {
										goto l1343
									}
									position++
								}
							l1348:
								{
									position1350, tokenIndex1350, depth1350 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1351
									}
									position++
									goto l1350
								l1351:
									position, tokenIndex, depth = position1350, tokenIndex1350, depth1350
									if buffer[position] != rune('R') {
										goto l1343
									}
									position++
								}
							l1350:
								{
									position1352, tokenIndex1352, depth1352 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1353
									}
									position++
									goto l1352
								l1353:
									position, tokenIndex, depth = position1352, tokenIndex1352, depth1352
									if buffer[position] != rune('U') {
										goto l1343
									}
									position++
								}
							l1352:
								{
									position1354, tokenIndex1354, depth1354 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1355
									}
									position++
									goto l1354
								l1355:
									position, tokenIndex, depth = position1354, tokenIndex1354, depth1354
									if buffer[position] != rune('U') {
										goto l1343
									}
									position++
								}
							l1354:
								{
									position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1357
									}
									position++
									goto l1356
								l1357:
									position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
									if buffer[position] != rune('I') {
										goto l1343
									}
									position++
								}
							l1356:
								{
									position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1359
									}
									position++
									goto l1358
								l1359:
									position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
									if buffer[position] != rune('D') {
										goto l1343
									}
									position++
								}
							l1358:
								if !rules[ruleskip]() {
									goto l1343
								}
								depth--
								add(ruleSTRUUID, position1345)
							}
							break
						case 'U', 'u':
							{
								position1360 := position
								depth++
								{
									position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1362
									}
									position++
									goto l1361
								l1362:
									position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
									if buffer[position] != rune('U') {
										goto l1343
									}
									position++
								}
							l1361:
								{
									position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1364
									}
									position++
									goto l1363
								l1364:
									position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
									if buffer[position] != rune('U') {
										goto l1343
									}
									position++
								}
							l1363:
								{
									position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1366
									}
									position++
									goto l1365
								l1366:
									position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
									if buffer[position] != rune('I') {
										goto l1343
									}
									position++
								}
							l1365:
								{
									position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1368
									}
									position++
									goto l1367
								l1368:
									position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
									if buffer[position] != rune('D') {
										goto l1343
									}
									position++
								}
							l1367:
								if !rules[ruleskip]() {
									goto l1343
								}
								depth--
								add(ruleUUID, position1360)
							}
							break
						case 'N', 'n':
							{
								position1369 := position
								depth++
								{
									position1370, tokenIndex1370, depth1370 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1371
									}
									position++
									goto l1370
								l1371:
									position, tokenIndex, depth = position1370, tokenIndex1370, depth1370
									if buffer[position] != rune('N') {
										goto l1343
									}
									position++
								}
							l1370:
								{
									position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1373
									}
									position++
									goto l1372
								l1373:
									position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
									if buffer[position] != rune('O') {
										goto l1343
									}
									position++
								}
							l1372:
								{
									position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1375
									}
									position++
									goto l1374
								l1375:
									position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
									if buffer[position] != rune('W') {
										goto l1343
									}
									position++
								}
							l1374:
								if !rules[ruleskip]() {
									goto l1343
								}
								depth--
								add(ruleNOW, position1369)
							}
							break
						default:
							{
								position1376 := position
								depth++
								{
									position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1378
									}
									position++
									goto l1377
								l1378:
									position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
									if buffer[position] != rune('R') {
										goto l1343
									}
									position++
								}
							l1377:
								{
									position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1380
									}
									position++
									goto l1379
								l1380:
									position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
									if buffer[position] != rune('A') {
										goto l1343
									}
									position++
								}
							l1379:
								{
									position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1382
									}
									position++
									goto l1381
								l1382:
									position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
									if buffer[position] != rune('N') {
										goto l1343
									}
									position++
								}
							l1381:
								{
									position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1384
									}
									position++
									goto l1383
								l1384:
									position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
									if buffer[position] != rune('D') {
										goto l1343
									}
									position++
								}
							l1383:
								if !rules[ruleskip]() {
									goto l1343
								}
								depth--
								add(ruleRAND, position1376)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1343
					}
					goto l805
				l1343:
					position, tokenIndex, depth = position805, tokenIndex805, depth805
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
								{
									position1388 := position
									depth++
									{
										position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1390
										}
										position++
										goto l1389
									l1390:
										position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
										if buffer[position] != rune('E') {
											goto l1387
										}
										position++
									}
								l1389:
									{
										position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1392
										}
										position++
										goto l1391
									l1392:
										position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
										if buffer[position] != rune('X') {
											goto l1387
										}
										position++
									}
								l1391:
									{
										position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1394
										}
										position++
										goto l1393
									l1394:
										position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
										if buffer[position] != rune('I') {
											goto l1387
										}
										position++
									}
								l1393:
									{
										position1395, tokenIndex1395, depth1395 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1396
										}
										position++
										goto l1395
									l1396:
										position, tokenIndex, depth = position1395, tokenIndex1395, depth1395
										if buffer[position] != rune('S') {
											goto l1387
										}
										position++
									}
								l1395:
									{
										position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1398
										}
										position++
										goto l1397
									l1398:
										position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
										if buffer[position] != rune('T') {
											goto l1387
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
											goto l1387
										}
										position++
									}
								l1399:
									if !rules[ruleskip]() {
										goto l1387
									}
									depth--
									add(ruleEXISTS, position1388)
								}
								goto l1386
							l1387:
								position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
								{
									position1401 := position
									depth++
									{
										position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1403
										}
										position++
										goto l1402
									l1403:
										position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
										if buffer[position] != rune('N') {
											goto l803
										}
										position++
									}
								l1402:
									{
										position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1405
										}
										position++
										goto l1404
									l1405:
										position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
										if buffer[position] != rune('O') {
											goto l803
										}
										position++
									}
								l1404:
									{
										position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1407
										}
										position++
										goto l1406
									l1407:
										position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
										if buffer[position] != rune('T') {
											goto l803
										}
										position++
									}
								l1406:
									if buffer[position] != rune(' ') {
										goto l803
									}
									position++
									{
										position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1409
										}
										position++
										goto l1408
									l1409:
										position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
										if buffer[position] != rune('E') {
											goto l803
										}
										position++
									}
								l1408:
									{
										position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1411
										}
										position++
										goto l1410
									l1411:
										position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
										if buffer[position] != rune('X') {
											goto l803
										}
										position++
									}
								l1410:
									{
										position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1413
										}
										position++
										goto l1412
									l1413:
										position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
										if buffer[position] != rune('I') {
											goto l803
										}
										position++
									}
								l1412:
									{
										position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1415
										}
										position++
										goto l1414
									l1415:
										position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
										if buffer[position] != rune('S') {
											goto l803
										}
										position++
									}
								l1414:
									{
										position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1417
										}
										position++
										goto l1416
									l1417:
										position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
										if buffer[position] != rune('T') {
											goto l803
										}
										position++
									}
								l1416:
									{
										position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1419
										}
										position++
										goto l1418
									l1419:
										position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
										if buffer[position] != rune('S') {
											goto l803
										}
										position++
									}
								l1418:
									if !rules[ruleskip]() {
										goto l803
									}
									depth--
									add(ruleNOTEXIST, position1401)
								}
							}
						l1386:
							if !rules[rulegroupGraphPattern]() {
								goto l803
							}
							break
						case 'I', 'i':
							{
								position1420 := position
								depth++
								{
									position1421, tokenIndex1421, depth1421 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1422
									}
									position++
									goto l1421
								l1422:
									position, tokenIndex, depth = position1421, tokenIndex1421, depth1421
									if buffer[position] != rune('I') {
										goto l803
									}
									position++
								}
							l1421:
								{
									position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1424
									}
									position++
									goto l1423
								l1424:
									position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
									if buffer[position] != rune('F') {
										goto l803
									}
									position++
								}
							l1423:
								if !rules[ruleskip]() {
									goto l803
								}
								depth--
								add(ruleIF, position1420)
							}
							if !rules[ruleLPAREN]() {
								goto l803
							}
							if !rules[ruleexpression]() {
								goto l803
							}
							if !rules[ruleCOMMA]() {
								goto l803
							}
							if !rules[ruleexpression]() {
								goto l803
							}
							if !rules[ruleCOMMA]() {
								goto l803
							}
							if !rules[ruleexpression]() {
								goto l803
							}
							if !rules[ruleRPAREN]() {
								goto l803
							}
							break
						case 'C', 'c':
							{
								position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
								{
									position1427 := position
									depth++
									{
										position1428, tokenIndex1428, depth1428 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1429
										}
										position++
										goto l1428
									l1429:
										position, tokenIndex, depth = position1428, tokenIndex1428, depth1428
										if buffer[position] != rune('C') {
											goto l1426
										}
										position++
									}
								l1428:
									{
										position1430, tokenIndex1430, depth1430 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1431
										}
										position++
										goto l1430
									l1431:
										position, tokenIndex, depth = position1430, tokenIndex1430, depth1430
										if buffer[position] != rune('O') {
											goto l1426
										}
										position++
									}
								l1430:
									{
										position1432, tokenIndex1432, depth1432 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1433
										}
										position++
										goto l1432
									l1433:
										position, tokenIndex, depth = position1432, tokenIndex1432, depth1432
										if buffer[position] != rune('N') {
											goto l1426
										}
										position++
									}
								l1432:
									{
										position1434, tokenIndex1434, depth1434 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1435
										}
										position++
										goto l1434
									l1435:
										position, tokenIndex, depth = position1434, tokenIndex1434, depth1434
										if buffer[position] != rune('C') {
											goto l1426
										}
										position++
									}
								l1434:
									{
										position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1437
										}
										position++
										goto l1436
									l1437:
										position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
										if buffer[position] != rune('A') {
											goto l1426
										}
										position++
									}
								l1436:
									{
										position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1439
										}
										position++
										goto l1438
									l1439:
										position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
										if buffer[position] != rune('T') {
											goto l1426
										}
										position++
									}
								l1438:
									if !rules[ruleskip]() {
										goto l1426
									}
									depth--
									add(ruleCONCAT, position1427)
								}
								goto l1425
							l1426:
								position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
								{
									position1440 := position
									depth++
									{
										position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1442
										}
										position++
										goto l1441
									l1442:
										position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
										if buffer[position] != rune('C') {
											goto l803
										}
										position++
									}
								l1441:
									{
										position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1444
										}
										position++
										goto l1443
									l1444:
										position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
										if buffer[position] != rune('O') {
											goto l803
										}
										position++
									}
								l1443:
									{
										position1445, tokenIndex1445, depth1445 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1446
										}
										position++
										goto l1445
									l1446:
										position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
										if buffer[position] != rune('A') {
											goto l803
										}
										position++
									}
								l1445:
									{
										position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1448
										}
										position++
										goto l1447
									l1448:
										position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
										if buffer[position] != rune('L') {
											goto l803
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
											goto l803
										}
										position++
									}
								l1449:
									{
										position1451, tokenIndex1451, depth1451 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1452
										}
										position++
										goto l1451
									l1452:
										position, tokenIndex, depth = position1451, tokenIndex1451, depth1451
										if buffer[position] != rune('S') {
											goto l803
										}
										position++
									}
								l1451:
									{
										position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1454
										}
										position++
										goto l1453
									l1454:
										position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
										if buffer[position] != rune('C') {
											goto l803
										}
										position++
									}
								l1453:
									{
										position1455, tokenIndex1455, depth1455 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1456
										}
										position++
										goto l1455
									l1456:
										position, tokenIndex, depth = position1455, tokenIndex1455, depth1455
										if buffer[position] != rune('E') {
											goto l803
										}
										position++
									}
								l1455:
									if !rules[ruleskip]() {
										goto l803
									}
									depth--
									add(ruleCOALESCE, position1440)
								}
							}
						l1425:
							if !rules[ruleargList]() {
								goto l803
							}
							break
						case 'B', 'b':
							{
								position1457 := position
								depth++
								{
									position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1459
									}
									position++
									goto l1458
								l1459:
									position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
									if buffer[position] != rune('B') {
										goto l803
									}
									position++
								}
							l1458:
								{
									position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1461
									}
									position++
									goto l1460
								l1461:
									position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
									if buffer[position] != rune('N') {
										goto l803
									}
									position++
								}
							l1460:
								{
									position1462, tokenIndex1462, depth1462 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1463
									}
									position++
									goto l1462
								l1463:
									position, tokenIndex, depth = position1462, tokenIndex1462, depth1462
									if buffer[position] != rune('O') {
										goto l803
									}
									position++
								}
							l1462:
								{
									position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1465
									}
									position++
									goto l1464
								l1465:
									position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
									if buffer[position] != rune('D') {
										goto l803
									}
									position++
								}
							l1464:
								{
									position1466, tokenIndex1466, depth1466 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1467
									}
									position++
									goto l1466
								l1467:
									position, tokenIndex, depth = position1466, tokenIndex1466, depth1466
									if buffer[position] != rune('E') {
										goto l803
									}
									position++
								}
							l1466:
								if !rules[ruleskip]() {
									goto l803
								}
								depth--
								add(ruleBNODE, position1457)
							}
							{
								position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1469
								}
								if !rules[ruleexpression]() {
									goto l1469
								}
								if !rules[ruleRPAREN]() {
									goto l1469
								}
								goto l1468
							l1469:
								position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
								if !rules[rulenil]() {
									goto l803
								}
							}
						l1468:
							break
						default:
							{
								position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
								{
									position1472 := position
									depth++
									{
										position1473, tokenIndex1473, depth1473 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1474
										}
										position++
										goto l1473
									l1474:
										position, tokenIndex, depth = position1473, tokenIndex1473, depth1473
										if buffer[position] != rune('S') {
											goto l1471
										}
										position++
									}
								l1473:
									{
										position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1476
										}
										position++
										goto l1475
									l1476:
										position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
										if buffer[position] != rune('U') {
											goto l1471
										}
										position++
									}
								l1475:
									{
										position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1478
										}
										position++
										goto l1477
									l1478:
										position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
										if buffer[position] != rune('B') {
											goto l1471
										}
										position++
									}
								l1477:
									{
										position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1480
										}
										position++
										goto l1479
									l1480:
										position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
										if buffer[position] != rune('S') {
											goto l1471
										}
										position++
									}
								l1479:
									{
										position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1482
										}
										position++
										goto l1481
									l1482:
										position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
										if buffer[position] != rune('T') {
											goto l1471
										}
										position++
									}
								l1481:
									{
										position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1484
										}
										position++
										goto l1483
									l1484:
										position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
										if buffer[position] != rune('R') {
											goto l1471
										}
										position++
									}
								l1483:
									if !rules[ruleskip]() {
										goto l1471
									}
									depth--
									add(ruleSUBSTR, position1472)
								}
								goto l1470
							l1471:
								position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
								{
									position1486 := position
									depth++
									{
										position1487, tokenIndex1487, depth1487 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1488
										}
										position++
										goto l1487
									l1488:
										position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
										if buffer[position] != rune('R') {
											goto l1485
										}
										position++
									}
								l1487:
									{
										position1489, tokenIndex1489, depth1489 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1490
										}
										position++
										goto l1489
									l1490:
										position, tokenIndex, depth = position1489, tokenIndex1489, depth1489
										if buffer[position] != rune('E') {
											goto l1485
										}
										position++
									}
								l1489:
									{
										position1491, tokenIndex1491, depth1491 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1492
										}
										position++
										goto l1491
									l1492:
										position, tokenIndex, depth = position1491, tokenIndex1491, depth1491
										if buffer[position] != rune('P') {
											goto l1485
										}
										position++
									}
								l1491:
									{
										position1493, tokenIndex1493, depth1493 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1494
										}
										position++
										goto l1493
									l1494:
										position, tokenIndex, depth = position1493, tokenIndex1493, depth1493
										if buffer[position] != rune('L') {
											goto l1485
										}
										position++
									}
								l1493:
									{
										position1495, tokenIndex1495, depth1495 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1496
										}
										position++
										goto l1495
									l1496:
										position, tokenIndex, depth = position1495, tokenIndex1495, depth1495
										if buffer[position] != rune('A') {
											goto l1485
										}
										position++
									}
								l1495:
									{
										position1497, tokenIndex1497, depth1497 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1498
										}
										position++
										goto l1497
									l1498:
										position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
										if buffer[position] != rune('C') {
											goto l1485
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
											goto l1485
										}
										position++
									}
								l1499:
									if !rules[ruleskip]() {
										goto l1485
									}
									depth--
									add(ruleREPLACE, position1486)
								}
								goto l1470
							l1485:
								position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
								{
									position1501 := position
									depth++
									{
										position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1503
										}
										position++
										goto l1502
									l1503:
										position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
										if buffer[position] != rune('R') {
											goto l803
										}
										position++
									}
								l1502:
									{
										position1504, tokenIndex1504, depth1504 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1505
										}
										position++
										goto l1504
									l1505:
										position, tokenIndex, depth = position1504, tokenIndex1504, depth1504
										if buffer[position] != rune('E') {
											goto l803
										}
										position++
									}
								l1504:
									{
										position1506, tokenIndex1506, depth1506 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1507
										}
										position++
										goto l1506
									l1507:
										position, tokenIndex, depth = position1506, tokenIndex1506, depth1506
										if buffer[position] != rune('G') {
											goto l803
										}
										position++
									}
								l1506:
									{
										position1508, tokenIndex1508, depth1508 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1509
										}
										position++
										goto l1508
									l1509:
										position, tokenIndex, depth = position1508, tokenIndex1508, depth1508
										if buffer[position] != rune('E') {
											goto l803
										}
										position++
									}
								l1508:
									{
										position1510, tokenIndex1510, depth1510 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1511
										}
										position++
										goto l1510
									l1511:
										position, tokenIndex, depth = position1510, tokenIndex1510, depth1510
										if buffer[position] != rune('X') {
											goto l803
										}
										position++
									}
								l1510:
									if !rules[ruleskip]() {
										goto l803
									}
									depth--
									add(ruleREGEX, position1501)
								}
							}
						l1470:
							if !rules[ruleLPAREN]() {
								goto l803
							}
							if !rules[ruleexpression]() {
								goto l803
							}
							if !rules[ruleCOMMA]() {
								goto l803
							}
							if !rules[ruleexpression]() {
								goto l803
							}
							{
								position1512, tokenIndex1512, depth1512 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1512
								}
								if !rules[ruleexpression]() {
									goto l1512
								}
								goto l1513
							l1512:
								position, tokenIndex, depth = position1512, tokenIndex1512, depth1512
							}
						l1513:
							if !rules[ruleRPAREN]() {
								goto l803
							}
							break
						}
					}

				}
			l805:
				depth--
				add(rulebuiltinCall, position804)
			}
			return true
		l803:
			position, tokenIndex, depth = position803, tokenIndex803, depth803
			return false
		},
		/* 70 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1514, tokenIndex1514, depth1514 := position, tokenIndex, depth
			{
				position1515 := position
				depth++
				{
					position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1517
					}
					position++
					goto l1516
				l1517:
					position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
					if buffer[position] != rune('$') {
						goto l1514
					}
					position++
				}
			l1516:
				{
					position1518 := position
					depth++
					{
						position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
						if !rules[rulepnCharsU]() {
							goto l1520
						}
						goto l1519
					l1520:
						position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1514
						}
						position++
					}
				l1519:
				l1521:
					{
						position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
						{
							position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
							if !rules[rulepnCharsU]() {
								goto l1524
							}
							goto l1523
						l1524:
							position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
							{
								switch buffer[position] {
								case 'â':
									if c := buffer[position]; c < rune('‿') || c > rune('⁀') {
										goto l1522
									}
									position++
									break
								case 'Ì', 'Í':
									if c := buffer[position]; c < rune('̀') || c > rune('ͯ') {
										goto l1522
									}
									position++
									break
								case 'Â':
									if buffer[position] != rune('·') {
										goto l1522
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1522
									}
									position++
									break
								}
							}

						}
					l1523:
						goto l1521
					l1522:
						position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
					}
					depth--
					add(ruleVARNAME, position1518)
				}
				if !rules[ruleskip]() {
					goto l1514
				}
				depth--
				add(rulevar, position1515)
			}
			return true
		l1514:
			position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
			return false
		},
		/* 71 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
			{
				position1527 := position
				depth++
				{
					position1528, tokenIndex1528, depth1528 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1529
					}
					goto l1528
				l1529:
					position, tokenIndex, depth = position1528, tokenIndex1528, depth1528
					{
						position1530 := position
						depth++
						{
							position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
							if !rules[rulepnPrefix]() {
								goto l1531
							}
							goto l1532
						l1531:
							position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
						}
					l1532:
						if buffer[position] != rune(':') {
							goto l1526
						}
						position++
						{
							position1533 := position
							depth++
							{
								switch buffer[position] {
								case '%', '\\':
									{
										position1537 := position
										depth++
										{
											position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
											{
												position1540 := position
												depth++
												if buffer[position] != rune('%') {
													goto l1539
												}
												position++
												if !rules[rulehex]() {
													goto l1539
												}
												if !rules[rulehex]() {
													goto l1539
												}
												depth--
												add(rulepercent, position1540)
											}
											goto l1538
										l1539:
											position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
											{
												position1541 := position
												depth++
												if buffer[position] != rune('\\') {
													goto l1526
												}
												position++
												{
													switch buffer[position] {
													case '%':
														if buffer[position] != rune('%') {
															goto l1526
														}
														position++
														break
													case '@':
														if buffer[position] != rune('@') {
															goto l1526
														}
														position++
														break
													case '#':
														if buffer[position] != rune('#') {
															goto l1526
														}
														position++
														break
													case '?':
														if buffer[position] != rune('?') {
															goto l1526
														}
														position++
														break
													case '/':
														if buffer[position] != rune('/') {
															goto l1526
														}
														position++
														break
													case '=':
														if buffer[position] != rune('=') {
															goto l1526
														}
														position++
														break
													case ';':
														if buffer[position] != rune(';') {
															goto l1526
														}
														position++
														break
													case ',':
														if buffer[position] != rune(',') {
															goto l1526
														}
														position++
														break
													case '+':
														if buffer[position] != rune('+') {
															goto l1526
														}
														position++
														break
													case '*':
														if buffer[position] != rune('*') {
															goto l1526
														}
														position++
														break
													case ')':
														if buffer[position] != rune(')') {
															goto l1526
														}
														position++
														break
													case '(':
														if buffer[position] != rune('(') {
															goto l1526
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l1526
														}
														position++
														break
													case '&':
														if buffer[position] != rune('&') {
															goto l1526
														}
														position++
														break
													case '$':
														if buffer[position] != rune('$') {
															goto l1526
														}
														position++
														break
													case '!':
														if buffer[position] != rune('!') {
															goto l1526
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l1526
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l1526
														}
														position++
														break
													case '~':
														if buffer[position] != rune('~') {
															goto l1526
														}
														position++
														break
													default:
														if buffer[position] != rune('_') {
															goto l1526
														}
														position++
														break
													}
												}

												depth--
												add(rulepnLocalEsc, position1541)
											}
										}
									l1538:
										depth--
										add(ruleplx, position1537)
									}
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1526
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l1526
									}
									position++
									break
								default:
									if !rules[rulepnCharsU]() {
										goto l1526
									}
									break
								}
							}

						l1534:
							{
								position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '%', '\\':
										{
											position1544 := position
											depth++
											{
												position1545, tokenIndex1545, depth1545 := position, tokenIndex, depth
												{
													position1547 := position
													depth++
													if buffer[position] != rune('%') {
														goto l1546
													}
													position++
													if !rules[rulehex]() {
														goto l1546
													}
													if !rules[rulehex]() {
														goto l1546
													}
													depth--
													add(rulepercent, position1547)
												}
												goto l1545
											l1546:
												position, tokenIndex, depth = position1545, tokenIndex1545, depth1545
												{
													position1548 := position
													depth++
													if buffer[position] != rune('\\') {
														goto l1535
													}
													position++
													{
														switch buffer[position] {
														case '%':
															if buffer[position] != rune('%') {
																goto l1535
															}
															position++
															break
														case '@':
															if buffer[position] != rune('@') {
																goto l1535
															}
															position++
															break
														case '#':
															if buffer[position] != rune('#') {
																goto l1535
															}
															position++
															break
														case '?':
															if buffer[position] != rune('?') {
																goto l1535
															}
															position++
															break
														case '/':
															if buffer[position] != rune('/') {
																goto l1535
															}
															position++
															break
														case '=':
															if buffer[position] != rune('=') {
																goto l1535
															}
															position++
															break
														case ';':
															if buffer[position] != rune(';') {
																goto l1535
															}
															position++
															break
														case ',':
															if buffer[position] != rune(',') {
																goto l1535
															}
															position++
															break
														case '+':
															if buffer[position] != rune('+') {
																goto l1535
															}
															position++
															break
														case '*':
															if buffer[position] != rune('*') {
																goto l1535
															}
															position++
															break
														case ')':
															if buffer[position] != rune(')') {
																goto l1535
															}
															position++
															break
														case '(':
															if buffer[position] != rune('(') {
																goto l1535
															}
															position++
															break
														case '\'':
															if buffer[position] != rune('\'') {
																goto l1535
															}
															position++
															break
														case '&':
															if buffer[position] != rune('&') {
																goto l1535
															}
															position++
															break
														case '$':
															if buffer[position] != rune('$') {
																goto l1535
															}
															position++
															break
														case '!':
															if buffer[position] != rune('!') {
																goto l1535
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l1535
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l1535
															}
															position++
															break
														case '~':
															if buffer[position] != rune('~') {
																goto l1535
															}
															position++
															break
														default:
															if buffer[position] != rune('_') {
																goto l1535
															}
															position++
															break
														}
													}

													depth--
													add(rulepnLocalEsc, position1548)
												}
											}
										l1545:
											depth--
											add(ruleplx, position1544)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1535
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1535
										}
										position++
										break
									default:
										if !rules[rulepnCharsU]() {
											goto l1535
										}
										break
									}
								}

								goto l1534
							l1535:
								position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
							}
							depth--
							add(rulepnLocal, position1533)
						}
						if !rules[ruleskip]() {
							goto l1526
						}
						depth--
						add(ruleprefixedName, position1530)
					}
				}
			l1528:
				depth--
				add(ruleiriref, position1527)
			}
			return true
		l1526:
			position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
			return false
		},
		/* 72 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1550, tokenIndex1550, depth1550 := position, tokenIndex, depth
			{
				position1551 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1550
				}
				position++
			l1552:
				{
					position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
					{
						position1554, tokenIndex1554, depth1554 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1554
						}
						position++
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
				if buffer[position] != rune('>') {
					goto l1550
				}
				position++
				if !rules[ruleskip]() {
					goto l1550
				}
				depth--
				add(ruleiri, position1551)
			}
			return true
		l1550:
			position, tokenIndex, depth = position1550, tokenIndex1550, depth1550
			return false
		},
		/* 73 prefixedName <- <(pnPrefix? ':' pnLocal skip)> */
		nil,
		/* 74 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
			{
				position1557 := position
				depth++
				if !rules[rulestring]() {
					goto l1556
				}
				{
					position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
					{
						position1560, tokenIndex1560, depth1560 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1561
						}
						position++
						{
							position1564, tokenIndex1564, depth1564 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1565
							}
							position++
							goto l1564
						l1565:
							position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1561
							}
							position++
						}
					l1564:
					l1562:
						{
							position1563, tokenIndex1563, depth1563 := position, tokenIndex, depth
							{
								position1566, tokenIndex1566, depth1566 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1567
								}
								position++
								goto l1566
							l1567:
								position, tokenIndex, depth = position1566, tokenIndex1566, depth1566
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1563
								}
								position++
							}
						l1566:
							goto l1562
						l1563:
							position, tokenIndex, depth = position1563, tokenIndex1563, depth1563
						}
					l1568:
						{
							position1569, tokenIndex1569, depth1569 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1569
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1569
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1569
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1569
									}
									position++
									break
								}
							}

						l1570:
							{
								position1571, tokenIndex1571, depth1571 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1571
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1571
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1571
										}
										position++
										break
									}
								}

								goto l1570
							l1571:
								position, tokenIndex, depth = position1571, tokenIndex1571, depth1571
							}
							goto l1568
						l1569:
							position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
						}
						goto l1560
					l1561:
						position, tokenIndex, depth = position1560, tokenIndex1560, depth1560
						if buffer[position] != rune('^') {
							goto l1558
						}
						position++
						if buffer[position] != rune('^') {
							goto l1558
						}
						position++
						if !rules[ruleiriref]() {
							goto l1558
						}
					}
				l1560:
					goto l1559
				l1558:
					position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
				}
			l1559:
				if !rules[ruleskip]() {
					goto l1556
				}
				depth--
				add(ruleliteral, position1557)
			}
			return true
		l1556:
			position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
			return false
		},
		/* 75 string <- <(stringLiteralA / stringLiteralB / stringLiteralLongA / stringLiteralLongB)> */
		func() bool {
			position1574, tokenIndex1574, depth1574 := position, tokenIndex, depth
			{
				position1575 := position
				depth++
				{
					position1576, tokenIndex1576, depth1576 := position, tokenIndex, depth
					{
						position1578 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1577
						}
						position++
					l1579:
						{
							position1580, tokenIndex1580, depth1580 := position, tokenIndex, depth
							{
								position1581, tokenIndex1581, depth1581 := position, tokenIndex, depth
								{
									position1583, tokenIndex1583, depth1583 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1583
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1583
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1583
											}
											position++
											break
										default:
											if buffer[position] != rune('\'') {
												goto l1583
											}
											position++
											break
										}
									}

									goto l1582
								l1583:
									position, tokenIndex, depth = position1583, tokenIndex1583, depth1583
								}
								if !matchDot() {
									goto l1582
								}
								goto l1581
							l1582:
								position, tokenIndex, depth = position1581, tokenIndex1581, depth1581
								if !rules[ruleechar]() {
									goto l1580
								}
							}
						l1581:
							goto l1579
						l1580:
							position, tokenIndex, depth = position1580, tokenIndex1580, depth1580
						}
						if buffer[position] != rune('\'') {
							goto l1577
						}
						position++
						depth--
						add(rulestringLiteralA, position1578)
					}
					goto l1576
				l1577:
					position, tokenIndex, depth = position1576, tokenIndex1576, depth1576
					{
						position1586 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1585
						}
						position++
					l1587:
						{
							position1588, tokenIndex1588, depth1588 := position, tokenIndex, depth
							{
								position1589, tokenIndex1589, depth1589 := position, tokenIndex, depth
								{
									position1591, tokenIndex1591, depth1591 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1591
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1591
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1591
											}
											position++
											break
										default:
											if buffer[position] != rune('"') {
												goto l1591
											}
											position++
											break
										}
									}

									goto l1590
								l1591:
									position, tokenIndex, depth = position1591, tokenIndex1591, depth1591
								}
								if !matchDot() {
									goto l1590
								}
								goto l1589
							l1590:
								position, tokenIndex, depth = position1589, tokenIndex1589, depth1589
								if !rules[ruleechar]() {
									goto l1588
								}
							}
						l1589:
							goto l1587
						l1588:
							position, tokenIndex, depth = position1588, tokenIndex1588, depth1588
						}
						if buffer[position] != rune('"') {
							goto l1585
						}
						position++
						depth--
						add(rulestringLiteralB, position1586)
					}
					goto l1576
				l1585:
					position, tokenIndex, depth = position1576, tokenIndex1576, depth1576
					{
						position1594 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1593
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1593
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1593
						}
						position++
					l1595:
						{
							position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
							{
								position1597, tokenIndex1597, depth1597 := position, tokenIndex, depth
								{
									position1599, tokenIndex1599, depth1599 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l1600
									}
									position++
									goto l1599
								l1600:
									position, tokenIndex, depth = position1599, tokenIndex1599, depth1599
									if buffer[position] != rune('\'') {
										goto l1597
									}
									position++
									if buffer[position] != rune('\'') {
										goto l1597
									}
									position++
								}
							l1599:
								goto l1598
							l1597:
								position, tokenIndex, depth = position1597, tokenIndex1597, depth1597
							}
						l1598:
							{
								position1601, tokenIndex1601, depth1601 := position, tokenIndex, depth
								{
									position1603, tokenIndex1603, depth1603 := position, tokenIndex, depth
									{
										position1604, tokenIndex1604, depth1604 := position, tokenIndex, depth
										if buffer[position] != rune('\'') {
											goto l1605
										}
										position++
										goto l1604
									l1605:
										position, tokenIndex, depth = position1604, tokenIndex1604, depth1604
										if buffer[position] != rune('\\') {
											goto l1603
										}
										position++
									}
								l1604:
									goto l1602
								l1603:
									position, tokenIndex, depth = position1603, tokenIndex1603, depth1603
								}
								if !matchDot() {
									goto l1602
								}
								goto l1601
							l1602:
								position, tokenIndex, depth = position1601, tokenIndex1601, depth1601
								if !rules[ruleechar]() {
									goto l1596
								}
							}
						l1601:
							goto l1595
						l1596:
							position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
						}
						if buffer[position] != rune('\'') {
							goto l1593
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1593
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1593
						}
						position++
						depth--
						add(rulestringLiteralLongA, position1594)
					}
					goto l1576
				l1593:
					position, tokenIndex, depth = position1576, tokenIndex1576, depth1576
					{
						position1606 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1574
						}
						position++
						if buffer[position] != rune('"') {
							goto l1574
						}
						position++
						if buffer[position] != rune('"') {
							goto l1574
						}
						position++
					l1607:
						{
							position1608, tokenIndex1608, depth1608 := position, tokenIndex, depth
							{
								position1609, tokenIndex1609, depth1609 := position, tokenIndex, depth
								{
									position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l1612
									}
									position++
									goto l1611
								l1612:
									position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
									if buffer[position] != rune('"') {
										goto l1609
									}
									position++
									if buffer[position] != rune('"') {
										goto l1609
									}
									position++
								}
							l1611:
								goto l1610
							l1609:
								position, tokenIndex, depth = position1609, tokenIndex1609, depth1609
							}
						l1610:
							{
								position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
								{
									position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
									{
										position1616, tokenIndex1616, depth1616 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											goto l1617
										}
										position++
										goto l1616
									l1617:
										position, tokenIndex, depth = position1616, tokenIndex1616, depth1616
										if buffer[position] != rune('\\') {
											goto l1615
										}
										position++
									}
								l1616:
									goto l1614
								l1615:
									position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
								}
								if !matchDot() {
									goto l1614
								}
								goto l1613
							l1614:
								position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
								if !rules[ruleechar]() {
									goto l1608
								}
							}
						l1613:
							goto l1607
						l1608:
							position, tokenIndex, depth = position1608, tokenIndex1608, depth1608
						}
						if buffer[position] != rune('"') {
							goto l1574
						}
						position++
						if buffer[position] != rune('"') {
							goto l1574
						}
						position++
						if buffer[position] != rune('"') {
							goto l1574
						}
						position++
						depth--
						add(rulestringLiteralLongB, position1606)
					}
				}
			l1576:
				depth--
				add(rulestring, position1575)
			}
			return true
		l1574:
			position, tokenIndex, depth = position1574, tokenIndex1574, depth1574
			return false
		},
		/* 76 stringLiteralA <- <('\'' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('\'') '\'')) .) / echar)* '\'')> */
		nil,
		/* 77 stringLiteralB <- <('"' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .) / echar)* '"')> */
		nil,
		/* 78 stringLiteralLongA <- <('\'' '\'' '\'' (('\'' / ('\'' '\''))? ((!('\'' / '\\') .) / echar))* ('\'' '\'' '\''))> */
		nil,
		/* 79 stringLiteralLongB <- <('"' '"' '"' (('"' / ('"' '"'))? ((!('"' / '\\') .) / echar))* ('"' '"' '"'))> */
		nil,
		/* 80 echar <- <('\\' ((&('\'') '\'') | (&('"') '"') | (&('\\') '\\') | (&('f') 'f') | (&('r') 'r') | (&('n') 'n') | (&('b') 'b') | (&('t') 't')))> */
		func() bool {
			position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
			{
				position1623 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l1622
				}
				position++
				{
					switch buffer[position] {
					case '\'':
						if buffer[position] != rune('\'') {
							goto l1622
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l1622
						}
						position++
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l1622
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l1622
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l1622
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l1622
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l1622
						}
						position++
						break
					default:
						if buffer[position] != rune('t') {
							goto l1622
						}
						position++
						break
					}
				}

				depth--
				add(ruleechar, position1623)
			}
			return true
		l1622:
			position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
			return false
		},
		/* 81 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1625, tokenIndex1625, depth1625 := position, tokenIndex, depth
			{
				position1626 := position
				depth++
				{
					position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
					{
						position1629, tokenIndex1629, depth1629 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1630
						}
						position++
						goto l1629
					l1630:
						position, tokenIndex, depth = position1629, tokenIndex1629, depth1629
						if buffer[position] != rune('-') {
							goto l1627
						}
						position++
					}
				l1629:
					goto l1628
				l1627:
					position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
				}
			l1628:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1625
				}
				position++
			l1631:
				{
					position1632, tokenIndex1632, depth1632 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1632
					}
					position++
					goto l1631
				l1632:
					position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
				}
				{
					position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1633
					}
					position++
				l1635:
					{
						position1636, tokenIndex1636, depth1636 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1636
						}
						position++
						goto l1635
					l1636:
						position, tokenIndex, depth = position1636, tokenIndex1636, depth1636
					}
					goto l1634
				l1633:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
				}
			l1634:
				if !rules[ruleskip]() {
					goto l1625
				}
				depth--
				add(rulenumericLiteral, position1626)
			}
			return true
		l1625:
			position, tokenIndex, depth = position1625, tokenIndex1625, depth1625
			return false
		},
		/* 82 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 83 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1638, tokenIndex1638, depth1638 := position, tokenIndex, depth
			{
				position1639 := position
				depth++
				{
					position1640, tokenIndex1640, depth1640 := position, tokenIndex, depth
					{
						position1642 := position
						depth++
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
								goto l1641
							}
							position++
						}
					l1643:
						{
							position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1646
							}
							position++
							goto l1645
						l1646:
							position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
							if buffer[position] != rune('R') {
								goto l1641
							}
							position++
						}
					l1645:
						{
							position1647, tokenIndex1647, depth1647 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1648
							}
							position++
							goto l1647
						l1648:
							position, tokenIndex, depth = position1647, tokenIndex1647, depth1647
							if buffer[position] != rune('U') {
								goto l1641
							}
							position++
						}
					l1647:
						{
							position1649, tokenIndex1649, depth1649 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1650
							}
							position++
							goto l1649
						l1650:
							position, tokenIndex, depth = position1649, tokenIndex1649, depth1649
							if buffer[position] != rune('E') {
								goto l1641
							}
							position++
						}
					l1649:
						if !rules[ruleskip]() {
							goto l1641
						}
						depth--
						add(ruleTRUE, position1642)
					}
					goto l1640
				l1641:
					position, tokenIndex, depth = position1640, tokenIndex1640, depth1640
					{
						position1651 := position
						depth++
						{
							position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1653
							}
							position++
							goto l1652
						l1653:
							position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
							if buffer[position] != rune('F') {
								goto l1638
							}
							position++
						}
					l1652:
						{
							position1654, tokenIndex1654, depth1654 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1655
							}
							position++
							goto l1654
						l1655:
							position, tokenIndex, depth = position1654, tokenIndex1654, depth1654
							if buffer[position] != rune('A') {
								goto l1638
							}
							position++
						}
					l1654:
						{
							position1656, tokenIndex1656, depth1656 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1657
							}
							position++
							goto l1656
						l1657:
							position, tokenIndex, depth = position1656, tokenIndex1656, depth1656
							if buffer[position] != rune('L') {
								goto l1638
							}
							position++
						}
					l1656:
						{
							position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1659
							}
							position++
							goto l1658
						l1659:
							position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
							if buffer[position] != rune('S') {
								goto l1638
							}
							position++
						}
					l1658:
						{
							position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1661
							}
							position++
							goto l1660
						l1661:
							position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
							if buffer[position] != rune('E') {
								goto l1638
							}
							position++
						}
					l1660:
						if !rules[ruleskip]() {
							goto l1638
						}
						depth--
						add(ruleFALSE, position1651)
					}
				}
			l1640:
				depth--
				add(rulebooleanLiteral, position1639)
			}
			return true
		l1638:
			position, tokenIndex, depth = position1638, tokenIndex1638, depth1638
			return false
		},
		/* 84 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 85 blankNodeLabel <- <('_' ':' (pnCharsU / [0-9]) (((pnCharsU / ((&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))* pnCharsU) / ([0-9] / '-'))? skip)> */
		nil,
		/* 86 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 87 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1665, tokenIndex1665, depth1665 := position, tokenIndex, depth
			{
				position1666 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1665
				}
				position++
			l1667:
				{
					position1668, tokenIndex1668, depth1668 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1668
					}
					goto l1667
				l1668:
					position, tokenIndex, depth = position1668, tokenIndex1668, depth1668
				}
				if buffer[position] != rune(')') {
					goto l1665
				}
				position++
				if !rules[ruleskip]() {
					goto l1665
				}
				depth--
				add(rulenil, position1666)
			}
			return true
		l1665:
			position, tokenIndex, depth = position1665, tokenIndex1665, depth1665
			return false
		},
		/* 88 VARNAME <- <((pnCharsU / [0-9]) (pnCharsU / ((&('â') [‿-⁀]) | (&('Ì' | 'Í') [̀-ͯ]) | (&('Â') '·') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))*)> */
		nil,
		/* 89 pnPrefix <- <(pnCharsBase pnChars*)> */
		func() bool {
			position1670, tokenIndex1670, depth1670 := position, tokenIndex, depth
			{
				position1671 := position
				depth++
				if !rules[rulepnCharsBase]() {
					goto l1670
				}
			l1672:
				{
					position1673, tokenIndex1673, depth1673 := position, tokenIndex, depth
					{
						position1674 := position
						depth++
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1673
								}
								position++
								break
							case '-':
								if buffer[position] != rune('-') {
									goto l1673
								}
								position++
								break
							default:
								if !rules[rulepnCharsU]() {
									goto l1673
								}
								break
							}
						}

						depth--
						add(rulepnChars, position1674)
					}
					goto l1672
				l1673:
					position, tokenIndex, depth = position1673, tokenIndex1673, depth1673
				}
				depth--
				add(rulepnPrefix, position1671)
			}
			return true
		l1670:
			position, tokenIndex, depth = position1670, tokenIndex1670, depth1670
			return false
		},
		/* 90 pnLocal <- <((&('%' | '\\') plx) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&(':') ':') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))+> */
		nil,
		/* 91 pnChars <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))> */
		nil,
		/* 92 pnCharsU <- <(pnCharsBase / '_')> */
		func() bool {
			position1678, tokenIndex1678, depth1678 := position, tokenIndex, depth
			{
				position1679 := position
				depth++
				{
					position1680, tokenIndex1680, depth1680 := position, tokenIndex, depth
					if !rules[rulepnCharsBase]() {
						goto l1681
					}
					goto l1680
				l1681:
					position, tokenIndex, depth = position1680, tokenIndex1680, depth1680
					if buffer[position] != rune('_') {
						goto l1678
					}
					position++
				}
			l1680:
				depth--
				add(rulepnCharsU, position1679)
			}
			return true
		l1678:
			position, tokenIndex, depth = position1678, tokenIndex1678, depth1678
			return false
		},
		/* 93 pnCharsBase <- <([À-Ö] / [Ø-ö] / [Ͱ-ͽ] / [‌-‍] / [⁰-↏] / [豈-﷏] / ((&('ð' | 'ñ' | 'ò' | 'ó') [𐀀-󯿿]) | (&('ï') [ﷰ-�]) | (&('ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í') [、-퟿]) | (&('â') [Ⰰ-⿯]) | (&('Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á') [Ϳ-῿]) | (&('Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë') [ø-˿]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))> */
		func() bool {
			position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
			{
				position1683 := position
				depth++
				{
					position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('À') || c > rune('Ö') {
						goto l1685
					}
					position++
					goto l1684
				l1685:
					position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
					if c := buffer[position]; c < rune('Ø') || c > rune('ö') {
						goto l1686
					}
					position++
					goto l1684
				l1686:
					position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
					if c := buffer[position]; c < rune('Ͱ') || c > rune('ͽ') {
						goto l1687
					}
					position++
					goto l1684
				l1687:
					position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
					if c := buffer[position]; c < rune('\u200c') || c > rune('\u200d') {
						goto l1688
					}
					position++
					goto l1684
				l1688:
					position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
					if c := buffer[position]; c < rune('⁰') || c > rune('\u218f') {
						goto l1689
					}
					position++
					goto l1684
				l1689:
					position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
					if c := buffer[position]; c < rune('豈') || c > rune('\ufdcf') {
						goto l1690
					}
					position++
					goto l1684
				l1690:
					position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
					{
						switch buffer[position] {
						case 'ð', 'ñ', 'ò', 'ó':
							if c := buffer[position]; c < rune('𐀀') || c > rune('\U000effff') {
								goto l1682
							}
							position++
							break
						case 'ï':
							if c := buffer[position]; c < rune('ﷰ') || c > rune('�') {
								goto l1682
							}
							position++
							break
						case 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í':
							if c := buffer[position]; c < rune('、') || c > rune('\ud7ff') {
								goto l1682
							}
							position++
							break
						case 'â':
							if c := buffer[position]; c < rune('Ⰰ') || c > rune('\u2fef') {
								goto l1682
							}
							position++
							break
						case 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'Þ', 'ß', 'à', 'á':
							if c := buffer[position]; c < rune('\u037f') || c > rune('\u1fff') {
								goto l1682
							}
							position++
							break
						case 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ê', 'Ë':
							if c := buffer[position]; c < rune('ø') || c > rune('˿') {
								goto l1682
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1682
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1682
							}
							position++
							break
						}
					}

				}
			l1684:
				depth--
				add(rulepnCharsBase, position1683)
			}
			return true
		l1682:
			position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
			return false
		},
		/* 94 plx <- <(percent / pnLocalEsc)> */
		nil,
		/* 95 percent <- <('%' hex hex)> */
		nil,
		/* 96 hex <- <((&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))> */
		func() bool {
			position1694, tokenIndex1694, depth1694 := position, tokenIndex, depth
			{
				position1695 := position
				depth++
				{
					switch buffer[position] {
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l1694
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1694
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1694
						}
						position++
						break
					}
				}

				depth--
				add(rulehex, position1695)
			}
			return true
		l1694:
			position, tokenIndex, depth = position1694, tokenIndex1694, depth1694
			return false
		},
		/* 97 pnLocalEsc <- <('\\' ((&('%') '%') | (&('@') '@') | (&('#') '#') | (&('?') '?') | (&('/') '/') | (&('=') '=') | (&(';') ';') | (&(',') ',') | (&('+') '+') | (&('*') '*') | (&(')') ')') | (&('(') '(') | (&('\'') '\'') | (&('&') '&') | (&('$') '$') | (&('!') '!') | (&('-') '-') | (&('.') '.') | (&('~') '~') | (&('_') '_')))> */
		nil,
		/* 98 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 99 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 100 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 101 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 102 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 103 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 104 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1704, tokenIndex1704, depth1704 := position, tokenIndex, depth
			{
				position1705 := position
				depth++
				{
					position1706, tokenIndex1706, depth1706 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1707
					}
					position++
					goto l1706
				l1707:
					position, tokenIndex, depth = position1706, tokenIndex1706, depth1706
					if buffer[position] != rune('D') {
						goto l1704
					}
					position++
				}
			l1706:
				{
					position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1709
					}
					position++
					goto l1708
				l1709:
					position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
					if buffer[position] != rune('I') {
						goto l1704
					}
					position++
				}
			l1708:
				{
					position1710, tokenIndex1710, depth1710 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1711
					}
					position++
					goto l1710
				l1711:
					position, tokenIndex, depth = position1710, tokenIndex1710, depth1710
					if buffer[position] != rune('S') {
						goto l1704
					}
					position++
				}
			l1710:
				{
					position1712, tokenIndex1712, depth1712 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1713
					}
					position++
					goto l1712
				l1713:
					position, tokenIndex, depth = position1712, tokenIndex1712, depth1712
					if buffer[position] != rune('T') {
						goto l1704
					}
					position++
				}
			l1712:
				{
					position1714, tokenIndex1714, depth1714 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1715
					}
					position++
					goto l1714
				l1715:
					position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
					if buffer[position] != rune('I') {
						goto l1704
					}
					position++
				}
			l1714:
				{
					position1716, tokenIndex1716, depth1716 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1717
					}
					position++
					goto l1716
				l1717:
					position, tokenIndex, depth = position1716, tokenIndex1716, depth1716
					if buffer[position] != rune('N') {
						goto l1704
					}
					position++
				}
			l1716:
				{
					position1718, tokenIndex1718, depth1718 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1719
					}
					position++
					goto l1718
				l1719:
					position, tokenIndex, depth = position1718, tokenIndex1718, depth1718
					if buffer[position] != rune('C') {
						goto l1704
					}
					position++
				}
			l1718:
				{
					position1720, tokenIndex1720, depth1720 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1721
					}
					position++
					goto l1720
				l1721:
					position, tokenIndex, depth = position1720, tokenIndex1720, depth1720
					if buffer[position] != rune('T') {
						goto l1704
					}
					position++
				}
			l1720:
				if !rules[ruleskip]() {
					goto l1704
				}
				depth--
				add(ruleDISTINCT, position1705)
			}
			return true
		l1704:
			position, tokenIndex, depth = position1704, tokenIndex1704, depth1704
			return false
		},
		/* 105 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 106 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 107 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 108 LBRACE <- <('{' skip)> */
		func() bool {
			position1725, tokenIndex1725, depth1725 := position, tokenIndex, depth
			{
				position1726 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1725
				}
				position++
				if !rules[ruleskip]() {
					goto l1725
				}
				depth--
				add(ruleLBRACE, position1726)
			}
			return true
		l1725:
			position, tokenIndex, depth = position1725, tokenIndex1725, depth1725
			return false
		},
		/* 109 RBRACE <- <('}' skip)> */
		func() bool {
			position1727, tokenIndex1727, depth1727 := position, tokenIndex, depth
			{
				position1728 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1727
				}
				position++
				if !rules[ruleskip]() {
					goto l1727
				}
				depth--
				add(ruleRBRACE, position1728)
			}
			return true
		l1727:
			position, tokenIndex, depth = position1727, tokenIndex1727, depth1727
			return false
		},
		/* 110 LBRACK <- <('[' skip)> */
		nil,
		/* 111 RBRACK <- <(']' skip)> */
		nil,
		/* 112 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1731, tokenIndex1731, depth1731 := position, tokenIndex, depth
			{
				position1732 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1731
				}
				position++
				if !rules[ruleskip]() {
					goto l1731
				}
				depth--
				add(ruleSEMICOLON, position1732)
			}
			return true
		l1731:
			position, tokenIndex, depth = position1731, tokenIndex1731, depth1731
			return false
		},
		/* 113 COMMA <- <(',' skip)> */
		func() bool {
			position1733, tokenIndex1733, depth1733 := position, tokenIndex, depth
			{
				position1734 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1733
				}
				position++
				if !rules[ruleskip]() {
					goto l1733
				}
				depth--
				add(ruleCOMMA, position1734)
			}
			return true
		l1733:
			position, tokenIndex, depth = position1733, tokenIndex1733, depth1733
			return false
		},
		/* 114 DOT <- <('.' skip)> */
		func() bool {
			position1735, tokenIndex1735, depth1735 := position, tokenIndex, depth
			{
				position1736 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1735
				}
				position++
				if !rules[ruleskip]() {
					goto l1735
				}
				depth--
				add(ruleDOT, position1736)
			}
			return true
		l1735:
			position, tokenIndex, depth = position1735, tokenIndex1735, depth1735
			return false
		},
		/* 115 COLON <- <(':' skip)> */
		nil,
		/* 116 PIPE <- <('|' skip)> */
		func() bool {
			position1738, tokenIndex1738, depth1738 := position, tokenIndex, depth
			{
				position1739 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1738
				}
				position++
				if !rules[ruleskip]() {
					goto l1738
				}
				depth--
				add(rulePIPE, position1739)
			}
			return true
		l1738:
			position, tokenIndex, depth = position1738, tokenIndex1738, depth1738
			return false
		},
		/* 117 SLASH <- <('/' skip)> */
		func() bool {
			position1740, tokenIndex1740, depth1740 := position, tokenIndex, depth
			{
				position1741 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1740
				}
				position++
				if !rules[ruleskip]() {
					goto l1740
				}
				depth--
				add(ruleSLASH, position1741)
			}
			return true
		l1740:
			position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
			return false
		},
		/* 118 INVERSE <- <('^' skip)> */
		func() bool {
			position1742, tokenIndex1742, depth1742 := position, tokenIndex, depth
			{
				position1743 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1742
				}
				position++
				if !rules[ruleskip]() {
					goto l1742
				}
				depth--
				add(ruleINVERSE, position1743)
			}
			return true
		l1742:
			position, tokenIndex, depth = position1742, tokenIndex1742, depth1742
			return false
		},
		/* 119 LPAREN <- <('(' skip)> */
		func() bool {
			position1744, tokenIndex1744, depth1744 := position, tokenIndex, depth
			{
				position1745 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1744
				}
				position++
				if !rules[ruleskip]() {
					goto l1744
				}
				depth--
				add(ruleLPAREN, position1745)
			}
			return true
		l1744:
			position, tokenIndex, depth = position1744, tokenIndex1744, depth1744
			return false
		},
		/* 120 RPAREN <- <(')' skip)> */
		func() bool {
			position1746, tokenIndex1746, depth1746 := position, tokenIndex, depth
			{
				position1747 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1746
				}
				position++
				if !rules[ruleskip]() {
					goto l1746
				}
				depth--
				add(ruleRPAREN, position1747)
			}
			return true
		l1746:
			position, tokenIndex, depth = position1746, tokenIndex1746, depth1746
			return false
		},
		/* 121 ISA <- <('a' skip)> */
		func() bool {
			position1748, tokenIndex1748, depth1748 := position, tokenIndex, depth
			{
				position1749 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1748
				}
				position++
				if !rules[ruleskip]() {
					goto l1748
				}
				depth--
				add(ruleISA, position1749)
			}
			return true
		l1748:
			position, tokenIndex, depth = position1748, tokenIndex1748, depth1748
			return false
		},
		/* 122 NOT <- <('!' skip)> */
		func() bool {
			position1750, tokenIndex1750, depth1750 := position, tokenIndex, depth
			{
				position1751 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1750
				}
				position++
				if !rules[ruleskip]() {
					goto l1750
				}
				depth--
				add(ruleNOT, position1751)
			}
			return true
		l1750:
			position, tokenIndex, depth = position1750, tokenIndex1750, depth1750
			return false
		},
		/* 123 STAR <- <('*' skip)> */
		func() bool {
			position1752, tokenIndex1752, depth1752 := position, tokenIndex, depth
			{
				position1753 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1752
				}
				position++
				if !rules[ruleskip]() {
					goto l1752
				}
				depth--
				add(ruleSTAR, position1753)
			}
			return true
		l1752:
			position, tokenIndex, depth = position1752, tokenIndex1752, depth1752
			return false
		},
		/* 124 QUESTION <- <('?' skip)> */
		nil,
		/* 125 PLUS <- <('+' skip)> */
		func() bool {
			position1755, tokenIndex1755, depth1755 := position, tokenIndex, depth
			{
				position1756 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1755
				}
				position++
				if !rules[ruleskip]() {
					goto l1755
				}
				depth--
				add(rulePLUS, position1756)
			}
			return true
		l1755:
			position, tokenIndex, depth = position1755, tokenIndex1755, depth1755
			return false
		},
		/* 126 MINUS <- <('-' skip)> */
		func() bool {
			position1757, tokenIndex1757, depth1757 := position, tokenIndex, depth
			{
				position1758 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1757
				}
				position++
				if !rules[ruleskip]() {
					goto l1757
				}
				depth--
				add(ruleMINUS, position1758)
			}
			return true
		l1757:
			position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
			return false
		},
		/* 127 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 128 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 129 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 130 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 131 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1763, tokenIndex1763, depth1763 := position, tokenIndex, depth
			{
				position1764 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1763
				}
				position++
			l1765:
				{
					position1766, tokenIndex1766, depth1766 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1766
					}
					position++
					goto l1765
				l1766:
					position, tokenIndex, depth = position1766, tokenIndex1766, depth1766
				}
				if !rules[ruleskip]() {
					goto l1763
				}
				depth--
				add(ruleINTEGER, position1764)
			}
			return true
		l1763:
			position, tokenIndex, depth = position1763, tokenIndex1763, depth1763
			return false
		},
		/* 132 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 133 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 134 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 135 OR <- <('|' '|' skip)> */
		nil,
		/* 136 AND <- <('&' '&' skip)> */
		nil,
		/* 137 EQ <- <('=' skip)> */
		func() bool {
			position1772, tokenIndex1772, depth1772 := position, tokenIndex, depth
			{
				position1773 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1772
				}
				position++
				if !rules[ruleskip]() {
					goto l1772
				}
				depth--
				add(ruleEQ, position1773)
			}
			return true
		l1772:
			position, tokenIndex, depth = position1772, tokenIndex1772, depth1772
			return false
		},
		/* 138 NE <- <('!' '=' skip)> */
		nil,
		/* 139 GT <- <('>' skip)> */
		nil,
		/* 140 LT <- <('<' skip)> */
		nil,
		/* 141 LE <- <('<' '=' skip)> */
		nil,
		/* 142 GE <- <('>' '=' skip)> */
		nil,
		/* 143 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 144 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 145 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1781, tokenIndex1781, depth1781 := position, tokenIndex, depth
			{
				position1782 := position
				depth++
				{
					position1783, tokenIndex1783, depth1783 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1784
					}
					position++
					goto l1783
				l1784:
					position, tokenIndex, depth = position1783, tokenIndex1783, depth1783
					if buffer[position] != rune('A') {
						goto l1781
					}
					position++
				}
			l1783:
				{
					position1785, tokenIndex1785, depth1785 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1786
					}
					position++
					goto l1785
				l1786:
					position, tokenIndex, depth = position1785, tokenIndex1785, depth1785
					if buffer[position] != rune('S') {
						goto l1781
					}
					position++
				}
			l1785:
				if !rules[ruleskip]() {
					goto l1781
				}
				depth--
				add(ruleAS, position1782)
			}
			return true
		l1781:
			position, tokenIndex, depth = position1781, tokenIndex1781, depth1781
			return false
		},
		/* 146 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 147 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 148 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 149 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 150 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 151 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 152 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 153 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 154 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 155 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 156 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 157 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 158 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 159 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 160 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 161 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 162 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 163 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 164 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 165 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 166 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 167 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 168 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 169 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 170 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 171 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 172 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 173 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 174 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 175 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 176 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 177 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 178 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 179 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 180 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 181 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 182 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 183 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 184 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 185 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 186 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 187 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 188 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 189 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 190 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 191 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 192 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 193 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 194 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 195 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 196 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 197 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 198 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 199 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 200 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 201 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 202 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 203 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 204 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 205 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 206 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 207 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 208 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 209 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 210 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 211 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 212 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 213 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 214 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1855, tokenIndex1855, depth1855 := position, tokenIndex, depth
			{
				position1856 := position
				depth++
				{
					position1857, tokenIndex1857, depth1857 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1858
					}
					position++
					goto l1857
				l1858:
					position, tokenIndex, depth = position1857, tokenIndex1857, depth1857
					if buffer[position] != rune('B') {
						goto l1855
					}
					position++
				}
			l1857:
				{
					position1859, tokenIndex1859, depth1859 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1860
					}
					position++
					goto l1859
				l1860:
					position, tokenIndex, depth = position1859, tokenIndex1859, depth1859
					if buffer[position] != rune('Y') {
						goto l1855
					}
					position++
				}
			l1859:
				if !rules[ruleskip]() {
					goto l1855
				}
				depth--
				add(ruleBY, position1856)
			}
			return true
		l1855:
			position, tokenIndex, depth = position1855, tokenIndex1855, depth1855
			return false
		},
		/* 215 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 216 GRAPH <- <(('g' / 'G') ('r' / 'R') ('a' / 'A') ('p' / 'P') ('h' / 'H') skip)> */
		nil,
		/* 217 MINUSSETOPER <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('s' / 'S') skip)> */
		nil,
		/* 218 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1865 := position
				depth++
			l1866:
				{
					position1867, tokenIndex1867, depth1867 := position, tokenIndex, depth
					{
						position1868, tokenIndex1868, depth1868 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1869
						}
						goto l1868
					l1869:
						position, tokenIndex, depth = position1868, tokenIndex1868, depth1868
						{
							position1870 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1867
							}
							position++
						l1871:
							{
								position1872, tokenIndex1872, depth1872 := position, tokenIndex, depth
								{
									position1873, tokenIndex1873, depth1873 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1873
									}
									goto l1872
								l1873:
									position, tokenIndex, depth = position1873, tokenIndex1873, depth1873
								}
								if !matchDot() {
									goto l1872
								}
								goto l1871
							l1872:
								position, tokenIndex, depth = position1872, tokenIndex1872, depth1872
							}
							if !rules[ruleendOfLine]() {
								goto l1867
							}
							depth--
							add(rulecomment, position1870)
						}
					}
				l1868:
					goto l1866
				l1867:
					position, tokenIndex, depth = position1867, tokenIndex1867, depth1867
				}
				depth--
				add(ruleskip, position1865)
			}
			return true
		},
		/* 219 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1874, tokenIndex1874, depth1874 := position, tokenIndex, depth
			{
				position1875 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1874
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1874
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1874
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1874
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1874
						}
						break
					}
				}

				depth--
				add(rulews, position1875)
			}
			return true
		l1874:
			position, tokenIndex, depth = position1874, tokenIndex1874, depth1874
			return false
		},
		/* 220 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 221 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1878, tokenIndex1878, depth1878 := position, tokenIndex, depth
			{
				position1879 := position
				depth++
				{
					position1880, tokenIndex1880, depth1880 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1881
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1881
					}
					position++
					goto l1880
				l1881:
					position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
					if buffer[position] != rune('\n') {
						goto l1882
					}
					position++
					goto l1880
				l1882:
					position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
					if buffer[position] != rune('\r') {
						goto l1878
					}
					position++
				}
			l1880:
				depth--
				add(ruleendOfLine, position1879)
			}
			return true
		l1878:
			position, tokenIndex, depth = position1878, tokenIndex1878, depth1878
			return false
		},
		/* 223 SERVICE <- <> */
		nil,
		/* 224 SILENT <- <> */
		nil,
	}
	p.rules = rules
}
