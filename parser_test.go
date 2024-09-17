package interpolate

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestParser(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		input string
		want  Expression
	}{
		{
			input: "Buildkite... ${HELLO_WORLD} ${ANOTHER_VAR:-🏖}",
			want: Expression{
				{Text: "Buildkite... "},
				{Expansion: VariableExpansion{
					Identifier: "HELLO_WORLD",
				}},
				{Text: " "},
				{Expansion: EmptyValueExpansion{
					Identifier: "ANOTHER_VAR",
					Content: Expression{{
						Text: "🏖",
					}},
				}},
			},
		},
		{
			input: "${TEST1:- ${TEST2:-$TEST3}}",
			want: Expression{
				{Expansion: EmptyValueExpansion{
					Identifier: "TEST1",
					Content: Expression{
						{Text: " "},
						{Expansion: EmptyValueExpansion{
							Identifier: "TEST2",
							Content: Expression{
								{Expansion: VariableExpansion{
									Identifier: "TEST3",
								}},
							},
						}},
					},
				}},
			},
		},
		{
			input: "${HELLO_WORLD-blah}",
			want: Expression{
				{Expansion: UnsetValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: Expression{{
						Text: "blah",
					}},
				}},
			},
		},
		{
			input: `\\${HELLO_WORLD-blah}`,
			want: Expression{
				{Text: `\\`},
				{Expansion: UnsetValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: Expression{{
						Text: "blah",
					}},
				}},
			},
		},
		{
			input: `\${HELLO_WORLD-blah}`,
			want: Expression{
				{Expansion: EscapedExpansion{
					PotentialIdentifier: "{HELLO_WORLD-blah}",
				}},
				{Text: "{HELLO_WORLD-blah}"},
			},
		},
		{
			input: `Test \\\${HELLO_WORLD-blah}`,
			want: Expression{
				{Text: "Test "},
				{Text: `\\`},
				{Expansion: EscapedExpansion{
					PotentialIdentifier: "{HELLO_WORLD-blah}",
				}},
				{Text: "{HELLO_WORLD-blah}"},
			},
		},
		{
			input: `${HELLO_WORLD:1}`,
			want: Expression{
				{Expansion: SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
				}},
			},
		},
		{
			input: "${HELLO_WORLD: -1}",
			want: Expression{
				{Expansion: SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     -1,
				}},
			},
		},
		{
			input: `${HELLO_WORLD:-1}`,
			want: Expression{
				{Expansion: EmptyValueExpansion{
					Identifier: "HELLO_WORLD",
					Content: Expression{{
						Text: "1",
					}},
				}},
			},
		},
		{
			input: `${HELLO_WORLD:1:7}`,
			want: Expression{
				{Expansion: SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
					Length:     7,
					HasLength:  true,
				}},
			},
		},
		{
			input: `${HELLO_WORLD:1:-7}`,
			want: Expression{
				{Expansion: SubstringExpansion{
					Identifier: "HELLO_WORLD",
					Offset:     1,
					Length:     -7,
					HasLength:  true,
				}},
			},
		},
		{
			input: `${HELLO_WORLD?Required}`,
			want: Expression{
				{Expansion: RequiredExpansion{
					Identifier: "HELLO_WORLD",
					Message: Expression{
						{Text: "Required"},
					},
				}},
			},
		},
		{
			input: `$${not actually a brace expression`,
			want: Expression{
				{Expansion: EscapedExpansion{}},
				{Text: "{not actually a brace expression"},
			},
		},
		{
			input: "$",
			want: Expression{
				{Text: "$"},
			},
		},
		{
			input: `\`,
			want: Expression{
				{Text: `\`},
			},
		},
		{
			input: "$(echo hello world)",
			want: Expression{
				{Text: "$("},
				{Text: "echo hello world)"},
			},
		},
		{
			input: "$$MOUNTAIN",
			want: Expression{
				{Expansion: EscapedExpansion{PotentialIdentifier: "MOUNTAIN"}},
				{Text: "MOUNTAIN"},
			},
		},
		{
			input: "this is a regex! /^start.*end$$/", // the dollar sign at the end of the regex has to be escaped to be treated as a literal dollar sign by this library
			want: Expression{
				{Text: "this is a regex! /^start.*end"},
				{Expansion: EscapedExpansion{}},
				{Text: "/"},
			},
		},
		{
			input: `this is a more different regex! /^start.*end\$/`, // the dollar sign at the end of the regex has to be escaped to be treated as a literal dollar sign by this library
			want: Expression{
				{Text: "this is a more different regex! /^start.*end"},
				{Expansion: EscapedExpansion{}},
				{Text: "/"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.input, func(t *testing.T) {
			t.Parallel()

			got, err := NewParser(tc.input).Parse()
			if err != nil {
				t.Fatalf("NewParser(%q).Parse() error = %v", tc.input, err)
			}

			if diff := cmp.Diff(got, tc.want); diff != "" {
				t.Errorf("parsed expression diff (-got +want):\n%s", diff)
			}
		})
	}
}
