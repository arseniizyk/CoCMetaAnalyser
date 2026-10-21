package cli

import (
	"testing"
)

func TestNormalizeSeason(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "already normalized YYYY-MM",
			input:    "2025-01",
			expected: "2025-01",
		},
		{
			name:     "mid-year YYYY-MM",
			input:    "2025-05",
			expected: "2025-05",
		},
		{
			name:     "legacy weekly date closing October season",
			input:    "2025-11-03",
			expected: "2025-10",
		},
		{
			name:     "legacy weekly date closing November season",
			input:    "2025-12-01",
			expected: "2025-11",
		},
		{
			name:     "legacy weekly date closing December season",
			input:    "2025-12-29",
			expected: "2025-12",
		},
		{
			name:     "legacy weekly date closing January 2026",
			input:    "2026-01-26",
			expected: "2026-01",
		},
		{
			name:     "supercell v2 May 2026 season",
			input:    "v2-2026-05-11T05:00:00Z",
			expected: "2026-05",
		},
		{
			name:     "supercell v2 June 2026 season",
			input:    "v2-2026-06-08T05:00:00Z",
			expected: "2026-06",
		},
		{
			name:     "supercell v2 late August covering September",
			input:    "v2-2026-08-31T05:00:00Z",
			expected: "2026-09",
		},
		{
			name:     "supercell v2 late September covering October",
			input:    "v2-2026-09-28T05:00:00Z",
			expected: "2026-10",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NormalizeSeason(tt.input)
			if got != tt.expected {
				t.Errorf("NormalizeSeason(%q) = %q; want %q", tt.input, got, tt.expected)
			}
		})
	}
}

func TestResolveBaseDir(t *testing.T) {
	// Custom dir should be respected
	custom := "/tmp/custom-meta"
	got := resolveBaseDir(custom)
	if got != custom {
		t.Errorf("resolveBaseDir(%q) = %q; want %q", custom, got, custom)
	}

	// Default resolution should fallback to a valid string
	def := resolveBaseDir("")
	if def == "" {
		t.Errorf("resolveBaseDir(\"\") returned empty string")
	}
}
