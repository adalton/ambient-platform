package handlers

import (
	"ambient-code-backend/types"
	"testing"
)

func TestParseRepoMap_V2Format(t *testing.T) {
	tests := []struct {
		name    string
		input   map[string]interface{}
		want    types.SimpleRepo
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid V2 format with input only",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
			},
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: types.StringPtr("main"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid V2 format with input and output",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/fork",
					"branch": "feature",
				},
			},
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: types.StringPtr("main"),
				},
				Output: &types.RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: types.StringPtr("feature"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid V2 format with autoPush",
			input: map[string]interface{}{
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
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: types.StringPtr("main"),
				},
				Output: &types.RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: types.StringPtr("feature"),
				},
				AutoPush: types.BoolPtr(true),
			},
			wantErr: false,
		},
		{
			name: "valid V2 format with autoPush false",
			input: map[string]interface{}{
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
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: types.StringPtr("main"),
				},
				Output: &types.RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: types.StringPtr("feature"),
				},
				AutoPush: types.BoolPtr(false),
			},
			wantErr: false,
		},
		{
			name: "valid V2 format without branch in input",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url": "https://github.com/user/repo",
				},
			},
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: nil,
				},
			},
			wantErr: false,
		},
		{
			name: "valid V2 format with empty branch string (treated as nil)",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "",
				},
			},
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: nil, // Empty string is normalized to nil
				},
			},
			wantErr: false,
		},
		{
			name: "valid V2 format with whitespace-only branch (treated as nil)",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "   ",
				},
			},
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: nil, // Whitespace-only string is normalized to nil
				},
			},
			wantErr: false,
		},
		{
			name: "missing input field",
			input: map[string]interface{}{
				"output": map[string]interface{}{
					"url": "https://github.com/user/fork",
				},
			},
			wantErr: true,
			errMsg:  "input is required in repository configuration",
		},
		{
			name:    "empty repository map",
			input:   map[string]interface{}{},
			wantErr: true,
			errMsg:  "input is required in repository configuration",
		},
		{
			name: "input is not a map",
			input: map[string]interface{}{
				"input": "not-a-map",
			},
			wantErr: true,
			errMsg:  "input is required in repository configuration",
		},
		{
			name: "missing URL in input",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"branch": "main",
				},
			},
			wantErr: true,
			errMsg:  "input.url is required",
		},
		{
			name: "empty URL in input",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "",
					"branch": "main",
				},
			},
			wantErr: true,
			errMsg:  "input.url is required",
		},
		{
			name: "whitespace-only URL in input",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "   ",
					"branch": "main",
				},
			},
			wantErr: true,
			errMsg:  "input.url is required",
		},
		{
			name: "identical input and output (same URL and branch)",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
			},
			wantErr: true,
			errMsg:  "output repository must differ from input (different URL or branch required)",
		},
		{
			name: "identical input and output (same URL, no branches)",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url": "https://github.com/user/repo",
				},
				"output": map[string]interface{}{
					"url": "https://github.com/user/repo",
				},
			},
			wantErr: true,
			errMsg:  "output repository must differ from input (different URL or branch required)",
		},
		{
			name: "valid: same URL but different branch",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "feature",
				},
			},
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: types.StringPtr("main"),
				},
				Output: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: types.StringPtr("feature"),
				},
			},
			wantErr: false,
		},
		{
			name: "valid: different URL, same branch",
			input: map[string]interface{}{
				"input": map[string]interface{}{
					"url":    "https://github.com/user/repo",
					"branch": "main",
				},
				"output": map[string]interface{}{
					"url":    "https://github.com/user/fork",
					"branch": "main",
				},
			},
			want: types.SimpleRepo{
				Input: &types.RepoLocation{
					URL:    "https://github.com/user/repo",
					Branch: types.StringPtr("main"),
				},
				Output: &types.RepoLocation{
					URL:    "https://github.com/user/fork",
					Branch: types.StringPtr("main"),
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ParseRepoMap(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Errorf("ParseRepoMap() expected error, got nil")
					return
				}
				if tt.errMsg != "" && err.Error() != tt.errMsg {
					t.Errorf("ParseRepoMap() error = %v, want error containing %v", err, tt.errMsg)
				}
				return
			}

			if err != nil {
				t.Errorf("ParseRepoMap() unexpected error = %v", err)
				return
			}

			// Compare Input
			if got.Input == nil {
				if tt.want.Input != nil {
					t.Errorf("ParseRepoMap() got.Input = nil, want %+v", tt.want.Input)
				}
			} else {
				if tt.want.Input == nil {
					t.Errorf("ParseRepoMap() got.Input = %+v, want nil", got.Input)
				} else {
					if got.Input.URL != tt.want.Input.URL {
						t.Errorf("ParseRepoMap() got.Input.URL = %v, want %v", got.Input.URL, tt.want.Input.URL)
					}
					if !stringPtrEqual(got.Input.Branch, tt.want.Input.Branch) {
						t.Errorf("ParseRepoMap() got.Input.Branch = %v, want %v", stringPtrValue(got.Input.Branch), stringPtrValue(tt.want.Input.Branch))
					}
				}
			}

			// Compare Output
			if got.Output == nil {
				if tt.want.Output != nil {
					t.Errorf("ParseRepoMap() got.Output = nil, want %+v", tt.want.Output)
				}
			} else {
				if tt.want.Output == nil {
					t.Errorf("ParseRepoMap() got.Output = %+v, want nil", got.Output)
				} else {
					if got.Output.URL != tt.want.Output.URL {
						t.Errorf("ParseRepoMap() got.Output.URL = %v, want %v", got.Output.URL, tt.want.Output.URL)
					}
					if !stringPtrEqual(got.Output.Branch, tt.want.Output.Branch) {
						t.Errorf("ParseRepoMap() got.Output.Branch = %v, want %v", stringPtrValue(got.Output.Branch), stringPtrValue(tt.want.Output.Branch))
					}
				}
			}

			// Compare AutoPush
			if !boolPtrEqual(got.AutoPush, tt.want.AutoPush) {
				t.Errorf("ParseRepoMap() got.AutoPush = %v, want %v", boolPtrValue(got.AutoPush), boolPtrValue(tt.want.AutoPush))
			}
		})
	}
}

// Helper functions for pointer comparisons
func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func stringPtrValue(p *string) string {
	if p == nil {
		return "<nil>"
	}
	return *p
}

func boolPtrEqual(a, b *bool) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func boolPtrValue(p *bool) string {
	if p == nil {
		return "<nil>"
	}
	if *p {
		return "true"
	}
	return "false"
}
