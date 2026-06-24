package httpstream_test

import (
	"testing"

	"github.com/burpheart/cursor-tap/internal/httpstream"
	"github.com/stretchr/testify/require"
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
		t.Run(tt.want, func(t *testing.T) {
			require.Equal(t, tt.want, tt.dir.String())
		})
	}
}
