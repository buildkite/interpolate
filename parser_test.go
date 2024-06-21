package interpolate_test

import (
	"reflect"
	"testing"

	"github.com/buildkite/interpolate"
)

func TestParser(t *testing.T) {
	t.Parallel()

	var testCases = []struct {
		String   string
		Expected []interpolate.ExpressionItem
	}{
		{
			String: `Buildkite... ${HELLO_WORLD} ${ANOTHER_VAR:-🏖}`,
			Expected: []interpolate.ExpressionItem{
				{Text: "Buildkite... "},
				{Expansion: interpolate.VariableExpansion{
					Identifier: "HELLO_WORLD",
				}},
				{Text: " "},
				{Expansion: interpolate.EmptyValueExpansion{
					Identifier: "ANOTHER_VAR",
					Content: interpolate.Expression([]interpolate.ExpressionItem{{
						Text: "🏖",
					}}),
				}},
			},
		},
		{
			String: `${TEST1:- ${TEST2:-$TEST3}}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.EmptyValueExpansion{
					Identifier: "TEST1",
					Content: interpolate.Expression([]interpolate.ExpressionItem{
						{Text: " "},
						{Expansion: interpolate.EmptyValueExpansion{
							Identifier: "TEST2",
							Content: interpolate.Expression([]interpolate.ExpressionItem{
								{Expansion: interpolate.VariableExpansion{
									Identifier: "TEST3",
								}},
							}),
						}},
					}),
				}},
			},
		},
		{
			String: `${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.UnsetValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: interpolate.Expression([]interpolate.ExpressionItem{{
						Text: "blah",
					}}),
				}},
			},
		},
		{
			String: `\\${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{Text: `\\`},
				{Expansion: interpolate.UnsetValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: interpolate.Expression([]interpolate.ExpressionItem{{
						Text: "blah",
					}}),
				}},
			},
		},
		{
			String: `\${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.EscapedExpansion{Identifier: "{HELLO_WORLD-blah}"}},
			},
		},
		{
			String: `Test \\\${HELLO_WORLD-blah}`,
			Expected: []interpolate.ExpressionItem{
				{Text: `Test `},
				{Text: `\\`},
				{Expansion: interpolate.EscapedExpansion{Identifier: "{HELLO_WORLD-blah}"}},
			},
		},
		{
			String: `${HELLO_WORLD:1}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
				}},
			},
		},
		{
			String: `${HELLO_WORLD: -1}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     -1,
				}},
			},
		},
		{
			String: `${HELLO_WORLD:-1}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.EmptyValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: interpolate.Expression([]interpolate.ExpressionItem{{
						Text: "1",
					}}),
				}},
			},
		},
		{
			String: `${HELLO_WORLD:1:7}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
					Length:     7,
					HasLength:  true,
				}},
			},
		},
		{
			String: `${HELLO_WORLD:1:-7}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
					Length:     -7,
					HasLength:  true,
				}},
			},
		},
		{
			String: `${HELLO_WORLD?Required}`,
			Expected: []interpolate.ExpressionItem{
				{Expansion: interpolate.RequiredExpansion{
					Identifier: "HELLO_WORLD",
					Message: interpolate.Expression([]interpolate.ExpressionItem{
						{Text: "Required"},
					}),
				}},
			},
		},
		{
			String: `$`,
			Expected: []interpolate.ExpressionItem{
				{Text: `$`},
			},
		},
		{
			String: `\`,
			Expected: []interpolate.ExpressionItem{
				{Text: `\`},
			},
		},
		{
			String: `$(echo hello world)`,
			Expected: []interpolate.ExpressionItem{
				{Text: `$(`},
				{Text: `echo hello world)`},
			},
		},
		{
			String:   "$$MOUNTAIN",
			Expected: []interpolate.ExpressionItem{{Expansion: interpolate.EscapedExpansion{Identifier: "MOUNTAIN"}}},
		},
		{
			String: "this is a regex! /^start.*end$$/", // the dollar sign at the end of the regex has to be escaped to be treated as a literal dollar sign by this library
			Expected: []interpolate.ExpressionItem{
				{Text: "this is a regex! /^start.*end"}, {Text: "$"}, {Text: "/"},
			},
		},
		{
			String: `this is a more different regex! /^start.*end\$/`, // the dollar sign at the end of the regex has to be escaped to be treated as a literal dollar sign by this library
			Expected: []interpolate.ExpressionItem{
				{Text: "this is a more different regex! /^start.*end"}, {Text: "$"}, {Text: "/"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.String, func(t *testing.T) {
			t.Parallel()

			actual, err := interpolate.NewParser(tc.String).Parse()
			if err != nil {
				t.Fatal(err)
			}

			expected := interpolate.Expression(tc.Expected)
			if !reflect.DeepEqual(expected, actual) {
				t.Fatalf("Expected vs Actual: \n%s\n\n%s", expected, actual)
			}
		})
	}
}
