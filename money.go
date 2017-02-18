/*
Package dough provides arithmetic for monetary amounts. The aim of the package is to make handling money more explicit, and a little safer, while trying to avoid making your code look awkward.

Assumptions

This package makes a few assumptions to keep things simple. Firstly, you're dealing with sub-units (cents, pennies, etc.), and you won't have fractions of a sub-unit. You won't have mixed currencies floating around in your code. If you do, you should take steps to ensure that, for example, you can't add USD to GBP. You don't care about producing a user-friendly view of the amount, e.g. "£9.99". If you do, I'm sure there are other packages for that.

Implementation

Since we're dealing with integers, the usual arithmetic is straightforward. `+`, `-` and `==` work as you'd expect. Multiplication is an odd one. Yes, you can multiply two monetary values, which is a little odd, but it's fine. If you want to multiply a `Money` by an `int` factor, use `Money.Multiply()`:

	var price Money = 100
	var quantity int = 3
	total := price.Multiply(quantity)

The main reason for this package's existence is to provide a method for division of monetary amounts. `Money.Share()` takes a slice of `uint` ratios and returns a slice of `Money`, shared proportionately to match the ratios. The package ensures that the total of the allocations exactly matches the original value, with spare pennies distributed among allocations, favouring the first.

For example, to calculate a 10% discount on a price, you can do this:

	var price Money = 99
	shares := price.Share([]unit{90,10})
	discounted := shares[0]
	saving := shares[1]

Percentage discounting is such a common use of `Share()`, that there's also a `PercentageDiscount()` method, which takes an `int` 0 <= p <=100.

	var price Money = 99
	discounted, error := price.PercentageDiscount(10)
*/
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
