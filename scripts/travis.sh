#!/bin/bash
# Script that Travis CI will run to test the project, it is conditional on the branch.

set -e

# Check that the VERSION file has been bumped on release branches
if [[ $TRAVIS_BRANCH == release/* ]]
then
  make check-version
fi

DIFF_CHANGES=$(go run ./ci/diff_changes.go)
echo "Result of diff_changes.go: $DIFF_CHANGES"

if [ "$DIFF_CHANGES" = "both" ] || [ "$DIFF_CHANGES" = "go" ]
then 
  # Build and test the CLI
  echo "Testing Go CLI"
  cd ./cli
  dep ensure -v
  make build
  make test
  cd ..
fi

if [ "$DIFF_CHANGES" = "both" ] || [ "$DIFF_CHANGES" = "other" ]
then
  # Test the Docker Compose stack
  echo "Testing Docker Compose stack"
  make integration-tests
fi
