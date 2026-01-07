package types

import (
	"reflect"
	"testing"
)

func TestSimpleRepo_ValidateRepo(t *testing.T) {
	tests := []struct {
		name    string
		repo    SimpleRepo
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid repo with input only",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid repo with input and output (different URLs)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid repo with input and output (same URL, different branch)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("feature"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid repo with input and output (different URL, same branch)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("main"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid repo with input only (no branch)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL: "https://github.com/user/repo",
				},
			},
			wantErr: false,
		},
		{
			name: "valid repo with autoPush",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
				AutoPush: BoolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "missing input",
			repo: SimpleRepo{
				Output: &RepoLocation{
					URL: "https://github.com/user/fork",
				},
			},
			wantErr: true,
			errMsg:  "input is required",
		},
		{
			name: "nil input",
			repo: SimpleRepo{
				Input: nil,
			},
			wantErr: true,
			errMsg:  "input is required",
		},
		{
			name: "empty input URL",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "",
					Branch: StringPtr("main"),
				},
			},
			wantErr: true,
			errMsg:  "input.url is required",
		},
		{
			name: "whitespace-only input URL",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "   ",
					Branch: StringPtr("main"),
				},
			},
			wantErr: true,
			errMsg:  "input.url is required",
		},
		{
			name: "identical input and output (same URL and branch)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
			},
			wantErr: true,
			errMsg:  "output repository must differ from input (different URL or branch required)",
		},
		{
			name: "identical input and output (same URL, no branches)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL: "https://github.com/user/repo",
				},
				Output: &RepoLocation{
					URL: "https://github.com/user/repo",
				},
			},
			wantErr: true,
			errMsg:  "output repository must differ from input (different URL or branch required)",
		},
		{
			name: "identical input and output (same URL, both nil branches)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: nil,
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: nil,
				},
			},
			wantErr: true,
			errMsg:  "output repository must differ from input (different URL or branch required)",
		},
		{
			name: "identical input and output (same URL, empty string branches)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr(""),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr(""),
				},
			},
			wantErr: true,
			errMsg:  "output repository must differ from input (different URL or branch required)",
		},
		{
			name: "identical input and output (whitespace branches treated as identical)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("  "),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("  "),
				},
			},
			wantErr: true,
			errMsg:  "output repository must differ from input (different URL or branch required)",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.repo.ValidateRepo()

			if tt.wantErr {
				if err == nil {
					t.Errorf("ValidateRepo() expected error, got nil")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ValidateRepo() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ValidateRepo() unexpected error = %v", err)
			}
		})
	}
}

func TestSimpleRepo_ToMapForCR(t *testing.T) {
	tests := []struct {
		name string
		repo SimpleRepo
		want map[string]interface{}
	}{
		{
			name: "input only with branch",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
			},
		},
		{
			name: "input only without branch",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL: "https://github.com/user/repo",
				},
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"url": "https://github.com/user/repo",
				},
			},
		},
		{
			name: "input and output with branches",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/fork",
					"branch": "feature",
				},
			},
		},
		{
			name: "input and output without branches",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL: "https://github.com/user/repo",
				},
				Output: &RepoLocation{
					URL: "https://github.com/user/fork",
				},
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"url": "https://github.com/user/repo",
				},
				"output": map[string]interface{}{
					"url": "https://github.com/user/fork",
				},
			},
		},
		{
			name: "with autoPush true",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
				AutoPush: BoolPtr(true),
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/fork",
					"branch": "feature",
				},
				"autoPush": true,
			},
		},
		{
			name: "with autoPush false",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
				AutoPush: BoolPtr(false),
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/fork",
					"branch": "feature",
				},
				"autoPush": false,
			},
		},
		{
			name: "with nil autoPush (omitted)",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
				AutoPush: nil,
			},
			want: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/fork",
					"branch": "feature",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.repo.ToMapForCR()

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ToMapForCR() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSimpleRepo_RoundTrip(t *testing.T) {
	// Test that converting SimpleRepo -> map -> SimpleRepo preserves data
	tests := []struct {
		name string
		repo SimpleRepo
	}{
		{
			name: "input only",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
			},
		},
		{
			name: "input and output",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
			},
		},
		{
			name: "with autoPush true",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
				AutoPush: BoolPtr(true),
			},
		},
		{
			name: "with autoPush false",
			repo: SimpleRepo{
				Input: &RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: StringPtr("main"),
				},
				Output: &RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: StringPtr("feature"),
				},
				AutoPush: BoolPtr(false),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Convert to map
			m := tt.repo.ToMapForCR()

			// Convert back to SimpleRepo (using ParseRepoMap equivalent logic)
			var reconstructed SimpleRepo

			// Parse input
			if inputMap, ok := m["input"].(map[string]interface{}); ok {
				input := &RepoLocation{}
				if url, ok := inputMap["url"].(string); ok {
					input.URL = url
				}
				if branch, ok := inputMap["branch"].(string); ok {
					input.Branch = StringPtr(branch)
				}
				reconstructed.Input = input
			}

			// Parse output
			if outputMap, ok := m["output"].(map[string]interface{}); ok {
				output := &RepoLocation{}
				if url, ok := outputMap["url"].(string); ok {
					output.URL = url
				}
				if branch, ok := outputMap["branch"].(string); ok {
					output.Branch = StringPtr(branch)
				}
				reconstructed.Output = output
			}

			// Parse autoPush
			if autoPush, ok := m["autoPush"].(bool); ok {
				reconstructed.AutoPush = BoolPtr(autoPush)
			}

			// Compare original and reconstructed
			if !reposEqual(tt.repo, reconstructed) {
				t.Errorf("Round-trip failed:\nOriginal:      %+v\nReconstructed: %+v", tt.repo, reconstructed)
			}
		})
	}
}

// Helper function to compare SimpleRepo structs
func reposEqual(a, b SimpleRepo) bool {
	// Compare Input
	if !repoLocationsEqual(a.Input, b.Input) {
		return false
	}

	// Compare Output
	if !repoLocationsEqual(a.Output, b.Output) {
		return false
	}

	// Compare AutoPush
	if (a.AutoPush == nil) != (b.AutoPush == nil) {
		return false
	}
	if a.AutoPush != nil && b.AutoPush != nil && *a.AutoPush != *b.AutoPush {
		return false
	}

	return true
}

func repoLocationsEqual(a, b *RepoLocation) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	if a.URL != b.URL {
		return false
	}

	// Compare branches
	if (a.Branch == nil) != (b.Branch == nil) {
		return false
	}
	if a.Branch != nil && b.Branch != nil && *a.Branch != *b.Branch {
		return false
	}

	return true
}
