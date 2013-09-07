// This file is subject to a 1-clause BSD license.
// Its contents can be found in the enclosed LICENSE file.

package evdev

import (
	"testing"
)

func TestBitset(t *testing.T) {
	bs := NewBitset(80)
	bs.Set(0)
	bs.Set(2)
	bs.Set(4)
	bs.Set(13)
	bs.Set(76)

	want := []struct {
		Index int
		Value bool
	}{
		{0, true},
		{1, false},
		{2, true},
		{3, false},
		{4, true},
		{5, false},
		{9, false},
		{50, false},
		{76, true},
		{80, false},
	}

	for _, w := range want {
		if bs.Test(w.Index) != w.Value {
			t.Fatalf("Index %d: Want %v", w.Index, w.Value)
		}
	}
}
