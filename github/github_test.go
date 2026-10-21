// Copyright 2026 OpenSSF Scorecard Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package github

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestParseFromURL(t *testing.T) {
	const token = "ghs_token"
	t.Setenv("GITHUB_AUTH_TOKEN", token)

	var auth string
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		auth = r.Header.Get("Authorization")
		fmt.Fprint(w, `{"default_branch":"main","fork":false,"private":true}`)
	}))
	t.Cleanup(server.Close)

	client := NewClient(context.Background())
	info, err := client.ParseFromURL(server.URL, "ossf-tests/scorecard-action")
	if err != nil {
		t.Fatalf("Unexpected error ParseFromURL: %v", err)
	}
	if want := "Bearer " + token; auth != want {
		t.Errorf("ParseFromURL() sent Authorization %q, wanted %q", auth, want)
	}
	if info.Repo.DefaultBranch == nil || *info.Repo.DefaultBranch != "main" {
		t.Errorf("ParseFromURL() default branch: got %v, wanted main", info.Repo.DefaultBranch)
	}
}
