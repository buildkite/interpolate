#!/usr/bin/env bash

set -Eeufo pipefail

echo --- :go: Checking go mod tidiness
go mod tidy
if ! git diff --no-ext-diff --exit-code; then
  echo ^^^ +++
  echo "The go.mod or go.sum files are out of sync with the source code"
  echo "Please run \`go mod tidy\` locally, and commit the result."

  exit 1
fi

echo --- :go: Checking go formatting
gofmt -w .
if ! git diff --no-ext-diff --exit-code; then
  echo ^^^ +++
  echo "Files have not been formatted with gofmt."
  echo "Fix this by running \`go fmt ./...\` locally, and committing the result."

  exit 1
fi

echo +++ Everything is clean and tidy! 🎉
