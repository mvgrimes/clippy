package cmd

import (
	"fmt"
	"strconv"
	"strings"
)

// parseSize parses a human-readable size string (e.g. "10M", "1G", "512K", "1024")
// into a byte count. Supported suffixes: B, K/KB, M/MB, G/GB (case-insensitive).
func parseSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" || s == "0" {
		return 0, nil
	}

	// Find where the numeric part ends.
	i := len(s)
	for j, c := range s {
		if (c < '0' || c > '9') && c != '.' {
			i = j
			break
		}
	}

	numStr := s[:i]
	suffix := strings.ToUpper(strings.TrimSpace(s[i:]))

	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size %q: %w", s, err)
	}

	var multiplier float64
	switch suffix {
	case "", "B":
		multiplier = 1
	case "K", "KB":
		multiplier = 1024
	case "M", "MB":
		multiplier = 1024 * 1024
	case "G", "GB":
		multiplier = 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("unknown size suffix %q in %q", suffix, s)
	}

	return int64(num * multiplier), nil
}
