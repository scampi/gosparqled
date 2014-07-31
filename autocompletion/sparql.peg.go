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
	rules  [233]func() bool
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
		/* 26 triplesBlock <- <triplesSameSubjectPath+> */
		func() bool {
			position327, tokenIndex327, depth327 := position, tokenIndex, depth
			{
				position328 := position
				depth++
				{
					position331 := position
					depth++
					{
						position332, tokenIndex332, depth332 := position, tokenIndex, depth
						{
							position334 := position
							depth++
							{
								position335, tokenIndex335, depth335 := position, tokenIndex, depth
								{
									position337 := position
									depth++
									if !rules[rulevar]() {
										goto l336
									}
									depth--
									add(rulePegText, position337)
								}
								{
									add(ruleAction1, position)
								}
								goto l335
							l336:
								position, tokenIndex, depth = position335, tokenIndex335, depth335
								{
									position340 := position
									depth++
									if !rules[rulegraphTerm]() {
										goto l339
									}
									depth--
									add(rulePegText, position340)
								}
								{
									add(ruleAction2, position)
								}
								goto l335
							l339:
								position, tokenIndex, depth = position335, tokenIndex335, depth335
								if !rules[rulepof]() {
									goto l333
								}
								{
									add(ruleAction3, position)
								}
							}
						l335:
							depth--
							add(rulevarOrTerm, position334)
						}
						if !rules[rulepropertyListPath]() {
							goto l333
						}
						goto l332
					l333:
						position, tokenIndex, depth = position332, tokenIndex332, depth332
						if !rules[ruletriplesNodePath]() {
							goto l327
						}
						{
							position343, tokenIndex343, depth343 := position, tokenIndex, depth
							if !rules[rulepropertyListPath]() {
								goto l343
							}
							goto l344
						l343:
							position, tokenIndex, depth = position343, tokenIndex343, depth343
						}
					l344:
					}
				l332:
					{
						position345, tokenIndex345, depth345 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l345
						}
						goto l346
					l345:
						position, tokenIndex, depth = position345, tokenIndex345, depth345
					}
				l346:
					depth--
					add(ruletriplesSameSubjectPath, position331)
				}
			l329:
				{
					position330, tokenIndex330, depth330 := position, tokenIndex, depth
					{
						position347 := position
						depth++
						{
							position348, tokenIndex348, depth348 := position, tokenIndex, depth
							{
								position350 := position
								depth++
								{
									position351, tokenIndex351, depth351 := position, tokenIndex, depth
									{
										position353 := position
										depth++
										if !rules[rulevar]() {
											goto l352
										}
										depth--
										add(rulePegText, position353)
									}
									{
										add(ruleAction1, position)
									}
									goto l351
								l352:
									position, tokenIndex, depth = position351, tokenIndex351, depth351
									{
										position356 := position
										depth++
										if !rules[rulegraphTerm]() {
											goto l355
										}
										depth--
										add(rulePegText, position356)
									}
									{
										add(ruleAction2, position)
									}
									goto l351
								l355:
									position, tokenIndex, depth = position351, tokenIndex351, depth351
									if !rules[rulepof]() {
										goto l349
									}
									{
										add(ruleAction3, position)
									}
								}
							l351:
								depth--
								add(rulevarOrTerm, position350)
							}
							if !rules[rulepropertyListPath]() {
								goto l349
							}
							goto l348
						l349:
							position, tokenIndex, depth = position348, tokenIndex348, depth348
							if !rules[ruletriplesNodePath]() {
								goto l330
							}
							{
								position359, tokenIndex359, depth359 := position, tokenIndex, depth
								if !rules[rulepropertyListPath]() {
									goto l359
								}
								goto l360
							l359:
								position, tokenIndex, depth = position359, tokenIndex359, depth359
							}
						l360:
						}
					l348:
						{
							position361, tokenIndex361, depth361 := position, tokenIndex, depth
							if !rules[ruleDOT]() {
								goto l361
							}
							goto l362
						l361:
							position, tokenIndex, depth = position361, tokenIndex361, depth361
						}
					l362:
						depth--
						add(ruletriplesSameSubjectPath, position347)
					}
					goto l329
				l330:
					position, tokenIndex, depth = position330, tokenIndex330, depth330
				}
				depth--
				add(ruletriplesBlock, position328)
			}
			return true
		l327:
			position, tokenIndex, depth = position327, tokenIndex327, depth327
			return false
		},
		/* 27 triplesSameSubjectPath <- <(((varOrTerm propertyListPath) / (triplesNodePath propertyListPath?)) DOT?)> */
		nil,
		/* 28 varOrTerm <- <((<var> Action1) / (<graphTerm> Action2) / (pof Action3))> */
		nil,
		/* 29 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		func() bool {
			position365, tokenIndex365, depth365 := position, tokenIndex, depth
			{
				position366 := position
				depth++
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l368
					}
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					{
						switch buffer[position] {
						case '(':
							if !rules[rulenil]() {
								goto l365
							}
							break
						case '[', '_':
							{
								position370 := position
								depth++
								{
									position371, tokenIndex371, depth371 := position, tokenIndex, depth
									{
										position373 := position
										depth++
										if buffer[position] != rune('_') {
											goto l372
										}
										position++
										if buffer[position] != rune(':') {
											goto l372
										}
										position++
										{
											position374, tokenIndex374, depth374 := position, tokenIndex, depth
											if !rules[rulepnCharsU]() {
												goto l375
											}
											goto l374
										l375:
											position, tokenIndex, depth = position374, tokenIndex374, depth374
											if c := buffer[position]; c < rune('0') || c > rune('9') {
												goto l372
											}
											position++
										}
									l374:
										{
											position376, tokenIndex376, depth376 := position, tokenIndex, depth
											{
												position378, tokenIndex378, depth378 := position, tokenIndex, depth
											l380:
												{
													position381, tokenIndex381, depth381 := position, tokenIndex, depth
													{
														position382, tokenIndex382, depth382 := position, tokenIndex, depth
														if !rules[rulepnCharsU]() {
															goto l383
														}
														goto l382
													l383:
														position, tokenIndex, depth = position382, tokenIndex382, depth382
														{
															switch buffer[position] {
															case '.':
																if buffer[position] != rune('.') {
																	goto l381
																}
																position++
																break
															case '-':
																if buffer[position] != rune('-') {
																	goto l381
																}
																position++
																break
															default:
																if c := buffer[position]; c < rune('0') || c > rune('9') {
																	goto l381
																}
																position++
																break
															}
														}

													}
												l382:
													goto l380
												l381:
													position, tokenIndex, depth = position381, tokenIndex381, depth381
												}
												if !rules[rulepnCharsU]() {
													goto l379
												}
												goto l378
											l379:
												position, tokenIndex, depth = position378, tokenIndex378, depth378
												{
													position385, tokenIndex385, depth385 := position, tokenIndex, depth
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l386
													}
													position++
													goto l385
												l386:
													position, tokenIndex, depth = position385, tokenIndex385, depth385
													if buffer[position] != rune('-') {
														goto l376
													}
													position++
												}
											l385:
											}
										l378:
											goto l377
										l376:
											position, tokenIndex, depth = position376, tokenIndex376, depth376
										}
									l377:
										if !rules[ruleskip]() {
											goto l372
										}
										depth--
										add(ruleblankNodeLabel, position373)
									}
									goto l371
								l372:
									position, tokenIndex, depth = position371, tokenIndex371, depth371
									{
										position387 := position
										depth++
										if buffer[position] != rune('[') {
											goto l365
										}
										position++
									l388:
										{
											position389, tokenIndex389, depth389 := position, tokenIndex, depth
											if !rules[rulews]() {
												goto l389
											}
											goto l388
										l389:
											position, tokenIndex, depth = position389, tokenIndex389, depth389
										}
										if buffer[position] != rune(']') {
											goto l365
										}
										position++
										if !rules[ruleskip]() {
											goto l365
										}
										depth--
										add(ruleanon, position387)
									}
								}
							l371:
								depth--
								add(ruleblankNode, position370)
							}
							break
						case 'F', 'T', 'f', 't':
							if !rules[rulebooleanLiteral]() {
								goto l365
							}
							break
						case '"':
							if !rules[ruleliteral]() {
								goto l365
							}
							break
						default:
							if !rules[rulenumericLiteral]() {
								goto l365
							}
							break
						}
					}

				}
			l367:
				depth--
				add(rulegraphTerm, position366)
			}
			return true
		l365:
			position, tokenIndex, depth = position365, tokenIndex365, depth365
			return false
		},
		/* 30 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		func() bool {
			position390, tokenIndex390, depth390 := position, tokenIndex, depth
			{
				position391 := position
				depth++
				{
					position392, tokenIndex392, depth392 := position, tokenIndex, depth
					{
						position394 := position
						depth++
						if !rules[ruleLPAREN]() {
							goto l393
						}
						if !rules[rulegraphNodePath]() {
							goto l393
						}
					l395:
						{
							position396, tokenIndex396, depth396 := position, tokenIndex, depth
							if !rules[rulegraphNodePath]() {
								goto l396
							}
							goto l395
						l396:
							position, tokenIndex, depth = position396, tokenIndex396, depth396
						}
						if !rules[ruleRPAREN]() {
							goto l393
						}
						depth--
						add(rulecollectionPath, position394)
					}
					goto l392
				l393:
					position, tokenIndex, depth = position392, tokenIndex392, depth392
					{
						position397 := position
						depth++
						{
							position398 := position
							depth++
							if buffer[position] != rune('[') {
								goto l390
							}
							position++
							if !rules[ruleskip]() {
								goto l390
							}
							depth--
							add(ruleLBRACK, position398)
						}
						if !rules[rulepropertyListPath]() {
							goto l390
						}
						{
							position399 := position
							depth++
							if buffer[position] != rune(']') {
								goto l390
							}
							position++
							if !rules[ruleskip]() {
								goto l390
							}
							depth--
							add(ruleRBRACK, position399)
						}
						depth--
						add(ruleblankNodePropertyListPath, position397)
					}
				}
			l392:
				depth--
				add(ruletriplesNodePath, position391)
			}
			return true
		l390:
			position, tokenIndex, depth = position390, tokenIndex390, depth390
			return false
		},
		/* 31 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 32 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 33 propertyListPath <- <(((pof Action4) / (<var> Action5) / verbPath) objectListPath (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position402, tokenIndex402, depth402 := position, tokenIndex, depth
			{
				position403 := position
				depth++
				{
					position404, tokenIndex404, depth404 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l405
					}
					{
						add(ruleAction4, position)
					}
					goto l404
				l405:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					{
						position408 := position
						depth++
						if !rules[rulevar]() {
							goto l407
						}
						depth--
						add(rulePegText, position408)
					}
					{
						add(ruleAction5, position)
					}
					goto l404
				l407:
					position, tokenIndex, depth = position404, tokenIndex404, depth404
					{
						position410 := position
						depth++
						if !rules[rulepath]() {
							goto l402
						}
						depth--
						add(ruleverbPath, position410)
					}
				}
			l404:
				{
					position411 := position
					depth++
					if !rules[ruleobjectPath]() {
						goto l402
					}
				l412:
					{
						position413, tokenIndex413, depth413 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l413
						}
						if !rules[ruleobjectPath]() {
							goto l413
						}
						goto l412
					l413:
						position, tokenIndex, depth = position413, tokenIndex413, depth413
					}
					depth--
					add(ruleobjectListPath, position411)
				}
				{
					position414, tokenIndex414, depth414 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l414
					}
					{
						position416, tokenIndex416, depth416 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l416
						}
						goto l417
					l416:
						position, tokenIndex, depth = position416, tokenIndex416, depth416
					}
				l417:
					goto l415
				l414:
					position, tokenIndex, depth = position414, tokenIndex414, depth414
				}
			l415:
				depth--
				add(rulepropertyListPath, position403)
			}
			return true
		l402:
			position, tokenIndex, depth = position402, tokenIndex402, depth402
			return false
		},
		/* 34 verbPath <- <path> */
		nil,
		/* 35 path <- <pathAlternative> */
		func() bool {
			position419, tokenIndex419, depth419 := position, tokenIndex, depth
			{
				position420 := position
				depth++
				{
					position421 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l419
					}
				l422:
					{
						position423, tokenIndex423, depth423 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l423
						}
						if !rules[rulepathSequence]() {
							goto l423
						}
						goto l422
					l423:
						position, tokenIndex, depth = position423, tokenIndex423, depth423
					}
					depth--
					add(rulepathAlternative, position421)
				}
				depth--
				add(rulepath, position420)
			}
			return true
		l419:
			position, tokenIndex, depth = position419, tokenIndex419, depth419
			return false
		},
		/* 36 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 37 pathSequence <- <(<pathElt> Action6 (SLASH pathSequence)*)> */
		func() bool {
			position425, tokenIndex425, depth425 := position, tokenIndex, depth
			{
				position426 := position
				depth++
				{
					position427 := position
					depth++
					{
						position428 := position
						depth++
						{
							position429, tokenIndex429, depth429 := position, tokenIndex, depth
							if !rules[ruleINVERSE]() {
								goto l429
							}
							goto l430
						l429:
							position, tokenIndex, depth = position429, tokenIndex429, depth429
						}
					l430:
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
								{
									switch buffer[position] {
									case '(':
										if !rules[ruleLPAREN]() {
											goto l425
										}
										if !rules[rulepath]() {
											goto l425
										}
										if !rules[ruleRPAREN]() {
											goto l425
										}
										break
									case '!':
										if !rules[ruleNOT]() {
											goto l425
										}
										{
											position435 := position
											depth++
											{
												position436, tokenIndex436, depth436 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l437
												}
												goto l436
											l437:
												position, tokenIndex, depth = position436, tokenIndex436, depth436
												if !rules[ruleLPAREN]() {
													goto l425
												}
												{
													position438, tokenIndex438, depth438 := position, tokenIndex, depth
													if !rules[rulepathOneInPropertySet]() {
														goto l438
													}
												l440:
													{
														position441, tokenIndex441, depth441 := position, tokenIndex, depth
														if !rules[rulePIPE]() {
															goto l441
														}
														if !rules[rulepathOneInPropertySet]() {
															goto l441
														}
														goto l440
													l441:
														position, tokenIndex, depth = position441, tokenIndex441, depth441
													}
													goto l439
												l438:
													position, tokenIndex, depth = position438, tokenIndex438, depth438
												}
											l439:
												if !rules[ruleRPAREN]() {
													goto l425
												}
											}
										l436:
											depth--
											add(rulepathNegatedPropertySet, position435)
										}
										break
									default:
										if !rules[ruleISA]() {
											goto l425
										}
										break
									}
								}

							}
						l432:
							depth--
							add(rulepathPrimary, position431)
						}
						{
							position442, tokenIndex442, depth442 := position, tokenIndex, depth
							{
								position444 := position
								depth++
								{
									switch buffer[position] {
									case '+':
										if !rules[rulePLUS]() {
											goto l442
										}
										break
									case '?':
										{
											position446 := position
											depth++
											if buffer[position] != rune('?') {
												goto l442
											}
											position++
											if !rules[ruleskip]() {
												goto l442
											}
											depth--
											add(ruleQUESTION, position446)
										}
										break
									default:
										if !rules[ruleSTAR]() {
											goto l442
										}
										break
									}
								}

								{
									position447, tokenIndex447, depth447 := position, tokenIndex, depth
									if !rules[ruleskip]() {
										goto l447
									}
									goto l442
								l447:
									position, tokenIndex, depth = position447, tokenIndex447, depth447
								}
								depth--
								add(rulepathMod, position444)
							}
							goto l443
						l442:
							position, tokenIndex, depth = position442, tokenIndex442, depth442
						}
					l443:
						depth--
						add(rulepathElt, position428)
					}
					depth--
					add(rulePegText, position427)
				}
				{
					add(ruleAction6, position)
				}
			l449:
				{
					position450, tokenIndex450, depth450 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l450
					}
					if !rules[rulepathSequence]() {
						goto l450
					}
					goto l449
				l450:
					position, tokenIndex, depth = position450, tokenIndex450, depth450
				}
				depth--
				add(rulepathSequence, position426)
			}
			return true
		l425:
			position, tokenIndex, depth = position425, tokenIndex425, depth425
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
			position454, tokenIndex454, depth454 := position, tokenIndex, depth
			{
				position455 := position
				depth++
				{
					position456, tokenIndex456, depth456 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l457
					}
					goto l456
				l457:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if !rules[ruleISA]() {
						goto l458
					}
					goto l456
				l458:
					position, tokenIndex, depth = position456, tokenIndex456, depth456
					if !rules[ruleINVERSE]() {
						goto l454
					}
					{
						position459, tokenIndex459, depth459 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l460
						}
						goto l459
					l460:
						position, tokenIndex, depth = position459, tokenIndex459, depth459
						if !rules[ruleISA]() {
							goto l454
						}
					}
				l459:
				}
			l456:
				depth--
				add(rulepathOneInPropertySet, position455)
			}
			return true
		l454:
			position, tokenIndex, depth = position454, tokenIndex454, depth454
			return false
		},
		/* 42 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !skip)> */
		nil,
		/* 43 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 44 objectPath <- <((pof Action7) / (<graphNodePath> Action8) / Action9)> */
		func() bool {
			{
				position464 := position
				depth++
				{
					position465, tokenIndex465, depth465 := position, tokenIndex, depth
					if !rules[rulepof]() {
						goto l466
					}
					{
						add(ruleAction7, position)
					}
					goto l465
				l466:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					{
						position469 := position
						depth++
						if !rules[rulegraphNodePath]() {
							goto l468
						}
						depth--
						add(rulePegText, position469)
					}
					{
						add(ruleAction8, position)
					}
					goto l465
				l468:
					position, tokenIndex, depth = position465, tokenIndex465, depth465
					{
						add(ruleAction9, position)
					}
				}
			l465:
				depth--
				add(ruleobjectPath, position464)
			}
			return true
		},
		/* 45 graphNodePath <- <(var / graphTerm / triplesNodePath)> */
		func() bool {
			position472, tokenIndex472, depth472 := position, tokenIndex, depth
			{
				position473 := position
				depth++
				{
					position474, tokenIndex474, depth474 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l475
					}
					goto l474
				l475:
					position, tokenIndex, depth = position474, tokenIndex474, depth474
					if !rules[rulegraphTerm]() {
						goto l476
					}
					goto l474
				l476:
					position, tokenIndex, depth = position474, tokenIndex474, depth474
					if !rules[ruletriplesNodePath]() {
						goto l472
					}
				}
			l474:
				depth--
				add(rulegraphNodePath, position473)
			}
			return true
		l472:
			position, tokenIndex, depth = position472, tokenIndex472, depth472
			return false
		},
		/* 46 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position478 := position
				depth++
				{
					position479, tokenIndex479, depth479 := position, tokenIndex, depth
					{
						position481, tokenIndex481, depth481 := position, tokenIndex, depth
						{
							position483 := position
							depth++
							{
								position484, tokenIndex484, depth484 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l485
								}
								position++
								goto l484
							l485:
								position, tokenIndex, depth = position484, tokenIndex484, depth484
								if buffer[position] != rune('O') {
									goto l482
								}
								position++
							}
						l484:
							{
								position486, tokenIndex486, depth486 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l487
								}
								position++
								goto l486
							l487:
								position, tokenIndex, depth = position486, tokenIndex486, depth486
								if buffer[position] != rune('R') {
									goto l482
								}
								position++
							}
						l486:
							{
								position488, tokenIndex488, depth488 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l489
								}
								position++
								goto l488
							l489:
								position, tokenIndex, depth = position488, tokenIndex488, depth488
								if buffer[position] != rune('D') {
									goto l482
								}
								position++
							}
						l488:
							{
								position490, tokenIndex490, depth490 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l491
								}
								position++
								goto l490
							l491:
								position, tokenIndex, depth = position490, tokenIndex490, depth490
								if buffer[position] != rune('E') {
									goto l482
								}
								position++
							}
						l490:
							{
								position492, tokenIndex492, depth492 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l493
								}
								position++
								goto l492
							l493:
								position, tokenIndex, depth = position492, tokenIndex492, depth492
								if buffer[position] != rune('R') {
									goto l482
								}
								position++
							}
						l492:
							if !rules[ruleskip]() {
								goto l482
							}
							depth--
							add(ruleORDER, position483)
						}
						if !rules[ruleBY]() {
							goto l482
						}
						{
							position496 := position
							depth++
							{
								position497, tokenIndex497, depth497 := position, tokenIndex, depth
								{
									position499, tokenIndex499, depth499 := position, tokenIndex, depth
									{
										position501, tokenIndex501, depth501 := position, tokenIndex, depth
										{
											position503 := position
											depth++
											{
												position504, tokenIndex504, depth504 := position, tokenIndex, depth
												if buffer[position] != rune('a') {
													goto l505
												}
												position++
												goto l504
											l505:
												position, tokenIndex, depth = position504, tokenIndex504, depth504
												if buffer[position] != rune('A') {
													goto l502
												}
												position++
											}
										l504:
											{
												position506, tokenIndex506, depth506 := position, tokenIndex, depth
												if buffer[position] != rune('s') {
													goto l507
												}
												position++
												goto l506
											l507:
												position, tokenIndex, depth = position506, tokenIndex506, depth506
												if buffer[position] != rune('S') {
													goto l502
												}
												position++
											}
										l506:
											{
												position508, tokenIndex508, depth508 := position, tokenIndex, depth
												if buffer[position] != rune('c') {
													goto l509
												}
												position++
												goto l508
											l509:
												position, tokenIndex, depth = position508, tokenIndex508, depth508
												if buffer[position] != rune('C') {
													goto l502
												}
												position++
											}
										l508:
											if !rules[ruleskip]() {
												goto l502
											}
											depth--
											add(ruleASC, position503)
										}
										goto l501
									l502:
										position, tokenIndex, depth = position501, tokenIndex501, depth501
										{
											position510 := position
											depth++
											{
												position511, tokenIndex511, depth511 := position, tokenIndex, depth
												if buffer[position] != rune('d') {
													goto l512
												}
												position++
												goto l511
											l512:
												position, tokenIndex, depth = position511, tokenIndex511, depth511
												if buffer[position] != rune('D') {
													goto l499
												}
												position++
											}
										l511:
											{
												position513, tokenIndex513, depth513 := position, tokenIndex, depth
												if buffer[position] != rune('e') {
													goto l514
												}
												position++
												goto l513
											l514:
												position, tokenIndex, depth = position513, tokenIndex513, depth513
												if buffer[position] != rune('E') {
													goto l499
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
													goto l499
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
													goto l499
												}
												position++
											}
										l517:
											if !rules[ruleskip]() {
												goto l499
											}
											depth--
											add(ruleDESC, position510)
										}
									}
								l501:
									goto l500
								l499:
									position, tokenIndex, depth = position499, tokenIndex499, depth499
								}
							l500:
								if !rules[rulebrackettedExpression]() {
									goto l498
								}
								goto l497
							l498:
								position, tokenIndex, depth = position497, tokenIndex497, depth497
								if !rules[rulefunctionCall]() {
									goto l519
								}
								goto l497
							l519:
								position, tokenIndex, depth = position497, tokenIndex497, depth497
								if !rules[rulebuiltinCall]() {
									goto l520
								}
								goto l497
							l520:
								position, tokenIndex, depth = position497, tokenIndex497, depth497
								if !rules[rulevar]() {
									goto l482
								}
							}
						l497:
							depth--
							add(ruleorderCondition, position496)
						}
					l494:
						{
							position495, tokenIndex495, depth495 := position, tokenIndex, depth
							{
								position521 := position
								depth++
								{
									position522, tokenIndex522, depth522 := position, tokenIndex, depth
									{
										position524, tokenIndex524, depth524 := position, tokenIndex, depth
										{
											position526, tokenIndex526, depth526 := position, tokenIndex, depth
											{
												position528 := position
												depth++
												{
													position529, tokenIndex529, depth529 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l530
													}
													position++
													goto l529
												l530:
													position, tokenIndex, depth = position529, tokenIndex529, depth529
													if buffer[position] != rune('A') {
														goto l527
													}
													position++
												}
											l529:
												{
													position531, tokenIndex531, depth531 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l532
													}
													position++
													goto l531
												l532:
													position, tokenIndex, depth = position531, tokenIndex531, depth531
													if buffer[position] != rune('S') {
														goto l527
													}
													position++
												}
											l531:
												{
													position533, tokenIndex533, depth533 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l534
													}
													position++
													goto l533
												l534:
													position, tokenIndex, depth = position533, tokenIndex533, depth533
													if buffer[position] != rune('C') {
														goto l527
													}
													position++
												}
											l533:
												if !rules[ruleskip]() {
													goto l527
												}
												depth--
												add(ruleASC, position528)
											}
											goto l526
										l527:
											position, tokenIndex, depth = position526, tokenIndex526, depth526
											{
												position535 := position
												depth++
												{
													position536, tokenIndex536, depth536 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l537
													}
													position++
													goto l536
												l537:
													position, tokenIndex, depth = position536, tokenIndex536, depth536
													if buffer[position] != rune('D') {
														goto l524
													}
													position++
												}
											l536:
												{
													position538, tokenIndex538, depth538 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l539
													}
													position++
													goto l538
												l539:
													position, tokenIndex, depth = position538, tokenIndex538, depth538
													if buffer[position] != rune('E') {
														goto l524
													}
													position++
												}
											l538:
												{
													position540, tokenIndex540, depth540 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l541
													}
													position++
													goto l540
												l541:
													position, tokenIndex, depth = position540, tokenIndex540, depth540
													if buffer[position] != rune('S') {
														goto l524
													}
													position++
												}
											l540:
												{
													position542, tokenIndex542, depth542 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l543
													}
													position++
													goto l542
												l543:
													position, tokenIndex, depth = position542, tokenIndex542, depth542
													if buffer[position] != rune('C') {
														goto l524
													}
													position++
												}
											l542:
												if !rules[ruleskip]() {
													goto l524
												}
												depth--
												add(ruleDESC, position535)
											}
										}
									l526:
										goto l525
									l524:
										position, tokenIndex, depth = position524, tokenIndex524, depth524
									}
								l525:
									if !rules[rulebrackettedExpression]() {
										goto l523
									}
									goto l522
								l523:
									position, tokenIndex, depth = position522, tokenIndex522, depth522
									if !rules[rulefunctionCall]() {
										goto l544
									}
									goto l522
								l544:
									position, tokenIndex, depth = position522, tokenIndex522, depth522
									if !rules[rulebuiltinCall]() {
										goto l545
									}
									goto l522
								l545:
									position, tokenIndex, depth = position522, tokenIndex522, depth522
									if !rules[rulevar]() {
										goto l495
									}
								}
							l522:
								depth--
								add(ruleorderCondition, position521)
							}
							goto l494
						l495:
							position, tokenIndex, depth = position495, tokenIndex495, depth495
						}
						goto l481
					l482:
						position, tokenIndex, depth = position481, tokenIndex481, depth481
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position547 := position
									depth++
									{
										position548, tokenIndex548, depth548 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l549
										}
										position++
										goto l548
									l549:
										position, tokenIndex, depth = position548, tokenIndex548, depth548
										if buffer[position] != rune('H') {
											goto l479
										}
										position++
									}
								l548:
									{
										position550, tokenIndex550, depth550 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l551
										}
										position++
										goto l550
									l551:
										position, tokenIndex, depth = position550, tokenIndex550, depth550
										if buffer[position] != rune('A') {
											goto l479
										}
										position++
									}
								l550:
									{
										position552, tokenIndex552, depth552 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l553
										}
										position++
										goto l552
									l553:
										position, tokenIndex, depth = position552, tokenIndex552, depth552
										if buffer[position] != rune('V') {
											goto l479
										}
										position++
									}
								l552:
									{
										position554, tokenIndex554, depth554 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l555
										}
										position++
										goto l554
									l555:
										position, tokenIndex, depth = position554, tokenIndex554, depth554
										if buffer[position] != rune('I') {
											goto l479
										}
										position++
									}
								l554:
									{
										position556, tokenIndex556, depth556 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l557
										}
										position++
										goto l556
									l557:
										position, tokenIndex, depth = position556, tokenIndex556, depth556
										if buffer[position] != rune('N') {
											goto l479
										}
										position++
									}
								l556:
									{
										position558, tokenIndex558, depth558 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l559
										}
										position++
										goto l558
									l559:
										position, tokenIndex, depth = position558, tokenIndex558, depth558
										if buffer[position] != rune('G') {
											goto l479
										}
										position++
									}
								l558:
									if !rules[ruleskip]() {
										goto l479
									}
									depth--
									add(ruleHAVING, position547)
								}
								if !rules[ruleconstraint]() {
									goto l479
								}
								break
							case 'G', 'g':
								{
									position560 := position
									depth++
									{
										position561, tokenIndex561, depth561 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l562
										}
										position++
										goto l561
									l562:
										position, tokenIndex, depth = position561, tokenIndex561, depth561
										if buffer[position] != rune('G') {
											goto l479
										}
										position++
									}
								l561:
									{
										position563, tokenIndex563, depth563 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l564
										}
										position++
										goto l563
									l564:
										position, tokenIndex, depth = position563, tokenIndex563, depth563
										if buffer[position] != rune('R') {
											goto l479
										}
										position++
									}
								l563:
									{
										position565, tokenIndex565, depth565 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l566
										}
										position++
										goto l565
									l566:
										position, tokenIndex, depth = position565, tokenIndex565, depth565
										if buffer[position] != rune('O') {
											goto l479
										}
										position++
									}
								l565:
									{
										position567, tokenIndex567, depth567 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l568
										}
										position++
										goto l567
									l568:
										position, tokenIndex, depth = position567, tokenIndex567, depth567
										if buffer[position] != rune('U') {
											goto l479
										}
										position++
									}
								l567:
									{
										position569, tokenIndex569, depth569 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l570
										}
										position++
										goto l569
									l570:
										position, tokenIndex, depth = position569, tokenIndex569, depth569
										if buffer[position] != rune('P') {
											goto l479
										}
										position++
									}
								l569:
									if !rules[ruleskip]() {
										goto l479
									}
									depth--
									add(ruleGROUP, position560)
								}
								if !rules[ruleBY]() {
									goto l479
								}
								{
									position573 := position
									depth++
									{
										position574, tokenIndex574, depth574 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l575
										}
										goto l574
									l575:
										position, tokenIndex, depth = position574, tokenIndex574, depth574
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l479
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l479
												}
												if !rules[ruleexpression]() {
													goto l479
												}
												{
													position577, tokenIndex577, depth577 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l577
													}
													if !rules[rulevar]() {
														goto l577
													}
													goto l578
												l577:
													position, tokenIndex, depth = position577, tokenIndex577, depth577
												}
											l578:
												if !rules[ruleRPAREN]() {
													goto l479
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l479
												}
												break
											}
										}

									}
								l574:
									depth--
									add(rulegroupCondition, position573)
								}
							l571:
								{
									position572, tokenIndex572, depth572 := position, tokenIndex, depth
									{
										position579 := position
										depth++
										{
											position580, tokenIndex580, depth580 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l581
											}
											goto l580
										l581:
											position, tokenIndex, depth = position580, tokenIndex580, depth580
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l572
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l572
													}
													if !rules[ruleexpression]() {
														goto l572
													}
													{
														position583, tokenIndex583, depth583 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l583
														}
														if !rules[rulevar]() {
															goto l583
														}
														goto l584
													l583:
														position, tokenIndex, depth = position583, tokenIndex583, depth583
													}
												l584:
													if !rules[ruleRPAREN]() {
														goto l572
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l572
													}
													break
												}
											}

										}
									l580:
										depth--
										add(rulegroupCondition, position579)
									}
									goto l571
								l572:
									position, tokenIndex, depth = position572, tokenIndex572, depth572
								}
								break
							default:
								{
									position585 := position
									depth++
									{
										position586, tokenIndex586, depth586 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l587
										}
										{
											position588, tokenIndex588, depth588 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l588
											}
											goto l589
										l588:
											position, tokenIndex, depth = position588, tokenIndex588, depth588
										}
									l589:
										goto l586
									l587:
										position, tokenIndex, depth = position586, tokenIndex586, depth586
										if !rules[ruleoffset]() {
											goto l479
										}
										{
											position590, tokenIndex590, depth590 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l590
											}
											goto l591
										l590:
											position, tokenIndex, depth = position590, tokenIndex590, depth590
										}
									l591:
									}
								l586:
									depth--
									add(rulelimitOffsetClauses, position585)
								}
								break
							}
						}

					}
				l481:
					goto l480
				l479:
					position, tokenIndex, depth = position479, tokenIndex479, depth479
				}
			l480:
				depth--
				add(rulesolutionModifier, position478)
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
			position595, tokenIndex595, depth595 := position, tokenIndex, depth
			{
				position596 := position
				depth++
				{
					position597 := position
					depth++
					{
						position598, tokenIndex598, depth598 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l599
						}
						position++
						goto l598
					l599:
						position, tokenIndex, depth = position598, tokenIndex598, depth598
						if buffer[position] != rune('L') {
							goto l595
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
							goto l595
						}
						position++
					}
				l600:
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l603
						}
						position++
						goto l602
					l603:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
						if buffer[position] != rune('M') {
							goto l595
						}
						position++
					}
				l602:
					{
						position604, tokenIndex604, depth604 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l605
						}
						position++
						goto l604
					l605:
						position, tokenIndex, depth = position604, tokenIndex604, depth604
						if buffer[position] != rune('I') {
							goto l595
						}
						position++
					}
				l604:
					{
						position606, tokenIndex606, depth606 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l607
						}
						position++
						goto l606
					l607:
						position, tokenIndex, depth = position606, tokenIndex606, depth606
						if buffer[position] != rune('T') {
							goto l595
						}
						position++
					}
				l606:
					if !rules[ruleskip]() {
						goto l595
					}
					depth--
					add(ruleLIMIT, position597)
				}
				if !rules[ruleINTEGER]() {
					goto l595
				}
				depth--
				add(rulelimit, position596)
			}
			return true
		l595:
			position, tokenIndex, depth = position595, tokenIndex595, depth595
			return false
		},
		/* 51 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position608, tokenIndex608, depth608 := position, tokenIndex, depth
			{
				position609 := position
				depth++
				{
					position610 := position
					depth++
					{
						position611, tokenIndex611, depth611 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l612
						}
						position++
						goto l611
					l612:
						position, tokenIndex, depth = position611, tokenIndex611, depth611
						if buffer[position] != rune('O') {
							goto l608
						}
						position++
					}
				l611:
					{
						position613, tokenIndex613, depth613 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l614
						}
						position++
						goto l613
					l614:
						position, tokenIndex, depth = position613, tokenIndex613, depth613
						if buffer[position] != rune('F') {
							goto l608
						}
						position++
					}
				l613:
					{
						position615, tokenIndex615, depth615 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l616
						}
						position++
						goto l615
					l616:
						position, tokenIndex, depth = position615, tokenIndex615, depth615
						if buffer[position] != rune('F') {
							goto l608
						}
						position++
					}
				l615:
					{
						position617, tokenIndex617, depth617 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l618
						}
						position++
						goto l617
					l618:
						position, tokenIndex, depth = position617, tokenIndex617, depth617
						if buffer[position] != rune('S') {
							goto l608
						}
						position++
					}
				l617:
					{
						position619, tokenIndex619, depth619 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l620
						}
						position++
						goto l619
					l620:
						position, tokenIndex, depth = position619, tokenIndex619, depth619
						if buffer[position] != rune('E') {
							goto l608
						}
						position++
					}
				l619:
					{
						position621, tokenIndex621, depth621 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l622
						}
						position++
						goto l621
					l622:
						position, tokenIndex, depth = position621, tokenIndex621, depth621
						if buffer[position] != rune('T') {
							goto l608
						}
						position++
					}
				l621:
					if !rules[ruleskip]() {
						goto l608
					}
					depth--
					add(ruleOFFSET, position610)
				}
				if !rules[ruleINTEGER]() {
					goto l608
				}
				depth--
				add(ruleoffset, position609)
			}
			return true
		l608:
			position, tokenIndex, depth = position608, tokenIndex608, depth608
			return false
		},
		/* 52 expression <- <conditionalOrExpression> */
		func() bool {
			position623, tokenIndex623, depth623 := position, tokenIndex, depth
			{
				position624 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l623
				}
				depth--
				add(ruleexpression, position624)
			}
			return true
		l623:
			position, tokenIndex, depth = position623, tokenIndex623, depth623
			return false
		},
		/* 53 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position625, tokenIndex625, depth625 := position, tokenIndex, depth
			{
				position626 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l625
				}
				{
					position627, tokenIndex627, depth627 := position, tokenIndex, depth
					{
						position629 := position
						depth++
						if buffer[position] != rune('|') {
							goto l627
						}
						position++
						if buffer[position] != rune('|') {
							goto l627
						}
						position++
						if !rules[ruleskip]() {
							goto l627
						}
						depth--
						add(ruleOR, position629)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l627
					}
					goto l628
				l627:
					position, tokenIndex, depth = position627, tokenIndex627, depth627
				}
			l628:
				depth--
				add(ruleconditionalOrExpression, position626)
			}
			return true
		l625:
			position, tokenIndex, depth = position625, tokenIndex625, depth625
			return false
		},
		/* 54 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position630, tokenIndex630, depth630 := position, tokenIndex, depth
			{
				position631 := position
				depth++
				{
					position632 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l630
					}
					{
						position633, tokenIndex633, depth633 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position636 := position
									depth++
									{
										position637 := position
										depth++
										{
											position638, tokenIndex638, depth638 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l639
											}
											position++
											goto l638
										l639:
											position, tokenIndex, depth = position638, tokenIndex638, depth638
											if buffer[position] != rune('N') {
												goto l633
											}
											position++
										}
									l638:
										{
											position640, tokenIndex640, depth640 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l641
											}
											position++
											goto l640
										l641:
											position, tokenIndex, depth = position640, tokenIndex640, depth640
											if buffer[position] != rune('O') {
												goto l633
											}
											position++
										}
									l640:
										{
											position642, tokenIndex642, depth642 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l643
											}
											position++
											goto l642
										l643:
											position, tokenIndex, depth = position642, tokenIndex642, depth642
											if buffer[position] != rune('T') {
												goto l633
											}
											position++
										}
									l642:
										if buffer[position] != rune(' ') {
											goto l633
										}
										position++
										{
											position644, tokenIndex644, depth644 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l645
											}
											position++
											goto l644
										l645:
											position, tokenIndex, depth = position644, tokenIndex644, depth644
											if buffer[position] != rune('I') {
												goto l633
											}
											position++
										}
									l644:
										{
											position646, tokenIndex646, depth646 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l647
											}
											position++
											goto l646
										l647:
											position, tokenIndex, depth = position646, tokenIndex646, depth646
											if buffer[position] != rune('N') {
												goto l633
											}
											position++
										}
									l646:
										if !rules[ruleskip]() {
											goto l633
										}
										depth--
										add(ruleNOTIN, position637)
									}
									if !rules[ruleargList]() {
										goto l633
									}
									depth--
									add(rulenotin, position636)
								}
								break
							case 'I', 'i':
								{
									position648 := position
									depth++
									{
										position649 := position
										depth++
										{
											position650, tokenIndex650, depth650 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l651
											}
											position++
											goto l650
										l651:
											position, tokenIndex, depth = position650, tokenIndex650, depth650
											if buffer[position] != rune('I') {
												goto l633
											}
											position++
										}
									l650:
										{
											position652, tokenIndex652, depth652 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l653
											}
											position++
											goto l652
										l653:
											position, tokenIndex, depth = position652, tokenIndex652, depth652
											if buffer[position] != rune('N') {
												goto l633
											}
											position++
										}
									l652:
										if !rules[ruleskip]() {
											goto l633
										}
										depth--
										add(ruleIN, position649)
									}
									if !rules[ruleargList]() {
										goto l633
									}
									depth--
									add(rulein, position648)
								}
								break
							default:
								{
									position654, tokenIndex654, depth654 := position, tokenIndex, depth
									{
										position656 := position
										depth++
										if buffer[position] != rune('<') {
											goto l655
										}
										position++
										if !rules[ruleskip]() {
											goto l655
										}
										depth--
										add(ruleLT, position656)
									}
									goto l654
								l655:
									position, tokenIndex, depth = position654, tokenIndex654, depth654
									{
										position658 := position
										depth++
										if buffer[position] != rune('>') {
											goto l657
										}
										position++
										if buffer[position] != rune('=') {
											goto l657
										}
										position++
										if !rules[ruleskip]() {
											goto l657
										}
										depth--
										add(ruleGE, position658)
									}
									goto l654
								l657:
									position, tokenIndex, depth = position654, tokenIndex654, depth654
									{
										switch buffer[position] {
										case '>':
											{
												position660 := position
												depth++
												if buffer[position] != rune('>') {
													goto l633
												}
												position++
												if !rules[ruleskip]() {
													goto l633
												}
												depth--
												add(ruleGT, position660)
											}
											break
										case '<':
											{
												position661 := position
												depth++
												if buffer[position] != rune('<') {
													goto l633
												}
												position++
												if buffer[position] != rune('=') {
													goto l633
												}
												position++
												if !rules[ruleskip]() {
													goto l633
												}
												depth--
												add(ruleLE, position661)
											}
											break
										case '!':
											{
												position662 := position
												depth++
												if buffer[position] != rune('!') {
													goto l633
												}
												position++
												if buffer[position] != rune('=') {
													goto l633
												}
												position++
												if !rules[ruleskip]() {
													goto l633
												}
												depth--
												add(ruleNE, position662)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l633
											}
											break
										}
									}

								}
							l654:
								if !rules[rulenumericExpression]() {
									goto l633
								}
								break
							}
						}

						goto l634
					l633:
						position, tokenIndex, depth = position633, tokenIndex633, depth633
					}
				l634:
					depth--
					add(rulevalueLogical, position632)
				}
				{
					position663, tokenIndex663, depth663 := position, tokenIndex, depth
					{
						position665 := position
						depth++
						if buffer[position] != rune('&') {
							goto l663
						}
						position++
						if buffer[position] != rune('&') {
							goto l663
						}
						position++
						if !rules[ruleskip]() {
							goto l663
						}
						depth--
						add(ruleAND, position665)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l663
					}
					goto l664
				l663:
					position, tokenIndex, depth = position663, tokenIndex663, depth663
				}
			l664:
				depth--
				add(ruleconditionalAndExpression, position631)
			}
			return true
		l630:
			position, tokenIndex, depth = position630, tokenIndex630, depth630
			return false
		},
		/* 55 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 56 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position667, tokenIndex667, depth667 := position, tokenIndex, depth
			{
				position668 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l667
				}
			l669:
				{
					position670, tokenIndex670, depth670 := position, tokenIndex, depth
					{
						position671, tokenIndex671, depth671 := position, tokenIndex, depth
						{
							position673, tokenIndex673, depth673 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l674
							}
							goto l673
						l674:
							position, tokenIndex, depth = position673, tokenIndex673, depth673
							if !rules[ruleMINUS]() {
								goto l672
							}
						}
					l673:
						if !rules[rulemultiplicativeExpression]() {
							goto l672
						}
						goto l671
					l672:
						position, tokenIndex, depth = position671, tokenIndex671, depth671
						{
							position675 := position
							depth++
							{
								position676, tokenIndex676, depth676 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l677
								}
								position++
								goto l676
							l677:
								position, tokenIndex, depth = position676, tokenIndex676, depth676
								if buffer[position] != rune('-') {
									goto l670
								}
								position++
							}
						l676:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l670
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
							{
								position680, tokenIndex680, depth680 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l680
								}
								position++
							l682:
								{
									position683, tokenIndex683, depth683 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l683
									}
									position++
									goto l682
								l683:
									position, tokenIndex, depth = position683, tokenIndex683, depth683
								}
								goto l681
							l680:
								position, tokenIndex, depth = position680, tokenIndex680, depth680
							}
						l681:
							if !rules[ruleskip]() {
								goto l670
							}
							depth--
							add(rulesignedNumericLiteral, position675)
						}
					}
				l671:
					goto l669
				l670:
					position, tokenIndex, depth = position670, tokenIndex670, depth670
				}
				depth--
				add(rulenumericExpression, position668)
			}
			return true
		l667:
			position, tokenIndex, depth = position667, tokenIndex667, depth667
			return false
		},
		/* 57 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position684, tokenIndex684, depth684 := position, tokenIndex, depth
			{
				position685 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l684
				}
			l686:
				{
					position687, tokenIndex687, depth687 := position, tokenIndex, depth
					{
						position688, tokenIndex688, depth688 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l689
						}
						goto l688
					l689:
						position, tokenIndex, depth = position688, tokenIndex688, depth688
						if !rules[ruleSLASH]() {
							goto l687
						}
					}
				l688:
					if !rules[ruleunaryExpression]() {
						goto l687
					}
					goto l686
				l687:
					position, tokenIndex, depth = position687, tokenIndex687, depth687
				}
				depth--
				add(rulemultiplicativeExpression, position685)
			}
			return true
		l684:
			position, tokenIndex, depth = position684, tokenIndex684, depth684
			return false
		},
		/* 58 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position690, tokenIndex690, depth690 := position, tokenIndex, depth
			{
				position691 := position
				depth++
				{
					position692, tokenIndex692, depth692 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l692
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l692
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l692
							}
							break
						}
					}

					goto l693
				l692:
					position, tokenIndex, depth = position692, tokenIndex692, depth692
				}
			l693:
				{
					position695 := position
					depth++
					{
						position696, tokenIndex696, depth696 := position, tokenIndex, depth
						if !rules[rulebuiltinCall]() {
							goto l697
						}
						goto l696
					l697:
						position, tokenIndex, depth = position696, tokenIndex696, depth696
						if !rules[rulefunctionCall]() {
							goto l698
						}
						goto l696
					l698:
						position, tokenIndex, depth = position696, tokenIndex696, depth696
						if !rules[ruleiriref]() {
							goto l699
						}
						goto l696
					l699:
						position, tokenIndex, depth = position696, tokenIndex696, depth696
						{
							switch buffer[position] {
							case 'A', 'C', 'G', 'M', 'S', 'a', 'c', 'g', 'm', 's':
								{
									position701 := position
									depth++
									{
										switch buffer[position] {
										case 'G', 'g':
											{
												position703 := position
												depth++
												{
													position704 := position
													depth++
													{
														position705, tokenIndex705, depth705 := position, tokenIndex, depth
														if buffer[position] != rune('g') {
															goto l706
														}
														position++
														goto l705
													l706:
														position, tokenIndex, depth = position705, tokenIndex705, depth705
														if buffer[position] != rune('G') {
															goto l690
														}
														position++
													}
												l705:
													{
														position707, tokenIndex707, depth707 := position, tokenIndex, depth
														if buffer[position] != rune('r') {
															goto l708
														}
														position++
														goto l707
													l708:
														position, tokenIndex, depth = position707, tokenIndex707, depth707
														if buffer[position] != rune('R') {
															goto l690
														}
														position++
													}
												l707:
													{
														position709, tokenIndex709, depth709 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l710
														}
														position++
														goto l709
													l710:
														position, tokenIndex, depth = position709, tokenIndex709, depth709
														if buffer[position] != rune('O') {
															goto l690
														}
														position++
													}
												l709:
													{
														position711, tokenIndex711, depth711 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l712
														}
														position++
														goto l711
													l712:
														position, tokenIndex, depth = position711, tokenIndex711, depth711
														if buffer[position] != rune('U') {
															goto l690
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
															goto l690
														}
														position++
													}
												l713:
													if buffer[position] != rune('_') {
														goto l690
													}
													position++
													{
														position715, tokenIndex715, depth715 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l716
														}
														position++
														goto l715
													l716:
														position, tokenIndex, depth = position715, tokenIndex715, depth715
														if buffer[position] != rune('C') {
															goto l690
														}
														position++
													}
												l715:
													{
														position717, tokenIndex717, depth717 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l718
														}
														position++
														goto l717
													l718:
														position, tokenIndex, depth = position717, tokenIndex717, depth717
														if buffer[position] != rune('O') {
															goto l690
														}
														position++
													}
												l717:
													{
														position719, tokenIndex719, depth719 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l720
														}
														position++
														goto l719
													l720:
														position, tokenIndex, depth = position719, tokenIndex719, depth719
														if buffer[position] != rune('N') {
															goto l690
														}
														position++
													}
												l719:
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
															goto l690
														}
														position++
													}
												l721:
													{
														position723, tokenIndex723, depth723 := position, tokenIndex, depth
														if buffer[position] != rune('a') {
															goto l724
														}
														position++
														goto l723
													l724:
														position, tokenIndex, depth = position723, tokenIndex723, depth723
														if buffer[position] != rune('A') {
															goto l690
														}
														position++
													}
												l723:
													{
														position725, tokenIndex725, depth725 := position, tokenIndex, depth
														if buffer[position] != rune('t') {
															goto l726
														}
														position++
														goto l725
													l726:
														position, tokenIndex, depth = position725, tokenIndex725, depth725
														if buffer[position] != rune('T') {
															goto l690
														}
														position++
													}
												l725:
													if !rules[ruleskip]() {
														goto l690
													}
													depth--
													add(ruleGROUPCONCAT, position704)
												}
												if !rules[ruleLPAREN]() {
													goto l690
												}
												{
													position727, tokenIndex727, depth727 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l727
													}
													goto l728
												l727:
													position, tokenIndex, depth = position727, tokenIndex727, depth727
												}
											l728:
												if !rules[ruleexpression]() {
													goto l690
												}
												{
													position729, tokenIndex729, depth729 := position, tokenIndex, depth
													if !rules[ruleSEMICOLON]() {
														goto l729
													}
													{
														position731 := position
														depth++
														{
															position732, tokenIndex732, depth732 := position, tokenIndex, depth
															if buffer[position] != rune('s') {
																goto l733
															}
															position++
															goto l732
														l733:
															position, tokenIndex, depth = position732, tokenIndex732, depth732
															if buffer[position] != rune('S') {
																goto l729
															}
															position++
														}
													l732:
														{
															position734, tokenIndex734, depth734 := position, tokenIndex, depth
															if buffer[position] != rune('e') {
																goto l735
															}
															position++
															goto l734
														l735:
															position, tokenIndex, depth = position734, tokenIndex734, depth734
															if buffer[position] != rune('E') {
																goto l729
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
																goto l729
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
																goto l729
															}
															position++
														}
													l738:
														{
															position740, tokenIndex740, depth740 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l741
															}
															position++
															goto l740
														l741:
															position, tokenIndex, depth = position740, tokenIndex740, depth740
															if buffer[position] != rune('R') {
																goto l729
															}
															position++
														}
													l740:
														{
															position742, tokenIndex742, depth742 := position, tokenIndex, depth
															if buffer[position] != rune('a') {
																goto l743
															}
															position++
															goto l742
														l743:
															position, tokenIndex, depth = position742, tokenIndex742, depth742
															if buffer[position] != rune('A') {
																goto l729
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
																goto l729
															}
															position++
														}
													l744:
														{
															position746, tokenIndex746, depth746 := position, tokenIndex, depth
															if buffer[position] != rune('o') {
																goto l747
															}
															position++
															goto l746
														l747:
															position, tokenIndex, depth = position746, tokenIndex746, depth746
															if buffer[position] != rune('O') {
																goto l729
															}
															position++
														}
													l746:
														{
															position748, tokenIndex748, depth748 := position, tokenIndex, depth
															if buffer[position] != rune('r') {
																goto l749
															}
															position++
															goto l748
														l749:
															position, tokenIndex, depth = position748, tokenIndex748, depth748
															if buffer[position] != rune('R') {
																goto l729
															}
															position++
														}
													l748:
														if !rules[ruleskip]() {
															goto l729
														}
														depth--
														add(ruleSEPARATOR, position731)
													}
													if !rules[ruleEQ]() {
														goto l729
													}
													if !rules[rulestring]() {
														goto l729
													}
													goto l730
												l729:
													position, tokenIndex, depth = position729, tokenIndex729, depth729
												}
											l730:
												if !rules[ruleRPAREN]() {
													goto l690
												}
												depth--
												add(rulegroupConcat, position703)
											}
											break
										case 'C', 'c':
											{
												position750 := position
												depth++
												{
													position751 := position
													depth++
													{
														position752, tokenIndex752, depth752 := position, tokenIndex, depth
														if buffer[position] != rune('c') {
															goto l753
														}
														position++
														goto l752
													l753:
														position, tokenIndex, depth = position752, tokenIndex752, depth752
														if buffer[position] != rune('C') {
															goto l690
														}
														position++
													}
												l752:
													{
														position754, tokenIndex754, depth754 := position, tokenIndex, depth
														if buffer[position] != rune('o') {
															goto l755
														}
														position++
														goto l754
													l755:
														position, tokenIndex, depth = position754, tokenIndex754, depth754
														if buffer[position] != rune('O') {
															goto l690
														}
														position++
													}
												l754:
													{
														position756, tokenIndex756, depth756 := position, tokenIndex, depth
														if buffer[position] != rune('u') {
															goto l757
														}
														position++
														goto l756
													l757:
														position, tokenIndex, depth = position756, tokenIndex756, depth756
														if buffer[position] != rune('U') {
															goto l690
														}
														position++
													}
												l756:
													{
														position758, tokenIndex758, depth758 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l759
														}
														position++
														goto l758
													l759:
														position, tokenIndex, depth = position758, tokenIndex758, depth758
														if buffer[position] != rune('N') {
															goto l690
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
															goto l690
														}
														position++
													}
												l760:
													if !rules[ruleskip]() {
														goto l690
													}
													depth--
													add(ruleCOUNT, position751)
												}
												if !rules[ruleLPAREN]() {
													goto l690
												}
												{
													position762, tokenIndex762, depth762 := position, tokenIndex, depth
													if !rules[ruleDISTINCT]() {
														goto l762
													}
													goto l763
												l762:
													position, tokenIndex, depth = position762, tokenIndex762, depth762
												}
											l763:
												{
													position764, tokenIndex764, depth764 := position, tokenIndex, depth
													if !rules[ruleSTAR]() {
														goto l765
													}
													goto l764
												l765:
													position, tokenIndex, depth = position764, tokenIndex764, depth764
													if !rules[ruleexpression]() {
														goto l690
													}
												}
											l764:
												if !rules[ruleRPAREN]() {
													goto l690
												}
												depth--
												add(rulecount, position750)
											}
											break
										default:
											{
												position766, tokenIndex766, depth766 := position, tokenIndex, depth
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
														if buffer[position] != rune('u') {
															goto l772
														}
														position++
														goto l771
													l772:
														position, tokenIndex, depth = position771, tokenIndex771, depth771
														if buffer[position] != rune('U') {
															goto l767
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
															goto l767
														}
														position++
													}
												l773:
													if !rules[ruleskip]() {
														goto l767
													}
													depth--
													add(ruleSUM, position768)
												}
												goto l766
											l767:
												position, tokenIndex, depth = position766, tokenIndex766, depth766
												{
													position776 := position
													depth++
													{
														position777, tokenIndex777, depth777 := position, tokenIndex, depth
														if buffer[position] != rune('m') {
															goto l778
														}
														position++
														goto l777
													l778:
														position, tokenIndex, depth = position777, tokenIndex777, depth777
														if buffer[position] != rune('M') {
															goto l775
														}
														position++
													}
												l777:
													{
														position779, tokenIndex779, depth779 := position, tokenIndex, depth
														if buffer[position] != rune('i') {
															goto l780
														}
														position++
														goto l779
													l780:
														position, tokenIndex, depth = position779, tokenIndex779, depth779
														if buffer[position] != rune('I') {
															goto l775
														}
														position++
													}
												l779:
													{
														position781, tokenIndex781, depth781 := position, tokenIndex, depth
														if buffer[position] != rune('n') {
															goto l782
														}
														position++
														goto l781
													l782:
														position, tokenIndex, depth = position781, tokenIndex781, depth781
														if buffer[position] != rune('N') {
															goto l775
														}
														position++
													}
												l781:
													if !rules[ruleskip]() {
														goto l775
													}
													depth--
													add(ruleMIN, position776)
												}
												goto l766
											l775:
												position, tokenIndex, depth = position766, tokenIndex766, depth766
												{
													switch buffer[position] {
													case 'S', 's':
														{
															position784 := position
															depth++
															{
																position785, tokenIndex785, depth785 := position, tokenIndex, depth
																if buffer[position] != rune('s') {
																	goto l786
																}
																position++
																goto l785
															l786:
																position, tokenIndex, depth = position785, tokenIndex785, depth785
																if buffer[position] != rune('S') {
																	goto l690
																}
																position++
															}
														l785:
															{
																position787, tokenIndex787, depth787 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l788
																}
																position++
																goto l787
															l788:
																position, tokenIndex, depth = position787, tokenIndex787, depth787
																if buffer[position] != rune('A') {
																	goto l690
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
																	goto l690
																}
																position++
															}
														l789:
															{
																position791, tokenIndex791, depth791 := position, tokenIndex, depth
																if buffer[position] != rune('p') {
																	goto l792
																}
																position++
																goto l791
															l792:
																position, tokenIndex, depth = position791, tokenIndex791, depth791
																if buffer[position] != rune('P') {
																	goto l690
																}
																position++
															}
														l791:
															{
																position793, tokenIndex793, depth793 := position, tokenIndex, depth
																if buffer[position] != rune('l') {
																	goto l794
																}
																position++
																goto l793
															l794:
																position, tokenIndex, depth = position793, tokenIndex793, depth793
																if buffer[position] != rune('L') {
																	goto l690
																}
																position++
															}
														l793:
															{
																position795, tokenIndex795, depth795 := position, tokenIndex, depth
																if buffer[position] != rune('e') {
																	goto l796
																}
																position++
																goto l795
															l796:
																position, tokenIndex, depth = position795, tokenIndex795, depth795
																if buffer[position] != rune('E') {
																	goto l690
																}
																position++
															}
														l795:
															if !rules[ruleskip]() {
																goto l690
															}
															depth--
															add(ruleSAMPLE, position784)
														}
														break
													case 'A', 'a':
														{
															position797 := position
															depth++
															{
																position798, tokenIndex798, depth798 := position, tokenIndex, depth
																if buffer[position] != rune('a') {
																	goto l799
																}
																position++
																goto l798
															l799:
																position, tokenIndex, depth = position798, tokenIndex798, depth798
																if buffer[position] != rune('A') {
																	goto l690
																}
																position++
															}
														l798:
															{
																position800, tokenIndex800, depth800 := position, tokenIndex, depth
																if buffer[position] != rune('v') {
																	goto l801
																}
																position++
																goto l800
															l801:
																position, tokenIndex, depth = position800, tokenIndex800, depth800
																if buffer[position] != rune('V') {
																	goto l690
																}
																position++
															}
														l800:
															{
																position802, tokenIndex802, depth802 := position, tokenIndex, depth
																if buffer[position] != rune('g') {
																	goto l803
																}
																position++
																goto l802
															l803:
																position, tokenIndex, depth = position802, tokenIndex802, depth802
																if buffer[position] != rune('G') {
																	goto l690
																}
																position++
															}
														l802:
															if !rules[ruleskip]() {
																goto l690
															}
															depth--
															add(ruleAVG, position797)
														}
														break
													default:
														{
															position804 := position
															depth++
															{
																position805, tokenIndex805, depth805 := position, tokenIndex, depth
																if buffer[position] != rune('m') {
																	goto l806
																}
																position++
																goto l805
															l806:
																position, tokenIndex, depth = position805, tokenIndex805, depth805
																if buffer[position] != rune('M') {
																	goto l690
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
																	goto l690
																}
																position++
															}
														l807:
															{
																position809, tokenIndex809, depth809 := position, tokenIndex, depth
																if buffer[position] != rune('x') {
																	goto l810
																}
																position++
																goto l809
															l810:
																position, tokenIndex, depth = position809, tokenIndex809, depth809
																if buffer[position] != rune('X') {
																	goto l690
																}
																position++
															}
														l809:
															if !rules[ruleskip]() {
																goto l690
															}
															depth--
															add(ruleMAX, position804)
														}
														break
													}
												}

											}
										l766:
											if !rules[ruleLPAREN]() {
												goto l690
											}
											{
												position811, tokenIndex811, depth811 := position, tokenIndex, depth
												if !rules[ruleDISTINCT]() {
													goto l811
												}
												goto l812
											l811:
												position, tokenIndex, depth = position811, tokenIndex811, depth811
											}
										l812:
											if !rules[ruleexpression]() {
												goto l690
											}
											if !rules[ruleRPAREN]() {
												goto l690
											}
											break
										}
									}

									depth--
									add(ruleaggregate, position701)
								}
								break
							case '$', '?':
								if !rules[rulevar]() {
									goto l690
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l690
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l690
								}
								break
							case '(':
								if !rules[rulebrackettedExpression]() {
									goto l690
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l690
								}
								break
							}
						}

					}
				l696:
					depth--
					add(ruleprimaryExpression, position695)
				}
				depth--
				add(ruleunaryExpression, position691)
			}
			return true
		l690:
			position, tokenIndex, depth = position690, tokenIndex690, depth690
			return false
		},
		/* 59 primaryExpression <- <(builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('(') brackettedExpression) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 60 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position814, tokenIndex814, depth814 := position, tokenIndex, depth
			{
				position815 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l814
				}
				if !rules[ruleexpression]() {
					goto l814
				}
				if !rules[ruleRPAREN]() {
					goto l814
				}
				depth--
				add(rulebrackettedExpression, position815)
			}
			return true
		l814:
			position, tokenIndex, depth = position814, tokenIndex814, depth814
			return false
		},
		/* 61 functionCall <- <(iriref argList)> */
		func() bool {
			position816, tokenIndex816, depth816 := position, tokenIndex, depth
			{
				position817 := position
				depth++
				if !rules[ruleiriref]() {
					goto l816
				}
				if !rules[ruleargList]() {
					goto l816
				}
				depth--
				add(rulefunctionCall, position817)
			}
			return true
		l816:
			position, tokenIndex, depth = position816, tokenIndex816, depth816
			return false
		},
		/* 62 in <- <(IN argList)> */
		nil,
		/* 63 notin <- <(NOTIN argList)> */
		nil,
		/* 64 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position820, tokenIndex820, depth820 := position, tokenIndex, depth
			{
				position821 := position
				depth++
				{
					position822, tokenIndex822, depth822 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l823
					}
					goto l822
				l823:
					position, tokenIndex, depth = position822, tokenIndex822, depth822
					if !rules[ruleLPAREN]() {
						goto l820
					}
					if !rules[ruleexpression]() {
						goto l820
					}
				l824:
					{
						position825, tokenIndex825, depth825 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l825
						}
						if !rules[ruleexpression]() {
							goto l825
						}
						goto l824
					l825:
						position, tokenIndex, depth = position825, tokenIndex825, depth825
					}
					if !rules[ruleRPAREN]() {
						goto l820
					}
				}
			l822:
				depth--
				add(ruleargList, position821)
			}
			return true
		l820:
			position, tokenIndex, depth = position820, tokenIndex820, depth820
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
			position829, tokenIndex829, depth829 := position, tokenIndex, depth
			{
				position830 := position
				depth++
				{
					position831, tokenIndex831, depth831 := position, tokenIndex, depth
					{
						position833, tokenIndex833, depth833 := position, tokenIndex, depth
						{
							position835 := position
							depth++
							{
								position836, tokenIndex836, depth836 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l837
								}
								position++
								goto l836
							l837:
								position, tokenIndex, depth = position836, tokenIndex836, depth836
								if buffer[position] != rune('S') {
									goto l834
								}
								position++
							}
						l836:
							{
								position838, tokenIndex838, depth838 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l839
								}
								position++
								goto l838
							l839:
								position, tokenIndex, depth = position838, tokenIndex838, depth838
								if buffer[position] != rune('T') {
									goto l834
								}
								position++
							}
						l838:
							{
								position840, tokenIndex840, depth840 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l841
								}
								position++
								goto l840
							l841:
								position, tokenIndex, depth = position840, tokenIndex840, depth840
								if buffer[position] != rune('R') {
									goto l834
								}
								position++
							}
						l840:
							if !rules[ruleskip]() {
								goto l834
							}
							depth--
							add(ruleSTR, position835)
						}
						goto l833
					l834:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position843 := position
							depth++
							{
								position844, tokenIndex844, depth844 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l845
								}
								position++
								goto l844
							l845:
								position, tokenIndex, depth = position844, tokenIndex844, depth844
								if buffer[position] != rune('L') {
									goto l842
								}
								position++
							}
						l844:
							{
								position846, tokenIndex846, depth846 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l847
								}
								position++
								goto l846
							l847:
								position, tokenIndex, depth = position846, tokenIndex846, depth846
								if buffer[position] != rune('A') {
									goto l842
								}
								position++
							}
						l846:
							{
								position848, tokenIndex848, depth848 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l849
								}
								position++
								goto l848
							l849:
								position, tokenIndex, depth = position848, tokenIndex848, depth848
								if buffer[position] != rune('N') {
									goto l842
								}
								position++
							}
						l848:
							{
								position850, tokenIndex850, depth850 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l851
								}
								position++
								goto l850
							l851:
								position, tokenIndex, depth = position850, tokenIndex850, depth850
								if buffer[position] != rune('G') {
									goto l842
								}
								position++
							}
						l850:
							if !rules[ruleskip]() {
								goto l842
							}
							depth--
							add(ruleLANG, position843)
						}
						goto l833
					l842:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position853 := position
							depth++
							{
								position854, tokenIndex854, depth854 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l855
								}
								position++
								goto l854
							l855:
								position, tokenIndex, depth = position854, tokenIndex854, depth854
								if buffer[position] != rune('D') {
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
								if buffer[position] != rune('t') {
									goto l859
								}
								position++
								goto l858
							l859:
								position, tokenIndex, depth = position858, tokenIndex858, depth858
								if buffer[position] != rune('T') {
									goto l852
								}
								position++
							}
						l858:
							{
								position860, tokenIndex860, depth860 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l861
								}
								position++
								goto l860
							l861:
								position, tokenIndex, depth = position860, tokenIndex860, depth860
								if buffer[position] != rune('A') {
									goto l852
								}
								position++
							}
						l860:
							{
								position862, tokenIndex862, depth862 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l863
								}
								position++
								goto l862
							l863:
								position, tokenIndex, depth = position862, tokenIndex862, depth862
								if buffer[position] != rune('T') {
									goto l852
								}
								position++
							}
						l862:
							{
								position864, tokenIndex864, depth864 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l865
								}
								position++
								goto l864
							l865:
								position, tokenIndex, depth = position864, tokenIndex864, depth864
								if buffer[position] != rune('Y') {
									goto l852
								}
								position++
							}
						l864:
							{
								position866, tokenIndex866, depth866 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l867
								}
								position++
								goto l866
							l867:
								position, tokenIndex, depth = position866, tokenIndex866, depth866
								if buffer[position] != rune('P') {
									goto l852
								}
								position++
							}
						l866:
							{
								position868, tokenIndex868, depth868 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l869
								}
								position++
								goto l868
							l869:
								position, tokenIndex, depth = position868, tokenIndex868, depth868
								if buffer[position] != rune('E') {
									goto l852
								}
								position++
							}
						l868:
							if !rules[ruleskip]() {
								goto l852
							}
							depth--
							add(ruleDATATYPE, position853)
						}
						goto l833
					l852:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position871 := position
							depth++
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
									goto l870
								}
								position++
							}
						l872:
							{
								position874, tokenIndex874, depth874 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l875
								}
								position++
								goto l874
							l875:
								position, tokenIndex, depth = position874, tokenIndex874, depth874
								if buffer[position] != rune('R') {
									goto l870
								}
								position++
							}
						l874:
							{
								position876, tokenIndex876, depth876 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l877
								}
								position++
								goto l876
							l877:
								position, tokenIndex, depth = position876, tokenIndex876, depth876
								if buffer[position] != rune('I') {
									goto l870
								}
								position++
							}
						l876:
							if !rules[ruleskip]() {
								goto l870
							}
							depth--
							add(ruleIRI, position871)
						}
						goto l833
					l870:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position879 := position
							depth++
							{
								position880, tokenIndex880, depth880 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l881
								}
								position++
								goto l880
							l881:
								position, tokenIndex, depth = position880, tokenIndex880, depth880
								if buffer[position] != rune('U') {
									goto l878
								}
								position++
							}
						l880:
							{
								position882, tokenIndex882, depth882 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l883
								}
								position++
								goto l882
							l883:
								position, tokenIndex, depth = position882, tokenIndex882, depth882
								if buffer[position] != rune('R') {
									goto l878
								}
								position++
							}
						l882:
							{
								position884, tokenIndex884, depth884 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l885
								}
								position++
								goto l884
							l885:
								position, tokenIndex, depth = position884, tokenIndex884, depth884
								if buffer[position] != rune('I') {
									goto l878
								}
								position++
							}
						l884:
							if !rules[ruleskip]() {
								goto l878
							}
							depth--
							add(ruleURI, position879)
						}
						goto l833
					l878:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position887 := position
							depth++
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
									goto l886
								}
								position++
							}
						l888:
							{
								position890, tokenIndex890, depth890 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l891
								}
								position++
								goto l890
							l891:
								position, tokenIndex, depth = position890, tokenIndex890, depth890
								if buffer[position] != rune('T') {
									goto l886
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
									goto l886
								}
								position++
							}
						l892:
							{
								position894, tokenIndex894, depth894 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l895
								}
								position++
								goto l894
							l895:
								position, tokenIndex, depth = position894, tokenIndex894, depth894
								if buffer[position] != rune('L') {
									goto l886
								}
								position++
							}
						l894:
							{
								position896, tokenIndex896, depth896 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l897
								}
								position++
								goto l896
							l897:
								position, tokenIndex, depth = position896, tokenIndex896, depth896
								if buffer[position] != rune('E') {
									goto l886
								}
								position++
							}
						l896:
							{
								position898, tokenIndex898, depth898 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l899
								}
								position++
								goto l898
							l899:
								position, tokenIndex, depth = position898, tokenIndex898, depth898
								if buffer[position] != rune('N') {
									goto l886
								}
								position++
							}
						l898:
							if !rules[ruleskip]() {
								goto l886
							}
							depth--
							add(ruleSTRLEN, position887)
						}
						goto l833
					l886:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position901 := position
							depth++
							{
								position902, tokenIndex902, depth902 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l903
								}
								position++
								goto l902
							l903:
								position, tokenIndex, depth = position902, tokenIndex902, depth902
								if buffer[position] != rune('M') {
									goto l900
								}
								position++
							}
						l902:
							{
								position904, tokenIndex904, depth904 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l905
								}
								position++
								goto l904
							l905:
								position, tokenIndex, depth = position904, tokenIndex904, depth904
								if buffer[position] != rune('O') {
									goto l900
								}
								position++
							}
						l904:
							{
								position906, tokenIndex906, depth906 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l907
								}
								position++
								goto l906
							l907:
								position, tokenIndex, depth = position906, tokenIndex906, depth906
								if buffer[position] != rune('N') {
									goto l900
								}
								position++
							}
						l906:
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
									goto l900
								}
								position++
							}
						l908:
							{
								position910, tokenIndex910, depth910 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l911
								}
								position++
								goto l910
							l911:
								position, tokenIndex, depth = position910, tokenIndex910, depth910
								if buffer[position] != rune('H') {
									goto l900
								}
								position++
							}
						l910:
							if !rules[ruleskip]() {
								goto l900
							}
							depth--
							add(ruleMONTH, position901)
						}
						goto l833
					l900:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position913 := position
							depth++
							{
								position914, tokenIndex914, depth914 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l915
								}
								position++
								goto l914
							l915:
								position, tokenIndex, depth = position914, tokenIndex914, depth914
								if buffer[position] != rune('M') {
									goto l912
								}
								position++
							}
						l914:
							{
								position916, tokenIndex916, depth916 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l917
								}
								position++
								goto l916
							l917:
								position, tokenIndex, depth = position916, tokenIndex916, depth916
								if buffer[position] != rune('I') {
									goto l912
								}
								position++
							}
						l916:
							{
								position918, tokenIndex918, depth918 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l919
								}
								position++
								goto l918
							l919:
								position, tokenIndex, depth = position918, tokenIndex918, depth918
								if buffer[position] != rune('N') {
									goto l912
								}
								position++
							}
						l918:
							{
								position920, tokenIndex920, depth920 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l921
								}
								position++
								goto l920
							l921:
								position, tokenIndex, depth = position920, tokenIndex920, depth920
								if buffer[position] != rune('U') {
									goto l912
								}
								position++
							}
						l920:
							{
								position922, tokenIndex922, depth922 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l923
								}
								position++
								goto l922
							l923:
								position, tokenIndex, depth = position922, tokenIndex922, depth922
								if buffer[position] != rune('T') {
									goto l912
								}
								position++
							}
						l922:
							{
								position924, tokenIndex924, depth924 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l925
								}
								position++
								goto l924
							l925:
								position, tokenIndex, depth = position924, tokenIndex924, depth924
								if buffer[position] != rune('E') {
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
							add(ruleMINUTES, position913)
						}
						goto l833
					l912:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position929 := position
							depth++
							{
								position930, tokenIndex930, depth930 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l931
								}
								position++
								goto l930
							l931:
								position, tokenIndex, depth = position930, tokenIndex930, depth930
								if buffer[position] != rune('S') {
									goto l928
								}
								position++
							}
						l930:
							{
								position932, tokenIndex932, depth932 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l933
								}
								position++
								goto l932
							l933:
								position, tokenIndex, depth = position932, tokenIndex932, depth932
								if buffer[position] != rune('E') {
									goto l928
								}
								position++
							}
						l932:
							{
								position934, tokenIndex934, depth934 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l935
								}
								position++
								goto l934
							l935:
								position, tokenIndex, depth = position934, tokenIndex934, depth934
								if buffer[position] != rune('C') {
									goto l928
								}
								position++
							}
						l934:
							{
								position936, tokenIndex936, depth936 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l937
								}
								position++
								goto l936
							l937:
								position, tokenIndex, depth = position936, tokenIndex936, depth936
								if buffer[position] != rune('O') {
									goto l928
								}
								position++
							}
						l936:
							{
								position938, tokenIndex938, depth938 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l939
								}
								position++
								goto l938
							l939:
								position, tokenIndex, depth = position938, tokenIndex938, depth938
								if buffer[position] != rune('N') {
									goto l928
								}
								position++
							}
						l938:
							{
								position940, tokenIndex940, depth940 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l941
								}
								position++
								goto l940
							l941:
								position, tokenIndex, depth = position940, tokenIndex940, depth940
								if buffer[position] != rune('D') {
									goto l928
								}
								position++
							}
						l940:
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
									goto l928
								}
								position++
							}
						l942:
							if !rules[ruleskip]() {
								goto l928
							}
							depth--
							add(ruleSECONDS, position929)
						}
						goto l833
					l928:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position945 := position
							depth++
							{
								position946, tokenIndex946, depth946 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l947
								}
								position++
								goto l946
							l947:
								position, tokenIndex, depth = position946, tokenIndex946, depth946
								if buffer[position] != rune('T') {
									goto l944
								}
								position++
							}
						l946:
							{
								position948, tokenIndex948, depth948 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l949
								}
								position++
								goto l948
							l949:
								position, tokenIndex, depth = position948, tokenIndex948, depth948
								if buffer[position] != rune('I') {
									goto l944
								}
								position++
							}
						l948:
							{
								position950, tokenIndex950, depth950 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l951
								}
								position++
								goto l950
							l951:
								position, tokenIndex, depth = position950, tokenIndex950, depth950
								if buffer[position] != rune('M') {
									goto l944
								}
								position++
							}
						l950:
							{
								position952, tokenIndex952, depth952 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l953
								}
								position++
								goto l952
							l953:
								position, tokenIndex, depth = position952, tokenIndex952, depth952
								if buffer[position] != rune('E') {
									goto l944
								}
								position++
							}
						l952:
							{
								position954, tokenIndex954, depth954 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l955
								}
								position++
								goto l954
							l955:
								position, tokenIndex, depth = position954, tokenIndex954, depth954
								if buffer[position] != rune('Z') {
									goto l944
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
									goto l944
								}
								position++
							}
						l956:
							{
								position958, tokenIndex958, depth958 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l959
								}
								position++
								goto l958
							l959:
								position, tokenIndex, depth = position958, tokenIndex958, depth958
								if buffer[position] != rune('N') {
									goto l944
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
									goto l944
								}
								position++
							}
						l960:
							if !rules[ruleskip]() {
								goto l944
							}
							depth--
							add(ruleTIMEZONE, position945)
						}
						goto l833
					l944:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
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
							if buffer[position] != rune('1') {
								goto l962
							}
							position++
							if !rules[ruleskip]() {
								goto l962
							}
							depth--
							add(ruleSHA1, position963)
						}
						goto l833
					l962:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position971 := position
							depth++
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
									goto l970
								}
								position++
							}
						l972:
							{
								position974, tokenIndex974, depth974 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l975
								}
								position++
								goto l974
							l975:
								position, tokenIndex, depth = position974, tokenIndex974, depth974
								if buffer[position] != rune('H') {
									goto l970
								}
								position++
							}
						l974:
							{
								position976, tokenIndex976, depth976 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l977
								}
								position++
								goto l976
							l977:
								position, tokenIndex, depth = position976, tokenIndex976, depth976
								if buffer[position] != rune('A') {
									goto l970
								}
								position++
							}
						l976:
							if buffer[position] != rune('2') {
								goto l970
							}
							position++
							if buffer[position] != rune('5') {
								goto l970
							}
							position++
							if buffer[position] != rune('6') {
								goto l970
							}
							position++
							if !rules[ruleskip]() {
								goto l970
							}
							depth--
							add(ruleSHA256, position971)
						}
						goto l833
					l970:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position979 := position
							depth++
							{
								position980, tokenIndex980, depth980 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l981
								}
								position++
								goto l980
							l981:
								position, tokenIndex, depth = position980, tokenIndex980, depth980
								if buffer[position] != rune('S') {
									goto l978
								}
								position++
							}
						l980:
							{
								position982, tokenIndex982, depth982 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l983
								}
								position++
								goto l982
							l983:
								position, tokenIndex, depth = position982, tokenIndex982, depth982
								if buffer[position] != rune('H') {
									goto l978
								}
								position++
							}
						l982:
							{
								position984, tokenIndex984, depth984 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l985
								}
								position++
								goto l984
							l985:
								position, tokenIndex, depth = position984, tokenIndex984, depth984
								if buffer[position] != rune('A') {
									goto l978
								}
								position++
							}
						l984:
							if buffer[position] != rune('3') {
								goto l978
							}
							position++
							if buffer[position] != rune('8') {
								goto l978
							}
							position++
							if buffer[position] != rune('4') {
								goto l978
							}
							position++
							if !rules[ruleskip]() {
								goto l978
							}
							depth--
							add(ruleSHA384, position979)
						}
						goto l833
					l978:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position987 := position
							depth++
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
									goto l986
								}
								position++
							}
						l988:
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
									goto l986
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
									goto l986
								}
								position++
							}
						l992:
							{
								position994, tokenIndex994, depth994 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l995
								}
								position++
								goto l994
							l995:
								position, tokenIndex, depth = position994, tokenIndex994, depth994
								if buffer[position] != rune('R') {
									goto l986
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
									goto l986
								}
								position++
							}
						l996:
							if !rules[ruleskip]() {
								goto l986
							}
							depth--
							add(ruleISIRI, position987)
						}
						goto l833
					l986:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position999 := position
							depth++
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
									goto l998
								}
								position++
							}
						l1000:
							{
								position1002, tokenIndex1002, depth1002 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1003
								}
								position++
								goto l1002
							l1003:
								position, tokenIndex, depth = position1002, tokenIndex1002, depth1002
								if buffer[position] != rune('S') {
									goto l998
								}
								position++
							}
						l1002:
							{
								position1004, tokenIndex1004, depth1004 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l1005
								}
								position++
								goto l1004
							l1005:
								position, tokenIndex, depth = position1004, tokenIndex1004, depth1004
								if buffer[position] != rune('U') {
									goto l998
								}
								position++
							}
						l1004:
							{
								position1006, tokenIndex1006, depth1006 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1007
								}
								position++
								goto l1006
							l1007:
								position, tokenIndex, depth = position1006, tokenIndex1006, depth1006
								if buffer[position] != rune('R') {
									goto l998
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
									goto l998
								}
								position++
							}
						l1008:
							if !rules[ruleskip]() {
								goto l998
							}
							depth--
							add(ruleISURI, position999)
						}
						goto l833
					l998:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
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
								if buffer[position] != rune('b') {
									goto l1017
								}
								position++
								goto l1016
							l1017:
								position, tokenIndex, depth = position1016, tokenIndex1016, depth1016
								if buffer[position] != rune('B') {
									goto l1010
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
									goto l1010
								}
								position++
							}
						l1018:
							{
								position1020, tokenIndex1020, depth1020 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1021
								}
								position++
								goto l1020
							l1021:
								position, tokenIndex, depth = position1020, tokenIndex1020, depth1020
								if buffer[position] != rune('A') {
									goto l1010
								}
								position++
							}
						l1020:
							{
								position1022, tokenIndex1022, depth1022 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1023
								}
								position++
								goto l1022
							l1023:
								position, tokenIndex, depth = position1022, tokenIndex1022, depth1022
								if buffer[position] != rune('N') {
									goto l1010
								}
								position++
							}
						l1022:
							{
								position1024, tokenIndex1024, depth1024 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l1025
								}
								position++
								goto l1024
							l1025:
								position, tokenIndex, depth = position1024, tokenIndex1024, depth1024
								if buffer[position] != rune('K') {
									goto l1010
								}
								position++
							}
						l1024:
							if !rules[ruleskip]() {
								goto l1010
							}
							depth--
							add(ruleISBLANK, position1011)
						}
						goto l833
					l1010:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							position1027 := position
							depth++
							{
								position1028, tokenIndex1028, depth1028 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1029
								}
								position++
								goto l1028
							l1029:
								position, tokenIndex, depth = position1028, tokenIndex1028, depth1028
								if buffer[position] != rune('I') {
									goto l1026
								}
								position++
							}
						l1028:
							{
								position1030, tokenIndex1030, depth1030 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1031
								}
								position++
								goto l1030
							l1031:
								position, tokenIndex, depth = position1030, tokenIndex1030, depth1030
								if buffer[position] != rune('S') {
									goto l1026
								}
								position++
							}
						l1030:
							{
								position1032, tokenIndex1032, depth1032 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1033
								}
								position++
								goto l1032
							l1033:
								position, tokenIndex, depth = position1032, tokenIndex1032, depth1032
								if buffer[position] != rune('L') {
									goto l1026
								}
								position++
							}
						l1032:
							{
								position1034, tokenIndex1034, depth1034 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l1035
								}
								position++
								goto l1034
							l1035:
								position, tokenIndex, depth = position1034, tokenIndex1034, depth1034
								if buffer[position] != rune('I') {
									goto l1026
								}
								position++
							}
						l1034:
							{
								position1036, tokenIndex1036, depth1036 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1037
								}
								position++
								goto l1036
							l1037:
								position, tokenIndex, depth = position1036, tokenIndex1036, depth1036
								if buffer[position] != rune('T') {
									goto l1026
								}
								position++
							}
						l1036:
							{
								position1038, tokenIndex1038, depth1038 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1039
								}
								position++
								goto l1038
							l1039:
								position, tokenIndex, depth = position1038, tokenIndex1038, depth1038
								if buffer[position] != rune('E') {
									goto l1026
								}
								position++
							}
						l1038:
							{
								position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1041
								}
								position++
								goto l1040
							l1041:
								position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
								if buffer[position] != rune('R') {
									goto l1026
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
									goto l1026
								}
								position++
							}
						l1042:
							{
								position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1045
								}
								position++
								goto l1044
							l1045:
								position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
								if buffer[position] != rune('L') {
									goto l1026
								}
								position++
							}
						l1044:
							if !rules[ruleskip]() {
								goto l1026
							}
							depth--
							add(ruleISLITERAL, position1027)
						}
						goto l833
					l1026:
						position, tokenIndex, depth = position833, tokenIndex833, depth833
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position1047 := position
									depth++
									{
										position1048, tokenIndex1048, depth1048 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1049
										}
										position++
										goto l1048
									l1049:
										position, tokenIndex, depth = position1048, tokenIndex1048, depth1048
										if buffer[position] != rune('I') {
											goto l832
										}
										position++
									}
								l1048:
									{
										position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1051
										}
										position++
										goto l1050
									l1051:
										position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
										if buffer[position] != rune('S') {
											goto l832
										}
										position++
									}
								l1050:
									{
										position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1053
										}
										position++
										goto l1052
									l1053:
										position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
										if buffer[position] != rune('N') {
											goto l832
										}
										position++
									}
								l1052:
									{
										position1054, tokenIndex1054, depth1054 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1055
										}
										position++
										goto l1054
									l1055:
										position, tokenIndex, depth = position1054, tokenIndex1054, depth1054
										if buffer[position] != rune('U') {
											goto l832
										}
										position++
									}
								l1054:
									{
										position1056, tokenIndex1056, depth1056 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1057
										}
										position++
										goto l1056
									l1057:
										position, tokenIndex, depth = position1056, tokenIndex1056, depth1056
										if buffer[position] != rune('M') {
											goto l832
										}
										position++
									}
								l1056:
									{
										position1058, tokenIndex1058, depth1058 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1059
										}
										position++
										goto l1058
									l1059:
										position, tokenIndex, depth = position1058, tokenIndex1058, depth1058
										if buffer[position] != rune('E') {
											goto l832
										}
										position++
									}
								l1058:
									{
										position1060, tokenIndex1060, depth1060 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1061
										}
										position++
										goto l1060
									l1061:
										position, tokenIndex, depth = position1060, tokenIndex1060, depth1060
										if buffer[position] != rune('R') {
											goto l832
										}
										position++
									}
								l1060:
									{
										position1062, tokenIndex1062, depth1062 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1063
										}
										position++
										goto l1062
									l1063:
										position, tokenIndex, depth = position1062, tokenIndex1062, depth1062
										if buffer[position] != rune('I') {
											goto l832
										}
										position++
									}
								l1062:
									{
										position1064, tokenIndex1064, depth1064 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1065
										}
										position++
										goto l1064
									l1065:
										position, tokenIndex, depth = position1064, tokenIndex1064, depth1064
										if buffer[position] != rune('C') {
											goto l832
										}
										position++
									}
								l1064:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleISNUMERIC, position1047)
								}
								break
							case 'S', 's':
								{
									position1066 := position
									depth++
									{
										position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1068
										}
										position++
										goto l1067
									l1068:
										position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
										if buffer[position] != rune('S') {
											goto l832
										}
										position++
									}
								l1067:
									{
										position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1070
										}
										position++
										goto l1069
									l1070:
										position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
										if buffer[position] != rune('H') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1071:
									if buffer[position] != rune('5') {
										goto l832
									}
									position++
									if buffer[position] != rune('1') {
										goto l832
									}
									position++
									if buffer[position] != rune('2') {
										goto l832
									}
									position++
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleSHA512, position1066)
								}
								break
							case 'M', 'm':
								{
									position1073 := position
									depth++
									{
										position1074, tokenIndex1074, depth1074 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1075
										}
										position++
										goto l1074
									l1075:
										position, tokenIndex, depth = position1074, tokenIndex1074, depth1074
										if buffer[position] != rune('M') {
											goto l832
										}
										position++
									}
								l1074:
									{
										position1076, tokenIndex1076, depth1076 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1077
										}
										position++
										goto l1076
									l1077:
										position, tokenIndex, depth = position1076, tokenIndex1076, depth1076
										if buffer[position] != rune('D') {
											goto l832
										}
										position++
									}
								l1076:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleMD5, position1073)
								}
								break
							case 'T', 't':
								{
									position1078 := position
									depth++
									{
										position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1080
										}
										position++
										goto l1079
									l1080:
										position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
										if buffer[position] != rune('T') {
											goto l832
										}
										position++
									}
								l1079:
									{
										position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l1082
										}
										position++
										goto l1081
									l1082:
										position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
										if buffer[position] != rune('Z') {
											goto l832
										}
										position++
									}
								l1081:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleTZ, position1078)
								}
								break
							case 'H', 'h':
								{
									position1083 := position
									depth++
									{
										position1084, tokenIndex1084, depth1084 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1085
										}
										position++
										goto l1084
									l1085:
										position, tokenIndex, depth = position1084, tokenIndex1084, depth1084
										if buffer[position] != rune('H') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1086:
									{
										position1088, tokenIndex1088, depth1088 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1089
										}
										position++
										goto l1088
									l1089:
										position, tokenIndex, depth = position1088, tokenIndex1088, depth1088
										if buffer[position] != rune('U') {
											goto l832
										}
										position++
									}
								l1088:
									{
										position1090, tokenIndex1090, depth1090 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1091
										}
										position++
										goto l1090
									l1091:
										position, tokenIndex, depth = position1090, tokenIndex1090, depth1090
										if buffer[position] != rune('R') {
											goto l832
										}
										position++
									}
								l1090:
									{
										position1092, tokenIndex1092, depth1092 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1093
										}
										position++
										goto l1092
									l1093:
										position, tokenIndex, depth = position1092, tokenIndex1092, depth1092
										if buffer[position] != rune('S') {
											goto l832
										}
										position++
									}
								l1092:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleHOURS, position1083)
								}
								break
							case 'D', 'd':
								{
									position1094 := position
									depth++
									{
										position1095, tokenIndex1095, depth1095 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1096
										}
										position++
										goto l1095
									l1096:
										position, tokenIndex, depth = position1095, tokenIndex1095, depth1095
										if buffer[position] != rune('D') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1097:
									{
										position1099, tokenIndex1099, depth1099 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1100
										}
										position++
										goto l1099
									l1100:
										position, tokenIndex, depth = position1099, tokenIndex1099, depth1099
										if buffer[position] != rune('Y') {
											goto l832
										}
										position++
									}
								l1099:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleDAY, position1094)
								}
								break
							case 'Y', 'y':
								{
									position1101 := position
									depth++
									{
										position1102, tokenIndex1102, depth1102 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l1103
										}
										position++
										goto l1102
									l1103:
										position, tokenIndex, depth = position1102, tokenIndex1102, depth1102
										if buffer[position] != rune('Y') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1104:
									{
										position1106, tokenIndex1106, depth1106 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1107
										}
										position++
										goto l1106
									l1107:
										position, tokenIndex, depth = position1106, tokenIndex1106, depth1106
										if buffer[position] != rune('A') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1108:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleYEAR, position1101)
								}
								break
							case 'E', 'e':
								{
									position1110 := position
									depth++
									{
										position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1112
										}
										position++
										goto l1111
									l1112:
										position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
										if buffer[position] != rune('E') {
											goto l832
										}
										position++
									}
								l1111:
									{
										position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1114
										}
										position++
										goto l1113
									l1114:
										position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
										if buffer[position] != rune('N') {
											goto l832
										}
										position++
									}
								l1113:
									{
										position1115, tokenIndex1115, depth1115 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1116
										}
										position++
										goto l1115
									l1116:
										position, tokenIndex, depth = position1115, tokenIndex1115, depth1115
										if buffer[position] != rune('C') {
											goto l832
										}
										position++
									}
								l1115:
									{
										position1117, tokenIndex1117, depth1117 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1118
										}
										position++
										goto l1117
									l1118:
										position, tokenIndex, depth = position1117, tokenIndex1117, depth1117
										if buffer[position] != rune('O') {
											goto l832
										}
										position++
									}
								l1117:
									{
										position1119, tokenIndex1119, depth1119 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1120
										}
										position++
										goto l1119
									l1120:
										position, tokenIndex, depth = position1119, tokenIndex1119, depth1119
										if buffer[position] != rune('D') {
											goto l832
										}
										position++
									}
								l1119:
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
											goto l832
										}
										position++
									}
								l1121:
									if buffer[position] != rune('_') {
										goto l832
									}
									position++
									{
										position1123, tokenIndex1123, depth1123 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1124
										}
										position++
										goto l1123
									l1124:
										position, tokenIndex, depth = position1123, tokenIndex1123, depth1123
										if buffer[position] != rune('F') {
											goto l832
										}
										position++
									}
								l1123:
									{
										position1125, tokenIndex1125, depth1125 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1126
										}
										position++
										goto l1125
									l1126:
										position, tokenIndex, depth = position1125, tokenIndex1125, depth1125
										if buffer[position] != rune('O') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1127:
									if buffer[position] != rune('_') {
										goto l832
									}
									position++
									{
										position1129, tokenIndex1129, depth1129 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1130
										}
										position++
										goto l1129
									l1130:
										position, tokenIndex, depth = position1129, tokenIndex1129, depth1129
										if buffer[position] != rune('U') {
											goto l832
										}
										position++
									}
								l1129:
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
											goto l832
										}
										position++
									}
								l1131:
									{
										position1133, tokenIndex1133, depth1133 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1134
										}
										position++
										goto l1133
									l1134:
										position, tokenIndex, depth = position1133, tokenIndex1133, depth1133
										if buffer[position] != rune('I') {
											goto l832
										}
										position++
									}
								l1133:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleENCODEFORURI, position1110)
								}
								break
							case 'L', 'l':
								{
									position1135 := position
									depth++
									{
										position1136, tokenIndex1136, depth1136 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1137
										}
										position++
										goto l1136
									l1137:
										position, tokenIndex, depth = position1136, tokenIndex1136, depth1136
										if buffer[position] != rune('L') {
											goto l832
										}
										position++
									}
								l1136:
									{
										position1138, tokenIndex1138, depth1138 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1139
										}
										position++
										goto l1138
									l1139:
										position, tokenIndex, depth = position1138, tokenIndex1138, depth1138
										if buffer[position] != rune('C') {
											goto l832
										}
										position++
									}
								l1138:
									{
										position1140, tokenIndex1140, depth1140 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1141
										}
										position++
										goto l1140
									l1141:
										position, tokenIndex, depth = position1140, tokenIndex1140, depth1140
										if buffer[position] != rune('A') {
											goto l832
										}
										position++
									}
								l1140:
									{
										position1142, tokenIndex1142, depth1142 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1143
										}
										position++
										goto l1142
									l1143:
										position, tokenIndex, depth = position1142, tokenIndex1142, depth1142
										if buffer[position] != rune('S') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1144:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleLCASE, position1135)
								}
								break
							case 'U', 'u':
								{
									position1146 := position
									depth++
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
											goto l832
										}
										position++
									}
								l1147:
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
											goto l832
										}
										position++
									}
								l1149:
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
											goto l832
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
											goto l832
										}
										position++
									}
								l1153:
									{
										position1155, tokenIndex1155, depth1155 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1156
										}
										position++
										goto l1155
									l1156:
										position, tokenIndex, depth = position1155, tokenIndex1155, depth1155
										if buffer[position] != rune('E') {
											goto l832
										}
										position++
									}
								l1155:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleUCASE, position1146)
								}
								break
							case 'F', 'f':
								{
									position1157 := position
									depth++
									{
										position1158, tokenIndex1158, depth1158 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l1159
										}
										position++
										goto l1158
									l1159:
										position, tokenIndex, depth = position1158, tokenIndex1158, depth1158
										if buffer[position] != rune('F') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1160:
									{
										position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1163
										}
										position++
										goto l1162
									l1163:
										position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
										if buffer[position] != rune('O') {
											goto l832
										}
										position++
									}
								l1162:
									{
										position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1165
										}
										position++
										goto l1164
									l1165:
										position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
										if buffer[position] != rune('O') {
											goto l832
										}
										position++
									}
								l1164:
									{
										position1166, tokenIndex1166, depth1166 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1167
										}
										position++
										goto l1166
									l1167:
										position, tokenIndex, depth = position1166, tokenIndex1166, depth1166
										if buffer[position] != rune('R') {
											goto l832
										}
										position++
									}
								l1166:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleFLOOR, position1157)
								}
								break
							case 'R', 'r':
								{
									position1168 := position
									depth++
									{
										position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1170
										}
										position++
										goto l1169
									l1170:
										position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
										if buffer[position] != rune('R') {
											goto l832
										}
										position++
									}
								l1169:
									{
										position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1172
										}
										position++
										goto l1171
									l1172:
										position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
										if buffer[position] != rune('O') {
											goto l832
										}
										position++
									}
								l1171:
									{
										position1173, tokenIndex1173, depth1173 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1174
										}
										position++
										goto l1173
									l1174:
										position, tokenIndex, depth = position1173, tokenIndex1173, depth1173
										if buffer[position] != rune('U') {
											goto l832
										}
										position++
									}
								l1173:
									{
										position1175, tokenIndex1175, depth1175 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1176
										}
										position++
										goto l1175
									l1176:
										position, tokenIndex, depth = position1175, tokenIndex1175, depth1175
										if buffer[position] != rune('N') {
											goto l832
										}
										position++
									}
								l1175:
									{
										position1177, tokenIndex1177, depth1177 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l1178
										}
										position++
										goto l1177
									l1178:
										position, tokenIndex, depth = position1177, tokenIndex1177, depth1177
										if buffer[position] != rune('D') {
											goto l832
										}
										position++
									}
								l1177:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleROUND, position1168)
								}
								break
							case 'C', 'c':
								{
									position1179 := position
									depth++
									{
										position1180, tokenIndex1180, depth1180 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1181
										}
										position++
										goto l1180
									l1181:
										position, tokenIndex, depth = position1180, tokenIndex1180, depth1180
										if buffer[position] != rune('C') {
											goto l832
										}
										position++
									}
								l1180:
									{
										position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1183
										}
										position++
										goto l1182
									l1183:
										position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
										if buffer[position] != rune('E') {
											goto l832
										}
										position++
									}
								l1182:
									{
										position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1185
										}
										position++
										goto l1184
									l1185:
										position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
										if buffer[position] != rune('I') {
											goto l832
										}
										position++
									}
								l1184:
									{
										position1186, tokenIndex1186, depth1186 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1187
										}
										position++
										goto l1186
									l1187:
										position, tokenIndex, depth = position1186, tokenIndex1186, depth1186
										if buffer[position] != rune('L') {
											goto l832
										}
										position++
									}
								l1186:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleCEIL, position1179)
								}
								break
							default:
								{
									position1188 := position
									depth++
									{
										position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1190
										}
										position++
										goto l1189
									l1190:
										position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
										if buffer[position] != rune('A') {
											goto l832
										}
										position++
									}
								l1189:
									{
										position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1192
										}
										position++
										goto l1191
									l1192:
										position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
										if buffer[position] != rune('B') {
											goto l832
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
											goto l832
										}
										position++
									}
								l1193:
									if !rules[ruleskip]() {
										goto l832
									}
									depth--
									add(ruleABS, position1188)
								}
								break
							}
						}

					}
				l833:
					if !rules[ruleLPAREN]() {
						goto l832
					}
					if !rules[ruleexpression]() {
						goto l832
					}
					if !rules[ruleRPAREN]() {
						goto l832
					}
					goto l831
				l832:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					{
						position1196, tokenIndex1196, depth1196 := position, tokenIndex, depth
						{
							position1198 := position
							depth++
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
									goto l1197
								}
								position++
							}
						l1199:
							{
								position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1202
								}
								position++
								goto l1201
							l1202:
								position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
								if buffer[position] != rune('T') {
									goto l1197
								}
								position++
							}
						l1201:
							{
								position1203, tokenIndex1203, depth1203 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1204
								}
								position++
								goto l1203
							l1204:
								position, tokenIndex, depth = position1203, tokenIndex1203, depth1203
								if buffer[position] != rune('R') {
									goto l1197
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
									goto l1197
								}
								position++
							}
						l1205:
							{
								position1207, tokenIndex1207, depth1207 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1208
								}
								position++
								goto l1207
							l1208:
								position, tokenIndex, depth = position1207, tokenIndex1207, depth1207
								if buffer[position] != rune('T') {
									goto l1197
								}
								position++
							}
						l1207:
							{
								position1209, tokenIndex1209, depth1209 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1210
								}
								position++
								goto l1209
							l1210:
								position, tokenIndex, depth = position1209, tokenIndex1209, depth1209
								if buffer[position] != rune('A') {
									goto l1197
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
									goto l1197
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
									goto l1197
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
									goto l1197
								}
								position++
							}
						l1215:
							if !rules[ruleskip]() {
								goto l1197
							}
							depth--
							add(ruleSTRSTARTS, position1198)
						}
						goto l1196
					l1197:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
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
								if buffer[position] != rune('e') {
									goto l1226
								}
								position++
								goto l1225
							l1226:
								position, tokenIndex, depth = position1225, tokenIndex1225, depth1225
								if buffer[position] != rune('E') {
									goto l1217
								}
								position++
							}
						l1225:
							{
								position1227, tokenIndex1227, depth1227 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1228
								}
								position++
								goto l1227
							l1228:
								position, tokenIndex, depth = position1227, tokenIndex1227, depth1227
								if buffer[position] != rune('N') {
									goto l1217
								}
								position++
							}
						l1227:
							{
								position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1230
								}
								position++
								goto l1229
							l1230:
								position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
								if buffer[position] != rune('D') {
									goto l1217
								}
								position++
							}
						l1229:
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
									goto l1217
								}
								position++
							}
						l1231:
							if !rules[ruleskip]() {
								goto l1217
							}
							depth--
							add(ruleSTRENDS, position1218)
						}
						goto l1196
					l1217:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
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
								if buffer[position] != rune('b') {
									goto l1242
								}
								position++
								goto l1241
							l1242:
								position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
								if buffer[position] != rune('B') {
									goto l1233
								}
								position++
							}
						l1241:
							{
								position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1244
								}
								position++
								goto l1243
							l1244:
								position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
								if buffer[position] != rune('E') {
									goto l1233
								}
								position++
							}
						l1243:
							{
								position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1246
								}
								position++
								goto l1245
							l1246:
								position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
								if buffer[position] != rune('F') {
									goto l1233
								}
								position++
							}
						l1245:
							{
								position1247, tokenIndex1247, depth1247 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l1248
								}
								position++
								goto l1247
							l1248:
								position, tokenIndex, depth = position1247, tokenIndex1247, depth1247
								if buffer[position] != rune('O') {
									goto l1233
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
									goto l1233
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
									goto l1233
								}
								position++
							}
						l1251:
							if !rules[ruleskip]() {
								goto l1233
							}
							depth--
							add(ruleSTRBEFORE, position1234)
						}
						goto l1196
					l1233:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
						{
							position1254 := position
							depth++
							{
								position1255, tokenIndex1255, depth1255 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1256
								}
								position++
								goto l1255
							l1256:
								position, tokenIndex, depth = position1255, tokenIndex1255, depth1255
								if buffer[position] != rune('S') {
									goto l1253
								}
								position++
							}
						l1255:
							{
								position1257, tokenIndex1257, depth1257 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1258
								}
								position++
								goto l1257
							l1258:
								position, tokenIndex, depth = position1257, tokenIndex1257, depth1257
								if buffer[position] != rune('T') {
									goto l1253
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
									goto l1253
								}
								position++
							}
						l1259:
							{
								position1261, tokenIndex1261, depth1261 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1262
								}
								position++
								goto l1261
							l1262:
								position, tokenIndex, depth = position1261, tokenIndex1261, depth1261
								if buffer[position] != rune('A') {
									goto l1253
								}
								position++
							}
						l1261:
							{
								position1263, tokenIndex1263, depth1263 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l1264
								}
								position++
								goto l1263
							l1264:
								position, tokenIndex, depth = position1263, tokenIndex1263, depth1263
								if buffer[position] != rune('F') {
									goto l1253
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
									goto l1253
								}
								position++
							}
						l1265:
							{
								position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l1268
								}
								position++
								goto l1267
							l1268:
								position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
								if buffer[position] != rune('E') {
									goto l1253
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
									goto l1253
								}
								position++
							}
						l1269:
							if !rules[ruleskip]() {
								goto l1253
							}
							depth--
							add(ruleSTRAFTER, position1254)
						}
						goto l1196
					l1253:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
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
								if buffer[position] != rune('l') {
									goto l1280
								}
								position++
								goto l1279
							l1280:
								position, tokenIndex, depth = position1279, tokenIndex1279, depth1279
								if buffer[position] != rune('L') {
									goto l1271
								}
								position++
							}
						l1279:
							{
								position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1282
								}
								position++
								goto l1281
							l1282:
								position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
								if buffer[position] != rune('A') {
									goto l1271
								}
								position++
							}
						l1281:
							{
								position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1284
								}
								position++
								goto l1283
							l1284:
								position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
								if buffer[position] != rune('N') {
									goto l1271
								}
								position++
							}
						l1283:
							{
								position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1286
								}
								position++
								goto l1285
							l1286:
								position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
								if buffer[position] != rune('G') {
									goto l1271
								}
								position++
							}
						l1285:
							if !rules[ruleskip]() {
								goto l1271
							}
							depth--
							add(ruleSTRLANG, position1272)
						}
						goto l1196
					l1271:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
						{
							position1288 := position
							depth++
							{
								position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l1290
								}
								position++
								goto l1289
							l1290:
								position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
								if buffer[position] != rune('S') {
									goto l1287
								}
								position++
							}
						l1289:
							{
								position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1292
								}
								position++
								goto l1291
							l1292:
								position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
								if buffer[position] != rune('T') {
									goto l1287
								}
								position++
							}
						l1291:
							{
								position1293, tokenIndex1293, depth1293 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1294
								}
								position++
								goto l1293
							l1294:
								position, tokenIndex, depth = position1293, tokenIndex1293, depth1293
								if buffer[position] != rune('R') {
									goto l1287
								}
								position++
							}
						l1293:
							{
								position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1296
								}
								position++
								goto l1295
							l1296:
								position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
								if buffer[position] != rune('D') {
									goto l1287
								}
								position++
							}
						l1295:
							{
								position1297, tokenIndex1297, depth1297 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1298
								}
								position++
								goto l1297
							l1298:
								position, tokenIndex, depth = position1297, tokenIndex1297, depth1297
								if buffer[position] != rune('T') {
									goto l1287
								}
								position++
							}
						l1297:
							if !rules[ruleskip]() {
								goto l1287
							}
							depth--
							add(ruleSTRDT, position1288)
						}
						goto l1196
					l1287:
						position, tokenIndex, depth = position1196, tokenIndex1196, depth1196
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1300 := position
									depth++
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
											goto l1195
										}
										position++
									}
								l1301:
									{
										position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1304
										}
										position++
										goto l1303
									l1304:
										position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
										if buffer[position] != rune('A') {
											goto l1195
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
											goto l1195
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
											goto l1195
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
											goto l1195
										}
										position++
									}
								l1309:
									{
										position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1312
										}
										position++
										goto l1311
									l1312:
										position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
										if buffer[position] != rune('E') {
											goto l1195
										}
										position++
									}
								l1311:
									{
										position1313, tokenIndex1313, depth1313 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1314
										}
										position++
										goto l1313
									l1314:
										position, tokenIndex, depth = position1313, tokenIndex1313, depth1313
										if buffer[position] != rune('R') {
											goto l1195
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
											goto l1195
										}
										position++
									}
								l1315:
									if !rules[ruleskip]() {
										goto l1195
									}
									depth--
									add(ruleSAMETERM, position1300)
								}
								break
							case 'C', 'c':
								{
									position1317 := position
									depth++
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
											goto l1195
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
											goto l1195
										}
										position++
									}
								l1320:
									{
										position1322, tokenIndex1322, depth1322 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1323
										}
										position++
										goto l1322
									l1323:
										position, tokenIndex, depth = position1322, tokenIndex1322, depth1322
										if buffer[position] != rune('N') {
											goto l1195
										}
										position++
									}
								l1322:
									{
										position1324, tokenIndex1324, depth1324 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1325
										}
										position++
										goto l1324
									l1325:
										position, tokenIndex, depth = position1324, tokenIndex1324, depth1324
										if buffer[position] != rune('T') {
											goto l1195
										}
										position++
									}
								l1324:
									{
										position1326, tokenIndex1326, depth1326 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1327
										}
										position++
										goto l1326
									l1327:
										position, tokenIndex, depth = position1326, tokenIndex1326, depth1326
										if buffer[position] != rune('A') {
											goto l1195
										}
										position++
									}
								l1326:
									{
										position1328, tokenIndex1328, depth1328 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1329
										}
										position++
										goto l1328
									l1329:
										position, tokenIndex, depth = position1328, tokenIndex1328, depth1328
										if buffer[position] != rune('I') {
											goto l1195
										}
										position++
									}
								l1328:
									{
										position1330, tokenIndex1330, depth1330 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1331
										}
										position++
										goto l1330
									l1331:
										position, tokenIndex, depth = position1330, tokenIndex1330, depth1330
										if buffer[position] != rune('N') {
											goto l1195
										}
										position++
									}
								l1330:
									{
										position1332, tokenIndex1332, depth1332 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1333
										}
										position++
										goto l1332
									l1333:
										position, tokenIndex, depth = position1332, tokenIndex1332, depth1332
										if buffer[position] != rune('S') {
											goto l1195
										}
										position++
									}
								l1332:
									if !rules[ruleskip]() {
										goto l1195
									}
									depth--
									add(ruleCONTAINS, position1317)
								}
								break
							default:
								{
									position1334 := position
									depth++
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
											goto l1195
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
											goto l1195
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
											goto l1195
										}
										position++
									}
								l1339:
									{
										position1341, tokenIndex1341, depth1341 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1342
										}
										position++
										goto l1341
									l1342:
										position, tokenIndex, depth = position1341, tokenIndex1341, depth1341
										if buffer[position] != rune('G') {
											goto l1195
										}
										position++
									}
								l1341:
									{
										position1343, tokenIndex1343, depth1343 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1344
										}
										position++
										goto l1343
									l1344:
										position, tokenIndex, depth = position1343, tokenIndex1343, depth1343
										if buffer[position] != rune('M') {
											goto l1195
										}
										position++
									}
								l1343:
									{
										position1345, tokenIndex1345, depth1345 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1346
										}
										position++
										goto l1345
									l1346:
										position, tokenIndex, depth = position1345, tokenIndex1345, depth1345
										if buffer[position] != rune('A') {
											goto l1195
										}
										position++
									}
								l1345:
									{
										position1347, tokenIndex1347, depth1347 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1348
										}
										position++
										goto l1347
									l1348:
										position, tokenIndex, depth = position1347, tokenIndex1347, depth1347
										if buffer[position] != rune('T') {
											goto l1195
										}
										position++
									}
								l1347:
									{
										position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1350
										}
										position++
										goto l1349
									l1350:
										position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
										if buffer[position] != rune('C') {
											goto l1195
										}
										position++
									}
								l1349:
									{
										position1351, tokenIndex1351, depth1351 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1352
										}
										position++
										goto l1351
									l1352:
										position, tokenIndex, depth = position1351, tokenIndex1351, depth1351
										if buffer[position] != rune('H') {
											goto l1195
										}
										position++
									}
								l1351:
									{
										position1353, tokenIndex1353, depth1353 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1354
										}
										position++
										goto l1353
									l1354:
										position, tokenIndex, depth = position1353, tokenIndex1353, depth1353
										if buffer[position] != rune('E') {
											goto l1195
										}
										position++
									}
								l1353:
									{
										position1355, tokenIndex1355, depth1355 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1356
										}
										position++
										goto l1355
									l1356:
										position, tokenIndex, depth = position1355, tokenIndex1355, depth1355
										if buffer[position] != rune('S') {
											goto l1195
										}
										position++
									}
								l1355:
									if !rules[ruleskip]() {
										goto l1195
									}
									depth--
									add(ruleLANGMATCHES, position1334)
								}
								break
							}
						}

					}
				l1196:
					if !rules[ruleLPAREN]() {
						goto l1195
					}
					if !rules[ruleexpression]() {
						goto l1195
					}
					if !rules[ruleCOMMA]() {
						goto l1195
					}
					if !rules[ruleexpression]() {
						goto l1195
					}
					if !rules[ruleRPAREN]() {
						goto l1195
					}
					goto l831
				l1195:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					{
						position1358 := position
						depth++
						{
							position1359, tokenIndex1359, depth1359 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1360
							}
							position++
							goto l1359
						l1360:
							position, tokenIndex, depth = position1359, tokenIndex1359, depth1359
							if buffer[position] != rune('B') {
								goto l1357
							}
							position++
						}
					l1359:
						{
							position1361, tokenIndex1361, depth1361 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1362
							}
							position++
							goto l1361
						l1362:
							position, tokenIndex, depth = position1361, tokenIndex1361, depth1361
							if buffer[position] != rune('O') {
								goto l1357
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
								goto l1357
							}
							position++
						}
					l1363:
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
								goto l1357
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
								goto l1357
							}
							position++
						}
					l1367:
						if !rules[ruleskip]() {
							goto l1357
						}
						depth--
						add(ruleBOUND, position1358)
					}
					if !rules[ruleLPAREN]() {
						goto l1357
					}
					if !rules[rulevar]() {
						goto l1357
					}
					if !rules[ruleRPAREN]() {
						goto l1357
					}
					goto l831
				l1357:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1371 := position
								depth++
								{
									position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
									if buffer[position] != rune('s') {
										goto l1373
									}
									position++
									goto l1372
								l1373:
									position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
									if buffer[position] != rune('S') {
										goto l1369
									}
									position++
								}
							l1372:
								{
									position1374, tokenIndex1374, depth1374 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1375
									}
									position++
									goto l1374
								l1375:
									position, tokenIndex, depth = position1374, tokenIndex1374, depth1374
									if buffer[position] != rune('T') {
										goto l1369
									}
									position++
								}
							l1374:
								{
									position1376, tokenIndex1376, depth1376 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1377
									}
									position++
									goto l1376
								l1377:
									position, tokenIndex, depth = position1376, tokenIndex1376, depth1376
									if buffer[position] != rune('R') {
										goto l1369
									}
									position++
								}
							l1376:
								{
									position1378, tokenIndex1378, depth1378 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1379
									}
									position++
									goto l1378
								l1379:
									position, tokenIndex, depth = position1378, tokenIndex1378, depth1378
									if buffer[position] != rune('U') {
										goto l1369
									}
									position++
								}
							l1378:
								{
									position1380, tokenIndex1380, depth1380 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1381
									}
									position++
									goto l1380
								l1381:
									position, tokenIndex, depth = position1380, tokenIndex1380, depth1380
									if buffer[position] != rune('U') {
										goto l1369
									}
									position++
								}
							l1380:
								{
									position1382, tokenIndex1382, depth1382 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1383
									}
									position++
									goto l1382
								l1383:
									position, tokenIndex, depth = position1382, tokenIndex1382, depth1382
									if buffer[position] != rune('I') {
										goto l1369
									}
									position++
								}
							l1382:
								{
									position1384, tokenIndex1384, depth1384 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1385
									}
									position++
									goto l1384
								l1385:
									position, tokenIndex, depth = position1384, tokenIndex1384, depth1384
									if buffer[position] != rune('D') {
										goto l1369
									}
									position++
								}
							l1384:
								if !rules[ruleskip]() {
									goto l1369
								}
								depth--
								add(ruleSTRUUID, position1371)
							}
							break
						case 'U', 'u':
							{
								position1386 := position
								depth++
								{
									position1387, tokenIndex1387, depth1387 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1388
									}
									position++
									goto l1387
								l1388:
									position, tokenIndex, depth = position1387, tokenIndex1387, depth1387
									if buffer[position] != rune('U') {
										goto l1369
									}
									position++
								}
							l1387:
								{
									position1389, tokenIndex1389, depth1389 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1390
									}
									position++
									goto l1389
								l1390:
									position, tokenIndex, depth = position1389, tokenIndex1389, depth1389
									if buffer[position] != rune('U') {
										goto l1369
									}
									position++
								}
							l1389:
								{
									position1391, tokenIndex1391, depth1391 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1392
									}
									position++
									goto l1391
								l1392:
									position, tokenIndex, depth = position1391, tokenIndex1391, depth1391
									if buffer[position] != rune('I') {
										goto l1369
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
										goto l1369
									}
									position++
								}
							l1393:
								if !rules[ruleskip]() {
									goto l1369
								}
								depth--
								add(ruleUUID, position1386)
							}
							break
						case 'N', 'n':
							{
								position1395 := position
								depth++
								{
									position1396, tokenIndex1396, depth1396 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1397
									}
									position++
									goto l1396
								l1397:
									position, tokenIndex, depth = position1396, tokenIndex1396, depth1396
									if buffer[position] != rune('N') {
										goto l1369
									}
									position++
								}
							l1396:
								{
									position1398, tokenIndex1398, depth1398 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1399
									}
									position++
									goto l1398
								l1399:
									position, tokenIndex, depth = position1398, tokenIndex1398, depth1398
									if buffer[position] != rune('O') {
										goto l1369
									}
									position++
								}
							l1398:
								{
									position1400, tokenIndex1400, depth1400 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1401
									}
									position++
									goto l1400
								l1401:
									position, tokenIndex, depth = position1400, tokenIndex1400, depth1400
									if buffer[position] != rune('W') {
										goto l1369
									}
									position++
								}
							l1400:
								if !rules[ruleskip]() {
									goto l1369
								}
								depth--
								add(ruleNOW, position1395)
							}
							break
						default:
							{
								position1402 := position
								depth++
								{
									position1403, tokenIndex1403, depth1403 := position, tokenIndex, depth
									if buffer[position] != rune('r') {
										goto l1404
									}
									position++
									goto l1403
								l1404:
									position, tokenIndex, depth = position1403, tokenIndex1403, depth1403
									if buffer[position] != rune('R') {
										goto l1369
									}
									position++
								}
							l1403:
								{
									position1405, tokenIndex1405, depth1405 := position, tokenIndex, depth
									if buffer[position] != rune('a') {
										goto l1406
									}
									position++
									goto l1405
								l1406:
									position, tokenIndex, depth = position1405, tokenIndex1405, depth1405
									if buffer[position] != rune('A') {
										goto l1369
									}
									position++
								}
							l1405:
								{
									position1407, tokenIndex1407, depth1407 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1408
									}
									position++
									goto l1407
								l1408:
									position, tokenIndex, depth = position1407, tokenIndex1407, depth1407
									if buffer[position] != rune('N') {
										goto l1369
									}
									position++
								}
							l1407:
								{
									position1409, tokenIndex1409, depth1409 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1410
									}
									position++
									goto l1409
								l1410:
									position, tokenIndex, depth = position1409, tokenIndex1409, depth1409
									if buffer[position] != rune('D') {
										goto l1369
									}
									position++
								}
							l1409:
								if !rules[ruleskip]() {
									goto l1369
								}
								depth--
								add(ruleRAND, position1402)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1369
					}
					goto l831
				l1369:
					position, tokenIndex, depth = position831, tokenIndex831, depth831
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
								{
									position1414 := position
									depth++
									{
										position1415, tokenIndex1415, depth1415 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1416
										}
										position++
										goto l1415
									l1416:
										position, tokenIndex, depth = position1415, tokenIndex1415, depth1415
										if buffer[position] != rune('E') {
											goto l1413
										}
										position++
									}
								l1415:
									{
										position1417, tokenIndex1417, depth1417 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1418
										}
										position++
										goto l1417
									l1418:
										position, tokenIndex, depth = position1417, tokenIndex1417, depth1417
										if buffer[position] != rune('X') {
											goto l1413
										}
										position++
									}
								l1417:
									{
										position1419, tokenIndex1419, depth1419 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1420
										}
										position++
										goto l1419
									l1420:
										position, tokenIndex, depth = position1419, tokenIndex1419, depth1419
										if buffer[position] != rune('I') {
											goto l1413
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
											goto l1413
										}
										position++
									}
								l1421:
									{
										position1423, tokenIndex1423, depth1423 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1424
										}
										position++
										goto l1423
									l1424:
										position, tokenIndex, depth = position1423, tokenIndex1423, depth1423
										if buffer[position] != rune('T') {
											goto l1413
										}
										position++
									}
								l1423:
									{
										position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1426
										}
										position++
										goto l1425
									l1426:
										position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
										if buffer[position] != rune('S') {
											goto l1413
										}
										position++
									}
								l1425:
									if !rules[ruleskip]() {
										goto l1413
									}
									depth--
									add(ruleEXISTS, position1414)
								}
								goto l1412
							l1413:
								position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
								{
									position1427 := position
									depth++
									{
										position1428, tokenIndex1428, depth1428 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1429
										}
										position++
										goto l1428
									l1429:
										position, tokenIndex, depth = position1428, tokenIndex1428, depth1428
										if buffer[position] != rune('N') {
											goto l829
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
											goto l829
										}
										position++
									}
								l1430:
									{
										position1432, tokenIndex1432, depth1432 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1433
										}
										position++
										goto l1432
									l1433:
										position, tokenIndex, depth = position1432, tokenIndex1432, depth1432
										if buffer[position] != rune('T') {
											goto l829
										}
										position++
									}
								l1432:
									if buffer[position] != rune(' ') {
										goto l829
									}
									position++
									{
										position1434, tokenIndex1434, depth1434 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1435
										}
										position++
										goto l1434
									l1435:
										position, tokenIndex, depth = position1434, tokenIndex1434, depth1434
										if buffer[position] != rune('E') {
											goto l829
										}
										position++
									}
								l1434:
									{
										position1436, tokenIndex1436, depth1436 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1437
										}
										position++
										goto l1436
									l1437:
										position, tokenIndex, depth = position1436, tokenIndex1436, depth1436
										if buffer[position] != rune('X') {
											goto l829
										}
										position++
									}
								l1436:
									{
										position1438, tokenIndex1438, depth1438 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1439
										}
										position++
										goto l1438
									l1439:
										position, tokenIndex, depth = position1438, tokenIndex1438, depth1438
										if buffer[position] != rune('I') {
											goto l829
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
											goto l829
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
											goto l829
										}
										position++
									}
								l1442:
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
											goto l829
										}
										position++
									}
								l1444:
									if !rules[ruleskip]() {
										goto l829
									}
									depth--
									add(ruleNOTEXIST, position1427)
								}
							}
						l1412:
							if !rules[rulegroupGraphPattern]() {
								goto l829
							}
							break
						case 'I', 'i':
							{
								position1446 := position
								depth++
								{
									position1447, tokenIndex1447, depth1447 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1448
									}
									position++
									goto l1447
								l1448:
									position, tokenIndex, depth = position1447, tokenIndex1447, depth1447
									if buffer[position] != rune('I') {
										goto l829
									}
									position++
								}
							l1447:
								{
									position1449, tokenIndex1449, depth1449 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1450
									}
									position++
									goto l1449
								l1450:
									position, tokenIndex, depth = position1449, tokenIndex1449, depth1449
									if buffer[position] != rune('F') {
										goto l829
									}
									position++
								}
							l1449:
								if !rules[ruleskip]() {
									goto l829
								}
								depth--
								add(ruleIF, position1446)
							}
							if !rules[ruleLPAREN]() {
								goto l829
							}
							if !rules[ruleexpression]() {
								goto l829
							}
							if !rules[ruleCOMMA]() {
								goto l829
							}
							if !rules[ruleexpression]() {
								goto l829
							}
							if !rules[ruleCOMMA]() {
								goto l829
							}
							if !rules[ruleexpression]() {
								goto l829
							}
							if !rules[ruleRPAREN]() {
								goto l829
							}
							break
						case 'C', 'c':
							{
								position1451, tokenIndex1451, depth1451 := position, tokenIndex, depth
								{
									position1453 := position
									depth++
									{
										position1454, tokenIndex1454, depth1454 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1455
										}
										position++
										goto l1454
									l1455:
										position, tokenIndex, depth = position1454, tokenIndex1454, depth1454
										if buffer[position] != rune('C') {
											goto l1452
										}
										position++
									}
								l1454:
									{
										position1456, tokenIndex1456, depth1456 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1457
										}
										position++
										goto l1456
									l1457:
										position, tokenIndex, depth = position1456, tokenIndex1456, depth1456
										if buffer[position] != rune('O') {
											goto l1452
										}
										position++
									}
								l1456:
									{
										position1458, tokenIndex1458, depth1458 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1459
										}
										position++
										goto l1458
									l1459:
										position, tokenIndex, depth = position1458, tokenIndex1458, depth1458
										if buffer[position] != rune('N') {
											goto l1452
										}
										position++
									}
								l1458:
									{
										position1460, tokenIndex1460, depth1460 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1461
										}
										position++
										goto l1460
									l1461:
										position, tokenIndex, depth = position1460, tokenIndex1460, depth1460
										if buffer[position] != rune('C') {
											goto l1452
										}
										position++
									}
								l1460:
									{
										position1462, tokenIndex1462, depth1462 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1463
										}
										position++
										goto l1462
									l1463:
										position, tokenIndex, depth = position1462, tokenIndex1462, depth1462
										if buffer[position] != rune('A') {
											goto l1452
										}
										position++
									}
								l1462:
									{
										position1464, tokenIndex1464, depth1464 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1465
										}
										position++
										goto l1464
									l1465:
										position, tokenIndex, depth = position1464, tokenIndex1464, depth1464
										if buffer[position] != rune('T') {
											goto l1452
										}
										position++
									}
								l1464:
									if !rules[ruleskip]() {
										goto l1452
									}
									depth--
									add(ruleCONCAT, position1453)
								}
								goto l1451
							l1452:
								position, tokenIndex, depth = position1451, tokenIndex1451, depth1451
								{
									position1466 := position
									depth++
									{
										position1467, tokenIndex1467, depth1467 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1468
										}
										position++
										goto l1467
									l1468:
										position, tokenIndex, depth = position1467, tokenIndex1467, depth1467
										if buffer[position] != rune('C') {
											goto l829
										}
										position++
									}
								l1467:
									{
										position1469, tokenIndex1469, depth1469 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1470
										}
										position++
										goto l1469
									l1470:
										position, tokenIndex, depth = position1469, tokenIndex1469, depth1469
										if buffer[position] != rune('O') {
											goto l829
										}
										position++
									}
								l1469:
									{
										position1471, tokenIndex1471, depth1471 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1472
										}
										position++
										goto l1471
									l1472:
										position, tokenIndex, depth = position1471, tokenIndex1471, depth1471
										if buffer[position] != rune('A') {
											goto l829
										}
										position++
									}
								l1471:
									{
										position1473, tokenIndex1473, depth1473 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1474
										}
										position++
										goto l1473
									l1474:
										position, tokenIndex, depth = position1473, tokenIndex1473, depth1473
										if buffer[position] != rune('L') {
											goto l829
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
											goto l829
										}
										position++
									}
								l1475:
									{
										position1477, tokenIndex1477, depth1477 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1478
										}
										position++
										goto l1477
									l1478:
										position, tokenIndex, depth = position1477, tokenIndex1477, depth1477
										if buffer[position] != rune('S') {
											goto l829
										}
										position++
									}
								l1477:
									{
										position1479, tokenIndex1479, depth1479 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1480
										}
										position++
										goto l1479
									l1480:
										position, tokenIndex, depth = position1479, tokenIndex1479, depth1479
										if buffer[position] != rune('C') {
											goto l829
										}
										position++
									}
								l1479:
									{
										position1481, tokenIndex1481, depth1481 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1482
										}
										position++
										goto l1481
									l1482:
										position, tokenIndex, depth = position1481, tokenIndex1481, depth1481
										if buffer[position] != rune('E') {
											goto l829
										}
										position++
									}
								l1481:
									if !rules[ruleskip]() {
										goto l829
									}
									depth--
									add(ruleCOALESCE, position1466)
								}
							}
						l1451:
							if !rules[ruleargList]() {
								goto l829
							}
							break
						case 'B', 'b':
							{
								position1483 := position
								depth++
								{
									position1484, tokenIndex1484, depth1484 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1485
									}
									position++
									goto l1484
								l1485:
									position, tokenIndex, depth = position1484, tokenIndex1484, depth1484
									if buffer[position] != rune('B') {
										goto l829
									}
									position++
								}
							l1484:
								{
									position1486, tokenIndex1486, depth1486 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1487
									}
									position++
									goto l1486
								l1487:
									position, tokenIndex, depth = position1486, tokenIndex1486, depth1486
									if buffer[position] != rune('N') {
										goto l829
									}
									position++
								}
							l1486:
								{
									position1488, tokenIndex1488, depth1488 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1489
									}
									position++
									goto l1488
								l1489:
									position, tokenIndex, depth = position1488, tokenIndex1488, depth1488
									if buffer[position] != rune('O') {
										goto l829
									}
									position++
								}
							l1488:
								{
									position1490, tokenIndex1490, depth1490 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1491
									}
									position++
									goto l1490
								l1491:
									position, tokenIndex, depth = position1490, tokenIndex1490, depth1490
									if buffer[position] != rune('D') {
										goto l829
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
										goto l829
									}
									position++
								}
							l1492:
								if !rules[ruleskip]() {
									goto l829
								}
								depth--
								add(ruleBNODE, position1483)
							}
							{
								position1494, tokenIndex1494, depth1494 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1495
								}
								if !rules[ruleexpression]() {
									goto l1495
								}
								if !rules[ruleRPAREN]() {
									goto l1495
								}
								goto l1494
							l1495:
								position, tokenIndex, depth = position1494, tokenIndex1494, depth1494
								if !rules[rulenil]() {
									goto l829
								}
							}
						l1494:
							break
						default:
							{
								position1496, tokenIndex1496, depth1496 := position, tokenIndex, depth
								{
									position1498 := position
									depth++
									{
										position1499, tokenIndex1499, depth1499 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1500
										}
										position++
										goto l1499
									l1500:
										position, tokenIndex, depth = position1499, tokenIndex1499, depth1499
										if buffer[position] != rune('S') {
											goto l1497
										}
										position++
									}
								l1499:
									{
										position1501, tokenIndex1501, depth1501 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l1502
										}
										position++
										goto l1501
									l1502:
										position, tokenIndex, depth = position1501, tokenIndex1501, depth1501
										if buffer[position] != rune('U') {
											goto l1497
										}
										position++
									}
								l1501:
									{
										position1503, tokenIndex1503, depth1503 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1504
										}
										position++
										goto l1503
									l1504:
										position, tokenIndex, depth = position1503, tokenIndex1503, depth1503
										if buffer[position] != rune('B') {
											goto l1497
										}
										position++
									}
								l1503:
									{
										position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1506
										}
										position++
										goto l1505
									l1506:
										position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
										if buffer[position] != rune('S') {
											goto l1497
										}
										position++
									}
								l1505:
									{
										position1507, tokenIndex1507, depth1507 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1508
										}
										position++
										goto l1507
									l1508:
										position, tokenIndex, depth = position1507, tokenIndex1507, depth1507
										if buffer[position] != rune('T') {
											goto l1497
										}
										position++
									}
								l1507:
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
											goto l1497
										}
										position++
									}
								l1509:
									if !rules[ruleskip]() {
										goto l1497
									}
									depth--
									add(ruleSUBSTR, position1498)
								}
								goto l1496
							l1497:
								position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
								{
									position1512 := position
									depth++
									{
										position1513, tokenIndex1513, depth1513 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1514
										}
										position++
										goto l1513
									l1514:
										position, tokenIndex, depth = position1513, tokenIndex1513, depth1513
										if buffer[position] != rune('R') {
											goto l1511
										}
										position++
									}
								l1513:
									{
										position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1516
										}
										position++
										goto l1515
									l1516:
										position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
										if buffer[position] != rune('E') {
											goto l1511
										}
										position++
									}
								l1515:
									{
										position1517, tokenIndex1517, depth1517 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1518
										}
										position++
										goto l1517
									l1518:
										position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
										if buffer[position] != rune('P') {
											goto l1511
										}
										position++
									}
								l1517:
									{
										position1519, tokenIndex1519, depth1519 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1520
										}
										position++
										goto l1519
									l1520:
										position, tokenIndex, depth = position1519, tokenIndex1519, depth1519
										if buffer[position] != rune('L') {
											goto l1511
										}
										position++
									}
								l1519:
									{
										position1521, tokenIndex1521, depth1521 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1522
										}
										position++
										goto l1521
									l1522:
										position, tokenIndex, depth = position1521, tokenIndex1521, depth1521
										if buffer[position] != rune('A') {
											goto l1511
										}
										position++
									}
								l1521:
									{
										position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1524
										}
										position++
										goto l1523
									l1524:
										position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
										if buffer[position] != rune('C') {
											goto l1511
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
											goto l1511
										}
										position++
									}
								l1525:
									if !rules[ruleskip]() {
										goto l1511
									}
									depth--
									add(ruleREPLACE, position1512)
								}
								goto l1496
							l1511:
								position, tokenIndex, depth = position1496, tokenIndex1496, depth1496
								{
									position1527 := position
									depth++
									{
										position1528, tokenIndex1528, depth1528 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1529
										}
										position++
										goto l1528
									l1529:
										position, tokenIndex, depth = position1528, tokenIndex1528, depth1528
										if buffer[position] != rune('R') {
											goto l829
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
											goto l829
										}
										position++
									}
								l1530:
									{
										position1532, tokenIndex1532, depth1532 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1533
										}
										position++
										goto l1532
									l1533:
										position, tokenIndex, depth = position1532, tokenIndex1532, depth1532
										if buffer[position] != rune('G') {
											goto l829
										}
										position++
									}
								l1532:
									{
										position1534, tokenIndex1534, depth1534 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1535
										}
										position++
										goto l1534
									l1535:
										position, tokenIndex, depth = position1534, tokenIndex1534, depth1534
										if buffer[position] != rune('E') {
											goto l829
										}
										position++
									}
								l1534:
									{
										position1536, tokenIndex1536, depth1536 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1537
										}
										position++
										goto l1536
									l1537:
										position, tokenIndex, depth = position1536, tokenIndex1536, depth1536
										if buffer[position] != rune('X') {
											goto l829
										}
										position++
									}
								l1536:
									if !rules[ruleskip]() {
										goto l829
									}
									depth--
									add(ruleREGEX, position1527)
								}
							}
						l1496:
							if !rules[ruleLPAREN]() {
								goto l829
							}
							if !rules[ruleexpression]() {
								goto l829
							}
							if !rules[ruleCOMMA]() {
								goto l829
							}
							if !rules[ruleexpression]() {
								goto l829
							}
							{
								position1538, tokenIndex1538, depth1538 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1538
								}
								if !rules[ruleexpression]() {
									goto l1538
								}
								goto l1539
							l1538:
								position, tokenIndex, depth = position1538, tokenIndex1538, depth1538
							}
						l1539:
							if !rules[ruleRPAREN]() {
								goto l829
							}
							break
						}
					}

				}
			l831:
				depth--
				add(rulebuiltinCall, position830)
			}
			return true
		l829:
			position, tokenIndex, depth = position829, tokenIndex829, depth829
			return false
		},
		/* 69 pof <- <(((<([a-z] / [A-Z])*> ':' Action10) / (<([2-9] [0-9]*)> '/' Action11) / (<((&('+') '+') | (&('_') '_') | (&('-') '-') | (&('.') '.') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))*> Action12)) '<' ws skip)> */
		func() bool {
			position1540, tokenIndex1540, depth1540 := position, tokenIndex, depth
			{
				position1541 := position
				depth++
				{
					position1542, tokenIndex1542, depth1542 := position, tokenIndex, depth
					{
						position1544 := position
						depth++
					l1545:
						{
							position1546, tokenIndex1546, depth1546 := position, tokenIndex, depth
							{
								position1547, tokenIndex1547, depth1547 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1548
								}
								position++
								goto l1547
							l1548:
								position, tokenIndex, depth = position1547, tokenIndex1547, depth1547
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1546
								}
								position++
							}
						l1547:
							goto l1545
						l1546:
							position, tokenIndex, depth = position1546, tokenIndex1546, depth1546
						}
						depth--
						add(rulePegText, position1544)
					}
					if buffer[position] != rune(':') {
						goto l1543
					}
					position++
					{
						add(ruleAction10, position)
					}
					goto l1542
				l1543:
					position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
					{
						position1551 := position
						depth++
						if c := buffer[position]; c < rune('2') || c > rune('9') {
							goto l1550
						}
						position++
					l1552:
						{
							position1553, tokenIndex1553, depth1553 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1553
							}
							position++
							goto l1552
						l1553:
							position, tokenIndex, depth = position1553, tokenIndex1553, depth1553
						}
						depth--
						add(rulePegText, position1551)
					}
					if buffer[position] != rune('/') {
						goto l1550
					}
					position++
					{
						add(ruleAction11, position)
					}
					goto l1542
				l1550:
					position, tokenIndex, depth = position1542, tokenIndex1542, depth1542
					{
						position1555 := position
						depth++
					l1556:
						{
							position1557, tokenIndex1557, depth1557 := position, tokenIndex, depth
							{
								switch buffer[position] {
								case '+':
									if buffer[position] != rune('+') {
										goto l1557
									}
									position++
									break
								case '_':
									if buffer[position] != rune('_') {
										goto l1557
									}
									position++
									break
								case '-':
									if buffer[position] != rune('-') {
										goto l1557
									}
									position++
									break
								case '.':
									if buffer[position] != rune('.') {
										goto l1557
									}
									position++
									break
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

							goto l1556
						l1557:
							position, tokenIndex, depth = position1557, tokenIndex1557, depth1557
						}
						depth--
						add(rulePegText, position1555)
					}
					{
						add(ruleAction12, position)
					}
				}
			l1542:
				if buffer[position] != rune('<') {
					goto l1540
				}
				position++
				if !rules[rulews]() {
					goto l1540
				}
				if !rules[ruleskip]() {
					goto l1540
				}
				depth--
				add(rulepof, position1541)
			}
			return true
		l1540:
			position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
			return false
		},
		/* 70 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1560, tokenIndex1560, depth1560 := position, tokenIndex, depth
			{
				position1561 := position
				depth++
				{
					position1562, tokenIndex1562, depth1562 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1563
					}
					position++
					goto l1562
				l1563:
					position, tokenIndex, depth = position1562, tokenIndex1562, depth1562
					if buffer[position] != rune('$') {
						goto l1560
					}
					position++
				}
			l1562:
				{
					position1564 := position
					depth++
					{
						position1565, tokenIndex1565, depth1565 := position, tokenIndex, depth
						if !rules[rulepnCharsU]() {
							goto l1566
						}
						goto l1565
					l1566:
						position, tokenIndex, depth = position1565, tokenIndex1565, depth1565
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1560
						}
						position++
					}
				l1565:
				l1567:
					{
						position1568, tokenIndex1568, depth1568 := position, tokenIndex, depth
						{
							position1569, tokenIndex1569, depth1569 := position, tokenIndex, depth
							if !rules[rulepnCharsU]() {
								goto l1570
							}
							goto l1569
						l1570:
							position, tokenIndex, depth = position1569, tokenIndex1569, depth1569
							{
								switch buffer[position] {
								case 'â':
									if c := buffer[position]; c < rune('‿') || c > rune('⁀') {
										goto l1568
									}
									position++
									break
								case 'Ì', 'Í':
									if c := buffer[position]; c < rune('̀') || c > rune('ͯ') {
										goto l1568
									}
									position++
									break
								case 'Â':
									if buffer[position] != rune('·') {
										goto l1568
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1568
									}
									position++
									break
								}
							}

						}
					l1569:
						goto l1567
					l1568:
						position, tokenIndex, depth = position1568, tokenIndex1568, depth1568
					}
					depth--
					add(ruleVARNAME, position1564)
				}
				if !rules[ruleskip]() {
					goto l1560
				}
				depth--
				add(rulevar, position1561)
			}
			return true
		l1560:
			position, tokenIndex, depth = position1560, tokenIndex1560, depth1560
			return false
		},
		/* 71 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1572, tokenIndex1572, depth1572 := position, tokenIndex, depth
			{
				position1573 := position
				depth++
				{
					position1574, tokenIndex1574, depth1574 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1575
					}
					goto l1574
				l1575:
					position, tokenIndex, depth = position1574, tokenIndex1574, depth1574
					{
						position1576 := position
						depth++
						{
							position1577, tokenIndex1577, depth1577 := position, tokenIndex, depth
							if !rules[rulepnPrefix]() {
								goto l1577
							}
							goto l1578
						l1577:
							position, tokenIndex, depth = position1577, tokenIndex1577, depth1577
						}
					l1578:
						if buffer[position] != rune(':') {
							goto l1572
						}
						position++
						{
							position1579 := position
							depth++
							{
								switch buffer[position] {
								case '%', '\\':
									{
										position1583 := position
										depth++
										{
											position1584, tokenIndex1584, depth1584 := position, tokenIndex, depth
											{
												position1586 := position
												depth++
												if buffer[position] != rune('%') {
													goto l1585
												}
												position++
												if !rules[rulehex]() {
													goto l1585
												}
												if !rules[rulehex]() {
													goto l1585
												}
												depth--
												add(rulepercent, position1586)
											}
											goto l1584
										l1585:
											position, tokenIndex, depth = position1584, tokenIndex1584, depth1584
											{
												position1587 := position
												depth++
												if buffer[position] != rune('\\') {
													goto l1572
												}
												position++
												{
													switch buffer[position] {
													case '%':
														if buffer[position] != rune('%') {
															goto l1572
														}
														position++
														break
													case '@':
														if buffer[position] != rune('@') {
															goto l1572
														}
														position++
														break
													case '#':
														if buffer[position] != rune('#') {
															goto l1572
														}
														position++
														break
													case '?':
														if buffer[position] != rune('?') {
															goto l1572
														}
														position++
														break
													case '/':
														if buffer[position] != rune('/') {
															goto l1572
														}
														position++
														break
													case '=':
														if buffer[position] != rune('=') {
															goto l1572
														}
														position++
														break
													case ';':
														if buffer[position] != rune(';') {
															goto l1572
														}
														position++
														break
													case ',':
														if buffer[position] != rune(',') {
															goto l1572
														}
														position++
														break
													case '+':
														if buffer[position] != rune('+') {
															goto l1572
														}
														position++
														break
													case '*':
														if buffer[position] != rune('*') {
															goto l1572
														}
														position++
														break
													case ')':
														if buffer[position] != rune(')') {
															goto l1572
														}
														position++
														break
													case '(':
														if buffer[position] != rune('(') {
															goto l1572
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l1572
														}
														position++
														break
													case '&':
														if buffer[position] != rune('&') {
															goto l1572
														}
														position++
														break
													case '$':
														if buffer[position] != rune('$') {
															goto l1572
														}
														position++
														break
													case '!':
														if buffer[position] != rune('!') {
															goto l1572
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l1572
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l1572
														}
														position++
														break
													case '~':
														if buffer[position] != rune('~') {
															goto l1572
														}
														position++
														break
													default:
														if buffer[position] != rune('_') {
															goto l1572
														}
														position++
														break
													}
												}

												depth--
												add(rulepnLocalEsc, position1587)
											}
										}
									l1584:
										depth--
										add(ruleplx, position1583)
									}
									break
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1572
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l1572
									}
									position++
									break
								default:
									if !rules[rulepnCharsU]() {
										goto l1572
									}
									break
								}
							}

						l1580:
							{
								position1581, tokenIndex1581, depth1581 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '%', '\\':
										{
											position1590 := position
											depth++
											{
												position1591, tokenIndex1591, depth1591 := position, tokenIndex, depth
												{
													position1593 := position
													depth++
													if buffer[position] != rune('%') {
														goto l1592
													}
													position++
													if !rules[rulehex]() {
														goto l1592
													}
													if !rules[rulehex]() {
														goto l1592
													}
													depth--
													add(rulepercent, position1593)
												}
												goto l1591
											l1592:
												position, tokenIndex, depth = position1591, tokenIndex1591, depth1591
												{
													position1594 := position
													depth++
													if buffer[position] != rune('\\') {
														goto l1581
													}
													position++
													{
														switch buffer[position] {
														case '%':
															if buffer[position] != rune('%') {
																goto l1581
															}
															position++
															break
														case '@':
															if buffer[position] != rune('@') {
																goto l1581
															}
															position++
															break
														case '#':
															if buffer[position] != rune('#') {
																goto l1581
															}
															position++
															break
														case '?':
															if buffer[position] != rune('?') {
																goto l1581
															}
															position++
															break
														case '/':
															if buffer[position] != rune('/') {
																goto l1581
															}
															position++
															break
														case '=':
															if buffer[position] != rune('=') {
																goto l1581
															}
															position++
															break
														case ';':
															if buffer[position] != rune(';') {
																goto l1581
															}
															position++
															break
														case ',':
															if buffer[position] != rune(',') {
																goto l1581
															}
															position++
															break
														case '+':
															if buffer[position] != rune('+') {
																goto l1581
															}
															position++
															break
														case '*':
															if buffer[position] != rune('*') {
																goto l1581
															}
															position++
															break
														case ')':
															if buffer[position] != rune(')') {
																goto l1581
															}
															position++
															break
														case '(':
															if buffer[position] != rune('(') {
																goto l1581
															}
															position++
															break
														case '\'':
															if buffer[position] != rune('\'') {
																goto l1581
															}
															position++
															break
														case '&':
															if buffer[position] != rune('&') {
																goto l1581
															}
															position++
															break
														case '$':
															if buffer[position] != rune('$') {
																goto l1581
															}
															position++
															break
														case '!':
															if buffer[position] != rune('!') {
																goto l1581
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l1581
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l1581
															}
															position++
															break
														case '~':
															if buffer[position] != rune('~') {
																goto l1581
															}
															position++
															break
														default:
															if buffer[position] != rune('_') {
																goto l1581
															}
															position++
															break
														}
													}

													depth--
													add(rulepnLocalEsc, position1594)
												}
											}
										l1591:
											depth--
											add(ruleplx, position1590)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1581
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1581
										}
										position++
										break
									default:
										if !rules[rulepnCharsU]() {
											goto l1581
										}
										break
									}
								}

								goto l1580
							l1581:
								position, tokenIndex, depth = position1581, tokenIndex1581, depth1581
							}
							depth--
							add(rulepnLocal, position1579)
						}
						if !rules[ruleskip]() {
							goto l1572
						}
						depth--
						add(ruleprefixedName, position1576)
					}
				}
			l1574:
				depth--
				add(ruleiriref, position1573)
			}
			return true
		l1572:
			position, tokenIndex, depth = position1572, tokenIndex1572, depth1572
			return false
		},
		/* 72 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1596, tokenIndex1596, depth1596 := position, tokenIndex, depth
			{
				position1597 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1596
				}
				position++
			l1598:
				{
					position1599, tokenIndex1599, depth1599 := position, tokenIndex, depth
					{
						position1600, tokenIndex1600, depth1600 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1600
						}
						position++
						goto l1599
					l1600:
						position, tokenIndex, depth = position1600, tokenIndex1600, depth1600
					}
					if !matchDot() {
						goto l1599
					}
					goto l1598
				l1599:
					position, tokenIndex, depth = position1599, tokenIndex1599, depth1599
				}
				if buffer[position] != rune('>') {
					goto l1596
				}
				position++
				if !rules[ruleskip]() {
					goto l1596
				}
				depth--
				add(ruleiri, position1597)
			}
			return true
		l1596:
			position, tokenIndex, depth = position1596, tokenIndex1596, depth1596
			return false
		},
		/* 73 prefixedName <- <(pnPrefix? ':' pnLocal skip)> */
		nil,
		/* 74 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1602, tokenIndex1602, depth1602 := position, tokenIndex, depth
			{
				position1603 := position
				depth++
				if !rules[rulestring]() {
					goto l1602
				}
				{
					position1604, tokenIndex1604, depth1604 := position, tokenIndex, depth
					{
						position1606, tokenIndex1606, depth1606 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1607
						}
						position++
						{
							position1610, tokenIndex1610, depth1610 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1611
							}
							position++
							goto l1610
						l1611:
							position, tokenIndex, depth = position1610, tokenIndex1610, depth1610
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1607
							}
							position++
						}
					l1610:
					l1608:
						{
							position1609, tokenIndex1609, depth1609 := position, tokenIndex, depth
							{
								position1612, tokenIndex1612, depth1612 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1613
								}
								position++
								goto l1612
							l1613:
								position, tokenIndex, depth = position1612, tokenIndex1612, depth1612
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1609
								}
								position++
							}
						l1612:
							goto l1608
						l1609:
							position, tokenIndex, depth = position1609, tokenIndex1609, depth1609
						}
					l1614:
						{
							position1615, tokenIndex1615, depth1615 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1615
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1615
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1615
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1615
									}
									position++
									break
								}
							}

						l1616:
							{
								position1617, tokenIndex1617, depth1617 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1617
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1617
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1617
										}
										position++
										break
									}
								}

								goto l1616
							l1617:
								position, tokenIndex, depth = position1617, tokenIndex1617, depth1617
							}
							goto l1614
						l1615:
							position, tokenIndex, depth = position1615, tokenIndex1615, depth1615
						}
						goto l1606
					l1607:
						position, tokenIndex, depth = position1606, tokenIndex1606, depth1606
						if buffer[position] != rune('^') {
							goto l1604
						}
						position++
						if buffer[position] != rune('^') {
							goto l1604
						}
						position++
						if !rules[ruleiriref]() {
							goto l1604
						}
					}
				l1606:
					goto l1605
				l1604:
					position, tokenIndex, depth = position1604, tokenIndex1604, depth1604
				}
			l1605:
				if !rules[ruleskip]() {
					goto l1602
				}
				depth--
				add(ruleliteral, position1603)
			}
			return true
		l1602:
			position, tokenIndex, depth = position1602, tokenIndex1602, depth1602
			return false
		},
		/* 75 string <- <('"' (!'"' .)* '"')> */
		func() bool {
			position1620, tokenIndex1620, depth1620 := position, tokenIndex, depth
			{
				position1621 := position
				depth++
				if buffer[position] != rune('"') {
					goto l1620
				}
				position++
			l1622:
				{
					position1623, tokenIndex1623, depth1623 := position, tokenIndex, depth
					{
						position1624, tokenIndex1624, depth1624 := position, tokenIndex, depth
						if buffer[position] != rune('"') {
							goto l1624
						}
						position++
						goto l1623
					l1624:
						position, tokenIndex, depth = position1624, tokenIndex1624, depth1624
					}
					if !matchDot() {
						goto l1623
					}
					goto l1622
				l1623:
					position, tokenIndex, depth = position1623, tokenIndex1623, depth1623
				}
				if buffer[position] != rune('"') {
					goto l1620
				}
				position++
				depth--
				add(rulestring, position1621)
			}
			return true
		l1620:
			position, tokenIndex, depth = position1620, tokenIndex1620, depth1620
			return false
		},
		/* 76 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
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
		/* 77 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 78 booleanLiteral <- <(TRUE / FALSE)> */
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
		/* 79 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 80 blankNodeLabel <- <('_' ':' (pnCharsU / [0-9]) (((pnCharsU / ((&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))* pnCharsU) / ([0-9] / '-'))? skip)> */
		nil,
		/* 81 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 82 nil <- <('(' ws* ')' skip)> */
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
		/* 83 VARNAME <- <((pnCharsU / [0-9]) (pnCharsU / ((&('â') [‿-⁀]) | (&('Ì' | 'Í') [̀-ͯ]) | (&('Â') '·') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))*)> */
		nil,
		/* 84 pnPrefix <- <(pnCharsBase pnChars*)> */
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
		/* 85 pnLocal <- <((&('%' | '\\') plx) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&(':') ':') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))+> */
		nil,
		/* 86 pnChars <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))> */
		nil,
		/* 87 pnCharsU <- <(pnCharsBase / '_')> */
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
		/* 88 pnCharsBase <- <([À-Ö] / [Ø-ö] / [Ͱ-ͽ] / [‌-‍] / [⁰-↏] / [豈-﷏] / ((&('ð' | 'ñ' | 'ò' | 'ó') [𐀀-󯿿]) | (&('ï') [ﷰ-�]) | (&('ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í') [、-퟿]) | (&('â') [Ⰰ-⿯]) | (&('Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á') [Ϳ-῿]) | (&('Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë') [ø-˿]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))> */
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
		/* 89 plx <- <(percent / pnLocalEsc)> */
		nil,
		/* 90 percent <- <('%' hex hex)> */
		nil,
		/* 91 hex <- <((&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))> */
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
		/* 92 pnLocalEsc <- <('\\' ((&('%') '%') | (&('@') '@') | (&('#') '#') | (&('?') '?') | (&('/') '/') | (&('=') '=') | (&(';') ';') | (&(',') ',') | (&('+') '+') | (&('*') '*') | (&(')') ')') | (&('(') '(') | (&('\'') '\'') | (&('&') '&') | (&('$') '$') | (&('!') '!') | (&('-') '-') | (&('.') '.') | (&('~') '~') | (&('_') '_')))> */
		nil,
		/* 93 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 94 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 95 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 96 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 97 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 98 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 99 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
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
		/* 100 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 101 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 102 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 103 LBRACE <- <('{' skip)> */
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
		/* 104 RBRACE <- <('}' skip)> */
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
		/* 105 LBRACK <- <('[' skip)> */
		nil,
		/* 106 RBRACK <- <(']' skip)> */
		nil,
		/* 107 SEMICOLON <- <(';' skip)> */
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
		/* 108 COMMA <- <(',' skip)> */
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
		/* 109 DOT <- <('.' skip)> */
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
		/* 110 COLON <- <(':' skip)> */
		nil,
		/* 111 PIPE <- <('|' skip)> */
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
		/* 112 SLASH <- <('/' skip)> */
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
		/* 113 INVERSE <- <('^' skip)> */
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
		/* 114 LPAREN <- <('(' skip)> */
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
		/* 115 RPAREN <- <(')' skip)> */
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
		/* 116 ISA <- <('a' skip)> */
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
		/* 117 NOT <- <('!' skip)> */
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
		/* 118 STAR <- <('*' skip)> */
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
		/* 119 QUESTION <- <('?' skip)> */
		nil,
		/* 120 PLUS <- <('+' skip)> */
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
		/* 121 MINUS <- <('-' skip)> */
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
		/* 122 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 123 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 124 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 125 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 126 INTEGER <- <([0-9]+ skip)> */
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
		/* 127 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 128 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 129 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 130 OR <- <('|' '|' skip)> */
		nil,
		/* 131 AND <- <('&' '&' skip)> */
		nil,
		/* 132 EQ <- <('=' skip)> */
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
		/* 133 NE <- <('!' '=' skip)> */
		nil,
		/* 134 GT <- <('>' skip)> */
		nil,
		/* 135 LT <- <('<' skip)> */
		nil,
		/* 136 LE <- <('<' '=' skip)> */
		nil,
		/* 137 GE <- <('>' '=' skip)> */
		nil,
		/* 138 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 139 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 140 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
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
		/* 141 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 142 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 143 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 144 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 145 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 146 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 147 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 148 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 149 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 150 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 151 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 152 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 153 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 154 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 155 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 156 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 157 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 158 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 159 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 160 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 161 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 162 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 163 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 164 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 165 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 166 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 167 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 168 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 169 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 170 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 171 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 172 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 173 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 174 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 175 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 176 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 177 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 178 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 179 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 180 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 181 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 182 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 183 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 184 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 185 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 186 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 187 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 188 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 189 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 190 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 191 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 192 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 193 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 194 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 195 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 196 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 197 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 198 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 199 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 200 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 201 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 202 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 203 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 204 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 205 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 206 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 207 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 208 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 209 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
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
		/* 210 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 211 GRAPH <- <(('g' / 'G') ('r' / 'R') ('a' / 'A') ('p' / 'P') ('h' / 'H') skip)> */
		nil,
		/* 212 MINUSSETOPER <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('s' / 'S') skip)> */
		nil,
		/* 213 skip <- <(ws / comment)*> */
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
							{
								position1871 := position
								depth++
								if buffer[position] != rune('#') {
									goto l1867
								}
								position++
							l1872:
								{
									position1873, tokenIndex1873, depth1873 := position, tokenIndex, depth
									{
										position1874, tokenIndex1874, depth1874 := position, tokenIndex, depth
										if !rules[ruleendOfLine]() {
											goto l1874
										}
										goto l1873
									l1874:
										position, tokenIndex, depth = position1874, tokenIndex1874, depth1874
									}
									if !matchDot() {
										goto l1873
									}
									goto l1872
								l1873:
									position, tokenIndex, depth = position1873, tokenIndex1873, depth1873
								}
								if !rules[ruleendOfLine]() {
									goto l1867
								}
								depth--
								add(rulePegText, position1871)
							}
							{
								add(ruleAction13, position)
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
		/* 214 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1876, tokenIndex1876, depth1876 := position, tokenIndex, depth
			{
				position1877 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1876
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1876
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1876
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1876
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1876
						}
						break
					}
				}

				depth--
				add(rulews, position1877)
			}
			return true
		l1876:
			position, tokenIndex, depth = position1876, tokenIndex1876, depth1876
			return false
		},
		/* 215 comment <- <(<('#' (!endOfLine .)* endOfLine)> Action13)> */
		nil,
		/* 216 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1880, tokenIndex1880, depth1880 := position, tokenIndex, depth
			{
				position1881 := position
				depth++
				{
					position1882, tokenIndex1882, depth1882 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1883
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1883
					}
					position++
					goto l1882
				l1883:
					position, tokenIndex, depth = position1882, tokenIndex1882, depth1882
					if buffer[position] != rune('\n') {
						goto l1884
					}
					position++
					goto l1882
				l1884:
					position, tokenIndex, depth = position1882, tokenIndex1882, depth1882
					if buffer[position] != rune('\r') {
						goto l1880
					}
					position++
				}
			l1882:
				depth--
				add(ruleendOfLine, position1881)
			}
			return true
		l1880:
			position, tokenIndex, depth = position1880, tokenIndex1880, depth1880
			return false
		},
		nil,
		/* 219 Action0 <- <{ p.addPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 220 Action1 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 221 Action2 <- <{ p.setSubject(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 222 Action3 <- <{ p.setSubject("?POF") }> */
		nil,
		/* 223 Action4 <- <{ p.setPredicate("?POF") }> */
		nil,
		/* 224 Action5 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 225 Action6 <- <{ p.setPredicate(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 226 Action7 <- <{ p.setObject("?POF"); p.addTriplePattern() }> */
		nil,
		/* 227 Action8 <- <{ p.setObject(p.skipComment(buffer, begin, end)); p.addTriplePattern() }> */
		nil,
		/* 228 Action9 <- <{ p.setObject("?FillVar"); p.addTriplePattern() }> */
		nil,
		/* 229 Action10 <- <{ p.setPrefix(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 230 Action11 <- <{ p.setPathLength(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 231 Action12 <- <{ p.setKeyword(p.skipComment(buffer, begin, end)) }> */
		nil,
		/* 232 Action13 <- <{ p.commentBegin = begin }> */
		nil,
	}
	p.rules = rules
}
