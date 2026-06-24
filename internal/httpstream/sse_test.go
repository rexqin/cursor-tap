package httpstream

import (
	"strings"
	"testing"
)

func TestParseSSEField(t *testing.T) {
	tests := []struct {
		line      string
		wantField string
		wantValue string
	}{
		{"data: hello", "data", "hello"},
		{"event:message", "event", "message"},
		{"id: 42", "id", "42"},
		{"retry: 3000", "retry", "3000"},
		{"comment only", "comment", "only"},
	}

	for _, tt := range tests {
		field, value := parseSSEField([]byte(tt.line))
		if field != tt.wantField || value != tt.wantValue {
			t.Errorf("parseSSEField(%q) = (%q, %q), want (%q, %q)",
				tt.line, field, value, tt.wantField, tt.wantValue)
		}
	}
}

func TestSSEParserNext(t *testing.T) {
	input := strings.Join([]string{
		"event: message",
		"data: hello",
		"data: world",
		"",
		": comment",
		"data: second",
		"",
	}, "\n")

	parser := NewSSEParser(strings.NewReader(input))
	event, err := parser.Next()
	if err != nil {
		t.Fatalf("Next: %v", err)
	}
	if event.Event != "message" {
		t.Errorf("Event = %q, want message", event.Event)
	}
	if event.Data != "hello\nworld" {
		t.Errorf("Data = %q, want hello\\nworld", event.Data)
	}

	event, err = parser.Next()
	if err != nil {
		t.Fatalf("Next second: %v", err)
	}
	if strings.TrimSuffix(event.Data, "\n") != "second" {
		t.Errorf("Data = %q, want second", event.Data)
	}
}
