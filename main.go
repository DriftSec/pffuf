package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"net/url"
	"os"
	"regexp"
	"strings"
	"time"

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
	Cl          string
)

func init() {
	initfunc()
}

func initfunc() {
	inputFiles = inputFiles[:0]
	curScreen = curScreen[:0]
	results = results[:0]
	origresults = origresults[:0]
	commands = commands[:0]

	var inputPath string

	flag.Usage = func() {
		fmt.Println("")
		fmt.Println("Usage:")
		fmt.Println("./pffuf [-cl 'cmds'] [path to ffuf JSON files]")
		fmt.Println("      -cl [cmds]    Run commands and exit, used for scripting  (i.e. -cl 'mr .*?\\.php;u' to regex for php and print urls)")
		fmt.Println("      -r            Recursively search for JSON files")

		fmt.Println("")
	}

	commandLine := flag.String("cl", "", "Run command line and exit")
	recurse := flag.Bool("r", false, "Recursively search for .json files")
	flag.Parse()

	if flag.NArg() > 0 {
		inputPath = flag.Args()[0]
		if !strings.HasSuffix(inputPath, "/") {
			inputPath = inputPath + "/"
		}
	} else {
		inputPath = "./"
	}
	if *recurse {
		inputFiles = getFilelistRecursive(inputPath)
	} else {
		inputFiles = getFilelist(inputPath + "*.json")
	}

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

	if *commandLine != "" {
		cmds := strings.Split(*commandLine, ";")
		for _, cmd := range cmds {
			parseCommand(strings.TrimLeft(cmd, " "))
		}
		os.Exit(0)
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

func selectSort(arg string) {
	if arg == "" {
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

	} else {
		tmp := []string{"status", "length", "words", "lines", "url", "endpoint", "none"}
		if !containsStr(tmp, arg) {
			fmt.Println(arg, "is not a valid sort option.(", tmp, ")")
			return
		}
		sorting = arg
	}
	results = results[:0]
	results = origresults
	doSort()
	doFilter()

}

func joinResults(filename string) {
	var (
		sz   int
		outs FfufOutput
	)

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

	outs.CommandLine = "Created by pffuf"
	currentTime := time.Now()

	outs.Time = currentTime.Format("2006-01-02T15:04:05Z") //"2021-04-02T17:41:37Z"
	// outs.Config.URL            string               `json:"url"`
	// outs.Config.Method         string               `json:"method"`
	outs.Config.Outputfile = filename
	// outs.Config.InputProviders []FfufInputProviders `json:"inputproviders"`

	for _, ep := range results {
		var ffufresult FfufResult
		ffufresult.Input.FUZZ = ep.Endpoint
		ffufresult.Position = 0 //<<<<<<<<<<<<<<<<<<???????????????????????
		ffufresult.Status = ep.Status
		ffufresult.Length = ep.Length
		ffufresult.Words = ep.Words
		ffufresult.Lines = ep.Lines
		ffufresult.URL = ep.URL
		u, _ := url.Parse(ep.URL)
		ffufresult.Host = u.Host

		outs.Results = append(outs.Results, ffufresult)
	}
	u, err := json.Marshal(outs)
	if err != nil {
		fmt.Println("failed to marshal JSON:", err)
	}
	sz, err = f.Write(u)
	if err != nil {
		fmt.Println("Failed to create ", filename, "!!!")
	} else {
		fmt.Println("Wrote", sz, "bytes to ", filename)
	}

}

func grep(expr string) {
	var tmpscreen []string
	tmpscreen = curScreen //append(tmpscreen, curScreen)
	curScreen = curScreen[:0]
	for _, line := range tmpscreen {
		matched, err := regexp.MatchString(expr, line)

		if err != nil {
			fmt.Println("Invalid Expression !!!")
		} else if matched {
			fmt.Println(line)
			curScreen = append(curScreen, line)
		}
	}
}

func grepv(expr string) {
	var tmpscreen []string
	tmpscreen = curScreen //append(tmpscreen, curScreen)
	curScreen = curScreen[:0]
	for _, line := range tmpscreen {
		matched, err := regexp.MatchString(expr, line)

		if err != nil {
			fmt.Println("Invalid Expression !!!")
		} else if !matched {
			fmt.Println(line)
			curScreen = append(curScreen, line)
		}
	}
}

func help() {

	fmt.Println("Commands:")
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
	fmt.Println("fr [regex1,regex2]  Filter URL using regex")
	fmt.Println("mc [val,val2]       Match status code")
	fmt.Println("mw [val,val2]       Match words")
	fmt.Println("ml [val,val2]       Match lines")
	fmt.Println("ms [val,val2]       Match length")
	fmt.Println("mr [regex1,regex2]  Match URL using regex")
	fmt.Println("s|sort              Sort options")
	fmt.Println("g|grep [expr]       Run grep on last output")
	fmt.Println("gv|grepv [expr]     Run grep exclude on last output")
	fmt.Println("r|reload            Reparse input path for new files")
	fmt.Println("j|join              Combine all filtered results and export to ffuf JSON file.")

}

func main() {
	for {
		validate := func(input string) error {
			validCMDs := []string{"j", "join", "r", "reload", "gv", "grepv", "g", "grep", "t", "tree", "s", "sort", "h", "help", "sf", "show-filters", "cf", "clear-filters", "fc", "fw", "fl", "fs", "fr", "mc", "mw", "ml", "ms", "mr", "c", "commands", "x", "exit", "e", "endpoints", "u", "urls", "d", "details", "w", "write"}
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

		parseCommand(result)
	}
}

func parseCommand(singleCommand string) {
	result := singleCommand
	// commands with args go here (except filters)
	parts := strings.Split(result, " ")
	cmd := parts[0]
	args := ""
	if len(parts) > 1 {
		args = parts[1]
	}

	if cmd == "s" || cmd == "sort" {
		selectSort(args)
	}

	if cmd == "w" || cmd == "write" {
		if args == "" || args == " " {
			fmt.Println("write: you must specify an output file !!")
		} else {
			writeFile(args)

		}
	}

	if cmd == "j" || cmd == "join" {
		if args == "" || args == " " {
			fmt.Println("write: you must specify an output file !!")
		} else {
			joinResults(args)

		}
	}

	if cmd == "g" || cmd == "grep" {
		if args == "" || args == " " {
			fmt.Println("write: you must specify an expression !!")
		} else {
			grep(args)

		}
	}
	if cmd == "gv" || cmd == "grepv" {
		if args == "" || args == " " {
			fmt.Println("write: you must specify an expression !!")
		} else {
			grepv(args)

		}
	}

	if cmd == "r" || cmd == "reload" {
		initfunc()
	}
	if cmd == "h" || cmd == "help" {
		help()
	}

	if cmd == "t" || cmd == "tree" {
		doTreePlain(results)
	}

	// if cmd == "s" || cmd == "sort" {
	// 	selectSort("")
	// }

	if cmd == "x" || cmd == "exit" {
		os.Exit(0)
	}

	if cmd == "x" || cmd == "exit" {
		os.Exit(0)
	}

	if cmd == "sf" || cmd == "show-filters" {
		showFilters()
	}

	if cmd == "cf" || cmd == "clear-filters" {
		clearFilters()
	}

	filterCMDs := []string{"fc", "fw", "fl", "fs", "fr", "mc", "mw", "ml", "ms", "mr"}

	if containsStr(filterCMDs, cmd) {
		setFilter(cmd, args)
	}

	if cmd == "c" || cmd == "commands" {
		curScreen = curScreen[:0]
		for _, cmd := range commands {
			fmt.Println(cmd)
			curScreen = append(curScreen, cmd)
		}
	}

	if cmd == "e" || cmd == "endpoints" {
		curScreen = curScreen[:0]
		for _, ep := range results {
			// if !checkFilter(ep) {
			// 	continue
			// }
			fmt.Println(ep.Endpoint)
			curScreen = append(curScreen, ep.Endpoint)
		}
	}

	if cmd == "u" || cmd == "urls" {
		curScreen = curScreen[:0]
		for _, ep := range results {
			// if !checkFilter(ep) {
			// 	continue
			// }
			fmt.Println(ep.URL)
			curScreen = append(curScreen, ep.URL)
		}
	}

	if cmd == "d" || cmd == "details" {
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
