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
	"context"
	"fmt"
	"google.golang.org/api/compute/v1"
	"google.golang.org/api/option"
)

func aggregatedList(projectID string) (*compute.InstanceAggregatedList, error) {
	ctx := context.Background()
	computeService, err := compute.NewService(ctx, option.WithEndpoint("https://compute.mtls.googleapis.com/compute/v1/projects/"))
	if err != nil {
		return nil, fmt.Errorf("compute.NewService: %v", err)
	}
	resp, err := computeService.Instances.AggregatedList(projectID).Do()
	if err != nil {
		return nil, fmt.Errorf("computeService.Instances.AggregatedList: %v", err)
	}
	return resp, nil
}
