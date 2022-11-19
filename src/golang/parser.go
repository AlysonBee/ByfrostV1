package golang

import (
	"fmt"
	"strings"
)

var importName string

// GoParser is the structure that'll control the parsing process for Go tokens
type GoParser struct {
	tokens []*Token
	length int
	ptr    int
	save   int
}

type localVariable struct {
	variableName    string
	datatype        string
	locationPackage string
}

var localVars = make(map[string]localVariable)

func repeatingChar(c string, times int) {
	counter := 0

	for counter < times {
		fmt.Printf(c)
		counter++
	}
}

func newline(linenumber int) {
	fmt.Printf("\n%d ", linenumber)
}

func visualize(tokens []*Token) {
	currLine := 1

	for _, tk := range tokens {
		if tk.Line > currLine {
			newline(tk.Line)
			currLine = tk.Line
		}
		if tk.Spaces > 0 {
			repeatingChar(" ", tk.Spaces)
		}
		if tk.Tabs > 0 {
			repeatingChar("\t", tk.Tabs)
		}
		// if tk.label == "DECL" {
		// 	fmt.Println("DECL")
		// }
		// if tk.label == "STRUCT" {
		// 	fmt.Println("STRUCT")
		// }
		fmt.Printf(tk.Name)
	}
}

var parser *GoParser

func (p *GoParser) next() *Token {
	if p.ptr < p.length {
		t := p.tokens[p.ptr]
		if p.tokens[p.ptr].Name == "HelloIndexer" {
			p.tokens[p.ptr].Typ = "FUNCTION"
			p.tokens[p.ptr].Label = "FUNCTION"
		}
		t.SymbolTab = SnapshotTable(SymTab)
		p.ptr++
		return t
	}
	return nil
}

func (p *GoParser) current() *Token {
	return p.tokens[p.ptr]
}

// DEBUG
func (p *GoParser) whereAreWe() {
	fmt.Printf("%d of %d\n", p.ptr, p.length)
}

func (p *GoParser) prev() *Token {
	if p.ptr > 0 {
		return p.tokens[p.ptr-1]
	}
	return nil
}

func (p *GoParser) snapshot() {
	p.save = p.ptr
}

func (p *GoParser) rollback() {
	p.ptr = p.save
}

func (p *GoParser) eof() bool {
	return p.ptr >= p.length-1
}

func (p *GoParser) lookAhead() *Token {
	if p.ptr+1 >= p.length {
		return nil
	}

	return p.tokens[p.ptr+1]
}

func (p *GoParser) moveBy(steps int) {
	for steps > 0 {
		p.next()
		steps--
	}
}

func newParser(tokens []*Token) *GoParser {
	return &GoParser{
		tokens: tokens,
		length: len(tokens),
		ptr:    0,
		save:   0,
	}
}

func joinTokens(t1 []*Token, t2 []*Token) []*Token {
	for _, tk := range t2 {
		t1 = append(t1, tk)
	}
	return t1
}

func getVarDecl() []*Token {
	var (
		varDecl []*Token
	)
	parser.snapshot()

	lookahead := parser.current()
	if lookahead.Typ == "OPENBRACKET" {
		parser.next()
		newLineChecker := lookahead.Line
		lineIndex := 0
		for {
			if parser.eof() {
				break
			}

			token := parser.next()
			if token.Typ == "CLOSEBRACKET" {
				break
			}
			if token.Line > newLineChecker {
				newLineChecker = token.Line
				lineIndex = 0
			}

			if lineIndex > 0 {
				token.Typ = "PARAM_TYPE"
				token.Label = "PARAM_TYPE"
			}
			varDecl = append(varDecl, token)
		}
	} else {
		newLineChecker := lookahead.Line
		for {
			if parser.eof() {
				break
			}

			token := parser.next()
			varDecl = append(varDecl, token)
			if token.Line > newLineChecker {
				break
			}
		}
	}
	fmt.Println("VARIABLE\n")
	// visualize(varDecl)
	fmt.Println()
	return varDecl
}

// Make a unique deref token
// TODO: Move this into the Lexer
func getDeref(derefId *Token, assist string) *Token {
	var (
		derefString strings.Builder
		derefTokens []*Token
	)
	parser.snapshot()

	// ID
	derefString.WriteString(derefId.Name)

	// DOT
	token := parser.next()
	derefString.WriteString(token.Name)

	// ID
	token = parser.next()
	derefString.WriteString(token.Name)
	typ := "FUNCTION"

	derefRet := &Token{
		Name:   derefString.String(),
		Typ:    typ,
		Line:   derefId.Line,
		Spaces: derefId.Spaces,
		Tabs:   derefId.Tabs,

		Label: typ,
	}
	derefTokens = append(derefTokens, derefRet)
	if parser.current().Typ == "OPENBRACKET" {

	}

	return derefRet
}

func getFCall() []*Token {

	var (
		fCall []*Token
	)

	parser.snapshot()
	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		nextToken := parser.current()

		if nextToken != nil {
			// Overwrite ID token with new DEREF token
			if nextToken.Typ == "DOT" && token.Typ == "ID" {
				token = getDeref(token, resolveImportCall(token.Name))
				fCall = append(fCall, token)
			} else if nextToken.Typ == "OPENBRACKET" && token.Typ == "ID" {
				token.Label = "FUNCTION"
				token.Typ = "FUNCTION"
				fCall = append(fCall, token)
				fCall = joinTokens(fCall, getFCall())
			} else if token.Typ == "ID" {
				fCall = append(fCall, token)
			} else {
				fCall = append(fCall, token)
			}

		}

		if token.Typ == "CLOSEBRACKET" {
			break
		}
	}
	return fCall
}

func getEquAssign() []*Token {
	var (
		equAssig []*Token
	)

	lookahead := parser.current()
	newLineChecker := lookahead.Line
	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		if token.Line > newLineChecker || token.Typ == "OPENCURLY" {
			break
		}
		equAssig = append(equAssig, token)
	}

	return equAssig
}

func getDeclAssign(token *Token) []*Token {
	var (
		declAssign []*Token
	)
	lookahead := parser.current()
	newLineChecker := lookahead.Line
	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		if token.Line > newLineChecker || token.Typ == "OPENCURLY" {
			break
		}
		declAssign = append(declAssign, token)
	}

	return declAssign
}

func parseParameters() []*Token {
	var (
		params   []*Token
		imported string
	)

	valueType := []string{
		"PARAM_VALUE",
		"PARAM_TYPE",
	}
	paramIndex := 0
	brackets := 1
	token := parser.next()
	params = append(params, token)
	if token.Name == ")" {
		return params
	}

	for {
		if parser.eof() {
			break
		}

		token = parser.next()
		nextToken := parser.current()
		if nextToken.Typ == "ID" && token.Typ == "DOT" {
			fmt.Println("nextToken is ", nextToken.Name)
			parser.current().Attribute = imported
		}

		if token.Typ == "ID" && nextToken.Typ == "DOT" {
			//	fmt.Println("Resolve Parameter ", token.Name, parser.lookAhead().Name)
			imported = resolveImportCall(token.Name)
		} else {
			imported = ""
		}

		if token.Typ == "CLOSEBRACKET" {
			brackets--
			if brackets == 0 {
				params = append(params, token)
				break
			}
		} else if token.Typ == "ID" || token.Typ == "DOT" {
			token.labelAndType(valueType[paramIndex])
		}

		params = append(params, token)
		if token.Typ == "COMMA" {
			paramIndex = 0
		} else {
			paramIndex = 1
		}

	}
	fmt.Println("\n")
	visualize(params)
	fmt.Println("\n")
	return params
}

func getFunction(filepath string) {
	var (
		functionBody []*Token
		depth        int
	)

	if parser.lookAhead().Typ != "ID" {
		parser.next()
		return
	}

	token := parser.next()
	functionBody = append(functionBody, token)

	token = parser.next()
	functionBody = append(functionBody, token)

	fName := token.Name
	functionBody = append(functionBody, parseParameters()...)

	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		functionBody = append(functionBody, token)
		if token.Typ == "OPENCURLY" {
			depth++
			break
		}
	}

	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		nextToken := parser.current()

		if nextToken != nil {

			// Overwrite ID token with new DEREF token
			if nextToken.Typ == "DOT" && token.Typ == "ID" {
				token = getDeref(token, resolveImportCall(token.Name))
				functionBody = append(functionBody, token)
			} else if nextToken.Typ == "OPENBRACKET" && token.Typ == "ID" {
				token.labelAndType("FUNCITON")
				token.Label = "FUNCTION"
				token.Typ = "FUNCTION"
				functionBody = append(functionBody, token)
				functionBody = append(functionBody, getFCall()...)
			} else {
				functionBody = append(functionBody, token)
			}
		}

		if token.Typ == "OPENCURLY" {
			depth++
		} else if token.Typ == "CLOSECURLY" {
			depth--
		}

		if depth == 0 {
			break
		}
	}

	PushFunctionItem(fName, functionBody, filepath)
}

func createPath(path string) string {
	x := 0
	for i, c := range path {
		if c == '/' {
			x = i
			break
		}
	}

	return importName + path[x:]
}

func resolveImportCall(libraryName string) string {
	for _, imp := range Imports {
		if imp == nil {
			continue
		}
		if libraryName == imp.Alias || libraryName == imp.BaseName {
			absPath := createPath(imp.Path[1 : len(imp.Path)-1])
			// fmt.Println("absPath found ", absPath)
			return absPath
		}
	}
	//fmt.Println("Not found : basename is ", libraryName)
	return ""
}

func parseType(filepath string) []*Token {
	var (
		typeBody []*Token
	)

	token := parser.next()
	typeBody = append(typeBody, token)

	token = parser.next()
	typeName := token.Name

	typeBody = append(typeBody, token)

	token = parser.next()

	typeBody = append(typeBody, token)

	if parser.lookAhead().Typ == "OPENCURLY" {
		visualize(typeBody)
		return typeBody
	}

	if token.Name == "struct" || token.Name == "interface" {
		for {
			if parser.eof() {
				break
			}

			token = parser.next()
			nextToken := parser.current()
			if nextToken != nil {
				// Overwrite ID token with new DEREF token
				if nextToken.Typ == "DOT" && token.Typ == "ID" {
					resolveImportCall(token.Name)
					token = getDeref(token, resolveImportCall(token.Name))
					typeBody = append(typeBody, token)
				} else {
					typeBody = append(typeBody, token)
				}
			}
			//			typeBody = append(typeBody, token)
			if token.Typ == "CLOSECURLY" {
				break
			}
		}
	} else if token.Name == "func" {
		token = parser.next()
		typeBody = append(typeBody, token)
		openbracket := 1
		// get function body
		for {
			if parser.eof() {
				break
			}

			token = parser.next()
			nextToken := parser.current()
			if nextToken != nil {
				// Overwrite ID token with new DEREF token
				if nextToken.Typ == "DOT" && token.Typ == "ID" {
					token = getDeref(token, resolveImportCall(token.Name))
					typeBody = append(typeBody, token)
				} else {
					typeBody = append(typeBody, token)
				}
			}

			if token.Typ == "CLOSEBRACKET" {
				openbracket--
				if openbracket == 0 {
					break
				}
			}
		}

		// has return type?
		token = parser.next()
		if token.Typ == "OPENBRACKET" {
			typeBody = append(typeBody, token)
			openbracket = 1
			for {
				if parser.eof() {
					break
				}

				token = parser.next()
				typeBody = append(typeBody, token)
				if token.Typ == "CLOSEBRACKET" {
					openbracket--
					if openbracket == 0 {
						break
					}
				}

			}
		} else {
			currline := token.Line
			typeBody = append(typeBody, token)
			for {
				if parser.eof() {
					break
				}

				token = parser.next()
				if currline < token.Line {
					break
				}

				typeBody = append(typeBody, token)
			}
		}
	}
	PushType(typeName, filepath, typeBody)
	return typeBody
}

func ParseTypes(tokens []*Token, filepath string) {
	parser = newParser(tokens)

	for {
		if parser.eof() {
			break
		}

		currTok := parser.current()
		if currTok.Typ == "TYPE" {
			parseType(filepath)
			continue
		}

		parser.next()
	}
}

func getMethodStructName() ([]*Token, string) {
	var (
		methodStructName []*Token
		objectName       string
	)

	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		nextToken := parser.current()
		if nextToken.Typ == "CLOSEBRACKET" {
			objectName = token.Name
		}

		methodStructName = append(methodStructName, token)
		if token.Typ == "CLOSEBRACKET" {
			break
		}
	}

	return methodStructName, objectName
}

func getMethod(filepath string) {
	var (
		functionBody []*Token
		depth        int
	)

	token := parser.next()
	functionBody = append(functionBody, token)

	objectCode, objectName := getMethodStructName()
	functionBody = append(functionBody, objectCode...)

	token = parser.next()
	fName := token.Name
	functionBody = append(functionBody, token)

	token = parser.next()
	functionBody = append(functionBody, token)

	// fName := token.Name
	functionBody = append(functionBody, parseParameters()...)
	// visualize(functionBody)
	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		functionBody = append(functionBody, token)
		if token.Typ == "OPENCURLY" {
			depth++
			break
		}
	}

	for {
		if parser.eof() {
			break
		}

		token := parser.next()
		nextToken := parser.current()

		if nextToken != nil {

			// Overwrite ID token with new DEREF token
			if nextToken.Typ == "DOT" && token.Typ == "ID" {
				token = getDeref(token, resolveImportCall(token.Name))
				functionBody = append(functionBody, token)
			} else if nextToken.Typ == "OPENBRACKET" && token.Typ == "ID" {
				token.labelAndType("FUNCITON")
				token.Label = "FUNCTION"
				token.Typ = "FUNCTION"
				functionBody = append(functionBody, token)
				functionBody = append(functionBody, getFCall()...)
			} else {
				functionBody = append(functionBody, token)
			}
		}

		if token.Typ == "OPENCURLY" {
			depth++
		} else if token.Typ == "CLOSECURLY" {
			depth--
		}

		if depth == 0 {
			break
		}
	}
	// visualize(functionBody)
	// fmt.Println("\n\n")
	PushMethod(objectName, fName, functionBody)
	// fmt.Println("\n\n")
}

func Parser(tokens []*Token, filepath, absoluteProjectPath string) {
	parser = newParser(tokens)
	importName = absoluteProjectPath
	for {
		if parser.eof() {
			break
		}

		currTok := parser.current()
		if currTok.Typ == "FUNC" {
			if parser.lookAhead().Name == "(" {
				getMethod(filepath)
			} else {
				getFunction(filepath)
			}
			continue
		}

		parser.next()
	}
	// PrintAllMethods()
}
