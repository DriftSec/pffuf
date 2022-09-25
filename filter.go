package main

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/bradfitz/slice"
)

func doSort() {
	sorting = strings.ToLower(sorting)
	if sorting == "status" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Status < results[j].Status
		})
	}

	if sorting == "length" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Length < results[j].Length
		})
	}

	if sorting == "words" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Words < results[j].Words
		})
	}

	if sorting == "lines" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Lines < results[j].Lines
		})
	}

	if sorting == "url" {
		slice.Sort(results[:], func(i, j int) bool { return results[i].URL < results[j].URL })
	}

	if sorting == "endpoint" {
		slice.Sort(results[:], func(i, j int) bool { return results[i].Endpoint < results[j].Endpoint })
	}

	if sorting == "none" {
		results = results[:0]
		results = origresults
		doFilter()

	}
}

func setFilter(cmd, args string) {
	// parts := strings.Split(cmdline, " ")
	// cmd := parts[0]
	// args := strings.Join(parts[1:], "")

	argsi := []int{}
	argsiRegEx := []*regexp.Regexp{}
	for _, arg := range strings.Split(args, ",") {
		if arg == " " || arg == "" {
			continue
		}
		if cmd == "mr" || cmd == "fr" {
			r, err := regexp.Compile(arg)
			if err != nil {
				fmt.Printf("%s is not a valid regular expresssion\n", arg)
				return
			}
			argsiRegEx = append(argsiRegEx, r)
		} else {
			argi, err := strconv.Atoi(arg)
			if err != nil {
				fmt.Printf("%s is not a valid filter condition\n", arg)
				return
			}
			argsi = append(argsi, argi)
		}
	}
	if cmd == "fc" {
		filters.StatusHide = argsi
	}
	if cmd == "fw" {
		filters.WordsHide = argsi
	}
	if cmd == "fl" {
		filters.LinesHide = argsi
	}
	if cmd == "fs" {
		filters.LenHide = argsi
	}
	if cmd == "mc" {
		filters.StatusMatch = argsi
	}
	if cmd == "mw" {
		filters.WordsMatch = argsi
	}
	if cmd == "ml" {
		filters.LenMatch = argsi
	}
	if cmd == "ms" {
		filters.LinesMatch = argsi
	}
	if cmd == "mr" {
		filters.RegExMatch = argsiRegEx
	}
	if cmd == "fr" {
		filters.RegExHide = argsiRegEx
	}
	doFilter()
}

func showFilters() {
	fmt.Println("Current Filters:")
	fmt.Println("Hide STATUS: ", filters.StatusHide)
	fmt.Println("Hide WORDS: ", filters.WordsHide)
	fmt.Println("Hide LINES: ", filters.LinesHide)
	fmt.Println("Hide LENGTH: ", filters.LenHide)
	fmt.Println("Hide REGEX: ", filters.RegExHide)
	fmt.Println("Match STATUS: ", filters.StatusMatch)
	fmt.Println("Match WORDS: ", filters.WordsMatch)
	fmt.Println("Match LINES: ", filters.LinesMatch)
	fmt.Println("Match LENGTH: ", filters.LenMatch)
	fmt.Println("Match REGEX: ", filters.RegExMatch)
}
func clearFilters() {
	filters.StatusHide = filters.StatusHide[:0]
	filters.WordsHide = filters.WordsHide[:0]
	filters.LinesHide = filters.LinesHide[:0]
	filters.LenHide = filters.LenHide[:0]
	filters.StatusMatch = filters.StatusMatch[:0]
	filters.WordsMatch = filters.WordsMatch[:0]
	filters.LinesMatch = filters.LinesMatch[:0]
	filters.LenMatch = filters.LenMatch[:0]
	filters.RegExMatch = filters.RegExMatch[:0]
	filters.RegExHide = filters.RegExHide[:0]
	results = results[:0]
	results = origresults
	doSort()
	fmt.Println("All Filters Cleared !")
}

func ifMatch(line NavResults) bool {
	if (len(filters.RegExMatch) == 0 || containsRegExp(filters.RegExMatch, line.URL)) && (len(filters.StatusMatch) == 0 || containsInt(filters.StatusMatch, line.Status)) && (len(filters.LenMatch) == 0 || containsInt(filters.LenMatch, line.Length)) && (len(filters.WordsMatch) == 0 || containsInt(filters.WordsMatch, line.Words)) && (len(filters.LinesMatch) == 0 || containsInt(filters.LinesMatch, line.Lines)) {
		return true
	}

	return false
}

func ifHide(line NavResults) bool {
	if len(filters.RegExHide) > 0 && containsRegExp(filters.RegExHide, line.URL) {
		return true
	}
	if len(filters.StatusHide) > 0 && containsInt(filters.StatusHide, line.Status) {
		return true
	}
	if len(filters.LenHide) > 0 && containsInt(filters.LenHide, line.Length) {
		return true
	}
	if len(filters.WordsHide) > 0 && containsInt(filters.WordsHide, line.Words) {
		return true
	}
	if len(filters.LinesHide) > 0 && containsInt(filters.LinesHide, line.Lines) {
		return true
	}
	return false
}

func doFilter() {
	results = results[:0]
	results = origresults

	var resultstmp []NavResults

	for _, cur := range results {
		if len(filters.StatusMatch) > 0 || len(filters.LenMatch) > 0 || len(filters.WordsMatch) > 0 || len(filters.LinesMatch) > 0 || len(filters.RegExMatch) > 0 {
			if ifMatch(cur) {
				resultstmp = append(resultstmp, cur)
			}
		} else {
			if !ifHide(cur) {
				resultstmp = append(resultstmp, cur)
			}
		}
	}
	results = resultstmp
	doSort()
}
