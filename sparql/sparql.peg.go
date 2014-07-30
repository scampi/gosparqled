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
	rules  [185]func() bool
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
						{
							position128 := position
							depth++
							{
								position129, tokenIndex129, depth129 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l130
								}
								position++
								goto l129
							l130:
								position, tokenIndex, depth = position129, tokenIndex129, depth129
								if buffer[position] != rune('D') {
									goto l127
								}
								position++
							}
						l129:
							{
								position131, tokenIndex131, depth131 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l132
								}
								position++
								goto l131
							l132:
								position, tokenIndex, depth = position131, tokenIndex131, depth131
								if buffer[position] != rune('I') {
									goto l127
								}
								position++
							}
						l131:
							{
								position133, tokenIndex133, depth133 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l134
								}
								position++
								goto l133
							l134:
								position, tokenIndex, depth = position133, tokenIndex133, depth133
								if buffer[position] != rune('S') {
									goto l127
								}
								position++
							}
						l133:
							{
								position135, tokenIndex135, depth135 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l136
								}
								position++
								goto l135
							l136:
								position, tokenIndex, depth = position135, tokenIndex135, depth135
								if buffer[position] != rune('T') {
									goto l127
								}
								position++
							}
						l135:
							{
								position137, tokenIndex137, depth137 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l138
								}
								position++
								goto l137
							l138:
								position, tokenIndex, depth = position137, tokenIndex137, depth137
								if buffer[position] != rune('I') {
									goto l127
								}
								position++
							}
						l137:
							{
								position139, tokenIndex139, depth139 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l140
								}
								position++
								goto l139
							l140:
								position, tokenIndex, depth = position139, tokenIndex139, depth139
								if buffer[position] != rune('N') {
									goto l127
								}
								position++
							}
						l139:
							{
								position141, tokenIndex141, depth141 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l142
								}
								position++
								goto l141
							l142:
								position, tokenIndex, depth = position141, tokenIndex141, depth141
								if buffer[position] != rune('C') {
									goto l127
								}
								position++
							}
						l141:
							{
								position143, tokenIndex143, depth143 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l144
								}
								position++
								goto l143
							l144:
								position, tokenIndex, depth = position143, tokenIndex143, depth143
								if buffer[position] != rune('T') {
									goto l127
								}
								position++
							}
						l143:
							if !rules[ruleskip]() {
								goto l127
							}
							depth--
							add(ruleDISTINCT, position128)
						}
						goto l126
					l127:
						position, tokenIndex, depth = position126, tokenIndex126, depth126
						{
							position145 := position
							depth++
							{
								position146, tokenIndex146, depth146 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l147
								}
								position++
								goto l146
							l147:
								position, tokenIndex, depth = position146, tokenIndex146, depth146
								if buffer[position] != rune('R') {
									goto l124
								}
								position++
							}
						l146:
							{
								position148, tokenIndex148, depth148 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l149
								}
								position++
								goto l148
							l149:
								position, tokenIndex, depth = position148, tokenIndex148, depth148
								if buffer[position] != rune('E') {
									goto l124
								}
								position++
							}
						l148:
							{
								position150, tokenIndex150, depth150 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l151
								}
								position++
								goto l150
							l151:
								position, tokenIndex, depth = position150, tokenIndex150, depth150
								if buffer[position] != rune('D') {
									goto l124
								}
								position++
							}
						l150:
							{
								position152, tokenIndex152, depth152 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l153
								}
								position++
								goto l152
							l153:
								position, tokenIndex, depth = position152, tokenIndex152, depth152
								if buffer[position] != rune('U') {
									goto l124
								}
								position++
							}
						l152:
							{
								position154, tokenIndex154, depth154 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l155
								}
								position++
								goto l154
							l155:
								position, tokenIndex, depth = position154, tokenIndex154, depth154
								if buffer[position] != rune('C') {
									goto l124
								}
								position++
							}
						l154:
							{
								position156, tokenIndex156, depth156 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l157
								}
								position++
								goto l156
							l157:
								position, tokenIndex, depth = position156, tokenIndex156, depth156
								if buffer[position] != rune('E') {
									goto l124
								}
								position++
							}
						l156:
							{
								position158, tokenIndex158, depth158 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l159
								}
								position++
								goto l158
							l159:
								position, tokenIndex, depth = position158, tokenIndex158, depth158
								if buffer[position] != rune('D') {
									goto l124
								}
								position++
							}
						l158:
							if !rules[ruleskip]() {
								goto l124
							}
							depth--
							add(ruleREDUCED, position145)
						}
					}
				l126:
					goto l125
				l124:
					position, tokenIndex, depth = position124, tokenIndex124, depth124
				}
			l125:
				{
					position160, tokenIndex160, depth160 := position, tokenIndex, depth
					if !rules[ruleSTAR]() {
						goto l161
					}
					goto l160
				l161:
					position, tokenIndex, depth = position160, tokenIndex160, depth160
					{
						position164 := position
						depth++
						{
							position165, tokenIndex165, depth165 := position, tokenIndex, depth
							if !rules[rulevar]() {
								goto l166
							}
							goto l165
						l166:
							position, tokenIndex, depth = position165, tokenIndex165, depth165
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
					l165:
						depth--
						add(ruleprojectionElem, position164)
					}
				l162:
					{
						position163, tokenIndex163, depth163 := position, tokenIndex, depth
						{
							position167 := position
							depth++
							{
								position168, tokenIndex168, depth168 := position, tokenIndex, depth
								if !rules[rulevar]() {
									goto l169
								}
								goto l168
							l169:
								position, tokenIndex, depth = position168, tokenIndex168, depth168
								if !rules[ruleLPAREN]() {
									goto l163
								}
								if !rules[ruleexpression]() {
									goto l163
								}
								if !rules[ruleAS]() {
									goto l163
								}
								if !rules[rulevar]() {
									goto l163
								}
								if !rules[ruleRPAREN]() {
									goto l163
								}
							}
						l168:
							depth--
							add(ruleprojectionElem, position167)
						}
						goto l162
					l163:
						position, tokenIndex, depth = position163, tokenIndex163, depth163
					}
				}
			l160:
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
			position170, tokenIndex170, depth170 := position, tokenIndex, depth
			{
				position171 := position
				depth++
				if !rules[ruleselect]() {
					goto l170
				}
				if !rules[rulewhereClause]() {
					goto l170
				}
				depth--
				add(rulesubSelect, position171)
			}
			return true
		l170:
			position, tokenIndex, depth = position170, tokenIndex170, depth170
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
			position178, tokenIndex178, depth178 := position, tokenIndex, depth
			{
				position179 := position
				depth++
				{
					position180 := position
					depth++
					{
						position181, tokenIndex181, depth181 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l182
						}
						position++
						goto l181
					l182:
						position, tokenIndex, depth = position181, tokenIndex181, depth181
						if buffer[position] != rune('F') {
							goto l178
						}
						position++
					}
				l181:
					{
						position183, tokenIndex183, depth183 := position, tokenIndex, depth
						if buffer[position] != rune('r') {
							goto l184
						}
						position++
						goto l183
					l184:
						position, tokenIndex, depth = position183, tokenIndex183, depth183
						if buffer[position] != rune('R') {
							goto l178
						}
						position++
					}
				l183:
					{
						position185, tokenIndex185, depth185 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l186
						}
						position++
						goto l185
					l186:
						position, tokenIndex, depth = position185, tokenIndex185, depth185
						if buffer[position] != rune('O') {
							goto l178
						}
						position++
					}
				l185:
					{
						position187, tokenIndex187, depth187 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l188
						}
						position++
						goto l187
					l188:
						position, tokenIndex, depth = position187, tokenIndex187, depth187
						if buffer[position] != rune('M') {
							goto l178
						}
						position++
					}
				l187:
					if !rules[ruleskip]() {
						goto l178
					}
					depth--
					add(ruleFROM, position180)
				}
				{
					position189, tokenIndex189, depth189 := position, tokenIndex, depth
					{
						position191 := position
						depth++
						{
							position192, tokenIndex192, depth192 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l193
							}
							position++
							goto l192
						l193:
							position, tokenIndex, depth = position192, tokenIndex192, depth192
							if buffer[position] != rune('N') {
								goto l189
							}
							position++
						}
					l192:
						{
							position194, tokenIndex194, depth194 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l195
							}
							position++
							goto l194
						l195:
							position, tokenIndex, depth = position194, tokenIndex194, depth194
							if buffer[position] != rune('A') {
								goto l189
							}
							position++
						}
					l194:
						{
							position196, tokenIndex196, depth196 := position, tokenIndex, depth
							if buffer[position] != rune('m') {
								goto l197
							}
							position++
							goto l196
						l197:
							position, tokenIndex, depth = position196, tokenIndex196, depth196
							if buffer[position] != rune('M') {
								goto l189
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
								goto l189
							}
							position++
						}
					l198:
						{
							position200, tokenIndex200, depth200 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l201
							}
							position++
							goto l200
						l201:
							position, tokenIndex, depth = position200, tokenIndex200, depth200
							if buffer[position] != rune('D') {
								goto l189
							}
							position++
						}
					l200:
						if !rules[ruleskip]() {
							goto l189
						}
						depth--
						add(ruleNAMED, position191)
					}
					goto l190
				l189:
					position, tokenIndex, depth = position189, tokenIndex189, depth189
				}
			l190:
				if !rules[ruleiriref]() {
					goto l178
				}
				depth--
				add(ruledatasetClause, position179)
			}
			return true
		l178:
			position, tokenIndex, depth = position178, tokenIndex178, depth178
			return false
		},
		/* 15 whereClause <- <(WHERE? groupGraphPattern)> */
		func() bool {
			position202, tokenIndex202, depth202 := position, tokenIndex, depth
			{
				position203 := position
				depth++
				{
					position204, tokenIndex204, depth204 := position, tokenIndex, depth
					{
						position206 := position
						depth++
						{
							position207, tokenIndex207, depth207 := position, tokenIndex, depth
							if buffer[position] != rune('w') {
								goto l208
							}
							position++
							goto l207
						l208:
							position, tokenIndex, depth = position207, tokenIndex207, depth207
							if buffer[position] != rune('W') {
								goto l204
							}
							position++
						}
					l207:
						{
							position209, tokenIndex209, depth209 := position, tokenIndex, depth
							if buffer[position] != rune('h') {
								goto l210
							}
							position++
							goto l209
						l210:
							position, tokenIndex, depth = position209, tokenIndex209, depth209
							if buffer[position] != rune('H') {
								goto l204
							}
							position++
						}
					l209:
						{
							position211, tokenIndex211, depth211 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l212
							}
							position++
							goto l211
						l212:
							position, tokenIndex, depth = position211, tokenIndex211, depth211
							if buffer[position] != rune('E') {
								goto l204
							}
							position++
						}
					l211:
						{
							position213, tokenIndex213, depth213 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l214
							}
							position++
							goto l213
						l214:
							position, tokenIndex, depth = position213, tokenIndex213, depth213
							if buffer[position] != rune('R') {
								goto l204
							}
							position++
						}
					l213:
						{
							position215, tokenIndex215, depth215 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l216
							}
							position++
							goto l215
						l216:
							position, tokenIndex, depth = position215, tokenIndex215, depth215
							if buffer[position] != rune('E') {
								goto l204
							}
							position++
						}
					l215:
						if !rules[ruleskip]() {
							goto l204
						}
						depth--
						add(ruleWHERE, position206)
					}
					goto l205
				l204:
					position, tokenIndex, depth = position204, tokenIndex204, depth204
				}
			l205:
				if !rules[rulegroupGraphPattern]() {
					goto l202
				}
				depth--
				add(rulewhereClause, position203)
			}
			return true
		l202:
			position, tokenIndex, depth = position202, tokenIndex202, depth202
			return false
		},
		/* 16 groupGraphPattern <- <(LBRACE (subSelect / graphPattern) RBRACE)> */
		func() bool {
			position217, tokenIndex217, depth217 := position, tokenIndex, depth
			{
				position218 := position
				depth++
				if !rules[ruleLBRACE]() {
					goto l217
				}
				{
					position219, tokenIndex219, depth219 := position, tokenIndex, depth
					if !rules[rulesubSelect]() {
						goto l220
					}
					goto l219
				l220:
					position, tokenIndex, depth = position219, tokenIndex219, depth219
					if !rules[rulegraphPattern]() {
						goto l217
					}
				}
			l219:
				if !rules[ruleRBRACE]() {
					goto l217
				}
				depth--
				add(rulegroupGraphPattern, position218)
			}
			return true
		l217:
			position, tokenIndex, depth = position217, tokenIndex217, depth217
			return false
		},
		/* 17 graphPattern <- <(basicGraphPattern? (graphPatternNotTriples DOT? graphPattern)?)> */
		func() bool {
			{
				position222 := position
				depth++
				{
					position223, tokenIndex223, depth223 := position, tokenIndex, depth
					{
						position225 := position
						depth++
						{
							position226, tokenIndex226, depth226 := position, tokenIndex, depth
							if !rules[ruletriplesBlock]() {
								goto l227
							}
						l228:
							{
								position229, tokenIndex229, depth229 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l229
								}
								{
									position230, tokenIndex230, depth230 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l230
									}
									goto l231
								l230:
									position, tokenIndex, depth = position230, tokenIndex230, depth230
								}
							l231:
								{
									position232, tokenIndex232, depth232 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l232
									}
									goto l233
								l232:
									position, tokenIndex, depth = position232, tokenIndex232, depth232
								}
							l233:
								goto l228
							l229:
								position, tokenIndex, depth = position229, tokenIndex229, depth229
							}
							goto l226
						l227:
							position, tokenIndex, depth = position226, tokenIndex226, depth226
							if !rules[rulefilterOrBind]() {
								goto l223
							}
							{
								position236, tokenIndex236, depth236 := position, tokenIndex, depth
								if !rules[ruleDOT]() {
									goto l236
								}
								goto l237
							l236:
								position, tokenIndex, depth = position236, tokenIndex236, depth236
							}
						l237:
							{
								position238, tokenIndex238, depth238 := position, tokenIndex, depth
								if !rules[ruletriplesBlock]() {
									goto l238
								}
								goto l239
							l238:
								position, tokenIndex, depth = position238, tokenIndex238, depth238
							}
						l239:
						l234:
							{
								position235, tokenIndex235, depth235 := position, tokenIndex, depth
								if !rules[rulefilterOrBind]() {
									goto l235
								}
								{
									position240, tokenIndex240, depth240 := position, tokenIndex, depth
									if !rules[ruleDOT]() {
										goto l240
									}
									goto l241
								l240:
									position, tokenIndex, depth = position240, tokenIndex240, depth240
								}
							l241:
								{
									position242, tokenIndex242, depth242 := position, tokenIndex, depth
									if !rules[ruletriplesBlock]() {
										goto l242
									}
									goto l243
								l242:
									position, tokenIndex, depth = position242, tokenIndex242, depth242
								}
							l243:
								goto l234
							l235:
								position, tokenIndex, depth = position235, tokenIndex235, depth235
							}
						}
					l226:
						depth--
						add(rulebasicGraphPattern, position225)
					}
					goto l224
				l223:
					position, tokenIndex, depth = position223, tokenIndex223, depth223
				}
			l224:
				{
					position244, tokenIndex244, depth244 := position, tokenIndex, depth
					{
						position246 := position
						depth++
						{
							position247, tokenIndex247, depth247 := position, tokenIndex, depth
							{
								position249 := position
								depth++
								{
									position250 := position
									depth++
									{
										position251, tokenIndex251, depth251 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l252
										}
										position++
										goto l251
									l252:
										position, tokenIndex, depth = position251, tokenIndex251, depth251
										if buffer[position] != rune('O') {
											goto l248
										}
										position++
									}
								l251:
									{
										position253, tokenIndex253, depth253 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l254
										}
										position++
										goto l253
									l254:
										position, tokenIndex, depth = position253, tokenIndex253, depth253
										if buffer[position] != rune('P') {
											goto l248
										}
										position++
									}
								l253:
									{
										position255, tokenIndex255, depth255 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l256
										}
										position++
										goto l255
									l256:
										position, tokenIndex, depth = position255, tokenIndex255, depth255
										if buffer[position] != rune('T') {
											goto l248
										}
										position++
									}
								l255:
									{
										position257, tokenIndex257, depth257 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l258
										}
										position++
										goto l257
									l258:
										position, tokenIndex, depth = position257, tokenIndex257, depth257
										if buffer[position] != rune('I') {
											goto l248
										}
										position++
									}
								l257:
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
											goto l248
										}
										position++
									}
								l259:
									{
										position261, tokenIndex261, depth261 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l262
										}
										position++
										goto l261
									l262:
										position, tokenIndex, depth = position261, tokenIndex261, depth261
										if buffer[position] != rune('N') {
											goto l248
										}
										position++
									}
								l261:
									{
										position263, tokenIndex263, depth263 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l264
										}
										position++
										goto l263
									l264:
										position, tokenIndex, depth = position263, tokenIndex263, depth263
										if buffer[position] != rune('A') {
											goto l248
										}
										position++
									}
								l263:
									{
										position265, tokenIndex265, depth265 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l266
										}
										position++
										goto l265
									l266:
										position, tokenIndex, depth = position265, tokenIndex265, depth265
										if buffer[position] != rune('L') {
											goto l248
										}
										position++
									}
								l265:
									if !rules[ruleskip]() {
										goto l248
									}
									depth--
									add(ruleOPTIONAL, position250)
								}
								if !rules[ruleLBRACE]() {
									goto l248
								}
								{
									position267, tokenIndex267, depth267 := position, tokenIndex, depth
									if !rules[rulesubSelect]() {
										goto l268
									}
									goto l267
								l268:
									position, tokenIndex, depth = position267, tokenIndex267, depth267
									if !rules[rulegraphPattern]() {
										goto l248
									}
								}
							l267:
								if !rules[ruleRBRACE]() {
									goto l248
								}
								depth--
								add(ruleoptionalGraphPattern, position249)
							}
							goto l247
						l248:
							position, tokenIndex, depth = position247, tokenIndex247, depth247
							if !rules[rulegroupOrUnionGraphPattern]() {
								goto l244
							}
						}
					l247:
						depth--
						add(rulegraphPatternNotTriples, position246)
					}
					{
						position269, tokenIndex269, depth269 := position, tokenIndex, depth
						if !rules[ruleDOT]() {
							goto l269
						}
						goto l270
					l269:
						position, tokenIndex, depth = position269, tokenIndex269, depth269
					}
				l270:
					if !rules[rulegraphPattern]() {
						goto l244
					}
					goto l245
				l244:
					position, tokenIndex, depth = position244, tokenIndex244, depth244
				}
			l245:
				depth--
				add(rulegraphPattern, position222)
			}
			return true
		},
		/* 18 graphPatternNotTriples <- <(optionalGraphPattern / groupOrUnionGraphPattern)> */
		nil,
		/* 19 optionalGraphPattern <- <(OPTIONAL LBRACE (subSelect / graphPattern) RBRACE)> */
		nil,
		/* 20 groupOrUnionGraphPattern <- <(groupGraphPattern (UNION groupOrUnionGraphPattern)?)> */
		func() bool {
			position273, tokenIndex273, depth273 := position, tokenIndex, depth
			{
				position274 := position
				depth++
				if !rules[rulegroupGraphPattern]() {
					goto l273
				}
				{
					position275, tokenIndex275, depth275 := position, tokenIndex, depth
					{
						position277 := position
						depth++
						{
							position278, tokenIndex278, depth278 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l279
							}
							position++
							goto l278
						l279:
							position, tokenIndex, depth = position278, tokenIndex278, depth278
							if buffer[position] != rune('U') {
								goto l275
							}
							position++
						}
					l278:
						{
							position280, tokenIndex280, depth280 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l281
							}
							position++
							goto l280
						l281:
							position, tokenIndex, depth = position280, tokenIndex280, depth280
							if buffer[position] != rune('N') {
								goto l275
							}
							position++
						}
					l280:
						{
							position282, tokenIndex282, depth282 := position, tokenIndex, depth
							if buffer[position] != rune('i') {
								goto l283
							}
							position++
							goto l282
						l283:
							position, tokenIndex, depth = position282, tokenIndex282, depth282
							if buffer[position] != rune('I') {
								goto l275
							}
							position++
						}
					l282:
						{
							position284, tokenIndex284, depth284 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l285
							}
							position++
							goto l284
						l285:
							position, tokenIndex, depth = position284, tokenIndex284, depth284
							if buffer[position] != rune('O') {
								goto l275
							}
							position++
						}
					l284:
						{
							position286, tokenIndex286, depth286 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l287
							}
							position++
							goto l286
						l287:
							position, tokenIndex, depth = position286, tokenIndex286, depth286
							if buffer[position] != rune('N') {
								goto l275
							}
							position++
						}
					l286:
						if !rules[ruleskip]() {
							goto l275
						}
						depth--
						add(ruleUNION, position277)
					}
					if !rules[rulegroupOrUnionGraphPattern]() {
						goto l275
					}
					goto l276
				l275:
					position, tokenIndex, depth = position275, tokenIndex275, depth275
				}
			l276:
				depth--
				add(rulegroupOrUnionGraphPattern, position274)
			}
			return true
		l273:
			position, tokenIndex, depth = position273, tokenIndex273, depth273
			return false
		},
		/* 21 basicGraphPattern <- <((triplesBlock (filterOrBind DOT? triplesBlock?)*) / (filterOrBind DOT? triplesBlock?)+)> */
		nil,
		/* 22 filterOrBind <- <((FILTER constraint) / (BIND LPAREN expression AS var RPAREN))> */
		func() bool {
			position289, tokenIndex289, depth289 := position, tokenIndex, depth
			{
				position290 := position
				depth++
				{
					position291, tokenIndex291, depth291 := position, tokenIndex, depth
					{
						position293 := position
						depth++
						{
							position294, tokenIndex294, depth294 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l295
							}
							position++
							goto l294
						l295:
							position, tokenIndex, depth = position294, tokenIndex294, depth294
							if buffer[position] != rune('F') {
								goto l292
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
								goto l292
							}
							position++
						}
					l296:
						{
							position298, tokenIndex298, depth298 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l299
							}
							position++
							goto l298
						l299:
							position, tokenIndex, depth = position298, tokenIndex298, depth298
							if buffer[position] != rune('L') {
								goto l292
							}
							position++
						}
					l298:
						{
							position300, tokenIndex300, depth300 := position, tokenIndex, depth
							if buffer[position] != rune('t') {
								goto l301
							}
							position++
							goto l300
						l301:
							position, tokenIndex, depth = position300, tokenIndex300, depth300
							if buffer[position] != rune('T') {
								goto l292
							}
							position++
						}
					l300:
						{
							position302, tokenIndex302, depth302 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l303
							}
							position++
							goto l302
						l303:
							position, tokenIndex, depth = position302, tokenIndex302, depth302
							if buffer[position] != rune('E') {
								goto l292
							}
							position++
						}
					l302:
						{
							position304, tokenIndex304, depth304 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l305
							}
							position++
							goto l304
						l305:
							position, tokenIndex, depth = position304, tokenIndex304, depth304
							if buffer[position] != rune('R') {
								goto l292
							}
							position++
						}
					l304:
						if !rules[ruleskip]() {
							goto l292
						}
						depth--
						add(ruleFILTER, position293)
					}
					{
						position306 := position
						depth++
						{
							position307, tokenIndex307, depth307 := position, tokenIndex, depth
							if !rules[rulebrackettedExpression]() {
								goto l308
							}
							goto l307
						l308:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if !rules[rulebuiltinCall]() {
								goto l309
							}
							goto l307
						l309:
							position, tokenIndex, depth = position307, tokenIndex307, depth307
							if !rules[rulefunctionCall]() {
								goto l292
							}
						}
					l307:
						depth--
						add(ruleconstraint, position306)
					}
					goto l291
				l292:
					position, tokenIndex, depth = position291, tokenIndex291, depth291
					{
						position310 := position
						depth++
						{
							position311, tokenIndex311, depth311 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l312
							}
							position++
							goto l311
						l312:
							position, tokenIndex, depth = position311, tokenIndex311, depth311
							if buffer[position] != rune('B') {
								goto l289
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
								goto l289
							}
							position++
						}
					l313:
						{
							position315, tokenIndex315, depth315 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l316
							}
							position++
							goto l315
						l316:
							position, tokenIndex, depth = position315, tokenIndex315, depth315
							if buffer[position] != rune('N') {
								goto l289
							}
							position++
						}
					l315:
						{
							position317, tokenIndex317, depth317 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l318
							}
							position++
							goto l317
						l318:
							position, tokenIndex, depth = position317, tokenIndex317, depth317
							if buffer[position] != rune('D') {
								goto l289
							}
							position++
						}
					l317:
						if !rules[ruleskip]() {
							goto l289
						}
						depth--
						add(ruleBIND, position310)
					}
					if !rules[ruleLPAREN]() {
						goto l289
					}
					if !rules[ruleexpression]() {
						goto l289
					}
					if !rules[ruleAS]() {
						goto l289
					}
					if !rules[rulevar]() {
						goto l289
					}
					if !rules[ruleRPAREN]() {
						goto l289
					}
				}
			l291:
				depth--
				add(rulefilterOrBind, position290)
			}
			return true
		l289:
			position, tokenIndex, depth = position289, tokenIndex289, depth289
			return false
		},
		/* 23 constraint <- <(brackettedExpression / builtinCall / functionCall)> */
		nil,
		/* 24 triplesBlock <- <(triplesSameSubjectPath (DOT triplesSameSubjectPath)* DOT?)> */
		func() bool {
			position320, tokenIndex320, depth320 := position, tokenIndex, depth
			{
				position321 := position
				depth++
				if !rules[ruletriplesSameSubjectPath]() {
					goto l320
				}
			l322:
				{
					position323, tokenIndex323, depth323 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l323
					}
					if !rules[ruletriplesSameSubjectPath]() {
						goto l323
					}
					goto l322
				l323:
					position, tokenIndex, depth = position323, tokenIndex323, depth323
				}
				{
					position324, tokenIndex324, depth324 := position, tokenIndex, depth
					if !rules[ruleDOT]() {
						goto l324
					}
					goto l325
				l324:
					position, tokenIndex, depth = position324, tokenIndex324, depth324
				}
			l325:
				depth--
				add(ruletriplesBlock, position321)
			}
			return true
		l320:
			position, tokenIndex, depth = position320, tokenIndex320, depth320
			return false
		},
		/* 25 triplesSameSubjectPath <- <((varOrTerm propertyListPath) / (triplesNodePath propertyListPath))> */
		func() bool {
			position326, tokenIndex326, depth326 := position, tokenIndex, depth
			{
				position327 := position
				depth++
				{
					position328, tokenIndex328, depth328 := position, tokenIndex, depth
					if !rules[rulevarOrTerm]() {
						goto l329
					}
					if !rules[rulepropertyListPath]() {
						goto l329
					}
					goto l328
				l329:
					position, tokenIndex, depth = position328, tokenIndex328, depth328
					{
						position330 := position
						depth++
						{
							position331, tokenIndex331, depth331 := position, tokenIndex, depth
							{
								position333 := position
								depth++
								if !rules[ruleLPAREN]() {
									goto l332
								}
								if !rules[rulegraphNodePath]() {
									goto l332
								}
							l334:
								{
									position335, tokenIndex335, depth335 := position, tokenIndex, depth
									if !rules[rulegraphNodePath]() {
										goto l335
									}
									goto l334
								l335:
									position, tokenIndex, depth = position335, tokenIndex335, depth335
								}
								if !rules[ruleRPAREN]() {
									goto l332
								}
								depth--
								add(rulecollectionPath, position333)
							}
							goto l331
						l332:
							position, tokenIndex, depth = position331, tokenIndex331, depth331
							{
								position336 := position
								depth++
								{
									position337 := position
									depth++
									if buffer[position] != rune('[') {
										goto l326
									}
									position++
									if !rules[ruleskip]() {
										goto l326
									}
									depth--
									add(ruleLBRACK, position337)
								}
								if !rules[rulepropertyListPath]() {
									goto l326
								}
								{
									position338 := position
									depth++
									if buffer[position] != rune(']') {
										goto l326
									}
									position++
									if !rules[ruleskip]() {
										goto l326
									}
									depth--
									add(ruleRBRACK, position338)
								}
								depth--
								add(ruleblankNodePropertyListPath, position336)
							}
						}
					l331:
						depth--
						add(ruletriplesNodePath, position330)
					}
					if !rules[rulepropertyListPath]() {
						goto l326
					}
				}
			l328:
				depth--
				add(ruletriplesSameSubjectPath, position327)
			}
			return true
		l326:
			position, tokenIndex, depth = position326, tokenIndex326, depth326
			return false
		},
		/* 26 varOrTerm <- <(var / graphTerm)> */
		func() bool {
			position339, tokenIndex339, depth339 := position, tokenIndex, depth
			{
				position340 := position
				depth++
				{
					position341, tokenIndex341, depth341 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l342
					}
					goto l341
				l342:
					position, tokenIndex, depth = position341, tokenIndex341, depth341
					{
						position343 := position
						depth++
						{
							position344, tokenIndex344, depth344 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l345
							}
							goto l344
						l345:
							position, tokenIndex, depth = position344, tokenIndex344, depth344
							{
								switch buffer[position] {
								case '(':
									if !rules[rulenil]() {
										goto l339
									}
									break
								case '[', '_':
									{
										position347 := position
										depth++
										{
											position348, tokenIndex348, depth348 := position, tokenIndex, depth
											{
												position350 := position
												depth++
												if buffer[position] != rune('_') {
													goto l349
												}
												position++
												if buffer[position] != rune(':') {
													goto l349
												}
												position++
												{
													switch buffer[position] {
													case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l349
														}
														position++
														break
													case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l349
														}
														position++
														break
													default:
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l349
														}
														position++
														break
													}
												}

												{
													position352, tokenIndex352, depth352 := position, tokenIndex, depth
													{
														position354, tokenIndex354, depth354 := position, tokenIndex, depth
														if c := buffer[position]; c < rune('a') || c > rune('z') {
															goto l355
														}
														position++
														goto l354
													l355:
														position, tokenIndex, depth = position354, tokenIndex354, depth354
														if c := buffer[position]; c < rune('A') || c > rune('Z') {
															goto l356
														}
														position++
														goto l354
													l356:
														position, tokenIndex, depth = position354, tokenIndex354, depth354
														if c := buffer[position]; c < rune('0') || c > rune('9') {
															goto l357
														}
														position++
														goto l354
													l357:
														position, tokenIndex, depth = position354, tokenIndex354, depth354
														if c := buffer[position]; c < rune('.') || c > rune('_') {
															goto l352
														}
														position++
													}
												l354:
													goto l353
												l352:
													position, tokenIndex, depth = position352, tokenIndex352, depth352
												}
											l353:
												if !rules[ruleskip]() {
													goto l349
												}
												depth--
												add(ruleblankNodeLabel, position350)
											}
											goto l348
										l349:
											position, tokenIndex, depth = position348, tokenIndex348, depth348
											{
												position358 := position
												depth++
												if buffer[position] != rune('[') {
													goto l339
												}
												position++
											l359:
												{
													position360, tokenIndex360, depth360 := position, tokenIndex, depth
													if !rules[rulews]() {
														goto l360
													}
													goto l359
												l360:
													position, tokenIndex, depth = position360, tokenIndex360, depth360
												}
												if buffer[position] != rune(']') {
													goto l339
												}
												position++
												if !rules[ruleskip]() {
													goto l339
												}
												depth--
												add(ruleanon, position358)
											}
										}
									l348:
										depth--
										add(ruleblankNode, position347)
									}
									break
								case 'F', 'T', 'f', 't':
									if !rules[rulebooleanLiteral]() {
										goto l339
									}
									break
								case '"':
									if !rules[ruleliteral]() {
										goto l339
									}
									break
								default:
									if !rules[rulenumericLiteral]() {
										goto l339
									}
									break
								}
							}

						}
					l344:
						depth--
						add(rulegraphTerm, position343)
					}
				}
			l341:
				depth--
				add(rulevarOrTerm, position340)
			}
			return true
		l339:
			position, tokenIndex, depth = position339, tokenIndex339, depth339
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
			position365, tokenIndex365, depth365 := position, tokenIndex, depth
			{
				position366 := position
				depth++
				{
					position367, tokenIndex367, depth367 := position, tokenIndex, depth
					if !rules[rulevar]() {
						goto l368
					}
					goto l367
				l368:
					position, tokenIndex, depth = position367, tokenIndex367, depth367
					{
						position369 := position
						depth++
						if !rules[rulepath]() {
							goto l365
						}
						depth--
						add(ruleverbPath, position369)
					}
				}
			l367:
				if !rules[ruleobjectListPath]() {
					goto l365
				}
				{
					position370, tokenIndex370, depth370 := position, tokenIndex, depth
					{
						position372 := position
						depth++
						if buffer[position] != rune(';') {
							goto l370
						}
						position++
						if !rules[ruleskip]() {
							goto l370
						}
						depth--
						add(ruleSEMICOLON, position372)
					}
					if !rules[rulepropertyListPath]() {
						goto l370
					}
					goto l371
				l370:
					position, tokenIndex, depth = position370, tokenIndex370, depth370
				}
			l371:
				depth--
				add(rulepropertyListPath, position366)
			}
			return true
		l365:
			position, tokenIndex, depth = position365, tokenIndex365, depth365
			return false
		},
		/* 32 verbPath <- <path> */
		nil,
		/* 33 path <- <pathAlternative> */
		func() bool {
			position374, tokenIndex374, depth374 := position, tokenIndex, depth
			{
				position375 := position
				depth++
				if !rules[rulepathAlternative]() {
					goto l374
				}
				depth--
				add(rulepath, position375)
			}
			return true
		l374:
			position, tokenIndex, depth = position374, tokenIndex374, depth374
			return false
		},
		/* 34 pathAlternative <- <(pathSequence (PIPE pathAlternative)*)> */
		func() bool {
			position376, tokenIndex376, depth376 := position, tokenIndex, depth
			{
				position377 := position
				depth++
				if !rules[rulepathSequence]() {
					goto l376
				}
			l378:
				{
					position379, tokenIndex379, depth379 := position, tokenIndex, depth
					if !rules[rulePIPE]() {
						goto l379
					}
					if !rules[rulepathAlternative]() {
						goto l379
					}
					goto l378
				l379:
					position, tokenIndex, depth = position379, tokenIndex379, depth379
				}
				depth--
				add(rulepathAlternative, position377)
			}
			return true
		l376:
			position, tokenIndex, depth = position376, tokenIndex376, depth376
			return false
		},
		/* 35 pathSequence <- <(pathElt (SLASH pathSequence)*)> */
		func() bool {
			position380, tokenIndex380, depth380 := position, tokenIndex, depth
			{
				position381 := position
				depth++
				{
					position382 := position
					depth++
					{
						position383, tokenIndex383, depth383 := position, tokenIndex, depth
						if !rules[ruleINVERSE]() {
							goto l383
						}
						goto l384
					l383:
						position, tokenIndex, depth = position383, tokenIndex383, depth383
					}
				l384:
					{
						position385 := position
						depth++
						{
							position386, tokenIndex386, depth386 := position, tokenIndex, depth
							if !rules[ruleiriref]() {
								goto l387
							}
							goto l386
						l387:
							position, tokenIndex, depth = position386, tokenIndex386, depth386
							{
								switch buffer[position] {
								case '(':
									if !rules[ruleLPAREN]() {
										goto l380
									}
									if !rules[rulepath]() {
										goto l380
									}
									if !rules[ruleRPAREN]() {
										goto l380
									}
									break
								case '!':
									if !rules[ruleNOT]() {
										goto l380
									}
									{
										position389 := position
										depth++
										{
											position390, tokenIndex390, depth390 := position, tokenIndex, depth
											if !rules[rulepathOneInPropertySet]() {
												goto l391
											}
											goto l390
										l391:
											position, tokenIndex, depth = position390, tokenIndex390, depth390
											if !rules[ruleLPAREN]() {
												goto l380
											}
											{
												position392, tokenIndex392, depth392 := position, tokenIndex, depth
												if !rules[rulepathOneInPropertySet]() {
													goto l392
												}
											l394:
												{
													position395, tokenIndex395, depth395 := position, tokenIndex, depth
													if !rules[rulePIPE]() {
														goto l395
													}
													if !rules[rulepathOneInPropertySet]() {
														goto l395
													}
													goto l394
												l395:
													position, tokenIndex, depth = position395, tokenIndex395, depth395
												}
												goto l393
											l392:
												position, tokenIndex, depth = position392, tokenIndex392, depth392
											}
										l393:
											if !rules[ruleRPAREN]() {
												goto l380
											}
										}
									l390:
										depth--
										add(rulepathNegatedPropertySet, position389)
									}
									break
								default:
									if !rules[ruleISA]() {
										goto l380
									}
									break
								}
							}

						}
					l386:
						depth--
						add(rulepathPrimary, position385)
					}
					depth--
					add(rulepathElt, position382)
				}
			l396:
				{
					position397, tokenIndex397, depth397 := position, tokenIndex, depth
					if !rules[ruleSLASH]() {
						goto l397
					}
					if !rules[rulepathSequence]() {
						goto l397
					}
					goto l396
				l397:
					position, tokenIndex, depth = position397, tokenIndex397, depth397
				}
				depth--
				add(rulepathSequence, position381)
			}
			return true
		l380:
			position, tokenIndex, depth = position380, tokenIndex380, depth380
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
			position401, tokenIndex401, depth401 := position, tokenIndex, depth
			{
				position402 := position
				depth++
				{
					position403, tokenIndex403, depth403 := position, tokenIndex, depth
					if !rules[ruleiriref]() {
						goto l404
					}
					goto l403
				l404:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
					if !rules[ruleISA]() {
						goto l405
					}
					goto l403
				l405:
					position, tokenIndex, depth = position403, tokenIndex403, depth403
					if !rules[ruleINVERSE]() {
						goto l401
					}
					{
						position406, tokenIndex406, depth406 := position, tokenIndex, depth
						if !rules[ruleiriref]() {
							goto l407
						}
						goto l406
					l407:
						position, tokenIndex, depth = position406, tokenIndex406, depth406
						if !rules[ruleISA]() {
							goto l401
						}
					}
				l406:
				}
			l403:
				depth--
				add(rulepathOneInPropertySet, position402)
			}
			return true
		l401:
			position, tokenIndex, depth = position401, tokenIndex401, depth401
			return false
		},
		/* 40 objectListPath <- <(objectPath (COMMA objectListPath)*)> */
		func() bool {
			position408, tokenIndex408, depth408 := position, tokenIndex, depth
			{
				position409 := position
				depth++
				{
					position410 := position
					depth++
					if !rules[rulegraphNodePath]() {
						goto l408
					}
					depth--
					add(ruleobjectPath, position410)
				}
			l411:
				{
					position412, tokenIndex412, depth412 := position, tokenIndex, depth
					if !rules[ruleCOMMA]() {
						goto l412
					}
					if !rules[ruleobjectListPath]() {
						goto l412
					}
					goto l411
				l412:
					position, tokenIndex, depth = position412, tokenIndex412, depth412
				}
				depth--
				add(ruleobjectListPath, position409)
			}
			return true
		l408:
			position, tokenIndex, depth = position408, tokenIndex408, depth408
			return false
		},
		/* 41 objectPath <- <graphNodePath> */
		nil,
		/* 42 graphNodePath <- <varOrTerm> */
		func() bool {
			position414, tokenIndex414, depth414 := position, tokenIndex, depth
			{
				position415 := position
				depth++
				if !rules[rulevarOrTerm]() {
					goto l414
				}
				depth--
				add(rulegraphNodePath, position415)
			}
			return true
		l414:
			position, tokenIndex, depth = position414, tokenIndex414, depth414
			return false
		},
		/* 43 solutionModifier <- <limitOffsetClauses?> */
		func() bool {
			{
				position417 := position
				depth++
				{
					position418, tokenIndex418, depth418 := position, tokenIndex, depth
					{
						position420 := position
						depth++
						{
							position421, tokenIndex421, depth421 := position, tokenIndex, depth
							if !rules[rulelimit]() {
								goto l422
							}
							{
								position423, tokenIndex423, depth423 := position, tokenIndex, depth
								if !rules[ruleoffset]() {
									goto l423
								}
								goto l424
							l423:
								position, tokenIndex, depth = position423, tokenIndex423, depth423
							}
						l424:
							goto l421
						l422:
							position, tokenIndex, depth = position421, tokenIndex421, depth421
							if !rules[ruleoffset]() {
								goto l418
							}
							{
								position425, tokenIndex425, depth425 := position, tokenIndex, depth
								if !rules[rulelimit]() {
									goto l425
								}
								goto l426
							l425:
								position, tokenIndex, depth = position425, tokenIndex425, depth425
							}
						l426:
						}
					l421:
						depth--
						add(rulelimitOffsetClauses, position420)
					}
					goto l419
				l418:
					position, tokenIndex, depth = position418, tokenIndex418, depth418
				}
			l419:
				depth--
				add(rulesolutionModifier, position417)
			}
			return true
		},
		/* 44 limitOffsetClauses <- <((limit offset?) / (offset limit?))> */
		nil,
		/* 45 limit <- <(LIMIT INTEGER)> */
		func() bool {
			position428, tokenIndex428, depth428 := position, tokenIndex, depth
			{
				position429 := position
				depth++
				{
					position430 := position
					depth++
					{
						position431, tokenIndex431, depth431 := position, tokenIndex, depth
						if buffer[position] != rune('l') {
							goto l432
						}
						position++
						goto l431
					l432:
						position, tokenIndex, depth = position431, tokenIndex431, depth431
						if buffer[position] != rune('L') {
							goto l428
						}
						position++
					}
				l431:
					{
						position433, tokenIndex433, depth433 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l434
						}
						position++
						goto l433
					l434:
						position, tokenIndex, depth = position433, tokenIndex433, depth433
						if buffer[position] != rune('I') {
							goto l428
						}
						position++
					}
				l433:
					{
						position435, tokenIndex435, depth435 := position, tokenIndex, depth
						if buffer[position] != rune('m') {
							goto l436
						}
						position++
						goto l435
					l436:
						position, tokenIndex, depth = position435, tokenIndex435, depth435
						if buffer[position] != rune('M') {
							goto l428
						}
						position++
					}
				l435:
					{
						position437, tokenIndex437, depth437 := position, tokenIndex, depth
						if buffer[position] != rune('i') {
							goto l438
						}
						position++
						goto l437
					l438:
						position, tokenIndex, depth = position437, tokenIndex437, depth437
						if buffer[position] != rune('I') {
							goto l428
						}
						position++
					}
				l437:
					{
						position439, tokenIndex439, depth439 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l440
						}
						position++
						goto l439
					l440:
						position, tokenIndex, depth = position439, tokenIndex439, depth439
						if buffer[position] != rune('T') {
							goto l428
						}
						position++
					}
				l439:
					if !rules[ruleskip]() {
						goto l428
					}
					depth--
					add(ruleLIMIT, position430)
				}
				if !rules[ruleINTEGER]() {
					goto l428
				}
				depth--
				add(rulelimit, position429)
			}
			return true
		l428:
			position, tokenIndex, depth = position428, tokenIndex428, depth428
			return false
		},
		/* 46 offset <- <(OFFSET INTEGER)> */
		func() bool {
			position441, tokenIndex441, depth441 := position, tokenIndex, depth
			{
				position442 := position
				depth++
				{
					position443 := position
					depth++
					{
						position444, tokenIndex444, depth444 := position, tokenIndex, depth
						if buffer[position] != rune('o') {
							goto l445
						}
						position++
						goto l444
					l445:
						position, tokenIndex, depth = position444, tokenIndex444, depth444
						if buffer[position] != rune('O') {
							goto l441
						}
						position++
					}
				l444:
					{
						position446, tokenIndex446, depth446 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l447
						}
						position++
						goto l446
					l447:
						position, tokenIndex, depth = position446, tokenIndex446, depth446
						if buffer[position] != rune('F') {
							goto l441
						}
						position++
					}
				l446:
					{
						position448, tokenIndex448, depth448 := position, tokenIndex, depth
						if buffer[position] != rune('f') {
							goto l449
						}
						position++
						goto l448
					l449:
						position, tokenIndex, depth = position448, tokenIndex448, depth448
						if buffer[position] != rune('F') {
							goto l441
						}
						position++
					}
				l448:
					{
						position450, tokenIndex450, depth450 := position, tokenIndex, depth
						if buffer[position] != rune('s') {
							goto l451
						}
						position++
						goto l450
					l451:
						position, tokenIndex, depth = position450, tokenIndex450, depth450
						if buffer[position] != rune('S') {
							goto l441
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
							goto l441
						}
						position++
					}
				l452:
					{
						position454, tokenIndex454, depth454 := position, tokenIndex, depth
						if buffer[position] != rune('t') {
							goto l455
						}
						position++
						goto l454
					l455:
						position, tokenIndex, depth = position454, tokenIndex454, depth454
						if buffer[position] != rune('T') {
							goto l441
						}
						position++
					}
				l454:
					if !rules[ruleskip]() {
						goto l441
					}
					depth--
					add(ruleOFFSET, position443)
				}
				if !rules[ruleINTEGER]() {
					goto l441
				}
				depth--
				add(ruleoffset, position442)
			}
			return true
		l441:
			position, tokenIndex, depth = position441, tokenIndex441, depth441
			return false
		},
		/* 47 expression <- <conditionalOrExpression> */
		func() bool {
			position456, tokenIndex456, depth456 := position, tokenIndex, depth
			{
				position457 := position
				depth++
				if !rules[ruleconditionalOrExpression]() {
					goto l456
				}
				depth--
				add(ruleexpression, position457)
			}
			return true
		l456:
			position, tokenIndex, depth = position456, tokenIndex456, depth456
			return false
		},
		/* 48 conditionalOrExpression <- <(conditionalAndExpression (OR conditionalOrExpression)?)> */
		func() bool {
			position458, tokenIndex458, depth458 := position, tokenIndex, depth
			{
				position459 := position
				depth++
				if !rules[ruleconditionalAndExpression]() {
					goto l458
				}
				{
					position460, tokenIndex460, depth460 := position, tokenIndex, depth
					{
						position462 := position
						depth++
						if buffer[position] != rune('|') {
							goto l460
						}
						position++
						if buffer[position] != rune('|') {
							goto l460
						}
						position++
						if !rules[ruleskip]() {
							goto l460
						}
						depth--
						add(ruleOR, position462)
					}
					if !rules[ruleconditionalOrExpression]() {
						goto l460
					}
					goto l461
				l460:
					position, tokenIndex, depth = position460, tokenIndex460, depth460
				}
			l461:
				depth--
				add(ruleconditionalOrExpression, position459)
			}
			return true
		l458:
			position, tokenIndex, depth = position458, tokenIndex458, depth458
			return false
		},
		/* 49 conditionalAndExpression <- <(valueLogical (AND conditionalAndExpression)?)> */
		func() bool {
			position463, tokenIndex463, depth463 := position, tokenIndex, depth
			{
				position464 := position
				depth++
				{
					position465 := position
					depth++
					if !rules[rulenumericExpression]() {
						goto l463
					}
					{
						position466, tokenIndex466, depth466 := position, tokenIndex, depth
						{
							switch buffer[position] {
							case 'N', 'n':
								{
									position469 := position
									depth++
									{
										position470 := position
										depth++
										{
											position471, tokenIndex471, depth471 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l472
											}
											position++
											goto l471
										l472:
											position, tokenIndex, depth = position471, tokenIndex471, depth471
											if buffer[position] != rune('N') {
												goto l466
											}
											position++
										}
									l471:
										{
											position473, tokenIndex473, depth473 := position, tokenIndex, depth
											if buffer[position] != rune('o') {
												goto l474
											}
											position++
											goto l473
										l474:
											position, tokenIndex, depth = position473, tokenIndex473, depth473
											if buffer[position] != rune('O') {
												goto l466
											}
											position++
										}
									l473:
										{
											position475, tokenIndex475, depth475 := position, tokenIndex, depth
											if buffer[position] != rune('t') {
												goto l476
											}
											position++
											goto l475
										l476:
											position, tokenIndex, depth = position475, tokenIndex475, depth475
											if buffer[position] != rune('T') {
												goto l466
											}
											position++
										}
									l475:
										if buffer[position] != rune(' ') {
											goto l466
										}
										position++
										{
											position477, tokenIndex477, depth477 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l478
											}
											position++
											goto l477
										l478:
											position, tokenIndex, depth = position477, tokenIndex477, depth477
											if buffer[position] != rune('I') {
												goto l466
											}
											position++
										}
									l477:
										{
											position479, tokenIndex479, depth479 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l480
											}
											position++
											goto l479
										l480:
											position, tokenIndex, depth = position479, tokenIndex479, depth479
											if buffer[position] != rune('N') {
												goto l466
											}
											position++
										}
									l479:
										if !rules[ruleskip]() {
											goto l466
										}
										depth--
										add(ruleNOTIN, position470)
									}
									if !rules[ruleargList]() {
										goto l466
									}
									depth--
									add(rulenotin, position469)
								}
								break
							case 'I', 'i':
								{
									position481 := position
									depth++
									{
										position482 := position
										depth++
										{
											position483, tokenIndex483, depth483 := position, tokenIndex, depth
											if buffer[position] != rune('i') {
												goto l484
											}
											position++
											goto l483
										l484:
											position, tokenIndex, depth = position483, tokenIndex483, depth483
											if buffer[position] != rune('I') {
												goto l466
											}
											position++
										}
									l483:
										{
											position485, tokenIndex485, depth485 := position, tokenIndex, depth
											if buffer[position] != rune('n') {
												goto l486
											}
											position++
											goto l485
										l486:
											position, tokenIndex, depth = position485, tokenIndex485, depth485
											if buffer[position] != rune('N') {
												goto l466
											}
											position++
										}
									l485:
										if !rules[ruleskip]() {
											goto l466
										}
										depth--
										add(ruleIN, position482)
									}
									if !rules[ruleargList]() {
										goto l466
									}
									depth--
									add(rulein, position481)
								}
								break
							default:
								{
									position487, tokenIndex487, depth487 := position, tokenIndex, depth
									{
										position489 := position
										depth++
										if buffer[position] != rune('<') {
											goto l488
										}
										position++
										if !rules[ruleskip]() {
											goto l488
										}
										depth--
										add(ruleLT, position489)
									}
									goto l487
								l488:
									position, tokenIndex, depth = position487, tokenIndex487, depth487
									{
										position491 := position
										depth++
										if buffer[position] != rune('>') {
											goto l490
										}
										position++
										if buffer[position] != rune('=') {
											goto l490
										}
										position++
										if !rules[ruleskip]() {
											goto l490
										}
										depth--
										add(ruleGE, position491)
									}
									goto l487
								l490:
									position, tokenIndex, depth = position487, tokenIndex487, depth487
									{
										switch buffer[position] {
										case '>':
											{
												position493 := position
												depth++
												if buffer[position] != rune('>') {
													goto l466
												}
												position++
												if !rules[ruleskip]() {
													goto l466
												}
												depth--
												add(ruleGT, position493)
											}
											break
										case '<':
											{
												position494 := position
												depth++
												if buffer[position] != rune('<') {
													goto l466
												}
												position++
												if buffer[position] != rune('=') {
													goto l466
												}
												position++
												if !rules[ruleskip]() {
													goto l466
												}
												depth--
												add(ruleLE, position494)
											}
											break
										case '!':
											{
												position495 := position
												depth++
												if buffer[position] != rune('!') {
													goto l466
												}
												position++
												if buffer[position] != rune('=') {
													goto l466
												}
												position++
												if !rules[ruleskip]() {
													goto l466
												}
												depth--
												add(ruleNE, position495)
											}
											break
										default:
											{
												position496 := position
												depth++
												if buffer[position] != rune('=') {
													goto l466
												}
												position++
												if !rules[ruleskip]() {
													goto l466
												}
												depth--
												add(ruleEQ, position496)
											}
											break
										}
									}

								}
							l487:
								if !rules[rulenumericExpression]() {
									goto l466
								}
								break
							}
						}

						goto l467
					l466:
						position, tokenIndex, depth = position466, tokenIndex466, depth466
					}
				l467:
					depth--
					add(rulevalueLogical, position465)
				}
				{
					position497, tokenIndex497, depth497 := position, tokenIndex, depth
					{
						position499 := position
						depth++
						if buffer[position] != rune('&') {
							goto l497
						}
						position++
						if buffer[position] != rune('&') {
							goto l497
						}
						position++
						if !rules[ruleskip]() {
							goto l497
						}
						depth--
						add(ruleAND, position499)
					}
					if !rules[ruleconditionalAndExpression]() {
						goto l497
					}
					goto l498
				l497:
					position, tokenIndex, depth = position497, tokenIndex497, depth497
				}
			l498:
				depth--
				add(ruleconditionalAndExpression, position464)
			}
			return true
		l463:
			position, tokenIndex, depth = position463, tokenIndex463, depth463
			return false
		},
		/* 50 valueLogical <- <(numericExpression ((&('N' | 'n') notin) | (&('I' | 'i') in) | (&('!' | '<' | '=' | '>') ((LT / GE / ((&('>') GT) | (&('<') LE) | (&('!') NE) | (&('=') EQ))) numericExpression)))?)> */
		nil,
		/* 51 numericExpression <- <(multiplicativeExpression (((PLUS / MINUS) multiplicativeExpression) / signedNumericLiteral)*)> */
		func() bool {
			position501, tokenIndex501, depth501 := position, tokenIndex, depth
			{
				position502 := position
				depth++
				if !rules[rulemultiplicativeExpression]() {
					goto l501
				}
			l503:
				{
					position504, tokenIndex504, depth504 := position, tokenIndex, depth
					{
						position505, tokenIndex505, depth505 := position, tokenIndex, depth
						{
							position507, tokenIndex507, depth507 := position, tokenIndex, depth
							if !rules[rulePLUS]() {
								goto l508
							}
							goto l507
						l508:
							position, tokenIndex, depth = position507, tokenIndex507, depth507
							if !rules[ruleMINUS]() {
								goto l506
							}
						}
					l507:
						if !rules[rulemultiplicativeExpression]() {
							goto l506
						}
						goto l505
					l506:
						position, tokenIndex, depth = position505, tokenIndex505, depth505
						{
							position509 := position
							depth++
							{
								position510, tokenIndex510, depth510 := position, tokenIndex, depth
								if buffer[position] != rune('+') {
									goto l511
								}
								position++
								goto l510
							l511:
								position, tokenIndex, depth = position510, tokenIndex510, depth510
								if buffer[position] != rune('-') {
									goto l504
								}
								position++
							}
						l510:
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l504
							}
							position++
						l512:
							{
								position513, tokenIndex513, depth513 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l513
								}
								position++
								goto l512
							l513:
								position, tokenIndex, depth = position513, tokenIndex513, depth513
							}
							{
								position514, tokenIndex514, depth514 := position, tokenIndex, depth
								if buffer[position] != rune('.') {
									goto l514
								}
								position++
							l516:
								{
									position517, tokenIndex517, depth517 := position, tokenIndex, depth
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l517
									}
									position++
									goto l516
								l517:
									position, tokenIndex, depth = position517, tokenIndex517, depth517
								}
								goto l515
							l514:
								position, tokenIndex, depth = position514, tokenIndex514, depth514
							}
						l515:
							if !rules[ruleskip]() {
								goto l504
							}
							depth--
							add(rulesignedNumericLiteral, position509)
						}
					}
				l505:
					goto l503
				l504:
					position, tokenIndex, depth = position504, tokenIndex504, depth504
				}
				depth--
				add(rulenumericExpression, position502)
			}
			return true
		l501:
			position, tokenIndex, depth = position501, tokenIndex501, depth501
			return false
		},
		/* 52 multiplicativeExpression <- <(unaryExpression ((STAR / SLASH) unaryExpression)*)> */
		func() bool {
			position518, tokenIndex518, depth518 := position, tokenIndex, depth
			{
				position519 := position
				depth++
				if !rules[ruleunaryExpression]() {
					goto l518
				}
			l520:
				{
					position521, tokenIndex521, depth521 := position, tokenIndex, depth
					{
						position522, tokenIndex522, depth522 := position, tokenIndex, depth
						if !rules[ruleSTAR]() {
							goto l523
						}
						goto l522
					l523:
						position, tokenIndex, depth = position522, tokenIndex522, depth522
						if !rules[ruleSLASH]() {
							goto l521
						}
					}
				l522:
					if !rules[ruleunaryExpression]() {
						goto l521
					}
					goto l520
				l521:
					position, tokenIndex, depth = position521, tokenIndex521, depth521
				}
				depth--
				add(rulemultiplicativeExpression, position519)
			}
			return true
		l518:
			position, tokenIndex, depth = position518, tokenIndex518, depth518
			return false
		},
		/* 53 unaryExpression <- <(((&('+') PLUS) | (&('-') MINUS) | (&('!') NOT))? primaryExpression)> */
		func() bool {
			position524, tokenIndex524, depth524 := position, tokenIndex, depth
			{
				position525 := position
				depth++
				{
					position526, tokenIndex526, depth526 := position, tokenIndex, depth
					{
						switch buffer[position] {
						case '+':
							if !rules[rulePLUS]() {
								goto l526
							}
							break
						case '-':
							if !rules[ruleMINUS]() {
								goto l526
							}
							break
						default:
							if !rules[ruleNOT]() {
								goto l526
							}
							break
						}
					}

					goto l527
				l526:
					position, tokenIndex, depth = position526, tokenIndex526, depth526
				}
			l527:
				{
					position529 := position
					depth++
					{
						position530, tokenIndex530, depth530 := position, tokenIndex, depth
						if !rules[rulebrackettedExpression]() {
							goto l531
						}
						goto l530
					l531:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
						if !rules[rulebuiltinCall]() {
							goto l532
						}
						goto l530
					l532:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
						if !rules[rulefunctionCall]() {
							goto l533
						}
						goto l530
					l533:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
						if !rules[ruleiriref]() {
							goto l534
						}
						goto l530
					l534:
						position, tokenIndex, depth = position530, tokenIndex530, depth530
						{
							switch buffer[position] {
							case '$', '?':
								if !rules[rulevar]() {
									goto l524
								}
								break
							case 'F', 'T', 'f', 't':
								if !rules[rulebooleanLiteral]() {
									goto l524
								}
								break
							case '"':
								if !rules[ruleliteral]() {
									goto l524
								}
								break
							default:
								if !rules[rulenumericLiteral]() {
									goto l524
								}
								break
							}
						}

					}
				l530:
					depth--
					add(ruleprimaryExpression, position529)
				}
				depth--
				add(ruleunaryExpression, position525)
			}
			return true
		l524:
			position, tokenIndex, depth = position524, tokenIndex524, depth524
			return false
		},
		/* 54 primaryExpression <- <(brackettedExpression / builtinCall / functionCall / iriref / ((&('$' | '?') var) | (&('F' | 'T' | 'f' | 't') booleanLiteral) | (&('"') literal) | (&('+' | '-' | '0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') numericLiteral)))> */
		nil,
		/* 55 brackettedExpression <- <(LPAREN expression RPAREN)> */
		func() bool {
			position537, tokenIndex537, depth537 := position, tokenIndex, depth
			{
				position538 := position
				depth++
				if !rules[ruleLPAREN]() {
					goto l537
				}
				if !rules[ruleexpression]() {
					goto l537
				}
				if !rules[ruleRPAREN]() {
					goto l537
				}
				depth--
				add(rulebrackettedExpression, position538)
			}
			return true
		l537:
			position, tokenIndex, depth = position537, tokenIndex537, depth537
			return false
		},
		/* 56 functionCall <- <(iriref argList)> */
		func() bool {
			position539, tokenIndex539, depth539 := position, tokenIndex, depth
			{
				position540 := position
				depth++
				if !rules[ruleiriref]() {
					goto l539
				}
				if !rules[ruleargList]() {
					goto l539
				}
				depth--
				add(rulefunctionCall, position540)
			}
			return true
		l539:
			position, tokenIndex, depth = position539, tokenIndex539, depth539
			return false
		},
		/* 57 in <- <(IN argList)> */
		nil,
		/* 58 notin <- <(NOTIN argList)> */
		nil,
		/* 59 argList <- <(nil / (LPAREN expression (COMMA expression)* RPAREN))> */
		func() bool {
			position543, tokenIndex543, depth543 := position, tokenIndex, depth
			{
				position544 := position
				depth++
				{
					position545, tokenIndex545, depth545 := position, tokenIndex, depth
					if !rules[rulenil]() {
						goto l546
					}
					goto l545
				l546:
					position, tokenIndex, depth = position545, tokenIndex545, depth545
					if !rules[ruleLPAREN]() {
						goto l543
					}
					if !rules[ruleexpression]() {
						goto l543
					}
				l547:
					{
						position548, tokenIndex548, depth548 := position, tokenIndex, depth
						if !rules[ruleCOMMA]() {
							goto l548
						}
						if !rules[ruleexpression]() {
							goto l548
						}
						goto l547
					l548:
						position, tokenIndex, depth = position548, tokenIndex548, depth548
					}
					if !rules[ruleRPAREN]() {
						goto l543
					}
				}
			l545:
				depth--
				add(ruleargList, position544)
			}
			return true
		l543:
			position, tokenIndex, depth = position543, tokenIndex543, depth543
			return false
		},
		/* 60 builtinCall <- <(((STR / LANG / DATATYPE / IRI / URI / STRLEN / MONTH / MINUTES / SECONDS / TIMEZONE / SHA1 / SHA256 / SHA384 / ISIRI / ISURI / ISBLANK / ISLITERAL / ((&('I' | 'i') ISNUMERIC) | (&('S' | 's') SHA512) | (&('M' | 'm') MD5) | (&('T' | 't') TZ) | (&('H' | 'h') HOURS) | (&('D' | 'd') DAY) | (&('Y' | 'y') YEAR) | (&('E' | 'e') ENCODEFORURI) | (&('L' | 'l') LCASE) | (&('U' | 'u') UCASE) | (&('F' | 'f') FLOOR) | (&('R' | 'r') ROUND) | (&('C' | 'c') CEIL) | (&('A' | 'a') ABS))) LPAREN expression RPAREN) / ((STRSTARTS / STRENDS / STRBEFORE / STRAFTER / STRLANG / STRDT / ((&('S' | 's') SAMETERM) | (&('C' | 'c') CONTAINS) | (&('L' | 'l') LANGMATCHES))) LPAREN expression COMMA expression RPAREN) / (BOUND LPAREN var RPAREN) / (((&('S' | 's') STRUUID) | (&('U' | 'u') UUID) | (&('N' | 'n') NOW) | (&('R' | 'r') RAND)) nil) / ((&('E' | 'N' | 'e' | 'n') ((EXISTS / NOTEXIST) groupGraphPattern)) | (&('I' | 'i') (IF LPAREN expression COMMA expression COMMA expression RPAREN)) | (&('C' | 'c') ((CONCAT / COALESCE) argList)) | (&('B' | 'b') (BNODE ((LPAREN expression RPAREN) / nil))) | (&('R' | 'S' | 'r' | 's') ((SUBSTR / REPLACE / REGEX) LPAREN expression COMMA expression (COMMA expression)? RPAREN))))> */
		func() bool {
			position549, tokenIndex549, depth549 := position, tokenIndex, depth
			{
				position550 := position
				depth++
				{
					position551, tokenIndex551, depth551 := position, tokenIndex, depth
					{
						position553, tokenIndex553, depth553 := position, tokenIndex, depth
						{
							position555 := position
							depth++
							{
								position556, tokenIndex556, depth556 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l557
								}
								position++
								goto l556
							l557:
								position, tokenIndex, depth = position556, tokenIndex556, depth556
								if buffer[position] != rune('S') {
									goto l554
								}
								position++
							}
						l556:
							{
								position558, tokenIndex558, depth558 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l559
								}
								position++
								goto l558
							l559:
								position, tokenIndex, depth = position558, tokenIndex558, depth558
								if buffer[position] != rune('T') {
									goto l554
								}
								position++
							}
						l558:
							{
								position560, tokenIndex560, depth560 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l561
								}
								position++
								goto l560
							l561:
								position, tokenIndex, depth = position560, tokenIndex560, depth560
								if buffer[position] != rune('R') {
									goto l554
								}
								position++
							}
						l560:
							if !rules[ruleskip]() {
								goto l554
							}
							depth--
							add(ruleSTR, position555)
						}
						goto l553
					l554:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position563 := position
							depth++
							{
								position564, tokenIndex564, depth564 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l565
								}
								position++
								goto l564
							l565:
								position, tokenIndex, depth = position564, tokenIndex564, depth564
								if buffer[position] != rune('L') {
									goto l562
								}
								position++
							}
						l564:
							{
								position566, tokenIndex566, depth566 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l567
								}
								position++
								goto l566
							l567:
								position, tokenIndex, depth = position566, tokenIndex566, depth566
								if buffer[position] != rune('A') {
									goto l562
								}
								position++
							}
						l566:
							{
								position568, tokenIndex568, depth568 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l569
								}
								position++
								goto l568
							l569:
								position, tokenIndex, depth = position568, tokenIndex568, depth568
								if buffer[position] != rune('N') {
									goto l562
								}
								position++
							}
						l568:
							{
								position570, tokenIndex570, depth570 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l571
								}
								position++
								goto l570
							l571:
								position, tokenIndex, depth = position570, tokenIndex570, depth570
								if buffer[position] != rune('G') {
									goto l562
								}
								position++
							}
						l570:
							if !rules[ruleskip]() {
								goto l562
							}
							depth--
							add(ruleLANG, position563)
						}
						goto l553
					l562:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position573 := position
							depth++
							{
								position574, tokenIndex574, depth574 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l575
								}
								position++
								goto l574
							l575:
								position, tokenIndex, depth = position574, tokenIndex574, depth574
								if buffer[position] != rune('D') {
									goto l572
								}
								position++
							}
						l574:
							{
								position576, tokenIndex576, depth576 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l577
								}
								position++
								goto l576
							l577:
								position, tokenIndex, depth = position576, tokenIndex576, depth576
								if buffer[position] != rune('A') {
									goto l572
								}
								position++
							}
						l576:
							{
								position578, tokenIndex578, depth578 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l579
								}
								position++
								goto l578
							l579:
								position, tokenIndex, depth = position578, tokenIndex578, depth578
								if buffer[position] != rune('T') {
									goto l572
								}
								position++
							}
						l578:
							{
								position580, tokenIndex580, depth580 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l581
								}
								position++
								goto l580
							l581:
								position, tokenIndex, depth = position580, tokenIndex580, depth580
								if buffer[position] != rune('A') {
									goto l572
								}
								position++
							}
						l580:
							{
								position582, tokenIndex582, depth582 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l583
								}
								position++
								goto l582
							l583:
								position, tokenIndex, depth = position582, tokenIndex582, depth582
								if buffer[position] != rune('T') {
									goto l572
								}
								position++
							}
						l582:
							{
								position584, tokenIndex584, depth584 := position, tokenIndex, depth
								if buffer[position] != rune('y') {
									goto l585
								}
								position++
								goto l584
							l585:
								position, tokenIndex, depth = position584, tokenIndex584, depth584
								if buffer[position] != rune('Y') {
									goto l572
								}
								position++
							}
						l584:
							{
								position586, tokenIndex586, depth586 := position, tokenIndex, depth
								if buffer[position] != rune('p') {
									goto l587
								}
								position++
								goto l586
							l587:
								position, tokenIndex, depth = position586, tokenIndex586, depth586
								if buffer[position] != rune('P') {
									goto l572
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
									goto l572
								}
								position++
							}
						l588:
							if !rules[ruleskip]() {
								goto l572
							}
							depth--
							add(ruleDATATYPE, position573)
						}
						goto l553
					l572:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position591 := position
							depth++
							{
								position592, tokenIndex592, depth592 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l593
								}
								position++
								goto l592
							l593:
								position, tokenIndex, depth = position592, tokenIndex592, depth592
								if buffer[position] != rune('I') {
									goto l590
								}
								position++
							}
						l592:
							{
								position594, tokenIndex594, depth594 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l595
								}
								position++
								goto l594
							l595:
								position, tokenIndex, depth = position594, tokenIndex594, depth594
								if buffer[position] != rune('R') {
									goto l590
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
									goto l590
								}
								position++
							}
						l596:
							if !rules[ruleskip]() {
								goto l590
							}
							depth--
							add(ruleIRI, position591)
						}
						goto l553
					l590:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position599 := position
							depth++
							{
								position600, tokenIndex600, depth600 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l601
								}
								position++
								goto l600
							l601:
								position, tokenIndex, depth = position600, tokenIndex600, depth600
								if buffer[position] != rune('U') {
									goto l598
								}
								position++
							}
						l600:
							{
								position602, tokenIndex602, depth602 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l603
								}
								position++
								goto l602
							l603:
								position, tokenIndex, depth = position602, tokenIndex602, depth602
								if buffer[position] != rune('R') {
									goto l598
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
									goto l598
								}
								position++
							}
						l604:
							if !rules[ruleskip]() {
								goto l598
							}
							depth--
							add(ruleURI, position599)
						}
						goto l553
					l598:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position607 := position
							depth++
							{
								position608, tokenIndex608, depth608 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l609
								}
								position++
								goto l608
							l609:
								position, tokenIndex, depth = position608, tokenIndex608, depth608
								if buffer[position] != rune('S') {
									goto l606
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
									goto l606
								}
								position++
							}
						l610:
							{
								position612, tokenIndex612, depth612 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l613
								}
								position++
								goto l612
							l613:
								position, tokenIndex, depth = position612, tokenIndex612, depth612
								if buffer[position] != rune('R') {
									goto l606
								}
								position++
							}
						l612:
							{
								position614, tokenIndex614, depth614 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l615
								}
								position++
								goto l614
							l615:
								position, tokenIndex, depth = position614, tokenIndex614, depth614
								if buffer[position] != rune('L') {
									goto l606
								}
								position++
							}
						l614:
							{
								position616, tokenIndex616, depth616 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l617
								}
								position++
								goto l616
							l617:
								position, tokenIndex, depth = position616, tokenIndex616, depth616
								if buffer[position] != rune('E') {
									goto l606
								}
								position++
							}
						l616:
							{
								position618, tokenIndex618, depth618 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l619
								}
								position++
								goto l618
							l619:
								position, tokenIndex, depth = position618, tokenIndex618, depth618
								if buffer[position] != rune('N') {
									goto l606
								}
								position++
							}
						l618:
							if !rules[ruleskip]() {
								goto l606
							}
							depth--
							add(ruleSTRLEN, position607)
						}
						goto l553
					l606:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position621 := position
							depth++
							{
								position622, tokenIndex622, depth622 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l623
								}
								position++
								goto l622
							l623:
								position, tokenIndex, depth = position622, tokenIndex622, depth622
								if buffer[position] != rune('M') {
									goto l620
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
									goto l620
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
									goto l620
								}
								position++
							}
						l626:
							{
								position628, tokenIndex628, depth628 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l629
								}
								position++
								goto l628
							l629:
								position, tokenIndex, depth = position628, tokenIndex628, depth628
								if buffer[position] != rune('T') {
									goto l620
								}
								position++
							}
						l628:
							{
								position630, tokenIndex630, depth630 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l631
								}
								position++
								goto l630
							l631:
								position, tokenIndex, depth = position630, tokenIndex630, depth630
								if buffer[position] != rune('H') {
									goto l620
								}
								position++
							}
						l630:
							if !rules[ruleskip]() {
								goto l620
							}
							depth--
							add(ruleMONTH, position621)
						}
						goto l553
					l620:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position633 := position
							depth++
							{
								position634, tokenIndex634, depth634 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l635
								}
								position++
								goto l634
							l635:
								position, tokenIndex, depth = position634, tokenIndex634, depth634
								if buffer[position] != rune('M') {
									goto l632
								}
								position++
							}
						l634:
							{
								position636, tokenIndex636, depth636 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l637
								}
								position++
								goto l636
							l637:
								position, tokenIndex, depth = position636, tokenIndex636, depth636
								if buffer[position] != rune('I') {
									goto l632
								}
								position++
							}
						l636:
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
									goto l632
								}
								position++
							}
						l638:
							{
								position640, tokenIndex640, depth640 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l641
								}
								position++
								goto l640
							l641:
								position, tokenIndex, depth = position640, tokenIndex640, depth640
								if buffer[position] != rune('U') {
									goto l632
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
									goto l632
								}
								position++
							}
						l642:
							{
								position644, tokenIndex644, depth644 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l645
								}
								position++
								goto l644
							l645:
								position, tokenIndex, depth = position644, tokenIndex644, depth644
								if buffer[position] != rune('E') {
									goto l632
								}
								position++
							}
						l644:
							{
								position646, tokenIndex646, depth646 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l647
								}
								position++
								goto l646
							l647:
								position, tokenIndex, depth = position646, tokenIndex646, depth646
								if buffer[position] != rune('S') {
									goto l632
								}
								position++
							}
						l646:
							if !rules[ruleskip]() {
								goto l632
							}
							depth--
							add(ruleMINUTES, position633)
						}
						goto l553
					l632:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position649 := position
							depth++
							{
								position650, tokenIndex650, depth650 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l651
								}
								position++
								goto l650
							l651:
								position, tokenIndex, depth = position650, tokenIndex650, depth650
								if buffer[position] != rune('S') {
									goto l648
								}
								position++
							}
						l650:
							{
								position652, tokenIndex652, depth652 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l653
								}
								position++
								goto l652
							l653:
								position, tokenIndex, depth = position652, tokenIndex652, depth652
								if buffer[position] != rune('E') {
									goto l648
								}
								position++
							}
						l652:
							{
								position654, tokenIndex654, depth654 := position, tokenIndex, depth
								if buffer[position] != rune('c') {
									goto l655
								}
								position++
								goto l654
							l655:
								position, tokenIndex, depth = position654, tokenIndex654, depth654
								if buffer[position] != rune('C') {
									goto l648
								}
								position++
							}
						l654:
							{
								position656, tokenIndex656, depth656 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l657
								}
								position++
								goto l656
							l657:
								position, tokenIndex, depth = position656, tokenIndex656, depth656
								if buffer[position] != rune('O') {
									goto l648
								}
								position++
							}
						l656:
							{
								position658, tokenIndex658, depth658 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l659
								}
								position++
								goto l658
							l659:
								position, tokenIndex, depth = position658, tokenIndex658, depth658
								if buffer[position] != rune('N') {
									goto l648
								}
								position++
							}
						l658:
							{
								position660, tokenIndex660, depth660 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l661
								}
								position++
								goto l660
							l661:
								position, tokenIndex, depth = position660, tokenIndex660, depth660
								if buffer[position] != rune('D') {
									goto l648
								}
								position++
							}
						l660:
							{
								position662, tokenIndex662, depth662 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l663
								}
								position++
								goto l662
							l663:
								position, tokenIndex, depth = position662, tokenIndex662, depth662
								if buffer[position] != rune('S') {
									goto l648
								}
								position++
							}
						l662:
							if !rules[ruleskip]() {
								goto l648
							}
							depth--
							add(ruleSECONDS, position649)
						}
						goto l553
					l648:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position665 := position
							depth++
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
									goto l664
								}
								position++
							}
						l666:
							{
								position668, tokenIndex668, depth668 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l669
								}
								position++
								goto l668
							l669:
								position, tokenIndex, depth = position668, tokenIndex668, depth668
								if buffer[position] != rune('I') {
									goto l664
								}
								position++
							}
						l668:
							{
								position670, tokenIndex670, depth670 := position, tokenIndex, depth
								if buffer[position] != rune('m') {
									goto l671
								}
								position++
								goto l670
							l671:
								position, tokenIndex, depth = position670, tokenIndex670, depth670
								if buffer[position] != rune('M') {
									goto l664
								}
								position++
							}
						l670:
							{
								position672, tokenIndex672, depth672 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l673
								}
								position++
								goto l672
							l673:
								position, tokenIndex, depth = position672, tokenIndex672, depth672
								if buffer[position] != rune('E') {
									goto l664
								}
								position++
							}
						l672:
							{
								position674, tokenIndex674, depth674 := position, tokenIndex, depth
								if buffer[position] != rune('z') {
									goto l675
								}
								position++
								goto l674
							l675:
								position, tokenIndex, depth = position674, tokenIndex674, depth674
								if buffer[position] != rune('Z') {
									goto l664
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
									goto l664
								}
								position++
							}
						l676:
							{
								position678, tokenIndex678, depth678 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l679
								}
								position++
								goto l678
							l679:
								position, tokenIndex, depth = position678, tokenIndex678, depth678
								if buffer[position] != rune('N') {
									goto l664
								}
								position++
							}
						l678:
							{
								position680, tokenIndex680, depth680 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l681
								}
								position++
								goto l680
							l681:
								position, tokenIndex, depth = position680, tokenIndex680, depth680
								if buffer[position] != rune('E') {
									goto l664
								}
								position++
							}
						l680:
							if !rules[ruleskip]() {
								goto l664
							}
							depth--
							add(ruleTIMEZONE, position665)
						}
						goto l553
					l664:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position683 := position
							depth++
							{
								position684, tokenIndex684, depth684 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l685
								}
								position++
								goto l684
							l685:
								position, tokenIndex, depth = position684, tokenIndex684, depth684
								if buffer[position] != rune('S') {
									goto l682
								}
								position++
							}
						l684:
							{
								position686, tokenIndex686, depth686 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l687
								}
								position++
								goto l686
							l687:
								position, tokenIndex, depth = position686, tokenIndex686, depth686
								if buffer[position] != rune('H') {
									goto l682
								}
								position++
							}
						l686:
							{
								position688, tokenIndex688, depth688 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l689
								}
								position++
								goto l688
							l689:
								position, tokenIndex, depth = position688, tokenIndex688, depth688
								if buffer[position] != rune('A') {
									goto l682
								}
								position++
							}
						l688:
							if buffer[position] != rune('1') {
								goto l682
							}
							position++
							if !rules[ruleskip]() {
								goto l682
							}
							depth--
							add(ruleSHA1, position683)
						}
						goto l553
					l682:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position691 := position
							depth++
							{
								position692, tokenIndex692, depth692 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l693
								}
								position++
								goto l692
							l693:
								position, tokenIndex, depth = position692, tokenIndex692, depth692
								if buffer[position] != rune('S') {
									goto l690
								}
								position++
							}
						l692:
							{
								position694, tokenIndex694, depth694 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l695
								}
								position++
								goto l694
							l695:
								position, tokenIndex, depth = position694, tokenIndex694, depth694
								if buffer[position] != rune('H') {
									goto l690
								}
								position++
							}
						l694:
							{
								position696, tokenIndex696, depth696 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l697
								}
								position++
								goto l696
							l697:
								position, tokenIndex, depth = position696, tokenIndex696, depth696
								if buffer[position] != rune('A') {
									goto l690
								}
								position++
							}
						l696:
							if buffer[position] != rune('2') {
								goto l690
							}
							position++
							if buffer[position] != rune('5') {
								goto l690
							}
							position++
							if buffer[position] != rune('6') {
								goto l690
							}
							position++
							if !rules[ruleskip]() {
								goto l690
							}
							depth--
							add(ruleSHA256, position691)
						}
						goto l553
					l690:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position699 := position
							depth++
							{
								position700, tokenIndex700, depth700 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l701
								}
								position++
								goto l700
							l701:
								position, tokenIndex, depth = position700, tokenIndex700, depth700
								if buffer[position] != rune('S') {
									goto l698
								}
								position++
							}
						l700:
							{
								position702, tokenIndex702, depth702 := position, tokenIndex, depth
								if buffer[position] != rune('h') {
									goto l703
								}
								position++
								goto l702
							l703:
								position, tokenIndex, depth = position702, tokenIndex702, depth702
								if buffer[position] != rune('H') {
									goto l698
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
									goto l698
								}
								position++
							}
						l704:
							if buffer[position] != rune('3') {
								goto l698
							}
							position++
							if buffer[position] != rune('8') {
								goto l698
							}
							position++
							if buffer[position] != rune('4') {
								goto l698
							}
							position++
							if !rules[ruleskip]() {
								goto l698
							}
							depth--
							add(ruleSHA384, position699)
						}
						goto l553
					l698:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position707 := position
							depth++
							{
								position708, tokenIndex708, depth708 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l709
								}
								position++
								goto l708
							l709:
								position, tokenIndex, depth = position708, tokenIndex708, depth708
								if buffer[position] != rune('I') {
									goto l706
								}
								position++
							}
						l708:
							{
								position710, tokenIndex710, depth710 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l711
								}
								position++
								goto l710
							l711:
								position, tokenIndex, depth = position710, tokenIndex710, depth710
								if buffer[position] != rune('S') {
									goto l706
								}
								position++
							}
						l710:
							{
								position712, tokenIndex712, depth712 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l713
								}
								position++
								goto l712
							l713:
								position, tokenIndex, depth = position712, tokenIndex712, depth712
								if buffer[position] != rune('I') {
									goto l706
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
									goto l706
								}
								position++
							}
						l714:
							{
								position716, tokenIndex716, depth716 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l717
								}
								position++
								goto l716
							l717:
								position, tokenIndex, depth = position716, tokenIndex716, depth716
								if buffer[position] != rune('I') {
									goto l706
								}
								position++
							}
						l716:
							if !rules[ruleskip]() {
								goto l706
							}
							depth--
							add(ruleISIRI, position707)
						}
						goto l553
					l706:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position719 := position
							depth++
							{
								position720, tokenIndex720, depth720 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l721
								}
								position++
								goto l720
							l721:
								position, tokenIndex, depth = position720, tokenIndex720, depth720
								if buffer[position] != rune('I') {
									goto l718
								}
								position++
							}
						l720:
							{
								position722, tokenIndex722, depth722 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l723
								}
								position++
								goto l722
							l723:
								position, tokenIndex, depth = position722, tokenIndex722, depth722
								if buffer[position] != rune('S') {
									goto l718
								}
								position++
							}
						l722:
							{
								position724, tokenIndex724, depth724 := position, tokenIndex, depth
								if buffer[position] != rune('u') {
									goto l725
								}
								position++
								goto l724
							l725:
								position, tokenIndex, depth = position724, tokenIndex724, depth724
								if buffer[position] != rune('U') {
									goto l718
								}
								position++
							}
						l724:
							{
								position726, tokenIndex726, depth726 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l727
								}
								position++
								goto l726
							l727:
								position, tokenIndex, depth = position726, tokenIndex726, depth726
								if buffer[position] != rune('R') {
									goto l718
								}
								position++
							}
						l726:
							{
								position728, tokenIndex728, depth728 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l729
								}
								position++
								goto l728
							l729:
								position, tokenIndex, depth = position728, tokenIndex728, depth728
								if buffer[position] != rune('I') {
									goto l718
								}
								position++
							}
						l728:
							if !rules[ruleskip]() {
								goto l718
							}
							depth--
							add(ruleISURI, position719)
						}
						goto l553
					l718:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position731 := position
							depth++
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
									goto l730
								}
								position++
							}
						l732:
							{
								position734, tokenIndex734, depth734 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l735
								}
								position++
								goto l734
							l735:
								position, tokenIndex, depth = position734, tokenIndex734, depth734
								if buffer[position] != rune('S') {
									goto l730
								}
								position++
							}
						l734:
							{
								position736, tokenIndex736, depth736 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l737
								}
								position++
								goto l736
							l737:
								position, tokenIndex, depth = position736, tokenIndex736, depth736
								if buffer[position] != rune('B') {
									goto l730
								}
								position++
							}
						l736:
							{
								position738, tokenIndex738, depth738 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l739
								}
								position++
								goto l738
							l739:
								position, tokenIndex, depth = position738, tokenIndex738, depth738
								if buffer[position] != rune('L') {
									goto l730
								}
								position++
							}
						l738:
							{
								position740, tokenIndex740, depth740 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l741
								}
								position++
								goto l740
							l741:
								position, tokenIndex, depth = position740, tokenIndex740, depth740
								if buffer[position] != rune('A') {
									goto l730
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
									goto l730
								}
								position++
							}
						l742:
							{
								position744, tokenIndex744, depth744 := position, tokenIndex, depth
								if buffer[position] != rune('k') {
									goto l745
								}
								position++
								goto l744
							l745:
								position, tokenIndex, depth = position744, tokenIndex744, depth744
								if buffer[position] != rune('K') {
									goto l730
								}
								position++
							}
						l744:
							if !rules[ruleskip]() {
								goto l730
							}
							depth--
							add(ruleISBLANK, position731)
						}
						goto l553
					l730:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							position747 := position
							depth++
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
									goto l746
								}
								position++
							}
						l748:
							{
								position750, tokenIndex750, depth750 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l751
								}
								position++
								goto l750
							l751:
								position, tokenIndex, depth = position750, tokenIndex750, depth750
								if buffer[position] != rune('S') {
									goto l746
								}
								position++
							}
						l750:
							{
								position752, tokenIndex752, depth752 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l753
								}
								position++
								goto l752
							l753:
								position, tokenIndex, depth = position752, tokenIndex752, depth752
								if buffer[position] != rune('L') {
									goto l746
								}
								position++
							}
						l752:
							{
								position754, tokenIndex754, depth754 := position, tokenIndex, depth
								if buffer[position] != rune('i') {
									goto l755
								}
								position++
								goto l754
							l755:
								position, tokenIndex, depth = position754, tokenIndex754, depth754
								if buffer[position] != rune('I') {
									goto l746
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
									goto l746
								}
								position++
							}
						l756:
							{
								position758, tokenIndex758, depth758 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l759
								}
								position++
								goto l758
							l759:
								position, tokenIndex, depth = position758, tokenIndex758, depth758
								if buffer[position] != rune('E') {
									goto l746
								}
								position++
							}
						l758:
							{
								position760, tokenIndex760, depth760 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l761
								}
								position++
								goto l760
							l761:
								position, tokenIndex, depth = position760, tokenIndex760, depth760
								if buffer[position] != rune('R') {
									goto l746
								}
								position++
							}
						l760:
							{
								position762, tokenIndex762, depth762 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l763
								}
								position++
								goto l762
							l763:
								position, tokenIndex, depth = position762, tokenIndex762, depth762
								if buffer[position] != rune('A') {
									goto l746
								}
								position++
							}
						l762:
							{
								position764, tokenIndex764, depth764 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l765
								}
								position++
								goto l764
							l765:
								position, tokenIndex, depth = position764, tokenIndex764, depth764
								if buffer[position] != rune('L') {
									goto l746
								}
								position++
							}
						l764:
							if !rules[ruleskip]() {
								goto l746
							}
							depth--
							add(ruleISLITERAL, position747)
						}
						goto l553
					l746:
						position, tokenIndex, depth = position553, tokenIndex553, depth553
						{
							switch buffer[position] {
							case 'I', 'i':
								{
									position767 := position
									depth++
									{
										position768, tokenIndex768, depth768 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l769
										}
										position++
										goto l768
									l769:
										position, tokenIndex, depth = position768, tokenIndex768, depth768
										if buffer[position] != rune('I') {
											goto l552
										}
										position++
									}
								l768:
									{
										position770, tokenIndex770, depth770 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l771
										}
										position++
										goto l770
									l771:
										position, tokenIndex, depth = position770, tokenIndex770, depth770
										if buffer[position] != rune('S') {
											goto l552
										}
										position++
									}
								l770:
									{
										position772, tokenIndex772, depth772 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l773
										}
										position++
										goto l772
									l773:
										position, tokenIndex, depth = position772, tokenIndex772, depth772
										if buffer[position] != rune('N') {
											goto l552
										}
										position++
									}
								l772:
									{
										position774, tokenIndex774, depth774 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l775
										}
										position++
										goto l774
									l775:
										position, tokenIndex, depth = position774, tokenIndex774, depth774
										if buffer[position] != rune('U') {
											goto l552
										}
										position++
									}
								l774:
									{
										position776, tokenIndex776, depth776 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l777
										}
										position++
										goto l776
									l777:
										position, tokenIndex, depth = position776, tokenIndex776, depth776
										if buffer[position] != rune('M') {
											goto l552
										}
										position++
									}
								l776:
									{
										position778, tokenIndex778, depth778 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l779
										}
										position++
										goto l778
									l779:
										position, tokenIndex, depth = position778, tokenIndex778, depth778
										if buffer[position] != rune('E') {
											goto l552
										}
										position++
									}
								l778:
									{
										position780, tokenIndex780, depth780 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l781
										}
										position++
										goto l780
									l781:
										position, tokenIndex, depth = position780, tokenIndex780, depth780
										if buffer[position] != rune('R') {
											goto l552
										}
										position++
									}
								l780:
									{
										position782, tokenIndex782, depth782 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l783
										}
										position++
										goto l782
									l783:
										position, tokenIndex, depth = position782, tokenIndex782, depth782
										if buffer[position] != rune('I') {
											goto l552
										}
										position++
									}
								l782:
									{
										position784, tokenIndex784, depth784 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l785
										}
										position++
										goto l784
									l785:
										position, tokenIndex, depth = position784, tokenIndex784, depth784
										if buffer[position] != rune('C') {
											goto l552
										}
										position++
									}
								l784:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleISNUMERIC, position767)
								}
								break
							case 'S', 's':
								{
									position786 := position
									depth++
									{
										position787, tokenIndex787, depth787 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l788
										}
										position++
										goto l787
									l788:
										position, tokenIndex, depth = position787, tokenIndex787, depth787
										if buffer[position] != rune('S') {
											goto l552
										}
										position++
									}
								l787:
									{
										position789, tokenIndex789, depth789 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l790
										}
										position++
										goto l789
									l790:
										position, tokenIndex, depth = position789, tokenIndex789, depth789
										if buffer[position] != rune('H') {
											goto l552
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
											goto l552
										}
										position++
									}
								l791:
									if buffer[position] != rune('5') {
										goto l552
									}
									position++
									if buffer[position] != rune('1') {
										goto l552
									}
									position++
									if buffer[position] != rune('2') {
										goto l552
									}
									position++
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleSHA512, position786)
								}
								break
							case 'M', 'm':
								{
									position793 := position
									depth++
									{
										position794, tokenIndex794, depth794 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l795
										}
										position++
										goto l794
									l795:
										position, tokenIndex, depth = position794, tokenIndex794, depth794
										if buffer[position] != rune('M') {
											goto l552
										}
										position++
									}
								l794:
									{
										position796, tokenIndex796, depth796 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l797
										}
										position++
										goto l796
									l797:
										position, tokenIndex, depth = position796, tokenIndex796, depth796
										if buffer[position] != rune('D') {
											goto l552
										}
										position++
									}
								l796:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleMD5, position793)
								}
								break
							case 'T', 't':
								{
									position798 := position
									depth++
									{
										position799, tokenIndex799, depth799 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l800
										}
										position++
										goto l799
									l800:
										position, tokenIndex, depth = position799, tokenIndex799, depth799
										if buffer[position] != rune('T') {
											goto l552
										}
										position++
									}
								l799:
									{
										position801, tokenIndex801, depth801 := position, tokenIndex, depth
										if buffer[position] != rune('z') {
											goto l802
										}
										position++
										goto l801
									l802:
										position, tokenIndex, depth = position801, tokenIndex801, depth801
										if buffer[position] != rune('Z') {
											goto l552
										}
										position++
									}
								l801:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleTZ, position798)
								}
								break
							case 'H', 'h':
								{
									position803 := position
									depth++
									{
										position804, tokenIndex804, depth804 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l805
										}
										position++
										goto l804
									l805:
										position, tokenIndex, depth = position804, tokenIndex804, depth804
										if buffer[position] != rune('H') {
											goto l552
										}
										position++
									}
								l804:
									{
										position806, tokenIndex806, depth806 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l807
										}
										position++
										goto l806
									l807:
										position, tokenIndex, depth = position806, tokenIndex806, depth806
										if buffer[position] != rune('O') {
											goto l552
										}
										position++
									}
								l806:
									{
										position808, tokenIndex808, depth808 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l809
										}
										position++
										goto l808
									l809:
										position, tokenIndex, depth = position808, tokenIndex808, depth808
										if buffer[position] != rune('U') {
											goto l552
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
											goto l552
										}
										position++
									}
								l810:
									{
										position812, tokenIndex812, depth812 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l813
										}
										position++
										goto l812
									l813:
										position, tokenIndex, depth = position812, tokenIndex812, depth812
										if buffer[position] != rune('S') {
											goto l552
										}
										position++
									}
								l812:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleHOURS, position803)
								}
								break
							case 'D', 'd':
								{
									position814 := position
									depth++
									{
										position815, tokenIndex815, depth815 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l816
										}
										position++
										goto l815
									l816:
										position, tokenIndex, depth = position815, tokenIndex815, depth815
										if buffer[position] != rune('D') {
											goto l552
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
											goto l552
										}
										position++
									}
								l817:
									{
										position819, tokenIndex819, depth819 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l820
										}
										position++
										goto l819
									l820:
										position, tokenIndex, depth = position819, tokenIndex819, depth819
										if buffer[position] != rune('Y') {
											goto l552
										}
										position++
									}
								l819:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleDAY, position814)
								}
								break
							case 'Y', 'y':
								{
									position821 := position
									depth++
									{
										position822, tokenIndex822, depth822 := position, tokenIndex, depth
										if buffer[position] != rune('y') {
											goto l823
										}
										position++
										goto l822
									l823:
										position, tokenIndex, depth = position822, tokenIndex822, depth822
										if buffer[position] != rune('Y') {
											goto l552
										}
										position++
									}
								l822:
									{
										position824, tokenIndex824, depth824 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l825
										}
										position++
										goto l824
									l825:
										position, tokenIndex, depth = position824, tokenIndex824, depth824
										if buffer[position] != rune('E') {
											goto l552
										}
										position++
									}
								l824:
									{
										position826, tokenIndex826, depth826 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l827
										}
										position++
										goto l826
									l827:
										position, tokenIndex, depth = position826, tokenIndex826, depth826
										if buffer[position] != rune('A') {
											goto l552
										}
										position++
									}
								l826:
									{
										position828, tokenIndex828, depth828 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l829
										}
										position++
										goto l828
									l829:
										position, tokenIndex, depth = position828, tokenIndex828, depth828
										if buffer[position] != rune('R') {
											goto l552
										}
										position++
									}
								l828:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleYEAR, position821)
								}
								break
							case 'E', 'e':
								{
									position830 := position
									depth++
									{
										position831, tokenIndex831, depth831 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l832
										}
										position++
										goto l831
									l832:
										position, tokenIndex, depth = position831, tokenIndex831, depth831
										if buffer[position] != rune('E') {
											goto l552
										}
										position++
									}
								l831:
									{
										position833, tokenIndex833, depth833 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l834
										}
										position++
										goto l833
									l834:
										position, tokenIndex, depth = position833, tokenIndex833, depth833
										if buffer[position] != rune('N') {
											goto l552
										}
										position++
									}
								l833:
									{
										position835, tokenIndex835, depth835 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l836
										}
										position++
										goto l835
									l836:
										position, tokenIndex, depth = position835, tokenIndex835, depth835
										if buffer[position] != rune('C') {
											goto l552
										}
										position++
									}
								l835:
									{
										position837, tokenIndex837, depth837 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l838
										}
										position++
										goto l837
									l838:
										position, tokenIndex, depth = position837, tokenIndex837, depth837
										if buffer[position] != rune('O') {
											goto l552
										}
										position++
									}
								l837:
									{
										position839, tokenIndex839, depth839 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l840
										}
										position++
										goto l839
									l840:
										position, tokenIndex, depth = position839, tokenIndex839, depth839
										if buffer[position] != rune('D') {
											goto l552
										}
										position++
									}
								l839:
									{
										position841, tokenIndex841, depth841 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l842
										}
										position++
										goto l841
									l842:
										position, tokenIndex, depth = position841, tokenIndex841, depth841
										if buffer[position] != rune('E') {
											goto l552
										}
										position++
									}
								l841:
									if buffer[position] != rune('_') {
										goto l552
									}
									position++
									{
										position843, tokenIndex843, depth843 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l844
										}
										position++
										goto l843
									l844:
										position, tokenIndex, depth = position843, tokenIndex843, depth843
										if buffer[position] != rune('F') {
											goto l552
										}
										position++
									}
								l843:
									{
										position845, tokenIndex845, depth845 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l846
										}
										position++
										goto l845
									l846:
										position, tokenIndex, depth = position845, tokenIndex845, depth845
										if buffer[position] != rune('O') {
											goto l552
										}
										position++
									}
								l845:
									{
										position847, tokenIndex847, depth847 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l848
										}
										position++
										goto l847
									l848:
										position, tokenIndex, depth = position847, tokenIndex847, depth847
										if buffer[position] != rune('R') {
											goto l552
										}
										position++
									}
								l847:
									if buffer[position] != rune('_') {
										goto l552
									}
									position++
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
											goto l552
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
											goto l552
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
											goto l552
										}
										position++
									}
								l853:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleENCODEFORURI, position830)
								}
								break
							case 'L', 'l':
								{
									position855 := position
									depth++
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
											goto l552
										}
										position++
									}
								l856:
									{
										position858, tokenIndex858, depth858 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l859
										}
										position++
										goto l858
									l859:
										position, tokenIndex, depth = position858, tokenIndex858, depth858
										if buffer[position] != rune('C') {
											goto l552
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
											goto l552
										}
										position++
									}
								l860:
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
											goto l552
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
											goto l552
										}
										position++
									}
								l864:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleLCASE, position855)
								}
								break
							case 'U', 'u':
								{
									position866 := position
									depth++
									{
										position867, tokenIndex867, depth867 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l868
										}
										position++
										goto l867
									l868:
										position, tokenIndex, depth = position867, tokenIndex867, depth867
										if buffer[position] != rune('U') {
											goto l552
										}
										position++
									}
								l867:
									{
										position869, tokenIndex869, depth869 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l870
										}
										position++
										goto l869
									l870:
										position, tokenIndex, depth = position869, tokenIndex869, depth869
										if buffer[position] != rune('C') {
											goto l552
										}
										position++
									}
								l869:
									{
										position871, tokenIndex871, depth871 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l872
										}
										position++
										goto l871
									l872:
										position, tokenIndex, depth = position871, tokenIndex871, depth871
										if buffer[position] != rune('A') {
											goto l552
										}
										position++
									}
								l871:
									{
										position873, tokenIndex873, depth873 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l874
										}
										position++
										goto l873
									l874:
										position, tokenIndex, depth = position873, tokenIndex873, depth873
										if buffer[position] != rune('S') {
											goto l552
										}
										position++
									}
								l873:
									{
										position875, tokenIndex875, depth875 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l876
										}
										position++
										goto l875
									l876:
										position, tokenIndex, depth = position875, tokenIndex875, depth875
										if buffer[position] != rune('E') {
											goto l552
										}
										position++
									}
								l875:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleUCASE, position866)
								}
								break
							case 'F', 'f':
								{
									position877 := position
									depth++
									{
										position878, tokenIndex878, depth878 := position, tokenIndex, depth
										if buffer[position] != rune('f') {
											goto l879
										}
										position++
										goto l878
									l879:
										position, tokenIndex, depth = position878, tokenIndex878, depth878
										if buffer[position] != rune('F') {
											goto l552
										}
										position++
									}
								l878:
									{
										position880, tokenIndex880, depth880 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l881
										}
										position++
										goto l880
									l881:
										position, tokenIndex, depth = position880, tokenIndex880, depth880
										if buffer[position] != rune('L') {
											goto l552
										}
										position++
									}
								l880:
									{
										position882, tokenIndex882, depth882 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l883
										}
										position++
										goto l882
									l883:
										position, tokenIndex, depth = position882, tokenIndex882, depth882
										if buffer[position] != rune('O') {
											goto l552
										}
										position++
									}
								l882:
									{
										position884, tokenIndex884, depth884 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l885
										}
										position++
										goto l884
									l885:
										position, tokenIndex, depth = position884, tokenIndex884, depth884
										if buffer[position] != rune('O') {
											goto l552
										}
										position++
									}
								l884:
									{
										position886, tokenIndex886, depth886 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l887
										}
										position++
										goto l886
									l887:
										position, tokenIndex, depth = position886, tokenIndex886, depth886
										if buffer[position] != rune('R') {
											goto l552
										}
										position++
									}
								l886:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleFLOOR, position877)
								}
								break
							case 'R', 'r':
								{
									position888 := position
									depth++
									{
										position889, tokenIndex889, depth889 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l890
										}
										position++
										goto l889
									l890:
										position, tokenIndex, depth = position889, tokenIndex889, depth889
										if buffer[position] != rune('R') {
											goto l552
										}
										position++
									}
								l889:
									{
										position891, tokenIndex891, depth891 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l892
										}
										position++
										goto l891
									l892:
										position, tokenIndex, depth = position891, tokenIndex891, depth891
										if buffer[position] != rune('O') {
											goto l552
										}
										position++
									}
								l891:
									{
										position893, tokenIndex893, depth893 := position, tokenIndex, depth
										if buffer[position] != rune('u') {
											goto l894
										}
										position++
										goto l893
									l894:
										position, tokenIndex, depth = position893, tokenIndex893, depth893
										if buffer[position] != rune('U') {
											goto l552
										}
										position++
									}
								l893:
									{
										position895, tokenIndex895, depth895 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l896
										}
										position++
										goto l895
									l896:
										position, tokenIndex, depth = position895, tokenIndex895, depth895
										if buffer[position] != rune('N') {
											goto l552
										}
										position++
									}
								l895:
									{
										position897, tokenIndex897, depth897 := position, tokenIndex, depth
										if buffer[position] != rune('d') {
											goto l898
										}
										position++
										goto l897
									l898:
										position, tokenIndex, depth = position897, tokenIndex897, depth897
										if buffer[position] != rune('D') {
											goto l552
										}
										position++
									}
								l897:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleROUND, position888)
								}
								break
							case 'C', 'c':
								{
									position899 := position
									depth++
									{
										position900, tokenIndex900, depth900 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l901
										}
										position++
										goto l900
									l901:
										position, tokenIndex, depth = position900, tokenIndex900, depth900
										if buffer[position] != rune('C') {
											goto l552
										}
										position++
									}
								l900:
									{
										position902, tokenIndex902, depth902 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l903
										}
										position++
										goto l902
									l903:
										position, tokenIndex, depth = position902, tokenIndex902, depth902
										if buffer[position] != rune('E') {
											goto l552
										}
										position++
									}
								l902:
									{
										position904, tokenIndex904, depth904 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l905
										}
										position++
										goto l904
									l905:
										position, tokenIndex, depth = position904, tokenIndex904, depth904
										if buffer[position] != rune('I') {
											goto l552
										}
										position++
									}
								l904:
									{
										position906, tokenIndex906, depth906 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l907
										}
										position++
										goto l906
									l907:
										position, tokenIndex, depth = position906, tokenIndex906, depth906
										if buffer[position] != rune('L') {
											goto l552
										}
										position++
									}
								l906:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleCEIL, position899)
								}
								break
							default:
								{
									position908 := position
									depth++
									{
										position909, tokenIndex909, depth909 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l910
										}
										position++
										goto l909
									l910:
										position, tokenIndex, depth = position909, tokenIndex909, depth909
										if buffer[position] != rune('A') {
											goto l552
										}
										position++
									}
								l909:
									{
										position911, tokenIndex911, depth911 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l912
										}
										position++
										goto l911
									l912:
										position, tokenIndex, depth = position911, tokenIndex911, depth911
										if buffer[position] != rune('B') {
											goto l552
										}
										position++
									}
								l911:
									{
										position913, tokenIndex913, depth913 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l914
										}
										position++
										goto l913
									l914:
										position, tokenIndex, depth = position913, tokenIndex913, depth913
										if buffer[position] != rune('S') {
											goto l552
										}
										position++
									}
								l913:
									if !rules[ruleskip]() {
										goto l552
									}
									depth--
									add(ruleABS, position908)
								}
								break
							}
						}

					}
				l553:
					if !rules[ruleLPAREN]() {
						goto l552
					}
					if !rules[ruleexpression]() {
						goto l552
					}
					if !rules[ruleRPAREN]() {
						goto l552
					}
					goto l551
				l552:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					{
						position916, tokenIndex916, depth916 := position, tokenIndex, depth
						{
							position918 := position
							depth++
							{
								position919, tokenIndex919, depth919 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l920
								}
								position++
								goto l919
							l920:
								position, tokenIndex, depth = position919, tokenIndex919, depth919
								if buffer[position] != rune('S') {
									goto l917
								}
								position++
							}
						l919:
							{
								position921, tokenIndex921, depth921 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l922
								}
								position++
								goto l921
							l922:
								position, tokenIndex, depth = position921, tokenIndex921, depth921
								if buffer[position] != rune('T') {
									goto l917
								}
								position++
							}
						l921:
							{
								position923, tokenIndex923, depth923 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l924
								}
								position++
								goto l923
							l924:
								position, tokenIndex, depth = position923, tokenIndex923, depth923
								if buffer[position] != rune('R') {
									goto l917
								}
								position++
							}
						l923:
							{
								position925, tokenIndex925, depth925 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l926
								}
								position++
								goto l925
							l926:
								position, tokenIndex, depth = position925, tokenIndex925, depth925
								if buffer[position] != rune('S') {
									goto l917
								}
								position++
							}
						l925:
							{
								position927, tokenIndex927, depth927 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l928
								}
								position++
								goto l927
							l928:
								position, tokenIndex, depth = position927, tokenIndex927, depth927
								if buffer[position] != rune('T') {
									goto l917
								}
								position++
							}
						l927:
							{
								position929, tokenIndex929, depth929 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l930
								}
								position++
								goto l929
							l930:
								position, tokenIndex, depth = position929, tokenIndex929, depth929
								if buffer[position] != rune('A') {
									goto l917
								}
								position++
							}
						l929:
							{
								position931, tokenIndex931, depth931 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l932
								}
								position++
								goto l931
							l932:
								position, tokenIndex, depth = position931, tokenIndex931, depth931
								if buffer[position] != rune('R') {
									goto l917
								}
								position++
							}
						l931:
							{
								position933, tokenIndex933, depth933 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l934
								}
								position++
								goto l933
							l934:
								position, tokenIndex, depth = position933, tokenIndex933, depth933
								if buffer[position] != rune('T') {
									goto l917
								}
								position++
							}
						l933:
							{
								position935, tokenIndex935, depth935 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l936
								}
								position++
								goto l935
							l936:
								position, tokenIndex, depth = position935, tokenIndex935, depth935
								if buffer[position] != rune('S') {
									goto l917
								}
								position++
							}
						l935:
							if !rules[ruleskip]() {
								goto l917
							}
							depth--
							add(ruleSTRSTARTS, position918)
						}
						goto l916
					l917:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						{
							position938 := position
							depth++
							{
								position939, tokenIndex939, depth939 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l940
								}
								position++
								goto l939
							l940:
								position, tokenIndex, depth = position939, tokenIndex939, depth939
								if buffer[position] != rune('S') {
									goto l937
								}
								position++
							}
						l939:
							{
								position941, tokenIndex941, depth941 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l942
								}
								position++
								goto l941
							l942:
								position, tokenIndex, depth = position941, tokenIndex941, depth941
								if buffer[position] != rune('T') {
									goto l937
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
									goto l937
								}
								position++
							}
						l943:
							{
								position945, tokenIndex945, depth945 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l946
								}
								position++
								goto l945
							l946:
								position, tokenIndex, depth = position945, tokenIndex945, depth945
								if buffer[position] != rune('E') {
									goto l937
								}
								position++
							}
						l945:
							{
								position947, tokenIndex947, depth947 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l948
								}
								position++
								goto l947
							l948:
								position, tokenIndex, depth = position947, tokenIndex947, depth947
								if buffer[position] != rune('N') {
									goto l937
								}
								position++
							}
						l947:
							{
								position949, tokenIndex949, depth949 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l950
								}
								position++
								goto l949
							l950:
								position, tokenIndex, depth = position949, tokenIndex949, depth949
								if buffer[position] != rune('D') {
									goto l937
								}
								position++
							}
						l949:
							{
								position951, tokenIndex951, depth951 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l952
								}
								position++
								goto l951
							l952:
								position, tokenIndex, depth = position951, tokenIndex951, depth951
								if buffer[position] != rune('S') {
									goto l937
								}
								position++
							}
						l951:
							if !rules[ruleskip]() {
								goto l937
							}
							depth--
							add(ruleSTRENDS, position938)
						}
						goto l916
					l937:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						{
							position954 := position
							depth++
							{
								position955, tokenIndex955, depth955 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l956
								}
								position++
								goto l955
							l956:
								position, tokenIndex, depth = position955, tokenIndex955, depth955
								if buffer[position] != rune('S') {
									goto l953
								}
								position++
							}
						l955:
							{
								position957, tokenIndex957, depth957 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l958
								}
								position++
								goto l957
							l958:
								position, tokenIndex, depth = position957, tokenIndex957, depth957
								if buffer[position] != rune('T') {
									goto l953
								}
								position++
							}
						l957:
							{
								position959, tokenIndex959, depth959 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l960
								}
								position++
								goto l959
							l960:
								position, tokenIndex, depth = position959, tokenIndex959, depth959
								if buffer[position] != rune('R') {
									goto l953
								}
								position++
							}
						l959:
							{
								position961, tokenIndex961, depth961 := position, tokenIndex, depth
								if buffer[position] != rune('b') {
									goto l962
								}
								position++
								goto l961
							l962:
								position, tokenIndex, depth = position961, tokenIndex961, depth961
								if buffer[position] != rune('B') {
									goto l953
								}
								position++
							}
						l961:
							{
								position963, tokenIndex963, depth963 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l964
								}
								position++
								goto l963
							l964:
								position, tokenIndex, depth = position963, tokenIndex963, depth963
								if buffer[position] != rune('E') {
									goto l953
								}
								position++
							}
						l963:
							{
								position965, tokenIndex965, depth965 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l966
								}
								position++
								goto l965
							l966:
								position, tokenIndex, depth = position965, tokenIndex965, depth965
								if buffer[position] != rune('F') {
									goto l953
								}
								position++
							}
						l965:
							{
								position967, tokenIndex967, depth967 := position, tokenIndex, depth
								if buffer[position] != rune('o') {
									goto l968
								}
								position++
								goto l967
							l968:
								position, tokenIndex, depth = position967, tokenIndex967, depth967
								if buffer[position] != rune('O') {
									goto l953
								}
								position++
							}
						l967:
							{
								position969, tokenIndex969, depth969 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l970
								}
								position++
								goto l969
							l970:
								position, tokenIndex, depth = position969, tokenIndex969, depth969
								if buffer[position] != rune('R') {
									goto l953
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
									goto l953
								}
								position++
							}
						l971:
							if !rules[ruleskip]() {
								goto l953
							}
							depth--
							add(ruleSTRBEFORE, position954)
						}
						goto l916
					l953:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						{
							position974 := position
							depth++
							{
								position975, tokenIndex975, depth975 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l976
								}
								position++
								goto l975
							l976:
								position, tokenIndex, depth = position975, tokenIndex975, depth975
								if buffer[position] != rune('S') {
									goto l973
								}
								position++
							}
						l975:
							{
								position977, tokenIndex977, depth977 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l978
								}
								position++
								goto l977
							l978:
								position, tokenIndex, depth = position977, tokenIndex977, depth977
								if buffer[position] != rune('T') {
									goto l973
								}
								position++
							}
						l977:
							{
								position979, tokenIndex979, depth979 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l980
								}
								position++
								goto l979
							l980:
								position, tokenIndex, depth = position979, tokenIndex979, depth979
								if buffer[position] != rune('R') {
									goto l973
								}
								position++
							}
						l979:
							{
								position981, tokenIndex981, depth981 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l982
								}
								position++
								goto l981
							l982:
								position, tokenIndex, depth = position981, tokenIndex981, depth981
								if buffer[position] != rune('A') {
									goto l973
								}
								position++
							}
						l981:
							{
								position983, tokenIndex983, depth983 := position, tokenIndex, depth
								if buffer[position] != rune('f') {
									goto l984
								}
								position++
								goto l983
							l984:
								position, tokenIndex, depth = position983, tokenIndex983, depth983
								if buffer[position] != rune('F') {
									goto l973
								}
								position++
							}
						l983:
							{
								position985, tokenIndex985, depth985 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l986
								}
								position++
								goto l985
							l986:
								position, tokenIndex, depth = position985, tokenIndex985, depth985
								if buffer[position] != rune('T') {
									goto l973
								}
								position++
							}
						l985:
							{
								position987, tokenIndex987, depth987 := position, tokenIndex, depth
								if buffer[position] != rune('e') {
									goto l988
								}
								position++
								goto l987
							l988:
								position, tokenIndex, depth = position987, tokenIndex987, depth987
								if buffer[position] != rune('E') {
									goto l973
								}
								position++
							}
						l987:
							{
								position989, tokenIndex989, depth989 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l990
								}
								position++
								goto l989
							l990:
								position, tokenIndex, depth = position989, tokenIndex989, depth989
								if buffer[position] != rune('R') {
									goto l973
								}
								position++
							}
						l989:
							if !rules[ruleskip]() {
								goto l973
							}
							depth--
							add(ruleSTRAFTER, position974)
						}
						goto l916
					l973:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						{
							position992 := position
							depth++
							{
								position993, tokenIndex993, depth993 := position, tokenIndex, depth
								if buffer[position] != rune('s') {
									goto l994
								}
								position++
								goto l993
							l994:
								position, tokenIndex, depth = position993, tokenIndex993, depth993
								if buffer[position] != rune('S') {
									goto l991
								}
								position++
							}
						l993:
							{
								position995, tokenIndex995, depth995 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l996
								}
								position++
								goto l995
							l996:
								position, tokenIndex, depth = position995, tokenIndex995, depth995
								if buffer[position] != rune('T') {
									goto l991
								}
								position++
							}
						l995:
							{
								position997, tokenIndex997, depth997 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l998
								}
								position++
								goto l997
							l998:
								position, tokenIndex, depth = position997, tokenIndex997, depth997
								if buffer[position] != rune('R') {
									goto l991
								}
								position++
							}
						l997:
							{
								position999, tokenIndex999, depth999 := position, tokenIndex, depth
								if buffer[position] != rune('l') {
									goto l1000
								}
								position++
								goto l999
							l1000:
								position, tokenIndex, depth = position999, tokenIndex999, depth999
								if buffer[position] != rune('L') {
									goto l991
								}
								position++
							}
						l999:
							{
								position1001, tokenIndex1001, depth1001 := position, tokenIndex, depth
								if buffer[position] != rune('a') {
									goto l1002
								}
								position++
								goto l1001
							l1002:
								position, tokenIndex, depth = position1001, tokenIndex1001, depth1001
								if buffer[position] != rune('A') {
									goto l991
								}
								position++
							}
						l1001:
							{
								position1003, tokenIndex1003, depth1003 := position, tokenIndex, depth
								if buffer[position] != rune('n') {
									goto l1004
								}
								position++
								goto l1003
							l1004:
								position, tokenIndex, depth = position1003, tokenIndex1003, depth1003
								if buffer[position] != rune('N') {
									goto l991
								}
								position++
							}
						l1003:
							{
								position1005, tokenIndex1005, depth1005 := position, tokenIndex, depth
								if buffer[position] != rune('g') {
									goto l1006
								}
								position++
								goto l1005
							l1006:
								position, tokenIndex, depth = position1005, tokenIndex1005, depth1005
								if buffer[position] != rune('G') {
									goto l991
								}
								position++
							}
						l1005:
							if !rules[ruleskip]() {
								goto l991
							}
							depth--
							add(ruleSTRLANG, position992)
						}
						goto l916
					l991:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						{
							position1008 := position
							depth++
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
									goto l1007
								}
								position++
							}
						l1009:
							{
								position1011, tokenIndex1011, depth1011 := position, tokenIndex, depth
								if buffer[position] != rune('t') {
									goto l1012
								}
								position++
								goto l1011
							l1012:
								position, tokenIndex, depth = position1011, tokenIndex1011, depth1011
								if buffer[position] != rune('T') {
									goto l1007
								}
								position++
							}
						l1011:
							{
								position1013, tokenIndex1013, depth1013 := position, tokenIndex, depth
								if buffer[position] != rune('r') {
									goto l1014
								}
								position++
								goto l1013
							l1014:
								position, tokenIndex, depth = position1013, tokenIndex1013, depth1013
								if buffer[position] != rune('R') {
									goto l1007
								}
								position++
							}
						l1013:
							{
								position1015, tokenIndex1015, depth1015 := position, tokenIndex, depth
								if buffer[position] != rune('d') {
									goto l1016
								}
								position++
								goto l1015
							l1016:
								position, tokenIndex, depth = position1015, tokenIndex1015, depth1015
								if buffer[position] != rune('D') {
									goto l1007
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
									goto l1007
								}
								position++
							}
						l1017:
							if !rules[ruleskip]() {
								goto l1007
							}
							depth--
							add(ruleSTRDT, position1008)
						}
						goto l916
					l1007:
						position, tokenIndex, depth = position916, tokenIndex916, depth916
						{
							switch buffer[position] {
							case 'S', 's':
								{
									position1020 := position
									depth++
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
											goto l915
										}
										position++
									}
								l1021:
									{
										position1023, tokenIndex1023, depth1023 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1024
										}
										position++
										goto l1023
									l1024:
										position, tokenIndex, depth = position1023, tokenIndex1023, depth1023
										if buffer[position] != rune('A') {
											goto l915
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
											goto l915
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
											goto l915
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
											goto l915
										}
										position++
									}
								l1029:
									{
										position1031, tokenIndex1031, depth1031 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1032
										}
										position++
										goto l1031
									l1032:
										position, tokenIndex, depth = position1031, tokenIndex1031, depth1031
										if buffer[position] != rune('E') {
											goto l915
										}
										position++
									}
								l1031:
									{
										position1033, tokenIndex1033, depth1033 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1034
										}
										position++
										goto l1033
									l1034:
										position, tokenIndex, depth = position1033, tokenIndex1033, depth1033
										if buffer[position] != rune('R') {
											goto l915
										}
										position++
									}
								l1033:
									{
										position1035, tokenIndex1035, depth1035 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1036
										}
										position++
										goto l1035
									l1036:
										position, tokenIndex, depth = position1035, tokenIndex1035, depth1035
										if buffer[position] != rune('M') {
											goto l915
										}
										position++
									}
								l1035:
									if !rules[ruleskip]() {
										goto l915
									}
									depth--
									add(ruleSAMETERM, position1020)
								}
								break
							case 'C', 'c':
								{
									position1037 := position
									depth++
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
											goto l915
										}
										position++
									}
								l1038:
									{
										position1040, tokenIndex1040, depth1040 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1041
										}
										position++
										goto l1040
									l1041:
										position, tokenIndex, depth = position1040, tokenIndex1040, depth1040
										if buffer[position] != rune('O') {
											goto l915
										}
										position++
									}
								l1040:
									{
										position1042, tokenIndex1042, depth1042 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1043
										}
										position++
										goto l1042
									l1043:
										position, tokenIndex, depth = position1042, tokenIndex1042, depth1042
										if buffer[position] != rune('N') {
											goto l915
										}
										position++
									}
								l1042:
									{
										position1044, tokenIndex1044, depth1044 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1045
										}
										position++
										goto l1044
									l1045:
										position, tokenIndex, depth = position1044, tokenIndex1044, depth1044
										if buffer[position] != rune('T') {
											goto l915
										}
										position++
									}
								l1044:
									{
										position1046, tokenIndex1046, depth1046 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1047
										}
										position++
										goto l1046
									l1047:
										position, tokenIndex, depth = position1046, tokenIndex1046, depth1046
										if buffer[position] != rune('A') {
											goto l915
										}
										position++
									}
								l1046:
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
											goto l915
										}
										position++
									}
								l1048:
									{
										position1050, tokenIndex1050, depth1050 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1051
										}
										position++
										goto l1050
									l1051:
										position, tokenIndex, depth = position1050, tokenIndex1050, depth1050
										if buffer[position] != rune('N') {
											goto l915
										}
										position++
									}
								l1050:
									{
										position1052, tokenIndex1052, depth1052 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1053
										}
										position++
										goto l1052
									l1053:
										position, tokenIndex, depth = position1052, tokenIndex1052, depth1052
										if buffer[position] != rune('S') {
											goto l915
										}
										position++
									}
								l1052:
									if !rules[ruleskip]() {
										goto l915
									}
									depth--
									add(ruleCONTAINS, position1037)
								}
								break
							default:
								{
									position1054 := position
									depth++
									{
										position1055, tokenIndex1055, depth1055 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1056
										}
										position++
										goto l1055
									l1056:
										position, tokenIndex, depth = position1055, tokenIndex1055, depth1055
										if buffer[position] != rune('L') {
											goto l915
										}
										position++
									}
								l1055:
									{
										position1057, tokenIndex1057, depth1057 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1058
										}
										position++
										goto l1057
									l1058:
										position, tokenIndex, depth = position1057, tokenIndex1057, depth1057
										if buffer[position] != rune('A') {
											goto l915
										}
										position++
									}
								l1057:
									{
										position1059, tokenIndex1059, depth1059 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1060
										}
										position++
										goto l1059
									l1060:
										position, tokenIndex, depth = position1059, tokenIndex1059, depth1059
										if buffer[position] != rune('N') {
											goto l915
										}
										position++
									}
								l1059:
									{
										position1061, tokenIndex1061, depth1061 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1062
										}
										position++
										goto l1061
									l1062:
										position, tokenIndex, depth = position1061, tokenIndex1061, depth1061
										if buffer[position] != rune('G') {
											goto l915
										}
										position++
									}
								l1061:
									{
										position1063, tokenIndex1063, depth1063 := position, tokenIndex, depth
										if buffer[position] != rune('m') {
											goto l1064
										}
										position++
										goto l1063
									l1064:
										position, tokenIndex, depth = position1063, tokenIndex1063, depth1063
										if buffer[position] != rune('M') {
											goto l915
										}
										position++
									}
								l1063:
									{
										position1065, tokenIndex1065, depth1065 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1066
										}
										position++
										goto l1065
									l1066:
										position, tokenIndex, depth = position1065, tokenIndex1065, depth1065
										if buffer[position] != rune('A') {
											goto l915
										}
										position++
									}
								l1065:
									{
										position1067, tokenIndex1067, depth1067 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1068
										}
										position++
										goto l1067
									l1068:
										position, tokenIndex, depth = position1067, tokenIndex1067, depth1067
										if buffer[position] != rune('T') {
											goto l915
										}
										position++
									}
								l1067:
									{
										position1069, tokenIndex1069, depth1069 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1070
										}
										position++
										goto l1069
									l1070:
										position, tokenIndex, depth = position1069, tokenIndex1069, depth1069
										if buffer[position] != rune('C') {
											goto l915
										}
										position++
									}
								l1069:
									{
										position1071, tokenIndex1071, depth1071 := position, tokenIndex, depth
										if buffer[position] != rune('h') {
											goto l1072
										}
										position++
										goto l1071
									l1072:
										position, tokenIndex, depth = position1071, tokenIndex1071, depth1071
										if buffer[position] != rune('H') {
											goto l915
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
											goto l915
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
											goto l915
										}
										position++
									}
								l1075:
									if !rules[ruleskip]() {
										goto l915
									}
									depth--
									add(ruleLANGMATCHES, position1054)
								}
								break
							}
						}

					}
				l916:
					if !rules[ruleLPAREN]() {
						goto l915
					}
					if !rules[ruleexpression]() {
						goto l915
					}
					if !rules[ruleCOMMA]() {
						goto l915
					}
					if !rules[ruleexpression]() {
						goto l915
					}
					if !rules[ruleRPAREN]() {
						goto l915
					}
					goto l551
				l915:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					{
						position1078 := position
						depth++
						{
							position1079, tokenIndex1079, depth1079 := position, tokenIndex, depth
							if buffer[position] != rune('b') {
								goto l1080
							}
							position++
							goto l1079
						l1080:
							position, tokenIndex, depth = position1079, tokenIndex1079, depth1079
							if buffer[position] != rune('B') {
								goto l1077
							}
							position++
						}
					l1079:
						{
							position1081, tokenIndex1081, depth1081 := position, tokenIndex, depth
							if buffer[position] != rune('o') {
								goto l1082
							}
							position++
							goto l1081
						l1082:
							position, tokenIndex, depth = position1081, tokenIndex1081, depth1081
							if buffer[position] != rune('O') {
								goto l1077
							}
							position++
						}
					l1081:
						{
							position1083, tokenIndex1083, depth1083 := position, tokenIndex, depth
							if buffer[position] != rune('u') {
								goto l1084
							}
							position++
							goto l1083
						l1084:
							position, tokenIndex, depth = position1083, tokenIndex1083, depth1083
							if buffer[position] != rune('U') {
								goto l1077
							}
							position++
						}
					l1083:
						{
							position1085, tokenIndex1085, depth1085 := position, tokenIndex, depth
							if buffer[position] != rune('n') {
								goto l1086
							}
							position++
							goto l1085
						l1086:
							position, tokenIndex, depth = position1085, tokenIndex1085, depth1085
							if buffer[position] != rune('N') {
								goto l1077
							}
							position++
						}
					l1085:
						{
							position1087, tokenIndex1087, depth1087 := position, tokenIndex, depth
							if buffer[position] != rune('d') {
								goto l1088
							}
							position++
							goto l1087
						l1088:
							position, tokenIndex, depth = position1087, tokenIndex1087, depth1087
							if buffer[position] != rune('D') {
								goto l1077
							}
							position++
						}
					l1087:
						if !rules[ruleskip]() {
							goto l1077
						}
						depth--
						add(ruleBOUND, position1078)
					}
					if !rules[ruleLPAREN]() {
						goto l1077
					}
					if !rules[rulevar]() {
						goto l1077
					}
					if !rules[ruleRPAREN]() {
						goto l1077
					}
					goto l551
				l1077:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					{
						switch buffer[position] {
						case 'S', 's':
							{
								position1091 := position
								depth++
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
										goto l1089
									}
									position++
								}
							l1092:
								{
									position1094, tokenIndex1094, depth1094 := position, tokenIndex, depth
									if buffer[position] != rune('t') {
										goto l1095
									}
									position++
									goto l1094
								l1095:
									position, tokenIndex, depth = position1094, tokenIndex1094, depth1094
									if buffer[position] != rune('T') {
										goto l1089
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
										goto l1089
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
										goto l1089
									}
									position++
								}
							l1098:
								{
									position1100, tokenIndex1100, depth1100 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1101
									}
									position++
									goto l1100
								l1101:
									position, tokenIndex, depth = position1100, tokenIndex1100, depth1100
									if buffer[position] != rune('U') {
										goto l1089
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
										goto l1089
									}
									position++
								}
							l1102:
								{
									position1104, tokenIndex1104, depth1104 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1105
									}
									position++
									goto l1104
								l1105:
									position, tokenIndex, depth = position1104, tokenIndex1104, depth1104
									if buffer[position] != rune('D') {
										goto l1089
									}
									position++
								}
							l1104:
								if !rules[ruleskip]() {
									goto l1089
								}
								depth--
								add(ruleSTRUUID, position1091)
							}
							break
						case 'U', 'u':
							{
								position1106 := position
								depth++
								{
									position1107, tokenIndex1107, depth1107 := position, tokenIndex, depth
									if buffer[position] != rune('u') {
										goto l1108
									}
									position++
									goto l1107
								l1108:
									position, tokenIndex, depth = position1107, tokenIndex1107, depth1107
									if buffer[position] != rune('U') {
										goto l1089
									}
									position++
								}
							l1107:
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
										goto l1089
									}
									position++
								}
							l1109:
								{
									position1111, tokenIndex1111, depth1111 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1112
									}
									position++
									goto l1111
								l1112:
									position, tokenIndex, depth = position1111, tokenIndex1111, depth1111
									if buffer[position] != rune('I') {
										goto l1089
									}
									position++
								}
							l1111:
								{
									position1113, tokenIndex1113, depth1113 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1114
									}
									position++
									goto l1113
								l1114:
									position, tokenIndex, depth = position1113, tokenIndex1113, depth1113
									if buffer[position] != rune('D') {
										goto l1089
									}
									position++
								}
							l1113:
								if !rules[ruleskip]() {
									goto l1089
								}
								depth--
								add(ruleUUID, position1106)
							}
							break
						case 'N', 'n':
							{
								position1115 := position
								depth++
								{
									position1116, tokenIndex1116, depth1116 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1117
									}
									position++
									goto l1116
								l1117:
									position, tokenIndex, depth = position1116, tokenIndex1116, depth1116
									if buffer[position] != rune('N') {
										goto l1089
									}
									position++
								}
							l1116:
								{
									position1118, tokenIndex1118, depth1118 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1119
									}
									position++
									goto l1118
								l1119:
									position, tokenIndex, depth = position1118, tokenIndex1118, depth1118
									if buffer[position] != rune('O') {
										goto l1089
									}
									position++
								}
							l1118:
								{
									position1120, tokenIndex1120, depth1120 := position, tokenIndex, depth
									if buffer[position] != rune('w') {
										goto l1121
									}
									position++
									goto l1120
								l1121:
									position, tokenIndex, depth = position1120, tokenIndex1120, depth1120
									if buffer[position] != rune('W') {
										goto l1089
									}
									position++
								}
							l1120:
								if !rules[ruleskip]() {
									goto l1089
								}
								depth--
								add(ruleNOW, position1115)
							}
							break
						default:
							{
								position1122 := position
								depth++
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
										goto l1089
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
										goto l1089
									}
									position++
								}
							l1125:
								{
									position1127, tokenIndex1127, depth1127 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1128
									}
									position++
									goto l1127
								l1128:
									position, tokenIndex, depth = position1127, tokenIndex1127, depth1127
									if buffer[position] != rune('N') {
										goto l1089
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
										goto l1089
									}
									position++
								}
							l1129:
								if !rules[ruleskip]() {
									goto l1089
								}
								depth--
								add(ruleRAND, position1122)
							}
							break
						}
					}

					if !rules[rulenil]() {
						goto l1089
					}
					goto l551
				l1089:
					position, tokenIndex, depth = position551, tokenIndex551, depth551
					{
						switch buffer[position] {
						case 'E', 'N', 'e', 'n':
							{
								position1132, tokenIndex1132, depth1132 := position, tokenIndex, depth
								{
									position1134 := position
									depth++
									{
										position1135, tokenIndex1135, depth1135 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1136
										}
										position++
										goto l1135
									l1136:
										position, tokenIndex, depth = position1135, tokenIndex1135, depth1135
										if buffer[position] != rune('E') {
											goto l1133
										}
										position++
									}
								l1135:
									{
										position1137, tokenIndex1137, depth1137 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1138
										}
										position++
										goto l1137
									l1138:
										position, tokenIndex, depth = position1137, tokenIndex1137, depth1137
										if buffer[position] != rune('X') {
											goto l1133
										}
										position++
									}
								l1137:
									{
										position1139, tokenIndex1139, depth1139 := position, tokenIndex, depth
										if buffer[position] != rune('i') {
											goto l1140
										}
										position++
										goto l1139
									l1140:
										position, tokenIndex, depth = position1139, tokenIndex1139, depth1139
										if buffer[position] != rune('I') {
											goto l1133
										}
										position++
									}
								l1139:
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
											goto l1133
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
											goto l1133
										}
										position++
									}
								l1143:
									{
										position1145, tokenIndex1145, depth1145 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1146
										}
										position++
										goto l1145
									l1146:
										position, tokenIndex, depth = position1145, tokenIndex1145, depth1145
										if buffer[position] != rune('S') {
											goto l1133
										}
										position++
									}
								l1145:
									if !rules[ruleskip]() {
										goto l1133
									}
									depth--
									add(ruleEXISTS, position1134)
								}
								goto l1132
							l1133:
								position, tokenIndex, depth = position1132, tokenIndex1132, depth1132
								{
									position1147 := position
									depth++
									{
										position1148, tokenIndex1148, depth1148 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1149
										}
										position++
										goto l1148
									l1149:
										position, tokenIndex, depth = position1148, tokenIndex1148, depth1148
										if buffer[position] != rune('N') {
											goto l549
										}
										position++
									}
								l1148:
									{
										position1150, tokenIndex1150, depth1150 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1151
										}
										position++
										goto l1150
									l1151:
										position, tokenIndex, depth = position1150, tokenIndex1150, depth1150
										if buffer[position] != rune('O') {
											goto l549
										}
										position++
									}
								l1150:
									{
										position1152, tokenIndex1152, depth1152 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1153
										}
										position++
										goto l1152
									l1153:
										position, tokenIndex, depth = position1152, tokenIndex1152, depth1152
										if buffer[position] != rune('T') {
											goto l549
										}
										position++
									}
								l1152:
									if buffer[position] != rune(' ') {
										goto l549
									}
									position++
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
											goto l549
										}
										position++
									}
								l1154:
									{
										position1156, tokenIndex1156, depth1156 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1157
										}
										position++
										goto l1156
									l1157:
										position, tokenIndex, depth = position1156, tokenIndex1156, depth1156
										if buffer[position] != rune('X') {
											goto l549
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
											goto l549
										}
										position++
									}
								l1158:
									{
										position1160, tokenIndex1160, depth1160 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1161
										}
										position++
										goto l1160
									l1161:
										position, tokenIndex, depth = position1160, tokenIndex1160, depth1160
										if buffer[position] != rune('S') {
											goto l549
										}
										position++
									}
								l1160:
									{
										position1162, tokenIndex1162, depth1162 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1163
										}
										position++
										goto l1162
									l1163:
										position, tokenIndex, depth = position1162, tokenIndex1162, depth1162
										if buffer[position] != rune('T') {
											goto l549
										}
										position++
									}
								l1162:
									{
										position1164, tokenIndex1164, depth1164 := position, tokenIndex, depth
										if buffer[position] != rune('s') {
											goto l1165
										}
										position++
										goto l1164
									l1165:
										position, tokenIndex, depth = position1164, tokenIndex1164, depth1164
										if buffer[position] != rune('S') {
											goto l549
										}
										position++
									}
								l1164:
									if !rules[ruleskip]() {
										goto l549
									}
									depth--
									add(ruleNOTEXIST, position1147)
								}
							}
						l1132:
							if !rules[rulegroupGraphPattern]() {
								goto l549
							}
							break
						case 'I', 'i':
							{
								position1166 := position
								depth++
								{
									position1167, tokenIndex1167, depth1167 := position, tokenIndex, depth
									if buffer[position] != rune('i') {
										goto l1168
									}
									position++
									goto l1167
								l1168:
									position, tokenIndex, depth = position1167, tokenIndex1167, depth1167
									if buffer[position] != rune('I') {
										goto l549
									}
									position++
								}
							l1167:
								{
									position1169, tokenIndex1169, depth1169 := position, tokenIndex, depth
									if buffer[position] != rune('f') {
										goto l1170
									}
									position++
									goto l1169
								l1170:
									position, tokenIndex, depth = position1169, tokenIndex1169, depth1169
									if buffer[position] != rune('F') {
										goto l549
									}
									position++
								}
							l1169:
								if !rules[ruleskip]() {
									goto l549
								}
								depth--
								add(ruleIF, position1166)
							}
							if !rules[ruleLPAREN]() {
								goto l549
							}
							if !rules[ruleexpression]() {
								goto l549
							}
							if !rules[ruleCOMMA]() {
								goto l549
							}
							if !rules[ruleexpression]() {
								goto l549
							}
							if !rules[ruleCOMMA]() {
								goto l549
							}
							if !rules[ruleexpression]() {
								goto l549
							}
							if !rules[ruleRPAREN]() {
								goto l549
							}
							break
						case 'C', 'c':
							{
								position1171, tokenIndex1171, depth1171 := position, tokenIndex, depth
								{
									position1173 := position
									depth++
									{
										position1174, tokenIndex1174, depth1174 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1175
										}
										position++
										goto l1174
									l1175:
										position, tokenIndex, depth = position1174, tokenIndex1174, depth1174
										if buffer[position] != rune('C') {
											goto l1172
										}
										position++
									}
								l1174:
									{
										position1176, tokenIndex1176, depth1176 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1177
										}
										position++
										goto l1176
									l1177:
										position, tokenIndex, depth = position1176, tokenIndex1176, depth1176
										if buffer[position] != rune('O') {
											goto l1172
										}
										position++
									}
								l1176:
									{
										position1178, tokenIndex1178, depth1178 := position, tokenIndex, depth
										if buffer[position] != rune('n') {
											goto l1179
										}
										position++
										goto l1178
									l1179:
										position, tokenIndex, depth = position1178, tokenIndex1178, depth1178
										if buffer[position] != rune('N') {
											goto l1172
										}
										position++
									}
								l1178:
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
											goto l1172
										}
										position++
									}
								l1180:
									{
										position1182, tokenIndex1182, depth1182 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1183
										}
										position++
										goto l1182
									l1183:
										position, tokenIndex, depth = position1182, tokenIndex1182, depth1182
										if buffer[position] != rune('A') {
											goto l1172
										}
										position++
									}
								l1182:
									{
										position1184, tokenIndex1184, depth1184 := position, tokenIndex, depth
										if buffer[position] != rune('t') {
											goto l1185
										}
										position++
										goto l1184
									l1185:
										position, tokenIndex, depth = position1184, tokenIndex1184, depth1184
										if buffer[position] != rune('T') {
											goto l1172
										}
										position++
									}
								l1184:
									if !rules[ruleskip]() {
										goto l1172
									}
									depth--
									add(ruleCONCAT, position1173)
								}
								goto l1171
							l1172:
								position, tokenIndex, depth = position1171, tokenIndex1171, depth1171
								{
									position1186 := position
									depth++
									{
										position1187, tokenIndex1187, depth1187 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1188
										}
										position++
										goto l1187
									l1188:
										position, tokenIndex, depth = position1187, tokenIndex1187, depth1187
										if buffer[position] != rune('C') {
											goto l549
										}
										position++
									}
								l1187:
									{
										position1189, tokenIndex1189, depth1189 := position, tokenIndex, depth
										if buffer[position] != rune('o') {
											goto l1190
										}
										position++
										goto l1189
									l1190:
										position, tokenIndex, depth = position1189, tokenIndex1189, depth1189
										if buffer[position] != rune('O') {
											goto l549
										}
										position++
									}
								l1189:
									{
										position1191, tokenIndex1191, depth1191 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1192
										}
										position++
										goto l1191
									l1192:
										position, tokenIndex, depth = position1191, tokenIndex1191, depth1191
										if buffer[position] != rune('A') {
											goto l549
										}
										position++
									}
								l1191:
									{
										position1193, tokenIndex1193, depth1193 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1194
										}
										position++
										goto l1193
									l1194:
										position, tokenIndex, depth = position1193, tokenIndex1193, depth1193
										if buffer[position] != rune('L') {
											goto l549
										}
										position++
									}
								l1193:
									{
										position1195, tokenIndex1195, depth1195 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1196
										}
										position++
										goto l1195
									l1196:
										position, tokenIndex, depth = position1195, tokenIndex1195, depth1195
										if buffer[position] != rune('E') {
											goto l549
										}
										position++
									}
								l1195:
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
											goto l549
										}
										position++
									}
								l1197:
									{
										position1199, tokenIndex1199, depth1199 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1200
										}
										position++
										goto l1199
									l1200:
										position, tokenIndex, depth = position1199, tokenIndex1199, depth1199
										if buffer[position] != rune('C') {
											goto l549
										}
										position++
									}
								l1199:
									{
										position1201, tokenIndex1201, depth1201 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1202
										}
										position++
										goto l1201
									l1202:
										position, tokenIndex, depth = position1201, tokenIndex1201, depth1201
										if buffer[position] != rune('E') {
											goto l549
										}
										position++
									}
								l1201:
									if !rules[ruleskip]() {
										goto l549
									}
									depth--
									add(ruleCOALESCE, position1186)
								}
							}
						l1171:
							if !rules[ruleargList]() {
								goto l549
							}
							break
						case 'B', 'b':
							{
								position1203 := position
								depth++
								{
									position1204, tokenIndex1204, depth1204 := position, tokenIndex, depth
									if buffer[position] != rune('b') {
										goto l1205
									}
									position++
									goto l1204
								l1205:
									position, tokenIndex, depth = position1204, tokenIndex1204, depth1204
									if buffer[position] != rune('B') {
										goto l549
									}
									position++
								}
							l1204:
								{
									position1206, tokenIndex1206, depth1206 := position, tokenIndex, depth
									if buffer[position] != rune('n') {
										goto l1207
									}
									position++
									goto l1206
								l1207:
									position, tokenIndex, depth = position1206, tokenIndex1206, depth1206
									if buffer[position] != rune('N') {
										goto l549
									}
									position++
								}
							l1206:
								{
									position1208, tokenIndex1208, depth1208 := position, tokenIndex, depth
									if buffer[position] != rune('o') {
										goto l1209
									}
									position++
									goto l1208
								l1209:
									position, tokenIndex, depth = position1208, tokenIndex1208, depth1208
									if buffer[position] != rune('O') {
										goto l549
									}
									position++
								}
							l1208:
								{
									position1210, tokenIndex1210, depth1210 := position, tokenIndex, depth
									if buffer[position] != rune('d') {
										goto l1211
									}
									position++
									goto l1210
								l1211:
									position, tokenIndex, depth = position1210, tokenIndex1210, depth1210
									if buffer[position] != rune('D') {
										goto l549
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
										goto l549
									}
									position++
								}
							l1212:
								if !rules[ruleskip]() {
									goto l549
								}
								depth--
								add(ruleBNODE, position1203)
							}
							{
								position1214, tokenIndex1214, depth1214 := position, tokenIndex, depth
								if !rules[ruleLPAREN]() {
									goto l1215
								}
								if !rules[ruleexpression]() {
									goto l1215
								}
								if !rules[ruleRPAREN]() {
									goto l1215
								}
								goto l1214
							l1215:
								position, tokenIndex, depth = position1214, tokenIndex1214, depth1214
								if !rules[rulenil]() {
									goto l549
								}
							}
						l1214:
							break
						default:
							{
								position1216, tokenIndex1216, depth1216 := position, tokenIndex, depth
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
										if buffer[position] != rune('u') {
											goto l1222
										}
										position++
										goto l1221
									l1222:
										position, tokenIndex, depth = position1221, tokenIndex1221, depth1221
										if buffer[position] != rune('U') {
											goto l1217
										}
										position++
									}
								l1221:
									{
										position1223, tokenIndex1223, depth1223 := position, tokenIndex, depth
										if buffer[position] != rune('b') {
											goto l1224
										}
										position++
										goto l1223
									l1224:
										position, tokenIndex, depth = position1223, tokenIndex1223, depth1223
										if buffer[position] != rune('B') {
											goto l1217
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
											goto l1217
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
											goto l1217
										}
										position++
									}
								l1227:
									{
										position1229, tokenIndex1229, depth1229 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1230
										}
										position++
										goto l1229
									l1230:
										position, tokenIndex, depth = position1229, tokenIndex1229, depth1229
										if buffer[position] != rune('R') {
											goto l1217
										}
										position++
									}
								l1229:
									if !rules[ruleskip]() {
										goto l1217
									}
									depth--
									add(ruleSUBSTR, position1218)
								}
								goto l1216
							l1217:
								position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
								{
									position1232 := position
									depth++
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
											goto l1231
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
											goto l1231
										}
										position++
									}
								l1235:
									{
										position1237, tokenIndex1237, depth1237 := position, tokenIndex, depth
										if buffer[position] != rune('p') {
											goto l1238
										}
										position++
										goto l1237
									l1238:
										position, tokenIndex, depth = position1237, tokenIndex1237, depth1237
										if buffer[position] != rune('P') {
											goto l1231
										}
										position++
									}
								l1237:
									{
										position1239, tokenIndex1239, depth1239 := position, tokenIndex, depth
										if buffer[position] != rune('l') {
											goto l1240
										}
										position++
										goto l1239
									l1240:
										position, tokenIndex, depth = position1239, tokenIndex1239, depth1239
										if buffer[position] != rune('L') {
											goto l1231
										}
										position++
									}
								l1239:
									{
										position1241, tokenIndex1241, depth1241 := position, tokenIndex, depth
										if buffer[position] != rune('a') {
											goto l1242
										}
										position++
										goto l1241
									l1242:
										position, tokenIndex, depth = position1241, tokenIndex1241, depth1241
										if buffer[position] != rune('A') {
											goto l1231
										}
										position++
									}
								l1241:
									{
										position1243, tokenIndex1243, depth1243 := position, tokenIndex, depth
										if buffer[position] != rune('c') {
											goto l1244
										}
										position++
										goto l1243
									l1244:
										position, tokenIndex, depth = position1243, tokenIndex1243, depth1243
										if buffer[position] != rune('C') {
											goto l1231
										}
										position++
									}
								l1243:
									{
										position1245, tokenIndex1245, depth1245 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1246
										}
										position++
										goto l1245
									l1246:
										position, tokenIndex, depth = position1245, tokenIndex1245, depth1245
										if buffer[position] != rune('E') {
											goto l1231
										}
										position++
									}
								l1245:
									if !rules[ruleskip]() {
										goto l1231
									}
									depth--
									add(ruleREPLACE, position1232)
								}
								goto l1216
							l1231:
								position, tokenIndex, depth = position1216, tokenIndex1216, depth1216
								{
									position1247 := position
									depth++
									{
										position1248, tokenIndex1248, depth1248 := position, tokenIndex, depth
										if buffer[position] != rune('r') {
											goto l1249
										}
										position++
										goto l1248
									l1249:
										position, tokenIndex, depth = position1248, tokenIndex1248, depth1248
										if buffer[position] != rune('R') {
											goto l549
										}
										position++
									}
								l1248:
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
											goto l549
										}
										position++
									}
								l1250:
									{
										position1252, tokenIndex1252, depth1252 := position, tokenIndex, depth
										if buffer[position] != rune('g') {
											goto l1253
										}
										position++
										goto l1252
									l1253:
										position, tokenIndex, depth = position1252, tokenIndex1252, depth1252
										if buffer[position] != rune('G') {
											goto l549
										}
										position++
									}
								l1252:
									{
										position1254, tokenIndex1254, depth1254 := position, tokenIndex, depth
										if buffer[position] != rune('e') {
											goto l1255
										}
										position++
										goto l1254
									l1255:
										position, tokenIndex, depth = position1254, tokenIndex1254, depth1254
										if buffer[position] != rune('E') {
											goto l549
										}
										position++
									}
								l1254:
									{
										position1256, tokenIndex1256, depth1256 := position, tokenIndex, depth
										if buffer[position] != rune('x') {
											goto l1257
										}
										position++
										goto l1256
									l1257:
										position, tokenIndex, depth = position1256, tokenIndex1256, depth1256
										if buffer[position] != rune('X') {
											goto l549
										}
										position++
									}
								l1256:
									if !rules[ruleskip]() {
										goto l549
									}
									depth--
									add(ruleREGEX, position1247)
								}
							}
						l1216:
							if !rules[ruleLPAREN]() {
								goto l549
							}
							if !rules[ruleexpression]() {
								goto l549
							}
							if !rules[ruleCOMMA]() {
								goto l549
							}
							if !rules[ruleexpression]() {
								goto l549
							}
							{
								position1258, tokenIndex1258, depth1258 := position, tokenIndex, depth
								if !rules[ruleCOMMA]() {
									goto l1258
								}
								if !rules[ruleexpression]() {
									goto l1258
								}
								goto l1259
							l1258:
								position, tokenIndex, depth = position1258, tokenIndex1258, depth1258
							}
						l1259:
							if !rules[ruleRPAREN]() {
								goto l549
							}
							break
						}
					}

				}
			l551:
				depth--
				add(rulebuiltinCall, position550)
			}
			return true
		l549:
			position, tokenIndex, depth = position549, tokenIndex549, depth549
			return false
		},
		/* 61 var <- <(('?' / '$') VARNAME skip)> */
		func() bool {
			position1260, tokenIndex1260, depth1260 := position, tokenIndex, depth
			{
				position1261 := position
				depth++
				{
					position1262, tokenIndex1262, depth1262 := position, tokenIndex, depth
					if buffer[position] != rune('?') {
						goto l1263
					}
					position++
					goto l1262
				l1263:
					position, tokenIndex, depth = position1262, tokenIndex1262, depth1262
					if buffer[position] != rune('$') {
						goto l1260
					}
					position++
				}
			l1262:
				{
					position1264 := position
					depth++
					{
						position1267, tokenIndex1267, depth1267 := position, tokenIndex, depth
						{
							position1269 := position
							depth++
							{
								position1270, tokenIndex1270, depth1270 := position, tokenIndex, depth
								{
									position1272 := position
									depth++
									{
										position1273, tokenIndex1273, depth1273 := position, tokenIndex, depth
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1274
										}
										position++
										goto l1273
									l1274:
										position, tokenIndex, depth = position1273, tokenIndex1273, depth1273
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1271
										}
										position++
									}
								l1273:
									depth--
									add(rulePN_CHARS_BASE, position1272)
								}
								goto l1270
							l1271:
								position, tokenIndex, depth = position1270, tokenIndex1270, depth1270
								if buffer[position] != rune('_') {
									goto l1268
								}
								position++
							}
						l1270:
							depth--
							add(rulePN_CHARS_U, position1269)
						}
						goto l1267
					l1268:
						position, tokenIndex, depth = position1267, tokenIndex1267, depth1267
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1260
						}
						position++
					}
				l1267:
				l1265:
					{
						position1266, tokenIndex1266, depth1266 := position, tokenIndex, depth
						{
							position1275, tokenIndex1275, depth1275 := position, tokenIndex, depth
							{
								position1277 := position
								depth++
								{
									position1278, tokenIndex1278, depth1278 := position, tokenIndex, depth
									{
										position1280 := position
										depth++
										{
											position1281, tokenIndex1281, depth1281 := position, tokenIndex, depth
											if c := buffer[position]; c < rune('a') || c > rune('z') {
												goto l1282
											}
											position++
											goto l1281
										l1282:
											position, tokenIndex, depth = position1281, tokenIndex1281, depth1281
											if c := buffer[position]; c < rune('A') || c > rune('Z') {
												goto l1279
											}
											position++
										}
									l1281:
										depth--
										add(rulePN_CHARS_BASE, position1280)
									}
									goto l1278
								l1279:
									position, tokenIndex, depth = position1278, tokenIndex1278, depth1278
									if buffer[position] != rune('_') {
										goto l1276
									}
									position++
								}
							l1278:
								depth--
								add(rulePN_CHARS_U, position1277)
							}
							goto l1275
						l1276:
							position, tokenIndex, depth = position1275, tokenIndex1275, depth1275
							if c := buffer[position]; c < rune('0') || c > rune('9') {
								goto l1266
							}
							position++
						}
					l1275:
						goto l1265
					l1266:
						position, tokenIndex, depth = position1266, tokenIndex1266, depth1266
					}
					depth--
					add(ruleVARNAME, position1264)
				}
				if !rules[ruleskip]() {
					goto l1260
				}
				depth--
				add(rulevar, position1261)
			}
			return true
		l1260:
			position, tokenIndex, depth = position1260, tokenIndex1260, depth1260
			return false
		},
		/* 62 iriref <- <(iri / prefixedName)> */
		func() bool {
			position1283, tokenIndex1283, depth1283 := position, tokenIndex, depth
			{
				position1284 := position
				depth++
				{
					position1285, tokenIndex1285, depth1285 := position, tokenIndex, depth
					if !rules[ruleiri]() {
						goto l1286
					}
					goto l1285
				l1286:
					position, tokenIndex, depth = position1285, tokenIndex1285, depth1285
					{
						position1287 := position
						depth++
					l1288:
						{
							position1289, tokenIndex1289, depth1289 := position, tokenIndex, depth
							{
								position1290, tokenIndex1290, depth1290 := position, tokenIndex, depth
								{
									position1291, tokenIndex1291, depth1291 := position, tokenIndex, depth
									if buffer[position] != rune(':') {
										goto l1292
									}
									position++
									goto l1291
								l1292:
									position, tokenIndex, depth = position1291, tokenIndex1291, depth1291
									if buffer[position] != rune(' ') {
										goto l1290
									}
									position++
								}
							l1291:
								goto l1289
							l1290:
								position, tokenIndex, depth = position1290, tokenIndex1290, depth1290
							}
							if !matchDot() {
								goto l1289
							}
							goto l1288
						l1289:
							position, tokenIndex, depth = position1289, tokenIndex1289, depth1289
						}
						if buffer[position] != rune(':') {
							goto l1283
						}
						position++
					l1293:
						{
							position1294, tokenIndex1294, depth1294 := position, tokenIndex, depth
							{
								position1295, tokenIndex1295, depth1295 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1296
								}
								position++
								goto l1295
							l1296:
								position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
								if c := buffer[position]; c < rune('0') || c > rune('9') {
									goto l1297
								}
								position++
								goto l1295
							l1297:
								position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
								if c := buffer[position]; c < rune('.') || c > rune('_') {
									goto l1298
								}
								position++
								goto l1295
							l1298:
								position, tokenIndex, depth = position1295, tokenIndex1295, depth1295
								{
									switch buffer[position] {
									case '%':
										if buffer[position] != rune('%') {
											goto l1294
										}
										position++
										break
									case ':':
										if buffer[position] != rune(':') {
											goto l1294
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1294
										}
										position++
										break
									}
								}

							}
						l1295:
							goto l1293
						l1294:
							position, tokenIndex, depth = position1294, tokenIndex1294, depth1294
						}
						if !rules[ruleskip]() {
							goto l1283
						}
						depth--
						add(ruleprefixedName, position1287)
					}
				}
			l1285:
				depth--
				add(ruleiriref, position1284)
			}
			return true
		l1283:
			position, tokenIndex, depth = position1283, tokenIndex1283, depth1283
			return false
		},
		/* 63 iri <- <('<' (!'>' .)* '>' skip)> */
		func() bool {
			position1300, tokenIndex1300, depth1300 := position, tokenIndex, depth
			{
				position1301 := position
				depth++
				if buffer[position] != rune('<') {
					goto l1300
				}
				position++
			l1302:
				{
					position1303, tokenIndex1303, depth1303 := position, tokenIndex, depth
					{
						position1304, tokenIndex1304, depth1304 := position, tokenIndex, depth
						if buffer[position] != rune('>') {
							goto l1304
						}
						position++
						goto l1303
					l1304:
						position, tokenIndex, depth = position1304, tokenIndex1304, depth1304
					}
					if !matchDot() {
						goto l1303
					}
					goto l1302
				l1303:
					position, tokenIndex, depth = position1303, tokenIndex1303, depth1303
				}
				if buffer[position] != rune('>') {
					goto l1300
				}
				position++
				if !rules[ruleskip]() {
					goto l1300
				}
				depth--
				add(ruleiri, position1301)
			}
			return true
		l1300:
			position, tokenIndex, depth = position1300, tokenIndex1300, depth1300
			return false
		},
		/* 64 prefixedName <- <((!(':' / ' ') .)* ':' ([A-Z] / [0-9] / [.-_] / ((&('%') '%') | (&(':') ':') | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])))* skip)> */
		nil,
		/* 65 literal <- <(string (('@' ([a-z] / [A-Z])+ ('-' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z]))+)*) / ('^' '^' iriref))? skip)> */
		func() bool {
			position1306, tokenIndex1306, depth1306 := position, tokenIndex, depth
			{
				position1307 := position
				depth++
				{
					position1308 := position
					depth++
					if buffer[position] != rune('"') {
						goto l1306
					}
					position++
				l1309:
					{
						position1310, tokenIndex1310, depth1310 := position, tokenIndex, depth
						{
							position1311, tokenIndex1311, depth1311 := position, tokenIndex, depth
							if buffer[position] != rune('"') {
								goto l1311
							}
							position++
							goto l1310
						l1311:
							position, tokenIndex, depth = position1311, tokenIndex1311, depth1311
						}
						if !matchDot() {
							goto l1310
						}
						goto l1309
					l1310:
						position, tokenIndex, depth = position1310, tokenIndex1310, depth1310
					}
					if buffer[position] != rune('"') {
						goto l1306
					}
					position++
					depth--
					add(rulestring, position1308)
				}
				{
					position1312, tokenIndex1312, depth1312 := position, tokenIndex, depth
					{
						position1314, tokenIndex1314, depth1314 := position, tokenIndex, depth
						if buffer[position] != rune('@') {
							goto l1315
						}
						position++
						{
							position1318, tokenIndex1318, depth1318 := position, tokenIndex, depth
							if c := buffer[position]; c < rune('a') || c > rune('z') {
								goto l1319
							}
							position++
							goto l1318
						l1319:
							position, tokenIndex, depth = position1318, tokenIndex1318, depth1318
							if c := buffer[position]; c < rune('A') || c > rune('Z') {
								goto l1315
							}
							position++
						}
					l1318:
					l1316:
						{
							position1317, tokenIndex1317, depth1317 := position, tokenIndex, depth
							{
								position1320, tokenIndex1320, depth1320 := position, tokenIndex, depth
								if c := buffer[position]; c < rune('a') || c > rune('z') {
									goto l1321
								}
								position++
								goto l1320
							l1321:
								position, tokenIndex, depth = position1320, tokenIndex1320, depth1320
								if c := buffer[position]; c < rune('A') || c > rune('Z') {
									goto l1317
								}
								position++
							}
						l1320:
							goto l1316
						l1317:
							position, tokenIndex, depth = position1317, tokenIndex1317, depth1317
						}
					l1322:
						{
							position1323, tokenIndex1323, depth1323 := position, tokenIndex, depth
							if buffer[position] != rune('-') {
								goto l1323
							}
							position++
							{
								switch buffer[position] {
								case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
									if c := buffer[position]; c < rune('0') || c > rune('9') {
										goto l1323
									}
									position++
									break
								case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
									if c := buffer[position]; c < rune('A') || c > rune('Z') {
										goto l1323
									}
									position++
									break
								default:
									if c := buffer[position]; c < rune('a') || c > rune('z') {
										goto l1323
									}
									position++
									break
								}
							}

						l1324:
							{
								position1325, tokenIndex1325, depth1325 := position, tokenIndex, depth
								{
									switch buffer[position] {
									case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
										if c := buffer[position]; c < rune('0') || c > rune('9') {
											goto l1325
										}
										position++
										break
									case 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z':
										if c := buffer[position]; c < rune('A') || c > rune('Z') {
											goto l1325
										}
										position++
										break
									default:
										if c := buffer[position]; c < rune('a') || c > rune('z') {
											goto l1325
										}
										position++
										break
									}
								}

								goto l1324
							l1325:
								position, tokenIndex, depth = position1325, tokenIndex1325, depth1325
							}
							goto l1322
						l1323:
							position, tokenIndex, depth = position1323, tokenIndex1323, depth1323
						}
						goto l1314
					l1315:
						position, tokenIndex, depth = position1314, tokenIndex1314, depth1314
						if buffer[position] != rune('^') {
							goto l1312
						}
						position++
						if buffer[position] != rune('^') {
							goto l1312
						}
						position++
						if !rules[ruleiriref]() {
							goto l1312
						}
					}
				l1314:
					goto l1313
				l1312:
					position, tokenIndex, depth = position1312, tokenIndex1312, depth1312
				}
			l1313:
				if !rules[ruleskip]() {
					goto l1306
				}
				depth--
				add(ruleliteral, position1307)
			}
			return true
		l1306:
			position, tokenIndex, depth = position1306, tokenIndex1306, depth1306
			return false
		},
		/* 66 string <- <('"' (!'"' .)* '"')> */
		nil,
		/* 67 numericLiteral <- <(('+' / '-')? [0-9]+ ('.' [0-9]*)? skip)> */
		func() bool {
			position1329, tokenIndex1329, depth1329 := position, tokenIndex, depth
			{
				position1330 := position
				depth++
				{
					position1331, tokenIndex1331, depth1331 := position, tokenIndex, depth
					{
						position1333, tokenIndex1333, depth1333 := position, tokenIndex, depth
						if buffer[position] != rune('+') {
							goto l1334
						}
						position++
						goto l1333
					l1334:
						position, tokenIndex, depth = position1333, tokenIndex1333, depth1333
						if buffer[position] != rune('-') {
							goto l1331
						}
						position++
					}
				l1333:
					goto l1332
				l1331:
					position, tokenIndex, depth = position1331, tokenIndex1331, depth1331
				}
			l1332:
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1329
				}
				position++
			l1335:
				{
					position1336, tokenIndex1336, depth1336 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1336
					}
					position++
					goto l1335
				l1336:
					position, tokenIndex, depth = position1336, tokenIndex1336, depth1336
				}
				{
					position1337, tokenIndex1337, depth1337 := position, tokenIndex, depth
					if buffer[position] != rune('.') {
						goto l1337
					}
					position++
				l1339:
					{
						position1340, tokenIndex1340, depth1340 := position, tokenIndex, depth
						if c := buffer[position]; c < rune('0') || c > rune('9') {
							goto l1340
						}
						position++
						goto l1339
					l1340:
						position, tokenIndex, depth = position1340, tokenIndex1340, depth1340
					}
					goto l1338
				l1337:
					position, tokenIndex, depth = position1337, tokenIndex1337, depth1337
				}
			l1338:
				if !rules[ruleskip]() {
					goto l1329
				}
				depth--
				add(rulenumericLiteral, position1330)
			}
			return true
		l1329:
			position, tokenIndex, depth = position1329, tokenIndex1329, depth1329
			return false
		},
		/* 68 signedNumericLiteral <- <(('+' / '-') [0-9]+ ('.' [0-9]*)? skip)> */
		nil,
		/* 69 booleanLiteral <- <(TRUE / FALSE)> */
		func() bool {
			position1342, tokenIndex1342, depth1342 := position, tokenIndex, depth
			{
				position1343 := position
				depth++
				{
					position1344, tokenIndex1344, depth1344 := position, tokenIndex, depth
					{
						position1346 := position
						depth++
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
								goto l1345
							}
							position++
						}
					l1347:
						{
							position1349, tokenIndex1349, depth1349 := position, tokenIndex, depth
							if buffer[position] != rune('r') {
								goto l1350
							}
							position++
							goto l1349
						l1350:
							position, tokenIndex, depth = position1349, tokenIndex1349, depth1349
							if buffer[position] != rune('R') {
								goto l1345
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
								goto l1345
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
								goto l1345
							}
							position++
						}
					l1353:
						if !rules[ruleskip]() {
							goto l1345
						}
						depth--
						add(ruleTRUE, position1346)
					}
					goto l1344
				l1345:
					position, tokenIndex, depth = position1344, tokenIndex1344, depth1344
					{
						position1355 := position
						depth++
						{
							position1356, tokenIndex1356, depth1356 := position, tokenIndex, depth
							if buffer[position] != rune('f') {
								goto l1357
							}
							position++
							goto l1356
						l1357:
							position, tokenIndex, depth = position1356, tokenIndex1356, depth1356
							if buffer[position] != rune('F') {
								goto l1342
							}
							position++
						}
					l1356:
						{
							position1358, tokenIndex1358, depth1358 := position, tokenIndex, depth
							if buffer[position] != rune('a') {
								goto l1359
							}
							position++
							goto l1358
						l1359:
							position, tokenIndex, depth = position1358, tokenIndex1358, depth1358
							if buffer[position] != rune('A') {
								goto l1342
							}
							position++
						}
					l1358:
						{
							position1360, tokenIndex1360, depth1360 := position, tokenIndex, depth
							if buffer[position] != rune('l') {
								goto l1361
							}
							position++
							goto l1360
						l1361:
							position, tokenIndex, depth = position1360, tokenIndex1360, depth1360
							if buffer[position] != rune('L') {
								goto l1342
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
								goto l1342
							}
							position++
						}
					l1362:
						{
							position1364, tokenIndex1364, depth1364 := position, tokenIndex, depth
							if buffer[position] != rune('e') {
								goto l1365
							}
							position++
							goto l1364
						l1365:
							position, tokenIndex, depth = position1364, tokenIndex1364, depth1364
							if buffer[position] != rune('E') {
								goto l1342
							}
							position++
						}
					l1364:
						if !rules[ruleskip]() {
							goto l1342
						}
						depth--
						add(ruleFALSE, position1355)
					}
				}
			l1344:
				depth--
				add(rulebooleanLiteral, position1343)
			}
			return true
		l1342:
			position, tokenIndex, depth = position1342, tokenIndex1342, depth1342
			return false
		},
		/* 70 blankNode <- <(blankNodeLabel / anon)> */
		nil,
		/* 71 blankNodeLabel <- <('_' ':' ((&('0' | '1' | '2' | '3' | '4' | '5' | '6' | '7' | '8' | '9') [0-9]) | (&('A' | 'B' | 'C' | 'D' | 'E' | 'F' | 'G' | 'H' | 'I' | 'J' | 'K' | 'L' | 'M' | 'N' | 'O' | 'P' | 'Q' | 'R' | 'S' | 'T' | 'U' | 'V' | 'W' | 'X' | 'Y' | 'Z') [A-Z]) | (&('a' | 'b' | 'c' | 'd' | 'e' | 'f' | 'g' | 'h' | 'i' | 'j' | 'k' | 'l' | 'm' | 'n' | 'o' | 'p' | 'q' | 'r' | 's' | 't' | 'u' | 'v' | 'w' | 'x' | 'y' | 'z') [a-z])) ([a-z] / [A-Z] / [0-9] / [.-_])? skip)> */
		nil,
		/* 72 anon <- <('[' ws* ']' skip)> */
		nil,
		/* 73 nil <- <('(' ws* ')' skip)> */
		func() bool {
			position1369, tokenIndex1369, depth1369 := position, tokenIndex, depth
			{
				position1370 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1369
				}
				position++
			l1371:
				{
					position1372, tokenIndex1372, depth1372 := position, tokenIndex, depth
					if !rules[rulews]() {
						goto l1372
					}
					goto l1371
				l1372:
					position, tokenIndex, depth = position1372, tokenIndex1372, depth1372
				}
				if buffer[position] != rune(')') {
					goto l1369
				}
				position++
				if !rules[ruleskip]() {
					goto l1369
				}
				depth--
				add(rulenil, position1370)
			}
			return true
		l1369:
			position, tokenIndex, depth = position1369, tokenIndex1369, depth1369
			return false
		},
		/* 74 VARNAME <- <(PN_CHARS_U / [0-9])+> */
		nil,
		/* 75 PN_CHARS_U <- <(PN_CHARS_BASE / '_')> */
		nil,
		/* 76 PN_CHARS_BASE <- <([a-z] / [A-Z])> */
		nil,
		/* 77 PREFIX <- <(('p' / 'P') ('r' / 'R') ('e' / 'E') ('f' / 'F') ('i' / 'I') ('x' / 'X') skip)> */
		nil,
		/* 78 TRUE <- <(('t' / 'T') ('r' / 'R') ('u' / 'U') ('e' / 'E') skip)> */
		nil,
		/* 79 FALSE <- <(('f' / 'F') ('a' / 'A') ('l' / 'L') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 80 BASE <- <(('b' / 'B') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 81 SELECT <- <(('s' / 'S') ('e' / 'E') ('l' / 'L') ('e' / 'E') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 82 REDUCED <- <(('r' / 'R') ('e' / 'E') ('d' / 'D') ('u' / 'U') ('c' / 'C') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 83 DISTINCT <- <(('d' / 'D') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('i' / 'I') ('n' / 'N') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 84 FROM <- <(('f' / 'F') ('r' / 'R') ('o' / 'O') ('m' / 'M') skip)> */
		nil,
		/* 85 NAMED <- <(('n' / 'N') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('d' / 'D') skip)> */
		nil,
		/* 86 WHERE <- <(('w' / 'W') ('h' / 'H') ('e' / 'E') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 87 LBRACE <- <('{' skip)> */
		func() bool {
			position1386, tokenIndex1386, depth1386 := position, tokenIndex, depth
			{
				position1387 := position
				depth++
				if buffer[position] != rune('{') {
					goto l1386
				}
				position++
				if !rules[ruleskip]() {
					goto l1386
				}
				depth--
				add(ruleLBRACE, position1387)
			}
			return true
		l1386:
			position, tokenIndex, depth = position1386, tokenIndex1386, depth1386
			return false
		},
		/* 88 RBRACE <- <('}' skip)> */
		func() bool {
			position1388, tokenIndex1388, depth1388 := position, tokenIndex, depth
			{
				position1389 := position
				depth++
				if buffer[position] != rune('}') {
					goto l1388
				}
				position++
				if !rules[ruleskip]() {
					goto l1388
				}
				depth--
				add(ruleRBRACE, position1389)
			}
			return true
		l1388:
			position, tokenIndex, depth = position1388, tokenIndex1388, depth1388
			return false
		},
		/* 89 LBRACK <- <('[' skip)> */
		nil,
		/* 90 RBRACK <- <(']' skip)> */
		nil,
		/* 91 SEMICOLON <- <(';' skip)> */
		nil,
		/* 92 COMMA <- <(',' skip)> */
		func() bool {
			position1393, tokenIndex1393, depth1393 := position, tokenIndex, depth
			{
				position1394 := position
				depth++
				if buffer[position] != rune(',') {
					goto l1393
				}
				position++
				if !rules[ruleskip]() {
					goto l1393
				}
				depth--
				add(ruleCOMMA, position1394)
			}
			return true
		l1393:
			position, tokenIndex, depth = position1393, tokenIndex1393, depth1393
			return false
		},
		/* 93 DOT <- <('.' skip)> */
		func() bool {
			position1395, tokenIndex1395, depth1395 := position, tokenIndex, depth
			{
				position1396 := position
				depth++
				if buffer[position] != rune('.') {
					goto l1395
				}
				position++
				if !rules[ruleskip]() {
					goto l1395
				}
				depth--
				add(ruleDOT, position1396)
			}
			return true
		l1395:
			position, tokenIndex, depth = position1395, tokenIndex1395, depth1395
			return false
		},
		/* 94 COLON <- <(':' skip)> */
		nil,
		/* 95 PIPE <- <('|' skip)> */
		func() bool {
			position1398, tokenIndex1398, depth1398 := position, tokenIndex, depth
			{
				position1399 := position
				depth++
				if buffer[position] != rune('|') {
					goto l1398
				}
				position++
				if !rules[ruleskip]() {
					goto l1398
				}
				depth--
				add(rulePIPE, position1399)
			}
			return true
		l1398:
			position, tokenIndex, depth = position1398, tokenIndex1398, depth1398
			return false
		},
		/* 96 SLASH <- <('/' skip)> */
		func() bool {
			position1400, tokenIndex1400, depth1400 := position, tokenIndex, depth
			{
				position1401 := position
				depth++
				if buffer[position] != rune('/') {
					goto l1400
				}
				position++
				if !rules[ruleskip]() {
					goto l1400
				}
				depth--
				add(ruleSLASH, position1401)
			}
			return true
		l1400:
			position, tokenIndex, depth = position1400, tokenIndex1400, depth1400
			return false
		},
		/* 97 INVERSE <- <('^' skip)> */
		func() bool {
			position1402, tokenIndex1402, depth1402 := position, tokenIndex, depth
			{
				position1403 := position
				depth++
				if buffer[position] != rune('^') {
					goto l1402
				}
				position++
				if !rules[ruleskip]() {
					goto l1402
				}
				depth--
				add(ruleINVERSE, position1403)
			}
			return true
		l1402:
			position, tokenIndex, depth = position1402, tokenIndex1402, depth1402
			return false
		},
		/* 98 LPAREN <- <('(' skip)> */
		func() bool {
			position1404, tokenIndex1404, depth1404 := position, tokenIndex, depth
			{
				position1405 := position
				depth++
				if buffer[position] != rune('(') {
					goto l1404
				}
				position++
				if !rules[ruleskip]() {
					goto l1404
				}
				depth--
				add(ruleLPAREN, position1405)
			}
			return true
		l1404:
			position, tokenIndex, depth = position1404, tokenIndex1404, depth1404
			return false
		},
		/* 99 RPAREN <- <(')' skip)> */
		func() bool {
			position1406, tokenIndex1406, depth1406 := position, tokenIndex, depth
			{
				position1407 := position
				depth++
				if buffer[position] != rune(')') {
					goto l1406
				}
				position++
				if !rules[ruleskip]() {
					goto l1406
				}
				depth--
				add(ruleRPAREN, position1407)
			}
			return true
		l1406:
			position, tokenIndex, depth = position1406, tokenIndex1406, depth1406
			return false
		},
		/* 100 ISA <- <('a' skip)> */
		func() bool {
			position1408, tokenIndex1408, depth1408 := position, tokenIndex, depth
			{
				position1409 := position
				depth++
				if buffer[position] != rune('a') {
					goto l1408
				}
				position++
				if !rules[ruleskip]() {
					goto l1408
				}
				depth--
				add(ruleISA, position1409)
			}
			return true
		l1408:
			position, tokenIndex, depth = position1408, tokenIndex1408, depth1408
			return false
		},
		/* 101 NOT <- <('!' skip)> */
		func() bool {
			position1410, tokenIndex1410, depth1410 := position, tokenIndex, depth
			{
				position1411 := position
				depth++
				if buffer[position] != rune('!') {
					goto l1410
				}
				position++
				if !rules[ruleskip]() {
					goto l1410
				}
				depth--
				add(ruleNOT, position1411)
			}
			return true
		l1410:
			position, tokenIndex, depth = position1410, tokenIndex1410, depth1410
			return false
		},
		/* 102 STAR <- <('*' skip)> */
		func() bool {
			position1412, tokenIndex1412, depth1412 := position, tokenIndex, depth
			{
				position1413 := position
				depth++
				if buffer[position] != rune('*') {
					goto l1412
				}
				position++
				if !rules[ruleskip]() {
					goto l1412
				}
				depth--
				add(ruleSTAR, position1413)
			}
			return true
		l1412:
			position, tokenIndex, depth = position1412, tokenIndex1412, depth1412
			return false
		},
		/* 103 PLUS <- <('+' skip)> */
		func() bool {
			position1414, tokenIndex1414, depth1414 := position, tokenIndex, depth
			{
				position1415 := position
				depth++
				if buffer[position] != rune('+') {
					goto l1414
				}
				position++
				if !rules[ruleskip]() {
					goto l1414
				}
				depth--
				add(rulePLUS, position1415)
			}
			return true
		l1414:
			position, tokenIndex, depth = position1414, tokenIndex1414, depth1414
			return false
		},
		/* 104 MINUS <- <('-' skip)> */
		func() bool {
			position1416, tokenIndex1416, depth1416 := position, tokenIndex, depth
			{
				position1417 := position
				depth++
				if buffer[position] != rune('-') {
					goto l1416
				}
				position++
				if !rules[ruleskip]() {
					goto l1416
				}
				depth--
				add(ruleMINUS, position1417)
			}
			return true
		l1416:
			position, tokenIndex, depth = position1416, tokenIndex1416, depth1416
			return false
		},
		/* 105 OPTIONAL <- <(('o' / 'O') ('p' / 'P') ('t' / 'T') ('i' / 'I') ('o' / 'O') ('n' / 'N') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 106 UNION <- <(('u' / 'U') ('n' / 'N') ('i' / 'I') ('o' / 'O') ('n' / 'N') skip)> */
		nil,
		/* 107 LIMIT <- <(('l' / 'L') ('i' / 'I') ('m' / 'M') ('i' / 'I') ('t' / 'T') skip)> */
		nil,
		/* 108 OFFSET <- <(('o' / 'O') ('f' / 'F') ('f' / 'F') ('s' / 'S') ('e' / 'E') ('t' / 'T') skip)> */
		nil,
		/* 109 INTEGER <- <([0-9]+ skip)> */
		func() bool {
			position1422, tokenIndex1422, depth1422 := position, tokenIndex, depth
			{
				position1423 := position
				depth++
				if c := buffer[position]; c < rune('0') || c > rune('9') {
					goto l1422
				}
				position++
			l1424:
				{
					position1425, tokenIndex1425, depth1425 := position, tokenIndex, depth
					if c := buffer[position]; c < rune('0') || c > rune('9') {
						goto l1425
					}
					position++
					goto l1424
				l1425:
					position, tokenIndex, depth = position1425, tokenIndex1425, depth1425
				}
				if !rules[ruleskip]() {
					goto l1422
				}
				depth--
				add(ruleINTEGER, position1423)
			}
			return true
		l1422:
			position, tokenIndex, depth = position1422, tokenIndex1422, depth1422
			return false
		},
		/* 110 CONSTRUCT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('c' / 'C') ('t' / 'T') skip)> */
		nil,
		/* 111 DESCRIBE <- <(('d' / 'D') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('r' / 'R') ('i' / 'I') ('b' / 'B') ('e' / 'E') skip)> */
		nil,
		/* 112 ASK <- <(('a' / 'A') ('s' / 'S') ('k' / 'K') skip)> */
		nil,
		/* 113 OR <- <('|' '|' skip)> */
		nil,
		/* 114 AND <- <('&' '&' skip)> */
		nil,
		/* 115 EQ <- <('=' skip)> */
		nil,
		/* 116 NE <- <('!' '=' skip)> */
		nil,
		/* 117 GT <- <('>' skip)> */
		nil,
		/* 118 LT <- <('<' skip)> */
		nil,
		/* 119 LE <- <('<' '=' skip)> */
		nil,
		/* 120 GE <- <('>' '=' skip)> */
		nil,
		/* 121 IN <- <(('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 122 NOTIN <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('i' / 'I') ('n' / 'N') skip)> */
		nil,
		/* 123 AS <- <(('a' / 'A') ('s' / 'S') skip)> */
		func() bool {
			position1439, tokenIndex1439, depth1439 := position, tokenIndex, depth
			{
				position1440 := position
				depth++
				{
					position1441, tokenIndex1441, depth1441 := position, tokenIndex, depth
					if buffer[position] != rune('a') {
						goto l1442
					}
					position++
					goto l1441
				l1442:
					position, tokenIndex, depth = position1441, tokenIndex1441, depth1441
					if buffer[position] != rune('A') {
						goto l1439
					}
					position++
				}
			l1441:
				{
					position1443, tokenIndex1443, depth1443 := position, tokenIndex, depth
					if buffer[position] != rune('s') {
						goto l1444
					}
					position++
					goto l1443
				l1444:
					position, tokenIndex, depth = position1443, tokenIndex1443, depth1443
					if buffer[position] != rune('S') {
						goto l1439
					}
					position++
				}
			l1443:
				if !rules[ruleskip]() {
					goto l1439
				}
				depth--
				add(ruleAS, position1440)
			}
			return true
		l1439:
			position, tokenIndex, depth = position1439, tokenIndex1439, depth1439
			return false
		},
		/* 124 STR <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 125 LANG <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 126 DATATYPE <- <(('d' / 'D') ('a' / 'A') ('t' / 'T') ('a' / 'A') ('t' / 'T') ('y' / 'Y') ('p' / 'P') ('e' / 'E') skip)> */
		nil,
		/* 127 IRI <- <(('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 128 URI <- <(('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 129 ABS <- <(('a' / 'A') ('b' / 'B') ('s' / 'S') skip)> */
		nil,
		/* 130 CEIL <- <(('c' / 'C') ('e' / 'E') ('i' / 'I') ('l' / 'L') skip)> */
		nil,
		/* 131 ROUND <- <(('r' / 'R') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 132 FLOOR <- <(('f' / 'F') ('l' / 'L') ('o' / 'O') ('o' / 'O') ('r' / 'R') skip)> */
		nil,
		/* 133 STRLEN <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('e' / 'E') ('n' / 'N') skip)> */
		nil,
		/* 134 UCASE <- <(('u' / 'U') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 135 LCASE <- <(('l' / 'L') ('c' / 'C') ('a' / 'A') ('s' / 'S') ('e' / 'E') skip)> */
		nil,
		/* 136 ENCODEFORURI <- <(('e' / 'E') ('n' / 'N') ('c' / 'C') ('o' / 'O') ('d' / 'D') ('e' / 'E') '_' ('f' / 'F') ('o' / 'O') ('r' / 'R') '_' ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 137 YEAR <- <(('y' / 'Y') ('e' / 'E') ('a' / 'A') ('r' / 'R') skip)> */
		nil,
		/* 138 MONTH <- <(('m' / 'M') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('h' / 'H') skip)> */
		nil,
		/* 139 DAY <- <(('d' / 'D') ('a' / 'A') ('y' / 'Y') skip)> */
		nil,
		/* 140 HOURS <- <(('h' / 'H') ('o' / 'O') ('u' / 'U') ('r' / 'R') ('s' / 'S') skip)> */
		nil,
		/* 141 MINUTES <- <(('m' / 'M') ('i' / 'I') ('n' / 'N') ('u' / 'U') ('t' / 'T') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 142 SECONDS <- <(('s' / 'S') ('e' / 'E') ('c' / 'C') ('o' / 'O') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 143 TIMEZONE <- <(('t' / 'T') ('i' / 'I') ('m' / 'M') ('e' / 'E') ('z' / 'Z') ('o' / 'O') ('n' / 'N') ('e' / 'E') skip)> */
		nil,
		/* 144 TZ <- <(('t' / 'T') ('z' / 'Z') skip)> */
		nil,
		/* 145 MD5 <- <(('m' / 'M') ('d' / 'D') skip)> */
		nil,
		/* 146 SHA1 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '1' skip)> */
		nil,
		/* 147 SHA256 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '2' '5' '6' skip)> */
		nil,
		/* 148 SHA384 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '3' '8' '4' skip)> */
		nil,
		/* 149 SHA512 <- <(('s' / 'S') ('h' / 'H') ('a' / 'A') '5' '1' '2' skip)> */
		nil,
		/* 150 ISIRI <- <(('i' / 'I') ('s' / 'S') ('i' / 'I') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 151 ISURI <- <(('i' / 'I') ('s' / 'S') ('u' / 'U') ('r' / 'R') ('i' / 'I') skip)> */
		nil,
		/* 152 ISBLANK <- <(('i' / 'I') ('s' / 'S') ('b' / 'B') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('k' / 'K') skip)> */
		nil,
		/* 153 ISLITERAL <- <(('i' / 'I') ('s' / 'S') ('l' / 'L') ('i' / 'I') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('a' / 'A') ('l' / 'L') skip)> */
		nil,
		/* 154 ISNUMERIC <- <(('i' / 'I') ('s' / 'S') ('n' / 'N') ('u' / 'U') ('m' / 'M') ('e' / 'E') ('r' / 'R') ('i' / 'I') ('c' / 'C') skip)> */
		nil,
		/* 155 LANGMATCHES <- <(('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') ('m' / 'M') ('a' / 'A') ('t' / 'T') ('c' / 'C') ('h' / 'H') ('e' / 'E') ('s' / 'S') skip)> */
		nil,
		/* 156 CONTAINS <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('t' / 'T') ('a' / 'A') ('i' / 'I') ('n' / 'N') ('s' / 'S') skip)> */
		nil,
		/* 157 STRSTARTS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('s' / 'S') ('t' / 'T') ('a' / 'A') ('r' / 'R') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 158 STRENDS <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('e' / 'E') ('n' / 'N') ('d' / 'D') ('s' / 'S') skip)> */
		nil,
		/* 159 STRBEFORE <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('b' / 'B') ('e' / 'E') ('f' / 'F') ('o' / 'O') ('r' / 'R') ('e' / 'E') skip)> */
		nil,
		/* 160 STRAFTER <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('a' / 'A') ('f' / 'F') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 161 STRLANG <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('l' / 'L') ('a' / 'A') ('n' / 'N') ('g' / 'G') skip)> */
		nil,
		/* 162 STRDT <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('d' / 'D') ('t' / 'T') skip)> */
		nil,
		/* 163 SAMETERM <- <(('s' / 'S') ('a' / 'A') ('m' / 'M') ('e' / 'E') ('t' / 'T') ('e' / 'E') ('r' / 'R') ('m' / 'M') skip)> */
		nil,
		/* 164 BOUND <- <(('b' / 'B') ('o' / 'O') ('u' / 'U') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 165 BNODE <- <(('b' / 'B') ('n' / 'N') ('o' / 'O') ('d' / 'D') ('e' / 'E') skip)> */
		nil,
		/* 166 RAND <- <(('r' / 'R') ('a' / 'A') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 167 NOW <- <(('n' / 'N') ('o' / 'O') ('w' / 'W') skip)> */
		nil,
		/* 168 UUID <- <(('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 169 STRUUID <- <(('s' / 'S') ('t' / 'T') ('r' / 'R') ('u' / 'U') ('u' / 'U') ('i' / 'I') ('d' / 'D') skip)> */
		nil,
		/* 170 CONCAT <- <(('c' / 'C') ('o' / 'O') ('n' / 'N') ('c' / 'C') ('a' / 'A') ('t' / 'T') skip)> */
		nil,
		/* 171 SUBSTR <- <(('s' / 'S') ('u' / 'U') ('b' / 'B') ('s' / 'S') ('t' / 'T') ('r' / 'R') skip)> */
		nil,
		/* 172 REPLACE <- <(('r' / 'R') ('e' / 'E') ('p' / 'P') ('l' / 'L') ('a' / 'A') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 173 REGEX <- <(('r' / 'R') ('e' / 'E') ('g' / 'G') ('e' / 'E') ('x' / 'X') skip)> */
		nil,
		/* 174 IF <- <(('i' / 'I') ('f' / 'F') skip)> */
		nil,
		/* 175 EXISTS <- <(('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 176 NOTEXIST <- <(('n' / 'N') ('o' / 'O') ('t' / 'T') ' ' ('e' / 'E') ('x' / 'X') ('i' / 'I') ('s' / 'S') ('t' / 'T') ('s' / 'S') skip)> */
		nil,
		/* 177 COALESCE <- <(('c' / 'C') ('o' / 'O') ('a' / 'A') ('l' / 'L') ('e' / 'E') ('s' / 'S') ('c' / 'C') ('e' / 'E') skip)> */
		nil,
		/* 178 FILTER <- <(('f' / 'F') ('i' / 'I') ('l' / 'L') ('t' / 'T') ('e' / 'E') ('r' / 'R') skip)> */
		nil,
		/* 179 BIND <- <(('b' / 'B') ('i' / 'I') ('n' / 'N') ('d' / 'D') skip)> */
		nil,
		/* 180 skip <- <(ws / comment)*> */
		func() bool {
			{
				position1502 := position
				depth++
			l1503:
				{
					position1504, tokenIndex1504, depth1504 := position, tokenIndex, depth
					{
						position1505, tokenIndex1505, depth1505 := position, tokenIndex, depth
						if !rules[rulews]() {
							goto l1506
						}
						goto l1505
					l1506:
						position, tokenIndex, depth = position1505, tokenIndex1505, depth1505
						{
							position1507 := position
							depth++
							if buffer[position] != rune('#') {
								goto l1504
							}
							position++
						l1508:
							{
								position1509, tokenIndex1509, depth1509 := position, tokenIndex, depth
								{
									position1510, tokenIndex1510, depth1510 := position, tokenIndex, depth
									if !rules[ruleendOfLine]() {
										goto l1510
									}
									goto l1509
								l1510:
									position, tokenIndex, depth = position1510, tokenIndex1510, depth1510
								}
								if !matchDot() {
									goto l1509
								}
								goto l1508
							l1509:
								position, tokenIndex, depth = position1509, tokenIndex1509, depth1509
							}
							if !rules[ruleendOfLine]() {
								goto l1504
							}
							depth--
							add(rulecomment, position1507)
						}
					}
				l1505:
					goto l1503
				l1504:
					position, tokenIndex, depth = position1504, tokenIndex1504, depth1504
				}
				depth--
				add(ruleskip, position1502)
			}
			return true
		},
		/* 181 ws <- <((&('\v') '\v') | (&('\f') '\f') | (&('\t') '\t') | (&(' ') ' ') | (&('\n' | '\r') endOfLine))> */
		func() bool {
			position1511, tokenIndex1511, depth1511 := position, tokenIndex, depth
			{
				position1512 := position
				depth++
				{
					switch buffer[position] {
					case '\v':
						if buffer[position] != rune('\v') {
							goto l1511
						}
						position++
						break
					case '\f':
						if buffer[position] != rune('\f') {
							goto l1511
						}
						position++
						break
					case '\t':
						if buffer[position] != rune('\t') {
							goto l1511
						}
						position++
						break
					case ' ':
						if buffer[position] != rune(' ') {
							goto l1511
						}
						position++
						break
					default:
						if !rules[ruleendOfLine]() {
							goto l1511
						}
						break
					}
				}

				depth--
				add(rulews, position1512)
			}
			return true
		l1511:
			position, tokenIndex, depth = position1511, tokenIndex1511, depth1511
			return false
		},
		/* 182 comment <- <('#' (!endOfLine .)* endOfLine)> */
		nil,
		/* 183 endOfLine <- <(('\r' '\n') / '\n' / '\r')> */
		func() bool {
			position1515, tokenIndex1515, depth1515 := position, tokenIndex, depth
			{
				position1516 := position
				depth++
				{
					position1517, tokenIndex1517, depth1517 := position, tokenIndex, depth
					if buffer[position] != rune('\r') {
						goto l1518
					}
					position++
					if buffer[position] != rune('\n') {
						goto l1518
					}
					position++
					goto l1517
				l1518:
					position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
					if buffer[position] != rune('\n') {
						goto l1519
					}
					position++
					goto l1517
				l1519:
					position, tokenIndex, depth = position1517, tokenIndex1517, depth1517
					if buffer[position] != rune('\r') {
						goto l1515
					}
					position++
				}
			l1517:
				depth--
				add(ruleendOfLine, position1516)
			}
			return true
		l1515:
			position, tokenIndex, depth = position1515, tokenIndex1515, depth1515
			return false
		},
	}
	p.rules = rules
}
