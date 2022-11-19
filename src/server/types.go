package server

type CoordinatesJSON struct {
	Height int `json:"height"`
	Width  int `json:"width"`
	XCoord int `json:"xcoord"`
	YCoord int `json:"ycoord"`
}

type BodyJSON struct {
	Label  string `json:"label"`
	Path   string `json:"path"`
	Tokens []*TokenJSON
	Coord  *CoordinatesJSON
}

type TokenJSON struct {
	Name     string `json:"name"`
	Filepath string `json:"filepath"`
	Typ      string `json:"typ"`
	Line     int    `json:"line"`
	Spaces   int    `json:"spaces"`
	Tabs     int    `json:"tabs"`

	// For function call tokens only
	Label    string `json:"label"`
	Id       string `json:"id"`
	StyleTag string `json:"styletag"`
}
