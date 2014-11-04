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
	ruleSERVICE
	ruleSILENT
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
	"SERVICE",
	"SILENT",
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

	// The start offset of a comment
	skipBegin int
	triplePattern

	*Scope

	Buffer string
	buffer []rune
	rules  [241]func() bool
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
			p.addPrefix(p.skipped(buffer, begin, end))
		case ruleAction1:
			p.S = p.skipped(buffer, begin, end)
		case ruleAction2:
			p.S = p.skipped(buffer, begin, end)
		case ruleAction3:
			p.S = "?POF"
		case ruleAction4:
			p.P = "?POF"
		case ruleAction5:
			p.P = p.skipped(buffer, begin, end)
		case ruleAction6:
			p.P = p.skipped(buffer, begin, end)
		case ruleAction7:
			p.O = "?POF"
			p.addTriplePattern()
		case ruleAction8:
			p.O = p.skipped(buffer, begin, end)
			p.addTriplePattern()
		case ruleAction9:
			p.O = "?FillVar"
			p.addTriplePattern()
		case ruleAction10:
			p.setPrefix(p.skipped(buffer, begin, end))
		case ruleAction11:
			p.setPathLength(p.skipped(buffer, begin, end))
		case ruleAction12:
			p.setKeyword(p.skipped(buffer, begin, end))
		case ruleAction13:
			p.skipBegin = begin

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
							position227, tokenIndex227, depth227 := position, tokenIndex, depth
							{
								position229 := position
								depth++
								{
									position230 := position
									depth++
									{
										position231, tokenIndex231, depth231 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l232
										}
										position++
										goto l231
									l232:
										position, tokenIndex, depth = position231, tokenIndex231, depth231
										if buffer[position] != rune('O') {
											goto l228
										}
										position++
									}
								l231:
									{
										position233, tokenIndex233, depth233 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l234
										}
										position++
										goto l233
									l234:
										position, tokenIndex, depth = position233, tokenIndex233, depth233
										if buffer[position] != rune('P') {
											goto l228
										}
										position++
									}
								l233:
									{
										position235, tokenIndex235, depth235 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l236
										}
										position++
										goto l235
									l236:
										position, tokenIndex, depth = position235, tokenIndex235, depth235
										if buffer[position] != rune('T') {
											goto l228
										}
										position++
									}
								l235:
									{
										position237, tokenIndex237, depth237 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l238
										}
										position++
										goto l237
									l238:
										position, tokenIndex, depth = position237, tokenIndex237, depth237
										if buffer[position] != rune('I') {
											goto l228
										}
										position++
									}
								l237:
									{
										position239, tokenIndex239, depth239 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l240
										}
										position++
										goto l239
									l240:
										position, tokenIndex, depth = position239, tokenIndex239, depth239
										if buffer[position] != rune('O') {
											goto l228
										}
										position++
									}
								l239:
									{
										position241, tokenIndex241, depth241 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l242
										}
										position++
										goto l241
									l242:
										position, tokenIndex, depth = position241, tokenIndex241, depth241
										if buffer[position] != rune('N') {
											goto l228
										}
										position++
									}
								l241:
									{
										position243, tokenIndex243, depth243 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l244
										}
										position++
										goto l243
									l244:
										position, tokenIndex, depth = position243, tokenIndex243, depth243
										if buffer[position] != rune('A') {
											goto l228
										}
										position++
									}
								l243:
									{
										position245, tokenIndex245, depth245 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l246
										}
										position++
										goto l245
									l246:
										position, tokenIndex, depth = position245, tokenIndex245, depth245
										if buffer[position] != rune('L') {
											goto l228
										}
										position++
									}
								l245:
									if !rules[ruleskip]() {
										goto l228
									}
									depth--
									add(ruleOPTIONAL, position230)
								}
								if !rules[ruleLBRACE]() {
									goto l228
								}
								{
									position247, tokenIndex247, depth247 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l248
									}
									goto l247
								l248:
									position, tokenIndex, depth = position247, tokenIndex247, depth247
									if !rules[rulegraphPattern]() {
										goto l228
									}
								}
							l247:
								if !rules[ruleRBRACE]() {
									goto l228
								}
								depth--
								add(ruleoptionalGraphPattern, position229)
							}
							goto l227
						l228:
							position, tokenIndex, depth = position227, tokenIndex227, depth227
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l249
							}
							goto l227
						l249:
							position, tokenIndex, depth = position227, tokenIndex227, depth227
							{
								position251 := position
								depth++
								{
									position252 := position
									depth++
									{
										position253, tokenIndex253, depth253 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l254
										}
										position++
										goto l253
									l254:
										position, tokenIndex, depth = position253, tokenIndex253, depth253
										if buffer[position] != rune('G') {
											goto l250
										}
										position++
									}
								l253:
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l256
										}
										position++
										goto l255
									l256:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
										if buffer[position] != rune('R') {
											goto l250
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
											goto l250
										}
										position++
									}
								l257:
									{
										position259, tokenIndex259, depth259 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l260
										}
										position++
										goto l259
									l260:
										position, tokenIndex, depth = position259, tokenIndex259, depth259
										if buffer[position] != rune('P') {
											goto l250
										}
										position++
									}
								l259:
									{
										position261, tokenIndex261, depth261 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l262
										}
										position++
										goto l261
									l262:
										position, tokenIndex, depth = position261, tokenIndex261, depth261
										if buffer[position] != rune('H') {
											goto l250
										}
										position++
									}
								l261:
									if !rules[ruleskip]() {
										goto l250
									}
									depth--
									add(ruleGRAPH, position252)
								}
								{
									position263, tokenIndex263, depth263 := position, tokenIndex, depth
									if !rules[rulevar]() {
										goto l264
									}
									goto l263
								l264:
									position, tokenIndex, depth = position263, tokenIndex263, depth263
									if !rules[ruleiriref]() {
										goto l250
									}
								}
							l263:
								if !rules[rulegroupGraphPattern]() {
									goto l250
								}
								depth--
								add(rulegraphGraphPattern, position251)
							}
							goto l227
						l250:
							position, tokenIndex, depth = position227, tokenIndex227, depth227
							{
								position266 := position
								depth++
								{
									position267 := position
									depth++
									{
										position268, tokenIndex268, depth268 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l269
										}
										position++
										goto l268
									l269:
										position, tokenIndex, depth = position268, tokenIndex268, depth268
										if buffer[position] != rune('M') {
											goto l265
										}
										position++
									}
								l268:
									{
										position270, tokenIndex270, depth270 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l271
										}
										position++
										goto l270
									l271:
										position, tokenIndex, depth = position270, tokenIndex270, depth270
										if buffer[position] != rune('I') {
											goto l265
										}
										position++
									}
								l270:
									{
										position272, tokenIndex272, depth272 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l273
										}
										position++
										goto l272
									l273:
										position, tokenIndex, depth = position272, tokenIndex272, depth272
										if buffer[position] != rune('N') {
											goto l265
										}
										position++
									}
								l272:
									{
										position274, tokenIndex274, depth274 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l275
										}
										position++
										goto l274
									l275:
										position, tokenIndex, depth = position274, tokenIndex274, depth274
										if buffer[position] != rune('U') {
											goto l265
										}
										position++
									}
								l274:
									{
										position276, tokenIndex276, depth276 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l277
										}
										position++
										goto l276
									l277:
										position, tokenIndex, depth = position276, tokenIndex276, depth276
										if buffer[position] != rune('S') {
											goto l265
										}
										position++
									}
								l276:
									if !rules[ruleskip]() {
										goto l265
									}
									depth--
									add(ruleMINUSSETOPER, position267)
								}
								if !rules[rulegroupGraphPattern]() {
									goto l265
								}
								depth--
								add(ruleminusGraphPattern, position266)
							}
							goto l227
						l265:
							position, tokenIndex, depth = position227, tokenIndex227, depth227
							{
								position278 := position
								depth++
								{
									position279 := position
									depth++
									depth--
									add(ruleSERVICE, position279)
								}
								{
									position280, tokenIndex280, depth280 := position, tokenIndex, depth
									{
										position282 := position
										depth++
										depth--
										add(ruleSILENT, position282)
									}
									goto l281

									position, tokenIndex, depth = position280, tokenIndex280, depth280
								}
							l281:
								{
									position283, tokenIndex283, depth283 := position, tokenIndex, depth
									if !rules[rulevar]() {
										goto l284
									}
									goto l283
								l284:
									position, tokenIndex, depth = position283, tokenIndex283, depth283
									if !rules[ruleiriref]() {
										goto l224
									}
								}
							l283:
								if !rules[rulegroupGraphPattern]() {
									goto l224
								}
								depth--
								add(ruleserviceGraphPattern, position278)
							}
						}
					l227:
						depth--
						add(rulegraphPatternNotTriples, position226)
					}
					{
						position285, tokenIndex285, depth285 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l285
						}
						goto l286
					l285:
						position, tokenIndex, depth = position285, tokenIndex285, depth285
					}
				l286:
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
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern / graphGraphPattern / minusGraphPattern / serviceGraphPattern)> */
		nil,
		/* 19 serviceGraphPattern <- <(SERVICE SILENT? (var / iriref) groupGraphPattern)> */
		nil,
		/* 20 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 21 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position290, tokenIndex290, depth290 := position, tokenIndex, depth
			{
				position291 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l290
				}
				{
					position292, tokenIndex292, depth292 := position, tokenIndex, depth
					{
						position294 := position
						depth++
						{
							position295, tokenIndex295, depth295 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l296
							}
							position++
							goto l295
						l296:
							position, tokenIndex, depth = position295, tokenIndex295, depth295
							if buffer[position] != rune('U') {
								goto l292
							}
							position++
						}
					l295:
						{
							position297, tokenIndex297, depth297 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l298
							}
							position++
							goto l297
						l298:
							position, tokenIndex, depth = position297, tokenIndex297, depth297
							if buffer[position] != rune('N') {
								goto l292
							}
							position++
						}
					l297:
						{
							position299, tokenIndex299, depth299 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l300
							}
							position++
							goto l299
						l300:
							position, tokenIndex, depth = position299, tokenIndex299, depth299
							if buffer[position] != rune('I') {
								goto l292
							}
							position++
						}
					l299:
						{
							position301, tokenIndex301, depth301 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l302
							}
							position++
							goto l301
						l302:
							position, tokenIndex, depth = position301, tokenIndex301, depth301
							if buffer[position] != rune('O') {
								goto l292
							}
							position++
						}
					l301:
						{
							position303, tokenIndex303, depth303 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l304
							}
							position++
							goto l303
						l304:
							position, tokenIndex, depth = position303, tokenIndex303, depth303
							if buffer[position] != rune('N') {
								goto l292
							}
							position++
						}
					l303:
						if !rules[ruleskip]() {
							goto l292
						}
						depth--
						add(ruleUNION, position294)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l292
					}
					goto l293
				l292:
					position, tokenIndex, depth = position292, tokenIndex292, depth292
				}
			l293:
				depth--
				add(rulegroupOrUnionGraphPattern, position291)
			}
			return true
		l290:
			position, tokenIndex, depth = position290, tokenIndex290, depth290
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
			position308, tokenIndex308, depth308 := position, tokenIndex, depth
			{
				position309 := position
				depth++
				{
					position310, tokenIndex310, depth310 := position, tokenIndex, depth
					{
						position312 := position
						depth++
						{
							position313, tokenIndex313, depth313 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l314
							}
							position++
							goto l313
						l314:
							position, tokenIndex, depth = position313, tokenIndex313, depth313
							if buffer[position] != rune('F') {
								goto l311
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
								goto l311
							}
							position++
						}
					l315:
						{
							position317, tokenIndex317, depth317 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l318
							}
							position++
							goto l317
						l318:
							position, tokenIndex, depth = position317, tokenIndex317, depth317
							if buffer[position] != rune('L') {
								goto l311
							}
							position++
						}
					l317:
						{
							position319, tokenIndex319, depth319 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l320
							}
							position++
							goto l319
						l320:
							position, tokenIndex, depth = position319, tokenIndex319, depth319
							if buffer[position] != rune('T') {
								goto l311
							}
							position++
						}
					l319:
						{
							position321, tokenIndex321, depth321 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l322
							}
							position++
							goto l321
						l322:
							position, tokenIndex, depth = position321, tokenIndex321, depth321
							if buffer[position] != rune('E') {
								goto l311
							}
							position++
						}
					l321:
						{
							position323, tokenIndex323, depth323 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l324
							}
							position++
							goto l323
						l324:
							position, tokenIndex, depth = position323, tokenIndex323, depth323
							if buffer[position] != rune('R') {
								goto l311
							}
							position++
						}
					l323:
						if !rules[ruleskip]() {
							goto l311
						}
						depth--
						add(ruleFILTER, position312)
					}
					if !rules[ruleconstraint]() {
						goto l311
					}
					goto l310
				l311:
					position, tokenIndex, depth = position310, tokenIndex310, depth310
					{
						position325 := position
						depth++
						{
							position326, tokenIndex326, depth326 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l327
							}
							position++
							goto l326
						l327:
							position, tokenIndex, depth = position326, tokenIndex326, depth326
							if buffer[position] != rune('B') {
								goto l308
							}
							position++
						}
					l326:
						{
							position328, tokenIndex328, depth328 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l329
							}
							position++
							goto l328
						l329:
							position, tokenIndex, depth = position328, tokenIndex328, depth328
							if buffer[position] != rune('I') {
								goto l308
							}
							position++
						}
					l328:
						{
							position330, tokenIndex330, depth330 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l331
							}
							position++
							goto l330
						l331:
							position, tokenIndex, depth = position330, tokenIndex330, depth330
							if buffer[position] != rune('N') {
								goto l308
							}
							position++
						}
					l330:
						{
							position332, tokenIndex332, depth332 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l333
							}
							position++
							goto l332
						l333:
							position, tokenIndex, depth = position332, tokenIndex332, depth332
							if buffer[position] != rune('D') {
								goto l308
							}
							position++
						}
					l332:
						if !rules[ruleskip]() {
							goto l308
						}
						depth--
						add(ruleBIND, position325)
					}
					if !rules[ruleLPAREN]() {
						goto l308
					}
					if !rules[ruleexpression]() {
						goto l308
					}
					if !rules[ruleAS]() {
						goto l308
					}
					if !rules[rulevar]() {
						goto l308
					}
					if !rules[ruleRPAREN]() {
						goto l308
					}
				}
			l310:
				depth--
				add(rulefilterOrBind, position309)
			}
			return true
		l308:
			position, tokenIndex, depth = position308, tokenIndex308, depth308
			return false
		},
		/* 26 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		func() bool {
			position334, tokenIndex334, depth334 := position, tokenIndex, depth
			{
				position335 := position
				depth++
				{
					position336, tokenIndex336, depth336 := position, tokenIndex, depth
					if !rules[rulebrackettedExpression]() {
						goto l337
					}
					goto l336
				l337:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if !rules[rulebuiltinCall]() {
						goto l338
					}
					goto l336
				l338:
					position, tokenIndex, depth = position336, tokenIndex336, depth336
					if !rules[rulefunctionCall]() {
						goto l334
					}
				}
			l336:
				depth--
				add(ruleconstraint, position335)
			}
			return true
		l334:
			position, tokenIndex, depth = position334, tokenIndex334, depth334
			return false
		},
		/* 27 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position339, tokenIndex339, depth339 := position, tokenIndex, depth
			{
				position340 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l339
				}
			l341:
				{
					position342, tokenIndex342, depth342 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l342
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l342
					}
					goto l341
				l342:
					position, tokenIndex, depth = position342, tokenIndex342, depth342
				}
				{
					position343, tokenIndex343, depth343 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l343
					}
					goto l344
				l343:
					position, tokenIndex, depth = position343, tokenIndex343, depth343
				}
			l344:
				depth--
				add(ruletriplesBlock, position340)
			}
			return true
		l339:
			position, tokenIndex, depth = position339, tokenIndex339, depth339
			return false
		},
		/* 28 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath?))> */
		func() bool {
			position345, tokenIndex345, depth345 := position, tokenIndex, depth
			{
				position346 := position
				depth++
				{
					position347, tokenIndex347, depth347 := position, tokenIndex, depth
					{
						position349 := position
						depth++
						{
							position350, tokenIndex350, depth350 := position, tokenIndex, depth
							{
								position352 := position
								depth++
								if !rules[rulevar]() {
									goto l351
								}
								depth--
								add(rulePegText, position352)
							}
							{
								add(ruleAction1, position)
							}
							goto l350
						l351:
							position, tokenIndex, depth = position350, tokenIndex350, depth350
							{
								position355 := position
								depth++
								if !rules[rulegraphTerm]() {
									goto l354
								}
								depth--
								add(rulePegText, position355)
							}
							{
								add(ruleAction2, position)
							}
							goto l350
						l354:
							position, tokenIndex, depth = position350, tokenIndex350, depth350
							if !rules[rulepof]() {
								goto l348
							}
							{
								add(ruleAction3, position)
							}
						}
					l350:
						depth--
						add(rulevarOrTerm, position349)
					}
					if !rules[rulepropertyListPath]() {
						goto l348
					}
					goto l347
				l348:
					position, tokenIndex, depth = position347, tokenIndex347, depth347
					if !rules[ruletriplesNodePath]() {
						goto l345
					}
					{
						position358, tokenIndex358, depth358 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l358
						}
						goto l359
					l358:
						position, tokenIndex, depth = position358, tokenIndex358, depth358
					}
				l359:
				}
			l347:
				depth--
				add(ruletriplesSameSubjectPath, position346)
			}
			return true
		l345:
			position, tokenIndex, depth = position345, tokenIndex345, depth345
			return false
		},
		/* 29 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 30 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position361, tokenIndex361, depth361 := position, tokenIndex, depth
			{
				position362 := position
				depth++
				{
					position363, tokenIndex363, depth363 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l364
					}
					goto l363
				l364:
					position, tokenIndex, depth = position363, tokenIndex363, depth363
					{
						switch buffer[position] {
						case '(':
							if !rules[rulenil]() {
								goto l361
							}
							break
						case '[', '_':
							{
								position366 := position
								depth++
								{
									position367, tokenIndex367, depth367 := position, tokenIndex, depth
									{
										position369 := position
										depth++
										if buffer[position] != rune('_') {
											goto l368
										}
										position++
										if buffer[position] != rune(':') {
											goto l368
										}
										position++
										{
											position370, tokenIndex370, depth370 := position, tokenIndex, depth
											if !rules[rulepnCharsU]() {
												goto l371
											}
											goto l370
										l371:
											position, tokenIndex, depth = position370, tokenIndex370, depth370
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l368
											}
											position++
										}
									l370:
										{
											position372, tokenIndex372, depth372 := position, tokenIndex, depth
											{
												position374, tokenIndex374, depth374 := position, tokenIndex, depth
											l376:
												{
													position377, tokenIndex377, depth377 := position, tokenIndex, depth
													{
														position378, tokenIndex378, depth378 := position, tokenIndex, depth
														if !rules[rulepnCharsU]() {
															goto l379
														}
														goto l378
													l379:
														position, tokenIndex, depth = position378, tokenIndex378, depth378
														{
															switch buffer[position] {
															case '.':
																if buffer[position] != rune('.') {
																	goto l377
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l377
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l377
																}
																position++
																break
															}
														}

													}
												l378:
													goto l376
												l377:
													position, tokenIndex, depth = position377, tokenIndex377, depth377
												}
												if !rules[rulepnCharsU]() {
													goto l375
												}
												goto l374
											l375:
												position, tokenIndex, depth = position374, tokenIndex374, depth374
												{
													position381, tokenIndex381, depth381 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l382
													}
													position++
													goto l381
												l382:
													position, tokenIndex, depth = position381, tokenIndex381, depth381
													if buffer[position] != rune('-') {
														goto l372
													}
													position++
												}
											l381:
											}
										l374:
											goto l373
										l372:
											position, tokenIndex, depth = position372, tokenIndex372, depth372
										}
									l373:
										if !rules[ruleskip]() {
											goto l368
										}
										depth--
										add(ruleblankNodeLabel, position369)
									}
									goto l367
								l368:
									position, tokenIndex, depth = position367, tokenIndex367, depth367
									{
										position383 := position
										depth++
										if buffer[position] != rune('[') {
											goto l361
										}
										position++
									l384:
										{
											position385, tokenIndex385, depth385 := position, tokenIndex, depth
											if !rules[rulews]() {
												goto l385
											}
											goto l384
										l385:
											position, tokenIndex, depth = position385, tokenIndex385, depth385
										}
										if buffer[position] != rune(']') {
											goto l361
										}
										position++
										if !rules[ruleskip]() {
											goto l361
										}
										depth--
										add(ruleanon, position383)
									}
								}
							l367:
								depth--
								add(ruleblankNode, position366)
							}
							break
						case 'F', 'T', 'f', 't':
							if !rules[rulebooleanLiteral]() {
								goto l361
							}
							break
						case '"', '\'':
							if !rules[ruleliteral]() {
								goto l361
							}
							break
						default:
							if !rules[rulenumericLiteral]() {
								goto l361
							}
							break
						}
					}

				}
			l363:
				depth--
				add(rulegraphTerm, position362)
			}
			return true
		l361:
			position, tokenIndex, depth = position361, tokenIndex361, depth361
			return false
		},
		/* 31 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		func() bool {
			position386, tokenIndex386, depth386 := position, tokenIndex, depth
			{
				position387 := position
				depth++
				{
					position388, tokenIndex388, depth388 := position, tokenIndex, depth
					{
						position390 := position
						depth++
						if !rules[ruleLPAREN]() {
							goto l389
						}
						if !rules[rulegraphNodePath]() {
							goto l389
						}
					l391:
						{
							position392, tokenIndex392, depth392 := position, tokenIndex, depth
							if !rules[rulegraphNodePath]() {
								goto l392
							}
							goto l391
						l392:
							position, tokenIndex, depth = position392, tokenIndex392, depth392
						}
						if !rules[ruleRPAREN]() {
							goto l389
						}
						depth--
						add(rulecollectionPath, position390)
					}
					goto l388
				l389:
					position, tokenIndex, depth = position388, tokenIndex388, depth388
					{
						position393 := position
						depth++
						{
							position394 := position
							depth++
							if buffer[position] != rune('[') {
								goto l386
							}
							position++
							if !rules[ruleskip]() {
								goto l386
							}
							depth--
							add(ruleLBRACK, position394)
						}
						if !rules[rulepropertyListPath]() {
							goto l386
						}
						{
							position395 := position
							depth++
							if buffer[position] != rune(']') {
								goto l386
							}
							position++
							if !rules[ruleskip]() {
								goto l386
							}
							depth--
							add(ruleRBRACK, position395)
						}
						depth--
						add(ruleblankNodePropertyListPath, position393)
					}
				}
			l388:
				depth--
				add(ruletriplesNodePath, position387)
			}
			return true
		l386:
			position, tokenIndex, depth = position386, tokenIndex386, depth386
			return false
		},
		/* 32 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 33 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 34 propertyListPath <- <(((pof Action4) / (<var> Action5) / verbPath) objectListPath (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position398, tokenIndex398, depth398 := position, tokenIndex, depth
			{
				position399 := position
				depth++
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l401
					}
					{
						add(ruleAction4, position)
					}
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					{
						position404 := position
						depth++
						if !rules[rulevar]() {
							goto l403
						}
						depth--
						add(rulePegText, position404)
					}
					{
						add(ruleAction5, position)
					}
					goto l400
				l403:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					{
						position406 := position
						depth++
						if !rules[rulepath]() {
							goto l398
						}
						depth--
						add(ruleverbPath, position406)
					}
				}
			l400:
				{
					position407 := position
					depth++
					if !rules[ruleobjectPath]() {
						goto l398
					}
				l408:
					{
						position409, tokenIndex409, depth409 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l409
						}
						if !rules[ruleobjectPath]() {
							goto l409
						}
						goto l408
					l409:
						position, tokenIndex, depth = position409, tokenIndex409, depth409
					}
					depth--
					add(ruleobjectListPath, position407)
				}
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l410
					}
					{
						position412, tokenIndex412, depth412 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l412
						}
						goto l413
					l412:
						position, tokenIndex, depth = position412, tokenIndex412, depth412
					}
				l413:
					goto l411
				l410:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
				}
			l411:
				depth--
				add(rulepropertyListPath, position399)
			}
			return true
		l398:
			position, tokenIndex, depth = position398, tokenIndex398, depth398
			return false
		},
		/* 35 verbPath <- <path> */
		nil,
		/* 36 path <- <pathAlternative> */
		func() bool {
			position415, tokenIndex415, depth415 := position, tokenIndex, depth
			{
				position416 := position
				depth++
				{
					position417 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l415
					}
				l418:
					{
						position419, tokenIndex419, depth419 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l419
						}
						if !rules[rulepathSequence]() {
							goto l419
						}
						goto l418
					l419:
						position, tokenIndex, depth = position419, tokenIndex419, depth419
					}
					depth--
					add(rulepathAlternative, position417)
				}
				depth--
				add(rulepath, position416)
			}
			return true
		l415:
			position, tokenIndex, depth = position415, tokenIndex415, depth415
			return false
		},
		/* 37 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 38 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position421, tokenIndex421, depth421 := position, tokenIndex, depth
			{
				position422 := position
				depth++
				{
					position423 := position
					depth++
					{
						position424 := position
						depth++
						{
							position425, tokenIndex425, depth425 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l425
							}
							goto l426
						l425:
							position, tokenIndex, depth = position425, tokenIndex425, depth425
						}
					l426:
						{
							position427 := position
							depth++
							{
								position428, tokenIndex428, depth428 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l429
								}
								goto l428
							l429:
								position, tokenIndex, depth = position428, tokenIndex428, depth428
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l421
										}
										if !rules[rulepath]() {
											goto l421
										}
										if !rules[ruleRPAREN]() {
											goto l421
										}
										break
									case '!':
										if !rules[ruleNOT]() {
											goto l421
										}
										{
											position431 := position
											depth++
											{
												position432, tokenIndex432, depth432 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l433
												}
												goto l432
											l433:
												position, tokenIndex, depth = position432, tokenIndex432, depth432
												if !rules[ruleLPAREN]() {
													goto l421
												}
												{
													position434, tokenIndex434, depth434 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l434
													}
												l436:
													{
														position437, tokenIndex437, depth437 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l437
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l437
														}
														goto l436
													l437:
														position, tokenIndex, depth = position437, tokenIndex437, depth437
													}
													goto l435
												l434:
													position, tokenIndex, depth = position434, tokenIndex434, depth434
												}
											l435:
												if !rules[ruleRPAREN]() {
													goto l421
												}
											}
										l432:
											depth--
											add(rulepathNegatedPropertySet, position431)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l421
										}
										break
									}
								}

							}
						l428:
							depth--
							add(rulepathPrimary, position427)
						}
						{
							position438, tokenIndex438, depth438 := position, tokenIndex, depth
							{
								position440 := position
								depth++
								{
									switch buffer[position] {
									case '+':
										if !rules[rulePLUS]() {
											goto l438
										}
										break
									case '?':
										{
											position442 := position
											depth++
											if buffer[position] != rune('?') {
												goto l438
											}
											position++
											if !rules[ruleskip]() {
												goto l438
											}
											depth--
											add(ruleQUESTION, position442)
										}
										break
									default:
										if !rules[ruleSTAR]() {
											goto l438
										}
										break
									}
								}

								{
									position443, tokenIndex443, depth443 := position, tokenIndex, depth
									if !matchDot() {
										goto l443
									}
									goto l438
								l443:
									position, tokenIndex, depth = position443, tokenIndex443, depth443
								}
								depth--
								add(rulepathMod, position440)
							}
							goto l439
						l438:
							position, tokenIndex, depth = position438, tokenIndex438, depth438
						}
					l439:
						depth--
						add(rulepathElt, position424)
					}
					depth--
					add(rulePegText, position423)
				}
				{
					add(ruleAction6, position)
				}
			l445:
				{
					position446, tokenIndex446, depth446 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l446
					}
					if !rules[rulepathSequence]() {
						goto l446
					}
					goto l445
				l446:
					position, tokenIndex, depth = position446, tokenIndex446, depth446
				}
				depth--
				add(rulepathSequence, position422)
			}
			return true
		l421:
			position, tokenIndex, depth = position421, tokenIndex421, depth421
			return false
		},
		/* 39 pathElt <- <(INVERSE? pathPrimary pathMod?)> */
		nil,
		/* 40 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 41 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 42 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position450, tokenIndex450, depth450 := position, tokenIndex, depth
			{
				position451 := position
				depth++
				{
					position452, tokenIndex452, depth452 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l453
					}
					goto l452
				l453:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if !rules[ruleISA]() {
						goto l454
					}
					goto l452
				l454:
					position, tokenIndex, depth = position452, tokenIndex452, depth452
					if !rules[ruleINVERSE]() {
						goto l450
					}
					{
						position455, tokenIndex455, depth455 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l456
						}
						goto l455
					l456:
						position, tokenIndex, depth = position455, tokenIndex455, depth455
						if !rules[ruleISA]() {
							goto l450
						}
					}
				l455:
				}
			l452:
				depth--
				add(rulepathOneInPropertySet, position451)
			}
			return true
		l450:
			position, tokenIndex, depth = position450, tokenIndex450, depth450
			return false
		},
		/* 43 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !.)> */
		nil,
		/* 44 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 45 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		func() bool {
			{
				position460 := position
				depth++
				{
					position461, tokenIndex461, depth461 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l462
					}
					{
						add(ruleAction7, position)
					}
					goto l461
				l462:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					{
						position465 := position
						depth++
						if !rules[rulegraphNodePath]() {
							goto l464
						}
						depth--
						add(rulePegText, position465)
					}
					{
						add(ruleAction8, position)
					}
					goto l461
				l464:
					position, tokenIndex, depth = position461, tokenIndex461, depth461
					{
						add(ruleAction9, position)
					}
				}
			l461:
				depth--
				add(ruleobjectPath, position460)
			}
			return true
		},
		/* 46 graphNodePath <- <(var / graphTerm / triplesNodePath)> */
		func() bool {
			position468, tokenIndex468, depth468 := position, tokenIndex, depth
			{
				position469 := position
				depth++
				{
					position470, tokenIndex470, depth470 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l471
					}
					goto l470
				l471:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if !rules[rulegraphTerm]() {
						goto l472
					}
					goto l470
				l472:
					position, tokenIndex, depth = position470, tokenIndex470, depth470
					if !rules[ruletriplesNodePath]() {
						goto l468
					}
				}
			l470:
				depth--
				add(rulegraphNodePath, position469)
			}
			return true
		l468:
			position, tokenIndex, depth = position468, tokenIndex468, depth468
			return false
		},
		/* 47 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position474 := position
				depth++
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					{
						position477, tokenIndex477, depth477 := position, tokenIndex, depth
						{
							position479 := position
							depth++
							{
								position480, tokenIndex480, depth480 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l481
								}
								position++
								goto l480
							l481:
								position, tokenIndex, depth = position480, tokenIndex480, depth480
								if buffer[position] != rune('O') {
									goto l478
								}
								position++
							}
						l480:
							{
								position482, tokenIndex482, depth482 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l483
								}
								position++
								goto l482
							l483:
								position, tokenIndex, depth = position482, tokenIndex482, depth482
								if buffer[position] != rune('R') {
									goto l478
								}
								position++
							}
						l482:
							{
								position484, tokenIndex484, depth484 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l485
								}
								position++
								goto l484
							l485:
								position, tokenIndex, depth = position484, tokenIndex484, depth484
								if buffer[position] != rune('D') {
									goto l478
								}
								position++
							}
						l484:
							{
								position486, tokenIndex486, depth486 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l487
								}
								position++
								goto l486
							l487:
								position, tokenIndex, depth = position486, tokenIndex486, depth486
								if buffer[position] != rune('E') {
									goto l478
								}
								position++
							}
						l486:
							{
								position488, tokenIndex488, depth488 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l489
								}
								position++
								goto l488
							l489:
								position, tokenIndex, depth = position488, tokenIndex488, depth488
								if buffer[position] != rune('R') {
									goto l478
								}
								position++
							}
						l488:
							if !rules[ruleskip]() {
								goto l478
							}
							depth--
							add(ruleORDER, position479)
						}
						if !rules[ruleBY]() {
							goto l478
						}
						{
							position492 := position
							depth++
							{
								position493, tokenIndex493, depth493 := position, tokenIndex, depth
								{
									position495, tokenIndex495, depth495 := position, tokenIndex, depth
									{
										position497, tokenIndex497, depth497 := position, tokenIndex, depth
										{
											position499 := position
											depth++
											{
												position500, tokenIndex500, depth500 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l501
												}
												position++
												goto l500
											l501:
												position, tokenIndex, depth = position500, tokenIndex500, depth500
												if buffer[position] != rune('A') {
													goto l498
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
													goto l498
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
													goto l498
												}
												position++
											}
										l504:
											if !rules[ruleskip]() {
												goto l498
											}
											depth--
											add(ruleASC, position499)
										}
										goto l497
									l498:
										position, tokenIndex, depth = position497, tokenIndex497, depth497
										{
											position506 := position
											depth++
											{
												position507, tokenIndex507, depth507 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l508
												}
												position++
												goto l507
											l508:
												position, tokenIndex, depth = position507, tokenIndex507, depth507
												if buffer[position] != rune('D') {
													goto l495
												}
												position++
											}
										l507:
											{
												position509, tokenIndex509, depth509 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l510
												}
												position++
												goto l509
											l510:
												position, tokenIndex, depth = position509, tokenIndex509, depth509
												if buffer[position] != rune('E') {
													goto l495
												}
												position++
											}
										l509:
											{
												position511, tokenIndex511, depth511 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l512
												}
												position++
												goto l511
											l512:
												position, tokenIndex, depth = position511, tokenIndex511, depth511
												if buffer[position] != rune('S') {
													goto l495
												}
												position++
											}
										l511:
											{
												position513, tokenIndex513, depth513 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l514
												}
												position++
												goto l513
											l514:
												position, tokenIndex, depth = position513, tokenIndex513, depth513
												if buffer[position] != rune('C') {
													goto l495
												}
												position++
											}
										l513:
											if !rules[ruleskip]() {
												goto l495
											}
											depth--
											add(ruleDESC, position506)
										}
									}
								l497:
									goto l496
								l495:
									position, tokenIndex, depth = position495, tokenIndex495, depth495
								}
							l496:
								if !rules[rulebrackettedExpression]() {
									goto l494
								}
								goto l493
							l494:
								position, tokenIndex, depth = position493, tokenIndex493, depth493
								if !rules[rulefunctionCall]() {
									goto l515
								}
								goto l493
							l515:
								position, tokenIndex, depth = position493, tokenIndex493, depth493
								if !rules[rulebuiltinCall]() {
									goto l516
								}
								goto l493
							l516:
								position, tokenIndex, depth = position493, tokenIndex493, depth493
								if !rules[rulevar]() {
									goto l478
								}
							}
						l493:
							depth--
							add(ruleorderCondition, position492)
						}
					l490:
						{
							position491, tokenIndex491, depth491 := position, tokenIndex, depth
							{
								position517 := position
								depth++
								{
									position518, tokenIndex518, depth518 := position, tokenIndex, depth
									{
										position520, tokenIndex520, depth520 := position, tokenIndex, depth
										{
											position522, tokenIndex522, depth522 := position, tokenIndex, depth
											{
												position524 := position
												depth++
												{
													position525, tokenIndex525, depth525 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l526
													}
													position++
													goto l525
												l526:
													position, tokenIndex, depth = position525, tokenIndex525, depth525
													if buffer[position] != rune('A') {
														goto l523
													}
													position++
												}
											l525:
												{
													position527, tokenIndex527, depth527 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l528
													}
													position++
													goto l527
												l528:
													position, tokenIndex, depth = position527, tokenIndex527, depth527
													if buffer[position] != rune('S') {
														goto l523
													}
													position++
												}
											l527:
												{
													position529, tokenIndex529, depth529 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l530
													}
													position++
													goto l529
												l530:
													position, tokenIndex, depth = position529, tokenIndex529, depth529
													if buffer[position] != rune('C') {
														goto l523
													}
													position++
												}
											l529:
												if !rules[ruleskip]() {
													goto l523
												}
												depth--
												add(ruleASC, position524)
											}
											goto l522
										l523:
											position, tokenIndex, depth = position522, tokenIndex522, depth522
											{
												position531 := position
												depth++
												{
													position532, tokenIndex532, depth532 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l533
													}
													position++
													goto l532
												l533:
													position, tokenIndex, depth = position532, tokenIndex532, depth532
													if buffer[position] != rune('D') {
														goto l520
													}
													position++
												}
											l532:
												{
													position534, tokenIndex534, depth534 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l535
													}
													position++
													goto l534
												l535:
													position, tokenIndex, depth = position534, tokenIndex534, depth534
													if buffer[position] != rune('E') {
														goto l520
													}
													position++
												}
											l534:
												{
													position536, tokenIndex536, depth536 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l537
													}
													position++
													goto l536
												l537:
													position, tokenIndex, depth = position536, tokenIndex536, depth536
													if buffer[position] != rune('S') {
														goto l520
													}
													position++
												}
											l536:
												{
													position538, tokenIndex538, depth538 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l539
													}
													position++
													goto l538
												l539:
													position, tokenIndex, depth = position538, tokenIndex538, depth538
													if buffer[position] != rune('C') {
														goto l520
													}
													position++
												}
											l538:
												if !rules[ruleskip]() {
													goto l520
												}
												depth--
												add(ruleDESC, position531)
											}
										}
									l522:
										goto l521
									l520:
										position, tokenIndex, depth = position520, tokenIndex520, depth520
									}
								l521:
									if !rules[rulebrackettedExpression]() {
										goto l519
									}
									goto l518
								l519:
									position, tokenIndex, depth = position518, tokenIndex518, depth518
									if !rules[rulefunctionCall]() {
										goto l540
									}
									goto l518
								l540:
									position, tokenIndex, depth = position518, tokenIndex518, depth518
									if !rules[rulebuiltinCall]() {
										goto l541
									}
									goto l518
								l541:
									position, tokenIndex, depth = position518, tokenIndex518, depth518
									if !rules[rulevar]() {
										goto l491
									}
								}
							l518:
								depth--
								add(ruleorderCondition, position517)
							}
							goto l490
						l491:
							position, tokenIndex, depth = position491, tokenIndex491, depth491
						}
						goto l477
					l478:
						position, tokenIndex, depth = position477, tokenIndex477, depth477
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position543 := position
									depth++
									{
										position544, tokenIndex544, depth544 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l545
										}
										position++
										goto l544
									l545:
										position, tokenIndex, depth = position544, tokenIndex544, depth544
										if buffer[position] != rune('H') {
											goto l475
										}
										position++
									}
								l544:
									{
										position546, tokenIndex546, depth546 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l547
										}
										position++
										goto l546
									l547:
										position, tokenIndex, depth = position546, tokenIndex546, depth546
										if buffer[position] != rune('A') {
											goto l475
										}
										position++
									}
								l546:
									{
										position548, tokenIndex548, depth548 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l549
										}
										position++
										goto l548
									l549:
										position, tokenIndex, depth = position548, tokenIndex548, depth548
										if buffer[position] != rune('V') {
											goto l475
										}
										position++
									}
								l548:
									{
										position550, tokenIndex550, depth550 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l551
										}
										position++
										goto l550
									l551:
										position, tokenIndex, depth = position550, tokenIndex550, depth550
										if buffer[position] != rune('I') {
											goto l475
										}
										position++
									}
								l550:
									{
										position552, tokenIndex552, depth552 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l553
										}
										position++
										goto l552
									l553:
										position, tokenIndex, depth = position552, tokenIndex552, depth552
										if buffer[position] != rune('N') {
											goto l475
										}
										position++
									}
								l552:
									{
										position554, tokenIndex554, depth554 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l555
										}
										position++
										goto l554
									l555:
										position, tokenIndex, depth = position554, tokenIndex554, depth554
										if buffer[position] != rune('G') {
											goto l475
										}
										position++
									}
								l554:
									if !rules[ruleskip]() {
										goto l475
									}
									depth--
									add(ruleHAVING, position543)
								}
								if !rules[ruleconstraint]() {
									goto l475
								}
								break
							case 'G', 'g':
								{
									position556 := position
									depth++
									{
										position557, tokenIndex557, depth557 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l558
										}
										position++
										goto l557
									l558:
										position, tokenIndex, depth = position557, tokenIndex557, depth557
										if buffer[position] != rune('G') {
											goto l475
										}
										position++
									}
								l557:
									{
										position559, tokenIndex559, depth559 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l560
										}
										position++
										goto l559
									l560:
										position, tokenIndex, depth = position559, tokenIndex559, depth559
										if buffer[position] != rune('R') {
											goto l475
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
											goto l475
										}
										position++
									}
								l561:
									{
										position563, tokenIndex563, depth563 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l564
										}
										position++
										goto l563
									l564:
										position, tokenIndex, depth = position563, tokenIndex563, depth563
										if buffer[position] != rune('U') {
											goto l475
										}
										position++
									}
								l563:
									{
										position565, tokenIndex565, depth565 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l566
										}
										position++
										goto l565
									l566:
										position, tokenIndex, depth = position565, tokenIndex565, depth565
										if buffer[position] != rune('P') {
											goto l475
										}
										position++
									}
								l565:
									if !rules[ruleskip]() {
										goto l475
									}
									depth--
									add(ruleGROUP, position556)
								}
								if !rules[ruleBY]() {
									goto l475
								}
								{
									position569 := position
									depth++
									{
										position570, tokenIndex570, depth570 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l571
										}
										goto l570
									l571:
										position, tokenIndex, depth = position570, tokenIndex570, depth570
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l475
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l475
												}
												if !rules[ruleexpression]() {
													goto l475
												}
												{
													position573, tokenIndex573, depth573 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l573
													}
													if !rules[rulevar]() {
														goto l573
													}
													goto l574
												l573:
													position, tokenIndex, depth = position573, tokenIndex573, depth573
												}
											l574:
												if !rules[ruleRPAREN]() {
													goto l475
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l475
												}
												break
											}
										}

									}
								l570:
									depth--
									add(rulegroupCondition, position569)
								}
							l567:
								{
									position568, tokenIndex568, depth568 := position, tokenIndex, depth
									{
										position575 := position
										depth++
										{
											position576, tokenIndex576, depth576 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l577
											}
											goto l576
										l577:
											position, tokenIndex, depth = position576, tokenIndex576, depth576
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l568
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l568
													}
													if !rules[ruleexpression]() {
														goto l568
													}
													{
														position579, tokenIndex579, depth579 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l579
														}
														if !rules[rulevar]() {
															goto l579
														}
														goto l580
													l579:
														position, tokenIndex, depth = position579, tokenIndex579, depth579
													}
												l580:
													if !rules[ruleRPAREN]() {
														goto l568
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l568
													}
													break
												}
											}

										}
									l576:
										depth--
										add(rulegroupCondition, position575)
									}
									goto l567
								l568:
									position, tokenIndex, depth = position568, tokenIndex568, depth568
								}
								break
							default:
								{
									position581 := position
									depth++
									{
										position582, tokenIndex582, depth582 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l583
										}
										{
											position584, tokenIndex584, depth584 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l584
											}
											goto l585
										l584:
											position, tokenIndex, depth = position584, tokenIndex584, depth584
										}
									l585:
										goto l582
									l583:
										position, tokenIndex, depth = position582, tokenIndex582, depth582
										if !rules[ruleoffset]() {
											goto l475
										}
										{
											position586, tokenIndex586, depth586 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l586
											}
											goto l587
										l586:
											position, tokenIndex, depth = position586, tokenIndex586, depth586
										}
									l587:
									}
								l582:
									depth--
									add(rulelimitOffsetClauses, position581)
								}
								break
							}
						}

					}
				l477:
					goto l476
				l475:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
				}
			l476:
				depth--
				add(rulesolutionModifier, position474)
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
			position591, tokenIndex591, depth591 := position, tokenIndex, depth
			{
				position592 := position
				depth++
				{
					position593 := position
					depth++
					{
						position594, tokenIndex594, depth594 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l595
						}
						position++
						goto l594
					l595:
						position, tokenIndex, depth = position594, tokenIndex594, depth594
						if buffer[position] != rune('L') {
							goto l591
						}
						position++
					}
				l594:
					{
						position596, tokenIndex596, depth596 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l597
						}
						position++
						goto l596
					l597:
						position, tokenIndex, depth = position596, tokenIndex596, depth596
						if buffer[position] != rune('I') {
							goto l591
						}
						position++
					}
				l596:
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l599
						}
						position++
						goto l598
					l599:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
						if buffer[position] != rune('M') {
							goto l591
						}
						position++
					}
				l598:
					{
						position600, tokenIndex600, depth600 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l601
						}
						position++
						goto l600
					l601:
						position, tokenIndex, depth = position600, tokenIndex600, depth600
						if buffer[position] != rune('I') {
							goto l591
						}
						position++
					}
				l600:
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l603
						}
						position++
						goto l602
					l603:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if buffer[position] != rune('T') {
							goto l591
						}
						position++
					}
				l602:
					if !rules[ruleskip]() {
						goto l591
					}
					depth--
					add(ruleLIMIT, position593)
				}
				if !rules[ruleINTEGER]() {
					goto l591
				}
				depth--
				add(rulelimit, position592)
			}
			return true
		l591:
			position, tokenIndex, depth = position591, tokenIndex591, depth591
			return false
		},
		/* 52 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position604, tokenIndex604, depth604 := position, tokenIndex, depth
			{
				position605 := position
				depth++
				{
					position606 := position
					depth++
					{
						position607, tokenIndex607, depth607 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l608
						}
						position++
						goto l607
					l608:
						position, tokenIndex, depth = position607, tokenIndex607, depth607
						if buffer[position] != rune('O') {
							goto l604
						}
						position++
					}
				l607:
					{
						position609, tokenIndex609, depth609 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l610
						}
						position++
						goto l609
					l610:
						position, tokenIndex, depth = position609, tokenIndex609, depth609
						if buffer[position] != rune('F') {
							goto l604
						}
						position++
					}
				l609:
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('F') {
							goto l604
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('S') {
							goto l604
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('E') {
							goto l604
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('T') {
							goto l604
						}
						position++
					}
				l617:
					if !rules[ruleskip]() {
						goto l604
					}
					depth--
					add(ruleOFFSET, position606)
				}
				if !rules[ruleINTEGER]() {
					goto l604
				}
				depth--
				add(ruleoffset, position605)
			}
			return true
		l604:
			position, tokenIndex, depth = position604, tokenIndex604, depth604
			return false
		},
		/* 53 expression <- <conditionalOrExpression> */
		func() bool {
			position619, tokenIndex619, depth619 := position, tokenIndex, depth
			{
				position620 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l619
				}
				depth--
				add(ruleexpression, position620)
			}
			return true
		l619:
			position, tokenIndex, depth = position619, tokenIndex619, depth619
			return false
		},
		/* 54 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position621, tokenIndex621, depth621 := position, tokenIndex, depth
			{
				position622 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l621
				}
				{
					position623, tokenIndex623, depth623 := position, tokenIndex, depth
					{
						position625 := position
						depth++
						if buffer[position] != rune('|') {
							goto l623
						}
						position++
						if buffer[position] != rune('|') {
							goto l623
						}
						position++
						if !rules[ruleskip]() {
							goto l623
						}
						depth--
						add(ruleOR, position625)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l623
					}
					goto l624
				l623:
					position, tokenIndex, depth = position623, tokenIndex623, depth623
				}
			l624:
				depth--
				add(ruleconditionalOrExpression, position622)
			}
			return true
		l621:
			position, tokenIndex, depth = position621, tokenIndex621, depth621
			return false
		},
		/* 55 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position626, tokenIndex626, depth626 := position, tokenIndex, depth
			{
				position627 := position
				depth++
				{
					position628 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l626
					}
					{
						position629, tokenIndex629, depth629 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position632 := position
									depth++
									{
										position633 := position
										depth++
										{
											position634, tokenIndex634, depth634 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l635
											}
											position++
											goto l634
										l635:
											position, tokenIndex, depth = position634, tokenIndex634, depth634
											if buffer[position] != rune('N') {
												goto l629
											}
											position++
										}
									l634:
										{
											position636, tokenIndex636, depth636 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l637
											}
											position++
											goto l636
										l637:
											position, tokenIndex, depth = position636, tokenIndex636, depth636
											if buffer[position] != rune('O') {
												goto l629
											}
											position++
										}
									l636:
										{
											position638, tokenIndex638, depth638 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l639
											}
											position++
											goto l638
										l639:
											position, tokenIndex, depth = position638, tokenIndex638, depth638
											if buffer[position] != rune('T') {
												goto l629
											}
											position++
										}
									l638:
										if buffer[position] != rune(' ') {
											goto l629
										}
										position++
										{
											position640, tokenIndex640, depth640 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l641
											}
											position++
											goto l640
										l641:
											position, tokenIndex, depth = position640, tokenIndex640, depth640
											if buffer[position] != rune('I') {
												goto l629
											}
											position++
										}
									l640:
										{
											position642, tokenIndex642, depth642 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l643
											}
											position++
											goto l642
										l643:
											position, tokenIndex, depth = position642, tokenIndex642, depth642
											if buffer[position] != rune('N') {
												goto l629
											}
											position++
										}
									l642:
										if !rules[ruleskip]() {
											goto l629
										}
										depth--
										add(ruleNOTIN, position633)
									}
									if !rules[ruleargList]() {
										goto l629
									}
									depth--
									add(rulenotin, position632)
								}
								break
							case 'I', 'i':
								{
									position644 := position
									depth++
									{
										position645 := position
										depth++
										{
											position646, tokenIndex646, depth646 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l647
											}
											position++
											goto l646
										l647:
											position, tokenIndex, depth = position646, tokenIndex646, depth646
											if buffer[position] != rune('I') {
												goto l629
											}
											position++
										}
									l646:
										{
											position648, tokenIndex648, depth648 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l649
											}
											position++
											goto l648
										l649:
											position, tokenIndex, depth = position648, tokenIndex648, depth648
											if buffer[position] != rune('N') {
												goto l629
											}
											position++
										}
									l648:
										if !rules[ruleskip]() {
											goto l629
										}
										depth--
										add(ruleIN, position645)
									}
									if !rules[ruleargList]() {
										goto l629
									}
									depth--
									add(rulein, position644)
								}
								break
							default:
								{
									position650, tokenIndex650, depth650 := position, tokenIndex, depth
									{
										position652 := position
										depth++
										if buffer[position] != rune('<') {
											goto l651
										}
										position++
										if !rules[ruleskip]() {
											goto l651
										}
										depth--
										add(ruleLT, position652)
									}
									goto l650
								l651:
									position, tokenIndex, depth = position650, tokenIndex650, depth650
									{
										position654 := position
										depth++
										if buffer[position] != rune('>') {
											goto l653
										}
										position++
										if buffer[position] != rune('=') {
											goto l653
										}
										position++
										if !rules[ruleskip]() {
											goto l653
										}
										depth--
										add(ruleGE, position654)
									}
									goto l650
								l653:
									position, tokenIndex, depth = position650, tokenIndex650, depth650
									{
										switch buffer[position] {
										case '>':
											{
												position656 := position
												depth++
												if buffer[position] != rune('>') {
													goto l629
												}
												position++
												if !rules[ruleskip]() {
													goto l629
												}
												depth--
												add(ruleGT, position656)
											}
											break
										case '<':
											{
												position657 := position
												depth++
												if buffer[position] != rune('<') {
													goto l629
												}
												position++
												if buffer[position] != rune('=') {
													goto l629
												}
												position++
												if !rules[ruleskip]() {
													goto l629
												}
												depth--
												add(ruleLE, position657)
											}
											break
										case '!':
											{
												position658 := position
												depth++
												if buffer[position] != rune('!') {
													goto l629
												}
												position++
												if buffer[position] != rune('=') {
													goto l629
												}
												position++
												if !rules[ruleskip]() {
													goto l629
												}
												depth--
												add(ruleNE, position658)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l629
											}
											break
										}
									}

								}
							l650:
								if !rules[rulenumericExpression]() {
									goto l629
								}
								break
							}
						}

						goto l630
					l629:
						position, tokenIndex, depth = position629, tokenIndex629, depth629
					}
				l630:
					depth--
					add(rulevalueLogical, position628)
				}
				{
					position659, tokenIndex659, depth659 := position, tokenIndex, depth
					{
						position661 := position
						depth++
						if buffer[position] != rune('&') {
							goto l659
						}
						position++
						if buffer[position] != rune('&') {
							goto l659
						}
						position++
						if !rules[ruleskip]() {
							goto l659
						}
						depth--
						add(ruleAND, position661)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l659
					}
					goto l660
				l659:
					position, tokenIndex, depth = position659, tokenIndex659, depth659
				}
			l660:
				depth--
				add(ruleconditionalAndExpression, position627)
			}
			return true
		l626:
			position, tokenIndex, depth = position626, tokenIndex626, depth626
			return false
		},
		/* 56 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 57 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position663, tokenIndex663, depth663 := position, tokenIndex, depth
			{
				position664 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l663
				}
			l665:
				{
					position666, tokenIndex666, depth666 := position, tokenIndex, depth
					{
						position667, tokenIndex667, depth667 := position, tokenIndex, depth
						{
							position669, tokenIndex669, depth669 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l670
							}
							goto l669
						l670:
							position, tokenIndex, depth = position669, tokenIndex669, depth669
							if !rules[ruleMINUS]() {
								goto l668
							}
						}
					l669:
						if !rules[rulemultiplicativeExpression]() {
							goto l668
						}
						goto l667
					l668:
						position, tokenIndex, depth = position667, tokenIndex667, depth667
						{
							position671 := position
							depth++
							{
								position672, tokenIndex672, depth672 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l673
								}
								position++
								goto l672
							l673:
								position, tokenIndex, depth = position672, tokenIndex672, depth672
								if buffer[position] != rune('-') {
									goto l666
								}
								position++
							}
						l672:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l666
							}
							position++
						l674:
							{
								position675, tokenIndex675, depth675 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l675
								}
								position++
								goto l674
							l675:
								position, tokenIndex, depth = position675, tokenIndex675, depth675
							}
							{
								position676, tokenIndex676, depth676 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l676
								}
								position++
							l678:
								{
									position679, tokenIndex679, depth679 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l679
									}
									position++
									goto l678
								l679:
									position, tokenIndex, depth = position679, tokenIndex679, depth679
								}
								goto l677
							l676:
								position, tokenIndex, depth = position676, tokenIndex676, depth676
							}
						l677:
							if !rules[ruleskip]() {
								goto l666
							}
							depth--
							add(rulesignedNumericLiteral, position671)
						}
					}
				l667:
					goto l665
				l666:
					position, tokenIndex, depth = position666, tokenIndex666, depth666
				}
				depth--
				add(rulenumericExpression, position664)
			}
			return true
		l663:
			position, tokenIndex, depth = position663, tokenIndex663, depth663
			return false
		},
		/* 58 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position680, tokenIndex680, depth680 := position, tokenIndex, depth
			{
				position681 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l680
				}
			l682:
				{
					position683, tokenIndex683, depth683 := position, tokenIndex, depth
					{
						position684, tokenIndex684, depth684 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l685
						}
						goto l684
					l685:
						position, tokenIndex, depth = position684, tokenIndex684, depth684
						if !rules[ruleSLASH]() {
							goto l683
						}
					}
				l684:
					if !rules[ruleunaryExpression]() {
						goto l683
					}
					goto l682
				l683:
					position, tokenIndex, depth = position683, tokenIndex683, depth683
				}
				depth--
				add(rulemultiplicativeExpression, position681)
			}
			return true
		l680:
			position, tokenIndex, depth = position680, tokenIndex680, depth680
			return false
		},
		/* 59 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position686, tokenIndex686, depth686 := position, tokenIndex, depth
			{
				position687 := position
				depth++
				{
					position688, tokenIndex688, depth688 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l688
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l688
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l688
							}
							break
						}
					}

					goto l689
				l688:
					position, tokenIndex, depth = position688, tokenIndex688, depth688
				}
			l689:
				{
					position691 := position
					depth++
					{
						position692, tokenIndex692, depth692 := position, tokenIndex, depth
						if !rules[rulebuiltinCall]() {
							goto l693
						}
						goto l692
					l693:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						if !rules[rulefunctionCall]() {
							goto l694
						}
						goto l692
					l694:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						if !rules[ruleiriref]() {
							goto l695
						}
						goto l692
					l695:
						position, tokenIndex, depth = position692, tokenIndex692, depth692
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position697 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position699 := position
												depth++
												{
													position700 := position
													depth++
													{
														position701, tokenIndex701, depth701 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l702
														}
														position++
														goto l701
													l702:
														position, tokenIndex, depth = position701, tokenIndex701, depth701
														if buffer[position] != rune('G') {
															goto l686
														}
														position++
													}
												l701:
													{
														position703, tokenIndex703, depth703 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l704
														}
														position++
														goto l703
													l704:
														position, tokenIndex, depth = position703, tokenIndex703, depth703
														if buffer[position] != rune('R') {
															goto l686
														}
														position++
													}
												l703:
													{
														position705, tokenIndex705, depth705 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l706
														}
														position++
														goto l705
													l706:
														position, tokenIndex, depth = position705, tokenIndex705, depth705
														if buffer[position] != rune('O') {
															goto l686
														}
														position++
													}
												l705:
													{
														position707, tokenIndex707, depth707 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l708
														}
														position++
														goto l707
													l708:
														position, tokenIndex, depth = position707, tokenIndex707, depth707
														if buffer[position] != rune('U') {
															goto l686
														}
														position++
													}
												l707:
													{
														position709, tokenIndex709, depth709 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l710
														}
														position++
														goto l709
													l710:
														position, tokenIndex, depth = position709, tokenIndex709, depth709
														if buffer[position] != rune('P') {
															goto l686
														}
														position++
													}
												l709:
													if buffer[position] != rune('_') {
														goto l686
													}
													position++
													{
														position711, tokenIndex711, depth711 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l712
														}
														position++
														goto l711
													l712:
														position, tokenIndex, depth = position711, tokenIndex711, depth711
														if buffer[position] != rune('C') {
															goto l686
														}
														position++
													}
												l711:
													{
														position713, tokenIndex713, depth713 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l714
														}
														position++
														goto l713
													l714:
														position, tokenIndex, depth = position713, tokenIndex713, depth713
														if buffer[position] != rune('O') {
															goto l686
														}
														position++
													}
												l713:
													{
														position715, tokenIndex715, depth715 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l716
														}
														position++
														goto l715
													l716:
														position, tokenIndex, depth = position715, tokenIndex715, depth715
														if buffer[position] != rune('N') {
															goto l686
														}
														position++
													}
												l715:
													{
														position717, tokenIndex717, depth717 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l718
														}
														position++
														goto l717
													l718:
														position, tokenIndex, depth = position717, tokenIndex717, depth717
														if buffer[position] != rune('C') {
															goto l686
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
															goto l686
														}
														position++
													}
												l719:
													{
														position721, tokenIndex721, depth721 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l722
														}
														position++
														goto l721
													l722:
														position, tokenIndex, depth = position721, tokenIndex721, depth721
														if buffer[position] != rune('T') {
															goto l686
														}
														position++
													}
												l721:
													if !rules[ruleskip]() {
														goto l686
													}
													depth--
													add(ruleGROUPCONCAT, position700)
												}
												if !rules[ruleLPAREN]() {
													goto l686
												}
												{
													position723, tokenIndex723, depth723 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l723
													}
													goto l724
												l723:
													position, tokenIndex, depth = position723, tokenIndex723, depth723
												}
											l724:
												if !rules[ruleexpression]() {
													goto l686
												}
												{
													position725, tokenIndex725, depth725 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l725
													}
													{
														position727 := position
														depth++
														{
															position728, tokenIndex728, depth728 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l729
															}
															position++
															goto l728
														l729:
															position, tokenIndex, depth = position728, tokenIndex728, depth728
															if buffer[position] != rune('S') {
																goto l725
															}
															position++
														}
													l728:
														{
															position730, tokenIndex730, depth730 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l731
															}
															position++
															goto l730
														l731:
															position, tokenIndex, depth = position730, tokenIndex730, depth730
															if buffer[position] != rune('E') {
																goto l725
															}
															position++
														}
													l730:
														{
															position732, tokenIndex732, depth732 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l733
															}
															position++
															goto l732
														l733:
															position, tokenIndex, depth = position732, tokenIndex732, depth732
															if buffer[position] != rune('P') {
																goto l725
															}
															position++
														}
													l732:
														{
															position734, tokenIndex734, depth734 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l735
															}
															position++
															goto l734
														l735:
															position, tokenIndex, depth = position734, tokenIndex734, depth734
															if buffer[position] != rune('A') {
																goto l725
															}
															position++
														}
													l734:
														{
															position736, tokenIndex736, depth736 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l737
															}
															position++
															goto l736
														l737:
															position, tokenIndex, depth = position736, tokenIndex736, depth736
															if buffer[position] != rune('R') {
																goto l725
															}
															position++
														}
													l736:
														{
															position738, tokenIndex738, depth738 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l739
															}
															position++
															goto l738
														l739:
															position, tokenIndex, depth = position738, tokenIndex738, depth738
															if buffer[position] != rune('A') {
																goto l725
															}
															position++
														}
													l738:
														{
															position740, tokenIndex740, depth740 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l741
															}
															position++
															goto l740
														l741:
															position, tokenIndex, depth = position740, tokenIndex740, depth740
															if buffer[position] != rune('T') {
																goto l725
															}
															position++
														}
													l740:
														{
															position742, tokenIndex742, depth742 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l743
															}
															position++
															goto l742
														l743:
															position, tokenIndex, depth = position742, tokenIndex742, depth742
															if buffer[position] != rune('O') {
																goto l725
															}
															position++
														}
													l742:
														{
															position744, tokenIndex744, depth744 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l745
															}
															position++
															goto l744
														l745:
															position, tokenIndex, depth = position744, tokenIndex744, depth744
															if buffer[position] != rune('R') {
																goto l725
															}
															position++
														}
													l744:
														if !rules[ruleskip]() {
															goto l725
														}
														depth--
														add(ruleSEPARATOR, position727)
													}
													if !rules[ruleEQ]() {
														goto l725
													}
													if !rules[rulestring]() {
														goto l725
													}
													goto l726
												l725:
													position, tokenIndex, depth = position725, tokenIndex725, depth725
												}
											l726:
												if !rules[ruleRPAREN]() {
													goto l686
												}
												depth--
												add(rulegroupConcat, position699)
											}
											break
										case 'C', 'c':
											{
												position746 := position
												depth++
												{
													position747 := position
													depth++
													{
														position748, tokenIndex748, depth748 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l749
														}
														position++
														goto l748
													l749:
														position, tokenIndex, depth = position748, tokenIndex748, depth748
														if buffer[position] != rune('C') {
															goto l686
														}
														position++
													}
												l748:
													{
														position750, tokenIndex750, depth750 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l751
														}
														position++
														goto l750
													l751:
														position, tokenIndex, depth = position750, tokenIndex750, depth750
														if buffer[position] != rune('O') {
															goto l686
														}
														position++
													}
												l750:
													{
														position752, tokenIndex752, depth752 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l753
														}
														position++
														goto l752
													l753:
														position, tokenIndex, depth = position752, tokenIndex752, depth752
														if buffer[position] != rune('U') {
															goto l686
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
															goto l686
														}
														position++
													}
												l754:
													{
														position756, tokenIndex756, depth756 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l757
														}
														position++
														goto l756
													l757:
														position, tokenIndex, depth = position756, tokenIndex756, depth756
														if buffer[position] != rune('T') {
															goto l686
														}
														position++
													}
												l756:
													if !rules[ruleskip]() {
														goto l686
													}
													depth--
													add(ruleCOUNT, position747)
												}
												if !rules[ruleLPAREN]() {
													goto l686
												}
												{
													position758, tokenIndex758, depth758 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l758
													}
													goto l759
												l758:
													position, tokenIndex, depth = position758, tokenIndex758, depth758
												}
											l759:
												{
													position760, tokenIndex760, depth760 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l761
													}
													goto l760
												l761:
													position, tokenIndex, depth = position760, tokenIndex760, depth760
													if !rules[ruleexpression]() {
														goto l686
													}
												}
											l760:
												if !rules[ruleRPAREN]() {
													goto l686
												}
												depth--
												add(rulecount, position746)
											}
											break
										default:
											{
												position762, tokenIndex762, depth762 := position, tokenIndex, depth
												{
													position764 := position
													depth++
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
															goto l763
														}
														position++
													}
												l765:
													{
														position767, tokenIndex767, depth767 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l768
														}
														position++
														goto l767
													l768:
														position, tokenIndex, depth = position767, tokenIndex767, depth767
														if buffer[position] != rune('U') {
															goto l763
														}
														position++
													}
												l767:
													{
														position769, tokenIndex769, depth769 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l770
														}
														position++
														goto l769
													l770:
														position, tokenIndex, depth = position769, tokenIndex769, depth769
														if buffer[position] != rune('M') {
															goto l763
														}
														position++
													}
												l769:
													if !rules[ruleskip]() {
														goto l763
													}
													depth--
													add(ruleSUM, position764)
												}
												goto l762
											l763:
												position, tokenIndex, depth = position762, tokenIndex762, depth762
												{
													position772 := position
													depth++
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
															goto l771
														}
														position++
													}
												l773:
													{
														position775, tokenIndex775, depth775 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l776
														}
														position++
														goto l775
													l776:
														position, tokenIndex, depth = position775, tokenIndex775, depth775
														if buffer[position] != rune('I') {
															goto l771
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
															goto l771
														}
														position++
													}
												l777:
													if !rules[ruleskip]() {
														goto l771
													}
													depth--
													add(ruleMIN, position772)
												}
												goto l762
											l771:
												position, tokenIndex, depth = position762, tokenIndex762, depth762
												{
													switch buffer[position] {
													case 'S', 's':
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
																	goto l686
																}
																position++
															}
														l781:
															{
																position783, tokenIndex783, depth783 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l784
																}
																position++
																goto l783
															l784:
																position, tokenIndex, depth = position783, tokenIndex783, depth783
																if buffer[position] != rune('A') {
																	goto l686
																}
																position++
															}
														l783:
															{
																position785, tokenIndex785, depth785 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l786
																}
																position++
																goto l785
															l786:
																position, tokenIndex, depth = position785, tokenIndex785, depth785
																if buffer[position] != rune('M') {
																	goto l686
																}
																position++
															}
														l785:
															{
																position787, tokenIndex787, depth787 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l788
																}
																position++
																goto l787
															l788:
																position, tokenIndex, depth = position787, tokenIndex787, depth787
																if buffer[position] != rune('P') {
																	goto l686
																}
																position++
															}
														l787:
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
																	goto l686
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
																	goto l686
																}
																position++
															}
														l791:
															if !rules[ruleskip]() {
																goto l686
															}
															depth--
															add(ruleSAMPLE, position780)
														}
														break
													case 'A', 'a':
														{
															position793 := position
															depth++
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
																	goto l686
																}
																position++
															}
														l794:
															{
																position796, tokenIndex796, depth796 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l797
																}
																position++
																goto l796
															l797:
																position, tokenIndex, depth = position796, tokenIndex796, depth796
																if buffer[position] != rune('V') {
																	goto l686
																}
																position++
															}
														l796:
															{
																position798, tokenIndex798, depth798 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l799
																}
																position++
																goto l798
															l799:
																position, tokenIndex, depth = position798, tokenIndex798, depth798
																if buffer[position] != rune('G') {
																	goto l686
																}
																position++
															}
														l798:
															if !rules[ruleskip]() {
																goto l686
															}
															depth--
															add(ruleAVG, position793)
														}
														break
													default:
														{
															position800 := position
															depth++
															{
																position801, tokenIndex801, depth801 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l802
																}
																position++
																goto l801
															l802:
																position, tokenIndex, depth = position801, tokenIndex801, depth801
																if buffer[position] != rune('M') {
																	goto l686
																}
																position++
															}
														l801:
															{
																position803, tokenIndex803, depth803 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l804
																}
																position++
																goto l803
															l804:
																position, tokenIndex, depth = position803, tokenIndex803, depth803
																if buffer[position] != rune('A') {
																	goto l686
																}
																position++
															}
														l803:
															{
																position805, tokenIndex805, depth805 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l806
																}
																position++
																goto l805
															l806:
																position, tokenIndex, depth = position805, tokenIndex805, depth805
																if buffer[position] != rune('X') {
																	goto l686
																}
																position++
															}
														l805:
															if !rules[ruleskip]() {
																goto l686
															}
															depth--
															add(ruleMAX, position800)
														}
														break
													}
												}

											}
										l762:
											if !rules[ruleLPAREN]() {
												goto l686
											}
											{
												position807, tokenIndex807, depth807 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l807
												}
												goto l808
											l807:
												position, tokenIndex, depth = position807, tokenIndex807, depth807
											}
										l808:
											if !rules[ruleexpression]() {
												goto l686
											}
											if !rules[ruleRPAREN]() {
												goto l686
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position697)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l686
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l686
								}
								break
							case '(':
								if !rules[rulebrackettedExpression]() {
									goto l686
								}
								break
							case '"', '\'':
								if !rules[ruleliteral]() {
									goto l686
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l686
								}
								break
							}
						}

					}
				l692:
					depth--
					add(ruleprimaryExpression, position691)
				}
				depth--
				add(ruleunaryExpression, position687)
			}
			return true
		l686:
			position, tokenIndex, depth = position686, tokenIndex686, depth686
			return false
		},
		/* 60 primaryExpression <- <(builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('(') brackettedExpression) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 61 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position810, tokenIndex810, depth810 := position, tokenIndex, depth
			{
				position811 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l810
				}
				if !rules[ruleexpression]() {
					goto l810
				}
				if !rules[ruleRPAREN]() {
					goto l810
				}
				depth--
				add(rulebrackettedExpression, position811)
			}
			return true
		l810:
			position, tokenIndex, depth = position810, tokenIndex810, depth810
			return false
		},
		/* 62 functionCall <- <(iriref argList)> */
		func() bool {
			position812, tokenIndex812, depth812 := position, tokenIndex, depth
			{
				position813 := position
				depth++
				if !rules[ruleiriref]() {
					goto l812
				}
				if !rules[ruleargList]() {
					goto l812
				}
				depth--
				add(rulefunctionCall, position813)
			}
			return true
		l812:
			position, tokenIndex, depth = position812, tokenIndex812, depth812
			return false
		},
		/* 63 in <- <(IN argList)> */
		nil,
		/* 64 notin <- <(NOTIN argList)> */
		nil,
		/* 65 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position816, tokenIndex816, depth816 := position, tokenIndex, depth
			{
				position817 := position
				depth++
				{
					position818, tokenIndex818, depth818 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l819
					}
					goto l818
				l819:
					position, tokenIndex, depth = position818, tokenIndex818, depth818
					if !rules[ruleLPAREN]() {
						goto l816
					}
					if !rules[ruleexpression]() {
						goto l816
					}
				l820:
					{
						position821, tokenIndex821, depth821 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l821
						}
						if !rules[ruleexpression]() {
							goto l821
						}
						goto l820
					l821:
						position, tokenIndex, depth = position821, tokenIndex821, depth821
					}
					if !rules[ruleRPAREN]() {
						goto l816
					}
				}
			l818:
				depth--
				add(ruleargList, position817)
			}
			return true
		l816:
			position, tokenIndex, depth = position816, tokenIndex816, depth816
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
			position825, tokenIndex825, depth825 := position, tokenIndex, depth
			{
				position826 := position
				depth++
				{
					position827, tokenIndex827, depth827 := position, tokenIndex, depth
					{
						position829, tokenIndex829, depth829 := position, tokenIndex, depth
						{
							position831 := position
							depth++
							{
								position832, tokenIndex832, depth832 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l833
								}
								position++
								goto l832
							l833:
								position, tokenIndex, depth = position832, tokenIndex832, depth832
								if buffer[position] != rune('S') {
									goto l830
								}
								position++
							}
						l832:
							{
								position834, tokenIndex834, depth834 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l835
								}
								position++
								goto l834
							l835:
								position, tokenIndex, depth = position834, tokenIndex834, depth834
								if buffer[position] != rune('T') {
									goto l830
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
									goto l830
								}
								position++
							}
						l836:
							if !rules[ruleskip]() {
								goto l830
							}
							depth--
							add(ruleSTR, position831)
						}
						goto l829
					l830:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position839 := position
							depth++
							{
								position840, tokenIndex840, depth840 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l841
								}
								position++
								goto l840
							l841:
								position, tokenIndex, depth = position840, tokenIndex840, depth840
								if buffer[position] != rune('L') {
									goto l838
								}
								position++
							}
						l840:
							{
								position842, tokenIndex842, depth842 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l843
								}
								position++
								goto l842
							l843:
								position, tokenIndex, depth = position842, tokenIndex842, depth842
								if buffer[position] != rune('A') {
									goto l838
								}
								position++
							}
						l842:
							{
								position844, tokenIndex844, depth844 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l845
								}
								position++
								goto l844
							l845:
								position, tokenIndex, depth = position844, tokenIndex844, depth844
								if buffer[position] != rune('N') {
									goto l838
								}
								position++
							}
						l844:
							{
								position846, tokenIndex846, depth846 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l847
								}
								position++
								goto l846
							l847:
								position, tokenIndex, depth = position846, tokenIndex846, depth846
								if buffer[position] != rune('G') {
									goto l838
								}
								position++
							}
						l846:
							if !rules[ruleskip]() {
								goto l838
							}
							depth--
							add(ruleLANG, position839)
						}
						goto l829
					l838:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position849 := position
							depth++
							{
								position850, tokenIndex850, depth850 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l851
								}
								position++
								goto l850
							l851:
								position, tokenIndex, depth = position850, tokenIndex850, depth850
								if buffer[position] != rune('D') {
									goto l848
								}
								position++
							}
						l850:
							{
								position852, tokenIndex852, depth852 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l853
								}
								position++
								goto l852
							l853:
								position, tokenIndex, depth = position852, tokenIndex852, depth852
								if buffer[position] != rune('A') {
									goto l848
								}
								position++
							}
						l852:
							{
								position854, tokenIndex854, depth854 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l855
								}
								position++
								goto l854
							l855:
								position, tokenIndex, depth = position854, tokenIndex854, depth854
								if buffer[position] != rune('T') {
									goto l848
								}
								position++
							}
						l854:
							{
								position856, tokenIndex856, depth856 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l857
								}
								position++
								goto l856
							l857:
								position, tokenIndex, depth = position856, tokenIndex856, depth856
								if buffer[position] != rune('A') {
									goto l848
								}
								position++
							}
						l856:
							{
								position858, tokenIndex858, depth858 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('T') {
									goto l848
								}
								position++
							}
						l858:
							{
								position860, tokenIndex860, depth860 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l861
								}
								position++
								goto l860
							l861:
								position, tokenIndex, depth = position860, tokenIndex860, depth860
								if buffer[position] != rune('Y') {
									goto l848
								}
								position++
							}
						l860:
							{
								position862, tokenIndex862, depth862 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l863
								}
								position++
								goto l862
							l863:
								position, tokenIndex, depth = position862, tokenIndex862, depth862
								if buffer[position] != rune('P') {
									goto l848
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
									goto l848
								}
								position++
							}
						l864:
							if !rules[ruleskip]() {
								goto l848
							}
							depth--
							add(ruleDATATYPE, position849)
						}
						goto l829
					l848:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position867 := position
							depth++
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
									goto l866
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
									goto l866
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
									goto l866
								}
								position++
							}
						l872:
							if !rules[ruleskip]() {
								goto l866
							}
							depth--
							add(ruleIRI, position867)
						}
						goto l829
					l866:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position875 := position
							depth++
							{
								position876, tokenIndex876, depth876 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l877
								}
								position++
								goto l876
							l877:
								position, tokenIndex, depth = position876, tokenIndex876, depth876
								if buffer[position] != rune('U') {
									goto l874
								}
								position++
							}
						l876:
							{
								position878, tokenIndex878, depth878 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l879
								}
								position++
								goto l878
							l879:
								position, tokenIndex, depth = position878, tokenIndex878, depth878
								if buffer[position] != rune('R') {
									goto l874
								}
								position++
							}
						l878:
							{
								position880, tokenIndex880, depth880 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l881
								}
								position++
								goto l880
							l881:
								position, tokenIndex, depth = position880, tokenIndex880, depth880
								if buffer[position] != rune('I') {
									goto l874
								}
								position++
							}
						l880:
							if !rules[ruleskip]() {
								goto l874
							}
							depth--
							add(ruleURI, position875)
						}
						goto l829
					l874:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position883 := position
							depth++
							{
								position884, tokenIndex884, depth884 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l885
								}
								position++
								goto l884
							l885:
								position, tokenIndex, depth = position884, tokenIndex884, depth884
								if buffer[position] != rune('S') {
									goto l882
								}
								position++
							}
						l884:
							{
								position886, tokenIndex886, depth886 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l887
								}
								position++
								goto l886
							l887:
								position, tokenIndex, depth = position886, tokenIndex886, depth886
								if buffer[position] != rune('T') {
									goto l882
								}
								position++
							}
						l886:
							{
								position888, tokenIndex888, depth888 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l889
								}
								position++
								goto l888
							l889:
								position, tokenIndex, depth = position888, tokenIndex888, depth888
								if buffer[position] != rune('R') {
									goto l882
								}
								position++
							}
						l888:
							{
								position890, tokenIndex890, depth890 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l891
								}
								position++
								goto l890
							l891:
								position, tokenIndex, depth = position890, tokenIndex890, depth890
								if buffer[position] != rune('L') {
									goto l882
								}
								position++
							}
						l890:
							{
								position892, tokenIndex892, depth892 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l893
								}
								position++
								goto l892
							l893:
								position, tokenIndex, depth = position892, tokenIndex892, depth892
								if buffer[position] != rune('E') {
									goto l882
								}
								position++
							}
						l892:
							{
								position894, tokenIndex894, depth894 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l895
								}
								position++
								goto l894
							l895:
								position, tokenIndex, depth = position894, tokenIndex894, depth894
								if buffer[position] != rune('N') {
									goto l882
								}
								position++
							}
						l894:
							if !rules[ruleskip]() {
								goto l882
							}
							depth--
							add(ruleSTRLEN, position883)
						}
						goto l829
					l882:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
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
								if buffer[position] != rune('o') {
									goto l901
								}
								position++
								goto l900
							l901:
								position, tokenIndex, depth = position900, tokenIndex900, depth900
								if buffer[position] != rune('O') {
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
								if buffer[position] != rune('t') {
									goto l905
								}
								position++
								goto l904
							l905:
								position, tokenIndex, depth = position904, tokenIndex904, depth904
								if buffer[position] != rune('T') {
									goto l896
								}
								position++
							}
						l904:
							{
								position906, tokenIndex906, depth906 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l907
								}
								position++
								goto l906
							l907:
								position, tokenIndex, depth = position906, tokenIndex906, depth906
								if buffer[position] != rune('H') {
									goto l896
								}
								position++
							}
						l906:
							if !rules[ruleskip]() {
								goto l896
							}
							depth--
							add(ruleMONTH, position897)
						}
						goto l829
					l896:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position909 := position
							depth++
							{
								position910, tokenIndex910, depth910 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l911
								}
								position++
								goto l910
							l911:
								position, tokenIndex, depth = position910, tokenIndex910, depth910
								if buffer[position] != rune('M') {
									goto l908
								}
								position++
							}
						l910:
							{
								position912, tokenIndex912, depth912 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l913
								}
								position++
								goto l912
							l913:
								position, tokenIndex, depth = position912, tokenIndex912, depth912
								if buffer[position] != rune('I') {
									goto l908
								}
								position++
							}
						l912:
							{
								position914, tokenIndex914, depth914 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l915
								}
								position++
								goto l914
							l915:
								position, tokenIndex, depth = position914, tokenIndex914, depth914
								if buffer[position] != rune('N') {
									goto l908
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
									goto l908
								}
								position++
							}
						l916:
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
									goto l908
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
									goto l908
								}
								position++
							}
						l920:
							{
								position922, tokenIndex922, depth922 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l923
								}
								position++
								goto l922
							l923:
								position, tokenIndex, depth = position922, tokenIndex922, depth922
								if buffer[position] != rune('S') {
									goto l908
								}
								position++
							}
						l922:
							if !rules[ruleskip]() {
								goto l908
							}
							depth--
							add(ruleMINUTES, position909)
						}
						goto l829
					l908:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
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
								if buffer[position] != rune('e') {
									goto l929
								}
								position++
								goto l928
							l929:
								position, tokenIndex, depth = position928, tokenIndex928, depth928
								if buffer[position] != rune('E') {
									goto l924
								}
								position++
							}
						l928:
							{
								position930, tokenIndex930, depth930 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l931
								}
								position++
								goto l930
							l931:
								position, tokenIndex, depth = position930, tokenIndex930, depth930
								if buffer[position] != rune('C') {
									goto l924
								}
								position++
							}
						l930:
							{
								position932, tokenIndex932, depth932 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l933
								}
								position++
								goto l932
							l933:
								position, tokenIndex, depth = position932, tokenIndex932, depth932
								if buffer[position] != rune('O') {
									goto l924
								}
								position++
							}
						l932:
							{
								position934, tokenIndex934, depth934 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l935
								}
								position++
								goto l934
							l935:
								position, tokenIndex, depth = position934, tokenIndex934, depth934
								if buffer[position] != rune('N') {
									goto l924
								}
								position++
							}
						l934:
							{
								position936, tokenIndex936, depth936 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l937
								}
								position++
								goto l936
							l937:
								position, tokenIndex, depth = position936, tokenIndex936, depth936
								if buffer[position] != rune('D') {
									goto l924
								}
								position++
							}
						l936:
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
									goto l924
								}
								position++
							}
						l938:
							if !rules[ruleskip]() {
								goto l924
							}
							depth--
							add(ruleSECONDS, position925)
						}
						goto l829
					l924:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position941 := position
							depth++
							{
								position942, tokenIndex942, depth942 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l943
								}
								position++
								goto l942
							l943:
								position, tokenIndex, depth = position942, tokenIndex942, depth942
								if buffer[position] != rune('T') {
									goto l940
								}
								position++
							}
						l942:
							{
								position944, tokenIndex944, depth944 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l945
								}
								position++
								goto l944
							l945:
								position, tokenIndex, depth = position944, tokenIndex944, depth944
								if buffer[position] != rune('I') {
									goto l940
								}
								position++
							}
						l944:
							{
								position946, tokenIndex946, depth946 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l947
								}
								position++
								goto l946
							l947:
								position, tokenIndex, depth = position946, tokenIndex946, depth946
								if buffer[position] != rune('M') {
									goto l940
								}
								position++
							}
						l946:
							{
								position948, tokenIndex948, depth948 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l949
								}
								position++
								goto l948
							l949:
								position, tokenIndex, depth = position948, tokenIndex948, depth948
								if buffer[position] != rune('E') {
									goto l940
								}
								position++
							}
						l948:
							{
								position950, tokenIndex950, depth950 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l951
								}
								position++
								goto l950
							l951:
								position, tokenIndex, depth = position950, tokenIndex950, depth950
								if buffer[position] != rune('Z') {
									goto l940
								}
								position++
							}
						l950:
							{
								position952, tokenIndex952, depth952 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l953
								}
								position++
								goto l952
							l953:
								position, tokenIndex, depth = position952, tokenIndex952, depth952
								if buffer[position] != rune('O') {
									goto l940
								}
								position++
							}
						l952:
							{
								position954, tokenIndex954, depth954 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l955
								}
								position++
								goto l954
							l955:
								position, tokenIndex, depth = position954, tokenIndex954, depth954
								if buffer[position] != rune('N') {
									goto l940
								}
								position++
							}
						l954:
							{
								position956, tokenIndex956, depth956 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l957
								}
								position++
								goto l956
							l957:
								position, tokenIndex, depth = position956, tokenIndex956, depth956
								if buffer[position] != rune('E') {
									goto l940
								}
								position++
							}
						l956:
							if !rules[ruleskip]() {
								goto l940
							}
							depth--
							add(ruleTIMEZONE, position941)
						}
						goto l829
					l940:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position959 := position
							depth++
							{
								position960, tokenIndex960, depth960 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l961
								}
								position++
								goto l960
							l961:
								position, tokenIndex, depth = position960, tokenIndex960, depth960
								if buffer[position] != rune('S') {
									goto l958
								}
								position++
							}
						l960:
							{
								position962, tokenIndex962, depth962 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l963
								}
								position++
								goto l962
							l963:
								position, tokenIndex, depth = position962, tokenIndex962, depth962
								if buffer[position] != rune('H') {
									goto l958
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
									goto l958
								}
								position++
							}
						l964:
							if buffer[position] != rune('1') {
								goto l958
							}
							position++
							if !rules[ruleskip]() {
								goto l958
							}
							depth--
							add(ruleSHA1, position959)
						}
						goto l829
					l958:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position967 := position
							depth++
							{
								position968, tokenIndex968, depth968 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l969
								}
								position++
								goto l968
							l969:
								position, tokenIndex, depth = position968, tokenIndex968, depth968
								if buffer[position] != rune('S') {
									goto l966
								}
								position++
							}
						l968:
							{
								position970, tokenIndex970, depth970 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l971
								}
								position++
								goto l970
							l971:
								position, tokenIndex, depth = position970, tokenIndex970, depth970
								if buffer[position] != rune('H') {
									goto l966
								}
								position++
							}
						l970:
							{
								position972, tokenIndex972, depth972 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l973
								}
								position++
								goto l972
							l973:
								position, tokenIndex, depth = position972, tokenIndex972, depth972
								if buffer[position] != rune('A') {
									goto l966
								}
								position++
							}
						l972:
							if buffer[position] != rune('2') {
								goto l966
							}
							position++
							if buffer[position] != rune('5') {
								goto l966
							}
							position++
							if buffer[position] != rune('6') {
								goto l966
							}
							position++
							if !rules[ruleskip]() {
								goto l966
							}
							depth--
							add(ruleSHA256, position967)
						}
						goto l829
					l966:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position975 := position
							depth++
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
									goto l974
								}
								position++
							}
						l976:
							{
								position978, tokenIndex978, depth978 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l979
								}
								position++
								goto l978
							l979:
								position, tokenIndex, depth = position978, tokenIndex978, depth978
								if buffer[position] != rune('H') {
									goto l974
								}
								position++
							}
						l978:
							{
								position980, tokenIndex980, depth980 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l981
								}
								position++
								goto l980
							l981:
								position, tokenIndex, depth = position980, tokenIndex980, depth980
								if buffer[position] != rune('A') {
									goto l974
								}
								position++
							}
						l980:
							if buffer[position] != rune('3') {
								goto l974
							}
							position++
							if buffer[position] != rune('8') {
								goto l974
							}
							position++
							if buffer[position] != rune('4') {
								goto l974
							}
							position++
							if !rules[ruleskip]() {
								goto l974
							}
							depth--
							add(ruleSHA384, position975)
						}
						goto l829
					l974:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
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
								if buffer[position] != rune('i') {
									goto l989
								}
								position++
								goto l988
							l989:
								position, tokenIndex, depth = position988, tokenIndex988, depth988
								if buffer[position] != rune('I') {
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
							add(ruleISIRI, position983)
						}
						goto l829
					l982:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
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
								if buffer[position] != rune('u') {
									goto l1001
								}
								position++
								goto l1000
							l1001:
								position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
								if buffer[position] != rune('U') {
									goto l994
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
									goto l994
								}
								position++
							}
						l1002:
							{
								position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1005
								}
								position++
								goto l1004
							l1005:
								position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
								if buffer[position] != rune('I') {
									goto l994
								}
								position++
							}
						l1004:
							if !rules[ruleskip]() {
								goto l994
							}
							depth--
							add(ruleISURI, position995)
						}
						goto l829
					l994:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position1007 := position
							depth++
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
									goto l1006
								}
								position++
							}
						l1008:
							{
								position1010, tokenIndex1010, depth1010 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1011
								}
								position++
								goto l1010
							l1011:
								position, tokenIndex, depth = position1010, tokenIndex1010, depth1010
								if buffer[position] != rune('S') {
									goto l1006
								}
								position++
							}
						l1010:
							{
								position1012, tokenIndex1012, depth1012 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1013
								}
								position++
								goto l1012
							l1013:
								position, tokenIndex, depth = position1012, tokenIndex1012, depth1012
								if buffer[position] != rune('B') {
									goto l1006
								}
								position++
							}
						l1012:
							{
								position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1015
								}
								position++
								goto l1014
							l1015:
								position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
								if buffer[position] != rune('L') {
									goto l1006
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
									goto l1006
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
									goto l1006
								}
								position++
							}
						l1018:
							{
								position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l1021
								}
								position++
								goto l1020
							l1021:
								position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
								if buffer[position] != rune('K') {
									goto l1006
								}
								position++
							}
						l1020:
							if !rules[ruleskip]() {
								goto l1006
							}
							depth--
							add(ruleISBLANK, position1007)
						}
						goto l829
					l1006:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							position1023 := position
							depth++
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
									goto l1022
								}
								position++
							}
						l1024:
							{
								position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1027
								}
								position++
								goto l1026
							l1027:
								position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
								if buffer[position] != rune('S') {
									goto l1022
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
									goto l1022
								}
								position++
							}
						l1028:
							{
								position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1031
								}
								position++
								goto l1030
							l1031:
								position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
								if buffer[position] != rune('I') {
									goto l1022
								}
								position++
							}
						l1030:
							{
								position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1033
								}
								position++
								goto l1032
							l1033:
								position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
								if buffer[position] != rune('T') {
									goto l1022
								}
								position++
							}
						l1032:
							{
								position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1035
								}
								position++
								goto l1034
							l1035:
								position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
								if buffer[position] != rune('E') {
									goto l1022
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
									goto l1022
								}
								position++
							}
						l1036:
							{
								position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1039
								}
								position++
								goto l1038
							l1039:
								position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
								if buffer[position] != rune('A') {
									goto l1022
								}
								position++
							}
						l1038:
							{
								position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1041
								}
								position++
								goto l1040
							l1041:
								position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
								if buffer[position] != rune('L') {
									goto l1022
								}
								position++
							}
						l1040:
							if !rules[ruleskip]() {
								goto l1022
							}
							depth--
							add(ruleISLITERAL, position1023)
						}
						goto l829
					l1022:
						position, tokenIndex, depth = position829, tokenIndex829, depth829
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position1043 := position
									depth++
									{
										position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1045
										}
										position++
										goto l1044
									l1045:
										position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
										if buffer[position] != rune('I') {
											goto l828
										}
										position++
									}
								l1044:
									{
										position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1047
										}
										position++
										goto l1046
									l1047:
										position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
										if buffer[position] != rune('S') {
											goto l828
										}
										position++
									}
								l1046:
									{
										position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1049
										}
										position++
										goto l1048
									l1049:
										position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
										if buffer[position] != rune('N') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1050:
									{
										position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1053
										}
										position++
										goto l1052
									l1053:
										position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
										if buffer[position] != rune('M') {
											goto l828
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
											goto l828
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
											goto l828
										}
										position++
									}
								l1056:
									{
										position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1059
										}
										position++
										goto l1058
									l1059:
										position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
										if buffer[position] != rune('I') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1060:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleISNUMERIC, position1043)
								}
								break
							case 'S', 's':
								{
									position1062 := position
									depth++
									{
										position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1064
										}
										position++
										goto l1063
									l1064:
										position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
										if buffer[position] != rune('S') {
											goto l828
										}
										position++
									}
								l1063:
									{
										position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1066
										}
										position++
										goto l1065
									l1066:
										position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
										if buffer[position] != rune('H') {
											goto l828
										}
										position++
									}
								l1065:
									{
										position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1068
										}
										position++
										goto l1067
									l1068:
										position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
										if buffer[position] != rune('A') {
											goto l828
										}
										position++
									}
								l1067:
									if buffer[position] != rune('5') {
										goto l828
									}
									position++
									if buffer[position] != rune('1') {
										goto l828
									}
									position++
									if buffer[position] != rune('2') {
										goto l828
									}
									position++
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleSHA512, position1062)
								}
								break
							case 'M', 'm':
								{
									position1069 := position
									depth++
									{
										position1070, tokenIndex1070, depth1070 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1071
										}
										position++
										goto l1070
									l1071:
										position, tokenIndex, depth = position1070, tokenIndex1070, depth1070
										if buffer[position] != rune('M') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1072:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleMD5, position1069)
								}
								break
							case 'T', 't':
								{
									position1074 := position
									depth++
									{
										position1075, tokenIndex1075, depth1075 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1076
										}
										position++
										goto l1075
									l1076:
										position, tokenIndex, depth = position1075, tokenIndex1075, depth1075
										if buffer[position] != rune('T') {
											goto l828
										}
										position++
									}
								l1075:
									{
										position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1078
										}
										position++
										goto l1077
									l1078:
										position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
										if buffer[position] != rune('Z') {
											goto l828
										}
										position++
									}
								l1077:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleTZ, position1074)
								}
								break
							case 'H', 'h':
								{
									position1079 := position
									depth++
									{
										position1080, tokenIndex1080, depth1080 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1081
										}
										position++
										goto l1080
									l1081:
										position, tokenIndex, depth = position1080, tokenIndex1080, depth1080
										if buffer[position] != rune('H') {
											goto l828
										}
										position++
									}
								l1080:
									{
										position1082, tokenIndex1082, depth1082 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1083
										}
										position++
										goto l1082
									l1083:
										position, tokenIndex, depth = position1082, tokenIndex1082, depth1082
										if buffer[position] != rune('O') {
											goto l828
										}
										position++
									}
								l1082:
									{
										position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1085
										}
										position++
										goto l1084
									l1085:
										position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
										if buffer[position] != rune('U') {
											goto l828
										}
										position++
									}
								l1084:
									{
										position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1087
										}
										position++
										goto l1086
									l1087:
										position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
										if buffer[position] != rune('R') {
											goto l828
										}
										position++
									}
								l1086:
									{
										position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1089
										}
										position++
										goto l1088
									l1089:
										position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
										if buffer[position] != rune('S') {
											goto l828
										}
										position++
									}
								l1088:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleHOURS, position1079)
								}
								break
							case 'D', 'd':
								{
									position1090 := position
									depth++
									{
										position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1092
										}
										position++
										goto l1091
									l1092:
										position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
										if buffer[position] != rune('D') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1093:
									{
										position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1096
										}
										position++
										goto l1095
									l1096:
										position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
										if buffer[position] != rune('Y') {
											goto l828
										}
										position++
									}
								l1095:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleDAY, position1090)
								}
								break
							case 'Y', 'y':
								{
									position1097 := position
									depth++
									{
										position1098, tokenIndex1098, depth1098 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1099
										}
										position++
										goto l1098
									l1099:
										position, tokenIndex, depth = position1098, tokenIndex1098, depth1098
										if buffer[position] != rune('Y') {
											goto l828
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
											goto l828
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
											goto l828
										}
										position++
									}
								l1102:
									{
										position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1105
										}
										position++
										goto l1104
									l1105:
										position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
										if buffer[position] != rune('R') {
											goto l828
										}
										position++
									}
								l1104:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleYEAR, position1097)
								}
								break
							case 'E', 'e':
								{
									position1106 := position
									depth++
									{
										position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1108
										}
										position++
										goto l1107
									l1108:
										position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
										if buffer[position] != rune('E') {
											goto l828
										}
										position++
									}
								l1107:
									{
										position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1110
										}
										position++
										goto l1109
									l1110:
										position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
										if buffer[position] != rune('N') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1111:
									{
										position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1114
										}
										position++
										goto l1113
									l1114:
										position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
										if buffer[position] != rune('O') {
											goto l828
										}
										position++
									}
								l1113:
									{
										position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1116
										}
										position++
										goto l1115
									l1116:
										position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
										if buffer[position] != rune('D') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1117:
									if buffer[position] != rune('_') {
										goto l828
									}
									position++
									{
										position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1120
										}
										position++
										goto l1119
									l1120:
										position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
										if buffer[position] != rune('F') {
											goto l828
										}
										position++
									}
								l1119:
									{
										position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1122
										}
										position++
										goto l1121
									l1122:
										position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
										if buffer[position] != rune('O') {
											goto l828
										}
										position++
									}
								l1121:
									{
										position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1124
										}
										position++
										goto l1123
									l1124:
										position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
										if buffer[position] != rune('R') {
											goto l828
										}
										position++
									}
								l1123:
									if buffer[position] != rune('_') {
										goto l828
									}
									position++
									{
										position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1126
										}
										position++
										goto l1125
									l1126:
										position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
										if buffer[position] != rune('U') {
											goto l828
										}
										position++
									}
								l1125:
									{
										position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1128
										}
										position++
										goto l1127
									l1128:
										position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
										if buffer[position] != rune('R') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1129:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleENCODEFORURI, position1106)
								}
								break
							case 'L', 'l':
								{
									position1131 := position
									depth++
									{
										position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1133
										}
										position++
										goto l1132
									l1133:
										position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
										if buffer[position] != rune('L') {
											goto l828
										}
										position++
									}
								l1132:
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
											goto l828
										}
										position++
									}
								l1134:
									{
										position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1137
										}
										position++
										goto l1136
									l1137:
										position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
										if buffer[position] != rune('A') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1138:
									{
										position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1141
										}
										position++
										goto l1140
									l1141:
										position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
										if buffer[position] != rune('E') {
											goto l828
										}
										position++
									}
								l1140:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleLCASE, position1131)
								}
								break
							case 'U', 'u':
								{
									position1142 := position
									depth++
									{
										position1143, tokenIndex1143, depth1143 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1144
										}
										position++
										goto l1143
									l1144:
										position, tokenIndex, depth = position1143, tokenIndex1143, depth1143
										if buffer[position] != rune('U') {
											goto l828
										}
										position++
									}
								l1143:
									{
										position1145, tokenIndex1145, depth1145 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1146
										}
										position++
										goto l1145
									l1146:
										position, tokenIndex, depth = position1145, tokenIndex1145, depth1145
										if buffer[position] != rune('C') {
											goto l828
										}
										position++
									}
								l1145:
									{
										position1147, tokenIndex1147, depth1147 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1148
										}
										position++
										goto l1147
									l1148:
										position, tokenIndex, depth = position1147, tokenIndex1147, depth1147
										if buffer[position] != rune('A') {
											goto l828
										}
										position++
									}
								l1147:
									{
										position1149, tokenIndex1149, depth1149 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1150
										}
										position++
										goto l1149
									l1150:
										position, tokenIndex, depth = position1149, tokenIndex1149, depth1149
										if buffer[position] != rune('S') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1151:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleUCASE, position1142)
								}
								break
							case 'F', 'f':
								{
									position1153 := position
									depth++
									{
										position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1155
										}
										position++
										goto l1154
									l1155:
										position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
										if buffer[position] != rune('F') {
											goto l828
										}
										position++
									}
								l1154:
									{
										position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1157
										}
										position++
										goto l1156
									l1157:
										position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
										if buffer[position] != rune('L') {
											goto l828
										}
										position++
									}
								l1156:
									{
										position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1159
										}
										position++
										goto l1158
									l1159:
										position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
										if buffer[position] != rune('O') {
											goto l828
										}
										position++
									}
								l1158:
									{
										position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1161
										}
										position++
										goto l1160
									l1161:
										position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
										if buffer[position] != rune('O') {
											goto l828
										}
										position++
									}
								l1160:
									{
										position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1163
										}
										position++
										goto l1162
									l1163:
										position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
										if buffer[position] != rune('R') {
											goto l828
										}
										position++
									}
								l1162:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleFLOOR, position1153)
								}
								break
							case 'R', 'r':
								{
									position1164 := position
									depth++
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
											goto l828
										}
										position++
									}
								l1165:
									{
										position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1168
										}
										position++
										goto l1167
									l1168:
										position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
										if buffer[position] != rune('O') {
											goto l828
										}
										position++
									}
								l1167:
									{
										position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1170
										}
										position++
										goto l1169
									l1170:
										position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
										if buffer[position] != rune('U') {
											goto l828
										}
										position++
									}
								l1169:
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
											goto l828
										}
										position++
									}
								l1171:
									{
										position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1174
										}
										position++
										goto l1173
									l1174:
										position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
										if buffer[position] != rune('D') {
											goto l828
										}
										position++
									}
								l1173:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleROUND, position1164)
								}
								break
							case 'C', 'c':
								{
									position1175 := position
									depth++
									{
										position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1177
										}
										position++
										goto l1176
									l1177:
										position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
										if buffer[position] != rune('C') {
											goto l828
										}
										position++
									}
								l1176:
									{
										position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1179
										}
										position++
										goto l1178
									l1179:
										position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
										if buffer[position] != rune('E') {
											goto l828
										}
										position++
									}
								l1178:
									{
										position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1181
										}
										position++
										goto l1180
									l1181:
										position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
										if buffer[position] != rune('I') {
											goto l828
										}
										position++
									}
								l1180:
									{
										position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1183
										}
										position++
										goto l1182
									l1183:
										position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
										if buffer[position] != rune('L') {
											goto l828
										}
										position++
									}
								l1182:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleCEIL, position1175)
								}
								break
							default:
								{
									position1184 := position
									depth++
									{
										position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1186
										}
										position++
										goto l1185
									l1186:
										position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
										if buffer[position] != rune('A') {
											goto l828
										}
										position++
									}
								l1185:
									{
										position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1188
										}
										position++
										goto l1187
									l1188:
										position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
										if buffer[position] != rune('B') {
											goto l828
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
											goto l828
										}
										position++
									}
								l1189:
									if !rules[ruleskip]() {
										goto l828
									}
									depth--
									add(ruleABS, position1184)
								}
								break
							}
						}

					}
				l829:
					if !rules[ruleLPAREN]() {
						goto l828
					}
					if !rules[ruleexpression]() {
						goto l828
					}
					if !rules[ruleRPAREN]() {
						goto l828
					}
					goto l827
				l828:
					position, tokenIndex, depth = position827, tokenIndex827, depth827
					{
						position1192, tokenIndex1192, depth1192 := position, tokenIndex, depth
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
								if buffer[position] != rune('s') {
									goto l1202
								}
								position++
								goto l1201
							l1202:
								position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
								if buffer[position] != rune('S') {
									goto l1193
								}
								position++
							}
						l1201:
							{
								position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1204
								}
								position++
								goto l1203
							l1204:
								position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
								if buffer[position] != rune('T') {
									goto l1193
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
									goto l1193
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
									goto l1193
								}
								position++
							}
						l1207:
							{
								position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1210
								}
								position++
								goto l1209
							l1210:
								position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
								if buffer[position] != rune('T') {
									goto l1193
								}
								position++
							}
						l1209:
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
									goto l1193
								}
								position++
							}
						l1211:
							if !rules[ruleskip]() {
								goto l1193
							}
							depth--
							add(ruleSTRSTARTS, position1194)
						}
						goto l1192
					l1193:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						{
							position1214 := position
							depth++
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
									goto l1213
								}
								position++
							}
						l1215:
							{
								position1217, tokenIndex1217, depth1217 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1218
								}
								position++
								goto l1217
							l1218:
								position, tokenIndex, depth = position1217, tokenIndex1217, depth1217
								if buffer[position] != rune('T') {
									goto l1213
								}
								position++
							}
						l1217:
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
									goto l1213
								}
								position++
							}
						l1219:
							{
								position1221, tokenIndex1221, depth1221 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1222
								}
								position++
								goto l1221
							l1222:
								position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
								if buffer[position] != rune('E') {
									goto l1213
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
									goto l1213
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
									goto l1213
								}
								position++
							}
						l1225:
							{
								position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1228
								}
								position++
								goto l1227
							l1228:
								position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
								if buffer[position] != rune('S') {
									goto l1213
								}
								position++
							}
						l1227:
							if !rules[ruleskip]() {
								goto l1213
							}
							depth--
							add(ruleSTRENDS, position1214)
						}
						goto l1192
					l1213:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						{
							position1230 := position
							depth++
							{
								position1231, tokenIndex1231, depth1231 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1232
								}
								position++
								goto l1231
							l1232:
								position, tokenIndex, depth = position1231, tokenIndex1231, depth1231
								if buffer[position] != rune('S') {
									goto l1229
								}
								position++
							}
						l1231:
							{
								position1233, tokenIndex1233, depth1233 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1234
								}
								position++
								goto l1233
							l1234:
								position, tokenIndex, depth = position1233, tokenIndex1233, depth1233
								if buffer[position] != rune('T') {
									goto l1229
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
									goto l1229
								}
								position++
							}
						l1235:
							{
								position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1238
								}
								position++
								goto l1237
							l1238:
								position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
								if buffer[position] != rune('B') {
									goto l1229
								}
								position++
							}
						l1237:
							{
								position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1240
								}
								position++
								goto l1239
							l1240:
								position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
								if buffer[position] != rune('E') {
									goto l1229
								}
								position++
							}
						l1239:
							{
								position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1242
								}
								position++
								goto l1241
							l1242:
								position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
								if buffer[position] != rune('F') {
									goto l1229
								}
								position++
							}
						l1241:
							{
								position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1244
								}
								position++
								goto l1243
							l1244:
								position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
								if buffer[position] != rune('O') {
									goto l1229
								}
								position++
							}
						l1243:
							{
								position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1246
								}
								position++
								goto l1245
							l1246:
								position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
								if buffer[position] != rune('R') {
									goto l1229
								}
								position++
							}
						l1245:
							{
								position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1248
								}
								position++
								goto l1247
							l1248:
								position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
								if buffer[position] != rune('E') {
									goto l1229
								}
								position++
							}
						l1247:
							if !rules[ruleskip]() {
								goto l1229
							}
							depth--
							add(ruleSTRBEFORE, position1230)
						}
						goto l1192
					l1229:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
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
								if buffer[position] != rune('a') {
									goto l1258
								}
								position++
								goto l1257
							l1258:
								position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
								if buffer[position] != rune('A') {
									goto l1249
								}
								position++
							}
						l1257:
							{
								position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1260
								}
								position++
								goto l1259
							l1260:
								position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
								if buffer[position] != rune('F') {
									goto l1249
								}
								position++
							}
						l1259:
							{
								position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1262
								}
								position++
								goto l1261
							l1262:
								position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
								if buffer[position] != rune('T') {
									goto l1249
								}
								position++
							}
						l1261:
							{
								position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1264
								}
								position++
								goto l1263
							l1264:
								position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
								if buffer[position] != rune('E') {
									goto l1249
								}
								position++
							}
						l1263:
							{
								position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1266
								}
								position++
								goto l1265
							l1266:
								position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
								if buffer[position] != rune('R') {
									goto l1249
								}
								position++
							}
						l1265:
							if !rules[ruleskip]() {
								goto l1249
							}
							depth--
							add(ruleSTRAFTER, position1250)
						}
						goto l1192
					l1249:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						{
							position1268 := position
							depth++
							{
								position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1270
								}
								position++
								goto l1269
							l1270:
								position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
								if buffer[position] != rune('S') {
									goto l1267
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
									goto l1267
								}
								position++
							}
						l1271:
							{
								position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1274
								}
								position++
								goto l1273
							l1274:
								position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
								if buffer[position] != rune('R') {
									goto l1267
								}
								position++
							}
						l1273:
							{
								position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1276
								}
								position++
								goto l1275
							l1276:
								position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
								if buffer[position] != rune('L') {
									goto l1267
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
									goto l1267
								}
								position++
							}
						l1277:
							{
								position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1280
								}
								position++
								goto l1279
							l1280:
								position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
								if buffer[position] != rune('N') {
									goto l1267
								}
								position++
							}
						l1279:
							{
								position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1282
								}
								position++
								goto l1281
							l1282:
								position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
								if buffer[position] != rune('G') {
									goto l1267
								}
								position++
							}
						l1281:
							if !rules[ruleskip]() {
								goto l1267
							}
							depth--
							add(ruleSTRLANG, position1268)
						}
						goto l1192
					l1267:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
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
									goto l1283
								}
								position++
							}
						l1285:
							{
								position1287, tokenIndex1287, depth1287 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1288
								}
								position++
								goto l1287
							l1288:
								position, tokenIndex, depth = position1287, tokenIndex1287, depth1287
								if buffer[position] != rune('T') {
									goto l1283
								}
								position++
							}
						l1287:
							{
								position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1290
								}
								position++
								goto l1289
							l1290:
								position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
								if buffer[position] != rune('R') {
									goto l1283
								}
								position++
							}
						l1289:
							{
								position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1292
								}
								position++
								goto l1291
							l1292:
								position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
								if buffer[position] != rune('D') {
									goto l1283
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
									goto l1283
								}
								position++
							}
						l1293:
							if !rules[ruleskip]() {
								goto l1283
							}
							depth--
							add(ruleSTRDT, position1284)
						}
						goto l1192
					l1283:
						position, tokenIndex, depth = position1192, tokenIndex1192, depth1192
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1296 := position
									depth++
									{
										position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1298
										}
										position++
										goto l1297
									l1298:
										position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
										if buffer[position] != rune('S') {
											goto l1191
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
											goto l1191
										}
										position++
									}
								l1299:
									{
										position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1302
										}
										position++
										goto l1301
									l1302:
										position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
										if buffer[position] != rune('M') {
											goto l1191
										}
										position++
									}
								l1301:
									{
										position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1304
										}
										position++
										goto l1303
									l1304:
										position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
										if buffer[position] != rune('E') {
											goto l1191
										}
										position++
									}
								l1303:
									{
										position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1306
										}
										position++
										goto l1305
									l1306:
										position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
										if buffer[position] != rune('T') {
											goto l1191
										}
										position++
									}
								l1305:
									{
										position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1308
										}
										position++
										goto l1307
									l1308:
										position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
										if buffer[position] != rune('E') {
											goto l1191
										}
										position++
									}
								l1307:
									{
										position1309, tokenIndex1309, depth1309 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1310
										}
										position++
										goto l1309
									l1310:
										position, tokenIndex, depth = position1309, tokenIndex1309, depth1309
										if buffer[position] != rune('R') {
											goto l1191
										}
										position++
									}
								l1309:
									{
										position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1312
										}
										position++
										goto l1311
									l1312:
										position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
										if buffer[position] != rune('M') {
											goto l1191
										}
										position++
									}
								l1311:
									if !rules[ruleskip]() {
										goto l1191
									}
									depth--
									add(ruleSAMETERM, position1296)
								}
								break
							case 'C', 'c':
								{
									position1313 := position
									depth++
									{
										position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1315
										}
										position++
										goto l1314
									l1315:
										position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
										if buffer[position] != rune('C') {
											goto l1191
										}
										position++
									}
								l1314:
									{
										position1316, tokenIndex1316, depth1316 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1317
										}
										position++
										goto l1316
									l1317:
										position, tokenIndex, depth = position1316, tokenIndex1316, depth1316
										if buffer[position] != rune('O') {
											goto l1191
										}
										position++
									}
								l1316:
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
											goto l1191
										}
										position++
									}
								l1318:
									{
										position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1321
										}
										position++
										goto l1320
									l1321:
										position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
										if buffer[position] != rune('T') {
											goto l1191
										}
										position++
									}
								l1320:
									{
										position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1323
										}
										position++
										goto l1322
									l1323:
										position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
										if buffer[position] != rune('A') {
											goto l1191
										}
										position++
									}
								l1322:
									{
										position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1325
										}
										position++
										goto l1324
									l1325:
										position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
										if buffer[position] != rune('I') {
											goto l1191
										}
										position++
									}
								l1324:
									{
										position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1327
										}
										position++
										goto l1326
									l1327:
										position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
										if buffer[position] != rune('N') {
											goto l1191
										}
										position++
									}
								l1326:
									{
										position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1329
										}
										position++
										goto l1328
									l1329:
										position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
										if buffer[position] != rune('S') {
											goto l1191
										}
										position++
									}
								l1328:
									if !rules[ruleskip]() {
										goto l1191
									}
									depth--
									add(ruleCONTAINS, position1313)
								}
								break
							default:
								{
									position1330 := position
									depth++
									{
										position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1332
										}
										position++
										goto l1331
									l1332:
										position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
										if buffer[position] != rune('L') {
											goto l1191
										}
										position++
									}
								l1331:
									{
										position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1334
										}
										position++
										goto l1333
									l1334:
										position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
										if buffer[position] != rune('A') {
											goto l1191
										}
										position++
									}
								l1333:
									{
										position1335, tokenIndex1335, depth1335 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1336
										}
										position++
										goto l1335
									l1336:
										position, tokenIndex, depth = position1335, tokenIndex1335, depth1335
										if buffer[position] != rune('N') {
											goto l1191
										}
										position++
									}
								l1335:
									{
										position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1338
										}
										position++
										goto l1337
									l1338:
										position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
										if buffer[position] != rune('G') {
											goto l1191
										}
										position++
									}
								l1337:
									{
										position1339, tokenIndex1339, depth1339 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1340
										}
										position++
										goto l1339
									l1340:
										position, tokenIndex, depth = position1339, tokenIndex1339, depth1339
										if buffer[position] != rune('M') {
											goto l1191
										}
										position++
									}
								l1339:
									{
										position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1342
										}
										position++
										goto l1341
									l1342:
										position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
										if buffer[position] != rune('A') {
											goto l1191
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
											goto l1191
										}
										position++
									}
								l1343:
									{
										position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1346
										}
										position++
										goto l1345
									l1346:
										position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
										if buffer[position] != rune('C') {
											goto l1191
										}
										position++
									}
								l1345:
									{
										position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1348
										}
										position++
										goto l1347
									l1348:
										position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
										if buffer[position] != rune('H') {
											goto l1191
										}
										position++
									}
								l1347:
									{
										position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1350
										}
										position++
										goto l1349
									l1350:
										position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
										if buffer[position] != rune('E') {
											goto l1191
										}
										position++
									}
								l1349:
									{
										position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1352
										}
										position++
										goto l1351
									l1352:
										position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
										if buffer[position] != rune('S') {
											goto l1191
										}
										position++
									}
								l1351:
									if !rules[ruleskip]() {
										goto l1191
									}
									depth--
									add(ruleLANGMATCHES, position1330)
								}
								break
							}
						}

					}
				l1192:
					if !rules[ruleLPAREN]() {
						goto l1191
					}
					if !rules[ruleexpression]() {
						goto l1191
					}
					if !rules[ruleCOMMA]() {
						goto l1191
					}
					if !rules[ruleexpression]() {
						goto l1191
					}
					if !rules[ruleRPAREN]() {
						goto l1191
					}
					goto l827
				l1191:
					position, tokenIndex, depth = position827, tokenIndex827, depth827
					{
						position1354 := position
						depth++
						{
							position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1356
							}
							position++
							goto l1355
						l1356:
							position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
							if buffer[position] != rune('B') {
								goto l1353
							}
							position++
						}
					l1355:
						{
							position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1358
							}
							position++
							goto l1357
						l1358:
							position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
							if buffer[position] != rune('O') {
								goto l1353
							}
							position++
						}
					l1357:
						{
							position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1360
							}
							position++
							goto l1359
						l1360:
							position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
							if buffer[position] != rune('U') {
								goto l1353
							}
							position++
						}
					l1359:
						{
							position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1362
							}
							position++
							goto l1361
						l1362:
							position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
							if buffer[position] != rune('N') {
								goto l1353
							}
							position++
						}
					l1361:
						{
							position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1364
							}
							position++
							goto l1363
						l1364:
							position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
							if buffer[position] != rune('D') {
								goto l1353
							}
							position++
						}
					l1363:
						if !rules[ruleskip]() {
							goto l1353
						}
						depth--
						add(ruleBOUND, position1354)
					}
					if !rules[ruleLPAREN]() {
						goto l1353
					}
					if !rules[rulevar]() {
						goto l1353
					}
					if !rules[ruleRPAREN]() {
						goto l1353
					}
					goto l827
				l1353:
					position, tokenIndex, depth = position827, tokenIndex827, depth827
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1367 := position
								depth++
								{
									position1368, tokenIndex1368, depth1368 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1369
									}
									position++
									goto l1368
								l1369:
									position, tokenIndex, depth = position1368, tokenIndex1368, depth1368
									if buffer[position] != rune('S') {
										goto l1365
									}
									position++
								}
							l1368:
								{
									position1370, tokenIndex1370, depth1370 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1371
									}
									position++
									goto l1370
								l1371:
									position, tokenIndex, depth = position1370, tokenIndex1370, depth1370
									if buffer[position] != rune('T') {
										goto l1365
									}
									position++
								}
							l1370:
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
										goto l1365
									}
									position++
								}
							l1372:
								{
									position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1375
									}
									position++
									goto l1374
								l1375:
									position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
									if buffer[position] != rune('U') {
										goto l1365
									}
									position++
								}
							l1374:
								{
									position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1377
									}
									position++
									goto l1376
								l1377:
									position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
									if buffer[position] != rune('U') {
										goto l1365
									}
									position++
								}
							l1376:
								{
									position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1379
									}
									position++
									goto l1378
								l1379:
									position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
									if buffer[position] != rune('I') {
										goto l1365
									}
									position++
								}
							l1378:
								{
									position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1381
									}
									position++
									goto l1380
								l1381:
									position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
									if buffer[position] != rune('D') {
										goto l1365
									}
									position++
								}
							l1380:
								if !rules[ruleskip]() {
									goto l1365
								}
								depth--
								add(ruleSTRUUID, position1367)
							}
							break
						case 'U', 'u':
							{
								position1382 := position
								depth++
								{
									position1383, tokenIndex1383, depth1383 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1384
									}
									position++
									goto l1383
								l1384:
									position, tokenIndex, depth = position1383, tokenIndex1383, depth1383
									if buffer[position] != rune('U') {
										goto l1365
									}
									position++
								}
							l1383:
								{
									position1385, tokenIndex1385, depth1385 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1386
									}
									position++
									goto l1385
								l1386:
									position, tokenIndex, depth = position1385, tokenIndex1385, depth1385
									if buffer[position] != rune('U') {
										goto l1365
									}
									position++
								}
							l1385:
								{
									position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1388
									}
									position++
									goto l1387
								l1388:
									position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
									if buffer[position] != rune('I') {
										goto l1365
									}
									position++
								}
							l1387:
								{
									position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1390
									}
									position++
									goto l1389
								l1390:
									position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
									if buffer[position] != rune('D') {
										goto l1365
									}
									position++
								}
							l1389:
								if !rules[ruleskip]() {
									goto l1365
								}
								depth--
								add(ruleUUID, position1382)
							}
							break
						case 'N', 'n':
							{
								position1391 := position
								depth++
								{
									position1392, tokenIndex1392, depth1392 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1393
									}
									position++
									goto l1392
								l1393:
									position, tokenIndex, depth = position1392, tokenIndex1392, depth1392
									if buffer[position] != rune('N') {
										goto l1365
									}
									position++
								}
							l1392:
								{
									position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1395
									}
									position++
									goto l1394
								l1395:
									position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
									if buffer[position] != rune('O') {
										goto l1365
									}
									position++
								}
							l1394:
								{
									position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1397
									}
									position++
									goto l1396
								l1397:
									position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
									if buffer[position] != rune('W') {
										goto l1365
									}
									position++
								}
							l1396:
								if !rules[ruleskip]() {
									goto l1365
								}
								depth--
								add(ruleNOW, position1391)
							}
							break
						default:
							{
								position1398 := position
								depth++
								{
									position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1400
									}
									position++
									goto l1399
								l1400:
									position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
									if buffer[position] != rune('R') {
										goto l1365
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
										goto l1365
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
										goto l1365
									}
									position++
								}
							l1403:
								{
									position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1406
									}
									position++
									goto l1405
								l1406:
									position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
									if buffer[position] != rune('D') {
										goto l1365
									}
									position++
								}
							l1405:
								if !rules[ruleskip]() {
									goto l1365
								}
								depth--
								add(ruleRAND, position1398)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1365
					}
					goto l827
				l1365:
					position, tokenIndex, depth = position827, tokenIndex827, depth827
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
								{
									position1410 := position
									depth++
									{
										position1411, tokenIndex1411, depth1411 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1412
										}
										position++
										goto l1411
									l1412:
										position, tokenIndex, depth = position1411, tokenIndex1411, depth1411
										if buffer[position] != rune('E') {
											goto l1409
										}
										position++
									}
								l1411:
									{
										position1413, tokenIndex1413, depth1413 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1414
										}
										position++
										goto l1413
									l1414:
										position, tokenIndex, depth = position1413, tokenIndex1413, depth1413
										if buffer[position] != rune('X') {
											goto l1409
										}
										position++
									}
								l1413:
									{
										position1415, tokenIndex1415, depth1415 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1416
										}
										position++
										goto l1415
									l1416:
										position, tokenIndex, depth = position1415, tokenIndex1415, depth1415
										if buffer[position] != rune('I') {
											goto l1409
										}
										position++
									}
								l1415:
									{
										position1417, tokenIndex1417, depth1417 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1418
										}
										position++
										goto l1417
									l1418:
										position, tokenIndex, depth = position1417, tokenIndex1417, depth1417
										if buffer[position] != rune('S') {
											goto l1409
										}
										position++
									}
								l1417:
									{
										position1419, tokenIndex1419, depth1419 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1420
										}
										position++
										goto l1419
									l1420:
										position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
										if buffer[position] != rune('T') {
											goto l1409
										}
										position++
									}
								l1419:
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
											goto l1409
										}
										position++
									}
								l1421:
									if !rules[ruleskip]() {
										goto l1409
									}
									depth--
									add(ruleEXISTS, position1410)
								}
								goto l1408
							l1409:
								position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
								{
									position1423 := position
									depth++
									{
										position1424, tokenIndex1424, depth1424 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1425
										}
										position++
										goto l1424
									l1425:
										position, tokenIndex, depth = position1424, tokenIndex1424, depth1424
										if buffer[position] != rune('N') {
											goto l825
										}
										position++
									}
								l1424:
									{
										position1426, tokenIndex1426, depth1426 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1427
										}
										position++
										goto l1426
									l1427:
										position, tokenIndex, depth = position1426, tokenIndex1426, depth1426
										if buffer[position] != rune('O') {
											goto l825
										}
										position++
									}
								l1426:
									{
										position1428, tokenIndex1428, depth1428 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1429
										}
										position++
										goto l1428
									l1429:
										position, tokenIndex, depth = position1428, tokenIndex1428, depth1428
										if buffer[position] != rune('T') {
											goto l825
										}
										position++
									}
								l1428:
									if buffer[position] != rune(' ') {
										goto l825
									}
									position++
									{
										position1430, tokenIndex1430, depth1430 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1431
										}
										position++
										goto l1430
									l1431:
										position, tokenIndex, depth = position1430, tokenIndex1430, depth1430
										if buffer[position] != rune('E') {
											goto l825
										}
										position++
									}
								l1430:
									{
										position1432, tokenIndex1432, depth1432 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1433
										}
										position++
										goto l1432
									l1433:
										position, tokenIndex, depth = position1432, tokenIndex1432, depth1432
										if buffer[position] != rune('X') {
											goto l825
										}
										position++
									}
								l1432:
									{
										position1434, tokenIndex1434, depth1434 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1435
										}
										position++
										goto l1434
									l1435:
										position, tokenIndex, depth = position1434, tokenIndex1434, depth1434
										if buffer[position] != rune('I') {
											goto l825
										}
										position++
									}
								l1434:
									{
										position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1437
										}
										position++
										goto l1436
									l1437:
										position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
										if buffer[position] != rune('S') {
											goto l825
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
											goto l825
										}
										position++
									}
								l1438:
									{
										position1440, tokenIndex1440, depth1440 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1441
										}
										position++
										goto l1440
									l1441:
										position, tokenIndex, depth = position1440, tokenIndex1440, depth1440
										if buffer[position] != rune('S') {
											goto l825
										}
										position++
									}
								l1440:
									if !rules[ruleskip]() {
										goto l825
									}
									depth--
									add(ruleNOTEXIST, position1423)
								}
							}
						l1408:
							if !rules[rulegroupGraphPattern]() {
								goto l825
							}
							break
						case 'I', 'i':
							{
								position1442 := position
								depth++
								{
									position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1444
									}
									position++
									goto l1443
								l1444:
									position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
									if buffer[position] != rune('I') {
										goto l825
									}
									position++
								}
							l1443:
								{
									position1445, tokenIndex1445, depth1445 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1446
									}
									position++
									goto l1445
								l1446:
									position, tokenIndex, depth = position1445, tokenIndex1445, depth1445
									if buffer[position] != rune('F') {
										goto l825
									}
									position++
								}
							l1445:
								if !rules[ruleskip]() {
									goto l825
								}
								depth--
								add(ruleIF, position1442)
							}
							if !rules[ruleLPAREN]() {
								goto l825
							}
							if !rules[ruleexpression]() {
								goto l825
							}
							if !rules[ruleCOMMA]() {
								goto l825
							}
							if !rules[ruleexpression]() {
								goto l825
							}
							if !rules[ruleCOMMA]() {
								goto l825
							}
							if !rules[ruleexpression]() {
								goto l825
							}
							if !rules[ruleRPAREN]() {
								goto l825
							}
							break
						case 'C', 'c':
							{
								position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
								{
									position1449 := position
									depth++
									{
										position1450, tokenIndex1450, depth1450 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1451
										}
										position++
										goto l1450
									l1451:
										position, tokenIndex, depth = position1450, tokenIndex1450, depth1450
										if buffer[position] != rune('C') {
											goto l1448
										}
										position++
									}
								l1450:
									{
										position1452, tokenIndex1452, depth1452 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1453
										}
										position++
										goto l1452
									l1453:
										position, tokenIndex, depth = position1452, tokenIndex1452, depth1452
										if buffer[position] != rune('O') {
											goto l1448
										}
										position++
									}
								l1452:
									{
										position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1455
										}
										position++
										goto l1454
									l1455:
										position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
										if buffer[position] != rune('N') {
											goto l1448
										}
										position++
									}
								l1454:
									{
										position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1457
										}
										position++
										goto l1456
									l1457:
										position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
										if buffer[position] != rune('C') {
											goto l1448
										}
										position++
									}
								l1456:
									{
										position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1459
										}
										position++
										goto l1458
									l1459:
										position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
										if buffer[position] != rune('A') {
											goto l1448
										}
										position++
									}
								l1458:
									{
										position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1461
										}
										position++
										goto l1460
									l1461:
										position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
										if buffer[position] != rune('T') {
											goto l1448
										}
										position++
									}
								l1460:
									if !rules[ruleskip]() {
										goto l1448
									}
									depth--
									add(ruleCONCAT, position1449)
								}
								goto l1447
							l1448:
								position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
								{
									position1462 := position
									depth++
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
											goto l825
										}
										position++
									}
								l1463:
									{
										position1465, tokenIndex1465, depth1465 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1466
										}
										position++
										goto l1465
									l1466:
										position, tokenIndex, depth = position1465, tokenIndex1465, depth1465
										if buffer[position] != rune('O') {
											goto l825
										}
										position++
									}
								l1465:
									{
										position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1468
										}
										position++
										goto l1467
									l1468:
										position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
										if buffer[position] != rune('A') {
											goto l825
										}
										position++
									}
								l1467:
									{
										position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1470
										}
										position++
										goto l1469
									l1470:
										position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
										if buffer[position] != rune('L') {
											goto l825
										}
										position++
									}
								l1469:
									{
										position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1472
										}
										position++
										goto l1471
									l1472:
										position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
										if buffer[position] != rune('E') {
											goto l825
										}
										position++
									}
								l1471:
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
											goto l825
										}
										position++
									}
								l1473:
									{
										position1475, tokenIndex1475, depth1475 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1476
										}
										position++
										goto l1475
									l1476:
										position, tokenIndex, depth = position1475, tokenIndex1475, depth1475
										if buffer[position] != rune('C') {
											goto l825
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
											goto l825
										}
										position++
									}
								l1477:
									if !rules[ruleskip]() {
										goto l825
									}
									depth--
									add(ruleCOALESCE, position1462)
								}
							}
						l1447:
							if !rules[ruleargList]() {
								goto l825
							}
							break
						case 'B', 'b':
							{
								position1479 := position
								depth++
								{
									position1480, tokenIndex1480, depth1480 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1481
									}
									position++
									goto l1480
								l1481:
									position, tokenIndex, depth = position1480, tokenIndex1480, depth1480
									if buffer[position] != rune('B') {
										goto l825
									}
									position++
								}
							l1480:
								{
									position1482, tokenIndex1482, depth1482 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1483
									}
									position++
									goto l1482
								l1483:
									position, tokenIndex, depth = position1482, tokenIndex1482, depth1482
									if buffer[position] != rune('N') {
										goto l825
									}
									position++
								}
							l1482:
								{
									position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1485
									}
									position++
									goto l1484
								l1485:
									position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
									if buffer[position] != rune('O') {
										goto l825
									}
									position++
								}
							l1484:
								{
									position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1487
									}
									position++
									goto l1486
								l1487:
									position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
									if buffer[position] != rune('D') {
										goto l825
									}
									position++
								}
							l1486:
								{
									position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1489
									}
									position++
									goto l1488
								l1489:
									position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
									if buffer[position] != rune('E') {
										goto l825
									}
									position++
								}
							l1488:
								if !rules[ruleskip]() {
									goto l825
								}
								depth--
								add(ruleBNODE, position1479)
							}
							{
								position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1491
								}
								if !rules[ruleexpression]() {
									goto l1491
								}
								if !rules[ruleRPAREN]() {
									goto l1491
								}
								goto l1490
							l1491:
								position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
								if !rules[rulenil]() {
									goto l825
								}
							}
						l1490:
							break
						default:
							{
								position1492, tokenIndex1492, depth1492 := position, tokenIndex, depth
								{
									position1494 := position
									depth++
									{
										position1495, tokenIndex1495, depth1495 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1496
										}
										position++
										goto l1495
									l1496:
										position, tokenIndex, depth = position1495, tokenIndex1495, depth1495
										if buffer[position] != rune('S') {
											goto l1493
										}
										position++
									}
								l1495:
									{
										position1497, tokenIndex1497, depth1497 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1498
										}
										position++
										goto l1497
									l1498:
										position, tokenIndex, depth = position1497, tokenIndex1497, depth1497
										if buffer[position] != rune('U') {
											goto l1493
										}
										position++
									}
								l1497:
									{
										position1499, tokenIndex1499, depth1499 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1500
										}
										position++
										goto l1499
									l1500:
										position, tokenIndex, depth = position1499, tokenIndex1499, depth1499
										if buffer[position] != rune('B') {
											goto l1493
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
											goto l1493
										}
										position++
									}
								l1501:
									{
										position1503, tokenIndex1503, depth1503 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1504
										}
										position++
										goto l1503
									l1504:
										position, tokenIndex, depth = position1503, tokenIndex1503, depth1503
										if buffer[position] != rune('T') {
											goto l1493
										}
										position++
									}
								l1503:
									{
										position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1506
										}
										position++
										goto l1505
									l1506:
										position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
										if buffer[position] != rune('R') {
											goto l1493
										}
										position++
									}
								l1505:
									if !rules[ruleskip]() {
										goto l1493
									}
									depth--
									add(ruleSUBSTR, position1494)
								}
								goto l1492
							l1493:
								position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
								{
									position1508 := position
									depth++
									{
										position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1510
										}
										position++
										goto l1509
									l1510:
										position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
										if buffer[position] != rune('R') {
											goto l1507
										}
										position++
									}
								l1509:
									{
										position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1512
										}
										position++
										goto l1511
									l1512:
										position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
										if buffer[position] != rune('E') {
											goto l1507
										}
										position++
									}
								l1511:
									{
										position1513, tokenIndex1513, depth1513 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1514
										}
										position++
										goto l1513
									l1514:
										position, tokenIndex, depth = position1513, tokenIndex1513, depth1513
										if buffer[position] != rune('P') {
											goto l1507
										}
										position++
									}
								l1513:
									{
										position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1516
										}
										position++
										goto l1515
									l1516:
										position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
										if buffer[position] != rune('L') {
											goto l1507
										}
										position++
									}
								l1515:
									{
										position1517, tokenIndex1517, depth1517 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1518
										}
										position++
										goto l1517
									l1518:
										position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
										if buffer[position] != rune('A') {
											goto l1507
										}
										position++
									}
								l1517:
									{
										position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1520
										}
										position++
										goto l1519
									l1520:
										position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
										if buffer[position] != rune('C') {
											goto l1507
										}
										position++
									}
								l1519:
									{
										position1521, tokenIndex1521, depth1521 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1522
										}
										position++
										goto l1521
									l1522:
										position, tokenIndex, depth = position1521, tokenIndex1521, depth1521
										if buffer[position] != rune('E') {
											goto l1507
										}
										position++
									}
								l1521:
									if !rules[ruleskip]() {
										goto l1507
									}
									depth--
									add(ruleREPLACE, position1508)
								}
								goto l1492
							l1507:
								position, tokenIndex, depth = position1492, tokenIndex1492, depth1492
								{
									position1523 := position
									depth++
									{
										position1524, tokenIndex1524, depth1524 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1525
										}
										position++
										goto l1524
									l1525:
										position, tokenIndex, depth = position1524, tokenIndex1524, depth1524
										if buffer[position] != rune('R') {
											goto l825
										}
										position++
									}
								l1524:
									{
										position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1527
										}
										position++
										goto l1526
									l1527:
										position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
										if buffer[position] != rune('E') {
											goto l825
										}
										position++
									}
								l1526:
									{
										position1528, tokenIndex1528, depth1528 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1529
										}
										position++
										goto l1528
									l1529:
										position, tokenIndex, depth = position1528, tokenIndex1528, depth1528
										if buffer[position] != rune('G') {
											goto l825
										}
										position++
									}
								l1528:
									{
										position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1531
										}
										position++
										goto l1530
									l1531:
										position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
										if buffer[position] != rune('E') {
											goto l825
										}
										position++
									}
								l1530:
									{
										position1532, tokenIndex1532, depth1532 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1533
										}
										position++
										goto l1532
									l1533:
										position, tokenIndex, depth = position1532, tokenIndex1532, depth1532
										if buffer[position] != rune('X') {
											goto l825
										}
										position++
									}
								l1532:
									if !rules[ruleskip]() {
										goto l825
									}
									depth--
									add(ruleREGEX, position1523)
								}
							}
						l1492:
							if !rules[ruleLPAREN]() {
								goto l825
							}
							if !rules[ruleexpression]() {
								goto l825
							}
							if !rules[ruleCOMMA]() {
								goto l825
							}
							if !rules[ruleexpression]() {
								goto l825
							}
							{
								position1534, tokenIndex1534, depth1534 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1534
								}
								if !rules[ruleexpression]() {
									goto l1534
								}
								goto l1535
							l1534:
								position, tokenIndex, depth = position1534, tokenIndex1534, depth1534
							}
						l1535:
							if !rules[ruleRPAREN]() {
								goto l825
							}
							break
						}
					}

				}
			l827:
				depth--
				add(rulebuiltinCall, position826)
			}
			return true
		l825:
			position, tokenIndex, depth = position825, tokenIndex825, depth825
			return false
		},
		/* 70 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1536, tokenIndex1536, depth1536 := position, tokenIndex, depth
			{
				position1537 := position
				depth++
				{
					position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
					{
						position1540 := position
						depth++
					l1541:
						{
							position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
							{
								position1543, tokenIndex1543, depth1543 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1544
								}
								position++
								goto l1543
							l1544:
								position, tokenIndex, depth = position1543, tokenIndex1543, depth1543
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1542
								}
								position++
							}
						l1543:
							goto l1541
						l1542:
							position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
						}
						depth--
						add(rulePegText, position1540)
					}
					if buffer[position] != rune(':') {
						goto l1539
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1538
				l1539:
					position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
					{
						position1547 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1546
						}
						position++
					l1548:
						{
							position1549, tokenIndex1549, depth1549 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1549
							}
							position++
							goto l1548
						l1549:
							position, tokenIndex, depth = position1549, tokenIndex1549, depth1549
						}
						depth--
						add(rulePegText, position1547)
					}
					if buffer[position] != rune('/') {
						goto l1546
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1538
				l1546:
					position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
					{
						position1551 := position
						depth++
					l1552:
						{
							position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1553
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1553
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1553
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1553
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1553
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1553
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1553
									}
									position++
									break
								}
							}

							goto l1552
						l1553:
							position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
						}
						depth--
						add(rulePegText, position1551)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1538:
				if buffer[position] != rune('<') {
					goto l1536
				}
				position++
				if !rules[rulews]() {
					goto l1536
				}
				if !rules[ruleskip]() {
					goto l1536
				}
				depth--
				add(rulepof, position1537)
			}
			return true
		l1536:
			position, tokenIndex, depth = position1536, tokenIndex1536, depth1536
			return false
		},
		/* 71 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
			{
				position1557 := position
				depth++
				{
					position1558, tokenIndex1558, depth1558 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1559
					}
					position++
					goto l1558
				l1559:
					position, tokenIndex, depth = position1558, tokenIndex1558, depth1558
					if buffer[position] != rune('$') {
						goto l1556
					}
					position++
				}
			l1558:
				{
					position1560 := position
					depth++
					{
						position1561, tokenIndex1561, depth1561 := position, tokenIndex, depth
						if !rules[rulepnCharsU]() {
							goto l1562
						}
						goto l1561
					l1562:
						position, tokenIndex, depth = position1561, tokenIndex1561, depth1561
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1556
						}
						position++
					}
				l1561:
				l1563:
					{
						position1564, tokenIndex1564, depth1564 := position, tokenIndex, depth
						{
							position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
							if !rules[rulepnCharsU]() {
								goto l1566
							}
							goto l1565
						l1566:
							position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
							{
								switch buffer[position] {
								case 'â':
									if c := buffer[position]; c < rune('‿') || c > rune('⁀') {
										goto l1564
									}
									position++
									break
								case 'Ì', 'Í':
									if c := buffer[position]; c < rune('̀') || c > rune('ͯ') {
										goto l1564
									}
									position++
									break
								case 'Â':
									if buffer[position] != rune('·') {
										goto l1564
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1564
									}
									position++
									break
								}
							}

						}
					l1565:
						goto l1563
					l1564:
						position, tokenIndex, depth = position1564, tokenIndex1564, depth1564
					}
					depth--
					add(ruleVARNAME, position1560)
				}
				if !rules[ruleskip]() {
					goto l1556
				}
				depth--
				add(rulevar, position1557)
			}
			return true
		l1556:
			position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
			return false
		},
		/* 72 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
			{
				position1569 := position
				depth++
				{
					position1570, tokenIndex1570, depth1570 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1571
					}
					goto l1570
				l1571:
					position, tokenIndex, depth = position1570, tokenIndex1570, depth1570
					{
						position1572 := position
						depth++
						{
							position1573, tokenIndex1573, depth1573 := position, tokenIndex, depth
							if !rules[rulepnPrefix]() {
								goto l1573
							}
							goto l1574
						l1573:
							position, tokenIndex, depth = position1573, tokenIndex1573, depth1573
						}
					l1574:
						if buffer[position] != rune(':') {
							goto l1568
						}
						position++
						{
							position1575 := position
							depth++
							{
								switch buffer[position] {
								case '%', '\\':
									{
										position1579 := position
										depth++
										{
											position1580, tokenIndex1580, depth1580 := position, tokenIndex, depth
											{
												position1582 := position
												depth++
												if buffer[position] != rune('%') {
													goto l1581
												}
												position++
												if !rules[rulehex]() {
													goto l1581
												}
												if !rules[rulehex]() {
													goto l1581
												}
												depth--
												add(rulepercent, position1582)
											}
											goto l1580
										l1581:
											position, tokenIndex, depth = position1580, tokenIndex1580, depth1580
											{
												position1583 := position
												depth++
												if buffer[position] != rune('\\') {
													goto l1568
												}
												position++
												{
													switch buffer[position] {
													case '%':
														if buffer[position] != rune('%') {
															goto l1568
														}
														position++
														break
													case '@':
														if buffer[position] != rune('@') {
															goto l1568
														}
														position++
														break
													case '#':
														if buffer[position] != rune('#') {
															goto l1568
														}
														position++
														break
													case '?':
														if buffer[position] != rune('?') {
															goto l1568
														}
														position++
														break
													case '/':
														if buffer[position] != rune('/') {
															goto l1568
														}
														position++
														break
													case '=':
														if buffer[position] != rune('=') {
															goto l1568
														}
														position++
														break
													case ';':
														if buffer[position] != rune(';') {
															goto l1568
														}
														position++
														break
													case ',':
														if buffer[position] != rune(',') {
															goto l1568
														}
														position++
														break
													case '+':
														if buffer[position] != rune('+') {
															goto l1568
														}
														position++
														break
													case '*':
														if buffer[position] != rune('*') {
															goto l1568
														}
														position++
														break
													case ')':
														if buffer[position] != rune(')') {
															goto l1568
														}
														position++
														break
													case '(':
														if buffer[position] != rune('(') {
															goto l1568
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l1568
														}
														position++
														break
													case '&':
														if buffer[position] != rune('&') {
															goto l1568
														}
														position++
														break
													case '$':
														if buffer[position] != rune('$') {
															goto l1568
														}
														position++
														break
													case '!':
														if buffer[position] != rune('!') {
															goto l1568
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l1568
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l1568
														}
														position++
														break
													case '~':
														if buffer[position] != rune('~') {
															goto l1568
														}
														position++
														break
													default:
														if buffer[position] != rune('_') {
															goto l1568
														}
														position++
														break
													}
												}

												depth--
												add(rulepnLocalEsc, position1583)
											}
										}
									l1580:
										depth--
										add(ruleplx, position1579)
									}
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1568
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l1568
									}
									position++
									break
								default:
									if !rules[rulepnCharsU]() {
										goto l1568
									}
									break
								}
							}

						l1576:
							{
								position1577, tokenIndex1577, depth1577 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '%', '\\':
										{
											position1586 := position
											depth++
											{
												position1587, tokenIndex1587, depth1587 := position, tokenIndex, depth
												{
													position1589 := position
													depth++
													if buffer[position] != rune('%') {
														goto l1588
													}
													position++
													if !rules[rulehex]() {
														goto l1588
													}
													if !rules[rulehex]() {
														goto l1588
													}
													depth--
													add(rulepercent, position1589)
												}
												goto l1587
											l1588:
												position, tokenIndex, depth = position1587, tokenIndex1587, depth1587
												{
													position1590 := position
													depth++
													if buffer[position] != rune('\\') {
														goto l1577
													}
													position++
													{
														switch buffer[position] {
														case '%':
															if buffer[position] != rune('%') {
																goto l1577
															}
															position++
															break
														case '@':
															if buffer[position] != rune('@') {
																goto l1577
															}
															position++
															break
														case '#':
															if buffer[position] != rune('#') {
																goto l1577
															}
															position++
															break
														case '?':
															if buffer[position] != rune('?') {
																goto l1577
															}
															position++
															break
														case '/':
															if buffer[position] != rune('/') {
																goto l1577
															}
															position++
															break
														case '=':
															if buffer[position] != rune('=') {
																goto l1577
															}
															position++
															break
														case ';':
															if buffer[position] != rune(';') {
																goto l1577
															}
															position++
															break
														case ',':
															if buffer[position] != rune(',') {
																goto l1577
															}
															position++
															break
														case '+':
															if buffer[position] != rune('+') {
																goto l1577
															}
															position++
															break
														case '*':
															if buffer[position] != rune('*') {
																goto l1577
															}
															position++
															break
														case ')':
															if buffer[position] != rune(')') {
																goto l1577
															}
															position++
															break
														case '(':
															if buffer[position] != rune('(') {
																goto l1577
															}
															position++
															break
														case '\'':
															if buffer[position] != rune('\'') {
																goto l1577
															}
															position++
															break
														case '&':
															if buffer[position] != rune('&') {
																goto l1577
															}
															position++
															break
														case '$':
															if buffer[position] != rune('$') {
																goto l1577
															}
															position++
															break
														case '!':
															if buffer[position] != rune('!') {
																goto l1577
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l1577
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l1577
															}
															position++
															break
														case '~':
															if buffer[position] != rune('~') {
																goto l1577
															}
															position++
															break
														default:
															if buffer[position] != rune('_') {
																goto l1577
															}
															position++
															break
														}
													}

													depth--
													add(rulepnLocalEsc, position1590)
												}
											}
										l1587:
											depth--
											add(ruleplx, position1586)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1577
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1577
										}
										position++
										break
									default:
										if !rules[rulepnCharsU]() {
											goto l1577
										}
										break
									}
								}

								goto l1576
							l1577:
								position, tokenIndex, depth = position1577, tokenIndex1577, depth1577
							}
							depth--
							add(rulepnLocal, position1575)
						}
						if !rules[ruleskip]() {
							goto l1568
						}
						depth--
						add(ruleprefixedName, position1572)
					}
				}
			l1570:
				depth--
				add(ruleiriref, position1569)
			}
			return true
		l1568:
			position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
			return false
		},
		/* 73 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1592, tokenIndex1592, depth1592 := position, tokenIndex, depth
			{
				position1593 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1592
				}
				position++
			l1594:
				{
					position1595, tokenIndex1595, depth1595 := position, tokenIndex, depth
					{
						position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1596
						}
						position++
						goto l1595
					l1596:
						position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
					}
					if !matchDot() {
						goto l1595
					}
					goto l1594
				l1595:
					position, tokenIndex, depth = position1595, tokenIndex1595, depth1595
				}
				if buffer[position] != rune('>') {
					goto l1592
				}
				position++
				if !rules[ruleskip]() {
					goto l1592
				}
				depth--
				add(ruleiri, position1593)
			}
			return true
		l1592:
			position, tokenIndex, depth = position1592, tokenIndex1592, depth1592
			return false
		},
		/* 74 prefixedName <- <(pnPrefix? ':' pnLocal skip)> */
		nil,
		/* 75 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1598, tokenIndex1598, depth1598 := position, tokenIndex, depth
			{
				position1599 := position
				depth++
				if !rules[rulestring]() {
					goto l1598
				}
				{
					position1600, tokenIndex1600, depth1600 := position, tokenIndex, depth
					{
						position1602, tokenIndex1602, depth1602 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1603
						}
						position++
						{
							position1606, tokenIndex1606, depth1606 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1607
							}
							position++
							goto l1606
						l1607:
							position, tokenIndex, depth = position1606, tokenIndex1606, depth1606
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1603
							}
							position++
						}
					l1606:
					l1604:
						{
							position1605, tokenIndex1605, depth1605 := position, tokenIndex, depth
							{
								position1608, tokenIndex1608, depth1608 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1609
								}
								position++
								goto l1608
							l1609:
								position, tokenIndex, depth = position1608, tokenIndex1608, depth1608
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1605
								}
								position++
							}
						l1608:
							goto l1604
						l1605:
							position, tokenIndex, depth = position1605, tokenIndex1605, depth1605
						}
					l1610:
						{
							position1611, tokenIndex1611, depth1611 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1611
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1611
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1611
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1611
									}
									position++
									break
								}
							}

						l1612:
							{
								position1613, tokenIndex1613, depth1613 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1613
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1613
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
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
							goto l1610
						l1611:
							position, tokenIndex, depth = position1611, tokenIndex1611, depth1611
						}
						goto l1602
					l1603:
						position, tokenIndex, depth = position1602, tokenIndex1602, depth1602
						if buffer[position] != rune('^') {
							goto l1600
						}
						position++
						if buffer[position] != rune('^') {
							goto l1600
						}
						position++
						if !rules[ruleiriref]() {
							goto l1600
						}
					}
				l1602:
					goto l1601
				l1600:
					position, tokenIndex, depth = position1600, tokenIndex1600, depth1600
				}
			l1601:
				if !rules[ruleskip]() {
					goto l1598
				}
				depth--
				add(ruleliteral, position1599)
			}
			return true
		l1598:
			position, tokenIndex, depth = position1598, tokenIndex1598, depth1598
			return false
		},
		/* 76 string <- <(stringLiteralA / stringLiteralB / stringLiteralLongA / stringLiteralLongB)> */
		func() bool {
			position1616, tokenIndex1616, depth1616 := position, tokenIndex, depth
			{
				position1617 := position
				depth++
				{
					position1618, tokenIndex1618, depth1618 := position, tokenIndex, depth
					{
						position1620 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1619
						}
						position++
					l1621:
						{
							position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
							{
								position1623, tokenIndex1623, depth1623 := position, tokenIndex, depth
								{
									position1625, tokenIndex1625, depth1625 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1625
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1625
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1625
											}
											position++
											break
										default:
											if buffer[position] != rune('\'') {
												goto l1625
											}
											position++
											break
										}
									}

									goto l1624
								l1625:
									position, tokenIndex, depth = position1625, tokenIndex1625, depth1625
								}
								if !matchDot() {
									goto l1624
								}
								goto l1623
							l1624:
								position, tokenIndex, depth = position1623, tokenIndex1623, depth1623
								if !rules[ruleechar]() {
									goto l1622
								}
							}
						l1623:
							goto l1621
						l1622:
							position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
						}
						if buffer[position] != rune('\'') {
							goto l1619
						}
						position++
						depth--
						add(rulestringLiteralA, position1620)
					}
					goto l1618
				l1619:
					position, tokenIndex, depth = position1618, tokenIndex1618, depth1618
					{
						position1628 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1627
						}
						position++
					l1629:
						{
							position1630, tokenIndex1630, depth1630 := position, tokenIndex, depth
							{
								position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
								{
									position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1633
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1633
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1633
											}
											position++
											break
										default:
											if buffer[position] != rune('"') {
												goto l1633
											}
											position++
											break
										}
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
								position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
								if !rules[ruleechar]() {
									goto l1630
								}
							}
						l1631:
							goto l1629
						l1630:
							position, tokenIndex, depth = position1630, tokenIndex1630, depth1630
						}
						if buffer[position] != rune('"') {
							goto l1627
						}
						position++
						depth--
						add(rulestringLiteralB, position1628)
					}
					goto l1618
				l1627:
					position, tokenIndex, depth = position1618, tokenIndex1618, depth1618
					{
						position1636 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1635
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1635
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1635
						}
						position++
					l1637:
						{
							position1638, tokenIndex1638, depth1638 := position, tokenIndex, depth
							{
								position1639, tokenIndex1639, depth1639 := position, tokenIndex, depth
								{
									position1641, tokenIndex1641, depth1641 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l1642
									}
									position++
									goto l1641
								l1642:
									position, tokenIndex, depth = position1641, tokenIndex1641, depth1641
									if buffer[position] != rune('\'') {
										goto l1639
									}
									position++
									if buffer[position] != rune('\'') {
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
										if buffer[position] != rune('\'') {
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
						if buffer[position] != rune('\'') {
							goto l1635
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1635
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1635
						}
						position++
						depth--
						add(rulestringLiteralLongA, position1636)
					}
					goto l1618
				l1635:
					position, tokenIndex, depth = position1618, tokenIndex1618, depth1618
					{
						position1648 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1616
						}
						position++
						if buffer[position] != rune('"') {
							goto l1616
						}
						position++
						if buffer[position] != rune('"') {
							goto l1616
						}
						position++
					l1649:
						{
							position1650, tokenIndex1650, depth1650 := position, tokenIndex, depth
							{
								position1651, tokenIndex1651, depth1651 := position, tokenIndex, depth
								{
									position1653, tokenIndex1653, depth1653 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l1654
									}
									position++
									goto l1653
								l1654:
									position, tokenIndex, depth = position1653, tokenIndex1653, depth1653
									if buffer[position] != rune('"') {
										goto l1651
									}
									position++
									if buffer[position] != rune('"') {
										goto l1651
									}
									position++
								}
							l1653:
								goto l1652
							l1651:
								position, tokenIndex, depth = position1651, tokenIndex1651, depth1651
							}
						l1652:
							{
								position1655, tokenIndex1655, depth1655 := position, tokenIndex, depth
								{
									position1657, tokenIndex1657, depth1657 := position, tokenIndex, depth
									{
										position1658, tokenIndex1658, depth1658 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											goto l1659
										}
										position++
										goto l1658
									l1659:
										position, tokenIndex, depth = position1658, tokenIndex1658, depth1658
										if buffer[position] != rune('\\') {
											goto l1657
										}
										position++
									}
								l1658:
									goto l1656
								l1657:
									position, tokenIndex, depth = position1657, tokenIndex1657, depth1657
								}
								if !matchDot() {
									goto l1656
								}
								goto l1655
							l1656:
								position, tokenIndex, depth = position1655, tokenIndex1655, depth1655
								if !rules[ruleechar]() {
									goto l1650
								}
							}
						l1655:
							goto l1649
						l1650:
							position, tokenIndex, depth = position1650, tokenIndex1650, depth1650
						}
						if buffer[position] != rune('"') {
							goto l1616
						}
						position++
						if buffer[position] != rune('"') {
							goto l1616
						}
						position++
						if buffer[position] != rune('"') {
							goto l1616
						}
						position++
						depth--
						add(rulestringLiteralLongB, position1648)
					}
				}
			l1618:
				depth--
				add(rulestring, position1617)
			}
			return true
		l1616:
			position, tokenIndex, depth = position1616, tokenIndex1616, depth1616
			return false
		},
		/* 77 stringLiteralA <- <('\'' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('\'') '\'')) .) / echar)* '\'')> */
		nil,
		/* 78 stringLiteralB <- <('"' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .) / echar)* '"')> */
		nil,
		/* 79 stringLiteralLongA <- <('\'' '\'' '\'' (('\'' / ('\'' '\''))? ((!('\'' / '\\') .) / echar))* ('\'' '\'' '\''))> */
		nil,
		/* 80 stringLiteralLongB <- <('"' '"' '"' (('"' / ('"' '"'))? ((!('"' / '\\') .) / echar))* ('"' '"' '"'))> */
		nil,
		/* 81 echar <- <('\\' ((&('\'') '\'') | (&('"') '"') | (&('\\') '\\') | (&('f') 'f') | (&('r') 'r') | (&('n') 'n') | (&('b') 'b') | (&('t') 't')))> */
		func() bool {
			position1664, tokenIndex1664, depth1664 := position, tokenIndex, depth
			{
				position1665 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l1664
				}
				position++
				{
					switch buffer[position] {
					case '\'':
						if buffer[position] != rune('\'') {
							goto l1664
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l1664
						}
						position++
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l1664
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l1664
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l1664
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l1664
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l1664
						}
						position++
						break
					default:
						if buffer[position] != rune('t') {
							goto l1664
						}
						position++
						break
					}
				}

				depth--
				add(ruleechar, position1665)
			}
			return true
		l1664:
			position, tokenIndex, depth = position1664, tokenIndex1664, depth1664
			return false
		},
		/* 82 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1667, tokenIndex1667, depth1667 := position, tokenIndex, depth
			{
				position1668 := position
				depth++
				{
					position1669, tokenIndex1669, depth1669 := position, tokenIndex, depth
					{
						position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1672
						}
						position++
						goto l1671
					l1672:
						position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
						if buffer[position] != rune('-') {
							goto l1669
						}
						position++
					}
				l1671:
					goto l1670
				l1669:
					position, tokenIndex, depth = position1669, tokenIndex1669, depth1669
				}
			l1670:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1667
				}
				position++
			l1673:
				{
					position1674, tokenIndex1674, depth1674 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1674
					}
					position++
					goto l1673
				l1674:
					position, tokenIndex, depth = position1674, tokenIndex1674, depth1674
				}
				{
					position1675, tokenIndex1675, depth1675 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1675
					}
					position++
				l1677:
					{
						position1678, tokenIndex1678, depth1678 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1678
						}
						position++
						goto l1677
					l1678:
						position, tokenIndex, depth = position1678, tokenIndex1678, depth1678
					}
					goto l1676
				l1675:
					position, tokenIndex, depth = position1675, tokenIndex1675, depth1675
				}
			l1676:
				if !rules[ruleskip]() {
					goto l1667
				}
				depth--
				add(rulenumericLiteral, position1668)
			}
			return true
		l1667:
			position, tokenIndex, depth = position1667, tokenIndex1667, depth1667
			return false
		},
		/* 83 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 84 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1680, tokenIndex1680, depth1680 := position, tokenIndex, depth
			{
				position1681 := position
				depth++
				{
					position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
					{
						position1684 := position
						depth++
						{
							position1685, tokenIndex1685, depth1685 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1686
							}
							position++
							goto l1685
						l1686:
							position, tokenIndex, depth = position1685, tokenIndex1685, depth1685
							if buffer[position] != rune('T') {
								goto l1683
							}
							position++
						}
					l1685:
						{
							position1687, tokenIndex1687, depth1687 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1688
							}
							position++
							goto l1687
						l1688:
							position, tokenIndex, depth = position1687, tokenIndex1687, depth1687
							if buffer[position] != rune('R') {
								goto l1683
							}
							position++
						}
					l1687:
						{
							position1689, tokenIndex1689, depth1689 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1690
							}
							position++
							goto l1689
						l1690:
							position, tokenIndex, depth = position1689, tokenIndex1689, depth1689
							if buffer[position] != rune('U') {
								goto l1683
							}
							position++
						}
					l1689:
						{
							position1691, tokenIndex1691, depth1691 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1692
							}
							position++
							goto l1691
						l1692:
							position, tokenIndex, depth = position1691, tokenIndex1691, depth1691
							if buffer[position] != rune('E') {
								goto l1683
							}
							position++
						}
					l1691:
						if !rules[ruleskip]() {
							goto l1683
						}
						depth--
						add(ruleTRUE, position1684)
					}
					goto l1682
				l1683:
					position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
					{
						position1693 := position
						depth++
						{
							position1694, tokenIndex1694, depth1694 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1695
							}
							position++
							goto l1694
						l1695:
							position, tokenIndex, depth = position1694, tokenIndex1694, depth1694
							if buffer[position] != rune('F') {
								goto l1680
							}
							position++
						}
					l1694:
						{
							position1696, tokenIndex1696, depth1696 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1697
							}
							position++
							goto l1696
						l1697:
							position, tokenIndex, depth = position1696, tokenIndex1696, depth1696
							if buffer[position] != rune('A') {
								goto l1680
							}
							position++
						}
					l1696:
						{
							position1698, tokenIndex1698, depth1698 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1699
							}
							position++
							goto l1698
						l1699:
							position, tokenIndex, depth = position1698, tokenIndex1698, depth1698
							if buffer[position] != rune('L') {
								goto l1680
							}
							position++
						}
					l1698:
						{
							position1700, tokenIndex1700, depth1700 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1701
							}
							position++
							goto l1700
						l1701:
							position, tokenIndex, depth = position1700, tokenIndex1700, depth1700
							if buffer[position] != rune('S') {
								goto l1680
							}
							position++
						}
					l1700:
						{
							position1702, tokenIndex1702, depth1702 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1703
							}
							position++
							goto l1702
						l1703:
							position, tokenIndex, depth = position1702, tokenIndex1702, depth1702
							if buffer[position] != rune('E') {
								goto l1680
							}
							position++
						}
					l1702:
						if !rules[ruleskip]() {
							goto l1680
						}
						depth--
						add(ruleFALSE, position1693)
					}
				}
			l1682:
				depth--
				add(rulebooleanLiteral, position1681)
			}
			return true
		l1680:
			position, tokenIndex, depth = position1680, tokenIndex1680, depth1680
			return false
		},
		/* 85 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 86 blankNodeLabel <- <('_' ':' (pnCharsU / [0-9]) (((pnCharsU / ((&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))* pnCharsU) / ([0-9] / '-'))? skip)> */
		nil,
		/* 87 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 88 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1707, tokenIndex1707, depth1707 := position, tokenIndex, depth
			{
				position1708 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1707
				}
				position++
			l1709:
				{
					position1710, tokenIndex1710, depth1710 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1710
					}
					goto l1709
				l1710:
					position, tokenIndex, depth = position1710, tokenIndex1710, depth1710
				}
				if buffer[position] != rune(')') {
					goto l1707
				}
				position++
				if !rules[ruleskip]() {
					goto l1707
				}
				depth--
				add(rulenil, position1708)
			}
			return true
		l1707:
			position, tokenIndex, depth = position1707, tokenIndex1707, depth1707
			return false
		},
		/* 89 VARNAME <- <((pnCharsU / [0-9]) (pnCharsU / ((&('â') [‿-⁀]) | (&('Ì' | 'Í') [̀-ͯ]) | (&('Â') '·') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))*)> */
		nil,
		/* 90 pnPrefix <- <(pnCharsBase pnChars*)> */
		func() bool {
			position1712, tokenIndex1712, depth1712 := position, tokenIndex, depth
			{
				position1713 := position
				depth++
				if !rules[rulepnCharsBase]() {
					goto l1712
				}
			l1714:
				{
					position1715, tokenIndex1715, depth1715 := position, tokenIndex, depth
					{
						position1716 := position
						depth++
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1715
								}
								position++
								break
							case '-':
								if buffer[position] != rune('-') {
									goto l1715
								}
								position++
								break
							default:
								if !rules[rulepnCharsU]() {
									goto l1715
								}
								break
							}
						}

						depth--
						add(rulepnChars, position1716)
					}
					goto l1714
				l1715:
					position, tokenIndex, depth = position1715, tokenIndex1715, depth1715
				}
				depth--
				add(rulepnPrefix, position1713)
			}
			return true
		l1712:
			position, tokenIndex, depth = position1712, tokenIndex1712, depth1712
			return false
		},
		/* 91 pnLocal <- <((&('%' | '\\') plx) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&(':') ':') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))+> */
		nil,
		/* 92 pnChars <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))> */
		nil,
		/* 93 pnCharsU <- <(pnCharsBase / '_')> */
		func() bool {
			position1720, tokenIndex1720, depth1720 := position, tokenIndex, depth
			{
				position1721 := position
				depth++
				{
					position1722, tokenIndex1722, depth1722 := position, tokenIndex, depth
					if !rules[rulepnCharsBase]() {
						goto l1723
					}
					goto l1722
				l1723:
					position, tokenIndex, depth = position1722, tokenIndex1722, depth1722
					if buffer[position] != rune('_') {
						goto l1720
					}
					position++
				}
			l1722:
				depth--
				add(rulepnCharsU, position1721)
			}
			return true
		l1720:
			position, tokenIndex, depth = position1720, tokenIndex1720, depth1720
			return false
		},
		/* 94 pnCharsBase <- <([À-Ö] / [Ø-ö] / [Ͱ-ͽ] / [‌-‍] / [⁰-↏] / [豈-﷏] / ((&('ð' | 'ñ' | 'ò' | 'ó') [𐀀-󯿿]) | (&('ï') [ﷰ-�]) | (&('ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í') [、-퟿]) | (&('â') [Ⰰ-⿯]) | (&('Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á') [Ϳ-῿]) | (&('Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë') [ø-˿]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))> */
		func() bool {
			position1724, tokenIndex1724, depth1724 := position, tokenIndex, depth
			{
				position1725 := position
				depth++
				{
					position1726, tokenIndex1726, depth1726 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('À') || c > rune('Ö') {
						goto l1727
					}
					position++
					goto l1726
				l1727:
					position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
					if c := buffer[position]; c < rune('Ø') || c > rune('ö') {
						goto l1728
					}
					position++
					goto l1726
				l1728:
					position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
					if c := buffer[position]; c < rune('Ͱ') || c > rune('ͽ') {
						goto l1729
					}
					position++
					goto l1726
				l1729:
					position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
					if c := buffer[position]; c < rune('\u200c') || c > rune('\u200d') {
						goto l1730
					}
					position++
					goto l1726
				l1730:
					position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
					if c := buffer[position]; c < rune('⁰') || c > rune('\u218f') {
						goto l1731
					}
					position++
					goto l1726
				l1731:
					position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
					if c := buffer[position]; c < rune('豈') || c > rune('\ufdcf') {
						goto l1732
					}
					position++
					goto l1726
				l1732:
					position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
					{
						switch buffer[position] {
						case 'ð', 'ñ', 'ò', 'ó':
							if c := buffer[position]; c < rune('𐀀') || c > rune('\U000effff') {
								goto l1724
							}
							position++
							break
						case 'ï':
							if c := buffer[position]; c < rune('ﷰ') || c > rune('�') {
								goto l1724
							}
							position++
							break
						case 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í':
							if c := buffer[position]; c < rune('、') || c > rune('\ud7ff') {
								goto l1724
							}
							position++
							break
						case 'â':
							if c := buffer[position]; c < rune('Ⰰ') || c > rune('\u2fef') {
								goto l1724
							}
							position++
							break
						case 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'Þ', 'ß', 'à', 'á':
							if c := buffer[position]; c < rune('\u037f') || c > rune('\u1fff') {
								goto l1724
							}
							position++
							break
						case 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ê', 'Ë':
							if c := buffer[position]; c < rune('ø') || c > rune('˿') {
								goto l1724
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1724
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1724
							}
							position++
							break
						}
					}

				}
			l1726:
				depth--
				add(rulepnCharsBase, position1725)
			}
			return true
		l1724:
			position, tokenIndex, depth = position1724, tokenIndex1724, depth1724
			return false
		},
		/* 95 plx <- <(percent / pnLocalEsc)> */
		nil,
		/* 96 percent <- <('%' hex hex)> */
		nil,
		/* 97 hex <- <((&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))> */
		func() bool {
			position1736, tokenIndex1736, depth1736 := position, tokenIndex, depth
			{
				position1737 := position
				depth++
				{
					switch buffer[position] {
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l1736
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1736
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1736
						}
						position++
						break
					}
				}

				depth--
				add(rulehex, position1737)
			}
			return true
		l1736:
			position, tokenIndex, depth = position1736, tokenIndex1736, depth1736
			return false
		},
		/* 98 pnLocalEsc <- <('\\' ((&('%') '%') | (&('@') '@') | (&('#') '#') | (&('?') '?') | (&('/') '/') | (&('=') '=') | (&(';') ';') | (&(',') ',') | (&('+') '+') | (&('*') '*') | (&(')') ')') | (&('(') '(') | (&('\'') '\'') | (&('&') '&') | (&('$') '$') | (&('!') '!') | (&('-') '-') | (&('.') '.') | (&('~') '~') | (&('_') '_')))> */
		nil,
		/* 99 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 100 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 101 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 102 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 103 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 104 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 105 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1746, tokenIndex1746, depth1746 := position, tokenIndex, depth
			{
				position1747 := position
				depth++
				{
					position1748, tokenIndex1748, depth1748 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1749
					}
					position++
					goto l1748
				l1749:
					position, tokenIndex, depth = position1748, tokenIndex1748, depth1748
					if buffer[position] != rune('D') {
						goto l1746
					}
					position++
				}
			l1748:
				{
					position1750, tokenIndex1750, depth1750 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1751
					}
					position++
					goto l1750
				l1751:
					position, tokenIndex, depth = position1750, tokenIndex1750, depth1750
					if buffer[position] != rune('I') {
						goto l1746
					}
					position++
				}
			l1750:
				{
					position1752, tokenIndex1752, depth1752 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1753
					}
					position++
					goto l1752
				l1753:
					position, tokenIndex, depth = position1752, tokenIndex1752, depth1752
					if buffer[position] != rune('S') {
						goto l1746
					}
					position++
				}
			l1752:
				{
					position1754, tokenIndex1754, depth1754 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1755
					}
					position++
					goto l1754
				l1755:
					position, tokenIndex, depth = position1754, tokenIndex1754, depth1754
					if buffer[position] != rune('T') {
						goto l1746
					}
					position++
				}
			l1754:
				{
					position1756, tokenIndex1756, depth1756 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1757
					}
					position++
					goto l1756
				l1757:
					position, tokenIndex, depth = position1756, tokenIndex1756, depth1756
					if buffer[position] != rune('I') {
						goto l1746
					}
					position++
				}
			l1756:
				{
					position1758, tokenIndex1758, depth1758 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1759
					}
					position++
					goto l1758
				l1759:
					position, tokenIndex, depth = position1758, tokenIndex1758, depth1758
					if buffer[position] != rune('N') {
						goto l1746
					}
					position++
				}
			l1758:
				{
					position1760, tokenIndex1760, depth1760 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1761
					}
					position++
					goto l1760
				l1761:
					position, tokenIndex, depth = position1760, tokenIndex1760, depth1760
					if buffer[position] != rune('C') {
						goto l1746
					}
					position++
				}
			l1760:
				{
					position1762, tokenIndex1762, depth1762 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1763
					}
					position++
					goto l1762
				l1763:
					position, tokenIndex, depth = position1762, tokenIndex1762, depth1762
					if buffer[position] != rune('T') {
						goto l1746
					}
					position++
				}
			l1762:
				if !rules[ruleskip]() {
					goto l1746
				}
				depth--
				add(ruleDISTINCT, position1747)
			}
			return true
		l1746:
			position, tokenIndex, depth = position1746, tokenIndex1746, depth1746
			return false
		},
		/* 106 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 107 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 108 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 109 LBRACE <- <('{' skip)> */
		func() bool {
			position1767, tokenIndex1767, depth1767 := position, tokenIndex, depth
			{
				position1768 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1767
				}
				position++
				if !rules[ruleskip]() {
					goto l1767
				}
				depth--
				add(ruleLBRACE, position1768)
			}
			return true
		l1767:
			position, tokenIndex, depth = position1767, tokenIndex1767, depth1767
			return false
		},
		/* 110 RBRACE <- <('}' skip)> */
		func() bool {
			position1769, tokenIndex1769, depth1769 := position, tokenIndex, depth
			{
				position1770 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1769
				}
				position++
				if !rules[ruleskip]() {
					goto l1769
				}
				depth--
				add(ruleRBRACE, position1770)
			}
			return true
		l1769:
			position, tokenIndex, depth = position1769, tokenIndex1769, depth1769
			return false
		},
		/* 111 LBRACK <- <('[' skip)> */
		nil,
		/* 112 RBRACK <- <(']' skip)> */
		nil,
		/* 113 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1773, tokenIndex1773, depth1773 := position, tokenIndex, depth
			{
				position1774 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1773
				}
				position++
				if !rules[ruleskip]() {
					goto l1773
				}
				depth--
				add(ruleSEMICOLON, position1774)
			}
			return true
		l1773:
			position, tokenIndex, depth = position1773, tokenIndex1773, depth1773
			return false
		},
		/* 114 COMMA <- <(',' skip)> */
		func() bool {
			position1775, tokenIndex1775, depth1775 := position, tokenIndex, depth
			{
				position1776 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1775
				}
				position++
				if !rules[ruleskip]() {
					goto l1775
				}
				depth--
				add(ruleCOMMA, position1776)
			}
			return true
		l1775:
			position, tokenIndex, depth = position1775, tokenIndex1775, depth1775
			return false
		},
		/* 115 DOT <- <('.' skip)> */
		func() bool {
			position1777, tokenIndex1777, depth1777 := position, tokenIndex, depth
			{
				position1778 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1777
				}
				position++
				if !rules[ruleskip]() {
					goto l1777
				}
				depth--
				add(ruleDOT, position1778)
			}
			return true
		l1777:
			position, tokenIndex, depth = position1777, tokenIndex1777, depth1777
			return false
		},
		/* 116 COLON <- <(':' skip)> */
		nil,
		/* 117 PIPE <- <('|' skip)> */
		func() bool {
			position1780, tokenIndex1780, depth1780 := position, tokenIndex, depth
			{
				position1781 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1780
				}
				position++
				if !rules[ruleskip]() {
					goto l1780
				}
				depth--
				add(rulePIPE, position1781)
			}
			return true
		l1780:
			position, tokenIndex, depth = position1780, tokenIndex1780, depth1780
			return false
		},
		/* 118 SLASH <- <('/' skip)> */
		func() bool {
			position1782, tokenIndex1782, depth1782 := position, tokenIndex, depth
			{
				position1783 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1782
				}
				position++
				if !rules[ruleskip]() {
					goto l1782
				}
				depth--
				add(ruleSLASH, position1783)
			}
			return true
		l1782:
			position, tokenIndex, depth = position1782, tokenIndex1782, depth1782
			return false
		},
		/* 119 INVERSE <- <('^' skip)> */
		func() bool {
			position1784, tokenIndex1784, depth1784 := position, tokenIndex, depth
			{
				position1785 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1784
				}
				position++
				if !rules[ruleskip]() {
					goto l1784
				}
				depth--
				add(ruleINVERSE, position1785)
			}
			return true
		l1784:
			position, tokenIndex, depth = position1784, tokenIndex1784, depth1784
			return false
		},
		/* 120 LPAREN <- <('(' skip)> */
		func() bool {
			position1786, tokenIndex1786, depth1786 := position, tokenIndex, depth
			{
				position1787 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1786
				}
				position++
				if !rules[ruleskip]() {
					goto l1786
				}
				depth--
				add(ruleLPAREN, position1787)
			}
			return true
		l1786:
			position, tokenIndex, depth = position1786, tokenIndex1786, depth1786
			return false
		},
		/* 121 RPAREN <- <(')' skip)> */
		func() bool {
			position1788, tokenIndex1788, depth1788 := position, tokenIndex, depth
			{
				position1789 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1788
				}
				position++
				if !rules[ruleskip]() {
					goto l1788
				}
				depth--
				add(ruleRPAREN, position1789)
			}
			return true
		l1788:
			position, tokenIndex, depth = position1788, tokenIndex1788, depth1788
			return false
		},
		/* 122 ISA <- <('a' skip)> */
		func() bool {
			position1790, tokenIndex1790, depth1790 := position, tokenIndex, depth
			{
				position1791 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1790
				}
				position++
				if !rules[ruleskip]() {
					goto l1790
				}
				depth--
				add(ruleISA, position1791)
			}
			return true
		l1790:
			position, tokenIndex, depth = position1790, tokenIndex1790, depth1790
			return false
		},
		/* 123 NOT <- <('!' skip)> */
		func() bool {
			position1792, tokenIndex1792, depth1792 := position, tokenIndex, depth
			{
				position1793 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1792
				}
				position++
				if !rules[ruleskip]() {
					goto l1792
				}
				depth--
				add(ruleNOT, position1793)
			}
			return true
		l1792:
			position, tokenIndex, depth = position1792, tokenIndex1792, depth1792
			return false
		},
		/* 124 STAR <- <('*' skip)> */
		func() bool {
			position1794, tokenIndex1794, depth1794 := position, tokenIndex, depth
			{
				position1795 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1794
				}
				position++
				if !rules[ruleskip]() {
					goto l1794
				}
				depth--
				add(ruleSTAR, position1795)
			}
			return true
		l1794:
			position, tokenIndex, depth = position1794, tokenIndex1794, depth1794
			return false
		},
		/* 125 QUESTION <- <('?' skip)> */
		nil,
		/* 126 PLUS <- <('+' skip)> */
		func() bool {
			position1797, tokenIndex1797, depth1797 := position, tokenIndex, depth
			{
				position1798 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1797
				}
				position++
				if !rules[ruleskip]() {
					goto l1797
				}
				depth--
				add(rulePLUS, position1798)
			}
			return true
		l1797:
			position, tokenIndex, depth = position1797, tokenIndex1797, depth1797
			return false
		},
		/* 127 MINUS <- <('-' skip)> */
		func() bool {
			position1799, tokenIndex1799, depth1799 := position, tokenIndex, depth
			{
				position1800 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1799
				}
				position++
				if !rules[ruleskip]() {
					goto l1799
				}
				depth--
				add(ruleMINUS, position1800)
			}
			return true
		l1799:
			position, tokenIndex, depth = position1799, tokenIndex1799, depth1799
			return false
		},
		/* 128 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 129 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 130 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 131 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 132 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1805, tokenIndex1805, depth1805 := position, tokenIndex, depth
			{
				position1806 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1805
				}
				position++
			l1807:
				{
					position1808, tokenIndex1808, depth1808 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1808
					}
					position++
					goto l1807
				l1808:
					position, tokenIndex, depth = position1808, tokenIndex1808, depth1808
				}
				if !rules[ruleskip]() {
					goto l1805
				}
				depth--
				add(ruleINTEGER, position1806)
			}
			return true
		l1805:
			position, tokenIndex, depth = position1805, tokenIndex1805, depth1805
			return false
		},
		/* 133 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 134 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 135 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 136 OR <- <('|' '|' skip)> */
		nil,
		/* 137 AND <- <('&' '&' skip)> */
		nil,
		/* 138 EQ <- <('=' skip)> */
		func() bool {
			position1814, tokenIndex1814, depth1814 := position, tokenIndex, depth
			{
				position1815 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1814
				}
				position++
				if !rules[ruleskip]() {
					goto l1814
				}
				depth--
				add(ruleEQ, position1815)
			}
			return true
		l1814:
			position, tokenIndex, depth = position1814, tokenIndex1814, depth1814
			return false
		},
		/* 139 NE <- <('!' '=' skip)> */
		nil,
		/* 140 GT <- <('>' skip)> */
		nil,
		/* 141 LT <- <('<' skip)> */
		nil,
		/* 142 LE <- <('<' '=' skip)> */
		nil,
		/* 143 GE <- <('>' '=' skip)> */
		nil,
		/* 144 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 145 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 146 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1823, tokenIndex1823, depth1823 := position, tokenIndex, depth
			{
				position1824 := position
				depth++
				{
					position1825, tokenIndex1825, depth1825 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1826
					}
					position++
					goto l1825
				l1826:
					position, tokenIndex, depth = position1825, tokenIndex1825, depth1825
					if buffer[position] != rune('A') {
						goto l1823
					}
					position++
				}
			l1825:
				{
					position1827, tokenIndex1827, depth1827 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1828
					}
					position++
					goto l1827
				l1828:
					position, tokenIndex, depth = position1827, tokenIndex1827, depth1827
					if buffer[position] != rune('S') {
						goto l1823
					}
					position++
				}
			l1827:
				if !rules[ruleskip]() {
					goto l1823
				}
				depth--
				add(ruleAS, position1824)
			}
			return true
		l1823:
			position, tokenIndex, depth = position1823, tokenIndex1823, depth1823
			return false
		},
		/* 147 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 148 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 149 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 150 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 151 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 152 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 153 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 154 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 155 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 156 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 157 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 158 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 159 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 160 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 161 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 162 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 163 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 164 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 165 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 166 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 167 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 168 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 169 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 170 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 171 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 172 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 173 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 174 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 175 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 176 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 177 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 178 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 179 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 180 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 181 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 182 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 183 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 184 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 185 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 186 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 187 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 188 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 189 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 190 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 191 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 192 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 193 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 194 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 195 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 196 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 197 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 198 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 199 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 200 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 201 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 202 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 203 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 204 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 205 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 206 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 207 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 208 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 209 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 210 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 211 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 212 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 213 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 214 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 215 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1897, tokenIndex1897, depth1897 := position, tokenIndex, depth
			{
				position1898 := position
				depth++
				{
					position1899, tokenIndex1899, depth1899 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1900
					}
					position++
					goto l1899
				l1900:
					position, tokenIndex, depth = position1899, tokenIndex1899, depth1899
					if buffer[position] != rune('B') {
						goto l1897
					}
					position++
				}
			l1899:
				{
					position1901, tokenIndex1901, depth1901 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1902
					}
					position++
					goto l1901
				l1902:
					position, tokenIndex, depth = position1901, tokenIndex1901, depth1901
					if buffer[position] != rune('Y') {
						goto l1897
					}
					position++
				}
			l1901:
				if !rules[ruleskip]() {
					goto l1897
				}
				depth--
				add(ruleBY, position1898)
			}
			return true
		l1897:
			position, tokenIndex, depth = position1897, tokenIndex1897, depth1897
			return false
		},
		/* 216 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 217 GRAPH <- <(('g' / 'G') ('r' / 'R') ('a' / 'A') ('p' / 'P') ('h' / 'H') skip)> */
		nil,
		/* 218 MINUSSETOPER <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('s' / 'S') skip)> */
		nil,
		/* 219 skip <- <(<(ws / comment)*> Action13)> */
		func() bool {
			{
				position1907 := position
				depth++
				{
					position1908 := position
					depth++
				l1909:
					{
						position1910, tokenIndex1910, depth1910 := position, tokenIndex, depth
						{
							position1911, tokenIndex1911, depth1911 := position, tokenIndex, depth
							if !rules[rulews]() {
								goto l1912
							}
							goto l1911
						l1912:
							position, tokenIndex, depth = position1911, tokenIndex1911, depth1911
							{
								position1913 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1910
								}
								position++
							l1914:
								{
									position1915, tokenIndex1915, depth1915 := position, tokenIndex, depth
									{
										position1916, tokenIndex1916, depth1916 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1916
										}
										goto l1915
									l1916:
										position, tokenIndex, depth = position1916, tokenIndex1916, depth1916
									}
									if !matchDot() {
										goto l1915
									}
									goto l1914
								l1915:
									position, tokenIndex, depth = position1915, tokenIndex1915, depth1915
								}
								if !rules[ruleendOfLine]() {
									goto l1910
								}
								depth--
								add(rulecomment, position1913)
							}
						}
					l1911:
						goto l1909
					l1910:
						position, tokenIndex, depth = position1910, tokenIndex1910, depth1910
					}
					depth--
					add(rulePegText, position1908)
				}
				{
					add(ruleAction13, position)
				}
				depth--
				add(ruleskip, position1907)
			}
			return true
		},
		/* 220 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1918, tokenIndex1918, depth1918 := position, tokenIndex, depth
			{
				position1919 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1918
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1918
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1918
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1918
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1918
						}
						break
					}
				}

				depth--
				add(rulews, position1919)
			}
			return true
		l1918:
			position, tokenIndex, depth = position1918, tokenIndex1918, depth1918
			return false
		},
		/* 221 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 222 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1922, tokenIndex1922, depth1922 := position, tokenIndex, depth
			{
				position1923 := position
				depth++
				{
					position1924, tokenIndex1924, depth1924 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1925
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1925
					}
					position++
					goto l1924
				l1925:
					position, tokenIndex, depth = position1924, tokenIndex1924, depth1924
					if buffer[position] != rune('\n') {
						goto l1926
					}
					position++
					goto l1924
				l1926:
					position, tokenIndex, depth = position1924, tokenIndex1924, depth1924
					if buffer[position] != rune('\r') {
						goto l1922
					}
					position++
				}
			l1924:
				depth--
				add(ruleendOfLine, position1923)
			}
			return true
		l1922:
			position, tokenIndex, depth = position1922, tokenIndex1922, depth1922
			return false
		},
		nil,
		/* 225 Action0 <- <{ p.addPrefix(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 226 SERVICE <- <> */
		nil,
		/* 227 SILENT <- <> */
		nil,
		/* 228 Action1 <- <{ p.S = p.skipped(buffer, begin, end) }> */
		nil,
		/* 229 Action2 <- <{ p.S = p.skipped(buffer, begin, end) }> */
		nil,
		/* 230 Action3 <- <{ p.S = "?POF" }> */
		nil,
		/* 231 Action4 <- <{ p.P = "?POF" }> */
		nil,
		/* 232 Action5 <- <{ p.P = p.skipped(buffer, begin, end) }> */
		nil,
		/* 233 Action6 <- <{ p.P = p.skipped(buffer, begin, end) }> */
		nil,
		/* 234 Action7 <- <{ p.O = "?POF"; p.addTriplePattern() }> */
		nil,
		/* 235 Action8 <- <{ p.O = p.skipped(buffer, begin, end); p.addTriplePattern() }> */
		nil,
		/* 236 Action9 <- <{ p.O = "?FillVar"; p.addTriplePattern() }> */
		nil,
		/* 237 Action10 <- <{ p.setPrefix(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 238 Action11 <- <{ p.setPathLength(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 239 Action12 <- <{ p.setKeyword(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 240 Action13 <- <{ p.skipBegin = begin }> */
		nil,
	}
	p.rules = rules
}
