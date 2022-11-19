package main

// import "fmt"

// // getFuntionName exploits the simplicity of C having the very first
// // ID token be the function name.
// func getFunctionName(tokens []*Token) string {
// 	for _, token := range tokens {
// 		if token.typ == "ID" {
// 			return token.name
// 		}
// 	}
// 	return ""
// }

// func lazyVriableGrepper() []*Token {
// 	Ps.rollback()
// 	var variableBody []*Token
// 	for {
// 		token := Ps.next()
// 		if token == nil {
// 			break
// 		}

// 		variableBody = append(variableBody, token)

// 		if token.typ == "SEMICOLON" {
// 			break
// 		}
// 	}
// 	Ps.snapshot()
// 	return variableBody
// }

// // lazyFunctionGrepper is a naive C function extractor. It works by
// // using OPENCURLY and CLOSECURLY tokens to determine when a function has been
// // completely scanned from top to bottom.
// func lazyFunctionGrepper() []*Token {
// 	var (
// 		functionBody []*Token
// 		depth        int
// 	)

// 	functionName := getFunctionName(Ps.tokens)

// 	// Go to the head of the function before parsing took place.
// 	Ps.rollback()
// 	// funtionFound := false
// 	for {
// 		token := Ps.next()
// 		if token == nil {
// 			break
// 		}
// 		// if token.typ == "ID" && !funtionFound {
// 		// 	token.label = "DECL"
// 		// 	funtionFound = true
// 		// }

// 		functionBody = append(functionBody, token)
// 		if token.typ == "OPENCURLY" {
// 			depth = 1
// 			break
// 		}
// 	}

// 	for {
// 		token := Ps.next()
// 		if token == nil {
// 			break
// 		}

// 		if token.typ == "OPENCURLY" {
// 			depth++
// 		} else if token.typ == "CLOSECURLY" {
// 			depth--
// 		}

// 		functionBody = append(functionBody, token)
// 		if depth == 0 {
// 			break
// 		}
// 	}

// 	functionName = getFunctionName(functionBody)
// 	fmt.Println(functionName)
// 	Ps.snapshot()
// 	pushFtable(functionName, functionBody[0].filepath, functionBody)
// 	return functionBody
// }

// func lazyStructGrepper() []*Token {
// 	var (
// 		structBody []*Token
// 		depth      int
// 	)

// 	structName := getFunctionName(Ps.tokens)
// 	Ps.rollback()
// 	for {
// 		token := Ps.next()
// 		if token == nil {
// 			break
// 		}

// 		structBody = append(structBody, token)
// 		if token.typ == "OPENCURLY" {
// 			depth = 1
// 			break
// 		}
// 	}

// 	for {

// 		token := Ps.next()
// 		if token == nil {
// 			break
// 		}
// 		fmt.Printf("%s ", token.name)
// 		if token.typ == "OPENCURLY" {
// 			depth++
// 		} else if token.typ == "CLOSECURLY" {
// 			depth--
// 		}

// 		structBody = append(structBody, token)
// 		if depth == 0 {
// 			break
// 		}
// 	}

// 	if TypedefStructSet {
// 		TypedefStructSet = false
// 		addTypedef(Ps.next().name, structName)
// 		Ps.next()
// 	}
// 	pushSTable(structName, structBody[0].filepath, structBody)
// 	return structBody
// }
