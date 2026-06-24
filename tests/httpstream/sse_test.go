package httpstream_test

import (
	"strings"
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
	"github.com/stretchr/testify/require"
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
	require.NoError(t, err)
	require.Equal(t, "message", event.Event)
	require.Equal(t, "hello\nworld", event.Data)

	event, err = parser.Next()
	require.NoError(t, err)
	require.Equal(t, "second", strings.TrimSuffix(event.Data, "\n"))
}

func TestSSEParserReadAll(t *testing.T) {
	input := "data: one\n\ndata: two\n\n"
	events, err := httpstream.NewSSEParser(strings.NewReader(input)).ReadAll()
	require.NoError(t, err)
	require.Len(t, events, 2)
	require.Equal(t, "one", events[0].Data)
	require.Equal(t, "two", events[1].Data)
}
