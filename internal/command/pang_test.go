package command

import (
	"testing"
)

func TestGetPangReply(t *testing.T) {
	if GetPangReply() != "PANG" {
		t.Errorf("We didn't get expected reply")
	}
}
