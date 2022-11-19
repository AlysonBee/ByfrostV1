package mapping

type Coordinate struct {
	Height int
	Width  int
	XCoord int
	YCoord int
}

var UsedPositions = [][]int{}

func NewCoordinate(height int, width int, xCoordinate int, yCoordinate int) *Coordinate {
	return &Coordinate{
		Height: height,
		Width:  width,
		XCoord: xCoordinate,
		YCoord: yCoordinate,
	}
}

func checkIfTaken(xcoord int, ycoord int) bool {
	for _, pos := range UsedPositions {
		if len(pos) == 2 && pos[0] != 0 && pos[1] != 0 {
			if xcoord == pos[0] && ycoord == pos[1] {
				return true
			}
		}
	}
	return false
}

func AddCoordinate(direction string, functionName string, height int,
	width int, parentName string, lastItem Coordinate) *Coordinate {
	var (
		xCoord int
		yCoord int
	)
	xCoord = lastItem.XCoord + lastItem.Width
	yCoord = lastItem.YCoord + lastItem.Height

	for {
		if direction == "TOP" {
			xCoord += 0
			yCoord -= 0
		} else if direction == "BOT" {
			xCoord += 400
			yCoord -= 0
		}

		if !checkIfTaken(xCoord, yCoord) {
			UsedPositions = append(UsedPositions, []int{xCoord, yCoord})
			break
		}

		if direction == "TOP" {
			direction = "BOT"
		} else {
			direction = "TOP"
		}
	}

	return NewCoordinate(height, width, xCoord, yCoord)
}
