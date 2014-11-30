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
	rulenoPofPropertyListPath
	rulepofPropertyListPath
	ruleverbPath
	rulepath
	rulepathAlternative
	rulepathSequence
	rulepathElt
	rulepathPrimary
	rulepathNegatedPropertySet
	rulepathOneInPropertySet
	rulepathMod
	rulefillObjectListPath
	rulefillObjectPath
	ruleobjectListPath
	ruleobjectPath
	ruleobject
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
	"noPofPropertyListPath",
	"pofPropertyListPath",
	"verbPath",
	"path",
	"pathAlternative",
	"pathSequence",
	"pathElt",
	"pathPrimary",
	"pathNegatedPropertySet",
	"pathOneInPropertySet",
	"pathMod",
	"fillObjectListPath",
	"fillObjectPath",
	"objectListPath",
	"objectPath",
	"object",
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
	rules  [246]func() bool
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
			p.P = p.skipped(buffer, begin, end)
		case ruleAction5:
			p.P = "?POF"
		case ruleAction6:
			p.P = p.skipped(buffer, begin, end)
		case ruleAction7:
			p.O = "?FillVar"
			p.addTriplePattern()
		case ruleAction8:
			p.O = "?POF"
			p.addTriplePattern()
		case ruleAction9:
			p.O = p.skipped(buffer, begin, end)
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
		/* 34 propertyListPath <- <((pofPropertyListPath / noPofPropertyListPath) (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position398, tokenIndex398, depth398 := position, tokenIndex, depth
			{
				position399 := position
				depth++
				{
					position400, tokenIndex400, depth400 := position, tokenIndex, depth
					{
						position402 := position
						depth++
						if !rules[rulepof]() {
							goto l401
						}
						{
							add(ruleAction5, position)
						}
						{
							position404 := position
							depth++
							if !rules[rulefillObjectPath]() {
								goto l401
							}
						l405:
							{
								position406, tokenIndex406, depth406 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l406
								}
								if !rules[rulefillObjectPath]() {
									goto l406
								}
								goto l405
							l406:
								position, tokenIndex, depth = position406, tokenIndex406, depth406
							}
							depth--
							add(rulefillObjectListPath, position404)
						}
						depth--
						add(rulepofPropertyListPath, position402)
					}
					goto l400
				l401:
					position, tokenIndex, depth = position400, tokenIndex400, depth400
					{
						position407 := position
						depth++
						{
							position408, tokenIndex408, depth408 := position, tokenIndex, depth
							{
								position410 := position
								depth++
								if !rules[rulevar]() {
									goto l409
								}
								depth--
								add(rulePegText, position410)
							}
							{
								add(ruleAction4, position)
							}
							goto l408
						l409:
							position, tokenIndex, depth = position408, tokenIndex408, depth408
							{
								position412 := position
								depth++
								if !rules[rulepath]() {
									goto l398
								}
								depth--
								add(ruleverbPath, position412)
							}
						}
					l408:
						{
							position413 := position
							depth++
							if !rules[ruleobjectPath]() {
								goto l398
							}
						l414:
							{
								position415, tokenIndex415, depth415 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l415
								}
								if !rules[ruleobjectPath]() {
									goto l415
								}
								goto l414
							l415:
								position, tokenIndex, depth = position415, tokenIndex415, depth415
							}
							depth--
							add(ruleobjectListPath, position413)
						}
						depth--
						add(rulenoPofPropertyListPath, position407)
					}
				}
			l400:
				{
					position416, tokenIndex416, depth416 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l416
					}
					{
						position418, tokenIndex418, depth418 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l418
						}
						goto l419
					l418:
						position, tokenIndex, depth = position418, tokenIndex418, depth418
					}
				l419:
					goto l417
				l416:
					position, tokenIndex, depth = position416, tokenIndex416, depth416
				}
			l417:
				depth--
				add(rulepropertyListPath, position399)
			}
			return true
		l398:
			position, tokenIndex, depth = position398, tokenIndex398, depth398
			return false
		},
		/* 35 noPofPropertyListPath <- <(((<var> Action4) / verbPath) objectListPath)> */
		nil,
		/* 36 pofPropertyListPath <- <(pof Action5 fillObjectListPath)> */
		nil,
		/* 37 verbPath <- <path> */
		nil,
		/* 38 path <- <pathAlternative> */
		func() bool {
			position423, tokenIndex423, depth423 := position, tokenIndex, depth
			{
				position424 := position
				depth++
				{
					position425 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l423
					}
				l426:
					{
						position427, tokenIndex427, depth427 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l427
						}
						if !rules[rulepathSequence]() {
							goto l427
						}
						goto l426
					l427:
						position, tokenIndex, depth = position427, tokenIndex427, depth427
					}
					depth--
					add(rulepathAlternative, position425)
				}
				depth--
				add(rulepath, position424)
			}
			return true
		l423:
			position, tokenIndex, depth = position423, tokenIndex423, depth423
			return false
		},
		/* 39 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 40 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position429, tokenIndex429, depth429 := position, tokenIndex, depth
			{
				position430 := position
				depth++
				{
					position431 := position
					depth++
					{
						position432 := position
						depth++
						{
							position433, tokenIndex433, depth433 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l433
							}
							goto l434
						l433:
							position, tokenIndex, depth = position433, tokenIndex433, depth433
						}
					l434:
						{
							position435 := position
							depth++
							{
								position436, tokenIndex436, depth436 := position, tokenIndex, depth
								if !rules[ruleiriref]() {
									goto l437
								}
								goto l436
							l437:
								position, tokenIndex, depth = position436, tokenIndex436, depth436
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l429
										}
										if !rules[rulepath]() {
											goto l429
										}
										if !rules[ruleRPAREN]() {
											goto l429
										}
										break
									case '!':
										if !rules[ruleNOT]() {
											goto l429
										}
										{
											position439 := position
											depth++
											{
												position440, tokenIndex440, depth440 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l441
												}
												goto l440
											l441:
												position, tokenIndex, depth = position440, tokenIndex440, depth440
												if !rules[ruleLPAREN]() {
													goto l429
												}
												{
													position442, tokenIndex442, depth442 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l442
													}
												l444:
													{
														position445, tokenIndex445, depth445 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l445
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l445
														}
														goto l444
													l445:
														position, tokenIndex, depth = position445, tokenIndex445, depth445
													}
													goto l443
												l442:
													position, tokenIndex, depth = position442, tokenIndex442, depth442
												}
											l443:
												if !rules[ruleRPAREN]() {
													goto l429
												}
											}
										l440:
											depth--
											add(rulepathNegatedPropertySet, position439)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l429
										}
										break
									}
								}

							}
						l436:
							depth--
							add(rulepathPrimary, position435)
						}
						{
							position446, tokenIndex446, depth446 := position, tokenIndex, depth
							{
								position448 := position
								depth++
								{
									switch buffer[position] {
									case '+':
										if !rules[rulePLUS]() {
											goto l446
										}
										break
									case '?':
										{
											position450 := position
											depth++
											if buffer[position] != rune('?') {
												goto l446
											}
											position++
											if !rules[ruleskip]() {
												goto l446
											}
											depth--
											add(ruleQUESTION, position450)
										}
										break
									default:
										if !rules[ruleSTAR]() {
											goto l446
										}
										break
									}
								}

								{
									position451, tokenIndex451, depth451 := position, tokenIndex, depth
									if !matchDot() {
										goto l451
									}
									goto l446
								l451:
									position, tokenIndex, depth = position451, tokenIndex451, depth451
								}
								depth--
								add(rulepathMod, position448)
							}
							goto l447
						l446:
							position, tokenIndex, depth = position446, tokenIndex446, depth446
						}
					l447:
						depth--
						add(rulepathElt, position432)
					}
					depth--
					add(rulePegText, position431)
				}
				{
					add(ruleAction6, position)
				}
			l453:
				{
					position454, tokenIndex454, depth454 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l454
					}
					if !rules[rulepathSequence]() {
						goto l454
					}
					goto l453
				l454:
					position, tokenIndex, depth = position454, tokenIndex454, depth454
				}
				depth--
				add(rulepathSequence, position430)
			}
			return true
		l429:
			position, tokenIndex, depth = position429, tokenIndex429, depth429
			return false
		},
		/* 41 pathElt <- <(INVERSE? pathPrimary pathMod?)> */
		nil,
		/* 42 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 43 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 44 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position458, tokenIndex458, depth458 := position, tokenIndex, depth
			{
				position459 := position
				depth++
				{
					position460, tokenIndex460, depth460 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l461
					}
					goto l460
				l461:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
					if !rules[ruleISA]() {
						goto l462
					}
					goto l460
				l462:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
					if !rules[ruleINVERSE]() {
						goto l458
					}
					{
						position463, tokenIndex463, depth463 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l464
						}
						goto l463
					l464:
						position, tokenIndex, depth = position463, tokenIndex463, depth463
						if !rules[ruleISA]() {
							goto l458
						}
					}
				l463:
				}
			l460:
				depth--
				add(rulepathOneInPropertySet, position459)
			}
			return true
		l458:
			position, tokenIndex, depth = position458, tokenIndex458, depth458
			return false
		},
		/* 45 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !.)> */
		nil,
		/* 46 fillObjectListPath <- <(fillObjectPath (COMMA fillObjectPath)*)> */
		nil,
		/* 47 fillObjectPath <- <(object / Action7)> */
		func() bool {
			{
				position468 := position
				depth++
				{
					position469, tokenIndex469, depth469 := position, tokenIndex, depth
					if !rules[ruleobject]() {
						goto l470
					}
					goto l469
				l470:
					position, tokenIndex, depth = position469, tokenIndex469, depth469
					{
						add(ruleAction7, position)
					}
				}
			l469:
				depth--
				add(rulefillObjectPath, position468)
			}
			return true
		},
		/* 48 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 49 objectPath <- <((pof Action8) / object)> */
		func() bool {
			position473, tokenIndex473, depth473 := position, tokenIndex, depth
			{
				position474 := position
				depth++
				{
					position475, tokenIndex475, depth475 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l476
					}
					{
						add(ruleAction8, position)
					}
					goto l475
				l476:
					position, tokenIndex, depth = position475, tokenIndex475, depth475
					if !rules[ruleobject]() {
						goto l473
					}
				}
			l475:
				depth--
				add(ruleobjectPath, position474)
			}
			return true
		l473:
			position, tokenIndex, depth = position473, tokenIndex473, depth473
			return false
		},
		/* 50 object <- <(<graphNodePath> Action9)> */
		func() bool {
			position478, tokenIndex478, depth478 := position, tokenIndex, depth
			{
				position479 := position
				depth++
				{
					position480 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l478
					}
					depth--
					add(rulePegText, position480)
				}
				{
					add(ruleAction9, position)
				}
				depth--
				add(ruleobject, position479)
			}
			return true
		l478:
			position, tokenIndex, depth = position478, tokenIndex478, depth478
			return false
		},
		/* 51 graphNodePath <- <(var / graphTerm / triplesNodePath)> */
		func() bool {
			position482, tokenIndex482, depth482 := position, tokenIndex, depth
			{
				position483 := position
				depth++
				{
					position484, tokenIndex484, depth484 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l485
					}
					goto l484
				l485:
					position, tokenIndex, depth = position484, tokenIndex484, depth484
					if !rules[rulegraphTerm]() {
						goto l486
					}
					goto l484
				l486:
					position, tokenIndex, depth = position484, tokenIndex484, depth484
					if !rules[ruletriplesNodePath]() {
						goto l482
					}
				}
			l484:
				depth--
				add(rulegraphNodePath, position483)
			}
			return true
		l482:
			position, tokenIndex, depth = position482, tokenIndex482, depth482
			return false
		},
		/* 52 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position488 := position
				depth++
				{
					position489, tokenIndex489, depth489 := position, tokenIndex, depth
					{
						position491, tokenIndex491, depth491 := position, tokenIndex, depth
						{
							position493 := position
							depth++
							{
								position494, tokenIndex494, depth494 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l495
								}
								position++
								goto l494
							l495:
								position, tokenIndex, depth = position494, tokenIndex494, depth494
								if buffer[position] != rune('O') {
									goto l492
								}
								position++
							}
						l494:
							{
								position496, tokenIndex496, depth496 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l497
								}
								position++
								goto l496
							l497:
								position, tokenIndex, depth = position496, tokenIndex496, depth496
								if buffer[position] != rune('R') {
									goto l492
								}
								position++
							}
						l496:
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
									goto l492
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
									goto l492
								}
								position++
							}
						l500:
							{
								position502, tokenIndex502, depth502 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l503
								}
								position++
								goto l502
							l503:
								position, tokenIndex, depth = position502, tokenIndex502, depth502
								if buffer[position] != rune('R') {
									goto l492
								}
								position++
							}
						l502:
							if !rules[ruleskip]() {
								goto l492
							}
							depth--
							add(ruleORDER, position493)
						}
						if !rules[ruleBY]() {
							goto l492
						}
						{
							position506 := position
							depth++
							{
								position507, tokenIndex507, depth507 := position, tokenIndex, depth
								{
									position509, tokenIndex509, depth509 := position, tokenIndex, depth
									{
										position511, tokenIndex511, depth511 := position, tokenIndex, depth
										{
											position513 := position
											depth++
											{
												position514, tokenIndex514, depth514 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l515
												}
												position++
												goto l514
											l515:
												position, tokenIndex, depth = position514, tokenIndex514, depth514
												if buffer[position] != rune('A') {
													goto l512
												}
												position++
											}
										l514:
											{
												position516, tokenIndex516, depth516 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l517
												}
												position++
												goto l516
											l517:
												position, tokenIndex, depth = position516, tokenIndex516, depth516
												if buffer[position] != rune('S') {
													goto l512
												}
												position++
											}
										l516:
											{
												position518, tokenIndex518, depth518 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l519
												}
												position++
												goto l518
											l519:
												position, tokenIndex, depth = position518, tokenIndex518, depth518
												if buffer[position] != rune('C') {
													goto l512
												}
												position++
											}
										l518:
											if !rules[ruleskip]() {
												goto l512
											}
											depth--
											add(ruleASC, position513)
										}
										goto l511
									l512:
										position, tokenIndex, depth = position511, tokenIndex511, depth511
										{
											position520 := position
											depth++
											{
												position521, tokenIndex521, depth521 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l522
												}
												position++
												goto l521
											l522:
												position, tokenIndex, depth = position521, tokenIndex521, depth521
												if buffer[position] != rune('D') {
													goto l509
												}
												position++
											}
										l521:
											{
												position523, tokenIndex523, depth523 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l524
												}
												position++
												goto l523
											l524:
												position, tokenIndex, depth = position523, tokenIndex523, depth523
												if buffer[position] != rune('E') {
													goto l509
												}
												position++
											}
										l523:
											{
												position525, tokenIndex525, depth525 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l526
												}
												position++
												goto l525
											l526:
												position, tokenIndex, depth = position525, tokenIndex525, depth525
												if buffer[position] != rune('S') {
													goto l509
												}
												position++
											}
										l525:
											{
												position527, tokenIndex527, depth527 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l528
												}
												position++
												goto l527
											l528:
												position, tokenIndex, depth = position527, tokenIndex527, depth527
												if buffer[position] != rune('C') {
													goto l509
												}
												position++
											}
										l527:
											if !rules[ruleskip]() {
												goto l509
											}
											depth--
											add(ruleDESC, position520)
										}
									}
								l511:
									goto l510
								l509:
									position, tokenIndex, depth = position509, tokenIndex509, depth509
								}
							l510:
								if !rules[rulebrackettedExpression]() {
									goto l508
								}
								goto l507
							l508:
								position, tokenIndex, depth = position507, tokenIndex507, depth507
								if !rules[rulefunctionCall]() {
									goto l529
								}
								goto l507
							l529:
								position, tokenIndex, depth = position507, tokenIndex507, depth507
								if !rules[rulebuiltinCall]() {
									goto l530
								}
								goto l507
							l530:
								position, tokenIndex, depth = position507, tokenIndex507, depth507
								if !rules[rulevar]() {
									goto l492
								}
							}
						l507:
							depth--
							add(ruleorderCondition, position506)
						}
					l504:
						{
							position505, tokenIndex505, depth505 := position, tokenIndex, depth
							{
								position531 := position
								depth++
								{
									position532, tokenIndex532, depth532 := position, tokenIndex, depth
									{
										position534, tokenIndex534, depth534 := position, tokenIndex, depth
										{
											position536, tokenIndex536, depth536 := position, tokenIndex, depth
											{
												position538 := position
												depth++
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
														goto l537
													}
													position++
												}
											l539:
												{
													position541, tokenIndex541, depth541 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l542
													}
													position++
													goto l541
												l542:
													position, tokenIndex, depth = position541, tokenIndex541, depth541
													if buffer[position] != rune('S') {
														goto l537
													}
													position++
												}
											l541:
												{
													position543, tokenIndex543, depth543 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l544
													}
													position++
													goto l543
												l544:
													position, tokenIndex, depth = position543, tokenIndex543, depth543
													if buffer[position] != rune('C') {
														goto l537
													}
													position++
												}
											l543:
												if !rules[ruleskip]() {
													goto l537
												}
												depth--
												add(ruleASC, position538)
											}
											goto l536
										l537:
											position, tokenIndex, depth = position536, tokenIndex536, depth536
											{
												position545 := position
												depth++
												{
													position546, tokenIndex546, depth546 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l547
													}
													position++
													goto l546
												l547:
													position, tokenIndex, depth = position546, tokenIndex546, depth546
													if buffer[position] != rune('D') {
														goto l534
													}
													position++
												}
											l546:
												{
													position548, tokenIndex548, depth548 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l549
													}
													position++
													goto l548
												l549:
													position, tokenIndex, depth = position548, tokenIndex548, depth548
													if buffer[position] != rune('E') {
														goto l534
													}
													position++
												}
											l548:
												{
													position550, tokenIndex550, depth550 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l551
													}
													position++
													goto l550
												l551:
													position, tokenIndex, depth = position550, tokenIndex550, depth550
													if buffer[position] != rune('S') {
														goto l534
													}
													position++
												}
											l550:
												{
													position552, tokenIndex552, depth552 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l553
													}
													position++
													goto l552
												l553:
													position, tokenIndex, depth = position552, tokenIndex552, depth552
													if buffer[position] != rune('C') {
														goto l534
													}
													position++
												}
											l552:
												if !rules[ruleskip]() {
													goto l534
												}
												depth--
												add(ruleDESC, position545)
											}
										}
									l536:
										goto l535
									l534:
										position, tokenIndex, depth = position534, tokenIndex534, depth534
									}
								l535:
									if !rules[rulebrackettedExpression]() {
										goto l533
									}
									goto l532
								l533:
									position, tokenIndex, depth = position532, tokenIndex532, depth532
									if !rules[rulefunctionCall]() {
										goto l554
									}
									goto l532
								l554:
									position, tokenIndex, depth = position532, tokenIndex532, depth532
									if !rules[rulebuiltinCall]() {
										goto l555
									}
									goto l532
								l555:
									position, tokenIndex, depth = position532, tokenIndex532, depth532
									if !rules[rulevar]() {
										goto l505
									}
								}
							l532:
								depth--
								add(ruleorderCondition, position531)
							}
							goto l504
						l505:
							position, tokenIndex, depth = position505, tokenIndex505, depth505
						}
						goto l491
					l492:
						position, tokenIndex, depth = position491, tokenIndex491, depth491
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position557 := position
									depth++
									{
										position558, tokenIndex558, depth558 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l559
										}
										position++
										goto l558
									l559:
										position, tokenIndex, depth = position558, tokenIndex558, depth558
										if buffer[position] != rune('H') {
											goto l489
										}
										position++
									}
								l558:
									{
										position560, tokenIndex560, depth560 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l561
										}
										position++
										goto l560
									l561:
										position, tokenIndex, depth = position560, tokenIndex560, depth560
										if buffer[position] != rune('A') {
											goto l489
										}
										position++
									}
								l560:
									{
										position562, tokenIndex562, depth562 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l563
										}
										position++
										goto l562
									l563:
										position, tokenIndex, depth = position562, tokenIndex562, depth562
										if buffer[position] != rune('V') {
											goto l489
										}
										position++
									}
								l562:
									{
										position564, tokenIndex564, depth564 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l565
										}
										position++
										goto l564
									l565:
										position, tokenIndex, depth = position564, tokenIndex564, depth564
										if buffer[position] != rune('I') {
											goto l489
										}
										position++
									}
								l564:
									{
										position566, tokenIndex566, depth566 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l567
										}
										position++
										goto l566
									l567:
										position, tokenIndex, depth = position566, tokenIndex566, depth566
										if buffer[position] != rune('N') {
											goto l489
										}
										position++
									}
								l566:
									{
										position568, tokenIndex568, depth568 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l569
										}
										position++
										goto l568
									l569:
										position, tokenIndex, depth = position568, tokenIndex568, depth568
										if buffer[position] != rune('G') {
											goto l489
										}
										position++
									}
								l568:
									if !rules[ruleskip]() {
										goto l489
									}
									depth--
									add(ruleHAVING, position557)
								}
								if !rules[ruleconstraint]() {
									goto l489
								}
								break
							case 'G', 'g':
								{
									position570 := position
									depth++
									{
										position571, tokenIndex571, depth571 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l572
										}
										position++
										goto l571
									l572:
										position, tokenIndex, depth = position571, tokenIndex571, depth571
										if buffer[position] != rune('G') {
											goto l489
										}
										position++
									}
								l571:
									{
										position573, tokenIndex573, depth573 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l574
										}
										position++
										goto l573
									l574:
										position, tokenIndex, depth = position573, tokenIndex573, depth573
										if buffer[position] != rune('R') {
											goto l489
										}
										position++
									}
								l573:
									{
										position575, tokenIndex575, depth575 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l576
										}
										position++
										goto l575
									l576:
										position, tokenIndex, depth = position575, tokenIndex575, depth575
										if buffer[position] != rune('O') {
											goto l489
										}
										position++
									}
								l575:
									{
										position577, tokenIndex577, depth577 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l578
										}
										position++
										goto l577
									l578:
										position, tokenIndex, depth = position577, tokenIndex577, depth577
										if buffer[position] != rune('U') {
											goto l489
										}
										position++
									}
								l577:
									{
										position579, tokenIndex579, depth579 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l580
										}
										position++
										goto l579
									l580:
										position, tokenIndex, depth = position579, tokenIndex579, depth579
										if buffer[position] != rune('P') {
											goto l489
										}
										position++
									}
								l579:
									if !rules[ruleskip]() {
										goto l489
									}
									depth--
									add(ruleGROUP, position570)
								}
								if !rules[ruleBY]() {
									goto l489
								}
								{
									position583 := position
									depth++
									{
										position584, tokenIndex584, depth584 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l585
										}
										goto l584
									l585:
										position, tokenIndex, depth = position584, tokenIndex584, depth584
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l489
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l489
												}
												if !rules[ruleexpression]() {
													goto l489
												}
												{
													position587, tokenIndex587, depth587 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l587
													}
													if !rules[rulevar]() {
														goto l587
													}
													goto l588
												l587:
													position, tokenIndex, depth = position587, tokenIndex587, depth587
												}
											l588:
												if !rules[ruleRPAREN]() {
													goto l489
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l489
												}
												break
											}
										}

									}
								l584:
									depth--
									add(rulegroupCondition, position583)
								}
							l581:
								{
									position582, tokenIndex582, depth582 := position, tokenIndex, depth
									{
										position589 := position
										depth++
										{
											position590, tokenIndex590, depth590 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l591
											}
											goto l590
										l591:
											position, tokenIndex, depth = position590, tokenIndex590, depth590
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l582
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l582
													}
													if !rules[ruleexpression]() {
														goto l582
													}
													{
														position593, tokenIndex593, depth593 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l593
														}
														if !rules[rulevar]() {
															goto l593
														}
														goto l594
													l593:
														position, tokenIndex, depth = position593, tokenIndex593, depth593
													}
												l594:
													if !rules[ruleRPAREN]() {
														goto l582
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l582
													}
													break
												}
											}

										}
									l590:
										depth--
										add(rulegroupCondition, position589)
									}
									goto l581
								l582:
									position, tokenIndex, depth = position582, tokenIndex582, depth582
								}
								break
							default:
								{
									position595 := position
									depth++
									{
										position596, tokenIndex596, depth596 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l597
										}
										{
											position598, tokenIndex598, depth598 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l598
											}
											goto l599
										l598:
											position, tokenIndex, depth = position598, tokenIndex598, depth598
										}
									l599:
										goto l596
									l597:
										position, tokenIndex, depth = position596, tokenIndex596, depth596
										if !rules[ruleoffset]() {
											goto l489
										}
										{
											position600, tokenIndex600, depth600 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l600
											}
											goto l601
										l600:
											position, tokenIndex, depth = position600, tokenIndex600, depth600
										}
									l601:
									}
								l596:
									depth--
									add(rulelimitOffsetClauses, position595)
								}
								break
							}
						}

					}
				l491:
					goto l490
				l489:
					position, tokenIndex, depth = position489, tokenIndex489, depth489
				}
			l490:
				depth--
				add(rulesolutionModifier, position488)
			}
			return true
		},
		/* 53 groupCondition <- <(functionCall / ((&('$' | '?') var) | (&('(') (LPAREN expression (AS var)? RPAREN)) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'H' | 'I' | 'L' | 'M' | 'N' | 'R' | 'S' | 'T' | 'U' | 'Y' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'h' | 'i' | 'l' | 'm' | 'n' | 'r' | 's' | 't' | 'u' | 'y') builtinCall)))> */
		nil,
		/* 54 orderCondition <- <(((ASC / DESC)? brackettedExpression) / functionCall / builtinCall / var)> */
		nil,
		/* 55 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 56 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position605, tokenIndex605, depth605 := position, tokenIndex, depth
			{
				position606 := position
				depth++
				{
					position607 := position
					depth++
					{
						position608, tokenIndex608, depth608 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l609
						}
						position++
						goto l608
					l609:
						position, tokenIndex, depth = position608, tokenIndex608, depth608
						if buffer[position] != rune('L') {
							goto l605
						}
						position++
					}
				l608:
					{
						position610, tokenIndex610, depth610 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l611
						}
						position++
						goto l610
					l611:
						position, tokenIndex, depth = position610, tokenIndex610, depth610
						if buffer[position] != rune('I') {
							goto l605
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
							goto l605
						}
						position++
					}
				l612:
					{
						position614, tokenIndex614, depth614 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l615
						}
						position++
						goto l614
					l615:
						position, tokenIndex, depth = position614, tokenIndex614, depth614
						if buffer[position] != rune('I') {
							goto l605
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
							goto l605
						}
						position++
					}
				l616:
					if !rules[ruleskip]() {
						goto l605
					}
					depth--
					add(ruleLIMIT, position607)
				}
				if !rules[ruleINTEGER]() {
					goto l605
				}
				depth--
				add(rulelimit, position606)
			}
			return true
		l605:
			position, tokenIndex, depth = position605, tokenIndex605, depth605
			return false
		},
		/* 57 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position618, tokenIndex618, depth618 := position, tokenIndex, depth
			{
				position619 := position
				depth++
				{
					position620 := position
					depth++
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l622
						}
						position++
						goto l621
					l622:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
						if buffer[position] != rune('O') {
							goto l618
						}
						position++
					}
				l621:
					{
						position623, tokenIndex623, depth623 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l624
						}
						position++
						goto l623
					l624:
						position, tokenIndex, depth = position623, tokenIndex623, depth623
						if buffer[position] != rune('F') {
							goto l618
						}
						position++
					}
				l623:
					{
						position625, tokenIndex625, depth625 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l626
						}
						position++
						goto l625
					l626:
						position, tokenIndex, depth = position625, tokenIndex625, depth625
						if buffer[position] != rune('F') {
							goto l618
						}
						position++
					}
				l625:
					{
						position627, tokenIndex627, depth627 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l628
						}
						position++
						goto l627
					l628:
						position, tokenIndex, depth = position627, tokenIndex627, depth627
						if buffer[position] != rune('S') {
							goto l618
						}
						position++
					}
				l627:
					{
						position629, tokenIndex629, depth629 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l630
						}
						position++
						goto l629
					l630:
						position, tokenIndex, depth = position629, tokenIndex629, depth629
						if buffer[position] != rune('E') {
							goto l618
						}
						position++
					}
				l629:
					{
						position631, tokenIndex631, depth631 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l632
						}
						position++
						goto l631
					l632:
						position, tokenIndex, depth = position631, tokenIndex631, depth631
						if buffer[position] != rune('T') {
							goto l618
						}
						position++
					}
				l631:
					if !rules[ruleskip]() {
						goto l618
					}
					depth--
					add(ruleOFFSET, position620)
				}
				if !rules[ruleINTEGER]() {
					goto l618
				}
				depth--
				add(ruleoffset, position619)
			}
			return true
		l618:
			position, tokenIndex, depth = position618, tokenIndex618, depth618
			return false
		},
		/* 58 expression <- <conditionalOrExpression> */
		func() bool {
			position633, tokenIndex633, depth633 := position, tokenIndex, depth
			{
				position634 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l633
				}
				depth--
				add(ruleexpression, position634)
			}
			return true
		l633:
			position, tokenIndex, depth = position633, tokenIndex633, depth633
			return false
		},
		/* 59 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position635, tokenIndex635, depth635 := position, tokenIndex, depth
			{
				position636 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l635
				}
				{
					position637, tokenIndex637, depth637 := position, tokenIndex, depth
					{
						position639 := position
						depth++
						if buffer[position] != rune('|') {
							goto l637
						}
						position++
						if buffer[position] != rune('|') {
							goto l637
						}
						position++
						if !rules[ruleskip]() {
							goto l637
						}
						depth--
						add(ruleOR, position639)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l637
					}
					goto l638
				l637:
					position, tokenIndex, depth = position637, tokenIndex637, depth637
				}
			l638:
				depth--
				add(ruleconditionalOrExpression, position636)
			}
			return true
		l635:
			position, tokenIndex, depth = position635, tokenIndex635, depth635
			return false
		},
		/* 60 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position640, tokenIndex640, depth640 := position, tokenIndex, depth
			{
				position641 := position
				depth++
				{
					position642 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l640
					}
					{
						position643, tokenIndex643, depth643 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position646 := position
									depth++
									{
										position647 := position
										depth++
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
												goto l643
											}
											position++
										}
									l648:
										{
											position650, tokenIndex650, depth650 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l651
											}
											position++
											goto l650
										l651:
											position, tokenIndex, depth = position650, tokenIndex650, depth650
											if buffer[position] != rune('O') {
												goto l643
											}
											position++
										}
									l650:
										{
											position652, tokenIndex652, depth652 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l653
											}
											position++
											goto l652
										l653:
											position, tokenIndex, depth = position652, tokenIndex652, depth652
											if buffer[position] != rune('T') {
												goto l643
											}
											position++
										}
									l652:
										if buffer[position] != rune(' ') {
											goto l643
										}
										position++
										{
											position654, tokenIndex654, depth654 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l655
											}
											position++
											goto l654
										l655:
											position, tokenIndex, depth = position654, tokenIndex654, depth654
											if buffer[position] != rune('I') {
												goto l643
											}
											position++
										}
									l654:
										{
											position656, tokenIndex656, depth656 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l657
											}
											position++
											goto l656
										l657:
											position, tokenIndex, depth = position656, tokenIndex656, depth656
											if buffer[position] != rune('N') {
												goto l643
											}
											position++
										}
									l656:
										if !rules[ruleskip]() {
											goto l643
										}
										depth--
										add(ruleNOTIN, position647)
									}
									if !rules[ruleargList]() {
										goto l643
									}
									depth--
									add(rulenotin, position646)
								}
								break
							case 'I', 'i':
								{
									position658 := position
									depth++
									{
										position659 := position
										depth++
										{
											position660, tokenIndex660, depth660 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l661
											}
											position++
											goto l660
										l661:
											position, tokenIndex, depth = position660, tokenIndex660, depth660
											if buffer[position] != rune('I') {
												goto l643
											}
											position++
										}
									l660:
										{
											position662, tokenIndex662, depth662 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l663
											}
											position++
											goto l662
										l663:
											position, tokenIndex, depth = position662, tokenIndex662, depth662
											if buffer[position] != rune('N') {
												goto l643
											}
											position++
										}
									l662:
										if !rules[ruleskip]() {
											goto l643
										}
										depth--
										add(ruleIN, position659)
									}
									if !rules[ruleargList]() {
										goto l643
									}
									depth--
									add(rulein, position658)
								}
								break
							default:
								{
									position664, tokenIndex664, depth664 := position, tokenIndex, depth
									{
										position666 := position
										depth++
										if buffer[position] != rune('<') {
											goto l665
										}
										position++
										if !rules[ruleskip]() {
											goto l665
										}
										depth--
										add(ruleLT, position666)
									}
									goto l664
								l665:
									position, tokenIndex, depth = position664, tokenIndex664, depth664
									{
										position668 := position
										depth++
										if buffer[position] != rune('>') {
											goto l667
										}
										position++
										if buffer[position] != rune('=') {
											goto l667
										}
										position++
										if !rules[ruleskip]() {
											goto l667
										}
										depth--
										add(ruleGE, position668)
									}
									goto l664
								l667:
									position, tokenIndex, depth = position664, tokenIndex664, depth664
									{
										switch buffer[position] {
										case '>':
											{
												position670 := position
												depth++
												if buffer[position] != rune('>') {
													goto l643
												}
												position++
												if !rules[ruleskip]() {
													goto l643
												}
												depth--
												add(ruleGT, position670)
											}
											break
										case '<':
											{
												position671 := position
												depth++
												if buffer[position] != rune('<') {
													goto l643
												}
												position++
												if buffer[position] != rune('=') {
													goto l643
												}
												position++
												if !rules[ruleskip]() {
													goto l643
												}
												depth--
												add(ruleLE, position671)
											}
											break
										case '!':
											{
												position672 := position
												depth++
												if buffer[position] != rune('!') {
													goto l643
												}
												position++
												if buffer[position] != rune('=') {
													goto l643
												}
												position++
												if !rules[ruleskip]() {
													goto l643
												}
												depth--
												add(ruleNE, position672)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l643
											}
											break
										}
									}

								}
							l664:
								if !rules[rulenumericExpression]() {
									goto l643
								}
								break
							}
						}

						goto l644
					l643:
						position, tokenIndex, depth = position643, tokenIndex643, depth643
					}
				l644:
					depth--
					add(rulevalueLogical, position642)
				}
				{
					position673, tokenIndex673, depth673 := position, tokenIndex, depth
					{
						position675 := position
						depth++
						if buffer[position] != rune('&') {
							goto l673
						}
						position++
						if buffer[position] != rune('&') {
							goto l673
						}
						position++
						if !rules[ruleskip]() {
							goto l673
						}
						depth--
						add(ruleAND, position675)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l673
					}
					goto l674
				l673:
					position, tokenIndex, depth = position673, tokenIndex673, depth673
				}
			l674:
				depth--
				add(ruleconditionalAndExpression, position641)
			}
			return true
		l640:
			position, tokenIndex, depth = position640, tokenIndex640, depth640
			return false
		},
		/* 61 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 62 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position677, tokenIndex677, depth677 := position, tokenIndex, depth
			{
				position678 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l677
				}
			l679:
				{
					position680, tokenIndex680, depth680 := position, tokenIndex, depth
					{
						position681, tokenIndex681, depth681 := position, tokenIndex, depth
						{
							position683, tokenIndex683, depth683 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l684
							}
							goto l683
						l684:
							position, tokenIndex, depth = position683, tokenIndex683, depth683
							if !rules[ruleMINUS]() {
								goto l682
							}
						}
					l683:
						if !rules[rulemultiplicativeExpression]() {
							goto l682
						}
						goto l681
					l682:
						position, tokenIndex, depth = position681, tokenIndex681, depth681
						{
							position685 := position
							depth++
							{
								position686, tokenIndex686, depth686 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l687
								}
								position++
								goto l686
							l687:
								position, tokenIndex, depth = position686, tokenIndex686, depth686
								if buffer[position] != rune('-') {
									goto l680
								}
								position++
							}
						l686:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l680
							}
							position++
						l688:
							{
								position689, tokenIndex689, depth689 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l689
								}
								position++
								goto l688
							l689:
								position, tokenIndex, depth = position689, tokenIndex689, depth689
							}
							{
								position690, tokenIndex690, depth690 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l690
								}
								position++
							l692:
								{
									position693, tokenIndex693, depth693 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l693
									}
									position++
									goto l692
								l693:
									position, tokenIndex, depth = position693, tokenIndex693, depth693
								}
								goto l691
							l690:
								position, tokenIndex, depth = position690, tokenIndex690, depth690
							}
						l691:
							if !rules[ruleskip]() {
								goto l680
							}
							depth--
							add(rulesignedNumericLiteral, position685)
						}
					}
				l681:
					goto l679
				l680:
					position, tokenIndex, depth = position680, tokenIndex680, depth680
				}
				depth--
				add(rulenumericExpression, position678)
			}
			return true
		l677:
			position, tokenIndex, depth = position677, tokenIndex677, depth677
			return false
		},
		/* 63 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position694, tokenIndex694, depth694 := position, tokenIndex, depth
			{
				position695 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l694
				}
			l696:
				{
					position697, tokenIndex697, depth697 := position, tokenIndex, depth
					{
						position698, tokenIndex698, depth698 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l699
						}
						goto l698
					l699:
						position, tokenIndex, depth = position698, tokenIndex698, depth698
						if !rules[ruleSLASH]() {
							goto l697
						}
					}
				l698:
					if !rules[ruleunaryExpression]() {
						goto l697
					}
					goto l696
				l697:
					position, tokenIndex, depth = position697, tokenIndex697, depth697
				}
				depth--
				add(rulemultiplicativeExpression, position695)
			}
			return true
		l694:
			position, tokenIndex, depth = position694, tokenIndex694, depth694
			return false
		},
		/* 64 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position700, tokenIndex700, depth700 := position, tokenIndex, depth
			{
				position701 := position
				depth++
				{
					position702, tokenIndex702, depth702 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l702
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l702
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l702
							}
							break
						}
					}

					goto l703
				l702:
					position, tokenIndex, depth = position702, tokenIndex702, depth702
				}
			l703:
				{
					position705 := position
					depth++
					{
						position706, tokenIndex706, depth706 := position, tokenIndex, depth
						if !rules[rulebuiltinCall]() {
							goto l707
						}
						goto l706
					l707:
						position, tokenIndex, depth = position706, tokenIndex706, depth706
						if !rules[rulefunctionCall]() {
							goto l708
						}
						goto l706
					l708:
						position, tokenIndex, depth = position706, tokenIndex706, depth706
						if !rules[ruleiriref]() {
							goto l709
						}
						goto l706
					l709:
						position, tokenIndex, depth = position706, tokenIndex706, depth706
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position711 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position713 := position
												depth++
												{
													position714 := position
													depth++
													{
														position715, tokenIndex715, depth715 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l716
														}
														position++
														goto l715
													l716:
														position, tokenIndex, depth = position715, tokenIndex715, depth715
														if buffer[position] != rune('G') {
															goto l700
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
															goto l700
														}
														position++
													}
												l717:
													{
														position719, tokenIndex719, depth719 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l720
														}
														position++
														goto l719
													l720:
														position, tokenIndex, depth = position719, tokenIndex719, depth719
														if buffer[position] != rune('O') {
															goto l700
														}
														position++
													}
												l719:
													{
														position721, tokenIndex721, depth721 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l722
														}
														position++
														goto l721
													l722:
														position, tokenIndex, depth = position721, tokenIndex721, depth721
														if buffer[position] != rune('U') {
															goto l700
														}
														position++
													}
												l721:
													{
														position723, tokenIndex723, depth723 := position, tokenIndex, depth
														if buffer[position] != rune('p') {
															goto l724
														}
														position++
														goto l723
													l724:
														position, tokenIndex, depth = position723, tokenIndex723, depth723
														if buffer[position] != rune('P') {
															goto l700
														}
														position++
													}
												l723:
													if buffer[position] != rune('_') {
														goto l700
													}
													position++
													{
														position725, tokenIndex725, depth725 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l726
														}
														position++
														goto l725
													l726:
														position, tokenIndex, depth = position725, tokenIndex725, depth725
														if buffer[position] != rune('C') {
															goto l700
														}
														position++
													}
												l725:
													{
														position727, tokenIndex727, depth727 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l728
														}
														position++
														goto l727
													l728:
														position, tokenIndex, depth = position727, tokenIndex727, depth727
														if buffer[position] != rune('O') {
															goto l700
														}
														position++
													}
												l727:
													{
														position729, tokenIndex729, depth729 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l730
														}
														position++
														goto l729
													l730:
														position, tokenIndex, depth = position729, tokenIndex729, depth729
														if buffer[position] != rune('N') {
															goto l700
														}
														position++
													}
												l729:
													{
														position731, tokenIndex731, depth731 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l732
														}
														position++
														goto l731
													l732:
														position, tokenIndex, depth = position731, tokenIndex731, depth731
														if buffer[position] != rune('C') {
															goto l700
														}
														position++
													}
												l731:
													{
														position733, tokenIndex733, depth733 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l734
														}
														position++
														goto l733
													l734:
														position, tokenIndex, depth = position733, tokenIndex733, depth733
														if buffer[position] != rune('A') {
															goto l700
														}
														position++
													}
												l733:
													{
														position735, tokenIndex735, depth735 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l736
														}
														position++
														goto l735
													l736:
														position, tokenIndex, depth = position735, tokenIndex735, depth735
														if buffer[position] != rune('T') {
															goto l700
														}
														position++
													}
												l735:
													if !rules[ruleskip]() {
														goto l700
													}
													depth--
													add(ruleGROUPCONCAT, position714)
												}
												if !rules[ruleLPAREN]() {
													goto l700
												}
												{
													position737, tokenIndex737, depth737 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l737
													}
													goto l738
												l737:
													position, tokenIndex, depth = position737, tokenIndex737, depth737
												}
											l738:
												if !rules[ruleexpression]() {
													goto l700
												}
												{
													position739, tokenIndex739, depth739 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l739
													}
													{
														position741 := position
														depth++
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
																goto l739
															}
															position++
														}
													l742:
														{
															position744, tokenIndex744, depth744 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l745
															}
															position++
															goto l744
														l745:
															position, tokenIndex, depth = position744, tokenIndex744, depth744
															if buffer[position] != rune('E') {
																goto l739
															}
															position++
														}
													l744:
														{
															position746, tokenIndex746, depth746 := position, tokenIndex, depth
															if buffer[position] != rune('p') {
																goto l747
															}
															position++
															goto l746
														l747:
															position, tokenIndex, depth = position746, tokenIndex746, depth746
															if buffer[position] != rune('P') {
																goto l739
															}
															position++
														}
													l746:
														{
															position748, tokenIndex748, depth748 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l749
															}
															position++
															goto l748
														l749:
															position, tokenIndex, depth = position748, tokenIndex748, depth748
															if buffer[position] != rune('A') {
																goto l739
															}
															position++
														}
													l748:
														{
															position750, tokenIndex750, depth750 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l751
															}
															position++
															goto l750
														l751:
															position, tokenIndex, depth = position750, tokenIndex750, depth750
															if buffer[position] != rune('R') {
																goto l739
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
																goto l739
															}
															position++
														}
													l752:
														{
															position754, tokenIndex754, depth754 := position, tokenIndex, depth
															if buffer[position] != rune('t') {
																goto l755
															}
															position++
															goto l754
														l755:
															position, tokenIndex, depth = position754, tokenIndex754, depth754
															if buffer[position] != rune('T') {
																goto l739
															}
															position++
														}
													l754:
														{
															position756, tokenIndex756, depth756 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l757
															}
															position++
															goto l756
														l757:
															position, tokenIndex, depth = position756, tokenIndex756, depth756
															if buffer[position] != rune('O') {
																goto l739
															}
															position++
														}
													l756:
														{
															position758, tokenIndex758, depth758 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l759
															}
															position++
															goto l758
														l759:
															position, tokenIndex, depth = position758, tokenIndex758, depth758
															if buffer[position] != rune('R') {
																goto l739
															}
															position++
														}
													l758:
														if !rules[ruleskip]() {
															goto l739
														}
														depth--
														add(ruleSEPARATOR, position741)
													}
													if !rules[ruleEQ]() {
														goto l739
													}
													if !rules[rulestring]() {
														goto l739
													}
													goto l740
												l739:
													position, tokenIndex, depth = position739, tokenIndex739, depth739
												}
											l740:
												if !rules[ruleRPAREN]() {
													goto l700
												}
												depth--
												add(rulegroupConcat, position713)
											}
											break
										case 'C', 'c':
											{
												position760 := position
												depth++
												{
													position761 := position
													depth++
													{
														position762, tokenIndex762, depth762 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l763
														}
														position++
														goto l762
													l763:
														position, tokenIndex, depth = position762, tokenIndex762, depth762
														if buffer[position] != rune('C') {
															goto l700
														}
														position++
													}
												l762:
													{
														position764, tokenIndex764, depth764 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l765
														}
														position++
														goto l764
													l765:
														position, tokenIndex, depth = position764, tokenIndex764, depth764
														if buffer[position] != rune('O') {
															goto l700
														}
														position++
													}
												l764:
													{
														position766, tokenIndex766, depth766 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l767
														}
														position++
														goto l766
													l767:
														position, tokenIndex, depth = position766, tokenIndex766, depth766
														if buffer[position] != rune('U') {
															goto l700
														}
														position++
													}
												l766:
													{
														position768, tokenIndex768, depth768 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l769
														}
														position++
														goto l768
													l769:
														position, tokenIndex, depth = position768, tokenIndex768, depth768
														if buffer[position] != rune('N') {
															goto l700
														}
														position++
													}
												l768:
													{
														position770, tokenIndex770, depth770 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l771
														}
														position++
														goto l770
													l771:
														position, tokenIndex, depth = position770, tokenIndex770, depth770
														if buffer[position] != rune('T') {
															goto l700
														}
														position++
													}
												l770:
													if !rules[ruleskip]() {
														goto l700
													}
													depth--
													add(ruleCOUNT, position761)
												}
												if !rules[ruleLPAREN]() {
													goto l700
												}
												{
													position772, tokenIndex772, depth772 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l772
													}
													goto l773
												l772:
													position, tokenIndex, depth = position772, tokenIndex772, depth772
												}
											l773:
												{
													position774, tokenIndex774, depth774 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l775
													}
													goto l774
												l775:
													position, tokenIndex, depth = position774, tokenIndex774, depth774
													if !rules[ruleexpression]() {
														goto l700
													}
												}
											l774:
												if !rules[ruleRPAREN]() {
													goto l700
												}
												depth--
												add(rulecount, position760)
											}
											break
										default:
											{
												position776, tokenIndex776, depth776 := position, tokenIndex, depth
												{
													position778 := position
													depth++
													{
														position779, tokenIndex779, depth779 := position, tokenIndex, depth
														if buffer[position] != rune('s') {
															goto l780
														}
														position++
														goto l779
													l780:
														position, tokenIndex, depth = position779, tokenIndex779, depth779
														if buffer[position] != rune('S') {
															goto l777
														}
														position++
													}
												l779:
													{
														position781, tokenIndex781, depth781 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l782
														}
														position++
														goto l781
													l782:
														position, tokenIndex, depth = position781, tokenIndex781, depth781
														if buffer[position] != rune('U') {
															goto l777
														}
														position++
													}
												l781:
													{
														position783, tokenIndex783, depth783 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l784
														}
														position++
														goto l783
													l784:
														position, tokenIndex, depth = position783, tokenIndex783, depth783
														if buffer[position] != rune('M') {
															goto l777
														}
														position++
													}
												l783:
													if !rules[ruleskip]() {
														goto l777
													}
													depth--
													add(ruleSUM, position778)
												}
												goto l776
											l777:
												position, tokenIndex, depth = position776, tokenIndex776, depth776
												{
													position786 := position
													depth++
													{
														position787, tokenIndex787, depth787 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l788
														}
														position++
														goto l787
													l788:
														position, tokenIndex, depth = position787, tokenIndex787, depth787
														if buffer[position] != rune('M') {
															goto l785
														}
														position++
													}
												l787:
													{
														position789, tokenIndex789, depth789 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l790
														}
														position++
														goto l789
													l790:
														position, tokenIndex, depth = position789, tokenIndex789, depth789
														if buffer[position] != rune('I') {
															goto l785
														}
														position++
													}
												l789:
													{
														position791, tokenIndex791, depth791 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l792
														}
														position++
														goto l791
													l792:
														position, tokenIndex, depth = position791, tokenIndex791, depth791
														if buffer[position] != rune('N') {
															goto l785
														}
														position++
													}
												l791:
													if !rules[ruleskip]() {
														goto l785
													}
													depth--
													add(ruleMIN, position786)
												}
												goto l776
											l785:
												position, tokenIndex, depth = position776, tokenIndex776, depth776
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position794 := position
															depth++
															{
																position795, tokenIndex795, depth795 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l796
																}
																position++
																goto l795
															l796:
																position, tokenIndex, depth = position795, tokenIndex795, depth795
																if buffer[position] != rune('S') {
																	goto l700
																}
																position++
															}
														l795:
															{
																position797, tokenIndex797, depth797 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l798
																}
																position++
																goto l797
															l798:
																position, tokenIndex, depth = position797, tokenIndex797, depth797
																if buffer[position] != rune('A') {
																	goto l700
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
																	goto l700
																}
																position++
															}
														l799:
															{
																position801, tokenIndex801, depth801 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l802
																}
																position++
																goto l801
															l802:
																position, tokenIndex, depth = position801, tokenIndex801, depth801
																if buffer[position] != rune('P') {
																	goto l700
																}
																position++
															}
														l801:
															{
																position803, tokenIndex803, depth803 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l804
																}
																position++
																goto l803
															l804:
																position, tokenIndex, depth = position803, tokenIndex803, depth803
																if buffer[position] != rune('L') {
																	goto l700
																}
																position++
															}
														l803:
															{
																position805, tokenIndex805, depth805 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l806
																}
																position++
																goto l805
															l806:
																position, tokenIndex, depth = position805, tokenIndex805, depth805
																if buffer[position] != rune('E') {
																	goto l700
																}
																position++
															}
														l805:
															if !rules[ruleskip]() {
																goto l700
															}
															depth--
															add(ruleSAMPLE, position794)
														}
														break
													case 'A', 'a':
														{
															position807 := position
															depth++
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
																	goto l700
																}
																position++
															}
														l808:
															{
																position810, tokenIndex810, depth810 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l811
																}
																position++
																goto l810
															l811:
																position, tokenIndex, depth = position810, tokenIndex810, depth810
																if buffer[position] != rune('V') {
																	goto l700
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
																	goto l700
																}
																position++
															}
														l812:
															if !rules[ruleskip]() {
																goto l700
															}
															depth--
															add(ruleAVG, position807)
														}
														break
													default:
														{
															position814 := position
															depth++
															{
																position815, tokenIndex815, depth815 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l816
																}
																position++
																goto l815
															l816:
																position, tokenIndex, depth = position815, tokenIndex815, depth815
																if buffer[position] != rune('M') {
																	goto l700
																}
																position++
															}
														l815:
															{
																position817, tokenIndex817, depth817 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l818
																}
																position++
																goto l817
															l818:
																position, tokenIndex, depth = position817, tokenIndex817, depth817
																if buffer[position] != rune('A') {
																	goto l700
																}
																position++
															}
														l817:
															{
																position819, tokenIndex819, depth819 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l820
																}
																position++
																goto l819
															l820:
																position, tokenIndex, depth = position819, tokenIndex819, depth819
																if buffer[position] != rune('X') {
																	goto l700
																}
																position++
															}
														l819:
															if !rules[ruleskip]() {
																goto l700
															}
															depth--
															add(ruleMAX, position814)
														}
														break
													}
												}

											}
										l776:
											if !rules[ruleLPAREN]() {
												goto l700
											}
											{
												position821, tokenIndex821, depth821 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l821
												}
												goto l822
											l821:
												position, tokenIndex, depth = position821, tokenIndex821, depth821
											}
										l822:
											if !rules[ruleexpression]() {
												goto l700
											}
											if !rules[ruleRPAREN]() {
												goto l700
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position711)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l700
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l700
								}
								break
							case '(':
								if !rules[rulebrackettedExpression]() {
									goto l700
								}
								break
							case '"', '\'':
								if !rules[ruleliteral]() {
									goto l700
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l700
								}
								break
							}
						}

					}
				l706:
					depth--
					add(ruleprimaryExpression, position705)
				}
				depth--
				add(ruleunaryExpression, position701)
			}
			return true
		l700:
			position, tokenIndex, depth = position700, tokenIndex700, depth700
			return false
		},
		/* 65 primaryExpression <- <(builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('(') brackettedExpression) | (&('"' | '\'') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 66 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position824, tokenIndex824, depth824 := position, tokenIndex, depth
			{
				position825 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l824
				}
				if !rules[ruleexpression]() {
					goto l824
				}
				if !rules[ruleRPAREN]() {
					goto l824
				}
				depth--
				add(rulebrackettedExpression, position825)
			}
			return true
		l824:
			position, tokenIndex, depth = position824, tokenIndex824, depth824
			return false
		},
		/* 67 functionCall <- <(iriref argList)> */
		func() bool {
			position826, tokenIndex826, depth826 := position, tokenIndex, depth
			{
				position827 := position
				depth++
				if !rules[ruleiriref]() {
					goto l826
				}
				if !rules[ruleargList]() {
					goto l826
				}
				depth--
				add(rulefunctionCall, position827)
			}
			return true
		l826:
			position, tokenIndex, depth = position826, tokenIndex826, depth826
			return false
		},
		/* 68 in <- <(IN argList)> */
		nil,
		/* 69 notin <- <(NOTIN argList)> */
		nil,
		/* 70 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position830, tokenIndex830, depth830 := position, tokenIndex, depth
			{
				position831 := position
				depth++
				{
					position832, tokenIndex832, depth832 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l833
					}
					goto l832
				l833:
					position, tokenIndex, depth = position832, tokenIndex832, depth832
					if !rules[ruleLPAREN]() {
						goto l830
					}
					if !rules[ruleexpression]() {
						goto l830
					}
				l834:
					{
						position835, tokenIndex835, depth835 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l835
						}
						if !rules[ruleexpression]() {
							goto l835
						}
						goto l834
					l835:
						position, tokenIndex, depth = position835, tokenIndex835, depth835
					}
					if !rules[ruleRPAREN]() {
						goto l830
					}
				}
			l832:
				depth--
				add(ruleargList, position831)
			}
			return true
		l830:
			position, tokenIndex, depth = position830, tokenIndex830, depth830
			return false
		},
		/* 71 aggregate <- <((&('G' | 'g') groupConcat) | (&('C' | 'c') count) | (&('A' | 'M' | 'S' | 'a' | 'm' | 's') ((SUM / MIN / ((&('S' | 's') SAMPLE) | (&('A' | 'a') AVG) | (&('M' | 'm') MAX))) LPAREN DISTINCT? expression RPAREN)))> */
		nil,
		/* 72 count <- <(COUNT LPAREN DISTINCT? (STAR / expression) RPAREN)> */
		nil,
		/* 73 groupConcat <- <(GROUPCONCAT LPAREN DISTINCT? expression (SEMICOLON SEPARATOR EQ string)? RPAREN)> */
		nil,
		/* 74 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		func() bool {
			position839, tokenIndex839, depth839 := position, tokenIndex, depth
			{
				position840 := position
				depth++
				{
					position841, tokenIndex841, depth841 := position, tokenIndex, depth
					{
						position843, tokenIndex843, depth843 := position, tokenIndex, depth
						{
							position845 := position
							depth++
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
									goto l844
								}
								position++
							}
						l846:
							{
								position848, tokenIndex848, depth848 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l849
								}
								position++
								goto l848
							l849:
								position, tokenIndex, depth = position848, tokenIndex848, depth848
								if buffer[position] != rune('T') {
									goto l844
								}
								position++
							}
						l848:
							{
								position850, tokenIndex850, depth850 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l851
								}
								position++
								goto l850
							l851:
								position, tokenIndex, depth = position850, tokenIndex850, depth850
								if buffer[position] != rune('R') {
									goto l844
								}
								position++
							}
						l850:
							if !rules[ruleskip]() {
								goto l844
							}
							depth--
							add(ruleSTR, position845)
						}
						goto l843
					l844:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position853 := position
							depth++
							{
								position854, tokenIndex854, depth854 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l855
								}
								position++
								goto l854
							l855:
								position, tokenIndex, depth = position854, tokenIndex854, depth854
								if buffer[position] != rune('L') {
									goto l852
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
									goto l852
								}
								position++
							}
						l856:
							{
								position858, tokenIndex858, depth858 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('N') {
									goto l852
								}
								position++
							}
						l858:
							{
								position860, tokenIndex860, depth860 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l861
								}
								position++
								goto l860
							l861:
								position, tokenIndex, depth = position860, tokenIndex860, depth860
								if buffer[position] != rune('G') {
									goto l852
								}
								position++
							}
						l860:
							if !rules[ruleskip]() {
								goto l852
							}
							depth--
							add(ruleLANG, position853)
						}
						goto l843
					l852:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position863 := position
							depth++
							{
								position864, tokenIndex864, depth864 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l865
								}
								position++
								goto l864
							l865:
								position, tokenIndex, depth = position864, tokenIndex864, depth864
								if buffer[position] != rune('D') {
									goto l862
								}
								position++
							}
						l864:
							{
								position866, tokenIndex866, depth866 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l867
								}
								position++
								goto l866
							l867:
								position, tokenIndex, depth = position866, tokenIndex866, depth866
								if buffer[position] != rune('A') {
									goto l862
								}
								position++
							}
						l866:
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
									goto l862
								}
								position++
							}
						l868:
							{
								position870, tokenIndex870, depth870 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l871
								}
								position++
								goto l870
							l871:
								position, tokenIndex, depth = position870, tokenIndex870, depth870
								if buffer[position] != rune('A') {
									goto l862
								}
								position++
							}
						l870:
							{
								position872, tokenIndex872, depth872 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l873
								}
								position++
								goto l872
							l873:
								position, tokenIndex, depth = position872, tokenIndex872, depth872
								if buffer[position] != rune('T') {
									goto l862
								}
								position++
							}
						l872:
							{
								position874, tokenIndex874, depth874 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l875
								}
								position++
								goto l874
							l875:
								position, tokenIndex, depth = position874, tokenIndex874, depth874
								if buffer[position] != rune('Y') {
									goto l862
								}
								position++
							}
						l874:
							{
								position876, tokenIndex876, depth876 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l877
								}
								position++
								goto l876
							l877:
								position, tokenIndex, depth = position876, tokenIndex876, depth876
								if buffer[position] != rune('P') {
									goto l862
								}
								position++
							}
						l876:
							{
								position878, tokenIndex878, depth878 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l879
								}
								position++
								goto l878
							l879:
								position, tokenIndex, depth = position878, tokenIndex878, depth878
								if buffer[position] != rune('E') {
									goto l862
								}
								position++
							}
						l878:
							if !rules[ruleskip]() {
								goto l862
							}
							depth--
							add(ruleDATATYPE, position863)
						}
						goto l843
					l862:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position881 := position
							depth++
							{
								position882, tokenIndex882, depth882 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l883
								}
								position++
								goto l882
							l883:
								position, tokenIndex, depth = position882, tokenIndex882, depth882
								if buffer[position] != rune('I') {
									goto l880
								}
								position++
							}
						l882:
							{
								position884, tokenIndex884, depth884 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l885
								}
								position++
								goto l884
							l885:
								position, tokenIndex, depth = position884, tokenIndex884, depth884
								if buffer[position] != rune('R') {
									goto l880
								}
								position++
							}
						l884:
							{
								position886, tokenIndex886, depth886 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l887
								}
								position++
								goto l886
							l887:
								position, tokenIndex, depth = position886, tokenIndex886, depth886
								if buffer[position] != rune('I') {
									goto l880
								}
								position++
							}
						l886:
							if !rules[ruleskip]() {
								goto l880
							}
							depth--
							add(ruleIRI, position881)
						}
						goto l843
					l880:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
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
									goto l888
								}
								position++
							}
						l890:
							{
								position892, tokenIndex892, depth892 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l893
								}
								position++
								goto l892
							l893:
								position, tokenIndex, depth = position892, tokenIndex892, depth892
								if buffer[position] != rune('R') {
									goto l888
								}
								position++
							}
						l892:
							{
								position894, tokenIndex894, depth894 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l895
								}
								position++
								goto l894
							l895:
								position, tokenIndex, depth = position894, tokenIndex894, depth894
								if buffer[position] != rune('I') {
									goto l888
								}
								position++
							}
						l894:
							if !rules[ruleskip]() {
								goto l888
							}
							depth--
							add(ruleURI, position889)
						}
						goto l843
					l888:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position897 := position
							depth++
							{
								position898, tokenIndex898, depth898 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l899
								}
								position++
								goto l898
							l899:
								position, tokenIndex, depth = position898, tokenIndex898, depth898
								if buffer[position] != rune('S') {
									goto l896
								}
								position++
							}
						l898:
							{
								position900, tokenIndex900, depth900 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l901
								}
								position++
								goto l900
							l901:
								position, tokenIndex, depth = position900, tokenIndex900, depth900
								if buffer[position] != rune('T') {
									goto l896
								}
								position++
							}
						l900:
							{
								position902, tokenIndex902, depth902 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l903
								}
								position++
								goto l902
							l903:
								position, tokenIndex, depth = position902, tokenIndex902, depth902
								if buffer[position] != rune('R') {
									goto l896
								}
								position++
							}
						l902:
							{
								position904, tokenIndex904, depth904 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l905
								}
								position++
								goto l904
							l905:
								position, tokenIndex, depth = position904, tokenIndex904, depth904
								if buffer[position] != rune('L') {
									goto l896
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
									goto l896
								}
								position++
							}
						l906:
							{
								position908, tokenIndex908, depth908 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l909
								}
								position++
								goto l908
							l909:
								position, tokenIndex, depth = position908, tokenIndex908, depth908
								if buffer[position] != rune('N') {
									goto l896
								}
								position++
							}
						l908:
							if !rules[ruleskip]() {
								goto l896
							}
							depth--
							add(ruleSTRLEN, position897)
						}
						goto l843
					l896:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position911 := position
							depth++
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
									goto l910
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
									goto l910
								}
								position++
							}
						l914:
							{
								position916, tokenIndex916, depth916 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l917
								}
								position++
								goto l916
							l917:
								position, tokenIndex, depth = position916, tokenIndex916, depth916
								if buffer[position] != rune('N') {
									goto l910
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
									goto l910
								}
								position++
							}
						l918:
							{
								position920, tokenIndex920, depth920 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l921
								}
								position++
								goto l920
							l921:
								position, tokenIndex, depth = position920, tokenIndex920, depth920
								if buffer[position] != rune('H') {
									goto l910
								}
								position++
							}
						l920:
							if !rules[ruleskip]() {
								goto l910
							}
							depth--
							add(ruleMONTH, position911)
						}
						goto l843
					l910:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position923 := position
							depth++
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
									goto l922
								}
								position++
							}
						l924:
							{
								position926, tokenIndex926, depth926 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l927
								}
								position++
								goto l926
							l927:
								position, tokenIndex, depth = position926, tokenIndex926, depth926
								if buffer[position] != rune('I') {
									goto l922
								}
								position++
							}
						l926:
							{
								position928, tokenIndex928, depth928 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l929
								}
								position++
								goto l928
							l929:
								position, tokenIndex, depth = position928, tokenIndex928, depth928
								if buffer[position] != rune('N') {
									goto l922
								}
								position++
							}
						l928:
							{
								position930, tokenIndex930, depth930 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l931
								}
								position++
								goto l930
							l931:
								position, tokenIndex, depth = position930, tokenIndex930, depth930
								if buffer[position] != rune('U') {
									goto l922
								}
								position++
							}
						l930:
							{
								position932, tokenIndex932, depth932 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l933
								}
								position++
								goto l932
							l933:
								position, tokenIndex, depth = position932, tokenIndex932, depth932
								if buffer[position] != rune('T') {
									goto l922
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
									goto l922
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
									goto l922
								}
								position++
							}
						l936:
							if !rules[ruleskip]() {
								goto l922
							}
							depth--
							add(ruleMINUTES, position923)
						}
						goto l843
					l922:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position939 := position
							depth++
							{
								position940, tokenIndex940, depth940 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l941
								}
								position++
								goto l940
							l941:
								position, tokenIndex, depth = position940, tokenIndex940, depth940
								if buffer[position] != rune('S') {
									goto l938
								}
								position++
							}
						l940:
							{
								position942, tokenIndex942, depth942 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l943
								}
								position++
								goto l942
							l943:
								position, tokenIndex, depth = position942, tokenIndex942, depth942
								if buffer[position] != rune('E') {
									goto l938
								}
								position++
							}
						l942:
							{
								position944, tokenIndex944, depth944 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l945
								}
								position++
								goto l944
							l945:
								position, tokenIndex, depth = position944, tokenIndex944, depth944
								if buffer[position] != rune('C') {
									goto l938
								}
								position++
							}
						l944:
							{
								position946, tokenIndex946, depth946 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l947
								}
								position++
								goto l946
							l947:
								position, tokenIndex, depth = position946, tokenIndex946, depth946
								if buffer[position] != rune('O') {
									goto l938
								}
								position++
							}
						l946:
							{
								position948, tokenIndex948, depth948 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l949
								}
								position++
								goto l948
							l949:
								position, tokenIndex, depth = position948, tokenIndex948, depth948
								if buffer[position] != rune('N') {
									goto l938
								}
								position++
							}
						l948:
							{
								position950, tokenIndex950, depth950 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l951
								}
								position++
								goto l950
							l951:
								position, tokenIndex, depth = position950, tokenIndex950, depth950
								if buffer[position] != rune('D') {
									goto l938
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
									goto l938
								}
								position++
							}
						l952:
							if !rules[ruleskip]() {
								goto l938
							}
							depth--
							add(ruleSECONDS, position939)
						}
						goto l843
					l938:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position955 := position
							depth++
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
									goto l954
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
									goto l954
								}
								position++
							}
						l958:
							{
								position960, tokenIndex960, depth960 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l961
								}
								position++
								goto l960
							l961:
								position, tokenIndex, depth = position960, tokenIndex960, depth960
								if buffer[position] != rune('M') {
									goto l954
								}
								position++
							}
						l960:
							{
								position962, tokenIndex962, depth962 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l963
								}
								position++
								goto l962
							l963:
								position, tokenIndex, depth = position962, tokenIndex962, depth962
								if buffer[position] != rune('E') {
									goto l954
								}
								position++
							}
						l962:
							{
								position964, tokenIndex964, depth964 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l965
								}
								position++
								goto l964
							l965:
								position, tokenIndex, depth = position964, tokenIndex964, depth964
								if buffer[position] != rune('Z') {
									goto l954
								}
								position++
							}
						l964:
							{
								position966, tokenIndex966, depth966 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l967
								}
								position++
								goto l966
							l967:
								position, tokenIndex, depth = position966, tokenIndex966, depth966
								if buffer[position] != rune('O') {
									goto l954
								}
								position++
							}
						l966:
							{
								position968, tokenIndex968, depth968 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l969
								}
								position++
								goto l968
							l969:
								position, tokenIndex, depth = position968, tokenIndex968, depth968
								if buffer[position] != rune('N') {
									goto l954
								}
								position++
							}
						l968:
							{
								position970, tokenIndex970, depth970 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l971
								}
								position++
								goto l970
							l971:
								position, tokenIndex, depth = position970, tokenIndex970, depth970
								if buffer[position] != rune('E') {
									goto l954
								}
								position++
							}
						l970:
							if !rules[ruleskip]() {
								goto l954
							}
							depth--
							add(ruleTIMEZONE, position955)
						}
						goto l843
					l954:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position973 := position
							depth++
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
									goto l972
								}
								position++
							}
						l974:
							{
								position976, tokenIndex976, depth976 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l977
								}
								position++
								goto l976
							l977:
								position, tokenIndex, depth = position976, tokenIndex976, depth976
								if buffer[position] != rune('H') {
									goto l972
								}
								position++
							}
						l976:
							{
								position978, tokenIndex978, depth978 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l979
								}
								position++
								goto l978
							l979:
								position, tokenIndex, depth = position978, tokenIndex978, depth978
								if buffer[position] != rune('A') {
									goto l972
								}
								position++
							}
						l978:
							if buffer[position] != rune('1') {
								goto l972
							}
							position++
							if !rules[ruleskip]() {
								goto l972
							}
							depth--
							add(ruleSHA1, position973)
						}
						goto l843
					l972:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position981 := position
							depth++
							{
								position982, tokenIndex982, depth982 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l983
								}
								position++
								goto l982
							l983:
								position, tokenIndex, depth = position982, tokenIndex982, depth982
								if buffer[position] != rune('S') {
									goto l980
								}
								position++
							}
						l982:
							{
								position984, tokenIndex984, depth984 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l985
								}
								position++
								goto l984
							l985:
								position, tokenIndex, depth = position984, tokenIndex984, depth984
								if buffer[position] != rune('H') {
									goto l980
								}
								position++
							}
						l984:
							{
								position986, tokenIndex986, depth986 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l987
								}
								position++
								goto l986
							l987:
								position, tokenIndex, depth = position986, tokenIndex986, depth986
								if buffer[position] != rune('A') {
									goto l980
								}
								position++
							}
						l986:
							if buffer[position] != rune('2') {
								goto l980
							}
							position++
							if buffer[position] != rune('5') {
								goto l980
							}
							position++
							if buffer[position] != rune('6') {
								goto l980
							}
							position++
							if !rules[ruleskip]() {
								goto l980
							}
							depth--
							add(ruleSHA256, position981)
						}
						goto l843
					l980:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position989 := position
							depth++
							{
								position990, tokenIndex990, depth990 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l991
								}
								position++
								goto l990
							l991:
								position, tokenIndex, depth = position990, tokenIndex990, depth990
								if buffer[position] != rune('S') {
									goto l988
								}
								position++
							}
						l990:
							{
								position992, tokenIndex992, depth992 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l993
								}
								position++
								goto l992
							l993:
								position, tokenIndex, depth = position992, tokenIndex992, depth992
								if buffer[position] != rune('H') {
									goto l988
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
									goto l988
								}
								position++
							}
						l994:
							if buffer[position] != rune('3') {
								goto l988
							}
							position++
							if buffer[position] != rune('8') {
								goto l988
							}
							position++
							if buffer[position] != rune('4') {
								goto l988
							}
							position++
							if !rules[ruleskip]() {
								goto l988
							}
							depth--
							add(ruleSHA384, position989)
						}
						goto l843
					l988:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position997 := position
							depth++
							{
								position998, tokenIndex998, depth998 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l999
								}
								position++
								goto l998
							l999:
								position, tokenIndex, depth = position998, tokenIndex998, depth998
								if buffer[position] != rune('I') {
									goto l996
								}
								position++
							}
						l998:
							{
								position1000, tokenIndex1000, depth1000 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1001
								}
								position++
								goto l1000
							l1001:
								position, tokenIndex, depth = position1000, tokenIndex1000, depth1000
								if buffer[position] != rune('S') {
									goto l996
								}
								position++
							}
						l1000:
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
									goto l996
								}
								position++
							}
						l1002:
							{
								position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1005
								}
								position++
								goto l1004
							l1005:
								position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
								if buffer[position] != rune('R') {
									goto l996
								}
								position++
							}
						l1004:
							{
								position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1007
								}
								position++
								goto l1006
							l1007:
								position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
								if buffer[position] != rune('I') {
									goto l996
								}
								position++
							}
						l1006:
							if !rules[ruleskip]() {
								goto l996
							}
							depth--
							add(ruleISIRI, position997)
						}
						goto l843
					l996:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
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
									goto l1008
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
									goto l1008
								}
								position++
							}
						l1012:
							{
								position1014, tokenIndex1014, depth1014 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l1015
								}
								position++
								goto l1014
							l1015:
								position, tokenIndex, depth = position1014, tokenIndex1014, depth1014
								if buffer[position] != rune('U') {
									goto l1008
								}
								position++
							}
						l1014:
							{
								position1016, tokenIndex1016, depth1016 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1017
								}
								position++
								goto l1016
							l1017:
								position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
								if buffer[position] != rune('R') {
									goto l1008
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
									goto l1008
								}
								position++
							}
						l1018:
							if !rules[ruleskip]() {
								goto l1008
							}
							depth--
							add(ruleISURI, position1009)
						}
						goto l843
					l1008:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
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
									goto l1020
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
									goto l1020
								}
								position++
							}
						l1024:
							{
								position1026, tokenIndex1026, depth1026 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1027
								}
								position++
								goto l1026
							l1027:
								position, tokenIndex, depth = position1026, tokenIndex1026, depth1026
								if buffer[position] != rune('B') {
									goto l1020
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
									goto l1020
								}
								position++
							}
						l1028:
							{
								position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1031
								}
								position++
								goto l1030
							l1031:
								position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
								if buffer[position] != rune('A') {
									goto l1020
								}
								position++
							}
						l1030:
							{
								position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1033
								}
								position++
								goto l1032
							l1033:
								position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
								if buffer[position] != rune('N') {
									goto l1020
								}
								position++
							}
						l1032:
							{
								position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l1035
								}
								position++
								goto l1034
							l1035:
								position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
								if buffer[position] != rune('K') {
									goto l1020
								}
								position++
							}
						l1034:
							if !rules[ruleskip]() {
								goto l1020
							}
							depth--
							add(ruleISBLANK, position1021)
						}
						goto l843
					l1020:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							position1037 := position
							depth++
							{
								position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1039
								}
								position++
								goto l1038
							l1039:
								position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
								if buffer[position] != rune('I') {
									goto l1036
								}
								position++
							}
						l1038:
							{
								position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1041
								}
								position++
								goto l1040
							l1041:
								position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
								if buffer[position] != rune('S') {
									goto l1036
								}
								position++
							}
						l1040:
							{
								position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1043
								}
								position++
								goto l1042
							l1043:
								position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
								if buffer[position] != rune('L') {
									goto l1036
								}
								position++
							}
						l1042:
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
								if buffer[position] != rune('e') {
									goto l1049
								}
								position++
								goto l1048
							l1049:
								position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
								if buffer[position] != rune('E') {
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
								if buffer[position] != rune('a') {
									goto l1053
								}
								position++
								goto l1052
							l1053:
								position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
								if buffer[position] != rune('A') {
									goto l1036
								}
								position++
							}
						l1052:
							{
								position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1055
								}
								position++
								goto l1054
							l1055:
								position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
								if buffer[position] != rune('L') {
									goto l1036
								}
								position++
							}
						l1054:
							if !rules[ruleskip]() {
								goto l1036
							}
							depth--
							add(ruleISLITERAL, position1037)
						}
						goto l843
					l1036:
						position, tokenIndex, depth = position843, tokenIndex843, depth843
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position1057 := position
									depth++
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
											goto l842
										}
										position++
									}
								l1058:
									{
										position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1061
										}
										position++
										goto l1060
									l1061:
										position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
										if buffer[position] != rune('S') {
											goto l842
										}
										position++
									}
								l1060:
									{
										position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1063
										}
										position++
										goto l1062
									l1063:
										position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
										if buffer[position] != rune('N') {
											goto l842
										}
										position++
									}
								l1062:
									{
										position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1065
										}
										position++
										goto l1064
									l1065:
										position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
										if buffer[position] != rune('U') {
											goto l842
										}
										position++
									}
								l1064:
									{
										position1066, tokenIndex1066, depth1066 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1067
										}
										position++
										goto l1066
									l1067:
										position, tokenIndex, depth = position1066, tokenIndex1066, depth1066
										if buffer[position] != rune('M') {
											goto l842
										}
										position++
									}
								l1066:
									{
										position1068, tokenIndex1068, depth1068 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1069
										}
										position++
										goto l1068
									l1069:
										position, tokenIndex, depth = position1068, tokenIndex1068, depth1068
										if buffer[position] != rune('E') {
											goto l842
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
											goto l842
										}
										position++
									}
								l1070:
									{
										position1072, tokenIndex1072, depth1072 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1073
										}
										position++
										goto l1072
									l1073:
										position, tokenIndex, depth = position1072, tokenIndex1072, depth1072
										if buffer[position] != rune('I') {
											goto l842
										}
										position++
									}
								l1072:
									{
										position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1075
										}
										position++
										goto l1074
									l1075:
										position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
										if buffer[position] != rune('C') {
											goto l842
										}
										position++
									}
								l1074:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleISNUMERIC, position1057)
								}
								break
							case 'S', 's':
								{
									position1076 := position
									depth++
									{
										position1077, tokenIndex1077, depth1077 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1078
										}
										position++
										goto l1077
									l1078:
										position, tokenIndex, depth = position1077, tokenIndex1077, depth1077
										if buffer[position] != rune('S') {
											goto l842
										}
										position++
									}
								l1077:
									{
										position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1080
										}
										position++
										goto l1079
									l1080:
										position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
										if buffer[position] != rune('H') {
											goto l842
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
											goto l842
										}
										position++
									}
								l1081:
									if buffer[position] != rune('5') {
										goto l842
									}
									position++
									if buffer[position] != rune('1') {
										goto l842
									}
									position++
									if buffer[position] != rune('2') {
										goto l842
									}
									position++
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleSHA512, position1076)
								}
								break
							case 'M', 'm':
								{
									position1083 := position
									depth++
									{
										position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1085
										}
										position++
										goto l1084
									l1085:
										position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
										if buffer[position] != rune('M') {
											goto l842
										}
										position++
									}
								l1084:
									{
										position1086, tokenIndex1086, depth1086 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1087
										}
										position++
										goto l1086
									l1087:
										position, tokenIndex, depth = position1086, tokenIndex1086, depth1086
										if buffer[position] != rune('D') {
											goto l842
										}
										position++
									}
								l1086:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleMD5, position1083)
								}
								break
							case 'T', 't':
								{
									position1088 := position
									depth++
									{
										position1089, tokenIndex1089, depth1089 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1090
										}
										position++
										goto l1089
									l1090:
										position, tokenIndex, depth = position1089, tokenIndex1089, depth1089
										if buffer[position] != rune('T') {
											goto l842
										}
										position++
									}
								l1089:
									{
										position1091, tokenIndex1091, depth1091 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1092
										}
										position++
										goto l1091
									l1092:
										position, tokenIndex, depth = position1091, tokenIndex1091, depth1091
										if buffer[position] != rune('Z') {
											goto l842
										}
										position++
									}
								l1091:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleTZ, position1088)
								}
								break
							case 'H', 'h':
								{
									position1093 := position
									depth++
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
											goto l842
										}
										position++
									}
								l1094:
									{
										position1096, tokenIndex1096, depth1096 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1097
										}
										position++
										goto l1096
									l1097:
										position, tokenIndex, depth = position1096, tokenIndex1096, depth1096
										if buffer[position] != rune('O') {
											goto l842
										}
										position++
									}
								l1096:
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
											goto l842
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
											goto l842
										}
										position++
									}
								l1100:
									{
										position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1103
										}
										position++
										goto l1102
									l1103:
										position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
										if buffer[position] != rune('S') {
											goto l842
										}
										position++
									}
								l1102:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleHOURS, position1093)
								}
								break
							case 'D', 'd':
								{
									position1104 := position
									depth++
									{
										position1105, tokenIndex1105, depth1105 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1106
										}
										position++
										goto l1105
									l1106:
										position, tokenIndex, depth = position1105, tokenIndex1105, depth1105
										if buffer[position] != rune('D') {
											goto l842
										}
										position++
									}
								l1105:
									{
										position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1108
										}
										position++
										goto l1107
									l1108:
										position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
										if buffer[position] != rune('A') {
											goto l842
										}
										position++
									}
								l1107:
									{
										position1109, tokenIndex1109, depth1109 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1110
										}
										position++
										goto l1109
									l1110:
										position, tokenIndex, depth = position1109, tokenIndex1109, depth1109
										if buffer[position] != rune('Y') {
											goto l842
										}
										position++
									}
								l1109:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleDAY, position1104)
								}
								break
							case 'Y', 'y':
								{
									position1111 := position
									depth++
									{
										position1112, tokenIndex1112, depth1112 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1113
										}
										position++
										goto l1112
									l1113:
										position, tokenIndex, depth = position1112, tokenIndex1112, depth1112
										if buffer[position] != rune('Y') {
											goto l842
										}
										position++
									}
								l1112:
									{
										position1114, tokenIndex1114, depth1114 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1115
										}
										position++
										goto l1114
									l1115:
										position, tokenIndex, depth = position1114, tokenIndex1114, depth1114
										if buffer[position] != rune('E') {
											goto l842
										}
										position++
									}
								l1114:
									{
										position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1117
										}
										position++
										goto l1116
									l1117:
										position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
										if buffer[position] != rune('A') {
											goto l842
										}
										position++
									}
								l1116:
									{
										position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1119
										}
										position++
										goto l1118
									l1119:
										position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
										if buffer[position] != rune('R') {
											goto l842
										}
										position++
									}
								l1118:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleYEAR, position1111)
								}
								break
							case 'E', 'e':
								{
									position1120 := position
									depth++
									{
										position1121, tokenIndex1121, depth1121 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1122
										}
										position++
										goto l1121
									l1122:
										position, tokenIndex, depth = position1121, tokenIndex1121, depth1121
										if buffer[position] != rune('E') {
											goto l842
										}
										position++
									}
								l1121:
									{
										position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1124
										}
										position++
										goto l1123
									l1124:
										position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
										if buffer[position] != rune('N') {
											goto l842
										}
										position++
									}
								l1123:
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
											goto l842
										}
										position++
									}
								l1125:
									{
										position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1128
										}
										position++
										goto l1127
									l1128:
										position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
										if buffer[position] != rune('O') {
											goto l842
										}
										position++
									}
								l1127:
									{
										position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1130
										}
										position++
										goto l1129
									l1130:
										position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
										if buffer[position] != rune('D') {
											goto l842
										}
										position++
									}
								l1129:
									{
										position1131, tokenIndex1131, depth1131 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1132
										}
										position++
										goto l1131
									l1132:
										position, tokenIndex, depth = position1131, tokenIndex1131, depth1131
										if buffer[position] != rune('E') {
											goto l842
										}
										position++
									}
								l1131:
									if buffer[position] != rune('_') {
										goto l842
									}
									position++
									{
										position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1134
										}
										position++
										goto l1133
									l1134:
										position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
										if buffer[position] != rune('F') {
											goto l842
										}
										position++
									}
								l1133:
									{
										position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1136
										}
										position++
										goto l1135
									l1136:
										position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
										if buffer[position] != rune('O') {
											goto l842
										}
										position++
									}
								l1135:
									{
										position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1138
										}
										position++
										goto l1137
									l1138:
										position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
										if buffer[position] != rune('R') {
											goto l842
										}
										position++
									}
								l1137:
									if buffer[position] != rune('_') {
										goto l842
									}
									position++
									{
										position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1140
										}
										position++
										goto l1139
									l1140:
										position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
										if buffer[position] != rune('U') {
											goto l842
										}
										position++
									}
								l1139:
									{
										position1141, tokenIndex1141, depth1141 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1142
										}
										position++
										goto l1141
									l1142:
										position, tokenIndex, depth = position1141, tokenIndex1141, depth1141
										if buffer[position] != rune('R') {
											goto l842
										}
										position++
									}
								l1141:
									{
										position1143, tokenIndex1143, depth1143 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1144
										}
										position++
										goto l1143
									l1144:
										position, tokenIndex, depth = position1143, tokenIndex1143, depth1143
										if buffer[position] != rune('I') {
											goto l842
										}
										position++
									}
								l1143:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleENCODEFORURI, position1120)
								}
								break
							case 'L', 'l':
								{
									position1145 := position
									depth++
									{
										position1146, tokenIndex1146, depth1146 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1147
										}
										position++
										goto l1146
									l1147:
										position, tokenIndex, depth = position1146, tokenIndex1146, depth1146
										if buffer[position] != rune('L') {
											goto l842
										}
										position++
									}
								l1146:
									{
										position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1149
										}
										position++
										goto l1148
									l1149:
										position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
										if buffer[position] != rune('C') {
											goto l842
										}
										position++
									}
								l1148:
									{
										position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1151
										}
										position++
										goto l1150
									l1151:
										position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
										if buffer[position] != rune('A') {
											goto l842
										}
										position++
									}
								l1150:
									{
										position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1153
										}
										position++
										goto l1152
									l1153:
										position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
										if buffer[position] != rune('S') {
											goto l842
										}
										position++
									}
								l1152:
									{
										position1154, tokenIndex1154, depth1154 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1155
										}
										position++
										goto l1154
									l1155:
										position, tokenIndex, depth = position1154, tokenIndex1154, depth1154
										if buffer[position] != rune('E') {
											goto l842
										}
										position++
									}
								l1154:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleLCASE, position1145)
								}
								break
							case 'U', 'u':
								{
									position1156 := position
									depth++
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
											goto l842
										}
										position++
									}
								l1157:
									{
										position1159, tokenIndex1159, depth1159 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1160
										}
										position++
										goto l1159
									l1160:
										position, tokenIndex, depth = position1159, tokenIndex1159, depth1159
										if buffer[position] != rune('C') {
											goto l842
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
											goto l842
										}
										position++
									}
								l1161:
									{
										position1163, tokenIndex1163, depth1163 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1164
										}
										position++
										goto l1163
									l1164:
										position, tokenIndex, depth = position1163, tokenIndex1163, depth1163
										if buffer[position] != rune('S') {
											goto l842
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
											goto l842
										}
										position++
									}
								l1165:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleUCASE, position1156)
								}
								break
							case 'F', 'f':
								{
									position1167 := position
									depth++
									{
										position1168, tokenIndex1168, depth1168 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1169
										}
										position++
										goto l1168
									l1169:
										position, tokenIndex, depth = position1168, tokenIndex1168, depth1168
										if buffer[position] != rune('F') {
											goto l842
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
											goto l842
										}
										position++
									}
								l1170:
									{
										position1172, tokenIndex1172, depth1172 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1173
										}
										position++
										goto l1172
									l1173:
										position, tokenIndex, depth = position1172, tokenIndex1172, depth1172
										if buffer[position] != rune('O') {
											goto l842
										}
										position++
									}
								l1172:
									{
										position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1175
										}
										position++
										goto l1174
									l1175:
										position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
										if buffer[position] != rune('O') {
											goto l842
										}
										position++
									}
								l1174:
									{
										position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1177
										}
										position++
										goto l1176
									l1177:
										position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
										if buffer[position] != rune('R') {
											goto l842
										}
										position++
									}
								l1176:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleFLOOR, position1167)
								}
								break
							case 'R', 'r':
								{
									position1178 := position
									depth++
									{
										position1179, tokenIndex1179, depth1179 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1180
										}
										position++
										goto l1179
									l1180:
										position, tokenIndex, depth = position1179, tokenIndex1179, depth1179
										if buffer[position] != rune('R') {
											goto l842
										}
										position++
									}
								l1179:
									{
										position1181, tokenIndex1181, depth1181 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1182
										}
										position++
										goto l1181
									l1182:
										position, tokenIndex, depth = position1181, tokenIndex1181, depth1181
										if buffer[position] != rune('O') {
											goto l842
										}
										position++
									}
								l1181:
									{
										position1183, tokenIndex1183, depth1183 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1184
										}
										position++
										goto l1183
									l1184:
										position, tokenIndex, depth = position1183, tokenIndex1183, depth1183
										if buffer[position] != rune('U') {
											goto l842
										}
										position++
									}
								l1183:
									{
										position1185, tokenIndex1185, depth1185 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1186
										}
										position++
										goto l1185
									l1186:
										position, tokenIndex, depth = position1185, tokenIndex1185, depth1185
										if buffer[position] != rune('N') {
											goto l842
										}
										position++
									}
								l1185:
									{
										position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1188
										}
										position++
										goto l1187
									l1188:
										position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
										if buffer[position] != rune('D') {
											goto l842
										}
										position++
									}
								l1187:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleROUND, position1178)
								}
								break
							case 'C', 'c':
								{
									position1189 := position
									depth++
									{
										position1190, tokenIndex1190, depth1190 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1191
										}
										position++
										goto l1190
									l1191:
										position, tokenIndex, depth = position1190, tokenIndex1190, depth1190
										if buffer[position] != rune('C') {
											goto l842
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
											goto l842
										}
										position++
									}
								l1192:
									{
										position1194, tokenIndex1194, depth1194 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1195
										}
										position++
										goto l1194
									l1195:
										position, tokenIndex, depth = position1194, tokenIndex1194, depth1194
										if buffer[position] != rune('I') {
											goto l842
										}
										position++
									}
								l1194:
									{
										position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1197
										}
										position++
										goto l1196
									l1197:
										position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
										if buffer[position] != rune('L') {
											goto l842
										}
										position++
									}
								l1196:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleCEIL, position1189)
								}
								break
							default:
								{
									position1198 := position
									depth++
									{
										position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1200
										}
										position++
										goto l1199
									l1200:
										position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
										if buffer[position] != rune('A') {
											goto l842
										}
										position++
									}
								l1199:
									{
										position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1202
										}
										position++
										goto l1201
									l1202:
										position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
										if buffer[position] != rune('B') {
											goto l842
										}
										position++
									}
								l1201:
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
											goto l842
										}
										position++
									}
								l1203:
									if !rules[ruleskip]() {
										goto l842
									}
									depth--
									add(ruleABS, position1198)
								}
								break
							}
						}

					}
				l843:
					if !rules[ruleLPAREN]() {
						goto l842
					}
					if !rules[ruleexpression]() {
						goto l842
					}
					if !rules[ruleRPAREN]() {
						goto l842
					}
					goto l841
				l842:
					position, tokenIndex, depth = position841, tokenIndex841, depth841
					{
						position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
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
								if buffer[position] != rune('s') {
									goto l1216
								}
								position++
								goto l1215
							l1216:
								position, tokenIndex, depth = position1215, tokenIndex1215, depth1215
								if buffer[position] != rune('S') {
									goto l1207
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
									goto l1207
								}
								position++
							}
						l1217:
							{
								position1219, tokenIndex1219, depth1219 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1220
								}
								position++
								goto l1219
							l1220:
								position, tokenIndex, depth = position1219, tokenIndex1219, depth1219
								if buffer[position] != rune('A') {
									goto l1207
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
									goto l1207
								}
								position++
							}
						l1221:
							{
								position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1224
								}
								position++
								goto l1223
							l1224:
								position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
								if buffer[position] != rune('T') {
									goto l1207
								}
								position++
							}
						l1223:
							{
								position1225, tokenIndex1225, depth1225 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1226
								}
								position++
								goto l1225
							l1226:
								position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
								if buffer[position] != rune('S') {
									goto l1207
								}
								position++
							}
						l1225:
							if !rules[ruleskip]() {
								goto l1207
							}
							depth--
							add(ruleSTRSTARTS, position1208)
						}
						goto l1206
					l1207:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
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
								if buffer[position] != rune('e') {
									goto l1236
								}
								position++
								goto l1235
							l1236:
								position, tokenIndex, depth = position1235, tokenIndex1235, depth1235
								if buffer[position] != rune('E') {
									goto l1227
								}
								position++
							}
						l1235:
							{
								position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1238
								}
								position++
								goto l1237
							l1238:
								position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
								if buffer[position] != rune('N') {
									goto l1227
								}
								position++
							}
						l1237:
							{
								position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1240
								}
								position++
								goto l1239
							l1240:
								position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
								if buffer[position] != rune('D') {
									goto l1227
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
									goto l1227
								}
								position++
							}
						l1241:
							if !rules[ruleskip]() {
								goto l1227
							}
							depth--
							add(ruleSTRENDS, position1228)
						}
						goto l1206
					l1227:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						{
							position1244 := position
							depth++
							{
								position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1246
								}
								position++
								goto l1245
							l1246:
								position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
								if buffer[position] != rune('S') {
									goto l1243
								}
								position++
							}
						l1245:
							{
								position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1248
								}
								position++
								goto l1247
							l1248:
								position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
								if buffer[position] != rune('T') {
									goto l1243
								}
								position++
							}
						l1247:
							{
								position1249, tokenIndex1249, depth1249 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1250
								}
								position++
								goto l1249
							l1250:
								position, tokenIndex, depth = position1249, tokenIndex1249, depth1249
								if buffer[position] != rune('R') {
									goto l1243
								}
								position++
							}
						l1249:
							{
								position1251, tokenIndex1251, depth1251 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l1252
								}
								position++
								goto l1251
							l1252:
								position, tokenIndex, depth = position1251, tokenIndex1251, depth1251
								if buffer[position] != rune('B') {
									goto l1243
								}
								position++
							}
						l1251:
							{
								position1253, tokenIndex1253, depth1253 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1254
								}
								position++
								goto l1253
							l1254:
								position, tokenIndex, depth = position1253, tokenIndex1253, depth1253
								if buffer[position] != rune('E') {
									goto l1243
								}
								position++
							}
						l1253:
							{
								position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1256
								}
								position++
								goto l1255
							l1256:
								position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
								if buffer[position] != rune('F') {
									goto l1243
								}
								position++
							}
						l1255:
							{
								position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1258
								}
								position++
								goto l1257
							l1258:
								position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
								if buffer[position] != rune('O') {
									goto l1243
								}
								position++
							}
						l1257:
							{
								position1259, tokenIndex1259, depth1259 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1260
								}
								position++
								goto l1259
							l1260:
								position, tokenIndex, depth = position1259, tokenIndex1259, depth1259
								if buffer[position] != rune('R') {
									goto l1243
								}
								position++
							}
						l1259:
							{
								position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1262
								}
								position++
								goto l1261
							l1262:
								position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
								if buffer[position] != rune('E') {
									goto l1243
								}
								position++
							}
						l1261:
							if !rules[ruleskip]() {
								goto l1243
							}
							depth--
							add(ruleSTRBEFORE, position1244)
						}
						goto l1206
					l1243:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						{
							position1264 := position
							depth++
							{
								position1265, tokenIndex1265, depth1265 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1266
								}
								position++
								goto l1265
							l1266:
								position, tokenIndex, depth = position1265, tokenIndex1265, depth1265
								if buffer[position] != rune('S') {
									goto l1263
								}
								position++
							}
						l1265:
							{
								position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1268
								}
								position++
								goto l1267
							l1268:
								position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
								if buffer[position] != rune('T') {
									goto l1263
								}
								position++
							}
						l1267:
							{
								position1269, tokenIndex1269, depth1269 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1270
								}
								position++
								goto l1269
							l1270:
								position, tokenIndex, depth = position1269, tokenIndex1269, depth1269
								if buffer[position] != rune('R') {
									goto l1263
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
									goto l1263
								}
								position++
							}
						l1271:
							{
								position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1274
								}
								position++
								goto l1273
							l1274:
								position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
								if buffer[position] != rune('F') {
									goto l1263
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
									goto l1263
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
									goto l1263
								}
								position++
							}
						l1277:
							{
								position1279, tokenIndex1279, depth1279 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1280
								}
								position++
								goto l1279
							l1280:
								position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
								if buffer[position] != rune('R') {
									goto l1263
								}
								position++
							}
						l1279:
							if !rules[ruleskip]() {
								goto l1263
							}
							depth--
							add(ruleSTRAFTER, position1264)
						}
						goto l1206
					l1263:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						{
							position1282 := position
							depth++
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
									goto l1281
								}
								position++
							}
						l1283:
							{
								position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1286
								}
								position++
								goto l1285
							l1286:
								position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
								if buffer[position] != rune('T') {
									goto l1281
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
									goto l1281
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
									goto l1281
								}
								position++
							}
						l1289:
							{
								position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1292
								}
								position++
								goto l1291
							l1292:
								position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
								if buffer[position] != rune('A') {
									goto l1281
								}
								position++
							}
						l1291:
							{
								position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1294
								}
								position++
								goto l1293
							l1294:
								position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
								if buffer[position] != rune('N') {
									goto l1281
								}
								position++
							}
						l1293:
							{
								position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1296
								}
								position++
								goto l1295
							l1296:
								position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
								if buffer[position] != rune('G') {
									goto l1281
								}
								position++
							}
						l1295:
							if !rules[ruleskip]() {
								goto l1281
							}
							depth--
							add(ruleSTRLANG, position1282)
						}
						goto l1206
					l1281:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						{
							position1298 := position
							depth++
							{
								position1299, tokenIndex1299, depth1299 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1300
								}
								position++
								goto l1299
							l1300:
								position, tokenIndex, depth = position1299, tokenIndex1299, depth1299
								if buffer[position] != rune('S') {
									goto l1297
								}
								position++
							}
						l1299:
							{
								position1301, tokenIndex1301, depth1301 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1302
								}
								position++
								goto l1301
							l1302:
								position, tokenIndex, depth = position1301, tokenIndex1301, depth1301
								if buffer[position] != rune('T') {
									goto l1297
								}
								position++
							}
						l1301:
							{
								position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1304
								}
								position++
								goto l1303
							l1304:
								position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
								if buffer[position] != rune('R') {
									goto l1297
								}
								position++
							}
						l1303:
							{
								position1305, tokenIndex1305, depth1305 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1306
								}
								position++
								goto l1305
							l1306:
								position, tokenIndex, depth = position1305, tokenIndex1305, depth1305
								if buffer[position] != rune('D') {
									goto l1297
								}
								position++
							}
						l1305:
							{
								position1307, tokenIndex1307, depth1307 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1308
								}
								position++
								goto l1307
							l1308:
								position, tokenIndex, depth = position1307, tokenIndex1307, depth1307
								if buffer[position] != rune('T') {
									goto l1297
								}
								position++
							}
						l1307:
							if !rules[ruleskip]() {
								goto l1297
							}
							depth--
							add(ruleSTRDT, position1298)
						}
						goto l1206
					l1297:
						position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1310 := position
									depth++
									{
										position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1312
										}
										position++
										goto l1311
									l1312:
										position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
										if buffer[position] != rune('S') {
											goto l1205
										}
										position++
									}
								l1311:
									{
										position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1314
										}
										position++
										goto l1313
									l1314:
										position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
										if buffer[position] != rune('A') {
											goto l1205
										}
										position++
									}
								l1313:
									{
										position1315, tokenIndex1315, depth1315 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1316
										}
										position++
										goto l1315
									l1316:
										position, tokenIndex, depth = position1315, tokenIndex1315, depth1315
										if buffer[position] != rune('M') {
											goto l1205
										}
										position++
									}
								l1315:
									{
										position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1318
										}
										position++
										goto l1317
									l1318:
										position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
										if buffer[position] != rune('E') {
											goto l1205
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
											goto l1205
										}
										position++
									}
								l1319:
									{
										position1321, tokenIndex1321, depth1321 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1322
										}
										position++
										goto l1321
									l1322:
										position, tokenIndex, depth = position1321, tokenIndex1321, depth1321
										if buffer[position] != rune('E') {
											goto l1205
										}
										position++
									}
								l1321:
									{
										position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1324
										}
										position++
										goto l1323
									l1324:
										position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
										if buffer[position] != rune('R') {
											goto l1205
										}
										position++
									}
								l1323:
									{
										position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1326
										}
										position++
										goto l1325
									l1326:
										position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
										if buffer[position] != rune('M') {
											goto l1205
										}
										position++
									}
								l1325:
									if !rules[ruleskip]() {
										goto l1205
									}
									depth--
									add(ruleSAMETERM, position1310)
								}
								break
							case 'C', 'c':
								{
									position1327 := position
									depth++
									{
										position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1329
										}
										position++
										goto l1328
									l1329:
										position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
										if buffer[position] != rune('C') {
											goto l1205
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
											goto l1205
										}
										position++
									}
								l1330:
									{
										position1332, tokenIndex1332, depth1332 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1333
										}
										position++
										goto l1332
									l1333:
										position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
										if buffer[position] != rune('N') {
											goto l1205
										}
										position++
									}
								l1332:
									{
										position1334, tokenIndex1334, depth1334 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1335
										}
										position++
										goto l1334
									l1335:
										position, tokenIndex, depth = position1334, tokenIndex1334, depth1334
										if buffer[position] != rune('T') {
											goto l1205
										}
										position++
									}
								l1334:
									{
										position1336, tokenIndex1336, depth1336 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1337
										}
										position++
										goto l1336
									l1337:
										position, tokenIndex, depth = position1336, tokenIndex1336, depth1336
										if buffer[position] != rune('A') {
											goto l1205
										}
										position++
									}
								l1336:
									{
										position1338, tokenIndex1338, depth1338 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1339
										}
										position++
										goto l1338
									l1339:
										position, tokenIndex, depth = position1338, tokenIndex1338, depth1338
										if buffer[position] != rune('I') {
											goto l1205
										}
										position++
									}
								l1338:
									{
										position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1341
										}
										position++
										goto l1340
									l1341:
										position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
										if buffer[position] != rune('N') {
											goto l1205
										}
										position++
									}
								l1340:
									{
										position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1343
										}
										position++
										goto l1342
									l1343:
										position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
										if buffer[position] != rune('S') {
											goto l1205
										}
										position++
									}
								l1342:
									if !rules[ruleskip]() {
										goto l1205
									}
									depth--
									add(ruleCONTAINS, position1327)
								}
								break
							default:
								{
									position1344 := position
									depth++
									{
										position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1346
										}
										position++
										goto l1345
									l1346:
										position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
										if buffer[position] != rune('L') {
											goto l1205
										}
										position++
									}
								l1345:
									{
										position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1348
										}
										position++
										goto l1347
									l1348:
										position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
										if buffer[position] != rune('A') {
											goto l1205
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
											goto l1205
										}
										position++
									}
								l1349:
									{
										position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1352
										}
										position++
										goto l1351
									l1352:
										position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
										if buffer[position] != rune('G') {
											goto l1205
										}
										position++
									}
								l1351:
									{
										position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1354
										}
										position++
										goto l1353
									l1354:
										position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
										if buffer[position] != rune('M') {
											goto l1205
										}
										position++
									}
								l1353:
									{
										position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1356
										}
										position++
										goto l1355
									l1356:
										position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
										if buffer[position] != rune('A') {
											goto l1205
										}
										position++
									}
								l1355:
									{
										position1357, tokenIndex1357, depth1357 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1358
										}
										position++
										goto l1357
									l1358:
										position, tokenIndex, depth = position1357, tokenIndex1357, depth1357
										if buffer[position] != rune('T') {
											goto l1205
										}
										position++
									}
								l1357:
									{
										position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1360
										}
										position++
										goto l1359
									l1360:
										position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
										if buffer[position] != rune('C') {
											goto l1205
										}
										position++
									}
								l1359:
									{
										position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1362
										}
										position++
										goto l1361
									l1362:
										position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
										if buffer[position] != rune('H') {
											goto l1205
										}
										position++
									}
								l1361:
									{
										position1363, tokenIndex1363, depth1363 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1364
										}
										position++
										goto l1363
									l1364:
										position, tokenIndex, depth = position1363, tokenIndex1363, depth1363
										if buffer[position] != rune('E') {
											goto l1205
										}
										position++
									}
								l1363:
									{
										position1365, tokenIndex1365, depth1365 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1366
										}
										position++
										goto l1365
									l1366:
										position, tokenIndex, depth = position1365, tokenIndex1365, depth1365
										if buffer[position] != rune('S') {
											goto l1205
										}
										position++
									}
								l1365:
									if !rules[ruleskip]() {
										goto l1205
									}
									depth--
									add(ruleLANGMATCHES, position1344)
								}
								break
							}
						}

					}
				l1206:
					if !rules[ruleLPAREN]() {
						goto l1205
					}
					if !rules[ruleexpression]() {
						goto l1205
					}
					if !rules[ruleCOMMA]() {
						goto l1205
					}
					if !rules[ruleexpression]() {
						goto l1205
					}
					if !rules[ruleRPAREN]() {
						goto l1205
					}
					goto l841
				l1205:
					position, tokenIndex, depth = position841, tokenIndex841, depth841
					{
						position1368 := position
						depth++
						{
							position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1370
							}
							position++
							goto l1369
						l1370:
							position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
							if buffer[position] != rune('B') {
								goto l1367
							}
							position++
						}
					l1369:
						{
							position1371, tokenIndex1371, depth1371 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1372
							}
							position++
							goto l1371
						l1372:
							position, tokenIndex, depth = position1371, tokenIndex1371, depth1371
							if buffer[position] != rune('O') {
								goto l1367
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
								goto l1367
							}
							position++
						}
					l1373:
						{
							position1375, tokenIndex1375, depth1375 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1376
							}
							position++
							goto l1375
						l1376:
							position, tokenIndex, depth = position1375, tokenIndex1375, depth1375
							if buffer[position] != rune('N') {
								goto l1367
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
								goto l1367
							}
							position++
						}
					l1377:
						if !rules[ruleskip]() {
							goto l1367
						}
						depth--
						add(ruleBOUND, position1368)
					}
					if !rules[ruleLPAREN]() {
						goto l1367
					}
					if !rules[rulevar]() {
						goto l1367
					}
					if !rules[ruleRPAREN]() {
						goto l1367
					}
					goto l841
				l1367:
					position, tokenIndex, depth = position841, tokenIndex841, depth841
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1381 := position
								depth++
								{
									position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1383
									}
									position++
									goto l1382
								l1383:
									position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
									if buffer[position] != rune('S') {
										goto l1379
									}
									position++
								}
							l1382:
								{
									position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1385
									}
									position++
									goto l1384
								l1385:
									position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
									if buffer[position] != rune('T') {
										goto l1379
									}
									position++
								}
							l1384:
								{
									position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1387
									}
									position++
									goto l1386
								l1387:
									position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
									if buffer[position] != rune('R') {
										goto l1379
									}
									position++
								}
							l1386:
								{
									position1388, tokenIndex1388, depth1388 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1389
									}
									position++
									goto l1388
								l1389:
									position, tokenIndex, depth = position1388, tokenIndex1388, depth1388
									if buffer[position] != rune('U') {
										goto l1379
									}
									position++
								}
							l1388:
								{
									position1390, tokenIndex1390, depth1390 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1391
									}
									position++
									goto l1390
								l1391:
									position, tokenIndex, depth = position1390, tokenIndex1390, depth1390
									if buffer[position] != rune('U') {
										goto l1379
									}
									position++
								}
							l1390:
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
										goto l1379
									}
									position++
								}
							l1392:
								{
									position1394, tokenIndex1394, depth1394 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1395
									}
									position++
									goto l1394
								l1395:
									position, tokenIndex, depth = position1394, tokenIndex1394, depth1394
									if buffer[position] != rune('D') {
										goto l1379
									}
									position++
								}
							l1394:
								if !rules[ruleskip]() {
									goto l1379
								}
								depth--
								add(ruleSTRUUID, position1381)
							}
							break
						case 'U', 'u':
							{
								position1396 := position
								depth++
								{
									position1397, tokenIndex1397, depth1397 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1398
									}
									position++
									goto l1397
								l1398:
									position, tokenIndex, depth = position1397, tokenIndex1397, depth1397
									if buffer[position] != rune('U') {
										goto l1379
									}
									position++
								}
							l1397:
								{
									position1399, tokenIndex1399, depth1399 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1400
									}
									position++
									goto l1399
								l1400:
									position, tokenIndex, depth = position1399, tokenIndex1399, depth1399
									if buffer[position] != rune('U') {
										goto l1379
									}
									position++
								}
							l1399:
								{
									position1401, tokenIndex1401, depth1401 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1402
									}
									position++
									goto l1401
								l1402:
									position, tokenIndex, depth = position1401, tokenIndex1401, depth1401
									if buffer[position] != rune('I') {
										goto l1379
									}
									position++
								}
							l1401:
								{
									position1403, tokenIndex1403, depth1403 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1404
									}
									position++
									goto l1403
								l1404:
									position, tokenIndex, depth = position1403, tokenIndex1403, depth1403
									if buffer[position] != rune('D') {
										goto l1379
									}
									position++
								}
							l1403:
								if !rules[ruleskip]() {
									goto l1379
								}
								depth--
								add(ruleUUID, position1396)
							}
							break
						case 'N', 'n':
							{
								position1405 := position
								depth++
								{
									position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1407
									}
									position++
									goto l1406
								l1407:
									position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
									if buffer[position] != rune('N') {
										goto l1379
									}
									position++
								}
							l1406:
								{
									position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1409
									}
									position++
									goto l1408
								l1409:
									position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
									if buffer[position] != rune('O') {
										goto l1379
									}
									position++
								}
							l1408:
								{
									position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1411
									}
									position++
									goto l1410
								l1411:
									position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
									if buffer[position] != rune('W') {
										goto l1379
									}
									position++
								}
							l1410:
								if !rules[ruleskip]() {
									goto l1379
								}
								depth--
								add(ruleNOW, position1405)
							}
							break
						default:
							{
								position1412 := position
								depth++
								{
									position1413, tokenIndex1413, depth1413 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1414
									}
									position++
									goto l1413
								l1414:
									position, tokenIndex, depth = position1413, tokenIndex1413, depth1413
									if buffer[position] != rune('R') {
										goto l1379
									}
									position++
								}
							l1413:
								{
									position1415, tokenIndex1415, depth1415 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1416
									}
									position++
									goto l1415
								l1416:
									position, tokenIndex, depth = position1415, tokenIndex1415, depth1415
									if buffer[position] != rune('A') {
										goto l1379
									}
									position++
								}
							l1415:
								{
									position1417, tokenIndex1417, depth1417 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1418
									}
									position++
									goto l1417
								l1418:
									position, tokenIndex, depth = position1417, tokenIndex1417, depth1417
									if buffer[position] != rune('N') {
										goto l1379
									}
									position++
								}
							l1417:
								{
									position1419, tokenIndex1419, depth1419 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1420
									}
									position++
									goto l1419
								l1420:
									position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
									if buffer[position] != rune('D') {
										goto l1379
									}
									position++
								}
							l1419:
								if !rules[ruleskip]() {
									goto l1379
								}
								depth--
								add(ruleRAND, position1412)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1379
					}
					goto l841
				l1379:
					position, tokenIndex, depth = position841, tokenIndex841, depth841
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1422, tokenIndex1422, depth1422 := position, tokenIndex, depth
								{
									position1424 := position
									depth++
									{
										position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1426
										}
										position++
										goto l1425
									l1426:
										position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
										if buffer[position] != rune('E') {
											goto l1423
										}
										position++
									}
								l1425:
									{
										position1427, tokenIndex1427, depth1427 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1428
										}
										position++
										goto l1427
									l1428:
										position, tokenIndex, depth = position1427, tokenIndex1427, depth1427
										if buffer[position] != rune('X') {
											goto l1423
										}
										position++
									}
								l1427:
									{
										position1429, tokenIndex1429, depth1429 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1430
										}
										position++
										goto l1429
									l1430:
										position, tokenIndex, depth = position1429, tokenIndex1429, depth1429
										if buffer[position] != rune('I') {
											goto l1423
										}
										position++
									}
								l1429:
									{
										position1431, tokenIndex1431, depth1431 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1432
										}
										position++
										goto l1431
									l1432:
										position, tokenIndex, depth = position1431, tokenIndex1431, depth1431
										if buffer[position] != rune('S') {
											goto l1423
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
											goto l1423
										}
										position++
									}
								l1433:
									{
										position1435, tokenIndex1435, depth1435 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1436
										}
										position++
										goto l1435
									l1436:
										position, tokenIndex, depth = position1435, tokenIndex1435, depth1435
										if buffer[position] != rune('S') {
											goto l1423
										}
										position++
									}
								l1435:
									if !rules[ruleskip]() {
										goto l1423
									}
									depth--
									add(ruleEXISTS, position1424)
								}
								goto l1422
							l1423:
								position, tokenIndex, depth = position1422, tokenIndex1422, depth1422
								{
									position1437 := position
									depth++
									{
										position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1439
										}
										position++
										goto l1438
									l1439:
										position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
										if buffer[position] != rune('N') {
											goto l839
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
											goto l839
										}
										position++
									}
								l1440:
									{
										position1442, tokenIndex1442, depth1442 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1443
										}
										position++
										goto l1442
									l1443:
										position, tokenIndex, depth = position1442, tokenIndex1442, depth1442
										if buffer[position] != rune('T') {
											goto l839
										}
										position++
									}
								l1442:
									if buffer[position] != rune(' ') {
										goto l839
									}
									position++
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
											goto l839
										}
										position++
									}
								l1444:
									{
										position1446, tokenIndex1446, depth1446 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1447
										}
										position++
										goto l1446
									l1447:
										position, tokenIndex, depth = position1446, tokenIndex1446, depth1446
										if buffer[position] != rune('X') {
											goto l839
										}
										position++
									}
								l1446:
									{
										position1448, tokenIndex1448, depth1448 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1449
										}
										position++
										goto l1448
									l1449:
										position, tokenIndex, depth = position1448, tokenIndex1448, depth1448
										if buffer[position] != rune('I') {
											goto l839
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
											goto l839
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
											goto l839
										}
										position++
									}
								l1452:
									{
										position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1455
										}
										position++
										goto l1454
									l1455:
										position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
										if buffer[position] != rune('S') {
											goto l839
										}
										position++
									}
								l1454:
									if !rules[ruleskip]() {
										goto l839
									}
									depth--
									add(ruleNOTEXIST, position1437)
								}
							}
						l1422:
							if !rules[rulegroupGraphPattern]() {
								goto l839
							}
							break
						case 'I', 'i':
							{
								position1456 := position
								depth++
								{
									position1457, tokenIndex1457, depth1457 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1458
									}
									position++
									goto l1457
								l1458:
									position, tokenIndex, depth = position1457, tokenIndex1457, depth1457
									if buffer[position] != rune('I') {
										goto l839
									}
									position++
								}
							l1457:
								{
									position1459, tokenIndex1459, depth1459 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1460
									}
									position++
									goto l1459
								l1460:
									position, tokenIndex, depth = position1459, tokenIndex1459, depth1459
									if buffer[position] != rune('F') {
										goto l839
									}
									position++
								}
							l1459:
								if !rules[ruleskip]() {
									goto l839
								}
								depth--
								add(ruleIF, position1456)
							}
							if !rules[ruleLPAREN]() {
								goto l839
							}
							if !rules[ruleexpression]() {
								goto l839
							}
							if !rules[ruleCOMMA]() {
								goto l839
							}
							if !rules[ruleexpression]() {
								goto l839
							}
							if !rules[ruleCOMMA]() {
								goto l839
							}
							if !rules[ruleexpression]() {
								goto l839
							}
							if !rules[ruleRPAREN]() {
								goto l839
							}
							break
						case 'C', 'c':
							{
								position1461, tokenIndex1461, depth1461 := position, tokenIndex, depth
								{
									position1463 := position
									depth++
									{
										position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1465
										}
										position++
										goto l1464
									l1465:
										position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
										if buffer[position] != rune('C') {
											goto l1462
										}
										position++
									}
								l1464:
									{
										position1466, tokenIndex1466, depth1466 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1467
										}
										position++
										goto l1466
									l1467:
										position, tokenIndex, depth = position1466, tokenIndex1466, depth1466
										if buffer[position] != rune('O') {
											goto l1462
										}
										position++
									}
								l1466:
									{
										position1468, tokenIndex1468, depth1468 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1469
										}
										position++
										goto l1468
									l1469:
										position, tokenIndex, depth = position1468, tokenIndex1468, depth1468
										if buffer[position] != rune('N') {
											goto l1462
										}
										position++
									}
								l1468:
									{
										position1470, tokenIndex1470, depth1470 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1471
										}
										position++
										goto l1470
									l1471:
										position, tokenIndex, depth = position1470, tokenIndex1470, depth1470
										if buffer[position] != rune('C') {
											goto l1462
										}
										position++
									}
								l1470:
									{
										position1472, tokenIndex1472, depth1472 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1473
										}
										position++
										goto l1472
									l1473:
										position, tokenIndex, depth = position1472, tokenIndex1472, depth1472
										if buffer[position] != rune('A') {
											goto l1462
										}
										position++
									}
								l1472:
									{
										position1474, tokenIndex1474, depth1474 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1475
										}
										position++
										goto l1474
									l1475:
										position, tokenIndex, depth = position1474, tokenIndex1474, depth1474
										if buffer[position] != rune('T') {
											goto l1462
										}
										position++
									}
								l1474:
									if !rules[ruleskip]() {
										goto l1462
									}
									depth--
									add(ruleCONCAT, position1463)
								}
								goto l1461
							l1462:
								position, tokenIndex, depth = position1461, tokenIndex1461, depth1461
								{
									position1476 := position
									depth++
									{
										position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1478
										}
										position++
										goto l1477
									l1478:
										position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
										if buffer[position] != rune('C') {
											goto l839
										}
										position++
									}
								l1477:
									{
										position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1480
										}
										position++
										goto l1479
									l1480:
										position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
										if buffer[position] != rune('O') {
											goto l839
										}
										position++
									}
								l1479:
									{
										position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1482
										}
										position++
										goto l1481
									l1482:
										position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
										if buffer[position] != rune('A') {
											goto l839
										}
										position++
									}
								l1481:
									{
										position1483, tokenIndex1483, depth1483 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1484
										}
										position++
										goto l1483
									l1484:
										position, tokenIndex, depth = position1483, tokenIndex1483, depth1483
										if buffer[position] != rune('L') {
											goto l839
										}
										position++
									}
								l1483:
									{
										position1485, tokenIndex1485, depth1485 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1486
										}
										position++
										goto l1485
									l1486:
										position, tokenIndex, depth = position1485, tokenIndex1485, depth1485
										if buffer[position] != rune('E') {
											goto l839
										}
										position++
									}
								l1485:
									{
										position1487, tokenIndex1487, depth1487 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1488
										}
										position++
										goto l1487
									l1488:
										position, tokenIndex, depth = position1487, tokenIndex1487, depth1487
										if buffer[position] != rune('S') {
											goto l839
										}
										position++
									}
								l1487:
									{
										position1489, tokenIndex1489, depth1489 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1490
										}
										position++
										goto l1489
									l1490:
										position, tokenIndex, depth = position1489, tokenIndex1489, depth1489
										if buffer[position] != rune('C') {
											goto l839
										}
										position++
									}
								l1489:
									{
										position1491, tokenIndex1491, depth1491 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1492
										}
										position++
										goto l1491
									l1492:
										position, tokenIndex, depth = position1491, tokenIndex1491, depth1491
										if buffer[position] != rune('E') {
											goto l839
										}
										position++
									}
								l1491:
									if !rules[ruleskip]() {
										goto l839
									}
									depth--
									add(ruleCOALESCE, position1476)
								}
							}
						l1461:
							if !rules[ruleargList]() {
								goto l839
							}
							break
						case 'B', 'b':
							{
								position1493 := position
								depth++
								{
									position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1495
									}
									position++
									goto l1494
								l1495:
									position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
									if buffer[position] != rune('B') {
										goto l839
									}
									position++
								}
							l1494:
								{
									position1496, tokenIndex1496, depth1496 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1497
									}
									position++
									goto l1496
								l1497:
									position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
									if buffer[position] != rune('N') {
										goto l839
									}
									position++
								}
							l1496:
								{
									position1498, tokenIndex1498, depth1498 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1499
									}
									position++
									goto l1498
								l1499:
									position, tokenIndex, depth = position1498, tokenIndex1498, depth1498
									if buffer[position] != rune('O') {
										goto l839
									}
									position++
								}
							l1498:
								{
									position1500, tokenIndex1500, depth1500 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1501
									}
									position++
									goto l1500
								l1501:
									position, tokenIndex, depth = position1500, tokenIndex1500, depth1500
									if buffer[position] != rune('D') {
										goto l839
									}
									position++
								}
							l1500:
								{
									position1502, tokenIndex1502, depth1502 := position, tokenIndex, depth
									if buffer[position] != rune('e') {
										goto l1503
									}
									position++
									goto l1502
								l1503:
									position, tokenIndex, depth = position1502, tokenIndex1502, depth1502
									if buffer[position] != rune('E') {
										goto l839
									}
									position++
								}
							l1502:
								if !rules[ruleskip]() {
									goto l839
								}
								depth--
								add(ruleBNODE, position1493)
							}
							{
								position1504, tokenIndex1504, depth1504 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1505
								}
								if !rules[ruleexpression]() {
									goto l1505
								}
								if !rules[ruleRPAREN]() {
									goto l1505
								}
								goto l1504
							l1505:
								position, tokenIndex, depth = position1504, tokenIndex1504, depth1504
								if !rules[rulenil]() {
									goto l839
								}
							}
						l1504:
							break
						default:
							{
								position1506, tokenIndex1506, depth1506 := position, tokenIndex, depth
								{
									position1508 := position
									depth++
									{
										position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1510
										}
										position++
										goto l1509
									l1510:
										position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
										if buffer[position] != rune('S') {
											goto l1507
										}
										position++
									}
								l1509:
									{
										position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1512
										}
										position++
										goto l1511
									l1512:
										position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
										if buffer[position] != rune('U') {
											goto l1507
										}
										position++
									}
								l1511:
									{
										position1513, tokenIndex1513, depth1513 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1514
										}
										position++
										goto l1513
									l1514:
										position, tokenIndex, depth = position1513, tokenIndex1513, depth1513
										if buffer[position] != rune('B') {
											goto l1507
										}
										position++
									}
								l1513:
									{
										position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1516
										}
										position++
										goto l1515
									l1516:
										position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
										if buffer[position] != rune('S') {
											goto l1507
										}
										position++
									}
								l1515:
									{
										position1517, tokenIndex1517, depth1517 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1518
										}
										position++
										goto l1517
									l1518:
										position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
										if buffer[position] != rune('T') {
											goto l1507
										}
										position++
									}
								l1517:
									{
										position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1520
										}
										position++
										goto l1519
									l1520:
										position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
										if buffer[position] != rune('R') {
											goto l1507
										}
										position++
									}
								l1519:
									if !rules[ruleskip]() {
										goto l1507
									}
									depth--
									add(ruleSUBSTR, position1508)
								}
								goto l1506
							l1507:
								position, tokenIndex, depth = position1506, tokenIndex1506, depth1506
								{
									position1522 := position
									depth++
									{
										position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1524
										}
										position++
										goto l1523
									l1524:
										position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
										if buffer[position] != rune('R') {
											goto l1521
										}
										position++
									}
								l1523:
									{
										position1525, tokenIndex1525, depth1525 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1526
										}
										position++
										goto l1525
									l1526:
										position, tokenIndex, depth = position1525, tokenIndex1525, depth1525
										if buffer[position] != rune('E') {
											goto l1521
										}
										position++
									}
								l1525:
									{
										position1527, tokenIndex1527, depth1527 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1528
										}
										position++
										goto l1527
									l1528:
										position, tokenIndex, depth = position1527, tokenIndex1527, depth1527
										if buffer[position] != rune('P') {
											goto l1521
										}
										position++
									}
								l1527:
									{
										position1529, tokenIndex1529, depth1529 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1530
										}
										position++
										goto l1529
									l1530:
										position, tokenIndex, depth = position1529, tokenIndex1529, depth1529
										if buffer[position] != rune('L') {
											goto l1521
										}
										position++
									}
								l1529:
									{
										position1531, tokenIndex1531, depth1531 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1532
										}
										position++
										goto l1531
									l1532:
										position, tokenIndex, depth = position1531, tokenIndex1531, depth1531
										if buffer[position] != rune('A') {
											goto l1521
										}
										position++
									}
								l1531:
									{
										position1533, tokenIndex1533, depth1533 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1534
										}
										position++
										goto l1533
									l1534:
										position, tokenIndex, depth = position1533, tokenIndex1533, depth1533
										if buffer[position] != rune('C') {
											goto l1521
										}
										position++
									}
								l1533:
									{
										position1535, tokenIndex1535, depth1535 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1536
										}
										position++
										goto l1535
									l1536:
										position, tokenIndex, depth = position1535, tokenIndex1535, depth1535
										if buffer[position] != rune('E') {
											goto l1521
										}
										position++
									}
								l1535:
									if !rules[ruleskip]() {
										goto l1521
									}
									depth--
									add(ruleREPLACE, position1522)
								}
								goto l1506
							l1521:
								position, tokenIndex, depth = position1506, tokenIndex1506, depth1506
								{
									position1537 := position
									depth++
									{
										position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1539
										}
										position++
										goto l1538
									l1539:
										position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
										if buffer[position] != rune('R') {
											goto l839
										}
										position++
									}
								l1538:
									{
										position1540, tokenIndex1540, depth1540 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1541
										}
										position++
										goto l1540
									l1541:
										position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
										if buffer[position] != rune('E') {
											goto l839
										}
										position++
									}
								l1540:
									{
										position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1543
										}
										position++
										goto l1542
									l1543:
										position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
										if buffer[position] != rune('G') {
											goto l839
										}
										position++
									}
								l1542:
									{
										position1544, tokenIndex1544, depth1544 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1545
										}
										position++
										goto l1544
									l1545:
										position, tokenIndex, depth = position1544, tokenIndex1544, depth1544
										if buffer[position] != rune('E') {
											goto l839
										}
										position++
									}
								l1544:
									{
										position1546, tokenIndex1546, depth1546 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1547
										}
										position++
										goto l1546
									l1547:
										position, tokenIndex, depth = position1546, tokenIndex1546, depth1546
										if buffer[position] != rune('X') {
											goto l839
										}
										position++
									}
								l1546:
									if !rules[ruleskip]() {
										goto l839
									}
									depth--
									add(ruleREGEX, position1537)
								}
							}
						l1506:
							if !rules[ruleLPAREN]() {
								goto l839
							}
							if !rules[ruleexpression]() {
								goto l839
							}
							if !rules[ruleCOMMA]() {
								goto l839
							}
							if !rules[ruleexpression]() {
								goto l839
							}
							{
								position1548, tokenIndex1548, depth1548 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1548
								}
								if !rules[ruleexpression]() {
									goto l1548
								}
								goto l1549
							l1548:
								position, tokenIndex, depth = position1548, tokenIndex1548, depth1548
							}
						l1549:
							if !rules[ruleRPAREN]() {
								goto l839
							}
							break
						}
					}

				}
			l841:
				depth--
				add(rulebuiltinCall, position840)
			}
			return true
		l839:
			position, tokenIndex, depth = position839, tokenIndex839, depth839
			return false
		},
		/* 75 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1550, tokenIndex1550, depth1550 := position, tokenIndex, depth
			{
				position1551 := position
				depth++
				{
					position1552, tokenIndex1552, depth1552 := position, tokenIndex, depth
					{
						position1554 := position
						depth++
					l1555:
						{
							position1556, tokenIndex1556, depth1556 := position, tokenIndex, depth
							{
								position1557, tokenIndex1557, depth1557 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1558
								}
								position++
								goto l1557
							l1558:
								position, tokenIndex, depth = position1557, tokenIndex1557, depth1557
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1556
								}
								position++
							}
						l1557:
							goto l1555
						l1556:
							position, tokenIndex, depth = position1556, tokenIndex1556, depth1556
						}
						depth--
						add(rulePegText, position1554)
					}
					if buffer[position] != rune(':') {
						goto l1553
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1552
				l1553:
					position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
					{
						position1561 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1560
						}
						position++
					l1562:
						{
							position1563, tokenIndex1563, depth1563 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1563
							}
							position++
							goto l1562
						l1563:
							position, tokenIndex, depth = position1563, tokenIndex1563, depth1563
						}
						depth--
						add(rulePegText, position1561)
					}
					if buffer[position] != rune('/') {
						goto l1560
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1552
				l1560:
					position, tokenIndex, depth = position1552, tokenIndex1552, depth1552
					{
						position1565 := position
						depth++
					l1566:
						{
							position1567, tokenIndex1567, depth1567 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1567
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1567
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1567
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1567
									}
									position++
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1567
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1567
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1567
									}
									position++
									break
								}
							}

							goto l1566
						l1567:
							position, tokenIndex, depth = position1567, tokenIndex1567, depth1567
						}
						depth--
						add(rulePegText, position1565)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1552:
				if buffer[position] != rune('<') {
					goto l1550
				}
				position++
				if !rules[rulews]() {
					goto l1550
				}
				if !rules[ruleskip]() {
					goto l1550
				}
				depth--
				add(rulepof, position1551)
			}
			return true
		l1550:
			position, tokenIndex, depth = position1550, tokenIndex1550, depth1550
			return false
		},
		/* 76 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1570, tokenIndex1570, depth1570 := position, tokenIndex, depth
			{
				position1571 := position
				depth++
				{
					position1572, tokenIndex1572, depth1572 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1573
					}
					position++
					goto l1572
				l1573:
					position, tokenIndex, depth = position1572, tokenIndex1572, depth1572
					if buffer[position] != rune('$') {
						goto l1570
					}
					position++
				}
			l1572:
				{
					position1574 := position
					depth++
					{
						position1575, tokenIndex1575, depth1575 := position, tokenIndex, depth
						if !rules[rulepnCharsU]() {
							goto l1576
						}
						goto l1575
					l1576:
						position, tokenIndex, depth = position1575, tokenIndex1575, depth1575
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1570
						}
						position++
					}
				l1575:
				l1577:
					{
						position1578, tokenIndex1578, depth1578 := position, tokenIndex, depth
						{
							position1579, tokenIndex1579, depth1579 := position, tokenIndex, depth
							if !rules[rulepnCharsU]() {
								goto l1580
							}
							goto l1579
						l1580:
							position, tokenIndex, depth = position1579, tokenIndex1579, depth1579
							{
								switch buffer[position] {
								case 'â':
									if c := buffer[position]; c < rune('‿') || c > rune('⁀') {
										goto l1578
									}
									position++
									break
								case 'Ì', 'Í':
									if c := buffer[position]; c < rune('̀') || c > rune('ͯ') {
										goto l1578
									}
									position++
									break
								case 'Â':
									if buffer[position] != rune('·') {
										goto l1578
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1578
									}
									position++
									break
								}
							}

						}
					l1579:
						goto l1577
					l1578:
						position, tokenIndex, depth = position1578, tokenIndex1578, depth1578
					}
					depth--
					add(ruleVARNAME, position1574)
				}
				if !rules[ruleskip]() {
					goto l1570
				}
				depth--
				add(rulevar, position1571)
			}
			return true
		l1570:
			position, tokenIndex, depth = position1570, tokenIndex1570, depth1570
			return false
		},
		/* 77 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1582, tokenIndex1582, depth1582 := position, tokenIndex, depth
			{
				position1583 := position
				depth++
				{
					position1584, tokenIndex1584, depth1584 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1585
					}
					goto l1584
				l1585:
					position, tokenIndex, depth = position1584, tokenIndex1584, depth1584
					{
						position1586 := position
						depth++
						{
							position1587, tokenIndex1587, depth1587 := position, tokenIndex, depth
							if !rules[rulepnPrefix]() {
								goto l1587
							}
							goto l1588
						l1587:
							position, tokenIndex, depth = position1587, tokenIndex1587, depth1587
						}
					l1588:
						if buffer[position] != rune(':') {
							goto l1582
						}
						position++
						{
							position1589 := position
							depth++
							{
								switch buffer[position] {
								case '%', '\\':
									{
										position1593 := position
										depth++
										{
											position1594, tokenIndex1594, depth1594 := position, tokenIndex, depth
											{
												position1596 := position
												depth++
												if buffer[position] != rune('%') {
													goto l1595
												}
												position++
												if !rules[rulehex]() {
													goto l1595
												}
												if !rules[rulehex]() {
													goto l1595
												}
												depth--
												add(rulepercent, position1596)
											}
											goto l1594
										l1595:
											position, tokenIndex, depth = position1594, tokenIndex1594, depth1594
											{
												position1597 := position
												depth++
												if buffer[position] != rune('\\') {
													goto l1582
												}
												position++
												{
													switch buffer[position] {
													case '%':
														if buffer[position] != rune('%') {
															goto l1582
														}
														position++
														break
													case '@':
														if buffer[position] != rune('@') {
															goto l1582
														}
														position++
														break
													case '#':
														if buffer[position] != rune('#') {
															goto l1582
														}
														position++
														break
													case '?':
														if buffer[position] != rune('?') {
															goto l1582
														}
														position++
														break
													case '/':
														if buffer[position] != rune('/') {
															goto l1582
														}
														position++
														break
													case '=':
														if buffer[position] != rune('=') {
															goto l1582
														}
														position++
														break
													case ';':
														if buffer[position] != rune(';') {
															goto l1582
														}
														position++
														break
													case ',':
														if buffer[position] != rune(',') {
															goto l1582
														}
														position++
														break
													case '+':
														if buffer[position] != rune('+') {
															goto l1582
														}
														position++
														break
													case '*':
														if buffer[position] != rune('*') {
															goto l1582
														}
														position++
														break
													case ')':
														if buffer[position] != rune(')') {
															goto l1582
														}
														position++
														break
													case '(':
														if buffer[position] != rune('(') {
															goto l1582
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l1582
														}
														position++
														break
													case '&':
														if buffer[position] != rune('&') {
															goto l1582
														}
														position++
														break
													case '$':
														if buffer[position] != rune('$') {
															goto l1582
														}
														position++
														break
													case '!':
														if buffer[position] != rune('!') {
															goto l1582
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l1582
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l1582
														}
														position++
														break
													case '~':
														if buffer[position] != rune('~') {
															goto l1582
														}
														position++
														break
													default:
														if buffer[position] != rune('_') {
															goto l1582
														}
														position++
														break
													}
												}

												depth--
												add(rulepnLocalEsc, position1597)
											}
										}
									l1594:
										depth--
										add(ruleplx, position1593)
									}
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1582
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l1582
									}
									position++
									break
								default:
									if !rules[rulepnCharsU]() {
										goto l1582
									}
									break
								}
							}

						l1590:
							{
								position1591, tokenIndex1591, depth1591 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '%', '\\':
										{
											position1600 := position
											depth++
											{
												position1601, tokenIndex1601, depth1601 := position, tokenIndex, depth
												{
													position1603 := position
													depth++
													if buffer[position] != rune('%') {
														goto l1602
													}
													position++
													if !rules[rulehex]() {
														goto l1602
													}
													if !rules[rulehex]() {
														goto l1602
													}
													depth--
													add(rulepercent, position1603)
												}
												goto l1601
											l1602:
												position, tokenIndex, depth = position1601, tokenIndex1601, depth1601
												{
													position1604 := position
													depth++
													if buffer[position] != rune('\\') {
														goto l1591
													}
													position++
													{
														switch buffer[position] {
														case '%':
															if buffer[position] != rune('%') {
																goto l1591
															}
															position++
															break
														case '@':
															if buffer[position] != rune('@') {
																goto l1591
															}
															position++
															break
														case '#':
															if buffer[position] != rune('#') {
																goto l1591
															}
															position++
															break
														case '?':
															if buffer[position] != rune('?') {
																goto l1591
															}
															position++
															break
														case '/':
															if buffer[position] != rune('/') {
																goto l1591
															}
															position++
															break
														case '=':
															if buffer[position] != rune('=') {
																goto l1591
															}
															position++
															break
														case ';':
															if buffer[position] != rune(';') {
																goto l1591
															}
															position++
															break
														case ',':
															if buffer[position] != rune(',') {
																goto l1591
															}
															position++
															break
														case '+':
															if buffer[position] != rune('+') {
																goto l1591
															}
															position++
															break
														case '*':
															if buffer[position] != rune('*') {
																goto l1591
															}
															position++
															break
														case ')':
															if buffer[position] != rune(')') {
																goto l1591
															}
															position++
															break
														case '(':
															if buffer[position] != rune('(') {
																goto l1591
															}
															position++
															break
														case '\'':
															if buffer[position] != rune('\'') {
																goto l1591
															}
															position++
															break
														case '&':
															if buffer[position] != rune('&') {
																goto l1591
															}
															position++
															break
														case '$':
															if buffer[position] != rune('$') {
																goto l1591
															}
															position++
															break
														case '!':
															if buffer[position] != rune('!') {
																goto l1591
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l1591
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l1591
															}
															position++
															break
														case '~':
															if buffer[position] != rune('~') {
																goto l1591
															}
															position++
															break
														default:
															if buffer[position] != rune('_') {
																goto l1591
															}
															position++
															break
														}
													}

													depth--
													add(rulepnLocalEsc, position1604)
												}
											}
										l1601:
											depth--
											add(ruleplx, position1600)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1591
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1591
										}
										position++
										break
									default:
										if !rules[rulepnCharsU]() {
											goto l1591
										}
										break
									}
								}

								goto l1590
							l1591:
								position, tokenIndex, depth = position1591, tokenIndex1591, depth1591
							}
							depth--
							add(rulepnLocal, position1589)
						}
						if !rules[ruleskip]() {
							goto l1582
						}
						depth--
						add(ruleprefixedName, position1586)
					}
				}
			l1584:
				depth--
				add(ruleiriref, position1583)
			}
			return true
		l1582:
			position, tokenIndex, depth = position1582, tokenIndex1582, depth1582
			return false
		},
		/* 78 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1606, tokenIndex1606, depth1606 := position, tokenIndex, depth
			{
				position1607 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1606
				}
				position++
			l1608:
				{
					position1609, tokenIndex1609, depth1609 := position, tokenIndex, depth
					{
						position1610, tokenIndex1610, depth1610 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1610
						}
						position++
						goto l1609
					l1610:
						position, tokenIndex, depth = position1610, tokenIndex1610, depth1610
					}
					if !matchDot() {
						goto l1609
					}
					goto l1608
				l1609:
					position, tokenIndex, depth = position1609, tokenIndex1609, depth1609
				}
				if buffer[position] != rune('>') {
					goto l1606
				}
				position++
				if !rules[ruleskip]() {
					goto l1606
				}
				depth--
				add(ruleiri, position1607)
			}
			return true
		l1606:
			position, tokenIndex, depth = position1606, tokenIndex1606, depth1606
			return false
		},
		/* 79 prefixedName <- <(pnPrefix? ':' pnLocal skip)> */
		nil,
		/* 80 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1612, tokenIndex1612, depth1612 := position, tokenIndex, depth
			{
				position1613 := position
				depth++
				if !rules[rulestring]() {
					goto l1612
				}
				{
					position1614, tokenIndex1614, depth1614 := position, tokenIndex, depth
					{
						position1616, tokenIndex1616, depth1616 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1617
						}
						position++
						{
							position1620, tokenIndex1620, depth1620 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1621
							}
							position++
							goto l1620
						l1621:
							position, tokenIndex, depth = position1620, tokenIndex1620, depth1620
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1617
							}
							position++
						}
					l1620:
					l1618:
						{
							position1619, tokenIndex1619, depth1619 := position, tokenIndex, depth
							{
								position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1623
								}
								position++
								goto l1622
							l1623:
								position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1619
								}
								position++
							}
						l1622:
							goto l1618
						l1619:
							position, tokenIndex, depth = position1619, tokenIndex1619, depth1619
						}
					l1624:
						{
							position1625, tokenIndex1625, depth1625 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1625
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1625
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1625
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1625
									}
									position++
									break
								}
							}

						l1626:
							{
								position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1627
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1627
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1627
										}
										position++
										break
									}
								}

								goto l1626
							l1627:
								position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
							}
							goto l1624
						l1625:
							position, tokenIndex, depth = position1625, tokenIndex1625, depth1625
						}
						goto l1616
					l1617:
						position, tokenIndex, depth = position1616, tokenIndex1616, depth1616
						if buffer[position] != rune('^') {
							goto l1614
						}
						position++
						if buffer[position] != rune('^') {
							goto l1614
						}
						position++
						if !rules[ruleiriref]() {
							goto l1614
						}
					}
				l1616:
					goto l1615
				l1614:
					position, tokenIndex, depth = position1614, tokenIndex1614, depth1614
				}
			l1615:
				if !rules[ruleskip]() {
					goto l1612
				}
				depth--
				add(ruleliteral, position1613)
			}
			return true
		l1612:
			position, tokenIndex, depth = position1612, tokenIndex1612, depth1612
			return false
		},
		/* 81 string <- <(stringLiteralA / stringLiteralB / stringLiteralLongA / stringLiteralLongB)> */
		func() bool {
			position1630, tokenIndex1630, depth1630 := position, tokenIndex, depth
			{
				position1631 := position
				depth++
				{
					position1632, tokenIndex1632, depth1632 := position, tokenIndex, depth
					{
						position1634 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1633
						}
						position++
					l1635:
						{
							position1636, tokenIndex1636, depth1636 := position, tokenIndex, depth
							{
								position1637, tokenIndex1637, depth1637 := position, tokenIndex, depth
								{
									position1639, tokenIndex1639, depth1639 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1639
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1639
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1639
											}
											position++
											break
										default:
											if buffer[position] != rune('\'') {
												goto l1639
											}
											position++
											break
										}
									}

									goto l1638
								l1639:
									position, tokenIndex, depth = position1639, tokenIndex1639, depth1639
								}
								if !matchDot() {
									goto l1638
								}
								goto l1637
							l1638:
								position, tokenIndex, depth = position1637, tokenIndex1637, depth1637
								if !rules[ruleechar]() {
									goto l1636
								}
							}
						l1637:
							goto l1635
						l1636:
							position, tokenIndex, depth = position1636, tokenIndex1636, depth1636
						}
						if buffer[position] != rune('\'') {
							goto l1633
						}
						position++
						depth--
						add(rulestringLiteralA, position1634)
					}
					goto l1632
				l1633:
					position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
					{
						position1642 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1641
						}
						position++
					l1643:
						{
							position1644, tokenIndex1644, depth1644 := position, tokenIndex, depth
							{
								position1645, tokenIndex1645, depth1645 := position, tokenIndex, depth
								{
									position1647, tokenIndex1647, depth1647 := position, tokenIndex, depth
									{
										switch buffer[position] {
										case '\r':
											if buffer[position] != rune('\r') {
												goto l1647
											}
											position++
											break
										case '\n':
											if buffer[position] != rune('\n') {
												goto l1647
											}
											position++
											break
										case '\\':
											if buffer[position] != rune('\\') {
												goto l1647
											}
											position++
											break
										default:
											if buffer[position] != rune('"') {
												goto l1647
											}
											position++
											break
										}
									}

									goto l1646
								l1647:
									position, tokenIndex, depth = position1647, tokenIndex1647, depth1647
								}
								if !matchDot() {
									goto l1646
								}
								goto l1645
							l1646:
								position, tokenIndex, depth = position1645, tokenIndex1645, depth1645
								if !rules[ruleechar]() {
									goto l1644
								}
							}
						l1645:
							goto l1643
						l1644:
							position, tokenIndex, depth = position1644, tokenIndex1644, depth1644
						}
						if buffer[position] != rune('"') {
							goto l1641
						}
						position++
						depth--
						add(rulestringLiteralB, position1642)
					}
					goto l1632
				l1641:
					position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
					{
						position1650 := position
						depth++
						if buffer[position] != rune('\'') {
							goto l1649
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1649
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1649
						}
						position++
					l1651:
						{
							position1652, tokenIndex1652, depth1652 := position, tokenIndex, depth
							{
								position1653, tokenIndex1653, depth1653 := position, tokenIndex, depth
								{
									position1655, tokenIndex1655, depth1655 := position, tokenIndex, depth
									if buffer[position] != rune('\'') {
										goto l1656
									}
									position++
									goto l1655
								l1656:
									position, tokenIndex, depth = position1655, tokenIndex1655, depth1655
									if buffer[position] != rune('\'') {
										goto l1653
									}
									position++
									if buffer[position] != rune('\'') {
										goto l1653
									}
									position++
								}
							l1655:
								goto l1654
							l1653:
								position, tokenIndex, depth = position1653, tokenIndex1653, depth1653
							}
						l1654:
							{
								position1657, tokenIndex1657, depth1657 := position, tokenIndex, depth
								{
									position1659, tokenIndex1659, depth1659 := position, tokenIndex, depth
									{
										position1660, tokenIndex1660, depth1660 := position, tokenIndex, depth
										if buffer[position] != rune('\'') {
											goto l1661
										}
										position++
										goto l1660
									l1661:
										position, tokenIndex, depth = position1660, tokenIndex1660, depth1660
										if buffer[position] != rune('\\') {
											goto l1659
										}
										position++
									}
								l1660:
									goto l1658
								l1659:
									position, tokenIndex, depth = position1659, tokenIndex1659, depth1659
								}
								if !matchDot() {
									goto l1658
								}
								goto l1657
							l1658:
								position, tokenIndex, depth = position1657, tokenIndex1657, depth1657
								if !rules[ruleechar]() {
									goto l1652
								}
							}
						l1657:
							goto l1651
						l1652:
							position, tokenIndex, depth = position1652, tokenIndex1652, depth1652
						}
						if buffer[position] != rune('\'') {
							goto l1649
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1649
						}
						position++
						if buffer[position] != rune('\'') {
							goto l1649
						}
						position++
						depth--
						add(rulestringLiteralLongA, position1650)
					}
					goto l1632
				l1649:
					position, tokenIndex, depth = position1632, tokenIndex1632, depth1632
					{
						position1662 := position
						depth++
						if buffer[position] != rune('"') {
							goto l1630
						}
						position++
						if buffer[position] != rune('"') {
							goto l1630
						}
						position++
						if buffer[position] != rune('"') {
							goto l1630
						}
						position++
					l1663:
						{
							position1664, tokenIndex1664, depth1664 := position, tokenIndex, depth
							{
								position1665, tokenIndex1665, depth1665 := position, tokenIndex, depth
								{
									position1667, tokenIndex1667, depth1667 := position, tokenIndex, depth
									if buffer[position] != rune('"') {
										goto l1668
									}
									position++
									goto l1667
								l1668:
									position, tokenIndex, depth = position1667, tokenIndex1667, depth1667
									if buffer[position] != rune('"') {
										goto l1665
									}
									position++
									if buffer[position] != rune('"') {
										goto l1665
									}
									position++
								}
							l1667:
								goto l1666
							l1665:
								position, tokenIndex, depth = position1665, tokenIndex1665, depth1665
							}
						l1666:
							{
								position1669, tokenIndex1669, depth1669 := position, tokenIndex, depth
								{
									position1671, tokenIndex1671, depth1671 := position, tokenIndex, depth
									{
										position1672, tokenIndex1672, depth1672 := position, tokenIndex, depth
										if buffer[position] != rune('"') {
											goto l1673
										}
										position++
										goto l1672
									l1673:
										position, tokenIndex, depth = position1672, tokenIndex1672, depth1672
										if buffer[position] != rune('\\') {
											goto l1671
										}
										position++
									}
								l1672:
									goto l1670
								l1671:
									position, tokenIndex, depth = position1671, tokenIndex1671, depth1671
								}
								if !matchDot() {
									goto l1670
								}
								goto l1669
							l1670:
								position, tokenIndex, depth = position1669, tokenIndex1669, depth1669
								if !rules[ruleechar]() {
									goto l1664
								}
							}
						l1669:
							goto l1663
						l1664:
							position, tokenIndex, depth = position1664, tokenIndex1664, depth1664
						}
						if buffer[position] != rune('"') {
							goto l1630
						}
						position++
						if buffer[position] != rune('"') {
							goto l1630
						}
						position++
						if buffer[position] != rune('"') {
							goto l1630
						}
						position++
						depth--
						add(rulestringLiteralLongB, position1662)
					}
				}
			l1632:
				depth--
				add(rulestring, position1631)
			}
			return true
		l1630:
			position, tokenIndex, depth = position1630, tokenIndex1630, depth1630
			return false
		},
		/* 82 stringLiteralA <- <('\'' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('\'') '\'')) .) / echar)* '\'')> */
		nil,
		/* 83 stringLiteralB <- <('"' ((!((&('\r') '\r') | (&('\n') '\n') | (&('\\') '\\') | (&('"') '"')) .) / echar)* '"')> */
		nil,
		/* 84 stringLiteralLongA <- <('\'' '\'' '\'' (('\'' / ('\'' '\''))? ((!('\'' / '\\') .) / echar))* ('\'' '\'' '\''))> */
		nil,
		/* 85 stringLiteralLongB <- <('"' '"' '"' (('"' / ('"' '"'))? ((!('"' / '\\') .) / echar))* ('"' '"' '"'))> */
		nil,
		/* 86 echar <- <('\\' ((&('\'') '\'') | (&('"') '"') | (&('\\') '\\') | (&('f') 'f') | (&('r') 'r') | (&('n') 'n') | (&('b') 'b') | (&('t') 't')))> */
		func() bool {
			position1678, tokenIndex1678, depth1678 := position, tokenIndex, depth
			{
				position1679 := position
				depth++
				if buffer[position] != rune('\\') {
					goto l1678
				}
				position++
				{
					switch buffer[position] {
					case '\'':
						if buffer[position] != rune('\'') {
							goto l1678
						}
						position++
						break
					case '"':
						if buffer[position] != rune('"') {
							goto l1678
						}
						position++
						break
					case '\\':
						if buffer[position] != rune('\\') {
							goto l1678
						}
						position++
						break
					case 'f':
						if buffer[position] != rune('f') {
							goto l1678
						}
						position++
						break
					case 'r':
						if buffer[position] != rune('r') {
							goto l1678
						}
						position++
						break
					case 'n':
						if buffer[position] != rune('n') {
							goto l1678
						}
						position++
						break
					case 'b':
						if buffer[position] != rune('b') {
							goto l1678
						}
						position++
						break
					default:
						if buffer[position] != rune('t') {
							goto l1678
						}
						position++
						break
					}
				}

				depth--
				add(ruleechar, position1679)
			}
			return true
		l1678:
			position, tokenIndex, depth = position1678, tokenIndex1678, depth1678
			return false
		},
		/* 87 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1681, tokenIndex1681, depth1681 := position, tokenIndex, depth
			{
				position1682 := position
				depth++
				{
					position1683, tokenIndex1683, depth1683 := position, tokenIndex, depth
					{
						position1685, tokenIndex1685, depth1685 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1686
						}
						position++
						goto l1685
					l1686:
						position, tokenIndex, depth = position1685, tokenIndex1685, depth1685
						if buffer[position] != rune('-') {
							goto l1683
						}
						position++
					}
				l1685:
					goto l1684
				l1683:
					position, tokenIndex, depth = position1683, tokenIndex1683, depth1683
				}
			l1684:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1681
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
				{
					position1689, tokenIndex1689, depth1689 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1689
					}
					position++
				l1691:
					{
						position1692, tokenIndex1692, depth1692 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1692
						}
						position++
						goto l1691
					l1692:
						position, tokenIndex, depth = position1692, tokenIndex1692, depth1692
					}
					goto l1690
				l1689:
					position, tokenIndex, depth = position1689, tokenIndex1689, depth1689
				}
			l1690:
				if !rules[ruleskip]() {
					goto l1681
				}
				depth--
				add(rulenumericLiteral, position1682)
			}
			return true
		l1681:
			position, tokenIndex, depth = position1681, tokenIndex1681, depth1681
			return false
		},
		/* 88 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 89 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1694, tokenIndex1694, depth1694 := position, tokenIndex, depth
			{
				position1695 := position
				depth++
				{
					position1696, tokenIndex1696, depth1696 := position, tokenIndex, depth
					{
						position1698 := position
						depth++
						{
							position1699, tokenIndex1699, depth1699 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l1700
							}
							position++
							goto l1699
						l1700:
							position, tokenIndex, depth = position1699, tokenIndex1699, depth1699
							if buffer[position] != rune('T') {
								goto l1697
							}
							position++
						}
					l1699:
						{
							position1701, tokenIndex1701, depth1701 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1702
							}
							position++
							goto l1701
						l1702:
							position, tokenIndex, depth = position1701, tokenIndex1701, depth1701
							if buffer[position] != rune('R') {
								goto l1697
							}
							position++
						}
					l1701:
						{
							position1703, tokenIndex1703, depth1703 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1704
							}
							position++
							goto l1703
						l1704:
							position, tokenIndex, depth = position1703, tokenIndex1703, depth1703
							if buffer[position] != rune('U') {
								goto l1697
							}
							position++
						}
					l1703:
						{
							position1705, tokenIndex1705, depth1705 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1706
							}
							position++
							goto l1705
						l1706:
							position, tokenIndex, depth = position1705, tokenIndex1705, depth1705
							if buffer[position] != rune('E') {
								goto l1697
							}
							position++
						}
					l1705:
						if !rules[ruleskip]() {
							goto l1697
						}
						depth--
						add(ruleTRUE, position1698)
					}
					goto l1696
				l1697:
					position, tokenIndex, depth = position1696, tokenIndex1696, depth1696
					{
						position1707 := position
						depth++
						{
							position1708, tokenIndex1708, depth1708 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1709
							}
							position++
							goto l1708
						l1709:
							position, tokenIndex, depth = position1708, tokenIndex1708, depth1708
							if buffer[position] != rune('F') {
								goto l1694
							}
							position++
						}
					l1708:
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
								goto l1694
							}
							position++
						}
					l1710:
						{
							position1712, tokenIndex1712, depth1712 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1713
							}
							position++
							goto l1712
						l1713:
							position, tokenIndex, depth = position1712, tokenIndex1712, depth1712
							if buffer[position] != rune('L') {
								goto l1694
							}
							position++
						}
					l1712:
						{
							position1714, tokenIndex1714, depth1714 := position, tokenIndex, depth
							if buffer[position] != rune('s') {
								goto l1715
							}
							position++
							goto l1714
						l1715:
							position, tokenIndex, depth = position1714, tokenIndex1714, depth1714
							if buffer[position] != rune('S') {
								goto l1694
							}
							position++
						}
					l1714:
						{
							position1716, tokenIndex1716, depth1716 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1717
							}
							position++
							goto l1716
						l1717:
							position, tokenIndex, depth = position1716, tokenIndex1716, depth1716
							if buffer[position] != rune('E') {
								goto l1694
							}
							position++
						}
					l1716:
						if !rules[ruleskip]() {
							goto l1694
						}
						depth--
						add(ruleFALSE, position1707)
					}
				}
			l1696:
				depth--
				add(rulebooleanLiteral, position1695)
			}
			return true
		l1694:
			position, tokenIndex, depth = position1694, tokenIndex1694, depth1694
			return false
		},
		/* 90 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 91 blankNodeLabel <- <('_' ':' (pnCharsU / [0-9]) (((pnCharsU / ((&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))* pnCharsU) / ([0-9] / '-'))? skip)> */
		nil,
		/* 92 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 93 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1721, tokenIndex1721, depth1721 := position, tokenIndex, depth
			{
				position1722 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1721
				}
				position++
			l1723:
				{
					position1724, tokenIndex1724, depth1724 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1724
					}
					goto l1723
				l1724:
					position, tokenIndex, depth = position1724, tokenIndex1724, depth1724
				}
				if buffer[position] != rune(')') {
					goto l1721
				}
				position++
				if !rules[ruleskip]() {
					goto l1721
				}
				depth--
				add(rulenil, position1722)
			}
			return true
		l1721:
			position, tokenIndex, depth = position1721, tokenIndex1721, depth1721
			return false
		},
		/* 94 VARNAME <- <((pnCharsU / [0-9]) (pnCharsU / ((&('â') [‿-⁀]) | (&('Ì' | 'Í') [̀-ͯ]) | (&('Â') '·') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))*)> */
		nil,
		/* 95 pnPrefix <- <(pnCharsBase pnChars*)> */
		func() bool {
			position1726, tokenIndex1726, depth1726 := position, tokenIndex, depth
			{
				position1727 := position
				depth++
				if !rules[rulepnCharsBase]() {
					goto l1726
				}
			l1728:
				{
					position1729, tokenIndex1729, depth1729 := position, tokenIndex, depth
					{
						position1730 := position
						depth++
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1729
								}
								position++
								break
							case '-':
								if buffer[position] != rune('-') {
									goto l1729
								}
								position++
								break
							default:
								if !rules[rulepnCharsU]() {
									goto l1729
								}
								break
							}
						}

						depth--
						add(rulepnChars, position1730)
					}
					goto l1728
				l1729:
					position, tokenIndex, depth = position1729, tokenIndex1729, depth1729
				}
				depth--
				add(rulepnPrefix, position1727)
			}
			return true
		l1726:
			position, tokenIndex, depth = position1726, tokenIndex1726, depth1726
			return false
		},
		/* 96 pnLocal <- <((&('%' | '\\') plx) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&(':') ':') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))+> */
		nil,
		/* 97 pnChars <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))> */
		nil,
		/* 98 pnCharsU <- <(pnCharsBase / '_')> */
		func() bool {
			position1734, tokenIndex1734, depth1734 := position, tokenIndex, depth
			{
				position1735 := position
				depth++
				{
					position1736, tokenIndex1736, depth1736 := position, tokenIndex, depth
					if !rules[rulepnCharsBase]() {
						goto l1737
					}
					goto l1736
				l1737:
					position, tokenIndex, depth = position1736, tokenIndex1736, depth1736
					if buffer[position] != rune('_') {
						goto l1734
					}
					position++
				}
			l1736:
				depth--
				add(rulepnCharsU, position1735)
			}
			return true
		l1734:
			position, tokenIndex, depth = position1734, tokenIndex1734, depth1734
			return false
		},
		/* 99 pnCharsBase <- <([À-Ö] / [Ø-ö] / [Ͱ-ͽ] / [‌-‍] / [⁰-↏] / [豈-﷏] / ((&('ð' | 'ñ' | 'ò' | 'ó') [𐀀-󯿿]) | (&('ï') [ﷰ-�]) | (&('ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í') [、-퟿]) | (&('â') [Ⰰ-⿯]) | (&('Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á') [Ϳ-῿]) | (&('Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë') [ø-˿]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))> */
		func() bool {
			position1738, tokenIndex1738, depth1738 := position, tokenIndex, depth
			{
				position1739 := position
				depth++
				{
					position1740, tokenIndex1740, depth1740 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('À') || c > rune('Ö') {
						goto l1741
					}
					position++
					goto l1740
				l1741:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
					if c := buffer[position]; c < rune('Ø') || c > rune('ö') {
						goto l1742
					}
					position++
					goto l1740
				l1742:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
					if c := buffer[position]; c < rune('Ͱ') || c > rune('ͽ') {
						goto l1743
					}
					position++
					goto l1740
				l1743:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
					if c := buffer[position]; c < rune('\u200c') || c > rune('\u200d') {
						goto l1744
					}
					position++
					goto l1740
				l1744:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
					if c := buffer[position]; c < rune('⁰') || c > rune('\u218f') {
						goto l1745
					}
					position++
					goto l1740
				l1745:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
					if c := buffer[position]; c < rune('豈') || c > rune('\ufdcf') {
						goto l1746
					}
					position++
					goto l1740
				l1746:
					position, tokenIndex, depth = position1740, tokenIndex1740, depth1740
					{
						switch buffer[position] {
						case 'ð', 'ñ', 'ò', 'ó':
							if c := buffer[position]; c < rune('𐀀') || c > rune('\U000effff') {
								goto l1738
							}
							position++
							break
						case 'ï':
							if c := buffer[position]; c < rune('ﷰ') || c > rune('�') {
								goto l1738
							}
							position++
							break
						case 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í':
							if c := buffer[position]; c < rune('、') || c > rune('\ud7ff') {
								goto l1738
							}
							position++
							break
						case 'â':
							if c := buffer[position]; c < rune('Ⰰ') || c > rune('\u2fef') {
								goto l1738
							}
							position++
							break
						case 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'Þ', 'ß', 'à', 'á':
							if c := buffer[position]; c < rune('\u037f') || c > rune('\u1fff') {
								goto l1738
							}
							position++
							break
						case 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ê', 'Ë':
							if c := buffer[position]; c < rune('ø') || c > rune('˿') {
								goto l1738
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1738
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1738
							}
							position++
							break
						}
					}

				}
			l1740:
				depth--
				add(rulepnCharsBase, position1739)
			}
			return true
		l1738:
			position, tokenIndex, depth = position1738, tokenIndex1738, depth1738
			return false
		},
		/* 100 plx <- <(percent / pnLocalEsc)> */
		nil,
		/* 101 percent <- <('%' hex hex)> */
		nil,
		/* 102 hex <- <((&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))> */
		func() bool {
			position1750, tokenIndex1750, depth1750 := position, tokenIndex, depth
			{
				position1751 := position
				depth++
				{
					switch buffer[position] {
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l1750
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1750
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1750
						}
						position++
						break
					}
				}

				depth--
				add(rulehex, position1751)
			}
			return true
		l1750:
			position, tokenIndex, depth = position1750, tokenIndex1750, depth1750
			return false
		},
		/* 103 pnLocalEsc <- <('\\' ((&('%') '%') | (&('@') '@') | (&('#') '#') | (&('?') '?') | (&('/') '/') | (&('=') '=') | (&(';') ';') | (&(',') ',') | (&('+') '+') | (&('*') '*') | (&(')') ')') | (&('(') '(') | (&('\'') '\'') | (&('&') '&') | (&('$') '$') | (&('!') '!') | (&('-') '-') | (&('.') '.') | (&('~') '~') | (&('_') '_')))> */
		nil,
		/* 104 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 105 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 106 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 107 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 108 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 109 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 110 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1760, tokenIndex1760, depth1760 := position, tokenIndex, depth
			{
				position1761 := position
				depth++
				{
					position1762, tokenIndex1762, depth1762 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1763
					}
					position++
					goto l1762
				l1763:
					position, tokenIndex, depth = position1762, tokenIndex1762, depth1762
					if buffer[position] != rune('D') {
						goto l1760
					}
					position++
				}
			l1762:
				{
					position1764, tokenIndex1764, depth1764 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1765
					}
					position++
					goto l1764
				l1765:
					position, tokenIndex, depth = position1764, tokenIndex1764, depth1764
					if buffer[position] != rune('I') {
						goto l1760
					}
					position++
				}
			l1764:
				{
					position1766, tokenIndex1766, depth1766 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1767
					}
					position++
					goto l1766
				l1767:
					position, tokenIndex, depth = position1766, tokenIndex1766, depth1766
					if buffer[position] != rune('S') {
						goto l1760
					}
					position++
				}
			l1766:
				{
					position1768, tokenIndex1768, depth1768 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1769
					}
					position++
					goto l1768
				l1769:
					position, tokenIndex, depth = position1768, tokenIndex1768, depth1768
					if buffer[position] != rune('T') {
						goto l1760
					}
					position++
				}
			l1768:
				{
					position1770, tokenIndex1770, depth1770 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1771
					}
					position++
					goto l1770
				l1771:
					position, tokenIndex, depth = position1770, tokenIndex1770, depth1770
					if buffer[position] != rune('I') {
						goto l1760
					}
					position++
				}
			l1770:
				{
					position1772, tokenIndex1772, depth1772 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1773
					}
					position++
					goto l1772
				l1773:
					position, tokenIndex, depth = position1772, tokenIndex1772, depth1772
					if buffer[position] != rune('N') {
						goto l1760
					}
					position++
				}
			l1772:
				{
					position1774, tokenIndex1774, depth1774 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1775
					}
					position++
					goto l1774
				l1775:
					position, tokenIndex, depth = position1774, tokenIndex1774, depth1774
					if buffer[position] != rune('C') {
						goto l1760
					}
					position++
				}
			l1774:
				{
					position1776, tokenIndex1776, depth1776 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1777
					}
					position++
					goto l1776
				l1777:
					position, tokenIndex, depth = position1776, tokenIndex1776, depth1776
					if buffer[position] != rune('T') {
						goto l1760
					}
					position++
				}
			l1776:
				if !rules[ruleskip]() {
					goto l1760
				}
				depth--
				add(ruleDISTINCT, position1761)
			}
			return true
		l1760:
			position, tokenIndex, depth = position1760, tokenIndex1760, depth1760
			return false
		},
		/* 111 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 112 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 113 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 114 LBRACE <- <('{' skip)> */
		func() bool {
			position1781, tokenIndex1781, depth1781 := position, tokenIndex, depth
			{
				position1782 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1781
				}
				position++
				if !rules[ruleskip]() {
					goto l1781
				}
				depth--
				add(ruleLBRACE, position1782)
			}
			return true
		l1781:
			position, tokenIndex, depth = position1781, tokenIndex1781, depth1781
			return false
		},
		/* 115 RBRACE <- <('}' skip)> */
		func() bool {
			position1783, tokenIndex1783, depth1783 := position, tokenIndex, depth
			{
				position1784 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1783
				}
				position++
				if !rules[ruleskip]() {
					goto l1783
				}
				depth--
				add(ruleRBRACE, position1784)
			}
			return true
		l1783:
			position, tokenIndex, depth = position1783, tokenIndex1783, depth1783
			return false
		},
		/* 116 LBRACK <- <('[' skip)> */
		nil,
		/* 117 RBRACK <- <(']' skip)> */
		nil,
		/* 118 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1787, tokenIndex1787, depth1787 := position, tokenIndex, depth
			{
				position1788 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1787
				}
				position++
				if !rules[ruleskip]() {
					goto l1787
				}
				depth--
				add(ruleSEMICOLON, position1788)
			}
			return true
		l1787:
			position, tokenIndex, depth = position1787, tokenIndex1787, depth1787
			return false
		},
		/* 119 COMMA <- <(',' skip)> */
		func() bool {
			position1789, tokenIndex1789, depth1789 := position, tokenIndex, depth
			{
				position1790 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1789
				}
				position++
				if !rules[ruleskip]() {
					goto l1789
				}
				depth--
				add(ruleCOMMA, position1790)
			}
			return true
		l1789:
			position, tokenIndex, depth = position1789, tokenIndex1789, depth1789
			return false
		},
		/* 120 DOT <- <('.' skip)> */
		func() bool {
			position1791, tokenIndex1791, depth1791 := position, tokenIndex, depth
			{
				position1792 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1791
				}
				position++
				if !rules[ruleskip]() {
					goto l1791
				}
				depth--
				add(ruleDOT, position1792)
			}
			return true
		l1791:
			position, tokenIndex, depth = position1791, tokenIndex1791, depth1791
			return false
		},
		/* 121 COLON <- <(':' skip)> */
		nil,
		/* 122 PIPE <- <('|' skip)> */
		func() bool {
			position1794, tokenIndex1794, depth1794 := position, tokenIndex, depth
			{
				position1795 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1794
				}
				position++
				if !rules[ruleskip]() {
					goto l1794
				}
				depth--
				add(rulePIPE, position1795)
			}
			return true
		l1794:
			position, tokenIndex, depth = position1794, tokenIndex1794, depth1794
			return false
		},
		/* 123 SLASH <- <('/' skip)> */
		func() bool {
			position1796, tokenIndex1796, depth1796 := position, tokenIndex, depth
			{
				position1797 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1796
				}
				position++
				if !rules[ruleskip]() {
					goto l1796
				}
				depth--
				add(ruleSLASH, position1797)
			}
			return true
		l1796:
			position, tokenIndex, depth = position1796, tokenIndex1796, depth1796
			return false
		},
		/* 124 INVERSE <- <('^' skip)> */
		func() bool {
			position1798, tokenIndex1798, depth1798 := position, tokenIndex, depth
			{
				position1799 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1798
				}
				position++
				if !rules[ruleskip]() {
					goto l1798
				}
				depth--
				add(ruleINVERSE, position1799)
			}
			return true
		l1798:
			position, tokenIndex, depth = position1798, tokenIndex1798, depth1798
			return false
		},
		/* 125 LPAREN <- <('(' skip)> */
		func() bool {
			position1800, tokenIndex1800, depth1800 := position, tokenIndex, depth
			{
				position1801 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1800
				}
				position++
				if !rules[ruleskip]() {
					goto l1800
				}
				depth--
				add(ruleLPAREN, position1801)
			}
			return true
		l1800:
			position, tokenIndex, depth = position1800, tokenIndex1800, depth1800
			return false
		},
		/* 126 RPAREN <- <(')' skip)> */
		func() bool {
			position1802, tokenIndex1802, depth1802 := position, tokenIndex, depth
			{
				position1803 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1802
				}
				position++
				if !rules[ruleskip]() {
					goto l1802
				}
				depth--
				add(ruleRPAREN, position1803)
			}
			return true
		l1802:
			position, tokenIndex, depth = position1802, tokenIndex1802, depth1802
			return false
		},
		/* 127 ISA <- <('a' skip)> */
		func() bool {
			position1804, tokenIndex1804, depth1804 := position, tokenIndex, depth
			{
				position1805 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1804
				}
				position++
				if !rules[ruleskip]() {
					goto l1804
				}
				depth--
				add(ruleISA, position1805)
			}
			return true
		l1804:
			position, tokenIndex, depth = position1804, tokenIndex1804, depth1804
			return false
		},
		/* 128 NOT <- <('!' skip)> */
		func() bool {
			position1806, tokenIndex1806, depth1806 := position, tokenIndex, depth
			{
				position1807 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1806
				}
				position++
				if !rules[ruleskip]() {
					goto l1806
				}
				depth--
				add(ruleNOT, position1807)
			}
			return true
		l1806:
			position, tokenIndex, depth = position1806, tokenIndex1806, depth1806
			return false
		},
		/* 129 STAR <- <('*' skip)> */
		func() bool {
			position1808, tokenIndex1808, depth1808 := position, tokenIndex, depth
			{
				position1809 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1808
				}
				position++
				if !rules[ruleskip]() {
					goto l1808
				}
				depth--
				add(ruleSTAR, position1809)
			}
			return true
		l1808:
			position, tokenIndex, depth = position1808, tokenIndex1808, depth1808
			return false
		},
		/* 130 QUESTION <- <('?' skip)> */
		nil,
		/* 131 PLUS <- <('+' skip)> */
		func() bool {
			position1811, tokenIndex1811, depth1811 := position, tokenIndex, depth
			{
				position1812 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1811
				}
				position++
				if !rules[ruleskip]() {
					goto l1811
				}
				depth--
				add(rulePLUS, position1812)
			}
			return true
		l1811:
			position, tokenIndex, depth = position1811, tokenIndex1811, depth1811
			return false
		},
		/* 132 MINUS <- <('-' skip)> */
		func() bool {
			position1813, tokenIndex1813, depth1813 := position, tokenIndex, depth
			{
				position1814 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1813
				}
				position++
				if !rules[ruleskip]() {
					goto l1813
				}
				depth--
				add(ruleMINUS, position1814)
			}
			return true
		l1813:
			position, tokenIndex, depth = position1813, tokenIndex1813, depth1813
			return false
		},
		/* 133 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 134 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 135 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 136 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 137 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1819, tokenIndex1819, depth1819 := position, tokenIndex, depth
			{
				position1820 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1819
				}
				position++
			l1821:
				{
					position1822, tokenIndex1822, depth1822 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1822
					}
					position++
					goto l1821
				l1822:
					position, tokenIndex, depth = position1822, tokenIndex1822, depth1822
				}
				if !rules[ruleskip]() {
					goto l1819
				}
				depth--
				add(ruleINTEGER, position1820)
			}
			return true
		l1819:
			position, tokenIndex, depth = position1819, tokenIndex1819, depth1819
			return false
		},
		/* 138 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 139 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 140 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 141 OR <- <('|' '|' skip)> */
		nil,
		/* 142 AND <- <('&' '&' skip)> */
		nil,
		/* 143 EQ <- <('=' skip)> */
		func() bool {
			position1828, tokenIndex1828, depth1828 := position, tokenIndex, depth
			{
				position1829 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1828
				}
				position++
				if !rules[ruleskip]() {
					goto l1828
				}
				depth--
				add(ruleEQ, position1829)
			}
			return true
		l1828:
			position, tokenIndex, depth = position1828, tokenIndex1828, depth1828
			return false
		},
		/* 144 NE <- <('!' '=' skip)> */
		nil,
		/* 145 GT <- <('>' skip)> */
		nil,
		/* 146 LT <- <('<' skip)> */
		nil,
		/* 147 LE <- <('<' '=' skip)> */
		nil,
		/* 148 GE <- <('>' '=' skip)> */
		nil,
		/* 149 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 150 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 151 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1837, tokenIndex1837, depth1837 := position, tokenIndex, depth
			{
				position1838 := position
				depth++
				{
					position1839, tokenIndex1839, depth1839 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1840
					}
					position++
					goto l1839
				l1840:
					position, tokenIndex, depth = position1839, tokenIndex1839, depth1839
					if buffer[position] != rune('A') {
						goto l1837
					}
					position++
				}
			l1839:
				{
					position1841, tokenIndex1841, depth1841 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1842
					}
					position++
					goto l1841
				l1842:
					position, tokenIndex, depth = position1841, tokenIndex1841, depth1841
					if buffer[position] != rune('S') {
						goto l1837
					}
					position++
				}
			l1841:
				if !rules[ruleskip]() {
					goto l1837
				}
				depth--
				add(ruleAS, position1838)
			}
			return true
		l1837:
			position, tokenIndex, depth = position1837, tokenIndex1837, depth1837
			return false
		},
		/* 152 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 153 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 154 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 155 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 156 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 157 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 158 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 159 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 160 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 161 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 162 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 163 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 164 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 165 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 166 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 167 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 168 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 169 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 170 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 171 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 172 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 173 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 174 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 175 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 176 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 177 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 178 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 179 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 180 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 181 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 182 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 183 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 184 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 185 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 186 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 187 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 188 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 189 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 190 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 191 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 192 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 193 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 194 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 195 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 196 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 197 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 198 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 199 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 200 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 201 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 202 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 203 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 204 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 205 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 206 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 207 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 208 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 209 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 210 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 211 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 212 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 213 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 214 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 215 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 216 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 217 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 218 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 219 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 220 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1911, tokenIndex1911, depth1911 := position, tokenIndex, depth
			{
				position1912 := position
				depth++
				{
					position1913, tokenIndex1913, depth1913 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1914
					}
					position++
					goto l1913
				l1914:
					position, tokenIndex, depth = position1913, tokenIndex1913, depth1913
					if buffer[position] != rune('B') {
						goto l1911
					}
					position++
				}
			l1913:
				{
					position1915, tokenIndex1915, depth1915 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1916
					}
					position++
					goto l1915
				l1916:
					position, tokenIndex, depth = position1915, tokenIndex1915, depth1915
					if buffer[position] != rune('Y') {
						goto l1911
					}
					position++
				}
			l1915:
				if !rules[ruleskip]() {
					goto l1911
				}
				depth--
				add(ruleBY, position1912)
			}
			return true
		l1911:
			position, tokenIndex, depth = position1911, tokenIndex1911, depth1911
			return false
		},
		/* 221 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 222 GRAPH <- <(('g' / 'G') ('r' / 'R') ('a' / 'A') ('p' / 'P') ('h' / 'H') skip)> */
		nil,
		/* 223 MINUSSETOPER <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('s' / 'S') skip)> */
		nil,
		/* 224 skip <- <(<(ws / comment)*> Action13)> */
		func() bool {
			{
				position1921 := position
				depth++
				{
					position1922 := position
					depth++
				l1923:
					{
						position1924, tokenIndex1924, depth1924 := position, tokenIndex, depth
						{
							position1925, tokenIndex1925, depth1925 := position, tokenIndex, depth
							if !rules[rulews]() {
								goto l1926
							}
							goto l1925
						l1926:
							position, tokenIndex, depth = position1925, tokenIndex1925, depth1925
							{
								position1927 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1924
								}
								position++
							l1928:
								{
									position1929, tokenIndex1929, depth1929 := position, tokenIndex, depth
									{
										position1930, tokenIndex1930, depth1930 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1930
										}
										goto l1929
									l1930:
										position, tokenIndex, depth = position1930, tokenIndex1930, depth1930
									}
									if !matchDot() {
										goto l1929
									}
									goto l1928
								l1929:
									position, tokenIndex, depth = position1929, tokenIndex1929, depth1929
								}
								if !rules[ruleendOfLine]() {
									goto l1924
								}
								depth--
								add(rulecomment, position1927)
							}
						}
					l1925:
						goto l1923
					l1924:
						position, tokenIndex, depth = position1924, tokenIndex1924, depth1924
					}
					depth--
					add(rulePegText, position1922)
				}
				{
					add(ruleAction13, position)
				}
				depth--
				add(ruleskip, position1921)
			}
			return true
		},
		/* 225 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1932, tokenIndex1932, depth1932 := position, tokenIndex, depth
			{
				position1933 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1932
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1932
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1932
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1932
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1932
						}
						break
					}
				}

				depth--
				add(rulews, position1933)
			}
			return true
		l1932:
			position, tokenIndex, depth = position1932, tokenIndex1932, depth1932
			return false
		},
		/* 226 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 227 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1936, tokenIndex1936, depth1936 := position, tokenIndex, depth
			{
				position1937 := position
				depth++
				{
					position1938, tokenIndex1938, depth1938 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1939
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1939
					}
					position++
					goto l1938
				l1939:
					position, tokenIndex, depth = position1938, tokenIndex1938, depth1938
					if buffer[position] != rune('\n') {
						goto l1940
					}
					position++
					goto l1938
				l1940:
					position, tokenIndex, depth = position1938, tokenIndex1938, depth1938
					if buffer[position] != rune('\r') {
						goto l1936
					}
					position++
				}
			l1938:
				depth--
				add(ruleendOfLine, position1937)
			}
			return true
		l1936:
			position, tokenIndex, depth = position1936, tokenIndex1936, depth1936
			return false
		},
		nil,
		/* 230 Action0 <- <{ p.addPrefix(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 231 SERVICE <- <> */
		nil,
		/* 232 SILENT <- <> */
		nil,
		/* 233 Action1 <- <{ p.S = p.skipped(buffer, begin, end) }> */
		nil,
		/* 234 Action2 <- <{ p.S = p.skipped(buffer, begin, end) }> */
		nil,
		/* 235 Action3 <- <{ p.S = "?POF" }> */
		nil,
		/* 236 Action4 <- <{ p.P = p.skipped(buffer, begin, end) }> */
		nil,
		/* 237 Action5 <- <{ p.P = "?POF" }> */
		nil,
		/* 238 Action6 <- <{ p.P = p.skipped(buffer, begin, end) }> */
		nil,
		/* 239 Action7 <- <{ p.O = "?FillVar"; p.addTriplePattern() }> */
		nil,
		/* 240 Action8 <- <{ p.O = "?POF"; p.addTriplePattern() }> */
		nil,
		/* 241 Action9 <- <{ p.O = p.skipped(buffer, begin, end); p.addTriplePattern() }> */
		nil,
		/* 242 Action10 <- <{ p.setPrefix(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 243 Action11 <- <{ p.setPathLength(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 244 Action12 <- <{ p.setKeyword(p.skipped(buffer, begin, end)) }> */
		nil,
		/* 245 Action13 <- <{ p.skipBegin = begin }> */
		nil,
	}
	p.rules = rules
}
