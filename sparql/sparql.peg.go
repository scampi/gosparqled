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
	rules  [222]func() bool
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
							switch buffer[position] {
							case 'M', 'm':
								{
									position226 := position
									depth++
									{
										position227 := position
										depth++
										{
											position228, tokenIndex228, depth228 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l229
											}
											position++
											goto l228
										l229:
											position, tokenIndex, depth = position228, tokenIndex228, depth228
											if buffer[position] != rune('M') {
												goto l222
											}
											position++
										}
									l228:
										{
											position230, tokenIndex230, depth230 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l231
											}
											position++
											goto l230
										l231:
											position, tokenIndex, depth = position230, tokenIndex230, depth230
											if buffer[position] != rune('I') {
												goto l222
											}
											position++
										}
									l230:
										{
											position232, tokenIndex232, depth232 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l233
											}
											position++
											goto l232
										l233:
											position, tokenIndex, depth = position232, tokenIndex232, depth232
											if buffer[position] != rune('N') {
												goto l222
											}
											position++
										}
									l232:
										{
											position234, tokenIndex234, depth234 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l235
											}
											position++
											goto l234
										l235:
											position, tokenIndex, depth = position234, tokenIndex234, depth234
											if buffer[position] != rune('U') {
												goto l222
											}
											position++
										}
									l234:
										{
											position236, tokenIndex236, depth236 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l237
											}
											position++
											goto l236
										l237:
											position, tokenIndex, depth = position236, tokenIndex236, depth236
											if buffer[position] != rune('S') {
												goto l222
											}
											position++
										}
									l236:
										if !rules[ruleskip]() {
											goto l222
										}
										depth--
										add(ruleMINUSSETOPER, position227)
									}
									if !rules[rulegroupGraphPattern]() {
										goto l222
									}
									depth--
									add(ruleminusGraphPattern, position226)
								}
								break
							case 'G', 'g':
								{
									position238 := position
									depth++
									{
										position239 := position
										depth++
										{
											position240, tokenIndex240, depth240 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l241
											}
											position++
											goto l240
										l241:
											position, tokenIndex, depth = position240, tokenIndex240, depth240
											if buffer[position] != rune('G') {
												goto l222
											}
											position++
										}
									l240:
										{
											position242, tokenIndex242, depth242 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l243
											}
											position++
											goto l242
										l243:
											position, tokenIndex, depth = position242, tokenIndex242, depth242
											if buffer[position] != rune('R') {
												goto l222
											}
											position++
										}
									l242:
										{
											position244, tokenIndex244, depth244 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l245
											}
											position++
											goto l244
										l245:
											position, tokenIndex, depth = position244, tokenIndex244, depth244
											if buffer[position] != rune('A') {
												goto l222
											}
											position++
										}
									l244:
										{
											position246, tokenIndex246, depth246 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l247
											}
											position++
											goto l246
										l247:
											position, tokenIndex, depth = position246, tokenIndex246, depth246
											if buffer[position] != rune('P') {
												goto l222
											}
											position++
										}
									l246:
										{
											position248, tokenIndex248, depth248 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l249
											}
											position++
											goto l248
										l249:
											position, tokenIndex, depth = position248, tokenIndex248, depth248
											if buffer[position] != rune('H') {
												goto l222
											}
											position++
										}
									l248:
										if !rules[ruleskip]() {
											goto l222
										}
										depth--
										add(ruleGRAPH, position239)
									}
									{
										position250, tokenIndex250, depth250 := position, tokenIndex, depth
										if !rules[rulevar]() {
											goto l251
										}
										goto l250
									l251:
										position, tokenIndex, depth = position250, tokenIndex250, depth250
										if !rules[ruleiriref]() {
											goto l222
										}
									}
								l250:
									if !rules[rulegroupGraphPattern]() {
										goto l222
									}
									depth--
									add(rulegraphGraphPattern, position238)
								}
								break
							case '{':
								if !rules[rulegroupOrUnionGraphPattern]() {
									goto l222
								}
								break
							default:
								{
									position252 := position
									depth++
									{
										position253 := position
										depth++
										{
											position254, tokenIndex254, depth254 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l255
											}
											position++
											goto l254
										l255:
											position, tokenIndex, depth = position254, tokenIndex254, depth254
											if buffer[position] != rune('O') {
												goto l222
											}
											position++
										}
									l254:
										{
											position256, tokenIndex256, depth256 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l257
											}
											position++
											goto l256
										l257:
											position, tokenIndex, depth = position256, tokenIndex256, depth256
											if buffer[position] != rune('P') {
												goto l222
											}
											position++
										}
									l256:
										{
											position258, tokenIndex258, depth258 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l259
											}
											position++
											goto l258
										l259:
											position, tokenIndex, depth = position258, tokenIndex258, depth258
											if buffer[position] != rune('T') {
												goto l222
											}
											position++
										}
									l258:
										{
											position260, tokenIndex260, depth260 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l261
											}
											position++
											goto l260
										l261:
											position, tokenIndex, depth = position260, tokenIndex260, depth260
											if buffer[position] != rune('I') {
												goto l222
											}
											position++
										}
									l260:
										{
											position262, tokenIndex262, depth262 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l263
											}
											position++
											goto l262
										l263:
											position, tokenIndex, depth = position262, tokenIndex262, depth262
											if buffer[position] != rune('O') {
												goto l222
											}
											position++
										}
									l262:
										{
											position264, tokenIndex264, depth264 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l265
											}
											position++
											goto l264
										l265:
											position, tokenIndex, depth = position264, tokenIndex264, depth264
											if buffer[position] != rune('N') {
												goto l222
											}
											position++
										}
									l264:
										{
											position266, tokenIndex266, depth266 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l267
											}
											position++
											goto l266
										l267:
											position, tokenIndex, depth = position266, tokenIndex266, depth266
											if buffer[position] != rune('A') {
												goto l222
											}
											position++
										}
									l266:
										{
											position268, tokenIndex268, depth268 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l269
											}
											position++
											goto l268
										l269:
											position, tokenIndex, depth = position268, tokenIndex268, depth268
											if buffer[position] != rune('L') {
												goto l222
											}
											position++
										}
									l268:
										if !rules[ruleskip]() {
											goto l222
										}
										depth--
										add(ruleOPTIONAL, position253)
									}
									if !rules[ruleLBRACE]() {
										goto l222
									}
									{
										position270, tokenIndex270, depth270 := position, tokenIndex, depth
										if !rules[rulesubSelect]() {
											goto l271
										}
										goto l270
									l271:
										position, tokenIndex, depth = position270, tokenIndex270, depth270
										if !rules[rulegraphPattern]() {
											goto l222
										}
									}
								l270:
									if !rules[ruleRBRACE]() {
										goto l222
									}
									depth--
									add(ruleoptionalGraphPattern, position252)
								}
								break
							}
						}

						depth--
						add(rulegraphPatternNotTriples, position224)
					}
					{
						position272, tokenIndex272, depth272 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l272
						}
						goto l273
					l272:
						position, tokenIndex, depth = position272, tokenIndex272, depth272
					}
				l273:
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
		/* 18 graphPatternNotTriples <- <((&('M' | 'm') minusGraphPattern) | (&('G' | 'g') graphGraphPattern) | (&('{') groupOrUnionGraphPattern) | (&('O' | 'o') optionalGraphPattern))> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position276, tokenIndex276, depth276 := position, tokenIndex, depth
			{
				position277 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l276
				}
				{
					position278, tokenIndex278, depth278 := position, tokenIndex, depth
					{
						position280 := position
						depth++
						{
							position281, tokenIndex281, depth281 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l282
							}
							position++
							goto l281
						l282:
							position, tokenIndex, depth = position281, tokenIndex281, depth281
							if buffer[position] != rune('U') {
								goto l278
							}
							position++
						}
					l281:
						{
							position283, tokenIndex283, depth283 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l284
							}
							position++
							goto l283
						l284:
							position, tokenIndex, depth = position283, tokenIndex283, depth283
							if buffer[position] != rune('N') {
								goto l278
							}
							position++
						}
					l283:
						{
							position285, tokenIndex285, depth285 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l286
							}
							position++
							goto l285
						l286:
							position, tokenIndex, depth = position285, tokenIndex285, depth285
							if buffer[position] != rune('I') {
								goto l278
							}
							position++
						}
					l285:
						{
							position287, tokenIndex287, depth287 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l288
							}
							position++
							goto l287
						l288:
							position, tokenIndex, depth = position287, tokenIndex287, depth287
							if buffer[position] != rune('O') {
								goto l278
							}
							position++
						}
					l287:
						{
							position289, tokenIndex289, depth289 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l290
							}
							position++
							goto l289
						l290:
							position, tokenIndex, depth = position289, tokenIndex289, depth289
							if buffer[position] != rune('N') {
								goto l278
							}
							position++
						}
					l289:
						if !rules[ruleskip]() {
							goto l278
						}
						depth--
						add(ruleUNION, position280)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l278
					}
					goto l279
				l278:
					position, tokenIndex, depth = position278, tokenIndex278, depth278
				}
			l279:
				depth--
				add(rulegroupOrUnionGraphPattern, position277)
			}
			return true
		l276:
			position, tokenIndex, depth = position276, tokenIndex276, depth276
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
			position294, tokenIndex294, depth294 := position, tokenIndex, depth
			{
				position295 := position
				depth++
				{
					position296, tokenIndex296, depth296 := position, tokenIndex, depth
					{
						position298 := position
						depth++
						{
							position299, tokenIndex299, depth299 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l300
							}
							position++
							goto l299
						l300:
							position, tokenIndex, depth = position299, tokenIndex299, depth299
							if buffer[position] != rune('F') {
								goto l297
							}
							position++
						}
					l299:
						{
							position301, tokenIndex301, depth301 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l302
							}
							position++
							goto l301
						l302:
							position, tokenIndex, depth = position301, tokenIndex301, depth301
							if buffer[position] != rune('I') {
								goto l297
							}
							position++
						}
					l301:
						{
							position303, tokenIndex303, depth303 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l304
							}
							position++
							goto l303
						l304:
							position, tokenIndex, depth = position303, tokenIndex303, depth303
							if buffer[position] != rune('L') {
								goto l297
							}
							position++
						}
					l303:
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l306
							}
							position++
							goto l305
						l306:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
							if buffer[position] != rune('T') {
								goto l297
							}
							position++
						}
					l305:
						{
							position307, tokenIndex307, depth307 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l308
							}
							position++
							goto l307
						l308:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if buffer[position] != rune('E') {
								goto l297
							}
							position++
						}
					l307:
						{
							position309, tokenIndex309, depth309 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l310
							}
							position++
							goto l309
						l310:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if buffer[position] != rune('R') {
								goto l297
							}
							position++
						}
					l309:
						if !rules[ruleskip]() {
							goto l297
						}
						depth--
						add(ruleFILTER, position298)
					}
					if !rules[ruleconstraint]() {
						goto l297
					}
					goto l296
				l297:
					position, tokenIndex, depth = position296, tokenIndex296, depth296
					{
						position311 := position
						depth++
						{
							position312, tokenIndex312, depth312 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l313
							}
							position++
							goto l312
						l313:
							position, tokenIndex, depth = position312, tokenIndex312, depth312
							if buffer[position] != rune('B') {
								goto l294
							}
							position++
						}
					l312:
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('I') {
								goto l294
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('N') {
								goto l294
							}
							position++
						}
					l316:
						{
							position318, tokenIndex318, depth318 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l319
							}
							position++
							goto l318
						l319:
							position, tokenIndex, depth = position318, tokenIndex318, depth318
							if buffer[position] != rune('D') {
								goto l294
							}
							position++
						}
					l318:
						if !rules[ruleskip]() {
							goto l294
						}
						depth--
						add(ruleBIND, position311)
					}
					if !rules[ruleLPAREN]() {
						goto l294
					}
					if !rules[ruleexpression]() {
						goto l294
					}
					if !rules[ruleAS]() {
						goto l294
					}
					if !rules[rulevar]() {
						goto l294
					}
					if !rules[ruleRPAREN]() {
						goto l294
					}
				}
			l296:
				depth--
				add(rulefilterOrBind, position295)
			}
			return true
		l294:
			position, tokenIndex, depth = position294, tokenIndex294, depth294
			return false
		},
		/* 25 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		func() bool {
			position320, tokenIndex320, depth320 := position, tokenIndex, depth
			{
				position321 := position
				depth++
				{
					position322, tokenIndex322, depth322 := position, tokenIndex, depth
					if !rules[rulebrackettedExpression]() {
						goto l323
					}
					goto l322
				l323:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if !rules[rulebuiltinCall]() {
						goto l324
					}
					goto l322
				l324:
					position, tokenIndex, depth = position322, tokenIndex322, depth322
					if !rules[rulefunctionCall]() {
						goto l320
					}
				}
			l322:
				depth--
				add(ruleconstraint, position321)
			}
			return true
		l320:
			position, tokenIndex, depth = position320, tokenIndex320, depth320
			return false
		},
		/* 26 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l325
				}
			l327:
				{
					position328, tokenIndex328, depth328 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l328
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l328
					}
					goto l327
				l328:
					position, tokenIndex, depth = position328, tokenIndex328, depth328
				}
				{
					position329, tokenIndex329, depth329 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l329
					}
					goto l330
				l329:
					position, tokenIndex, depth = position329, tokenIndex329, depth329
				}
			l330:
				depth--
				add(ruletriplesBlock, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
			return false
		},
		/* 27 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath?))> */
		func() bool {
			position331, tokenIndex331, depth331 := position, tokenIndex, depth
			{
				position332 := position
				depth++
				{
					position333, tokenIndex333, depth333 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l334
					}
					if !rules[rulepropertyListPath]() {
						goto l334
					}
					goto l333
				l334:
					position, tokenIndex, depth = position333, tokenIndex333, depth333
					if !rules[ruletriplesNodePath]() {
						goto l331
					}
					{
						position335, tokenIndex335, depth335 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l335
						}
						goto l336
					l335:
						position, tokenIndex, depth = position335, tokenIndex335, depth335
					}
				l336:
				}
			l333:
				depth--
				add(ruletriplesSameSubjectPath, position332)
			}
			return true
		l331:
			position, tokenIndex, depth = position331, tokenIndex331, depth331
			return false
		},
		/* 28 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position337, tokenIndex337, depth337 := position, tokenIndex, depth
			{
				position338 := position
				depth++
				{
					position339, tokenIndex339, depth339 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l340
					}
					goto l339
				l340:
					position, tokenIndex, depth = position339, tokenIndex339, depth339
					{
						position341 := position
						depth++
						{
							position342, tokenIndex342, depth342 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l343
							}
							goto l342
						l343:
							position, tokenIndex, depth = position342, tokenIndex342, depth342
							{
								switch buffer[position] {
								case '(':
									if !rules[rulenil]() {
										goto l337
									}
									break
								case '[', '_':
									{
										position345 := position
										depth++
										{
											position346, tokenIndex346, depth346 := position, tokenIndex, depth
											{
												position348 := position
												depth++
												if buffer[position] != rune('_') {
													goto l347
												}
												position++
												if buffer[position] != rune(':') {
													goto l347
												}
												position++
												{
													position349, tokenIndex349, depth349 := position, tokenIndex, depth
													if !rules[rulepnCharsU]() {
														goto l350
													}
													goto l349
												l350:
													position, tokenIndex, depth = position349, tokenIndex349, depth349
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l347
													}
													position++
												}
											l349:
												{
													position351, tokenIndex351, depth351 := position, tokenIndex, depth
													{
														position353, tokenIndex353, depth353 := position, tokenIndex, depth
													l355:
														{
															position356, tokenIndex356, depth356 := position, tokenIndex, depth
															{
																position357, tokenIndex357, depth357 := position, tokenIndex, depth
																if !rules[rulepnCharsU]() {
																	goto l358
																}
																goto l357
															l358:
																position, tokenIndex, depth = position357, tokenIndex357, depth357
																{
																	switch buffer[position] {
																	case '.':
																		if buffer[position] != rune('.') {
																			goto l356
																		}
																		position++
																		break
																	case '-':
																		if buffer[position] != rune('-') {
																			goto l356
																		}
																		position++
																		break
																	default:
																		if c := buffer[position]; c < rune('0') || c > rune('9') {
																			goto l356
																		}
																		position++
																		break
																	}
																}

															}
														l357:
															goto l355
														l356:
															position, tokenIndex, depth = position356, tokenIndex356, depth356
														}
														if !rules[rulepnCharsU]() {
															goto l354
														}
														goto l353
													l354:
														position, tokenIndex, depth = position353, tokenIndex353, depth353
														{
															position360, tokenIndex360, depth360 := position, tokenIndex, depth
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l361
															}
															position++
															goto l360
														l361:
															position, tokenIndex, depth = position360, tokenIndex360, depth360
															if buffer[position] != rune('-') {
																goto l351
															}
															position++
														}
													l360:
													}
												l353:
													goto l352
												l351:
													position, tokenIndex, depth = position351, tokenIndex351, depth351
												}
											l352:
												if !rules[ruleskip]() {
													goto l347
												}
												depth--
												add(ruleblankNodeLabel, position348)
											}
											goto l346
										l347:
											position, tokenIndex, depth = position346, tokenIndex346, depth346
											{
												position362 := position
												depth++
												if buffer[position] != rune('[') {
													goto l337
												}
												position++
											l363:
												{
													position364, tokenIndex364, depth364 := position, tokenIndex, depth
													if !rules[rulews]() {
														goto l364
													}
													goto l363
												l364:
													position, tokenIndex, depth = position364, tokenIndex364, depth364
												}
												if buffer[position] != rune(']') {
													goto l337
												}
												position++
												if !rules[ruleskip]() {
													goto l337
												}
												depth--
												add(ruleanon, position362)
											}
										}
									l346:
										depth--
										add(ruleblankNode, position345)
									}
									break
								case 'F', 'T', 'f', 't':
									if !rules[rulebooleanLiteral]() {
										goto l337
									}
									break
								case '"', '\'':
									if !rules[ruleliteral]() {
										goto l337
									}
									break
								default:
									if !rules[rulenumericLiteral]() {
										goto l337
									}
									break
								}
							}

						}
					l342:
						depth--
						add(rulegraphTerm, position341)
					}
				}
			l339:
				depth--
				add(rulevarOrTerm, position338)
			}
			return true
		l337:
			position, tokenIndex, depth = position337, tokenIndex337, depth337
			return false
		},
		/* 29 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 30 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		func() bool {
			position366, tokenIndex366, depth366 := position, tokenIndex, depth
			{
				position367 := position
				depth++
				{
					position368, tokenIndex368, depth368 := position, tokenIndex, depth
					{
						position370 := position
						depth++
						if !rules[ruleLPAREN]() {
							goto l369
						}
						if !rules[rulegraphNodePath]() {
							goto l369
						}
					l371:
						{
							position372, tokenIndex372, depth372 := position, tokenIndex, depth
							if !rules[rulegraphNodePath]() {
								goto l372
							}
							goto l371
						l372:
							position, tokenIndex, depth = position372, tokenIndex372, depth372
						}
						if !rules[ruleRPAREN]() {
							goto l369
						}
						depth--
						add(rulecollectionPath, position370)
					}
					goto l368
				l369:
					position, tokenIndex, depth = position368, tokenIndex368, depth368
					{
						position373 := position
						depth++
						{
							position374 := position
							depth++
							if buffer[position] != rune('[') {
								goto l366
							}
							position++
							if !rules[ruleskip]() {
								goto l366
							}
							depth--
							add(ruleLBRACK, position374)
						}
						if !rules[rulepropertyListPath]() {
							goto l366
						}
						{
							position375 := position
							depth++
							if buffer[position] != rune(']') {
								goto l366
							}
							position++
							if !rules[ruleskip]() {
								goto l366
							}
							depth--
							add(ruleRBRACK, position375)
						}
						depth--
						add(ruleblankNodePropertyListPath, position373)
					}
				}
			l368:
				depth--
				add(ruletriplesNodePath, position367)
			}
			return true
		l366:
			position, tokenIndex, depth = position366, tokenIndex366, depth366
			return false
		},
		/* 31 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 32 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 33 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position378, tokenIndex378, depth378 := position, tokenIndex, depth
			{
				position379 := position
				depth++
				{
					position380, tokenIndex380, depth380 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l381
					}
					goto l380
				l381:
					position, tokenIndex, depth = position380, tokenIndex380, depth380
					{
						position382 := position
						depth++
						if !rules[rulepath]() {
							goto l378
						}
						depth--
						add(ruleverbPath, position382)
					}
				}
			l380:
				{
					position383 := position
					depth++
					if !rules[ruleobjectPath]() {
						goto l378
					}
				l384:
					{
						position385, tokenIndex385, depth385 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l385
						}
						if !rules[ruleobjectPath]() {
							goto l385
						}
						goto l384
					l385:
						position, tokenIndex, depth = position385, tokenIndex385, depth385
					}
					depth--
					add(ruleobjectListPath, position383)
				}
				{
					position386, tokenIndex386, depth386 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l386
					}
					{
						position388, tokenIndex388, depth388 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l388
						}
						goto l389
					l388:
						position, tokenIndex, depth = position388, tokenIndex388, depth388
					}
				l389:
					goto l387
				l386:
					position, tokenIndex, depth = position386, tokenIndex386, depth386
				}
			l387:
				depth--
				add(rulepropertyListPath, position379)
			}
			return true
		l378:
			position, tokenIndex, depth = position378, tokenIndex378, depth378
			return false
		},
		/* 34 verbPath <- <path> */
		nil,
		/* 35 path <- <pathAlternative> */
		func() bool {
			position391, tokenIndex391, depth391 := position, tokenIndex, depth
			{
				position392 := position
				depth++
				{
					position393 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l391
					}
				l394:
					{
						position395, tokenIndex395, depth395 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l395
						}
						if !rules[rulepathSequence]() {
							goto l395
						}
						goto l394
					l395:
						position, tokenIndex, depth = position395, tokenIndex395, depth395
					}
					depth--
					add(rulepathAlternative, position393)
				}
				depth--
				add(rulepath, position392)
			}
			return true
		l391:
			position, tokenIndex, depth = position391, tokenIndex391, depth391
			return false
		},
		/* 36 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 37 pathSequence <- <(pathElt (SLASH pathElt)*)> */
		func() bool {
			position397, tokenIndex397, depth397 := position, tokenIndex, depth
			{
				position398 := position
				depth++
				if !rules[rulepathElt]() {
					goto l397
				}
			l399:
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l400
					}
					if !rules[rulepathElt]() {
						goto l400
					}
					goto l399
				l400:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
				}
				depth--
				add(rulepathSequence, position398)
			}
			return true
		l397:
			position, tokenIndex, depth = position397, tokenIndex397, depth397
			return false
		},
		/* 38 pathElt <- <(INVERSE? pathPrimary pathMod?)> */
		func() bool {
			position401, tokenIndex401, depth401 := position, tokenIndex, depth
			{
				position402 := position
				depth++
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if !rules[ruleINVERSE]() {
						goto l403
					}
					goto l404
				l403:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
				}
			l404:
				{
					position405 := position
					depth++
					{
						position406, tokenIndex406, depth406 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l407
						}
						goto l406
					l407:
						position, tokenIndex, depth = position406, tokenIndex406, depth406
						{
							switch buffer[position] {
							case '(':
								if !rules[ruleLPAREN]() {
									goto l401
								}
								if !rules[rulepath]() {
									goto l401
								}
								if !rules[ruleRPAREN]() {
									goto l401
								}
								break
							case '!':
								if !rules[ruleNOT]() {
									goto l401
								}
								{
									position409 := position
									depth++
									{
										position410, tokenIndex410, depth410 := position, tokenIndex, depth
										if !rules[rulepathOneInPropertySet]() {
											goto l411
										}
										goto l410
									l411:
										position, tokenIndex, depth = position410, tokenIndex410, depth410
										if !rules[ruleLPAREN]() {
											goto l401
										}
										{
											position412, tokenIndex412, depth412 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l412
											}
										l414:
											{
												position415, tokenIndex415, depth415 := position, tokenIndex, depth
												if !rules[rulePIPE]() {
													goto l415
												}
												if !rules[rulepathOneInPropertySet]() {
													goto l415
												}
												goto l414
											l415:
												position, tokenIndex, depth = position415, tokenIndex415, depth415
											}
											goto l413
										l412:
											position, tokenIndex, depth = position412, tokenIndex412, depth412
										}
									l413:
										if !rules[ruleRPAREN]() {
											goto l401
										}
									}
								l410:
									depth--
									add(rulepathNegatedPropertySet, position409)
								}
								break
							default:
								if !rules[ruleISA]() {
									goto l401
								}
								break
							}
						}

					}
				l406:
					depth--
					add(rulepathPrimary, position405)
				}
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					{
						position418 := position
						depth++
						{
							switch buffer[position] {
							case '+':
								if !rules[rulePLUS]() {
									goto l416
								}
								break
							case '?':
								{
									position420 := position
									depth++
									if buffer[position] != rune('?') {
										goto l416
									}
									position++
									if !rules[ruleskip]() {
										goto l416
									}
									depth--
									add(ruleQUESTION, position420)
								}
								break
							default:
								if !rules[ruleSTAR]() {
									goto l416
								}
								break
							}
						}

						{
							position421, tokenIndex421, depth421 := position, tokenIndex, depth
							if !matchDot() {
								goto l421
							}
							goto l416
						l421:
							position, tokenIndex, depth = position421, tokenIndex421, depth421
						}
						depth--
						add(rulepathMod, position418)
					}
					goto l417
				l416:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
				}
			l417:
				depth--
				add(rulepathElt, position402)
			}
			return true
		l401:
			position, tokenIndex, depth = position401, tokenIndex401, depth401
			return false
		},
		/* 39 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 40 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 41 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position424, tokenIndex424, depth424 := position, tokenIndex, depth
			{
				position425 := position
				depth++
				{
					position426, tokenIndex426, depth426 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l427
					}
					goto l426
				l427:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if !rules[ruleISA]() {
						goto l428
					}
					goto l426
				l428:
					position, tokenIndex, depth = position426, tokenIndex426, depth426
					if !rules[ruleINVERSE]() {
						goto l424
					}
					{
						position429, tokenIndex429, depth429 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l430
						}
						goto l429
					l430:
						position, tokenIndex, depth = position429, tokenIndex429, depth429
						if !rules[ruleISA]() {
							goto l424
						}
					}
				l429:
				}
			l426:
				depth--
				add(rulepathOneInPropertySet, position425)
			}
			return true
		l424:
			position, tokenIndex, depth = position424, tokenIndex424, depth424
			return false
		},
		/* 42 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !.)> */
		nil,
		/* 43 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 44 objectPath <- <graphNodePath> */
		func() bool {
			position433, tokenIndex433, depth433 := position, tokenIndex, depth
			{
				position434 := position
				depth++
				if !rules[rulegraphNodePath]() {
					goto l433
				}
				depth--
				add(ruleobjectPath, position434)
			}
			return true
		l433:
			position, tokenIndex, depth = position433, tokenIndex433, depth433
			return false
		},
		/* 45 graphNodePath <- <(varOrTerm / triplesNodePath)> */
		func() bool {
			position435, tokenIndex435, depth435 := position, tokenIndex, depth
			{
				position436 := position
				depth++
				{
					position437, tokenIndex437, depth437 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l438
					}
					goto l437
				l438:
					position, tokenIndex, depth = position437, tokenIndex437, depth437
					if !rules[ruletriplesNodePath]() {
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
		/* 46 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position440 := position
				depth++
				{
					position441, tokenIndex441, depth441 := position, tokenIndex, depth
					{
						position443, tokenIndex443, depth443 := position, tokenIndex, depth
						{
							position445 := position
							depth++
							{
								position446, tokenIndex446, depth446 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l447
								}
								position++
								goto l446
							l447:
								position, tokenIndex, depth = position446, tokenIndex446, depth446
								if buffer[position] != rune('O') {
									goto l444
								}
								position++
							}
						l446:
							{
								position448, tokenIndex448, depth448 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l449
								}
								position++
								goto l448
							l449:
								position, tokenIndex, depth = position448, tokenIndex448, depth448
								if buffer[position] != rune('R') {
									goto l444
								}
								position++
							}
						l448:
							{
								position450, tokenIndex450, depth450 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l451
								}
								position++
								goto l450
							l451:
								position, tokenIndex, depth = position450, tokenIndex450, depth450
								if buffer[position] != rune('D') {
									goto l444
								}
								position++
							}
						l450:
							{
								position452, tokenIndex452, depth452 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l453
								}
								position++
								goto l452
							l453:
								position, tokenIndex, depth = position452, tokenIndex452, depth452
								if buffer[position] != rune('E') {
									goto l444
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
									goto l444
								}
								position++
							}
						l454:
							if !rules[ruleskip]() {
								goto l444
							}
							depth--
							add(ruleORDER, position445)
						}
						if !rules[ruleBY]() {
							goto l444
						}
						{
							position458 := position
							depth++
							{
								position459, tokenIndex459, depth459 := position, tokenIndex, depth
								{
									position461, tokenIndex461, depth461 := position, tokenIndex, depth
									{
										position463, tokenIndex463, depth463 := position, tokenIndex, depth
										{
											position465 := position
											depth++
											{
												position466, tokenIndex466, depth466 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l467
												}
												position++
												goto l466
											l467:
												position, tokenIndex, depth = position466, tokenIndex466, depth466
												if buffer[position] != rune('A') {
													goto l464
												}
												position++
											}
										l466:
											{
												position468, tokenIndex468, depth468 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l469
												}
												position++
												goto l468
											l469:
												position, tokenIndex, depth = position468, tokenIndex468, depth468
												if buffer[position] != rune('S') {
													goto l464
												}
												position++
											}
										l468:
											{
												position470, tokenIndex470, depth470 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l471
												}
												position++
												goto l470
											l471:
												position, tokenIndex, depth = position470, tokenIndex470, depth470
												if buffer[position] != rune('C') {
													goto l464
												}
												position++
											}
										l470:
											if !rules[ruleskip]() {
												goto l464
											}
											depth--
											add(ruleASC, position465)
										}
										goto l463
									l464:
										position, tokenIndex, depth = position463, tokenIndex463, depth463
										{
											position472 := position
											depth++
											{
												position473, tokenIndex473, depth473 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l474
												}
												position++
												goto l473
											l474:
												position, tokenIndex, depth = position473, tokenIndex473, depth473
												if buffer[position] != rune('D') {
													goto l461
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
													goto l461
												}
												position++
											}
										l475:
											{
												position477, tokenIndex477, depth477 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l478
												}
												position++
												goto l477
											l478:
												position, tokenIndex, depth = position477, tokenIndex477, depth477
												if buffer[position] != rune('S') {
													goto l461
												}
												position++
											}
										l477:
											{
												position479, tokenIndex479, depth479 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l480
												}
												position++
												goto l479
											l480:
												position, tokenIndex, depth = position479, tokenIndex479, depth479
												if buffer[position] != rune('C') {
													goto l461
												}
												position++
											}
										l479:
											if !rules[ruleskip]() {
												goto l461
											}
											depth--
											add(ruleDESC, position472)
										}
									}
								l463:
									goto l462
								l461:
									position, tokenIndex, depth = position461, tokenIndex461, depth461
								}
							l462:
								if !rules[rulebrackettedExpression]() {
									goto l460
								}
								goto l459
							l460:
								position, tokenIndex, depth = position459, tokenIndex459, depth459
								if !rules[rulefunctionCall]() {
									goto l481
								}
								goto l459
							l481:
								position, tokenIndex, depth = position459, tokenIndex459, depth459
								if !rules[rulebuiltinCall]() {
									goto l482
								}
								goto l459
							l482:
								position, tokenIndex, depth = position459, tokenIndex459, depth459
								if !rules[rulevar]() {
									goto l444
								}
							}
						l459:
							depth--
							add(ruleorderCondition, position458)
						}
					l456:
						{
							position457, tokenIndex457, depth457 := position, tokenIndex, depth
							{
								position483 := position
								depth++
								{
									position484, tokenIndex484, depth484 := position, tokenIndex, depth
									{
										position486, tokenIndex486, depth486 := position, tokenIndex, depth
										{
											position488, tokenIndex488, depth488 := position, tokenIndex, depth
											{
												position490 := position
												depth++
												{
													position491, tokenIndex491, depth491 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l492
													}
													position++
													goto l491
												l492:
													position, tokenIndex, depth = position491, tokenIndex491, depth491
													if buffer[position] != rune('A') {
														goto l489
													}
													position++
												}
											l491:
												{
													position493, tokenIndex493, depth493 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l494
													}
													position++
													goto l493
												l494:
													position, tokenIndex, depth = position493, tokenIndex493, depth493
													if buffer[position] != rune('S') {
														goto l489
													}
													position++
												}
											l493:
												{
													position495, tokenIndex495, depth495 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l496
													}
													position++
													goto l495
												l496:
													position, tokenIndex, depth = position495, tokenIndex495, depth495
													if buffer[position] != rune('C') {
														goto l489
													}
													position++
												}
											l495:
												if !rules[ruleskip]() {
													goto l489
												}
												depth--
												add(ruleASC, position490)
											}
											goto l488
										l489:
											position, tokenIndex, depth = position488, tokenIndex488, depth488
											{
												position497 := position
												depth++
												{
													position498, tokenIndex498, depth498 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l499
													}
													position++
													goto l498
												l499:
													position, tokenIndex, depth = position498, tokenIndex498, depth498
													if buffer[position] != rune('D') {
														goto l486
													}
													position++
												}
											l498:
												{
													position500, tokenIndex500, depth500 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l501
													}
													position++
													goto l500
												l501:
													position, tokenIndex, depth = position500, tokenIndex500, depth500
													if buffer[position] != rune('E') {
														goto l486
													}
													position++
												}
											l500:
												{
													position502, tokenIndex502, depth502 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l503
													}
													position++
													goto l502
												l503:
													position, tokenIndex, depth = position502, tokenIndex502, depth502
													if buffer[position] != rune('S') {
														goto l486
													}
													position++
												}
											l502:
												{
													position504, tokenIndex504, depth504 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l505
													}
													position++
													goto l504
												l505:
													position, tokenIndex, depth = position504, tokenIndex504, depth504
													if buffer[position] != rune('C') {
														goto l486
													}
													position++
												}
											l504:
												if !rules[ruleskip]() {
													goto l486
												}
												depth--
												add(ruleDESC, position497)
											}
										}
									l488:
										goto l487
									l486:
										position, tokenIndex, depth = position486, tokenIndex486, depth486
									}
								l487:
									if !rules[rulebrackettedExpression]() {
										goto l485
									}
									goto l484
								l485:
									position, tokenIndex, depth = position484, tokenIndex484, depth484
									if !rules[rulefunctionCall]() {
										goto l506
									}
									goto l484
								l506:
									position, tokenIndex, depth = position484, tokenIndex484, depth484
									if !rules[rulebuiltinCall]() {
										goto l507
									}
									goto l484
								l507:
									position, tokenIndex, depth = position484, tokenIndex484, depth484
									if !rules[rulevar]() {
										goto l457
									}
								}
							l484:
								depth--
								add(ruleorderCondition, position483)
							}
							goto l456
						l457:
							position, tokenIndex, depth = position457, tokenIndex457, depth457
						}
						goto l443
					l444:
						position, tokenIndex, depth = position443, tokenIndex443, depth443
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position509 := position
									depth++
									{
										position510, tokenIndex510, depth510 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l511
										}
										position++
										goto l510
									l511:
										position, tokenIndex, depth = position510, tokenIndex510, depth510
										if buffer[position] != rune('H') {
											goto l441
										}
										position++
									}
								l510:
									{
										position512, tokenIndex512, depth512 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l513
										}
										position++
										goto l512
									l513:
										position, tokenIndex, depth = position512, tokenIndex512, depth512
										if buffer[position] != rune('A') {
											goto l441
										}
										position++
									}
								l512:
									{
										position514, tokenIndex514, depth514 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l515
										}
										position++
										goto l514
									l515:
										position, tokenIndex, depth = position514, tokenIndex514, depth514
										if buffer[position] != rune('V') {
											goto l441
										}
										position++
									}
								l514:
									{
										position516, tokenIndex516, depth516 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l517
										}
										position++
										goto l516
									l517:
										position, tokenIndex, depth = position516, tokenIndex516, depth516
										if buffer[position] != rune('I') {
											goto l441
										}
										position++
									}
								l516:
									{
										position518, tokenIndex518, depth518 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l519
										}
										position++
										goto l518
									l519:
										position, tokenIndex, depth = position518, tokenIndex518, depth518
										if buffer[position] != rune('N') {
											goto l441
										}
										position++
									}
								l518:
									{
										position520, tokenIndex520, depth520 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l521
										}
										position++
										goto l520
									l521:
										position, tokenIndex, depth = position520, tokenIndex520, depth520
										if buffer[position] != rune('G') {
											goto l441
										}
										position++
									}
								l520:
									if !rules[ruleskip]() {
										goto l441
									}
									depth--
									add(ruleHAVING, position509)
								}
								if !rules[ruleconstraint]() {
									goto l441
								}
								break
							case 'G', 'g':
								{
									position522 := position
									depth++
									{
										position523, tokenIndex523, depth523 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l524
										}
										position++
										goto l523
									l524:
										position, tokenIndex, depth = position523, tokenIndex523, depth523
										if buffer[position] != rune('G') {
											goto l441
										}
										position++
									}
								l523:
									{
										position525, tokenIndex525, depth525 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l526
										}
										position++
										goto l525
									l526:
										position, tokenIndex, depth = position525, tokenIndex525, depth525
										if buffer[position] != rune('R') {
											goto l441
										}
										position++
									}
								l525:
									{
										position527, tokenIndex527, depth527 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l528
										}
										position++
										goto l527
									l528:
										position, tokenIndex, depth = position527, tokenIndex527, depth527
										if buffer[position] != rune('O') {
											goto l441
										}
										position++
									}
								l527:
									{
										position529, tokenIndex529, depth529 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l530
										}
										position++
										goto l529
									l530:
										position, tokenIndex, depth = position529, tokenIndex529, depth529
										if buffer[position] != rune('U') {
											goto l441
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
											goto l441
										}
										position++
									}
								l531:
									if !rules[ruleskip]() {
										goto l441
									}
									depth--
									add(ruleGROUP, position522)
								}
								if !rules[ruleBY]() {
									goto l441
								}
								{
									position535 := position
									depth++
									{
										position536, tokenIndex536, depth536 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l537
										}
										goto l536
									l537:
										position, tokenIndex, depth = position536, tokenIndex536, depth536
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l441
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l441
												}
												if !rules[ruleexpression]() {
													goto l441
												}
												{
													position539, tokenIndex539, depth539 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l539
													}
													if !rules[rulevar]() {
														goto l539
													}
													goto l540
												l539:
													position, tokenIndex, depth = position539, tokenIndex539, depth539
												}
											l540:
												if !rules[ruleRPAREN]() {
													goto l441
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l441
												}
												break
											}
										}

									}
								l536:
									depth--
									add(rulegroupCondition, position535)
								}
							l533:
								{
									position534, tokenIndex534, depth534 := position, tokenIndex, depth
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
														goto l534
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l534
													}
													if !rules[ruleexpression]() {
														goto l534
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
														goto l534
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l534
													}
													break
												}
											}

										}
									l542:
										depth--
										add(rulegroupCondition, position541)
									}
									goto l533
								l534:
									position, tokenIndex, depth = position534, tokenIndex534, depth534
								}
								break
							default:
								{
									position547 := position
									depth++
									{
										position548, tokenIndex548, depth548 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l549
										}
										{
											position550, tokenIndex550, depth550 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l550
											}
											goto l551
										l550:
											position, tokenIndex, depth = position550, tokenIndex550, depth550
										}
									l551:
										goto l548
									l549:
										position, tokenIndex, depth = position548, tokenIndex548, depth548
										if !rules[ruleoffset]() {
											goto l441
										}
										{
											position552, tokenIndex552, depth552 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l552
											}
											goto l553
										l552:
											position, tokenIndex, depth = position552, tokenIndex552, depth552
										}
									l553:
									}
								l548:
									depth--
									add(rulelimitOffsetClauses, position547)
								}
								break
							}
						}

					}
				l443:
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
		/* 47 groupCondition <- <(functionCall / ((&('$' | '?') var) | (&('(') (LPAREN expression (AS var)? RPAREN)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'H' | 'I' | 'L' | 'M' | 'N' | 'R' | 'S' | 'T' | 'U' | 'Y' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'h' | 'i' | 'l' | 'm' | 'n' | 'r' | 's' | 't' | 'u' | 'y') builtinCall)))> */
		nil,
		/* 48 orderCondition <- <(((ASC / DESC)? brackettedExpression) / functionCall / builtinCall / var)> */
		nil,
		/* 49 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 50 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position557, tokenIndex557, depth557 := position, tokenIndex, depth
			{
				position558 := position
				depth++
				{
					position559 := position
					depth++
					{
						position560, tokenIndex560, depth560 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l561
						}
						position++
						goto l560
					l561:
						position, tokenIndex, depth = position560, tokenIndex560, depth560
						if buffer[position] != rune('L') {
							goto l557
						}
						position++
					}
				l560:
					{
						position562, tokenIndex562, depth562 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l563
						}
						position++
						goto l562
					l563:
						position, tokenIndex, depth = position562, tokenIndex562, depth562
						if buffer[position] != rune('I') {
							goto l557
						}
						position++
					}
				l562:
					{
						position564, tokenIndex564, depth564 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l565
						}
						position++
						goto l564
					l565:
						position, tokenIndex, depth = position564, tokenIndex564, depth564
						if buffer[position] != rune('M') {
							goto l557
						}
						position++
					}
				l564:
					{
						position566, tokenIndex566, depth566 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l567
						}
						position++
						goto l566
					l567:
						position, tokenIndex, depth = position566, tokenIndex566, depth566
						if buffer[position] != rune('I') {
							goto l557
						}
						position++
					}
				l566:
					{
						position568, tokenIndex568, depth568 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l569
						}
						position++
						goto l568
					l569:
						position, tokenIndex, depth = position568, tokenIndex568, depth568
						if buffer[position] != rune('T') {
							goto l557
						}
						position++
					}
				l568:
					if !rules[ruleskip]() {
						goto l557
					}
					depth--
					add(ruleLIMIT, position559)
				}
				if !rules[ruleINTEGER]() {
					goto l557
				}
				depth--
				add(rulelimit, position558)
			}
			return true
		l557:
			position, tokenIndex, depth = position557, tokenIndex557, depth557
			return false
		},
		/* 51 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position570, tokenIndex570, depth570 := position, tokenIndex, depth
			{
				position571 := position
				depth++
				{
					position572 := position
					depth++
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l574
						}
						position++
						goto l573
					l574:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
						if buffer[position] != rune('O') {
							goto l570
						}
						position++
					}
				l573:
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l576
						}
						position++
						goto l575
					l576:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
						if buffer[position] != rune('F') {
							goto l570
						}
						position++
					}
				l575:
					{
						position577, tokenIndex577, depth577 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l578
						}
						position++
						goto l577
					l578:
						position, tokenIndex, depth = position577, tokenIndex577, depth577
						if buffer[position] != rune('F') {
							goto l570
						}
						position++
					}
				l577:
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
							goto l570
						}
						position++
					}
				l579:
					{
						position581, tokenIndex581, depth581 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l582
						}
						position++
						goto l581
					l582:
						position, tokenIndex, depth = position581, tokenIndex581, depth581
						if buffer[position] != rune('E') {
							goto l570
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
							goto l570
						}
						position++
					}
				l583:
					if !rules[ruleskip]() {
						goto l570
					}
					depth--
					add(ruleOFFSET, position572)
				}
				if !rules[ruleINTEGER]() {
					goto l570
				}
				depth--
				add(ruleoffset, position571)
			}
			return true
		l570:
			position, tokenIndex, depth = position570, tokenIndex570, depth570
			return false
		},
		/* 52 expression <- <conditionalOrExpression> */
		func() bool {
			position585, tokenIndex585, depth585 := position, tokenIndex, depth
			{
				position586 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l585
				}
				depth--
				add(ruleexpression, position586)
			}
			return true
		l585:
			position, tokenIndex, depth = position585, tokenIndex585, depth585
			return false
		},
		/* 53 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position587, tokenIndex587, depth587 := position, tokenIndex, depth
			{
				position588 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l587
				}
				{
					position589, tokenIndex589, depth589 := position, tokenIndex, depth
					{
						position591 := position
						depth++
						if buffer[position] != rune('|') {
							goto l589
						}
						position++
						if buffer[position] != rune('|') {
							goto l589
						}
						position++
						if !rules[ruleskip]() {
							goto l589
						}
						depth--
						add(ruleOR, position591)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l589
					}
					goto l590
				l589:
					position, tokenIndex, depth = position589, tokenIndex589, depth589
				}
			l590:
				depth--
				add(ruleconditionalOrExpression, position588)
			}
			return true
		l587:
			position, tokenIndex, depth = position587, tokenIndex587, depth587
			return false
		},
		/* 54 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position592, tokenIndex592, depth592 := position, tokenIndex, depth
			{
				position593 := position
				depth++
				{
					position594 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l592
					}
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position598 := position
									depth++
									{
										position599 := position
										depth++
										{
											position600, tokenIndex600, depth600 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l601
											}
											position++
											goto l600
										l601:
											position, tokenIndex, depth = position600, tokenIndex600, depth600
											if buffer[position] != rune('N') {
												goto l595
											}
											position++
										}
									l600:
										{
											position602, tokenIndex602, depth602 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l603
											}
											position++
											goto l602
										l603:
											position, tokenIndex, depth = position602, tokenIndex602, depth602
											if buffer[position] != rune('O') {
												goto l595
											}
											position++
										}
									l602:
										{
											position604, tokenIndex604, depth604 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l605
											}
											position++
											goto l604
										l605:
											position, tokenIndex, depth = position604, tokenIndex604, depth604
											if buffer[position] != rune('T') {
												goto l595
											}
											position++
										}
									l604:
										if buffer[position] != rune(' ') {
											goto l595
										}
										position++
										{
											position606, tokenIndex606, depth606 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l607
											}
											position++
											goto l606
										l607:
											position, tokenIndex, depth = position606, tokenIndex606, depth606
											if buffer[position] != rune('I') {
												goto l595
											}
											position++
										}
									l606:
										{
											position608, tokenIndex608, depth608 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l609
											}
											position++
											goto l608
										l609:
											position, tokenIndex, depth = position608, tokenIndex608, depth608
											if buffer[position] != rune('N') {
												goto l595
											}
											position++
										}
									l608:
										if !rules[ruleskip]() {
											goto l595
										}
										depth--
										add(ruleNOTIN, position599)
									}
									if !rules[ruleargList]() {
										goto l595
									}
									depth--
									add(rulenotin, position598)
								}
								break
							case 'I', 'i':
								{
									position610 := position
									depth++
									{
										position611 := position
										depth++
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
												goto l595
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
												goto l595
											}
											position++
										}
									l614:
										if !rules[ruleskip]() {
											goto l595
										}
										depth--
										add(ruleIN, position611)
									}
									if !rules[ruleargList]() {
										goto l595
									}
									depth--
									add(rulein, position610)
								}
								break
							default:
								{
									position616, tokenIndex616, depth616 := position, tokenIndex, depth
									{
										position618 := position
										depth++
										if buffer[position] != rune('<') {
											goto l617
										}
										position++
										if !rules[ruleskip]() {
											goto l617
										}
										depth--
										add(ruleLT, position618)
									}
									goto l616
								l617:
									position, tokenIndex, depth = position616, tokenIndex616, depth616
									{
										position620 := position
										depth++
										if buffer[position] != rune('>') {
											goto l619
										}
										position++
										if buffer[position] != rune('=') {
											goto l619
										}
										position++
										if !rules[ruleskip]() {
											goto l619
										}
										depth--
										add(ruleGE, position620)
									}
									goto l616
								l619:
									position, tokenIndex, depth = position616, tokenIndex616, depth616
									{
										switch buffer[position] {
										case '>':
											{
												position622 := position
												depth++
												if buffer[position] != rune('>') {
													goto l595
												}
												position++
												if !rules[ruleskip]() {
													goto l595
												}
												depth--
												add(ruleGT, position622)
											}
											break
										case '<':
											{
												position623 := position
												depth++
												if buffer[position] != rune('<') {
													goto l595
												}
												position++
												if buffer[position] != rune('=') {
													goto l595
												}
												position++
												if !rules[ruleskip]() {
													goto l595
												}
												depth--
												add(ruleLE, position623)
											}
											break
										case '!':
											{
												position624 := position
												depth++
												if buffer[position] != rune('!') {
													goto l595
												}
												position++
												if buffer[position] != rune('=') {
													goto l595
												}
												position++
												if !rules[ruleskip]() {
													goto l595
												}
												depth--
												add(ruleNE, position624)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l595
											}
											break
										}
									}

								}
							l616:
								if !rules[rulenumericExpression]() {
									goto l595
								}
								break
							}
						}

						goto l596
					l595:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
					}
				l596:
					depth--
					add(rulevalueLogical, position594)
				}
				{
					position625, tokenIndex625, depth625 := position, tokenIndex, depth
					{
						position627 := position
						depth++
						if buffer[position] != rune('&') {
							goto l625
						}
						position++
						if buffer[position] != rune('&') {
							goto l625
						}
						position++
						if !rules[ruleskip]() {
							goto l625
						}
						depth--
						add(ruleAND, position627)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l625
					}
					goto l626
				l625:
					position, tokenIndex, depth = position625, tokenIndex625, depth625
				}
			l626:
				depth--
				add(ruleconditionalAndExpression, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 55 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 56 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position629, tokenIndex629, depth629 := position, tokenIndex, depth
			{
				position630 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l629
				}
			l631:
				{
					position632, tokenIndex632, depth632 := position, tokenIndex, depth
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						{
							position635, tokenIndex635, depth635 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l636
							}
							goto l635
						l636:
							position, tokenIndex, depth = position635, tokenIndex635, depth635
							if !rules[ruleMINUS]() {
								goto l634
							}
						}
					l635:
						if !rules[rulemultiplicativeExpression]() {
							goto l634
						}
						goto l633
					l634:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
						{
							position637 := position
							depth++
							{
								position638, tokenIndex638, depth638 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l639
								}
								position++
								goto l638
							l639:
								position, tokenIndex, depth = position638, tokenIndex638, depth638
								if buffer[position] != rune('-') {
									goto l632
								}
								position++
							}
						l638:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l632
							}
							position++
						l640:
							{
								position641, tokenIndex641, depth641 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l641
								}
								position++
								goto l640
							l641:
								position, tokenIndex, depth = position641, tokenIndex641, depth641
							}
							{
								position642, tokenIndex642, depth642 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l642
								}
								position++
							l644:
								{
									position645, tokenIndex645, depth645 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l645
									}
									position++
									goto l644
								l645:
									position, tokenIndex, depth = position645, tokenIndex645, depth645
								}
								goto l643
							l642:
								position, tokenIndex, depth = position642, tokenIndex642, depth642
							}
						l643:
							if !rules[ruleskip]() {
								goto l632
							}
							depth--
							add(rulesignedNumericLiteral, position637)
						}
					}
				l633:
					goto l631
				l632:
					position, tokenIndex, depth = position632, tokenIndex632, depth632
				}
				depth--
				add(rulenumericExpression, position630)
			}
			return true
		l629:
			position, tokenIndex, depth = position629, tokenIndex629, depth629
			return false
		},
		/* 57 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position646, tokenIndex646, depth646 := position, tokenIndex, depth
			{
				position647 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l646
				}
			l648:
				{
					position649, tokenIndex649, depth649 := position, tokenIndex, depth
					{
						position650, tokenIndex650, depth650 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l651
						}
						goto l650
					l651:
						position, tokenIndex, depth = position650, tokenIndex650, depth650
						if !rules[ruleSLASH]() {
							goto l649
						}
					}
				l650:
					if !rules[ruleunaryExpression]() {
						goto l649
					}
					goto l648
				l649:
					position, tokenIndex, depth = position649, tokenIndex649, depth649
				}
				depth--
				add(rulemultiplicativeExpression, position647)
			}
			return true
		l646:
			position, tokenIndex, depth = position646, tokenIndex646, depth646
			return false
		},
		/* 58 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position652, tokenIndex652, depth652 := position, tokenIndex, depth
			{
				position653 := position
				depth++
				{
					position654, tokenIndex654, depth654 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l654
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l654
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l654
							}
							break
						}
					}

					goto l655
				l654:
					position, tokenIndex, depth = position654, tokenIndex654, depth654
				}
			l655:
				{
					position657 := position
					depth++
					{
						position658, tokenIndex658, depth658 := position, tokenIndex, depth
						if !rules[rulebuiltinCall]() {
							goto l659
						}
						goto l658
					l659:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
						if !rules[rulefunctionCall]() {
							goto l660
						}
						goto l658
					l660:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
						if !rules[ruleiriref]() {
							goto l661
						}
						goto l658
					l661:
						position, tokenIndex, depth = position658, tokenIndex658, depth658
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position663 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position665 := position
												depth++
												{
													position666 := position
													depth++
													{
														position667, tokenIndex667, depth667 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l668
														}
														position++
														goto l667
													l668:
														position, tokenIndex, depth = position667, tokenIndex667, depth667
														if buffer[position] != rune('G') {
															goto l652
														}
														position++
													}
												l667:
													{
														position669, tokenIndex669, depth669 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l670
														}
														position++
														goto l669
													l670:
														position, tokenIndex, depth = position669, tokenIndex669, depth669
														if buffer[position] != rune('R') {
															goto l652
														}
														position++
													}
												l669:
													{
														position671, tokenIndex671, depth671 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l672
														}
														position++
														goto l671
													l672:
														position, tokenIndex, depth = position671, tokenIndex671, depth671
														if buffer[position] != rune('O') {
															goto l652
														}
														position++
													}
												l671:
													{
														position673, tokenIndex673, depth673 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l674
														}
														position++
														goto l673
													l674:
														position, tokenIndex, depth = position673, tokenIndex673, depth673
														if buffer[position] != rune('U') {
															goto l652
														}
														position++
													}
												l673:
													{
														position675, tokenIndex675, depth675 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l676
														}
														position++
														goto l675
													l676:
														position, tokenIndex, depth = position675, tokenIndex675, depth675
														if buffer[position] != rune('P') {
															goto l652
														}
														position++
													}
												l675:
													if buffer[position] != rune('_') {
														goto l652
													}
													position++
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
															goto l652
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
															goto l652
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
															goto l652
														}
														position++
													}
												l681:
													{
														position683, tokenIndex683, depth683 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l684
														}
														position++
														goto l683
													l684:
														position, tokenIndex, depth = position683, tokenIndex683, depth683
														if buffer[position] != rune('C') {
															goto l652
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
															goto l652
														}
														position++
													}
												l685:
													{
														position687, tokenIndex687, depth687 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l688
														}
														position++
														goto l687
													l688:
														position, tokenIndex, depth = position687, tokenIndex687, depth687
														if buffer[position] != rune('T') {
															goto l652
														}
														position++
													}
												l687:
													if !rules[ruleskip]() {
														goto l652
													}
													depth--
													add(ruleGROUPCONCAT, position666)
												}
												if !rules[ruleLPAREN]() {
													goto l652
												}
												{
													position689, tokenIndex689, depth689 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l689
													}
													goto l690
												l689:
													position, tokenIndex, depth = position689, tokenIndex689, depth689
												}
											l690:
												if !rules[ruleexpression]() {
													goto l652
												}
												{
													position691, tokenIndex691, depth691 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l691
													}
													{
														position693 := position
														depth++
														{
															position694, tokenIndex694, depth694 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l695
															}
															position++
															goto l694
														l695:
															position, tokenIndex, depth = position694, tokenIndex694, depth694
															if buffer[position] != rune('S') {
																goto l691
															}
															position++
														}
													l694:
														{
															position696, tokenIndex696, depth696 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l697
															}
															position++
															goto l696
														l697:
															position, tokenIndex, depth = position696, tokenIndex696, depth696
															if buffer[position] != rune('E') {
																goto l691
															}
															position++
														}
													l696:
														{
															position698, tokenIndex698, depth698 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l699
															}
															position++
															goto l698
														l699:
															position, tokenIndex, depth = position698, tokenIndex698, depth698
															if buffer[position] != rune('P') {
																goto l691
															}
															position++
														}
													l698:
														{
															position700, tokenIndex700, depth700 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l701
															}
															position++
															goto l700
														l701:
															position, tokenIndex, depth = position700, tokenIndex700, depth700
															if buffer[position] != rune('A') {
																goto l691
															}
															position++
														}
													l700:
														{
															position702, tokenIndex702, depth702 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l703
															}
															position++
															goto l702
														l703:
															position, tokenIndex, depth = position702, tokenIndex702, depth702
															if buffer[position] != rune('R') {
																goto l691
															}
															position++
														}
													l702:
														{
															position704, tokenIndex704, depth704 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l705
															}
															position++
															goto l704
														l705:
															position, tokenIndex, depth = position704, tokenIndex704, depth704
															if buffer[position] != rune('A') {
																goto l691
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
																goto l691
															}
															position++
														}
													l706:
														{
															position708, tokenIndex708, depth708 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l709
															}
															position++
															goto l708
														l709:
															position, tokenIndex, depth = position708, tokenIndex708, depth708
															if buffer[position] != rune('O') {
																goto l691
															}
															position++
														}
													l708:
														{
															position710, tokenIndex710, depth710 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l711
															}
															position++
															goto l710
														l711:
															position, tokenIndex, depth = position710, tokenIndex710, depth710
															if buffer[position] != rune('R') {
																goto l691
															}
															position++
														}
													l710:
														if !rules[ruleskip]() {
															goto l691
														}
														depth--
														add(ruleSEPARATOR, position693)
													}
													if !rules[ruleEQ]() {
														goto l691
													}
													if !rules[rulestring]() {
														goto l691
													}
													goto l692
												l691:
													position, tokenIndex, depth = position691, tokenIndex691, depth691
												}
											l692:
												if !rules[ruleRPAREN]() {
													goto l652
												}
												depth--
												add(rulegroupConcat, position665)
											}
											break
										case 'C', 'c':
											{
												position712 := position
												depth++
												{
													position713 := position
													depth++
													{
														position714, tokenIndex714, depth714 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l715
														}
														position++
														goto l714
													l715:
														position, tokenIndex, depth = position714, tokenIndex714, depth714
														if buffer[position] != rune('C') {
															goto l652
														}
														position++
													}
												l714:
													{
														position716, tokenIndex716, depth716 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l717
														}
														position++
														goto l716
													l717:
														position, tokenIndex, depth = position716, tokenIndex716, depth716
														if buffer[position] != rune('O') {
															goto l652
														}
														position++
													}
												l716:
													{
														position718, tokenIndex718, depth718 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l719
														}
														position++
														goto l718
													l719:
														position, tokenIndex, depth = position718, tokenIndex718, depth718
														if buffer[position] != rune('U') {
															goto l652
														}
														position++
													}
												l718:
													{
														position720, tokenIndex720, depth720 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l721
														}
														position++
														goto l720
													l721:
														position, tokenIndex, depth = position720, tokenIndex720, depth720
														if buffer[position] != rune('N') {
															goto l652
														}
														position++
													}
												l720:
													{
														position722, tokenIndex722, depth722 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l723
														}
														position++
														goto l722
													l723:
														position, tokenIndex, depth = position722, tokenIndex722, depth722
														if buffer[position] != rune('T') {
															goto l652
														}
														position++
													}
												l722:
													if !rules[ruleskip]() {
														goto l652
													}
													depth--
													add(ruleCOUNT, position713)
												}
												if !rules[ruleLPAREN]() {
													goto l652
												}
												{
													position724, tokenIndex724, depth724 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l724
													}
													goto l725
												l724:
													position, tokenIndex, depth = position724, tokenIndex724, depth724
												}
											l725:
												{
													position726, tokenIndex726, depth726 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l727
													}
													goto l726
												l727:
													position, tokenIndex, depth = position726, tokenIndex726, depth726
													if !rules[ruleexpression]() {
														goto l652
													}
												}
											l726:
												if !rules[ruleRPAREN]() {
													goto l652
												}
												depth--
												add(rulecount, position712)
											}
											break
										default:
											{
												position728, tokenIndex728, depth728 := position, tokenIndex, depth
												{
													position730 := position
													depth++
													{
														position731, tokenIndex731, depth731 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l732
														}
														position++
														goto l731
													l732:
														position, tokenIndex, depth = position731, tokenIndex731, depth731
														if buffer[position] != rune('S') {
															goto l729
														}
														position++
													}
												l731:
													{
														position733, tokenIndex733, depth733 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l734
														}
														position++
														goto l733
													l734:
														position, tokenIndex, depth = position733, tokenIndex733, depth733
														if buffer[position] != rune('U') {
															goto l729
														}
														position++
													}
												l733:
													{
														position735, tokenIndex735, depth735 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l736
														}
														position++
														goto l735
													l736:
														position, tokenIndex, depth = position735, tokenIndex735, depth735
														if buffer[position] != rune('M') {
															goto l729
														}
														position++
													}
												l735:
													if !rules[ruleskip]() {
														goto l729
													}
													depth--
													add(ruleSUM, position730)
												}
												goto l728
											l729:
												position, tokenIndex, depth = position728, tokenIndex728, depth728
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
															goto l737
														}
														position++
													}
												l739:
													{
														position741, tokenIndex741, depth741 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l742
														}
														position++
														goto l741
													l742:
														position, tokenIndex, depth = position741, tokenIndex741, depth741
														if buffer[position] != rune('I') {
															goto l737
														}
														position++
													}
												l741:
													{
														position743, tokenIndex743, depth743 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l744
														}
														position++
														goto l743
													l744:
														position, tokenIndex, depth = position743, tokenIndex743, depth743
														if buffer[position] != rune('N') {
															goto l737
														}
														position++
													}
												l743:
													if !rules[ruleskip]() {
														goto l737
													}
													depth--
													add(ruleMIN, position738)
												}
												goto l728
											l737:
												position, tokenIndex, depth = position728, tokenIndex728, depth728
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position746 := position
															depth++
															{
																position747, tokenIndex747, depth747 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l748
																}
																position++
																goto l747
															l748:
																position, tokenIndex, depth = position747, tokenIndex747, depth747
																if buffer[position] != rune('S') {
																	goto l652
																}
																position++
															}
														l747:
															{
																position749, tokenIndex749, depth749 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l750
																}
																position++
																goto l749
															l750:
																position, tokenIndex, depth = position749, tokenIndex749, depth749
																if buffer[position] != rune('A') {
																	goto l652
																}
																position++
															}
														l749:
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
																	goto l652
																}
																position++
															}
														l751:
															{
																position753, tokenIndex753, depth753 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l754
																}
																position++
																goto l753
															l754:
																position, tokenIndex, depth = position753, tokenIndex753, depth753
																if buffer[position] != rune('P') {
																	goto l652
																}
																position++
															}
														l753:
															{
																position755, tokenIndex755, depth755 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l756
																}
																position++
																goto l755
															l756:
																position, tokenIndex, depth = position755, tokenIndex755, depth755
																if buffer[position] != rune('L') {
																	goto l652
																}
																position++
															}
														l755:
															{
																position757, tokenIndex757, depth757 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l758
																}
																position++
																goto l757
															l758:
																position, tokenIndex, depth = position757, tokenIndex757, depth757
																if buffer[position] != rune('E') {
																	goto l652
																}
																position++
															}
														l757:
															if !rules[ruleskip]() {
																goto l652
															}
															depth--
															add(ruleSAMPLE, position746)
														}
														break
													case 'A', 'a':
														{
															position759 := position
															depth++
															{
																position760, tokenIndex760, depth760 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l761
																}
																position++
																goto l760
															l761:
																position, tokenIndex, depth = position760, tokenIndex760, depth760
																if buffer[position] != rune('A') {
																	goto l652
																}
																position++
															}
														l760:
															{
																position762, tokenIndex762, depth762 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l763
																}
																position++
																goto l762
															l763:
																position, tokenIndex, depth = position762, tokenIndex762, depth762
																if buffer[position] != rune('V') {
																	goto l652
																}
																position++
															}
														l762:
															{
																position764, tokenIndex764, depth764 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l765
																}
																position++
																goto l764
															l765:
																position, tokenIndex, depth = position764, tokenIndex764, depth764
																if buffer[position] != rune('G') {
																	goto l652
																}
																position++
															}
														l764:
															if !rules[ruleskip]() {
																goto l652
															}
															depth--
															add(ruleAVG, position759)
														}
														break
													default:
														{
															position766 := position
															depth++
															{
																position767, tokenIndex767, depth767 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l768
																}
																position++
																goto l767
															l768:
																position, tokenIndex, depth = position767, tokenIndex767, depth767
																if buffer[position] != rune('M') {
																	goto l652
																}
																position++
															}
														l767:
															{
																position769, tokenIndex769, depth769 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l770
																}
																position++
																goto l769
															l770:
																position, tokenIndex, depth = position769, tokenIndex769, depth769
																if buffer[position] != rune('A') {
																	goto l652
																}
																position++
															}
														l769:
															{
																position771, tokenIndex771, depth771 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l772
																}
																position++
																goto l771
															l772:
																position, tokenIndex, depth = position771, tokenIndex771, depth771
																if buffer[position] != rune('X') {
																	goto l652
																}
																position++
															}
														l771:
															if !rules[ruleskip]() {
																goto l652
															}
															depth--
															add(ruleMAX, position766)
														}
														break
													}
												}

											}
										l728:
											if !rules[ruleLPAREN]() {
												goto l652
											}
											{
												position773, tokenIndex773, depth773 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l773
												}
												goto l774
											l773:
												position, tokenIndex, depth = position773, tokenIndex773, depth773
											}
										l774:
											if !rules[ruleexpression]() {
												goto l652
											}
											if !rules[ruleRPAREN]() {
												goto l652
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position663)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l652
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l652
								}
								break
							case '(':
								if !rules[rulebrackettedExpression]() {
									goto l652
								}
								break
							case '"', '\'':
								if !rules[ruleliteral]() {
									goto l652
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l652
								}
								break
							}
						}

					}
				l658:
					depth--
					add(ruleprimaryExpression, position657)
				}
				depth--
				add(ruleunaryExpression, position653)
			}
			return true
		l652:
			position, tokenIndex, depth = position652, tokenIndex652, depth652
			return false
		},
		/* 59 primaryExpression <- <(builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('(') brackettedExpression) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 60 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position776, tokenIndex776, depth776 := position, tokenIndex, depth
			{
				position777 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l776
				}
				if !rules[ruleexpression]() {
					goto l776
				}
				if !rules[ruleRPAREN]() {
					goto l776
				}
				depth--
				add(rulebrackettedExpression, position777)
			}
			return true
		l776:
			position, tokenIndex, depth = position776, tokenIndex776, depth776
			return false
		},
		/* 61 functionCall <- <(iriref argList)> */
		func() bool {
			position778, tokenIndex778, depth778 := position, tokenIndex, depth
			{
				position779 := position
				depth++
				if !rules[ruleiriref]() {
					goto l778
				}
				if !rules[ruleargList]() {
					goto l778
				}
				depth--
				add(rulefunctionCall, position779)
			}
			return true
		l778:
			position, tokenIndex, depth = position778, tokenIndex778, depth778
			return false
		},
		/* 62 in <- <(IN argList)> */
		nil,
		/* 63 notin <- <(NOTIN argList)> */
		nil,
		/* 64 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position782, tokenIndex782, depth782 := position, tokenIndex, depth
			{
				position783 := position
				depth++
				{
					position784, tokenIndex784, depth784 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l785
					}
					goto l784
				l785:
					position, tokenIndex, depth = position784, tokenIndex784, depth784
					if !rules[ruleLPAREN]() {
						goto l782
					}
					if !rules[ruleexpression]() {
						goto l782
					}
				l786:
					{
						position787, tokenIndex787, depth787 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l787
						}
						if !rules[ruleexpression]() {
							goto l787
						}
						goto l786
					l787:
						position, tokenIndex, depth = position787, tokenIndex787, depth787
					}
					if !rules[ruleRPAREN]() {
						goto l782
					}
				}
			l784:
				depth--
				add(ruleargList, position783)
			}
			return true
		l782:
			position, tokenIndex, depth = position782, tokenIndex782, depth782
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
			position791, tokenIndex791, depth791 := position, tokenIndex, depth
			{
				position792 := position
				depth++
				{
					position793, tokenIndex793, depth793 := position, tokenIndex, depth
					{
						position795, tokenIndex795, depth795 := position, tokenIndex, depth
						{
							position797 := position
							depth++
							{
								position798, tokenIndex798, depth798 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l799
								}
								position++
								goto l798
							l799:
								position, tokenIndex, depth = position798, tokenIndex798, depth798
								if buffer[position] != rune('S') {
									goto l796
								}
								position++
							}
						l798:
							{
								position800, tokenIndex800, depth800 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l801
								}
								position++
								goto l800
							l801:
								position, tokenIndex, depth = position800, tokenIndex800, depth800
								if buffer[position] != rune('T') {
									goto l796
								}
								position++
							}
						l800:
							{
								position802, tokenIndex802, depth802 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l803
								}
								position++
								goto l802
							l803:
								position, tokenIndex, depth = position802, tokenIndex802, depth802
								if buffer[position] != rune('R') {
									goto l796
								}
								position++
							}
						l802:
							if !rules[ruleskip]() {
								goto l796
							}
							depth--
							add(ruleSTR, position797)
						}
						goto l795
					l796:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position805 := position
							depth++
							{
								position806, tokenIndex806, depth806 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l807
								}
								position++
								goto l806
							l807:
								position, tokenIndex, depth = position806, tokenIndex806, depth806
								if buffer[position] != rune('L') {
									goto l804
								}
								position++
							}
						l806:
							{
								position808, tokenIndex808, depth808 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l809
								}
								position++
								goto l808
							l809:
								position, tokenIndex, depth = position808, tokenIndex808, depth808
								if buffer[position] != rune('A') {
									goto l804
								}
								position++
							}
						l808:
							{
								position810, tokenIndex810, depth810 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l811
								}
								position++
								goto l810
							l811:
								position, tokenIndex, depth = position810, tokenIndex810, depth810
								if buffer[position] != rune('N') {
									goto l804
								}
								position++
							}
						l810:
							{
								position812, tokenIndex812, depth812 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l813
								}
								position++
								goto l812
							l813:
								position, tokenIndex, depth = position812, tokenIndex812, depth812
								if buffer[position] != rune('G') {
									goto l804
								}
								position++
							}
						l812:
							if !rules[ruleskip]() {
								goto l804
							}
							depth--
							add(ruleLANG, position805)
						}
						goto l795
					l804:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position815 := position
							depth++
							{
								position816, tokenIndex816, depth816 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l817
								}
								position++
								goto l816
							l817:
								position, tokenIndex, depth = position816, tokenIndex816, depth816
								if buffer[position] != rune('D') {
									goto l814
								}
								position++
							}
						l816:
							{
								position818, tokenIndex818, depth818 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l819
								}
								position++
								goto l818
							l819:
								position, tokenIndex, depth = position818, tokenIndex818, depth818
								if buffer[position] != rune('A') {
									goto l814
								}
								position++
							}
						l818:
							{
								position820, tokenIndex820, depth820 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l821
								}
								position++
								goto l820
							l821:
								position, tokenIndex, depth = position820, tokenIndex820, depth820
								if buffer[position] != rune('T') {
									goto l814
								}
								position++
							}
						l820:
							{
								position822, tokenIndex822, depth822 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l823
								}
								position++
								goto l822
							l823:
								position, tokenIndex, depth = position822, tokenIndex822, depth822
								if buffer[position] != rune('A') {
									goto l814
								}
								position++
							}
						l822:
							{
								position824, tokenIndex824, depth824 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l825
								}
								position++
								goto l824
							l825:
								position, tokenIndex, depth = position824, tokenIndex824, depth824
								if buffer[position] != rune('T') {
									goto l814
								}
								position++
							}
						l824:
							{
								position826, tokenIndex826, depth826 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l827
								}
								position++
								goto l826
							l827:
								position, tokenIndex, depth = position826, tokenIndex826, depth826
								if buffer[position] != rune('Y') {
									goto l814
								}
								position++
							}
						l826:
							{
								position828, tokenIndex828, depth828 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l829
								}
								position++
								goto l828
							l829:
								position, tokenIndex, depth = position828, tokenIndex828, depth828
								if buffer[position] != rune('P') {
									goto l814
								}
								position++
							}
						l828:
							{
								position830, tokenIndex830, depth830 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l831
								}
								position++
								goto l830
							l831:
								position, tokenIndex, depth = position830, tokenIndex830, depth830
								if buffer[position] != rune('E') {
									goto l814
								}
								position++
							}
						l830:
							if !rules[ruleskip]() {
								goto l814
							}
							depth--
							add(ruleDATATYPE, position815)
						}
						goto l795
					l814:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position833 := position
							depth++
							{
								position834, tokenIndex834, depth834 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l835
								}
								position++
								goto l834
							l835:
								position, tokenIndex, depth = position834, tokenIndex834, depth834
								if buffer[position] != rune('I') {
									goto l832
								}
								position++
							}
						l834:
							{
								position836, tokenIndex836, depth836 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l837
								}
								position++
								goto l836
							l837:
								position, tokenIndex, depth = position836, tokenIndex836, depth836
								if buffer[position] != rune('R') {
									goto l832
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
									goto l832
								}
								position++
							}
						l838:
							if !rules[ruleskip]() {
								goto l832
							}
							depth--
							add(ruleIRI, position833)
						}
						goto l795
					l832:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position841 := position
							depth++
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
									goto l840
								}
								position++
							}
						l842:
							{
								position844, tokenIndex844, depth844 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l845
								}
								position++
								goto l844
							l845:
								position, tokenIndex, depth = position844, tokenIndex844, depth844
								if buffer[position] != rune('R') {
									goto l840
								}
								position++
							}
						l844:
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
									goto l840
								}
								position++
							}
						l846:
							if !rules[ruleskip]() {
								goto l840
							}
							depth--
							add(ruleURI, position841)
						}
						goto l795
					l840:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position849 := position
							depth++
							{
								position850, tokenIndex850, depth850 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l851
								}
								position++
								goto l850
							l851:
								position, tokenIndex, depth = position850, tokenIndex850, depth850
								if buffer[position] != rune('S') {
									goto l848
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
									goto l848
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
									goto l848
								}
								position++
							}
						l854:
							{
								position856, tokenIndex856, depth856 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l857
								}
								position++
								goto l856
							l857:
								position, tokenIndex, depth = position856, tokenIndex856, depth856
								if buffer[position] != rune('L') {
									goto l848
								}
								position++
							}
						l856:
							{
								position858, tokenIndex858, depth858 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('E') {
									goto l848
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
									goto l848
								}
								position++
							}
						l860:
							if !rules[ruleskip]() {
								goto l848
							}
							depth--
							add(ruleSTRLEN, position849)
						}
						goto l795
					l848:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position863 := position
							depth++
							{
								position864, tokenIndex864, depth864 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l865
								}
								position++
								goto l864
							l865:
								position, tokenIndex, depth = position864, tokenIndex864, depth864
								if buffer[position] != rune('M') {
									goto l862
								}
								position++
							}
						l864:
							{
								position866, tokenIndex866, depth866 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l867
								}
								position++
								goto l866
							l867:
								position, tokenIndex, depth = position866, tokenIndex866, depth866
								if buffer[position] != rune('O') {
									goto l862
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
									goto l862
								}
								position++
							}
						l868:
							{
								position870, tokenIndex870, depth870 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l871
								}
								position++
								goto l870
							l871:
								position, tokenIndex, depth = position870, tokenIndex870, depth870
								if buffer[position] != rune('T') {
									goto l862
								}
								position++
							}
						l870:
							{
								position872, tokenIndex872, depth872 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l873
								}
								position++
								goto l872
							l873:
								position, tokenIndex, depth = position872, tokenIndex872, depth872
								if buffer[position] != rune('H') {
									goto l862
								}
								position++
							}
						l872:
							if !rules[ruleskip]() {
								goto l862
							}
							depth--
							add(ruleMONTH, position863)
						}
						goto l795
					l862:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
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
								if buffer[position] != rune('i') {
									goto l879
								}
								position++
								goto l878
							l879:
								position, tokenIndex, depth = position878, tokenIndex878, depth878
								if buffer[position] != rune('I') {
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
								if buffer[position] != rune('u') {
									goto l883
								}
								position++
								goto l882
							l883:
								position, tokenIndex, depth = position882, tokenIndex882, depth882
								if buffer[position] != rune('U') {
									goto l874
								}
								position++
							}
						l882:
							{
								position884, tokenIndex884, depth884 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l885
								}
								position++
								goto l884
							l885:
								position, tokenIndex, depth = position884, tokenIndex884, depth884
								if buffer[position] != rune('T') {
									goto l874
								}
								position++
							}
						l884:
							{
								position886, tokenIndex886, depth886 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l887
								}
								position++
								goto l886
							l887:
								position, tokenIndex, depth = position886, tokenIndex886, depth886
								if buffer[position] != rune('E') {
									goto l874
								}
								position++
							}
						l886:
							{
								position888, tokenIndex888, depth888 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l889
								}
								position++
								goto l888
							l889:
								position, tokenIndex, depth = position888, tokenIndex888, depth888
								if buffer[position] != rune('S') {
									goto l874
								}
								position++
							}
						l888:
							if !rules[ruleskip]() {
								goto l874
							}
							depth--
							add(ruleMINUTES, position875)
						}
						goto l795
					l874:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position891 := position
							depth++
							{
								position892, tokenIndex892, depth892 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l893
								}
								position++
								goto l892
							l893:
								position, tokenIndex, depth = position892, tokenIndex892, depth892
								if buffer[position] != rune('S') {
									goto l890
								}
								position++
							}
						l892:
							{
								position894, tokenIndex894, depth894 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l895
								}
								position++
								goto l894
							l895:
								position, tokenIndex, depth = position894, tokenIndex894, depth894
								if buffer[position] != rune('E') {
									goto l890
								}
								position++
							}
						l894:
							{
								position896, tokenIndex896, depth896 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l897
								}
								position++
								goto l896
							l897:
								position, tokenIndex, depth = position896, tokenIndex896, depth896
								if buffer[position] != rune('C') {
									goto l890
								}
								position++
							}
						l896:
							{
								position898, tokenIndex898, depth898 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l899
								}
								position++
								goto l898
							l899:
								position, tokenIndex, depth = position898, tokenIndex898, depth898
								if buffer[position] != rune('O') {
									goto l890
								}
								position++
							}
						l898:
							{
								position900, tokenIndex900, depth900 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l901
								}
								position++
								goto l900
							l901:
								position, tokenIndex, depth = position900, tokenIndex900, depth900
								if buffer[position] != rune('N') {
									goto l890
								}
								position++
							}
						l900:
							{
								position902, tokenIndex902, depth902 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l903
								}
								position++
								goto l902
							l903:
								position, tokenIndex, depth = position902, tokenIndex902, depth902
								if buffer[position] != rune('D') {
									goto l890
								}
								position++
							}
						l902:
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
									goto l890
								}
								position++
							}
						l904:
							if !rules[ruleskip]() {
								goto l890
							}
							depth--
							add(ruleSECONDS, position891)
						}
						goto l795
					l890:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position907 := position
							depth++
							{
								position908, tokenIndex908, depth908 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l909
								}
								position++
								goto l908
							l909:
								position, tokenIndex, depth = position908, tokenIndex908, depth908
								if buffer[position] != rune('T') {
									goto l906
								}
								position++
							}
						l908:
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
									goto l906
								}
								position++
							}
						l910:
							{
								position912, tokenIndex912, depth912 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l913
								}
								position++
								goto l912
							l913:
								position, tokenIndex, depth = position912, tokenIndex912, depth912
								if buffer[position] != rune('M') {
									goto l906
								}
								position++
							}
						l912:
							{
								position914, tokenIndex914, depth914 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l915
								}
								position++
								goto l914
							l915:
								position, tokenIndex, depth = position914, tokenIndex914, depth914
								if buffer[position] != rune('E') {
									goto l906
								}
								position++
							}
						l914:
							{
								position916, tokenIndex916, depth916 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l917
								}
								position++
								goto l916
							l917:
								position, tokenIndex, depth = position916, tokenIndex916, depth916
								if buffer[position] != rune('Z') {
									goto l906
								}
								position++
							}
						l916:
							{
								position918, tokenIndex918, depth918 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l919
								}
								position++
								goto l918
							l919:
								position, tokenIndex, depth = position918, tokenIndex918, depth918
								if buffer[position] != rune('O') {
									goto l906
								}
								position++
							}
						l918:
							{
								position920, tokenIndex920, depth920 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l921
								}
								position++
								goto l920
							l921:
								position, tokenIndex, depth = position920, tokenIndex920, depth920
								if buffer[position] != rune('N') {
									goto l906
								}
								position++
							}
						l920:
							{
								position922, tokenIndex922, depth922 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l923
								}
								position++
								goto l922
							l923:
								position, tokenIndex, depth = position922, tokenIndex922, depth922
								if buffer[position] != rune('E') {
									goto l906
								}
								position++
							}
						l922:
							if !rules[ruleskip]() {
								goto l906
							}
							depth--
							add(ruleTIMEZONE, position907)
						}
						goto l795
					l906:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position925 := position
							depth++
							{
								position926, tokenIndex926, depth926 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l927
								}
								position++
								goto l926
							l927:
								position, tokenIndex, depth = position926, tokenIndex926, depth926
								if buffer[position] != rune('S') {
									goto l924
								}
								position++
							}
						l926:
							{
								position928, tokenIndex928, depth928 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l929
								}
								position++
								goto l928
							l929:
								position, tokenIndex, depth = position928, tokenIndex928, depth928
								if buffer[position] != rune('H') {
									goto l924
								}
								position++
							}
						l928:
							{
								position930, tokenIndex930, depth930 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l931
								}
								position++
								goto l930
							l931:
								position, tokenIndex, depth = position930, tokenIndex930, depth930
								if buffer[position] != rune('A') {
									goto l924
								}
								position++
							}
						l930:
							if buffer[position] != rune('1') {
								goto l924
							}
							position++
							if !rules[ruleskip]() {
								goto l924
							}
							depth--
							add(ruleSHA1, position925)
						}
						goto l795
					l924:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position933 := position
							depth++
							{
								position934, tokenIndex934, depth934 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l935
								}
								position++
								goto l934
							l935:
								position, tokenIndex, depth = position934, tokenIndex934, depth934
								if buffer[position] != rune('S') {
									goto l932
								}
								position++
							}
						l934:
							{
								position936, tokenIndex936, depth936 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l937
								}
								position++
								goto l936
							l937:
								position, tokenIndex, depth = position936, tokenIndex936, depth936
								if buffer[position] != rune('H') {
									goto l932
								}
								position++
							}
						l936:
							{
								position938, tokenIndex938, depth938 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l939
								}
								position++
								goto l938
							l939:
								position, tokenIndex, depth = position938, tokenIndex938, depth938
								if buffer[position] != rune('A') {
									goto l932
								}
								position++
							}
						l938:
							if buffer[position] != rune('2') {
								goto l932
							}
							position++
							if buffer[position] != rune('5') {
								goto l932
							}
							position++
							if buffer[position] != rune('6') {
								goto l932
							}
							position++
							if !rules[ruleskip]() {
								goto l932
							}
							depth--
							add(ruleSHA256, position933)
						}
						goto l795
					l932:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
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
								if buffer[position] != rune('h') {
									goto l945
								}
								position++
								goto l944
							l945:
								position, tokenIndex, depth = position944, tokenIndex944, depth944
								if buffer[position] != rune('H') {
									goto l940
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
									goto l940
								}
								position++
							}
						l946:
							if buffer[position] != rune('3') {
								goto l940
							}
							position++
							if buffer[position] != rune('8') {
								goto l940
							}
							position++
							if buffer[position] != rune('4') {
								goto l940
							}
							position++
							if !rules[ruleskip]() {
								goto l940
							}
							depth--
							add(ruleSHA384, position941)
						}
						goto l795
					l940:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
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
								if buffer[position] != rune('i') {
									goto l955
								}
								position++
								goto l954
							l955:
								position, tokenIndex, depth = position954, tokenIndex954, depth954
								if buffer[position] != rune('I') {
									goto l948
								}
								position++
							}
						l954:
							{
								position956, tokenIndex956, depth956 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l957
								}
								position++
								goto l956
							l957:
								position, tokenIndex, depth = position956, tokenIndex956, depth956
								if buffer[position] != rune('R') {
									goto l948
								}
								position++
							}
						l956:
							{
								position958, tokenIndex958, depth958 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l959
								}
								position++
								goto l958
							l959:
								position, tokenIndex, depth = position958, tokenIndex958, depth958
								if buffer[position] != rune('I') {
									goto l948
								}
								position++
							}
						l958:
							if !rules[ruleskip]() {
								goto l948
							}
							depth--
							add(ruleISIRI, position949)
						}
						goto l795
					l948:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
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
								if buffer[position] != rune('u') {
									goto l967
								}
								position++
								goto l966
							l967:
								position, tokenIndex, depth = position966, tokenIndex966, depth966
								if buffer[position] != rune('U') {
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
							add(ruleISURI, position961)
						}
						goto l795
					l960:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
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
								if buffer[position] != rune('b') {
									goto l979
								}
								position++
								goto l978
							l979:
								position, tokenIndex, depth = position978, tokenIndex978, depth978
								if buffer[position] != rune('B') {
									goto l972
								}
								position++
							}
						l978:
							{
								position980, tokenIndex980, depth980 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l981
								}
								position++
								goto l980
							l981:
								position, tokenIndex, depth = position980, tokenIndex980, depth980
								if buffer[position] != rune('L') {
									goto l972
								}
								position++
							}
						l980:
							{
								position982, tokenIndex982, depth982 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l983
								}
								position++
								goto l982
							l983:
								position, tokenIndex, depth = position982, tokenIndex982, depth982
								if buffer[position] != rune('A') {
									goto l972
								}
								position++
							}
						l982:
							{
								position984, tokenIndex984, depth984 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l985
								}
								position++
								goto l984
							l985:
								position, tokenIndex, depth = position984, tokenIndex984, depth984
								if buffer[position] != rune('N') {
									goto l972
								}
								position++
							}
						l984:
							{
								position986, tokenIndex986, depth986 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l987
								}
								position++
								goto l986
							l987:
								position, tokenIndex, depth = position986, tokenIndex986, depth986
								if buffer[position] != rune('K') {
									goto l972
								}
								position++
							}
						l986:
							if !rules[ruleskip]() {
								goto l972
							}
							depth--
							add(ruleISBLANK, position973)
						}
						goto l795
					l972:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							position989 := position
							depth++
							{
								position990, tokenIndex990, depth990 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l991
								}
								position++
								goto l990
							l991:
								position, tokenIndex, depth = position990, tokenIndex990, depth990
								if buffer[position] != rune('I') {
									goto l988
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
									goto l988
								}
								position++
							}
						l992:
							{
								position994, tokenIndex994, depth994 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l995
								}
								position++
								goto l994
							l995:
								position, tokenIndex, depth = position994, tokenIndex994, depth994
								if buffer[position] != rune('L') {
									goto l988
								}
								position++
							}
						l994:
							{
								position996, tokenIndex996, depth996 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l997
								}
								position++
								goto l996
							l997:
								position, tokenIndex, depth = position996, tokenIndex996, depth996
								if buffer[position] != rune('I') {
									goto l988
								}
								position++
							}
						l996:
							{
								position998, tokenIndex998, depth998 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l999
								}
								position++
								goto l998
							l999:
								position, tokenIndex, depth = position998, tokenIndex998, depth998
								if buffer[position] != rune('T') {
									goto l988
								}
								position++
							}
						l998:
							{
								position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1001
								}
								position++
								goto l1000
							l1001:
								position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
								if buffer[position] != rune('E') {
									goto l988
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
									goto l988
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
									goto l988
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
									goto l988
								}
								position++
							}
						l1006:
							if !rules[ruleskip]() {
								goto l988
							}
							depth--
							add(ruleISLITERAL, position989)
						}
						goto l795
					l988:
						position, tokenIndex, depth = position795, tokenIndex795, depth795
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position1009 := position
									depth++
									{
										position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1011
										}
										position++
										goto l1010
									l1011:
										position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
										if buffer[position] != rune('I') {
											goto l794
										}
										position++
									}
								l1010:
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
											goto l794
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
											goto l794
										}
										position++
									}
								l1014:
									{
										position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1017
										}
										position++
										goto l1016
									l1017:
										position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
										if buffer[position] != rune('U') {
											goto l794
										}
										position++
									}
								l1016:
									{
										position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1019
										}
										position++
										goto l1018
									l1019:
										position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
										if buffer[position] != rune('M') {
											goto l794
										}
										position++
									}
								l1018:
									{
										position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1021
										}
										position++
										goto l1020
									l1021:
										position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
										if buffer[position] != rune('E') {
											goto l794
										}
										position++
									}
								l1020:
									{
										position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1023
										}
										position++
										goto l1022
									l1023:
										position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
										if buffer[position] != rune('R') {
											goto l794
										}
										position++
									}
								l1022:
									{
										position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1025
										}
										position++
										goto l1024
									l1025:
										position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
										if buffer[position] != rune('I') {
											goto l794
										}
										position++
									}
								l1024:
									{
										position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1027
										}
										position++
										goto l1026
									l1027:
										position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
										if buffer[position] != rune('C') {
											goto l794
										}
										position++
									}
								l1026:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleISNUMERIC, position1009)
								}
								break
							case 'S', 's':
								{
									position1028 := position
									depth++
									{
										position1029, tokenIndex1029, depth1029 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1030
										}
										position++
										goto l1029
									l1030:
										position, tokenIndex, depth = position1029, tokenIndex1029, depth1029
										if buffer[position] != rune('S') {
											goto l794
										}
										position++
									}
								l1029:
									{
										position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1032
										}
										position++
										goto l1031
									l1032:
										position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
										if buffer[position] != rune('H') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1033:
									if buffer[position] != rune('5') {
										goto l794
									}
									position++
									if buffer[position] != rune('1') {
										goto l794
									}
									position++
									if buffer[position] != rune('2') {
										goto l794
									}
									position++
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleSHA512, position1028)
								}
								break
							case 'M', 'm':
								{
									position1035 := position
									depth++
									{
										position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1037
										}
										position++
										goto l1036
									l1037:
										position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
										if buffer[position] != rune('M') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1038:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleMD5, position1035)
								}
								break
							case 'T', 't':
								{
									position1040 := position
									depth++
									{
										position1041, tokenIndex1041, depth1041 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1042
										}
										position++
										goto l1041
									l1042:
										position, tokenIndex, depth = position1041, tokenIndex1041, depth1041
										if buffer[position] != rune('T') {
											goto l794
										}
										position++
									}
								l1041:
									{
										position1043, tokenIndex1043, depth1043 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1044
										}
										position++
										goto l1043
									l1044:
										position, tokenIndex, depth = position1043, tokenIndex1043, depth1043
										if buffer[position] != rune('Z') {
											goto l794
										}
										position++
									}
								l1043:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleTZ, position1040)
								}
								break
							case 'H', 'h':
								{
									position1045 := position
									depth++
									{
										position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1047
										}
										position++
										goto l1046
									l1047:
										position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
										if buffer[position] != rune('H') {
											goto l794
										}
										position++
									}
								l1046:
									{
										position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1049
										}
										position++
										goto l1048
									l1049:
										position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
										if buffer[position] != rune('O') {
											goto l794
										}
										position++
									}
								l1048:
									{
										position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1051
										}
										position++
										goto l1050
									l1051:
										position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
										if buffer[position] != rune('U') {
											goto l794
										}
										position++
									}
								l1050:
									{
										position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1053
										}
										position++
										goto l1052
									l1053:
										position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
										if buffer[position] != rune('R') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1054:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleHOURS, position1045)
								}
								break
							case 'D', 'd':
								{
									position1056 := position
									depth++
									{
										position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1058
										}
										position++
										goto l1057
									l1058:
										position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
										if buffer[position] != rune('D') {
											goto l794
										}
										position++
									}
								l1057:
									{
										position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1060
										}
										position++
										goto l1059
									l1060:
										position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
										if buffer[position] != rune('A') {
											goto l794
										}
										position++
									}
								l1059:
									{
										position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1062
										}
										position++
										goto l1061
									l1062:
										position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
										if buffer[position] != rune('Y') {
											goto l794
										}
										position++
									}
								l1061:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleDAY, position1056)
								}
								break
							case 'Y', 'y':
								{
									position1063 := position
									depth++
									{
										position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1065
										}
										position++
										goto l1064
									l1065:
										position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
										if buffer[position] != rune('Y') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1066:
									{
										position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1069
										}
										position++
										goto l1068
									l1069:
										position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
										if buffer[position] != rune('A') {
											goto l794
										}
										position++
									}
								l1068:
									{
										position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1071
										}
										position++
										goto l1070
									l1071:
										position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
										if buffer[position] != rune('R') {
											goto l794
										}
										position++
									}
								l1070:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleYEAR, position1063)
								}
								break
							case 'E', 'e':
								{
									position1072 := position
									depth++
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
											goto l794
										}
										position++
									}
								l1073:
									{
										position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1076
										}
										position++
										goto l1075
									l1076:
										position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
										if buffer[position] != rune('N') {
											goto l794
										}
										position++
									}
								l1075:
									{
										position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1078
										}
										position++
										goto l1077
									l1078:
										position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
										if buffer[position] != rune('C') {
											goto l794
										}
										position++
									}
								l1077:
									{
										position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1080
										}
										position++
										goto l1079
									l1080:
										position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
										if buffer[position] != rune('O') {
											goto l794
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
											goto l794
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
											goto l794
										}
										position++
									}
								l1083:
									if buffer[position] != rune('_') {
										goto l794
									}
									position++
									{
										position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1086
										}
										position++
										goto l1085
									l1086:
										position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
										if buffer[position] != rune('F') {
											goto l794
										}
										position++
									}
								l1085:
									{
										position1087, tokenIndex1087, depth1087 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1088
										}
										position++
										goto l1087
									l1088:
										position, tokenIndex, depth = position1087, tokenIndex1087, depth1087
										if buffer[position] != rune('O') {
											goto l794
										}
										position++
									}
								l1087:
									{
										position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1090
										}
										position++
										goto l1089
									l1090:
										position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
										if buffer[position] != rune('R') {
											goto l794
										}
										position++
									}
								l1089:
									if buffer[position] != rune('_') {
										goto l794
									}
									position++
									{
										position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1092
										}
										position++
										goto l1091
									l1092:
										position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
										if buffer[position] != rune('U') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1093:
									{
										position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1096
										}
										position++
										goto l1095
									l1096:
										position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
										if buffer[position] != rune('I') {
											goto l794
										}
										position++
									}
								l1095:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleENCODEFORURI, position1072)
								}
								break
							case 'L', 'l':
								{
									position1097 := position
									depth++
									{
										position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1099
										}
										position++
										goto l1098
									l1099:
										position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
										if buffer[position] != rune('L') {
											goto l794
										}
										position++
									}
								l1098:
									{
										position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1101
										}
										position++
										goto l1100
									l1101:
										position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
										if buffer[position] != rune('C') {
											goto l794
										}
										position++
									}
								l1100:
									{
										position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1103
										}
										position++
										goto l1102
									l1103:
										position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
										if buffer[position] != rune('A') {
											goto l794
										}
										position++
									}
								l1102:
									{
										position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1105
										}
										position++
										goto l1104
									l1105:
										position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
										if buffer[position] != rune('S') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1106:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleLCASE, position1097)
								}
								break
							case 'U', 'u':
								{
									position1108 := position
									depth++
									{
										position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1110
										}
										position++
										goto l1109
									l1110:
										position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
										if buffer[position] != rune('U') {
											goto l794
										}
										position++
									}
								l1109:
									{
										position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1112
										}
										position++
										goto l1111
									l1112:
										position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
										if buffer[position] != rune('C') {
											goto l794
										}
										position++
									}
								l1111:
									{
										position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1114
										}
										position++
										goto l1113
									l1114:
										position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
										if buffer[position] != rune('A') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1115:
									{
										position1117, tokenIndex1117, depth1117 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1118
										}
										position++
										goto l1117
									l1118:
										position, tokenIndex, depth = position1117, tokenIndex1117, depth1117
										if buffer[position] != rune('E') {
											goto l794
										}
										position++
									}
								l1117:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleUCASE, position1108)
								}
								break
							case 'F', 'f':
								{
									position1119 := position
									depth++
									{
										position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1121
										}
										position++
										goto l1120
									l1121:
										position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
										if buffer[position] != rune('F') {
											goto l794
										}
										position++
									}
								l1120:
									{
										position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1123
										}
										position++
										goto l1122
									l1123:
										position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
										if buffer[position] != rune('L') {
											goto l794
										}
										position++
									}
								l1122:
									{
										position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1125
										}
										position++
										goto l1124
									l1125:
										position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
										if buffer[position] != rune('O') {
											goto l794
										}
										position++
									}
								l1124:
									{
										position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1127
										}
										position++
										goto l1126
									l1127:
										position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
										if buffer[position] != rune('O') {
											goto l794
										}
										position++
									}
								l1126:
									{
										position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1129
										}
										position++
										goto l1128
									l1129:
										position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
										if buffer[position] != rune('R') {
											goto l794
										}
										position++
									}
								l1128:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleFLOOR, position1119)
								}
								break
							case 'R', 'r':
								{
									position1130 := position
									depth++
									{
										position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1132
										}
										position++
										goto l1131
									l1132:
										position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
										if buffer[position] != rune('R') {
											goto l794
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
											goto l794
										}
										position++
									}
								l1133:
									{
										position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1136
										}
										position++
										goto l1135
									l1136:
										position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
										if buffer[position] != rune('U') {
											goto l794
										}
										position++
									}
								l1135:
									{
										position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1138
										}
										position++
										goto l1137
									l1138:
										position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
										if buffer[position] != rune('N') {
											goto l794
										}
										position++
									}
								l1137:
									{
										position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1140
										}
										position++
										goto l1139
									l1140:
										position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
										if buffer[position] != rune('D') {
											goto l794
										}
										position++
									}
								l1139:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleROUND, position1130)
								}
								break
							case 'C', 'c':
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
											goto l794
										}
										position++
									}
								l1142:
									{
										position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1145
										}
										position++
										goto l1144
									l1145:
										position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
										if buffer[position] != rune('E') {
											goto l794
										}
										position++
									}
								l1144:
									{
										position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1147
										}
										position++
										goto l1146
									l1147:
										position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
										if buffer[position] != rune('I') {
											goto l794
										}
										position++
									}
								l1146:
									{
										position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1149
										}
										position++
										goto l1148
									l1149:
										position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
										if buffer[position] != rune('L') {
											goto l794
										}
										position++
									}
								l1148:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleCEIL, position1141)
								}
								break
							default:
								{
									position1150 := position
									depth++
									{
										position1151, tokenIndex1151, depth1151 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1152
										}
										position++
										goto l1151
									l1152:
										position, tokenIndex, depth = position1151, tokenIndex1151, depth1151
										if buffer[position] != rune('A') {
											goto l794
										}
										position++
									}
								l1151:
									{
										position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1154
										}
										position++
										goto l1153
									l1154:
										position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
										if buffer[position] != rune('B') {
											goto l794
										}
										position++
									}
								l1153:
									{
										position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1156
										}
										position++
										goto l1155
									l1156:
										position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
										if buffer[position] != rune('S') {
											goto l794
										}
										position++
									}
								l1155:
									if !rules[ruleskip]() {
										goto l794
									}
									depth--
									add(ruleABS, position1150)
								}
								break
							}
						}

					}
				l795:
					if !rules[ruleLPAREN]() {
						goto l794
					}
					if !rules[ruleexpression]() {
						goto l794
					}
					if !rules[ruleRPAREN]() {
						goto l794
					}
					goto l793
				l794:
					position, tokenIndex, depth = position793, tokenIndex793, depth793
					{
						position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
						{
							position1160 := position
							depth++
							{
								position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1162
								}
								position++
								goto l1161
							l1162:
								position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
								if buffer[position] != rune('S') {
									goto l1159
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
									goto l1159
								}
								position++
							}
						l1163:
							{
								position1165, tokenIndex1165, depth1165 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1166
								}
								position++
								goto l1165
							l1166:
								position, tokenIndex, depth = position1165, tokenIndex1165, depth1165
								if buffer[position] != rune('R') {
									goto l1159
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
									goto l1159
								}
								position++
							}
						l1167:
							{
								position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1170
								}
								position++
								goto l1169
							l1170:
								position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
								if buffer[position] != rune('T') {
									goto l1159
								}
								position++
							}
						l1169:
							{
								position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1172
								}
								position++
								goto l1171
							l1172:
								position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
								if buffer[position] != rune('A') {
									goto l1159
								}
								position++
							}
						l1171:
							{
								position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1174
								}
								position++
								goto l1173
							l1174:
								position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
								if buffer[position] != rune('R') {
									goto l1159
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
									goto l1159
								}
								position++
							}
						l1175:
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
									goto l1159
								}
								position++
							}
						l1177:
							if !rules[ruleskip]() {
								goto l1159
							}
							depth--
							add(ruleSTRSTARTS, position1160)
						}
						goto l1158
					l1159:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						{
							position1180 := position
							depth++
							{
								position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1182
								}
								position++
								goto l1181
							l1182:
								position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
								if buffer[position] != rune('S') {
									goto l1179
								}
								position++
							}
						l1181:
							{
								position1183, tokenIndex1183, depth1183 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1184
								}
								position++
								goto l1183
							l1184:
								position, tokenIndex, depth = position1183, tokenIndex1183, depth1183
								if buffer[position] != rune('T') {
									goto l1179
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
									goto l1179
								}
								position++
							}
						l1185:
							{
								position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1188
								}
								position++
								goto l1187
							l1188:
								position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
								if buffer[position] != rune('E') {
									goto l1179
								}
								position++
							}
						l1187:
							{
								position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1190
								}
								position++
								goto l1189
							l1190:
								position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
								if buffer[position] != rune('N') {
									goto l1179
								}
								position++
							}
						l1189:
							{
								position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1192
								}
								position++
								goto l1191
							l1192:
								position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
								if buffer[position] != rune('D') {
									goto l1179
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
									goto l1179
								}
								position++
							}
						l1193:
							if !rules[ruleskip]() {
								goto l1179
							}
							depth--
							add(ruleSTRENDS, position1180)
						}
						goto l1158
					l1179:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						{
							position1196 := position
							depth++
							{
								position1197, tokenIndex1197, depth1197 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1198
								}
								position++
								goto l1197
							l1198:
								position, tokenIndex, depth = position1197, tokenIndex1197, depth1197
								if buffer[position] != rune('S') {
									goto l1195
								}
								position++
							}
						l1197:
							{
								position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1200
								}
								position++
								goto l1199
							l1200:
								position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
								if buffer[position] != rune('T') {
									goto l1195
								}
								position++
							}
						l1199:
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
									goto l1195
								}
								position++
							}
						l1201:
							{
								position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1204
								}
								position++
								goto l1203
							l1204:
								position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
								if buffer[position] != rune('B') {
									goto l1195
								}
								position++
							}
						l1203:
							{
								position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1206
								}
								position++
								goto l1205
							l1206:
								position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
								if buffer[position] != rune('E') {
									goto l1195
								}
								position++
							}
						l1205:
							{
								position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1208
								}
								position++
								goto l1207
							l1208:
								position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
								if buffer[position] != rune('F') {
									goto l1195
								}
								position++
							}
						l1207:
							{
								position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1210
								}
								position++
								goto l1209
							l1210:
								position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
								if buffer[position] != rune('O') {
									goto l1195
								}
								position++
							}
						l1209:
							{
								position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1212
								}
								position++
								goto l1211
							l1212:
								position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
								if buffer[position] != rune('R') {
									goto l1195
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
									goto l1195
								}
								position++
							}
						l1213:
							if !rules[ruleskip]() {
								goto l1195
							}
							depth--
							add(ruleSTRBEFORE, position1196)
						}
						goto l1158
					l1195:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						{
							position1216 := position
							depth++
							{
								position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1218
								}
								position++
								goto l1217
							l1218:
								position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
								if buffer[position] != rune('S') {
									goto l1215
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
									goto l1215
								}
								position++
							}
						l1219:
							{
								position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1222
								}
								position++
								goto l1221
							l1222:
								position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
								if buffer[position] != rune('R') {
									goto l1215
								}
								position++
							}
						l1221:
							{
								position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1224
								}
								position++
								goto l1223
							l1224:
								position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
								if buffer[position] != rune('A') {
									goto l1215
								}
								position++
							}
						l1223:
							{
								position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1226
								}
								position++
								goto l1225
							l1226:
								position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
								if buffer[position] != rune('F') {
									goto l1215
								}
								position++
							}
						l1225:
							{
								position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1228
								}
								position++
								goto l1227
							l1228:
								position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
								if buffer[position] != rune('T') {
									goto l1215
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
									goto l1215
								}
								position++
							}
						l1229:
							{
								position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1232
								}
								position++
								goto l1231
							l1232:
								position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
								if buffer[position] != rune('R') {
									goto l1215
								}
								position++
							}
						l1231:
							if !rules[ruleskip]() {
								goto l1215
							}
							depth--
							add(ruleSTRAFTER, position1216)
						}
						goto l1158
					l1215:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						{
							position1234 := position
							depth++
							{
								position1235, tokenIndex1235, depth1235 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1236
								}
								position++
								goto l1235
							l1236:
								position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
								if buffer[position] != rune('S') {
									goto l1233
								}
								position++
							}
						l1235:
							{
								position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1238
								}
								position++
								goto l1237
							l1238:
								position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
								if buffer[position] != rune('T') {
									goto l1233
								}
								position++
							}
						l1237:
							{
								position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1240
								}
								position++
								goto l1239
							l1240:
								position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
								if buffer[position] != rune('R') {
									goto l1233
								}
								position++
							}
						l1239:
							{
								position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1242
								}
								position++
								goto l1241
							l1242:
								position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
								if buffer[position] != rune('L') {
									goto l1233
								}
								position++
							}
						l1241:
							{
								position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1244
								}
								position++
								goto l1243
							l1244:
								position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
								if buffer[position] != rune('A') {
									goto l1233
								}
								position++
							}
						l1243:
							{
								position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1246
								}
								position++
								goto l1245
							l1246:
								position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
								if buffer[position] != rune('N') {
									goto l1233
								}
								position++
							}
						l1245:
							{
								position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1248
								}
								position++
								goto l1247
							l1248:
								position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
								if buffer[position] != rune('G') {
									goto l1233
								}
								position++
							}
						l1247:
							if !rules[ruleskip]() {
								goto l1233
							}
							depth--
							add(ruleSTRLANG, position1234)
						}
						goto l1158
					l1233:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						{
							position1250 := position
							depth++
							{
								position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1252
								}
								position++
								goto l1251
							l1252:
								position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
								if buffer[position] != rune('S') {
									goto l1249
								}
								position++
							}
						l1251:
							{
								position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1254
								}
								position++
								goto l1253
							l1254:
								position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
								if buffer[position] != rune('T') {
									goto l1249
								}
								position++
							}
						l1253:
							{
								position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1256
								}
								position++
								goto l1255
							l1256:
								position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
								if buffer[position] != rune('R') {
									goto l1249
								}
								position++
							}
						l1255:
							{
								position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1258
								}
								position++
								goto l1257
							l1258:
								position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
								if buffer[position] != rune('D') {
									goto l1249
								}
								position++
							}
						l1257:
							{
								position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1260
								}
								position++
								goto l1259
							l1260:
								position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
								if buffer[position] != rune('T') {
									goto l1249
								}
								position++
							}
						l1259:
							if !rules[ruleskip]() {
								goto l1249
							}
							depth--
							add(ruleSTRDT, position1250)
						}
						goto l1158
					l1249:
						position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
						{
							switch buffer[position] {
							case 'S', 's':
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
											goto l1157
										}
										position++
									}
								l1263:
									{
										position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1266
										}
										position++
										goto l1265
									l1266:
										position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
										if buffer[position] != rune('A') {
											goto l1157
										}
										position++
									}
								l1265:
									{
										position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1268
										}
										position++
										goto l1267
									l1268:
										position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
										if buffer[position] != rune('M') {
											goto l1157
										}
										position++
									}
								l1267:
									{
										position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1270
										}
										position++
										goto l1269
									l1270:
										position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
										if buffer[position] != rune('E') {
											goto l1157
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
											goto l1157
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
											goto l1157
										}
										position++
									}
								l1273:
									{
										position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1276
										}
										position++
										goto l1275
									l1276:
										position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
										if buffer[position] != rune('R') {
											goto l1157
										}
										position++
									}
								l1275:
									{
										position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1278
										}
										position++
										goto l1277
									l1278:
										position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
										if buffer[position] != rune('M') {
											goto l1157
										}
										position++
									}
								l1277:
									if !rules[ruleskip]() {
										goto l1157
									}
									depth--
									add(ruleSAMETERM, position1262)
								}
								break
							case 'C', 'c':
								{
									position1279 := position
									depth++
									{
										position1280, tokenIndex1280, depth1280 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1281
										}
										position++
										goto l1280
									l1281:
										position, tokenIndex, depth = position1280, tokenIndex1280, depth1280
										if buffer[position] != rune('C') {
											goto l1157
										}
										position++
									}
								l1280:
									{
										position1282, tokenIndex1282, depth1282 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1283
										}
										position++
										goto l1282
									l1283:
										position, tokenIndex, depth = position1282, tokenIndex1282, depth1282
										if buffer[position] != rune('O') {
											goto l1157
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
											goto l1157
										}
										position++
									}
								l1284:
									{
										position1286, tokenIndex1286, depth1286 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1287
										}
										position++
										goto l1286
									l1287:
										position, tokenIndex, depth = position1286, tokenIndex1286, depth1286
										if buffer[position] != rune('T') {
											goto l1157
										}
										position++
									}
								l1286:
									{
										position1288, tokenIndex1288, depth1288 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1289
										}
										position++
										goto l1288
									l1289:
										position, tokenIndex, depth = position1288, tokenIndex1288, depth1288
										if buffer[position] != rune('A') {
											goto l1157
										}
										position++
									}
								l1288:
									{
										position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1291
										}
										position++
										goto l1290
									l1291:
										position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
										if buffer[position] != rune('I') {
											goto l1157
										}
										position++
									}
								l1290:
									{
										position1292, tokenIndex1292, depth1292 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1293
										}
										position++
										goto l1292
									l1293:
										position, tokenIndex, depth = position1292, tokenIndex1292, depth1292
										if buffer[position] != rune('N') {
											goto l1157
										}
										position++
									}
								l1292:
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
											goto l1157
										}
										position++
									}
								l1294:
									if !rules[ruleskip]() {
										goto l1157
									}
									depth--
									add(ruleCONTAINS, position1279)
								}
								break
							default:
								{
									position1296 := position
									depth++
									{
										position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1298
										}
										position++
										goto l1297
									l1298:
										position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
										if buffer[position] != rune('L') {
											goto l1157
										}
										position++
									}
								l1297:
									{
										position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1300
										}
										position++
										goto l1299
									l1300:
										position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
										if buffer[position] != rune('A') {
											goto l1157
										}
										position++
									}
								l1299:
									{
										position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1302
										}
										position++
										goto l1301
									l1302:
										position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
										if buffer[position] != rune('N') {
											goto l1157
										}
										position++
									}
								l1301:
									{
										position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1304
										}
										position++
										goto l1303
									l1304:
										position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
										if buffer[position] != rune('G') {
											goto l1157
										}
										position++
									}
								l1303:
									{
										position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1306
										}
										position++
										goto l1305
									l1306:
										position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
										if buffer[position] != rune('M') {
											goto l1157
										}
										position++
									}
								l1305:
									{
										position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1308
										}
										position++
										goto l1307
									l1308:
										position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
										if buffer[position] != rune('A') {
											goto l1157
										}
										position++
									}
								l1307:
									{
										position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1310
										}
										position++
										goto l1309
									l1310:
										position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
										if buffer[position] != rune('T') {
											goto l1157
										}
										position++
									}
								l1309:
									{
										position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1312
										}
										position++
										goto l1311
									l1312:
										position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
										if buffer[position] != rune('C') {
											goto l1157
										}
										position++
									}
								l1311:
									{
										position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1314
										}
										position++
										goto l1313
									l1314:
										position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
										if buffer[position] != rune('H') {
											goto l1157
										}
										position++
									}
								l1313:
									{
										position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1316
										}
										position++
										goto l1315
									l1316:
										position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
										if buffer[position] != rune('E') {
											goto l1157
										}
										position++
									}
								l1315:
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
											goto l1157
										}
										position++
									}
								l1317:
									if !rules[ruleskip]() {
										goto l1157
									}
									depth--
									add(ruleLANGMATCHES, position1296)
								}
								break
							}
						}

					}
				l1158:
					if !rules[ruleLPAREN]() {
						goto l1157
					}
					if !rules[ruleexpression]() {
						goto l1157
					}
					if !rules[ruleCOMMA]() {
						goto l1157
					}
					if !rules[ruleexpression]() {
						goto l1157
					}
					if !rules[ruleRPAREN]() {
						goto l1157
					}
					goto l793
				l1157:
					position, tokenIndex, depth = position793, tokenIndex793, depth793
					{
						position1320 := position
						depth++
						{
							position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1322
							}
							position++
							goto l1321
						l1322:
							position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
							if buffer[position] != rune('B') {
								goto l1319
							}
							position++
						}
					l1321:
						{
							position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1324
							}
							position++
							goto l1323
						l1324:
							position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
							if buffer[position] != rune('O') {
								goto l1319
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
								goto l1319
							}
							position++
						}
					l1325:
						{
							position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1328
							}
							position++
							goto l1327
						l1328:
							position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
							if buffer[position] != rune('N') {
								goto l1319
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
								goto l1319
							}
							position++
						}
					l1329:
						if !rules[ruleskip]() {
							goto l1319
						}
						depth--
						add(ruleBOUND, position1320)
					}
					if !rules[ruleLPAREN]() {
						goto l1319
					}
					if !rules[rulevar]() {
						goto l1319
					}
					if !rules[ruleRPAREN]() {
						goto l1319
					}
					goto l793
				l1319:
					position, tokenIndex, depth = position793, tokenIndex793, depth793
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1333 := position
								depth++
								{
									position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1335
									}
									position++
									goto l1334
								l1335:
									position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
									if buffer[position] != rune('S') {
										goto l1331
									}
									position++
								}
							l1334:
								{
									position1336, tokenIndex1336, depth1336 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1337
									}
									position++
									goto l1336
								l1337:
									position, tokenIndex, depth = position1336, tokenIndex1336, depth1336
									if buffer[position] != rune('T') {
										goto l1331
									}
									position++
								}
							l1336:
								{
									position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1339
									}
									position++
									goto l1338
								l1339:
									position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
									if buffer[position] != rune('R') {
										goto l1331
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
										goto l1331
									}
									position++
								}
							l1340:
								{
									position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1343
									}
									position++
									goto l1342
								l1343:
									position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
									if buffer[position] != rune('U') {
										goto l1331
									}
									position++
								}
							l1342:
								{
									position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1345
									}
									position++
									goto l1344
								l1345:
									position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
									if buffer[position] != rune('I') {
										goto l1331
									}
									position++
								}
							l1344:
								{
									position1346, tokenIndex1346, depth1346 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1347
									}
									position++
									goto l1346
								l1347:
									position, tokenIndex, depth = position1346, tokenIndex1346, depth1346
									if buffer[position] != rune('D') {
										goto l1331
									}
									position++
								}
							l1346:
								if !rules[ruleskip]() {
									goto l1331
								}
								depth--
								add(ruleSTRUUID, position1333)
							}
							break
						case 'U', 'u':
							{
								position1348 := position
								depth++
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
										goto l1331
									}
									position++
								}
							l1349:
								{
									position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1352
									}
									position++
									goto l1351
								l1352:
									position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
									if buffer[position] != rune('U') {
										goto l1331
									}
									position++
								}
							l1351:
								{
									position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1354
									}
									position++
									goto l1353
								l1354:
									position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
									if buffer[position] != rune('I') {
										goto l1331
									}
									position++
								}
							l1353:
								{
									position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1356
									}
									position++
									goto l1355
								l1356:
									position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
									if buffer[position] != rune('D') {
										goto l1331
									}
									position++
								}
							l1355:
								if !rules[ruleskip]() {
									goto l1331
								}
								depth--
								add(ruleUUID, position1348)
							}
							break
						case 'N', 'n':
							{
								position1357 := position
								depth++
								{
									position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1359
									}
									position++
									goto l1358
								l1359:
									position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
									if buffer[position] != rune('N') {
										goto l1331
									}
									position++
								}
							l1358:
								{
									position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1361
									}
									position++
									goto l1360
								l1361:
									position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
									if buffer[position] != rune('O') {
										goto l1331
									}
									position++
								}
							l1360:
								{
									position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1363
									}
									position++
									goto l1362
								l1363:
									position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
									if buffer[position] != rune('W') {
										goto l1331
									}
									position++
								}
							l1362:
								if !rules[ruleskip]() {
									goto l1331
								}
								depth--
								add(ruleNOW, position1357)
							}
							break
						default:
							{
								position1364 := position
								depth++
								{
									position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1366
									}
									position++
									goto l1365
								l1366:
									position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
									if buffer[position] != rune('R') {
										goto l1331
									}
									position++
								}
							l1365:
								{
									position1367, tokenIndex1367, depth1367 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1368
									}
									position++
									goto l1367
								l1368:
									position, tokenIndex, depth = position1367, tokenIndex1367, depth1367
									if buffer[position] != rune('A') {
										goto l1331
									}
									position++
								}
							l1367:
								{
									position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1370
									}
									position++
									goto l1369
								l1370:
									position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
									if buffer[position] != rune('N') {
										goto l1331
									}
									position++
								}
							l1369:
								{
									position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1372
									}
									position++
									goto l1371
								l1372:
									position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
									if buffer[position] != rune('D') {
										goto l1331
									}
									position++
								}
							l1371:
								if !rules[ruleskip]() {
									goto l1331
								}
								depth--
								add(ruleRAND, position1364)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1331
					}
					goto l793
				l1331:
					position, tokenIndex, depth = position793, tokenIndex793, depth793
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
								{
									position1376 := position
									depth++
									{
										position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1378
										}
										position++
										goto l1377
									l1378:
										position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
										if buffer[position] != rune('E') {
											goto l1375
										}
										position++
									}
								l1377:
									{
										position1379, tokenIndex1379, depth1379 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1380
										}
										position++
										goto l1379
									l1380:
										position, tokenIndex, depth = position1379, tokenIndex1379, depth1379
										if buffer[position] != rune('X') {
											goto l1375
										}
										position++
									}
								l1379:
									{
										position1381, tokenIndex1381, depth1381 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1382
										}
										position++
										goto l1381
									l1382:
										position, tokenIndex, depth = position1381, tokenIndex1381, depth1381
										if buffer[position] != rune('I') {
											goto l1375
										}
										position++
									}
								l1381:
									{
										position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1384
										}
										position++
										goto l1383
									l1384:
										position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
										if buffer[position] != rune('S') {
											goto l1375
										}
										position++
									}
								l1383:
									{
										position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1386
										}
										position++
										goto l1385
									l1386:
										position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
										if buffer[position] != rune('T') {
											goto l1375
										}
										position++
									}
								l1385:
									{
										position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1388
										}
										position++
										goto l1387
									l1388:
										position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
										if buffer[position] != rune('S') {
											goto l1375
										}
										position++
									}
								l1387:
									if !rules[ruleskip]() {
										goto l1375
									}
									depth--
									add(ruleEXISTS, position1376)
								}
								goto l1374
							l1375:
								position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
								{
									position1389 := position
									depth++
									{
										position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1391
										}
										position++
										goto l1390
									l1391:
										position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
										if buffer[position] != rune('N') {
											goto l791
										}
										position++
									}
								l1390:
									{
										position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1393
										}
										position++
										goto l1392
									l1393:
										position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
										if buffer[position] != rune('O') {
											goto l791
										}
										position++
									}
								l1392:
									{
										position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1395
										}
										position++
										goto l1394
									l1395:
										position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
										if buffer[position] != rune('T') {
											goto l791
										}
										position++
									}
								l1394:
									if buffer[position] != rune(' ') {
										goto l791
									}
									position++
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
											goto l791
										}
										position++
									}
								l1396:
									{
										position1398, tokenIndex1398, depth1398 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1399
										}
										position++
										goto l1398
									l1399:
										position, tokenIndex, depth = position1398, tokenIndex1398, depth1398
										if buffer[position] != rune('X') {
											goto l791
										}
										position++
									}
								l1398:
									{
										position1400, tokenIndex1400, depth1400 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1401
										}
										position++
										goto l1400
									l1401:
										position, tokenIndex, depth = position1400, tokenIndex1400, depth1400
										if buffer[position] != rune('I') {
											goto l791
										}
										position++
									}
								l1400:
									{
										position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1403
										}
										position++
										goto l1402
									l1403:
										position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
										if buffer[position] != rune('S') {
											goto l791
										}
										position++
									}
								l1402:
									{
										position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1405
										}
										position++
										goto l1404
									l1405:
										position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
										if buffer[position] != rune('T') {
											goto l791
										}
										position++
									}
								l1404:
									{
										position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1407
										}
										position++
										goto l1406
									l1407:
										position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
										if buffer[position] != rune('S') {
											goto l791
										}
										position++
									}
								l1406:
									if !rules[ruleskip]() {
										goto l791
									}
									depth--
									add(ruleNOTEXIST, position1389)
								}
							}
						l1374:
							if !rules[rulegroupGraphPattern]() {
								goto l791
							}
							break
						case 'I', 'i':
							{
								position1408 := position
								depth++
								{
									position1409, tokenIndex1409, depth1409 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1410
									}
									position++
									goto l1409
								l1410:
									position, tokenIndex, depth = position1409, tokenIndex1409, depth1409
									if buffer[position] != rune('I') {
										goto l791
									}
									position++
								}
							l1409:
								{
									position1411, tokenIndex1411, depth1411 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1412
									}
									position++
									goto l1411
								l1412:
									position, tokenIndex, depth = position1411, tokenIndex1411, depth1411
									if buffer[position] != rune('F') {
										goto l791
									}
									position++
								}
							l1411:
								if !rules[ruleskip]() {
									goto l791
								}
								depth--
								add(ruleIF, position1408)
							}
							if !rules[ruleLPAREN]() {
								goto l791
							}
							if !rules[ruleexpression]() {
								goto l791
							}
							if !rules[ruleCOMMA]() {
								goto l791
							}
							if !rules[ruleexpression]() {
								goto l791
							}
							if !rules[ruleCOMMA]() {
								goto l791
							}
							if !rules[ruleexpression]() {
								goto l791
							}
							if !rules[ruleRPAREN]() {
								goto l791
							}
							break
						case 'C', 'c':
							{
								position1413, tokenIndex1413, depth1413 := position, tokenIndex, depth
								{
									position1415 := position
									depth++
									{
										position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1417
										}
										position++
										goto l1416
									l1417:
										position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
										if buffer[position] != rune('C') {
											goto l1414
										}
										position++
									}
								l1416:
									{
										position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1419
										}
										position++
										goto l1418
									l1419:
										position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
										if buffer[position] != rune('O') {
											goto l1414
										}
										position++
									}
								l1418:
									{
										position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1421
										}
										position++
										goto l1420
									l1421:
										position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
										if buffer[position] != rune('N') {
											goto l1414
										}
										position++
									}
								l1420:
									{
										position1422, tokenIndex1422, depth1422 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1423
										}
										position++
										goto l1422
									l1423:
										position, tokenIndex, depth = position1422, tokenIndex1422, depth1422
										if buffer[position] != rune('C') {
											goto l1414
										}
										position++
									}
								l1422:
									{
										position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1425
										}
										position++
										goto l1424
									l1425:
										position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
										if buffer[position] != rune('A') {
											goto l1414
										}
										position++
									}
								l1424:
									{
										position1426, tokenIndex1426, depth1426 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1427
										}
										position++
										goto l1426
									l1427:
										position, tokenIndex, depth = position1426, tokenIndex1426, depth1426
										if buffer[position] != rune('T') {
											goto l1414
										}
										position++
									}
								l1426:
									if !rules[ruleskip]() {
										goto l1414
									}
									depth--
									add(ruleCONCAT, position1415)
								}
								goto l1413
							l1414:
								position, tokenIndex, depth = position1413, tokenIndex1413, depth1413
								{
									position1428 := position
									depth++
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
											goto l791
										}
										position++
									}
								l1429:
									{
										position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1432
										}
										position++
										goto l1431
									l1432:
										position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
										if buffer[position] != rune('O') {
											goto l791
										}
										position++
									}
								l1431:
									{
										position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1434
										}
										position++
										goto l1433
									l1434:
										position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
										if buffer[position] != rune('A') {
											goto l791
										}
										position++
									}
								l1433:
									{
										position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1436
										}
										position++
										goto l1435
									l1436:
										position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
										if buffer[position] != rune('L') {
											goto l791
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
											goto l791
										}
										position++
									}
								l1437:
									{
										position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1440
										}
										position++
										goto l1439
									l1440:
										position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
										if buffer[position] != rune('S') {
											goto l791
										}
										position++
									}
								l1439:
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
											goto l791
										}
										position++
									}
								l1441:
									{
										position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1444
										}
										position++
										goto l1443
									l1444:
										position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
										if buffer[position] != rune('E') {
											goto l791
										}
										position++
									}
								l1443:
									if !rules[ruleskip]() {
										goto l791
									}
									depth--
									add(ruleCOALESCE, position1428)
								}
							}
						l1413:
							if !rules[ruleargList]() {
								goto l791
							}
							break
						case 'B', 'b':
							{
								position1445 := position
								depth++
								{
									position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1447
									}
									position++
									goto l1446
								l1447:
									position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
									if buffer[position] != rune('B') {
										goto l791
									}
									position++
								}
							l1446:
								{
									position1448, tokenIndex1448, depth1448 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1449
									}
									position++
									goto l1448
								l1449:
									position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
									if buffer[position] != rune('N') {
										goto l791
									}
									position++
								}
							l1448:
								{
									position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1451
									}
									position++
									goto l1450
								l1451:
									position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
									if buffer[position] != rune('O') {
										goto l791
									}
									position++
								}
							l1450:
								{
									position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1453
									}
									position++
									goto l1452
								l1453:
									position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
									if buffer[position] != rune('D') {
										goto l791
									}
									position++
								}
							l1452:
								{
									position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1455
									}
									position++
									goto l1454
								l1455:
									position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
									if buffer[position] != rune('E') {
										goto l791
									}
									position++
								}
							l1454:
								if !rules[ruleskip]() {
									goto l791
								}
								depth--
								add(ruleBNODE, position1445)
							}
							{
								position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1457
								}
								if !rules[ruleexpression]() {
									goto l1457
								}
								if !rules[ruleRPAREN]() {
									goto l1457
								}
								goto l1456
							l1457:
								position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
								if !rules[rulenil]() {
									goto l791
								}
							}
						l1456:
							break
						default:
							{
								position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
								{
									position1460 := position
									depth++
									{
										position1461, tokenIndex1461, depth1461 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1462
										}
										position++
										goto l1461
									l1462:
										position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
										if buffer[position] != rune('S') {
											goto l1459
										}
										position++
									}
								l1461:
									{
										position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1464
										}
										position++
										goto l1463
									l1464:
										position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
										if buffer[position] != rune('U') {
											goto l1459
										}
										position++
									}
								l1463:
									{
										position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1466
										}
										position++
										goto l1465
									l1466:
										position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
										if buffer[position] != rune('B') {
											goto l1459
										}
										position++
									}
								l1465:
									{
										position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1468
										}
										position++
										goto l1467
									l1468:
										position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
										if buffer[position] != rune('S') {
											goto l1459
										}
										position++
									}
								l1467:
									{
										position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1470
										}
										position++
										goto l1469
									l1470:
										position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
										if buffer[position] != rune('T') {
											goto l1459
										}
										position++
									}
								l1469:
									{
										position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1472
										}
										position++
										goto l1471
									l1472:
										position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
										if buffer[position] != rune('R') {
											goto l1459
										}
										position++
									}
								l1471:
									if !rules[ruleskip]() {
										goto l1459
									}
									depth--
									add(ruleSUBSTR, position1460)
								}
								goto l1458
							l1459:
								position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
								{
									position1474 := position
									depth++
									{
										position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1476
										}
										position++
										goto l1475
									l1476:
										position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
										if buffer[position] != rune('R') {
											goto l1473
										}
										position++
									}
								l1475:
									{
										position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1478
										}
										position++
										goto l1477
									l1478:
										position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
										if buffer[position] != rune('E') {
											goto l1473
										}
										position++
									}
								l1477:
									{
										position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1480
										}
										position++
										goto l1479
									l1480:
										position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
										if buffer[position] != rune('P') {
											goto l1473
										}
										position++
									}
								l1479:
									{
										position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1482
										}
										position++
										goto l1481
									l1482:
										position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
										if buffer[position] != rune('L') {
											goto l1473
										}
										position++
									}
								l1481:
									{
										position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1484
										}
										position++
										goto l1483
									l1484:
										position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
										if buffer[position] != rune('A') {
											goto l1473
										}
										position++
									}
								l1483:
									{
										position1485, tokenIndex1485, depth1485 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1486
										}
										position++
										goto l1485
									l1486:
										position, tokenIndex, depth = position1485, tokenIndex1485, depth1485
										if buffer[position] != rune('C') {
											goto l1473
										}
										position++
									}
								l1485:
									{
										position1487, tokenIndex1487, depth1487 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1488
										}
										position++
										goto l1487
									l1488:
										position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
										if buffer[position] != rune('E') {
											goto l1473
										}
										position++
									}
								l1487:
									if !rules[ruleskip]() {
										goto l1473
									}
									depth--
									add(ruleREPLACE, position1474)
								}
								goto l1458
							l1473:
								position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
								{
									position1489 := position
									depth++
									{
										position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1491
										}
										position++
										goto l1490
									l1491:
										position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
										if buffer[position] != rune('R') {
											goto l791
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
											goto l791
										}
										position++
									}
								l1492:
									{
										position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1495
										}
										position++
										goto l1494
									l1495:
										position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
										if buffer[position] != rune('G') {
											goto l791
										}
										position++
									}
								l1494:
									{
										position1496, tokenIndex1496, depth1496 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1497
										}
										position++
										goto l1496
									l1497:
										position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
										if buffer[position] != rune('E') {
											goto l791
										}
										position++
									}
								l1496:
									{
										position1498, tokenIndex1498, depth1498 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1499
										}
										position++
										goto l1498
									l1499:
										position, tokenIndex, depth = position1498, tokenIndex1498, depth1498
										if buffer[position] != rune('X') {
											goto l791
										}
										position++
									}
								l1498:
									if !rules[ruleskip]() {
										goto l791
									}
									depth--
									add(ruleREGEX, position1489)
								}
							}
						l1458:
							if !rules[ruleLPAREN]() {
								goto l791
							}
							if !rules[ruleexpression]() {
								goto l791
							}
							if !rules[ruleCOMMA]() {
								goto l791
							}
							if !rules[ruleexpression]() {
								goto l791
							}
							{
								position1500, tokenIndex1500, depth1500 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1500
								}
								if !rules[ruleexpression]() {
									goto l1500
								}
								goto l1501
							l1500:
								position, tokenIndex, depth = position1500, tokenIndex1500, depth1500
							}
						l1501:
							if !rules[ruleRPAREN]() {
								goto l791
							}
							break
						}
					}

				}
			l793:
				depth--
				add(rulebuiltinCall, position792)
			}
			return true
		l791:
			position, tokenIndex, depth = position791, tokenIndex791, depth791
			return false
		},
		/* 69 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
			{
				position1503 := position
				depth++
				{
					position1504, tokenIndex1504, depth1504 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1505
					}
					position++
					goto l1504
				l1505:
					position, tokenIndex, depth = position1504, tokenIndex1504, depth1504
					if buffer[position] != rune('$') {
						goto l1502
					}
					position++
				}
			l1504:
				{
					position1506 := position
					depth++
					{
						position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
						if !rules[rulepnCharsU]() {
							goto l1508
						}
						goto l1507
					l1508:
						position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1502
						}
						position++
					}
				l1507:
				l1509:
					{
						position1510, tokenIndex1510, depth1510 := position, tokenIndex, depth
						{
							position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
							if !rules[rulepnCharsU]() {
								goto l1512
							}
							goto l1511
						l1512:
							position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
							{
								switch buffer[position] {
								case 'â':
									if c := buffer[position]; c < rune('‿') || c > rune('⁀') {
										goto l1510
									}
									position++
									break
								case 'Ì', 'Í':
									if c := buffer[position]; c < rune('̀') || c > rune('ͯ') {
										goto l1510
									}
									position++
									break
								case 'Â':
									if buffer[position] != rune('·') {
										goto l1510
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1510
									}
									position++
									break
								}
							}

						}
					l1511:
						goto l1509
					l1510:
						position, tokenIndex, depth = position1510, tokenIndex1510, depth1510
					}
					depth--
					add(ruleVARNAME, position1506)
				}
				if !rules[ruleskip]() {
					goto l1502
				}
				depth--
				add(rulevar, position1503)
			}
			return true
		l1502:
			position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
			return false
		},
		/* 70 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1514, tokenIndex1514, depth1514 := position, tokenIndex, depth
			{
				position1515 := position
				depth++
				{
					position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1517
					}
					goto l1516
				l1517:
					position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
					{
						position1518 := position
						depth++
						{
							position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
							if !rules[rulepnPrefix]() {
								goto l1519
							}
							goto l1520
						l1519:
							position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
						}
					l1520:
						if buffer[position] != rune(':') {
							goto l1514
						}
						position++
						{
							position1521 := position
							depth++
							{
								switch buffer[position] {
								case '%', '\\':
									{
										position1525 := position
										depth++
										{
											position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
											{
												position1528 := position
												depth++
												if buffer[position] != rune('%') {
													goto l1527
												}
												position++
												if !rules[rulehex]() {
													goto l1527
												}
												if !rules[rulehex]() {
													goto l1527
												}
												depth--
												add(rulepercent, position1528)
											}
											goto l1526
										l1527:
											position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
											{
												position1529 := position
												depth++
												if buffer[position] != rune('\\') {
													goto l1514
												}
												position++
												{
													switch buffer[position] {
													case '%':
														if buffer[position] != rune('%') {
															goto l1514
														}
														position++
														break
													case '@':
														if buffer[position] != rune('@') {
															goto l1514
														}
														position++
														break
													case '#':
														if buffer[position] != rune('#') {
															goto l1514
														}
														position++
														break
													case '?':
														if buffer[position] != rune('?') {
															goto l1514
														}
														position++
														break
													case '/':
														if buffer[position] != rune('/') {
															goto l1514
														}
														position++
														break
													case '=':
														if buffer[position] != rune('=') {
															goto l1514
														}
														position++
														break
													case ';':
														if buffer[position] != rune(';') {
															goto l1514
														}
														position++
														break
													case ',':
														if buffer[position] != rune(',') {
															goto l1514
														}
														position++
														break
													case '+':
														if buffer[position] != rune('+') {
															goto l1514
														}
														position++
														break
													case '*':
														if buffer[position] != rune('*') {
															goto l1514
														}
														position++
														break
													case ')':
														if buffer[position] != rune(')') {
															goto l1514
														}
														position++
														break
													case '(':
														if buffer[position] != rune('(') {
															goto l1514
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l1514
														}
														position++
														break
													case '&':
														if buffer[position] != rune('&') {
															goto l1514
														}
														position++
														break
													case '$':
														if buffer[position] != rune('$') {
															goto l1514
														}
														position++
														break
													case '!':
														if buffer[position] != rune('!') {
															goto l1514
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l1514
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l1514
														}
														position++
														break
													case '~':
														if buffer[position] != rune('~') {
															goto l1514
														}
														position++
														break
													default:
														if buffer[position] != rune('_') {
															goto l1514
														}
														position++
														break
													}
												}

												depth--
												add(rulepnLocalEsc, position1529)
											}
										}
									l1526:
										depth--
										add(ruleplx, position1525)
									}
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1514
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l1514
									}
									position++
									break
								default:
									if !rules[rulepnCharsU]() {
										goto l1514
									}
									break
								}
							}

						l1522:
							{
								position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '%', '\\':
										{
											position1532 := position
											depth++
											{
												position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
												{
													position1535 := position
													depth++
													if buffer[position] != rune('%') {
														goto l1534
													}
													position++
													if !rules[rulehex]() {
														goto l1534
													}
													if !rules[rulehex]() {
														goto l1534
													}
													depth--
													add(rulepercent, position1535)
												}
												goto l1533
											l1534:
												position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
												{
													position1536 := position
													depth++
													if buffer[position] != rune('\\') {
														goto l1523
													}
													position++
													{
														switch buffer[position] {
														case '%':
															if buffer[position] != rune('%') {
																goto l1523
															}
															position++
															break
														case '@':
															if buffer[position] != rune('@') {
																goto l1523
															}
															position++
															break
														case '#':
															if buffer[position] != rune('#') {
																goto l1523
															}
															position++
															break
														case '?':
															if buffer[position] != rune('?') {
																goto l1523
															}
															position++
															break
														case '/':
															if buffer[position] != rune('/') {
																goto l1523
															}
															position++
															break
														case '=':
															if buffer[position] != rune('=') {
																goto l1523
															}
															position++
															break
														case ';':
															if buffer[position] != rune(';') {
																goto l1523
															}
															position++
															break
														case ',':
															if buffer[position] != rune(',') {
																goto l1523
															}
															position++
															break
														case '+':
															if buffer[position] != rune('+') {
																goto l1523
															}
															position++
															break
														case '*':
															if buffer[position] != rune('*') {
																goto l1523
															}
															position++
															break
														case ')':
															if buffer[position] != rune(')') {
																goto l1523
															}
															position++
															break
														case '(':
															if buffer[position] != rune('(') {
																goto l1523
															}
															position++
															break
														case '\'':
															if buffer[position] != rune('\'') {
																goto l1523
															}
															position++
															break
														case '&':
															if buffer[position] != rune('&') {
																goto l1523
															}
															position++
															break
														case '$':
															if buffer[position] != rune('$') {
																goto l1523
															}
															position++
															break
														case '!':
															if buffer[position] != rune('!') {
																goto l1523
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l1523
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l1523
															}
															position++
															break
														case '~':
															if buffer[position] != rune('~') {
																goto l1523
															}
															position++
															break
														default:
															if buffer[position] != rune('_') {
																goto l1523
															}
															position++
															break
														}
													}

													depth--
													add(rulepnLocalEsc, position1536)
												}
											}
										l1533:
											depth--
											add(ruleplx, position1532)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1523
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1523
										}
										position++
										break
									default:
										if !rules[rulepnCharsU]() {
											goto l1523
										}
										break
									}
								}

								goto l1522
							l1523:
								position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
							}
							depth--
							add(rulepnLocal, position1521)
						}
						if !rules[ruleskip]() {
							goto l1514
						}
						depth--
						add(ruleprefixedName, position1518)
					}
				}
			l1516:
				depth--
				add(ruleiriref, position1515)
			}
			return true
		l1514:
			position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
			return false
		},
		/* 71 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
			{
				position1539 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1538
				}
				position++
			l1540:
				{
					position1541, tokenIndex1541, depth1541 := position, tokenIndex, depth
					{
						position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1542
						}
						position++
						goto l1541
					l1542:
						position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
					}
					if !matchDot() {
						goto l1541
					}
					goto l1540
				l1541:
					position, tokenIndex, depth = position1541, tokenIndex1541, depth1541
				}
				if buffer[position] != rune('>') {
					goto l1538
				}
				position++
				if !rules[ruleskip]() {
					goto l1538
				}
				depth--
				add(ruleiri, position1539)
			}
			return true
		l1538:
			position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
			return false
		},
		/* 72 prefixedName <- <(pnPrefix? ':' pnLocal skip)> */
		nil,
		/* 73 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1544, tokenIndex1544, depth1544 := position, tokenIndex, depth
			{
				position1545 := position
				depth++
				if !rules[rulestring]() {
					goto l1544
				}
				{
					position1546, tokenIndex1546, depth1546 := position, tokenIndex, depth
					{
						position1548, tokenIndex1548, depth1548 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1549
						}
						position++
						{
							position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1553
							}
							position++
							goto l1552
						l1553:
							position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1549
							}
							position++
						}
					l1552:
					l1550:
						{
							position1551, tokenIndex1551, depth1551 := position, tokenIndex, depth
							{
								position1554, tokenIndex1554, depth1554 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1555
								}
								position++
								goto l1554
							l1555:
								position, tokenIndex, depth = position1554, tokenIndex1554, depth1554
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1551
								}
								position++
							}
						l1554:
							goto l1550
						l1551:
							position, tokenIndex, depth = position1551, tokenIndex1551, depth1551
						}
					l1556:
						{
							position1557, tokenIndex1557, depth1557 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1557
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1557
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1557
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1557
									}
									position++
									break
								}
							}

						l1558:
							{
								position1559, tokenIndex1559, depth1559 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1559
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1559
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1559
										}
										position++
										break
									}
								}

								goto l1558
							l1559:
								position, tokenIndex, depth = position1559, tokenIndex1559, depth1559
							}
							goto l1556
						l1557:
							position, tokenIndex, depth = position1557, tokenIndex1557, depth1557
						}
						goto l1548
					l1549:
						position, tokenIndex, depth = position1548, tokenIndex1548, depth1548
						if buffer[position] != rune('^') {
							goto l1546
						}
						position++
						if buffer[position] != rune('^') {
							goto l1546
						}
						position++
						if !rules[ruleiriref]() {
							goto l1546
						}
					}
				l1548:
					goto l1547
				l1546:
					position, tokenIndex, depth = position1546, tokenIndex1546, depth1546
				}
			l1547:
				if !rules[ruleskip]() {
					goto l1544
				}
				depth--
				add(ruleliteral, position1545)
			}
			return true
		l1544:
			position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
			return false
		},
		/* 74 string <- <(stringLiteralA / stringLiteralB / stringLiteralLongA / stringLiteralLongB)> */
		func() bool {
			position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
			{
				position1563 := position
				depth++
				{
					position1564, tokenIndex1564, depth1564 := position, tokenIndex, depth
					{
						position1566 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1565
						}
						position++
					l1567:
						{
							position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
							{
								position1569, tokenIndex1569, depth1569 := position, tokenIndex, depth
								{
									position1571, tokenIndex1571, depth1571 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1571
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1571
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1571
											}
											position++
											break
										default:
											if buffer[position] != rune('\'') {
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
								if !matchDot() {
									goto l1570
								}
								goto l1569
							l1570:
								position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
								if !rules[ruleechar]() {
									goto l1568
								}
							}
						l1569:
							goto l1567
						l1568:
							position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
						}
						if buffer[position] != rune('\'') {
							goto l1565
						}
						position++
						depth--
						add(rulestringLiteralA, position1566)
					}
					goto l1564
				l1565:
					position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
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
								{
									position1579, tokenIndex1579, depth1579 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1579
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1579
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1579
											}
											position++
											break
										default:
											if buffer[position] != rune('"') {
												goto l1579
											}
											position++
											break
										}
									}

									goto l1578
								l1579:
									position, tokenIndex, depth = position1579, tokenIndex1579, depth1579
								}
								if !matchDot() {
									goto l1578
								}
								goto l1577
							l1578:
								position, tokenIndex, depth = position1577, tokenIndex1577, depth1577
								if !rules[ruleechar]() {
									goto l1576
								}
							}
						l1577:
							goto l1575
						l1576:
							position, tokenIndex, depth = position1576, tokenIndex1576, depth1576
						}
						if buffer[position] != rune('"') {
							goto l1573
						}
						position++
						depth--
						add(rulestringLiteralB, position1574)
					}
					goto l1564
				l1573:
					position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
					{
						position1582 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1581
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1581
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1581
						}
						position++
					l1583:
						{
							position1584, tokenIndex1584, depth1584 := position, tokenIndex, depth
							{
								position1585, tokenIndex1585, depth1585 := position, tokenIndex, depth
								{
									position1587, tokenIndex1587, depth1587 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l1588
									}
									position++
									goto l1587
								l1588:
									position, tokenIndex, depth = position1587, tokenIndex1587, depth1587
									if buffer[position] != rune('\'') {
										goto l1585
									}
									position++
									if buffer[position] != rune('\'') {
										goto l1585
									}
									position++
								}
							l1587:
								goto l1586
							l1585:
								position, tokenIndex, depth = position1585, tokenIndex1585, depth1585
							}
						l1586:
							{
								position1589, tokenIndex1589, depth1589 := position, tokenIndex, depth
								{
									position1591, tokenIndex1591, depth1591 := position, tokenIndex, depth
									{
										position1592, tokenIndex1592, depth1592 := position, tokenIndex, depth
										if buffer[position] != rune('\'') {
											goto l1593
										}
										position++
										goto l1592
									l1593:
										position, tokenIndex, depth = position1592, tokenIndex1592, depth1592
										if buffer[position] != rune('\\') {
											goto l1591
										}
										position++
									}
								l1592:
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
									goto l1584
								}
							}
						l1589:
							goto l1583
						l1584:
							position, tokenIndex, depth = position1584, tokenIndex1584, depth1584
						}
						if buffer[position] != rune('\'') {
							goto l1581
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1581
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1581
						}
						position++
						depth--
						add(rulestringLiteralLongA, position1582)
					}
					goto l1564
				l1581:
					position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
					{
						position1594 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1562
						}
						position++
						if buffer[position] != rune('"') {
							goto l1562
						}
						position++
						if buffer[position] != rune('"') {
							goto l1562
						}
						position++
					l1595:
						{
							position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
							{
								position1597, tokenIndex1597, depth1597 := position, tokenIndex, depth
								{
									position1599, tokenIndex1599, depth1599 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l1600
									}
									position++
									goto l1599
								l1600:
									position, tokenIndex, depth = position1599, tokenIndex1599, depth1599
									if buffer[position] != rune('"') {
										goto l1597
									}
									position++
									if buffer[position] != rune('"') {
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
										if buffer[position] != rune('"') {
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
						if buffer[position] != rune('"') {
							goto l1562
						}
						position++
						if buffer[position] != rune('"') {
							goto l1562
						}
						position++
						if buffer[position] != rune('"') {
							goto l1562
						}
						position++
						depth--
						add(rulestringLiteralLongB, position1594)
					}
				}
			l1564:
				depth--
				add(rulestring, position1563)
			}
			return true
		l1562:
			position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
			return false
		},
		/* 75 stringLiteralA <- <('\'' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('\'') '\'')) .) / echar)* '\'')> */
		nil,
		/* 76 stringLiteralB <- <('"' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .) / echar)* '"')> */
		nil,
		/* 77 stringLiteralLongA <- <('\'' '\'' '\'' (('\'' / ('\'' '\''))? ((!('\'' / '\\') .) / echar))* ('\'' '\'' '\''))> */
		nil,
		/* 78 stringLiteralLongB <- <('"' '"' '"' (('"' / ('"' '"'))? ((!('"' / '\\') .) / echar))* ('"' '"' '"'))> */
		nil,
		/* 79 echar <- <('\\' ((&('\'') '\'') | (&('"') '"') | (&('\\') '\\') | (&('f') 'f') | (&('r') 'r') | (&('n') 'n') | (&('b') 'b') | (&('t') 't')))> */
		func() bool {
			position1610, tokenIndex1610, depth1610 := position, tokenIndex, depth
			{
				position1611 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l1610
				}
				position++
				{
					switch buffer[position] {
					case '\'':
						if buffer[position] != rune('\'') {
							goto l1610
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l1610
						}
						position++
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l1610
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l1610
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l1610
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l1610
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l1610
						}
						position++
						break
					default:
						if buffer[position] != rune('t') {
							goto l1610
						}
						position++
						break
					}
				}

				depth--
				add(ruleechar, position1611)
			}
			return true
		l1610:
			position, tokenIndex, depth = position1610, tokenIndex1610, depth1610
			return false
		},
		/* 80 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
			{
				position1614 := position
				depth++
				{
					position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
					{
						position1617, tokenIndex1617, depth1617 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1618
						}
						position++
						goto l1617
					l1618:
						position, tokenIndex, depth = position1617, tokenIndex1617, depth1617
						if buffer[position] != rune('-') {
							goto l1615
						}
						position++
					}
				l1617:
					goto l1616
				l1615:
					position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
				}
			l1616:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1613
				}
				position++
			l1619:
				{
					position1620, tokenIndex1620, depth1620 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1620
					}
					position++
					goto l1619
				l1620:
					position, tokenIndex, depth = position1620, tokenIndex1620, depth1620
				}
				{
					position1621, tokenIndex1621, depth1621 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1621
					}
					position++
				l1623:
					{
						position1624, tokenIndex1624, depth1624 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1624
						}
						position++
						goto l1623
					l1624:
						position, tokenIndex, depth = position1624, tokenIndex1624, depth1624
					}
					goto l1622
				l1621:
					position, tokenIndex, depth = position1621, tokenIndex1621, depth1621
				}
			l1622:
				if !rules[ruleskip]() {
					goto l1613
				}
				depth--
				add(rulenumericLiteral, position1614)
			}
			return true
		l1613:
			position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
			return false
		},
		/* 81 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 82 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1626, tokenIndex1626, depth1626 := position, tokenIndex, depth
			{
				position1627 := position
				depth++
				{
					position1628, tokenIndex1628, depth1628 := position, tokenIndex, depth
					{
						position1630 := position
						depth++
						{
							position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1632
							}
							position++
							goto l1631
						l1632:
							position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
							if buffer[position] != rune('T') {
								goto l1629
							}
							position++
						}
					l1631:
						{
							position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1634
							}
							position++
							goto l1633
						l1634:
							position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
							if buffer[position] != rune('R') {
								goto l1629
							}
							position++
						}
					l1633:
						{
							position1635, tokenIndex1635, depth1635 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1636
							}
							position++
							goto l1635
						l1636:
							position, tokenIndex, depth = position1635, tokenIndex1635, depth1635
							if buffer[position] != rune('U') {
								goto l1629
							}
							position++
						}
					l1635:
						{
							position1637, tokenIndex1637, depth1637 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1638
							}
							position++
							goto l1637
						l1638:
							position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
							if buffer[position] != rune('E') {
								goto l1629
							}
							position++
						}
					l1637:
						if !rules[ruleskip]() {
							goto l1629
						}
						depth--
						add(ruleTRUE, position1630)
					}
					goto l1628
				l1629:
					position, tokenIndex, depth = position1628, tokenIndex1628, depth1628
					{
						position1639 := position
						depth++
						{
							position1640, tokenIndex1640, depth1640 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1641
							}
							position++
							goto l1640
						l1641:
							position, tokenIndex, depth = position1640, tokenIndex1640, depth1640
							if buffer[position] != rune('F') {
								goto l1626
							}
							position++
						}
					l1640:
						{
							position1642, tokenIndex1642, depth1642 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1643
							}
							position++
							goto l1642
						l1643:
							position, tokenIndex, depth = position1642, tokenIndex1642, depth1642
							if buffer[position] != rune('A') {
								goto l1626
							}
							position++
						}
					l1642:
						{
							position1644, tokenIndex1644, depth1644 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1645
							}
							position++
							goto l1644
						l1645:
							position, tokenIndex, depth = position1644, tokenIndex1644, depth1644
							if buffer[position] != rune('L') {
								goto l1626
							}
							position++
						}
					l1644:
						{
							position1646, tokenIndex1646, depth1646 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1647
							}
							position++
							goto l1646
						l1647:
							position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
							if buffer[position] != rune('S') {
								goto l1626
							}
							position++
						}
					l1646:
						{
							position1648, tokenIndex1648, depth1648 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1649
							}
							position++
							goto l1648
						l1649:
							position, tokenIndex, depth = position1648, tokenIndex1648, depth1648
							if buffer[position] != rune('E') {
								goto l1626
							}
							position++
						}
					l1648:
						if !rules[ruleskip]() {
							goto l1626
						}
						depth--
						add(ruleFALSE, position1639)
					}
				}
			l1628:
				depth--
				add(rulebooleanLiteral, position1627)
			}
			return true
		l1626:
			position, tokenIndex, depth = position1626, tokenIndex1626, depth1626
			return false
		},
		/* 83 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 84 blankNodeLabel <- <('_' ':' (pnCharsU / [0-9]) (((pnCharsU / ((&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))* pnCharsU) / ([0-9] / '-'))? skip)> */
		nil,
		/* 85 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 86 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1653, tokenIndex1653, depth1653 := position, tokenIndex, depth
			{
				position1654 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1653
				}
				position++
			l1655:
				{
					position1656, tokenIndex1656, depth1656 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1656
					}
					goto l1655
				l1656:
					position, tokenIndex, depth = position1656, tokenIndex1656, depth1656
				}
				if buffer[position] != rune(')') {
					goto l1653
				}
				position++
				if !rules[ruleskip]() {
					goto l1653
				}
				depth--
				add(rulenil, position1654)
			}
			return true
		l1653:
			position, tokenIndex, depth = position1653, tokenIndex1653, depth1653
			return false
		},
		/* 87 VARNAME <- <((pnCharsU / [0-9]) (pnCharsU / ((&('â') [‿-⁀]) | (&('Ì' | 'Í') [̀-ͯ]) | (&('Â') '·') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))*)> */
		nil,
		/* 88 pnPrefix <- <(pnCharsBase pnChars*)> */
		func() bool {
			position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
			{
				position1659 := position
				depth++
				if !rules[rulepnCharsBase]() {
					goto l1658
				}
			l1660:
				{
					position1661, tokenIndex1661, depth1661 := position, tokenIndex, depth
					{
						position1662 := position
						depth++
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1661
								}
								position++
								break
							case '-':
								if buffer[position] != rune('-') {
									goto l1661
								}
								position++
								break
							default:
								if !rules[rulepnCharsU]() {
									goto l1661
								}
								break
							}
						}

						depth--
						add(rulepnChars, position1662)
					}
					goto l1660
				l1661:
					position, tokenIndex, depth = position1661, tokenIndex1661, depth1661
				}
				depth--
				add(rulepnPrefix, position1659)
			}
			return true
		l1658:
			position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
			return false
		},
		/* 89 pnLocal <- <((&('%' | '\\') plx) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&(':') ':') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))+> */
		nil,
		/* 90 pnChars <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))> */
		nil,
		/* 91 pnCharsU <- <(pnCharsBase / '_')> */
		func() bool {
			position1666, tokenIndex1666, depth1666 := position, tokenIndex, depth
			{
				position1667 := position
				depth++
				{
					position1668, tokenIndex1668, depth1668 := position, tokenIndex, depth
					if !rules[rulepnCharsBase]() {
						goto l1669
					}
					goto l1668
				l1669:
					position, tokenIndex, depth = position1668, tokenIndex1668, depth1668
					if buffer[position] != rune('_') {
						goto l1666
					}
					position++
				}
			l1668:
				depth--
				add(rulepnCharsU, position1667)
			}
			return true
		l1666:
			position, tokenIndex, depth = position1666, tokenIndex1666, depth1666
			return false
		},
		/* 92 pnCharsBase <- <([À-Ö] / [Ø-ö] / [Ͱ-ͽ] / [‌-‍] / [⁰-↏] / [豈-﷏] / ((&('ð' | 'ñ' | 'ò' | 'ó') [𐀀-󯿿]) | (&('ï') [ﷰ-�]) | (&('ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í') [、-퟿]) | (&('â') [Ⰰ-⿯]) | (&('Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á') [Ϳ-῿]) | (&('Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë') [ø-˿]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))> */
		func() bool {
			position1670, tokenIndex1670, depth1670 := position, tokenIndex, depth
			{
				position1671 := position
				depth++
				{
					position1672, tokenIndex1672, depth1672 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('À') || c > rune('Ö') {
						goto l1673
					}
					position++
					goto l1672
				l1673:
					position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
					if c := buffer[position]; c < rune('Ø') || c > rune('ö') {
						goto l1674
					}
					position++
					goto l1672
				l1674:
					position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
					if c := buffer[position]; c < rune('Ͱ') || c > rune('ͽ') {
						goto l1675
					}
					position++
					goto l1672
				l1675:
					position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
					if c := buffer[position]; c < rune('\u200c') || c > rune('\u200d') {
						goto l1676
					}
					position++
					goto l1672
				l1676:
					position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
					if c := buffer[position]; c < rune('⁰') || c > rune('\u218f') {
						goto l1677
					}
					position++
					goto l1672
				l1677:
					position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
					if c := buffer[position]; c < rune('豈') || c > rune('\ufdcf') {
						goto l1678
					}
					position++
					goto l1672
				l1678:
					position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
					{
						switch buffer[position] {
						case 'ð', 'ñ', 'ò', 'ó':
							if c := buffer[position]; c < rune('𐀀') || c > rune('\U000effff') {
								goto l1670
							}
							position++
							break
						case 'ï':
							if c := buffer[position]; c < rune('ﷰ') || c > rune('�') {
								goto l1670
							}
							position++
							break
						case 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í':
							if c := buffer[position]; c < rune('、') || c > rune('\ud7ff') {
								goto l1670
							}
							position++
							break
						case 'â':
							if c := buffer[position]; c < rune('Ⰰ') || c > rune('\u2fef') {
								goto l1670
							}
							position++
							break
						case 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'Þ', 'ß', 'à', 'á':
							if c := buffer[position]; c < rune('\u037f') || c > rune('\u1fff') {
								goto l1670
							}
							position++
							break
						case 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ê', 'Ë':
							if c := buffer[position]; c < rune('ø') || c > rune('˿') {
								goto l1670
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1670
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1670
							}
							position++
							break
						}
					}

				}
			l1672:
				depth--
				add(rulepnCharsBase, position1671)
			}
			return true
		l1670:
			position, tokenIndex, depth = position1670, tokenIndex1670, depth1670
			return false
		},
		/* 93 plx <- <(percent / pnLocalEsc)> */
		nil,
		/* 94 percent <- <('%' hex hex)> */
		nil,
		/* 95 hex <- <((&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))> */
		func() bool {
			position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
			{
				position1683 := position
				depth++
				{
					switch buffer[position] {
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l1682
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1682
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1682
						}
						position++
						break
					}
				}

				depth--
				add(rulehex, position1683)
			}
			return true
		l1682:
			position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
			return false
		},
		/* 96 pnLocalEsc <- <('\\' ((&('%') '%') | (&('@') '@') | (&('#') '#') | (&('?') '?') | (&('/') '/') | (&('=') '=') | (&(';') ';') | (&(',') ',') | (&('+') '+') | (&('*') '*') | (&(')') ')') | (&('(') '(') | (&('\'') '\'') | (&('&') '&') | (&('$') '$') | (&('!') '!') | (&('-') '-') | (&('.') '.') | (&('~') '~') | (&('_') '_')))> */
		nil,
		/* 97 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 98 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 99 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 100 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 101 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 102 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 103 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1692, tokenIndex1692, depth1692 := position, tokenIndex, depth
			{
				position1693 := position
				depth++
				{
					position1694, tokenIndex1694, depth1694 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1695
					}
					position++
					goto l1694
				l1695:
					position, tokenIndex, depth = position1694, tokenIndex1694, depth1694
					if buffer[position] != rune('D') {
						goto l1692
					}
					position++
				}
			l1694:
				{
					position1696, tokenIndex1696, depth1696 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1697
					}
					position++
					goto l1696
				l1697:
					position, tokenIndex, depth = position1696, tokenIndex1696, depth1696
					if buffer[position] != rune('I') {
						goto l1692
					}
					position++
				}
			l1696:
				{
					position1698, tokenIndex1698, depth1698 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1699
					}
					position++
					goto l1698
				l1699:
					position, tokenIndex, depth = position1698, tokenIndex1698, depth1698
					if buffer[position] != rune('S') {
						goto l1692
					}
					position++
				}
			l1698:
				{
					position1700, tokenIndex1700, depth1700 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1701
					}
					position++
					goto l1700
				l1701:
					position, tokenIndex, depth = position1700, tokenIndex1700, depth1700
					if buffer[position] != rune('T') {
						goto l1692
					}
					position++
				}
			l1700:
				{
					position1702, tokenIndex1702, depth1702 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1703
					}
					position++
					goto l1702
				l1703:
					position, tokenIndex, depth = position1702, tokenIndex1702, depth1702
					if buffer[position] != rune('I') {
						goto l1692
					}
					position++
				}
			l1702:
				{
					position1704, tokenIndex1704, depth1704 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1705
					}
					position++
					goto l1704
				l1705:
					position, tokenIndex, depth = position1704, tokenIndex1704, depth1704
					if buffer[position] != rune('N') {
						goto l1692
					}
					position++
				}
			l1704:
				{
					position1706, tokenIndex1706, depth1706 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1707
					}
					position++
					goto l1706
				l1707:
					position, tokenIndex, depth = position1706, tokenIndex1706, depth1706
					if buffer[position] != rune('C') {
						goto l1692
					}
					position++
				}
			l1706:
				{
					position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1709
					}
					position++
					goto l1708
				l1709:
					position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
					if buffer[position] != rune('T') {
						goto l1692
					}
					position++
				}
			l1708:
				if !rules[ruleskip]() {
					goto l1692
				}
				depth--
				add(ruleDISTINCT, position1693)
			}
			return true
		l1692:
			position, tokenIndex, depth = position1692, tokenIndex1692, depth1692
			return false
		},
		/* 104 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 105 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 106 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 107 LBRACE <- <('{' skip)> */
		func() bool {
			position1713, tokenIndex1713, depth1713 := position, tokenIndex, depth
			{
				position1714 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1713
				}
				position++
				if !rules[ruleskip]() {
					goto l1713
				}
				depth--
				add(ruleLBRACE, position1714)
			}
			return true
		l1713:
			position, tokenIndex, depth = position1713, tokenIndex1713, depth1713
			return false
		},
		/* 108 RBRACE <- <('}' skip)> */
		func() bool {
			position1715, tokenIndex1715, depth1715 := position, tokenIndex, depth
			{
				position1716 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1715
				}
				position++
				if !rules[ruleskip]() {
					goto l1715
				}
				depth--
				add(ruleRBRACE, position1716)
			}
			return true
		l1715:
			position, tokenIndex, depth = position1715, tokenIndex1715, depth1715
			return false
		},
		/* 109 LBRACK <- <('[' skip)> */
		nil,
		/* 110 RBRACK <- <(']' skip)> */
		nil,
		/* 111 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1719, tokenIndex1719, depth1719 := position, tokenIndex, depth
			{
				position1720 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1719
				}
				position++
				if !rules[ruleskip]() {
					goto l1719
				}
				depth--
				add(ruleSEMICOLON, position1720)
			}
			return true
		l1719:
			position, tokenIndex, depth = position1719, tokenIndex1719, depth1719
			return false
		},
		/* 112 COMMA <- <(',' skip)> */
		func() bool {
			position1721, tokenIndex1721, depth1721 := position, tokenIndex, depth
			{
				position1722 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1721
				}
				position++
				if !rules[ruleskip]() {
					goto l1721
				}
				depth--
				add(ruleCOMMA, position1722)
			}
			return true
		l1721:
			position, tokenIndex, depth = position1721, tokenIndex1721, depth1721
			return false
		},
		/* 113 DOT <- <('.' skip)> */
		func() bool {
			position1723, tokenIndex1723, depth1723 := position, tokenIndex, depth
			{
				position1724 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1723
				}
				position++
				if !rules[ruleskip]() {
					goto l1723
				}
				depth--
				add(ruleDOT, position1724)
			}
			return true
		l1723:
			position, tokenIndex, depth = position1723, tokenIndex1723, depth1723
			return false
		},
		/* 114 COLON <- <(':' skip)> */
		nil,
		/* 115 PIPE <- <('|' skip)> */
		func() bool {
			position1726, tokenIndex1726, depth1726 := position, tokenIndex, depth
			{
				position1727 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1726
				}
				position++
				if !rules[ruleskip]() {
					goto l1726
				}
				depth--
				add(rulePIPE, position1727)
			}
			return true
		l1726:
			position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
			return false
		},
		/* 116 SLASH <- <('/' skip)> */
		func() bool {
			position1728, tokenIndex1728, depth1728 := position, tokenIndex, depth
			{
				position1729 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1728
				}
				position++
				if !rules[ruleskip]() {
					goto l1728
				}
				depth--
				add(ruleSLASH, position1729)
			}
			return true
		l1728:
			position, tokenIndex, depth = position1728, tokenIndex1728, depth1728
			return false
		},
		/* 117 INVERSE <- <('^' skip)> */
		func() bool {
			position1730, tokenIndex1730, depth1730 := position, tokenIndex, depth
			{
				position1731 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1730
				}
				position++
				if !rules[ruleskip]() {
					goto l1730
				}
				depth--
				add(ruleINVERSE, position1731)
			}
			return true
		l1730:
			position, tokenIndex, depth = position1730, tokenIndex1730, depth1730
			return false
		},
		/* 118 LPAREN <- <('(' skip)> */
		func() bool {
			position1732, tokenIndex1732, depth1732 := position, tokenIndex, depth
			{
				position1733 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1732
				}
				position++
				if !rules[ruleskip]() {
					goto l1732
				}
				depth--
				add(ruleLPAREN, position1733)
			}
			return true
		l1732:
			position, tokenIndex, depth = position1732, tokenIndex1732, depth1732
			return false
		},
		/* 119 RPAREN <- <(')' skip)> */
		func() bool {
			position1734, tokenIndex1734, depth1734 := position, tokenIndex, depth
			{
				position1735 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1734
				}
				position++
				if !rules[ruleskip]() {
					goto l1734
				}
				depth--
				add(ruleRPAREN, position1735)
			}
			return true
		l1734:
			position, tokenIndex, depth = position1734, tokenIndex1734, depth1734
			return false
		},
		/* 120 ISA <- <('a' skip)> */
		func() bool {
			position1736, tokenIndex1736, depth1736 := position, tokenIndex, depth
			{
				position1737 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1736
				}
				position++
				if !rules[ruleskip]() {
					goto l1736
				}
				depth--
				add(ruleISA, position1737)
			}
			return true
		l1736:
			position, tokenIndex, depth = position1736, tokenIndex1736, depth1736
			return false
		},
		/* 121 NOT <- <('!' skip)> */
		func() bool {
			position1738, tokenIndex1738, depth1738 := position, tokenIndex, depth
			{
				position1739 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1738
				}
				position++
				if !rules[ruleskip]() {
					goto l1738
				}
				depth--
				add(ruleNOT, position1739)
			}
			return true
		l1738:
			position, tokenIndex, depth = position1738, tokenIndex1738, depth1738
			return false
		},
		/* 122 STAR <- <('*' skip)> */
		func() bool {
			position1740, tokenIndex1740, depth1740 := position, tokenIndex, depth
			{
				position1741 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1740
				}
				position++
				if !rules[ruleskip]() {
					goto l1740
				}
				depth--
				add(ruleSTAR, position1741)
			}
			return true
		l1740:
			position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
			return false
		},
		/* 123 QUESTION <- <('?' skip)> */
		nil,
		/* 124 PLUS <- <('+' skip)> */
		func() bool {
			position1743, tokenIndex1743, depth1743 := position, tokenIndex, depth
			{
				position1744 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1743
				}
				position++
				if !rules[ruleskip]() {
					goto l1743
				}
				depth--
				add(rulePLUS, position1744)
			}
			return true
		l1743:
			position, tokenIndex, depth = position1743, tokenIndex1743, depth1743
			return false
		},
		/* 125 MINUS <- <('-' skip)> */
		func() bool {
			position1745, tokenIndex1745, depth1745 := position, tokenIndex, depth
			{
				position1746 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1745
				}
				position++
				if !rules[ruleskip]() {
					goto l1745
				}
				depth--
				add(ruleMINUS, position1746)
			}
			return true
		l1745:
			position, tokenIndex, depth = position1745, tokenIndex1745, depth1745
			return false
		},
		/* 126 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 127 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 128 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 129 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 130 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1751, tokenIndex1751, depth1751 := position, tokenIndex, depth
			{
				position1752 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1751
				}
				position++
			l1753:
				{
					position1754, tokenIndex1754, depth1754 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1754
					}
					position++
					goto l1753
				l1754:
					position, tokenIndex, depth = position1754, tokenIndex1754, depth1754
				}
				if !rules[ruleskip]() {
					goto l1751
				}
				depth--
				add(ruleINTEGER, position1752)
			}
			return true
		l1751:
			position, tokenIndex, depth = position1751, tokenIndex1751, depth1751
			return false
		},
		/* 131 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 132 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 133 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 134 OR <- <('|' '|' skip)> */
		nil,
		/* 135 AND <- <('&' '&' skip)> */
		nil,
		/* 136 EQ <- <('=' skip)> */
		func() bool {
			position1760, tokenIndex1760, depth1760 := position, tokenIndex, depth
			{
				position1761 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1760
				}
				position++
				if !rules[ruleskip]() {
					goto l1760
				}
				depth--
				add(ruleEQ, position1761)
			}
			return true
		l1760:
			position, tokenIndex, depth = position1760, tokenIndex1760, depth1760
			return false
		},
		/* 137 NE <- <('!' '=' skip)> */
		nil,
		/* 138 GT <- <('>' skip)> */
		nil,
		/* 139 LT <- <('<' skip)> */
		nil,
		/* 140 LE <- <('<' '=' skip)> */
		nil,
		/* 141 GE <- <('>' '=' skip)> */
		nil,
		/* 142 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 143 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 144 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1769, tokenIndex1769, depth1769 := position, tokenIndex, depth
			{
				position1770 := position
				depth++
				{
					position1771, tokenIndex1771, depth1771 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1772
					}
					position++
					goto l1771
				l1772:
					position, tokenIndex, depth = position1771, tokenIndex1771, depth1771
					if buffer[position] != rune('A') {
						goto l1769
					}
					position++
				}
			l1771:
				{
					position1773, tokenIndex1773, depth1773 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1774
					}
					position++
					goto l1773
				l1774:
					position, tokenIndex, depth = position1773, tokenIndex1773, depth1773
					if buffer[position] != rune('S') {
						goto l1769
					}
					position++
				}
			l1773:
				if !rules[ruleskip]() {
					goto l1769
				}
				depth--
				add(ruleAS, position1770)
			}
			return true
		l1769:
			position, tokenIndex, depth = position1769, tokenIndex1769, depth1769
			return false
		},
		/* 145 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 146 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 147 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 148 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 149 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 150 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 151 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 152 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 153 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 154 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 155 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 156 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 157 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 158 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 159 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 160 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 161 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 162 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 163 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 164 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 165 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 166 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 167 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 168 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 169 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 170 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 171 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 172 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 173 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 174 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 175 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 176 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 177 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 178 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 179 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 180 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 181 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 182 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 183 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 184 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 185 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 186 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 187 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 188 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 189 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 190 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 191 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 192 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 193 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 194 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 195 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 196 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 197 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 198 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 199 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 200 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 201 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 202 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 203 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 204 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 205 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 206 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 207 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 208 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 209 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 210 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 211 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 212 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 213 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1843, tokenIndex1843, depth1843 := position, tokenIndex, depth
			{
				position1844 := position
				depth++
				{
					position1845, tokenIndex1845, depth1845 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1846
					}
					position++
					goto l1845
				l1846:
					position, tokenIndex, depth = position1845, tokenIndex1845, depth1845
					if buffer[position] != rune('B') {
						goto l1843
					}
					position++
				}
			l1845:
				{
					position1847, tokenIndex1847, depth1847 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1848
					}
					position++
					goto l1847
				l1848:
					position, tokenIndex, depth = position1847, tokenIndex1847, depth1847
					if buffer[position] != rune('Y') {
						goto l1843
					}
					position++
				}
			l1847:
				if !rules[ruleskip]() {
					goto l1843
				}
				depth--
				add(ruleBY, position1844)
			}
			return true
		l1843:
			position, tokenIndex, depth = position1843, tokenIndex1843, depth1843
			return false
		},
		/* 214 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 215 GRAPH <- <(('g' / 'G') ('r' / 'R') ('a' / 'A') ('p' / 'P') ('h' / 'H') skip)> */
		nil,
		/* 216 MINUSSETOPER <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('s' / 'S') skip)> */
		nil,
		/* 217 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1853 := position
				depth++
			l1854:
				{
					position1855, tokenIndex1855, depth1855 := position, tokenIndex, depth
					{
						position1856, tokenIndex1856, depth1856 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1857
						}
						goto l1856
					l1857:
						position, tokenIndex, depth = position1856, tokenIndex1856, depth1856
						{
							position1858 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1855
							}
							position++
						l1859:
							{
								position1860, tokenIndex1860, depth1860 := position, tokenIndex, depth
								{
									position1861, tokenIndex1861, depth1861 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1861
									}
									goto l1860
								l1861:
									position, tokenIndex, depth = position1861, tokenIndex1861, depth1861
								}
								if !matchDot() {
									goto l1860
								}
								goto l1859
							l1860:
								position, tokenIndex, depth = position1860, tokenIndex1860, depth1860
							}
							if !rules[ruleendOfLine]() {
								goto l1855
							}
							depth--
							add(rulecomment, position1858)
						}
					}
				l1856:
					goto l1854
				l1855:
					position, tokenIndex, depth = position1855, tokenIndex1855, depth1855
				}
				depth--
				add(ruleskip, position1853)
			}
			return true
		},
		/* 218 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1862, tokenIndex1862, depth1862 := position, tokenIndex, depth
			{
				position1863 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1862
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1862
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1862
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1862
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1862
						}
						break
					}
				}

				depth--
				add(rulews, position1863)
			}
			return true
		l1862:
			position, tokenIndex, depth = position1862, tokenIndex1862, depth1862
			return false
		},
		/* 219 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 220 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1866, tokenIndex1866, depth1866 := position, tokenIndex, depth
			{
				position1867 := position
				depth++
				{
					position1868, tokenIndex1868, depth1868 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1869
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1869
					}
					position++
					goto l1868
				l1869:
					position, tokenIndex, depth = position1868, tokenIndex1868, depth1868
					if buffer[position] != rune('\n') {
						goto l1870
					}
					position++
					goto l1868
				l1870:
					position, tokenIndex, depth = position1868, tokenIndex1868, depth1868
					if buffer[position] != rune('\r') {
						goto l1866
					}
					position++
				}
			l1868:
				depth--
				add(ruleendOfLine, position1867)
			}
			return true
		l1866:
			position, tokenIndex, depth = position1866, tokenIndex1866, depth1866
			return false
		},
	}
	p.rules = rules
}
