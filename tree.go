package main

import (
	"fmt"
	"strings"
)

// Node tree struct
type Node struct {
	Name     string
	Children []Node
}

var out = ""

func parseNode(nde Node, level int) {
	for i, b := range nde.Children {
		if len(b.Name) > 0 {
			if i+1 == len(nde.Children) {
				fmt.Println(indentStr(level, true), b.Name)
				out = out + fmt.Sprintln(indentStr(level, true)+b.Name)
			} else {
				fmt.Println(indentStr(level, false), b.Name)
				out = out + fmt.Sprintln(indentStr(level, false)+b.Name)
			}
		}
		if len(b.Children) > 0 {

			parseNode(b, level+1)
		}
	}
}

func indentStr(level int, last bool) string {
	baseindent := "│     "
	indent := ""
	for i := 0; i < level-1; i++ {
		indent = indent + baseindent
	}
	if !last {
		indent = indent + "├──"
	} else {
		indent = indent + "└──"
	}
	return "  " + indent
}

func doTreePlain(treeresult []NavResults) {
	s := []string{}

	for _, res := range treeresult {
		line := res.URL
		line = strings.Replace(line, "http://", "", 1)
		line = strings.Replace(line, "https://", "", 1)
		s = append(s, line)

	}

	var tree []Node
	for i := range s {
		tree = addToTree(tree, strings.Split(s[i], "/"))
	}

	fmt.Println(tree[0].Name)
	out = out + fmt.Sprintln(tree[0].Name)
	for i, a := range tree[0].Children {
		if len(a.Name) > 0 {
			if i+1 == len(tree[0].Children) {
				fmt.Println("  └──", a.Name)

				out = out + fmt.Sprintln("  └── "+a.Name)
			} else {
				out = out + fmt.Sprintln("  ├── "+a.Name)
				fmt.Println("  ├──", a.Name)
			}

		}
		parseNode(a, 2)

	}

}

func addToTree(root []Node, names []string) []Node {
	if len(names) > 0 {
		var i int
		for i = 0; i < len(root); i++ {
			if root[i].Name == names[0] { //already in tree
				break
			}
		}
		if i == len(root) {
			root = append(root, Node{Name: names[0]})
		}
		root[i].Children = addToTree(root[i].Children, names[1:])
	}
	return root
}
