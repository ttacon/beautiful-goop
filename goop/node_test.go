package goop

import (
	"code.google.com/p/go.net/html"
	"strings"
	"testing"
)

func Test_FindAllElement(t *testing.T) {

}

func Test_NewGoopNode(t *testing.T) {

}

type goopNodeTest struct {
	input string
	node  *GoopNode
}

type goopTest struct {
	input string
	goop  *Goop
}

func Test_BuildGoop(t *testing.T) {
	parent := &html.Node{
		Type: 0x2,
	}
	child := &html.Node{
		Type:     0x3,
		DataAtom: 0x27604,
		Data:     "html",
	}
	head := &html.Node{
		Parent:   child,
		Type:     0x3,
		DataAtom: 0x2fa04,
		Data:     "head",
	}

	body := &html.Node{
		Parent:      child,
		PrevSibling: head,
		Type:        0x3,
		DataAtom:    0x2f04,
		Data:        "body",
	}
	head.NextSibling = body

	div := &html.Node{
		Parent:   body,
		Type:     0x3,
		DataAtom: 0x10703,
		Data:     "div",
	}
	body.FirstChild = div
	body.LastChild = div

	foo := &html.Node{
		Parent: div,
		Type:   0x1,
		Data:   "Foo",
	}
	div.FirstChild = foo
	div.LastChild = foo

	parent.FirstChild = child
	parent.LastChild = child
	child.FirstChild = head
	child.LastChild = body
	child.Parent = parent

	tests := []goopTest{
		goopTest{
			"<div>Foo</div>",
			&Goop{Root: &GoopNode{parent}},
		},
	}

	for _, test := range tests {
		g, err := BuildGoop(strings.NewReader(test.input))

		if err != nil {
			t.Errorf("error occured while building some tasty goop: %v", err)
			continue
		}

		if !nodeEqual(test.goop.Root.Node, g.Root.Node) {
			t.Errorf("goop built: %v doesnt match expected %v\n", g, test.goop)
		}
	}
}

func nodeEqual(n1, n2 *html.Node) bool {
	if n1 == nil || n2 == nil {
		return true
	}
	if (n1 != nil && n2 == nil) || (n1 == nil && n2 != nil) {
		return false
	}

	// TODO(ttacon): go through node's own siblings

	c1 := n1.FirstChild
	c2 := n2.FirstChild
	for c1 != nil && c2 != nil {
		if c1 == nil || c2 == nil {
			return false
		}
		if !nodeEqual(c1, c2) {
			return false
		}
		c1 = c1.NextSibling
		c2 = c2.NextSibling
	}

	return n1.Type == n2.Type &&
		n1.Data == n2.Data &&
		n1.DataAtom == n2.DataAtom
}

func Test_GoopFind(t *testing.T) {
}

type tokenizeTest struct {
	input  string
	output [][]string
}

func Test_tokenize(t *testing.T) {
	tests := []tokenizeTest{
		tokenizeTest{
			"div#id.class0.class1.class2",
			[][]string{
				[]string{
					"div",
				},
				[]string{
					"id",
				},
				[]string{
					"class0",
					"class1",
					"class2",
				},
			},
		},
		tokenizeTest{
			"#id.class0.class1.class2",
			[][]string{
				[]string{},
				[]string{
					"id",
				},
				[]string{
					"class0",
					"class1",
					"class2",
				},
			},
		},
		tokenizeTest{
			".class0.class1.class2",
			[][]string{
				[]string{},
				[]string{},
				[]string{
					"class0",
					"class1",
					"class2",
				},
			},
		},
		/*		tokenizeTest{
				"div#id#id2.class0.class1.class2",
				[][]string{
					[]string{
						"div",
					},
					[]string{
						"id",
					},
					[]string{
						"class0",
						"class1",
						"class2",
					},
				},
			},*/
		tokenizeTest{
			"a.class0",
			[][]string{
				[]string{
					"a",
				},
				[]string{},
				[]string{
					"class0",
				},
			},
		},
	}

	for _, test := range tests {
		vals := tokenize(test.input)
		if !sliceEquality(vals, test.output) {
			t.Errorf("tokenization failed, expected: %v, got: %v", test.output, vals)
		}
	}
}

func sliceEquality(s1 [][]string, s2 [][]string) bool {
	if len(s1) != len(s2) {
		return false
	}
	for i, val1 := range s1 {
		val2 := s2[i]
		if len(val1) != len(val2) {
			return false
		}
		for j, v1 := range val1 {
			if v1 != val2[j] {
				return false
			}
		}
	}
	return true
}

func Test_GoopNodeFind(t *testing.T) {

}

func Test_GoopNodeHasClasses(t *testing.T) {

}

func Test_GoopNodeIsElement(t *testing.T) {

}

func Test_GoopFindAllElements(t *testing.T) {

}

func Test_GoopNodeFindAllElements(t *testing.T) {

}

func Test_GoopNodeSearchByElement(t *testing.T) {

}

func Test_GoopFindAllWithClass(t *testing.T) {

}

func Test_GoopNodeSearchByClass(t *testing.T) {

}

func Test_GoopFindById(t *testing.T) {

}

func Test_GoopNodeFindById(t *testing.T) {

}

func Test_Attributes(t *testing.T) {

}
