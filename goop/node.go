package goop

import (
	"code.google.com/p/go.net/html"
	"code.google.com/p/go.net/html/atom"
	"io"
	"strings"
)

// A GoopNode is a wrapper for html.Node to add extended functionality.
type GoopNode struct {
	*html.Node
}

// Goop is represents a parsed webpage.
type Goop struct {
	Root *GoopNode // The root node of the document
}

// NewGoopNode creates a new GoopNode from its core html.Node
func NewGoopNode(n *html.Node) *GoopNode {
	return &GoopNode{n}
}

// BuildGoop constructs a Goop struct from a webpage
func BuildGoop(r io.Reader) (*Goop, error) {
	doc, err := html.Parse(r)
	if err != nil {
		return nil, err
	}
	return &Goop{NewGoopNode(doc)}, nil
}

// Find takes one or more comma separated selectors (as in jQuery) and returns
// a slice of GoopNodes satisfying the query.
func (g *Goop) Find(query string) []*GoopNode {
	return g.Root.Find(query)
}

// tokenizes a query and returns a [][]string of the form:
// vals[0] = html element (len(vals[0]) <= 1)
// vals[1] = element id (len(vals[1]) <= 1)
// vals[2] = element classes (len(vals[2]) >= 0)
func tokenize(query string) [][]string {
	vals := make([][]string, 3)

	firstClass := strings.Index(query, ".")
	id := strings.Index(query, "#")
	endElement := firstClass
	if id < firstClass && id != -1 {
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

// Find takes one or more comma separated selectors (as in jQuery) and returns
// a slice of GoopNodes satisfying the query.
func (g *GoopNode) Find(query string) []*GoopNode {
	// parse query for element, classes and id
	queries := strings.Split(query, ",")

	var toReturn []*GoopNode
	for _, q := range queries {
		searchFrom := []*GoopNode{g}

		for _, p := range strings.Split(q, " ") {
			vals := tokenize(p)
			var found []*GoopNode

			for _, s := range searchFrom {
				// if we have an id, search by it
				if len(vals[1]) > 0 {
					fs := s.FindById(vals[1][0])
					if fs == nil {
						continue
					}

					// now need to validate it has classes and type
					if g.HasClasses(vals[2]) && g.IsElement(vals[0]) {
						found = append(found, fs)
					}
					continue
				}

				// get all elements of specific type if it exists
				if len(vals[0]) > 0 {
					// TODO(ttacon): deal w/ case when len(vals[0]) > 0
					fs := s.FindAllElements(vals[0][0])

					for _, f := range fs {
						if f.HasClasses(vals[2]) {
							found = append(found, f)
						}
					}
					continue
				}

				// just get all elements by first class then verify
				if len(vals[2]) == 0 {
					continue
				}

				firstClass := vals[2][0]
				fs := s.SearchByClass(firstClass)
				if len(fs) == 0 {
					continue
				}

				if len(vals[2]) > 1 {
					found = append(found, fs...)
					continue
				}

				for _, f := range fs {
					if f.HasClasses(vals[2][1:]) {
						found = append(found, f)
					}
				}
			}
			searchFrom = found
		}
		if len(searchFrom) > 0 {
			toReturn = append(toReturn, searchFrom...)
		}
	}

	return toReturn
}

// Determine whether or not the receiving node has the given classes.
func (g *GoopNode) HasClasses(classes []string) bool {
	classMap := make(map[string]bool)
	for _, attr := range g.Attr {
		if attr.Key == "class" {
			classMap[attr.Val] = true
		}
	}

	for _, class := range classes {
		if _, ok := classMap[class]; !ok {
			return false
		}
	}
	return true
}

// Determine whether or not the receiving node is of the given element type.
func (g *GoopNode) IsElement(eles []string) bool {
	if len(eles) == 0 {
		return true
	}
	// TODO(ttacon): deal w/ > 1 ele
	ele := strings.Title(eles[0])
	eleAtom := atom.Lookup([]byte(ele))
	return g.DataAtom == eleAtom
}

// Returns all elements of the given type in the webpage.
func (g *Goop) FindAllElements(ele string) []*GoopNode {
	eleAtom := atom.Lookup([]byte(ele))
	return g.Root.SearchByElement(eleAtom)
}

// Finds all elements of the given type which are children of the receiving node
// (the recieving node may also be returned if it is of the given type).
func (g *GoopNode) FindAllElements(ele string) []*GoopNode {
	eleAtom := atom.Lookup([]byte(ele))
	return g.SearchByElement(eleAtom)
}

// SearchByElement has the same functionality of FindAllElements but takes in
// a html.Atom instead of the element type as a string
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

// Finds all elements with the given class in the current document.
func (g *Goop) FindAllWithClass(class string) []*GoopNode {
	return g.Root.SearchByClass(class)
}

// Finds all elements with the given class which are descended from the
// current node (the receiving node may also be returned if it has the
// given class).
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

// Finds the element with the given id in the current doc. FindById
// expects there to only be one element with the id - if the page you are
// searching in has more than one element with the same id, FindById will
// return the first one it encounters.
func (g *Goop) FindById(id string) *GoopNode {
	return g.Root.FindById(id)
}

// Finds an element by id from the current node (with the same constraints
// as when calling FindById from a Goop struct).
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

// Returns the attributes of a node in a friendlier format than stored in
// an html.Node.
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
