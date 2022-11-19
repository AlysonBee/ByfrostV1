package main

import (
	"byfrostV1/golang"
)

var targetLanguage = "Go"
var supportedLanguages = map[string]func(content string, filepath string) *golang.Token{
	"Go": golang.GOscan,
}

func BuildTokenList(content string, filepath string) []*Token {
	var token *Token
	var tokenList []*Token
	golang.FilePtr = 0
	golang.LineNumber = 1

	scan := supportedLanguages[targetLanguage]
	for golang.FilePtr < len(content) {
		token = GolangTokenToToken(scan(content, filepath))
		tokenList = append(tokenList, token)
	}
	golang.GetImports(GenericTokensToGolangTokens(tokenList))
	return tokenList
}
