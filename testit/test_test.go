package test

import "testing"

func TestAddition(t *testing.T) {
	if 2+2 != 4 {
		t.Error("Addition test failed")
	}
}
