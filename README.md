Interpolate
===========

[![GoDoc](https://godoc.org/github.com/buildkite/interpolate?status.svg)](https://godoc.org/github.com/buildkite/interpolate)

A golang library for parameter expansion (like `${BLAH}` or `$BLAH`) in strings from environment variables. An implementation of [POSIX Parameter Expansion](http://pubs.opengroup.org/onlinepubs/9699919799/utilities/V3_chap02.html#tag_18_06_02), plus some other basic operations that you'd expect in a shell scripting environment like bash.

## Installation

```
go get -u github.com/buildkite/interpolate
```

## Usage

```go
package main

import (
  "github.com/buildkite/interpolate"
  "fmt"
)

func main() {
	env := interpolate.NewSliceEnv([]string{
		"HELLO_WORLD=🦀",
	})

	output, _ := interpolate.Interpolate(env, "Buildkite... ${HELLO_WORLD} ${ANOTHER_VAR:-🏖}")
	fmt.Println(output)
}

// Output: Buildkite... 🦀 🏖

```

## Supported Expansions

#### *`${parameter}`* or *`$parameter`*
**Use value.** If parameter is set, then it shall be substituted; otherwise it will be blank

#### *`${parameter:-[word]}`*
**Use default values.** If parameter is unset or null, the expansion of word (or an empty string if word is omitted) shall be substituted; otherwise, the value of parameter shall be substituted.

#### *`${parameter-[word]}`*
**Use default values when not set.** If parameter is unset, the expansion of word (or an empty string if word is omitted) shall be substituted; otherwise, the value of parameter shall be substituted.

#### *`${parameter:[offset]}`*
**Use the substring of parameter after offset.** A negative number will select from the end of the string. If the value is out of bounds, an empty string will be substituted.

#### *`${parameter:[offset]:[length]}`*
**Use the substring of parameter after offset of given length.** A negative number will select from the end of the string. If the offset is out of bounds, an empty string will be substituted. If the length is greater than the length then the entire string will be returned.

#### *`${parameter:?[word]}`*
**Indicate Error if Null or Unset.** If parameter is unset or null, the expansion of word (or a message indicating it is unset if word is omitted) shall be returned as an error.

#### *`${parameter:%[word]}`*
**Remove Smallest Suffix Pattern.** The word shall be expanded to produce a pattern. The parameter expansion shall then result in parameter, with the smallest portion of the suffix matched by the pattern deleted. If present, word shall not begin with an unquoted '%'.

#### *`${parameter:%%[word]}`*
**Remove Largest Suffix Pattern.** The word shall be expanded to produce a pattern. The parameter expansion shall then result in parameter, with the largest portion of the suffix matched by the pattern deleted.

#### *`${parameter:#[word]}`*
**Remove Smallest Prefix Pattern.** The word shall be expanded to produce a pattern. The parameter expansion shall then result in parameter, with the smallest portion of the prefix matched by the pattern deleted. If present, word shall not begin with an unquoted '#'.

#### *`${parameter:##[word]}`*
**Remove Largest Prefix Pattern.** The word shall be expanded to produce a pattern. The parameter expansion shall then result in parameter, with the largest portion of the prefix matched by the pattern deleted.

## License

Licensed under MIT license, in `LICENSE`.
