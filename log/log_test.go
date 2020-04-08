package log

import "testing"

func TestLog(t *testing.T) {
	Info("test log")

	t.Log("X")
}
