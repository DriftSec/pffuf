package main

import "regexp"

// global debug mode used with env variable
var DEBUGMODE bool

/*
 * Ffuf JSON Structs
 */

type FfufOutput struct {
	CommandLine string       `json:"commandline"`
	Time        string       `json:"time"`
	Config      FfufMetaData `json:"config"`
	Results     []FfufResult `json:"results"`
}

type FfufInputProviders struct {
	Keyword string `json:"keyword"`
	Name    string `json:"name"`
	Value   string `json:"value"`
}

// ffuf json 'commandline' struct
type FfufMetaData struct {
	URL            string               `json:"url"`
	Method         string               `json:"method"`
	Outputfile     string               `json:"outputfile"`
	InputProviders []FfufInputProviders `json:"inputproviders"`
}

// ffuf json 'results' struct - status, size, words, lines
type FfufResult struct {
	Input struct {
		FUZZ string `json:"FUZZ"`
	} `json:"input"`
	// Input    string `json:"input"`
	Position int    `json:"position"`
	Status   int    `json:"status"`
	Length   int    `json:"length"`
	Words    int    `json:"words"`
	Lines    int    `json:"lines"`
	Host     string `json:"host"`
	URL      string `json:"url"`
}

/*
 * FuzzNav Structs
 */

type NavResults struct {
	Endpoint   string
	Status     int
	Length     int
	Words      int
	Lines      int
	URL        string
	Wordlist   string
	Outputfile string
}

type Filters struct {
	StatusMatch []int
	StatusHide  []int
	LenMatch    []int
	LenHide     []int
	WordsMatch  []int
	WordsHide   []int
	LinesHide   []int
	LinesMatch  []int
	RegExMatch  []*regexp.Regexp
	RegExHide   []*regexp.Regexp
}
