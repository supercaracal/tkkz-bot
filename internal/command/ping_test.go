package command

import (
	"testing"
)

func TestGetPingReply(t *testing.T) {
	if GetPingReply() != "PONG" {
		t.Errorf("We didn't get expected reply")
	}
}
