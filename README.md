Pffuf is a tool for parsing FFUF json files. 
- filter, sort results from multiple files
- combine multiple files into one
- various outputs full urls, details, endpoints only,tree view, commands 
- interactive or command line modes


### Usage
```
Usage:
./pffuf [-cl 'cmds'] [path to ffuf JSON files]
      -cl [cmds]    Run commands and exit, used for scripting  (i.e. -cl 'mr .*?\.php;u' to regex for php and print urls)
```
### Available Commands
```
Commands:
c|commands          List ffuf commands that have been run
x|exit              Quit
e|endpoints         List endpoints
u|urls              List full URLs
d|details           Show endpoint details (status,words,lines,length)
t|tree              Show a treeview of endpoints (glitchy with multiple vhosts, no write to file)
w|write [filename]  write last output to file
sf|show-filters     Show current filters
cf|clear-filters    Clear all filters
fc [val,val2]       Filter Status code
fw [val,val2]       Filter by words
fl [val,val2]       Filter lines
fs [val,val2]       Filter lenght
fr [regex1,regex2]  Filter URL using regex
mc [val,val2]       Match status code
mw [val,val2]       Match words
ml [val,val2]       Match lines
ms [val,val2]       Match length
mr [regex1,regex2]  Match URL using regex
s|sort              Sort options
g|grep [expr]       Run grep on last output
gv|grepv [expr]     Run grep exclude on last output
r|reload            Reparse input path for new files
j|join              Combine all filtered results and export to ffuf JSON file.
```

### Command Line scripting
Parse all files in ./data, match with regex only php endpoints, sort by status code, and write full urls to out.txt
```
./pffuf -cl 'mr .*?\.php; s status; urls; w ./out.txt' ./data
```

