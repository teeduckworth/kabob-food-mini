package http

import (
	"reflect"
	"testing"
)

func TestNormalizeOrigins(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "normalizes case, trims whitespace and slashes",
			input:    []string{" https://Example.com/ ", "https://example.com//", "http://Localhost:3000/"},
			expected: []string{"https://example.com", "http://localhost:3000"},
		},
		{
			name:     "drops empty entries",
			input:    []string{"", "   ", "/"},
			expected: nil,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizeOrigins(tc.input); !reflect.DeepEqual(got, tc.expected) {
				t.Fatalf("normalizeOrigins() = %v, expected %v", got, tc.expected)
			}
		})
	}
}

func TestBuildCORSConfig(t *testing.T) {
	t.Parallel()

	cfg := buildCORSConfig([]string{" https://Example.com/ "})
	if cfg.AllowAllOrigins {
		t.Fatal("expected specific origins to disable AllowAllOrigins")
	}
	if len(cfg.AllowOrigins) != 1 || cfg.AllowOrigins[0] != "https://example.com" {
		t.Fatalf("unexpected AllowOrigins: %v", cfg.AllowOrigins)
	}

	cfg = buildCORSConfig(nil)
	if !cfg.AllowAllOrigins {
		t.Fatal("expected nil origins to enable AllowAllOrigins")
	}
}
