package log

import "testing"


func TestLog(t *testing.T) {

	log := Named("test")

	log.Info("test log")

	t.Log("X")
}
