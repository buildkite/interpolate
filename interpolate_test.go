package interpolate_test

import (
	"fmt"
	"testing"

	"github.com/buildkite/interpolate"
)

func ExampleInterpolate() {
	env := interpolate.EnvFromSlice([]string{
		"HELLO_WORLD=🦀",
	})

	output, _ := interpolate.Interpolate(env, "Buildkite... ${HELLO_WORLD} ${ANOTHER_VAR:-🏖}")
	fmt.Println(output)

	// Output: Buildkite... 🦀 🏖
}

func TestBasicInterpolation(t *testing.T) {
	environ := map[string]string{
		"TEST1": "A test",
		"TEST2": "Another",
		"TEST3": "Llamas",
		"TEST4": "Only one level of $TEST3 interpolation",
	}

	for _, tc := range []struct {
		Str      string
		Expected string
	}{
		{``, ``},
		{`foo`, `foo`},
		{`test1`, `test1`},
		{`TEST1`, `TEST1`},
		{`$TEST1`, `A test`},
		{`${TEST1}`, `A test`},
		{`$TEST1, $TEST2, $TEST3`, `A test, Another, Llamas`},
		{`$Test1, $Test2, $TeST3`, `, , `},
		{`${TEST1}, ${Test2}, ${tEST3}`, `A test, , `},
		{`my$TEST1`, `myA test`},
		{`$TEST4`, "Only one level of $TEST3 interpolation"},

		// currently failing
		//{`${TEST4}`, "Only one level of $TEST3 interpolation"},
	} {
		result, err := interpolate.Interpolate(environ, tc.Str)
		if err != nil {
			t.Fatal(err)
		}
		if result != tc.Expected {
			t.Fatalf("Test %q failed: Expected substring %q, got %q", tc.Str, tc.Expected, result)
		}
	}
}

func TestVariablesMustStartWithLetters(t *testing.T) {
	for _, str := range []string{
		`$1 burgers`,
		`$99bottles`,
	} {
		_, err := interpolate.Interpolate(nil, str)
		if err == nil {
			t.Fatalf("Test %q should have resulted in an error", str)
		}
	}
}

func TestMissingParameterValuesReturnEmptyStrings(t *testing.T) {
	for _, str := range []string{
		`$BUILDKITE_COMMIT`,
		`${BUILDKITE_COMMIT}`,
		`${BUILDKITE_COMMIT:0:7}`,
		`${BUILDKITE_COMMIT:7}`,
		`${BUILDKITE_COMMIT:0:7}`,
		`${BUILDKITE_COMMIT:7:14}`,
	} {
		result, err := interpolate.Interpolate(nil, str)
		if err != nil {
			t.Fatal(err)
		}
		if result != "" {
			t.Fatalf("Expected empty string, got %q", result)
		}
	}
}

func TestSubstringsWithOffsets(t *testing.T) {
	environ := map[string]string{"BUILDKITE_COMMIT": "1adf998e39f647b4b25842f107c6ed9d30a3a7c7"}

	for _, tc := range []struct {
		Str      string
		Expected string
	}{
		// in range offsets, no lengths
		{`${BUILDKITE_COMMIT:0}`, `1adf998e39f647b4b25842f107c6ed9d30a3a7c7`},
		{`${BUILDKITE_COMMIT:7}`, `e39f647b4b25842f107c6ed9d30a3a7c7`},
		{`${BUILDKITE_COMMIT:-7}`, `0a3a7c7`},

		// out of range offsets, no lengths
		{`${BUILDKITE_COMMIT:-128}`, `1adf998e39f647b4b25842f107c6ed9d30a3a7c7`},
		{`${BUILDKITE_COMMIT:128}`, ``},

		// in range offsets and lengths
		{`${BUILDKITE_COMMIT:0:7}`, `1adf998`},
		{`${BUILDKITE_COMMIT:7:7}`, `e39f647`},
		{`${BUILDKITE_COMMIT:7:-7}`, `e39f647b4b25842f107c6ed9d3`},

		// zero lengths
		{`${BUILDKITE_COMMIT:0:0}`, ``},
		{`${BUILDKITE_COMMIT:7:0}`, ``},

		// in range offsets and out of range lengths
		{`${BUILDKITE_COMMIT:0:128}`, `1adf998e39f647b4b25842f107c6ed9d30a3a7c7`},
		{`${BUILDKITE_COMMIT:7:128}`, `e39f647b4b25842f107c6ed9d30a3a7c7`},
		{`${BUILDKITE_COMMIT:0:-128}`, ``},
		{`${BUILDKITE_COMMIT:7:-128}`, ``},
	} {
		result, err := interpolate.Interpolate(environ, tc.Str)
		if err != nil {
			t.Fatal(err)
		}
		if result != tc.Expected {
			t.Fatalf("Expected substring %q, got %q", tc.Expected, result)
		}
	}
}

func TestInterpolateIsntGreedy(t *testing.T) {
	environ := map[string]string{
		"BUILDKITE_COMMIT":       "cfeeee3fa7fa1a6311723f5cbff95b738ec6e683",
		"BUILDKITE_PARALLEL_JOB": "456",
	}

	for _, tc := range []struct {
		Str      string
		Expected string
	}{
		{`echo "ENV_1=test_$BUILDKITE_COMMIT_$BUILDKITE_PARALLEL_JOB"`, `echo "ENV_1=test_456"`},
		{`echo "ENV_1=test-$BUILDKITE_COMMIT-$BUILDKITE_PARALLEL_JOB"`, `echo "ENV_1=test_cfeeee3fa7fa1a6311723f5cbff95b738ec6e683-456"`},
		{`echo "ENV_1=test_${BUILDKITE_COMMIT}_${BUILDKITE_PARALLEL_JOB}"`, `echo "ENV_2=test_cfeeee3fa7fa1a6311723f5cbff95b738ec6e683_456"`},
	} {
		result, err := interpolate.Interpolate(environ, tc.Str)
		if err != nil {
			t.Fatal(err)
		}
		if result != tc.Expected {
			t.Fatalf("Expected substring %q, got %q", tc.Expected, result)
		}
	}
}

func TestDefaultValues(t *testing.T) {
	environ := map[string]string{
		"DAY":       "Blarghday",
		"EMPTY_DAY": "",
	}

	for _, tc := range []struct {
		Str      string
		Expected string
	}{
		{`Today is ${TODAY-Tuesday}`, `Today is Tuesday`},
		{`Tomorrow is ${TOMORROW-Wednesday}`, `Tomorrow is Wednesday`},
		{`Today is ${DAY-Wednesday}`, `Today is Blarghday`},
		{`Today is ${EMPTY_DAY-Wednesday}`, `Today is `},
		{`Today is ${EMPTY_DAY:-Wednesday}`, `Today is Wednesday`},
		{`${EMPTY_DAY:--:{}}`, `-:{}`},
	} {
		result, err := interpolate.Interpolate(environ, tc.Str)
		if err != nil {
			t.Fatal(err)
		}
		if result != tc.Expected {
			t.Fatalf("Test %q failed: Expected substring %q, got %q", tc.Str, tc.Expected, result)
		}
	}
}

func TestRequiredVariables(t *testing.T) {
	for _, tc := range []struct {
		Str         string
		ExpectedErr string
	}{
		{`Hello ${REQUIRED_VAR?}`, `$REQUIRED_VAR: not set`},
		{`Hello ${REQUIRED_VAR?y u no set me? :-{}`, `$REQUIRED_VAR: y u no set me? :-{`},
		{`Hello ${REQUIRED_VAR?{}}`, `$REQUIRED_VAR: {`},
	} {
		_, err := interpolate.Interpolate(nil, tc.Str)
		if err == nil || err.Error() != tc.ExpectedErr {
			t.Fatalf("Test %q should have failed with error %q, got %v", tc.Str, tc.ExpectedErr, err)
		}
	}
}

func TestEscapingVariables(t *testing.T) {
	for _, tc := range []struct {
		Str      string
		Expected string
	}{
		{`Do this $$ESCAPE_PARTY`, `Do this $ESCAPE_PARTY`},
		{`Do this \$ESCAPE_PARTY`, `Do this $ESCAPE_PARTY`},
		{`Do this $${SUCH_ESCAPE}`, `Do this ${SUCH_ESCAPE}`},
		{`Do this \${SUCH_ESCAPE}`, `Do this ${SUCH_ESCAPE}`},
	} {
		result, err := interpolate.Interpolate(nil, tc.Str)
		if err != nil {
			t.Fatal(err)
		}
		if result != tc.Expected {
			t.Fatalf("Test %q failed: Expected substring %q, got %q", tc.Str, tc.Expected, result)
		}
	}
}
