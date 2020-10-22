#!/bin/bash

# Copyright 2019 Google LLC
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     https://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.

##
# system_tests.sh
# Runs CI checks for entire repository.
#
# Parameters
#
# [ARG 1]: Directory for the samples. Default: github/golang-samples.
# KOKORO_GFILE_DIR: Persistent filesystem location. (environment variable)
# KOKORO_KEYSTORE_DIR: Secret storage location. (environment variable)
# GOLANG_SAMPLES_GO_VET: If set, run code analysis checks. (environment variable)
##

set -ex

go version
date

export GO111MODULE=on # Always use modules.
export GOPROXY=https://proxy.golang.org
TIMEOUT=60m

# Don't print environment variables in case there are secrets.
# If you need a secret, use a keystore_resource in common.cfg.
set +x

export GOLANG_SAMPLES_KMS_KEYRING=ring1
export GOLANG_SAMPLES_KMS_CRYPTOKEY=key1

export GCLOUD_ORGANIZATION=1081635000895
export SCC_PUBSUB_PROJECT="project-a-id"
export SCC_PUBSUB_TOPIC="projects/project-a-id/topics/notifications-sample-topic"
export SCC_PUBSUB_SUBSCRIPTION="notification-sample-subscription"

export GOLANG_SAMPLES_SPANNER=projects/golang-samples-tests/instances/golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_PROJECT=golang-samples-tests
export GOLANG_SAMPLES_BIGTABLE_INSTANCE=testing-instance

export GOLANG_SAMPLES_FIRESTORE_PROJECT=golang-samples-fire-0

echo $CREDENTIALS > ~/secret.json
export GOOGLE_APPLICATION_CREDENTIALS=~/secret.json

set -x

pwd

mkdir -p ~/.secureConnect/
cp ./testing/dca/context_aware_metadata.json ~/.secureConnect/context_aware_metadata.json
cp ./testing/dca/cert ~/cert

# exit_code collects all of the exit codes of the tests, and is used to set the
# exit code at the end of the script.
exit_code=0
set +e # Don't exit on errors to make sure we run all tests.

# runSamples runs the tests in the current directory. If an argument is specified,
# it is used as the argument to `go test`.
runSamples() {
  set +x
  echo "Running 'go test' in '$(pwd)'..."
  set -x
  2>&1 go test -timeout $TIMEOUT -v "${1:-./...}"
  exit_code=$((exit_code + $?))
  set +x
}

# Returns 0 if the test should be skipped because the current Go
# version is too old for the current module.
goVersionShouldSkip() {
  modVersion="$(go list -m -f '{{.GoVersion}}')"
  if [ -z "$modVersion" ]; then
    # Not in a module or minimum Go version not specified, don't skip.
    return 1
  fi

  go list -f "{{context.ReleaseTags}}" | grep -q -v "go$modVersion\b"
}

if [[ $RUN_ALL_TESTS = "1" ]]; then
  echo "Running all tests"
  # shellcheck disable=SC2044
  for i in $(find . -name go.mod); do
    pushd "$(dirname "$i")" > /dev/null;
      runSamples
    popd > /dev/null;
  done
else
  echo "Running tests in the following directory: $TESTING_DIR"
  for d in $TESTING_DIR; do
    mods=$(find "$d" -name go.mod)
    # If there are no modules, just run the tests directly.
    if [[ -z "$mods" ]]; then
      pushd "$d" > /dev/null;
        runSamples
      popd > /dev/null;
    # Otherwise, run the tests in all Go directories. This way, we don't have to
    # check to see if there are tests that aren't in a sub-module.
    else
      goDirectories="$(find "$d" -name "*.go" -printf "%h\n" | sort -u)"
      if [[ -n "$goDirectories" ]]; then
        for gd in $goDirectories; do
          pushd "$gd" > /dev/null;
            runSamples .
          popd > /dev/null;
        done
      fi
    fi
  done
fi

exit $exit_code
