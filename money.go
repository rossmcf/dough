// Package dough provides arithmetic for monetary amounts.
package dough

import (
	"fmt"
	"math"
)

// Money is a value object representing a monetary amount.
type Money int

// Multiply returns a new Money, with the value of the original multiplied by the factor.
func (x Money) Multiply(f int) Money {
	return x * Money(f)
}

// Share allocates portions of a Money's value between parties based on weightings given.
// Spare pennies are distributed among parties evenly, from first to last.
func (x Money) Share(weightings []uint) []Money {
	n := len(weightings)
	var sum uint
	for _, w := range weightings {
		sum += w
	}
	if sum == 0 {
		for i := range weightings {
			weightings[i] = 1
		}
		sum = uint(n)
	}
	ratios := make([]float64, n)
	for i := range weightings {
		ratios[i] = float64(weightings[i]) / float64(sum)
	}

	allocations := make([]Money, n)
	fa := float64(x)
	rem := x
	for i := range ratios {
		a := Money(math.Trunc(ratios[i] * fa))
		allocations[i] = a
		rem -= a
	}
	var d Money = 1
	if rem < 0 {
		d = -1
	}
	for i := 0; rem != 0; i++ {
		ind := i % n
		if weightings[ind] == 0 {
			continue
		}
		allocations[ind] += Money(d)
		rem += (-d)
	}

	// Double-check allocation to make sure we haven't made or lost pennies.
	// It would be _very_ bad to get this wrong.
	var total Money
	for i := range allocations {
		total += allocations[i]
	}
	if total != x {
		panic(fmt.Sprintf("dough package: bad allocation. Started with %d atoms, allocated %d as %v. Weightings=%v", x, total, allocations, weightings))
	}

	return allocations
}

// PercentageDiscount discounts a Money by the given percentage, returning the discounted amount.
func (x Money) PercentageDiscount(p uint) (Money, error) {
	if p < 0 || p > 100 {
		return x, fmt.Errorf("Percentage must be ≥0, ≤100. %d given.", p)
	}
	all := x.Share([]uint{100 - p, p})
	return all[0], nil
}
