package util

import "testing"

// this can cause issues on some more obscure build systems.
func TestPrngIsAvailable(t *testing.T) {
	prng := PrngIsAvailable()
	if prng != true {
		t.Error("rand assertion failed")
	}
}
