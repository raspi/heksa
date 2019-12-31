package display

import "testing"

func TestNearest(t *testing.T) {

	type test struct {
		Width    uint8
		Expected uint8
		Actual   uint8
	}

	tests := []test{
		{
			Width:    1,
			Expected: 8,
		},
		{
			Width:    7,
			Expected: 8,
		},
		{
			Width:    8,
			Expected: 8,
		},
		{
			Width:    9,
			Expected: 16,
		},
		{
			Width:    15,
			Expected: 16,
		},
		{
			Width:    16,
			Expected: 16,
		},
		{
			Width:    17,
			Expected: 32,
		},
		{
			Width:    31,
			Expected: 32,
		},
		{
			Width:    32,
			Expected: 32,
		},
		{
			Width:    33,
			Expected: 64,
		},
		{
			Width:    63,
			Expected: 64,
		},
		{
			Width:    64,
			Expected: 64,
		},
	}

	failed := false

	for _, tst := range tests {
		tst.Actual = nearest(tst.Width)
		if tst.Actual != tst.Expected {
			t.Logf(`expected %d got %d`, tst.Expected, tst.Actual)
			failed = true
		}
	}

	if failed {
		t.Fail()
	}

}
