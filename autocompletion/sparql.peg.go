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
	rulepof
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
	"pof",
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
	rules  [238]func() bool
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
									{
										position22, tokenIndex22, depth22 := position, tokenIndex, depth
										if !rules[rulepnPrefix]() {
											goto l22
										}
										goto l23
									l22:
										position, tokenIndex, depth = position22, tokenIndex22, depth22
									}
								l23:
									{
										position24 := position
										depth++
										if buffer[position] != rune(':') {
											goto l6
										}
										position++
										if !rules[ruleskip]() {
											goto l6
										}
										depth--
										add(ruleCOLON, position24)
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
								position26 := position
								depth++
								{
									position27 := position
									depth++
									{
										position28, tokenIndex28, depth28 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l29
										}
										position++
										goto l28
									l29:
										position, tokenIndex, depth = position28, tokenIndex28, depth28
										if buffer[position] != rune('B') {
											goto l4
										}
										position++
									}
								l28:
									{
										position30, tokenIndex30, depth30 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l31
										}
										position++
										goto l30
									l31:
										position, tokenIndex, depth = position30, tokenIndex30, depth30
										if buffer[position] != rune('A') {
											goto l4
										}
										position++
									}
								l30:
									{
										position32, tokenIndex32, depth32 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l33
										}
										position++
										goto l32
									l33:
										position, tokenIndex, depth = position32, tokenIndex32, depth32
										if buffer[position] != rune('S') {
											goto l4
										}
										position++
									}
								l32:
									{
										position34, tokenIndex34, depth34 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l35
										}
										position++
										goto l34
									l35:
										position, tokenIndex, depth = position34, tokenIndex34, depth34
										if buffer[position] != rune('E') {
											goto l4
										}
										position++
									}
								l34:
									if !rules[ruleskip]() {
										goto l4
									}
									depth--
									add(ruleBASE, position27)
								}
								if !rules[ruleiri]() {
									goto l4
								}
								depth--
								add(rulebaseDecl, position26)
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
					position36 := position
					depth++
					{
						switch buffer[position] {
						case 'A', 'a':
							{
								position38 := position
								depth++
								{
									position39 := position
									depth++
									{
										position40, tokenIndex40, depth40 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l41
										}
										position++
										goto l40
									l41:
										position, tokenIndex, depth = position40, tokenIndex40, depth40
										if buffer[position] != rune('A') {
											goto l0
										}
										position++
									}
								l40:
									{
										position42, tokenIndex42, depth42 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l43
										}
										position++
										goto l42
									l43:
										position, tokenIndex, depth = position42, tokenIndex42, depth42
										if buffer[position] != rune('S') {
											goto l0
										}
										position++
									}
								l42:
									{
										position44, tokenIndex44, depth44 := position, tokenIndex, depth
										if buffer[position] != rune('k') {
											goto l45
										}
										position++
										goto l44
									l45:
										position, tokenIndex, depth = position44, tokenIndex44, depth44
										if buffer[position] != rune('K') {
											goto l0
										}
										position++
									}
								l44:
									if !rules[ruleskip]() {
										goto l0
									}
									depth--
									add(ruleASK, position39)
								}
							l46:
								{
									position47, tokenIndex47, depth47 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l47
									}
									goto l46
								l47:
									position, tokenIndex, depth = position47, tokenIndex47, depth47
								}
								if !rules[rulewhereClause]() {
									goto l0
								}
								depth--
								add(ruleaskQuery, position38)
							}
							break
						case 'D', 'd':
							{
								position48 := position
								depth++
								{
									position49 := position
									depth++
									{
										position50 := position
										depth++
										{
											position51, tokenIndex51, depth51 := position, tokenIndex, depth
											if buffer[position] != rune('d') {
												goto l52
											}
											position++
											goto l51
										l52:
											position, tokenIndex, depth = position51, tokenIndex51, depth51
											if buffer[position] != rune('D') {
												goto l0
											}
											position++
										}
									l51:
										{
											position53, tokenIndex53, depth53 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l54
											}
											position++
											goto l53
										l54:
											position, tokenIndex, depth = position53, tokenIndex53, depth53
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l53:
										{
											position55, tokenIndex55, depth55 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l56
											}
											position++
											goto l55
										l56:
											position, tokenIndex, depth = position55, tokenIndex55, depth55
											if buffer[position] != rune('S') {
												goto l0
											}
											position++
										}
									l55:
										{
											position57, tokenIndex57, depth57 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l58
											}
											position++
											goto l57
										l58:
											position, tokenIndex, depth = position57, tokenIndex57, depth57
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l57:
										{
											position59, tokenIndex59, depth59 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l60
											}
											position++
											goto l59
										l60:
											position, tokenIndex, depth = position59, tokenIndex59, depth59
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l59:
										{
											position61, tokenIndex61, depth61 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l62
											}
											position++
											goto l61
										l62:
											position, tokenIndex, depth = position61, tokenIndex61, depth61
											if buffer[position] != rune('I') {
												goto l0
											}
											position++
										}
									l61:
										{
											position63, tokenIndex63, depth63 := position, tokenIndex, depth
											if buffer[position] != rune('b') {
												goto l64
											}
											position++
											goto l63
										l64:
											position, tokenIndex, depth = position63, tokenIndex63, depth63
											if buffer[position] != rune('B') {
												goto l0
											}
											position++
										}
									l63:
										{
											position65, tokenIndex65, depth65 := position, tokenIndex, depth
											if buffer[position] != rune('e') {
												goto l66
											}
											position++
											goto l65
										l66:
											position, tokenIndex, depth = position65, tokenIndex65, depth65
											if buffer[position] != rune('E') {
												goto l0
											}
											position++
										}
									l65:
										if !rules[ruleskip]() {
											goto l0
										}
										depth--
										add(ruleDESCRIBE, position50)
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
									add(ruledescribe, position49)
								}
							l68:
								{
									position69, tokenIndex69, depth69 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l69
									}
									goto l68
								l69:
									position, tokenIndex, depth = position69, tokenIndex69, depth69
								}
								{
									position70, tokenIndex70, depth70 := position, tokenIndex, depth
									if !rules[rulewhereClause]() {
										goto l70
									}
									goto l71
								l70:
									position, tokenIndex, depth = position70, tokenIndex70, depth70
								}
							l71:
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruledescribeQuery, position48)
							}
							break
						case 'C', 'c':
							{
								position72 := position
								depth++
								{
									position73 := position
									depth++
									{
										position74 := position
										depth++
										{
											position75, tokenIndex75, depth75 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l76
											}
											position++
											goto l75
										l76:
											position, tokenIndex, depth = position75, tokenIndex75, depth75
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l75:
										{
											position77, tokenIndex77, depth77 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l78
											}
											position++
											goto l77
										l78:
											position, tokenIndex, depth = position77, tokenIndex77, depth77
											if buffer[position] != rune('O') {
												goto l0
											}
											position++
										}
									l77:
										{
											position79, tokenIndex79, depth79 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l80
											}
											position++
											goto l79
										l80:
											position, tokenIndex, depth = position79, tokenIndex79, depth79
											if buffer[position] != rune('N') {
												goto l0
											}
											position++
										}
									l79:
										{
											position81, tokenIndex81, depth81 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l82
											}
											position++
											goto l81
										l82:
											position, tokenIndex, depth = position81, tokenIndex81, depth81
											if buffer[position] != rune('S') {
												goto l0
											}
											position++
										}
									l81:
										{
											position83, tokenIndex83, depth83 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l84
											}
											position++
											goto l83
										l84:
											position, tokenIndex, depth = position83, tokenIndex83, depth83
											if buffer[position] != rune('T') {
												goto l0
											}
											position++
										}
									l83:
										{
											position85, tokenIndex85, depth85 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l86
											}
											position++
											goto l85
										l86:
											position, tokenIndex, depth = position85, tokenIndex85, depth85
											if buffer[position] != rune('R') {
												goto l0
											}
											position++
										}
									l85:
										{
											position87, tokenIndex87, depth87 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l88
											}
											position++
											goto l87
										l88:
											position, tokenIndex, depth = position87, tokenIndex87, depth87
											if buffer[position] != rune('U') {
												goto l0
											}
											position++
										}
									l87:
										{
											position89, tokenIndex89, depth89 := position, tokenIndex, depth
											if buffer[position] != rune('c') {
												goto l90
											}
											position++
											goto l89
										l90:
											position, tokenIndex, depth = position89, tokenIndex89, depth89
											if buffer[position] != rune('C') {
												goto l0
											}
											position++
										}
									l89:
										{
											position91, tokenIndex91, depth91 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l92
											}
											position++
											goto l91
										l92:
											position, tokenIndex, depth = position91, tokenIndex91, depth91
											if buffer[position] != rune('T') {
												goto l0
											}
											position++
										}
									l91:
										if !rules[ruleskip]() {
											goto l0
										}
										depth--
										add(ruleCONSTRUCT, position74)
									}
									if !rules[ruleLBRACE]() {
										goto l0
									}
									{
										position93, tokenIndex93, depth93 := position, tokenIndex, depth
										if !rules[ruletriplesBlock]() {
											goto l93
										}
										goto l94
									l93:
										position, tokenIndex, depth = position93, tokenIndex93, depth93
									}
								l94:
									if !rules[ruleRBRACE]() {
										goto l0
									}
									depth--
									add(ruleconstruct, position73)
								}
							l95:
								{
									position96, tokenIndex96, depth96 := position, tokenIndex, depth
									if !rules[ruledatasetClause]() {
										goto l96
									}
									goto l95
								l96:
									position, tokenIndex, depth = position96, tokenIndex96, depth96
								}
								if !rules[rulewhereClause]() {
									goto l0
								}
								if !rules[rulesolutionModifier]() {
									goto l0
								}
								depth--
								add(ruleconstructQuery, position72)
							}
							break
						default:
							{
								position97 := position
								depth++
								if !rules[ruleselect]() {
									goto l0
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
								add(ruleselectQuery, position97)
							}
							break
						}
					}

					depth--
					add(rulequery, position36)
				}
				{
					position100, tokenIndex100, depth100 := position, tokenIndex, depth
					if !matchDot() {
						goto l100
					}
					goto l0
				l100:
					position, tokenIndex, depth = position100, tokenIndex100, depth100
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
		/* 2 prefixDecl <- <(PREFIX <(pnPrefix? COLON iri)> Action0)> */
		nil,
		/* 3 baseDecl <- <(BASE iri)> */
		nil,
		/* 4 query <- <((&('A' | 'a') askQuery) | (&('D' | 'd') describeQuery) | (&('C' | 'c') constructQuery) | (&('S' | 's') selectQuery))> */
		nil,
		/* 5 selectQuery <- <(select datasetClause* whereClause solutionModifier)> */
		nil,
		/* 6 select <- <(SELECT (DISTINCT / REDUCED)? (STAR / projectionElem+))> */
		func() bool {
			position106, tokenIndex106, depth106 := position, tokenIndex, depth
			{
				position107 := position
				depth++
				{
					position108 := position
					depth++
					{
						position109, tokenIndex109, depth109 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l110
						}
						position++
						goto l109
					l110:
						position, tokenIndex, depth = position109, tokenIndex109, depth109
						if buffer[position] != rune('S') {
							goto l106
						}
						position++
					}
				l109:
					{
						position111, tokenIndex111, depth111 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l112
						}
						position++
						goto l111
					l112:
						position, tokenIndex, depth = position111, tokenIndex111, depth111
						if buffer[position] != rune('E') {
							goto l106
						}
						position++
					}
				l111:
					{
						position113, tokenIndex113, depth113 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l114
						}
						position++
						goto l113
					l114:
						position, tokenIndex, depth = position113, tokenIndex113, depth113
						if buffer[position] != rune('L') {
							goto l106
						}
						position++
					}
				l113:
					{
						position115, tokenIndex115, depth115 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l116
						}
						position++
						goto l115
					l116:
						position, tokenIndex, depth = position115, tokenIndex115, depth115
						if buffer[position] != rune('E') {
							goto l106
						}
						position++
					}
				l115:
					{
						position117, tokenIndex117, depth117 := position, tokenIndex, depth
						if buffer[position] != rune('c') {
							goto l118
						}
						position++
						goto l117
					l118:
						position, tokenIndex, depth = position117, tokenIndex117, depth117
						if buffer[position] != rune('C') {
							goto l106
						}
						position++
					}
				l117:
					{
						position119, tokenIndex119, depth119 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l120
						}
						position++
						goto l119
					l120:
						position, tokenIndex, depth = position119, tokenIndex119, depth119
						if buffer[position] != rune('T') {
							goto l106
						}
						position++
					}
				l119:
					if !rules[ruleskip]() {
						goto l106
					}
					depth--
					add(ruleSELECT, position108)
				}
				{
					position121, tokenIndex121, depth121 := position, tokenIndex, depth
					{
						position123, tokenIndex123, depth123 := position, tokenIndex, depth
						if !rules[ruleDISTINCT]() {
							goto l124
						}
						goto l123
					l124:
						position, tokenIndex, depth = position123, tokenIndex123, depth123
						{
							position125 := position
							depth++
							{
								position126, tokenIndex126, depth126 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l127
								}
								position++
								goto l126
							l127:
								position, tokenIndex, depth = position126, tokenIndex126, depth126
								if buffer[position] != rune('R') {
									goto l121
								}
								position++
							}
						l126:
							{
								position128, tokenIndex128, depth128 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l129
								}
								position++
								goto l128
							l129:
								position, tokenIndex, depth = position128, tokenIndex128, depth128
								if buffer[position] != rune('E') {
									goto l121
								}
								position++
							}
						l128:
							{
								position130, tokenIndex130, depth130 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l131
								}
								position++
								goto l130
							l131:
								position, tokenIndex, depth = position130, tokenIndex130, depth130
								if buffer[position] != rune('D') {
									goto l121
								}
								position++
							}
						l130:
							{
								position132, tokenIndex132, depth132 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l133
								}
								position++
								goto l132
							l133:
								position, tokenIndex, depth = position132, tokenIndex132, depth132
								if buffer[position] != rune('U') {
									goto l121
								}
								position++
							}
						l132:
							{
								position134, tokenIndex134, depth134 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l135
								}
								position++
								goto l134
							l135:
								position, tokenIndex, depth = position134, tokenIndex134, depth134
								if buffer[position] != rune('C') {
									goto l121
								}
								position++
							}
						l134:
							{
								position136, tokenIndex136, depth136 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l137
								}
								position++
								goto l136
							l137:
								position, tokenIndex, depth = position136, tokenIndex136, depth136
								if buffer[position] != rune('E') {
									goto l121
								}
								position++
							}
						l136:
							{
								position138, tokenIndex138, depth138 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l139
								}
								position++
								goto l138
							l139:
								position, tokenIndex, depth = position138, tokenIndex138, depth138
								if buffer[position] != rune('D') {
									goto l121
								}
								position++
							}
						l138:
							if !rules[ruleskip]() {
								goto l121
							}
							depth--
							add(ruleREDUCED, position125)
						}
					}
				l123:
					goto l122
				l121:
					position, tokenIndex, depth = position121, tokenIndex121, depth121
				}
			l122:
				{
					position140, tokenIndex140, depth140 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l141
					}
					goto l140
				l141:
					position, tokenIndex, depth = position140, tokenIndex140, depth140
					{
						position144 := position
						depth++
						{
							position145, tokenIndex145, depth145 := position, tokenIndex, depth
							if !rules[rulevar]() {
								goto l146
							}
							goto l145
						l146:
							position, tokenIndex, depth = position145, tokenIndex145, depth145
							if !rules[ruleLPAREN]() {
								goto l106
							}
							if !rules[ruleexpression]() {
								goto l106
							}
							if !rules[ruleAS]() {
								goto l106
							}
							if !rules[rulevar]() {
								goto l106
							}
							if !rules[ruleRPAREN]() {
								goto l106
							}
						}
					l145:
						depth--
						add(ruleprojectionElem, position144)
					}
				l142:
					{
						position143, tokenIndex143, depth143 := position, tokenIndex, depth
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
									goto l143
								}
								if !rules[ruleexpression]() {
									goto l143
								}
								if !rules[ruleAS]() {
									goto l143
								}
								if !rules[rulevar]() {
									goto l143
								}
								if !rules[ruleRPAREN]() {
									goto l143
								}
							}
						l148:
							depth--
							add(ruleprojectionElem, position147)
						}
						goto l142
					l143:
						position, tokenIndex, depth = position143, tokenIndex143, depth143
					}
				}
			l140:
				depth--
				add(ruleselect, position107)
			}
			return true
		l106:
			position, tokenIndex, depth = position106, tokenIndex106, depth106
			return false
		},
		/* 7 subSelect <- <(select whereClause solutionModifier)> */
		func() bool {
			position150, tokenIndex150, depth150 := position, tokenIndex, depth
			{
				position151 := position
				depth++
				if !rules[ruleselect]() {
					goto l150
				}
				if !rules[rulewhereClause]() {
					goto l150
				}
				if !rules[rulesolutionModifier]() {
					goto l150
				}
				depth--
				add(rulesubSelect, position151)
			}
			return true
		l150:
			position, tokenIndex, depth = position150, tokenIndex150, depth150
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
			position158, tokenIndex158, depth158 := position, tokenIndex, depth
			{
				position159 := position
				depth++
				{
					position160 := position
					depth++
					{
						position161, tokenIndex161, depth161 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l162
						}
						position++
						goto l161
					l162:
						position, tokenIndex, depth = position161, tokenIndex161, depth161
						if buffer[position] != rune('F') {
							goto l158
						}
						position++
					}
				l161:
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l164
						}
						position++
						goto l163
					l164:
						position, tokenIndex, depth = position163, tokenIndex163, depth163
						if buffer[position] != rune('R') {
							goto l158
						}
						position++
					}
				l163:
					{
						position165, tokenIndex165, depth165 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l166
						}
						position++
						goto l165
					l166:
						position, tokenIndex, depth = position165, tokenIndex165, depth165
						if buffer[position] != rune('O') {
							goto l158
						}
						position++
					}
				l165:
					{
						position167, tokenIndex167, depth167 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l168
						}
						position++
						goto l167
					l168:
						position, tokenIndex, depth = position167, tokenIndex167, depth167
						if buffer[position] != rune('M') {
							goto l158
						}
						position++
					}
				l167:
					if !rules[ruleskip]() {
						goto l158
					}
					depth--
					add(ruleFROM, position160)
				}
				{
					position169, tokenIndex169, depth169 := position, tokenIndex, depth
					{
						position171 := position
						depth++
						{
							position172, tokenIndex172, depth172 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l173
							}
							position++
							goto l172
						l173:
							position, tokenIndex, depth = position172, tokenIndex172, depth172
							if buffer[position] != rune('N') {
								goto l169
							}
							position++
						}
					l172:
						{
							position174, tokenIndex174, depth174 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l175
							}
							position++
							goto l174
						l175:
							position, tokenIndex, depth = position174, tokenIndex174, depth174
							if buffer[position] != rune('A') {
								goto l169
							}
							position++
						}
					l174:
						{
							position176, tokenIndex176, depth176 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l177
							}
							position++
							goto l176
						l177:
							position, tokenIndex, depth = position176, tokenIndex176, depth176
							if buffer[position] != rune('M') {
								goto l169
							}
							position++
						}
					l176:
						{
							position178, tokenIndex178, depth178 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l179
							}
							position++
							goto l178
						l179:
							position, tokenIndex, depth = position178, tokenIndex178, depth178
							if buffer[position] != rune('E') {
								goto l169
							}
							position++
						}
					l178:
						{
							position180, tokenIndex180, depth180 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l181
							}
							position++
							goto l180
						l181:
							position, tokenIndex, depth = position180, tokenIndex180, depth180
							if buffer[position] != rune('D') {
								goto l169
							}
							position++
						}
					l180:
						if !rules[ruleskip]() {
							goto l169
						}
						depth--
						add(ruleNAMED, position171)
					}
					goto l170
				l169:
					position, tokenIndex, depth = position169, tokenIndex169, depth169
				}
			l170:
				if !rules[ruleiriref]() {
					goto l158
				}
				depth--
				add(ruledatasetClause, position159)
			}
			return true
		l158:
			position, tokenIndex, depth = position158, tokenIndex158, depth158
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position182, tokenIndex182, depth182 := position, tokenIndex, depth
			{
				position183 := position
				depth++
				{
					position184, tokenIndex184, depth184 := position, tokenIndex, depth
					{
						position186 := position
						depth++
						{
							position187, tokenIndex187, depth187 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l188
							}
							position++
							goto l187
						l188:
							position, tokenIndex, depth = position187, tokenIndex187, depth187
							if buffer[position] != rune('W') {
								goto l184
							}
							position++
						}
					l187:
						{
							position189, tokenIndex189, depth189 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l190
							}
							position++
							goto l189
						l190:
							position, tokenIndex, depth = position189, tokenIndex189, depth189
							if buffer[position] != rune('H') {
								goto l184
							}
							position++
						}
					l189:
						{
							position191, tokenIndex191, depth191 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l192
							}
							position++
							goto l191
						l192:
							position, tokenIndex, depth = position191, tokenIndex191, depth191
							if buffer[position] != rune('E') {
								goto l184
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
								goto l184
							}
							position++
						}
					l193:
						{
							position195, tokenIndex195, depth195 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l196
							}
							position++
							goto l195
						l196:
							position, tokenIndex, depth = position195, tokenIndex195, depth195
							if buffer[position] != rune('E') {
								goto l184
							}
							position++
						}
					l195:
						if !rules[ruleskip]() {
							goto l184
						}
						depth--
						add(ruleWHERE, position186)
					}
					goto l185
				l184:
					position, tokenIndex, depth = position184, tokenIndex184, depth184
				}
			l185:
				if !rules[rulegroupGraphPattern]() {
					goto l182
				}
				depth--
				add(rulewhereClause, position183)
			}
			return true
		l182:
			position, tokenIndex, depth = position182, tokenIndex182, depth182
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position197, tokenIndex197, depth197 := position, tokenIndex, depth
			{
				position198 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l197
				}
				{
					position199, tokenIndex199, depth199 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l200
					}
					goto l199
				l200:
					position, tokenIndex, depth = position199, tokenIndex199, depth199
					if !rules[rulegraphPattern]() {
						goto l197
					}
				}
			l199:
				if !rules[ruleRBRACE]() {
					goto l197
				}
				depth--
				add(rulegroupGraphPattern, position198)
			}
			return true
		l197:
			position, tokenIndex, depth = position197, tokenIndex197, depth197
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position202 := position
				depth++
				{
					position203, tokenIndex203, depth203 := position, tokenIndex, depth
					{
						position205 := position
						depth++
						{
							position206, tokenIndex206, depth206 := position, tokenIndex, depth
							if !rules[ruletriplesBlock]() {
								goto l207
							}
						l208:
							{
								position209, tokenIndex209, depth209 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l209
								}
								{
									position210, tokenIndex210, depth210 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l210
									}
									goto l211
								l210:
									position, tokenIndex, depth = position210, tokenIndex210, depth210
								}
							l211:
								{
									position212, tokenIndex212, depth212 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l212
									}
									goto l213
								l212:
									position, tokenIndex, depth = position212, tokenIndex212, depth212
								}
							l213:
								goto l208
							l209:
								position, tokenIndex, depth = position209, tokenIndex209, depth209
							}
							goto l206
						l207:
							position, tokenIndex, depth = position206, tokenIndex206, depth206
							if !rules[rulefilterOrBind]() {
								goto l203
							}
							{
								position216, tokenIndex216, depth216 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l216
								}
								goto l217
							l216:
								position, tokenIndex, depth = position216, tokenIndex216, depth216
							}
						l217:
							{
								position218, tokenIndex218, depth218 := position, tokenIndex, depth
								if !rules[ruletriplesBlock]() {
									goto l218
								}
								goto l219
							l218:
								position, tokenIndex, depth = position218, tokenIndex218, depth218
							}
						l219:
						l214:
							{
								position215, tokenIndex215, depth215 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l215
								}
								{
									position220, tokenIndex220, depth220 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l220
									}
									goto l221
								l220:
									position, tokenIndex, depth = position220, tokenIndex220, depth220
								}
							l221:
								{
									position222, tokenIndex222, depth222 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l222
									}
									goto l223
								l222:
									position, tokenIndex, depth = position222, tokenIndex222, depth222
								}
							l223:
								goto l214
							l215:
								position, tokenIndex, depth = position215, tokenIndex215, depth215
							}
						}
					l206:
						depth--
						add(rulebasicGraphPattern, position205)
					}
					goto l204
				l203:
					position, tokenIndex, depth = position203, tokenIndex203, depth203
				}
			l204:
				{
					position224, tokenIndex224, depth224 := position, tokenIndex, depth
					{
						position226 := position
						depth++
						{
							switch buffer[position] {
							case 'M', 'm':
								{
									position228 := position
									depth++
									{
										position229 := position
										depth++
										{
											position230, tokenIndex230, depth230 := position, tokenIndex, depth
											if buffer[position] != rune('m') {
												goto l231
											}
											position++
											goto l230
										l231:
											position, tokenIndex, depth = position230, tokenIndex230, depth230
											if buffer[position] != rune('M') {
												goto l224
											}
											position++
										}
									l230:
										{
											position232, tokenIndex232, depth232 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l233
											}
											position++
											goto l232
										l233:
											position, tokenIndex, depth = position232, tokenIndex232, depth232
											if buffer[position] != rune('I') {
												goto l224
											}
											position++
										}
									l232:
										{
											position234, tokenIndex234, depth234 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l235
											}
											position++
											goto l234
										l235:
											position, tokenIndex, depth = position234, tokenIndex234, depth234
											if buffer[position] != rune('N') {
												goto l224
											}
											position++
										}
									l234:
										{
											position236, tokenIndex236, depth236 := position, tokenIndex, depth
											if buffer[position] != rune('u') {
												goto l237
											}
											position++
											goto l236
										l237:
											position, tokenIndex, depth = position236, tokenIndex236, depth236
											if buffer[position] != rune('U') {
												goto l224
											}
											position++
										}
									l236:
										{
											position238, tokenIndex238, depth238 := position, tokenIndex, depth
											if buffer[position] != rune('s') {
												goto l239
											}
											position++
											goto l238
										l239:
											position, tokenIndex, depth = position238, tokenIndex238, depth238
											if buffer[position] != rune('S') {
												goto l224
											}
											position++
										}
									l238:
										if !rules[ruleskip]() {
											goto l224
										}
										depth--
										add(ruleMINUSSETOPER, position229)
									}
									if !rules[rulegroupGraphPattern]() {
										goto l224
									}
									depth--
									add(ruleminusGraphPattern, position228)
								}
								break
							case 'G', 'g':
								{
									position240 := position
									depth++
									{
										position241 := position
										depth++
										{
											position242, tokenIndex242, depth242 := position, tokenIndex, depth
											if buffer[position] != rune('g') {
												goto l243
											}
											position++
											goto l242
										l243:
											position, tokenIndex, depth = position242, tokenIndex242, depth242
											if buffer[position] != rune('G') {
												goto l224
											}
											position++
										}
									l242:
										{
											position244, tokenIndex244, depth244 := position, tokenIndex, depth
											if buffer[position] != rune('r') {
												goto l245
											}
											position++
											goto l244
										l245:
											position, tokenIndex, depth = position244, tokenIndex244, depth244
											if buffer[position] != rune('R') {
												goto l224
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
												goto l224
											}
											position++
										}
									l246:
										{
											position248, tokenIndex248, depth248 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l249
											}
											position++
											goto l248
										l249:
											position, tokenIndex, depth = position248, tokenIndex248, depth248
											if buffer[position] != rune('P') {
												goto l224
											}
											position++
										}
									l248:
										{
											position250, tokenIndex250, depth250 := position, tokenIndex, depth
											if buffer[position] != rune('h') {
												goto l251
											}
											position++
											goto l250
										l251:
											position, tokenIndex, depth = position250, tokenIndex250, depth250
											if buffer[position] != rune('H') {
												goto l224
											}
											position++
										}
									l250:
										if !rules[ruleskip]() {
											goto l224
										}
										depth--
										add(ruleGRAPH, position241)
									}
									{
										position252, tokenIndex252, depth252 := position, tokenIndex, depth
										if !rules[rulevar]() {
											goto l253
										}
										goto l252
									l253:
										position, tokenIndex, depth = position252, tokenIndex252, depth252
										if !rules[ruleiriref]() {
											goto l224
										}
									}
								l252:
									if !rules[rulegroupGraphPattern]() {
										goto l224
									}
									depth--
									add(rulegraphGraphPattern, position240)
								}
								break
							case '{':
								if !rules[rulegroupOrUnionGraphPattern]() {
									goto l224
								}
								break
							default:
								{
									position254 := position
									depth++
									{
										position255 := position
										depth++
										{
											position256, tokenIndex256, depth256 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l257
											}
											position++
											goto l256
										l257:
											position, tokenIndex, depth = position256, tokenIndex256, depth256
											if buffer[position] != rune('O') {
												goto l224
											}
											position++
										}
									l256:
										{
											position258, tokenIndex258, depth258 := position, tokenIndex, depth
											if buffer[position] != rune('p') {
												goto l259
											}
											position++
											goto l258
										l259:
											position, tokenIndex, depth = position258, tokenIndex258, depth258
											if buffer[position] != rune('P') {
												goto l224
											}
											position++
										}
									l258:
										{
											position260, tokenIndex260, depth260 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l261
											}
											position++
											goto l260
										l261:
											position, tokenIndex, depth = position260, tokenIndex260, depth260
											if buffer[position] != rune('T') {
												goto l224
											}
											position++
										}
									l260:
										{
											position262, tokenIndex262, depth262 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l263
											}
											position++
											goto l262
										l263:
											position, tokenIndex, depth = position262, tokenIndex262, depth262
											if buffer[position] != rune('I') {
												goto l224
											}
											position++
										}
									l262:
										{
											position264, tokenIndex264, depth264 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l265
											}
											position++
											goto l264
										l265:
											position, tokenIndex, depth = position264, tokenIndex264, depth264
											if buffer[position] != rune('O') {
												goto l224
											}
											position++
										}
									l264:
										{
											position266, tokenIndex266, depth266 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l267
											}
											position++
											goto l266
										l267:
											position, tokenIndex, depth = position266, tokenIndex266, depth266
											if buffer[position] != rune('N') {
												goto l224
											}
											position++
										}
									l266:
										{
											position268, tokenIndex268, depth268 := position, tokenIndex, depth
											if buffer[position] != rune('a') {
												goto l269
											}
											position++
											goto l268
										l269:
											position, tokenIndex, depth = position268, tokenIndex268, depth268
											if buffer[position] != rune('A') {
												goto l224
											}
											position++
										}
									l268:
										{
											position270, tokenIndex270, depth270 := position, tokenIndex, depth
											if buffer[position] != rune('l') {
												goto l271
											}
											position++
											goto l270
										l271:
											position, tokenIndex, depth = position270, tokenIndex270, depth270
											if buffer[position] != rune('L') {
												goto l224
											}
											position++
										}
									l270:
										if !rules[ruleskip]() {
											goto l224
										}
										depth--
										add(ruleOPTIONAL, position255)
									}
									if !rules[ruleLBRACE]() {
										goto l224
									}
									{
										position272, tokenIndex272, depth272 := position, tokenIndex, depth
										if !rules[rulesubSelect]() {
											goto l273
										}
										goto l272
									l273:
										position, tokenIndex, depth = position272, tokenIndex272, depth272
										if !rules[rulegraphPattern]() {
											goto l224
										}
									}
								l272:
									if !rules[ruleRBRACE]() {
										goto l224
									}
									depth--
									add(ruleoptionalGraphPattern, position254)
								}
								break
							}
						}

						depth--
						add(rulegraphPatternNotTriples, position226)
					}
					{
						position274, tokenIndex274, depth274 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l274
						}
						goto l275
					l274:
						position, tokenIndex, depth = position274, tokenIndex274, depth274
					}
				l275:
					if !rules[rulegraphPattern]() {
						goto l224
					}
					goto l225
				l224:
					position, tokenIndex, depth = position224, tokenIndex224, depth224
				}
			l225:
				depth--
				add(rulegraphPattern, position202)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <((&('M' | 'm') minusGraphPattern) | (&('G' | 'g') graphGraphPattern) | (&('{') groupOrUnionGraphPattern) | (&('O' | 'o') optionalGraphPattern))> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position278, tokenIndex278, depth278 := position, tokenIndex, depth
			{
				position279 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l278
				}
				{
					position280, tokenIndex280, depth280 := position, tokenIndex, depth
					{
						position282 := position
						depth++
						{
							position283, tokenIndex283, depth283 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l284
							}
							position++
							goto l283
						l284:
							position, tokenIndex, depth = position283, tokenIndex283, depth283
							if buffer[position] != rune('U') {
								goto l280
							}
							position++
						}
					l283:
						{
							position285, tokenIndex285, depth285 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l286
							}
							position++
							goto l285
						l286:
							position, tokenIndex, depth = position285, tokenIndex285, depth285
							if buffer[position] != rune('N') {
								goto l280
							}
							position++
						}
					l285:
						{
							position287, tokenIndex287, depth287 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l288
							}
							position++
							goto l287
						l288:
							position, tokenIndex, depth = position287, tokenIndex287, depth287
							if buffer[position] != rune('I') {
								goto l280
							}
							position++
						}
					l287:
						{
							position289, tokenIndex289, depth289 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l290
							}
							position++
							goto l289
						l290:
							position, tokenIndex, depth = position289, tokenIndex289, depth289
							if buffer[position] != rune('O') {
								goto l280
							}
							position++
						}
					l289:
						{
							position291, tokenIndex291, depth291 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l292
							}
							position++
							goto l291
						l292:
							position, tokenIndex, depth = position291, tokenIndex291, depth291
							if buffer[position] != rune('N') {
								goto l280
							}
							position++
						}
					l291:
						if !rules[ruleskip]() {
							goto l280
						}
						depth--
						add(ruleUNION, position282)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l280
					}
					goto l281
				l280:
					position, tokenIndex, depth = position280, tokenIndex280, depth280
				}
			l281:
				depth--
				add(rulegroupOrUnionGraphPattern, position279)
			}
			return true
		l278:
			position, tokenIndex, depth = position278, tokenIndex278, depth278
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
			position296, tokenIndex296, depth296 := position, tokenIndex, depth
			{
				position297 := position
				depth++
				{
					position298, tokenIndex298, depth298 := position, tokenIndex, depth
					{
						position300 := position
						depth++
						{
							position301, tokenIndex301, depth301 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l302
							}
							position++
							goto l301
						l302:
							position, tokenIndex, depth = position301, tokenIndex301, depth301
							if buffer[position] != rune('F') {
								goto l299
							}
							position++
						}
					l301:
						{
							position303, tokenIndex303, depth303 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l304
							}
							position++
							goto l303
						l304:
							position, tokenIndex, depth = position303, tokenIndex303, depth303
							if buffer[position] != rune('I') {
								goto l299
							}
							position++
						}
					l303:
						{
							position305, tokenIndex305, depth305 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l306
							}
							position++
							goto l305
						l306:
							position, tokenIndex, depth = position305, tokenIndex305, depth305
							if buffer[position] != rune('L') {
								goto l299
							}
							position++
						}
					l305:
						{
							position307, tokenIndex307, depth307 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l308
							}
							position++
							goto l307
						l308:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if buffer[position] != rune('T') {
								goto l299
							}
							position++
						}
					l307:
						{
							position309, tokenIndex309, depth309 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l310
							}
							position++
							goto l309
						l310:
							position, tokenIndex, depth = position309, tokenIndex309, depth309
							if buffer[position] != rune('E') {
								goto l299
							}
							position++
						}
					l309:
						{
							position311, tokenIndex311, depth311 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l312
							}
							position++
							goto l311
						l312:
							position, tokenIndex, depth = position311, tokenIndex311, depth311
							if buffer[position] != rune('R') {
								goto l299
							}
							position++
						}
					l311:
						if !rules[ruleskip]() {
							goto l299
						}
						depth--
						add(ruleFILTER, position300)
					}
					if !rules[ruleconstraint]() {
						goto l299
					}
					goto l298
				l299:
					position, tokenIndex, depth = position298, tokenIndex298, depth298
					{
						position313 := position
						depth++
						{
							position314, tokenIndex314, depth314 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l315
							}
							position++
							goto l314
						l315:
							position, tokenIndex, depth = position314, tokenIndex314, depth314
							if buffer[position] != rune('B') {
								goto l296
							}
							position++
						}
					l314:
						{
							position316, tokenIndex316, depth316 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l317
							}
							position++
							goto l316
						l317:
							position, tokenIndex, depth = position316, tokenIndex316, depth316
							if buffer[position] != rune('I') {
								goto l296
							}
							position++
						}
					l316:
						{
							position318, tokenIndex318, depth318 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l319
							}
							position++
							goto l318
						l319:
							position, tokenIndex, depth = position318, tokenIndex318, depth318
							if buffer[position] != rune('N') {
								goto l296
							}
							position++
						}
					l318:
						{
							position320, tokenIndex320, depth320 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l321
							}
							position++
							goto l320
						l321:
							position, tokenIndex, depth = position320, tokenIndex320, depth320
							if buffer[position] != rune('D') {
								goto l296
							}
							position++
						}
					l320:
						if !rules[ruleskip]() {
							goto l296
						}
						depth--
						add(ruleBIND, position313)
					}
					if !rules[ruleLPAREN]() {
						goto l296
					}
					if !rules[ruleexpression]() {
						goto l296
					}
					if !rules[ruleAS]() {
						goto l296
					}
					if !rules[rulevar]() {
						goto l296
					}
					if !rules[ruleRPAREN]() {
						goto l296
					}
				}
			l298:
				depth--
				add(rulefilterOrBind, position297)
			}
			return true
		l296:
			position, tokenIndex, depth = position296, tokenIndex296, depth296
			return false
		},
		/* 25 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		func() bool {
			position322, tokenIndex322, depth322 := position, tokenIndex, depth
			{
				position323 := position
				depth++
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if !rules[rulebrackettedExpression]() {
						goto l325
					}
					goto l324
				l325:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if !rules[rulebuiltinCall]() {
						goto l326
					}
					goto l324
				l326:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
					if !rules[rulefunctionCall]() {
						goto l322
					}
				}
			l324:
				depth--
				add(ruleconstraint, position323)
			}
			return true
		l322:
			position, tokenIndex, depth = position322, tokenIndex322, depth322
			return false
		},
		/* 26 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position327, tokenIndex327, depth327 := position, tokenIndex, depth
			{
				position328 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l327
				}
			l329:
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l330
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l330
					}
					goto l329
				l330:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
				}
				{
					position331, tokenIndex331, depth331 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l331
					}
					goto l332
				l331:
					position, tokenIndex, depth = position331, tokenIndex331, depth331
				}
			l332:
				depth--
				add(ruletriplesBlock, position328)
			}
			return true
		l327:
			position, tokenIndex, depth = position327, tokenIndex327, depth327
			return false
		},
		/* 27 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath?))> */
		func() bool {
			position333, tokenIndex333, depth333 := position, tokenIndex, depth
			{
				position334 := position
				depth++
				{
					position335, tokenIndex335, depth335 := position, tokenIndex, depth
					{
						position337 := position
						depth++
						{
							position338, tokenIndex338, depth338 := position, tokenIndex, depth
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
								add(ruleAction1, position)
							}
							goto l338
						l339:
							position, tokenIndex, depth = position338, tokenIndex338, depth338
							{
								position343 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l342
								}
								depth--
								add(rulePegText, position343)
							}
							{
								add(ruleAction2, position)
							}
							goto l338
						l342:
							position, tokenIndex, depth = position338, tokenIndex338, depth338
							if !rules[rulepof]() {
								goto l336
							}
							{
								add(ruleAction3, position)
							}
						}
					l338:
						depth--
						add(rulevarOrTerm, position337)
					}
					if !rules[rulepropertyListPath]() {
						goto l336
					}
					goto l335
				l336:
					position, tokenIndex, depth = position335, tokenIndex335, depth335
					if !rules[ruletriplesNodePath]() {
						goto l333
					}
					{
						position346, tokenIndex346, depth346 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l346
						}
						goto l347
					l346:
						position, tokenIndex, depth = position346, tokenIndex346, depth346
					}
				l347:
				}
			l335:
				depth--
				add(ruletriplesSameSubjectPath, position334)
			}
			return true
		l333:
			position, tokenIndex, depth = position333, tokenIndex333, depth333
			return false
		},
		/* 28 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 29 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position349, tokenIndex349, depth349 := position, tokenIndex, depth
			{
				position350 := position
				depth++
				{
					position351, tokenIndex351, depth351 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l352
					}
					goto l351
				l352:
					position, tokenIndex, depth = position351, tokenIndex351, depth351
					{
						switch buffer[position] {
						case '(':
							if !rules[rulenil]() {
								goto l349
							}
							break
						case '[', '_':
							{
								position354 := position
								depth++
								{
									position355, tokenIndex355, depth355 := position, tokenIndex, depth
									{
										position357 := position
										depth++
										if buffer[position] != rune('_') {
											goto l356
										}
										position++
										if buffer[position] != rune(':') {
											goto l356
										}
										position++
										{
											position358, tokenIndex358, depth358 := position, tokenIndex, depth
											if !rules[rulepnCharsU]() {
												goto l359
											}
											goto l358
										l359:
											position, tokenIndex, depth = position358, tokenIndex358, depth358
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l356
											}
											position++
										}
									l358:
										{
											position360, tokenIndex360, depth360 := position, tokenIndex, depth
											{
												position362, tokenIndex362, depth362 := position, tokenIndex, depth
											l364:
												{
													position365, tokenIndex365, depth365 := position, tokenIndex, depth
													{
														position366, tokenIndex366, depth366 := position, tokenIndex, depth
														if !rules[rulepnCharsU]() {
															goto l367
														}
														goto l366
													l367:
														position, tokenIndex, depth = position366, tokenIndex366, depth366
														{
															switch buffer[position] {
															case '.':
																if buffer[position] != rune('.') {
																	goto l365
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l365
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l365
																}
																position++
																break
															}
														}

													}
												l366:
													goto l364
												l365:
													position, tokenIndex, depth = position365, tokenIndex365, depth365
												}
												if !rules[rulepnCharsU]() {
													goto l363
												}
												goto l362
											l363:
												position, tokenIndex, depth = position362, tokenIndex362, depth362
												{
													position369, tokenIndex369, depth369 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l370
													}
													position++
													goto l369
												l370:
													position, tokenIndex, depth = position369, tokenIndex369, depth369
													if buffer[position] != rune('-') {
														goto l360
													}
													position++
												}
											l369:
											}
										l362:
											goto l361
										l360:
											position, tokenIndex, depth = position360, tokenIndex360, depth360
										}
									l361:
										if !rules[ruleskip]() {
											goto l356
										}
										depth--
										add(ruleblankNodeLabel, position357)
									}
									goto l355
								l356:
									position, tokenIndex, depth = position355, tokenIndex355, depth355
									{
										position371 := position
										depth++
										if buffer[position] != rune('[') {
											goto l349
										}
										position++
									l372:
										{
											position373, tokenIndex373, depth373 := position, tokenIndex, depth
											if !rules[rulews]() {
												goto l373
											}
											goto l372
										l373:
											position, tokenIndex, depth = position373, tokenIndex373, depth373
										}
										if buffer[position] != rune(']') {
											goto l349
										}
										position++
										if !rules[ruleskip]() {
											goto l349
										}
										depth--
										add(ruleanon, position371)
									}
								}
							l355:
								depth--
								add(ruleblankNode, position354)
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
			l351:
				depth--
				add(rulegraphTerm, position350)
			}
			return true
		l349:
			position, tokenIndex, depth = position349, tokenIndex349, depth349
			return false
		},
		/* 30 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				{
					position376, tokenIndex376, depth376 := position, tokenIndex, depth
					{
						position378 := position
						depth++
						if !rules[ruleLPAREN]() {
							goto l377
						}
						if !rules[rulegraphNodePath]() {
							goto l377
						}
					l379:
						{
							position380, tokenIndex380, depth380 := position, tokenIndex, depth
							if !rules[rulegraphNodePath]() {
								goto l380
							}
							goto l379
						l380:
							position, tokenIndex, depth = position380, tokenIndex380, depth380
						}
						if !rules[ruleRPAREN]() {
							goto l377
						}
						depth--
						add(rulecollectionPath, position378)
					}
					goto l376
				l377:
					position, tokenIndex, depth = position376, tokenIndex376, depth376
					{
						position381 := position
						depth++
						{
							position382 := position
							depth++
							if buffer[position] != rune('[') {
								goto l374
							}
							position++
							if !rules[ruleskip]() {
								goto l374
							}
							depth--
							add(ruleLBRACK, position382)
						}
						if !rules[rulepropertyListPath]() {
							goto l374
						}
						{
							position383 := position
							depth++
							if buffer[position] != rune(']') {
								goto l374
							}
							position++
							if !rules[ruleskip]() {
								goto l374
							}
							depth--
							add(ruleRBRACK, position383)
						}
						depth--
						add(ruleblankNodePropertyListPath, position381)
					}
				}
			l376:
				depth--
				add(ruletriplesNodePath, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
			return false
		},
		/* 31 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 32 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 33 propertyListPath <- <(((pof Action4) / (<var> Action5) / verbPath) objectListPath (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position386, tokenIndex386, depth386 := position, tokenIndex, depth
			{
				position387 := position
				depth++
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l389
					}
					{
						add(ruleAction4, position)
					}
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					{
						position392 := position
						depth++
						if !rules[rulevar]() {
							goto l391
						}
						depth--
						add(rulePegText, position392)
					}
					{
						add(ruleAction5, position)
					}
					goto l388
				l391:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					{
						position394 := position
						depth++
						if !rules[rulepath]() {
							goto l386
						}
						depth--
						add(ruleverbPath, position394)
					}
				}
			l388:
				{
					position395 := position
					depth++
					if !rules[ruleobjectPath]() {
						goto l386
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
				add(rulepropertyListPath, position387)
			}
			return true
		l386:
			position, tokenIndex, depth = position386, tokenIndex386, depth386
			return false
		},
		/* 34 verbPath <- <path> */
		nil,
		/* 35 path <- <pathAlternative> */
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
		/* 36 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 37 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position409, tokenIndex409, depth409 := position, tokenIndex, depth
			{
				position410 := position
				depth++
				{
					position411 := position
					depth++
					{
						position412 := position
						depth++
						{
							position413, tokenIndex413, depth413 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l413
							}
							goto l414
						l413:
							position, tokenIndex, depth = position413, tokenIndex413, depth413
						}
					l414:
						{
							position415 := position
							depth++
							{
								position416, tokenIndex416, depth416 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l417
								}
								goto l416
							l417:
								position, tokenIndex, depth = position416, tokenIndex416, depth416
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l409
										}
										if !rules[rulepath]() {
											goto l409
										}
										if !rules[ruleRPAREN]() {
											goto l409
										}
										break
									case '!':
										if !rules[ruleNOT]() {
											goto l409
										}
										{
											position419 := position
											depth++
											{
												position420, tokenIndex420, depth420 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l421
												}
												goto l420
											l421:
												position, tokenIndex, depth = position420, tokenIndex420, depth420
												if !rules[ruleLPAREN]() {
													goto l409
												}
												{
													position422, tokenIndex422, depth422 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l422
													}
												l424:
													{
														position425, tokenIndex425, depth425 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l425
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l425
														}
														goto l424
													l425:
														position, tokenIndex, depth = position425, tokenIndex425, depth425
													}
													goto l423
												l422:
													position, tokenIndex, depth = position422, tokenIndex422, depth422
												}
											l423:
												if !rules[ruleRPAREN]() {
													goto l409
												}
											}
										l420:
											depth--
											add(rulepathNegatedPropertySet, position419)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l409
										}
										break
									}
								}

							}
						l416:
							depth--
							add(rulepathPrimary, position415)
						}
						{
							position426, tokenIndex426, depth426 := position, tokenIndex, depth
							{
								position428 := position
								depth++
								{
									switch buffer[position] {
									case '+':
										if !rules[rulePLUS]() {
											goto l426
										}
										break
									case '?':
										{
											position430 := position
											depth++
											if buffer[position] != rune('?') {
												goto l426
											}
											position++
											if !rules[ruleskip]() {
												goto l426
											}
											depth--
											add(ruleQUESTION, position430)
										}
										break
									default:
										if !rules[ruleSTAR]() {
											goto l426
										}
										break
									}
								}

								{
									position431, tokenIndex431, depth431 := position, tokenIndex, depth
									if !matchDot() {
										goto l431
									}
									goto l426
								l431:
									position, tokenIndex, depth = position431, tokenIndex431, depth431
								}
								depth--
								add(rulepathMod, position428)
							}
							goto l427
						l426:
							position, tokenIndex, depth = position426, tokenIndex426, depth426
						}
					l427:
						depth--
						add(rulepathElt, position412)
					}
					depth--
					add(rulePegText, position411)
				}
				{
					add(ruleAction6, position)
				}
			l433:
				{
					position434, tokenIndex434, depth434 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l434
					}
					if !rules[rulepathSequence]() {
						goto l434
					}
					goto l433
				l434:
					position, tokenIndex, depth = position434, tokenIndex434, depth434
				}
				depth--
				add(rulepathSequence, position410)
			}
			return true
		l409:
			position, tokenIndex, depth = position409, tokenIndex409, depth409
			return false
		},
		/* 38 pathElt <- <(INVERSE? pathPrimary pathMod?)> */
		nil,
		/* 39 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 40 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 41 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position438, tokenIndex438, depth438 := position, tokenIndex, depth
			{
				position439 := position
				depth++
				{
					position440, tokenIndex440, depth440 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l441
					}
					goto l440
				l441:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if !rules[ruleISA]() {
						goto l442
					}
					goto l440
				l442:
					position, tokenIndex, depth = position440, tokenIndex440, depth440
					if !rules[ruleINVERSE]() {
						goto l438
					}
					{
						position443, tokenIndex443, depth443 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l444
						}
						goto l443
					l444:
						position, tokenIndex, depth = position443, tokenIndex443, depth443
						if !rules[ruleISA]() {
							goto l438
						}
					}
				l443:
				}
			l440:
				depth--
				add(rulepathOneInPropertySet, position439)
			}
			return true
		l438:
			position, tokenIndex, depth = position438, tokenIndex438, depth438
			return false
		},
		/* 42 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !.)> */
		nil,
		/* 43 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 44 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		func() bool {
			{
				position448 := position
				depth++
				{
					position449, tokenIndex449, depth449 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l450
					}
					{
						add(ruleAction7, position)
					}
					goto l449
				l450:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					{
						position453 := position
						depth++
						if !rules[rulegraphNodePath]() {
							goto l452
						}
						depth--
						add(rulePegText, position453)
					}
					{
						add(ruleAction8, position)
					}
					goto l449
				l452:
					position, tokenIndex, depth = position449, tokenIndex449, depth449
					{
						add(ruleAction9, position)
					}
				}
			l449:
				depth--
				add(ruleobjectPath, position448)
			}
			return true
		},
		/* 45 graphNodePath <- <(var / graphTerm / triplesNodePath)> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				{
					position458, tokenIndex458, depth458 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l459
					}
					goto l458
				l459:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !rules[rulegraphTerm]() {
						goto l460
					}
					goto l458
				l460:
					position, tokenIndex, depth = position458, tokenIndex458, depth458
					if !rules[ruletriplesNodePath]() {
						goto l456
					}
				}
			l458:
				depth--
				add(rulegraphNodePath, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 46 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position462 := position
				depth++
				{
					position463, tokenIndex463, depth463 := position, tokenIndex, depth
					{
						position465, tokenIndex465, depth465 := position, tokenIndex, depth
						{
							position467 := position
							depth++
							{
								position468, tokenIndex468, depth468 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l469
								}
								position++
								goto l468
							l469:
								position, tokenIndex, depth = position468, tokenIndex468, depth468
								if buffer[position] != rune('O') {
									goto l466
								}
								position++
							}
						l468:
							{
								position470, tokenIndex470, depth470 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l471
								}
								position++
								goto l470
							l471:
								position, tokenIndex, depth = position470, tokenIndex470, depth470
								if buffer[position] != rune('R') {
									goto l466
								}
								position++
							}
						l470:
							{
								position472, tokenIndex472, depth472 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l473
								}
								position++
								goto l472
							l473:
								position, tokenIndex, depth = position472, tokenIndex472, depth472
								if buffer[position] != rune('D') {
									goto l466
								}
								position++
							}
						l472:
							{
								position474, tokenIndex474, depth474 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l475
								}
								position++
								goto l474
							l475:
								position, tokenIndex, depth = position474, tokenIndex474, depth474
								if buffer[position] != rune('E') {
									goto l466
								}
								position++
							}
						l474:
							{
								position476, tokenIndex476, depth476 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l477
								}
								position++
								goto l476
							l477:
								position, tokenIndex, depth = position476, tokenIndex476, depth476
								if buffer[position] != rune('R') {
									goto l466
								}
								position++
							}
						l476:
							if !rules[ruleskip]() {
								goto l466
							}
							depth--
							add(ruleORDER, position467)
						}
						if !rules[ruleBY]() {
							goto l466
						}
						{
							position480 := position
							depth++
							{
								position481, tokenIndex481, depth481 := position, tokenIndex, depth
								{
									position483, tokenIndex483, depth483 := position, tokenIndex, depth
									{
										position485, tokenIndex485, depth485 := position, tokenIndex, depth
										{
											position487 := position
											depth++
											{
												position488, tokenIndex488, depth488 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l489
												}
												position++
												goto l488
											l489:
												position, tokenIndex, depth = position488, tokenIndex488, depth488
												if buffer[position] != rune('A') {
													goto l486
												}
												position++
											}
										l488:
											{
												position490, tokenIndex490, depth490 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l491
												}
												position++
												goto l490
											l491:
												position, tokenIndex, depth = position490, tokenIndex490, depth490
												if buffer[position] != rune('S') {
													goto l486
												}
												position++
											}
										l490:
											{
												position492, tokenIndex492, depth492 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l493
												}
												position++
												goto l492
											l493:
												position, tokenIndex, depth = position492, tokenIndex492, depth492
												if buffer[position] != rune('C') {
													goto l486
												}
												position++
											}
										l492:
											if !rules[ruleskip]() {
												goto l486
											}
											depth--
											add(ruleASC, position487)
										}
										goto l485
									l486:
										position, tokenIndex, depth = position485, tokenIndex485, depth485
										{
											position494 := position
											depth++
											{
												position495, tokenIndex495, depth495 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l496
												}
												position++
												goto l495
											l496:
												position, tokenIndex, depth = position495, tokenIndex495, depth495
												if buffer[position] != rune('D') {
													goto l483
												}
												position++
											}
										l495:
											{
												position497, tokenIndex497, depth497 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l498
												}
												position++
												goto l497
											l498:
												position, tokenIndex, depth = position497, tokenIndex497, depth497
												if buffer[position] != rune('E') {
													goto l483
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
													goto l483
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
													goto l483
												}
												position++
											}
										l501:
											if !rules[ruleskip]() {
												goto l483
											}
											depth--
											add(ruleDESC, position494)
										}
									}
								l485:
									goto l484
								l483:
									position, tokenIndex, depth = position483, tokenIndex483, depth483
								}
							l484:
								if !rules[rulebrackettedExpression]() {
									goto l482
								}
								goto l481
							l482:
								position, tokenIndex, depth = position481, tokenIndex481, depth481
								if !rules[rulefunctionCall]() {
									goto l503
								}
								goto l481
							l503:
								position, tokenIndex, depth = position481, tokenIndex481, depth481
								if !rules[rulebuiltinCall]() {
									goto l504
								}
								goto l481
							l504:
								position, tokenIndex, depth = position481, tokenIndex481, depth481
								if !rules[rulevar]() {
									goto l466
								}
							}
						l481:
							depth--
							add(ruleorderCondition, position480)
						}
					l478:
						{
							position479, tokenIndex479, depth479 := position, tokenIndex, depth
							{
								position505 := position
								depth++
								{
									position506, tokenIndex506, depth506 := position, tokenIndex, depth
									{
										position508, tokenIndex508, depth508 := position, tokenIndex, depth
										{
											position510, tokenIndex510, depth510 := position, tokenIndex, depth
											{
												position512 := position
												depth++
												{
													position513, tokenIndex513, depth513 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l514
													}
													position++
													goto l513
												l514:
													position, tokenIndex, depth = position513, tokenIndex513, depth513
													if buffer[position] != rune('A') {
														goto l511
													}
													position++
												}
											l513:
												{
													position515, tokenIndex515, depth515 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l516
													}
													position++
													goto l515
												l516:
													position, tokenIndex, depth = position515, tokenIndex515, depth515
													if buffer[position] != rune('S') {
														goto l511
													}
													position++
												}
											l515:
												{
													position517, tokenIndex517, depth517 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l518
													}
													position++
													goto l517
												l518:
													position, tokenIndex, depth = position517, tokenIndex517, depth517
													if buffer[position] != rune('C') {
														goto l511
													}
													position++
												}
											l517:
												if !rules[ruleskip]() {
													goto l511
												}
												depth--
												add(ruleASC, position512)
											}
											goto l510
										l511:
											position, tokenIndex, depth = position510, tokenIndex510, depth510
											{
												position519 := position
												depth++
												{
													position520, tokenIndex520, depth520 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l521
													}
													position++
													goto l520
												l521:
													position, tokenIndex, depth = position520, tokenIndex520, depth520
													if buffer[position] != rune('D') {
														goto l508
													}
													position++
												}
											l520:
												{
													position522, tokenIndex522, depth522 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l523
													}
													position++
													goto l522
												l523:
													position, tokenIndex, depth = position522, tokenIndex522, depth522
													if buffer[position] != rune('E') {
														goto l508
													}
													position++
												}
											l522:
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
														goto l508
													}
													position++
												}
											l524:
												{
													position526, tokenIndex526, depth526 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l527
													}
													position++
													goto l526
												l527:
													position, tokenIndex, depth = position526, tokenIndex526, depth526
													if buffer[position] != rune('C') {
														goto l508
													}
													position++
												}
											l526:
												if !rules[ruleskip]() {
													goto l508
												}
												depth--
												add(ruleDESC, position519)
											}
										}
									l510:
										goto l509
									l508:
										position, tokenIndex, depth = position508, tokenIndex508, depth508
									}
								l509:
									if !rules[rulebrackettedExpression]() {
										goto l507
									}
									goto l506
								l507:
									position, tokenIndex, depth = position506, tokenIndex506, depth506
									if !rules[rulefunctionCall]() {
										goto l528
									}
									goto l506
								l528:
									position, tokenIndex, depth = position506, tokenIndex506, depth506
									if !rules[rulebuiltinCall]() {
										goto l529
									}
									goto l506
								l529:
									position, tokenIndex, depth = position506, tokenIndex506, depth506
									if !rules[rulevar]() {
										goto l479
									}
								}
							l506:
								depth--
								add(ruleorderCondition, position505)
							}
							goto l478
						l479:
							position, tokenIndex, depth = position479, tokenIndex479, depth479
						}
						goto l465
					l466:
						position, tokenIndex, depth = position465, tokenIndex465, depth465
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position531 := position
									depth++
									{
										position532, tokenIndex532, depth532 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l533
										}
										position++
										goto l532
									l533:
										position, tokenIndex, depth = position532, tokenIndex532, depth532
										if buffer[position] != rune('H') {
											goto l463
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
											goto l463
										}
										position++
									}
								l534:
									{
										position536, tokenIndex536, depth536 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l537
										}
										position++
										goto l536
									l537:
										position, tokenIndex, depth = position536, tokenIndex536, depth536
										if buffer[position] != rune('V') {
											goto l463
										}
										position++
									}
								l536:
									{
										position538, tokenIndex538, depth538 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l539
										}
										position++
										goto l538
									l539:
										position, tokenIndex, depth = position538, tokenIndex538, depth538
										if buffer[position] != rune('I') {
											goto l463
										}
										position++
									}
								l538:
									{
										position540, tokenIndex540, depth540 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l541
										}
										position++
										goto l540
									l541:
										position, tokenIndex, depth = position540, tokenIndex540, depth540
										if buffer[position] != rune('N') {
											goto l463
										}
										position++
									}
								l540:
									{
										position542, tokenIndex542, depth542 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l543
										}
										position++
										goto l542
									l543:
										position, tokenIndex, depth = position542, tokenIndex542, depth542
										if buffer[position] != rune('G') {
											goto l463
										}
										position++
									}
								l542:
									if !rules[ruleskip]() {
										goto l463
									}
									depth--
									add(ruleHAVING, position531)
								}
								if !rules[ruleconstraint]() {
									goto l463
								}
								break
							case 'G', 'g':
								{
									position544 := position
									depth++
									{
										position545, tokenIndex545, depth545 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l546
										}
										position++
										goto l545
									l546:
										position, tokenIndex, depth = position545, tokenIndex545, depth545
										if buffer[position] != rune('G') {
											goto l463
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
											goto l463
										}
										position++
									}
								l547:
									{
										position549, tokenIndex549, depth549 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l550
										}
										position++
										goto l549
									l550:
										position, tokenIndex, depth = position549, tokenIndex549, depth549
										if buffer[position] != rune('O') {
											goto l463
										}
										position++
									}
								l549:
									{
										position551, tokenIndex551, depth551 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l552
										}
										position++
										goto l551
									l552:
										position, tokenIndex, depth = position551, tokenIndex551, depth551
										if buffer[position] != rune('U') {
											goto l463
										}
										position++
									}
								l551:
									{
										position553, tokenIndex553, depth553 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l554
										}
										position++
										goto l553
									l554:
										position, tokenIndex, depth = position553, tokenIndex553, depth553
										if buffer[position] != rune('P') {
											goto l463
										}
										position++
									}
								l553:
									if !rules[ruleskip]() {
										goto l463
									}
									depth--
									add(ruleGROUP, position544)
								}
								if !rules[ruleBY]() {
									goto l463
								}
								{
									position557 := position
									depth++
									{
										position558, tokenIndex558, depth558 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l559
										}
										goto l558
									l559:
										position, tokenIndex, depth = position558, tokenIndex558, depth558
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l463
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l463
												}
												if !rules[ruleexpression]() {
													goto l463
												}
												{
													position561, tokenIndex561, depth561 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l561
													}
													if !rules[rulevar]() {
														goto l561
													}
													goto l562
												l561:
													position, tokenIndex, depth = position561, tokenIndex561, depth561
												}
											l562:
												if !rules[ruleRPAREN]() {
													goto l463
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l463
												}
												break
											}
										}

									}
								l558:
									depth--
									add(rulegroupCondition, position557)
								}
							l555:
								{
									position556, tokenIndex556, depth556 := position, tokenIndex, depth
									{
										position563 := position
										depth++
										{
											position564, tokenIndex564, depth564 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l565
											}
											goto l564
										l565:
											position, tokenIndex, depth = position564, tokenIndex564, depth564
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l556
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l556
													}
													if !rules[ruleexpression]() {
														goto l556
													}
													{
														position567, tokenIndex567, depth567 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l567
														}
														if !rules[rulevar]() {
															goto l567
														}
														goto l568
													l567:
														position, tokenIndex, depth = position567, tokenIndex567, depth567
													}
												l568:
													if !rules[ruleRPAREN]() {
														goto l556
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l556
													}
													break
												}
											}

										}
									l564:
										depth--
										add(rulegroupCondition, position563)
									}
									goto l555
								l556:
									position, tokenIndex, depth = position556, tokenIndex556, depth556
								}
								break
							default:
								{
									position569 := position
									depth++
									{
										position570, tokenIndex570, depth570 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l571
										}
										{
											position572, tokenIndex572, depth572 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l572
											}
											goto l573
										l572:
											position, tokenIndex, depth = position572, tokenIndex572, depth572
										}
									l573:
										goto l570
									l571:
										position, tokenIndex, depth = position570, tokenIndex570, depth570
										if !rules[ruleoffset]() {
											goto l463
										}
										{
											position574, tokenIndex574, depth574 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l574
											}
											goto l575
										l574:
											position, tokenIndex, depth = position574, tokenIndex574, depth574
										}
									l575:
									}
								l570:
									depth--
									add(rulelimitOffsetClauses, position569)
								}
								break
							}
						}

					}
				l465:
					goto l464
				l463:
					position, tokenIndex, depth = position463, tokenIndex463, depth463
				}
			l464:
				depth--
				add(rulesolutionModifier, position462)
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
			position579, tokenIndex579, depth579 := position, tokenIndex, depth
			{
				position580 := position
				depth++
				{
					position581 := position
					depth++
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
							goto l579
						}
						position++
					}
				l582:
					{
						position584, tokenIndex584, depth584 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l585
						}
						position++
						goto l584
					l585:
						position, tokenIndex, depth = position584, tokenIndex584, depth584
						if buffer[position] != rune('I') {
							goto l579
						}
						position++
					}
				l584:
					{
						position586, tokenIndex586, depth586 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l587
						}
						position++
						goto l586
					l587:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
						if buffer[position] != rune('M') {
							goto l579
						}
						position++
					}
				l586:
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
							goto l579
						}
						position++
					}
				l588:
					{
						position590, tokenIndex590, depth590 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l591
						}
						position++
						goto l590
					l591:
						position, tokenIndex, depth = position590, tokenIndex590, depth590
						if buffer[position] != rune('T') {
							goto l579
						}
						position++
					}
				l590:
					if !rules[ruleskip]() {
						goto l579
					}
					depth--
					add(ruleLIMIT, position581)
				}
				if !rules[ruleINTEGER]() {
					goto l579
				}
				depth--
				add(rulelimit, position580)
			}
			return true
		l579:
			position, tokenIndex, depth = position579, tokenIndex579, depth579
			return false
		},
		/* 51 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position592, tokenIndex592, depth592 := position, tokenIndex, depth
			{
				position593 := position
				depth++
				{
					position594 := position
					depth++
					{
						position595, tokenIndex595, depth595 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l596
						}
						position++
						goto l595
					l596:
						position, tokenIndex, depth = position595, tokenIndex595, depth595
						if buffer[position] != rune('O') {
							goto l592
						}
						position++
					}
				l595:
					{
						position597, tokenIndex597, depth597 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l598
						}
						position++
						goto l597
					l598:
						position, tokenIndex, depth = position597, tokenIndex597, depth597
						if buffer[position] != rune('F') {
							goto l592
						}
						position++
					}
				l597:
					{
						position599, tokenIndex599, depth599 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l600
						}
						position++
						goto l599
					l600:
						position, tokenIndex, depth = position599, tokenIndex599, depth599
						if buffer[position] != rune('F') {
							goto l592
						}
						position++
					}
				l599:
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
							goto l592
						}
						position++
					}
				l601:
					{
						position603, tokenIndex603, depth603 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l604
						}
						position++
						goto l603
					l604:
						position, tokenIndex, depth = position603, tokenIndex603, depth603
						if buffer[position] != rune('E') {
							goto l592
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
							goto l592
						}
						position++
					}
				l605:
					if !rules[ruleskip]() {
						goto l592
					}
					depth--
					add(ruleOFFSET, position594)
				}
				if !rules[ruleINTEGER]() {
					goto l592
				}
				depth--
				add(ruleoffset, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 52 expression <- <conditionalOrExpression> */
		func() bool {
			position607, tokenIndex607, depth607 := position, tokenIndex, depth
			{
				position608 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l607
				}
				depth--
				add(ruleexpression, position608)
			}
			return true
		l607:
			position, tokenIndex, depth = position607, tokenIndex607, depth607
			return false
		},
		/* 53 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position609, tokenIndex609, depth609 := position, tokenIndex, depth
			{
				position610 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l609
				}
				{
					position611, tokenIndex611, depth611 := position, tokenIndex, depth
					{
						position613 := position
						depth++
						if buffer[position] != rune('|') {
							goto l611
						}
						position++
						if buffer[position] != rune('|') {
							goto l611
						}
						position++
						if !rules[ruleskip]() {
							goto l611
						}
						depth--
						add(ruleOR, position613)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l611
					}
					goto l612
				l611:
					position, tokenIndex, depth = position611, tokenIndex611, depth611
				}
			l612:
				depth--
				add(ruleconditionalOrExpression, position610)
			}
			return true
		l609:
			position, tokenIndex, depth = position609, tokenIndex609, depth609
			return false
		},
		/* 54 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position614, tokenIndex614, depth614 := position, tokenIndex, depth
			{
				position615 := position
				depth++
				{
					position616 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l614
					}
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position620 := position
									depth++
									{
										position621 := position
										depth++
										{
											position622, tokenIndex622, depth622 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l623
											}
											position++
											goto l622
										l623:
											position, tokenIndex, depth = position622, tokenIndex622, depth622
											if buffer[position] != rune('N') {
												goto l617
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
												goto l617
											}
											position++
										}
									l624:
										{
											position626, tokenIndex626, depth626 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l627
											}
											position++
											goto l626
										l627:
											position, tokenIndex, depth = position626, tokenIndex626, depth626
											if buffer[position] != rune('T') {
												goto l617
											}
											position++
										}
									l626:
										if buffer[position] != rune(' ') {
											goto l617
										}
										position++
										{
											position628, tokenIndex628, depth628 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l629
											}
											position++
											goto l628
										l629:
											position, tokenIndex, depth = position628, tokenIndex628, depth628
											if buffer[position] != rune('I') {
												goto l617
											}
											position++
										}
									l628:
										{
											position630, tokenIndex630, depth630 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l631
											}
											position++
											goto l630
										l631:
											position, tokenIndex, depth = position630, tokenIndex630, depth630
											if buffer[position] != rune('N') {
												goto l617
											}
											position++
										}
									l630:
										if !rules[ruleskip]() {
											goto l617
										}
										depth--
										add(ruleNOTIN, position621)
									}
									if !rules[ruleargList]() {
										goto l617
									}
									depth--
									add(rulenotin, position620)
								}
								break
							case 'I', 'i':
								{
									position632 := position
									depth++
									{
										position633 := position
										depth++
										{
											position634, tokenIndex634, depth634 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l635
											}
											position++
											goto l634
										l635:
											position, tokenIndex, depth = position634, tokenIndex634, depth634
											if buffer[position] != rune('I') {
												goto l617
											}
											position++
										}
									l634:
										{
											position636, tokenIndex636, depth636 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l637
											}
											position++
											goto l636
										l637:
											position, tokenIndex, depth = position636, tokenIndex636, depth636
											if buffer[position] != rune('N') {
												goto l617
											}
											position++
										}
									l636:
										if !rules[ruleskip]() {
											goto l617
										}
										depth--
										add(ruleIN, position633)
									}
									if !rules[ruleargList]() {
										goto l617
									}
									depth--
									add(rulein, position632)
								}
								break
							default:
								{
									position638, tokenIndex638, depth638 := position, tokenIndex, depth
									{
										position640 := position
										depth++
										if buffer[position] != rune('<') {
											goto l639
										}
										position++
										if !rules[ruleskip]() {
											goto l639
										}
										depth--
										add(ruleLT, position640)
									}
									goto l638
								l639:
									position, tokenIndex, depth = position638, tokenIndex638, depth638
									{
										position642 := position
										depth++
										if buffer[position] != rune('>') {
											goto l641
										}
										position++
										if buffer[position] != rune('=') {
											goto l641
										}
										position++
										if !rules[ruleskip]() {
											goto l641
										}
										depth--
										add(ruleGE, position642)
									}
									goto l638
								l641:
									position, tokenIndex, depth = position638, tokenIndex638, depth638
									{
										switch buffer[position] {
										case '>':
											{
												position644 := position
												depth++
												if buffer[position] != rune('>') {
													goto l617
												}
												position++
												if !rules[ruleskip]() {
													goto l617
												}
												depth--
												add(ruleGT, position644)
											}
											break
										case '<':
											{
												position645 := position
												depth++
												if buffer[position] != rune('<') {
													goto l617
												}
												position++
												if buffer[position] != rune('=') {
													goto l617
												}
												position++
												if !rules[ruleskip]() {
													goto l617
												}
												depth--
												add(ruleLE, position645)
											}
											break
										case '!':
											{
												position646 := position
												depth++
												if buffer[position] != rune('!') {
													goto l617
												}
												position++
												if buffer[position] != rune('=') {
													goto l617
												}
												position++
												if !rules[ruleskip]() {
													goto l617
												}
												depth--
												add(ruleNE, position646)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l617
											}
											break
										}
									}

								}
							l638:
								if !rules[rulenumericExpression]() {
									goto l617
								}
								break
							}
						}

						goto l618
					l617:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
					}
				l618:
					depth--
					add(rulevalueLogical, position616)
				}
				{
					position647, tokenIndex647, depth647 := position, tokenIndex, depth
					{
						position649 := position
						depth++
						if buffer[position] != rune('&') {
							goto l647
						}
						position++
						if buffer[position] != rune('&') {
							goto l647
						}
						position++
						if !rules[ruleskip]() {
							goto l647
						}
						depth--
						add(ruleAND, position649)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l647
					}
					goto l648
				l647:
					position, tokenIndex, depth = position647, tokenIndex647, depth647
				}
			l648:
				depth--
				add(ruleconditionalAndExpression, position615)
			}
			return true
		l614:
			position, tokenIndex, depth = position614, tokenIndex614, depth614
			return false
		},
		/* 55 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 56 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position651, tokenIndex651, depth651 := position, tokenIndex, depth
			{
				position652 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l651
				}
			l653:
				{
					position654, tokenIndex654, depth654 := position, tokenIndex, depth
					{
						position655, tokenIndex655, depth655 := position, tokenIndex, depth
						{
							position657, tokenIndex657, depth657 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l658
							}
							goto l657
						l658:
							position, tokenIndex, depth = position657, tokenIndex657, depth657
							if !rules[ruleMINUS]() {
								goto l656
							}
						}
					l657:
						if !rules[rulemultiplicativeExpression]() {
							goto l656
						}
						goto l655
					l656:
						position, tokenIndex, depth = position655, tokenIndex655, depth655
						{
							position659 := position
							depth++
							{
								position660, tokenIndex660, depth660 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l661
								}
								position++
								goto l660
							l661:
								position, tokenIndex, depth = position660, tokenIndex660, depth660
								if buffer[position] != rune('-') {
									goto l654
								}
								position++
							}
						l660:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l654
							}
							position++
						l662:
							{
								position663, tokenIndex663, depth663 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l663
								}
								position++
								goto l662
							l663:
								position, tokenIndex, depth = position663, tokenIndex663, depth663
							}
							{
								position664, tokenIndex664, depth664 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l664
								}
								position++
							l666:
								{
									position667, tokenIndex667, depth667 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l667
									}
									position++
									goto l666
								l667:
									position, tokenIndex, depth = position667, tokenIndex667, depth667
								}
								goto l665
							l664:
								position, tokenIndex, depth = position664, tokenIndex664, depth664
							}
						l665:
							if !rules[ruleskip]() {
								goto l654
							}
							depth--
							add(rulesignedNumericLiteral, position659)
						}
					}
				l655:
					goto l653
				l654:
					position, tokenIndex, depth = position654, tokenIndex654, depth654
				}
				depth--
				add(rulenumericExpression, position652)
			}
			return true
		l651:
			position, tokenIndex, depth = position651, tokenIndex651, depth651
			return false
		},
		/* 57 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position668, tokenIndex668, depth668 := position, tokenIndex, depth
			{
				position669 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l668
				}
			l670:
				{
					position671, tokenIndex671, depth671 := position, tokenIndex, depth
					{
						position672, tokenIndex672, depth672 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l673
						}
						goto l672
					l673:
						position, tokenIndex, depth = position672, tokenIndex672, depth672
						if !rules[ruleSLASH]() {
							goto l671
						}
					}
				l672:
					if !rules[ruleunaryExpression]() {
						goto l671
					}
					goto l670
				l671:
					position, tokenIndex, depth = position671, tokenIndex671, depth671
				}
				depth--
				add(rulemultiplicativeExpression, position669)
			}
			return true
		l668:
			position, tokenIndex, depth = position668, tokenIndex668, depth668
			return false
		},
		/* 58 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position674, tokenIndex674, depth674 := position, tokenIndex, depth
			{
				position675 := position
				depth++
				{
					position676, tokenIndex676, depth676 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l676
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l676
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l676
							}
							break
						}
					}

					goto l677
				l676:
					position, tokenIndex, depth = position676, tokenIndex676, depth676
				}
			l677:
				{
					position679 := position
					depth++
					{
						position680, tokenIndex680, depth680 := position, tokenIndex, depth
						if !rules[rulebuiltinCall]() {
							goto l681
						}
						goto l680
					l681:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if !rules[rulefunctionCall]() {
							goto l682
						}
						goto l680
					l682:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						if !rules[ruleiriref]() {
							goto l683
						}
						goto l680
					l683:
						position, tokenIndex, depth = position680, tokenIndex680, depth680
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position685 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position687 := position
												depth++
												{
													position688 := position
													depth++
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
															goto l674
														}
														position++
													}
												l689:
													{
														position691, tokenIndex691, depth691 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l692
														}
														position++
														goto l691
													l692:
														position, tokenIndex, depth = position691, tokenIndex691, depth691
														if buffer[position] != rune('R') {
															goto l674
														}
														position++
													}
												l691:
													{
														position693, tokenIndex693, depth693 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l694
														}
														position++
														goto l693
													l694:
														position, tokenIndex, depth = position693, tokenIndex693, depth693
														if buffer[position] != rune('O') {
															goto l674
														}
														position++
													}
												l693:
													{
														position695, tokenIndex695, depth695 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l696
														}
														position++
														goto l695
													l696:
														position, tokenIndex, depth = position695, tokenIndex695, depth695
														if buffer[position] != rune('U') {
															goto l674
														}
														position++
													}
												l695:
													{
														position697, tokenIndex697, depth697 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l698
														}
														position++
														goto l697
													l698:
														position, tokenIndex, depth = position697, tokenIndex697, depth697
														if buffer[position] != rune('P') {
															goto l674
														}
														position++
													}
												l697:
													if buffer[position] != rune('_') {
														goto l674
													}
													position++
													{
														position699, tokenIndex699, depth699 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l700
														}
														position++
														goto l699
													l700:
														position, tokenIndex, depth = position699, tokenIndex699, depth699
														if buffer[position] != rune('C') {
															goto l674
														}
														position++
													}
												l699:
													{
														position701, tokenIndex701, depth701 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l702
														}
														position++
														goto l701
													l702:
														position, tokenIndex, depth = position701, tokenIndex701, depth701
														if buffer[position] != rune('O') {
															goto l674
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
															goto l674
														}
														position++
													}
												l703:
													{
														position705, tokenIndex705, depth705 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l706
														}
														position++
														goto l705
													l706:
														position, tokenIndex, depth = position705, tokenIndex705, depth705
														if buffer[position] != rune('C') {
															goto l674
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
															goto l674
														}
														position++
													}
												l707:
													{
														position709, tokenIndex709, depth709 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l710
														}
														position++
														goto l709
													l710:
														position, tokenIndex, depth = position709, tokenIndex709, depth709
														if buffer[position] != rune('T') {
															goto l674
														}
														position++
													}
												l709:
													if !rules[ruleskip]() {
														goto l674
													}
													depth--
													add(ruleGROUPCONCAT, position688)
												}
												if !rules[ruleLPAREN]() {
													goto l674
												}
												{
													position711, tokenIndex711, depth711 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l711
													}
													goto l712
												l711:
													position, tokenIndex, depth = position711, tokenIndex711, depth711
												}
											l712:
												if !rules[ruleexpression]() {
													goto l674
												}
												{
													position713, tokenIndex713, depth713 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l713
													}
													{
														position715 := position
														depth++
														{
															position716, tokenIndex716, depth716 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l717
															}
															position++
															goto l716
														l717:
															position, tokenIndex, depth = position716, tokenIndex716, depth716
															if buffer[position] != rune('S') {
																goto l713
															}
															position++
														}
													l716:
														{
															position718, tokenIndex718, depth718 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l719
															}
															position++
															goto l718
														l719:
															position, tokenIndex, depth = position718, tokenIndex718, depth718
															if buffer[position] != rune('E') {
																goto l713
															}
															position++
														}
													l718:
														{
															position720, tokenIndex720, depth720 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l721
															}
															position++
															goto l720
														l721:
															position, tokenIndex, depth = position720, tokenIndex720, depth720
															if buffer[position] != rune('P') {
																goto l713
															}
															position++
														}
													l720:
														{
															position722, tokenIndex722, depth722 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l723
															}
															position++
															goto l722
														l723:
															position, tokenIndex, depth = position722, tokenIndex722, depth722
															if buffer[position] != rune('A') {
																goto l713
															}
															position++
														}
													l722:
														{
															position724, tokenIndex724, depth724 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l725
															}
															position++
															goto l724
														l725:
															position, tokenIndex, depth = position724, tokenIndex724, depth724
															if buffer[position] != rune('R') {
																goto l713
															}
															position++
														}
													l724:
														{
															position726, tokenIndex726, depth726 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l727
															}
															position++
															goto l726
														l727:
															position, tokenIndex, depth = position726, tokenIndex726, depth726
															if buffer[position] != rune('A') {
																goto l713
															}
															position++
														}
													l726:
														{
															position728, tokenIndex728, depth728 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l729
															}
															position++
															goto l728
														l729:
															position, tokenIndex, depth = position728, tokenIndex728, depth728
															if buffer[position] != rune('T') {
																goto l713
															}
															position++
														}
													l728:
														{
															position730, tokenIndex730, depth730 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l731
															}
															position++
															goto l730
														l731:
															position, tokenIndex, depth = position730, tokenIndex730, depth730
															if buffer[position] != rune('O') {
																goto l713
															}
															position++
														}
													l730:
														{
															position732, tokenIndex732, depth732 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l733
															}
															position++
															goto l732
														l733:
															position, tokenIndex, depth = position732, tokenIndex732, depth732
															if buffer[position] != rune('R') {
																goto l713
															}
															position++
														}
													l732:
														if !rules[ruleskip]() {
															goto l713
														}
														depth--
														add(ruleSEPARATOR, position715)
													}
													if !rules[ruleEQ]() {
														goto l713
													}
													if !rules[rulestring]() {
														goto l713
													}
													goto l714
												l713:
													position, tokenIndex, depth = position713, tokenIndex713, depth713
												}
											l714:
												if !rules[ruleRPAREN]() {
													goto l674
												}
												depth--
												add(rulegroupConcat, position687)
											}
											break
										case 'C', 'c':
											{
												position734 := position
												depth++
												{
													position735 := position
													depth++
													{
														position736, tokenIndex736, depth736 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l737
														}
														position++
														goto l736
													l737:
														position, tokenIndex, depth = position736, tokenIndex736, depth736
														if buffer[position] != rune('C') {
															goto l674
														}
														position++
													}
												l736:
													{
														position738, tokenIndex738, depth738 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l739
														}
														position++
														goto l738
													l739:
														position, tokenIndex, depth = position738, tokenIndex738, depth738
														if buffer[position] != rune('O') {
															goto l674
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
															goto l674
														}
														position++
													}
												l740:
													{
														position742, tokenIndex742, depth742 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l743
														}
														position++
														goto l742
													l743:
														position, tokenIndex, depth = position742, tokenIndex742, depth742
														if buffer[position] != rune('N') {
															goto l674
														}
														position++
													}
												l742:
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
															goto l674
														}
														position++
													}
												l744:
													if !rules[ruleskip]() {
														goto l674
													}
													depth--
													add(ruleCOUNT, position735)
												}
												if !rules[ruleLPAREN]() {
													goto l674
												}
												{
													position746, tokenIndex746, depth746 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l746
													}
													goto l747
												l746:
													position, tokenIndex, depth = position746, tokenIndex746, depth746
												}
											l747:
												{
													position748, tokenIndex748, depth748 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l749
													}
													goto l748
												l749:
													position, tokenIndex, depth = position748, tokenIndex748, depth748
													if !rules[ruleexpression]() {
														goto l674
													}
												}
											l748:
												if !rules[ruleRPAREN]() {
													goto l674
												}
												depth--
												add(rulecount, position734)
											}
											break
										default:
											{
												position750, tokenIndex750, depth750 := position, tokenIndex, depth
												{
													position752 := position
													depth++
													{
														position753, tokenIndex753, depth753 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l754
														}
														position++
														goto l753
													l754:
														position, tokenIndex, depth = position753, tokenIndex753, depth753
														if buffer[position] != rune('S') {
															goto l751
														}
														position++
													}
												l753:
													{
														position755, tokenIndex755, depth755 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l756
														}
														position++
														goto l755
													l756:
														position, tokenIndex, depth = position755, tokenIndex755, depth755
														if buffer[position] != rune('U') {
															goto l751
														}
														position++
													}
												l755:
													{
														position757, tokenIndex757, depth757 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l758
														}
														position++
														goto l757
													l758:
														position, tokenIndex, depth = position757, tokenIndex757, depth757
														if buffer[position] != rune('M') {
															goto l751
														}
														position++
													}
												l757:
													if !rules[ruleskip]() {
														goto l751
													}
													depth--
													add(ruleSUM, position752)
												}
												goto l750
											l751:
												position, tokenIndex, depth = position750, tokenIndex750, depth750
												{
													position760 := position
													depth++
													{
														position761, tokenIndex761, depth761 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l762
														}
														position++
														goto l761
													l762:
														position, tokenIndex, depth = position761, tokenIndex761, depth761
														if buffer[position] != rune('M') {
															goto l759
														}
														position++
													}
												l761:
													{
														position763, tokenIndex763, depth763 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l764
														}
														position++
														goto l763
													l764:
														position, tokenIndex, depth = position763, tokenIndex763, depth763
														if buffer[position] != rune('I') {
															goto l759
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
															goto l759
														}
														position++
													}
												l765:
													if !rules[ruleskip]() {
														goto l759
													}
													depth--
													add(ruleMIN, position760)
												}
												goto l750
											l759:
												position, tokenIndex, depth = position750, tokenIndex750, depth750
												{
													switch buffer[position] {
													case 'S', 's':
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
																	goto l674
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
																	goto l674
																}
																position++
															}
														l771:
															{
																position773, tokenIndex773, depth773 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l774
																}
																position++
																goto l773
															l774:
																position, tokenIndex, depth = position773, tokenIndex773, depth773
																if buffer[position] != rune('M') {
																	goto l674
																}
																position++
															}
														l773:
															{
																position775, tokenIndex775, depth775 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l776
																}
																position++
																goto l775
															l776:
																position, tokenIndex, depth = position775, tokenIndex775, depth775
																if buffer[position] != rune('P') {
																	goto l674
																}
																position++
															}
														l775:
															{
																position777, tokenIndex777, depth777 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l778
																}
																position++
																goto l777
															l778:
																position, tokenIndex, depth = position777, tokenIndex777, depth777
																if buffer[position] != rune('L') {
																	goto l674
																}
																position++
															}
														l777:
															{
																position779, tokenIndex779, depth779 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l780
																}
																position++
																goto l779
															l780:
																position, tokenIndex, depth = position779, tokenIndex779, depth779
																if buffer[position] != rune('E') {
																	goto l674
																}
																position++
															}
														l779:
															if !rules[ruleskip]() {
																goto l674
															}
															depth--
															add(ruleSAMPLE, position768)
														}
														break
													case 'A', 'a':
														{
															position781 := position
															depth++
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
																	goto l674
																}
																position++
															}
														l782:
															{
																position784, tokenIndex784, depth784 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l785
																}
																position++
																goto l784
															l785:
																position, tokenIndex, depth = position784, tokenIndex784, depth784
																if buffer[position] != rune('V') {
																	goto l674
																}
																position++
															}
														l784:
															{
																position786, tokenIndex786, depth786 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l787
																}
																position++
																goto l786
															l787:
																position, tokenIndex, depth = position786, tokenIndex786, depth786
																if buffer[position] != rune('G') {
																	goto l674
																}
																position++
															}
														l786:
															if !rules[ruleskip]() {
																goto l674
															}
															depth--
															add(ruleAVG, position781)
														}
														break
													default:
														{
															position788 := position
															depth++
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
																	goto l674
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
																	goto l674
																}
																position++
															}
														l791:
															{
																position793, tokenIndex793, depth793 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l794
																}
																position++
																goto l793
															l794:
																position, tokenIndex, depth = position793, tokenIndex793, depth793
																if buffer[position] != rune('X') {
																	goto l674
																}
																position++
															}
														l793:
															if !rules[ruleskip]() {
																goto l674
															}
															depth--
															add(ruleMAX, position788)
														}
														break
													}
												}

											}
										l750:
											if !rules[ruleLPAREN]() {
												goto l674
											}
											{
												position795, tokenIndex795, depth795 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l795
												}
												goto l796
											l795:
												position, tokenIndex, depth = position795, tokenIndex795, depth795
											}
										l796:
											if !rules[ruleexpression]() {
												goto l674
											}
											if !rules[ruleRPAREN]() {
												goto l674
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position685)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l674
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l674
								}
								break
							case '(':
								if !rules[rulebrackettedExpression]() {
									goto l674
								}
								break
							case '"', '\'':
								if !rules[ruleliteral]() {
									goto l674
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l674
								}
								break
							}
						}

					}
				l680:
					depth--
					add(ruleprimaryExpression, position679)
				}
				depth--
				add(ruleunaryExpression, position675)
			}
			return true
		l674:
			position, tokenIndex, depth = position674, tokenIndex674, depth674
			return false
		},
		/* 59 primaryExpression <- <(builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('(') brackettedExpression) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 60 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position798, tokenIndex798, depth798 := position, tokenIndex, depth
			{
				position799 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l798
				}
				if !rules[ruleexpression]() {
					goto l798
				}
				if !rules[ruleRPAREN]() {
					goto l798
				}
				depth--
				add(rulebrackettedExpression, position799)
			}
			return true
		l798:
			position, tokenIndex, depth = position798, tokenIndex798, depth798
			return false
		},
		/* 61 functionCall <- <(iriref argList)> */
		func() bool {
			position800, tokenIndex800, depth800 := position, tokenIndex, depth
			{
				position801 := position
				depth++
				if !rules[ruleiriref]() {
					goto l800
				}
				if !rules[ruleargList]() {
					goto l800
				}
				depth--
				add(rulefunctionCall, position801)
			}
			return true
		l800:
			position, tokenIndex, depth = position800, tokenIndex800, depth800
			return false
		},
		/* 62 in <- <(IN argList)> */
		nil,
		/* 63 notin <- <(NOTIN argList)> */
		nil,
		/* 64 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position804, tokenIndex804, depth804 := position, tokenIndex, depth
			{
				position805 := position
				depth++
				{
					position806, tokenIndex806, depth806 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l807
					}
					goto l806
				l807:
					position, tokenIndex, depth = position806, tokenIndex806, depth806
					if !rules[ruleLPAREN]() {
						goto l804
					}
					if !rules[ruleexpression]() {
						goto l804
					}
				l808:
					{
						position809, tokenIndex809, depth809 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l809
						}
						if !rules[ruleexpression]() {
							goto l809
						}
						goto l808
					l809:
						position, tokenIndex, depth = position809, tokenIndex809, depth809
					}
					if !rules[ruleRPAREN]() {
						goto l804
					}
				}
			l806:
				depth--
				add(ruleargList, position805)
			}
			return true
		l804:
			position, tokenIndex, depth = position804, tokenIndex804, depth804
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
			position813, tokenIndex813, depth813 := position, tokenIndex, depth
			{
				position814 := position
				depth++
				{
					position815, tokenIndex815, depth815 := position, tokenIndex, depth
					{
						position817, tokenIndex817, depth817 := position, tokenIndex, depth
						{
							position819 := position
							depth++
							{
								position820, tokenIndex820, depth820 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l821
								}
								position++
								goto l820
							l821:
								position, tokenIndex, depth = position820, tokenIndex820, depth820
								if buffer[position] != rune('S') {
									goto l818
								}
								position++
							}
						l820:
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
									goto l818
								}
								position++
							}
						l822:
							{
								position824, tokenIndex824, depth824 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l825
								}
								position++
								goto l824
							l825:
								position, tokenIndex, depth = position824, tokenIndex824, depth824
								if buffer[position] != rune('R') {
									goto l818
								}
								position++
							}
						l824:
							if !rules[ruleskip]() {
								goto l818
							}
							depth--
							add(ruleSTR, position819)
						}
						goto l817
					l818:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position827 := position
							depth++
							{
								position828, tokenIndex828, depth828 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l829
								}
								position++
								goto l828
							l829:
								position, tokenIndex, depth = position828, tokenIndex828, depth828
								if buffer[position] != rune('L') {
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
								if buffer[position] != rune('n') {
									goto l833
								}
								position++
								goto l832
							l833:
								position, tokenIndex, depth = position832, tokenIndex832, depth832
								if buffer[position] != rune('N') {
									goto l826
								}
								position++
							}
						l832:
							{
								position834, tokenIndex834, depth834 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l835
								}
								position++
								goto l834
							l835:
								position, tokenIndex, depth = position834, tokenIndex834, depth834
								if buffer[position] != rune('G') {
									goto l826
								}
								position++
							}
						l834:
							if !rules[ruleskip]() {
								goto l826
							}
							depth--
							add(ruleLANG, position827)
						}
						goto l817
					l826:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
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
									goto l836
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
									goto l836
								}
								position++
							}
						l840:
							{
								position842, tokenIndex842, depth842 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l843
								}
								position++
								goto l842
							l843:
								position, tokenIndex, depth = position842, tokenIndex842, depth842
								if buffer[position] != rune('T') {
									goto l836
								}
								position++
							}
						l842:
							{
								position844, tokenIndex844, depth844 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l845
								}
								position++
								goto l844
							l845:
								position, tokenIndex, depth = position844, tokenIndex844, depth844
								if buffer[position] != rune('A') {
									goto l836
								}
								position++
							}
						l844:
							{
								position846, tokenIndex846, depth846 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l847
								}
								position++
								goto l846
							l847:
								position, tokenIndex, depth = position846, tokenIndex846, depth846
								if buffer[position] != rune('T') {
									goto l836
								}
								position++
							}
						l846:
							{
								position848, tokenIndex848, depth848 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l849
								}
								position++
								goto l848
							l849:
								position, tokenIndex, depth = position848, tokenIndex848, depth848
								if buffer[position] != rune('Y') {
									goto l836
								}
								position++
							}
						l848:
							{
								position850, tokenIndex850, depth850 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l851
								}
								position++
								goto l850
							l851:
								position, tokenIndex, depth = position850, tokenIndex850, depth850
								if buffer[position] != rune('P') {
									goto l836
								}
								position++
							}
						l850:
							{
								position852, tokenIndex852, depth852 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l853
								}
								position++
								goto l852
							l853:
								position, tokenIndex, depth = position852, tokenIndex852, depth852
								if buffer[position] != rune('E') {
									goto l836
								}
								position++
							}
						l852:
							if !rules[ruleskip]() {
								goto l836
							}
							depth--
							add(ruleDATATYPE, position837)
						}
						goto l817
					l836:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position855 := position
							depth++
							{
								position856, tokenIndex856, depth856 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l857
								}
								position++
								goto l856
							l857:
								position, tokenIndex, depth = position856, tokenIndex856, depth856
								if buffer[position] != rune('I') {
									goto l854
								}
								position++
							}
						l856:
							{
								position858, tokenIndex858, depth858 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('R') {
									goto l854
								}
								position++
							}
						l858:
							{
								position860, tokenIndex860, depth860 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l861
								}
								position++
								goto l860
							l861:
								position, tokenIndex, depth = position860, tokenIndex860, depth860
								if buffer[position] != rune('I') {
									goto l854
								}
								position++
							}
						l860:
							if !rules[ruleskip]() {
								goto l854
							}
							depth--
							add(ruleIRI, position855)
						}
						goto l817
					l854:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position863 := position
							depth++
							{
								position864, tokenIndex864, depth864 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l865
								}
								position++
								goto l864
							l865:
								position, tokenIndex, depth = position864, tokenIndex864, depth864
								if buffer[position] != rune('U') {
									goto l862
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
									goto l862
								}
								position++
							}
						l866:
							{
								position868, tokenIndex868, depth868 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l869
								}
								position++
								goto l868
							l869:
								position, tokenIndex, depth = position868, tokenIndex868, depth868
								if buffer[position] != rune('I') {
									goto l862
								}
								position++
							}
						l868:
							if !rules[ruleskip]() {
								goto l862
							}
							depth--
							add(ruleURI, position863)
						}
						goto l817
					l862:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position871 := position
							depth++
							{
								position872, tokenIndex872, depth872 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l873
								}
								position++
								goto l872
							l873:
								position, tokenIndex, depth = position872, tokenIndex872, depth872
								if buffer[position] != rune('S') {
									goto l870
								}
								position++
							}
						l872:
							{
								position874, tokenIndex874, depth874 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l875
								}
								position++
								goto l874
							l875:
								position, tokenIndex, depth = position874, tokenIndex874, depth874
								if buffer[position] != rune('T') {
									goto l870
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
									goto l870
								}
								position++
							}
						l876:
							{
								position878, tokenIndex878, depth878 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l879
								}
								position++
								goto l878
							l879:
								position, tokenIndex, depth = position878, tokenIndex878, depth878
								if buffer[position] != rune('L') {
									goto l870
								}
								position++
							}
						l878:
							{
								position880, tokenIndex880, depth880 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l881
								}
								position++
								goto l880
							l881:
								position, tokenIndex, depth = position880, tokenIndex880, depth880
								if buffer[position] != rune('E') {
									goto l870
								}
								position++
							}
						l880:
							{
								position882, tokenIndex882, depth882 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l883
								}
								position++
								goto l882
							l883:
								position, tokenIndex, depth = position882, tokenIndex882, depth882
								if buffer[position] != rune('N') {
									goto l870
								}
								position++
							}
						l882:
							if !rules[ruleskip]() {
								goto l870
							}
							depth--
							add(ruleSTRLEN, position871)
						}
						goto l817
					l870:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position885 := position
							depth++
							{
								position886, tokenIndex886, depth886 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l887
								}
								position++
								goto l886
							l887:
								position, tokenIndex, depth = position886, tokenIndex886, depth886
								if buffer[position] != rune('M') {
									goto l884
								}
								position++
							}
						l886:
							{
								position888, tokenIndex888, depth888 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l889
								}
								position++
								goto l888
							l889:
								position, tokenIndex, depth = position888, tokenIndex888, depth888
								if buffer[position] != rune('O') {
									goto l884
								}
								position++
							}
						l888:
							{
								position890, tokenIndex890, depth890 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l891
								}
								position++
								goto l890
							l891:
								position, tokenIndex, depth = position890, tokenIndex890, depth890
								if buffer[position] != rune('N') {
									goto l884
								}
								position++
							}
						l890:
							{
								position892, tokenIndex892, depth892 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l893
								}
								position++
								goto l892
							l893:
								position, tokenIndex, depth = position892, tokenIndex892, depth892
								if buffer[position] != rune('T') {
									goto l884
								}
								position++
							}
						l892:
							{
								position894, tokenIndex894, depth894 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l895
								}
								position++
								goto l894
							l895:
								position, tokenIndex, depth = position894, tokenIndex894, depth894
								if buffer[position] != rune('H') {
									goto l884
								}
								position++
							}
						l894:
							if !rules[ruleskip]() {
								goto l884
							}
							depth--
							add(ruleMONTH, position885)
						}
						goto l817
					l884:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position897 := position
							depth++
							{
								position898, tokenIndex898, depth898 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l899
								}
								position++
								goto l898
							l899:
								position, tokenIndex, depth = position898, tokenIndex898, depth898
								if buffer[position] != rune('M') {
									goto l896
								}
								position++
							}
						l898:
							{
								position900, tokenIndex900, depth900 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l901
								}
								position++
								goto l900
							l901:
								position, tokenIndex, depth = position900, tokenIndex900, depth900
								if buffer[position] != rune('I') {
									goto l896
								}
								position++
							}
						l900:
							{
								position902, tokenIndex902, depth902 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l903
								}
								position++
								goto l902
							l903:
								position, tokenIndex, depth = position902, tokenIndex902, depth902
								if buffer[position] != rune('N') {
									goto l896
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
									goto l896
								}
								position++
							}
						l904:
							{
								position906, tokenIndex906, depth906 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l907
								}
								position++
								goto l906
							l907:
								position, tokenIndex, depth = position906, tokenIndex906, depth906
								if buffer[position] != rune('T') {
									goto l896
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
									goto l896
								}
								position++
							}
						l908:
							{
								position910, tokenIndex910, depth910 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l911
								}
								position++
								goto l910
							l911:
								position, tokenIndex, depth = position910, tokenIndex910, depth910
								if buffer[position] != rune('S') {
									goto l896
								}
								position++
							}
						l910:
							if !rules[ruleskip]() {
								goto l896
							}
							depth--
							add(ruleMINUTES, position897)
						}
						goto l817
					l896:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position913 := position
							depth++
							{
								position914, tokenIndex914, depth914 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l915
								}
								position++
								goto l914
							l915:
								position, tokenIndex, depth = position914, tokenIndex914, depth914
								if buffer[position] != rune('S') {
									goto l912
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
									goto l912
								}
								position++
							}
						l916:
							{
								position918, tokenIndex918, depth918 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l919
								}
								position++
								goto l918
							l919:
								position, tokenIndex, depth = position918, tokenIndex918, depth918
								if buffer[position] != rune('C') {
									goto l912
								}
								position++
							}
						l918:
							{
								position920, tokenIndex920, depth920 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l921
								}
								position++
								goto l920
							l921:
								position, tokenIndex, depth = position920, tokenIndex920, depth920
								if buffer[position] != rune('O') {
									goto l912
								}
								position++
							}
						l920:
							{
								position922, tokenIndex922, depth922 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l923
								}
								position++
								goto l922
							l923:
								position, tokenIndex, depth = position922, tokenIndex922, depth922
								if buffer[position] != rune('N') {
									goto l912
								}
								position++
							}
						l922:
							{
								position924, tokenIndex924, depth924 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l925
								}
								position++
								goto l924
							l925:
								position, tokenIndex, depth = position924, tokenIndex924, depth924
								if buffer[position] != rune('D') {
									goto l912
								}
								position++
							}
						l924:
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
									goto l912
								}
								position++
							}
						l926:
							if !rules[ruleskip]() {
								goto l912
							}
							depth--
							add(ruleSECONDS, position913)
						}
						goto l817
					l912:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position929 := position
							depth++
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
									goto l928
								}
								position++
							}
						l930:
							{
								position932, tokenIndex932, depth932 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l933
								}
								position++
								goto l932
							l933:
								position, tokenIndex, depth = position932, tokenIndex932, depth932
								if buffer[position] != rune('I') {
									goto l928
								}
								position++
							}
						l932:
							{
								position934, tokenIndex934, depth934 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l935
								}
								position++
								goto l934
							l935:
								position, tokenIndex, depth = position934, tokenIndex934, depth934
								if buffer[position] != rune('M') {
									goto l928
								}
								position++
							}
						l934:
							{
								position936, tokenIndex936, depth936 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l937
								}
								position++
								goto l936
							l937:
								position, tokenIndex, depth = position936, tokenIndex936, depth936
								if buffer[position] != rune('E') {
									goto l928
								}
								position++
							}
						l936:
							{
								position938, tokenIndex938, depth938 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l939
								}
								position++
								goto l938
							l939:
								position, tokenIndex, depth = position938, tokenIndex938, depth938
								if buffer[position] != rune('Z') {
									goto l928
								}
								position++
							}
						l938:
							{
								position940, tokenIndex940, depth940 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l941
								}
								position++
								goto l940
							l941:
								position, tokenIndex, depth = position940, tokenIndex940, depth940
								if buffer[position] != rune('O') {
									goto l928
								}
								position++
							}
						l940:
							{
								position942, tokenIndex942, depth942 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l943
								}
								position++
								goto l942
							l943:
								position, tokenIndex, depth = position942, tokenIndex942, depth942
								if buffer[position] != rune('N') {
									goto l928
								}
								position++
							}
						l942:
							{
								position944, tokenIndex944, depth944 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l945
								}
								position++
								goto l944
							l945:
								position, tokenIndex, depth = position944, tokenIndex944, depth944
								if buffer[position] != rune('E') {
									goto l928
								}
								position++
							}
						l944:
							if !rules[ruleskip]() {
								goto l928
							}
							depth--
							add(ruleTIMEZONE, position929)
						}
						goto l817
					l928:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position947 := position
							depth++
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
									goto l946
								}
								position++
							}
						l948:
							{
								position950, tokenIndex950, depth950 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l951
								}
								position++
								goto l950
							l951:
								position, tokenIndex, depth = position950, tokenIndex950, depth950
								if buffer[position] != rune('H') {
									goto l946
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
									goto l946
								}
								position++
							}
						l952:
							if buffer[position] != rune('1') {
								goto l946
							}
							position++
							if !rules[ruleskip]() {
								goto l946
							}
							depth--
							add(ruleSHA1, position947)
						}
						goto l817
					l946:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position955 := position
							depth++
							{
								position956, tokenIndex956, depth956 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l957
								}
								position++
								goto l956
							l957:
								position, tokenIndex, depth = position956, tokenIndex956, depth956
								if buffer[position] != rune('S') {
									goto l954
								}
								position++
							}
						l956:
							{
								position958, tokenIndex958, depth958 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l959
								}
								position++
								goto l958
							l959:
								position, tokenIndex, depth = position958, tokenIndex958, depth958
								if buffer[position] != rune('H') {
									goto l954
								}
								position++
							}
						l958:
							{
								position960, tokenIndex960, depth960 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l961
								}
								position++
								goto l960
							l961:
								position, tokenIndex, depth = position960, tokenIndex960, depth960
								if buffer[position] != rune('A') {
									goto l954
								}
								position++
							}
						l960:
							if buffer[position] != rune('2') {
								goto l954
							}
							position++
							if buffer[position] != rune('5') {
								goto l954
							}
							position++
							if buffer[position] != rune('6') {
								goto l954
							}
							position++
							if !rules[ruleskip]() {
								goto l954
							}
							depth--
							add(ruleSHA256, position955)
						}
						goto l817
					l954:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position963 := position
							depth++
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
									goto l962
								}
								position++
							}
						l964:
							{
								position966, tokenIndex966, depth966 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l967
								}
								position++
								goto l966
							l967:
								position, tokenIndex, depth = position966, tokenIndex966, depth966
								if buffer[position] != rune('H') {
									goto l962
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
									goto l962
								}
								position++
							}
						l968:
							if buffer[position] != rune('3') {
								goto l962
							}
							position++
							if buffer[position] != rune('8') {
								goto l962
							}
							position++
							if buffer[position] != rune('4') {
								goto l962
							}
							position++
							if !rules[ruleskip]() {
								goto l962
							}
							depth--
							add(ruleSHA384, position963)
						}
						goto l817
					l962:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position971 := position
							depth++
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
									goto l970
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
									goto l970
								}
								position++
							}
						l974:
							{
								position976, tokenIndex976, depth976 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l977
								}
								position++
								goto l976
							l977:
								position, tokenIndex, depth = position976, tokenIndex976, depth976
								if buffer[position] != rune('I') {
									goto l970
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
									goto l970
								}
								position++
							}
						l978:
							{
								position980, tokenIndex980, depth980 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l981
								}
								position++
								goto l980
							l981:
								position, tokenIndex, depth = position980, tokenIndex980, depth980
								if buffer[position] != rune('I') {
									goto l970
								}
								position++
							}
						l980:
							if !rules[ruleskip]() {
								goto l970
							}
							depth--
							add(ruleISIRI, position971)
						}
						goto l817
					l970:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position983 := position
							depth++
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
									goto l982
								}
								position++
							}
						l984:
							{
								position986, tokenIndex986, depth986 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l987
								}
								position++
								goto l986
							l987:
								position, tokenIndex, depth = position986, tokenIndex986, depth986
								if buffer[position] != rune('S') {
									goto l982
								}
								position++
							}
						l986:
							{
								position988, tokenIndex988, depth988 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l989
								}
								position++
								goto l988
							l989:
								position, tokenIndex, depth = position988, tokenIndex988, depth988
								if buffer[position] != rune('U') {
									goto l982
								}
								position++
							}
						l988:
							{
								position990, tokenIndex990, depth990 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l991
								}
								position++
								goto l990
							l991:
								position, tokenIndex, depth = position990, tokenIndex990, depth990
								if buffer[position] != rune('R') {
									goto l982
								}
								position++
							}
						l990:
							{
								position992, tokenIndex992, depth992 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l993
								}
								position++
								goto l992
							l993:
								position, tokenIndex, depth = position992, tokenIndex992, depth992
								if buffer[position] != rune('I') {
									goto l982
								}
								position++
							}
						l992:
							if !rules[ruleskip]() {
								goto l982
							}
							depth--
							add(ruleISURI, position983)
						}
						goto l817
					l982:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position995 := position
							depth++
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
									goto l994
								}
								position++
							}
						l996:
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
									goto l994
								}
								position++
							}
						l998:
							{
								position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1001
								}
								position++
								goto l1000
							l1001:
								position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
								if buffer[position] != rune('B') {
									goto l994
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
									goto l994
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
									goto l994
								}
								position++
							}
						l1004:
							{
								position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1007
								}
								position++
								goto l1006
							l1007:
								position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
								if buffer[position] != rune('N') {
									goto l994
								}
								position++
							}
						l1006:
							{
								position1008, tokenIndex1008, depth1008 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l1009
								}
								position++
								goto l1008
							l1009:
								position, tokenIndex, depth = position1008, tokenIndex1008, depth1008
								if buffer[position] != rune('K') {
									goto l994
								}
								position++
							}
						l1008:
							if !rules[ruleskip]() {
								goto l994
							}
							depth--
							add(ruleISBLANK, position995)
						}
						goto l817
					l994:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							position1011 := position
							depth++
							{
								position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1013
								}
								position++
								goto l1012
							l1013:
								position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
								if buffer[position] != rune('I') {
									goto l1010
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
									goto l1010
								}
								position++
							}
						l1014:
							{
								position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1017
								}
								position++
								goto l1016
							l1017:
								position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
								if buffer[position] != rune('L') {
									goto l1010
								}
								position++
							}
						l1016:
							{
								position1018, tokenIndex1018, depth1018 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1019
								}
								position++
								goto l1018
							l1019:
								position, tokenIndex, depth = position1018, tokenIndex1018, depth1018
								if buffer[position] != rune('I') {
									goto l1010
								}
								position++
							}
						l1018:
							{
								position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1021
								}
								position++
								goto l1020
							l1021:
								position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
								if buffer[position] != rune('T') {
									goto l1010
								}
								position++
							}
						l1020:
							{
								position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1023
								}
								position++
								goto l1022
							l1023:
								position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
								if buffer[position] != rune('E') {
									goto l1010
								}
								position++
							}
						l1022:
							{
								position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1025
								}
								position++
								goto l1024
							l1025:
								position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
								if buffer[position] != rune('R') {
									goto l1010
								}
								position++
							}
						l1024:
							{
								position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1027
								}
								position++
								goto l1026
							l1027:
								position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
								if buffer[position] != rune('A') {
									goto l1010
								}
								position++
							}
						l1026:
							{
								position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1029
								}
								position++
								goto l1028
							l1029:
								position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
								if buffer[position] != rune('L') {
									goto l1010
								}
								position++
							}
						l1028:
							if !rules[ruleskip]() {
								goto l1010
							}
							depth--
							add(ruleISLITERAL, position1011)
						}
						goto l817
					l1010:
						position, tokenIndex, depth = position817, tokenIndex817, depth817
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position1031 := position
									depth++
									{
										position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1033
										}
										position++
										goto l1032
									l1033:
										position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
										if buffer[position] != rune('I') {
											goto l816
										}
										position++
									}
								l1032:
									{
										position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1035
										}
										position++
										goto l1034
									l1035:
										position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
										if buffer[position] != rune('S') {
											goto l816
										}
										position++
									}
								l1034:
									{
										position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1037
										}
										position++
										goto l1036
									l1037:
										position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
										if buffer[position] != rune('N') {
											goto l816
										}
										position++
									}
								l1036:
									{
										position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1039
										}
										position++
										goto l1038
									l1039:
										position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
										if buffer[position] != rune('U') {
											goto l816
										}
										position++
									}
								l1038:
									{
										position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1041
										}
										position++
										goto l1040
									l1041:
										position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
										if buffer[position] != rune('M') {
											goto l816
										}
										position++
									}
								l1040:
									{
										position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1043
										}
										position++
										goto l1042
									l1043:
										position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
										if buffer[position] != rune('E') {
											goto l816
										}
										position++
									}
								l1042:
									{
										position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1045
										}
										position++
										goto l1044
									l1045:
										position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
										if buffer[position] != rune('R') {
											goto l816
										}
										position++
									}
								l1044:
									{
										position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1047
										}
										position++
										goto l1046
									l1047:
										position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
										if buffer[position] != rune('I') {
											goto l816
										}
										position++
									}
								l1046:
									{
										position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1049
										}
										position++
										goto l1048
									l1049:
										position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
										if buffer[position] != rune('C') {
											goto l816
										}
										position++
									}
								l1048:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleISNUMERIC, position1031)
								}
								break
							case 'S', 's':
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
											goto l816
										}
										position++
									}
								l1051:
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
											goto l816
										}
										position++
									}
								l1053:
									{
										position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1056
										}
										position++
										goto l1055
									l1056:
										position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
										if buffer[position] != rune('A') {
											goto l816
										}
										position++
									}
								l1055:
									if buffer[position] != rune('5') {
										goto l816
									}
									position++
									if buffer[position] != rune('1') {
										goto l816
									}
									position++
									if buffer[position] != rune('2') {
										goto l816
									}
									position++
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleSHA512, position1050)
								}
								break
							case 'M', 'm':
								{
									position1057 := position
									depth++
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
											goto l816
										}
										position++
									}
								l1058:
									{
										position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1061
										}
										position++
										goto l1060
									l1061:
										position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
										if buffer[position] != rune('D') {
											goto l816
										}
										position++
									}
								l1060:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleMD5, position1057)
								}
								break
							case 'T', 't':
								{
									position1062 := position
									depth++
									{
										position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1064
										}
										position++
										goto l1063
									l1064:
										position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
										if buffer[position] != rune('T') {
											goto l816
										}
										position++
									}
								l1063:
									{
										position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1066
										}
										position++
										goto l1065
									l1066:
										position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
										if buffer[position] != rune('Z') {
											goto l816
										}
										position++
									}
								l1065:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleTZ, position1062)
								}
								break
							case 'H', 'h':
								{
									position1067 := position
									depth++
									{
										position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1069
										}
										position++
										goto l1068
									l1069:
										position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
										if buffer[position] != rune('H') {
											goto l816
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
											goto l816
										}
										position++
									}
								l1070:
									{
										position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1073
										}
										position++
										goto l1072
									l1073:
										position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
										if buffer[position] != rune('U') {
											goto l816
										}
										position++
									}
								l1072:
									{
										position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1075
										}
										position++
										goto l1074
									l1075:
										position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
										if buffer[position] != rune('R') {
											goto l816
										}
										position++
									}
								l1074:
									{
										position1076, tokenIndex1076, depth1076 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1077
										}
										position++
										goto l1076
									l1077:
										position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
										if buffer[position] != rune('S') {
											goto l816
										}
										position++
									}
								l1076:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleHOURS, position1067)
								}
								break
							case 'D', 'd':
								{
									position1078 := position
									depth++
									{
										position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1080
										}
										position++
										goto l1079
									l1080:
										position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
										if buffer[position] != rune('D') {
											goto l816
										}
										position++
									}
								l1079:
									{
										position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1082
										}
										position++
										goto l1081
									l1082:
										position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
										if buffer[position] != rune('A') {
											goto l816
										}
										position++
									}
								l1081:
									{
										position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1084
										}
										position++
										goto l1083
									l1084:
										position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
										if buffer[position] != rune('Y') {
											goto l816
										}
										position++
									}
								l1083:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleDAY, position1078)
								}
								break
							case 'Y', 'y':
								{
									position1085 := position
									depth++
									{
										position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1087
										}
										position++
										goto l1086
									l1087:
										position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
										if buffer[position] != rune('Y') {
											goto l816
										}
										position++
									}
								l1086:
									{
										position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1089
										}
										position++
										goto l1088
									l1089:
										position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
										if buffer[position] != rune('E') {
											goto l816
										}
										position++
									}
								l1088:
									{
										position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1091
										}
										position++
										goto l1090
									l1091:
										position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
										if buffer[position] != rune('A') {
											goto l816
										}
										position++
									}
								l1090:
									{
										position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1093
										}
										position++
										goto l1092
									l1093:
										position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
										if buffer[position] != rune('R') {
											goto l816
										}
										position++
									}
								l1092:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleYEAR, position1085)
								}
								break
							case 'E', 'e':
								{
									position1094 := position
									depth++
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
											goto l816
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
											goto l816
										}
										position++
									}
								l1097:
									{
										position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1100
										}
										position++
										goto l1099
									l1100:
										position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
										if buffer[position] != rune('C') {
											goto l816
										}
										position++
									}
								l1099:
									{
										position1101, tokenIndex1101, depth1101 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1102
										}
										position++
										goto l1101
									l1102:
										position, tokenIndex, depth = position1101, tokenIndex1101, depth1101
										if buffer[position] != rune('O') {
											goto l816
										}
										position++
									}
								l1101:
									{
										position1103, tokenIndex1103, depth1103 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1104
										}
										position++
										goto l1103
									l1104:
										position, tokenIndex, depth = position1103, tokenIndex1103, depth1103
										if buffer[position] != rune('D') {
											goto l816
										}
										position++
									}
								l1103:
									{
										position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1106
										}
										position++
										goto l1105
									l1106:
										position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
										if buffer[position] != rune('E') {
											goto l816
										}
										position++
									}
								l1105:
									if buffer[position] != rune('_') {
										goto l816
									}
									position++
									{
										position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1108
										}
										position++
										goto l1107
									l1108:
										position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
										if buffer[position] != rune('F') {
											goto l816
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
											goto l816
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
											goto l816
										}
										position++
									}
								l1111:
									if buffer[position] != rune('_') {
										goto l816
									}
									position++
									{
										position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1114
										}
										position++
										goto l1113
									l1114:
										position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
										if buffer[position] != rune('U') {
											goto l816
										}
										position++
									}
								l1113:
									{
										position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1116
										}
										position++
										goto l1115
									l1116:
										position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
										if buffer[position] != rune('R') {
											goto l816
										}
										position++
									}
								l1115:
									{
										position1117, tokenIndex1117, depth1117 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1118
										}
										position++
										goto l1117
									l1118:
										position, tokenIndex, depth = position1117, tokenIndex1117, depth1117
										if buffer[position] != rune('I') {
											goto l816
										}
										position++
									}
								l1117:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleENCODEFORURI, position1094)
								}
								break
							case 'L', 'l':
								{
									position1119 := position
									depth++
									{
										position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1121
										}
										position++
										goto l1120
									l1121:
										position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
										if buffer[position] != rune('L') {
											goto l816
										}
										position++
									}
								l1120:
									{
										position1122, tokenIndex1122, depth1122 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1123
										}
										position++
										goto l1122
									l1123:
										position, tokenIndex, depth = position1122, tokenIndex1122, depth1122
										if buffer[position] != rune('C') {
											goto l816
										}
										position++
									}
								l1122:
									{
										position1124, tokenIndex1124, depth1124 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1125
										}
										position++
										goto l1124
									l1125:
										position, tokenIndex, depth = position1124, tokenIndex1124, depth1124
										if buffer[position] != rune('A') {
											goto l816
										}
										position++
									}
								l1124:
									{
										position1126, tokenIndex1126, depth1126 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1127
										}
										position++
										goto l1126
									l1127:
										position, tokenIndex, depth = position1126, tokenIndex1126, depth1126
										if buffer[position] != rune('S') {
											goto l816
										}
										position++
									}
								l1126:
									{
										position1128, tokenIndex1128, depth1128 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1129
										}
										position++
										goto l1128
									l1129:
										position, tokenIndex, depth = position1128, tokenIndex1128, depth1128
										if buffer[position] != rune('E') {
											goto l816
										}
										position++
									}
								l1128:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleLCASE, position1119)
								}
								break
							case 'U', 'u':
								{
									position1130 := position
									depth++
									{
										position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1132
										}
										position++
										goto l1131
									l1132:
										position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
										if buffer[position] != rune('U') {
											goto l816
										}
										position++
									}
								l1131:
									{
										position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1134
										}
										position++
										goto l1133
									l1134:
										position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
										if buffer[position] != rune('C') {
											goto l816
										}
										position++
									}
								l1133:
									{
										position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1136
										}
										position++
										goto l1135
									l1136:
										position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
										if buffer[position] != rune('A') {
											goto l816
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
											goto l816
										}
										position++
									}
								l1137:
									{
										position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1140
										}
										position++
										goto l1139
									l1140:
										position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
										if buffer[position] != rune('E') {
											goto l816
										}
										position++
									}
								l1139:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleUCASE, position1130)
								}
								break
							case 'F', 'f':
								{
									position1141 := position
									depth++
									{
										position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1143
										}
										position++
										goto l1142
									l1143:
										position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
										if buffer[position] != rune('F') {
											goto l816
										}
										position++
									}
								l1142:
									{
										position1144, tokenIndex1144, depth1144 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1145
										}
										position++
										goto l1144
									l1145:
										position, tokenIndex, depth = position1144, tokenIndex1144, depth1144
										if buffer[position] != rune('L') {
											goto l816
										}
										position++
									}
								l1144:
									{
										position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1147
										}
										position++
										goto l1146
									l1147:
										position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
										if buffer[position] != rune('O') {
											goto l816
										}
										position++
									}
								l1146:
									{
										position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1149
										}
										position++
										goto l1148
									l1149:
										position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
										if buffer[position] != rune('O') {
											goto l816
										}
										position++
									}
								l1148:
									{
										position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1151
										}
										position++
										goto l1150
									l1151:
										position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
										if buffer[position] != rune('R') {
											goto l816
										}
										position++
									}
								l1150:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleFLOOR, position1141)
								}
								break
							case 'R', 'r':
								{
									position1152 := position
									depth++
									{
										position1153, tokenIndex1153, depth1153 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1154
										}
										position++
										goto l1153
									l1154:
										position, tokenIndex, depth = position1153, tokenIndex1153, depth1153
										if buffer[position] != rune('R') {
											goto l816
										}
										position++
									}
								l1153:
									{
										position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1156
										}
										position++
										goto l1155
									l1156:
										position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
										if buffer[position] != rune('O') {
											goto l816
										}
										position++
									}
								l1155:
									{
										position1157, tokenIndex1157, depth1157 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1158
										}
										position++
										goto l1157
									l1158:
										position, tokenIndex, depth = position1157, tokenIndex1157, depth1157
										if buffer[position] != rune('U') {
											goto l816
										}
										position++
									}
								l1157:
									{
										position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1160
										}
										position++
										goto l1159
									l1160:
										position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
										if buffer[position] != rune('N') {
											goto l816
										}
										position++
									}
								l1159:
									{
										position1161, tokenIndex1161, depth1161 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1162
										}
										position++
										goto l1161
									l1162:
										position, tokenIndex, depth = position1161, tokenIndex1161, depth1161
										if buffer[position] != rune('D') {
											goto l816
										}
										position++
									}
								l1161:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleROUND, position1152)
								}
								break
							case 'C', 'c':
								{
									position1163 := position
									depth++
									{
										position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1165
										}
										position++
										goto l1164
									l1165:
										position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
										if buffer[position] != rune('C') {
											goto l816
										}
										position++
									}
								l1164:
									{
										position1166, tokenIndex1166, depth1166 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1167
										}
										position++
										goto l1166
									l1167:
										position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
										if buffer[position] != rune('E') {
											goto l816
										}
										position++
									}
								l1166:
									{
										position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1169
										}
										position++
										goto l1168
									l1169:
										position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
										if buffer[position] != rune('I') {
											goto l816
										}
										position++
									}
								l1168:
									{
										position1170, tokenIndex1170, depth1170 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1171
										}
										position++
										goto l1170
									l1171:
										position, tokenIndex, depth = position1170, tokenIndex1170, depth1170
										if buffer[position] != rune('L') {
											goto l816
										}
										position++
									}
								l1170:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleCEIL, position1163)
								}
								break
							default:
								{
									position1172 := position
									depth++
									{
										position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1174
										}
										position++
										goto l1173
									l1174:
										position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
										if buffer[position] != rune('A') {
											goto l816
										}
										position++
									}
								l1173:
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
											goto l816
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
											goto l816
										}
										position++
									}
								l1177:
									if !rules[ruleskip]() {
										goto l816
									}
									depth--
									add(ruleABS, position1172)
								}
								break
							}
						}

					}
				l817:
					if !rules[ruleLPAREN]() {
						goto l816
					}
					if !rules[ruleexpression]() {
						goto l816
					}
					if !rules[ruleRPAREN]() {
						goto l816
					}
					goto l815
				l816:
					position, tokenIndex, depth = position815, tokenIndex815, depth815
					{
						position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
						{
							position1182 := position
							depth++
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
									goto l1181
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
									goto l1181
								}
								position++
							}
						l1185:
							{
								position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1188
								}
								position++
								goto l1187
							l1188:
								position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
								if buffer[position] != rune('R') {
									goto l1181
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
									goto l1181
								}
								position++
							}
						l1189:
							{
								position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1192
								}
								position++
								goto l1191
							l1192:
								position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
								if buffer[position] != rune('T') {
									goto l1181
								}
								position++
							}
						l1191:
							{
								position1193, tokenIndex1193, depth1193 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1194
								}
								position++
								goto l1193
							l1194:
								position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
								if buffer[position] != rune('A') {
									goto l1181
								}
								position++
							}
						l1193:
							{
								position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1196
								}
								position++
								goto l1195
							l1196:
								position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
								if buffer[position] != rune('R') {
									goto l1181
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
									goto l1181
								}
								position++
							}
						l1197:
							{
								position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1200
								}
								position++
								goto l1199
							l1200:
								position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
								if buffer[position] != rune('S') {
									goto l1181
								}
								position++
							}
						l1199:
							if !rules[ruleskip]() {
								goto l1181
							}
							depth--
							add(ruleSTRSTARTS, position1182)
						}
						goto l1180
					l1181:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						{
							position1202 := position
							depth++
							{
								position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1204
								}
								position++
								goto l1203
							l1204:
								position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
								if buffer[position] != rune('S') {
									goto l1201
								}
								position++
							}
						l1203:
							{
								position1205, tokenIndex1205, depth1205 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1206
								}
								position++
								goto l1205
							l1206:
								position, tokenIndex, depth = position1205, tokenIndex1205, depth1205
								if buffer[position] != rune('T') {
									goto l1201
								}
								position++
							}
						l1205:
							{
								position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1208
								}
								position++
								goto l1207
							l1208:
								position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
								if buffer[position] != rune('R') {
									goto l1201
								}
								position++
							}
						l1207:
							{
								position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1210
								}
								position++
								goto l1209
							l1210:
								position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
								if buffer[position] != rune('E') {
									goto l1201
								}
								position++
							}
						l1209:
							{
								position1211, tokenIndex1211, depth1211 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1212
								}
								position++
								goto l1211
							l1212:
								position, tokenIndex, depth = position1211, tokenIndex1211, depth1211
								if buffer[position] != rune('N') {
									goto l1201
								}
								position++
							}
						l1211:
							{
								position1213, tokenIndex1213, depth1213 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1214
								}
								position++
								goto l1213
							l1214:
								position, tokenIndex, depth = position1213, tokenIndex1213, depth1213
								if buffer[position] != rune('D') {
									goto l1201
								}
								position++
							}
						l1213:
							{
								position1215, tokenIndex1215, depth1215 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1216
								}
								position++
								goto l1215
							l1216:
								position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
								if buffer[position] != rune('S') {
									goto l1201
								}
								position++
							}
						l1215:
							if !rules[ruleskip]() {
								goto l1201
							}
							depth--
							add(ruleSTRENDS, position1202)
						}
						goto l1180
					l1201:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						{
							position1218 := position
							depth++
							{
								position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1220
								}
								position++
								goto l1219
							l1220:
								position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
								if buffer[position] != rune('S') {
									goto l1217
								}
								position++
							}
						l1219:
							{
								position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1222
								}
								position++
								goto l1221
							l1222:
								position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
								if buffer[position] != rune('T') {
									goto l1217
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
									goto l1217
								}
								position++
							}
						l1223:
							{
								position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1226
								}
								position++
								goto l1225
							l1226:
								position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
								if buffer[position] != rune('B') {
									goto l1217
								}
								position++
							}
						l1225:
							{
								position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1228
								}
								position++
								goto l1227
							l1228:
								position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
								if buffer[position] != rune('E') {
									goto l1217
								}
								position++
							}
						l1227:
							{
								position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1230
								}
								position++
								goto l1229
							l1230:
								position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
								if buffer[position] != rune('F') {
									goto l1217
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
									goto l1217
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
									goto l1217
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
									goto l1217
								}
								position++
							}
						l1235:
							if !rules[ruleskip]() {
								goto l1217
							}
							depth--
							add(ruleSTRBEFORE, position1218)
						}
						goto l1180
					l1217:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						{
							position1238 := position
							depth++
							{
								position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1240
								}
								position++
								goto l1239
							l1240:
								position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
								if buffer[position] != rune('S') {
									goto l1237
								}
								position++
							}
						l1239:
							{
								position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1242
								}
								position++
								goto l1241
							l1242:
								position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
								if buffer[position] != rune('T') {
									goto l1237
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
									goto l1237
								}
								position++
							}
						l1243:
							{
								position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1246
								}
								position++
								goto l1245
							l1246:
								position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
								if buffer[position] != rune('A') {
									goto l1237
								}
								position++
							}
						l1245:
							{
								position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1248
								}
								position++
								goto l1247
							l1248:
								position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
								if buffer[position] != rune('F') {
									goto l1237
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
									goto l1237
								}
								position++
							}
						l1249:
							{
								position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1252
								}
								position++
								goto l1251
							l1252:
								position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
								if buffer[position] != rune('E') {
									goto l1237
								}
								position++
							}
						l1251:
							{
								position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1254
								}
								position++
								goto l1253
							l1254:
								position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
								if buffer[position] != rune('R') {
									goto l1237
								}
								position++
							}
						l1253:
							if !rules[ruleskip]() {
								goto l1237
							}
							depth--
							add(ruleSTRAFTER, position1238)
						}
						goto l1180
					l1237:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						{
							position1256 := position
							depth++
							{
								position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1258
								}
								position++
								goto l1257
							l1258:
								position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
								if buffer[position] != rune('S') {
									goto l1255
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
									goto l1255
								}
								position++
							}
						l1259:
							{
								position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1262
								}
								position++
								goto l1261
							l1262:
								position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
								if buffer[position] != rune('R') {
									goto l1255
								}
								position++
							}
						l1261:
							{
								position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1264
								}
								position++
								goto l1263
							l1264:
								position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
								if buffer[position] != rune('L') {
									goto l1255
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
									goto l1255
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
									goto l1255
								}
								position++
							}
						l1267:
							{
								position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1270
								}
								position++
								goto l1269
							l1270:
								position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
								if buffer[position] != rune('G') {
									goto l1255
								}
								position++
							}
						l1269:
							if !rules[ruleskip]() {
								goto l1255
							}
							depth--
							add(ruleSTRLANG, position1256)
						}
						goto l1180
					l1255:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						{
							position1272 := position
							depth++
							{
								position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1274
								}
								position++
								goto l1273
							l1274:
								position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
								if buffer[position] != rune('S') {
									goto l1271
								}
								position++
							}
						l1273:
							{
								position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1276
								}
								position++
								goto l1275
							l1276:
								position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
								if buffer[position] != rune('T') {
									goto l1271
								}
								position++
							}
						l1275:
							{
								position1277, tokenIndex1277, depth1277 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1278
								}
								position++
								goto l1277
							l1278:
								position, tokenIndex, depth = position1277, tokenIndex1277, depth1277
								if buffer[position] != rune('R') {
									goto l1271
								}
								position++
							}
						l1277:
							{
								position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1280
								}
								position++
								goto l1279
							l1280:
								position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
								if buffer[position] != rune('D') {
									goto l1271
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
									goto l1271
								}
								position++
							}
						l1281:
							if !rules[ruleskip]() {
								goto l1271
							}
							depth--
							add(ruleSTRDT, position1272)
						}
						goto l1180
					l1271:
						position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1284 := position
									depth++
									{
										position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1286
										}
										position++
										goto l1285
									l1286:
										position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
										if buffer[position] != rune('S') {
											goto l1179
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
											goto l1179
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
											goto l1179
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
											goto l1179
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
											goto l1179
										}
										position++
									}
								l1293:
									{
										position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1296
										}
										position++
										goto l1295
									l1296:
										position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
										if buffer[position] != rune('E') {
											goto l1179
										}
										position++
									}
								l1295:
									{
										position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1298
										}
										position++
										goto l1297
									l1298:
										position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
										if buffer[position] != rune('R') {
											goto l1179
										}
										position++
									}
								l1297:
									{
										position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1300
										}
										position++
										goto l1299
									l1300:
										position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
										if buffer[position] != rune('M') {
											goto l1179
										}
										position++
									}
								l1299:
									if !rules[ruleskip]() {
										goto l1179
									}
									depth--
									add(ruleSAMETERM, position1284)
								}
								break
							case 'C', 'c':
								{
									position1301 := position
									depth++
									{
										position1302, tokenIndex1302, depth1302 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1303
										}
										position++
										goto l1302
									l1303:
										position, tokenIndex, depth = position1302, tokenIndex1302, depth1302
										if buffer[position] != rune('C') {
											goto l1179
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
											goto l1179
										}
										position++
									}
								l1304:
									{
										position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1307
										}
										position++
										goto l1306
									l1307:
										position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
										if buffer[position] != rune('N') {
											goto l1179
										}
										position++
									}
								l1306:
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
											goto l1179
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
											goto l1179
										}
										position++
									}
								l1310:
									{
										position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1313
										}
										position++
										goto l1312
									l1313:
										position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
										if buffer[position] != rune('I') {
											goto l1179
										}
										position++
									}
								l1312:
									{
										position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1315
										}
										position++
										goto l1314
									l1315:
										position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
										if buffer[position] != rune('N') {
											goto l1179
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
											goto l1179
										}
										position++
									}
								l1316:
									if !rules[ruleskip]() {
										goto l1179
									}
									depth--
									add(ruleCONTAINS, position1301)
								}
								break
							default:
								{
									position1318 := position
									depth++
									{
										position1319, tokenIndex1319, depth1319 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1320
										}
										position++
										goto l1319
									l1320:
										position, tokenIndex, depth = position1319, tokenIndex1319, depth1319
										if buffer[position] != rune('L') {
											goto l1179
										}
										position++
									}
								l1319:
									{
										position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1322
										}
										position++
										goto l1321
									l1322:
										position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
										if buffer[position] != rune('A') {
											goto l1179
										}
										position++
									}
								l1321:
									{
										position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1324
										}
										position++
										goto l1323
									l1324:
										position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
										if buffer[position] != rune('N') {
											goto l1179
										}
										position++
									}
								l1323:
									{
										position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1326
										}
										position++
										goto l1325
									l1326:
										position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
										if buffer[position] != rune('G') {
											goto l1179
										}
										position++
									}
								l1325:
									{
										position1327, tokenIndex1327, depth1327 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1328
										}
										position++
										goto l1327
									l1328:
										position, tokenIndex, depth = position1327, tokenIndex1327, depth1327
										if buffer[position] != rune('M') {
											goto l1179
										}
										position++
									}
								l1327:
									{
										position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1330
										}
										position++
										goto l1329
									l1330:
										position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
										if buffer[position] != rune('A') {
											goto l1179
										}
										position++
									}
								l1329:
									{
										position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1332
										}
										position++
										goto l1331
									l1332:
										position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
										if buffer[position] != rune('T') {
											goto l1179
										}
										position++
									}
								l1331:
									{
										position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1334
										}
										position++
										goto l1333
									l1334:
										position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
										if buffer[position] != rune('C') {
											goto l1179
										}
										position++
									}
								l1333:
									{
										position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1336
										}
										position++
										goto l1335
									l1336:
										position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
										if buffer[position] != rune('H') {
											goto l1179
										}
										position++
									}
								l1335:
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
											goto l1179
										}
										position++
									}
								l1337:
									{
										position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1340
										}
										position++
										goto l1339
									l1340:
										position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
										if buffer[position] != rune('S') {
											goto l1179
										}
										position++
									}
								l1339:
									if !rules[ruleskip]() {
										goto l1179
									}
									depth--
									add(ruleLANGMATCHES, position1318)
								}
								break
							}
						}

					}
				l1180:
					if !rules[ruleLPAREN]() {
						goto l1179
					}
					if !rules[ruleexpression]() {
						goto l1179
					}
					if !rules[ruleCOMMA]() {
						goto l1179
					}
					if !rules[ruleexpression]() {
						goto l1179
					}
					if !rules[ruleRPAREN]() {
						goto l1179
					}
					goto l815
				l1179:
					position, tokenIndex, depth = position815, tokenIndex815, depth815
					{
						position1342 := position
						depth++
						{
							position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1344
							}
							position++
							goto l1343
						l1344:
							position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
							if buffer[position] != rune('B') {
								goto l1341
							}
							position++
						}
					l1343:
						{
							position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1346
							}
							position++
							goto l1345
						l1346:
							position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
							if buffer[position] != rune('O') {
								goto l1341
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
								goto l1341
							}
							position++
						}
					l1347:
						{
							position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1350
							}
							position++
							goto l1349
						l1350:
							position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
							if buffer[position] != rune('N') {
								goto l1341
							}
							position++
						}
					l1349:
						{
							position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1352
							}
							position++
							goto l1351
						l1352:
							position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
							if buffer[position] != rune('D') {
								goto l1341
							}
							position++
						}
					l1351:
						if !rules[ruleskip]() {
							goto l1341
						}
						depth--
						add(ruleBOUND, position1342)
					}
					if !rules[ruleLPAREN]() {
						goto l1341
					}
					if !rules[rulevar]() {
						goto l1341
					}
					if !rules[ruleRPAREN]() {
						goto l1341
					}
					goto l815
				l1341:
					position, tokenIndex, depth = position815, tokenIndex815, depth815
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1355 := position
								depth++
								{
									position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1357
									}
									position++
									goto l1356
								l1357:
									position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
									if buffer[position] != rune('S') {
										goto l1353
									}
									position++
								}
							l1356:
								{
									position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1359
									}
									position++
									goto l1358
								l1359:
									position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
									if buffer[position] != rune('T') {
										goto l1353
									}
									position++
								}
							l1358:
								{
									position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1361
									}
									position++
									goto l1360
								l1361:
									position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
									if buffer[position] != rune('R') {
										goto l1353
									}
									position++
								}
							l1360:
								{
									position1362, tokenIndex1362, depth1362 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1363
									}
									position++
									goto l1362
								l1363:
									position, tokenIndex, depth = position1362, tokenIndex1362, depth1362
									if buffer[position] != rune('U') {
										goto l1353
									}
									position++
								}
							l1362:
								{
									position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1365
									}
									position++
									goto l1364
								l1365:
									position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
									if buffer[position] != rune('U') {
										goto l1353
									}
									position++
								}
							l1364:
								{
									position1366, tokenIndex1366, depth1366 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1367
									}
									position++
									goto l1366
								l1367:
									position, tokenIndex, depth = position1366, tokenIndex1366, depth1366
									if buffer[position] != rune('I') {
										goto l1353
									}
									position++
								}
							l1366:
								{
									position1368, tokenIndex1368, depth1368 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1369
									}
									position++
									goto l1368
								l1369:
									position, tokenIndex, depth = position1368, tokenIndex1368, depth1368
									if buffer[position] != rune('D') {
										goto l1353
									}
									position++
								}
							l1368:
								if !rules[ruleskip]() {
									goto l1353
								}
								depth--
								add(ruleSTRUUID, position1355)
							}
							break
						case 'U', 'u':
							{
								position1370 := position
								depth++
								{
									position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1372
									}
									position++
									goto l1371
								l1372:
									position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
									if buffer[position] != rune('U') {
										goto l1353
									}
									position++
								}
							l1371:
								{
									position1373, tokenIndex1373, depth1373 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1374
									}
									position++
									goto l1373
								l1374:
									position, tokenIndex, depth = position1373, tokenIndex1373, depth1373
									if buffer[position] != rune('U') {
										goto l1353
									}
									position++
								}
							l1373:
								{
									position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1376
									}
									position++
									goto l1375
								l1376:
									position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
									if buffer[position] != rune('I') {
										goto l1353
									}
									position++
								}
							l1375:
								{
									position1377, tokenIndex1377, depth1377 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1378
									}
									position++
									goto l1377
								l1378:
									position, tokenIndex, depth = position1377, tokenIndex1377, depth1377
									if buffer[position] != rune('D') {
										goto l1353
									}
									position++
								}
							l1377:
								if !rules[ruleskip]() {
									goto l1353
								}
								depth--
								add(ruleUUID, position1370)
							}
							break
						case 'N', 'n':
							{
								position1379 := position
								depth++
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
										goto l1353
									}
									position++
								}
							l1380:
								{
									position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1383
									}
									position++
									goto l1382
								l1383:
									position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
									if buffer[position] != rune('O') {
										goto l1353
									}
									position++
								}
							l1382:
								{
									position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1385
									}
									position++
									goto l1384
								l1385:
									position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
									if buffer[position] != rune('W') {
										goto l1353
									}
									position++
								}
							l1384:
								if !rules[ruleskip]() {
									goto l1353
								}
								depth--
								add(ruleNOW, position1379)
							}
							break
						default:
							{
								position1386 := position
								depth++
								{
									position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1388
									}
									position++
									goto l1387
								l1388:
									position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
									if buffer[position] != rune('R') {
										goto l1353
									}
									position++
								}
							l1387:
								{
									position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1390
									}
									position++
									goto l1389
								l1390:
									position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
									if buffer[position] != rune('A') {
										goto l1353
									}
									position++
								}
							l1389:
								{
									position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1392
									}
									position++
									goto l1391
								l1392:
									position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
									if buffer[position] != rune('N') {
										goto l1353
									}
									position++
								}
							l1391:
								{
									position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1394
									}
									position++
									goto l1393
								l1394:
									position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
									if buffer[position] != rune('D') {
										goto l1353
									}
									position++
								}
							l1393:
								if !rules[ruleskip]() {
									goto l1353
								}
								depth--
								add(ruleRAND, position1386)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1353
					}
					goto l815
				l1353:
					position, tokenIndex, depth = position815, tokenIndex815, depth815
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
								{
									position1398 := position
									depth++
									{
										position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1400
										}
										position++
										goto l1399
									l1400:
										position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
										if buffer[position] != rune('E') {
											goto l1397
										}
										position++
									}
								l1399:
									{
										position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1402
										}
										position++
										goto l1401
									l1402:
										position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
										if buffer[position] != rune('X') {
											goto l1397
										}
										position++
									}
								l1401:
									{
										position1403, tokenIndex1403, depth1403 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1404
										}
										position++
										goto l1403
									l1404:
										position, tokenIndex, depth = position1403, tokenIndex1403, depth1403
										if buffer[position] != rune('I') {
											goto l1397
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
											goto l1397
										}
										position++
									}
								l1405:
									{
										position1407, tokenIndex1407, depth1407 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1408
										}
										position++
										goto l1407
									l1408:
										position, tokenIndex, depth = position1407, tokenIndex1407, depth1407
										if buffer[position] != rune('T') {
											goto l1397
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
											goto l1397
										}
										position++
									}
								l1409:
									if !rules[ruleskip]() {
										goto l1397
									}
									depth--
									add(ruleEXISTS, position1398)
								}
								goto l1396
							l1397:
								position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
								{
									position1411 := position
									depth++
									{
										position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1413
										}
										position++
										goto l1412
									l1413:
										position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
										if buffer[position] != rune('N') {
											goto l813
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
											goto l813
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
											goto l813
										}
										position++
									}
								l1416:
									if buffer[position] != rune(' ') {
										goto l813
									}
									position++
									{
										position1418, tokenIndex1418, depth1418 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1419
										}
										position++
										goto l1418
									l1419:
										position, tokenIndex, depth = position1418, tokenIndex1418, depth1418
										if buffer[position] != rune('E') {
											goto l813
										}
										position++
									}
								l1418:
									{
										position1420, tokenIndex1420, depth1420 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1421
										}
										position++
										goto l1420
									l1421:
										position, tokenIndex, depth = position1420, tokenIndex1420, depth1420
										if buffer[position] != rune('X') {
											goto l813
										}
										position++
									}
								l1420:
									{
										position1422, tokenIndex1422, depth1422 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1423
										}
										position++
										goto l1422
									l1423:
										position, tokenIndex, depth = position1422, tokenIndex1422, depth1422
										if buffer[position] != rune('I') {
											goto l813
										}
										position++
									}
								l1422:
									{
										position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1425
										}
										position++
										goto l1424
									l1425:
										position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
										if buffer[position] != rune('S') {
											goto l813
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
											goto l813
										}
										position++
									}
								l1426:
									{
										position1428, tokenIndex1428, depth1428 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1429
										}
										position++
										goto l1428
									l1429:
										position, tokenIndex, depth = position1428, tokenIndex1428, depth1428
										if buffer[position] != rune('S') {
											goto l813
										}
										position++
									}
								l1428:
									if !rules[ruleskip]() {
										goto l813
									}
									depth--
									add(ruleNOTEXIST, position1411)
								}
							}
						l1396:
							if !rules[rulegroupGraphPattern]() {
								goto l813
							}
							break
						case 'I', 'i':
							{
								position1430 := position
								depth++
								{
									position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1432
									}
									position++
									goto l1431
								l1432:
									position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
									if buffer[position] != rune('I') {
										goto l813
									}
									position++
								}
							l1431:
								{
									position1433, tokenIndex1433, depth1433 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1434
									}
									position++
									goto l1433
								l1434:
									position, tokenIndex, depth = position1433, tokenIndex1433, depth1433
									if buffer[position] != rune('F') {
										goto l813
									}
									position++
								}
							l1433:
								if !rules[ruleskip]() {
									goto l813
								}
								depth--
								add(ruleIF, position1430)
							}
							if !rules[ruleLPAREN]() {
								goto l813
							}
							if !rules[ruleexpression]() {
								goto l813
							}
							if !rules[ruleCOMMA]() {
								goto l813
							}
							if !rules[ruleexpression]() {
								goto l813
							}
							if !rules[ruleCOMMA]() {
								goto l813
							}
							if !rules[ruleexpression]() {
								goto l813
							}
							if !rules[ruleRPAREN]() {
								goto l813
							}
							break
						case 'C', 'c':
							{
								position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
								{
									position1437 := position
									depth++
									{
										position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1439
										}
										position++
										goto l1438
									l1439:
										position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
										if buffer[position] != rune('C') {
											goto l1436
										}
										position++
									}
								l1438:
									{
										position1440, tokenIndex1440, depth1440 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1441
										}
										position++
										goto l1440
									l1441:
										position, tokenIndex, depth = position1440, tokenIndex1440, depth1440
										if buffer[position] != rune('O') {
											goto l1436
										}
										position++
									}
								l1440:
									{
										position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1443
										}
										position++
										goto l1442
									l1443:
										position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
										if buffer[position] != rune('N') {
											goto l1436
										}
										position++
									}
								l1442:
									{
										position1444, tokenIndex1444, depth1444 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1445
										}
										position++
										goto l1444
									l1445:
										position, tokenIndex, depth = position1444, tokenIndex1444, depth1444
										if buffer[position] != rune('C') {
											goto l1436
										}
										position++
									}
								l1444:
									{
										position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1447
										}
										position++
										goto l1446
									l1447:
										position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
										if buffer[position] != rune('A') {
											goto l1436
										}
										position++
									}
								l1446:
									{
										position1448, tokenIndex1448, depth1448 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1449
										}
										position++
										goto l1448
									l1449:
										position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
										if buffer[position] != rune('T') {
											goto l1436
										}
										position++
									}
								l1448:
									if !rules[ruleskip]() {
										goto l1436
									}
									depth--
									add(ruleCONCAT, position1437)
								}
								goto l1435
							l1436:
								position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
								{
									position1450 := position
									depth++
									{
										position1451, tokenIndex1451, depth1451 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1452
										}
										position++
										goto l1451
									l1452:
										position, tokenIndex, depth = position1451, tokenIndex1451, depth1451
										if buffer[position] != rune('C') {
											goto l813
										}
										position++
									}
								l1451:
									{
										position1453, tokenIndex1453, depth1453 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1454
										}
										position++
										goto l1453
									l1454:
										position, tokenIndex, depth = position1453, tokenIndex1453, depth1453
										if buffer[position] != rune('O') {
											goto l813
										}
										position++
									}
								l1453:
									{
										position1455, tokenIndex1455, depth1455 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1456
										}
										position++
										goto l1455
									l1456:
										position, tokenIndex, depth = position1455, tokenIndex1455, depth1455
										if buffer[position] != rune('A') {
											goto l813
										}
										position++
									}
								l1455:
									{
										position1457, tokenIndex1457, depth1457 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1458
										}
										position++
										goto l1457
									l1458:
										position, tokenIndex, depth = position1457, tokenIndex1457, depth1457
										if buffer[position] != rune('L') {
											goto l813
										}
										position++
									}
								l1457:
									{
										position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1460
										}
										position++
										goto l1459
									l1460:
										position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
										if buffer[position] != rune('E') {
											goto l813
										}
										position++
									}
								l1459:
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
											goto l813
										}
										position++
									}
								l1461:
									{
										position1463, tokenIndex1463, depth1463 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1464
										}
										position++
										goto l1463
									l1464:
										position, tokenIndex, depth = position1463, tokenIndex1463, depth1463
										if buffer[position] != rune('C') {
											goto l813
										}
										position++
									}
								l1463:
									{
										position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1466
										}
										position++
										goto l1465
									l1466:
										position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
										if buffer[position] != rune('E') {
											goto l813
										}
										position++
									}
								l1465:
									if !rules[ruleskip]() {
										goto l813
									}
									depth--
									add(ruleCOALESCE, position1450)
								}
							}
						l1435:
							if !rules[ruleargList]() {
								goto l813
							}
							break
						case 'B', 'b':
							{
								position1467 := position
								depth++
								{
									position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1469
									}
									position++
									goto l1468
								l1469:
									position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
									if buffer[position] != rune('B') {
										goto l813
									}
									position++
								}
							l1468:
								{
									position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1471
									}
									position++
									goto l1470
								l1471:
									position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
									if buffer[position] != rune('N') {
										goto l813
									}
									position++
								}
							l1470:
								{
									position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1473
									}
									position++
									goto l1472
								l1473:
									position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
									if buffer[position] != rune('O') {
										goto l813
									}
									position++
								}
							l1472:
								{
									position1474, tokenIndex1474, depth1474 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1475
									}
									position++
									goto l1474
								l1475:
									position, tokenIndex, depth = position1474, tokenIndex1474, depth1474
									if buffer[position] != rune('D') {
										goto l813
									}
									position++
								}
							l1474:
								{
									position1476, tokenIndex1476, depth1476 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1477
									}
									position++
									goto l1476
								l1477:
									position, tokenIndex, depth = position1476, tokenIndex1476, depth1476
									if buffer[position] != rune('E') {
										goto l813
									}
									position++
								}
							l1476:
								if !rules[ruleskip]() {
									goto l813
								}
								depth--
								add(ruleBNODE, position1467)
							}
							{
								position1478, tokenIndex1478, depth1478 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1479
								}
								if !rules[ruleexpression]() {
									goto l1479
								}
								if !rules[ruleRPAREN]() {
									goto l1479
								}
								goto l1478
							l1479:
								position, tokenIndex, depth = position1478, tokenIndex1478, depth1478
								if !rules[rulenil]() {
									goto l813
								}
							}
						l1478:
							break
						default:
							{
								position1480, tokenIndex1480, depth1480 := position, tokenIndex, depth
								{
									position1482 := position
									depth++
									{
										position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1484
										}
										position++
										goto l1483
									l1484:
										position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
										if buffer[position] != rune('S') {
											goto l1481
										}
										position++
									}
								l1483:
									{
										position1485, tokenIndex1485, depth1485 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1486
										}
										position++
										goto l1485
									l1486:
										position, tokenIndex, depth = position1485, tokenIndex1485, depth1485
										if buffer[position] != rune('U') {
											goto l1481
										}
										position++
									}
								l1485:
									{
										position1487, tokenIndex1487, depth1487 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1488
										}
										position++
										goto l1487
									l1488:
										position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
										if buffer[position] != rune('B') {
											goto l1481
										}
										position++
									}
								l1487:
									{
										position1489, tokenIndex1489, depth1489 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1490
										}
										position++
										goto l1489
									l1490:
										position, tokenIndex, depth = position1489, tokenIndex1489, depth1489
										if buffer[position] != rune('S') {
											goto l1481
										}
										position++
									}
								l1489:
									{
										position1491, tokenIndex1491, depth1491 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1492
										}
										position++
										goto l1491
									l1492:
										position, tokenIndex, depth = position1491, tokenIndex1491, depth1491
										if buffer[position] != rune('T') {
											goto l1481
										}
										position++
									}
								l1491:
									{
										position1493, tokenIndex1493, depth1493 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1494
										}
										position++
										goto l1493
									l1494:
										position, tokenIndex, depth = position1493, tokenIndex1493, depth1493
										if buffer[position] != rune('R') {
											goto l1481
										}
										position++
									}
								l1493:
									if !rules[ruleskip]() {
										goto l1481
									}
									depth--
									add(ruleSUBSTR, position1482)
								}
								goto l1480
							l1481:
								position, tokenIndex, depth = position1480, tokenIndex1480, depth1480
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
											goto l1495
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
											goto l1495
										}
										position++
									}
								l1499:
									{
										position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1502
										}
										position++
										goto l1501
									l1502:
										position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
										if buffer[position] != rune('P') {
											goto l1495
										}
										position++
									}
								l1501:
									{
										position1503, tokenIndex1503, depth1503 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1504
										}
										position++
										goto l1503
									l1504:
										position, tokenIndex, depth = position1503, tokenIndex1503, depth1503
										if buffer[position] != rune('L') {
											goto l1495
										}
										position++
									}
								l1503:
									{
										position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1506
										}
										position++
										goto l1505
									l1506:
										position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
										if buffer[position] != rune('A') {
											goto l1495
										}
										position++
									}
								l1505:
									{
										position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1508
										}
										position++
										goto l1507
									l1508:
										position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
										if buffer[position] != rune('C') {
											goto l1495
										}
										position++
									}
								l1507:
									{
										position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1510
										}
										position++
										goto l1509
									l1510:
										position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
										if buffer[position] != rune('E') {
											goto l1495
										}
										position++
									}
								l1509:
									if !rules[ruleskip]() {
										goto l1495
									}
									depth--
									add(ruleREPLACE, position1496)
								}
								goto l1480
							l1495:
								position, tokenIndex, depth = position1480, tokenIndex1480, depth1480
								{
									position1511 := position
									depth++
									{
										position1512, tokenIndex1512, depth1512 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1513
										}
										position++
										goto l1512
									l1513:
										position, tokenIndex, depth = position1512, tokenIndex1512, depth1512
										if buffer[position] != rune('R') {
											goto l813
										}
										position++
									}
								l1512:
									{
										position1514, tokenIndex1514, depth1514 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1515
										}
										position++
										goto l1514
									l1515:
										position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
										if buffer[position] != rune('E') {
											goto l813
										}
										position++
									}
								l1514:
									{
										position1516, tokenIndex1516, depth1516 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1517
										}
										position++
										goto l1516
									l1517:
										position, tokenIndex, depth = position1516, tokenIndex1516, depth1516
										if buffer[position] != rune('G') {
											goto l813
										}
										position++
									}
								l1516:
									{
										position1518, tokenIndex1518, depth1518 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1519
										}
										position++
										goto l1518
									l1519:
										position, tokenIndex, depth = position1518, tokenIndex1518, depth1518
										if buffer[position] != rune('E') {
											goto l813
										}
										position++
									}
								l1518:
									{
										position1520, tokenIndex1520, depth1520 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1521
										}
										position++
										goto l1520
									l1521:
										position, tokenIndex, depth = position1520, tokenIndex1520, depth1520
										if buffer[position] != rune('X') {
											goto l813
										}
										position++
									}
								l1520:
									if !rules[ruleskip]() {
										goto l813
									}
									depth--
									add(ruleREGEX, position1511)
								}
							}
						l1480:
							if !rules[ruleLPAREN]() {
								goto l813
							}
							if !rules[ruleexpression]() {
								goto l813
							}
							if !rules[ruleCOMMA]() {
								goto l813
							}
							if !rules[ruleexpression]() {
								goto l813
							}
							{
								position1522, tokenIndex1522, depth1522 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1522
								}
								if !rules[ruleexpression]() {
									goto l1522
								}
								goto l1523
							l1522:
								position, tokenIndex, depth = position1522, tokenIndex1522, depth1522
							}
						l1523:
							if !rules[ruleRPAREN]() {
								goto l813
							}
							break
						}
					}

				}
			l815:
				depth--
				add(rulebuiltinCall, position814)
			}
			return true
		l813:
			position, tokenIndex, depth = position813, tokenIndex813, depth813
			return false
		},
		/* 69 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
			{
				position1525 := position
				depth++
				{
					position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
					{
						position1528 := position
						depth++
					l1529:
						{
							position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
							{
								position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1532
								}
								position++
								goto l1531
							l1532:
								position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1530
								}
								position++
							}
						l1531:
							goto l1529
						l1530:
							position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
						}
						depth--
						add(rulePegText, position1528)
					}
					if buffer[position] != rune(':') {
						goto l1527
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1526
				l1527:
					position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
					{
						position1535 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1534
						}
						position++
					l1536:
						{
							position1537, tokenIndex1537, depth1537 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1537
							}
							position++
							goto l1536
						l1537:
							position, tokenIndex, depth = position1537, tokenIndex1537, depth1537
						}
						depth--
						add(rulePegText, position1535)
					}
					if buffer[position] != rune('/') {
						goto l1534
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1526
				l1534:
					position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
					{
						position1539 := position
						depth++
					l1540:
						{
							position1541, tokenIndex1541, depth1541 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1541
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1541
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1541
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1541
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1541
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1541
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1541
									}
									position++
									break
								}
							}

							goto l1540
						l1541:
							position, tokenIndex, depth = position1541, tokenIndex1541, depth1541
						}
						depth--
						add(rulePegText, position1539)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1526:
				if buffer[position] != rune('<') {
					goto l1524
				}
				position++
				if !rules[rulews]() {
					goto l1524
				}
				if !rules[ruleskip]() {
					goto l1524
				}
				depth--
				add(rulepof, position1525)
			}
			return true
		l1524:
			position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
			return false
		},
		/* 70 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1544, tokenIndex1544, depth1544 := position, tokenIndex, depth
			{
				position1545 := position
				depth++
				{
					position1546, tokenIndex1546, depth1546 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1547
					}
					position++
					goto l1546
				l1547:
					position, tokenIndex, depth = position1546, tokenIndex1546, depth1546
					if buffer[position] != rune('$') {
						goto l1544
					}
					position++
				}
			l1546:
				{
					position1548 := position
					depth++
					{
						position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
						if !rules[rulepnCharsU]() {
							goto l1550
						}
						goto l1549
					l1550:
						position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1544
						}
						position++
					}
				l1549:
				l1551:
					{
						position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
						{
							position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
							if !rules[rulepnCharsU]() {
								goto l1554
							}
							goto l1553
						l1554:
							position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
							{
								switch buffer[position] {
								case 'â':
									if c := buffer[position]; c < rune('‿') || c > rune('⁀') {
										goto l1552
									}
									position++
									break
								case 'Ì', 'Í':
									if c := buffer[position]; c < rune('̀') || c > rune('ͯ') {
										goto l1552
									}
									position++
									break
								case 'Â':
									if buffer[position] != rune('·') {
										goto l1552
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1552
									}
									position++
									break
								}
							}

						}
					l1553:
						goto l1551
					l1552:
						position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
					}
					depth--
					add(ruleVARNAME, position1548)
				}
				if !rules[ruleskip]() {
					goto l1544
				}
				depth--
				add(rulevar, position1545)
			}
			return true
		l1544:
			position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
			return false
		},
		/* 71 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
			{
				position1557 := position
				depth++
				{
					position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1559
					}
					goto l1558
				l1559:
					position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
					{
						position1560 := position
						depth++
						{
							position1561, tokenIndex1561, depth1561 := position, tokenIndex, depth
							if !rules[rulepnPrefix]() {
								goto l1561
							}
							goto l1562
						l1561:
							position, tokenIndex, depth = position1561, tokenIndex1561, depth1561
						}
					l1562:
						if buffer[position] != rune(':') {
							goto l1556
						}
						position++
						{
							position1563 := position
							depth++
							{
								switch buffer[position] {
								case '%', '\\':
									{
										position1567 := position
										depth++
										{
											position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
											{
												position1570 := position
												depth++
												if buffer[position] != rune('%') {
													goto l1569
												}
												position++
												if !rules[rulehex]() {
													goto l1569
												}
												if !rules[rulehex]() {
													goto l1569
												}
												depth--
												add(rulepercent, position1570)
											}
											goto l1568
										l1569:
											position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
											{
												position1571 := position
												depth++
												if buffer[position] != rune('\\') {
													goto l1556
												}
												position++
												{
													switch buffer[position] {
													case '%':
														if buffer[position] != rune('%') {
															goto l1556
														}
														position++
														break
													case '@':
														if buffer[position] != rune('@') {
															goto l1556
														}
														position++
														break
													case '#':
														if buffer[position] != rune('#') {
															goto l1556
														}
														position++
														break
													case '?':
														if buffer[position] != rune('?') {
															goto l1556
														}
														position++
														break
													case '/':
														if buffer[position] != rune('/') {
															goto l1556
														}
														position++
														break
													case '=':
														if buffer[position] != rune('=') {
															goto l1556
														}
														position++
														break
													case ';':
														if buffer[position] != rune(';') {
															goto l1556
														}
														position++
														break
													case ',':
														if buffer[position] != rune(',') {
															goto l1556
														}
														position++
														break
													case '+':
														if buffer[position] != rune('+') {
															goto l1556
														}
														position++
														break
													case '*':
														if buffer[position] != rune('*') {
															goto l1556
														}
														position++
														break
													case ')':
														if buffer[position] != rune(')') {
															goto l1556
														}
														position++
														break
													case '(':
														if buffer[position] != rune('(') {
															goto l1556
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l1556
														}
														position++
														break
													case '&':
														if buffer[position] != rune('&') {
															goto l1556
														}
														position++
														break
													case '$':
														if buffer[position] != rune('$') {
															goto l1556
														}
														position++
														break
													case '!':
														if buffer[position] != rune('!') {
															goto l1556
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l1556
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l1556
														}
														position++
														break
													case '~':
														if buffer[position] != rune('~') {
															goto l1556
														}
														position++
														break
													default:
														if buffer[position] != rune('_') {
															goto l1556
														}
														position++
														break
													}
												}

												depth--
												add(rulepnLocalEsc, position1571)
											}
										}
									l1568:
										depth--
										add(ruleplx, position1567)
									}
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1556
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l1556
									}
									position++
									break
								default:
									if !rules[rulepnCharsU]() {
										goto l1556
									}
									break
								}
							}

						l1564:
							{
								position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '%', '\\':
										{
											position1574 := position
											depth++
											{
												position1575, tokenIndex1575, depth1575 := position, tokenIndex, depth
												{
													position1577 := position
													depth++
													if buffer[position] != rune('%') {
														goto l1576
													}
													position++
													if !rules[rulehex]() {
														goto l1576
													}
													if !rules[rulehex]() {
														goto l1576
													}
													depth--
													add(rulepercent, position1577)
												}
												goto l1575
											l1576:
												position, tokenIndex, depth = position1575, tokenIndex1575, depth1575
												{
													position1578 := position
													depth++
													if buffer[position] != rune('\\') {
														goto l1565
													}
													position++
													{
														switch buffer[position] {
														case '%':
															if buffer[position] != rune('%') {
																goto l1565
															}
															position++
															break
														case '@':
															if buffer[position] != rune('@') {
																goto l1565
															}
															position++
															break
														case '#':
															if buffer[position] != rune('#') {
																goto l1565
															}
															position++
															break
														case '?':
															if buffer[position] != rune('?') {
																goto l1565
															}
															position++
															break
														case '/':
															if buffer[position] != rune('/') {
																goto l1565
															}
															position++
															break
														case '=':
															if buffer[position] != rune('=') {
																goto l1565
															}
															position++
															break
														case ';':
															if buffer[position] != rune(';') {
																goto l1565
															}
															position++
															break
														case ',':
															if buffer[position] != rune(',') {
																goto l1565
															}
															position++
															break
														case '+':
															if buffer[position] != rune('+') {
																goto l1565
															}
															position++
															break
														case '*':
															if buffer[position] != rune('*') {
																goto l1565
															}
															position++
															break
														case ')':
															if buffer[position] != rune(')') {
																goto l1565
															}
															position++
															break
														case '(':
															if buffer[position] != rune('(') {
																goto l1565
															}
															position++
															break
														case '\'':
															if buffer[position] != rune('\'') {
																goto l1565
															}
															position++
															break
														case '&':
															if buffer[position] != rune('&') {
																goto l1565
															}
															position++
															break
														case '$':
															if buffer[position] != rune('$') {
																goto l1565
															}
															position++
															break
														case '!':
															if buffer[position] != rune('!') {
																goto l1565
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l1565
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l1565
															}
															position++
															break
														case '~':
															if buffer[position] != rune('~') {
																goto l1565
															}
															position++
															break
														default:
															if buffer[position] != rune('_') {
																goto l1565
															}
															position++
															break
														}
													}

													depth--
													add(rulepnLocalEsc, position1578)
												}
											}
										l1575:
											depth--
											add(ruleplx, position1574)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1565
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1565
										}
										position++
										break
									default:
										if !rules[rulepnCharsU]() {
											goto l1565
										}
										break
									}
								}

								goto l1564
							l1565:
								position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
							}
							depth--
							add(rulepnLocal, position1563)
						}
						if !rules[ruleskip]() {
							goto l1556
						}
						depth--
						add(ruleprefixedName, position1560)
					}
				}
			l1558:
				depth--
				add(ruleiriref, position1557)
			}
			return true
		l1556:
			position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
			return false
		},
		/* 72 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1580, tokenIndex1580, depth1580 := position, tokenIndex, depth
			{
				position1581 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1580
				}
				position++
			l1582:
				{
					position1583, tokenIndex1583, depth1583 := position, tokenIndex, depth
					{
						position1584, tokenIndex1584, depth1584 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1584
						}
						position++
						goto l1583
					l1584:
						position, tokenIndex, depth = position1584, tokenIndex1584, depth1584
					}
					if !matchDot() {
						goto l1583
					}
					goto l1582
				l1583:
					position, tokenIndex, depth = position1583, tokenIndex1583, depth1583
				}
				if buffer[position] != rune('>') {
					goto l1580
				}
				position++
				if !rules[ruleskip]() {
					goto l1580
				}
				depth--
				add(ruleiri, position1581)
			}
			return true
		l1580:
			position, tokenIndex, depth = position1580, tokenIndex1580, depth1580
			return false
		},
		/* 73 prefixedName <- <(pnPrefix? ':' pnLocal skip)> */
		nil,
		/* 74 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1586, tokenIndex1586, depth1586 := position, tokenIndex, depth
			{
				position1587 := position
				depth++
				if !rules[rulestring]() {
					goto l1586
				}
				{
					position1588, tokenIndex1588, depth1588 := position, tokenIndex, depth
					{
						position1590, tokenIndex1590, depth1590 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1591
						}
						position++
						{
							position1594, tokenIndex1594, depth1594 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1595
							}
							position++
							goto l1594
						l1595:
							position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1591
							}
							position++
						}
					l1594:
					l1592:
						{
							position1593, tokenIndex1593, depth1593 := position, tokenIndex, depth
							{
								position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1597
								}
								position++
								goto l1596
							l1597:
								position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1593
								}
								position++
							}
						l1596:
							goto l1592
						l1593:
							position, tokenIndex, depth = position1593, tokenIndex1593, depth1593
						}
					l1598:
						{
							position1599, tokenIndex1599, depth1599 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1599
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1599
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1599
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1599
									}
									position++
									break
								}
							}

						l1600:
							{
								position1601, tokenIndex1601, depth1601 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1601
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1601
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1601
										}
										position++
										break
									}
								}

								goto l1600
							l1601:
								position, tokenIndex, depth = position1601, tokenIndex1601, depth1601
							}
							goto l1598
						l1599:
							position, tokenIndex, depth = position1599, tokenIndex1599, depth1599
						}
						goto l1590
					l1591:
						position, tokenIndex, depth = position1590, tokenIndex1590, depth1590
						if buffer[position] != rune('^') {
							goto l1588
						}
						position++
						if buffer[position] != rune('^') {
							goto l1588
						}
						position++
						if !rules[ruleiriref]() {
							goto l1588
						}
					}
				l1590:
					goto l1589
				l1588:
					position, tokenIndex, depth = position1588, tokenIndex1588, depth1588
				}
			l1589:
				if !rules[ruleskip]() {
					goto l1586
				}
				depth--
				add(ruleliteral, position1587)
			}
			return true
		l1586:
			position, tokenIndex, depth = position1586, tokenIndex1586, depth1586
			return false
		},
		/* 75 string <- <(stringLiteralA / stringLiteralB / stringLiteralLongA / stringLiteralLongB)> */
		func() bool {
			position1604, tokenIndex1604, depth1604 := position, tokenIndex, depth
			{
				position1605 := position
				depth++
				{
					position1606, tokenIndex1606, depth1606 := position, tokenIndex, depth
					{
						position1608 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1607
						}
						position++
					l1609:
						{
							position1610, tokenIndex1610, depth1610 := position, tokenIndex, depth
							{
								position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
								{
									position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1613
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1613
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1613
											}
											position++
											break
										default:
											if buffer[position] != rune('\'') {
												goto l1613
											}
											position++
											break
										}
									}

									goto l1612
								l1613:
									position, tokenIndex, depth = position1613, tokenIndex1613, depth1613
								}
								if !matchDot() {
									goto l1612
								}
								goto l1611
							l1612:
								position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
								if !rules[ruleechar]() {
									goto l1610
								}
							}
						l1611:
							goto l1609
						l1610:
							position, tokenIndex, depth = position1610, tokenIndex1610, depth1610
						}
						if buffer[position] != rune('\'') {
							goto l1607
						}
						position++
						depth--
						add(rulestringLiteralA, position1608)
					}
					goto l1606
				l1607:
					position, tokenIndex, depth = position1606, tokenIndex1606, depth1606
					{
						position1616 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1615
						}
						position++
					l1617:
						{
							position1618, tokenIndex1618, depth1618 := position, tokenIndex, depth
							{
								position1619, tokenIndex1619, depth1619 := position, tokenIndex, depth
								{
									position1621, tokenIndex1621, depth1621 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1621
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1621
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1621
											}
											position++
											break
										default:
											if buffer[position] != rune('"') {
												goto l1621
											}
											position++
											break
										}
									}

									goto l1620
								l1621:
									position, tokenIndex, depth = position1621, tokenIndex1621, depth1621
								}
								if !matchDot() {
									goto l1620
								}
								goto l1619
							l1620:
								position, tokenIndex, depth = position1619, tokenIndex1619, depth1619
								if !rules[ruleechar]() {
									goto l1618
								}
							}
						l1619:
							goto l1617
						l1618:
							position, tokenIndex, depth = position1618, tokenIndex1618, depth1618
						}
						if buffer[position] != rune('"') {
							goto l1615
						}
						position++
						depth--
						add(rulestringLiteralB, position1616)
					}
					goto l1606
				l1615:
					position, tokenIndex, depth = position1606, tokenIndex1606, depth1606
					{
						position1624 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1623
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1623
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1623
						}
						position++
					l1625:
						{
							position1626, tokenIndex1626, depth1626 := position, tokenIndex, depth
							{
								position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
								{
									position1629, tokenIndex1629, depth1629 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l1630
									}
									position++
									goto l1629
								l1630:
									position, tokenIndex, depth = position1629, tokenIndex1629, depth1629
									if buffer[position] != rune('\'') {
										goto l1627
									}
									position++
									if buffer[position] != rune('\'') {
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
							{
								position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
								{
									position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
									{
										position1634, tokenIndex1634, depth1634 := position, tokenIndex, depth
										if buffer[position] != rune('\'') {
											goto l1635
										}
										position++
										goto l1634
									l1635:
										position, tokenIndex, depth = position1634, tokenIndex1634, depth1634
										if buffer[position] != rune('\\') {
											goto l1633
										}
										position++
									}
								l1634:
									goto l1632
								l1633:
									position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
								}
								if !matchDot() {
									goto l1632
								}
								goto l1631
							l1632:
								position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
								if !rules[ruleechar]() {
									goto l1626
								}
							}
						l1631:
							goto l1625
						l1626:
							position, tokenIndex, depth = position1626, tokenIndex1626, depth1626
						}
						if buffer[position] != rune('\'') {
							goto l1623
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1623
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1623
						}
						position++
						depth--
						add(rulestringLiteralLongA, position1624)
					}
					goto l1606
				l1623:
					position, tokenIndex, depth = position1606, tokenIndex1606, depth1606
					{
						position1636 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1604
						}
						position++
						if buffer[position] != rune('"') {
							goto l1604
						}
						position++
						if buffer[position] != rune('"') {
							goto l1604
						}
						position++
					l1637:
						{
							position1638, tokenIndex1638, depth1638 := position, tokenIndex, depth
							{
								position1639, tokenIndex1639, depth1639 := position, tokenIndex, depth
								{
									position1641, tokenIndex1641, depth1641 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l1642
									}
									position++
									goto l1641
								l1642:
									position, tokenIndex, depth = position1641, tokenIndex1641, depth1641
									if buffer[position] != rune('"') {
										goto l1639
									}
									position++
									if buffer[position] != rune('"') {
										goto l1639
									}
									position++
								}
							l1641:
								goto l1640
							l1639:
								position, tokenIndex, depth = position1639, tokenIndex1639, depth1639
							}
						l1640:
							{
								position1643, tokenIndex1643, depth1643 := position, tokenIndex, depth
								{
									position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
									{
										position1646, tokenIndex1646, depth1646 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											goto l1647
										}
										position++
										goto l1646
									l1647:
										position, tokenIndex, depth = position1646, tokenIndex1646, depth1646
										if buffer[position] != rune('\\') {
											goto l1645
										}
										position++
									}
								l1646:
									goto l1644
								l1645:
									position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
								}
								if !matchDot() {
									goto l1644
								}
								goto l1643
							l1644:
								position, tokenIndex, depth = position1643, tokenIndex1643, depth1643
								if !rules[ruleechar]() {
									goto l1638
								}
							}
						l1643:
							goto l1637
						l1638:
							position, tokenIndex, depth = position1638, tokenIndex1638, depth1638
						}
						if buffer[position] != rune('"') {
							goto l1604
						}
						position++
						if buffer[position] != rune('"') {
							goto l1604
						}
						position++
						if buffer[position] != rune('"') {
							goto l1604
						}
						position++
						depth--
						add(rulestringLiteralLongB, position1636)
					}
				}
			l1606:
				depth--
				add(rulestring, position1605)
			}
			return true
		l1604:
			position, tokenIndex, depth = position1604, tokenIndex1604, depth1604
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
			position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
			{
				position1653 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l1652
				}
				position++
				{
					switch buffer[position] {
					case '\'':
						if buffer[position] != rune('\'') {
							goto l1652
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l1652
						}
						position++
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l1652
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l1652
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l1652
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l1652
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l1652
						}
						position++
						break
					default:
						if buffer[position] != rune('t') {
							goto l1652
						}
						position++
						break
					}
				}

				depth--
				add(ruleechar, position1653)
			}
			return true
		l1652:
			position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
			return false
		},
		/* 81 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1655, tokenIndex1655, depth1655 := position, tokenIndex, depth
			{
				position1656 := position
				depth++
				{
					position1657, tokenIndex1657, depth1657 := position, tokenIndex, depth
					{
						position1659, tokenIndex1659, depth1659 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1660
						}
						position++
						goto l1659
					l1660:
						position, tokenIndex, depth = position1659, tokenIndex1659, depth1659
						if buffer[position] != rune('-') {
							goto l1657
						}
						position++
					}
				l1659:
					goto l1658
				l1657:
					position, tokenIndex, depth = position1657, tokenIndex1657, depth1657
				}
			l1658:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1655
				}
				position++
			l1661:
				{
					position1662, tokenIndex1662, depth1662 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1662
					}
					position++
					goto l1661
				l1662:
					position, tokenIndex, depth = position1662, tokenIndex1662, depth1662
				}
				{
					position1663, tokenIndex1663, depth1663 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1663
					}
					position++
				l1665:
					{
						position1666, tokenIndex1666, depth1666 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1666
						}
						position++
						goto l1665
					l1666:
						position, tokenIndex, depth = position1666, tokenIndex1666, depth1666
					}
					goto l1664
				l1663:
					position, tokenIndex, depth = position1663, tokenIndex1663, depth1663
				}
			l1664:
				if !rules[ruleskip]() {
					goto l1655
				}
				depth--
				add(rulenumericLiteral, position1656)
			}
			return true
		l1655:
			position, tokenIndex, depth = position1655, tokenIndex1655, depth1655
			return false
		},
		/* 82 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 83 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1668, tokenIndex1668, depth1668 := position, tokenIndex, depth
			{
				position1669 := position
				depth++
				{
					position1670, tokenIndex1670, depth1670 := position, tokenIndex, depth
					{
						position1672 := position
						depth++
						{
							position1673, tokenIndex1673, depth1673 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1674
							}
							position++
							goto l1673
						l1674:
							position, tokenIndex, depth = position1673, tokenIndex1673, depth1673
							if buffer[position] != rune('T') {
								goto l1671
							}
							position++
						}
					l1673:
						{
							position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1676
							}
							position++
							goto l1675
						l1676:
							position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
							if buffer[position] != rune('R') {
								goto l1671
							}
							position++
						}
					l1675:
						{
							position1677, tokenIndex1677, depth1677 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1678
							}
							position++
							goto l1677
						l1678:
							position, tokenIndex, depth = position1677, tokenIndex1677, depth1677
							if buffer[position] != rune('U') {
								goto l1671
							}
							position++
						}
					l1677:
						{
							position1679, tokenIndex1679, depth1679 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1680
							}
							position++
							goto l1679
						l1680:
							position, tokenIndex, depth = position1679, tokenIndex1679, depth1679
							if buffer[position] != rune('E') {
								goto l1671
							}
							position++
						}
					l1679:
						if !rules[ruleskip]() {
							goto l1671
						}
						depth--
						add(ruleTRUE, position1672)
					}
					goto l1670
				l1671:
					position, tokenIndex, depth = position1670, tokenIndex1670, depth1670
					{
						position1681 := position
						depth++
						{
							position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1683
							}
							position++
							goto l1682
						l1683:
							position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
							if buffer[position] != rune('F') {
								goto l1668
							}
							position++
						}
					l1682:
						{
							position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1685
							}
							position++
							goto l1684
						l1685:
							position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
							if buffer[position] != rune('A') {
								goto l1668
							}
							position++
						}
					l1684:
						{
							position1686, tokenIndex1686, depth1686 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1687
							}
							position++
							goto l1686
						l1687:
							position, tokenIndex, depth = position1686, tokenIndex1686, depth1686
							if buffer[position] != rune('L') {
								goto l1668
							}
							position++
						}
					l1686:
						{
							position1688, tokenIndex1688, depth1688 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1689
							}
							position++
							goto l1688
						l1689:
							position, tokenIndex, depth = position1688, tokenIndex1688, depth1688
							if buffer[position] != rune('S') {
								goto l1668
							}
							position++
						}
					l1688:
						{
							position1690, tokenIndex1690, depth1690 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1691
							}
							position++
							goto l1690
						l1691:
							position, tokenIndex, depth = position1690, tokenIndex1690, depth1690
							if buffer[position] != rune('E') {
								goto l1668
							}
							position++
						}
					l1690:
						if !rules[ruleskip]() {
							goto l1668
						}
						depth--
						add(ruleFALSE, position1681)
					}
				}
			l1670:
				depth--
				add(rulebooleanLiteral, position1669)
			}
			return true
		l1668:
			position, tokenIndex, depth = position1668, tokenIndex1668, depth1668
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
			position1695, tokenIndex1695, depth1695 := position, tokenIndex, depth
			{
				position1696 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1695
				}
				position++
			l1697:
				{
					position1698, tokenIndex1698, depth1698 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1698
					}
					goto l1697
				l1698:
					position, tokenIndex, depth = position1698, tokenIndex1698, depth1698
				}
				if buffer[position] != rune(')') {
					goto l1695
				}
				position++
				if !rules[ruleskip]() {
					goto l1695
				}
				depth--
				add(rulenil, position1696)
			}
			return true
		l1695:
			position, tokenIndex, depth = position1695, tokenIndex1695, depth1695
			return false
		},
		/* 88 VARNAME <- <((pnCharsU / [0-9]) (pnCharsU / ((&('â') [‿-⁀]) | (&('Ì' | 'Í') [̀-ͯ]) | (&('Â') '·') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))*)> */
		nil,
		/* 89 pnPrefix <- <(pnCharsBase pnChars*)> */
		func() bool {
			position1700, tokenIndex1700, depth1700 := position, tokenIndex, depth
			{
				position1701 := position
				depth++
				if !rules[rulepnCharsBase]() {
					goto l1700
				}
			l1702:
				{
					position1703, tokenIndex1703, depth1703 := position, tokenIndex, depth
					{
						position1704 := position
						depth++
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1703
								}
								position++
								break
							case '-':
								if buffer[position] != rune('-') {
									goto l1703
								}
								position++
								break
							default:
								if !rules[rulepnCharsU]() {
									goto l1703
								}
								break
							}
						}

						depth--
						add(rulepnChars, position1704)
					}
					goto l1702
				l1703:
					position, tokenIndex, depth = position1703, tokenIndex1703, depth1703
				}
				depth--
				add(rulepnPrefix, position1701)
			}
			return true
		l1700:
			position, tokenIndex, depth = position1700, tokenIndex1700, depth1700
			return false
		},
		/* 90 pnLocal <- <((&('%' | '\\') plx) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&(':') ':') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))+> */
		nil,
		/* 91 pnChars <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))> */
		nil,
		/* 92 pnCharsU <- <(pnCharsBase / '_')> */
		func() bool {
			position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
			{
				position1709 := position
				depth++
				{
					position1710, tokenIndex1710, depth1710 := position, tokenIndex, depth
					if !rules[rulepnCharsBase]() {
						goto l1711
					}
					goto l1710
				l1711:
					position, tokenIndex, depth = position1710, tokenIndex1710, depth1710
					if buffer[position] != rune('_') {
						goto l1708
					}
					position++
				}
			l1710:
				depth--
				add(rulepnCharsU, position1709)
			}
			return true
		l1708:
			position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
			return false
		},
		/* 93 pnCharsBase <- <([À-Ö] / [Ø-ö] / [Ͱ-ͽ] / [‌-‍] / [⁰-↏] / [豈-﷏] / ((&('ð' | 'ñ' | 'ò' | 'ó') [𐀀-󯿿]) | (&('ï') [ﷰ-�]) | (&('ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í') [、-퟿]) | (&('â') [Ⰰ-⿯]) | (&('Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á') [Ϳ-῿]) | (&('Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë') [ø-˿]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))> */
		func() bool {
			position1712, tokenIndex1712, depth1712 := position, tokenIndex, depth
			{
				position1713 := position
				depth++
				{
					position1714, tokenIndex1714, depth1714 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('À') || c > rune('Ö') {
						goto l1715
					}
					position++
					goto l1714
				l1715:
					position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
					if c := buffer[position]; c < rune('Ø') || c > rune('ö') {
						goto l1716
					}
					position++
					goto l1714
				l1716:
					position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
					if c := buffer[position]; c < rune('Ͱ') || c > rune('ͽ') {
						goto l1717
					}
					position++
					goto l1714
				l1717:
					position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
					if c := buffer[position]; c < rune('\u200c') || c > rune('\u200d') {
						goto l1718
					}
					position++
					goto l1714
				l1718:
					position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
					if c := buffer[position]; c < rune('⁰') || c > rune('\u218f') {
						goto l1719
					}
					position++
					goto l1714
				l1719:
					position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
					if c := buffer[position]; c < rune('豈') || c > rune('\ufdcf') {
						goto l1720
					}
					position++
					goto l1714
				l1720:
					position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
					{
						switch buffer[position] {
						case 'ð', 'ñ', 'ò', 'ó':
							if c := buffer[position]; c < rune('𐀀') || c > rune('\U000effff') {
								goto l1712
							}
							position++
							break
						case 'ï':
							if c := buffer[position]; c < rune('ﷰ') || c > rune('�') {
								goto l1712
							}
							position++
							break
						case 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í':
							if c := buffer[position]; c < rune('、') || c > rune('\ud7ff') {
								goto l1712
							}
							position++
							break
						case 'â':
							if c := buffer[position]; c < rune('Ⰰ') || c > rune('\u2fef') {
								goto l1712
							}
							position++
							break
						case 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'Þ', 'ß', 'à', 'á':
							if c := buffer[position]; c < rune('\u037f') || c > rune('\u1fff') {
								goto l1712
							}
							position++
							break
						case 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ê', 'Ë':
							if c := buffer[position]; c < rune('ø') || c > rune('˿') {
								goto l1712
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1712
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1712
							}
							position++
							break
						}
					}

				}
			l1714:
				depth--
				add(rulepnCharsBase, position1713)
			}
			return true
		l1712:
			position, tokenIndex, depth = position1712, tokenIndex1712, depth1712
			return false
		},
		/* 94 plx <- <(percent / pnLocalEsc)> */
		nil,
		/* 95 percent <- <('%' hex hex)> */
		nil,
		/* 96 hex <- <((&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))> */
		func() bool {
			position1724, tokenIndex1724, depth1724 := position, tokenIndex, depth
			{
				position1725 := position
				depth++
				{
					switch buffer[position] {
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l1724
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1724
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1724
						}
						position++
						break
					}
				}

				depth--
				add(rulehex, position1725)
			}
			return true
		l1724:
			position, tokenIndex, depth = position1724, tokenIndex1724, depth1724
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
			position1734, tokenIndex1734, depth1734 := position, tokenIndex, depth
			{
				position1735 := position
				depth++
				{
					position1736, tokenIndex1736, depth1736 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1737
					}
					position++
					goto l1736
				l1737:
					position, tokenIndex, depth = position1736, tokenIndex1736, depth1736
					if buffer[position] != rune('D') {
						goto l1734
					}
					position++
				}
			l1736:
				{
					position1738, tokenIndex1738, depth1738 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1739
					}
					position++
					goto l1738
				l1739:
					position, tokenIndex, depth = position1738, tokenIndex1738, depth1738
					if buffer[position] != rune('I') {
						goto l1734
					}
					position++
				}
			l1738:
				{
					position1740, tokenIndex1740, depth1740 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1741
					}
					position++
					goto l1740
				l1741:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
					if buffer[position] != rune('S') {
						goto l1734
					}
					position++
				}
			l1740:
				{
					position1742, tokenIndex1742, depth1742 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1743
					}
					position++
					goto l1742
				l1743:
					position, tokenIndex, depth = position1742, tokenIndex1742, depth1742
					if buffer[position] != rune('T') {
						goto l1734
					}
					position++
				}
			l1742:
				{
					position1744, tokenIndex1744, depth1744 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1745
					}
					position++
					goto l1744
				l1745:
					position, tokenIndex, depth = position1744, tokenIndex1744, depth1744
					if buffer[position] != rune('I') {
						goto l1734
					}
					position++
				}
			l1744:
				{
					position1746, tokenIndex1746, depth1746 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1747
					}
					position++
					goto l1746
				l1747:
					position, tokenIndex, depth = position1746, tokenIndex1746, depth1746
					if buffer[position] != rune('N') {
						goto l1734
					}
					position++
				}
			l1746:
				{
					position1748, tokenIndex1748, depth1748 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1749
					}
					position++
					goto l1748
				l1749:
					position, tokenIndex, depth = position1748, tokenIndex1748, depth1748
					if buffer[position] != rune('C') {
						goto l1734
					}
					position++
				}
			l1748:
				{
					position1750, tokenIndex1750, depth1750 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1751
					}
					position++
					goto l1750
				l1751:
					position, tokenIndex, depth = position1750, tokenIndex1750, depth1750
					if buffer[position] != rune('T') {
						goto l1734
					}
					position++
				}
			l1750:
				if !rules[ruleskip]() {
					goto l1734
				}
				depth--
				add(ruleDISTINCT, position1735)
			}
			return true
		l1734:
			position, tokenIndex, depth = position1734, tokenIndex1734, depth1734
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
			position1755, tokenIndex1755, depth1755 := position, tokenIndex, depth
			{
				position1756 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1755
				}
				position++
				if !rules[ruleskip]() {
					goto l1755
				}
				depth--
				add(ruleLBRACE, position1756)
			}
			return true
		l1755:
			position, tokenIndex, depth = position1755, tokenIndex1755, depth1755
			return false
		},
		/* 109 RBRACE <- <('}' skip)> */
		func() bool {
			position1757, tokenIndex1757, depth1757 := position, tokenIndex, depth
			{
				position1758 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1757
				}
				position++
				if !rules[ruleskip]() {
					goto l1757
				}
				depth--
				add(ruleRBRACE, position1758)
			}
			return true
		l1757:
			position, tokenIndex, depth = position1757, tokenIndex1757, depth1757
			return false
		},
		/* 110 LBRACK <- <('[' skip)> */
		nil,
		/* 111 RBRACK <- <(']' skip)> */
		nil,
		/* 112 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1761, tokenIndex1761, depth1761 := position, tokenIndex, depth
			{
				position1762 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1761
				}
				position++
				if !rules[ruleskip]() {
					goto l1761
				}
				depth--
				add(ruleSEMICOLON, position1762)
			}
			return true
		l1761:
			position, tokenIndex, depth = position1761, tokenIndex1761, depth1761
			return false
		},
		/* 113 COMMA <- <(',' skip)> */
		func() bool {
			position1763, tokenIndex1763, depth1763 := position, tokenIndex, depth
			{
				position1764 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1763
				}
				position++
				if !rules[ruleskip]() {
					goto l1763
				}
				depth--
				add(ruleCOMMA, position1764)
			}
			return true
		l1763:
			position, tokenIndex, depth = position1763, tokenIndex1763, depth1763
			return false
		},
		/* 114 DOT <- <('.' skip)> */
		func() bool {
			position1765, tokenIndex1765, depth1765 := position, tokenIndex, depth
			{
				position1766 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1765
				}
				position++
				if !rules[ruleskip]() {
					goto l1765
				}
				depth--
				add(ruleDOT, position1766)
			}
			return true
		l1765:
			position, tokenIndex, depth = position1765, tokenIndex1765, depth1765
			return false
		},
		/* 115 COLON <- <(':' skip)> */
		nil,
		/* 116 PIPE <- <('|' skip)> */
		func() bool {
			position1768, tokenIndex1768, depth1768 := position, tokenIndex, depth
			{
				position1769 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1768
				}
				position++
				if !rules[ruleskip]() {
					goto l1768
				}
				depth--
				add(rulePIPE, position1769)
			}
			return true
		l1768:
			position, tokenIndex, depth = position1768, tokenIndex1768, depth1768
			return false
		},
		/* 117 SLASH <- <('/' skip)> */
		func() bool {
			position1770, tokenIndex1770, depth1770 := position, tokenIndex, depth
			{
				position1771 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1770
				}
				position++
				if !rules[ruleskip]() {
					goto l1770
				}
				depth--
				add(ruleSLASH, position1771)
			}
			return true
		l1770:
			position, tokenIndex, depth = position1770, tokenIndex1770, depth1770
			return false
		},
		/* 118 INVERSE <- <('^' skip)> */
		func() bool {
			position1772, tokenIndex1772, depth1772 := position, tokenIndex, depth
			{
				position1773 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1772
				}
				position++
				if !rules[ruleskip]() {
					goto l1772
				}
				depth--
				add(ruleINVERSE, position1773)
			}
			return true
		l1772:
			position, tokenIndex, depth = position1772, tokenIndex1772, depth1772
			return false
		},
		/* 119 LPAREN <- <('(' skip)> */
		func() bool {
			position1774, tokenIndex1774, depth1774 := position, tokenIndex, depth
			{
				position1775 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1774
				}
				position++
				if !rules[ruleskip]() {
					goto l1774
				}
				depth--
				add(ruleLPAREN, position1775)
			}
			return true
		l1774:
			position, tokenIndex, depth = position1774, tokenIndex1774, depth1774
			return false
		},
		/* 120 RPAREN <- <(')' skip)> */
		func() bool {
			position1776, tokenIndex1776, depth1776 := position, tokenIndex, depth
			{
				position1777 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1776
				}
				position++
				if !rules[ruleskip]() {
					goto l1776
				}
				depth--
				add(ruleRPAREN, position1777)
			}
			return true
		l1776:
			position, tokenIndex, depth = position1776, tokenIndex1776, depth1776
			return false
		},
		/* 121 ISA <- <('a' skip)> */
		func() bool {
			position1778, tokenIndex1778, depth1778 := position, tokenIndex, depth
			{
				position1779 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1778
				}
				position++
				if !rules[ruleskip]() {
					goto l1778
				}
				depth--
				add(ruleISA, position1779)
			}
			return true
		l1778:
			position, tokenIndex, depth = position1778, tokenIndex1778, depth1778
			return false
		},
		/* 122 NOT <- <('!' skip)> */
		func() bool {
			position1780, tokenIndex1780, depth1780 := position, tokenIndex, depth
			{
				position1781 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1780
				}
				position++
				if !rules[ruleskip]() {
					goto l1780
				}
				depth--
				add(ruleNOT, position1781)
			}
			return true
		l1780:
			position, tokenIndex, depth = position1780, tokenIndex1780, depth1780
			return false
		},
		/* 123 STAR <- <('*' skip)> */
		func() bool {
			position1782, tokenIndex1782, depth1782 := position, tokenIndex, depth
			{
				position1783 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1782
				}
				position++
				if !rules[ruleskip]() {
					goto l1782
				}
				depth--
				add(ruleSTAR, position1783)
			}
			return true
		l1782:
			position, tokenIndex, depth = position1782, tokenIndex1782, depth1782
			return false
		},
		/* 124 QUESTION <- <('?' skip)> */
		nil,
		/* 125 PLUS <- <('+' skip)> */
		func() bool {
			position1785, tokenIndex1785, depth1785 := position, tokenIndex, depth
			{
				position1786 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1785
				}
				position++
				if !rules[ruleskip]() {
					goto l1785
				}
				depth--
				add(rulePLUS, position1786)
			}
			return true
		l1785:
			position, tokenIndex, depth = position1785, tokenIndex1785, depth1785
			return false
		},
		/* 126 MINUS <- <('-' skip)> */
		func() bool {
			position1787, tokenIndex1787, depth1787 := position, tokenIndex, depth
			{
				position1788 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1787
				}
				position++
				if !rules[ruleskip]() {
					goto l1787
				}
				depth--
				add(ruleMINUS, position1788)
			}
			return true
		l1787:
			position, tokenIndex, depth = position1787, tokenIndex1787, depth1787
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
			position1793, tokenIndex1793, depth1793 := position, tokenIndex, depth
			{
				position1794 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1793
				}
				position++
			l1795:
				{
					position1796, tokenIndex1796, depth1796 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1796
					}
					position++
					goto l1795
				l1796:
					position, tokenIndex, depth = position1796, tokenIndex1796, depth1796
				}
				if !rules[ruleskip]() {
					goto l1793
				}
				depth--
				add(ruleINTEGER, position1794)
			}
			return true
		l1793:
			position, tokenIndex, depth = position1793, tokenIndex1793, depth1793
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
			position1802, tokenIndex1802, depth1802 := position, tokenIndex, depth
			{
				position1803 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1802
				}
				position++
				if !rules[ruleskip]() {
					goto l1802
				}
				depth--
				add(ruleEQ, position1803)
			}
			return true
		l1802:
			position, tokenIndex, depth = position1802, tokenIndex1802, depth1802
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
			position1811, tokenIndex1811, depth1811 := position, tokenIndex, depth
			{
				position1812 := position
				depth++
				{
					position1813, tokenIndex1813, depth1813 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1814
					}
					position++
					goto l1813
				l1814:
					position, tokenIndex, depth = position1813, tokenIndex1813, depth1813
					if buffer[position] != rune('A') {
						goto l1811
					}
					position++
				}
			l1813:
				{
					position1815, tokenIndex1815, depth1815 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1816
					}
					position++
					goto l1815
				l1816:
					position, tokenIndex, depth = position1815, tokenIndex1815, depth1815
					if buffer[position] != rune('S') {
						goto l1811
					}
					position++
				}
			l1815:
				if !rules[ruleskip]() {
					goto l1811
				}
				depth--
				add(ruleAS, position1812)
			}
			return true
		l1811:
			position, tokenIndex, depth = position1811, tokenIndex1811, depth1811
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
			position1885, tokenIndex1885, depth1885 := position, tokenIndex, depth
			{
				position1886 := position
				depth++
				{
					position1887, tokenIndex1887, depth1887 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1888
					}
					position++
					goto l1887
				l1888:
					position, tokenIndex, depth = position1887, tokenIndex1887, depth1887
					if buffer[position] != rune('B') {
						goto l1885
					}
					position++
				}
			l1887:
				{
					position1889, tokenIndex1889, depth1889 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1890
					}
					position++
					goto l1889
				l1890:
					position, tokenIndex, depth = position1889, tokenIndex1889, depth1889
					if buffer[position] != rune('Y') {
						goto l1885
					}
					position++
				}
			l1889:
				if !rules[ruleskip]() {
					goto l1885
				}
				depth--
				add(ruleBY, position1886)
			}
			return true
		l1885:
			position, tokenIndex, depth = position1885, tokenIndex1885, depth1885
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
				position1895 := position
				depth++
			l1896:
				{
					position1897, tokenIndex1897, depth1897 := position, tokenIndex, depth
					{
						position1898, tokenIndex1898, depth1898 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1899
						}
						goto l1898
					l1899:
						position, tokenIndex, depth = position1898, tokenIndex1898, depth1898
						{
							position1900 := position
							depth++
							{
								position1901 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1897
								}
								position++
							l1902:
								{
									position1903, tokenIndex1903, depth1903 := position, tokenIndex, depth
									{
										position1904, tokenIndex1904, depth1904 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1904
										}
										goto l1903
									l1904:
										position, tokenIndex, depth = position1904, tokenIndex1904, depth1904
									}
									if !matchDot() {
										goto l1903
									}
									goto l1902
								l1903:
									position, tokenIndex, depth = position1903, tokenIndex1903, depth1903
								}
								if !rules[ruleendOfLine]() {
									goto l1897
								}
								depth--
								add(rulePegText, position1901)
							}
							{
								add(ruleAction13, position)
							}
							depth--
							add(rulecomment, position1900)
						}
					}
				l1898:
					goto l1896
				l1897:
					position, tokenIndex, depth = position1897, tokenIndex1897, depth1897
				}
				depth--
				add(ruleskip, position1895)
			}
			return true
		},
		/* 219 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1906, tokenIndex1906, depth1906 := position, tokenIndex, depth
			{
				position1907 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1906
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1906
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1906
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1906
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1906
						}
						break
					}
				}

				depth--
				add(rulews, position1907)
			}
			return true
		l1906:
			position, tokenIndex, depth = position1906, tokenIndex1906, depth1906
			return false
		},
		/* 220 comment <- <(<('#' (!endOfLine .)* endOfLine)> Action13)> */
		nil,
		/* 221 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1910, tokenIndex1910, depth1910 := position, tokenIndex, depth
			{
				position1911 := position
				depth++
				{
					position1912, tokenIndex1912, depth1912 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1913
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1913
					}
					position++
					goto l1912
				l1913:
					position, tokenIndex, depth = position1912, tokenIndex1912, depth1912
					if buffer[position] != rune('\n') {
						goto l1914
					}
					position++
					goto l1912
				l1914:
					position, tokenIndex, depth = position1912, tokenIndex1912, depth1912
					if buffer[position] != rune('\r') {
						goto l1910
					}
					position++
				}
			l1912:
				depth--
				add(ruleendOfLine, position1911)
			}
			return true
		l1910:
			position, tokenIndex, depth = position1910, tokenIndex1910, depth1910
			return false
		},
		nil,
		/* 224 Action0 <- <{ p.addPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 225 Action1 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 226 Action2 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 227 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 228 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 229 Action5 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 230 Action6 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 231 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 232 Action8 <- <{ p.setObject(p.skipComment(buffer, begin, end)); p.addTriplePattern() }> */
		nil,
		/* 233 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 234 Action10 <- <{ p.setPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 235 Action11 <- <{ p.setPathLength(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 236 Action12 <- <{ p.setKeyword(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 237 Action13 <- <{ p.commentBegin = begin }> */
		nil,
	}
	p.rules = rules
}
