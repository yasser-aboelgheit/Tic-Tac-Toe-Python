package metric

import "testing"

func TestNoStatsIncr(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("unexpected panic: %v", err)
		}
	}()

	Incr("test", nil, 1)
}

func TestStopNoPanic(t *testing.T) {
	defer func() {
		if err := recover(); err != nil {
			t.Fatalf("unexpected panic: %v", err)
		}
	}()

	Stop()
}
