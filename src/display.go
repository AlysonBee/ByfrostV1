package main

import (
	"fmt"
	"byfrostV1/errors"
	"byfrostV1/mapping"
	"byfrostV1/server"
	"byfrostV1/styling"
	"strings"
)

type Display struct {
	name string
	path string

	namespace map[string]ImportNamespace

	active      bool
	rawTokens   []*Token
	coordinates mapping.Coordinate

	prev          *Display
	childDisplays map[string]*Display
}

var RootPath *Display

func idName(idPath string, parentString string) string {
	var fullID strings.Builder
	fullID.WriteString(idPath)
	fullID.WriteString(".")
	fullID.WriteString(parentString)
	return fullID.String()
}

func verifyNode(toPath string) *Display {
	resetPath := RootPath
	segments := strings.Split(toPath, "/")

	for _, segment := range segments {
		if len(segment) == 0 {
			continue
		}

		if doesDispalyExist(segment, resetPath) {
			resetPath = resetPath.childDisplays[segment]
		} else {
			fmt.Printf("Error: %s: path not found\n", toPath)
			return nil
		}
	}
	return resetPath
}

func targetNamespace(derefFunction string, derefToken string) string {
	namespaceTargeter := strings.Split(derefFunction, derefToken)
	if len(namespaceTargeter) > 1 {
		return namespaceTargeter[0]
	}
	return "."
}

func buildID(tokens []*Token, idPath string) []*Token {

	var idBuilder strings.Builder

	for i := range tokens {
		if tokens[i].label == "FUNCTION" || tokens[i].label == "DEREF" || tokens[i].label == "PARAM_TYPE" {
			idBuilder.WriteString(idPath)
			idBuilder.WriteString("-")
			idBuilder.WriteString(tokens[i].name)
			//	idBuilder.WriteString("-")
			//	idBuilder.WriteString(GlobalNamespace) // Make a global token library
			tokens[i].id = idBuilder.String()
			idBuilder.Reset()
		}

		tokens[i].styleTag = styling.AssignStyleTag(tokens[i].typ)
		i++
	}
	return tokens
}

func functinHeightAndWidth(tokens []*Token) (int, int) {
	var (
		height int
		width  int
		tmpRow int
	)

	for _, tk := range tokens {
		if tk.line > height {
			height = tk.line
		}

		if tk.spaces > 0 {
			spaces := 0
			for spaces < tk.spaces {
				tmpRow += 1
				spaces += 1
			}
		}

		if tk.tabs > 0 {
			tabs := 0
			for tabs < tk.tabs {
				tmpRow += 4
				tabs += 1
			}
		}

		tmpRow += len(tk.name)
		if tmpRow > width {
			width = tmpRow
		}
	}
	return height, width
}

func NewDisplay(name string, path string, html []*Token) *Display {
	height, width := functinHeightAndWidth(html)

	return &Display{
		name:      name,
		path:      path,
		namespace: Namespace,
		rawTokens: buildID(html, path),
		active:    true,
		coordinates: mapping.Coordinate{
			Height: height,
			Width:  width,
		},
		childDisplays: make(map[string]*Display),
	}
}

func doesDispalyExist(pathSegment string, currentDisplay *Display) bool {
	if currentDisplay == nil {
		return false
	}

	if _, found := currentDisplay.childDisplays[pathSegment]; found {
		return true
	}
	return false
}

func initDisplay() *Display {
	RootPath = NewDisplay("__init", "__init", []*Token{})
	RootPath.prev = &Display{}

	writable := RootPath
	return writable
}

func allToState(d *Display, state bool) {
	resetPath := d
	allPaths := []*Display{resetPath}
	tmpPaths := []*Display{}
	collapseList := []string{}

	i := 0
	for i < len(allPaths) {

		path := allPaths[i]
		path.active = state

		collapseList = append(collapseList, path.path) //+"-"+GlobalNamespace)
		for key := range path.childDisplays {
			tmpPaths = append(tmpPaths, path.childDisplays[key])
		}

		i++
		if i == len(allPaths) && len(tmpPaths) > 0 {
			allPaths = tmpPaths
			tmpPaths = []*Display{}
			i = 0
		}
	}
	Collapse = true
	// fmt.Println("collapseList is ", collapseList)
	CollapseList = collapseList
}

func isActive(toPath string, label string) (*Display, bool) {
	resetPath := verifyNode(toPath)
	if resetPath == nil {
		return nil, false
	}

	p, ok := resetPath.childDisplays[label]
	if ok {
		if p.active {
			p.active = false
			resetPath.childDisplays[label].active = false
		} else {
			p.active = true
			resetPath.childDisplays[label].active = true
		}
		return p, true
	}
	return nil, false
}

func AddDisplayNode(d *Display, toPath string, toAdd string, html []*Token) (bool, *Display) {
	// Always start from root.
	resetPath := verifyNode(toPath)
	if resetPath == nil {
		return false, nil
	}

	var newPath strings.Builder

	newPath.WriteString(toPath)
	newPath.WriteString("/")
	newPath.WriteString(toAdd)
	newNode := NewDisplay(toAdd, newPath.String(), html)
	resetPath.childDisplays[toAdd] = newNode
	resetPath.childDisplays[toAdd].prev = resetPath // Add backtrack to parent.

	parent := resetPath.path
	lastCoord := resetPath.coordinates

	resetPath.childDisplays[toAdd].coordinates = *mapping.AddCoordinate(
		"BOT", // TODO: Default for now just for testing.
		newNode.path,
		newNode.coordinates.Height,
		newNode.coordinates.Width,
		parent, lastCoord,
	)

	return true, resetPath.childDisplays[toAdd]
}

func getDisplay(label string, toPath string) (*server.BodyJSON, errors.ReturnCode) {
	display := RootPath

	deref := parseDeref(label, ".")
	finfo, ok := FTable[deref]
	if !ok {
		fmt.Println("ok is first ", ok)
		fmt.Println("deref is ", deref)
		pinfo, ok := TTable[deref]
		if ok {
			fmt.Println("okay is now ", ok)
			finfo = &FunctionTable{}
			finfo.file = pinfo.file
			finfo.tokens = pinfo.tokens
			visualize(finfo.tokens)
			fmt.Println("==============")
			body := BodyToJson(label, toPath, nil, nil)
			p, ok := isActive(toPath, label)

			if ok {
				if !p.active {
					allToState(p, false)
				}
				body.Tokens = TokenToJson(p.rawTokens)
				body.Path = finfo.file
				body.Coord = CoordinatesToJson(p.coordinates)
				return body, errors.ReturnCodeStateChange
			}

			success, p := AddDisplayNode(display, toPath, label, finfo.tokens)
			if !success {
				return nil, errors.ErrorCodeAddDisplayNode
			}

			body.Tokens = TokenToJson(p.rawTokens)
			body.Path = finfo.file
			body.Coord = CoordinatesToJson(p.coordinates)
			return body, errors.ReturnCodeDisplayDoesntExistNewDisplayAdded
		}
		fmt.Println("and finally it is now ", ok)
	}

	if ok {
		body := BodyToJson(label, toPath, nil, nil)
		p, ok := isActive(toPath, label)

		if ok {
			if !p.active {
				allToState(p, false)
			}
			body.Tokens = TokenToJson(p.rawTokens)
			body.Path = finfo.file
			body.Coord = CoordinatesToJson(p.coordinates)
			return body, errors.ReturnCodeStateChange
		}

		success, p := AddDisplayNode(display, toPath, label, finfo.tokens)
		if !success {
			return nil, errors.ErrorCodeAddDisplayNode
		}

		body.Tokens = TokenToJson(p.rawTokens)
		body.Path = finfo.file
		body.Coord = CoordinatesToJson(p.coordinates)
		return body, errors.ReturnCodeDisplayDoesntExistNewDisplayAdded
	}

	return nil, errors.ReturnCodeFunctionToDisplayNotFound
}
