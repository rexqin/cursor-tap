package httpstream_test

import (
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
)

func TestDirectionString(t *testing.T) {
	tests := []struct {
		dir  httpstream.Direction
		want string
	}{
		{httpstream.ClientToServer, "C2S"},
		{httpstream.ServerToClient, "S2C"},
		{httpstream.Direction(99), "S2C"},
	}

	for _, tt := range tests {
		if got := tt.dir.String(); got != tt.want {
			t.Errorf("Direction(%d).String() = %q, want %q", tt.dir, got, tt.want)
		}
	}
}
