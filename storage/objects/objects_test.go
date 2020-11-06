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

package objects

import (
	"bytes"
	"context"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"cloud.google.com/go/storage"
	"github.com/shinfan/google-dca-test/internal/testutil"
	"google.golang.org/api/option"

)

// TestObjects runs all samples tests of the package.
func TestObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithEndpoint("https://storage.mtls.googleapis.com/storage/v1/"))
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	var (
		bucket                = tc.ProjectID + "-samples-object-bucket-1"
		dstBucket             = tc.ProjectID + "-samples-object-bucket-2"
		bucketVersioning      = tc.ProjectID + "-bucket-versioning-enabled"
		object1               = "foo.txt"
		object2               = "foo/a.txt"
		object3               = "bar.txt"
		allAuthenticatedUsers = storage.AllAuthenticatedUsers
		roleReader            = storage.RoleReader
	)

	testutil.CleanBucket(ctx, t, tc.ProjectID, bucket)
	testutil.CleanBucket(ctx, t, tc.ProjectID, dstBucket)
	testutil.CleanBucket(ctx, t, tc.ProjectID, bucketVersioning)

	if err := enableVersioning(ioutil.Discard, bucketVersioning); err != nil {
		t.Fatalf("enableVersioning: %v", err)
	}

	if err := uploadFile(ioutil.Discard, bucket, object1); err != nil {
		t.Fatalf("uploadFile(%q): %v", object1, err)
	}
	if err := uploadFile(ioutil.Discard, bucket, object2); err != nil {
		t.Fatalf("uploadFile(%q): %v", object2, err)
	}
	if err := uploadFile(ioutil.Discard, bucketVersioning, object1); err != nil {
		t.Fatalf("uploadFile(%q): %v", object1, err)
	}
	// Check enableVersioning correctly work.
	bkt := client.Bucket(bucketVersioning)
	bAttrs, err := bkt.Attrs(ctx)
	if !bAttrs.VersioningEnabled {
		t.Fatalf("object versioning is not enabled")
	}
	obj := bkt.Object(object1)
	attrs, err := obj.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Object(%q).Attrs: %v", bucketVersioning, object1, err)
	}
	// Keep the original generation of object1 before re-uploading
	// to use in the versioning samples.
	gen := attrs.Generation
	if err := uploadFile(ioutil.Discard, bucketVersioning, object1); err != nil {
		t.Fatalf("uploadFile(%q): %v", object1, err)
	}

	{
		// Should only show "foo/a.txt", not "foo.txt"
		var buf bytes.Buffer
		if err := listFiles(&buf, bucket); err != nil {
			t.Fatalf("listFiles: %v", err)
		}
		if got, want := buf.String(), object1; !strings.Contains(got, want) {
			t.Errorf("List() got %q; want to contain %q", got, want)
		}
		if got, want := buf.String(), object2; !strings.Contains(got, want) {
			t.Errorf("List() got %q; want to contain %q", got, want)
		}
	}

	{
		// Should only show "foo/a.txt", not "foo.txt"
		const prefix = "foo/"
		var buf bytes.Buffer
		if err := listFilesWithPrefix(&buf, bucket, prefix, ""); err != nil {
			t.Fatalf("listFilesWithPrefix: %v", err)
		}
		if got, want := buf.String(), object1; strings.Contains(got, want) {
			t.Errorf("List(%q) got %q; want NOT to contain %q", prefix, got, want)
		}
		if got, want := buf.String(), object2; !strings.Contains(got, want) {
			t.Errorf("List(%q) got %q; want to contain %q", prefix, got, want)
		}
	}

	{
		// Should show 2 versions of foo.txt
		var buf bytes.Buffer
		if err := listFilesAllVersion(&buf, bucketVersioning); err != nil {
			t.Fatalf("listFilesAllVersion: %v", err)
		}

		i := 0
		for _, line := range strings.Split(strings.TrimSuffix(buf.String(), "\n"), "\n") {
			if got, want := line, object1; !strings.Contains(got, want) {
				t.Errorf("List(Versions: true) got %q; want to contain %q", got, want)
			}
			i++
		}
		if i != 2 {
			t.Errorf("listFilesAllVersion should show 2 versions of foo.txt; got %d", i)
		}
	}

	if err := copyOldVersionOfObject(ioutil.Discard, bucketVersioning, object1, object3, gen); err != nil {
		t.Fatalf("copyOldVersionOfObject: %v", err)
	}

	data, err := downloadFile(ioutil.Discard, bucket, object1)
	if err != nil {
		println("bucket downloading: " + bucket)
		t.Fatalf("downloadFile: %v", err)
	}
	if got, want := string(data), "Hello\nworld"; got != want {
		t.Errorf("contents = %q; want %q", got, want)
	}

	_, err = getMetadata(ioutil.Discard, bucket, object1)
	if err != nil {
		t.Errorf("getMetadata: %v", err)
	}
	if err := makePublic(ioutil.Discard, bucket, object1, allAuthenticatedUsers, roleReader); err != nil {
		t.Errorf("makePublic: %v", err)
	}

	err = moveFile(ioutil.Discard, bucket, object1)
	if err != nil {
		t.Fatalf("moveFile: %v", err)
	}
	// object1's new name.
	object1 = object1 + "-rename"

	if err := copyFile(ioutil.Discard, dstBucket, bucket, object1); err != nil {
		t.Errorf("copyFile: %v", err)
	}

	key := []byte("my-secret-AES-256-encryption-key")
	newKey := []byte("My-secret-AES-256-encryption-key")

	if err := generateEncryptionKey(ioutil.Discard); err != nil {
		t.Errorf("generateEncryptionKey: %v", err)
	}
	if err := uploadEncryptedFile(ioutil.Discard, bucket, object1, key); err != nil {
		t.Errorf("uploadEncryptedFile: %v", err)
	}
	if err := rotateEncryptionKey(ioutil.Discard, bucket, object1, key, newKey); err != nil {
		t.Errorf("rotateEncryptionKey: %v", err)
	}
	if err := deleteFile(ioutil.Discard, bucket, object1); err != nil {
		t.Errorf("deleteFile: %v", err)
	}
	if err := deleteFile(ioutil.Discard, bucket, object2); err != nil {
		t.Errorf("deleteFile: %v", err)
	}
	if err := disableVersioning(ioutil.Discard, bucketVersioning); err != nil {
		t.Fatalf("disableVersioning: %v", err)
	}

	bAttrs, err = bkt.Attrs(ctx)
	if err != nil {
		t.Fatalf("Bucket(%q).Attrs: %v", bucketVersioning, err)
	}
	if bAttrs.VersioningEnabled {
		t.Fatalf("object versioning is not disabled")
	}
	testutil.Retry(t, 2, time.Second, func(r *testutil.R) {
		// Cleanup, this part won't be executed if Fatal happens.
		// TODO(jbd): Implement garbage cleaning.
		if err := client.Bucket(bucket).Delete(ctx); err != nil {
			r.Errorf("Bucket(%q).Delete: %v", bucket, err)
		}
	})

	testutil.Retry(t, 2, time.Second, func(r *testutil.R) {
		if err := deleteFile(ioutil.Discard, dstBucket, object1+"-copy"); err != nil {
			r.Errorf("deleteFile: %v", err)
		}
	})

	testutil.Retry(t, 2, time.Second, func(r *testutil.R) {
		if err := client.Bucket(dstBucket).Delete(ctx); err != nil {
			r.Errorf("Bucket(%q).Delete: %v", dstBucket, err)
		}
	})

	// CleanBucket to delete versioned objects in bucket
	testutil.CleanBucket(ctx, t, tc.ProjectID, bucketVersioning)
	testutil.Retry(t, 2, time.Second, func(r *testutil.R) {
		if err := client.Bucket(bucketVersioning).Delete(ctx); err != nil {
			r.Errorf("Bucket(%q).Delete: %v", bucketVersioning, err)
		}
	})
}

func TestKMSObjects(t *testing.T) {
	tc := testutil.SystemTest(t)
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithEndpoint("https://storage.mtls.googleapis.com/storage/v1/"))
	if err != nil {
		t.Fatalf("storage.NewClient: %v", err)
	}
	defer client.Close()

	keyRingID := os.Getenv("GOLANG_SAMPLES_KMS_KEYRING")
	cryptoKeyID := os.Getenv("GOLANG_SAMPLES_KMS_CRYPTOKEY")
	if keyRingID == "" || cryptoKeyID == "" {
		t.Skip("GOLANG_SAMPLES_KMS_KEYRING and GOLANG_SAMPLES_KMS_CRYPTOKEY must be set")
	}

	bucket := tc.ProjectID + "-samples-object-bucket-1"
	object := "foo.txt"

	testutil.CleanBucket(ctx, t, tc.ProjectID, bucket)

	kmsKeyName := fmt.Sprintf("projects/%s/locations/%s/keyRings/%s/cryptoKeys/%s", tc.ProjectID, "global", keyRingID, cryptoKeyID)

	testutil.Retry(t, 10, time.Second, func(r *testutil.R) {
		if err := uploadWithKMSKey(ioutil.Discard, bucket, object, kmsKeyName); err != nil {
			r.Errorf("uploadWithKMSKey: %v", err)
		}
	})
}
