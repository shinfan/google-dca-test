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
	"testing"
	"time"

	"google.golang.org/api/option"
	"cloud.google.com/go/logging"
	"cloud.google.com/go/logging/logadmin"

	"github.com/shinfan/google-dca-test/internal/testutil"
)

func TestSimplelog(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()

	client, err := logging.NewClient(ctx, tc.ProjectID, option.WithEndpoint("logging.mtls.googleapis.com:443"))
	if err != nil {
		t.Fatalf("logging.NewClient: %v", err)
	}
	adminClient, err := logadmin.NewClient(ctx, tc.ProjectID, option.WithEndpoint("logging.mtls.googleapis.com:443"))
	if err != nil {
		t.Fatalf("logadmin.NewClient: %v", err)
	}
	defer func() {
		println("Close Client")
		if err := client.Close(); err != nil {
			t.Errorf("Close: %v", err)
		}
	}()

	defer func() {
		println("Deleting log")
		testutil.Retry(t, 1, 5*time.Second, func(r *testutil.R) {
			if err := deleteLog(adminClient); err != nil {
				r.Errorf("deleteLog: %v", err)
			}
		})
	}()

	client.OnError = func(err error) {
		t.Errorf("OnError: %v", err)
	}

	writeEntry(client)
	structuredWrite(client)
}
