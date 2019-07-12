package main

import (
	"bufio"
	"fmt"
	"log"
	"os"

	"github.com/buildkite/interpolate"
)

func main() {
	stdin := bufio.NewScanner(os.Stdin)
	stdout := bufio.NewWriter(os.Stdout)

	for stdin.Scan() {
		env := interpolate.NewSliceEnv(os.Environ())
		line, err := interpolate.Interpolate(env, stdin.Text())
		if err != nil {
			log.Fatalf("Error while interpolating: %v", err)
		}
		_, err = fmt.Fprintln(stdout, line)
		if err != nil {
			log.Fatalf("Error while writing to stdout: %v", err)
		}
		stdout.Flush()
	}
}
