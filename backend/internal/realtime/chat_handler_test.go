package realtime

import "testing"

func TestTargetRoomAllowed(t *testing.T) {
	cases := []struct {
		connRoom, envRoom string
		want              bool
	}{
		{"room-a", "room-a", true},
		// Empty envelope room_id is valid — the FE can send it before its
		// store hydrates; the connection's own room is used. Regression lock:
		// this case was invisible while the check was disabled in dev.
		{"room-a", "", true},
		{"room-a", "room-b", false},
	}
	for _, c := range cases {
		if got := targetRoomAllowed(c.connRoom, c.envRoom); got != c.want {
			t.Errorf("targetRoomAllowed(%q, %q) = %v, want %v", c.connRoom, c.envRoom, got, c.want)
		}
	}
}
