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
package main

import (
	"context"
	"fmt"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
	"log"
)

// Test program for interacting with Compute service.
func main() {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithEndpoint("https://compute.mtls.googleapis.com/compute/v1/projects/"))
	if err != nil {
		log.Fatal(err)
	}
	aggregatedList, err := computeService.Instances.AggregatedList(
		"shinfan-mtls-demo").Do()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("SUCCESS")
	fmt.Println(aggregatedList)
}
