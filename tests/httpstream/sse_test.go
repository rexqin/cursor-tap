package httpstream_test

import (
	"strings"
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
)

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

	parser := httpstream.NewSSEParser(strings.NewReader(input))
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

func TestSSEParserReadAll(t *testing.T) {
	input := "data: one\n\ndata: two\n\n"
	events, err := httpstream.NewSSEParser(strings.NewReader(input)).ReadAll()
	if err != nil {
		t.Fatalf("ReadAll: %v", err)
	}
	if len(events) != 2 {
		t.Fatalf("got %d events, want 2", len(events))
	}
	if events[0].Data != "one" || events[1].Data != "two" {
		t.Fatalf("unexpected events: %+v", events)
	}
}
