package cmd

import "testing"

func TestParseSize(t *testing.T) {
	tests := []struct {
		input string
		want  int64
	}{
		{"0", 0},
		{"", 0},
		{"100", 100},
		{"100B", 100},
		{"1K", 1024},
		{"1KB", 1024},
		{"1k", 1024},
		{"1M", 1024 * 1024},
		{"1MB", 1024 * 1024},
		{"1m", 1024 * 1024},
		{"1G", 1024 * 1024 * 1024},
		{"10M", 10 * 1024 * 1024},
		{"1.5M", int64(1.5 * 1024 * 1024)},
		{"512K", 512 * 1024},
	}

	for _, tt := range tests {
		got, err := parseSize(tt.input)
		if err != nil {
			t.Errorf("parseSize(%q) error: %v", tt.input, err)
			continue
		}
		if got != tt.want {
			t.Errorf("parseSize(%q) = %d, want %d", tt.input, got, tt.want)
		}
	}
}

func TestParseSize_Errors(t *testing.T) {
	bad := []string{"abc", "1X", "1TB"}
	for _, s := range bad {
		_, err := parseSize(s)
		if err == nil {
			t.Errorf("parseSize(%q) expected error, got nil", s)
		}
	}
}
