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
	rules  [217]func() bool
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
		/* 26 triplesBlock <- <triplesSameSubjectPath+> */
		func() bool {
			position325, tokenIndex325, depth325 := position, tokenIndex, depth
			{
				position326 := position
				depth++
				{
					position329 := position
					depth++
					{
						position330, tokenIndex330, depth330 := position, tokenIndex, depth
						if !rules[rulevarOrTerm]() {
							goto l331
						}
						if !rules[rulepropertyListPath]() {
							goto l331
						}
						goto l330
					l331:
						position, tokenIndex, depth = position330, tokenIndex330, depth330
						if !rules[ruletriplesNodePath]() {
							goto l325
						}
						{
							position332, tokenIndex332, depth332 := position, tokenIndex, depth
							if !rules[rulepropertyListPath]() {
								goto l332
							}
							goto l333
						l332:
							position, tokenIndex, depth = position332, tokenIndex332, depth332
						}
					l333:
					}
				l330:
					{
						position334, tokenIndex334, depth334 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l334
						}
						goto l335
					l334:
						position, tokenIndex, depth = position334, tokenIndex334, depth334
					}
				l335:
					depth--
					add(ruletriplesSameSubjectPath, position329)
				}
			l327:
				{
					position328, tokenIndex328, depth328 := position, tokenIndex, depth
					{
						position336 := position
						depth++
						{
							position337, tokenIndex337, depth337 := position, tokenIndex, depth
							if !rules[rulevarOrTerm]() {
								goto l338
							}
							if !rules[rulepropertyListPath]() {
								goto l338
							}
							goto l337
						l338:
							position, tokenIndex, depth = position337, tokenIndex337, depth337
							if !rules[ruletriplesNodePath]() {
								goto l328
							}
							{
								position339, tokenIndex339, depth339 := position, tokenIndex, depth
								if !rules[rulepropertyListPath]() {
									goto l339
								}
								goto l340
							l339:
								position, tokenIndex, depth = position339, tokenIndex339, depth339
							}
						l340:
						}
					l337:
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
						add(ruletriplesSameSubjectPath, position336)
					}
					goto l327
				l328:
					position, tokenIndex, depth = position328, tokenIndex328, depth328
				}
				depth--
				add(ruletriplesBlock, position326)
			}
			return true
		l325:
			position, tokenIndex, depth = position325, tokenIndex325, depth325
			return false
		},
		/* 27 triplesSameSubjectPath <- <(((varOrTerm propertyListPath) / (triplesNodePath propertyListPath?)) DOT?)> */
		nil,
		/* 28 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position344, tokenIndex344, depth344 := position, tokenIndex, depth
			{
				position345 := position
				depth++
				{
					position346, tokenIndex346, depth346 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l347
					}
					goto l346
				l347:
					position, tokenIndex, depth = position346, tokenIndex346, depth346
					{
						position348 := position
						depth++
						{
							position349, tokenIndex349, depth349 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l350
							}
							goto l349
						l350:
							position, tokenIndex, depth = position349, tokenIndex349, depth349
							{
								switch buffer[position] {
								case '(':
									if !rules[rulenil]() {
										goto l344
									}
									break
								case '[', '_':
									{
										position352 := position
										depth++
										{
											position353, tokenIndex353, depth353 := position, tokenIndex, depth
											{
												position355 := position
												depth++
												if buffer[position] != rune('_') {
													goto l354
												}
												position++
												if buffer[position] != rune(':') {
													goto l354
												}
												position++
												{
													position356, tokenIndex356, depth356 := position, tokenIndex, depth
													if !rules[rulepnCharsU]() {
														goto l357
													}
													goto l356
												l357:
													position, tokenIndex, depth = position356, tokenIndex356, depth356
													if c := buffer[position]; c < rune('0') || c > rune('9') {
														goto l354
													}
													position++
												}
											l356:
												{
													position358, tokenIndex358, depth358 := position, tokenIndex, depth
													{
														position360, tokenIndex360, depth360 := position, tokenIndex, depth
													l362:
														{
															position363, tokenIndex363, depth363 := position, tokenIndex, depth
															{
																position364, tokenIndex364, depth364 := position, tokenIndex, depth
																if !rules[rulepnCharsU]() {
																	goto l365
																}
																goto l364
															l365:
																position, tokenIndex, depth = position364, tokenIndex364, depth364
																{
																	switch buffer[position] {
																	case '.':
																		if buffer[position] != rune('.') {
																			goto l363
																		}
																		position++
																		break
																	case '-':
																		if buffer[position] != rune('-') {
																			goto l363
																		}
																		position++
																		break
																	default:
																		if c := buffer[position]; c < rune('0') || c > rune('9') {
																			goto l363
																		}
																		position++
																		break
																	}
																}

															}
														l364:
															goto l362
														l363:
															position, tokenIndex, depth = position363, tokenIndex363, depth363
														}
														if !rules[rulepnCharsU]() {
															goto l361
														}
														goto l360
													l361:
														position, tokenIndex, depth = position360, tokenIndex360, depth360
														{
															position367, tokenIndex367, depth367 := position, tokenIndex, depth
															if c := buffer[position]; c < rune('0') || c > rune('9') {
																goto l368
															}
															position++
															goto l367
														l368:
															position, tokenIndex, depth = position367, tokenIndex367, depth367
															if buffer[position] != rune('-') {
																goto l358
															}
															position++
														}
													l367:
													}
												l360:
													goto l359
												l358:
													position, tokenIndex, depth = position358, tokenIndex358, depth358
												}
											l359:
												if !rules[ruleskip]() {
													goto l354
												}
												depth--
												add(ruleblankNodeLabel, position355)
											}
											goto l353
										l354:
											position, tokenIndex, depth = position353, tokenIndex353, depth353
											{
												position369 := position
												depth++
												if buffer[position] != rune('[') {
													goto l344
												}
												position++
											l370:
												{
													position371, tokenIndex371, depth371 := position, tokenIndex, depth
													if !rules[rulews]() {
														goto l371
													}
													goto l370
												l371:
													position, tokenIndex, depth = position371, tokenIndex371, depth371
												}
												if buffer[position] != rune(']') {
													goto l344
												}
												position++
												if !rules[ruleskip]() {
													goto l344
												}
												depth--
												add(ruleanon, position369)
											}
										}
									l353:
										depth--
										add(ruleblankNode, position352)
									}
									break
								case 'F', 'T', 'f', 't':
									if !rules[rulebooleanLiteral]() {
										goto l344
									}
									break
								case '"':
									if !rules[ruleliteral]() {
										goto l344
									}
									break
								default:
									if !rules[rulenumericLiteral]() {
										goto l344
									}
									break
								}
							}

						}
					l349:
						depth--
						add(rulegraphTerm, position348)
					}
				}
			l346:
				depth--
				add(rulevarOrTerm, position345)
			}
			return true
		l344:
			position, tokenIndex, depth = position344, tokenIndex344, depth344
			return false
		},
		/* 29 graphTerm <- <(iriref / ((&('(') nil) | (&('[' | '_') blankNode) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 30 triplesNodePath <- <(collectionPath / blankNodePropertyListPath)> */
		func() bool {
			position373, tokenIndex373, depth373 := position, tokenIndex, depth
			{
				position374 := position
				depth++
				{
					position375, tokenIndex375, depth375 := position, tokenIndex, depth
					{
						position377 := position
						depth++
						if !rules[ruleLPAREN]() {
							goto l376
						}
						if !rules[rulegraphNodePath]() {
							goto l376
						}
					l378:
						{
							position379, tokenIndex379, depth379 := position, tokenIndex, depth
							if !rules[rulegraphNodePath]() {
								goto l379
							}
							goto l378
						l379:
							position, tokenIndex, depth = position379, tokenIndex379, depth379
						}
						if !rules[ruleRPAREN]() {
							goto l376
						}
						depth--
						add(rulecollectionPath, position377)
					}
					goto l375
				l376:
					position, tokenIndex, depth = position375, tokenIndex375, depth375
					{
						position380 := position
						depth++
						{
							position381 := position
							depth++
							if buffer[position] != rune('[') {
								goto l373
							}
							position++
							if !rules[ruleskip]() {
								goto l373
							}
							depth--
							add(ruleLBRACK, position381)
						}
						if !rules[rulepropertyListPath]() {
							goto l373
						}
						{
							position382 := position
							depth++
							if buffer[position] != rune(']') {
								goto l373
							}
							position++
							if !rules[ruleskip]() {
								goto l373
							}
							depth--
							add(ruleRBRACK, position382)
						}
						depth--
						add(ruleblankNodePropertyListPath, position380)
					}
				}
			l375:
				depth--
				add(ruletriplesNodePath, position374)
			}
			return true
		l373:
			position, tokenIndex, depth = position373, tokenIndex373, depth373
			return false
		},
		/* 31 collectionPath <- <(LPAREN graphNodePath+ RPAREN)> */
		nil,
		/* 32 blankNodePropertyListPath <- <(LBRACK propertyListPath RBRACK)> */
		nil,
		/* 33 propertyListPath <- <((var / verbPath) objectListPath (SEMICOLON propertyListPath?)?)> */
		func() bool {
			position385, tokenIndex385, depth385 := position, tokenIndex, depth
			{
				position386 := position
				depth++
				{
					position387, tokenIndex387, depth387 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l388
					}
					goto l387
				l388:
					position, tokenIndex, depth = position387, tokenIndex387, depth387
					{
						position389 := position
						depth++
						if !rules[rulepath]() {
							goto l385
						}
						depth--
						add(ruleverbPath, position389)
					}
				}
			l387:
				{
					position390 := position
					depth++
					if !rules[ruleobjectPath]() {
						goto l385
					}
				l391:
					{
						position392, tokenIndex392, depth392 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l392
						}
						if !rules[ruleobjectPath]() {
							goto l392
						}
						goto l391
					l392:
						position, tokenIndex, depth = position392, tokenIndex392, depth392
					}
					depth--
					add(ruleobjectListPath, position390)
				}
				{
					position393, tokenIndex393, depth393 := position, tokenIndex, depth
					if !rules[ruleSEMICOLON]() {
						goto l393
					}
					{
						position395, tokenIndex395, depth395 := position, tokenIndex, depth
						if !rules[rulepropertyListPath]() {
							goto l395
						}
						goto l396
					l395:
						position, tokenIndex, depth = position395, tokenIndex395, depth395
					}
				l396:
					goto l394
				l393:
					position, tokenIndex, depth = position393, tokenIndex393, depth393
				}
			l394:
				depth--
				add(rulepropertyListPath, position386)
			}
			return true
		l385:
			position, tokenIndex, depth = position385, tokenIndex385, depth385
			return false
		},
		/* 34 verbPath <- <path> */
		nil,
		/* 35 path <- <pathAlternative> */
		func() bool {
			position398, tokenIndex398, depth398 := position, tokenIndex, depth
			{
				position399 := position
				depth++
				{
					position400 := position
					depth++
					if !rules[rulepathSequence]() {
						goto l398
					}
				l401:
					{
						position402, tokenIndex402, depth402 := position, tokenIndex, depth
						if !rules[rulePIPE]() {
							goto l402
						}
						if !rules[rulepathSequence]() {
							goto l402
						}
						goto l401
					l402:
						position, tokenIndex, depth = position402, tokenIndex402, depth402
					}
					depth--
					add(rulepathAlternative, position400)
				}
				depth--
				add(rulepath, position399)
			}
			return true
		l398:
			position, tokenIndex, depth = position398, tokenIndex398, depth398
			return false
		},
		/* 36 pathAlternative <- <(pathSequence (PIPE pathSequence)*)> */
		nil,
		/* 37 pathSequence <- <(pathElt (SLASH pathElt)*)> */
		func() bool {
			position404, tokenIndex404, depth404 := position, tokenIndex, depth
			{
				position405 := position
				depth++
				if !rules[rulepathElt]() {
					goto l404
				}
			l406:
				{
					position407, tokenIndex407, depth407 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l407
					}
					if !rules[rulepathElt]() {
						goto l407
					}
					goto l406
				l407:
					position, tokenIndex, depth = position407, tokenIndex407, depth407
				}
				depth--
				add(rulepathSequence, position405)
			}
			return true
		l404:
			position, tokenIndex, depth = position404, tokenIndex404, depth404
			return false
		},
		/* 38 pathElt <- <(INVERSE? pathPrimary pathMod?)> */
		func() bool {
			position408, tokenIndex408, depth408 := position, tokenIndex, depth
			{
				position409 := position
				depth++
				{
					position410, tokenIndex410, depth410 := position, tokenIndex, depth
					if !rules[ruleINVERSE]() {
						goto l410
					}
					goto l411
				l410:
					position, tokenIndex, depth = position410, tokenIndex410, depth410
				}
			l411:
				{
					position412 := position
					depth++
					{
						position413, tokenIndex413, depth413 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l414
						}
						goto l413
					l414:
						position, tokenIndex, depth = position413, tokenIndex413, depth413
						{
							switch buffer[position] {
							case '(':
								if !rules[ruleLPAREN]() {
									goto l408
								}
								if !rules[rulepath]() {
									goto l408
								}
								if !rules[ruleRPAREN]() {
									goto l408
								}
								break
							case '!':
								if !rules[ruleNOT]() {
									goto l408
								}
								{
									position416 := position
									depth++
									{
										position417, tokenIndex417, depth417 := position, tokenIndex, depth
										if !rules[rulepathOneInPropertySet]() {
											goto l418
										}
										goto l417
									l418:
										position, tokenIndex, depth = position417, tokenIndex417, depth417
										if !rules[ruleLPAREN]() {
											goto l408
										}
										{
											position419, tokenIndex419, depth419 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l419
											}
										l421:
											{
												position422, tokenIndex422, depth422 := position, tokenIndex, depth
												if !rules[rulePIPE]() {
													goto l422
												}
												if !rules[rulepathOneInPropertySet]() {
													goto l422
												}
												goto l421
											l422:
												position, tokenIndex, depth = position422, tokenIndex422, depth422
											}
											goto l420
										l419:
											position, tokenIndex, depth = position419, tokenIndex419, depth419
										}
									l420:
										if !rules[ruleRPAREN]() {
											goto l408
										}
									}
								l417:
									depth--
									add(rulepathNegatedPropertySet, position416)
								}
								break
							default:
								if !rules[ruleISA]() {
									goto l408
								}
								break
							}
						}

					}
				l413:
					depth--
					add(rulepathPrimary, position412)
				}
				{
					position423, tokenIndex423, depth423 := position, tokenIndex, depth
					{
						position425 := position
						depth++
						{
							switch buffer[position] {
							case '+':
								if !rules[rulePLUS]() {
									goto l423
								}
								break
							case '?':
								{
									position427 := position
									depth++
									if buffer[position] != rune('?') {
										goto l423
									}
									position++
									if !rules[ruleskip]() {
										goto l423
									}
									depth--
									add(ruleQUESTION, position427)
								}
								break
							default:
								if !rules[ruleSTAR]() {
									goto l423
								}
								break
							}
						}

						{
							position428, tokenIndex428, depth428 := position, tokenIndex, depth
							if !rules[ruleskip]() {
								goto l428
							}
							goto l423
						l428:
							position, tokenIndex, depth = position428, tokenIndex428, depth428
						}
						depth--
						add(rulepathMod, position425)
					}
					goto l424
				l423:
					position, tokenIndex, depth = position423, tokenIndex423, depth423
				}
			l424:
				depth--
				add(rulepathElt, position409)
			}
			return true
		l408:
			position, tokenIndex, depth = position408, tokenIndex408, depth408
			return false
		},
		/* 39 pathPrimary <- <(iriref / ((&('(') (LPAREN path RPAREN)) | (&('!') (NOT pathNegatedPropertySet)) | (&('a') ISA)))> */
		nil,
		/* 40 pathNegatedPropertySet <- <(pathOneInPropertySet / (LPAREN (pathOneInPropertySet (PIPE pathOneInPropertySet)*)? RPAREN))> */
		nil,
		/* 41 pathOneInPropertySet <- <(iriref / ISA / (INVERSE (iriref / ISA)))> */
		func() bool {
			position431, tokenIndex431, depth431 := position, tokenIndex, depth
			{
				position432 := position
				depth++
				{
					position433, tokenIndex433, depth433 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l434
					}
					goto l433
				l434:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if !rules[ruleISA]() {
						goto l435
					}
					goto l433
				l435:
					position, tokenIndex, depth = position433, tokenIndex433, depth433
					if !rules[ruleINVERSE]() {
						goto l431
					}
					{
						position436, tokenIndex436, depth436 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l437
						}
						goto l436
					l437:
						position, tokenIndex, depth = position436, tokenIndex436, depth436
						if !rules[ruleISA]() {
							goto l431
						}
					}
				l436:
				}
			l433:
				depth--
				add(rulepathOneInPropertySet, position432)
			}
			return true
		l431:
			position, tokenIndex, depth = position431, tokenIndex431, depth431
			return false
		},
		/* 42 pathMod <- <(((&('+') PLUS) | (&('?') QUESTION) | (&('*') STAR)) !skip)> */
		nil,
		/* 43 objectListPath <- <(objectPath (COMMA objectPath)*)> */
		nil,
		/* 44 objectPath <- <graphNodePath> */
		func() bool {
			position440, tokenIndex440, depth440 := position, tokenIndex, depth
			{
				position441 := position
				depth++
				if !rules[rulegraphNodePath]() {
					goto l440
				}
				depth--
				add(ruleobjectPath, position441)
			}
			return true
		l440:
			position, tokenIndex, depth = position440, tokenIndex440, depth440
			return false
		},
		/* 45 graphNodePath <- <(varOrTerm / triplesNodePath)> */
		func() bool {
			position442, tokenIndex442, depth442 := position, tokenIndex, depth
			{
				position443 := position
				depth++
				{
					position444, tokenIndex444, depth444 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l445
					}
					goto l444
				l445:
					position, tokenIndex, depth = position444, tokenIndex444, depth444
					if !rules[ruletriplesNodePath]() {
						goto l442
					}
				}
			l444:
				depth--
				add(rulegraphNodePath, position443)
			}
			return true
		l442:
			position, tokenIndex, depth = position442, tokenIndex442, depth442
			return false
		},
		/* 46 solutionModifier <- <((ORDER BY orderCondition+) / ((&('H' | 'h') (HAVING constraint)) | (&('G' | 'g') (GROUP BY groupCondition+)) | (&('L' | 'O' | 'l' | 'o') limitOffsetClauses)))?> */
		func() bool {
			{
				position447 := position
				depth++
				{
					position448, tokenIndex448, depth448 := position, tokenIndex, depth
					{
						position450, tokenIndex450, depth450 := position, tokenIndex, depth
						{
							position452 := position
							depth++
							{
								position453, tokenIndex453, depth453 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l454
								}
								position++
								goto l453
							l454:
								position, tokenIndex, depth = position453, tokenIndex453, depth453
								if buffer[position] != rune('O') {
									goto l451
								}
								position++
							}
						l453:
							{
								position455, tokenIndex455, depth455 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l456
								}
								position++
								goto l455
							l456:
								position, tokenIndex, depth = position455, tokenIndex455, depth455
								if buffer[position] != rune('R') {
									goto l451
								}
								position++
							}
						l455:
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
									goto l451
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
									goto l451
								}
								position++
							}
						l459:
							{
								position461, tokenIndex461, depth461 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l462
								}
								position++
								goto l461
							l462:
								position, tokenIndex, depth = position461, tokenIndex461, depth461
								if buffer[position] != rune('R') {
									goto l451
								}
								position++
							}
						l461:
							if !rules[ruleskip]() {
								goto l451
							}
							depth--
							add(ruleORDER, position452)
						}
						if !rules[ruleBY]() {
							goto l451
						}
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
									goto l451
								}
							}
						l466:
							depth--
							add(ruleorderCondition, position465)
						}
					l463:
						{
							position464, tokenIndex464, depth464 := position, tokenIndex, depth
							{
								position490 := position
								depth++
								{
									position491, tokenIndex491, depth491 := position, tokenIndex, depth
									{
										position493, tokenIndex493, depth493 := position, tokenIndex, depth
										{
											position495, tokenIndex495, depth495 := position, tokenIndex, depth
											{
												position497 := position
												depth++
												{
													position498, tokenIndex498, depth498 := position, tokenIndex, depth
													if buffer[position] != rune('a') {
														goto l499
													}
													position++
													goto l498
												l499:
													position, tokenIndex, depth = position498, tokenIndex498, depth498
													if buffer[position] != rune('A') {
														goto l496
													}
													position++
												}
											l498:
												{
													position500, tokenIndex500, depth500 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l501
													}
													position++
													goto l500
												l501:
													position, tokenIndex, depth = position500, tokenIndex500, depth500
													if buffer[position] != rune('S') {
														goto l496
													}
													position++
												}
											l500:
												{
													position502, tokenIndex502, depth502 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l503
													}
													position++
													goto l502
												l503:
													position, tokenIndex, depth = position502, tokenIndex502, depth502
													if buffer[position] != rune('C') {
														goto l496
													}
													position++
												}
											l502:
												if !rules[ruleskip]() {
													goto l496
												}
												depth--
												add(ruleASC, position497)
											}
											goto l495
										l496:
											position, tokenIndex, depth = position495, tokenIndex495, depth495
											{
												position504 := position
												depth++
												{
													position505, tokenIndex505, depth505 := position, tokenIndex, depth
													if buffer[position] != rune('d') {
														goto l506
													}
													position++
													goto l505
												l506:
													position, tokenIndex, depth = position505, tokenIndex505, depth505
													if buffer[position] != rune('D') {
														goto l493
													}
													position++
												}
											l505:
												{
													position507, tokenIndex507, depth507 := position, tokenIndex, depth
													if buffer[position] != rune('e') {
														goto l508
													}
													position++
													goto l507
												l508:
													position, tokenIndex, depth = position507, tokenIndex507, depth507
													if buffer[position] != rune('E') {
														goto l493
													}
													position++
												}
											l507:
												{
													position509, tokenIndex509, depth509 := position, tokenIndex, depth
													if buffer[position] != rune('s') {
														goto l510
													}
													position++
													goto l509
												l510:
													position, tokenIndex, depth = position509, tokenIndex509, depth509
													if buffer[position] != rune('S') {
														goto l493
													}
													position++
												}
											l509:
												{
													position511, tokenIndex511, depth511 := position, tokenIndex, depth
													if buffer[position] != rune('c') {
														goto l512
													}
													position++
													goto l511
												l512:
													position, tokenIndex, depth = position511, tokenIndex511, depth511
													if buffer[position] != rune('C') {
														goto l493
													}
													position++
												}
											l511:
												if !rules[ruleskip]() {
													goto l493
												}
												depth--
												add(ruleDESC, position504)
											}
										}
									l495:
										goto l494
									l493:
										position, tokenIndex, depth = position493, tokenIndex493, depth493
									}
								l494:
									if !rules[rulebrackettedExpression]() {
										goto l492
									}
									goto l491
								l492:
									position, tokenIndex, depth = position491, tokenIndex491, depth491
									if !rules[rulefunctionCall]() {
										goto l513
									}
									goto l491
								l513:
									position, tokenIndex, depth = position491, tokenIndex491, depth491
									if !rules[rulebuiltinCall]() {
										goto l514
									}
									goto l491
								l514:
									position, tokenIndex, depth = position491, tokenIndex491, depth491
									if !rules[rulevar]() {
										goto l464
									}
								}
							l491:
								depth--
								add(ruleorderCondition, position490)
							}
							goto l463
						l464:
							position, tokenIndex, depth = position464, tokenIndex464, depth464
						}
						goto l450
					l451:
						position, tokenIndex, depth = position450, tokenIndex450, depth450
						{
							switch buffer[position] {
							case 'H', 'h':
								{
									position516 := position
									depth++
									{
										position517, tokenIndex517, depth517 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l518
										}
										position++
										goto l517
									l518:
										position, tokenIndex, depth = position517, tokenIndex517, depth517
										if buffer[position] != rune('H') {
											goto l448
										}
										position++
									}
								l517:
									{
										position519, tokenIndex519, depth519 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l520
										}
										position++
										goto l519
									l520:
										position, tokenIndex, depth = position519, tokenIndex519, depth519
										if buffer[position] != rune('A') {
											goto l448
										}
										position++
									}
								l519:
									{
										position521, tokenIndex521, depth521 := position, tokenIndex, depth
										if buffer[position] != rune('v') {
											goto l522
										}
										position++
										goto l521
									l522:
										position, tokenIndex, depth = position521, tokenIndex521, depth521
										if buffer[position] != rune('V') {
											goto l448
										}
										position++
									}
								l521:
									{
										position523, tokenIndex523, depth523 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l524
										}
										position++
										goto l523
									l524:
										position, tokenIndex, depth = position523, tokenIndex523, depth523
										if buffer[position] != rune('I') {
											goto l448
										}
										position++
									}
								l523:
									{
										position525, tokenIndex525, depth525 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l526
										}
										position++
										goto l525
									l526:
										position, tokenIndex, depth = position525, tokenIndex525, depth525
										if buffer[position] != rune('N') {
											goto l448
										}
										position++
									}
								l525:
									{
										position527, tokenIndex527, depth527 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l528
										}
										position++
										goto l527
									l528:
										position, tokenIndex, depth = position527, tokenIndex527, depth527
										if buffer[position] != rune('G') {
											goto l448
										}
										position++
									}
								l527:
									if !rules[ruleskip]() {
										goto l448
									}
									depth--
									add(ruleHAVING, position516)
								}
								if !rules[ruleconstraint]() {
									goto l448
								}
								break
							case 'G', 'g':
								{
									position529 := position
									depth++
									{
										position530, tokenIndex530, depth530 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l531
										}
										position++
										goto l530
									l531:
										position, tokenIndex, depth = position530, tokenIndex530, depth530
										if buffer[position] != rune('G') {
											goto l448
										}
										position++
									}
								l530:
									{
										position532, tokenIndex532, depth532 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l533
										}
										position++
										goto l532
									l533:
										position, tokenIndex, depth = position532, tokenIndex532, depth532
										if buffer[position] != rune('R') {
											goto l448
										}
										position++
									}
								l532:
									{
										position534, tokenIndex534, depth534 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l535
										}
										position++
										goto l534
									l535:
										position, tokenIndex, depth = position534, tokenIndex534, depth534
										if buffer[position] != rune('O') {
											goto l448
										}
										position++
									}
								l534:
									{
										position536, tokenIndex536, depth536 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l537
										}
										position++
										goto l536
									l537:
										position, tokenIndex, depth = position536, tokenIndex536, depth536
										if buffer[position] != rune('U') {
											goto l448
										}
										position++
									}
								l536:
									{
										position538, tokenIndex538, depth538 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l539
										}
										position++
										goto l538
									l539:
										position, tokenIndex, depth = position538, tokenIndex538, depth538
										if buffer[position] != rune('P') {
											goto l448
										}
										position++
									}
								l538:
									if !rules[ruleskip]() {
										goto l448
									}
									depth--
									add(ruleGROUP, position529)
								}
								if !rules[ruleBY]() {
									goto l448
								}
								{
									position542 := position
									depth++
									{
										position543, tokenIndex543, depth543 := position, tokenIndex, depth
										if !rules[rulefunctionCall]() {
											goto l544
										}
										goto l543
									l544:
										position, tokenIndex, depth = position543, tokenIndex543, depth543
										{
											switch buffer[position] {
											case '$', '?':
												if !rules[rulevar]() {
													goto l448
												}
												break
											case '(':
												if !rules[ruleLPAREN]() {
													goto l448
												}
												if !rules[ruleexpression]() {
													goto l448
												}
												{
													position546, tokenIndex546, depth546 := position, tokenIndex, depth
													if !rules[ruleAS]() {
														goto l546
													}
													if !rules[rulevar]() {
														goto l546
													}
													goto l547
												l546:
													position, tokenIndex, depth = position546, tokenIndex546, depth546
												}
											l547:
												if !rules[ruleRPAREN]() {
													goto l448
												}
												break
											default:
												if !rules[rulebuiltinCall]() {
													goto l448
												}
												break
											}
										}

									}
								l543:
									depth--
									add(rulegroupCondition, position542)
								}
							l540:
								{
									position541, tokenIndex541, depth541 := position, tokenIndex, depth
									{
										position548 := position
										depth++
										{
											position549, tokenIndex549, depth549 := position, tokenIndex, depth
											if !rules[rulefunctionCall]() {
												goto l550
											}
											goto l549
										l550:
											position, tokenIndex, depth = position549, tokenIndex549, depth549
											{
												switch buffer[position] {
												case '$', '?':
													if !rules[rulevar]() {
														goto l541
													}
													break
												case '(':
													if !rules[ruleLPAREN]() {
														goto l541
													}
													if !rules[ruleexpression]() {
														goto l541
													}
													{
														position552, tokenIndex552, depth552 := position, tokenIndex, depth
														if !rules[ruleAS]() {
															goto l552
														}
														if !rules[rulevar]() {
															goto l552
														}
														goto l553
													l552:
														position, tokenIndex, depth = position552, tokenIndex552, depth552
													}
												l553:
													if !rules[ruleRPAREN]() {
														goto l541
													}
													break
												default:
													if !rules[rulebuiltinCall]() {
														goto l541
													}
													break
												}
											}

										}
									l549:
										depth--
										add(rulegroupCondition, position548)
									}
									goto l540
								l541:
									position, tokenIndex, depth = position541, tokenIndex541, depth541
								}
								break
							default:
								{
									position554 := position
									depth++
									{
										position555, tokenIndex555, depth555 := position, tokenIndex, depth
										if !rules[rulelimit]() {
											goto l556
										}
										{
											position557, tokenIndex557, depth557 := position, tokenIndex, depth
											if !rules[ruleoffset]() {
												goto l557
											}
											goto l558
										l557:
											position, tokenIndex, depth = position557, tokenIndex557, depth557
										}
									l558:
										goto l555
									l556:
										position, tokenIndex, depth = position555, tokenIndex555, depth555
										if !rules[ruleoffset]() {
											goto l448
										}
										{
											position559, tokenIndex559, depth559 := position, tokenIndex, depth
											if !rules[rulelimit]() {
												goto l559
											}
											goto l560
										l559:
											position, tokenIndex, depth = position559, tokenIndex559, depth559
										}
									l560:
									}
								l555:
									depth--
									add(rulelimitOffsetClauses, position554)
								}
								break
							}
						}

					}
				l450:
					goto l449
				l448:
					position, tokenIndex, depth = position448, tokenIndex448, depth448
				}
			l449:
				depth--
				add(rulesolutionModifier, position447)
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
			position564, tokenIndex564, depth564 := position, tokenIndex, depth
			{
				position565 := position
				depth++
				{
					position566 := position
					depth++
					{
						position567, tokenIndex567, depth567 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l568
						}
						position++
						goto l567
					l568:
						position, tokenIndex, depth = position567, tokenIndex567, depth567
						if buffer[position] != rune('L') {
							goto l564
						}
						position++
					}
				l567:
					{
						position569, tokenIndex569, depth569 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l570
						}
						position++
						goto l569
					l570:
						position, tokenIndex, depth = position569, tokenIndex569, depth569
						if buffer[position] != rune('I') {
							goto l564
						}
						position++
					}
				l569:
					{
						position571, tokenIndex571, depth571 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l572
						}
						position++
						goto l571
					l572:
						position, tokenIndex, depth = position571, tokenIndex571, depth571
						if buffer[position] != rune('M') {
							goto l564
						}
						position++
					}
				l571:
					{
						position573, tokenIndex573, depth573 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l574
						}
						position++
						goto l573
					l574:
						position, tokenIndex, depth = position573, tokenIndex573, depth573
						if buffer[position] != rune('I') {
							goto l564
						}
						position++
					}
				l573:
					{
						position575, tokenIndex575, depth575 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l576
						}
						position++
						goto l575
					l576:
						position, tokenIndex, depth = position575, tokenIndex575, depth575
						if buffer[position] != rune('T') {
							goto l564
						}
						position++
					}
				l575:
					if !rules[ruleskip]() {
						goto l564
					}
					depth--
					add(ruleLIMIT, position566)
				}
				if !rules[ruleINTEGER]() {
					goto l564
				}
				depth--
				add(rulelimit, position565)
			}
			return true
		l564:
			position, tokenIndex, depth = position564, tokenIndex564, depth564
			return false
		},
		/* 51 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position577, tokenIndex577, depth577 := position, tokenIndex, depth
			{
				position578 := position
				depth++
				{
					position579 := position
					depth++
					{
						position580, tokenIndex580, depth580 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l581
						}
						position++
						goto l580
					l581:
						position, tokenIndex, depth = position580, tokenIndex580, depth580
						if buffer[position] != rune('O') {
							goto l577
						}
						position++
					}
				l580:
					{
						position582, tokenIndex582, depth582 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l583
						}
						position++
						goto l582
					l583:
						position, tokenIndex, depth = position582, tokenIndex582, depth582
						if buffer[position] != rune('F') {
							goto l577
						}
						position++
					}
				l582:
					{
						position584, tokenIndex584, depth584 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l585
						}
						position++
						goto l584
					l585:
						position, tokenIndex, depth = position584, tokenIndex584, depth584
						if buffer[position] != rune('F') {
							goto l577
						}
						position++
					}
				l584:
					{
						position586, tokenIndex586, depth586 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l587
						}
						position++
						goto l586
					l587:
						position, tokenIndex, depth = position586, tokenIndex586, depth586
						if buffer[position] != rune('S') {
							goto l577
						}
						position++
					}
				l586:
					{
						position588, tokenIndex588, depth588 := position, tokenIndex, depth
						if buffer[position] != rune('e') {
							goto l589
						}
						position++
						goto l588
					l589:
						position, tokenIndex, depth = position588, tokenIndex588, depth588
						if buffer[position] != rune('E') {
							goto l577
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
							goto l577
						}
						position++
					}
				l590:
					if !rules[ruleskip]() {
						goto l577
					}
					depth--
					add(ruleOFFSET, position579)
				}
				if !rules[ruleINTEGER]() {
					goto l577
				}
				depth--
				add(ruleoffset, position578)
			}
			return true
		l577:
			position, tokenIndex, depth = position577, tokenIndex577, depth577
			return false
		},
		/* 52 expression <- <conditionalOrExpression> */
		func() bool {
			position592, tokenIndex592, depth592 := position, tokenIndex, depth
			{
				position593 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l592
				}
				depth--
				add(ruleexpression, position593)
			}
			return true
		l592:
			position, tokenIndex, depth = position592, tokenIndex592, depth592
			return false
		},
		/* 53 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position594, tokenIndex594, depth594 := position, tokenIndex, depth
			{
				position595 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l594
				}
				{
					position596, tokenIndex596, depth596 := position, tokenIndex, depth
					{
						position598 := position
						depth++
						if buffer[position] != rune('|') {
							goto l596
						}
						position++
						if buffer[position] != rune('|') {
							goto l596
						}
						position++
						if !rules[ruleskip]() {
							goto l596
						}
						depth--
						add(ruleOR, position598)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l596
					}
					goto l597
				l596:
					position, tokenIndex, depth = position596, tokenIndex596, depth596
				}
			l597:
				depth--
				add(ruleconditionalOrExpression, position595)
			}
			return true
		l594:
			position, tokenIndex, depth = position594, tokenIndex594, depth594
			return false
		},
		/* 54 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position599, tokenIndex599, depth599 := position, tokenIndex, depth
			{
				position600 := position
				depth++
				{
					position601 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l599
					}
					{
						position602, tokenIndex602, depth602 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position605 := position
									depth++
									{
										position606 := position
										depth++
										{
											position607, tokenIndex607, depth607 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l608
											}
											position++
											goto l607
										l608:
											position, tokenIndex, depth = position607, tokenIndex607, depth607
											if buffer[position] != rune('N') {
												goto l602
											}
											position++
										}
									l607:
										{
											position609, tokenIndex609, depth609 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l610
											}
											position++
											goto l609
										l610:
											position, tokenIndex, depth = position609, tokenIndex609, depth609
											if buffer[position] != rune('O') {
												goto l602
											}
											position++
										}
									l609:
										{
											position611, tokenIndex611, depth611 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l612
											}
											position++
											goto l611
										l612:
											position, tokenIndex, depth = position611, tokenIndex611, depth611
											if buffer[position] != rune('T') {
												goto l602
											}
											position++
										}
									l611:
										if buffer[position] != rune(' ') {
											goto l602
										}
										position++
										{
											position613, tokenIndex613, depth613 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l614
											}
											position++
											goto l613
										l614:
											position, tokenIndex, depth = position613, tokenIndex613, depth613
											if buffer[position] != rune('I') {
												goto l602
											}
											position++
										}
									l613:
										{
											position615, tokenIndex615, depth615 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l616
											}
											position++
											goto l615
										l616:
											position, tokenIndex, depth = position615, tokenIndex615, depth615
											if buffer[position] != rune('N') {
												goto l602
											}
											position++
										}
									l615:
										if !rules[ruleskip]() {
											goto l602
										}
										depth--
										add(ruleNOTIN, position606)
									}
									if !rules[ruleargList]() {
										goto l602
									}
									depth--
									add(rulenotin, position605)
								}
								break
							case 'I', 'i':
								{
									position617 := position
									depth++
									{
										position618 := position
										depth++
										{
											position619, tokenIndex619, depth619 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l620
											}
											position++
											goto l619
										l620:
											position, tokenIndex, depth = position619, tokenIndex619, depth619
											if buffer[position] != rune('I') {
												goto l602
											}
											position++
										}
									l619:
										{
											position621, tokenIndex621, depth621 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l622
											}
											position++
											goto l621
										l622:
											position, tokenIndex, depth = position621, tokenIndex621, depth621
											if buffer[position] != rune('N') {
												goto l602
											}
											position++
										}
									l621:
										if !rules[ruleskip]() {
											goto l602
										}
										depth--
										add(ruleIN, position618)
									}
									if !rules[ruleargList]() {
										goto l602
									}
									depth--
									add(rulein, position617)
								}
								break
							default:
								{
									position623, tokenIndex623, depth623 := position, tokenIndex, depth
									{
										position625 := position
										depth++
										if buffer[position] != rune('<') {
											goto l624
										}
										position++
										if !rules[ruleskip]() {
											goto l624
										}
										depth--
										add(ruleLT, position625)
									}
									goto l623
								l624:
									position, tokenIndex, depth = position623, tokenIndex623, depth623
									{
										position627 := position
										depth++
										if buffer[position] != rune('>') {
											goto l626
										}
										position++
										if buffer[position] != rune('=') {
											goto l626
										}
										position++
										if !rules[ruleskip]() {
											goto l626
										}
										depth--
										add(ruleGE, position627)
									}
									goto l623
								l626:
									position, tokenIndex, depth = position623, tokenIndex623, depth623
									{
										switch buffer[position] {
										case '>':
											{
												position629 := position
												depth++
												if buffer[position] != rune('>') {
													goto l602
												}
												position++
												if !rules[ruleskip]() {
													goto l602
												}
												depth--
												add(ruleGT, position629)
											}
											break
										case '<':
											{
												position630 := position
												depth++
												if buffer[position] != rune('<') {
													goto l602
												}
												position++
												if buffer[position] != rune('=') {
													goto l602
												}
												position++
												if !rules[ruleskip]() {
													goto l602
												}
												depth--
												add(ruleLE, position630)
											}
											break
										case '!':
											{
												position631 := position
												depth++
												if buffer[position] != rune('!') {
													goto l602
												}
												position++
												if buffer[position] != rune('=') {
													goto l602
												}
												position++
												if !rules[ruleskip]() {
													goto l602
												}
												depth--
												add(ruleNE, position631)
											}
											break
										default:
											if !rules[ruleEQ]() {
												goto l602
											}
											break
										}
									}

								}
							l623:
								if !rules[rulenumericExpression]() {
									goto l602
								}
								break
							}
						}

						goto l603
					l602:
						position, tokenIndex, depth = position602, tokenIndex602, depth602
					}
				l603:
					depth--
					add(rulevalueLogical, position601)
				}
				{
					position632, tokenIndex632, depth632 := position, tokenIndex, depth
					{
						position634 := position
						depth++
						if buffer[position] != rune('&') {
							goto l632
						}
						position++
						if buffer[position] != rune('&') {
							goto l632
						}
						position++
						if !rules[ruleskip]() {
							goto l632
						}
						depth--
						add(ruleAND, position634)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l632
					}
					goto l633
				l632:
					position, tokenIndex, depth = position632, tokenIndex632, depth632
				}
			l633:
				depth--
				add(ruleconditionalAndExpression, position600)
			}
			return true
		l599:
			position, tokenIndex, depth = position599, tokenIndex599, depth599
			return false
		},
		/* 55 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 56 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position636, tokenIndex636, depth636 := position, tokenIndex, depth
			{
				position637 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l636
				}
			l638:
				{
					position639, tokenIndex639, depth639 := position, tokenIndex, depth
					{
						position640, tokenIndex640, depth640 := position, tokenIndex, depth
						{
							position642, tokenIndex642, depth642 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l643
							}
							goto l642
						l643:
							position, tokenIndex, depth = position642, tokenIndex642, depth642
							if !rules[ruleMINUS]() {
								goto l641
							}
						}
					l642:
						if !rules[rulemultiplicativeExpression]() {
							goto l641
						}
						goto l640
					l641:
						position, tokenIndex, depth = position640, tokenIndex640, depth640
						{
							position644 := position
							depth++
							{
								position645, tokenIndex645, depth645 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l646
								}
								position++
								goto l645
							l646:
								position, tokenIndex, depth = position645, tokenIndex645, depth645
								if buffer[position] != rune('-') {
									goto l639
								}
								position++
							}
						l645:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l639
							}
							position++
						l647:
							{
								position648, tokenIndex648, depth648 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l648
								}
								position++
								goto l647
							l648:
								position, tokenIndex, depth = position648, tokenIndex648, depth648
							}
							{
								position649, tokenIndex649, depth649 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l649
								}
								position++
							l651:
								{
									position652, tokenIndex652, depth652 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l652
									}
									position++
									goto l651
								l652:
									position, tokenIndex, depth = position652, tokenIndex652, depth652
								}
								goto l650
							l649:
								position, tokenIndex, depth = position649, tokenIndex649, depth649
							}
						l650:
							if !rules[ruleskip]() {
								goto l639
							}
							depth--
							add(rulesignedNumericLiteral, position644)
						}
					}
				l640:
					goto l638
				l639:
					position, tokenIndex, depth = position639, tokenIndex639, depth639
				}
				depth--
				add(rulenumericExpression, position637)
			}
			return true
		l636:
			position, tokenIndex, depth = position636, tokenIndex636, depth636
			return false
		},
		/* 57 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position653, tokenIndex653, depth653 := position, tokenIndex, depth
			{
				position654 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l653
				}
			l655:
				{
					position656, tokenIndex656, depth656 := position, tokenIndex, depth
					{
						position657, tokenIndex657, depth657 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l658
						}
						goto l657
					l658:
						position, tokenIndex, depth = position657, tokenIndex657, depth657
						if !rules[ruleSLASH]() {
							goto l656
						}
					}
				l657:
					if !rules[ruleunaryExpression]() {
						goto l656
					}
					goto l655
				l656:
					position, tokenIndex, depth = position656, tokenIndex656, depth656
				}
				depth--
				add(rulemultiplicativeExpression, position654)
			}
			return true
		l653:
			position, tokenIndex, depth = position653, tokenIndex653, depth653
			return false
		},
		/* 58 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position659, tokenIndex659, depth659 := position, tokenIndex, depth
			{
				position660 := position
				depth++
				{
					position661, tokenIndex661, depth661 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l661
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l661
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l661
							}
							break
						}
					}

					goto l662
				l661:
					position, tokenIndex, depth = position661, tokenIndex661, depth661
				}
			l662:
				{
					position664 := position
					depth++
					{
						position665, tokenIndex665, depth665 := position, tokenIndex, depth
						if !rules[rulebuiltinCall]() {
							goto l666
						}
						goto l665
					l666:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if !rules[rulefunctionCall]() {
							goto l667
						}
						goto l665
					l667:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
						if !rules[ruleiriref]() {
							goto l668
						}
						goto l665
					l668:
						position, tokenIndex, depth = position665, tokenIndex665, depth665
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
															goto l659
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
															goto l659
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
															goto l659
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
															goto l659
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
															goto l659
														}
														position++
													}
												l682:
													if buffer[position] != rune('_') {
														goto l659
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
															goto l659
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
															goto l659
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
															goto l659
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
															goto l659
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
															goto l659
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
															goto l659
														}
														position++
													}
												l694:
													if !rules[ruleskip]() {
														goto l659
													}
													depth--
													add(ruleGROUPCONCAT, position673)
												}
												if !rules[ruleLPAREN]() {
													goto l659
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
													goto l659
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
													goto l659
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
															goto l659
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
															goto l659
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
															goto l659
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
															goto l659
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
															goto l659
														}
														position++
													}
												l729:
													if !rules[ruleskip]() {
														goto l659
													}
													depth--
													add(ruleCOUNT, position720)
												}
												if !rules[ruleLPAREN]() {
													goto l659
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
														goto l659
													}
												}
											l733:
												if !rules[ruleRPAREN]() {
													goto l659
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
																	goto l659
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
																	goto l659
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
																	goto l659
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
																	goto l659
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
																	goto l659
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
																	goto l659
																}
																position++
															}
														l764:
															if !rules[ruleskip]() {
																goto l659
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
																	goto l659
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
																	goto l659
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
																	goto l659
																}
																position++
															}
														l771:
															if !rules[ruleskip]() {
																goto l659
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
																	goto l659
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
																	goto l659
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
																	goto l659
																}
																position++
															}
														l778:
															if !rules[ruleskip]() {
																goto l659
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
												goto l659
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
												goto l659
											}
											if !rules[ruleRPAREN]() {
												goto l659
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
									goto l659
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l659
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l659
								}
								break
							case '(':
								if !rules[rulebrackettedExpression]() {
									goto l659
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l659
								}
								break
							}
						}

					}
				l665:
					depth--
					add(ruleprimaryExpression, position664)
				}
				depth--
				add(ruleunaryExpression, position660)
			}
			return true
		l659:
			position, tokenIndex, depth = position659, tokenIndex659, depth659
			return false
		},
		/* 59 primaryExpression <- <(builtinCall / functionCall / iriref / ((&('A' | 'C' | 'G' | 'M' | 'S' | 'a' | 'c' | 'g' | 'm' | 's') aggregate) | (&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('(') brackettedExpression) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
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
						position1514, tokenIndex1514, depth1514 := position, tokenIndex, depth
						if !rules[rulepnCharsU]() {
							goto l1515
						}
						goto l1514
					l1515:
						position, tokenIndex, depth = position1514, tokenIndex1514, depth1514
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1509
						}
						position++
					}
				l1514:
				l1516:
					{
						position1517, tokenIndex1517, depth1517 := position, tokenIndex, depth
						{
							position1518, tokenIndex1518, depth1518 := position, tokenIndex, depth
							if !rules[rulepnCharsU]() {
								goto l1519
							}
							goto l1518
						l1519:
							position, tokenIndex, depth = position1518, tokenIndex1518, depth1518
							{
								switch buffer[position] {
								case 'â':
									if c := buffer[position]; c < rune('‿') || c > rune('⁀') {
										goto l1517
									}
									position++
									break
								case 'Ì', 'Í':
									if c := buffer[position]; c < rune('̀') || c > rune('ͯ') {
										goto l1517
									}
									position++
									break
								case 'Â':
									if buffer[position] != rune('·') {
										goto l1517
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1517
									}
									position++
									break
								}
							}

						}
					l1518:
						goto l1516
					l1517:
						position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
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
			position1521, tokenIndex1521, depth1521 := position, tokenIndex, depth
			{
				position1522 := position
				depth++
				{
					position1523, tokenIndex1523, depth1523 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1524
					}
					goto l1523
				l1524:
					position, tokenIndex, depth = position1523, tokenIndex1523, depth1523
					{
						position1525 := position
						depth++
						{
							position1526, tokenIndex1526, depth1526 := position, tokenIndex, depth
							if !rules[rulepnPrefix]() {
								goto l1526
							}
							goto l1527
						l1526:
							position, tokenIndex, depth = position1526, tokenIndex1526, depth1526
						}
					l1527:
						if buffer[position] != rune(':') {
							goto l1521
						}
						position++
						{
							position1528 := position
							depth++
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
													goto l1521
												}
												position++
												{
													switch buffer[position] {
													case '%':
														if buffer[position] != rune('%') {
															goto l1521
														}
														position++
														break
													case '@':
														if buffer[position] != rune('@') {
															goto l1521
														}
														position++
														break
													case '#':
														if buffer[position] != rune('#') {
															goto l1521
														}
														position++
														break
													case '?':
														if buffer[position] != rune('?') {
															goto l1521
														}
														position++
														break
													case '/':
														if buffer[position] != rune('/') {
															goto l1521
														}
														position++
														break
													case '=':
														if buffer[position] != rune('=') {
															goto l1521
														}
														position++
														break
													case ';':
														if buffer[position] != rune(';') {
															goto l1521
														}
														position++
														break
													case ',':
														if buffer[position] != rune(',') {
															goto l1521
														}
														position++
														break
													case '+':
														if buffer[position] != rune('+') {
															goto l1521
														}
														position++
														break
													case '*':
														if buffer[position] != rune('*') {
															goto l1521
														}
														position++
														break
													case ')':
														if buffer[position] != rune(')') {
															goto l1521
														}
														position++
														break
													case '(':
														if buffer[position] != rune('(') {
															goto l1521
														}
														position++
														break
													case '\'':
														if buffer[position] != rune('\'') {
															goto l1521
														}
														position++
														break
													case '&':
														if buffer[position] != rune('&') {
															goto l1521
														}
														position++
														break
													case '$':
														if buffer[position] != rune('$') {
															goto l1521
														}
														position++
														break
													case '!':
														if buffer[position] != rune('!') {
															goto l1521
														}
														position++
														break
													case '-':
														if buffer[position] != rune('-') {
															goto l1521
														}
														position++
														break
													case '.':
														if buffer[position] != rune('.') {
															goto l1521
														}
														position++
														break
													case '~':
														if buffer[position] != rune('~') {
															goto l1521
														}
														position++
														break
													default:
														if buffer[position] != rune('_') {
															goto l1521
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
										goto l1521
									}
									position++
									break
								case ':':
									if buffer[position] != rune(':') {
										goto l1521
									}
									position++
									break
								default:
									if !rules[rulepnCharsU]() {
										goto l1521
									}
									break
								}
							}

						l1529:
							{
								position1530, tokenIndex1530, depth1530 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '%', '\\':
										{
											position1539 := position
											depth++
											{
												position1540, tokenIndex1540, depth1540 := position, tokenIndex, depth
												{
													position1542 := position
													depth++
													if buffer[position] != rune('%') {
														goto l1541
													}
													position++
													if !rules[rulehex]() {
														goto l1541
													}
													if !rules[rulehex]() {
														goto l1541
													}
													depth--
													add(rulepercent, position1542)
												}
												goto l1540
											l1541:
												position, tokenIndex, depth = position1540, tokenIndex1540, depth1540
												{
													position1543 := position
													depth++
													if buffer[position] != rune('\\') {
														goto l1530
													}
													position++
													{
														switch buffer[position] {
														case '%':
															if buffer[position] != rune('%') {
																goto l1530
															}
															position++
															break
														case '@':
															if buffer[position] != rune('@') {
																goto l1530
															}
															position++
															break
														case '#':
															if buffer[position] != rune('#') {
																goto l1530
															}
															position++
															break
														case '?':
															if buffer[position] != rune('?') {
																goto l1530
															}
															position++
															break
														case '/':
															if buffer[position] != rune('/') {
																goto l1530
															}
															position++
															break
														case '=':
															if buffer[position] != rune('=') {
																goto l1530
															}
															position++
															break
														case ';':
															if buffer[position] != rune(';') {
																goto l1530
															}
															position++
															break
														case ',':
															if buffer[position] != rune(',') {
																goto l1530
															}
															position++
															break
														case '+':
															if buffer[position] != rune('+') {
																goto l1530
															}
															position++
															break
														case '*':
															if buffer[position] != rune('*') {
																goto l1530
															}
															position++
															break
														case ')':
															if buffer[position] != rune(')') {
																goto l1530
															}
															position++
															break
														case '(':
															if buffer[position] != rune('(') {
																goto l1530
															}
															position++
															break
														case '\'':
															if buffer[position] != rune('\'') {
																goto l1530
															}
															position++
															break
														case '&':
															if buffer[position] != rune('&') {
																goto l1530
															}
															position++
															break
														case '$':
															if buffer[position] != rune('$') {
																goto l1530
															}
															position++
															break
														case '!':
															if buffer[position] != rune('!') {
																goto l1530
															}
															position++
															break
														case '-':
															if buffer[position] != rune('-') {
																goto l1530
															}
															position++
															break
														case '.':
															if buffer[position] != rune('.') {
																goto l1530
															}
															position++
															break
														case '~':
															if buffer[position] != rune('~') {
																goto l1530
															}
															position++
															break
														default:
															if buffer[position] != rune('_') {
																goto l1530
															}
															position++
															break
														}
													}

													depth--
													add(rulepnLocalEsc, position1543)
												}
											}
										l1540:
											depth--
											add(ruleplx, position1539)
										}
										break
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1530
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1530
										}
										position++
										break
									default:
										if !rules[rulepnCharsU]() {
											goto l1530
										}
										break
									}
								}

								goto l1529
							l1530:
								position, tokenIndex, depth = position1530, tokenIndex1530, depth1530
							}
							depth--
							add(rulepnLocal, position1528)
						}
						if !rules[ruleskip]() {
							goto l1521
						}
						depth--
						add(ruleprefixedName, position1525)
					}
				}
			l1523:
				depth--
				add(ruleiriref, position1522)
			}
			return true
		l1521:
			position, tokenIndex, depth = position1521, tokenIndex1521, depth1521
			return false
		},
		/* 71 iri <- <('<' (!'>' .)* '>' skip)> */
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
		/* 72 prefixedName <- <(pnPrefix? ':' pnLocal skip)> */
		nil,
		/* 73 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
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
		/* 74 string <- <('"' (!'"' .)* '"')> */
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
		/* 75 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
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
		/* 76 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 77 booleanLiteral <- <(TRUE / FALSE)> */
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
		/* 78 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 79 blankNodeLabel <- <('_' ':' (pnCharsU / [0-9]) (((pnCharsU / ((&('.') '.') | (&('-') '-') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))* pnCharsU) / ([0-9] / '-'))? skip)> */
		nil,
		/* 80 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 81 nil <- <('(' ws* ')' skip)> */
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
		/* 82 VARNAME <- <((pnCharsU / [0-9]) (pnCharsU / ((&('â') [‿-⁀]) | (&('Ì' | 'Í') [̀-ͯ]) | (&('Â') '·') | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9])))*)> */
		nil,
		/* 83 pnPrefix <- <(pnCharsBase pnChars*)> */
		func() bool {
			position1619, tokenIndex1619, depth1619 := position, tokenIndex, depth
			{
				position1620 := position
				depth++
				if !rules[rulepnCharsBase]() {
					goto l1619
				}
			l1621:
				{
					position1622, tokenIndex1622, depth1622 := position, tokenIndex, depth
					{
						position1623 := position
						depth++
						{
							switch buffer[position] {
							case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1622
								}
								position++
								break
							case '-':
								if buffer[position] != rune('-') {
									goto l1622
								}
								position++
								break
							default:
								if !rules[rulepnCharsU]() {
									goto l1622
								}
								break
							}
						}

						depth--
						add(rulepnChars, position1623)
					}
					goto l1621
				l1622:
					position, tokenIndex, depth = position1622, tokenIndex1622, depth1622
				}
				depth--
				add(rulepnPrefix, position1620)
			}
			return true
		l1619:
			position, tokenIndex, depth = position1619, tokenIndex1619, depth1619
			return false
		},
		/* 84 pnLocal <- <((&('%' | '\\') plx) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&(':') ':') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))+> */
		nil,
		/* 85 pnChars <- <((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('-') '-') | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z' | '_' | 'a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z' | 'Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë' | 'Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á' | 'â' | 'ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í' | 'ï' | 'ð' | 'ñ' | 'ò' | 'ó') pnCharsU))> */
		nil,
		/* 86 pnCharsU <- <(pnCharsBase / '_')> */
		func() bool {
			position1627, tokenIndex1627, depth1627 := position, tokenIndex, depth
			{
				position1628 := position
				depth++
				{
					position1629, tokenIndex1629, depth1629 := position, tokenIndex, depth
					if !rules[rulepnCharsBase]() {
						goto l1630
					}
					goto l1629
				l1630:
					position, tokenIndex, depth = position1629, tokenIndex1629, depth1629
					if buffer[position] != rune('_') {
						goto l1627
					}
					position++
				}
			l1629:
				depth--
				add(rulepnCharsU, position1628)
			}
			return true
		l1627:
			position, tokenIndex, depth = position1627, tokenIndex1627, depth1627
			return false
		},
		/* 87 pnCharsBase <- <([À-Ö] / [Ø-ö] / [Ͱ-ͽ] / [‌-‍] / [⁰-↏] / [豈-﷏] / ((&('ð' | 'ñ' | 'ò' | 'ó') [𐀀-󯿿]) | (&('ï') [ﷰ-�]) | (&('ã' | 'ä' | 'å' | 'æ' | 'ç' | 'è' | 'é' | 'ê' | 'ë' | 'ì' | 'í') [、-퟿]) | (&('â') [Ⰰ-⿯]) | (&('Í' | 'Î' | 'Ï' | 'Ð' | 'Ñ' | 'Ò' | 'Ó' | 'Ô' | 'Õ' | 'Ö' | '×' | 'Ø' | 'Ù' | 'Ú' | 'Û' | 'Ü' | 'Ý' | 'Þ' | 'ß' | 'à' | 'á') [Ϳ-῿]) | (&('Ã' | 'Ä' | 'Å' | 'Æ' | 'Ç' | 'È' | 'É' | 'Ê' | 'Ë') [ø-˿]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))> */
		func() bool {
			position1631, tokenIndex1631, depth1631 := position, tokenIndex, depth
			{
				position1632 := position
				depth++
				{
					position1633, tokenIndex1633, depth1633 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('À') || c > rune('Ö') {
						goto l1634
					}
					position++
					goto l1633
				l1634:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					if c := buffer[position]; c < rune('Ø') || c > rune('ö') {
						goto l1635
					}
					position++
					goto l1633
				l1635:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					if c := buffer[position]; c < rune('Ͱ') || c > rune('ͽ') {
						goto l1636
					}
					position++
					goto l1633
				l1636:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					if c := buffer[position]; c < rune('\u200c') || c > rune('\u200d') {
						goto l1637
					}
					position++
					goto l1633
				l1637:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					if c := buffer[position]; c < rune('⁰') || c > rune('\u218f') {
						goto l1638
					}
					position++
					goto l1633
				l1638:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					if c := buffer[position]; c < rune('豈') || c > rune('\ufdcf') {
						goto l1639
					}
					position++
					goto l1633
				l1639:
					position, tokenIndex, depth = position1633, tokenIndex1633, depth1633
					{
						switch buffer[position] {
						case 'ð', 'ñ', 'ò', 'ó':
							if c := buffer[position]; c < rune('𐀀') || c > rune('\U000effff') {
								goto l1631
							}
							position++
							break
						case 'ï':
							if c := buffer[position]; c < rune('ﷰ') || c > rune('�') {
								goto l1631
							}
							position++
							break
						case 'ã', 'ä', 'å', 'æ', 'ç', 'è', 'é', 'ê', 'ë', 'ì', 'í':
							if c := buffer[position]; c < rune('、') || c > rune('\ud7ff') {
								goto l1631
							}
							position++
							break
						case 'â':
							if c := buffer[position]; c < rune('Ⰰ') || c > rune('\u2fef') {
								goto l1631
							}
							position++
							break
						case 'Í', 'Î', 'Ï', 'Ð', 'Ñ', 'Ò', 'Ó', 'Ô', 'Õ', 'Ö', '×', 'Ø', 'Ù', 'Ú', 'Û', 'Ü', 'Ý', 'Þ', 'ß', 'à', 'á':
							if c := buffer[position]; c < rune('\u037f') || c > rune('\u1fff') {
								goto l1631
							}
							position++
							break
						case 'Ã', 'Ä', 'Å', 'Æ', 'Ç', 'È', 'É', 'Ê', 'Ë':
							if c := buffer[position]; c < rune('ø') || c > rune('˿') {
								goto l1631
							}
							position++
							break
						case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1631
							}
							position++
							break
						default:
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1631
							}
							position++
							break
						}
					}

				}
			l1633:
				depth--
				add(rulepnCharsBase, position1632)
			}
			return true
		l1631:
			position, tokenIndex, depth = position1631, tokenIndex1631, depth1631
			return false
		},
		/* 88 plx <- <(percent / pnLocalEsc)> */
		nil,
		/* 89 percent <- <('%' hex hex)> */
		nil,
		/* 90 hex <- <((&('a' | 'b' | 'c' | 'd' | 'e' | 'f') [a-f]) | (&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]))> */
		func() bool {
			position1643, tokenIndex1643, depth1643 := position, tokenIndex, depth
			{
				position1644 := position
				depth++
				{
					switch buffer[position] {
					case 'a', 'b', 'c', 'd', 'e', 'f':
						if c := buffer[position]; c < rune('a') || c > rune('f') {
							goto l1643
						}
						position++
						break
					case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1643
						}
						position++
						break
					default:
						if c := buffer[position]; c < rune('A') || c > rune('Z') {
							goto l1643
						}
						position++
						break
					}
				}

				depth--
				add(rulehex, position1644)
			}
			return true
		l1643:
			position, tokenIndex, depth = position1643, tokenIndex1643, depth1643
			return false
		},
		/* 91 pnLocalEsc <- <('\\' ((&('%') '%') | (&('@') '@') | (&('#') '#') | (&('?') '?') | (&('/') '/') | (&('=') '=') | (&(';') ';') | (&(',') ',') | (&('+') '+') | (&('*') '*') | (&(')') ')') | (&('(') '(') | (&('\'') '\'') | (&('&') '&') | (&('$') '$') | (&('!') '!') | (&('-') '-') | (&('.') '.') | (&('~') '~') | (&('_') '_')))> */
		nil,
		/* 92 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 93 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 94 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 95 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 96 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 97 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 98 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		func() bool {
			position1653, tokenIndex1653, depth1653 := position, tokenIndex, depth
			{
				position1654 := position
				depth++
				{
					position1655, tokenIndex1655, depth1655 := position, tokenIndex, depth
					if buffer[position] != rune('d') {
						goto l1656
					}
					position++
					goto l1655
				l1656:
					position, tokenIndex, depth = position1655, tokenIndex1655, depth1655
					if buffer[position] != rune('D') {
						goto l1653
					}
					position++
				}
			l1655:
				{
					position1657, tokenIndex1657, depth1657 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1658
					}
					position++
					goto l1657
				l1658:
					position, tokenIndex, depth = position1657, tokenIndex1657, depth1657
					if buffer[position] != rune('I') {
						goto l1653
					}
					position++
				}
			l1657:
				{
					position1659, tokenIndex1659, depth1659 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1660
					}
					position++
					goto l1659
				l1660:
					position, tokenIndex, depth = position1659, tokenIndex1659, depth1659
					if buffer[position] != rune('S') {
						goto l1653
					}
					position++
				}
			l1659:
				{
					position1661, tokenIndex1661, depth1661 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1662
					}
					position++
					goto l1661
				l1662:
					position, tokenIndex, depth = position1661, tokenIndex1661, depth1661
					if buffer[position] != rune('T') {
						goto l1653
					}
					position++
				}
			l1661:
				{
					position1663, tokenIndex1663, depth1663 := position, tokenIndex, depth
					if buffer[position] != rune('i') {
						goto l1664
					}
					position++
					goto l1663
				l1664:
					position, tokenIndex, depth = position1663, tokenIndex1663, depth1663
					if buffer[position] != rune('I') {
						goto l1653
					}
					position++
				}
			l1663:
				{
					position1665, tokenIndex1665, depth1665 := position, tokenIndex, depth
					if buffer[position] != rune('n') {
						goto l1666
					}
					position++
					goto l1665
				l1666:
					position, tokenIndex, depth = position1665, tokenIndex1665, depth1665
					if buffer[position] != rune('N') {
						goto l1653
					}
					position++
				}
			l1665:
				{
					position1667, tokenIndex1667, depth1667 := position, tokenIndex, depth
					if buffer[position] != rune('c') {
						goto l1668
					}
					position++
					goto l1667
				l1668:
					position, tokenIndex, depth = position1667, tokenIndex1667, depth1667
					if buffer[position] != rune('C') {
						goto l1653
					}
					position++
				}
			l1667:
				{
					position1669, tokenIndex1669, depth1669 := position, tokenIndex, depth
					if buffer[position] != rune('t') {
						goto l1670
					}
					position++
					goto l1669
				l1670:
					position, tokenIndex, depth = position1669, tokenIndex1669, depth1669
					if buffer[position] != rune('T') {
						goto l1653
					}
					position++
				}
			l1669:
				if !rules[ruleskip]() {
					goto l1653
				}
				depth--
				add(ruleDISTINCT, position1654)
			}
			return true
		l1653:
			position, tokenIndex, depth = position1653, tokenIndex1653, depth1653
			return false
		},
		/* 99 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 100 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 101 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 102 LBRACE <- <('{' skip)> */
		func() bool {
			position1674, tokenIndex1674, depth1674 := position, tokenIndex, depth
			{
				position1675 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1674
				}
				position++
				if !rules[ruleskip]() {
					goto l1674
				}
				depth--
				add(ruleLBRACE, position1675)
			}
			return true
		l1674:
			position, tokenIndex, depth = position1674, tokenIndex1674, depth1674
			return false
		},
		/* 103 RBRACE <- <('}' skip)> */
		func() bool {
			position1676, tokenIndex1676, depth1676 := position, tokenIndex, depth
			{
				position1677 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1676
				}
				position++
				if !rules[ruleskip]() {
					goto l1676
				}
				depth--
				add(ruleRBRACE, position1677)
			}
			return true
		l1676:
			position, tokenIndex, depth = position1676, tokenIndex1676, depth1676
			return false
		},
		/* 104 LBRACK <- <('[' skip)> */
		nil,
		/* 105 RBRACK <- <(']' skip)> */
		nil,
		/* 106 SEMICOLON <- <(';' skip)> */
		func() bool {
			position1680, tokenIndex1680, depth1680 := position, tokenIndex, depth
			{
				position1681 := position
				depth++
				if buffer[position] != rune(';') {
					goto l1680
				}
				position++
				if !rules[ruleskip]() {
					goto l1680
				}
				depth--
				add(ruleSEMICOLON, position1681)
			}
			return true
		l1680:
			position, tokenIndex, depth = position1680, tokenIndex1680, depth1680
			return false
		},
		/* 107 COMMA <- <(',' skip)> */
		func() bool {
			position1682, tokenIndex1682, depth1682 := position, tokenIndex, depth
			{
				position1683 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1682
				}
				position++
				if !rules[ruleskip]() {
					goto l1682
				}
				depth--
				add(ruleCOMMA, position1683)
			}
			return true
		l1682:
			position, tokenIndex, depth = position1682, tokenIndex1682, depth1682
			return false
		},
		/* 108 DOT <- <('.' skip)> */
		func() bool {
			position1684, tokenIndex1684, depth1684 := position, tokenIndex, depth
			{
				position1685 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1684
				}
				position++
				if !rules[ruleskip]() {
					goto l1684
				}
				depth--
				add(ruleDOT, position1685)
			}
			return true
		l1684:
			position, tokenIndex, depth = position1684, tokenIndex1684, depth1684
			return false
		},
		/* 109 COLON <- <(':' skip)> */
		nil,
		/* 110 PIPE <- <('|' skip)> */
		func() bool {
			position1687, tokenIndex1687, depth1687 := position, tokenIndex, depth
			{
				position1688 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1687
				}
				position++
				if !rules[ruleskip]() {
					goto l1687
				}
				depth--
				add(rulePIPE, position1688)
			}
			return true
		l1687:
			position, tokenIndex, depth = position1687, tokenIndex1687, depth1687
			return false
		},
		/* 111 SLASH <- <('/' skip)> */
		func() bool {
			position1689, tokenIndex1689, depth1689 := position, tokenIndex, depth
			{
				position1690 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1689
				}
				position++
				if !rules[ruleskip]() {
					goto l1689
				}
				depth--
				add(ruleSLASH, position1690)
			}
			return true
		l1689:
			position, tokenIndex, depth = position1689, tokenIndex1689, depth1689
			return false
		},
		/* 112 INVERSE <- <('^' skip)> */
		func() bool {
			position1691, tokenIndex1691, depth1691 := position, tokenIndex, depth
			{
				position1692 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1691
				}
				position++
				if !rules[ruleskip]() {
					goto l1691
				}
				depth--
				add(ruleINVERSE, position1692)
			}
			return true
		l1691:
			position, tokenIndex, depth = position1691, tokenIndex1691, depth1691
			return false
		},
		/* 113 LPAREN <- <('(' skip)> */
		func() bool {
			position1693, tokenIndex1693, depth1693 := position, tokenIndex, depth
			{
				position1694 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1693
				}
				position++
				if !rules[ruleskip]() {
					goto l1693
				}
				depth--
				add(ruleLPAREN, position1694)
			}
			return true
		l1693:
			position, tokenIndex, depth = position1693, tokenIndex1693, depth1693
			return false
		},
		/* 114 RPAREN <- <(')' skip)> */
		func() bool {
			position1695, tokenIndex1695, depth1695 := position, tokenIndex, depth
			{
				position1696 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1695
				}
				position++
				if !rules[ruleskip]() {
					goto l1695
				}
				depth--
				add(ruleRPAREN, position1696)
			}
			return true
		l1695:
			position, tokenIndex, depth = position1695, tokenIndex1695, depth1695
			return false
		},
		/* 115 ISA <- <('a' skip)> */
		func() bool {
			position1697, tokenIndex1697, depth1697 := position, tokenIndex, depth
			{
				position1698 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1697
				}
				position++
				if !rules[ruleskip]() {
					goto l1697
				}
				depth--
				add(ruleISA, position1698)
			}
			return true
		l1697:
			position, tokenIndex, depth = position1697, tokenIndex1697, depth1697
			return false
		},
		/* 116 NOT <- <('!' skip)> */
		func() bool {
			position1699, tokenIndex1699, depth1699 := position, tokenIndex, depth
			{
				position1700 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1699
				}
				position++
				if !rules[ruleskip]() {
					goto l1699
				}
				depth--
				add(ruleNOT, position1700)
			}
			return true
		l1699:
			position, tokenIndex, depth = position1699, tokenIndex1699, depth1699
			return false
		},
		/* 117 STAR <- <('*' skip)> */
		func() bool {
			position1701, tokenIndex1701, depth1701 := position, tokenIndex, depth
			{
				position1702 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1701
				}
				position++
				if !rules[ruleskip]() {
					goto l1701
				}
				depth--
				add(ruleSTAR, position1702)
			}
			return true
		l1701:
			position, tokenIndex, depth = position1701, tokenIndex1701, depth1701
			return false
		},
		/* 118 QUESTION <- <('?' skip)> */
		nil,
		/* 119 PLUS <- <('+' skip)> */
		func() bool {
			position1704, tokenIndex1704, depth1704 := position, tokenIndex, depth
			{
				position1705 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1704
				}
				position++
				if !rules[ruleskip]() {
					goto l1704
				}
				depth--
				add(rulePLUS, position1705)
			}
			return true
		l1704:
			position, tokenIndex, depth = position1704, tokenIndex1704, depth1704
			return false
		},
		/* 120 MINUS <- <('-' skip)> */
		func() bool {
			position1706, tokenIndex1706, depth1706 := position, tokenIndex, depth
			{
				position1707 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1706
				}
				position++
				if !rules[ruleskip]() {
					goto l1706
				}
				depth--
				add(ruleMINUS, position1707)
			}
			return true
		l1706:
			position, tokenIndex, depth = position1706, tokenIndex1706, depth1706
			return false
		},
		/* 121 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 122 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 123 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 124 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 125 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1712, tokenIndex1712, depth1712 := position, tokenIndex, depth
			{
				position1713 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1712
				}
				position++
			l1714:
				{
					position1715, tokenIndex1715, depth1715 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1715
					}
					position++
					goto l1714
				l1715:
					position, tokenIndex, depth = position1715, tokenIndex1715, depth1715
				}
				if !rules[ruleskip]() {
					goto l1712
				}
				depth--
				add(ruleINTEGER, position1713)
			}
			return true
		l1712:
			position, tokenIndex, depth = position1712, tokenIndex1712, depth1712
			return false
		},
		/* 126 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 127 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 128 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 129 OR <- <('|' '|' skip)> */
		nil,
		/* 130 AND <- <('&' '&' skip)> */
		nil,
		/* 131 EQ <- <('=' skip)> */
		func() bool {
			position1721, tokenIndex1721, depth1721 := position, tokenIndex, depth
			{
				position1722 := position
				depth++
				if buffer[position] != rune('=') {
					goto l1721
				}
				position++
				if !rules[ruleskip]() {
					goto l1721
				}
				depth--
				add(ruleEQ, position1722)
			}
			return true
		l1721:
			position, tokenIndex, depth = position1721, tokenIndex1721, depth1721
			return false
		},
		/* 132 NE <- <('!' '=' skip)> */
		nil,
		/* 133 GT <- <('>' skip)> */
		nil,
		/* 134 LT <- <('<' skip)> */
		nil,
		/* 135 LE <- <('<' '=' skip)> */
		nil,
		/* 136 GE <- <('>' '=' skip)> */
		nil,
		/* 137 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 138 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 139 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1730, tokenIndex1730, depth1730 := position, tokenIndex, depth
			{
				position1731 := position
				depth++
				{
					position1732, tokenIndex1732, depth1732 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1733
					}
					position++
					goto l1732
				l1733:
					position, tokenIndex, depth = position1732, tokenIndex1732, depth1732
					if buffer[position] != rune('A') {
						goto l1730
					}
					position++
				}
			l1732:
				{
					position1734, tokenIndex1734, depth1734 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1735
					}
					position++
					goto l1734
				l1735:
					position, tokenIndex, depth = position1734, tokenIndex1734, depth1734
					if buffer[position] != rune('S') {
						goto l1730
					}
					position++
				}
			l1734:
				if !rules[ruleskip]() {
					goto l1730
				}
				depth--
				add(ruleAS, position1731)
			}
			return true
		l1730:
			position, tokenIndex, depth = position1730, tokenIndex1730, depth1730
			return false
		},
		/* 140 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 141 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 142 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 143 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 144 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 145 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 146 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 147 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 148 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 149 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 150 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 151 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 152 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 153 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 154 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 155 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 156 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 157 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 158 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 159 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 160 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 161 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 162 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 163 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 164 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 165 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 166 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 167 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 168 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 169 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 170 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 171 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 172 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 173 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 174 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 175 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 176 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 177 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 178 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 179 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 180 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 181 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 182 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 183 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 184 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 185 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 186 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 187 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 188 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 189 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 190 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 191 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 192 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 193 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 194 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 195 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 196 SUM <- <(('s' / 'S') ('u' / 'U') ('m' / 'M') skip)> */
		nil,
		/* 197 MIN <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 198 MAX <- <(('m' / 'M') ('a' / 'A') ('x' / 'X') skip)> */
		nil,
		/* 199 AVG <- <(('a' / 'A') ('v' / 'V') ('g' / 'G') skip)> */
		nil,
		/* 200 SAMPLE <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('p' / 'P') ('l' / 'L') ('e' / 'E') skip)> */
		nil,
		/* 201 COUNT <- <(('c' / 'C') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('t' / 'T') skip)> */
		nil,
		/* 202 GROUPCONCAT <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') '_' ('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 203 SEPARATOR <- <(('s' / 'S') ('e' / 'E') ('p' / 'P') ('a' / 'A') ('r' / 'R') ('a' / 'A') ('t' / 'T') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 204 ASC <- <(('a' / 'A') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 205 DESC <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') skip)> */
		nil,
		/* 206 ORDER <- <(('o' / 'O') ('r' / 'R') ('d' / 'D') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 207 GROUP <- <(('g' / 'G') ('r' / 'R') ('o' / 'O') ('u' / 'U') ('p' / 'P') skip)> */
		nil,
		/* 208 BY <- <(('b' / 'B') ('y' / 'Y') skip)> */
		func() bool {
			position1804, tokenIndex1804, depth1804 := position, tokenIndex, depth
			{
				position1805 := position
				depth++
				{
					position1806, tokenIndex1806, depth1806 := position, tokenIndex, depth
					if buffer[position] != rune('b') {
						goto l1807
					}
					position++
					goto l1806
				l1807:
					position, tokenIndex, depth = position1806, tokenIndex1806, depth1806
					if buffer[position] != rune('B') {
						goto l1804
					}
					position++
				}
			l1806:
				{
					position1808, tokenIndex1808, depth1808 := position, tokenIndex, depth
					if buffer[position] != rune('y') {
						goto l1809
					}
					position++
					goto l1808
				l1809:
					position, tokenIndex, depth = position1808, tokenIndex1808, depth1808
					if buffer[position] != rune('Y') {
						goto l1804
					}
					position++
				}
			l1808:
				if !rules[ruleskip]() {
					goto l1804
				}
				depth--
				add(ruleBY, position1805)
			}
			return true
		l1804:
			position, tokenIndex, depth = position1804, tokenIndex1804, depth1804
			return false
		},
		/* 209 HAVING <- <(('h' / 'H') ('a' / 'A') ('v' / 'V') ('i' / 'I') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 210 GRAPH <- <(('g' / 'G') ('r' / 'R') ('a' / 'A') ('p' / 'P') ('h' / 'H') skip)> */
		nil,
		/* 211 MINUSSETOPER <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('s' / 'S') skip)> */
		nil,
		/* 212 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1814 := position
				depth++
			l1815:
				{
					position1816, tokenIndex1816, depth1816 := position, tokenIndex, depth
					{
						position1817, tokenIndex1817, depth1817 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1818
						}
						goto l1817
					l1818:
						position, tokenIndex, depth = position1817, tokenIndex1817, depth1817
						{
							position1819 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1816
							}
							position++
						l1820:
							{
								position1821, tokenIndex1821, depth1821 := position, tokenIndex, depth
								{
									position1822, tokenIndex1822, depth1822 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1822
									}
									goto l1821
								l1822:
									position, tokenIndex, depth = position1822, tokenIndex1822, depth1822
								}
								if !matchDot() {
									goto l1821
								}
								goto l1820
							l1821:
								position, tokenIndex, depth = position1821, tokenIndex1821, depth1821
							}
							if !rules[ruleendOfLine]() {
								goto l1816
							}
							depth--
							add(rulecomment, position1819)
						}
					}
				l1817:
					goto l1815
				l1816:
					position, tokenIndex, depth = position1816, tokenIndex1816, depth1816
				}
				depth--
				add(ruleskip, position1814)
			}
			return true
		},
		/* 213 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1823, tokenIndex1823, depth1823 := position, tokenIndex, depth
			{
				position1824 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1823
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1823
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1823
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1823
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1823
						}
						break
					}
				}

				depth--
				add(rulews, position1824)
			}
			return true
		l1823:
			position, tokenIndex, depth = position1823, tokenIndex1823, depth1823
			return false
		},
		/* 214 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 215 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1827, tokenIndex1827, depth1827 := position, tokenIndex, depth
			{
				position1828 := position
				depth++
				{
					position1829, tokenIndex1829, depth1829 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1830
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1830
					}
					position++
					goto l1829
				l1830:
					position, tokenIndex, depth = position1829, tokenIndex1829, depth1829
					if buffer[position] != rune('\n') {
						goto l1831
					}
					position++
					goto l1829
				l1831:
					position, tokenIndex, depth = position1829, tokenIndex1829, depth1829
					if buffer[position] != rune('\r') {
						goto l1827
					}
					position++
				}
			l1829:
				depth--
				add(ruleendOfLine, position1828)
			}
			return true
		l1827:
			position, tokenIndex, depth = position1827, tokenIndex1827, depth1827
			return false
		},
	}
	p.rules = rules
}
