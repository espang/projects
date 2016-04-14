package boolParser

import "testing"

type test struct {
	eval     string
	literals map[string]int
	testmap  map[int]bool
}

var tests = []test{
	{
		eval:     "A and B",
		literals: map[string]int{"A": 1, "B": 2},
		testmap: map[int]bool{
			0: false,
			1: false,
			2: false,
			3: true,
		},
	},
	{
		eval:     "A or B",
		literals: map[string]int{"A": 1, "B": 2},
		testmap: map[int]bool{
			0: false,
			1: true,
			2: true,
			3: true,
		},
	},
	{
		eval:     "A and B and C",
		literals: map[string]int{"A": 1, "B": 2, "C": 4},
		testmap: map[int]bool{
			0: false,
			1: false,
			2: false,
			3: false,
			4: false,
			5: false,
			6: false,
			7: true,
		},
	},
	{
		eval:     "A or B and C",
		literals: map[string]int{"A": 1, "B": 2, "C": 4},
		testmap: map[int]bool{
			0: false,
			1: true,
			2: false,
			3: true,
			4: false,
			5: true,
			6: true,
			7: true,
		},
	},
	{
		eval:     "( A or B ) and C",
		literals: map[string]int{"A": 1, "B": 2, "C": 4},
		testmap: map[int]bool{
			0: false,
			1: false,
			2: false,
			3: false,
			4: false,
			5: true,
			6: true,
			7: true,
		},
	},
	{
		eval:     "not A",
		literals: map[string]int{"A": 1},
		testmap: map[int]bool{
			0: true,
			1: false,
		},
	},
	{
		eval:     "A or not A",
		literals: map[string]int{"A": 1},
		testmap: map[int]bool{
			0: true,
			1: true,
		},
	},
	{
		eval:     "not ( A and B ) or not C",
		literals: map[string]int{"A": 1, "B": 2, "C": 4},
		testmap: map[int]bool{
			0: true,
			1: true,
			2: true,
			3: true,
			4: true,
			5: true,
			6: true,
			7: false,
		},
	},
}

func TestAll(t *testing.T) {
	for _, testcase := range tests {
		expr, err := NewExpression(testcase.eval)
		if err != nil {
			t.Errorf("error creating expression: %v", err)
		}

		for value, expect := range testcase.testmap {
			vMap := map[string]bool{}
			for lit, bit := range testcase.literals {
				vMap[lit] = value&bit == bit
			}

			result := expr.evaluate(vMap)
			if result != expect {
				t.Errorf(
					"error evaluating '%s' got %v, expect %v. Values: %v",
					testcase.eval,
					result,
					expect,
					vMap,
				)
			}
		}

	}
}
