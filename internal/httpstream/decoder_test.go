package httpstream

import "testing"

func TestParseContentEncoding(t *testing.T) {
	tests := []struct {
		name  string
		value string
		want  []string
	}{
		{"empty", "", nil},
		{"gzip", "gzip", []string{"gzip"}},
		{"multiple", "gzip, br", []string{"gzip", "br"}},
		{"identity skipped", "identity, gzip", []string{"gzip"}},
		{"whitespace", "  deflate , gzip  ", []string{"deflate", "gzip"}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseContentEncoding(tt.value)
			if len(got) != len(tt.want) {
				t.Fatalf("parseContentEncoding(%q) = %v, want %v", tt.value, got, tt.want)
			}
			for i := range got {
				if got[i] != tt.want[i] {
					t.Fatalf("parseContentEncoding(%q) = %v, want %v", tt.value, got, tt.want)
				}
			}
		})
	}
}
