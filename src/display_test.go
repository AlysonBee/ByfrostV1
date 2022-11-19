package main

import (
	"fmt"
	"indexer/styling"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func seedFTable() {
	pushFtable("main", "main.c", []*Token{})
	pushFtable("read_file", "main.c", []*Token{})
	pushSTable("write", "main.c", []*Token{})
	pushSTable("helloworld", "main.c", []*Token{})
}

func evalPath(expected string) (bool, string) {
	resetPath := RootPath
	segments := strings.Split(expected, ".")

	var expectedString strings.Builder

	for _, segment := range segments {
		if len(segment) == 0 {
			continue
		}

		if doesDispalyExist(segment, resetPath) {
			resetPath = resetPath.childDisplays[segment]
			expectedString.WriteString(".")
			expectedString.WriteString(resetPath.name)
		} else {
			return false, ""
		}
	}
	fmt.Print(expectedString.String())
	return true, expectedString.String()
}

func getNode(path string) *Display {
	node := RootPath
	segments := strings.Split(path, ".")

	for _, segment := range segments {
		if len(segment) == 0 {
			continue
		}

		if doesDispalyExist(segment, node) {
			node = node.childDisplays[segment]
		} else {
			fmt.Printf("Error: %s: path not found\n", path)
			return nil
		}
	}
	return node
}

func TestIsActive(t *testing.T) {
	// Test function:
	// fn read() {
	//		return
	// }
	//
	// fn read_file() {
	//		helloworld()
	// }
	//
	// fn main() {
	// 	 	read_file()
	// 	 	write()
	// }

	RootPath = &Display{
		name:   "__init",
		path:   "",
		active: true,
		childDisplays: map[string]*Display{
			"main": {
				name:   "main",
				path:   "",
				active: true,
				childDisplays: map[string]*Display{
					"read_file": {
						name:   "read_file",
						path:   ".main",
						active: true,
						childDisplays: map[string]*Display{
							"helloworld": {
								name:          "helloworld",
								path:          ".main.read_file",
								active:        true,
								childDisplays: map[string]*Display{},
							},
						},
					},
					"write": {
						name:   "write",
						path:   ".main",
						active: true,
						childDisplays: map[string]*Display{
							"read": {
								name:          "read",
								path:          ".main.read_file",
								childDisplays: map[string]*Display{},
							},
						},
					},
				},
			},
		},
	}

	var tests = []struct {
		name         string
		path         string
		functionName string
		absPath      string
		isExpanded   bool
		expected     bool
	}{
		{
			name:         "Collapse test starting as true  - collapse helloworld",
			path:         ".main.read_file",
			functionName: "helloworld",
			absPath:      ".main.read_file.helloworld",
			isExpanded:   false,
			expected:     true,
		},
		{
			name:         "Collapse test starting as false - expand helloworld",
			path:         ".main.read_file",
			functionName: "helloworld",
			absPath:      ".main.read_file.helloworld",
			isExpanded:   true,
			expected:     true,
		},
		{
			name:         "Collapse test nonexistent node",
			path:         ".main.read_file",
			functionName: "nonexistent",
			expected:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual, boolean := isActive(tt.path, tt.functionName)
			require.Equal(t, boolean, tt.expected)

			if tt.expected {
				require.Equal(t, actual.active, tt.isExpanded)
				require.Equal(t, actual.path, tt.path)

				testNode := getNode(tt.absPath)
				require.Equal(t, tt.isExpanded, testNode.active)
			}
		})
	}
}

func TestAllToState(t *testing.T) {

	//	fn main() {
	//		read_file();
	//		write();
	//	}
	//
	//	fn read_file() {
	//		helloworld();
	//	}
	//
	// 	fn helloworld() {
	//		testing();
	//	}
	//
	//	fn testing() {
	// 		even_deeper();
	// 	}
	//
	//	fn even_deeper() {
	//		return ;
	//	}
	//
	//	fn write() {
	//		read();
	//	}
	//
	// 	fn read() {
	// 		malloc();
	//	}

	RootPath = &Display{
		name:   "__init",
		path:   "",
		active: true,
		childDisplays: map[string]*Display{
			"main": {
				name:   "main",
				path:   "",
				active: true,
				childDisplays: map[string]*Display{
					"read_file": {
						name:   "read_file",
						path:   ".main",
						active: true,
						childDisplays: map[string]*Display{
							"helloworld": {
								name:   "helloworld",
								path:   ".main.read_file",
								active: true,
								childDisplays: map[string]*Display{
									"testing": {
										name:   "testing",
										path:   ".main.read_file.helloworld",
										active: true,
										childDisplays: map[string]*Display{
											"even_deeper": {
												name:          "evendeeper",
												path:          ".main.read_file.helloworld.testing",
												active:        true,
												childDisplays: map[string]*Display{},
											},
										},
									},
								},
							},
						},
					},
					"write": {
						name:   "write",
						path:   ".main",
						active: true,
						childDisplays: map[string]*Display{
							"read": {
								name:   "read",
								path:   ".main.write",
								active: true,
								childDisplays: map[string]*Display{
									"malloc": {
										name:          "malloc",
										path:          ".main.write.read",
										active:        true,
										childDisplays: map[string]*Display{},
									},
								},
							},
						},
					},
				},
			},
		},
	}

	var tests = []struct {
		name          string
		toCollapse    string
		stateToSet    bool
		expectedPaths []string
	}{
		{
			name:       "Collapse everything - main",
			toCollapse: "main",
			stateToSet: false,
			expectedPaths: []string{
				".main",
				".main.read_file",
				".main.read_file.helloworld",
				".main.read_file.helloworld.testing",
				".main.read_file.helloworld.testing.even_deeper",
				".main.write",
				".main.write.read",
				".main.write.read.malloc",
			},
		},
		{
			name:       "Expand everything - main",
			toCollapse: "main",
			stateToSet: false,
			expectedPaths: []string{
				".main",
				".main.read_file",
				".main.read_file.helloworld",
				".main.read_file.helloworld.testing",
				".main.read_file.helloworld.testing.even_deeper",
				".main.write",
				".main.write.read",
				".main.write.read.malloc",
			},
		},
		{
			name:       "Expand .main.read_file.helloworld - helloword",
			toCollapse: ".main.read_file.helloworld",
			stateToSet: true,
			expectedPaths: []string{
				".main.read_file.helloworld",
				".main.read_file.helloworld.testing",
				".main.read_file.helloworld.testing.even_deeper",
			},
		},
		{
			name:       "Collapse .main.read_file.helloworld - helloword",
			toCollapse: ".main.read_file.helloworld",
			stateToSet: false,
			expectedPaths: []string{
				".main.read_file.helloworld",
				".main.read_file.helloworld.testing",
				".main.read_file.helloworld.testing.even_deeper",
			},
		},
		{
			name:       "Collapse .main.read_file.helloworld.testing.even_deeper - even_deeper",
			toCollapse: ".main.read_file.helloworld.testing.even_deeper",
			stateToSet: false,
			expectedPaths: []string{
				".main.read_file.helloworld.testing.even_deeper",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testNode := getNode(tt.toCollapse)
			allToState(testNode, tt.stateToSet)

			for _, expected := range tt.expectedPaths {
				check := getNode(expected)
				require.Equal(t, tt.stateToSet, check.active)
			}
		})
	}

}

func TestDisplay(t *testing.T) {
	var tests = []struct {
		name         string
		functionName string
		path         string
		absPath      string
		tokens       []*Token
		parentName   string

		expected string
	}{
		{
			name: `
				First Degree:
				path: .
				example:
					fn main() {
						return 0
					}
				`,
			functionName: "main",
			path:         "",
			tokens:       []*Token{},
			absPath:      ".main",
			parentName:   "__init",
			expected:     "",
		},
		{
			name: `
				Second Degree:
				path: main
				example:
					fn main() { 
						helloworld() 
					}
					
					fn helloworld() {
						return 
					}
				`,
			functionName: "helloworld",
			path:         ".main",
			absPath:      ".main.helloworld",
			tokens:       []*Token{},
			parentName:   "main",
			expected:     ".main",
		},
		{
			name: `
				Third Degree:
				path: main.helloworld.read_file
				example:
					fn read_file() {
						return 
					}

					fn helloworld() {
						read_file()
					}
					
					fn main() {
						helloworld()
					}
				`,
			functionName: "read_file",
			path:         ".main.helloworld",
			absPath:      ".main.helloworld.read_file",
			tokens:       []*Token{},
			parentName:   "helloworld",
			expected:     ".main.helloworld",
		},
		{
			name: `
				Second degree:
				path: main.test_file 
				example:
					fn test_file() {
						return
					},

					fn helloworld() {
						read_file()
					}

					fn main() { 
						test_file()
						helloworld()
					}
			`,
			functionName: "test_file",
			path:         ".main",
			absPath:      ".main.test_file",
			tokens:       []*Token{},
			parentName:   "main",
			expected:     ".main",
		},
	}

	styling.PrepSyntaxHighlighting("styling/configs/c.json")
	initDisplay()
	d := RootPath

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			AddDisplayNode(d, tt.path, tt.functionName, tt.tokens)
			boolean, actual := evalPath(tt.expected)

			require.Equal(t, actual, tt.expected)
			require.True(t, boolean)

			node := getNode(tt.absPath)
			require.Equal(t, node.prev.name, tt.parentName)
		})
	}
}
