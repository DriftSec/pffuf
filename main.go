package main

import (
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/manifoldco/promptui"
)

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
	var inputPath string

	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("     ./pffuf [path to ffuf JSON files]")
		fmt.Println("")
	}

	flag.Parse()

	if flag.NArg() > 0 {
		inputPath = flag.Args()[0]
		if !strings.HasSuffix(inputPath, "/") {
			inputPath = inputPath + "/"
		}
	} else {
		inputPath = "./"
	}

	inputFiles = getFilelist(inputPath + "*.json")

	if len(inputFiles) <= 0 {
		fmt.Println("No JSON files found !!!")
		flag.Usage()
		os.Exit(1)
	}

	for _, curfile := range inputFiles {
		// read the whole file at once
		b, err := ioutil.ReadFile(curfile)
		if err != nil {
			panic(err)
		}
		s := string(b)
		// //check whether s contains substring text
		if !strings.Contains(s, "ffuf") && !strings.Contains(s, "results") {
			fmt.Println(curfile, "does appear to be an ffuf JSON file, skipping")
			continue
		}

		getEndpoints(curfile)
		getCommands(curfile)

	}

	if len(commands) <= 0 && len(results) <= 0 {
		fmt.Println("Nothing to parse, exiting")
		os.Exit(1)
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

	results = results[:0]
	results = origresults
	doSort()
	doFilter()

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
				// if !checkFilter(ep) {
				// 	continue
				// }
				fmt.Println(ep.Endpoint)
				curScreen = append(curScreen, ep.Endpoint)
			}
		}

		if result == "u" || result == "urls" {
			curScreen = curScreen[:0]
			for _, ep := range results {
				// if !checkFilter(ep) {
				// 	continue
				// }
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
				// if !checkFilter(ep) {
				// 	continue
				// }
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
