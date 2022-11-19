package styling

// Setting new stylesheets needs a better config experience.

import (
	"encoding/json"
	"io/ioutil"
	"log"
)

type Stylesheet struct {
	StylingConfig    map[string]string
	EditorBackground string
}

var StylesheetConfig *Stylesheet

func newStylesheet(stylingConfig map[string]string, editorBackground string) *Stylesheet {
	return &Stylesheet{
		StylingConfig:    stylingConfig,
		EditorBackground: editorBackground,
	}
}

func readConfig(conrfigFile string) string {
	content, err := ioutil.ReadFile(conrfigFile)
	if err != nil {
		log.Fatal(err)
	}
	return string(content)
}

func PrepSyntaxHighlighting(languageJSON string) {
	var result map[string]string

	config := readConfig(languageJSON)
	if config == "" {
		return
	}

	json.Unmarshal([]byte(config), &result)
	for key, tag := range result {
		result[key] = tag
	}

	StylesheetConfig = newStylesheet(
		result,
		"black",
	)
}
