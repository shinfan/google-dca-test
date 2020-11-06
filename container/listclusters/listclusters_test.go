// Copyright 2019 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package main

import (
	"context"
	"log"
	"testing"

	"github.com/shinfan/google-dca-test/internal/testutil"
	container "google.golang.org/api/container/v1"
)

func TestSample(t *testing.T) {
	tc := testutil.EndToEndTest(t)
	println("Start test")
	svc, err := container.NewService(context.Background())
	if err != nil {
		log.Fatalf("Could not initialize gke client: %v", err)
	}

	if err := listClusters(svc, tc.ProjectID, "us-central1-c"); err != nil {
		log.Fatal(err)
	}
}
