package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

const (
	TERMINAL_CLEAR_LINE = "\r\x1b[2K"
	ANSI_CLEAR          = "\x1b[0m"
	ANSI_RED            = "\x1b[31m"
	ANSI_GREEN          = "\x1b[32m"
	ANSI_BLUE           = "\x1b[34m"
	ANSI_YELLOW         = "\x1b[33m"
)

func getFilelist(inputArg string) []string {
	// pattern := inputArg
	var retval []string
	matches, _ := filepath.Glob(inputArg)
	for _, match := range matches {
		retval = append(retval, match)
	}
	return retval
}

func getFilelistRecursive(inputArg string) []string {
	var retval []string
	err := filepath.Walk(inputArg,
		func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if strings.HasSuffix(path, ".json") {
				retval = append(retval, path)
			}
			return nil
		})
	if err != nil {
		log.Fatal(err)
	}

	return retval
}

// func strcontains(filterlst []string, cur string) bool {
// 	for _, v := range filterlst {
// 		if strings.Contains(cur, v) {
// 			return true
// 		}
// 	}

// 	return false
// }

func containsRegExp(filterlst []*regexp.Regexp, cur string) bool {
	for _, r := range filterlst {
		if r.MatchString(cur) {
			return true
		}
	}
	return false
}
func containsStr(filterlst []string, cur string) bool {
	for _, v := range filterlst {
		if v == cur {
			return true
		}
	}

	return false
}
func containsInt(filterlst []int, cur int) bool {
	for _, v := range filterlst {
		if v == cur {
			return true
		}
	}

	return false
}

func colorize(input string, status int) string {

	colorCode := ANSI_CLEAR
	if status >= 200 && status < 300 {
		colorCode = ANSI_GREEN
	}
	if status >= 300 && status < 400 {
		colorCode = ANSI_BLUE
	}
	if status >= 400 && status < 500 {
		colorCode = ANSI_YELLOW
	}
	if status >= 500 && status < 600 {
		colorCode = ANSI_RED
	}
	return fmt.Sprintf("%s%s%s", colorCode, input, ANSI_CLEAR)
}

func writeFile(filename string) {
	var sz int
	_, err := os.Stat(filename)
	if !os.IsNotExist(err) {
		if !promptYesNo("File exsists, overwrite?") {
			return
		}
	}

	f, err := os.Create(filename)

	if err != nil {
		fmt.Println("Failed to create ", filename, "!!!")
	}

	defer f.Close()

	for _, line := range curScreen {
		outstr := fmt.Sprintf("%s\n", line)
		sz += len(outstr)
		_, err2 := f.WriteString(outstr)
		if err2 != nil {
			fmt.Println("Failed writing to ", filename, "!!!")
		}
	}

	fmt.Println("Wrote ", sz, "bytes to", filename)

}
