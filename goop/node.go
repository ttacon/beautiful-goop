package goop

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"io"
	"regexp"
	"strings"
)

type GoopNode struct {
	*html.Node
}

type Goop struct {
	Root *GoopNode
}

func NewGoopNode(n *html.Node) *GoopNode {
	return &GoopNode{n}
}

func BuildGoop(r io.Reader) (*Goop, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Goop{NewGoopNode(doc)}, nil
}

func (g *Goop) Find(query string) []*GoopNode {
	return g.Root.Find(query)
}

var ele = regexp.MustCompile(`^([a-zA-Z\d]*)[\.|#]+.*$`)
var classes = regexp.MustCompile(`(\.[a-zA-Z0-9_-]*)*`)
var id = regexp.MustCompile(`(#[\w]*)`)
var reg = regexp.MustCompile(`([a-zA-Z\d]*([\.|#][\da-zA-Z_-]*])*)*`)
var nodeReg = regexp.MustCompile(`([\s]*)(#|\.[\w-]*)*`)

func tokenize(query string) [][]string {
	vals := make([][]string, 3)

	firstClass := strings.Index(query, ".")
	id := strings.Index(query, "#")
	endElement := firstClass
	if id < firstClass {
		endElement = id
	}

	// get html element if exists
	if id == 0 || firstClass == 0 {
		vals[0] = []string{}
	} else {
		vals[0] = []string{query[0:endElement]}
		query = query[endElement:]
	}

	// get id if any
	if id > -1 {
		id = strings.Index(query, "#")
		part := query[id:]
		endId := strings.Index(part, ".")
		if endId == -1 {
			vals[1] = []string{query[id+1:]}
			query = query[0:id]
		} else {
			vals[1] = []string{query[id+1 : id+endId]}
			query = strings.Join([]string{query[0:id], query[id+endId:]}, "")
		}
	} else {
		vals[1] = []string{}
	}

	// get all classes
	if firstClass > -1 {
		vals[2] = strings.Split(strings.TrimSpace(query), ".")[1:]
	} else {
		vals[2] = []string{}
	}
	return vals
}

func (g *GoopNode) Find(query string) []*GoopNode {
	// parse query for element, classes and id
	queries := strings.Split(query, ",")
	searchFrom := []*GoopNode{g}
	for _, q := range queries {
		for _, p := range strings.Split(q, " ") {
			vals := tokenize(p)
			// var found []*GoopNode
			for _, s := range searchFrom {
				// if we have an id, search by it
				if len(vals[1]) > 0 {
					s.FindById(vals[1][0])
				}
				// TODO(stopped here)
			}
		}
	}

	return nil
}

func (g *Goop) FindAllElements(ele string) []*GoopNode {
	ele = strings.Title(ele)
	eleAtom := atom.Lookup([]byte(ele))
	return g.Root.SearchByElement(eleAtom)
}

func (g *GoopNode) SearchByElement(a atom.Atom) []*GoopNode {
	var found []*GoopNode
	if g.DataAtom == a {
		found = append(found, g)
	}

	for child := g.FirstChild; child != nil; child = child.NextSibling {
		if gns := (&GoopNode{child}).SearchByElement(a); len(gns) > 0 {
			found = append(found, gns...)
		}
	}
	return found
}

func (g *Goop) FindAllWithClass(class string) []*GoopNode {
	return g.Root.SearchByClass(class)
}

func (g *GoopNode) SearchByClass(class string) []*GoopNode {
	var found []*GoopNode
	for _, attr := range g.Attr {
		if attr.Key == "class" && attr.Val == class {
			found = append(found, g)
		}
	}

	for child := g.FirstChild; child != nil; child = child.NextSibling {
		if gns := (&GoopNode{child}).SearchByClass(class); len(gns) > 0 {
			found = append(found, gns...)
		}
	}

	return found
}

func (g *Goop) FindById(id string) *GoopNode {
	return g.Root.FindById(id)
}

func (g *GoopNode) FindById(id string) *GoopNode {
	for _, attr := range g.Attr {
		if attr.Key == "id" && attr.Val == id {
			return g
		}
	}

	for child := g.FirstChild; child != nil; child = child.NextSibling {
		if n := (&GoopNode{child}).FindById(id); n != nil {
			return n
		}
	}

	return nil
}

func (g *GoopNode) Attributes() map[string][]string {
	attrs := make(map[string][]string)
	for _, attr := range g.Attr {
		if vals, ok := attrs[attr.Key]; ok {
			attrs[attr.Key] = append(vals, attr.Val)
		} else {
			attrs[attr.Key] = []string{attr.Val}
		}
	}
	return attrs
}
