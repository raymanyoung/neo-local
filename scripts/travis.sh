#!/bin/bash
# Script that Travis CI will run to test the project, it is conditional on the branch.

version=$(cat ./VERSION)

set -e

# Check that the VERSION file has been bumped on release branches
if [[ $TRAVIS_BRANCH == release/* ]]
then
  make check-version
fi

diff_changes=$(go run ./ci/diff_changes.go)
echo $diff_changes

if [ "$diff_changes" = "both" ] || [ "$diff_changes" = "go" ]
then 
  # Build and test the CLI
  cd ./cli
  dep ensure -v
  make build
  make test
  cd ..
fi

if [ "$diff_changes" = "both" ] || [ "$diff_changes" = "other" ]
then 
  # Test the Docker Compose stack
  make integration-tests
fi
