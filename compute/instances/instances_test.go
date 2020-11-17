// Copyright 2020 Google LLC
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

// Package instances is a set of utilities for managing Google Cloud VM instances.
package instances

import (
	"fmt"
	"testing"
	"time"

	"github.com/shinfan/google-dca-test/internal/testutil"
)

func TestAggregatedList(t *testing.T) {
	tc := testutil.SystemTest(t)
	_, err := aggregatedList(tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to get aggregated list: %v", err)
	}
}

func TestList(t *testing.T) {
	tc := testutil.SystemTest(t)
	_, err := list(tc.ProjectID, "us-central1-a")
	if err != nil {
		t.Fatalf("failed to get instance list: %v", err)
	}
}

func TestGet(t *testing.T) {
	tc := testutil.SystemTest(t)
	_, err := get(tc.ProjectID, "us-central1-a", "dca-test-instance-1")
	if err != nil {
		t.Fatalf("failed to get instance: %v", err)
	}
}

func TestInsertThenDelete(t *testing.T) {
	tc := testutil.SystemTest(t)
	testInsert(t)
	// Wait for instance to be created before deleting.
	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		_, err := get(tc.ProjectID, "us-central1-b", "dca-test-instance-2")
		if err == nil {
			break
		} else if i == 9 {
			t.Fatalf("failed to get instance: %v", err)
		}
	}
	testDelete(t)
}

func testInsert(t *testing.T) {
	tc := testutil.SystemTest(t)
	operation, err := insert(tc.ProjectID, "us-central1-b", "dca-test-instance-2")
	if err != nil {
		t.Fatalf("failed to insert new instance: %v", err)
	}
	fmt.Println(operation)
}

func testDelete(t *testing.T) {
	tc := testutil.SystemTest(t)
	operation, err := delete(tc.ProjectID, "us-central1-b", "dca-test-instance-2")
	if err != nil {
		t.Fatalf("failed to delete instance: %v", err)
	}
	fmt.Println(operation)
}
