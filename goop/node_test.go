package goop

import "testing"

func Test_FindAllElement(t *testing.T) {

}

func Test_NewGoopNode(t *testing.T) {

}

func Test_BuildGoop(t *testing.T) {

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
