# Dough
A _really_ simple money library, for when you don't care about currencies, string formatting, or units. The aim of Dough is to make handling money more explicit, and a little safer, while trying to avoid making your code look awkward. 

### Assumptions
This package makes a few assumptions to keep things simple:

- You're dealing with sub-units (cents, pennies, etc.)
- You won't have fractions of a sub-unit. If you do, you'll round and allocate (more on that below).
- You won't have mixed currencies floating around in your code. If you do, you should take steps to ensure that, for example, you can't add USD to GBP.
- You don't care about producing a user-friendly view of the amount, e.g. "£9.99". If you do, I'm sure there are other packages for that.

## Implementation
A `Money` is just an `int`. The type lets you make sure you're not accidentally mixing up monetary amounts with other numeric types. It's effectively Hungarian notation, but enforced by the compiler.

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

