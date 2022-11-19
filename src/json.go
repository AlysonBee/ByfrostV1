package main

import (
	"bytes"
	"encoding/json"
	"byfrostV1/mapping"
	"byfrostV1/server"
)

func move(token *Token) *server.TokenJSON {
	return &server.TokenJSON{
		Name:     token.name,
		Filepath: token.filepath,
		Typ:      token.typ,
		Line:     token.line,
		Spaces:   token.spaces,
		Tabs:     token.tabs,
		Label:    token.label,
		Id:       token.id,
		StyleTag: token.styleTag,
	}
}

func structToJSON(data interface{}) ([]byte, error) {
	buf := new(bytes.Buffer)

	if err := json.NewEncoder(buf).Encode(data); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func TokenToJson(tokens []*Token) []*server.TokenJSON {
	var toWeb []*server.TokenJSON

	for _, t := range tokens {
		toWeb = append(toWeb, move(t))
	}

	return toWeb
}

func BodyToJson(label string, path string,
	jsonTokens []*server.TokenJSON, coordinates *server.CoordinatesJSON) *server.BodyJSON {

	return &server.BodyJSON{
		Label:  label,
		Path:   path,
		Coord:  coordinates,
		Tokens: jsonTokens,
	}
}

func CoordinatesToJson(coordinate mapping.Coordinate) *server.CoordinatesJSON {
	return &server.CoordinatesJSON{
		Height: coordinate.Height,
		Width:  coordinate.Width,
		XCoord: coordinate.XCoord,
		YCoord: coordinate.YCoord,
	}
}
