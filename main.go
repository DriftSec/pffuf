package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"

	"sort"

	"github.com/bradfitz/slice"
	"github.com/manifoldco/promptui"
)

const (
	TERMINAL_CLEAR_LINE = "\r\x1b[2K"
	ANSI_CLEAR          = "\x1b[0m"
	ANSI_RED            = "\x1b[31m"
	ANSI_GREEN          = "\x1b[32m"
	ANSI_BLUE           = "\x1b[34m"
	ANSI_YELLOW         = "\x1b[33m"
)

type Filters struct {
	StatusMatch []int
	StatusHide  []int
	LenMatch    []int
	LenHide     []int
	WordsMatch  []int
	WordsHide   []int
	LinesHide   []int
	LinesMatch  []int
}

var (
	filters     Filters
	inputFiles  []string
	curScreen   []string
	results     []NavResults
	origresults []NavResults
	commands    []string
	sorting     string
)

func init() {
	var inputFile string
	// flag.StringVar(&inputFile, "f", "*.json", "json files")
	flag.Parse()
	inputFile = flag.Args()[0]
	inputFiles = getFilelist(inputFile)

	for _, curfile := range inputFiles {
		getEndpoints(curfile)
		getCommands(curfile)

	}

	fmt.Println("Loaded", len(inputFiles), "JSON files.")
	fmt.Println(len(commands), "ffuf commands")
	fmt.Println(len(results), "endpoints discovered")
	fmt.Println("")

}
func promptYesNo(question string) bool {
	validate := func(input string) error {
		if input != "y" && input != "n" && input != "Y" && input != "N" {
			return errors.New("Must answer with y/n/Y/N")
		}
		return nil
	}

	prompt := promptui.Prompt{
		Label:    question,
		Validate: validate,
	}

	result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
	}
	if result == "y" || result == "Y" {
		return true
	}
	return false

}

func selectSort() {
	prompt := promptui.Select{
		Label: "Sort by:",
		Items: []string{"Status", "Length", "Words", "Lines", "URL", "Endpoint", "None"},
	}

	_, result, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return
	}

	sorting = result

	if sorting == "Status" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Status < results[j].Status
		})
	}

	if sorting == "Length" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Length < results[j].Length
		})
	}

	if sorting == "Words" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Words < results[j].Words
		})
	}

	if sorting == "Lines" {
		sort.Slice(results[:], func(i, j int) bool {
			return results[i].Lines < results[j].Lines
		})
	}

	if sorting == "URL" {
		slice.Sort(results[:], func(i, j int) bool { return results[i].URL < results[j].URL })
	}

	if sorting == "Endpoint" {
		slice.Sort(results[:], func(i, j int) bool { return results[i].Endpoint < results[j].Endpoint })
	}

	if sorting == "None" {
		results = results[:0]
		results = origresults

	}

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

func help() {

	fmt.Println("c|commands          List ffuf commands that have been run")
	fmt.Println("x|exit              Quit")
	fmt.Println("e|endpoints         List endpoints")
	fmt.Println("u|urls              List full URLs")
	fmt.Println("d|details           Show endpoint details (status,words,lines,length)")
	fmt.Println("t|tree              Show a treeview of endpoints (glitchy with multiple vhosts, no write to file)")
	fmt.Println("w|write [filename]  write last output to file")
	fmt.Println("sf|show-filters     Show current filters")
	fmt.Println("cf|clear-filters    Clear all filters")
	fmt.Println("fc [val,val2]       Filter Status code")
	fmt.Println("fw [val,val2]       Filter by words")
	fmt.Println("fl [val,val2]       Filter lines")
	fmt.Println("fs [val,val2]       Filter lenght")
	fmt.Println("mc [val,val2]       Match status code")
	fmt.Println("mw [val,val2]       Match words")
	fmt.Println("ml [val,val2]       Match lines")
	fmt.Println("ms [val,val2]       Match length")
	fmt.Println("s|sort                sort ")

}

func main() {
	for {
		validate := func(input string) error {
			validCMDs := []string{"t", "tree", "s", "sort", "h", "help", "sf", "show-filters", "cf", "clear-filters", "fc", "fw", "fl", "fs", "mc", "mw", "ml", "ms", "c", "commands", "x", "exit", "e", "endpoints", "u", "urls", "d", "details", "w", "write"}
			cmd := strings.Split(input, " ")[0]
			for _, a := range validCMDs {
				if a == cmd {
					return nil
				}
			}

			return errors.New("Invalid command (try help)")

		}

		prompt := promptui.Prompt{
			Label:    ">",
			Validate: validate,
		}

		result, err := prompt.Run()

		if err != nil {
			fmt.Printf("Prompt failed %v\n", err)
			return
		}

		// commands with args go here (except filters)
		cmd := ""
		args := ""
		if strings.Contains(result, " ") {
			parts := strings.Split(result, " ")
			cmd = parts[0]
			args = parts[1]

			if cmd == "w" || cmd == "write" {
				if args == "" || args == " " {
					fmt.Println("write: you must specify an output file !!")
				} else {
					writeFile(args)

				}
			}
		}

		if result == "w" || result == "write" {
			fmt.Println("write: you must specify an output file !!")
		}

		if result == "h" || result == "help" {
			help()
		}

		if result == "t" || result == "tree" {
			doTreePlain(results)
		}

		if result == "s" || result == "sort" {
			selectSort()
		}

		if result == "x" || result == "exit" {
			os.Exit(0)
		}

		if result == "x" || result == "exit" {
			os.Exit(0)
		}

		if result == "sf" || result == "show-filters" {
			showFilters()
		}

		if result == "cf" || result == "clear-filters" {
			clearFilters()
		}

		filterCMDs := []string{"fc", "fw", "fl", "fs", "mc", "mw", "ml", "ms"}
		for _, fc := range filterCMDs {
			if fc == strings.Split(result, " ")[0] {
				setFilter(result)
			}
		}

		if result == "c" || result == "commands" {
			for _, cmd := range commands {
				fmt.Println(cmd)
			}
		}

		if result == "e" || result == "endpoints" {
			curScreen = curScreen[:0]
			for _, ep := range results {
				if !checkFilter(ep) {
					continue
				}
				fmt.Println(ep.Endpoint)
				curScreen = append(curScreen, ep.Endpoint)
			}
		}

		if result == "u" || result == "urls" {
			curScreen = curScreen[:0]
			for _, ep := range results {
				if !checkFilter(ep) {
					continue
				}
				fmt.Println(ep.URL)
				curScreen = append(curScreen, ep.URL)
			}
		}

		if result == "d" || result == "details" {
			curScreen = curScreen[:0]
			maxLen := 0
			// sorts results by status

			for _, eptest := range results {
				if len(eptest.URL) > maxLen {
					maxLen = len(eptest.URL)
				}
			}

			for _, ep := range results {
				if !checkFilter(ep) {
					continue
				}
				indent := strings.Repeat(" ", (maxLen + 5 - len(ep.URL)))
				var res_hdr string
				res_hdr = fmt.Sprintf("%s%s[Status: %d, Size: %d, Words: %d, Lines: %d]", ep.URL, indent, ep.Status, ep.Length, ep.Words, ep.Lines)
				res_hdr = colorize(res_hdr, ep.Status)
				fmt.Printf("%s\n", res_hdr)
				curScreen = append(curScreen, res_hdr)
			}
		}

	}
}

func setFilter(cmdline string) {
	parts := strings.Split(cmdline, " ")
	cmd := parts[0]
	args := strings.Join(parts[1:], "")

	argsi := []int{}
	for _, arg := range strings.Split(args, ",") {
		if arg == " " || arg == "" {
			continue
		}
		argi, err := strconv.Atoi(arg)
		if err != nil {
			fmt.Printf("%s is not a valid filter condition\n", arg)
			return
		}
		argsi = append(argsi, argi)
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
}

func showFilters() {
	fmt.Println("Current Filters:")
	fmt.Println("Hide STATUS: ", filters.StatusHide)
	fmt.Println("Hide WORDS: ", filters.WordsHide)
	fmt.Println("Hide LINES: ", filters.LinesHide)
	fmt.Println("Hide LENGTH: ", filters.LenHide)
	fmt.Println("Match STATUS: ", filters.StatusMatch)
	fmt.Println("Match WORDS: ", filters.WordsMatch)
	fmt.Println("Match LINES: ", filters.LinesMatch)
	fmt.Println("Match LENGTH: ", filters.LenMatch)
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
	fmt.Println("All Filters Cleared !")
}

func ifMatch(line NavResults) bool {
	if contains(filters.StatusMatch, line.Status) {
		return true
	}
	if contains(filters.LenMatch, line.Length) {
		return true
	}
	if contains(filters.WordsMatch, line.Words) {
		return true
	}
	if contains(filters.LinesMatch, line.Lines) {
		return true
	}

	return false
}

// checks the current filters, true = show it, false = hide it
func checkFilter(line NavResults) bool {
	if ifMatch(line) {
		return true
	} else {

		if contains(filters.StatusHide, line.Status) {
			return false
		}
		if contains(filters.LenHide, line.Length) {
			return false
		}
		if contains(filters.WordsHide, line.Words) {
			return false
		}
		if contains(filters.LinesHide, line.Lines) {
			return false
		}
	}
	if len(filters.LinesMatch) == 0 && len(filters.StatusMatch) == 0 && len(filters.WordsMatch) == 0 && len(filters.LenMatch) == 0 {
		return true
	}

	return false
}

func contains(filterlst []int, cur int) bool {
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
