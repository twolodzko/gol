package parser

// func TestParse(t *testing.T) {
// 	var testCases = []struct {
// 		input    string
// 		expected list.List
// 	}{
// 		{"", list.List{}},
// 		{"a", list.NewList("a")},
// 	}

// 	for _, tt := range testCases {
// 		result, err := Parse(tt.input)

// 		if err != nil {
// 			t.Errorf("unexpected error: %s", err)
// 		}

// 		if same, _ := list.AreSame(result, tt.expected); !same {
// 			t.Errorf("expected: %v, got: %s", tt.expected, result)
// 		}
// 	}
// }
