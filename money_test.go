package dough

import "testing"

// This package was largely taken from itsoneiota/dough-go.
// dough-go has a ton of tests for addition, subtraction, multiplication, and comparison.
// Since we've switch to an int, rather than a struct, there's not much point testing go's arithmetic.

func TestCanAllocate(t *testing.T) {
	var cases = []struct {
		a      Money
		ratios []uint
		want   []Money
	}{
		{0, []uint{1, 1, 1}, []Money{0, 0, 0}},
		{1, []uint{1, 1, 1}, []Money{1, 0, 0}},
		{2, []uint{1, 1, 1}, []Money{1, 1, 0}},
		{3, []uint{1, 1, 1}, []Money{1, 1, 1}},
		{4, []uint{1, 1, 1}, []Money{2, 1, 1}},
		{5, []uint{1, 1, 1}, []Money{2, 2, 1}},
		{100, []uint{0, 1, 0}, []Money{0, 100, 0}},
		{3, []uint{0, 5, 0}, []Money{0, 3, 0}},
		{300, []uint{1, 1, 1}, []Money{100, 100, 100}},
		{100, []uint{1, 1, 1}, []Money{34, 33, 33}},
		{3, []uint{0, 5, 0}, []Money{0, 3, 0}},
		{3, []uint{0, 4, 2}, []Money{0, 2, 1}},

		// Allocate spare pennies and skip zero weighting.
		{7, []uint{0, 1, 1}, []Money{0, 4, 3}},

		// Copied from MoneyTest.php
		{105, []uint{3, 7}, []Money{32, 73}},
		{5, []uint{1, 1}, []Money{3, 2}},
		{30000, []uint{122, 878}, []Money{3660, 26340}},
		{30000, []uint{122, 0, 878}, []Money{3660, 0, 26340}},
		{12000, []uint{20, 100}, []Money{2000, 10000}},

		// If weightings are equal, the amount will be shared.
		{30000, []uint{0}, []Money{30000}},
		{30000, []uint{0, 0, 0}, []Money{10000, 10000, 10000}},

		// Repeat all of the above with negatives.
		{-0, []uint{1, 1, 1}, []Money{0, 0, 0}},
		{-1, []uint{1, 1, 1}, []Money{-1, 0, 0}},
		{-2, []uint{1, 1, 1}, []Money{-1, -1, 0}},
		{-3, []uint{1, 1, 1}, []Money{-1, -1, -1}},
		{-4, []uint{1, 1, 1}, []Money{-2, -1, -1}},
		{-5, []uint{1, 1, 1}, []Money{-2, -2, -1}},
		{-100, []uint{0, 1, 0}, []Money{0, -100, 0}},
		{-3, []uint{0, 5, 0}, []Money{0, -3, 0}},
		{-300, []uint{1, 1, 1}, []Money{-100, -100, -100}},
		{-100, []uint{1, 1, 1}, []Money{-34, -33, -33}},
		{-3, []uint{0, 5, 0}, []Money{0, -3, 0}},
		{-3, []uint{0, 4, 2}, []Money{0, -2, -1}},
		{-105, []uint{3, 7}, []Money{-32, -73}},
		{-5, []uint{1, 1}, []Money{-3, -2}},
		{-30000, []uint{122, 878}, []Money{-3660, -26340}},
		{-30000, []uint{122, 0, 878}, []Money{-3660, 0, -26340}},
		{-12000, []uint{20, 100}, []Money{-2000, -10000}},
		{-30000, []uint{0}, []Money{-30000}},
		{-30000, []uint{0, 0, 0}, []Money{-10000, -10000, -10000}},
	}
	for ci, c := range cases {
		res := c.a.Share(c.ratios)
		if len(c.ratios) != len(res) {
			t.Errorf("Case %d. Incorrect number of allocations returned. Expected %d, got %d: %v", ci, len(c.ratios), len(res), res)
			return
		}
		for i := range c.want {
			if c.want[i] != res[i] {
				t.Errorf("Case %d: Sharing %d into (%v), portion %d: Expected %d, got %d", ci, c, c.ratios, i, c.want[i], res[i])
			}
		}
	}
}
