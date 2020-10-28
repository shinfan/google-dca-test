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

// Package topics is a tool to manage Google Cloud Pub/Sub topics by using the Pub/Sub API.
// See more about Google Cloud Pub/Sub at https://cloud.google.com/pubsub/docs/overview.package topics
package topics

import (
	"bytes"
	"context"
	"sync"
	"testing"

	"cloud.google.com/go/pubsub"
	"github.com/shinfan/google-dca-test/internal/testutil"
)

var topicID string

// once guards cleanup related operations in setup. No need to set up and tear
// down every time, so this speeds things up.
var once sync.Once

func setup(t *testing.T) *pubsub.Client {
	ctx := context.Background()
	tc := testutil.SystemTest(t)

	topicID = "test-topic"
	var err error
	println(tc.ProjectID)
	client, err := pubsub.NewClient(ctx, tc.ProjectID)
	if err != nil {
		t.Fatalf("failed to create client: %v", err.Error())
	}

	// Cleanup resources from the previous tests.
	once.Do(func() {
		topic := client.Topic(topicID)
		ok, err := topic.Exists(ctx)
		if err != nil {
			t.Fatalf("failed to check if topic exists: %v", err)
		}
		if ok {
			if err := topic.Delete(ctx); err != nil {
				t.Fatalf("failed to cleanup the topic (%q): %v", topicID, err)
			}
		}
	})

	return client
}

func TestCreate(t *testing.T) {
	println("start test")
	client := setup(t)
	tc := testutil.SystemTest(t)
	buf := new(bytes.Buffer)
	if err := create(buf, tc.ProjectID, topicID); err != nil {
		t.Fatalf("failed to create a topic: %v", err)
	}
	ok, err := client.Topic(topicID).Exists(context.Background())
	if err != nil {
		t.Fatalf("failed to check if topic exists: %v", err)
	}
	if !ok {
		t.Fatalf("got none; want topic = %q", topicID)
	}
}
