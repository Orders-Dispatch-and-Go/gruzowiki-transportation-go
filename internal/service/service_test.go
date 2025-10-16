package service

import "testing"

func TestExample(t *testing.T) {
	t.Parallel()

	if 2+2 != 4 {
		t.Fatal("2 + 2 != 4")
	}
}
