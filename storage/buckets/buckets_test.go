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

package buckets

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/shinfan/google-dca-test/internal/testutil"
)

func TestCreate(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	// Clean up bucket before running tests.
	deleteBucket(ioutil.Discard, bucketName)
	if err := createBucket(ioutil.Discard, tc.ProjectID, bucketName); err != nil {
		t.Fatalf("createBucket: %v", err)
	}
}

func TestCreateBucketClassLocation(t *testing.T) {
	tc := testutil.SystemTest(t)
	name := tc.ProjectID + "-storage-buckets-tests-attrs"

	// Clean up bucket before running the test.
	deleteBucket(ioutil.Discard, name)
	if err := createBucketClassLocation(ioutil.Discard, tc.ProjectID, name); err != nil {
		t.Fatalf("createBucketClassLocation: %v", err)
	}
	if err := deleteBucket(ioutil.Discard, name); err != nil {
		t.Fatalf("deleteBucket: %v", err)
	}
}

func TestListBuckets(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	buckets, err := listBuckets(ioutil.Discard, tc.ProjectID)
	if err != nil {
		t.Fatalf("listBuckets: %v", err)
	}

	var ok bool
	testutil.Retry(t, 5, 2*time.Second, func(r *testutil.R) { // for eventual consistency
		for _, b := range buckets {
			if b == bucketName {
				ok = true
				break
			}
		}
		if !ok {
			r.Errorf("got bucket list: %v; want %q in the list", buckets, bucketName)
		}
	})
}

func TestGetBucketMetadata(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	buf := new(bytes.Buffer)
	if _, err := getBucketMetadata(buf, bucketName); err != nil {
		t.Errorf("getBucketMetadata: %#v", err)
	}

	got := buf.String()
	if want := "BucketName:"; !strings.Contains(got, want) {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestIAM(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	if _, err := getBucketPolicy(ioutil.Discard, bucketName); err != nil {
		t.Errorf("getBucketPolicy: %#v", err)
	}
	if err := addBucketIAMMember(ioutil.Discard, bucketName); err != nil {
		t.Errorf("addBucketIAMMember: %v", err)
	}
	if err := removeBucketIAMMember(ioutil.Discard, bucketName); err != nil {
		t.Errorf("removeBucketIAMMember: %v", err)
	}

	// Uniform bucket-level access is required to use IAM with conditions.
	if err := enableUniformBucketLevelAccess(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("enableUniformBucketLevelAccess:  %v", err)
	}

	role := "roles/storage.objectViewer"
	member := "group:cloud-logs@google.com"
	title := "title"
	description := "description"
	expression := "resource.name.startsWith(\"projects/_/buckets/bucket-name/objects/prefix-a-\")"

	if err := addBucketConditionalIAMBinding(ioutil.Discard, bucketName, role, member, title, description, expression); err != nil {
		t.Errorf("addBucketConditionalIAMBinding: %v", err)
	}
	if err := removeBucketConditionalIAMBinding(ioutil.Discard, bucketName, role, title, description, expression); err != nil {
		t.Errorf("removeBucketConditionalIAMBinding: %v", err)
	}
}

func TestRequesterPays(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	if err := enableRequesterPays(ioutil.Discard, bucketName); err != nil {
		t.Errorf("enableRequesterPays: %#v", err)
	}
	if err := disableRequesterPays(ioutil.Discard, bucketName); err != nil {
		t.Errorf("disableRequesterPays: %#v", err)
	}
	if err := getRequesterPaysStatus(ioutil.Discard, bucketName); err != nil {
		t.Errorf("getRequesterPaysStatus: %#v", err)
	}
}

func TestKMS(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")

	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)
	if err := setBucketDefaultKMSKey(ioutil.Discard, bucketName, kmsKeyName); err != nil {
		t.Fatalf("setBucketDefaultKmsKey: failed to enable default kms key (%q): %v", kmsKeyName, err)
	}
}


func TestUniformBucketLevelAccess(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	testutil.Retry(t, 10, 10*time.Second, func(r *testutil.R) {
		if err := enableUniformBucketLevelAccess(ioutil.Discard, bucketName); err != nil {
			r.Errorf("enableUniformBucketLevelAccess: %v", err)
		}
	})

	attrs, err := getUniformBucketLevelAccess(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getUniformBucketLevelAccess: %v", err)
	}
	if !attrs.UniformBucketLevelAccess.Enabled {
		t.Fatalf("Uniform bucket-level access was not enabled for (%q).", bucketName)
	}

	testutil.Retry(t, 2, 10*time.Second, func(r *testutil.R) {
		if err := disableUniformBucketLevelAccess(ioutil.Discard, bucketName); err != nil {
			r.Errorf("disableUniformBucketLevelAccess: %v", err)
		}
	})

	attrs, err = getUniformBucketLevelAccess(ioutil.Discard, bucketName)
	if err != nil {
		t.Fatalf("getUniformBucketLevelAccess: %v", err)
	}
	if attrs.UniformBucketLevelAccess.Enabled {
		t.Fatalf("Uniform bucket-level access was not disabled for (%q).", bucketName)
	}
}

func TestDelete(t *testing.T) {
	tc := testutil.SystemTest(t)
	bucketName := tc.ProjectID + "-storage-buckets-tests"

	if err := deleteBucket(ioutil.Discard, bucketName); err != nil {
		t.Fatalf("deleteBucket: %v", err)
	}
}
