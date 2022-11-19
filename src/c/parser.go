package main

// import (
// 	"fmt"
// )

// var Ps *Parser

// func clojureOneOrMany(brackets []string, index int) {
// 	var subFlow []string

// 	subFlow = append(subFlow, brackets[index-1])
// 	evalSubFlow(subFlow)
// 	for {
// 		if !evalSubFlow(subFlow) {
// 			break
// 		}
// 	}
// }

// func clojureZeroOrOne(brackets []string, index int) bool {
// 	presentValue := brackets[index-1]

// 	if index-2 > 0 {
// 		// There's a list of more than 1 of the same value repeating
// 		if presentValue == brackets[index-2] {
// 			return false
// 		}
// 	}
// 	return true
// }

// func clojureRule(brackets []string, index int) {
// 	var subFlow []string

// 	subFlow = append(subFlow, brackets[index-1])
// 	for {
// 		if !evalSubFlow(subFlow) {
// 			break
// 		}
// 	}
// }

// func evalSubFlow(brackets []string) bool {
// 	index := 0

// 	for index < len(brackets) {
// 		fmt.Printf(">> %s|%s ---> %s\n", Ps.current().typ, Ps.current().name, brackets[index])
// 		if brackets[index] == "(" {
// 			index = evalBracket(index, brackets)
// 			index++
// 			continue
// 		}

// 		if brackets[index] == "?" {
// 			repeating := clojureZeroOrOne(brackets, index)
// 			if !repeating {
// 				fmt.Println("Error : Repeating values")
// 				return false
// 			}
// 			index++
// 			continue
// 		}

// 		if brackets[index] == "*" {
// 			clojureRule(brackets, index)
// 			index++
// 			continue
// 		}

// 		if Ps.current().typ != brackets[index] {
// 			if index+1 < len(brackets) && (brackets[index+1] == "*" || brackets[index+1] == "+") {
// 				fmt.Println("brackets[index] is ", brackets[index+1])
// 				index++
// 				continue
// 			} else {
// 				Ps.whereAreWe()
// 				fmt.Println("FALSE :RETURN from ", brackets)
// 				return false
// 			}
// 		}

// 		index++
// 		Ps.next()
// 	}

// 	fmt.Println("RETURN from ", brackets)
// 	return true
// }

// func evalBracket(index int, rule []string) int {
// 	var brackets []string
// 	subFlowHeight := 1

// 	index++
// 	for index < len(rule) {

// 		if rule[index] == ")" {
// 			subFlowHeight--
// 		} else if rule[index] == "(" {
// 			subFlowHeight++
// 		}

// 		if subFlowHeight == 0 {
// 			index++
// 			break
// 		}

// 		brackets = append(brackets, rule[index])
// 		index++
// 	}

// 	fmt.Println("brackets is ", brackets)
// 	evalSubFlow(brackets)
// 	if index < len(rule) && (rule[index] == "*" || rule[index] == "+") {
// 		for {
// 			if !evalSubFlow(brackets) {
// 				break
// 			}
// 		}
// 	}
// 	return index
// }

// func eval(rule []string) bool {
// 	Ps.snapshot()

// 	var toEval []*Token
// 	toEval = append(toEval, Ps.current())
// 	index := 0

// 	for index < len(rule) {
// 		fmt.Printf("%s|%s ---> %s|%d\n", Ps.current().typ, Ps.current().name, rule[index], index)
// 		if rule[index] == "(" {
// 			index = evalBracket(index, rule)
// 			index++
// 			continue
// 		}

// 		if rule[index] == "?" {
// 			repeating := clojureZeroOrOne(rule, index)
// 			if !repeating {
// 				fmt.Println("Error : Repeating values")
// 				return false
// 			}
// 			index++
// 			continue
// 		}

// 		if rule[index] == "*" {
// 			clojureRule(rule, index)
// 			index++
// 			continue
// 		}

// 		if Ps.current().typ != rule[index] {
// 			if index+1 < len(rule) && (rule[index+1] == "*" || rule[index+1] == "+") {
// 				fmt.Println("brackets[index] is ", rule[index+1])
// 				index++
// 				continue
// 			}
// 			fmt.Printf("Error, expected: %s but got %s\n", rule[index], Ps.current().typ)
// 			return false
// 		}
// 		index++
// 		Ps.next()
// 	}
// 	fmt.Println("PASS")
// 	return true
// }

// func parser(tokens []*Token, datatypes []string) {
// 	Ps = newParser(tokens)
// 	datatypeToEval := 0

// 	var pass bool

// 	for {
// 		if Ps.eof() || datatypeToEval == len(datatypes) {
// 			break
// 		}

// 		CSpecificCheck()

// 		fmt.Println("datatypeToEval", datatypeToEval)
// 		r, ok := regexRules[datatypes[datatypeToEval]]
// 		if ok {
// 			rc := 0
// 			for rc < len(r) {
// 				pass = eval(r[rc])
// 				if pass {
// 					break
// 				}

// 				// Rule eval fails. Reset to the start of evaluation and check next rule.
// 				rc++
// 				Ps.rollback()
// 			}

// 			if pass {
// 				fmt.Println(">>> ", datatypes[datatypeToEval])
// 				toViz := lazyStorage[datatypes[datatypeToEval]]
// 				visualize(toViz())
// 				datatypeToEval = 0
// 			} else if rc == len(r) {
// 				datatypeToEval++
// 			}
// 		}
// 	}
// }
