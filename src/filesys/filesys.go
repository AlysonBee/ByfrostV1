package filesys

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

var Files = make(map[string]bool)

func pushFile(dirname string) {
	Files[dirname] = true
}

func String() {
	for file := range Files {
		fmt.Println(file)
	}
}

func checkExtensions(file string, ext []string) bool {
	for _, e := range ext {
		if strings.HasSuffix(file, e) {
			return true
		}
	}
	return false
}

func RecurseProject(dirName string, ext []string) {
	var (
		directories []string
		currentList []string
		outerIndex  int
		fileName    os.FileInfo
	)

	outerIndex = 0
	currentList = append(currentList, dirName)
	for outerIndex < len(currentList) {

		dir := currentList[outerIndex]

		index := 0
		files, err := ioutil.ReadDir(dir)
		if err != nil {
			outerIndex++
			continue
		}

		for index < len(files) {
			fileName = files[index]
			fi, err := os.Stat(dir + "/" + fileName.Name())
			if err != nil {
				index++
				continue
			}

			mode := fi.Mode()
			if mode.IsDir() {
				directories = append(directories, dir+"/"+fileName.Name())
			} else {
				if checkExtensions(dir+"/"+fileName.Name(), ext) {
					pushFile(dir + "/" + fileName.Name())
				}
			}
			index++
		}

		outerIndex++

		if outerIndex == len(currentList) && len(directories) > 0 {
			currentList = directories
			directories = []string{}
			outerIndex = 0
		}
	}
}
