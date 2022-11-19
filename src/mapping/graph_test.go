package mapping

import (
	"fmt"
	"testing"
)

// func TestCheckCrosses(t *testing.T) {
// 	var tests = []struct {
// 		name    string
// 		x       int
// 		y       int
// 		height  int
// 		width   int
// 		x2      int
// 		y2      int
// 		height2 int
// 		width2  int
// 	}{
// 		{
// 			name:    "Both blocks are in the same position",
// 			x:       10,
// 			y:       10,
// 			height:  60,
// 			width:   60,
// 			x2:      10,
// 			y2:      10,
// 			height2: 5,
// 			width2:  5,
// 		},
// 	}

// 	for _, tc := range tests {
// 		t.Run(tc.name, func(t *testing.T) {
// 			checkCrosses(tc.x, tc.y, tc.height, tc.width, tc.x2, tc.y2,
// 				tc.height2, tc.width2)
// 			t.Fail()
// 		})
// 	}
// }

func TestCollisionDetection(t *testing.T) {
	CoordinatesList = []*OccupiedBlock{
		{
			functionName: "test1",
			corner1:      []int{0, 0},
			corner2:      []int{0, 30},
			corner3:      []int{10, 30},
			corner4:      []int{10, 0},
		},
		{
			functionName: "test2",
			corner1:      []int{12, 31},
			corner2:      []int{12, 40},
			corner3:      []int{60, 40},
			corner4:      []int{60, 31},
		},
		{
			functionName: "test3",
			corner1:      []int{12, 20},
			corner2:      []int{12, -5},
			corner3:      []int{33, -5},
			corner4:      []int{33, 20},
		},
	}

	var tests = []struct {
		name   string
		height int
		width  int
	}{
		{
			name:   "A small function",
			height: 11,
			width:  22,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			fmt.Println("=================")
			collisionDetection(tc.height, tc.width)
			fmt.Println("================")
			t.Fail()
		})
	}
}
