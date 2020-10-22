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

// Package kms contains samples for asymmetric keys feature of Cloud Key Management Service
// https://cloud.google.com/kms/
package kms

import (
	"bytes"
	"log"
	"os"
	"strings"
	"testing"

	"github.com/shinfan/google-dca-test/internal/testutil"
)

var fixture *kmsFixture

func TestMain(m *testing.M) {
	tc, ok := testutil.ContextMain(m)
	if !ok {
		log.Print("skipping - unset GOLANG_SAMPLES_PROJECT_ID?")
		return
	}

	var err error
	fixture, err = NewKMSFixture(tc.ProjectID)
	if err != nil {
		log.Fatalf("failed to create fixture: %s", err)
	}

	exitCode := m.Run()

	if err := fixture.Cleanup(); err != nil {
		log.Fatalf("failed to cleanup resources: %s", err)
	}

	os.Exit(exitCode)
}

func TestCreateKeyAsymmetricDecrypt(t *testing.T) {
	testutil.SystemTest(t)

	parent, id := fixture.KeyRingName, fixture.RandomID()

	var b bytes.Buffer
	if err := createKeyAsymmetricDecrypt(&b, parent, id); err != nil {
		t.Fatal(err)
	}

	if got, want := b.String(), "Created key:"; !strings.Contains(got, want) {
		t.Errorf("createKeyAsymmetricDecrypt: expected %q to contain %q", got, want)
	}
}
